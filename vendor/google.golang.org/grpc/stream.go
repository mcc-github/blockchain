

package grpc

import (
	"context"
	"errors"
	"io"
	"math"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/trace"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/balancerload"
	"google.golang.org/grpc/internal/binarylog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcrand"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)






type StreamHandler func(srv interface{}, stream ServerStream) error


type StreamDesc struct {
	StreamName string
	Handler    StreamHandler

	
	ServerStreams bool
	ClientStreams bool
}




type Stream interface {
	
	Context() context.Context
	
	SendMsg(m interface{}) error
	
	RecvMsg(m interface{}) error
}





type ClientStream interface {
	
	
	Header() (metadata.MD, error)
	
	
	
	Trailer() metadata.MD
	
	
	
	CloseSend() error
	
	
	
	
	Context() context.Context
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	SendMsg(m interface{}) error
	
	
	
	
	
	
	
	
	RecvMsg(m interface{}) error
}

















func (cc *ClientConn) NewStream(ctx context.Context, desc *StreamDesc, method string, opts ...CallOption) (ClientStream, error) {
	
	
	opts = combine(cc.dopts.callOptions, opts)

	if cc.dopts.streamInt != nil {
		return cc.dopts.streamInt(ctx, desc, cc, method, newClientStream, opts...)
	}
	return newClientStream(ctx, desc, cc, method, opts...)
}


func NewClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (ClientStream, error) {
	return cc.NewStream(ctx, desc, method, opts...)
}

func newClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (_ ClientStream, err error) {
	if channelz.IsOn() {
		cc.incrCallsStarted()
		defer func() {
			if err != nil {
				cc.incrCallsFailed()
			}
		}()
	}
	c := defaultCallInfo()
	
	
	if err := cc.waitForResolvedAddrs(ctx); err != nil {
		return nil, err
	}
	mc := cc.GetMethodConfig(method)
	if mc.WaitForReady != nil {
		c.failFast = !*mc.WaitForReady
	}

	
	
	
	
	
	var cancel context.CancelFunc
	if mc.Timeout != nil && *mc.Timeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, *mc.Timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	for _, o := range opts {
		if err := o.before(c); err != nil {
			return nil, toRPCErr(err)
		}
	}
	c.maxSendMessageSize = getMaxSize(mc.MaxReqSize, c.maxSendMessageSize, defaultClientMaxSendMessageSize)
	c.maxReceiveMessageSize = getMaxSize(mc.MaxRespSize, c.maxReceiveMessageSize, defaultClientMaxReceiveMessageSize)
	if err := setCallInfoCodec(c); err != nil {
		return nil, err
	}

	callHdr := &transport.CallHdr{
		Host:           cc.authority,
		Method:         method,
		ContentSubtype: c.contentSubtype,
	}

	
	
	
	
	var cp Compressor
	var comp encoding.Compressor
	if ct := c.compressorType; ct != "" {
		callHdr.SendCompress = ct
		if ct != encoding.Identity {
			comp = encoding.GetCompressor(ct)
			if comp == nil {
				return nil, status.Errorf(codes.Internal, "grpc: Compressor is not installed for requested grpc-encoding %q", ct)
			}
		}
	} else if cc.dopts.cp != nil {
		callHdr.SendCompress = cc.dopts.cp.Type()
		cp = cc.dopts.cp
	}
	if c.creds != nil {
		callHdr.Creds = c.creds
	}
	var trInfo *traceInfo
	if EnableTracing {
		trInfo = &traceInfo{
			tr: trace.New("grpc.Sent."+methodFamily(method), method),
			firstLine: firstLine{
				client: true,
			},
		}
		if deadline, ok := ctx.Deadline(); ok {
			trInfo.firstLine.deadline = time.Until(deadline)
		}
		trInfo.tr.LazyLog(&trInfo.firstLine, false)
		ctx = trace.NewContext(ctx, trInfo.tr)
	}
	ctx = newContextWithRPCInfo(ctx, c.failFast)
	sh := cc.dopts.copts.StatsHandler
	var beginTime time.Time
	if sh != nil {
		ctx = sh.TagRPC(ctx, &stats.RPCTagInfo{FullMethodName: method, FailFast: c.failFast})
		beginTime = time.Now()
		begin := &stats.Begin{
			Client:    true,
			BeginTime: beginTime,
			FailFast:  c.failFast,
		}
		sh.HandleRPC(ctx, begin)
	}

	cs := &clientStream{
		callHdr:      callHdr,
		ctx:          ctx,
		methodConfig: &mc,
		opts:         opts,
		callInfo:     c,
		cc:           cc,
		desc:         desc,
		codec:        c.codec,
		cp:           cp,
		comp:         comp,
		cancel:       cancel,
		beginTime:    beginTime,
		firstAttempt: true,
	}
	if !cc.dopts.disableRetry {
		cs.retryThrottler = cc.retryThrottler.Load().(*retryThrottler)
	}
	cs.binlog = binarylog.GetMethodLogger(method)

	cs.callInfo.stream = cs
	
	
	if err := cs.newAttemptLocked(sh, trInfo); err != nil {
		cs.finish(err)
		return nil, err
	}

	op := func(a *csAttempt) error { return a.newStream() }
	if err := cs.withRetry(op, func() { cs.bufferForRetryLocked(0, op) }); err != nil {
		cs.finish(err)
		return nil, err
	}

	if cs.binlog != nil {
		md, _ := metadata.FromOutgoingContext(ctx)
		logEntry := &binarylog.ClientHeader{
			OnClientSide: true,
			Header:       md,
			MethodName:   method,
			Authority:    cs.cc.authority,
		}
		if deadline, ok := ctx.Deadline(); ok {
			logEntry.Timeout = time.Until(deadline)
			if logEntry.Timeout < 0 {
				logEntry.Timeout = 0
			}
		}
		cs.binlog.Log(logEntry)
	}

	if desc != unaryStreamDesc {
		
		
		
		
		
		go func() {
			select {
			case <-cc.ctx.Done():
				cs.finish(ErrClientConnClosing)
			case <-ctx.Done():
				cs.finish(toRPCErr(ctx.Err()))
			}
		}()
	}
	return cs, nil
}

