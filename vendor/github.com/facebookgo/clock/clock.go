package clock

import (
	"runtime"
	"sort"
	"sync"
	"time"
)






type Clock interface {
	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) *Timer
	Now() time.Time
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	Ticker(d time.Duration) *Ticker
	Timer(d time.Duration) *Timer
}


func New() Clock {
	return &clock{}
}


type clock struct{}

func (c *clock) After(d time.Duration) <-chan time.Time { return time.After(d) }

func (c *clock) AfterFunc(d time.Duration, f func()) *Timer {
	return &Timer{timer: time.AfterFunc(d, f)}
}

func (c *clock) Now() time.Time { return time.Now() }

func (c *clock) Sleep(d time.Duration) { time.Sleep(d) }

func (c *clock) Tick(d time.Duration) <-chan time.Time { return time.Tick(d) }

func (c *clock) Ticker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{C: t.C, ticker: t}
}

func (c *clock) Timer(d time.Duration) *Timer {
	t := time.NewTimer(d)
	return &Timer{C: t.C, timer: t}
}



type Mock struct {
	mu     sync.Mutex
	now    time.Time   
	timers clockTimers 

	calls      Calls
	waiting    []waiting
	callsMutex sync.Mutex
}



func NewMock() *Mock {
	return &Mock{now: time.Unix(0, 0)}
}



func (m *Mock) Add(d time.Duration) {
	
	t := m.now.Add(d)

	
	for {
		if !m.runNextTimer(t) {
			break
		}
	}

	
	m.mu.Lock()
	m.now = t
	m.mu.Unlock()

	
	gosched()
}




func (m *Mock) runNextTimer(max time.Time) bool {
	m.mu.Lock()

	
	sort.Sort(m.timers)

	
	if len(m.timers) == 0 {
		m.mu.Unlock()
		return false
	}

	
	t := m.timers[0]
	if t.Next().After(max) {
		m.mu.Unlock()
		return false
	}

	
	m.now = t.Next()
	m.mu.Unlock()

	
	t.Tick(m.now)
	return true
}


func (m *Mock) After(d time.Duration) <-chan time.Time {
	defer m.inc(&m.calls.After)
	return m.Timer(d).C
}



func (m *Mock) AfterFunc(d time.Duration, f func()) *Timer {
	defer m.inc(&m.calls.AfterFunc)
	t := m.Timer(d)
	t.C = nil
	t.fn = f
	return t
}


func (m *Mock) Now() time.Time {
	defer m.inc(&m.calls.Now)
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.now
}



func (m *Mock) Sleep(d time.Duration) {
	defer m.inc(&m.calls.Sleep)
	<-m.After(d)
}



func (m *Mock) Tick(d time.Duration) <-chan time.Time {
	defer m.inc(&m.calls.Tick)
	return m.Ticker(d).C
}


func (m *Mock) Ticker(d time.Duration) *Ticker {
	defer m.inc(&m.calls.Ticker)
	m.mu.Lock()
	defer m.mu.Unlock()
	ch := make(chan time.Time)
	t := &Ticker{
		C:    ch,
		c:    ch,
		mock: m,
		d:    d,
		next: m.now.Add(d),
	}
	m.timers = append(m.timers, (*internalTicker)(t))
	return t
}


func (m *Mock) Timer(d time.Duration) *Timer {
	defer m.inc(&m.calls.Timer)
	m.mu.Lock()
	defer m.mu.Unlock()
	ch := make(chan time.Time)
	t := &Timer{
		C:    ch,
		c:    ch,
		mock: m,
		next: m.now.Add(d),
	}
	m.timers = append(m.timers, (*internalTimer)(t))
	return t
}

func (m *Mock) removeClockTimer(t clockTimer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, timer := range m.timers {
		if timer == t {
			copy(m.timers[i:], m.timers[i+1:])
			m.timers[len(m.timers)-1] = nil
			m.timers = m.timers[:len(m.timers)-1]
			break
		}
	}
	sort.Sort(m.timers)
}

func (m *Mock) inc(addr *uint32) {
	m.callsMutex.Lock()
	defer m.callsMutex.Unlock()
	*addr++
	var newWaiting []waiting
	for _, w := range m.waiting {
		if m.calls.atLeast(w.expected) {
			close(w.done)
			continue
		}
		newWaiting = append(newWaiting, w)
	}
	m.waiting = newWaiting
}





func (m *Mock) Wait(s Calls) {
	m.callsMutex.Lock()
	if m.calls.atLeast(s) {
		m.callsMutex.Unlock()
		return
	}
	done := make(chan struct{})
	m.waiting = append(m.waiting, waiting{expected: s, done: done})
	m.callsMutex.Unlock()
	<-done
}


type clockTimer interface {
	Next() time.Time
	Tick(time.Time)
}


type clockTimers []clockTimer

func (a clockTimers) Len() int           { return len(a) }
func (a clockTimers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a clockTimers) Less(i, j int) bool { return a[i].Next().Before(a[j].Next()) }



type Timer struct {
	C     <-chan time.Time
	c     chan time.Time
	timer *time.Timer 
	next  time.Time   
	mock  *Mock       
	fn    func()      
}


func (t *Timer) Stop() {
	if t.timer != nil {
		t.timer.Stop()
	} else {
		t.mock.removeClockTimer((*internalTimer)(t))
	}
}

type internalTimer Timer

func (t *internalTimer) Next() time.Time { return t.next }
func (t *internalTimer) Tick(now time.Time) {
	if t.fn != nil {
		t.fn()
	} else {
		t.c <- now
	}
	t.mock.removeClockTimer((*internalTimer)(t))
	gosched()
}


type Ticker struct {
	C      <-chan time.Time
	c      chan time.Time
	ticker *time.Ticker  
	next   time.Time     
	mock   *Mock         
	d      time.Duration 
}


func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
	} else {
		t.mock.removeClockTimer((*internalTicker)(t))
	}
}

type internalTicker Ticker

func (t *internalTicker) Next() time.Time { return t.next }
func (t *internalTicker) Tick(now time.Time) {
	select {
	case t.c <- now:
	case <-time.After(1 * time.Millisecond):
	}
	t.next = now.Add(t.d)
	gosched()
}


func gosched() { runtime.Gosched() }



type Calls struct {
	After     uint32
	AfterFunc uint32
	Now       uint32
	Sleep     uint32
	Tick      uint32
	Ticker    uint32
	Timer     uint32
}


func (c Calls) atLeast(o Calls) bool {
	if c.After < o.After {
		return false
	}
	if c.AfterFunc < o.AfterFunc {
		return false
	}
	if c.Now < o.Now {
		return false
	}
	if c.Sleep < o.Sleep {
		return false
	}
	if c.Tick < o.Tick {
		return false
	}
	if c.Ticker < o.Ticker {
		return false
	}
	if c.Timer < o.Timer {
		return false
	}
	return true
}

type waiting struct {
	expected Calls
	done     chan struct{}
}
