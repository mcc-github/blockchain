
























package http2

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http2/hpack"
)

const (
	prefaceTimeout        = 10 * time.Second
	firstSettingsTimeout  = 2 * time.Second 
	handlerChunkWriteSize = 4 << 10
	defaultMaxStreams     = 250 
)

var (
	errClientDisconnected = errors.New("client disconnected")
	errClosedBody         = errors.New("body closed by handler")
	errHandlerComplete    = errors.New("http2: request body closed due to handler exiting")
	errStreamClosed       = errors.New("http2: stream closed")
)

var responseWriterStatePool = sync.Pool{
	New: func() interface{} {
		rws := &responseWriterState{}
		rws.bw = bufio.NewWriterSize(chunkWriter{rws}, handlerChunkWriteSize)
		return rws
	},
}


var (
	testHookOnConn        func()
	testHookGetServerConn func(*serverConn)
	testHookOnPanicMu     *sync.Mutex 
	testHookOnPanic       func(sc *serverConn, panicVal interface{}) (rePanic bool)
)


type Server struct {
	
	
	
	
	MaxHandlers int

	
	
	
	
	
	
	MaxConcurrentStreams uint32

	
	
	
	
	MaxReadFrameSize uint32

	
	
	PermitProhibitedCipherSuites bool

	
	
	
	IdleTimeout time.Duration

	
	
	
	
	
	MaxUploadBufferPerConnection int32

	
	
	
	
	MaxUploadBufferPerStream int32

	
	
	NewWriteScheduler func() WriteScheduler

	
	
	
	state *serverInternalState
}

func (s *Server) initialConnRecvWindowSize() int32 {
	if s.MaxUploadBufferPerConnection > initialWindowSize {
		return s.MaxUploadBufferPerConnection
	}
	return 1 << 20
}

func (s *Server) initialStreamRecvWindowSize() int32 {
	if s.MaxUploadBufferPerStream > 0 {
		return s.MaxUploadBufferPerStream
	}
	return 1 << 20
}

func (s *Server) maxReadFrameSize() uint32 {
	if v := s.MaxReadFrameSize; v >= minMaxFrameSize && v <= maxFrameSize {
		return v
	}
	return defaultMaxReadFrameSize
}

func (s *Server) maxConcurrentStreams() uint32 {
	if v := s.MaxConcurrentStreams; v > 0 {
		return v
	}
	return defaultMaxStreams
}

type serverInternalState struct {
	mu          sync.Mutex
	activeConns map[*serverConn]struct{}
}

func (s *serverInternalState) registerConn(sc *serverConn) {
	if s == nil {
		return 
	}
	s.mu.Lock()
	s.activeConns[sc] = struct{}{}
	s.mu.Unlock()
}

func (s *serverInternalState) unregisterConn(sc *serverConn) {
	if s == nil {
		return 
	}
	s.mu.Lock()
	delete(s.activeConns, sc)
	s.mu.Unlock()
}

func (s *serverInternalState) startGracefulShutdown() {
	if s == nil {
		return 
	}
	s.mu.Lock()
	for sc := range s.activeConns {
		sc.startGracefulShutdown()
	}
	s.mu.Unlock()
}






func ConfigureServer(s *http.Server, conf *Server) error {
	if s == nil {
		panic("nil *http.Server")
	}
	if conf == nil {
		conf = new(Server)
	}
	conf.state = &serverInternalState{activeConns: make(map[*serverConn]struct{})}
	if err := configureServer18(s, conf); err != nil {
		return err
	}
	if err := configureServer19(s, conf); err != nil {
		return err
	}

	if s.TLSConfig == nil {
		s.TLSConfig = new(tls.Config)
	} else if s.TLSConfig.CipherSuites != nil {
		
		
		
		haveRequired := false
		sawBad := false
		for i, cs := range s.TLSConfig.CipherSuites {
			switch cs {
			case tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				
				
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:
				haveRequired = true
			}
			if isBadCipher(cs) {
				sawBad = true
			} else if sawBad {
				return fmt.Errorf("http2: TLSConfig.CipherSuites index %d contains an HTTP/2-approved cipher suite (%#04x), but it comes after unapproved cipher suites. With this configuration, clients that don't support previous, approved cipher suites may be given an unapproved one and reject the connection.", i, cs)
			}
		}
		if !haveRequired {
			return fmt.Errorf("http2: TLSConfig.CipherSuites is missing an HTTP/2-required AES_128_GCM_SHA256 cipher.")
		}
	}

	
	
	
	
	
	

	s.TLSConfig.PreferServerCipherSuites = true

	haveNPN := false
	for _, p := range s.TLSConfig.NextProtos {
		if p == NextProtoTLS {
			haveNPN = true
			break
		}
	}
	if !haveNPN {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, NextProtoTLS)
	}

	if s.TLSNextProto == nil {
		s.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler){}
	}
	protoHandler := func(hs *http.Server, c *tls.Conn, h http.Handler) {
		if testHookOnConn != nil {
			testHookOnConn()
		}
		conf.ServeConn(c, &ServeConnOpts{
			Handler:    h,
			BaseConfig: hs,
		})
	}
	s.TLSNextProto[NextProtoTLS] = protoHandler
	return nil
}


type ServeConnOpts struct {
	
	
	BaseConfig *http.Server

	
	
	
	Handler http.Handler
}

func (o *ServeConnOpts) baseConfig() *http.Server {
	if o != nil && o.BaseConfig != nil {
		return o.BaseConfig
	}
	return new(http.Server)
}

func (o *ServeConnOpts) handler() http.Handler {
	if o != nil {
		if o.Handler != nil {
			return o.Handler
		}
		if o.BaseConfig != nil && o.BaseConfig.Handler != nil {
			return o.BaseConfig.Handler
		}
	}
	return http.DefaultServeMux
}















func (s *Server) ServeConn(c net.Conn, opts *ServeConnOpts) {
	baseCtx, cancel := serverConnBaseContext(c, opts)
	defer cancel()

	sc := &serverConn{
		srv:                         s,
		hs:                          opts.baseConfig(),
		conn:                        c,
		baseCtx:                     baseCtx,
		remoteAddrStr:               c.RemoteAddr().String(),
		bw:                          newBufferedWriter(c),
		handler:                     opts.handler(),
		streams:                     make(map[uint32]*stream),
		readFrameCh:                 make(chan readFrameResult),
		wantWriteFrameCh:            make(chan FrameWriteRequest, 8),
		serveMsgCh:                  make(chan interface{}, 8),
		wroteFrameCh:                make(chan frameWriteResult, 1), 
		bodyReadCh:                  make(chan bodyReadMsg),         
		doneServing:                 make(chan struct{}),
		clientMaxStreams:            math.MaxUint32, 
		advMaxStreams:               s.maxConcurrentStreams(),
		initialStreamSendWindowSize: initialWindowSize,
		maxFrameSize:                initialMaxFrameSize,
		headerTableSize:             initialHeaderTableSize,
		serveG:                      newGoroutineLock(),
		pushEnabled:                 true,
	}

	s.state.registerConn(sc)
	defer s.state.unregisterConn(sc)

	
	
	
	
	
	if sc.hs.WriteTimeout != 0 {
		sc.conn.SetWriteDeadline(time.Time{})
	}

	if s.NewWriteScheduler != nil {
		sc.writeSched = s.NewWriteScheduler()
	} else {
		sc.writeSched = NewRandomWriteScheduler()
	}

	
	
	
	sc.flow.add(initialWindowSize)
	sc.inflow.add(initialWindowSize)
	sc.hpackEncoder = hpack.NewEncoder(&sc.headerWriteBuf)

	fr := NewFramer(sc.bw, c)
	fr.ReadMetaHeaders = hpack.NewDecoder(initialHeaderTableSize, nil)
	fr.MaxHeaderListSize = sc.maxHeaderListSize()
	fr.SetMaxReadFrameSize(s.maxReadFrameSize())
	sc.framer = fr

	if tc, ok := c.(connectionStater); ok {
		sc.tlsState = new(tls.ConnectionState)
		*sc.tlsState = tc.ConnectionState()
		
		
		
		
		
		
		
		
		
		
		if sc.tlsState.Version < tls.VersionTLS12 {
			sc.rejectConn(ErrCodeInadequateSecurity, "TLS version too low")
			return
		}

		if sc.tlsState.ServerName == "" {
			
			
			
			
			
			
			
			
			
		}

		if !s.PermitProhibitedCipherSuites && isBadCipher(sc.tlsState.CipherSuite) {
			
			
			
			
			
			
			
			
			
			
			sc.rejectConn(ErrCodeInadequateSecurity, fmt.Sprintf("Prohibited TLS 1.2 Cipher Suite: %x", sc.tlsState.CipherSuite))
			return
		}
	}

	if hook := testHookGetServerConn; hook != nil {
		hook(sc)
	}
	sc.serve()
}

