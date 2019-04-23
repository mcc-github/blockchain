


package token 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type TokenOwner_Type int32

const (
	TokenOwner_MSP_IDENTIFIER TokenOwner_Type = 0
)

var TokenOwner_Type_name = map[int32]string{
	0: "MSP_IDENTIFIER",
}
var TokenOwner_Type_value = map[string]int32{
	"MSP_IDENTIFIER": 0,
}

func (x TokenOwner_Type) String() string {
	return proto.EnumName(TokenOwner_Type_name, int32(x))
}
func (TokenOwner_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f976d1453a3f05a8, []int{2, 0}
}




type TokenTransaction struct {
	
	
	
	
	Action               isTokenTransaction_Action `protobuf_oneof:"action"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *TokenTransaction) Reset()         { *m = TokenTransaction{} }
func (m *TokenTransaction) String() string { return proto.CompactTextString(m) }
func (*TokenTransaction) ProtoMessage()    {}
func (*TokenTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f976d1453a3f05a8, []int{0}
}
func (m *TokenTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenTransaction.Unmarshal(m, b)
}
func (m *TokenTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenTransaction.Marshal(b, m, deterministic)
}
func (dst *TokenTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenTransaction.Merge(dst, src)
}
func (m *TokenTransaction) XXX_Size() int {
	return xxx_messageInfo_TokenTransaction.Size(m)
}
func (m *TokenTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_TokenTransaction proto.InternalMessageInfo

type isTokenTransaction_Action interface {
	isTokenTransaction_Action()
}

type TokenTransaction_TokenAction struct {
	TokenAction *TokenAction `protobuf:"bytes,1,opt,name=token_action,json=tokenAction,proto3,oneof"`
}

func (*TokenTransaction_TokenAction) isTokenTransaction_Action() {}

func (m *TokenTransaction) GetAction() isTokenTransaction_Action {
	if m != nil {
		return m.Action
	}
	return nil
}

func (m *TokenTransaction) GetTokenAction() *TokenAction {
	if x, ok := m.GetAction().(*TokenTransaction_TokenAction); ok {
		return x.TokenAction
	}
	return nil
}


func (*TokenTransaction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TokenTransaction_OneofMarshaler, _TokenTransaction_OneofUnmarshaler, _TokenTransaction_OneofSizer, []interface{}{
		(*TokenTransaction_TokenAction)(nil),
	}
}

func _TokenTransaction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TokenTransaction)
	
	switch x := m.Action.(type) {
	case *TokenTransaction_TokenAction:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TokenAction); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("TokenTransaction.Action has unexpected type %T", x)
	}
	return nil
}

func _TokenTransaction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TokenTransaction)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TokenAction)
		err := b.DecodeMessage(msg)
		m.Action = &TokenTransaction_TokenAction{msg}
		return true, err
	default:
		return false, nil
	}
}

func _TokenTransaction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TokenTransaction)
	
	switch x := m.Action.(type) {
	case *TokenTransaction_TokenAction:
		s := proto.Size(x.TokenAction)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type TokenAction struct {
	
	
	
	
	
	
	Data                 isTokenAction_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *TokenAction) Reset()         { *m = TokenAction{} }
func (m *TokenAction) String() string { return proto.CompactTextString(m) }
func (*TokenAction) ProtoMessage()    {}
func (*TokenAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f976d1453a3f05a8, []int{1}
}
func (m *TokenAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenAction.Unmarshal(m, b)
}
func (m *TokenAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenAction.Marshal(b, m, deterministic)
}
func (dst *TokenAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenAction.Merge(dst, src)
}
func (m *TokenAction) XXX_Size() int {
	return xxx_messageInfo_TokenAction.Size(m)
}
func (m *TokenAction) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenAction.DiscardUnknown(m)
}

var xxx_messageInfo_TokenAction proto.InternalMessageInfo

type isTokenAction_Data interface {
	isTokenAction_Data()
}

type TokenAction_Issue struct {
	Issue *Issue `protobuf:"bytes,1,opt,name=issue,proto3,oneof"`
}

type TokenAction_Transfer struct {
	Transfer *Transfer `protobuf:"bytes,2,opt,name=transfer,proto3,oneof"`
}

type TokenAction_Redeem struct {
	Redeem *Transfer `protobuf:"bytes,3,opt,name=redeem,proto3,oneof"`
}

func (*TokenAction_Issue) isTokenAction_Data() {}

func (*TokenAction_Transfer) isTokenAction_Data() {}

func (*TokenAction_Redeem) isTokenAction_Data() {}

func (m *TokenAction) GetData() isTokenAction_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *TokenAction) GetIssue() *Issue {
	if x, ok := m.GetData().(*TokenAction_Issue); ok {
		return x.Issue
	}
	return nil
}

func (m *TokenAction) GetTransfer() *Transfer {
	if x, ok := m.GetData().(*TokenAction_Transfer); ok {
		return x.Transfer
	}
	return nil
}

func (m *TokenAction) GetRedeem() *Transfer {
	if x, ok := m.GetData().(*TokenAction_Redeem); ok {
		return x.Redeem
	}
	return nil
}


func (*TokenAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TokenAction_OneofMarshaler, _TokenAction_OneofUnmarshaler, _TokenAction_OneofSizer, []interface{}{
		(*TokenAction_Issue)(nil),
		(*TokenAction_Transfer)(nil),
		(*TokenAction_Redeem)(nil),
	}
}

func _TokenAction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TokenAction)
	
	switch x := m.Data.(type) {
	case *TokenAction_Issue:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Issue); err != nil {
			return err
		}
	case *TokenAction_Transfer:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Transfer); err != nil {
			return err
		}
	case *TokenAction_Redeem:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Redeem); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("TokenAction.Data has unexpected type %T", x)
	}
	return nil
}

func _TokenAction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TokenAction)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Issue)
		err := b.DecodeMessage(msg)
		m.Data = &TokenAction_Issue{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Transfer)
		err := b.DecodeMessage(msg)
		m.Data = &TokenAction_Transfer{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Transfer)
		err := b.DecodeMessage(msg)
		m.Data = &TokenAction_Redeem{msg}
		return true, err
	default:
		return false, nil
	}
}

func _TokenAction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TokenAction)
	
	switch x := m.Data.(type) {
	case *TokenAction_Issue:
		s := proto.Size(x.Issue)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TokenAction_Transfer:
		s := proto.Size(x.Transfer)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TokenAction_Redeem:
		s := proto.Size(x.Redeem)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type TokenOwner struct {
	
	Type TokenOwner_Type `protobuf:"varint,1,opt,name=type,proto3,enum=token.TokenOwner_Type" json:"type,omitempty"`
	
	Raw                  []byte   `protobuf:"bytes,2,opt,name=raw,proto3" json:"raw,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenOwner) Reset()         { *m = TokenOwner{} }
func (m *TokenOwner) String() string { return proto.CompactTextString(m) }
func (*TokenOwner) ProtoMessage()    {}
func (*TokenOwner) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f976d1453a3f05a8, []int{2}
}
func (m *TokenOwner) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenOwner.Unmarshal(m, b)
}
func (m *TokenOwner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenOwner.Marshal(b, m, deterministic)
}
func (dst *TokenOwner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenOwner.Merge(dst, src)
}
func (m *TokenOwner) XXX_Size() int {
	return xxx_messageInfo_TokenOwner.Size(m)
}
func (m *TokenOwner) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenOwner.DiscardUnknown(m)
}

