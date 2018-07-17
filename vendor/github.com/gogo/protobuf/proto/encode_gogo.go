



































package proto

import (
	"reflect"
)

func NewRequiredNotSetError(field string) *RequiredNotSetError {
	return &RequiredNotSetError{field}
}

type Sizer interface {
	Size() int
}

func (o *Buffer) enc_ext_slice_byte(p *Properties, base structPointer) error {
	s := *structPointer_Bytes(base, p.field)
	if s == nil {
		return ErrNil
	}
	o.buf = append(o.buf, s...)
	return nil
}

func size_ext_slice_byte(p *Properties, base structPointer) (n int) {
	s := *structPointer_Bytes(base, p.field)
	if s == nil {
		return 0
	}
	n += len(s)
	return
}


func (o *Buffer) enc_ref_bool(p *Properties, base structPointer) error {
	v := *structPointer_BoolVal(base, p.field)
	x := 0
	if v {
		x = 1
	}
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
}

func size_ref_bool(p *Properties, base structPointer) int {
	return len(p.tagcode) + 1 
}


func (o *Buffer) enc_ref_int32(p *Properties, base structPointer) error {
	v := structPointer_Word32Val(base, p.field)
	x := int32(word32Val_Get(v))
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
}

func size_ref_int32(p *Properties, base structPointer) (n int) {
	v := structPointer_Word32Val(base, p.field)
	x := int32(word32Val_Get(v))
	n += len(p.tagcode)
	n += p.valSize(uint64(x))
	return
}

func (o *Buffer) enc_ref_uint32(p *Properties, base structPointer) error {
	v := structPointer_Word32Val(base, p.field)
	x := word32Val_Get(v)
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
}

func size_ref_uint32(p *Properties, base structPointer) (n int) {
	v := structPointer_Word32Val(base, p.field)
	x := word32Val_Get(v)
	n += len(p.tagcode)
	n += p.valSize(uint64(x))
	return
}


func (o *Buffer) enc_ref_int64(p *Properties, base structPointer) error {
	v := structPointer_Word64Val(base, p.field)
	x := word64Val_Get(v)
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, x)
	return nil
}

func size_ref_int64(p *Properties, base structPointer) (n int) {
	v := structPointer_Word64Val(base, p.field)
	x := word64Val_Get(v)
	n += len(p.tagcode)
	n += p.valSize(x)
	return
}


func (o *Buffer) enc_ref_string(p *Properties, base structPointer) error {
	v := *structPointer_StringVal(base, p.field)
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeStringBytes(v)
	return nil
}

func size_ref_string(p *Properties, base structPointer) (n int) {
	v := *structPointer_StringVal(base, p.field)
	n += len(p.tagcode)
	n += sizeStringBytes(v)
	return
}


func (o *Buffer) enc_ref_struct_message(p *Properties, base structPointer) error {
	var state errorState
	structp := structPointer_GetRefStructPointer(base, p.field)
	if structPointer_IsNil(structp) {
		return ErrNil
	}

	
	if p.isMarshaler {
		m := structPointer_Interface(structp, p.stype).(Marshaler)
		data, err := m.Marshal()
		if err != nil && !state.shouldContinue(err, nil) {
			return err
		}
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
		return nil
	}

	o.buf = append(o.buf, p.tagcode...)
	return o.enc_len_struct(p.sprop, structp, &state)
}


func size_ref_struct_message(p *Properties, base structPointer) int {
	structp := structPointer_GetRefStructPointer(base, p.field)
	if structPointer_IsNil(structp) {
		return 0
	}

	
	if p.isMarshaler {
		m := structPointer_Interface(structp, p.stype).(Marshaler)
		data, _ := m.Marshal()
		n0 := len(p.tagcode)
		n1 := sizeRawBytes(data)
		return n0 + n1
	}

	n0 := len(p.tagcode)
	n1 := size_struct(p.sprop, structp)
	n2 := sizeVarint(uint64(n1)) 
	return n0 + n1 + n2
}


