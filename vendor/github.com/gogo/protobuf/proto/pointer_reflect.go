




































package proto

import (
	"math"
	"reflect"
)


type structPointer struct {
	v reflect.Value
}



func toStructPointer(v reflect.Value) structPointer {
	return structPointer{v}
}


func structPointer_IsNil(p structPointer) bool {
	return p.v.IsNil()
}


func structPointer_Interface(p structPointer, _ reflect.Type) interface{} {
	return p.v.Interface()
}




type field []int


func toField(f *reflect.StructField) field {
	return f.Index
}


var invalidField = field(nil)


func (f field) IsValid() bool { return f != nil }


func structPointer_field(p structPointer, f field) reflect.Value {
	
	
	
	
	
	
	if f == nil {
		return p.v.Elem()
	}

	return p.v.Elem().FieldByIndex(f)
}


func structPointer_ifield(p structPointer, f field) interface{} {
	return structPointer_field(p, f).Addr().Interface()
}


func structPointer_Bytes(p structPointer, f field) *[]byte {
	return structPointer_ifield(p, f).(*[]byte)
}


func structPointer_BytesSlice(p structPointer, f field) *[][]byte {
	return structPointer_ifield(p, f).(*[][]byte)
}


func structPointer_Bool(p structPointer, f field) **bool {
	return structPointer_ifield(p, f).(**bool)
}


func structPointer_BoolVal(p structPointer, f field) *bool {
	return structPointer_ifield(p, f).(*bool)
}


func structPointer_BoolSlice(p structPointer, f field) *[]bool {
	return structPointer_ifield(p, f).(*[]bool)
}


func structPointer_String(p structPointer, f field) **string {
	return structPointer_ifield(p, f).(**string)
}


func structPointer_StringVal(p structPointer, f field) *string {
	return structPointer_ifield(p, f).(*string)
}


func structPointer_StringSlice(p structPointer, f field) *[]string {
	return structPointer_ifield(p, f).(*[]string)
}


func structPointer_Extensions(p structPointer, f field) *XXX_InternalExtensions {
	return structPointer_ifield(p, f).(*XXX_InternalExtensions)
}


func structPointer_ExtMap(p structPointer, f field) *map[int32]Extension {
	return structPointer_ifield(p, f).(*map[int32]Extension)
}


func structPointer_NewAt(p structPointer, f field, typ reflect.Type) reflect.Value {
	return structPointer_field(p, f).Addr()
}


func structPointer_SetStructPointer(p structPointer, f field, q structPointer) {
	structPointer_field(p, f).Set(q.v)
}


func structPointer_GetStructPointer(p structPointer, f field) structPointer {
	return structPointer{structPointer_field(p, f)}
}


func structPointer_StructPointerSlice(p structPointer, f field) structPointerSlice {
	return structPointerSlice{structPointer_field(p, f)}
}



type structPointerSlice struct {
	v reflect.Value
}

func (p structPointerSlice) Len() int                  { return p.v.Len() }
func (p structPointerSlice) Index(i int) structPointer { return structPointer{p.v.Index(i)} }
func (p structPointerSlice) Append(q structPointer) {
	p.v.Set(reflect.Append(p.v, q.v))
}

var (
	int32Type   = reflect.TypeOf(int32(0))
	uint32Type  = reflect.TypeOf(uint32(0))
	float32Type = reflect.TypeOf(float32(0))
	int64Type   = reflect.TypeOf(int64(0))
	uint64Type  = reflect.TypeOf(uint64(0))
	float64Type = reflect.TypeOf(float64(0))
)



type word32 struct {
	v reflect.Value
}


func word32_IsNil(p word32) bool {
	return p.v.IsNil()
}


func word32_Set(p word32, o *Buffer, x uint32) {
	t := p.v.Type().Elem()
	switch t {
	case int32Type:
		if len(o.int32s) == 0 {
			o.int32s = make([]int32, uint32PoolSize)
		}
		o.int32s[0] = int32(x)
		p.v.Set(reflect.ValueOf(&o.int32s[0]))
		o.int32s = o.int32s[1:]
		return
	case uint32Type:
		if len(o.uint32s) == 0 {
			o.uint32s = make([]uint32, uint32PoolSize)
		}
		o.uint32s[0] = x
		p.v.Set(reflect.ValueOf(&o.uint32s[0]))
		o.uint32s = o.uint32s[1:]
		return
	case float32Type:
		if len(o.float32s) == 0 {
			o.float32s = make([]float32, uint32PoolSize)
		}
		o.float32s[0] = math.Float32frombits(x)
		p.v.Set(reflect.ValueOf(&o.float32s[0]))
		o.float32s = o.float32s[1:]
		return
	}

	
	p.v.Set(reflect.New(t))
	p.v.Elem().SetInt(int64(int32(x)))
}


func word32_Get(p word32) uint32 {
	elem := p.v.Elem()
	switch elem.Kind() {
	case reflect.Int32:
		return uint32(elem.Int())
	case reflect.Uint32:
		return uint32(elem.Uint())
	case reflect.Float32:
		return math.Float32bits(float32(elem.Float()))
	}
	panic("unreachable")
}


func structPointer_Word32(p structPointer, f field) word32 {
	return word32{structPointer_field(p, f)}
}



type word32Val struct {
	v reflect.Value
}


func word32Val_Set(p word32Val, x uint32) {
	switch p.v.Type() {
	case int32Type:
		p.v.SetInt(int64(x))
		return
	case uint32Type:
		p.v.SetUint(uint64(x))
		return
	case float32Type:
		p.v.SetFloat(float64(math.Float32frombits(x)))
		return
	}

	
	p.v.SetInt(int64(int32(x)))
}


