



package html

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html/atom"
)


type TokenType uint32

const (
	
	ErrorToken TokenType = iota
	
	TextToken
	
	StartTagToken
	
	EndTagToken
	
	SelfClosingTagToken
	
	CommentToken
	
	DoctypeToken
)


var ErrBufferExceeded = errors.New("max buffer exceeded")


func (t TokenType) String() string {
	switch t {
	case ErrorToken:
		return "Error"
	case TextToken:
		return "Text"
	case StartTagToken:
		return "StartTag"
	case EndTagToken:
		return "EndTag"
	case SelfClosingTagToken:
		return "SelfClosingTag"
	case CommentToken:
		return "Comment"
	case DoctypeToken:
		return "Doctype"
	}
	return "Invalid(" + strconv.Itoa(int(t)) + ")"
}







type Attribute struct {
	Namespace, Key, Val string
}






type Token struct {
	Type     TokenType
	DataAtom atom.Atom
	Data     string
	Attr     []Attribute
}


func (t Token) tagString() string {
	if len(t.Attr) == 0 {
		return t.Data
	}
	buf := bytes.NewBufferString(t.Data)
	for _, a := range t.Attr {
		buf.WriteByte(' ')
		buf.WriteString(a.Key)
		buf.WriteString(`="`)
		escape(buf, a.Val)
		buf.WriteByte('"')
	}
	return buf.String()
}


func (t Token) String() string {
	switch t.Type {
	case ErrorToken:
		return ""
	case TextToken:
		return EscapeString(t.Data)
	case StartTagToken:
		return "<" + t.tagString() + ">"
	case EndTagToken:
		return "</" + t.tagString() + ">"
	case SelfClosingTagToken:
		return "<" + t.tagString() + "/>"
	case CommentToken:
		return "<!--" + t.Data + "-->"
	case DoctypeToken:
		return "<!DOCTYPE " + t.Data + ">"
	}
	return "Invalid(" + strconv.Itoa(int(t.Type)) + ")"
}



type span struct {
	start, end int
}


type Tokenizer struct {
	
	r io.Reader
	
	tt TokenType
	
	
	
	
	
	
	
	err error
	
	
	
	
	readErr error
	
	
	raw span
	buf []byte
	
	maxBuf int
	
	
	data span
	
	
	
	pendingAttr   [2]span
	attr          [][2]span
	nAttrReturned int
	
	
	
	
	rawTag string
	
	textIsRaw bool
	
	
	convertNUL bool
	
	allowCDATA bool
}















func (z *Tokenizer) AllowCDATA(allowCDATA bool) {
	z.allowCDATA = allowCDATA
}

























func (z *Tokenizer) NextIsNotRawText() {
	z.rawTag = ""
}



func (z *Tokenizer) Err() error {
	if z.tt != ErrorToken {
		return nil
	}
	return z.err
}






func (z *Tokenizer) readByte() byte {
	if z.raw.end >= len(z.buf) {
		
		
		if z.readErr != nil {
			z.err = z.readErr
			return 0
		}
		
		
		
		c := cap(z.buf)
		d := z.raw.end - z.raw.start
		var buf1 []byte
		if 2*d > c {
			buf1 = make([]byte, d, 2*c)
		} else {
			buf1 = z.buf[:d]
		}
		copy(buf1, z.buf[z.raw.start:z.raw.end])
		if x := z.raw.start; x != 0 {
			
			z.data.start -= x
			z.data.end -= x
			z.pendingAttr[0].start -= x
			z.pendingAttr[0].end -= x
			z.pendingAttr[1].start -= x
			z.pendingAttr[1].end -= x
			for i := range z.attr {
				z.attr[i][0].start -= x
				z.attr[i][0].end -= x
				z.attr[i][1].start -= x
				z.attr[i][1].end -= x
			}
		}
		z.raw.start, z.raw.end, z.buf = 0, d, buf1[:d]
		
		
		var n int
		n, z.readErr = readAtLeastOneByte(z.r, buf1[d:cap(buf1)])
		if n == 0 {
			z.err = z.readErr
			return 0
		}
		z.buf = buf1[:d+n]
	}
	x := z.buf[z.raw.end]
	z.raw.end++
	if z.maxBuf > 0 && z.raw.end-z.raw.start >= z.maxBuf {
		z.err = ErrBufferExceeded
		return 0
	}
	return x
}


func (z *Tokenizer) Buffered() []byte {
	return z.buf[z.raw.end:]
}




func readAtLeastOneByte(r io.Reader, b []byte) (int, error) {
	for i := 0; i < 100; i++ {
		n, err := r.Read(b)
		if n != 0 || err != nil {
			return n, err
		}
	}
	return 0, io.ErrNoProgress
}


