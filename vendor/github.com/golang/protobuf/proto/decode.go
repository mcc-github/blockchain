






























package proto



import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)


var errOverflow = errors.New("proto: integer overflow")



var ErrInternalBadWireType = errors.New("proto: internal error: bad wiretype for oneof")











func DecodeVarint(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	
	return 0, 0
}

func (p *Buffer) decodeVarintSlow() (x uint64, err error) {
	i := p.index
	l := len(p.buf)

	for shift := uint(0); shift < 64; shift += 7 {
		if i >= l {
			err = io.ErrUnexpectedEOF
			return
		}
		b := p.buf[i]
		i++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			p.index = i
			return
		}
	}

	
	err = errOverflow
	return
}





func (p *Buffer) DecodeVarint() (x uint64, err error) {
	i := p.index
	buf := p.buf

	if i >= len(buf) {
		return 0, io.ErrUnexpectedEOF
	} else if buf[i] < 0x80 {
		p.index++
		return uint64(buf[i]), nil
	} else if len(buf)-i < 10 {
		return p.decodeVarintSlow()
	}

	var b uint64
	
	x = uint64(buf[i]) - 0x80
	i++

	b = uint64(buf[i])
	i++
	x += b << 7
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 7

	b = uint64(buf[i])
	i++
	x += b << 14
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 14

	b = uint64(buf[i])
	i++
	x += b << 21
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 21

	b = uint64(buf[i])
	i++
	x += b << 28
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 28

	b = uint64(buf[i])
	i++
	x += b << 35
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 35

	b = uint64(buf[i])
	i++
	x += b << 42
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 42

	b = uint64(buf[i])
	i++
	x += b << 49
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 49

	b = uint64(buf[i])
	i++
	x += b << 56
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 56

	b = uint64(buf[i])
	i++
	x += b << 63
	if b&0x80 == 0 {
		goto done
	}
	

	return 0, errOverflow

done:
	p.index = i
	return x, nil
}




func (p *Buffer) DecodeFixed64() (x uint64, err error) {
	
	i := p.index + 8
	if i < 0 || i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-8])
	x |= uint64(p.buf[i-7]) << 8
	x |= uint64(p.buf[i-6]) << 16
	x |= uint64(p.buf[i-5]) << 24
	x |= uint64(p.buf[i-4]) << 32
	x |= uint64(p.buf[i-3]) << 40
	x |= uint64(p.buf[i-2]) << 48
	x |= uint64(p.buf[i-1]) << 56
	return
}




func (p *Buffer) DecodeFixed32() (x uint64, err error) {
	
	i := p.index + 4
	if i < 0 || i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-4])
	x |= uint64(p.buf[i-3]) << 8
	x |= uint64(p.buf[i-2]) << 16
	x |= uint64(p.buf[i-1]) << 24
	return
}




func (p *Buffer) DecodeZigzag64() (x uint64, err error) {
	x, err = p.DecodeVarint()
	if err != nil {
		return
	}
	x = (x >> 1) ^ uint64((int64(x&1)<<63)>>63)
	return
}




func (p *Buffer) DecodeZigzag32() (x uint64, err error) {
	x, err = p.DecodeVarint()
	if err != nil {
		return
	}
	x = uint64((uint32(x) >> 1) ^ uint32((int32(x&1)<<31)>>31))
	return
}







func (p *Buffer) DecodeRawBytes(alloc bool) (buf []byte, err error) {
	n, err := p.DecodeVarint()
	if err != nil {
		return nil, err
	}

	nb := int(n)
	if nb < 0 {
		return nil, fmt.Errorf("proto: bad byte length %d", nb)
	}
	end := p.index + nb
	if end < p.index || end > len(p.buf) {
		return nil, io.ErrUnexpectedEOF
	}

	if !alloc {
		
		buf = p.buf[p.index:end]
		p.index += nb
		return
	}

	buf = make([]byte, nb)
	copy(buf, p.buf[p.index:])
	p.index += nb
	return
}



func (p *Buffer) DecodeStringBytes() (s string, err error) {
	buf, err := p.DecodeRawBytes(false)
	if err != nil {
		return
	}
	return string(buf), nil
}




func (o *Buffer) skipAndSave(t reflect.Type, tag, wire int, base structPointer, unrecField field) error {
	oi := o.index

	err := o.skip(t, tag, wire)
	if err != nil {
		return err
	}

	if !unrecField.IsValid() {
		return nil
	}

	ptr := structPointer_Bytes(base, unrecField)

	
	obuf := o.buf

	o.buf = *ptr
	o.EncodeVarint(uint64(tag<<3 | wire))
	*ptr = append(o.buf, obuf[oi:o.index]...)

	o.buf = obuf

	return nil
}


