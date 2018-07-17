



package statsd

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

var senderPool = newBufferPool()



type BufferedSender struct {
	sender        Sender
	flushBytes    int
	flushInterval time.Duration
	
	bufmx  sync.Mutex
	buffer *bytes.Buffer
	bufs   chan *bytes.Buffer
	
	runmx    sync.RWMutex
	shutdown chan chan error
	running  bool
}


func (s *BufferedSender) Send(data []byte) (int, error) {
	s.runmx.RLock()
	if !s.running {
		s.runmx.RUnlock()
		return 0, fmt.Errorf("BufferedSender is not running")
	}

	s.withBufferLock(func() {
		blen := s.buffer.Len()
		if blen > 0 && blen+len(data)+1 >= s.flushBytes {
			s.swapnqueue()
		}

		s.buffer.Write(data)
		s.buffer.WriteByte('\n')

		if s.buffer.Len() >= s.flushBytes {
			s.swapnqueue()
		}
	})
	s.runmx.RUnlock()
	return len(data), nil
}


func (s *BufferedSender) Close() error {
	
	s.runmx.Lock()
	defer s.runmx.Unlock()
	if !s.running {
		return nil
	}

	errChan := make(chan error)
	s.running = false
	s.shutdown <- errChan
	return <-errChan
}



func (s *BufferedSender) Start() {
	
	s.runmx.Lock()
	defer s.runmx.Unlock()
	if s.running {
		return
	}

	s.running = true
	s.bufs = make(chan *bytes.Buffer, 32)
	go s.run()
}

func (s *BufferedSender) withBufferLock(fn func()) {
	s.bufmx.Lock()
	fn()
	s.bufmx.Unlock()
}

func (s *BufferedSender) swapnqueue() {
	if s.buffer.Len() == 0 {
		return
	}
	ob := s.buffer
	nb := senderPool.Get()
	s.buffer = nb
	s.bufs <- ob
}

func (s *BufferedSender) run() {
	ticker := time.NewTicker(s.flushInterval)
	defer ticker.Stop()

	doneChan := make(chan bool)
	go func() {
		for buf := range s.bufs {
			s.flush(buf)
			senderPool.Put(buf)
		}
		doneChan <- true
	}()

	for {
		select {
		case <-ticker.C:
			s.withBufferLock(func() {
				s.swapnqueue()
			})
		case errChan := <-s.shutdown:
			s.withBufferLock(func() {
				s.swapnqueue()
			})
			close(s.bufs)
			<-doneChan
			errChan <- s.sender.Close()
			return
		}
	}
}


func (s *BufferedSender) flush(b *bytes.Buffer) (int, error) {
	bb := b.Bytes()
	bbl := len(bb)
	if bb[bbl-1] == '\n' {
		bb = bb[:bbl-1]
	}
	
	n, err := s.sender.Send(bb)
	b.Truncate(0) 
	return n, err
}













func NewBufferedSender(addr string, flushInterval time.Duration, flushBytes int) (Sender, error) {
	simpleSender, err := NewSimpleSender(addr)
	if err != nil {
		return nil, err
	}

	sender := &BufferedSender{
		flushBytes:    flushBytes,
		flushInterval: flushInterval,
		sender:        simpleSender,
		buffer:        senderPool.Get(),
		shutdown:      make(chan chan error),
	}

	sender.Start()
	return sender, nil
}
