


package orderer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type KafkaMessageRegular_Class int32

const (
	KafkaMessageRegular_UNKNOWN KafkaMessageRegular_Class = 0
	KafkaMessageRegular_NORMAL  KafkaMessageRegular_Class = 1
	KafkaMessageRegular_CONFIG  KafkaMessageRegular_Class = 2
)

var KafkaMessageRegular_Class_name = map[int32]string{
	0: "UNKNOWN",
	1: "NORMAL",
	2: "CONFIG",
}
var KafkaMessageRegular_Class_value = map[string]int32{
	"UNKNOWN": 0,
	"NORMAL":  1,
	"CONFIG":  2,
}

func (x KafkaMessageRegular_Class) String() string {
	return proto.EnumName(KafkaMessageRegular_Class_name, int32(x))
}
func (KafkaMessageRegular_Class) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_kafka_d6fc97ca114dc767, []int{1, 0}
}



type KafkaMessage struct {
	
	
	
	
	Type                 isKafkaMessage_Type `protobuf_oneof:"Type"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *KafkaMessage) Reset()         { *m = KafkaMessage{} }
func (m *KafkaMessage) String() string { return proto.CompactTextString(m) }
func (*KafkaMessage) ProtoMessage()    {}
func (*KafkaMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_kafka_d6fc97ca114dc767, []int{0}
}
func (m *KafkaMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaMessage.Unmarshal(m, b)
}
func (m *KafkaMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaMessage.Marshal(b, m, deterministic)
}
func (dst *KafkaMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaMessage.Merge(dst, src)
}
func (m *KafkaMessage) XXX_Size() int {
	return xxx_messageInfo_KafkaMessage.Size(m)
}
func (m *KafkaMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaMessage.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaMessage proto.InternalMessageInfo

type isKafkaMessage_Type interface {
	isKafkaMessage_Type()
}

type KafkaMessage_Regular struct {
	Regular *KafkaMessageRegular `protobuf:"bytes,1,opt,name=regular,oneof"`
}
type KafkaMessage_TimeToCut struct {
	TimeToCut *KafkaMessageTimeToCut `protobuf:"bytes,2,opt,name=time_to_cut,json=timeToCut,oneof"`
}
type KafkaMessage_Connect struct {
	Connect *KafkaMessageConnect `protobuf:"bytes,3,opt,name=connect,oneof"`
}

func (*KafkaMessage_Regular) isKafkaMessage_Type()   {}
func (*KafkaMessage_TimeToCut) isKafkaMessage_Type() {}
func (*KafkaMessage_Connect) isKafkaMessage_Type()   {}

func (m *KafkaMessage) GetType() isKafkaMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *KafkaMessage) GetRegular() *KafkaMessageRegular {
	if x, ok := m.GetType().(*KafkaMessage_Regular); ok {
		return x.Regular
	}
	return nil
}

func (m *KafkaMessage) GetTimeToCut() *KafkaMessageTimeToCut {
	if x, ok := m.GetType().(*KafkaMessage_TimeToCut); ok {
		return x.TimeToCut
	}
	return nil
}

func (m *KafkaMessage) GetConnect() *KafkaMessageConnect {
	if x, ok := m.GetType().(*KafkaMessage_Connect); ok {
		return x.Connect
	}
	return nil
}


func (*KafkaMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _KafkaMessage_OneofMarshaler, _KafkaMessage_OneofUnmarshaler, _KafkaMessage_OneofSizer, []interface{}{
		(*KafkaMessage_Regular)(nil),
		(*KafkaMessage_TimeToCut)(nil),
		(*KafkaMessage_Connect)(nil),
	}
}

func _KafkaMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*KafkaMessage)
	
	switch x := m.Type.(type) {
	case *KafkaMessage_Regular:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Regular); err != nil {
			return err
		}
	case *KafkaMessage_TimeToCut:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TimeToCut); err != nil {
			return err
		}
	case *KafkaMessage_Connect:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Connect); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("KafkaMessage.Type has unexpected type %T", x)
	}
	return nil
}

func _KafkaMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*KafkaMessage)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(KafkaMessageRegular)
		err := b.DecodeMessage(msg)
		m.Type = &KafkaMessage_Regular{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(KafkaMessageTimeToCut)
		err := b.DecodeMessage(msg)
		m.Type = &KafkaMessage_TimeToCut{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(KafkaMessageConnect)
		err := b.DecodeMessage(msg)
		m.Type = &KafkaMessage_Connect{msg}
		return true, err
	default:
		return false, nil
	}
}

func _KafkaMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*KafkaMessage)
	
	switch x := m.Type.(type) {
	case *KafkaMessage_Regular:
		s := proto.Size(x.Regular)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *KafkaMessage_TimeToCut:
		s := proto.Size(x.TimeToCut)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *KafkaMessage_Connect:
		s := proto.Size(x.Connect)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type KafkaMessageRegular struct {
	Payload              []byte                    `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	ConfigSeq            uint64                    `protobuf:"varint,2,opt,name=config_seq,json=configSeq" json:"config_seq,omitempty"`
	Class                KafkaMessageRegular_Class `protobuf:"varint,3,opt,name=class,enum=orderer.KafkaMessageRegular_Class" json:"class,omitempty"`
	OriginalOffset       int64                     `protobuf:"varint,4,opt,name=original_offset,json=originalOffset" json:"original_offset,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *KafkaMessageRegular) Reset()         { *m = KafkaMessageRegular{} }
func (m *KafkaMessageRegular) String() string { return proto.CompactTextString(m) }
func (*KafkaMessageRegular) ProtoMessage()    {}
func (*KafkaMessageRegular) Descriptor() ([]byte, []int) {
	return fileDescriptor_kafka_d6fc97ca114dc767, []int{1}
}
func (m *KafkaMessageRegular) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaMessageRegular.Unmarshal(m, b)
}
func (m *KafkaMessageRegular) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaMessageRegular.Marshal(b, m, deterministic)
}
func (dst *KafkaMessageRegular) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaMessageRegular.Merge(dst, src)
}
func (m *KafkaMessageRegular) XXX_Size() int {
	return xxx_messageInfo_KafkaMessageRegular.Size(m)
}
func (m *KafkaMessageRegular) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaMessageRegular.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaMessageRegular proto.InternalMessageInfo

func (m *KafkaMessageRegular) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *KafkaMessageRegular) GetConfigSeq() uint64 {
	if m != nil {
		return m.ConfigSeq
	}
	return 0
}

func (m *KafkaMessageRegular) GetClass() KafkaMessageRegular_Class {
	if m != nil {
		return m.Class
	}
	return KafkaMessageRegular_UNKNOWN
}

func (m *KafkaMessageRegular) GetOriginalOffset() int64 {
	if m != nil {
		return m.OriginalOffset
	}
	return 0
}



type KafkaMessageTimeToCut struct {
	BlockNumber          uint64   `protobuf:"varint,1,opt,name=block_number,json=blockNumber" json:"block_number,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KafkaMessageTimeToCut) Reset()         { *m = KafkaMessageTimeToCut{} }
func (m *KafkaMessageTimeToCut) String() string { return proto.CompactTextString(m) }
func (*KafkaMessageTimeToCut) ProtoMessage()    {}
func (*KafkaMessageTimeToCut) Descriptor() ([]byte, []int) {
	return fileDescriptor_kafka_d6fc97ca114dc767, []int{2}
}
func (m *KafkaMessageTimeToCut) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaMessageTimeToCut.Unmarshal(m, b)
}
func (m *KafkaMessageTimeToCut) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaMessageTimeToCut.Marshal(b, m, deterministic)
}
func (dst *KafkaMessageTimeToCut) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaMessageTimeToCut.Merge(dst, src)
}
func (m *KafkaMessageTimeToCut) XXX_Size() int {
	return xxx_messageInfo_KafkaMessageTimeToCut.Size(m)
}
func (m *KafkaMessageTimeToCut) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaMessageTimeToCut.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaMessageTimeToCut proto.InternalMessageInfo

