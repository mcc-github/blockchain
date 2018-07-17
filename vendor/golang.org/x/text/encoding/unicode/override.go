



package unicode

import (
	"golang.org/x/text/transform"
)

















func BOMOverride(fallback transform.Transformer) transform.Transformer {
	
	
	
	
	return &bomOverride{fallback: fallback}
}

type bomOverride struct {
	fallback transform.Transformer
	current  transform.Transformer
}

func (d *bomOverride) Reset() {
	d.current = nil
	d.fallback.Reset()
}

var (
	
	
	utf16le = UTF16(LittleEndian, IgnoreBOM)
	utf16be = UTF16(BigEndian, IgnoreBOM)
)

const utf8BOM = "\ufeff"

func (d *bomOverride) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if d.current != nil {
		return d.current.Transform(dst, src, atEOF)
	}
	if len(src) < 3 && !atEOF {
		return 0, 0, transform.ErrShortSrc
	}
	d.current = d.fallback
	bomSize := 0
	if len(src) >= 2 {
		if src[0] == 0xFF && src[1] == 0xFE {
			d.current = utf16le.NewDecoder()
			bomSize = 2
		} else if src[0] == 0xFE && src[1] == 0xFF {
			d.current = utf16be.NewDecoder()
			bomSize = 2
		} else if len(src) >= 3 &&
			src[0] == utf8BOM[0] &&
			src[1] == utf8BOM[1] &&
			src[2] == utf8BOM[2] {
			d.current = transform.Nop
			bomSize = 3
		}
	}
	if bomSize < len(src) {
		nDst, nSrc, err = d.current.Transform(dst, src[bomSize:], atEOF)
	}
	return nDst, nSrc + bomSize, err
}
