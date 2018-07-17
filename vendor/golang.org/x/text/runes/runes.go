




package runes 

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
)


type Set interface {
	
	Contains(r rune) bool
}

type setFunc func(rune) bool

func (s setFunc) Contains(r rune) bool {
	return s(r)
}






func In(rt *unicode.RangeTable) Set {
	return setFunc(func(r rune) bool { return unicode.Is(rt, r) })
}



func NotIn(rt *unicode.RangeTable) Set {
	return setFunc(func(r rune) bool { return !unicode.Is(rt, r) })
}


func Predicate(f func(rune) bool) Set {
	return setFunc(f)
}


type Transformer struct {
	t transform.SpanningTransformer
}

func (t Transformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	return t.t.Transform(dst, src, atEOF)
}

func (t Transformer) Span(b []byte, atEOF bool) (n int, err error) {
	return t.t.Span(b, atEOF)
}

func (t Transformer) Reset() { t.t.Reset() }




func (t Transformer) Bytes(b []byte) []byte {
	b, _, err := transform.Bytes(t, b)
	if err != nil {
		return nil
	}
	return b
}




func (t Transformer) String(s string) string {
	s, _, err := transform.String(t, s)
	if err != nil {
		return ""
	}
	return s
}






const runeErrorString = string(utf8.RuneError)



func Remove(s Set) Transformer {
	if f, ok := s.(setFunc); ok {
		
		
		
		return Transformer{remove(f)}
	}
	return Transformer{remove(s.Contains)}
}



type remove func(r rune) bool

func (remove) Reset() {}


func (t remove) Span(src []byte, atEOF bool) (n int, err error) {
	for r, size := rune(0), 0; n < len(src); {
		if r = rune(src[n]); r < utf8.RuneSelf {
			size = 1
		} else if r, size = utf8.DecodeRune(src[n:]); size == 1 {
			
			if !atEOF && !utf8.FullRune(src[n:]) {
				err = transform.ErrShortSrc
			} else {
				err = transform.ErrEndOfSpan
			}
			break
		}
		if t(r) {
			err = transform.ErrEndOfSpan
			break
		}
		n += size
	}
	return
}


func (t remove) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for r, size := rune(0), 0; nSrc < len(src); {
		if r = rune(src[nSrc]); r < utf8.RuneSelf {
			size = 1
		} else if r, size = utf8.DecodeRune(src[nSrc:]); size == 1 {
			
			if !atEOF && !utf8.FullRune(src[nSrc:]) {
				err = transform.ErrShortSrc
				break
			}
			
			
			
			
			if !t(utf8.RuneError) {
				if nDst+3 > len(dst) {
					err = transform.ErrShortDst
					break
				}
				dst[nDst+0] = runeErrorString[0]
				dst[nDst+1] = runeErrorString[1]
				dst[nDst+2] = runeErrorString[2]
				nDst += 3
			}
			nSrc++
			continue
		}
		if t(r) {
			nSrc += size
			continue
		}
		if nDst+size > len(dst) {
			err = transform.ErrShortDst
			break
		}
		for i := 0; i < size; i++ {
			dst[nDst] = src[nSrc]
			nDst++
			nSrc++
		}
	}
	return
}




func Map(mapping func(rune) rune) Transformer {
	return Transformer{mapper(mapping)}
}

type mapper func(rune) rune

func (mapper) Reset() {}


func (t mapper) Span(src []byte, atEOF bool) (n int, err error) {
	for r, size := rune(0), 0; n < len(src); n += size {
		if r = rune(src[n]); r < utf8.RuneSelf {
			size = 1
		} else if r, size = utf8.DecodeRune(src[n:]); size == 1 {
			
			if !atEOF && !utf8.FullRune(src[n:]) {
				err = transform.ErrShortSrc
			} else {
				err = transform.ErrEndOfSpan
			}
			break
		}
		if t(r) != r {
			err = transform.ErrEndOfSpan
			break
		}
	}
	return n, err
}


func (t mapper) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	var replacement rune
	var b [utf8.UTFMax]byte

	for r, size := rune(0), 0; nSrc < len(src); {
		if r = rune(src[nSrc]); r < utf8.RuneSelf {
			if replacement = t(r); replacement < utf8.RuneSelf {
				if nDst == len(dst) {
					err = transform.ErrShortDst
					break
				}
				dst[nDst] = byte(replacement)
				nDst++
				nSrc++
				continue
			}
			size = 1
		} else if r, size = utf8.DecodeRune(src[nSrc:]); size == 1 {
			
			if !atEOF && !utf8.FullRune(src[nSrc:]) {
				err = transform.ErrShortSrc
				break
			}

			if replacement = t(utf8.RuneError); replacement == utf8.RuneError {
				if nDst+3 > len(dst) {
					err = transform.ErrShortDst
					break
				}
				dst[nDst+0] = runeErrorString[0]
				dst[nDst+1] = runeErrorString[1]
				dst[nDst+2] = runeErrorString[2]
				nDst += 3
				nSrc++
				continue
			}
		} else if replacement = t(r); replacement == r {
			if nDst+size > len(dst) {
				err = transform.ErrShortDst
				break
			}
			for i := 0; i < size; i++ {
				dst[nDst] = src[nSrc]
				nDst++
				nSrc++
			}
			continue
		}

		n := utf8.EncodeRune(b[:], replacement)

		if nDst+n > len(dst) {
			err = transform.ErrShortDst
			break
		}
		for i := 0; i < n; i++ {
			dst[nDst] = b[i]
			nDst++
		}
		nSrc += size
	}
	return
}



func ReplaceIllFormed() Transformer {
	return Transformer{&replaceIllFormed{}}
}

type replaceIllFormed struct{ transform.NopResetter }

func (t replaceIllFormed) Span(src []byte, atEOF bool) (n int, err error) {
	for n < len(src) {
		
		if src[n] < utf8.RuneSelf {
			n++
			continue
		}

		r, size := utf8.DecodeRune(src[n:])

		
		if r != utf8.RuneError || size != 1 {
			n += size
			continue
		}

		
		if !atEOF && !utf8.FullRune(src[n:]) {
			err = transform.ErrShortSrc
			break
		}

		
		err = transform.ErrEndOfSpan
		break
	}
	return n, err
}

func (t replaceIllFormed) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nSrc < len(src) {
		
		if r := src[nSrc]; r < utf8.RuneSelf {
			if nDst == len(dst) {
				err = transform.ErrShortDst
				break
			}
			dst[nDst] = r
			nDst++
			nSrc++
			continue
		}

		
		if _, size := utf8.DecodeRune(src[nSrc:]); size != 1 {
			if size != copy(dst[nDst:], src[nSrc:nSrc+size]) {
				err = transform.ErrShortDst
				break
			}
			nDst += size
			nSrc += size
			continue
		}

		
		if !atEOF && !utf8.FullRune(src[nSrc:]) {
			err = transform.ErrShortSrc
			break
		}

		
		if nDst+3 > len(dst) {
			err = transform.ErrShortDst
			break
		}
		dst[nDst+0] = runeErrorString[0]
		dst[nDst+1] = runeErrorString[1]
		dst[nDst+2] = runeErrorString[2]
		nDst += 3
		nSrc++
	}
	return nDst, nSrc, err
}
