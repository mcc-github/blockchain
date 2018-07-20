






























package proto

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)







func (a *InternalMessageInfo) Unmarshal(msg Message, b []byte) error {
	
	
	u := atomicLoadUnmarshalInfo(&a.unmarshal)
	if u == nil {
		
		u = getUnmarshalInfo(reflect.TypeOf(msg).Elem())
		atomicStoreUnmarshalInfo(&a.unmarshal, u)
	}
	
	err := u.unmarshal(toPointer(&msg), b)
	return err
}

type unmarshalInfo struct {
	typ reflect.Type 

	
	
	initialized     int32
	lock            sync.Mutex                    
	dense           []unmarshalFieldInfo          
	sparse          map[uint64]unmarshalFieldInfo 
	reqFields       []string                      
	reqMask         uint64                        
	unrecognized    field                         
	extensions      field                         
	oldExtensions   field                         
	extensionRanges []ExtensionRange              
	isMessageSet    bool                          
}





type unmarshaler func(b []byte, f pointer, w int) ([]byte, error)

type unmarshalFieldInfo struct {
	
	field field

	
	unmarshal unmarshaler

	
	reqMask uint64
}

var (
	unmarshalInfoMap  = map[reflect.Type]*unmarshalInfo{}
	unmarshalInfoLock sync.Mutex
)




func getUnmarshalInfo(t reflect.Type) *unmarshalInfo {
	
	
	
	
	unmarshalInfoLock.Lock()
	defer unmarshalInfoLock.Unlock()
	u := unmarshalInfoMap[t]
	if u == nil {
		u = &unmarshalInfo{typ: t}
		
		
		unmarshalInfoMap[t] = u
	}
	return u
}






func (u *unmarshalInfo) unmarshal(m pointer, b []byte) error {
	if atomic.LoadInt32(&u.initialized) == 0 {
		u.computeUnmarshalInfo()
	}
	if u.isMessageSet {
		return UnmarshalMessageSet(b, m.offset(u.extensions).toExtensions())
	}
	var reqMask uint64            
	var rnse *RequiredNotSetError 
	for len(b) > 0 {
		
		
		var x uint64
		if b[0] < 128 {
			x = uint64(b[0])
			b = b[1:]
		} else if len(b) >= 2 && b[1] < 128 {
			x = uint64(b[0]&0x7f) + uint64(b[1])<<7
			b = b[2:]
		} else {
			var n int
			x, n = decodeVarint(b)
			if n == 0 {
				return io.ErrUnexpectedEOF
			}
			b = b[n:]
		}
		tag := x >> 3
		wire := int(x) & 7

		
		var f unmarshalFieldInfo
		if tag < uint64(len(u.dense)) {
			f = u.dense[tag]
		} else {
			f = u.sparse[tag]
		}
		if fn := f.unmarshal; fn != nil {
			var err error
			b, err = fn(b, m.offset(f.field), wire)
			if err == nil {
				reqMask |= f.reqMask
				continue
			}
			if r, ok := err.(*RequiredNotSetError); ok {
				
				
				rnse = r
				reqMask |= f.reqMask
				continue
			}
			if err != errInternalBadWireType {
				return err
			}
			
		}

		
		if !u.unrecognized.IsValid() {
			
			var err error
			b, err = skipField(b, wire)
			if err != nil {
				return err
			}
			continue
		}
		
		
		z := m.offset(u.unrecognized).toBytes()
		var emap map[int32]Extension
		var e Extension
		for _, r := range u.extensionRanges {
			if uint64(r.Start) <= tag && tag <= uint64(r.End) {
				if u.extensions.IsValid() {
					mp := m.offset(u.extensions).toExtensions()
					emap = mp.extensionsWrite()
					e = emap[int32(tag)]
					z = &e.enc
					break
				}
				if u.oldExtensions.IsValid() {
					p := m.offset(u.oldExtensions).toOldExtensions()
					emap = *p
					if emap == nil {
						emap = map[int32]Extension{}
						*p = emap
					}
					e = emap[int32(tag)]
					z = &e.enc
					break
				}
				panic("no extensions field available")
			}
		}

		
		var err error
		b0 := b
		b, err = skipField(b, wire)
		if err != nil {
			return err
		}
		*z = encodeVarint(*z, tag<<3|uint64(wire))
		*z = append(*z, b0[:len(b0)-len(b)]...)

		if emap != nil {
			emap[int32(tag)] = e
		}
	}
	if rnse != nil {
		
		return rnse
	}
	if reqMask != u.reqMask {
		
		for _, n := range u.reqFields {
			if reqMask&1 == 0 {
				return &RequiredNotSetError{n}
			}
			reqMask >>= 1
		}
	}
	return nil
}



