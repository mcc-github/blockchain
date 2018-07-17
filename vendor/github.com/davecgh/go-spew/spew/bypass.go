



















package spew

import (
	"reflect"
	"unsafe"
)

const (
	
	
	UnsafeDisabled = false

	
	ptrSize = unsafe.Sizeof((*byte)(nil))
)

var (
	
	
	
	
	
	
	offsetPtr    = uintptr(ptrSize)
	offsetScalar = uintptr(0)
	offsetFlag   = uintptr(ptrSize * 2)

	
	
	
	
	
	
	
	
	
	
	
	
	flagKindWidth = uintptr(5)
	flagKindShift = uintptr(flagKindWidth - 1)
	flagRO        = uintptr(1 << 0)
	flagIndir     = uintptr(1 << 1)
)

func init() {
	
	
	
	
	
	
	
	
	
	
	vv := reflect.ValueOf(0xf00)
	if unsafe.Sizeof(vv) == (ptrSize * 4) {
		offsetScalar = ptrSize * 2
		offsetFlag = ptrSize * 3
	}

	
	
	
	
	
	upf := unsafe.Pointer(uintptr(unsafe.Pointer(&vv)) + offsetFlag)
	upfv := *(*uintptr)(upf)
	flagKindMask := uintptr((1<<flagKindWidth - 1) << flagKindShift)
	if (upfv&flagKindMask)>>flagKindShift != uintptr(reflect.Int) {
		flagKindShift = 0
		flagRO = 1 << 5
		flagIndir = 1 << 6

		
		
		
		
		
		
		
		
		
		
		if upfv&flagIndir == 0 {
			flagRO = 3 << 5
			flagIndir = 1 << 7
		}
	}
}










func unsafeReflectValue(v reflect.Value) (rv reflect.Value) {
	indirects := 1
	vt := v.Type()
	upv := unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + offsetPtr)
	rvf := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + offsetFlag))
	if rvf&flagIndir != 0 {
		vt = reflect.PtrTo(v.Type())
		indirects++
	} else if offsetScalar != 0 {
		
		
		switch vt.Kind() {
		case reflect.Uintptr:
		case reflect.Chan:
		case reflect.Func:
		case reflect.Map:
		case reflect.Ptr:
		case reflect.UnsafePointer:
		default:
			upv = unsafe.Pointer(uintptr(unsafe.Pointer(&v)) +
				offsetScalar)
		}
	}

	pv := reflect.NewAt(vt, upv)
	rv = pv
	for i := 0; i < indirects; i++ {
		rv = rv.Elem()
	}
	return rv
}