func (z *Tokenizer) skipWhiteSpace() {
	if z.err != nil {
		return
	}
	for {
		c := z.readByte()
		if z.err != nil {
			return
		}
		switch c {
		case ' ', '\n', '\r', '\t', '\f':
			
		default:
			z.raw.end--
			return
		}
	}
}



func (z *Tokenizer) readRawOrRCDATA() {
	if z.rawTag == "script" {
		z.readScript()
		z.textIsRaw = true
		z.rawTag = ""
		return
	}
loop:
	for {
		c := z.readByte()
		if z.err != nil {
			break loop
		}
		if c != '<' {
			continue loop
		}
		c = z.readByte()
		if z.err != nil {
			break loop
		}
		if c != '/' {
			continue loop
		}
		if z.readRawEndTag() || z.err != nil {
			break loop
		}
	}
	z.data.end = z.raw.end
	
	z.textIsRaw = z.rawTag != "textarea" && z.rawTag != "title"
	z.rawTag = ""
}





func (z *Tokenizer) readRawEndTag() bool {
	for i := 0; i < len(z.rawTag); i++ {
		c := z.readByte()
		if z.err != nil {
			return false
		}
		if c != z.rawTag[i] && c != z.rawTag[i]-('a'-'A') {
			z.raw.end--
			return false
		}
	}
	c := z.readByte()
	if z.err != nil {
		return false
	}
	switch c {
	case ' ', '\n', '\r', '\t', '\f', '/', '>':
		
		z.raw.end -= 3 + len(z.rawTag)
		return true
	}
	z.raw.end--
	return false
}



func (z *Tokenizer) readScript() {
	defer func() {
		z.data.end = z.raw.end
	}()
	var c byte

scriptData:
	c = z.readByte()
	if z.err != nil {
		return
	}
	if c == '<' {
		goto scriptDataLessThanSign
	}
	goto scriptData

scriptDataLessThanSign:
	c = z.readByte()
	if z.err != nil {
		return
	}
	switch c {
	case '/':
		goto scriptDataEndTagOpen
	case '!':
		goto scriptDataEscapeStart
	}
	z.raw.end--
	goto scriptData

scriptDataEndTagOpen:
	if z.readRawEndTag() || z.err != nil {
		return
	}
	goto scriptData

scriptDataEscapeStart:
	c = z.readByte()
	if z.err != nil {
		return
	}
	if c == '-' {
		goto scriptDataEscapeStartDash
	}
	z.raw.end--
	goto scriptData

scriptDataEscapeStartDash:
	c = z.readByte()
	if z.err != nil {
		return
	}
	if c == '-' {
		goto scriptDataEscapedDashDash
	}
	z.raw.end--
	goto scriptData

scriptDataEscaped:
	c = z.readByte()
	if z.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataEscapedDash
	case '<':
		goto scriptDataEscapedLessThanSign
	}
	goto scriptDataEscaped

scriptDataEscapedDash:
	c = z.readByte()
	if z.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataEscapedDashDash
	case '<':
		goto scriptDataEscapedLessThanSign
	}
	goto scriptDataEscaped

scriptDataEscapedDashDash:
	c = z.readByte()
	if z.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataEscapedDashDash
	case '<':
		goto scriptDataEscapedLessThanSign
	case '>':
		goto scriptData
	}
	goto scriptDataEscaped

scriptDataEscapedLessThanSign:
	c = z.readByte()
	if z.err != nil {
		return
	}
	if c == '/' {
		goto scriptDataEscapedEndTagOpen
	}
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
		goto scriptDataDoubleEscapeStart
	}
	z.raw.end--
	goto scriptData

scriptDataEscapedEndTagOpen:
	if z.readRawEndTag() || z.err != nil {
		return
	}
	goto scriptDataEscaped

scriptDataDoubleEscapeStart:
	z.raw.end--
	for i := 0; i < len("script"); i++ {
		c = z.readByte()
		if z.err != nil {
			return
		}
		if c != "script"[i] && c != "SCRIPT"[i] {
			z.raw.end--
			goto scriptDataEscaped
		}
	}
	c = z.readByte()
	if z.err != nil {
		return
	}
	switch c {
	case ' ', '\n', '\r', '\t', '\f', '/', '>':
		goto scriptDataDoubleEscaped
	}
	z.raw.end--
	goto scriptDataEscaped

