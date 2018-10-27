


package msp 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

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
	return fileDescriptor_msp_principal_9016cf1a8a7156cd, []int{0, 0}
}

type MSPRole_MSPRoleType int32

const (
	MSPRole_MEMBER MSPRole_MSPRoleType = 0
	MSPRole_ADMIN  MSPRole_MSPRoleType = 1
	MSPRole_CLIENT MSPRole_MSPRoleType = 2
	MSPRole_PEER   MSPRole_MSPRoleType = 3
)

var MSPRole_MSPRoleType_name = map[int32]string{
	0: "MEMBER",
	1: "ADMIN",
	2: "CLIENT",
	3: "PEER",
}
var MSPRole_MSPRoleType_value = map[string]int32{
	"MEMBER": 0,
	"ADMIN":  1,
	"CLIENT": 2,
	"PEER":   3,
}

func (x MSPRole_MSPRoleType) String() string {
	return proto.EnumName(MSPRole_MSPRoleType_name, int32(x))
}
func (MSPRole_MSPRoleType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_msp_principal_9016cf1a8a7156cd, []int{2, 0}
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
	return fileDescriptor_msp_principal_9016cf1a8a7156cd, []int{3, 0}
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
	return fileDescriptor_msp_principal_9016cf1a8a7156cd, []int{0}
}
func (m *MSPPrincipal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MSPPrincipal.Unmarshal(m, b)
}
func (m *MSPPrincipal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MSPPrincipal.Marshal(b, m, deterministic)
}
func (dst *MSPPrincipal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MSPPrincipal.Merge(dst, src)
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
	return fileDescriptor_msp_principal_9016cf1a8a7156cd, []int{1}
}
func (m *OrganizationUnit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrganizationUnit.Unmarshal(m, b)
}
func (m *OrganizationUnit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrganizationUnit.Marshal(b, m, deterministic)
}
func (dst *OrganizationUnit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrganizationUnit.Merge(dst, src)
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
	return fileDescriptor_msp_principal_9016cf1a8a7156cd, []int{2}
}
func (m *MSPRole) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MSPRole.Unmarshal(m, b)
}
func (m *MSPRole) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MSPRole.Marshal(b, m, deterministic)
}
func (dst *MSPRole) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MSPRole.Merge(dst, src)
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
	return fileDescriptor_msp_principal_9016cf1a8a7156cd, []int{3}
}
func (m *MSPIdentityAnonymity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MSPIdentityAnonymity.Unmarshal(m, b)
}
func (m *MSPIdentityAnonymity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MSPIdentityAnonymity.Marshal(b, m, deterministic)
}
func (dst *MSPIdentityAnonymity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MSPIdentityAnonymity.Merge(dst, src)
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
	return fileDescriptor_msp_principal_9016cf1a8a7156cd, []int{4}
}
func (m *CombinedPrincipal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CombinedPrincipal.Unmarshal(m, b)
}
func (m *CombinedPrincipal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CombinedPrincipal.Marshal(b, m, deterministic)
}
func (dst *CombinedPrincipal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CombinedPrincipal.Merge(dst, src)
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
	proto.RegisterType((*MSPPrincipal)(nil), "common.MSPPrincipal")
	proto.RegisterType((*OrganizationUnit)(nil), "common.OrganizationUnit")
	proto.RegisterType((*MSPRole)(nil), "common.MSPRole")
	proto.RegisterType((*MSPIdentityAnonymity)(nil), "common.MSPIdentityAnonymity")
	proto.RegisterType((*CombinedPrincipal)(nil), "common.CombinedPrincipal")
	proto.RegisterEnum("common.MSPPrincipal_Classification", MSPPrincipal_Classification_name, MSPPrincipal_Classification_value)
	proto.RegisterEnum("common.MSPRole_MSPRoleType", MSPRole_MSPRoleType_name, MSPRole_MSPRoleType_value)
	proto.RegisterEnum("common.MSPIdentityAnonymity_MSPIdentityAnonymityType", MSPIdentityAnonymity_MSPIdentityAnonymityType_name, MSPIdentityAnonymity_MSPIdentityAnonymityType_value)
}

func init() {
	proto.RegisterFile("msp/msp_principal.proto", fileDescriptor_msp_principal_9016cf1a8a7156cd)
}

