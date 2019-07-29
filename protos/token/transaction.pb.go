


package token

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 

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
	return fileDescriptor_fadc60fa5929c0a6, []int{2, 0}
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
	return fileDescriptor_fadc60fa5929c0a6, []int{0}
}

func (m *TokenTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenTransaction.Unmarshal(m, b)
}
func (m *TokenTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenTransaction.Marshal(b, m, deterministic)
}
func (m *TokenTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenTransaction.Merge(m, src)
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


func (*TokenTransaction) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TokenTransaction_TokenAction)(nil),
	}
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
	return fileDescriptor_fadc60fa5929c0a6, []int{1}
}

func (m *TokenAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenAction.Unmarshal(m, b)
}
func (m *TokenAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenAction.Marshal(b, m, deterministic)
}
func (m *TokenAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenAction.Merge(m, src)
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


func (*TokenAction) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TokenAction_Issue)(nil),
		(*TokenAction_Transfer)(nil),
		(*TokenAction_Redeem)(nil),
	}
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
	return fileDescriptor_fadc60fa5929c0a6, []int{2}
}

func (m *TokenOwner) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenOwner.Unmarshal(m, b)
}
func (m *TokenOwner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenOwner.Marshal(b, m, deterministic)
}
func (m *TokenOwner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenOwner.Merge(m, src)
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
	return fileDescriptor_fadc60fa5929c0a6, []int{3}
}

func (m *Issue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Issue.Unmarshal(m, b)
}
func (m *Issue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Issue.Marshal(b, m, deterministic)
}
func (m *Issue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Issue.Merge(m, src)
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
	return fileDescriptor_fadc60fa5929c0a6, []int{4}
}

func (m *Transfer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transfer.Unmarshal(m, b)
}
func (m *Transfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transfer.Marshal(b, m, deterministic)
}
func (m *Transfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transfer.Merge(m, src)
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
	return fileDescriptor_fadc60fa5929c0a6, []int{5}
}

func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (m *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(m, src)
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
	return fileDescriptor_fadc60fa5929c0a6, []int{6}
}

func (m *TokenId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenId.Unmarshal(m, b)
}
func (m *TokenId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenId.Marshal(b, m, deterministic)
}
func (m *TokenId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenId.Merge(m, src)
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
	proto.RegisterEnum("token.TokenOwner_Type", TokenOwner_Type_name, TokenOwner_Type_value)
	proto.RegisterType((*TokenTransaction)(nil), "token.TokenTransaction")
	proto.RegisterType((*TokenAction)(nil), "token.TokenAction")
	proto.RegisterType((*TokenOwner)(nil), "token.TokenOwner")
	proto.RegisterType((*Issue)(nil), "token.Issue")
	proto.RegisterType((*Transfer)(nil), "token.Transfer")
	proto.RegisterType((*Token)(nil), "token.Token")
	proto.RegisterType((*TokenId)(nil), "token.TokenId")
}

func init() { proto.RegisterFile("token/transaction.proto", fileDescriptor_fadc60fa5929c0a6) }

