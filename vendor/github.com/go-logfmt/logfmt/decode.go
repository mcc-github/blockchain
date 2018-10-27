package logfmt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)


type Decoder struct {
	pos     int
	key     []byte
	value   []byte
	lineNum int
	s       *bufio.Scanner
	err     error
}





func NewDecoder(r io.Reader) *Decoder {
	dec := &Decoder{
		s: bufio.NewScanner(r),
	}
	return dec
}






func (dec *Decoder) ScanRecord() bool {
	if dec.err != nil {
		return false
	}
	if !dec.s.Scan() {
		dec.err = dec.s.Err()
		return false
	}
	dec.lineNum++
	dec.pos = 0
	return true
}





func (dec *Decoder) ScanKeyval() bool {
	dec.key, dec.value = nil, nil
	if dec.err != nil {
		return false
	}

	line := dec.s.Bytes()

	
	for p, c := range line[dec.pos:] {
		if c > ' ' {
			dec.pos += p
			goto key
		}
	}
	dec.pos = len(line)
	return false

key:
	const invalidKeyError = "invalid key"

	start, multibyte := dec.pos, false
	for p, c := range line[dec.pos:] {
		switch {
		case c == '=':
			dec.pos += p
			if dec.pos > start {
				dec.key = line[start:dec.pos]
				if multibyte && bytes.IndexRune(dec.key, utf8.RuneError) != -1 {
					dec.syntaxError(invalidKeyError)
					return false
				}
			}
			if dec.key == nil {
				dec.unexpectedByte(c)
				return false
			}
			goto equal
		case c == '"':
			dec.pos += p
			dec.unexpectedByte(c)
			return false
		case c <= ' ':
			dec.pos += p
			if dec.pos > start {
				dec.key = line[start:dec.pos]
				if multibyte && bytes.IndexRune(dec.key, utf8.RuneError) != -1 {
					dec.syntaxError(invalidKeyError)
					return false
				}
			}
			return true
		case c >= utf8.RuneSelf:
			multibyte = true
		}
	}
	dec.pos = len(line)
	if dec.pos > start {
		dec.key = line[start:dec.pos]
		if multibyte && bytes.IndexRune(dec.key, utf8.RuneError) != -1 {
			dec.syntaxError(invalidKeyError)
			return false
		}
	}
	return true

equal:
	dec.pos++
	if dec.pos >= len(line) {
		return true
	}
	switch c := line[dec.pos]; {
	case c <= ' ':
		return true
	case c == '"':
		goto qvalue
	}

	
	start = dec.pos
	for p, c := range line[dec.pos:] {
		switch {
		case c == '=' || c == '"':
			dec.pos += p
			dec.unexpectedByte(c)
			return false
		case c <= ' ':
			dec.pos += p
			if dec.pos > start {
				dec.value = line[start:dec.pos]
			}
			return true
		}
	}
	dec.pos = len(line)
	if dec.pos > start {
		dec.value = line[start:dec.pos]
	}
	return true

qvalue:
	const (
		untermQuote  = "unterminated quoted value"
		invalidQuote = "invalid quoted value"
	)

	hasEsc, esc := false, false
	start = dec.pos
	for p, c := range line[dec.pos+1:] {
		switch {
		case esc:
			esc = false
		case c == '\\':
			hasEsc, esc = true, true
		case c == '"':
			dec.pos += p + 2
			if hasEsc {
				v, ok := unquoteBytes(line[start:dec.pos])
				if !ok {
					dec.syntaxError(invalidQuote)
					return false
				}
				dec.value = v
			} else {
				start++
				end := dec.pos - 1
				if end > start {
					dec.value = line[start:end]
				}
			}
			return true
		}
	}
	dec.pos = len(line)
	dec.syntaxError(untermQuote)
	return false
}




func (dec *Decoder) Key() []byte {
	return dec.key
}





func (dec *Decoder) Value() []byte {
	return dec.value
}


func (dec *Decoder) Err() error {
	return dec.err
}

func (dec *Decoder) syntaxError(msg string) {
	dec.err = &SyntaxError{
		Msg:  msg,
		Line: dec.lineNum,
		Pos:  dec.pos + 1,
	}
}

func (dec *Decoder) unexpectedByte(c byte) {
	dec.err = &SyntaxError{
		Msg:  fmt.Sprintf("unexpected %q", c),
		Line: dec.lineNum,
		Pos:  dec.pos + 1,
	}
}


type SyntaxError struct {
	Msg  string
	Line int
	Pos  int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("logfmt syntax error at pos %d on line %d: %s", e.Pos, e.Line, e.Msg)
}
