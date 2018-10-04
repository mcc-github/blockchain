

package grpc

import (
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc/balancer"
	_ "google.golang.org/grpc/balancer/roundrobin" 
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	_ "google.golang.org/grpc/resolver/dns"         
	_ "google.golang.org/grpc/resolver/passthrough" 
	"google.golang.org/grpc/status"
)

const (
	
	minConnectTimeout = 20 * time.Second
	
	grpclbName = "grpclb"
)

var (
	
	
	
	
	
	ErrClientConnClosing = status.Error(codes.Canceled, "grpc: the client connection is closing")
	
	errConnDrain = errors.New("grpc: the connection is drained")
	
	errConnClosing = errors.New("grpc: the connection is closing")
	
	errBalancerClosed = errors.New("grpc: balancer is closed")
	
	
	getMinConnectTimeout = func() time.Duration {
		return minConnectTimeout
	}
)


var (
	
	
	
	errNoTransportSecurity = errors.New("grpc: no transport security set (use grpc.WithInsecure() explicitly or set credentials)")
	
	
	
	errTransportCredentialsMissing = errors.New("grpc: the credentials require transport level security (use grpc.WithTransportCredentials() to set)")
	
	
	errCredentialsConflict = errors.New("grpc: transport credentials are set for an insecure connection (grpc.WithTransportCredentials() and grpc.WithInsecure() are both called)")
)

const (
	defaultClientMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultClientMaxSendMessageSize    = math.MaxInt32
	
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)


func Dial(target string, opts ...DialOption) (*ClientConn, error) {
	return DialContext(context.Background(), target, opts...)
}

