func (cs *clientStream) newAttemptLocked(sh stats.Handler, trInfo *traceInfo) error {
	cs.attempt = &csAttempt{
		cs:           cs,
		dc:           cs.cc.dopts.dc,
		statsHandler: sh,
		trInfo:       trInfo,
	}

	if err := cs.ctx.Err(); err != nil {
		return toRPCErr(err)
	}
	t, done, err := cs.cc.getTransport(cs.ctx, cs.callInfo.failFast, cs.callHdr.Method)
	if err != nil {
		return err
	}
	if trInfo != nil {
		trInfo.firstLine.SetRemoteAddr(t.RemoteAddr())
	}
	cs.attempt.t = t
	cs.attempt.done = done
	return nil
}

func (a *csAttempt) newStream() error {
	cs := a.cs
	cs.callHdr.PreviousAttempts = cs.numRetries
	s, err := a.t.NewStream(cs.ctx, cs.callHdr)
	if err != nil {
		return toRPCErr(err)
	}
	cs.attempt.s = s
	cs.attempt.p = &parser{r: s}
	return nil
}


type clientStream struct {
	callHdr  *transport.CallHdr
	opts     []CallOption
	callInfo *callInfo
	cc       *ClientConn
	desc     *StreamDesc

	codec baseCodec
	cp    Compressor
	comp  encoding.Compressor

	cancel context.CancelFunc 

	sentLast  bool 
	beginTime time.Time

	methodConfig *MethodConfig

	ctx context.Context 

	retryThrottler *retryThrottler 

	binlog *binarylog.MethodLogger 
	
	
	
	
	
	
	serverHeaderBinlogged bool

	mu                      sync.Mutex
	firstAttempt            bool       
	numRetries              int        
	numRetriesSincePushback int        
	finished                bool       
	attempt                 *csAttempt 
	
	committed  bool                       
	buffer     []func(a *csAttempt) error 
	bufferSize int                        
}



type csAttempt struct {
	cs   *clientStream
	t    transport.ClientTransport
	s    *transport.Stream
	p    *parser
	done func(balancer.DoneInfo)

	finished  bool
	dc        Decompressor
	decomp    encoding.Compressor
	decompSet bool

	mu sync.Mutex 
	
	
	
	trInfo *traceInfo

	statsHandler stats.Handler
}

func (cs *clientStream) commitAttemptLocked() {
	cs.committed = true
	cs.buffer = nil
}

func (cs *clientStream) commitAttempt() {
	cs.mu.Lock()
	cs.commitAttemptLocked()
	cs.mu.Unlock()
}



