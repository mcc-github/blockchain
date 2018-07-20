






























package proto

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)



type sizer func(pointer, int) int



type marshaler func(b []byte, ptr pointer, wiretag uint64, deterministic bool) ([]byte, error)


type marshalInfo struct {
	typ          reflect.Type
	fields       []*marshalFieldInfo
	unrecognized field                      
	extensions   field                      
	v1extensions field                      
	sizecache    field                      
	initialized  int32                      
	messageset   bool                       
	hasmarshaler bool                       
	sync.RWMutex                            
	extElems     map[int32]*marshalElemInfo 
}


type marshalFieldInfo struct {
	field      field
	wiretag    uint64 
	tagsize    int    
	sizer      sizer
	marshaler  marshaler
	isPointer  bool
	required   bool                              
	name       string                            
	oneofElems map[reflect.Type]*marshalElemInfo 
}


type marshalElemInfo struct {
	wiretag   uint64 
	tagsize   int    
	sizer     sizer
	marshaler marshaler
	isptr     bool 
}

var (
	marshalInfoMap  = map[reflect.Type]*marshalInfo{}
	marshalInfoLock sync.Mutex
)




func getMarshalInfo(t reflect.Type) *marshalInfo {
	marshalInfoLock.Lock()
	u, ok := marshalInfoMap[t]
	if !ok {
		u = &marshalInfo{typ: t}
		marshalInfoMap[t] = u
	}
	marshalInfoLock.Unlock()
	return u
}





func (a *InternalMessageInfo) Size(msg Message) int {
	u := getMessageMarshalInfo(msg, a)
	ptr := toPointer(&msg)
	if ptr.isNil() {
		
		
		
		return 0
	}
	return u.size(ptr)
}





func (a *InternalMessageInfo) Marshal(b []byte, msg Message, deterministic bool) ([]byte, error) {
	u := getMessageMarshalInfo(msg, a)
	ptr := toPointer(&msg)
	if ptr.isNil() {
		
		
		
		return b, ErrNil
	}
	return u.marshal(b, ptr, deterministic)
}

func getMessageMarshalInfo(msg interface{}, a *InternalMessageInfo) *marshalInfo {
	
	
	u := atomicLoadMarshalInfo(&a.marshal)
	if u == nil {
		
		t := reflect.ValueOf(msg).Type()
		if t.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("cannot handle non-pointer message type %v", t))
		}
		u = getMarshalInfo(t.Elem())
		
		
		atomicStoreMarshalInfo(&a.marshal, u)
	}
	return u
}



func (u *marshalInfo) size(ptr pointer) int {
	if atomic.LoadInt32(&u.initialized) == 0 {
		u.computeMarshalInfo()
	}

	
	
	if u.hasmarshaler {
		m := ptr.asPointerTo(u.typ).Interface().(Marshaler)
		b, _ := m.Marshal()
		return len(b)
	}

	n := 0
	for _, f := range u.fields {
		if f.isPointer && ptr.offset(f.field).getPointer().isNil() {
			
			continue
		}
		n += f.sizer(ptr.offset(f.field), f.tagsize)
	}
	if u.extensions.IsValid() {
		e := ptr.offset(u.extensions).toExtensions()
		if u.messageset {
			n += u.sizeMessageSet(e)
		} else {
			n += u.sizeExtensions(e)
		}
	}
	if u.v1extensions.IsValid() {
		m := *ptr.offset(u.v1extensions).toOldExtensions()
		n += u.sizeV1Extensions(m)
	}
	if u.unrecognized.IsValid() {
		s := *ptr.offset(u.unrecognized).toBytes()
		n += len(s)
	}
	
	if u.sizecache.IsValid() {
		atomic.StoreInt32(ptr.offset(u.sizecache).toInt32(), int32(n))
	}
	return n
}



func (u *marshalInfo) cachedsize(ptr pointer) int {
	if u.sizecache.IsValid() {
		return int(atomic.LoadInt32(ptr.offset(u.sizecache).toInt32()))
	}
	return u.size(ptr)
}





func (u *marshalInfo) marshal(b []byte, ptr pointer, deterministic bool) ([]byte, error) {
	if atomic.LoadInt32(&u.initialized) == 0 {
		u.computeMarshalInfo()
	}

	
	
	if u.hasmarshaler {
		m := ptr.asPointerTo(u.typ).Interface().(Marshaler)
		b1, err := m.Marshal()
		b = append(b, b1...)
		return b, err
	}

	var err, errreq error
	
	if u.extensions.IsValid() {
		e := ptr.offset(u.extensions).toExtensions()
		if u.messageset {
			b, err = u.appendMessageSet(b, e, deterministic)
		} else {
			b, err = u.appendExtensions(b, e, deterministic)
		}
		if err != nil {
			return b, err
		}
	}
	if u.v1extensions.IsValid() {
		m := *ptr.offset(u.v1extensions).toOldExtensions()
		b, err = u.appendV1Extensions(b, m, deterministic)
		if err != nil {
			return b, err
		}
	}
	for _, f := range u.fields {
		if f.required && errreq == nil {
			if ptr.offset(f.field).getPointer().isNil() {
				
				
				errreq = &RequiredNotSetError{f.name}
				continue
			}
		}
		if f.isPointer && ptr.offset(f.field).getPointer().isNil() {
			
			continue
		}
		b, err = f.marshaler(b, ptr.offset(f.field), f.wiretag, deterministic)
		if err != nil {
			if err1, ok := err.(*RequiredNotSetError); ok {
				
				
				if errreq == nil {
					errreq = &RequiredNotSetError{f.name + "." + err1.field}
				}
				continue
			}
			if err == errRepeatedHasNil {
				err = errors.New("proto: repeated field " + f.name + " has nil element")
			}
			return b, err
		}
	}
	if u.unrecognized.IsValid() {
		s := *ptr.offset(u.unrecognized).toBytes()
		b = append(b, s...)
	}
	return b, errreq
}


func (u *marshalInfo) computeMarshalInfo() {
	u.Lock()
	defer u.Unlock()
	if u.initialized != 0 { 
		return
	}

	t := u.typ
	u.unrecognized = invalidField
	u.extensions = invalidField
	u.v1extensions = invalidField
	u.sizecache = invalidField

	
	
	if reflect.PtrTo(t).Implements(marshalerType) {
		u.hasmarshaler = true
		atomic.StoreInt32(&u.initialized, 1)
		return
	}

	
	var oneofImplementers []interface{}
	if m, ok := reflect.Zero(reflect.PtrTo(t)).Interface().(oneofMessage); ok {
		_, _, _, oneofImplementers = m.XXX_OneofFuncs()
	}

	n := t.NumField()

	
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !strings.HasPrefix(f.Name, "XXX_") {
			continue
		}
		switch f.Name {
		case "XXX_sizecache":
			u.sizecache = toField(&f)
		case "XXX_unrecognized":
			u.unrecognized = toField(&f)
		case "XXX_InternalExtensions":
			u.extensions = toField(&f)
			u.messageset = f.Tag.Get("protobuf_messageset") == "1"
		case "XXX_extensions":
			u.v1extensions = toField(&f)
		case "XXX_NoUnkeyedLiteral":
			
		default:
			panic("unknown XXX field: " + f.Name)
		}
		n--
	}

	
	fields := make([]marshalFieldInfo, n) 
	u.fields = make([]*marshalFieldInfo, 0, n)
	for i, j := 0, 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if strings.HasPrefix(f.Name, "XXX_") {
			continue
		}
		field := &fields[j]
		j++
		field.name = f.Name
		u.fields = append(u.fields, field)
		if f.Tag.Get("protobuf_oneof") != "" {
			field.computeOneofFieldInfo(&f, oneofImplementers)
			continue
		}
		if f.Tag.Get("protobuf") == "" {
			
			u.fields = u.fields[:len(u.fields)-1]
			j--
			continue
		}
		field.computeMarshalFieldInfo(&f)
	}

	
	sort.Sort(byTag(u.fields))

	atomic.StoreInt32(&u.initialized, 1)
}


