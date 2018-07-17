


package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common1 "github.com/mcc-github/blockchain/protos/msp"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Policy_PolicyType int32

const (
	Policy_UNKNOWN       Policy_PolicyType = 0
	Policy_SIGNATURE     Policy_PolicyType = 1
	Policy_MSP           Policy_PolicyType = 2
	Policy_IMPLICIT_META Policy_PolicyType = 3
)

var Policy_PolicyType_name = map[int32]string{
	0: "UNKNOWN",
	1: "SIGNATURE",
	2: "MSP",
	3: "IMPLICIT_META",
}
var Policy_PolicyType_value = map[string]int32{
	"UNKNOWN":       0,
	"SIGNATURE":     1,
	"MSP":           2,
	"IMPLICIT_META": 3,
}

func (x Policy_PolicyType) String() string {
	return proto.EnumName(Policy_PolicyType_name, int32(x))
}
func (Policy_PolicyType) EnumDescriptor() ([]byte, []int) { return fileDescriptor5, []int{0, 0} }

type ImplicitMetaPolicy_Rule int32

const (
	ImplicitMetaPolicy_ANY      ImplicitMetaPolicy_Rule = 0
	ImplicitMetaPolicy_ALL      ImplicitMetaPolicy_Rule = 1
	ImplicitMetaPolicy_MAJORITY ImplicitMetaPolicy_Rule = 2
)

var ImplicitMetaPolicy_Rule_name = map[int32]string{
	0: "ANY",
	1: "ALL",
	2: "MAJORITY",
}
var ImplicitMetaPolicy_Rule_value = map[string]int32{
	"ANY":      0,
	"ALL":      1,
	"MAJORITY": 2,
}

func (x ImplicitMetaPolicy_Rule) String() string {
	return proto.EnumName(ImplicitMetaPolicy_Rule_name, int32(x))
}
func (ImplicitMetaPolicy_Rule) EnumDescriptor() ([]byte, []int) { return fileDescriptor5, []int{3, 0} }



type Policy struct {
	Type  int32  `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Policy) Reset()                    { *m = Policy{} }
func (m *Policy) String() string            { return proto.CompactTextString(m) }
func (*Policy) ProtoMessage()               {}
func (*Policy) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{0} }

func (m *Policy) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Policy) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}


type SignaturePolicyEnvelope struct {
	Version    int32                   `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Rule       *SignaturePolicy        `protobuf:"bytes,2,opt,name=rule" json:"rule,omitempty"`
	Identities []*common1.MSPPrincipal `protobuf:"bytes,3,rep,name=identities" json:"identities,omitempty"`
}

func (m *SignaturePolicyEnvelope) Reset()                    { *m = SignaturePolicyEnvelope{} }
func (m *SignaturePolicyEnvelope) String() string            { return proto.CompactTextString(m) }
func (*SignaturePolicyEnvelope) ProtoMessage()               {}
func (*SignaturePolicyEnvelope) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{1} }

func (m *SignaturePolicyEnvelope) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *SignaturePolicyEnvelope) GetRule() *SignaturePolicy {
	if m != nil {
		return m.Rule
	}
	return nil
}

func (m *SignaturePolicyEnvelope) GetIdentities() []*common1.MSPPrincipal {
	if m != nil {
		return m.Identities
	}
	return nil
}







type SignaturePolicy struct {
	
	
	
	Type isSignaturePolicy_Type `protobuf_oneof:"Type"`
}

func (m *SignaturePolicy) Reset()                    { *m = SignaturePolicy{} }
func (m *SignaturePolicy) String() string            { return proto.CompactTextString(m) }
func (*SignaturePolicy) ProtoMessage()               {}
func (*SignaturePolicy) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{2} }

type isSignaturePolicy_Type interface{ isSignaturePolicy_Type() }

