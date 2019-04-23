

package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcrand"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

var (
	
	
	ErrIllegalHeaderWrite = errors.New("transport: the stream is done or WriteHeader was already called")
	
	
	ErrHeaderListSizeLimitViolation = errors.New("transport: trying to send header list size larger than the limit set by peer")
)


type http2Server struct {
	ctx         context.Context
	ctxDone     <-chan struct{} 
	cancel      context.CancelFunc
	conn        net.Conn
	loopy       *loopyWriter
	readerDone  chan struct{} 
	writerDone  chan struct{} 
	remoteAddr  net.Addr
	localAddr   net.Addr
	maxStreamID uint32               
	authInfo    credentials.AuthInfo 
	inTapHandle tap.ServerInHandle
	framer      *framer
	
	maxStreams uint32
	
	
	controlBuf *controlBuffer
	fc         *trInFlow
	stats      stats.Handler
	
	
	activity uint32 
	
	kp keepalive.ServerParameters

	
	kep keepalive.EnforcementPolicy
	
	lastPingAt time.Time
	
	pingStrikes uint8
	
	
	
	resetPingStrikes      uint32 
	initialWindowSize     int32
	bdpEst                *bdpEstimator
	maxSendHeaderListSize *uint32

	mu sync.Mutex 

	
	
	
	
	
	
	drainChan     chan struct{}
	state         transportState
	activeStreams map[uint32]*Stream
	
	
	
	
	idle time.Time

	
	channelzID int64 
	czData     *channelzData
}



