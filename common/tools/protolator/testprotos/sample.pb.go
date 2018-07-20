


package testprotos 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



type SimpleMsg struct {
	PlainField           string            `protobuf:"bytes,1,opt,name=plain_field,json=plainField" json:"plain_field,omitempty"`
	MapField             map[string]string `protobuf:"bytes,2,rep,name=map_field,json=mapField" json:"map_field,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	SliceField           []string          `protobuf:"bytes,3,rep,name=slice_field,json=sliceField" json:"slice_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SimpleMsg) Reset()         { *m = SimpleMsg{} }
func (m *SimpleMsg) String() string { return proto.CompactTextString(m) }
func (*SimpleMsg) ProtoMessage()    {}
func (*SimpleMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_fde2d2951f3e63f3, []int{0}
}
func (m *SimpleMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SimpleMsg.Unmarshal(m, b)
}
func (m *SimpleMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SimpleMsg.Marshal(b, m, deterministic)
}
func (dst *SimpleMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimpleMsg.Merge(dst, src)
}
func (m *SimpleMsg) XXX_Size() int {
	return xxx_messageInfo_SimpleMsg.Size(m)
}
func (m *SimpleMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_SimpleMsg.DiscardUnknown(m)
}

var xxx_messageInfo_SimpleMsg proto.InternalMessageInfo

func (m *SimpleMsg) GetPlainField() string {
	if m != nil {
		return m.PlainField
	}
	return ""
}

func (m *SimpleMsg) GetMapField() map[string]string {
	if m != nil {
		return m.MapField
	}
	return nil
}

func (m *SimpleMsg) GetSliceField() []string {
	if m != nil {
		return m.SliceField
	}
	return nil
}


type NestedMsg struct {
	PlainNestedField     *SimpleMsg            `protobuf:"bytes,1,opt,name=plain_nested_field,json=plainNestedField" json:"plain_nested_field,omitempty"`
	MapNestedField       map[string]*SimpleMsg `protobuf:"bytes,2,rep,name=map_nested_field,json=mapNestedField" json:"map_nested_field,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	SliceNestedField     []*SimpleMsg          `protobuf:"bytes,3,rep,name=slice_nested_field,json=sliceNestedField" json:"slice_nested_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *NestedMsg) Reset()         { *m = NestedMsg{} }
func (m *NestedMsg) String() string { return proto.CompactTextString(m) }
func (*NestedMsg) ProtoMessage()    {}
func (*NestedMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_fde2d2951f3e63f3, []int{1}
}
func (m *NestedMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NestedMsg.Unmarshal(m, b)
}
func (m *NestedMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NestedMsg.Marshal(b, m, deterministic)
}
func (dst *NestedMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NestedMsg.Merge(dst, src)
}
func (m *NestedMsg) XXX_Size() int {
	return xxx_messageInfo_NestedMsg.Size(m)
}
func (m *NestedMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_NestedMsg.DiscardUnknown(m)
}

var xxx_messageInfo_NestedMsg proto.InternalMessageInfo

func (m *NestedMsg) GetPlainNestedField() *SimpleMsg {
	if m != nil {
		return m.PlainNestedField
	}
	return nil
}

func (m *NestedMsg) GetMapNestedField() map[string]*SimpleMsg {
	if m != nil {
		return m.MapNestedField
	}
	return nil
}

func (m *NestedMsg) GetSliceNestedField() []*SimpleMsg {
	if m != nil {
		return m.SliceNestedField
	}
	return nil
}



type StaticallyOpaqueMsg struct {
	PlainOpaqueField     []byte            `protobuf:"bytes,1,opt,name=plain_opaque_field,json=plainOpaqueField,proto3" json:"plain_opaque_field,omitempty"`
	MapOpaqueField       map[string][]byte `protobuf:"bytes,2,rep,name=map_opaque_field,json=mapOpaqueField" json:"map_opaque_field,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SliceOpaqueField     [][]byte          `protobuf:"bytes,3,rep,name=slice_opaque_field,json=sliceOpaqueField,proto3" json:"slice_opaque_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *StaticallyOpaqueMsg) Reset()         { *m = StaticallyOpaqueMsg{} }
func (m *StaticallyOpaqueMsg) String() string { return proto.CompactTextString(m) }
func (*StaticallyOpaqueMsg) ProtoMessage()    {}
func (*StaticallyOpaqueMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_fde2d2951f3e63f3, []int{2}
}
func (m *StaticallyOpaqueMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StaticallyOpaqueMsg.Unmarshal(m, b)
}
func (m *StaticallyOpaqueMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StaticallyOpaqueMsg.Marshal(b, m, deterministic)
}
func (dst *StaticallyOpaqueMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StaticallyOpaqueMsg.Merge(dst, src)
}
func (m *StaticallyOpaqueMsg) XXX_Size() int {
	return xxx_messageInfo_StaticallyOpaqueMsg.Size(m)
}
func (m *StaticallyOpaqueMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_StaticallyOpaqueMsg.DiscardUnknown(m)
}

var xxx_messageInfo_StaticallyOpaqueMsg proto.InternalMessageInfo

func (m *StaticallyOpaqueMsg) GetPlainOpaqueField() []byte {
	if m != nil {
		return m.PlainOpaqueField
	}
	return nil
}

func (m *StaticallyOpaqueMsg) GetMapOpaqueField() map[string][]byte {
	if m != nil {
		return m.MapOpaqueField
	}
	return nil
}

func (m *StaticallyOpaqueMsg) GetSliceOpaqueField() [][]byte {
	if m != nil {
		return m.SliceOpaqueField
	}
	return nil
}



type VariablyOpaqueMsg struct {
	OpaqueType           string            `protobuf:"bytes,1,opt,name=opaque_type,json=opaqueType" json:"opaque_type,omitempty"`
	PlainOpaqueField     []byte            `protobuf:"bytes,2,opt,name=plain_opaque_field,json=plainOpaqueField,proto3" json:"plain_opaque_field,omitempty"`
	MapOpaqueField       map[string][]byte `protobuf:"bytes,3,rep,name=map_opaque_field,json=mapOpaqueField" json:"map_opaque_field,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SliceOpaqueField     [][]byte          `protobuf:"bytes,4,rep,name=slice_opaque_field,json=sliceOpaqueField,proto3" json:"slice_opaque_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *VariablyOpaqueMsg) Reset()         { *m = VariablyOpaqueMsg{} }
func (m *VariablyOpaqueMsg) String() string { return proto.CompactTextString(m) }
func (*VariablyOpaqueMsg) ProtoMessage()    {}
func (*VariablyOpaqueMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_fde2d2951f3e63f3, []int{3}
}
func (m *VariablyOpaqueMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VariablyOpaqueMsg.Unmarshal(m, b)
}
func (m *VariablyOpaqueMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VariablyOpaqueMsg.Marshal(b, m, deterministic)
}
func (dst *VariablyOpaqueMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VariablyOpaqueMsg.Merge(dst, src)
}
func (m *VariablyOpaqueMsg) XXX_Size() int {
	return xxx_messageInfo_VariablyOpaqueMsg.Size(m)
}
func (m *VariablyOpaqueMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_VariablyOpaqueMsg.DiscardUnknown(m)
}

var xxx_messageInfo_VariablyOpaqueMsg proto.InternalMessageInfo

func (m *VariablyOpaqueMsg) GetOpaqueType() string {
	if m != nil {
		return m.OpaqueType
	}
	return ""
}

func (m *VariablyOpaqueMsg) GetPlainOpaqueField() []byte {
	if m != nil {
		return m.PlainOpaqueField
	}
	return nil
}

func (m *VariablyOpaqueMsg) GetMapOpaqueField() map[string][]byte {
	if m != nil {
		return m.MapOpaqueField
	}
	return nil
}

func (m *VariablyOpaqueMsg) GetSliceOpaqueField() [][]byte {
	if m != nil {
		return m.SliceOpaqueField
	}
	return nil
}




type DynamicMsg struct {
	DynamicType          string                     `protobuf:"bytes,1,opt,name=dynamic_type,json=dynamicType" json:"dynamic_type,omitempty"`
	PlainDynamicField    *ContextlessMsg            `protobuf:"bytes,2,opt,name=plain_dynamic_field,json=plainDynamicField" json:"plain_dynamic_field,omitempty"`
	MapDynamicField      map[string]*ContextlessMsg `protobuf:"bytes,3,rep,name=map_dynamic_field,json=mapDynamicField" json:"map_dynamic_field,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	SliceDynamicField    []*ContextlessMsg          `protobuf:"bytes,4,rep,name=slice_dynamic_field,json=sliceDynamicField" json:"slice_dynamic_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *DynamicMsg) Reset()         { *m = DynamicMsg{} }
func (m *DynamicMsg) String() string { return proto.CompactTextString(m) }
func (*DynamicMsg) ProtoMessage()    {}
func (*DynamicMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_fde2d2951f3e63f3, []int{4}
}
func (m *DynamicMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DynamicMsg.Unmarshal(m, b)
}
func (m *DynamicMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DynamicMsg.Marshal(b, m, deterministic)
}
func (dst *DynamicMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DynamicMsg.Merge(dst, src)
}
func (m *DynamicMsg) XXX_Size() int {
	return xxx_messageInfo_DynamicMsg.Size(m)
}
func (m *DynamicMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_DynamicMsg.DiscardUnknown(m)
}

var xxx_messageInfo_DynamicMsg proto.InternalMessageInfo

func (m *DynamicMsg) GetDynamicType() string {
	if m != nil {
		return m.DynamicType
	}
	return ""
}

func (m *DynamicMsg) GetPlainDynamicField() *ContextlessMsg {
	if m != nil {
		return m.PlainDynamicField
	}
	return nil
}

func (m *DynamicMsg) GetMapDynamicField() map[string]*ContextlessMsg {
	if m != nil {
		return m.MapDynamicField
	}
	return nil
}

func (m *DynamicMsg) GetSliceDynamicField() []*ContextlessMsg {
	if m != nil {
		return m.SliceDynamicField
	}
	return nil
}




type ContextlessMsg struct {
	OpaqueField          []byte   `protobuf:"bytes,1,opt,name=opaque_field,json=opaqueField,proto3" json:"opaque_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContextlessMsg) Reset()         { *m = ContextlessMsg{} }
func (m *ContextlessMsg) String() string { return proto.CompactTextString(m) }
func (*ContextlessMsg) ProtoMessage()    {}
func (*ContextlessMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_fde2d2951f3e63f3, []int{5}
}
func (m *ContextlessMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContextlessMsg.Unmarshal(m, b)
}
func (m *ContextlessMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContextlessMsg.Marshal(b, m, deterministic)
}
func (dst *ContextlessMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContextlessMsg.Merge(dst, src)
}
func (m *ContextlessMsg) XXX_Size() int {
	return xxx_messageInfo_ContextlessMsg.Size(m)
}
func (m *ContextlessMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ContextlessMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ContextlessMsg proto.InternalMessageInfo

func (m *ContextlessMsg) GetOpaqueField() []byte {
	if m != nil {
		return m.OpaqueField
	}
	return nil
}

func init() {
	proto.RegisterType((*SimpleMsg)(nil), "testprotos.SimpleMsg")
	proto.RegisterMapType((map[string]string)(nil), "testprotos.SimpleMsg.MapFieldEntry")
	proto.RegisterType((*NestedMsg)(nil), "testprotos.NestedMsg")
	proto.RegisterMapType((map[string]*SimpleMsg)(nil), "testprotos.NestedMsg.MapNestedFieldEntry")
	proto.RegisterType((*StaticallyOpaqueMsg)(nil), "testprotos.StaticallyOpaqueMsg")
	proto.RegisterMapType((map[string][]byte)(nil), "testprotos.StaticallyOpaqueMsg.MapOpaqueFieldEntry")
	proto.RegisterType((*VariablyOpaqueMsg)(nil), "testprotos.VariablyOpaqueMsg")
	proto.RegisterMapType((map[string][]byte)(nil), "testprotos.VariablyOpaqueMsg.MapOpaqueFieldEntry")
	proto.RegisterType((*DynamicMsg)(nil), "testprotos.DynamicMsg")
	proto.RegisterMapType((map[string]*ContextlessMsg)(nil), "testprotos.DynamicMsg.MapDynamicFieldEntry")
	proto.RegisterType((*ContextlessMsg)(nil), "testprotos.ContextlessMsg")
}

func init() { proto.RegisterFile("sample.proto", fileDescriptor_sample_fde2d2951f3e63f3) }

var fileDescriptor_sample_fde2d2951f3e63f3 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0xd5, 0x64, 0x20, 0xfa, 0x12, 0x46, 0x9b, 0x0e, 0x69, 0xca, 0x65, 0x65, 0x5c, 0x86,
	0x36, 0x25, 0xd0, 0x5e, 0x10, 0x5c, 0xc6, 0x06, 0x1c, 0x90, 0x06, 0x52, 0x8b, 0x00, 0x81, 0x00,
	0xb9, 0xa9, 0xd7, 0x45, 0x38, 0x71, 0x48, 0x5c, 0x44, 0x6e, 0x7c, 0x07, 0xbe, 0x08, 0x1f, 0x82,
	0x23, 0x1f, 0x0a, 0xf9, 0xd9, 0x5d, 0x6c, 0x96, 0x8e, 0x03, 0xe2, 0xd4, 0xfa, 0xd9, 0xef, 0xff,
	0xff, 0xbf, 0x5f, 0x6c, 0xf0, 0x2b, 0x92, 0x15, 0x8c, 0x46, 0x45, 0xc9, 0x05, 0x0f, 0x40, 0xd0,
	0x4a, 0xe0, 0xdf, 0x6a, 0xf7, 0x57, 0x07, 0xba, 0xd3, 0x54, 0x6e, 0x9e, 0x54, 0x8b, 0x60, 0x07,
	0xbc, 0x82, 0x91, 0x34, 0xff, 0x78, 0x9a, 0x52, 0x36, 0xdf, 0xee, 0x0c, 0x3b, 0x7b, 0xdd, 0x09,
	0x60, 0xe9, 0xa9, 0xac, 0x04, 0x87, 0xd0, 0xcd, 0x48, 0xa1, 0xb7, 0x9d, 0xa1, 0xbb, 0xe7, 0x8d,
	0x6e, 0x47, 0x8d, 0x5c, 0x74, 0x2e, 0x15, 0x9d, 0x90, 0x02, 0x5b, 0x9e, 0xe4, 0xa2, 0xac, 0x27,
	0xd7, 0x32, 0xbd, 0x94, 0x16, 0x15, 0x4b, 0x13, 0xaa, 0x35, 0xdc, 0xa1, 0x2b, 0x2d, 0xb0, 0x84,
	0x07, 0xc2, 0x87, 0x70, 0xdd, 0xea, 0x0d, 0x7a, 0xe0, 0x7e, 0xa2, 0xb5, 0x0e, 0x23, 0xff, 0x06,
	0x5b, 0x70, 0xe5, 0x0b, 0x61, 0x4b, 0xba, 0xed, 0x60, 0x4d, 0x2d, 0x1e, 0x38, 0xf7, 0x3b, 0xbb,
	0x3f, 0x1d, 0xe8, 0x3e, 0xa7, 0x95, 0xa0, 0x73, 0x39, 0xce, 0x31, 0x04, 0x6a, 0x9c, 0x1c, 0x4b,
	0xc6, 0x54, 0xde, 0xe8, 0x66, 0x6b, 0xec, 0x49, 0x0f, 0x1b, 0x94, 0x84, 0x0a, 0x3c, 0x85, 0x9e,
	0x1c, 0xd9, 0x92, 0x50, 0x93, 0xdf, 0x31, 0x25, 0xce, 0x5d, 0xe5, 0xe4, 0x46, 0xbf, 0x9a, 0x7f,
	0x33, 0xb3, 0x8a, 0x32, 0x99, 0xa2, 0x60, 0xc9, 0xba, 0x28, 0xbb, 0x2e, 0x19, 0x36, 0x18, 0x22,
	0xe1, 0x1b, 0x18, 0xb4, 0x78, 0xb5, 0xf0, 0xda, 0x37, 0x79, 0xad, 0x35, 0x30, 0x30, 0x7e, 0x77,
	0x60, 0x30, 0x15, 0x44, 0xa4, 0x09, 0x61, 0xac, 0x7e, 0x51, 0x90, 0xcf, 0x4b, 0xbc, 0x1f, 0x07,
	0x2b, 0xa0, 0x1c, 0x4b, 0x06, 0x50, 0x5f, 0x93, 0x53, 0x67, 0xd5, 0x90, 0xef, 0x15, 0x39, 0xeb,
	0xac, 0x22, 0x37, 0xb6, 0x12, 0x5c, 0x34, 0x92, 0x0c, 0x0d, 0xa5, 0x86, 0xa1, 0x29, 0x7f, 0xb0,
	0x62, 0x68, 0x19, 0x48, 0x86, 0xbe, 0x86, 0x65, 0x9c, 0x0e, 0x1f, 0x21, 0xac, 0x3f, 0x45, 0xff,
	0x76, 0xb9, 0x7c, 0x93, 0xca, 0x0f, 0x07, 0xfa, 0xaf, 0x48, 0x99, 0x92, 0x99, 0xc9, 0x64, 0x07,
	0x3c, 0x1d, 0x40, 0xd4, 0x05, 0x5d, 0xbd, 0x19, 0x55, 0x7a, 0x59, 0x17, 0x74, 0x0d, 0x34, 0x67,
	0x0d, 0xb4, 0x77, 0x2d, 0xd0, 0xd4, 0xbd, 0xb8, 0x67, 0x42, 0xbb, 0x90, 0xe3, 0x1f, 0x90, 0x6d,
	0xfc, 0x3f, 0x64, 0xdf, 0x5c, 0x80, 0xc7, 0x75, 0x4e, 0xb2, 0x34, 0x91, 0xac, 0x6e, 0x81, 0x3f,
	0x57, 0x2b, 0x13, 0x96, 0xa7, 0x6b, 0x48, 0xeb, 0x19, 0x0c, 0x14, 0xad, 0xd5, 0xc1, 0x06, 0x97,
	0x37, 0x0a, 0x4d, 0x04, 0xc7, 0x3c, 0x17, 0xf4, 0xab, 0x60, 0xb4, 0xaa, 0xe4, 0xf5, 0xed, 0x63,
	0x9b, 0x36, 0x53, 0xe3, 0xbe, 0x86, 0xbe, 0x64, 0x69, 0x2b, 0x29, 0x98, 0xfb, 0xa6, 0x52, 0x93,
	0x50, 0x52, 0x34, 0x25, 0x14, 0xc6, 0x1b, 0x99, 0x5d, 0x95, 0x21, 0x15, 0x47, 0x5b, 0x7a, 0x03,
	0xa5, 0x2f, 0x0d, 0x89, 0x6d, 0xa6, 0x56, 0xf8, 0x01, 0xb6, 0xda, 0x4c, 0x5b, 0x30, 0xdf, 0xb5,
	0x9f, 0xf1, 0x65, 0x3e, 0xc6, 0x27, 0x18, 0xc3, 0xa6, 0xbd, 0x29, 0xbf, 0x42, 0xcb, 0xfb, 0xd5,
	0xb7, 0x18, 0x13, 0x1c, 0x1d, 0xbd, 0x3d, 0x5c, 0xa4, 0xe2, 0x6c, 0x39, 0x8b, 0x12, 0x9e, 0xc5,
	0x67, 0x75, 0x41, 0x4b, 0x46, 0xe7, 0x0b, 0x5a, 0xc6, 0xa7, 0x64, 0x56, 0xa6, 0x49, 0x9c, 0xf0,
	0x2c, 0xe3, 0x79, 0x2c, 0x38, 0x67, 0x55, 0x8c, 0x19, 0x18, 0x11, 0xbc, 0x8c, 0x9b, 0x48, 0xb3,
	0xab, 0xf8, 0x3b, 0xfe, 0x1d, 0x00, 0x00, 0xff, 0xff, 0xb4, 0x56, 0x77, 0x65, 0x7d, 0x06, 0x00,
	0x00,
}
