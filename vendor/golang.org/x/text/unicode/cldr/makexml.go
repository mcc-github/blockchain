






package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/internal/gen"
)

var outputFile = flag.String("output", "xml.go", "output file name")

func main() {
	flag.Parse()

	r := gen.OpenCLDRCoreZip()
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal("Could not read zip file")
	}
	r.Close()
	z, err := zip.NewReader(bytes.NewReader(buffer), int64(len(buffer)))
	if err != nil {
		log.Fatalf("Could not read zip archive: %v", err)
	}

	var buf bytes.Buffer

	version := gen.CLDRVersion()

	for _, dtd := range files {
		for _, f := range z.File {
			if strings.HasSuffix(f.Name, dtd.file+".dtd") {
				r, err := f.Open()
				failOnError(err)

				b := makeBuilder(&buf, dtd)
				b.parseDTD(r)
				b.resolve(b.index[dtd.top[0]])
				b.write()
				if b.version != "" && version != b.version {
					println(f.Name)
					log.Fatalf("main: inconsistent versions: found %s; want %s", b.version, version)
				}
				break
			}
		}
	}
	fmt.Fprintln(&buf, "
	fmt.Fprintf(&buf, "const Version = %q\n", version)

	gen.WriteGoFile(*outputFile, "cldr", buf.Bytes())
}

func failOnError(err error) {
	if err != nil {
		log.New(os.Stderr, "", log.Lshortfile).Output(2, err.Error())
		os.Exit(1)
	}
}


type dtd struct {
	file string   
	root string   
	top  []string 

	skipElem    []string 
	skipAttr    []string 
	predefined  []string 
	forceRepeat []string 
}

var files = []dtd{
	{
		file: "ldmlBCP47",
		root: "LDMLBCP47",
		top:  []string{"ldmlBCP47"},
		skipElem: []string{
			"cldrVersion", 
		},
	},
	{
		file: "ldmlSupplemental",
		root: "SupplementalData",
		top:  []string{"supplementalData"},
		skipElem: []string{
			"cldrVersion", 
		},
		forceRepeat: []string{
			"plurals", 
		},
	},
	{
		file: "ldml",
		root: "LDML",
		top: []string{
			"ldml", "collation", "calendar", "timeZoneNames", "localeDisplayNames", "numbers",
		},
		skipElem: []string{
			"cp",       
			"special",  
			"fallback", 
			"alias",    
			"default",  
		},
		skipAttr: []string{
			"hiraganaQuarternary", 
		},
		predefined: []string{"rules"},
	},
}

var comments = map[string]string{
	"ldmlBCP47": `

`,
	"supplementalData": `


`,
	"ldml": `

`,
	"collation": `




`,
	"calendar": `




`,
	"dates": `

`,
	"localeDisplayNames": `


`,
	"numbers": `

`,
}

type element struct {
	name      string 
	category  string 
	signature string 

	attr []*attribute 
	sub  []struct {   
		e      *element
		repeat bool 
	}

	resolved bool 
}

type attribute struct {
	name string
	key  string
	list []string

	tag string 
}

var (
	reHead  = regexp.MustCompile(` *(\w+) +([\w\-]+)`)
	reAttr  = regexp.MustCompile(` *(\w+) *(?:(\w+)|\(([\w\- \|]+)\)) *(?:#([A-Z]*) *(?:\"([\.\d+])\")?)? *("[\w\-:]*")?`)
	reElem  = regexp.MustCompile(`^ *(EMPTY|ANY|\(.*\)[\*\+\?]?) *$`)
	reToken = regexp.MustCompile(`\w\-`)
)



type builder struct {
	w       io.Writer
	index   map[string]*element
	elem    []*element
	info    dtd
	version string
}

func makeBuilder(w io.Writer, d dtd) builder {
	return builder{
		w:     w,
		index: make(map[string]*element),
		elem:  []*element{},
		info:  d,
	}
}


