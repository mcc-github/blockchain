



























package proto

import (
	"reflect"
)


func (o *Buffer) dec_ref_struct_message(p *Properties, base structPointer) (err error) {
	raw, e := o.DecodeRawBytes(false)
	if e != nil {
		return e
	}

	
	if p.isUnmarshaler {
		panic("not supported, since this is a pointer receiver")
	}

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	bas := structPointer_FieldPointer(base, p.field)

	err = o.unmarshalType(p.stype, p.sprop, false, bas)
	o.buf = obuf
	o.index = oi

	return err
}


func (o *Buffer) dec_slice_ref_struct(p *Properties, is_group bool, base structPointer) error {
	newBas := appendStructPointer(base, p.field, p.sstype)

	if is_group {
		panic("not supported, maybe in future, if requested.")
	}

	raw, err := o.DecodeRawBytes(false)
	if err != nil {
		return err
	}

	
	if p.isUnmarshaler {
		panic("not supported, since this is not a pointer receiver.")
	}

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	err = o.unmarshalType(p.stype, p.sprop, is_group, newBas)

	o.buf = obuf
	o.index = oi

	return err
}


func (o *Buffer) dec_slice_ref_struct_message(p *Properties, base structPointer) error {
	return o.dec_slice_ref_struct(p, false, base)
}

func setPtrCustomType(base structPointer, f field, v interface{}) {
	if v == nil {
		return
	}
	structPointer_SetStructPointer(base, f, toStructPointer(reflect.ValueOf(v)))
}

func setCustomType(base structPointer, f field, value interface{}) {
	if value == nil {
		return
	}
	v := reflect.ValueOf(value).Elem()
	t := reflect.TypeOf(value).Elem()
	kind := t.Kind()
	switch kind {
	case reflect.Slice:
		slice := reflect.MakeSlice(t, v.Len(), v.Cap())
		reflect.Copy(slice, v)
		oldHeader := structPointer_GetSliceHeader(base, f)
		oldHeader.Data = slice.Pointer()
		oldHeader.Len = v.Len()
		oldHeader.Cap = v.Cap()
	default:
		size := reflect.TypeOf(value).Elem().Size()
		structPointer_Copy(toStructPointer(reflect.ValueOf(value)), structPointer_Add(base, f), int(size))
	}
}

func (o *Buffer) dec_custom_bytes(p *Properties, base structPointer) error {
	b, err := o.DecodeRawBytes(true)
	if err != nil {
		return err
	}
	i := reflect.New(p.ctype.Elem()).Interface()
	custom := (i).(Unmarshaler)
	if err := custom.Unmarshal(b); err != nil {
		return err
	}
	setPtrCustomType(base, p.field, custom)
	return nil
}

func (o *Buffer) dec_custom_ref_bytes(p *Properties, base structPointer) error {
	b, err := o.DecodeRawBytes(true)
	if err != nil {
		return err
	}
	i := reflect.New(p.ctype).Interface()
	custom := (i).(Unmarshaler)
	if err := custom.Unmarshal(b); err != nil {
		return err
	}
	if custom != nil {
		setCustomType(base, p.field, custom)
	}
	return nil
}


func (o *Buffer) dec_custom_slice_bytes(p *Properties, base structPointer) error {
	b, err := o.DecodeRawBytes(true)
	if err != nil {
		return err
	}
	i := reflect.New(p.ctype.Elem()).Interface()
	custom := (i).(Unmarshaler)
	if err := custom.Unmarshal(b); err != nil {
		return err
	}
	newBas := appendStructPointer(base, p.field, p.ctype)

	var zero field
	setCustomType(newBas, zero, custom)

	return nil
}
