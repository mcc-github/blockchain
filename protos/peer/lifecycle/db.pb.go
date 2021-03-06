


package lifecycle

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 




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
	return fileDescriptor_389b29a8aaf0aebb, []int{0}
}

func (m *StateMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateMetadata.Unmarshal(m, b)
}
func (m *StateMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateMetadata.Marshal(b, m, deterministic)
}
func (m *StateMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateMetadata.Merge(m, src)
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
	return fileDescriptor_389b29a8aaf0aebb, []int{1}
}

func (m *StateData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateData.Unmarshal(m, b)
}
func (m *StateData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateData.Marshal(b, m, deterministic)
}
func (m *StateData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateData.Merge(m, src)
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

type StateData_Int64 struct {
	Int64 int64 `protobuf:"varint,1,opt,name=Int64,proto3,oneof"`
}

type StateData_Bytes struct {
	Bytes []byte `protobuf:"bytes,2,opt,name=Bytes,proto3,oneof"`
}

type StateData_String_ struct {
	String_ string `protobuf:"bytes,3,opt,name=String,proto3,oneof"`
}

func (*StateData_Int64) isStateData_Type() {}

func (*StateData_Bytes) isStateData_Type() {}

func (*StateData_String_) isStateData_Type() {}

func (m *StateData) GetType() isStateData_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *StateData) GetInt64() int64 {
	if x, ok := m.GetType().(*StateData_Int64); ok {
		return x.Int64
	}
	return 0
}

func (m *StateData) GetBytes() []byte {
	if x, ok := m.GetType().(*StateData_Bytes); ok {
		return x.Bytes
	}
	return nil
}

func (m *StateData) GetString_() string {
	if x, ok := m.GetType().(*StateData_String_); ok {
		return x.String_
	}
	return ""
}


func (*StateData) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*StateData_Int64)(nil),
		(*StateData_Bytes)(nil),
		(*StateData_String_)(nil),
	}
}

func init() {
	proto.RegisterType((*StateMetadata)(nil), "lifecycle.StateMetadata")
	proto.RegisterType((*StateData)(nil), "lifecycle.StateData")
}

func init() { proto.RegisterFile("peer/lifecycle/db.proto", fileDescriptor_389b29a8aaf0aebb) }

var fileDescriptor_389b29a8aaf0aebb = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x1b, 0x02, 0x11, 0x39, 0xc1, 0x92, 0xa1, 0x44, 0x4c, 0x55, 0xa7, 0x0c, 0xc8, 0x1e,
	0x8a, 0xf8, 0x01, 0x81, 0xa1, 0x0c, 0x2c, 0x29, 0x13, 0x12, 0x83, 0xe3, 0x5c, 0x52, 0x4b, 0xa6,
	0xb6, 0xdc, 0x63, 0xf0, 0xbf, 0x47, 0xb6, 0xa3, 0x88, 0x4e, 0xd6, 0xf7, 0xa4, 0xef, 0x59, 0xef,
	0xe0, 0xc1, 0x22, 0x3a, 0xae, 0xd5, 0x88, 0xd2, 0x4b, 0x8d, 0x7c, 0xe8, 0x99, 0x75, 0x86, 0x4c,
	0x55, 0x2e, 0xd9, 0xf6, 0x15, 0xee, 0x0f, 0x24, 0x08, 0x3f, 0x90, 0xc4, 0x20, 0x48, 0x54, 0x8f,
	0x70, 0x1b, 0x5e, 0xf2, 0x16, 0xeb, 0x6c, 0x93, 0x35, 0x65, 0xb7, 0x70, 0xb5, 0x86, 0x62, 0x54,
	0xa8, 0x87, 0x73, 0x7d, 0xb5, 0xc9, 0x9b, 0xb2, 0x9b, 0x69, 0xfb, 0x0d, 0x65, 0x2c, 0x79, 0x0b,
	0x05, 0x6b, 0xb8, 0x79, 0x3f, 0xd1, 0xcb, 0x73, 0xb4, 0xf3, 0xfd, 0xaa, 0x4b, 0x18, 0xf2, 0xd6,
	0x13, 0x06, 0x37, 0x6b, 0xee, 0x42, 0x1e, 0xb1, 0xaa, 0xa1, 0x38, 0x90, 0x53, 0xa7, 0xa9, 0xce,
	0xc3, 0x77, 0xfb, 0x55, 0x37, 0x73, 0x5b, 0xc0, 0xf5, 0xa7, 0xb7, 0xd8, 0x4a, 0x78, 0x32, 0x6e,
	0x62, 0x47, 0x6f, 0xd1, 0x69, 0x1c, 0x26, 0x74, 0x6c, 0x14, 0xbd, 0x53, 0x32, 0xcd, 0x39, 0xb3,
	0xb0, 0x93, 0x2d, 0x9b, 0xbe, 0x76, 0x93, 0xa2, 0xe3, 0x6f, 0xcf, 0xa4, 0xf9, 0xe1, 0xff, 0x24,
	0x9e, 0x24, 0x9e, 0x24, 0x7e, 0x79, 0x9c, 0xbe, 0x88, 0xf1, 0xee, 0x2f, 0x00, 0x00, 0xff, 0xff,
	0xf3, 0x30, 0xef, 0x11, 0x35, 0x01, 0x00, 0x00,
}
