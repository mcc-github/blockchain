



package properties




import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const maxExpansionDepth = 64




type ErrorHandlerFunc func(error)



var ErrorHandler ErrorHandlerFunc = LogFatalHandler


type LogHandlerFunc func(fmt string, args ...interface{})


var LogPrintf LogHandlerFunc = log.Printf


func LogFatalHandler(err error) {
	log.Fatal(err)
}


func PanicHandler(err error) {
	panic(err)
}





type Properties struct {
	
	Prefix  string
	Postfix string

	
	
	
	
	DisableExpansion bool

	
	m map[string]string

	
	c map[string][]string

	
	k []string
}



func NewProperties() *Properties {
	return &Properties{
		Prefix:  "${",
		Postfix: "}",
		m:       map[string]string{},
		c:       map[string][]string{},
		k:       []string{},
	}
}


func (p *Properties) Load(buf []byte, enc Encoding) error {
	l := &Loader{Encoding: enc, DisableExpansion: p.DisableExpansion}
	newProperties, err := l.LoadBytes(buf)
	if err != nil {
		return err
	}
	p.Merge(newProperties)
	return nil
}



func (p *Properties) Get(key string) (value string, ok bool) {
	v, ok := p.m[key]
	if p.DisableExpansion {
		return v, ok
	}
	if !ok {
		return "", false
	}

	expanded, err := p.expand(key, v)

	
	
	
	if err != nil {
		ErrorHandler(fmt.Errorf("%s in %q", err, key+" = "+v))
	}

	return expanded, true
}



func (p *Properties) MustGet(key string) string {
	if v, ok := p.Get(key); ok {
		return v
	}
	ErrorHandler(invalidKeyError(key))
	panic("ErrorHandler should exit")
}




func (p *Properties) ClearComments() {
	p.c = map[string][]string{}
}




func (p *Properties) GetComment(key string) string {
	comments, ok := p.c[key]
	if !ok || len(comments) == 0 {
		return ""
	}
	return comments[len(comments)-1]
}




func (p *Properties) GetComments(key string) []string {
	if comments, ok := p.c[key]; ok {
		return comments
	}
	return nil
}




func (p *Properties) SetComment(key, comment string) {
	p.c[key] = []string{comment}
}





func (p *Properties) SetComments(key string, comments []string) {
	if comments == nil {
		delete(p.c, key)
		return
	}
	p.c[key] = comments
}






func (p *Properties) GetBool(key string, def bool) bool {
	v, err := p.getBool(key)
	if err != nil {
		return def
	}
	return v
}




func (p *Properties) MustGetBool(key string) bool {
	v, err := p.getBool(key)
	if err != nil {
		ErrorHandler(err)
	}
	return v
}

func (p *Properties) getBool(key string) (value bool, err error) {
	if v, ok := p.Get(key); ok {
		return boolVal(v), nil
	}
	return false, invalidKeyError(key)
}

func boolVal(v string) bool {
	v = strings.ToLower(v)
	return v == "1" || v == "true" || v == "yes" || v == "on"
}






func (p *Properties) GetDuration(key string, def time.Duration) time.Duration {
	v, err := p.getInt64(key)
	if err != nil {
		return def
	}
	return time.Duration(v)
}




func (p *Properties) MustGetDuration(key string) time.Duration {
	v, err := p.getInt64(key)
	if err != nil {
		ErrorHandler(err)
	}
	return time.Duration(v)
}






func (p *Properties) GetParsedDuration(key string, def time.Duration) time.Duration {
	s, ok := p.Get(key)
	if !ok {
		return def
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return def
	}
	return v
}



func (p *Properties) MustGetParsedDuration(key string) time.Duration {
	s, ok := p.Get(key)
	if !ok {
		ErrorHandler(invalidKeyError(key))
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		ErrorHandler(err)
	}
	return v
}






func (p *Properties) GetFloat64(key string, def float64) float64 {
	v, err := p.getFloat64(key)
	if err != nil {
		return def
	}
	return v
}



func (p *Properties) MustGetFloat64(key string) float64 {
	v, err := p.getFloat64(key)
	if err != nil {
		ErrorHandler(err)
	}
	return v
}

func (p *Properties) getFloat64(key string) (value float64, err error) {
	if v, ok := p.Get(key); ok {
		value, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		return value, nil
	}
	return 0, invalidKeyError(key)
}







func (p *Properties) GetInt(key string, def int) int {
	v, err := p.getInt64(key)
	if err != nil {
		return def
	}
	return intRangeCheck(key, v)
}





