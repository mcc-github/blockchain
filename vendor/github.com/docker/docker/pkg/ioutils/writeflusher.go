package ioutils 

import (
	"io"
	"sync"
)




type WriteFlusher struct {
	w           io.Writer
	flusher     flusher
	flushed     chan struct{}
	flushedOnce sync.Once
	closed      chan struct{}
	closeLock   sync.Mutex
}

type flusher interface {
	Flush()
}

var errWriteFlusherClosed = io.EOF

func (wf *WriteFlusher) Write(b []byte) (n int, err error) {
	select {
	case <-wf.closed:
		return 0, errWriteFlusherClosed
	default:
	}

	n, err = wf.w.Write(b)
	wf.Flush() 
	return n, err
}


func (wf *WriteFlusher) Flush() {
	select {
	case <-wf.closed:
		return
	default:
	}

	wf.flushedOnce.Do(func() {
		close(wf.flushed)
	})
	wf.flusher.Flush()
}



func (wf *WriteFlusher) Flushed() bool {
	
	
	
	var flushed bool
	select {
	case <-wf.flushed:
		flushed = true
	default:
	}
	return flushed
}




func (wf *WriteFlusher) Close() error {
	wf.closeLock.Lock()
	defer wf.closeLock.Unlock()

	select {
	case <-wf.closed:
		return errWriteFlusherClosed
	default:
		close(wf.closed)
	}
	return nil
}


func NewWriteFlusher(w io.Writer) *WriteFlusher {
	var fl flusher
	if f, ok := w.(flusher); ok {
		fl = f
	} else {
		fl = &NopFlusher{}
	}
	return &WriteFlusher{w: w, flusher: fl, closed: make(chan struct{}), flushed: make(chan struct{})}
}
