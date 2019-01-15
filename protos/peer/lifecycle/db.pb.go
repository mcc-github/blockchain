


package lifecycle 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 




type StateMetadata struct {
	Datatype             string   `protobuf:"bytes,1,opt,name=datatype,proto3" json:"datatype,omitempty"`
	Fields               []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StateMetadata) Reset()         { *m = StateMetadata{} }
func (m *StateMetadata) String() string { return proto.CompactTextString(m) }
func (*StateMetadata) ProtoMessage()    {}
func (*StateMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_db_7342255f8aba90fe, []int{0}
}
func (m *StateMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateMetadata.Unmarshal(m, b)
}
func (m *StateMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateMetadata.Marshal(b, m, deterministic)
}
func (dst *StateMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateMetadata.Merge(dst, src)
}
func (m *StateMetadata) XXX_Size() int {
	return xxx_messageInfo_StateMetadata.Size(m)
}
func (m *StateMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_StateMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_StateMetadata proto.InternalMessageInfo

func (m *StateMetadata) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *StateMetadata) GetFields() []string {
	if m != nil {
		return m.Fields
	}
	return nil
}


type StateData struct {
	
	
	
	
	
	Type                 isStateData_Type `protobuf_oneof:"Type"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *StateData) Reset()         { *m = StateData{} }
func (m *StateData) String() string { return proto.CompactTextString(m) }
func (*StateData) ProtoMessage()    {}
func (*StateData) Descriptor() ([]byte, []int) {
	return fileDescriptor_db_7342255f8aba90fe, []int{1}
}
func (m *StateData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateData.Unmarshal(m, b)
}
func (m *StateData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateData.Marshal(b, m, deterministic)
}
func (dst *StateData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateData.Merge(dst, src)
}
func (m *StateData) XXX_Size() int {
	return xxx_messageInfo_StateData.Size(m)
}
func (m *StateData) XXX_DiscardUnknown() {
	xxx_messageInfo_StateData.DiscardUnknown(m)
}

var xxx_messageInfo_StateData proto.InternalMessageInfo

type isStateData_Type interface {
	isStateData_Type()
}

type StateData_String_ struct {
	String_ string `protobuf:"bytes,1,opt,name=String,proto3,oneof"`
}

type StateData_Bytes struct {
	Bytes []byte `protobuf:"bytes,2,opt,name=Bytes,proto3,oneof"`
}

type StateData_Uint64 struct {
	Uint64 uint64 `protobuf:"varint,3,opt,name=Uint64,proto3,oneof"`
}

type StateData_Int64 struct {
	Int64 int64 `protobuf:"varint,4,opt,name=Int64,proto3,oneof"`
}

func (*StateData_String_) isStateData_Type() {}

func (*StateData_Bytes) isStateData_Type() {}

func (*StateData_Uint64) isStateData_Type() {}

func (*StateData_Int64) isStateData_Type() {}

func (m *StateData) GetType() isStateData_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *StateData) GetString_() string {
	if x, ok := m.GetType().(*StateData_String_); ok {
		return x.String_
	}
	return ""
}

func (m *StateData) GetBytes() []byte {
	if x, ok := m.GetType().(*StateData_Bytes); ok {
		return x.Bytes
	}
	return nil
}

func (m *StateData) GetUint64() uint64 {
	if x, ok := m.GetType().(*StateData_Uint64); ok {
		return x.Uint64
	}
	return 0
}

func (m *StateData) GetInt64() int64 {
	if x, ok := m.GetType().(*StateData_Int64); ok {
		return x.Int64
	}
	return 0
}


func (*StateData) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _StateData_OneofMarshaler, _StateData_OneofUnmarshaler, _StateData_OneofSizer, []interface{}{
		(*StateData_String_)(nil),
		(*StateData_Bytes)(nil),
		(*StateData_Uint64)(nil),
		(*StateData_Int64)(nil),
	}
}

func _StateData_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*StateData)
	
	switch x := m.Type.(type) {
	case *StateData_String_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.String_)
	case *StateData_Bytes:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeRawBytes(x.Bytes)
	case *StateData_Uint64:
		b.EncodeVarint(3<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.Uint64))
	case *StateData_Int64:
		b.EncodeVarint(4<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.Int64))
	case nil:
	default:
		return fmt.Errorf("StateData.Type has unexpected type %T", x)
	}
	return nil
}

func _StateData_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*StateData)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Type = &StateData_String_{x}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.Type = &StateData_Bytes{x}
		return true, err
	case 3: 
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Type = &StateData_Uint64{x}
		return true, err
	case 4: 
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Type = &StateData_Int64{int64(x)}
		return true, err
	default:
		return false, nil
	}
}

func _StateData_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*StateData)
	
	switch x := m.Type.(type) {
	case *StateData_String_:
		n += 1 
		n += proto.SizeVarint(uint64(len(x.String_)))
		n += len(x.String_)
	case *StateData_Bytes:
		n += 1 
		n += proto.SizeVarint(uint64(len(x.Bytes)))
		n += len(x.Bytes)
	case *StateData_Uint64:
		n += 1 
		n += proto.SizeVarint(uint64(x.Uint64))
	case *StateData_Int64:
		n += 1 
		n += proto.SizeVarint(uint64(x.Int64))
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*StateMetadata)(nil), "lifecycle.StateMetadata")
	proto.RegisterType((*StateData)(nil), "lifecycle.StateData")
}

func init() { proto.RegisterFile("peer/lifecycle/db.proto", fileDescriptor_db_7342255f8aba90fe) }

var fileDescriptor_db_7342255f8aba90fe = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x6b, 0x12, 0x22, 0x72, 0x82, 0x25, 0x43, 0x89, 0x98, 0xa2, 0x4e, 0x1e, 0x90, 0x3d,
	0x14, 0xf1, 0x03, 0x02, 0x43, 0x19, 0x58, 0x52, 0x58, 0xd8, 0x1c, 0xe7, 0x92, 0x5a, 0x0a, 0x8d,
	0xe5, 0x1e, 0x42, 0xfe, 0xf7, 0xc8, 0x71, 0x15, 0xd1, 0xc9, 0xfa, 0xce, 0xfa, 0x9e, 0xf4, 0x1e,
	0xdc, 0x5b, 0x44, 0x27, 0x47, 0xd3, 0xa3, 0xf6, 0x7a, 0x44, 0xd9, 0xb5, 0xc2, 0xba, 0x89, 0xa6,
	0x22, 0x5f, 0x6e, 0x9b, 0x17, 0xb8, 0xdb, 0x93, 0x22, 0x7c, 0x47, 0x52, 0x9d, 0x22, 0x55, 0x3c,
	0xc0, 0x4d, 0x78, 0xc9, 0x5b, 0x2c, 0x59, 0xc5, 0x78, 0xde, 0x2c, 0x5c, 0xac, 0x21, 0xeb, 0x0d,
	0x8e, 0xdd, 0xa9, 0xbc, 0xaa, 0x12, 0x9e, 0x37, 0x67, 0xda, 0xfc, 0x42, 0x3e, 0x87, 0xbc, 0x86,
	0x80, 0x12, 0xb2, 0x3d, 0x39, 0x73, 0x1c, 0xa2, 0xbe, 0x5b, 0x35, 0x67, 0x2e, 0xd6, 0x70, 0x5d,
	0x7b, 0xc2, 0x60, 0x33, 0x7e, 0xbb, 0x5b, 0x35, 0x11, 0x83, 0xf1, 0x69, 0x8e, 0xf4, 0xfc, 0x54,
	0x26, 0x15, 0xe3, 0x69, 0x30, 0x22, 0x07, 0xe3, 0x6d, 0xfe, 0x48, 0x2b, 0xc6, 0x93, 0x60, 0xcc,
	0x58, 0x67, 0x90, 0x7e, 0x78, 0x8b, 0xb5, 0x86, 0xc7, 0xc9, 0x0d, 0xe2, 0xe0, 0x2d, 0xba, 0x11,
	0xbb, 0x01, 0x9d, 0xe8, 0x55, 0xeb, 0x8c, 0x8e, 0x45, 0x4f, 0x22, 0x2c, 0x20, 0x96, 0xb6, 0x5f,
	0xdb, 0xc1, 0xd0, 0xe1, 0xa7, 0x15, 0x7a, 0xfa, 0x96, 0xff, 0x24, 0x19, 0x25, 0x19, 0x25, 0x79,
	0x39, 0x5b, 0x9b, 0xcd, 0xe7, 0xed, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1b, 0x93, 0xac, 0x49,
	0x4f, 0x01, 0x00, 0x00,
}
