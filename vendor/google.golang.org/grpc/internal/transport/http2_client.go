

package transport

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/syscall"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)


type http2Client struct {
	ctx        context.Context
	cancel     context.CancelFunc
	ctxDone    <-chan struct{} 
	userAgent  string
	md         interface{}
	conn       net.Conn 
	loopy      *loopyWriter
	remoteAddr net.Addr
	localAddr  net.Addr
	authInfo   credentials.AuthInfo 

	readerDone chan struct{} 
	writerDone chan struct{} 
	
	
	goAway chan struct{}
	
	awakenKeepalive chan struct{}

	framer *framer
	
	
	controlBuf *controlBuffer
	fc         *trInFlow
	
	scheme string

	isSecure bool

	perRPCCreds []credentials.PerRPCCredentials

	
	
	activity         uint32 
	kp               keepalive.ClientParameters
	keepaliveEnabled bool

	statsHandler stats.Handler

	initialWindowSize int32

	
	maxSendHeaderListSize *uint32

	bdpEst *bdpEstimator
	
	
	
	onSuccess func()

	maxConcurrentStreams  uint32
	streamQuota           int64
	streamsQuotaAvailable chan struct{}
	waitingStreams        uint32
	nextID                uint32

	mu            sync.Mutex 
	state         transportState
	activeStreams map[uint32]*Stream
	
	prevGoAwayID uint32
	
	
	goAwayReason GoAwayReason

	
	channelzID int64 
	czData     *channelzData

	onGoAway func(GoAwayReason)
	onClose  func()
}

func dial(ctx context.Context, fn func(context.Context, string) (net.Conn, error), addr string) (net.Conn, error) {
	if fn != nil {
		return fn(ctx, addr)
	}
	return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
}

func isTemporary(err error) bool {
	switch err := err.(type) {
	case interface {
		Temporary() bool
	}:
		return err.Temporary()
	case interface {
		Timeout() bool
	}:
		
		
		return err.Timeout()
	}
	return true
}




