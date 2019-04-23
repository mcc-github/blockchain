


package lifecycle 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



type ChaincodeEndorsementInfo struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	InitRequired         bool     `protobuf:"varint,2,opt,name=init_required,json=initRequired,proto3" json:"init_required,omitempty"`
	EndorsementPlugin    string   `protobuf:"bytes,3,opt,name=endorsement_plugin,json=endorsementPlugin,proto3" json:"endorsement_plugin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeEndorsementInfo) Reset()         { *m = ChaincodeEndorsementInfo{} }
func (m *ChaincodeEndorsementInfo) String() string { return proto.CompactTextString(m) }
func (*ChaincodeEndorsementInfo) ProtoMessage()    {}
func (*ChaincodeEndorsementInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_definition_bcb7540be0cd1d31, []int{0}
}
func (m *ChaincodeEndorsementInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeEndorsementInfo.Unmarshal(m, b)
}
func (m *ChaincodeEndorsementInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeEndorsementInfo.Marshal(b, m, deterministic)
}
func (dst *ChaincodeEndorsementInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeEndorsementInfo.Merge(dst, src)
}
func (m *ChaincodeEndorsementInfo) XXX_Size() int {
	return xxx_messageInfo_ChaincodeEndorsementInfo.Size(m)
}
func (m *ChaincodeEndorsementInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeEndorsementInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeEndorsementInfo proto.InternalMessageInfo

func (m *ChaincodeEndorsementInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ChaincodeEndorsementInfo) GetInitRequired() bool {
	if m != nil {
		return m.InitRequired
	}
	return false
}

func (m *ChaincodeEndorsementInfo) GetEndorsementPlugin() string {
	if m != nil {
		return m.EndorsementPlugin
	}
	return ""
}



type ChaincodeValidationInfo struct {
	ValidationPlugin     string   `protobuf:"bytes,1,opt,name=validation_plugin,json=validationPlugin,proto3" json:"validation_plugin,omitempty"`
	ValidationParameter  []byte   `protobuf:"bytes,2,opt,name=validation_parameter,json=validationParameter,proto3" json:"validation_parameter,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeValidationInfo) Reset()         { *m = ChaincodeValidationInfo{} }
func (m *ChaincodeValidationInfo) String() string { return proto.CompactTextString(m) }
func (*ChaincodeValidationInfo) ProtoMessage()    {}
func (*ChaincodeValidationInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_definition_bcb7540be0cd1d31, []int{1}
}
func (m *ChaincodeValidationInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeValidationInfo.Unmarshal(m, b)
}
func (m *ChaincodeValidationInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeValidationInfo.Marshal(b, m, deterministic)
}
func (dst *ChaincodeValidationInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeValidationInfo.Merge(dst, src)
}
func (m *ChaincodeValidationInfo) XXX_Size() int {
	return xxx_messageInfo_ChaincodeValidationInfo.Size(m)
}
func (m *ChaincodeValidationInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeValidationInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeValidationInfo proto.InternalMessageInfo

func (m *ChaincodeValidationInfo) GetValidationPlugin() string {
	if m != nil {
		return m.ValidationPlugin
	}
	return ""
}

func (m *ChaincodeValidationInfo) GetValidationParameter() []byte {
	if m != nil {
		return m.ValidationParameter
	}
	return nil
}

func init() {
	proto.RegisterType((*ChaincodeEndorsementInfo)(nil), "lifecycle.ChaincodeEndorsementInfo")
	proto.RegisterType((*ChaincodeValidationInfo)(nil), "lifecycle.ChaincodeValidationInfo")
}

func init() {
	proto.RegisterFile("peer/lifecycle/chaincode_definition.proto", fileDescriptor_chaincode_definition_bcb7540be0cd1d31)
}

var fileDescriptor_chaincode_definition_bcb7540be0cd1d31 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0x41, 0x4b, 0xc4, 0x30,
	0x14, 0x84, 0xa9, 0x82, 0xba, 0x61, 0x05, 0x37, 0x0a, 0xf6, 0xb8, 0xac, 0x97, 0x15, 0x35, 0x41,
	0xf6, 0x1f, 0x28, 0x1e, 0xbc, 0x49, 0x0f, 0x1e, 0xbc, 0x2c, 0x69, 0xf2, 0xda, 0x3e, 0x68, 0x93,
	0xfa, 0xda, 0x2e, 0xf4, 0x1f, 0xf8, 0xb3, 0xa5, 0x89, 0xed, 0xd6, 0x63, 0x66, 0xe6, 0x0d, 0x1f,
	0x13, 0x76, 0x5f, 0x03, 0x90, 0x2c, 0x31, 0x03, 0xdd, 0xeb, 0x12, 0xa4, 0x2e, 0x14, 0x5a, 0xed,
	0x0c, 0xec, 0x0d, 0x64, 0x68, 0xb1, 0x45, 0x67, 0x45, 0x4d, 0xae, 0x75, 0x7c, 0x31, 0xa5, 0x36,
	0x3f, 0x11, 0x8b, 0x5f, 0xc7, 0xe4, 0x9b, 0x35, 0x8e, 0x1a, 0xa8, 0xc0, 0xb6, 0xef, 0x36, 0x73,
	0x3c, 0x66, 0xe7, 0x07, 0xa0, 0x06, 0x9d, 0x8d, 0xa3, 0x75, 0xb4, 0x5d, 0x24, 0xe3, 0x93, 0xdf,
	0xb1, 0xcb, 0xa1, 0x72, 0x4f, 0xf0, 0xdd, 0x21, 0x81, 0x89, 0x4f, 0xd6, 0xd1, 0xf6, 0x22, 0x59,
	0x0e, 0x62, 0xf2, 0xa7, 0xf1, 0x27, 0xc6, 0xe1, 0xd8, 0xb8, 0xaf, 0xcb, 0x2e, 0x47, 0x1b, 0x9f,
	0xfa, 0xa6, 0xd5, 0xcc, 0xf9, 0xf0, 0xc6, 0xa6, 0x67, 0xb7, 0x13, 0xc9, 0xa7, 0x2a, 0xd1, 0xa8,
	0x01, 0xd9, 0x83, 0x3c, 0xb0, 0xd5, 0x61, 0x52, 0xc6, 0xa2, 0x80, 0x74, 0x75, 0x34, 0x42, 0x0f,
	0x7f, 0x66, 0x37, 0xf3, 0xb0, 0x22, 0x55, 0x41, 0x0b, 0xe4, 0x11, 0x97, 0xc9, 0xf5, 0x2c, 0x3f,
	0x5a, 0x2f, 0x9a, 0x3d, 0x3a, 0xca, 0x45, 0xd1, 0xd7, 0x40, 0x25, 0x98, 0x1c, 0x48, 0x64, 0x2a,
	0x25, 0xd4, 0x61, 0xb0, 0x46, 0x0c, 0xdb, 0x8a, 0x69, 0xb5, 0xaf, 0x5d, 0x8e, 0x6d, 0xd1, 0xa5,
	0x42, 0xbb, 0x4a, 0xce, 0x8e, 0x64, 0x38, 0x92, 0xe1, 0x48, 0xfe, 0xff, 0x90, 0xf4, 0xcc, 0xcb,
	0xbb, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x93, 0xc8, 0xb4, 0xe8, 0xa9, 0x01, 0x00, 0x00,
}