func word32Val_Get(p word32Val) uint32 {
	elem := p.v
	switch elem.Kind() {
	case reflect.Int32:
		return uint32(elem.Int())
	case reflect.Uint32:
		return uint32(elem.Uint())
	case reflect.Float32:
		return math.Float32bits(float32(elem.Float()))
	}
	panic("unreachable")
}


func structPointer_Word32Val(p structPointer, f field) word32Val {
	return word32Val{structPointer_field(p, f)}
}



type word32Slice struct {
	v reflect.Value
}

func (p word32Slice) Append(x uint32) {
	n, m := p.v.Len(), p.v.Cap()
	if n < m {
		p.v.SetLen(n + 1)
	} else {
		t := p.v.Type().Elem()
		p.v.Set(reflect.Append(p.v, reflect.Zero(t)))
	}
	elem := p.v.Index(n)
	switch elem.Kind() {
	case reflect.Int32:
		elem.SetInt(int64(int32(x)))
	case reflect.Uint32:
		elem.SetUint(uint64(x))
	case reflect.Float32:
		elem.SetFloat(float64(math.Float32frombits(x)))
	}
}

func (p word32Slice) Len() int {
	return p.v.Len()
}

func (p word32Slice) Index(i int) uint32 {
	elem := p.v.Index(i)
	switch elem.Kind() {
	case reflect.Int32:
		return uint32(elem.Int())
	case reflect.Uint32:
		return uint32(elem.Uint())
	case reflect.Float32:
		return math.Float32bits(float32(elem.Float()))
	}
	panic("unreachable")
}


func structPointer_Word32Slice(p structPointer, f field) word32Slice {
	return word32Slice{structPointer_field(p, f)}
}


type word64 struct {
	v reflect.Value
}

func word64_Set(p word64, o *Buffer, x uint64) {
	t := p.v.Type().Elem()
	switch t {
	case int64Type:
		if len(o.int64s) == 0 {
			o.int64s = make([]int64, uint64PoolSize)
		}
		o.int64s[0] = int64(x)
		p.v.Set(reflect.ValueOf(&o.int64s[0]))
		o.int64s = o.int64s[1:]
		return
	case uint64Type:
		if len(o.uint64s) == 0 {
			o.uint64s = make([]uint64, uint64PoolSize)
		}
		o.uint64s[0] = x
		p.v.Set(reflect.ValueOf(&o.uint64s[0]))
		o.uint64s = o.uint64s[1:]
		return
	case float64Type:
		if len(o.float64s) == 0 {
			o.float64s = make([]float64, uint64PoolSize)
		}
		o.float64s[0] = math.Float64frombits(x)
		p.v.Set(reflect.ValueOf(&o.float64s[0]))
		o.float64s = o.float64s[1:]
		return
	}
	panic("unreachable")
}

func word64_IsNil(p word64) bool {
	return p.v.IsNil()
}

func word64_Get(p word64) uint64 {
	elem := p.v.Elem()
	switch elem.Kind() {
	case reflect.Int64:
		return uint64(elem.Int())
	case reflect.Uint64:
		return elem.Uint()
	case reflect.Float64:
		return math.Float64bits(elem.Float())
	}
	panic("unreachable")
}

func structPointer_Word64(p structPointer, f field) word64 {
	return word64{structPointer_field(p, f)}
}


type word64Val struct {
	v reflect.Value
}

func word64Val_Set(p word64Val, o *Buffer, x uint64) {
	switch p.v.Type() {
	case int64Type:
		p.v.SetInt(int64(x))
		return
	case uint64Type:
		p.v.SetUint(x)
		return
	case float64Type:
		p.v.SetFloat(math.Float64frombits(x))
		return
	}
	panic("unreachable")
}

func word64Val_Get(p word64Val) uint64 {
	elem := p.v
	switch elem.Kind() {
	case reflect.Int64:
		return uint64(elem.Int())
	case reflect.Uint64:
		return elem.Uint()
	case reflect.Float64:
		return math.Float64bits(elem.Float())
	}
	panic("unreachable")
}

func structPointer_Word64Val(p structPointer, f field) word64Val {
	return word64Val{structPointer_field(p, f)}
}

type word64Slice struct {
	v reflect.Value
}

func (p word64Slice) Append(x uint64) {
	n, m := p.v.Len(), p.v.Cap()
	if n < m {
		p.v.SetLen(n + 1)
	} else {
		t := p.v.Type().Elem()
		p.v.Set(reflect.Append(p.v, reflect.Zero(t)))
	}
	elem := p.v.Index(n)
	switch elem.Kind() {
	case reflect.Int64:
		elem.SetInt(int64(int64(x)))
	case reflect.Uint64:
		elem.SetUint(uint64(x))
	case reflect.Float64:
		elem.SetFloat(float64(math.Float64frombits(x)))
	}
}

func (p word64Slice) Len() int {
	return p.v.Len()
}

func (p word64Slice) Index(i int) uint64 {
	elem := p.v.Index(i)
	switch elem.Kind() {
	case reflect.Int64:
		return uint64(elem.Int())
	case reflect.Uint64:
		return uint64(elem.Uint())
	case reflect.Float64:
		return math.Float64bits(float64(elem.Float()))
	}
	panic("unreachable")
}

func structPointer_Word64Slice(p structPointer, f field) word64Slice {
	return word64Slice{structPointer_field(p, f)}
}