func newHTTP2Client(connectCtx, ctx context.Context, addr TargetInfo, opts ConnectOptions, onSuccess func(), onGoAway func(GoAwayReason), onClose func()) (_ *http2Client, err error) {
	scheme := "http"
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	conn, err := dial(connectCtx, opts.Dialer, addr.Addr)
	if err != nil {
		if opts.FailOnNonTempDialError {
			return nil, connectionErrorf(isTemporary(err), err, "transport: error while dialing: %v", err)
		}
		return nil, connectionErrorf(true, err, "transport: Error while dialing %v", err)
	}
	
	defer func(conn net.Conn) {
		if err != nil {
			conn.Close()
		}
	}(conn)
	kp := opts.KeepaliveParams
	
	if kp.Time == 0 {
		kp.Time = defaultClientKeepaliveTime
	}
	if kp.Timeout == 0 {
		kp.Timeout = defaultClientKeepaliveTimeout
	}
	keepaliveEnabled := false
	if kp.Time != infinity {
		if err = syscall.SetTCPUserTimeout(conn, kp.Timeout); err != nil {
			return nil, connectionErrorf(false, err, "transport: failed to set TCP_USER_TIMEOUT: %v", err)
		}
		keepaliveEnabled = true
	}
	var (
		isSecure bool
		authInfo credentials.AuthInfo
	)
	transportCreds := opts.TransportCredentials
	perRPCCreds := opts.PerRPCCredentials

	if b := opts.CredsBundle; b != nil {
		if t := b.TransportCredentials(); t != nil {
			transportCreds = t
		}
		if t := b.PerRPCCredentials(); t != nil {
			perRPCCreds = append(perRPCCreds, t)
		}
	}
	if transportCreds != nil {
		scheme = "https"
		conn, authInfo, err = transportCreds.ClientHandshake(connectCtx, addr.Authority, conn)
		if err != nil {
			return nil, connectionErrorf(isTemporary(err), err, "transport: authentication handshake failed: %v", err)
		}
		isSecure = true
	}
	dynamicWindow := true
	icwz := int32(initialWindowSize)
	if opts.InitialConnWindowSize >= defaultWindowSize {
		icwz = opts.InitialConnWindowSize
		dynamicWindow = false
	}
	writeBufSize := opts.WriteBufferSize
	readBufSize := opts.ReadBufferSize
	maxHeaderListSize := defaultClientMaxHeaderListSize
	if opts.MaxHeaderListSize != nil {
		maxHeaderListSize = *opts.MaxHeaderListSize
	}
	t := &http2Client{
		ctx:                   ctx,
		ctxDone:               ctx.Done(), 
		cancel:                cancel,
		userAgent:             opts.UserAgent,
		md:                    addr.Metadata,
		conn:                  conn,
		remoteAddr:            conn.RemoteAddr(),
		localAddr:             conn.LocalAddr(),
		authInfo:              authInfo,
		readerDone:            make(chan struct{}),
		writerDone:            make(chan struct{}),
		goAway:                make(chan struct{}),
		awakenKeepalive:       make(chan struct{}, 1),
		framer:                newFramer(conn, writeBufSize, readBufSize, maxHeaderListSize),
		fc:                    &trInFlow{limit: uint32(icwz)},
		scheme:                scheme,
		activeStreams:         make(map[uint32]*Stream),
		isSecure:              isSecure,
		perRPCCreds:           perRPCCreds,
		kp:                    kp,
		statsHandler:          opts.StatsHandler,
		initialWindowSize:     initialWindowSize,
		onSuccess:             onSuccess,
		nextID:                1,
		maxConcurrentStreams:  defaultMaxStreamsClient,
		streamQuota:           defaultMaxStreamsClient,
		streamsQuotaAvailable: make(chan struct{}, 1),
		czData:                new(channelzData),
		onGoAway:              onGoAway,
		onClose:               onClose,
		keepaliveEnabled:      keepaliveEnabled,
	}
	t.controlBuf = newControlBuffer(t.ctxDone)
	if opts.InitialWindowSize >= defaultWindowSize {
		t.initialWindowSize = opts.InitialWindowSize
		dynamicWindow = false
	}
	if dynamicWindow {
		t.bdpEst = &bdpEstimator{
			bdp:               initialWindowSize,
			updateFlowControl: t.updateFlowControl,
		}
	}
	
	
	t.awakenKeepalive <- struct{}{}
	if t.statsHandler != nil {
		t.ctx = t.statsHandler.TagConn(t.ctx, &stats.ConnTagInfo{
			RemoteAddr: t.remoteAddr,
			LocalAddr:  t.localAddr,
		})
		connBegin := &stats.ConnBegin{
			Client: true,
		}
		t.statsHandler.HandleConn(t.ctx, connBegin)
	}
	if channelz.IsOn() {
		t.channelzID = channelz.RegisterNormalSocket(t, opts.ChannelzParentID, fmt.Sprintf("%s -> %s", t.localAddr, t.remoteAddr))
	}
	if t.keepaliveEnabled {
		go t.keepalive()
	}
	
	
	
	go t.reader()

	
	n, err := t.conn.Write(clientPreface)
	if err != nil {
		t.Close()
		return nil, connectionErrorf(true, err, "transport: failed to write client preface: %v", err)
	}
	if n != len(clientPreface) {
		t.Close()
		return nil, connectionErrorf(true, err, "transport: preface mismatch, wrote %d bytes; want %d", n, len(clientPreface))
	}
	var ss []http2.Setting

	if t.initialWindowSize != defaultWindowSize {
		ss = append(ss, http2.Setting{
			ID:  http2.SettingInitialWindowSize,
			Val: uint32(t.initialWindowSize),
		})
	}
	if opts.MaxHeaderListSize != nil {
		ss = append(ss, http2.Setting{
			ID:  http2.SettingMaxHeaderListSize,
			Val: *opts.MaxHeaderListSize,
		})
	}
	err = t.framer.fr.WriteSettings(ss...)
	if err != nil {
		t.Close()
		return nil, connectionErrorf(true, err, "transport: failed to write initial settings frame: %v", err)
	}
	
	if delta := uint32(icwz - defaultWindowSize); delta > 0 {
		if err := t.framer.fr.WriteWindowUpdate(0, delta); err != nil {
			t.Close()
			return nil, connectionErrorf(true, err, "transport: failed to write window update: %v", err)
		}
	}

	t.framer.writer.Flush()
	go func() {
		t.loopy = newLoopyWriter(clientSide, t.framer, t.controlBuf, t.bdpEst)
		err := t.loopy.run()
		if err != nil {
			errorf("transport: loopyWriter.run returning. Err: %v", err)
		}
		
		
		if _, ok := err.(net.Error); !ok {
			t.conn.Close()
		}
		close(t.writerDone)
	}()
	return t, nil
}

func (t *http2Client) newStream(ctx context.Context, callHdr *CallHdr) *Stream {
	
	s := &Stream{
		done:           make(chan struct{}),
		method:         callHdr.Method,
		sendCompress:   callHdr.SendCompress,
		buf:            newRecvBuffer(),
		headerChan:     make(chan struct{}),
		contentSubtype: callHdr.ContentSubtype,
	}
	s.wq = newWriteQuota(defaultWriteQuota, s.done)
	s.requestRead = func(n int) {
		t.adjustWindow(s, uint32(n))
	}
	
	
	
	s.ctx = ctx
	s.trReader = &transportReader{
		reader: &recvBufferReader{
			ctx:     s.ctx,
			ctxDone: s.ctx.Done(),
			recv:    s.buf,
		},
		windowHandler: func(n int) {
			t.updateWindow(s, uint32(n))
		},
	}
	return s
}

