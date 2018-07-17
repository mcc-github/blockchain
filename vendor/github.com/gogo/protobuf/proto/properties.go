



































package proto



import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const debug bool = false


const (
	WireVarint     = 0
	WireFixed64    = 1
	WireBytes      = 2
	WireStartGroup = 3
	WireEndGroup   = 4
	WireFixed32    = 5
)

const startSize = 10 




type encoder func(p *Buffer, prop *Properties, base structPointer) error


type valueEncoder func(o *Buffer, x uint64) error




type sizer func(prop *Properties, base structPointer) int



type valueSizer func(x uint64) int




type decoder func(p *Buffer, prop *Properties, base structPointer) error


type valueDecoder func(o *Buffer) (x uint64, err error)


type oneofMarshaler func(Message, *Buffer) error


type oneofUnmarshaler func(Message, int, int, *Buffer) (bool, error)


type oneofSizer func(Message) int




type tagMap struct {
	fastTags []int
	slowTags map[int]int
}



const tagMapFastLimit = 1024

func (p *tagMap) get(t int) (int, bool) {
	if t > 0 && t < tagMapFastLimit {
		if t >= len(p.fastTags) {
			return 0, false
		}
		fi := p.fastTags[t]
		return fi, fi >= 0
	}
	fi, ok := p.slowTags[t]
	return fi, ok
}

func (p *tagMap) put(t int, fi int) {
	if t > 0 && t < tagMapFastLimit {
		for len(p.fastTags) < t+1 {
			p.fastTags = append(p.fastTags, -1)
		}
		p.fastTags[t] = fi
		return
	}
	if p.slowTags == nil {
		p.slowTags = make(map[int]int)
	}
	p.slowTags[t] = fi
}



type StructProperties struct {
	Prop             []*Properties  
	reqCount         int            
	decoderTags      tagMap         
	decoderOrigNames map[string]int 
	order            []int          
	unrecField       field          
	extendable       bool           

	oneofMarshaler   oneofMarshaler
	oneofUnmarshaler oneofUnmarshaler
	oneofSizer       oneofSizer
	stype            reflect.Type

	
	
	OneofTypes map[string]*OneofProperties
}


type OneofProperties struct {
	Type  reflect.Type 
	Field int          
	Prop  *Properties
}




func (sp *StructProperties) Len() int { return len(sp.order) }
func (sp *StructProperties) Less(i, j int) bool {
	return sp.Prop[sp.order[i]].Tag < sp.Prop[sp.order[j]].Tag
}
func (sp *StructProperties) Swap(i, j int) { sp.order[i], sp.order[j] = sp.order[j], sp.order[i] }


type Properties struct {
	Name     string 
	OrigName string 
	JSONName string 
	Wire     string
	WireType int
	Tag      int
	Required bool
	Optional bool
	Repeated bool
	Packed   bool   
	Enum     string 
	proto3   bool   
	oneof    bool   

	Default     string 
	HasDefault  bool   
	CustomType  string
	CastType    string
	StdTime     bool
	StdDuration bool

	enc           encoder
	valEnc        valueEncoder 
	field         field
	tagcode       []byte 
	tagbuf        [8]byte
	stype         reflect.Type      
	sstype        reflect.Type      
	ctype         reflect.Type      
	sprop         *StructProperties 
	isMarshaler   bool
	isUnmarshaler bool

	mtype    reflect.Type 
	mkeyprop *Properties  
	mvalprop *Properties  

	size    sizer
	valSize valueSizer 

	dec    decoder
	valDec valueDecoder 

	
	packedDec decoder
}


func (p *Properties) String() string {
	s := p.Wire
	s = ","
	s += strconv.Itoa(p.Tag)
	if p.Required {
		s += ",req"
	}
	if p.Optional {
		s += ",opt"
	}
	if p.Repeated {
		s += ",rep"
	}
	if p.Packed {
		s += ",packed"
	}
	s += ",name=" + p.OrigName
	if p.JSONName != p.OrigName {
		s += ",json=" + p.JSONName
	}
	if p.proto3 {
		s += ",proto3"
	}
	if p.oneof {
		s += ",oneof"
	}
	if len(p.Enum) > 0 {
		s += ",enum=" + p.Enum
	}
	if p.HasDefault {
		s += ",def=" + p.Default
	}
	return s
}


