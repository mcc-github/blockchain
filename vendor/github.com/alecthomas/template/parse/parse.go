







package parse

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)


type Tree struct {
	Name      string    
	ParseName string    
	Root      *ListNode 
	text      string    
	
	funcs     []map[string]interface{}
	lex       *lexer
	token     [3]item 
	peekCount int
	vars      []string 
}


func (t *Tree) Copy() *Tree {
	if t == nil {
		return nil
	}
	return &Tree{
		Name:      t.Name,
		ParseName: t.ParseName,
		Root:      t.Root.CopyList(),
		text:      t.text,
	}
}





func Parse(name, text, leftDelim, rightDelim string, funcs ...map[string]interface{}) (treeSet map[string]*Tree, err error) {
	treeSet = make(map[string]*Tree)
	t := New(name)
	t.text = text
	_, err = t.Parse(text, leftDelim, rightDelim, treeSet, funcs...)
	return
}


func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}


func (t *Tree) backup() {
	t.peekCount++
}



func (t *Tree) backup2(t1 item) {
	t.token[1] = t1
	t.peekCount = 2
}



func (t *Tree) backup3(t2, t1 item) { 
	t.token[1] = t1
	t.token[2] = t2
	t.peekCount = 3
}


func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}


func (t *Tree) nextNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	return token
}


func (t *Tree) peekNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	t.backup()
	return token
}




func New(name string, funcs ...map[string]interface{}) *Tree {
	return &Tree{
		Name:  name,
		funcs: funcs,
	}
}




func (t *Tree) ErrorContext(n Node) (location, context string) {
	pos := int(n.Position())
	tree := n.tree()
	if tree == nil {
		tree = t
	}
	text := tree.text[:pos]
	byteNum := strings.LastIndex(text, "\n")
	if byteNum == -1 {
		byteNum = pos 
	} else {
		byteNum++ 
		byteNum = pos - byteNum
	}
	lineNum := 1 + strings.Count(text, "\n")
	context = n.String()
	if len(context) > 20 {
		context = fmt.Sprintf("%.20s...", context)
	}
	return fmt.Sprintf("%s:%d:%d", tree.ParseName, lineNum, byteNum), context
}


func (t *Tree) errorf(format string, args ...interface{}) {
	t.Root = nil
	format = fmt.Sprintf("template: %s:%d: %s", t.ParseName, t.lex.lineNumber(), format)
	panic(fmt.Errorf(format, args...))
}


func (t *Tree) error(err error) {
	t.errorf("%s", err)
}


func (t *Tree) expect(expected itemType, context string) item {
	token := t.nextNonSpace()
	if token.typ != expected {
		t.unexpected(token, context)
	}
	return token
}


func (t *Tree) expectOneOf(expected1, expected2 itemType, context string) item {
	token := t.nextNonSpace()
	if token.typ != expected1 && token.typ != expected2 {
		t.unexpected(token, context)
	}
	return token
}


func (t *Tree) unexpected(token item, context string) {
	t.errorf("unexpected %s in %s", token, context)
}


func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.stopParse()
		}
		*errp = e.(error)
	}
	return
}


func (t *Tree) startParse(funcs []map[string]interface{}, lex *lexer) {
	t.Root = nil
	t.lex = lex
	t.vars = []string{"$"}
	t.funcs = funcs
}


func (t *Tree) stopParse() {
	t.lex = nil
	t.vars = nil
	t.funcs = nil
}





func (t *Tree) Parse(text, leftDelim, rightDelim string, treeSet map[string]*Tree, funcs ...map[string]interface{}) (tree *Tree, err error) {
	defer t.recover(&err)
	t.ParseName = t.Name
	t.startParse(funcs, lex(t.Name, text, leftDelim, rightDelim))
	t.text = text
	t.parse(treeSet)
	t.add(treeSet)
	t.stopParse()
	return t, nil
}


func (t *Tree) add(treeSet map[string]*Tree) {
	tree := treeSet[t.Name]
	if tree == nil || IsEmptyTree(tree.Root) {
		treeSet[t.Name] = t
		return
	}
	if !IsEmptyTree(t.Root) {
		t.errorf("template: multiple definition of template %q", t.Name)
	}
}


func IsEmptyTree(n Node) bool {
	switch n := n.(type) {
	case nil:
		return true
	case *ActionNode:
	case *IfNode:
	case *ListNode:
		for _, node := range n.Nodes {
			if !IsEmptyTree(node) {
				return false
			}
		}
		return true
	case *RangeNode:
	case *TemplateNode:
	case *TextNode:
		return len(bytes.TrimSpace(n.Text)) == 0
	case *WithNode:
	default:
		panic("unknown node: " + n.String())
	}
	return false
}