func (cs *clientStream) shouldRetry(err error) error {
	if cs.attempt.s == nil && !cs.callInfo.failFast {
		
		
		
		return nil
	}
	if cs.finished || cs.committed {
		
		return err
	}
	
	if cs.attempt.s != nil {
		<-cs.attempt.s.Done()
	}
	if cs.firstAttempt && !cs.callInfo.failFast && (cs.attempt.s == nil || cs.attempt.s.Unprocessed()) {
		
		cs.firstAttempt = false
		return nil
	}
	cs.firstAttempt = false
	if cs.cc.dopts.disableRetry {
		return err
	}

	pushback := 0
	hasPushback := false
	if cs.attempt.s != nil {
		if to, toErr := cs.attempt.s.TrailersOnly(); toErr != nil || !to {
			return err
		}

		
		
		sps := cs.attempt.s.Trailer()["grpc-retry-pushback-ms"]
		if len(sps) == 1 {
			var e error
			if pushback, e = strconv.Atoi(sps[0]); e != nil || pushback < 0 {
				grpclog.Infof("Server retry pushback specified to abort (%q).", sps[0])
				cs.retryThrottler.throttle() 
				return err
			}
			hasPushback = true
		} else if len(sps) > 1 {
			grpclog.Warningf("Server retry pushback specified multiple values (%q); not retrying.", sps)
			cs.retryThrottler.throttle() 
			return err
		}
	}

	var code codes.Code
	if cs.attempt.s != nil {
		code = cs.attempt.s.Status().Code()
	} else {
		code = status.Convert(err).Code()
	}

	rp := cs.methodConfig.retryPolicy
	if rp == nil || !rp.retryableStatusCodes[code] {
		return err
	}

	
	
	if cs.retryThrottler.throttle() {
		return err
	}
	if cs.numRetries+1 >= rp.maxAttempts {
		return err
	}

	var dur time.Duration
	if hasPushback {
		dur = time.Millisecond * time.Duration(pushback)
		cs.numRetriesSincePushback = 0
	} else {
		fact := math.Pow(rp.backoffMultiplier, float64(cs.numRetriesSincePushback))
		cur := float64(rp.initialBackoff) * fact
		if max := float64(rp.maxBackoff); cur > max {
			cur = max
		}
		dur = time.Duration(grpcrand.Int63n(int64(cur)))
		cs.numRetriesSincePushback++
	}

	
	
	t := time.NewTimer(dur)
	select {
	case <-t.C:
		cs.numRetries++
		return nil
	case <-cs.ctx.Done():
		t.Stop()
		return status.FromContextError(cs.ctx.Err()).Err()
	}
}


func (cs *clientStream) retryLocked(lastErr error) error {
	for {
		cs.attempt.finish(lastErr)
		if err := cs.shouldRetry(lastErr); err != nil {
			cs.commitAttemptLocked()
			return err
		}
		if err := cs.newAttemptLocked(nil, nil); err != nil {
			return err
		}
		if lastErr = cs.replayBufferLocked(); lastErr == nil {
			return nil
		}
	}
}

func (cs *clientStream) Context() context.Context {
	cs.commitAttempt()
	
	
	return cs.attempt.s.Context()
}

func (cs *clientStream) withRetry(op func(a *csAttempt) error, onSuccess func()) error {
	cs.mu.Lock()
	for {
		if cs.committed {
			cs.mu.Unlock()
			return op(cs.attempt)
		}
		a := cs.attempt
		cs.mu.Unlock()
		err := op(a)
		cs.mu.Lock()
		if a != cs.attempt {
			
			continue
		}
		if err == io.EOF {
			<-a.s.Done()
		}
		if err == nil || (err == io.EOF && a.s.Status().Code() == codes.OK) {
			onSuccess()
			cs.mu.Unlock()
			return err
		}
		if err := cs.retryLocked(err); err != nil {
			cs.mu.Unlock()
			return err
		}
	}
}

func (cs *clientStream) Header() (metadata.MD, error) {
	var m metadata.MD
	err := cs.withRetry(func(a *csAttempt) error {
		var err error
		m, err = a.s.Header()
		return toRPCErr(err)
	}, cs.commitAttemptLocked)
	if err != nil {
		cs.finish(err)
		return nil, err
	}
	if cs.binlog != nil && !cs.serverHeaderBinlogged {
		
		logEntry := &binarylog.ServerHeader{
			OnClientSide: true,
			Header:       m,
			PeerAddr:     nil,
		}
		if peer, ok := peer.FromContext(cs.Context()); ok {
			logEntry.PeerAddr = peer.Addr
		}
		cs.binlog.Log(logEntry)
		cs.serverHeaderBinlogged = true
	}
	return m, err
}

