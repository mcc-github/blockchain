


package common

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 






















type ConfigEnvelope struct {
	Config               *Config   `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	LastUpdate           *Envelope `protobuf:"bytes,2,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ConfigEnvelope) Reset()         { *m = ConfigEnvelope{} }
func (m *ConfigEnvelope) String() string { return proto.CompactTextString(m) }
func (*ConfigEnvelope) ProtoMessage()    {}
func (*ConfigEnvelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{0}
}

func (m *ConfigEnvelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigEnvelope.Unmarshal(m, b)
}
func (m *ConfigEnvelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigEnvelope.Marshal(b, m, deterministic)
}
func (m *ConfigEnvelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigEnvelope.Merge(m, src)
}
func (m *ConfigEnvelope) XXX_Size() int {
	return xxx_messageInfo_ConfigEnvelope.Size(m)
}
func (m *ConfigEnvelope) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigEnvelope.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigEnvelope proto.InternalMessageInfo

func (m *ConfigEnvelope) GetConfig() *Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *ConfigEnvelope) GetLastUpdate() *Envelope {
	if m != nil {
		return m.LastUpdate
	}
	return nil
}

type ConfigGroupSchema struct {
	Groups               map[string]*ConfigGroupSchema  `protobuf:"bytes,1,rep,name=groups,proto3" json:"groups,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Values               map[string]*ConfigValueSchema  `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Policies             map[string]*ConfigPolicySchema `protobuf:"bytes,3,rep,name=policies,proto3" json:"policies,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *ConfigGroupSchema) Reset()         { *m = ConfigGroupSchema{} }
func (m *ConfigGroupSchema) String() string { return proto.CompactTextString(m) }
func (*ConfigGroupSchema) ProtoMessage()    {}
func (*ConfigGroupSchema) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{1}
}

func (m *ConfigGroupSchema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigGroupSchema.Unmarshal(m, b)
}
func (m *ConfigGroupSchema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigGroupSchema.Marshal(b, m, deterministic)
}
func (m *ConfigGroupSchema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigGroupSchema.Merge(m, src)
}
func (m *ConfigGroupSchema) XXX_Size() int {
	return xxx_messageInfo_ConfigGroupSchema.Size(m)
}
func (m *ConfigGroupSchema) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigGroupSchema.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigGroupSchema proto.InternalMessageInfo

func (m *ConfigGroupSchema) GetGroups() map[string]*ConfigGroupSchema {
	if m != nil {
		return m.Groups
	}
	return nil
}

func (m *ConfigGroupSchema) GetValues() map[string]*ConfigValueSchema {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *ConfigGroupSchema) GetPolicies() map[string]*ConfigPolicySchema {
	if m != nil {
		return m.Policies
	}
	return nil
}

type ConfigValueSchema struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigValueSchema) Reset()         { *m = ConfigValueSchema{} }
func (m *ConfigValueSchema) String() string { return proto.CompactTextString(m) }
func (*ConfigValueSchema) ProtoMessage()    {}
func (*ConfigValueSchema) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{2}
}

func (m *ConfigValueSchema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigValueSchema.Unmarshal(m, b)
}
func (m *ConfigValueSchema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigValueSchema.Marshal(b, m, deterministic)
}
func (m *ConfigValueSchema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigValueSchema.Merge(m, src)
}
func (m *ConfigValueSchema) XXX_Size() int {
	return xxx_messageInfo_ConfigValueSchema.Size(m)
}
func (m *ConfigValueSchema) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigValueSchema.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigValueSchema proto.InternalMessageInfo

type ConfigPolicySchema struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigPolicySchema) Reset()         { *m = ConfigPolicySchema{} }
func (m *ConfigPolicySchema) String() string { return proto.CompactTextString(m) }
func (*ConfigPolicySchema) ProtoMessage()    {}
func (*ConfigPolicySchema) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{3}
}

func (m *ConfigPolicySchema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigPolicySchema.Unmarshal(m, b)
}
func (m *ConfigPolicySchema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigPolicySchema.Marshal(b, m, deterministic)
}
func (m *ConfigPolicySchema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigPolicySchema.Merge(m, src)
}
func (m *ConfigPolicySchema) XXX_Size() int {
	return xxx_messageInfo_ConfigPolicySchema.Size(m)
}
func (m *ConfigPolicySchema) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigPolicySchema.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigPolicySchema proto.InternalMessageInfo


type Config struct {
	Sequence             uint64       `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	ChannelGroup         *ConfigGroup `protobuf:"bytes,2,opt,name=channel_group,json=channelGroup,proto3" json:"channel_group,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{4}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *Config) GetChannelGroup() *ConfigGroup {
	if m != nil {
		return m.ChannelGroup
	}
	return nil
}

type ConfigUpdateEnvelope struct {
	ConfigUpdate         []byte             `protobuf:"bytes,1,opt,name=config_update,json=configUpdate,proto3" json:"config_update,omitempty"`
	Signatures           []*ConfigSignature `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ConfigUpdateEnvelope) Reset()         { *m = ConfigUpdateEnvelope{} }
