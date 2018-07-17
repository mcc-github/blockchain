














package idna 

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/secure/bidirule"
	"golang.org/x/text/unicode/bidi"
	"golang.org/x/text/unicode/norm"
)

















func ToASCII(s string) (string, error) {
	return Punycode.process(s, true)
}


func ToUnicode(s string) (string, error) {
	return Punycode.process(s, false)
}


type Option func(*options)






func Transitional(transitional bool) Option {
	return func(o *options) { o.transitional = true }
}



func VerifyDNSLength(verify bool) Option {
	return func(o *options) { o.verifyDNSLength = verify }
}






func RemoveLeadingDots(remove bool) Option {
	return func(o *options) { o.removeLeadingDots = remove }
}




func ValidateLabels(enable bool) Option {
	return func(o *options) {
		
		
		if o.mapping == nil && enable {
			o.mapping = normalize
		}
		o.trie = trie
		o.validateLabels = enable
		o.fromPuny = validateFromPunycode
	}
}









func StrictDomainName(use bool) Option {
	return func(o *options) {
		o.trie = trie
		o.useSTD3Rules = use
		o.fromPuny = validateFromPunycode
	}
}






func BidiRule() Option {
	return func(o *options) { o.bidirule = bidirule.ValidString }
}



func ValidateForRegistration() Option {
	return func(o *options) {
		o.mapping = validateRegistration
		StrictDomainName(true)(o)
		ValidateLabels(true)(o)
		VerifyDNSLength(true)(o)
		BidiRule()(o)
	}
}









func MapForLookup() Option {
	return func(o *options) {
		o.mapping = validateAndMap
		StrictDomainName(true)(o)
		ValidateLabels(true)(o)
	}
}

type options struct {
	transitional      bool
	useSTD3Rules      bool
	validateLabels    bool
	verifyDNSLength   bool
	removeLeadingDots bool

	trie *idnaTrie

	
	fromPuny func(p *Profile, s string) error

	
	
	mapping func(p *Profile, s string) (mapped string, isBidi bool, err error)

	
	
	bidirule func(s string) bool
}


type Profile struct {
	options
}

func apply(o *options, opts []Option) {
	for _, f := range opts {
		f(o)
	}
}









func New(o ...Option) *Profile {
	p := &Profile{}
	apply(&p.options, o)
	return p
}





func (p *Profile) ToASCII(s string) (string, error) {
	return p.process(s, true)
}





func (p *Profile) ToUnicode(s string) (string, error) {
	pp := *p
	pp.transitional = false
	return pp.process(s, false)
}



func (p *Profile) String() string {
	s := ""
	if p.transitional {
		s = "Transitional"
	} else {
		s = "NonTransitional"
	}
	if p.useSTD3Rules {
		s += ":UseSTD3Rules"
	}
	if p.validateLabels {
		s += ":ValidateLabels"
	}
	if p.verifyDNSLength {
		s += ":VerifyDNSLength"
	}
	return s
}

var (
	
	
	Punycode *Profile = punycode

	
	
	
	Lookup *Profile = lookup

	
	
	Display *Profile = display

	
	
	Registration *Profile = registration

	punycode = &Profile{}
	lookup   = &Profile{options{
		transitional:   true,
		useSTD3Rules:   true,
		validateLabels: true,
		trie:           trie,
		fromPuny:       validateFromPunycode,
		mapping:        validateAndMap,
		bidirule:       bidirule.ValidString,
	}}
	display = &Profile{options{
		useSTD3Rules:   true,
		validateLabels: true,
		trie:           trie,
		fromPuny:       validateFromPunycode,
		mapping:        validateAndMap,
		bidirule:       bidirule.ValidString,
	}}
	registration = &Profile{options{
		useSTD3Rules:    true,
		validateLabels:  true,
		verifyDNSLength: true,
		trie:            trie,
		fromPuny:        validateFromPunycode,
		mapping:         validateRegistration,
		bidirule:        bidirule.ValidString,
	}}

	
	
	
)

type labelError struct{ label, code_ string }

func (e labelError) code() string { return e.code_ }
func (e labelError) Error() string {
	return fmt.Sprintf("idna: invalid label %q", e.label)
}

type runeError rune

func (e runeError) code() string { return "P1" }
func (e runeError) Error() string {
	return fmt.Sprintf("idna: disallowed rune %U", e)
}



