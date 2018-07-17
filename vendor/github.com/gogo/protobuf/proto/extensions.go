






























package proto



import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)


var ErrMissingExtension = errors.New("proto: missing extension")



type ExtensionRange struct {
	Start, End int32 
}



type extendableProto interface {
	Message
	ExtensionRangeArray() []ExtensionRange
	extensionsWrite() map[int32]Extension
	extensionsRead() (map[int32]Extension, sync.Locker)
}



type extendableProtoV1 interface {
	Message
	ExtensionRangeArray() []ExtensionRange
	ExtensionMap() map[int32]Extension
}

type extensionsBytes interface {
	Message
	ExtensionRangeArray() []ExtensionRange
	GetExtensions() *[]byte
}


type extensionAdapter struct {
	extendableProtoV1
}

func (e extensionAdapter) extensionsWrite() map[int32]Extension {
	return e.ExtensionMap()
}

func (e extensionAdapter) extensionsRead() (map[int32]Extension, sync.Locker) {
	return e.ExtensionMap(), notLocker{}
}


type notLocker struct{}

func (n notLocker) Lock()   {}
func (n notLocker) Unlock() {}




func extendable(p interface{}) (extendableProto, bool) {
	if ep, ok := p.(extendableProto); ok {
		return ep, ok
	}
	if ep, ok := p.(extendableProtoV1); ok {
		return extensionAdapter{ep}, ok
	}
	return nil, false
}








type XXX_InternalExtensions struct {
	
	
	
	
	
	
	
	p *struct {
		mu           sync.Mutex
		extensionMap map[int32]Extension
	}
}


func (e *XXX_InternalExtensions) extensionsWrite() map[int32]Extension {
	if e.p == nil {
		e.p = new(struct {
			mu           sync.Mutex
			extensionMap map[int32]Extension
		})
		e.p.extensionMap = make(map[int32]Extension)
	}
	return e.p.extensionMap
}



func (e *XXX_InternalExtensions) extensionsRead() (map[int32]Extension, sync.Locker) {
	if e.p == nil {
		return nil, nil
	}
	return e.p.extensionMap, &e.p.mu
}

type extensionRange interface {
	Message
	ExtensionRangeArray() []ExtensionRange
}

var extendableProtoType = reflect.TypeOf((*extendableProto)(nil)).Elem()
var extendableProtoV1Type = reflect.TypeOf((*extendableProtoV1)(nil)).Elem()
var extendableBytesType = reflect.TypeOf((*extensionsBytes)(nil)).Elem()
var extensionRangeType = reflect.TypeOf((*extensionRange)(nil)).Elem()



type ExtensionDesc struct {
	ExtendedType  Message     
	ExtensionType interface{} 
	Field         int32       
	Name          string      
	Tag           string      
	Filename      string      
}

func (ed *ExtensionDesc) repeated() bool {
	t := reflect.TypeOf(ed.ExtensionType)
	return t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8
}


type Extension struct {
	
	
	
	
	
	
	
	
	desc  *ExtensionDesc
	value interface{}
	enc   []byte
}


func SetRawExtension(base Message, id int32, b []byte) {
	if ebase, ok := base.(extensionsBytes); ok {
		clearExtension(base, id)
		ext := ebase.GetExtensions()
		*ext = append(*ext, b...)
		return
	}
	epb, ok := extendable(base)
	if !ok {
		return
	}
	extmap := epb.extensionsWrite()
	extmap[id] = Extension{enc: b}
}


func isExtensionField(pb extensionRange, field int32) bool {
	for _, er := range pb.ExtensionRangeArray() {
		if er.Start <= field && field <= er.End {
			return true
		}
	}
	return false
}


func checkExtensionTypes(pb extendableProto, extension *ExtensionDesc) error {
	var pbi interface{} = pb
	
	if ea, ok := pbi.(extensionAdapter); ok {
		pbi = ea.extendableProtoV1
	}
	if a, b := reflect.TypeOf(pbi), reflect.TypeOf(extension.ExtendedType); a != b {
		return errors.New("proto: bad extended type; " + b.String() + " does not extend " + a.String())
	}
	
	if !isExtensionField(pb, extension.Field) {
		return errors.New("proto: bad extension number; not in declared ranges")
	}
	return nil
}


type extPropKey struct {
	base  reflect.Type
	field int32
}

var extProp = struct {
	sync.RWMutex
	m map[extPropKey]*Properties
}{
	m: make(map[extPropKey]*Properties),
}

func extensionProperties(ed *ExtensionDesc) *Properties {
	key := extPropKey{base: reflect.TypeOf(ed.ExtendedType), field: ed.Field}

	extProp.RLock()
	if prop, ok := extProp.m[key]; ok {
		extProp.RUnlock()
		return prop
	}
	extProp.RUnlock()

	extProp.Lock()
	defer extProp.Unlock()
	
	if prop, ok := extProp.m[key]; ok {
		return prop
	}

	prop := new(Properties)
	prop.Init(reflect.TypeOf(ed.ExtensionType), "unknown_name", ed.Tag, nil)
	extProp.m[key] = prop
	return prop
}