func (sc *serverConn) rejectConn(err ErrCode, debug string) {
	sc.vlogf("http2: server rejecting conn: %v, %s", err, debug)
	
	sc.framer.WriteGoAway(0, err, []byte(debug))
	sc.bw.Flush()
	sc.conn.Close()
}

type serverConn struct {
	
	srv              *Server
	hs               *http.Server
	conn             net.Conn
	bw               *bufferedWriter 
	handler          http.Handler
	baseCtx          contextContext
	framer           *Framer
	doneServing      chan struct{}          
	readFrameCh      chan readFrameResult   
	wantWriteFrameCh chan FrameWriteRequest 
	wroteFrameCh     chan frameWriteResult  
	bodyReadCh       chan bodyReadMsg       
	serveMsgCh       chan interface{}       
	flow             flow                   
	inflow           flow                   
	tlsState         *tls.ConnectionState   
	remoteAddrStr    string
	writeSched       WriteScheduler

	
	serveG                      goroutineLock 
	pushEnabled                 bool
	sawFirstSettings            bool 
	needToSendSettingsAck       bool
	unackedSettings             int    
	clientMaxStreams            uint32 
	advMaxStreams               uint32 
	curClientStreams            uint32 
	curPushedStreams            uint32 
	maxClientStreamID           uint32 
	maxPushPromiseID            uint32 
	streams                     map[uint32]*stream
	initialStreamSendWindowSize int32
	maxFrameSize                int32
	headerTableSize             uint32
	peerMaxHeaderListSize       uint32            
	canonHeader                 map[string]string 
	writingFrame                bool              
	writingFrameAsync           bool              
	needsFrameFlush             bool              
	inGoAway                    bool              
	inFrameScheduleLoop         bool              
	needToSendGoAway            bool              
	goAwayCode                  ErrCode
	shutdownTimer               *time.Timer 
	idleTimer                   *time.Timer 

	
	headerWriteBuf bytes.Buffer
	hpackEncoder   *hpack.Encoder

	
	shutdownOnce sync.Once
}

func (sc *serverConn) maxHeaderListSize() uint32 {
	n := sc.hs.MaxHeaderBytes
	if n <= 0 {
		n = http.DefaultMaxHeaderBytes
	}
	
	
	const perFieldOverhead = 32 
	const typicalHeaders = 10   
	return uint32(n + typicalHeaders*perFieldOverhead)
}

func (sc *serverConn) curOpenStreams() uint32 {
	sc.serveG.check()
	return sc.curClientStreams + sc.curPushedStreams
}








type stream struct {
	
	sc        *serverConn
	id        uint32
	body      *pipe       
	cw        closeWaiter 
	ctx       contextContext
	cancelCtx func()

	
	bodyBytes        int64   
	declBodyBytes    int64   
	flow             flow    
	inflow           flow    
	parent           *stream 
	numTrailerValues int64
	weight           uint8
	state            streamState
	resetQueued      bool        
	gotTrailerHeader bool        
	wroteHeaders     bool        
	writeDeadline    *time.Timer 

	trailer    http.Header 
	reqTrailer http.Header 
}

func (sc *serverConn) Framer() *Framer  { return sc.framer }
func (sc *serverConn) CloseConn() error { return sc.conn.Close() }
func (sc *serverConn) Flush() error     { return sc.bw.Flush() }
func (sc *serverConn) HeaderEncoder() (*hpack.Encoder, *bytes.Buffer) {
	return sc.hpackEncoder, &sc.headerWriteBuf
}

func (sc *serverConn) state(streamID uint32) (streamState, *stream) {
	sc.serveG.check()
	
	if st, ok := sc.streams[streamID]; ok {
		return st.state, st
	}
	
	
	
	
	
	
	if streamID%2 == 1 {
		if streamID <= sc.maxClientStreamID {
			return stateClosed, nil
		}
	} else {
		if streamID <= sc.maxPushPromiseID {
			return stateClosed, nil
		}
	}
	return stateIdle, nil
}




func (sc *serverConn) setConnState(state http.ConnState) {
	if sc.hs.ConnState != nil {
		sc.hs.ConnState(sc.conn, state)
	}
}

func (sc *serverConn) vlogf(format string, args ...interface{}) {
	if VerboseLogs {
		sc.logf(format, args...)
	}
}

func (sc *serverConn) logf(format string, args ...interface{}) {
	if lg := sc.hs.ErrorLog; lg != nil {
		lg.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}





func errno(v error) uintptr {
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Uintptr {
		return uintptr(rv.Uint())
	}
	return 0
}



func isClosedConnError(err error) bool {
	if err == nil {
		return false
	}

	
	
	
	str := err.Error()
	if strings.Contains(str, "use of closed network connection") {
		return true
	}

	
	
	
	
	if runtime.GOOS == "windows" {
		if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
			if se, ok := oe.Err.(*os.SyscallError); ok && se.Syscall == "wsarecv" {
				const WSAECONNABORTED = 10053
				const WSAECONNRESET = 10054
				if n := errno(se.Err); n == WSAECONNRESET || n == WSAECONNABORTED {
					return true
				}
			}
		}
	}
	return false
}

func (sc *serverConn) condlogf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF || isClosedConnError(err) || err == errPrefaceTimeout {
		
		sc.vlogf(format, args...)
	} else {
		sc.logf(format, args...)
	}
}

func (sc *serverConn) canonicalHeader(v string) string {
	sc.serveG.check()
	cv, ok := commonCanonHeader[v]
	if ok {
		return cv
	}
	cv, ok = sc.canonHeader[v]
	if ok {
		return cv
	}
	if sc.canonHeader == nil {
		sc.canonHeader = make(map[string]string)
	}
	cv = http.CanonicalHeaderKey(v)
	sc.canonHeader[v] = cv
	return cv
}

type readFrameResult struct {
	f   Frame 
	err error

	
	
	
	readMore func()
}





func (sc *serverConn) readFrames() {
	gate := make(gate)
	gateDone := gate.Done
	for {
		f, err := sc.framer.ReadFrame()
		select {
		case sc.readFrameCh <- readFrameResult{f, err, gateDone}:
		case <-sc.doneServing:
			return
		}
		select {
		case <-gate:
		case <-sc.doneServing:
			return
		}
		if terminalReadFrameError(err) {
			return
		}
	}
}


type frameWriteResult struct {
	wr  FrameWriteRequest 
	err error             
}





func (sc *serverConn) writeFrameAsync(wr FrameWriteRequest) {
	err := wr.write.writeFrame(sc)
	sc.wroteFrameCh <- frameWriteResult{wr, err}
}

func (sc *serverConn) closeAllStreamsOnConnClose() {
	sc.serveG.check()
	for _, st := range sc.streams {
		sc.closeStream(st, errClientDisconnected)
	}
}

func (sc *serverConn) stopShutdownTimer() {
	sc.serveG.check()
	if t := sc.shutdownTimer; t != nil {
		t.Stop()
	}
}

func (sc *serverConn) notePanic() {
	
	if testHookOnPanicMu != nil {
		testHookOnPanicMu.Lock()
		defer testHookOnPanicMu.Unlock()
	}
	if testHookOnPanic != nil {
		if e := recover(); e != nil {
			if testHookOnPanic(sc, e) {
				panic(e)
			}
		}
	}
}

