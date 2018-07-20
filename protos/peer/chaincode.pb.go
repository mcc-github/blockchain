


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 


type ConfidentialityLevel int32

const (
	ConfidentialityLevel_PUBLIC       ConfidentialityLevel = 0
	ConfidentialityLevel_CONFIDENTIAL ConfidentialityLevel = 1
)

var ConfidentialityLevel_name = map[int32]string{
	0: "PUBLIC",
	1: "CONFIDENTIAL",
}
var ConfidentialityLevel_value = map[string]int32{
	"PUBLIC":       0,
	"CONFIDENTIAL": 1,
}

func (x ConfidentialityLevel) String() string {
	return proto.EnumName(ConfidentialityLevel_name, int32(x))
}
func (ConfidentialityLevel) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{0}
}

type ChaincodeSpec_Type int32

const (
	ChaincodeSpec_UNDEFINED ChaincodeSpec_Type = 0
	ChaincodeSpec_GOLANG    ChaincodeSpec_Type = 1
	ChaincodeSpec_NODE      ChaincodeSpec_Type = 2
	ChaincodeSpec_CAR       ChaincodeSpec_Type = 3
	ChaincodeSpec_JAVA      ChaincodeSpec_Type = 4
)

var ChaincodeSpec_Type_name = map[int32]string{
	0: "UNDEFINED",
	1: "GOLANG",
	2: "NODE",
	3: "CAR",
	4: "JAVA",
}
var ChaincodeSpec_Type_value = map[string]int32{
	"UNDEFINED": 0,
	"GOLANG":    1,
	"NODE":      2,
	"CAR":       3,
	"JAVA":      4,
}

func (x ChaincodeSpec_Type) String() string {
	return proto.EnumName(ChaincodeSpec_Type_name, int32(x))
}
func (ChaincodeSpec_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{2, 0}
}

type ChaincodeDeploymentSpec_ExecutionEnvironment int32

const (
	ChaincodeDeploymentSpec_DOCKER ChaincodeDeploymentSpec_ExecutionEnvironment = 0
	ChaincodeDeploymentSpec_SYSTEM ChaincodeDeploymentSpec_ExecutionEnvironment = 1
)

var ChaincodeDeploymentSpec_ExecutionEnvironment_name = map[int32]string{
	0: "DOCKER",
	1: "SYSTEM",
}
var ChaincodeDeploymentSpec_ExecutionEnvironment_value = map[string]int32{
	"DOCKER": 0,
	"SYSTEM": 1,
}

func (x ChaincodeDeploymentSpec_ExecutionEnvironment) String() string {
	return proto.EnumName(ChaincodeDeploymentSpec_ExecutionEnvironment_name, int32(x))
}
func (ChaincodeDeploymentSpec_ExecutionEnvironment) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{3, 0}
}








