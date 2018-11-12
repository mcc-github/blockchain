


package token 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 




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
	return fileDescriptor_transaction_e0328d3afd4f7f06, []int{0}
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

type TokenTransaction_PlainAction struct {
	PlainAction *PlainTokenAction `protobuf:"bytes,1,opt,name=plain_action,json=plainAction,proto3,oneof"`
}

func (*TokenTransaction_PlainAction) isTokenTransaction_Action() {}

func (m *TokenTransaction) GetAction() isTokenTransaction_Action {
	if m != nil {
		return m.Action
	}
	return nil
}

func (m *TokenTransaction) GetPlainAction() *PlainTokenAction {
	if x, ok := m.GetAction().(*TokenTransaction_PlainAction); ok {
		return x.PlainAction
	}
	return nil
}


func (*TokenTransaction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TokenTransaction_OneofMarshaler, _TokenTransaction_OneofUnmarshaler, _TokenTransaction_OneofSizer, []interface{}{
		(*TokenTransaction_PlainAction)(nil),
	}
}

func _TokenTransaction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TokenTransaction)
	
	switch x := m.Action.(type) {
	case *TokenTransaction_PlainAction:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PlainAction); err != nil {
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
		msg := new(PlainTokenAction)
		err := b.DecodeMessage(msg)
		m.Action = &TokenTransaction_PlainAction{msg}
		return true, err
	default:
		return false, nil
	}
}

func _TokenTransaction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TokenTransaction)
	
	switch x := m.Action.(type) {
	case *TokenTransaction_PlainAction:
		s := proto.Size(x.PlainAction)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}



type PlainTokenAction struct {
	
	
	
	Data                 isPlainTokenAction_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *PlainTokenAction) Reset()         { *m = PlainTokenAction{} }
func (m *PlainTokenAction) String() string { return proto.CompactTextString(m) }
func (*PlainTokenAction) ProtoMessage()    {}
func (*PlainTokenAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_e0328d3afd4f7f06, []int{1}
}
func (m *PlainTokenAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainTokenAction.Unmarshal(m, b)
}
func (m *PlainTokenAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainTokenAction.Marshal(b, m, deterministic)
}
func (dst *PlainTokenAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainTokenAction.Merge(dst, src)
}
func (m *PlainTokenAction) XXX_Size() int {
	return xxx_messageInfo_PlainTokenAction.Size(m)
}
func (m *PlainTokenAction) XXX_DiscardUnknown() {
	xxx_messageInfo_PlainTokenAction.DiscardUnknown(m)
}

var xxx_messageInfo_PlainTokenAction proto.InternalMessageInfo

type isPlainTokenAction_Data interface {
	isPlainTokenAction_Data()
}

type PlainTokenAction_PlainImport struct {
	PlainImport *PlainImport `protobuf:"bytes,1,opt,name=plain_import,json=plainImport,proto3,oneof"`
}

type PlainTokenAction_PlainTransfer struct {
	PlainTransfer *PlainTransfer `protobuf:"bytes,2,opt,name=plain_transfer,json=plainTransfer,proto3,oneof"`
}

func (*PlainTokenAction_PlainImport) isPlainTokenAction_Data() {}

func (*PlainTokenAction_PlainTransfer) isPlainTokenAction_Data() {}

func (m *PlainTokenAction) GetData() isPlainTokenAction_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *PlainTokenAction) GetPlainImport() *PlainImport {
	if x, ok := m.GetData().(*PlainTokenAction_PlainImport); ok {
		return x.PlainImport
	}
	return nil
}

func (m *PlainTokenAction) GetPlainTransfer() *PlainTransfer {
	if x, ok := m.GetData().(*PlainTokenAction_PlainTransfer); ok {
		return x.PlainTransfer
	}
	return nil
}


func (*PlainTokenAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _PlainTokenAction_OneofMarshaler, _PlainTokenAction_OneofUnmarshaler, _PlainTokenAction_OneofSizer, []interface{}{
		(*PlainTokenAction_PlainImport)(nil),
		(*PlainTokenAction_PlainTransfer)(nil),
	}
}

func _PlainTokenAction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*PlainTokenAction)
	
	switch x := m.Data.(type) {
	case *PlainTokenAction_PlainImport:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PlainImport); err != nil {
			return err
		}
	case *PlainTokenAction_PlainTransfer:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PlainTransfer); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("PlainTokenAction.Data has unexpected type %T", x)
	}
	return nil
}

func _PlainTokenAction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*PlainTokenAction)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PlainImport)
		err := b.DecodeMessage(msg)
		m.Data = &PlainTokenAction_PlainImport{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PlainTransfer)
		err := b.DecodeMessage(msg)
		m.Data = &PlainTokenAction_PlainTransfer{msg}
		return true, err
	default:
		return false, nil
	}
}