var xxx_messageInfo_TokenOwner proto.InternalMessageInfo

func (m *TokenOwner) GetType() TokenOwner_Type {
	if m != nil {
		return m.Type
	}
	return TokenOwner_MSP_IDENTIFIER
}

func (m *TokenOwner) GetRaw() []byte {
	if m != nil {
		return m.Raw
	}
	return nil
}


type Issue struct {
	
	Outputs              []*Token `protobuf:"bytes,1,rep,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Issue) Reset()         { *m = Issue{} }
func (m *Issue) String() string { return proto.CompactTextString(m) }
func (*Issue) ProtoMessage()    {}
func (*Issue) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f976d1453a3f05a8, []int{3}
}
func (m *Issue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Issue.Unmarshal(m, b)
}
func (m *Issue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Issue.Marshal(b, m, deterministic)
}
func (dst *Issue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Issue.Merge(dst, src)
}
func (m *Issue) XXX_Size() int {
	return xxx_messageInfo_Issue.Size(m)
}
func (m *Issue) XXX_DiscardUnknown() {
	xxx_messageInfo_Issue.DiscardUnknown(m)
}

var xxx_messageInfo_Issue proto.InternalMessageInfo

func (m *Issue) GetOutputs() []*Token {
	if m != nil {
		return m.Outputs
	}
	return nil
}


type Transfer struct {
	
	Inputs []*TokenId `protobuf:"bytes,1,rep,name=inputs,proto3" json:"inputs,omitempty"`
	
	Outputs              []*Token `protobuf:"bytes,2,rep,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transfer) Reset()         { *m = Transfer{} }
func (m *Transfer) String() string { return proto.CompactTextString(m) }
func (*Transfer) ProtoMessage()    {}
func (*Transfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f976d1453a3f05a8, []int{4}
}
func (m *Transfer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transfer.Unmarshal(m, b)
}
func (m *Transfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transfer.Marshal(b, m, deterministic)
}
func (dst *Transfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transfer.Merge(dst, src)
}
func (m *Transfer) XXX_Size() int {
	return xxx_messageInfo_Transfer.Size(m)
}
func (m *Transfer) XXX_DiscardUnknown() {
	xxx_messageInfo_Transfer.DiscardUnknown(m)
}

var xxx_messageInfo_Transfer proto.InternalMessageInfo

func (m *Transfer) GetInputs() []*TokenId {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *Transfer) GetOutputs() []*Token {
	if m != nil {
		return m.Outputs
	}
	return nil
}


type Token struct {
	
	Owner *TokenOwner `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	
	
	Quantity             string   `protobuf:"bytes,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f976d1453a3f05a8, []int{5}
}
func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (dst *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(dst, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetOwner() *TokenOwner {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *Token) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Token) GetQuantity() string {
	if m != nil {
		return m.Quantity
	}
	return ""
}



type TokenId struct {
	
	TxId string `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	
	Index                uint32   `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenId) Reset()         { *m = TokenId{} }
func (m *TokenId) String() string { return proto.CompactTextString(m) }
func (*TokenId) ProtoMessage()    {}
func (*TokenId) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f976d1453a3f05a8, []int{6}
}
func (m *TokenId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenId.Unmarshal(m, b)
}
func (m *TokenId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenId.Marshal(b, m, deterministic)
}
func (dst *TokenId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenId.Merge(dst, src)
}
func (m *TokenId) XXX_Size() int {
	return xxx_messageInfo_TokenId.Size(m)
}
func (m *TokenId) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenId.DiscardUnknown(m)
}

