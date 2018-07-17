

package naming

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
)

const (
	defaultPort = "443"
	defaultFreq = time.Minute * 30
)

var (
	errMissingAddr  = errors.New("missing address")
	errWatcherClose = errors.New("watcher has been closed")
)



func NewDNSResolverWithFreq(freq time.Duration) (Resolver, error) {
	return &dnsResolver{freq: freq}, nil
}



func NewDNSResolver() (Resolver, error) {
	return NewDNSResolverWithFreq(defaultFreq)
}


type dnsResolver struct {
	
	freq time.Duration
}




func formatIP(addr string) (addrIP string, ok bool) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", false
	}
	if ip.To4() != nil {
		return addr, true
	}
	return "[" + addr + "]", true
}











func parseTarget(target string) (host, port string, err error) {
	if target == "" {
		return "", "", errMissingAddr
	}

	if ip := net.ParseIP(target); ip != nil {
		
		return target, defaultPort, nil
	}
	if host, port, err := net.SplitHostPort(target); err == nil {
		
		if host == "" {
			
			host = "localhost"
		}
		if port == "" {
			
			port = defaultPort
		}
		return host, port, nil
	}
	if host, port, err := net.SplitHostPort(target + ":" + defaultPort); err == nil {
		
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v", target)
}


func (r *dnsResolver) Resolve(target string) (Watcher, error) {
	host, port, err := parseTarget(target)
	if err != nil {
		return nil, err
	}

	if net.ParseIP(host) != nil {
		ipWatcher := &ipWatcher{
			updateChan: make(chan *Update, 1),
		}
		host, _ = formatIP(host)
		ipWatcher.updateChan <- &Update{Op: Add, Addr: host + ":" + port}
		return ipWatcher, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &dnsWatcher{
		r:      r,
		host:   host,
		port:   port,
		ctx:    ctx,
		cancel: cancel,
		t:      time.NewTimer(0),
	}, nil
}


type dnsWatcher struct {
	r    *dnsResolver
	host string
	port string
	
	curAddrs map[string]*Update
	ctx      context.Context
	cancel   context.CancelFunc
	t        *time.Timer
}


type ipWatcher struct {
	updateChan chan *Update
}





func (i *ipWatcher) Next() ([]*Update, error) {
	u, ok := <-i.updateChan
	if !ok {
		return nil, errWatcherClose
	}
	return []*Update{u}, nil
}


func (i *ipWatcher) Close() {
	close(i.updateChan)
}


type AddressType uint8

const (
	
	Backend AddressType = iota
	
	GRPCLB
)




type AddrMetadataGRPCLB struct {
	
	AddrType AddressType
	
	ServerName string
}



func (w *dnsWatcher) compileUpdate(newAddrs map[string]*Update) []*Update {
	var res []*Update
	for a, u := range w.curAddrs {
		if _, ok := newAddrs[a]; !ok {
			u.Op = Delete
			res = append(res, u)
		}
	}
	for a, u := range newAddrs {
		if _, ok := w.curAddrs[a]; !ok {
			res = append(res, u)
		}
	}
	return res
}

func (w *dnsWatcher) lookupSRV() map[string]*Update {
	newAddrs := make(map[string]*Update)
	_, srvs, err := lookupSRV(w.ctx, "grpclb", "tcp", w.host)
	if err != nil {
		grpclog.Infof("grpc: failed dns SRV record lookup due to %v.\n", err)
		return nil
	}
	for _, s := range srvs {
		lbAddrs, err := lookupHost(w.ctx, s.Target)
		if err != nil {
			grpclog.Warningf("grpc: failed load banlacer address dns lookup due to %v.\n", err)
			continue
		}
		for _, a := range lbAddrs {
			a, ok := formatIP(a)
			if !ok {
				grpclog.Errorf("grpc: failed IP parsing due to %v.\n", err)
				continue
			}
			addr := a + ":" + strconv.Itoa(int(s.Port))
			newAddrs[addr] = &Update{Addr: addr,
				Metadata: AddrMetadataGRPCLB{AddrType: GRPCLB, ServerName: s.Target}}
		}
	}
	return newAddrs
}

func (w *dnsWatcher) lookupHost() map[string]*Update {
	newAddrs := make(map[string]*Update)
	addrs, err := lookupHost(w.ctx, w.host)
	if err != nil {
		grpclog.Warningf("grpc: failed dns A record lookup due to %v.\n", err)
		return nil
	}
	for _, a := range addrs {
		a, ok := formatIP(a)
		if !ok {
			grpclog.Errorf("grpc: failed IP parsing due to %v.\n", err)
			continue
		}
		addr := a + ":" + w.port
		newAddrs[addr] = &Update{Addr: addr}
	}
	return newAddrs
}

func (w *dnsWatcher) lookup() []*Update {
	newAddrs := w.lookupSRV()
	if newAddrs == nil {
		
		
		
		newAddrs = w.lookupHost()
	}
	result := w.compileUpdate(newAddrs)
	w.curAddrs = newAddrs
	return result
}



func (w *dnsWatcher) Next() ([]*Update, error) {
	for {
		select {
		case <-w.ctx.Done():
			return nil, errWatcherClose
		case <-w.t.C:
		}
		result := w.lookup()
		
		w.t.Reset(w.r.freq)
		if len(result) > 0 {
			return result, nil
		}
	}
}

func (w *dnsWatcher) Close() {
	w.cancel()
}
