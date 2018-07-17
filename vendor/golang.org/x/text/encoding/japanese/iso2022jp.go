



package japanese

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)


var ISO2022JP encoding.Encoding = &iso2022JP

var iso2022JP = internal.Encoding{
	internal.FuncEncoding{iso2022JPNewDecoder, iso2022JPNewEncoder},
	"ISO-2022-JP",
	identifier.ISO2022JP,
}

func iso2022JPNewDecoder() transform.Transformer {
	return new(iso2022JPDecoder)
}

func iso2022JPNewEncoder() transform.Transformer {
	return new(iso2022JPEncoder)
}

const (
	asciiState = iota
	katakanaState
	jis0208State
	jis0212State
)

const asciiEsc = 0x1b

type iso2022JPDecoder int

func (d *iso2022JPDecoder) Reset() {
	*d = asciiState
}

func (d *iso2022JPDecoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	r, size := rune(0), 0
	for ; nSrc < len(src); nSrc += size {
		c0 := src[nSrc]
		if c0 >= utf8.RuneSelf {
			r, size = '\ufffd', 1
			goto write
		}

		if c0 == asciiEsc {
			if nSrc+2 >= len(src) {
				if !atEOF {
					return nDst, nSrc, transform.ErrShortSrc
				}
				
				r, size = '\ufffd', 1
				goto write
			}
			size = 3
			c1 := src[nSrc+1]
			c2 := src[nSrc+2]
			switch {
			case c1 == '$' && (c2 == '@' || c2 == 'B'): 
				*d = jis0208State
				continue
			case c1 == '$' && c2 == '(': 
				if nSrc+3 >= len(src) {
					if !atEOF {
						return nDst, nSrc, transform.ErrShortSrc
					}
					r, size = '\ufffd', 1
					goto write
				}
				size = 4
				if src[nSrc+3] == 'D' {
					*d = jis0212State
					continue
				}
			case c1 == '(' && (c2 == 'B' || c2 == 'J'): 
				*d = asciiState
				continue
			case c1 == '(' && c2 == 'I': 
				*d = katakanaState
				continue
			}
			r, size = '\ufffd', 1
			goto write
		}

		switch *d {
		case asciiState:
			r, size = rune(c0), 1

		case katakanaState:
			if c0 < 0x21 || 0x60 <= c0 {
				r, size = '\ufffd', 1
				goto write
			}
			r, size = rune(c0)+(0xff61-0x21), 1

		default:
			if c0 == 0x0a {
				*d = asciiState
				r, size = rune(c0), 1
				goto write
			}
			if nSrc+1 >= len(src) {
				if !atEOF {
					return nDst, nSrc, transform.ErrShortSrc
				}
				r, size = '\ufffd', 1
				goto write
			}
			size = 2
			c1 := src[nSrc+1]
			i := int(c0-0x21)*94 + int(c1-0x21)
			if *d == jis0208State && i < len(jis0208Decode) {
				r = rune(jis0208Decode[i])
			} else if *d == jis0212State && i < len(jis0212Decode) {
				r = rune(jis0212Decode[i])
			} else {
				r = '\ufffd'
				goto write
			}
			if r == 0 {
				r = '\ufffd'
			}
		}

	write:
		if nDst+utf8.RuneLen(r) > len(dst) {
			return nDst, nSrc, transform.ErrShortDst
		}
		nDst += utf8.EncodeRune(dst[nDst:], r)
	}
	return nDst, nSrc, err
}

type iso2022JPEncoder int

func (e *iso2022JPEncoder) Reset() {
	*e = asciiState
}

func (e *iso2022JPEncoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
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
			}

			
			
			
			
			
			
			
			
			
			switch {
			case encode0Low <= r && r < encode0High:
				if r = rune(encode0[r-encode0Low]); r>>tableShift == jis0208 {
					goto writeJIS
				}
			case encode1Low <= r && r < encode1High:
				if r = rune(encode1[r-encode1Low]); r>>tableShift == jis0208 {
					goto writeJIS
				}
			case encode2Low <= r && r < encode2High:
				if r = rune(encode2[r-encode2Low]); r>>tableShift == jis0208 {
					goto writeJIS
				}
			case encode3Low <= r && r < encode3High:
				if r = rune(encode3[r-encode3Low]); r>>tableShift == jis0208 {
					goto writeJIS
				}
			case encode4Low <= r && r < encode4High:
				if r = rune(encode4[r-encode4Low]); r>>tableShift == jis0208 {
					goto writeJIS
				}
			case encode5Low <= r && r < encode5High:
				if 0xff61 <= r && r < 0xffa0 {
					goto writeKatakana
				}
				if r = rune(encode5[r-encode5Low]); r>>tableShift == jis0208 {
					goto writeJIS
				}
			}

			
			
			if *e != asciiState {
				if nDst+3 > len(dst) {
					err = transform.ErrShortDst
					break
				}
				*e = asciiState
				dst[nDst+0] = asciiEsc
				dst[nDst+1] = '('
				dst[nDst+2] = 'B'
				nDst += 3
			}
			err = internal.ErrASCIIReplacement
			break
		}

		if *e != asciiState {
			if nDst+4 > len(dst) {
				err = transform.ErrShortDst
				break
			}
			*e = asciiState
			dst[nDst+0] = asciiEsc
			dst[nDst+1] = '('
			dst[nDst+2] = 'B'
			nDst += 3
		} else if nDst >= len(dst) {
			err = transform.ErrShortDst
			break
		}
		dst[nDst] = uint8(r)
		nDst++
		continue

	writeJIS:
		if *e != jis0208State {
			if nDst+5 > len(dst) {
				err = transform.ErrShortDst
				break
			}
			*e = jis0208State
			dst[nDst+0] = asciiEsc
			dst[nDst+1] = '$'
			dst[nDst+2] = 'B'
			nDst += 3
		} else if nDst+2 > len(dst) {
			err = transform.ErrShortDst
			break
		}
		dst[nDst+0] = 0x21 + uint8(r>>codeShift)&codeMask
		dst[nDst+1] = 0x21 + uint8(r)&codeMask
		nDst += 2
		continue

	writeKatakana:
		if *e != katakanaState {
			if nDst+4 > len(dst) {
				err = transform.ErrShortDst
				break
			}
			*e = katakanaState
			dst[nDst+0] = asciiEsc
			dst[nDst+1] = '('
			dst[nDst+2] = 'I'
			nDst += 3
		} else if nDst >= len(dst) {
			err = transform.ErrShortDst
			break
		}
		dst[nDst] = uint8(r - (0xff61 - 0x21))
		nDst++
		continue
	}
	if atEOF && err == nil && *e != asciiState {
		if nDst+3 > len(dst) {
			err = transform.ErrShortDst
		} else {
			*e = asciiState
			dst[nDst+0] = asciiEsc
			dst[nDst+1] = '('
			dst[nDst+2] = 'B'
			nDst += 3
		}
	}
	return nDst, nSrc, err
}
