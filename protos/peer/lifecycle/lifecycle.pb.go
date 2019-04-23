


package lifecycle 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "github.com/mcc-github/blockchain/protos/common"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{0}
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{1}
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{2}
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{3}
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{4}
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
	InstalledChaincodes  []*QueryInstalledChaincodesResult_InstalledChaincode `protobuf:"bytes,1,rep,name=installed_chaincodes,json=installedChaincodes,proto3" json:"installed_chaincodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                             `json:"-"`
	XXX_unrecognized     []byte                                               `json:"-"`
	XXX_sizecache        int32                                                `json:"-"`
}

func (m *QueryInstalledChaincodesResult) Reset()         { *m = QueryInstalledChaincodesResult{} }
func (m *QueryInstalledChaincodesResult) String() string { return proto.CompactTextString(m) }
func (*QueryInstalledChaincodesResult) ProtoMessage()    {}
func (*QueryInstalledChaincodesResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{5}
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
	PackageId            string   `protobuf:"bytes,1,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	Label                string   `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{5, 0}
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{6}
}
func (m *ApproveChaincodeDefinitionForMyOrgArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs.Unmarshal(m, b)
}
func (m *ApproveChaincodeDefinitionForMyOrgArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs.Marshal(b, m, deterministic)
}
func (dst *ApproveChaincodeDefinitionForMyOrgArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgArgs.Merge(dst, src)
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{7}
}
func (m *ChaincodeSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeSource.Unmarshal(m, b)
}
func (m *ChaincodeSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeSource.Marshal(b, m, deterministic)
}
func (dst *ChaincodeSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeSource.Merge(dst, src)
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


func (*ChaincodeSource) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ChaincodeSource_OneofMarshaler, _ChaincodeSource_OneofUnmarshaler, _ChaincodeSource_OneofSizer, []interface{}{
		(*ChaincodeSource_Unavailable_)(nil),
		(*ChaincodeSource_LocalPackage)(nil),
	}
}

func _ChaincodeSource_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ChaincodeSource)
	
	switch x := m.Type.(type) {
	case *ChaincodeSource_Unavailable_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Unavailable); err != nil {
			return err
		}
	case *ChaincodeSource_LocalPackage:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LocalPackage); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ChaincodeSource.Type has unexpected type %T", x)
	}
	return nil
}

func _ChaincodeSource_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ChaincodeSource)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeSource_Unavailable)
		err := b.DecodeMessage(msg)
		m.Type = &ChaincodeSource_Unavailable_{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeSource_Local)
		err := b.DecodeMessage(msg)
		m.Type = &ChaincodeSource_LocalPackage{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ChaincodeSource_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ChaincodeSource)
	
	switch x := m.Type.(type) {
	case *ChaincodeSource_Unavailable_:
		s := proto.Size(x.Unavailable)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ChaincodeSource_LocalPackage:
		s := proto.Size(x.LocalPackage)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{7, 0}
}
func (m *ChaincodeSource_Unavailable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeSource_Unavailable.Unmarshal(m, b)
}
func (m *ChaincodeSource_Unavailable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeSource_Unavailable.Marshal(b, m, deterministic)
}
func (dst *ChaincodeSource_Unavailable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeSource_Unavailable.Merge(dst, src)
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{7, 1}
}
func (m *ChaincodeSource_Local) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeSource_Local.Unmarshal(m, b)
}
func (m *ChaincodeSource_Local) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeSource_Local.Marshal(b, m, deterministic)
}
func (dst *ChaincodeSource_Local) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeSource_Local.Merge(dst, src)
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{8}
}
func (m *ApproveChaincodeDefinitionForMyOrgResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult.Unmarshal(m, b)
}
func (m *ApproveChaincodeDefinitionForMyOrgResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult.Marshal(b, m, deterministic)
}
func (dst *ApproveChaincodeDefinitionForMyOrgResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApproveChaincodeDefinitionForMyOrgResult.Merge(dst, src)
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{9}
}
func (m *CommitChaincodeDefinitionArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommitChaincodeDefinitionArgs.Unmarshal(m, b)
}
func (m *CommitChaincodeDefinitionArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommitChaincodeDefinitionArgs.Marshal(b, m, deterministic)
}
func (dst *CommitChaincodeDefinitionArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitChaincodeDefinitionArgs.Merge(dst, src)
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{10}
}
func (m *CommitChaincodeDefinitionResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommitChaincodeDefinitionResult.Unmarshal(m, b)
}
func (m *CommitChaincodeDefinitionResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommitChaincodeDefinitionResult.Marshal(b, m, deterministic)
}
func (dst *CommitChaincodeDefinitionResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitChaincodeDefinitionResult.Merge(dst, src)
}
func (m *CommitChaincodeDefinitionResult) XXX_Size() int {
	return xxx_messageInfo_CommitChaincodeDefinitionResult.Size(m)
}
func (m *CommitChaincodeDefinitionResult) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitChaincodeDefinitionResult.DiscardUnknown(m)
}