var fileDescriptor_msp_principal_9016cf1a8a7156cd = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xdf, 0x6a, 0xdb, 0x30,
	0x14, 0xc6, 0xeb, 0xa4, 0x6b, 0x9b, 0x93, 0x3f, 0xa8, 0x22, 0xa5, 0x81, 0x95, 0x11, 0xbc, 0x0d,
	0x72, 0xe5, 0x40, 0xba, 0xed, 0x62, 0x77, 0x4e, 0x62, 0x86, 0x20, 0x96, 0x8d, 0xe3, 0x5c, 0xb4,
	0x94, 0x05, 0xc7, 0x51, 0x52, 0x81, 0x6d, 0x19, 0xdb, 0xbd, 0xf0, 0xde, 0x65, 0x6f, 0xb0, 0xcb,
	0x3d, 0xd5, 0x9e, 0x62, 0xd8, 0x6e, 0x12, 0x65, 0xeb, 0x60, 0x57, 0xf6, 0x39, 0xe7, 0xf7, 0x1d,
	0x1d, 0x49, 0x9f, 0xe0, 0x3a, 0x4c, 0xe3, 0x61, 0x98, 0xc6, 0xcb, 0x38, 0xe1, 0x91, 0xcf, 0x63,
	0x2f, 0xd0, 0xe2, 0x44, 0x64, 0x02, 0x9f, 0xf9, 0x22, 0x0c, 0x45, 0xa4, 0xfe, 0x52, 0xa0, 0x65,
	0xce, 0x6d, 0x7b, 0x57, 0xc6, 0x5f, 0xa1, 0xb7, 0x67, 0x97, 0x7e, 0xe0, 0xa5, 0x29, 0xdf, 0x70,
	0xdf, 0xcb, 0xb8, 0x88, 0x7a, 0x4a, 0x5f, 0x19, 0x74, 0x46, 0x6f, 0xb5, 0x4a, 0xab, 0xc9, 0x3a,
	0x6d, 0x72, 0x84, 0x3a, 0xd7, 0xfb, 0x26, 0xc7, 0x05, 0x7c, 0x03, 0x8d, 0x7d, 0xa9, 0x57, 0xeb,
	0x2b, 0x83, 0x96, 0x73, 0x48, 0xa8, 0x0f, 0xd0, 0xf9, 0x83, 0xbf, 0x80, 0x53, 0xc7, 0x9a, 0x19,
	0xe8, 0x04, 0x5f, 0xc1, 0xa5, 0xe5, 0x7c, 0xd1, 0x29, 0xb9, 0xd7, 0x5d, 0x62, 0xd1, 0xe5, 0x82,
	0x12, 0x17, 0x29, 0xb8, 0x05, 0x17, 0x64, 0x6a, 0x50, 0x97, 0xb8, 0x77, 0xa8, 0x86, 0xdb, 0xd0,
	0xd0, 0xa9, 0x45, 0xef, 0xcc, 0x22, 0xac, 0x17, 0xc5, 0x89, 0x65, 0x8e, 0x09, 0x35, 0xa6, 0xe8,
	0x54, 0xfd, 0xa9, 0x00, 0xb2, 0x92, 0xad, 0x17, 0xf1, 0x6f, 0x65, 0xf3, 0x45, 0xc4, 0x33, 0xfc,
	0x1e, 0x3a, 0xc5, 0x01, 0xf1, 0x35, 0x8b, 0x32, 0xbe, 0xe1, 0x2c, 0x29, 0xb7, 0xd9, 0x70, 0xda,
	0x61, 0x1a, 0x93, 0x7d, 0x12, 0x4f, 0xe1, 0x8d, 0x90, 0xa4, 0x5e, 0xb0, 0x7c, 0x8a, 0x78, 0x26,
	0xcb, 0x6a, 0xa5, 0xec, 0xe6, 0x98, 0x2a, 0x96, 0x90, 0xba, 0xdc, 0xc2, 0x95, 0xcf, 0x92, 0x2a,
	0x48, 0x65, 0x71, 0xbd, 0x3c, 0x89, 0xee, 0xa1, 0x78, 0x10, 0xa9, 0xdf, 0x15, 0x38, 0x37, 0xe7,
	0xb6, 0x23, 0x02, 0xf6, 0xbf, 0xd3, 0x0e, 0xe1, 0x34, 0x11, 0x01, 0x2b, 0x67, 0xea, 0x8c, 0x5e,
	0x4b, 0x37, 0x56, 0x74, 0xd9, 0x7d, 0xdd, 0x3c, 0x66, 0x4e, 0x09, 0xaa, 0x9f, 0xa1, 0x29, 0x25,
	0x31, 0xc0, 0x99, 0x69, 0x98, 0x63, 0xc3, 0x41, 0x27, 0xb8, 0x01, 0xaf, 0xf4, 0xa9, 0x49, 0x28,
	0x52, 0x8a, 0xf4, 0x64, 0x46, 0x0c, 0xea, 0xa2, 0x5a, 0x71, 0x31, 0xb6, 0x61, 0x38, 0xa8, 0xae,
	0xfe, 0x50, 0xa0, 0x6b, 0xce, 0xed, 0x6a, 0xf9, 0x2c, 0xd7, 0x23, 0x11, 0xe5, 0x21, 0xcf, 0x72,
	0xfc, 0x00, 0x1d, 0x6f, 0x17, 0x2c, 0xb3, 0x3c, 0x66, 0xcf, 0x0e, 0xfa, 0x28, 0xcd, 0xf3, 0x97,
	0xea, 0xc5, 0x64, 0x39, 0x69, 0xdb, 0x93, 0x43, 0xf5, 0x13, 0xf4, 0xfe, 0x85, 0xe2, 0x26, 0x9c,
	0x53, 0xcb, 0x24, 0x54, 0x9f, 0xa1, 0x93, 0x83, 0x27, 0xac, 0xc5, 0x1c, 0x29, 0x2a, 0x81, 0xcb,
	0x89, 0x08, 0x57, 0x3c, 0x62, 0xeb, 0x83, 0xed, 0x3f, 0x00, 0xec, 0x5d, 0x98, 0xf6, 0x94, 0x7e,
	0x7d, 0xd0, 0x1c, 0x75, 0x5f, 0x32, 0xba, 0x23, 0x71, 0x63, 0x1b, 0xde, 0x89, 0x64, 0xab, 0x3d,
	0xe6, 0x31, 0x4b, 0x02, 0xb6, 0xde, 0xb2, 0x44, 0xdb, 0x78, 0xab, 0x84, 0xfb, 0xd5, 0x2b, 0x4b,
	0x9f, 0x1b, 0xdc, 0x0f, 0xb6, 0x3c, 0x7b, 0x7c, 0x5a, 0x15, 0xe1, 0x50, 0x82, 0x87, 0x15, 0x3c,
	0xac, 0xe0, 0xe2, 0x9d, 0xae, 0xce, 0xca, 0xff, 0xdb, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x40,
	0x36, 0xd2, 0xf9, 0xb9, 0x03, 0x00, 0x00,
}
