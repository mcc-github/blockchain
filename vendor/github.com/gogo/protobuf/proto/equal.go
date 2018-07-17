
































package proto

import (
	"bytes"
	"log"
	"reflect"
	"strings"
)


func Equal(a, b Message) bool {
	if a == nil || b == nil {
		return a == b
	}
	v1, v2 := reflect.ValueOf(a), reflect.ValueOf(b)
	if v1.Type() != v2.Type() {
		return false
	}
	if v1.Kind() == reflect.Ptr {
		if v1.IsNil() {
			return v2.IsNil()
		}
		if v2.IsNil() {
			return false
		}
		v1, v2 = v1.Elem(), v2.Elem()
	}
	if v1.Kind() != reflect.Struct {
		return false
	}
	return equalStruct(v1, v2)
}


func equalStruct(v1, v2 reflect.Value) bool {
	sprop := GetProperties(v1.Type())
	for i := 0; i < v1.NumField(); i++ {
		f := v1.Type().Field(i)
		if strings.HasPrefix(f.Name, "XXX_") {
			continue
		}
		f1, f2 := v1.Field(i), v2.Field(i)
		if f.Type.Kind() == reflect.Ptr {
			if n1, n2 := f1.IsNil(), f2.IsNil(); n1 && n2 {
				
				continue
			} else if n1 != n2 {
				
				return false
			}
			b1, ok := f1.Interface().(raw)
			if ok {
				b2 := f2.Interface().(raw)
				
				if !bytes.Equal(b1.Bytes(), b2.Bytes()) {
					return false
				}
				continue
			}
			f1, f2 = f1.Elem(), f2.Elem()
		}
		if !equalAny(f1, f2, sprop.Prop[i]) {
			return false
		}
	}

	if em1 := v1.FieldByName("XXX_InternalExtensions"); em1.IsValid() {
		em2 := v2.FieldByName("XXX_InternalExtensions")
		if !equalExtensions(v1.Type(), em1.Interface().(XXX_InternalExtensions), em2.Interface().(XXX_InternalExtensions)) {
			return false
		}
	}

	if em1 := v1.FieldByName("XXX_extensions"); em1.IsValid() {
		em2 := v2.FieldByName("XXX_extensions")
		if !equalExtMap(v1.Type(), em1.Interface().(map[int32]Extension), em2.Interface().(map[int32]Extension)) {
			return false
		}
	}

	uf := v1.FieldByName("XXX_unrecognized")
	if !uf.IsValid() {
		return true
	}

	u1 := uf.Bytes()
	u2 := v2.FieldByName("XXX_unrecognized").Bytes()
	if !bytes.Equal(u1, u2) {
		return false
	}

	return true
}



func equalAny(v1, v2 reflect.Value, prop *Properties) bool {
	if v1.Type() == protoMessageType {
		m1, _ := v1.Interface().(Message)
		m2, _ := v2.Interface().(Message)
		return Equal(m1, m2)
	}
	switch v1.Kind() {
	case reflect.Bool:
		return v1.Bool() == v2.Bool()
	case reflect.Float32, reflect.Float64:
		return v1.Float() == v2.Float()
	case reflect.Int32, reflect.Int64:
		return v1.Int() == v2.Int()
	case reflect.Interface:
		
		n1, n2 := v1.IsNil(), v2.IsNil()
		if n1 || n2 {
			return n1 == n2
		}
		e1, e2 := v1.Elem(), v2.Elem()
		if e1.Type() != e2.Type() {
			return false
		}
		return equalAny(e1, e2, nil)
	case reflect.Map:
		if v1.Len() != v2.Len() {
			return false
		}
		for _, key := range v1.MapKeys() {
			val2 := v2.MapIndex(key)
			if !val2.IsValid() {
				
				return false
			}
			if !equalAny(v1.MapIndex(key), val2, nil) {
				return false
			}
		}
		return true
	case reflect.Ptr:
		
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		return equalAny(v1.Elem(), v2.Elem(), prop)
	case reflect.Slice:
		if v1.Type().Elem().Kind() == reflect.Uint8 {
			

			
			
			if prop != nil && prop.proto3 && v1.Len() == 0 && v2.Len() == 0 {
				return true
			}
			if v1.IsNil() != v2.IsNil() {
				return false
			}
			return bytes.Equal(v1.Interface().([]byte), v2.Interface().([]byte))
		}

		if v1.Len() != v2.Len() {
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if !equalAny(v1.Index(i), v2.Index(i), prop) {
				return false
			}
		}
		return true
	case reflect.String:
		return v1.Interface().(string) == v2.Interface().(string)
	case reflect.Struct:
		return equalStruct(v1, v2)
	case reflect.Uint32, reflect.Uint64:
		return v1.Uint() == v2.Uint()
	}

	
	log.Printf("proto: don't know how to compare %v", v1)
	return false
}



func equalExtensions(base reflect.Type, x1, x2 XXX_InternalExtensions) bool {
	em1, _ := x1.extensionsRead()
	em2, _ := x2.extensionsRead()
	return equalExtMap(base, em1, em2)
}

func equalExtMap(base reflect.Type, em1, em2 map[int32]Extension) bool {
	if len(em1) != len(em2) {
		return false
	}

	for extNum, e1 := range em1 {
		e2, ok := em2[extNum]
		if !ok {
			return false
		}

		m1, m2 := e1.value, e2.value

		if m1 != nil && m2 != nil {
			
			if !equalAny(reflect.ValueOf(m1), reflect.ValueOf(m2), nil) {
				return false
			}
			continue
		}

		
		
		var desc *ExtensionDesc
		if m := extensionMaps[base]; m != nil {
			desc = m[extNum]
		}
		if desc == nil {
			log.Printf("proto: don't know how to compare extension %d of %v", extNum, base)
			continue
		}
		var err error
		if m1 == nil {
			m1, err = decodeExtension(e1.enc, desc)
		}
		if m2 == nil && err == nil {
			m2, err = decodeExtension(e2.enc, desc)
		}
		if err != nil {
			
			log.Printf("proto: badly encoded extension %d of %v: %v", extNum, base, err)
			return false
		}
		if !equalAny(reflect.ValueOf(m1), reflect.ValueOf(m2), nil) {
			return false
		}
	}

	return true
}
