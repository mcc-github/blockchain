

package grpc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

const (
	defaultServerMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultServerMaxSendMessageSize    = math.MaxInt32
)

type methodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor UnaryServerInterceptor) (interface{}, error)


type MethodDesc struct {
	MethodName string
	Handler    methodHandler
}


type ServiceDesc struct {
	ServiceName string
	
	
	HandlerType interface{}
	Methods     []MethodDesc
	Streams     []StreamDesc
	Metadata    interface{}
}



type service struct {
	server interface{} 
	md     map[string]*MethodDesc
	sd     map[string]*StreamDesc
	mdata  interface{}
}


type Server struct {
	opts options

	mu     sync.Mutex 
	lis    map[net.Listener]bool
	conns  map[io.Closer]bool
	serve  bool
	drain  bool
	cv     *sync.Cond          
	m      map[string]*service 
	events trace.EventLog

	quit               chan struct{}
	done               chan struct{}
	quitOnce           sync.Once
	doneOnce           sync.Once
	channelzRemoveOnce sync.Once
	serveWG            sync.WaitGroup 

	channelzID int64 
	czData     *channelzData
}

type options struct {
	creds                 credentials.TransportCredentials
	codec                 baseCodec
	cp                    Compressor
	dc                    Decompressor
	unaryInt              UnaryServerInterceptor
	streamInt             StreamServerInterceptor
	inTapHandle           tap.ServerInHandle
	statsHandler          stats.Handler
	maxConcurrentStreams  uint32
	maxReceiveMessageSize int
	maxSendMessageSize    int
	unknownStreamDesc     *StreamDesc
	keepaliveParams       keepalive.ServerParameters
	keepalivePolicy       keepalive.EnforcementPolicy
	initialWindowSize     int32
	initialConnWindowSize int32
	writeBufferSize       int
	readBufferSize        int
	connectionTimeout     time.Duration
	maxHeaderListSize     *uint32
}

var defaultServerOptions = options{
	maxReceiveMessageSize: defaultServerMaxReceiveMessageSize,
	maxSendMessageSize:    defaultServerMaxSendMessageSize,
	connectionTimeout:     120 * time.Second,
	writeBufferSize:       defaultWriteBufSize,
	readBufferSize:        defaultReadBufSize,
}


type ServerOption func(*options)






func WriteBufferSize(s int) ServerOption {
	return func(o *options) {
		o.writeBufferSize = s
	}
}






func ReadBufferSize(s int) ServerOption {
	return func(o *options) {
		o.readBufferSize = s
	}
}



func InitialWindowSize(s int32) ServerOption {
	return func(o *options) {
		o.initialWindowSize = s
	}
}



func InitialConnWindowSize(s int32) ServerOption {
	return func(o *options) {
		o.initialConnWindowSize = s
	}
}


func KeepaliveParams(kp keepalive.ServerParameters) ServerOption {
	return func(o *options) {
		o.keepaliveParams = kp
	}
}


func KeepaliveEnforcementPolicy(kep keepalive.EnforcementPolicy) ServerOption {
	return func(o *options) {
		o.keepalivePolicy = kep
	}
}




func CustomCodec(codec Codec) ServerOption {
	return func(o *options) {
		o.codec = codec
	}
}








func RPCCompressor(cp Compressor) ServerOption {
	return func(o *options) {
		o.cp = cp
	}
}






func RPCDecompressor(dc Decompressor) ServerOption {
	return func(o *options) {
		o.dc = dc
	}
}





func MaxMsgSize(m int) ServerOption {
	return MaxRecvMsgSize(m)
}



func MaxRecvMsgSize(m int) ServerOption {
	return func(o *options) {
		o.maxReceiveMessageSize = m
	}
}



func MaxSendMsgSize(m int) ServerOption {
	return func(o *options) {
		o.maxSendMessageSize = m
	}
}



func MaxConcurrentStreams(n uint32) ServerOption {
	return func(o *options) {
		o.maxConcurrentStreams = n
	}
}


func Creds(c credentials.TransportCredentials) ServerOption {
	return func(o *options) {
		o.creds = c
	}
}