func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
	cc := &ClientConn{
		target:         target,
		csMgr:          &connectivityStateManager{},
		conns:          make(map[*addrConn]struct{}),
		dopts:          defaultDialOptions(),
		blockingpicker: newPickerWrapper(),
		czData:         new(channelzData),
	}
	cc.retryThrottler.Store((*retryThrottler)(nil))
	cc.ctx, cc.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt.apply(&cc.dopts)
	}

	if channelz.IsOn() {
		if cc.dopts.channelzParentID != 0 {
			cc.channelzID = channelz.RegisterChannel(&channelzChannel{cc}, cc.dopts.channelzParentID, target)
		} else {
			cc.channelzID = channelz.RegisterChannel(&channelzChannel{cc}, 0, target)
		}
	}

	if !cc.dopts.insecure {
		if cc.dopts.copts.TransportCredentials == nil {
			return nil, errNoTransportSecurity
		}
	} else {
		if cc.dopts.copts.TransportCredentials != nil {
			return nil, errCredentialsConflict
		}
		for _, cd := range cc.dopts.copts.PerRPCCredentials {
			if cd.RequireTransportSecurity() {
				return nil, errTransportCredentialsMissing
			}
		}
	}

	cc.mkp = cc.dopts.copts.KeepaliveParams

	if cc.dopts.copts.Dialer == nil {
		cc.dopts.copts.Dialer = newProxyDialer(
			func(ctx context.Context, addr string) (net.Conn, error) {
				network, addr := parseDialTarget(addr)
				return dialContext(ctx, network, addr)
			},
		)
	}

	if cc.dopts.copts.UserAgent != "" {
		cc.dopts.copts.UserAgent += " " + grpcUA
	} else {
		cc.dopts.copts.UserAgent = grpcUA
	}

	if cc.dopts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cc.dopts.timeout)
		defer cancel()
	}

	defer func() {
		select {
		case <-ctx.Done():
			conn, err = nil, ctx.Err()
		default:
		}

		if err != nil {
			cc.Close()
		}
	}()

	scSet := false
	if cc.dopts.scChan != nil {
		
		select {
		case sc, ok := <-cc.dopts.scChan:
			if ok {
				cc.sc = sc
				scSet = true
			}
		default:
		}
	}
	if cc.dopts.bs == nil {
		cc.dopts.bs = backoff.Exponential{
			MaxDelay: DefaultBackoffConfig.MaxDelay,
		}
	}
	if cc.dopts.resolverBuilder == nil {
		
		cc.parsedTarget = parseTarget(cc.target)
		grpclog.Infof("parsed scheme: %q", cc.parsedTarget.Scheme)
		cc.dopts.resolverBuilder = resolver.Get(cc.parsedTarget.Scheme)
		if cc.dopts.resolverBuilder == nil {
			
			
			
			grpclog.Infof("scheme %q not registered, fallback to default scheme", cc.parsedTarget.Scheme)
			cc.parsedTarget = resolver.Target{
				Scheme:   resolver.GetDefaultScheme(),
				Endpoint: target,
			}
			cc.dopts.resolverBuilder = resolver.Get(cc.parsedTarget.Scheme)
		}
	} else {
		cc.parsedTarget = resolver.Target{Endpoint: target}
	}
	creds := cc.dopts.copts.TransportCredentials
	if creds != nil && creds.Info().ServerName != "" {
		cc.authority = creds.Info().ServerName
	} else if cc.dopts.insecure && cc.dopts.authority != "" {
		cc.authority = cc.dopts.authority
	} else {
		
		
		cc.authority = cc.parsedTarget.Endpoint
	}

	if cc.dopts.scChan != nil && !scSet {
		
		select {
		case sc, ok := <-cc.dopts.scChan:
			if ok {
				cc.sc = sc
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if cc.dopts.scChan != nil {
		go cc.scWatcher()
	}

	var credsClone credentials.TransportCredentials
	if creds := cc.dopts.copts.TransportCredentials; creds != nil {
		credsClone = creds.Clone()
	}
	cc.balancerBuildOpts = balancer.BuildOptions{
		DialCreds:        credsClone,
		Dialer:           cc.dopts.copts.Dialer,
		ChannelzParentID: cc.channelzID,
	}

	
	cc.resolverWrapper, err = newCCResolverWrapper(cc)
	if err != nil {
		return nil, fmt.Errorf("failed to build resolver: %v", err)
	}
	
	
	
	
	
	
	
	cc.resolverWrapper.start()

	
	if cc.dopts.block {
		for {
			s := cc.GetState()
			if s == connectivity.Ready {
				break
			} else if cc.dopts.copts.FailOnNonTempDialError && s == connectivity.TransientFailure {
				if err = cc.blockingpicker.connectionError(); err != nil {
					terr, ok := err.(interface{ Temporary() bool })
					if ok && !terr.Temporary() {
						return nil, err
					}
				}
			}
			if !cc.WaitForStateChange(ctx, s) {
				
				return nil, ctx.Err()
			}
		}
	}

	return cc, nil
}



type connectivityStateManager struct {
	mu         sync.Mutex
	state      connectivity.State
	notifyChan chan struct{}
}




func (csm *connectivityStateManager) updateState(state connectivity.State) {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	if csm.state == connectivity.Shutdown {
		return
	}
	if csm.state == state {
		return
	}
	csm.state = state
	if csm.notifyChan != nil {
		
		close(csm.notifyChan)
		csm.notifyChan = nil
	}
}

func (csm *connectivityStateManager) getState() connectivity.State {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	return csm.state
}

func (csm *connectivityStateManager) getNotifyChan() <-chan struct{} {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	if csm.notifyChan == nil {
		csm.notifyChan = make(chan struct{})
	}
	return csm.notifyChan
}


type ClientConn struct {
	ctx    context.Context
	cancel context.CancelFunc

	target       string
	parsedTarget resolver.Target
	authority    string
	dopts        dialOptions
	csMgr        *connectivityStateManager

	balancerBuildOpts balancer.BuildOptions
	resolverWrapper   *ccResolverWrapper
	blockingpicker    *pickerWrapper

	mu    sync.RWMutex
	sc    ServiceConfig
	scRaw string
	conns map[*addrConn]struct{}
	
	mkp             keepalive.ClientParameters
	curBalancerName string
	preBalancerName string 
	curAddresses    []resolver.Address
	balancerWrapper *ccBalancerWrapper
	retryThrottler  atomic.Value

	channelzID int64 
	czData     *channelzData
}




func (cc *ClientConn) WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool {
	ch := cc.csMgr.getNotifyChan()
	if cc.csMgr.getState() != sourceState {
		return true
	}
	select {
	case <-ctx.Done():
		return false
	case <-ch:
		return true
	}
}



func (cc *ClientConn) GetState() connectivity.State {
	return cc.csMgr.getState()
}

func (cc *ClientConn) scWatcher() {
	for {
		select {
		case sc, ok := <-cc.dopts.scChan:
			if !ok {
				return
			}
			cc.mu.Lock()
			
			
			cc.sc = sc
			cc.scRaw = ""
			cc.mu.Unlock()
		case <-cc.ctx.Done():
			return
		}
	}
}

func (cc *ClientConn) handleResolvedAddrs(addrs []resolver.Address, err error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if cc.conns == nil {
		
		return
	}

	if reflect.DeepEqual(cc.curAddresses, addrs) {
		return
	}

	cc.curAddresses = addrs

	if cc.dopts.balancerBuilder == nil {
		
		
		var isGRPCLB bool
		for _, a := range addrs {
			if a.Type == resolver.GRPCLB {
				isGRPCLB = true
				break
			}
		}
		var newBalancerName string
		if isGRPCLB {
			newBalancerName = grpclbName
		} else {
			
			
			newBalancerName = cc.curBalancerName
			
			if newBalancerName == grpclbName {
				newBalancerName = cc.preBalancerName
			}
			
			
			
			
			
			if newBalancerName == "" {
				newBalancerName = PickFirstBalancerName
			}
		}
		cc.switchBalancer(newBalancerName)
	} else if cc.balancerWrapper == nil {
		
		
		cc.balancerWrapper = newCCBalancerWrapper(cc, cc.dopts.balancerBuilder, cc.balancerBuildOpts)
	}

	cc.balancerWrapper.handleResolvedAddrs(addrs, nil)
}









func (cc *ClientConn) switchBalancer(name string) {
	if cc.conns == nil {
		return
	}

	if strings.ToLower(cc.curBalancerName) == strings.ToLower(name) {
		return
	}

	grpclog.Infof("ClientConn switching balancer to %q", name)
	if cc.dopts.balancerBuilder != nil {
		grpclog.Infoln("ignoring balancer switching: Balancer DialOption used instead")
		return
	}
	
	
	if cc.balancerWrapper != nil {
		cc.balancerWrapper.close()
	}

	builder := balancer.Get(name)
	if builder == nil {
		grpclog.Infof("failed to get balancer builder for: %v, using pick_first instead", name)
		builder = newPickfirstBuilder()
	}
	cc.preBalancerName = cc.curBalancerName
	cc.curBalancerName = builder.Name()
	cc.balancerWrapper = newCCBalancerWrapper(cc, builder, cc.balancerBuildOpts)
}

func (cc *ClientConn) handleSubConnStateChange(sc balancer.SubConn, s connectivity.State) {
	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return
	}
	
	
	cc.balancerWrapper.handleSubConnStateChange(sc, s)
	cc.mu.Unlock()
}




func (cc *ClientConn) newAddrConn(addrs []resolver.Address) (*addrConn, error) {
	ac := &addrConn{
		cc:           cc,
		addrs:        addrs,
		dopts:        cc.dopts,
		czData:       new(channelzData),
		resetBackoff: make(chan struct{}),
	}
	ac.ctx, ac.cancel = context.WithCancel(cc.ctx)
	
	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return nil, ErrClientConnClosing
	}
	if channelz.IsOn() {
		ac.channelzID = channelz.RegisterSubChannel(ac, cc.channelzID, "")
	}
	cc.conns[ac] = struct{}{}
	cc.mu.Unlock()
	return ac, nil
}



