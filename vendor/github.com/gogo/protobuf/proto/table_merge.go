






























package proto

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)



func (a *InternalMessageInfo) Merge(dst, src Message) {
	mi := atomicLoadMergeInfo(&a.merge)
	if mi == nil {
		mi = getMergeInfo(reflect.TypeOf(dst).Elem())
		atomicStoreMergeInfo(&a.merge, mi)
	}
	mi.merge(toPointer(&dst), toPointer(&src))
}

type mergeInfo struct {
	typ reflect.Type

	initialized int32 
	lock        sync.Mutex

	fields       []mergeFieldInfo
	unrecognized field 
}

type mergeFieldInfo struct {
	field field 

	
	
	
	
	
	
	isPointer bool

	
	
	
	
	
	
	
	basicWidth int

	
	merge func(dst, src pointer)
}

var (
	mergeInfoMap  = map[reflect.Type]*mergeInfo{}
	mergeInfoLock sync.Mutex
)

func getMergeInfo(t reflect.Type) *mergeInfo {
	mergeInfoLock.Lock()
	defer mergeInfoLock.Unlock()
	mi := mergeInfoMap[t]
	if mi == nil {
		mi = &mergeInfo{typ: t}
		mergeInfoMap[t] = mi
	}
	return mi
}


func (mi *mergeInfo) merge(dst, src pointer) {
	if dst.isNil() {
		panic("proto: nil destination")
	}
	if src.isNil() {
		return 
	}

	if atomic.LoadInt32(&mi.initialized) == 0 {
		mi.computeMergeInfo()
	}

	for _, fi := range mi.fields {
		sfp := src.offset(fi.field)

		
		
		
		if unsafeAllowed {
			if fi.isPointer && sfp.getPointer().isNil() { 
				continue
			}
			if fi.basicWidth > 0 {
				switch {
				case fi.basicWidth == 1 && !*sfp.toBool():
					continue
				case fi.basicWidth == 4 && *sfp.toUint32() == 0:
					continue
				case fi.basicWidth == 8 && *sfp.toUint64() == 0:
					continue
				}
			}
		}

		dfp := dst.offset(fi.field)
		fi.merge(dfp, sfp)
	}

	
	out := dst.asPointerTo(mi.typ).Elem()
	in := src.asPointerTo(mi.typ).Elem()
	if emIn, err := extendable(in.Addr().Interface()); err == nil {
		emOut, _ := extendable(out.Addr().Interface())
		mIn, muIn := emIn.extensionsRead()
		if mIn != nil {
			mOut := emOut.extensionsWrite()
			muIn.Lock()
			mergeExtension(mOut, mIn)
			muIn.Unlock()
		}
	}

	if mi.unrecognized.IsValid() {
		if b := *src.offset(mi.unrecognized).toBytes(); len(b) > 0 {
			*dst.offset(mi.unrecognized).toBytes() = append([]byte(nil), b...)
		}
	}
}

