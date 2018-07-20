




































package proto

import (
	"reflect"
	"sync"
)

const unsafeAllowed = false




type field []int


func toField(f *reflect.StructField) field {
	return f.Index
}


var invalidField = field(nil)


var zeroField = field([]int{})


func (f field) IsValid() bool { return f != nil }





type pointer struct {
	v reflect.Value
}



func toPointer(i *Message) pointer {
	return pointer{v: reflect.ValueOf(*i)}
}



func toAddrPointer(i *interface{}, isptr bool) pointer {
	v := reflect.ValueOf(*i)
	u := reflect.New(v.Type())
	u.Elem().Set(v)
	return pointer{v: u}
}


func valToPointer(v reflect.Value) pointer {
	return pointer{v: v}
}



func (p pointer) offset(f field) pointer {
	return pointer{v: p.v.Elem().FieldByIndex(f).Addr()}
}

func (p pointer) isNil() bool {
	return p.v.IsNil()
}




func grow(s reflect.Value) reflect.Value {
	n, m := s.Len(), s.Cap()
	if n < m {
		s.SetLen(n + 1)
	} else {
		s.Set(reflect.Append(s, reflect.Zero(s.Type().Elem())))
	}
	return s.Index(n)
}

func (p pointer) toInt64() *int64 {
	return p.v.Interface().(*int64)
}
func (p pointer) toInt64Ptr() **int64 {
	return p.v.Interface().(**int64)
}
func (p pointer) toInt64Slice() *[]int64 {
	return p.v.Interface().(*[]int64)
}

var int32ptr = reflect.TypeOf((*int32)(nil))

func (p pointer) toInt32() *int32 {
	return p.v.Convert(int32ptr).Interface().(*int32)
}




func (p pointer) getInt32Ptr() *int32 {
	if p.v.Type().Elem().Elem() == reflect.TypeOf(int32(0)) {
		
		return p.v.Elem().Interface().(*int32)
	}
	
	return p.v.Elem().Convert(int32PtrType).Interface().(*int32)
}
func (p pointer) setInt32Ptr(v int32) {
	
	
	
	
	p.v.Elem().Set(reflect.ValueOf(&v).Convert(p.v.Type().Elem()))
}



func (p pointer) getInt32Slice() []int32 {
	if p.v.Type().Elem().Elem() == reflect.TypeOf(int32(0)) {
		
		return p.v.Elem().Interface().([]int32)
	}
	
	
	
	slice := p.v.Elem()
	s := make([]int32, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		s[i] = int32(slice.Index(i).Int())
	}
	return s
}



func (p pointer) setInt32Slice(v []int32) {
	if p.v.Type().Elem().Elem() == reflect.TypeOf(int32(0)) {
		
		p.v.Elem().Set(reflect.ValueOf(v))
		return
	}
	
	
	
	slice := reflect.MakeSlice(p.v.Type().Elem(), len(v), cap(v))
	for i, x := range v {
		slice.Index(i).SetInt(int64(x))
	}
	p.v.Elem().Set(slice)
}
func (p pointer) appendInt32Slice(v int32) {
	grow(p.v.Elem()).SetInt(int64(v))
}