func UnaryInterceptor(i UnaryServerInterceptor) ServerOption {
	return func(o *options) {
		if o.unaryInt != nil {
			panic("The unary server interceptor was already set and may not be reset.")
		}
		o.unaryInt = i
	}
}



func StreamInterceptor(i StreamServerInterceptor) ServerOption {
	return func(o *options) {
		if o.streamInt != nil {
			panic("The stream server interceptor was already set and may not be reset.")
		}
		o.streamInt = i
	}
}



func InTapHandle(h tap.ServerInHandle) ServerOption {
	return func(o *options) {
		if o.inTapHandle != nil {
			panic("The tap handle was already set and may not be reset.")
		}
		o.inTapHandle = h
	}
}


func StatsHandler(h stats.Handler) ServerOption {
	return func(o *options) {
		o.statsHandler = h
	}
}







func UnknownServiceHandler(streamHandler StreamHandler) ServerOption {
	return func(o *options) {
		o.unknownStreamDesc = &StreamDesc{
			StreamName: "unknown_service_handler",
			Handler:    streamHandler,
			
			ClientStreams: true,
			ServerStreams: true,
		}
	}
}







func ConnectionTimeout(d time.Duration) ServerOption {
	return func(o *options) {
		o.connectionTimeout = d
	}
}



func MaxHeaderListSize(s uint32) ServerOption {
	return func(o *options) {
		o.maxHeaderListSize = &s
	}
}



func NewServer(opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o(&opts)
	}
	s := &Server{
		lis:    make(map[net.Listener]bool),
		opts:   opts,
		conns:  make(map[io.Closer]bool),
		m:      make(map[string]*service),
		quit:   make(chan struct{}),
		done:   make(chan struct{}),
		czData: new(channelzData),
	}
	s.cv = sync.NewCond(&s.mu)
	if EnableTracing {
		_, file, line, _ := runtime.Caller(1)
		s.events = trace.NewEventLog("grpc.Server", fmt.Sprintf("%s:%d", file, line))
	}

	if channelz.IsOn() {
		s.channelzID = channelz.RegisterServer(&channelzServer{s}, "")
	}
	return s
}



func (s *Server) printf(format string, a ...interface{}) {
	if s.events != nil {
		s.events.Printf(format, a...)
	}
}



func (s *Server) errorf(format string, a ...interface{}) {
	if s.events != nil {
		s.events.Errorf(format, a...)
	}
}




func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	st := reflect.TypeOf(ss)
	if !st.Implements(ht) {
		grpclog.Fatalf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
	}
	s.register(sd, ss)
}

func (s *Server) register(sd *ServiceDesc, ss interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.printf("RegisterService(%q)", sd.ServiceName)
	if s.serve {
		grpclog.Fatalf("grpc: Server.RegisterService after Server.Serve for %q", sd.ServiceName)
	}
	if _, ok := s.m[sd.ServiceName]; ok {
		grpclog.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
	}
	srv := &service{
		server: ss,
		md:     make(map[string]*MethodDesc),
		sd:     make(map[string]*StreamDesc),
		mdata:  sd.Metadata,
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		srv.md[d.MethodName] = d
	}
	for i := range sd.Streams {
		d := &sd.Streams[i]
		srv.sd[d.StreamName] = d
	}
	s.m[sd.ServiceName] = srv
}


type MethodInfo struct {
	
	Name string
	
	IsClientStream bool
	
	IsServerStream bool
}


type ServiceInfo struct {
	Methods []MethodInfo
	
	Metadata interface{}
}



func (s *Server) GetServiceInfo() map[string]ServiceInfo {
	ret := make(map[string]ServiceInfo)
	for n, srv := range s.m {
		methods := make([]MethodInfo, 0, len(srv.md)+len(srv.sd))
		for m := range srv.md {
			methods = append(methods, MethodInfo{
				Name:           m,
				IsClientStream: false,
				IsServerStream: false,
			})
		}
		for m, d := range srv.sd {
			methods = append(methods, MethodInfo{
				Name:           m,
				IsClientStream: d.ClientStreams,
				IsServerStream: d.ServerStreams,
			})
		}

		ret[n] = ServiceInfo{
			Methods:  methods,
			Metadata: srv.mdata,
		}
	}
	return ret
}