func (t *http2Client) getPeer() *peer.Peer {
	pr := &peer.Peer{
		Addr: t.remoteAddr,
	}
	
	if t.authInfo != nil {
		pr.AuthInfo = t.authInfo
	}
	return pr
}

func (t *http2Client) createHeaderFields(ctx context.Context, callHdr *CallHdr) ([]hpack.HeaderField, error) {
	aud := t.createAudience(callHdr)
	authData, err := t.getTrAuthData(ctx, aud)
	if err != nil {
		return nil, err
	}
	callAuthData, err := t.getCallAuthData(ctx, aud, callHdr)
	if err != nil {
		return nil, err
	}
	
	
	
	hfLen := 7 
	hfLen += len(authData) + len(callAuthData)
	headerFields := make([]hpack.HeaderField, 0, hfLen)
	headerFields = append(headerFields, hpack.HeaderField{Name: ":method", Value: "POST"})
	headerFields = append(headerFields, hpack.HeaderField{Name: ":scheme", Value: t.scheme})
	headerFields = append(headerFields, hpack.HeaderField{Name: ":path", Value: callHdr.Method})
	headerFields = append(headerFields, hpack.HeaderField{Name: ":authority", Value: callHdr.Host})
	headerFields = append(headerFields, hpack.HeaderField{Name: "content-type", Value: contentType(callHdr.ContentSubtype)})
	headerFields = append(headerFields, hpack.HeaderField{Name: "user-agent", Value: t.userAgent})
	headerFields = append(headerFields, hpack.HeaderField{Name: "te", Value: "trailers"})
	if callHdr.PreviousAttempts > 0 {
		headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-previous-rpc-attempts", Value: strconv.Itoa(callHdr.PreviousAttempts)})
	}

	if callHdr.SendCompress != "" {
		headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-encoding", Value: callHdr.SendCompress})
	}
	if dl, ok := ctx.Deadline(); ok {
		
		
		timeout := dl.Sub(time.Now())
		headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-timeout", Value: encodeTimeout(timeout)})
	}
	for k, v := range authData {
		headerFields = append(headerFields, hpack.HeaderField{Name: k, Value: encodeMetadataHeader(k, v)})
	}
	for k, v := range callAuthData {
		headerFields = append(headerFields, hpack.HeaderField{Name: k, Value: encodeMetadataHeader(k, v)})
	}
	if b := stats.OutgoingTags(ctx); b != nil {
		headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-tags-bin", Value: encodeBinHeader(b)})
	}
	if b := stats.OutgoingTrace(ctx); b != nil {
		headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-trace-bin", Value: encodeBinHeader(b)})
	}

	if md, added, ok := metadata.FromOutgoingContextRaw(ctx); ok {
		var k string
		for _, vv := range added {
			for i, v := range vv {
				if i%2 == 0 {
					k = v
					continue
				}
				
				if isReservedHeader(k) {
					continue
				}
				headerFields = append(headerFields, hpack.HeaderField{Name: strings.ToLower(k), Value: encodeMetadataHeader(k, v)})
			}
		}
		for k, vv := range md {
			
			if isReservedHeader(k) {
				continue
			}
			for _, v := range vv {
				headerFields = append(headerFields, hpack.HeaderField{Name: k, Value: encodeMetadataHeader(k, v)})
			}
		}
	}
	if md, ok := t.md.(*metadata.MD); ok {
		for k, vv := range *md {
			if isReservedHeader(k) {
				continue
			}
			for _, v := range vv {
				headerFields = append(headerFields, hpack.HeaderField{Name: k, Value: encodeMetadataHeader(k, v)})
			}
		}
	}
	return headerFields, nil
}

func (t *http2Client) createAudience(callHdr *CallHdr) string {
	
	if len(t.perRPCCreds) == 0 && callHdr.Creds == nil {
		return ""
	}
	
	
	host := strings.TrimSuffix(callHdr.Host, ":443")
	pos := strings.LastIndex(callHdr.Method, "/")
	if pos == -1 {
		pos = len(callHdr.Method)
	}
	return "https://" + host + callHdr.Method[:pos]
}

func (t *http2Client) getTrAuthData(ctx context.Context, audience string) (map[string]string, error) {
	authData := map[string]string{}
	for _, c := range t.perRPCCreds {
		data, err := c.GetRequestMetadata(ctx, audience)
		if err != nil {
			if _, ok := status.FromError(err); ok {
				return nil, err
			}

			return nil, status.Errorf(codes.Unauthenticated, "transport: %v", err)
		}
		for k, v := range data {
			
			k = strings.ToLower(k)
			authData[k] = v
		}
	}
	return authData, nil
}