func (cs *clientStream) Trailer() metadata.MD {
	
	
	
	
	
	
	
	cs.commitAttempt()
	if cs.attempt.s == nil {
		return nil
	}
	return cs.attempt.s.Trailer()
}

func (cs *clientStream) replayBufferLocked() error {
	a := cs.attempt
	for _, f := range cs.buffer {
		if err := f(a); err != nil {
			return err
		}
	}
	return nil
}

func (cs *clientStream) bufferForRetryLocked(sz int, op func(a *csAttempt) error) {
	
	if cs.committed {
		return
	}
	cs.bufferSize += sz
	if cs.bufferSize > cs.callInfo.maxRetryRPCBufferSize {
		cs.commitAttemptLocked()
		return
	}
	cs.buffer = append(cs.buffer, op)
}

func (cs *clientStream) SendMsg(m interface{}) (err error) {
	defer func() {
		if err != nil && err != io.EOF {
			
			
			
			
			
			cs.finish(err)
		}
	}()
	if cs.sentLast {
		return status.Errorf(codes.Internal, "SendMsg called after CloseSend")
	}
	if !cs.desc.ClientStreams {
		cs.sentLast = true
	}
	data, err := encode(cs.codec, m)
	if err != nil {
		return err
	}
	compData, err := compress(data, cs.cp, cs.comp)
	if err != nil {
		return err
	}
	hdr, payload := msgHeader(data, compData)
	
	if len(payload) > *cs.callInfo.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "trying to send message larger than max (%d vs. %d)", len(payload), *cs.callInfo.maxSendMessageSize)
	}
	msgBytes := data 
	op := func(a *csAttempt) error {
		err := a.sendMsg(m, hdr, payload, data)
		
		
		m, data = nil, nil
		return err
	}
	err = cs.withRetry(op, func() { cs.bufferForRetryLocked(len(hdr)+len(payload), op) })
	if cs.binlog != nil && err == nil {
		cs.binlog.Log(&binarylog.ClientMessage{
			OnClientSide: true,
			Message:      msgBytes,
		})
	}
	return
}

func (cs *clientStream) RecvMsg(m interface{}) error {
	if cs.binlog != nil && !cs.serverHeaderBinlogged {
		
		cs.Header()
	}
	var recvInfo *payloadInfo
	if cs.binlog != nil {
		recvInfo = &payloadInfo{}
	}
	err := cs.withRetry(func(a *csAttempt) error {
		return a.recvMsg(m, recvInfo)
	}, cs.commitAttemptLocked)
	if cs.binlog != nil && err == nil {
		cs.binlog.Log(&binarylog.ServerMessage{
			OnClientSide: true,
			Message:      recvInfo.uncompressedBytes,
		})
	}
	if err != nil || !cs.desc.ServerStreams {
		
		cs.finish(err)

		if cs.binlog != nil {
			
			logEntry := &binarylog.ServerTrailer{
				OnClientSide: true,
				Trailer:      cs.Trailer(),
				Err:          err,
			}
			if logEntry.Err == io.EOF {
				logEntry.Err = nil
			}
			if peer, ok := peer.FromContext(cs.Context()); ok {
				logEntry.PeerAddr = peer.Addr
			}
			cs.binlog.Log(logEntry)
		}
	}
	return err
}

func (cs *clientStream) CloseSend() error {
	if cs.sentLast {
		
		return nil
	}
	cs.sentLast = true
	op := func(a *csAttempt) error {
		a.t.Write(a.s, nil, nil, &transport.Options{Last: true})
		
		
		
		
		return nil
	}
	cs.withRetry(op, func() { cs.bufferForRetryLocked(0, op) })
	if cs.binlog != nil {
		cs.binlog.Log(&binarylog.ClientHalfClose{
			OnClientSide: true,
		})
	}
	
	return nil
}

