



package language

import "errors"


type MatchOption func(*matcher)




func PreferSameScript(preferSame bool) MatchOption {
	return func(m *matcher) { m.preferSameScript = preferSame }
}










func MatchStrings(m Matcher, lang ...string) (tag Tag, index int) {
	for _, accept := range lang {
		desired, _, err := ParseAcceptLanguage(accept)
		if err != nil {
			continue
		}
		if tag, index, conf := m.Match(desired...); conf != No {
			return tag, index
		}
	}
	tag, index, _ = m.Match()
	return
}






type Matcher interface {
	Match(t ...Tag) (tag Tag, index int, c Confidence)
}



func Comprehends(speaker, alternative Tag) Confidence {
	_, _, c := NewMatcher([]Tag{alternative}).Match(speaker)
	return c
}















func NewMatcher(t []Tag, options ...MatchOption) Matcher {
	return newMatcher(t, options)
}

func (m *matcher) Match(want ...Tag) (t Tag, index int, c Confidence) {
	match, w, c := m.getBest(want...)
	if match != nil {
		t, index = match.tag, match.index
	} else {
		
		t = m.default_.tag
		if m.preferSameScript {
		outer:
			for _, w := range want {
				script, _ := w.Script()
				if script.scriptID == 0 {
					
					
					continue
				}
				for i, h := range m.supported {
					if script.scriptID == h.maxScript {
						t, index = h.tag, i
						break outer
					}
				}
			}
		}
		
	}
	if w.region != 0 && t.region != 0 && t.region.contains(w.region) {
		t, _ = Raw.Compose(t, Region{w.region})
	}
	
	
	
	
	if e := w.Extensions(); len(e) > 0 {
		t, _ = Raw.Compose(t, e)
	}
	return t, index, c
}

type scriptRegionFlags uint8

const (
	isList = 1 << iota
	scriptInFrom
	regionInFrom
)

func (t *Tag) setUndefinedLang(id langID) {
	if t.lang == 0 {
		t.lang = id
	}
}

func (t *Tag) setUndefinedScript(id scriptID) {
	if t.script == 0 {
		t.script = id
	}
}

func (t *Tag) setUndefinedRegion(id regionID) {
	if t.region == 0 || t.region.contains(id) {
		t.region = id
	}
}



var ErrMissingLikelyTagsData = errors.New("missing likely tags data")





func (t Tag) addLikelySubtags() (Tag, error) {
	id, err := addTags(t)
	if err != nil {
		return t, err
	} else if id.equalTags(t) {
		return t, nil
	}
	id.remakeString()
	return id, nil
}


func specializeRegion(t *Tag) bool {
	if i := regionInclusion[t.region]; i < nRegionGroups {
		x := likelyRegionGroup[i]
		if langID(x.lang) == t.lang && scriptID(x.script) == t.script {
			t.region = regionID(x.region)
		}
		return true
	}
	return false
}

