




package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)



type recvMsg struct {
	data []byte
	
	
	
	err error
}






type recvBuffer struct {
	c       chan recvMsg
	mu      sync.Mutex
	backlog []recvMsg
	err     error
}

func newRecvBuffer() *recvBuffer {
	b := &recvBuffer{
		c: make(chan recvMsg, 1),
	}
	return b
}

func (b *recvBuffer) put(r recvMsg) {
	b.mu.Lock()
	if b.err != nil {
		b.mu.Unlock()
		
		
		return
	}
	b.err = r.err
	if len(b.backlog) == 0 {
		select {
		case b.c <- r:
			b.mu.Unlock()
			return
		default:
		}
	}
	b.backlog = append(b.backlog, r)
	b.mu.Unlock()
}

func (b *recvBuffer) load() {
	b.mu.Lock()
	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			b.backlog[0] = recvMsg{}
			b.backlog = b.backlog[1:]
		default:
		}
	}
	b.mu.Unlock()
}





func (b *recvBuffer) get() <-chan recvMsg {
	return b.c
}




type recvBufferReader struct {
	ctx     context.Context
	ctxDone <-chan struct{} 
	recv    *recvBuffer
	last    []byte 
	err     error
}




func (r *recvBufferReader) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	n, r.err = r.read(p)
	return n, r.err
}

func (r *recvBufferReader) read(p []byte) (n int, err error) {
	if r.last != nil && len(r.last) > 0 {
		
		copied := copy(p, r.last)
		r.last = r.last[copied:]
		return copied, nil
	}
	select {
	case <-r.ctxDone:
		return 0, ContextErr(r.ctx.Err())
	case m := <-r.recv.get():
		r.recv.load()
		if m.err != nil {
			return 0, m.err
		}
		copied := copy(p, m.data)
		r.last = m.data[copied:]
		return copied, nil
	}
}

type streamState uint32

const (
	streamActive    streamState = iota
	streamWriteDone             
	streamReadDone              
	streamDone                  
)


type Stream struct {
	id           uint32
	st           ServerTransport    
	ctx          context.Context    
	cancel       context.CancelFunc 
	done         chan struct{}      
	ctxDone      <-chan struct{}    
	method       string             
	recvCompress string
	sendCompress string
	buf          *recvBuffer
	trReader     io.Reader
	fc           *inFlow
	wq           *writeQuota

	
	
	requestRead func(int)

	headerChan chan struct{} 
	headerDone uint32        

	
	hdrMu sync.Mutex
	
	
	
	
	header  metadata.MD
	trailer metadata.MD 

	noHeaders bool 

	
	headerSent uint32

	state streamState

	
	
	status *status.Status

	bytesReceived uint32 
	unprocessed   uint32 

	
	
	contentSubtype string
}


func (s *Stream) isHeaderSent() bool {
	return atomic.LoadUint32(&s.headerSent) == 1
}



func (s *Stream) updateHeaderSent() bool {
	return atomic.SwapUint32(&s.headerSent, 1) == 1
}

func (s *Stream) swapState(st streamState) streamState {
	return streamState(atomic.SwapUint32((*uint32)(&s.state), uint32(st)))
}

func (s *Stream) compareAndSwapState(oldState, newState streamState) bool {
	return atomic.CompareAndSwapUint32((*uint32)(&s.state), uint32(oldState), uint32(newState))
}

func (s *Stream) getState() streamState {
	return streamState(atomic.LoadUint32((*uint32)(&s.state)))
}

func (s *Stream) waitOnHeader() error {
	if s.headerChan == nil {
		
		
		return nil
	}
	select {
	case <-s.ctx.Done():
		return ContextErr(s.ctx.Err())
	case <-s.headerChan:
		return nil
	}
}



func (s *Stream) RecvCompress() string {
	if err := s.waitOnHeader(); err != nil {
		return ""
	}
	return s.recvCompress
}


func (s *Stream) SetSendCompress(str string) {
	s.sendCompress = str
}



func (s *Stream) Done() <-chan struct{} {
	return s.done
}








func (s *Stream) Header() (metadata.MD, error) {
	if s.headerChan == nil && s.header != nil {
		
		
		return s.header.Copy(), nil
	}
	err := s.waitOnHeader()
	
	select {
	case <-s.headerChan:
		if s.header == nil {
			return nil, nil
		}
		return s.header.Copy(), nil
	default:
	}
	return nil, err
}





func (s *Stream) TrailersOnly() (bool, error) {
	err := s.waitOnHeader()
	if err != nil {
		return false, err
	}
	
	return s.noHeaders && atomic.LoadUint32(&s.headerDone) == 1, nil
}






func (s *Stream) Trailer() metadata.MD {
	c := s.trailer.Copy()
	return c
}






func (s *Stream) ContentSubtype() string {
	return s.contentSubtype
}


func (s *Stream) Context() context.Context {
	return s.ctx
}


func (s *Stream) Method() string {
	return s.method
}




func (s *Stream) Status() *status.Status {
	return s.status
}




func (s *Stream) SetHeader(md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	if s.isHeaderSent() || s.getState() == streamDone {
		return ErrIllegalHeaderWrite
	}
	s.hdrMu.Lock()
	s.header = metadata.Join(s.header, md)
	s.hdrMu.Unlock()
	return nil
}




func (s *Stream) SendHeader(md metadata.MD) error {
	return s.st.WriteHeader(s, md)
}




func (s *Stream) SetTrailer(md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	if s.getState() == streamDone {
		return ErrIllegalHeaderWrite
	}
	s.hdrMu.Lock()
	s.trailer = metadata.Join(s.trailer, md)
	s.hdrMu.Unlock()
	return nil
}