func (mi *mergeInfo) computeMergeInfo() {
	mi.lock.Lock()
	defer mi.lock.Unlock()
	if mi.initialized != 0 {
		return
	}
	t := mi.typ
	n := t.NumField()

	props := GetProperties(t)
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if strings.HasPrefix(f.Name, "XXX_") {
			continue
		}

		mfi := mergeFieldInfo{field: toField(&f)}
		tf := f.Type

		
		
		
		if unsafeAllowed {
			switch tf.Kind() {
			case reflect.Ptr, reflect.Slice, reflect.String:
				
				
				
				mfi.isPointer = true
			case reflect.Bool:
				mfi.basicWidth = 1
			case reflect.Int32, reflect.Uint32, reflect.Float32:
				mfi.basicWidth = 4
			case reflect.Int64, reflect.Uint64, reflect.Float64:
				mfi.basicWidth = 8
			}
		}

		
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
			panic("both pointer and slice for basic type in " + tf.Name())
		}

		switch tf.Kind() {
		case reflect.Int32:
			switch {
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					
					
					sfs := src.getInt32Slice()
					if sfs != nil {
						dfs := dst.getInt32Slice()
						dfs = append(dfs, sfs...)
						if dfs == nil {
							dfs = []int32{}
						}
						dst.setInt32Slice(dfs)
					}
				}
			case isPointer: 
				mfi.merge = func(dst, src pointer) {
					
					
					sfp := src.getInt32Ptr()
					if sfp != nil {
						dfp := dst.getInt32Ptr()
						if dfp == nil {
							dst.setInt32Ptr(*sfp)
						} else {
							*dfp = *sfp
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					if v := *src.toInt32(); v != 0 {
						*dst.toInt32() = v
					}
				}
			}
		case reflect.Int64:
			switch {
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					sfsp := src.toInt64Slice()
					if *sfsp != nil {
						dfsp := dst.toInt64Slice()
						*dfsp = append(*dfsp, *sfsp...)
						if *dfsp == nil {
							*dfsp = []int64{}
						}
					}
				}
			case isPointer: 
				mfi.merge = func(dst, src pointer) {
					sfpp := src.toInt64Ptr()
					if *sfpp != nil {
						dfpp := dst.toInt64Ptr()
						if *dfpp == nil {
							*dfpp = Int64(**sfpp)
						} else {
							**dfpp = **sfpp
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					if v := *src.toInt64(); v != 0 {
						*dst.toInt64() = v
					}
				}
			}
		case reflect.Uint32:
			switch {
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					sfsp := src.toUint32Slice()
					if *sfsp != nil {
						dfsp := dst.toUint32Slice()
						*dfsp = append(*dfsp, *sfsp...)
						if *dfsp == nil {
							*dfsp = []uint32{}
						}
					}
				}
			case isPointer: 
				mfi.merge = func(dst, src pointer) {
					sfpp := src.toUint32Ptr()
					if *sfpp != nil {
						dfpp := dst.toUint32Ptr()
						if *dfpp == nil {
							*dfpp = Uint32(**sfpp)
						} else {
							**dfpp = **sfpp
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					if v := *src.toUint32(); v != 0 {
						*dst.toUint32() = v
					}
				}
			}
		case reflect.Uint64:
			switch {
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					sfsp := src.toUint64Slice()
					if *sfsp != nil {
						dfsp := dst.toUint64Slice()
						*dfsp = append(*dfsp, *sfsp...)
						if *dfsp == nil {
							*dfsp = []uint64{}
						}
					}
				}
			case isPointer: 
				mfi.merge = func(dst, src pointer) {
					sfpp := src.toUint64Ptr()
					if *sfpp != nil {
						dfpp := dst.toUint64Ptr()
						if *dfpp == nil {
							*dfpp = Uint64(**sfpp)
						} else {
							**dfpp = **sfpp
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					if v := *src.toUint64(); v != 0 {
						*dst.toUint64() = v
					}
				}
			}
		case reflect.Float32:
			switch {
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					sfsp := src.toFloat32Slice()
					if *sfsp != nil {
						dfsp := dst.toFloat32Slice()
						*dfsp = append(*dfsp, *sfsp...)
						if *dfsp == nil {
							*dfsp = []float32{}
						}
					}
				}
			case isPointer: 
				mfi.merge = func(dst, src pointer) {
					sfpp := src.toFloat32Ptr()
					if *sfpp != nil {
						dfpp := dst.toFloat32Ptr()
						if *dfpp == nil {
							*dfpp = Float32(**sfpp)
						} else {
							**dfpp = **sfpp
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					if v := *src.toFloat32(); v != 0 {
						*dst.toFloat32() = v
					}
				}
			}
		case reflect.Float64:
			switch {
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					sfsp := src.toFloat64Slice()
					if *sfsp != nil {
						dfsp := dst.toFloat64Slice()
						*dfsp = append(*dfsp, *sfsp...)
						if *dfsp == nil {
							*dfsp = []float64{}
						}
					}
				}
			case isPointer: 
				mfi.merge = func(dst, src pointer) {
					sfpp := src.toFloat64Ptr()
					if *sfpp != nil {
						dfpp := dst.toFloat64Ptr()
						if *dfpp == nil {
							*dfpp = Float64(**sfpp)
						} else {
							**dfpp = **sfpp
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					if v := *src.toFloat64(); v != 0 {
						*dst.toFloat64() = v
					}
				}
			}
		case reflect.Bool:
			switch {
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					sfsp := src.toBoolSlice()
					if *sfsp != nil {
						dfsp := dst.toBoolSlice()
						*dfsp = append(*dfsp, *sfsp...)
						if *dfsp == nil {
							*dfsp = []bool{}
						}
					}
				}
			case isPointer: 
				mfi.merge = func(dst, src pointer) {
					sfpp := src.toBoolPtr()
					if *sfpp != nil {
						dfpp := dst.toBoolPtr()
						if *dfpp == nil {
							*dfpp = Bool(**sfpp)
						} else {
							**dfpp = **sfpp
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					if v := *src.toBool(); v {
						*dst.toBool() = v
					}
				}
			}
		case reflect.String:
			switch {
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					sfsp := src.toStringSlice()
					if *sfsp != nil {
						dfsp := dst.toStringSlice()
						*dfsp = append(*dfsp, *sfsp...)
						if *dfsp == nil {
							*dfsp = []string{}
						}
					}
				}
			case isPointer: 
				mfi.merge = func(dst, src pointer) {
					sfpp := src.toStringPtr()
					if *sfpp != nil {
						dfpp := dst.toStringPtr()
						if *dfpp == nil {
							*dfpp = String(**sfpp)
						} else {
							**dfpp = **sfpp
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					if v := *src.toString(); v != "" {
						*dst.toString() = v
					}
				}
			}
		case reflect.Slice:
			isProto3 := props.Prop[i].proto3
			switch {
			case isPointer:
				panic("bad pointer in byte slice case in " + tf.Name())
			case tf.Elem().Kind() != reflect.Uint8:
				panic("bad element kind in byte slice case in " + tf.Name())
			case isSlice: 
				mfi.merge = func(dst, src pointer) {
					sbsp := src.toBytesSlice()
					if *sbsp != nil {
						dbsp := dst.toBytesSlice()
						for _, sb := range *sbsp {
							if sb == nil {
								*dbsp = append(*dbsp, nil)
							} else {
								*dbsp = append(*dbsp, append([]byte{}, sb...))
							}
						}
						if *dbsp == nil {
							*dbsp = [][]byte{}
						}
					}
				}
			default: 
				mfi.merge = func(dst, src pointer) {
					sbp := src.toBytes()
					if *sbp != nil {
						dbp := dst.toBytes()
						if !isProto3 || len(*sbp) > 0 {
							*dbp = append([]byte{}, *sbp...)
						}
					}
				}
			}
		case reflect.Struct:
			switch {
			case !isPointer:
				mergeInfo := getMergeInfo(tf)
				mfi.merge = func(dst, src pointer) {
					mergeInfo.merge(dst, src)
				}
			case isSlice: 
				mergeInfo := getMergeInfo(tf)
				mfi.merge = func(dst, src pointer) {
					sps := src.getPointerSlice()
					if sps != nil {
						dps := dst.getPointerSlice()
						for _, sp := range sps {
							var dp pointer
							if !sp.isNil() {
								dp = valToPointer(reflect.New(tf))
								mergeInfo.merge(dp, sp)
							}
							dps = append(dps, dp)
						}
						if dps == nil {
							dps = []pointer{}
						}
						dst.setPointerSlice(dps)
					}
				}
			default: 
				mergeInfo := getMergeInfo(tf)
				mfi.merge = func(dst, src pointer) {
					sp := src.getPointer()
					if !sp.isNil() {
						dp := dst.getPointer()
						if dp.isNil() {
							dp = valToPointer(reflect.New(tf))
							dst.setPointer(dp)
						}
						mergeInfo.merge(dp, sp)
					}
				}
			}
		case reflect.Map:
			switch {
			case isPointer || isSlice:
				panic("bad pointer or slice in map case in " + tf.Name())
			default: 
				mfi.merge = func(dst, src pointer) {
					sm := src.asPointerTo(tf).Elem()
					if sm.Len() == 0 {
						return
					}
					dm := dst.asPointerTo(tf).Elem()
					if dm.IsNil() {
						dm.Set(reflect.MakeMap(tf))
					}

					switch tf.Elem().Kind() {
					case reflect.Ptr: 
						for _, key := range sm.MapKeys() {
							val := sm.MapIndex(key)
							val = reflect.ValueOf(Clone(val.Interface().(Message)))
							dm.SetMapIndex(key, val)
						}
					case reflect.Slice: 
						for _, key := range sm.MapKeys() {
							val := sm.MapIndex(key)
							val = reflect.ValueOf(append([]byte{}, val.Bytes()...))
							dm.SetMapIndex(key, val)
						}
					default: 
						for _, key := range sm.MapKeys() {
							val := sm.MapIndex(key)
							dm.SetMapIndex(key, val)
						}
					}
				}
			}
		case reflect.Interface:
			
			switch {
			case isPointer || isSlice:
				panic("bad pointer or slice in interface case in " + tf.Name())
			default: 
				
				mfi.merge = func(dst, src pointer) {
					su := src.asPointerTo(tf).Elem()
					if !su.IsNil() {
						du := dst.asPointerTo(tf).Elem()
						typ := su.Elem().Type()
						if du.IsNil() || du.Elem().Type() != typ {
							du.Set(reflect.New(typ.Elem())) 
						}
						sv := su.Elem().Elem().Field(0)
						if sv.Kind() == reflect.Ptr && sv.IsNil() {
							return
						}
						dv := du.Elem().Elem().Field(0)
						if dv.Kind() == reflect.Ptr && dv.IsNil() {
							dv.Set(reflect.New(sv.Type().Elem())) 
						}
						switch sv.Type().Kind() {
						case reflect.Ptr: 
							Merge(dv.Interface().(Message), sv.Interface().(Message))
						case reflect.Slice: 
							dv.Set(reflect.ValueOf(append([]byte{}, sv.Bytes()...)))
						default: 
							dv.Set(sv)
						}
					}
				}
			}
		default:
			panic(fmt.Sprintf("merger not found for type:%s", tf))
		}
		mi.fields = append(mi.fields, mfi)
	}

	mi.unrecognized = invalidField
	if f, ok := t.FieldByName("XXX_unrecognized"); ok {
		if f.Type != reflect.TypeOf([]byte{}) {
			panic("expected XXX_unrecognized to be of type []byte")
		}
		mi.unrecognized = toField(&f)
	}

	atomic.StoreInt32(&mi.initialized, 1)
}