func newHTTP2Server(conn net.Conn, config *ServerConfig) (_ ServerTransport, err error) {
	writeBufSize := config.WriteBufferSize
	readBufSize := config.ReadBufferSize
	maxHeaderListSize := defaultServerMaxHeaderListSize
	if config.MaxHeaderListSize != nil {
		maxHeaderListSize = *config.MaxHeaderListSize
	}
	framer := newFramer(conn, writeBufSize, readBufSize, maxHeaderListSize)
	
	var isettings []http2.Setting
	
	
	maxStreams := config.MaxStreams
	if maxStreams == 0 {
		maxStreams = math.MaxUint32
	} else {
		isettings = append(isettings, http2.Setting{
			ID:  http2.SettingMaxConcurrentStreams,
			Val: maxStreams,
		})
	}
	dynamicWindow := true
	iwz := int32(initialWindowSize)
	if config.InitialWindowSize >= defaultWindowSize {
		iwz = config.InitialWindowSize
		dynamicWindow = false
	}
	icwz := int32(initialWindowSize)
	if config.InitialConnWindowSize >= defaultWindowSize {
		icwz = config.InitialConnWindowSize
		dynamicWindow = false
	}
	if iwz != defaultWindowSize {
		isettings = append(isettings, http2.Setting{
			ID:  http2.SettingInitialWindowSize,
			Val: uint32(iwz)})
	}
	if config.MaxHeaderListSize != nil {
		isettings = append(isettings, http2.Setting{
			ID:  http2.SettingMaxHeaderListSize,
			Val: *config.MaxHeaderListSize,
		})
	}
	if err := framer.fr.WriteSettings(isettings...); err != nil {
		return nil, connectionErrorf(false, err, "transport: %v", err)
	}
	
	if delta := uint32(icwz - defaultWindowSize); delta > 0 {
		if err := framer.fr.WriteWindowUpdate(0, delta); err != nil {
			return nil, connectionErrorf(false, err, "transport: %v", err)
		}
	}
	kp := config.KeepaliveParams
	if kp.MaxConnectionIdle == 0 {
		kp.MaxConnectionIdle = defaultMaxConnectionIdle
	}
	if kp.MaxConnectionAge == 0 {
		kp.MaxConnectionAge = defaultMaxConnectionAge
	}
	
	kp.MaxConnectionAge += getJitter(kp.MaxConnectionAge)
	if kp.MaxConnectionAgeGrace == 0 {
		kp.MaxConnectionAgeGrace = defaultMaxConnectionAgeGrace
	}
	if kp.Time == 0 {
		kp.Time = defaultServerKeepaliveTime
	}
	if kp.Timeout == 0 {
		kp.Timeout = defaultServerKeepaliveTimeout
	}
	kep := config.KeepalivePolicy
	if kep.MinTime == 0 {
		kep.MinTime = defaultKeepalivePolicyMinTime
	}
	ctx, cancel := context.WithCancel(context.Background())
	t := &http2Server{
		ctx:               ctx,
		cancel:            cancel,
		ctxDone:           ctx.Done(),
		conn:              conn,
		remoteAddr:        conn.RemoteAddr(),
		localAddr:         conn.LocalAddr(),
		authInfo:          config.AuthInfo,
		framer:            framer,
		readerDone:        make(chan struct{}),
		writerDone:        make(chan struct{}),
		maxStreams:        maxStreams,
		inTapHandle:       config.InTapHandle,
		fc:                &trInFlow{limit: uint32(icwz)},
		state:             reachable,
		activeStreams:     make(map[uint32]*Stream),
		stats:             config.StatsHandler,
		kp:                kp,
		idle:              time.Now(),
		kep:               kep,
		initialWindowSize: iwz,
		czData:            new(channelzData),
	}
	t.controlBuf = newControlBuffer(t.ctxDone)
	if dynamicWindow {
		t.bdpEst = &bdpEstimator{
			bdp:               initialWindowSize,
			updateFlowControl: t.updateFlowControl,
		}
	}
	if t.stats != nil {
		t.ctx = t.stats.TagConn(t.ctx, &stats.ConnTagInfo{
			RemoteAddr: t.remoteAddr,
			LocalAddr:  t.localAddr,
		})
		connBegin := &stats.ConnBegin{}
		t.stats.HandleConn(t.ctx, connBegin)
	}
	if channelz.IsOn() {
		t.channelzID = channelz.RegisterNormalSocket(t, config.ChannelzParentID, fmt.Sprintf("%s -> %s", t.remoteAddr, t.localAddr))
	}
	t.framer.writer.Flush()

	defer func() {
		if err != nil {
			t.Close()
		}
	}()

	
	preface := make([]byte, len(clientPreface))
	if _, err := io.ReadFull(t.conn, preface); err != nil {
		return nil, connectionErrorf(false, err, "transport: http2Server.HandleStreams failed to receive the preface from client: %v", err)
	}
	if !bytes.Equal(preface, clientPreface) {
		return nil, connectionErrorf(false, nil, "transport: http2Server.HandleStreams received bogus greeting from client: %q", preface)
	}

	frame, err := t.framer.fr.ReadFrame()
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return nil, err
	}
	if err != nil {
		return nil, connectionErrorf(false, err, "transport: http2Server.HandleStreams failed to read initial settings frame: %v", err)
	}
	atomic.StoreUint32(&t.activity, 1)
	sf, ok := frame.(*http2.SettingsFrame)
	if !ok {
		return nil, connectionErrorf(false, nil, "transport: http2Server.HandleStreams saw invalid preface type %T from client", frame)
	}
	t.handleSettings(sf)

	go func() {
		t.loopy = newLoopyWriter(serverSide, t.framer, t.controlBuf, t.bdpEst)
		t.loopy.ssGoAwayHandler = t.outgoingGoAwayHandler
		if err := t.loopy.run(); err != nil {
			errorf("transport: loopyWriter.run returning. Err: %v", err)
		}
		t.conn.Close()
		close(t.writerDone)
	}()
	go t.keepalive()
	return t, nil
}


