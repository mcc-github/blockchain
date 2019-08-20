


package lifecycle

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common "github.com/mcc-github/blockchain/protos/common"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 



type InstallChaincodeArgs struct {
	ChaincodeInstallPackage []byte   `protobuf:"bytes,1,opt,name=chaincode_install_package,json=chaincodeInstallPackage,proto3" json:"chaincode_install_package,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *InstallChaincodeArgs) Reset()         { *m = InstallChaincodeArgs{} }
func (m *InstallChaincodeArgs) String() string { return proto.CompactTextString(m) }
func (*InstallChaincodeArgs) ProtoMessage()    {}
func (*InstallChaincodeArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{0}
}

func (m *InstallChaincodeArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallChaincodeArgs.Unmarshal(m, b)
}
func (m *InstallChaincodeArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallChaincodeArgs.Marshal(b, m, deterministic)
}
func (m *InstallChaincodeArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallChaincodeArgs.Merge(m, src)
}
func (m *InstallChaincodeArgs) XXX_Size() int {
	return xxx_messageInfo_InstallChaincodeArgs.Size(m)
}
func (m *InstallChaincodeArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallChaincodeArgs.DiscardUnknown(m)
}

var xxx_messageInfo_InstallChaincodeArgs proto.InternalMessageInfo

func (m *InstallChaincodeArgs) GetChaincodeInstallPackage() []byte {
	if m != nil {
		return m.ChaincodeInstallPackage
	}
	return nil
}



type InstallChaincodeResult struct {
	PackageId            string   `protobuf:"bytes,1,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	Label                string   `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallChaincodeResult) Reset()         { *m = InstallChaincodeResult{} }
func (m *InstallChaincodeResult) String() string { return proto.CompactTextString(m) }
func (*InstallChaincodeResult) ProtoMessage()    {}
func (*InstallChaincodeResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{1}
}

func (m *InstallChaincodeResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallChaincodeResult.Unmarshal(m, b)
}
func (m *InstallChaincodeResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallChaincodeResult.Marshal(b, m, deterministic)
}
func (m *InstallChaincodeResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallChaincodeResult.Merge(m, src)
}
func (m *InstallChaincodeResult) XXX_Size() int {
	return xxx_messageInfo_InstallChaincodeResult.Size(m)
}
func (m *InstallChaincodeResult) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallChaincodeResult.DiscardUnknown(m)
}

var xxx_messageInfo_InstallChaincodeResult proto.InternalMessageInfo

func (m *InstallChaincodeResult) GetPackageId() string {
	if m != nil {
		return m.PackageId
	}
	return ""
}

func (m *InstallChaincodeResult) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}



type QueryInstalledChaincodeArgs struct {
	PackageId            string   `protobuf:"bytes,1,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodeArgs) Reset()         { *m = QueryInstalledChaincodeArgs{} }
func (m *QueryInstalledChaincodeArgs) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodeArgs) ProtoMessage()    {}
func (*QueryInstalledChaincodeArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{2}
}

func (m *QueryInstalledChaincodeArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodeArgs.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodeArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodeArgs.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodeArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodeArgs.Merge(m, src)
}
func (m *QueryInstalledChaincodeArgs) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodeArgs.Size(m)
}
func (m *QueryInstalledChaincodeArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodeArgs.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodeArgs proto.InternalMessageInfo

func (m *QueryInstalledChaincodeArgs) GetPackageId() string {
	if m != nil {
		return m.PackageId
	}
	return ""
}



type QueryInstalledChaincodeResult struct {
	PackageId            string                                               `protobuf:"bytes,1,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	Label                string                                               `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	References           map[string]*QueryInstalledChaincodeResult_References `protobuf:"bytes,3,rep,name=references,proto3" json:"references,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                                             `json:"-"`
	XXX_unrecognized     []byte                                               `json:"-"`
	XXX_sizecache        int32                                                `json:"-"`
}

func (m *QueryInstalledChaincodeResult) Reset()         { *m = QueryInstalledChaincodeResult{} }
func (m *QueryInstalledChaincodeResult) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodeResult) ProtoMessage()    {}
func (*QueryInstalledChaincodeResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{3}
}

func (m *QueryInstalledChaincodeResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodeResult.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodeResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodeResult.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodeResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodeResult.Merge(m, src)
}
func (m *QueryInstalledChaincodeResult) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodeResult.Size(m)
}
func (m *QueryInstalledChaincodeResult) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodeResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodeResult proto.InternalMessageInfo

func (m *QueryInstalledChaincodeResult) GetPackageId() string {
	if m != nil {
		return m.PackageId
	}
	return ""
}

func (m *QueryInstalledChaincodeResult) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *QueryInstalledChaincodeResult) GetReferences() map[string]*QueryInstalledChaincodeResult_References {
	if m != nil {
		return m.References
	}
	return nil
}

