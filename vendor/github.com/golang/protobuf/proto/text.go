






























package proto



import (
	"bufio"
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
	"sort"
	"strings"
)

var (
	newline         = []byte("\n")
	spaces          = []byte("                                        ")
	endBraceNewline = []byte("}\n")
	backslashN      = []byte{'\\', 'n'}
	backslashR      = []byte{'\\', 'r'}
	backslashT      = []byte{'\\', 't'}
	backslashDQ     = []byte{'\\', '"'}
	backslashBS     = []byte{'\\', '\\'}
	posInf          = []byte("inf")
	negInf          = []byte("-inf")
	nan             = []byte("nan")
)

type writer interface {
	io.Writer
	WriteByte(byte) error
}


type textWriter struct {
	ind      int
	complete bool 
	compact  bool 
	w        writer
}

func (w *textWriter) WriteString(s string) (n int, err error) {
	if !strings.Contains(s, "\n") {
		if !w.compact && w.complete {
			w.writeIndent()
		}
		w.complete = false
		return io.WriteString(w.w, s)
	}
	
	
	
	return w.Write([]byte(s))
}

func (w *textWriter) Write(p []byte) (n int, err error) {
	newlines := bytes.Count(p, newline)
	if newlines == 0 {
		if !w.compact && w.complete {
			w.writeIndent()
		}
		n, err = w.w.Write(p)
		w.complete = false
		return n, err
	}

	frags := bytes.SplitN(p, newline, newlines+1)
	if w.compact {
		for i, frag := range frags {
			if i > 0 {
				if err := w.w.WriteByte(' '); err != nil {
					return n, err
				}
				n++
			}
			nn, err := w.w.Write(frag)
			n += nn
			if err != nil {
				return n, err
			}
		}
		return n, nil
	}

	for i, frag := range frags {
		if w.complete {
			w.writeIndent()
		}
		nn, err := w.w.Write(frag)
		n += nn
		if err != nil {
			return n, err
		}
		if i+1 < len(frags) {
			if err := w.w.WriteByte('\n'); err != nil {
				return n, err
			}
			n++
		}
	}
	w.complete = len(frags[len(frags)-1]) == 0
	return n, nil
}

func (w *textWriter) WriteByte(c byte) error {
	if w.compact && c == '\n' {
		c = ' '
	}
	if !w.compact && w.complete {
		w.writeIndent()
	}
	err := w.w.WriteByte(c)
	w.complete = c == '\n'
	return err
}

func (w *textWriter) indent() { w.ind++ }

func (w *textWriter) unindent() {
	if w.ind == 0 {
		log.Print("proto: textWriter unindented too far")
		return
	}
	w.ind--
}

func writeName(w *textWriter, props *Properties) error {
	if _, err := w.WriteString(props.OrigName); err != nil {
		return err
	}
	if props.Wire != "group" {
		return w.WriteByte(':')
	}
	return nil
}

func requiresQuotes(u string) bool {
	
	for _, ch := range u {
		switch {
		case ch == '.' || ch == '/' || ch == '_':
			continue
		case '0' <= ch && ch <= '9':
			continue
		case 'A' <= ch && ch <= 'Z':
			continue
		case 'a' <= ch && ch <= 'z':
			continue
		default:
			return true
		}
	}
	return false
}


func isAny(sv reflect.Value) bool {
	type wkt interface {
		XXX_WellKnownType() string
	}
	t, ok := sv.Addr().Interface().(wkt)
	return ok && t.XXX_WellKnownType() == "Any"
}








func (tm *TextMarshaler) writeProto3Any(w *textWriter, sv reflect.Value) (bool, error) {
	turl := sv.FieldByName("TypeUrl")
	val := sv.FieldByName("Value")
	if !turl.IsValid() || !val.IsValid() {
		return true, errors.New("proto: invalid google.protobuf.Any message")
	}

	b, ok := val.Interface().([]byte)
	if !ok {
		return true, errors.New("proto: invalid google.protobuf.Any message")
	}

	parts := strings.Split(turl.String(), "/")
	mt := MessageType(parts[len(parts)-1])
	if mt == nil {
		return false, nil
	}
	m := reflect.New(mt.Elem())
	if err := Unmarshal(b, m.Interface().(Message)); err != nil {
		return false, nil
	}
	w.Write([]byte("["))
	u := turl.String()
	if requiresQuotes(u) {
		writeString(w, u)
	} else {
		w.Write([]byte(u))
	}
	if w.compact {
		w.Write([]byte("]:<"))
	} else {
		w.Write([]byte("]: <\n"))
		w.ind++
	}
	if err := tm.writeStruct(w, m.Elem()); err != nil {
		return true, err
	}
	if w.compact {
		w.Write([]byte("> "))
	} else {
		w.ind--
		w.Write([]byte(">\n"))
	}
	return true, nil
}