func (cc *ClientConn) removeAddrConn(ac *addrConn, err error) {
	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return
	}
	delete(cc.conns, ac)
	cc.mu.Unlock()
	ac.tearDown(err)
}

func (cc *ClientConn) channelzMetric() *channelz.ChannelInternalMetric {
	return &channelz.ChannelInternalMetric{
		State:                    cc.GetState(),
		Target:                   cc.target,
		CallsStarted:             atomic.LoadInt64(&cc.czData.callsStarted),
		CallsSucceeded:           atomic.LoadInt64(&cc.czData.callsSucceeded),
		CallsFailed:              atomic.LoadInt64(&cc.czData.callsFailed),
		LastCallStartedTimestamp: time.Unix(0, atomic.LoadInt64(&cc.czData.lastCallStartedTime)),
	}
}



func (cc *ClientConn) Target() string {
	return cc.target
}

func (cc *ClientConn) incrCallsStarted() {
	atomic.AddInt64(&cc.czData.callsStarted, 1)
	atomic.StoreInt64(&cc.czData.lastCallStartedTime, time.Now().UnixNano())
}

func (cc *ClientConn) incrCallsSucceeded() {
	atomic.AddInt64(&cc.czData.callsSucceeded, 1)
}

func (cc *ClientConn) incrCallsFailed() {
	atomic.AddInt64(&cc.czData.callsFailed, 1)
}