type byTag []*marshalFieldInfo

func (a byTag) Len() int           { return len(a) }
func (a byTag) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTag) Less(i, j int) bool { return a[i].wiretag < a[j].wiretag }



func (u *marshalInfo) getExtElemInfo(desc *ExtensionDesc) *marshalElemInfo {
	
	u.RLock()
	e, ok := u.extElems[desc.Field]
	u.RUnlock()
	if ok {
		return e
	}

	t := reflect.TypeOf(desc.ExtensionType) 
	tags := strings.Split(desc.Tag, ",")
	tag, err := strconv.Atoi(tags[1])
	if err != nil {
		panic("tag is not an integer")
	}
	wt := wiretype(tags[0])
	sizer, marshaler := typeMarshaler(t, tags, false, false)
	e = &marshalElemInfo{
		wiretag:   uint64(tag)<<3 | wt,
		tagsize:   SizeVarint(uint64(tag) << 3),
		sizer:     sizer,
		marshaler: marshaler,
		isptr:     t.Kind() == reflect.Ptr,
	}

	
	u.Lock()
	if u.extElems == nil {
		u.extElems = make(map[int32]*marshalElemInfo)
	}
	u.extElems[desc.Field] = e
	u.Unlock()
	return e
}


func (fi *marshalFieldInfo) computeMarshalFieldInfo(f *reflect.StructField) {
	
	
	tags := strings.Split(f.Tag.Get("protobuf"), ",")
	if tags[0] == "" {
		return
	}
	tag, err := strconv.Atoi(tags[1])
	if err != nil {
		panic("tag is not an integer")
	}
	wt := wiretype(tags[0])
	if tags[2] == "req" {
		fi.required = true
	}
	fi.setTag(f, tag, wt)
	fi.setMarshaler(f, tags)
}

func (fi *marshalFieldInfo) computeOneofFieldInfo(f *reflect.StructField, oneofImplementers []interface{}) {
	fi.field = toField(f)
	fi.wiretag = 1<<31 - 1 
	fi.isPointer = true
	fi.sizer, fi.marshaler = makeOneOfMarshaler(fi, f)
	fi.oneofElems = make(map[reflect.Type]*marshalElemInfo)

	ityp := f.Type 
	for _, o := range oneofImplementers {
		t := reflect.TypeOf(o)
		if !t.Implements(ityp) {
			continue
		}
		sf := t.Elem().Field(0) 
		tags := strings.Split(sf.Tag.Get("protobuf"), ",")
		tag, err := strconv.Atoi(tags[1])
		if err != nil {
			panic("tag is not an integer")
		}
		wt := wiretype(tags[0])
		sizer, marshaler := typeMarshaler(sf.Type, tags, false, true) 
		fi.oneofElems[t.Elem()] = &marshalElemInfo{
			wiretag:   uint64(tag)<<3 | wt,
			tagsize:   SizeVarint(uint64(tag) << 3),
			sizer:     sizer,
			marshaler: marshaler,
		}
	}
}

type oneofMessage interface {
	XXX_OneofFuncs() (func(Message, *Buffer) error, func(Message, int, int, *Buffer) (bool, error), func(Message) int, []interface{})
}


func wiretype(encoding string) uint64 {
	switch encoding {
	case "fixed32":
		return WireFixed32
	case "fixed64":
		return WireFixed64
	case "varint", "zigzag32", "zigzag64":
		return WireVarint
	case "bytes":
		return WireBytes
	case "group":
		return WireStartGroup
	}
	panic("unknown wire type " + encoding)
}


func (fi *marshalFieldInfo) setTag(f *reflect.StructField, tag int, wt uint64) {
	fi.field = toField(f)
	fi.wiretag = uint64(tag)<<3 | wt
	fi.tagsize = SizeVarint(uint64(tag) << 3)
}


func (fi *marshalFieldInfo) setMarshaler(f *reflect.StructField, tags []string) {
	switch f.Type.Kind() {
	case reflect.Map:
		
		fi.isPointer = true
		fi.sizer, fi.marshaler = makeMapMarshaler(f)
		return
	case reflect.Ptr, reflect.Slice:
		fi.isPointer = true
	}
	fi.sizer, fi.marshaler = typeMarshaler(f.Type, tags, true, false)
}