func (m *ConfigUpdateEnvelope) String() string { return proto.CompactTextString(m) }
func (*ConfigUpdateEnvelope) ProtoMessage()    {}
func (*ConfigUpdateEnvelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{5}
}

func (m *ConfigUpdateEnvelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigUpdateEnvelope.Unmarshal(m, b)
}
func (m *ConfigUpdateEnvelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigUpdateEnvelope.Marshal(b, m, deterministic)
}
func (m *ConfigUpdateEnvelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigUpdateEnvelope.Merge(m, src)
}
func (m *ConfigUpdateEnvelope) XXX_Size() int {
	return xxx_messageInfo_ConfigUpdateEnvelope.Size(m)
}
func (m *ConfigUpdateEnvelope) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigUpdateEnvelope.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigUpdateEnvelope proto.InternalMessageInfo

func (m *ConfigUpdateEnvelope) GetConfigUpdate() []byte {
	if m != nil {
		return m.ConfigUpdate
	}
	return nil
}

func (m *ConfigUpdateEnvelope) GetSignatures() []*ConfigSignature {
	if m != nil {
		return m.Signatures
	}
	return nil
}










type ConfigUpdate struct {
	ChannelId            string            `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	ReadSet              *ConfigGroup      `protobuf:"bytes,2,opt,name=read_set,json=readSet,proto3" json:"read_set,omitempty"`
	WriteSet             *ConfigGroup      `protobuf:"bytes,3,opt,name=write_set,json=writeSet,proto3" json:"write_set,omitempty"`
	IsolatedData         map[string][]byte `protobuf:"bytes,5,rep,name=isolated_data,json=isolatedData,proto3" json:"isolated_data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ConfigUpdate) Reset()         { *m = ConfigUpdate{} }
func (m *ConfigUpdate) String() string { return proto.CompactTextString(m) }
func (*ConfigUpdate) ProtoMessage()    {}
func (*ConfigUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{6}
}

func (m *ConfigUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigUpdate.Unmarshal(m, b)
}
func (m *ConfigUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigUpdate.Marshal(b, m, deterministic)
}
func (m *ConfigUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigUpdate.Merge(m, src)
}
func (m *ConfigUpdate) XXX_Size() int {
	return xxx_messageInfo_ConfigUpdate.Size(m)
}
func (m *ConfigUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigUpdate proto.InternalMessageInfo

func (m *ConfigUpdate) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *ConfigUpdate) GetReadSet() *ConfigGroup {
	if m != nil {
		return m.ReadSet
	}
	return nil
}

func (m *ConfigUpdate) GetWriteSet() *ConfigGroup {
	if m != nil {
		return m.WriteSet
	}
	return nil
}

func (m *ConfigUpdate) GetIsolatedData() map[string][]byte {
	if m != nil {
		return m.IsolatedData
	}
	return nil
}


type ConfigGroup struct {
	Version              uint64                   `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Groups               map[string]*ConfigGroup  `protobuf:"bytes,2,rep,name=groups,proto3" json:"groups,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Values               map[string]*ConfigValue  `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Policies             map[string]*ConfigPolicy `protobuf:"bytes,4,rep,name=policies,proto3" json:"policies,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ModPolicy            string                   `protobuf:"bytes,5,opt,name=mod_policy,json=modPolicy,proto3" json:"mod_policy,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ConfigGroup) Reset()         { *m = ConfigGroup{} }
func (m *ConfigGroup) String() string { return proto.CompactTextString(m) }
func (*ConfigGroup) ProtoMessage()    {}
func (*ConfigGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{7}
}

func (m *ConfigGroup) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigGroup.Unmarshal(m, b)
}
func (m *ConfigGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigGroup.Marshal(b, m, deterministic)
}
func (m *ConfigGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigGroup.Merge(m, src)
}
func (m *ConfigGroup) XXX_Size() int {
	return xxx_messageInfo_ConfigGroup.Size(m)
}
func (m *ConfigGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigGroup.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigGroup proto.InternalMessageInfo

func (m *ConfigGroup) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ConfigGroup) GetGroups() map[string]*ConfigGroup {
	if m != nil {
		return m.Groups
	}
	return nil
}

func (m *ConfigGroup) GetValues() map[string]*ConfigValue {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *ConfigGroup) GetPolicies() map[string]*ConfigPolicy {
	if m != nil {
		return m.Policies
	}
	return nil
}

func (m *ConfigGroup) GetModPolicy() string {
	if m != nil {
		return m.ModPolicy
	}
	return ""
}


type ConfigValue struct {
	Version              uint64   `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	ModPolicy            string   `protobuf:"bytes,3,opt,name=mod_policy,json=modPolicy,proto3" json:"mod_policy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigValue) Reset()         { *m = ConfigValue{} }
func (m *ConfigValue) String() string { return proto.CompactTextString(m) }
func (*ConfigValue) ProtoMessage()    {}
func (*ConfigValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{8}
}

func (m *ConfigValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigValue.Unmarshal(m, b)
}
func (m *ConfigValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigValue.Marshal(b, m, deterministic)
}
func (m *ConfigValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigValue.Merge(m, src)
}
func (m *ConfigValue) XXX_Size() int {
	return xxx_messageInfo_ConfigValue.Size(m)
}
func (m *ConfigValue) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigValue.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigValue proto.InternalMessageInfo

func (m *ConfigValue) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ConfigValue) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *ConfigValue) GetModPolicy() string {
	if m != nil {
		return m.ModPolicy
	}
	return ""
}

type ConfigPolicy struct {
	Version              uint64   `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Policy               *Policy  `protobuf:"bytes,2,opt,name=policy,proto3" json:"policy,omitempty"`
	ModPolicy            string   `protobuf:"bytes,3,opt,name=mod_policy,json=modPolicy,proto3" json:"mod_policy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigPolicy) Reset()         { *m = ConfigPolicy{} }
func (m *ConfigPolicy) String() string { return proto.CompactTextString(m) }
func (*ConfigPolicy) ProtoMessage()    {}
func (*ConfigPolicy) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{9}
}

func (m *ConfigPolicy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigPolicy.Unmarshal(m, b)
}
func (m *ConfigPolicy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigPolicy.Marshal(b, m, deterministic)
}
func (m *ConfigPolicy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigPolicy.Merge(m, src)
}
func (m *ConfigPolicy) XXX_Size() int {
	return xxx_messageInfo_ConfigPolicy.Size(m)
}
func (m *ConfigPolicy) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigPolicy.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigPolicy proto.InternalMessageInfo

func (m *ConfigPolicy) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ConfigPolicy) GetPolicy() *Policy {
	if m != nil {
		return m.Policy
	}
	return nil
}

func (m *ConfigPolicy) GetModPolicy() string {
	if m != nil {
		return m.ModPolicy
	}
	return ""
}

type ConfigSignature struct {
	SignatureHeader      []byte   `protobuf:"bytes,1,opt,name=signature_header,json=signatureHeader,proto3" json:"signature_header,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigSignature) Reset()         { *m = ConfigSignature{} }
func (m *ConfigSignature) String() string { return proto.CompactTextString(m) }
func (*ConfigSignature) ProtoMessage()    {}
func (*ConfigSignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{10}
}

func (m *ConfigSignature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigSignature.Unmarshal(m, b)
}
func (m *ConfigSignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigSignature.Marshal(b, m, deterministic)
}
func (m *ConfigSignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigSignature.Merge(m, src)
}
func (m *ConfigSignature) XXX_Size() int {
	return xxx_messageInfo_ConfigSignature.Size(m)
}
func (m *ConfigSignature) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigSignature.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigSignature proto.InternalMessageInfo

func (m *ConfigSignature) GetSignatureHeader() []byte {
	if m != nil {
		return m.SignatureHeader
	}
	return nil
}

func (m *ConfigSignature) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*ConfigEnvelope)(nil), "common.ConfigEnvelope")
	proto.RegisterType((*ConfigGroupSchema)(nil), "common.ConfigGroupSchema")
	proto.RegisterMapType((map[string]*ConfigGroupSchema)(nil), "common.ConfigGroupSchema.GroupsEntry")
	proto.RegisterMapType((map[string]*ConfigPolicySchema)(nil), "common.ConfigGroupSchema.PoliciesEntry")
	proto.RegisterMapType((map[string]*ConfigValueSchema)(nil), "common.ConfigGroupSchema.ValuesEntry")
	proto.RegisterType((*ConfigValueSchema)(nil), "common.ConfigValueSchema")
	proto.RegisterType((*ConfigPolicySchema)(nil), "common.ConfigPolicySchema")
	proto.RegisterType((*Config)(nil), "common.Config")
	proto.RegisterType((*ConfigUpdateEnvelope)(nil), "common.ConfigUpdateEnvelope")
	proto.RegisterType((*ConfigUpdate)(nil), "common.ConfigUpdate")
	proto.RegisterMapType((map[string][]byte)(nil), "common.ConfigUpdate.IsolatedDataEntry")
	proto.RegisterType((*ConfigGroup)(nil), "common.ConfigGroup")
	proto.RegisterMapType((map[string]*ConfigGroup)(nil), "common.ConfigGroup.GroupsEntry")
	proto.RegisterMapType((map[string]*ConfigPolicy)(nil), "common.ConfigGroup.PoliciesEntry")
	proto.RegisterMapType((map[string]*ConfigValue)(nil), "common.ConfigGroup.ValuesEntry")
	proto.RegisterType((*ConfigValue)(nil), "common.ConfigValue")
	proto.RegisterType((*ConfigPolicy)(nil), "common.ConfigPolicy")
	proto.RegisterType((*ConfigSignature)(nil), "common.ConfigSignature")
}

