package toml

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int

const (
	itemError itemType = iota
	itemNIL            
	itemEOF
	itemText
	itemString
	itemRawString
	itemMultilineString
	itemRawMultilineString
	itemBool
	itemInteger
	itemFloat
	itemDatetime
	itemArray 
	itemArrayEnd
	itemTableStart
	itemTableEnd
	itemArrayTableStart
	itemArrayTableEnd
	itemKeyStart
	itemCommentStart
	itemInlineTableStart
	itemInlineTableEnd
)

const (
	eof              = 0
	comma            = ','
	tableStart       = '['
	tableEnd         = ']'
	arrayTableStart  = '['
	arrayTableEnd    = ']'
	tableSep         = '.'
	keySep           = '='
	arrayStart       = '['
	arrayEnd         = ']'
	commentStart     = '#'
	stringStart      = '"'
	stringEnd        = '"'
	rawStringStart   = '\''
	rawStringEnd     = '\''
	inlineTableStart = '{'
	inlineTableEnd   = '}'
)

type stateFn func(lx *lexer) stateFn

type lexer struct {
	input string
	start int
	pos   int
	line  int
	state stateFn
	items chan item

	
	
	prevWidths [3]int
	nprev      int 
	
	
	atEOF bool

	
	
	
	
	
	stack []stateFn
}

type item struct {
	typ  itemType
	val  string
	line int
}

func (lx *lexer) nextItem() item {
	for {
		select {
		case item := <-lx.items:
			return item
		default:
			lx.state = lx.state(lx)
		}
	}
}

func lex(input string) *lexer {
	lx := &lexer{
		input: input,
		state: lexTop,
		line:  1,
		items: make(chan item, 10),
		stack: make([]stateFn, 0, 10),
	}
	return lx
}

func (lx *lexer) push(state stateFn) {
	lx.stack = append(lx.stack, state)
}

func (lx *lexer) pop() stateFn {
	if len(lx.stack) == 0 {
		return lx.errorf("BUG in lexer: no states to pop")
	}
	last := lx.stack[len(lx.stack)-1]
	lx.stack = lx.stack[0 : len(lx.stack)-1]
	return last
}

func (lx *lexer) current() string {
	return lx.input[lx.start:lx.pos]
}

func (lx *lexer) emit(typ itemType) {
	lx.items <- item{typ, lx.current(), lx.line}
	lx.start = lx.pos
}

func (lx *lexer) emitTrim(typ itemType) {
	lx.items <- item{typ, strings.TrimSpace(lx.current()), lx.line}
	lx.start = lx.pos
}

func (lx *lexer) next() (r rune) {
	if lx.atEOF {
		panic("next called after EOF")
	}
	if lx.pos >= len(lx.input) {
		lx.atEOF = true
		return eof
	}

	if lx.input[lx.pos] == '\n' {
		lx.line++
	}
	lx.prevWidths[2] = lx.prevWidths[1]
	lx.prevWidths[1] = lx.prevWidths[0]
	if lx.nprev < 3 {
		lx.nprev++
	}
	r, w := utf8.DecodeRuneInString(lx.input[lx.pos:])
	lx.prevWidths[0] = w
	lx.pos += w
	return r
}


func (lx *lexer) ignore() {
	lx.start = lx.pos
}


func (lx *lexer) backup() {
	if lx.atEOF {
		lx.atEOF = false
		return
	}
	if lx.nprev < 1 {
		panic("backed up too far")
	}
	w := lx.prevWidths[0]
	lx.prevWidths[0] = lx.prevWidths[1]
	lx.prevWidths[1] = lx.prevWidths[2]
	lx.nprev--
	lx.pos -= w
	if lx.pos < len(lx.input) && lx.input[lx.pos] == '\n' {
		lx.line--
	}
}


func (lx *lexer) accept(valid rune) bool {
	if lx.next() == valid {
		return true
	}
	lx.backup()
	return false
}


func (lx *lexer) peek() rune {
	r := lx.next()
	lx.backup()
	return r
}