func (ac *addrConn) connect() error {
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		return errConnClosing
	}
	if ac.state != connectivity.Idle {
		ac.mu.Unlock()
		return nil
	}
	ac.state = connectivity.Connecting
	ac.cc.handleSubConnStateChange(ac.acbw, ac.state)
	ac.mu.Unlock()

	
	go func() {
		if err := ac.resetTransport(); err != nil {
			grpclog.Warningf("Failed to dial %s: %v; please retry.", ac.addrs[0].Addr, err)
			if err != errConnClosing {
				
				ac.tearDown(err)
			}
			return
		}
		ac.transportMonitor()
	}()
	return nil
}







func (ac *addrConn) tryUpdateAddrs(addrs []resolver.Address) bool {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	grpclog.Infof("addrConn: tryUpdateAddrs curAddr: %v, addrs: %v", ac.curAddr, addrs)
	if ac.state == connectivity.Shutdown {
		ac.addrs = addrs
		return true
	}

	var curAddrFound bool
	for _, a := range addrs {
		if reflect.DeepEqual(ac.curAddr, a) {
			curAddrFound = true
			break
		}
	}
	grpclog.Infof("addrConn: tryUpdateAddrs curAddrFound: %v", curAddrFound)
	if curAddrFound {
		ac.addrs = addrs
		ac.reconnectIdx = 0 
	}

	return curAddrFound
}








func (cc *ClientConn) GetMethodConfig(method string) MethodConfig {
	
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	m, ok := cc.sc.Methods[method]
	if !ok {
		i := strings.LastIndex(method, "/")
		m = cc.sc.Methods[method[:i+1]]
	}
	return m
}

func (cc *ClientConn) getTransport(ctx context.Context, failfast bool, method string) (transport.ClientTransport, func(balancer.DoneInfo), error) {
	t, done, err := cc.blockingpicker.pick(ctx, failfast, balancer.PickOptions{
		FullMethodName: method,
	})
	if err != nil {
		return nil, nil, toRPCErr(err)
	}
	return t, done, nil
}



func (cc *ClientConn) handleServiceConfig(js string) error {
	if cc.dopts.disableServiceConfig {
		return nil
	}
	sc, err := parseServiceConfig(js)
	if err != nil {
		return err
	}
	cc.mu.Lock()
	cc.scRaw = js
	cc.sc = sc

	if sc.retryThrottling != nil {
		newThrottler := &retryThrottler{
			tokens: sc.retryThrottling.MaxTokens,
			max:    sc.retryThrottling.MaxTokens,
			thresh: sc.retryThrottling.MaxTokens / 2,
			ratio:  sc.retryThrottling.TokenRatio,
		}
		cc.retryThrottler.Store(newThrottler)
	} else {
		cc.retryThrottler.Store((*retryThrottler)(nil))
	}

	if sc.LB != nil && *sc.LB != grpclbName { 
		if cc.curBalancerName == grpclbName {
			
			
			
			
			
			cc.preBalancerName = *sc.LB
		} else {
			cc.switchBalancer(*sc.LB)
			cc.balancerWrapper.handleResolvedAddrs(cc.curAddresses, nil)
		}
	}

	cc.mu.Unlock()
	return nil
}

func (cc *ClientConn) resolveNow(o resolver.ResolveNowOption) {
	cc.mu.RLock()
	r := cc.resolverWrapper
	cc.mu.RUnlock()
	if r == nil {
		return
	}
	go r.resolveNow(o)
}











func (cc *ClientConn) ResetConnectBackoff() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for ac := range cc.conns {
		ac.resetConnectBackoff()
	}
}