type QueryInstalledChaincodeResult_References struct {
	Chaincodes           []*QueryInstalledChaincodeResult_Chaincode `protobuf:"bytes,1,rep,name=chaincodes,proto3" json:"chaincodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                   `json:"-"`
	XXX_unrecognized     []byte                                     `json:"-"`
	XXX_sizecache        int32                                      `json:"-"`
}

func (m *QueryInstalledChaincodeResult_References) Reset() {
	*m = QueryInstalledChaincodeResult_References{}
}
func (m *QueryInstalledChaincodeResult_References) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodeResult_References) ProtoMessage()    {}
func (*QueryInstalledChaincodeResult_References) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{3, 1}
}

func (m *QueryInstalledChaincodeResult_References) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodeResult_References.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodeResult_References) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodeResult_References.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodeResult_References) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodeResult_References.Merge(m, src)
}
func (m *QueryInstalledChaincodeResult_References) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodeResult_References.Size(m)
}
func (m *QueryInstalledChaincodeResult_References) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodeResult_References.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodeResult_References proto.InternalMessageInfo

func (m *QueryInstalledChaincodeResult_References) GetChaincodes() []*QueryInstalledChaincodeResult_Chaincode {
	if m != nil {
		return m.Chaincodes
	}
	return nil
}

type QueryInstalledChaincodeResult_Chaincode struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodeResult_Chaincode) Reset() {
	*m = QueryInstalledChaincodeResult_Chaincode{}
}
func (m *QueryInstalledChaincodeResult_Chaincode) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodeResult_Chaincode) ProtoMessage()    {}
func (*QueryInstalledChaincodeResult_Chaincode) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{3, 2}
}

func (m *QueryInstalledChaincodeResult_Chaincode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodeResult_Chaincode.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodeResult_Chaincode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodeResult_Chaincode.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodeResult_Chaincode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodeResult_Chaincode.Merge(m, src)
}
func (m *QueryInstalledChaincodeResult_Chaincode) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodeResult_Chaincode.Size(m)
}
func (m *QueryInstalledChaincodeResult_Chaincode) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodeResult_Chaincode.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodeResult_Chaincode proto.InternalMessageInfo

func (m *QueryInstalledChaincodeResult_Chaincode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *QueryInstalledChaincodeResult_Chaincode) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}



type GetInstalledChaincodePackageArgs struct {
	PackageId            string   `protobuf:"bytes,1,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInstalledChaincodePackageArgs) Reset()         { *m = GetInstalledChaincodePackageArgs{} }
func (m *GetInstalledChaincodePackageArgs) String() string { return proto.CompactTextString(m) }
func (*GetInstalledChaincodePackageArgs) ProtoMessage()    {}
func (*GetInstalledChaincodePackageArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{4}
}

func (m *GetInstalledChaincodePackageArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetInstalledChaincodePackageArgs.Unmarshal(m, b)
}
func (m *GetInstalledChaincodePackageArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetInstalledChaincodePackageArgs.Marshal(b, m, deterministic)
}
func (m *GetInstalledChaincodePackageArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetInstalledChaincodePackageArgs.Merge(m, src)
}
func (m *GetInstalledChaincodePackageArgs) XXX_Size() int {
	return xxx_messageInfo_GetInstalledChaincodePackageArgs.Size(m)
}
func (m *GetInstalledChaincodePackageArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_GetInstalledChaincodePackageArgs.DiscardUnknown(m)
}

var xxx_messageInfo_GetInstalledChaincodePackageArgs proto.InternalMessageInfo

func (m *GetInstalledChaincodePackageArgs) GetPackageId() string {
	if m != nil {
		return m.PackageId
	}
	return ""
}



type GetInstalledChaincodePackageResult struct {
	ChaincodeInstallPackage []byte   `protobuf:"bytes,1,opt,name=chaincode_install_package,json=chaincodeInstallPackage,proto3" json:"chaincode_install_package,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *GetInstalledChaincodePackageResult) Reset()         { *m = GetInstalledChaincodePackageResult{} }
func (m *GetInstalledChaincodePackageResult) String() string { return proto.CompactTextString(m) }
func (*GetInstalledChaincodePackageResult) ProtoMessage()    {}
func (*GetInstalledChaincodePackageResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{5}
}

func (m *GetInstalledChaincodePackageResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetInstalledChaincodePackageResult.Unmarshal(m, b)
}
func (m *GetInstalledChaincodePackageResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetInstalledChaincodePackageResult.Marshal(b, m, deterministic)
}
func (m *GetInstalledChaincodePackageResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetInstalledChaincodePackageResult.Merge(m, src)
}
func (m *GetInstalledChaincodePackageResult) XXX_Size() int {
	return xxx_messageInfo_GetInstalledChaincodePackageResult.Size(m)
}
func (m *GetInstalledChaincodePackageResult) XXX_DiscardUnknown() {
	xxx_messageInfo_GetInstalledChaincodePackageResult.DiscardUnknown(m)
}

var xxx_messageInfo_GetInstalledChaincodePackageResult proto.InternalMessageInfo

func (m *GetInstalledChaincodePackageResult) GetChaincodeInstallPackage() []byte {
	if m != nil {
		return m.ChaincodeInstallPackage
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
	return fileDescriptor_6625a5b20951add3, []int{6}
}

func (m *QueryInstalledChaincodesArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodesArgs.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodesArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodesArgs.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodesArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodesArgs.Merge(m, src)
}
func (m *QueryInstalledChaincodesArgs) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodesArgs.Size(m)
}
func (m *QueryInstalledChaincodesArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodesArgs.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodesArgs proto.InternalMessageInfo





type QueryInstalledChaincodesResult struct {
	InstalledChaincodes  []*QueryInstalledChaincodesResult_InstalledChaincode `protobuf:"bytes,1,rep,name=installed_chaincodes,json=installedChaincodes,proto3" json:"installed_chaincodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                             `json:"-"`
	XXX_unrecognized     []byte                                               `json:"-"`
	XXX_sizecache        int32                                                `json:"-"`
}

func (m *QueryInstalledChaincodesResult) Reset()         { *m = QueryInstalledChaincodesResult{} }
func (m *QueryInstalledChaincodesResult) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodesResult) ProtoMessage()    {}
func (*QueryInstalledChaincodesResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{7}
}

func (m *QueryInstalledChaincodesResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodesResult.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodesResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodesResult.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodesResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodesResult.Merge(m, src)
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
	PackageId            string                                                `protobuf:"bytes,1,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	Label                string                                                `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	References           map[string]*QueryInstalledChaincodesResult_References `protobuf:"bytes,3,rep,name=references,proto3" json:"references,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                                              `json:"-"`
	XXX_unrecognized     []byte                                                `json:"-"`
	XXX_sizecache        int32                                                 `json:"-"`
}

func (m *QueryInstalledChaincodesResult_InstalledChaincode) Reset() {
	*m = QueryInstalledChaincodesResult_InstalledChaincode{}
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) String() string {
	return proto.CompactTextString(m)
}
func (*QueryInstalledChaincodesResult_InstalledChaincode) ProtoMessage() {}
func (*QueryInstalledChaincodesResult_InstalledChaincode) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{7, 0}
}

func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.Merge(m, src)
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.Size(m)
}
func (m *QueryInstalledChaincodesResult_InstalledChaincode) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodesResult_InstalledChaincode proto.InternalMessageInfo

func (m *QueryInstalledChaincodesResult_InstalledChaincode) GetPackageId() string {
	if m != nil {
		return m.PackageId
	}
	return ""
}

func (m *QueryInstalledChaincodesResult_InstalledChaincode) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *QueryInstalledChaincodesResult_InstalledChaincode) GetReferences() map[string]*QueryInstalledChaincodesResult_References {
	if m != nil {
		return m.References
	}
	return nil
}

type QueryInstalledChaincodesResult_References struct {
	Chaincodes           []*QueryInstalledChaincodesResult_Chaincode `protobuf:"bytes,1,rep,name=chaincodes,proto3" json:"chaincodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *QueryInstalledChaincodesResult_References) Reset() {
	*m = QueryInstalledChaincodesResult_References{}
}
func (m *QueryInstalledChaincodesResult_References) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodesResult_References) ProtoMessage()    {}
func (*QueryInstalledChaincodesResult_References) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{7, 1}
}

func (m *QueryInstalledChaincodesResult_References) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodesResult_References.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodesResult_References) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodesResult_References.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodesResult_References) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodesResult_References.Merge(m, src)
}
func (m *QueryInstalledChaincodesResult_References) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodesResult_References.Size(m)
}
func (m *QueryInstalledChaincodesResult_References) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodesResult_References.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodesResult_References proto.InternalMessageInfo

func (m *QueryInstalledChaincodesResult_References) GetChaincodes() []*QueryInstalledChaincodesResult_Chaincode {
	if m != nil {
		return m.Chaincodes
	}
	return nil
}

type QueryInstalledChaincodesResult_Chaincode struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodesResult_Chaincode) Reset() {
	*m = QueryInstalledChaincodesResult_Chaincode{}
}
func (m *QueryInstalledChaincodesResult_Chaincode) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodesResult_Chaincode) ProtoMessage()    {}
func (*QueryInstalledChaincodesResult_Chaincode) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{7, 2}
}

func (m *QueryInstalledChaincodesResult_Chaincode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryInstalledChaincodesResult_Chaincode.Unmarshal(m, b)
}
func (m *QueryInstalledChaincodesResult_Chaincode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryInstalledChaincodesResult_Chaincode.Marshal(b, m, deterministic)
}
func (m *QueryInstalledChaincodesResult_Chaincode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryInstalledChaincodesResult_Chaincode.Merge(m, src)
}
func (m *QueryInstalledChaincodesResult_Chaincode) XXX_Size() int {
	return xxx_messageInfo_QueryInstalledChaincodesResult_Chaincode.Size(m)
}
func (m *QueryInstalledChaincodesResult_Chaincode) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryInstalledChaincodesResult_Chaincode.DiscardUnknown(m)
}