type SignaturePolicy_SignedBy struct {
	SignedBy int32 `protobuf:"varint,1,opt,name=signed_by,json=signedBy,oneof"`
}
type SignaturePolicy_NOutOf_ struct {
	NOutOf *SignaturePolicy_NOutOf `protobuf:"bytes,2,opt,name=n_out_of,json=nOutOf,oneof"`
}

func (*SignaturePolicy_SignedBy) isSignaturePolicy_Type() {}
func (*SignaturePolicy_NOutOf_) isSignaturePolicy_Type()  {}

func (m *SignaturePolicy) GetType() isSignaturePolicy_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *SignaturePolicy) GetSignedBy() int32 {
	if x, ok := m.GetType().(*SignaturePolicy_SignedBy); ok {
		return x.SignedBy
	}
	return 0
}

func (m *SignaturePolicy) GetNOutOf() *SignaturePolicy_NOutOf {
	if x, ok := m.GetType().(*SignaturePolicy_NOutOf_); ok {
		return x.NOutOf
	}
	return nil
}


func (*SignaturePolicy) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _SignaturePolicy_OneofMarshaler, _SignaturePolicy_OneofUnmarshaler, _SignaturePolicy_OneofSizer, []interface{}{
		(*SignaturePolicy_SignedBy)(nil),
		(*SignaturePolicy_NOutOf_)(nil),
	}
}

func _SignaturePolicy_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*SignaturePolicy)
	
	switch x := m.Type.(type) {
	case *SignaturePolicy_SignedBy:
		b.EncodeVarint(1<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.SignedBy))
	case *SignaturePolicy_NOutOf_:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.NOutOf); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("SignaturePolicy.Type has unexpected type %T", x)
	}
	return nil
}

func _SignaturePolicy_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*SignaturePolicy)
	switch tag {
	case 1: 
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Type = &SignaturePolicy_SignedBy{int32(x)}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SignaturePolicy_NOutOf)
		err := b.DecodeMessage(msg)
		m.Type = &SignaturePolicy_NOutOf_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _SignaturePolicy_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*SignaturePolicy)
	
	switch x := m.Type.(type) {
	case *SignaturePolicy_SignedBy:
		n += proto.SizeVarint(1<<3 | proto.WireVarint)
		n += proto.SizeVarint(uint64(x.SignedBy))
	case *SignaturePolicy_NOutOf_:
		s := proto.Size(x.NOutOf)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type SignaturePolicy_NOutOf struct {
	N     int32              `protobuf:"varint,1,opt,name=n" json:"n,omitempty"`
	Rules []*SignaturePolicy `protobuf:"bytes,2,rep,name=rules" json:"rules,omitempty"`
}

func (m *SignaturePolicy_NOutOf) Reset()                    { *m = SignaturePolicy_NOutOf{} }
func (m *SignaturePolicy_NOutOf) String() string            { return proto.CompactTextString(m) }
func (*SignaturePolicy_NOutOf) ProtoMessage()               {}
func (*SignaturePolicy_NOutOf) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{2, 0} }

func (m *SignaturePolicy_NOutOf) GetN() int32 {
	if m != nil {
		return m.N
	}
	return 0
}

func (m *SignaturePolicy_NOutOf) GetRules() []*SignaturePolicy {
	if m != nil {
		return m.Rules
	}
	return nil
}









type ImplicitMetaPolicy struct {
	SubPolicy string                  `protobuf:"bytes,1,opt,name=sub_policy,json=subPolicy" json:"sub_policy,omitempty"`
	Rule      ImplicitMetaPolicy_Rule `protobuf:"varint,2,opt,name=rule,enum=common.ImplicitMetaPolicy_Rule" json:"rule,omitempty"`
}

func (m *ImplicitMetaPolicy) Reset()                    { *m = ImplicitMetaPolicy{} }
func (m *ImplicitMetaPolicy) String() string            { return proto.CompactTextString(m) }
func (*ImplicitMetaPolicy) ProtoMessage()               {}
func (*ImplicitMetaPolicy) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{3} }