var ErrServerStopped = errors.New("grpc: the server has been stopped")

func (s *Server) useTransportAuthenticator(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	if s.opts.creds == nil {
		return rawConn, nil, nil
	}
	return s.opts.creds.ServerHandshake(rawConn)
}

type listenSocket struct {
	net.Listener
	channelzID int64
}

func (l *listenSocket) ChannelzMetric() *channelz.SocketInternalMetric {
	return &channelz.SocketInternalMetric{
		SocketOptions: channelz.GetSocketOption(l.Listener),
		LocalAddr:     l.Listener.Addr(),
	}
}

func (l *listenSocket) Close() error {
	err := l.Listener.Close()
	if channelz.IsOn() {
		channelz.RemoveEntry(l.channelzID)
	}
	return err
}







func (s *Server) Serve(lis net.Listener) error {
	s.mu.Lock()
	s.printf("serving")
	s.serve = true
	if s.lis == nil {
		
		s.mu.Unlock()
		lis.Close()
		return ErrServerStopped
	}

	s.serveWG.Add(1)
	defer func() {
		s.serveWG.Done()
		select {
		
		case <-s.quit:
			<-s.done
		default:
		}
	}()

	ls := &listenSocket{Listener: lis}
	s.lis[ls] = true

	if channelz.IsOn() {
		ls.channelzID = channelz.RegisterListenSocket(ls, s.channelzID, "")
	}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		if s.lis != nil && s.lis[ls] {
			ls.Close()
			delete(s.lis, ls)
		}
		s.mu.Unlock()
	}()

	var tempDelay time.Duration 

	for {
		rawConn, err := lis.Accept()
		if err != nil {
			if ne, ok := err.(interface {
				Temporary() bool
			}); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				s.mu.Lock()
				s.printf("Accept error: %v; retrying in %v", err, tempDelay)
				s.mu.Unlock()
				timer := time.NewTimer(tempDelay)
				select {
				case <-timer.C:
				case <-s.quit:
					timer.Stop()
					return nil
				}
				continue
			}
			s.mu.Lock()
			s.printf("done serving; Accept = %v", err)
			s.mu.Unlock()

			select {
			case <-s.quit:
				return nil
			default:
			}
			return err
		}
		tempDelay = 0
		
		
		
		
		
		s.serveWG.Add(1)
		go func() {
			s.handleRawConn(rawConn)
			s.serveWG.Done()
		}()
	}
}



func (s *Server) handleRawConn(rawConn net.Conn) {
	rawConn.SetDeadline(time.Now().Add(s.opts.connectionTimeout))
	conn, authInfo, err := s.useTransportAuthenticator(rawConn)
	if err != nil {
		s.mu.Lock()
		s.errorf("ServerHandshake(%q) failed: %v", rawConn.RemoteAddr(), err)
		s.mu.Unlock()
		grpclog.Warningf("grpc: Server.Serve failed to complete security handshake from %q: %v", rawConn.RemoteAddr(), err)
		
		if err != credentials.ErrConnDispatched {
			rawConn.Close()
		}
		rawConn.SetDeadline(time.Time{})
		return
	}

	s.mu.Lock()
	if s.conns == nil {
		s.mu.Unlock()
		conn.Close()
		return
	}
	s.mu.Unlock()

	
	st := s.newHTTP2Transport(conn, authInfo)
	if st == nil {
		return
	}

	rawConn.SetDeadline(time.Time{})
	if !s.addConn(st) {
		return
	}
	go func() {
		s.serveStreams(st)
		s.removeConn(st)
	}()
}



