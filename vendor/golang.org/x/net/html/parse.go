



package html

import (
	"errors"
	"fmt"
	"io"
	"strings"

	a "golang.org/x/net/html/atom"
)



type parser struct {
	
	tokenizer *Tokenizer
	
	tok Token
	
	
	hasSelfClosingToken bool
	
	doc *Node
	
	
	oe, afe nodeStack
	
	head, form *Node
	
	scripting, framesetOK bool
	
	templateStack insertionModeStack
	
	im insertionMode
	
	
	originalIM insertionMode
	
	
	fosterParenting bool
	
	quirks bool
	
	fragment bool
	
	
	context *Node
}

func (p *parser) top() *Node {
	if n := p.oe.top(); n != nil {
		return n
	}
	return p.doc
}


var (
	defaultScopeStopTags = map[string][]a.Atom{
		"":     {a.Applet, a.Caption, a.Html, a.Table, a.Td, a.Th, a.Marquee, a.Object, a.Template},
		"math": {a.AnnotationXml, a.Mi, a.Mn, a.Mo, a.Ms, a.Mtext},
		"svg":  {a.Desc, a.ForeignObject, a.Title},
	}
)

type scope int

const (
	defaultScope scope = iota
	listItemScope
	buttonScope
	tableScope
	tableRowScope
	tableBodyScope
	selectScope
)


















func (p *parser) popUntil(s scope, matchTags ...a.Atom) bool {
	if i := p.indexOfElementInScope(s, matchTags...); i != -1 {
		p.oe = p.oe[:i]
		return true
	}
	return false
}




func (p *parser) indexOfElementInScope(s scope, matchTags ...a.Atom) int {
	for i := len(p.oe) - 1; i >= 0; i-- {
		tagAtom := p.oe[i].DataAtom
		if p.oe[i].Namespace == "" {
			for _, t := range matchTags {
				if t == tagAtom {
					return i
				}
			}
			switch s {
			case defaultScope:
				
			case listItemScope:
				if tagAtom == a.Ol || tagAtom == a.Ul {
					return -1
				}
			case buttonScope:
				if tagAtom == a.Button {
					return -1
				}
			case tableScope:
				if tagAtom == a.Html || tagAtom == a.Table || tagAtom == a.Template {
					return -1
				}
			case selectScope:
				if tagAtom != a.Optgroup && tagAtom != a.Option {
					return -1
				}
			default:
				panic("unreachable")
			}
		}
		switch s {
		case defaultScope, listItemScope, buttonScope:
			for _, t := range defaultScopeStopTags[p.oe[i].Namespace] {
				if t == tagAtom {
					return -1
				}
			}
		}
	}
	return -1
}



func (p *parser) elementInScope(s scope, matchTags ...a.Atom) bool {
	return p.indexOfElementInScope(s, matchTags...) != -1
}



func (p *parser) clearStackToContext(s scope) {
	for i := len(p.oe) - 1; i >= 0; i-- {
		tagAtom := p.oe[i].DataAtom
		switch s {
		case tableScope:
			if tagAtom == a.Html || tagAtom == a.Table || tagAtom == a.Template {
				p.oe = p.oe[:i+1]
				return
			}
		case tableRowScope:
			if tagAtom == a.Html || tagAtom == a.Tr || tagAtom == a.Template {
				p.oe = p.oe[:i+1]
				return
			}
		case tableBodyScope:
			if tagAtom == a.Html || tagAtom == a.Tbody || tagAtom == a.Tfoot || tagAtom == a.Thead || tagAtom == a.Template {
				p.oe = p.oe[:i+1]
				return
			}
		default:
			panic("unreachable")
		}
	}
}




func (p *parser) generateImpliedEndTags(exceptions ...string) {
	var i int
loop:
	for i = len(p.oe) - 1; i >= 0; i-- {
		n := p.oe[i]
		if n.Type == ElementNode {
			switch n.DataAtom {
			case a.Dd, a.Dt, a.Li, a.Optgroup, a.Option, a.P, a.Rb, a.Rp, a.Rt, a.Rtc:
				for _, except := range exceptions {
					if n.Data == except {
						break loop
					}
				}
				continue
			}
		}
		break
	}

	p.oe = p.oe[:i+1]
}



func (p *parser) addChild(n *Node) {
	if p.shouldFosterParent() {
		p.fosterParent(n)
	} else {
		p.top().AppendChild(n)
	}

	if n.Type == ElementNode {
		p.oe = append(p.oe, n)
	}
}



func (p *parser) shouldFosterParent() bool {
	if p.fosterParenting {
		switch p.top().DataAtom {
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			return true
		}
	}
	return false
}



func (p *parser) fosterParent(n *Node) {
	var table, parent, prev, template *Node
	var i int
	for i = len(p.oe) - 1; i >= 0; i-- {
		if p.oe[i].DataAtom == a.Table {
			table = p.oe[i]
			break
		}
	}

	var j int
	for j = len(p.oe) - 1; j >= 0; j-- {
		if p.oe[j].DataAtom == a.Template {
			template = p.oe[j]
			break
		}
	}

	if template != nil && (table == nil || j > i) {
		template.AppendChild(n)
		return
	}

	if table == nil {
		
		parent = p.oe[0]
	} else {
		parent = table.Parent
	}
	if parent == nil {
		parent = p.oe[i-1]
	}

	if table != nil {
		prev = table.PrevSibling
	} else {
		prev = parent.LastChild
	}
	if prev != nil && prev.Type == TextNode && n.Type == TextNode {
		prev.Data += n.Data
		return
	}

	parent.InsertBefore(n, table)
}



func (p *parser) addText(text string) {
	if text == "" {
		return
	}

	if p.shouldFosterParent() {
		p.fosterParent(&Node{
			Type: TextNode,
			Data: text,
		})
		return
	}

	t := p.top()
	if n := t.LastChild; n != nil && n.Type == TextNode {
		n.Data += text
		return
	}
	p.addChild(&Node{
		Type: TextNode,
		Data: text,
	})
}


