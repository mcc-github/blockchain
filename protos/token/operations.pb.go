


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


type TokenOperation struct {
	
	
	Operation            isTokenOperation_Operation `protobuf_oneof:"Operation"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *TokenOperation) Reset()         { *m = TokenOperation{} }
func (m *TokenOperation) String() string { return proto.CompactTextString(m) }
func (*TokenOperation) ProtoMessage()    {}
func (*TokenOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf84c62b8e84f69b, []int{0}
}

func (m *TokenOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenOperation.Unmarshal(m, b)
}
func (m *TokenOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenOperation.Marshal(b, m, deterministic)
}
func (m *TokenOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenOperation.Merge(m, src)
}
func (m *TokenOperation) XXX_Size() int {
	return xxx_messageInfo_TokenOperation.Size(m)
}
func (m *TokenOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenOperation.DiscardUnknown(m)
}

var xxx_messageInfo_TokenOperation proto.InternalMessageInfo

type isTokenOperation_Operation interface {
	isTokenOperation_Operation()
}

type TokenOperation_Action struct {
	Action *TokenOperationAction `protobuf:"bytes,1,opt,name=Action,proto3,oneof"`
}

func (*TokenOperation_Action) isTokenOperation_Operation() {}

func (m *TokenOperation) GetOperation() isTokenOperation_Operation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (m *TokenOperation) GetAction() *TokenOperationAction {
	if x, ok := m.GetOperation().(*TokenOperation_Action); ok {
		return x.Action
	}
	return nil
}


func (*TokenOperation) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TokenOperation_Action)(nil),
	}
}


type TokenOperationAction struct {
	
	
	
	Payload              isTokenOperationAction_Payload `protobuf_oneof:"Payload"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *TokenOperationAction) Reset()         { *m = TokenOperationAction{} }
func (m *TokenOperationAction) String() string { return proto.CompactTextString(m) }
func (*TokenOperationAction) ProtoMessage()    {}
func (*TokenOperationAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf84c62b8e84f69b, []int{1}
}

func (m *TokenOperationAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenOperationAction.Unmarshal(m, b)
}
func (m *TokenOperationAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenOperationAction.Marshal(b, m, deterministic)
}
func (m *TokenOperationAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenOperationAction.Merge(m, src)
}
func (m *TokenOperationAction) XXX_Size() int {
	return xxx_messageInfo_TokenOperationAction.Size(m)
}
func (m *TokenOperationAction) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenOperationAction.DiscardUnknown(m)
}

var xxx_messageInfo_TokenOperationAction proto.InternalMessageInfo

type isTokenOperationAction_Payload interface {
	isTokenOperationAction_Payload()
}

type TokenOperationAction_Issue struct {
	Issue *TokenActionTerms `protobuf:"bytes,1,opt,name=Issue,proto3,oneof"`
}

type TokenOperationAction_Transfer struct {
	Transfer *TokenActionTerms `protobuf:"bytes,2,opt,name=Transfer,proto3,oneof"`
}

func (*TokenOperationAction_Issue) isTokenOperationAction_Payload() {}

func (*TokenOperationAction_Transfer) isTokenOperationAction_Payload() {}

func (m *TokenOperationAction) GetPayload() isTokenOperationAction_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *TokenOperationAction) GetIssue() *TokenActionTerms {
	if x, ok := m.GetPayload().(*TokenOperationAction_Issue); ok {
		return x.Issue
	}
	return nil
}

func (m *TokenOperationAction) GetTransfer() *TokenActionTerms {
	if x, ok := m.GetPayload().(*TokenOperationAction_Transfer); ok {
		return x.Transfer
	}
	return nil
}


func (*TokenOperationAction) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TokenOperationAction_Issue)(nil),
		(*TokenOperationAction_Transfer)(nil),
	}
}