func (s *Server) newHTTP2Transport(c net.Conn, authInfo credentials.AuthInfo) transport.ServerTransport {
	config := &transport.ServerConfig{
		MaxStreams:            s.opts.maxConcurrentStreams,
		AuthInfo:              authInfo,
		InTapHandle:           s.opts.inTapHandle,
		StatsHandler:          s.opts.statsHandler,
		KeepaliveParams:       s.opts.keepaliveParams,
		KeepalivePolicy:       s.opts.keepalivePolicy,
		InitialWindowSize:     s.opts.initialWindowSize,
		InitialConnWindowSize: s.opts.initialConnWindowSize,
		WriteBufferSize:       s.opts.writeBufferSize,
		ReadBufferSize:        s.opts.readBufferSize,
		ChannelzParentID:      s.channelzID,
		MaxHeaderListSize:     s.opts.maxHeaderListSize,
	}
	st, err := transport.NewServerTransport("http2", c, config)
	if err != nil {
		s.mu.Lock()
		s.errorf("NewServerTransport(%q) failed: %v", c.RemoteAddr(), err)
		s.mu.Unlock()
		c.Close()
		grpclog.Warningln("grpc: Server.Serve failed to create ServerTransport: ", err)
		return nil
	}

	return st
}

func (s *Server) serveStreams(st transport.ServerTransport) {
	defer st.Close()
	var wg sync.WaitGroup
	st.HandleStreams(func(stream *transport.Stream) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.handleStream(st, stream, s.traceInfo(st, stream))
		}()
	}, func(ctx context.Context, method string) context.Context {
		if !EnableTracing {
			return ctx
		}
		tr := trace.New("grpc.Recv."+methodFamily(method), method)
		return trace.NewContext(ctx, tr)
	})
	wg.Wait()
}

var _ http.Handler = (*Server)(nil)

























func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	st, err := transport.NewServerHandlerTransport(w, r, s.opts.statsHandler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !s.addConn(st) {
		return
	}
	defer s.removeConn(st)
	s.serveStreams(st)
}



func (s *Server) traceInfo(st transport.ServerTransport, stream *transport.Stream) (trInfo *traceInfo) {
	tr, ok := trace.FromContext(stream.Context())
	if !ok {
		return nil
	}

	trInfo = &traceInfo{
		tr: tr,
	}
	trInfo.firstLine.client = false
	trInfo.firstLine.remoteAddr = st.RemoteAddr()

	if dl, ok := stream.Context().Deadline(); ok {
		trInfo.firstLine.deadline = dl.Sub(time.Now())
	}
	return trInfo
}

func (s *Server) addConn(c io.Closer) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil {
		c.Close()
		return false
	}
	if s.drain {
		
		
		c.(transport.ServerTransport).Drain()
	}
	s.conns[c] = true
	return true
}

func (s *Server) removeConn(c io.Closer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns != nil {
		delete(s.conns, c)
		s.cv.Broadcast()
	}
}

func (s *Server) channelzMetric() *channelz.ServerInternalMetric {
	return &channelz.ServerInternalMetric{
		CallsStarted:             atomic.LoadInt64(&s.czData.callsStarted),
		CallsSucceeded:           atomic.LoadInt64(&s.czData.callsSucceeded),
		CallsFailed:              atomic.LoadInt64(&s.czData.callsFailed),
		LastCallStartedTimestamp: time.Unix(0, atomic.LoadInt64(&s.czData.lastCallStartedTime)),
	}
}

func (s *Server) incrCallsStarted() {
	atomic.AddInt64(&s.czData.callsStarted, 1)
	atomic.StoreInt64(&s.czData.lastCallStartedTime, time.Now().UnixNano())
}

func (s *Server) incrCallsSucceeded() {
	atomic.AddInt64(&s.czData.callsSucceeded, 1)
}

func (s *Server) incrCallsFailed() {
	atomic.AddInt64(&s.czData.callsFailed, 1)
}

func (s *Server) sendResponse(t transport.ServerTransport, stream *transport.Stream, msg interface{}, cp Compressor, opts *transport.Options, comp encoding.Compressor) error {
	data, err := encode(s.getCodec(stream.ContentSubtype()), msg)
	if err != nil {
		grpclog.Errorln("grpc: server failed to encode response: ", err)
		return err
	}
	compData, err := compress(data, cp, comp)
	if err != nil {
		grpclog.Errorln("grpc: server failed to compress response: ", err)
		return err
	}
	hdr, payload := msgHeader(data, compData)
	
	if len(payload) > s.opts.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "grpc: trying to send message larger than max (%d vs. %d)", len(payload), s.opts.maxSendMessageSize)
	}
	err = t.Write(stream, hdr, payload, opts)
	if err == nil && s.opts.statsHandler != nil {
		s.opts.statsHandler.HandleRPC(stream.Context(), outPayload(false, msg, data, payload, time.Now()))
	}
	return err
}