func (cc *ClientConn) Close() error {
	defer cc.cancel()

	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return ErrClientConnClosing
	}
	conns := cc.conns
	cc.conns = nil
	cc.csMgr.updateState(connectivity.Shutdown)

	rWrapper := cc.resolverWrapper
	cc.resolverWrapper = nil
	bWrapper := cc.balancerWrapper
	cc.balancerWrapper = nil
	cc.mu.Unlock()

	cc.blockingpicker.close()

	if rWrapper != nil {
		rWrapper.close()
	}
	if bWrapper != nil {
		bWrapper.close()
	}

	for ac := range conns {
		ac.tearDown(ErrClientConnClosing)
	}
	if channelz.IsOn() {
		channelz.RemoveEntry(cc.channelzID)
	}
	return nil
}


type addrConn struct {
	ctx    context.Context
	cancel context.CancelFunc

	cc     *ClientConn
	addrs  []resolver.Address
	dopts  dialOptions
	events trace.EventLog
	acbw   balancer.SubConn

	mu           sync.Mutex
	curAddr      resolver.Address
	reconnectIdx int 
	state        connectivity.State
	
	
	ready     chan struct{}
	transport transport.ClientTransport

	
	tearDownErr error

	connectRetryNum int
	
	
	backoffDeadline time.Time
	
	
	connectDeadline time.Time

	resetBackoff chan struct{}

	channelzID int64 
	czData     *channelzData
}



func (ac *addrConn) adjustParams(r transport.GoAwayReason) {
	switch r {
	case transport.GoAwayTooManyPings:
		v := 2 * ac.dopts.copts.KeepaliveParams.Time
		ac.cc.mu.Lock()
		if v > ac.cc.mkp.Time {
			ac.cc.mkp.Time = v
		}
		ac.cc.mu.Unlock()
	}
}



func (ac *addrConn) printf(format string, a ...interface{}) {
	if ac.events != nil {
		ac.events.Printf(format, a...)
	}
}














func (ac *addrConn) resetTransport() error {
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		return errConnClosing
	}
	if ac.ready != nil {
		close(ac.ready)
		ac.ready = nil
	}
	ac.transport = nil
	ridx := ac.reconnectIdx
	ac.mu.Unlock()
	ac.cc.mu.RLock()
	ac.dopts.copts.KeepaliveParams = ac.cc.mkp
	ac.cc.mu.RUnlock()
	var backoffDeadline, connectDeadline time.Time
	var resetBackoff chan struct{}
	for connectRetryNum := 0; ; connectRetryNum++ {
		ac.mu.Lock()
		if ac.backoffDeadline.IsZero() {
			
			
			
			backoffFor := ac.dopts.bs.Backoff(connectRetryNum) 
			resetBackoff = ac.resetBackoff
			
			dialDuration := getMinConnectTimeout()
			if backoffFor > dialDuration {
				
				dialDuration = backoffFor
			}
			start := time.Now()
			backoffDeadline = start.Add(backoffFor)
			connectDeadline = start.Add(dialDuration)
			ridx = 0 
		} else {
			
			connectRetryNum = ac.connectRetryNum
			backoffDeadline = ac.backoffDeadline
			connectDeadline = ac.connectDeadline
			ac.backoffDeadline = time.Time{}
			ac.connectDeadline = time.Time{}
			ac.connectRetryNum = 0
		}
		if ac.state == connectivity.Shutdown {
			ac.mu.Unlock()
			return errConnClosing
		}
		ac.printf("connecting")
		if ac.state != connectivity.Connecting {
			ac.state = connectivity.Connecting
			ac.cc.handleSubConnStateChange(ac.acbw, ac.state)
		}
		
		addrsIter := make([]resolver.Address, len(ac.addrs))
		copy(addrsIter, ac.addrs)
		copts := ac.dopts.copts
		ac.mu.Unlock()
		connected, err := ac.createTransport(connectRetryNum, ridx, backoffDeadline, connectDeadline, addrsIter, copts, resetBackoff)
		if err != nil {
			return err
		}
		if connected {
			return nil
		}
	}
}