var xxx_messageInfo_CommitChaincodeDefinitionResult proto.InternalMessageInfo



type QueryApprovalStatusArgs struct {
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

func (m *QueryApprovalStatusArgs) Reset()         { *m = QueryApprovalStatusArgs{} }
func (m *QueryApprovalStatusArgs) String() string { return proto.CompactTextString(m) }
func (*QueryApprovalStatusArgs) ProtoMessage()    {}
func (*QueryApprovalStatusArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{11}
}
func (m *QueryApprovalStatusArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryApprovalStatusArgs.Unmarshal(m, b)
}
func (m *QueryApprovalStatusArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryApprovalStatusArgs.Marshal(b, m, deterministic)
}
func (dst *QueryApprovalStatusArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryApprovalStatusArgs.Merge(dst, src)
}
func (m *QueryApprovalStatusArgs) XXX_Size() int {
	return xxx_messageInfo_QueryApprovalStatusArgs.Size(m)
}
func (m *QueryApprovalStatusArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryApprovalStatusArgs.DiscardUnknown(m)
}

var xxx_messageInfo_QueryApprovalStatusArgs proto.InternalMessageInfo

func (m *QueryApprovalStatusArgs) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *QueryApprovalStatusArgs) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *QueryApprovalStatusArgs) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *QueryApprovalStatusArgs) GetEndorsementPlugin() string {
	if m != nil {
		return m.EndorsementPlugin
	}
	return ""
}

func (m *QueryApprovalStatusArgs) GetValidationPlugin() string {
	if m != nil {
		return m.ValidationPlugin
	}
	return ""
}

func (m *QueryApprovalStatusArgs) GetValidationParameter() []byte {
	if m != nil {
		return m.ValidationParameter
	}
	return nil
}

func (m *QueryApprovalStatusArgs) GetCollections() *common.CollectionConfigPackage {
	if m != nil {
		return m.Collections
	}
	return nil
}

func (m *QueryApprovalStatusArgs) GetInitRequired() bool {
	if m != nil {
		return m.InitRequired
	}
	return false
}