func (m *ImplicitMetaPolicy) GetSubPolicy() string {
	if m != nil {
		return m.SubPolicy
	}
	return ""
}

func (m *ImplicitMetaPolicy) GetRule() ImplicitMetaPolicy_Rule {
	if m != nil {
		return m.Rule
	}
	return ImplicitMetaPolicy_ANY
}

func init() {
	proto.RegisterType((*Policy)(nil), "common.Policy")
	proto.RegisterType((*SignaturePolicyEnvelope)(nil), "common.SignaturePolicyEnvelope")
	proto.RegisterType((*SignaturePolicy)(nil), "common.SignaturePolicy")
	proto.RegisterType((*SignaturePolicy_NOutOf)(nil), "common.SignaturePolicy.NOutOf")
	proto.RegisterType((*ImplicitMetaPolicy)(nil), "common.ImplicitMetaPolicy")
	proto.RegisterEnum("common.Policy_PolicyType", Policy_PolicyType_name, Policy_PolicyType_value)
	proto.RegisterEnum("common.ImplicitMetaPolicy_Rule", ImplicitMetaPolicy_Rule_name, ImplicitMetaPolicy_Rule_value)
}

func init() { proto.RegisterFile("common/policies.proto", fileDescriptor5) }

var fileDescriptor5 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0xdf, 0x8b, 0xda, 0x40,
	0x10, 0x76, 0xfd, 0x11, 0x75, 0xf4, 0xda, 0x74, 0xb9, 0xa2, 0x1c, 0xb4, 0x95, 0x50, 0x8a, 0x70,
	0x34, 0x01, 0xaf, 0x4f, 0x7d, 0xd3, 0x56, 0x7a, 0x69, 0x4d, 0x94, 0xd5, 0xa3, 0x5c, 0x5f, 0x82,
	0xd1, 0xd5, 0x5b, 0x88, 0xbb, 0x4b, 0x76, 0x23, 0xf5, 0xbf, 0xe8, 0x53, 0xff, 0x99, 0xfe, 0x73,
	0x25, 0x59, 0x2d, 0x72, 0xe5, 0xde, 0xe6, 0x9b, 0xfd, 0x66, 0xf6, 0xfb, 0x66, 0x06, 0x5e, 0xae,
	0xc4, 0x6e, 0x27, 0xb8, 0x27, 0x45, 0xc2, 0x56, 0x8c, 0x2a, 0x57, 0xa6, 0x42, 0x0b, 0x6c, 0x99,
	0xf4, 0x55, 0x67, 0xa7, 0xa4, 0xb7, 0x53, 0x32, 0x92, 0x29, 0xe3, 0x2b, 0x26, 0x97, 0x89, 0x21,
	0x38, 0x3f, 0xc1, 0x9a, 0xe5, 0x25, 0x07, 0x8c, 0xa1, 0xaa, 0x0f, 0x92, 0x76, 0x51, 0x0f, 0xf5,
	0x6b, 0xa4, 0x88, 0xf1, 0x25, 0xd4, 0xf6, 0xcb, 0x24, 0xa3, 0xdd, 0x72, 0x0f, 0xf5, 0xdb, 0xc4,
	0x00, 0xe7, 0x33, 0x80, 0xa9, 0x59, 0xe4, 0x9c, 0x16, 0xd4, 0xef, 0xc2, 0x6f, 0xe1, 0xf4, 0x7b,
	0x68, 0x97, 0xf0, 0x05, 0x34, 0xe7, 0xfe, 0x97, 0x70, 0xb8, 0xb8, 0x23, 0x63, 0x1b, 0xe1, 0x3a,
	0x54, 0x82, 0xf9, 0xcc, 0x2e, 0xe3, 0x17, 0x70, 0xe1, 0x07, 0xb3, 0x89, 0xff, 0xc9, 0x5f, 0x44,
	0xc1, 0x78, 0x31, 0xb4, 0x2b, 0xce, 0x6f, 0x04, 0x9d, 0x39, 0xdb, 0xf2, 0xa5, 0xce, 0x52, 0x6a,
	0xfa, 0x8d, 0xf9, 0x9e, 0x26, 0x42, 0x52, 0xdc, 0x85, 0xfa, 0x9e, 0xa6, 0x8a, 0x09, 0x7e, 0x94,
	0x73, 0x82, 0xf8, 0x1a, 0xaa, 0x69, 0x96, 0x18, 0x41, 0xad, 0x41, 0xc7, 0x35, 0xfe, 0xdc, 0x47,
	0x8d, 0x48, 0x41, 0xc2, 0x1f, 0x00, 0xd8, 0x9a, 0x72, 0xcd, 0x34, 0xa3, 0xaa, 0x5b, 0xe9, 0x55,
	0xfa, 0xad, 0xc1, 0xe5, 0xa9, 0x24, 0x98, 0xcf, 0x66, 0xa7, 0x61, 0x90, 0x33, 0x9e, 0xf3, 0x07,
	0xc1, 0xf3, 0x47, 0xfd, 0xf0, 0x2b, 0x68, 0x2a, 0xb6, 0xe5, 0x74, 0x1d, 0xc5, 0x07, 0x23, 0xe9,
	0xb6, 0x44, 0x1a, 0x26, 0x35, 0x3a, 0xe0, 0x8f, 0xd0, 0xe0, 0x91, 0xc8, 0x74, 0x24, 0x36, 0x47,
	0x65, 0xaf, 0x9f, 0x50, 0xe6, 0x86, 0xd3, 0x4c, 0x4f, 0x37, 0xb7, 0x25, 0x62, 0xf1, 0x22, 0xba,
	0x1a, 0x83, 0x65, 0x72, 0xb8, 0x0d, 0xe8, 0xe4, 0x17, 0x71, 0xfc, 0x1e, 0x6a, 0xb9, 0x09, 0xd5,
	0x2d, 0x17, 0xba, 0x9f, 0xb4, 0x6a, 0x58, 0x23, 0x0b, 0xaa, 0xf9, 0x3a, 0x9c, 0x5f, 0x08, 0xb0,
	0xbf, 0x93, 0xf9, 0x15, 0xe8, 0x80, 0xea, 0xe5, 0x3f, 0x03, 0xa0, 0xb2, 0x38, 0x2a, 0xce, 0xc3,
	0x38, 0x68, 0x92, 0xa6, 0xca, 0xe2, 0xe3, 0xf3, 0xcd, 0xd9, 0x58, 0x9f, 0x0d, 0xde, 0x9c, 0xfe,
	0xfa, 0xbf, 0x91, 0x4b, 0xb2, 0x84, 0x9a, 0xf1, 0x3a, 0xef, 0xa0, 0x9a, 0xa3, 0x7c, 0xcb, 0xc3,
	0xf0, 0xde, 0x2e, 0x15, 0xc1, 0x64, 0x62, 0x23, 0xdc, 0x86, 0x46, 0x30, 0xfc, 0x3a, 0x25, 0xfe,
	0xe2, 0xde, 0x2e, 0x8f, 0xe6, 0xf0, 0x56, 0xa4, 0x5b, 0xf7, 0xe1, 0x20, 0x69, 0x9a, 0xd0, 0xf5,
	0x96, 0xa6, 0xee, 0x66, 0x19, 0xa7, 0x6c, 0x65, 0x6e, 0x50, 0x1d, 0x7f, 0xfb, 0x71, 0xbd, 0x65,
	0xfa, 0x21, 0x8b, 0x73, 0xe8, 0x9d, 0x91, 0x3d, 0x43, 0xf6, 0x0c, 0xd9, 0x33, 0xe4, 0xd8, 0x2a,
	0xe0, 0xcd, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x34, 0xab, 0x8a, 0xb5, 0xf9, 0x02, 0x00, 0x00,
}