func (cs *clientStream) finish(err error) {
	if err == io.EOF {
		
		err = nil
	}
	cs.mu.Lock()
	if cs.finished {
		cs.mu.Unlock()
		return
	}
	cs.finished = true
	cs.commitAttemptLocked()
	cs.mu.Unlock()
	
	
	
	
	
	if cs.binlog != nil && status.Code(err) == codes.Canceled {
		cs.binlog.Log(&binarylog.Cancel{
			OnClientSide: true,
		})
	}
	if err == nil {
		cs.retryThrottler.successfulRPC()
	}
	if channelz.IsOn() {
		if err != nil {
			cs.cc.incrCallsFailed()
		} else {
			cs.cc.incrCallsSucceeded()
		}
	}
	if cs.attempt != nil {
		cs.attempt.finish(err)
	}
	
	if cs.attempt.s != nil {
		for _, o := range cs.opts {
			o.after(cs.callInfo)
		}
	}
	cs.cancel()
}

func (a *csAttempt) sendMsg(m interface{}, hdr, payld, data []byte) error {
	cs := a.cs
	if a.trInfo != nil {
		a.mu.Lock()
		if a.trInfo.tr != nil {
			a.trInfo.tr.LazyLog(&payload{sent: true, msg: m}, true)
		}
		a.mu.Unlock()
	}
	if err := a.t.Write(a.s, hdr, payld, &transport.Options{Last: !cs.desc.ClientStreams}); err != nil {
		if !cs.desc.ClientStreams {
			
			
			
			return nil
		}
		return io.EOF
	}
	if a.statsHandler != nil {
		a.statsHandler.HandleRPC(cs.ctx, outPayload(true, m, data, payld, time.Now()))
	}
	if channelz.IsOn() {
		a.t.IncrMsgSent()
	}
	return nil
}

func (a *csAttempt) recvMsg(m interface{}, payInfo *payloadInfo) (err error) {
	cs := a.cs
	if a.statsHandler != nil && payInfo == nil {
		payInfo = &payloadInfo{}
	}

	if !a.decompSet {
		
		if ct := a.s.RecvCompress(); ct != "" && ct != encoding.Identity {
			if a.dc == nil || a.dc.Type() != ct {
				
				
				a.dc = nil
				a.decomp = encoding.GetCompressor(ct)
			}
		} else {
			
			a.dc = nil
		}
		
		a.decompSet = true
	}
	err = recv(a.p, cs.codec, a.s, a.dc, m, *cs.callInfo.maxReceiveMessageSize, payInfo, a.decomp)
	if err != nil {
		if err == io.EOF {
			if statusErr := a.s.Status().Err(); statusErr != nil {
				return statusErr
			}
			return io.EOF 
		}
		return toRPCErr(err)
	}
	if a.trInfo != nil {
		a.mu.Lock()
		if a.trInfo.tr != nil {
			a.trInfo.tr.LazyLog(&payload{sent: false, msg: m}, true)
		}
		a.mu.Unlock()
	}
	if a.statsHandler != nil {
		a.statsHandler.HandleRPC(cs.ctx, &stats.InPayload{
			Client:   true,
			RecvTime: time.Now(),
			Payload:  m,
			
			Data:       payInfo.uncompressedBytes,
			WireLength: payInfo.wireLength,
			Length:     len(payInfo.uncompressedBytes),
		})
	}
	if channelz.IsOn() {
		a.t.IncrMsgRecv()
	}
	if cs.desc.ServerStreams {
		
		return nil
	}
	
	
	err = recv(a.p, cs.codec, a.s, a.dc, m, *cs.callInfo.maxReceiveMessageSize, nil, a.decomp)
	if err == nil {
		return toRPCErr(errors.New("grpc: client streaming protocol violation: get <nil>, want <EOF>"))
	}
	if err == io.EOF {
		return a.s.Status().Err() 
	}
	return toRPCErr(err)
}

func (a *csAttempt) finish(err error) {
	a.mu.Lock()
	if a.finished {
		a.mu.Unlock()
		return
	}
	a.finished = true
	if err == io.EOF {
		
		err = nil
	}
	var tr metadata.MD
	if a.s != nil {
		a.t.CloseStream(a.s, err)
		tr = a.s.Trailer()
	}

	if a.done != nil {
		br := false
		if a.s != nil {
			br = a.s.BytesReceived()
		}
		a.done(balancer.DoneInfo{
			Err:           err,
			Trailer:       tr,
			BytesSent:     a.s != nil,
			BytesReceived: br,
			ServerLoad:    balancerload.Parse(tr),
		})
	}
	if a.statsHandler != nil {
		end := &stats.End{
			Client:    true,
			BeginTime: a.cs.beginTime,
			EndTime:   time.Now(),
			Trailer:   tr,
			Error:     err,
		}
		a.statsHandler.HandleRPC(a.cs.ctx, end)
	}
	if a.trInfo != nil && a.trInfo.tr != nil {
		if err == nil {
			a.trInfo.tr.LazyPrintf("RPC: [OK]")
		} else {
			a.trInfo.tr.LazyPrintf("RPC: [%v]", err)
			a.trInfo.tr.SetError()
		}
		a.trInfo.tr.Finish()
		a.trInfo.tr = nil
	}
	a.mu.Unlock()
}