func (t *Tree) parse(treeSet map[string]*Tree) (next Node) {
	t.Root = t.newList(t.peek().pos)
	for t.peek().typ != itemEOF {
		if t.peek().typ == itemLeftDelim {
			delim := t.next()
			if t.nextNonSpace().typ == itemDefine {
				newT := New("definition") 
				newT.text = t.text
				newT.ParseName = t.ParseName
				newT.startParse(t.funcs, t.lex)
				newT.parseDefinition(treeSet)
				continue
			}
			t.backup2(delim)
		}
		n := t.textOrAction()
		if n.Type() == nodeEnd {
			t.errorf("unexpected %s", n)
		}
		t.Root.append(n)
	}
	return nil
}




func (t *Tree) parseDefinition(treeSet map[string]*Tree) {
	const context = "define clause"
	name := t.expectOneOf(itemString, itemRawString, context)
	var err error
	t.Name, err = strconv.Unquote(name.val)
	if err != nil {
		t.error(err)
	}
	t.expect(itemRightDelim, context)
	var end Node
	t.Root, end = t.itemList()
	if end.Type() != nodeEnd {
		t.errorf("unexpected %s in %s", end, context)
	}
	t.add(treeSet)
	t.stopParse()
}




func (t *Tree) itemList() (list *ListNode, next Node) {
	list = t.newList(t.peekNonSpace().pos)
	for t.peekNonSpace().typ != itemEOF {
		n := t.textOrAction()
		switch n.Type() {
		case nodeEnd, nodeElse:
			return list, n
		}
		list.append(n)
	}
	t.errorf("unexpected EOF")
	return
}



func (t *Tree) textOrAction() Node {
	switch token := t.nextNonSpace(); token.typ {
	case itemElideNewline:
		return t.elideNewline()
	case itemText:
		return t.newText(token.pos, token.val)
	case itemLeftDelim:
		return t.action()
	default:
		t.unexpected(token, "input")
	}
	return nil
}



func (t *Tree) elideNewline() Node {
	token := t.peek()
	if token.typ != itemText {
		t.unexpected(token, "input")
		return nil
	}

	t.next()
	stripped := strings.TrimLeft(token.val, "\n\r")
	diff := len(token.val) - len(stripped)
	if diff > 0 {
		
		
		token.pos += Pos(diff)
		token.val = stripped
	}
	return t.newText(token.pos, token.val)
}






func (t *Tree) action() (n Node) {
	switch token := t.nextNonSpace(); token.typ {
	case itemElse:
		return t.elseControl()
	case itemEnd:
		return t.endControl()
	case itemIf:
		return t.ifControl()
	case itemRange:
		return t.rangeControl()
	case itemTemplate:
		return t.templateControl()
	case itemWith:
		return t.withControl()
	}
	t.backup()
	
	return t.newAction(t.peek().pos, t.lex.lineNumber(), t.pipeline("command"))
}



func (t *Tree) pipeline(context string) (pipe *PipeNode) {
	var decl []*VariableNode
	pos := t.peekNonSpace().pos
	
	for {
		if v := t.peekNonSpace(); v.typ == itemVariable {
			t.next()
			
			
			
			
			tokenAfterVariable := t.peek()
			if next := t.peekNonSpace(); next.typ == itemColonEquals || (next.typ == itemChar && next.val == ",") {
				t.nextNonSpace()
				variable := t.newVariable(v.pos, v.val)
				decl = append(decl, variable)
				t.vars = append(t.vars, v.val)
				if next.typ == itemChar && next.val == "," {
					if context == "range" && len(decl) < 2 {
						continue
					}
					t.errorf("too many declarations in %s", context)
				}
			} else if tokenAfterVariable.typ == itemSpace {
				t.backup3(v, tokenAfterVariable)
			} else {
				t.backup2(v)
			}
		}
		break
	}
	pipe = t.newPipeline(pos, t.lex.lineNumber(), decl)
	for {
		switch token := t.nextNonSpace(); token.typ {
		case itemRightDelim, itemRightParen:
			if len(pipe.Cmds) == 0 {
				t.errorf("missing value for %s", context)
			}
			if token.typ == itemRightParen {
				t.backup()
			}
			return
		case itemBool, itemCharConstant, itemComplex, itemDot, itemField, itemIdentifier,
			itemNumber, itemNil, itemRawString, itemString, itemVariable, itemLeftParen:
			t.backup()
			pipe.append(t.command())
		default:
			t.unexpected(token, context)
		}
	}
}