type ChaincodeID struct {
	
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	
	
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	
	Version              string   `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeID) Reset()         { *m = ChaincodeID{} }
func (m *ChaincodeID) String() string { return proto.CompactTextString(m) }
func (*ChaincodeID) ProtoMessage()    {}
func (*ChaincodeID) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{0}
}
func (m *ChaincodeID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeID.Unmarshal(m, b)
}
func (m *ChaincodeID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeID.Marshal(b, m, deterministic)
}
func (dst *ChaincodeID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeID.Merge(dst, src)
}
func (m *ChaincodeID) XXX_Size() int {
	return xxx_messageInfo_ChaincodeID.Size(m)
}
func (m *ChaincodeID) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeID.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeID proto.InternalMessageInfo

func (m *ChaincodeID) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *ChaincodeID) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ChaincodeID) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}




type ChaincodeInput struct {
	Args                 [][]byte          `protobuf:"bytes,1,rep,name=args,proto3" json:"args,omitempty"`
	Decorations          map[string][]byte `protobuf:"bytes,2,rep,name=decorations" json:"decorations,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ChaincodeInput) Reset()         { *m = ChaincodeInput{} }
func (m *ChaincodeInput) String() string { return proto.CompactTextString(m) }
func (*ChaincodeInput) ProtoMessage()    {}
func (*ChaincodeInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{1}
}
func (m *ChaincodeInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeInput.Unmarshal(m, b)
}
func (m *ChaincodeInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeInput.Marshal(b, m, deterministic)
}
func (dst *ChaincodeInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeInput.Merge(dst, src)
}
func (m *ChaincodeInput) XXX_Size() int {
	return xxx_messageInfo_ChaincodeInput.Size(m)
}
func (m *ChaincodeInput) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeInput.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeInput proto.InternalMessageInfo

func (m *ChaincodeInput) GetArgs() [][]byte {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *ChaincodeInput) GetDecorations() map[string][]byte {
	if m != nil {
		return m.Decorations
	}
	return nil
}



type ChaincodeSpec struct {
	Type                 ChaincodeSpec_Type `protobuf:"varint,1,opt,name=type,enum=protos.ChaincodeSpec_Type" json:"type,omitempty"`
	ChaincodeId          *ChaincodeID       `protobuf:"bytes,2,opt,name=chaincode_id,json=chaincodeId" json:"chaincode_id,omitempty"`
	Input                *ChaincodeInput    `protobuf:"bytes,3,opt,name=input" json:"input,omitempty"`
	Timeout              int32              `protobuf:"varint,4,opt,name=timeout" json:"timeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ChaincodeSpec) Reset()         { *m = ChaincodeSpec{} }
func (m *ChaincodeSpec) String() string { return proto.CompactTextString(m) }
func (*ChaincodeSpec) ProtoMessage()    {}
func (*ChaincodeSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{2}
}
func (m *ChaincodeSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeSpec.Unmarshal(m, b)
}
func (m *ChaincodeSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeSpec.Marshal(b, m, deterministic)
}
func (dst *ChaincodeSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeSpec.Merge(dst, src)
}
func (m *ChaincodeSpec) XXX_Size() int {
	return xxx_messageInfo_ChaincodeSpec.Size(m)
}
func (m *ChaincodeSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeSpec.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeSpec proto.InternalMessageInfo

func (m *ChaincodeSpec) GetType() ChaincodeSpec_Type {
	if m != nil {
		return m.Type
	}
	return ChaincodeSpec_UNDEFINED
}

func (m *ChaincodeSpec) GetChaincodeId() *ChaincodeID {
	if m != nil {
		return m.ChaincodeId
	}
	return nil
}

func (m *ChaincodeSpec) GetInput() *ChaincodeInput {
	if m != nil {
		return m.Input
	}
	return nil
}

func (m *ChaincodeSpec) GetTimeout() int32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}



type ChaincodeDeploymentSpec struct {
	ChaincodeSpec        *ChaincodeSpec                               `protobuf:"bytes,1,opt,name=chaincode_spec,json=chaincodeSpec" json:"chaincode_spec,omitempty"`
	CodePackage          []byte                                       `protobuf:"bytes,3,opt,name=code_package,json=codePackage,proto3" json:"code_package,omitempty"`
	ExecEnv              ChaincodeDeploymentSpec_ExecutionEnvironment `protobuf:"varint,4,opt,name=exec_env,json=execEnv,enum=protos.ChaincodeDeploymentSpec_ExecutionEnvironment" json:"exec_env,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                     `json:"-"`
	XXX_unrecognized     []byte                                       `json:"-"`
	XXX_sizecache        int32                                        `json:"-"`
}

func (m *ChaincodeDeploymentSpec) Reset()         { *m = ChaincodeDeploymentSpec{} }
func (m *ChaincodeDeploymentSpec) String() string { return proto.CompactTextString(m) }
func (*ChaincodeDeploymentSpec) ProtoMessage()    {}
func (*ChaincodeDeploymentSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{3}
}
func (m *ChaincodeDeploymentSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeDeploymentSpec.Unmarshal(m, b)
}
func (m *ChaincodeDeploymentSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeDeploymentSpec.Marshal(b, m, deterministic)
}
func (dst *ChaincodeDeploymentSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeDeploymentSpec.Merge(dst, src)
}
func (m *ChaincodeDeploymentSpec) XXX_Size() int {
	return xxx_messageInfo_ChaincodeDeploymentSpec.Size(m)
}
func (m *ChaincodeDeploymentSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeDeploymentSpec.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeDeploymentSpec proto.InternalMessageInfo

func (m *ChaincodeDeploymentSpec) GetChaincodeSpec() *ChaincodeSpec {
	if m != nil {
		return m.ChaincodeSpec
	}
	return nil
}

func (m *ChaincodeDeploymentSpec) GetCodePackage() []byte {
	if m != nil {
		return m.CodePackage
	}
	return nil
}

func (m *ChaincodeDeploymentSpec) GetExecEnv() ChaincodeDeploymentSpec_ExecutionEnvironment {
	if m != nil {
		return m.ExecEnv
	}
	return ChaincodeDeploymentSpec_DOCKER
}


type ChaincodeInvocationSpec struct {
	ChaincodeSpec        *ChaincodeSpec `protobuf:"bytes,1,opt,name=chaincode_spec,json=chaincodeSpec" json:"chaincode_spec,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ChaincodeInvocationSpec) Reset()         { *m = ChaincodeInvocationSpec{} }
func (m *ChaincodeInvocationSpec) String() string { return proto.CompactTextString(m) }
func (*ChaincodeInvocationSpec) ProtoMessage()    {}
func (*ChaincodeInvocationSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{4}
}
func (m *ChaincodeInvocationSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeInvocationSpec.Unmarshal(m, b)
}
func (m *ChaincodeInvocationSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeInvocationSpec.Marshal(b, m, deterministic)
}
func (dst *ChaincodeInvocationSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeInvocationSpec.Merge(dst, src)
}
func (m *ChaincodeInvocationSpec) XXX_Size() int {
	return xxx_messageInfo_ChaincodeInvocationSpec.Size(m)
}
func (m *ChaincodeInvocationSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeInvocationSpec.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeInvocationSpec proto.InternalMessageInfo

func (m *ChaincodeInvocationSpec) GetChaincodeSpec() *ChaincodeSpec {
	if m != nil {
		return m.ChaincodeSpec
	}
	return nil
}


type LifecycleEvent struct {
	ChaincodeName        string   `protobuf:"bytes,1,opt,name=chaincode_name,json=chaincodeName" json:"chaincode_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LifecycleEvent) Reset()         { *m = LifecycleEvent{} }
func (m *LifecycleEvent) String() string { return proto.CompactTextString(m) }
func (*LifecycleEvent) ProtoMessage()    {}
func (*LifecycleEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{5}
}
func (m *LifecycleEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LifecycleEvent.Unmarshal(m, b)
}
func (m *LifecycleEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LifecycleEvent.Marshal(b, m, deterministic)
}
func (dst *LifecycleEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LifecycleEvent.Merge(dst, src)
}
func (m *LifecycleEvent) XXX_Size() int {
	return xxx_messageInfo_LifecycleEvent.Size(m)
}
func (m *LifecycleEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_LifecycleEvent.DiscardUnknown(m)
}

var xxx_messageInfo_LifecycleEvent proto.InternalMessageInfo

func (m *LifecycleEvent) GetChaincodeName() string {
	if m != nil {
		return m.ChaincodeName
	}
	return ""
}


type ChaincodeInstallPackage struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Path                 string   `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	CodePackage          []byte   `protobuf:"bytes,3,opt,name=code_package,json=codePackage,proto3" json:"code_package,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeInstallPackage) Reset()         { *m = ChaincodeInstallPackage{} }
func (m *ChaincodeInstallPackage) String() string { return proto.CompactTextString(m) }
func (*ChaincodeInstallPackage) ProtoMessage()    {}
func (*ChaincodeInstallPackage) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_97aac793bf11f629, []int{6}
}
func (m *ChaincodeInstallPackage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeInstallPackage.Unmarshal(m, b)
}
func (m *ChaincodeInstallPackage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeInstallPackage.Marshal(b, m, deterministic)
}
func (dst *ChaincodeInstallPackage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeInstallPackage.Merge(dst, src)
}
func (m *ChaincodeInstallPackage) XXX_Size() int {
	return xxx_messageInfo_ChaincodeInstallPackage.Size(m)
}
func (m *ChaincodeInstallPackage) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeInstallPackage.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeInstallPackage proto.InternalMessageInfo

func (m *ChaincodeInstallPackage) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ChaincodeInstallPackage) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *ChaincodeInstallPackage) GetCodePackage() []byte {
	if m != nil {
		return m.CodePackage
	}
	return nil
}