func (sc *serverConn) serve() {
	sc.serveG.check()
	defer sc.notePanic()
	defer sc.conn.Close()
	defer sc.closeAllStreamsOnConnClose()
	defer sc.stopShutdownTimer()
	defer close(sc.doneServing) 

	if VerboseLogs {
		sc.vlogf("http2: server connection from %v on %p", sc.conn.RemoteAddr(), sc.hs)
	}

	sc.writeFrame(FrameWriteRequest{
		write: writeSettings{
			{SettingMaxFrameSize, sc.srv.maxReadFrameSize()},
			{SettingMaxConcurrentStreams, sc.advMaxStreams},
			{SettingMaxHeaderListSize, sc.maxHeaderListSize()},
			{SettingInitialWindowSize, uint32(sc.srv.initialStreamRecvWindowSize())},
		},
	})
	sc.unackedSettings++

	
	
	if diff := sc.srv.initialConnRecvWindowSize() - initialWindowSize; diff > 0 {
		sc.sendWindowUpdate(nil, int(diff))
	}

	if err := sc.readPreface(); err != nil {
		sc.condlogf(err, "http2: server: error reading preface from client %v: %v", sc.conn.RemoteAddr(), err)
		return
	}
	
	
	
	
	sc.setConnState(http.StateActive)
	sc.setConnState(http.StateIdle)

	if sc.srv.IdleTimeout != 0 {
		sc.idleTimer = time.AfterFunc(sc.srv.IdleTimeout, sc.onIdleTimer)
		defer sc.idleTimer.Stop()
	}

	go sc.readFrames() 

	settingsTimer := time.AfterFunc(firstSettingsTimeout, sc.onSettingsTimer)
	defer settingsTimer.Stop()

	loopNum := 0
	for {
		loopNum++
		select {
		case wr := <-sc.wantWriteFrameCh:
			if se, ok := wr.write.(StreamError); ok {
				sc.resetStream(se)
				break
			}
			sc.writeFrame(wr)
		case res := <-sc.wroteFrameCh:
			sc.wroteFrame(res)
		case res := <-sc.readFrameCh:
			if !sc.processFrameFromReader(res) {
				return
			}
			res.readMore()
			if settingsTimer != nil {
				settingsTimer.Stop()
				settingsTimer = nil
			}
		case m := <-sc.bodyReadCh:
			sc.noteBodyRead(m.st, m.n)
		case msg := <-sc.serveMsgCh:
			switch v := msg.(type) {
			case func(int):
				v(loopNum) 
			case *serverMessage:
				switch v {
				case settingsTimerMsg:
					sc.logf("timeout waiting for SETTINGS frames from %v", sc.conn.RemoteAddr())
					return
				case idleTimerMsg:
					sc.vlogf("connection is idle")
					sc.goAway(ErrCodeNo)
				case shutdownTimerMsg:
					sc.vlogf("GOAWAY close timer fired; closing conn from %v", sc.conn.RemoteAddr())
					return
				case gracefulShutdownMsg:
					sc.startGracefulShutdownInternal()
				default:
					panic("unknown timer")
				}
			case *startPushRequest:
				sc.startPush(v)
			default:
				panic(fmt.Sprintf("unexpected type %T", v))
			}
		}

		
		
		
		sentGoAway := sc.inGoAway && !sc.needToSendGoAway && !sc.writingFrame
		gracefulShutdownComplete := sc.goAwayCode == ErrCodeNo && sc.curOpenStreams() == 0
		if sentGoAway && sc.shutdownTimer == nil && (sc.goAwayCode != ErrCodeNo || gracefulShutdownComplete) {
			sc.shutDownIn(goAwayTimeout)
		}
	}
}

func (sc *serverConn) awaitGracefulShutdown(sharedCh <-chan struct{}, privateCh chan struct{}) {
	select {
	case <-sc.doneServing:
	case <-sharedCh:
		close(privateCh)
	}
}

type serverMessage int


var (
	settingsTimerMsg    = new(serverMessage)
	idleTimerMsg        = new(serverMessage)
	shutdownTimerMsg    = new(serverMessage)
	gracefulShutdownMsg = new(serverMessage)
)

func (sc *serverConn) onSettingsTimer() { sc.sendServeMsg(settingsTimerMsg) }
func (sc *serverConn) onIdleTimer()     { sc.sendServeMsg(idleTimerMsg) }
func (sc *serverConn) onShutdownTimer() { sc.sendServeMsg(shutdownTimerMsg) }

func (sc *serverConn) sendServeMsg(msg interface{}) {
	sc.serveG.checkNotOn() 
	select {
	case sc.serveMsgCh <- msg:
	case <-sc.doneServing:
	}
}

var errPrefaceTimeout = errors.New("timeout waiting for client preface")




func (sc *serverConn) readPreface() error {
	errc := make(chan error, 1)
	go func() {
		
		buf := make([]byte, len(ClientPreface))
		if _, err := io.ReadFull(sc.conn, buf); err != nil {
			errc <- err
		} else if !bytes.Equal(buf, clientPreface) {
			errc <- fmt.Errorf("bogus greeting %q", buf)
		} else {
			errc <- nil
		}
	}()
	timer := time.NewTimer(prefaceTimeout) 
	defer timer.Stop()
	select {
	case <-timer.C:
		return errPrefaceTimeout
	case err := <-errc:
		if err == nil {
			if VerboseLogs {
				sc.vlogf("http2: server: client %v said hello", sc.conn.RemoteAddr())
			}
		}
		return err
	}
}

var errChanPool = sync.Pool{
	New: func() interface{} { return make(chan error, 1) },
}

var writeDataPool = sync.Pool{
	New: func() interface{} { return new(writeData) },
}



func (sc *serverConn) writeDataFromHandler(stream *stream, data []byte, endStream bool) error {
	ch := errChanPool.Get().(chan error)
	writeArg := writeDataPool.Get().(*writeData)
	*writeArg = writeData{stream.id, data, endStream}
	err := sc.writeFrameFromHandler(FrameWriteRequest{
		write:  writeArg,
		stream: stream,
		done:   ch,
	})
	if err != nil {
		return err
	}
	var frameWriteDone bool 
	select {
	case err = <-ch:
		frameWriteDone = true
	case <-sc.doneServing:
		return errClientDisconnected
	case <-stream.cw:
		
		
		
		
		
		
		
		select {
		case err = <-ch:
			frameWriteDone = true
		default:
			return errStreamClosed
		}
	}
	errChanPool.Put(ch)
	if frameWriteDone {
		writeDataPool.Put(writeArg)
	}
	return err
}








func (sc *serverConn) writeFrameFromHandler(wr FrameWriteRequest) error {
	sc.serveG.checkNotOn() 
	select {
	case sc.wantWriteFrameCh <- wr:
		return nil
	case <-sc.doneServing:
		
		
		return errClientDisconnected
	}
}









func (sc *serverConn) writeFrame(wr FrameWriteRequest) {
	sc.serveG.check()

	
	var ignoreWrite bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	if wr.StreamID() != 0 {
		_, isReset := wr.write.(StreamError)
		if state, _ := sc.state(wr.StreamID()); state == stateClosed && !isReset {
			ignoreWrite = true
		}
	}

	
	
	switch wr.write.(type) {
	case *writeResHeaders:
		wr.stream.wroteHeaders = true
	case write100ContinueHeadersFrame:
		if wr.stream.wroteHeaders {
			
			
			if wr.done != nil {
				panic("wr.done != nil for write100ContinueHeadersFrame")
			}
			ignoreWrite = true
		}
	}

	if !ignoreWrite {
		sc.writeSched.Push(wr)
	}
	sc.scheduleFrameWrite()
}




func (sc *serverConn) startFrameWrite(wr FrameWriteRequest) {
	sc.serveG.check()
	if sc.writingFrame {
		panic("internal error: can only be writing one frame at a time")
	}

	st := wr.stream
	if st != nil {
		switch st.state {
		case stateHalfClosedLocal:
			switch wr.write.(type) {
			case StreamError, handlerPanicRST, writeWindowUpdate:
				
				
			default:
				panic(fmt.Sprintf("internal error: attempt to send frame on a half-closed-local stream: %v", wr))
			}
		case stateClosed:
			panic(fmt.Sprintf("internal error: attempt to send frame on a closed stream: %v", wr))
		}
	}
	if wpp, ok := wr.write.(*writePushPromise); ok {
		var err error
		wpp.promisedID, err = wpp.allocatePromisedID()
		if err != nil {
			sc.writingFrameAsync = false
			wr.replyToWriter(err)
			return
		}
	}

	sc.writingFrame = true
	sc.needsFrameFlush = true
	if wr.write.staysWithinBuffer(sc.bw.Available()) {
		sc.writingFrameAsync = false
		err := wr.write.writeFrame(sc)
		sc.wroteFrame(frameWriteResult{wr, err})
	} else {
		sc.writingFrameAsync = true
		go sc.writeFrameAsync(wr)
	}
}