func (u *unmarshalInfo) computeUnmarshalInfo() {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.initialized != 0 {
		return
	}
	t := u.typ
	n := t.NumField()

	
	
	u.unrecognized = invalidField
	u.extensions = invalidField
	u.oldExtensions = invalidField

	
	type oneofField struct {
		ityp  reflect.Type 
		field field        
	}
	var oneofFields []oneofField

	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.Name == "XXX_unrecognized" {
			
			if f.Type != reflect.TypeOf(([]byte)(nil)) {
				panic("bad type for XXX_unrecognized field: " + f.Type.Name())
			}
			u.unrecognized = toField(&f)
			continue
		}
		if f.Name == "XXX_InternalExtensions" {
			
			if f.Type != reflect.TypeOf(XXX_InternalExtensions{}) {
				panic("bad type for XXX_InternalExtensions field: " + f.Type.Name())
			}
			u.extensions = toField(&f)
			if f.Tag.Get("protobuf_messageset") == "1" {
				u.isMessageSet = true
			}
			continue
		}
		if f.Name == "XXX_extensions" {
			
			if f.Type != reflect.TypeOf((map[int32]Extension)(nil)) {
				panic("bad type for XXX_extensions field: " + f.Type.Name())
			}
			u.oldExtensions = toField(&f)
			continue
		}
		if f.Name == "XXX_NoUnkeyedLiteral" || f.Name == "XXX_sizecache" {
			continue
		}

		oneof := f.Tag.Get("protobuf_oneof")
		if oneof != "" {
			oneofFields = append(oneofFields, oneofField{f.Type, toField(&f)})
			
			continue
		}

		tags := f.Tag.Get("protobuf")
		tagArray := strings.Split(tags, ",")
		if len(tagArray) < 2 {
			panic("protobuf tag not enough fields in " + t.Name() + "." + f.Name + ": " + tags)
		}
		tag, err := strconv.Atoi(tagArray[1])
		if err != nil {
			panic("protobuf tag field not an integer: " + tagArray[1])
		}

		name := ""
		for _, tag := range tagArray[3:] {
			if strings.HasPrefix(tag, "name=") {
				name = tag[5:]
			}
		}

		
		unmarshal := fieldUnmarshaler(&f)

		
		var reqMask uint64
		if tagArray[2] == "req" {
			bit := len(u.reqFields)
			u.reqFields = append(u.reqFields, name)
			reqMask = uint64(1) << uint(bit)
			
			
			
		}

		
		u.setTag(tag, toField(&f), unmarshal, reqMask)
	}

	
	
	fn := reflect.Zero(reflect.PtrTo(t)).MethodByName("XXX_OneofFuncs")
	if fn.IsValid() {
		res := fn.Call(nil)[3] 
		for i := res.Len() - 1; i >= 0; i-- {
			v := res.Index(i)                             
			tptr := reflect.ValueOf(v.Interface()).Type() 
			typ := tptr.Elem()                            

			f := typ.Field(0) 
			baseUnmarshal := fieldUnmarshaler(&f)
			tagstr := strings.Split(f.Tag.Get("protobuf"), ",")[1]
			tag, err := strconv.Atoi(tagstr)
			if err != nil {
				panic("protobuf tag field not an integer: " + tagstr)
			}

			
			
			for _, of := range oneofFields {
				if tptr.Implements(of.ityp) {
					
					
					
					unmarshal := makeUnmarshalOneof(typ, of.ityp, baseUnmarshal)
					u.setTag(tag, of.field, unmarshal, 0)
				}
			}
		}
	}

	
	fn = reflect.Zero(reflect.PtrTo(t)).MethodByName("ExtensionRangeArray")
	if fn.IsValid() {
		if !u.extensions.IsValid() && !u.oldExtensions.IsValid() {
			panic("a message with extensions, but no extensions field in " + t.Name())
		}
		u.extensionRanges = fn.Call(nil)[0].Interface().([]ExtensionRange)
	}

	
	
	
	
	u.setTag(0, zeroField, func(b []byte, f pointer, w int) ([]byte, error) {
		return nil, fmt.Errorf("proto: %s: illegal tag 0 (wire type %d)", t, w)
	}, 0)

	
	u.reqMask = uint64(1)<<uint(len(u.reqFields)) - 1

	atomic.StoreInt32(&u.initialized, 1)
}





