







package transform 

import (
	"bytes"
	"errors"
	"io"
	"unicode/utf8"
)

var (
	
	
	ErrShortDst = errors.New("transform: short destination buffer")

	
	
	ErrShortSrc = errors.New("transform: short source buffer")

	
	
	ErrEndOfSpan = errors.New("transform: input and output are not identical")

	
	
	errInconsistentByteCount = errors.New("transform: inconsistent byte count returned")

	
	
	errShortInternal = errors.New("transform: short internal buffer")
)


type Transformer interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error)

	
	Reset()
}



type SpanningTransformer interface {
	Transformer

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Span(src []byte, atEOF bool) (n int, err error)
}



type NopResetter struct{}


func (NopResetter) Reset() {}


type Reader struct {
	r   io.Reader
	t   Transformer
	err error

	
	
	dst        []byte
	dst0, dst1 int

	
	
	src        []byte
	src0, src1 int

	
	
	transformComplete bool
}

const defaultBufSize = 4096



func NewReader(r io.Reader, t Transformer) *Reader {
	t.Reset()
	return &Reader{
		r:   r,
		t:   t,
		dst: make([]byte, defaultBufSize),
		src: make([]byte, defaultBufSize),
	}
}


func (r *Reader) Read(p []byte) (int, error) {
	n, err := 0, error(nil)
	for {
		
		if r.dst0 != r.dst1 {
			n = copy(p, r.dst[r.dst0:r.dst1])
			r.dst0 += n
			if r.dst0 == r.dst1 && r.transformComplete {
				return n, r.err
			}
			return n, nil
		} else if r.transformComplete {
			return 0, r.err
		}

		
		
		
		
		if r.src0 != r.src1 || r.err != nil {
			r.dst0 = 0
			r.dst1, n, err = r.t.Transform(r.dst, r.src[r.src0:r.src1], r.err == io.EOF)
			r.src0 += n

			switch {
			case err == nil:
				if r.src0 != r.src1 {
					r.err = errInconsistentByteCount
				}
				
				
				r.transformComplete = r.err != nil
				continue
			case err == ErrShortDst && (r.dst1 != 0 || n != 0):
				
				continue
			case err == ErrShortSrc && r.src1-r.src0 != len(r.src) && r.err == nil:
				
			default:
				r.transformComplete = true
				
				
				if r.err == nil || r.err == io.EOF {
					r.err = err
				}
				continue
			}
		}

		
		
		if r.src0 != 0 {
			r.src0, r.src1 = 0, copy(r.src, r.src[r.src0:r.src1])
		}
		n, r.err = r.r.Read(r.src[r.src1:])
		r.src1 += n
	}
}






type Writer struct {
	w   io.Writer
	t   Transformer
	dst []byte

	
	src []byte
	n   int
}



func NewWriter(w io.Writer, t Transformer) *Writer {
	t.Reset()
	return &Writer{
		w:   w,
		t:   t,
		dst: make([]byte, defaultBufSize),
		src: make([]byte, defaultBufSize),
	}
}




func (w *Writer) Write(data []byte) (n int, err error) {
	src := data
	if w.n > 0 {
		
		
		n = copy(w.src[w.n:], data)
		w.n += n
		src = w.src[:w.n]
	}
	for {
		nDst, nSrc, err := w.t.Transform(w.dst, src, false)
		if _, werr := w.w.Write(w.dst[:nDst]); werr != nil {
			return n, werr
		}
		src = src[nSrc:]
		if w.n == 0 {
			n += nSrc
		} else if len(src) <= n {
			
			
			w.n = 0
			n -= len(src)
			src = data[n:]
			if n < len(data) && (err == nil || err == ErrShortSrc) {
				continue
			}
		}
		switch err {
		case ErrShortDst:
			
			if nDst > 0 || nSrc > 0 {
				continue
			}
		case ErrShortSrc:
			if len(src) < len(w.src) {
				m := copy(w.src, src)
				
				
				if w.n == 0 {
					n += m
				}
				w.n = m
				err = nil
			} else if nDst > 0 || nSrc > 0 {
				
				
				
				
				
				
				continue
			}
		case nil:
			if w.n > 0 {
				err = errInconsistentByteCount
			}
		}
		return n, err
	}
}


