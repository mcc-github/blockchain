

package grpc

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)




type Compressor interface {
	
	Do(w io.Writer, p []byte) error
	
	Type() string
}

type gzipCompressor struct {
	pool sync.Pool
}




func NewGZIPCompressor() Compressor {
	c, _ := NewGZIPCompressorWithLevel(gzip.DefaultCompression)
	return c
}







func NewGZIPCompressorWithLevel(level int) (Compressor, error) {
	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		return nil, fmt.Errorf("grpc: invalid compression level: %d", level)
	}
	return &gzipCompressor{
		pool: sync.Pool{
			New: func() interface{} {
				w, err := gzip.NewWriterLevel(ioutil.Discard, level)
				if err != nil {
					panic(err)
				}
				return w
			},
		},
	}, nil
}

func (c *gzipCompressor) Do(w io.Writer, p []byte) error {
	z := c.pool.Get().(*gzip.Writer)
	defer c.pool.Put(z)
	z.Reset(w)
	if _, err := z.Write(p); err != nil {
		return err
	}
	return z.Close()
}

func (c *gzipCompressor) Type() string {
	return "gzip"
}




type Decompressor interface {
	
	Do(r io.Reader) ([]byte, error)
	
	Type() string
}

type gzipDecompressor struct {
	pool sync.Pool
}




func NewGZIPDecompressor() Decompressor {
	return &gzipDecompressor{}
}

func (d *gzipDecompressor) Do(r io.Reader) ([]byte, error) {
	var z *gzip.Reader
	switch maybeZ := d.pool.Get().(type) {
	case nil:
		newZ, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		z = newZ
	case *gzip.Reader:
		z = maybeZ
		if err := z.Reset(r); err != nil {
			d.pool.Put(z)
			return nil, err
		}
	}

	defer func() {
		z.Close()
		d.pool.Put(z)
	}()
	return ioutil.ReadAll(z)
}

func (d *gzipDecompressor) Type() string {
	return "gzip"
}


type callInfo struct {
	compressorType        string
	failFast              bool
	stream                *clientStream
	maxReceiveMessageSize *int
	maxSendMessageSize    *int
	creds                 credentials.PerRPCCredentials
	contentSubtype        string
	codec                 baseCodec
	maxRetryRPCBufferSize int
}

func defaultCallInfo() *callInfo {
	return &callInfo{
		failFast:              true,
		maxRetryRPCBufferSize: 256 * 1024, 
	}
}



type CallOption interface {
	
	
	before(*callInfo) error

	
	
	after(*callInfo)
}




type EmptyCallOption struct{}

func (EmptyCallOption) before(*callInfo) error { return nil }
func (EmptyCallOption) after(*callInfo)        {}



func Header(md *metadata.MD) CallOption {
	return HeaderCallOption{HeaderAddr: md}
}




type HeaderCallOption struct {
	HeaderAddr *metadata.MD
}

func (o HeaderCallOption) before(c *callInfo) error { return nil }
func (o HeaderCallOption) after(c *callInfo) {
	if c.stream != nil {
		*o.HeaderAddr, _ = c.stream.Header()
	}
}



func Trailer(md *metadata.MD) CallOption {
	return TrailerCallOption{TrailerAddr: md}
}




type TrailerCallOption struct {
	TrailerAddr *metadata.MD
}

func (o TrailerCallOption) before(c *callInfo) error { return nil }
func (o TrailerCallOption) after(c *callInfo) {
	if c.stream != nil {
		*o.TrailerAddr = c.stream.Trailer()
	}
}



func Peer(p *peer.Peer) CallOption {
	return PeerCallOption{PeerAddr: p}
}




type PeerCallOption struct {
	PeerAddr *peer.Peer
}

func (o PeerCallOption) before(c *callInfo) error { return nil }
func (o PeerCallOption) after(c *callInfo) {
	if c.stream != nil {
		if x, ok := peer.FromContext(c.stream.Context()); ok {
			*o.PeerAddr = *x
		}
	}
}











func FailFast(failFast bool) CallOption {
	return FailFastCallOption{FailFast: failFast}
}




type FailFastCallOption struct {
	FailFast bool
}

func (o FailFastCallOption) before(c *callInfo) error {
	c.failFast = o.FailFast
	return nil
}
func (o FailFastCallOption) after(c *callInfo) {}


func MaxCallRecvMsgSize(s int) CallOption {
	return MaxRecvMsgSizeCallOption{MaxRecvMsgSize: s}
}




type MaxRecvMsgSizeCallOption struct {
	MaxRecvMsgSize int
}

func (o MaxRecvMsgSizeCallOption) before(c *callInfo) error {
	c.maxReceiveMessageSize = &o.MaxRecvMsgSize
	return nil
}
func (o MaxRecvMsgSizeCallOption) after(c *callInfo) {}


func MaxCallSendMsgSize(s int) CallOption {
	return MaxSendMsgSizeCallOption{MaxSendMsgSize: s}
}