var xxx_messageInfo_QueryInstalledChaincodesResult_Chaincode proto.InternalMessageInfo

func (m *QueryInstalledChaincodesResult_Chaincode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *QueryInstalledChaincodesResult_Chaincode) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}



type ApproveChaincodeDefinitionForMyOrgArgs struct {
	Sequence             int64                           `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Name                 string                          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Version              string                          `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	EndorsementPlugin    string                          `protobuf:"bytes,4,opt,name=endorsement_plugin,json=endorsementPlugin,proto3" json:"endorsement_plugin,omitempty"`
	ValidationPlugin     string                          `protobuf:"bytes,5,opt,name=validation_plugin,json=validationPlugin,proto3" json:"validation_plugin,omitempty"`
	ValidationParameter  []byte                          `protobuf:"bytes,6,opt,name=validation_parameter,json=validationParameter,proto3" json:"validation_parameter,omitempty"`
	Collections          *common.CollectionConfigPackage `protobuf:"bytes,7,opt,name=collections,proto3" json:"collections,omitempty"`
	InitRequired         bool                            `protobuf:"varint,8,opt,name=init_required,json=initRequired,proto3" json:"init_required,omitempty"`
	Source               *ChaincodeSource                `protobuf:"bytes,9,opt,name=source,proto3" json:"source,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) Reset() {
	*m = ApproveChaincodeDefinitionForMyOrgArgs{}
}
func (m *ApproveChaincodeDefinitionForMyOrgArgs) String() string { return proto.CompactTextString(m) }
func (*ApproveChaincodeDefinitionForMyOrgArgs) ProtoMessage()    {}
func (*ApproveChaincodeDefinitionForMyOrgArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{8}
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs.Unmarshal(m, b)
}
func (m *ApproveChaincodeDefinitionForMyOrgArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs.Marshal(b, m, deterministic)
}
func (m *ApproveChaincodeDefinitionForMyOrgArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs.Merge(m, src)
}
func (m *ApproveChaincodeDefinitionForMyOrgArgs) XXX_Size() int {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs.Size(m)
}
func (m *ApproveChaincodeDefinitionForMyOrgArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs.DiscardUnknown(m)
}

var xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs proto.InternalMessageInfo

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetEndorsementPlugin() string {
	if m != nil {
		return m.EndorsementPlugin
	}
	return ""
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetValidationPlugin() string {
	if m != nil {
		return m.ValidationPlugin
	}
	return ""
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetValidationParameter() []byte {
	if m != nil {
		return m.ValidationParameter
	}
	return nil
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetCollections() *common.CollectionConfigPackage {
	if m != nil {
		return m.Collections
	}
	return nil
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetInitRequired() bool {
	if m != nil {
		return m.InitRequired
	}
	return false
}

func (m *ApproveChaincodeDefinitionForMyOrgArgs) GetSource() *ChaincodeSource {
	if m != nil {
		return m.Source
	}
	return nil
}

type ChaincodeSource struct {
	
	
	
	Type                 isChaincodeSource_Type `protobuf_oneof:"Type"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *ChaincodeSource) Reset()         { *m = ChaincodeSource{} }
func (m *ChaincodeSource) String() string { return proto.CompactTextString(m) }
func (*ChaincodeSource) ProtoMessage()    {}
func (*ChaincodeSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{9}
}

func (m *ChaincodeSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeSource.Unmarshal(m, b)
}
func (m *ChaincodeSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeSource.Marshal(b, m, deterministic)
}
func (m *ChaincodeSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeSource.Merge(m, src)
}
func (m *ChaincodeSource) XXX_Size() int {
	return xxx_messageInfo_ChaincodeSource.Size(m)
}
func (m *ChaincodeSource) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeSource.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeSource proto.InternalMessageInfo

type isChaincodeSource_Type interface {
	isChaincodeSource_Type()
}

type ChaincodeSource_Unavailable_ struct {
	Unavailable *ChaincodeSource_Unavailable `protobuf:"bytes,1,opt,name=unavailable,proto3,oneof"`
}

type ChaincodeSource_LocalPackage struct {
	LocalPackage *ChaincodeSource_Local `protobuf:"bytes,2,opt,name=local_package,json=localPackage,proto3,oneof"`
}

func (*ChaincodeSource_Unavailable_) isChaincodeSource_Type() {}

func (*ChaincodeSource_LocalPackage) isChaincodeSource_Type() {}

func (m *ChaincodeSource) GetType() isChaincodeSource_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *ChaincodeSource) GetUnavailable() *ChaincodeSource_Unavailable {
	if x, ok := m.GetType().(*ChaincodeSource_Unavailable_); ok {
		return x.Unavailable
	}
	return nil
}

func (m *ChaincodeSource) GetLocalPackage() *ChaincodeSource_Local {
	if x, ok := m.GetType().(*ChaincodeSource_LocalPackage); ok {
		return x.LocalPackage
	}
	return nil
}


func (*ChaincodeSource) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ChaincodeSource_Unavailable_)(nil),
		(*ChaincodeSource_LocalPackage)(nil),
	}
}

type ChaincodeSource_Unavailable struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeSource_Unavailable) Reset()         { *m = ChaincodeSource_Unavailable{} }
func (m *ChaincodeSource_Unavailable) String() string { return proto.CompactTextString(m) }
func (*ChaincodeSource_Unavailable) ProtoMessage()    {}
func (*ChaincodeSource_Unavailable) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{9, 0}
}

func (m *ChaincodeSource_Unavailable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeSource_Unavailable.Unmarshal(m, b)
}
func (m *ChaincodeSource_Unavailable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeSource_Unavailable.Marshal(b, m, deterministic)
}
func (m *ChaincodeSource_Unavailable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeSource_Unavailable.Merge(m, src)
}
func (m *ChaincodeSource_Unavailable) XXX_Size() int {
	return xxx_messageInfo_ChaincodeSource_Unavailable.Size(m)
}
func (m *ChaincodeSource_Unavailable) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeSource_Unavailable.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeSource_Unavailable proto.InternalMessageInfo

