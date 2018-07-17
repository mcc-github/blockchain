







package bidirule

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/bidi"
)




























var ErrInvalid = errors.New("bidirule: failed Bidi Rule")

type ruleState uint8

const (
	ruleInitial ruleState = iota
	ruleLTR
	ruleLTRFinal
	ruleRTL
	ruleRTLFinal
	ruleInvalid
)

type ruleTransition struct {
	next ruleState
	mask uint16
}

var transitions = [...][2]ruleTransition{
	
	
	
	ruleInitial: {
		{ruleLTRFinal, 1 << bidi.L},
		{ruleRTLFinal, 1<<bidi.R | 1<<bidi.AL},
	},
	ruleRTL: {
		
		
		
		{ruleRTLFinal, 1<<bidi.R | 1<<bidi.AL | 1<<bidi.EN | 1<<bidi.AN},

		
		
		
		{ruleRTL, 1<<bidi.ES | 1<<bidi.CS | 1<<bidi.ET | 1<<bidi.ON | 1<<bidi.BN | 1<<bidi.NSM},
	},
	ruleRTLFinal: {
		
		
		
		{ruleRTLFinal, 1<<bidi.R | 1<<bidi.AL | 1<<bidi.EN | 1<<bidi.AN | 1<<bidi.NSM},

		
		
		
		{ruleRTL, 1<<bidi.ES | 1<<bidi.CS | 1<<bidi.ET | 1<<bidi.ON | 1<<bidi.BN},
	},
	ruleLTR: {
		
		
		
		{ruleLTRFinal, 1<<bidi.L | 1<<bidi.EN},

		
		
		
		{ruleLTR, 1<<bidi.ES | 1<<bidi.CS | 1<<bidi.ET | 1<<bidi.ON | 1<<bidi.BN | 1<<bidi.NSM},
	},
	ruleLTRFinal: {
		
		
		
		{ruleLTRFinal, 1<<bidi.L | 1<<bidi.EN | 1<<bidi.NSM},

		
		
		
		{ruleLTR, 1<<bidi.ES | 1<<bidi.CS | 1<<bidi.ET | 1<<bidi.ON | 1<<bidi.BN},
	},
	ruleInvalid: {
		{ruleInvalid, 0},
		{ruleInvalid, 0},
	},
}



const exclusiveRTL = uint16(1<<bidi.EN | 1<<bidi.AN)










func Direction(b []byte) bidi.Direction {
	for i := 0; i < len(b); {
		e, sz := bidi.Lookup(b[i:])
		if sz == 0 {
			i++
		}
		c := e.Class()
		if c == bidi.R || c == bidi.AL || c == bidi.AN {
			return bidi.RightToLeft
		}
		i += sz
	}
	return bidi.LeftToRight
}




func DirectionString(s string) bidi.Direction {
	for i := 0; i < len(s); {
		e, sz := bidi.LookupString(s[i:])
		if sz == 0 {
			i++
			continue
		}
		c := e.Class()
		if c == bidi.R || c == bidi.AL || c == bidi.AN {
			return bidi.RightToLeft
		}
		i += sz
	}
	return bidi.LeftToRight
}


func Valid(b []byte) bool {
	var t Transformer
	if n, ok := t.advance(b); !ok || n < len(b) {
		return false
	}
	return t.isFinal()
}


func ValidString(s string) bool {
	var t Transformer
	if n, ok := t.advanceString(s); !ok || n < len(s) {
		return false
	}
	return t.isFinal()
}


func New() *Transformer {
	return &Transformer{}
}


type Transformer struct {
	state  ruleState
	hasRTL bool
	seen   uint16
}



func (t *Transformer) isRTL() bool {
	const isRTL = 1<<bidi.R | 1<<bidi.AL | 1<<bidi.AN
	return t.seen&isRTL != 0
}


func (t *Transformer) Reset() { *t = Transformer{} }



func (t *Transformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if len(dst) < len(src) {
		src = src[:len(dst)]
		atEOF = false
		err = transform.ErrShortDst
	}
	n, err1 := t.Span(src, atEOF)
	copy(dst, src[:n])
	if err == nil || err1 != nil && err1 != transform.ErrShortSrc {
		err = err1
	}
	return n, n, err
}


func (t *Transformer) Span(src []byte, atEOF bool) (n int, err error) {
	if t.state == ruleInvalid && t.isRTL() {
		return 0, ErrInvalid
	}
	n, ok := t.advance(src)
	switch {
	case !ok:
		err = ErrInvalid
	case n < len(src):
		if !atEOF {
			err = transform.ErrShortSrc
			break
		}
		err = ErrInvalid
	case !t.isFinal():
		err = ErrInvalid
	}
	return n, err
}



var asciiTable [128]bidi.Properties

func init() {
	for i := range asciiTable {
		p, _ := bidi.LookupRune(rune(i))
		asciiTable[i] = p
	}
}

func (t *Transformer) advance(s []byte) (n int, ok bool) {
	var e bidi.Properties
	var sz int
	for n < len(s) {
		if s[n] < utf8.RuneSelf {
			e, sz = asciiTable[s[n]], 1
		} else {
			e, sz = bidi.Lookup(s[n:])
			if sz <= 1 {
				if sz == 1 {
					
					
					
					return n, false
				}
				return n, true 
			}
		}
		
		
		c := uint16(1 << e.Class())
		t.seen |= c
		if t.seen&exclusiveRTL == exclusiveRTL {
			t.state = ruleInvalid
			return n, false
		}
		switch tr := transitions[t.state]; {
		case tr[0].mask&c != 0:
			t.state = tr[0].next
		case tr[1].mask&c != 0:
			t.state = tr[1].next
		default:
			t.state = ruleInvalid
			if t.isRTL() {
				return n, false
			}
		}
		n += sz
	}
	return n, true
}

func (t *Transformer) advanceString(s string) (n int, ok bool) {
	var e bidi.Properties
	var sz int
	for n < len(s) {
		if s[n] < utf8.RuneSelf {
			e, sz = asciiTable[s[n]], 1
		} else {
			e, sz = bidi.LookupString(s[n:])
			if sz <= 1 {
				if sz == 1 {
					return n, false 
				}
				return n, true 
			}
		}
		
		
		c := uint16(1 << e.Class())
		t.seen |= c
		if t.seen&exclusiveRTL == exclusiveRTL {
			t.state = ruleInvalid
			return n, false
		}
		switch tr := transitions[t.state]; {
		case tr[0].mask&c != 0:
			t.state = tr[0].next
		case tr[1].mask&c != 0:
			t.state = tr[1].next
		default:
			t.state = ruleInvalid
			if t.isRTL() {
				return n, false
			}
		}
		n += sz
	}
	return n, true
}
