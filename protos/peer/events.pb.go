


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import common "github.com/mcc-github/blockchain/protos/common"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type EventType int32

const (
	EventType_REGISTER      EventType = 0
	EventType_BLOCK         EventType = 1
	EventType_CHAINCODE     EventType = 2
	EventType_REJECTION     EventType = 3
	EventType_FILTEREDBLOCK EventType = 4
)

var EventType_name = map[int32]string{
	0: "REGISTER",
	1: "BLOCK",
	2: "CHAINCODE",
	3: "REJECTION",
	4: "FILTEREDBLOCK",
}
var EventType_value = map[string]int32{
	"REGISTER":      0,
	"BLOCK":         1,
	"CHAINCODE":     2,
	"REJECTION":     3,
	"FILTEREDBLOCK": 4,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}
func (EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{0}
}



type ChaincodeReg struct {
	ChaincodeId          string   `protobuf:"bytes,1,opt,name=chaincode_id,json=chaincodeId" json:"chaincode_id,omitempty"`
	EventName            string   `protobuf:"bytes,2,opt,name=event_name,json=eventName" json:"event_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeReg) Reset()         { *m = ChaincodeReg{} }
func (m *ChaincodeReg) String() string { return proto.CompactTextString(m) }
func (*ChaincodeReg) ProtoMessage()    {}
func (*ChaincodeReg) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{0}
}
func (m *ChaincodeReg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeReg.Unmarshal(m, b)
}
func (m *ChaincodeReg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeReg.Marshal(b, m, deterministic)
}
func (dst *ChaincodeReg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeReg.Merge(dst, src)
}
func (m *ChaincodeReg) XXX_Size() int {
	return xxx_messageInfo_ChaincodeReg.Size(m)
}
func (m *ChaincodeReg) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeReg.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeReg proto.InternalMessageInfo

func (m *ChaincodeReg) GetChaincodeId() string {
	if m != nil {
		return m.ChaincodeId
	}
	return ""
}

func (m *ChaincodeReg) GetEventName() string {
	if m != nil {
		return m.EventName
	}
	return ""
}

type Interest struct {
	EventType EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,enum=protos.EventType" json:"event_type,omitempty"`
	
	
	
	
	
	
	
	RegInfo              isInterest_RegInfo `protobuf_oneof:"RegInfo"`
	ChainID              string             `protobuf:"bytes,3,opt,name=chainID" json:"chainID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Interest) Reset()         { *m = Interest{} }
func (m *Interest) String() string { return proto.CompactTextString(m) }
func (*Interest) ProtoMessage()    {}
func (*Interest) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{1}
}
func (m *Interest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Interest.Unmarshal(m, b)
}
func (m *Interest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Interest.Marshal(b, m, deterministic)
}
func (dst *Interest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Interest.Merge(dst, src)
}
func (m *Interest) XXX_Size() int {
	return xxx_messageInfo_Interest.Size(m)
}
func (m *Interest) XXX_DiscardUnknown() {
	xxx_messageInfo_Interest.DiscardUnknown(m)
}

var xxx_messageInfo_Interest proto.InternalMessageInfo

type isInterest_RegInfo interface {
	isInterest_RegInfo()
}

type Interest_ChaincodeRegInfo struct {
	ChaincodeRegInfo *ChaincodeReg `protobuf:"bytes,2,opt,name=chaincode_reg_info,json=chaincodeRegInfo,oneof"`
}

func (*Interest_ChaincodeRegInfo) isInterest_RegInfo() {}

func (m *Interest) GetRegInfo() isInterest_RegInfo {
	if m != nil {
		return m.RegInfo
	}
	return nil
}

func (m *Interest) GetEventType() EventType {
	if m != nil {
		return m.EventType
	}
	return EventType_REGISTER
}

func (m *Interest) GetChaincodeRegInfo() *ChaincodeReg {
	if x, ok := m.GetRegInfo().(*Interest_ChaincodeRegInfo); ok {
		return x.ChaincodeRegInfo
	}
	return nil
}

func (m *Interest) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}


func (*Interest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Interest_OneofMarshaler, _Interest_OneofUnmarshaler, _Interest_OneofSizer, []interface{}{
		(*Interest_ChaincodeRegInfo)(nil),
	}
}

func _Interest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Interest)
	
	switch x := m.RegInfo.(type) {
	case *Interest_ChaincodeRegInfo:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChaincodeRegInfo); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Interest.RegInfo has unexpected type %T", x)
	}
	return nil
}

func _Interest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Interest)
	switch tag {
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeReg)
		err := b.DecodeMessage(msg)
		m.RegInfo = &Interest_ChaincodeRegInfo{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Interest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Interest)
	
	switch x := m.RegInfo.(type) {
	case *Interest_ChaincodeRegInfo:
		s := proto.Size(x.ChaincodeRegInfo)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}




type Register struct {
	Events               []*Interest `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Register) Reset()         { *m = Register{} }
func (m *Register) String() string { return proto.CompactTextString(m) }
func (*Register) ProtoMessage()    {}
func (*Register) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{2}
}
func (m *Register) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Register.Unmarshal(m, b)
}
func (m *Register) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Register.Marshal(b, m, deterministic)
}
func (dst *Register) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Register.Merge(dst, src)
}
func (m *Register) XXX_Size() int {
	return xxx_messageInfo_Register.Size(m)
}
func (m *Register) XXX_DiscardUnknown() {
	xxx_messageInfo_Register.DiscardUnknown(m)
}

var xxx_messageInfo_Register proto.InternalMessageInfo

func (m *Register) GetEvents() []*Interest {
	if m != nil {
		return m.Events
	}
	return nil
}



