

package grpc

import (
	"fmt"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)


type scStateUpdate struct {
	sc    balancer.SubConn
	state connectivity.State
}



type scStateUpdateBuffer struct {
	c       chan *scStateUpdate
	mu      sync.Mutex
	backlog []*scStateUpdate
}

func newSCStateUpdateBuffer() *scStateUpdateBuffer {
	return &scStateUpdateBuffer{
		c: make(chan *scStateUpdate, 1),
	}
}

func (b *scStateUpdateBuffer) put(t *scStateUpdate) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.backlog) == 0 {
		select {
		case b.c <- t:
			return
		default:
		}
	}
	b.backlog = append(b.backlog, t)
}

func (b *scStateUpdateBuffer) load() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			b.backlog[0] = nil
			b.backlog = b.backlog[1:]
		default:
		}
	}
}





func (b *scStateUpdateBuffer) get() <-chan *scStateUpdate {
	return b.c
}



type resolverUpdate struct {
	addrs []resolver.Address
	err   error
}



type ccBalancerWrapper struct {
	cc               *ClientConn
	balancer         balancer.Balancer
	stateChangeQueue *scStateUpdateBuffer
	resolverUpdateCh chan *resolverUpdate
	done             chan struct{}

	mu       sync.Mutex
	subConns map[*acBalancerWrapper]struct{}
}

func newCCBalancerWrapper(cc *ClientConn, b balancer.Builder, bopts balancer.BuildOptions) *ccBalancerWrapper {
	ccb := &ccBalancerWrapper{
		cc:               cc,
		stateChangeQueue: newSCStateUpdateBuffer(),
		resolverUpdateCh: make(chan *resolverUpdate, 1),
		done:             make(chan struct{}),
		subConns:         make(map[*acBalancerWrapper]struct{}),
	}
	go ccb.watcher()
	ccb.balancer = b.Build(ccb, bopts)
	return ccb
}



func (ccb *ccBalancerWrapper) watcher() {
	for {
		select {
		case t := <-ccb.stateChangeQueue.get():
			ccb.stateChangeQueue.load()
			select {
			case <-ccb.done:
				ccb.balancer.Close()
				return
			default:
			}
			ccb.balancer.HandleSubConnStateChange(t.sc, t.state)
		case t := <-ccb.resolverUpdateCh:
			select {
			case <-ccb.done:
				ccb.balancer.Close()
				return
			default:
			}
			ccb.balancer.HandleResolvedAddrs(t.addrs, t.err)
		case <-ccb.done:
		}

		select {
		case <-ccb.done:
			ccb.balancer.Close()
			ccb.mu.Lock()
			scs := ccb.subConns
			ccb.subConns = nil
			ccb.mu.Unlock()
			for acbw := range scs {
				ccb.cc.removeAddrConn(acbw.getAddrConn(), errConnDrain)
			}
			return
		default:
		}
	}
}

func (ccb *ccBalancerWrapper) close() {
	close(ccb.done)
}

func (ccb *ccBalancerWrapper) handleSubConnStateChange(sc balancer.SubConn, s connectivity.State) {
	
	
	
	
	
	
	
	if sc == nil {
		return
	}
	ccb.stateChangeQueue.put(&scStateUpdate{
		sc:    sc,
		state: s,
	})
}

func (ccb *ccBalancerWrapper) handleResolvedAddrs(addrs []resolver.Address, err error) {
	select {
	case <-ccb.resolverUpdateCh:
	default:
	}
	ccb.resolverUpdateCh <- &resolverUpdate{
		addrs: addrs,
		err:   err,
	}
}

func (ccb *ccBalancerWrapper) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
	if len(addrs) <= 0 {
		return nil, fmt.Errorf("grpc: cannot create SubConn with empty address list")
	}
	ccb.mu.Lock()
	defer ccb.mu.Unlock()
	if ccb.subConns == nil {
		return nil, fmt.Errorf("grpc: ClientConn balancer wrapper was closed")
	}
	ac, err := ccb.cc.newAddrConn(addrs)
	if err != nil {
		return nil, err
	}
	acbw := &acBalancerWrapper{ac: ac}
	acbw.ac.mu.Lock()
	ac.acbw = acbw
	acbw.ac.mu.Unlock()
	ccb.subConns[acbw] = struct{}{}
	return acbw, nil
}

func (ccb *ccBalancerWrapper) RemoveSubConn(sc balancer.SubConn) {
	acbw, ok := sc.(*acBalancerWrapper)
	if !ok {
		return
	}
	ccb.mu.Lock()
	defer ccb.mu.Unlock()
	if ccb.subConns == nil {
		return
	}
	delete(ccb.subConns, acbw)
	ccb.cc.removeAddrConn(acbw.getAddrConn(), errConnDrain)
}

func (ccb *ccBalancerWrapper) UpdateBalancerState(s connectivity.State, p balancer.Picker) {
	ccb.mu.Lock()
	defer ccb.mu.Unlock()
	if ccb.subConns == nil {
		return
	}
	ccb.cc.csMgr.updateState(s)
	ccb.cc.blockingpicker.updatePicker(p)
}

func (ccb *ccBalancerWrapper) ResolveNow(o resolver.ResolveNowOption) {
	ccb.cc.resolveNow(o)
}

func (ccb *ccBalancerWrapper) Target() string {
	return ccb.cc.target
}



type acBalancerWrapper struct {
	mu sync.Mutex
	ac *addrConn
}

func (acbw *acBalancerWrapper) UpdateAddresses(addrs []resolver.Address) {
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	if len(addrs) <= 0 {
		acbw.ac.tearDown(errConnDrain)
		return
	}
	if !acbw.ac.tryUpdateAddrs(addrs) {
		cc := acbw.ac.cc
		acbw.ac.mu.Lock()
		
		
		
		
		
		acbw.ac.acbw = nil
		acbw.ac.mu.Unlock()
		acState := acbw.ac.getState()
		acbw.ac.tearDown(errConnDrain)

		if acState == connectivity.Shutdown {
			return
		}

		ac, err := cc.newAddrConn(addrs)
		if err != nil {
			grpclog.Warningf("acBalancerWrapper: UpdateAddresses: failed to newAddrConn: %v", err)
			return
		}
		acbw.ac = ac
		ac.mu.Lock()
		ac.acbw = acbw
		ac.mu.Unlock()
		if acState != connectivity.Idle {
			ac.connect()
		}
	}
}

func (acbw *acBalancerWrapper) Connect() {
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	acbw.ac.connect()
}

func (acbw *acBalancerWrapper) getAddrConn() *addrConn {
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	return acbw.ac
}