type TokenActionTerms struct {
	
	Sender *TokenOwner `protobuf:"bytes,1,opt,name=Sender,proto3" json:"Sender,omitempty"`
	
	Outputs              []*Token `protobuf:"bytes,2,rep,name=Outputs,proto3" json:"Outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenActionTerms) Reset()         { *m = TokenActionTerms{} }
func (m *TokenActionTerms) String() string { return proto.CompactTextString(m) }
func (*TokenActionTerms) ProtoMessage()    {}
func (*TokenActionTerms) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf84c62b8e84f69b, []int{2}
}

func (m *TokenActionTerms) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenActionTerms.Unmarshal(m, b)
}
func (m *TokenActionTerms) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenActionTerms.Marshal(b, m, deterministic)
}
func (m *TokenActionTerms) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenActionTerms.Merge(m, src)
}
func (m *TokenActionTerms) XXX_Size() int {
	return xxx_messageInfo_TokenActionTerms.Size(m)
}
func (m *TokenActionTerms) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenActionTerms.DiscardUnknown(m)
}

var xxx_messageInfo_TokenActionTerms proto.InternalMessageInfo

func (m *TokenActionTerms) GetSender() *TokenOwner {
	if m != nil {
		return m.Sender
	}
	return nil
}

func (m *TokenActionTerms) GetOutputs() []*Token {
	if m != nil {
		return m.Outputs
	}
	return nil
}

func init() {
	proto.RegisterType((*TokenOperation)(nil), "token.TokenOperation")
	proto.RegisterType((*TokenOperationAction)(nil), "token.TokenOperationAction")
	proto.RegisterType((*TokenActionTerms)(nil), "token.TokenActionTerms")
}

func init() { proto.RegisterFile("token/operations.proto", fileDescriptor_bf84c62b8e84f69b) }

var fileDescriptor_bf84c62b8e84f69b = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xd0, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0x07, 0xf0, 0x6d, 0xb2, 0xcd, 0xbd, 0x8a, 0x68, 0x10, 0x57, 0xf4, 0x32, 0x2a, 0xc8, 0xf4,
	0x90, 0xc0, 0x64, 0x1f, 0xc0, 0x9e, 0xe6, 0x69, 0x5a, 0x7b, 0xf2, 0x96, 0xb6, 0x6f, 0x5d, 0x71,
	0x6b, 0xca, 0x4b, 0x8a, 0xec, 0x23, 0xf8, 0xad, 0xa5, 0x4d, 0x37, 0x3a, 0x11, 0x3c, 0x26, 0xff,
	0xdf, 0xff, 0xf1, 0x12, 0xb8, 0x36, 0xea, 0x13, 0x73, 0xa1, 0x0a, 0x24, 0x69, 0x32, 0x95, 0x6b,
	0x5e, 0x90, 0x32, 0x8a, 0xf5, 0xeb, 0xfb, 0x9b, 0xb1, 0x8d, 0x0d, 0xc9, 0x5c, 0xcb, 0xb8, 0x02,
	0x36, 0xf7, 0x42, 0x38, 0x0f, 0xab, 0x68, 0xb9, 0x2f, 0xb2, 0x39, 0x0c, 0x9e, 0x6b, 0xe1, 0x76,
	0x27, 0xdd, 0xa9, 0x33, 0xbb, 0xe5, 0x75, 0x97, 0x1f, 0x33, 0x4b, 0x16, 0x9d, 0xa0, 0xc1, 0xbe,
	0x03, 0xa3, 0x43, 0xe8, 0x7d, 0x77, 0xe1, 0xea, 0x2f, 0xcf, 0x04, 0xf4, 0x5f, 0xb4, 0x2e, 0xb1,
	0x99, 0x3d, 0x6e, 0xcf, 0xb6, 0x24, 0x44, 0xda, 0xea, 0x45, 0x27, 0xb0, 0x8e, 0xcd, 0xe1, 0x34,
	0xac, 0x96, 0x5e, 0x21, 0xb9, 0xbd, 0xff, 0x3a, 0x07, 0xea, 0x8f, 0x60, 0xf8, 0x2a, 0x77, 0x1b,
	0x25, 0x13, 0x0f, 0xe1, 0xe2, 0x37, 0x65, 0x0f, 0x30, 0x78, 0xc7, 0x3c, 0x41, 0x6a, 0xf6, 0xb8,
	0x3c, 0x7a, 0xe3, 0x57, 0x8e, 0x14, 0x34, 0x80, 0xdd, 0xc3, 0x70, 0x59, 0x9a, 0xa2, 0x34, 0xda,
	0xed, 0x4d, 0x4e, 0xa6, 0xce, 0xec, 0xac, 0x6d, 0x83, 0x7d, 0xe8, 0xbf, 0xc1, 0x9d, 0xa2, 0x94,
	0xaf, 0x77, 0x05, 0xd2, 0x06, 0x93, 0x14, 0x89, 0xaf, 0x64, 0x44, 0x59, 0x6c, 0x3f, 0x5a, 0xdb,
	0xd6, 0xc7, 0x63, 0x9a, 0x99, 0x75, 0x19, 0xf1, 0x58, 0x6d, 0x45, 0xcb, 0x0a, 0x6b, 0x85, 0xb5,
	0xa2, 0xb6, 0xd1, 0xa0, 0x3e, 0x3d, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0xd1, 0xfc, 0xf1, 0x11,
	0xdc, 0x01, 0x00, 0x00,
}
