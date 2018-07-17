

package pb

import (
	"io"
	"sync"
	"time"
)



func StartPool(pbs ...*ProgressBar) (pool *Pool, err error) {
	pool = new(Pool)
	if err = pool.Start(); err != nil {
		return
	}
	pool.Add(pbs...)
	return
}



func NewPool(pbs ...*ProgressBar) (pool *Pool) {
	pool = new(Pool)
	pool.Add(pbs...)
	return
}

type Pool struct {
	Output        io.Writer
	RefreshRate   time.Duration
	bars          []*ProgressBar
	lastBarsCount int
	shutdownCh    chan struct{}
	workerCh      chan struct{}
	m             sync.Mutex
	finishOnce    sync.Once
}


func (p *Pool) Add(pbs ...*ProgressBar) {
	p.m.Lock()
	defer p.m.Unlock()
	for _, bar := range pbs {
		bar.ManualUpdate = true
		bar.NotPrint = true
		bar.Start()
		p.bars = append(p.bars, bar)
	}
}

func (p *Pool) Start() (err error) {
	p.RefreshRate = DefaultRefreshRate
	p.shutdownCh, err = lockEcho()
	if err != nil {
		return
	}
	p.workerCh = make(chan struct{})
	go p.writer()
	return
}

func (p *Pool) writer() {
	var first = true
	defer func() {
		if first == false {
			p.print(false)
		} else {
			p.print(true)
			p.print(false)
		}
		close(p.workerCh)
	}()

	for {
		select {
		case <-time.After(p.RefreshRate):
			if p.print(first) {
				p.print(false)
				return
			}
			first = false
		case <-p.shutdownCh:
			return
		}
	}
}


func (p *Pool) Stop() error {
	p.finishOnce.Do(func() {
		close(p.shutdownCh)
	})

	
	select {
	case <-p.workerCh:
	}

	return unlockEcho()
}