func (lx *lexer) skip(pred func(rune) bool) {
	for {
		r := lx.next()
		if pred(r) {
			continue
		}
		lx.backup()
		lx.ignore()
		return
	}
}




func (lx *lexer) errorf(format string, values ...interface{}) stateFn {
	lx.items <- item{
		itemError,
		fmt.Sprintf(format, values...),
		lx.line,
	}
	return nil
}


func lexTop(lx *lexer) stateFn {
	r := lx.next()
	if isWhitespace(r) || isNL(r) {
		return lexSkip(lx, lexTop)
	}
	switch r {
	case commentStart:
		lx.push(lexTop)
		return lexCommentStart
	case tableStart:
		return lexTableStart
	case eof:
		if lx.pos > lx.start {
			return lx.errorf("unexpected EOF")
		}
		lx.emit(itemEOF)
		return nil
	}

	
	
	lx.backup()
	lx.push(lexTopEnd)
	return lexKeyStart
}




func lexTopEnd(lx *lexer) stateFn {
	r := lx.next()
	switch {
	case r == commentStart:
		
		lx.push(lexTop)
		return lexCommentStart
	case isWhitespace(r):
		return lexTopEnd
	case isNL(r):
		lx.ignore()
		return lexTop
	case r == eof:
		lx.emit(itemEOF)
		return nil
	}
	return lx.errorf("expected a top-level item to end with a newline, "+
		"comment, or EOF, but got %q instead", r)
}






func lexTableStart(lx *lexer) stateFn {
	if lx.peek() == arrayTableStart {
		lx.next()
		lx.emit(itemArrayTableStart)
		lx.push(lexArrayTableEnd)
	} else {
		lx.emit(itemTableStart)
		lx.push(lexTableEnd)
	}
	return lexTableNameStart
}

func lexTableEnd(lx *lexer) stateFn {
	lx.emit(itemTableEnd)
	return lexTopEnd
}

func lexArrayTableEnd(lx *lexer) stateFn {
	if r := lx.next(); r != arrayTableEnd {
		return lx.errorf("expected end of table array name delimiter %q, "+
			"but got %q instead", arrayTableEnd, r)
	}
	lx.emit(itemArrayTableEnd)
	return lexTopEnd
}

func lexTableNameStart(lx *lexer) stateFn {
	lx.skip(isWhitespace)
	switch r := lx.peek(); {
	case r == tableEnd || r == eof:
		return lx.errorf("unexpected end of table name " +
			"(table names cannot be empty)")
	case r == tableSep:
		return lx.errorf("unexpected table separator " +
			"(table names cannot be empty)")
	case r == stringStart || r == rawStringStart:
		lx.ignore()
		lx.push(lexTableNameEnd)
		return lexValue 
	default:
		return lexBareTableName
	}
}



func lexBareTableName(lx *lexer) stateFn {
	r := lx.next()
	if isBareKeyChar(r) {
		return lexBareTableName
	}
	lx.backup()
	lx.emit(itemText)
	return lexTableNameEnd
}



func lexTableNameEnd(lx *lexer) stateFn {
	lx.skip(isWhitespace)
	switch r := lx.next(); {
	case isWhitespace(r):
		return lexTableNameEnd
	case r == tableSep:
		lx.ignore()
		return lexTableNameStart
	case r == tableEnd:
		return lx.pop()
	default:
		return lx.errorf("expected '.' or ']' to end table name, "+
			"but got %q instead", r)
	}
}



func lexKeyStart(lx *lexer) stateFn {
	r := lx.peek()
	switch {
	case r == keySep:
		return lx.errorf("unexpected key separator %q", keySep)
	case isWhitespace(r) || isNL(r):
		lx.next()
		return lexSkip(lx, lexKeyStart)
	case r == stringStart || r == rawStringStart:
		lx.ignore()
		lx.emit(itemKeyStart)
		lx.push(lexKeyEnd)
		return lexValue 
	default:
		lx.ignore()
		lx.emit(itemKeyStart)
		return lexBareKey
	}
}