func addTags(t Tag) (Tag, error) {
	
	if t.private() {
		return t, nil
	}
	if t.script != 0 && t.region != 0 {
		if t.lang != 0 {
			
			specializeRegion(&t)
			return t, nil
		}
		
		
		list := likelyRegion[t.region : t.region+1]
		if x := list[0]; x.flags&isList != 0 {
			list = likelyRegionList[x.lang : x.lang+uint16(x.script)]
		}
		for _, x := range list {
			
			if scriptID(x.script) == t.script {
				t.setUndefinedLang(langID(x.lang))
				return t, nil
			}
		}
	}
	if t.lang != 0 {
		
		if t.lang < langNoIndexOffset {
			x := likelyLang[t.lang]
			if x.flags&isList != 0 {
				list := likelyLangList[x.region : x.region+uint16(x.script)]
				if t.script != 0 {
					for _, x := range list {
						if scriptID(x.script) == t.script && x.flags&scriptInFrom != 0 {
							t.setUndefinedRegion(regionID(x.region))
							return t, nil
						}
					}
				} else if t.region != 0 {
					count := 0
					goodScript := true
					tt := t
					for _, x := range list {
						
						
						
						
						if x.flags&scriptInFrom == 0 && t.region.contains(regionID(x.region)) {
							tt.region = regionID(x.region)
							tt.setUndefinedScript(scriptID(x.script))
							goodScript = goodScript && tt.script == scriptID(x.script)
							count++
						}
					}
					if count == 1 {
						return tt, nil
					}
					
					
					if goodScript {
						t.script = tt.script
					}
				}
			}
		}
	} else {
		
		if t.script != 0 {
			x := likelyScript[t.script]
			if x.region != 0 {
				t.setUndefinedRegion(regionID(x.region))
				t.setUndefinedLang(langID(x.lang))
				return t, nil
			}
		}
		
		
		if t.region != 0 {
			if i := regionInclusion[t.region]; i < nRegionGroups {
				x := likelyRegionGroup[i]
				if x.region != 0 {
					t.setUndefinedLang(langID(x.lang))
					t.setUndefinedScript(scriptID(x.script))
					t.region = regionID(x.region)
				}
			} else {
				x := likelyRegion[t.region]
				if x.flags&isList != 0 {
					x = likelyRegionList[x.lang]
				}
				if x.script != 0 && x.flags != scriptInFrom {
					t.setUndefinedLang(langID(x.lang))
					t.setUndefinedScript(scriptID(x.script))
					return t, nil
				}
			}
		}
	}

	
	if t.lang < langNoIndexOffset {
		x := likelyLang[t.lang]
		if x.flags&isList != 0 {
			x = likelyLangList[x.region]
		}
		if x.region != 0 {
			t.setUndefinedScript(scriptID(x.script))
			t.setUndefinedRegion(regionID(x.region))
		}
		specializeRegion(&t)
		if t.lang == 0 {
			t.lang = _en 
		}
		return t, nil
	}
	return t, ErrMissingLikelyTagsData
}

func (t *Tag) setTagsFrom(id Tag) {
	t.lang = id.lang
	t.script = id.script
	t.region = id.region
}



func (t Tag) minimize() (Tag, error) {
	t, err := minimizeTags(t)
	if err != nil {
		return t, err
	}
	t.remakeString()
	return t, nil
}


func minimizeTags(t Tag) (Tag, error) {
	if t.equalTags(und) {
		return t, nil
	}
	max, err := addTags(t)
	if err != nil {
		return t, err
	}
	for _, id := range [...]Tag{
		{lang: t.lang},
		{lang: t.lang, region: t.region},
		{lang: t.lang, script: t.script},
	} {
		if x, err := addTags(id); err == nil && max.equalTags(x) {
			t.setTagsFrom(id)
			break
		}
	}
	return t, nil
}

















































































































type matcher struct {
	default_         *haveTag
	supported        []*haveTag
	index            map[langID]*matchHeader
	passSettings     bool
	preferSameScript bool
}



type matchHeader struct {
	haveTags []*haveTag
	original bool
}



type haveTag struct {
	tag Tag

	
	index int

	
	
	conf Confidence

	
	maxRegion regionID
	maxScript scriptID

	
	
	
	altScript scriptID

	
	nextMax uint16
}

func makeHaveTag(tag Tag, index int) (haveTag, langID) {
	max := tag
	if tag.lang != 0 || tag.region != 0 || tag.script != 0 {
		max, _ = max.canonicalize(All)
		max, _ = addTags(max)
		max.remakeString()
	}
	return haveTag{tag, index, Exact, max.region, max.script, altScript(max.lang, max.script), 0}, max.lang
}




func altScript(l langID, s scriptID) scriptID {
	for _, alt := range matchScript {
		
		if (langID(alt.wantLang) == l || langID(alt.haveLang) == l) &&
			scriptID(alt.haveScript) == s {
			return scriptID(alt.wantScript)
		}
	}
	return 0
}



func (h *matchHeader) addIfNew(n haveTag, exact bool) {
	h.original = h.original || exact
	
	for _, v := range h.haveTags {
		if v.tag.equalsRest(n.tag) {
			return
		}
	}
	
	
	for i, v := range h.haveTags {
		if v.maxScript == n.maxScript &&
			v.maxRegion == n.maxRegion &&
			v.tag.variantOrPrivateTagStr() == n.tag.variantOrPrivateTagStr() {
			for h.haveTags[i].nextMax != 0 {
				i = int(h.haveTags[i].nextMax)
			}
			h.haveTags[i].nextMax = uint16(len(h.haveTags))
			break
		}
	}
	h.haveTags = append(h.haveTags, &n)
}