scriptDataDoubleEscaped:
	c = z.readByte()
	if z.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataDoubleEscapedDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	}
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedDash:
	c = z.readByte()
	if z.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataDoubleEscapedDashDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	}
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedDashDash:
	c = z.readByte()
	if z.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataDoubleEscapedDashDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	case '>':
		goto scriptData
	}
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedLessThanSign:
	c = z.readByte()
	if z.err != nil {
		return
	}
	if c == '/' {
		goto scriptDataDoubleEscapeEnd
	}
	z.raw.end--
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapeEnd:
	if z.readRawEndTag() {
		z.raw.end += len("</script>")
		goto scriptDataEscaped
	}
	if z.err != nil {
		return
	}
	goto scriptDataDoubleEscaped
}



func (z *Tokenizer) readComment() {
	z.data.start = z.raw.end
	defer func() {
		if z.data.end < z.data.start {
			
			z.data.end = z.data.start
		}
	}()
	for dashCount := 2; ; {
		c := z.readByte()
		if z.err != nil {
			
			if dashCount > 2 {
				dashCount = 2
			}
			z.data.end = z.raw.end - dashCount
			return
		}
		switch c {
		case '-':
			dashCount++
			continue
		case '>':
			if dashCount >= 2 {
				z.data.end = z.raw.end - len("-->")
				return
			}
		case '!':
			if dashCount >= 2 {
				c = z.readByte()
				if z.err != nil {
					z.data.end = z.raw.end
					return
				}
				if c == '>' {
					z.data.end = z.raw.end - len("--!>")
					return
				}
			}
		}
		dashCount = 0
	}
}


func (z *Tokenizer) readUntilCloseAngle() {
	z.data.start = z.raw.end
	for {
		c := z.readByte()
		if z.err != nil {
			z.data.end = z.raw.end
			return
		}
		if c == '>' {
			z.data.end = z.raw.end - len(">")
			return
		}
	}
}




func (z *Tokenizer) readMarkupDeclaration() TokenType {
	z.data.start = z.raw.end
	var c [2]byte
	for i := 0; i < 2; i++ {
		c[i] = z.readByte()
		if z.err != nil {
			z.data.end = z.raw.end
			return CommentToken
		}
	}
	if c[0] == '-' && c[1] == '-' {
		z.readComment()
		return CommentToken
	}
	z.raw.end -= 2
	if z.readDoctype() {
		return DoctypeToken
	}
	if z.allowCDATA && z.readCDATA() {
		z.convertNUL = true
		return TextToken
	}
	
	z.readUntilCloseAngle()
	return CommentToken
}



func (z *Tokenizer) readDoctype() bool {
	const s = "DOCTYPE"
	for i := 0; i < len(s); i++ {
		c := z.readByte()
		if z.err != nil {
			z.data.end = z.raw.end
			return false
		}
		if c != s[i] && c != s[i]+('a'-'A') {
			
			z.raw.end = z.data.start
			return false
		}
	}
	if z.skipWhiteSpace(); z.err != nil {
		z.data.start = z.raw.end
		z.data.end = z.raw.end
		return true
	}
	z.readUntilCloseAngle()
	return true
}



func (z *Tokenizer) readCDATA() bool {
	const s = "[CDATA["
	for i := 0; i < len(s); i++ {
		c := z.readByte()
		if z.err != nil {
			z.data.end = z.raw.end
			return false
		}
		if c != s[i] {
			
			z.raw.end = z.data.start
			return false
		}
	}
	z.data.start = z.raw.end
	brackets := 0
	for {
		c := z.readByte()
		if z.err != nil {
			z.data.end = z.raw.end
			return true
		}
		switch c {
		case ']':
			brackets++
		case '>':
			if brackets >= 2 {
				z.data.end = z.raw.end - len("]]>")
				return true
			}
			brackets = 0
		default:
			brackets = 0
		}
	}
}



func (z *Tokenizer) startTagIn(ss ...string) bool {
loop:
	for _, s := range ss {
		if z.data.end-z.data.start != len(s) {
			continue loop
		}
		for i := 0; i < len(s); i++ {
			c := z.buf[z.data.start+i]
			if 'A' <= c && c <= 'Z' {
				c += 'a' - 'A'
			}
			if c != s[i] {
				continue loop
			}
		}
		return true
	}
	return false
}



func (z *Tokenizer) readStartTag() TokenType {
	z.readTag(true)
	if z.err != nil {
		return ErrorToken
	}
	
	c, raw := z.buf[z.data.start], false
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	switch c {
	case 'i':
		raw = z.startTagIn("iframe")
	case 'n':
		raw = z.startTagIn("noembed", "noframes", "noscript")
	case 'p':
		raw = z.startTagIn("plaintext")
	case 's':
		raw = z.startTagIn("script", "style")
	case 't':
		raw = z.startTagIn("textarea", "title")
	case 'x':
		raw = z.startTagIn("xmp")
	}
	if raw {
		z.rawTag = strings.ToLower(string(z.buf[z.data.start:z.data.end]))
	}
	
	if z.err == nil && z.buf[z.raw.end-2] == '/' {
		return SelfClosingTagToken
	}
	return StartTagToken
}





