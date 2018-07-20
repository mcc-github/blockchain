


package orderer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type ConsensusType struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsensusType) Reset()         { *m = ConsensusType{} }
func (m *ConsensusType) String() string { return proto.CompactTextString(m) }
func (*ConsensusType) ProtoMessage()    {}
func (*ConsensusType) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_fcdec24820b11e98, []int{0}
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

type BatchSize struct {
	
	
	MaxMessageCount uint32 `protobuf:"varint,1,opt,name=max_message_count,json=maxMessageCount" json:"max_message_count,omitempty"`
	
	
	AbsoluteMaxBytes uint32 `protobuf:"varint,2,opt,name=absolute_max_bytes,json=absoluteMaxBytes" json:"absolute_max_bytes,omitempty"`
	
	
	PreferredMaxBytes    uint32   `protobuf:"varint,3,opt,name=preferred_max_bytes,json=preferredMaxBytes" json:"preferred_max_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchSize) Reset()         { *m = BatchSize{} }
func (m *BatchSize) String() string { return proto.CompactTextString(m) }
func (*BatchSize) ProtoMessage()    {}
func (*BatchSize) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_fcdec24820b11e98, []int{1}
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
	
	
	Timeout              string   `protobuf:"bytes,1,opt,name=timeout" json:"timeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchTimeout) Reset()         { *m = BatchTimeout{} }
func (m *BatchTimeout) String() string { return proto.CompactTextString(m) }
func (*BatchTimeout) ProtoMessage()    {}
func (*BatchTimeout) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_fcdec24820b11e98, []int{2}
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
	
	
	Brokers              []string `protobuf:"bytes,1,rep,name=brokers" json:"brokers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KafkaBrokers) Reset()         { *m = KafkaBrokers{} }
func (m *KafkaBrokers) String() string { return proto.CompactTextString(m) }
func (*KafkaBrokers) ProtoMessage()    {}
func (*KafkaBrokers) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_fcdec24820b11e98, []int{3}
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
	MaxCount             uint64   `protobuf:"varint,1,opt,name=max_count,json=maxCount" json:"max_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChannelRestrictions) Reset()         { *m = ChannelRestrictions{} }
func (m *ChannelRestrictions) String() string { return proto.CompactTextString(m) }
func (*ChannelRestrictions) ProtoMessage()    {}
func (*ChannelRestrictions) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_fcdec24820b11e98, []int{4}
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
}

func init() {
	proto.RegisterFile("orderer/configuration.proto", fileDescriptor_configuration_fcdec24820b11e98)
}

var fileDescriptor_configuration_fcdec24820b11e98 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xd0, 0xcd, 0x4a, 0xc3, 0x40,
	0x10, 0x07, 0x70, 0x62, 0x8b, 0xb5, 0x8b, 0x45, 0xbb, 0xbd, 0x04, 0x7a, 0x29, 0x11, 0xa1, 0x48,
	0x49, 0x40, 0xdf, 0x20, 0x3d, 0x4a, 0x2f, 0xb1, 0x5e, 0xbc, 0x94, 0x4d, 0x3a, 0x49, 0x96, 0x36,
	0x3b, 0x61, 0x76, 0x03, 0x89, 0xef, 0xe1, 0xfb, 0xca, 0x6e, 0x52, 0xed, 0x6d, 0x3e, 0x7e, 0x0b,
	0xb3, 0x7f, 0xb6, 0x44, 0x3a, 0x02, 0x01, 0x45, 0x19, 0xaa, 0x5c, 0x16, 0x0d, 0x09, 0x23, 0x51,
	0x85, 0x35, 0xa1, 0x41, 0x3e, 0x19, 0x96, 0xc1, 0x13, 0x9b, 0x6d, 0x51, 0x69, 0x50, 0xba, 0xd1,
	0xfb, 0xae, 0x06, 0xce, 0xd9, 0xd8, 0x74, 0x35, 0xf8, 0xde, 0xca, 0x5b, 0x4f, 0x13, 0x57, 0x07,
	0x3f, 0x1e, 0x9b, 0xc6, 0xc2, 0x64, 0xe5, 0x87, 0xfc, 0x06, 0xfe, 0xc2, 0xe6, 0x95, 0x68, 0x0f,
	0x15, 0x68, 0x2d, 0x0a, 0x38, 0x64, 0xd8, 0x28, 0xe3, 0xf8, 0x2c, 0x79, 0xa8, 0x44, 0xbb, 0xeb,
	0xe7, 0x5b, 0x3b, 0xe6, 0x1b, 0xc6, 0x45, 0xaa, 0xf1, 0xdc, 0x18, 0x38, 0xd8, 0x47, 0x69, 0x67,
	0x40, 0xfb, 0x37, 0x0e, 0x3f, 0x5e, 0x36, 0x3b, 0xd1, 0xc6, 0x76, 0xce, 0x43, 0xb6, 0xa8, 0x09,
	0x72, 0x20, 0x82, 0xe3, 0x15, 0x1f, 0x39, 0x3e, 0xff, 0x5b, 0x5d, 0x7c, 0xb0, 0x66, 0xf7, 0xee,
	0xac, 0xbd, 0xac, 0x00, 0x1b, 0xc3, 0x7d, 0x36, 0x31, 0x7d, 0x39, 0x9c, 0x7f, 0x69, 0xad, 0x7c,
	0x17, 0xf9, 0x49, 0xc4, 0x84, 0x27, 0x20, 0x6d, 0x65, 0xda, 0x97, 0xbe, 0xb7, 0x1a, 0x59, 0x39,
	0xb4, 0xc1, 0x2b, 0x5b, 0x6c, 0x4b, 0xa1, 0x14, 0x9c, 0x13, 0xd0, 0x86, 0x64, 0x66, 0x53, 0xd3,
	0x7c, 0xc9, 0xa6, 0xf6, 0xa0, 0xff, 0xcf, 0x8e, 0x93, 0xbb, 0x4a, 0xb4, 0xee, 0x97, 0xf1, 0x27,
	0x7b, 0x46, 0x2a, 0xc2, 0xb2, 0xab, 0x81, 0xce, 0x70, 0x2c, 0x80, 0xc2, 0x5c, 0xa4, 0x24, 0xb3,
	0x3e, 0x6d, 0x1d, 0x0e, 0x69, 0x7f, 0x6d, 0x0a, 0x69, 0xca, 0x26, 0x0d, 0x33, 0xac, 0xa2, 0x2b,
	0x1d, 0xf5, 0x3a, 0xea, 0x75, 0x34, 0xe8, 0xf4, 0xd6, 0xf5, 0x6f, 0xbf, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x0c, 0x47, 0xa2, 0x55, 0xca, 0x01, 0x00, 0x00,
}