type Rejection struct {
	Tx                   *Transaction `protobuf:"bytes,1,opt,name=tx" json:"tx,omitempty"`
	ErrorMsg             string       `protobuf:"bytes,2,opt,name=error_msg,json=errorMsg" json:"error_msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Rejection) Reset()         { *m = Rejection{} }
func (m *Rejection) String() string { return proto.CompactTextString(m) }
func (*Rejection) ProtoMessage()    {}
func (*Rejection) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{3}
}
func (m *Rejection) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Rejection.Unmarshal(m, b)
}
func (m *Rejection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Rejection.Marshal(b, m, deterministic)
}
func (dst *Rejection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Rejection.Merge(dst, src)
}
func (m *Rejection) XXX_Size() int {
	return xxx_messageInfo_Rejection.Size(m)
}
func (m *Rejection) XXX_DiscardUnknown() {
	xxx_messageInfo_Rejection.DiscardUnknown(m)
}

var xxx_messageInfo_Rejection proto.InternalMessageInfo

func (m *Rejection) GetTx() *Transaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *Rejection) GetErrorMsg() string {
	if m != nil {
		return m.ErrorMsg
	}
	return ""
}


type Unregister struct {
	Events               []*Interest `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Unregister) Reset()         { *m = Unregister{} }
func (m *Unregister) String() string { return proto.CompactTextString(m) }
func (*Unregister) ProtoMessage()    {}
func (*Unregister) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{4}
}
func (m *Unregister) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Unregister.Unmarshal(m, b)
}
func (m *Unregister) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Unregister.Marshal(b, m, deterministic)
}
func (dst *Unregister) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Unregister.Merge(dst, src)
}
func (m *Unregister) XXX_Size() int {
	return xxx_messageInfo_Unregister.Size(m)
}
func (m *Unregister) XXX_DiscardUnknown() {
	xxx_messageInfo_Unregister.DiscardUnknown(m)
}

var xxx_messageInfo_Unregister proto.InternalMessageInfo

func (m *Unregister) GetEvents() []*Interest {
	if m != nil {
		return m.Events
	}
	return nil
}



type FilteredBlock struct {
	ChannelId            string                 `protobuf:"bytes,1,opt,name=channel_id,json=channelId" json:"channel_id,omitempty"`
	Number               uint64                 `protobuf:"varint,2,opt,name=number" json:"number,omitempty"`
	FilteredTransactions []*FilteredTransaction `protobuf:"bytes,4,rep,name=filtered_transactions,json=filteredTransactions" json:"filtered_transactions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *FilteredBlock) Reset()         { *m = FilteredBlock{} }
func (m *FilteredBlock) String() string { return proto.CompactTextString(m) }
func (*FilteredBlock) ProtoMessage()    {}
func (*FilteredBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{5}
}
func (m *FilteredBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredBlock.Unmarshal(m, b)
}
func (m *FilteredBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredBlock.Marshal(b, m, deterministic)
}
func (dst *FilteredBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredBlock.Merge(dst, src)
}
func (m *FilteredBlock) XXX_Size() int {
	return xxx_messageInfo_FilteredBlock.Size(m)
}
func (m *FilteredBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_FilteredBlock.DiscardUnknown(m)
}

var xxx_messageInfo_FilteredBlock proto.InternalMessageInfo

func (m *FilteredBlock) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *FilteredBlock) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *FilteredBlock) GetFilteredTransactions() []*FilteredTransaction {
	if m != nil {
		return m.FilteredTransactions
	}
	return nil
}



type FilteredTransaction struct {
	Txid             string            `protobuf:"bytes,1,opt,name=txid" json:"txid,omitempty"`
	Type             common.HeaderType `protobuf:"varint,2,opt,name=type,enum=common.HeaderType" json:"type,omitempty"`
	TxValidationCode TxValidationCode  `protobuf:"varint,3,opt,name=tx_validation_code,json=txValidationCode,enum=protos.TxValidationCode" json:"tx_validation_code,omitempty"`
	
	
	Data                 isFilteredTransaction_Data `protobuf_oneof:"Data"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *FilteredTransaction) Reset()         { *m = FilteredTransaction{} }
func (m *FilteredTransaction) String() string { return proto.CompactTextString(m) }
func (*FilteredTransaction) ProtoMessage()    {}
func (*FilteredTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{6}
}
func (m *FilteredTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredTransaction.Unmarshal(m, b)
}
func (m *FilteredTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredTransaction.Marshal(b, m, deterministic)
}
func (dst *FilteredTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredTransaction.Merge(dst, src)
}
func (m *FilteredTransaction) XXX_Size() int {
	return xxx_messageInfo_FilteredTransaction.Size(m)
}
func (m *FilteredTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_FilteredTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_FilteredTransaction proto.InternalMessageInfo

type isFilteredTransaction_Data interface {
	isFilteredTransaction_Data()
}

type FilteredTransaction_TransactionActions struct {
	TransactionActions *FilteredTransactionActions `protobuf:"bytes,4,opt,name=transaction_actions,json=transactionActions,oneof"`
}

func (*FilteredTransaction_TransactionActions) isFilteredTransaction_Data() {}

func (m *FilteredTransaction) GetData() isFilteredTransaction_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *FilteredTransaction) GetTxid() string {
	if m != nil {
		return m.Txid
	}
	return ""
}

func (m *FilteredTransaction) GetType() common.HeaderType {
	if m != nil {
		return m.Type
	}
	return common.HeaderType_MESSAGE
}