func (tm *TextMarshaler) writeStruct(w *textWriter, sv reflect.Value) error {
	if tm.ExpandAny && isAny(sv) {
		if canExpand, err := tm.writeProto3Any(w, sv); canExpand {
			return err
		}
	}
	st := sv.Type()
	sprops := GetProperties(st)
	for i := 0; i < sv.NumField(); i++ {
		fv := sv.Field(i)
		props := sprops.Prop[i]
		name := st.Field(i).Name

		if name == "XXX_NoUnkeyedLiteral" {
			continue
		}

		if strings.HasPrefix(name, "XXX_") {
			
			
			
			
			
			if name == "XXX_unrecognized" && !fv.IsNil() {
				if err := writeUnknownStruct(w, fv.Interface().([]byte)); err != nil {
					return err
				}
			}
			continue
		}
		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			
			
			
			continue
		}
		if fv.Kind() == reflect.Slice && fv.IsNil() {
			
			continue
		}

		if props.Repeated && fv.Kind() == reflect.Slice {
			
			for j := 0; j < fv.Len(); j++ {
				if err := writeName(w, props); err != nil {
					return err
				}
				if !w.compact {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
				}
				v := fv.Index(j)
				if v.Kind() == reflect.Ptr && v.IsNil() {
					
					
					if _, err := w.Write([]byte("<nil>\n")); err != nil {
						return err
					}
					continue
				}
				if err := tm.writeAny(w, v, props); err != nil {
					return err
				}
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
			continue
		}
		if fv.Kind() == reflect.Map {
			
			keys := fv.MapKeys()
			sort.Sort(mapKeys(keys))
			for _, key := range keys {
				val := fv.MapIndex(key)
				if err := writeName(w, props); err != nil {
					return err
				}
				if !w.compact {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
				}
				
				if err := w.WriteByte('<'); err != nil {
					return err
				}
				if !w.compact {
					if err := w.WriteByte('\n'); err != nil {
						return err
					}
				}
				w.indent()
				
				if _, err := w.WriteString("key:"); err != nil {
					return err
				}
				if !w.compact {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
				}
				if err := tm.writeAny(w, key, props.mkeyprop); err != nil {
					return err
				}
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
				
				if val.Kind() != reflect.Ptr || !val.IsNil() {
					
					if _, err := w.WriteString("value:"); err != nil {
						return err
					}
					if !w.compact {
						if err := w.WriteByte(' '); err != nil {
							return err
						}
					}
					if err := tm.writeAny(w, val, props.mvalprop); err != nil {
						return err
					}
					if err := w.WriteByte('\n'); err != nil {
						return err
					}
				}
				
				w.unindent()
				if err := w.WriteByte('>'); err != nil {
					return err
				}
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
			continue
		}
		if props.proto3 && fv.Kind() == reflect.Slice && fv.Len() == 0 {
			
			continue
		}
		if fv.Kind() != reflect.Ptr && fv.Kind() != reflect.Slice {
			
			if isProto3Zero(fv) {
				continue
			}
		}

		if fv.Kind() == reflect.Interface {
			
			if st.Field(i).Tag.Get("protobuf_oneof") != "" {
				
				
				
				if fv.IsNil() {
					continue
				}
				inner := fv.Elem().Elem() 
				tag := inner.Type().Field(0).Tag.Get("protobuf")
				props = new(Properties) 
				props.Parse(tag)
				
				fv = inner.Field(0)

				
				
				
				if fv.Kind() == reflect.Ptr && fv.IsNil() {
					
					msg := errors.New("")
					fv = reflect.ValueOf(&msg).Elem()
				}
			}
		}

		if err := writeName(w, props); err != nil {
			return err
		}
		if !w.compact {
			if err := w.WriteByte(' '); err != nil {
				return err
			}
		}

		
		if err := tm.writeAny(w, fv, props); err != nil {
			return err
		}

		if err := w.WriteByte('\n'); err != nil {
			return err
		}
	}

	
	pv := sv.Addr()
	if _, err := extendable(pv.Interface()); err == nil {
		if err := tm.writeExtensions(w, pv); err != nil {
			return err
		}
	}

	return nil
}


func (tm *TextMarshaler) writeAny(w *textWriter, v reflect.Value, props *Properties) error {
	v = reflect.Indirect(v)

	
	if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
		x := v.Float()
		var b []byte
		switch {
		case math.IsInf(x, 1):
			b = posInf
		case math.IsInf(x, -1):
			b = negInf
		case math.IsNaN(x):
			b = nan
		}
		if b != nil {
			_, err := w.Write(b)
			return err
		}
		
	}

	
	
	switch v.Kind() {
	case reflect.Slice:
		
		if err := writeString(w, string(v.Bytes())); err != nil {
			return err
		}
	case reflect.String:
		if err := writeString(w, v.String()); err != nil {
			return err
		}
	case reflect.Struct:
		
		var bra, ket byte = '<', '>'
		if props != nil && props.Wire == "group" {
			bra, ket = '{', '}'
		}
		if err := w.WriteByte(bra); err != nil {
			return err
		}
		if !w.compact {
			if err := w.WriteByte('\n'); err != nil {
				return err
			}
		}
		w.indent()
		if v.CanAddr() {
			
			
			
			
			
			
			
			
			
			
			v = v.Addr()
		}
		if etm, ok := v.Interface().(encoding.TextMarshaler); ok {
			text, err := etm.MarshalText()
			if err != nil {
				return err
			}
			if _, err = w.Write(text); err != nil {
				return err
			}
		} else {
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if err := tm.writeStruct(w, v); err != nil {
				return err
			}
		}
		w.unindent()
		if err := w.WriteByte(ket); err != nil {
			return err
		}
	default:
		_, err := fmt.Fprint(w, v.Interface())
		return err
	}
	return nil
}