var fileDescriptor_fadc60fa5929c0a6 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x4f, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0xed, 0x26, 0x76, 0xdd, 0x49, 0x08, 0x61, 0x40, 0x60, 0xf5, 0x54, 0x19, 0x54, 0xa0,
	0x12, 0xb6, 0x54, 0x90, 0x38, 0x13, 0x51, 0x14, 0x1f, 0xf8, 0xb7, 0x98, 0x4b, 0x2f, 0xc1, 0x89,
	0x37, 0xe9, 0x0a, 0xea, 0x35, 0xeb, 0xb5, 0x1a, 0x7f, 0x10, 0xbe, 0x2f, 0xf2, 0xec, 0xa6, 0xb1,
	0x84, 0xe8, 0x6d, 0x67, 0xde, 0x6f, 0x66, 0xdf, 0x5b, 0x1b, 0x9e, 0x68, 0xf9, 0x93, 0x97, 0x89,
	0x56, 0x79, 0x59, 0xe7, 0x2b, 0x2d, 0x64, 0x19, 0x57, 0x4a, 0x6a, 0x89, 0x1e, 0x09, 0xd1, 0x77,
	0x98, 0x66, 0xdd, 0x21, 0xdb, 0x03, 0xf8, 0x16, 0xc6, 0x24, 0x2e, 0x4c, 0x1d, 0xba, 0x27, 0xee,
	0x8b, 0xd1, 0x39, 0xc6, 0xd4, 0x8c, 0x09, 0x7f, 0x47, 0xca, 0xdc, 0x61, 0x23, 0xbd, 0x2f, 0x67,
	0x01, 0xf8, 0x66, 0x24, 0xfa, 0xe3, 0xc2, 0xa8, 0x07, 0xe2, 0x33, 0xf0, 0x44, 0x5d, 0x37, 0xdc,
	0xee, 0x1a, 0xdb, 0x5d, 0x69, 0xd7, 0x9b, 0x3b, 0xcc, 0x88, 0xf8, 0x0a, 0x02, 0x32, 0xba, 0xe6,
	0x2a, 0x3c, 0x20, 0xf0, 0xfe, 0xee, 0x52, 0xdb, 0x9e, 0x3b, 0xec, 0x16, 0xc1, 0x97, 0xe0, 0x2b,
	0x5e, 0x70, 0x7e, 0x1d, 0x0e, 0xfe, 0x07, 0x5b, 0x60, 0xe6, 0xc3, 0xb0, 0xc8, 0x75, 0x1e, 0xad,
	0x01, 0xc8, 0xd6, 0xe7, 0x9b, 0x92, 0x2b, 0x3c, 0x83, 0xa1, 0x6e, 0x2b, 0x63, 0x6a, 0x72, 0xfe,
	0xb8, 0x1f, 0x90, 0x80, 0x38, 0x6b, 0x2b, 0xce, 0x88, 0xc1, 0x29, 0x0c, 0x54, 0x7e, 0x43, 0xb6,
	0xc6, 0xac, 0x3b, 0x46, 0xc7, 0x30, 0xec, 0x74, 0x44, 0x98, 0x7c, 0xfc, 0xf6, 0x65, 0x91, 0xbe,
	0xbf, 0xf8, 0x94, 0xa5, 0x1f, 0xd2, 0x0b, 0x36, 0x75, 0xa2, 0x04, 0x3c, 0xca, 0x86, 0xa7, 0x70,
	0x28, 0x1b, 0x5d, 0x35, 0xba, 0x0e, 0xdd, 0x93, 0x41, 0x2f, 0x3a, 0xdd, 0xc2, 0x76, 0x62, 0x74,
	0x09, 0xc1, 0xce, 0x36, 0x9e, 0x82, 0x2f, 0xca, 0xde, 0xc8, 0xa4, 0x3f, 0x92, 0x16, 0xcc, 0xaa,
	0xfd, 0xdd, 0x07, 0x77, 0xed, 0xfe, 0x01, 0x1e, 0x75, 0xf0, 0x39, 0x78, 0xb2, 0xcb, 0x65, 0xbf,
	0xc2, 0x83, 0x7f, 0x02, 0x33, 0xa3, 0x23, 0xda, 0x87, 0xe9, 0xd2, 0x1e, 0xd9, 0x07, 0x38, 0x86,
	0xe0, 0x77, 0x93, 0x97, 0x5a, 0xe8, 0x96, 0xde, 0xfb, 0x88, 0xdd, 0xd6, 0xd1, 0x1b, 0x38, 0xb4,
	0xe6, 0xf0, 0x21, 0x78, 0x7a, 0xbb, 0x10, 0x05, 0xdd, 0xd1, 0xcd, 0x6e, 0xd3, 0x02, 0x1f, 0x81,
	0x27, 0xca, 0x82, 0x6f, 0x69, 0xe1, 0x3d, 0x66, 0x8a, 0xd9, 0x57, 0x78, 0x2a, 0xd5, 0x26, 0xbe,
	0x6a, 0x2b, 0xae, 0x7e, 0xf1, 0x62, 0xc3, 0x55, 0xbc, 0xce, 0x97, 0x4a, 0xac, 0xcc, 0x2f, 0x5a,
	0x1b, 0x7b, 0x97, 0x67, 0x1b, 0xa1, 0xaf, 0x9a, 0x65, 0xbc, 0x92, 0xd7, 0x49, 0x8f, 0x4d, 0x0c,
	0x9b, 0x18, 0x36, 0x21, 0x76, 0xe9, 0x53, 0xf5, 0xfa, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x42,
	0x2d, 0xf5, 0x41, 0xf7, 0x02, 0x00, 0x00,
}