func (t *http2Client) getCallAuthData(ctx context.Context, audience string, callHdr *CallHdr) (map[string]string, error) {
	callAuthData := map[string]string{}
	
	
	
	if callCreds := callHdr.Creds; callCreds != nil {
		if !t.isSecure && callCreds.RequireTransportSecurity() {
			return nil, status.Error(codes.Unauthenticated, "transport: cannot send secure credentials on an insecure connection")
		}
		data, err := callCreds.GetRequestMetadata(ctx, audience)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transport: %v", err)
		}
		for k, v := range data {
			
			k = strings.ToLower(k)
			callAuthData[k] = v
		}
	}
	return callAuthData, nil
}



func (t *http2Client) NewStream(ctx context.Context, callHdr *CallHdr) (_ *Stream, err error) {
	ctx = peer.NewContext(ctx, t.getPeer())
	headerFields, err := t.createHeaderFields(ctx, callHdr)
	if err != nil {
		return nil, err
	}
	s := t.newStream(ctx, callHdr)
	cleanup := func(err error) {
		if s.swapState(streamDone) == streamDone {
			
			return
		}
		
		atomic.StoreUint32(&s.unprocessed, 1)
		s.write(recvMsg{err: err})
		close(s.done)
		
		if atomic.SwapUint32(&s.headerDone, 1) == 0 {
			close(s.headerChan)
		}

	}
	hdr := &headerFrame{
		hf:        headerFields,
		endStream: false,
		initStream: func(id uint32) (bool, error) {
			t.mu.Lock()
			if state := t.state; state != reachable {
				t.mu.Unlock()
				
				err := error(errStreamDrain)
				if state == closing {
					err = ErrConnClosing
				}
				cleanup(err)
				return false, err
			}
			t.activeStreams[id] = s
			if channelz.IsOn() {
				atomic.AddInt64(&t.czData.streamsStarted, 1)
				atomic.StoreInt64(&t.czData.lastStreamCreatedTime, time.Now().UnixNano())
			}
			var sendPing bool
			
			
			if len(t.activeStreams) == 1 && t.keepaliveEnabled {
				select {
				case t.awakenKeepalive <- struct{}{}:
					sendPing = true
					
					
					
					t.awakenKeepalive <- struct{}{}
				default:
				}
			}
			t.mu.Unlock()
			return sendPing, nil
		},
		onOrphaned: cleanup,
		wq:         s.wq,
	}
	firstTry := true
	var ch chan struct{}
	checkForStreamQuota := func(it interface{}) bool {
		if t.streamQuota <= 0 { 
			if firstTry {
				t.waitingStreams++
			}
			ch = t.streamsQuotaAvailable
			return false
		}
		if !firstTry {
			t.waitingStreams--
		}
		t.streamQuota--
		h := it.(*headerFrame)
		h.streamID = t.nextID
		t.nextID += 2
		s.id = h.streamID
		s.fc = &inFlow{limit: uint32(t.initialWindowSize)}
		if t.streamQuota > 0 && t.waitingStreams > 0 {
			select {
			case t.streamsQuotaAvailable <- struct{}{}:
			default:
			}
		}
		return true
	}
	var hdrListSizeErr error
	checkForHeaderListSize := func(it interface{}) bool {
		if t.maxSendHeaderListSize == nil {
			return true
		}
		hdrFrame := it.(*headerFrame)
		var sz int64
		for _, f := range hdrFrame.hf {
			if sz += int64(f.Size()); sz > int64(*t.maxSendHeaderListSize) {
				hdrListSizeErr = status.Errorf(codes.Internal, "header list size to send violates the maximum size (%d bytes) set by server", *t.maxSendHeaderListSize)
				return false
			}
		}
		return true
	}
	for {
		success, err := t.controlBuf.executeAndPut(func(it interface{}) bool {
			if !checkForStreamQuota(it) {
				return false
			}
			if !checkForHeaderListSize(it) {
				return false
			}
			return true
		}, hdr)
		if err != nil {
			return nil, err
		}
		if success {
			break
		}
		if hdrListSizeErr != nil {
			return nil, hdrListSizeErr
		}
		firstTry = false
		select {
		case <-ch:
		case <-s.ctx.Done():
			return nil, ContextErr(s.ctx.Err())
		case <-t.goAway:
			return nil, errStreamDrain
		case <-t.ctx.Done():
			return nil, ErrConnClosing
		}
	}
	if t.statsHandler != nil {
		outHeader := &stats.OutHeader{
			Client:      true,
			FullMethod:  callHdr.Method,
			RemoteAddr:  t.remoteAddr,
			LocalAddr:   t.localAddr,
			Compression: callHdr.SendCompress,
		}
		t.statsHandler.HandleRPC(s.ctx, outHeader)
	}
	return s, nil
}



