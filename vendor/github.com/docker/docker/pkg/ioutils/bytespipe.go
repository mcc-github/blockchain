package ioutils 

import (
	"errors"
	"io"
	"sync"
)


const maxCap = 1e6


const minCap = 64



const blockThreshold = 1e6

var (
	
	ErrClosed = errors.New("write to closed BytesPipe")

	bufPools     = make(map[int]*sync.Pool)
	bufPoolsLock sync.Mutex
)





type BytesPipe struct {
	mu       sync.Mutex
	wait     *sync.Cond
	buf      []*fixedBuffer
	bufLen   int
	closeErr error 
}




func NewBytesPipe() *BytesPipe {
	bp := &BytesPipe{}
	bp.buf = append(bp.buf, getBuffer(minCap))
	bp.wait = sync.NewCond(&bp.mu)
	return bp
}



func (bp *BytesPipe) Write(p []byte) (int, error) {
	bp.mu.Lock()

	written := 0
loop0:
	for {
		if bp.closeErr != nil {
			bp.mu.Unlock()
			return written, ErrClosed
		}

		if len(bp.buf) == 0 {
			bp.buf = append(bp.buf, getBuffer(64))
		}
		
		b := bp.buf[len(bp.buf)-1]

		n, err := b.Write(p)
		written += n
		bp.bufLen += n

		
		if err != nil && err != errBufferFull {
			bp.wait.Broadcast()
			bp.mu.Unlock()
			return written, err
		}

		
		if len(p) == n {
			break
		}

		
		p = p[n:]

		
		for bp.bufLen >= blockThreshold {
			bp.wait.Wait()
			if bp.closeErr != nil {
				continue loop0
			}
		}

		
		nextCap := b.Cap() * 2
		if nextCap > maxCap {
			nextCap = maxCap
		}
		bp.buf = append(bp.buf, getBuffer(nextCap))
	}
	bp.wait.Broadcast()
	bp.mu.Unlock()
	return written, nil
}


func (bp *BytesPipe) CloseWithError(err error) error {
	bp.mu.Lock()
	if err != nil {
		bp.closeErr = err
	} else {
		bp.closeErr = io.EOF
	}
	bp.wait.Broadcast()
	bp.mu.Unlock()
	return nil
}


func (bp *BytesPipe) Close() error {
	return bp.CloseWithError(nil)
}



func (bp *BytesPipe) Read(p []byte) (n int, err error) {
	bp.mu.Lock()
	if bp.bufLen == 0 {
		if bp.closeErr != nil {
			bp.mu.Unlock()
			return 0, bp.closeErr
		}
		bp.wait.Wait()
		if bp.bufLen == 0 && bp.closeErr != nil {
			err := bp.closeErr
			bp.mu.Unlock()
			return 0, err
		}
	}

	for bp.bufLen > 0 {
		b := bp.buf[0]
		read, _ := b.Read(p) 
		n += read
		bp.bufLen -= read

		if b.Len() == 0 {
			
			returnBuffer(b)
			bp.buf[0] = nil
			bp.buf = bp.buf[1:]
		}

		if len(p) == read {
			break
		}

		p = p[read:]
	}

	bp.wait.Broadcast()
	bp.mu.Unlock()
	return
}

func returnBuffer(b *fixedBuffer) {
	b.Reset()
	bufPoolsLock.Lock()
	pool := bufPools[b.Cap()]
	bufPoolsLock.Unlock()
	if pool != nil {
		pool.Put(b)
	}
}

func getBuffer(size int) *fixedBuffer {
	bufPoolsLock.Lock()
	pool, ok := bufPools[size]
	if !ok {
		pool = &sync.Pool{New: func() interface{} { return &fixedBuffer{buf: make([]byte, 0, size)} }}
		bufPools[size] = pool
	}
	bufPoolsLock.Unlock()
	return pool.Get().(*fixedBuffer)
}
