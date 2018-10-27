


package lifecycle 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



type InstallChaincodeArgs struct {
	Name                    string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version                 string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	ChaincodeInstallPackage []byte   `protobuf:"bytes,3,opt,name=chaincode_install_package,json=chaincodeInstallPackage,proto3" json:"chaincode_install_package,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *InstallChaincodeArgs) Reset()         { *m = InstallChaincodeArgs{} }
func (m *InstallChaincodeArgs) String() string { return proto.CompactTextString(m) }
func (*InstallChaincodeArgs) ProtoMessage()    {}
func (*InstallChaincodeArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_5fcb277bc76852f1, []int{0}
}
func (m *InstallChaincodeArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallChaincodeArgs.Unmarshal(m, b)
}
func (m *InstallChaincodeArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallChaincodeArgs.Marshal(b, m, deterministic)
}
func (dst *InstallChaincodeArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallChaincodeArgs.Merge(dst, src)
}
func (m *InstallChaincodeArgs) XXX_Size() int {
	return xxx_messageInfo_InstallChaincodeArgs.Size(m)
}
func (m *InstallChaincodeArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallChaincodeArgs.DiscardUnknown(m)
}

var xxx_messageInfo_InstallChaincodeArgs proto.InternalMessageInfo

func (m *InstallChaincodeArgs) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *InstallChaincodeArgs) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *InstallChaincodeArgs) GetChaincodeInstallPackage() []byte {
	if m != nil {
		return m.ChaincodeInstallPackage
	}
	return nil
}



type InstallChaincodeResult struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallChaincodeResult) Reset()         { *m = InstallChaincodeResult{} }
func (m *InstallChaincodeResult) String() string { return proto.CompactTextString(m) }
func (*InstallChaincodeResult) ProtoMessage()    {}
func (*InstallChaincodeResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_5fcb277bc76852f1, []int{1}
}
func (m *InstallChaincodeResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallChaincodeResult.Unmarshal(m, b)
}
func (m *InstallChaincodeResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallChaincodeResult.Marshal(b, m, deterministic)
}
func (dst *InstallChaincodeResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallChaincodeResult.Merge(dst, src)
}
func (m *InstallChaincodeResult) XXX_Size() int {
	return xxx_messageInfo_InstallChaincodeResult.Size(m)
}
func (m *InstallChaincodeResult) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallChaincodeResult.DiscardUnknown(m)
}

var xxx_messageInfo_InstallChaincodeResult proto.InternalMessageInfo

func (m *InstallChaincodeResult) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func init() {
	proto.RegisterType((*InstallChaincodeArgs)(nil), "lifecycle.InstallChaincodeArgs")
	proto.RegisterType((*InstallChaincodeResult)(nil), "lifecycle.InstallChaincodeResult")
}

func init() {
	proto.RegisterFile("peer/lifecycle/lifecycle.proto", fileDescriptor_lifecycle_5fcb277bc76852f1)
}

var fileDescriptor_lifecycle_5fcb277bc76852f1 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x3f, 0x4f, 0x43, 0x21,
	0x14, 0xc5, 0xf3, 0xd4, 0x68, 0x4a, 0x3a, 0x11, 0xa3, 0xcf, 0xc5, 0x34, 0x9d, 0x3a, 0x34, 0x30,
	0x74, 0x73, 0x53, 0x27, 0x37, 0xc3, 0xe8, 0xd2, 0xf0, 0x6e, 0x6f, 0x81, 0x48, 0x81, 0x5c, 0xa8,
	0x49, 0x37, 0x3f, 0xba, 0x29, 0xd8, 0xe7, 0x9f, 0xed, 0x70, 0xe0, 0xc7, 0x3d, 0xe7, 0xb2, 0xfb,
	0x84, 0x48, 0xd2, 0xbb, 0x2d, 0xc2, 0x01, 0x3c, 0xfe, 0x28, 0x91, 0x28, 0x96, 0xc8, 0x27, 0xa3,
	0x31, 0xff, 0xec, 0xd8, 0xf5, 0x4b, 0xc8, 0x45, 0x7b, 0xff, 0x6c, 0xb5, 0x0b, 0x10, 0x37, 0xf8,
	0x48, 0x26, 0x73, 0xce, 0x2e, 0x82, 0xde, 0x61, 0xdf, 0xcd, 0xba, 0xc5, 0x44, 0x55, 0xcd, 0x7b,
	0x76, 0xf5, 0x81, 0x94, 0x5d, 0x0c, 0xfd, 0x59, 0xb5, 0x4f, 0x47, 0xfe, 0xc0, 0xee, 0xe0, 0x84,
	0xaf, 0x5d, 0xfb, 0x6f, 0x9d, 0x34, 0xbc, 0x6b, 0x83, 0xfd, 0xf9, 0xac, 0x5b, 0x4c, 0xd5, 0xed,
	0xf8, 0xe0, 0x7b, 0xde, 0x6b, 0xbb, 0x9e, 0x2f, 0xd9, 0xcd, 0xff, 0x04, 0x0a, 0xf3, 0xde, 0x97,
	0x63, 0x06, 0xab, 0xb3, 0xad, 0x19, 0xa6, 0xaa, 0xea, 0x27, 0x60, 0xcb, 0x48, 0x46, 0xd8, 0x43,
	0x42, 0xf2, 0xb8, 0x31, 0x48, 0x62, 0xab, 0x07, 0x72, 0xd0, 0xba, 0x65, 0x71, 0xec, 0x2e, 0xc6,
	0x82, 0x6f, 0x2b, 0xe3, 0x8a, 0xdd, 0x0f, 0x02, 0xe2, 0x4e, 0xfe, 0x82, 0x64, 0x83, 0x64, 0x83,
	0xe4, 0xdf, 0x85, 0x0d, 0x97, 0xd5, 0x5e, 0x7d, 0x05, 0x00, 0x00, 0xff, 0xff, 0x59, 0xab, 0xe4,
	0xee, 0x49, 0x01, 0x00, 0x00,
}
