package metrics

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)



type Meter interface {
	Count() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	Snapshot() Meter
	Stop()
}





func GetOrRegisterMeter(name string, r Registry) Meter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewMeter).(Meter)
}



func NewMeter() Meter {
	if UseNilMetrics {
		return NilMeter{}
	}
	m := newStandardMeter()
	arbiter.Lock()
	defer arbiter.Unlock()
	arbiter.meters[m] = struct{}{}
	if !arbiter.started {
		arbiter.started = true
		go arbiter.tick()
	}
	return m
}





func NewRegisteredMeter(name string, r Registry) Meter {
	c := NewMeter()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}


type MeterSnapshot struct {
	count                          int64
	rate1, rate5, rate15, rateMean uint64
}


func (m *MeterSnapshot) Count() int64 { return m.count }


func (*MeterSnapshot) Mark(n int64) {
	panic("Mark called on a MeterSnapshot")
}



func (m *MeterSnapshot) Rate1() float64 { return math.Float64frombits(m.rate1) }



func (m *MeterSnapshot) Rate5() float64 { return math.Float64frombits(m.rate5) }



func (m *MeterSnapshot) Rate15() float64 { return math.Float64frombits(m.rate15) }



func (m *MeterSnapshot) RateMean() float64 { return math.Float64frombits(m.rateMean) }


func (m *MeterSnapshot) Snapshot() Meter { return m }


func (m *MeterSnapshot) Stop() {}


type NilMeter struct{}


func (NilMeter) Count() int64 { return 0 }


func (NilMeter) Mark(n int64) {}


func (NilMeter) Rate1() float64 { return 0.0 }


func (NilMeter) Rate5() float64 { return 0.0 }


func (NilMeter) Rate15() float64 { return 0.0 }


func (NilMeter) RateMean() float64 { return 0.0 }


func (NilMeter) Snapshot() Meter { return NilMeter{} }


func (NilMeter) Stop() {}


type StandardMeter struct {
	snapshot    *MeterSnapshot
	a1, a5, a15 EWMA
	startTime   time.Time
	stopped     uint32
}

func newStandardMeter() *StandardMeter {
	return &StandardMeter{
		snapshot:  &MeterSnapshot{},
		a1:        NewEWMA1(),
		a5:        NewEWMA5(),
		a15:       NewEWMA15(),
		startTime: time.Now(),
	}
}


func (m *StandardMeter) Stop() {
	if atomic.CompareAndSwapUint32(&m.stopped, 0, 1) {
		arbiter.Lock()
		delete(arbiter.meters, m)
		arbiter.Unlock()
	}
}


func (m *StandardMeter) Count() int64 {
	return atomic.LoadInt64(&m.snapshot.count)
}


func (m *StandardMeter) Mark(n int64) {
	if atomic.LoadUint32(&m.stopped) == 1 {
		return
	}

	atomic.AddInt64(&m.snapshot.count, n)

	m.a1.Update(n)
	m.a5.Update(n)
	m.a15.Update(n)
	m.updateSnapshot()
}


func (m *StandardMeter) Rate1() float64 {
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rate1))
}


func (m *StandardMeter) Rate5() float64 {
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rate5))
}


func (m *StandardMeter) Rate15() float64 {
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rate15))
}


func (m *StandardMeter) RateMean() float64 {
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rateMean))
}


func (m *StandardMeter) Snapshot() Meter {
	copiedSnapshot := MeterSnapshot{
		count:    atomic.LoadInt64(&m.snapshot.count),
		rate1:    atomic.LoadUint64(&m.snapshot.rate1),
		rate5:    atomic.LoadUint64(&m.snapshot.rate5),
		rate15:   atomic.LoadUint64(&m.snapshot.rate15),
		rateMean: atomic.LoadUint64(&m.snapshot.rateMean),
	}
	return &copiedSnapshot
}

func (m *StandardMeter) updateSnapshot() {
	rate1 := math.Float64bits(m.a1.Rate())
	rate5 := math.Float64bits(m.a5.Rate())
	rate15 := math.Float64bits(m.a15.Rate())
	rateMean := math.Float64bits(float64(m.Count()) / time.Since(m.startTime).Seconds())

	atomic.StoreUint64(&m.snapshot.rate1, rate1)
	atomic.StoreUint64(&m.snapshot.rate5, rate5)
	atomic.StoreUint64(&m.snapshot.rate15, rate15)
	atomic.StoreUint64(&m.snapshot.rateMean, rateMean)
}

func (m *StandardMeter) tick() {
	m.a1.Tick()
	m.a5.Tick()
	m.a15.Tick()
	m.updateSnapshot()
}



type meterArbiter struct {
	sync.RWMutex
	started bool
	meters  map[*StandardMeter]struct{}
	ticker  *time.Ticker
}

var arbiter = meterArbiter{ticker: time.NewTicker(5e9), meters: make(map[*StandardMeter]struct{})}


func (ma *meterArbiter) tick() {
	for {
		select {
		case <-ma.ticker.C:
			ma.tickMeters()
		}
	}
}

func (ma *meterArbiter) tickMeters() {
	ma.RLock()
	defer ma.RUnlock()
	for meter := range ma.meters {
		meter.tick()
	}
}
