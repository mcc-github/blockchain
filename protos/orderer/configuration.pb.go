


package orderer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



type ConsensusType_MigrationState int32

const (
	ConsensusType_MIG_STATE_NONE    ConsensusType_MigrationState = 0
	ConsensusType_MIG_STATE_START   ConsensusType_MigrationState = 1
	ConsensusType_MIG_STATE_COMMIT  ConsensusType_MigrationState = 2
	ConsensusType_MIG_STATE_ABORT   ConsensusType_MigrationState = 3
	ConsensusType_MIG_STATE_CONTEXT ConsensusType_MigrationState = 4
)

var ConsensusType_MigrationState_name = map[int32]string{
	0: "MIG_STATE_NONE",
	1: "MIG_STATE_START",
	2: "MIG_STATE_COMMIT",
	3: "MIG_STATE_ABORT",
	4: "MIG_STATE_CONTEXT",
}
var ConsensusType_MigrationState_value = map[string]int32{
	"MIG_STATE_NONE":    0,
	"MIG_STATE_START":   1,
	"MIG_STATE_COMMIT":  2,
	"MIG_STATE_ABORT":   3,
	"MIG_STATE_CONTEXT": 4,
}

func (x ConsensusType_MigrationState) String() string {
	return proto.EnumName(ConsensusType_MigrationState_name, int32(x))
}
func (ConsensusType_MigrationState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_configuration_9ce28addc6cd8ab7, []int{0, 0}
}

type ConsensusType struct {
	
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	
	Metadata []byte `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	
	
	
	
	MigrationState ConsensusType_MigrationState `protobuf:"varint,3,opt,name=migration_state,json=migrationState,proto3,enum=orderer.ConsensusType_MigrationState" json:"migration_state,omitempty"`
	
	
	
	
	MigrationContext     uint64   `protobuf:"varint,4,opt,name=migration_context,json=migrationContext,proto3" json:"migration_context,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsensusType) Reset()         { *m = ConsensusType{} }
func (m *ConsensusType) String() string { return proto.CompactTextString(m) }
func (*ConsensusType) ProtoMessage()    {}
func (*ConsensusType) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_9ce28addc6cd8ab7, []int{0}
}
func (m *ConsensusType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsensusType.Unmarshal(m, b)
}
func (m *ConsensusType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsensusType.Marshal(b, m, deterministic)
}
func (dst *ConsensusType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusType.Merge(dst, src)
}
func (m *ConsensusType) XXX_Size() int {
	return xxx_messageInfo_ConsensusType.Size(m)
}
func (m *ConsensusType) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusType.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusType proto.InternalMessageInfo

func (m *ConsensusType) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ConsensusType) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *ConsensusType) GetMigrationState() ConsensusType_MigrationState {
	if m != nil {
		return m.MigrationState
	}
	return ConsensusType_MIG_STATE_NONE
}

func (m *ConsensusType) GetMigrationContext() uint64 {
	if m != nil {
		return m.MigrationContext
	}
	return 0
}

type BatchSize struct {
	
	
	MaxMessageCount uint32 `protobuf:"varint,1,opt,name=max_message_count,json=maxMessageCount,proto3" json:"max_message_count,omitempty"`
	
	
	AbsoluteMaxBytes uint32 `protobuf:"varint,2,opt,name=absolute_max_bytes,json=absoluteMaxBytes,proto3" json:"absolute_max_bytes,omitempty"`
	
	
	PreferredMaxBytes    uint32   `protobuf:"varint,3,opt,name=preferred_max_bytes,json=preferredMaxBytes,proto3" json:"preferred_max_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchSize) Reset()         { *m = BatchSize{} }
func (m *BatchSize) String() string { return proto.CompactTextString(m) }
func (*BatchSize) ProtoMessage()    {}
func (*BatchSize) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_9ce28addc6cd8ab7, []int{1}
}
func (m *BatchSize) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BatchSize.Unmarshal(m, b)
}
func (m *BatchSize) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BatchSize.Marshal(b, m, deterministic)
}
func (dst *BatchSize) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchSize.Merge(dst, src)
}
func (m *BatchSize) XXX_Size() int {
	return xxx_messageInfo_BatchSize.Size(m)
}
func (m *BatchSize) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchSize.DiscardUnknown(m)
}

var xxx_messageInfo_BatchSize proto.InternalMessageInfo

func (m *BatchSize) GetMaxMessageCount() uint32 {
	if m != nil {
		return m.MaxMessageCount
	}
	return 0
}

func (m *BatchSize) GetAbsoluteMaxBytes() uint32 {
	if m != nil {
		return m.AbsoluteMaxBytes
	}
	return 0
}

func (m *BatchSize) GetPreferredMaxBytes() uint32 {
	if m != nil {
		return m.PreferredMaxBytes
	}
	return 0
}

type BatchTimeout struct {
	
	
	Timeout              string   `protobuf:"bytes,1,opt,name=timeout,proto3" json:"timeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchTimeout) Reset()         { *m = BatchTimeout{} }
func (m *BatchTimeout) String() string { return proto.CompactTextString(m) }
func (*BatchTimeout) ProtoMessage()    {}
func (*BatchTimeout) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_9ce28addc6cd8ab7, []int{2}
}
func (m *BatchTimeout) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BatchTimeout.Unmarshal(m, b)
}
func (m *BatchTimeout) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BatchTimeout.Marshal(b, m, deterministic)
}
func (dst *BatchTimeout) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchTimeout.Merge(dst, src)
}
func (m *BatchTimeout) XXX_Size() int {
	return xxx_messageInfo_BatchTimeout.Size(m)
}
func (m *BatchTimeout) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchTimeout.DiscardUnknown(m)
}