func typeMarshaler(t reflect.Type, tags []string, nozero, oneof bool) (sizer, marshaler) {
	encoding := tags[0]

	pointer := false
	slice := false
	if t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8 {
		slice = true
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		pointer = true
		t = t.Elem()
	}

	packed := false
	proto3 := false
	for i := 2; i < len(tags); i++ {
		if tags[i] == "packed" {
			packed = true
		}
		if tags[i] == "proto3" {
			proto3 = true
		}
	}

	switch t.Kind() {
	case reflect.Bool:
		if pointer {
			return sizeBoolPtr, appendBoolPtr
		}
		if slice {
			if packed {
				return sizeBoolPackedSlice, appendBoolPackedSlice
			}
			return sizeBoolSlice, appendBoolSlice
		}
		if nozero {
			return sizeBoolValueNoZero, appendBoolValueNoZero
		}
		return sizeBoolValue, appendBoolValue
	case reflect.Uint32:
		switch encoding {
		case "fixed32":
			if pointer {
				return sizeFixed32Ptr, appendFixed32Ptr
			}
			if slice {
				if packed {
					return sizeFixed32PackedSlice, appendFixed32PackedSlice
				}
				return sizeFixed32Slice, appendFixed32Slice
			}
			if nozero {
				return sizeFixed32ValueNoZero, appendFixed32ValueNoZero
			}
			return sizeFixed32Value, appendFixed32Value
		case "varint":
			if pointer {
				return sizeVarint32Ptr, appendVarint32Ptr
			}
			if slice {
				if packed {
					return sizeVarint32PackedSlice, appendVarint32PackedSlice
				}
				return sizeVarint32Slice, appendVarint32Slice
			}
			if nozero {
				return sizeVarint32ValueNoZero, appendVarint32ValueNoZero
			}
			return sizeVarint32Value, appendVarint32Value
		}
	case reflect.Int32:
		switch encoding {
		case "fixed32":
			if pointer {
				return sizeFixedS32Ptr, appendFixedS32Ptr
			}
			if slice {
				if packed {
					return sizeFixedS32PackedSlice, appendFixedS32PackedSlice
				}
				return sizeFixedS32Slice, appendFixedS32Slice
			}
			if nozero {
				return sizeFixedS32ValueNoZero, appendFixedS32ValueNoZero
			}
			return sizeFixedS32Value, appendFixedS32Value
		case "varint":
			if pointer {
				return sizeVarintS32Ptr, appendVarintS32Ptr
			}
			if slice {
				if packed {
					return sizeVarintS32PackedSlice, appendVarintS32PackedSlice
				}
				return sizeVarintS32Slice, appendVarintS32Slice
			}
			if nozero {
				return sizeVarintS32ValueNoZero, appendVarintS32ValueNoZero
			}
			return sizeVarintS32Value, appendVarintS32Value
		case "zigzag32":
			if pointer {
				return sizeZigzag32Ptr, appendZigzag32Ptr
			}
			if slice {
				if packed {
					return sizeZigzag32PackedSlice, appendZigzag32PackedSlice
				}
				return sizeZigzag32Slice, appendZigzag32Slice
			}
			if nozero {
				return sizeZigzag32ValueNoZero, appendZigzag32ValueNoZero
			}
			return sizeZigzag32Value, appendZigzag32Value
		}
	case reflect.Uint64:
		switch encoding {
		case "fixed64":
			if pointer {
				return sizeFixed64Ptr, appendFixed64Ptr
			}
			if slice {
				if packed {
					return sizeFixed64PackedSlice, appendFixed64PackedSlice
				}
				return sizeFixed64Slice, appendFixed64Slice
			}
			if nozero {
				return sizeFixed64ValueNoZero, appendFixed64ValueNoZero
			}
			return sizeFixed64Value, appendFixed64Value
		case "varint":
			if pointer {
				return sizeVarint64Ptr, appendVarint64Ptr
			}
			if slice {
				if packed {
					return sizeVarint64PackedSlice, appendVarint64PackedSlice
				}
				return sizeVarint64Slice, appendVarint64Slice
			}
			if nozero {
				return sizeVarint64ValueNoZero, appendVarint64ValueNoZero
			}
			return sizeVarint64Value, appendVarint64Value
		}
	case reflect.Int64:
		switch encoding {
		case "fixed64":
			if pointer {
				return sizeFixedS64Ptr, appendFixedS64Ptr
			}
			if slice {
				if packed {
					return sizeFixedS64PackedSlice, appendFixedS64PackedSlice
				}
				return sizeFixedS64Slice, appendFixedS64Slice
			}
			if nozero {
				return sizeFixedS64ValueNoZero, appendFixedS64ValueNoZero
			}
			return sizeFixedS64Value, appendFixedS64Value
		case "varint":
			if pointer {
				return sizeVarintS64Ptr, appendVarintS64Ptr
			}
			if slice {
				if packed {
					return sizeVarintS64PackedSlice, appendVarintS64PackedSlice
				}
				return sizeVarintS64Slice, appendVarintS64Slice
			}
			if nozero {
				return sizeVarintS64ValueNoZero, appendVarintS64ValueNoZero
			}
			return sizeVarintS64Value, appendVarintS64Value
		case "zigzag64":
			if pointer {
				return sizeZigzag64Ptr, appendZigzag64Ptr
			}
			if slice {
				if packed {
					return sizeZigzag64PackedSlice, appendZigzag64PackedSlice
				}
				return sizeZigzag64Slice, appendZigzag64Slice
			}
			if nozero {
				return sizeZigzag64ValueNoZero, appendZigzag64ValueNoZero
			}
			return sizeZigzag64Value, appendZigzag64Value
		}
	case reflect.Float32:
		if pointer {
			return sizeFloat32Ptr, appendFloat32Ptr
		}
		if slice {
			if packed {
				return sizeFloat32PackedSlice, appendFloat32PackedSlice
			}
			return sizeFloat32Slice, appendFloat32Slice
		}
		if nozero {
			return sizeFloat32ValueNoZero, appendFloat32ValueNoZero
		}
		return sizeFloat32Value, appendFloat32Value
	case reflect.Float64:
		if pointer {
			return sizeFloat64Ptr, appendFloat64Ptr
		}
		if slice {
			if packed {
				return sizeFloat64PackedSlice, appendFloat64PackedSlice
			}
			return sizeFloat64Slice, appendFloat64Slice
		}
		if nozero {
			return sizeFloat64ValueNoZero, appendFloat64ValueNoZero
		}
		return sizeFloat64Value, appendFloat64Value
	case reflect.String:
		if pointer {
			return sizeStringPtr, appendStringPtr
		}
		if slice {
			return sizeStringSlice, appendStringSlice
		}
		if nozero {
			return sizeStringValueNoZero, appendStringValueNoZero
		}
		return sizeStringValue, appendStringValue
	case reflect.Slice:
		if slice {
			return sizeBytesSlice, appendBytesSlice
		}
		if oneof {
			
			
			
			return sizeBytesOneof, appendBytesOneof
		}
		if proto3 {
			return sizeBytes3, appendBytes3
		}
		return sizeBytes, appendBytes
	case reflect.Struct:
		switch encoding {
		case "group":
			if slice {
				return makeGroupSliceMarshaler(getMarshalInfo(t))
			}
			return makeGroupMarshaler(getMarshalInfo(t))
		case "bytes":
			if slice {
				return makeMessageSliceMarshaler(getMarshalInfo(t))
			}
			return makeMessageMarshaler(getMarshalInfo(t))
		}
	}
	panic(fmt.Sprintf("unknown or mismatched type: type: %v, wire type: %v", t, encoding))
}