func (p *Properties) Parse(s string) {
	
	fields := strings.Split(s, ",") 
	if len(fields) < 2 {
		fmt.Fprintf(os.Stderr, "proto: tag has too few fields: %q\n", s)
		return
	}

	p.Wire = fields[0]
	switch p.Wire {
	case "varint":
		p.WireType = WireVarint
		p.valEnc = (*Buffer).EncodeVarint
		p.valDec = (*Buffer).DecodeVarint
		p.valSize = sizeVarint
	case "fixed32":
		p.WireType = WireFixed32
		p.valEnc = (*Buffer).EncodeFixed32
		p.valDec = (*Buffer).DecodeFixed32
		p.valSize = sizeFixed32
	case "fixed64":
		p.WireType = WireFixed64
		p.valEnc = (*Buffer).EncodeFixed64
		p.valDec = (*Buffer).DecodeFixed64
		p.valSize = sizeFixed64
	case "zigzag32":
		p.WireType = WireVarint
		p.valEnc = (*Buffer).EncodeZigzag32
		p.valDec = (*Buffer).DecodeZigzag32
		p.valSize = sizeZigzag32
	case "zigzag64":
		p.WireType = WireVarint
		p.valEnc = (*Buffer).EncodeZigzag64
		p.valDec = (*Buffer).DecodeZigzag64
		p.valSize = sizeZigzag64
	case "bytes", "group":
		p.WireType = WireBytes
		
	default:
		fmt.Fprintf(os.Stderr, "proto: tag has unknown wire type: %q\n", s)
		return
	}

	var err error
	p.Tag, err = strconv.Atoi(fields[1])
	if err != nil {
		return
	}

	for i := 2; i < len(fields); i++ {
		f := fields[i]
		switch {
		case f == "req":
			p.Required = true
		case f == "opt":
			p.Optional = true
		case f == "rep":
			p.Repeated = true
		case f == "packed":
			p.Packed = true
		case strings.HasPrefix(f, "name="):
			p.OrigName = f[5:]
		case strings.HasPrefix(f, "json="):
			p.JSONName = f[5:]
		case strings.HasPrefix(f, "enum="):
			p.Enum = f[5:]
		case f == "proto3":
			p.proto3 = true
		case f == "oneof":
			p.oneof = true
		case strings.HasPrefix(f, "def="):
			p.HasDefault = true
			p.Default = f[4:] 
			if i+1 < len(fields) {
				
				p.Default += "," + strings.Join(fields[i+1:], ",")
				break
			}
		case strings.HasPrefix(f, "embedded="):
			p.OrigName = strings.Split(f, "=")[1]
		case strings.HasPrefix(f, "customtype="):
			p.CustomType = strings.Split(f, "=")[1]
		case strings.HasPrefix(f, "casttype="):
			p.CastType = strings.Split(f, "=")[1]
		case f == "stdtime":
			p.StdTime = true
		case f == "stdduration":
			p.StdDuration = true
		}
	}
}

func logNoSliceEnc(t1, t2 reflect.Type) {
	fmt.Fprintf(os.Stderr, "proto: no slice oenc for %T = []%T\n", t1, t2)
}

var protoMessageType = reflect.TypeOf((*Message)(nil)).Elem()


