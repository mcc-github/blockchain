



package gen

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)









type CodeWriter struct {
	buf  bytes.Buffer
	Size int
	Hash hash.Hash32 
	gob  *gob.Encoder
	
	
	skipSep bool
}

func (w *CodeWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}


func NewCodeWriter() *CodeWriter {
	h := fnv.New32()
	return &CodeWriter{Hash: h, gob: gob.NewEncoder(h)}
}



func (w *CodeWriter) WriteGoFile(filename, pkg string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file %s: %v", filename, err)
	}
	defer f.Close()
	if _, err = w.WriteGo(f, pkg, ""); err != nil {
		log.Fatalf("Error writing file %s: %v", filename, err)
	}
}




func (w *CodeWriter) WriteVersionedGoFile(filename, pkg string) {
	tags := buildTags()
	if tags != "" {
		filename = insertVersion(filename, UnicodeVersion())
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file %s: %v", filename, err)
	}
	defer f.Close()
	if _, err = w.WriteGo(f, pkg, tags); err != nil {
		log.Fatalf("Error writing file %s: %v", filename, err)
	}
}



func (w *CodeWriter) WriteGo(out io.Writer, pkg, tags string) (n int, err error) {
	sz := w.Size
	w.WriteComment("Total table size %d bytes (%dKiB); checksum: %X\n", sz, sz/1024, w.Hash.Sum32())
	defer w.buf.Reset()
	return WriteGo(out, pkg, tags, w.buf.Bytes())
}

func (w *CodeWriter) printf(f string, x ...interface{}) {
	fmt.Fprintf(w, f, x...)
}

func (w *CodeWriter) insertSep() {
	if w.skipSep {
		w.skipSep = false
		return
	}
	
	
	w.printf("\n\n")
}