func sizeFixed32Value(_ pointer, tagsize int) int {
	return 4 + tagsize
}
func sizeFixed32ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toUint32()
	if v == 0 {
		return 0
	}
	return 4 + tagsize
}
func sizeFixed32Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toUint32Ptr()
	if p == nil {
		return 0
	}
	return 4 + tagsize
}
func sizeFixed32Slice(ptr pointer, tagsize int) int {
	s := *ptr.toUint32Slice()
	return (4 + tagsize) * len(s)
}
func sizeFixed32PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toUint32Slice()
	if len(s) == 0 {
		return 0
	}
	return 4*len(s) + SizeVarint(uint64(4*len(s))) + tagsize
}
func sizeFixedS32Value(_ pointer, tagsize int) int {
	return 4 + tagsize
}
func sizeFixedS32ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toInt32()
	if v == 0 {
		return 0
	}
	return 4 + tagsize
}
func sizeFixedS32Ptr(ptr pointer, tagsize int) int {
	p := ptr.getInt32Ptr()
	if p == nil {
		return 0
	}
	return 4 + tagsize
}
func sizeFixedS32Slice(ptr pointer, tagsize int) int {
	s := ptr.getInt32Slice()
	return (4 + tagsize) * len(s)
}
func sizeFixedS32PackedSlice(ptr pointer, tagsize int) int {
	s := ptr.getInt32Slice()
	if len(s) == 0 {
		return 0
	}
	return 4*len(s) + SizeVarint(uint64(4*len(s))) + tagsize
}
func sizeFloat32Value(_ pointer, tagsize int) int {
	return 4 + tagsize
}
func sizeFloat32ValueNoZero(ptr pointer, tagsize int) int {
	v := math.Float32bits(*ptr.toFloat32())
	if v == 0 {
		return 0
	}
	return 4 + tagsize
}
func sizeFloat32Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toFloat32Ptr()
	if p == nil {
		return 0
	}
	return 4 + tagsize
}
func sizeFloat32Slice(ptr pointer, tagsize int) int {
	s := *ptr.toFloat32Slice()
	return (4 + tagsize) * len(s)
}
func sizeFloat32PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toFloat32Slice()
	if len(s) == 0 {
		return 0
	}
	return 4*len(s) + SizeVarint(uint64(4*len(s))) + tagsize
}
func sizeFixed64Value(_ pointer, tagsize int) int {
	return 8 + tagsize
}
func sizeFixed64ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toUint64()
	if v == 0 {
		return 0
	}
	return 8 + tagsize
}
func sizeFixed64Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toUint64Ptr()
	if p == nil {
		return 0
	}
	return 8 + tagsize
}
func sizeFixed64Slice(ptr pointer, tagsize int) int {
	s := *ptr.toUint64Slice()
	return (8 + tagsize) * len(s)
}
func sizeFixed64PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toUint64Slice()
	if len(s) == 0 {
		return 0
	}
	return 8*len(s) + SizeVarint(uint64(8*len(s))) + tagsize
}
func sizeFixedS64Value(_ pointer, tagsize int) int {
	return 8 + tagsize
}
func sizeFixedS64ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toInt64()
	if v == 0 {
		return 0
	}
	return 8 + tagsize
}
func sizeFixedS64Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toInt64Ptr()
	if p == nil {
		return 0
	}
	return 8 + tagsize
}
func sizeFixedS64Slice(ptr pointer, tagsize int) int {
	s := *ptr.toInt64Slice()
	return (8 + tagsize) * len(s)
}
func sizeFixedS64PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toInt64Slice()
	if len(s) == 0 {
		return 0
	}
	return 8*len(s) + SizeVarint(uint64(8*len(s))) + tagsize
}
func sizeFloat64Value(_ pointer, tagsize int) int {
	return 8 + tagsize
}
func sizeFloat64ValueNoZero(ptr pointer, tagsize int) int {
	v := math.Float64bits(*ptr.toFloat64())
	if v == 0 {
		return 0
	}
	return 8 + tagsize
}
func sizeFloat64Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toFloat64Ptr()
	if p == nil {
		return 0
	}
	return 8 + tagsize
}
func sizeFloat64Slice(ptr pointer, tagsize int) int {
	s := *ptr.toFloat64Slice()
	return (8 + tagsize) * len(s)
}
func sizeFloat64PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toFloat64Slice()
	if len(s) == 0 {
		return 0
	}
	return 8*len(s) + SizeVarint(uint64(8*len(s))) + tagsize
}
func sizeVarint32Value(ptr pointer, tagsize int) int {
	v := *ptr.toUint32()
	return SizeVarint(uint64(v)) + tagsize
}
func sizeVarint32ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toUint32()
	if v == 0 {
		return 0
	}
	return SizeVarint(uint64(v)) + tagsize
}
func sizeVarint32Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toUint32Ptr()
	if p == nil {
		return 0
	}
	return SizeVarint(uint64(*p)) + tagsize
}
func sizeVarint32Slice(ptr pointer, tagsize int) int {
	s := *ptr.toUint32Slice()
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v)) + tagsize
	}
	return n
}
func sizeVarint32PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toUint32Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v))
	}
	return n + SizeVarint(uint64(n)) + tagsize
}
func sizeVarintS32Value(ptr pointer, tagsize int) int {
	v := *ptr.toInt32()
	return SizeVarint(uint64(v)) + tagsize
}
func sizeVarintS32ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toInt32()
	if v == 0 {
		return 0
	}
	return SizeVarint(uint64(v)) + tagsize
}
func sizeVarintS32Ptr(ptr pointer, tagsize int) int {
	p := ptr.getInt32Ptr()
	if p == nil {
		return 0
	}
	return SizeVarint(uint64(*p)) + tagsize
}
func sizeVarintS32Slice(ptr pointer, tagsize int) int {
	s := ptr.getInt32Slice()
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v)) + tagsize
	}
	return n
}
func sizeVarintS32PackedSlice(ptr pointer, tagsize int) int {
	s := ptr.getInt32Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v))
	}
	return n + SizeVarint(uint64(n)) + tagsize
}
func sizeVarint64Value(ptr pointer, tagsize int) int {
	v := *ptr.toUint64()
	return SizeVarint(v) + tagsize
}
func sizeVarint64ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toUint64()
	if v == 0 {
		return 0
	}
	return SizeVarint(v) + tagsize
}
func sizeVarint64Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toUint64Ptr()
	if p == nil {
		return 0
	}
	return SizeVarint(*p) + tagsize
}
func sizeVarint64Slice(ptr pointer, tagsize int) int {
	s := *ptr.toUint64Slice()
	n := 0
	for _, v := range s {
		n += SizeVarint(v) + tagsize
	}
	return n
}
func sizeVarint64PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toUint64Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += SizeVarint(v)
	}
	return n + SizeVarint(uint64(n)) + tagsize
}
func sizeVarintS64Value(ptr pointer, tagsize int) int {
	v := *ptr.toInt64()
	return SizeVarint(uint64(v)) + tagsize
}
func sizeVarintS64ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toInt64()
	if v == 0 {
		return 0
	}
	return SizeVarint(uint64(v)) + tagsize
}
func sizeVarintS64Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toInt64Ptr()
	if p == nil {
		return 0
	}
	return SizeVarint(uint64(*p)) + tagsize
}
func sizeVarintS64Slice(ptr pointer, tagsize int) int {
	s := *ptr.toInt64Slice()
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v)) + tagsize
	}
	return n
}
func sizeVarintS64PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toInt64Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v))
	}
	return n + SizeVarint(uint64(n)) + tagsize
}
func sizeZigzag32Value(ptr pointer, tagsize int) int {
	v := *ptr.toInt32()
	return SizeVarint(uint64((uint32(v)<<1)^uint32((int32(v)>>31)))) + tagsize
}
func sizeZigzag32ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toInt32()
	if v == 0 {
		return 0
	}
	return SizeVarint(uint64((uint32(v)<<1)^uint32((int32(v)>>31)))) + tagsize
}
func sizeZigzag32Ptr(ptr pointer, tagsize int) int {
	p := ptr.getInt32Ptr()
	if p == nil {
		return 0
	}
	v := *p
	return SizeVarint(uint64((uint32(v)<<1)^uint32((int32(v)>>31)))) + tagsize
}
func sizeZigzag32Slice(ptr pointer, tagsize int) int {
	s := ptr.getInt32Slice()
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64((uint32(v)<<1)^uint32((int32(v)>>31)))) + tagsize
	}
	return n
}
func sizeZigzag32PackedSlice(ptr pointer, tagsize int) int {
	s := ptr.getInt32Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64((uint32(v) << 1) ^ uint32((int32(v) >> 31))))
	}
	return n + SizeVarint(uint64(n)) + tagsize
}
func sizeZigzag64Value(ptr pointer, tagsize int) int {
	v := *ptr.toInt64()
	return SizeVarint(uint64(v<<1)^uint64((int64(v)>>63))) + tagsize
}
func sizeZigzag64ValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toInt64()
	if v == 0 {
		return 0
	}
	return SizeVarint(uint64(v<<1)^uint64((int64(v)>>63))) + tagsize
}
func sizeZigzag64Ptr(ptr pointer, tagsize int) int {
	p := *ptr.toInt64Ptr()
	if p == nil {
		return 0
	}
	v := *p
	return SizeVarint(uint64(v<<1)^uint64((int64(v)>>63))) + tagsize
}
func sizeZigzag64Slice(ptr pointer, tagsize int) int {
	s := *ptr.toInt64Slice()
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v<<1)^uint64((int64(v)>>63))) + tagsize
	}
	return n
}
func sizeZigzag64PackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toInt64Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v<<1) ^ uint64((int64(v) >> 63)))
	}
	return n + SizeVarint(uint64(n)) + tagsize
}
func sizeBoolValue(_ pointer, tagsize int) int {
	return 1 + tagsize
}
func sizeBoolValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toBool()
	if !v {
		return 0
	}
	return 1 + tagsize
}
func sizeBoolPtr(ptr pointer, tagsize int) int {
	p := *ptr.toBoolPtr()
	if p == nil {
		return 0
	}
	return 1 + tagsize
}
func sizeBoolSlice(ptr pointer, tagsize int) int {
	s := *ptr.toBoolSlice()
	return (1 + tagsize) * len(s)
}
func sizeBoolPackedSlice(ptr pointer, tagsize int) int {
	s := *ptr.toBoolSlice()
	if len(s) == 0 {
		return 0
	}
	return len(s) + SizeVarint(uint64(len(s))) + tagsize
}
func sizeStringValue(ptr pointer, tagsize int) int {
	v := *ptr.toString()
	return len(v) + SizeVarint(uint64(len(v))) + tagsize
}
func sizeStringValueNoZero(ptr pointer, tagsize int) int {
	v := *ptr.toString()
	if v == "" {
		return 0
	}
	return len(v) + SizeVarint(uint64(len(v))) + tagsize
}
func sizeStringPtr(ptr pointer, tagsize int) int {
	p := *ptr.toStringPtr()
	if p == nil {
		return 0
	}
	v := *p
	return len(v) + SizeVarint(uint64(len(v))) + tagsize
}
func sizeStringSlice(ptr pointer, tagsize int) int {
	s := *ptr.toStringSlice()
	n := 0
	for _, v := range s {
		n += len(v) + SizeVarint(uint64(len(v))) + tagsize
	}
	return n
}
func sizeBytes(ptr pointer, tagsize int) int {
	v := *ptr.toBytes()
	if v == nil {
		return 0
	}
	return len(v) + SizeVarint(uint64(len(v))) + tagsize
}
func sizeBytes3(ptr pointer, tagsize int) int {
	v := *ptr.toBytes()
	if len(v) == 0 {
		return 0
	}
	return len(v) + SizeVarint(uint64(len(v))) + tagsize
}
func sizeBytesOneof(ptr pointer, tagsize int) int {
	v := *ptr.toBytes()
	return len(v) + SizeVarint(uint64(len(v))) + tagsize
}
func sizeBytesSlice(ptr pointer, tagsize int) int {
	s := *ptr.toBytesSlice()
	n := 0
	for _, v := range s {
		n += len(v) + SizeVarint(uint64(len(v))) + tagsize
	}
	return n
}


