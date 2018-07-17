







package charset 

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)






func Lookup(label string) (e encoding.Encoding, name string) {
	e, err := htmlindex.Get(label)
	if err != nil {
		return nil, ""
	}
	name, _ = htmlindex.Name(e)
	return &htmlEncoding{e}, name
}

type htmlEncoding struct{ encoding.Encoding }

func (h *htmlEncoding) NewEncoder() *encoding.Encoder {
	
	
	return encoding.HTMLEscapeUnsupported(h.Encoding.NewEncoder())
}





func DetermineEncoding(content []byte, contentType string) (e encoding.Encoding, name string, certain bool) {
	if len(content) > 1024 {
		content = content[:1024]
	}

	for _, b := range boms {
		if bytes.HasPrefix(content, b.bom) {
			e, name = Lookup(b.enc)
			return e, name, true
		}
	}

	if _, params, err := mime.ParseMediaType(contentType); err == nil {
		if cs, ok := params["charset"]; ok {
			if e, name = Lookup(cs); e != nil {
				return e, name, true
			}
		}
	}

	if len(content) > 0 {
		e, name = prescan(content)
		if e != nil {
			return e, name, false
		}
	}

	
	
	for i := len(content) - 1; i >= 0 && i > len(content)-4; i-- {
		b := content[i]
		if b < 0x80 {
			break
		}
		if utf8.RuneStart(b) {
			content = content[:i]
			break
		}
	}
	hasHighBit := false
	for _, c := range content {
		if c >= 0x80 {
			hasHighBit = true
			break
		}
	}
	if hasHighBit && utf8.Valid(content) {
		return encoding.Nop, "utf-8", false
	}

	
	return charmap.Windows1252, "windows-1252", false
}



func NewReader(r io.Reader, contentType string) (io.Reader, error) {
	preview := make([]byte, 1024)
	n, err := io.ReadFull(r, preview)
	switch {
	case err == io.ErrUnexpectedEOF:
		preview = preview[:n]
		r = bytes.NewReader(preview)
	case err != nil:
		return nil, err
	default:
		r = io.MultiReader(bytes.NewReader(preview), r)
	}

	if e, _, _ := DetermineEncoding(preview, contentType); e != encoding.Nop {
		r = transform.NewReader(r, e.NewDecoder())
	}
	return r, nil
}





func NewReaderLabel(label string, input io.Reader) (io.Reader, error) {
	e, _ := Lookup(label)
	if e == nil {
		return nil, fmt.Errorf("unsupported charset: %q", label)
	}
	return transform.NewReader(input, e.NewDecoder()), nil
}

func prescan(content []byte) (e encoding.Encoding, name string) {
	z := html.NewTokenizer(bytes.NewReader(content))
	for {
		switch z.Next() {
		case html.ErrorToken:
			return nil, ""

		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, hasAttr := z.TagName()
			if !bytes.Equal(tagName, []byte("meta")) {
				continue
			}
			attrList := make(map[string]bool)
			gotPragma := false

			const (
				dontKnow = iota
				doNeedPragma
				doNotNeedPragma
			)
			needPragma := dontKnow

			name = ""
			e = nil
			for hasAttr {
				var key, val []byte
				key, val, hasAttr = z.TagAttr()
				ks := string(key)
				if attrList[ks] {
					continue
				}
				attrList[ks] = true
				for i, c := range val {
					if 'A' <= c && c <= 'Z' {
						val[i] = c + 0x20
					}
				}

				switch ks {
				case "http-equiv":
					if bytes.Equal(val, []byte("content-type")) {
						gotPragma = true
					}

				case "content":
					if e == nil {
						name = fromMetaElement(string(val))
						if name != "" {
							e, name = Lookup(name)
							if e != nil {
								needPragma = doNeedPragma
							}
						}
					}

				case "charset":
					e, name = Lookup(string(val))
					needPragma = doNotNeedPragma
				}
			}

			if needPragma == dontKnow || needPragma == doNeedPragma && !gotPragma {
				continue
			}

			if strings.HasPrefix(name, "utf-16") {
				name = "utf-8"
				e = encoding.Nop
			}

			if e != nil {
				return e, name
			}
		}
	}
}

func fromMetaElement(s string) string {
	for s != "" {
		csLoc := strings.Index(s, "charset")
		if csLoc == -1 {
			return ""
		}
		s = s[csLoc+len("charset"):]
		s = strings.TrimLeft(s, " \t\n\f\r")
		if !strings.HasPrefix(s, "=") {
			continue
		}
		s = s[1:]
		s = strings.TrimLeft(s, " \t\n\f\r")
		if s == "" {
			return ""
		}
		if q := s[0]; q == '"' || q == '\'' {
			s = s[1:]
			closeQuote := strings.IndexRune(s, rune(q))
			if closeQuote == -1 {
				return ""
			}
			return s[:closeQuote]
		}

		end := strings.IndexAny(s, "; \t\n\f\r")
		if end == -1 {
			end = len(s)
		}
		return s[:end]
	}
	return ""
}

var boms = []struct {
	bom []byte
	enc string
}{
	{[]byte{0xfe, 0xff}, "utf-16be"},
	{[]byte{0xff, 0xfe}, "utf-16le"},
	{[]byte{0xef, 0xbb, 0xbf}, "utf-8"},
}