func (u *unmarshalInfo) setTag(tag int, field field, unmarshal unmarshaler, reqMask uint64) {
	i := unmarshalFieldInfo{field: field, unmarshal: unmarshal, reqMask: reqMask}
	n := u.typ.NumField()
	if tag >= 0 && (tag < 16 || tag < 2*n) { 
		for len(u.dense) <= tag {
			u.dense = append(u.dense, unmarshalFieldInfo{})
		}
		u.dense[tag] = i
		return
	}
	if u.sparse == nil {
		u.sparse = map[uint64]unmarshalFieldInfo{}
	}
	u.sparse[uint64(tag)] = i
}


func fieldUnmarshaler(f *reflect.StructField) unmarshaler {
	if f.Type.Kind() == reflect.Map {
		return makeUnmarshalMap(f)
	}
	return typeUnmarshaler(f.Type, f.Tag.Get("protobuf"))
}


func typeUnmarshaler(t reflect.Type, tags string) unmarshaler {
	tagArray := strings.Split(tags, ",")
	encoding := tagArray[0]
	name := "unknown"
	for _, tag := range tagArray[3:] {
		if strings.HasPrefix(tag, "name=") {
			name = tag[5:]
		}
	}

	
	slice := false
	pointer := false
	if t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8 {
		slice = true
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		pointer = true
		t = t.Elem()
	}

	
	if pointer && slice && t.Kind() != reflect.Struct {
		panic("both pointer and slice for basic type in " + t.Name())
	}

	switch t.Kind() {
	case reflect.Bool:
		if pointer {
			return unmarshalBoolPtr
		}
		if slice {
			return unmarshalBoolSlice
		}
		return unmarshalBoolValue
	case reflect.Int32:
		switch encoding {
		case "fixed32":
			if pointer {
				return unmarshalFixedS32Ptr
			}
			if slice {
				return unmarshalFixedS32Slice
			}
			return unmarshalFixedS32Value
		case "varint":
			
			if pointer {
				return unmarshalInt32Ptr
			}
			if slice {
				return unmarshalInt32Slice
			}
			return unmarshalInt32Value
		case "zigzag32":
			if pointer {
				return unmarshalSint32Ptr
			}
			if slice {
				return unmarshalSint32Slice
			}
			return unmarshalSint32Value
		}
	case reflect.Int64:
		switch encoding {
		case "fixed64":
			if pointer {
				return unmarshalFixedS64Ptr
			}
			if slice {
				return unmarshalFixedS64Slice
			}
			return unmarshalFixedS64Value
		case "varint":
			if pointer {
				return unmarshalInt64Ptr
			}
			if slice {
				return unmarshalInt64Slice
			}
			return unmarshalInt64Value
		case "zigzag64":
			if pointer {
				return unmarshalSint64Ptr
			}
			if slice {
				return unmarshalSint64Slice
			}
			return unmarshalSint64Value
		}
	case reflect.Uint32:
		switch encoding {
		case "fixed32":
			if pointer {
				return unmarshalFixed32Ptr
			}
			if slice {
				return unmarshalFixed32Slice
			}
			return unmarshalFixed32Value
		case "varint":
			if pointer {
				return unmarshalUint32Ptr
			}
			if slice {
				return unmarshalUint32Slice
			}
			return unmarshalUint32Value
		}
	case reflect.Uint64:
		switch encoding {
		case "fixed64":
			if pointer {
				return unmarshalFixed64Ptr
			}
			if slice {
				return unmarshalFixed64Slice
			}
			return unmarshalFixed64Value
		case "varint":
			if pointer {
				return unmarshalUint64Ptr
			}
			if slice {
				return unmarshalUint64Slice
			}
			return unmarshalUint64Value
		}
	case reflect.Float32:
		if pointer {
			return unmarshalFloat32Ptr
		}
		if slice {
			return unmarshalFloat32Slice
		}
		return unmarshalFloat32Value
	case reflect.Float64:
		if pointer {
			return unmarshalFloat64Ptr
		}
		if slice {
			return unmarshalFloat64Slice
		}
		return unmarshalFloat64Value
	case reflect.Map:
		panic("map type in typeUnmarshaler in " + t.Name())
	case reflect.Slice:
		if pointer {
			panic("bad pointer in slice case in " + t.Name())
		}
		if slice {
			return unmarshalBytesSlice
		}
		return unmarshalBytesValue
	case reflect.String:
		if pointer {
			return unmarshalStringPtr
		}
		if slice {
			return unmarshalStringSlice
		}
		return unmarshalStringValue
	case reflect.Struct:
		
		if !pointer {
			panic(fmt.Sprintf("message/group field %s:%s without pointer", t, encoding))
		}
		switch encoding {
		case "bytes":
			if slice {
				return makeUnmarshalMessageSlicePtr(getUnmarshalInfo(t), name)
			}
			return makeUnmarshalMessagePtr(getUnmarshalInfo(t), name)
		case "group":
			if slice {
				return makeUnmarshalGroupSlicePtr(getUnmarshalInfo(t), name)
			}
			return makeUnmarshalGroupPtr(getUnmarshalInfo(t), name)
		}
	}
	panic(fmt.Sprintf("unmarshaler not found type:%s encoding:%s", t, encoding))
}