func appendFixed32(b []byte, v uint32) []byte {
	b = append(b,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24))
	return b
}


func appendFixed64(b []byte, v uint64) []byte {
	b = append(b,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56))
	return b
}


func appendVarint(b []byte, v uint64) []byte {
	
	
	switch {
	case v < 1<<7:
		b = append(b, byte(v))
	case v < 1<<14:
		b = append(b,
			byte(v&0x7f|0x80),
			byte(v>>7))
	case v < 1<<21:
		b = append(b,
			byte(v&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte(v>>14))
	case v < 1<<28:
		b = append(b,
			byte(v&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte(v>>21))
	case v < 1<<35:
		b = append(b,
			byte(v&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte(v>>28))
	case v < 1<<42:
		b = append(b,
			byte(v&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte(v>>35))
	case v < 1<<49:
		b = append(b,
			byte(v&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte(v>>42))
	case v < 1<<56:
		b = append(b,
			byte(v&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte(v>>49))
	case v < 1<<63:
		b = append(b,
			byte(v&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte(v>>56))
	default:
		b = append(b,
			byte(v&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte((v>>56)&0x7f|0x80),
			1)
	}
	return b
}

func appendFixed32Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toUint32()
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, v)
	return b, nil
}
func appendFixed32ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toUint32()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, v)
	return b, nil
}
func appendFixed32Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toUint32Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, *p)
	return b, nil
}
func appendFixed32Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toUint32Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendFixed32(b, v)
	}
	return b, nil
}
func appendFixed32PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toUint32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	b = appendVarint(b, uint64(4*len(s)))
	for _, v := range s {
		b = appendFixed32(b, v)
	}
	return b, nil
}
func appendFixedS32Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt32()
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, uint32(v))
	return b, nil
}
func appendFixedS32ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt32()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, uint32(v))
	return b, nil
}
func appendFixedS32Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := ptr.getInt32Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, uint32(*p))
	return b, nil
}
func appendFixedS32Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := ptr.getInt32Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendFixed32(b, uint32(v))
	}
	return b, nil
}
func appendFixedS32PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := ptr.getInt32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	b = appendVarint(b, uint64(4*len(s)))
	for _, v := range s {
		b = appendFixed32(b, uint32(v))
	}
	return b, nil
}
func appendFloat32Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := math.Float32bits(*ptr.toFloat32())
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, v)
	return b, nil
}
func appendFloat32ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := math.Float32bits(*ptr.toFloat32())
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, v)
	return b, nil
}
func appendFloat32Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toFloat32Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed32(b, math.Float32bits(*p))
	return b, nil
}
func appendFloat32Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toFloat32Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendFixed32(b, math.Float32bits(v))
	}
	return b, nil
}
func appendFloat32PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toFloat32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	b = appendVarint(b, uint64(4*len(s)))
	for _, v := range s {
		b = appendFixed32(b, math.Float32bits(v))
	}
	return b, nil
}
func appendFixed64Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toUint64()
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, v)
	return b, nil
}
func appendFixed64ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toUint64()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, v)
	return b, nil
}
func appendFixed64Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toUint64Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, *p)
	return b, nil
}
func appendFixed64Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toUint64Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendFixed64(b, v)
	}
	return b, nil
}
func appendFixed64PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toUint64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	b = appendVarint(b, uint64(8*len(s)))
	for _, v := range s {
		b = appendFixed64(b, v)
	}
	return b, nil
}
func appendFixedS64Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt64()
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, uint64(v))
	return b, nil
}
func appendFixedS64ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt64()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, uint64(v))
	return b, nil
}
func appendFixedS64Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toInt64Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, uint64(*p))
	return b, nil
}
func appendFixedS64Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toInt64Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendFixed64(b, uint64(v))
	}
	return b, nil
}
func appendFixedS64PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toInt64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	b = appendVarint(b, uint64(8*len(s)))
	for _, v := range s {
		b = appendFixed64(b, uint64(v))
	}
	return b, nil
}
func appendFloat64Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := math.Float64bits(*ptr.toFloat64())
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, v)
	return b, nil
}
func appendFloat64ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := math.Float64bits(*ptr.toFloat64())
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, v)
	return b, nil
}
func appendFloat64Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toFloat64Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendFixed64(b, math.Float64bits(*p))
	return b, nil
}
func appendFloat64Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toFloat64Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendFixed64(b, math.Float64bits(v))
	}
	return b, nil
}
func appendFloat64PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toFloat64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	b = appendVarint(b, uint64(8*len(s)))
	for _, v := range s {
		b = appendFixed64(b, math.Float64bits(v))
	}
	return b, nil
}
func appendVarint32Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toUint32()
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(v))
	return b, nil
}
func appendVarint32ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toUint32()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(v))
	return b, nil
}
func appendVarint32Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toUint32Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(*p))
	return b, nil
}
func appendVarint32Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toUint32Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendVarint(b, uint64(v))
	}
	return b, nil
}
func appendVarint32PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toUint32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v))
	}
	b = appendVarint(b, uint64(n))
	for _, v := range s {
		b = appendVarint(b, uint64(v))
	}
	return b, nil
}
func appendVarintS32Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt32()
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(v))
	return b, nil
}
func appendVarintS32ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt32()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(v))
	return b, nil
}
func appendVarintS32Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := ptr.getInt32Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(*p))
	return b, nil
}
func appendVarintS32Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := ptr.getInt32Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendVarint(b, uint64(v))
	}
	return b, nil
}
func appendVarintS32PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := ptr.getInt32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v))
	}
	b = appendVarint(b, uint64(n))
	for _, v := range s {
		b = appendVarint(b, uint64(v))
	}
	return b, nil
}
func appendVarint64Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toUint64()
	b = appendVarint(b, wiretag)
	b = appendVarint(b, v)
	return b, nil
}
func appendVarint64ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toUint64()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, v)
	return b, nil
}
func appendVarint64Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toUint64Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, *p)
	return b, nil
}
func appendVarint64Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toUint64Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendVarint(b, v)
	}
	return b, nil
}
func appendVarint64PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toUint64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	
	n := 0
	for _, v := range s {
		n += SizeVarint(v)
	}
	b = appendVarint(b, uint64(n))
	for _, v := range s {
		b = appendVarint(b, v)
	}
	return b, nil
}
func appendVarintS64Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt64()
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(v))
	return b, nil
}
func appendVarintS64ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt64()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(v))
	return b, nil
}
func appendVarintS64Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toInt64Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(*p))
	return b, nil
}
func appendVarintS64Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toInt64Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendVarint(b, uint64(v))
	}
	return b, nil
}
func appendVarintS64PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toInt64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v))
	}
	b = appendVarint(b, uint64(n))
	for _, v := range s {
		b = appendVarint(b, uint64(v))
	}
	return b, nil
}
func appendZigzag32Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt32()
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64((uint32(v)<<1)^uint32((int32(v)>>31))))
	return b, nil
}
func appendZigzag32ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt32()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64((uint32(v)<<1)^uint32((int32(v)>>31))))
	return b, nil
}
func appendZigzag32Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := ptr.getInt32Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	v := *p
	b = appendVarint(b, uint64((uint32(v)<<1)^uint32((int32(v)>>31))))
	return b, nil
}
func appendZigzag32Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := ptr.getInt32Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendVarint(b, uint64((uint32(v)<<1)^uint32((int32(v)>>31))))
	}
	return b, nil
}
func appendZigzag32PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := ptr.getInt32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64((uint32(v) << 1) ^ uint32((int32(v) >> 31))))
	}
	b = appendVarint(b, uint64(n))
	for _, v := range s {
		b = appendVarint(b, uint64((uint32(v)<<1)^uint32((int32(v)>>31))))
	}
	return b, nil
}
func appendZigzag64Value(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt64()
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(v<<1)^uint64((int64(v)>>63)))
	return b, nil
}
func appendZigzag64ValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toInt64()
	if v == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(v<<1)^uint64((int64(v)>>63)))
	return b, nil
}
func appendZigzag64Ptr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toInt64Ptr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	v := *p
	b = appendVarint(b, uint64(v<<1)^uint64((int64(v)>>63)))
	return b, nil
}
func appendZigzag64Slice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toInt64Slice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendVarint(b, uint64(v<<1)^uint64((int64(v)>>63)))
	}
	return b, nil
}
func appendZigzag64PackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toInt64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	
	n := 0
	for _, v := range s {
		n += SizeVarint(uint64(v<<1) ^ uint64((int64(v) >> 63)))
	}
	b = appendVarint(b, uint64(n))
	for _, v := range s {
		b = appendVarint(b, uint64(v<<1)^uint64((int64(v)>>63)))
	}
	return b, nil
}
func appendBoolValue(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toBool()
	b = appendVarint(b, wiretag)
	if v {
		b = append(b, 1)
	} else {
		b = append(b, 0)
	}
	return b, nil
}
func appendBoolValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toBool()
	if !v {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = append(b, 1)
	return b, nil
}

