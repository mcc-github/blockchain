






package language




import (
	"errors"
	"fmt"
	"strings"
)

const (
	
	
	maxCoreSize = 12

	
	
	max99thPercentileSize = 32

	
	
	maxSimpleUExtensionSize = 14
)




type Tag struct {
	lang   langID
	region regionID
	
	
	
	
	
	
	
	script   scriptID
	pVariant byte   
	pExt     uint16 

	
	
	str string
}



func Make(s string) Tag {
	return Default.Make(s)
}



func (c CanonType) Make(s string) Tag {
	t, _ := c.Parse(s)
	return t
}



func (t Tag) Raw() (b Base, s Script, r Region) {
	return Base{t.lang}, Script{t.script}, Region{t.region}
}


func (t Tag) equalTags(a Tag) bool {
	return t.lang == a.lang && t.script == a.script && t.region == a.region
}


func (t Tag) IsRoot() bool {
	if int(t.pVariant) < len(t.str) {
		return false
	}
	return t.equalTags(und)
}


func (t Tag) private() bool {
	return t.str != "" && t.pVariant == 0
}


type CanonType int

const (
	
	DeprecatedBase CanonType = 1 << iota
	
	DeprecatedScript
	
	DeprecatedRegion
	
	SuppressScript
	
	
	Legacy
	
	
	Macro
	
	
	
	CLDR

	
	Raw CanonType = 0

	
	Deprecated = DeprecatedBase | DeprecatedScript | DeprecatedRegion

	
	BCP47 = Deprecated | SuppressScript

	
	All = BCP47 | Legacy | Macro

	
	
	
	
	
	Default = Deprecated | Legacy

	canonLang = DeprecatedBase | Legacy | Macro

	
)



func (t Tag) canonicalize(c CanonType) (Tag, bool) {
	if c == Raw {
		return t, false
	}
	changed := false
	if c&SuppressScript != 0 {
		if t.lang < langNoIndexOffset && uint8(t.script) == suppressScript[t.lang] {
			t.script = 0
			changed = true
		}
	}
	if c&canonLang != 0 {
		for {
			if l, aliasType := normLang(t.lang); l != t.lang {
				switch aliasType {
				case langLegacy:
					if c&Legacy != 0 {
						if t.lang == _sh && t.script == 0 {
							t.script = _Latn
						}
						t.lang = l
						changed = true
					}
				case langMacro:
					if c&Macro != 0 {
						
						
						
						
						
						
						
						
						
						if c&CLDR == 0 || t.lang != _nb {
							changed = true
							t.lang = l
						}
					}
				case langDeprecated:
					if c&DeprecatedBase != 0 {
						if t.lang == _mo && t.region == 0 {
							t.region = _MD
						}
						t.lang = l
						changed = true
						
						continue
					}
				}
			} else if c&Legacy != 0 && t.lang == _no && c&CLDR != 0 {
				t.lang = _nb
				changed = true
			}
			break
		}
	}
	if c&DeprecatedScript != 0 {
		if t.script == _Qaai {
			changed = true
			t.script = _Zinh
		}
	}
	if c&DeprecatedRegion != 0 {
		if r := normRegion(t.region); r != 0 {
			changed = true
			t.region = r
		}
	}
	return t, changed
}


func (c CanonType) Canonicalize(t Tag) (Tag, error) {
	t, changed := t.canonicalize(c)
	if changed {
		t.remakeString()
	}
	return t, nil
}






type Confidence int

const (
	No    Confidence = iota 
	Low                     
	High                    
	Exact                   
)

var confName = []string{"No", "Low", "High", "Exact"}

func (c Confidence) String() string {
	return confName[c]
}




func (t *Tag) remakeString() {
	if t.str == "" {
		return
	}
	extra := t.str[t.pVariant:]
	if t.pVariant > 0 {
		extra = extra[1:]
	}
	if t.equalTags(und) && strings.HasPrefix(extra, "x-") {
		t.str = extra
		t.pVariant = 0
		t.pExt = 0
		return
	}
	var buf [max99thPercentileSize]byte 
	b := buf[:t.genCoreBytes(buf[:])]
	if extra != "" {
		diff := len(b) - int(t.pVariant)
		b = append(b, '-')
		b = append(b, extra...)
		t.pVariant = uint8(int(t.pVariant) + diff)
		t.pExt = uint16(int(t.pExt) + diff)
	} else {
		t.pVariant = uint8(len(b))
		t.pExt = uint16(len(b))
	}
	t.str = string(b)
}