func (p pointer) toUint64() *uint64 {
	return p.v.Interface().(*uint64)
}
func (p pointer) toUint64Ptr() **uint64 {
	return p.v.Interface().(**uint64)
}
func (p pointer) toUint64Slice() *[]uint64 {
	return p.v.Interface().(*[]uint64)
}
func (p pointer) toUint32() *uint32 {
	return p.v.Interface().(*uint32)
}
func (p pointer) toUint32Ptr() **uint32 {
	return p.v.Interface().(**uint32)
}
func (p pointer) toUint32Slice() *[]uint32 {
	return p.v.Interface().(*[]uint32)
}
func (p pointer) toBool() *bool {
	return p.v.Interface().(*bool)
}
func (p pointer) toBoolPtr() **bool {
	return p.v.Interface().(**bool)
}
func (p pointer) toBoolSlice() *[]bool {
	return p.v.Interface().(*[]bool)
}
func (p pointer) toFloat64() *float64 {
	return p.v.Interface().(*float64)
}
func (p pointer) toFloat64Ptr() **float64 {
	return p.v.Interface().(**float64)
}
func (p pointer) toFloat64Slice() *[]float64 {
	return p.v.Interface().(*[]float64)
}
func (p pointer) toFloat32() *float32 {
	return p.v.Interface().(*float32)
}
func (p pointer) toFloat32Ptr() **float32 {
	return p.v.Interface().(**float32)
}
func (p pointer) toFloat32Slice() *[]float32 {
	return p.v.Interface().(*[]float32)
}
func (p pointer) toString() *string {
	return p.v.Interface().(*string)
}
func (p pointer) toStringPtr() **string {
	return p.v.Interface().(**string)
}
func (p pointer) toStringSlice() *[]string {
	return p.v.Interface().(*[]string)
}
func (p pointer) toBytes() *[]byte {
	return p.v.Interface().(*[]byte)
}
func (p pointer) toBytesSlice() *[][]byte {
	return p.v.Interface().(*[][]byte)
}
func (p pointer) toExtensions() *XXX_InternalExtensions {
	return p.v.Interface().(*XXX_InternalExtensions)
}
func (p pointer) toOldExtensions() *map[int32]Extension {
	return p.v.Interface().(*map[int32]Extension)
}
func (p pointer) getPointer() pointer {
	return pointer{v: p.v.Elem()}
}
func (p pointer) setPointer(q pointer) {
	p.v.Elem().Set(q.v)
}
func (p pointer) appendPointer(q pointer) {
	grow(p.v.Elem()).Set(q.v)
}



func (p pointer) getPointerSlice() []pointer {
	if p.v.IsNil() {
		return nil
	}
	n := p.v.Elem().Len()
	s := make([]pointer, n)
	for i := 0; i < n; i++ {
		s[i] = pointer{v: p.v.Elem().Index(i)}
	}
	return s
}



func (p pointer) setPointerSlice(v []pointer) {
	if v == nil {
		p.v.Elem().Set(reflect.New(p.v.Elem().Type()).Elem())
		return
	}
	s := reflect.MakeSlice(p.v.Elem().Type(), 0, len(v))
	for _, p := range v {
		s = reflect.Append(s, p.v)
	}
	p.v.Elem().Set(s)
}



func (p pointer) getInterfacePointer() pointer {
	if p.v.Elem().IsNil() {
		return pointer{v: p.v.Elem()}
	}
	return pointer{v: p.v.Elem().Elem().Elem().Field(0).Addr()} 
}

func (p pointer) asPointerTo(t reflect.Type) reflect.Value {
	
	return p.v
}

func atomicLoadUnmarshalInfo(p **unmarshalInfo) *unmarshalInfo {
	atomicLock.Lock()
	defer atomicLock.Unlock()
	return *p
}
func atomicStoreUnmarshalInfo(p **unmarshalInfo, v *unmarshalInfo) {
	atomicLock.Lock()
	defer atomicLock.Unlock()
	*p = v
}
func atomicLoadMarshalInfo(p **marshalInfo) *marshalInfo {
	atomicLock.Lock()
	defer atomicLock.Unlock()
	return *p
}
func atomicStoreMarshalInfo(p **marshalInfo, v *marshalInfo) {
	atomicLock.Lock()
	defer atomicLock.Unlock()
	*p = v
}
func atomicLoadMergeInfo(p **mergeInfo) *mergeInfo {
	atomicLock.Lock()
	defer atomicLock.Unlock()
	return *p
}
func atomicStoreMergeInfo(p **mergeInfo, v *mergeInfo) {
	atomicLock.Lock()
	defer atomicLock.Unlock()
	*p = v
}
func atomicLoadDiscardInfo(p **discardInfo) *discardInfo {
	atomicLock.Lock()
	defer atomicLock.Unlock()
	return *p
}
func atomicStoreDiscardInfo(p **discardInfo, v *discardInfo) {
	atomicLock.Lock()
	defer atomicLock.Unlock()
	*p = v
}

var atomicLock sync.Mutex