func (s *Server) processUnaryRPC(t transport.ServerTransport, stream *transport.Stream, srv *service, md *MethodDesc, trInfo *traceInfo) (err error) {
	if channelz.IsOn() {
		s.incrCallsStarted()
		defer func() {
			if err != nil && err != io.EOF {
				s.incrCallsFailed()
			} else {
				s.incrCallsSucceeded()
			}
		}()
	}
	sh := s.opts.statsHandler
	if sh != nil {
		beginTime := time.Now()
		begin := &stats.Begin{
			BeginTime: beginTime,
		}
		sh.HandleRPC(stream.Context(), begin)
		defer func() {
			end := &stats.End{
				BeginTime: beginTime,
				EndTime:   time.Now(),
			}
			if err != nil && err != io.EOF {
				end.Error = toRPCErr(err)
			}
			sh.HandleRPC(stream.Context(), end)
		}()
	}
	if trInfo != nil {
		defer trInfo.tr.Finish()
		trInfo.firstLine.client = false
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
		defer func() {
			if err != nil && err != io.EOF {
				trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
				trInfo.tr.SetError()
			}
		}()
	}

	
	
	
	
	var comp, decomp encoding.Compressor
	var cp Compressor
	var dc Decompressor

	
	
	if rc := stream.RecvCompress(); s.opts.dc != nil && s.opts.dc.Type() == rc {
		dc = s.opts.dc
	} else if rc != "" && rc != encoding.Identity {
		decomp = encoding.GetCompressor(rc)
		if decomp == nil {
			st := status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", rc)
			t.WriteStatus(stream, st)
			return st.Err()
		}
	}

	
	
	
	
	if s.opts.cp != nil {
		cp = s.opts.cp
		stream.SetSendCompress(cp.Type())
	} else if rc := stream.RecvCompress(); rc != "" && rc != encoding.Identity {
		
		comp = encoding.GetCompressor(rc)
		if comp != nil {
			stream.SetSendCompress(rc)
		}
	}

	p := &parser{r: stream}
	pf, req, err := p.recvMsg(s.opts.maxReceiveMessageSize)
	if err == io.EOF {
		
		return err
	}
	if err == io.ErrUnexpectedEOF {
		err = status.Errorf(codes.Internal, io.ErrUnexpectedEOF.Error())
	}
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if e := t.WriteStatus(stream, st); e != nil {
				grpclog.Warningf("grpc: Server.processUnaryRPC failed to write status %v", e)
			}
		} else {
			switch st := err.(type) {
			case transport.ConnectionError:
				
			default:
				panic(fmt.Sprintf("grpc: Unexpected error (%T) from recvMsg: %v", st, st))
			}
		}
		return err
	}
	if channelz.IsOn() {
		t.IncrMsgRecv()
	}
	if st := checkRecvPayload(pf, stream.RecvCompress(), dc != nil || decomp != nil); st != nil {
		if e := t.WriteStatus(stream, st); e != nil {
			grpclog.Warningf("grpc: Server.processUnaryRPC failed to write status %v", e)
		}
		return st.Err()
	}
	var inPayload *stats.InPayload
	if sh != nil {
		inPayload = &stats.InPayload{
			RecvTime: time.Now(),
		}
	}
	df := func(v interface{}) error {
		if inPayload != nil {
			inPayload.WireLength = len(req)
		}
		if pf == compressionMade {
			var err error
			if dc != nil {
				req, err = dc.Do(bytes.NewReader(req))
				if err != nil {
					return status.Errorf(codes.Internal, err.Error())
				}
			} else {
				tmp, _ := decomp.Decompress(bytes.NewReader(req))
				req, err = ioutil.ReadAll(tmp)
				if err != nil {
					return status.Errorf(codes.Internal, "grpc: failed to decompress the received message %v", err)
				}
			}
		}
		if len(req) > s.opts.maxReceiveMessageSize {
			
			
			return status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max (%d vs. %d)", len(req), s.opts.maxReceiveMessageSize)
		}
		if err := s.getCodec(stream.ContentSubtype()).Unmarshal(req, v); err != nil {
			return status.Errorf(codes.Internal, "grpc: error unmarshalling request: %v", err)
		}
		if inPayload != nil {
			inPayload.Payload = v
			inPayload.Data = req
			inPayload.Length = len(req)
			sh.HandleRPC(stream.Context(), inPayload)
		}
		if trInfo != nil {
			trInfo.tr.LazyLog(&payload{sent: false, msg: v}, true)
		}
		return nil
	}
	ctx := NewContextWithServerTransportStream(stream.Context(), stream)
	reply, appErr := md.Handler(srv.server, ctx, df, s.opts.unaryInt)
	if appErr != nil {
		appStatus, ok := status.FromError(appErr)
		if !ok {
			
			appErr = status.Error(codes.Unknown, appErr.Error())
			appStatus, _ = status.FromError(appErr)
		}
		if trInfo != nil {
			trInfo.tr.LazyLog(stringer(appStatus.Message()), true)
			trInfo.tr.SetError()
		}
		if e := t.WriteStatus(stream, appStatus); e != nil {
			grpclog.Warningf("grpc: Server.processUnaryRPC failed to write status: %v", e)
		}
		return appErr
	}
	if trInfo != nil {
		trInfo.tr.LazyLog(stringer("OK"), false)
	}
	opts := &transport.Options{Last: true}

	if err := s.sendResponse(t, stream, reply, cp, opts, comp); err != nil {
		if err == io.EOF {
			
			return err
		}
		if s, ok := status.FromError(err); ok {
			if e := t.WriteStatus(stream, s); e != nil {
				grpclog.Warningf("grpc: Server.processUnaryRPC failed to write status: %v", e)
			}
		} else {
			switch st := err.(type) {
			case transport.ConnectionError:
				
			default:
				panic(fmt.Sprintf("grpc: Unexpected error (%T) from sendResponse: %v", st, st))
			}
		}
		return err
	}
	if channelz.IsOn() {
		t.IncrMsgSent()
	}
	if trInfo != nil {
		trInfo.tr.LazyLog(&payload{sent: true, msg: reply}, true)
	}
	
	
	
	return t.WriteStatus(stream, status.New(codes.OK, ""))
}

