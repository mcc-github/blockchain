

package base

import (
	"context"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

type baseBuilder struct {
	name          string
	pickerBuilder PickerBuilder
	config        Config
}

func (bb *baseBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
	return &baseBalancer{
		cc:            cc,
		pickerBuilder: bb.pickerBuilder,

		subConns: make(map[resolver.Address]balancer.SubConn),
		scStates: make(map[balancer.SubConn]connectivity.State),
		csEvltr:  &connectivityStateEvaluator{},
		
		
		
		picker: NewErrPicker(balancer.ErrNoSubConnAvailable),
		config: bb.config,
	}
}

func (bb *baseBuilder) Name() string {
	return bb.name
}

type baseBalancer struct {
	cc            balancer.ClientConn
	pickerBuilder PickerBuilder

	csEvltr *connectivityStateEvaluator
	state   connectivity.State

	subConns map[resolver.Address]balancer.SubConn
	scStates map[balancer.SubConn]connectivity.State
	picker   balancer.Picker
	config   Config
}

func (b *baseBalancer) HandleResolvedAddrs(addrs []resolver.Address, err error) {
	if err != nil {
		grpclog.Infof("base.baseBalancer: HandleResolvedAddrs called with error %v", err)
		return
	}
	grpclog.Infoln("base.baseBalancer: got new resolved addresses: ", addrs)
	
	addrsSet := make(map[resolver.Address]struct{})
	for _, a := range addrs {
		addrsSet[a] = struct{}{}
		if _, ok := b.subConns[a]; !ok {
			
			sc, err := b.cc.NewSubConn([]resolver.Address{a}, balancer.NewSubConnOptions{HealthCheckEnabled: b.config.HealthCheck})
			if err != nil {
				grpclog.Warningf("base.baseBalancer: failed to create new SubConn: %v", err)
				continue
			}
			b.subConns[a] = sc
			b.scStates[sc] = connectivity.Idle
			sc.Connect()
		}
	}
	for a, sc := range b.subConns {
		
		if _, ok := addrsSet[a]; !ok {
			b.cc.RemoveSubConn(sc)
			delete(b.subConns, a)
			
			
		}
	}
}





func (b *baseBalancer) regeneratePicker() {
	if b.state == connectivity.TransientFailure {
		b.picker = NewErrPicker(balancer.ErrTransientFailure)
		return
	}
	readySCs := make(map[resolver.Address]balancer.SubConn)

	
	for addr, sc := range b.subConns {
		if st, ok := b.scStates[sc]; ok && st == connectivity.Ready {
			readySCs[addr] = sc
		}
	}
	b.picker = b.pickerBuilder.Build(readySCs)
}

func (b *baseBalancer) HandleSubConnStateChange(sc balancer.SubConn, s connectivity.State) {
	grpclog.Infof("base.baseBalancer: handle SubConn state change: %p, %v", sc, s)
	oldS, ok := b.scStates[sc]
	if !ok {
		grpclog.Infof("base.baseBalancer: got state changes for an unknown SubConn: %p, %v", sc, s)
		return
	}
	b.scStates[sc] = s
	switch s {
	case connectivity.Idle:
		sc.Connect()
	case connectivity.Shutdown:
		
		
		delete(b.scStates, sc)
	}

	oldAggrState := b.state
	b.state = b.csEvltr.recordTransition(oldS, s)

	
	
	
	
	
	if (s == connectivity.Ready) != (oldS == connectivity.Ready) ||
		(b.state == connectivity.TransientFailure) != (oldAggrState == connectivity.TransientFailure) {
		b.regeneratePicker()
	}

	b.cc.UpdateBalancerState(b.state, b.picker)
}



func (b *baseBalancer) Close() {
}


func NewErrPicker(err error) balancer.Picker {
	return &errPicker{err: err}
}

type errPicker struct {
	err error 
}

func (p *errPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	return nil, nil, p.err
}




type connectivityStateEvaluator struct {
	numReady            uint64 
	numConnecting       uint64 
	numTransientFailure uint64 
}









func (cse *connectivityStateEvaluator) recordTransition(oldState, newState connectivity.State) connectivity.State {
	
	for idx, state := range []connectivity.State{oldState, newState} {
		updateVal := 2*uint64(idx) - 1 
		switch state {
		case connectivity.Ready:
			cse.numReady += updateVal
		case connectivity.Connecting:
			cse.numConnecting += updateVal
		case connectivity.TransientFailure:
			cse.numTransientFailure += updateVal
		}
	}

	
	if cse.numReady > 0 {
		return connectivity.Ready
	}
	if cse.numConnecting > 0 {
		return connectivity.Connecting
	}
	return connectivity.TransientFailure
}