func appendBoolPtr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toBoolPtr()
	if p == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	if *p {
		b = append(b, 1)
	} else {
		b = append(b, 0)
	}
	return b, nil
}
func appendBoolSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toBoolSlice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		if v {
			b = append(b, 1)
		} else {
			b = append(b, 0)
		}
	}
	return b, nil
}
func appendBoolPackedSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toBoolSlice()
	if len(s) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag&^7|WireBytes)
	b = appendVarint(b, uint64(len(s)))
	for _, v := range s {
		if v {
			b = append(b, 1)
		} else {
			b = append(b, 0)
		}
	}
	return b, nil
}
func appendStringValue(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toString()
	if !utf8.ValidString(v) {
		return nil, errInvalidUTF8
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(len(v)))
	b = append(b, v...)
	return b, nil
}
func appendStringValueNoZero(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toString()
	if v == "" {
		return b, nil
	}
	if !utf8.ValidString(v) {
		return nil, errInvalidUTF8
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(len(v)))
	b = append(b, v...)
	return b, nil
}
func appendStringPtr(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	p := *ptr.toStringPtr()
	if p == nil {
		return b, nil
	}
	v := *p
	if !utf8.ValidString(v) {
		return nil, errInvalidUTF8
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(len(v)))
	b = append(b, v...)
	return b, nil
}
func appendStringSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toStringSlice()
	for _, v := range s {
		if !utf8.ValidString(v) {
			return nil, errInvalidUTF8
		}
		b = appendVarint(b, wiretag)
		b = appendVarint(b, uint64(len(v)))
		b = append(b, v...)
	}
	return b, nil
}
func appendBytes(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toBytes()
	if v == nil {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(len(v)))
	b = append(b, v...)
	return b, nil
}
func appendBytes3(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toBytes()
	if len(v) == 0 {
		return b, nil
	}
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(len(v)))
	b = append(b, v...)
	return b, nil
}
func appendBytesOneof(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	v := *ptr.toBytes()
	b = appendVarint(b, wiretag)
	b = appendVarint(b, uint64(len(v)))
	b = append(b, v...)
	return b, nil
}
func appendBytesSlice(b []byte, ptr pointer, wiretag uint64, _ bool) ([]byte, error) {
	s := *ptr.toBytesSlice()
	for _, v := range s {
		b = appendVarint(b, wiretag)
		b = appendVarint(b, uint64(len(v)))
		b = append(b, v...)
	}
	return b, nil
}



