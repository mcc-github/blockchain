


package peer

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common "github.com/mcc-github/blockchain/protos/common"
	rwset "github.com/mcc-github/blockchain/protos/ledger/rwset"
	grpc "google.golang.org/grpc"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 


type FilteredBlock struct {
	ChannelId            string                 `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Number               uint64                 `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
	FilteredTransactions []*FilteredTransaction `protobuf:"bytes,4,rep,name=filtered_transactions,json=filteredTransactions,proto3" json:"filtered_transactions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *FilteredBlock) Reset()         { *m = FilteredBlock{} }
func (m *FilteredBlock) String() string { return proto.CompactTextString(m) }
func (*FilteredBlock) ProtoMessage()    {}
func (*FilteredBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_5eedcc5fab2714e6, []int{0}
}

func (m *FilteredBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredBlock.Unmarshal(m, b)
}
func (m *FilteredBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredBlock.Marshal(b, m, deterministic)
}
func (m *FilteredBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredBlock.Merge(m, src)
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
	Txid             string            `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Type             common.HeaderType `protobuf:"varint,2,opt,name=type,proto3,enum=common.HeaderType" json:"type,omitempty"`
	TxValidationCode TxValidationCode  `protobuf:"varint,3,opt,name=tx_validation_code,json=txValidationCode,proto3,enum=protos.TxValidationCode" json:"tx_validation_code,omitempty"`
	
	
	Data                 isFilteredTransaction_Data `protobuf_oneof:"Data"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *FilteredTransaction) Reset()         { *m = FilteredTransaction{} }
func (m *FilteredTransaction) String() string { return proto.CompactTextString(m) }
func (*FilteredTransaction) ProtoMessage()    {}
func (*FilteredTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_5eedcc5fab2714e6, []int{1}
}

func (m *FilteredTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredTransaction.Unmarshal(m, b)
}
func (m *FilteredTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredTransaction.Marshal(b, m, deterministic)
}
func (m *FilteredTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredTransaction.Merge(m, src)
}
func (m *FilteredTransaction) XXX_Size() int {
	return xxx_messageInfo_FilteredTransaction.Size(m)
}
func (m *FilteredTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_FilteredTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_FilteredTransaction proto.InternalMessageInfo

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

type isFilteredTransaction_Data interface {
	isFilteredTransaction_Data()
}

type FilteredTransaction_TransactionActions struct {
	TransactionActions *FilteredTransactionActions `protobuf:"bytes,4,opt,name=transaction_actions,json=transactionActions,proto3,oneof"`
}

func (*FilteredTransaction_TransactionActions) isFilteredTransaction_Data() {}

func (m *FilteredTransaction) GetData() isFilteredTransaction_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *FilteredTransaction) GetTransactionActions() *FilteredTransactionActions {
	if x, ok := m.GetData().(*FilteredTransaction_TransactionActions); ok {
		return x.TransactionActions
	}
	return nil
}


func (*FilteredTransaction) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*FilteredTransaction_TransactionActions)(nil),
	}
}



type FilteredTransactionActions struct {
	ChaincodeActions     []*FilteredChaincodeAction `protobuf:"bytes,1,rep,name=chaincode_actions,json=chaincodeActions,proto3" json:"chaincode_actions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *FilteredTransactionActions) Reset()         { *m = FilteredTransactionActions{} }
func (m *FilteredTransactionActions) String() string { return proto.CompactTextString(m) }
func (*FilteredTransactionActions) ProtoMessage()    {}
func (*FilteredTransactionActions) Descriptor() ([]byte, []int) {
	return fileDescriptor_5eedcc5fab2714e6, []int{2}
}

func (m *FilteredTransactionActions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredTransactionActions.Unmarshal(m, b)
}
func (m *FilteredTransactionActions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredTransactionActions.Marshal(b, m, deterministic)
}
func (m *FilteredTransactionActions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredTransactionActions.Merge(m, src)
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
	ChaincodeEvent       *ChaincodeEvent `protobuf:"bytes,1,opt,name=chaincode_event,json=chaincodeEvent,proto3" json:"chaincode_event,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *FilteredChaincodeAction) Reset()         { *m = FilteredChaincodeAction{} }
func (m *FilteredChaincodeAction) String() string { return proto.CompactTextString(m) }
func (*FilteredChaincodeAction) ProtoMessage()    {}
func (*FilteredChaincodeAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_5eedcc5fab2714e6, []int{3}
}

func (m *FilteredChaincodeAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredChaincodeAction.Unmarshal(m, b)
}
func (m *FilteredChaincodeAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredChaincodeAction.Marshal(b, m, deterministic)
}
func (m *FilteredChaincodeAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredChaincodeAction.Merge(m, src)
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


type BlockAndPrivateData struct {
	Block *common.Block `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	
	PrivateDataMap       map[uint64]*rwset.TxPvtReadWriteSet `protobuf:"bytes,2,rep,name=private_data_map,json=privateDataMap,proto3" json:"private_data_map,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_unrecognized     []byte                              `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *BlockAndPrivateData) Reset()         { *m = BlockAndPrivateData{} }
func (m *BlockAndPrivateData) String() string { return proto.CompactTextString(m) }
func (*BlockAndPrivateData) ProtoMessage()    {}
func (*BlockAndPrivateData) Descriptor() ([]byte, []int) {
	return fileDescriptor_5eedcc5fab2714e6, []int{4}
}

func (m *BlockAndPrivateData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockAndPrivateData.Unmarshal(m, b)
}
func (m *BlockAndPrivateData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockAndPrivateData.Marshal(b, m, deterministic)
}
func (m *BlockAndPrivateData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockAndPrivateData.Merge(m, src)
}
func (m *BlockAndPrivateData) XXX_Size() int {
	return xxx_messageInfo_BlockAndPrivateData.Size(m)
}
func (m *BlockAndPrivateData) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockAndPrivateData.DiscardUnknown(m)
}

