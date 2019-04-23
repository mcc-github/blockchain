


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "github.com/mcc-github/blockchain/protos/common"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



type ApplicationPolicy struct {
	
	
	
	Type                 isApplicationPolicy_Type `protobuf_oneof:"Type"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ApplicationPolicy) Reset()         { *m = ApplicationPolicy{} }
func (m *ApplicationPolicy) String() string { return proto.CompactTextString(m) }
func (*ApplicationPolicy) ProtoMessage()    {}
func (*ApplicationPolicy) Descriptor() ([]byte, []int) {
	return fileDescriptor_policy_cb2dcc38242d4048, []int{0}
}
func (m *ApplicationPolicy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplicationPolicy.Unmarshal(m, b)
}
func (m *ApplicationPolicy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplicationPolicy.Marshal(b, m, deterministic)
}
func (dst *ApplicationPolicy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplicationPolicy.Merge(dst, src)
}
func (m *ApplicationPolicy) XXX_Size() int {
	return xxx_messageInfo_ApplicationPolicy.Size(m)
}
func (m *ApplicationPolicy) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplicationPolicy.DiscardUnknown(m)
}

var xxx_messageInfo_ApplicationPolicy proto.InternalMessageInfo

type isApplicationPolicy_Type interface {
	isApplicationPolicy_Type()
}

type ApplicationPolicy_SignaturePolicy struct {
	SignaturePolicy *common.SignaturePolicyEnvelope `protobuf:"bytes,1,opt,name=signature_policy,json=signaturePolicy,proto3,oneof"`
}

type ApplicationPolicy_ChannelConfigPolicyReference struct {
	ChannelConfigPolicyReference string `protobuf:"bytes,2,opt,name=channel_config_policy_reference,json=channelConfigPolicyReference,proto3,oneof"`
}

func (*ApplicationPolicy_SignaturePolicy) isApplicationPolicy_Type() {}

func (*ApplicationPolicy_ChannelConfigPolicyReference) isApplicationPolicy_Type() {}

func (m *ApplicationPolicy) GetType() isApplicationPolicy_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *ApplicationPolicy) GetSignaturePolicy() *common.SignaturePolicyEnvelope {
	if x, ok := m.GetType().(*ApplicationPolicy_SignaturePolicy); ok {
		return x.SignaturePolicy
	}
	return nil
}

func (m *ApplicationPolicy) GetChannelConfigPolicyReference() string {
	if x, ok := m.GetType().(*ApplicationPolicy_ChannelConfigPolicyReference); ok {
		return x.ChannelConfigPolicyReference
	}
	return ""
}


func (*ApplicationPolicy) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ApplicationPolicy_OneofMarshaler, _ApplicationPolicy_OneofUnmarshaler, _ApplicationPolicy_OneofSizer, []interface{}{
		(*ApplicationPolicy_SignaturePolicy)(nil),
		(*ApplicationPolicy_ChannelConfigPolicyReference)(nil),
	}
}

func _ApplicationPolicy_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ApplicationPolicy)
	
	switch x := m.Type.(type) {
	case *ApplicationPolicy_SignaturePolicy:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SignaturePolicy); err != nil {
			return err
		}
	case *ApplicationPolicy_ChannelConfigPolicyReference:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.ChannelConfigPolicyReference)
	case nil:
	default:
		return fmt.Errorf("ApplicationPolicy.Type has unexpected type %T", x)
	}
	return nil
}

func _ApplicationPolicy_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ApplicationPolicy)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(common.SignaturePolicyEnvelope)
		err := b.DecodeMessage(msg)
		m.Type = &ApplicationPolicy_SignaturePolicy{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Type = &ApplicationPolicy_ChannelConfigPolicyReference{x}
		return true, err
	default:
		return false, nil
	}
}

func _ApplicationPolicy_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ApplicationPolicy)
	
	switch x := m.Type.(type) {
	case *ApplicationPolicy_SignaturePolicy:
		s := proto.Size(x.SignaturePolicy)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ApplicationPolicy_ChannelConfigPolicyReference:
		n += 1 
		n += proto.SizeVarint(uint64(len(x.ChannelConfigPolicyReference)))
		n += len(x.ChannelConfigPolicyReference)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ApplicationPolicy)(nil), "protos.ApplicationPolicy")
}

func init() { proto.RegisterFile("peer/policy.proto", fileDescriptor_policy_cb2dcc38242d4048) }

var fileDescriptor_policy_cb2dcc38242d4048 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x1b, 0x91, 0x82, 0xeb, 0x41, 0x1b, 0x10, 0x8a, 0x08, 0x2d, 0x3d, 0xd5, 0xcb, 0x2e,
	0xe8, 0x13, 0x58, 0x11, 0x7b, 0x10, 0x94, 0xe8, 0xc9, 0x4b, 0x48, 0xd6, 0xc9, 0x66, 0x61, 0xbb,
	0x33, 0xcc, 0xa6, 0x42, 0x5e, 0xcb, 0x27, 0x94, 0x64, 0x5a, 0xd0, 0xd3, 0x1e, 0xbe, 0xef, 0xff,
	0xd9, 0xf9, 0xd5, 0x8c, 0x00, 0xd8, 0x10, 0x06, 0x6f, 0x7b, 0x4d, 0x8c, 0x1d, 0xe6, 0xd3, 0xf1,
	0x49, 0xd7, 0x57, 0x16, 0x77, 0x3b, 0x8c, 0x02, 0x3d, 0x24, 0xc1, 0xab, 0x9f, 0x4c, 0xcd, 0x1e,
	0x88, 0x82, 0xb7, 0x55, 0xe7, 0x31, 0xbe, 0x8d, 0xd1, 0xfc, 0x45, 0x5d, 0x26, 0xef, 0x62, 0xd5,
	0xed, 0x19, 0x4a, 0xa9, 0x9b, 0x67, 0xcb, 0x6c, 0x7d, 0x7e, 0xb7, 0xd0, 0xd2, 0xa3, 0xdf, 0x8f,
	0x5c, 0x22, 0x4f, 0xf1, 0x1b, 0x02, 0x12, 0x6c, 0x27, 0xc5, 0x45, 0xfa, 0x8f, 0xf2, 0x67, 0xb5,
	0xb0, 0x6d, 0x15, 0x23, 0x84, 0xd2, 0x62, 0x6c, 0xbc, 0x3b, 0x54, 0x96, 0x0c, 0x0d, 0x30, 0x44,
	0x0b, 0xf3, 0x93, 0x65, 0xb6, 0x3e, 0xdb, 0x4e, 0x8a, 0x9b, 0x83, 0xf8, 0x38, 0x7a, 0x92, 0x2f,
	0x8e, 0xd6, 0x66, 0xaa, 0x4e, 0x3f, 0x7a, 0x82, 0xcd, 0xab, 0x5a, 0x21, 0x3b, 0xdd, 0xf6, 0x04,
	0x1c, 0xe0, 0xcb, 0x01, 0xeb, 0xa6, 0xaa, 0xd9, 0x5b, 0x39, 0x2a, 0xe9, 0x61, 0x86, 0xcf, 0x5b,
	0xe7, 0xbb, 0x76, 0x5f, 0x0f, 0x1f, 0x36, 0x7f, 0x54, 0x23, 0xaa, 0x11, 0xd5, 0x0c, 0x6a, 0x2d,
	0x23, 0xdd, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xd3, 0x7d, 0xd7, 0x44, 0x40, 0x01, 0x00, 0x00,
}