var errHandlerPanicked = errors.New("http2: handler panicked")



func (sc *serverConn) wroteFrame(res frameWriteResult) {
	sc.serveG.check()
	if !sc.writingFrame {
		panic("internal error: expected to be already writing a frame")
	}
	sc.writingFrame = false
	sc.writingFrameAsync = false

	wr := res.wr

	if writeEndsStream(wr.write) {
		st := wr.stream
		if st == nil {
			panic("internal error: expecting non-nil stream")
		}
		switch st.state {
		case stateOpen:
			
			
			
			
			
			
			
			
			
			
			st.state = stateHalfClosedLocal
			
			
			
			
			sc.resetStream(streamError(st.id, ErrCodeNo))
		case stateHalfClosedRemote:
			sc.closeStream(st, errHandlerComplete)
		}
	} else {
		switch v := wr.write.(type) {
		case StreamError:
			
			if st, ok := sc.streams[v.StreamID]; ok {
				sc.closeStream(st, v)
			}
		case handlerPanicRST:
			sc.closeStream(wr.stream, errHandlerPanicked)
		}
	}

	
	wr.replyToWriter(res.err)

	sc.scheduleFrameWrite()
}













func (sc *serverConn) scheduleFrameWrite() {
	sc.serveG.check()
	if sc.writingFrame || sc.inFrameScheduleLoop {
		return
	}
	sc.inFrameScheduleLoop = true
	for !sc.writingFrameAsync {
		if sc.needToSendGoAway {
			sc.needToSendGoAway = false
			sc.startFrameWrite(FrameWriteRequest{
				write: &writeGoAway{
					maxStreamID: sc.maxClientStreamID,
					code:        sc.goAwayCode,
				},
			})
			continue
		}
		if sc.needToSendSettingsAck {
			sc.needToSendSettingsAck = false
			sc.startFrameWrite(FrameWriteRequest{write: writeSettingsAck{}})
			continue
		}
		if !sc.inGoAway || sc.goAwayCode == ErrCodeNo {
			if wr, ok := sc.writeSched.Pop(); ok {
				sc.startFrameWrite(wr)
				continue
			}
		}
		if sc.needsFrameFlush {
			sc.startFrameWrite(FrameWriteRequest{write: flushFrameWriter{}})
			sc.needsFrameFlush = false 
			continue
		}
		break
	}
	sc.inFrameScheduleLoop = false
}








func (sc *serverConn) startGracefulShutdown() {
	sc.serveG.checkNotOn() 
	sc.shutdownOnce.Do(func() { sc.sendServeMsg(gracefulShutdownMsg) })
}















var goAwayTimeout = 1 * time.Second

func (sc *serverConn) startGracefulShutdownInternal() {
	sc.goAway(ErrCodeNo)
}

func (sc *serverConn) goAway(code ErrCode) {
	sc.serveG.check()
	if sc.inGoAway {
		return
	}
	sc.inGoAway = true
	sc.needToSendGoAway = true
	sc.goAwayCode = code
	sc.scheduleFrameWrite()
}

func (sc *serverConn) shutDownIn(d time.Duration) {
	sc.serveG.check()
	sc.shutdownTimer = time.AfterFunc(d, sc.onShutdownTimer)
}

func (sc *serverConn) resetStream(se StreamError) {
	sc.serveG.check()
	sc.writeFrame(FrameWriteRequest{write: se})
	if st, ok := sc.streams[se.StreamID]; ok {
		st.resetQueued = true
	}
}




func (sc *serverConn) processFrameFromReader(res readFrameResult) bool {
	sc.serveG.check()
	err := res.err
	if err != nil {
		if err == ErrFrameTooLarge {
			sc.goAway(ErrCodeFrameSize)
			return true 
		}
		clientGone := err == io.EOF || err == io.ErrUnexpectedEOF || isClosedConnError(err)
		if clientGone {
			
			
			
			
			
			
			
			
			return false
		}
	} else {
		f := res.f
		if VerboseLogs {
			sc.vlogf("http2: server read frame %v", summarizeFrame(f))
		}
		err = sc.processFrame(f)
		if err == nil {
			return true
		}
	}

	switch ev := err.(type) {
	case StreamError:
		sc.resetStream(ev)
		return true
	case goAwayFlowError:
		sc.goAway(ErrCodeFlowControl)
		return true
	case ConnectionError:
		sc.logf("http2: server connection error from %v: %v", sc.conn.RemoteAddr(), ev)
		sc.goAway(ErrCode(ev))
		return true 
	default:
		if res.err != nil {
			sc.vlogf("http2: server closing client connection; error reading frame from client %s: %v", sc.conn.RemoteAddr(), err)
		} else {
			sc.logf("http2: server closing client connection: %v", err)
		}
		return false
	}
}

func (sc *serverConn) processFrame(f Frame) error {
	sc.serveG.check()

	
	if !sc.sawFirstSettings {
		if _, ok := f.(*SettingsFrame); !ok {
			return ConnectionError(ErrCodeProtocol)
		}
		sc.sawFirstSettings = true
	}

	switch f := f.(type) {
	case *SettingsFrame:
		return sc.processSettings(f)
	case *MetaHeadersFrame:
		return sc.processHeaders(f)
	case *WindowUpdateFrame:
		return sc.processWindowUpdate(f)
	case *PingFrame:
		return sc.processPing(f)
	case *DataFrame:
		return sc.processData(f)
	case *RSTStreamFrame:
		return sc.processResetStream(f)
	case *PriorityFrame:
		return sc.processPriority(f)
	case *GoAwayFrame:
		return sc.processGoAway(f)
	case *PushPromiseFrame:
		
		
		return ConnectionError(ErrCodeProtocol)
	default:
		sc.vlogf("http2: server ignoring frame: %v", f.Header())
		return nil
	}
}

func (sc *serverConn) processPing(f *PingFrame) error {
	sc.serveG.check()
	if f.IsAck() {
		
		
		return nil
	}
	if f.StreamID != 0 {
		
		
		
		
		
		return ConnectionError(ErrCodeProtocol)
	}
	if sc.inGoAway && sc.goAwayCode != ErrCodeNo {
		return nil
	}
	sc.writeFrame(FrameWriteRequest{write: writePingAck{f}})
	return nil
}

func (sc *serverConn) processWindowUpdate(f *WindowUpdateFrame) error {
	sc.serveG.check()
	switch {
	case f.StreamID != 0: 
		state, st := sc.state(f.StreamID)
		if state == stateIdle {
			
			
			
			
			return ConnectionError(ErrCodeProtocol)
		}
		if st == nil {
			
			
			
			
			
			return nil
		}
		if !st.flow.add(int32(f.Increment)) {
			return streamError(f.StreamID, ErrCodeFlowControl)
		}
	default: 
		if !sc.flow.add(int32(f.Increment)) {
			return goAwayFlowError{}
		}
	}
	sc.scheduleFrameWrite()
	return nil
}

func (sc *serverConn) processResetStream(f *RSTStreamFrame) error {
	sc.serveG.check()

	state, st := sc.state(f.StreamID)
	if state == stateIdle {
		
		
		
		
		
		return ConnectionError(ErrCodeProtocol)
	}
	if st != nil {
		st.cancelCtx()
		sc.closeStream(st, streamError(f.StreamID, f.ErrCode))
	}
	return nil
}

func (sc *serverConn) closeStream(st *stream, err error) {
	sc.serveG.check()
	if st.state == stateIdle || st.state == stateClosed {
		panic(fmt.Sprintf("invariant; can't close stream in state %v", st.state))
	}
	st.state = stateClosed
	if st.writeDeadline != nil {
		st.writeDeadline.Stop()
	}
	if st.isPushed() {
		sc.curPushedStreams--
	} else {
		sc.curClientStreams--
	}
	delete(sc.streams, st.id)
	if len(sc.streams) == 0 {
		sc.setConnState(http.StateIdle)
		if sc.srv.IdleTimeout != 0 {
			sc.idleTimer.Reset(sc.srv.IdleTimeout)
		}
		if h1ServerKeepAlivesDisabled(sc.hs) {
			sc.startGracefulShutdownInternal()
		}
	}
	if p := st.body; p != nil {
		
		
		sc.sendWindowUpdate(nil, p.Len())

		p.CloseWithError(err)
	}
	st.cw.Close() 
	sc.writeSched.CloseStream(st.id)
}

