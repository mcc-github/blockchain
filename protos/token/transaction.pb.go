


package token

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)


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
	return fileDescriptor_fadc60fa5929c0a6, []int{1}
}

func (m *PlainTokenAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainTokenAction.Unmarshal(m, b)
}
func (m *PlainTokenAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainTokenAction.Marshal(b, m, deterministic)
}
func (m *PlainTokenAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainTokenAction.Merge(m, src)
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

type PlainTokenAction_PlainRedeem struct {
	PlainRedeem *PlainTransfer `protobuf:"bytes,3,opt,name=plain_redeem,json=plainRedeem,proto3,oneof"`
}

type PlainTokenAction_PlainApprove struct {
	PlainApprove *PlainApprove `protobuf:"bytes,4,opt,name=plain_approve,json=plainApprove,proto3,oneof"`
}

func (*PlainTokenAction_PlainImport) isPlainTokenAction_Data() {}

func (*PlainTokenAction_PlainTransfer) isPlainTokenAction_Data() {}

func (*PlainTokenAction_PlainRedeem) isPlainTokenAction_Data() {}

func (*PlainTokenAction_PlainApprove) isPlainTokenAction_Data() {}

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

func (m *PlainTokenAction) GetPlainApprove() *PlainApprove {
	if x, ok := m.GetData().(*PlainTokenAction_PlainApprove); ok {
		return x.PlainApprove
	}
	return nil
}


func (*PlainTokenAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _PlainTokenAction_OneofMarshaler, _PlainTokenAction_OneofUnmarshaler, _PlainTokenAction_OneofSizer, []interface{}{
		(*PlainTokenAction_PlainImport)(nil),
		(*PlainTokenAction_PlainTransfer)(nil),
		(*PlainTokenAction_PlainRedeem)(nil),
		(*PlainTokenAction_PlainApprove)(nil),
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
	case *PlainTokenAction_PlainApprove:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PlainApprove); err != nil {
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
	case 4: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PlainApprove)
		err := b.DecodeMessage(msg)
		m.Data = &PlainTokenAction_PlainApprove{msg}
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
	case *PlainTokenAction_PlainApprove:
		s := proto.Size(x.PlainApprove)
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
	return fileDescriptor_fadc60fa5929c0a6, []int{2}
}

func (m *PlainImport) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainImport.Unmarshal(m, b)
}
func (m *PlainImport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainImport.Marshal(b, m, deterministic)
}
func (m *PlainImport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainImport.Merge(m, src)
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
	return fileDescriptor_fadc60fa5929c0a6, []int{3}
}