var xxx_messageInfo_BatchTimeout proto.InternalMessageInfo

func (m *BatchTimeout) GetTimeout() string {
	if m != nil {
		return m.Timeout
	}
	return ""
}



type KafkaBrokers struct {
	
	
	Brokers              []string `protobuf:"bytes,1,rep,name=brokers,proto3" json:"brokers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KafkaBrokers) Reset()         { *m = KafkaBrokers{} }
func (m *KafkaBrokers) String() string { return proto.CompactTextString(m) }
func (*KafkaBrokers) ProtoMessage()    {}
func (*KafkaBrokers) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_9ce28addc6cd8ab7, []int{3}
}
func (m *KafkaBrokers) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaBrokers.Unmarshal(m, b)
}
func (m *KafkaBrokers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaBrokers.Marshal(b, m, deterministic)
}
func (dst *KafkaBrokers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaBrokers.Merge(dst, src)
}
func (m *KafkaBrokers) XXX_Size() int {
	return xxx_messageInfo_KafkaBrokers.Size(m)
}
func (m *KafkaBrokers) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaBrokers.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaBrokers proto.InternalMessageInfo

func (m *KafkaBrokers) GetBrokers() []string {
	if m != nil {
		return m.Brokers
	}
	return nil
}


type ChannelRestrictions struct {
	MaxCount             uint64   `protobuf:"varint,1,opt,name=max_count,json=maxCount,proto3" json:"max_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChannelRestrictions) Reset()         { *m = ChannelRestrictions{} }
func (m *ChannelRestrictions) String() string { return proto.CompactTextString(m) }
func (*ChannelRestrictions) ProtoMessage()    {}
func (*ChannelRestrictions) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_9ce28addc6cd8ab7, []int{4}
}
func (m *ChannelRestrictions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChannelRestrictions.Unmarshal(m, b)
}
func (m *ChannelRestrictions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChannelRestrictions.Marshal(b, m, deterministic)
}
func (dst *ChannelRestrictions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChannelRestrictions.Merge(dst, src)
}
func (m *ChannelRestrictions) XXX_Size() int {
	return xxx_messageInfo_ChannelRestrictions.Size(m)
}
func (m *ChannelRestrictions) XXX_DiscardUnknown() {
	xxx_messageInfo_ChannelRestrictions.DiscardUnknown(m)
}

var xxx_messageInfo_ChannelRestrictions proto.InternalMessageInfo

func (m *ChannelRestrictions) GetMaxCount() uint64 {
	if m != nil {
		return m.MaxCount
	}
	return 0
}

func init() {
	proto.RegisterType((*ConsensusType)(nil), "orderer.ConsensusType")
	proto.RegisterType((*BatchSize)(nil), "orderer.BatchSize")
	proto.RegisterType((*BatchTimeout)(nil), "orderer.BatchTimeout")
	proto.RegisterType((*KafkaBrokers)(nil), "orderer.KafkaBrokers")
	proto.RegisterType((*ChannelRestrictions)(nil), "orderer.ChannelRestrictions")
	proto.RegisterEnum("orderer.ConsensusType_MigrationState", ConsensusType_MigrationState_name, ConsensusType_MigrationState_value)
}

func init() {
	proto.RegisterFile("orderer/configuration.proto", fileDescriptor_configuration_9ce28addc6cd8ab7)
}