var xxx_messageInfo_BlockAndPrivateData proto.InternalMessageInfo

func (m *BlockAndPrivateData) GetBlock() *common.Block {
	if m != nil {
		return m.Block
	}
	return nil
}

func (m *BlockAndPrivateData) GetPrivateDataMap() map[uint64]*rwset.TxPvtReadWriteSet {
	if m != nil {
		return m.PrivateDataMap
	}
	return nil
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
	return fileDescriptor_5eedcc5fab2714e6, []int{5}
}

func (m *DeliverResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeliverResponse.Unmarshal(m, b)
}
func (m *DeliverResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeliverResponse.Marshal(b, m, deterministic)
}
func (m *DeliverResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeliverResponse.Merge(m, src)
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
	Status common.Status `protobuf:"varint,1,opt,name=status,proto3,enum=common.Status,oneof"`
}

type DeliverResponse_Block struct {
	Block *common.Block `protobuf:"bytes,2,opt,name=block,proto3,oneof"`
}

type DeliverResponse_FilteredBlock struct {
	FilteredBlock *FilteredBlock `protobuf:"bytes,3,opt,name=filtered_block,json=filteredBlock,proto3,oneof"`
}

type DeliverResponse_BlockAndPrivateData struct {
	BlockAndPrivateData *BlockAndPrivateData `protobuf:"bytes,4,opt,name=block_and_private_data,json=blockAndPrivateData,proto3,oneof"`
}

func (*DeliverResponse_Status) isDeliverResponse_Type() {}

func (*DeliverResponse_Block) isDeliverResponse_Type() {}

func (*DeliverResponse_FilteredBlock) isDeliverResponse_Type() {}

func (*DeliverResponse_BlockAndPrivateData) isDeliverResponse_Type() {}

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

func (m *DeliverResponse) GetBlockAndPrivateData() *BlockAndPrivateData {
	if x, ok := m.GetType().(*DeliverResponse_BlockAndPrivateData); ok {
		return x.BlockAndPrivateData
	}
	return nil
}


func (*DeliverResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*DeliverResponse_Status)(nil),
		(*DeliverResponse_Block)(nil),
		(*DeliverResponse_FilteredBlock)(nil),
		(*DeliverResponse_BlockAndPrivateData)(nil),
	}
}

