


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
	return fileDescriptor_lifecycle_e9f4a4fdde6b4b81, []int{0}
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
	return fileDescriptor_lifecycle_e9f4a4fdde6b4b81, []int{1}
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
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodeArgs) Reset()         { *m = QueryInstalledChaincodeArgs{} }
func (m *QueryInstalledChaincodeArgs) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodeArgs) ProtoMessage()    {}
func (*QueryInstalledChaincodeArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_e9f4a4fdde6b4b81, []int{2}
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
	return fileDescriptor_lifecycle_e9f4a4fdde6b4b81, []int{3}
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




type QueryInstalledChaincodesArgs struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodesArgs) Reset()         { *m = QueryInstalledChaincodesArgs{} }
func (m *QueryInstalledChaincodesArgs) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodesArgs) ProtoMessage()    {}
func (*QueryInstalledChaincodesArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_e9f4a4fdde6b4b81, []int{4}
}
func (m *QueryInstalledChaincodesArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodesArgs.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodesArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodesArgs.Marshal(b, m, deterministic)
}
func (dst *QueryInstalledChaincodesArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodesArgs.Merge(dst, src)
}
func (m *QueryInstalledChaincodesArgs) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodesArgs.Size(m)
}
func (m *QueryInstalledChaincodesArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodesArgs.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodesArgs proto.InternalMessageInfo




type QueryInstalledChaincodesResult struct {
	InstalledChaincodes  []*QueryInstalledChaincodesResult_InstalledChaincode `protobuf:"bytes,1,rep,name=installed_chaincodes,json=installedChaincodes" json:"installed_chaincodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                             `json:"-"`
	XXX_unrecognized     []byte                                               `json:"-"`
	XXX_sizecache        int32                                                `json:"-"`
}

func (m *QueryInstalledChaincodesResult) Reset()         { *m = QueryInstalledChaincodesResult{} }
func (m *QueryInstalledChaincodesResult) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodesResult) ProtoMessage()    {}
func (*QueryInstalledChaincodesResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_e9f4a4fdde6b4b81, []int{5}
}
func (m *QueryInstalledChaincodesResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodesResult.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodesResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodesResult.Marshal(b, m, deterministic)
}
func (dst *QueryInstalledChaincodesResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodesResult.Merge(dst, src)
}
func (m *QueryInstalledChaincodesResult) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodesResult.Size(m)
}
func (m *QueryInstalledChaincodesResult) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodesResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodesResult proto.InternalMessageInfo

func (m *QueryInstalledChaincodesResult) GetInstalledChaincodes() []*QueryInstalledChaincodesResult_InstalledChaincode {
	if m != nil {
		return m.InstalledChaincodes
	}
	return nil
}

type QueryInstalledChaincodesResult_InstalledChaincode struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
	Hash                 []byte   `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodesResult_InstalledChaincode) Reset() {
	*m = QueryInstalledChaincodesResult_InstalledChaincode{}
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) String() string {
	return proto.CompactTextString(m)
}
func (*QueryInstalledChaincodesResult_InstalledChaincode) ProtoMessage() {}
func (*QueryInstalledChaincodesResult_InstalledChaincode) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_e9f4a4fdde6b4b81, []int{5, 0}
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.Marshal(b, m, deterministic)
}
func (dst *QueryInstalledChaincodesResult_InstalledChaincode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.Merge(dst, src)
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.Size(m)
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode proto.InternalMessageInfo

func (m *QueryInstalledChaincodesResult_InstalledChaincode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *QueryInstalledChaincodesResult_InstalledChaincode) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *QueryInstalledChaincodesResult_InstalledChaincode) GetHash() []byte {
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
	proto.RegisterType((*QueryInstalledChaincodesArgs)(nil), "lifecycle.QueryInstalledChaincodesArgs")
	proto.RegisterType((*QueryInstalledChaincodesResult)(nil), "lifecycle.QueryInstalledChaincodesResult")
	proto.RegisterType((*QueryInstalledChaincodesResult_InstalledChaincode)(nil), "lifecycle.QueryInstalledChaincodesResult.InstalledChaincode")
}

func init() {
	proto.RegisterFile("peer/lifecycle/lifecycle.proto", fileDescriptor_lifecycle_e9f4a4fdde6b4b81)
}

