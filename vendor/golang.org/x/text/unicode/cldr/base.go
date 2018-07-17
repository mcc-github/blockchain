



package cldr

import (
	"encoding/xml"
	"regexp"
	"strconv"
)


type Elem interface {
	setEnclosing(Elem)
	setName(string)
	enclosing() Elem

	GetCommon() *Common
}

type hidden struct {
	CharData string `xml:",chardata"`
	Alias    *struct {
		Common
		Source string `xml:"source,attr"`
		Path   string `xml:"path,attr"`
	} `xml:"alias"`
	Def *struct {
		Common
		Choice string `xml:"choice,attr,omitempty"`
		Type   string `xml:"type,attr,omitempty"`
	} `xml:"default"`
}



type Common struct {
	XMLName         xml.Name
	name            string
	enclElem        Elem
	Type            string `xml:"type,attr,omitempty"`
	Reference       string `xml:"reference,attr,omitempty"`
	Alt             string `xml:"alt,attr,omitempty"`
	ValidSubLocales string `xml:"validSubLocales,attr,omitempty"`
	Draft           string `xml:"draft,attr,omitempty"`
	hidden
}



func (e *Common) Default() string {
	if e.Def == nil {
		return ""
	}
	if e.Def.Choice != "" {
		return e.Def.Choice
	} else if e.Def.Type != "" {
		
		return e.Def.Type
	}
	return ""
}


func (e *Common) Element() string {
	return e.name
}


func (e *Common) GetCommon() *Common {
	return e
}


func (e *Common) Data() string {
	e.CharData = charRe.ReplaceAllStringFunc(e.CharData, replaceUnicode)
	return e.CharData
}

func (e *Common) setName(s string) {
	e.name = s
}

func (e *Common) enclosing() Elem {
	return e.enclElem
}

func (e *Common) setEnclosing(en Elem) {
	e.enclElem = en
}


var charRe = regexp.MustCompile(`&#x[0-9a-fA-F]*;|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8}|\\x[0-9a-fA-F]{2}|\\[0-7]{3}|\\[abtnvfr]`)



func replaceUnicode(s string) string {
	if s[1] == '#' {
		r, _ := strconv.ParseInt(s[3:len(s)-1], 16, 32)
		return string(r)
	}
	r, _, _, _ := strconv.UnquoteChar(s, 0)
	return string(r)
}