func init() {
	proto.RegisterType((*FilteredBlock)(nil), "protos.FilteredBlock")
	proto.RegisterType((*FilteredTransaction)(nil), "protos.FilteredTransaction")
	proto.RegisterType((*FilteredTransactionActions)(nil), "protos.FilteredTransactionActions")
	proto.RegisterType((*FilteredChaincodeAction)(nil), "protos.FilteredChaincodeAction")
	proto.RegisterType((*BlockAndPrivateData)(nil), "protos.BlockAndPrivateData")
	proto.RegisterMapType((map[uint64]*rwset.TxPvtReadWriteSet)(nil), "protos.BlockAndPrivateData.PrivateDataMapEntry")
	proto.RegisterType((*DeliverResponse)(nil), "protos.DeliverResponse")
}

func init() { proto.RegisterFile("peer/events.proto", fileDescriptor_5eedcc5fab2714e6) }

var fileDescriptor_5eedcc5fab2714e6 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4d, 0x6f, 0xda, 0x4a,
	0x14, 0xc5, 0x40, 0x78, 0xca, 0x20, 0x08, 0x19, 0x5e, 0x88, 0x45, 0xf4, 0xf4, 0x22, 0x57, 0xad,
	0xe8, 0xc6, 0xae, 0xe8, 0xa6, 0xca, 0xa2, 0x55, 0xc8, 0x87, 0x88, 0xd4, 0x4a, 0x68, 0x42, 0x1b,
	0x35, 0x95, 0x6a, 0x0d, 0xf6, 0x05, 0xdc, 0x18, 0xdb, 0x1a, 0x0f, 0x2e, 0xfc, 0x93, 0xfe, 0xb0,
	0xfe, 0x92, 0xae, 0xba, 0xaa, 0x2a, 0xcf, 0x78, 0xf8, 0x0a, 0x89, 0x94, 0x8d, 0x3d, 0xbe, 0xf7,
	0x9c, 0x73, 0xe7, 0x1e, 0xdf, 0x19, 0xb4, 0x1f, 0x01, 0x30, 0x0b, 0x12, 0x08, 0x78, 0x6c, 0x46,
	0x2c, 0xe4, 0x21, 0x2e, 0x89, 0x57, 0xdc, 0xac, 0x3b, 0xe1, 0x64, 0x12, 0x06, 0x96, 0x7c, 0xc9,
	0x64, 0x53, 0xf7, 0xc1, 0x1d, 0x01, 0xb3, 0xd8, 0xf7, 0x18, 0xb8, 0x7c, 0x66, 0x99, 0xa6, 0x50,
	0x72, 0xc6, 0xd4, 0x0b, 0x9c, 0xd0, 0x05, 0x5b, 0x68, 0x66, 0xb9, 0x86, 0xc8, 0x71, 0x46, 0x83,
	0x98, 0x3a, 0xdc, 0x53, 0x6a, 0xc6, 0x0f, 0x0d, 0x55, 0x2e, 0x3d, 0x9f, 0x03, 0x03, 0xb7, 0xe3,
	0x87, 0xce, 0x1d, 0xfe, 0x0f, 0x21, 0x67, 0x4c, 0x83, 0x00, 0x7c, 0xdb, 0x73, 0x75, 0xed, 0x58,
	0x6b, 0xed, 0x92, 0xdd, 0x2c, 0x72, 0xe5, 0xe2, 0x06, 0x2a, 0x05, 0xd3, 0xc9, 0x00, 0x98, 0x9e,
	0x3f, 0xd6, 0x5a, 0x45, 0x92, 0x7d, 0xe1, 0x1e, 0x3a, 0x18, 0x66, 0x3a, 0xf6, 0x4a, 0x99, 0x58,
	0x2f, 0x1e, 0x17, 0x5a, 0xe5, 0xf6, 0x91, 0xac, 0x17, 0x9b, 0xaa, 0x58, 0x7f, 0x89, 0x21, 0xff,
	0x0e, 0xef, 0x07, 0x63, 0xe3, 0xb7, 0x86, 0xea, 0x5b, 0xd0, 0x18, 0xa3, 0x22, 0x9f, 0x2d, 0xb6,
	0x26, 0xd6, 0xf8, 0x05, 0x2a, 0xf2, 0x79, 0x04, 0x62, 0x4f, 0xd5, 0x36, 0x36, 0x33, 0xc7, 0xba,
	0x40, 0x5d, 0x60, 0xfd, 0x79, 0x04, 0x44, 0xe4, 0xf1, 0x25, 0xc2, 0x7c, 0x66, 0x27, 0xd4, 0xf7,
	0x5c, 0x9a, 0x8a, 0xd9, 0xa9, 0x51, 0x7a, 0x41, 0xb0, 0x74, 0xb5, 0xc5, 0xfe, 0xec, 0xd3, 0x02,
	0x70, 0x16, 0xba, 0x40, 0x6a, 0x7c, 0x23, 0x82, 0x3f, 0xa2, 0xfa, 0x4a, 0x93, 0xf6, 0xb2, 0x57,
	0xad, 0x55, 0x6e, 0x1b, 0x8f, 0xf4, 0x7a, 0x2a, 0x91, 0xdd, 0x1c, 0xc1, 0xfc, 0x5e, 0xb4, 0x53,
	0x42, 0xc5, 0x73, 0xca, 0xa9, 0xf1, 0x0d, 0x35, 0x1f, 0xe6, 0xe2, 0xf7, 0x68, 0x7f, 0xf9, 0x93,
	0x55, 0x69, 0x4d, 0xd8, 0xfc, 0xff, 0x66, 0xe9, 0x33, 0x05, 0x94, 0x64, 0x52, 0x73, 0xd6, 0x03,
	0xb1, 0x71, 0x8b, 0x0e, 0x1f, 0x00, 0xe3, 0x77, 0x68, 0x6f, 0x63, 0x9a, 0x84, 0xe9, 0xe5, 0x76,
	0x43, 0x95, 0x59, 0x30, 0x2e, 0xd2, 0x2c, 0xa9, 0x3a, 0x6b, 0xdf, 0xc6, 0x2f, 0x0d, 0xd5, 0xc5,
	0x54, 0x9d, 0x06, 0x6e, 0x8f, 0x79, 0x09, 0xe5, 0x90, 0xf6, 0x87, 0x9f, 0xa1, 0x9d, 0x41, 0x1a,
	0xce, 0xe4, 0x2a, 0xea, 0x7f, 0x09, 0x2c, 0x91, 0x39, 0xfc, 0x19, 0xd5, 0x22, 0xc9, 0xb1, 0x5d,
	0xca, 0xa9, 0x3d, 0xa1, 0x91, 0x9e, 0x17, 0x5d, 0x5a, 0xaa, 0xfc, 0x16, 0x6d, 0x73, 0x65, 0xfd,
	0x81, 0x46, 0x17, 0x01, 0x67, 0x73, 0x52, 0x8d, 0xd6, 0x82, 0xcd, 0x2f, 0xa8, 0xbe, 0x05, 0x86,
	0x6b, 0xa8, 0x70, 0x07, 0x73, 0xb1, 0xa9, 0x22, 0x49, 0x97, 0xd8, 0x44, 0x3b, 0x09, 0xf5, 0xa7,
	0x72, 0xb0, 0xca, 0x6d, 0xdd, 0x94, 0xe7, 0xad, 0x3f, 0xeb, 0x25, 0x9c, 0x00, 0x75, 0x6f, 0x98,
	0xc7, 0xe1, 0x1a, 0x38, 0x91, 0xb0, 0x93, 0xfc, 0x1b, 0xcd, 0xf8, 0xa3, 0xa1, 0xbd, 0x73, 0xf0,
	0xbd, 0x04, 0x18, 0x81, 0x38, 0x0a, 0x83, 0x18, 0x70, 0x0b, 0x95, 0x62, 0x4e, 0xf9, 0x34, 0x16,
	0xe2, 0xd5, 0x76, 0x55, 0x75, 0x7c, 0x2d, 0xa2, 0xdd, 0x1c, 0xc9, 0xf2, 0xf8, 0xb9, 0xb2, 0x26,
	0xbf, 0xc5, 0x9a, 0x6e, 0x4e, 0x99, 0xf3, 0x16, 0x55, 0x17, 0xc7, 0x4d, 0xe2, 0x0b, 0x02, 0x7f,
	0xb0, 0x39, 0x00, 0x8a, 0x57, 0x19, 0xae, 0x9d, 0x72, 0x82, 0x1a, 0x82, 0x66, 0xd3, 0xc0, 0xb5,
	0x57, 0x6d, 0xce, 0x66, 0xf8, 0xe8, 0x11, 0x8b, 0xbb, 0x39, 0x52, 0x1f, 0xdc, 0x0f, 0xa7, 0xd3,
	0x9b, 0x1e, 0xb5, 0xf6, 0x4f, 0x0d, 0xfd, 0x93, 0x19, 0x80, 0x4f, 0x96, 0xcb, 0x9a, 0x6a, 0xe5,
	0x22, 0x48, 0xc0, 0x0f, 0x23, 0x68, 0x1e, 0xaa, 0x22, 0x1b, 0x76, 0x19, 0xb9, 0x96, 0xf6, 0x4a,
	0xc3, 0x9d, 0x85, 0x8f, 0xaa, 0x99, 0xa7, 0x6b, 0x5c, 0xa1, 0x46, 0x96, 0xb8, 0xf1, 0xf8, 0x78,
	0x75, 0x06, 0x9f, 0x2a, 0xd5, 0xf9, 0x8a, 0x8c, 0x90, 0x8d, 0xcc, 0xf1, 0x3c, 0x02, 0x26, 0xef,
	0x60, 0x73, 0x48, 0x07, 0xcc, 0x73, 0x14, 0x2d, 0xbd, 0x62, 0x3b, 0x15, 0x31, 0xf9, 0x71, 0x8f,
	0x3a, 0x77, 0x74, 0x04, 0xb7, 0x2f, 0x47, 0x1e, 0x1f, 0x4f, 0x07, 0x69, 0x2d, 0x6b, 0x85, 0x69,
	0x49, 0xa6, 0x25, 0x99, 0x56, 0xca, 0x1c, 0xc8, 0x5b, 0xff, 0xf5, 0xdf, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xf9, 0xaa, 0x6e, 0x86, 0x11, 0x06, 0x00, 0x00,
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type DeliverClient interface {
	
	
	
	Deliver(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverClient, error)
	
	
	
	DeliverFiltered(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverFilteredClient, error)
	
	
	
	DeliverWithPrivateData(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverWithPrivateDataClient, error)
}

