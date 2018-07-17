package toml

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type parser struct {
	mapping map[string]interface{}
	types   map[string]tomlType
	lx      *lexer

	
	ordered []Key

	
	context Key

	
	currentKey string

	
	approxLine int

	
	implicits map[string]bool
}

type parseError string

func (pe parseError) Error() string {
	return string(pe)
}

func parse(data string) (p *parser, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(parseError); ok {
				return
			}
			panic(r)
		}
	}()

	p = &parser{
		mapping:   make(map[string]interface{}),
		types:     make(map[string]tomlType),
		lx:        lex(data),
		ordered:   make([]Key, 0),
		implicits: make(map[string]bool),
	}
	for {
		item := p.next()
		if item.typ == itemEOF {
			break
		}
		p.topLevel(item)
	}

	return p, nil
}

func (p *parser) panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf("Near line %d (last key parsed '%s'): %s",
		p.approxLine, p.current(), fmt.Sprintf(format, v...))
	panic(parseError(msg))
}

func (p *parser) next() item {
	it := p.lx.nextItem()
	if it.typ == itemError {
		p.panicf("%s", it.val)
	}
	return it
}

func (p *parser) bug(format string, v ...interface{}) {
	panic(fmt.Sprintf("BUG: "+format+"\n\n", v...))
}

func (p *parser) expect(typ itemType) item {
	it := p.next()
	p.assertEqual(typ, it.typ)
	return it
}

func (p *parser) assertEqual(expected, got itemType) {
	if expected != got {
		p.bug("Expected '%s' but got '%s'.", expected, got)
	}
}

func (p *parser) topLevel(item item) {
	switch item.typ {
	case itemCommentStart:
		p.approxLine = item.line
		p.expect(itemText)
	case itemTableStart:
		kg := p.next()
		p.approxLine = kg.line

		var key Key
		for ; kg.typ != itemTableEnd && kg.typ != itemEOF; kg = p.next() {
			key = append(key, p.keyString(kg))
		}
		p.assertEqual(itemTableEnd, kg.typ)

		p.establishContext(key, false)
		p.setType("", tomlHash)
		p.ordered = append(p.ordered, key)
	case itemArrayTableStart:
		kg := p.next()
		p.approxLine = kg.line

		var key Key
		for ; kg.typ != itemArrayTableEnd && kg.typ != itemEOF; kg = p.next() {
			key = append(key, p.keyString(kg))
		}
		p.assertEqual(itemArrayTableEnd, kg.typ)

		p.establishContext(key, true)
		p.setType("", tomlArrayHash)
		p.ordered = append(p.ordered, key)
	case itemKeyStart:
		kname := p.next()
		p.approxLine = kname.line
		p.currentKey = p.keyString(kname)

		val, typ := p.value(p.next())
		p.setValue(p.currentKey, val)
		p.setType(p.currentKey, typ)
		p.ordered = append(p.ordered, p.context.add(p.currentKey))
		p.currentKey = ""
	default:
		p.bug("Unexpected type at top level: %s", item.typ)
	}
}


func (p *parser) keyString(it item) string {
	switch it.typ {
	case itemText:
		return it.val
	case itemString, itemMultilineString,
		itemRawString, itemRawMultilineString:
		s, _ := p.value(it)
		return s.(string)
	default:
		p.bug("Unexpected key type: %s", it.typ)
		panic("unreachable")
	}
}