func (m *PlainTransfer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainTransfer.Unmarshal(m, b)
}
func (m *PlainTransfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainTransfer.Marshal(b, m, deterministic)
}
func (m *PlainTransfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainTransfer.Merge(m, src)
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


type PlainApprove struct {
	
	Inputs []*InputId `protobuf:"bytes,1,rep,name=inputs,proto3" json:"inputs,omitempty"`
	
	DelegatedOutputs []*PlainDelegatedOutput `protobuf:"bytes,2,rep,name=delegated_outputs,json=delegatedOutputs,proto3" json:"delegated_outputs,omitempty"`
	
	Output               *PlainOutput `protobuf:"bytes,3,opt,name=output,proto3" json:"output,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *PlainApprove) Reset()         { *m = PlainApprove{} }
func (m *PlainApprove) String() string { return proto.CompactTextString(m) }
func (*PlainApprove) ProtoMessage()    {}
func (*PlainApprove) Descriptor() ([]byte, []int) {
	return fileDescriptor_fadc60fa5929c0a6, []int{4}
}

func (m *PlainApprove) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainApprove.Unmarshal(m, b)
}
func (m *PlainApprove) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainApprove.Marshal(b, m, deterministic)
}
func (m *PlainApprove) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainApprove.Merge(m, src)
}
func (m *PlainApprove) XXX_Size() int {
	return xxx_messageInfo_PlainApprove.Size(m)
}
func (m *PlainApprove) XXX_DiscardUnknown() {
	xxx_messageInfo_PlainApprove.DiscardUnknown(m)
}

var xxx_messageInfo_PlainApprove proto.InternalMessageInfo

func (m *PlainApprove) GetInputs() []*InputId {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *PlainApprove) GetDelegatedOutputs() []*PlainDelegatedOutput {
	if m != nil {
		return m.DelegatedOutputs
	}
	return nil
}

func (m *PlainApprove) GetOutput() *PlainOutput {
	if m != nil {
		return m.Output
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
	return fileDescriptor_fadc60fa5929c0a6, []int{5}
}

func (m *PlainOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainOutput.Unmarshal(m, b)
}
func (m *PlainOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainOutput.Marshal(b, m, deterministic)
}
func (m *PlainOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainOutput.Merge(m, src)
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
	return fileDescriptor_fadc60fa5929c0a6, []int{6}
}

func (m *InputId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InputId.Unmarshal(m, b)
}
func (m *InputId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InputId.Marshal(b, m, deterministic)
}
func (m *InputId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InputId.Merge(m, src)
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


type PlainDelegatedOutput struct {
	
	Owner []byte `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	
	
	Delegatees [][]byte `protobuf:"bytes,2,rep,name=delegatees,proto3" json:"delegatees,omitempty"`
	
	Type string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	
	Quantity             uint64   `protobuf:"varint,4,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlainDelegatedOutput) Reset()         { *m = PlainDelegatedOutput{} }
func (m *PlainDelegatedOutput) String() string { return proto.CompactTextString(m) }
func (*PlainDelegatedOutput) ProtoMessage()    {}
func (*PlainDelegatedOutput) Descriptor() ([]byte, []int) {
	return fileDescriptor_fadc60fa5929c0a6, []int{7}
}

func (m *PlainDelegatedOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlainDelegatedOutput.Unmarshal(m, b)
}
func (m *PlainDelegatedOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlainDelegatedOutput.Marshal(b, m, deterministic)
}
func (m *PlainDelegatedOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlainDelegatedOutput.Merge(m, src)
}
func (m *PlainDelegatedOutput) XXX_Size() int {
	return xxx_messageInfo_PlainDelegatedOutput.Size(m)
}
func (m *PlainDelegatedOutput) XXX_DiscardUnknown() {
	xxx_messageInfo_PlainDelegatedOutput.DiscardUnknown(m)
}

var xxx_messageInfo_PlainDelegatedOutput proto.InternalMessageInfo

func (m *PlainDelegatedOutput) GetOwner() []byte {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *PlainDelegatedOutput) GetDelegatees() [][]byte {
	if m != nil {
		return m.Delegatees
	}
	return nil
}

func (m *PlainDelegatedOutput) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *PlainDelegatedOutput) GetQuantity() uint64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func init() {
	proto.RegisterType((*TokenTransaction)(nil), "TokenTransaction")
	proto.RegisterType((*PlainTokenAction)(nil), "PlainTokenAction")
	proto.RegisterType((*PlainImport)(nil), "PlainImport")
	proto.RegisterType((*PlainTransfer)(nil), "PlainTransfer")
	proto.RegisterType((*PlainApprove)(nil), "PlainApprove")
	proto.RegisterType((*PlainOutput)(nil), "PlainOutput")
	proto.RegisterType((*InputId)(nil), "InputId")
	proto.RegisterType((*PlainDelegatedOutput)(nil), "PlainDelegatedOutput")
}

func init() { proto.RegisterFile("token/transaction.proto", fileDescriptor_fadc60fa5929c0a6) }

var fileDescriptor_fadc60fa5929c0a6 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x4f, 0x8f, 0x12, 0x31,
	0x14, 0x67, 0x96, 0xd9, 0x59, 0x7c, 0xcc, 0x6c, 0xd8, 0xba, 0xc6, 0x89, 0x07, 0x43, 0x26, 0xc6,
	0x18, 0x63, 0x98, 0xe8, 0xae, 0x7a, 0x5e, 0xe2, 0x01, 0x4e, 0x9a, 0xca, 0x45, 0x2f, 0x64, 0xa0,
	0x5d, 0xb6, 0x11, 0xda, 0x5a, 0x8a, 0x42, 0xe2, 0x27, 0xf1, 0x5b, 0xfa, 0x0d, 0xcc, 0xbc, 0x76,
	0x60, 0xc6, 0x88, 0xf1, 0xc6, 0xef, 0xf7, 0xfa, 0xfb, 0xf3, 0xca, 0x14, 0x1e, 0x5a, 0xf5, 0x85,
	0xcb, 0xdc, 0x9a, 0x42, 0xae, 0x8b, 0xb9, 0x15, 0x4a, 0x0e, 0xb4, 0x51, 0x56, 0x65, 0x13, 0xe8,
	0x4d, 0xca, 0xd1, 0xe4, 0x30, 0x21, 0x6f, 0x20, 0xd6, 0xcb, 0x42, 0xc8, 0xa9, 0xc3, 0x69, 0xd0,
	0x0f, 0x9e, 0x75, 0x5f, 0x5d, 0x0c, 0x3e, 0x94, 0x24, 0x9e, 0xbe, 0xc1, 0xc1, 0xa8, 0x45, 0xbb,
	0x78, 0xd0, 0xc1, 0x61, 0x07, 0x22, 0xa7, 0xc8, 0x7e, 0x05, 0xd0, 0xfb, 0xf3, 0x34, 0x79, 0x59,
	0xd9, 0x8a, 0x95, 0x56, 0xc6, 0x7a, 0xdb, 0xd8, 0xd9, 0x8e, 0x91, 0xdb, 0x3b, 0x3a, 0x48, 0xde,
	0xc2, 0xb9, 0x93, 0x60, 0xf1, 0x5b, 0x6e, 0xd2, 0x13, 0x14, 0x9d, 0xfb, 0x2e, 0x9e, 0x1d, 0xb5,
	0x68, 0xa2, 0xeb, 0x04, 0xb9, 0xaa, 0xb2, 0x0c, 0x67, 0x9c, 0xaf, 0xd2, 0xf6, 0x11, 0x99, 0x4b,
	0xa3, 0x78, 0x88, 0x5c, 0x43, 0xe2, 0xf7, 0xd6, 0xda, 0xa8, 0x6f, 0x3c, 0x0d, 0x51, 0x95, 0x38,
	0xd5, 0x8d, 0x23, 0x47, 0x2d, 0xea, 0xac, 0x3d, 0x1e, 0x46, 0x10, 0xb2, 0xc2, 0x16, 0xd9, 0x6b,
	0xe8, 0xd6, 0x36, 0x21, 0x4f, 0xe1, 0x4c, 0x6d, 0xac, 0xde, 0xd8, 0x75, 0x1a, 0xf4, 0xdb, 0x87,
	0x45, 0xdf, 0x23, 0x49, 0xab, 0x61, 0xf6, 0x09, 0x92, 0x46, 0x29, 0xd2, 0x87, 0x48, 0xc8, 0x9a,
	0xae, 0x33, 0x18, 0x97, 0x70, 0xcc, 0xa8, 0xe7, 0xeb, 0xd6, 0x27, 0xff, 0xb2, 0xfe, 0x19, 0x40,
	0x5c, 0xaf, 0xfe, 0x1f, 0xd6, 0x43, 0xb8, 0x60, 0x7c, 0xc9, 0x17, 0x85, 0xe5, 0x6c, 0xda, 0x0c,
	0x79, 0xe0, 0x42, 0xde, 0x55, 0x63, 0x9f, 0xd6, 0x63, 0x4d, 0x62, 0x4d, 0x9e, 0x40, 0xe4, 0x94,
	0xfe, 0xd6, 0x9b, 0xed, 0xfc, 0x2c, 0xfb, 0xe8, 0xaf, 0xcb, 0xd1, 0xe4, 0x12, 0x4e, 0xd5, 0x77,
	0xc9, 0x0d, 0x7e, 0x15, 0x31, 0x75, 0x80, 0x10, 0x08, 0xed, 0x4e, 0x73, 0xfc, 0xd7, 0xef, 0x51,
	0xfc, 0x4d, 0x1e, 0x41, 0xe7, 0xeb, 0xa6, 0x90, 0x56, 0xd8, 0x1d, 0x06, 0x84, 0x74, 0x8f, 0xb3,
	0x6b, 0x38, 0xf3, 0x1b, 0x91, 0xfb, 0x70, 0x6a, 0xb7, 0x53, 0xc1, 0xd0, 0xb0, 0xd4, 0x6e, 0xc7,
	0xac, 0x4c, 0x11, 0x92, 0xf1, 0x2d, 0x1a, 0x26, 0xd4, 0x81, 0xec, 0x07, 0x5c, 0xfe, 0x6d, 0xb5,
	0x23, 0x9d, 0x1e, 0x03, 0x54, 0x2b, 0x73, 0x77, 0x37, 0x31, 0xad, 0x31, 0xfb, 0xce, 0xed, 0x23,
	0x9d, 0xc3, 0x66, 0xe7, 0xe1, 0x8b, 0xcf, 0xcf, 0x17, 0xc2, 0xde, 0x6d, 0x66, 0x83, 0xb9, 0x5a,
	0xe5, 0x77, 0x3b, 0xcd, 0xcd, 0x92, 0xb3, 0x05, 0x37, 0xf9, 0x6d, 0x31, 0x33, 0x62, 0x9e, 0xe3,
	0x43, 0x5d, 0xe7, 0xf8, 0x82, 0x67, 0x11, 0xa2, 0xab, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x4d,
	0x22, 0xa4, 0x13, 0xd1, 0x03, 0x00, 0x00,
}