func (t *http2Client) CloseStream(s *Stream, err error) {
	var (
		rst     bool
		rstCode http2.ErrCode
	)
	if err != nil {
		rst = true
		rstCode = http2.ErrCodeCancel
	}
	t.closeStream(s, err, rst, rstCode, status.Convert(err), nil, false)
}

func (t *http2Client) closeStream(s *Stream, err error, rst bool, rstCode http2.ErrCode, st *status.Status, mdata map[string][]string, eosReceived bool) {
	
	if s.swapState(streamDone) == streamDone {
		
		
		<-s.done
		return
	}
	
	
	
	s.status = st
	if len(mdata) > 0 {
		s.trailer = mdata
	}
	if err != nil {
		
		s.write(recvMsg{err: err})
	}
	
	if atomic.SwapUint32(&s.headerDone, 1) == 0 {
		s.noHeaders = true
		close(s.headerChan)
	}
	cleanup := &cleanupStream{
		streamID: s.id,
		onWrite: func() {
			t.mu.Lock()
			if t.activeStreams != nil {
				delete(t.activeStreams, s.id)
			}
			t.mu.Unlock()
			if channelz.IsOn() {
				if eosReceived {
					atomic.AddInt64(&t.czData.streamsSucceeded, 1)
				} else {
					atomic.AddInt64(&t.czData.streamsFailed, 1)
				}
			}
		},
		rst:     rst,
		rstCode: rstCode,
	}
	addBackStreamQuota := func(interface{}) bool {
		t.streamQuota++
		if t.streamQuota > 0 && t.waitingStreams > 0 {
			select {
			case t.streamsQuotaAvailable <- struct{}{}:
			default:
			}
		}
		return true
	}
	t.controlBuf.executeAndPut(addBackStreamQuota, cleanup)
	
	close(s.done)
}








func (t *http2Client) Close() error {
	t.mu.Lock()
	
	if t.state == closing {
		t.mu.Unlock()
		return nil
	}
	t.state = closing
	streams := t.activeStreams
	t.activeStreams = nil
	t.mu.Unlock()
	t.controlBuf.finish()
	t.cancel()
	err := t.conn.Close()
	if channelz.IsOn() {
		channelz.RemoveEntry(t.channelzID)
	}
	
	for _, s := range streams {
		t.closeStream(s, ErrConnClosing, false, http2.ErrCodeNo, status.New(codes.Unavailable, ErrConnClosing.Desc), nil, false)
	}
	if t.statsHandler != nil {
		connEnd := &stats.ConnEnd{
			Client: true,
		}
		t.statsHandler.HandleConn(t.ctx, connEnd)
	}
	go t.onClose()
	return err
}






func (t *http2Client) GracefulClose() error {
	t.mu.Lock()
	
	if t.state == draining || t.state == closing {
		t.mu.Unlock()
		return nil
	}
	t.state = draining
	active := len(t.activeStreams)
	t.mu.Unlock()
	if active == 0 {
		return t.Close()
	}
	t.controlBuf.put(&incomingGoAway{})
	return nil
}



func (t *http2Client) Write(s *Stream, hdr []byte, data []byte, opts *Options) error {
	if opts.Last {
		
		if !s.compareAndSwapState(streamActive, streamWriteDone) {
			return errStreamDone
		}
	} else if s.getState() != streamActive {
		return errStreamDone
	}
	df := &dataFrame{
		streamID:  s.id,
		endStream: opts.Last,
	}
	if hdr != nil || data != nil { 
		
		
		emptyLen := http2MaxFrameLen - len(hdr)
		if emptyLen > len(data) {
			emptyLen = len(data)
		}
		hdr = append(hdr, data[:emptyLen]...)
		data = data[emptyLen:]
		df.h, df.d = hdr, data
		
		if err := s.wq.get(int32(len(hdr) + len(data))); err != nil {
			return err
		}
	}
	return t.controlBuf.put(df)
}

func (t *http2Client) getStream(f http2.Frame) (*Stream, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s, ok := t.activeStreams[f.Header().StreamID]
	return s, ok
}




func (t *http2Client) adjustWindow(s *Stream, n uint32) {
	if w := s.fc.maybeAdjust(n); w > 0 {
		t.controlBuf.put(&outgoingWindowUpdate{streamID: s.id, increment: w})
	}
}




func (t *http2Client) updateWindow(s *Stream, n uint32) {
	if w := s.fc.onRead(n); w > 0 {
		t.controlBuf.put(&outgoingWindowUpdate{streamID: s.id, increment: w})
	}
}




