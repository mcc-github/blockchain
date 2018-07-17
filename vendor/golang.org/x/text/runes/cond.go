



package runes

import (
	"unicode/utf8"

	"golang.org/x/text/transform"
)






















func If(s Set, tIn, tNotIn transform.Transformer) Transformer {
	if tIn == nil && tNotIn == nil {
		return Transformer{transform.Nop}
	}
	if tIn == nil {
		tIn = transform.Nop
	}
	if tNotIn == nil {
		tNotIn = transform.Nop
	}
	sIn, ok := tIn.(transform.SpanningTransformer)
	if !ok {
		sIn = dummySpan{tIn}
	}
	sNotIn, ok := tNotIn.(transform.SpanningTransformer)
	if !ok {
		sNotIn = dummySpan{tNotIn}
	}

	a := &cond{
		tIn:    sIn,
		tNotIn: sNotIn,
		f:      s.Contains,
	}
	a.Reset()
	return Transformer{a}
}

type dummySpan struct{ transform.Transformer }

func (d dummySpan) Span(src []byte, atEOF bool) (n int, err error) {
	return 0, transform.ErrEndOfSpan
}

type cond struct {
	tIn, tNotIn transform.SpanningTransformer
	f           func(rune) bool
	check       func(rune) bool               
	t           transform.SpanningTransformer 
}


func (t *cond) Reset() {
	t.check = t.is
	t.t = t.tIn
	t.t.Reset() 
}

func (t *cond) is(r rune) bool {
	if t.f(r) {
		return true
	}
	t.check = t.isNot
	t.t = t.tNotIn
	t.tNotIn.Reset()
	return false
}

func (t *cond) isNot(r rune) bool {
	if !t.f(r) {
		return true
	}
	t.check = t.is
	t.t = t.tIn
	t.tIn.Reset()
	return false
}






func (t *cond) Span(src []byte, atEOF bool) (n int, err error) {
	p := 0
	for n < len(src) && err == nil {
		
		
		const maxChunk = 4096
		max := len(src)
		if v := n + maxChunk; v < max {
			max = v
		}
		atEnd := false
		size := 0
		current := t.t
		for ; p < max; p += size {
			r := rune(src[p])
			if r < utf8.RuneSelf {
				size = 1
			} else if r, size = utf8.DecodeRune(src[p:]); size == 1 {
				if !atEOF && !utf8.FullRune(src[p:]) {
					err = transform.ErrShortSrc
					break
				}
			}
			if !t.check(r) {
				
				atEnd = true
				break
			}
		}
		n2, err2 := current.Span(src[n:p], atEnd || (atEOF && p == len(src)))
		n += n2
		if err2 != nil {
			return n, err2
		}
		
		p = n + size
	}
	return n, err
}

func (t *cond) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	p := 0
	for nSrc < len(src) && err == nil {
		
		
		
		const maxChunk = 4096
		max := len(src)
		if n := nSrc + maxChunk; n < len(src) {
			max = n
		}
		atEnd := false
		size := 0
		current := t.t
		for ; p < max; p += size {
			r := rune(src[p])
			if r < utf8.RuneSelf {
				size = 1
			} else if r, size = utf8.DecodeRune(src[p:]); size == 1 {
				if !atEOF && !utf8.FullRune(src[p:]) {
					err = transform.ErrShortSrc
					break
				}
			}
			if !t.check(r) {
				
				atEnd = true
				break
			}
		}
		nDst2, nSrc2, err2 := current.Transform(dst[nDst:], src[nSrc:p], atEnd || (atEOF && p == len(src)))
		nDst += nDst2
		nSrc += nSrc2
		if err2 != nil {
			return nDst, nSrc, err2
		}
		
		p = nSrc + size
	}
	return nDst, nSrc, err
}