func (ac *addrConn) newClientStream(ctx context.Context, desc *StreamDesc, method string, t transport.ClientTransport, opts ...CallOption) (_ ClientStream, err error) {
	ac.mu.Lock()
	if ac.transport != t {
		ac.mu.Unlock()
		return nil, status.Error(codes.Canceled, "the provided transport is no longer valid to use")
	}
	
	if ac.state != connectivity.Connecting {
		ac.updateConnectivityState(connectivity.Connecting)
		ac.cc.handleSubConnStateChange(ac.acbw, ac.state)
	}
	ac.mu.Unlock()

	if t == nil {
		
		return nil, errors.New("transport provided is nil")
	}
	
	c := &callInfo{}

	for _, o := range opts {
		if err := o.before(c); err != nil {
			return nil, toRPCErr(err)
		}
	}
	c.maxReceiveMessageSize = getMaxSize(nil, c.maxReceiveMessageSize, defaultClientMaxReceiveMessageSize)
	c.maxSendMessageSize = getMaxSize(nil, c.maxSendMessageSize, defaultServerMaxSendMessageSize)

	
	
	
	
	
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	if err := setCallInfoCodec(c); err != nil {
		return nil, err
	}

	callHdr := &transport.CallHdr{
		Host:           ac.cc.authority,
		Method:         method,
		ContentSubtype: c.contentSubtype,
	}

	
	
	
	
	var cp Compressor
	var comp encoding.Compressor
	if ct := c.compressorType; ct != "" {
		callHdr.SendCompress = ct
		if ct != encoding.Identity {
			comp = encoding.GetCompressor(ct)
			if comp == nil {
				return nil, status.Errorf(codes.Internal, "grpc: Compressor is not installed for requested grpc-encoding %q", ct)
			}
		}
	} else if ac.cc.dopts.cp != nil {
		callHdr.SendCompress = ac.cc.dopts.cp.Type()
		cp = ac.cc.dopts.cp
	}
	if c.creds != nil {
		callHdr.Creds = c.creds
	}

	as := &addrConnStream{
		callHdr:  callHdr,
		ac:       ac,
		ctx:      ctx,
		cancel:   cancel,
		opts:     opts,
		callInfo: c,
		desc:     desc,
		codec:    c.codec,
		cp:       cp,
		comp:     comp,
		t:        t,
	}

	as.callInfo.stream = as
	s, err := as.t.NewStream(as.ctx, as.callHdr)
	if err != nil {
		err = toRPCErr(err)
		return nil, err
	}
	as.s = s
	as.p = &parser{r: s}
	ac.incrCallsStarted()
	if desc != unaryStreamDesc {
		
		
		
		
		
		go func() {
			select {
			case <-ac.ctx.Done():
				as.finish(status.Error(codes.Canceled, "grpc: the SubConn is closing"))
			case <-ctx.Done():
				as.finish(toRPCErr(ctx.Err()))
			}
		}()
	}
	return as, nil
}

type addrConnStream struct {
	s         *transport.Stream
	ac        *addrConn
	callHdr   *transport.CallHdr
	cancel    context.CancelFunc
	opts      []CallOption
	callInfo  *callInfo
	t         transport.ClientTransport
	ctx       context.Context
	sentLast  bool
	desc      *StreamDesc
	codec     baseCodec
	cp        Compressor
	comp      encoding.Compressor
	decompSet bool
	dc        Decompressor
	decomp    encoding.Compressor
	p         *parser
	mu        sync.Mutex
	finished  bool
}

func (as *addrConnStream) Header() (metadata.MD, error) {
	m, err := as.s.Header()
	if err != nil {
		as.finish(toRPCErr(err))
	}
	return m, err
}

func (as *addrConnStream) Trailer() metadata.MD {
	return as.s.Trailer()
}

