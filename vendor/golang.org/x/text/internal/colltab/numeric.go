



package colltab

import (
	"unicode"
	"unicode/utf8"
)







func NewNumericWeighter(w Weighter) Weighter {
	getElem := func(s string) Elem {
		elems, _ := w.AppendNextString(nil, s)
		return elems[0]
	}
	nine := getElem("9")

	
	
	ns, _ := MakeElem(nine.Primary()+1, nine.Secondary(), int(nine.Tertiary()), 0)

	return &numericWeighter{
		Weighter: w,

		
		
		
		
		
		zero:          getElem("0"),
		zeroSpecialLo: getElem("０"), 
		zeroSpecialHi: getElem("₀"), 
		nine:          nine,
		nineSpecialHi: getElem("₉"), 
		numberStart:   ns,
	}
}



type numericWeighter struct {
	Weighter

	
	
	
	
	
	
	
	

	zero          Elem 
	zeroSpecialLo Elem 
	zeroSpecialHi Elem 
	nine          Elem 
	nineSpecialHi Elem 
	numberStart   Elem
}



func (nw *numericWeighter) AppendNext(buf []Elem, s []byte) (ce []Elem, n int) {
	ce, n = nw.Weighter.AppendNext(buf, s)
	nc := numberConverter{
		elems: buf,
		w:     nw,
		b:     s,
	}
	isZero, ok := nc.checkNextDigit(ce)
	if !ok {
		return ce, n
	}
	
	nc.init(ce, len(buf), isZero)
	for n < len(s) {
		ce, sz := nw.Weighter.AppendNext(nc.elems, s[n:])
		nc.b = s
		n += sz
		if !nc.update(ce) {
			break
		}
	}
	return nc.result(), n
}



func (nw *numericWeighter) AppendNextString(buf []Elem, s string) (ce []Elem, n int) {
	ce, n = nw.Weighter.AppendNextString(buf, s)
	nc := numberConverter{
		elems: buf,
		w:     nw,
		s:     s,
	}
	isZero, ok := nc.checkNextDigit(ce)
	if !ok {
		return ce, n
	}
	nc.init(ce, len(buf), isZero)
	for n < len(s) {
		ce, sz := nw.Weighter.AppendNextString(nc.elems, s[n:])
		nc.s = s
		n += sz
		if !nc.update(ce) {
			break
		}
	}
	return nc.result(), n
}

type numberConverter struct {
	w *numericWeighter

	elems    []Elem
	nDigits  int
	lenIndex int

	s string 
	b []byte 
}



func (nc *numberConverter) init(elems []Elem, oldLen int, isZero bool) {
	
	
	if isZero {
		elems = append(elems[:oldLen], nc.w.numberStart, 0)
	} else {
		elems = append(elems, 0, 0)
		copy(elems[oldLen+2:], elems[oldLen:])
		elems[oldLen] = nc.w.numberStart
		elems[oldLen+1] = 0

		nc.nDigits = 1
	}
	nc.elems = elems
	nc.lenIndex = oldLen + 1
}



func (nc *numberConverter) checkNextDigit(bufNew []Elem) (isZero, ok bool) {
	if len(nc.elems) >= len(bufNew) {
		return false, false
	}
	e := bufNew[len(nc.elems)]
	if e < nc.w.zeroSpecialLo || nc.w.nine < e {
		
		return false, false
	}
	if e < nc.w.zero {
		if e > nc.w.nineSpecialHi {
			
			return false, false
		}
		if !nc.isDigit() {
			return false, false
		}
		isZero = e <= nc.w.zeroSpecialHi
	} else {
		
		isZero = e == nc.w.zero
	}
	
	if n := len(bufNew) - len(nc.elems); n > 1 {
		for i := len(nc.elems) + 1; i < len(bufNew); i++ {
			if bufNew[i].Primary() != 0 {
				return false, false
			}
		}
		
		
		
		
		
		
		
		
		
		
		if !nc.isDigit() {
			return false, false
		}
	}
	return isZero, true
}

func (nc *numberConverter) isDigit() bool {
	if nc.b != nil {
		r, _ := utf8.DecodeRune(nc.b)
		return unicode.In(r, unicode.Nd)
	}
	r, _ := utf8.DecodeRuneInString(nc.s)
	return unicode.In(r, unicode.Nd)
}










const maxDigits = 1<<maxPrimaryBits - 1

func (nc *numberConverter) update(elems []Elem) bool {
	isZero, ok := nc.checkNextDigit(elems)
	if nc.nDigits == 0 && isZero {
		return true
	}
	nc.elems = elems
	if !ok {
		return false
	}
	nc.nDigits++
	return nc.nDigits < maxDigits
}



func (nc *numberConverter) result() []Elem {
	e, _ := MakeElem(nc.nDigits, defaultSecondary, defaultTertiary, 0)
	nc.elems[nc.lenIndex] = e
	return nc.elems
}