func lexBareKey(lx *lexer) stateFn {
	switch r := lx.next(); {
	case isBareKeyChar(r):
		return lexBareKey
	case isWhitespace(r):
		lx.backup()
		lx.emit(itemText)
		return lexKeyEnd
	case r == keySep:
		lx.backup()
		lx.emit(itemText)
		return lexKeyEnd
	default:
		return lx.errorf("bare keys cannot contain %q", r)
	}
}



func lexKeyEnd(lx *lexer) stateFn {
	switch r := lx.next(); {
	case r == keySep:
		return lexSkip(lx, lexValue)
	case isWhitespace(r):
		return lexSkip(lx, lexKeyEnd)
	default:
		return lx.errorf("expected key separator %q, but got %q instead",
			keySep, r)
	}
}




func lexValue(lx *lexer) stateFn {
	
	
	r := lx.next()
	switch {
	case isWhitespace(r):
		return lexSkip(lx, lexValue)
	case isDigit(r):
		lx.backup() 
		return lexNumberOrDateStart
	}
	switch r {
	case arrayStart:
		lx.ignore()
		lx.emit(itemArray)
		return lexArrayValue
	case inlineTableStart:
		lx.ignore()
		lx.emit(itemInlineTableStart)
		return lexInlineTableValue
	case stringStart:
		if lx.accept(stringStart) {
			if lx.accept(stringStart) {
				lx.ignore() 
				return lexMultilineString
			}
			lx.backup()
		}
		lx.ignore() 
		return lexString
	case rawStringStart:
		if lx.accept(rawStringStart) {
			if lx.accept(rawStringStart) {
				lx.ignore() 
				return lexMultilineRawString
			}
			lx.backup()
		}
		lx.ignore() 
		return lexRawString
	case '+', '-':
		return lexNumberStart
	case '.': 
		return lx.errorf("floats must start with a digit, not '.'")
	}
	if unicode.IsLetter(r) {
		
		
		
		
		lx.backup()
		return lexBool
	}
	return lx.errorf("expected value but found %q instead", r)
}



func lexArrayValue(lx *lexer) stateFn {
	r := lx.next()
	switch {
	case isWhitespace(r) || isNL(r):
		return lexSkip(lx, lexArrayValue)
	case r == commentStart:
		lx.push(lexArrayValue)
		return lexCommentStart
	case r == comma:
		return lx.errorf("unexpected comma")
	case r == arrayEnd:
		
		
		return lexArrayEnd
	}

	lx.backup()
	lx.push(lexArrayValueEnd)
	return lexValue
}




func lexArrayValueEnd(lx *lexer) stateFn {
	r := lx.next()
	switch {
	case isWhitespace(r) || isNL(r):
		return lexSkip(lx, lexArrayValueEnd)
	case r == commentStart:
		lx.push(lexArrayValueEnd)
		return lexCommentStart
	case r == comma:
		lx.ignore()
		return lexArrayValue 
	case r == arrayEnd:
		return lexArrayEnd
	}
	return lx.errorf(
		"expected a comma or array terminator %q, but got %q instead",
		arrayEnd, r,
	)
}



func lexArrayEnd(lx *lexer) stateFn {
	lx.ignore()
	lx.emit(itemArrayEnd)
	return lx.pop()
}



func lexInlineTableValue(lx *lexer) stateFn {
	r := lx.next()
	switch {
	case isWhitespace(r):
		return lexSkip(lx, lexInlineTableValue)
	case isNL(r):
		return lx.errorf("newlines not allowed within inline tables")
	case r == commentStart:
		lx.push(lexInlineTableValue)
		return lexCommentStart
	case r == comma:
		return lx.errorf("unexpected comma")
	case r == inlineTableEnd:
		return lexInlineTableEnd
	}
	lx.backup()
	lx.push(lexInlineTableValueEnd)
	return lexKeyStart
}




func lexInlineTableValueEnd(lx *lexer) stateFn {
	r := lx.next()
	switch {
	case isWhitespace(r):
		return lexSkip(lx, lexInlineTableValueEnd)
	case isNL(r):
		return lx.errorf("newlines not allowed within inline tables")
	case r == commentStart:
		lx.push(lexInlineTableValueEnd)
		return lexCommentStart
	case r == comma:
		lx.ignore()
		return lexInlineTableValue
	case r == inlineTableEnd:
		return lexInlineTableEnd
	}
	return lx.errorf("expected a comma or an inline table terminator %q, "+
		"but got %q instead", inlineTableEnd, r)
}



