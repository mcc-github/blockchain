



package collate

import (
	"sort"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)


func newCollator(t colltab.Weighter) *Collator {
	
	c := &Collator{
		options: options{
			ignore: [colltab.NumLevels]bool{
				colltab.Quaternary: true,
				colltab.Identity:   true,
			},
			f: norm.NFD,
			t: t,
		},
	}

	
	c.variableTop = t.Top()

	return c
}



type Option struct {
	priority int
	f        func(o *options)
}

type prioritizedOptions []Option

func (p prioritizedOptions) Len() int {
	return len(p)
}

func (p prioritizedOptions) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p prioritizedOptions) Less(i, j int) bool {
	return p[i].priority < p[j].priority
}

type options struct {
	
	ignore [colltab.NumLevels]bool

	
	
	caseLevel bool

	
	
	backwards bool

	
	
	
	
	numeric bool

	
	alternate alternateHandling

	
	
	variableTop uint32

	t colltab.Weighter

	f norm.Form
}

func (o *options) setOptions(opts []Option) {
	sort.Sort(prioritizedOptions(opts))
	for _, x := range opts {
		x.f(o)
	}
}




func OptionsFromTag(t language.Tag) Option {
	return Option{0, func(o *options) {
		o.setFromTag(t)
	}}
}

func (o *options) setFromTag(t language.Tag) {
	o.caseLevel = ldmlBool(t, o.caseLevel, "kc")
	o.backwards = ldmlBool(t, o.backwards, "kb")
	o.numeric = ldmlBool(t, o.numeric, "kn")

	
	switch t.TypeForKey("ks") { 
	case "level1":
		o.ignore[colltab.Secondary] = true
		o.ignore[colltab.Tertiary] = true
	case "level2":
		o.ignore[colltab.Tertiary] = true
	case "level3", "":
		
	case "level4":
		o.ignore[colltab.Quaternary] = false
	case "identic":
		o.ignore[colltab.Quaternary] = false
		o.ignore[colltab.Identity] = false
	}

	switch t.TypeForKey("ka") {
	case "shifted":
		o.alternate = altShifted
	
	
	
	
	case "blanked":
		o.alternate = altBlanked
	case "posix":
		o.alternate = altShiftTrimmed
	}

	

	
	
	
}

func ldmlBool(t language.Tag, old bool, key string) bool {
	switch t.TypeForKey(key) {
	case "true":
		return true
	case "false":
		return false
	default:
		return old
	}
}

var (
	
	IgnoreCase Option = ignoreCase
	ignoreCase        = Option{3, ignoreCaseF}

	
	IgnoreDiacritics Option = ignoreDiacritics
	ignoreDiacritics        = Option{3, ignoreDiacriticsF}

	
	
	IgnoreWidth Option = ignoreWidth
	ignoreWidth        = Option{2, ignoreWidthF}

	
	Loose Option = loose
	loose        = Option{4, looseF}

	
	Force Option = force
	force        = Option{5, forceF}

	
	Numeric Option = numeric
	numeric        = Option{5, numericF}
)

func ignoreWidthF(o *options) {
	o.ignore[colltab.Tertiary] = true
	o.caseLevel = true
}

func ignoreDiacriticsF(o *options) {
	o.ignore[colltab.Secondary] = true
}

func ignoreCaseF(o *options) {
	o.ignore[colltab.Tertiary] = true
	o.caseLevel = false
}

func looseF(o *options) {
	ignoreWidthF(o)
	ignoreDiacriticsF(o)
	ignoreCaseF(o)
}

func forceF(o *options) {
	o.ignore[colltab.Identity] = false
}

func numericF(o *options) { o.numeric = true }


func Reorder(s ...string) Option {
	
	panic("TODO: implement")
}









type alternateHandling int

const (
	
	altNonIgnorable alternateHandling = iota

	
	
	
	altBlanked

	
	
	altShifted

	
	
	altShiftTrimmed
)
