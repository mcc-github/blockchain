



package dns

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(NewBuilder())
}

const (
	defaultPort = "443"
	defaultFreq = time.Minute * 30
	golang      = "GO"
	
	
	txtAttribute = "grpc_config="
)

var errMissingAddr = errors.New("missing address")


func NewBuilder() resolver.Builder {
	return &dnsBuilder{freq: defaultFreq}
}

type dnsBuilder struct {
	
	freq time.Duration
}


func (b *dnsBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	host, port, err := parseTarget(target.Endpoint)
	if err != nil {
		return nil, err
	}

	
	if net.ParseIP(host) != nil {
		host, _ = formatIP(host)
		addr := []resolver.Address{{Addr: host + ":" + port}}
		i := &ipResolver{
			cc: cc,
			ip: addr,
			rn: make(chan struct{}, 1),
			q:  make(chan struct{}),
		}
		cc.NewAddress(addr)
		go i.watcher()
		return i, nil
	}

	
	ctx, cancel := context.WithCancel(context.Background())
	d := &dnsResolver{
		freq:   b.freq,
		host:   host,
		port:   port,
		ctx:    ctx,
		cancel: cancel,
		cc:     cc,
		t:      time.NewTimer(0),
		rn:     make(chan struct{}, 1),
	}

	d.wg.Add(1)
	go d.watcher()
	return d, nil
}


func (b *dnsBuilder) Scheme() string {
	return "dns"
}


type ipResolver struct {
	cc resolver.ClientConn
	ip []resolver.Address
	
	rn chan struct{}
	q  chan struct{}
}


func (i *ipResolver) ResolveNow(opt resolver.ResolveNowOption) {
	select {
	case i.rn <- struct{}{}:
	default:
	}
}


func (i *ipResolver) Close() {
	close(i.q)
}

func (i *ipResolver) watcher() {
	for {
		select {
		case <-i.rn:
			i.cc.NewAddress(i.ip)
		case <-i.q:
			return
		}
	}
}


type dnsResolver struct {
	freq   time.Duration
	host   string
	port   string
	ctx    context.Context
	cancel context.CancelFunc
	cc     resolver.ClientConn
	
	rn chan struct{}
	t  *time.Timer
	
	
	
	
	
	
	wg sync.WaitGroup
}


func (d *dnsResolver) ResolveNow(opt resolver.ResolveNowOption) {
	select {
	case d.rn <- struct{}{}:
	default:
	}
}


func (d *dnsResolver) Close() {
	d.cancel()
	d.wg.Wait()
	d.t.Stop()
}

func (d *dnsResolver) watcher() {
	defer d.wg.Done()
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-d.t.C:
		case <-d.rn:
		}
		result, sc := d.lookup()
		
		d.t.Reset(d.freq)
		d.cc.NewServiceConfig(string(sc))
		d.cc.NewAddress(result)
	}
}

func (d *dnsResolver) lookupSRV() []resolver.Address {
	var newAddrs []resolver.Address
	_, srvs, err := lookupSRV(d.ctx, "grpclb", "tcp", d.host)
	if err != nil {
		grpclog.Infof("grpc: failed dns SRV record lookup due to %v.\n", err)
		return nil
	}
	for _, s := range srvs {
		lbAddrs, err := lookupHost(d.ctx, s.Target)
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
			newAddrs = append(newAddrs, resolver.Address{Addr: addr, Type: resolver.GRPCLB, ServerName: s.Target})
		}
	}
	return newAddrs
}

func (d *dnsResolver) lookupTXT() string {
	ss, err := lookupTXT(d.ctx, d.host)
	if err != nil {
		grpclog.Warningf("grpc: failed dns TXT record lookup due to %v.\n", err)
		return ""
	}
	var res string
	for _, s := range ss {
		res += s
	}

	
	if !strings.HasPrefix(res, txtAttribute) {
		grpclog.Warningf("grpc: TXT record %v missing %v attribute", res, txtAttribute)
		return ""
	}
	return strings.TrimPrefix(res, txtAttribute)
}

func (d *dnsResolver) lookupHost() []resolver.Address {
	var newAddrs []resolver.Address
	addrs, err := lookupHost(d.ctx, d.host)
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
		addr := a + ":" + d.port
		newAddrs = append(newAddrs, resolver.Address{Addr: addr})
	}
	return newAddrs
}

func (d *dnsResolver) lookup() ([]resolver.Address, string) {
	var newAddrs []resolver.Address
	newAddrs = d.lookupSRV()
	
	newAddrs = append(newAddrs, d.lookupHost()...)
	sc := d.lookupTXT()
	return newAddrs, canaryingSC(sc)
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
	if host, port, err = net.SplitHostPort(target); err == nil {
		
		if host == "" {
			
			host = "localhost"
		}
		if port == "" {
			
			port = defaultPort
		}
		return host, port, nil
	}
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil {
		
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
}

type rawChoice struct {
	ClientLanguage *[]string        `json:"clientLanguage,omitempty"`
	Percentage     *int             `json:"percentage,omitempty"`
	ClientHostName *[]string        `json:"clientHostName,omitempty"`
	ServiceConfig  *json.RawMessage `json:"serviceConfig,omitempty"`
}

func containsString(a *[]string, b string) bool {
	if a == nil {
		return true
	}
	for _, c := range *a {
		if c == b {
			return true
		}
	}
	return false
}

func chosenByPercentage(a *int) bool {
	if a == nil {
		return true
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	if r.Intn(100)+1 > *a {
		return false
	}
	return true
}

func canaryingSC(js string) string {
	if js == "" {
		return ""
	}
	var rcs []rawChoice
	err := json.Unmarshal([]byte(js), &rcs)
	if err != nil {
		grpclog.Warningf("grpc: failed to parse service config json string due to %v.\n", err)
		return ""
	}
	cliHostname, err := os.Hostname()
	if err != nil {
		grpclog.Warningf("grpc: failed to get client hostname due to %v.\n", err)
		return ""
	}
	var sc string
	for _, c := range rcs {
		if !containsString(c.ClientLanguage, golang) ||
			!chosenByPercentage(c.Percentage) ||
			!containsString(c.ClientHostName, cliHostname) ||
			c.ServiceConfig == nil {
			continue
		}
		sc = string(*c.ServiceConfig)
		break
	}
	return sc
}
