



package language

import (
	"fmt"
	"sort"
)





type Coverage interface {
	
	Tags() []Tag

	
	BaseLanguages() []Base

	
	Scripts() []Script

	
	Regions() []Region
}

var (
	
	
	Supported Coverage = allSubtags{}
)






type allSubtags struct{}




func (s allSubtags) Regions() []Region {
	reg := make([]Region, numRegions)
	for i := range reg {
		reg[i] = Region{regionID(i + 1)}
	}
	return reg
}




func (s allSubtags) Scripts() []Script {
	scr := make([]Script, numScripts)
	for i := range scr {
		scr[i] = Script{scriptID(i + 1)}
	}
	return scr
}



func (s allSubtags) BaseLanguages() []Base {
	base := make([]Base, 0, numLanguages)
	for i := 0; i < langNoIndexOffset; i++ {
		
		if i != nonCanonicalUnd {
			base = append(base, Base{langID(i)})
		}
	}
	i := langNoIndexOffset
	for _, v := range langNoIndex {
		for k := 0; k < 8; k++ {
			if v&1 == 1 {
				base = append(base, Base{langID(i)})
			}
			v >>= 1
			i++
		}
	}
	return base
}


func (s allSubtags) Tags() []Tag {
	return nil
}






type coverage struct {
	tags    func() []Tag
	bases   func() []Base
	scripts func() []Script
	regions func() []Region
}

func (s *coverage) Tags() []Tag {
	if s.tags == nil {
		return nil
	}
	return s.tags()
}


type bases []Base

func (b bases) Len() int {
	return len(b)
}

func (b bases) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b bases) Less(i, j int) bool {
	return b[i].langID < b[j].langID
}



func (s *coverage) BaseLanguages() []Base {
	if s.bases == nil {
		tags := s.Tags()
		if len(tags) == 0 {
			return nil
		}
		a := make([]Base, len(tags))
		for i, t := range tags {
			a[i] = Base{langID(t.lang)}
		}
		sort.Sort(bases(a))
		k := 0
		for i := 1; i < len(a); i++ {
			if a[k] != a[i] {
				k++
				a[k] = a[i]
			}
		}
		return a[:k+1]
	}
	return s.bases()
}

func (s *coverage) Scripts() []Script {
	if s.scripts == nil {
		return nil
	}
	return s.scripts()
}

func (s *coverage) Regions() []Region {
	if s.regions == nil {
		return nil
	}
	return s.regions()
}







func NewCoverage(list ...interface{}) Coverage {
	s := &coverage{}
	for _, x := range list {
		switch v := x.(type) {
		case func() []Base:
			s.bases = v
		case func() []Script:
			s.scripts = v
		case func() []Region:
			s.regions = v
		case func() []Tag:
			s.tags = v
		case []Base:
			s.bases = func() []Base { return v }
		case []Script:
			s.scripts = func() []Script { return v }
		case []Region:
			s.regions = func() []Region { return v }
		case []Tag:
			s.tags = func() []Tag { return v }
		default:
			panic(fmt.Sprintf("language: unsupported set type %T", v))
		}
	}
	return s
}