func unmarshalInt64Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int64(x)
	*f.toInt64() = v
	return b, nil
}

func unmarshalInt64Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int64(x)
	*f.toInt64Ptr() = &v
	return b, nil
}

func unmarshalInt64Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			x, n = decodeVarint(b)
			if n == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			b = b[n:]
			v := int64(x)
			s := f.toInt64Slice()
			*s = append(*s, v)
		}
		return res, nil
	}
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int64(x)
	s := f.toInt64Slice()
	*s = append(*s, v)
	return b, nil
}

func unmarshalSint64Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int64(x>>1) ^ int64(x)<<63>>63
	*f.toInt64() = v
	return b, nil
}

func unmarshalSint64Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int64(x>>1) ^ int64(x)<<63>>63
	*f.toInt64Ptr() = &v
	return b, nil
}

func unmarshalSint64Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			x, n = decodeVarint(b)
			if n == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			b = b[n:]
			v := int64(x>>1) ^ int64(x)<<63>>63
			s := f.toInt64Slice()
			*s = append(*s, v)
		}
		return res, nil
	}
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int64(x>>1) ^ int64(x)<<63>>63
	s := f.toInt64Slice()
	*s = append(*s, v)
	return b, nil
}

func unmarshalUint64Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := uint64(x)
	*f.toUint64() = v
	return b, nil
}

func unmarshalUint64Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := uint64(x)
	*f.toUint64Ptr() = &v
	return b, nil
}

func unmarshalUint64Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			x, n = decodeVarint(b)
			if n == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			b = b[n:]
			v := uint64(x)
			s := f.toUint64Slice()
			*s = append(*s, v)
		}
		return res, nil
	}
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := uint64(x)
	s := f.toUint64Slice()
	*s = append(*s, v)
	return b, nil
}

func unmarshalInt32Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int32(x)
	*f.toInt32() = v
	return b, nil
}

func unmarshalInt32Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int32(x)
	f.setInt32Ptr(v)
	return b, nil
}

func unmarshalInt32Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			x, n = decodeVarint(b)
			if n == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			b = b[n:]
			v := int32(x)
			f.appendInt32Slice(v)
		}
		return res, nil
	}
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int32(x)
	f.appendInt32Slice(v)
	return b, nil
}