type MaxSendMsgSizeCallOption struct {
	MaxSendMsgSize int
}

func (o MaxSendMsgSizeCallOption) before(c *callInfo) error {
	c.maxSendMessageSize = &o.MaxSendMsgSize
	return nil
}
func (o MaxSendMsgSizeCallOption) after(c *callInfo) {}



func PerRPCCredentials(creds credentials.PerRPCCredentials) CallOption {
	return PerRPCCredsCallOption{Creds: creds}
}




type PerRPCCredsCallOption struct {
	Creds credentials.PerRPCCredentials
}

func (o PerRPCCredsCallOption) before(c *callInfo) error {
	c.creds = o.Creds
	return nil
}
func (o PerRPCCredsCallOption) after(c *callInfo) {}






func UseCompressor(name string) CallOption {
	return CompressorCallOption{CompressorType: name}
}



type CompressorCallOption struct {
	CompressorType string
}

func (o CompressorCallOption) before(c *callInfo) error {
	c.compressorType = o.CompressorType
	return nil
}
func (o CompressorCallOption) after(c *callInfo) {}

















func CallContentSubtype(contentSubtype string) CallOption {
	return ContentSubtypeCallOption{ContentSubtype: strings.ToLower(contentSubtype)}
}




type ContentSubtypeCallOption struct {
	ContentSubtype string
}

func (o ContentSubtypeCallOption) before(c *callInfo) error {
	c.contentSubtype = o.ContentSubtype
	return nil
}
func (o ContentSubtypeCallOption) after(c *callInfo) {}













func CallCustomCodec(codec Codec) CallOption {
	return CustomCodecCallOption{Codec: codec}
}




type CustomCodecCallOption struct {
	Codec Codec
}

func (o CustomCodecCallOption) before(c *callInfo) error {
	c.codec = o.Codec
	return nil
}
func (o CustomCodecCallOption) after(c *callInfo) {}





func MaxRetryRPCBufferSize(bytes int) CallOption {
	return MaxRetryRPCBufferSizeCallOption{bytes}
}




type MaxRetryRPCBufferSizeCallOption struct {
	MaxRetryRPCBufferSize int
}

func (o MaxRetryRPCBufferSizeCallOption) before(c *callInfo) error {
	c.maxRetryRPCBufferSize = o.MaxRetryRPCBufferSize
	return nil
}
func (o MaxRetryRPCBufferSizeCallOption) after(c *callInfo) {}


type payloadFormat uint8

const (
	compressionNone payloadFormat = 0 
	compressionMade payloadFormat = 1 
)


type parser struct {
	
	
	
	r io.Reader

	
	
	header [5]byte
}














func (p *parser) recvMsg(maxReceiveMessageSize int) (pf payloadFormat, msg []byte, err error) {
	if _, err := p.r.Read(p.header[:]); err != nil {
		return 0, nil, err
	}

	pf = payloadFormat(p.header[0])
	length := binary.BigEndian.Uint32(p.header[1:])

	if length == 0 {
		return pf, nil, nil
	}
	if int64(length) > int64(maxInt) {
		return 0, nil, status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max length allowed on current machine (%d vs. %d)", length, maxInt)
	}
	if int(length) > maxReceiveMessageSize {
		return 0, nil, status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max (%d vs. %d)", length, maxReceiveMessageSize)
	}
	
	
	msg = make([]byte, int(length))
	if _, err := p.r.Read(msg); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return 0, nil, err
	}
	return pf, msg, nil
}




func encode(c baseCodec, msg interface{}) ([]byte, error) {
	if msg == nil { 
		return nil, nil
	}
	b, err := c.Marshal(msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "grpc: error while marshaling: %v", err.Error())
	}
	if uint(len(b)) > math.MaxUint32 {
		return nil, status.Errorf(codes.ResourceExhausted, "grpc: message too large (%d bytes)", len(b))
	}
	return b, nil
}





func compress(in []byte, cp Compressor, compressor encoding.Compressor) ([]byte, error) {
	if compressor == nil && cp == nil {
		return nil, nil
	}
	wrapErr := func(err error) error {
		return status.Errorf(codes.Internal, "grpc: error while compressing: %v", err.Error())
	}
	cbuf := &bytes.Buffer{}
	if compressor != nil {
		z, _ := compressor.Compress(cbuf)
		if _, err := z.Write(in); err != nil {
			return nil, wrapErr(err)
		}
		if err := z.Close(); err != nil {
			return nil, wrapErr(err)
		}
	} else {
		if err := cp.Do(cbuf, in); err != nil {
			return nil, wrapErr(err)
		}
	}
	return cbuf.Bytes(), nil
}

const (
	payloadLen = 1
	sizeLen    = 4
	headerLen  = payloadLen + sizeLen
)



func msgHeader(data, compData []byte) (hdr []byte, payload []byte) {
	hdr = make([]byte, headerLen)
	if compData != nil {
		hdr[0] = byte(compressionMade)
		data = compData
	} else {
		hdr[0] = byte(compressionNone)
	}

	
	binary.BigEndian.PutUint32(hdr[payloadLen:], uint32(len(data)))
	return hdr, data
}

