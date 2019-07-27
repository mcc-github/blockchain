


package msp

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 

type MSPPrincipal_Classification int32

const (
	MSPPrincipal_ROLE MSPPrincipal_Classification = 0
	
	
	MSPPrincipal_ORGANIZATION_UNIT MSPPrincipal_Classification = 1
	
	
	
	MSPPrincipal_IDENTITY MSPPrincipal_Classification = 2
	
	MSPPrincipal_ANONYMITY MSPPrincipal_Classification = 3
	
	MSPPrincipal_COMBINED MSPPrincipal_Classification = 4
)

var MSPPrincipal_Classification_name = map[int32]string{
	0: "ROLE",
	1: "ORGANIZATION_UNIT",
	2: "IDENTITY",
	3: "ANONYMITY",
	4: "COMBINED",
}

var MSPPrincipal_Classification_value = map[string]int32{
	"ROLE":              0,
	"ORGANIZATION_UNIT": 1,
	"IDENTITY":          2,
	"ANONYMITY":         3,
	"COMBINED":          4,
}

func (x MSPPrincipal_Classification) String() string {
	return proto.EnumName(MSPPrincipal_Classification_name, int32(x))
}

func (MSPPrincipal_Classification) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_82e08b7ead29bd48, []int{0, 0}
}

type MSPRole_MSPRoleType int32

const (
	MSPRole_MEMBER  MSPRole_MSPRoleType = 0
	MSPRole_ADMIN   MSPRole_MSPRoleType = 1
	MSPRole_CLIENT  MSPRole_MSPRoleType = 2
	MSPRole_PEER    MSPRole_MSPRoleType = 3
	MSPRole_ORDERER MSPRole_MSPRoleType = 4
)

var MSPRole_MSPRoleType_name = map[int32]string{
	0: "MEMBER",
	1: "ADMIN",
	2: "CLIENT",
	3: "PEER",
	4: "ORDERER",
}

var MSPRole_MSPRoleType_value = map[string]int32{
	"MEMBER":  0,
	"ADMIN":   1,
	"CLIENT":  2,
	"PEER":    3,
	"ORDERER": 4,
}

func (x MSPRole_MSPRoleType) String() string {
	return proto.EnumName(MSPRole_MSPRoleType_name, int32(x))
}

func (MSPRole_MSPRoleType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_82e08b7ead29bd48, []int{2, 0}
}

type MSPIdentityAnonymity_MSPIdentityAnonymityType int32

const (
	MSPIdentityAnonymity_NOMINAL   MSPIdentityAnonymity_MSPIdentityAnonymityType = 0
	MSPIdentityAnonymity_ANONYMOUS MSPIdentityAnonymity_MSPIdentityAnonymityType = 1
)

var MSPIdentityAnonymity_MSPIdentityAnonymityType_name = map[int32]string{
	0: "NOMINAL",
	1: "ANONYMOUS",
}

var MSPIdentityAnonymity_MSPIdentityAnonymityType_value = map[string]int32{
	"NOMINAL":   0,
	"ANONYMOUS": 1,
}

func (x MSPIdentityAnonymity_MSPIdentityAnonymityType) String() string {
	return proto.EnumName(MSPIdentityAnonymity_MSPIdentityAnonymityType_name, int32(x))
}

func (MSPIdentityAnonymity_MSPIdentityAnonymityType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_82e08b7ead29bd48, []int{3, 0}
}



















