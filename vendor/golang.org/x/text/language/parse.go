



package language

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/internal/tag"
)



func isAlpha(b byte) bool {
	return b > '9'
}


func isAlphaNum(s []byte) bool {
	for _, c := range s {
		if !('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9') {
			return false
		}
	}
	return true
}




var errSyntax = errors.New("language: tag is not well-formed")




type ValueError struct {
	v [8]byte
}

func mkErrInvalid(s []byte) error {
	var e ValueError
	copy(e.v[:], s)
	return e
}

func (e ValueError) tag() []byte {
	n := bytes.IndexByte(e.v[:], 0)
	if n == -1 {
		n = 8
	}
	return e.v[:n]
}


func (e ValueError) Error() string {
	return fmt.Sprintf("language: subtag %q is well-formed but unknown", e.tag())
}


func (e ValueError) Subtag() string {
	return string(e.tag())
}


type scanner struct {
	b     []byte
	bytes [max99thPercentileSize]byte
	token []byte
	start int 
	end   int 
	next  int 
	err   error
	done  bool
}

func makeScannerString(s string) scanner {
	scan := scanner{}
	if len(s) <= len(scan.bytes) {
		scan.b = scan.bytes[:copy(scan.bytes[:], s)]
	} else {
		scan.b = []byte(s)
	}
	scan.init()
	return scan
}



func makeScanner(b []byte) scanner {
	scan := scanner{b: b}
	scan.init()
	return scan
}

func (s *scanner) init() {
	for i, c := range s.b {
		if c == '_' {
			s.b[i] = '-'
		}
	}
	s.scan()
}


func (s *scanner) toLower(start, end int) {
	for i := start; i < end; i++ {
		c := s.b[i]
		if 'A' <= c && c <= 'Z' {
			s.b[i] += 'a' - 'A'
		}
	}
}

func (s *scanner) setError(e error) {
	if s.err == nil || (e == errSyntax && s.err != errSyntax) {
		s.err = e
	}
}




func (s *scanner) resizeRange(oldStart, oldEnd, newSize int) {
	s.start = oldStart
	if end := oldStart + newSize; end != oldEnd {
		diff := end - oldEnd
		if end < cap(s.b) {
			b := make([]byte, len(s.b)+diff)
			copy(b, s.b[:oldStart])
			copy(b[end:], s.b[oldEnd:])
			s.b = b
		} else {
			s.b = append(s.b[end:], s.b[oldEnd:]...)
		}
		s.next = end + (s.next - s.end)
		s.end = end
	}
}


func (s *scanner) replace(repl string) {
	s.resizeRange(s.start, s.end, len(repl))
	copy(s.b[s.start:], repl)
}



func (s *scanner) gobble(e error) {
	s.setError(e)
	if s.start == 0 {
		s.b = s.b[:+copy(s.b, s.b[s.next:])]
		s.end = 0
	} else {
		s.b = s.b[:s.start-1+copy(s.b[s.start-1:], s.b[s.end:])]
		s.end = s.start - 1
	}
	s.next = s.start
}


func (s *scanner) deleteRange(start, end int) {
	s.setError(errSyntax)
	s.b = s.b[:start+copy(s.b[start:], s.b[end:])]
	diff := end - start
	s.next -= diff
	s.start -= diff
	s.end -= diff
}





func (s *scanner) scan() (end int) {
	end = s.end
	s.token = nil
	for s.start = s.next; s.next < len(s.b); {
		i := bytes.IndexByte(s.b[s.next:], '-')
		if i == -1 {
			s.end = len(s.b)
			s.next = len(s.b)
			i = s.end - s.start
		} else {
			s.end = s.next + i
			s.next = s.end + 1
		}
		token := s.b[s.start:s.end]
		if i < 1 || i > 8 || !isAlphaNum(token) {
			s.gobble(errSyntax)
			continue
		}
		s.token = token
		return end
	}
	if n := len(s.b); n > 0 && s.b[n-1] == '-' {
		s.setError(errSyntax)
		s.b = s.b[:len(s.b)-1]
	}
	s.done = true
	return end
}



func (s *scanner) acceptMinSize(min int) (end int) {
	end = s.end
	s.scan()
	for ; len(s.token) >= min; s.scan() {
		end = s.end
	}
	return end
}









func Parse(s string) (t Tag, err error) {
	return Default.Parse(s)
}









