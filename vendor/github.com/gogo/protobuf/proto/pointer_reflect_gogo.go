





























package proto

import (
	"reflect"
)

func structPointer_FieldPointer(p structPointer, f field) structPointer {
	panic("not implemented")
}

func appendStructPointer(base structPointer, f field, typ reflect.Type) structPointer {
	panic("not implemented")
}

func structPointer_InterfaceAt(p structPointer, f field, t reflect.Type) interface{} {
	panic("not implemented")
}

func structPointer_InterfaceRef(p structPointer, f field, t reflect.Type) interface{} {
	panic("not implemented")
}

func structPointer_GetRefStructPointer(p structPointer, f field) structPointer {
	panic("not implemented")
}

func structPointer_Add(p structPointer, size field) structPointer {
	panic("not implemented")
}

func structPointer_Len(p structPointer, f field) int {
	panic("not implemented")
}

func structPointer_GetSliceHeader(p structPointer, f field) *reflect.SliceHeader {
	panic("not implemented")
}

func structPointer_Copy(oldptr structPointer, newptr structPointer, size int) {
	panic("not implemented")
}

func structPointer_StructRefSlice(p structPointer, f field, size uintptr) *structRefSlice {
	panic("not implemented")
}

type structRefSlice struct{}

func (v *structRefSlice) Len() int {
	panic("not implemented")
}

func (v *structRefSlice) Index(i int) structPointer {
	panic("not implemented")
}