func (o *Buffer) skip(t reflect.Type, tag, wire int) error {

	var u uint64
	var err error

	switch wire {
	case WireVarint:
		_, err = o.DecodeVarint()
	case WireFixed64:
		_, err = o.DecodeFixed64()
	case WireBytes:
		_, err = o.DecodeRawBytes(false)
	case WireFixed32:
		_, err = o.DecodeFixed32()
	case WireStartGroup:
		for {
			u, err = o.DecodeVarint()
			if err != nil {
				break
			}
			fwire := int(u & 0x7)
			if fwire == WireEndGroup {
				break
			}
			ftag := int(u >> 3)
			err = o.skip(t, ftag, fwire)
			if err != nil {
				break
			}
		}
	default:
		err = fmt.Errorf("proto: can't skip unknown wire type %d for %s", wire, t)
	}
	return err
}






type Unmarshaler interface {
	Unmarshal([]byte) error
}








func Unmarshal(buf []byte, pb Message) error {
	pb.Reset()
	return UnmarshalMerge(buf, pb)
}







func UnmarshalMerge(buf []byte, pb Message) error {
	
	if u, ok := pb.(Unmarshaler); ok {
		return u.Unmarshal(buf)
	}
	return NewBuffer(buf).Unmarshal(pb)
}


func (p *Buffer) DecodeMessage(pb Message) error {
	enc, err := p.DecodeRawBytes(false)
	if err != nil {
		return err
	}
	return NewBuffer(enc).Unmarshal(pb)
}


func (p *Buffer) DecodeGroup(pb Message) error {
	typ, base, err := getbase(pb)
	if err != nil {
		return err
	}
	return p.unmarshalType(typ.Elem(), GetProperties(typ.Elem()), true, base)
}







func (p *Buffer) Unmarshal(pb Message) error {
	
	if u, ok := pb.(Unmarshaler); ok {
		err := u.Unmarshal(p.buf[p.index:])
		p.index = len(p.buf)
		return err
	}

	typ, base, err := getbase(pb)
	if err != nil {
		return err
	}

	err = p.unmarshalType(typ.Elem(), GetProperties(typ.Elem()), false, base)

	if collectStats {
		stats.Decode++
	}

	return err
}


func (o *Buffer) unmarshalType(st reflect.Type, prop *StructProperties, is_group bool, base structPointer) error {
	var state errorState
	required, reqFields := prop.reqCount, uint64(0)

	var err error
	for err == nil && o.index < len(o.buf) {
		oi := o.index
		var u uint64
		u, err = o.DecodeVarint()
		if err != nil {
			break
		}
		wire := int(u & 0x7)
		if wire == WireEndGroup {
			if is_group {
				if required > 0 {
					
					
					return &RequiredNotSetError{"{Unknown}"}
				}
				return nil 
			}
			return fmt.Errorf("proto: %s: wiretype end group for non-group", st)
		}
		tag := int(u >> 3)
		if tag <= 0 {
			return fmt.Errorf("proto: %s: illegal tag %d (wire type %d)", st, tag, wire)
		}
		fieldnum, ok := prop.decoderTags.get(tag)
		if !ok {
			
			if prop.extendable {
				if e, _ := extendable(structPointer_Interface(base, st)); isExtensionField(e, int32(tag)) {
					if err = o.skip(st, tag, wire); err == nil {
						extmap := e.extensionsWrite()
						ext := extmap[int32(tag)] 
						ext.enc = append(ext.enc, o.buf[oi:o.index]...)
						extmap[int32(tag)] = ext
					}
					continue
				}
			}
			
			if prop.oneofUnmarshaler != nil {
				m := structPointer_Interface(base, st).(Message)
				
				ok, err = prop.oneofUnmarshaler(m, tag, wire, o)
				if err == ErrInternalBadWireType {
					
					
					err = fmt.Errorf("bad wiretype for oneof field in %T", m)
				}
				if ok {
					continue
				}
			}
			err = o.skipAndSave(st, tag, wire, base, prop.unrecField)
			continue
		}
		p := prop.Prop[fieldnum]

		if p.dec == nil {
			fmt.Fprintf(os.Stderr, "proto: no protobuf decoder for %s.%s\n", st, st.Field(fieldnum).Name)
			continue
		}
		dec := p.dec
		if wire != WireStartGroup && wire != p.WireType {
			if wire == WireBytes && p.packedDec != nil {
				
				dec = p.packedDec
			} else {
				err = fmt.Errorf("proto: bad wiretype for field %s.%s: got wiretype %d, want %d", st, st.Field(fieldnum).Name, wire, p.WireType)
				continue
			}
		}
		decErr := dec(o, p, base)
		if decErr != nil && !state.shouldContinue(decErr, p) {
			err = decErr
		}
		if err == nil && p.Required {
			
			if tag <= 64 {
				
				var mask uint64 = 1 << uint64(tag-1)
				if reqFields&mask == 0 {
					
					reqFields |= mask
					required--
				}
			} else {
				
				
				
				
				required--
			}
		}
	}
	if err == nil {
		if is_group {
			return io.ErrUnexpectedEOF
		}
		if state.err != nil {
			return state.err
		}
		if required > 0 {
			
			
			
			return &RequiredNotSetError{"{Unknown}"}
		}
	}
	return err
}









