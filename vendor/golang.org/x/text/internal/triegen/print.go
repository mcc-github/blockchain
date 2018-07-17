



package triegen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)



func (b *builder) print(w io.Writer) error {
	b.Stats.NValueEntries = len(b.ValueBlocks) * blockSize
	b.Stats.NValueBytes = len(b.ValueBlocks) * blockSize * b.ValueSize
	b.Stats.NIndexEntries = len(b.IndexBlocks) * blockSize
	b.Stats.NIndexBytes = len(b.IndexBlocks) * blockSize * b.IndexSize
	b.Stats.NHandleBytes = len(b.Trie) * 2 * b.IndexSize

	
	
	if len(b.Trie) == 1 {
		
		b.ASCIIBlock = b.Name + "Values"
		b.StarterBlock = b.Name + "Index"
	} else {
		
		
		b.ASCIIBlock = "t.ascii"
		b.StarterBlock = "t.utf8Start"
	}

	b.SourceType = "[]byte"
	if err := lookupGen.Execute(w, b); err != nil {
		return err
	}

	b.SourceType = "string"
	if err := lookupGen.Execute(w, b); err != nil {
		return err
	}

	if err := trieGen.Execute(w, b); err != nil {
		return err
	}

	for _, c := range b.Compactions {
		if err := c.c.Print(w); err != nil {
			return err
		}
	}

	return nil
}

func printValues(n int, values []uint64) string {
	w := &bytes.Buffer{}
	boff := n * blockSize
	fmt.Fprintf(w, "\t
	var newline bool
	for i, v := range values {
		if i%6 == 0 {
			newline = true
		}
		if v != 0 {
			if newline {
				fmt.Fprintf(w, "\n")
				newline = false
			}
			fmt.Fprintf(w, "\t%#02x:%#04x, ", boff+i, v)
		}
	}
	return w.String()
}

func printIndex(b *builder, nr int, n *node) string {
	w := &bytes.Buffer{}
	boff := nr * blockSize
	fmt.Fprintf(w, "\t
	var newline bool
	for i, c := range n.children {
		if i%8 == 0 {
			newline = true
		}
		if c != nil {
			v := b.Compactions[c.index.compaction].Offset + uint32(c.index.index)
			if v != 0 {
				if newline {
					fmt.Fprintf(w, "\n")
					newline = false
				}
				fmt.Fprintf(w, "\t%#02x:%#02x, ", boff+i, v)
			}
		}
	}
	return w.String()
}

var (
	trieGen = template.Must(template.New("trie").Funcs(template.FuncMap{
		"printValues": printValues,
		"printIndex":  printIndex,
		"title":       strings.Title,
		"dec":         func(x int) int { return x - 1 },
		"psize": func(n int) string {
			return fmt.Sprintf("%d bytes (%.2f KiB)", n, float64(n)/1024)
		},
	}).Parse(trieTemplate))
	lookupGen = template.Must(template.New("lookup").Parse(lookupTemplate))
)




const trieTemplate = `{{$b := .}}{{$multi := gt (len .Trie) 1}}

type {{.Name}}Trie struct { {{if $multi}}
	ascii []{{.ValueType}} 
	utf8Start  []{{.IndexType}} 
{{end}}}

func new{{title .Name}}Trie(i int) *{{.Name}}Trie { {{if $multi}}
	h := {{.Name}}TrieHandles[i]
	return &{{.Name}}Trie{ {{.Name}}Values[uint32(h.ascii)<<6:], {{.Name}}Index[uint32(h.multi)<<6:] }
}

type {{.Name}}TrieHandle struct {
	ascii, multi {{.IndexType}}
}


var {{.Name}}TrieHandles = [{{len .Trie}}]{{.Name}}TrieHandle{
{{range .Trie}}	{ {{.ASCIIIndex}}, {{.StarterIndex}} }, 
{{end}}}{{else}}
	return &{{.Name}}Trie{}
}
{{end}}

func (t *{{.Name}}Trie) lookupValue(n uint32, b byte) {{.ValueType}}{{$last := dec (len .Compactions)}} {
	switch { {{range $i, $c := .Compactions}}
		{{if eq $i $last}}default{{else}}case n < {{$c.Cutoff}}{{end}}:{{if ne $i 0}}
			n -= {{$c.Offset}}{{end}}
			return {{print $b.ValueType}}({{$c.Handler}}){{end}}
	}
}



var {{.Name}}Values = [{{.Stats.NValueEntries}}]{{.ValueType}} {
{{range $i, $v := .ValueBlocks}}{{printValues $i $v}}
{{end}}}



var {{.Name}}Index = [{{.Stats.NIndexEntries}}]{{.IndexType}} {
{{range $i, $v := .IndexBlocks}}{{printIndex $b $i $v}}
{{end}}}
`



const lookupTemplate = `



func (t *{{.Name}}Trie) lookup{{if eq .SourceType "string"}}String{{end}}(s {{.SourceType}}) (v {{.ValueType}}, sz int) {
	c0 := s[0]
	switch {
	case c0 < 0x80: 
		return {{.ASCIIBlock}}[c0], 1
	case c0 < 0xC2:
		return 0, 1  
	case c0 < 0xE0: 
		if len(s) < 2 {
			return 0, 0
		}
		i := {{.StarterBlock}}[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 {
			return 0, 1 
		}
		return t.lookupValue(uint32(i), c1), 2
	case c0 < 0xF0: 
		if len(s) < 3 {
			return 0, 0
		}
		i := {{.StarterBlock}}[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 {
			return 0, 1 
		}
		o := uint32(i)<<6 + uint32(c1)
		i = {{.Name}}Index[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 {
			return 0, 2 
		}
		return t.lookupValue(uint32(i), c2), 3
	case c0 < 0xF8: 
		if len(s) < 4 {
			return 0, 0
		}
		i := {{.StarterBlock}}[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 {
			return 0, 1 
		}
		o := uint32(i)<<6 + uint32(c1)
		i = {{.Name}}Index[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 {
			return 0, 2 
		}
		o = uint32(i)<<6 + uint32(c2)
		i = {{.Name}}Index[o]
		c3 := s[3]
		if c3 < 0x80 || 0xC0 <= c3 {
			return 0, 3 
		}
		return t.lookupValue(uint32(i), c3), 4
	}
	
	return 0, 1
}



func (t *{{.Name}}Trie) lookup{{if eq .SourceType "string"}}String{{end}}Unsafe(s {{.SourceType}}) {{.ValueType}} {
	c0 := s[0]
	if c0 < 0x80 { 
		return {{.ASCIIBlock}}[c0]
	}
	i := {{.StarterBlock}}[c0]
	if c0 < 0xE0 { 
		return t.lookupValue(uint32(i), s[1])
	}
	i = {{.Name}}Index[uint32(i)<<6+uint32(s[1])]
	if c0 < 0xF0 { 
		return t.lookupValue(uint32(i), s[2])
	}
	i = {{.Name}}Index[uint32(i)<<6+uint32(s[2])]
	if c0 < 0xF8 { 
		return t.lookupValue(uint32(i), s[3])
	}
	return 0
}
`