var xxx_messageInfo_TokenId proto.InternalMessageInfo

func (m *TokenId) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *TokenId) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func init() {
	proto.RegisterType((*TokenTransaction)(nil), "token.TokenTransaction")
	proto.RegisterType((*TokenAction)(nil), "token.TokenAction")
	proto.RegisterType((*TokenOwner)(nil), "token.TokenOwner")
	proto.RegisterType((*Issue)(nil), "token.Issue")
	proto.RegisterType((*Transfer)(nil), "token.Transfer")
	proto.RegisterType((*Token)(nil), "token.Token")
	proto.RegisterType((*TokenId)(nil), "token.TokenId")
	proto.RegisterEnum("token.TokenOwner_Type", TokenOwner_Type_name, TokenOwner_Type_value)
}

func init() {
	proto.RegisterFile("token/transaction.proto", fileDescriptor_transaction_f976d1453a3f05a8)
}

var fileDescriptor_transaction_f976d1453a3f05a8 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x5d, 0xaf, 0xd2, 0x40,
	0x10, 0x6d, 0x2f, 0xb4, 0xb7, 0x77, 0x40, 0xc4, 0xd1, 0x68, 0x73, 0x9f, 0x6e, 0x1a, 0x83, 0x4a,
	0xb4, 0x4d, 0xd0, 0xc4, 0x67, 0x89, 0x18, 0xfa, 0xe0, 0x47, 0xd6, 0xfa, 0xc2, 0x0b, 0x16, 0xba,
	0xc0, 0x46, 0x69, 0xeb, 0x76, 0x1b, 0xe8, 0x0f, 0xf1, 0xff, 0x9a, 0x4e, 0x17, 0x68, 0x62, 0xf4,
	0x6d, 0x67, 0xce, 0x99, 0x33, 0xe7, 0x4c, 0x0b, 0x4f, 0x54, 0xf6, 0x83, 0xa7, 0x81, 0x92, 0x71,
	0x5a, 0xc4, 0x6b, 0x25, 0xb2, 0xd4, 0xcf, 0x65, 0xa6, 0x32, 0xb4, 0x08, 0xf0, 0xbe, 0xc1, 0x30,
	0xaa, 0x1f, 0xd1, 0x85, 0x80, 0x6f, 0xa1, 0x4f, 0xe0, 0xb2, 0xa9, 0x5d, 0xf3, 0xce, 0x7c, 0xde,
	0x9b, 0xa0, 0x4f, 0x4d, 0x9f, 0xe8, 0xef, 0x08, 0x99, 0x1b, 0xac, 0xa7, 0x2e, 0xe5, 0xd4, 0x01,
	0xbb, 0x19, 0xf1, 0x7e, 0x9b, 0xd0, 0x6b, 0x11, 0xf1, 0x29, 0x58, 0xa2, 0x28, 0x4a, 0xae, 0xb5,
	0xfa, 0x5a, 0x2b, 0xac, 0x7b, 0x73, 0x83, 0x35, 0x20, 0xbe, 0x02, 0x87, 0x8c, 0x6e, 0xb8, 0x74,
	0xaf, 0x88, 0x78, 0xff, 0xb4, 0x54, 0xb7, 0xe7, 0x06, 0x3b, 0x53, 0xf0, 0x05, 0xd8, 0x92, 0x27,
	0x9c, 0xef, 0xdd, 0xce, 0xbf, 0xc8, 0x9a, 0x30, 0xb5, 0xa1, 0x9b, 0xc4, 0x2a, 0xf6, 0x36, 0x00,
	0x64, 0xeb, 0xf3, 0x21, 0xe5, 0x12, 0xc7, 0xd0, 0x55, 0x55, 0xde, 0x98, 0x1a, 0x4c, 0x1e, 0xb7,
	0x03, 0x12, 0xc1, 0x8f, 0xaa, 0x9c, 0x33, 0xe2, 0xe0, 0x10, 0x3a, 0x32, 0x3e, 0x90, 0xad, 0x3e,
	0xab, 0x9f, 0xde, 0x2d, 0x74, 0x6b, 0x1c, 0x11, 0x06, 0x1f, 0xbf, 0x7e, 0x59, 0x86, 0xef, 0x67,
	0x9f, 0xa2, 0xf0, 0x43, 0x38, 0x63, 0x43, 0xc3, 0x0b, 0xc0, 0xa2, 0x6c, 0x38, 0x82, 0xeb, 0xac,
	0x54, 0x79, 0xa9, 0x0a, 0xd7, 0xbc, 0xeb, 0xb4, 0xa2, 0xd3, 0x16, 0x76, 0x02, 0xbd, 0x05, 0x38,
	0x27, 0xdb, 0x38, 0x02, 0x5b, 0xa4, 0xad, 0x91, 0x41, 0x7b, 0x24, 0x4c, 0x98, 0x46, 0xdb, 0xda,
	0x57, 0xff, 0xd3, 0xfe, 0x0e, 0x16, 0x75, 0xf0, 0x19, 0x58, 0x59, 0x9d, 0x4b, 0x7f, 0x85, 0x07,
	0x7f, 0x05, 0x66, 0x0d, 0x8e, 0xa8, 0x0f, 0x53, 0xa7, 0xbd, 0xd1, 0x07, 0xb8, 0x05, 0xe7, 0x57,
	0x19, 0xa7, 0x4a, 0xa8, 0x8a, 0xee, 0x7d, 0xc3, 0xce, 0xb5, 0xf7, 0x06, 0xae, 0xb5, 0x39, 0x7c,
	0x08, 0x96, 0x3a, 0x2e, 0x45, 0x42, 0x3b, 0xea, 0xd9, 0x63, 0x98, 0xe0, 0x23, 0xb0, 0x44, 0x9a,
	0xf0, 0x23, 0x09, 0xde, 0x63, 0x4d, 0x31, 0x7d, 0xb9, 0x18, 0x6f, 0x85, 0xda, 0x95, 0x2b, 0x7f,
	0x9d, 0xed, 0x83, 0x5d, 0x95, 0x73, 0xf9, 0x93, 0x27, 0x5b, 0x2e, 0x83, 0x4d, 0xbc, 0x92, 0x62,
	0x1d, 0xd0, 0x9f, 0x5a, 0x04, 0xe4, 0x72, 0x65, 0x53, 0xf5, 0xfa, 0x4f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x5a, 0xf7, 0x90, 0xb8, 0xd2, 0x02, 0x00, 0x00,
}