func (as *addrConnStream) CloseSend() error {
	if as.sentLast {
		
		return nil
	}
	as.sentLast = true

	as.t.Write(as.s, nil, nil, &transport.Options{Last: true})
	
	
	
	
	return nil
}

func (as *addrConnStream) Context() context.Context {
	return as.s.Context()
}

func (as *addrConnStream) SendMsg(m interface{}) (err error) {
	defer func() {
		if err != nil && err != io.EOF {
			
			
			
			
			
			as.finish(err)
		}
	}()
	if as.sentLast {
		return status.Errorf(codes.Internal, "SendMsg called after CloseSend")
	}
	if !as.desc.ClientStreams {
		as.sentLast = true
	}
	data, err := encode(as.codec, m)
	if err != nil {
		return err
	}
	compData, err := compress(data, as.cp, as.comp)
	if err != nil {
		return err
	}
	hdr, payld := msgHeader(data, compData)
	
	if len(payld) > *as.callInfo.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "trying to send message larger than max (%d vs. %d)", len(payld), *as.callInfo.maxSendMessageSize)
	}

	if err := as.t.Write(as.s, hdr, payld, &transport.Options{Last: !as.desc.ClientStreams}); err != nil {
		if !as.desc.ClientStreams {
			
			
			
			return nil
		}
		return io.EOF
	}

	if channelz.IsOn() {
		as.t.IncrMsgSent()
	}
	return nil
}

func (as *addrConnStream) RecvMsg(m interface{}) (err error) {
	defer func() {
		if err != nil || !as.desc.ServerStreams {
			
			as.finish(err)
		}
	}()

	if !as.decompSet {
		
		if ct := as.s.RecvCompress(); ct != "" && ct != encoding.Identity {
			if as.dc == nil || as.dc.Type() != ct {
				
				
				as.dc = nil
				as.decomp = encoding.GetCompressor(ct)
			}
		} else {
			
			as.dc = nil
		}
		
		as.decompSet = true
	}
	err = recv(as.p, as.codec, as.s, as.dc, m, *as.callInfo.maxReceiveMessageSize, nil, as.decomp)
	if err != nil {
		if err == io.EOF {
			if statusErr := as.s.Status().Err(); statusErr != nil {
				return statusErr
			}
			return io.EOF 
		}
		return toRPCErr(err)
	}

	if channelz.IsOn() {
		as.t.IncrMsgRecv()
	}
	if as.desc.ServerStreams {
		
		return nil
	}

	
	
	err = recv(as.p, as.codec, as.s, as.dc, m, *as.callInfo.maxReceiveMessageSize, nil, as.decomp)
	if err == nil {
		return toRPCErr(errors.New("grpc: client streaming protocol violation: get <nil>, want <EOF>"))
	}
	if err == io.EOF {
		return as.s.Status().Err() 
	}
	return toRPCErr(err)
}

func (as *addrConnStream) finish(err error) {
	as.mu.Lock()
	if as.finished {
		as.mu.Unlock()
		return
	}
	as.finished = true
	if err == io.EOF {
		
		err = nil
	}
	if as.s != nil {
		as.t.CloseStream(as.s, err)
	}

	if err != nil {
		as.ac.incrCallsFailed()
	} else {
		as.ac.incrCallsSucceeded()
	}
	as.cancel()
	as.mu.Unlock()
}





type ServerStream interface {
	
	
	
	
	
	
	SetHeader(metadata.MD) error
	
	
	
	SendHeader(metadata.MD) error
	
	
	SetTrailer(metadata.MD)
	
	Context() context.Context
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	SendMsg(m interface{}) error
	
	
	
	
	
	
	
	
	RecvMsg(m interface{}) error
}


type serverStream struct {
	ctx   context.Context
	t     transport.ServerTransport
	s     *transport.Stream
	p     *parser
	codec baseCodec

	cp     Compressor
	dc     Decompressor
	comp   encoding.Compressor
	decomp encoding.Compressor

	maxReceiveMessageSize int
	maxSendMessageSize    int
	trInfo                *traceInfo

	statsHandler stats.Handler

	binlog *binarylog.MethodLogger
	
	
	
	
	
	
	serverHeaderBinlogged bool

	mu sync.Mutex 
}

func (ss *serverStream) Context() context.Context {
	return ss.ctx
}

func (ss *serverStream) SetHeader(md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	return ss.s.SetHeader(md)
}

