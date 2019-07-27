


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
	PackageId            string   `protobuf:"bytes,1,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	Label                string   `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
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




type QueryInstalledChaincodesArgs struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryInstalledChaincodesArgs) Reset()         { *m = QueryInstalledChaincodesArgs{} }
func (m *QueryInstalledChaincodesArgs) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodesArgs) ProtoMessage()    {}
func (*QueryInstalledChaincodesArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{4}
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
	return fileDescriptor_6625a5b20951add3, []int{5}
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
	return fileDescriptor_6625a5b20951add3, []int{5, 0}
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
	return fileDescriptor_6625a5b20951add3, []int{5, 1}
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
	return fileDescriptor_6625a5b20951add3, []int{5, 2}
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
	return fileDescriptor_6625a5b20951add3, []int{6}
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
	return fileDescriptor_6625a5b20951add3, []int{7}
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
	return fileDescriptor_6625a5b20951add3, []int{7, 0}
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
	return fileDescriptor_6625a5b20951add3, []int{7, 1}
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
	return fileDescriptor_6625a5b20951add3, []int{8}
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
	return fileDescriptor_6625a5b20951add3, []int{9}
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
	return fileDescriptor_6625a5b20951add3, []int{10}
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



type SimulateCommitChaincodeDefinitionArgs struct {
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

func (m *SimulateCommitChaincodeDefinitionArgs) Reset()         { *m = SimulateCommitChaincodeDefinitionArgs{} }
func (m *SimulateCommitChaincodeDefinitionArgs) String() string { return proto.CompactTextString(m) }
func (*SimulateCommitChaincodeDefinitionArgs) ProtoMessage()    {}
func (*SimulateCommitChaincodeDefinitionArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{11}
}

func (m *SimulateCommitChaincodeDefinitionArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SimulateCommitChaincodeDefinitionArgs.Unmarshal(m, b)
}
func (m *SimulateCommitChaincodeDefinitionArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SimulateCommitChaincodeDefinitionArgs.Marshal(b, m, deterministic)
}
func (m *SimulateCommitChaincodeDefinitionArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimulateCommitChaincodeDefinitionArgs.Merge(m, src)
}
func (m *SimulateCommitChaincodeDefinitionArgs) XXX_Size() int {
	return xxx_messageInfo_SimulateCommitChaincodeDefinitionArgs.Size(m)
}
func (m *SimulateCommitChaincodeDefinitionArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_SimulateCommitChaincodeDefinitionArgs.DiscardUnknown(m)
}

var xxx_messageInfo_SimulateCommitChaincodeDefinitionArgs proto.InternalMessageInfo

func (m *SimulateCommitChaincodeDefinitionArgs) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *SimulateCommitChaincodeDefinitionArgs) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SimulateCommitChaincodeDefinitionArgs) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *SimulateCommitChaincodeDefinitionArgs) GetEndorsementPlugin() string {
	if m != nil {
		return m.EndorsementPlugin
	}
	return ""
}

func (m *SimulateCommitChaincodeDefinitionArgs) GetValidationPlugin() string {
	if m != nil {
		return m.ValidationPlugin
	}
	return ""
}

func (m *SimulateCommitChaincodeDefinitionArgs) GetValidationParameter() []byte {
	if m != nil {
		return m.ValidationParameter
	}
	return nil
}

func (m *SimulateCommitChaincodeDefinitionArgs) GetCollections() *common.CollectionConfigPackage {
	if m != nil {
		return m.Collections
	}
	return nil
}

func (m *SimulateCommitChaincodeDefinitionArgs) GetInitRequired() bool {
	if m != nil {
		return m.InitRequired
	}
	return false
}




