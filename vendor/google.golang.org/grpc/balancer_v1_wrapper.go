

package grpc

import (
	"context"
	"strings"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

type balancerWrapperBuilder struct {
	b Balancer 
}

func (bwb *balancerWrapperBuilder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	targetAddr := cc.Target()
	targetSplitted := strings.Split(targetAddr, ":/
	if len(targetSplitted) >= 2 {
		targetAddr = targetSplitted[1]
	}

	bwb.b.Start(targetAddr, BalancerConfig{
		DialCreds: opts.DialCreds,
		Dialer:    opts.Dialer,
	})
	_, pickfirst := bwb.b.(*pickFirst)
	bw := &balancerWrapper{
		balancer:   bwb.b,
		pickfirst:  pickfirst,
		cc:         cc,
		targetAddr: targetAddr,
		startCh:    make(chan struct{}),
		conns:      make(map[resolver.Address]balancer.SubConn),
		connSt:     make(map[balancer.SubConn]*scState),
		csEvltr:    &balancer.ConnectivityStateEvaluator{},
		state:      connectivity.Idle,
	}
	cc.UpdateBalancerState(connectivity.Idle, bw)
	go bw.lbWatcher()
	return bw
}

func (bwb *balancerWrapperBuilder) Name() string {
	return "wrapper"
}

type scState struct {
	addr Address 
	s    connectivity.State
	down func(error)
}

type balancerWrapper struct {
	balancer  Balancer 
	pickfirst bool

	cc         balancer.ClientConn
	targetAddr string 

	mu     sync.Mutex
	conns  map[resolver.Address]balancer.SubConn
	connSt map[balancer.SubConn]*scState
	
	
	
	
	startCh chan struct{}

	
	csEvltr *balancer.ConnectivityStateEvaluator
	state   connectivity.State
}



func (bw *balancerWrapper) lbWatcher() {
	<-bw.startCh
	notifyCh := bw.balancer.Notify()
	if notifyCh == nil {
		
		a := resolver.Address{
			Addr: bw.targetAddr,
			Type: resolver.Backend,
		}
		sc, err := bw.cc.NewSubConn([]resolver.Address{a}, balancer.NewSubConnOptions{})
		if err != nil {
			grpclog.Warningf("Error creating connection to %v. Err: %v", a, err)
		} else {
			bw.mu.Lock()
			bw.conns[a] = sc
			bw.connSt[sc] = &scState{
				addr: Address{Addr: bw.targetAddr},
				s:    connectivity.Idle,
			}
			bw.mu.Unlock()
			sc.Connect()
		}
		return
	}

	for addrs := range notifyCh {
		grpclog.Infof("balancerWrapper: got update addr from Notify: %v\n", addrs)
		if bw.pickfirst {
			var (
				oldA  resolver.Address
				oldSC balancer.SubConn
			)
			bw.mu.Lock()
			for oldA, oldSC = range bw.conns {
				break
			}
			bw.mu.Unlock()
			if len(addrs) <= 0 {
				if oldSC != nil {
					
					bw.mu.Lock()
					delete(bw.conns, oldA)
					delete(bw.connSt, oldSC)
					bw.mu.Unlock()
					bw.cc.RemoveSubConn(oldSC)
				}
				continue
			}

			var newAddrs []resolver.Address
			for _, a := range addrs {
				newAddr := resolver.Address{
					Addr:       a.Addr,
					Type:       resolver.Backend, 
					ServerName: "",
					Metadata:   a.Metadata,
				}
				newAddrs = append(newAddrs, newAddr)
			}
			if oldSC == nil {
				
				sc, err := bw.cc.NewSubConn(newAddrs, balancer.NewSubConnOptions{})
				if err != nil {
					grpclog.Warningf("Error creating connection to %v. Err: %v", newAddrs, err)
				} else {
					bw.mu.Lock()
					
					
					
					bw.conns[resolver.Address{}] = sc
					bw.connSt[sc] = &scState{
						addr: addrs[0], 
						s:    connectivity.Idle,
					}
					bw.mu.Unlock()
					sc.Connect()
				}
			} else {
				bw.mu.Lock()
				bw.connSt[oldSC].addr = addrs[0]
				bw.mu.Unlock()
				oldSC.UpdateAddresses(newAddrs)
			}
		} else {
			var (
				add []resolver.Address 
				del []balancer.SubConn 
			)
			resAddrs := make(map[resolver.Address]Address)
			for _, a := range addrs {
				resAddrs[resolver.Address{
					Addr:       a.Addr,
					Type:       resolver.Backend, 
					ServerName: "",
					Metadata:   a.Metadata,
				}] = a
			}
			bw.mu.Lock()
			for a := range resAddrs {
				if _, ok := bw.conns[a]; !ok {
					add = append(add, a)
				}
			}
			for a, c := range bw.conns {
				if _, ok := resAddrs[a]; !ok {
					del = append(del, c)
					delete(bw.conns, a)
					
				}
			}
			bw.mu.Unlock()
			for _, a := range add {
				sc, err := bw.cc.NewSubConn([]resolver.Address{a}, balancer.NewSubConnOptions{})
				if err != nil {
					grpclog.Warningf("Error creating connection to %v. Err: %v", a, err)
				} else {
					bw.mu.Lock()
					bw.conns[a] = sc
					bw.connSt[sc] = &scState{
						addr: resAddrs[a],
						s:    connectivity.Idle,
					}
					bw.mu.Unlock()
					sc.Connect()
				}
			}
			for _, c := range del {
				bw.cc.RemoveSubConn(c)
			}
		}
	}
}

func (bw *balancerWrapper) HandleSubConnStateChange(sc balancer.SubConn, s connectivity.State) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	scSt, ok := bw.connSt[sc]
	if !ok {
		return
	}
	if s == connectivity.Idle {
		sc.Connect()
	}
	oldS := scSt.s
	scSt.s = s
	if oldS != connectivity.Ready && s == connectivity.Ready {
		scSt.down = bw.balancer.Up(scSt.addr)
	} else if oldS == connectivity.Ready && s != connectivity.Ready {
		if scSt.down != nil {
			scSt.down(errConnClosing)
		}
	}
	sa := bw.csEvltr.RecordTransition(oldS, s)
	if bw.state != sa {
		bw.state = sa
	}
	bw.cc.UpdateBalancerState(bw.state, bw)
	if s == connectivity.Shutdown {
		
		delete(bw.connSt, sc)
	}
}

func (bw *balancerWrapper) HandleResolvedAddrs([]resolver.Address, error) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	select {
	case <-bw.startCh:
	default:
		close(bw.startCh)
	}
	
	
}