func (s *Server) processStreamingRPC(t transport.ServerTransport, stream *transport.Stream, srv *service, sd *StreamDesc, trInfo *traceInfo) (err error) {
	if channelz.IsOn() {
		s.incrCallsStarted()
		defer func() {
			if err != nil && err != io.EOF {
				s.incrCallsFailed()
			} else {
				s.incrCallsSucceeded()
			}
		}()
	}
	sh := s.opts.statsHandler
	if sh != nil {
		beginTime := time.Now()
		begin := &stats.Begin{
			BeginTime: beginTime,
		}
		sh.HandleRPC(stream.Context(), begin)
		defer func() {
			end := &stats.End{
				BeginTime: beginTime,
				EndTime:   time.Now(),
			}
			if err != nil && err != io.EOF {
				end.Error = toRPCErr(err)
			}
			sh.HandleRPC(stream.Context(), end)
		}()
	}
	ctx := NewContextWithServerTransportStream(stream.Context(), stream)
	ss := &serverStream{
		ctx:                   ctx,
		t:                     t,
		s:                     stream,
		p:                     &parser{r: stream},
		codec:                 s.getCodec(stream.ContentSubtype()),
		maxReceiveMessageSize: s.opts.maxReceiveMessageSize,
		maxSendMessageSize:    s.opts.maxSendMessageSize,
		trInfo:                trInfo,
		statsHandler:          sh,
	}

	
	
	if rc := stream.RecvCompress(); s.opts.dc != nil && s.opts.dc.Type() == rc {
		ss.dc = s.opts.dc
	} else if rc != "" && rc != encoding.Identity {
		ss.decomp = encoding.GetCompressor(rc)
		if ss.decomp == nil {
			st := status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", rc)
			t.WriteStatus(ss.s, st)
			return st.Err()
		}
	}

	
	
	
	
	if s.opts.cp != nil {
		ss.cp = s.opts.cp
		stream.SetSendCompress(s.opts.cp.Type())
	} else if rc := stream.RecvCompress(); rc != "" && rc != encoding.Identity {
		
		ss.comp = encoding.GetCompressor(rc)
		if ss.comp != nil {
			stream.SetSendCompress(rc)
		}
	}

	if trInfo != nil {
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
		defer func() {
			ss.mu.Lock()
			if err != nil && err != io.EOF {
				ss.trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
				ss.trInfo.tr.SetError()
			}
			ss.trInfo.tr.Finish()
			ss.trInfo.tr = nil
			ss.mu.Unlock()
		}()
	}
	var appErr error
	var server interface{}
	if srv != nil {
		server = srv.server
	}
	if s.opts.streamInt == nil {
		appErr = sd.Handler(server, ss)
	} else {
		info := &StreamServerInfo{
			FullMethod:     stream.Method(),
			IsClientStream: sd.ClientStreams,
			IsServerStream: sd.ServerStreams,
		}
		appErr = s.opts.streamInt(server, ss, info, sd.Handler)
	}
	if appErr != nil {
		appStatus, ok := status.FromError(appErr)
		if !ok {
			appStatus = status.New(codes.Unknown, appErr.Error())
			appErr = appStatus.Err()
		}
		if trInfo != nil {
			ss.mu.Lock()
			ss.trInfo.tr.LazyLog(stringer(appStatus.Message()), true)
			ss.trInfo.tr.SetError()
			ss.mu.Unlock()
		}
		t.WriteStatus(ss.s, appStatus)
		
		return appErr
	}
	if trInfo != nil {
		ss.mu.Lock()
		ss.trInfo.tr.LazyLog(stringer("OK"), false)
		ss.mu.Unlock()
	}
	return t.WriteStatus(ss.s, status.New(codes.OK, ""))
}

