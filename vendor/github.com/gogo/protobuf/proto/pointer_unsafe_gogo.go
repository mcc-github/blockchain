































package proto

import (
	"reflect"
	"unsafe"
)

func structPointer_InterfaceAt(p structPointer, f field, t reflect.Type) interface{} {
	point := unsafe.Pointer(uintptr(p) + uintptr(f))
	r := reflect.NewAt(t, point)
	return r.Interface()
}

func structPointer_InterfaceRef(p structPointer, f field, t reflect.Type) interface{} {
	point := unsafe.Pointer(uintptr(p) + uintptr(f))
	r := reflect.NewAt(t, point)
	if r.Elem().IsNil() {
		return nil
	}
	return r.Elem().Interface()
}

func copyUintPtr(oldptr, newptr uintptr, size int) {
	oldbytes := make([]byte, 0)
	oldslice := (*reflect.SliceHeader)(unsafe.Pointer(&oldbytes))
	oldslice.Data = oldptr
	oldslice.Len = size
	oldslice.Cap = size
	newbytes := make([]byte, 0)
	newslice := (*reflect.SliceHeader)(unsafe.Pointer(&newbytes))
	newslice.Data = newptr
	newslice.Len = size
	newslice.Cap = size
	copy(newbytes, oldbytes)
}

func structPointer_Copy(oldptr structPointer, newptr structPointer, size int) {
	copyUintPtr(uintptr(oldptr), uintptr(newptr), size)
}

func appendStructPointer(base structPointer, f field, typ reflect.Type) structPointer {
	size := typ.Elem().Size()

	oldHeader := structPointer_GetSliceHeader(base, f)
	oldSlice := reflect.NewAt(typ, unsafe.Pointer(oldHeader)).Elem()
	newLen := oldHeader.Len + 1
	newSlice := reflect.MakeSlice(typ, newLen, newLen)
	reflect.Copy(newSlice, oldSlice)
	bas := toStructPointer(newSlice)
	oldHeader.Data = uintptr(bas)
	oldHeader.Len = newLen
	oldHeader.Cap = newLen

	return structPointer(unsafe.Pointer(uintptr(unsafe.Pointer(bas)) + uintptr(uintptr(newLen-1)*size)))
}

func structPointer_FieldPointer(p structPointer, f field) structPointer {
	return structPointer(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

func structPointer_GetRefStructPointer(p structPointer, f field) structPointer {
	return structPointer((*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}

func structPointer_GetSliceHeader(p structPointer, f field) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

func structPointer_Add(p structPointer, size field) structPointer {
	return structPointer(unsafe.Pointer(uintptr(p) + uintptr(size)))
}

func structPointer_Len(p structPointer, f field) int {
	return len(*(*[]interface{})(unsafe.Pointer(structPointer_GetRefStructPointer(p, f))))
}

func structPointer_StructRefSlice(p structPointer, f field, size uintptr) *structRefSlice {
	return &structRefSlice{p: p, f: f, size: size}
}


type structRefSlice struct {
	p    structPointer
	f    field
	size uintptr
}

func (v *structRefSlice) Len() int {
	return structPointer_Len(v.p, v.f)
}

func (v *structRefSlice) Index(i int) structPointer {
	ss := structPointer_GetStructPointer(v.p, v.f)
	ss1 := structPointer_GetRefStructPointer(ss, 0)
	return structPointer_Add(ss1, field(uintptr(i)*v.size))
}