func (bw *balancerWrapper) Close() {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	select {
	case <-bw.startCh:
	default:
		close(bw.startCh)
	}
	bw.balancer.Close()
}




func (bw *balancerWrapper) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	failfast := true 
	if ss, ok := rpcInfoFromContext(ctx); ok {
		failfast = ss.failfast
	}
	a, p, err := bw.balancer.Get(ctx, BalancerGetOptions{BlockingWait: !failfast})
	if err != nil {
		return nil, nil, err
	}
	var done func(balancer.DoneInfo)
	if p != nil {
		done = func(i balancer.DoneInfo) { p() }
	}
	var sc balancer.SubConn
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if bw.pickfirst {
		
		for _, sc = range bw.conns {
			break
		}
	} else {
		var ok bool
		sc, ok = bw.conns[resolver.Address{
			Addr:       a.Addr,
			Type:       resolver.Backend,
			ServerName: "",
			Metadata:   a.Metadata,
		}]
		if !ok && failfast {
			return nil, nil, status.Errorf(codes.Unavailable, "there is no connection available")
		}
		if s, ok := bw.connSt[sc]; failfast && (!ok || s.s != connectivity.Ready) {
			
			
			return nil, nil, status.Errorf(codes.Unavailable, "there is no connection available")
		}
	}

	return sc, done, nil
}