func lexInlineTableEnd(lx *lexer) stateFn {
	lx.ignore()
	lx.emit(itemInlineTableEnd)
	return lx.pop()
}



func lexString(lx *lexer) stateFn {
	r := lx.next()
	switch {
	case r == eof:
		return lx.errorf("unexpected EOF")
	case isNL(r):
		return lx.errorf("strings cannot contain newlines")
	case r == '\\':
		lx.push(lexString)
		return lexStringEscape
	case r == stringEnd:
		lx.backup()
		lx.emit(itemString)
		lx.next()
		lx.ignore()
		return lx.pop()
	}
	return lexString
}



func lexMultilineString(lx *lexer) stateFn {
	switch lx.next() {
	case eof:
		return lx.errorf("unexpected EOF")
	case '\\':
		return lexMultilineStringEscape
	case stringEnd:
		if lx.accept(stringEnd) {
			if lx.accept(stringEnd) {
				lx.backup()
				lx.backup()
				lx.backup()
				lx.emit(itemMultilineString)
				lx.next()
				lx.next()
				lx.next()
				lx.ignore()
				return lx.pop()
			}
			lx.backup()
		}
	}
	return lexMultilineString
}



func lexRawString(lx *lexer) stateFn {
	r := lx.next()
	switch {
	case r == eof:
		return lx.errorf("unexpected EOF")
	case isNL(r):
		return lx.errorf("strings cannot contain newlines")
	case r == rawStringEnd:
		lx.backup()
		lx.emit(itemRawString)
		lx.next()
		lx.ignore()
		return lx.pop()
	}
	return lexRawString
}




func lexMultilineRawString(lx *lexer) stateFn {
	switch lx.next() {
	case eof:
		return lx.errorf("unexpected EOF")
	case rawStringEnd:
		if lx.accept(rawStringEnd) {
			if lx.accept(rawStringEnd) {
				lx.backup()
				lx.backup()
				lx.backup()
				lx.emit(itemRawMultilineString)
				lx.next()
				lx.next()
				lx.next()
				lx.ignore()
				return lx.pop()
			}
			lx.backup()
		}
	}
	return lexMultilineRawString
}



func lexMultilineStringEscape(lx *lexer) stateFn {
	
	if isNL(lx.next()) {
		return lexMultilineString
	}
	lx.backup()
	lx.push(lexMultilineString)
	return lexStringEscape(lx)
}

func lexStringEscape(lx *lexer) stateFn {
	r := lx.next()
	switch r {
	case 'b':
		fallthrough
	case 't':
		fallthrough
	case 'n':
		fallthrough
	case 'f':
		fallthrough
	case 'r':
		fallthrough
	case '"':
		fallthrough
	case '\\':
		return lx.pop()
	case 'u':
		return lexShortUnicodeEscape
	case 'U':
		return lexLongUnicodeEscape
	}
	return lx.errorf("invalid escape character %q; only the following "+
		"escape characters are allowed: "+
		`\b, \t, \n, \f, \r, \", \\, \uXXXX, and \UXXXXXXXX`, r)
}

func lexShortUnicodeEscape(lx *lexer) stateFn {
	var r rune
	for i := 0; i < 4; i++ {
		r = lx.next()
		if !isHexadecimal(r) {
			return lx.errorf(`expected four hexadecimal digits after '\u', `+
				"but got %q instead", lx.current())
		}
	}
	return lx.pop()
}

func lexLongUnicodeEscape(lx *lexer) stateFn {
	var r rune
	for i := 0; i < 8; i++ {
		r = lx.next()
		if !isHexadecimal(r) {
			return lx.errorf(`expected eight hexadecimal digits after '\U', `+
				"but got %q instead", lx.current())
		}
	}
	return lx.pop()
}


