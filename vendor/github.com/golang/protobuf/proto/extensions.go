






























package proto



import (
	"errors"
	"fmt"
	"io"
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




func extendable(p interface{}) (extendableProto, error) {
	switch p := p.(type) {
	case extendableProto:
		if isNilPtr(p) {
			return nil, fmt.Errorf("proto: nil %T is not extendable", p)
		}
		return p, nil
	case extendableProtoV1:
		if isNilPtr(p) {
			return nil, fmt.Errorf("proto: nil %T is not extendable", p)
		}
		return extensionAdapter{p}, nil
	}
	
	
	return nil, errNotExtendable
}

var errNotExtendable = errors.New("proto: not an extendable proto.Message")

func isNilPtr(x interface{}) bool {
	v := reflect.ValueOf(x)
	return v.Kind() == reflect.Ptr && v.IsNil()
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
	epb, err := extendable(base)
	if err != nil {
		return
	}
	extmap := epb.extensionsWrite()
	extmap[id] = Extension{enc: b}
}


func isExtensionField(pb extendableProto, field int32) bool {
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
		return fmt.Errorf("proto: bad extended type; %v does not extend %v", b, a)
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


func HasExtension(pb Message, extension *ExtensionDesc) bool {
	
	epb, err := extendable(pb)
	if err != nil {
		return false
	}
	extmap, mu := epb.extensionsRead()
	if extmap == nil {
		return false
	}
	mu.Lock()
	_, ok := extmap[extension.Field]
	mu.Unlock()
	return ok
}


func ClearExtension(pb Message, extension *ExtensionDesc) {
	epb, err := extendable(pb)
	if err != nil {
		return
	}
	
	extmap := epb.extensionsWrite()
	delete(extmap, extension.Field)
}










func GetExtension(pb Message, extension *ExtensionDesc) (interface{}, error) {
	epb, err := extendable(pb)
	if err != nil {
		return nil, err
	}

	if extension.ExtendedType != nil {
		
		if err := checkExtensionTypes(epb, extension); err != nil {
			return nil, err
		}
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

	if extension.ExtensionType == nil {
		
		return e.enc, nil
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
	if extension.ExtensionType == nil {
		
		return nil, ErrMissingExtension
	}

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
	t := reflect.TypeOf(extension.ExtensionType)
	unmarshal := typeUnmarshaler(t, extension.Tag)

	
	
	value := reflect.New(t).Elem()

	var err error
	for {
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		wire := int(x) & 7

		b, err = unmarshal(b, valToPointer(value.Addr()), wire)
		if err != nil {
			return nil, err
		}

		if len(b) == 0 {
			break
		}
	}
	return value.Interface(), nil
}



func GetExtensions(pb Message, es []*ExtensionDesc) (extensions []interface{}, err error) {
	epb, err := extendable(pb)
	if err != nil {
		return nil, err
	}
	extensions = make([]interface{}, len(es))
	for i, e := range es {
		extensions[i], err = GetExtension(epb, e)
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
	epb, err := extendable(pb)
	if err != nil {
		return nil, err
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
	epb, err := extendable(pb)
	if err != nil {
		return err
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
	epb, err := extendable(pb)
	if err != nil {
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