func (p *parser) value(it item) (interface{}, tomlType) {
	switch it.typ {
	case itemString:
		return p.replaceEscapes(it.val), p.typeOfPrimitive(it)
	case itemMultilineString:
		trimmed := stripFirstNewline(stripEscapedWhitespace(it.val))
		return p.replaceEscapes(trimmed), p.typeOfPrimitive(it)
	case itemRawString:
		return it.val, p.typeOfPrimitive(it)
	case itemRawMultilineString:
		return stripFirstNewline(it.val), p.typeOfPrimitive(it)
	case itemBool:
		switch it.val {
		case "true":
			return true, p.typeOfPrimitive(it)
		case "false":
			return false, p.typeOfPrimitive(it)
		}
		p.bug("Expected boolean value, but got '%s'.", it.val)
	case itemInteger:
		if !numUnderscoresOK(it.val) {
			p.panicf("Invalid integer %q: underscores must be surrounded by digits",
				it.val)
		}
		val := strings.Replace(it.val, "_", "", -1)
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			
			
			
			
			
			if e, ok := err.(*strconv.NumError); ok &&
				e.Err == strconv.ErrRange {

				p.panicf("Integer '%s' is out of the range of 64-bit "+
					"signed integers.", it.val)
			} else {
				p.bug("Expected integer value, but got '%s'.", it.val)
			}
		}
		return num, p.typeOfPrimitive(it)
	case itemFloat:
		parts := strings.FieldsFunc(it.val, func(r rune) bool {
			switch r {
			case '.', 'e', 'E':
				return true
			}
			return false
		})
		for _, part := range parts {
			if !numUnderscoresOK(part) {
				p.panicf("Invalid float %q: underscores must be "+
					"surrounded by digits", it.val)
			}
		}
		if !numPeriodsOK(it.val) {
			
			
			
			
			p.panicf("Invalid float %q: '.' must be followed "+
				"by one or more digits", it.val)
		}
		val := strings.Replace(it.val, "_", "", -1)
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			if e, ok := err.(*strconv.NumError); ok &&
				e.Err == strconv.ErrRange {

				p.panicf("Float '%s' is out of the range of 64-bit "+
					"IEEE-754 floating-point numbers.", it.val)
			} else {
				p.panicf("Invalid float value: %q", it.val)
			}
		}
		return num, p.typeOfPrimitive(it)
	case itemDatetime:
		var t time.Time
		var ok bool
		var err error
		for _, format := range []string{
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02",
		} {
			t, err = time.ParseInLocation(format, it.val, time.Local)
			if err == nil {
				ok = true
				break
			}
		}
		if !ok {
			p.panicf("Invalid TOML Datetime: %q.", it.val)
		}
		return t, p.typeOfPrimitive(it)
	case itemArray:
		array := make([]interface{}, 0)
		types := make([]tomlType, 0)

		for it = p.next(); it.typ != itemArrayEnd; it = p.next() {
			if it.typ == itemCommentStart {
				p.expect(itemText)
				continue
			}

			val, typ := p.value(it)
			array = append(array, val)
			types = append(types, typ)
		}
		return array, p.typeOfArray(types)
	case itemInlineTableStart:
		var (
			hash         = make(map[string]interface{})
			outerContext = p.context
			outerKey     = p.currentKey
		)

		p.context = append(p.context, p.currentKey)
		p.currentKey = ""
		for it := p.next(); it.typ != itemInlineTableEnd; it = p.next() {
			if it.typ != itemKeyStart {
				p.bug("Expected key start but instead found %q, around line %d",
					it.val, p.approxLine)
			}
			if it.typ == itemCommentStart {
				p.expect(itemText)
				continue
			}

			
			k := p.next()
			p.approxLine = k.line
			kname := p.keyString(k)

			
			p.currentKey = kname
			val, typ := p.value(p.next())
			
			p.setType(kname, typ)
			p.ordered = append(p.ordered, p.context.add(p.currentKey))
			hash[kname] = val
		}
		p.context = outerContext
		p.currentKey = outerKey
		return hash, tomlHash
	}
	p.bug("Unexpected value type: %s", it.typ)
	panic("unreachable")
}



func numUnderscoresOK(s string) bool {
	accept := false
	for _, r := range s {
		if r == '_' {
			if !accept {
				return false
			}
			accept = false
			continue
		}
		accept = true
	}
	return accept
}


func numPeriodsOK(s string) bool {
	period := false
	for _, r := range s {
		if period && !isDigit(r) {
			return false
		}
		period = r == '.'
	}
	return !period
}







func (p *parser) establishContext(key Key, array bool) {
	var ok bool

	
	hashContext := p.mapping
	keyContext := make(Key, 0)

	
	for _, k := range key[0 : len(key)-1] {
		_, ok = hashContext[k]
		keyContext = append(keyContext, k)

		
		if !ok {
			p.addImplicit(keyContext)
			hashContext[k] = make(map[string]interface{})
		}

		
		
		
		
		
		switch t := hashContext[k].(type) {
		case []map[string]interface{}:
			hashContext = t[len(t)-1]
		case map[string]interface{}:
			hashContext = t
		default:
			p.panicf("Key '%s' was already created as a hash.", keyContext)
		}
	}

	p.context = keyContext
	if array {
		
		
		k := key[len(key)-1]
		if _, ok := hashContext[k]; !ok {
			hashContext[k] = make([]map[string]interface{}, 0, 5)
		}

		
		
		if hash, ok := hashContext[k].([]map[string]interface{}); ok {
			hashContext[k] = append(hash, make(map[string]interface{}))
		} else {
			p.panicf("Key '%s' was already created and cannot be used as "+
				"an array.", keyContext)
		}
	} else {
		p.setValue(key[len(key)-1], make(map[string]interface{}))
	}
	p.context = append(p.context, key[len(key)-1])
}




