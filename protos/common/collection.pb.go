



package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 




type CollectionConfigPackage struct {
	Config []*CollectionConfig `protobuf:"bytes,1,rep,name=config" json:"config,omitempty"`
}

func (m *CollectionConfigPackage) Reset()                    { *m = CollectionConfigPackage{} }
func (m *CollectionConfigPackage) String() string            { return proto.CompactTextString(m) }
func (*CollectionConfigPackage) ProtoMessage()               {}
func (*CollectionConfigPackage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CollectionConfigPackage) GetConfig() []*CollectionConfig {
	if m != nil {
		return m.Config
	}
	return nil
}




type CollectionConfig struct {
	
	
	Payload isCollectionConfig_Payload `protobuf_oneof:"payload"`
}

func (m *CollectionConfig) Reset()                    { *m = CollectionConfig{} }
func (m *CollectionConfig) String() string            { return proto.CompactTextString(m) }
func (*CollectionConfig) ProtoMessage()               {}
func (*CollectionConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isCollectionConfig_Payload interface{ isCollectionConfig_Payload() }

type CollectionConfig_StaticCollectionConfig struct {
	StaticCollectionConfig *StaticCollectionConfig `protobuf:"bytes,1,opt,name=static_collection_config,json=staticCollectionConfig,oneof"`
}

func (*CollectionConfig_StaticCollectionConfig) isCollectionConfig_Payload() {}

func (m *CollectionConfig) GetPayload() isCollectionConfig_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *CollectionConfig) GetStaticCollectionConfig() *StaticCollectionConfig {
	if x, ok := m.GetPayload().(*CollectionConfig_StaticCollectionConfig); ok {
		return x.StaticCollectionConfig
	}
	return nil
}


func (*CollectionConfig) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CollectionConfig_OneofMarshaler, _CollectionConfig_OneofUnmarshaler, _CollectionConfig_OneofSizer, []interface{}{
		(*CollectionConfig_StaticCollectionConfig)(nil),
	}
}

func _CollectionConfig_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*CollectionConfig)
	
	switch x := m.Payload.(type) {
	case *CollectionConfig_StaticCollectionConfig:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StaticCollectionConfig); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("CollectionConfig.Payload has unexpected type %T", x)
	}
	return nil
}

func _CollectionConfig_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*CollectionConfig)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(StaticCollectionConfig)
		err := b.DecodeMessage(msg)
		m.Payload = &CollectionConfig_StaticCollectionConfig{msg}
		return true, err
	default:
		return false, nil
	}
}

func _CollectionConfig_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*CollectionConfig)
	
	switch x := m.Payload.(type) {
	case *CollectionConfig_StaticCollectionConfig:
		s := proto.Size(x.StaticCollectionConfig)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}





type StaticCollectionConfig struct {
	
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	
	
	MemberOrgsPolicy *CollectionPolicyConfig `protobuf:"bytes,2,opt,name=member_orgs_policy,json=memberOrgsPolicy" json:"member_orgs_policy,omitempty"`
	
	
	
	RequiredPeerCount int32 `protobuf:"varint,3,opt,name=required_peer_count,json=requiredPeerCount" json:"required_peer_count,omitempty"`
	
	
	MaximumPeerCount int32 `protobuf:"varint,4,opt,name=maximum_peer_count,json=maximumPeerCount" json:"maximum_peer_count,omitempty"`
	
	
	
	BlockToLive uint64 `protobuf:"varint,5,opt,name=block_to_live,json=blockToLive" json:"block_to_live,omitempty"`
}

func (m *StaticCollectionConfig) Reset()                    { *m = StaticCollectionConfig{} }
func (m *StaticCollectionConfig) String() string            { return proto.CompactTextString(m) }
func (*StaticCollectionConfig) ProtoMessage()               {}
func (*StaticCollectionConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StaticCollectionConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StaticCollectionConfig) GetMemberOrgsPolicy() *CollectionPolicyConfig {
	if m != nil {
		return m.MemberOrgsPolicy
	}
	return nil
}

func (m *StaticCollectionConfig) GetRequiredPeerCount() int32 {
	if m != nil {
		return m.RequiredPeerCount
	}
	return 0
}

func (m *StaticCollectionConfig) GetMaximumPeerCount() int32 {
	if m != nil {
		return m.MaximumPeerCount
	}
	return 0
}

func (m *StaticCollectionConfig) GetBlockToLive() uint64 {
	if m != nil {
		return m.BlockToLive
	}
	return 0
}





type CollectionPolicyConfig struct {
	
	
	Payload isCollectionPolicyConfig_Payload `protobuf_oneof:"payload"`
}

