






























package proto

import (
	"fmt"
	"reflect"
	"strings"
)












func DiscardUnknown(m Message) {
	discardLegacy(m)
}

func discardLegacy(m Message) {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if strings.HasPrefix(f.Name, "XXX_") {
			continue
		}
		vf := v.Field(i)
		tf := f.Type

		
		var isPointer, isSlice bool
		if tf.Kind() == reflect.Slice && tf.Elem().Kind() != reflect.Uint8 {
			isSlice = true
			tf = tf.Elem()
		}
		if tf.Kind() == reflect.Ptr {
			isPointer = true
			tf = tf.Elem()
		}
		if isPointer && isSlice && tf.Kind() != reflect.Struct {
			panic(fmt.Sprintf("%T.%s cannot be a slice of pointers to primitive types", m, f.Name))
		}

		switch tf.Kind() {
		case reflect.Struct:
			switch {
			case !isPointer:
				panic(fmt.Sprintf("%T.%s cannot be a direct struct value", m, f.Name))
			case isSlice: 
				for j := 0; j < vf.Len(); j++ {
					discardLegacy(vf.Index(j).Interface().(Message))
				}
			default: 
				discardLegacy(vf.Interface().(Message))
			}
		case reflect.Map:
			switch {
			case isPointer || isSlice:
				panic(fmt.Sprintf("%T.%s cannot be a pointer to a map or a slice of map values", m, f.Name))
			default: 
				tv := vf.Type().Elem()
				if tv.Kind() == reflect.Ptr && tv.Implements(protoMessageType) { 
					for _, key := range vf.MapKeys() {
						val := vf.MapIndex(key)
						discardLegacy(val.Interface().(Message))
					}
				}
			}
		case reflect.Interface:
			
			switch {
			case isPointer || isSlice:
				panic(fmt.Sprintf("%T.%s cannot be a pointer to a interface or a slice of interface values", m, f.Name))
			default: 
				if !vf.IsNil() && f.Tag.Get("protobuf_oneof") != "" {
					vf = vf.Elem() 
					if !vf.IsNil() {
						vf = vf.Elem()   
						vf = vf.Field(0) 
						if vf.Kind() == reflect.Ptr {
							discardLegacy(vf.Interface().(Message))
						}
					}
				}
			}
		}
	}

	if vf := v.FieldByName("XXX_unrecognized"); vf.IsValid() {
		if vf.Type() != reflect.TypeOf([]byte{}) {
			panic("expected XXX_unrecognized to be of type []byte")
		}
		vf.Set(reflect.ValueOf([]byte(nil)))
	}

	
	
	if em, ok := extendable(m); ok {
		
		emm, _ := em.extensionsRead()
		for _, mx := range emm {
			if m, ok := mx.value.(Message); ok {
				discardLegacy(m)
			}
		}
	}
}