func lexNumberOrDateStart(lx *lexer) stateFn {
	r := lx.next()
	if isDigit(r) {
		return lexNumberOrDate
	}
	switch r {
	case '_':
		return lexNumber
	case 'e', 'E':
		return lexFloat
	case '.':
		return lx.errorf("floats must start with a digit, not '.'")
	}
	return lx.errorf("expected a digit but got %q", r)
}


func lexNumberOrDate(lx *lexer) stateFn {
	r := lx.next()
	if isDigit(r) {
		return lexNumberOrDate
	}
	switch r {
	case '-':
		return lexDatetime
	case '_':
		return lexNumber
	case '.', 'e', 'E':
		return lexFloat
	}

	lx.backup()
	lx.emit(itemInteger)
	return lx.pop()
}



func lexDatetime(lx *lexer) stateFn {
	r := lx.next()
	if isDigit(r) {
		return lexDatetime
	}
	switch r {
	case '-', 'T', ':', '.', 'Z':
		return lexDatetime
	}

	lx.backup()
	lx.emit(itemDatetime)
	return lx.pop()
}




func lexNumberStart(lx *lexer) stateFn {
	
	r := lx.next()
	if !isDigit(r) {
		if r == '.' {
			return lx.errorf("floats must start with a digit, not '.'")
		}
		return lx.errorf("expected a digit but got %q", r)
	}
	return lexNumber
}


func lexNumber(lx *lexer) stateFn {
	r := lx.next()
	if isDigit(r) {
		return lexNumber
	}
	switch r {
	case '_':
		return lexNumber
	case '.', 'e', 'E':
		return lexFloat
	}

	lx.backup()
	lx.emit(itemInteger)
	return lx.pop()
}




func lexFloat(lx *lexer) stateFn {
	r := lx.next()
	if isDigit(r) {
		return lexFloat
	}
	switch r {
	case '_', '.', '-', '+', 'e', 'E':
		return lexFloat
	}

	lx.backup()
	lx.emit(itemFloat)
	return lx.pop()
}


func lexBool(lx *lexer) stateFn {
	var rs []rune
	for {
		r := lx.next()
		if !unicode.IsLetter(r) {
			lx.backup()
			break
		}
		rs = append(rs, r)
	}
	s := string(rs)
	switch s {
	case "true", "false":
		lx.emit(itemBool)
		return lx.pop()
	}
	return lx.errorf("expected value but found %q instead", s)
}



func lexCommentStart(lx *lexer) stateFn {
	lx.ignore()
	lx.emit(itemCommentStart)
	return lexComment
}




func lexComment(lx *lexer) stateFn {
	r := lx.peek()
	if isNL(r) || r == eof {
		lx.emit(itemText)
		return lx.pop()
	}
	lx.next()
	return lexComment
}


func lexSkip(lx *lexer, nextState stateFn) stateFn {
	return func(lx *lexer) stateFn {
		lx.ignore()
		return nextState
	}
}



func isWhitespace(r rune) bool {
	return r == '\t' || r == ' '
}

func isNL(r rune) bool {
	return r == '\n' || r == '\r'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexadecimal(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

func isBareKeyChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') ||
		r == '_' ||
		r == '-'
}

func (itype itemType) String() string {
	switch itype {
	case itemError:
		return "Error"
	case itemNIL:
		return "NIL"
	case itemEOF:
		return "EOF"
	case itemText:
		return "Text"
	case itemString, itemRawString, itemMultilineString, itemRawMultilineString:
		return "String"
	case itemBool:
		return "Bool"
	case itemInteger:
		return "Integer"
	case itemFloat:
		return "Float"
	case itemDatetime:
		return "DateTime"
	case itemTableStart:
		return "TableStart"
	case itemTableEnd:
		return "TableEnd"
	case itemKeyStart:
		return "KeyStart"
	case itemArray:
		return "Array"
	case itemArrayEnd:
		return "ArrayEnd"
	case itemCommentStart:
		return "CommentStart"
	}
	panic(fmt.Sprintf("BUG: Unknown type '%d'.", int(itype)))
}

func (item item) String() string {
	return fmt.Sprintf("(%s, %s)", item.typ.String(), item.val)
}