func init() {
	proto.RegisterType((*ChaincodeID)(nil), "protos.ChaincodeID")
	proto.RegisterType((*ChaincodeInput)(nil), "protos.ChaincodeInput")
	proto.RegisterMapType((map[string][]byte)(nil), "protos.ChaincodeInput.DecorationsEntry")
	proto.RegisterType((*ChaincodeSpec)(nil), "protos.ChaincodeSpec")
	proto.RegisterType((*ChaincodeDeploymentSpec)(nil), "protos.ChaincodeDeploymentSpec")
	proto.RegisterType((*ChaincodeInvocationSpec)(nil), "protos.ChaincodeInvocationSpec")
	proto.RegisterType((*LifecycleEvent)(nil), "protos.LifecycleEvent")
	proto.RegisterType((*ChaincodeInstallPackage)(nil), "protos.ChaincodeInstallPackage")
	proto.RegisterEnum("protos.ConfidentialityLevel", ConfidentialityLevel_name, ConfidentialityLevel_value)
	proto.RegisterEnum("protos.ChaincodeSpec_Type", ChaincodeSpec_Type_name, ChaincodeSpec_Type_value)
	proto.RegisterEnum("protos.ChaincodeDeploymentSpec_ExecutionEnvironment", ChaincodeDeploymentSpec_ExecutionEnvironment_name, ChaincodeDeploymentSpec_ExecutionEnvironment_value)
}