func (m *matcher) header(l langID) *matchHeader {
	if h := m.index[l]; h != nil {
		return h
	}
	h := &matchHeader{}
	m.index[l] = h
	return h
}

func toConf(d uint8) Confidence {
	if d <= 10 {
		return High
	}
	if d < 30 {
		return Low
	}
	return No
}




func newMatcher(supported []Tag, options []MatchOption) *matcher {
	m := &matcher{
		index:            make(map[langID]*matchHeader),
		preferSameScript: true,
	}
	for _, o := range options {
		o(m)
	}
	if len(supported) == 0 {
		m.default_ = &haveTag{}
		return m
	}
	
	
	for i, tag := range supported {
		pair, _ := makeHaveTag(tag, i)
		m.header(tag.lang).addIfNew(pair, true)
		m.supported = append(m.supported, &pair)
	}
	m.default_ = m.header(supported[0].lang).haveTags[0]
	
	
	for i, tag := range supported {
		pair, max := makeHaveTag(tag, i)
		if max != tag.lang {
			m.header(max).addIfNew(pair, true)
		}
	}

	
	
	
	update := func(want, have uint16, conf Confidence) {
		if hh := m.index[langID(have)]; hh != nil {
			if !hh.original {
				return
			}
			hw := m.header(langID(want))
			for _, ht := range hh.haveTags {
				v := *ht
				if conf < v.conf {
					v.conf = conf
				}
				v.nextMax = 0 
				if v.altScript != 0 {
					v.altScript = altScript(langID(want), v.maxScript)
				}
				hw.addIfNew(v, conf == Exact && hh.original)
			}
		}
	}

	
	
	for _, ml := range matchLang {
		update(ml.want, ml.have, toConf(ml.distance))
		if !ml.oneway {
			update(ml.have, ml.want, toConf(ml.distance))
		}
	}

	
	
	
	
	
	for i, lm := range langAliasMap {
		
		
		conf := Exact
		if langAliasTypes[i] != langMacro {
			if !isExactEquivalent(langID(lm.from)) {
				conf = High
			}
			update(lm.to, lm.from, conf)
		}
		update(lm.from, lm.to, conf)
	}
	return m
}



func (m *matcher) getBest(want ...Tag) (got *haveTag, orig Tag, c Confidence) {
	best := bestMatch{}
	for i, w := range want {
		var max Tag
		
		h := m.index[w.lang]
		if w.lang != 0 {
			if h == nil {
				continue
			}
			
			max, _ = w.canonicalize(Legacy | Deprecated | Macro)
			
			
			if w.region != max.region {
				w.region = max.region
			}
			
			
			max, _ = addTags(max)
		} else {
			
			if h != nil {
				for i := range h.haveTags {
					have := h.haveTags[i]
					if have.tag.equalsRest(w) {
						return have, w, Exact
					}
				}
			}
			if w.script == 0 && w.region == 0 {
				
				
				continue
			}
			max, _ = addTags(w)
			if h = m.index[max.lang]; h == nil {
				continue
			}
		}
		pin := true
		for _, t := range want[i+1:] {
			if w.lang == t.lang {
				pin = false
				break
			}
		}
		
		for i := range h.haveTags {
			have := h.haveTags[i]
			best.update(have, w, max.script, max.region, pin)
			if best.conf == Exact {
				for have.nextMax != 0 {
					have = h.haveTags[have.nextMax]
					best.update(have, w, max.script, max.region, pin)
				}
				return best.have, best.want, best.conf
			}
		}
	}
	if best.conf <= No {
		if len(want) != 0 {
			return nil, want[0], No
		}
		return nil, Tag{}, No
	}
	return best.have, best.want, best.conf
}


type bestMatch struct {
	have            *haveTag
	want            Tag
	conf            Confidence
	pinnedRegion    regionID
	pinLanguage     bool
	sameRegionGroup bool
	
	origLang     bool
	origReg      bool
	paradigmReg  bool
	regGroupDist uint8
	origScript   bool
}















