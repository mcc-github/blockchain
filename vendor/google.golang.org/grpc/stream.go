

package grpc

import (
	"errors"
	"io"
	"math"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/grpcrand"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/metadata"
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
	var trInfo traceInfo
	if EnableTracing {
		trInfo.tr = trace.New("grpc.Sent."+methodFamily(method), method)
		trInfo.firstLine.client = true
		if deadline, ok := ctx.Deadline(); ok {
			trInfo.firstLine.deadline = deadline.Sub(time.Now())
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

func (cs *clientStream) newAttemptLocked(sh stats.Handler, trInfo traceInfo) error {
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
	
	
	trInfo traceInfo

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
		if to, toErr := cs.attempt.s.TrailersOnly(); toErr != nil {
			
			return toErr
		} else if !to {
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
		if err := cs.newAttemptLocked(nil, traceInfo{}); err != nil {
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
	op := func(a *csAttempt) error {
		err := a.sendMsg(m, hdr, payload, data)
		
		
		m, data = nil, nil
		return err
	}
	return cs.withRetry(op, func() { cs.bufferForRetryLocked(len(hdr)+len(payload), op) })
}

func (cs *clientStream) RecvMsg(m interface{}) error {
	err := cs.withRetry(func(a *csAttempt) error {
		return a.recvMsg(m)
	}, cs.commitAttemptLocked)
	if err != nil || !cs.desc.ServerStreams {
		
		cs.finish(err)
	}
	return err
}

func (cs *clientStream) CloseSend() error {
	if cs.sentLast {
		
		return nil
	}
	cs.sentLast = true
	op := func(a *csAttempt) error { return a.t.Write(a.s, nil, nil, &transport.Options{Last: true}) }
	cs.withRetry(op, func() { cs.bufferForRetryLocked(0, op) })
	
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
	if EnableTracing {
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

func (a *csAttempt) recvMsg(m interface{}) (err error) {
	cs := a.cs
	var inPayload *stats.InPayload
	if a.statsHandler != nil {
		inPayload = &stats.InPayload{
			Client: true,
		}
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
	err = recv(a.p, cs.codec, a.s, a.dc, m, *cs.callInfo.maxReceiveMessageSize, inPayload, a.decomp)
	if err != nil {
		if err == io.EOF {
			if statusErr := a.s.Status().Err(); statusErr != nil {
				return statusErr
			}
			return io.EOF 
		}
		return toRPCErr(err)
	}
	if EnableTracing {
		a.mu.Lock()
		if a.trInfo.tr != nil {
			a.trInfo.tr.LazyLog(&payload{sent: false, msg: m}, true)
		}
		a.mu.Unlock()
	}
	if inPayload != nil {
		a.statsHandler.HandleRPC(cs.ctx, inPayload)
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
	if a.s != nil {
		a.t.CloseStream(a.s, err)
	}

	if a.done != nil {
		br := false
		if a.s != nil {
			br = a.s.BytesReceived()
		}
		a.done(balancer.DoneInfo{
			Err:           err,
			BytesSent:     a.s != nil,
			BytesReceived: br,
		})
	}
	if a.statsHandler != nil {
		end := &stats.End{
			Client:    true,
			BeginTime: a.cs.beginTime,
			EndTime:   time.Now(),
			Error:     err,
		}
		a.statsHandler.HandleRPC(a.cs.ctx, end)
	}
	if a.trInfo.tr != nil {
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
	return ss.t.WriteHeader(ss.s, md)
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
	var inPayload *stats.InPayload
	if ss.statsHandler != nil {
		inPayload = &stats.InPayload{}
	}
	if err := recv(ss.p, ss.codec, ss.s, ss.dc, m, ss.maxReceiveMessageSize, inPayload, ss.decomp); err != nil {
		if err == io.EOF {
			return err
		}
		if err == io.ErrUnexpectedEOF {
			err = status.Errorf(codes.Internal, io.ErrUnexpectedEOF.Error())
		}
		return toRPCErr(err)
	}
	if inPayload != nil {
		ss.statsHandler.HandleRPC(ss.s.Context(), inPayload)
	}
	return nil
}



func MethodFromServerStream(stream ServerStream) (string, bool) {
	return Method(stream.Context())
}