func (p *parser) addElement() {
	p.addChild(&Node{
		Type:     ElementNode,
		DataAtom: p.tok.DataAtom,
		Data:     p.tok.Data,
		Attr:     p.tok.Attr,
	})
}


func (p *parser) addFormattingElement() {
	tagAtom, attr := p.tok.DataAtom, p.tok.Attr
	p.addElement()

	
	identicalElements := 0
findIdenticalElements:
	for i := len(p.afe) - 1; i >= 0; i-- {
		n := p.afe[i]
		if n.Type == scopeMarkerNode {
			break
		}
		if n.Type != ElementNode {
			continue
		}
		if n.Namespace != "" {
			continue
		}
		if n.DataAtom != tagAtom {
			continue
		}
		if len(n.Attr) != len(attr) {
			continue
		}
	compareAttributes:
		for _, t0 := range n.Attr {
			for _, t1 := range attr {
				if t0.Key == t1.Key && t0.Namespace == t1.Namespace && t0.Val == t1.Val {
					
					continue compareAttributes
				}
			}
			
			
			continue findIdenticalElements
		}

		identicalElements++
		if identicalElements >= 3 {
			p.afe.remove(n)
		}
	}

	p.afe = append(p.afe, p.top())
}


func (p *parser) clearActiveFormattingElements() {
	for {
		n := p.afe.pop()
		if len(p.afe) == 0 || n.Type == scopeMarkerNode {
			return
		}
	}
}


func (p *parser) reconstructActiveFormattingElements() {
	n := p.afe.top()
	if n == nil {
		return
	}
	if n.Type == scopeMarkerNode || p.oe.index(n) != -1 {
		return
	}
	i := len(p.afe) - 1
	for n.Type != scopeMarkerNode && p.oe.index(n) == -1 {
		if i == 0 {
			i = -1
			break
		}
		i--
		n = p.afe[i]
	}
	for {
		i++
		clone := p.afe[i].clone()
		p.addChild(clone)
		p.afe[i] = clone
		if i == len(p.afe)-1 {
			break
		}
	}
}


func (p *parser) acknowledgeSelfClosingTag() {
	p.hasSelfClosingToken = false
}





type insertionMode func(*parser) bool




func (p *parser) setOriginalIM() {
	if p.originalIM != nil {
		panic("html: bad parser state: originalIM was set twice")
	}
	p.originalIM = p.im
}


func (p *parser) resetInsertionMode() {
	for i := len(p.oe) - 1; i >= 0; i-- {
		n := p.oe[i]
		last := i == 0
		if last && p.context != nil {
			n = p.context
		}

		switch n.DataAtom {
		case a.Select:
			if !last {
				for ancestor, first := n, p.oe[0]; ancestor != first; {
					if ancestor == first {
						break
					}
					ancestor = p.oe[p.oe.index(ancestor)-1]
					switch ancestor.DataAtom {
					case a.Template:
						p.im = inSelectIM
						return
					case a.Table:
						p.im = inSelectInTableIM
						return
					}
				}
			}
			p.im = inSelectIM
		case a.Td, a.Th:
			
			
			
			p.im = inCellIM
		case a.Tr:
			p.im = inRowIM
		case a.Tbody, a.Thead, a.Tfoot:
			p.im = inTableBodyIM
		case a.Caption:
			p.im = inCaptionIM
		case a.Colgroup:
			p.im = inColumnGroupIM
		case a.Table:
			p.im = inTableIM
		case a.Template:
			
			if n.Namespace != "" {
				continue
			}
			p.im = p.templateStack.top()
		case a.Head:
			
			
			
			p.im = inHeadIM
		case a.Body:
			p.im = inBodyIM
		case a.Frameset:
			p.im = inFramesetIM
		case a.Html:
			if p.head == nil {
				p.im = beforeHeadIM
			} else {
				p.im = afterHeadIM
			}
		default:
			if last {
				p.im = inBodyIM
				return
			}
			continue
		}
		return
	}
}

const whitespace = " \t\r\n\f"


func initialIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		p.tok.Data = strings.TrimLeft(p.tok.Data, whitespace)
		if len(p.tok.Data) == 0 {
			
			return true
		}
	case CommentToken:
		p.doc.AppendChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	case DoctypeToken:
		n, quirks := parseDoctype(p.tok.Data)
		p.doc.AppendChild(n)
		p.quirks = quirks
		p.im = beforeHTMLIM
		return true
	}
	p.quirks = true
	p.im = beforeHTMLIM
	return false
}


func beforeHTMLIM(p *parser) bool {
	switch p.tok.Type {
	case DoctypeToken:
		
		return true
	case TextToken:
		p.tok.Data = strings.TrimLeft(p.tok.Data, whitespace)
		if len(p.tok.Data) == 0 {
			
			return true
		}
	case StartTagToken:
		if p.tok.DataAtom == a.Html {
			p.addElement()
			p.im = beforeHeadIM
			return true
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Head, a.Body, a.Html, a.Br:
			p.parseImpliedToken(StartTagToken, a.Html, a.Html.String())
			return false
		default:
			
			return true
		}
	case CommentToken:
		p.doc.AppendChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	}
	p.parseImpliedToken(StartTagToken, a.Html, a.Html.String())
	return false
}


func beforeHeadIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		p.tok.Data = strings.TrimLeft(p.tok.Data, whitespace)
		if len(p.tok.Data) == 0 {
			
			return true
		}
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Head:
			p.addElement()
			p.head = p.top()
			p.im = inHeadIM
			return true
		case a.Html:
			return inBodyIM(p)
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Head, a.Body, a.Html, a.Br:
			p.parseImpliedToken(StartTagToken, a.Head, a.Head.String())
			return false
		default:
			
			return true
		}
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	case DoctypeToken:
		
		return true
	}

	p.parseImpliedToken(StartTagToken, a.Head, a.Head.String())
	return false
}


func inHeadIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) < len(p.tok.Data) {
			
			p.addText(p.tok.Data[:len(p.tok.Data)-len(s)])
			if s == "" {
				return true
			}
			p.tok.Data = s
		}
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Base, a.Basefont, a.Bgsound, a.Command, a.Link, a.Meta:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			return true
		case a.Script, a.Title, a.Noscript, a.Noframes, a.Style:
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
			return true
		case a.Head:
			
			return true
		case a.Template:
			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.framesetOK = false
			p.im = inTemplateIM
			p.templateStack = append(p.templateStack, inTemplateIM)
			return true
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Head:
			p.oe.pop()
			p.im = afterHeadIM
			return true
		case a.Body, a.Html, a.Br:
			p.parseImpliedToken(EndTagToken, a.Head, a.Head.String())
			return false
		case a.Template:
			if !p.oe.contains(a.Template) {
				return true
			}
			
			
			
			p.generateImpliedEndTags()
			for i := len(p.oe) - 1; i >= 0; i-- {
				if n := p.oe[i]; n.Namespace == "" && n.DataAtom == a.Template {
					p.oe = p.oe[:i]
					break
				}
			}
			p.clearActiveFormattingElements()
			p.templateStack.pop()
			p.resetInsertionMode()
			return true
		default:
			
			return true
		}
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	case DoctypeToken:
		
		return true
	}

	p.parseImpliedToken(EndTagToken, a.Head, a.Head.String())
	return false
}


func afterHeadIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) < len(p.tok.Data) {
			
			p.addText(p.tok.Data[:len(p.tok.Data)-len(s)])
			if s == "" {
				return true
			}
			p.tok.Data = s
		}
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Body:
			p.addElement()
			p.framesetOK = false
			p.im = inBodyIM
			return true
		case a.Frameset:
			p.addElement()
			p.im = inFramesetIM
			return true
		case a.Base, a.Basefont, a.Bgsound, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			p.oe = append(p.oe, p.head)
			defer p.oe.remove(p.head)
			return inHeadIM(p)
		case a.Head:
			
			return true
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Body, a.Html, a.Br:
			
		case a.Template:
			return inHeadIM(p)
		default:
			
			return true
		}
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	case DoctypeToken:
		
		return true
	}

	p.parseImpliedToken(StartTagToken, a.Body, a.Body.String())
	p.framesetOK = true
	return false
}


func copyAttributes(dst *Node, src Token) {
	if len(src.Attr) == 0 {
		return
	}
	attr := map[string]string{}
	for _, t := range dst.Attr {
		attr[t.Key] = t.Val
	}
	for _, t := range src.Attr {
		if _, ok := attr[t.Key]; !ok {
			dst.Attr = append(dst.Attr, t)
			attr[t.Key] = t.Val
		}
	}
}


func inBodyIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		d := p.tok.Data
		switch n := p.oe.top(); n.DataAtom {
		case a.Pre, a.Listing:
			if n.FirstChild == nil {
				
				if d != "" && d[0] == '\r' {
					d = d[1:]
				}
				if d != "" && d[0] == '\n' {
					d = d[1:]
				}
			}
		}
		d = strings.Replace(d, "\x00", "", -1)
		if d == "" {
			return true
		}
		p.reconstructActiveFormattingElements()
		p.addText(d)
		if p.framesetOK && strings.TrimLeft(d, whitespace) != "" {
			
			p.framesetOK = false
		}
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			if p.oe.contains(a.Template) {
				return true
			}
			copyAttributes(p.oe[0], p.tok)
		case a.Base, a.Basefont, a.Bgsound, a.Command, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			return inHeadIM(p)
		case a.Body:
			if p.oe.contains(a.Template) {
				return true
			}
			if len(p.oe) >= 2 {
				body := p.oe[1]
				if body.Type == ElementNode && body.DataAtom == a.Body {
					p.framesetOK = false
					copyAttributes(body, p.tok)
				}
			}
		case a.Frameset:
			if !p.framesetOK || len(p.oe) < 2 || p.oe[1].DataAtom != a.Body {
				
				return true
			}
			body := p.oe[1]
			if body.Parent != nil {
				body.Parent.RemoveChild(body)
			}
			p.oe = p.oe[:1]
			p.addElement()
			p.im = inFramesetIM
			return true
		case a.Address, a.Article, a.Aside, a.Blockquote, a.Center, a.Details, a.Dir, a.Div, a.Dl, a.Fieldset, a.Figcaption, a.Figure, a.Footer, a.Header, a.Hgroup, a.Menu, a.Nav, a.Ol, a.P, a.Section, a.Summary, a.Ul:
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
			p.popUntil(buttonScope, a.P)
			switch n := p.top(); n.DataAtom {
			case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
				p.oe.pop()
			}
			p.addElement()
		case a.Pre, a.Listing:
			p.popUntil(buttonScope, a.P)
			p.addElement()
			
			p.framesetOK = false
		case a.Form:
			if p.form != nil && !p.oe.contains(a.Template) {
				
				return true
			}
			p.popUntil(buttonScope, a.P)
			p.addElement()
			if !p.oe.contains(a.Template) {
				p.form = p.top()
			}
		case a.Li:
			p.framesetOK = false
			for i := len(p.oe) - 1; i >= 0; i-- {
				node := p.oe[i]
				switch node.DataAtom {
				case a.Li:
					p.oe = p.oe[:i]
				case a.Address, a.Div, a.P:
					continue
				default:
					if !isSpecialElement(node) {
						continue
					}
				}
				break
			}
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Dd, a.Dt:
			p.framesetOK = false
			for i := len(p.oe) - 1; i >= 0; i-- {
				node := p.oe[i]
				switch node.DataAtom {
				case a.Dd, a.Dt:
					p.oe = p.oe[:i]
				case a.Address, a.Div, a.P:
					continue
				default:
					if !isSpecialElement(node) {
						continue
					}
				}
				break
			}
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Plaintext:
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Button:
			p.popUntil(defaultScope, a.Button)
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
		case a.A:
			for i := len(p.afe) - 1; i >= 0 && p.afe[i].Type != scopeMarkerNode; i-- {
				if n := p.afe[i]; n.Type == ElementNode && n.DataAtom == a.A {
					p.inBodyEndTagFormatting(a.A)
					p.oe.remove(n)
					p.afe.remove(n)
					break
				}
			}
			p.reconstructActiveFormattingElements()
			p.addFormattingElement()
		case a.B, a.Big, a.Code, a.Em, a.Font, a.I, a.S, a.Small, a.Strike, a.Strong, a.Tt, a.U:
			p.reconstructActiveFormattingElements()
			p.addFormattingElement()
		case a.Nobr:
			p.reconstructActiveFormattingElements()
			if p.elementInScope(defaultScope, a.Nobr) {
				p.inBodyEndTagFormatting(a.Nobr)
				p.reconstructActiveFormattingElements()
			}
			p.addFormattingElement()
		case a.Applet, a.Marquee, a.Object:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.framesetOK = false
		case a.Table:
			if !p.quirks {
				p.popUntil(buttonScope, a.P)
			}
			p.addElement()
			p.framesetOK = false
			p.im = inTableIM
			return true
		case a.Area, a.Br, a.Embed, a.Img, a.Input, a.Keygen, a.Wbr:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			if p.tok.DataAtom == a.Input {
				for _, t := range p.tok.Attr {
					if t.Key == "type" {
						if strings.ToLower(t.Val) == "hidden" {
							
							return true
						}
					}
				}
			}
			p.framesetOK = false
		case a.Param, a.Source, a.Track:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		case a.Hr:
			p.popUntil(buttonScope, a.P)
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			p.framesetOK = false
		case a.Image:
			p.tok.DataAtom = a.Img
			p.tok.Data = a.Img.String()
			return false
		case a.Isindex:
			if p.form != nil {
				
				return true
			}
			action := ""
			prompt := "This is a searchable index. Enter search keywords: "
			attr := []Attribute{{Key: "name", Val: "isindex"}}
			for _, t := range p.tok.Attr {
				switch t.Key {
				case "action":
					action = t.Val
				case "name":
					
				case "prompt":
					prompt = t.Val
				default:
					attr = append(attr, t)
				}
			}
			p.acknowledgeSelfClosingTag()
			p.popUntil(buttonScope, a.P)
			p.parseImpliedToken(StartTagToken, a.Form, a.Form.String())
			if p.form == nil {
				
				
				
				
				
				return true
			}
			if action != "" {
				p.form.Attr = []Attribute{{Key: "action", Val: action}}
			}
			p.parseImpliedToken(StartTagToken, a.Hr, a.Hr.String())
			p.parseImpliedToken(StartTagToken, a.Label, a.Label.String())
			p.addText(prompt)
			p.addChild(&Node{
				Type:     ElementNode,
				DataAtom: a.Input,
				Data:     a.Input.String(),
				Attr:     attr,
			})
			p.oe.pop()
			p.parseImpliedToken(EndTagToken, a.Label, a.Label.String())
			p.parseImpliedToken(StartTagToken, a.Hr, a.Hr.String())
			p.parseImpliedToken(EndTagToken, a.Form, a.Form.String())
		case a.Textarea:
			p.addElement()
			p.setOriginalIM()
			p.framesetOK = false
			p.im = textIM
		case a.Xmp:
			p.popUntil(buttonScope, a.P)
			p.reconstructActiveFormattingElements()
			p.framesetOK = false
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
		case a.Iframe:
			p.framesetOK = false
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
		case a.Noembed, a.Noscript:
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
		case a.Select:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
			p.im = inSelectIM
			return true
		case a.Optgroup, a.Option:
			if p.top().DataAtom == a.Option {
				p.oe.pop()
			}
			p.reconstructActiveFormattingElements()
			p.addElement()
		case a.Rb, a.Rtc:
			if p.elementInScope(defaultScope, a.Ruby) {
				p.generateImpliedEndTags()
			}
			p.addElement()
		case a.Rp, a.Rt:
			if p.elementInScope(defaultScope, a.Ruby) {
				p.generateImpliedEndTags("rtc")
			}
			p.addElement()
		case a.Math, a.Svg:
			p.reconstructActiveFormattingElements()
			if p.tok.DataAtom == a.Math {
				adjustAttributeNames(p.tok.Attr, mathMLAttributeAdjustments)
			} else {
				adjustAttributeNames(p.tok.Attr, svgAttributeAdjustments)
			}
			adjustForeignAttributes(p.tok.Attr)
			p.addElement()
			p.top().Namespace = p.tok.Data
			if p.hasSelfClosingToken {
				p.oe.pop()
				p.acknowledgeSelfClosingTag()
			}
			return true
		case a.Caption, a.Col, a.Colgroup, a.Frame, a.Head, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			
		default:
			p.reconstructActiveFormattingElements()
			p.addElement()
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Body:
			if p.elementInScope(defaultScope, a.Body) {
				p.im = afterBodyIM
			}
		case a.Html:
			if p.elementInScope(defaultScope, a.Body) {
				p.parseImpliedToken(EndTagToken, a.Body, a.Body.String())
				return false
			}
			return true
		case a.Address, a.Article, a.Aside, a.Blockquote, a.Button, a.Center, a.Details, a.Dir, a.Div, a.Dl, a.Fieldset, a.Figcaption, a.Figure, a.Footer, a.Header, a.Hgroup, a.Listing, a.Menu, a.Nav, a.Ol, a.Pre, a.Section, a.Summary, a.Ul:
			p.popUntil(defaultScope, p.tok.DataAtom)
		case a.Form:
			if p.oe.contains(a.Template) {
				i := p.indexOfElementInScope(defaultScope, a.Form)
				if i == -1 {
					
					return true
				}
				p.generateImpliedEndTags()
				if p.oe[i].DataAtom != a.Form {
					
					return true
				}
				p.popUntil(defaultScope, a.Form)
			} else {
				node := p.form
				p.form = nil
				i := p.indexOfElementInScope(defaultScope, a.Form)
				if node == nil || i == -1 || p.oe[i] != node {
					
					return true
				}
				p.generateImpliedEndTags()
				p.oe.remove(node)
			}
		case a.P:
			if !p.elementInScope(buttonScope, a.P) {
				p.parseImpliedToken(StartTagToken, a.P, a.P.String())
			}
			p.popUntil(buttonScope, a.P)
		case a.Li:
			p.popUntil(listItemScope, a.Li)
		case a.Dd, a.Dt:
			p.popUntil(defaultScope, p.tok.DataAtom)
		case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
			p.popUntil(defaultScope, a.H1, a.H2, a.H3, a.H4, a.H5, a.H6)
		case a.A, a.B, a.Big, a.Code, a.Em, a.Font, a.I, a.Nobr, a.S, a.Small, a.Strike, a.Strong, a.Tt, a.U:
			p.inBodyEndTagFormatting(p.tok.DataAtom)
		case a.Applet, a.Marquee, a.Object:
			if p.popUntil(defaultScope, p.tok.DataAtom) {
				p.clearActiveFormattingElements()
			}
		case a.Br:
			p.tok.Type = StartTagToken
			return false
		case a.Template:
			return inHeadIM(p)
		default:
			p.inBodyEndTagOther(p.tok.DataAtom)
		}
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
	case ErrorToken:
		
		if len(p.templateStack) > 0 {
			p.im = inTemplateIM
			return false
		} else {
			for _, e := range p.oe {
				switch e.DataAtom {
				case a.Dd, a.Dt, a.Li, a.Optgroup, a.Option, a.P, a.Rb, a.Rp, a.Rt, a.Rtc, a.Tbody, a.Td, a.Tfoot, a.Th,
					a.Thead, a.Tr, a.Body, a.Html:
				default:
					return true
				}
			}
		}
	}

	return true
}

