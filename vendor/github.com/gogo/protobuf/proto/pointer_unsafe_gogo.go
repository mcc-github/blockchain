































package proto

import (
	"reflect"
	"unsafe"
)

func (p pointer) getRef() pointer {
	return pointer{p: (unsafe.Pointer)(&p.p)}
}

func (p pointer) appendRef(v pointer, typ reflect.Type) {
	slice := p.getSlice(typ)
	elem := v.asPointerTo(typ).Elem()
	newSlice := reflect.Append(slice, elem)
	slice.Set(newSlice)
}

func (p pointer) getSlice(typ reflect.Type) reflect.Value {
	sliceTyp := reflect.SliceOf(typ)
	slice := p.asPointerTo(sliceTyp)
	slice = slice.Elem()
	return slice
}