func (z *Tokenizer) readTag(saveAttr bool) {
	z.attr = z.attr[:0]
	z.nAttrReturned = 0
	
	z.readTagName()
	if z.skipWhiteSpace(); z.err != nil {
		return
	}
	for {
		c := z.readByte()
		if z.err != nil || c == '>' {
			break
		}
		z.raw.end--
		z.readTagAttrKey()
		z.readTagAttrVal()
		
		if saveAttr && z.pendingAttr[0].start != z.pendingAttr[0].end {
			z.attr = append(z.attr, z.pendingAttr)
		}
		if z.skipWhiteSpace(); z.err != nil {
			break
		}
	}
}




func (z *Tokenizer) readTagName() {
	z.data.start = z.raw.end - 1
	for {
		c := z.readByte()
		if z.err != nil {
			z.data.end = z.raw.end
			return
		}
		switch c {
		case ' ', '\n', '\r', '\t', '\f':
			z.data.end = z.raw.end - 1
			return
		case '/', '>':
			z.raw.end--
			z.data.end = z.raw.end
			return
		}
	}
}



func (z *Tokenizer) readTagAttrKey() {
	z.pendingAttr[0].start = z.raw.end
	for {
		c := z.readByte()
		if z.err != nil {
			z.pendingAttr[0].end = z.raw.end
			return
		}
		switch c {
		case ' ', '\n', '\r', '\t', '\f', '/':
			z.pendingAttr[0].end = z.raw.end - 1
			return
		case '=', '>':
			z.raw.end--
			z.pendingAttr[0].end = z.raw.end
			return
		}
	}
}


func (z *Tokenizer) readTagAttrVal() {
	z.pendingAttr[1].start = z.raw.end
	z.pendingAttr[1].end = z.raw.end
	if z.skipWhiteSpace(); z.err != nil {
		return
	}
	c := z.readByte()
	if z.err != nil {
		return
	}
	if c != '=' {
		z.raw.end--
		return
	}
	if z.skipWhiteSpace(); z.err != nil {
		return
	}
	quote := z.readByte()
	if z.err != nil {
		return
	}
	switch quote {
	case '>':
		z.raw.end--
		return

	case '\'', '"':
		z.pendingAttr[1].start = z.raw.end
		for {
			c := z.readByte()
			if z.err != nil {
				z.pendingAttr[1].end = z.raw.end
				return
			}
			if c == quote {
				z.pendingAttr[1].end = z.raw.end - 1
				return
			}
		}

	default:
		z.pendingAttr[1].start = z.raw.end - 1
		for {
			c := z.readByte()
			if z.err != nil {
				z.pendingAttr[1].end = z.raw.end
				return
			}
			switch c {
			case ' ', '\n', '\r', '\t', '\f':
				z.pendingAttr[1].end = z.raw.end - 1
				return
			case '>':
				z.raw.end--
				z.pendingAttr[1].end = z.raw.end
				return
			}
		}
	}
}


func (z *Tokenizer) Next() TokenType {
	z.raw.start = z.raw.end
	z.data.start = z.raw.end
	z.data.end = z.raw.end
	if z.err != nil {
		z.tt = ErrorToken
		return z.tt
	}
	if z.rawTag != "" {
		if z.rawTag == "plaintext" {
			
			for z.err == nil {
				z.readByte()
			}
			z.data.end = z.raw.end
			z.textIsRaw = true
		} else {
			z.readRawOrRCDATA()
		}
		if z.data.end > z.data.start {
			z.tt = TextToken
			z.convertNUL = true
			return z.tt
		}
	}
	z.textIsRaw = false
	z.convertNUL = false

loop:
	for {
		c := z.readByte()
		if z.err != nil {
			break loop
		}
		if c != '<' {
			continue loop
		}

		
		
		c = z.readByte()
		if z.err != nil {
			break loop
		}
		var tokenType TokenType
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
			tokenType = StartTagToken
		case c == '/':
			tokenType = EndTagToken
		case c == '!' || c == '?':
			
			
			tokenType = CommentToken
		default:
			
			z.raw.end--
			continue
		}

		
		
		
		if x := z.raw.end - len("<a"); z.raw.start < x {
			z.raw.end = x
			z.data.end = x
			z.tt = TextToken
			return z.tt
		}
		switch tokenType {
		case StartTagToken:
			z.tt = z.readStartTag()
			return z.tt
		case EndTagToken:
			c = z.readByte()
			if z.err != nil {
				break loop
			}
			if c == '>' {
				
				
				
				z.tt = CommentToken
				return z.tt
			}
			if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
				z.readTag(false)
				if z.err != nil {
					z.tt = ErrorToken
				} else {
					z.tt = EndTagToken
				}
				return z.tt
			}
			z.raw.end--
			z.readUntilCloseAngle()
			z.tt = CommentToken
			return z.tt
		case CommentToken:
			if c == '!' {
				z.tt = z.readMarkupDeclaration()
				return z.tt
			}
			z.raw.end--
			z.readUntilCloseAngle()
			z.tt = CommentToken
			return z.tt
		}
	}
	if z.raw.start < z.raw.end {
		z.data.end = z.raw.end
		z.tt = TextToken
		return z.tt
	}
	z.tt = ErrorToken
	return z.tt
}