func (p *Properties) setEncAndDec(typ reflect.Type, f *reflect.StructField, lockGetProp bool) {
	p.enc = nil
	p.dec = nil
	p.size = nil
	isMap := typ.Kind() == reflect.Map
	if len(p.CustomType) > 0 && !isMap {
		p.setCustomEncAndDec(typ)
		p.setTag(lockGetProp)
		return
	}
	if p.StdTime && !isMap {
		p.setTimeEncAndDec(typ)
		p.setTag(lockGetProp)
		return
	}
	if p.StdDuration && !isMap {
		p.setDurationEncAndDec(typ)
		p.setTag(lockGetProp)
		return
	}
	switch t1 := typ; t1.Kind() {
	default:
		fmt.Fprintf(os.Stderr, "proto: no coders for %v\n", t1)

	

	case reflect.Bool:
		if p.proto3 {
			p.enc = (*Buffer).enc_proto3_bool
			p.dec = (*Buffer).dec_proto3_bool
			p.size = size_proto3_bool
		} else {
			p.enc = (*Buffer).enc_ref_bool
			p.dec = (*Buffer).dec_proto3_bool
			p.size = size_ref_bool
		}
	case reflect.Int32:
		if p.proto3 {
			p.enc = (*Buffer).enc_proto3_int32
			p.dec = (*Buffer).dec_proto3_int32
			p.size = size_proto3_int32
		} else {
			p.enc = (*Buffer).enc_ref_int32
			p.dec = (*Buffer).dec_proto3_int32
			p.size = size_ref_int32
		}
	case reflect.Uint32:
		if p.proto3 {
			p.enc = (*Buffer).enc_proto3_uint32
			p.dec = (*Buffer).dec_proto3_int32 
			p.size = size_proto3_uint32
		} else {
			p.enc = (*Buffer).enc_ref_uint32
			p.dec = (*Buffer).dec_proto3_int32 
			p.size = size_ref_uint32
		}
	case reflect.Int64, reflect.Uint64:
		if p.proto3 {
			p.enc = (*Buffer).enc_proto3_int64
			p.dec = (*Buffer).dec_proto3_int64
			p.size = size_proto3_int64
		} else {
			p.enc = (*Buffer).enc_ref_int64
			p.dec = (*Buffer).dec_proto3_int64
			p.size = size_ref_int64
		}
	case reflect.Float32:
		if p.proto3 {
			p.enc = (*Buffer).enc_proto3_uint32 
			p.dec = (*Buffer).dec_proto3_int32
			p.size = size_proto3_uint32
		} else {
			p.enc = (*Buffer).enc_ref_uint32 
			p.dec = (*Buffer).dec_proto3_int32
			p.size = size_ref_uint32
		}
	case reflect.Float64:
		if p.proto3 {
			p.enc = (*Buffer).enc_proto3_int64 
			p.dec = (*Buffer).dec_proto3_int64
			p.size = size_proto3_int64
		} else {
			p.enc = (*Buffer).enc_ref_int64 
			p.dec = (*Buffer).dec_proto3_int64
			p.size = size_ref_int64
		}
	case reflect.String:
		if p.proto3 {
			p.enc = (*Buffer).enc_proto3_string
			p.dec = (*Buffer).dec_proto3_string
			p.size = size_proto3_string
		} else {
			p.enc = (*Buffer).enc_ref_string
			p.dec = (*Buffer).dec_proto3_string
			p.size = size_ref_string
		}
	case reflect.Struct:
		p.stype = typ
		p.isMarshaler = isMarshaler(typ)
		p.isUnmarshaler = isUnmarshaler(typ)
		if p.Wire == "bytes" {
			p.enc = (*Buffer).enc_ref_struct_message
			p.dec = (*Buffer).dec_ref_struct_message
			p.size = size_ref_struct_message
		} else {
			fmt.Fprintf(os.Stderr, "proto: no coders for struct %T\n", typ)
		}

	case reflect.Ptr:
		switch t2 := t1.Elem(); t2.Kind() {
		default:
			fmt.Fprintf(os.Stderr, "proto: no encoder function for %v -> %v\n", t1, t2)
			break
		case reflect.Bool:
			p.enc = (*Buffer).enc_bool
			p.dec = (*Buffer).dec_bool
			p.size = size_bool
		case reflect.Int32:
			p.enc = (*Buffer).enc_int32
			p.dec = (*Buffer).dec_int32
			p.size = size_int32
		case reflect.Uint32:
			p.enc = (*Buffer).enc_uint32
			p.dec = (*Buffer).dec_int32 
			p.size = size_uint32
		case reflect.Int64, reflect.Uint64:
			p.enc = (*Buffer).enc_int64
			p.dec = (*Buffer).dec_int64
			p.size = size_int64
		case reflect.Float32:
			p.enc = (*Buffer).enc_uint32 
			p.dec = (*Buffer).dec_int32
			p.size = size_uint32
		case reflect.Float64:
			p.enc = (*Buffer).enc_int64 
			p.dec = (*Buffer).dec_int64
			p.size = size_int64
		case reflect.String:
			p.enc = (*Buffer).enc_string
			p.dec = (*Buffer).dec_string
			p.size = size_string
		case reflect.Struct:
			p.stype = t1.Elem()
			p.isMarshaler = isMarshaler(t1)
			p.isUnmarshaler = isUnmarshaler(t1)
			if p.Wire == "bytes" {
				p.enc = (*Buffer).enc_struct_message
				p.dec = (*Buffer).dec_struct_message
				p.size = size_struct_message
			} else {
				p.enc = (*Buffer).enc_struct_group
				p.dec = (*Buffer).dec_struct_group
				p.size = size_struct_group
			}
		}

	case reflect.Slice:
		switch t2 := t1.Elem(); t2.Kind() {
		default:
			logNoSliceEnc(t1, t2)
			break
		case reflect.Bool:
			if p.Packed {
				p.enc = (*Buffer).enc_slice_packed_bool
				p.size = size_slice_packed_bool
			} else {
				p.enc = (*Buffer).enc_slice_bool
				p.size = size_slice_bool
			}
			p.dec = (*Buffer).dec_slice_bool
			p.packedDec = (*Buffer).dec_slice_packed_bool
		case reflect.Int32:
			if p.Packed {
				p.enc = (*Buffer).enc_slice_packed_int32
				p.size = size_slice_packed_int32
			} else {
				p.enc = (*Buffer).enc_slice_int32
				p.size = size_slice_int32
			}
			p.dec = (*Buffer).dec_slice_int32
			p.packedDec = (*Buffer).dec_slice_packed_int32
		case reflect.Uint32:
			if p.Packed {
				p.enc = (*Buffer).enc_slice_packed_uint32
				p.size = size_slice_packed_uint32
			} else {
				p.enc = (*Buffer).enc_slice_uint32
				p.size = size_slice_uint32
			}
			p.dec = (*Buffer).dec_slice_int32
			p.packedDec = (*Buffer).dec_slice_packed_int32
		case reflect.Int64, reflect.Uint64:
			if p.Packed {
				p.enc = (*Buffer).enc_slice_packed_int64
				p.size = size_slice_packed_int64
			} else {
				p.enc = (*Buffer).enc_slice_int64
				p.size = size_slice_int64
			}
			p.dec = (*Buffer).dec_slice_int64
			p.packedDec = (*Buffer).dec_slice_packed_int64
		case reflect.Uint8:
			p.dec = (*Buffer).dec_slice_byte
			if p.proto3 {
				p.enc = (*Buffer).enc_proto3_slice_byte
				p.size = size_proto3_slice_byte
			} else {
				p.enc = (*Buffer).enc_slice_byte
				p.size = size_slice_byte
			}
		case reflect.Float32, reflect.Float64:
			switch t2.Bits() {
			case 32:
				
				if p.Packed {
					p.enc = (*Buffer).enc_slice_packed_uint32
					p.size = size_slice_packed_uint32
				} else {
					p.enc = (*Buffer).enc_slice_uint32
					p.size = size_slice_uint32
				}
				p.dec = (*Buffer).dec_slice_int32
				p.packedDec = (*Buffer).dec_slice_packed_int32
			case 64:
				
				if p.Packed {
					p.enc = (*Buffer).enc_slice_packed_int64
					p.size = size_slice_packed_int64
				} else {
					p.enc = (*Buffer).enc_slice_int64
					p.size = size_slice_int64
				}
				p.dec = (*Buffer).dec_slice_int64
				p.packedDec = (*Buffer).dec_slice_packed_int64
			default:
				logNoSliceEnc(t1, t2)
				break
			}
		case reflect.String:
			p.enc = (*Buffer).enc_slice_string
			p.dec = (*Buffer).dec_slice_string
			p.size = size_slice_string
		case reflect.Ptr:
			switch t3 := t2.Elem(); t3.Kind() {
			default:
				fmt.Fprintf(os.Stderr, "proto: no ptr oenc for %T -> %T -> %T\n", t1, t2, t3)
				break
			case reflect.Struct:
				p.stype = t2.Elem()
				p.isMarshaler = isMarshaler(t2)
				p.isUnmarshaler = isUnmarshaler(t2)
				if p.Wire == "bytes" {
					p.enc = (*Buffer).enc_slice_struct_message
					p.dec = (*Buffer).dec_slice_struct_message
					p.size = size_slice_struct_message
				} else {
					p.enc = (*Buffer).enc_slice_struct_group
					p.dec = (*Buffer).dec_slice_struct_group
					p.size = size_slice_struct_group
				}
			}
		case reflect.Slice:
			switch t2.Elem().Kind() {
			default:
				fmt.Fprintf(os.Stderr, "proto: no slice elem oenc for %T -> %T -> %T\n", t1, t2, t2.Elem())
				break
			case reflect.Uint8:
				p.enc = (*Buffer).enc_slice_slice_byte
				p.dec = (*Buffer).dec_slice_slice_byte
				p.size = size_slice_slice_byte
			}
		case reflect.Struct:
			p.setSliceOfNonPointerStructs(t1)
		}

	case reflect.Map:
		p.enc = (*Buffer).enc_new_map
		p.dec = (*Buffer).dec_new_map
		p.size = size_new_map

		p.mtype = t1
		p.mkeyprop = &Properties{}
		p.mkeyprop.init(reflect.PtrTo(p.mtype.Key()), "Key", f.Tag.Get("protobuf_key"), nil, lockGetProp)
		p.mvalprop = &Properties{}
		vtype := p.mtype.Elem()
		if vtype.Kind() != reflect.Ptr && vtype.Kind() != reflect.Slice {
			
			
			vtype = reflect.PtrTo(vtype)
		}

		p.mvalprop.CustomType = p.CustomType
		p.mvalprop.StdDuration = p.StdDuration
		p.mvalprop.StdTime = p.StdTime
		p.mvalprop.init(vtype, "Value", f.Tag.Get("protobuf_val"), nil, lockGetProp)
	}
	p.setTag(lockGetProp)
}

