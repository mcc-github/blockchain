

































package proto

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)


func Clone(src Message) Message {
	in := reflect.ValueOf(src)
	if in.IsNil() {
		return src
	}
	out := reflect.New(in.Type().Elem())
	dst := out.Interface().(Message)
	Merge(dst, src)
	return dst
}


type Merger interface {
	
	
	
	
	
	Merge(src Message)
}




type generatedMerger interface {
	XXX_Merge(src Message)
}





func Merge(dst, src Message) {
	if m, ok := dst.(Merger); ok {
		m.Merge(src)
		return
	}

	in := reflect.ValueOf(src)
	out := reflect.ValueOf(dst)
	if out.IsNil() {
		panic("proto: nil destination")
	}
	if in.Type() != out.Type() {
		panic(fmt.Sprintf("proto.Merge(%T, %T) type mismatch", dst, src))
	}
	if in.IsNil() {
		return 
	}
	if m, ok := dst.(generatedMerger); ok {
		m.XXX_Merge(src)
		return
	}
	mergeStruct(out.Elem(), in.Elem())
}

func mergeStruct(out, in reflect.Value) {
	sprop := GetProperties(in.Type())
	for i := 0; i < in.NumField(); i++ {
		f := in.Type().Field(i)
		if strings.HasPrefix(f.Name, "XXX_") {
			continue
		}
		mergeAny(out.Field(i), in.Field(i), false, sprop.Prop[i])
	}

	if emIn, ok := in.Addr().Interface().(extensionsBytes); ok {
		emOut := out.Addr().Interface().(extensionsBytes)
		bIn := emIn.GetExtensions()
		bOut := emOut.GetExtensions()
		*bOut = append(*bOut, *bIn...)
	} else if emIn, err := extendable(in.Addr().Interface()); err == nil {
		emOut, _ := extendable(out.Addr().Interface())
		mIn, muIn := emIn.extensionsRead()
		if mIn != nil {
			mOut := emOut.extensionsWrite()
			muIn.Lock()
			mergeExtension(mOut, mIn)
			muIn.Unlock()
		}
	}

	uf := in.FieldByName("XXX_unrecognized")
	if !uf.IsValid() {
		return
	}
	uin := uf.Bytes()
	if len(uin) > 0 {
		out.FieldByName("XXX_unrecognized").SetBytes(append([]byte(nil), uin...))
	}
}




func mergeAny(out, in reflect.Value, viaPtr bool, prop *Properties) {
	if in.Type() == protoMessageType {
		if !in.IsNil() {
			if out.IsNil() {
				out.Set(reflect.ValueOf(Clone(in.Interface().(Message))))
			} else {
				Merge(out.Interface().(Message), in.Interface().(Message))
			}
		}
		return
	}
	switch in.Kind() {
	case reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int32, reflect.Int64,
		reflect.String, reflect.Uint32, reflect.Uint64:
		if !viaPtr && isProto3Zero(in) {
			return
		}
		out.Set(in)
	case reflect.Interface:
		
		if in.IsNil() {
			return
		}
		
		
		if out.IsNil() || out.Elem().Type() != in.Elem().Type() {
			out.Set(reflect.New(in.Elem().Elem().Type())) 
		}
		mergeAny(out.Elem(), in.Elem(), false, nil)
	case reflect.Map:
		if in.Len() == 0 {
			return
		}
		if out.IsNil() {
			out.Set(reflect.MakeMap(in.Type()))
		}
		
		elemKind := in.Type().Elem().Kind()
		for _, key := range in.MapKeys() {
			var val reflect.Value
			switch elemKind {
			case reflect.Ptr:
				val = reflect.New(in.Type().Elem().Elem())
				mergeAny(val, in.MapIndex(key), false, nil)
			case reflect.Slice:
				val = in.MapIndex(key)
				val = reflect.ValueOf(append([]byte{}, val.Bytes()...))
			default:
				val = in.MapIndex(key)
			}
			out.SetMapIndex(key, val)
		}
	case reflect.Ptr:
		if in.IsNil() {
			return
		}
		if out.IsNil() {
			out.Set(reflect.New(in.Elem().Type()))
		}
		mergeAny(out.Elem(), in.Elem(), true, nil)
	case reflect.Slice:
		if in.IsNil() {
			return
		}
		if in.Type().Elem().Kind() == reflect.Uint8 {
			

			
			
			
			if prop != nil && prop.proto3 && in.Len() == 0 {
				return
			}

			
			
			
			out.SetBytes(append([]byte{}, in.Bytes()...))
			return
		}
		n := in.Len()
		if out.IsNil() {
			out.Set(reflect.MakeSlice(in.Type(), 0, n))
		}
		switch in.Type().Elem().Kind() {
		case reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int32, reflect.Int64,
			reflect.String, reflect.Uint32, reflect.Uint64:
			out.Set(reflect.AppendSlice(out, in))
		default:
			for i := 0; i < n; i++ {
				x := reflect.Indirect(reflect.New(in.Type().Elem()))
				mergeAny(x, in.Index(i), false, nil)
				out.Set(reflect.Append(out, x))
			}
		}
	case reflect.Struct:
		mergeStruct(out, in)
	default:
		
		log.Printf("proto: don't know how to copy %v", in)
	}
}

func mergeExtension(out, in map[int32]Extension) {
	for extNum, eIn := range in {
		eOut := Extension{desc: eIn.desc}
		if eIn.value != nil {
			v := reflect.New(reflect.TypeOf(eIn.value)).Elem()
			mergeAny(v, reflect.ValueOf(eIn.value), false, nil)
			eOut.value = v.Interface()
		}
		if eIn.enc != nil {
			eOut.enc = make([]byte, len(eIn.enc))
			copy(eOut.enc, eIn.enc)
		}

		out[extNum] = eOut
	}
}
