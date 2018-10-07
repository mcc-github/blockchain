





package http2

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	mathrand "math/rand"
	"net"
	"net/http"
	"net/textproto"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2/hpack"
	"golang.org/x/net/idna"
)

const (
	
	
	transportDefaultConnFlow = 1 << 30

	
	
	
	transportDefaultStreamFlow = 4 << 20

	
	
	transportDefaultStreamMinRefresh = 4 << 10

	defaultUserAgent = "Go-http-client/2.0"
)





type Transport struct {
	
	
	
	
	
	
	
	DialTLS func(network, addr string, cfg *tls.Config) (net.Conn, error)

	
	
	TLSClientConfig *tls.Config

	
	
	ConnPool ClientConnPool

	
	
	
	
	
	
	
	
	DisableCompression bool

	
	
	AllowHTTP bool

	
	
	
	
	
	
	
	MaxHeaderListSize uint32

	
	
	
	t1 *http.Transport

	connPoolOnce  sync.Once
	connPoolOrDef ClientConnPool 
}

func (t *Transport) maxHeaderListSize() uint32 {
	if t.MaxHeaderListSize == 0 {
		return 10 << 20
	}
	if t.MaxHeaderListSize == 0xffffffff {
		return 0
	}
	return t.MaxHeaderListSize
}

func (t *Transport) disableCompression() bool {
	return t.DisableCompression || (t.t1 != nil && t.t1.DisableCompression)
}

var errTransportVersion = errors.New("http2: ConfigureTransport is only supported starting at Go 1.6")




func ConfigureTransport(t1 *http.Transport) error {
	_, err := configureTransport(t1) 
	return err
}

func (t *Transport) connPool() ClientConnPool {
	t.connPoolOnce.Do(t.initConnPool)
	return t.connPoolOrDef
}

func (t *Transport) initConnPool() {
	if t.ConnPool != nil {
		t.connPoolOrDef = t.ConnPool
	} else {
		t.connPoolOrDef = &clientConnPool{t: t}
	}
}



type ClientConn struct {
	t         *Transport
	tconn     net.Conn             
	tlsState  *tls.ConnectionState 
	singleUse bool                 

	
	readerDone chan struct{} 
	readerErr  error         

	idleTimeout time.Duration 
	idleTimer   *time.Timer

	mu              sync.Mutex 
	cond            *sync.Cond 
	flow            flow       
	inflow          flow       
	closing         bool
	closed          bool
	wantSettingsAck bool                     
	goAway          *GoAwayFrame             
	goAwayDebug     string                   
	streams         map[uint32]*clientStream 
	nextStreamID    uint32
	pendingRequests int                       
	pings           map[[8]byte]chan struct{} 
	bw              *bufio.Writer
	br              *bufio.Reader
	fr              *Framer
	lastActive      time.Time
	
	maxFrameSize          uint32
	maxConcurrentStreams  uint32
	peerMaxHeaderListSize uint64
	initialWindowSize     uint32

	hbuf    bytes.Buffer 
	henc    *hpack.Encoder
	freeBuf [][]byte

	wmu  sync.Mutex 
	werr error      
}



type clientStream struct {
	cc            *ClientConn
	req           *http.Request
	trace         *clientTrace 
	ID            uint32
	resc          chan resAndError
	bufPipe       pipe 
	startedWrite  bool 
	requestedGzip bool
	on100         func() 

	flow        flow  
	inflow      flow  
	bytesRemain int64 
	readErr     error 
	stopReqBody error 
	didReset    bool  

	peerReset chan struct{} 
	resetErr  error         

	done chan struct{} 

	
	firstByte    bool  
	pastHeaders  bool  
	pastTrailers bool  
	num1xx       uint8 

	trailer    http.Header  
	resTrailer *http.Header 
}