func (p *Properties) setTag(lockGetProp bool) {
	
	wire := p.WireType
	if p.Packed {
		wire = WireBytes
	}
	x := uint32(p.Tag)<<3 | uint32(wire)
	i := 0
	for i = 0; x > 127; i++ {
		p.tagbuf[i] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	p.tagbuf[i] = uint8(x)
	p.tagcode = p.tagbuf[0 : i+1]

	if p.stype != nil {
		if lockGetProp {
			p.sprop = GetProperties(p.stype)
		} else {
			p.sprop = getPropertiesLocked(p.stype)
		}
	}
}

var (
	marshalerType   = reflect.TypeOf((*Marshaler)(nil)).Elem()
	unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)


func isMarshaler(t reflect.Type) bool {
	return t.Implements(marshalerType)
}


func isUnmarshaler(t reflect.Type) bool {
	return t.Implements(unmarshalerType)
}


func (p *Properties) Init(typ reflect.Type, name, tag string, f *reflect.StructField) {
	p.init(typ, name, tag, f, true)
}

func (p *Properties) init(typ reflect.Type, name, tag string, f *reflect.StructField, lockGetProp bool) {
	
	p.Name = name
	p.OrigName = name
	if f != nil {
		p.field = toField(f)
	}
	if tag == "" {
		return
	}
	p.Parse(tag)
	p.setEncAndDec(typ, f, lockGetProp)
}

var (
	propertiesMu  sync.RWMutex
	propertiesMap = make(map[reflect.Type]*StructProperties)
)



func GetProperties(t reflect.Type) *StructProperties {
	if t.Kind() != reflect.Struct {
		panic("proto: type must have kind struct")
	}

	
	
	propertiesMu.RLock()
	sprop, ok := propertiesMap[t]
	propertiesMu.RUnlock()
	if ok {
		if collectStats {
			stats.Chit++
		}
		return sprop
	}

	propertiesMu.Lock()
	sprop = getPropertiesLocked(t)
	propertiesMu.Unlock()
	return sprop
}


func getPropertiesLocked(t reflect.Type) *StructProperties {
	if prop, ok := propertiesMap[t]; ok {
		if collectStats {
			stats.Chit++
		}
		return prop
	}
	if collectStats {
		stats.Cmiss++
	}

	prop := new(StructProperties)
	
	propertiesMap[t] = prop

	
	prop.extendable = reflect.PtrTo(t).Implements(extendableProtoType) ||
		reflect.PtrTo(t).Implements(extendableProtoV1Type) ||
		reflect.PtrTo(t).Implements(extendableBytesType)
	prop.unrecField = invalidField
	prop.Prop = make([]*Properties, t.NumField())
	prop.order = make([]int, t.NumField())

	isOneofMessage := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		p := new(Properties)
		name := f.Name
		p.init(f.Type, name, f.Tag.Get("protobuf"), &f, false)

		if f.Name == "XXX_InternalExtensions" { 
			p.enc = (*Buffer).enc_exts
			p.dec = nil 
			p.size = size_exts
		} else if f.Name == "XXX_extensions" { 
			if len(f.Tag.Get("protobuf")) > 0 {
				p.enc = (*Buffer).enc_ext_slice_byte
				p.dec = nil 
				p.size = size_ext_slice_byte
			} else {
				p.enc = (*Buffer).enc_map
				p.dec = nil 
				p.size = size_map
			}
		} else if f.Name == "XXX_unrecognized" { 
			prop.unrecField = toField(&f)
		}
		oneof := f.Tag.Get("protobuf_oneof") 
		if oneof != "" {
			isOneofMessage = true
			
			p.OrigName = oneof
		}
		prop.Prop[i] = p
		prop.order[i] = i
		if debug {
			print(i, " ", f.Name, " ", t.String(), " ")
			if p.Tag > 0 {
				print(p.String())
			}
			print("\n")
		}
		if p.enc == nil && !strings.HasPrefix(f.Name, "XXX_") && oneof == "" {
			fmt.Fprintln(os.Stderr, "proto: no encoder for", f.Name, f.Type.String(), "[GetProperties]")
		}
	}

	
	sort.Sort(prop)

	type oneofMessage interface {
		XXX_OneofFuncs() (func(Message, *Buffer) error, func(Message, int, int, *Buffer) (bool, error), func(Message) int, []interface{})
	}
	if om, ok := reflect.Zero(reflect.PtrTo(t)).Interface().(oneofMessage); isOneofMessage && ok {
		var oots []interface{}
		prop.oneofMarshaler, prop.oneofUnmarshaler, prop.oneofSizer, oots = om.XXX_OneofFuncs()
		prop.stype = t

		
		prop.OneofTypes = make(map[string]*OneofProperties)
		for _, oot := range oots {
			oop := &OneofProperties{
				Type: reflect.ValueOf(oot).Type(), 
				Prop: new(Properties),
			}
			sft := oop.Type.Elem().Field(0)
			oop.Prop.Name = sft.Name
			oop.Prop.Parse(sft.Tag.Get("protobuf"))
			
			
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.Type.Kind() != reflect.Interface {
					continue
				}
				if !oop.Type.AssignableTo(f.Type) {
					continue
				}
				oop.Field = i
				break
			}
			prop.OneofTypes[oop.Prop.OrigName] = oop
		}
	}

	
	
	reqCount := 0
	prop.decoderOrigNames = make(map[string]int)
	for i, p := range prop.Prop {
		if strings.HasPrefix(p.Name, "XXX_") {
			
			
			continue
		}
		if p.Required {
			reqCount++
		}
		prop.decoderTags.put(p.Tag, i)
		prop.decoderOrigNames[p.OrigName] = i
	}
	prop.reqCount = reqCount

	return prop
}


