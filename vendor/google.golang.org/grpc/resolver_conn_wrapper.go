

package grpc

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/resolver"
)



type ccResolverWrapper struct {
	cc                 *ClientConn
	resolver           resolver.Resolver
	addrCh             chan []resolver.Address
	scCh               chan string
	done               chan struct{}
	lastAddressesCount int
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
		done:   make(chan struct{}),
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
	close(ccr.done)
}


func (ccr *ccResolverWrapper) NewAddress(addrs []resolver.Address) {
	select {
	case <-ccr.done:
		return
	default:
	}
	grpclog.Infof("ccResolverWrapper: sending new addresses to cc: %v", addrs)
	if channelz.IsOn() {
		ccr.addChannelzTraceEvent(addrs)
	}
	ccr.cc.handleResolvedAddrs(addrs, nil)
}



func (ccr *ccResolverWrapper) NewServiceConfig(sc string) {
	select {
	case <-ccr.done:
		return
	default:
	}
	grpclog.Infof("ccResolverWrapper: got new service config: %v", sc)
	ccr.cc.handleServiceConfig(sc)
}

func (ccr *ccResolverWrapper) addChannelzTraceEvent(addrs []resolver.Address) {
	if len(addrs) == 0 && ccr.lastAddressesCount != 0 {
		channelz.AddTraceEvent(ccr.cc.channelzID, &channelz.TraceEventDesc{
			Desc:     "Resolver returns an empty address list",
			Severity: channelz.CtWarning,
		})
	} else if len(addrs) != 0 && ccr.lastAddressesCount == 0 {
		var s string
		for i, a := range addrs {
			if a.ServerName != "" {
				s += a.Addr + "(" + a.ServerName + ")"
			} else {
				s += a.Addr
			}
			if i != len(addrs)-1 {
				s += " "
			}
		}
		channelz.AddTraceEvent(ccr.cc.channelzID, &channelz.TraceEventDesc{
			Desc:     fmt.Sprintf("Resolver returns a non-empty address list (previous one was empty) %q", s),
			Severity: channelz.CtINFO,
		})
	}
	ccr.lastAddressesCount = len(addrs)
}