type QueryApprovalStatusResults struct {
	Approved             map[string]bool `protobuf:"bytes,1,rep,name=approved,proto3" json:"approved,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *QueryApprovalStatusResults) Reset()         { *m = QueryApprovalStatusResults{} }
func (m *QueryApprovalStatusResults) String() string { return proto.CompactTextString(m) }
func (*QueryApprovalStatusResults) ProtoMessage()    {}
func (*QueryApprovalStatusResults) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{12}
}
func (m *QueryApprovalStatusResults) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryApprovalStatusResults.Unmarshal(m, b)
}
func (m *QueryApprovalStatusResults) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryApprovalStatusResults.Marshal(b, m, deterministic)
}
func (dst *QueryApprovalStatusResults) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryApprovalStatusResults.Merge(dst, src)
}
func (m *QueryApprovalStatusResults) XXX_Size() int {
	return xxx_messageInfo_QueryApprovalStatusResults.Size(m)
}
func (m *QueryApprovalStatusResults) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryApprovalStatusResults.DiscardUnknown(m)
}

var xxx_messageInfo_QueryApprovalStatusResults proto.InternalMessageInfo

func (m *QueryApprovalStatusResults) GetApproved() map[string]bool {
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{13}
}
func (m *QueryChaincodeDefinitionArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryChaincodeDefinitionArgs.Unmarshal(m, b)
}
func (m *QueryChaincodeDefinitionArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryChaincodeDefinitionArgs.Marshal(b, m, deterministic)
}
func (dst *QueryChaincodeDefinitionArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryChaincodeDefinitionArgs.Merge(dst, src)
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
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *QueryChaincodeDefinitionResult) Reset()         { *m = QueryChaincodeDefinitionResult{} }
func (m *QueryChaincodeDefinitionResult) String() string { return proto.CompactTextString(m) }
func (*QueryChaincodeDefinitionResult) ProtoMessage()    {}
func (*QueryChaincodeDefinitionResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{14}
}
func (m *QueryChaincodeDefinitionResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryChaincodeDefinitionResult.Unmarshal(m, b)
}
func (m *QueryChaincodeDefinitionResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryChaincodeDefinitionResult.Marshal(b, m, deterministic)
}
func (dst *QueryChaincodeDefinitionResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryChaincodeDefinitionResult.Merge(dst, src)
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



type QueryNamespaceDefinitionsArgs struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryNamespaceDefinitionsArgs) Reset()         { *m = QueryNamespaceDefinitionsArgs{} }
func (m *QueryNamespaceDefinitionsArgs) String() string { return proto.CompactTextString(m) }
func (*QueryNamespaceDefinitionsArgs) ProtoMessage()    {}
func (*QueryNamespaceDefinitionsArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{15}
}
func (m *QueryNamespaceDefinitionsArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryNamespaceDefinitionsArgs.Unmarshal(m, b)
}
func (m *QueryNamespaceDefinitionsArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryNamespaceDefinitionsArgs.Marshal(b, m, deterministic)
}
func (dst *QueryNamespaceDefinitionsArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryNamespaceDefinitionsArgs.Merge(dst, src)
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{16}
}
func (m *QueryNamespaceDefinitionsResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult.Unmarshal(m, b)
}
func (m *QueryNamespaceDefinitionsResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult.Marshal(b, m, deterministic)
}
func (dst *QueryNamespaceDefinitionsResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryNamespaceDefinitionsResult.Merge(dst, src)
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
	return fileDescriptor_lifecycle_ddf322ea8d928d52, []int{16, 0}
}
func (m *QueryNamespaceDefinitionsResult_Namespace) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace.Unmarshal(m, b)
}
func (m *QueryNamespaceDefinitionsResult_Namespace) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace.Marshal(b, m, deterministic)
}
func (dst *QueryNamespaceDefinitionsResult_Namespace) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryNamespaceDefinitionsResult_Namespace.Merge(dst, src)
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
	proto.RegisterType((*ApproveChaincodeDefinitionForMyOrgArgs)(nil), "lifecycle.ApproveChaincodeDefinitionForMyOrgArgs")
	proto.RegisterType((*ChaincodeSource)(nil), "lifecycle.ChaincodeSource")
	proto.RegisterType((*ChaincodeSource_Unavailable)(nil), "lifecycle.ChaincodeSource.Unavailable")
	proto.RegisterType((*ChaincodeSource_Local)(nil), "lifecycle.ChaincodeSource.Local")
	proto.RegisterType((*ApproveChaincodeDefinitionForMyOrgResult)(nil), "lifecycle.ApproveChaincodeDefinitionForMyOrgResult")
	proto.RegisterType((*CommitChaincodeDefinitionArgs)(nil), "lifecycle.CommitChaincodeDefinitionArgs")
	proto.RegisterType((*CommitChaincodeDefinitionResult)(nil), "lifecycle.CommitChaincodeDefinitionResult")
	proto.RegisterType((*QueryApprovalStatusArgs)(nil), "lifecycle.QueryApprovalStatusArgs")
	proto.RegisterType((*QueryApprovalStatusResults)(nil), "lifecycle.QueryApprovalStatusResults")
	proto.RegisterMapType((map[string]bool)(nil), "lifecycle.QueryApprovalStatusResults.ApprovedEntry")
	proto.RegisterType((*QueryChaincodeDefinitionArgs)(nil), "lifecycle.QueryChaincodeDefinitionArgs")
	proto.RegisterType((*QueryChaincodeDefinitionResult)(nil), "lifecycle.QueryChaincodeDefinitionResult")
	proto.RegisterType((*QueryNamespaceDefinitionsArgs)(nil), "lifecycle.QueryNamespaceDefinitionsArgs")
	proto.RegisterType((*QueryNamespaceDefinitionsResult)(nil), "lifecycle.QueryNamespaceDefinitionsResult")
	proto.RegisterMapType((map[string]*QueryNamespaceDefinitionsResult_Namespace)(nil), "lifecycle.QueryNamespaceDefinitionsResult.NamespacesEntry")
	proto.RegisterType((*QueryNamespaceDefinitionsResult_Namespace)(nil), "lifecycle.QueryNamespaceDefinitionsResult.Namespace")
}

func init() {
	proto.RegisterFile("peer/lifecycle/lifecycle.proto", fileDescriptor_lifecycle_ddf322ea8d928d52)
}

var fileDescriptor_lifecycle_ddf322ea8d928d52 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x56, 0xcb, 0x72, 0xdb, 0x36,
	0x14, 0x0d, 0x25, 0x5b, 0x96, 0xae, 0xec, 0x49, 0x82, 0x78, 0x6a, 0x96, 0xad, 0x2d, 0x95, 0x9d,
	0xf1, 0x68, 0xfa, 0xa0, 0xa6, 0x72, 0x17, 0x1d, 0x37, 0x1b, 0xc7, 0x7d, 0xc4, 0x99, 0xa6, 0x49,
	0x99, 0x74, 0x93, 0x8d, 0x06, 0x22, 0x61, 0x19, 0x13, 0x88, 0x60, 0x00, 0x52, 0x33, 0xfc, 0x8e,
	0xae, 0xfb, 0x03, 0xfd, 0x97, 0xfe, 0x40, 0xa7, 0x8b, 0x6e, 0xfa, 0x1f, 0x1d, 0x02, 0xe0, 0xc3,
	0x8a, 0xe8, 0x3c, 0xea, 0xa5, 0x77, 0x00, 0xce, 0xb9, 0x07, 0xe0, 0xbd, 0xe7, 0x12, 0x80, 0x83,
	0x98, 0x10, 0x31, 0x66, 0xf4, 0x9c, 0x04, 0x59, 0xc0, 0x48, 0x35, 0xf2, 0x62, 0xc1, 0x13, 0x8e,
	0x7a, 0xe5, 0x82, 0xb3, 0x17, 0xf0, 0xc5, 0x82, 0x47, 0xe3, 0x80, 0x33, 0x46, 0x82, 0x84, 0xf2,
	0x48, 0x73, 0x5c, 0x1f, 0x76, 0xcf, 0x22, 0x99, 0x60, 0xc6, 0x4e, 0x2f, 0x30, 0x8d, 0x02, 0x1e,
	0x92, 0x13, 0x31, 0x97, 0xe8, 0x18, 0x3e, 0x0c, 0x8a, 0x85, 0x29, 0xd5, 0x8c, 0x69, 0x8c, 0x83,
	0x97, 0x78, 0x4e, 0x6c, 0x6b, 0x68, 0x8d, 0xb6, 0xfd, 0xbd, 0x92, 0x60, 0x14, 0x9e, 0x6a, 0xd8,
	0x7d, 0x0c, 0x1f, 0xac, 0x6a, 0xfa, 0x44, 0xa6, 0x2c, 0x41, 0xfb, 0x00, 0x46, 0x63, 0x4a, 0x43,
	0x25, 0xd3, 0xf3, 0x7b, 0x66, 0xe5, 0x2c, 0x44, 0xbb, 0xb0, 0xc9, 0xf0, 0x8c, 0x30, 0xbb, 0xa5,
	0x10, 0x3d, 0x71, 0xef, 0xc3, 0x47, 0xbf, 0xa4, 0x44, 0x64, 0x46, 0x93, 0x84, 0x97, 0x4f, 0x7a,
	0xb5, 0xa6, 0xfb, 0x1c, 0xf6, 0x1b, 0xa2, 0xff, 0xcf, 0x99, 0x0e, 0xe0, 0xe3, 0x06, 0x55, 0x99,
	0x1f, 0xca, 0xfd, 0xdb, 0x82, 0x83, 0x26, 0x82, 0xd9, 0x97, 0xc3, 0x2e, 0x2d, 0xc0, 0x69, 0x99,
	0x4a, 0x69, 0x5b, 0xc3, 0xf6, 0xa8, 0x3f, 0xb9, 0xef, 0x55, 0xd5, 0xbc, 0x5a, 0xc8, 0x5b, 0xf3,
	0x65, 0xf7, 0xe8, 0xeb, 0x6c, 0xe7, 0x0c, 0xd0, 0xeb, 0xd4, 0xf7, 0xfb, 0xfc, 0xdf, 0xdb, 0x70,
	0x78, 0x12, 0xc7, 0x82, 0x2f, 0x49, 0xa9, 0xf4, 0x1d, 0x39, 0xa7, 0x11, 0xcd, 0xad, 0xf5, 0x03,
	0x17, 0x8f, 0xb3, 0x27, 0x62, 0xae, 0xca, 0xe3, 0x40, 0x57, 0x92, 0x57, 0x29, 0x89, 0x02, 0xed,
	0x9b, 0xb6, 0x5f, 0xce, 0x11, 0x82, 0x8d, 0x08, 0x2f, 0x88, 0xd1, 0x56, 0x63, 0x64, 0xc3, 0xd6,
	0x92, 0x08, 0x49, 0x79, 0x64, 0xb7, 0xd5, 0x72, 0x31, 0x45, 0x5f, 0x02, 0x22, 0x51, 0xc8, 0x85,
	0x24, 0x0b, 0x12, 0x25, 0xd3, 0x98, 0xa5, 0x73, 0x1a, 0xd9, 0x1b, 0x8a, 0x74, 0xb7, 0x86, 0x3c,
	0x55, 0x00, 0xfa, 0x1c, 0xee, 0x2e, 0x31, 0xa3, 0x21, 0xce, 0x8f, 0x54, 0xb0, 0x37, 0x15, 0xfb,
	0x4e, 0x05, 0x18, 0xf2, 0x57, 0xb0, 0x5b, 0x27, 0x63, 0x81, 0x17, 0x24, 0x21, 0xc2, 0xee, 0x28,
	0xa7, 0xdf, 0xab, 0xf1, 0x0b, 0x08, 0x9d, 0x40, 0xbf, 0xea, 0x26, 0x69, 0x6f, 0x0d, 0xad, 0x51,
	0x7f, 0x32, 0xf0, 0x74, 0xa3, 0x79, 0xa7, 0x25, 0x74, 0xca, 0xa3, 0x73, 0x3a, 0x37, 0xbd, 0xe1,
	0xd7, 0x63, 0xd0, 0xa7, 0xb0, 0x93, 0xa7, 0x6c, 0x2a, 0xc8, 0xab, 0x94, 0x0a, 0x12, 0xda, 0xdd,
	0xa1, 0x35, 0xea, 0xfa, 0xdb, 0xf9, 0xa2, 0x6f, 0xd6, 0xd0, 0x04, 0x3a, 0x92, 0xa7, 0x22, 0x20,
	0x76, 0x4f, 0x6d, 0xe1, 0xd4, 0x9c, 0x51, 0x26, 0xff, 0x99, 0x62, 0xf8, 0x86, 0xe9, 0xfe, 0x6b,
	0xc1, 0xed, 0x15, 0x0c, 0x3d, 0x82, 0x7e, 0x1a, 0xe1, 0x25, 0xa6, 0x0c, 0xcf, 0x98, 0xae, 0x45,
	0x7f, 0x72, 0xd8, 0x2c, 0xe6, 0xfd, 0x5a, 0xb1, 0x1f, 0xde, 0xf2, 0xeb, 0xc1, 0xe8, 0x47, 0xd8,
	0x61, 0x3c, 0xc0, 0xd5, 0x1f, 0xa1, 0xa5, 0xd4, 0x86, 0x57, 0xa8, 0xfd, 0x94, 0xf3, 0x1f, 0xde,
	0xf2, 0xb7, 0x55, 0xa0, 0x49, 0x87, 0xb3, 0x03, 0xfd, 0xda, 0x36, 0xce, 0x21, 0x6c, 0x2a, 0xde,
	0x1b, 0x5c, 0xf9, 0xa0, 0x03, 0x1b, 0xcf, 0xb3, 0x98, 0xb8, 0x9f, 0xc1, 0xe8, 0xcd, 0x36, 0xd4,
	0x6d, 0xe2, 0xfe, 0xd3, 0x82, 0xfd, 0x53, 0xbe, 0x58, 0xd0, 0x64, 0x0d, 0xf7, 0xc6, 0xaa, 0xd7,
	0x60, 0x55, 0xf7, 0x13, 0x18, 0x34, 0x66, 0xd8, 0x54, 0xe1, 0xaf, 0x16, 0xec, 0xa9, 0xff, 0x99,
	0xae, 0x1b, 0x66, 0xcf, 0x12, 0x9c, 0xa4, 0xf2, 0x26, 0xff, 0xd7, 0x91, 0xff, 0x3f, 0x2c, 0x70,
	0xd6, 0x24, 0x57, 0xa7, 0x5e, 0xa2, 0x27, 0xd0, 0xc5, 0xba, 0x5b, 0x42, 0x73, 0xcb, 0x1c, 0xad,
	0xde, 0x32, 0x6b, 0x03, 0x3d, 0xd3, 0x63, 0xe1, 0xf7, 0x51, 0x22, 0x32, 0xbf, 0x14, 0x71, 0xbe,
	0x85, 0x9d, 0x4b, 0x10, 0xba, 0x03, 0xed, 0x97, 0x24, 0x33, 0xfd, 0x9a, 0x0f, 0xf3, 0xfb, 0x63,
	0x89, 0x59, 0xaa, 0x0b, 0xd7, 0xf5, 0xf5, 0xe4, 0xb8, 0xf5, 0x8d, 0xe5, 0x4e, 0xcc, 0x15, 0xda,
	0xd4, 0x8d, 0x45, 0xc5, 0xad, 0xaa, 0xe2, 0xee, 0x9f, 0x2d, 0x73, 0xad, 0x36, 0x1a, 0xec, 0x4a,
	0x13, 0xd5, 0x0c, 0xd3, 0x7a, 0x1b, 0xc3, 0xb4, 0xdf, 0xc9, 0x30, 0x1b, 0xef, 0x68, 0x98, 0xcd,
	0xb7, 0x36, 0x4c, 0xe7, 0x3a, 0x0c, 0xb3, 0xb5, 0xc6, 0x30, 0x03, 0xf3, 0x38, 0xfa, 0x19, 0x2f,
	0x88, 0x8c, 0x71, 0x50, 0x4b, 0xa7, 0x7e, 0xc7, 0xfc, 0xd6, 0x82, 0x41, 0x23, 0xc3, 0x64, 0xfc,
	0x05, 0x40, 0x54, 0xa0, 0xc5, 0xf3, 0xe5, 0x78, 0xd5, 0x58, 0xcd, 0xf1, 0x5e, 0x09, 0x49, 0xed,
	0xaf, 0x9a, 0x9a, 0x33, 0x80, 0x5e, 0x09, 0xe7, 0x8e, 0x48, 0xb2, 0xb8, 0x74, 0x44, 0x3e, 0x76,
	0x24, 0xdc, 0x5e, 0x89, 0x5f, 0x63, 0xc2, 0x47, 0x75, 0x13, 0xf6, 0x27, 0x5f, 0xbf, 0xcf, 0xe1,
	0x6a, 0xd6, 0x7d, 0x10, 0xc0, 0x17, 0x5c, 0xcc, 0xbd, 0x8b, 0x2c, 0x26, 0x82, 0x91, 0x70, 0x4e,
	0x84, 0x77, 0x8e, 0x67, 0x82, 0x06, 0xfa, 0x51, 0x2d, 0xbd, 0xfc, 0x61, 0x5e, 0x6d, 0xf2, 0xe2,
	0x68, 0x4e, 0x93, 0x8b, 0x74, 0x96, 0xd7, 0x6f, 0x5c, 0x0b, 0x1a, 0xeb, 0xa0, 0xb1, 0x0e, 0x1a,
	0x5f, 0x7e, 0xcd, 0xcf, 0x3a, 0x6a, 0xf9, 0xe8, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x70, 0x0b,
	0x39, 0xd0, 0xe6, 0x0b, 0x00, 0x00,
}
