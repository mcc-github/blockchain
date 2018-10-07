





















package spew

import (
	"reflect"
	"unsafe"
)

const (
	
	
	UnsafeDisabled = false

	
	ptrSize = unsafe.Sizeof((*byte)(nil))
)

type flag uintptr

var (
	
	
	flagRO flag

	
	
	flagAddr flag
)




const flagKindMask = flag(0x1f)




var okFlags = []struct {
	ro, addr flag
}{{
	
	ro:   1 << 5,
	addr: 1 << 7,
}, {
	
	ro:   1<<5 | 1<<6,
	addr: 1 << 8,
}}

var flagValOffset = func() uintptr {
	field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag")
	if !ok {
		panic("reflect.Value has no flag field")
	}
	return field.Offset
}()


func flagField(v *reflect.Value) *flag {
	return (*flag)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + flagValOffset))
}










func unsafeReflectValue(v reflect.Value) reflect.Value {
	if !v.IsValid() || (v.CanInterface() && v.CanAddr()) {
		return v
	}
	flagFieldPtr := flagField(&v)
	*flagFieldPtr &^= flagRO
	*flagFieldPtr |= flagAddr
	return v
}



func init() {
	field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag")
	if !ok {
		panic("reflect.Value has no flag field")
	}
	if field.Type.Kind() != reflect.TypeOf(flag(0)).Kind() {
		panic("reflect.Value flag field has changed kind")
	}
	type t0 int
	var t struct {
		A t0
		
		t0
		
		a t0
	}
	vA := reflect.ValueOf(t).FieldByName("A")
	va := reflect.ValueOf(t).FieldByName("a")
	vt0 := reflect.ValueOf(t).FieldByName("t0")

	
	
	flagPublic := *flagField(&vA)
	flagWithRO := *flagField(&va) | *flagField(&vt0)
	flagRO = flagPublic ^ flagWithRO

	
	
	vPtrA := reflect.ValueOf(&t).Elem().FieldByName("A")
	flagNoPtr := *flagField(&vA)
	flagPtr := *flagField(&vPtrA)
	flagAddr = flagNoPtr ^ flagPtr

	
	for _, f := range okFlags {
		if flagRO == f.ro && flagAddr == f.addr {
			return
		}
	}
	panic("reflect.Value read-only flag has changed semantics")
}