type ChaincodeSource_Local struct {
	PackageId            string   `protobuf:"bytes,1,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeSource_Local) Reset()         { *m = ChaincodeSource_Local{} }
func (m *ChaincodeSource_Local) String() string { return proto.CompactTextString(m) }
func (*ChaincodeSource_Local) ProtoMessage()    {}
func (*ChaincodeSource_Local) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{9, 1}
}

func (m *ChaincodeSource_Local) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeSource_Local.Unmarshal(m, b)
}
func (m *ChaincodeSource_Local) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeSource_Local.Marshal(b, m, deterministic)
}
func (m *ChaincodeSource_Local) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeSource_Local.Merge(m, src)
}
func (m *ChaincodeSource_Local) XXX_Size() int {
	return xxx_messageInfo_ChaincodeSource_Local.Size(m)
}
func (m *ChaincodeSource_Local) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeSource_Local.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeSource_Local proto.InternalMessageInfo

func (m *ChaincodeSource_Local) GetPackageId() string {
	if m != nil {
		return m.PackageId
	}
	return ""
}




type ApproveChaincodeDefinitionForMyOrgResult struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApproveChaincodeDefinitionForMyOrgResult) Reset() {
	*m = ApproveChaincodeDefinitionForMyOrgResult{}
}
func (m *ApproveChaincodeDefinitionForMyOrgResult) String() string { return proto.CompactTextString(m) }
func (*ApproveChaincodeDefinitionForMyOrgResult) ProtoMessage()    {}
func (*ApproveChaincodeDefinitionForMyOrgResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{10}
}

func (m *ApproveChaincodeDefinitionForMyOrgResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult.Unmarshal(m, b)
}
func (m *ApproveChaincodeDefinitionForMyOrgResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult.Marshal(b, m, deterministic)
}
func (m *ApproveChaincodeDefinitionForMyOrgResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult.Merge(m, src)
}
func (m *ApproveChaincodeDefinitionForMyOrgResult) XXX_Size() int {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult.Size(m)
}
func (m *ApproveChaincodeDefinitionForMyOrgResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult.DiscardUnknown(m)
}

var xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult proto.InternalMessageInfo



type CommitChaincodeDefinitionArgs struct {
	Sequence             int64                           `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Name                 string                          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Version              string                          `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	EndorsementPlugin    string                          `protobuf:"bytes,4,opt,name=endorsement_plugin,json=endorsementPlugin,proto3" json:"endorsement_plugin,omitempty"`
	ValidationPlugin     string                          `protobuf:"bytes,5,opt,name=validation_plugin,json=validationPlugin,proto3" json:"validation_plugin,omitempty"`
	ValidationParameter  []byte                          `protobuf:"bytes,6,opt,name=validation_parameter,json=validationParameter,proto3" json:"validation_parameter,omitempty"`
	Collections          *common.CollectionConfigPackage `protobuf:"bytes,7,opt,name=collections,proto3" json:"collections,omitempty"`
	InitRequired         bool                            `protobuf:"varint,8,opt,name=init_required,json=initRequired,proto3" json:"init_required,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *CommitChaincodeDefinitionArgs) Reset()         { *m = CommitChaincodeDefinitionArgs{} }
func (m *CommitChaincodeDefinitionArgs) String() string { return proto.CompactTextString(m) }
func (*CommitChaincodeDefinitionArgs) ProtoMessage()    {}
func (*CommitChaincodeDefinitionArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{11}
}

func (m *CommitChaincodeDefinitionArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommitChaincodeDefinitionArgs.Unmarshal(m, b)
}
func (m *CommitChaincodeDefinitionArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommitChaincodeDefinitionArgs.Marshal(b, m, deterministic)
}
func (m *CommitChaincodeDefinitionArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitChaincodeDefinitionArgs.Merge(m, src)
}
func (m *CommitChaincodeDefinitionArgs) XXX_Size() int {
	return xxx_messageInfo_CommitChaincodeDefinitionArgs.Size(m)
}
func (m *CommitChaincodeDefinitionArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitChaincodeDefinitionArgs.DiscardUnknown(m)
}

var xxx_messageInfo_CommitChaincodeDefinitionArgs proto.InternalMessageInfo

func (m *CommitChaincodeDefinitionArgs) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *CommitChaincodeDefinitionArgs) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CommitChaincodeDefinitionArgs) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *CommitChaincodeDefinitionArgs) GetEndorsementPlugin() string {
	if m != nil {
		return m.EndorsementPlugin
	}
	return ""
}

func (m *CommitChaincodeDefinitionArgs) GetValidationPlugin() string {
	if m != nil {
		return m.ValidationPlugin
	}
	return ""
}

func (m *CommitChaincodeDefinitionArgs) GetValidationParameter() []byte {
	if m != nil {
		return m.ValidationParameter
	}
	return nil
}

func (m *CommitChaincodeDefinitionArgs) GetCollections() *common.CollectionConfigPackage {
	if m != nil {
		return m.Collections
	}
	return nil
}

func (m *CommitChaincodeDefinitionArgs) GetInitRequired() bool {
	if m != nil {
		return m.InitRequired
	}
	return false
}




type CommitChaincodeDefinitionResult struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommitChaincodeDefinitionResult) Reset()         { *m = CommitChaincodeDefinitionResult{} }
func (m *CommitChaincodeDefinitionResult) String() string { return proto.CompactTextString(m) }
func (*CommitChaincodeDefinitionResult) ProtoMessage()    {}
func (*CommitChaincodeDefinitionResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{12}
}

func (m *CommitChaincodeDefinitionResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommitChaincodeDefinitionResult.Unmarshal(m, b)
}
func (m *CommitChaincodeDefinitionResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommitChaincodeDefinitionResult.Marshal(b, m, deterministic)
}
func (m *CommitChaincodeDefinitionResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitChaincodeDefinitionResult.Merge(m, src)
}
func (m *CommitChaincodeDefinitionResult) XXX_Size() int {
	return xxx_messageInfo_CommitChaincodeDefinitionResult.Size(m)
}
func (m *CommitChaincodeDefinitionResult) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitChaincodeDefinitionResult.DiscardUnknown(m)
}

var xxx_messageInfo_CommitChaincodeDefinitionResult proto.InternalMessageInfo



type CheckCommitReadinessArgs struct {
	Sequence             int64                           `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Name                 string                          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Version              string                          `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	EndorsementPlugin    string                          `protobuf:"bytes,4,opt,name=endorsement_plugin,json=endorsementPlugin,proto3" json:"endorsement_plugin,omitempty"`
	ValidationPlugin     string                          `protobuf:"bytes,5,opt,name=validation_plugin,json=validationPlugin,proto3" json:"validation_plugin,omitempty"`
	ValidationParameter  []byte                          `protobuf:"bytes,6,opt,name=validation_parameter,json=validationParameter,proto3" json:"validation_parameter,omitempty"`
	Collections          *common.CollectionConfigPackage `protobuf:"bytes,7,opt,name=collections,proto3" json:"collections,omitempty"`
	InitRequired         bool                            `protobuf:"varint,8,opt,name=init_required,json=initRequired,proto3" json:"init_required,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *CheckCommitReadinessArgs) Reset()         { *m = CheckCommitReadinessArgs{} }
func (m *CheckCommitReadinessArgs) String() string { return proto.CompactTextString(m) }
func (*CheckCommitReadinessArgs) ProtoMessage()    {}
func (*CheckCommitReadinessArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{13}
}

func (m *CheckCommitReadinessArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckCommitReadinessArgs.Unmarshal(m, b)
}
func (m *CheckCommitReadinessArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckCommitReadinessArgs.Marshal(b, m, deterministic)
}
func (m *CheckCommitReadinessArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckCommitReadinessArgs.Merge(m, src)
}
func (m *CheckCommitReadinessArgs) XXX_Size() int {
	return xxx_messageInfo_CheckCommitReadinessArgs.Size(m)
}
func (m *CheckCommitReadinessArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckCommitReadinessArgs.DiscardUnknown(m)
}

var xxx_messageInfo_CheckCommitReadinessArgs proto.InternalMessageInfo

func (m *CheckCommitReadinessArgs) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *CheckCommitReadinessArgs) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CheckCommitReadinessArgs) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *CheckCommitReadinessArgs) GetEndorsementPlugin() string {
	if m != nil {
		return m.EndorsementPlugin
	}
	return ""
}

func (m *CheckCommitReadinessArgs) GetValidationPlugin() string {
	if m != nil {
		return m.ValidationPlugin
	}
	return ""
}

func (m *CheckCommitReadinessArgs) GetValidationParameter() []byte {
	if m != nil {
		return m.ValidationParameter
	}
	return nil
}

func (m *CheckCommitReadinessArgs) GetCollections() *common.CollectionConfigPackage {
	if m != nil {
		return m.Collections
	}
	return nil
}

func (m *CheckCommitReadinessArgs) GetInitRequired() bool {
	if m != nil {
		return m.InitRequired
	}
	return false
}





type CheckCommitReadinessResult struct {
	Approvals            map[string]bool `protobuf:"bytes,1,rep,name=approvals,proto3" json:"approvals,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CheckCommitReadinessResult) Reset()         { *m = CheckCommitReadinessResult{} }
func (m *CheckCommitReadinessResult) String() string { return proto.CompactTextString(m) }
func (*CheckCommitReadinessResult) ProtoMessage()    {}
func (*CheckCommitReadinessResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{14}
}

func (m *CheckCommitReadinessResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckCommitReadinessResult.Unmarshal(m, b)
}
func (m *CheckCommitReadinessResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckCommitReadinessResult.Marshal(b, m, deterministic)
}
func (m *CheckCommitReadinessResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckCommitReadinessResult.Merge(m, src)
}
func (m *CheckCommitReadinessResult) XXX_Size() int {
	return xxx_messageInfo_CheckCommitReadinessResult.Size(m)
}
func (m *CheckCommitReadinessResult) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckCommitReadinessResult.DiscardUnknown(m)
}