func (ac *addrConn) createTransport(connectRetryNum, ridx int, backoffDeadline, connectDeadline time.Time, addrs []resolver.Address, copts transport.ConnectOptions, resetBackoff chan struct{}) (bool, error) {
	for i := ridx; i < len(addrs); i++ {
		addr := addrs[i]
		target := transport.TargetInfo{
			Addr:      addr.Addr,
			Metadata:  addr.Metadata,
			Authority: ac.cc.authority,
		}
		done := make(chan struct{})
		onPrefaceReceipt := func() {
			ac.mu.Lock()
			close(done)
			if !ac.backoffDeadline.IsZero() {
				
				
				
				
				
				ac.backoffDeadline = time.Time{}
				ac.connectDeadline = time.Time{}
				ac.connectRetryNum = 0
			}
			ac.mu.Unlock()
		}
		
		
		connectCtx, cancel := context.WithDeadline(ac.ctx, connectDeadline)
		if channelz.IsOn() {
			copts.ChannelzParentID = ac.channelzID
		}
		newTr, err := transport.NewClientTransport(connectCtx, ac.cc.ctx, target, copts, onPrefaceReceipt)
		if err != nil {
			cancel()
			ac.cc.blockingpicker.updateConnectionError(err)
			ac.mu.Lock()
			if ac.state == connectivity.Shutdown {
				
				ac.mu.Unlock()
				return false, errConnClosing
			}
			ac.mu.Unlock()
			grpclog.Warningf("grpc: addrConn.createTransport failed to connect to %v. Err :%v. Reconnecting...", addr, err)
			continue
		}
		if ac.dopts.waitForHandshake {
			select {
			case <-done:
			case <-connectCtx.Done():
				
				grpclog.Warningf("grpc: addrConn.createTransport failed to receive server preface before deadline.")
				newTr.Close()
				continue
			case <-ac.ctx.Done():
			}
		}
		ac.mu.Lock()
		if ac.state == connectivity.Shutdown {
			ac.mu.Unlock()
			
			newTr.Close()
			return false, errConnClosing
		}
		ac.printf("ready")
		ac.state = connectivity.Ready
		ac.cc.handleSubConnStateChange(ac.acbw, ac.state)
		ac.transport = newTr
		ac.curAddr = addr
		if ac.ready != nil {
			close(ac.ready)
			ac.ready = nil
		}
		select {
		case <-done:
			
			
		default:
			ac.connectRetryNum = connectRetryNum
			ac.backoffDeadline = backoffDeadline
			ac.connectDeadline = connectDeadline
			ac.reconnectIdx = i + 1 
		}
		ac.mu.Unlock()
		return true, nil
	}
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		return false, errConnClosing
	}
	ac.state = connectivity.TransientFailure
	ac.cc.handleSubConnStateChange(ac.acbw, ac.state)
	ac.cc.resolveNow(resolver.ResolveNowOption{})
	if ac.ready != nil {
		close(ac.ready)
		ac.ready = nil
	}
	ac.mu.Unlock()
	timer := time.NewTimer(backoffDeadline.Sub(time.Now()))
	select {
	case <-timer.C:
	case <-resetBackoff:
		timer.Stop()
	case <-ac.ctx.Done():
		timer.Stop()
		return false, ac.ctx.Err()
	}
	return false, nil
}

func (ac *addrConn) resetConnectBackoff() {
	ac.mu.Lock()
	close(ac.resetBackoff)
	ac.resetBackoff = make(chan struct{})
	ac.connectRetryNum = 0
	ac.mu.Unlock()
}



