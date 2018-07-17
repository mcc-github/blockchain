




package semaphore 

import (
	"container/list"
	"sync"

	
	
	
	"golang.org/x/net/context"
)

type waiter struct {
	n     int64
	ready chan<- struct{} 
}



func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}



type Weighted struct {
	size    int64
	cur     int64
	mu      sync.Mutex
	waiters list.List
}






func (s *Weighted) Acquire(ctx context.Context, n int64) error {
	s.mu.Lock()
	if s.size-s.cur >= n && s.waiters.Len() == 0 {
		s.cur += n
		s.mu.Unlock()
		return nil
	}

	if n > s.size {
		
		s.mu.Unlock()
		<-ctx.Done()
		return ctx.Err()
	}

	ready := make(chan struct{})
	w := waiter{n: n, ready: ready}
	elem := s.waiters.PushBack(w)
	s.mu.Unlock()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		s.mu.Lock()
		select {
		case <-ready:
			
			
			err = nil
		default:
			s.waiters.Remove(elem)
		}
		s.mu.Unlock()
		return err

	case <-ready:
		return nil
	}
}



func (s *Weighted) TryAcquire(n int64) bool {
	s.mu.Lock()
	success := s.size-s.cur >= n && s.waiters.Len() == 0
	if success {
		s.cur += n
	}
	s.mu.Unlock()
	return success
}


func (s *Weighted) Release(n int64) {
	s.mu.Lock()
	s.cur -= n
	if s.cur < 0 {
		s.mu.Unlock()
		panic("semaphore: bad release")
	}
	for {
		next := s.waiters.Front()
		if next == nil {
			break 
		}

		w := next.Value.(waiter)
		if s.size-s.cur < w.n {
			
			
			
			
			
			
			
			
			
			
			
			break
		}

		s.cur += w.n
		s.waiters.Remove(next)
		close(w.ready)
	}
	s.mu.Unlock()
}