func (sc *serverConn) processSettings(f *SettingsFrame) error {
	sc.serveG.check()
	if f.IsAck() {
		sc.unackedSettings--
		if sc.unackedSettings < 0 {
			
			
			
			return ConnectionError(ErrCodeProtocol)
		}
		return nil
	}
	if err := f.ForeachSetting(sc.processSetting); err != nil {
		return err
	}
	sc.needToSendSettingsAck = true
	sc.scheduleFrameWrite()
	return nil
}

func (sc *serverConn) processSetting(s Setting) error {
	sc.serveG.check()
	if err := s.Valid(); err != nil {
		return err
	}
	if VerboseLogs {
		sc.vlogf("http2: server processing setting %v", s)
	}
	switch s.ID {
	case SettingHeaderTableSize:
		sc.headerTableSize = s.Val
		sc.hpackEncoder.SetMaxDynamicTableSize(s.Val)
	case SettingEnablePush:
		sc.pushEnabled = s.Val != 0
	case SettingMaxConcurrentStreams:
		sc.clientMaxStreams = s.Val
	case SettingInitialWindowSize:
		return sc.processSettingInitialWindowSize(s.Val)
	case SettingMaxFrameSize:
		sc.maxFrameSize = int32(s.Val) 
	case SettingMaxHeaderListSize:
		sc.peerMaxHeaderListSize = s.Val
	default:
		
		
		
		if VerboseLogs {
			sc.vlogf("http2: server ignoring unknown setting %v", s)
		}
	}
	return nil
}

func (sc *serverConn) processSettingInitialWindowSize(val uint32) error {
	sc.serveG.check()
	
	

	
	
	
	
	
	
	old := sc.initialStreamSendWindowSize
	sc.initialStreamSendWindowSize = int32(val)
	growth := int32(val) - old 
	for _, st := range sc.streams {
		if !st.flow.add(growth) {
			
			
			
			
			
			
			return ConnectionError(ErrCodeFlowControl)
		}
	}
	return nil
}

func (sc *serverConn) processData(f *DataFrame) error {
	sc.serveG.check()
	if sc.inGoAway && sc.goAwayCode != ErrCodeNo {
		return nil
	}
	data := f.Data()

	
	
	
	id := f.Header().StreamID
	state, st := sc.state(id)
	if id == 0 || state == stateIdle {
		
		
		
		
		return ConnectionError(ErrCodeProtocol)
	}
	if st == nil || state != stateOpen || st.gotTrailerHeader || st.resetQueued {
		
		
		
		
		

		
		
		
		if sc.inflow.available() < int32(f.Length) {
			return streamError(id, ErrCodeFlowControl)
		}
		
		
		
		
		sc.inflow.take(int32(f.Length))
		sc.sendWindowUpdate(nil, int(f.Length)) 

		if st != nil && st.resetQueued {
			
			return nil
		}
		return streamError(id, ErrCodeStreamClosed)
	}
	if st.body == nil {
		panic("internal error: should have a body in this state")
	}

	
	if st.declBodyBytes != -1 && st.bodyBytes+int64(len(data)) > st.declBodyBytes {
		st.body.CloseWithError(fmt.Errorf("sender tried to send more than declared Content-Length of %d bytes", st.declBodyBytes))
		return streamError(id, ErrCodeStreamClosed)
	}
	if f.Length > 0 {
		
		if st.inflow.available() < int32(f.Length) {
			return streamError(id, ErrCodeFlowControl)
		}
		st.inflow.take(int32(f.Length))

		if len(data) > 0 {
			wrote, err := st.body.Write(data)
			if err != nil {
				return streamError(id, ErrCodeStreamClosed)
			}
			if wrote != len(data) {
				panic("internal error: bad Writer")
			}
			st.bodyBytes += int64(len(data))
		}

		
		
		if pad := int32(f.Length) - int32(len(data)); pad > 0 {
			sc.sendWindowUpdate32(nil, pad)
			sc.sendWindowUpdate32(st, pad)
		}
	}
	if f.StreamEnded() {
		st.endStream()
	}
	return nil
}

func (sc *serverConn) processGoAway(f *GoAwayFrame) error {
	sc.serveG.check()
	if f.ErrCode != ErrCodeNo {
		sc.logf("http2: received GOAWAY %+v, starting graceful shutdown", f)
	} else {
		sc.vlogf("http2: received GOAWAY %+v, starting graceful shutdown", f)
	}
	sc.startGracefulShutdownInternal()
	
	
	sc.pushEnabled = false
	return nil
}


func (st *stream) isPushed() bool {
	return st.id%2 == 0
}



func (st *stream) endStream() {
	sc := st.sc
	sc.serveG.check()

	if st.declBodyBytes != -1 && st.declBodyBytes != st.bodyBytes {
		st.body.CloseWithError(fmt.Errorf("request declared a Content-Length of %d but only wrote %d bytes",
			st.declBodyBytes, st.bodyBytes))
	} else {
		st.body.closeWithErrorAndCode(io.EOF, st.copyTrailersToHandlerRequest)
		st.body.CloseWithError(io.EOF)
	}
	st.state = stateHalfClosedRemote
}



func (st *stream) copyTrailersToHandlerRequest() {
	for k, vv := range st.trailer {
		if _, ok := st.reqTrailer[k]; ok {
			
			st.reqTrailer[k] = vv
		}
	}
}



func (st *stream) onWriteTimeout() {
	st.sc.writeFrameFromHandler(FrameWriteRequest{write: streamError(st.id, ErrCodeInternal)})
}

func (sc *serverConn) processHeaders(f *MetaHeadersFrame) error {
	sc.serveG.check()
	id := f.StreamID
	if sc.inGoAway {
		
		return nil
	}
	
	
	
	
	
	if id%2 != 1 {
		return ConnectionError(ErrCodeProtocol)
	}
	
	
	
	
	if st := sc.streams[f.StreamID]; st != nil {
		if st.resetQueued {
			
			
			return nil
		}
		return st.processTrailerHeaders(f)
	}

	
	
	
	
	
	if id <= sc.maxClientStreamID {
		return ConnectionError(ErrCodeProtocol)
	}
	sc.maxClientStreamID = id

	if sc.idleTimer != nil {
		sc.idleTimer.Stop()
	}

	
	
	
	
	
	
	if sc.curClientStreams+1 > sc.advMaxStreams {
		if sc.unackedSettings == 0 {
			
			return streamError(id, ErrCodeProtocol)
		}
		
		
		
		
		
		return streamError(id, ErrCodeRefusedStream)
	}

	initialState := stateOpen
	if f.StreamEnded() {
		initialState = stateHalfClosedRemote
	}
	st := sc.newStream(id, 0, initialState)

	if f.HasPriority() {
		if err := checkPriority(f.StreamID, f.Priority); err != nil {
			return err
		}
		sc.writeSched.AdjustStream(st.id, f.Priority)
	}

	rw, req, err := sc.newWriterAndRequest(st, f)
	if err != nil {
		return err
	}
	st.reqTrailer = req.Trailer
	if st.reqTrailer != nil {
		st.trailer = make(http.Header)
	}
	st.body = req.Body.(*requestBody).pipe 
	st.declBodyBytes = req.ContentLength

	handler := sc.handler.ServeHTTP
	if f.Truncated {
		
		handler = handleHeaderListTooLong
	} else if err := checkValidHTTP2RequestHeaders(req.Header); err != nil {
		handler = new400Handler(err)
	}

	
	
	
	
	
	
	
	if sc.hs.ReadTimeout != 0 {
		sc.conn.SetReadDeadline(time.Time{})
	}

	go sc.runHandler(rw, req, handler)
	return nil
}

func (st *stream) processTrailerHeaders(f *MetaHeadersFrame) error {
	sc := st.sc
	sc.serveG.check()
	if st.gotTrailerHeader {
		return ConnectionError(ErrCodeProtocol)
	}
	st.gotTrailerHeader = true
	if !f.StreamEnded() {
		return streamError(st.id, ErrCodeProtocol)
	}

	if len(f.PseudoFields()) > 0 {
		return streamError(st.id, ErrCodeProtocol)
	}
	if st.trailer != nil {
		for _, hf := range f.RegularFields() {
			key := sc.canonicalHeader(hf.Name)
			if !ValidTrailerHeader(key) {
				
				
				
				return streamError(st.id, ErrCodeProtocol)
			}
			st.trailer[key] = append(st.trailer[key], hf.Value)
		}
	}
	st.endStream()
	return nil
}