func _PlainTokenAction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*PlainTokenAction)
	
	switch x := m.Data.(type) {
	case *PlainTokenAction_PlainImport:
		s := proto.Size(x.PlainImport)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *PlainTokenAction_PlainTransfer:
		s := proto.Size(x.PlainTransfer)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type PlainImport struct {
	
	Outputs              []*PlainOutput `protobuf:"bytes,1,rep,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *PlainImport) Reset()         { *m = PlainImport{} }
func (m *PlainImport) String() string { return proto.CompactTextString(m) }
func (*PlainImport) ProtoMessage()    {}
func (*PlainImport) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_e0328d3afd4f7f06, []int{2}
}
func (m *PlainImport) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainImport.Unmarshal(m, b)
}
func (m *PlainImport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainImport.Marshal(b, m, deterministic)
}
func (dst *PlainImport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainImport.Merge(dst, src)
}
func (m *PlainImport) XXX_Size() int {
	return xxx_messageInfo_PlainImport.Size(m)
}
func (m *PlainImport) XXX_DiscardUnknown() {
	xxx_messageInfo_PlainImport.DiscardUnknown(m)
}

var xxx_messageInfo_PlainImport proto.InternalMessageInfo

func (m *PlainImport) GetOutputs() []*PlainOutput {
	if m != nil {
		return m.Outputs
	}
	return nil
}


type PlainTransfer struct {
	
	Inputs []*InputId `protobuf:"bytes,1,rep,name=inputs,proto3" json:"inputs,omitempty"`
	
	Outputs              []*PlainOutput `protobuf:"bytes,2,rep,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *PlainTransfer) Reset()         { *m = PlainTransfer{} }
func (m *PlainTransfer) String() string { return proto.CompactTextString(m) }
func (*PlainTransfer) ProtoMessage()    {}
func (*PlainTransfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_e0328d3afd4f7f06, []int{3}
}
func (m *PlainTransfer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainTransfer.Unmarshal(m, b)
}
func (m *PlainTransfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainTransfer.Marshal(b, m, deterministic)
}
func (dst *PlainTransfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainTransfer.Merge(dst, src)
}
func (m *PlainTransfer) XXX_Size() int {
	return xxx_messageInfo_PlainTransfer.Size(m)
}
func (m *PlainTransfer) XXX_DiscardUnknown() {
	xxx_messageInfo_PlainTransfer.DiscardUnknown(m)
}

var xxx_messageInfo_PlainTransfer proto.InternalMessageInfo

func (m *PlainTransfer) GetInputs() []*InputId {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *PlainTransfer) GetOutputs() []*PlainOutput {
	if m != nil {
		return m.Outputs
	}
	return nil
}


type PlainOutput struct {
	
	Owner []byte `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	
	Quantity             uint64   `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlainOutput) Reset()         { *m = PlainOutput{} }
func (m *PlainOutput) String() string { return proto.CompactTextString(m) }
func (*PlainOutput) ProtoMessage()    {}
func (*PlainOutput) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_e0328d3afd4f7f06, []int{4}
}
func (m *PlainOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainOutput.Unmarshal(m, b)
}
func (m *PlainOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainOutput.Marshal(b, m, deterministic)
}
func (dst *PlainOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainOutput.Merge(dst, src)
}
func (m *PlainOutput) XXX_Size() int {
	return xxx_messageInfo_PlainOutput.Size(m)
}
func (m *PlainOutput) XXX_DiscardUnknown() {
	xxx_messageInfo_PlainOutput.DiscardUnknown(m)
}

var xxx_messageInfo_PlainOutput proto.InternalMessageInfo

func (m *PlainOutput) GetOwner() []byte {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *PlainOutput) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *PlainOutput) GetQuantity() uint64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}


type InputId struct {
	
	TxId string `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	
	Index                uint32   `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InputId) Reset()         { *m = InputId{} }
func (m *InputId) String() string { return proto.CompactTextString(m) }
func (*InputId) ProtoMessage()    {}
func (*InputId) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_e0328d3afd4f7f06, []int{5}
}
func (m *InputId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InputId.Unmarshal(m, b)
}
func (m *InputId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InputId.Marshal(b, m, deterministic)
}
func (dst *InputId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InputId.Merge(dst, src)
}
func (m *InputId) XXX_Size() int {
	return xxx_messageInfo_InputId.Size(m)
}
func (m *InputId) XXX_DiscardUnknown() {
	xxx_messageInfo_InputId.DiscardUnknown(m)
}

var xxx_messageInfo_InputId proto.InternalMessageInfo

func (m *InputId) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *InputId) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func init() {
	proto.RegisterType((*TokenTransaction)(nil), "TokenTransaction")
	proto.RegisterType((*PlainTokenAction)(nil), "PlainTokenAction")
	proto.RegisterType((*PlainImport)(nil), "PlainImport")
	proto.RegisterType((*PlainTransfer)(nil), "PlainTransfer")
	proto.RegisterType((*PlainOutput)(nil), "PlainOutput")
	proto.RegisterType((*InputId)(nil), "InputId")
}