func (c CanonType) Parse(s string) (t Tag, err error) {
	
	if s == "" {
		return und, errSyntax
	}
	if len(s) <= maxAltTaglen {
		b := [maxAltTaglen]byte{}
		for i, c := range s {
			
			if 'A' <= c && c <= 'Z' {
				c += 'a' - 'A'
			} else if c == '_' {
				c = '-'
			}
			b[i] = byte(c)
		}
		if t, ok := grandfathered(b); ok {
			return t, nil
		}
	}
	scan := makeScannerString(s)
	t, err = parse(&scan, s)
	t, changed := t.canonicalize(c)
	if changed {
		t.remakeString()
	}
	return t, err
}

func parse(scan *scanner, s string) (t Tag, err error) {
	t = und
	var end int
	if n := len(scan.token); n <= 1 {
		scan.toLower(0, len(scan.b))
		if n == 0 || scan.token[0] != 'x' {
			return t, errSyntax
		}
		end = parseExtensions(scan)
	} else if n >= 4 {
		return und, errSyntax
	} else { 
		t, end = parseTag(scan)
		if n := len(scan.token); n == 1 {
			t.pExt = uint16(end)
			end = parseExtensions(scan)
		} else if end < len(scan.b) {
			scan.setError(errSyntax)
			scan.b = scan.b[:end]
		}
	}
	if int(t.pVariant) < len(scan.b) {
		if end < len(s) {
			s = s[:end]
		}
		if len(s) > 0 && tag.Compare(s, scan.b) == 0 {
			t.str = s
		} else {
			t.str = string(scan.b)
		}
	} else {
		t.pVariant, t.pExt = 0, 0
	}
	return t, scan.err
}



func parseTag(scan *scanner) (t Tag, end int) {
	var e error
	
	t.lang, e = getLangID(scan.token)
	scan.setError(e)
	scan.replace(t.lang.String())
	langStart := scan.start
	end = scan.scan()
	for len(scan.token) == 3 && isAlpha(scan.token[0]) {
		
		
		lang, e := getLangID(scan.token)
		if lang != 0 {
			t.lang = lang
			copy(scan.b[langStart:], lang.String())
			scan.b[langStart+3] = '-'
			scan.start = langStart + 4
		}
		scan.gobble(e)
		end = scan.scan()
	}
	if len(scan.token) == 4 && isAlpha(scan.token[0]) {
		t.script, e = getScriptID(script, scan.token)
		if t.script == 0 {
			scan.gobble(e)
		}
		end = scan.scan()
	}
	if n := len(scan.token); n >= 2 && n <= 3 {
		t.region, e = getRegionID(scan.token)
		if t.region == 0 {
			scan.gobble(e)
		} else {
			scan.replace(t.region.String())
		}
		end = scan.scan()
	}
	scan.toLower(scan.start, len(scan.b))
	t.pVariant = byte(end)
	end = parseVariants(scan, end, t)
	t.pExt = uint16(end)
	return t, end
}

var separator = []byte{'-'}



func parseVariants(scan *scanner, end int, t Tag) int {
	start := scan.start
	varIDBuf := [4]uint8{}
	variantBuf := [4][]byte{}
	varID := varIDBuf[:0]
	variant := variantBuf[:0]
	last := -1
	needSort := false
	for ; len(scan.token) >= 4; scan.scan() {
		
		
		v, ok := variantIndex[string(scan.token)]
		if !ok {
			
			
			scan.gobble(mkErrInvalid(scan.token))
			continue
		}
		varID = append(varID, v)
		variant = append(variant, scan.token)
		if !needSort {
			if last < int(v) {
				last = int(v)
			} else {
				needSort = true
				
				
				const maxVariants = 8
				if len(varID) > maxVariants {
					break
				}
			}
		}
		end = scan.end
	}
	if needSort {
		sort.Sort(variantsSort{varID, variant})
		k, l := 0, -1
		for i, v := range varID {
			w := int(v)
			if l == w {
				
				continue
			}
			varID[k] = varID[i]
			variant[k] = variant[i]
			k++
			l = w
		}
		if str := bytes.Join(variant[:k], separator); len(str) == 0 {
			end = start - 1
		} else {
			scan.resizeRange(start, end, len(str))
			copy(scan.b[scan.start:], str)
			end = scan.end
		}
	}
	return end
}

type variantsSort struct {
	i []uint8
	v [][]byte
}

func (s variantsSort) Len() int {
	return len(s.i)
}

func (s variantsSort) Swap(i, j int) {
	s.i[i], s.i[j] = s.i[j], s.i[i]
	s.v[i], s.v[j] = s.v[j], s.v[i]
}

func (s variantsSort) Less(i, j int) bool {
	return s.i[i] < s.i[j]
}

type bytesSort [][]byte

func (b bytesSort) Len() int {
	return len(b)
}

func (b bytesSort) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b bytesSort) Less(i, j int) bool {
	return bytes.Compare(b[i], b[j]) == -1
}




