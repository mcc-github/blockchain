


































package proto

import (
	"reflect"
	"unsafe"
)











type structPointer unsafe.Pointer


func toStructPointer(v reflect.Value) structPointer {
	return structPointer(unsafe.Pointer(v.Pointer()))
}


func structPointer_IsNil(p structPointer) bool {
	return p == nil
}



func structPointer_Interface(p structPointer, t reflect.Type) interface{} {
	return reflect.NewAt(t, unsafe.Pointer(p)).Interface()
}



type field uintptr


func toField(f *reflect.StructField) field {
	return field(f.Offset)
}


const invalidField = ^field(0)


func (f field) IsValid() bool {
	return f != ^field(0)
}


func structPointer_Bytes(p structPointer, f field) *[]byte {
	return (*[]byte)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_BytesSlice(p structPointer, f field) *[][]byte {
	return (*[][]byte)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_Bool(p structPointer, f field) **bool {
	return (**bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_BoolVal(p structPointer, f field) *bool {
	return (*bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_BoolSlice(p structPointer, f field) *[]bool {
	return (*[]bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_String(p structPointer, f field) **string {
	return (**string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_StringVal(p structPointer, f field) *string {
	return (*string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_StringSlice(p structPointer, f field) *[]string {
	return (*[]string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_Extensions(p structPointer, f field) *XXX_InternalExtensions {
	return (*XXX_InternalExtensions)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

func structPointer_ExtMap(p structPointer, f field) *map[int32]Extension {
	return (*map[int32]Extension)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_NewAt(p structPointer, f field, typ reflect.Type) reflect.Value {
	return reflect.NewAt(typ, unsafe.Pointer(uintptr(p)+uintptr(f)))
}


func structPointer_SetStructPointer(p structPointer, f field, q structPointer) {
	*(*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f))) = q
}


func structPointer_GetStructPointer(p structPointer, f field) structPointer {
	return *(*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


func structPointer_StructPointerSlice(p structPointer, f field) *structPointerSlice {
	return (*structPointerSlice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


type structPointerSlice []structPointer

func (v *structPointerSlice) Len() int                  { return len(*v) }
func (v *structPointerSlice) Index(i int) structPointer { return (*v)[i] }
func (v *structPointerSlice) Append(p structPointer)    { *v = append(*v, p) }


type word32 **uint32


func word32_IsNil(p word32) bool {
	return *p == nil
}


func word32_Set(p word32, o *Buffer, x uint32) {
	if len(o.uint32s) == 0 {
		o.uint32s = make([]uint32, uint32PoolSize)
	}
	o.uint32s[0] = x
	*p = &o.uint32s[0]
	o.uint32s = o.uint32s[1:]
}


func word32_Get(p word32) uint32 {
	return **p
}


func structPointer_Word32(p structPointer, f field) word32 {
	return word32((**uint32)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}


type word32Val *uint32


func word32Val_Set(p word32Val, x uint32) {
	*p = x
}


func word32Val_Get(p word32Val) uint32 {
	return *p
}


func structPointer_Word32Val(p structPointer, f field) word32Val {
	return word32Val((*uint32)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}


type word32Slice []uint32

func (v *word32Slice) Append(x uint32)    { *v = append(*v, x) }
func (v *word32Slice) Len() int           { return len(*v) }
func (v *word32Slice) Index(i int) uint32 { return (*v)[i] }


func structPointer_Word32Slice(p structPointer, f field) *word32Slice {
	return (*word32Slice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}


type word64 **uint64

func word64_Set(p word64, o *Buffer, x uint64) {
	if len(o.uint64s) == 0 {
		o.uint64s = make([]uint64, uint64PoolSize)
	}
	o.uint64s[0] = x
	*p = &o.uint64s[0]
	o.uint64s = o.uint64s[1:]
}

func word64_IsNil(p word64) bool {
	return *p == nil
}

func word64_Get(p word64) uint64 {
	return **p
}

func structPointer_Word64(p structPointer, f field) word64 {
	return word64((**uint64)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}


type word64Val *uint64

func word64Val_Set(p word64Val, o *Buffer, x uint64) {
	*p = x
}

func word64Val_Get(p word64Val) uint64 {
	return *p
}

func structPointer_Word64Val(p structPointer, f field) word64Val {
	return word64Val((*uint64)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}


type word64Slice []uint64

func (v *word64Slice) Append(x uint64)    { *v = append(*v, x) }
func (v *word64Slice) Len() int           { return len(*v) }
func (v *word64Slice) Index(i int) uint64 { return (*v)[i] }

func structPointer_Word64Slice(p structPointer, f field) *word64Slice {
	return (*word64Slice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}