var xxx_messageInfo_CheckCommitReadinessResult proto.InternalMessageInfo

func (m *CheckCommitReadinessResult) GetApprovals() map[string]bool {
	if m != nil {
		return m.Approvals
	}
	return nil
}



type QueryChaincodeDefinitionArgs struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryChaincodeDefinitionArgs) Reset()         { *m = QueryChaincodeDefinitionArgs{} }
func (m *QueryChaincodeDefinitionArgs) String() string { return proto.CompactTextString(m) }
func (*QueryChaincodeDefinitionArgs) ProtoMessage()    {}
func (*QueryChaincodeDefinitionArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{15}
}

func (m *QueryChaincodeDefinitionArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryChaincodeDefinitionArgs.Unmarshal(m, b)
}
func (m *QueryChaincodeDefinitionArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryChaincodeDefinitionArgs.Marshal(b, m, deterministic)
}
func (m *QueryChaincodeDefinitionArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryChaincodeDefinitionArgs.Merge(m, src)
}
func (m *QueryChaincodeDefinitionArgs) XXX_Size() int {
	return xxx_messageInfo_QueryChaincodeDefinitionArgs.Size(m)
}
func (m *QueryChaincodeDefinitionArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryChaincodeDefinitionArgs.DiscardUnknown(m)
}

var xxx_messageInfo_QueryChaincodeDefinitionArgs proto.InternalMessageInfo

func (m *QueryChaincodeDefinitionArgs) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}



type QueryChaincodeDefinitionResult struct {
	Sequence             int64                           `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Version              string                          `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	EndorsementPlugin    string                          `protobuf:"bytes,3,opt,name=endorsement_plugin,json=endorsementPlugin,proto3" json:"endorsement_plugin,omitempty"`
	ValidationPlugin     string                          `protobuf:"bytes,4,opt,name=validation_plugin,json=validationPlugin,proto3" json:"validation_plugin,omitempty"`
	ValidationParameter  []byte                          `protobuf:"bytes,5,opt,name=validation_parameter,json=validationParameter,proto3" json:"validation_parameter,omitempty"`
	Collections          *common.CollectionConfigPackage `protobuf:"bytes,6,opt,name=collections,proto3" json:"collections,omitempty"`
	InitRequired         bool                            `protobuf:"varint,7,opt,name=init_required,json=initRequired,proto3" json:"init_required,omitempty"`
	Approvals            map[string]bool                 `protobuf:"bytes,8,rep,name=approvals,proto3" json:"approvals,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *QueryChaincodeDefinitionResult) Reset()         { *m = QueryChaincodeDefinitionResult{} }
func (m *QueryChaincodeDefinitionResult) String() string { return proto.CompactTextString(m) }
func (*QueryChaincodeDefinitionResult) ProtoMessage()    {}
func (*QueryChaincodeDefinitionResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{16}
}

func (m *QueryChaincodeDefinitionResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryChaincodeDefinitionResult.Unmarshal(m, b)
}
func (m *QueryChaincodeDefinitionResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryChaincodeDefinitionResult.Marshal(b, m, deterministic)
}
func (m *QueryChaincodeDefinitionResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryChaincodeDefinitionResult.Merge(m, src)
}
func (m *QueryChaincodeDefinitionResult) XXX_Size() int {
	return xxx_messageInfo_QueryChaincodeDefinitionResult.Size(m)
}
func (m *QueryChaincodeDefinitionResult) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryChaincodeDefinitionResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryChaincodeDefinitionResult proto.InternalMessageInfo

func (m *QueryChaincodeDefinitionResult) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *QueryChaincodeDefinitionResult) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *QueryChaincodeDefinitionResult) GetEndorsementPlugin() string {
	if m != nil {
		return m.EndorsementPlugin
	}
	return ""
}

func (m *QueryChaincodeDefinitionResult) GetValidationPlugin() string {
	if m != nil {
		return m.ValidationPlugin
	}
	return ""
}

func (m *QueryChaincodeDefinitionResult) GetValidationParameter() []byte {
	if m != nil {
		return m.ValidationParameter
	}
	return nil
}

func (m *QueryChaincodeDefinitionResult) GetCollections() *common.CollectionConfigPackage {
	if m != nil {
		return m.Collections
	}
	return nil
}

func (m *QueryChaincodeDefinitionResult) GetInitRequired() bool {
	if m != nil {
		return m.InitRequired
	}
	return false
}

func (m *QueryChaincodeDefinitionResult) GetApprovals() map[string]bool {
	if m != nil {
		return m.Approvals
	}
	return nil
}



type QueryChaincodeDefinitionsArgs struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryChaincodeDefinitionsArgs) Reset()         { *m = QueryChaincodeDefinitionsArgs{} }
func (m *QueryChaincodeDefinitionsArgs) String() string { return proto.CompactTextString(m) }
func (*QueryChaincodeDefinitionsArgs) ProtoMessage()    {}
func (*QueryChaincodeDefinitionsArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{17}
}

func (m *QueryChaincodeDefinitionsArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryChaincodeDefinitionsArgs.Unmarshal(m, b)
}
func (m *QueryChaincodeDefinitionsArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryChaincodeDefinitionsArgs.Marshal(b, m, deterministic)
}
func (m *QueryChaincodeDefinitionsArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryChaincodeDefinitionsArgs.Merge(m, src)
}
func (m *QueryChaincodeDefinitionsArgs) XXX_Size() int {
	return xxx_messageInfo_QueryChaincodeDefinitionsArgs.Size(m)
}
func (m *QueryChaincodeDefinitionsArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryChaincodeDefinitionsArgs.DiscardUnknown(m)
}

var xxx_messageInfo_QueryChaincodeDefinitionsArgs proto.InternalMessageInfo



type QueryChaincodeDefinitionsResult struct {
	ChaincodeDefinitions []*QueryChaincodeDefinitionsResult_ChaincodeDefinition `protobuf:"bytes,1,rep,name=chaincode_definitions,json=chaincodeDefinitions,proto3" json:"chaincode_definitions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                               `json:"-"`
	XXX_unrecognized     []byte                                                 `json:"-"`
	XXX_sizecache        int32                                                  `json:"-"`
}