func parseExtensions(scan *scanner) int {
	start := scan.start
	exts := [][]byte{}
	private := []byte{}
	end := scan.end
	for len(scan.token) == 1 {
		extStart := scan.start
		ext := scan.token[0]
		end = parseExtension(scan)
		extension := scan.b[extStart:end]
		if len(extension) < 3 || (ext != 'x' && len(extension) < 4) {
			scan.setError(errSyntax)
			end = extStart
			continue
		} else if start == extStart && (ext == 'x' || scan.start == len(scan.b)) {
			scan.b = scan.b[:end]
			return end
		} else if ext == 'x' {
			private = extension
			break
		}
		exts = append(exts, extension)
	}
	sort.Sort(bytesSort(exts))
	if len(private) > 0 {
		exts = append(exts, private)
	}
	scan.b = scan.b[:start]
	if len(exts) > 0 {
		scan.b = append(scan.b, bytes.Join(exts, separator)...)
	} else if start > 0 {
		
		scan.b = scan.b[:start-1]
	}
	return end
}



func parseExtension(scan *scanner) int {
	start, end := scan.start, scan.end
	switch scan.token[0] {
	case 'u':
		attrStart := end
		scan.scan()
		for last := []byte{}; len(scan.token) > 2; scan.scan() {
			if bytes.Compare(scan.token, last) != -1 {
				
				p := attrStart + 1
				scan.next = p
				attrs := [][]byte{}
				for scan.scan(); len(scan.token) > 2; scan.scan() {
					attrs = append(attrs, scan.token)
					end = scan.end
				}
				sort.Sort(bytesSort(attrs))
				copy(scan.b[p:], bytes.Join(attrs, separator))
				break
			}
			last = scan.token
			end = scan.end
		}
		var last, key []byte
		for attrEnd := end; len(scan.token) == 2; last = key {
			key = scan.token
			keyEnd := scan.end
			end = scan.acceptMinSize(3)
			
			if keyEnd == end || bytes.Compare(key, last) != 1 {
				
				
				p := attrEnd + 1
				scan.next = p
				keys := [][]byte{}
				for scan.scan(); len(scan.token) == 2; {
					keyStart, keyEnd := scan.start, scan.end
					end = scan.acceptMinSize(3)
					if keyEnd != end {
						keys = append(keys, scan.b[keyStart:end])
					} else {
						scan.setError(errSyntax)
						end = keyStart
					}
				}
				sort.Sort(bytesSort(keys))
				reordered := bytes.Join(keys, separator)
				if e := p + len(reordered); e < end {
					scan.deleteRange(e, end)
					end = e
				}
				copy(scan.b[p:], bytes.Join(keys, separator))
				break
			}
		}
	case 't':
		scan.scan()
		if n := len(scan.token); n >= 2 && n <= 3 && isAlpha(scan.token[1]) {
			_, end = parseTag(scan)
			scan.toLower(start, end)
		}
		for len(scan.token) == 2 && !isAlpha(scan.token[1]) {
			end = scan.acceptMinSize(3)
		}
	case 'x':
		end = scan.acceptMinSize(1)
	default:
		end = scan.acceptMinSize(2)
	}
	return end
}










func Compose(part ...interface{}) (t Tag, err error) {
	return Default.Compose(part...)
}










func (c CanonType) Compose(part ...interface{}) (t Tag, err error) {
	var b builder
	if err = b.update(part...); err != nil {
		return und, err
	}
	t, _ = b.tag.canonicalize(c)

	if len(b.ext) > 0 || len(b.variant) > 0 {
		sort.Sort(sortVariant(b.variant))
		sort.Strings(b.ext)
		if b.private != "" {
			b.ext = append(b.ext, b.private)
		}
		n := maxCoreSize + tokenLen(b.variant...) + tokenLen(b.ext...)
		buf := make([]byte, n)
		p := t.genCoreBytes(buf)
		t.pVariant = byte(p)
		p += appendTokens(buf[p:], b.variant...)
		t.pExt = uint16(p)
		p += appendTokens(buf[p:], b.ext...)
		t.str = string(buf[:p])
	} else if b.private != "" {
		t.str = b.private
		t.remakeString()
	}
	return
}

type builder struct {
	tag Tag

	private string 
	ext     []string
	variant []string

	err error
}

func (b *builder) addExt(e string) {
	if e == "" {
	} else if e[0] == 'x' {
		b.private = e
	} else {
		b.ext = append(b.ext, e)
	}
}

var errInvalidArgument = errors.New("invalid Extension or Variant")