func (m *FilteredTransaction) GetTxValidationCode() TxValidationCode {
	if m != nil {
		return m.TxValidationCode
	}
	return TxValidationCode_VALID
}

func (m *FilteredTransaction) GetTransactionActions() *FilteredTransactionActions {
	if x, ok := m.GetData().(*FilteredTransaction_TransactionActions); ok {
		return x.TransactionActions
	}
	return nil
}


func (*FilteredTransaction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FilteredTransaction_OneofMarshaler, _FilteredTransaction_OneofUnmarshaler, _FilteredTransaction_OneofSizer, []interface{}{
		(*FilteredTransaction_TransactionActions)(nil),
	}
}

func _FilteredTransaction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FilteredTransaction)
	
	switch x := m.Data.(type) {
	case *FilteredTransaction_TransactionActions:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TransactionActions); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("FilteredTransaction.Data has unexpected type %T", x)
	}
	return nil
}

func _FilteredTransaction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*FilteredTransaction)
	switch tag {
	case 4: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FilteredTransactionActions)
		err := b.DecodeMessage(msg)
		m.Data = &FilteredTransaction_TransactionActions{msg}
		return true, err
	default:
		return false, nil
	}
}

func _FilteredTransaction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FilteredTransaction)
	
	switch x := m.Data.(type) {
	case *FilteredTransaction_TransactionActions:
		s := proto.Size(x.TransactionActions)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}



type FilteredTransactionActions struct {
	ChaincodeActions     []*FilteredChaincodeAction `protobuf:"bytes,1,rep,name=chaincode_actions,json=chaincodeActions" json:"chaincode_actions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *FilteredTransactionActions) Reset()         { *m = FilteredTransactionActions{} }
func (m *FilteredTransactionActions) String() string { return proto.CompactTextString(m) }
func (*FilteredTransactionActions) ProtoMessage()    {}
func (*FilteredTransactionActions) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{7}
}
func (m *FilteredTransactionActions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredTransactionActions.Unmarshal(m, b)
}
func (m *FilteredTransactionActions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredTransactionActions.Marshal(b, m, deterministic)
}
func (dst *FilteredTransactionActions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredTransactionActions.Merge(dst, src)
}
func (m *FilteredTransactionActions) XXX_Size() int {
	return xxx_messageInfo_FilteredTransactionActions.Size(m)
}
func (m *FilteredTransactionActions) XXX_DiscardUnknown() {
	xxx_messageInfo_FilteredTransactionActions.DiscardUnknown(m)
}

var xxx_messageInfo_FilteredTransactionActions proto.InternalMessageInfo

func (m *FilteredTransactionActions) GetChaincodeActions() []*FilteredChaincodeAction {
	if m != nil {
		return m.ChaincodeActions
	}
	return nil
}



type FilteredChaincodeAction struct {
	ChaincodeEvent       *ChaincodeEvent `protobuf:"bytes,1,opt,name=chaincode_event,json=chaincodeEvent" json:"chaincode_event,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *FilteredChaincodeAction) Reset()         { *m = FilteredChaincodeAction{} }
func (m *FilteredChaincodeAction) String() string { return proto.CompactTextString(m) }
func (*FilteredChaincodeAction) ProtoMessage()    {}
func (*FilteredChaincodeAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{8}
}
func (m *FilteredChaincodeAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredChaincodeAction.Unmarshal(m, b)
}
func (m *FilteredChaincodeAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredChaincodeAction.Marshal(b, m, deterministic)
}
func (dst *FilteredChaincodeAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredChaincodeAction.Merge(dst, src)
}
func (m *FilteredChaincodeAction) XXX_Size() int {
	return xxx_messageInfo_FilteredChaincodeAction.Size(m)
}
func (m *FilteredChaincodeAction) XXX_DiscardUnknown() {
	xxx_messageInfo_FilteredChaincodeAction.DiscardUnknown(m)
}

var xxx_messageInfo_FilteredChaincodeAction proto.InternalMessageInfo

func (m *FilteredChaincodeAction) GetChaincodeEvent() *ChaincodeEvent {
	if m != nil {
		return m.ChaincodeEvent
	}
	return nil
}


type SignedEvent struct {
	
	Signature []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	
	EventBytes           []byte   `protobuf:"bytes,2,opt,name=eventBytes,proto3" json:"eventBytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedEvent) Reset()         { *m = SignedEvent{} }
func (m *SignedEvent) String() string { return proto.CompactTextString(m) }
func (*SignedEvent) ProtoMessage()    {}
func (*SignedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{9}
}
func (m *SignedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedEvent.Unmarshal(m, b)
}
func (m *SignedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedEvent.Marshal(b, m, deterministic)
}
func (dst *SignedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedEvent.Merge(dst, src)
}
func (m *SignedEvent) XXX_Size() int {
	return xxx_messageInfo_SignedEvent.Size(m)
}
func (m *SignedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_SignedEvent proto.InternalMessageInfo

func (m *SignedEvent) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *SignedEvent) GetEventBytes() []byte {
	if m != nil {
		return m.EventBytes
	}
	return nil
}




type Event struct {
	
	
	
	
	
	
	
	Event isEvent_Event `protobuf_oneof:"Event"`
	
	Creator []byte `protobuf:"bytes,6,opt,name=creator,proto3" json:"creator,omitempty"`
	
	Timestamp *timestamp.Timestamp `protobuf:"bytes,8,opt,name=timestamp" json:"timestamp,omitempty"`
	
	
	TlsCertHash          []byte   `protobuf:"bytes,9,opt,name=tls_cert_hash,json=tlsCertHash,proto3" json:"tls_cert_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{10}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (dst *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(dst, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

type isEvent_Event interface {
	isEvent_Event()
}

type Event_Register struct {
	Register *Register `protobuf:"bytes,1,opt,name=register,oneof"`
}
type Event_Block struct {
	Block *common.Block `protobuf:"bytes,2,opt,name=block,oneof"`
}
type Event_ChaincodeEvent struct {
	ChaincodeEvent *ChaincodeEvent `protobuf:"bytes,3,opt,name=chaincode_event,json=chaincodeEvent,oneof"`
}
type Event_Rejection struct {
	Rejection *Rejection `protobuf:"bytes,4,opt,name=rejection,oneof"`
}
type Event_Unregister struct {
	Unregister *Unregister `protobuf:"bytes,5,opt,name=unregister,oneof"`
}
type Event_FilteredBlock struct {
	FilteredBlock *FilteredBlock `protobuf:"bytes,7,opt,name=filtered_block,json=filteredBlock,oneof"`
}

func (*Event_Register) isEvent_Event()       {}
func (*Event_Block) isEvent_Event()          {}
func (*Event_ChaincodeEvent) isEvent_Event() {}
func (*Event_Rejection) isEvent_Event()      {}
func (*Event_Unregister) isEvent_Event()     {}
func (*Event_FilteredBlock) isEvent_Event()  {}

func (m *Event) GetEvent() isEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *Event) GetRegister() *Register {
	if x, ok := m.GetEvent().(*Event_Register); ok {
		return x.Register
	}
	return nil
}

func (m *Event) GetBlock() *common.Block {
	if x, ok := m.GetEvent().(*Event_Block); ok {
		return x.Block
	}
	return nil
}

func (m *Event) GetChaincodeEvent() *ChaincodeEvent {
	if x, ok := m.GetEvent().(*Event_ChaincodeEvent); ok {
		return x.ChaincodeEvent
	}
	return nil
}

func (m *Event) GetRejection() *Rejection {
	if x, ok := m.GetEvent().(*Event_Rejection); ok {
		return x.Rejection
	}
	return nil
}

func (m *Event) GetUnregister() *Unregister {
	if x, ok := m.GetEvent().(*Event_Unregister); ok {
		return x.Unregister
	}
	return nil
}

func (m *Event) GetFilteredBlock() *FilteredBlock {
	if x, ok := m.GetEvent().(*Event_FilteredBlock); ok {
		return x.FilteredBlock
	}
	return nil
}

func (m *Event) GetCreator() []byte {
	if m != nil {
		return m.Creator
	}
	return nil
}

func (m *Event) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Event) GetTlsCertHash() []byte {
	if m != nil {
		return m.TlsCertHash
	}
	return nil
}