func (w *Writer) Close() error {
	src := w.src[:w.n]
	for {
		nDst, nSrc, err := w.t.Transform(w.dst, src, true)
		if _, werr := w.w.Write(w.dst[:nDst]); werr != nil {
			return werr
		}
		if err != ErrShortDst {
			return err
		}
		src = src[nSrc:]
	}
}

type nop struct{ NopResetter }

func (nop) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	n := copy(dst, src)
	if n < len(src) {
		err = ErrShortDst
	}
	return n, n, err
}

func (nop) Span(src []byte, atEOF bool) (n int, err error) {
	return len(src), nil
}

type discard struct{ NopResetter }

func (discard) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	return 0, len(src), nil
}

var (
	
	
	Discard Transformer = discard{}

	
	Nop SpanningTransformer = nop{}
)







type chain struct {
	link []link
	err  error
	
	
	
	errStart int
}

func (c *chain) fatalError(errIndex int, err error) {
	if i := errIndex + 1; i > c.errStart {
		c.errStart = i
		c.err = err
	}
}

type link struct {
	t Transformer
	
	b []byte
	p int
	n int
}

func (l *link) src() []byte {
	return l.b[l.p:l.n]
}

func (l *link) dst() []byte {
	return l.b[l.n:]
}


func Chain(t ...Transformer) Transformer {
	if len(t) == 0 {
		return nop{}
	}
	c := &chain{link: make([]link, len(t)+1)}
	for i, tt := range t {
		c.link[i].t = tt
	}
	
	b := make([][defaultBufSize]byte, len(t)-1)
	for i := range b {
		c.link[i+1].b = b[i][:]
	}
	return c
}


func (c *chain) Reset() {
	for i, l := range c.link {
		if l.t != nil {
			l.t.Reset()
		}
		c.link[i].p, c.link[i].n = 0, 0
	}
}




func (c *chain) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	
	srcL := &c.link[0]
	dstL := &c.link[len(c.link)-1]
	srcL.b, srcL.p, srcL.n = src, 0, len(src)
	dstL.b, dstL.n = dst, 0
	var lastFull, needProgress bool 

	
	
	
	
	
	for low, i, high := c.errStart, c.errStart, len(c.link)-2; low <= i && i <= high; {
		in, out := &c.link[i], &c.link[i+1]
		nDst, nSrc, err0 := in.t.Transform(out.dst(), in.src(), atEOF && low == i)
		out.n += nDst
		in.p += nSrc
		if i > 0 && in.p == in.n {
			in.p, in.n = 0, 0
		}
		needProgress, lastFull = lastFull, false
		switch err0 {
		case ErrShortDst:
			
			
			if i == high {
				return dstL.n, srcL.p, ErrShortDst
			}
			if out.n != 0 {
				i++
				
				
				
				
				lastFull = true
				continue
			}
			
			
			c.fatalError(i, errShortInternal)
		case ErrShortSrc:
			if i == 0 {
				
				err = ErrShortSrc
				break
			}
			
			
			
			if needProgress && nSrc == 0 || in.n-in.p == len(in.b) {
				
				
				
				c.fatalError(i, errShortInternal)
				break
			}
			
			in.p, in.n = 0, copy(in.b, in.src())
			fallthrough
		case nil:
			
			
			
			if i > low {
				i--
				continue
			}
		default:
			c.fatalError(i, err0)
		}
		
		
		i++
		low = i
	}

	
	
	
	if c.errStart > 0 {
		for i := 1; i < c.errStart; i++ {
			c.link[i].p, c.link[i].n = 0, 0
		}
		err, c.errStart, c.err = c.err, 0, nil
	}
	return dstL.n, srcL.p, err
}


func RemoveFunc(f func(r rune) bool) Transformer {
	return removeF(f)
}

type removeF func(r rune) bool

func (removeF) Reset() {}