const (
	boolPoolSize   = 16
	uint32PoolSize = 8
	uint64PoolSize = 4
)


func (o *Buffer) dec_bool(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}
	if len(o.bools) == 0 {
		o.bools = make([]bool, boolPoolSize)
	}
	o.bools[0] = u != 0
	*structPointer_Bool(base, p.field) = &o.bools[0]
	o.bools = o.bools[1:]
	return nil
}

func (o *Buffer) dec_proto3_bool(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}
	*structPointer_BoolVal(base, p.field) = u != 0
	return nil
}


func (o *Buffer) dec_int32(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}
	word32_Set(structPointer_Word32(base, p.field), o, uint32(u))
	return nil
}

func (o *Buffer) dec_proto3_int32(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}
	word32Val_Set(structPointer_Word32Val(base, p.field), uint32(u))
	return nil
}


func (o *Buffer) dec_int64(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}
	word64_Set(structPointer_Word64(base, p.field), o, u)
	return nil
}

func (o *Buffer) dec_proto3_int64(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}
	word64Val_Set(structPointer_Word64Val(base, p.field), o, u)
	return nil
}


func (o *Buffer) dec_string(p *Properties, base structPointer) error {
	s, err := o.DecodeStringBytes()
	if err != nil {
		return err
	}
	*structPointer_String(base, p.field) = &s
	return nil
}

func (o *Buffer) dec_proto3_string(p *Properties, base structPointer) error {
	s, err := o.DecodeStringBytes()
	if err != nil {
		return err
	}
	*structPointer_StringVal(base, p.field) = s
	return nil
}


func (o *Buffer) dec_slice_byte(p *Properties, base structPointer) error {
	b, err := o.DecodeRawBytes(true)
	if err != nil {
		return err
	}
	*structPointer_Bytes(base, p.field) = b
	return nil
}


func (o *Buffer) dec_slice_bool(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}
	v := structPointer_BoolSlice(base, p.field)
	*v = append(*v, u != 0)
	return nil
}


func (o *Buffer) dec_slice_packed_bool(p *Properties, base structPointer) error {
	v := structPointer_BoolSlice(base, p.field)

	nn, err := o.DecodeVarint()
	if err != nil {
		return err
	}
	nb := int(nn) 
	fin := o.index + nb
	if fin < o.index {
		return errOverflow
	}

	y := *v
	for o.index < fin {
		u, err := p.valDec(o)
		if err != nil {
			return err
		}
		y = append(y, u != 0)
	}

	*v = y
	return nil
}


func (o *Buffer) dec_slice_int32(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}
	structPointer_Word32Slice(base, p.field).Append(uint32(u))
	return nil
}


func (o *Buffer) dec_slice_packed_int32(p *Properties, base structPointer) error {
	v := structPointer_Word32Slice(base, p.field)

	nn, err := o.DecodeVarint()
	if err != nil {
		return err
	}
	nb := int(nn) 

	fin := o.index + nb
	if fin < o.index {
		return errOverflow
	}
	for o.index < fin {
		u, err := p.valDec(o)
		if err != nil {
			return err
		}
		v.Append(uint32(u))
	}
	return nil
}


func (o *Buffer) dec_slice_int64(p *Properties, base structPointer) error {
	u, err := p.valDec(o)
	if err != nil {
		return err
	}

	structPointer_Word64Slice(base, p.field).Append(u)
	return nil
}


func (o *Buffer) dec_slice_packed_int64(p *Properties, base structPointer) error {
	v := structPointer_Word64Slice(base, p.field)

	nn, err := o.DecodeVarint()
	if err != nil {
		return err
	}
	nb := int(nn) 

	fin := o.index + nb
	if fin < o.index {
		return errOverflow
	}
	for o.index < fin {
		u, err := p.valDec(o)
		if err != nil {
			return err
		}
		v.Append(u)
	}
	return nil
}