func awaitRequestCancel(req *http.Request, done <-chan struct{}) error {
	ctx := reqContext(req)
	if req.Cancel == nil && ctx.Done() == nil {
		return nil
	}
	select {
	case <-req.Cancel:
		return errRequestCanceled
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

var got1xxFuncForTests func(int, textproto.MIMEHeader) error



func (cs *clientStream) get1xxTraceFunc() func(int, textproto.MIMEHeader) error {
	if fn := got1xxFuncForTests; fn != nil {
		return fn
	}
	return traceGot1xxResponseFunc(cs.trace)
}





func (cs *clientStream) awaitRequestCancel(req *http.Request) {
	if err := awaitRequestCancel(req, cs.done); err != nil {
		cs.cancelStream()
		cs.bufPipe.CloseWithError(err)
	}
}

func (cs *clientStream) cancelStream() {
	cc := cs.cc
	cc.mu.Lock()
	didReset := cs.didReset
	cs.didReset = true
	cc.mu.Unlock()

	if !didReset {
		cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
		cc.forgetStreamID(cs.ID)
	}
}



func (cs *clientStream) checkResetOrDone() error {
	select {
	case <-cs.peerReset:
		return cs.resetErr
	case <-cs.done:
		return errStreamClosed
	default:
		return nil
	}
}

func (cs *clientStream) getStartedWrite() bool {
	cc := cs.cc
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cs.startedWrite
}

func (cs *clientStream) abortRequestBodyWrite(err error) {
	if err == nil {
		panic("nil error")
	}
	cc := cs.cc
	cc.mu.Lock()
	cs.stopReqBody = err
	cc.cond.Broadcast()
	cc.mu.Unlock()
}

type stickyErrWriter struct {
	w   io.Writer
	err *error
}

func (sew stickyErrWriter) Write(p []byte) (n int, err error) {
	if *sew.err != nil {
		return 0, *sew.err
	}
	n, err = sew.w.Write(p)
	*sew.err = err
	return
}







type noCachedConnError struct{}

func (noCachedConnError) IsHTTP2NoCachedConnError() {}
func (noCachedConnError) Error() string             { return "http2: no cached connection was available" }




func isNoCachedConnError(err error) bool {
	_, ok := err.(interface{ IsHTTP2NoCachedConnError() })
	return ok
}

var ErrNoCachedConn error = noCachedConnError{}


type RoundTripOpt struct {
	
	
	
	
	OnlyCachedConn bool
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripOpt(req, RoundTripOpt{})
}



func authorityAddr(scheme string, authority string) (addr string) {
	host, port, err := net.SplitHostPort(authority)
	if err != nil { 
		port = "443"
		if scheme == "http" {
			port = "80"
		}
		host = authority
	}
	if a, err := idna.ToASCII(host); err == nil {
		host = a
	}
	
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return host + ":" + port
	}
	return net.JoinHostPort(host, port)
}


func (t *Transport) RoundTripOpt(req *http.Request, opt RoundTripOpt) (*http.Response, error) {
	if !(req.URL.Scheme == "https" || (req.URL.Scheme == "http" && t.AllowHTTP)) {
		return nil, errors.New("http2: unsupported scheme")
	}

	addr := authorityAddr(req.URL.Scheme, req.URL.Host)
	for retry := 0; ; retry++ {
		cc, err := t.connPool().GetClientConn(req, addr)
		if err != nil {
			t.vlogf("http2: Transport failed to get client conn for %s: %v", addr, err)
			return nil, err
		}
		traceGotConn(req, cc)
		res, gotErrAfterReqBodyWrite, err := cc.roundTrip(req)
		if err != nil && retry <= 6 {
			if req, err = shouldRetryRequest(req, err, gotErrAfterReqBodyWrite); err == nil {
				
				if retry == 0 {
					continue
				}
				backoff := float64(uint(1) << (uint(retry) - 1))
				backoff += backoff * (0.1 * mathrand.Float64())
				select {
				case <-time.After(time.Second * time.Duration(backoff)):
					continue
				case <-reqContext(req).Done():
					return nil, reqContext(req).Err()
				}
			}
		}
		if err != nil {
			t.vlogf("RoundTrip failure: %v", err)
			return nil, err
		}
		return res, nil
	}
}




func (t *Transport) CloseIdleConnections() {
	if cp, ok := t.connPool().(clientConnPoolIdleCloser); ok {
		cp.closeIdleConnections()
	}
}

var (
	errClientConnClosed    = errors.New("http2: client conn is closed")
	errClientConnUnusable  = errors.New("http2: client conn not usable")
	errClientConnGotGoAway = errors.New("http2: Transport received Server's graceful shutdown GOAWAY")
)





func shouldRetryRequest(req *http.Request, err error, afterBodyWrite bool) (*http.Request, error) {
	if !canRetryError(err) {
		return nil, err
	}
	
	
	if req.Body == nil || reqBodyIsNoBody(req.Body) {
		return req, nil
	}

	
	
	getBody := reqGetBody(req) 
	if getBody != nil {
		
		body, err := getBody()
		if err != nil {
			return nil, err
		}
		newReq := *req
		newReq.Body = body
		return &newReq, nil
	}

	
	
	
	
	
	if !afterBodyWrite {
		return req, nil
	}

	return nil, fmt.Errorf("http2: Transport: cannot retry err [%v] after Request.Body was written; define Request.GetBody to avoid this error", err)
}

func canRetryError(err error) bool {
	if err == errClientConnUnusable || err == errClientConnGotGoAway {
		return true
	}
	if se, ok := err.(StreamError); ok {
		return se.Code == ErrCodeRefusedStream
	}
	return false
}

func (t *Transport) dialClientConn(addr string, singleUse bool) (*ClientConn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	tconn, err := t.dialTLS()("tcp", addr, t.newTLSConfig(host))
	if err != nil {
		return nil, err
	}
	return t.newClientConn(tconn, singleUse)
}

func (t *Transport) newTLSConfig(host string) *tls.Config {
	cfg := new(tls.Config)
	if t.TLSClientConfig != nil {
		*cfg = *cloneTLSConfig(t.TLSClientConfig)
	}
	if !strSliceContains(cfg.NextProtos, NextProtoTLS) {
		cfg.NextProtos = append([]string{NextProtoTLS}, cfg.NextProtos...)
	}
	if cfg.ServerName == "" {
		cfg.ServerName = host
	}
	return cfg
}

func (t *Transport) dialTLS() func(string, string, *tls.Config) (net.Conn, error) {
	if t.DialTLS != nil {
		return t.DialTLS
	}
	return t.dialTLSDefault
}

func (t *Transport) dialTLSDefault(network, addr string, cfg *tls.Config) (net.Conn, error) {
	cn, err := tls.Dial(network, addr, cfg)
	if err != nil {
		return nil, err
	}
	if err := cn.Handshake(); err != nil {
		return nil, err
	}
	if !cfg.InsecureSkipVerify {
		if err := cn.VerifyHostname(cfg.ServerName); err != nil {
			return nil, err
		}
	}
	state := cn.ConnectionState()
	if p := state.NegotiatedProtocol; p != NextProtoTLS {
		return nil, fmt.Errorf("http2: unexpected ALPN protocol %q; want %q", p, NextProtoTLS)
	}
	if !state.NegotiatedProtocolIsMutual {
		return nil, errors.New("http2: could not negotiate protocol mutually")
	}
	return cn, nil
}



func (t *Transport) disableKeepAlives() bool {
	return t.t1 != nil && t.t1.DisableKeepAlives
}

func (t *Transport) expectContinueTimeout() time.Duration {
	if t.t1 == nil {
		return 0
	}
	return transportExpectContinueTimeout(t.t1)
}

func (t *Transport) NewClientConn(c net.Conn) (*ClientConn, error) {
	return t.newClientConn(c, false)
}

func (t *Transport) newClientConn(c net.Conn, singleUse bool) (*ClientConn, error) {
	cc := &ClientConn{
		t:                     t,
		tconn:                 c,
		readerDone:            make(chan struct{}),
		nextStreamID:          1,
		maxFrameSize:          16 << 10,           
		initialWindowSize:     65535,              
		maxConcurrentStreams:  1000,               
		peerMaxHeaderListSize: 0xffffffffffffffff, 
		streams:               make(map[uint32]*clientStream),
		singleUse:             singleUse,
		wantSettingsAck:       true,
		pings:                 make(map[[8]byte]chan struct{}),
	}
	if d := t.idleConnTimeout(); d != 0 {
		cc.idleTimeout = d
		cc.idleTimer = time.AfterFunc(d, cc.onIdleTimeout)
	}
	if VerboseLogs {
		t.vlogf("http2: Transport creating client conn %p to %v", cc, c.RemoteAddr())
	}

	cc.cond = sync.NewCond(&cc.mu)
	cc.flow.add(int32(initialWindowSize))

	
	
	cc.bw = bufio.NewWriter(stickyErrWriter{c, &cc.werr})
	cc.br = bufio.NewReader(c)
	cc.fr = NewFramer(cc.bw, cc.br)
	cc.fr.ReadMetaHeaders = hpack.NewDecoder(initialHeaderTableSize, nil)
	cc.fr.MaxHeaderListSize = t.maxHeaderListSize()

	
	
	cc.henc = hpack.NewEncoder(&cc.hbuf)

	if t.AllowHTTP {
		cc.nextStreamID = 3
	}

	if cs, ok := c.(connectionStater); ok {
		state := cs.ConnectionState()
		cc.tlsState = &state
	}

	initialSettings := []Setting{
		{ID: SettingEnablePush, Val: 0},
		{ID: SettingInitialWindowSize, Val: transportDefaultStreamFlow},
	}
	if max := t.maxHeaderListSize(); max != 0 {
		initialSettings = append(initialSettings, Setting{ID: SettingMaxHeaderListSize, Val: max})
	}

	cc.bw.Write(clientPreface)
	cc.fr.WriteSettings(initialSettings...)
	cc.fr.WriteWindowUpdate(0, transportDefaultConnFlow)
	cc.inflow.add(transportDefaultConnFlow + initialWindowSize)
	cc.bw.Flush()
	if cc.werr != nil {
		return nil, cc.werr
	}

	go cc.readLoop()
	return cc, nil
}

func (cc *ClientConn) setGoAway(f *GoAwayFrame) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	old := cc.goAway
	cc.goAway = f

	
	if cc.goAwayDebug == "" {
		cc.goAwayDebug = string(f.DebugData())
	}
	if old != nil && old.ErrCode != ErrCodeNo {
		cc.goAway.ErrCode = old.ErrCode
	}
	last := f.LastStreamID
	for streamID, cs := range cc.streams {
		if streamID > last {
			select {
			case cs.resc <- resAndError{err: errClientConnGotGoAway}:
			default:
			}
		}
	}
}



func (cc *ClientConn) CanTakeNewRequest() bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.canTakeNewRequestLocked()
}