func isprint(c byte) bool {
	return c >= 0x20 && c < 0x7f
}






func writeString(w *textWriter, s string) error {
	
	if err := w.WriteByte('"'); err != nil {
		return err
	}
	
	for i := 0; i < len(s); i++ {
		var err error
		
		
		
		switch c := s[i]; c {
		case '\n':
			_, err = w.w.Write(backslashN)
		case '\r':
			_, err = w.w.Write(backslashR)
		case '\t':
			_, err = w.w.Write(backslashT)
		case '"':
			_, err = w.w.Write(backslashDQ)
		case '\\':
			_, err = w.w.Write(backslashBS)
		default:
			if isprint(c) {
				err = w.w.WriteByte(c)
			} else {
				_, err = fmt.Fprintf(w.w, "\\%03o", c)
			}
		}
		if err != nil {
			return err
		}
	}
	return w.WriteByte('"')
}

func writeUnknownStruct(w *textWriter, data []byte) (err error) {
	if !w.compact {
		if _, err := fmt.Fprintf(w, "\n", len(data)); err != nil {
			return err
		}
	}
	b := NewBuffer(data)
	for b.index < len(b.buf) {
		x, err := b.DecodeVarint()
		if err != nil {
			_, err := fmt.Fprintf(w, "\n", err)
			return err
		}
		wire, tag := x&7, x>>3
		if wire == WireEndGroup {
			w.unindent()
			if _, err := w.Write(endBraceNewline); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprint(w, tag); err != nil {
			return err
		}
		if wire != WireStartGroup {
			if err := w.WriteByte(':'); err != nil {
				return err
			}
		}
		if !w.compact || wire == WireStartGroup {
			if err := w.WriteByte(' '); err != nil {
				return err
			}
		}
		switch wire {
		case WireBytes:
			buf, e := b.DecodeRawBytes(false)
			if e == nil {
				_, err = fmt.Fprintf(w, "%q", buf)
			} else {
				_, err = fmt.Fprintf(w, "", e)
			}
		case WireFixed32:
			x, err = b.DecodeFixed32()
			err = writeUnknownInt(w, x, err)
		case WireFixed64:
			x, err = b.DecodeFixed64()
			err = writeUnknownInt(w, x, err)
		case WireStartGroup:
			err = w.WriteByte('{')
			w.indent()
		case WireVarint:
			x, err = b.DecodeVarint()
			err = writeUnknownInt(w, x, err)
		default:
			_, err = fmt.Fprintf(w, "", wire)
		}
		if err != nil {
			return err
		}
		if err = w.WriteByte('\n'); err != nil {
			return err
		}
	}
	return nil
}

func writeUnknownInt(w *textWriter, x uint64, err error) error {
	if err == nil {
		_, err = fmt.Fprint(w, x)
	} else {
		_, err = fmt.Fprintf(w, "", err)
	}
	return err
}

type int32Slice []int32

func (s int32Slice) Len() int           { return len(s) }
func (s int32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s int32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }



func (tm *TextMarshaler) writeExtensions(w *textWriter, pv reflect.Value) error {
	emap := extensionMaps[pv.Type().Elem()]
	ep, _ := extendable(pv.Interface())

	
	
	
	m, mu := ep.extensionsRead()
	if m == nil {
		return nil
	}
	mu.Lock()
	ids := make([]int32, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	sort.Sort(int32Slice(ids))
	mu.Unlock()

	for _, extNum := range ids {
		ext := m[extNum]
		var desc *ExtensionDesc
		if emap != nil {
			desc = emap[extNum]
		}
		if desc == nil {
			
			if err := writeUnknownStruct(w, ext.enc); err != nil {
				return err
			}
			continue
		}

		pb, err := GetExtension(ep, desc)
		if err != nil {
			return fmt.Errorf("failed getting extension: %v", err)
		}

		
		if !desc.repeated() {
			if err := tm.writeExtension(w, desc.Name, pb); err != nil {
				return err
			}
		} else {
			v := reflect.ValueOf(pb)
			for i := 0; i < v.Len(); i++ {
				if err := tm.writeExtension(w, desc.Name, v.Index(i).Interface()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (tm *TextMarshaler) writeExtension(w *textWriter, name string, pb interface{}) error {
	if _, err := fmt.Fprintf(w, "[%s]:", name); err != nil {
		return err
	}
	if !w.compact {
		if err := w.WriteByte(' '); err != nil {
			return err
		}
	}
	if err := tm.writeAny(w, reflect.ValueOf(pb), nil); err != nil {
		return err
	}
	if err := w.WriteByte('\n'); err != nil {
		return err
	}
	return nil
}

func (w *textWriter) writeIndent() {
	if !w.complete {
		return
	}
	remain := w.ind * 2
	for remain > 0 {
		n := remain
		if n > len(spaces) {
			n = len(spaces)
		}
		w.w.Write(spaces[:n])
		remain -= n
	}
	w.complete = false
}


type TextMarshaler struct {
	Compact   bool 
	ExpandAny bool 
}



func (tm *TextMarshaler) Marshal(w io.Writer, pb Message) error {
	val := reflect.ValueOf(pb)
	if pb == nil || val.IsNil() {
		w.Write([]byte("<nil>"))
		return nil
	}
	var bw *bufio.Writer
	ww, ok := w.(writer)
	if !ok {
		bw = bufio.NewWriter(w)
		ww = bw
	}
	aw := &textWriter{
		w:        ww,
		complete: true,
		compact:  tm.Compact,
	}

	if etm, ok := pb.(encoding.TextMarshaler); ok {
		text, err := etm.MarshalText()
		if err != nil {
			return err
		}
		if _, err = aw.Write(text); err != nil {
			return err
		}
		if bw != nil {
			return bw.Flush()
		}
		return nil
	}
	
	v := reflect.Indirect(val)
	if err := tm.writeStruct(aw, v); err != nil {
		return err
	}
	if bw != nil {
		return bw.Flush()
	}
	return nil
}


func (tm *TextMarshaler) Text(pb Message) string {
	var buf bytes.Buffer
	tm.Marshal(&buf, pb)
	return buf.String()
}

var (
	defaultTextMarshaler = TextMarshaler{}
	compactTextMarshaler = TextMarshaler{Compact: true}
)





func MarshalText(w io.Writer, pb Message) error { return defaultTextMarshaler.Marshal(w, pb) }


func MarshalTextString(pb Message) string { return defaultTextMarshaler.Text(pb) }


func CompactText(w io.Writer, pb Message) error { return compactTextMarshaler.Marshal(w, pb) }


func CompactTextString(pb Message) string { return compactTextMarshaler.Text(pb) }
