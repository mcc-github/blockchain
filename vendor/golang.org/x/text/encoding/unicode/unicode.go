




package unicode 

import (
	"errors"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/internal/utf8internal"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)









var UTF8 encoding.Encoding = utf8enc

var utf8enc = &internal.Encoding{
	&internal.SimpleEncoding{utf8Decoder{}, runes.ReplaceIllFormed()},
	"UTF-8",
	identifier.UTF8,
}

type utf8Decoder struct{ transform.NopResetter }

func (utf8Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	var pSrc int 
	var accept utf8internal.AcceptRange

	
	n := len(src)
	if len(dst) < n {
		err = transform.ErrShortDst
		n = len(dst)
		atEOF = false
	}
	for nSrc < n {
		c := src[nSrc]
		if c < utf8.RuneSelf {
			nSrc++
			continue
		}
		first := utf8internal.First[c]
		size := int(first & utf8internal.SizeMask)
		if first == utf8internal.FirstInvalid {
			goto handleInvalid 
		}
		accept = utf8internal.AcceptRanges[first>>utf8internal.AcceptShift]
		if nSrc+size > n {
			if !atEOF {
				
				
				
				if err == nil {
					err = transform.ErrShortSrc
				}
				break
			}
			
			switch {
			case nSrc+1 >= n || src[nSrc+1] < accept.Lo || accept.Hi < src[nSrc+1]:
				size = 1
			case nSrc+2 >= n || src[nSrc+2] < utf8internal.LoCB || utf8internal.HiCB < src[nSrc+2]:
				size = 2
			default:
				size = 3 
			}
			goto handleInvalid
		}
		if c = src[nSrc+1]; c < accept.Lo || accept.Hi < c {
			size = 1
			goto handleInvalid 
		} else if size == 2 {
		} else if c = src[nSrc+2]; c < utf8internal.LoCB || utf8internal.HiCB < c {
			size = 2
			goto handleInvalid 
		} else if size == 3 {
		} else if c = src[nSrc+3]; c < utf8internal.LoCB || utf8internal.HiCB < c {
			size = 3
			goto handleInvalid 
		}
		nSrc += size
		continue

	handleInvalid:
		
		nDst += copy(dst[nDst:], src[pSrc:nSrc])

		
		const runeError = "\ufffd"
		if nDst+len(runeError) > len(dst) {
			return nDst, nSrc, transform.ErrShortDst
		}
		nDst += copy(dst[nDst:], runeError)

		
		
		
		nSrc += size
		pSrc = nSrc

		
		if sz := len(dst) - nDst; sz < len(src)-nSrc {
			err = transform.ErrShortDst
			n = nSrc + sz
			atEOF = false
		}
	}
	return nDst + copy(dst[nDst:], src[pSrc:nSrc]), nSrc, err
}




























func UTF16(e Endianness, b BOMPolicy) encoding.Encoding {
	return utf16Encoding{config{e, b}, mibValue[e][b&bomMask]}
}





var mibValue = map[Endianness][numBOMValues]identifier.MIB{
	BigEndian: [numBOMValues]identifier.MIB{
		IgnoreBOM: identifier.UTF16BE,
		UseBOM:    identifier.UTF16, 
		
	},
	LittleEndian: [numBOMValues]identifier.MIB{
		IgnoreBOM: identifier.UTF16LE,
		UseBOM:    identifier.UTF16, 
		
	},
	
}


var All = []encoding.Encoding{
	UTF8,
	UTF16(BigEndian, UseBOM),
	UTF16(BigEndian, IgnoreBOM),
	UTF16(LittleEndian, IgnoreBOM),
}


type BOMPolicy uint8

const (
	writeBOM   BOMPolicy = 0x01
	acceptBOM  BOMPolicy = 0x02
	requireBOM BOMPolicy = 0x04
	bomMask    BOMPolicy = 0x07

	
	
	
	
	
	numBOMValues = 8 + 1

	
	IgnoreBOM BOMPolicy = 0
	

	
	
	UseBOM BOMPolicy = writeBOM | acceptBOM
	

	
	
	ExpectBOM BOMPolicy = writeBOM | acceptBOM | requireBOM
	
	

	
	
	
	
	
	
)


type Endianness bool

const (
	
	BigEndian Endianness = false
	
	LittleEndian Endianness = true
)



var ErrMissingBOM = errors.New("encoding: missing byte order mark")

type utf16Encoding struct {
	config
	mib identifier.MIB
}

type config struct {
	endianness Endianness
	bomPolicy  BOMPolicy
}

func (u utf16Encoding) NewDecoder() *encoding.Decoder {
	return &encoding.Decoder{Transformer: &utf16Decoder{
		initial: u.config,
		current: u.config,
	}}
}

func (u utf16Encoding) NewEncoder() *encoding.Encoder {
	return &encoding.Encoder{Transformer: &utf16Encoder{
		endianness:       u.endianness,
		initialBOMPolicy: u.bomPolicy,
		currentBOMPolicy: u.bomPolicy,
	}}
}

func (u utf16Encoding) ID() (mib identifier.MIB, other string) {
	return u.mib, ""
}

func (u utf16Encoding) String() string {
	e, b := "B", ""
	if u.endianness == LittleEndian {
		e = "L"
	}
	switch u.bomPolicy {
	case ExpectBOM:
		b = "Expect"
	case UseBOM:
		b = "Use"
	case IgnoreBOM:
		b = "Ignore"
	}
	return "UTF-16" + e + "E (" + b + " BOM)"
}

type utf16Decoder struct {
	initial config
	current config
}

func (u *utf16Decoder) Reset() {
	u.current = u.initial
}

func (u *utf16Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if len(src) == 0 {
		if atEOF && u.current.bomPolicy&requireBOM != 0 {
			return 0, 0, ErrMissingBOM
		}
		return 0, 0, nil
	}
	if u.current.bomPolicy&acceptBOM != 0 {
		if len(src) < 2 {
			return 0, 0, transform.ErrShortSrc
		}
		switch {
		case src[0] == 0xfe && src[1] == 0xff:
			u.current.endianness = BigEndian
			nSrc = 2
		case src[0] == 0xff && src[1] == 0xfe:
			u.current.endianness = LittleEndian
			nSrc = 2
		default:
			if u.current.bomPolicy&requireBOM != 0 {
				return 0, 0, ErrMissingBOM
			}
		}
		u.current.bomPolicy = IgnoreBOM
	}

	var r rune
	var dSize, sSize int
	for nSrc < len(src) {
		if nSrc+1 < len(src) {
			x := uint16(src[nSrc+0])<<8 | uint16(src[nSrc+1])
			if u.current.endianness == LittleEndian {
				x = x>>8 | x<<8
			}
			r, sSize = rune(x), 2
			if utf16.IsSurrogate(r) {
				if nSrc+3 < len(src) {
					x = uint16(src[nSrc+2])<<8 | uint16(src[nSrc+3])
					if u.current.endianness == LittleEndian {
						x = x>>8 | x<<8
					}
					
					if isHighSurrogate(rune(x)) {
						r, sSize = utf16.DecodeRune(r, rune(x)), 4
					}
				} else if !atEOF {
					err = transform.ErrShortSrc
					break
				}
			}
			if dSize = utf8.RuneLen(r); dSize < 0 {
				r, dSize = utf8.RuneError, 3
			}
		} else if atEOF {
			
			r, dSize, sSize = utf8.RuneError, 3, 1
		} else {
			err = transform.ErrShortSrc
			break
		}
		if nDst+dSize > len(dst) {
			err = transform.ErrShortDst
			break
		}
		nDst += utf8.EncodeRune(dst[nDst:], r)
		nSrc += sSize
	}
	return nDst, nSrc, err
}

func isHighSurrogate(r rune) bool {
	return 0xDC00 <= r && r <= 0xDFFF
}

type utf16Encoder struct {
	endianness       Endianness
	initialBOMPolicy BOMPolicy
	currentBOMPolicy BOMPolicy
}

func (u *utf16Encoder) Reset() {
	u.currentBOMPolicy = u.initialBOMPolicy
}

func (u *utf16Encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if u.currentBOMPolicy&writeBOM != 0 {
		if len(dst) < 2 {
			return 0, 0, transform.ErrShortDst
		}
		dst[0], dst[1] = 0xfe, 0xff
		u.currentBOMPolicy = IgnoreBOM
		nDst = 2
	}

	r, size := rune(0), 0
	for nSrc < len(src) {
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
		}

		if r <= 0xffff {
			if nDst+2 > len(dst) {
				err = transform.ErrShortDst
				break
			}
			dst[nDst+0] = uint8(r >> 8)
			dst[nDst+1] = uint8(r)
			nDst += 2
		} else {
			if nDst+4 > len(dst) {
				err = transform.ErrShortDst
				break
			}
			r1, r2 := utf16.EncodeRune(r)
			dst[nDst+0] = uint8(r1 >> 8)
			dst[nDst+1] = uint8(r1)
			dst[nDst+2] = uint8(r2 >> 8)
			dst[nDst+3] = uint8(r2)
			nDst += 4
		}
		nSrc += size
	}

	if u.endianness == LittleEndian {
		for i := 0; i < nDst; i += 2 {
			dst[i], dst[i+1] = dst[i+1], dst[i]
		}
	}
	return nDst, nSrc, err
}
