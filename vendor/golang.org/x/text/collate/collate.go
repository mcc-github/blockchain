









package collate 

import (
	"bytes"
	"strings"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/language"
)



type Collator struct {
	options

	sorter sorter

	_iter [2]iter
}

func (c *Collator) iter(i int) *iter {
	
	return &c._iter[i]
}


func Supported() []language.Tag {
	

	t := make([]language.Tag, len(tags))
	copy(t, tags)
	return t
}

func init() {
	ids := strings.Split(availableLocales, ",")
	tags = make([]language.Tag, len(ids))
	for i, s := range ids {
		tags[i] = language.Raw.MustParse(s)
	}
}

var tags []language.Tag


func New(t language.Tag, o ...Option) *Collator {
	index := colltab.MatchLang(t, tags)
	c := newCollator(getTable(locales[index]))

	
	c.setFromTag(t)

	
	c.setOptions(o)

	c.init()
	return c
}


func NewFromTable(w colltab.Weighter, o ...Option) *Collator {
	c := newCollator(w)
	c.setOptions(o)
	c.init()
	return c
}

func (c *Collator) init() {
	if c.numeric {
		c.t = colltab.NewNumericWeighter(c.t)
	}
	c._iter[0].init(c)
	c._iter[1].init(c)
}


type Buffer struct {
	buf [4096]byte
	key []byte
}

func (b *Buffer) init() {
	if b.key == nil {
		b.key = b.buf[:0]
	}
}


func (b *Buffer) Reset() {
	b.key = b.key[:0]
}



func (c *Collator) Compare(a, b []byte) int {
	
	
	c.iter(0).SetInput(a)
	c.iter(1).SetInput(b)
	if res := c.compare(); res != 0 {
		return res
	}
	if !c.ignore[colltab.Identity] {
		return bytes.Compare(a, b)
	}
	return 0
}



func (c *Collator) CompareString(a, b string) int {
	
	
	c.iter(0).SetInputString(a)
	c.iter(1).SetInputString(b)
	if res := c.compare(); res != 0 {
		return res
	}
	if !c.ignore[colltab.Identity] {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
	}
	return 0
}

func compareLevel(f func(i *iter) int, a, b *iter) int {
	a.pce = 0
	b.pce = 0
	for {
		va := f(a)
		vb := f(b)
		if va != vb {
			if va < vb {
				return -1
			}
			return 1
		} else if va == 0 {
			break
		}
	}
	return 0
}

func (c *Collator) compare() int {
	ia, ib := c.iter(0), c.iter(1)
	
	if c.alternate != altShifted {
		
		if res := compareLevel((*iter).nextPrimary, ia, ib); res != 0 {
			return res
		}
	} else {
		
	}
	if !c.ignore[colltab.Secondary] {
		f := (*iter).nextSecondary
		if c.backwards {
			f = (*iter).prevSecondary
		}
		if res := compareLevel(f, ia, ib); res != 0 {
			return res
		}
	}
	
	if !c.ignore[colltab.Tertiary] || c.caseLevel {
		if res := compareLevel((*iter).nextTertiary, ia, ib); res != 0 {
			return res
		}
		if !c.ignore[colltab.Quaternary] {
			if res := compareLevel((*iter).nextQuaternary, ia, ib); res != 0 {
				return res
			}
		}
	}
	return 0
}





func (c *Collator) Key(buf *Buffer, str []byte) []byte {
	
	buf.init()
	return c.key(buf, c.getColElems(str))
}





func (c *Collator) KeyFromString(buf *Buffer, str string) []byte {
	
	buf.init()
	return c.key(buf, c.getColElemsString(str))
}

func (c *Collator) key(buf *Buffer, w []colltab.Elem) []byte {
	processWeights(c.alternate, c.t.Top(), w)
	kn := len(buf.key)
	c.keyFromElems(buf, w)
	return buf.key[kn:]
}

func (c *Collator) getColElems(str []byte) []colltab.Elem {
	i := c.iter(0)
	i.SetInput(str)
	for i.Next() {
	}
	return i.Elems
}

func (c *Collator) getColElemsString(str string) []colltab.Elem {
	i := c.iter(0)
	i.SetInputString(str)
	for i.Next() {
	}
	return i.Elems
}

type iter struct {
	wa [512]colltab.Elem

	colltab.Iter
	pce int
}

func (i *iter) init(c *Collator) {
	i.Weighter = c.t
	i.Elems = i.wa[:0]
}