func (t *Tree) parseControl(allowElseIf bool, context string) (pos Pos, line int, pipe *PipeNode, list, elseList *ListNode) {
	defer t.popVars(len(t.vars))
	line = t.lex.lineNumber()
	pipe = t.pipeline(context)
	var next Node
	list, next = t.itemList()
	switch next.Type() {
	case nodeEnd: 
	case nodeElse:
		if allowElseIf {
			
			
			
			
			
			
			
			
			if t.peek().typ == itemIf {
				t.next() 
				elseList = t.newList(next.Position())
				elseList.append(t.ifControl())
				
				break
			}
		}
		elseList, next = t.itemList()
		if next.Type() != nodeEnd {
			t.errorf("expected end; found %s", next)
		}
	}
	return pipe.Position(), line, pipe, list, elseList
}





func (t *Tree) ifControl() Node {
	return t.newIf(t.parseControl(true, "if"))
}





func (t *Tree) rangeControl() Node {
	return t.newRange(t.parseControl(false, "range"))
}





func (t *Tree) withControl() Node {
	return t.newWith(t.parseControl(false, "with"))
}




func (t *Tree) endControl() Node {
	return t.newEnd(t.expect(itemRightDelim, "end").pos)
}




func (t *Tree) elseControl() Node {
	
	peek := t.peekNonSpace()
	if peek.typ == itemIf {
		
		return t.newElse(peek.pos, t.lex.lineNumber())
	}
	return t.newElse(t.expect(itemRightDelim, "else").pos, t.lex.lineNumber())
}





func (t *Tree) templateControl() Node {
	var name string
	token := t.nextNonSpace()
	switch token.typ {
	case itemString, itemRawString:
		s, err := strconv.Unquote(token.val)
		if err != nil {
			t.error(err)
		}
		name = s
	default:
		t.unexpected(token, "template invocation")
	}
	var pipe *PipeNode
	if t.nextNonSpace().typ != itemRightDelim {
		t.backup()
		
		pipe = t.pipeline("template")
	}
	return t.newTemplate(token.pos, t.lex.lineNumber(), name, pipe)
}





func (t *Tree) command() *CommandNode {
	cmd := t.newCommand(t.peekNonSpace().pos)
	for {
		t.peekNonSpace() 
		operand := t.operand()
		if operand != nil {
			cmd.append(operand)
		}
		switch token := t.next(); token.typ {
		case itemSpace:
			continue
		case itemError:
			t.errorf("%s", token.val)
		case itemRightDelim, itemRightParen:
			t.backup()
		case itemPipe:
		default:
			t.errorf("unexpected %s in operand; missing space?", token)
		}
		break
	}
	if len(cmd.Args) == 0 {
		t.errorf("empty command")
	}
	return cmd
}






func (t *Tree) operand() Node {
	node := t.term()
	if node == nil {
		return nil
	}
	if t.peek().typ == itemField {
		chain := t.newChain(t.peek().pos, node)
		for t.peek().typ == itemField {
			chain.Add(t.next().val)
		}
		
		
		
		
		switch node.Type() {
		case NodeField:
			node = t.newField(chain.Position(), chain.String())
		case NodeVariable:
			node = t.newVariable(chain.Position(), chain.String())
		default:
			node = chain
		}
	}
	return node
}










func (t *Tree) term() Node {
	switch token := t.nextNonSpace(); token.typ {
	case itemError:
		t.errorf("%s", token.val)
	case itemIdentifier:
		if !t.hasFunction(token.val) {
			t.errorf("function %q not defined", token.val)
		}
		return NewIdentifier(token.val).SetTree(t).SetPos(token.pos)
	case itemDot:
		return t.newDot(token.pos)
	case itemNil:
		return t.newNil(token.pos)
	case itemVariable:
		return t.useVar(token.pos, token.val)
	case itemField:
		return t.newField(token.pos, token.val)
	case itemBool:
		return t.newBool(token.pos, token.val == "true")
	case itemCharConstant, itemComplex, itemNumber:
		number, err := t.newNumber(token.pos, token.val, token.typ)
		if err != nil {
			t.error(err)
		}
		return number
	case itemLeftParen:
		pipe := t.pipeline("parenthesized pipeline")
		if token := t.next(); token.typ != itemRightParen {
			t.errorf("unclosed right paren: unexpected %s", token)
		}
		return pipe
	case itemString, itemRawString:
		s, err := strconv.Unquote(token.val)
		if err != nil {
			t.error(err)
		}
		return t.newString(token.pos, token.val, s)
	}
	t.backup()
	return nil
}


func (t *Tree) hasFunction(name string) bool {
	for _, funcMap := range t.funcs {
		if funcMap == nil {
			continue
		}
		if funcMap[name] != nil {
			return true
		}
	}
	return false
}


func (t *Tree) popVars(n int) {
	t.vars = t.vars[:n]
}



func (t *Tree) useVar(pos Pos, name string) Node {
	v := t.newVariable(pos, name)
	for _, varName := range t.vars {
		if varName == v.Ident[0] {
			return v
		}
	}
	t.errorf("undefined variable %q", v.Ident[0])
	return nil
}