type clientConnIdleState struct {
	canTakeNewRequest bool
	freshConn         bool 
}

func (cc *ClientConn) idleState() clientConnIdleState {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.idleStateLocked()
}

func (cc *ClientConn) idleStateLocked() (st clientConnIdleState) {
	if cc.singleUse && cc.nextStreamID > 1 {
		return
	}
	st.canTakeNewRequest = cc.goAway == nil && !cc.closed && !cc.closing &&
		int64(cc.nextStreamID)+int64(cc.pendingRequests) < math.MaxInt32
	st.freshConn = cc.nextStreamID == 1 && st.canTakeNewRequest
	return
}

func (cc *ClientConn) canTakeNewRequestLocked() bool {
	st := cc.idleStateLocked()
	return st.canTakeNewRequest
}







func (cc *ClientConn) onIdleTimeout() {
	cc.closeIfIdle()
}

func (cc *ClientConn) closeIfIdle() {
	cc.mu.Lock()
	if len(cc.streams) > 0 {
		cc.mu.Unlock()
		return
	}
	cc.closed = true
	nextID := cc.nextStreamID
	
	cc.mu.Unlock()

	if VerboseLogs {
		cc.vlogf("http2: Transport closing idle conn %p (forSingleUse=%v, maxStream=%v)", cc, cc.singleUse, nextID-2)
	}
	cc.tconn.Close()
}

var shutdownEnterWaitStateHook = func() {}



func (cc *ClientConn) shutdown(ctx contextContext) error {
	if err := cc.sendGoAway(); err != nil {
		return err
	}
	
	done := make(chan error, 1)
	cancelled := false 
	go func() {
		cc.mu.Lock()
		defer cc.mu.Unlock()
		for {
			if len(cc.streams) == 0 || cc.closed {
				cc.closed = true
				done <- cc.tconn.Close()
				break
			}
			if cancelled {
				break
			}
			cc.cond.Wait()
		}
	}()
	shutdownEnterWaitStateHook()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		cc.mu.Lock()
		
		cancelled = true
		cc.cond.Broadcast()
		cc.mu.Unlock()
		return ctx.Err()
	}
}

func (cc *ClientConn) sendGoAway() error {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	if cc.closing {
		
		return nil
	}
	
	maxStreamID := cc.nextStreamID
	if err := cc.fr.WriteGoAway(maxStreamID, ErrCodeNo, nil); err != nil {
		return err
	}
	if err := cc.bw.Flush(); err != nil {
		return err
	}
	
	cc.closing = true
	return nil
}




func (cc *ClientConn) Close() error {
	cc.mu.Lock()
	defer cc.cond.Broadcast()
	defer cc.mu.Unlock()
	err := errors.New("http2: client connection force closed via ClientConn.Close")
	for id, cs := range cc.streams {
		select {
		case cs.resc <- resAndError{err: err}:
		default:
		}
		cs.bufPipe.CloseWithError(err)
		delete(cc.streams, id)
	}
	cc.closed = true
	return cc.tconn.Close()
}

const maxAllocFrameSize = 512 << 10





func (cc *ClientConn) frameScratchBuffer() []byte {
	cc.mu.Lock()
	size := cc.maxFrameSize
	if size > maxAllocFrameSize {
		size = maxAllocFrameSize
	}
	for i, buf := range cc.freeBuf {
		if len(buf) >= int(size) {
			cc.freeBuf[i] = nil
			cc.mu.Unlock()
			return buf[:size]
		}
	}
	cc.mu.Unlock()
	return make([]byte, size)
}

func (cc *ClientConn) putFrameScratchBuffer(buf []byte) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	const maxBufs = 4 
	if len(cc.freeBuf) < maxBufs {
		cc.freeBuf = append(cc.freeBuf, buf)
		return
	}
	for i, old := range cc.freeBuf {
		if old == nil {
			cc.freeBuf[i] = buf
			return
		}
	}
	
}



var errRequestCanceled = errors.New("net/http: request canceled")

func commaSeparatedTrailers(req *http.Request) (string, error) {
	keys := make([]string, 0, len(req.Trailer))
	for k := range req.Trailer {
		k = http.CanonicalHeaderKey(k)
		switch k {
		case "Transfer-Encoding", "Trailer", "Content-Length":
			return "", &badStringError{"invalid Trailer key", k}
		}
		keys = append(keys, k)
	}
	if len(keys) > 0 {
		sort.Strings(keys)
		return strings.Join(keys, ","), nil
	}
	return "", nil
}

func (cc *ClientConn) responseHeaderTimeout() time.Duration {
	if cc.t.t1 != nil {
		return cc.t.t1.ResponseHeaderTimeout
	}
	
	
	
	
	return 0
}




func checkConnHeaders(req *http.Request) error {
	if v := req.Header.Get("Upgrade"); v != "" {
		return fmt.Errorf("http2: invalid Upgrade request header: %q", req.Header["Upgrade"])
	}
	if vv := req.Header["Transfer-Encoding"]; len(vv) > 0 && (len(vv) > 1 || vv[0] != "" && vv[0] != "chunked") {
		return fmt.Errorf("http2: invalid Transfer-Encoding request header: %q", vv)
	}
	if vv := req.Header["Connection"]; len(vv) > 0 && (len(vv) > 1 || vv[0] != "" && !strings.EqualFold(vv[0], "close") && !strings.EqualFold(vv[0], "keep-alive")) {
		return fmt.Errorf("http2: invalid Connection request header: %q", vv)
	}
	return nil
}




func actualContentLength(req *http.Request) int64 {
	if req.Body == nil || reqBodyIsNoBody(req.Body) {
		return 0
	}
	if req.ContentLength != 0 {
		return req.ContentLength
	}
	return -1
}

func (cc *ClientConn) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, _, err := cc.roundTrip(req)
	return resp, err
}