func (m *QueryChaincodeDefinitionsResult) Reset()         { *m = QueryChaincodeDefinitionsResult{} }
func (m *QueryChaincodeDefinitionsResult) String() string { return proto.CompactTextString(m) }
func (*QueryChaincodeDefinitionsResult) ProtoMessage()    {}
func (*QueryChaincodeDefinitionsResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{18}
}

func (m *QueryChaincodeDefinitionsResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryChaincodeDefinitionsResult.Unmarshal(m, b)
}
func (m *QueryChaincodeDefinitionsResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryChaincodeDefinitionsResult.Marshal(b, m, deterministic)
}
func (m *QueryChaincodeDefinitionsResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryChaincodeDefinitionsResult.Merge(m, src)
}
func (m *QueryChaincodeDefinitionsResult) XXX_Size() int {
	return xxx_messageInfo_QueryChaincodeDefinitionsResult.Size(m)
}
func (m *QueryChaincodeDefinitionsResult) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryChaincodeDefinitionsResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryChaincodeDefinitionsResult proto.InternalMessageInfo

func (m *QueryChaincodeDefinitionsResult) GetChaincodeDefinitions() []*QueryChaincodeDefinitionsResult_ChaincodeDefinition {
	if m != nil {
		return m.ChaincodeDefinitions
	}
	return nil
}

type QueryChaincodeDefinitionsResult_ChaincodeDefinition struct {
	Name                 string                          `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Sequence             int64                           `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Version              string                          `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	EndorsementPlugin    string                          `protobuf:"bytes,4,opt,name=endorsement_plugin,json=endorsementPlugin,proto3" json:"endorsement_plugin,omitempty"`
	ValidationPlugin     string                          `protobuf:"bytes,5,opt,name=validation_plugin,json=validationPlugin,proto3" json:"validation_plugin,omitempty"`
	ValidationParameter  []byte                          `protobuf:"bytes,6,opt,name=validation_parameter,json=validationParameter,proto3" json:"validation_parameter,omitempty"`
	Collections          *common.CollectionConfigPackage `protobuf:"bytes,7,opt,name=collections,proto3" json:"collections,omitempty"`
	InitRequired         bool                            `protobuf:"varint,8,opt,name=init_required,json=initRequired,proto3" json:"init_required,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) Reset() {
	*m = QueryChaincodeDefinitionsResult_ChaincodeDefinition{}
}
func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) String() string {
	return proto.CompactTextString(m)
}
func (*QueryChaincodeDefinitionsResult_ChaincodeDefinition) ProtoMessage() {}
func (*QueryChaincodeDefinitionsResult_ChaincodeDefinition) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{18, 0}
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryChaincodeDefinitionsResult_ChaincodeDefinition.Unmarshal(m, b)
}
func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryChaincodeDefinitionsResult_ChaincodeDefinition.Marshal(b, m, deterministic)
}
func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryChaincodeDefinitionsResult_ChaincodeDefinition.Merge(m, src)
}
func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) XXX_Size() int {
	return xxx_messageInfo_QueryChaincodeDefinitionsResult_ChaincodeDefinition.Size(m)
}
func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryChaincodeDefinitionsResult_ChaincodeDefinition.DiscardUnknown(m)
}

var xxx_messageInfo_QueryChaincodeDefinitionsResult_ChaincodeDefinition proto.InternalMessageInfo

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) GetEndorsementPlugin() string {
	if m != nil {
		return m.EndorsementPlugin
	}
	return ""
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) GetValidationPlugin() string {
	if m != nil {
		return m.ValidationPlugin
	}
	return ""
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) GetValidationParameter() []byte {
	if m != nil {
		return m.ValidationParameter
	}
	return nil
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) GetCollections() *common.CollectionConfigPackage {
	if m != nil {
		return m.Collections
	}
	return nil
}

func (m *QueryChaincodeDefinitionsResult_ChaincodeDefinition) GetInitRequired() bool {
	if m != nil {
		return m.InitRequired
	}
	return false
}

func init() {
	proto.RegisterType((*InstallChaincodeArgs)(nil), "lifecycle.InstallChaincodeArgs")
	proto.RegisterType((*InstallChaincodeResult)(nil), "lifecycle.InstallChaincodeResult")
	proto.RegisterType((*QueryInstalledChaincodeArgs)(nil), "lifecycle.QueryInstalledChaincodeArgs")
	proto.RegisterType((*QueryInstalledChaincodeResult)(nil), "lifecycle.QueryInstalledChaincodeResult")
	proto.RegisterMapType((map[string]*QueryInstalledChaincodeResult_References)(nil), "lifecycle.QueryInstalledChaincodeResult.ReferencesEntry")
	proto.RegisterType((*QueryInstalledChaincodeResult_References)(nil), "lifecycle.QueryInstalledChaincodeResult.References")
	proto.RegisterType((*QueryInstalledChaincodeResult_Chaincode)(nil), "lifecycle.QueryInstalledChaincodeResult.Chaincode")
	proto.RegisterType((*GetInstalledChaincodePackageArgs)(nil), "lifecycle.GetInstalledChaincodePackageArgs")
	proto.RegisterType((*GetInstalledChaincodePackageResult)(nil), "lifecycle.GetInstalledChaincodePackageResult")
	proto.RegisterType((*QueryInstalledChaincodesArgs)(nil), "lifecycle.QueryInstalledChaincodesArgs")
	proto.RegisterType((*QueryInstalledChaincodesResult)(nil), "lifecycle.QueryInstalledChaincodesResult")
	proto.RegisterType((*QueryInstalledChaincodesResult_InstalledChaincode)(nil), "lifecycle.QueryInstalledChaincodesResult.InstalledChaincode")
	proto.RegisterMapType((map[string]*QueryInstalledChaincodesResult_References)(nil), "lifecycle.QueryInstalledChaincodesResult.InstalledChaincode.ReferencesEntry")
	proto.RegisterType((*QueryInstalledChaincodesResult_References)(nil), "lifecycle.QueryInstalledChaincodesResult.References")
	proto.RegisterType((*QueryInstalledChaincodesResult_Chaincode)(nil), "lifecycle.QueryInstalledChaincodesResult.Chaincode")
	proto.RegisterType((*ApproveChaincodeDefinitionForMyOrgArgs)(nil), "lifecycle.ApproveChaincodeDefinitionForMyOrgArgs")
	proto.RegisterType((*ChaincodeSource)(nil), "lifecycle.ChaincodeSource")
	proto.RegisterType((*ChaincodeSource_Unavailable)(nil), "lifecycle.ChaincodeSource.Unavailable")
	proto.RegisterType((*ChaincodeSource_Local)(nil), "lifecycle.ChaincodeSource.Local")
	proto.RegisterType((*ApproveChaincodeDefinitionForMyOrgResult)(nil), "lifecycle.ApproveChaincodeDefinitionForMyOrgResult")
	proto.RegisterType((*CommitChaincodeDefinitionArgs)(nil), "lifecycle.CommitChaincodeDefinitionArgs")
	proto.RegisterType((*CommitChaincodeDefinitionResult)(nil), "lifecycle.CommitChaincodeDefinitionResult")
	proto.RegisterType((*CheckCommitReadinessArgs)(nil), "lifecycle.CheckCommitReadinessArgs")
	proto.RegisterType((*CheckCommitReadinessResult)(nil), "lifecycle.CheckCommitReadinessResult")
	proto.RegisterMapType((map[string]bool)(nil), "lifecycle.CheckCommitReadinessResult.ApprovalsEntry")
	proto.RegisterType((*QueryChaincodeDefinitionArgs)(nil), "lifecycle.QueryChaincodeDefinitionArgs")
	proto.RegisterType((*QueryChaincodeDefinitionResult)(nil), "lifecycle.QueryChaincodeDefinitionResult")
	proto.RegisterMapType((map[string]bool)(nil), "lifecycle.QueryChaincodeDefinitionResult.ApprovalsEntry")
	proto.RegisterType((*QueryChaincodeDefinitionsArgs)(nil), "lifecycle.QueryChaincodeDefinitionsArgs")
	proto.RegisterType((*QueryChaincodeDefinitionsResult)(nil), "lifecycle.QueryChaincodeDefinitionsResult")
	proto.RegisterType((*QueryChaincodeDefinitionsResult_ChaincodeDefinition)(nil), "lifecycle.QueryChaincodeDefinitionsResult.ChaincodeDefinition")
}