func unmarshalSint32Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int32(x>>1) ^ int32(x)<<31>>31
	*f.toInt32() = v
	return b, nil
}

func unmarshalSint32Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int32(x>>1) ^ int32(x)<<31>>31
	f.setInt32Ptr(v)
	return b, nil
}

func unmarshalSint32Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			x, n = decodeVarint(b)
			if n == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			b = b[n:]
			v := int32(x>>1) ^ int32(x)<<31>>31
			f.appendInt32Slice(v)
		}
		return res, nil
	}
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := int32(x>>1) ^ int32(x)<<31>>31
	f.appendInt32Slice(v)
	return b, nil
}

func unmarshalUint32Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := uint32(x)
	*f.toUint32() = v
	return b, nil
}

func unmarshalUint32Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := uint32(x)
	*f.toUint32Ptr() = &v
	return b, nil
}

func unmarshalUint32Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			x, n = decodeVarint(b)
			if n == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			b = b[n:]
			v := uint32(x)
			s := f.toUint32Slice()
			*s = append(*s, v)
		}
		return res, nil
	}
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	v := uint32(x)
	s := f.toUint32Slice()
	*s = append(*s, v)
	return b, nil
}

func unmarshalFixed64Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	*f.toUint64() = v
	return b[8:], nil
}

func unmarshalFixed64Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	*f.toUint64Ptr() = &v
	return b[8:], nil
}

func unmarshalFixed64Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			if len(b) < 8 {
				return nil, io.ErrUnexpectedEOF
			}
			v := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
			s := f.toUint64Slice()
			*s = append(*s, v)
			b = b[8:]
		}
		return res, nil
	}
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	s := f.toUint64Slice()
	*s = append(*s, v)
	return b[8:], nil
}

func unmarshalFixedS64Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 | int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56
	*f.toInt64() = v
	return b[8:], nil
}

func unmarshalFixedS64Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 | int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56
	*f.toInt64Ptr() = &v
	return b[8:], nil
}

func unmarshalFixedS64Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			if len(b) < 8 {
				return nil, io.ErrUnexpectedEOF
			}
			v := int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 | int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56
			s := f.toInt64Slice()
			*s = append(*s, v)
			b = b[8:]
		}
		return res, nil
	}
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 | int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56
	s := f.toInt64Slice()
	*s = append(*s, v)
	return b[8:], nil
}

func unmarshalFixed32Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	*f.toUint32() = v
	return b[4:], nil
}

func unmarshalFixed32Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	*f.toUint32Ptr() = &v
	return b[4:], nil
}

func unmarshalFixed32Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			if len(b) < 4 {
				return nil, io.ErrUnexpectedEOF
			}
			v := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
			s := f.toUint32Slice()
			*s = append(*s, v)
			b = b[4:]
		}
		return res, nil
	}
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	s := f.toUint32Slice()
	*s = append(*s, v)
	return b[4:], nil
}

func unmarshalFixedS32Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
	*f.toInt32() = v
	return b[4:], nil
}

func unmarshalFixedS32Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
	f.setInt32Ptr(v)
	return b[4:], nil
}

func unmarshalFixedS32Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			if len(b) < 4 {
				return nil, io.ErrUnexpectedEOF
			}
			v := int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
			f.appendInt32Slice(v)
			b = b[4:]
		}
		return res, nil
	}
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
	f.appendInt32Slice(v)
	return b[4:], nil
}

func unmarshalBoolValue(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	
	
	
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	
	v := x != 0
	*f.toBool() = v
	return b[n:], nil
}

func unmarshalBoolPtr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	v := x != 0
	*f.toBoolPtr() = &v
	return b[n:], nil
}

func unmarshalBoolSlice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			x, n = decodeVarint(b)
			if n == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			v := x != 0
			s := f.toBoolSlice()
			*s = append(*s, v)
			b = b[n:]
		}
		return res, nil
	}
	if w != WireVarint {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	v := x != 0
	s := f.toBoolSlice()
	*s = append(*s, v)
	return b[n:], nil
}

func unmarshalFloat64Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := math.Float64frombits(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
	*f.toFloat64() = v
	return b[8:], nil
}