func (cc *ClientConn) roundTrip(req *http.Request) (res *http.Response, gotErrAfterReqBodyWrite bool, err error) {
	if err := checkConnHeaders(req); err != nil {
		return nil, false, err
	}
	if cc.idleTimer != nil {
		cc.idleTimer.Stop()
	}

	trailers, err := commaSeparatedTrailers(req)
	if err != nil {
		return nil, false, err
	}
	hasTrailers := trailers != ""

	cc.mu.Lock()
	if err := cc.awaitOpenSlotForRequest(req); err != nil {
		cc.mu.Unlock()
		return nil, false, err
	}

	body := req.Body
	contentLen := actualContentLength(req)
	hasBody := contentLen != 0

	
	var requestedGzip bool
	if !cc.t.disableCompression() &&
		req.Header.Get("Accept-Encoding") == "" &&
		req.Header.Get("Range") == "" &&
		req.Method != "HEAD" {
		
		
		
		
		
		
		
		
		
		
		
		
		requestedGzip = true
	}

	
	
	
	hdrs, err := cc.encodeHeaders(req, requestedGzip, trailers, contentLen)
	if err != nil {
		cc.mu.Unlock()
		return nil, false, err
	}

	cs := cc.newStream()
	cs.req = req
	cs.trace = requestTrace(req)
	cs.requestedGzip = requestedGzip
	bodyWriter := cc.t.getBodyWriterState(cs, body)
	cs.on100 = bodyWriter.on100

	cc.wmu.Lock()
	endStream := !hasBody && !hasTrailers
	werr := cc.writeHeaders(cs.ID, endStream, int(cc.maxFrameSize), hdrs)
	cc.wmu.Unlock()
	traceWroteHeaders(cs.trace)
	cc.mu.Unlock()

	if werr != nil {
		if hasBody {
			req.Body.Close() 
			bodyWriter.cancel()
		}
		cc.forgetStreamID(cs.ID)
		
		
		traceWroteRequest(cs.trace, werr)
		return nil, false, werr
	}

	var respHeaderTimer <-chan time.Time
	if hasBody {
		bodyWriter.scheduleBodyWrite()
	} else {
		traceWroteRequest(cs.trace, nil)
		if d := cc.responseHeaderTimeout(); d != 0 {
			timer := time.NewTimer(d)
			defer timer.Stop()
			respHeaderTimer = timer.C
		}
	}

	readLoopResCh := cs.resc
	bodyWritten := false
	ctx := reqContext(req)

	handleReadLoopResponse := func(re resAndError) (*http.Response, bool, error) {
		res := re.res
		if re.err != nil || res.StatusCode > 299 {
			
			
			
			
			
			
			
			
			
			bodyWriter.cancel()
			cs.abortRequestBodyWrite(errStopReqBodyWrite)
		}
		if re.err != nil {
			cc.forgetStreamID(cs.ID)
			return nil, cs.getStartedWrite(), re.err
		}
		res.Request = req
		res.TLS = cc.tlsState
		return res, false, nil
	}

	for {
		select {
		case re := <-readLoopResCh:
			return handleReadLoopResponse(re)
		case <-respHeaderTimer:
			if !hasBody || bodyWritten {
				cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
			} else {
				bodyWriter.cancel()
				cs.abortRequestBodyWrite(errStopReqBodyWriteAndCancel)
			}
			cc.forgetStreamID(cs.ID)
			return nil, cs.getStartedWrite(), errTimeout
		case <-ctx.Done():
			if !hasBody || bodyWritten {
				cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
			} else {
				bodyWriter.cancel()
				cs.abortRequestBodyWrite(errStopReqBodyWriteAndCancel)
			}
			cc.forgetStreamID(cs.ID)
			return nil, cs.getStartedWrite(), ctx.Err()
		case <-req.Cancel:
			if !hasBody || bodyWritten {
				cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
			} else {
				bodyWriter.cancel()
				cs.abortRequestBodyWrite(errStopReqBodyWriteAndCancel)
			}
			cc.forgetStreamID(cs.ID)
			return nil, cs.getStartedWrite(), errRequestCanceled
		case <-cs.peerReset:
			
			
			
			return nil, cs.getStartedWrite(), cs.resetErr
		case err := <-bodyWriter.resc:
			
			select {
			case re := <-readLoopResCh:
				return handleReadLoopResponse(re)
			default:
			}
			if err != nil {
				return nil, cs.getStartedWrite(), err
			}
			bodyWritten = true
			if d := cc.responseHeaderTimeout(); d != 0 {
				timer := time.NewTimer(d)
				defer timer.Stop()
				respHeaderTimer = timer.C
			}
		}
	}
}



func (cc *ClientConn) awaitOpenSlotForRequest(req *http.Request) error {
	var waitingForConn chan struct{}
	var waitingForConnErr error 
	for {
		cc.lastActive = time.Now()
		if cc.closed || !cc.canTakeNewRequestLocked() {
			if waitingForConn != nil {
				close(waitingForConn)
			}
			return errClientConnUnusable
		}
		if int64(len(cc.streams))+1 <= int64(cc.maxConcurrentStreams) {
			if waitingForConn != nil {
				close(waitingForConn)
			}
			return nil
		}
		
		
		
		if waitingForConn == nil {
			waitingForConn = make(chan struct{})
			go func() {
				if err := awaitRequestCancel(req, waitingForConn); err != nil {
					cc.mu.Lock()
					waitingForConnErr = err
					cc.cond.Broadcast()
					cc.mu.Unlock()
				}
			}()
		}
		cc.pendingRequests++
		cc.cond.Wait()
		cc.pendingRequests--
		if waitingForConnErr != nil {
			return waitingForConnErr
		}
	}
}


func (cc *ClientConn) writeHeaders(streamID uint32, endStream bool, maxFrameSize int, hdrs []byte) error {
	first := true 
	for len(hdrs) > 0 && cc.werr == nil {
		chunk := hdrs
		if len(chunk) > maxFrameSize {
			chunk = chunk[:maxFrameSize]
		}
		hdrs = hdrs[len(chunk):]
		endHeaders := len(hdrs) == 0
		if first {
			cc.fr.WriteHeaders(HeadersFrameParam{
				StreamID:      streamID,
				BlockFragment: chunk,
				EndStream:     endStream,
				EndHeaders:    endHeaders,
			})
			first = false
		} else {
			cc.fr.WriteContinuation(streamID, endHeaders, chunk)
		}
	}
	
	
	
	
	cc.bw.Flush()
	return cc.werr
}


var (
	
	errStopReqBodyWrite = errors.New("http2: aborting request body write")

	
	errStopReqBodyWriteAndCancel = errors.New("http2: canceling request")
)

func (cs *clientStream) writeRequestBody(body io.Reader, bodyCloser io.Closer) (err error) {
	cc := cs.cc
	sentEnd := false 
	buf := cc.frameScratchBuffer()
	defer cc.putFrameScratchBuffer(buf)

	defer func() {
		traceWroteRequest(cs.trace, err)
		
		
		
		
		cerr := bodyCloser.Close()
		if err == nil {
			err = cerr
		}
	}()

	req := cs.req
	hasTrailers := req.Trailer != nil

	var sawEOF bool
	for !sawEOF {
		n, err := body.Read(buf)
		if err == io.EOF {
			sawEOF = true
			err = nil
		} else if err != nil {
			return err
		}

		remain := buf[:n]
		for len(remain) > 0 && err == nil {
			var allowed int32
			allowed, err = cs.awaitFlowControl(len(remain))
			switch {
			case err == errStopReqBodyWrite:
				return err
			case err == errStopReqBodyWriteAndCancel:
				cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
				return err
			case err != nil:
				return err
			}
			cc.wmu.Lock()
			data := remain[:allowed]
			remain = remain[allowed:]
			sentEnd = sawEOF && len(remain) == 0 && !hasTrailers
			err = cc.fr.WriteData(cs.ID, sentEnd, data)
			if err == nil {
				
				
				
				
				
				
				err = cc.bw.Flush()
			}
			cc.wmu.Unlock()
		}
		if err != nil {
			return err
		}
	}

	if sentEnd {
		
		
		
		return nil
	}

	var trls []byte
	if hasTrailers {
		cc.mu.Lock()
		trls, err = cc.encodeTrailers(req)
		cc.mu.Unlock()
		if err != nil {
			cc.writeStreamReset(cs.ID, ErrCodeInternal, err)
			cc.forgetStreamID(cs.ID)
			return err
		}
	}

	cc.mu.Lock()
	maxFrameSize := int(cc.maxFrameSize)
	cc.mu.Unlock()

	cc.wmu.Lock()
	defer cc.wmu.Unlock()

	
	
	if len(trls) > 0 {
		err = cc.writeHeaders(cs.ID, true, maxFrameSize, trls)
	} else {
		err = cc.fr.WriteData(cs.ID, true, nil)
	}
	if ferr := cc.bw.Flush(); ferr != nil && err == nil {
		err = ferr
	}
	return err
}





