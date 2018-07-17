








package pools 

import (
	"bufio"
	"io"
	"sync"

	"github.com/docker/docker/pkg/ioutils"
)

const buffer32K = 32 * 1024

var (
	
	BufioReader32KPool = newBufioReaderPoolWithSize(buffer32K)
	
	BufioWriter32KPool = newBufioWriterPoolWithSize(buffer32K)
	buffer32KPool      = newBufferPoolWithSize(buffer32K)
)


type BufioReaderPool struct {
	pool sync.Pool
}



func newBufioReaderPoolWithSize(size int) *BufioReaderPool {
	return &BufioReaderPool{
		pool: sync.Pool{
			New: func() interface{} { return bufio.NewReaderSize(nil, size) },
		},
	}
}


func (bufPool *BufioReaderPool) Get(r io.Reader) *bufio.Reader {
	buf := bufPool.pool.Get().(*bufio.Reader)
	buf.Reset(r)
	return buf
}


func (bufPool *BufioReaderPool) Put(b *bufio.Reader) {
	b.Reset(nil)
	bufPool.pool.Put(b)
}

type bufferPool struct {
	pool sync.Pool
}

func newBufferPoolWithSize(size int) *bufferPool {
	return &bufferPool{
		pool: sync.Pool{
			New: func() interface{} { return make([]byte, size) },
		},
	}
}

func (bp *bufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

func (bp *bufferPool) Put(b []byte) {
	bp.pool.Put(b)
}


func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := buffer32KPool.Get()
	written, err = io.CopyBuffer(dst, src, buf)
	buffer32KPool.Put(buf)
	return
}



func (bufPool *BufioReaderPool) NewReadCloserWrapper(buf *bufio.Reader, r io.Reader) io.ReadCloser {
	return ioutils.NewReadCloserWrapper(r, func() error {
		if readCloser, ok := r.(io.ReadCloser); ok {
			readCloser.Close()
		}
		bufPool.Put(buf)
		return nil
	})
}


type BufioWriterPool struct {
	pool sync.Pool
}



func newBufioWriterPoolWithSize(size int) *BufioWriterPool {
	return &BufioWriterPool{
		pool: sync.Pool{
			New: func() interface{} { return bufio.NewWriterSize(nil, size) },
		},
	}
}


func (bufPool *BufioWriterPool) Get(w io.Writer) *bufio.Writer {
	buf := bufPool.pool.Get().(*bufio.Writer)
	buf.Reset(w)
	return buf
}


func (bufPool *BufioWriterPool) Put(b *bufio.Writer) {
	b.Reset(nil)
	bufPool.pool.Put(b)
}



func (bufPool *BufioWriterPool) NewWriteCloserWrapper(buf *bufio.Writer, w io.Writer) io.WriteCloser {
	return ioutils.NewWriteCloserWrapper(w, func() error {
		buf.Flush()
		if writeCloser, ok := w.(io.WriteCloser); ok {
			writeCloser.Close()
		}
		bufPool.Put(buf)
		return nil
	})
}