func init() { proto.RegisterFile("common/configtx.proto", fileDescriptor_5190bbf196fa7499) }

var fileDescriptor_5190bbf196fa7499 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x56, 0x6f, 0x6f, 0xd3, 0x3e,
	0x10, 0x56, 0x9b, 0xb6, 0x6b, 0xaf, 0xed, 0xd6, 0x79, 0xfd, 0xe9, 0x17, 0x22, 0x10, 0x23, 0xc0,
	0xd8, 0x40, 0x4a, 0xc7, 0x78, 0xb1, 0x09, 0x69, 0x42, 0x62, 0x4c, 0xb0, 0x21, 0x4d, 0x90, 0xf2,
	0x47, 0x9a, 0x90, 0x2a, 0x2f, 0xf1, 0xda, 0xb0, 0x34, 0x0e, 0x89, 0x3b, 0xe8, 0x47, 0xe2, 0x33,
	0xf1, 0x0d, 0xf8, 0x14, 0x28, 0xb6, 0x13, 0x9c, 0x35, 0x6d, 0xc5, 0xab, 0xf5, 0xce, 0xcf, 0xf3,
	0xdc, 0xf9, 0xee, 0x7c, 0x19, 0xfc, 0xe7, 0xd0, 0xf1, 0x98, 0x06, 0x3d, 0x87, 0x06, 0x97, 0xde,
	0x90, 0xfd, 0xb0, 0xc2, 0x88, 0x32, 0x8a, 0x6a, 0xc2, 0x6d, 0x6c, 0x64, 0xc7, 0xc9, 0x1f, 0x71,
	0x68, 0xa4, 0x9c, 0x90, 0xfa, 0x9e, 0xe3, 0x91, 0x58, 0xb8, 0xcd, 0x2b, 0x58, 0x3d, 0xe2, 0x2a,
	0xc7, 0xc1, 0x35, 0xf1, 0x69, 0x48, 0xd0, 0x16, 0xd4, 0x84, 0xae, 0x5e, 0xda, 0x2c, 0x6d, 0x37,
	0xf7, 0x56, 0x2d, 0xa9, 0x23, 0x70, 0xb6, 0x3c, 0x45, 0x4f, 0xa1, 0xe9, 0xe3, 0x98, 0x0d, 0x26,
	0xa1, 0x8b, 0x19, 0xd1, 0xcb, 0x1c, 0xdc, 0x49, 0xc1, 0xa9, 0x9c, 0x0d, 0x09, 0xe8, 0x23, 0xc7,
	0x98, 0xbf, 0x34, 0x58, 0x17, 0x2a, 0xaf, 0x23, 0x3a, 0x09, 0xfb, 0xce, 0x88, 0x8c, 0x31, 0x3a,
	0x84, 0xda, 0x30, 0x31, 0x63, 0xbd, 0xb4, 0xa9, 0x6d, 0x37, 0xf7, 0x1e, 0xe6, 0x03, 0x2a, 0x50,
	0x8b, 0xff, 0x8e, 0x8f, 0x03, 0x16, 0x4d, 0x6d, 0x49, 0x4a, 0xe8, 0xd7, 0xd8, 0x9f, 0x90, 0x58,
	0x2f, 0x2f, 0xa3, 0x7f, 0xe2, 0x38, 0x49, 0x17, 0x24, 0x74, 0x04, 0xf5, 0xb4, 0x24, 0xba, 0xc6,
	0x05, 0x1e, 0xcd, 0x17, 0x78, 0x27, 0x91, 0x42, 0x22, 0x23, 0x1a, 0x1f, 0xa0, 0xa9, 0xa4, 0x86,
	0x3a, 0xa0, 0x5d, 0x91, 0x29, 0xaf, 0x5f, 0xc3, 0x4e, 0x7e, 0xa2, 0x1e, 0x54, 0x79, 0x3c, 0x59,
	0xa6, 0x5b, 0x73, 0x43, 0xd8, 0x02, 0xf7, 0xbc, 0x7c, 0x50, 0x4a, 0x54, 0x95, 0x8c, 0xff, 0x59,
	0x95, 0x73, 0x67, 0x55, 0x3f, 0x43, 0x3b, 0x77, 0x8d, 0x02, 0xdd, 0xdd, 0xbc, 0xae, 0x91, 0xd7,
	0xe5, 0xec, 0xe9, 0x8c, 0xb0, 0xb9, 0x91, 0x36, 0x57, 0x09, 0x6c, 0x76, 0x01, 0xcd, 0xb2, 0xcc,
	0xaf, 0x50, 0x13, 0x5e, 0x64, 0x40, 0x3d, 0x26, 0xdf, 0x26, 0x24, 0x70, 0x08, 0xcf, 0xa0, 0x62,
	0x67, 0x36, 0x3a, 0x80, 0xb6, 0x33, 0xc2, 0x41, 0x40, 0xfc, 0x01, 0xef, 0xb5, 0x4c, 0x67, 0xa3,
	0xa0, 0x78, 0x76, 0x4b, 0x22, 0xb9, 0x75, 0x5a, 0xa9, 0x6b, 0x9d, 0x8a, 0x5d, 0x61, 0xd3, 0x90,
	0x98, 0x0c, 0xba, 0x02, 0x28, 0x86, 0x30, 0x9b, 0xf3, 0xfb, 0xd0, 0x16, 0x93, 0x9c, 0x4e, 0x70,
	0x12, 0xbe, 0x65, 0xb7, 0x1c, 0x05, 0x8c, 0xf6, 0x01, 0x62, 0x6f, 0x18, 0x60, 0x36, 0x89, 0xb2,
	0x01, 0xfb, 0x3f, 0x1f, 0xbf, 0x9f, 0x9e, 0xdb, 0x0a, 0xd4, 0xfc, 0x59, 0x86, 0x96, 0x1a, 0x16,
	0xdd, 0x01, 0x48, 0x2f, 0xe3, 0xb9, 0xb2, 0xd8, 0x0d, 0xe9, 0x39, 0x71, 0x91, 0x05, 0xf5, 0x88,
	0x60, 0x77, 0x10, 0x13, 0xb6, 0xe8, 0x9a, 0x2b, 0x09, 0xa8, 0x4f, 0x18, 0xda, 0x85, 0xc6, 0xf7,
	0xc8, 0x63, 0x84, 0x13, 0xb4, 0xf9, 0x84, 0x3a, 0x47, 0x25, 0x8c, 0xb7, 0xd0, 0xf6, 0x62, 0xea,
	0x63, 0x46, 0xdc, 0x81, 0x8b, 0x19, 0xd6, 0xab, 0xfc, 0x36, 0x5b, 0x79, 0x96, 0xc8, 0xd6, 0x3a,
	0x91, 0xc8, 0x57, 0x98, 0x61, 0x31, 0xec, 0x2d, 0x4f, 0x71, 0x19, 0x2f, 0x60, 0x7d, 0x06, 0x52,
	0x30, 0x48, 0x5d, 0x75, 0x90, 0x5a, 0xca, 0xb0, 0x9c, 0x56, 0xea, 0x95, 0x4e, 0x55, 0x76, 0xe8,
	0xb7, 0x06, 0x4d, 0x25, 0x67, 0xa4, 0xc3, 0xca, 0x35, 0x89, 0x62, 0x8f, 0x06, 0x72, 0x24, 0x52,
	0x13, 0xed, 0x67, 0xab, 0x42, 0xb4, 0xe2, 0x6e, 0xc1, 0x95, 0x0b, 0x97, 0xc4, 0x7e, 0xb6, 0x24,
	0xb4, 0xf9, 0xc4, 0xa2, 0xf5, 0x70, 0xa8, 0xac, 0x87, 0x0a, 0xa7, 0xde, 0x2b, 0xa2, 0xce, 0x59,
	0x0c, 0x49, 0xd7, 0xc7, 0xd4, 0x1d, 0x70, 0x7b, 0xaa, 0x57, 0x45, 0xd7, 0xc7, 0xd4, 0x15, 0xaf,
	0xc1, 0x38, 0x5b, 0xb6, 0x37, 0x76, 0xf2, 0x2f, 0xb1, 0xb0, 0xc5, 0xca, 0xdb, 0x3e, 0x5b, 0xb6,
	0x31, 0x16, 0xeb, 0x71, 0xae, 0xaa, 0xf7, 0x7e, 0xf9, 0xae, 0x78, 0x9c, 0x57, 0xec, 0x16, 0xed,
	0x0a, 0x75, 0x4b, 0x7c, 0x49, 0x7b, 0xcd, 0x83, 0x2d, 0xe8, 0x75, 0xe1, 0xec, 0xdc, 0x28, 0xa8,
	0x76, 0xa3, 0xa0, 0x26, 0x4d, 0x5f, 0x9d, 0xb0, 0x17, 0xc8, 0x6f, 0x41, 0x4d, 0x8a, 0x94, 0xf3,
	0x9f, 0x39, 0x99, 0xb2, 0x3c, 0x5d, 0x16, 0xf0, 0x1c, 0xd6, 0x6e, 0xac, 0x01, 0xb4, 0x03, 0x9d,
	0x6c, 0x11, 0x0c, 0x46, 0x04, 0xbb, 0x24, 0x92, 0xbb, 0x65, 0x2d, 0xf3, 0xbf, 0xe1, 0x6e, 0x74,
	0x1b, 0x1a, 0x99, 0x4b, 0xde, 0xf3, 0xaf, 0xe3, 0x65, 0x1f, 0x1e, 0xd0, 0x68, 0x68, 0x8d, 0xa6,
	0x21, 0x89, 0x7c, 0xe2, 0x0e, 0x49, 0x64, 0x5d, 0xe2, 0x8b, 0xc8, 0x73, 0xc4, 0xb7, 0x3b, 0x96,
	0x19, 0x9f, 0x3f, 0x19, 0x7a, 0x6c, 0x34, 0xb9, 0x48, 0xcc, 0x9e, 0x02, 0xee, 0x09, 0x70, 0x4f,
	0x80, 0xe5, 0x7f, 0x03, 0x17, 0x35, 0x6e, 0x3e, 0xfb, 0x13, 0x00, 0x00, 0xff, 0xff, 0xa8, 0xa5,
	0xb6, 0x86, 0x44, 0x08, 0x00, 0x00,
}