func (b *builder) parseDTD(r io.Reader) {
	for d := xml.NewDecoder(r); ; {
		t, err := d.Token()
		if t == nil {
			break
		}
		failOnError(err)
		dir, ok := t.(xml.Directive)
		if !ok {
			continue
		}
		m := reHead.FindSubmatch(dir)
		dir = dir[len(m[0]):]
		ename := string(m[2])
		el, elementFound := b.index[ename]
		switch string(m[1]) {
		case "ELEMENT":
			if elementFound {
				log.Fatal("parseDTD: duplicate entry for element %q", ename)
			}
			m := reElem.FindSubmatch(dir)
			if m == nil {
				log.Fatalf("parseDTD: invalid element %q", string(dir))
			}
			if len(m[0]) != len(dir) {
				log.Fatal("parseDTD: invalid element %q", string(dir), len(dir), len(m[0]), string(m[0]))
			}
			s := string(m[1])
			el = &element{
				name:     ename,
				category: s,
			}
			b.index[ename] = el
		case "ATTLIST":
			if !elementFound {
				log.Fatalf("parseDTD: unknown element %q", ename)
			}
			s := string(dir)
			m := reAttr.FindStringSubmatch(s)
			if m == nil {
				log.Fatal(fmt.Errorf("parseDTD: invalid attribute %q", string(dir)))
			}
			if m[4] == "FIXED" {
				b.version = m[5]
			} else {
				switch m[1] {
				case "draft", "references", "alt", "validSubLocales", "standard"  :
				case "type", "choice":
				default:
					el.attr = append(el.attr, &attribute{
						name: m[1],
						key:  s,
						list: reToken.FindAllString(m[3], -1),
					})
					el.signature = fmt.Sprintf("%s=%s+%s", el.signature, m[1], m[2])
				}
			}
		}
	}
}

var reCat = regexp.MustCompile(`[ ,\|]*(?:(\(|\)|\#?[\w_-]+)([\*\+\?]?))?`)



func (b *builder) resolve(e *element) {
	if e.resolved {
		return
	}
	b.elem = append(b.elem, e)
	e.resolved = true
	s := e.category
	found := make(map[string]bool)
	sequenceStart := []int{}
	for len(s) > 0 {
		m := reCat.FindStringSubmatch(s)
		if m == nil {
			log.Fatalf("%s: invalid category string %q", e.name, s)
		}
		repeat := m[2] == "*" || m[2] == "+" || in(b.info.forceRepeat, m[1])
		switch m[1] {
		case "":
		case "(":
			sequenceStart = append(sequenceStart, len(e.sub))
		case ")":
			if len(sequenceStart) == 0 {
				log.Fatalf("%s: unmatched closing parenthesis", e.name)
			}
			for i := sequenceStart[len(sequenceStart)-1]; i < len(e.sub); i++ {
				e.sub[i].repeat = e.sub[i].repeat || repeat
			}
			sequenceStart = sequenceStart[:len(sequenceStart)-1]
		default:
			if in(b.info.skipElem, m[1]) {
			} else if sub, ok := b.index[m[1]]; ok {
				if !found[sub.name] {
					e.sub = append(e.sub, struct {
						e      *element
						repeat bool
					}{sub, repeat})
					found[sub.name] = true
					b.resolve(sub)
				}
			} else if m[1] == "#PCDATA" || m[1] == "ANY" {
			} else if m[1] != "EMPTY" {
				log.Fatalf("resolve:%s: element %q not found", e.name, m[1])
			}
		}
		s = s[len(m[0]):]
	}
}


func in(set []string, s string) bool {
	for _, v := range set {
		if v == s {
			return true
		}
	}
	return false
}

var repl = strings.NewReplacer("-", " ", "_", " ")



func title(s string) string {
	return strings.Replace(strings.Title(repl.Replace(s)), " ", "", -1)
}


func (b *builder) writeElem(tab int, e *element) {
	p := func(f string, x ...interface{}) {
		f = strings.Replace(f, "\n", "\n"+strings.Repeat("\t", tab), -1)
		fmt.Fprintf(b.w, f, x...)
	}
	if len(e.sub) == 0 && len(e.attr) == 0 {
		p("Common")
		return
	}
	p("struct {")
	tab++
	p("\nCommon")
	for _, attr := range e.attr {
		if !in(b.info.skipAttr, attr.name) {
			p("\n%s string `xml:\"%s,attr\"`", title(attr.name), attr.name)
		}
	}
	for _, sub := range e.sub {
		if in(b.info.predefined, sub.e.name) {
			p("\n%sElem", sub.e.name)
			continue
		}
		if in(b.info.skipElem, sub.e.name) {
			continue
		}
		p("\n%s ", title(sub.e.name))
		if sub.repeat {
			p("[]")
		}
		p("*")
		if in(b.info.top, sub.e.name) {
			p(title(sub.e.name))
		} else {
			b.writeElem(tab, sub.e)
		}
		p(" `xml:\"%s\"`", sub.e.name)
	}
	tab--
	p("\n}")
}


func (b *builder) write() {
	for i, name := range b.info.top {
		e := b.index[name]
		if e != nil {
			fmt.Fprintf(b.w, comments[name])
			name := title(e.name)
			if i == 0 {
				name = b.info.root
			}
			fmt.Fprintf(b.w, "type %s ", name)
			b.writeElem(0, e)
			fmt.Fprint(b.w, "\n")
		}
	}
}