func checkPriority(streamID uint32, p PriorityParam) error {
	if streamID == p.StreamDep {
		
		
		
		
		return streamError(streamID, ErrCodeProtocol)
	}
	return nil
}

func (sc *serverConn) processPriority(f *PriorityFrame) error {
	if sc.inGoAway {
		return nil
	}
	if err := checkPriority(f.StreamID, f.PriorityParam); err != nil {
		return err
	}
	sc.writeSched.AdjustStream(f.StreamID, f.PriorityParam)
	return nil
}

func (sc *serverConn) newStream(id, pusherID uint32, state streamState) *stream {
	sc.serveG.check()
	if id == 0 {
		panic("internal error: cannot create stream with id 0")
	}

	ctx, cancelCtx := contextWithCancel(sc.baseCtx)
	st := &stream{
		sc:        sc,
		id:        id,
		state:     state,
		ctx:       ctx,
		cancelCtx: cancelCtx,
	}
	st.cw.Init()
	st.flow.conn = &sc.flow 
	st.flow.add(sc.initialStreamSendWindowSize)
	st.inflow.conn = &sc.inflow 
	st.inflow.add(sc.srv.initialStreamRecvWindowSize())
	if sc.hs.WriteTimeout != 0 {
		st.writeDeadline = time.AfterFunc(sc.hs.WriteTimeout, st.onWriteTimeout)
	}

	sc.streams[id] = st
	sc.writeSched.OpenStream(st.id, OpenStreamOptions{PusherID: pusherID})
	if st.isPushed() {
		sc.curPushedStreams++
	} else {
		sc.curClientStreams++
	}
	if sc.curOpenStreams() == 1 {
		sc.setConnState(http.StateActive)
	}

	return st
}

func (sc *serverConn) newWriterAndRequest(st *stream, f *MetaHeadersFrame) (*responseWriter, *http.Request, error) {
	sc.serveG.check()

	rp := requestParam{
		method:    f.PseudoValue("method"),
		scheme:    f.PseudoValue("scheme"),
		authority: f.PseudoValue("authority"),
		path:      f.PseudoValue("path"),
	}

	isConnect := rp.method == "CONNECT"
	if isConnect {
		if rp.path != "" || rp.scheme != "" || rp.authority == "" {
			return nil, nil, streamError(f.StreamID, ErrCodeProtocol)
		}
	} else if rp.method == "" || rp.path == "" || (rp.scheme != "https" && rp.scheme != "http") {
		
		
		
		
		
		
		
		
		
		
		return nil, nil, streamError(f.StreamID, ErrCodeProtocol)
	}

	bodyOpen := !f.StreamEnded()
	if rp.method == "HEAD" && bodyOpen {
		
		return nil, nil, streamError(f.StreamID, ErrCodeProtocol)
	}

	rp.header = make(http.Header)
	for _, hf := range f.RegularFields() {
		rp.header.Add(sc.canonicalHeader(hf.Name), hf.Value)
	}
	if rp.authority == "" {
		rp.authority = rp.header.Get("Host")
	}

	rw, req, err := sc.newWriterAndRequestNoBody(st, rp)
	if err != nil {
		return nil, nil, err
	}
	if bodyOpen {
		if vv, ok := rp.header["Content-Length"]; ok {
			req.ContentLength, _ = strconv.ParseInt(vv[0], 10, 64)
		} else {
			req.ContentLength = -1
		}
		req.Body.(*requestBody).pipe = &pipe{
			b: &dataBuffer{expected: req.ContentLength},
		}
	}
	return rw, req, nil
}

type requestParam struct {
	method                  string
	scheme, authority, path string
	header                  http.Header
}

func (sc *serverConn) newWriterAndRequestNoBody(st *stream, rp requestParam) (*responseWriter, *http.Request, error) {
	sc.serveG.check()

	var tlsState *tls.ConnectionState 
	if rp.scheme == "https" {
		tlsState = sc.tlsState
	}

	needsContinue := rp.header.Get("Expect") == "100-continue"
	if needsContinue {
		rp.header.Del("Expect")
	}
	
	if cookies := rp.header["Cookie"]; len(cookies) > 1 {
		rp.header.Set("Cookie", strings.Join(cookies, "; "))
	}

	
	var trailer http.Header
	for _, v := range rp.header["Trailer"] {
		for _, key := range strings.Split(v, ",") {
			key = http.CanonicalHeaderKey(strings.TrimSpace(key))
			switch key {
			case "Transfer-Encoding", "Trailer", "Content-Length":
				
				
			default:
				if trailer == nil {
					trailer = make(http.Header)
				}
				trailer[key] = nil
			}
		}
	}
	delete(rp.header, "Trailer")

	var url_ *url.URL
	var requestURI string
	if rp.method == "CONNECT" {
		url_ = &url.URL{Host: rp.authority}
		requestURI = rp.authority 
	} else {
		var err error
		url_, err = url.ParseRequestURI(rp.path)
		if err != nil {
			return nil, nil, streamError(st.id, ErrCodeProtocol)
		}
		requestURI = rp.path
	}

	body := &requestBody{
		conn:          sc,
		stream:        st,
		needsContinue: needsContinue,
	}
	req := &http.Request{
		Method:     rp.method,
		URL:        url_,
		RemoteAddr: sc.remoteAddrStr,
		Header:     rp.header,
		RequestURI: requestURI,
		Proto:      "HTTP/2.0",
		ProtoMajor: 2,
		ProtoMinor: 0,
		TLS:        tlsState,
		Host:       rp.authority,
		Body:       body,
		Trailer:    trailer,
	}
	req = requestWithContext(req, st.ctx)

	rws := responseWriterStatePool.Get().(*responseWriterState)
	bwSave := rws.bw
	*rws = responseWriterState{} 
	rws.conn = sc
	rws.bw = bwSave
	rws.bw.Reset(chunkWriter{rws})
	rws.stream = st
	rws.req = req
	rws.body = body

	rw := &responseWriter{rws: rws}
	return rw, req, nil
}


func (sc *serverConn) runHandler(rw *responseWriter, req *http.Request, handler func(http.ResponseWriter, *http.Request)) {
	didPanic := true
	defer func() {
		rw.rws.stream.cancelCtx()
		if didPanic {
			e := recover()
			sc.writeFrameFromHandler(FrameWriteRequest{
				write:  handlerPanicRST{rw.rws.stream.id},
				stream: rw.rws.stream,
			})
			
			if shouldLogPanic(e) {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				sc.logf("http2: panic serving %v: %v\n%s", sc.conn.RemoteAddr(), e, buf)
			}
			return
		}
		rw.handlerDone()
	}()
	handler(rw, req)
	didPanic = false
}

func handleHeaderListTooLong(w http.ResponseWriter, r *http.Request) {
	
	
	
	
	const statusRequestHeaderFieldsTooLarge = 431 
	w.WriteHeader(statusRequestHeaderFieldsTooLarge)
	io.WriteString(w, "<h1>HTTP Error 431</h1><p>Request Header Field(s) Too Large</p>")
}



func (sc *serverConn) writeHeaders(st *stream, headerData *writeResHeaders) error {
	sc.serveG.checkNotOn() 
	var errc chan error
	if headerData.h != nil {
		
		
		
		
		errc = errChanPool.Get().(chan error)
	}
	if err := sc.writeFrameFromHandler(FrameWriteRequest{
		write:  headerData,
		stream: st,
		done:   errc,
	}); err != nil {
		return err
	}
	if errc != nil {
		select {
		case err := <-errc:
			errChanPool.Put(errc)
			return err
		case <-sc.doneServing:
			return errClientDisconnected
		case <-st.cw:
			return errStreamClosed
		}
	}
	return nil
}


func (sc *serverConn) write100ContinueHeaders(st *stream) {
	sc.writeFrameFromHandler(FrameWriteRequest{
		write:  write100ContinueHeadersFrame{st.id},
		stream: st,
	})
}



type bodyReadMsg struct {
	st *stream
	n  int
}




