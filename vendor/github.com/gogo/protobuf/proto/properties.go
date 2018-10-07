



































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

	stype reflect.Type      
	ctype reflect.Type      
	sprop *StructProperties 

	mtype    reflect.Type 
	mkeyprop *Properties  
	mvalprop *Properties  
}


func (p *Properties) String() string {
	s := p.Wire
	s += ","
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
	case "fixed32":
		p.WireType = WireFixed32
	case "fixed64":
		p.WireType = WireFixed64
	case "zigzag32":
		p.WireType = WireVarint
	case "zigzag64":
		p.WireType = WireVarint
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

outer:
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
				break outer
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

var protoMessageType = reflect.TypeOf((*Message)(nil)).Elem()


func (p *Properties) setFieldProps(typ reflect.Type, f *reflect.StructField, lockGetProp bool) {
	isMap := typ.Kind() == reflect.Map
	if len(p.CustomType) > 0 && !isMap {
		p.ctype = typ
		p.setTag(lockGetProp)
		return
	}
	if p.StdTime && !isMap {
		p.setTag(lockGetProp)
		return
	}
	if p.StdDuration && !isMap {
		p.setTag(lockGetProp)
		return
	}
	switch t1 := typ; t1.Kind() {
	case reflect.Struct:
		p.stype = typ
	case reflect.Ptr:
		if t1.Elem().Kind() == reflect.Struct {
			p.stype = t1.Elem()
		}
	case reflect.Slice:
		switch t2 := t1.Elem(); t2.Kind() {
		case reflect.Ptr:
			switch t3 := t2.Elem(); t3.Kind() {
			case reflect.Struct:
				p.stype = t3
			}
		case reflect.Struct:
			p.stype = t2
		}

	case reflect.Map:

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
	if p.stype != nil {
		if lockGetProp {
			p.sprop = GetProperties(p.stype)
		} else {
			p.sprop = getPropertiesLocked(p.stype)
		}
	}
}

var (
	marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()
)


func (p *Properties) Init(typ reflect.Type, name, tag string, f *reflect.StructField) {
	p.init(typ, name, tag, f, true)
}

func (p *Properties) init(typ reflect.Type, name, tag string, f *reflect.StructField, lockGetProp bool) {
	
	p.Name = name
	p.OrigName = name
	if tag == "" {
		return
	}
	p.Parse(tag)
	p.setFieldProps(typ, f, lockGetProp)
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

	
	prop.Prop = make([]*Properties, t.NumField())
	prop.order = make([]int, t.NumField())

	isOneofMessage := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		p := new(Properties)
		name := f.Name
		p.init(f.Type, name, f.Tag.Get("protobuf"), &f, false)

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
	}

	
	sort.Sort(prop)

	type oneofMessage interface {
		XXX_OneofFuncs() (func(Message, *Buffer) error, func(Message, int, int, *Buffer) (bool, error), func(Message) int, []interface{})
	}
	if om, ok := reflect.Zero(reflect.PtrTo(t)).Interface().(oneofMessage); isOneofMessage && ok {
		var oots []interface{}
		_, _, _, oots = om.XXX_OneofFuncs()

		
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
	protoTypedNils = make(map[string]Message)      
	protoMapTypes  = make(map[string]reflect.Type) 
	revProtoTypes  = make(map[reflect.Type]string)
)



func RegisterType(x Message, name string) {
	if _, ok := protoTypedNils[name]; ok {
		
		log.Printf("proto: duplicate proto type registered: %s", name)
		return
	}
	t := reflect.TypeOf(x)
	if v := reflect.ValueOf(x); v.Kind() == reflect.Ptr && v.Pointer() == 0 {
		
		
		protoTypedNils[name] = x
	} else {
		protoTypedNils[name] = reflect.Zero(t).Interface().(Message)
	}
	revProtoTypes[t] = name
}



func RegisterMapType(x interface{}, name string) {
	if reflect.TypeOf(x).Kind() != reflect.Map {
		panic(fmt.Sprintf("RegisterMapType(%T, %q); want map", x, name))
	}
	if _, ok := protoMapTypes[name]; ok {
		log.Printf("proto: duplicate proto type registered: %s", name)
		return
	}
	t := reflect.TypeOf(x)
	protoMapTypes[name] = t
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




func MessageType(name string) reflect.Type {
	if t, ok := protoTypedNils[name]; ok {
		return reflect.TypeOf(t)
	}
	return protoMapTypes[name]
}


var (
	protoFiles = make(map[string][]byte) 
)



func RegisterFile(filename string, fileDescriptor []byte) {
	protoFiles[filename] = fileDescriptor
}


func FileDescriptor(filename string) []byte { return protoFiles[filename] }
