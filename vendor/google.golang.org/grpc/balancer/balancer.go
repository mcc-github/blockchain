



package balancer

import (
	"errors"
	"net"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
)

var (
	
	m = make(map[string]Builder)
)







func Register(b Builder) {
	m[strings.ToLower(b.Name())] = b
}




func Get(name string) Builder {
	if b, ok := m[strings.ToLower(name)]; ok {
		return b
	}
	return nil
}



















type SubConn interface {
	
	
	
	
	
	
	
	UpdateAddresses([]resolver.Address)
	
	Connect()
}


type NewSubConnOptions struct{}







type ClientConn interface {
	
	
	
	NewSubConn([]resolver.Address, NewSubConnOptions) (SubConn, error)
	
	
	RemoveSubConn(SubConn)

	
	
	
	
	
	UpdateBalancerState(s connectivity.State, p Picker)

	
	ResolveNow(resolver.ResolveNowOption)

	
	Target() string
}


type BuildOptions struct {
	
	
	
	DialCreds credentials.TransportCredentials
	
	
	
	Dialer func(context.Context, string) (net.Conn, error)
	
	ChannelzParentID int64
}


type Builder interface {
	
	Build(cc ClientConn, opts BuildOptions) Balancer
	
	
	Name() string
}


type PickOptions struct {
	
	
	FullMethodName string
}


type DoneInfo struct {
	
	Err error
	
	BytesSent bool
	
	BytesReceived bool
}

var (
	
	
	ErrNoSubConnAvailable = errors.New("no SubConn is available")
	
	
	ErrTransientFailure = errors.New("all SubConns are in TransientFailure")
)






type Picker interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Pick(ctx context.Context, opts PickOptions) (conn SubConn, done func(DoneInfo), err error)
}









type Balancer interface {
	
	
	
	
	
	
	HandleSubConnStateChange(sc SubConn, state connectivity.State)
	
	
	
	
	
	HandleResolvedAddrs([]resolver.Address, error)
	
	
	Close()
}





type ConnectivityStateEvaluator struct {
	numReady            uint64 
	numConnecting       uint64 
	numTransientFailure uint64 
}









func (cse *ConnectivityStateEvaluator) RecordTransition(oldState, newState connectivity.State) connectivity.State {
	
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
