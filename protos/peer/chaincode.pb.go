


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

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
	return fileDescriptor_chaincode_f10dd7a7cd7af370, []int{2, 0}
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
	return fileDescriptor_chaincode_f10dd7a7cd7af370, []int{3, 0}
}








type ChaincodeID struct {
	
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	
	
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	
	Version              string   `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeID) Reset()         { *m = ChaincodeID{} }
func (m *ChaincodeID) String() string { return proto.CompactTextString(m) }
func (*ChaincodeID) ProtoMessage()    {}
func (*ChaincodeID) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_f10dd7a7cd7af370, []int{0}
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
	Args        [][]byte          `protobuf:"bytes,1,rep,name=args,proto3" json:"args,omitempty"`
	Decorations map[string][]byte `protobuf:"bytes,2,rep,name=decorations,proto3" json:"decorations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	
	
	
	
	IsInit               bool     `protobuf:"varint,3,opt,name=is_init,json=isInit,proto3" json:"is_init,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeInput) Reset()         { *m = ChaincodeInput{} }
func (m *ChaincodeInput) String() string { return proto.CompactTextString(m) }
func (*ChaincodeInput) ProtoMessage()    {}
func (*ChaincodeInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_f10dd7a7cd7af370, []int{1}
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

func (m *ChaincodeInput) GetIsInit() bool {
	if m != nil {
		return m.IsInit
	}
	return false
}



type ChaincodeSpec struct {
	Type                 ChaincodeSpec_Type `protobuf:"varint,1,opt,name=type,proto3,enum=protos.ChaincodeSpec_Type" json:"type,omitempty"`
	ChaincodeId          *ChaincodeID       `protobuf:"bytes,2,opt,name=chaincode_id,json=chaincodeId,proto3" json:"chaincode_id,omitempty"`
	Input                *ChaincodeInput    `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`
	Timeout              int32              `protobuf:"varint,4,opt,name=timeout,proto3" json:"timeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ChaincodeSpec) Reset()         { *m = ChaincodeSpec{} }
func (m *ChaincodeSpec) String() string { return proto.CompactTextString(m) }
func (*ChaincodeSpec) ProtoMessage()    {}
func (*ChaincodeSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_f10dd7a7cd7af370, []int{2}
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
	ChaincodeSpec        *ChaincodeSpec                               `protobuf:"bytes,1,opt,name=chaincode_spec,json=chaincodeSpec,proto3" json:"chaincode_spec,omitempty"`
	CodePackage          []byte                                       `protobuf:"bytes,3,opt,name=code_package,json=codePackage,proto3" json:"code_package,omitempty"`
	ExecEnv              ChaincodeDeploymentSpec_ExecutionEnvironment `protobuf:"varint,4,opt,name=exec_env,json=execEnv,proto3,enum=protos.ChaincodeDeploymentSpec_ExecutionEnvironment" json:"exec_env,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                     `json:"-"`
	XXX_unrecognized     []byte                                       `json:"-"`
	XXX_sizecache        int32                                        `json:"-"`
}

func (m *ChaincodeDeploymentSpec) Reset()         { *m = ChaincodeDeploymentSpec{} }
func (m *ChaincodeDeploymentSpec) String() string { return proto.CompactTextString(m) }
func (*ChaincodeDeploymentSpec) ProtoMessage()    {}
func (*ChaincodeDeploymentSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_f10dd7a7cd7af370, []int{3}
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
	ChaincodeSpec        *ChaincodeSpec `protobuf:"bytes,1,opt,name=chaincode_spec,json=chaincodeSpec,proto3" json:"chaincode_spec,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ChaincodeInvocationSpec) Reset()         { *m = ChaincodeInvocationSpec{} }
func (m *ChaincodeInvocationSpec) String() string { return proto.CompactTextString(m) }
func (*ChaincodeInvocationSpec) ProtoMessage()    {}
func (*ChaincodeInvocationSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_f10dd7a7cd7af370, []int{4}
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
	ChaincodeName        string   `protobuf:"bytes,1,opt,name=chaincode_name,json=chaincodeName,proto3" json:"chaincode_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LifecycleEvent) Reset()         { *m = LifecycleEvent{} }
func (m *LifecycleEvent) String() string { return proto.CompactTextString(m) }
func (*LifecycleEvent) ProtoMessage()    {}
func (*LifecycleEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_f10dd7a7cd7af370, []int{5}
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

func init() {
	proto.RegisterType((*ChaincodeID)(nil), "protos.ChaincodeID")
	proto.RegisterType((*ChaincodeInput)(nil), "protos.ChaincodeInput")
	proto.RegisterMapType((map[string][]byte)(nil), "protos.ChaincodeInput.DecorationsEntry")
	proto.RegisterType((*ChaincodeSpec)(nil), "protos.ChaincodeSpec")
	proto.RegisterType((*ChaincodeDeploymentSpec)(nil), "protos.ChaincodeDeploymentSpec")
	proto.RegisterType((*ChaincodeInvocationSpec)(nil), "protos.ChaincodeInvocationSpec")
	proto.RegisterType((*LifecycleEvent)(nil), "protos.LifecycleEvent")
	proto.RegisterEnum("protos.ChaincodeSpec_Type", ChaincodeSpec_Type_name, ChaincodeSpec_Type_value)
	proto.RegisterEnum("protos.ChaincodeDeploymentSpec_ExecutionEnvironment", ChaincodeDeploymentSpec_ExecutionEnvironment_name, ChaincodeDeploymentSpec_ExecutionEnvironment_value)
}

func init() { proto.RegisterFile("peer/chaincode.proto", fileDescriptor_chaincode_f10dd7a7cd7af370) }

var fileDescriptor_chaincode_f10dd7a7cd7af370 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x5d, 0x6f, 0xd3, 0x4a,
	0x10, 0xbd, 0x4e, 0xdc, 0x26, 0x1d, 0xa7, 0x91, 0xef, 0xde, 0x5e, 0x6a, 0xf5, 0x29, 0x58, 0x42,
	0x04, 0x09, 0x39, 0x52, 0x40, 0x80, 0x10, 0xaa, 0x14, 0x6a, 0x53, 0xa5, 0x94, 0x04, 0x6d, 0x0b,
	0x12, 0xbc, 0x44, 0xee, 0x7a, 0xe2, 0xac, 0x9a, 0xac, 0x2d, 0x7b, 0x63, 0xd5, 0x3f, 0x83, 0x7f,
	0xc5, 0xbf, 0x02, 0xed, 0xba, 0xf9, 0x28, 0xed, 0x1b, 0x4f, 0x9e, 0x9d, 0x3d, 0x33, 0x73, 0xce,
	0xd9, 0xf5, 0xc2, 0x41, 0x8a, 0x98, 0xf5, 0xd8, 0x2c, 0xe4, 0x82, 0x25, 0x11, 0x7a, 0x69, 0x96,
	0xc8, 0x84, 0xec, 0xea, 0x4f, 0xee, 0x8e, 0xc1, 0x3a, 0x59, 0x6d, 0x0d, 0x7d, 0x42, 0xc0, 0x4c,
	0x43, 0x39, 0x73, 0x8c, 0x8e, 0xd1, 0xdd, 0xa3, 0x3a, 0x56, 0x39, 0x11, 0x2e, 0xd0, 0xa9, 0x55,
	0x39, 0x15, 0x13, 0x07, 0x1a, 0x05, 0x66, 0x39, 0x4f, 0x84, 0x53, 0xd7, 0xe9, 0xd5, 0xd2, 0xfd,
	0x69, 0x40, 0x7b, 0xd3, 0x51, 0xa4, 0x4b, 0xa9, 0x1a, 0x84, 0x59, 0x9c, 0x3b, 0x46, 0xa7, 0xde,
	0x6d, 0x51, 0x1d, 0x93, 0x21, 0x58, 0x11, 0xb2, 0x24, 0x0b, 0x25, 0x4f, 0x44, 0xee, 0xd4, 0x3a,
	0xf5, 0xae, 0xd5, 0x7f, 0x5a, 0x91, 0xcb, 0xbd, 0xbb, 0x0d, 0x3c, 0x7f, 0x83, 0x0c, 0x84, 0xcc,
	0x4a, 0xba, 0x5d, 0x4b, 0x0e, 0xa1, 0xc1, 0xf3, 0x09, 0x17, 0x5c, 0x6a, 0x2e, 0x4d, 0xba, 0xcb,
	0xf3, 0xa1, 0xe0, 0xf2, 0xe8, 0x18, 0xec, 0x3f, 0x2b, 0x89, 0x0d, 0xf5, 0x6b, 0x2c, 0x6f, 0xf5,
	0xa9, 0x90, 0x1c, 0xc0, 0x4e, 0x11, 0xce, 0x97, 0x95, 0xbe, 0x16, 0xad, 0x16, 0x6f, 0x6b, 0x6f,
	0x0c, 0xf7, 0x97, 0x01, 0xfb, 0x6b, 0x26, 0x17, 0x29, 0x32, 0xe2, 0x81, 0x29, 0xcb, 0x14, 0x75,
	0x79, 0xbb, 0x7f, 0x74, 0x8f, 0xae, 0x02, 0x79, 0x97, 0x65, 0x8a, 0x54, 0xe3, 0xc8, 0x2b, 0x68,
	0xad, 0x8d, 0x9f, 0xf0, 0x48, 0x8f, 0xb0, 0xfa, 0xff, 0xdd, 0x97, 0xe9, 0x53, 0x6b, 0x0d, 0x1c,
	0x46, 0xe4, 0x39, 0xec, 0x70, 0xa5, 0x5c, 0x0b, 0xb2, 0xfa, 0x8f, 0x1e, 0xf6, 0x85, 0x56, 0x20,
	0x75, 0x18, 0x92, 0x2f, 0x30, 0x59, 0x4a, 0xc7, 0xec, 0x18, 0xdd, 0x1d, 0xba, 0x5a, 0xba, 0xc7,
	0x60, 0x2a, 0x36, 0x64, 0x1f, 0xf6, 0xbe, 0x8c, 0xfc, 0xe0, 0xc3, 0x70, 0x14, 0xf8, 0xf6, 0x3f,
	0x04, 0x60, 0xf7, 0x74, 0x7c, 0x3e, 0x18, 0x9d, 0xda, 0x06, 0x69, 0x82, 0x39, 0x1a, 0xfb, 0x81,
	0x5d, 0x23, 0x0d, 0xa8, 0x9f, 0x0c, 0xa8, 0x5d, 0x57, 0xa9, 0xb3, 0xc1, 0xd7, 0x81, 0x6d, 0xba,
	0x3f, 0x6a, 0x70, 0xb8, 0x9e, 0xe9, 0x63, 0x3a, 0x4f, 0xca, 0x05, 0x0a, 0xa9, 0xbd, 0x78, 0x07,
	0xed, 0x8d, 0xb6, 0x3c, 0x45, 0xa6, 0x5d, 0xb1, 0xfa, 0xff, 0x3f, 0xe8, 0x0a, 0xdd, 0x67, 0x77,
	0x9c, 0x7c, 0x0c, 0x2d, 0x5d, 0x98, 0x86, 0xec, 0x3a, 0x8c, 0x51, 0x0b, 0x6d, 0x51, 0x4b, 0xe5,
	0x3e, 0x57, 0x29, 0x32, 0x86, 0x26, 0xde, 0x20, 0x9b, 0xa0, 0x28, 0xb4, 0xae, 0x76, 0xff, 0xe5,
	0xbd, 0xd6, 0x77, 0x39, 0x79, 0xc1, 0x0d, 0xb2, 0xa5, 0x3a, 0xed, 0x40, 0x14, 0x3c, 0x4b, 0x84,
	0xda, 0xa0, 0x0d, 0xd5, 0x25, 0x10, 0x85, 0xeb, 0xc1, 0xc1, 0x43, 0x00, 0x65, 0x87, 0x3f, 0x3e,
	0xf9, 0x18, 0xd0, 0xca, 0x9a, 0x8b, 0x6f, 0x17, 0x97, 0xc1, 0x27, 0xdb, 0x38, 0x33, 0x9b, 0x35,
	0xbb, 0x4e, 0xdb, 0x38, 0x9d, 0x22, 0x93, 0xbc, 0xc0, 0x49, 0x14, 0x4a, 0x74, 0xd3, 0x2d, 0x4b,
	0x86, 0xa2, 0x48, 0x98, 0xbe, 0x5e, 0x7f, 0x6f, 0xc9, 0xed, 0xb8, 0x7f, 0x79, 0x34, 0x89, 0x51,
	0x60, 0x75, 0x6b, 0x27, 0xe1, 0x3c, 0x76, 0x5f, 0x43, 0xfb, 0x9c, 0x4f, 0x91, 0x95, 0x6c, 0x8e,
	0x41, 0xa1, 0x18, 0x3f, 0xd9, 0x1e, 0xa4, 0x7f, 0xce, 0xea, 0x42, 0x6f, 0x3a, 0x8e, 0xc2, 0x05,
	0xbe, 0x1f, 0x83, 0x9b, 0x64, 0xb1, 0x37, 0x2b, 0x53, 0xcc, 0xe6, 0x18, 0xc5, 0x98, 0x79, 0xd3,
	0xf0, 0x2a, 0xe3, 0x6c, 0xc5, 0x47, 0x3d, 0x0d, 0xdf, 0x9f, 0xc5, 0x5c, 0xce, 0x96, 0x57, 0x1e,
	0x4b, 0x16, 0xbd, 0x2d, 0x68, 0xaf, 0x82, 0xf6, 0x2a, 0x68, 0x4f, 0x41, 0xaf, 0xaa, 0x57, 0xe3,
	0xc5, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x85, 0x3e, 0x7a, 0x54, 0x04, 0x00, 0x00,
}