func (t *http2Server) operateHeaders(frame *http2.MetaHeadersFrame, handle func(*Stream), traceCtx func(context.Context, string) context.Context) (fatal bool) {
	streamID := frame.Header().StreamID
	state := decodeState{serverSide: true}
	if err := state.decodeHeader(frame); err != nil {
		if se, ok := status.FromError(err); ok {
			t.controlBuf.put(&cleanupStream{
				streamID: streamID,
				rst:      true,
				rstCode:  statusCodeConvTab[se.Code()],
				onWrite:  func() {},
			})
		}
		return false
	}

	buf := newRecvBuffer()
	s := &Stream{
		id:             streamID,
		st:             t,
		buf:            buf,
		fc:             &inFlow{limit: uint32(t.initialWindowSize)},
		recvCompress:   state.encoding,
		method:         state.method,
		contentSubtype: state.contentSubtype,
	}
	if frame.StreamEnded() {
		
		s.state = streamReadDone
	}
	if state.timeoutSet {
		s.ctx, s.cancel = context.WithTimeout(t.ctx, state.timeout)
	} else {
		s.ctx, s.cancel = context.WithCancel(t.ctx)
	}
	pr := &peer.Peer{
		Addr: t.remoteAddr,
	}
	
	if t.authInfo != nil {
		pr.AuthInfo = t.authInfo
	}
	s.ctx = peer.NewContext(s.ctx, pr)
	
	if len(state.mdata) > 0 {
		s.ctx = metadata.NewIncomingContext(s.ctx, state.mdata)
	}
	if state.statsTags != nil {
		s.ctx = stats.SetIncomingTags(s.ctx, state.statsTags)
	}
	if state.statsTrace != nil {
		s.ctx = stats.SetIncomingTrace(s.ctx, state.statsTrace)
	}
	if t.inTapHandle != nil {
		var err error
		info := &tap.Info{
			FullMethodName: state.method,
		}
		s.ctx, err = t.inTapHandle(s.ctx, info)
		if err != nil {
			warningf("transport: http2Server.operateHeaders got an error from InTapHandle: %v", err)
			t.controlBuf.put(&cleanupStream{
				streamID: s.id,
				rst:      true,
				rstCode:  http2.ErrCodeRefusedStream,
				onWrite:  func() {},
			})
			return false
		}
	}
	t.mu.Lock()
	if t.state != reachable {
		t.mu.Unlock()
		return false
	}
	if uint32(len(t.activeStreams)) >= t.maxStreams {
		t.mu.Unlock()
		t.controlBuf.put(&cleanupStream{
			streamID: streamID,
			rst:      true,
			rstCode:  http2.ErrCodeRefusedStream,
			onWrite:  func() {},
		})
		return false
	}
	if streamID%2 != 1 || streamID <= t.maxStreamID {
		t.mu.Unlock()
		
		errorf("transport: http2Server.HandleStreams received an illegal stream id: %v", streamID)
		return true
	}
	t.maxStreamID = streamID
	t.activeStreams[streamID] = s
	if len(t.activeStreams) == 1 {
		t.idle = time.Time{}
	}
	t.mu.Unlock()
	if channelz.IsOn() {
		atomic.AddInt64(&t.czData.streamsStarted, 1)
		atomic.StoreInt64(&t.czData.lastStreamCreatedTime, time.Now().UnixNano())
	}
	s.requestRead = func(n int) {
		t.adjustWindow(s, uint32(n))
	}
	s.ctx = traceCtx(s.ctx, s.method)
	if t.stats != nil {
		s.ctx = t.stats.TagRPC(s.ctx, &stats.RPCTagInfo{FullMethodName: s.method})
		inHeader := &stats.InHeader{
			FullMethod:  s.method,
			RemoteAddr:  t.remoteAddr,
			LocalAddr:   t.localAddr,
			Compression: s.recvCompress,
			WireLength:  int(frame.Header().Length),
		}
		t.stats.HandleRPC(s.ctx, inHeader)
	}
	s.ctxDone = s.ctx.Done()
	s.wq = newWriteQuota(defaultWriteQuota, s.ctxDone)
	s.trReader = &transportReader{
		reader: &recvBufferReader{
			ctx:     s.ctx,
			ctxDone: s.ctxDone,
			recv:    s.buf,
		},
		windowHandler: func(n int) {
			t.updateWindow(s, uint32(n))
		},
	}
	
	t.controlBuf.put(&registerStream{
		streamID: s.id,
		wq:       s.wq,
	})
	handle(s)
	return false
}