type SimulateCommitChaincodeDefinitionResult struct {
	Approved             map[string]bool `protobuf:"bytes,1,rep,name=approved,proto3" json:"approved,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SimulateCommitChaincodeDefinitionResult) Reset() {
	*m = SimulateCommitChaincodeDefinitionResult{}
}
func (m *SimulateCommitChaincodeDefinitionResult) String() string { return proto.CompactTextString(m) }
func (*SimulateCommitChaincodeDefinitionResult) ProtoMessage()    {}
func (*SimulateCommitChaincodeDefinitionResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{12}
}

func (m *SimulateCommitChaincodeDefinitionResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SimulateCommitChaincodeDefinitionResult.Unmarshal(m, b)
}
func (m *SimulateCommitChaincodeDefinitionResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SimulateCommitChaincodeDefinitionResult.Marshal(b, m, deterministic)
}
func (m *SimulateCommitChaincodeDefinitionResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimulateCommitChaincodeDefinitionResult.Merge(m, src)
}
func (m *SimulateCommitChaincodeDefinitionResult) XXX_Size() int {
	return xxx_messageInfo_SimulateCommitChaincodeDefinitionResult.Size(m)
}
func (m *SimulateCommitChaincodeDefinitionResult) XXX_DiscardUnknown() {
	xxx_messageInfo_SimulateCommitChaincodeDefinitionResult.DiscardUnknown(m)
}

var xxx_messageInfo_SimulateCommitChaincodeDefinitionResult proto.InternalMessageInfo

func (m *SimulateCommitChaincodeDefinitionResult) GetApproved() map[string]bool {
	if m != nil {
		return m.Approved
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
	return fileDescriptor_6625a5b20951add3, []int{13}
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
	Approved             map[string]bool                 `protobuf:"bytes,8,rep,name=approved,proto3" json:"approved,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *QueryChaincodeDefinitionResult) Reset()         { *m = QueryChaincodeDefinitionResult{} }
func (m *QueryChaincodeDefinitionResult) String() string { return proto.CompactTextString(m) }
func (*QueryChaincodeDefinitionResult) ProtoMessage()    {}
func (*QueryChaincodeDefinitionResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{14}
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

func (m *QueryChaincodeDefinitionResult) GetApproved() map[string]bool {
	if m != nil {
		return m.Approved
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
	return fileDescriptor_6625a5b20951add3, []int{15}
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
	return fileDescriptor_6625a5b20951add3, []int{16}
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
	return fileDescriptor_6625a5b20951add3, []int{16, 0}
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



type QueryNamespaceDefinitionsArgs struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryNamespaceDefinitionsArgs) Reset()         { *m = QueryNamespaceDefinitionsArgs{} }
func (m *QueryNamespaceDefinitionsArgs) String() string { return proto.CompactTextString(m) }
func (*QueryNamespaceDefinitionsArgs) ProtoMessage()    {}
func (*QueryNamespaceDefinitionsArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{17}
}

func (m *QueryNamespaceDefinitionsArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryNamespaceDefinitionsArgs.Unmarshal(m, b)
}
func (m *QueryNamespaceDefinitionsArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryNamespaceDefinitionsArgs.Marshal(b, m, deterministic)
}
func (m *QueryNamespaceDefinitionsArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryNamespaceDefinitionsArgs.Merge(m, src)
}
func (m *QueryNamespaceDefinitionsArgs) XXX_Size() int {
	return xxx_messageInfo_QueryNamespaceDefinitionsArgs.Size(m)
}
func (m *QueryNamespaceDefinitionsArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryNamespaceDefinitionsArgs.DiscardUnknown(m)
}

var xxx_messageInfo_QueryNamespaceDefinitionsArgs proto.InternalMessageInfo



type QueryNamespaceDefinitionsResult struct {
	Namespaces           map[string]*QueryNamespaceDefinitionsResult_Namespace `protobuf:"bytes,1,rep,name=namespaces,proto3" json:"namespaces,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                                              `json:"-"`
	XXX_unrecognized     []byte                                                `json:"-"`
	XXX_sizecache        int32                                                 `json:"-"`
}

func (m *QueryNamespaceDefinitionsResult) Reset()         { *m = QueryNamespaceDefinitionsResult{} }
func (m *QueryNamespaceDefinitionsResult) String() string { return proto.CompactTextString(m) }
func (*QueryNamespaceDefinitionsResult) ProtoMessage()    {}
func (*QueryNamespaceDefinitionsResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{18}
}

func (m *QueryNamespaceDefinitionsResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult.Unmarshal(m, b)
}
func (m *QueryNamespaceDefinitionsResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult.Marshal(b, m, deterministic)
}
func (m *QueryNamespaceDefinitionsResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryNamespaceDefinitionsResult.Merge(m, src)
}
func (m *QueryNamespaceDefinitionsResult) XXX_Size() int {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult.Size(m)
}
func (m *QueryNamespaceDefinitionsResult) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryNamespaceDefinitionsResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryNamespaceDefinitionsResult proto.InternalMessageInfo

func (m *QueryNamespaceDefinitionsResult) GetNamespaces() map[string]*QueryNamespaceDefinitionsResult_Namespace {
	if m != nil {
		return m.Namespaces
	}
	return nil
}

type QueryNamespaceDefinitionsResult_Namespace struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryNamespaceDefinitionsResult_Namespace) Reset() {
	*m = QueryNamespaceDefinitionsResult_Namespace{}
}
func (m *QueryNamespaceDefinitionsResult_Namespace) String() string { return proto.CompactTextString(m) }
func (*QueryNamespaceDefinitionsResult_Namespace) ProtoMessage()    {}
func (*QueryNamespaceDefinitionsResult_Namespace) Descriptor() ([]byte, []int) {
	return fileDescriptor_6625a5b20951add3, []int{18, 0}
}

func (m *QueryNamespaceDefinitionsResult_Namespace) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace.Unmarshal(m, b)
}
func (m *QueryNamespaceDefinitionsResult_Namespace) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace.Marshal(b, m, deterministic)
}
func (m *QueryNamespaceDefinitionsResult_Namespace) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace.Merge(m, src)
}
func (m *QueryNamespaceDefinitionsResult_Namespace) XXX_Size() int {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace.Size(m)
}
func (m *QueryNamespaceDefinitionsResult_Namespace) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace.DiscardUnknown(m)
}

var xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace proto.InternalMessageInfo

func (m *QueryNamespaceDefinitionsResult_Namespace) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func init() {
	proto.RegisterType((*InstallChaincodeArgs)(nil), "lifecycle.InstallChaincodeArgs")
	proto.RegisterType((*InstallChaincodeResult)(nil), "lifecycle.InstallChaincodeResult")
	proto.RegisterType((*QueryInstalledChaincodeArgs)(nil), "lifecycle.QueryInstalledChaincodeArgs")
	proto.RegisterType((*QueryInstalledChaincodeResult)(nil), "lifecycle.QueryInstalledChaincodeResult")
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
	proto.RegisterType((*SimulateCommitChaincodeDefinitionArgs)(nil), "lifecycle.SimulateCommitChaincodeDefinitionArgs")
	proto.RegisterType((*SimulateCommitChaincodeDefinitionResult)(nil), "lifecycle.SimulateCommitChaincodeDefinitionResult")
	proto.RegisterMapType((map[string]bool)(nil), "lifecycle.SimulateCommitChaincodeDefinitionResult.ApprovedEntry")
	proto.RegisterType((*QueryChaincodeDefinitionArgs)(nil), "lifecycle.QueryChaincodeDefinitionArgs")
	proto.RegisterType((*QueryChaincodeDefinitionResult)(nil), "lifecycle.QueryChaincodeDefinitionResult")
	proto.RegisterMapType((map[string]bool)(nil), "lifecycle.QueryChaincodeDefinitionResult.ApprovedEntry")
	proto.RegisterType((*QueryChaincodeDefinitionsArgs)(nil), "lifecycle.QueryChaincodeDefinitionsArgs")
	proto.RegisterType((*QueryChaincodeDefinitionsResult)(nil), "lifecycle.QueryChaincodeDefinitionsResult")
	proto.RegisterType((*QueryChaincodeDefinitionsResult_ChaincodeDefinition)(nil), "lifecycle.QueryChaincodeDefinitionsResult.ChaincodeDefinition")
	proto.RegisterType((*QueryNamespaceDefinitionsArgs)(nil), "lifecycle.QueryNamespaceDefinitionsArgs")
	proto.RegisterType((*QueryNamespaceDefinitionsResult)(nil), "lifecycle.QueryNamespaceDefinitionsResult")
	proto.RegisterMapType((map[string]*QueryNamespaceDefinitionsResult_Namespace)(nil), "lifecycle.QueryNamespaceDefinitionsResult.NamespacesEntry")
	proto.RegisterType((*QueryNamespaceDefinitionsResult_Namespace)(nil), "lifecycle.QueryNamespaceDefinitionsResult.Namespace")
}

func init() { proto.RegisterFile("peer/lifecycle/lifecycle.proto", fileDescriptor_6625a5b20951add3) }

var fileDescriptor_6625a5b20951add3 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x58, 0xdd, 0x8e, 0xdb, 0x44,
	0x14, 0x6e, 0x7e, 0x9b, 0x9c, 0xec, 0xaa, 0xed, 0x6c, 0xa0, 0xc6, 0xb0, 0x9b, 0xc5, 0x88, 0x65,
	0xc5, 0x8f, 0x23, 0xb2, 0x48, 0xc0, 0x52, 0x21, 0xb6, 0xcb, 0x4f, 0x5b, 0xb5, 0x50, 0xbc, 0xe5,
	0xa6, 0x42, 0x8a, 0x26, 0xf6, 0x49, 0x76, 0xd4, 0x89, 0xed, 0x8e, 0xed, 0x48, 0x79, 0x0e, 0xae,
	0x79, 0x03, 0x1e, 0x83, 0x67, 0xe0, 0x12, 0xb8, 0x41, 0xe2, 0x8e, 0x57, 0x40, 0x1e, 0x4f, 0x6c,
	0x6f, 0x62, 0xef, 0xa6, 0x3f, 0x70, 0xb5, 0x77, 0xce, 0x9c, 0xef, 0x7c, 0xe7, 0xcc, 0x9c, 0x6f,
	0x66, 0xce, 0x04, 0x76, 0x7c, 0x44, 0xd1, 0xe7, 0x6c, 0x8c, 0xf6, 0xdc, 0xe6, 0x98, 0x7d, 0x99,
	0xbe, 0xf0, 0x42, 0x8f, 0xb4, 0xd3, 0x01, 0xfd, 0xa6, 0xed, 0x4d, 0xa7, 0x9e, 0xdb, 0xb7, 0x3d,
	0xce, 0xd1, 0x0e, 0x99, 0xe7, 0x26, 0x18, 0xc3, 0x82, 0xee, 0x5d, 0x37, 0x08, 0x29, 0xe7, 0xc7,
	0xa7, 0x94, 0xb9, 0xb6, 0xe7, 0xe0, 0x91, 0x98, 0x04, 0xe4, 0x10, 0x5e, 0xb3, 0x17, 0x03, 0x43,
	0x96, 0x20, 0x86, 0x3e, 0xb5, 0x9f, 0xd0, 0x09, 0x6a, 0x95, 0xdd, 0xca, 0xfe, 0x86, 0x75, 0x33,
	0x05, 0x28, 0x86, 0x87, 0x89, 0xd9, 0x78, 0x00, 0xaf, 0x2e, 0x73, 0x5a, 0x18, 0x44, 0x3c, 0x24,
	0xdb, 0x00, 0x8a, 0x63, 0xc8, 0x1c, 0x49, 0xd3, 0xb6, 0xda, 0x6a, 0xe4, 0xae, 0x43, 0xba, 0xd0,
	0xe0, 0x74, 0x84, 0x5c, 0xab, 0x4a, 0x4b, 0xf2, 0xc3, 0xb8, 0x05, 0xaf, 0x7f, 0x1f, 0xa1, 0x98,
	0x2b, 0x4e, 0x74, 0xce, 0x66, 0x7a, 0x3e, 0xa7, 0xf1, 0x08, 0xb6, 0x4b, 0xbc, 0x5f, 0x24, 0xa7,
	0x1d, 0x78, 0xa3, 0x84, 0x35, 0x88, 0x93, 0x32, 0x7e, 0xaf, 0xc3, 0x4e, 0x19, 0x40, 0xc5, 0xf5,
	0xa0, 0xcb, 0x16, 0xc6, 0x61, 0xba, 0x94, 0x81, 0x56, 0xd9, 0xad, 0xed, 0x77, 0x06, 0xb7, 0xcc,
	0xac, 0x9a, 0xe7, 0x13, 0x99, 0x05, 0x33, 0xdb, 0x62, 0xab, 0x68, 0xfd, 0x97, 0x2a, 0x90, 0x55,
	0xec, 0x73, 0xcd, 0x9f, 0x70, 0x00, 0x81, 0x63, 0x14, 0xe8, 0xda, 0x18, 0x68, 0x35, 0x99, 0xf2,
	0xfd, 0x17, 0x49, 0xd9, 0xb4, 0x52, 0xba, 0xaf, 0xdc, 0x50, 0xcc, 0xad, 0x1c, 0xbf, 0x1e, 0xc0,
	0xb5, 0x25, 0x33, 0xb9, 0x0e, 0xb5, 0x27, 0x38, 0x57, 0xe9, 0xc6, 0x9f, 0xe4, 0x1e, 0x34, 0x66,
	0x94, 0x47, 0x28, 0x13, 0xed, 0x0c, 0x3e, 0x5a, 0x3f, 0x9b, 0x8c, 0xdb, 0x4a, 0x28, 0x0e, 0xab,
	0x9f, 0x54, 0x74, 0x0a, 0x90, 0x19, 0xc8, 0x09, 0xc0, 0x4a, 0x8d, 0x0e, 0xd6, 0x0f, 0x91, 0x95,
	0x26, 0x47, 0xa3, 0x7f, 0x0a, 0xed, 0xac, 0x0e, 0x04, 0xea, 0x2e, 0x9d, 0xa2, 0x9a, 0x92, 0xfc,
	0x26, 0x1a, 0x5c, 0x9d, 0xa1, 0x08, 0x98, 0xe7, 0xaa, 0xe5, 0x5f, 0xfc, 0x34, 0x7e, 0xae, 0xc1,
	0xde, 0x91, 0xef, 0x0b, 0x6f, 0x86, 0x29, 0xc5, 0x97, 0x38, 0x66, 0x2e, 0x8b, 0x37, 0xf7, 0xd7,
	0x9e, 0x78, 0x30, 0xff, 0x4e, 0x4c, 0xe4, 0x06, 0xd1, 0xa1, 0x15, 0xe0, 0xd3, 0x28, 0x9e, 0x87,
	0x24, 0xaf, 0x59, 0xe9, 0xef, 0x34, 0x68, 0xb5, 0x38, 0x68, 0xed, 0x4c, 0x50, 0xf2, 0x01, 0x10,
	0x74, 0x1d, 0x4f, 0x04, 0x38, 0x45, 0x37, 0x1c, 0xfa, 0x3c, 0x9a, 0x30, 0x57, 0xab, 0x4b, 0xd0,
	0x8d, 0x9c, 0xe5, 0xa1, 0x34, 0x90, 0xf7, 0xe0, 0xc6, 0x8c, 0x72, 0xe6, 0xd0, 0x38, 0xa5, 0x05,
	0xba, 0x21, 0xd1, 0xd7, 0x33, 0x83, 0x02, 0x7f, 0x08, 0xdd, 0x3c, 0x98, 0x0a, 0x3a, 0xc5, 0x10,
	0x85, 0xd6, 0x94, 0x67, 0xcd, 0x56, 0x0e, 0xbf, 0x30, 0x91, 0x23, 0xe8, 0x64, 0xe7, 0x59, 0xa0,
	0x5d, 0x95, 0x75, 0xef, 0x99, 0xc9, 0x51, 0x67, 0x1e, 0xa7, 0xa6, 0x63, 0xcf, 0x1d, 0xb3, 0x89,
	0x3a, 0x9d, 0xac, 0xbc, 0x0f, 0x79, 0x0b, 0x36, 0xe3, 0x25, 0x1b, 0x0a, 0x7c, 0x1a, 0x31, 0x81,
	0x8e, 0xd6, 0xda, 0xad, 0xec, 0xb7, 0xac, 0x8d, 0x78, 0xd0, 0x52, 0x63, 0x64, 0x00, 0xcd, 0xc0,
	0x8b, 0x84, 0x8d, 0x5a, 0x5b, 0x86, 0xd0, 0x73, 0x75, 0x4f, 0x17, 0xff, 0x44, 0x22, 0x2c, 0x85,
	0x34, 0xfe, 0xaa, 0xc0, 0xb5, 0x25, 0x1b, 0xb9, 0x07, 0x9d, 0xc8, 0xa5, 0x33, 0xca, 0x38, 0x1d,
	0xf1, 0xa4, 0x16, 0x9d, 0xc1, 0x5e, 0x39, 0x99, 0xf9, 0x43, 0x86, 0xbe, 0x73, 0xc5, 0xca, 0x3b,
	0x93, 0x6f, 0x60, 0x93, 0x7b, 0x36, 0xcd, 0xce, 0xe4, 0x44, 0xf5, 0xbb, 0xe7, 0xb0, 0xdd, 0x8f,
	0xf1, 0x77, 0xae, 0x58, 0x1b, 0xd2, 0x51, 0x2d, 0x87, 0xbe, 0x09, 0x9d, 0x5c, 0x18, 0x7d, 0x0f,
	0x1a, 0x12, 0x77, 0xc1, 0xb1, 0x70, 0xbb, 0x09, 0xf5, 0x47, 0x73, 0x1f, 0x8d, 0x77, 0x61, 0xff,
	0x62, 0x19, 0x26, 0x9b, 0xc0, 0xf8, 0xb3, 0x0a, 0xdb, 0xc7, 0xde, 0x74, 0xca, 0xc2, 0x02, 0xec,
	0xa5, 0x54, 0x5f, 0x82, 0x54, 0x8d, 0x37, 0xa1, 0x57, 0xba, 0xc2, 0xaa, 0x0a, 0x7f, 0x57, 0xe1,
	0xed, 0x13, 0x36, 0x8d, 0x38, 0x0d, 0xf1, 0xb2, 0x1a, 0xff, 0x69, 0x35, 0x7e, 0xad, 0xc0, 0x3b,
	0x17, 0x2e, 0xb5, 0x6a, 0x07, 0x7e, 0x84, 0x16, 0x4d, 0x36, 0x92, 0xa3, 0xae, 0x97, 0x2f, 0x72,
	0x7b, 0x79, 0x4d, 0x16, 0x53, 0xed, 0x45, 0x27, 0xb9, 0x43, 0x53, 0x46, 0xfd, 0x33, 0xd8, 0x3c,
	0x63, 0x2a, 0xb8, 0x3f, 0xbb, 0xf9, 0xfb, 0xb3, 0x95, 0xbb, 0x09, 0x8d, 0x81, 0x6a, 0x76, 0xca,
	0x74, 0x52, 0x70, 0x73, 0x19, 0x7f, 0xd4, 0x54, 0x03, 0x54, 0x3e, 0xe3, 0xf3, 0xe4, 0x55, 0x7a,
	0xf1, 0x95, 0x48, 0xa9, 0xf6, 0x4c, 0x52, 0xaa, 0x3f, 0xa3, 0x94, 0x1a, 0x6b, 0x4b, 0xa9, 0xf9,
	0x32, 0xa4, 0x74, 0xb5, 0xe0, 0x0e, 0x3a, 0xc9, 0xc9, 0xa3, 0x25, 0xe5, 0xf1, 0xf1, 0x72, 0xf7,
	0xf1, 0x3f, 0xab, 0xa2, 0xa7, 0x1a, 0xeb, 0x82, 0xb0, 0x49, 0x0f, 0xfc, 0x4f, 0x0d, 0x7a, 0xa5,
	0x08, 0xa5, 0x81, 0x00, 0x5e, 0xc9, 0x9e, 0x19, 0x4e, 0x66, 0x56, 0x5b, 0xe0, 0xf3, 0x35, 0xe6,
	0xb8, 0xd2, 0x62, 0xe5, 0xa6, 0xdf, 0xb5, 0x0b, 0xf0, 0xfa, 0x6f, 0x55, 0xd8, 0x2a, 0x40, 0x17,
	0x76, 0x60, 0x79, 0x91, 0x56, 0xcb, 0x45, 0x7a, 0x79, 0xde, 0xc5, 0xe7, 0xdd, 0x42, 0x12, 0xdf,
	0xd2, 0x29, 0x06, 0x3e, 0xb5, 0x57, 0x24, 0xf1, 0x53, 0x55, 0x49, 0xa2, 0x08, 0xa1, 0x24, 0xf1,
	0x18, 0xc0, 0x5d, 0x58, 0x17, 0x3a, 0x38, 0x5c, 0xd6, 0x41, 0xb9, 0xbf, 0x99, 0x9a, 0x16, 0x0f,
	0x89, 0x8c, 0x4d, 0xef, 0x41, 0x3b, 0x35, 0xc7, 0xe5, 0x0e, 0xe7, 0x7e, 0x5a, 0xee, 0xf8, 0x3b,
	0x7e, 0x69, 0x2c, 0xf9, 0x3f, 0xc7, 0x4b, 0x63, 0x9d, 0xe4, 0x72, 0x3b, 0xe9, 0xb6, 0x0d, 0xef,
	0x7b, 0x62, 0x62, 0x9e, 0xce, 0x7d, 0x14, 0x1c, 0x9d, 0x09, 0x0a, 0x73, 0x4c, 0x47, 0x82, 0xd9,
	0xc9, 0x1b, 0x3d, 0x30, 0xe3, 0x77, 0x7e, 0x16, 0xe4, 0xf1, 0xc1, 0x84, 0x85, 0xa7, 0xd1, 0x28,
	0xae, 0x5f, 0x3f, 0xe7, 0xd4, 0x4f, 0x9c, 0xfa, 0x89, 0x53, 0xff, 0xec, 0x9f, 0x03, 0xa3, 0xa6,
	0x1c, 0x3e, 0xf8, 0x37, 0x00, 0x00, 0xff, 0xff, 0xc0, 0xf2, 0xc4, 0xdb, 0x35, 0x10, 0x00, 0x00,
}
