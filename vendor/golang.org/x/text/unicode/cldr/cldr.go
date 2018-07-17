















package cldr 

import (
	"fmt"
	"sort"
)


type CLDR struct {
	parent   map[string][]string
	locale   map[string]*LDML
	resolved map[string]*LDML
	bcp47    *LDMLBCP47
	supp     *SupplementalData
}

func makeCLDR() *CLDR {
	return &CLDR{
		parent:   make(map[string][]string),
		locale:   make(map[string]*LDML),
		resolved: make(map[string]*LDML),
		bcp47:    &LDMLBCP47{},
		supp:     &SupplementalData{},
	}
}


func (cldr *CLDR) BCP47() *LDMLBCP47 {
	return nil
}


type Draft int

const (
	Approved Draft = iota
	Contributed
	Provisional
	Unconfirmed
)

var drafts = []string{"unconfirmed", "provisional", "contributed", "approved", ""}



func ParseDraft(level string) (Draft, error) {
	if level == "" {
		return Approved, nil
	}
	for i, s := range drafts {
		if level == s {
			return Unconfirmed - Draft(i), nil
		}
	}
	return Approved, fmt.Errorf("cldr: unknown draft level %q", level)
}

func (d Draft) String() string {
	return drafts[len(drafts)-1-int(d)]
}







func (cldr *CLDR) SetDraftLevel(lev Draft, preferDraft bool) {
	
	cldr.resolved = make(map[string]*LDML)
}



func (cldr *CLDR) RawLDML(loc string) *LDML {
	return cldr.locale[loc]
}



func (cldr *CLDR) LDML(loc string) (*LDML, error) {
	return cldr.resolve(loc)
}



func (cldr *CLDR) Supplemental() *SupplementalData {
	return cldr.supp
}




func (cldr *CLDR) Locales() []string {
	loc := []string{"root"}
	hasRoot := false
	for l, _ := range cldr.locale {
		if l == "root" {
			hasRoot = true
			continue
		}
		loc = append(loc, l)
	}
	sort.Strings(loc[1:])
	if !hasRoot {
		return loc[1:]
	}
	return loc
}


func Get(e Elem, path string) (res Elem, err error) {
	return walkXPath(e, path)
}
