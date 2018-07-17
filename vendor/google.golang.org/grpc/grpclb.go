

package grpc

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	lbpb "google.golang.org/grpc/grpclb/grpc_lb_v1/messages"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

const (
	lbTokeyKey             = "lb-token"
	defaultFallbackTimeout = 10 * time.Second
	grpclbName             = "grpclb"
)

func convertDuration(d *lbpb.Duration) time.Duration {
	if d == nil {
		return 0
	}
	return time.Duration(d.Seconds)*time.Second + time.Duration(d.Nanos)*time.Nanosecond
}




type loadBalancerClient struct {
	cc *ClientConn
}

func (c *loadBalancerClient) BalanceLoad(ctx context.Context, opts ...CallOption) (*balanceLoadClientStream, error) {
	desc := &StreamDesc{
		StreamName:    "BalanceLoad",
		ServerStreams: true,
		ClientStreams: true,
	}
	stream, err := NewClientStream(ctx, desc, c.cc, "/grpc.lb.v1.LoadBalancer/BalanceLoad", opts...)
	if err != nil {
		return nil, err
	}
	x := &balanceLoadClientStream{stream}
	return x, nil
}

type balanceLoadClientStream struct {
	ClientStream
}

func (x *balanceLoadClientStream) Send(m *lbpb.LoadBalanceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *balanceLoadClientStream) Recv() (*lbpb.LoadBalanceResponse, error) {
	m := new(lbpb.LoadBalanceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func init() {
	balancer.Register(newLBBuilder())
}


func newLBBuilder() balancer.Builder {
	return NewLBBuilderWithFallbackTimeout(defaultFallbackTimeout)
}







func NewLBBuilderWithFallbackTimeout(fallbackTimeout time.Duration) balancer.Builder {
	return &lbBuilder{
		fallbackTimeout: fallbackTimeout,
	}
}

type lbBuilder struct {
	fallbackTimeout time.Duration
}

func (b *lbBuilder) Name() string {
	return grpclbName
}

func (b *lbBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
	
	
	
	scheme := "grpclb_internal_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	r := &lbManualResolver{scheme: scheme, ccb: cc}

	var target string
	targetSplitted := strings.Split(cc.Target(), ":/
	if len(targetSplitted) < 2 {
		target = cc.Target()
	} else {
		target = targetSplitted[1]
	}

	lb := &lbBalancer{
		cc:              cc,
		target:          target,
		opt:             opt,
		fallbackTimeout: b.fallbackTimeout,
		doneCh:          make(chan struct{}),

		manualResolver: r,
		csEvltr:        &connectivityStateEvaluator{},
		subConns:       make(map[resolver.Address]balancer.SubConn),
		scStates:       make(map[balancer.SubConn]connectivity.State),
		picker:         &errPicker{err: balancer.ErrNoSubConnAvailable},
		clientStats:    &rpcStats{},
	}

	return lb
}

type lbBalancer struct {
	cc              balancer.ClientConn
	target          string
	opt             balancer.BuildOptions
	fallbackTimeout time.Duration
	doneCh          chan struct{}

	
	
	
	manualResolver *lbManualResolver
	
	ccRemoteLB *ClientConn

	
	
	clientStats *rpcStats

	mu sync.Mutex 
	
	
	
	fullServerList []*lbpb.Server
	
	
	
	
	backendAddrs []resolver.Address
	
	csEvltr  *connectivityStateEvaluator
	state    connectivity.State
	subConns map[resolver.Address]balancer.SubConn   
	scStates map[balancer.SubConn]connectivity.State 
	picker   balancer.Picker
	
	
	fallbackTimerExpired bool
	serverListReceived   bool
	
	
	
	resolvedBackendAddrs []resolver.Address
}






func (lb *lbBalancer) regeneratePicker() {
	if lb.state == connectivity.TransientFailure {
		lb.picker = &errPicker{err: balancer.ErrTransientFailure}
		return
	}
	var readySCs []balancer.SubConn
	for _, a := range lb.backendAddrs {
		if sc, ok := lb.subConns[a]; ok {
			if st, ok := lb.scStates[sc]; ok && st == connectivity.Ready {
				readySCs = append(readySCs, sc)
			}
		}
	}

	if len(lb.fullServerList) <= 0 {
		if len(readySCs) <= 0 {
			lb.picker = &errPicker{err: balancer.ErrNoSubConnAvailable}
			return
		}
		lb.picker = &rrPicker{subConns: readySCs}
		return
	}
	lb.picker = &lbPicker{
		serverList: lb.fullServerList,
		subConns:   readySCs,
		stats:      lb.clientStats,
	}
	return
}

func (lb *lbBalancer) HandleSubConnStateChange(sc balancer.SubConn, s connectivity.State) {
	grpclog.Infof("lbBalancer: handle SubConn state change: %p, %v", sc, s)
	lb.mu.Lock()
	defer lb.mu.Unlock()

	oldS, ok := lb.scStates[sc]
	if !ok {
		grpclog.Infof("lbBalancer: got state changes for an unknown SubConn: %p, %v", sc, s)
		return
	}
	lb.scStates[sc] = s
	switch s {
	case connectivity.Idle:
		sc.Connect()
	case connectivity.Shutdown:
		
		
		delete(lb.scStates, sc)
	}

	oldAggrState := lb.state
	lb.state = lb.csEvltr.recordTransition(oldS, s)

	
	
	
	
	
	if (oldS == connectivity.Ready) != (s == connectivity.Ready) ||
		(lb.state == connectivity.TransientFailure) != (oldAggrState == connectivity.TransientFailure) {
		lb.regeneratePicker()
	}

	lb.cc.UpdateBalancerState(lb.state, lb.picker)
	return
}




func (lb *lbBalancer) fallbackToBackendsAfter(fallbackTimeout time.Duration) {
	timer := time.NewTimer(fallbackTimeout)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-lb.doneCh:
		return
	}
	lb.mu.Lock()
	if lb.serverListReceived {
		lb.mu.Unlock()
		return
	}
	lb.fallbackTimerExpired = true
	lb.refreshSubConns(lb.resolvedBackendAddrs)
	lb.mu.Unlock()
}




func (lb *lbBalancer) HandleResolvedAddrs(addrs []resolver.Address, err error) {
	grpclog.Infof("lbBalancer: handleResolvedResult: %+v", addrs)
	if len(addrs) <= 0 {
		return
	}

	var remoteBalancerAddrs, backendAddrs []resolver.Address
	for _, a := range addrs {
		if a.Type == resolver.GRPCLB {
			remoteBalancerAddrs = append(remoteBalancerAddrs, a)
		} else {
			backendAddrs = append(backendAddrs, a)
		}
	}

	if lb.ccRemoteLB == nil {
		if len(remoteBalancerAddrs) <= 0 {
			grpclog.Errorf("grpclb: no remote balancer address is available, should never happen")
			return
		}
		
		
		lb.dialRemoteLB(remoteBalancerAddrs[0].ServerName)
		
		go lb.fallbackToBackendsAfter(lb.fallbackTimeout)
	}

	
	
	lb.manualResolver.NewAddress(remoteBalancerAddrs)

	lb.mu.Lock()
	lb.resolvedBackendAddrs = backendAddrs
	
	
	
	if !lb.serverListReceived && lb.fallbackTimerExpired {
		
		
		
		lb.refreshSubConns(lb.resolvedBackendAddrs)
	}
	lb.mu.Unlock()
}

func (lb *lbBalancer) Close() {
	select {
	case <-lb.doneCh:
		return
	default:
	}
	close(lb.doneCh)
	if lb.ccRemoteLB != nil {
		lb.ccRemoteLB.Close()
	}
}
