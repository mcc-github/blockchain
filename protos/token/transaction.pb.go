


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
	return fileDescriptor_transaction_26436d18589330e1, []int{0}
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
	PlainAction *PlainTokenAction `protobuf:"bytes,1,opt,name=plain_action,json=plainAction,oneof"`
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
	return fileDescriptor_transaction_26436d18589330e1, []int{1}
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
	PlainImport *PlainImport `protobuf:"bytes,1,opt,name=plain_import,json=plainImport,oneof"`
}

type PlainTokenAction_PlainTransfer struct {
	PlainTransfer *PlainTransfer `protobuf:"bytes,2,opt,name=plain_transfer,json=plainTransfer,oneof"`
}
type PlainTokenAction_PlainRedeem struct {
	PlainRedeem *PlainTransfer `protobuf:"bytes,3,opt,name=plain_redeem,json=plainRedeem,oneof"`
}

func (*PlainTokenAction_PlainImport) isPlainTokenAction_Data() {}

func (*PlainTokenAction_PlainTransfer) isPlainTokenAction_Data() {}
func (*PlainTokenAction_PlainRedeem) isPlainTokenAction_Data()   {}

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

func (m *PlainTokenAction) GetPlainRedeem() *PlainTransfer {
	if x, ok := m.GetData().(*PlainTokenAction_PlainRedeem); ok {
		return x.PlainRedeem
	}
	return nil
}


func (*PlainTokenAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _PlainTokenAction_OneofMarshaler, _PlainTokenAction_OneofUnmarshaler, _PlainTokenAction_OneofSizer, []interface{}{
		(*PlainTokenAction_PlainImport)(nil),
		(*PlainTokenAction_PlainTransfer)(nil),
		(*PlainTokenAction_PlainRedeem)(nil),
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
	case *PlainTokenAction_PlainRedeem:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PlainRedeem); err != nil {
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
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PlainTransfer)
		err := b.DecodeMessage(msg)
		m.Data = &PlainTokenAction_PlainRedeem{msg}
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
	case *PlainTokenAction_PlainRedeem:
		s := proto.Size(x.PlainRedeem)
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
	
	Outputs              []*PlainOutput `protobuf:"bytes,1,rep,name=outputs" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *PlainImport) Reset()         { *m = PlainImport{} }
func (m *PlainImport) String() string { return proto.CompactTextString(m) }
func (*PlainImport) ProtoMessage()    {}
func (*PlainImport) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_26436d18589330e1, []int{2}
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
	
	Inputs []*InputId `protobuf:"bytes,1,rep,name=inputs" json:"inputs,omitempty"`
	
	Outputs              []*PlainOutput `protobuf:"bytes,2,rep,name=outputs" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *PlainTransfer) Reset()         { *m = PlainTransfer{} }
func (m *PlainTransfer) String() string { return proto.CompactTextString(m) }
func (*PlainTransfer) ProtoMessage()    {}
func (*PlainTransfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_26436d18589330e1, []int{3}
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
	
	Type string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	
	Quantity             uint64   `protobuf:"varint,3,opt,name=quantity" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlainOutput) Reset()         { *m = PlainOutput{} }
func (m *PlainOutput) String() string { return proto.CompactTextString(m) }
func (*PlainOutput) ProtoMessage()    {}
func (*PlainOutput) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_26436d18589330e1, []int{4}
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
	
	TxId string `protobuf:"bytes,1,opt,name=tx_id,json=txId" json:"tx_id,omitempty"`
	
	Index                uint32   `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InputId) Reset()         { *m = InputId{} }
func (m *InputId) String() string { return proto.CompactTextString(m) }
func (*InputId) ProtoMessage()    {}
func (*InputId) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_26436d18589330e1, []int{5}
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
	proto.RegisterFile("token/transaction.proto", fileDescriptor_transaction_26436d18589330e1)
}