func (b *builder) update(part ...interface{}) (err error) {
	replace := func(l *[]string, s string, eq func(a, b string) bool) bool {
		if s == "" {
			b.err = errInvalidArgument
			return true
		}
		for i, v := range *l {
			if eq(v, s) {
				(*l)[i] = s
				return true
			}
		}
		return false
	}
	for _, x := range part {
		switch v := x.(type) {
		case Tag:
			b.tag.lang = v.lang
			b.tag.region = v.region
			b.tag.script = v.script
			if v.str != "" {
				b.variant = nil
				for x, s := "", v.str[v.pVariant:v.pExt]; s != ""; {
					x, s = nextToken(s)
					b.variant = append(b.variant, x)
				}
				b.ext, b.private = nil, ""
				for i, e := int(v.pExt), ""; i < len(v.str); {
					i, e = getExtension(v.str, i)
					b.addExt(e)
				}
			}
		case Base:
			b.tag.lang = v.langID
		case Script:
			b.tag.script = v.scriptID
		case Region:
			b.tag.region = v.regionID
		case Variant:
			if !replace(&b.variant, v.variant, func(a, b string) bool { return a == b }) {
				b.variant = append(b.variant, v.variant)
			}
		case Extension:
			if !replace(&b.ext, v.s, func(a, b string) bool { return a[0] == b[0] }) {
				b.addExt(v.s)
			}
		case []Variant:
			b.variant = nil
			for _, x := range v {
				b.update(x)
			}
		case []Extension:
			b.ext, b.private = nil, ""
			for _, e := range v {
				b.update(e)
			}
		
		case error:
			err = v
		}
	}
	return
}

func tokenLen(token ...string) (n int) {
	for _, t := range token {
		n += len(t) + 1
	}
	return
}

func appendTokens(b []byte, token ...string) int {
	p := 0
	for _, t := range token {
		b[p] = '-'
		copy(b[p+1:], t)
		p += 1 + len(t)
	}
	return p
}

type sortVariant []string

func (s sortVariant) Len() int {
	return len(s)
}

func (s sortVariant) Swap(i, j int) {
	s[j], s[i] = s[i], s[j]
}

func (s sortVariant) Less(i, j int) bool {
	return variantIndex[s[i]] < variantIndex[s[j]]
}

func findExt(list []string, x byte) int {
	for i, e := range list {
		if e[0] == x {
			return i
		}
	}
	return -1
}


func getExtension(s string, p int) (end int, ext string) {
	if s[p] == '-' {
		p++
	}
	if s[p] == 'x' {
		return len(s), s[p:]
	}
	end = nextExtension(s, p)
	return end, s[p:end]
}





func nextExtension(s string, p int) int {
	for n := len(s) - 3; p < n; {
		if s[p] == '-' {
			if s[p+2] == '-' {
				return p
			}
			p += 3
		} else {
			p++
		}
	}
	return len(s)
}

var errInvalidWeight = errors.New("ParseAcceptLanguage: invalid weight")








func ParseAcceptLanguage(s string) (tag []Tag, q []float32, err error) {
	var entry string
	for s != "" {
		if entry, s = split(s, ','); entry == "" {
			continue
		}

		entry, weight := split(entry, ';')

		
		t, err := Parse(entry)
		if err != nil {
			id, ok := acceptFallback[entry]
			if !ok {
				return nil, nil, err
			}
			t = Tag{lang: id}
		}

		
		w := 1.0
		if weight != "" {
			weight = consume(weight, 'q')
			weight = consume(weight, '=')
			
			
			if w, err = strconv.ParseFloat(weight, 32); err != nil {
				return nil, nil, errInvalidWeight
			}
			
			if w <= 0 {
				continue
			}
		}

		tag = append(tag, t)
		q = append(q, float32(w))
	}
	sortStable(&tagSort{tag, q})
	return tag, q, nil
}



func consume(s string, c byte) string {
	if s == "" || s[0] != c {
		return ""
	}
	return strings.TrimSpace(s[1:])
}

func split(s string, c byte) (head, tail string) {
	if i := strings.IndexByte(s, c); i >= 0 {
		return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:])
	}
	return strings.TrimSpace(s), ""
}



var acceptFallback = map[string]langID{
	"english": _en,
	"deutsch": _de,
	"italian": _it,
	"french":  _fr,
	"*":       _mul, 
}

type tagSort struct {
	tag []Tag
	q   []float32
}

func (s *tagSort) Len() int {
	return len(s.q)
}

func (s *tagSort) Less(i, j int) bool {
	return s.q[i] > s.q[j]
}

func (s *tagSort) Swap(i, j int) {
	s.tag[i], s.tag[j] = s.tag[j], s.tag[i]
	s.q[i], s.q[j] = s.q[j], s.q[i]
}
