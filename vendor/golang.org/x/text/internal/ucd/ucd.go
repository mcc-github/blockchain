








package ucd 

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)


const (
	CodePoint = iota
	Name
	GeneralCategory
	CanonicalCombiningClass
	BidiClass
	DecompMapping
	DecimalValue
	DigitValue
	NumericValue
	BidiMirrored
	Unicode1Name
	ISOComment
	SimpleUppercaseMapping
	SimpleLowercaseMapping
	SimpleTitlecaseMapping
)





func Parse(r io.ReadCloser, f func(p *Parser)) {
	defer r.Close()

	p := New(r)
	for p.Next() {
		f(p)
	}
	if err := p.Err(); err != nil {
		r.Close() 
		log.Fatal(err)
	}
}


type Option func(p *Parser)

func keepRanges(p *Parser) {
	p.keepRanges = true
}

var (
	
	
	KeepRanges Option = keepRanges
)



func Part(f func(p *Parser)) Option {
	return func(p *Parser) {
		p.partHandler = f
	}
}



func CommentHandler(f func(s string)) Option {
	return func(p *Parser) {
		p.commentHandler = f
	}
}


type Parser struct {
	scanner *bufio.Scanner

	keepRanges bool 

	err     error
	comment string
	field   []string
	
	
	line                 int
	parsedRange          bool
	rangeStart, rangeEnd rune

	partHandler    func(p *Parser)
	commentHandler func(s string)
}

func (p *Parser) setError(err error, msg string) {
	if p.err == nil && err != nil {
		if msg == "" {
			p.err = fmt.Errorf("ucd:line:%d: %v", p.line, err)
		} else {
			p.err = fmt.Errorf("ucd:line:%d:%s: %v", p.line, msg, err)
		}
	}
}

func (p *Parser) getField(i int) string {
	if i >= len(p.field) {
		return ""
	}
	return p.field[i]
}


func (p *Parser) Err() error {
	return p.err
}


func New(r io.Reader, o ...Option) *Parser {
	p := &Parser{
		scanner: bufio.NewScanner(r),
	}
	for _, f := range o {
		f(p)
	}
	return p
}



func (p *Parser) Next() bool {
	if !p.keepRanges && p.rangeStart < p.rangeEnd {
		p.rangeStart++
		return true
	}
	p.comment = ""
	p.field = p.field[:0]
	p.parsedRange = false

	for p.scanner.Scan() && p.err == nil {
		p.line++
		s := p.scanner.Text()
		if s == "" {
			continue
		}
		if s[0] == '#' {
			if p.commentHandler != nil {
				p.commentHandler(strings.TrimSpace(s[1:]))
			}
			continue
		}

		
		if i := strings.IndexByte(s, '#'); i != -1 {
			p.comment = strings.TrimSpace(s[i+1:])
			s = s[:i]
		}
		if s[0] == '@' {
			if p.partHandler != nil {
				p.field = append(p.field, strings.TrimSpace(s[1:]))
				p.partHandler(p)
				p.field = p.field[:0]
			}
			p.comment = ""
			continue
		}
		for {
			i := strings.IndexByte(s, ';')
			if i == -1 {
				p.field = append(p.field, strings.TrimSpace(s))
				break
			}
			p.field = append(p.field, strings.TrimSpace(s[:i]))
			s = s[i+1:]
		}
		if !p.keepRanges {
			p.rangeStart, p.rangeEnd = p.getRange(0)
		}
		return true
	}
	p.setError(p.scanner.Err(), "scanner failed")
	return false
}

func parseRune(b string) (rune, error) {
	if len(b) > 2 && b[0] == 'U' && b[1] == '+' {
		b = b[2:]
	}
	x, err := strconv.ParseUint(b, 16, 32)
	return rune(x), err
}

func (p *Parser) parseRune(s string) rune {
	x, err := parseRune(s)
	p.setError(err, "failed to parse rune")
	return x
}


