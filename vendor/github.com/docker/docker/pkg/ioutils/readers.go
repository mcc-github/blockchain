package ioutils 

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"golang.org/x/net/context"
)




type ReadCloserWrapper struct {
	io.Reader
	closer func() error
}


func (r *ReadCloserWrapper) Close() error {
	return r.closer()
}


func NewReadCloserWrapper(r io.Reader, closer func() error) io.ReadCloser {
	return &ReadCloserWrapper{
		Reader: r,
		closer: closer,
	}
}

type readerErrWrapper struct {
	reader io.Reader
	closer func()
}

func (r *readerErrWrapper) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if err != nil {
		r.closer()
	}
	return n, err
}


func NewReaderErrWrapper(r io.Reader, closer func()) io.Reader {
	return &readerErrWrapper{
		reader: r,
		closer: closer,
	}
}


func HashData(src io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, src); err != nil {
		return "", err
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}



type OnEOFReader struct {
	Rc io.ReadCloser
	Fn func()
}

func (r *OnEOFReader) Read(p []byte) (n int, err error) {
	n, err = r.Rc.Read(p)
	if err == io.EOF {
		r.runFunc()
	}
	return
}


func (r *OnEOFReader) Close() error {
	err := r.Rc.Close()
	r.runFunc()
	return err
}

func (r *OnEOFReader) runFunc() {
	if fn := r.Fn; fn != nil {
		fn()
		r.Fn = nil
	}
}



type cancelReadCloser struct {
	cancel func()
	pR     *io.PipeReader 
	pW     *io.PipeWriter
}




func NewCancelReadCloser(ctx context.Context, in io.ReadCloser) io.ReadCloser {
	pR, pW := io.Pipe()

	
	doneCtx, cancel := context.WithCancel(context.Background())

	p := &cancelReadCloser{
		cancel: cancel,
		pR:     pR,
		pW:     pW,
	}

	go func() {
		_, err := io.Copy(pW, in)
		select {
		case <-ctx.Done():
			
			
			
		default:
			p.closeWithError(err)
		}
		in.Close()
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				p.closeWithError(ctx.Err())
			case <-doneCtx.Done():
				return
			}
		}
	}()

	return p
}



func (p *cancelReadCloser) Read(buf []byte) (n int, err error) {
	return p.pR.Read(buf)
}



func (p *cancelReadCloser) closeWithError(err error) {
	p.pW.CloseWithError(err)
	p.cancel()
}



func (p *cancelReadCloser) Close() error {
	p.closeWithError(io.EOF)
	return nil
}
