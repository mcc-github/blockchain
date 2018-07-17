

package grpc

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

































type lbManualResolver struct {
	scheme string
	ccr    resolver.ClientConn

	ccb balancer.ClientConn
}

func (r *lbManualResolver) Build(_ resolver.Target, cc resolver.ClientConn, _ resolver.BuildOption) (resolver.Resolver, error) {
	r.ccr = cc
	return r, nil
}

func (r *lbManualResolver) Scheme() string {
	return r.scheme
}


func (r *lbManualResolver) ResolveNow(o resolver.ResolveNowOption) {
	r.ccb.ResolveNow(o)
}


func (*lbManualResolver) Close() {}


func (r *lbManualResolver) NewAddress(addrs []resolver.Address) {
	r.ccr.NewAddress(addrs)
}


func (r *lbManualResolver) NewServiceConfig(sc string) {
	r.ccr.NewServiceConfig(sc)
}