func (sc *serverConn) noteBodyReadFromHandler(st *stream, n int, err error) {
	sc.serveG.checkNotOn() 
	if n > 0 {
		select {
		case sc.bodyReadCh <- bodyReadMsg{st, n}:
		case <-sc.doneServing:
		}
	}
}

func (sc *serverConn) noteBodyRead(st *stream, n int) {
	sc.serveG.check()
	sc.sendWindowUpdate(nil, n) 
	if st.state != stateHalfClosedRemote && st.state != stateClosed {
		
		
		sc.sendWindowUpdate(st, n)
	}
}


func (sc *serverConn) sendWindowUpdate(st *stream, n int) {
	sc.serveG.check()
	
	
	
	
	
	const maxUint31 = 1<<31 - 1
	for n >= maxUint31 {
		sc.sendWindowUpdate32(st, maxUint31)
		n -= maxUint31
	}
	sc.sendWindowUpdate32(st, int32(n))
}


func (sc *serverConn) sendWindowUpdate32(st *stream, n int32) {
	sc.serveG.check()
	if n == 0 {
		return
	}
	if n < 0 {
		panic("negative update")
	}
	var streamID uint32
	if st != nil {
		streamID = st.id
	}
	sc.writeFrame(FrameWriteRequest{
		write:  writeWindowUpdate{streamID: streamID, n: uint32(n)},
		stream: st,
	})
	var ok bool
	if st == nil {
		ok = sc.inflow.add(n)
	} else {
		ok = st.inflow.add(n)
	}
	if !ok {
		panic("internal error; sent too many window updates without decrements?")
	}
}



type requestBody struct {
	stream        *stream
	conn          *serverConn
	closed        bool  
	sawEOF        bool  
	pipe          *pipe 
	needsContinue bool  
}

func (b *requestBody) Close() error {
	if b.pipe != nil && !b.closed {
		b.pipe.BreakWithError(errClosedBody)
	}
	b.closed = true
	return nil
}

func (b *requestBody) Read(p []byte) (n int, err error) {
	if b.needsContinue {
		b.needsContinue = false
		b.conn.write100ContinueHeaders(b.stream)
	}
	if b.pipe == nil || b.sawEOF {
		return 0, io.EOF
	}
	n, err = b.pipe.Read(p)
	if err == io.EOF {
		b.sawEOF = true
	}
	if b.conn == nil && inTests {
		return
	}
	b.conn.noteBodyReadFromHandler(b.stream, n, err)
	return
}







type responseWriter struct {
	rws *responseWriterState
}


var (
	_ http.CloseNotifier = (*responseWriter)(nil)
	_ http.Flusher       = (*responseWriter)(nil)
	_ stringWriter       = (*responseWriter)(nil)
)

type responseWriterState struct {
	
	stream *stream
	req    *http.Request
	body   *requestBody 
	conn   *serverConn

	
	bw *bufio.Writer 

	
	handlerHeader http.Header 
	snapHeader    http.Header 
	trailers      []string    
	status        int         
	wroteHeader   bool        
	sentHeader    bool        
	handlerDone   bool        
	dirty         bool        

	sentContentLen int64 
	wroteBytes     int64

	closeNotifierMu sync.Mutex 
	closeNotifierCh chan bool  
}

type chunkWriter struct{ rws *responseWriterState }

func (cw chunkWriter) Write(p []byte) (n int, err error) { return cw.rws.writeChunk(p) }

func (rws *responseWriterState) hasTrailers() bool { return len(rws.trailers) != 0 }




func (rws *responseWriterState) declareTrailer(k string) {
	k = http.CanonicalHeaderKey(k)
	if !ValidTrailerHeader(k) {
		
		rws.conn.logf("ignoring invalid trailer %q", k)
		return
	}
	if !strSliceContains(rws.trailers, k) {
		rws.trailers = append(rws.trailers, k)
	}
}







func (rws *responseWriterState) writeChunk(p []byte) (n int, err error) {
	if !rws.wroteHeader {
		rws.writeHeader(200)
	}

	isHeadResp := rws.req.Method == "HEAD"
	if !rws.sentHeader {
		rws.sentHeader = true
		var ctype, clen string
		if clen = rws.snapHeader.Get("Content-Length"); clen != "" {
			rws.snapHeader.Del("Content-Length")
			clen64, err := strconv.ParseInt(clen, 10, 64)
			if err == nil && clen64 >= 0 {
				rws.sentContentLen = clen64
			} else {
				clen = ""
			}
		}
		if clen == "" && rws.handlerDone && bodyAllowedForStatus(rws.status) && (len(p) > 0 || !isHeadResp) {
			clen = strconv.Itoa(len(p))
		}
		_, hasContentType := rws.snapHeader["Content-Type"]
		if !hasContentType && bodyAllowedForStatus(rws.status) && len(p) > 0 {
			ctype = http.DetectContentType(p)
		}
		var date string
		if _, ok := rws.snapHeader["Date"]; !ok {
			
			date = time.Now().UTC().Format(http.TimeFormat)
		}

		for _, v := range rws.snapHeader["Trailer"] {
			foreachHeaderElement(v, rws.declareTrailer)
		}

		endStream := (rws.handlerDone && !rws.hasTrailers() && len(p) == 0) || isHeadResp
		err = rws.conn.writeHeaders(rws.stream, &writeResHeaders{
			streamID:      rws.stream.id,
			httpResCode:   rws.status,
			h:             rws.snapHeader,
			endStream:     endStream,
			contentType:   ctype,
			contentLength: clen,
			date:          date,
		})
		if err != nil {
			rws.dirty = true
			return 0, err
		}
		if endStream {
			return 0, nil
		}
	}
	if isHeadResp {
		return len(p), nil
	}
	if len(p) == 0 && !rws.handlerDone {
		return 0, nil
	}

	if rws.handlerDone {
		rws.promoteUndeclaredTrailers()
	}

	endStream := rws.handlerDone && !rws.hasTrailers()
	if len(p) > 0 || endStream {
		
		if err := rws.conn.writeDataFromHandler(rws.stream, p, endStream); err != nil {
			rws.dirty = true
			return 0, err
		}
	}

	if rws.handlerDone && rws.hasTrailers() {
		err = rws.conn.writeHeaders(rws.stream, &writeResHeaders{
			streamID:  rws.stream.id,
			h:         rws.handlerHeader,
			trailers:  rws.trailers,
			endStream: true,
		})
		if err != nil {
			rws.dirty = true
		}
		return len(p), err
	}
	return len(p), nil
}













const TrailerPrefix = "Trailer:"






















func (rws *responseWriterState) promoteUndeclaredTrailers() {
	for k, vv := range rws.handlerHeader {
		if !strings.HasPrefix(k, TrailerPrefix) {
			continue
		}
		trailerKey := strings.TrimPrefix(k, TrailerPrefix)
		rws.declareTrailer(trailerKey)
		rws.handlerHeader[http.CanonicalHeaderKey(trailerKey)] = vv
	}

	if len(rws.trailers) > 1 {
		sorter := sorterPool.Get().(*sorter)
		sorter.SortStrings(rws.trailers)
		sorterPool.Put(sorter)
	}
}

func (w *responseWriter) Flush() {
	rws := w.rws
	if rws == nil {
		panic("Header called after Handler finished")
	}
	if rws.bw.Buffered() > 0 {
		if err := rws.bw.Flush(); err != nil {
			
			return
		}
	} else {
		
		
		
		
		rws.writeChunk(nil)
	}
}

func (w *responseWriter) CloseNotify() <-chan bool {
	rws := w.rws
	if rws == nil {
		panic("CloseNotify called after Handler finished")
	}
	rws.closeNotifierMu.Lock()
	ch := rws.closeNotifierCh
	if ch == nil {
		ch = make(chan bool, 1)
		rws.closeNotifierCh = ch
		cw := rws.stream.cw
		go func() {
			cw.Wait() 
			ch <- true
		}()
	}
	rws.closeNotifierMu.Unlock()
	return ch
}

func (w *responseWriter) Header() http.Header {
	rws := w.rws
	if rws == nil {
		panic("Header called after Handler finished")
	}
	if rws.handlerHeader == nil {
		rws.handlerHeader = make(http.Header)
	}
	return rws.handlerHeader
}


