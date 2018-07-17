

package grpc

import (
	"sync"
	"sync/atomic"

	"golang.org/x/net/context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	lbpb "google.golang.org/grpc/grpclb/grpc_lb_v1/messages"
	"google.golang.org/grpc/status"
)

type rpcStats struct {
	NumCallsStarted                          int64
	NumCallsFinished                         int64
	NumCallsFinishedWithDropForRateLimiting  int64
	NumCallsFinishedWithDropForLoadBalancing int64
	NumCallsFinishedWithClientFailedToSend   int64
	NumCallsFinishedKnownReceived            int64
}


func (s *rpcStats) toClientStats() *lbpb.ClientStats {
	stats := &lbpb.ClientStats{
		NumCallsStarted:                          atomic.SwapInt64(&s.NumCallsStarted, 0),
		NumCallsFinished:                         atomic.SwapInt64(&s.NumCallsFinished, 0),
		NumCallsFinishedWithDropForRateLimiting:  atomic.SwapInt64(&s.NumCallsFinishedWithDropForRateLimiting, 0),
		NumCallsFinishedWithDropForLoadBalancing: atomic.SwapInt64(&s.NumCallsFinishedWithDropForLoadBalancing, 0),
		NumCallsFinishedWithClientFailedToSend:   atomic.SwapInt64(&s.NumCallsFinishedWithClientFailedToSend, 0),
		NumCallsFinishedKnownReceived:            atomic.SwapInt64(&s.NumCallsFinishedKnownReceived, 0),
	}
	return stats
}

func (s *rpcStats) dropForRateLimiting() {
	atomic.AddInt64(&s.NumCallsStarted, 1)
	atomic.AddInt64(&s.NumCallsFinishedWithDropForRateLimiting, 1)
	atomic.AddInt64(&s.NumCallsFinished, 1)
}

func (s *rpcStats) dropForLoadBalancing() {
	atomic.AddInt64(&s.NumCallsStarted, 1)
	atomic.AddInt64(&s.NumCallsFinishedWithDropForLoadBalancing, 1)
	atomic.AddInt64(&s.NumCallsFinished, 1)
}

func (s *rpcStats) failedToSend() {
	atomic.AddInt64(&s.NumCallsStarted, 1)
	atomic.AddInt64(&s.NumCallsFinishedWithClientFailedToSend, 1)
	atomic.AddInt64(&s.NumCallsFinished, 1)
}

func (s *rpcStats) knownReceived() {
	atomic.AddInt64(&s.NumCallsStarted, 1)
	atomic.AddInt64(&s.NumCallsFinishedKnownReceived, 1)
	atomic.AddInt64(&s.NumCallsFinished, 1)
}

type errPicker struct {
	
	err error
}

func (p *errPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	return nil, nil, p.err
}






type rrPicker struct {
	mu           sync.Mutex
	subConns     []balancer.SubConn 
	subConnsNext int
}

func (p *rrPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	sc := p.subConns[p.subConnsNext]
	p.subConnsNext = (p.subConnsNext + 1) % len(p.subConns)
	return sc, nil, nil
}










type lbPicker struct {
	mu             sync.Mutex
	serverList     []*lbpb.Server
	serverListNext int
	subConns       []balancer.SubConn 
	subConnsNext   int

	stats *rpcStats
}

func (p *lbPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	
	s := p.serverList[p.serverListNext]
	p.serverListNext = (p.serverListNext + 1) % len(p.serverList)

	
	if s.DropForRateLimiting {
		p.stats.dropForRateLimiting()
		return nil, nil, status.Errorf(codes.Unavailable, "request dropped by grpclb")
	}
	if s.DropForLoadBalancing {
		p.stats.dropForLoadBalancing()
		return nil, nil, status.Errorf(codes.Unavailable, "request dropped by grpclb")
	}

	
	if len(p.subConns) <= 0 {
		return nil, nil, balancer.ErrNoSubConnAvailable
	}

	
	sc := p.subConns[p.subConnsNext]
	p.subConnsNext = (p.subConnsNext + 1) % len(p.subConns)
	done := func(info balancer.DoneInfo) {
		if !info.BytesSent {
			p.stats.failedToSend()
		} else if info.BytesReceived {
			p.stats.knownReceived()
		}
	}
	return sc, done, nil
}
