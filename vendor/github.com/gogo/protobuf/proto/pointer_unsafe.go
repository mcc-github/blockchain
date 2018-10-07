


































package proto

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

const unsafeAllowed = true



type field uintptr


func toField(f *reflect.StructField) field {
	return field(f.Offset)
}


const invalidField = ^field(0)


const zeroField = field(0)


func (f field) IsValid() bool {
	return f != invalidField
}





type pointer struct {
	p unsafe.Pointer
}


var ptrSize = unsafe.Sizeof(uintptr(0))



func toPointer(i *Message) pointer {
	
	
	
	return pointer{p: (*[2]unsafe.Pointer)(unsafe.Pointer(i))[1]}
}



func toAddrPointer(i *interface{}, isptr bool) pointer {
	
	if isptr {
		
		
		return pointer{p: unsafe.Pointer(uintptr(unsafe.Pointer(i)) + ptrSize)}
	}
	
	
	return pointer{p: (*[2]unsafe.Pointer)(unsafe.Pointer(i))[1]}
}


func valToPointer(v reflect.Value) pointer {
	return pointer{p: unsafe.Pointer(v.Pointer())}
}



func (p pointer) offset(f field) pointer {
	
	
	
	return pointer{p: unsafe.Pointer(uintptr(p.p) + uintptr(f))}
}

func (p pointer) isNil() bool {
	return p.p == nil
}

func (p pointer) toInt64() *int64 {
	return (*int64)(p.p)
}
func (p pointer) toInt64Ptr() **int64 {
	return (**int64)(p.p)
}
func (p pointer) toInt64Slice() *[]int64 {
	return (*[]int64)(p.p)
}
func (p pointer) toInt32() *int32 {
	return (*int32)(p.p)
}



func (p pointer) getInt32Ptr() *int32 {
	return *(**int32)(p.p)
}
func (p pointer) setInt32Ptr(v int32) {
	*(**int32)(p.p) = &v
}




func (p pointer) getInt32Slice() []int32 {
	return *(*[]int32)(p.p)
}




func (p pointer) setInt32Slice(v []int32) {
	*(*[]int32)(p.p) = v
}


func (p pointer) appendInt32Slice(v int32) {
	s := (*[]int32)(p.p)
	*s = append(*s, v)
}

func (p pointer) toUint64() *uint64 {
	return (*uint64)(p.p)
}
func (p pointer) toUint64Ptr() **uint64 {
	return (**uint64)(p.p)
}
func (p pointer) toUint64Slice() *[]uint64 {
	return (*[]uint64)(p.p)
}
func (p pointer) toUint32() *uint32 {
	return (*uint32)(p.p)
}
func (p pointer) toUint32Ptr() **uint32 {
	return (**uint32)(p.p)
}
func (p pointer) toUint32Slice() *[]uint32 {
	return (*[]uint32)(p.p)
}
func (p pointer) toBool() *bool {
	return (*bool)(p.p)
}
func (p pointer) toBoolPtr() **bool {
	return (**bool)(p.p)
}
func (p pointer) toBoolSlice() *[]bool {
	return (*[]bool)(p.p)
}
func (p pointer) toFloat64() *float64 {
	return (*float64)(p.p)
}
func (p pointer) toFloat64Ptr() **float64 {
	return (**float64)(p.p)
}
func (p pointer) toFloat64Slice() *[]float64 {
	return (*[]float64)(p.p)
}
func (p pointer) toFloat32() *float32 {
	return (*float32)(p.p)
}
func (p pointer) toFloat32Ptr() **float32 {
	return (**float32)(p.p)
}
func (p pointer) toFloat32Slice() *[]float32 {
	return (*[]float32)(p.p)
}
func (p pointer) toString() *string {
	return (*string)(p.p)
}
func (p pointer) toStringPtr() **string {
	return (**string)(p.p)
}
func (p pointer) toStringSlice() *[]string {
	return (*[]string)(p.p)
}
func (p pointer) toBytes() *[]byte {
	return (*[]byte)(p.p)
}
func (p pointer) toBytesSlice() *[][]byte {
	return (*[][]byte)(p.p)
}
func (p pointer) toExtensions() *XXX_InternalExtensions {
	return (*XXX_InternalExtensions)(p.p)
}
func (p pointer) toOldExtensions() *map[int32]Extension {
	return (*map[int32]Extension)(p.p)
}




func (p pointer) getPointerSlice() []pointer {
	
	
	return *(*[]pointer)(p.p)
}




func (p pointer) setPointerSlice(v []pointer) {
	
	
	*(*[]pointer)(p.p) = v
}


func (p pointer) getPointer() pointer {
	return pointer{p: *(*unsafe.Pointer)(p.p)}
}


func (p pointer) setPointer(q pointer) {
	*(*unsafe.Pointer)(p.p) = q.p
}


func (p pointer) appendPointer(q pointer) {
	s := (*[]unsafe.Pointer)(p.p)
	*s = append(*s, q.p)
}



func (p pointer) getInterfacePointer() pointer {
	
	return pointer{p: (*(*[2]unsafe.Pointer)(p.p))[1]}
}



func (p pointer) asPointerTo(t reflect.Type) reflect.Value {
	return reflect.NewAt(t, p.p)
}

func atomicLoadUnmarshalInfo(p **unmarshalInfo) *unmarshalInfo {
	return (*unmarshalInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))
}
func atomicStoreUnmarshalInfo(p **unmarshalInfo, v *unmarshalInfo) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(v))
}
func atomicLoadMarshalInfo(p **marshalInfo) *marshalInfo {
	return (*marshalInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))
}
func atomicStoreMarshalInfo(p **marshalInfo, v *marshalInfo) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(v))
}
func atomicLoadMergeInfo(p **mergeInfo) *mergeInfo {
	return (*mergeInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))
}
func atomicStoreMergeInfo(p **mergeInfo, v *mergeInfo) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(v))
}
func atomicLoadDiscardInfo(p **discardInfo) *discardInfo {
	return (*discardInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))
}
func atomicStoreDiscardInfo(p **discardInfo, v *discardInfo) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(v))
}