type deliverClient struct {
	cc *grpc.ClientConn
}

func NewDeliverClient(cc *grpc.ClientConn) DeliverClient {
	return &deliverClient{cc}
}

func (c *deliverClient) Deliver(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Deliver_serviceDesc.Streams[0], "/protos.Deliver/Deliver", opts...)
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
	stream, err := c.cc.NewStream(ctx, &_Deliver_serviceDesc.Streams[1], "/protos.Deliver/DeliverFiltered", opts...)
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

func (c *deliverClient) DeliverWithPrivateData(ctx context.Context, opts ...grpc.CallOption) (Deliver_DeliverWithPrivateDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Deliver_serviceDesc.Streams[2], "/protos.Deliver/DeliverWithPrivateData", opts...)
	if err != nil {
		return nil, err
	}
	x := &deliverDeliverWithPrivateDataClient{stream}
	return x, nil
}

type Deliver_DeliverWithPrivateDataClient interface {
	Send(*common.Envelope) error
	Recv() (*DeliverResponse, error)
	grpc.ClientStream
}

type deliverDeliverWithPrivateDataClient struct {
	grpc.ClientStream
}

func (x *deliverDeliverWithPrivateDataClient) Send(m *common.Envelope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *deliverDeliverWithPrivateDataClient) Recv() (*DeliverResponse, error) {
	m := new(DeliverResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}


type DeliverServer interface {
	
	
	
	Deliver(Deliver_DeliverServer) error
	
	
	
	DeliverFiltered(Deliver_DeliverFilteredServer) error
	
	
	
	DeliverWithPrivateData(Deliver_DeliverWithPrivateDataServer) error
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

func _Deliver_DeliverWithPrivateData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DeliverServer).DeliverWithPrivateData(&deliverDeliverWithPrivateDataServer{stream})
}

type Deliver_DeliverWithPrivateDataServer interface {
	Send(*DeliverResponse) error
	Recv() (*common.Envelope, error)
	grpc.ServerStream
}

type deliverDeliverWithPrivateDataServer struct {
	grpc.ServerStream
}

func (x *deliverDeliverWithPrivateDataServer) Send(m *DeliverResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *deliverDeliverWithPrivateDataServer) Recv() (*common.Envelope, error) {
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
		{
			StreamName:    "DeliverWithPrivateData",
			Handler:       _Deliver_DeliverWithPrivateData_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "peer/events.proto",
}