func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_Register)(nil),
		(*Event_Block)(nil),
		(*Event_ChaincodeEvent)(nil),
		(*Event_Rejection)(nil),
		(*Event_Unregister)(nil),
		(*Event_FilteredBlock)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	
	switch x := m.Event.(type) {
	case *Event_Register:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Register); err != nil {
			return err
		}
	case *Event_Block:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Block); err != nil {
			return err
		}
	case *Event_ChaincodeEvent:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChaincodeEvent); err != nil {
			return err
		}
	case *Event_Rejection:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Rejection); err != nil {
			return err
		}
	case *Event_Unregister:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Unregister); err != nil {
			return err
		}
	case *Event_FilteredBlock:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.FilteredBlock); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.Event has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Register)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Register{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(common.Block)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Block{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeEvent)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ChaincodeEvent{msg}
		return true, err
	case 4: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Rejection)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Rejection{msg}
		return true, err
	case 5: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Unregister)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Unregister{msg}
		return true, err
	case 7: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FilteredBlock)
		err := b.DecodeMessage(msg)
		m.Event = &Event_FilteredBlock{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	
	switch x := m.Event.(type) {
	case *Event_Register:
		s := proto.Size(x.Register)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Block:
		s := proto.Size(x.Block)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ChaincodeEvent:
		s := proto.Size(x.ChaincodeEvent)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Rejection:
		s := proto.Size(x.Rejection)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Unregister:
		s := proto.Size(x.Unregister)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_FilteredBlock:
		s := proto.Size(x.FilteredBlock)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type DeliverResponse struct {
	
	
	
	
	Type                 isDeliverResponse_Type `protobuf_oneof:"Type"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *DeliverResponse) Reset()         { *m = DeliverResponse{} }
func (m *DeliverResponse) String() string { return proto.CompactTextString(m) }
func (*DeliverResponse) ProtoMessage()    {}
func (*DeliverResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_events_288a7d6350978185, []int{11}
}
func (m *DeliverResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeliverResponse.Unmarshal(m, b)
}
func (m *DeliverResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeliverResponse.Marshal(b, m, deterministic)
}
func (dst *DeliverResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeliverResponse.Merge(dst, src)
}
func (m *DeliverResponse) XXX_Size() int {
	return xxx_messageInfo_DeliverResponse.Size(m)
}
func (m *DeliverResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeliverResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeliverResponse proto.InternalMessageInfo

type isDeliverResponse_Type interface {
	isDeliverResponse_Type()
}

type DeliverResponse_Status struct {
	Status common.Status `protobuf:"varint,1,opt,name=status,enum=common.Status,oneof"`
}
type DeliverResponse_Block struct {
	Block *common.Block `protobuf:"bytes,2,opt,name=block,oneof"`
}
type DeliverResponse_FilteredBlock struct {
	FilteredBlock *FilteredBlock `protobuf:"bytes,3,opt,name=filtered_block,json=filteredBlock,oneof"`
}

func (*DeliverResponse_Status) isDeliverResponse_Type()        {}
func (*DeliverResponse_Block) isDeliverResponse_Type()         {}
func (*DeliverResponse_FilteredBlock) isDeliverResponse_Type() {}

func (m *DeliverResponse) GetType() isDeliverResponse_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *DeliverResponse) GetStatus() common.Status {
	if x, ok := m.GetType().(*DeliverResponse_Status); ok {
		return x.Status
	}
	return common.Status_UNKNOWN
}

func (m *DeliverResponse) GetBlock() *common.Block {
	if x, ok := m.GetType().(*DeliverResponse_Block); ok {
		return x.Block
	}
	return nil
}

func (m *DeliverResponse) GetFilteredBlock() *FilteredBlock {
	if x, ok := m.GetType().(*DeliverResponse_FilteredBlock); ok {
		return x.FilteredBlock
	}
	return nil
}


func (*DeliverResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _DeliverResponse_OneofMarshaler, _DeliverResponse_OneofUnmarshaler, _DeliverResponse_OneofSizer, []interface{}{
		(*DeliverResponse_Status)(nil),
		(*DeliverResponse_Block)(nil),
		(*DeliverResponse_FilteredBlock)(nil),
	}
}

func _DeliverResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*DeliverResponse)
	
	switch x := m.Type.(type) {
	case *DeliverResponse_Status:
		b.EncodeVarint(1<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.Status))
	case *DeliverResponse_Block:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Block); err != nil {
			return err
		}
	case *DeliverResponse_FilteredBlock:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.FilteredBlock); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("DeliverResponse.Type has unexpected type %T", x)
	}
	return nil
}

func _DeliverResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*DeliverResponse)
	switch tag {
	case 1: 
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Type = &DeliverResponse_Status{common.Status(x)}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(common.Block)
		err := b.DecodeMessage(msg)
		m.Type = &DeliverResponse_Block{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FilteredBlock)
		err := b.DecodeMessage(msg)
		m.Type = &DeliverResponse_FilteredBlock{msg}
		return true, err
	default:
		return false, nil
	}
}

func _DeliverResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*DeliverResponse)
	
	switch x := m.Type.(type) {
	case *DeliverResponse_Status:
		n += 1 
		n += proto.SizeVarint(uint64(x.Status))
	case *DeliverResponse_Block:
		s := proto.Size(x.Block)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *DeliverResponse_FilteredBlock:
		s := proto.Size(x.FilteredBlock)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ChaincodeReg)(nil), "protos.ChaincodeReg")
	proto.RegisterType((*Interest)(nil), "protos.Interest")
	proto.RegisterType((*Register)(nil), "protos.Register")
	proto.RegisterType((*Rejection)(nil), "protos.Rejection")
	proto.RegisterType((*Unregister)(nil), "protos.Unregister")
	proto.RegisterType((*FilteredBlock)(nil), "protos.FilteredBlock")
	proto.RegisterType((*FilteredTransaction)(nil), "protos.FilteredTransaction")
	proto.RegisterType((*FilteredTransactionActions)(nil), "protos.FilteredTransactionActions")
	proto.RegisterType((*FilteredChaincodeAction)(nil), "protos.FilteredChaincodeAction")
	proto.RegisterType((*SignedEvent)(nil), "protos.SignedEvent")
	proto.RegisterType((*Event)(nil), "protos.Event")
	proto.RegisterType((*DeliverResponse)(nil), "protos.DeliverResponse")
	proto.RegisterEnum("protos.EventType", EventType_name, EventType_value)
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4



type EventsClient interface {
	
	Chat(ctx context.Context, opts ...grpc.CallOption) (Events_ChatClient, error)
}

type eventsClient struct {
	cc *grpc.ClientConn
}

func NewEventsClient(cc *grpc.ClientConn) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) Chat(ctx context.Context, opts ...grpc.CallOption) (Events_ChatClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Events_serviceDesc.Streams[0], c.cc, "/protos.Events/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventsChatClient{stream}
	return x, nil
}

type Events_ChatClient interface {
	Send(*SignedEvent) error
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventsChatClient struct {
	grpc.ClientStream
}

func (x *eventsChatClient) Send(m *SignedEvent) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventsChatClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}



type EventsServer interface {
	
	Chat(Events_ChatServer) error
}

func RegisterEventsServer(s *grpc.Server, srv EventsServer) {
	s.RegisterService(&_Events_serviceDesc, srv)
}

func _Events_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventsServer).Chat(&eventsChatServer{stream})
}

type Events_ChatServer interface {
	Send(*Event) error
	Recv() (*SignedEvent, error)
	grpc.ServerStream
}

type eventsChatServer struct {
	grpc.ServerStream
}

func (x *eventsChatServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventsChatServer) Recv() (*SignedEvent, error) {
	m := new(SignedEvent)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Events_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Events",
	HandlerType: (*EventsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _Events_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "peer/events.proto",
}



type DeliverClient interface {
	
	
	Deliver(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverClient, error)
	
	
	DeliverFiltered(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverFilteredClient, error)
}

type deliverClient struct {
	cc *grpc.ClientConn
}

func NewDeliverClient(cc *grpc.ClientConn) DeliverClient {
	return &deliverClient{cc}
}

func (c *deliverClient) Deliver(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Deliver_serviceDesc.Streams[0], c.cc, "/protos.Deliver/Deliver", opts...)
	if err != nil {
		return nil, err
	}
	x := &deliverDeliverClient{stream}
	return x, nil
}

type Deliver_DeliverClient interface {
	Send(*common.Envelope) error
	Recv() (*DeliverResponse, error)
	grpc.ClientStream
}

type deliverDeliverClient struct {
	grpc.ClientStream
}

func (x *deliverDeliverClient) Send(m *common.Envelope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *deliverDeliverClient) Recv() (*DeliverResponse, error) {
	m := new(DeliverResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deliverClient) DeliverFiltered(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverFilteredClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Deliver_serviceDesc.Streams[1], c.cc, "/protos.Deliver/DeliverFiltered", opts...)
	if err != nil {
		return nil, err
	}
	x := &deliverDeliverFilteredClient{stream}
	return x, nil
}

type Deliver_DeliverFilteredClient interface {
	Send(*common.Envelope) error
	Recv() (*DeliverResponse, error)
	grpc.ClientStream
}

type deliverDeliverFilteredClient struct {
	grpc.ClientStream
}

func (x *deliverDeliverFilteredClient) Send(m *common.Envelope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *deliverDeliverFilteredClient) Recv() (*DeliverResponse, error) {
	m := new(DeliverResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}



type DeliverServer interface {
	
	
	Deliver(Deliver_DeliverServer) error
	
	
	DeliverFiltered(Deliver_DeliverFilteredServer) error
}

func RegisterDeliverServer(s *grpc.Server, srv DeliverServer) {
	s.RegisterService(&_Deliver_serviceDesc, srv)
}

func _Deliver_Deliver_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DeliverServer).Deliver(&deliverDeliverServer{stream})
}

type Deliver_DeliverServer interface {
	Send(*DeliverResponse) error
	Recv() (*common.Envelope, error)
	grpc.ServerStream
}

type deliverDeliverServer struct {
	grpc.ServerStream
}

func (x *deliverDeliverServer) Send(m *DeliverResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *deliverDeliverServer) Recv() (*common.Envelope, error) {
	m := new(common.Envelope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Deliver_DeliverFiltered_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DeliverServer).DeliverFiltered(&deliverDeliverFilteredServer{stream})
}

type Deliver_DeliverFilteredServer interface {
	Send(*DeliverResponse) error
	Recv() (*common.Envelope, error)
	grpc.ServerStream
}

type deliverDeliverFilteredServer struct {
	grpc.ServerStream
}

func (x *deliverDeliverFilteredServer) Send(m *DeliverResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *deliverDeliverFilteredServer) Recv() (*common.Envelope, error) {
	m := new(common.Envelope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Deliver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Deliver",
	HandlerType: (*DeliverServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Deliver",
			Handler:       _Deliver_Deliver_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "DeliverFiltered",
			Handler:       _Deliver_DeliverFiltered_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "peer/events.proto",
}

func init() { proto.RegisterFile("peer/events.proto", fileDescriptor_events_288a7d6350978185) }

var fileDescriptor_events_288a7d6350978185 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x56, 0xdd, 0x72, 0xdb, 0x44,
	0x14, 0xb6, 0x62, 0xc7, 0xb1, 0x8e, 0xe3, 0xd4, 0xd9, 0xb4, 0xa9, 0xc6, 0x05, 0x5a, 0xc4, 0xc0,
	0x04, 0x2e, 0xec, 0x60, 0x3a, 0x0c, 0xd3, 0x0b, 0x98, 0xf8, 0x27, 0xc8, 0x34, 0x4d, 0x32, 0x1b,
	0x87, 0x8b, 0x5e, 0xa0, 0x59, 0xcb, 0x6b, 0x59, 0xad, 0x2c, 0x79, 0x76, 0xd7, 0x99, 0xe4, 0x11,
	0x78, 0x03, 0xde, 0x80, 0x19, 0x5e, 0x85, 0x17, 0xe2, 0x92, 0xd1, 0x6a, 0x57, 0x52, 0x6c, 0xda,
	0x21, 0x57, 0xf6, 0xf9, 0xfb, 0xf6, 0xfc, 0x7c, 0x67, 0x57, 0xb0, 0xbf, 0xa4, 0x94, 0x75, 0xe8,
	0x0d, 0x8d, 0x04, 0x6f, 0x2f, 0x59, 0x2c, 0x62, 0x54, 0x95, 0x3f, 0xbc, 0x75, 0xe0, 0xc5, 0x8b,
	0x45, 0x1c, 0x75, 0xd2, 0x9f, 0xd4, 0xd8, 0x7a, 0xee, 0xc7, 0xb1, 0x1f, 0xd2, 0x8e, 0x94, 0x26,
	0xab, 0x59, 0x47, 0x04, 0x0b, 0xca, 0x05, 0x59, 0x2c, 0x95, 0x43, 0x4b, 0x02, 0x7a, 0x73, 0x12,
	0x44, 0x5e, 0x3c, 0xa5, 0xae, 0x84, 0x56, 0xb6, 0x43, 0x69, 0x13, 0x8c, 0x44, 0x9c, 0x78, 0x22,
	0xd0, 0xa0, 0xf6, 0x25, 0xec, 0xf6, 0x75, 0x00, 0xa6, 0x3e, 0xfa, 0x1c, 0x76, 0x73, 0x80, 0x60,
	0x6a, 0x19, 0x2f, 0x8c, 0x23, 0x13, 0xd7, 0x33, 0xdd, 0x68, 0x8a, 0x3e, 0x05, 0x90, 0xc8, 0x6e,
	0x44, 0x16, 0xd4, 0xda, 0x92, 0x0e, 0xa6, 0xd4, 0x9c, 0x93, 0x05, 0xb5, 0xff, 0x34, 0xa0, 0x36,
	0x8a, 0x04, 0x65, 0x94, 0x0b, 0x74, 0xac, 0x7d, 0xc5, 0xdd, 0x92, 0x4a, 0xb0, 0xbd, 0xee, 0x7e,
	0x7a, 0x34, 0x6f, 0x0f, 0x13, 0xcb, 0xf8, 0x6e, 0x49, 0x55, 0x78, 0xf2, 0x17, 0x0d, 0x00, 0xe5,
	0x09, 0x30, 0xea, 0xbb, 0x41, 0x34, 0x8b, 0xe5, 0x29, 0xf5, 0xee, 0x63, 0x1d, 0x59, 0x4c, 0xd9,
	0x29, 0xe1, 0xa6, 0x57, 0x90, 0x47, 0xd1, 0x2c, 0x46, 0x16, 0xec, 0x48, 0xdd, 0x68, 0x60, 0x95,
	0x65, 0x82, 0x5a, 0xec, 0x99, 0xb0, 0xa3, 0x9c, 0xec, 0x97, 0x50, 0xc3, 0xd4, 0x0f, 0xb8, 0xa0,
	0x0c, 0x1d, 0x41, 0x35, 0x9d, 0x84, 0x65, 0xbc, 0x28, 0x1f, 0xd5, 0xbb, 0x4d, 0x7d, 0x94, 0x2e,
	0x05, 0x2b, 0xbb, 0xfd, 0x06, 0x4c, 0x4c, 0xdf, 0x51, 0xd9, 0x44, 0xf4, 0x05, 0x6c, 0x89, 0x5b,
	0x59, 0x57, 0xbd, 0x7b, 0xa0, 0x43, 0xc6, 0x79, 0x97, 0xf1, 0x96, 0xb8, 0x45, 0xcf, 0xc0, 0xa4,
	0x8c, 0xc5, 0xcc, 0x5d, 0x70, 0x5f, 0xf5, 0xab, 0x26, 0x15, 0x6f, 0xb8, 0x6f, 0x7f, 0x0f, 0x70,
	0x1d, 0xb1, 0x87, 0xa7, 0xf1, 0x87, 0x01, 0x8d, 0xd3, 0x20, 0x4c, 0xb4, 0xd3, 0x5e, 0x18, 0x7b,
	0xef, 0x93, 0xb9, 0x78, 0x73, 0x12, 0x45, 0x34, 0xcc, 0x07, 0x67, 0x2a, 0xcd, 0x68, 0x8a, 0x0e,
	0xa1, 0x1a, 0xad, 0x16, 0x13, 0xca, 0x64, 0x0a, 0x15, 0xac, 0x24, 0x74, 0x09, 0x4f, 0x66, 0x0a,
	0xc7, 0x2d, 0xf0, 0x83, 0x5b, 0x15, 0x99, 0xc1, 0x33, 0x9d, 0x81, 0x3e, 0xac, 0x58, 0xdd, 0xe3,
	0xd9, 0xa6, 0x92, 0xdb, 0xff, 0x18, 0x70, 0xf0, 0x1f, 0xde, 0x08, 0x41, 0x45, 0xdc, 0x66, 0xa9,
	0xc9, 0xff, 0xe8, 0x2b, 0xa8, 0x48, 0x6a, 0x6c, 0x49, 0x6a, 0xa0, 0xb6, 0x62, 0xbc, 0x43, 0xc9,
	0x94, 0x32, 0xc9, 0x0d, 0x69, 0x47, 0xa7, 0x80, 0xc4, 0xad, 0x7b, 0x43, 0xc2, 0x60, 0x4a, 0x12,
	0x30, 0x37, 0x99, 0xb6, 0x9c, 0xed, 0x5e, 0xd7, 0xca, 0x1a, 0x7f, 0xfb, 0x6b, 0xe6, 0xd0, 0x4f,
	0xd8, 0xd0, 0x14, 0x6b, 0x1a, 0x74, 0x0d, 0x07, 0x85, 0x22, 0xdd, 0xbc, 0xd6, 0x64, 0x82, 0xf6,
	0x47, 0x6a, 0x3d, 0x49, 0x3d, 0x9d, 0x12, 0x46, 0x62, 0x43, 0xdb, 0xab, 0x42, 0x65, 0x40, 0x04,
	0xb1, 0xdf, 0x41, 0xeb, 0xc3, 0xb1, 0xe8, 0x0c, 0xf6, 0x73, 0x6e, 0xeb, 0xa3, 0xd3, 0x41, 0x3f,
	0x5f, 0x3f, 0x3a, 0xa3, 0x78, 0x1a, 0x5c, 0xe0, 0xb8, 0x42, 0xb3, 0xdf, 0xc2, 0xd3, 0x0f, 0x38,
	0xa3, 0x9f, 0xe0, 0xd1, 0xda, 0x35, 0xa0, 0x38, 0x7a, 0xb8, 0xb1, 0x41, 0x72, 0x09, 0xf1, 0x9e,
	0x77, 0x4f, 0xb6, 0x5f, 0x43, 0xfd, 0x2a, 0xf0, 0x23, 0x3a, 0x95, 0x22, 0xfa, 0x04, 0x4c, 0x1e,
	0xf8, 0x11, 0x11, 0x2b, 0x96, 0x6e, 0xf1, 0x2e, 0xce, 0x15, 0xe8, 0x33, 0xb5, 0xe4, 0xbd, 0x3b,
	0x41, 0xb9, 0x9c, 0xe4, 0x2e, 0x2e, 0x68, 0xec, 0xbf, 0xcb, 0xb0, 0x9d, 0xe2, 0xb4, 0xa1, 0xa6,
	0xa9, 0xae, 0x12, 0xca, 0x08, 0xae, 0x37, 0xd1, 0x29, 0xe1, 0xcc, 0x07, 0x7d, 0x09, 0xdb, 0x93,
	0x84, 0xdb, 0x6a, 0xff, 0x1b, 0x9a, 0x1e, 0x92, 0xf0, 0x4e, 0x09, 0xa7, 0x56, 0x74, 0xb2, 0x59,
	0x6e, 0xf9, 0x63, 0xe5, 0x3a, 0xa5, 0xf5, 0x82, 0xd1, 0xb7, 0x60, 0x32, 0xbd, 0xd5, 0x8a, 0x0d,
	0xfb, 0x79, 0x6a, 0xca, 0xe0, 0x94, 0x70, 0xee, 0x85, 0x5e, 0x02, 0xac, 0xb2, 0xcd, 0xb5, 0xb6,
	0x65, 0x0c, 0xd2, 0x31, 0xf9, 0x4e, 0x3b, 0x25, 0x5c, 0xf0, 0x43, 0x3f, 0xc2, 0x5e, 0xb6, 0x6e,
	0x69, 0x6d, 0x3b, 0x32, 0xf2, 0xc9, 0x3a, 0x01, 0x74, 0x8d, 0x8d, 0xd9, 0xbd, 0x2d, 0x4f, 0x6e,
	0x36, 0x46, 0x89, 0x88, 0x99, 0x55, 0x95, 0x9d, 0xd6, 0x22, 0xfa, 0x01, 0xcc, 0xec, 0x45, 0xb0,
	0x6a, 0x12, 0xb4, 0xd5, 0x4e, 0xdf, 0x8c, 0xb6, 0x7e, 0x33, 0xda, 0x63, 0xed, 0x81, 0x73, 0x67,
	0x64, 0x43, 0x43, 0x84, 0xdc, 0xf5, 0x28, 0x13, 0xee, 0x9c, 0xf0, 0xb9, 0x65, 0x4a, 0xe4, 0xba,
	0x08, 0x79, 0x9f, 0x32, 0xe1, 0x10, 0x3e, 0xef, 0xed, 0xa8, 0x19, 0xda, 0x7f, 0x19, 0xf0, 0x68,
	0x40, 0xc3, 0xe0, 0x86, 0x32, 0x4c, 0xf9, 0x32, 0x8e, 0x38, 0x4d, 0xae, 0x2d, 0x2e, 0x88, 0x58,
	0x71, 0x75, 0xc5, 0xef, 0xe9, 0x41, 0x5d, 0x49, 0xad, 0x53, 0xc2, 0xca, 0xfe, 0x7f, 0x27, 0xba,
	0xd9, 0xa5, 0xf2, 0x43, 0xba, 0x94, 0xec, 0x63, 0x72, 0x79, 0x7c, 0x73, 0x0d, 0x66, 0xf6, 0xca,
	0xa0, 0x5d, 0xa8, 0xe1, 0xe1, 0xcf, 0xa3, 0xab, 0xf1, 0x10, 0x37, 0x4b, 0xc8, 0x84, 0xed, 0xde,
	0xd9, 0x45, 0xff, 0x75, 0xd3, 0x40, 0x0d, 0x30, 0xfb, 0xce, 0xc9, 0xe8, 0xbc, 0x7f, 0x31, 0x18,
	0x36, 0xb7, 0x12, 0x11, 0x0f, 0x7f, 0x19, 0xf6, 0xc7, 0xa3, 0x8b, 0xf3, 0x66, 0x19, 0xed, 0x43,
	0xe3, 0x74, 0x74, 0x36, 0x1e, 0xe2, 0xe1, 0x20, 0x0d, 0xa8, 0x74, 0x5f, 0x41, 0x55, 0xc2, 0x72,
	0x74, 0x0c, 0x95, 0xfe, 0x9c, 0x08, 0x94, 0x5d, 0xfe, 0x85, 0xb5, 0x69, 0x35, 0xee, 0xbd, 0x74,
	0x76, 0xe9, 0xc8, 0x38, 0x36, 0xba, 0xbf, 0x1b, 0xb0, 0xa3, 0xfa, 0x87, 0x5e, 0xe5, 0x7f, 0x9b,
	0xba, 0x13, 0xc3, 0xe8, 0x86, 0x86, 0xf1, 0x92, 0xb6, 0x9e, 0xea, 0xe8, 0xb5, 0x6e, 0xa7, 0x38,
	0xa8, 0x97, 0x8d, 0x41, 0xf7, 0xe2, 0xc1, 0x18, 0xbd, 0xdf, 0xc0, 0x8e, 0x99, 0xdf, 0x9e, 0xdf,
	0x2d, 0x29, 0x0b, 0xe9, 0xd4, 0xa7, 0xac, 0x3d, 0x23, 0x13, 0x16, 0x78, 0x3a, 0x2c, 0xf9, 0x6a,
	0xe8, 0x35, 0xd2, 0x5a, 0x2f, 0x89, 0xf7, 0x9e, 0xf8, 0xf4, 0xed, 0xd7, 0x7e, 0x20, 0xe6, 0xab,
	0x49, 0x72, 0x56, 0xa7, 0x10, 0xd9, 0x49, 0x23, 0xd3, 0xcf, 0x13, 0xde, 0x49, 0x22, 0x27, 0xe9,
	0xf7, 0xcc, 0x77, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x7f, 0xfc, 0x9e, 0x06, 0xeb, 0x08, 0x00,
	0x00,
}
