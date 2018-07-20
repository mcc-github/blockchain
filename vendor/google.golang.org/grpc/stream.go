

package grpc

import (
	"errors"
	"io"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/transport"
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
	
	
	
	
	
	
	Stream
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
		Host:   cc.authority,
		Method: method,
		
		
		
		
		Flush:          desc.ClientStreams,
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
		defer func() {
			if err != nil {
				
				
				trInfo.tr.LazyPrintf("RPC: [%v]", err)
				trInfo.tr.SetError()
				trInfo.tr.Finish()
			}
		}()
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
		defer func() {
			if err != nil {
				
				end := &stats.End{
					Client:    true,
					Error:     err,
					BeginTime: beginTime,
					EndTime:   time.Now(),
				}
				sh.HandleRPC(ctx, end)
			}
		}()
	}

	var (
		t    transport.ClientTransport
		s    *transport.Stream
		done func(balancer.DoneInfo)
	)
	for {
		
		
		
		select {
		case <-ctx.Done():
			return nil, toRPCErr(ctx.Err())
		default:
		}

		t, done, err = cc.getTransport(ctx, c.failFast)
		if err != nil {
			return nil, err
		}

		s, err = t.NewStream(ctx, callHdr)
		if err != nil {
			if done != nil {
				done(balancer.DoneInfo{Err: err})
				done = nil
			}
			
			
			
			if !c.failFast {
				continue
			}
			return nil, toRPCErr(err)
		}
		break
	}

	cs := &clientStream{
		opts:   opts,
		c:      c,
		cc:     cc,
		desc:   desc,
		codec:  c.codec,
		cp:     cp,
		comp:   comp,
		cancel: cancel,
		attempt: &csAttempt{
			t:            t,
			s:            s,
			p:            &parser{r: s},
			done:         done,
			dc:           cc.dopts.dc,
			ctx:          ctx,
			trInfo:       trInfo,
			statsHandler: sh,
			beginTime:    beginTime,
		},
	}
	cs.c.stream = cs
	cs.attempt.cs = cs
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


type clientStream struct {
	opts []CallOption
	c    *callInfo
	cc   *ClientConn
	desc *StreamDesc

	codec baseCodec
	cp    Compressor
	comp  encoding.Compressor

	cancel context.CancelFunc 

	sentLast bool 

	mu       sync.Mutex 
	finished bool       

	attempt *csAttempt 
	
}



type csAttempt struct {
	cs   *clientStream
	t    transport.ClientTransport
	s    *transport.Stream
	p    *parser
	done func(balancer.DoneInfo)

	dc        Decompressor
	decomp    encoding.Compressor
	decompSet bool

	ctx context.Context 

	mu sync.Mutex 
	
	
	trInfo traceInfo

	statsHandler stats.Handler
	beginTime    time.Time
}

func (cs *clientStream) Context() context.Context {
	
	return cs.attempt.context()
}

func (cs *clientStream) Header() (metadata.MD, error) {
	m, err := cs.attempt.header()
	if err != nil {
		
		err = toRPCErr(err)
		cs.finish(err)
	}
	return m, err
}

func (cs *clientStream) Trailer() metadata.MD {
	
	return cs.attempt.trailer()
}

func (cs *clientStream) SendMsg(m interface{}) (err error) {
	
	return cs.attempt.sendMsg(m)
}

func (cs *clientStream) RecvMsg(m interface{}) (err error) {
	
	return cs.attempt.recvMsg(m)
}

func (cs *clientStream) CloseSend() error {
	cs.attempt.closeSend()
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
	cs.mu.Unlock()
	if channelz.IsOn() {
		if err != nil {
			cs.cc.incrCallsFailed()
		} else {
			cs.cc.incrCallsSucceeded()
		}
	}
	
	cs.attempt.finish(err)
	for _, o := range cs.opts {
		o.after(cs.c)
	}
	cs.cancel()
}

func (a *csAttempt) context() context.Context {
	return a.s.Context()
}

func (a *csAttempt) header() (metadata.MD, error) {
	return a.s.Header()
}

func (a *csAttempt) trailer() metadata.MD {
	return a.s.Trailer()
}

func (a *csAttempt) sendMsg(m interface{}) (err error) {
	
	
	cs := a.cs
	defer func() {
		
		
		
		if err == io.EOF && !cs.desc.ClientStreams {
			err = nil
		}
		if err != nil && err != io.EOF {
			
			
			
			
			cs.finish(err)
		}
	}()
	
	if EnableTracing {
		a.mu.Lock()
		if a.trInfo.tr != nil {
			a.trInfo.tr.LazyLog(&payload{sent: true, msg: m}, true)
		}
		a.mu.Unlock()
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
	
	if len(payload) > *cs.c.maxSendMessageSize {
		return status.Errorf(codes.ResourceExhausted, "trying to send message larger than max (%d vs. %d)", len(payload), *cs.c.maxSendMessageSize)
	}

	if !cs.desc.ClientStreams {
		cs.sentLast = true
	}
	err = a.t.Write(a.s, hdr, payload, &transport.Options{Last: !cs.desc.ClientStreams})
	if err == nil {
		if a.statsHandler != nil {
			a.statsHandler.HandleRPC(a.ctx, outPayload(true, m, data, payload, time.Now()))
		}
		if channelz.IsOn() {
			a.t.IncrMsgSent()
		}
		return nil
	}
	return io.EOF
}

func (a *csAttempt) recvMsg(m interface{}) (err error) {
	cs := a.cs
	defer func() {
		if err != nil || !cs.desc.ServerStreams {
			
			cs.finish(err)
		}
	}()
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
	err = recv(a.p, cs.codec, a.s, a.dc, m, *cs.c.maxReceiveMessageSize, inPayload, a.decomp)
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
		a.statsHandler.HandleRPC(a.ctx, inPayload)
	}
	if channelz.IsOn() {
		a.t.IncrMsgRecv()
	}
	if cs.desc.ServerStreams {
		
		return nil
	}

	
	
	err = recv(a.p, cs.codec, a.s, a.dc, m, *cs.c.maxReceiveMessageSize, nil, a.decomp)
	if err == nil {
		return toRPCErr(errors.New("grpc: client streaming protocol violation: get <nil>, want <EOF>"))
	}
	if err == io.EOF {
		return a.s.Status().Err() 
	}
	return toRPCErr(err)
}

func (a *csAttempt) closeSend() {
	cs := a.cs
	if cs.sentLast {
		return
	}
	cs.sentLast = true
	cs.attempt.t.Write(cs.attempt.s, nil, nil, &transport.Options{Last: true})
	
	
	
}

func (a *csAttempt) finish(err error) {
	a.mu.Lock()
	a.t.CloseStream(a.s, err)

	if a.done != nil {
		a.done(balancer.DoneInfo{
			Err:           err,
			BytesSent:     true,
			BytesReceived: a.s.BytesReceived(),
		})
	}
	if a.statsHandler != nil {
		end := &stats.End{
			Client:    true,
			BeginTime: a.beginTime,
			EndTime:   time.Now(),
			Error:     err,
		}
		a.statsHandler.HandleRPC(a.ctx, end)
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
	Stream
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
