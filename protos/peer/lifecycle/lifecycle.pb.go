


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
	return fileDescriptor_lifecycle_f98901bea638af10, []int{0}
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
	return fileDescriptor_lifecycle_f98901bea638af10, []int{1}
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



type QueryInstalledChaincodeArgs struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodeArgs) Reset()         { *m = QueryInstalledChaincodeArgs{} }
func (m *QueryInstalledChaincodeArgs) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodeArgs) ProtoMessage()    {}
func (*QueryInstalledChaincodeArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_f98901bea638af10, []int{2}
}
func (m *QueryInstalledChaincodeArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodeArgs.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodeArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodeArgs.Marshal(b, m, deterministic)
}
func (dst *QueryInstalledChaincodeArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodeArgs.Merge(dst, src)
}
func (m *QueryInstalledChaincodeArgs) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodeArgs.Size(m)
}
func (m *QueryInstalledChaincodeArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodeArgs.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodeArgs proto.InternalMessageInfo

func (m *QueryInstalledChaincodeArgs) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *QueryInstalledChaincodeArgs) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}



type QueryInstalledChaincodeResult struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodeResult) Reset()         { *m = QueryInstalledChaincodeResult{} }
func (m *QueryInstalledChaincodeResult) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodeResult) ProtoMessage()    {}
func (*QueryInstalledChaincodeResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_f98901bea638af10, []int{3}
}
func (m *QueryInstalledChaincodeResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodeResult.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodeResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodeResult.Marshal(b, m, deterministic)
}
func (dst *QueryInstalledChaincodeResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodeResult.Merge(dst, src)
}
func (m *QueryInstalledChaincodeResult) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodeResult.Size(m)
}
func (m *QueryInstalledChaincodeResult) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodeResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodeResult proto.InternalMessageInfo

func (m *QueryInstalledChaincodeResult) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func init() {
	proto.RegisterType((*InstallChaincodeArgs)(nil), "lifecycle.InstallChaincodeArgs")
	proto.RegisterType((*InstallChaincodeResult)(nil), "lifecycle.InstallChaincodeResult")
	proto.RegisterType((*QueryInstalledChaincodeArgs)(nil), "lifecycle.QueryInstalledChaincodeArgs")
	proto.RegisterType((*QueryInstalledChaincodeResult)(nil), "lifecycle.QueryInstalledChaincodeResult")
}

func init() {
	proto.RegisterFile("peer/lifecycle/lifecycle.proto", fileDescriptor_lifecycle_f98901bea638af10)
}

var fileDescriptor_lifecycle_f98901bea638af10 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x15, 0x40, 0xa0, 0x9e, 0x3a, 0x59, 0x08, 0x82, 0x10, 0xa8, 0xca, 0xd4, 0xa1, 0xb2,
	0x87, 0x6e, 0x6c, 0xc0, 0x84, 0x58, 0x20, 0x23, 0x4b, 0xe5, 0x38, 0x57, 0xdb, 0xc2, 0x8d, 0xa3,
	0x73, 0x82, 0x94, 0x8d, 0x8f, 0x8e, 0x6a, 0xb7, 0xe1, 0x8f, 0xd4, 0x89, 0xed, 0xf9, 0xd9, 0xef,
	0xfc, 0xd3, 0x3b, 0xb8, 0x6d, 0x11, 0x49, 0x38, 0xbb, 0x46, 0x35, 0x28, 0x87, 0xdf, 0x8a, 0xb7,
	0xe4, 0x3b, 0xcf, 0x26, 0xa3, 0x51, 0x7c, 0x66, 0x70, 0xfe, 0xd4, 0x84, 0x4e, 0x3a, 0xf7, 0x68,
	0xa4, 0x6d, 0x94, 0xaf, 0xf1, 0x9e, 0x74, 0x60, 0x0c, 0x4e, 0x1a, 0xb9, 0xc1, 0x3c, 0x9b, 0x65,
	0xf3, 0x49, 0x19, 0x35, 0xcb, 0xe1, 0xec, 0x03, 0x29, 0x58, 0xdf, 0xe4, 0x47, 0xd1, 0xde, 0x1f,
	0xd9, 0x1d, 0x5c, 0xa9, 0x7d, 0x7c, 0x65, 0xd3, 0xbc, 0x55, 0x2b, 0xd5, 0xbb, 0xd4, 0x98, 0x1f,
	0xcf, 0xb2, 0xf9, 0xb4, 0xbc, 0x1c, 0x1f, 0xec, 0xfe, 0x7b, 0x49, 0xd7, 0xc5, 0x02, 0x2e, 0xfe,
	0x12, 0x94, 0x18, 0x7a, 0xd7, 0x6d, 0x19, 0x8c, 0x0c, 0x26, 0x32, 0x4c, 0xcb, 0xa8, 0x8b, 0x67,
	0xb8, 0x7e, 0xed, 0x91, 0x86, 0x5d, 0x04, 0xeb, 0x7f, 0x60, 0x17, 0x4b, 0xb8, 0x39, 0x30, 0xec,
	0x30, 0xc1, 0x83, 0x82, 0x85, 0x27, 0xcd, 0xcd, 0xd0, 0x22, 0x39, 0xac, 0x35, 0x12, 0x5f, 0xcb,
	0x8a, 0xac, 0x4a, 0xed, 0x06, 0xbe, 0x6d, 0x9f, 0x8f, 0x15, 0xbf, 0x2d, 0xb5, 0xed, 0x4c, 0x5f,
	0x71, 0xe5, 0x37, 0xe2, 0x47, 0x48, 0xa4, 0x90, 0x48, 0x21, 0xf1, 0x7b, 0x65, 0xd5, 0x69, 0xb4,
	0x97, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd3, 0x6e, 0xec, 0xb6, 0xcb, 0x01, 0x00, 0x00,
}
