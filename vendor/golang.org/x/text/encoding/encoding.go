









package encoding 

import (
	"errors"
	"io"
	"strconv"
	"unicode/utf8"

	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)










type Encoding interface {
	
	NewDecoder() *Decoder

	
	NewEncoder() *Encoder
}






type Decoder struct {
	transform.Transformer

	
	
	
	_ struct{}
}



func (d *Decoder) Bytes(b []byte) ([]byte, error) {
	b, _, err := transform.Bytes(d, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}



func (d *Decoder) String(s string) (string, error) {
	s, _, err := transform.String(d, s)
	if err != nil {
		return "", err
	}
	return s, nil
}





func (d *Decoder) Reader(r io.Reader) io.Reader {
	return transform.NewReader(r, d)
}








type Encoder struct {
	transform.Transformer

	
	
	
	_ struct{}
}



func (e *Encoder) Bytes(b []byte) ([]byte, error) {
	b, _, err := transform.Bytes(e, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}



func (e *Encoder) String(s string) (string, error) {
	s, _, err := transform.String(e, s)
	if err != nil {
		return "", err
	}
	return s, nil
}





func (e *Encoder) Writer(w io.Writer) io.Writer {
	return transform.NewWriter(w, e)
}



const ASCIISub = '\x1a'



var Nop Encoding = nop{}

type nop struct{}

func (nop) NewDecoder() *Decoder {
	return &Decoder{Transformer: transform.Nop}
}
func (nop) NewEncoder() *Encoder {
	return &Encoder{Transformer: transform.Nop}
}







var Replacement Encoding = replacement{}

type replacement struct{}

func (replacement) NewDecoder() *Decoder {
	return &Decoder{Transformer: replacementDecoder{}}
}

func (replacement) NewEncoder() *Encoder {
	return &Encoder{Transformer: replacementEncoder{}}
}

func (replacement) ID() (mib identifier.MIB, other string) {
	return identifier.Replacement, ""
}

type replacementDecoder struct{ transform.NopResetter }

func (replacementDecoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if len(dst) < 3 {
		return 0, 0, transform.ErrShortDst
	}
	if atEOF {
		const fffd = "\ufffd"
		dst[0] = fffd[0]
		dst[1] = fffd[1]
		dst[2] = fffd[2]
		nDst = 3
	}
	return nDst, len(src), nil
}

type replacementEncoder struct{ transform.NopResetter }

func (replacementEncoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	r, size := rune(0), 0

	for ; nSrc < len(src); nSrc += size {
		r = rune(src[nSrc])

		
		if r < utf8.RuneSelf {
			size = 1

		} else {
			
			r, size = utf8.DecodeRune(src[nSrc:])
			if size == 1 {
				
				
				
				if !atEOF && !utf8.FullRune(src[nSrc:]) {
					err = transform.ErrShortSrc
					break
				}
				r = '\ufffd'
			}
		}

		if nDst+utf8.RuneLen(r) > len(dst) {
			err = transform.ErrShortDst
			break
		}
		nDst += utf8.EncodeRune(dst[nDst:], r)
	}
	return nDst, nSrc, err
}








func HTMLEscapeUnsupported(e *Encoder) *Encoder {
	return &Encoder{Transformer: &errorHandler{e, errorToHTML}}
}







func ReplaceUnsupported(e *Encoder) *Encoder {
	return &Encoder{Transformer: &errorHandler{e, errorToReplacement}}
}

type errorHandler struct {
	*Encoder
	handler func(dst []byte, r rune, err repertoireError) (n int, ok bool)
}


type repertoireError interface {
	Replacement() byte
}

func (h errorHandler) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	nDst, nSrc, err = h.Transformer.Transform(dst, src, atEOF)
	for err != nil {
		rerr, ok := err.(repertoireError)
		if !ok {
			return nDst, nSrc, err
		}
		r, sz := utf8.DecodeRune(src[nSrc:])
		n, ok := h.handler(dst[nDst:], r, rerr)
		if !ok {
			return nDst, nSrc, transform.ErrShortDst
		}
		err = nil
		nDst += n
		if nSrc += sz; nSrc < len(src) {
			var dn, sn int
			dn, sn, err = h.Transformer.Transform(dst[nDst:], src[nSrc:], atEOF)
			nDst += dn
			nSrc += sn
		}
	}
	return nDst, nSrc, err
}

func errorToHTML(dst []byte, r rune, err repertoireError) (n int, ok bool) {
	buf := [8]byte{}
	b := strconv.AppendUint(buf[:0], uint64(r), 10)
	if n = len(b) + len("&#;"); n >= len(dst) {
		return 0, false
	}
	dst[0] = '&'
	dst[1] = '#'
	dst[copy(dst[2:], b)+2] = ';'
	return n, true
}

func errorToReplacement(dst []byte, r rune, err repertoireError) (n int, ok bool) {
	if len(dst) == 0 {
		return 0, false
	}
	dst[0] = err.Replacement()
	return 1, true
}


var ErrInvalidUTF8 = errors.New("encoding: invalid UTF-8")



var UTF8Validator transform.Transformer = utf8Validator{}

type utf8Validator struct{ transform.NopResetter }

func (utf8Validator) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	n := len(src)
	if n > len(dst) {
		n = len(dst)
	}
	for i := 0; i < n; {
		if c := src[i]; c < utf8.RuneSelf {
			dst[i] = c
			i++
			continue
		}
		_, size := utf8.DecodeRune(src[i:])
		if size == 1 {
			
			
			
			err = ErrInvalidUTF8
			if !atEOF && !utf8.FullRune(src[i:]) {
				err = transform.ErrShortSrc
			}
			return i, i, err
		}
		if i+size > len(dst) {
			return i, i, transform.ErrShortDst
		}
		for ; size > 0; size-- {
			dst[i] = src[i]
			i++
		}
	}
	if len(src) > len(dst) {
		err = transform.ErrShortDst
	}
	return n, n, err
}