func (ac *addrConn) transportMonitor() {
	for {
		var timer *time.Timer
		var cdeadline <-chan time.Time
		ac.mu.Lock()
		t := ac.transport
		if !ac.connectDeadline.IsZero() {
			timer = time.NewTimer(ac.connectDeadline.Sub(time.Now()))
			cdeadline = timer.C
		}
		ac.mu.Unlock()
		
		select {
		case <-t.GoAway():
			done := t.Error()
			cleanup := t.Close
			
			
			
			go func() {
				<-done
				cleanup()
			}()
		case <-t.Error():
			
			
			
			t.Close()
		case <-cdeadline:
			ac.mu.Lock()
			
			if ac.backoffDeadline.IsZero() {
				ac.mu.Unlock()
				continue
			}
			ac.mu.Unlock()
			timer = nil
			
			
			grpclog.Warningf("grpc: addrConn.transportMonitor didn't get server preface after waiting. Closing the new transport now.")
			t.Close()
		}
		if timer != nil {
			timer.Stop()
		}
		
		
		select {
		case <-t.GoAway():
			ac.adjustParams(t.GetGoAwayReason())
		default:
		}
		ac.mu.Lock()
		if ac.state == connectivity.Shutdown {
			ac.mu.Unlock()
			return
		}
		
		
		ac.state = connectivity.TransientFailure
		ac.cc.handleSubConnStateChange(ac.acbw, ac.state)
		ac.cc.resolveNow(resolver.ResolveNowOption{})
		ac.curAddr = resolver.Address{}
		ac.mu.Unlock()
		if err := ac.resetTransport(); err != nil {
			ac.mu.Lock()
			ac.printf("transport exiting: %v", err)
			ac.mu.Unlock()
			grpclog.Warningf("grpc: addrConn.transportMonitor exits due to: %v", err)
			if err != errConnClosing {
				
				ac.tearDown(err)
			}
			return
		}
	}
}




func (ac *addrConn) getReadyTransport() (transport.ClientTransport, bool) {
	ac.mu.Lock()
	if ac.state == connectivity.Ready {
		t := ac.transport
		ac.mu.Unlock()
		return t, true
	}
	var idle bool
	if ac.state == connectivity.Idle {
		idle = true
	}
	ac.mu.Unlock()
	
	if idle {
		ac.connect()
	}
	return nil, false
}






func (ac *addrConn) tearDown(err error) {
	ac.cancel()
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if ac.state == connectivity.Shutdown {
		return
	}
	ac.curAddr = resolver.Address{}
	if err == errConnDrain && ac.transport != nil {
		
		
		
		
		ac.transport.GracefulClose()
	}
	ac.state = connectivity.Shutdown
	ac.tearDownErr = err
	ac.cc.handleSubConnStateChange(ac.acbw, ac.state)
	if ac.events != nil {
		ac.events.Finish()
		ac.events = nil
	}
	if ac.ready != nil {
		close(ac.ready)
		ac.ready = nil
	}
	if channelz.IsOn() {
		channelz.RemoveEntry(ac.channelzID)
	}
}

func (ac *addrConn) getState() connectivity.State {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.state
}

func (ac *addrConn) ChannelzMetric() *channelz.ChannelInternalMetric {
	ac.mu.Lock()
	addr := ac.curAddr.Addr
	ac.mu.Unlock()
	return &channelz.ChannelInternalMetric{
		State:                    ac.getState(),
		Target:                   addr,
		CallsStarted:             atomic.LoadInt64(&ac.czData.callsStarted),
		CallsSucceeded:           atomic.LoadInt64(&ac.czData.callsSucceeded),
		CallsFailed:              atomic.LoadInt64(&ac.czData.callsFailed),
		LastCallStartedTimestamp: time.Unix(0, atomic.LoadInt64(&ac.czData.lastCallStartedTime)),
	}
}

func (ac *addrConn) incrCallsStarted() {
	atomic.AddInt64(&ac.czData.callsStarted, 1)
	atomic.StoreInt64(&ac.czData.lastCallStartedTime, time.Now().UnixNano())
}

func (ac *addrConn) incrCallsSucceeded() {
	atomic.AddInt64(&ac.czData.callsSucceeded, 1)
}

func (ac *addrConn) incrCallsFailed() {
	atomic.AddInt64(&ac.czData.callsFailed, 1)
}

type retryThrottler struct {
	max    float64
	thresh float64
	ratio  float64

	mu     sync.Mutex
	tokens float64 
}




func (rt *retryThrottler) throttle() bool {
	if rt == nil {
		return false
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.tokens--
	if rt.tokens < 0 {
		rt.tokens = 0
	}
	return rt.tokens <= rt.thresh
}

func (rt *retryThrottler) successfulRPC() {
	if rt == nil {
		return
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.tokens += rt.ratio
	if rt.tokens > rt.max {
		rt.tokens = rt.max
	}
}

type channelzChannel struct {
	cc *ClientConn
}

func (c *channelzChannel) ChannelzMetric() *channelz.ChannelInternalMetric {
	return c.cc.channelzMetric()
}






var ErrClientConnTimeout = errors.New("grpc: timed out when dialing")