func (s *Stream) write(m recvMsg) {
	s.buf.put(m)
}


func (s *Stream) Read(p []byte) (n int, err error) {
	
	if er := s.trReader.(*transportReader).er; er != nil {
		return 0, er
	}
	s.requestRead(len(p))
	return io.ReadFull(s.trReader, p)
}





type transportReader struct {
	reader io.Reader
	
	
	windowHandler func(int)
	er            error
}

func (t *transportReader) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)
	if err != nil {
		t.er = err
		return
	}
	t.windowHandler(n)
	return
}


func (s *Stream) BytesReceived() bool {
	return atomic.LoadUint32(&s.bytesReceived) == 1
}



func (s *Stream) Unprocessed() bool {
	return atomic.LoadUint32(&s.unprocessed) == 1
}



func (s *Stream) GoString() string {
	return fmt.Sprintf("<stream: %p, %v>", s, s.method)
}


type transportState int

const (
	reachable transportState = iota
	closing
	draining
)


type ServerConfig struct {
	MaxStreams            uint32
	AuthInfo              credentials.AuthInfo
	InTapHandle           tap.ServerInHandle
	StatsHandler          stats.Handler
	KeepaliveParams       keepalive.ServerParameters
	KeepalivePolicy       keepalive.EnforcementPolicy
	InitialWindowSize     int32
	InitialConnWindowSize int32
	WriteBufferSize       int
	ReadBufferSize        int
	ChannelzParentID      int64
	MaxHeaderListSize     *uint32
}



func NewServerTransport(protocol string, conn net.Conn, config *ServerConfig) (ServerTransport, error) {
	return newHTTP2Server(conn, config)
}


type ConnectOptions struct {
	
	UserAgent string
	
	Dialer func(context.Context, string) (net.Conn, error)
	
	FailOnNonTempDialError bool
	
	PerRPCCredentials []credentials.PerRPCCredentials
	
	
	TransportCredentials credentials.TransportCredentials
	
	
	CredsBundle credentials.Bundle
	
	KeepaliveParams keepalive.ClientParameters
	
	StatsHandler stats.Handler
	
	InitialWindowSize int32
	
	InitialConnWindowSize int32
	
	WriteBufferSize int
	
	ReadBufferSize int
	
	ChannelzParentID int64
	
	MaxHeaderListSize *uint32
}


type TargetInfo struct {
	Addr      string
	Metadata  interface{}
	Authority string
}



func NewClientTransport(connectCtx, ctx context.Context, target TargetInfo, opts ConnectOptions, onSuccess func(), onGoAway func(GoAwayReason), onClose func()) (ClientTransport, error) {
	return newHTTP2Client(connectCtx, ctx, target, opts, onSuccess, onGoAway, onClose)
}



type Options struct {
	
	
	Last bool
}


type CallHdr struct {
	
	Host string

	
	Method string

	
	
	SendCompress string

	
	Creds credentials.PerRPCCredentials

	
	
	
	
	
	
	ContentSubtype string

	PreviousAttempts int 
}



type ClientTransport interface {
	
	
	
	Close() error

	
	
	GracefulClose() error

	
	
	Write(s *Stream, hdr []byte, data []byte, opts *Options) error

	
	NewStream(ctx context.Context, callHdr *CallHdr) (*Stream, error)

	
	
	
	
	CloseStream(stream *Stream, err error)

	
	
	
	
	
	Error() <-chan struct{}

	
	
	
	GoAway() <-chan struct{}

	
	GetGoAwayReason() GoAwayReason

	
	IncrMsgSent()

	
	IncrMsgRecv()
}






type ServerTransport interface {
	
	HandleStreams(func(*Stream), func(context.Context, string) context.Context)

	
	
	WriteHeader(s *Stream, md metadata.MD) error

	
	
	Write(s *Stream, hdr []byte, data []byte, opts *Options) error

	
	
	WriteStatus(s *Stream, st *status.Status) error

	
	
	
	Close() error

	
	RemoteAddr() net.Addr

	
	Drain()

	
	IncrMsgSent()

	
	IncrMsgRecv()
}


func connectionErrorf(temp bool, e error, format string, a ...interface{}) ConnectionError {
	return ConnectionError{
		Desc: fmt.Sprintf(format, a...),
		temp: temp,
		err:  e,
	}
}



type ConnectionError struct {
	Desc string
	temp bool
	err  error
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("connection error: desc = %q", e.Desc)
}


func (e ConnectionError) Temporary() bool {
	return e.temp
}


func (e ConnectionError) Origin() error {
	
	
	if e.err == nil {
		return e
	}
	return e.err
}

var (
	
	ErrConnClosing = connectionErrorf(true, nil, "transport is closing")
	
	
	
	errStreamDrain = status.Error(codes.Unavailable, "the connection is draining")
	
	
	errStreamDone = errors.New("the stream is done")
	
	
	statusGoAway = status.New(codes.Unavailable, "the stream is rejected because server is draining the connection")
)


type GoAwayReason uint8

const (
	
	GoAwayInvalid GoAwayReason = 0
	
	GoAwayNoReason GoAwayReason = 1
	
	
	
	GoAwayTooManyPings GoAwayReason = 2
)





type channelzData struct {
	kpCount int64
	
	streamsStarted int64
	
	
	
	
	streamsSucceeded int64
	streamsFailed    int64
	
	
	
	lastStreamCreatedTime int64
	msgSent               int64
	msgRecv               int64
	lastMsgSentTime       int64
	lastMsgRecvTime       int64
}


func ContextErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	}
	return status.Errorf(codes.Internal, "Unexpected error from context packet: %v", err)
}
