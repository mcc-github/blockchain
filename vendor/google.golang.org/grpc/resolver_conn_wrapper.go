

package grpc

import (
	"fmt"
	"strings"
	"sync/atomic"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/resolver"
)



type ccResolverWrapper struct {
	cc       *ClientConn
	resolver resolver.Resolver
	addrCh   chan []resolver.Address
	scCh     chan string
	done     uint32 
	curState resolver.State
}



func split2(s, sep string) (string, string, bool) {
	spl := strings.SplitN(s, sep, 2)
	if len(spl) < 2 {
		return "", "", false
	}
	return spl[0], spl[1], true
}






func parseTarget(target string) (ret resolver.Target) {
	var ok bool
	ret.Scheme, ret.Endpoint, ok = split2(target, "://")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	ret.Authority, ret.Endpoint, ok = split2(ret.Endpoint, "/")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	return ret
}







func newCCResolverWrapper(cc *ClientConn) (*ccResolverWrapper, error) {
	rb := cc.dopts.resolverBuilder
	if rb == nil {
		return nil, fmt.Errorf("could not get resolver for scheme: %q", cc.parsedTarget.Scheme)
	}

	ccr := &ccResolverWrapper{
		cc:     cc,
		addrCh: make(chan []resolver.Address, 1),
		scCh:   make(chan string, 1),
	}

	var err error
	ccr.resolver, err = rb.Build(cc.parsedTarget, ccr, resolver.BuildOption{DisableServiceConfig: cc.dopts.disableServiceConfig})
	if err != nil {
		return nil, err
	}
	return ccr, nil
}

func (ccr *ccResolverWrapper) resolveNow(o resolver.ResolveNowOption) {
	ccr.resolver.ResolveNow(o)
}

func (ccr *ccResolverWrapper) close() {
	ccr.resolver.Close()
	atomic.StoreUint32(&ccr.done, 1)
}

func (ccr *ccResolverWrapper) isDone() bool {
	return atomic.LoadUint32(&ccr.done) == 1
}

func (ccr *ccResolverWrapper) UpdateState(s resolver.State) {
	if ccr.isDone() {
		return
	}
	grpclog.Infof("ccResolverWrapper: sending update to cc: %v", s)
	if channelz.IsOn() {
		ccr.addChannelzTraceEvent(s)
	}
	ccr.cc.updateResolverState(s)
	ccr.curState = s
}


func (ccr *ccResolverWrapper) NewAddress(addrs []resolver.Address) {
	if ccr.isDone() {
		return
	}
	grpclog.Infof("ccResolverWrapper: sending new addresses to cc: %v", addrs)
	if channelz.IsOn() {
		ccr.addChannelzTraceEvent(resolver.State{Addresses: addrs, ServiceConfig: ccr.curState.ServiceConfig})
	}
	ccr.curState.Addresses = addrs
	ccr.cc.updateResolverState(ccr.curState)
}



func (ccr *ccResolverWrapper) NewServiceConfig(sc string) {
	if ccr.isDone() {
		return
	}
	grpclog.Infof("ccResolverWrapper: got new service config: %v", sc)
	if channelz.IsOn() {
		ccr.addChannelzTraceEvent(resolver.State{Addresses: ccr.curState.Addresses, ServiceConfig: sc})
	}
	ccr.curState.ServiceConfig = sc
	ccr.cc.updateResolverState(ccr.curState)
}

func (ccr *ccResolverWrapper) addChannelzTraceEvent(s resolver.State) {
	if s.ServiceConfig == ccr.curState.ServiceConfig && (len(ccr.curState.Addresses) == 0) == (len(s.Addresses) == 0) {
		return
	}
	var updates []string
	if s.ServiceConfig != ccr.curState.ServiceConfig {
		updates = append(updates, "service config updated")
	}
	if len(ccr.curState.Addresses) > 0 && len(s.Addresses) == 0 {
		updates = append(updates, "resolver returned an empty address list")
	} else if len(ccr.curState.Addresses) == 0 && len(s.Addresses) > 0 {
		updates = append(updates, "resolver returned new addresses")
	}
	channelz.AddTraceEvent(ccr.cc.channelzID, &channelz.TraceEventDesc{
		Desc:     fmt.Sprintf("Resolver state updated: %+v (%v)", s, strings.Join(updates, "; ")),
		Severity: channelz.CtINFO,
	})
}