func (t *Tag) genCoreBytes(buf []byte) int {
	n := t.lang.stringToBuf(buf[:])
	if t.script != 0 {
		n += copy(buf[n:], "-")
		n += copy(buf[n:], t.script.String())
	}
	if t.region != 0 {
		n += copy(buf[n:], "-")
		n += copy(buf[n:], t.region.String())
	}
	return n
}


func (t Tag) String() string {
	if t.str != "" {
		return t.str
	}
	if t.script == 0 && t.region == 0 {
		return t.lang.String()
	}
	buf := [maxCoreSize]byte{}
	return string(buf[:t.genCoreBytes(buf[:])])
}


func (t Tag) MarshalText() (text []byte, err error) {
	if t.str != "" {
		text = append(text, t.str...)
	} else if t.script == 0 && t.region == 0 {
		text = append(text, t.lang.String()...)
	} else {
		buf := [maxCoreSize]byte{}
		text = buf[:t.genCoreBytes(buf[:])]
	}
	return text, nil
}


func (t *Tag) UnmarshalText(text []byte) error {
	tag, err := Raw.Parse(string(text))
	*t = tag
	return err
}




func (t Tag) Base() (Base, Confidence) {
	if t.lang != 0 {
		return Base{t.lang}, Exact
	}
	c := High
	if t.script == 0 && !(Region{t.region}).IsCountry() {
		c = Low
	}
	if tag, err := addTags(t); err == nil && tag.lang != 0 {
		return Base{tag.lang}, c
	}
	return Base{0}, No
}















func (t Tag) Script() (Script, Confidence) {
	if t.script != 0 {
		return Script{t.script}, Exact
	}
	sc, c := scriptID(_Zzzz), No
	if t.lang < langNoIndexOffset {
		if scr := scriptID(suppressScript[t.lang]); scr != 0 {
			
			
			if t.region == 0 {
				return Script{scriptID(scr)}, High
			}
			sc, c = scr, High
		}
	}
	if tag, err := addTags(t); err == nil {
		if tag.script != sc {
			sc, c = tag.script, Low
		}
	} else {
		t, _ = (Deprecated | Macro).Canonicalize(t)
		if tag, err := addTags(t); err == nil && tag.script != sc {
			sc, c = tag.script, Low
		}
	}
	return Script{sc}, c
}




func (t Tag) Region() (Region, Confidence) {
	if t.region != 0 {
		return Region{t.region}, Exact
	}
	if t, err := addTags(t); err == nil {
		return Region{t.region}, Low 
	}
	t, _ = (Deprecated | Macro).Canonicalize(t)
	if tag, err := addTags(t); err == nil {
		return Region{tag.region}, Low
	}
	return Region{_ZZ}, No 
}



func (t Tag) Variants() []Variant {
	v := []Variant{}
	if int(t.pVariant) < int(t.pExt) {
		for x, str := "", t.str[t.pVariant:t.pExt]; str != ""; {
			x, str = nextToken(str)
			v = append(v, Variant{x})
		}
	}
	return v
}




func (t Tag) Parent() Tag {
	if t.str != "" {
		
		t, _ = Raw.Compose(t.Raw())
		if t.region == 0 && t.script != 0 && t.lang != 0 {
			base, _ := addTags(Tag{lang: t.lang})
			if base.script == t.script {
				return Tag{lang: t.lang}
			}
		}
		return t
	}
	if t.lang != 0 {
		if t.region != 0 {
			maxScript := t.script
			if maxScript == 0 {
				max, _ := addTags(t)
				maxScript = max.script
			}

			for i := range parents {
				if langID(parents[i].lang) == t.lang && scriptID(parents[i].maxScript) == maxScript {
					for _, r := range parents[i].fromRegion {
						if regionID(r) == t.region {
							return Tag{
								lang:   t.lang,
								script: scriptID(parents[i].script),
								region: regionID(parents[i].toRegion),
							}
						}
					}
				}
			}

			
			base, _ := addTags(Tag{lang: t.lang})
			if base.script != maxScript {
				return Tag{lang: t.lang, script: maxScript}
			}
			return Tag{lang: t.lang}
		} else if t.script != 0 {
			
			
			base, _ := addTags(Tag{lang: t.lang})
			if base.script != t.script {
				return und
			}
			return Tag{lang: t.lang}
		}
	}
	return und
}