func (s *Server) handleStream(t transport.ServerTransport, stream *transport.Stream, trInfo *traceInfo) {
	sm := stream.Method()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		if trInfo != nil {
			trInfo.tr.LazyLog(&fmtStringer{"Malformed method name %q", []interface{}{sm}}, true)
			trInfo.tr.SetError()
		}
		errDesc := fmt.Sprintf("malformed method name: %q", stream.Method())
		if err := t.WriteStatus(stream, status.New(codes.ResourceExhausted, errDesc)); err != nil {
			if trInfo != nil {
				trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
				trInfo.tr.SetError()
			}
			grpclog.Warningf("grpc: Server.handleStream failed to write status: %v", err)
		}
		if trInfo != nil {
			trInfo.tr.Finish()
		}
		return
	}
	service := sm[:pos]
	method := sm[pos+1:]
	srv, ok := s.m[service]
	if !ok {
		if unknownDesc := s.opts.unknownStreamDesc; unknownDesc != nil {
			s.processStreamingRPC(t, stream, nil, unknownDesc, trInfo)
			return
		}
		if trInfo != nil {
			trInfo.tr.LazyLog(&fmtStringer{"Unknown service %v", []interface{}{service}}, true)
			trInfo.tr.SetError()
		}
		errDesc := fmt.Sprintf("unknown service %v", service)
		if err := t.WriteStatus(stream, status.New(codes.Unimplemented, errDesc)); err != nil {
			if trInfo != nil {
				trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
				trInfo.tr.SetError()
			}
			grpclog.Warningf("grpc: Server.handleStream failed to write status: %v", err)
		}
		if trInfo != nil {
			trInfo.tr.Finish()
		}
		return
	}
	
	if md, ok := srv.md[method]; ok {
		s.processUnaryRPC(t, stream, srv, md, trInfo)
		return
	}
	if sd, ok := srv.sd[method]; ok {
		s.processStreamingRPC(t, stream, srv, sd, trInfo)
		return
	}
	if trInfo != nil {
		trInfo.tr.LazyLog(&fmtStringer{"Unknown method %v", []interface{}{method}}, true)
		trInfo.tr.SetError()
	}
	if unknownDesc := s.opts.unknownStreamDesc; unknownDesc != nil {
		s.processStreamingRPC(t, stream, nil, unknownDesc, trInfo)
		return
	}
	errDesc := fmt.Sprintf("unknown method %v", method)
	if err := t.WriteStatus(stream, status.New(codes.Unimplemented, errDesc)); err != nil {
		if trInfo != nil {
			trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
			trInfo.tr.SetError()
		}
		grpclog.Warningf("grpc: Server.handleStream failed to write status: %v", err)
	}
	if trInfo != nil {
		trInfo.tr.Finish()
	}
}