func unmarshalFloat64Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := math.Float64frombits(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
	*f.toFloat64Ptr() = &v
	return b[8:], nil
}

func unmarshalFloat64Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			if len(b) < 8 {
				return nil, io.ErrUnexpectedEOF
			}
			v := math.Float64frombits(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
			s := f.toFloat64Slice()
			*s = append(*s, v)
			b = b[8:]
		}
		return res, nil
	}
	if w != WireFixed64 {
		return b, errInternalBadWireType
	}
	if len(b) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	v := math.Float64frombits(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
	s := f.toFloat64Slice()
	*s = append(*s, v)
	return b[8:], nil
}

func unmarshalFloat32Value(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := math.Float32frombits(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
	*f.toFloat32() = v
	return b[4:], nil
}

func unmarshalFloat32Ptr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := math.Float32frombits(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
	*f.toFloat32Ptr() = &v
	return b[4:], nil
}

func unmarshalFloat32Slice(b []byte, f pointer, w int) ([]byte, error) {
	if w == WireBytes { 
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		res := b[x:]
		b = b[:x]
		for len(b) > 0 {
			if len(b) < 4 {
				return nil, io.ErrUnexpectedEOF
			}
			v := math.Float32frombits(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
			s := f.toFloat32Slice()
			*s = append(*s, v)
			b = b[4:]
		}
		return res, nil
	}
	if w != WireFixed32 {
		return b, errInternalBadWireType
	}
	if len(b) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	v := math.Float32frombits(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
	s := f.toFloat32Slice()
	*s = append(*s, v)
	return b[4:], nil
}

func unmarshalStringValue(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireBytes {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	if x > uint64(len(b)) {
		return nil, io.ErrUnexpectedEOF
	}
	v := string(b[:x])
	if !utf8.ValidString(v) {
		return nil, errInvalidUTF8
	}
	*f.toString() = v
	return b[x:], nil
}

func unmarshalStringPtr(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireBytes {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	if x > uint64(len(b)) {
		return nil, io.ErrUnexpectedEOF
	}
	v := string(b[:x])
	if !utf8.ValidString(v) {
		return nil, errInvalidUTF8
	}
	*f.toStringPtr() = &v
	return b[x:], nil
}

func unmarshalStringSlice(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireBytes {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	if x > uint64(len(b)) {
		return nil, io.ErrUnexpectedEOF
	}
	v := string(b[:x])
	if !utf8.ValidString(v) {
		return nil, errInvalidUTF8
	}
	s := f.toStringSlice()
	*s = append(*s, v)
	return b[x:], nil
}

var emptyBuf [0]byte

func unmarshalBytesValue(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireBytes {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	if x > uint64(len(b)) {
		return nil, io.ErrUnexpectedEOF
	}
	
	
	
	
	v := append(emptyBuf[:], b[:x]...)
	*f.toBytes() = v
	return b[x:], nil
}

func unmarshalBytesSlice(b []byte, f pointer, w int) ([]byte, error) {
	if w != WireBytes {
		return b, errInternalBadWireType
	}
	x, n := decodeVarint(b)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	b = b[n:]
	if x > uint64(len(b)) {
		return nil, io.ErrUnexpectedEOF
	}
	v := append(emptyBuf[:], b[:x]...)
	s := f.toBytesSlice()
	*s = append(*s, v)
	return b[x:], nil
}

func makeUnmarshalMessagePtr(sub *unmarshalInfo, name string) unmarshaler {
	return func(b []byte, f pointer, w int) ([]byte, error) {
		if w != WireBytes {
			return b, errInternalBadWireType
		}
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		
		
		
		
		v := f.getPointer()
		if v.isNil() {
			v = valToPointer(reflect.New(sub.typ))
			f.setPointer(v)
		}
		err := sub.unmarshal(v, b[:x])
		if err != nil {
			if r, ok := err.(*RequiredNotSetError); ok {
				r.field = name + "." + r.field
			} else {
				return nil, err
			}
		}
		return b[x:], err
	}
}

func makeUnmarshalMessageSlicePtr(sub *unmarshalInfo, name string) unmarshaler {
	return func(b []byte, f pointer, w int) ([]byte, error) {
		if w != WireBytes {
			return b, errInternalBadWireType
		}
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		v := valToPointer(reflect.New(sub.typ))
		err := sub.unmarshal(v, b[:x])
		if err != nil {
			if r, ok := err.(*RequiredNotSetError); ok {
				r.field = name + "." + r.field
			} else {
				return nil, err
			}
		}
		f.appendPointer(v)
		return b[x:], err
	}
}

func makeUnmarshalGroupPtr(sub *unmarshalInfo, name string) unmarshaler {
	return func(b []byte, f pointer, w int) ([]byte, error) {
		if w != WireStartGroup {
			return b, errInternalBadWireType
		}
		x, y := findEndGroup(b)
		if x < 0 {
			return nil, io.ErrUnexpectedEOF
		}
		v := f.getPointer()
		if v.isNil() {
			v = valToPointer(reflect.New(sub.typ))
			f.setPointer(v)
		}
		err := sub.unmarshal(v, b[:x])
		if err != nil {
			if r, ok := err.(*RequiredNotSetError); ok {
				r.field = name + "." + r.field
			} else {
				return nil, err
			}
		}
		return b[y:], err
	}
}

func makeUnmarshalGroupSlicePtr(sub *unmarshalInfo, name string) unmarshaler {
	return func(b []byte, f pointer, w int) ([]byte, error) {
		if w != WireStartGroup {
			return b, errInternalBadWireType
		}
		x, y := findEndGroup(b)
		if x < 0 {
			return nil, io.ErrUnexpectedEOF
		}
		v := valToPointer(reflect.New(sub.typ))
		err := sub.unmarshal(v, b[:x])
		if err != nil {
			if r, ok := err.(*RequiredNotSetError); ok {
				r.field = name + "." + r.field
			} else {
				return nil, err
			}
		}
		f.appendPointer(v)
		return b[y:], err
	}
}

func makeUnmarshalMap(f *reflect.StructField) unmarshaler {
	t := f.Type
	kt := t.Key()
	vt := t.Elem()
	unmarshalKey := typeUnmarshaler(kt, f.Tag.Get("protobuf_key"))
	unmarshalVal := typeUnmarshaler(vt, f.Tag.Get("protobuf_val"))
	return func(b []byte, f pointer, w int) ([]byte, error) {
		
		if w != WireBytes {
			return nil, fmt.Errorf("proto: bad wiretype for map field: got %d want %d", w, WireBytes)
		}
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		if x > uint64(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}
		r := b[x:] 
		b = b[:x]  

		
		
		

		
		k := reflect.New(kt)
		v := reflect.New(vt)
		for len(b) > 0 {
			x, n := decodeVarint(b)
			if n == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			wire := int(x) & 7
			b = b[n:]

			var err error
			switch x >> 3 {
			case 1:
				b, err = unmarshalKey(b, valToPointer(k), wire)
			case 2:
				b, err = unmarshalVal(b, valToPointer(v), wire)
			default:
				err = errInternalBadWireType 
			}

			if err == nil {
				continue
			}
			if err != errInternalBadWireType {
				return nil, err
			}

			
			b, err = skipField(b, wire)
			if err != nil {
				return nil, err
			}
		}

		
		m := f.asPointerTo(t).Elem() 
		if m.IsNil() {
			m.Set(reflect.MakeMap(t))
		}

		
		m.SetMapIndex(k.Elem(), v.Elem())

		return r, nil
	}
}













func makeUnmarshalOneof(typ, ityp reflect.Type, unmarshal unmarshaler) unmarshaler {
	sf := typ.Field(0)
	field0 := toField(&sf)
	return func(b []byte, f pointer, w int) ([]byte, error) {
		
		v := reflect.New(typ)

		
		
		var err error
		b, err = unmarshal(b, valToPointer(v).offset(field0), w)
		if err != nil {
			return nil, err
		}

		
		f.asPointerTo(ityp).Elem().Set(v)

		return b, nil
	}
}


var errInternalBadWireType = errors.New("proto: internal error: bad wiretype")


func skipField(b []byte, wire int) ([]byte, error) {
	switch wire {
	case WireVarint:
		_, k := decodeVarint(b)
		if k == 0 {
			return b, io.ErrUnexpectedEOF
		}
		b = b[k:]
	case WireFixed32:
		if len(b) < 4 {
			return b, io.ErrUnexpectedEOF
		}
		b = b[4:]
	case WireFixed64:
		if len(b) < 8 {
			return b, io.ErrUnexpectedEOF
		}
		b = b[8:]
	case WireBytes:
		m, k := decodeVarint(b)
		if k == 0 || uint64(len(b)-k) < m {
			return b, io.ErrUnexpectedEOF
		}
		b = b[uint64(k)+m:]
	case WireStartGroup:
		_, i := findEndGroup(b)
		if i == -1 {
			return b, io.ErrUnexpectedEOF
		}
		b = b[i:]
	default:
		return b, fmt.Errorf("proto: can't skip unknown wire type %d", wire)
	}
	return b, nil
}






func findEndGroup(b []byte) (int, int) {
	depth := 1
	i := 0
	for {
		x, n := decodeVarint(b[i:])
		if n == 0 {
			return -1, -1
		}
		j := i
		i += n
		switch x & 7 {
		case WireVarint:
			_, k := decodeVarint(b[i:])
			if k == 0 {
				return -1, -1
			}
			i += k
		case WireFixed32:
			if len(b)-4 < i {
				return -1, -1
			}
			i += 4
		case WireFixed64:
			if len(b)-8 < i {
				return -1, -1
			}
			i += 8
		case WireBytes:
			m, k := decodeVarint(b[i:])
			if k == 0 {
				return -1, -1
			}
			i += k
			if uint64(len(b)-i) < m {
				return -1, -1
			}
			i += int(m)
		case WireStartGroup:
			depth++
		case WireEndGroup:
			depth--
			if depth == 0 {
				return j, i
			}
		default:
			return -1, -1
		}
	}
}


func encodeVarint(b []byte, x uint64) []byte {
	for x >= 1<<7 {
		b = append(b, byte(x&0x7f|0x80))
		x >>= 7
	}
	return append(b, byte(x))
}




func decodeVarint(b []byte) (uint64, int) {
	var x, y uint64
	if len(b) <= 0 {
		goto bad
	}
	x = uint64(b[0])
	if x < 0x80 {
		return x, 1
	}
	x -= 0x80

	if len(b) <= 1 {
		goto bad
	}
	y = uint64(b[1])
	x += y << 7
	if y < 0x80 {
		return x, 2
	}
	x -= 0x80 << 7

	if len(b) <= 2 {
		goto bad
	}
	y = uint64(b[2])
	x += y << 14
	if y < 0x80 {
		return x, 3
	}
	x -= 0x80 << 14

	if len(b) <= 3 {
		goto bad
	}
	y = uint64(b[3])
	x += y << 21
	if y < 0x80 {
		return x, 4
	}
	x -= 0x80 << 21

	if len(b) <= 4 {
		goto bad
	}
	y = uint64(b[4])
	x += y << 28
	if y < 0x80 {
		return x, 5
	}
	x -= 0x80 << 28

	if len(b) <= 5 {
		goto bad
	}
	y = uint64(b[5])
	x += y << 35
	if y < 0x80 {
		return x, 6
	}
	x -= 0x80 << 35

	if len(b) <= 6 {
		goto bad
	}
	y = uint64(b[6])
	x += y << 42
	if y < 0x80 {
		return x, 7
	}
	x -= 0x80 << 42

	if len(b) <= 7 {
		goto bad
	}
	y = uint64(b[7])
	x += y << 49
	if y < 0x80 {
		return x, 8
	}
	x -= 0x80 << 49

	if len(b) <= 8 {
		goto bad
	}
	y = uint64(b[8])
	x += y << 56
	if y < 0x80 {
		return x, 9
	}
	x -= 0x80 << 56

	if len(b) <= 9 {
		goto bad
	}
	y = uint64(b[9])
	x += y << 63
	if y < 2 {
		return x, 10
	}

bad:
	return 0, 0
}