func (t *http2Client) updateFlowControl(n uint32) {
	t.mu.Lock()
	for _, s := range t.activeStreams {
		s.fc.newLimit(n)
	}
	t.mu.Unlock()
	updateIWS := func(interface{}) bool {
		t.initialWindowSize = int32(n)
		return true
	}
	t.controlBuf.executeAndPut(updateIWS, &outgoingWindowUpdate{streamID: 0, increment: t.fc.newLimit(n)})
	t.controlBuf.put(&outgoingSettings{
		ss: []http2.Setting{
			{
				ID:  http2.SettingInitialWindowSize,
				Val: n,
			},
		},
	})
}

func (t *http2Client) handleData(f *http2.DataFrame) {
	size := f.Header().Length
	var sendBDPPing bool
	if t.bdpEst != nil {
		sendBDPPing = t.bdpEst.add(size)
	}
	
	
	
	
	
	
	
	
	
	if w := t.fc.onData(size); w > 0 {
		t.controlBuf.put(&outgoingWindowUpdate{
			streamID:  0,
			increment: w,
		})
	}
	if sendBDPPing {
		
		

		if w := t.fc.reset(); w > 0 {
			t.controlBuf.put(&outgoingWindowUpdate{
				streamID:  0,
				increment: w,
			})
		}

		t.controlBuf.put(bdpPing)
	}
	
	s, ok := t.getStream(f)
	if !ok {
		return
	}
	if size > 0 {
		if err := s.fc.onData(size); err != nil {
			t.closeStream(s, io.EOF, true, http2.ErrCodeFlowControl, status.New(codes.Internal, err.Error()), nil, false)
			return
		}
		if f.Header().Flags.Has(http2.FlagDataPadded) {
			if w := s.fc.onRead(size - uint32(len(f.Data()))); w > 0 {
				t.controlBuf.put(&outgoingWindowUpdate{s.id, w})
			}
		}
		
		
		
		if len(f.Data()) > 0 {
			data := make([]byte, len(f.Data()))
			copy(data, f.Data())
			s.write(recvMsg{data: data})
		}
	}
	
	
	if f.FrameHeader.Flags.Has(http2.FlagDataEndStream) {
		t.closeStream(s, io.EOF, false, http2.ErrCodeNo, status.New(codes.Internal, "server closed the stream without sending trailers"), nil, true)
	}
}

func (t *http2Client) handleRSTStream(f *http2.RSTStreamFrame) {
	s, ok := t.getStream(f)
	if !ok {
		return
	}
	if f.ErrCode == http2.ErrCodeRefusedStream {
		
		atomic.StoreUint32(&s.unprocessed, 1)
	}
	statusCode, ok := http2ErrConvTab[f.ErrCode]
	if !ok {
		warningf("transport: http2Client.handleRSTStream found no mapped gRPC status for the received http2 error %v", f.ErrCode)
		statusCode = codes.Unknown
	}
	if statusCode == codes.Canceled {
		
		
		if d, ok := s.ctx.Deadline(); ok && d.After(time.Now()) {
			statusCode = codes.DeadlineExceeded
		}
	}
	t.closeStream(s, io.EOF, false, http2.ErrCodeNo, status.Newf(statusCode, "stream terminated by RST_STREAM with error code: %v", f.ErrCode), nil, false)
}

func (t *http2Client) handleSettings(f *http2.SettingsFrame, isFirst bool) {
	if f.IsAck() {
		return
	}
	var maxStreams *uint32
	var ss []http2.Setting
	var updateFuncs []func()
	f.ForeachSetting(func(s http2.Setting) error {
		switch s.ID {
		case http2.SettingMaxConcurrentStreams:
			maxStreams = new(uint32)
			*maxStreams = s.Val
		case http2.SettingMaxHeaderListSize:
			updateFuncs = append(updateFuncs, func() {
				t.maxSendHeaderListSize = new(uint32)
				*t.maxSendHeaderListSize = s.Val
			})
		default:
			ss = append(ss, s)
		}
		return nil
	})
	if isFirst && maxStreams == nil {
		maxStreams = new(uint32)
		*maxStreams = math.MaxUint32
	}
	sf := &incomingSettings{
		ss: ss,
	}
	if maxStreams != nil {
		updateStreamQuota := func() {
			delta := int64(*maxStreams) - int64(t.maxConcurrentStreams)
			t.maxConcurrentStreams = *maxStreams
			t.streamQuota += delta
			if delta > 0 && t.waitingStreams > 0 {
				close(t.streamsQuotaAvailable) 
				t.streamsQuotaAvailable = make(chan struct{}, 1)
			}
		}
		updateFuncs = append(updateFuncs, updateStreamQuota)
	}
	t.controlBuf.executeAndPut(func(interface{}) bool {
		for _, f := range updateFuncs {
			f()
		}
		return true
	}, sf)
}