func (t *http2Server) HandleStreams(handle func(*Stream), traceCtx func(context.Context, string) context.Context) {
	defer close(t.readerDone)
	for {
		frame, err := t.framer.fr.ReadFrame()
		atomic.StoreUint32(&t.activity, 1)
		if err != nil {
			if se, ok := err.(http2.StreamError); ok {
				warningf("transport: http2Server.HandleStreams encountered http2.StreamError: %v", se)
				t.mu.Lock()
				s := t.activeStreams[se.StreamID]
				t.mu.Unlock()
				if s != nil {
					t.closeStream(s, true, se.Code, nil, false)
				} else {
					t.controlBuf.put(&cleanupStream{
						streamID: se.StreamID,
						rst:      true,
						rstCode:  se.Code,
						onWrite:  func() {},
					})
				}
				continue
			}
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				t.Close()
				return
			}
			warningf("transport: http2Server.HandleStreams failed to read frame: %v", err)
			t.Close()
			return
		}
		switch frame := frame.(type) {
		case *http2.MetaHeadersFrame:
			if t.operateHeaders(frame, handle, traceCtx) {
				t.Close()
				break
			}
		case *http2.DataFrame:
			t.handleData(frame)
		case *http2.RSTStreamFrame:
			t.handleRSTStream(frame)
		case *http2.SettingsFrame:
			t.handleSettings(frame)
		case *http2.PingFrame:
			t.handlePing(frame)
		case *http2.WindowUpdateFrame:
			t.handleWindowUpdate(frame)
		case *http2.GoAwayFrame:
			
		default:
			errorf("transport: http2Server.HandleStreams found unhandled frame type %v.", frame)
		}
	}
}

func (t *http2Server) getStream(f http2.Frame) (*Stream, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.activeStreams == nil {
		
		return nil, false
	}
	s, ok := t.activeStreams[f.Header().StreamID]
	if !ok {
		
		return nil, false
	}
	return s, true
}




func (t *http2Server) adjustWindow(s *Stream, n uint32) {
	if w := s.fc.maybeAdjust(n); w > 0 {
		t.controlBuf.put(&outgoingWindowUpdate{streamID: s.id, increment: w})
	}

}




func (t *http2Server) updateWindow(s *Stream, n uint32) {
	if w := s.fc.onRead(n); w > 0 {
		t.controlBuf.put(&outgoingWindowUpdate{streamID: s.id,
			increment: w,
		})
	}
}




func (t *http2Server) updateFlowControl(n uint32) {
	t.mu.Lock()
	for _, s := range t.activeStreams {
		s.fc.newLimit(n)
	}
	t.initialWindowSize = int32(n)
	t.mu.Unlock()
	t.controlBuf.put(&outgoingWindowUpdate{
		streamID:  0,
		increment: t.fc.newLimit(n),
	})
	t.controlBuf.put(&outgoingSettings{
		ss: []http2.Setting{
			{
				ID:  http2.SettingInitialWindowSize,
				Val: n,
			},
		},
	})

}

func (t *http2Server) handleData(f *http2.DataFrame) {
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
			t.closeStream(s, true, http2.ErrCodeFlowControl, nil, false)
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
	if f.Header().Flags.Has(http2.FlagDataEndStream) {
		
		s.compareAndSwapState(streamActive, streamReadDone)
		s.write(recvMsg{err: io.EOF})
	}
}

func (t *http2Server) handleRSTStream(f *http2.RSTStreamFrame) {
	s, ok := t.getStream(f)
	if !ok {
		return
	}
	t.closeStream(s, false, 0, nil, false)
}