func (p *Parser) Rune(i int) rune {
	if i > 0 || p.keepRanges {
		return p.parseRune(p.getField(i))
	}
	return p.rangeStart
}


func (p *Parser) Runes(i int) (runes []rune) {
	add := func(s string) {
		if s = strings.TrimSpace(s); len(s) > 0 {
			runes = append(runes, p.parseRune(s))
		}
	}
	for b := p.getField(i); ; {
		i := strings.IndexByte(b, ' ')
		if i == -1 {
			add(b)
			break
		}
		add(b[:i])
		b = b[i+1:]
	}
	return
}

var (
	errIncorrectLegacyRange = errors.New("ucd: unmatched <* First>")

	
	reRange = regexp.MustCompile("^([0-9A-F]*);<([^,]*), ([^>]*)>(.*)$")
)




func (p *Parser) Range(i int) (first, last rune) {
	if !p.keepRanges {
		return p.rangeStart, p.rangeStart
	}
	return p.getRange(i)
}

func (p *Parser) getRange(i int) (first, last rune) {
	b := p.getField(i)
	if k := strings.Index(b, ".."); k != -1 {
		return p.parseRune(b[:k]), p.parseRune(b[k+2:])
	}
	
	
	x, err := parseRune(b)
	if err != nil {
		
		
		
		p.keepRanges = true
	}
	
	if i == 0 && len(p.field) > 1 && strings.HasSuffix(p.field[1], "First>") {
		if p.parsedRange {
			return p.rangeStart, p.rangeEnd
		}
		mf := reRange.FindStringSubmatch(p.scanner.Text())
		p.line++
		if mf == nil || !p.scanner.Scan() {
			p.setError(errIncorrectLegacyRange, "")
			return x, x
		}
		
		
		ml := reRange.FindStringSubmatch(p.scanner.Text())
		if ml == nil || mf[2] != ml[2] || ml[3] != "Last" || mf[4] != ml[4] {
			p.setError(errIncorrectLegacyRange, "")
			return x, x
		}
		p.rangeStart, p.rangeEnd = x, p.parseRune(p.scanner.Text()[:len(ml[1])])
		p.parsedRange = true
		return p.rangeStart, p.rangeEnd
	}
	return x, x
}


var bools = map[string]bool{
	"":      false,
	"N":     false,
	"No":    false,
	"F":     false,
	"False": false,
	"Y":     true,
	"Yes":   true,
	"T":     true,
	"True":  true,
}


func (p *Parser) Bool(i int) bool {
	f := p.getField(i)
	for s, v := range bools {
		if f == s {
			return v
		}
	}
	p.setError(strconv.ErrSyntax, "error parsing bool")
	return false
}


func (p *Parser) Int(i int) int {
	x, err := strconv.ParseInt(string(p.getField(i)), 10, 64)
	p.setError(err, "error parsing int")
	return int(x)
}


func (p *Parser) Uint(i int) uint {
	x, err := strconv.ParseUint(string(p.getField(i)), 10, 64)
	p.setError(err, "error parsing uint")
	return uint(x)
}


func (p *Parser) Float(i int) float64 {
	x, err := strconv.ParseFloat(string(p.getField(i)), 64)
	p.setError(err, "error parsing float")
	return x
}


func (p *Parser) String(i int) string {
	return string(p.getField(i))
}


func (p *Parser) Strings(i int) []string {
	ss := strings.Split(string(p.getField(i)), " ")
	for i, s := range ss {
		ss[i] = strings.TrimSpace(s)
	}
	return ss
}


func (p *Parser) Comment() string {
	return string(p.comment)
}

var errUndefinedEnum = errors.New("ucd: undefined enum value")



func (p *Parser) Enum(i int, enum ...string) string {
	f := p.getField(i)
	for _, s := range enum {
		if f == s {
			return s
		}
	}
	p.setError(errUndefinedEnum, "error parsing enum")
	return ""
}
