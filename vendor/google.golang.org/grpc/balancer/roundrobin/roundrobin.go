




package roundrobin

import (
	"context"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/grpcrand"
	"google.golang.org/grpc/resolver"
)


const Name = "round_robin"


func newBuilder() balancer.Builder {
	return base.NewBalancerBuilderWithConfig(Name, &rrPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct{}

func (*rrPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	grpclog.Infof("roundrobinPicker: newPicker called with readySCs: %v", readySCs)
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for _, sc := range readySCs {
		scs = append(scs, sc)
	}
	return &rrPicker{
		subConns: scs,
		
		
		
		next: grpcrand.Intn(len(scs)),
	}
}

type rrPicker struct {
	
	
	
	subConns []balancer.SubConn

	mu   sync.Mutex
	next int
}

func (p *rrPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	p.mu.Lock()
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return sc, nil, nil
}