type MSPPrincipal struct {
	
	
	
	
	
	
	
	PrincipalClassification MSPPrincipal_Classification `protobuf:"varint,1,opt,name=principal_classification,json=principalClassification,proto3,enum=common.MSPPrincipal_Classification" json:"principal_classification,omitempty"`
	
	
	
	
	
	
	
	Principal            []byte   `protobuf:"bytes,2,opt,name=principal,proto3" json:"principal,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MSPPrincipal) Reset()         { *m = MSPPrincipal{} }
func (m *MSPPrincipal) String() string { return proto.CompactTextString(m) }
func (*MSPPrincipal) ProtoMessage()    {}
func (*MSPPrincipal) Descriptor() ([]byte, []int) {
	return fileDescriptor_82e08b7ead29bd48, []int{0}
}

func (m *MSPPrincipal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MSPPrincipal.Unmarshal(m, b)
}
func (m *MSPPrincipal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MSPPrincipal.Marshal(b, m, deterministic)
}
func (m *MSPPrincipal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MSPPrincipal.Merge(m, src)
}
func (m *MSPPrincipal) XXX_Size() int {
	return xxx_messageInfo_MSPPrincipal.Size(m)
}
func (m *MSPPrincipal) XXX_DiscardUnknown() {
	xxx_messageInfo_MSPPrincipal.DiscardUnknown(m)
}

var xxx_messageInfo_MSPPrincipal proto.InternalMessageInfo

func (m *MSPPrincipal) GetPrincipalClassification() MSPPrincipal_Classification {
	if m != nil {
		return m.PrincipalClassification
	}
	return MSPPrincipal_ROLE
}

func (m *MSPPrincipal) GetPrincipal() []byte {
	if m != nil {
		return m.Principal
	}
	return nil
}




type OrganizationUnit struct {
	
	
	MspIdentifier string `protobuf:"bytes,1,opt,name=msp_identifier,json=mspIdentifier,proto3" json:"msp_identifier,omitempty"`
	
	
	OrganizationalUnitIdentifier string `protobuf:"bytes,2,opt,name=organizational_unit_identifier,json=organizationalUnitIdentifier,proto3" json:"organizational_unit_identifier,omitempty"`
	
	
	CertifiersIdentifier []byte   `protobuf:"bytes,3,opt,name=certifiers_identifier,json=certifiersIdentifier,proto3" json:"certifiers_identifier,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrganizationUnit) Reset()         { *m = OrganizationUnit{} }
func (m *OrganizationUnit) String() string { return proto.CompactTextString(m) }
func (*OrganizationUnit) ProtoMessage()    {}
func (*OrganizationUnit) Descriptor() ([]byte, []int) {
	return fileDescriptor_82e08b7ead29bd48, []int{1}
}

func (m *OrganizationUnit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrganizationUnit.Unmarshal(m, b)
}
func (m *OrganizationUnit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrganizationUnit.Marshal(b, m, deterministic)
}
func (m *OrganizationUnit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrganizationUnit.Merge(m, src)
}
func (m *OrganizationUnit) XXX_Size() int {
	return xxx_messageInfo_OrganizationUnit.Size(m)
}
func (m *OrganizationUnit) XXX_DiscardUnknown() {
	xxx_messageInfo_OrganizationUnit.DiscardUnknown(m)
}

var xxx_messageInfo_OrganizationUnit proto.InternalMessageInfo

func (m *OrganizationUnit) GetMspIdentifier() string {
	if m != nil {
		return m.MspIdentifier
	}
	return ""
}

func (m *OrganizationUnit) GetOrganizationalUnitIdentifier() string {
	if m != nil {
		return m.OrganizationalUnitIdentifier
	}
	return ""
}

func (m *OrganizationUnit) GetCertifiersIdentifier() []byte {
	if m != nil {
		return m.CertifiersIdentifier
	}
	return nil
}




type MSPRole struct {
	
	
	MspIdentifier string `protobuf:"bytes,1,opt,name=msp_identifier,json=mspIdentifier,proto3" json:"msp_identifier,omitempty"`
	
	
	Role                 MSPRole_MSPRoleType `protobuf:"varint,2,opt,name=role,proto3,enum=common.MSPRole_MSPRoleType" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *MSPRole) Reset()         { *m = MSPRole{} }
func (m *MSPRole) String() string { return proto.CompactTextString(m) }
func (*MSPRole) ProtoMessage()    {}
func (*MSPRole) Descriptor() ([]byte, []int) {
	return fileDescriptor_82e08b7ead29bd48, []int{2}
}

func (m *MSPRole) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MSPRole.Unmarshal(m, b)
}
func (m *MSPRole) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MSPRole.Marshal(b, m, deterministic)
}
func (m *MSPRole) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MSPRole.Merge(m, src)
}
func (m *MSPRole) XXX_Size() int {
	return xxx_messageInfo_MSPRole.Size(m)
}
func (m *MSPRole) XXX_DiscardUnknown() {
	xxx_messageInfo_MSPRole.DiscardUnknown(m)
}

var xxx_messageInfo_MSPRole proto.InternalMessageInfo

func (m *MSPRole) GetMspIdentifier() string {
	if m != nil {
		return m.MspIdentifier
	}
	return ""
}

func (m *MSPRole) GetRole() MSPRole_MSPRoleType {
	if m != nil {
		return m.Role
	}
	return MSPRole_MEMBER
}


type MSPIdentityAnonymity struct {
	AnonymityType        MSPIdentityAnonymity_MSPIdentityAnonymityType `protobuf:"varint,1,opt,name=anonymity_type,json=anonymityType,proto3,enum=common.MSPIdentityAnonymity_MSPIdentityAnonymityType" json:"anonymity_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *MSPIdentityAnonymity) Reset()         { *m = MSPIdentityAnonymity{} }
func (m *MSPIdentityAnonymity) String() string { return proto.CompactTextString(m) }
func (*MSPIdentityAnonymity) ProtoMessage()    {}
func (*MSPIdentityAnonymity) Descriptor() ([]byte, []int) {
	return fileDescriptor_82e08b7ead29bd48, []int{3}
}

func (m *MSPIdentityAnonymity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MSPIdentityAnonymity.Unmarshal(m, b)
}
func (m *MSPIdentityAnonymity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MSPIdentityAnonymity.Marshal(b, m, deterministic)
}
func (m *MSPIdentityAnonymity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MSPIdentityAnonymity.Merge(m, src)
}
func (m *MSPIdentityAnonymity) XXX_Size() int {
	return xxx_messageInfo_MSPIdentityAnonymity.Size(m)
}
func (m *MSPIdentityAnonymity) XXX_DiscardUnknown() {
	xxx_messageInfo_MSPIdentityAnonymity.DiscardUnknown(m)
}

var xxx_messageInfo_MSPIdentityAnonymity proto.InternalMessageInfo

func (m *MSPIdentityAnonymity) GetAnonymityType() MSPIdentityAnonymity_MSPIdentityAnonymityType {
	if m != nil {
		return m.AnonymityType
	}
	return MSPIdentityAnonymity_NOMINAL
}




type CombinedPrincipal struct {
	
	Principals           []*MSPPrincipal `protobuf:"bytes,1,rep,name=principals,proto3" json:"principals,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CombinedPrincipal) Reset()         { *m = CombinedPrincipal{} }
func (m *CombinedPrincipal) String() string { return proto.CompactTextString(m) }
func (*CombinedPrincipal) ProtoMessage()    {}
func (*CombinedPrincipal) Descriptor() ([]byte, []int) {
	return fileDescriptor_82e08b7ead29bd48, []int{4}
}

func (m *CombinedPrincipal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CombinedPrincipal.Unmarshal(m, b)
}
func (m *CombinedPrincipal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CombinedPrincipal.Marshal(b, m, deterministic)
}
func (m *CombinedPrincipal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CombinedPrincipal.Merge(m, src)
}
func (m *CombinedPrincipal) XXX_Size() int {
	return xxx_messageInfo_CombinedPrincipal.Size(m)
}
func (m *CombinedPrincipal) XXX_DiscardUnknown() {
	xxx_messageInfo_CombinedPrincipal.DiscardUnknown(m)
}

var xxx_messageInfo_CombinedPrincipal proto.InternalMessageInfo

func (m *CombinedPrincipal) GetPrincipals() []*MSPPrincipal {
	if m != nil {
		return m.Principals
	}
	return nil
}

func init() {
	proto.RegisterEnum("common.MSPPrincipal_Classification", MSPPrincipal_Classification_name, MSPPrincipal_Classification_value)
	proto.RegisterEnum("common.MSPRole_MSPRoleType", MSPRole_MSPRoleType_name, MSPRole_MSPRoleType_value)
	proto.RegisterEnum("common.MSPIdentityAnonymity_MSPIdentityAnonymityType", MSPIdentityAnonymity_MSPIdentityAnonymityType_name, MSPIdentityAnonymity_MSPIdentityAnonymityType_value)
	proto.RegisterType((*MSPPrincipal)(nil), "common.MSPPrincipal")
	proto.RegisterType((*OrganizationUnit)(nil), "common.OrganizationUnit")
	proto.RegisterType((*MSPRole)(nil), "common.MSPRole")
	proto.RegisterType((*MSPIdentityAnonymity)(nil), "common.MSPIdentityAnonymity")
	proto.RegisterType((*CombinedPrincipal)(nil), "common.CombinedPrincipal")
}

func init() { proto.RegisterFile("msp/msp_principal.proto", fileDescriptor_82e08b7ead29bd48) }

var fileDescriptor_82e08b7ead29bd48 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xc1, 0x6e, 0xda, 0x40,
	0x10, 0x86, 0x63, 0xa0, 0x49, 0x98, 0x00, 0xda, 0xac, 0x88, 0x82, 0xd4, 0xa8, 0x42, 0x6e, 0x2b,
	0x71, 0x32, 0x12, 0x69, 0x7b, 0x37, 0x60, 0x45, 0x96, 0xf0, 0xda, 0x5a, 0xcc, 0x21, 0x51, 0x54,
	0x64, 0xcc, 0x42, 0x56, 0xb2, 0xbd, 0x96, 0xed, 0x1c, 0xdc, 0x47, 0xaa, 0x7a, 0xec, 0x53, 0xf5,
	0x29, 0x2a, 0xdb, 0x01, 0x96, 0x36, 0x95, 0x7a, 0xb2, 0x67, 0xe6, 0xfb, 0xc7, 0xbf, 0x77, 0x67,
	0xe0, 0x3a, 0x4c, 0xe3, 0x61, 0x98, 0xc6, 0xcb, 0x38, 0xe1, 0x91, 0xcf, 0x63, 0x2f, 0xd0, 0xe2,
	0x44, 0x64, 0x02, 0x9f, 0xfa, 0x22, 0x0c, 0x45, 0xa4, 0xfe, 0x52, 0xa0, 0x65, 0xcd, 0x1d, 0x67,
	0x57, 0xc6, 0x5f, 0xa1, 0xb7, 0x67, 0x97, 0x7e, 0xe0, 0xa5, 0x29, 0xdf, 0x70, 0xdf, 0xcb, 0xb8,
	0x88, 0x7a, 0x4a, 0x5f, 0x19, 0x74, 0x46, 0xef, 0xb5, 0x4a, 0xab, 0xc9, 0x3a, 0x6d, 0x72, 0x84,
	0xd2, 0xeb, 0x7d, 0x93, 0xe3, 0x02, 0xbe, 0x81, 0xe6, 0xbe, 0xd4, 0xab, 0xf5, 0x95, 0x41, 0x8b,
	0x1e, 0x12, 0xea, 0x23, 0x74, 0xfe, 0xe0, 0xcf, 0xa1, 0x41, 0xed, 0x99, 0x81, 0x4e, 0xf0, 0x15,
	0x5c, 0xda, 0xf4, 0x4e, 0x27, 0xe6, 0x83, 0xee, 0x9a, 0x36, 0x59, 0x2e, 0x88, 0xe9, 0x22, 0x05,
	0xb7, 0xe0, 0xdc, 0x9c, 0x1a, 0xc4, 0x35, 0xdd, 0x7b, 0x54, 0xc3, 0x6d, 0x68, 0xea, 0xc4, 0x26,
	0xf7, 0x56, 0x11, 0xd6, 0x8b, 0xe2, 0xc4, 0xb6, 0xc6, 0x26, 0x31, 0xa6, 0xa8, 0xa1, 0xfe, 0x54,
	0x00, 0xd9, 0xc9, 0xd6, 0x8b, 0xf8, 0xb7, 0xb2, 0xf9, 0x22, 0xe2, 0x19, 0xfe, 0x08, 0x9d, 0xe2,
	0x80, 0xf8, 0x9a, 0x45, 0x19, 0xdf, 0x70, 0x96, 0x94, 0xbf, 0xd9, 0xa4, 0xed, 0x30, 0x8d, 0xcd,
	0x7d, 0x12, 0x4f, 0xe1, 0x9d, 0x90, 0xa4, 0x5e, 0xb0, 0x7c, 0x8e, 0x78, 0x26, 0xcb, 0x6a, 0xa5,
	0xec, 0xe6, 0x98, 0x2a, 0x3e, 0x21, 0x75, 0xb9, 0x85, 0x2b, 0x9f, 0x25, 0x55, 0x90, 0xca, 0xe2,
	0x7a, 0x79, 0x12, 0xdd, 0x43, 0xf1, 0x20, 0x52, 0xbf, 0x2b, 0x70, 0x66, 0xcd, 0x1d, 0x2a, 0x02,
	0xf6, 0xbf, 0x6e, 0x87, 0xd0, 0x48, 0x44, 0xc0, 0x4a, 0x4f, 0x9d, 0xd1, 0x5b, 0xe9, 0xc6, 0x8a,
	0x2e, 0xbb, 0xa7, 0x9b, 0xc7, 0x8c, 0x96, 0xa0, 0x7a, 0x07, 0x17, 0x52, 0x12, 0x03, 0x9c, 0x5a,
	0x86, 0x35, 0x36, 0x28, 0x3a, 0xc1, 0x4d, 0x78, 0xa3, 0x4f, 0x2d, 0x93, 0x20, 0xa5, 0x48, 0x4f,
	0x66, 0xa6, 0x41, 0x5c, 0x54, 0x2b, 0x2e, 0xc6, 0x31, 0x0c, 0x8a, 0xea, 0xf8, 0x02, 0xce, 0x6c,
	0x3a, 0x35, 0xa8, 0x41, 0x51, 0x43, 0xfd, 0xa1, 0x40, 0xd7, 0x9a, 0x3b, 0x95, 0x97, 0x2c, 0xd7,
	0x23, 0x11, 0xe5, 0x21, 0xcf, 0x72, 0xfc, 0x08, 0x1d, 0x6f, 0x17, 0x2c, 0xb3, 0x3c, 0x66, 0x2f,
	0xe3, 0xf4, 0x59, 0x32, 0xf7, 0x97, 0xea, 0xd5, 0x64, 0x69, 0xbb, 0xed, 0xc9, 0xa1, 0xfa, 0x05,
	0x7a, 0xff, 0x42, 0x0b, 0x7f, 0xc4, 0xb6, 0x4c, 0xa2, 0xcf, 0xd0, 0xc9, 0x61, 0x40, 0xec, 0xc5,
	0x1c, 0x29, 0xaa, 0x09, 0x97, 0x13, 0x11, 0xae, 0x78, 0xc4, 0xd6, 0x87, 0x1d, 0xf8, 0x04, 0xb0,
	0x1f, 0xc9, 0xb4, 0xa7, 0xf4, 0xeb, 0x83, 0x8b, 0x51, 0xf7, 0xb5, 0xa9, 0xa7, 0x12, 0x37, 0x76,
	0xe0, 0x83, 0x48, 0xb6, 0xda, 0x53, 0x1e, 0xb3, 0x24, 0x60, 0xeb, 0x2d, 0x4b, 0xb4, 0x8d, 0xb7,
	0x4a, 0xb8, 0x5f, 0xad, 0x5c, 0xfa, 0xd2, 0xe0, 0x61, 0xb0, 0xe5, 0xd9, 0xd3, 0xf3, 0xaa, 0x08,
	0x87, 0x12, 0x3c, 0xac, 0xe0, 0x61, 0x05, 0x17, 0x4b, 0xbb, 0x3a, 0x2d, 0xdf, 0x6f, 0x7f, 0x07,
	0x00, 0x00, 0xff, 0xff, 0x30, 0x91, 0xbc, 0x0d, 0xc6, 0x03, 0x00, 0x00,
}
