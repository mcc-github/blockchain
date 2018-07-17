



package norm

import "io"

type normWriter struct {
	rb  reorderBuffer
	w   io.Writer
	buf []byte
}




func (w *normWriter) Write(data []byte) (n int, err error) {
	
	const chunk = 4000

	for len(data) > 0 {
		
		m := len(data)
		if m > chunk {
			m = chunk
		}
		w.rb.src = inputBytes(data[:m])
		w.rb.nsrc = m
		w.buf = doAppend(&w.rb, w.buf, 0)
		data = data[m:]
		n += m

		
		
		i := lastBoundary(&w.rb.f, w.buf)
		if i == -1 {
			i = 0
		}
		if i > 0 {
			if _, err = w.w.Write(w.buf[:i]); err != nil {
				break
			}
			bn := copy(w.buf, w.buf[i:])
			w.buf = w.buf[:bn]
		}
	}
	return n, err
}


func (w *normWriter) Close() error {
	if len(w.buf) > 0 {
		_, err := w.w.Write(w.buf)
		if err != nil {
			return err
		}
	}
	return nil
}





func (f Form) Writer(w io.Writer) io.WriteCloser {
	wr := &normWriter{rb: reorderBuffer{}, w: w}
	wr.rb.init(f, nil)
	return wr
}

type normReader struct {
	rb           reorderBuffer
	r            io.Reader
	inbuf        []byte
	outbuf       []byte
	bufStart     int
	lastBoundary int
	err          error
}


func (r *normReader) Read(p []byte) (int, error) {
	for {
		if r.lastBoundary-r.bufStart > 0 {
			n := copy(p, r.outbuf[r.bufStart:r.lastBoundary])
			r.bufStart += n
			if r.lastBoundary-r.bufStart > 0 {
				return n, nil
			}
			return n, r.err
		}
		if r.err != nil {
			return 0, r.err
		}
		outn := copy(r.outbuf, r.outbuf[r.lastBoundary:])
		r.outbuf = r.outbuf[0:outn]
		r.bufStart = 0

		n, err := r.r.Read(r.inbuf)
		r.rb.src = inputBytes(r.inbuf[0:n])
		r.rb.nsrc, r.err = n, err
		if n > 0 {
			r.outbuf = doAppend(&r.rb, r.outbuf, 0)
		}
		if err == io.EOF {
			r.lastBoundary = len(r.outbuf)
		} else {
			r.lastBoundary = lastBoundary(&r.rb.f, r.outbuf)
			if r.lastBoundary == -1 {
				r.lastBoundary = 0
			}
		}
	}
}



func (f Form) Reader(r io.Reader) io.Reader {
	const chunk = 4000
	buf := make([]byte, chunk)
	rr := &normReader{rb: reorderBuffer{}, r: r, inbuf: buf}
	rr.rb.init(f, buf)
	return rr
}