func outPayload(client bool, msg interface{}, data, payload []byte, t time.Time) *stats.OutPayload {
	return &stats.OutPayload{
		Client:     client,
		Payload:    msg,
		Data:       data,
		Length:     len(data),
		WireLength: len(payload) + headerLen,
		SentTime:   t,
	}
}

func checkRecvPayload(pf payloadFormat, recvCompress string, haveCompressor bool) *status.Status {
	switch pf {
	case compressionNone:
	case compressionMade:
		if recvCompress == "" || recvCompress == encoding.Identity {
			return status.New(codes.Internal, "grpc: compressed flag set with identity or empty encoding")
		}
		if !haveCompressor {
			return status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", recvCompress)
		}
	default:
		return status.Newf(codes.Internal, "grpc: received unexpected payload format %d", pf)
	}
	return nil
}




func recv(p *parser, c baseCodec, s *transport.Stream, dc Decompressor, m interface{}, maxReceiveMessageSize int, inPayload *stats.InPayload, compressor encoding.Compressor) error {
	pf, d, err := p.recvMsg(maxReceiveMessageSize)
	if err != nil {
		return err
	}
	if inPayload != nil {
		inPayload.WireLength = len(d)
	}

	if st := checkRecvPayload(pf, s.RecvCompress(), compressor != nil || dc != nil); st != nil {
		return st.Err()
	}

	if pf == compressionMade {
		
		
		if dc != nil {
			d, err = dc.Do(bytes.NewReader(d))
			if err != nil {
				return status.Errorf(codes.Internal, "grpc: failed to decompress the received message %v", err)
			}
		} else {
			dcReader, err := compressor.Decompress(bytes.NewReader(d))
			if err != nil {
				return status.Errorf(codes.Internal, "grpc: failed to decompress the received message %v", err)
			}
			d, err = ioutil.ReadAll(dcReader)
			if err != nil {
				return status.Errorf(codes.Internal, "grpc: failed to decompress the received message %v", err)
			}
		}
	}
	if len(d) > maxReceiveMessageSize {
		
		
		return status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max (%d vs. %d)", len(d), maxReceiveMessageSize)
	}
	if err := c.Unmarshal(d, m); err != nil {
		return status.Errorf(codes.Internal, "grpc: failed to unmarshal the received message %v", err)
	}
	if inPayload != nil {
		inPayload.RecvTime = time.Now()
		inPayload.Payload = m
		
		inPayload.Data = d
		inPayload.Length = len(d)
	}
	return nil
}

type rpcInfo struct {
	failfast bool
}

type rpcInfoContextKey struct{}

func newContextWithRPCInfo(ctx context.Context, failfast bool) context.Context {
	return context.WithValue(ctx, rpcInfoContextKey{}, &rpcInfo{failfast: failfast})
}

func rpcInfoFromContext(ctx context.Context) (s *rpcInfo, ok bool) {
	s, ok = ctx.Value(rpcInfoContextKey{}).(*rpcInfo)
	return
}





func Code(err error) codes.Code {
	if s, ok := status.FromError(err); ok {
		return s.Code()
	}
	return codes.Unknown
}





func ErrorDesc(err error) string {
	if s, ok := status.FromError(err); ok {
		return s.Message()
	}
	return err.Error()
}





func Errorf(c codes.Code, format string, a ...interface{}) error {
	return status.Errorf(c, format, a...)
}


func setCallInfoCodec(c *callInfo) error {
	if c.codec != nil {
		
		return nil
	}

	if c.contentSubtype == "" {
		
		c.codec = encoding.GetCodec(proto.Name)
		return nil
	}

	
	c.codec = encoding.GetCodec(c.contentSubtype)
	if c.codec == nil {
		return status.Errorf(codes.Internal, "no codec registered for content-subtype %s", c.contentSubtype)
	}
	return nil
}


func parseDialTarget(target string) (net string, addr string) {
	net = "tcp"

	m1 := strings.Index(target, ":")
	m2 := strings.Index(target, ":/")

	
	if m1 >= 0 && m2 < 0 {
		if n := target[0:m1]; n == "unix" {
			net = n
			addr = target[m1+1:]
			return net, addr
		}
	}
	if m2 >= 0 {
		t, err := url.Parse(target)
		if err != nil {
			return net, target
		}
		scheme := t.Scheme
		addr = t.Path
		if scheme == "unix" {
			net = scheme
			if addr == "" {
				addr = t.Host
			}
			return net, addr
		}
	}

	return net, target
}





type channelzData struct {
	callsStarted   int64
	callsFailed    int64
	callsSucceeded int64
	
	
	lastCallStartedTime int64
}









const (
	SupportPackageIsVersion3 = true
	SupportPackageIsVersion4 = true
	SupportPackageIsVersion5 = true
)

const grpcUA = "grpc-go/" + Version
