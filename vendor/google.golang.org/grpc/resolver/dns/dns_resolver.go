



package dns

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/grpcrand"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(NewBuilder())
}

const (
	defaultPort       = "443"
	defaultFreq       = time.Minute * 30
	defaultDNSSvrPort = "53"
	golang            = "GO"
	
	txtPrefix = "_grpc_config."
	
	
	txtAttribute = "grpc_config="
)

var (
	errMissingAddr = errors.New("dns resolver: missing address")

	
	
	
	
	errEndsWithColon = errors.New("dns resolver: missing port after port-separator colon")
)

var (
	defaultResolver netResolver = net.DefaultResolver
)

var customAuthorityDialler = func(authority string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		var dialer net.Dialer
		return dialer.DialContext(ctx, network, authority)
	}
}

var customAuthorityResolver = func(authority string) (netResolver, error) {
	host, port, err := parseTarget(authority, defaultDNSSvrPort)
	if err != nil {
		return nil, err
	}

	authorityWithPort := net.JoinHostPort(host, port)

	return &net.Resolver{
		PreferGo: true,
		Dial:     customAuthorityDialler(authorityWithPort),
	}, nil
}


func NewBuilder() resolver.Builder {
	return &dnsBuilder{minFreq: defaultFreq}
}

type dnsBuilder struct {
	
	minFreq time.Duration
}


func (b *dnsBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	host, port, err := parseTarget(target.Endpoint, defaultPort)
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
		freq:                 b.minFreq,
		backoff:              backoff.Exponential{MaxDelay: b.minFreq},
		host:                 host,
		port:                 port,
		ctx:                  ctx,
		cancel:               cancel,
		cc:                   cc,
		t:                    time.NewTimer(0),
		rn:                   make(chan struct{}, 1),
		disableServiceConfig: opts.DisableServiceConfig,
	}

	if target.Authority == "" {
		d.resolver = defaultResolver
	} else {
		d.resolver, err = customAuthorityResolver(target.Authority)
		if err != nil {
			return nil, err
		}
	}

	d.wg.Add(1)
	go d.watcher()
	return d, nil
}


func (b *dnsBuilder) Scheme() string {
	return "dns"
}

type netResolver interface {
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
	LookupTXT(ctx context.Context, name string) (txts []string, err error)
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
	freq       time.Duration
	backoff    backoff.Exponential
	retryCount int
	host       string
	port       string
	resolver   netResolver
	ctx        context.Context
	cancel     context.CancelFunc
	cc         resolver.ClientConn
	
	rn chan struct{}
	t  *time.Timer
	
	
	
	
	
	
	wg                   sync.WaitGroup
	disableServiceConfig bool
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
		
		
		if len(result) == 0 {
			d.retryCount++
			d.t.Reset(d.backoff.Backoff(d.retryCount))
		} else {
			d.retryCount = 0
			d.t.Reset(d.freq)
		}
		d.cc.NewServiceConfig(sc)
		d.cc.NewAddress(result)
	}
}

func (d *dnsResolver) lookupSRV() []resolver.Address {
	var newAddrs []resolver.Address
	_, srvs, err := d.resolver.LookupSRV(d.ctx, "grpclb", "tcp", d.host)
	if err != nil {
		grpclog.Infof("grpc: failed dns SRV record lookup due to %v.\n", err)
		return nil
	}
	for _, s := range srvs {
		lbAddrs, err := d.resolver.LookupHost(d.ctx, s.Target)
		if err != nil {
			grpclog.Infof("grpc: failed load balancer address dns lookup due to %v.\n", err)
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
	ss, err := d.resolver.LookupTXT(d.ctx, txtPrefix+d.host)
	if err != nil {
		grpclog.Infof("grpc: failed dns TXT record lookup due to %v.\n", err)
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
	addrs, err := d.resolver.LookupHost(d.ctx, d.host)
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
	newAddrs := d.lookupSRV()
	
	newAddrs = append(newAddrs, d.lookupHost()...)
	if d.disableServiceConfig {
		return newAddrs, ""
	}
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










func parseTarget(target, defaultPort string) (host, port string, err error) {
	if target == "" {
		return "", "", errMissingAddr
	}
	if ip := net.ParseIP(target); ip != nil {
		
		return target, defaultPort, nil
	}
	if host, port, err = net.SplitHostPort(target); err == nil {
		if port == "" {
			
			return "", "", errEndsWithColon
		}
		
		if host == "" {
			
			host = "localhost"
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
	return grpcrand.Intn(100)+1 <= *a
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