func checkWriteHeaderCode(code int) {
	
	
	
	
	
	
	
	
	
	
	
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

func (w *responseWriter) WriteHeader(code int) {
	rws := w.rws
	if rws == nil {
		panic("WriteHeader called after Handler finished")
	}
	rws.writeHeader(code)
}

func (rws *responseWriterState) writeHeader(code int) {
	if !rws.wroteHeader {
		checkWriteHeaderCode(code)
		rws.wroteHeader = true
		rws.status = code
		if len(rws.handlerHeader) > 0 {
			rws.snapHeader = cloneHeader(rws.handlerHeader)
		}
	}
}

func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}









func (w *responseWriter) Write(p []byte) (n int, err error) {
	return w.write(len(p), p, "")
}

func (w *responseWriter) WriteString(s string) (n int, err error) {
	return w.write(len(s), nil, s)
}


func (w *responseWriter) write(lenData int, dataB []byte, dataS string) (n int, err error) {
	rws := w.rws
	if rws == nil {
		panic("Write called after Handler finished")
	}
	if !rws.wroteHeader {
		w.WriteHeader(200)
	}
	if !bodyAllowedForStatus(rws.status) {
		return 0, http.ErrBodyNotAllowed
	}
	rws.wroteBytes += int64(len(dataB)) + int64(len(dataS)) 
	if rws.sentContentLen != 0 && rws.wroteBytes > rws.sentContentLen {
		
		return 0, errors.New("http2: handler wrote more than declared Content-Length")
	}

	if dataB != nil {
		return rws.bw.Write(dataB)
	} else {
		return rws.bw.WriteString(dataS)
	}
}

func (w *responseWriter) handlerDone() {
	rws := w.rws
	dirty := rws.dirty
	rws.handlerDone = true
	w.Flush()
	w.rws = nil
	if !dirty {
		
		
		
		
		
		
		responseWriterStatePool.Put(rws)
	}
}


var (
	ErrRecursivePush    = errors.New("http2: recursive push not allowed")
	ErrPushLimitReached = errors.New("http2: push would exceed peer's SETTINGS_MAX_CONCURRENT_STREAMS")
)



type pushOptions struct {
	Method string
	Header http.Header
}

func (w *responseWriter) push(target string, opts pushOptions) error {
	st := w.rws.stream
	sc := st.sc
	sc.serveG.checkNotOn()

	
	
	if st.isPushed() {
		return ErrRecursivePush
	}

	
	if opts.Method == "" {
		opts.Method = "GET"
	}
	if opts.Header == nil {
		opts.Header = http.Header{}
	}
	wantScheme := "http"
	if w.rws.req.TLS != nil {
		wantScheme = "https"
	}

	
	u, err := url.Parse(target)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		if !strings.HasPrefix(target, "/") {
			return fmt.Errorf("target must be an absolute URL or an absolute path: %q", target)
		}
		u.Scheme = wantScheme
		u.Host = w.rws.req.Host
	} else {
		if u.Scheme != wantScheme {
			return fmt.Errorf("cannot push URL with scheme %q from request with scheme %q", u.Scheme, wantScheme)
		}
		if u.Host == "" {
			return errors.New("URL must have a host")
		}
	}
	for k := range opts.Header {
		if strings.HasPrefix(k, ":") {
			return fmt.Errorf("promised request headers cannot include pseudo header %q", k)
		}
		
		
		
		
		switch strings.ToLower(k) {
		case "content-length", "content-encoding", "trailer", "te", "expect", "host":
			return fmt.Errorf("promised request headers cannot include %q", k)
		}
	}
	if err := checkValidHTTP2RequestHeaders(opts.Header); err != nil {
		return err
	}

	
	
	
	if opts.Method != "GET" && opts.Method != "HEAD" {
		return fmt.Errorf("method %q must be GET or HEAD", opts.Method)
	}

	msg := &startPushRequest{
		parent: st,
		method: opts.Method,
		url:    u,
		header: cloneHeader(opts.Header),
		done:   errChanPool.Get().(chan error),
	}

	select {
	case <-sc.doneServing:
		return errClientDisconnected
	case <-st.cw:
		return errStreamClosed
	case sc.serveMsgCh <- msg:
	}

	select {
	case <-sc.doneServing:
		return errClientDisconnected
	case <-st.cw:
		return errStreamClosed
	case err := <-msg.done:
		errChanPool.Put(msg.done)
		return err
	}
}

type startPushRequest struct {
	parent *stream
	method string
	url    *url.URL
	header http.Header
	done   chan error
}

func (sc *serverConn) startPush(msg *startPushRequest) {
	sc.serveG.check()

	
	
	
	if msg.parent.state != stateOpen && msg.parent.state != stateHalfClosedRemote {
		
		msg.done <- errStreamClosed
		return
	}

	
	if !sc.pushEnabled {
		msg.done <- http.ErrNotSupported
		return
	}

	
	
	
	allocatePromisedID := func() (uint32, error) {
		sc.serveG.check()

		
		
		if !sc.pushEnabled {
			return 0, http.ErrNotSupported
		}
		
		if sc.curPushedStreams+1 > sc.clientMaxStreams {
			return 0, ErrPushLimitReached
		}

		
		
		
		
		if sc.maxPushPromiseID+2 >= 1<<31 {
			sc.startGracefulShutdownInternal()
			return 0, ErrPushLimitReached
		}
		sc.maxPushPromiseID += 2
		promisedID := sc.maxPushPromiseID

		
		
		
		
		
		promised := sc.newStream(promisedID, msg.parent.id, stateHalfClosedRemote)
		rw, req, err := sc.newWriterAndRequestNoBody(promised, requestParam{
			method:    msg.method,
			scheme:    msg.url.Scheme,
			authority: msg.url.Host,
			path:      msg.url.RequestURI(),
			header:    cloneHeader(msg.header), 
		})
		if err != nil {
			
			panic(fmt.Sprintf("newWriterAndRequestNoBody(%+v): %v", msg.url, err))
		}

		go sc.runHandler(rw, req, sc.handler.ServeHTTP)
		return promisedID, nil
	}

	sc.writeFrame(FrameWriteRequest{
		write: &writePushPromise{
			streamID:           msg.parent.id,
			method:             msg.method,
			url:                msg.url,
			h:                  msg.header,
			allocatePromisedID: allocatePromisedID,
		},
		stream: msg.parent,
		done:   msg.done,
	})
}



func foreachHeaderElement(v string, fn func(string)) {
	v = textproto.TrimString(v)
	if v == "" {
		return
	}
	if !strings.Contains(v, ",") {
		fn(v)
		return
	}
	for _, f := range strings.Split(v, ",") {
		if f = textproto.TrimString(f); f != "" {
			fn(f)
		}
	}
}


var connHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Connection",
	"Transfer-Encoding",
	"Upgrade",
}




func checkValidHTTP2RequestHeaders(h http.Header) error {
	for _, k := range connHeaders {
		if _, ok := h[k]; ok {
			return fmt.Errorf("request header %q is not valid in HTTP/2", k)
		}
	}
	te := h["Te"]
	if len(te) > 0 && (len(te) > 1 || (te[0] != "trailers" && te[0] != "")) {
		return errors.New(`request header "TE" may only be "trailers" in HTTP/2`)
	}
	return nil
}

func new400Handler(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}




func ValidTrailerHeader(name string) bool {
	name = http.CanonicalHeaderKey(name)
	if strings.HasPrefix(name, "If-") || badTrailer[name] {
		return false
	}
	return true
}

var badTrailer = map[string]bool{
	"Authorization":       true,
	"Cache-Control":       true,
	"Connection":          true,
	"Content-Encoding":    true,
	"Content-Length":      true,
	"Content-Range":       true,
	"Content-Type":        true,
	"Expect":              true,
	"Host":                true,
	"Keep-Alive":          true,
	"Max-Forwards":        true,
	"Pragma":              true,
	"Proxy-Authenticate":  true,
	"Proxy-Authorization": true,
	"Proxy-Connection":    true,
	"Range":               true,
	"Realm":               true,
	"Te":                  true,
	"Trailer":             true,
	"Transfer-Encoding":   true,
	"Www-Authenticate":    true,
}




func h1ServerKeepAlivesDisabled(hs *http.Server) bool {
	var x interface{} = hs
	type I interface {
		doKeepAlives() bool
	}
	if hs, ok := x.(I); ok {
		return !hs.doKeepAlives()
	}
	return false
}
