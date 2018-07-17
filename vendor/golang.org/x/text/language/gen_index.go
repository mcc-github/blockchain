





package main



import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var (
	test = flag.Bool("test", false,
		"test existing tables; can be used to compare web data with package data.")

	draft = flag.String("draft",
		"contributed",
		`Minimal draft requirements (approved, contributed, provisional, unconfirmed).`)
)

func main() {
	gen.Init()

	
	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder{}
	data, err := d.DecodeZip(r)
	if err != nil {
		log.Fatalf("DecodeZip: %v", err)
	}

	w := gen.NewCodeWriter()
	defer func() {
		buf := &bytes.Buffer{}

		if _, err = w.WriteGo(buf, "language", ""); err != nil {
			log.Fatalf("Error formatting file index.go: %v", err)
		}

		
		
		
		out := bytes.Replace(buf.Bytes(), []byte("language."), nil, -1)
		if err := ioutil.WriteFile("index.go", out, 0600); err != nil {
			log.Fatalf("Could not create file index.go: %v", err)
		}
	}()

	m := map[language.Tag]bool{}
	for _, lang := range data.Locales() {
		
		

		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		

		
		tag := language.Make(strings.Replace(lang, "_POSIX", "-u-va-posix", 1))
		m[tag] = true
		
	}
	
	for _, plurals := range data.Supplemental().Plurals {
		for _, rules := range plurals.PluralRules {
			for _, lang := range strings.Split(rules.Locales, " ") {
				m[language.Make(lang)] = true
			}
		}
	}

	var core, special []language.Tag

	for t := range m {
		if x := t.Extensions(); len(x) != 0 && fmt.Sprint(x) != "[u-va-posix]" {
			log.Fatalf("Unexpected extension %v in %v", x, t)
		}
		if len(t.Variants()) == 0 && len(t.Extensions()) == 0 {
			core = append(core, t)
		} else {
			special = append(special, t)
		}
	}

	w.WriteComment(`
	NumCompactTags is the number of common tags. The maximum tag is
	NumCompactTags-1.`)
	w.WriteConst("NumCompactTags", len(core)+len(special))

	sort.Sort(byAlpha(special))
	w.WriteVar("specialTags", special)

	
	sort.Sort(byAlpha(core))

	
	w.Size += int(reflect.TypeOf(map[uint32]uint16{}).Size())
	w.Size += len(core) * 6 

	fmt.Fprintln(w)
	fmt.Fprintln(w, "var coreTags = map[uint32]uint16{")
	fmt.Fprintln(w, "0x0: 0, 
	i := len(special) + 1 
	for _, t := range core {
		if t == language.Und {
			continue
		}
		fmt.Fprint(w.Hash, t, i)
		b, s, r := t.Raw()
		fmt.Fprintf(w, "0x%s%s%s: %d, 
			getIndex(b, 3), 
			getIndex(s, 2),
			getIndex(r, 3),
			i, t)
		i++
	}
	fmt.Fprintln(w, "}")
}



func getIndex(x interface{}, n int) string {
	s := fmt.Sprintf("%#v", x) 
	s = s[strings.Index(s, "0x")+2 : len(s)-1]
	return strings.Repeat("0", n-len(s)) + s
}

type byAlpha []language.Tag

func (a byAlpha) Len() int           { return len(a) }
func (a byAlpha) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAlpha) Less(i, j int) bool { return a[i].String() < a[j].String() }