func (p *Properties) MustGetInt(key string) int {
	v, err := p.getInt64(key)
	if err != nil {
		ErrorHandler(err)
	}
	return intRangeCheck(key, v)
}






func (p *Properties) GetInt64(key string, def int64) int64 {
	v, err := p.getInt64(key)
	if err != nil {
		return def
	}
	return v
}



func (p *Properties) MustGetInt64(key string) int64 {
	v, err := p.getInt64(key)
	if err != nil {
		ErrorHandler(err)
	}
	return v
}

func (p *Properties) getInt64(key string) (value int64, err error) {
	if v, ok := p.Get(key); ok {
		value, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return value, nil
	}
	return 0, invalidKeyError(key)
}







func (p *Properties) GetUint(key string, def uint) uint {
	v, err := p.getUint64(key)
	if err != nil {
		return def
	}
	return uintRangeCheck(key, v)
}





func (p *Properties) MustGetUint(key string) uint {
	v, err := p.getUint64(key)
	if err != nil {
		ErrorHandler(err)
	}
	return uintRangeCheck(key, v)
}






func (p *Properties) GetUint64(key string, def uint64) uint64 {
	v, err := p.getUint64(key)
	if err != nil {
		return def
	}
	return v
}



func (p *Properties) MustGetUint64(key string) uint64 {
	v, err := p.getUint64(key)
	if err != nil {
		ErrorHandler(err)
	}
	return v
}

func (p *Properties) getUint64(key string) (value uint64, err error) {
	if v, ok := p.Get(key); ok {
		value, err = strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return value, nil
	}
	return 0, invalidKeyError(key)
}





func (p *Properties) GetString(key, def string) string {
	if v, ok := p.Get(key); ok {
		return v
	}
	return def
}



func (p *Properties) MustGetString(key string) string {
	if v, ok := p.Get(key); ok {
		return v
	}
	ErrorHandler(invalidKeyError(key))
	panic("ErrorHandler should exit")
}





func (p *Properties) Filter(pattern string) (*Properties, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return p.FilterRegexp(re), nil
}



func (p *Properties) FilterRegexp(re *regexp.Regexp) *Properties {
	pp := NewProperties()
	for _, k := range p.k {
		if re.MatchString(k) {
			
			
			pp.Set(k, p.m[k])
		}
	}
	return pp
}



func (p *Properties) FilterPrefix(prefix string) *Properties {
	pp := NewProperties()
	for _, k := range p.k {
		if strings.HasPrefix(k, prefix) {
			
			
			pp.Set(k, p.m[k])
		}
	}
	return pp
}



func (p *Properties) FilterStripPrefix(prefix string) *Properties {
	pp := NewProperties()
	n := len(prefix)
	for _, k := range p.k {
		if len(k) > len(prefix) && strings.HasPrefix(k, prefix) {
			
			
			
			pp.Set(k[n:], p.m[k])
		}
	}
	return pp
}


func (p *Properties) Len() int {
	return len(p.m)
}


func (p *Properties) Keys() []string {
	keys := make([]string, len(p.k))
	copy(keys, p.k)
	return keys
}







func (p *Properties) Set(key, value string) (prev string, ok bool, err error) {
	if key == "" {
		return "", false, nil
	}

	
	if p.DisableExpansion {
		prev, ok = p.Get(key)
		p.m[key] = value
		if !ok {
			p.k = append(p.k, key)
		}
		return prev, ok, nil
	}

	
	
	
	
	prev, ok = p.Get(key)
	p.m[key] = value

	
	_, err = p.expand(key, value)
	if err != nil {

		
		if ok {
			p.m[key] = prev
		} else {
			delete(p.m, key)
		}

		return "", false, err
	}

	if !ok {
		p.k = append(p.k, key)
	}

	return prev, ok, nil
}



func (p *Properties) SetValue(key string, value interface{}) error {
	_, _, err := p.Set(key, fmt.Sprintf("%v", value))
	return err
}




func (p *Properties) MustSet(key, value string) (prev string, ok bool) {
	prev, ok, err := p.Set(key, value)
	if err != nil {
		ErrorHandler(err)
	}
	return prev, ok
}


func (p *Properties) String() string {
	var s string
	for _, key := range p.k {
		value, _ := p.Get(key)
		s = fmt.Sprintf("%s%s = %s\n", s, key, value)
	}
	return s
}



func (p *Properties) Write(w io.Writer, enc Encoding) (n int, err error) {
	return p.WriteComment(w, "", enc)
}