func (p *parser) setValue(key string, value interface{}) {
	var tmpHash interface{}
	var ok bool

	hash := p.mapping
	keyContext := make(Key, 0)
	for _, k := range p.context {
		keyContext = append(keyContext, k)
		if tmpHash, ok = hash[k]; !ok {
			p.bug("Context for key '%s' has not been established.", keyContext)
		}
		switch t := tmpHash.(type) {
		case []map[string]interface{}:
			
			
			hash = t[len(t)-1]
		case map[string]interface{}:
			hash = t
		default:
			p.bug("Expected hash to have type 'map[string]interface{}', but "+
				"it has '%T' instead.", tmpHash)
		}
	}
	keyContext = append(keyContext, key)

	if _, ok := hash[key]; ok {
		
		
		
		
		
		
		
		
		
		
		
		if p.isImplicit(keyContext) {
			p.removeImplicit(keyContext)
			return
		}

		
		
		p.panicf("Key '%s' has already been defined.", keyContext)
	}
	hash[key] = value
}






func (p *parser) setType(key string, typ tomlType) {
	keyContext := make(Key, 0, len(p.context)+1)
	for _, k := range p.context {
		keyContext = append(keyContext, k)
	}
	if len(key) > 0 { 
		keyContext = append(keyContext, key)
	}
	p.types[keyContext.String()] = typ
}


func (p *parser) addImplicit(key Key) {
	p.implicits[key.String()] = true
}



func (p *parser) removeImplicit(key Key) {
	p.implicits[key.String()] = false
}



func (p *parser) isImplicit(key Key) bool {
	return p.implicits[key.String()]
}


func (p *parser) current() string {
	if len(p.currentKey) == 0 {
		return p.context.String()
	}
	if len(p.context) == 0 {
		return p.currentKey
	}
	return fmt.Sprintf("%s.%s", p.context, p.currentKey)
}

func stripFirstNewline(s string) string {
	if len(s) == 0 || s[0] != '\n' {
		return s
	}
	return s[1:]
}

func stripEscapedWhitespace(s string) string {
	esc := strings.Split(s, "\\\n")
	if len(esc) > 1 {
		for i := 1; i < len(esc); i++ {
			esc[i] = strings.TrimLeftFunc(esc[i], unicode.IsSpace)
		}
	}
	return strings.Join(esc, "")
}

func (p *parser) replaceEscapes(str string) string {
	var replaced []rune
	s := []byte(str)
	r := 0
	for r < len(s) {
		if s[r] != '\\' {
			c, size := utf8.DecodeRune(s[r:])
			r += size
			replaced = append(replaced, c)
			continue
		}
		r += 1
		if r >= len(s) {
			p.bug("Escape sequence at end of string.")
			return ""
		}
		switch s[r] {
		default:
			p.bug("Expected valid escape code after \\, but got %q.", s[r])
			return ""
		case 'b':
			replaced = append(replaced, rune(0x0008))
			r += 1
		case 't':
			replaced = append(replaced, rune(0x0009))
			r += 1
		case 'n':
			replaced = append(replaced, rune(0x000A))
			r += 1
		case 'f':
			replaced = append(replaced, rune(0x000C))
			r += 1
		case 'r':
			replaced = append(replaced, rune(0x000D))
			r += 1
		case '"':
			replaced = append(replaced, rune(0x0022))
			r += 1
		case '\\':
			replaced = append(replaced, rune(0x005C))
			r += 1
		case 'u':
			
			
			
			escaped := p.asciiEscapeToUnicode(s[r+1 : r+5])
			replaced = append(replaced, escaped)
			r += 5
		case 'U':
			
			
			
			escaped := p.asciiEscapeToUnicode(s[r+1 : r+9])
			replaced = append(replaced, escaped)
			r += 9
		}
	}
	return string(replaced)
}

func (p *parser) asciiEscapeToUnicode(bs []byte) rune {
	s := string(bs)
	hex, err := strconv.ParseUint(strings.ToLower(s), 16, 32)
	if err != nil {
		p.bug("Could not parse '%s' as a hexadecimal number, but the "+
			"lexer claims it's OK: %s", s, err)
	}
	if !utf8.ValidRune(rune(hex)) {
		p.panicf("Escaped character '\\u%s' is not valid UTF-8.", s)
	}
	return rune(hex)
}

func isStringType(ty itemType) bool {
	return ty == itemString || ty == itemMultilineString ||
		ty == itemRawString || ty == itemRawMultilineString
}
