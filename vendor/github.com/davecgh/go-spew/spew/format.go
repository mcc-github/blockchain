

package spew

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)


const supportedFlags = "0-+# "





type formatState struct {
	value          interface{}
	fs             fmt.State
	depth          int
	pointers       map[uintptr]int
	ignoreNextType bool
	cs             *ConfigState
}





func (f *formatState) buildDefaultFormat() (format string) {
	buf := bytes.NewBuffer(percentBytes)

	for _, flag := range supportedFlags {
		if f.fs.Flag(int(flag)) {
			buf.WriteRune(flag)
		}
	}

	buf.WriteRune('v')

	format = buf.String()
	return format
}




func (f *formatState) constructOrigFormat(verb rune) (format string) {
	buf := bytes.NewBuffer(percentBytes)

	for _, flag := range supportedFlags {
		if f.fs.Flag(int(flag)) {
			buf.WriteRune(flag)
		}
	}

	if width, ok := f.fs.Width(); ok {
		buf.WriteString(strconv.Itoa(width))
	}

	if precision, ok := f.fs.Precision(); ok {
		buf.Write(precisionBytes)
		buf.WriteString(strconv.Itoa(precision))
	}

	buf.WriteRune(verb)

	format = buf.String()
	return format
}






func (f *formatState) unpackValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		f.ignoreNextType = false
		if !v.IsNil() {
			v = v.Elem()
		}
	}
	return v
}


func (f *formatState) formatPtr(v reflect.Value) {
	
	showTypes := f.fs.Flag('#')
	if v.IsNil() && (!showTypes || f.ignoreNextType) {
		f.fs.Write(nilAngleBytes)
		return
	}

	
	
	for k, depth := range f.pointers {
		if depth >= f.depth {
			delete(f.pointers, k)
		}
	}

	
	pointerChain := make([]uintptr, 0)

	
	
	
	nilFound := false
	cycleFound := false
	indirects := 0
	ve := v
	for ve.Kind() == reflect.Ptr {
		if ve.IsNil() {
			nilFound = true
			break
		}
		indirects++
		addr := ve.Pointer()
		pointerChain = append(pointerChain, addr)
		if pd, ok := f.pointers[addr]; ok && pd < f.depth {
			cycleFound = true
			indirects--
			break
		}
		f.pointers[addr] = f.depth

		ve = ve.Elem()
		if ve.Kind() == reflect.Interface {
			if ve.IsNil() {
				nilFound = true
				break
			}
			ve = ve.Elem()
		}
	}

	
	if showTypes && !f.ignoreNextType {
		f.fs.Write(openParenBytes)
		f.fs.Write(bytes.Repeat(asteriskBytes, indirects))
		f.fs.Write([]byte(ve.Type().String()))
		f.fs.Write(closeParenBytes)
	} else {
		if nilFound || cycleFound {
			indirects += strings.Count(ve.Type().String(), "*")
		}
		f.fs.Write(openAngleBytes)
		f.fs.Write([]byte(strings.Repeat("*", indirects)))
		f.fs.Write(closeAngleBytes)
	}

	
	if f.fs.Flag('+') && (len(pointerChain) > 0) {
		f.fs.Write(openParenBytes)
		for i, addr := range pointerChain {
			if i > 0 {
				f.fs.Write(pointerChainBytes)
			}
			printHexPtr(f.fs, addr)
		}
		f.fs.Write(closeParenBytes)
	}

	
	switch {
	case nilFound == true:
		f.fs.Write(nilAngleBytes)

	case cycleFound == true:
		f.fs.Write(circularShortBytes)

	default:
		f.ignoreNextType = true
		f.format(ve)
	}
}





