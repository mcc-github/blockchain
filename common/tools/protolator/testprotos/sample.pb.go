


package testprotos 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



type SimpleMsg struct {
	PlainField           string            `protobuf:"bytes,1,opt,name=plain_field,json=plainField,proto3" json:"plain_field,omitempty"`
	MapField             map[string]string `protobuf:"bytes,2,rep,name=map_field,json=mapField,proto3" json:"map_field,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SliceField           []string          `protobuf:"bytes,3,rep,name=slice_field,json=sliceField,proto3" json:"slice_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SimpleMsg) Reset()         { *m = SimpleMsg{} }
func (m *SimpleMsg) String() string { return proto.CompactTextString(m) }
func (*SimpleMsg) ProtoMessage()    {}
func (*SimpleMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_1b65d606a8063655, []int{0}
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
	PlainNestedField     *SimpleMsg            `protobuf:"bytes,1,opt,name=plain_nested_field,json=plainNestedField,proto3" json:"plain_nested_field,omitempty"`
	MapNestedField       map[string]*SimpleMsg `protobuf:"bytes,2,rep,name=map_nested_field,json=mapNestedField,proto3" json:"map_nested_field,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SliceNestedField     []*SimpleMsg          `protobuf:"bytes,3,rep,name=slice_nested_field,json=sliceNestedField,proto3" json:"slice_nested_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *NestedMsg) Reset()         { *m = NestedMsg{} }
func (m *NestedMsg) String() string { return proto.CompactTextString(m) }
func (*NestedMsg) ProtoMessage()    {}
func (*NestedMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_1b65d606a8063655, []int{1}
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
	MapOpaqueField       map[string][]byte `protobuf:"bytes,2,rep,name=map_opaque_field,json=mapOpaqueField,proto3" json:"map_opaque_field,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SliceOpaqueField     [][]byte          `protobuf:"bytes,3,rep,name=slice_opaque_field,json=sliceOpaqueField,proto3" json:"slice_opaque_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *StaticallyOpaqueMsg) Reset()         { *m = StaticallyOpaqueMsg{} }
func (m *StaticallyOpaqueMsg) String() string { return proto.CompactTextString(m) }
func (*StaticallyOpaqueMsg) ProtoMessage()    {}
func (*StaticallyOpaqueMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_1b65d606a8063655, []int{2}
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
	OpaqueType           string            `protobuf:"bytes,1,opt,name=opaque_type,json=opaqueType,proto3" json:"opaque_type,omitempty"`
	PlainOpaqueField     []byte            `protobuf:"bytes,2,opt,name=plain_opaque_field,json=plainOpaqueField,proto3" json:"plain_opaque_field,omitempty"`
	MapOpaqueField       map[string][]byte `protobuf:"bytes,3,rep,name=map_opaque_field,json=mapOpaqueField,proto3" json:"map_opaque_field,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SliceOpaqueField     [][]byte          `protobuf:"bytes,4,rep,name=slice_opaque_field,json=sliceOpaqueField,proto3" json:"slice_opaque_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *VariablyOpaqueMsg) Reset()         { *m = VariablyOpaqueMsg{} }
func (m *VariablyOpaqueMsg) String() string { return proto.CompactTextString(m) }
func (*VariablyOpaqueMsg) ProtoMessage()    {}
func (*VariablyOpaqueMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_1b65d606a8063655, []int{3}
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
	DynamicType          string                     `protobuf:"bytes,1,opt,name=dynamic_type,json=dynamicType,proto3" json:"dynamic_type,omitempty"`
	PlainDynamicField    *ContextlessMsg            `protobuf:"bytes,2,opt,name=plain_dynamic_field,json=plainDynamicField,proto3" json:"plain_dynamic_field,omitempty"`
	MapDynamicField      map[string]*ContextlessMsg `protobuf:"bytes,3,rep,name=map_dynamic_field,json=mapDynamicField,proto3" json:"map_dynamic_field,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SliceDynamicField    []*ContextlessMsg          `protobuf:"bytes,4,rep,name=slice_dynamic_field,json=sliceDynamicField,proto3" json:"slice_dynamic_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *DynamicMsg) Reset()         { *m = DynamicMsg{} }
func (m *DynamicMsg) String() string { return proto.CompactTextString(m) }
func (*DynamicMsg) ProtoMessage()    {}
func (*DynamicMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_1b65d606a8063655, []int{4}
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
	return fileDescriptor_sample_1b65d606a8063655, []int{5}
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



type UnmarshalableDeepFields struct {
	PlainOpaqueField     []byte            `protobuf:"bytes,1,opt,name=plain_opaque_field,json=plainOpaqueField,proto3" json:"plain_opaque_field,omitempty"`
	MapOpaqueField       map[string][]byte `protobuf:"bytes,2,rep,name=map_opaque_field,json=mapOpaqueField,proto3" json:"map_opaque_field,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SliceOpaqueField     [][]byte          `protobuf:"bytes,3,rep,name=slice_opaque_field,json=sliceOpaqueField,proto3" json:"slice_opaque_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UnmarshalableDeepFields) Reset()         { *m = UnmarshalableDeepFields{} }
func (m *UnmarshalableDeepFields) String() string { return proto.CompactTextString(m) }
func (*UnmarshalableDeepFields) ProtoMessage()    {}
func (*UnmarshalableDeepFields) Descriptor() ([]byte, []int) {
	return fileDescriptor_sample_1b65d606a8063655, []int{6}
}
func (m *UnmarshalableDeepFields) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnmarshalableDeepFields.Unmarshal(m, b)
}
func (m *UnmarshalableDeepFields) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnmarshalableDeepFields.Marshal(b, m, deterministic)
}
func (dst *UnmarshalableDeepFields) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnmarshalableDeepFields.Merge(dst, src)
}
func (m *UnmarshalableDeepFields) XXX_Size() int {
	return xxx_messageInfo_UnmarshalableDeepFields.Size(m)
}
func (m *UnmarshalableDeepFields) XXX_DiscardUnknown() {
	xxx_messageInfo_UnmarshalableDeepFields.DiscardUnknown(m)
}

var xxx_messageInfo_UnmarshalableDeepFields proto.InternalMessageInfo

func (m *UnmarshalableDeepFields) GetPlainOpaqueField() []byte {
	if m != nil {
		return m.PlainOpaqueField
	}
	return nil
}

func (m *UnmarshalableDeepFields) GetMapOpaqueField() map[string][]byte {
	if m != nil {
		return m.MapOpaqueField
	}
	return nil
}

func (m *UnmarshalableDeepFields) GetSliceOpaqueField() [][]byte {
	if m != nil {
		return m.SliceOpaqueField
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
	proto.RegisterType((*UnmarshalableDeepFields)(nil), "testprotos.UnmarshalableDeepFields")
	proto.RegisterMapType((map[string][]byte)(nil), "testprotos.UnmarshalableDeepFields.MapOpaqueFieldEntry")
}

func init() { proto.RegisterFile("sample.proto", fileDescriptor_sample_1b65d606a8063655) }

var fileDescriptor_sample_1b65d606a8063655 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x94, 0xcb, 0x6e, 0xd3, 0x40,
	0x14, 0x86, 0x15, 0xbb, 0x20, 0x72, 0x1c, 0x4a, 0xe2, 0x14, 0x51, 0x65, 0xd3, 0x10, 0x36, 0x45,
	0xad, 0x62, 0x48, 0x16, 0x20, 0xd8, 0x94, 0xb6, 0xb0, 0x40, 0x2a, 0x48, 0x09, 0x37, 0x81, 0x00,
	0x4d, 0x9c, 0x69, 0x62, 0x31, 0xbe, 0xe0, 0x99, 0x20, 0xbc, 0xe3, 0x1d, 0x58, 0xf2, 0x12, 0x3c,
	0x04, 0x4b, 0x1e, 0x0a, 0xcd, 0x25, 0xcd, 0x19, 0x6a, 0x17, 0x21, 0x84, 0xd4, 0x55, 0xe2, 0xe3,
	0x39, 0xff, 0xf9, 0xff, 0xcf, 0x33, 0x03, 0x0d, 0x4e, 0xe2, 0x8c, 0xd1, 0x7e, 0x96, 0xa7, 0x22,
	0xf5, 0x41, 0x50, 0x2e, 0xd4, 0x5f, 0xde, 0xfb, 0x59, 0x83, 0xfa, 0x38, 0x92, 0x2f, 0x8f, 0xf8,
	0xcc, 0xdf, 0x02, 0x2f, 0x63, 0x24, 0x4a, 0xde, 0x1f, 0x47, 0x94, 0x4d, 0x37, 0x6b, 0xdd, 0xda,
	0x76, 0x7d, 0x04, 0xaa, 0xf4, 0x48, 0x56, 0xfc, 0x3d, 0xa8, 0xc7, 0x24, 0x33, 0xaf, 0x9d, 0xae,
	0xbb, 0xed, 0x0d, 0x6e, 0xf4, 0x57, 0x72, 0xfd, 0x13, 0xa9, 0xfe, 0x11, 0xc9, 0x54, 0xcb, 0xc3,
	0x44, 0xe4, 0xc5, 0xe8, 0x52, 0x6c, 0x1e, 0xe5, 0x08, 0xce, 0xa2, 0x90, 0x1a, 0x0d, 0xb7, 0xeb,
	0xca, 0x11, 0xaa, 0xa4, 0x16, 0x74, 0xee, 0xc3, 0x65, 0xab, 0xd7, 0x6f, 0x82, 0xfb, 0x81, 0x16,
	0xc6, 0x8c, 0xfc, 0xeb, 0x6f, 0xc0, 0x85, 0x4f, 0x84, 0x2d, 0xe8, 0xa6, 0xa3, 0x6a, 0xfa, 0xe1,
	0x9e, 0x73, 0xb7, 0xd6, 0xfb, 0xe1, 0x40, 0xfd, 0x09, 0xe5, 0x82, 0x4e, 0x65, 0x9c, 0x03, 0xf0,
	0x75, 0x9c, 0x44, 0x95, 0x50, 0x2a, 0x6f, 0x70, 0xb5, 0xd4, 0xf6, 0xa8, 0xa9, 0x1a, 0xb4, 0x84,
	0x36, 0x3c, 0x86, 0xa6, 0x8c, 0x6c, 0x49, 0xe8, 0xe4, 0x37, 0xb1, 0xc4, 0xc9, 0x54, 0x99, 0x1c,
	0xf5, 0xeb, 0xfc, 0xeb, 0xb1, 0x55, 0x94, 0xce, 0x34, 0x05, 0x4b, 0xd6, 0x55, 0xb2, 0x55, 0xce,
	0x54, 0x03, 0x12, 0xe9, 0xbc, 0x82, 0x76, 0xc9, 0xac, 0x12, 0x5e, 0x3b, 0x98, 0x57, 0xe5, 0x00,
	0x84, 0xf1, 0xab, 0x03, 0xed, 0xb1, 0x20, 0x22, 0x0a, 0x09, 0x63, 0xc5, 0xd3, 0x8c, 0x7c, 0x5c,
	0xa8, 0xfd, 0xb1, 0xbb, 0x04, 0x9a, 0xaa, 0x12, 0x02, 0xda, 0x30, 0xe4, 0xf4, 0x5a, 0x1d, 0xf2,
	0xad, 0x26, 0x67, 0xad, 0xd5, 0xe4, 0x86, 0x96, 0x83, 0xd3, 0x83, 0x24, 0x43, 0xa4, 0xb4, 0x62,
	0x88, 0xe5, 0x77, 0x97, 0x0c, 0xad, 0x01, 0x92, 0x61, 0xc3, 0xc0, 0x42, 0xab, 0x3b, 0x0f, 0x14,
	0xac, 0xdf, 0x45, 0xff, 0xb4, 0xb9, 0x1a, 0x98, 0xca, 0x77, 0x07, 0x5a, 0x2f, 0x48, 0x1e, 0x91,
	0x09, 0x66, 0xb2, 0x05, 0x9e, 0x31, 0x20, 0x8a, 0x8c, 0x2e, 0xcf, 0x8c, 0x2e, 0x3d, 0x2b, 0x32,
	0x5a, 0x01, 0xcd, 0xa9, 0x80, 0xf6, 0xa6, 0x04, 0x9a, 0xde, 0x17, 0xb7, 0x31, 0xb4, 0x53, 0x3e,
	0xfe, 0x01, 0xd9, 0xda, 0xff, 0x43, 0xf6, 0xc5, 0x05, 0x38, 0x2c, 0x12, 0x12, 0x47, 0xa1, 0x64,
	0x75, 0x1d, 0x1a, 0x53, 0xfd, 0x84, 0x61, 0x79, 0xa6, 0xa6, 0x68, 0x3d, 0x86, 0xb6, 0xa6, 0xb5,
	0x5c, 0xb8, 0xc2, 0xe5, 0x0d, 0x3a, 0x18, 0xc1, 0x41, 0x9a, 0x08, 0xfa, 0x59, 0x30, 0xca, 0xb9,
	0xdc, 0xbe, 0x2d, 0xd5, 0x66, 0x86, 0xe9, 0xb8, 0x2f, 0xa1, 0x25, 0x59, 0xda, 0x4a, 0x1a, 0xe6,
	0x0e, 0x56, 0x5a, 0x39, 0x94, 0x14, 0xb1, 0x84, 0xc6, 0x78, 0x25, 0xb6, 0xab, 0xd2, 0xa4, 0xe6,
	0x68, 0x4b, 0xaf, 0x29, 0xe9, 0x33, 0x4d, 0xaa, 0x36, 0xac, 0xd5, 0x79, 0x07, 0x1b, 0x65, 0x43,
	0x4b, 0x30, 0xdf, 0xb2, 0x8f, 0xf1, 0x59, 0x73, 0xd0, 0x27, 0x18, 0xc2, 0xba, 0xfd, 0x52, 0x7e,
	0x85, 0x92, 0xf3, 0x6b, 0x76, 0xb1, 0x72, 0xd0, 0xfb, 0xe6, 0xc0, 0xb5, 0xe7, 0x49, 0x4c, 0x72,
	0x3e, 0x27, 0x8c, 0x4c, 0x18, 0x3d, 0xa4, 0x54, 0xdf, 0xc9, 0xfc, 0x2f, 0x2f, 0x01, 0x52, 0x79,
	0x09, 0xdc, 0xc1, 0xfe, 0x2b, 0x86, 0x9d, 0xcb, 0x8b, 0x60, 0x7f, 0xff, 0xf5, 0xde, 0x2c, 0x12,
	0xf3, 0xc5, 0xa4, 0x1f, 0xa6, 0x71, 0x30, 0x2f, 0x32, 0x9a, 0x33, 0x3a, 0x9d, 0xd1, 0x3c, 0x38,
	0x26, 0x93, 0x3c, 0x0a, 0x83, 0x30, 0x8d, 0xe3, 0x34, 0x09, 0x44, 0x9a, 0x32, 0x1e, 0xa8, 0x84,
	0x8c, 0x88, 0x34, 0x0f, 0x56, 0x81, 0x27, 0x17, 0xd5, 0xef, 0xf0, 0x57, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xbc, 0xe8, 0xa5, 0x47, 0x9b, 0x07, 0x00, 0x00,
}
