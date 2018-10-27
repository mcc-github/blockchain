package metrics

import (
	"math"
	"sync"
	"sync/atomic"
)



type EWMA interface {
	Rate() float64
	Snapshot() EWMA
	Tick()
	Update(int64)
}


func NewEWMA(alpha float64) EWMA {
	if UseNilMetrics {
		return NilEWMA{}
	}
	return &StandardEWMA{alpha: alpha}
}


func NewEWMA1() EWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/1))
}


func NewEWMA5() EWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/5))
}


func NewEWMA15() EWMA {
	return NewEWMA(1 - math.Exp(-5.0/60.0/15))
}


type EWMASnapshot float64



func (a EWMASnapshot) Rate() float64 { return float64(a) }


func (a EWMASnapshot) Snapshot() EWMA { return a }


func (EWMASnapshot) Tick() {
	panic("Tick called on an EWMASnapshot")
}


func (EWMASnapshot) Update(int64) {
	panic("Update called on an EWMASnapshot")
}


type NilEWMA struct{}


func (NilEWMA) Rate() float64 { return 0.0 }


func (NilEWMA) Snapshot() EWMA { return NilEWMA{} }


func (NilEWMA) Tick() {}


func (NilEWMA) Update(n int64) {}




type StandardEWMA struct {
	uncounted int64 
	alpha     float64
	rate      uint64
	init      uint32
	mutex     sync.Mutex
}


func (a *StandardEWMA) Rate() float64 {
	currentRate := math.Float64frombits(atomic.LoadUint64(&a.rate)) * float64(1e9)
	return currentRate
}


func (a *StandardEWMA) Snapshot() EWMA {
	return EWMASnapshot(a.Rate())
}



func (a *StandardEWMA) Tick() {
	
	if atomic.LoadUint32(&a.init) == 1 {
		a.updateRate(a.fetchInstantRate())
	} else {
		
		
		
		
		a.mutex.Lock()
		if atomic.LoadUint32(&a.init) == 1 {
			
			
			a.updateRate(a.fetchInstantRate())
		} else {
			atomic.StoreUint32(&a.init, 1)
			atomic.StoreUint64(&a.rate, math.Float64bits(a.fetchInstantRate()))
		}
		a.mutex.Unlock()
	}
}

func (a *StandardEWMA) fetchInstantRate() float64 {
	count := atomic.LoadInt64(&a.uncounted)
	atomic.AddInt64(&a.uncounted, -count)
	instantRate := float64(count) / float64(5e9)
	return instantRate
}

func (a *StandardEWMA) updateRate(instantRate float64) {
	currentRate := math.Float64frombits(atomic.LoadUint64(&a.rate))
	currentRate += a.alpha * (instantRate - currentRate)
	atomic.StoreUint64(&a.rate, math.Float64bits(currentRate))
}


func (a *StandardEWMA) Update(n int64) {
	atomic.AddInt64(&a.uncounted, n)
}