func (t *http2Client) handlePing(f *http2.PingFrame) {
	if f.IsAck() {
		
		if t.bdpEst != nil {
			t.bdpEst.calculate(f.Data)
		}
		return
	}
	pingAck := &ping{ack: true}
	copy(pingAck.data[:], f.Data[:])
	t.controlBuf.put(pingAck)
}

func (t *http2Client) handleGoAway(f *http2.GoAwayFrame) {
	t.mu.Lock()
	if t.state == closing {
		t.mu.Unlock()
		return
	}
	if f.ErrCode == http2.ErrCodeEnhanceYourCalm {
		infof("Client received GoAway with http2.ErrCodeEnhanceYourCalm.")
	}
	id := f.LastStreamID
	if id > 0 && id%2 != 1 {
		t.mu.Unlock()
		t.Close()
		return
	}
	
	
	
	
	
	
	
	
	
	
	select {
	case <-t.goAway: 
		
		if id > t.prevGoAwayID {
			t.mu.Unlock()
			t.Close()
			return
		}
	default:
		t.setGoAwayReason(f)
		close(t.goAway)
		t.state = draining
		t.controlBuf.put(&incomingGoAway{})

		
		t.onGoAway(t.goAwayReason)
	}
	
	
	upperLimit := t.prevGoAwayID
	if upperLimit == 0 { 
		upperLimit = math.MaxUint32 
	}
	for streamID, stream := range t.activeStreams {
		if streamID > id && streamID <= upperLimit {
			
			atomic.StoreUint32(&stream.unprocessed, 1)
			t.closeStream(stream, errStreamDrain, false, http2.ErrCodeNo, statusGoAway, nil, false)
		}
	}
	t.prevGoAwayID = id
	active := len(t.activeStreams)
	t.mu.Unlock()
	if active == 0 {
		t.Close()
	}
}





func (t *http2Client) setGoAwayReason(f *http2.GoAwayFrame) {
	t.goAwayReason = GoAwayNoReason
	switch f.ErrCode {
	case http2.ErrCodeEnhanceYourCalm:
		if string(f.DebugData()) == "too_many_pings" {
			t.goAwayReason = GoAwayTooManyPings
		}
	}
}

func (t *http2Client) GetGoAwayReason() GoAwayReason {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.goAwayReason
}

func (t *http2Client) handleWindowUpdate(f *http2.WindowUpdateFrame) {
	t.controlBuf.put(&incomingWindowUpdate{
		streamID:  f.Header().StreamID,
		increment: f.Increment,
	})
}


func (t *http2Client) operateHeaders(frame *http2.MetaHeadersFrame) {
	s, ok := t.getStream(frame)
	if !ok {
		return
	}
	atomic.StoreUint32(&s.bytesReceived, 1)
	var state decodeState
	if err := state.decodeHeader(frame); err != nil {
		t.closeStream(s, err, true, http2.ErrCodeProtocol, status.New(codes.Internal, err.Error()), nil, false)
		
		return
	}

	endStream := frame.StreamEnded()
	var isHeader bool
	defer func() {
		if t.statsHandler != nil {
			if isHeader {
				inHeader := &stats.InHeader{
					Client:     true,
					WireLength: int(frame.Header().Length),
				}
				t.statsHandler.HandleRPC(s.ctx, inHeader)
			} else {
				inTrailer := &stats.InTrailer{
					Client:     true,
					WireLength: int(frame.Header().Length),
				}
				t.statsHandler.HandleRPC(s.ctx, inTrailer)
			}
		}
	}()
	
	if atomic.SwapUint32(&s.headerDone, 1) == 0 {
		if !endStream {
			
			isHeader = true
			
			
			
			s.recvCompress = state.encoding
			if len(state.mdata) > 0 {
				s.header = state.mdata
			}
		} else {
			s.noHeaders = true
		}
		close(s.headerChan)
	}
	if !endStream {
		return
	}
	
	rst := s.getState() == streamActive
	t.closeStream(s, io.EOF, rst, http2.ErrCodeNo, state.status(), state.mdata, true)
}