type streamKey struct{}





func NewContextWithServerTransportStream(ctx context.Context, stream ServerTransportStream) context.Context {
	return context.WithValue(ctx, streamKey{}, stream)
}









type ServerTransportStream interface {
	Method() string
	SetHeader(md metadata.MD) error
	SendHeader(md metadata.MD) error
	SetTrailer(md metadata.MD) error
}






func ServerTransportStreamFromContext(ctx context.Context) ServerTransportStream {
	s, _ := ctx.Value(streamKey{}).(ServerTransportStream)
	return s
}






func (s *Server) Stop() {
	s.quitOnce.Do(func() {
		close(s.quit)
	})

	defer func() {
		s.serveWG.Wait()
		s.doneOnce.Do(func() {
			close(s.done)
		})
	}()

	s.channelzRemoveOnce.Do(func() {
		if channelz.IsOn() {
			channelz.RemoveEntry(s.channelzID)
		}
	})

	s.mu.Lock()
	listeners := s.lis
	s.lis = nil
	st := s.conns
	s.conns = nil
	
	s.cv.Broadcast()
	s.mu.Unlock()

	for lis := range listeners {
		lis.Close()
	}
	for c := range st {
		c.Close()
	}

	s.mu.Lock()
	if s.events != nil {
		s.events.Finish()
		s.events = nil
	}
	s.mu.Unlock()
}




func (s *Server) GracefulStop() {
	s.quitOnce.Do(func() {
		close(s.quit)
	})

	defer func() {
		s.doneOnce.Do(func() {
			close(s.done)
		})
	}()

	s.channelzRemoveOnce.Do(func() {
		if channelz.IsOn() {
			channelz.RemoveEntry(s.channelzID)
		}
	})
	s.mu.Lock()
	if s.conns == nil {
		s.mu.Unlock()
		return
	}

	for lis := range s.lis {
		lis.Close()
	}
	s.lis = nil
	if !s.drain {
		for c := range s.conns {
			c.(transport.ServerTransport).Drain()
		}
		s.drain = true
	}

	
	
	s.mu.Unlock()
	s.serveWG.Wait()
	s.mu.Lock()

	for len(s.conns) != 0 {
		s.cv.Wait()
	}
	s.conns = nil
	if s.events != nil {
		s.events.Finish()
		s.events = nil
	}
	s.mu.Unlock()
}



func (s *Server) getCodec(contentSubtype string) baseCodec {
	if s.opts.codec != nil {
		return s.opts.codec
	}
	if contentSubtype == "" {
		return encoding.GetCodec(proto.Name)
	}
	codec := encoding.GetCodec(contentSubtype)
	if codec == nil {
		return encoding.GetCodec(proto.Name)
	}
	return codec
}







func SetHeader(ctx context.Context, md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil {
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	}
	return stream.SetHeader(md)
}



func SendHeader(ctx context.Context, md metadata.MD) error {
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil {
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	}
	if err := stream.SendHeader(md); err != nil {
		return toRPCErr(err)
	}
	return nil
}



func SetTrailer(ctx context.Context, md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	stream := ServerTransportStreamFromContext(ctx)
	if stream == nil {
		return status.Errorf(codes.Internal, "grpc: failed to fetch the stream from the context %v", ctx)
	}
	return stream.SetTrailer(md)
}



func Method(ctx context.Context) (string, bool) {
	s := ServerTransportStreamFromContext(ctx)
	if s == nil {
		return "", false
	}
	return s.Method(), true
}

type channelzServer struct {
	s *Server
}

func (c *channelzServer) ChannelzMetric() *channelz.ServerInternalMetric {
	return c.s.channelzMetric()
}