func (m *bestMatch) update(have *haveTag, tag Tag, maxScript scriptID, maxRegion regionID, pin bool) {
	
	c := have.conf
	if c < m.conf {
		return
	}
	
	if m.pinLanguage && tag.lang != m.want.lang {
		return
	}
	
	if tag.lang == m.want.lang && m.sameRegionGroup {
		_, sameGroup := regionGroupDist(m.pinnedRegion, have.maxRegion, have.maxScript, m.want.lang)
		if !sameGroup {
			return
		}
	}
	if c == Exact && have.maxScript == maxScript {
		
		
		m.pinLanguage = pin
	}
	if have.tag.equalsRest(tag) {
	} else if have.maxScript != maxScript {
		
		
		
		if Low < m.conf || have.altScript != maxScript {
			return
		}
		c = Low
	} else if have.maxRegion != maxRegion {
		if High < c {
			
			c = High
		}
	}

	
	
	
	
	beaten := false 
	if c != m.conf {
		if c < m.conf {
			return
		}
		beaten = true
	}

	
	
	origLang := have.tag.lang == tag.lang && tag.lang != 0
	if !beaten && m.origLang != origLang {
		if m.origLang {
			return
		}
		beaten = true
	}

	
	origReg := have.tag.region == tag.region && tag.region != 0
	if !beaten && m.origReg != origReg {
		if m.origReg {
			return
		}
		beaten = true
	}

	regGroupDist, sameGroup := regionGroupDist(have.maxRegion, maxRegion, maxScript, tag.lang)
	if !beaten && m.regGroupDist != regGroupDist {
		if regGroupDist > m.regGroupDist {
			return
		}
		beaten = true
	}

	paradigmReg := isParadigmLocale(tag.lang, have.maxRegion)
	if !beaten && m.paradigmReg != paradigmReg {
		if !paradigmReg {
			return
		}
		beaten = true
	}

	
	origScript := have.tag.script == tag.script && tag.script != 0
	if !beaten && m.origScript != origScript {
		if m.origScript {
			return
		}
		beaten = true
	}

	
	if beaten {
		m.have = have
		m.want = tag
		m.conf = c
		m.pinnedRegion = maxRegion
		m.sameRegionGroup = sameGroup
		m.origLang = origLang
		m.origReg = origReg
		m.paradigmReg = paradigmReg
		m.origScript = origScript
		m.regGroupDist = regGroupDist
	}
}

func isParadigmLocale(lang langID, r regionID) bool {
	for _, e := range paradigmLocales {
		if langID(e[0]) == lang && (r == regionID(e[1]) || r == regionID(e[2])) {
			return true
		}
	}
	return false
}



func regionGroupDist(a, b regionID, script scriptID, lang langID) (dist uint8, same bool) {
	const defaultDistance = 4

	aGroup := uint(regionToGroups[a]) << 1
	bGroup := uint(regionToGroups[b]) << 1
	for _, ri := range matchRegion {
		if langID(ri.lang) == lang && (ri.script == 0 || scriptID(ri.script) == script) {
			group := uint(1 << (ri.group &^ 0x80))
			if 0x80&ri.group == 0 {
				if aGroup&bGroup&group != 0 { 
					return ri.distance, ri.distance == defaultDistance
				}
			} else {
				if (aGroup|bGroup)&group == 0 { 
					return ri.distance, ri.distance == defaultDistance
				}
			}
		}
	}
	return defaultDistance, true
}

func (t Tag) variants() string {
	if t.pVariant == 0 {
		return ""
	}
	return t.str[t.pVariant:t.pExt]
}


func (t Tag) variantOrPrivateTagStr() string {
	if t.pExt > 0 {
		return t.str[t.pVariant:t.pExt]
	}
	return t.str[t.pVariant:]
}


func (a Tag) equalsRest(b Tag) bool {
	
	
	return a.script == b.script && a.region == b.region && a.variantOrPrivateTagStr() == b.variantOrPrivateTagStr()
}



func isExactEquivalent(l langID) bool {
	for _, o := range notEquivalent {
		if o == l {
			return false
		}
	}
	return true
}

var notEquivalent []langID

func init() {
	
	
	for _, lm := range langAliasMap {
		tag := Tag{lang: langID(lm.from)}
		if tag, _ = tag.canonicalize(All); tag.script != 0 || tag.region != 0 {
			notEquivalent = append(notEquivalent, langID(lm.from))
		}
	}
	
	for i, v := range paradigmLocales {
		max, _ := addTags(Tag{lang: langID(v[0])})
		if v[1] == 0 {
			paradigmLocales[i][1] = uint16(max.region)
		}
		if v[2] == 0 {
			paradigmLocales[i][2] = uint16(max.region)
		}
	}
}