func (o *Buffer) dec_slice_string(p *Properties, base structPointer) error {
	s, err := o.DecodeStringBytes()
	if err != nil {
		return err
	}
	v := structPointer_StringSlice(base, p.field)
	*v = append(*v, s)
	return nil
}


func (o *Buffer) dec_slice_slice_byte(p *Properties, base structPointer) error {
	b, err := o.DecodeRawBytes(true)
	if err != nil {
		return err
	}
	v := structPointer_BytesSlice(base, p.field)
	*v = append(*v, b)
	return nil
}


func (o *Buffer) dec_new_map(p *Properties, base structPointer) error {
	raw, err := o.DecodeRawBytes(false)
	if err != nil {
		return err
	}
	oi := o.index       
	o.index -= len(raw) 

	mptr := structPointer_NewAt(base, p.field, p.mtype) 
	if mptr.Elem().IsNil() {
		mptr.Elem().Set(reflect.MakeMap(mptr.Type().Elem()))
	}
	v := mptr.Elem() 

	
	
	keyptr := reflect.New(reflect.PtrTo(p.mtype.Key())).Elem() 
	keybase := toStructPointer(keyptr.Addr())                  

	var valbase structPointer
	var valptr reflect.Value
	switch p.mtype.Elem().Kind() {
	case reflect.Slice:
		
		var dummy []byte
		valptr = reflect.ValueOf(&dummy)  
		valbase = toStructPointer(valptr) 
	case reflect.Ptr:
		
		valptr = reflect.New(reflect.PtrTo(p.mtype.Elem())).Elem() 
		valptr.Set(reflect.New(valptr.Type().Elem()))
		valbase = toStructPointer(valptr)
	default:
		
		valptr = reflect.New(reflect.PtrTo(p.mtype.Elem())).Elem() 
		valbase = toStructPointer(valptr.Addr())                   
	}

	
	
	
	for o.index < oi {
		
		
		tagcode := o.buf[o.index]
		o.index++
		switch tagcode {
		case p.mkeyprop.tagcode[0]:
			if err := p.mkeyprop.dec(o, p.mkeyprop, keybase); err != nil {
				return err
			}
		case p.mvalprop.tagcode[0]:
			if err := p.mvalprop.dec(o, p.mvalprop, valbase); err != nil {
				return err
			}
		default:
			
			return fmt.Errorf("proto: bad map data tag %d", raw[0])
		}
	}
	keyelem, valelem := keyptr.Elem(), valptr.Elem()
	if !keyelem.IsValid() {
		keyelem = reflect.Zero(p.mtype.Key())
	}
	if !valelem.IsValid() {
		valelem = reflect.Zero(p.mtype.Elem())
	}

	v.SetMapIndex(keyelem, valelem)
	return nil
}


func (o *Buffer) dec_struct_group(p *Properties, base structPointer) error {
	bas := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(bas) {
		
		bas = toStructPointer(reflect.New(p.stype))
		structPointer_SetStructPointer(base, p.field, bas)
	}
	return o.unmarshalType(p.stype, p.sprop, true, bas)
}


func (o *Buffer) dec_struct_message(p *Properties, base structPointer) (err error) {
	raw, e := o.DecodeRawBytes(false)
	if e != nil {
		return e
	}

	bas := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(bas) {
		
		bas = toStructPointer(reflect.New(p.stype))
		structPointer_SetStructPointer(base, p.field, bas)
	}

	
	if p.isUnmarshaler {
		iv := structPointer_Interface(bas, p.stype)
		return iv.(Unmarshaler).Unmarshal(raw)
	}

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	err = o.unmarshalType(p.stype, p.sprop, false, bas)
	o.buf = obuf
	o.index = oi

	return err
}


func (o *Buffer) dec_slice_struct_message(p *Properties, base structPointer) error {
	return o.dec_slice_struct(p, false, base)
}


func (o *Buffer) dec_slice_struct_group(p *Properties, base structPointer) error {
	return o.dec_slice_struct(p, true, base)
}


func (o *Buffer) dec_slice_struct(p *Properties, is_group bool, base structPointer) error {
	v := reflect.New(p.stype)
	bas := toStructPointer(v)
	structPointer_StructPointerSlice(base, p.field).Append(bas)

	if is_group {
		err := o.unmarshalType(p.stype, p.sprop, is_group, bas)
		return err
	}

	raw, err := o.DecodeRawBytes(false)
	if err != nil {
		return err
	}

	
	if p.isUnmarshaler {
		iv := v.Interface()
		return iv.(Unmarshaler).Unmarshal(raw)
	}

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	err = o.unmarshalType(p.stype, p.sprop, is_group, bas)

	o.buf = obuf
	o.index = oi

	return err
}