func (cs *clientStream) awaitFlowControl(maxBytes int) (taken int32, err error) {
	cc := cs.cc
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for {
		if cc.closed {
			return 0, errClientConnClosed
		}
		if cs.stopReqBody != nil {
			return 0, cs.stopReqBody
		}
		if err := cs.checkResetOrDone(); err != nil {
			return 0, err
		}
		if a := cs.flow.available(); a > 0 {
			take := a
			if int(take) > maxBytes {

				take = int32(maxBytes) 
			}
			if take > int32(cc.maxFrameSize) {
				take = int32(cc.maxFrameSize)
			}
			cs.flow.take(take)
			return take, nil
		}
		cc.cond.Wait()
	}
}

type badStringError struct {
	what string
	str  string
}

func (e *badStringError) Error() string { return fmt.Sprintf("%s %q", e.what, e.str) }


func (cc *ClientConn) encodeHeaders(req *http.Request, addGzipHeader bool, trailers string, contentLength int64) ([]byte, error) {
	cc.hbuf.Reset()

	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	host, err := httpguts.PunycodeHostPort(host)
	if err != nil {
		return nil, err
	}

	var path string
	if req.Method != "CONNECT" {
		path = req.URL.RequestURI()
		if !validPseudoPath(path) {
			orig := path
			path = strings.TrimPrefix(path, req.URL.Scheme+"://"+host)
			if !validPseudoPath(path) {
				if req.URL.Opaque != "" {
					return nil, fmt.Errorf("invalid request :path %q from URL.Opaque = %q", orig, req.URL.Opaque)
				} else {
					return nil, fmt.Errorf("invalid request :path %q", orig)
				}
			}
		}
	}

	
	
	
	for k, vv := range req.Header {
		if !httpguts.ValidHeaderFieldName(k) {
			return nil, fmt.Errorf("invalid HTTP header name %q", k)
		}
		for _, v := range vv {
			if !httpguts.ValidHeaderFieldValue(v) {
				return nil, fmt.Errorf("invalid HTTP header value %q for header %q", v, k)
			}
		}
	}

	enumerateHeaders := func(f func(name, value string)) {
		
		
		
		
		
		f(":authority", host)
		f(":method", req.Method)
		if req.Method != "CONNECT" {
			f(":path", path)
			f(":scheme", req.URL.Scheme)
		}
		if trailers != "" {
			f("trailer", trailers)
		}

		var didUA bool
		for k, vv := range req.Header {
			if strings.EqualFold(k, "host") || strings.EqualFold(k, "content-length") {
				
				
				continue
			} else if strings.EqualFold(k, "connection") || strings.EqualFold(k, "proxy-connection") ||
				strings.EqualFold(k, "transfer-encoding") || strings.EqualFold(k, "upgrade") ||
				strings.EqualFold(k, "keep-alive") {
				
				
				
				
				continue
			} else if strings.EqualFold(k, "user-agent") {
				
				
				
				
				didUA = true
				if len(vv) < 1 {
					continue
				}
				vv = vv[:1]
				if vv[0] == "" {
					continue
				}

			}

			for _, v := range vv {
				f(k, v)
			}
		}
		if shouldSendReqContentLength(req.Method, contentLength) {
			f("content-length", strconv.FormatInt(contentLength, 10))
		}
		if addGzipHeader {
			f("accept-encoding", "gzip")
		}
		if !didUA {
			f("user-agent", defaultUserAgent)
		}
	}

	
	
	
	
	hlSize := uint64(0)
	enumerateHeaders(func(name, value string) {
		hf := hpack.HeaderField{Name: name, Value: value}
		hlSize += uint64(hf.Size())
	})

	if hlSize > cc.peerMaxHeaderListSize {
		return nil, errRequestHeaderListSize
	}

	trace := requestTrace(req)
	traceHeaders := traceHasWroteHeaderField(trace)

	
	enumerateHeaders(func(name, value string) {
		name = strings.ToLower(name)
		cc.writeHeader(name, value)
		if traceHeaders {
			traceWroteHeaderField(trace, name, value)
		}
	})

	return cc.hbuf.Bytes(), nil
}






func shouldSendReqContentLength(method string, contentLength int64) bool {
	if contentLength > 0 {
		return true
	}
	if contentLength < 0 {
		return false
	}
	
	
	switch method {
	case "POST", "PUT", "PATCH":
		return true
	default:
		return false
	}
}


func (cc *ClientConn) encodeTrailers(req *http.Request) ([]byte, error) {
	cc.hbuf.Reset()

	hlSize := uint64(0)
	for k, vv := range req.Trailer {
		for _, v := range vv {
			hf := hpack.HeaderField{Name: k, Value: v}
			hlSize += uint64(hf.Size())
		}
	}
	if hlSize > cc.peerMaxHeaderListSize {
		return nil, errRequestHeaderListSize
	}

	for k, vv := range req.Trailer {
		
		
		lowKey := strings.ToLower(k)
		for _, v := range vv {
			cc.writeHeader(lowKey, v)
		}
	}
	return cc.hbuf.Bytes(), nil
}

func (cc *ClientConn) writeHeader(name, value string) {
	if VerboseLogs {
		log.Printf("http2: Transport encoding header %q = %q", name, value)
	}
	cc.henc.WriteField(hpack.HeaderField{Name: name, Value: value})
}

type resAndError struct {
	res *http.Response
	err error
}


func (cc *ClientConn) newStream() *clientStream {
	cs := &clientStream{
		cc:        cc,
		ID:        cc.nextStreamID,
		resc:      make(chan resAndError, 1),
		peerReset: make(chan struct{}),
		done:      make(chan struct{}),
	}
	cs.flow.add(int32(cc.initialWindowSize))
	cs.flow.setConnFlow(&cc.flow)
	cs.inflow.add(transportDefaultStreamFlow)
	cs.inflow.setConnFlow(&cc.inflow)
	cc.nextStreamID += 2
	cc.streams[cs.ID] = cs
	return cs
}

func (cc *ClientConn) forgetStreamID(id uint32) {
	cc.streamByID(id, true)
}

func (cc *ClientConn) streamByID(id uint32, andRemove bool) *clientStream {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cs := cc.streams[id]
	if andRemove && cs != nil && !cc.closed {
		cc.lastActive = time.Now()
		delete(cc.streams, id)
		if len(cc.streams) == 0 && cc.idleTimer != nil {
			cc.idleTimer.Reset(cc.idleTimeout)
		}
		close(cs.done)
		
		
		cc.cond.Broadcast()
	}
	return cs
}


type clientConnReadLoop struct {
	cc            *ClientConn
	closeWhenIdle bool
}


func (cc *ClientConn) readLoop() {
	rl := &clientConnReadLoop{cc: cc}
	defer rl.cleanup()
	cc.readerErr = rl.run()
	if ce, ok := cc.readerErr.(ConnectionError); ok {
		cc.wmu.Lock()
		cc.fr.WriteGoAway(0, ErrCode(ce), nil)
		cc.wmu.Unlock()
	}
}



type GoAwayError struct {
	LastStreamID uint32
	ErrCode      ErrCode
	DebugData    string
}

func (e GoAwayError) Error() string {
	return fmt.Sprintf("http2: server sent GOAWAY and closed the connection; LastStreamID=%v, ErrCode=%v, debug=%q",
		e.LastStreamID, e.ErrCode, e.DebugData)
}

func isEOFOrNetReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}

func (rl *clientConnReadLoop) cleanup() {
	cc := rl.cc
	defer cc.tconn.Close()
	defer cc.t.connPool().MarkDead(cc)
	defer close(cc.readerDone)

	if cc.idleTimer != nil {
		cc.idleTimer.Stop()
	}

	
	
	
	err := cc.readerErr
	cc.mu.Lock()
	if cc.goAway != nil && isEOFOrNetReadError(err) {
		err = GoAwayError{
			LastStreamID: cc.goAway.LastStreamID,
			ErrCode:      cc.goAway.ErrCode,
			DebugData:    cc.goAwayDebug,
		}
	} else if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	for _, cs := range cc.streams {
		cs.bufPipe.CloseWithError(err) 
		select {
		case cs.resc <- resAndError{err: err}:
		default:
		}
		close(cs.done)
	}
	cc.closed = true
	cc.cond.Broadcast()
	cc.mu.Unlock()
}

func (rl *clientConnReadLoop) run() error {
	cc := rl.cc
	rl.closeWhenIdle = cc.t.disableKeepAlives() || cc.singleUse
	gotReply := false 
	gotSettings := false
	for {
		f, err := cc.fr.ReadFrame()
		if err != nil {
			cc.vlogf("http2: Transport readFrame error on conn %p: (%T) %v", cc, err, err)
		}
		if se, ok := err.(StreamError); ok {
			if cs := cc.streamByID(se.StreamID, false); cs != nil {
				cs.cc.writeStreamReset(cs.ID, se.Code, err)
				cs.cc.forgetStreamID(cs.ID)
				if se.Cause == nil {
					se.Cause = cc.fr.errDetail
				}
				rl.endStreamError(cs, se)
			}
			continue
		} else if err != nil {
			return err
		}
		if VerboseLogs {
			cc.vlogf("http2: Transport received %s", summarizeFrame(f))
		}
		if !gotSettings {
			if _, ok := f.(*SettingsFrame); !ok {
				cc.logf("protocol error: received %T before a SETTINGS frame", f)
				return ConnectionError(ErrCodeProtocol)
			}
			gotSettings = true
		}
		maybeIdle := false 

		switch f := f.(type) {
		case *MetaHeadersFrame:
			err = rl.processHeaders(f)
			maybeIdle = true
			gotReply = true
		case *DataFrame:
			err = rl.processData(f)
			maybeIdle = true
		case *GoAwayFrame:
			err = rl.processGoAway(f)
			maybeIdle = true
		case *RSTStreamFrame:
			err = rl.processResetStream(f)
			maybeIdle = true
		case *SettingsFrame:
			err = rl.processSettings(f)
		case *PushPromiseFrame:
			err = rl.processPushPromise(f)
		case *WindowUpdateFrame:
			err = rl.processWindowUpdate(f)
		case *PingFrame:
			err = rl.processPing(f)
		default:
			cc.logf("Transport: unhandled response frame type %T", f)
		}
		if err != nil {
			if VerboseLogs {
				cc.vlogf("http2: Transport conn %p received error from processing frame %v: %v", cc, summarizeFrame(f), err)
			}
			return err
		}
		if rl.closeWhenIdle && gotReply && maybeIdle {
			cc.closeIfIdle()
		}
	}
}

func (rl *clientConnReadLoop) processHeaders(f *MetaHeadersFrame) error {
	cc := rl.cc
	cs := cc.streamByID(f.StreamID, false)
	if cs == nil {
		
		
		
		return nil
	}
	if f.StreamEnded() {
		
		
		
		
		
		
		
		
		
		
		
		if cs.req.Body != nil {
			defer cc.forgetStreamID(f.StreamID)
		} else {
			cc.forgetStreamID(f.StreamID)
		}
	}
	if !cs.firstByte {
		if cs.trace != nil {
			
			
			
			
			traceFirstResponseByte(cs.trace)
		}
		cs.firstByte = true
	}
	if !cs.pastHeaders {
		cs.pastHeaders = true
	} else {
		return rl.processTrailers(cs, f)
	}

	res, err := rl.handleResponse(cs, f)
	if err != nil {
		if _, ok := err.(ConnectionError); ok {
			return err
		}
		
		cs.cc.writeStreamReset(f.StreamID, ErrCodeProtocol, err)
		cc.forgetStreamID(cs.ID)
		cs.resc <- resAndError{err: err}
		return nil 
	}
	if res == nil {
		
		return nil
	}
	cs.resTrailer = &res.Trailer
	cs.resc <- resAndError{res: res}
	return nil
}







func (rl *clientConnReadLoop) handleResponse(cs *clientStream, f *MetaHeadersFrame) (*http.Response, error) {
	if f.Truncated {
		return nil, errResponseHeaderListSize
	}

	status := f.PseudoValue("status")
	if status == "" {
		return nil, errors.New("malformed response from server: missing status pseudo header")
	}
	statusCode, err := strconv.Atoi(status)
	if err != nil {
		return nil, errors.New("malformed response from server: malformed non-numeric status pseudo header")
	}

	header := make(http.Header)
	res := &http.Response{
		Proto:      "HTTP/2.0",
		ProtoMajor: 2,
		Header:     header,
		StatusCode: statusCode,
		Status:     status + " " + http.StatusText(statusCode),
	}
	for _, hf := range f.RegularFields() {
		key := http.CanonicalHeaderKey(hf.Name)
		if key == "Trailer" {
			t := res.Trailer
			if t == nil {
				t = make(http.Header)
				res.Trailer = t
			}
			foreachHeaderElement(hf.Value, func(v string) {
				t[http.CanonicalHeaderKey(v)] = nil
			})
		} else {
			header[key] = append(header[key], hf.Value)
		}
	}

	if statusCode >= 100 && statusCode <= 199 {
		cs.num1xx++
		const max1xxResponses = 5 
		if cs.num1xx > max1xxResponses {
			return nil, errors.New("http2: too many 1xx informational responses")
		}
		if fn := cs.get1xxTraceFunc(); fn != nil {
			if err := fn(statusCode, textproto.MIMEHeader(header)); err != nil {
				return nil, err
			}
		}
		if statusCode == 100 {
			traceGot100Continue(cs.trace)
			if cs.on100 != nil {
				cs.on100() 
			}
		}
		cs.pastHeaders = false 
		return nil, nil
	}

	streamEnded := f.StreamEnded()
	isHead := cs.req.Method == "HEAD"
	if !streamEnded || isHead {
		res.ContentLength = -1
		if clens := res.Header["Content-Length"]; len(clens) == 1 {
			if clen64, err := strconv.ParseInt(clens[0], 10, 64); err == nil {
				res.ContentLength = clen64
			} else {
				
				
			}
		} else if len(clens) > 1 {
			
			
		}
	}

	if streamEnded || isHead {
		res.Body = noBody
		return res, nil
	}

	cs.bufPipe = pipe{b: &dataBuffer{expected: res.ContentLength}}
	cs.bytesRemain = res.ContentLength
	res.Body = transportResponseBody{cs}
	go cs.awaitRequestCancel(cs.req)

	if cs.requestedGzip && res.Header.Get("Content-Encoding") == "gzip" {
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Length")
		res.ContentLength = -1
		res.Body = &gzipReader{body: res.Body}
		setResponseUncompressed(res)
	}
	return res, nil
}

