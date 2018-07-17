





package http2

import (
	"crypto/tls"
	"net/http"
	"sync"
)


type ClientConnPool interface {
	GetClientConn(req *http.Request, addr string) (*ClientConn, error)
	MarkDead(*ClientConn)
}



type clientConnPoolIdleCloser interface {
	ClientConnPool
	closeIdleConnections()
}

var (
	_ clientConnPoolIdleCloser = (*clientConnPool)(nil)
	_ clientConnPoolIdleCloser = noDialClientConnPool{}
)


type clientConnPool struct {
	t *Transport

	mu sync.Mutex 
	
	
	conns        map[string][]*ClientConn 
	dialing      map[string]*dialCall     
	keys         map[*ClientConn][]string
	addConnCalls map[string]*addConnCall 
}

func (p *clientConnPool) GetClientConn(req *http.Request, addr string) (*ClientConn, error) {
	return p.getClientConn(req, addr, dialOnMiss)
}

const (
	dialOnMiss   = true
	noDialOnMiss = false
)

func (p *clientConnPool) getClientConn(req *http.Request, addr string, dialOnMiss bool) (*ClientConn, error) {
	if isConnectionCloseRequest(req) && dialOnMiss {
		
		const singleUse = true
		cc, err := p.t.dialClientConn(addr, singleUse)
		if err != nil {
			return nil, err
		}
		return cc, nil
	}
	p.mu.Lock()
	for _, cc := range p.conns[addr] {
		if cc.CanTakeNewRequest() {
			p.mu.Unlock()
			return cc, nil
		}
	}
	if !dialOnMiss {
		p.mu.Unlock()
		return nil, ErrNoCachedConn
	}
	call := p.getStartDialLocked(addr)
	p.mu.Unlock()
	<-call.done
	return call.res, call.err
}


type dialCall struct {
	p    *clientConnPool
	done chan struct{} 
	res  *ClientConn   
	err  error         
}


func (p *clientConnPool) getStartDialLocked(addr string) *dialCall {
	if call, ok := p.dialing[addr]; ok {
		
		return call
	}
	call := &dialCall{p: p, done: make(chan struct{})}
	if p.dialing == nil {
		p.dialing = make(map[string]*dialCall)
	}
	p.dialing[addr] = call
	go call.dial(addr)
	return call
}


func (c *dialCall) dial(addr string) {
	const singleUse = false 
	c.res, c.err = c.p.t.dialClientConn(addr, singleUse)
	close(c.done)

	c.p.mu.Lock()
	delete(c.p.dialing, addr)
	if c.err == nil {
		c.p.addConnLocked(addr, c.res)
	}
	c.p.mu.Unlock()
}









func (p *clientConnPool) addConnIfNeeded(key string, t *Transport, c *tls.Conn) (used bool, err error) {
	p.mu.Lock()
	for _, cc := range p.conns[key] {
		if cc.CanTakeNewRequest() {
			p.mu.Unlock()
			return false, nil
		}
	}
	call, dup := p.addConnCalls[key]
	if !dup {
		if p.addConnCalls == nil {
			p.addConnCalls = make(map[string]*addConnCall)
		}
		call = &addConnCall{
			p:    p,
			done: make(chan struct{}),
		}
		p.addConnCalls[key] = call
		go call.run(t, key, c)
	}
	p.mu.Unlock()

	<-call.done
	if call.err != nil {
		return false, call.err
	}
	return !dup, nil
}

type addConnCall struct {
	p    *clientConnPool
	done chan struct{} 
	err  error
}

func (c *addConnCall) run(t *Transport, key string, tc *tls.Conn) {
	cc, err := t.NewClientConn(tc)

	p := c.p
	p.mu.Lock()
	if err != nil {
		c.err = err
	} else {
		p.addConnLocked(key, cc)
	}
	delete(p.addConnCalls, key)
	p.mu.Unlock()
	close(c.done)
}

func (p *clientConnPool) addConn(key string, cc *ClientConn) {
	p.mu.Lock()
	p.addConnLocked(key, cc)
	p.mu.Unlock()
}


func (p *clientConnPool) addConnLocked(key string, cc *ClientConn) {
	for _, v := range p.conns[key] {
		if v == cc {
			return
		}
	}
	if p.conns == nil {
		p.conns = make(map[string][]*ClientConn)
	}
	if p.keys == nil {
		p.keys = make(map[*ClientConn][]string)
	}
	p.conns[key] = append(p.conns[key], cc)
	p.keys[cc] = append(p.keys[cc], key)
}

func (p *clientConnPool) MarkDead(cc *ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, key := range p.keys[cc] {
		vv, ok := p.conns[key]
		if !ok {
			continue
		}
		newList := filterOutClientConn(vv, cc)
		if len(newList) > 0 {
			p.conns[key] = newList
		} else {
			delete(p.conns, key)
		}
	}
	delete(p.keys, cc)
}

func (p *clientConnPool) closeIdleConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	
	
	
	
	
	for _, vv := range p.conns {
		for _, cc := range vv {
			cc.closeIfIdle()
		}
	}
}

func filterOutClientConn(in []*ClientConn, exclude *ClientConn) []*ClientConn {
	out := in[:0]
	for _, v := range in {
		if v != exclude {
			out = append(out, v)
		}
	}
	
	
	if len(in) != len(out) {
		in[len(in)-1] = nil
	}
	return out
}




type noDialClientConnPool struct{ *clientConnPool }

func (p noDialClientConnPool) GetClientConn(req *http.Request, addr string) (*ClientConn, error) {
	return p.getClientConn(req, addr, noDialOnMiss)
}