func (i *iter) nextPrimary() int {
	for {
		for ; i.pce < i.N; i.pce++ {
			if v := i.Elems[i.pce].Primary(); v != 0 {
				i.pce++
				return v
			}
		}
		if !i.Next() {
			return 0
		}
	}
	panic("should not reach here")
}

func (i *iter) nextSecondary() int {
	for ; i.pce < len(i.Elems); i.pce++ {
		if v := i.Elems[i.pce].Secondary(); v != 0 {
			i.pce++
			return v
		}
	}
	return 0
}

func (i *iter) prevSecondary() int {
	for ; i.pce < len(i.Elems); i.pce++ {
		if v := i.Elems[len(i.Elems)-i.pce-1].Secondary(); v != 0 {
			i.pce++
			return v
		}
	}
	return 0
}

func (i *iter) nextTertiary() int {
	for ; i.pce < len(i.Elems); i.pce++ {
		if v := i.Elems[i.pce].Tertiary(); v != 0 {
			i.pce++
			return int(v)
		}
	}
	return 0
}

func (i *iter) nextQuaternary() int {
	for ; i.pce < len(i.Elems); i.pce++ {
		if v := i.Elems[i.pce].Quaternary(); v != 0 {
			i.pce++
			return v
		}
	}
	return 0
}

func appendPrimary(key []byte, p int) []byte {
	
	if p <= 0x7FFF {
		key = append(key, uint8(p>>8), uint8(p))
	} else {
		key = append(key, uint8(p>>16)|0x80, uint8(p>>8), uint8(p))
	}
	return key
}



func (c *Collator) keyFromElems(buf *Buffer, ws []colltab.Elem) {
	for _, v := range ws {
		if w := v.Primary(); w > 0 {
			buf.key = appendPrimary(buf.key, w)
		}
	}
	if !c.ignore[colltab.Secondary] {
		buf.key = append(buf.key, 0, 0)
		
		if !c.backwards {
			for _, v := range ws {
				if w := v.Secondary(); w > 0 {
					buf.key = append(buf.key, uint8(w>>8), uint8(w))
				}
			}
		} else {
			for i := len(ws) - 1; i >= 0; i-- {
				if w := ws[i].Secondary(); w > 0 {
					buf.key = append(buf.key, uint8(w>>8), uint8(w))
				}
			}
		}
	} else if c.caseLevel {
		buf.key = append(buf.key, 0, 0)
	}
	if !c.ignore[colltab.Tertiary] || c.caseLevel {
		buf.key = append(buf.key, 0, 0)
		for _, v := range ws {
			if w := v.Tertiary(); w > 0 {
				buf.key = append(buf.key, uint8(w))
			}
		}
		
		
		
		
		if !c.ignore[colltab.Quaternary] && c.alternate >= altShifted {
			if c.alternate == altShiftTrimmed {
				lastNonFFFF := len(buf.key)
				buf.key = append(buf.key, 0)
				for _, v := range ws {
					if w := v.Quaternary(); w == colltab.MaxQuaternary {
						buf.key = append(buf.key, 0xFF)
					} else if w > 0 {
						buf.key = appendPrimary(buf.key, w)
						lastNonFFFF = len(buf.key)
					}
				}
				buf.key = buf.key[:lastNonFFFF]
			} else {
				buf.key = append(buf.key, 0)
				for _, v := range ws {
					if w := v.Quaternary(); w == colltab.MaxQuaternary {
						buf.key = append(buf.key, 0xFF)
					} else if w > 0 {
						buf.key = appendPrimary(buf.key, w)
					}
				}
			}
		}
	}
}

func processWeights(vw alternateHandling, top uint32, wa []colltab.Elem) {
	ignore := false
	vtop := int(top)
	switch vw {
	case altShifted, altShiftTrimmed:
		for i := range wa {
			if p := wa[i].Primary(); p <= vtop && p != 0 {
				wa[i] = colltab.MakeQuaternary(p)
				ignore = true
			} else if p == 0 {
				if ignore {
					wa[i] = colltab.Ignore
				}
			} else {
				ignore = false
			}
		}
	case altBlanked:
		for i := range wa {
			if p := wa[i].Primary(); p <= vtop && (ignore || p != 0) {
				wa[i] = colltab.Ignore
				ignore = true
			} else {
				ignore = false
			}
		}
	}
}