func (rl *clientConnReadLoop) processTrailers(cs *clientStream, f *MetaHeadersFrame) error {
	if cs.pastTrailers {
		
		return ConnectionError(ErrCodeProtocol)
	}
	cs.pastTrailers = true
	if !f.StreamEnded() {
		
		
		return ConnectionError(ErrCodeProtocol)
	}
	if len(f.PseudoFields()) > 0 {
		
		
		return ConnectionError(ErrCodeProtocol)
	}

	trailer := make(http.Header)
	for _, hf := range f.RegularFields() {
		key := http.CanonicalHeaderKey(hf.Name)
		trailer[key] = append(trailer[key], hf.Value)
	}
	cs.trailer = trailer

	rl.endStream(cs)
	return nil
}




type transportResponseBody struct {
	cs *clientStream
}

func (b transportResponseBody) Read(p []byte) (n int, err error) {
	cs := b.cs
	cc := cs.cc

	if cs.readErr != nil {
		return 0, cs.readErr
	}
	n, err = b.cs.bufPipe.Read(p)
	if cs.bytesRemain != -1 {
		if int64(n) > cs.bytesRemain {
			n = int(cs.bytesRemain)
			if err == nil {
				err = errors.New("net/http: server replied with more than declared Content-Length; truncated")
				cc.writeStreamReset(cs.ID, ErrCodeProtocol, err)
			}
			cs.readErr = err
			return int(cs.bytesRemain), err
		}
		cs.bytesRemain -= int64(n)
		if err == io.EOF && cs.bytesRemain > 0 {
			err = io.ErrUnexpectedEOF
			cs.readErr = err
			return n, err
		}
	}
	if n == 0 {
		
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	var connAdd, streamAdd int32
	
	if v := cc.inflow.available(); v < transportDefaultConnFlow/2 {
		connAdd = transportDefaultConnFlow - v
		cc.inflow.add(connAdd)
	}
	if err == nil { 
		
		
		
		v := int(cs.inflow.available()) + cs.bufPipe.Len()
		if v < transportDefaultStreamFlow-transportDefaultStreamMinRefresh {
			streamAdd = int32(transportDefaultStreamFlow - v)
			cs.inflow.add(streamAdd)
		}
	}
	if connAdd != 0 || streamAdd != 0 {
		cc.wmu.Lock()
		defer cc.wmu.Unlock()
		if connAdd != 0 {
			cc.fr.WriteWindowUpdate(0, mustUint31(connAdd))
		}
		if streamAdd != 0 {
			cc.fr.WriteWindowUpdate(cs.ID, mustUint31(streamAdd))
		}
		cc.bw.Flush()
	}
	return
}

var errClosedResponseBody = errors.New("http2: response body closed")

func (b transportResponseBody) Close() error {
	cs := b.cs
	cc := cs.cc

	serverSentStreamEnd := cs.bufPipe.Err() == io.EOF
	unread := cs.bufPipe.Len()

	if unread > 0 || !serverSentStreamEnd {
		cc.mu.Lock()
		cc.wmu.Lock()
		if !serverSentStreamEnd {
			cc.fr.WriteRSTStream(cs.ID, ErrCodeCancel)
			cs.didReset = true
		}
		
		if unread > 0 {
			cc.inflow.add(int32(unread))
			cc.fr.WriteWindowUpdate(0, uint32(unread))
		}
		cc.bw.Flush()
		cc.wmu.Unlock()
		cc.mu.Unlock()
	}

	cs.bufPipe.BreakWithError(errClosedResponseBody)
	cc.forgetStreamID(cs.ID)
	return nil
}

func (rl *clientConnReadLoop) processData(f *DataFrame) error {
	cc := rl.cc
	cs := cc.streamByID(f.StreamID, f.StreamEnded())
	data := f.Data()
	if cs == nil {
		cc.mu.Lock()
		neverSent := cc.nextStreamID
		cc.mu.Unlock()
		if f.StreamID >= neverSent {
			
			cc.logf("http2: Transport received unsolicited DATA frame; closing connection")
			return ConnectionError(ErrCodeProtocol)
		}
		
		
		
		

		
		if f.Length > 0 {
			cc.mu.Lock()
			cc.inflow.add(int32(f.Length))
			cc.mu.Unlock()

			cc.wmu.Lock()
			cc.fr.WriteWindowUpdate(0, uint32(f.Length))
			cc.bw.Flush()
			cc.wmu.Unlock()
		}
		return nil
	}
	if !cs.firstByte {
		cc.logf("protocol error: received DATA before a HEADERS frame")
		rl.endStreamError(cs, StreamError{
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
		})
		return nil
	}
	if f.Length > 0 {
		if cs.req.Method == "HEAD" && len(data) > 0 {
			cc.logf("protocol error: received DATA on a HEAD request")
			rl.endStreamError(cs, StreamError{
				StreamID: f.StreamID,
				Code:     ErrCodeProtocol,
			})
			return nil
		}
		
		cc.mu.Lock()
		if cs.inflow.available() >= int32(f.Length) {
			cs.inflow.take(int32(f.Length))
		} else {
			cc.mu.Unlock()
			return ConnectionError(ErrCodeFlowControl)
		}
		
		
		var refund int
		if pad := int(f.Length) - len(data); pad > 0 {
			refund += pad
		}
		
		
		didReset := cs.didReset
		if didReset {
			refund += len(data)
		}
		if refund > 0 {
			cc.inflow.add(int32(refund))
			cc.wmu.Lock()
			cc.fr.WriteWindowUpdate(0, uint32(refund))
			if !didReset {
				cs.inflow.add(int32(refund))
				cc.fr.WriteWindowUpdate(cs.ID, uint32(refund))
			}
			cc.bw.Flush()
			cc.wmu.Unlock()
		}
		cc.mu.Unlock()

		if len(data) > 0 && !didReset {
			if _, err := cs.bufPipe.Write(data); err != nil {
				rl.endStreamError(cs, err)
				return err
			}
		}
	}

	if f.StreamEnded() {
		rl.endStream(cs)
	}
	return nil
}

var errInvalidTrailers = errors.New("http2: invalid trailers")

func (rl *clientConnReadLoop) endStream(cs *clientStream) {
	
	
	rl.endStreamError(cs, nil)
}

func (rl *clientConnReadLoop) endStreamError(cs *clientStream, err error) {
	var code func()
	if err == nil {
		err = io.EOF
		code = cs.copyTrailers
	}
	if isConnectionCloseRequest(cs.req) {
		rl.closeWhenIdle = true
	}
	cs.bufPipe.closeWithErrorAndCode(err, code)

	select {
	case cs.resc <- resAndError{err: err}:
	default:
	}
}

func (cs *clientStream) copyTrailers() {
	for k, vv := range cs.trailer {
		t := cs.resTrailer
		if *t == nil {
			*t = make(http.Header)
		}
		(*t)[k] = vv
	}
}

func (rl *clientConnReadLoop) processGoAway(f *GoAwayFrame) error {
	cc := rl.cc
	cc.t.connPool().MarkDead(cc)
	if f.ErrCode != 0 {
		
		cc.vlogf("transport got GOAWAY with error code = %v", f.ErrCode)
	}
	cc.setGoAway(f)
	return nil
}

func (rl *clientConnReadLoop) processSettings(f *SettingsFrame) error {
	cc := rl.cc
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if f.IsAck() {
		if cc.wantSettingsAck {
			cc.wantSettingsAck = false
			return nil
		}
		return ConnectionError(ErrCodeProtocol)
	}

	err := f.ForeachSetting(func(s Setting) error {
		switch s.ID {
		case SettingMaxFrameSize:
			cc.maxFrameSize = s.Val
		case SettingMaxConcurrentStreams:
			cc.maxConcurrentStreams = s.Val
		case SettingMaxHeaderListSize:
			cc.peerMaxHeaderListSize = uint64(s.Val)
		case SettingInitialWindowSize:
			
			
			
			
			if s.Val > math.MaxInt32 {
				return ConnectionError(ErrCodeFlowControl)
			}

			
			
			
			delta := int32(s.Val) - int32(cc.initialWindowSize)
			for _, cs := range cc.streams {
				cs.flow.add(delta)
			}
			cc.cond.Broadcast()

			cc.initialWindowSize = s.Val
		default:
			
			cc.vlogf("Unhandled Setting: %v", s)
		}
		return nil
	})
	if err != nil {
		return err
	}

	cc.wmu.Lock()
	defer cc.wmu.Unlock()

	cc.fr.WriteSettingsAck()
	cc.bw.Flush()
	return cc.werr
}

func (rl *clientConnReadLoop) processWindowUpdate(f *WindowUpdateFrame) error {
	cc := rl.cc
	cs := cc.streamByID(f.StreamID, false)
	if f.StreamID != 0 && cs == nil {
		return nil
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	fl := &cc.flow
	if cs != nil {
		fl = &cs.flow
	}
	if !fl.add(int32(f.Increment)) {
		return ConnectionError(ErrCodeFlowControl)
	}
	cc.cond.Broadcast()
	return nil
}

func (rl *clientConnReadLoop) processResetStream(f *RSTStreamFrame) error {
	cs := rl.cc.streamByID(f.StreamID, true)
	if cs == nil {
		
		return nil
	}
	select {
	case <-cs.peerReset:
		
		
		
		
	default:
		err := streamError(cs.ID, f.ErrCode)
		cs.resetErr = err
		close(cs.peerReset)
		cs.bufPipe.CloseWithError(err)
		cs.cc.cond.Broadcast() 
	}
	return nil
}



func (cc *ClientConn) ping(ctx contextContext) error {
	c := make(chan struct{})
	
	var p [8]byte
	for {
		if _, err := rand.Read(p[:]); err != nil {
			return err
		}
		cc.mu.Lock()
		
		if _, found := cc.pings[p]; !found {
			cc.pings[p] = c
			cc.mu.Unlock()
			break
		}
		cc.mu.Unlock()
	}
	cc.wmu.Lock()
	if err := cc.fr.WritePing(false, p); err != nil {
		cc.wmu.Unlock()
		return err
	}
	if err := cc.bw.Flush(); err != nil {
		cc.wmu.Unlock()
		return err
	}
	cc.wmu.Unlock()
	select {
	case <-c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-cc.readerDone:
		
		return cc.readerErr
	}
}

func (rl *clientConnReadLoop) processPing(f *PingFrame) error {
	if f.IsAck() {
		cc := rl.cc
		cc.mu.Lock()
		defer cc.mu.Unlock()
		
		if c, ok := cc.pings[f.Data]; ok {
			close(c)
			delete(cc.pings, f.Data)
		}
		return nil
	}
	cc := rl.cc
	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	if err := cc.fr.WritePing(true, f.Data); err != nil {
		return err
	}
	return cc.bw.Flush()
}

func (rl *clientConnReadLoop) processPushPromise(f *PushPromiseFrame) error {
	
	
	
	
	
	
	
	return ConnectionError(ErrCodeProtocol)
}

func (cc *ClientConn) writeStreamReset(streamID uint32, code ErrCode, err error) {
	
	
	
	
	cc.wmu.Lock()
	cc.fr.WriteRSTStream(streamID, code)
	cc.bw.Flush()
	cc.wmu.Unlock()
}

var (
	errResponseHeaderListSize = errors.New("http2: response header list larger than advertised limit")
	errRequestHeaderListSize  = errors.New("http2: request header list larger than peer's advertised limit")
	errPseudoTrailers         = errors.New("http2: invalid pseudo header in trailers")
)

func (cc *ClientConn) logf(format string, args ...interface{}) {
	cc.t.logf(format, args...)
}

func (cc *ClientConn) vlogf(format string, args ...interface{}) {
	cc.t.vlogf(format, args...)
}

func (t *Transport) vlogf(format string, args ...interface{}) {
	if VerboseLogs {
		t.logf(format, args...)
	}
}

func (t *Transport) logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

var noBody io.ReadCloser = ioutil.NopCloser(bytes.NewReader(nil))

func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

type erringRoundTripper struct{ err error }

func (rt erringRoundTripper) RoundTrip(*http.Request) (*http.Response, error) { return nil, rt.err }



type gzipReader struct {
	body io.ReadCloser 
	zr   *gzip.Reader  
	zerr error         
}

func (gz *gzipReader) Read(p []byte) (n int, err error) {
	if gz.zerr != nil {
		return 0, gz.zerr
	}
	if gz.zr == nil {
		gz.zr, err = gzip.NewReader(gz.body)
		if err != nil {
			gz.zerr = err
			return 0, err
		}
	}
	return gz.zr.Read(p)
}

func (gz *gzipReader) Close() error {
	return gz.body.Close()
}

type errorReader struct{ err error }

func (r errorReader) Read(p []byte) (int, error) { return 0, r.err }




type bodyWriterState struct {
	cs     *clientStream
	timer  *time.Timer   
	fnonce *sync.Once    
	fn     func()        
	resc   chan error    
	delay  time.Duration 
}

func (t *Transport) getBodyWriterState(cs *clientStream, body io.Reader) (s bodyWriterState) {
	s.cs = cs
	if body == nil {
		return
	}
	resc := make(chan error, 1)
	s.resc = resc
	s.fn = func() {
		cs.cc.mu.Lock()
		cs.startedWrite = true
		cs.cc.mu.Unlock()
		resc <- cs.writeRequestBody(body, cs.req.Body)
	}
	s.delay = t.expectContinueTimeout()
	if s.delay == 0 ||
		!httpguts.HeaderValuesContainsToken(
			cs.req.Header["Expect"],
			"100-continue") {
		return
	}
	s.fnonce = new(sync.Once)

	
	
	
	
	
	const hugeDuration = 365 * 24 * time.Hour
	s.timer = time.AfterFunc(hugeDuration, func() {
		s.fnonce.Do(s.fn)
	})
	return
}

func (s bodyWriterState) cancel() {
	if s.timer != nil {
		s.timer.Stop()
	}
}

func (s bodyWriterState) on100() {
	if s.timer == nil {
		
		
		return
	}
	s.timer.Stop()
	go func() { s.fnonce.Do(s.fn) }()
}




func (s bodyWriterState) scheduleBodyWrite() {
	if s.timer == nil {
		
		
		
		go s.fn()
		return
	}
	traceWait100Continue(s.cs.trace)
	if s.timer.Stop() {
		s.timer.Reset(s.delay)
	}
}



func isConnectionCloseRequest(req *http.Request) bool {
	return req.Close || httpguts.HeaderValuesContainsToken(req.Header["Connection"], "close")
}