func (p *Properties) WriteComment(w io.Writer, prefix string, enc Encoding) (n int, err error) {
	var x int

	for _, key := range p.k {
		value := p.m[key]

		if prefix != "" {
			if comments, ok := p.c[key]; ok {
				
				allEmpty := true
				for _, c := range comments {
					if c != "" {
						allEmpty = false
						break
					}
				}

				if !allEmpty {
					
					if len(comments) > 0 && n > 0 {
						x, err = fmt.Fprintln(w)
						if err != nil {
							return
						}
						n += x
					}

					for _, c := range comments {
						x, err = fmt.Fprintf(w, "%s%s\n", prefix, encode(c, "", enc))
						if err != nil {
							return
						}
						n += x
					}
				}
			}
		}

		x, err = fmt.Fprintf(w, "%s = %s\n", encode(key, " :", enc), encode(value, "", enc))
		if err != nil {
			return
		}
		n += x
	}
	return
}


func (p *Properties) Map() map[string]string {
	m := make(map[string]string)
	for k, v := range p.m {
		m[k] = v
	}
	return m
}


func (p *Properties) FilterFunc(filters ...func(k, v string) bool) *Properties {
	pp := NewProperties()
outer:
	for k, v := range p.m {
		for _, f := range filters {
			if !f(k, v) {
				continue outer
			}
			pp.Set(k, v)
		}
	}
	return pp
}




func (p *Properties) Delete(key string) {
	delete(p.m, key)
	delete(p.c, key)
	newKeys := []string{}
	for _, k := range p.k {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}
	p.k = newKeys
}


func (p *Properties) Merge(other *Properties) {
	for k, v := range other.m {
		p.m[k] = v
	}
	for k, v := range other.c {
		p.c[k] = v
	}

outer:
	for _, otherKey := range other.k {
		for _, key := range p.k {
			if otherKey == key {
				continue outer
			}
		}
		p.k = append(p.k, otherKey)
	}
}





func (p *Properties) check() error {
	for key, value := range p.m {
		if _, err := p.expand(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (p *Properties) expand(key, input string) (string, error) {
	
	if p.Prefix == "" && p.Postfix == "" {
		return input, nil
	}

	return expand(input, []string{key}, p.Prefix, p.Postfix, p.m)
}




func expand(s string, keys []string, prefix, postfix string, values map[string]string) (string, error) {
	if len(keys) > maxExpansionDepth {
		return "", fmt.Errorf("expansion too deep")
	}

	for {
		start := strings.Index(s, prefix)
		if start == -1 {
			return s, nil
		}

		keyStart := start + len(prefix)
		keyLen := strings.Index(s[keyStart:], postfix)
		if keyLen == -1 {
			return "", fmt.Errorf("malformed expression")
		}

		end := keyStart + keyLen + len(postfix) - 1
		key := s[keyStart : keyStart+keyLen]

		

		for _, k := range keys {
			if key == k {
				return "", fmt.Errorf("circular reference")
			}
		}

		val, ok := values[key]
		if !ok {
			val = os.Getenv(key)
		}
		new_val, err := expand(val, append(keys, key), prefix, postfix, values)
		if err != nil {
			return "", err
		}
		s = s[:start] + new_val + s[end+1:]
	}
	return s, nil
}


func encode(s string, special string, enc Encoding) string {
	switch enc {
	case UTF8:
		return encodeUtf8(s, special)
	case ISO_8859_1:
		return encodeIso(s, special)
	default:
		panic(fmt.Sprintf("unsupported encoding %v", enc))
	}
}

func encodeUtf8(s string, special string) string {
	v := ""
	for pos := 0; pos < len(s); {
		r, w := utf8.DecodeRuneInString(s[pos:])
		pos += w
		v += escape(r, special)
	}
	return v
}

func encodeIso(s string, special string) string {
	var r rune
	var w int
	var v string
	for pos := 0; pos < len(s); {
		switch r, w = utf8.DecodeRuneInString(s[pos:]); {
		case r < 1<<8: 
			v += escape(r, special)
		case r < 1<<16: 
			v += fmt.Sprintf("\\u%04x", r)
		default: 
			v += "?"
		}
		pos += w
	}
	return v
}

func escape(r rune, special string) string {
	switch r {
	case '\f':
		return "\\f"
	case '\n':
		return "\\n"
	case '\r':
		return "\\r"
	case '\t':
		return "\\t"
	default:
		if strings.ContainsRune(special, r) {
			return "\\" + string(r)
		}
		return string(r)
	}
}

func invalidKeyError(key string) error {
	return fmt.Errorf("unknown property: %s", key)
}