func (t *http2Server) handleSettings(f *http2.SettingsFrame) {
	if f.IsAck() {
		return
	}
	var ss []http2.Setting
	var updateFuncs []func()
	f.ForeachSetting(func(s http2.Setting) error {
		switch s.ID {
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
	t.controlBuf.executeAndPut(func(interface{}) bool {
		for _, f := range updateFuncs {
			f()
		}
		return true
	}, &incomingSettings{
		ss: ss,
	})
}

const (
	maxPingStrikes     = 2
	defaultPingTimeout = 2 * time.Hour
)

func (t *http2Server) handlePing(f *http2.PingFrame) {
	if f.IsAck() {
		if f.Data == goAwayPing.data && t.drainChan != nil {
			close(t.drainChan)
			return
		}
		
		if t.bdpEst != nil {
			t.bdpEst.calculate(f.Data)
		}
		return
	}
	pingAck := &ping{ack: true}
	copy(pingAck.data[:], f.Data[:])
	t.controlBuf.put(pingAck)

	now := time.Now()
	defer func() {
		t.lastPingAt = now
	}()
	
	
	
	if atomic.CompareAndSwapUint32(&t.resetPingStrikes, 1, 0) {
		t.pingStrikes = 0
		return
	}
	t.mu.Lock()
	ns := len(t.activeStreams)
	t.mu.Unlock()
	if ns < 1 && !t.kep.PermitWithoutStream {
		
		
		if t.lastPingAt.Add(defaultPingTimeout).After(now) {
			t.pingStrikes++
		}
	} else {
		
		if t.lastPingAt.Add(t.kep.MinTime).After(now) {
			t.pingStrikes++
		}
	}

	if t.pingStrikes > maxPingStrikes {
		
		errorf("transport: Got too many pings from the client, closing the connection.")
		t.controlBuf.put(&goAway{code: http2.ErrCodeEnhanceYourCalm, debugData: []byte("too_many_pings"), closeConn: true})
	}
}

func (t *http2Server) handleWindowUpdate(f *http2.WindowUpdateFrame) {
	t.controlBuf.put(&incomingWindowUpdate{
		streamID:  f.Header().StreamID,
		increment: f.Increment,
	})
}

func appendHeaderFieldsFromMD(headerFields []hpack.HeaderField, md metadata.MD) []hpack.HeaderField {
	for k, vv := range md {
		if isReservedHeader(k) {
			
			continue
		}
		for _, v := range vv {
			headerFields = append(headerFields, hpack.HeaderField{Name: k, Value: encodeMetadataHeader(k, v)})
		}
	}
	return headerFields
}

func (t *http2Server) checkForHeaderListSize(it interface{}) bool {
	if t.maxSendHeaderListSize == nil {
		return true
	}
	hdrFrame := it.(*headerFrame)
	var sz int64
	for _, f := range hdrFrame.hf {
		if sz += int64(f.Size()); sz > int64(*t.maxSendHeaderListSize) {
			errorf("header list size to send violates the maximum size (%d bytes) set by client", *t.maxSendHeaderListSize)
			return false
		}
	}
	return true
}


func (t *http2Server) WriteHeader(s *Stream, md metadata.MD) error {
	if s.updateHeaderSent() || s.getState() == streamDone {
		return ErrIllegalHeaderWrite
	}
	s.hdrMu.Lock()
	if md.Len() > 0 {
		if s.header.Len() > 0 {
			s.header = metadata.Join(s.header, md)
		} else {
			s.header = md
		}
	}
	if err := t.writeHeaderLocked(s); err != nil {
		s.hdrMu.Unlock()
		return err
	}
	s.hdrMu.Unlock()
	return nil
}

func (t *http2Server) writeHeaderLocked(s *Stream) error {
	
	
	headerFields := make([]hpack.HeaderField, 0, 2) 
	headerFields = append(headerFields, hpack.HeaderField{Name: ":status", Value: "200"})
	headerFields = append(headerFields, hpack.HeaderField{Name: "content-type", Value: contentType(s.contentSubtype)})
	if s.sendCompress != "" {
		headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-encoding", Value: s.sendCompress})
	}
	headerFields = appendHeaderFieldsFromMD(headerFields, s.header)
	success, err := t.controlBuf.executeAndPut(t.checkForHeaderListSize, &headerFrame{
		streamID:  s.id,
		hf:        headerFields,
		endStream: false,
		onWrite: func() {
			atomic.StoreUint32(&t.resetPingStrikes, 1)
		},
	})
	if !success {
		if err != nil {
			return err
		}
		t.closeStream(s, true, http2.ErrCodeInternal, nil, false)
		return ErrHeaderListSizeLimitViolation
	}
	if t.stats != nil {
		
		
		outHeader := &stats.OutHeader{}
		t.stats.HandleRPC(s.Context(), outHeader)
	}
	return nil
}





func (t *http2Server) WriteStatus(s *Stream, st *status.Status) error {
	if s.getState() == streamDone {
		return nil
	}
	s.hdrMu.Lock()
	
	
	headerFields := make([]hpack.HeaderField, 0, 2) 
	if !s.updateHeaderSent() {                      
		if len(s.header) > 0 { 
			if err := t.writeHeaderLocked(s); err != nil {
				s.hdrMu.Unlock()
				return err
			}
		} else { 
			headerFields = append(headerFields, hpack.HeaderField{Name: ":status", Value: "200"})
			headerFields = append(headerFields, hpack.HeaderField{Name: "content-type", Value: contentType(s.contentSubtype)})
		}
	}
	headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-status", Value: strconv.Itoa(int(st.Code()))})
	headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-message", Value: encodeGrpcMessage(st.Message())})

	if p := st.Proto(); p != nil && len(p.Details) > 0 {
		stBytes, err := proto.Marshal(p)
		if err != nil {
			
			grpclog.Errorf("transport: failed to marshal rpc status: %v, error: %v", p, err)
		} else {
			headerFields = append(headerFields, hpack.HeaderField{Name: "grpc-status-details-bin", Value: encodeBinHeader(stBytes)})
		}
	}

	
	headerFields = appendHeaderFieldsFromMD(headerFields, s.trailer)
	trailingHeader := &headerFrame{
		streamID:  s.id,
		hf:        headerFields,
		endStream: true,
		onWrite: func() {
			atomic.StoreUint32(&t.resetPingStrikes, 1)
		},
	}
	s.hdrMu.Unlock()
	success, err := t.controlBuf.execute(t.checkForHeaderListSize, trailingHeader)
	if !success {
		if err != nil {
			return err
		}
		t.closeStream(s, true, http2.ErrCodeInternal, nil, false)
		return ErrHeaderListSizeLimitViolation
	}
	t.closeStream(s, false, 0, trailingHeader, true)
	if t.stats != nil {
		t.stats.HandleRPC(s.Context(), &stats.OutTrailer{})
	}
	return nil
}



func (t *http2Server) Write(s *Stream, hdr []byte, data []byte, opts *Options) error {
	if !s.isHeaderSent() { 
		if err := t.WriteHeader(s, nil); err != nil {
			
			return status.Errorf(codes.Internal, "transport: %v", err)
		}
	} else {
		
		if s.getState() == streamDone {
			
			s.cancel()
			select {
			case <-t.ctx.Done():
				return ErrConnClosing
			default:
			}
			return ContextErr(s.ctx.Err())
		}
	}
	
	emptyLen := http2MaxFrameLen - len(hdr)
	if emptyLen > len(data) {
		emptyLen = len(data)
	}
	hdr = append(hdr, data[:emptyLen]...)
	data = data[emptyLen:]
	df := &dataFrame{
		streamID: s.id,
		h:        hdr,
		d:        data,
		onEachWrite: func() {
			atomic.StoreUint32(&t.resetPingStrikes, 1)
		},
	}
	if err := s.wq.get(int32(len(hdr) + len(data))); err != nil {
		select {
		case <-t.ctx.Done():
			return ErrConnClosing
		default:
		}
		return ContextErr(s.ctx.Err())
	}
	return t.controlBuf.put(df)
}







func (t *http2Server) keepalive() {
	p := &ping{}
	var pingSent bool
	maxIdle := time.NewTimer(t.kp.MaxConnectionIdle)
	maxAge := time.NewTimer(t.kp.MaxConnectionAge)
	keepalive := time.NewTimer(t.kp.Time)
	
	
	
	defer func() {
		if !maxIdle.Stop() {
			<-maxIdle.C
		}
		if !maxAge.Stop() {
			<-maxAge.C
		}
		if !keepalive.Stop() {
			<-keepalive.C
		}
	}()
	for {
		select {
		case <-maxIdle.C:
			t.mu.Lock()
			idle := t.idle
			if idle.IsZero() { 
				t.mu.Unlock()
				maxIdle.Reset(t.kp.MaxConnectionIdle)
				continue
			}
			val := t.kp.MaxConnectionIdle - time.Since(idle)
			t.mu.Unlock()
			if val <= 0 {
				
				
				t.drain(http2.ErrCodeNo, []byte{})
				
				maxIdle.Reset(infinity)
				return
			}
			maxIdle.Reset(val)
		case <-maxAge.C:
			t.drain(http2.ErrCodeNo, []byte{})
			maxAge.Reset(t.kp.MaxConnectionAgeGrace)
			select {
			case <-maxAge.C:
				
				t.Close()
				
				maxAge.Reset(infinity)
			case <-t.ctx.Done():
			}
			return
		case <-keepalive.C:
			if atomic.CompareAndSwapUint32(&t.activity, 1, 0) {
				pingSent = false
				keepalive.Reset(t.kp.Time)
				continue
			}
			if pingSent {
				t.Close()
				
				keepalive.Reset(infinity)
				return
			}
			pingSent = true
			if channelz.IsOn() {
				atomic.AddInt64(&t.czData.kpCount, 1)
			}
			t.controlBuf.put(p)
			keepalive.Reset(t.kp.Timeout)
		case <-t.ctx.Done():
			return
		}
	}
}




func (t *http2Server) Close() error {
	t.mu.Lock()
	if t.state == closing {
		t.mu.Unlock()
		return errors.New("transport: Close() was already called")
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
		s.cancel()
	}
	if t.stats != nil {
		connEnd := &stats.ConnEnd{}
		t.stats.HandleConn(t.ctx, connEnd)
	}
	return err
}


func (t *http2Server) deleteStream(s *Stream, eosReceived bool) {
	t.mu.Lock()
	if _, ok := t.activeStreams[s.id]; !ok {
		t.mu.Unlock()
		return
	}

	delete(t.activeStreams, s.id)
	if len(t.activeStreams) == 0 {
		t.idle = time.Now()
	}
	t.mu.Unlock()

	if channelz.IsOn() {
		if eosReceived {
			atomic.AddInt64(&t.czData.streamsSucceeded, 1)
		} else {
			atomic.AddInt64(&t.czData.streamsFailed, 1)
		}
	}
}



func (t *http2Server) closeStream(s *Stream, rst bool, rstCode http2.ErrCode, hdr *headerFrame, eosReceived bool) {
	
	oldState := s.swapState(streamDone)

	
	
	
	s.cancel()

	
	t.deleteStream(s, eosReceived)

	cleanup := &cleanupStream{
		streamID: s.id,
		rst:      rst,
		rstCode:  rstCode,
		onWrite:  func() {},
	}

	
	if hdr == nil {
		t.controlBuf.put(cleanup)
		return
	}

	
	
	
	
	
	
	
	
	
	

	
	if oldState == streamDone {
		return
	}

	hdr.cleanup = cleanup
	t.controlBuf.put(hdr)
}

func (t *http2Server) RemoteAddr() net.Addr {
	return t.remoteAddr
}

func (t *http2Server) Drain() {
	t.drain(http2.ErrCodeNo, []byte{})
}

func (t *http2Server) drain(code http2.ErrCode, debugData []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.drainChan != nil {
		return
	}
	t.drainChan = make(chan struct{})
	t.controlBuf.put(&goAway{code: code, debugData: debugData, headsUp: true})
}

var goAwayPing = &ping{data: [8]byte{1, 6, 1, 8, 0, 3, 3, 9}}



func (t *http2Server) outgoingGoAwayHandler(g *goAway) (bool, error) {
	t.mu.Lock()
	if t.state == closing { 
		t.mu.Unlock()
		
		return false, ErrConnClosing
	}
	sid := t.maxStreamID
	if !g.headsUp {
		
		t.state = draining
		if len(t.activeStreams) == 0 {
			g.closeConn = true
		}
		t.mu.Unlock()
		if err := t.framer.fr.WriteGoAway(sid, g.code, g.debugData); err != nil {
			return false, err
		}
		if g.closeConn {
			
			
			t.framer.writer.Flush()
			return false, fmt.Errorf("transport: Connection closing")
		}
		return true, nil
	}
	t.mu.Unlock()
	
	
	
	
	
	
	if err := t.framer.fr.WriteGoAway(math.MaxUint32, http2.ErrCodeNo, []byte{}); err != nil {
		return false, err
	}
	if err := t.framer.fr.WritePing(false, goAwayPing.data); err != nil {
		return false, err
	}
	go func() {
		timer := time.NewTimer(time.Minute)
		defer timer.Stop()
		select {
		case <-t.drainChan:
		case <-timer.C:
		case <-t.ctx.Done():
			return
		}
		t.controlBuf.put(&goAway{code: g.code, debugData: g.debugData})
	}()
	return false, nil
}

func (t *http2Server) ChannelzMetric() *channelz.SocketInternalMetric {
	s := channelz.SocketInternalMetric{
		StreamsStarted:                   atomic.LoadInt64(&t.czData.streamsStarted),
		StreamsSucceeded:                 atomic.LoadInt64(&t.czData.streamsSucceeded),
		StreamsFailed:                    atomic.LoadInt64(&t.czData.streamsFailed),
		MessagesSent:                     atomic.LoadInt64(&t.czData.msgSent),
		MessagesReceived:                 atomic.LoadInt64(&t.czData.msgRecv),
		KeepAlivesSent:                   atomic.LoadInt64(&t.czData.kpCount),
		LastRemoteStreamCreatedTimestamp: time.Unix(0, atomic.LoadInt64(&t.czData.lastStreamCreatedTime)),
		LastMessageSentTimestamp:         time.Unix(0, atomic.LoadInt64(&t.czData.lastMsgSentTime)),
		LastMessageReceivedTimestamp:     time.Unix(0, atomic.LoadInt64(&t.czData.lastMsgRecvTime)),
		LocalFlowControlWindow:           int64(t.fc.getSize()),
		SocketOptions:                    channelz.GetSocketOption(t.conn),
		LocalAddr:                        t.localAddr,
		RemoteAddr:                       t.remoteAddr,
		
	}
	if au, ok := t.authInfo.(credentials.ChannelzSecurityInfo); ok {
		s.Security = au.GetSecurityValue()
	}
	s.RemoteFlowControlWindow = t.getOutFlowWindow()
	return &s
}

func (t *http2Server) IncrMsgSent() {
	atomic.AddInt64(&t.czData.msgSent, 1)
	atomic.StoreInt64(&t.czData.lastMsgSentTime, time.Now().UnixNano())
}

func (t *http2Server) IncrMsgRecv() {
	atomic.AddInt64(&t.czData.msgRecv, 1)
	atomic.StoreInt64(&t.czData.lastMsgRecvTime, time.Now().UnixNano())
}

func (t *http2Server) getOutFlowWindow() int64 {
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

func getJitter(v time.Duration) time.Duration {
	if v == infinity {
		return 0
	}
	
	r := int64(v / 10)
	j := grpcrand.Int63n(2*r) - r
	return time.Duration(j)
}