func (m *CollectionPolicyConfig) Reset()                    { *m = CollectionPolicyConfig{} }
func (m *CollectionPolicyConfig) String() string            { return proto.CompactTextString(m) }
func (*CollectionPolicyConfig) ProtoMessage()               {}
func (*CollectionPolicyConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type isCollectionPolicyConfig_Payload interface{ isCollectionPolicyConfig_Payload() }

type CollectionPolicyConfig_SignaturePolicy struct {
	SignaturePolicy *SignaturePolicyEnvelope `protobuf:"bytes,1,opt,name=signature_policy,json=signaturePolicy,oneof"`
}

func (*CollectionPolicyConfig_SignaturePolicy) isCollectionPolicyConfig_Payload() {}

func (m *CollectionPolicyConfig) GetPayload() isCollectionPolicyConfig_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *CollectionPolicyConfig) GetSignaturePolicy() *SignaturePolicyEnvelope {
	if x, ok := m.GetPayload().(*CollectionPolicyConfig_SignaturePolicy); ok {
		return x.SignaturePolicy
	}
	return nil
}


func (*CollectionPolicyConfig) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CollectionPolicyConfig_OneofMarshaler, _CollectionPolicyConfig_OneofUnmarshaler, _CollectionPolicyConfig_OneofSizer, []interface{}{
		(*CollectionPolicyConfig_SignaturePolicy)(nil),
	}
}

func _CollectionPolicyConfig_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*CollectionPolicyConfig)
	
	switch x := m.Payload.(type) {
	case *CollectionPolicyConfig_SignaturePolicy:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SignaturePolicy); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("CollectionPolicyConfig.Payload has unexpected type %T", x)
	}
	return nil
}

func _CollectionPolicyConfig_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*CollectionPolicyConfig)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SignaturePolicyEnvelope)
		err := b.DecodeMessage(msg)
		m.Payload = &CollectionPolicyConfig_SignaturePolicy{msg}
		return true, err
	default:
		return false, nil
	}
}

func _CollectionPolicyConfig_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*CollectionPolicyConfig)
	
	switch x := m.Payload.(type) {
	case *CollectionPolicyConfig_SignaturePolicy:
		s := proto.Size(x.SignaturePolicy)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}