func encodeExtensions(e *XXX_InternalExtensions) error {
	m, mu := e.extensionsRead()
	if m == nil {
		return nil 
	}
	mu.Lock()
	defer mu.Unlock()
	return encodeExtensionsMap(m)
}


func encodeExtensionsMap(m map[int32]Extension) error {
	for k, e := range m {
		if e.value == nil || e.desc == nil {
			
			continue
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
		m[k] = e
	}
	return nil
}

func extensionsSize(e *XXX_InternalExtensions) (n int) {
	m, mu := e.extensionsRead()
	if m == nil {
		return 0
	}
	mu.Lock()
	defer mu.Unlock()
	return extensionsMapSize(m)
}

func extensionsMapSize(m map[int32]Extension) (n int) {
	for _, e := range m {
		if e.value == nil || e.desc == nil {
			
			n += len(e.enc)
			continue
		}

		
		
		

		et := reflect.TypeOf(e.desc.ExtensionType)
		props := extensionProperties(e.desc)

		
		
		x := reflect.New(et)
		x.Elem().Set(reflect.ValueOf(e.value))
		n += props.size(props, toStructPointer(x))
	}
	return
}


func HasExtension(pb Message, extension *ExtensionDesc) bool {
	if epb, doki := pb.(extensionsBytes); doki {
		ext := epb.GetExtensions()
		buf := *ext
		o := 0
		for o < len(buf) {
			tag, n := DecodeVarint(buf[o:])
			fieldNum := int32(tag >> 3)
			if int32(fieldNum) == extension.Field {
				return true
			}
			wireType := int(tag & 0x7)
			o += n
			l, err := size(buf[o:], wireType)
			if err != nil {
				return false
			}
			o += l
		}
		return false
	}
	
	epb, ok := extendable(pb)
	if !ok {
		return false
	}
	extmap, mu := epb.extensionsRead()
	if extmap == nil {
		return false
	}
	mu.Lock()
	_, ok = extmap[extension.Field]
	mu.Unlock()
	return ok
}

func deleteExtension(pb extensionsBytes, theFieldNum int32, offset int) int {
	ext := pb.GetExtensions()
	for offset < len(*ext) {
		tag, n1 := DecodeVarint((*ext)[offset:])
		fieldNum := int32(tag >> 3)
		wireType := int(tag & 0x7)
		n2, err := size((*ext)[offset+n1:], wireType)
		if err != nil {
			panic(err)
		}
		newOffset := offset + n1 + n2
		if fieldNum == theFieldNum {
			*ext = append((*ext)[:offset], (*ext)[newOffset:]...)
			return offset
		}
		offset = newOffset
	}
	return -1
}


func ClearExtension(pb Message, extension *ExtensionDesc) {
	clearExtension(pb, extension.Field)
}

func clearExtension(pb Message, fieldNum int32) {
	if epb, doki := pb.(extensionsBytes); doki {
		offset := 0
		for offset != -1 {
			offset = deleteExtension(epb, fieldNum, offset)
		}
		return
	}
	epb, ok := extendable(pb)
	if !ok {
		return
	}
	
	extmap := epb.extensionsWrite()
	delete(extmap, fieldNum)
}



func GetExtension(pb Message, extension *ExtensionDesc) (interface{}, error) {
	if epb, doki := pb.(extensionsBytes); doki {
		ext := epb.GetExtensions()
		o := 0
		for o < len(*ext) {
			tag, n := DecodeVarint((*ext)[o:])
			fieldNum := int32(tag >> 3)
			wireType := int(tag & 0x7)
			l, err := size((*ext)[o+n:], wireType)
			if err != nil {
				return nil, err
			}
			if int32(fieldNum) == extension.Field {
				v, err := decodeExtension((*ext)[o:o+n+l], extension)
				if err != nil {
					return nil, err
				}
				return v, nil
			}
			o += n + l
		}
		return defaultExtensionValue(extension)
	}
	epb, ok := extendable(pb)
	if !ok {
		return nil, errors.New("proto: not an extendable proto")
	}
	if err := checkExtensionTypes(epb, extension); err != nil {
		return nil, err
	}

	emap, mu := epb.extensionsRead()
	if emap == nil {
		return defaultExtensionValue(extension)
	}
	mu.Lock()
	defer mu.Unlock()
	e, ok := emap[extension.Field]
	if !ok {
		
		
		return defaultExtensionValue(extension)
	}

	if e.value != nil {
		
		if e.desc != extension {
			
			
			
			return nil, errors.New("proto: descriptor conflict")
		}
		return e.value, nil
	}

	v, err := decodeExtension(e.enc, extension)
	if err != nil {
		return nil, err
	}

	
	
	e.value = v
	e.desc = extension
	e.enc = nil
	emap[extension.Field] = e
	return e.value, nil
}



func defaultExtensionValue(extension *ExtensionDesc) (interface{}, error) {
	t := reflect.TypeOf(extension.ExtensionType)
	props := extensionProperties(extension)

	sf, _, err := fieldDefault(t, props)
	if err != nil {
		return nil, err
	}

	if sf == nil || sf.value == nil {
		
		return nil, ErrMissingExtension
	}

	if t.Kind() != reflect.Ptr {
		
		return sf.value, nil
	}

	
	value := reflect.New(t).Elem()
	value.Set(reflect.New(value.Type().Elem()))
	if sf.kind == reflect.Int32 {
		
		
		
		value.Elem().SetInt(int64(sf.value.(int32)))
	} else {
		value.Elem().Set(reflect.ValueOf(sf.value))
	}
	return value.Interface(), nil
}


func decodeExtension(b []byte, extension *ExtensionDesc) (interface{}, error) {
	o := NewBuffer(b)

	t := reflect.TypeOf(extension.ExtensionType)

	props := extensionProperties(extension)

	
	
	
	
	
	
	value := reflect.New(t).Elem()

	for {
		
		if _, err := o.DecodeVarint(); err != nil {
			return nil, err
		}

		if err := props.dec(o, props, toStructPointer(value.Addr())); err != nil {
			return nil, err
		}

		if o.index >= len(o.buf) {
			break
		}
	}
	return value.Interface(), nil
}



func GetExtensions(pb Message, es []*ExtensionDesc) (extensions []interface{}, err error) {
	extensions = make([]interface{}, len(es))
	for i, e := range es {
		extensions[i], err = GetExtension(pb, e)
		if err == ErrMissingExtension {
			err = nil
		}
		if err != nil {
			return
		}
	}
	return
}




func ExtensionDescs(pb Message) ([]*ExtensionDesc, error) {
	epb, ok := extendable(pb)
	if !ok {
		return nil, fmt.Errorf("proto: %T is not an extendable proto.Message", pb)
	}
	registeredExtensions := RegisteredExtensions(pb)

	emap, mu := epb.extensionsRead()
	if emap == nil {
		return nil, nil
	}
	mu.Lock()
	defer mu.Unlock()
	extensions := make([]*ExtensionDesc, 0, len(emap))
	for extid, e := range emap {
		desc := e.desc
		if desc == nil {
			desc = registeredExtensions[extid]
			if desc == nil {
				desc = &ExtensionDesc{Field: extid}
			}
		}

		extensions = append(extensions, desc)
	}
	return extensions, nil
}


func SetExtension(pb Message, extension *ExtensionDesc, value interface{}) error {
	if epb, doki := pb.(extensionsBytes); doki {
		ClearExtension(pb, extension)
		ext := epb.GetExtensions()
		et := reflect.TypeOf(extension.ExtensionType)
		props := extensionProperties(extension)
		p := NewBuffer(nil)
		x := reflect.New(et)
		x.Elem().Set(reflect.ValueOf(value))
		if err := props.enc(p, props, toStructPointer(x)); err != nil {
			return err
		}
		*ext = append(*ext, p.buf...)
		return nil
	}
	epb, ok := extendable(pb)
	if !ok {
		return errors.New("proto: not an extendable proto")
	}
	if err := checkExtensionTypes(epb, extension); err != nil {
		return err
	}
	typ := reflect.TypeOf(extension.ExtensionType)
	if typ != reflect.TypeOf(value) {
		return errors.New("proto: bad extension value type")
	}
	
	
	
	
	
	if reflect.ValueOf(value).IsNil() {
		return fmt.Errorf("proto: SetExtension called with nil value of type %T", value)
	}

	extmap := epb.extensionsWrite()
	extmap[extension.Field] = Extension{desc: extension, value: value}
	return nil
}


func ClearAllExtensions(pb Message) {
	if epb, doki := pb.(extensionsBytes); doki {
		ext := epb.GetExtensions()
		*ext = []byte{}
		return
	}
	epb, ok := extendable(pb)
	if !ok {
		return
	}
	m := epb.extensionsWrite()
	for k := range m {
		delete(m, k)
	}
}




var extensionMaps = make(map[reflect.Type]map[int32]*ExtensionDesc)


func RegisterExtension(desc *ExtensionDesc) {
	st := reflect.TypeOf(desc.ExtendedType).Elem()
	m := extensionMaps[st]
	if m == nil {
		m = make(map[int32]*ExtensionDesc)
		extensionMaps[st] = m
	}
	if _, ok := m[desc.Field]; ok {
		panic("proto: duplicate extension registered: " + st.String() + " " + strconv.Itoa(int(desc.Field)))
	}
	m[desc.Field] = desc
}




func RegisteredExtensions(pb Message) map[int32]*ExtensionDesc {
	return extensionMaps[reflect.TypeOf(pb).Elem()]
}