func (p *parser) inBodyEndTagFormatting(tagAtom a.Atom) {
	
	

	
	
	

	
	for i := 0; i < 8; i++ {
		
		var formattingElement *Node
		for j := len(p.afe) - 1; j >= 0; j-- {
			if p.afe[j].Type == scopeMarkerNode {
				break
			}
			if p.afe[j].DataAtom == tagAtom {
				formattingElement = p.afe[j]
				break
			}
		}
		if formattingElement == nil {
			p.inBodyEndTagOther(tagAtom)
			return
		}
		feIndex := p.oe.index(formattingElement)
		if feIndex == -1 {
			p.afe.remove(formattingElement)
			return
		}
		if !p.elementInScope(defaultScope, tagAtom) {
			
			return
		}

		
		var furthestBlock *Node
		for _, e := range p.oe[feIndex:] {
			if isSpecialElement(e) {
				furthestBlock = e
				break
			}
		}
		if furthestBlock == nil {
			e := p.oe.pop()
			for e != formattingElement {
				e = p.oe.pop()
			}
			p.afe.remove(e)
			return
		}

		
		commonAncestor := p.oe[feIndex-1]
		bookmark := p.afe.index(formattingElement)

		
		lastNode := furthestBlock
		node := furthestBlock
		x := p.oe.index(node)
		
		for j := 0; j < 3; j++ {
			
			x--
			node = p.oe[x]
			
			if p.afe.index(node) == -1 {
				p.oe.remove(node)
				continue
			}
			
			if node == formattingElement {
				break
			}
			
			clone := node.clone()
			p.afe[p.afe.index(node)] = clone
			p.oe[p.oe.index(node)] = clone
			node = clone
			
			if lastNode == furthestBlock {
				bookmark = p.afe.index(node) + 1
			}
			
			if lastNode.Parent != nil {
				lastNode.Parent.RemoveChild(lastNode)
			}
			node.AppendChild(lastNode)
			
			lastNode = node
		}

		
		
		if lastNode.Parent != nil {
			lastNode.Parent.RemoveChild(lastNode)
		}
		switch commonAncestor.DataAtom {
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			p.fosterParent(lastNode)
		default:
			commonAncestor.AppendChild(lastNode)
		}

		
		
		clone := formattingElement.clone()
		reparentChildren(clone, furthestBlock)
		furthestBlock.AppendChild(clone)

		
		if oldLoc := p.afe.index(formattingElement); oldLoc != -1 && oldLoc < bookmark {
			
			bookmark--
		}
		p.afe.remove(formattingElement)
		p.afe.insert(bookmark, clone)

		
		p.oe.remove(formattingElement)
		p.oe.insert(p.oe.index(furthestBlock)+1, clone)
	}
}




func (p *parser) inBodyEndTagOther(tagAtom a.Atom) {
	for i := len(p.oe) - 1; i >= 0; i-- {
		if p.oe[i].DataAtom == tagAtom {
			p.oe = p.oe[:i]
			break
		}
		if isSpecialElement(p.oe[i]) {
			break
		}
	}
}


func textIM(p *parser) bool {
	switch p.tok.Type {
	case ErrorToken:
		p.oe.pop()
	case TextToken:
		d := p.tok.Data
		if n := p.oe.top(); n.DataAtom == a.Textarea && n.FirstChild == nil {
			
			if d != "" && d[0] == '\r' {
				d = d[1:]
			}
			if d != "" && d[0] == '\n' {
				d = d[1:]
			}
		}
		if d == "" {
			return true
		}
		p.addText(d)
		return true
	case EndTagToken:
		p.oe.pop()
	}
	p.im = p.originalIM
	p.originalIM = nil
	return p.tok.Type == EndTagToken
}


func inTableIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		p.tok.Data = strings.Replace(p.tok.Data, "\x00", "", -1)
		switch p.oe.top().DataAtom {
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if strings.Trim(p.tok.Data, whitespace) == "" {
				p.addText(p.tok.Data)
				return true
			}
		}
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Caption:
			p.clearStackToContext(tableScope)
			p.afe = append(p.afe, &scopeMarker)
			p.addElement()
			p.im = inCaptionIM
			return true
		case a.Colgroup:
			p.clearStackToContext(tableScope)
			p.addElement()
			p.im = inColumnGroupIM
			return true
		case a.Col:
			p.parseImpliedToken(StartTagToken, a.Colgroup, a.Colgroup.String())
			return false
		case a.Tbody, a.Tfoot, a.Thead:
			p.clearStackToContext(tableScope)
			p.addElement()
			p.im = inTableBodyIM
			return true
		case a.Td, a.Th, a.Tr:
			p.parseImpliedToken(StartTagToken, a.Tbody, a.Tbody.String())
			return false
		case a.Table:
			if p.popUntil(tableScope, a.Table) {
				p.resetInsertionMode()
				return false
			}
			
			return true
		case a.Style, a.Script, a.Template:
			return inHeadIM(p)
		case a.Input:
			for _, t := range p.tok.Attr {
				if t.Key == "type" && strings.ToLower(t.Val) == "hidden" {
					p.addElement()
					p.oe.pop()
					return true
				}
			}
			
		case a.Form:
			if p.oe.contains(a.Template) || p.form != nil {
				
				return true
			}
			p.addElement()
			p.form = p.oe.pop()
		case a.Select:
			p.reconstructActiveFormattingElements()
			switch p.top().DataAtom {
			case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
				p.fosterParenting = true
			}
			p.addElement()
			p.fosterParenting = false
			p.framesetOK = false
			p.im = inSelectInTableIM
			return true
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Table:
			if p.popUntil(tableScope, a.Table) {
				p.resetInsertionMode()
				return true
			}
			
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			
			return true
		case a.Template:
			return inHeadIM(p)
		}
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	case DoctypeToken:
		
		return true
	case ErrorToken:
		return inBodyIM(p)
	}

	p.fosterParenting = true
	defer func() { p.fosterParenting = false }()

	return inBodyIM(p)
}


func inCaptionIM(p *parser) bool {
	switch p.tok.Type {
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Td, a.Tfoot, a.Thead, a.Tr:
			if p.popUntil(tableScope, a.Caption) {
				p.clearActiveFormattingElements()
				p.im = inTableIM
				return false
			} else {
				
				return true
			}
		case a.Select:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
			p.im = inSelectInTableIM
			return true
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Caption:
			if p.popUntil(tableScope, a.Caption) {
				p.clearActiveFormattingElements()
				p.im = inTableIM
			}
			return true
		case a.Table:
			if p.popUntil(tableScope, a.Caption) {
				p.clearActiveFormattingElements()
				p.im = inTableIM
				return false
			} else {
				
				return true
			}
		case a.Body, a.Col, a.Colgroup, a.Html, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			
			return true
		}
	}
	return inBodyIM(p)
}


func inColumnGroupIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) < len(p.tok.Data) {
			
			p.addText(p.tok.Data[:len(p.tok.Data)-len(s)])
			if s == "" {
				return true
			}
			p.tok.Data = s
		}
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	case DoctypeToken:
		
		return true
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Col:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			return true
		case a.Template:
			return inHeadIM(p)
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Colgroup:
			if p.oe.top().DataAtom == a.Colgroup {
				p.oe.pop()
				p.im = inTableIM
			}
			return true
		case a.Col:
			
			return true
		case a.Template:
			return inHeadIM(p)
		}
	case ErrorToken:
		return inBodyIM(p)
	}
	if p.oe.top().DataAtom != a.Colgroup {
		return true
	}
	p.oe.pop()
	p.im = inTableIM
	return false
}


func inTableBodyIM(p *parser) bool {
	switch p.tok.Type {
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Tr:
			p.clearStackToContext(tableBodyScope)
			p.addElement()
			p.im = inRowIM
			return true
		case a.Td, a.Th:
			p.parseImpliedToken(StartTagToken, a.Tr, a.Tr.String())
			return false
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Tfoot, a.Thead:
			if p.popUntil(tableScope, a.Tbody, a.Thead, a.Tfoot) {
				p.im = inTableIM
				return false
			}
			
			return true
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Tbody, a.Tfoot, a.Thead:
			if p.elementInScope(tableScope, p.tok.DataAtom) {
				p.clearStackToContext(tableBodyScope)
				p.oe.pop()
				p.im = inTableIM
			}
			return true
		case a.Table:
			if p.popUntil(tableScope, a.Tbody, a.Thead, a.Tfoot) {
				p.im = inTableIM
				return false
			}
			
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Td, a.Th, a.Tr:
			
			return true
		}
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	}

	return inTableIM(p)
}


func inRowIM(p *parser) bool {
	switch p.tok.Type {
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Td, a.Th:
			p.clearStackToContext(tableRowScope)
			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.im = inCellIM
			return true
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if p.popUntil(tableScope, a.Tr) {
				p.im = inTableBodyIM
				return false
			}
			
			return true
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Tr:
			if p.popUntil(tableScope, a.Tr) {
				p.im = inTableBodyIM
				return true
			}
			
			return true
		case a.Table:
			if p.popUntil(tableScope, a.Tr) {
				p.im = inTableBodyIM
				return false
			}
			
			return true
		case a.Tbody, a.Tfoot, a.Thead:
			if p.elementInScope(tableScope, p.tok.DataAtom) {
				p.parseImpliedToken(EndTagToken, a.Tr, a.Tr.String())
				return false
			}
			
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Td, a.Th:
			
			return true
		}
	}

	return inTableIM(p)
}


func inCellIM(p *parser) bool {
	switch p.tok.Type {
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			if p.popUntil(tableScope, a.Td, a.Th) {
				
				p.clearActiveFormattingElements()
				p.im = inRowIM
				return false
			}
			
			return true
		case a.Select:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
			p.im = inSelectInTableIM
			return true
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Td, a.Th:
			if !p.popUntil(tableScope, p.tok.DataAtom) {
				
				return true
			}
			p.clearActiveFormattingElements()
			p.im = inRowIM
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html:
			
			return true
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if !p.elementInScope(tableScope, p.tok.DataAtom) {
				
				return true
			}
			
			p.popUntil(tableScope, a.Td, a.Th)
			p.clearActiveFormattingElements()
			p.im = inRowIM
			return false
		}
	}
	return inBodyIM(p)
}


func inSelectIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		p.addText(strings.Replace(p.tok.Data, "\x00", "", -1))
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Option:
			if p.top().DataAtom == a.Option {
				p.oe.pop()
			}
			p.addElement()
		case a.Optgroup:
			if p.top().DataAtom == a.Option {
				p.oe.pop()
			}
			if p.top().DataAtom == a.Optgroup {
				p.oe.pop()
			}
			p.addElement()
		case a.Select:
			p.tok.Type = EndTagToken
			return false
		case a.Input, a.Keygen, a.Textarea:
			if p.elementInScope(selectScope, a.Select) {
				p.parseImpliedToken(EndTagToken, a.Select, a.Select.String())
				return false
			}
			
			p.tokenizer.NextIsNotRawText()
			
			return true
		case a.Script, a.Template:
			return inHeadIM(p)
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Option:
			if p.top().DataAtom == a.Option {
				p.oe.pop()
			}
		case a.Optgroup:
			i := len(p.oe) - 1
			if p.oe[i].DataAtom == a.Option {
				i--
			}
			if p.oe[i].DataAtom == a.Optgroup {
				p.oe = p.oe[:i]
			}
		case a.Select:
			if p.popUntil(selectScope, a.Select) {
				p.resetInsertionMode()
			}
		case a.Template:
			return inHeadIM(p)
		}
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
	case DoctypeToken:
		
		return true
	case ErrorToken:
		return inBodyIM(p)
	}

	return true
}


func inSelectInTableIM(p *parser) bool {
	switch p.tok.Type {
	case StartTagToken, EndTagToken:
		switch p.tok.DataAtom {
		case a.Caption, a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr, a.Td, a.Th:
			if p.tok.Type == StartTagToken || p.elementInScope(tableScope, p.tok.DataAtom) {
				p.parseImpliedToken(EndTagToken, a.Select, a.Select.String())
				return false
			} else {
				
				return true
			}
		}
	}
	return inSelectIM(p)
}


func inTemplateIM(p *parser) bool {
	switch p.tok.Type {
	case TextToken, CommentToken, DoctypeToken:
		return inBodyIM(p)
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Base, a.Basefont, a.Bgsound, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			return inHeadIM(p)
		case a.Caption, a.Colgroup, a.Tbody, a.Tfoot, a.Thead:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inTableIM)
			p.im = inTableIM
			return false
		case a.Col:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inColumnGroupIM)
			p.im = inColumnGroupIM
			return false
		case a.Tr:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inTableBodyIM)
			p.im = inTableBodyIM
			return false
		case a.Td, a.Th:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inRowIM)
			p.im = inRowIM
			return false
		default:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inBodyIM)
			p.im = inBodyIM
			return false
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Template:
			return inHeadIM(p)
		default:
			
			return true
		}
	case ErrorToken:
		if !p.oe.contains(a.Template) {
			
			return true
		}
		
		
		
		p.generateImpliedEndTags()
		for i := len(p.oe) - 1; i >= 0; i-- {
			if n := p.oe[i]; n.Namespace == "" && n.DataAtom == a.Template {
				p.oe = p.oe[:i]
				break
			}
		}
		p.clearActiveFormattingElements()
		p.templateStack.pop()
		p.resetInsertionMode()
		return false
	}
	return false
}


func afterBodyIM(p *parser) bool {
	switch p.tok.Type {
	case ErrorToken:
		
		return true
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) == 0 {
			
			return inBodyIM(p)
		}
	case StartTagToken:
		if p.tok.DataAtom == a.Html {
			return inBodyIM(p)
		}
	case EndTagToken:
		if p.tok.DataAtom == a.Html {
			if !p.fragment {
				p.im = afterAfterBodyIM
			}
			return true
		}
	case CommentToken:
		
		if len(p.oe) < 1 || p.oe[0].DataAtom != a.Html {
			panic("html: bad parser state: <html> element not found, in the after-body insertion mode")
		}
		p.oe[0].AppendChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	}
	p.im = inBodyIM
	return false
}


func inFramesetIM(p *parser) bool {
	switch p.tok.Type {
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
	case TextToken:
		
		s := strings.Map(func(c rune) rune {
			switch c {
			case ' ', '\t', '\n', '\f', '\r':
				return c
			}
			return -1
		}, p.tok.Data)
		if s != "" {
			p.addText(s)
		}
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Frameset:
			p.addElement()
		case a.Frame:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		case a.Noframes:
			return inHeadIM(p)
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Frameset:
			if p.oe.top().DataAtom != a.Html {
				p.oe.pop()
				if p.oe.top().DataAtom != a.Frameset {
					p.im = afterFramesetIM
					return true
				}
			}
		}
	default:
		
	}
	return true
}


func afterFramesetIM(p *parser) bool {
	switch p.tok.Type {
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
	case TextToken:
		
		s := strings.Map(func(c rune) rune {
			switch c {
			case ' ', '\t', '\n', '\f', '\r':
				return c
			}
			return -1
		}, p.tok.Data)
		if s != "" {
			p.addText(s)
		}
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Noframes:
			return inHeadIM(p)
		}
	case EndTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			p.im = afterAfterFramesetIM
			return true
		}
	default:
		
	}
	return true
}


func afterAfterBodyIM(p *parser) bool {
	switch p.tok.Type {
	case ErrorToken:
		
		return true
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) == 0 {
			
			return inBodyIM(p)
		}
	case StartTagToken:
		if p.tok.DataAtom == a.Html {
			return inBodyIM(p)
		}
	case CommentToken:
		p.doc.AppendChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
		return true
	case DoctypeToken:
		return inBodyIM(p)
	}
	p.im = inBodyIM
	return false
}


func afterAfterFramesetIM(p *parser) bool {
	switch p.tok.Type {
	case CommentToken:
		p.doc.AppendChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
	case TextToken:
		
		s := strings.Map(func(c rune) rune {
			switch c {
			case ' ', '\t', '\n', '\f', '\r':
				return c
			}
			return -1
		}, p.tok.Data)
		if s != "" {
			p.tok.Data = s
			return inBodyIM(p)
		}
	case StartTagToken:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Noframes:
			return inHeadIM(p)
		}
	case DoctypeToken:
		return inBodyIM(p)
	default:
		
	}
	return true
}

const whitespaceOrNUL = whitespace + "\x00"