func (o *Buffer) enc_slice_ref_struct_message(p *Properties, base structPointer) error {
	var state errorState
	ss := structPointer_StructRefSlice(base, p.field, p.stype.Size())
	l := ss.Len()
	for i := 0; i < l; i++ {
		structp := ss.Index(i)
		if structPointer_IsNil(structp) {
			return errRepeatedHasNil
		}

		
		if p.isMarshaler {
			m := structPointer_Interface(structp, p.stype).(Marshaler)
			data, err := m.Marshal()
			if err != nil && !state.shouldContinue(err, nil) {
				return err
			}
			o.buf = append(o.buf, p.tagcode...)
			o.EncodeRawBytes(data)
			continue
		}

		o.buf = append(o.buf, p.tagcode...)
		err := o.enc_len_struct(p.sprop, structp, &state)
		if err != nil && !state.shouldContinue(err, nil) {
			if err == ErrNil {
				return errRepeatedHasNil
			}
			return err
		}

	}
	return state.err
}


func size_slice_ref_struct_message(p *Properties, base structPointer) (n int) {
	ss := structPointer_StructRefSlice(base, p.field, p.stype.Size())
	l := ss.Len()
	n += l * len(p.tagcode)
	for i := 0; i < l; i++ {
		structp := ss.Index(i)
		if structPointer_IsNil(structp) {
			return 
		}

		
		if p.isMarshaler {
			m := structPointer_Interface(structp, p.stype).(Marshaler)
			data, _ := m.Marshal()
			n += len(p.tagcode)
			n += sizeRawBytes(data)
			continue
		}

		n0 := size_struct(p.sprop, structp)
		n1 := sizeVarint(uint64(n0)) 
		n += n0 + n1
	}
	return
}

func (o *Buffer) enc_custom_bytes(p *Properties, base structPointer) error {
	i := structPointer_InterfaceRef(base, p.field, p.ctype)
	if i == nil {
		return ErrNil
	}
	custom := i.(Marshaler)
	data, err := custom.Marshal()
	if err != nil {
		return err
	}
	if data == nil {
		return ErrNil
	}
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(data)
	return nil
}

func size_custom_bytes(p *Properties, base structPointer) (n int) {
	n += len(p.tagcode)
	i := structPointer_InterfaceRef(base, p.field, p.ctype)
	if i == nil {
		return 0
	}
	custom := i.(Marshaler)
	data, _ := custom.Marshal()
	n += sizeRawBytes(data)
	return
}

func (o *Buffer) enc_custom_ref_bytes(p *Properties, base structPointer) error {
	custom := structPointer_InterfaceAt(base, p.field, p.ctype).(Marshaler)
	data, err := custom.Marshal()
	if err != nil {
		return err
	}
	if data == nil {
		return ErrNil
	}
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(data)
	return nil
}

func size_custom_ref_bytes(p *Properties, base structPointer) (n int) {
	n += len(p.tagcode)
	i := structPointer_InterfaceAt(base, p.field, p.ctype)
	if i == nil {
		return 0
	}
	custom := i.(Marshaler)
	data, _ := custom.Marshal()
	n += sizeRawBytes(data)
	return
}

func (o *Buffer) enc_custom_slice_bytes(p *Properties, base structPointer) error {
	inter := structPointer_InterfaceRef(base, p.field, p.ctype)
	if inter == nil {
		return ErrNil
	}
	slice := reflect.ValueOf(inter)
	l := slice.Len()
	for i := 0; i < l; i++ {
		v := slice.Index(i)
		custom := v.Interface().(Marshaler)
		data, err := custom.Marshal()
		if err != nil {
			return err
		}
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
	}
	return nil
}

func size_custom_slice_bytes(p *Properties, base structPointer) (n int) {
	inter := structPointer_InterfaceRef(base, p.field, p.ctype)
	if inter == nil {
		return 0
	}
	slice := reflect.ValueOf(inter)
	l := slice.Len()
	n += l * len(p.tagcode)
	for i := 0; i < l; i++ {
		v := slice.Index(i)
		custom := v.Interface().(Marshaler)
		data, _ := custom.Marshal()
		n += sizeRawBytes(data)
	}
	return
}