func propByIndex(t reflect.Type, x []int) *Properties {
	if len(x) != 1 {
		fmt.Fprintf(os.Stderr, "proto: field index dimension %d (not 1) for type %s\n", len(x), t)
		return nil
	}
	prop := GetProperties(t)
	return prop.Prop[x[0]]
}


func getbase(pb Message) (t reflect.Type, b structPointer, err error) {
	if pb == nil {
		err = ErrNil
		return
	}
	
	t = reflect.TypeOf(pb)
	
	value := reflect.ValueOf(pb)
	b = toStructPointer(value)
	return
}




var enumValueMaps = make(map[string]map[string]int32)
var enumStringMaps = make(map[string]map[int32]string)



func RegisterEnum(typeName string, unusedNameMap map[int32]string, valueMap map[string]int32) {
	if _, ok := enumValueMaps[typeName]; ok {
		panic("proto: duplicate enum registered: " + typeName)
	}
	enumValueMaps[typeName] = valueMap
	if _, ok := enumStringMaps[typeName]; ok {
		panic("proto: duplicate enum registered: " + typeName)
	}
	enumStringMaps[typeName] = unusedNameMap
}



func EnumValueMap(enumType string) map[string]int32 {
	return enumValueMaps[enumType]
}



var (
	protoTypes    = make(map[string]reflect.Type)
	revProtoTypes = make(map[reflect.Type]string)
)



func RegisterType(x Message, name string) {
	if _, ok := protoTypes[name]; ok {
		
		log.Printf("proto: duplicate proto type registered: %s", name)
		return
	}
	t := reflect.TypeOf(x)
	protoTypes[name] = t
	revProtoTypes[t] = name
}


func MessageName(x Message) string {
	type xname interface {
		XXX_MessageName() string
	}
	if m, ok := x.(xname); ok {
		return m.XXX_MessageName()
	}
	return revProtoTypes[reflect.TypeOf(x)]
}


func MessageType(name string) reflect.Type { return protoTypes[name] }


var (
	protoFiles = make(map[string][]byte) 
)



func RegisterFile(filename string, fileDescriptor []byte) {
	protoFiles[filename] = fileDescriptor
}


func FileDescriptor(filename string) []byte { return protoFiles[filename] }
