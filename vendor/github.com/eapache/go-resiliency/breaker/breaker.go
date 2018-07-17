
package breaker

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)



var ErrBreakerOpen = errors.New("circuit breaker is open")

const (
	closed uint32 = iota
	open
	halfOpen
)


type Breaker struct {
	errorThreshold, successThreshold int
	timeout                          time.Duration

	lock              sync.Mutex
	state             uint32
	errors, successes int
	lastError         time.Time
}






func New(errorThreshold, successThreshold int, timeout time.Duration) *Breaker {
	return &Breaker{
		errorThreshold:   errorThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}




func (b *Breaker) Run(work func() error) error {
	state := atomic.LoadUint32(&b.state)

	if state == open {
		return ErrBreakerOpen
	}

	return b.doWork(state, work)
}






func (b *Breaker) Go(work func() error) error {
	state := atomic.LoadUint32(&b.state)

	if state == open {
		return ErrBreakerOpen
	}

	
	
	
	go b.doWork(state, work)

	return nil
}

func (b *Breaker) doWork(state uint32, work func() error) error {
	var panicValue interface{}

	result := func() error {
		defer func() {
			panicValue = recover()
		}()
		return work()
	}()

	if result == nil && panicValue == nil && state == closed {
		
		
		return nil
	}

	
	b.processResult(result, panicValue)

	if panicValue != nil {
		
		
		panic(panicValue)
	}

	return result
}

func (b *Breaker) processResult(result error, panicValue interface{}) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if result == nil && panicValue == nil {
		if b.state == halfOpen {
			b.successes++
			if b.successes == b.successThreshold {
				b.closeBreaker()
			}
		}
	} else {
		if b.errors > 0 {
			expiry := b.lastError.Add(b.timeout)
			if time.Now().After(expiry) {
				b.errors = 0
			}
		}

		switch b.state {
		case closed:
			b.errors++
			if b.errors == b.errorThreshold {
				b.openBreaker()
			} else {
				b.lastError = time.Now()
			}
		case halfOpen:
			b.openBreaker()
		}
	}
}

func (b *Breaker) openBreaker() {
	b.changeState(open)
	go b.timer()
}

func (b *Breaker) closeBreaker() {
	b.changeState(closed)
}

func (b *Breaker) timer() {
	time.Sleep(b.timeout)

	b.lock.Lock()
	defer b.lock.Unlock()

	b.changeState(halfOpen)
}

func (b *Breaker) changeState(newState uint32) {
	b.errors = 0
	b.successes = 0
	atomic.StoreUint32(&b.state, newState)
}