func (ss *serverStream) SendHeader(md metadata.MD) error {
	err := ss.t.WriteHeader(ss.s, md)
	if ss.binlog != nil && !ss.serverHeaderBinlogged {
		h, _ := ss.s.Header()
		ss.binlog.Log(&binarylog.ServerHeader{
			Header: h,
		})
		ss.serverHeaderBinlogged = true
	}
	return err
}

func (ss *serverStream) SetTrailer(md metadata.MD) {
	if md.Len() == 0 {
		return
	}
	ss.s.SetTrailer(md)
}

func (ss *serverStream) SendMsg(m interface{}) (err error) {
	defer func() {
		if ss.trInfo != nil {
			ss.mu.Lock()
			if ss.trInfo.tr != nil {
				if err == nil {
					ss.trInfo.tr.LazyLog(&payload{sent: true, msg: m}, true)
				} else {
					ss.trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
					ss.trInfo.tr.SetError()
				}
			}
			ss.mu.Unlock()
		}
		if err != nil && err != io.EOF {
			st, _ := status.FromError(toRPCErr(err))
			ss.t.WriteStatus(ss.s, st)
			
			
			
			
			
			
		}
		if channelz.IsOn() && err == nil {
			ss.t.IncrMsgSent()
		}
	}()
	data, err := encode(ss.codec, m)
	if err != nil {
		return err
	}
	compData, err := compress(data, ss.cp, ss.comp)
	if err != nil {
		return err
	}
	hdr, payload := msgHeader(data, compData)
	
	if len(payload) > ss.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "trying to send message larger than max (%d vs. %d)", len(payload), ss.maxSendMessageSize)
	}
	if err := ss.t.Write(ss.s, hdr, payload, &transport.Options{Last: false}); err != nil {
		return toRPCErr(err)
	}
	if ss.binlog != nil {
		if !ss.serverHeaderBinlogged {
			h, _ := ss.s.Header()
			ss.binlog.Log(&binarylog.ServerHeader{
				Header: h,
			})
			ss.serverHeaderBinlogged = true
		}
		ss.binlog.Log(&binarylog.ServerMessage{
			Message: data,
		})
	}
	if ss.statsHandler != nil {
		ss.statsHandler.HandleRPC(ss.s.Context(), outPayload(false, m, data, payload, time.Now()))
	}
	return nil
}

func (ss *serverStream) RecvMsg(m interface{}) (err error) {
	defer func() {
		if ss.trInfo != nil {
			ss.mu.Lock()
			if ss.trInfo.tr != nil {
				if err == nil {
					ss.trInfo.tr.LazyLog(&payload{sent: false, msg: m}, true)
				} else if err != io.EOF {
					ss.trInfo.tr.LazyLog(&fmtStringer{"%v", []interface{}{err}}, true)
					ss.trInfo.tr.SetError()
				}
			}
			ss.mu.Unlock()
		}
		if err != nil && err != io.EOF {
			st, _ := status.FromError(toRPCErr(err))
			ss.t.WriteStatus(ss.s, st)
			
			
			
			
			
			
		}
		if channelz.IsOn() && err == nil {
			ss.t.IncrMsgRecv()
		}
	}()
	var payInfo *payloadInfo
	if ss.statsHandler != nil || ss.binlog != nil {
		payInfo = &payloadInfo{}
	}
	if err := recv(ss.p, ss.codec, ss.s, ss.dc, m, ss.maxReceiveMessageSize, payInfo, ss.decomp); err != nil {
		if err == io.EOF {
			if ss.binlog != nil {
				ss.binlog.Log(&binarylog.ClientHalfClose{})
			}
			return err
		}
		if err == io.ErrUnexpectedEOF {
			err = status.Errorf(codes.Internal, io.ErrUnexpectedEOF.Error())
		}
		return toRPCErr(err)
	}
	if ss.statsHandler != nil {
		ss.statsHandler.HandleRPC(ss.s.Context(), &stats.InPayload{
			RecvTime: time.Now(),
			Payload:  m,
			
			Data:       payInfo.uncompressedBytes,
			WireLength: payInfo.wireLength,
			Length:     len(payInfo.uncompressedBytes),
		})
	}
	if ss.binlog != nil {
		ss.binlog.Log(&binarylog.ClientMessage{
			Message: payInfo.uncompressedBytes,
		})
	}
	return nil
}



func MethodFromServerStream(stream ServerStream) (string, bool) {
	return Method(stream.Context())
}