func makeGroupMarshaler(u *marshalInfo) (sizer, marshaler) {
	return func(ptr pointer, tagsize int) int {
			p := ptr.getPointer()
			if p.isNil() {
				return 0
			}
			return u.size(p) + 2*tagsize
		},
		func(b []byte, ptr pointer, wiretag uint64, deterministic bool) ([]byte, error) {
			p := ptr.getPointer()
			if p.isNil() {
				return b, nil
			}
			var err error
			b = appendVarint(b, wiretag) 
			b, err = u.marshal(b, p, deterministic)
			b = appendVarint(b, wiretag+(WireEndGroup-WireStartGroup)) 
			return b, err
		}
}



func makeGroupSliceMarshaler(u *marshalInfo) (sizer, marshaler) {
	return func(ptr pointer, tagsize int) int {
			s := ptr.getPointerSlice()
			n := 0
			for _, v := range s {
				if v.isNil() {
					continue
				}
				n += u.size(v) + 2*tagsize
			}
			return n
		},
		func(b []byte, ptr pointer, wiretag uint64, deterministic bool) ([]byte, error) {
			s := ptr.getPointerSlice()
			var err, errreq error
			for _, v := range s {
				if v.isNil() {
					return b, errRepeatedHasNil
				}
				b = appendVarint(b, wiretag) 
				b, err = u.marshal(b, v, deterministic)
				b = appendVarint(b, wiretag+(WireEndGroup-WireStartGroup)) 
				if err != nil {
					if _, ok := err.(*RequiredNotSetError); ok {
						
						
						if errreq == nil {
							errreq = err
						}
						continue
					}
					if err == ErrNil {
						err = errRepeatedHasNil
					}
					return b, err
				}
			}
			return b, errreq
		}
}



func makeMessageMarshaler(u *marshalInfo) (sizer, marshaler) {
	return func(ptr pointer, tagsize int) int {
			p := ptr.getPointer()
			if p.isNil() {
				return 0
			}
			siz := u.size(p)
			return siz + SizeVarint(uint64(siz)) + tagsize
		},
		func(b []byte, ptr pointer, wiretag uint64, deterministic bool) ([]byte, error) {
			p := ptr.getPointer()
			if p.isNil() {
				return b, nil
			}
			b = appendVarint(b, wiretag)
			siz := u.cachedsize(p)
			b = appendVarint(b, uint64(siz))
			return u.marshal(b, p, deterministic)
		}
}



func makeMessageSliceMarshaler(u *marshalInfo) (sizer, marshaler) {
	return func(ptr pointer, tagsize int) int {
			s := ptr.getPointerSlice()
			n := 0
			for _, v := range s {
				if v.isNil() {
					continue
				}
				siz := u.size(v)
				n += siz + SizeVarint(uint64(siz)) + tagsize
			}
			return n
		},
		func(b []byte, ptr pointer, wiretag uint64, deterministic bool) ([]byte, error) {
			s := ptr.getPointerSlice()
			var err, errreq error
			for _, v := range s {
				if v.isNil() {
					return b, errRepeatedHasNil
				}
				b = appendVarint(b, wiretag)
				siz := u.cachedsize(v)
				b = appendVarint(b, uint64(siz))
				b, err = u.marshal(b, v, deterministic)

				if err != nil {
					if _, ok := err.(*RequiredNotSetError); ok {
						
						
						if errreq == nil {
							errreq = err
						}
						continue
					}
					if err == ErrNil {
						err = errRepeatedHasNil
					}
					return b, err
				}
			}
			return b, errreq
		}
}



func makeMapMarshaler(f *reflect.StructField) (sizer, marshaler) {
	
	t := f.Type
	keyType := t.Key()
	valType := t.Elem()
	keyTags := strings.Split(f.Tag.Get("protobuf_key"), ",")
	valTags := strings.Split(f.Tag.Get("protobuf_val"), ",")
	keySizer, keyMarshaler := typeMarshaler(keyType, keyTags, false, false) 
	valSizer, valMarshaler := typeMarshaler(valType, valTags, false, false) 
	keyWireTag := 1<<3 | wiretype(keyTags[0])
	valWireTag := 2<<3 | wiretype(valTags[0])

	
	
	
	
	
	valIsPtr := valType.Kind() == reflect.Ptr
	return func(ptr pointer, tagsize int) int {
			m := ptr.asPointerTo(t).Elem() 
			n := 0
			for _, k := range m.MapKeys() {
				ki := k.Interface()
				vi := m.MapIndex(k).Interface()
				kaddr := toAddrPointer(&ki, false)             
				vaddr := toAddrPointer(&vi, valIsPtr)          
				siz := keySizer(kaddr, 1) + valSizer(vaddr, 1) 
				n += siz + SizeVarint(uint64(siz)) + tagsize
			}
			return n
		},
		func(b []byte, ptr pointer, tag uint64, deterministic bool) ([]byte, error) {
			m := ptr.asPointerTo(t).Elem() 
			var err error
			keys := m.MapKeys()
			if len(keys) > 1 && deterministic {
				sort.Sort(mapKeys(keys))
			}
			for _, k := range keys {
				ki := k.Interface()
				vi := m.MapIndex(k).Interface()
				kaddr := toAddrPointer(&ki, false)    
				vaddr := toAddrPointer(&vi, valIsPtr) 
				b = appendVarint(b, tag)
				siz := keySizer(kaddr, 1) + valSizer(vaddr, 1) 
				b = appendVarint(b, uint64(siz))
				b, err = keyMarshaler(b, kaddr, keyWireTag, deterministic)
				if err != nil {
					return b, err
				}
				b, err = valMarshaler(b, vaddr, valWireTag, deterministic)
				if err != nil && err != ErrNil { 
					return b, err
				}
			}
			return b, nil
		}
}




func makeOneOfMarshaler(fi *marshalFieldInfo, f *reflect.StructField) (sizer, marshaler) {
	
	t := f.Type
	return func(ptr pointer, _ int) int {
			p := ptr.getInterfacePointer()
			if p.isNil() {
				return 0
			}
			v := ptr.asPointerTo(t).Elem().Elem().Elem() 
			telem := v.Type()
			e := fi.oneofElems[telem]
			return e.sizer(p, e.tagsize)
		},
		func(b []byte, ptr pointer, _ uint64, deterministic bool) ([]byte, error) {
			p := ptr.getInterfacePointer()
			if p.isNil() {
				return b, nil
			}
			v := ptr.asPointerTo(t).Elem().Elem().Elem() 
			telem := v.Type()
			if telem.Field(0).Type.Kind() == reflect.Ptr && p.getPointer().isNil() {
				return b, errOneofHasNil
			}
			e := fi.oneofElems[telem]
			return e.marshaler(b, p, e.wiretag, deterministic)
		}
}


func (u *marshalInfo) sizeExtensions(ext *XXX_InternalExtensions) int {
	m, mu := ext.extensionsRead()
	if m == nil {
		return 0
	}
	mu.Lock()

	n := 0
	for _, e := range m {
		if e.value == nil || e.desc == nil {
			
			n += len(e.enc)
			continue
		}

		
		
		
		ei := u.getExtElemInfo(e.desc)
		v := e.value
		p := toAddrPointer(&v, ei.isptr)
		n += ei.sizer(p, ei.tagsize)
	}
	mu.Unlock()
	return n
}


