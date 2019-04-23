


package msp 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 




type SerializedIdentity struct {
	
	Mspid string `protobuf:"bytes,1,opt,name=mspid,proto3" json:"mspid,omitempty"`
	
	IdBytes              []byte   `protobuf:"bytes,2,opt,name=id_bytes,json=idBytes,proto3" json:"id_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SerializedIdentity) Reset()         { *m = SerializedIdentity{} }
func (m *SerializedIdentity) String() string { return proto.CompactTextString(m) }
func (*SerializedIdentity) ProtoMessage()    {}
func (*SerializedIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_identities_8a64c0f25eb30ee4, []int{0}
}
func (m *SerializedIdentity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SerializedIdentity.Unmarshal(m, b)
}
func (m *SerializedIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SerializedIdentity.Marshal(b, m, deterministic)
}
func (dst *SerializedIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SerializedIdentity.Merge(dst, src)
}
func (m *SerializedIdentity) XXX_Size() int {
	return xxx_messageInfo_SerializedIdentity.Size(m)
}
func (m *SerializedIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_SerializedIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_SerializedIdentity proto.InternalMessageInfo

func (m *SerializedIdentity) GetMspid() string {
	if m != nil {
		return m.Mspid
	}
	return ""
}

func (m *SerializedIdentity) GetIdBytes() []byte {
	if m != nil {
		return m.IdBytes
	}
	return nil
}





type SerializedIdemixIdentity struct {
	
	
	
	NymX []byte `protobuf:"bytes,1,opt,name=nym_x,json=nymX,proto3" json:"nym_x,omitempty"`
	
	
	
	NymY []byte `protobuf:"bytes,2,opt,name=nym_y,json=nymY,proto3" json:"nym_y,omitempty"`
	
	Ou []byte `protobuf:"bytes,3,opt,name=ou,proto3" json:"ou,omitempty"`
	
	Role []byte `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`
	
	Proof                []byte   `protobuf:"bytes,5,opt,name=proof,proto3" json:"proof,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SerializedIdemixIdentity) Reset()         { *m = SerializedIdemixIdentity{} }
func (m *SerializedIdemixIdentity) String() string { return proto.CompactTextString(m) }
func (*SerializedIdemixIdentity) ProtoMessage()    {}
func (*SerializedIdemixIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_identities_8a64c0f25eb30ee4, []int{1}
}
func (m *SerializedIdemixIdentity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SerializedIdemixIdentity.Unmarshal(m, b)
}
func (m *SerializedIdemixIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SerializedIdemixIdentity.Marshal(b, m, deterministic)
}
func (dst *SerializedIdemixIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SerializedIdemixIdentity.Merge(dst, src)
}
func (m *SerializedIdemixIdentity) XXX_Size() int {
	return xxx_messageInfo_SerializedIdemixIdentity.Size(m)
}
func (m *SerializedIdemixIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_SerializedIdemixIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_SerializedIdemixIdentity proto.InternalMessageInfo

func (m *SerializedIdemixIdentity) GetNymX() []byte {
	if m != nil {
		return m.NymX
	}
	return nil
}

func (m *SerializedIdemixIdentity) GetNymY() []byte {
	if m != nil {
		return m.NymY
	}
	return nil
}

func (m *SerializedIdemixIdentity) GetOu() []byte {
	if m != nil {
		return m.Ou
	}
	return nil
}

func (m *SerializedIdemixIdentity) GetRole() []byte {
	if m != nil {
		return m.Role
	}
	return nil
}

func (m *SerializedIdemixIdentity) GetProof() []byte {
	if m != nil {
		return m.Proof
	}
	return nil
}

func init() {
	proto.RegisterType((*SerializedIdentity)(nil), "msp.SerializedIdentity")
	proto.RegisterType((*SerializedIdemixIdentity)(nil), "msp.SerializedIdemixIdentity")
}

func init() { proto.RegisterFile("msp/identities.proto", fileDescriptor_identities_8a64c0f25eb30ee4) }

var fileDescriptor_identities_8a64c0f25eb30ee4 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0x47, 0x95, 0x34, 0xe1, 0x8f, 0x55, 0x31, 0x98, 0x0e, 0x66, 0x2b, 0x9d, 0x32, 0xc5, 0x03,
	0xdf, 0xa0, 0x12, 0x03, 0x03, 0x4b, 0x58, 0x80, 0xa5, 0x6a, 0xea, 0x6b, 0x7a, 0x52, 0x2e, 0x67,
	0xd9, 0x8e, 0x54, 0x33, 0xf0, 0xd9, 0x51, 0x62, 0x40, 0xb0, 0xdd, 0xef, 0xe9, 0xe9, 0xc9, 0x16,
	0x2b, 0xf2, 0x56, 0xa3, 0x81, 0x21, 0x60, 0x40, 0xf0, 0xb5, 0x75, 0x1c, 0x58, 0x2e, 0xc8, 0xdb,
	0xcd, 0xa3, 0x90, 0x2f, 0xe0, 0x70, 0xdf, 0xe3, 0x07, 0x98, 0xa7, 0xa4, 0x44, 0xb9, 0x12, 0x25,
	0x79, 0x8b, 0x46, 0x65, 0xeb, 0xac, 0xba, 0x6e, 0xd2, 0x90, 0x77, 0xe2, 0x0a, 0xcd, 0xae, 0x8d,
	0x01, 0xbc, 0xca, 0xd7, 0x59, 0xb5, 0x6c, 0x2e, 0xd1, 0x6c, 0xa7, 0xb9, 0xf9, 0x14, 0xea, 0x5f,
	0x86, 0xf0, 0xfc, 0x1b, 0xbb, 0x15, 0xe5, 0x10, 0x69, 0x77, 0x9e, 0x63, 0xcb, 0xa6, 0x18, 0x22,
	0xbd, 0xfe, 0xc0, 0xf8, 0x1d, 0x9a, 0xe0, 0x9b, 0xbc, 0x11, 0x39, 0x8f, 0x6a, 0x31, 0x93, 0x9c,
	0x47, 0x29, 0x45, 0xe1, 0xb8, 0x07, 0x55, 0x24, 0x67, 0xba, 0xa7, 0xa7, 0x59, 0xc7, 0x7c, 0x54,
	0xe5, 0x0c, 0xd3, 0xd8, 0x3e, 0x8b, 0x7b, 0x76, 0x5d, 0x7d, 0x8a, 0x16, 0x5c, 0x0f, 0xa6, 0x03,
	0x57, 0x1f, 0xf7, 0xad, 0xc3, 0x43, 0xfa, 0xab, 0xaf, 0xc9, 0xdb, 0xf7, 0xaa, 0xc3, 0x70, 0x1a,
	0xdb, 0xfa, 0xc0, 0xa4, 0xff, 0x98, 0x3a, 0x99, 0x3a, 0x99, 0x9a, 0xbc, 0x6d, 0x2f, 0xe6, 0xfb,
	0xe1, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x13, 0xdc, 0xc8, 0x62, 0x39, 0x01, 0x00, 0x00,
}