func (z *Tokenizer) Raw() []byte {
	return z.buf[z.raw.start:z.raw.end]
}



func convertNewlines(s []byte) []byte {
	for i, c := range s {
		if c != '\r' {
			continue
		}

		src := i + 1
		if src >= len(s) || s[src] != '\n' {
			s[i] = '\n'
			continue
		}

		dst := i
		for src < len(s) {
			if s[src] == '\r' {
				if src+1 < len(s) && s[src+1] == '\n' {
					src++
				}
				s[dst] = '\n'
			} else {
				s[dst] = s[src]
			}
			src++
			dst++
		}
		return s[:dst]
	}
	return s
}

var (
	nul         = []byte("\x00")
	replacement = []byte("\ufffd")
)



func (z *Tokenizer) Text() []byte {
	switch z.tt {
	case TextToken, CommentToken, DoctypeToken:
		s := z.buf[z.data.start:z.data.end]
		z.data.start = z.raw.end
		z.data.end = z.raw.end
		s = convertNewlines(s)
		if (z.convertNUL || z.tt == CommentToken) && bytes.Contains(s, nul) {
			s = bytes.Replace(s, nul, replacement, -1)
		}
		if !z.textIsRaw {
			s = unescape(s, false)
		}
		return s
	}
	return nil
}




func (z *Tokenizer) TagName() (name []byte, hasAttr bool) {
	if z.data.start < z.data.end {
		switch z.tt {
		case StartTagToken, EndTagToken, SelfClosingTagToken:
			s := z.buf[z.data.start:z.data.end]
			z.data.start = z.raw.end
			z.data.end = z.raw.end
			return lower(s), z.nAttrReturned < len(z.attr)
		}
	}
	return nil, false
}




func (z *Tokenizer) TagAttr() (key, val []byte, moreAttr bool) {
	if z.nAttrReturned < len(z.attr) {
		switch z.tt {
		case StartTagToken, SelfClosingTagToken:
			x := z.attr[z.nAttrReturned]
			z.nAttrReturned++
			key = z.buf[x[0].start:x[0].end]
			val = z.buf[x[1].start:x[1].end]
			return lower(key), unescape(convertNewlines(val), true), z.nAttrReturned < len(z.attr)
		}
	}
	return nil, nil, false
}



func (z *Tokenizer) Token() Token {
	t := Token{Type: z.tt}
	switch z.tt {
	case TextToken, CommentToken, DoctypeToken:
		t.Data = string(z.Text())
	case StartTagToken, SelfClosingTagToken, EndTagToken:
		name, moreAttr := z.TagName()
		for moreAttr {
			var key, val []byte
			key, val, moreAttr = z.TagAttr()
			t.Attr = append(t.Attr, Attribute{"", atom.String(key), string(val)})
		}
		if a := atom.Lookup(name); a != 0 {
			t.DataAtom, t.Data = a, a.String()
		} else {
			t.DataAtom, t.Data = 0, string(name)
		}
	}
	return t
}



func (z *Tokenizer) SetMaxBuf(n int) {
	z.maxBuf = n
}



func NewTokenizer(r io.Reader) *Tokenizer {
	return NewTokenizerFragment(r, "")
}









func NewTokenizerFragment(r io.Reader, contextTag string) *Tokenizer {
	z := &Tokenizer{
		r:   r,
		buf: make([]byte, 0, 4096),
	}
	if contextTag != "" {
		switch s := strings.ToLower(contextTag); s {
		case "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "title", "textarea", "xmp":
			z.rawTag = s
		}
	}
	return z
}