func (t *http2Client) reader() {
	defer close(t.readerDone)
	
	frame, err := t.framer.fr.ReadFrame()
	if err != nil {
		t.Close() 
		return
	}
	t.conn.SetReadDeadline(time.Time{}) 
	if t.keepaliveEnabled {
		atomic.CompareAndSwapUint32(&t.activity, 0, 1)
	}
	sf, ok := frame.(*http2.SettingsFrame)
	if !ok {
		t.Close() 
		return
	}
	t.onSuccess()
	t.handleSettings(sf, true)

	
	for {
		frame, err := t.framer.fr.ReadFrame()
		if t.keepaliveEnabled {
			atomic.CompareAndSwapUint32(&t.activity, 0, 1)
		}
		if err != nil {
			
			
			
			if se, ok := err.(http2.StreamError); ok {
				t.mu.Lock()
				s := t.activeStreams[se.StreamID]
				t.mu.Unlock()
				if s != nil {
					
					code := http2ErrConvTab[se.Code]
					msg := t.framer.fr.ErrorDetail().Error()
					t.closeStream(s, status.Error(code, msg), true, http2.ErrCodeProtocol, status.New(code, msg), nil, false)
				}
				continue
			} else {
				
				t.Close()
				return
			}
		}
		switch frame := frame.(type) {
		case *http2.MetaHeadersFrame:
			t.operateHeaders(frame)
		case *http2.DataFrame:
			t.handleData(frame)
		case *http2.RSTStreamFrame:
			t.handleRSTStream(frame)
		case *http2.SettingsFrame:
			t.handleSettings(frame, false)
		case *http2.PingFrame:
			t.handlePing(frame)
		case *http2.GoAwayFrame:
			t.handleGoAway(frame)
		case *http2.WindowUpdateFrame:
			t.handleWindowUpdate(frame)
		default:
			errorf("transport: http2Client.reader got unhandled frame type %v.", frame)
		}
	}
}


func (t *http2Client) keepalive() {
	p := &ping{data: [8]byte{}}
	timer := time.NewTimer(t.kp.Time)
	for {
		select {
		case <-timer.C:
			if atomic.CompareAndSwapUint32(&t.activity, 1, 0) {
				timer.Reset(t.kp.Time)
				continue
			}
			
			t.mu.Lock()
			if len(t.activeStreams) < 1 && !t.kp.PermitWithoutStream {
				
				<-t.awakenKeepalive
				t.mu.Unlock()
				select {
				case <-t.awakenKeepalive:
					
					
				case <-t.ctx.Done():
					return
				}
			} else {
				t.mu.Unlock()
				if channelz.IsOn() {
					atomic.AddInt64(&t.czData.kpCount, 1)
				}
				
				t.controlBuf.put(p)
			}

			
			timer.Reset(t.kp.Timeout)
			select {
			case <-timer.C:
				if atomic.CompareAndSwapUint32(&t.activity, 1, 0) {
					timer.Reset(t.kp.Time)
					continue
				}
				t.Close()
				return
			case <-t.ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return
			}
		case <-t.ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		}
	}
}

func (t *http2Client) Error() <-chan struct{} {
	return t.ctx.Done()
}

func (t *http2Client) GoAway() <-chan struct{} {
	return t.goAway
}

func (t *http2Client) ChannelzMetric() *channelz.SocketInternalMetric {
	s := channelz.SocketInternalMetric{
		StreamsStarted:                  atomic.LoadInt64(&t.czData.streamsStarted),
		StreamsSucceeded:                atomic.LoadInt64(&t.czData.streamsSucceeded),
		StreamsFailed:                   atomic.LoadInt64(&t.czData.streamsFailed),
		MessagesSent:                    atomic.LoadInt64(&t.czData.msgSent),
		MessagesReceived:                atomic.LoadInt64(&t.czData.msgRecv),
		KeepAlivesSent:                  atomic.LoadInt64(&t.czData.kpCount),
		LastLocalStreamCreatedTimestamp: time.Unix(0, atomic.LoadInt64(&t.czData.lastStreamCreatedTime)),
		LastMessageSentTimestamp:        time.Unix(0, atomic.LoadInt64(&t.czData.lastMsgSentTime)),
		LastMessageReceivedTimestamp:    time.Unix(0, atomic.LoadInt64(&t.czData.lastMsgRecvTime)),
		LocalFlowControlWindow:          int64(t.fc.getSize()),
		SocketOptions:                   channelz.GetSocketOption(t.conn),
		LocalAddr:                       t.localAddr,
		RemoteAddr:                      t.remoteAddr,
		
	}
	if au, ok := t.authInfo.(credentials.ChannelzSecurityInfo); ok {
		s.Security = au.GetSecurityValue()
	}
	s.RemoteFlowControlWindow = t.getOutFlowWindow()
	return &s
}

func (t *http2Client) IncrMsgSent() {
	atomic.AddInt64(&t.czData.msgSent, 1)
	atomic.StoreInt64(&t.czData.lastMsgSentTime, time.Now().UnixNano())
}

func (t *http2Client) IncrMsgRecv() {
	atomic.AddInt64(&t.czData.msgRecv, 1)
	atomic.StoreInt64(&t.czData.lastMsgRecvTime, time.Now().UnixNano())
}

func (t *http2Client) getOutFlowWindow() int64 {
	resp := make(chan uint32, 1)
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	t.controlBuf.put(&outFlowControlSizeRequest{resp})
	select {
	case sz := <-resp:
		return int64(sz)
	case <-t.ctxDone:
		return -1
	case <-timer.C:
		return -2
	}
}
