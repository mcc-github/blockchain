

package grpc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/balancer"
	_ "google.golang.org/grpc/balancer/roundrobin" 
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/envconfig"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
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
	
	
	errTransportCredsAndBundle = errors.New("grpc: credentials.Bundle may not be used with individual TransportCredentials")
	
	
	
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
		target:            target,
		csMgr:             &connectivityStateManager{},
		conns:             make(map[*addrConn]struct{}),
		dopts:             defaultDialOptions(),
		blockingpicker:    newPickerWrapper(),
		czData:            new(channelzData),
		firstResolveEvent: grpcsync.NewEvent(),
	}
	cc.retryThrottler.Store((*retryThrottler)(nil))
	cc.ctx, cc.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt.apply(&cc.dopts)
	}

	if channelz.IsOn() {
		if cc.dopts.channelzParentID != 0 {
			cc.channelzID = channelz.RegisterChannel(&channelzChannel{cc}, cc.dopts.channelzParentID, target)
			channelz.AddTraceEvent(cc.channelzID, &channelz.TraceEventDesc{
				Desc:     "Channel Created",
				Severity: channelz.CtINFO,
				Parent: &channelz.TraceEventDesc{
					Desc:     fmt.Sprintf("Nested Channel(id:%d) created", cc.channelzID),
					Severity: channelz.CtINFO,
				},
			})
		} else {
			cc.channelzID = channelz.RegisterChannel(&channelzChannel{cc}, 0, target)
			channelz.AddTraceEvent(cc.channelzID, &channelz.TraceEventDesc{
				Desc:     "Channel Created",
				Severity: channelz.CtINFO,
			})
		}
		cc.csMgr.channelzID = cc.channelzID
	}

	if !cc.dopts.insecure {
		if cc.dopts.copts.TransportCredentials == nil && cc.dopts.copts.CredsBundle == nil {
			return nil, errNoTransportSecurity
		}
		if cc.dopts.copts.TransportCredentials != nil && cc.dopts.copts.CredsBundle != nil {
			return nil, errTransportCredsAndBundle
		}
	} else {
		if cc.dopts.copts.TransportCredentials != nil || cc.dopts.copts.CredsBundle != nil {
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
				return (&net.Dialer{}).DialContext(ctx, network, addr)
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
		CredsBundle:      cc.dopts.copts.CredsBundle,
		Dialer:           cc.dopts.copts.Dialer,
		ChannelzParentID: cc.channelzID,
	}

	
	rWrapper, err := newCCResolverWrapper(cc)
	if err != nil {
		return nil, fmt.Errorf("failed to build resolver: %v", err)
	}

	cc.mu.Lock()
	cc.resolverWrapper = rWrapper
	cc.mu.Unlock()
	
	if cc.dopts.block {
		for {
			s := cc.GetState()
			if s == connectivity.Ready {
				break
			} else if cc.dopts.copts.FailOnNonTempDialError && s == connectivity.TransientFailure {
				if err = cc.blockingpicker.connectionError(); err != nil {
					terr, ok := err.(interface {
						Temporary() bool
					})
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
	channelzID int64
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
	if channelz.IsOn() {
		channelz.AddTraceEvent(csm.channelzID, &channelz.TraceEventDesc{
			Desc:     fmt.Sprintf("Channel Connectivity change to %v", state),
			Severity: channelz.CtINFO,
		})
	}
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
	blockingpicker    *pickerWrapper

	mu              sync.RWMutex
	resolverWrapper *ccResolverWrapper
	sc              ServiceConfig
	scRaw           string
	conns           map[*addrConn]struct{}
	
	mkp             keepalive.ClientParameters
	curBalancerName string
	preBalancerName string 
	curAddresses    []resolver.Address
	balancerWrapper *ccBalancerWrapper
	retryThrottler  atomic.Value

	firstResolveEvent *grpcsync.Event

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




func (cc *ClientConn) waitForResolvedAddrs(ctx context.Context) error {
	
	
	if cc.firstResolveEvent.HasFired() {
		return nil
	}
	select {
	case <-cc.firstResolveEvent.Done():
		return nil
	case <-ctx.Done():
		return status.FromContextError(ctx.Err()).Err()
	case <-cc.ctx.Done():
		return ErrClientConnClosing
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
	cc.firstResolveEvent.Fire()

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
	
	
	if channelz.IsOn() {
		if builder == nil {
			channelz.AddTraceEvent(cc.channelzID, &channelz.TraceEventDesc{
				Desc:     fmt.Sprintf("Channel switches to new LB policy %q due to fallback from invalid balancer name", PickFirstBalancerName),
				Severity: channelz.CtWarning,
			})
		} else {
			channelz.AddTraceEvent(cc.channelzID, &channelz.TraceEventDesc{
				Desc:     fmt.Sprintf("Channel switches to new LB policy %q", name),
				Severity: channelz.CtINFO,
			})
		}
	}
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




func (cc *ClientConn) newAddrConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (*addrConn, error) {
	ac := &addrConn{
		cc:           cc,
		addrs:        addrs,
		scopts:       opts,
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
		channelz.AddTraceEvent(ac.channelzID, &channelz.TraceEventDesc{
			Desc:     "Subchannel Created",
			Severity: channelz.CtINFO,
			Parent: &channelz.TraceEventDesc{
				Desc:     fmt.Sprintf("Subchannel(id:%d) created", ac.channelzID),
				Severity: channelz.CtINFO,
			},
		})
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
	ac.updateConnectivityState(connectivity.Connecting)
	ac.mu.Unlock()

	
	go ac.resetTransport()
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

	
	
	if ac.state != connectivity.Ready {
		return false
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

func (cc *ClientConn) healthCheckConfig() *healthCheckConfig {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.sc.healthCheckConfig
}

func (cc *ClientConn) getTransport(ctx context.Context, failfast bool, method string) (transport.ClientTransport, func(balancer.DoneInfo), error) {
	hdr, _ := metadata.FromOutgoingContext(ctx)
	t, done, err := cc.blockingpicker.pick(ctx, failfast, balancer.PickOptions{
		FullMethodName: method,
		Header:         hdr,
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
	if cc.scRaw == js {
		return nil
	}
	if channelz.IsOn() {
		channelz.AddTraceEvent(cc.channelzID, &channelz.TraceEventDesc{
			
			
			Desc:     fmt.Sprintf("Channel has a new service config \"%s\"", js),
			Severity: channelz.CtINFO,
		})
	}
	sc, err := parseServiceConfig(js)
	if err != nil {
		return err
	}
	cc.mu.Lock()
	
	
	
	if cc.conns == nil {
		cc.mu.Unlock()
		return nil
	}
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
		ted := &channelz.TraceEventDesc{
			Desc:     "Channel Deleted",
			Severity: channelz.CtINFO,
		}
		if cc.dopts.channelzParentID != 0 {
			ted.Parent = &channelz.TraceEventDesc{
				Desc:     fmt.Sprintf("Nested channel(id:%d) deleted", cc.channelzID),
				Severity: channelz.CtINFO,
			}
		}
		channelz.AddTraceEvent(cc.channelzID, ted)
		
		
		channelz.RemoveEntry(cc.channelzID)
	}
	return nil
}


type addrConn struct {
	ctx    context.Context
	cancel context.CancelFunc

	cc     *ClientConn
	dopts  dialOptions
	acbw   balancer.SubConn
	scopts balancer.NewSubConnOptions

	
	
	
	
	transport transport.ClientTransport 

	mu      sync.Mutex
	curAddr resolver.Address   
	addrs   []resolver.Address 

	
	state connectivity.State

	tearDownErr error 

	backoffIdx   int 
	resetBackoff chan struct{}

	channelzID int64 
	czData     *channelzData
}


func (ac *addrConn) updateConnectivityState(s connectivity.State) {
	if ac.state == s {
		return
	}

	updateMsg := fmt.Sprintf("Subchannel Connectivity change to %v", s)
	ac.state = s
	if channelz.IsOn() {
		channelz.AddTraceEvent(ac.channelzID, &channelz.TraceEventDesc{
			Desc:     updateMsg,
			Severity: channelz.CtINFO,
		})
	}
	ac.cc.handleSubConnStateChange(ac.acbw, s)
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

func (ac *addrConn) resetTransport() {
	for i := 0; ; i++ {
		tryNextAddrFromStart := grpcsync.NewEvent()

		ac.mu.Lock()
		if i > 0 {
			ac.cc.resolveNow(resolver.ResolveNowOption{})
		}
		addrs := ac.addrs
		backoffFor := ac.dopts.bs.Backoff(ac.backoffIdx)

		
		dialDuration := getMinConnectTimeout()
		if dialDuration < backoffFor {
			
			dialDuration = backoffFor
		}
		connectDeadline := time.Now().Add(dialDuration)
		ac.mu.Unlock()

	addrLoop:
		for _, addr := range addrs {
			ac.mu.Lock()

			if ac.state == connectivity.Shutdown {
				ac.mu.Unlock()
				return
			}
			ac.updateConnectivityState(connectivity.Connecting)
			ac.transport = nil

			ac.cc.mu.RLock()
			ac.dopts.copts.KeepaliveParams = ac.cc.mkp
			ac.cc.mu.RUnlock()

			if ac.state == connectivity.Shutdown {
				ac.mu.Unlock()
				return
			}

			copts := ac.dopts.copts
			if ac.scopts.CredsBundle != nil {
				copts.CredsBundle = ac.scopts.CredsBundle
			}
			hctx, hcancel := context.WithCancel(ac.ctx)
			defer hcancel()
			ac.mu.Unlock()

			if channelz.IsOn() {
				channelz.AddTraceEvent(ac.channelzID, &channelz.TraceEventDesc{
					Desc:     fmt.Sprintf("Subchannel picks a new address %q to connect", addr.Addr),
					Severity: channelz.CtINFO,
				})
			}

			reconnect := grpcsync.NewEvent()
			prefaceReceived := make(chan struct{})
			newTr, err := ac.createTransport(addr, copts, connectDeadline, reconnect, prefaceReceived)
			if err == nil {
				ac.mu.Lock()
				ac.curAddr = addr
				ac.transport = newTr
				ac.mu.Unlock()

				healthCheckConfig := ac.cc.healthCheckConfig()
				
				
				
				
				
				healthcheckManagingState := false
				if !ac.cc.dopts.disableHealthCheck && healthCheckConfig != nil && ac.scopts.HealthCheckEnabled {
					if ac.cc.dopts.healthCheckFunc == nil {
						
						grpclog.Error("the client side LB channel health check function has not been set.")
					} else {
						
						go ac.startHealthCheck(hctx, newTr, addr, healthCheckConfig.ServiceName)
						healthcheckManagingState = true
					}
				}
				if !healthcheckManagingState {
					ac.mu.Lock()
					ac.updateConnectivityState(connectivity.Ready)
					ac.mu.Unlock()
				}
			} else {
				hcancel()
				if err == errConnClosing {
					return
				}

				if tryNextAddrFromStart.HasFired() {
					break addrLoop
				}
				continue
			}

			ac.mu.Lock()
			reqHandshake := ac.dopts.reqHandshake
			ac.mu.Unlock()

			<-reconnect.Done()
			hcancel()

			if reqHandshake == envconfig.RequireHandshakeHybrid {
				
				
				
				
				
				select {
				case <-prefaceReceived:
					
					
					ac.mu.Lock()
					ac.backoffIdx = 0
					ac.mu.Unlock()
					break addrLoop
				default:
					
					
					
					ac.mu.Lock()
					if ac.state == connectivity.Shutdown {
						ac.mu.Unlock()
						return
					}

					ac.updateConnectivityState(connectivity.TransientFailure)
					ac.mu.Unlock()

					if tryNextAddrFromStart.HasFired() {
						break addrLoop
					}
				}
			} else {
				
				
				
				
				
				
				ac.mu.Lock()
				ac.backoffIdx = 0
				ac.mu.Unlock()
				break addrLoop
			}
		}

		
		
		ac.mu.Lock()
		if ac.state == connectivity.Shutdown {
			ac.mu.Unlock()
			return
		}
		ac.updateConnectivityState(connectivity.TransientFailure)

		
		b := ac.resetBackoff
		timer := time.NewTimer(backoffFor)
		acctx := ac.ctx
		ac.mu.Unlock()

		select {
		case <-timer.C:
			ac.mu.Lock()
			ac.backoffIdx++
			ac.mu.Unlock()
		case <-b:
			timer.Stop()
		case <-acctx.Done():
			timer.Stop()
			return
		}
	}
}






func (ac *addrConn) createTransport(addr resolver.Address, copts transport.ConnectOptions, connectDeadline time.Time, reconnect *grpcsync.Event, prefaceReceived chan struct{}) (transport.ClientTransport, error) {
	onCloseCalled := make(chan struct{})

	target := transport.TargetInfo{
		Addr:      addr.Addr,
		Metadata:  addr.Metadata,
		Authority: ac.cc.authority,
	}

	prefaceTimer := time.NewTimer(time.Until(connectDeadline))

	onGoAway := func(r transport.GoAwayReason) {
		ac.mu.Lock()
		ac.adjustParams(r)
		ac.mu.Unlock()
		reconnect.Fire()
	}

	onClose := func() {
		close(onCloseCalled)
		prefaceTimer.Stop()
		reconnect.Fire()
	}

	onPrefaceReceipt := func() {
		close(prefaceReceived)
		prefaceTimer.Stop()
	}

	connectCtx, cancel := context.WithDeadline(ac.ctx, connectDeadline)
	defer cancel()
	if channelz.IsOn() {
		copts.ChannelzParentID = ac.channelzID
	}

	newTr, err := transport.NewClientTransport(connectCtx, ac.cc.ctx, target, copts, onPrefaceReceipt, onGoAway, onClose)

	if err == nil {
		if ac.dopts.reqHandshake == envconfig.RequireHandshakeOn {
			select {
			case <-prefaceTimer.C:
				
				newTr.Close()
				err = errors.New("timed out waiting for server handshake")
			case <-prefaceReceived:
				
			case <-onCloseCalled:
				
				return nil, errors.New("connection closed")
			}
		} else if ac.dopts.reqHandshake == envconfig.RequireHandshakeHybrid {
			go func() {
				select {
				case <-prefaceTimer.C:
					
					newTr.Close()
				case <-prefaceReceived:
					
				case <-onCloseCalled:
					
				}
			}()
		}
	}

	if err != nil {
		
		ac.cc.blockingpicker.updateConnectionError(err)
		ac.mu.Lock()
		if ac.state == connectivity.Shutdown {
			
			ac.mu.Unlock()

			return nil, errConnClosing
		}
		ac.mu.Unlock()
		grpclog.Warningf("grpc: addrConn.createTransport failed to connect to %v. Err :%v. Reconnecting...", addr, err)
		return nil, err
	}

	
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		newTr.Close()
		return nil, errConnClosing
	}
	ac.mu.Unlock()

	
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		newTr.Close()
		return nil, errConnClosing
	}
	ac.mu.Unlock()

	return newTr, nil
}

func (ac *addrConn) startHealthCheck(ctx context.Context, newTr transport.ClientTransport, addr resolver.Address, serviceName string) {
	
	newStream := func() (interface{}, error) {
		return ac.newClientStream(ctx, &StreamDesc{ServerStreams: true}, "/grpc.health.v1.Health/Watch", newTr)
	}
	firstReady := true
	reportHealth := func(ok bool) {
		ac.mu.Lock()
		defer ac.mu.Unlock()
		if ac.transport != newTr {
			return
		}
		if ok {
			if firstReady {
				firstReady = false
				ac.curAddr = addr
			}
			ac.updateConnectivityState(connectivity.Ready)
		} else {
			ac.updateConnectivityState(connectivity.TransientFailure)
		}
	}
	err := ac.cc.dopts.healthCheckFunc(ctx, newStream, reportHealth, serviceName)
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			if channelz.IsOn() {
				channelz.AddTraceEvent(ac.channelzID, &channelz.TraceEventDesc{
					Desc:     "Subchannel health check is unimplemented at server side, thus health check is disabled",
					Severity: channelz.CtError,
				})
			}
			grpclog.Error("Subchannel health check is unimplemented at server side, thus health check is disabled")
		} else {
			grpclog.Errorf("HealthCheckFunc exits with unexpected error %v", err)
		}
	}
}

func (ac *addrConn) resetConnectBackoff() {
	ac.mu.Lock()
	close(ac.resetBackoff)
	ac.backoffIdx = 0
	ac.resetBackoff = make(chan struct{})
	ac.mu.Unlock()
}




func (ac *addrConn) getReadyTransport() (transport.ClientTransport, bool) {
	ac.mu.Lock()
	if ac.state == connectivity.Ready && ac.transport != nil {
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
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		return
	}
	curTr := ac.transport
	ac.transport = nil
	
	
	ac.updateConnectivityState(connectivity.Shutdown)
	ac.cancel()
	ac.tearDownErr = err
	ac.curAddr = resolver.Address{}
	if err == errConnDrain && curTr != nil {
		
		
		
		
		
		ac.mu.Unlock()
		curTr.GracefulClose()
		ac.mu.Lock()
	}
	if channelz.IsOn() {
		channelz.AddTraceEvent(ac.channelzID, &channelz.TraceEventDesc{
			Desc:     "Subchannel Deleted",
			Severity: channelz.CtINFO,
			Parent: &channelz.TraceEventDesc{
				Desc:     fmt.Sprintf("Subchanel(id:%d) deleted", ac.channelzID),
				Severity: channelz.CtINFO,
			},
		})
		
		
		channelz.RemoveEntry(ac.channelzID)
	}
	ac.mu.Unlock()
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
