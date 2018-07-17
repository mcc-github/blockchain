











package properties

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)


type item struct {
	typ itemType 
	pos int      
	val string   
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}


type itemType int

const (
	itemError itemType = iota 
	itemEOF
	itemKey     
	itemValue   
	itemComment 
)


const eof = -1


const whitespace = " \f\t"


type stateFn func(*lexer) stateFn


type lexer struct {
	input   string    
	state   stateFn   
	pos     int       
	start   int       
	width   int       
	lastPos int       
	runes   []rune    
	items   chan item 
}


func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}


func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}


func (l *lexer) backup() {
	l.pos -= l.width
}


func (l *lexer) emit(t itemType) {
	i := item{t, l.start, string(l.runes)}
	l.items <- i
	l.start = l.pos
	l.runes = l.runes[:0]
}


func (l *lexer) ignore() {
	l.start = l.pos
}


func (l *lexer) appendRune(r rune) {
	l.runes = append(l.runes, r)
}


func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}


func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}


func (l *lexer) acceptRunUntil(term rune) {
	for term != l.next() {
	}
	l.backup()
}


func (l *lexer) isNotEmpty() bool {
	return l.pos > l.start
}




func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}



func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}


func (l *lexer) nextItem() item {
	i := <-l.items
	l.lastPos = i.pos
	return i
}


func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
		runes: make([]rune, 0, 32),
	}
	go l.run()
	return l
}


func (l *lexer) run() {
	for l.state = lexBeforeKey(l); l.state != nil; {
		l.state = l.state(l)
	}
}




func lexBeforeKey(l *lexer) stateFn {
	switch r := l.next(); {
	case isEOF(r):
		l.emit(itemEOF)
		return nil

	case isEOL(r):
		l.ignore()
		return lexBeforeKey

	case isComment(r):
		return lexComment

	case isWhitespace(r):
		l.ignore()
		return lexBeforeKey

	default:
		l.backup()
		return lexKey
	}
}


func lexComment(l *lexer) stateFn {
	l.acceptRun(whitespace)
	l.ignore()
	for {
		switch r := l.next(); {
		case isEOF(r):
			l.ignore()
			l.emit(itemEOF)
			return nil
		case isEOL(r):
			l.emit(itemComment)
			return lexBeforeKey
		default:
			l.appendRune(r)
		}
	}
}


func lexKey(l *lexer) stateFn {
	var r rune

Loop:
	for {
		switch r = l.next(); {

		case isEscape(r):
			err := l.scanEscapeSequence()
			if err != nil {
				return l.errorf(err.Error())
			}

		case isEndOfKey(r):
			l.backup()
			break Loop

		case isEOF(r):
			break Loop

		default:
			l.appendRune(r)
		}
	}

	if len(l.runes) > 0 {
		l.emit(itemKey)
	}

	if isEOF(r) {
		l.emit(itemEOF)
		return nil
	}

	return lexBeforeValue
}




func lexBeforeValue(l *lexer) stateFn {
	l.acceptRun(whitespace)
	l.accept(":=")
	l.acceptRun(whitespace)
	l.ignore()
	return lexValue
}


func lexValue(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isEscape(r):
			if isEOL(l.peek()) {
				l.next()
				l.acceptRun(whitespace)
			} else {
				err := l.scanEscapeSequence()
				if err != nil {
					return l.errorf(err.Error())
				}
			}

		case isEOL(r):
			l.emit(itemValue)
			l.ignore()
			return lexBeforeKey

		case isEOF(r):
			l.emit(itemValue)
			l.emit(itemEOF)
			return nil

		default:
			l.appendRune(r)
		}
	}
}



func (l *lexer) scanEscapeSequence() error {
	switch r := l.next(); {

	case isEscapedCharacter(r):
		l.appendRune(decodeEscapedCharacter(r))
		return nil

	case atUnicodeLiteral(r):
		return l.scanUnicodeLiteral()

	case isEOF(r):
		return fmt.Errorf("premature EOF")

	
	default:
		l.appendRune(r)
		return nil
	}
}


func (l *lexer) scanUnicodeLiteral() error {
	
	d := make([]rune, 4)
	for i := 0; i < 4; i++ {
		d[i] = l.next()
		if d[i] == eof || !strings.ContainsRune("0123456789abcdefABCDEF", d[i]) {
			return fmt.Errorf("invalid unicode literal")
		}
	}

	
	r, err := strconv.ParseInt(string(d), 16, 0)
	if err != nil {
		return err
	}

	l.appendRune(rune(r))
	return nil
}


func decodeEscapedCharacter(r rune) rune {
	switch r {
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	default:
		return r
	}
}



func atUnicodeLiteral(r rune) bool {
	return r == 'u'
}


func isComment(r rune) bool {
	return r == '#' || r == '!'
}


func isEndOfKey(r rune) bool {
	return strings.ContainsRune(" \f\t\r\n:=", r)
}


func isEOF(r rune) bool {
	return r == eof
}


func isEOL(r rune) bool {
	return r == '\n' || r == '\r'
}



func isEscape(r rune) bool {
	return r == '\\'
}



func isEscapedCharacter(r rune) bool {
	return strings.ContainsRune(" :=fnrt", r)
}


func isWhitespace(r rune) bool {
	return strings.ContainsRune(whitespace, r)
}