func init() { proto.RegisterFile("peer/chaincode.proto", fileDescriptor_chaincode_97aac793bf11f629) }

var fileDescriptor_chaincode_97aac793bf11f629 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0x86, 0xeb, 0x24, 0xfd, 0x1b, 0xa7, 0x91, 0xbf, 0xfd, 0x02, 0x44, 0x3d, 0x0a, 0x96, 0x10,
	0x01, 0x21, 0x47, 0x0a, 0x15, 0x20, 0x84, 0x2a, 0xa5, 0xb1, 0x5b, 0xb9, 0x84, 0xa4, 0x72, 0x5b,
	0x24, 0x38, 0x89, 0xdc, 0xf5, 0x24, 0x59, 0xd5, 0x59, 0x5b, 0xf6, 0xc6, 0xaa, 0x2f, 0x83, 0x2b,
	0xe1, 0x12, 0x41, 0xbb, 0x6e, 0xfe, 0x68, 0x25, 0x0e, 0x38, 0xf2, 0xee, 0xf8, 0xdd, 0x99, 0x79,
	0x1f, 0xcf, 0x1a, 0xea, 0x31, 0x62, 0xd2, 0xa6, 0x53, 0x9f, 0x71, 0x1a, 0x05, 0x68, 0xc5, 0x49,
	0x24, 0x22, 0xb2, 0xa3, 0x1e, 0xa9, 0x39, 0x04, 0xbd, 0xb7, 0x78, 0xe5, 0xda, 0x84, 0x40, 0x25,
	0xf6, 0xc5, 0xb4, 0xa1, 0x35, 0xb5, 0xd6, 0xbe, 0xa7, 0xd6, 0x32, 0xc6, 0xfd, 0x19, 0x36, 0x4a,
	0x45, 0x4c, 0xae, 0x49, 0x03, 0x76, 0x33, 0x4c, 0x52, 0x16, 0xf1, 0x46, 0x59, 0x85, 0x17, 0x5b,
	0xf3, 0xa7, 0x06, 0xb5, 0x55, 0x46, 0x1e, 0xcf, 0x85, 0x4c, 0xe0, 0x27, 0x93, 0xb4, 0xa1, 0x35,
	0xcb, 0xad, 0xaa, 0xa7, 0xd6, 0xc4, 0x05, 0x3d, 0x40, 0x1a, 0x25, 0xbe, 0x60, 0x11, 0x4f, 0x1b,
	0xa5, 0x66, 0xb9, 0xa5, 0x77, 0x5e, 0x16, 0xcd, 0xa5, 0xd6, 0x66, 0x02, 0xcb, 0x5e, 0x29, 0x1d,
	0x2e, 0x92, 0xdc, 0x5b, 0x3f, 0x7b, 0x78, 0x0c, 0xc6, 0x9f, 0x02, 0x62, 0x40, 0xf9, 0x16, 0xf3,
	0x7b, 0x1b, 0x72, 0x49, 0xea, 0xb0, 0x9d, 0xf9, 0xe1, 0xbc, 0xb0, 0x51, 0xf5, 0x8a, 0xcd, 0xc7,
	0xd2, 0x07, 0xcd, 0xfc, 0xa5, 0xc1, 0xc1, 0xb2, 0xe0, 0x65, 0x8c, 0x94, 0x58, 0x50, 0x11, 0x79,
	0x8c, 0xea, 0x78, 0xad, 0x73, 0xf8, 0xa0, 0x2b, 0x29, 0xb2, 0xae, 0xf2, 0x18, 0x3d, 0xa5, 0x23,
	0xef, 0xa0, 0xba, 0xe4, 0x3b, 0x62, 0x81, 0x2a, 0xa1, 0x77, 0xfe, 0x7f, 0xe8, 0xc6, 0xf6, 0xf4,
	0xa5, 0xd0, 0x0d, 0xc8, 0x1b, 0xd8, 0x66, 0xd2, 0xa0, 0x62, 0xa8, 0x77, 0x9e, 0x3e, 0x6e, 0xdf,
	0x2b, 0x44, 0x92, 0xb9, 0x60, 0x33, 0x8c, 0xe6, 0xa2, 0x51, 0x69, 0x6a, 0xad, 0x6d, 0x6f, 0xb1,
	0x35, 0x8f, 0xa1, 0x22, 0xbb, 0x21, 0x07, 0xb0, 0x7f, 0x3d, 0xb0, 0x9d, 0x53, 0x77, 0xe0, 0xd8,
	0xc6, 0x16, 0x01, 0xd8, 0x39, 0x1b, 0xf6, 0xbb, 0x83, 0x33, 0x43, 0x23, 0x7b, 0x50, 0x19, 0x0c,
	0x6d, 0xc7, 0x28, 0x91, 0x5d, 0x28, 0xf7, 0xba, 0x9e, 0x51, 0x96, 0xa1, 0xf3, 0xee, 0xd7, 0xae,
	0x51, 0x31, 0x7f, 0x94, 0xe0, 0xd9, 0xb2, 0xa6, 0x8d, 0x71, 0x18, 0xe5, 0x33, 0xe4, 0x42, 0xb1,
	0xf8, 0x04, 0xb5, 0x95, 0xb7, 0x34, 0x46, 0xaa, 0xa8, 0xe8, 0x9d, 0x27, 0x8f, 0x52, 0xf1, 0x0e,
	0xe8, 0x06, 0xc9, 0xe7, 0x50, 0x55, 0x07, 0x63, 0x9f, 0xde, 0xfa, 0x13, 0x54, 0x46, 0xab, 0x9e,
	0x2e, 0x63, 0x17, 0x45, 0x88, 0x0c, 0x61, 0x0f, 0xef, 0x90, 0x8e, 0x90, 0x67, 0xca, 0x57, 0xad,
	0x73, 0xf4, 0x20, 0xf5, 0x66, 0x4f, 0x96, 0x73, 0x87, 0x74, 0x2e, 0xbf, 0xb6, 0xc3, 0x33, 0x96,
	0x44, 0x5c, 0xbe, 0xf0, 0x76, 0x65, 0x16, 0x87, 0x67, 0xa6, 0x05, 0xf5, 0xc7, 0x04, 0x12, 0x87,
	0x3d, 0xec, 0x7d, 0x76, 0xbc, 0x02, 0xcd, 0xe5, 0xb7, 0xcb, 0x2b, 0xe7, 0x8b, 0xa1, 0x9d, 0x57,
	0xf6, 0x4a, 0x46, 0xd9, 0xab, 0xe1, 0x78, 0x8c, 0x54, 0xb0, 0x0c, 0x47, 0x81, 0x2f, 0xd0, 0x8c,
	0xd7, 0x90, 0xb8, 0x3c, 0x8b, 0xa8, 0x1a, 0xaf, 0x7f, 0x47, 0x72, 0x5f, 0xee, 0x3f, 0x16, 0x8c,
	0x26, 0xc8, 0xb1, 0x98, 0xda, 0x91, 0x1f, 0x4e, 0xcc, 0xf7, 0x50, 0xeb, 0xb3, 0x31, 0xd2, 0x9c,
	0x86, 0xe8, 0x64, 0xb2, 0xe3, 0x17, 0xeb, 0x85, 0xd4, 0x1d, 0x2c, 0x06, 0x7a, 0x95, 0x71, 0xe0,
	0xcf, 0xd0, 0x0c, 0x36, 0x5a, 0x4d, 0x85, 0x1f, 0x86, 0x0b, 0xb8, 0x64, 0x6d, 0x92, 0xf7, 0xef,
	0xa7, 0x75, 0x71, 0xc7, 0x4b, 0x6b, 0x77, 0xfc, 0xef, 0xdf, 0xe9, 0xf5, 0x11, 0xd4, 0x7b, 0x11,
	0x1f, 0xb3, 0x00, 0xb9, 0x60, 0x7e, 0xc8, 0x44, 0xde, 0xc7, 0x0c, 0x43, 0x89, 0xf2, 0xe2, 0xfa,
	0xa4, 0xef, 0xf6, 0x8c, 0x2d, 0x62, 0x40, 0xb5, 0x37, 0x1c, 0x9c, 0xba, 0xb6, 0x33, 0xb8, 0x72,
	0xbb, 0x7d, 0x43, 0x3b, 0x19, 0x82, 0x19, 0x25, 0x13, 0x6b, 0x9a, 0xc7, 0x98, 0x84, 0x18, 0x4c,
	0x30, 0xb1, 0xc6, 0xfe, 0x4d, 0xc2, 0xe8, 0x82, 0x95, 0xfc, 0x3b, 0x7d, 0x7f, 0x35, 0x61, 0x62,
	0x3a, 0xbf, 0xb1, 0x68, 0x34, 0x6b, 0xaf, 0x49, 0xdb, 0x85, 0xb4, 0x5d, 0x48, 0xdb, 0x52, 0x7a,
	0x53, 0xfc, 0xb8, 0xde, 0xfe, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x4f, 0x23, 0xff, 0x1c, 0xd7, 0x04,
	0x00, 0x00,
}