func (t removeF) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for r, sz := rune(0), 0; len(src) > 0; src = src[sz:] {

		if r = rune(src[0]); r < utf8.RuneSelf {
			sz = 1
		} else {
			r, sz = utf8.DecodeRune(src)

			if sz == 1 {
				
				if !atEOF && !utf8.FullRune(src) {
					err = ErrShortSrc
					break
				}
				
				
				
				
				if !t(r) {
					if nDst+3 > len(dst) {
						err = ErrShortDst
						break
					}
					nDst += copy(dst[nDst:], "\uFFFD")
				}
				nSrc++
				continue
			}
		}

		if !t(r) {
			if nDst+sz > len(dst) {
				err = ErrShortDst
				break
			}
			nDst += copy(dst[nDst:], src[:sz])
		}
		nSrc += sz
	}
	return
}



func grow(b []byte, n int) []byte {
	m := len(b)
	if m <= 32 {
		m = 64
	} else if m <= 256 {
		m *= 2
	} else {
		m += m >> 1
	}
	buf := make([]byte, m)
	copy(buf, b[:n])
	return buf
}

const initialBufSize = 128



func String(t Transformer, s string) (result string, n int, err error) {
	t.Reset()
	if s == "" {
		
		
		if _, _, err := t.Transform(nil, nil, true); err == nil {
			return "", 0, nil
		}
	}

	
	
	buf := [2 * initialBufSize]byte{}
	dst := buf[:initialBufSize:initialBufSize]
	src := buf[initialBufSize : 2*initialBufSize]

	
	
	
	nDst, nSrc := 0, 0
	pDst, pSrc := 0, 0

	
	
	
	
	
	
	pPrefix := 0
	for {
		

		n := copy(src, s[pSrc:])
		nDst, nSrc, err = t.Transform(dst, src[:n], pSrc+n == len(s))
		pDst += nDst
		pSrc += nSrc

		
		
		if !bytes.Equal(dst[:nDst], src[:nSrc]) {
			break
		}
		pPrefix = pSrc
		if err == ErrShortDst {
			
			break
		} else if err == ErrShortSrc {
			if nSrc == 0 {
				
				break
			}
			
		} else if err != nil || pPrefix == len(s) {
			return string(s[:pPrefix]), pPrefix, err
		}
	}
	

	
	
	
	
	
	if pPrefix != 0 {
		newDst := dst
		if pDst > len(newDst) {
			newDst = make([]byte, len(s)+nDst-nSrc)
		}
		copy(newDst[pPrefix:pDst], dst[:nDst])
		copy(newDst[:pPrefix], s[:pPrefix])
		dst = newDst
	}

	
	
	if (err == nil && pSrc == len(s)) ||
		(err != nil && err != ErrShortDst && err != ErrShortSrc) {
		return string(dst[:pDst]), pSrc, err
	}

	
	for {
		n := copy(src, s[pSrc:])
		nDst, nSrc, err := t.Transform(dst[pDst:], src[:n], pSrc+n == len(s))
		pDst += nDst
		pSrc += nSrc

		
		
		if err == ErrShortDst {
			if nDst == 0 {
				dst = grow(dst, pDst)
			}
		} else if err == ErrShortSrc {
			if nSrc == 0 {
				src = grow(src, 0)
			}
		} else if err != nil || pSrc == len(s) {
			return string(dst[:pDst]), pSrc, err
		}
	}
}



func Bytes(t Transformer, b []byte) (result []byte, n int, err error) {
	return doAppend(t, 0, make([]byte, len(b)), b)
}



func Append(t Transformer, dst, src []byte) (result []byte, n int, err error) {
	if len(dst) == cap(dst) {
		n := len(src) + len(dst) 
		b := make([]byte, n)
		dst = b[:copy(b, dst)]
	}
	return doAppend(t, len(dst), dst[:cap(dst)], src)
}

func doAppend(t Transformer, pDst int, dst, src []byte) (result []byte, n int, err error) {
	t.Reset()
	pSrc := 0
	for {
		nDst, nSrc, err := t.Transform(dst[pDst:], src[pSrc:], true)
		pDst += nDst
		pSrc += nSrc
		if err != ErrShortDst {
			return dst[:pDst], pSrc, err
		}

		
		
		if nDst == 0 {
			dst = grow(dst, pDst)
		}
	}
}