func (u *marshalInfo) appendExtensions(b []byte, ext *XXX_InternalExtensions, deterministic bool) ([]byte, error) {
	m, mu := ext.extensionsRead()
	if m == nil {
		return b, nil
	}
	mu.Lock()
	defer mu.Unlock()

	var err error

	
	
	if len(m) <= 1 {
		for _, e := range m {
			if e.value == nil || e.desc == nil {
				
				b = append(b, e.enc...)
				continue
			}

			
			
			

			ei := u.getExtElemInfo(e.desc)
			v := e.value
			p := toAddrPointer(&v, ei.isptr)
			b, err = ei.marshaler(b, p, ei.wiretag, deterministic)
			if err != nil {
				return b, err
			}
		}
		return b, nil
	}

	
	
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, k := range keys {
		e := m[int32(k)]
		if e.value == nil || e.desc == nil {
			
			b = append(b, e.enc...)
			continue
		}

		
		
		

		ei := u.getExtElemInfo(e.desc)
		v := e.value
		p := toAddrPointer(&v, ei.isptr)
		b, err = ei.marshaler(b, p, ei.wiretag, deterministic)
		if err != nil {
			return b, err
		}
	}
	return b, nil
}











func (u *marshalInfo) sizeMessageSet(ext *XXX_InternalExtensions) int {
	m, mu := ext.extensionsRead()
	if m == nil {
		return 0
	}
	mu.Lock()

	n := 0
	for id, e := range m {
		n += 2                          
		n += SizeVarint(uint64(id)) + 1 

		if e.value == nil || e.desc == nil {
			
			msgWithLen := skipVarint(e.enc) 
			siz := len(msgWithLen)
			n += siz + 1 
			continue
		}

		
		
		

		ei := u.getExtElemInfo(e.desc)
		v := e.value
		p := toAddrPointer(&v, ei.isptr)
		n += ei.sizer(p, 1) 
	}
	mu.Unlock()
	return n
}



func (u *marshalInfo) appendMessageSet(b []byte, ext *XXX_InternalExtensions, deterministic bool) ([]byte, error) {
	m, mu := ext.extensionsRead()
	if m == nil {
		return b, nil
	}
	mu.Lock()
	defer mu.Unlock()

	var err error

	
	
	if len(m) <= 1 {
		for id, e := range m {
			b = append(b, 1<<3|WireStartGroup)
			b = append(b, 2<<3|WireVarint)
			b = appendVarint(b, uint64(id))

			if e.value == nil || e.desc == nil {
				
				msgWithLen := skipVarint(e.enc) 
				b = append(b, 3<<3|WireBytes)
				b = append(b, msgWithLen...)
				b = append(b, 1<<3|WireEndGroup)
				continue
			}

			
			
			

			ei := u.getExtElemInfo(e.desc)
			v := e.value
			p := toAddrPointer(&v, ei.isptr)
			b, err = ei.marshaler(b, p, 3<<3|WireBytes, deterministic)
			if err != nil {
				return b, err
			}
			b = append(b, 1<<3|WireEndGroup)
		}
		return b, nil
	}

	
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, id := range keys {
		e := m[int32(id)]
		b = append(b, 1<<3|WireStartGroup)
		b = append(b, 2<<3|WireVarint)
		b = appendVarint(b, uint64(id))

		if e.value == nil || e.desc == nil {
			
			msgWithLen := skipVarint(e.enc) 
			b = append(b, 3<<3|WireBytes)
			b = append(b, msgWithLen...)
			b = append(b, 1<<3|WireEndGroup)
			continue
		}

		
		
		

		ei := u.getExtElemInfo(e.desc)
		v := e.value
		p := toAddrPointer(&v, ei.isptr)
		b, err = ei.marshaler(b, p, 3<<3|WireBytes, deterministic)
		b = append(b, 1<<3|WireEndGroup)
		if err != nil {
			return b, err
		}
	}
	return b, nil
}


func (u *marshalInfo) sizeV1Extensions(m map[int32]Extension) int {
	if m == nil {
		return 0
	}

	n := 0
	for _, e := range m {
		if e.value == nil || e.desc == nil {
			
			n += len(e.enc)
			continue
		}

		
		
		

		ei := u.getExtElemInfo(e.desc)
		v := e.value
		p := toAddrPointer(&v, ei.isptr)
		n += ei.sizer(p, ei.tagsize)
	}
	return n
}


func (u *marshalInfo) appendV1Extensions(b []byte, m map[int32]Extension, deterministic bool) ([]byte, error) {
	if m == nil {
		return b, nil
	}

	
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	var err error
	for _, k := range keys {
		e := m[int32(k)]
		if e.value == nil || e.desc == nil {
			
			b = append(b, e.enc...)
			continue
		}

		
		
		

		ei := u.getExtElemInfo(e.desc)
		v := e.value
		p := toAddrPointer(&v, ei.isptr)
		b, err = ei.marshaler(b, p, ei.wiretag, deterministic)
		if err != nil {
			return b, err
		}
	}
	return b, nil
}







type newMarshaler interface {
	XXX_Size() int
	XXX_Marshal(b []byte, deterministic bool) ([]byte, error)
}



func Size(pb Message) int {
	if m, ok := pb.(newMarshaler); ok {
		return m.XXX_Size()
	}
	if m, ok := pb.(Marshaler); ok {
		
		
		b, _ := m.Marshal()
		return len(b)
	}
	
	if pb == nil {
		return 0
	}
	var info InternalMessageInfo
	return info.Size(pb)
}




func Marshal(pb Message) ([]byte, error) {
	if m, ok := pb.(newMarshaler); ok {
		siz := m.XXX_Size()
		b := make([]byte, 0, siz)
		return m.XXX_Marshal(b, false)
	}
	if m, ok := pb.(Marshaler); ok {
		
		
		return m.Marshal()
	}
	
	if pb == nil {
		return nil, ErrNil
	}
	var info InternalMessageInfo
	siz := info.Size(pb)
	b := make([]byte, 0, siz)
	return info.Marshal(b, pb, false)
}






func (p *Buffer) Marshal(pb Message) error {
	var err error
	if m, ok := pb.(newMarshaler); ok {
		siz := m.XXX_Size()
		p.grow(siz) 
		p.buf, err = m.XXX_Marshal(p.buf, p.deterministic)
		return err
	}
	if m, ok := pb.(Marshaler); ok {
		
		
		b, err := m.Marshal()
		p.buf = append(p.buf, b...)
		return err
	}
	
	if pb == nil {
		return ErrNil
	}
	var info InternalMessageInfo
	siz := info.Size(pb)
	p.grow(siz) 
	p.buf, err = info.Marshal(p.buf, pb, p.deterministic)
	return err
}




func (p *Buffer) grow(n int) {
	need := len(p.buf) + n
	if need <= cap(p.buf) {
		return
	}
	newCap := len(p.buf) * 2
	if newCap < need {
		newCap = need
	}
	p.buf = append(make([]byte, 0, newCap), p.buf...)
}