func (p *Profile) process(s string, toASCII bool) (string, error) {
	var err error
	var isBidi bool
	if p.mapping != nil {
		s, isBidi, err = p.mapping(p, s)
	}
	
	if p.removeLeadingDots {
		for ; len(s) > 0 && s[0] == '.'; s = s[1:] {
		}
	}
	
	
	
	if err == nil && p.verifyDNSLength && s == "" {
		err = &labelError{s, "A4"}
	}
	labels := labelIter{orig: s}
	for ; !labels.done(); labels.next() {
		label := labels.label()
		if label == "" {
			
			
			if err == nil && p.verifyDNSLength {
				err = &labelError{s, "A4"}
			}
			continue
		}
		if strings.HasPrefix(label, acePrefix) {
			u, err2 := decode(label[len(acePrefix):])
			if err2 != nil {
				if err == nil {
					err = err2
				}
				
				continue
			}
			isBidi = isBidi || bidirule.DirectionString(u) != bidi.LeftToRight
			labels.set(u)
			if err == nil && p.validateLabels {
				err = p.fromPuny(p, u)
			}
			if err == nil {
				
				
				
				err = p.validateLabel(u)
			}
		} else if err == nil {
			err = p.validateLabel(label)
		}
	}
	if isBidi && p.bidirule != nil && err == nil {
		for labels.reset(); !labels.done(); labels.next() {
			if !p.bidirule(labels.label()) {
				err = &labelError{s, "B"}
				break
			}
		}
	}
	if toASCII {
		for labels.reset(); !labels.done(); labels.next() {
			label := labels.label()
			if !ascii(label) {
				a, err2 := encode(acePrefix, label)
				if err == nil {
					err = err2
				}
				label = a
				labels.set(a)
			}
			n := len(label)
			if p.verifyDNSLength && err == nil && (n == 0 || n > 63) {
				err = &labelError{label, "A4"}
			}
		}
	}
	s = labels.result()
	if toASCII && p.verifyDNSLength && err == nil {
		
		n := len(s)
		if n > 0 && s[n-1] == '.' {
			n--
		}
		if len(s) < 1 || n > 253 {
			err = &labelError{s, "A4"}
		}
	}
	return s, err
}

func normalize(p *Profile, s string) (mapped string, isBidi bool, err error) {
	
	
	
	mapped = norm.NFC.String(s)
	isBidi = bidirule.DirectionString(mapped) == bidi.RightToLeft
	return mapped, isBidi, nil
}

func validateRegistration(p *Profile, s string) (idem string, bidi bool, err error) {
	
	if !norm.NFC.IsNormalString(s) {
		return s, false, &labelError{s, "V1"}
	}
	for i := 0; i < len(s); {
		v, sz := trie.lookupString(s[i:])
		if sz == 0 {
			return s, bidi, runeError(utf8.RuneError)
		}
		bidi = bidi || info(v).isBidi(s[i:])
		
		switch p.simplify(info(v).category()) {
		
		
		case valid, deviation:
		case disallowed, mapped, unknown, ignored:
			r, _ := utf8.DecodeRuneInString(s[i:])
			return s, bidi, runeError(r)
		}
		i += sz
	}
	return s, bidi, nil
}

func (c info) isBidi(s string) bool {
	if !c.isMapped() {
		return c&attributesMask == rtl
	}
	
	
	p, _ := bidi.LookupString(s)
	switch p.Class() {
	case bidi.R, bidi.AL, bidi.AN:
		return true
	}
	return false
}

func validateAndMap(p *Profile, s string) (vm string, bidi bool, err error) {
	var (
		b []byte
		k int
	)
	
	
	
	
	var combinedInfoBits info
	for i := 0; i < len(s); {
		v, sz := trie.lookupString(s[i:])
		if sz == 0 {
			b = append(b, s[k:i]...)
			b = append(b, "\ufffd"...)
			k = len(s)
			if err == nil {
				err = runeError(utf8.RuneError)
			}
			break
		}
		combinedInfoBits |= info(v)
		bidi = bidi || info(v).isBidi(s[i:])
		start := i
		i += sz
		
		switch p.simplify(info(v).category()) {
		case valid:
			continue
		case disallowed:
			if err == nil {
				r, _ := utf8.DecodeRuneInString(s[start:])
				err = runeError(r)
			}
			continue
		case mapped, deviation:
			b = append(b, s[k:start]...)
			b = info(v).appendMapping(b, s[start:i])
		case ignored:
			b = append(b, s[k:start]...)
			
		case unknown:
			b = append(b, s[k:start]...)
			b = append(b, "\ufffd"...)
		}
		k = i
	}
	if k == 0 {
		
		if combinedInfoBits&mayNeedNorm != 0 {
			s = norm.NFC.String(s)
		}
	} else {
		b = append(b, s[k:]...)
		if norm.NFC.QuickSpan(b) != len(b) {
			b = norm.NFC.Bytes(b)
		}
		
		s = string(b)
	}
	return s, bidi, err
}


type labelIter struct {
	orig     string
	slice    []string
	curStart int
	curEnd   int
	i        int
}

func (l *labelIter) reset() {
	l.curStart = 0
	l.curEnd = 0
	l.i = 0
}

func (l *labelIter) done() bool {
	return l.curStart >= len(l.orig)
}

func (l *labelIter) result() string {
	if l.slice != nil {
		return strings.Join(l.slice, ".")
	}
	return l.orig
}

