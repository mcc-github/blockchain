

package transport

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const (
	
	defaultWindowSize = 65535
	
	initialWindowSize             = defaultWindowSize 
	infinity                      = time.Duration(math.MaxInt64)
	defaultClientKeepaliveTime    = infinity
	defaultClientKeepaliveTimeout = 20 * time.Second
	defaultMaxStreamsClient       = 100
	defaultMaxConnectionIdle      = infinity
	defaultMaxConnectionAge       = infinity
	defaultMaxConnectionAgeGrace  = infinity
	defaultServerKeepaliveTime    = 2 * time.Hour
	defaultServerKeepaliveTimeout = 20 * time.Second
	defaultKeepalivePolicyMinTime = 5 * time.Minute
	
	maxWindowSize = math.MaxInt32
	
	
	
	defaultWriteQuota = 64 * 1024
)



type writeQuota struct {
	quota int32
	
	
	ch chan struct{}
	
	done <-chan struct{}
	
	
	
	replenish func(n int)
}

func newWriteQuota(sz int32, done <-chan struct{}) *writeQuota {
	w := &writeQuota{
		quota: sz,
		ch:    make(chan struct{}, 1),
		done:  done,
	}
	w.replenish = w.realReplenish
	return w
}

func (w *writeQuota) get(sz int32) error {
	for {
		if atomic.LoadInt32(&w.quota) > 0 {
			atomic.AddInt32(&w.quota, -sz)
			return nil
		}
		select {
		case <-w.ch:
			continue
		case <-w.done:
			return errStreamDone
		}
	}
}

func (w *writeQuota) realReplenish(n int) {
	sz := int32(n)
	a := atomic.AddInt32(&w.quota, sz)
	b := a - sz
	if b <= 0 && a > 0 {
		select {
		case w.ch <- struct{}{}:
		default:
		}
	}
}

type trInFlow struct {
	limit               uint32
	unacked             uint32
	effectiveWindowSize uint32
}

func (f *trInFlow) newLimit(n uint32) uint32 {
	d := n - f.limit
	f.limit = n
	f.updateEffectiveWindowSize()
	return d
}

func (f *trInFlow) onData(n uint32) uint32 {
	f.unacked += n
	if f.unacked >= f.limit/4 {
		w := f.unacked
		f.unacked = 0
		f.updateEffectiveWindowSize()
		return w
	}
	f.updateEffectiveWindowSize()
	return 0
}

func (f *trInFlow) reset() uint32 {
	w := f.unacked
	f.unacked = 0
	f.updateEffectiveWindowSize()
	return w
}

func (f *trInFlow) updateEffectiveWindowSize() {
	atomic.StoreUint32(&f.effectiveWindowSize, f.limit-f.unacked)
}

func (f *trInFlow) getSize() uint32 {
	return atomic.LoadUint32(&f.effectiveWindowSize)
}



type inFlow struct {
	mu sync.Mutex
	
	limit uint32
	
	
	pendingData uint32
	
	
	pendingUpdate uint32
	
	
	delta uint32
}



func (f *inFlow) newLimit(n uint32) uint32 {
	f.mu.Lock()
	d := n - f.limit
	f.limit = n
	f.mu.Unlock()
	return d
}

func (f *inFlow) maybeAdjust(n uint32) uint32 {
	if n > uint32(math.MaxInt32) {
		n = uint32(math.MaxInt32)
	}
	f.mu.Lock()
	
	
	estSenderQuota := int32(f.limit - (f.pendingData + f.pendingUpdate))
	
	
	
	estUntransmittedData := int32(n - f.pendingData) 
	
	
	
	if estUntransmittedData > estSenderQuota {
		
		if f.limit+n > maxWindowSize {
			f.delta = maxWindowSize - f.limit
		} else {
			
			
			
			f.delta = n
		}
		f.mu.Unlock()
		return f.delta
	}
	f.mu.Unlock()
	return 0
}


func (f *inFlow) onData(n uint32) error {
	f.mu.Lock()
	f.pendingData += n
	if f.pendingData+f.pendingUpdate > f.limit+f.delta {
		limit := f.limit
		rcvd := f.pendingData + f.pendingUpdate
		f.mu.Unlock()
		return fmt.Errorf("received %d-bytes data exceeding the limit %d bytes", rcvd, limit)
	}
	f.mu.Unlock()
	return nil
}



func (f *inFlow) onRead(n uint32) uint32 {
	f.mu.Lock()
	if f.pendingData == 0 {
		f.mu.Unlock()
		return 0
	}
	f.pendingData -= n
	if n > f.delta {
		n -= f.delta
		f.delta = 0
	} else {
		f.delta -= n
		n = 0
	}
	f.pendingUpdate += n
	if f.pendingUpdate >= f.limit/4 {
		wu := f.pendingUpdate
		f.pendingUpdate = 0
		f.mu.Unlock()
		return wu
	}
	f.mu.Unlock()
	return 0
}