func parseForeignContent(p *parser) bool {
	switch p.tok.Type {
	case TextToken:
		if p.framesetOK {
			p.framesetOK = strings.TrimLeft(p.tok.Data, whitespaceOrNUL) == ""
		}
		p.tok.Data = strings.Replace(p.tok.Data, "\x00", "\ufffd", -1)
		p.addText(p.tok.Data)
	case CommentToken:
		p.addChild(&Node{
			Type: CommentNode,
			Data: p.tok.Data,
		})
	case StartTagToken:
		b := breakout[p.tok.Data]
		if p.tok.DataAtom == a.Font {
		loop:
			for _, attr := range p.tok.Attr {
				switch attr.Key {
				case "color", "face", "size":
					b = true
					break loop
				}
			}
		}
		if b {
			for i := len(p.oe) - 1; i >= 0; i-- {
				n := p.oe[i]
				if n.Namespace == "" || htmlIntegrationPoint(n) || mathMLTextIntegrationPoint(n) {
					p.oe = p.oe[:i+1]
					break
				}
			}
			return false
		}
		switch p.top().Namespace {
		case "math":
			adjustAttributeNames(p.tok.Attr, mathMLAttributeAdjustments)
		case "svg":
			
			
			if x := svgTagNameAdjustments[p.tok.Data]; x != "" {
				p.tok.DataAtom = a.Lookup([]byte(x))
				p.tok.Data = x
			}
			adjustAttributeNames(p.tok.Attr, svgAttributeAdjustments)
		default:
			panic("html: bad parser state: unexpected namespace")
		}
		adjustForeignAttributes(p.tok.Attr)
		namespace := p.top().Namespace
		p.addElement()
		p.top().Namespace = namespace
		if namespace != "" {
			
			
			p.tokenizer.NextIsNotRawText()
		}
		if p.hasSelfClosingToken {
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		}
	case EndTagToken:
		for i := len(p.oe) - 1; i >= 0; i-- {
			if p.oe[i].Namespace == "" {
				return p.im(p)
			}
			if strings.EqualFold(p.oe[i].Data, p.tok.Data) {
				p.oe = p.oe[:i]
				break
			}
		}
		return true
	default:
		
	}
	return true
}


func (p *parser) inForeignContent() bool {
	if len(p.oe) == 0 {
		return false
	}
	n := p.oe[len(p.oe)-1]
	if n.Namespace == "" {
		return false
	}
	if mathMLTextIntegrationPoint(n) {
		if p.tok.Type == StartTagToken && p.tok.DataAtom != a.Mglyph && p.tok.DataAtom != a.Malignmark {
			return false
		}
		if p.tok.Type == TextToken {
			return false
		}
	}
	if n.Namespace == "math" && n.DataAtom == a.AnnotationXml && p.tok.Type == StartTagToken && p.tok.DataAtom == a.Svg {
		return false
	}
	if htmlIntegrationPoint(n) && (p.tok.Type == StartTagToken || p.tok.Type == TextToken) {
		return false
	}
	if p.tok.Type == ErrorToken {
		return false
	}
	return true
}



func (p *parser) parseImpliedToken(t TokenType, dataAtom a.Atom, data string) {
	realToken, selfClosing := p.tok, p.hasSelfClosingToken
	p.tok = Token{
		Type:     t,
		DataAtom: dataAtom,
		Data:     data,
	}
	p.hasSelfClosingToken = false
	p.parseCurrentToken()
	p.tok, p.hasSelfClosingToken = realToken, selfClosing
}



func (p *parser) parseCurrentToken() {
	if p.tok.Type == SelfClosingTagToken {
		p.hasSelfClosingToken = true
		p.tok.Type = StartTagToken
	}

	consumed := false
	for !consumed {
		if p.inForeignContent() {
			consumed = parseForeignContent(p)
		} else {
			consumed = p.im(p)
		}
	}

	if p.hasSelfClosingToken {
		
		p.hasSelfClosingToken = false
	}
}

func (p *parser) parse() error {
	
	var err error
	for err != io.EOF {
		
		n := p.oe.top()
		p.tokenizer.AllowCDATA(n != nil && n.Namespace != "")
		
		p.tokenizer.Next()
		p.tok = p.tokenizer.Token()
		if p.tok.Type == ErrorToken {
			err = p.tokenizer.Err()
			if err != nil && err != io.EOF {
				return err
			}
		}
		p.parseCurrentToken()
	}
	return nil
}












func Parse(r io.Reader) (*Node, error) {
	p := &parser{
		tokenizer: NewTokenizer(r),
		doc: &Node{
			Type: DocumentNode,
		},
		scripting:  true,
		framesetOK: true,
		im:         initialIM,
	}
	err := p.parse()
	if err != nil {
		return nil, err
	}
	return p.doc, nil
}






func ParseFragment(r io.Reader, context *Node) ([]*Node, error) {
	contextTag := ""
	if context != nil {
		if context.Type != ElementNode {
			return nil, errors.New("html: ParseFragment of non-element Node")
		}
		
		
		
		if context.DataAtom != a.Lookup([]byte(context.Data)) {
			return nil, fmt.Errorf("html: inconsistent Node: DataAtom=%q, Data=%q", context.DataAtom, context.Data)
		}
		contextTag = context.DataAtom.String()
	}
	p := &parser{
		tokenizer: NewTokenizerFragment(r, contextTag),
		doc: &Node{
			Type: DocumentNode,
		},
		scripting: true,
		fragment:  true,
		context:   context,
	}

	root := &Node{
		Type:     ElementNode,
		DataAtom: a.Html,
		Data:     a.Html.String(),
	}
	p.doc.AppendChild(root)
	p.oe = nodeStack{root}
	if context != nil && context.DataAtom == a.Template {
		p.templateStack = append(p.templateStack, inTemplateIM)
	}
	p.resetInsertionMode()

	for n := context; n != nil; n = n.Parent {
		if n.Type == ElementNode && n.DataAtom == a.Form {
			p.form = n
			break
		}
	}

	err := p.parse()
	if err != nil {
		return nil, err
	}

	parent := p.doc
	if context != nil {
		parent = root
	}

	var result []*Node
	for c := parent.FirstChild; c != nil; {
		next := c.NextSibling
		parent.RemoveChild(c)
		result = append(result, c)
		c = next
	}
	return result, nil
}