func (l *labelIter) label() string {
	if l.slice != nil {
		return l.slice[l.i]
	}
	p := strings.IndexByte(l.orig[l.curStart:], '.')
	l.curEnd = l.curStart + p
	if p == -1 {
		l.curEnd = len(l.orig)
	}
	return l.orig[l.curStart:l.curEnd]
}


func (l *labelIter) next() {
	l.i++
	if l.slice != nil {
		if l.i >= len(l.slice) || l.i == len(l.slice)-1 && l.slice[l.i] == "" {
			l.curStart = len(l.orig)
		}
	} else {
		l.curStart = l.curEnd + 1
		if l.curStart == len(l.orig)-1 && l.orig[l.curStart] == '.' {
			l.curStart = len(l.orig)
		}
	}
}

func (l *labelIter) set(s string) {
	if l.slice == nil {
		l.slice = strings.Split(l.orig, ".")
	}
	l.slice[l.i] = s
}


const acePrefix = "xn--"

func (p *Profile) simplify(cat category) category {
	switch cat {
	case disallowedSTD3Mapped:
		if p.useSTD3Rules {
			cat = disallowed
		} else {
			cat = mapped
		}
	case disallowedSTD3Valid:
		if p.useSTD3Rules {
			cat = disallowed
		} else {
			cat = valid
		}
	case deviation:
		if !p.transitional {
			cat = valid
		}
	case validNV8, validXV8:
		
		cat = valid
	}
	return cat
}

func validateFromPunycode(p *Profile, s string) error {
	if !norm.NFC.IsNormalString(s) {
		return &labelError{s, "V1"}
	}
	
	
	for i := 0; i < len(s); {
		v, sz := trie.lookupString(s[i:])
		if sz == 0 {
			return runeError(utf8.RuneError)
		}
		if c := p.simplify(info(v).category()); c != valid && c != deviation {
			return &labelError{s, "V6"}
		}
		i += sz
	}
	return nil
}

const (
	zwnj = "\u200c"
	zwj  = "\u200d"
)

type joinState int8

const (
	stateStart joinState = iota
	stateVirama
	stateBefore
	stateBeforeVirama
	stateAfter
	stateFAIL
)

var joinStates = [][numJoinTypes]joinState{
	stateStart: {
		joiningL:   stateBefore,
		joiningD:   stateBefore,
		joinZWNJ:   stateFAIL,
		joinZWJ:    stateFAIL,
		joinVirama: stateVirama,
	},
	stateVirama: {
		joiningL: stateBefore,
		joiningD: stateBefore,
	},
	stateBefore: {
		joiningL:   stateBefore,
		joiningD:   stateBefore,
		joiningT:   stateBefore,
		joinZWNJ:   stateAfter,
		joinZWJ:    stateFAIL,
		joinVirama: stateBeforeVirama,
	},
	stateBeforeVirama: {
		joiningL: stateBefore,
		joiningD: stateBefore,
		joiningT: stateBefore,
	},
	stateAfter: {
		joiningL:   stateFAIL,
		joiningD:   stateBefore,
		joiningT:   stateAfter,
		joiningR:   stateStart,
		joinZWNJ:   stateFAIL,
		joinZWJ:    stateFAIL,
		joinVirama: stateAfter, 
	},
	stateFAIL: {
		0:          stateFAIL,
		joiningL:   stateFAIL,
		joiningD:   stateFAIL,
		joiningT:   stateFAIL,
		joiningR:   stateFAIL,
		joinZWNJ:   stateFAIL,
		joinZWJ:    stateFAIL,
		joinVirama: stateFAIL,
	},
}



func (p *Profile) validateLabel(s string) (err error) {
	if s == "" {
		if p.verifyDNSLength {
			return &labelError{s, "A4"}
		}
		return nil
	}
	if !p.validateLabels {
		return nil
	}
	trie := p.trie 
	if len(s) > 4 && s[2] == '-' && s[3] == '-' {
		return &labelError{s, "V2"}
	}
	if s[0] == '-' || s[len(s)-1] == '-' {
		return &labelError{s, "V3"}
	}
	
	v, sz := trie.lookupString(s)
	x := info(v)
	if x.isModifier() {
		return &labelError{s, "V5"}
	}
	
	if strings.Index(s, zwj) == -1 && strings.Index(s, zwnj) == -1 {
		return nil
	}
	st := stateStart
	for i := 0; ; {
		jt := x.joinType()
		if s[i:i+sz] == zwj {
			jt = joinZWJ
		} else if s[i:i+sz] == zwnj {
			jt = joinZWNJ
		}
		st = joinStates[st][jt]
		if x.isViramaModifier() {
			st = joinStates[st][joinVirama]
		}
		if i += sz; i == len(s) {
			break
		}
		v, sz = trie.lookupString(s[i:])
		x = info(v)
	}
	if st == stateFAIL || st == stateAfter {
		return &labelError{s, "C"}
	}
	return nil
}

func ascii(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}