type CollectionCriteria struct {
	Channel    string `protobuf:"bytes,1,opt,name=channel" json:"channel,omitempty"`
	TxId       string `protobuf:"bytes,2,opt,name=tx_id,json=txId" json:"tx_id,omitempty"`
	Collection string `protobuf:"bytes,3,opt,name=collection" json:"collection,omitempty"`
	Namespace  string `protobuf:"bytes,4,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *CollectionCriteria) Reset()                    { *m = CollectionCriteria{} }
func (m *CollectionCriteria) String() string            { return proto.CompactTextString(m) }
func (*CollectionCriteria) ProtoMessage()               {}
func (*CollectionCriteria) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *CollectionCriteria) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *CollectionCriteria) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *CollectionCriteria) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

func (m *CollectionCriteria) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func init() {
	proto.RegisterType((*CollectionConfigPackage)(nil), "common.CollectionConfigPackage")
	proto.RegisterType((*CollectionConfig)(nil), "common.CollectionConfig")
	proto.RegisterType((*StaticCollectionConfig)(nil), "common.StaticCollectionConfig")
	proto.RegisterType((*CollectionPolicyConfig)(nil), "common.CollectionPolicyConfig")
	proto.RegisterType((*CollectionCriteria)(nil), "common.CollectionCriteria")
}

func init() { proto.RegisterFile("common/collection.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x41, 0x6b, 0xdb, 0x40,
	0x10, 0x85, 0xa3, 0xc6, 0x76, 0xd0, 0x98, 0x52, 0x77, 0x43, 0x1d, 0x51, 0x4a, 0x6a, 0x44, 0x0f,
	0x86, 0x16, 0xa9, 0xa4, 0xff, 0x20, 0xa6, 0x90, 0x52, 0x43, 0x8d, 0xd2, 0x53, 0x2e, 0x62, 0xb5,
	0x9a, 0xc8, 0x4b, 0x24, 0xad, 0xb2, 0xbb, 0x32, 0xf6, 0xb1, 0xff, 0xbb, 0x87, 0xe0, 0x5d, 0xc9,
	0x52, 0x8c, 0x6f, 0x9e, 0x79, 0xdf, 0x3c, 0xcf, 0x3c, 0x2d, 0x5c, 0x31, 0x51, 0x14, 0xa2, 0x0c,
	0x99, 0xc8, 0x73, 0x64, 0x9a, 0x8b, 0x32, 0xa8, 0xa4, 0xd0, 0x82, 0x8c, 0xac, 0xf0, 0xf1, 0x43,
	0x03, 0x54, 0x22, 0xe7, 0x8c, 0xa3, 0xb2, 0xb2, 0xff, 0x1b, 0xae, 0x16, 0x87, 0x91, 0x85, 0x28,
	0x1f, 0x79, 0xb6, 0xa2, 0xec, 0x89, 0x66, 0x48, 0xbe, 0xc3, 0x88, 0x99, 0x86, 0xe7, 0xcc, 0xce,
	0xe7, 0xe3, 0x1b, 0x2f, 0xb0, 0x16, 0xc1, 0xf1, 0x40, 0xd4, 0x70, 0xfe, 0x0e, 0x26, 0xc7, 0x1a,
	0x79, 0x00, 0x4f, 0x69, 0xaa, 0x39, 0x8b, 0xbb, 0xd5, 0xe2, 0x83, 0xaf, 0x33, 0x1f, 0xdf, 0x5c,
	0xb7, 0xbe, 0xf7, 0x86, 0x3b, 0x76, 0xb8, 0x3b, 0x8b, 0xa6, 0xea, 0xa4, 0x72, 0xeb, 0xc2, 0x45,
	0x45, 0x77, 0xb9, 0xa0, 0xa9, 0xff, 0xdf, 0x81, 0xe9, 0xe9, 0x79, 0x42, 0x60, 0x50, 0xd2, 0x02,
	0xcd, 0xbf, 0xb9, 0x91, 0xf9, 0x4d, 0x96, 0x40, 0x0a, 0x2c, 0x12, 0x94, 0xb1, 0x90, 0x99, 0x8a,
	0x4d, 0x28, 0x3b, 0xef, 0xcd, 0xeb, 0x7d, 0x3a, 0xa7, 0x95, 0xd1, 0x9b, 0x6b, 0x27, 0x76, 0xf2,
	0x8f, 0xcc, 0x94, 0xed, 0x93, 0x00, 0x2e, 0x25, 0x3e, 0xd7, 0x5c, 0x62, 0x1a, 0x57, 0x88, 0x32,
	0x66, 0xa2, 0x2e, 0xb5, 0x77, 0x3e, 0x73, 0xe6, 0xc3, 0xe8, 0x7d, 0x2b, 0xad, 0x10, 0xe5, 0x62,
	0x2f, 0x90, 0x6f, 0x40, 0x0a, 0xba, 0xe5, 0x45, 0x5d, 0xf4, 0xf1, 0x81, 0xc1, 0x27, 0x8d, 0xd2,
	0xd1, 0x3e, 0xbc, 0x4d, 0x72, 0xc1, 0x9e, 0x62, 0x2d, 0xe2, 0x9c, 0x6f, 0xd0, 0x1b, 0xce, 0x9c,
	0xf9, 0x20, 0x1a, 0x9b, 0xe6, 0x5f, 0xb1, 0xe4, 0x1b, 0xf4, 0x9f, 0x61, 0x7a, 0x7a, 0x5b, 0xb2,
	0x84, 0x89, 0xe2, 0x59, 0x49, 0x75, 0x2d, 0xb1, 0xbd, 0xd3, 0xe6, 0xfe, 0xf9, 0x90, 0x7b, 0xab,
	0xdb, 0xc1, 0x9f, 0xe5, 0x06, 0x73, 0x51, 0xe1, 0xdd, 0x59, 0xf4, 0x4e, 0xbd, 0x96, 0xfa, 0x89,
	0xff, 0x73, 0x80, 0xf4, 0xb2, 0x96, 0x5c, 0xa3, 0xe4, 0x94, 0x78, 0x70, 0xc1, 0xd6, 0xb4, 0x2c,
	0x31, 0x6f, 0x02, 0x6f, 0x4b, 0x72, 0x09, 0x43, 0xbd, 0x8d, 0x79, 0x6a, 0x62, 0x76, 0xa3, 0x81,
	0xde, 0xfe, 0x4a, 0xc9, 0x35, 0x40, 0xf7, 0x2e, 0x4c, 0x62, 0x6e, 0xd4, 0xeb, 0x90, 0x4f, 0xe0,
	0xee, 0x3f, 0x98, 0xaa, 0x28, 0x43, 0x93, 0x90, 0x1b, 0x75, 0x8d, 0xdb, 0x7b, 0xf8, 0x22, 0x64,
	0x16, 0xac, 0x77, 0x15, 0xca, 0x1c, 0xd3, 0x0c, 0x65, 0xf0, 0x48, 0x13, 0xc9, 0x99, 0x7d, 0xdd,
	0xaa, 0xb9, 0xf0, 0xe1, 0x6b, 0xc6, 0xf5, 0xba, 0x4e, 0xf6, 0x65, 0xd8, 0x83, 0x43, 0x0b, 0x87,
	0x16, 0x0e, 0x2d, 0x9c, 0x8c, 0x4c, 0xf9, 0xe3, 0x25, 0x00, 0x00, 0xff, 0xff, 0x04, 0x6f, 0x60,
	0x95, 0x53, 0x03, 0x00, 0x00,
}
