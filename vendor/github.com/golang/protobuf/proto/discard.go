






























package proto

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

type generatedDiscarder interface {
	XXX_DiscardUnknown()
}












func DiscardUnknown(m Message) {
	if m, ok := m.(generatedDiscarder); ok {
		m.XXX_DiscardUnknown()
		return
	}
	
	
	
	discardLegacy(m)
}


func (a *InternalMessageInfo) DiscardUnknown(m Message) {
	di := atomicLoadDiscardInfo(&a.discard)
	if di == nil {
		di = getDiscardInfo(reflect.TypeOf(m).Elem())
		atomicStoreDiscardInfo(&a.discard, di)
	}
	di.discard(toPointer(&m))
}

type discardInfo struct {
	typ reflect.Type

	initialized int32 
	lock        sync.Mutex

	fields       []discardFieldInfo
	unrecognized field
}

type discardFieldInfo struct {
	field   field 
	discard func(src pointer)
}

var (
	discardInfoMap  = map[reflect.Type]*discardInfo{}
	discardInfoLock sync.Mutex
)

func getDiscardInfo(t reflect.Type) *discardInfo {
	discardInfoLock.Lock()
	defer discardInfoLock.Unlock()
	di := discardInfoMap[t]
	if di == nil {
		di = &discardInfo{typ: t}
		discardInfoMap[t] = di
	}
	return di
}

func (di *discardInfo) discard(src pointer) {
	if src.isNil() {
		return 
	}

	if atomic.LoadInt32(&di.initialized) == 0 {
		di.computeDiscardInfo()
	}

	for _, fi := range di.fields {
		sfp := src.offset(fi.field)
		fi.discard(sfp)
	}

	
	
	if em, err := extendable(src.asPointerTo(di.typ).Interface()); err == nil {
		
		emm, _ := em.extensionsRead()
		for _, mx := range emm {
			if m, ok := mx.value.(Message); ok {
				DiscardUnknown(m)
			}
		}
	}

	if di.unrecognized.IsValid() {
		*src.offset(di.unrecognized).toBytes() = nil
	}
}

func (di *discardInfo) computeDiscardInfo() {
	di.lock.Lock()
	defer di.lock.Unlock()
	if di.initialized != 0 {
		return
	}
	t := di.typ
	n := t.NumField()

	for i := 0; i < n; i++ {
		f := t.Field(i)
		if strings.HasPrefix(f.Name, "XXX_") {
			continue
		}

		dfi := discardFieldInfo{field: toField(&f)}
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
			panic(fmt.Sprintf("%v.%s cannot be a slice of pointers to primitive types", t, f.Name))
		}

		switch tf.Kind() {
		case reflect.Struct:
			switch {
			case !isPointer:
				panic(fmt.Sprintf("%v.%s cannot be a direct struct value", t, f.Name))
			case isSlice: 
				di := getDiscardInfo(tf)
				dfi.discard = func(src pointer) {
					sps := src.getPointerSlice()
					for _, sp := range sps {
						if !sp.isNil() {
							di.discard(sp)
						}
					}
				}
			default: 
				di := getDiscardInfo(tf)
				dfi.discard = func(src pointer) {
					sp := src.getPointer()
					if !sp.isNil() {
						di.discard(sp)
					}
				}
			}
		case reflect.Map:
			switch {
			case isPointer || isSlice:
				panic(fmt.Sprintf("%v.%s cannot be a pointer to a map or a slice of map values", t, f.Name))
			default: 
				if tf.Elem().Kind() == reflect.Ptr { 
					dfi.discard = func(src pointer) {
						sm := src.asPointerTo(tf).Elem()
						if sm.Len() == 0 {
							return
						}
						for _, key := range sm.MapKeys() {
							val := sm.MapIndex(key)
							DiscardUnknown(val.Interface().(Message))
						}
					}
				} else {
					dfi.discard = func(pointer) {} 
				}
			}
		case reflect.Interface:
			
			switch {
			case isPointer || isSlice:
				panic(fmt.Sprintf("%v.%s cannot be a pointer to a interface or a slice of interface values", t, f.Name))
			default: 
				
				dfi.discard = func(src pointer) {
					su := src.asPointerTo(tf).Elem()
					if !su.IsNil() {
						sv := su.Elem().Elem().Field(0)
						if sv.Kind() == reflect.Ptr && sv.IsNil() {
							return
						}
						switch sv.Type().Kind() {
						case reflect.Ptr: 
							DiscardUnknown(sv.Interface().(Message))
						}
					}
				}
			}
		default:
			continue
		}
		di.fields = append(di.fields, dfi)
	}

	di.unrecognized = invalidField
	if f, ok := t.FieldByName("XXX_unrecognized"); ok {
		if f.Type != reflect.TypeOf([]byte{}) {
			panic("expected XXX_unrecognized to be of type []byte")
		}
		di.unrecognized = toField(&f)
	}

	atomic.StoreInt32(&di.initialized, 1)
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

	
	
	if em, err := extendable(m); err == nil {
		
		emm, _ := em.extensionsRead()
		for _, mx := range emm {
			if m, ok := mx.value.(Message); ok {
				discardLegacy(m)
			}
		}
	}
}
