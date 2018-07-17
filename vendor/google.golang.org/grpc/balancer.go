

package grpc

import (
	"fmt"
	"net"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/naming"
	"google.golang.org/grpc/status"
)



type Address struct {
	
	Addr string
	
	
	Metadata interface{}
}


type BalancerConfig struct {
	
	
	
	DialCreds credentials.TransportCredentials
	
	
	
	Dialer func(context.Context, string) (net.Conn, error)
}



type BalancerGetOptions struct {
	
	
	BlockingWait bool
}



type Balancer interface {
	
	
	
	Start(target string, config BalancerConfig) error
	
	
	
	
	
	Up(addr Address) (down func(error))
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Get(ctx context.Context, opts BalancerGetOptions) (addr Address, put func(), err error)
	
	
	
	
	
	
	
	
	
	Notify() <-chan []Address
	
	Close() error
}



type downErr struct {
	timeout   bool
	temporary bool
	desc      string
}

func (e downErr) Error() string   { return e.desc }
func (e downErr) Timeout() bool   { return e.timeout }
func (e downErr) Temporary() bool { return e.temporary }

func downErrorf(timeout, temporary bool, format string, a ...interface{}) downErr {
	return downErr{
		timeout:   timeout,
		temporary: temporary,
		desc:      fmt.Sprintf(format, a...),
	}
}



func RoundRobin(r naming.Resolver) Balancer {
	return &roundRobin{r: r}
}

type addrInfo struct {
	addr      Address
	connected bool
}

type roundRobin struct {
	r      naming.Resolver
	w      naming.Watcher
	addrs  []*addrInfo 
	mu     sync.Mutex
	addrCh chan []Address 
	next   int            
	waitCh chan struct{}  
	done   bool           
}

func (rr *roundRobin) watchAddrUpdates() error {
	updates, err := rr.w.Next()
	if err != nil {
		grpclog.Warningf("grpc: the naming watcher stops working due to %v.", err)
		return err
	}
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for _, update := range updates {
		addr := Address{
			Addr:     update.Addr,
			Metadata: update.Metadata,
		}
		switch update.Op {
		case naming.Add:
			var exist bool
			for _, v := range rr.addrs {
				if addr == v.addr {
					exist = true
					grpclog.Infoln("grpc: The name resolver wanted to add an existing address: ", addr)
					break
				}
			}
			if exist {
				continue
			}
			rr.addrs = append(rr.addrs, &addrInfo{addr: addr})
		case naming.Delete:
			for i, v := range rr.addrs {
				if addr == v.addr {
					copy(rr.addrs[i:], rr.addrs[i+1:])
					rr.addrs = rr.addrs[:len(rr.addrs)-1]
					break
				}
			}
		default:
			grpclog.Errorln("Unknown update.Op ", update.Op)
		}
	}
	
	open := make([]Address, len(rr.addrs))
	for i, v := range rr.addrs {
		open[i] = v.addr
	}
	if rr.done {
		return ErrClientConnClosing
	}
	select {
	case <-rr.addrCh:
	default:
	}
	rr.addrCh <- open
	return nil
}

func (rr *roundRobin) Start(target string, config BalancerConfig) error {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if rr.done {
		return ErrClientConnClosing
	}
	if rr.r == nil {
		
		
		
		rr.addrs = append(rr.addrs, &addrInfo{addr: Address{Addr: target}})
		return nil
	}
	w, err := rr.r.Resolve(target)
	if err != nil {
		return err
	}
	rr.w = w
	rr.addrCh = make(chan []Address, 1)
	go func() {
		for {
			if err := rr.watchAddrUpdates(); err != nil {
				return
			}
		}
	}()
	return nil
}



func (rr *roundRobin) Up(addr Address) func(error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	var cnt int
	for _, a := range rr.addrs {
		if a.addr == addr {
			if a.connected {
				return nil
			}
			a.connected = true
		}
		if a.connected {
			cnt++
		}
	}
	
	if cnt == 1 && rr.waitCh != nil {
		close(rr.waitCh)
		rr.waitCh = nil
	}
	return func(err error) {
		rr.down(addr, err)
	}
}


func (rr *roundRobin) down(addr Address, err error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for _, a := range rr.addrs {
		if addr == a.addr {
			a.connected = false
			break
		}
	}
}


func (rr *roundRobin) Get(ctx context.Context, opts BalancerGetOptions) (addr Address, put func(), err error) {
	var ch chan struct{}
	rr.mu.Lock()
	if rr.done {
		rr.mu.Unlock()
		err = ErrClientConnClosing
		return
	}

	if len(rr.addrs) > 0 {
		if rr.next >= len(rr.addrs) {
			rr.next = 0
		}
		next := rr.next
		for {
			a := rr.addrs[next]
			next = (next + 1) % len(rr.addrs)
			if a.connected {
				addr = a.addr
				rr.next = next
				rr.mu.Unlock()
				return
			}
			if next == rr.next {
				
				break
			}
		}
	}
	if !opts.BlockingWait {
		if len(rr.addrs) == 0 {
			rr.mu.Unlock()
			err = status.Errorf(codes.Unavailable, "there is no address available")
			return
		}
		
		addr = rr.addrs[rr.next].addr
		rr.next++
		rr.mu.Unlock()
		return
	}
	
	if rr.waitCh == nil {
		ch = make(chan struct{})
		rr.waitCh = ch
	} else {
		ch = rr.waitCh
	}
	rr.mu.Unlock()
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-ch:
			rr.mu.Lock()
			if rr.done {
				rr.mu.Unlock()
				err = ErrClientConnClosing
				return
			}

			if len(rr.addrs) > 0 {
				if rr.next >= len(rr.addrs) {
					rr.next = 0
				}
				next := rr.next
				for {
					a := rr.addrs[next]
					next = (next + 1) % len(rr.addrs)
					if a.connected {
						addr = a.addr
						rr.next = next
						rr.mu.Unlock()
						return
					}
					if next == rr.next {
						
						break
					}
				}
			}
			
			if rr.waitCh == nil {
				ch = make(chan struct{})
				rr.waitCh = ch
			} else {
				ch = rr.waitCh
			}
			rr.mu.Unlock()
		}
	}
}

func (rr *roundRobin) Notify() <-chan []Address {
	return rr.addrCh
}

func (rr *roundRobin) Close() error {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if rr.done {
		return errBalancerClosed
	}
	rr.done = true
	if rr.w != nil {
		rr.w.Close()
	}
	if rr.waitCh != nil {
		close(rr.waitCh)
		rr.waitCh = nil
	}
	if rr.addrCh != nil {
		close(rr.addrCh)
	}
	return nil
}




type pickFirst struct {
	*roundRobin
}

func pickFirstBalancerV1(r naming.Resolver) Balancer {
	return &pickFirst{&roundRobin{r: r}}
}