var fileDescriptor_lifecycle_e9f4a4fdde6b4b81 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x65, 0x8a, 0x40, 0x3d, 0x3a, 0x99, 0x0a, 0xc2, 0xbf, 0xaa, 0xca, 0xd4, 0xa1, 0x72,
	0x24, 0xba, 0x21, 0x16, 0x60, 0x42, 0x2c, 0x90, 0x81, 0x81, 0xa5, 0x72, 0xdd, 0x6b, 0x62, 0xe1,
	0xc6, 0x91, 0x9d, 0x20, 0x65, 0xe3, 0xeb, 0xf2, 0x2d, 0x50, 0xfe, 0x03, 0x25, 0x95, 0x10, 0xdb,
	0xe5, 0xfc, 0xde, 0xcb, 0xcf, 0xe7, 0x83, 0x51, 0x8c, 0x68, 0x3c, 0x25, 0x57, 0x28, 0x32, 0xa1,
	0xb0, 0xad, 0x58, 0x6c, 0x74, 0xa2, 0x69, 0xbf, 0x69, 0xb8, 0xef, 0x04, 0x86, 0xf7, 0x91, 0x4d,
	0xb8, 0x52, 0x77, 0x21, 0x97, 0x91, 0xd0, 0x4b, 0xbc, 0x31, 0x81, 0xa5, 0x14, 0x76, 0x23, 0xbe,
	0x46, 0x87, 0x8c, 0xc9, 0xa4, 0xef, 0x17, 0x35, 0x75, 0x60, 0xff, 0x0d, 0x8d, 0x95, 0x3a, 0x72,
	0x76, 0x8a, 0x76, 0xfd, 0x49, 0xaf, 0xe0, 0x44, 0xd4, 0xf6, 0xb9, 0x2c, 0xf3, 0xe6, 0x31, 0x17,
	0xaf, 0x3c, 0x40, 0xa7, 0x37, 0x26, 0x93, 0x81, 0x7f, 0xdc, 0x08, 0xaa, 0xff, 0x3d, 0x96, 0xc7,
	0xee, 0x14, 0x8e, 0x7e, 0x12, 0xf8, 0x68, 0x53, 0x95, 0xe4, 0x0c, 0x21, 0xb7, 0x61, 0xc1, 0x30,
	0xf0, 0x8b, 0xda, 0x7d, 0x80, 0xb3, 0xa7, 0x14, 0x4d, 0x56, 0x59, 0x70, 0xf9, 0x0f, 0x6c, 0x77,
	0x06, 0x17, 0x1d, 0x61, 0x5b, 0x08, 0x46, 0x70, 0xde, 0x61, 0xb2, 0x39, 0x82, 0xfb, 0x41, 0x60,
	0xd4, 0x25, 0xa8, 0x62, 0x35, 0x0c, 0x65, 0x7d, 0x38, 0x6f, 0xe6, 0x62, 0x1d, 0x32, 0xee, 0x4d,
	0x0e, 0x2e, 0xaf, 0x59, 0xfb, 0x60, 0xdb, 0x83, 0xd8, 0x2f, 0xe0, 0x87, 0x72, 0x53, 0x7d, 0xfa,
	0x0c, 0x74, 0x53, 0xfa, 0xc7, 0x37, 0xae, 0x67, 0xd1, 0x6b, 0x67, 0x71, 0x2b, 0x60, 0xaa, 0x4d,
	0xc0, 0xc2, 0x2c, 0x46, 0xa3, 0x70, 0x19, 0xa0, 0x61, 0x2b, 0xbe, 0x30, 0x52, 0x94, 0x9b, 0x66,
	0x59, 0xbe, 0x89, 0xed, 0x75, 0x5e, 0x66, 0x81, 0x4c, 0xc2, 0x74, 0xc1, 0x84, 0x5e, 0x7b, 0x5f,
	0x4c, 0x5e, 0x69, 0xf2, 0x4a, 0x93, 0xf7, 0x7d, 0x7d, 0x17, 0x7b, 0x45, 0x7b, 0xf6, 0x19, 0x00,
	0x00, 0xff, 0xff, 0x42, 0x2f, 0xd7, 0xc2, 0xd7, 0x02, 0x00, 0x00,
}