var fileDescriptor_configuration_9ce28addc6cd8ab7 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xdf, 0x8a, 0xda, 0x40,
	0x14, 0xc6, 0x1b, 0x95, 0xee, 0x7a, 0x58, 0x35, 0x8e, 0x2d, 0x84, 0xee, 0x8d, 0x08, 0x0b, 0xd2,
	0x2e, 0x11, 0xb6, 0x4f, 0x60, 0x82, 0x94, 0xa5, 0x44, 0x61, 0x9c, 0x42, 0xe9, 0x4d, 0x98, 0xc4,
	0x63, 0x0c, 0x6b, 0x32, 0x61, 0x66, 0x02, 0xb1, 0x7d, 0x8e, 0x3e, 0x4c, 0xdf, 0xae, 0xe4, 0x9f,
	0x7f, 0xee, 0xce, 0xf7, 0x7d, 0xbf, 0x9c, 0x99, 0x39, 0x39, 0xf0, 0x28, 0xe4, 0x0e, 0x25, 0xca,
	0x45, 0x28, 0xd2, 0x7d, 0x1c, 0xe5, 0x92, 0xeb, 0x58, 0xa4, 0x76, 0x26, 0x85, 0x16, 0xe4, 0xae,
	0x09, 0x67, 0xff, 0x3a, 0x30, 0x70, 0x45, 0xaa, 0x30, 0x55, 0xb9, 0x62, 0xa7, 0x0c, 0x09, 0x81,
	0x9e, 0x3e, 0x65, 0x68, 0x19, 0x53, 0x63, 0xde, 0xa7, 0x55, 0x4d, 0x3e, 0xc1, 0x7d, 0x82, 0x9a,
	0xef, 0xb8, 0xe6, 0x56, 0x67, 0x6a, 0xcc, 0x1f, 0xe8, 0x59, 0x93, 0x35, 0x8c, 0x92, 0x38, 0xaa,
	0xbb, 0xfb, 0x4a, 0x73, 0x8d, 0x56, 0x77, 0x6a, 0xcc, 0x87, 0x2f, 0x4f, 0x76, 0x73, 0x88, 0x7d,
	0x73, 0x80, 0xed, 0xb5, 0xf4, 0xb6, 0x84, 0xe9, 0x30, 0xb9, 0xd1, 0xe4, 0x0b, 0x8c, 0x2f, 0xfd,
	0x42, 0x91, 0x6a, 0x2c, 0xb4, 0xd5, 0x9b, 0x1a, 0xf3, 0x1e, 0x35, 0xcf, 0x81, 0x5b, 0xfb, 0xb3,
	0x3f, 0x30, 0xbc, 0x6d, 0x47, 0x08, 0x0c, 0xbd, 0xd7, 0x6f, 0xfe, 0x96, 0x2d, 0xd9, 0xca, 0x5f,
	0x6f, 0xd6, 0x2b, 0xf3, 0x1d, 0x99, 0xc0, 0xe8, 0xe2, 0x6d, 0xd9, 0x92, 0x32, 0xd3, 0x20, 0x1f,
	0xc0, 0xbc, 0x98, 0xee, 0xc6, 0xf3, 0x5e, 0x99, 0xd9, 0xb9, 0x45, 0x97, 0xce, 0x86, 0x32, 0xb3,
	0x4b, 0x3e, 0xc2, 0xf8, 0x1a, 0x5d, 0xb3, 0xd5, 0x4f, 0x66, 0xf6, 0x66, 0x7f, 0x0d, 0xe8, 0x3b,
	0x5c, 0x87, 0x87, 0x6d, 0xfc, 0x1b, 0xc9, 0x67, 0x18, 0x27, 0xbc, 0xf0, 0x13, 0x54, 0x8a, 0x47,
	0xe8, 0x87, 0x22, 0x4f, 0x75, 0x35, 0xc4, 0x01, 0x1d, 0x25, 0xbc, 0xf0, 0x6a, 0xdf, 0x2d, 0x6d,
	0xf2, 0x0c, 0x84, 0x07, 0x4a, 0x1c, 0x73, 0x8d, 0x7e, 0xf9, 0x51, 0x70, 0xd2, 0xa8, 0xaa, 0xc9,
	0x0e, 0xa8, 0xd9, 0x26, 0x1e, 0x2f, 0x9c, 0xd2, 0x27, 0x36, 0x4c, 0x32, 0x89, 0x7b, 0x94, 0x12,
	0x77, 0x57, 0x78, 0xb7, 0xc2, 0xc7, 0xe7, 0xa8, 0xe5, 0x67, 0x73, 0x78, 0xa8, 0xae, 0xc5, 0xe2,
	0x04, 0x45, 0xae, 0x89, 0x05, 0x77, 0xba, 0x2e, 0x9b, 0x9f, 0xda, 0xca, 0x92, 0xfc, 0xce, 0xf7,
	0x6f, 0xdc, 0x91, 0xe2, 0x0d, 0xa5, 0x2a, 0xc9, 0xa0, 0x2e, 0x2d, 0x63, 0xda, 0x2d, 0xc9, 0x46,
	0xce, 0x5e, 0x60, 0xe2, 0x1e, 0x78, 0x9a, 0xe2, 0x91, 0xa2, 0xd2, 0x32, 0x0e, 0xcb, 0x89, 0x2b,
	0xf2, 0x08, 0xfd, 0xf2, 0x42, 0x97, 0xc7, 0xf6, 0xe8, 0x7d, 0xc2, 0x8b, 0xea, 0x95, 0xce, 0x0f,
	0x78, 0x12, 0x32, 0xb2, 0x0f, 0xa7, 0x0c, 0xe5, 0x11, 0x77, 0x11, 0x4a, 0x7b, 0xcf, 0x03, 0x19,
	0x87, 0xf5, 0x12, 0xaa, 0x76, 0x3f, 0x7e, 0x3d, 0x47, 0xb1, 0x3e, 0xe4, 0x81, 0x1d, 0x8a, 0x64,
	0x71, 0x45, 0x2f, 0x6a, 0x7a, 0x51, 0xd3, 0x8b, 0x86, 0x0e, 0xde, 0x57, 0xfa, 0xeb, 0xff, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x9f, 0x9f, 0xe4, 0x13, 0xe1, 0x02, 0x00, 0x00,
}