func init() { proto.RegisterFile("peer/lifecycle/lifecycle.proto", fileDescriptor_6625a5b20951add3) }

var fileDescriptor_6625a5b20951add3 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x58, 0xdd, 0x8e, 0xdb, 0x44,
	0x14, 0x6e, 0xe2, 0x24, 0x4d, 0x4e, 0x76, 0x69, 0x3b, 0x1b, 0xa8, 0x31, 0xec, 0x6e, 0x30, 0xd2,
	0x6a, 0xc5, 0x8f, 0x23, 0xb2, 0xbd, 0x28, 0xd5, 0x0a, 0x29, 0x0d, 0xd0, 0x6e, 0xd5, 0x8a, 0xe2,
	0x02, 0x42, 0xdc, 0xa4, 0x13, 0xfb, 0x24, 0x3b, 0xda, 0x89, 0x9d, 0x8e, 0x9d, 0x48, 0x79, 0x18,
	0xde, 0x00, 0xf1, 0x0a, 0xbc, 0x05, 0x37, 0x48, 0x08, 0x09, 0x71, 0xcd, 0x2b, 0xa0, 0x8c, 0x27,
	0xb6, 0xb3, 0xb1, 0xb3, 0x69, 0x77, 0xb9, 0xdb, 0x3b, 0xcf, 0x9c, 0xef, 0xfc, 0xcc, 0x39, 0xdf,
	0x99, 0x33, 0x09, 0xec, 0x8d, 0x11, 0x45, 0x8b, 0xb3, 0x01, 0x3a, 0x33, 0x87, 0x63, 0xf2, 0x65,
	0x8d, 0x85, 0x1f, 0xfa, 0xa4, 0x16, 0x6f, 0x18, 0x77, 0x1d, 0x7f, 0x34, 0xf2, 0xbd, 0x96, 0xe3,
	0x73, 0x8e, 0x4e, 0xc8, 0x7c, 0x2f, 0xc2, 0x98, 0x36, 0x34, 0x4e, 0xbc, 0x20, 0xa4, 0x9c, 0x77,
	0x4f, 0x29, 0xf3, 0x1c, 0xdf, 0xc5, 0x8e, 0x18, 0x06, 0xe4, 0x01, 0xbc, 0xeb, 0x2c, 0x36, 0x7a,
	0x2c, 0x42, 0xf4, 0xc6, 0xd4, 0x39, 0xa3, 0x43, 0xd4, 0x0b, 0xcd, 0xc2, 0xe1, 0x96, 0x7d, 0x37,
	0x06, 0x28, 0x0b, 0xcf, 0x23, 0xb1, 0xf9, 0x0c, 0xde, 0x39, 0x6f, 0xd3, 0xc6, 0x60, 0xc2, 0x43,
	0xb2, 0x0b, 0xa0, 0x6c, 0xf4, 0x98, 0x2b, 0xcd, 0xd4, 0xec, 0x9a, 0xda, 0x39, 0x71, 0x49, 0x03,
	0xca, 0x9c, 0xf6, 0x91, 0xeb, 0x45, 0x29, 0x89, 0x16, 0xe6, 0x31, 0xbc, 0xf7, 0xed, 0x04, 0xc5,
	0x4c, 0xd9, 0x44, 0x77, 0x39, 0xd2, 0xf5, 0x36, 0xcd, 0xdf, 0x34, 0xd8, 0xcd, 0x51, 0xbf, 0x44,
	0x50, 0xe4, 0x47, 0x00, 0x81, 0x03, 0x14, 0xe8, 0x39, 0x18, 0xe8, 0x5a, 0x53, 0x3b, 0xac, 0xb7,
	0xef, 0x5b, 0x49, 0x05, 0xd6, 0xba, 0xb4, 0xec, 0x58, 0xf5, 0x2b, 0x2f, 0x14, 0x33, 0x3b, 0x65,
	0xcb, 0x10, 0x70, 0xeb, 0x9c, 0x98, 0xdc, 0x06, 0xed, 0x0c, 0x67, 0x2a, 0xb4, 0xf9, 0x27, 0x39,
	0x81, 0xf2, 0x94, 0xf2, 0x09, 0xca, 0xa0, 0xea, 0xed, 0xa3, 0x37, 0xf0, 0x6c, 0x47, 0x16, 0x1e,
	0x14, 0xef, 0x17, 0x8c, 0x97, 0x00, 0x89, 0x80, 0xd8, 0x00, 0x71, 0x69, 0x03, 0xbd, 0x20, 0xcf,
	0xd6, 0xde, 0xd8, 0x43, 0xb2, 0x4e, 0x59, 0x31, 0x3e, 0x87, 0x5a, 0x2c, 0x20, 0x04, 0x4a, 0x1e,
	0x1d, 0xa1, 0x3a, 0x90, 0xfc, 0x26, 0x3a, 0xdc, 0x9c, 0xa2, 0x08, 0x98, 0xef, 0xa9, 0x44, 0x2f,
	0x96, 0x66, 0x07, 0x9a, 0x8f, 0x30, 0x5c, 0xf5, 0xa7, 0xe8, 0xb6, 0x09, 0x09, 0x5e, 0x82, 0xb9,
	0xce, 0x84, 0x22, 0xc2, 0x65, 0x38, 0xbf, 0x07, 0xef, 0xe7, 0xa4, 0x25, 0x98, 0x07, 0x68, 0xfe,
	0x59, 0x82, 0xbd, 0x3c, 0x80, 0x72, 0xef, 0x43, 0x83, 0x2d, 0x84, 0xbd, 0x95, 0x02, 0x1c, 0x5f,
	0x5c, 0x00, 0x65, 0xc8, 0xca, 0x28, 0xcd, 0x0e, 0x5b, 0x45, 0x1b, 0xbf, 0x14, 0x81, 0xac, 0x62,
	0xdf, 0xac, 0x1f, 0x78, 0x46, 0x3f, 0x3c, 0xbd, 0x4c, 0xc8, 0x6b, 0x7b, 0x24, 0xd8, 0xa4, 0x47,
	0x9e, 0x2c, 0xf7, 0xc8, 0xbd, 0xcd, 0xa3, 0xc9, 0x6e, 0x12, 0xba, 0xd4, 0x24, 0x2f, 0x32, 0x9a,
	0xe4, 0x68, 0x73, 0x17, 0x57, 0xde, 0x25, 0x3f, 0x6b, 0x70, 0xd0, 0x19, 0x8f, 0x85, 0x3f, 0xc5,
	0xd8, 0xc4, 0x97, 0x38, 0x60, 0x1e, 0x9b, 0xdf, 0xf6, 0x5f, 0xfb, 0xe2, 0xd9, 0xec, 0x1b, 0x31,
	0x94, 0xcd, 0x62, 0x40, 0x35, 0xc0, 0x57, 0x93, 0xf9, 0x39, 0xa4, 0x71, 0xcd, 0x8e, 0xd7, 0xb1,
	0xd3, 0x62, 0xb6, 0x53, 0x6d, 0xc9, 0x29, 0xf9, 0x14, 0x08, 0x7a, 0xae, 0x2f, 0x02, 0x1c, 0xa1,
	0x17, 0xf6, 0xc6, 0x7c, 0x32, 0x64, 0x9e, 0x5e, 0x92, 0xa0, 0x3b, 0x29, 0xc9, 0x73, 0x29, 0x20,
	0x1f, 0xc3, 0x9d, 0x29, 0xe5, 0xcc, 0xa5, 0xf3, 0x90, 0x16, 0xe8, 0xb2, 0x44, 0xdf, 0x4e, 0x04,
	0x0a, 0xfc, 0x19, 0x34, 0xd2, 0x60, 0x2a, 0xe8, 0x08, 0x43, 0x14, 0x7a, 0x45, 0x36, 0xe2, 0x4e,
	0x0a, 0xbf, 0x10, 0x91, 0x0e, 0xd4, 0x93, 0x01, 0x17, 0xe8, 0x37, 0x65, 0xdd, 0xf7, 0xad, 0x68,
	0xf6, 0x59, 0xdd, 0x58, 0xd4, 0xf5, 0xbd, 0x01, 0x1b, 0x2e, 0x9a, 0x3f, 0xad, 0x43, 0x3e, 0x84,
	0xed, 0x79, 0xca, 0x7a, 0x02, 0x5f, 0x4d, 0x98, 0x40, 0x57, 0xaf, 0x36, 0x0b, 0x87, 0x55, 0x7b,
	0x6b, 0xbe, 0x69, 0xab, 0x3d, 0xd2, 0x86, 0x4a, 0xe0, 0x4f, 0x84, 0x83, 0x7a, 0x4d, 0xba, 0x30,
	0x52, 0x75, 0x8f, 0x93, 0xff, 0x42, 0x22, 0x6c, 0x85, 0x34, 0xff, 0x29, 0xc0, 0xad, 0x73, 0x32,
	0xf2, 0x04, 0xea, 0x13, 0x8f, 0x4e, 0x29, 0xe3, 0xb4, 0xcf, 0xa3, 0x5a, 0xd4, 0xdb, 0x07, 0xf9,
	0xc6, 0xac, 0xef, 0x13, 0xf4, 0xe3, 0x1b, 0x76, 0x5a, 0x99, 0x3c, 0x82, 0x6d, 0xee, 0x3b, 0x34,
	0xb9, 0xb0, 0x22, 0xd6, 0x37, 0xd7, 0x58, 0x7b, 0x3a, 0xc7, 0x3f, 0xbe, 0x61, 0x6f, 0x49, 0x45,
	0x95, 0x0e, 0x63, 0x1b, 0xea, 0x29, 0x37, 0xc6, 0x01, 0x94, 0x25, 0xee, 0x82, 0x6b, 0xe1, 0x61,
	0x05, 0x4a, 0xdf, 0xcd, 0xc6, 0x68, 0x7e, 0x04, 0x87, 0x17, 0xd3, 0x30, 0x6a, 0x02, 0xf3, 0xaf,
	0x22, 0xec, 0x76, 0xfd, 0xd1, 0x88, 0x85, 0x19, 0xd8, 0x6b, 0xaa, 0x5e, 0x01, 0x55, 0xcd, 0x0f,
	0x60, 0x3f, 0x37, 0xc3, 0xaa, 0x0a, 0x7f, 0x14, 0x41, 0xef, 0x9e, 0xa2, 0x73, 0x16, 0x01, 0x6d,
	0xa4, 0x2e, 0xf3, 0x30, 0x08, 0xae, 0x0b, 0x70, 0x15, 0x05, 0xf8, 0xb5, 0x00, 0x46, 0x56, 0x76,
	0xd5, 0xd0, 0xb7, 0xa1, 0x46, 0x65, 0xbb, 0x50, 0xbe, 0x98, 0x22, 0xf7, 0x96, 0x5a, 0x36, 0x4f,
	0xd3, 0xea, 0x2c, 0xd4, 0xa2, 0xf1, 0x98, 0x98, 0x31, 0x8e, 0xe1, 0xad, 0x65, 0x61, 0xc6, 0x70,
	0x6c, 0xa4, 0x87, 0x63, 0x35, 0x35, 0xe6, 0xcc, 0xb6, 0x7a, 0xc9, 0xe4, 0xb5, 0x64, 0xc6, 0x58,
	0x32, 0xff, 0xd6, 0xd4, 0xeb, 0x26, 0x97, 0x65, 0x6b, 0x89, 0x94, 0x3b, 0xd5, 0x72, 0x48, 0xa3,
	0xbd, 0x16, 0x69, 0x4a, 0xaf, 0x49, 0x9a, 0xf2, 0xc6, 0xa4, 0xa9, 0x5c, 0x05, 0x69, 0x6e, 0x66,
	0x0c, 0x98, 0x1f, 0xd2, 0xac, 0xa8, 0x66, 0xff, 0xb8, 0xc8, 0x4d, 0xf5, 0xff, 0xc6, 0x8c, 0x7d,
	0xf5, 0x4b, 0x2a, 0xc3, 0x73, 0xf4, 0xc8, 0xfd, 0x57, 0x83, 0xfd, 0x5c, 0x84, 0xe2, 0x41, 0x00,
	0x6f, 0x27, 0x8f, 0x6c, 0x37, 0x11, 0x2b, 0xf2, 0x7f, 0xb1, 0xc1, 0x31, 0x57, 0xde, 0x50, 0xa9,
	0x0c, 0x34, 0x9c, 0x0c, 0xbc, 0xf1, 0x7b, 0x11, 0x76, 0x32, 0xd0, 0x99, 0x4f, 0xac, 0x34, 0x51,
	0x8b, 0xf9, 0x44, 0xbd, 0xbe, 0xdd, 0x04, 0xba, 0x0f, 0x1d, 0xf8, 0xc4, 0x17, 0x43, 0xeb, 0x74,
	0x36, 0x46, 0xc1, 0xd1, 0x1d, 0xa2, 0xb0, 0x06, 0xb4, 0x2f, 0x98, 0x13, 0xfd, 0xbd, 0x10, 0x58,
	0x63, 0x44, 0x91, 0x94, 0xf4, 0xa7, 0xa3, 0x21, 0x0b, 0x4f, 0x27, 0xfd, 0x79, 0x20, 0xad, 0x94,
	0x52, 0x2b, 0x52, 0x6a, 0x45, 0x4a, 0xad, 0xe5, 0xff, 0x35, 0xfa, 0x15, 0xb9, 0x7d, 0xf4, 0x5f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xdf, 0xb0, 0xb1, 0x8b, 0xf0, 0x10, 0x00, 0x00,
}