func nextToken(s string) (t, tail string) {
	p := strings.Index(s[1:], "-")
	if p == -1 {
		return s[1:], ""
	}
	p++
	return s[1:p], s[p:]
}


type Extension struct {
	s string
}



func (e Extension) String() string {
	return e.s
}


func ParseExtension(s string) (e Extension, err error) {
	scan := makeScannerString(s)
	var end int
	if n := len(scan.token); n != 1 {
		return Extension{}, errSyntax
	}
	scan.toLower(0, len(scan.b))
	end = parseExtension(&scan)
	if end != len(s) {
		return Extension{}, errSyntax
	}
	return Extension{string(scan.b)}, nil
}



func (e Extension) Type() byte {
	if e.s == "" {
		return 0
	}
	return e.s[0]
}


func (e Extension) Tokens() []string {
	return strings.Split(e.s, "-")
}




func (t Tag) Extension(x byte) (ext Extension, ok bool) {
	for i := int(t.pExt); i < len(t.str)-1; {
		var ext string
		i, ext = getExtension(t.str, i)
		if ext[0] == x {
			return Extension{ext}, true
		}
	}
	return Extension{}, false
}


func (t Tag) Extensions() []Extension {
	e := []Extension{}
	for i := int(t.pExt); i < len(t.str)-1; {
		var ext string
		i, ext = getExtension(t.str, i)
		e = append(e, Extension{ext})
	}
	return e
}





func (t Tag) TypeForKey(key string) string {
	if start, end, _ := t.findTypeForKey(key); end != start {
		return t.str[start:end]
	}
	return ""
}

var (
	errPrivateUse       = errors.New("cannot set a key on a private use tag")
	errInvalidArguments = errors.New("invalid key or type")
)





func (t Tag) SetTypeForKey(key, value string) (Tag, error) {
	if t.private() {
		return t, errPrivateUse
	}
	if len(key) != 2 {
		return t, errInvalidArguments
	}

	
	if value == "" {
		start, end, _ := t.findTypeForKey(key)
		if start != end {
			
			start -= 4

			
			if (end == len(t.str) || t.str[end+2] == '-') && t.str[start-2] == '-' {
				start -= 2
			}
			if start == int(t.pVariant) && end == len(t.str) {
				t.str = ""
				t.pVariant, t.pExt = 0, 0
			} else {
				t.str = fmt.Sprintf("%s%s", t.str[:start], t.str[end:])
			}
		}
		return t, nil
	}

	if len(value) < 3 || len(value) > 8 {
		return t, errInvalidArguments
	}

	var (
		buf    [maxCoreSize + maxSimpleUExtensionSize]byte
		uStart int 
	)

	
	if t.str == "" {
		uStart = t.genCoreBytes(buf[:])
		buf[uStart] = '-'
		uStart++
	}

	
	b := buf[uStart:]
	copy(b, "u-")
	copy(b[2:], key)
	b[4] = '-'
	b = b[:5+copy(b[5:], value)]
	scan := makeScanner(b)
	if parseExtensions(&scan); scan.err != nil {
		return t, scan.err
	}

	
	if t.str == "" {
		t.pVariant, t.pExt = byte(uStart-1), uint16(uStart-1)
		t.str = string(buf[:uStart+len(b)])
	} else {
		s := t.str
		start, end, hasExt := t.findTypeForKey(key)
		if start == end {
			if hasExt {
				b = b[2:]
			}
			t.str = fmt.Sprintf("%s-%s%s", s[:start], b, s[end:])
		} else {
			t.str = fmt.Sprintf("%s%s%s", s[:start], value, s[end:])
		}
	}
	return t, nil
}