func (m *KafkaMessageTimeToCut) GetBlockNumber() uint64 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}





type KafkaMessageConnect struct {
	Payload              []byte   `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KafkaMessageConnect) Reset()         { *m = KafkaMessageConnect{} }
func (m *KafkaMessageConnect) String() string { return proto.CompactTextString(m) }
func (*KafkaMessageConnect) ProtoMessage()    {}
func (*KafkaMessageConnect) Descriptor() ([]byte, []int) {
	return fileDescriptor_kafka_d6fc97ca114dc767, []int{3}
}
func (m *KafkaMessageConnect) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaMessageConnect.Unmarshal(m, b)
}
func (m *KafkaMessageConnect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaMessageConnect.Marshal(b, m, deterministic)
}
func (dst *KafkaMessageConnect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaMessageConnect.Merge(dst, src)
}
func (m *KafkaMessageConnect) XXX_Size() int {
	return xxx_messageInfo_KafkaMessageConnect.Size(m)
}
func (m *KafkaMessageConnect) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaMessageConnect.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaMessageConnect proto.InternalMessageInfo

func (m *KafkaMessageConnect) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}



type KafkaMetadata struct {
	
	
	
	LastOffsetPersisted int64 `protobuf:"varint,1,opt,name=last_offset_persisted,json=lastOffsetPersisted" json:"last_offset_persisted,omitempty"`
	
	
	
	
	LastOriginalOffsetProcessed int64 `protobuf:"varint,2,opt,name=last_original_offset_processed,json=lastOriginalOffsetProcessed" json:"last_original_offset_processed,omitempty"`
	
	
	
	
	
	
	
	LastResubmittedConfigOffset int64    `protobuf:"varint,3,opt,name=last_resubmitted_config_offset,json=lastResubmittedConfigOffset" json:"last_resubmitted_config_offset,omitempty"`
	XXX_NoUnkeyedLiteral        struct{} `json:"-"`
	XXX_unrecognized            []byte   `json:"-"`
	XXX_sizecache               int32    `json:"-"`
}

func (m *KafkaMetadata) Reset()         { *m = KafkaMetadata{} }
func (m *KafkaMetadata) String() string { return proto.CompactTextString(m) }
func (*KafkaMetadata) ProtoMessage()    {}
func (*KafkaMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_kafka_d6fc97ca114dc767, []int{4}
}
func (m *KafkaMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaMetadata.Unmarshal(m, b)
}
func (m *KafkaMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaMetadata.Marshal(b, m, deterministic)
}
func (dst *KafkaMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaMetadata.Merge(dst, src)
}
func (m *KafkaMetadata) XXX_Size() int {
	return xxx_messageInfo_KafkaMetadata.Size(m)
}
func (m *KafkaMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaMetadata proto.InternalMessageInfo

func (m *KafkaMetadata) GetLastOffsetPersisted() int64 {
	if m != nil {
		return m.LastOffsetPersisted
	}
	return 0
}

func (m *KafkaMetadata) GetLastOriginalOffsetProcessed() int64 {
	if m != nil {
		return m.LastOriginalOffsetProcessed
	}
	return 0
}

func (m *KafkaMetadata) GetLastResubmittedConfigOffset() int64 {
	if m != nil {
		return m.LastResubmittedConfigOffset
	}
	return 0
}

func init() {
	proto.RegisterType((*KafkaMessage)(nil), "orderer.KafkaMessage")
	proto.RegisterType((*KafkaMessageRegular)(nil), "orderer.KafkaMessageRegular")
	proto.RegisterType((*KafkaMessageTimeToCut)(nil), "orderer.KafkaMessageTimeToCut")
	proto.RegisterType((*KafkaMessageConnect)(nil), "orderer.KafkaMessageConnect")
	proto.RegisterType((*KafkaMetadata)(nil), "orderer.KafkaMetadata")
	proto.RegisterEnum("orderer.KafkaMessageRegular_Class", KafkaMessageRegular_Class_name, KafkaMessageRegular_Class_value)
}

func init() { proto.RegisterFile("orderer/kafka.proto", fileDescriptor_kafka_d6fc97ca114dc767) }

var fileDescriptor_kafka_d6fc97ca114dc767 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x93, 0xd1, 0x6a, 0xdb, 0x30,
	0x14, 0x86, 0xe3, 0x26, 0x4d, 0xe8, 0x49, 0xd6, 0x05, 0x85, 0x42, 0x60, 0x5b, 0xe9, 0x0c, 0x63,
	0xbd, 0x28, 0x36, 0x64, 0x37, 0x65, 0x57, 0x5b, 0x0d, 0x5b, 0x47, 0x57, 0xa7, 0x68, 0x29, 0x83,
	0xdd, 0x18, 0xd9, 0x3e, 0x76, 0x4d, 0x6c, 0xcb, 0x95, 0xe4, 0x8b, 0xbc, 0xe3, 0xf6, 0x0c, 0x7b,
	0x95, 0x61, 0xc9, 0x5e, 0x1b, 0x48, 0x7b, 0x67, 0xfd, 0xfa, 0x7e, 0x9d, 0x73, 0xfe, 0x83, 0x61,
	0xc6, 0x45, 0x8c, 0x02, 0x85, 0xbb, 0x66, 0xc9, 0x9a, 0x39, 0x95, 0xe0, 0x8a, 0x93, 0x51, 0x2b,
	0xda, 0xbf, 0x2d, 0x98, 0x5c, 0x35, 0x17, 0xd7, 0x28, 0x25, 0x4b, 0x91, 0x9c, 0xc3, 0x48, 0x60,
	0x5a, 0xe7, 0x4c, 0xcc, 0xad, 0x13, 0xeb, 0x74, 0xbc, 0x78, 0xed, 0xb4, 0xac, 0xf3, 0x98, 0xa3,
	0x86, 0xb9, 0xec, 0xd1, 0x0e, 0x27, 0x9f, 0x60, 0xac, 0xb2, 0x02, 0x03, 0xc5, 0x83, 0xa8, 0x56,
	0xf3, 0x3d, 0xed, 0x3e, 0xde, 0xe9, 0x5e, 0x65, 0x05, 0xae, 0xb8, 0x57, 0xab, 0xcb, 0x1e, 0x3d,
	0x50, 0xdd, 0xa1, 0xa9, 0x1d, 0xf1, 0xb2, 0xc4, 0x48, 0xcd, 0xfb, 0xcf, 0xd4, 0xf6, 0x0c, 0xd3,
	0xd4, 0x6e, 0xf1, 0x8b, 0x21, 0x0c, 0x56, 0x9b, 0x0a, 0xed, 0xbf, 0x16, 0xcc, 0x76, 0xb4, 0x49,
	0xe6, 0x30, 0xaa, 0xd8, 0x26, 0xe7, 0x2c, 0xd6, 0x53, 0x4d, 0x68, 0x77, 0x24, 0x6f, 0x00, 0x22,
	0x5e, 0x26, 0x59, 0x1a, 0x48, 0xbc, 0xd7, 0x4d, 0x0f, 0xe8, 0x81, 0x51, 0x7e, 0xe0, 0x3d, 0x39,
	0x87, 0xfd, 0x28, 0x67, 0x52, 0xea, 0x86, 0x0e, 0x17, 0xf6, 0x73, 0x61, 0x38, 0x5e, 0x43, 0x52,
	0x63, 0x20, 0xef, 0xe1, 0x25, 0x17, 0x59, 0x9a, 0x95, 0x2c, 0x0f, 0x78, 0x92, 0x48, 0x54, 0xf3,
	0xc1, 0x89, 0x75, 0xda, 0xa7, 0x87, 0x9d, 0xbc, 0xd4, 0xaa, 0x7d, 0x06, 0xfb, 0xda, 0x48, 0xc6,
	0x30, 0xba, 0xf5, 0xaf, 0xfc, 0xe5, 0x4f, 0x7f, 0xda, 0x23, 0x00, 0x43, 0x7f, 0x49, 0xaf, 0x3f,
	0x7f, 0x9f, 0x5a, 0xcd, 0xb7, 0xb7, 0xf4, 0xbf, 0x7c, 0xfb, 0x3a, 0xdd, 0xb3, 0x3f, 0xc2, 0xd1,
	0xce, 0x24, 0xc9, 0x5b, 0x98, 0x84, 0x39, 0x8f, 0xd6, 0x41, 0x59, 0x17, 0x21, 0x9a, 0xed, 0x0d,
	0xe8, 0x58, 0x6b, 0xbe, 0x96, 0x6c, 0x77, 0x3b, 0x9c, 0x36, 0xc7, 0xa7, 0xc3, 0xb1, 0xff, 0x58,
	0xf0, 0xa2, 0x75, 0x28, 0x16, 0x33, 0xc5, 0xc8, 0x02, 0x8e, 0x72, 0x26, 0x55, 0x3b, 0x51, 0x50,
	0xa1, 0x90, 0x99, 0x54, 0x68, 0x9c, 0x7d, 0x3a, 0x6b, 0x2e, 0xcd, 0x5c, 0x37, 0xdd, 0x15, 0xf1,
	0xe0, 0xd8, 0x78, 0xb6, 0xe3, 0x08, 0x2a, 0xc1, 0x23, 0x94, 0x12, 0x63, 0x1d, 0x7b, 0x9f, 0xbe,
	0xd2, 0xe6, 0xad, 0x70, 0x6e, 0x3a, 0xe4, 0xff, 0x23, 0x02, 0x65, 0x1d, 0x16, 0x99, 0x52, 0x18,
	0x07, 0xed, 0xe2, 0xda, 0x74, 0xfb, 0x0f, 0x8f, 0xd0, 0x07, 0xc8, 0xd3, 0x8c, 0x79, 0xed, 0xe2,
	0x16, 0xde, 0x71, 0x91, 0x3a, 0x77, 0x9b, 0x0a, 0x45, 0x8e, 0x71, 0x8a, 0xc2, 0x49, 0x58, 0x28,
	0xb2, 0xc8, 0xfc, 0x16, 0xb2, 0xdb, 0xee, 0xaf, 0xb3, 0x34, 0x53, 0x77, 0x75, 0xe8, 0x44, 0xbc,
	0x70, 0x1f, 0xd1, 0xae, 0xa1, 0x5d, 0x43, 0xbb, 0x2d, 0x1d, 0x0e, 0xf5, 0xf9, 0xc3, 0xbf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x41, 0xa5, 0x96, 0xb1, 0x6b, 0x03, 0x00, 0x00,
}