func (w *CodeWriter) WriteComment(comment string, args ...interface{}) {
	s := fmt.Sprintf(comment, args...)
	s = strings.Trim(s, "\n")

	
	
	w.printf("\n\n
	w.skipSep = true

	
	sep := "\n"
	for ; len(s) > 0 && (s[0] == '\t' || s[0] == ' '); s = s[1:] {
		sep += s[:1]
	}

	strings.NewReplacer(sep, "\n

	w.printf("\n")
}

func (w *CodeWriter) writeSizeInfo(size int) {
	w.printf("
}


func (w *CodeWriter) WriteConst(name string, x interface{}) {
	w.insertSep()
	v := reflect.ValueOf(x)

	switch v.Type().Kind() {
	case reflect.String:
		w.printf("const %s %s = ", name, typeName(x))
		w.WriteString(v.String())
		w.printf("\n")
	default:
		w.printf("const %s = %#v\n", name, x)
	}
}


func (w *CodeWriter) WriteVar(name string, x interface{}) {
	w.insertSep()
	v := reflect.ValueOf(x)
	oldSize := w.Size
	sz := int(v.Type().Size())
	w.Size += sz

	switch v.Type().Kind() {
	case reflect.String:
		w.printf("var %s %s = ", name, typeName(x))
		w.WriteString(v.String())
	case reflect.Struct:
		w.gob.Encode(x)
		fallthrough
	case reflect.Slice, reflect.Array:
		w.printf("var %s = ", name)
		w.writeValue(v)
		w.writeSizeInfo(w.Size - oldSize)
	default:
		w.printf("var %s %s = ", name, typeName(x))
		w.gob.Encode(x)
		w.writeValue(v)
		w.writeSizeInfo(w.Size - oldSize)
	}
	w.printf("\n")
}

func (w *CodeWriter) writeValue(v reflect.Value) {
	x := v.Interface()
	switch v.Kind() {
	case reflect.String:
		w.WriteString(v.String())
	case reflect.Array:
		
		
		w.Size -= int(v.Type().Size())
		w.writeSlice(x, true)
	case reflect.Slice:
		w.writeSlice(x, false)
	case reflect.Struct:
		w.printf("%s{\n", typeName(v.Interface()))
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			w.printf("%s: ", t.Field(i).Name)
			w.writeValue(v.Field(i))
			w.printf(",\n")
		}
		w.printf("}")
	default:
		w.printf("%#v", x)
	}
}


func (w *CodeWriter) WriteString(s string) {
	s = strings.Replace(s, `\`, `\\`, -1)
	io.WriteString(w.Hash, s) 
	w.Size += len(s)

	const maxInline = 40
	if len(s) <= maxInline {
		w.printf("%q", s)
		return
	}

	
	const maxWidth = 80 - 4 - len(`"`) - len(`" +`)

	
	n, max := maxWidth, maxWidth-4

	
	
	
	
	
	explicitParens, extraComment := len(s) > 128*1024, ""
	if explicitParens {
		w.printf(`(`)
		extraComment = "; the redundant, explicit parens are for https://golang.org/issue/18078"
	}

	
	b := w.buf.Bytes()
	if p := len(bytes.TrimRight(b, " \t")); p > 0 && b[p-1] != '\n' {
		w.printf("\"\" + 
		n, max = maxWidth, maxWidth
	}

	w.printf(`"`)

	for sz, p, nLines := 0, 0, 0; p < len(s); {
		var r rune
		r, sz = utf8.DecodeRuneInString(s[p:])
		out := s[p : p+sz]
		chars := 1
		if !unicode.IsPrint(r) || r == utf8.RuneError || r == '"' {
			switch sz {
			case 1:
				out = fmt.Sprintf("\\x%02x", s[p])
			case 2, 3:
				out = fmt.Sprintf("\\u%04x", r)
			case 4:
				out = fmt.Sprintf("\\U%08x", r)
			}
			chars = len(out)
		}
		if n -= chars; n < 0 {
			nLines++
			if explicitParens && nLines&63 == 63 {
				w.printf("\") + (\"")
			}
			w.printf("\" +\n\"")
			n = max - len(out)
		}
		w.printf("%s", out)
		p += sz
	}
	w.printf(`"`)
	if explicitParens {
		w.printf(`)`)
	}
}


func (w *CodeWriter) WriteSlice(x interface{}) {
	w.writeSlice(x, false)
}


func (w *CodeWriter) WriteArray(x interface{}) {
	w.writeSlice(x, true)
}

func (w *CodeWriter) writeSlice(x interface{}, isArray bool) {
	v := reflect.ValueOf(x)
	w.gob.Encode(v.Len())
	w.Size += v.Len() * int(v.Type().Elem().Size())
	name := typeName(x)
	if isArray {
		name = fmt.Sprintf("[%d]%s", v.Len(), name[strings.Index(name, "]")+1:])
	}
	if isArray {
		w.printf("%s{\n", name)
	} else {
		w.printf("%s{ 
	}

	switch kind := v.Type().Elem().Kind(); kind {
	case reflect.String:
		for _, s := range x.([]string) {
			w.WriteString(s)
			w.printf(",\n")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		
		nLine, nBlock, format := 8, 64, "%d,"
		switch kind {
		case reflect.Uint8:
			format = "%#02x,"
		case reflect.Uint16:
			format = "%#04x,"
		case reflect.Uint32:
			nLine, nBlock, format = 4, 32, "%#08x,"
		case reflect.Uint, reflect.Uint64:
			nLine, nBlock, format = 4, 32, "%#016x,"
		case reflect.Int8:
			nLine = 16
		}
		n := nLine
		for i := 0; i < v.Len(); i++ {
			if i%nBlock == 0 && v.Len() > nBlock {
				w.printf("
			}
			x := v.Index(i).Interface()
			w.gob.Encode(x)
			w.printf(format, x)
			if n--; n == 0 {
				n = nLine
				w.printf("\n")
			}
		}
		w.printf("\n")
	case reflect.Struct:
		zero := reflect.Zero(v.Type().Elem()).Interface()
		for i := 0; i < v.Len(); i++ {
			x := v.Index(i).Interface()
			w.gob.EncodeValue(v)
			if !reflect.DeepEqual(zero, x) {
				line := fmt.Sprintf("%#v,\n", x)
				line = line[strings.IndexByte(line, '{'):]
				w.printf("%d: ", i)
				w.printf(line)
			}
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			w.printf("%d: %#v,\n", i, v.Index(i).Interface())
		}
	default:
		panic("gen: slice elem type not supported")
	}
	w.printf("}")
}



func (w *CodeWriter) WriteType(x interface{}) string {
	t := reflect.TypeOf(x)
	w.printf("type %s struct {\n", t.Name())
	for i := 0; i < t.NumField(); i++ {
		w.printf("\t%s %s\n", t.Field(i).Name, t.Field(i).Type)
	}
	w.printf("}\n")
	return t.Name()
}


func typeName(x interface{}) string {
	t := reflect.ValueOf(x).Type()
	return strings.Replace(fmt.Sprint(t), "main.", "", 1)
}
