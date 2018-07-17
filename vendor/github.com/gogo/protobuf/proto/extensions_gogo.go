



























package proto

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

func GetBoolExtension(pb Message, extension *ExtensionDesc, ifnotset bool) bool {
	if reflect.ValueOf(pb).IsNil() {
		return ifnotset
	}
	value, err := GetExtension(pb, extension)
	if err != nil {
		return ifnotset
	}
	if value == nil {
		return ifnotset
	}
	if value.(*bool) == nil {
		return ifnotset
	}
	return *(value.(*bool))
}

func (this *Extension) Equal(that *Extension) bool {
	return bytes.Equal(this.enc, that.enc)
}

func (this *Extension) Compare(that *Extension) int {
	return bytes.Compare(this.enc, that.enc)
}

func SizeOfInternalExtension(m extendableProto) (n int) {
	return SizeOfExtensionMap(m.extensionsWrite())
}

func SizeOfExtensionMap(m map[int32]Extension) (n int) {
	return extensionsMapSize(m)
}

type sortableMapElem struct {
	field int32
	ext   Extension
}

func newSortableExtensionsFromMap(m map[int32]Extension) sortableExtensions {
	s := make(sortableExtensions, 0, len(m))
	for k, v := range m {
		s = append(s, &sortableMapElem{field: k, ext: v})
	}
	return s
}

type sortableExtensions []*sortableMapElem

func (this sortableExtensions) Len() int { return len(this) }

func (this sortableExtensions) Swap(i, j int) { this[i], this[j] = this[j], this[i] }

func (this sortableExtensions) Less(i, j int) bool { return this[i].field < this[j].field }

func (this sortableExtensions) String() string {
	sort.Sort(this)
	ss := make([]string, len(this))
	for i := range this {
		ss[i] = fmt.Sprintf("%d: %v", this[i].field, this[i].ext)
	}
	return "map[" + strings.Join(ss, ",") + "]"
}

func StringFromInternalExtension(m extendableProto) string {
	return StringFromExtensionsMap(m.extensionsWrite())
}

func StringFromExtensionsMap(m map[int32]Extension) string {
	return newSortableExtensionsFromMap(m).String()
}

func StringFromExtensionsBytes(ext []byte) string {
	m, err := BytesToExtensionsMap(ext)
	if err != nil {
		panic(err)
	}
	return StringFromExtensionsMap(m)
}

func EncodeInternalExtension(m extendableProto, data []byte) (n int, err error) {
	return EncodeExtensionMap(m.extensionsWrite(), data)
}

func EncodeExtensionMap(m map[int32]Extension, data []byte) (n int, err error) {
	if err := encodeExtensionsMap(m); err != nil {
		return 0, err
	}
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
		n += copy(data[n:], m[int32(k)].enc)
	}
	return n, nil
}

func GetRawExtension(m map[int32]Extension, id int32) ([]byte, error) {
	if m[id].value == nil || m[id].desc == nil {
		return m[id].enc, nil
	}
	if err := encodeExtensionsMap(m); err != nil {
		return nil, err
	}
	return m[id].enc, nil
}

func size(buf []byte, wire int) (int, error) {
	switch wire {
	case WireVarint:
		_, n := DecodeVarint(buf)
		return n, nil
	case WireFixed64:
		return 8, nil
	case WireBytes:
		v, n := DecodeVarint(buf)
		return int(v) + n, nil
	case WireFixed32:
		return 4, nil
	case WireStartGroup:
		offset := 0
		for {
			u, n := DecodeVarint(buf[offset:])
			fwire := int(u & 0x7)
			offset += n
			if fwire == WireEndGroup {
				return offset, nil
			}
			s, err := size(buf[offset:], wire)
			if err != nil {
				return 0, err
			}
			offset += s
		}
	}
	return 0, fmt.Errorf("proto: can't get size for unknown wire type %d", wire)
}

func BytesToExtensionsMap(buf []byte) (map[int32]Extension, error) {
	m := make(map[int32]Extension)
	i := 0
	for i < len(buf) {
		tag, n := DecodeVarint(buf[i:])
		if n <= 0 {
			return nil, fmt.Errorf("unable to decode varint")
		}
		fieldNum := int32(tag >> 3)
		wireType := int(tag & 0x7)
		l, err := size(buf[i+n:], wireType)
		if err != nil {
			return nil, err
		}
		end := i + int(l) + n
		m[int32(fieldNum)] = Extension{enc: buf[i:end]}
		i = end
	}
	return m, nil
}

func NewExtension(e []byte) Extension {
	ee := Extension{enc: make([]byte, len(e))}
	copy(ee.enc, e)
	return ee
}

func AppendExtension(e Message, tag int32, buf []byte) {
	if ee, eok := e.(extensionsBytes); eok {
		ext := ee.GetExtensions()
		*ext = append(*ext, buf...)
		return
	}
	if ee, eok := e.(extendableProto); eok {
		m := ee.extensionsWrite()
		ext := m[int32(tag)] 
		ext.enc = append(ext.enc, buf...)
		m[int32(tag)] = ext
	}
}

func encodeExtension(e *Extension) error {
	if e.value == nil || e.desc == nil {
		
		return nil
	}
	
	
	

	et := reflect.TypeOf(e.desc.ExtensionType)
	props := extensionProperties(e.desc)

	p := NewBuffer(nil)
	
	
	x := reflect.New(et)
	x.Elem().Set(reflect.ValueOf(e.value))
	if err := props.enc(p, props, toStructPointer(x)); err != nil {
		return err
	}
	e.enc = p.buf
	return nil
}

func (this Extension) GoString() string {
	if this.enc == nil {
		if err := encodeExtension(&this); err != nil {
			panic(err)
		}
	}
	return fmt.Sprintf("proto.NewExtension(%#v)", this.enc)
}

func SetUnsafeExtension(pb Message, fieldNum int32, value interface{}) error {
	typ := reflect.TypeOf(pb).Elem()
	ext, ok := extensionMaps[typ]
	if !ok {
		return fmt.Errorf("proto: bad extended type; %s is not extendable", typ.String())
	}
	desc, ok := ext[fieldNum]
	if !ok {
		return errors.New("proto: bad extension number; not in declared ranges")
	}
	return SetExtension(pb, desc, value)
}

func GetUnsafeExtension(pb Message, fieldNum int32) (interface{}, error) {
	typ := reflect.TypeOf(pb).Elem()
	ext, ok := extensionMaps[typ]
	if !ok {
		return nil, fmt.Errorf("proto: bad extended type; %s is not extendable", typ.String())
	}
	desc, ok := ext[fieldNum]
	if !ok {
		return nil, fmt.Errorf("unregistered field number %d", fieldNum)
	}
	return GetExtension(pb, desc)
}

func NewUnsafeXXX_InternalExtensions(m map[int32]Extension) XXX_InternalExtensions {
	x := &XXX_InternalExtensions{
		p: new(struct {
			mu           sync.Mutex
			extensionMap map[int32]Extension
		}),
	}
	x.p.extensionMap = m
	return *x
}

func GetUnsafeExtensionsMap(extendable Message) map[int32]Extension {
	pb := extendable.(extendableProto)
	return pb.extensionsWrite()
}