func init() {
	proto.RegisterFile("token/transaction.proto", fileDescriptor_transaction_e0328d3afd4f7f06)
}

var fileDescriptor_transaction_e0328d3afd4f7f06 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x4b, 0x6f, 0xea, 0x30,
	0x10, 0x85, 0x09, 0x8f, 0xc0, 0x1d, 0x1e, 0xe2, 0xfa, 0x5e, 0xa9, 0x51, 0x57, 0x51, 0x16, 0x15,
	0xaa, 0xaa, 0x44, 0x7d, 0xaf, 0xcb, 0x8a, 0xac, 0x5a, 0xb9, 0x6c, 0xda, 0x0d, 0x0a, 0xc4, 0x80,
	0x55, 0xb0, 0x5d, 0xe3, 0xa8, 0xf0, 0x07, 0xfa, 0xbb, 0xab, 0x8c, 0x03, 0xa4, 0x5d, 0x74, 0x97,
	0x73, 0x66, 0xce, 0x37, 0x13, 0x6b, 0xe0, 0xc4, 0xc8, 0x37, 0x26, 0x22, 0xa3, 0x13, 0xb1, 0x49,
	0x66, 0x86, 0x4b, 0x11, 0x2a, 0x2d, 0x8d, 0x0c, 0xc6, 0xd0, 0x1f, 0xe7, 0xa5, 0xf1, 0xb1, 0x42,
	0xee, 0xa0, 0xa3, 0x56, 0x09, 0x17, 0x13, 0xab, 0x3d, 0xc7, 0x77, 0x06, 0xed, 0xab, 0xbf, 0xe1,
	0x53, 0x6e, 0x62, 0xf7, 0x03, 0x16, 0x46, 0x15, 0xda, 0xc6, 0x46, 0x2b, 0x87, 0x2d, 0x70, 0x6d,
	0x22, 0xf8, 0x74, 0xa0, 0xff, 0xb3, 0x9b, 0x5c, 0xee, 0xb1, 0x7c, 0xad, 0xa4, 0x36, 0x05, 0xb6,
	0x63, 0xb1, 0x31, 0x7a, 0x07, 0xa2, 0x95, 0xe4, 0x1e, 0x7a, 0x36, 0x82, 0x8b, 0xcf, 0x99, 0xf6,
	0xaa, 0x18, 0xea, 0x15, 0xbb, 0x14, 0xee, 0xa8, 0x42, 0xbb, 0xaa, 0x6c, 0x0c, 0x5d, 0xa8, 0xa7,
	0x89, 0x49, 0x82, 0x5b, 0x68, 0x97, 0xf0, 0xe4, 0x0c, 0x9a, 0x32, 0x33, 0x2a, 0x33, 0x1b, 0xcf,
	0xf1, 0x6b, 0xc7, 0xe9, 0x8f, 0x68, 0xd2, 0x7d, 0x31, 0x78, 0x81, 0xee, 0xb7, 0x01, 0xc4, 0x07,
	0x97, 0x8b, 0x52, 0xae, 0x15, 0xc6, 0xb9, 0x8c, 0x53, 0x5a, 0xf8, 0x65, 0x74, 0xf5, 0x37, 0xf4,
	0x73, 0xb1, 0x91, 0xf5, 0xc9, 0x7f, 0x68, 0xc8, 0x0f, 0xc1, 0x34, 0xbe, 0x46, 0x87, 0x5a, 0x41,
	0x08, 0xd4, 0xcd, 0x4e, 0x31, 0xfc, 0xdb, 0x3f, 0x14, 0xbf, 0xc9, 0x29, 0xb4, 0xde, 0xb3, 0x44,
	0x18, 0x6e, 0x76, 0x5e, 0xcd, 0x77, 0x06, 0x75, 0x7a, 0xd0, 0xc1, 0x0d, 0x34, 0x8b, 0x7d, 0xc8,
	0x3f, 0x68, 0x98, 0xed, 0x84, 0xa7, 0x08, 0xcc, 0xb3, 0xdb, 0x38, 0xcd, 0xa7, 0x70, 0x91, 0xb2,
	0x2d, 0x02, 0xbb, 0xd4, 0x8a, 0xe1, 0xc5, 0xeb, 0xf9, 0x82, 0x9b, 0x65, 0x36, 0x0d, 0x67, 0x72,
	0x1d, 0x2d, 0x77, 0x8a, 0xe9, 0x15, 0x4b, 0x17, 0x4c, 0x47, 0xf3, 0x64, 0xaa, 0xf9, 0x2c, 0xc2,
	0x13, 0xd9, 0x44, 0x78, 0x3b, 0x53, 0x17, 0xd5, 0xf5, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x17,
	0xe9, 0x14, 0xa0, 0x4b, 0x02, 0x00, 0x00,
}