func (t Tag) findTypeForKey(key string) (start, end int, hasExt bool) {
	p := int(t.pExt)
	if len(key) != 2 || p == len(t.str) || p == 0 {
		return p, p, false
	}
	s := t.str

	
	for p++; s[p] != 'u'; p++ {
		if s[p] > 'u' {
			p--
			return p, p, false
		}
		if p = nextExtension(s, p); p == len(s) {
			return len(s), len(s), false
		}
	}
	
	p++

	
	curKey := ""

	
	for {
		
		if p3 := p + 3; s[p3] == '-' {
			
			
			if curKey == key {
				return start, p, true
			}
			
			curKey = s[p+1 : p3]
			if curKey > key {
				return p, p, true
			}
			
			start = p + 4
			
			p += 7 
		} else {
			
			p += 4
		}
		
		max := p + 5 
		if len(s) < max {
			max = len(s)
		}
		for ; p < max && s[p] != '-'; p++ {
		}
		
		
		if p == len(s) || s[p+2] == '-' {
			if curKey == key {
				return start, p, true
			}
			return p, p, true
		}
	}
}






func CompactIndex(t Tag) (index int, ok bool) {
	
	
	
	b, s, r := t.Raw()
	if len(t.str) > 0 {
		if strings.HasPrefix(t.str, "x-") {
			
			return 0, false
		}
		if uint16(t.pVariant) != t.pExt {
			
			if t.TypeForKey("va") != "" {
				return 0, false
			}
			t, _ = Raw.Compose(b, s, r, t.Variants())
		} else if _, ok := t.Extension('u'); ok {
			
			variant := t.TypeForKey("va")
			t, _ = Raw.Compose(b, s, r)
			t, _ = t.SetTypeForKey("va", variant)
		}
		if len(t.str) > 0 {
			
			for i, s := range specialTags {
				if s == t {
					return i + 1, true
				}
			}
			return 0, false
		}
	}
	
	
	
	key := uint32(b.langID) << (8 + 12)
	key |= uint32(s.scriptID) << 12
	key |= uint32(r.regionID)
	x, ok := coreTags[key]
	return int(x), ok
}



type Base struct {
	langID
}




func ParseBase(s string) (Base, error) {
	if n := len(s); n < 2 || 3 < n {
		return Base{}, errSyntax
	}
	var buf [3]byte
	l, err := getLangID(buf[:copy(buf[:], s)])
	return Base{l}, err
}



type Script struct {
	scriptID
}




func ParseScript(s string) (Script, error) {
	if len(s) != 4 {
		return Script{}, errSyntax
	}
	var buf [4]byte
	sc, err := getScriptID(script, buf[:copy(buf[:], s)])
	return Script{sc}, err
}


type Region struct {
	regionID
}



func EncodeM49(r int) (Region, error) {
	rid, err := getRegionM49(r)
	return Region{rid}, err
}




func ParseRegion(s string) (Region, error) {
	if n := len(s); n < 2 || 3 < n {
		return Region{}, errSyntax
	}
	var buf [3]byte
	r, err := getRegionID(buf[:copy(buf[:], s)])
	return Region{r}, err
}



func (r Region) IsCountry() bool {
	if r.regionID == 0 || r.IsGroup() || r.IsPrivateUse() && r.regionID != _XK {
		return false
	}
	return true
}



func (r Region) IsGroup() bool {
	if r.regionID == 0 {
		return false
	}
	return int(regionInclusion[r.regionID]) < len(regionContainment)
}



func (r Region) Contains(c Region) bool {
	return r.regionID.contains(c.regionID)
}

func (r regionID) contains(c regionID) bool {
	if r == c {
		return true
	}
	g := regionInclusion[r]
	if g >= nRegionGroups {
		return false
	}
	m := regionContainment[g]

	d := regionInclusion[c]
	b := regionInclusionBits[d]

	
	
	
	if d >= nRegionGroups {
		return b&m != 0
	}
	return b&^m == 0
}

var errNoTLD = errors.New("language: region is not a valid ccTLD")








func (r Region) TLD() (Region, error) {
	
	
	if r.regionID == _GB {
		r = Region{_UK}
	}
	if (r.typ() & ccTLD) == 0 {
		return Region{}, errNoTLD
	}
	return r, nil
}




func (r Region) Canonicalize() Region {
	if cr := normRegion(r.regionID); cr != 0 {
		return Region{cr}
	}
	return r
}


type Variant struct {
	variant string
}



func ParseVariant(s string) (Variant, error) {
	s = strings.ToLower(s)
	if _, ok := variantIndex[s]; ok {
		return Variant{s}, nil
	}
	return Variant{}, mkErrInvalid([]byte(s))
}


func (v Variant) String() string {
	return v.variant
}
