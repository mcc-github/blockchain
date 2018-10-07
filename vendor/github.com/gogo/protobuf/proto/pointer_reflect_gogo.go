

































package proto

import (
	"reflect"
)



func (p pointer) getRef() pointer {
	return pointer{v: p.v.Addr()}
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