var fileDescriptor_transaction_26436d18589330e1 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x4b, 0x4f, 0xf2, 0x40,
	0x14, 0xa5, 0x3c, 0x0a, 0xdf, 0xe5, 0x11, 0xbe, 0xd1, 0xc4, 0xc6, 0x15, 0xe9, 0xc2, 0x18, 0x63,
	0xda, 0x28, 0x3e, 0xd6, 0xb2, 0xa2, 0x2b, 0xcd, 0xc8, 0x46, 0x37, 0xa4, 0xd0, 0x01, 0x26, 0xc2,
	0xcc, 0x38, 0x4c, 0x23, 0xfc, 0x36, 0xff, 0x9c, 0xe9, 0x9d, 0x02, 0xd5, 0x44, 0x77, 0x73, 0xce,
	0x3d, 0xe7, 0xdc, 0x93, 0x9b, 0x81, 0x13, 0x23, 0xdf, 0x98, 0x08, 0x8d, 0x8e, 0xc5, 0x3a, 0x9e,
	0x1a, 0x2e, 0x45, 0xa0, 0xb4, 0x34, 0xd2, 0x1f, 0x41, 0x77, 0x94, 0x8d, 0x46, 0x87, 0x09, 0xb9,
	0x83, 0x96, 0x5a, 0xc6, 0x5c, 0x8c, 0x2d, 0xf6, 0x9c, 0x9e, 0x73, 0xde, 0xbc, 0xfe, 0x1f, 0x3c,
	0x65, 0x24, 0xaa, 0x1f, 0x70, 0x30, 0x2c, 0xd1, 0x26, 0x0a, 0x2d, 0x1c, 0x34, 0xc0, 0xb5, 0x0e,
	0xff, 0xd3, 0x81, 0xee, 0x4f, 0x35, 0xb9, 0xda, 0xc5, 0xf2, 0x95, 0x92, 0xda, 0xe4, 0xb1, 0x2d,
	0x1b, 0x1b, 0x21, 0xb7, 0x4f, 0xb4, 0x90, 0xdc, 0x43, 0xc7, 0x5a, 0xb0, 0xf8, 0x8c, 0x69, 0xaf,
	0x8c, 0xa6, 0x4e, 0xde, 0x25, 0x67, 0x87, 0x25, 0xda, 0x56, 0x45, 0x82, 0xf4, 0x77, 0xbb, 0x34,
	0x4b, 0x18, 0x5b, 0x79, 0x95, 0x5f, 0x6c, 0x76, 0x1b, 0x45, 0xd1, 0xc0, 0x85, 0x6a, 0x12, 0x9b,
	0xd8, 0xbf, 0x85, 0x66, 0xa1, 0x13, 0x39, 0x83, 0xba, 0x4c, 0x8d, 0x4a, 0xcd, 0xda, 0x73, 0x7a,
	0x95, 0x43, 0xe5, 0x47, 0x24, 0xe9, 0x6e, 0xe8, 0xbf, 0x40, 0xfb, 0x5b, 0x3c, 0xe9, 0x81, 0xcb,
	0x45, 0xc1, 0xd7, 0x08, 0xa2, 0x0c, 0x46, 0x09, 0xcd, 0xf9, 0x62, 0x74, 0xf9, 0xaf, 0xe8, 0xe7,
	0xbc, 0x91, 0xe5, 0xc9, 0x31, 0xd4, 0xe4, 0x87, 0x60, 0x1a, 0x4f, 0xd8, 0xa2, 0x16, 0x10, 0x02,
	0x55, 0xb3, 0x55, 0x0c, 0x4f, 0xf4, 0x8f, 0xe2, 0x9b, 0x9c, 0x42, 0xe3, 0x3d, 0x8d, 0x85, 0xe1,
	0x66, 0x8b, 0x37, 0xa8, 0xd2, 0x3d, 0xf6, 0x6f, 0xa0, 0x9e, 0xf7, 0x21, 0x47, 0x50, 0x33, 0x9b,
	0x31, 0x4f, 0x30, 0x30, 0xf3, 0x6e, 0xa2, 0x24, 0xdb, 0xc2, 0x45, 0xc2, 0x36, 0x18, 0xd8, 0xa6,
	0x16, 0x0c, 0x2e, 0x5f, 0x2f, 0xe6, 0xdc, 0x2c, 0xd2, 0x49, 0x30, 0x95, 0xab, 0x70, 0xb1, 0x55,
	0x4c, 0x2f, 0x59, 0x32, 0x67, 0x3a, 0x9c, 0xc5, 0x13, 0xcd, 0xa7, 0x21, 0xfe, 0xab, 0x75, 0x88,
	0x1f, 0x6e, 0xe2, 0x22, 0xea, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0x56, 0x69, 0x33, 0x85, 0x80,
	0x02, 0x00, 0x00,
}