func (f *formatState) format(v reflect.Value) {
	
	kind := v.Kind()
	if kind == reflect.Invalid {
		f.fs.Write(invalidAngleBytes)
		return
	}

	
	if kind == reflect.Ptr {
		f.formatPtr(v)
		return
	}

	
	if !f.ignoreNextType && f.fs.Flag('#') {
		f.fs.Write(openParenBytes)
		f.fs.Write([]byte(v.Type().String()))
		f.fs.Write(closeParenBytes)
	}
	f.ignoreNextType = false

	
	
	if !f.cs.DisableMethods {
		if (kind != reflect.Invalid) && (kind != reflect.Interface) {
			if handled := handleMethods(f.cs, f.fs, v); handled {
				return
			}
		}
	}

	switch kind {
	case reflect.Invalid:
		
		

	case reflect.Bool:
		printBool(f.fs, v.Bool())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		printInt(f.fs, v.Int(), 10)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		printUint(f.fs, v.Uint(), 10)

	case reflect.Float32:
		printFloat(f.fs, v.Float(), 32)

	case reflect.Float64:
		printFloat(f.fs, v.Float(), 64)

	case reflect.Complex64:
		printComplex(f.fs, v.Complex(), 32)

	case reflect.Complex128:
		printComplex(f.fs, v.Complex(), 64)

	case reflect.Slice:
		if v.IsNil() {
			f.fs.Write(nilAngleBytes)
			break
		}
		fallthrough

	case reflect.Array:
		f.fs.Write(openBracketBytes)
		f.depth++
		if (f.cs.MaxDepth != 0) && (f.depth > f.cs.MaxDepth) {
			f.fs.Write(maxShortBytes)
		} else {
			numEntries := v.Len()
			for i := 0; i < numEntries; i++ {
				if i > 0 {
					f.fs.Write(spaceBytes)
				}
				f.ignoreNextType = true
				f.format(f.unpackValue(v.Index(i)))
			}
		}
		f.depth--
		f.fs.Write(closeBracketBytes)

	case reflect.String:
		f.fs.Write([]byte(v.String()))

	case reflect.Interface:
		
		
		if v.IsNil() {
			f.fs.Write(nilAngleBytes)
		}

	case reflect.Ptr:
		
		

	case reflect.Map:
		
		if v.IsNil() {
			f.fs.Write(nilAngleBytes)
			break
		}

		f.fs.Write(openMapBytes)
		f.depth++
		if (f.cs.MaxDepth != 0) && (f.depth > f.cs.MaxDepth) {
			f.fs.Write(maxShortBytes)
		} else {
			keys := v.MapKeys()
			if f.cs.SortKeys {
				sortValues(keys, f.cs)
			}
			for i, key := range keys {
				if i > 0 {
					f.fs.Write(spaceBytes)
				}
				f.ignoreNextType = true
				f.format(f.unpackValue(key))
				f.fs.Write(colonBytes)
				f.ignoreNextType = true
				f.format(f.unpackValue(v.MapIndex(key)))
			}
		}
		f.depth--
		f.fs.Write(closeMapBytes)

	case reflect.Struct:
		numFields := v.NumField()
		f.fs.Write(openBraceBytes)
		f.depth++
		if (f.cs.MaxDepth != 0) && (f.depth > f.cs.MaxDepth) {
			f.fs.Write(maxShortBytes)
		} else {
			vt := v.Type()
			for i := 0; i < numFields; i++ {
				if i > 0 {
					f.fs.Write(spaceBytes)
				}
				vtf := vt.Field(i)
				if f.fs.Flag('+') || f.fs.Flag('#') {
					f.fs.Write([]byte(vtf.Name))
					f.fs.Write(colonBytes)
				}
				f.format(f.unpackValue(v.Field(i)))
			}
		}
		f.depth--
		f.fs.Write(closeBraceBytes)

	case reflect.Uintptr:
		printHexPtr(f.fs, uintptr(v.Uint()))

	case reflect.UnsafePointer, reflect.Chan, reflect.Func:
		printHexPtr(f.fs, v.Pointer())

	
	
	default:
		format := f.buildDefaultFormat()
		if v.CanInterface() {
			fmt.Fprintf(f.fs, format, v.Interface())
		} else {
			fmt.Fprintf(f.fs, format, v.String())
		}
	}
}



func (f *formatState) Format(fs fmt.State, verb rune) {
	f.fs = fs

	
	if verb != 'v' {
		format := f.constructOrigFormat(verb)
		fmt.Fprintf(fs, format, f.value)
		return
	}

	if f.value == nil {
		if fs.Flag('#') {
			fs.Write(interfaceBytes)
		}
		fs.Write(nilAngleBytes)
		return
	}

	f.format(reflect.ValueOf(f.value))
}



func newFormatter(cs *ConfigState, v interface{}) fmt.Formatter {
	fs := &formatState{value: v, cs: cs}
	fs.pointers = make(map[uintptr]int)
	return fs
}


func NewFormatter(v interface{}) fmt.Formatter {
	return newFormatter(&Config, v)
}
