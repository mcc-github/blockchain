


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/golang/protobuf/ptypes/timestamp"
import common "github.com/mcc-github/blockchain/protos/common"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type TxValidationCode int32

const (
	TxValidationCode_VALID                        TxValidationCode = 0
	TxValidationCode_NIL_ENVELOPE                 TxValidationCode = 1
	TxValidationCode_BAD_PAYLOAD                  TxValidationCode = 2
	TxValidationCode_BAD_COMMON_HEADER            TxValidationCode = 3
	TxValidationCode_BAD_CREATOR_SIGNATURE        TxValidationCode = 4
	TxValidationCode_INVALID_ENDORSER_TRANSACTION TxValidationCode = 5
	TxValidationCode_INVALID_CONFIG_TRANSACTION   TxValidationCode = 6
	TxValidationCode_UNSUPPORTED_TX_PAYLOAD       TxValidationCode = 7
	TxValidationCode_BAD_PROPOSAL_TXID            TxValidationCode = 8
	TxValidationCode_DUPLICATE_TXID               TxValidationCode = 9
	TxValidationCode_ENDORSEMENT_POLICY_FAILURE   TxValidationCode = 10
	TxValidationCode_MVCC_READ_CONFLICT           TxValidationCode = 11
	TxValidationCode_PHANTOM_READ_CONFLICT        TxValidationCode = 12
	TxValidationCode_UNKNOWN_TX_TYPE              TxValidationCode = 13
	TxValidationCode_TARGET_CHAIN_NOT_FOUND       TxValidationCode = 14
	TxValidationCode_MARSHAL_TX_ERROR             TxValidationCode = 15
	TxValidationCode_NIL_TXACTION                 TxValidationCode = 16
	TxValidationCode_EXPIRED_CHAINCODE            TxValidationCode = 17
	TxValidationCode_CHAINCODE_VERSION_CONFLICT   TxValidationCode = 18
	TxValidationCode_BAD_HEADER_EXTENSION         TxValidationCode = 19
	TxValidationCode_BAD_CHANNEL_HEADER           TxValidationCode = 20
	TxValidationCode_BAD_RESPONSE_PAYLOAD         TxValidationCode = 21
	TxValidationCode_BAD_RWSET                    TxValidationCode = 22
	TxValidationCode_ILLEGAL_WRITESET             TxValidationCode = 23
	TxValidationCode_INVALID_WRITESET             TxValidationCode = 24
	TxValidationCode_NOT_VALIDATED                TxValidationCode = 254
	TxValidationCode_INVALID_OTHER_REASON         TxValidationCode = 255
)

var TxValidationCode_name = map[int32]string{
	0:   "VALID",
	1:   "NIL_ENVELOPE",
	2:   "BAD_PAYLOAD",
	3:   "BAD_COMMON_HEADER",
	4:   "BAD_CREATOR_SIGNATURE",
	5:   "INVALID_ENDORSER_TRANSACTION",
	6:   "INVALID_CONFIG_TRANSACTION",
	7:   "UNSUPPORTED_TX_PAYLOAD",
	8:   "BAD_PROPOSAL_TXID",
	9:   "DUPLICATE_TXID",
	10:  "ENDORSEMENT_POLICY_FAILURE",
	11:  "MVCC_READ_CONFLICT",
	12:  "PHANTOM_READ_CONFLICT",
	13:  "UNKNOWN_TX_TYPE",
	14:  "TARGET_CHAIN_NOT_FOUND",
	15:  "MARSHAL_TX_ERROR",
	16:  "NIL_TXACTION",
	17:  "EXPIRED_CHAINCODE",
	18:  "CHAINCODE_VERSION_CONFLICT",
	19:  "BAD_HEADER_EXTENSION",
	20:  "BAD_CHANNEL_HEADER",
	21:  "BAD_RESPONSE_PAYLOAD",
	22:  "BAD_RWSET",
	23:  "ILLEGAL_WRITESET",
	24:  "INVALID_WRITESET",
	254: "NOT_VALIDATED",
	255: "INVALID_OTHER_REASON",
}
var TxValidationCode_value = map[string]int32{
	"VALID":                        0,
	"NIL_ENVELOPE":                 1,
	"BAD_PAYLOAD":                  2,
	"BAD_COMMON_HEADER":            3,
	"BAD_CREATOR_SIGNATURE":        4,
	"INVALID_ENDORSER_TRANSACTION": 5,
	"INVALID_CONFIG_TRANSACTION":   6,
	"UNSUPPORTED_TX_PAYLOAD":       7,
	"BAD_PROPOSAL_TXID":            8,
	"DUPLICATE_TXID":               9,
	"ENDORSEMENT_POLICY_FAILURE":   10,
	"MVCC_READ_CONFLICT":           11,
	"PHANTOM_READ_CONFLICT":        12,
	"UNKNOWN_TX_TYPE":              13,
	"TARGET_CHAIN_NOT_FOUND":       14,
	"MARSHAL_TX_ERROR":             15,
	"NIL_TXACTION":                 16,
	"EXPIRED_CHAINCODE":            17,
	"CHAINCODE_VERSION_CONFLICT":   18,
	"BAD_HEADER_EXTENSION":         19,
	"BAD_CHANNEL_HEADER":           20,
	"BAD_RESPONSE_PAYLOAD":         21,
	"BAD_RWSET":                    22,
	"ILLEGAL_WRITESET":             23,
	"INVALID_WRITESET":             24,
	"NOT_VALIDATED":                254,
	"INVALID_OTHER_REASON":         255,
}

func (x TxValidationCode) String() string {
	return proto.EnumName(TxValidationCode_name, int32(x))
}
func (TxValidationCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_transaction_0011a51fc2808780, []int{0}
}




type SignedTransaction struct {
	
	TransactionBytes []byte `protobuf:"bytes,1,opt,name=transaction_bytes,json=transactionBytes,proto3" json:"transaction_bytes,omitempty"`
	
	
	
	
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedTransaction) Reset()         { *m = SignedTransaction{} }
func (m *SignedTransaction) String() string { return proto.CompactTextString(m) }
func (*SignedTransaction) ProtoMessage()    {}
func (*SignedTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_0011a51fc2808780, []int{0}
}
func (m *SignedTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedTransaction.Unmarshal(m, b)
}
func (m *SignedTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedTransaction.Marshal(b, m, deterministic)
}
func (dst *SignedTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedTransaction.Merge(dst, src)
}
func (m *SignedTransaction) XXX_Size() int {
	return xxx_messageInfo_SignedTransaction.Size(m)
}
func (m *SignedTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_SignedTransaction proto.InternalMessageInfo

func (m *SignedTransaction) GetTransactionBytes() []byte {
	if m != nil {
		return m.TransactionBytes
	}
	return nil
}

func (m *SignedTransaction) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}







type ProcessedTransaction struct {
	
	TransactionEnvelope *common.Envelope `protobuf:"bytes,1,opt,name=transactionEnvelope" json:"transactionEnvelope,omitempty"`
	
	ValidationCode       int32    `protobuf:"varint,2,opt,name=validationCode" json:"validationCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProcessedTransaction) Reset()         { *m = ProcessedTransaction{} }
func (m *ProcessedTransaction) String() string { return proto.CompactTextString(m) }
func (*ProcessedTransaction) ProtoMessage()    {}
func (*ProcessedTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_0011a51fc2808780, []int{1}
}
func (m *ProcessedTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessedTransaction.Unmarshal(m, b)
}
func (m *ProcessedTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessedTransaction.Marshal(b, m, deterministic)
}
func (dst *ProcessedTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessedTransaction.Merge(dst, src)
}
func (m *ProcessedTransaction) XXX_Size() int {
	return xxx_messageInfo_ProcessedTransaction.Size(m)
}
func (m *ProcessedTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessedTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessedTransaction proto.InternalMessageInfo

func (m *ProcessedTransaction) GetTransactionEnvelope() *common.Envelope {
	if m != nil {
		return m.TransactionEnvelope
	}
	return nil
}

func (m *ProcessedTransaction) GetValidationCode() int32 {
	if m != nil {
		return m.ValidationCode
	}
	return 0
}













type Transaction struct {
	
	
	Actions              []*TransactionAction `protobuf:"bytes,1,rep,name=actions" json:"actions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_0011a51fc2808780, []int{2}
}
func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (dst *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(dst, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetActions() []*TransactionAction {
	if m != nil {
		return m.Actions
	}
	return nil
}



type TransactionAction struct {
	
	Header []byte `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	
	
	Payload              []byte   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransactionAction) Reset()         { *m = TransactionAction{} }
func (m *TransactionAction) String() string { return proto.CompactTextString(m) }
func (*TransactionAction) ProtoMessage()    {}
func (*TransactionAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_0011a51fc2808780, []int{3}
}
func (m *TransactionAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionAction.Unmarshal(m, b)
}
func (m *TransactionAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionAction.Marshal(b, m, deterministic)
}
func (dst *TransactionAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionAction.Merge(dst, src)
}
func (m *TransactionAction) XXX_Size() int {
	return xxx_messageInfo_TransactionAction.Size(m)
}
func (m *TransactionAction) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionAction.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionAction proto.InternalMessageInfo

func (m *TransactionAction) GetHeader() []byte {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *TransactionAction) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}




type ChaincodeActionPayload struct {
	
	
	
	
	
	
	
	
	
	ChaincodeProposalPayload []byte `protobuf:"bytes,1,opt,name=chaincode_proposal_payload,json=chaincodeProposalPayload,proto3" json:"chaincode_proposal_payload,omitempty"`
	
	Action               *ChaincodeEndorsedAction `protobuf:"bytes,2,opt,name=action" json:"action,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ChaincodeActionPayload) Reset()         { *m = ChaincodeActionPayload{} }
func (m *ChaincodeActionPayload) String() string { return proto.CompactTextString(m) }
func (*ChaincodeActionPayload) ProtoMessage()    {}
func (*ChaincodeActionPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_0011a51fc2808780, []int{4}
}
func (m *ChaincodeActionPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeActionPayload.Unmarshal(m, b)
}
func (m *ChaincodeActionPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeActionPayload.Marshal(b, m, deterministic)
}
func (dst *ChaincodeActionPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeActionPayload.Merge(dst, src)
}
func (m *ChaincodeActionPayload) XXX_Size() int {
	return xxx_messageInfo_ChaincodeActionPayload.Size(m)
}
func (m *ChaincodeActionPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeActionPayload.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeActionPayload proto.InternalMessageInfo

func (m *ChaincodeActionPayload) GetChaincodeProposalPayload() []byte {
	if m != nil {
		return m.ChaincodeProposalPayload
	}
	return nil
}

func (m *ChaincodeActionPayload) GetAction() *ChaincodeEndorsedAction {
	if m != nil {
		return m.Action
	}
	return nil
}



type ChaincodeEndorsedAction struct {
	
	
	
	ProposalResponsePayload []byte `protobuf:"bytes,1,opt,name=proposal_response_payload,json=proposalResponsePayload,proto3" json:"proposal_response_payload,omitempty"`
	
	
	Endorsements         []*Endorsement `protobuf:"bytes,2,rep,name=endorsements" json:"endorsements,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ChaincodeEndorsedAction) Reset()         { *m = ChaincodeEndorsedAction{} }
func (m *ChaincodeEndorsedAction) String() string { return proto.CompactTextString(m) }
func (*ChaincodeEndorsedAction) ProtoMessage()    {}
func (*ChaincodeEndorsedAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_0011a51fc2808780, []int{5}
}
func (m *ChaincodeEndorsedAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeEndorsedAction.Unmarshal(m, b)
}
func (m *ChaincodeEndorsedAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeEndorsedAction.Marshal(b, m, deterministic)
}
func (dst *ChaincodeEndorsedAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeEndorsedAction.Merge(dst, src)
}
func (m *ChaincodeEndorsedAction) XXX_Size() int {
	return xxx_messageInfo_ChaincodeEndorsedAction.Size(m)
}
func (m *ChaincodeEndorsedAction) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeEndorsedAction.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeEndorsedAction proto.InternalMessageInfo

func (m *ChaincodeEndorsedAction) GetProposalResponsePayload() []byte {
	if m != nil {
		return m.ProposalResponsePayload
	}
	return nil
}

func (m *ChaincodeEndorsedAction) GetEndorsements() []*Endorsement {
	if m != nil {
		return m.Endorsements
	}
	return nil
}

func init() {
	proto.RegisterType((*SignedTransaction)(nil), "protos.SignedTransaction")
	proto.RegisterType((*ProcessedTransaction)(nil), "protos.ProcessedTransaction")
	proto.RegisterType((*Transaction)(nil), "protos.Transaction")
	proto.RegisterType((*TransactionAction)(nil), "protos.TransactionAction")
	proto.RegisterType((*ChaincodeActionPayload)(nil), "protos.ChaincodeActionPayload")
	proto.RegisterType((*ChaincodeEndorsedAction)(nil), "protos.ChaincodeEndorsedAction")
	proto.RegisterEnum("protos.TxValidationCode", TxValidationCode_name, TxValidationCode_value)
}

func init() { proto.RegisterFile("peer/transaction.proto", fileDescriptor_transaction_0011a51fc2808780) }

var fileDescriptor_transaction_0011a51fc2808780 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0x5d, 0x6f, 0x22, 0x37,
	0x14, 0x2d, 0xd9, 0x26, 0x69, 0x2e, 0xf9, 0x30, 0x86, 0x10, 0x82, 0xa2, 0xee, 0x8a, 0x87, 0x6a,
	0xdb, 0x4a, 0x20, 0x65, 0x1f, 0x2a, 0x55, 0x7d, 0x31, 0x33, 0x4e, 0x18, 0x75, 0xb0, 0x47, 0x1e,
	0x43, 0x48, 0x1f, 0x6a, 0x0d, 0xe0, 0x25, 0xa8, 0x30, 0x83, 0x66, 0xc8, 0xaa, 0x79, 0xed, 0x0f,
	0x68, 0x5f, 0xfa, 0x7b, 0xdb, 0xca, 0xf3, 0x01, 0x24, 0xe9, 0xbe, 0x30, 0xf8, 0xdc, 0xe3, 0x7b,
	0xce, 0xbd, 0xd7, 0xba, 0x50, 0x5f, 0x69, 0x1d, 0x77, 0xd6, 0x71, 0x10, 0x26, 0xc1, 0x64, 0x3d,
	0x8f, 0xc2, 0xf6, 0x2a, 0x8e, 0xd6, 0x11, 0x3e, 0x48, 0x3f, 0x49, 0xf3, 0xed, 0x2c, 0x8a, 0x66,
	0x0b, 0xdd, 0x49, 0x8f, 0xe3, 0xc7, 0x8f, 0x9d, 0xf5, 0x7c, 0xa9, 0x93, 0x75, 0xb0, 0x5c, 0x65,
	0xc4, 0xe6, 0x55, 0x9a, 0x60, 0x15, 0x47, 0xab, 0x28, 0x09, 0x16, 0x2a, 0xd6, 0xc9, 0x2a, 0x0a,
	0x13, 0x9d, 0x47, 0xab, 0x93, 0x68, 0xb9, 0x8c, 0xc2, 0x4e, 0xf6, 0xc9, 0xc0, 0xd6, 0xaf, 0x50,
	0xf1, 0xe7, 0xb3, 0x50, 0x4f, 0xe5, 0x56, 0x16, 0x7f, 0x0f, 0x95, 0x1d, 0x17, 0x6a, 0xfc, 0xb4,
	0xd6, 0x49, 0xa3, 0xf4, 0xae, 0xf4, 0xfe, 0x58, 0xa0, 0x9d, 0x40, 0xd7, 0xe0, 0xf8, 0x0a, 0x8e,
	0x92, 0xf9, 0x2c, 0x0c, 0xd6, 0x8f, 0xb1, 0x6e, 0xec, 0xa5, 0xa4, 0x2d, 0xd0, 0xfa, 0xa3, 0x04,
	0x35, 0x2f, 0x8e, 0x26, 0x3a, 0x49, 0x9e, 0x6b, 0x74, 0xa1, 0xba, 0x93, 0x8a, 0x86, 0x9f, 0xf4,
	0x22, 0x5a, 0xe9, 0x54, 0xa5, 0x7c, 0x8d, 0xda, 0xb9, 0xc9, 0x02, 0x17, 0xff, 0x47, 0xc6, 0xdf,
	0xc0, 0xe9, 0xa7, 0x60, 0x31, 0x9f, 0x06, 0x06, 0xb5, 0xa2, 0x69, 0xa6, 0xbf, 0x2f, 0x5e, 0xa0,
	0xad, 0x2e, 0x94, 0x77, 0xa5, 0x3f, 0xc0, 0x61, 0xf6, 0xcf, 0x14, 0xf5, 0xe6, 0x7d, 0xf9, 0xfa,
	0x32, 0x6b, 0x46, 0xd2, 0xde, 0x61, 0x91, 0xf4, 0x57, 0x14, 0xcc, 0x16, 0x85, 0xca, 0xab, 0x28,
	0xae, 0xc3, 0xc1, 0x83, 0x0e, 0xa6, 0x3a, 0xce, 0xbb, 0x93, 0x9f, 0x70, 0x03, 0x0e, 0x57, 0xc1,
	0xd3, 0x22, 0x0a, 0xa6, 0x79, 0x47, 0x8a, 0x63, 0xeb, 0xaf, 0x12, 0xd4, 0xad, 0x87, 0x60, 0x1e,
	0x4e, 0xa2, 0xa9, 0xce, 0xb2, 0x78, 0x59, 0x08, 0xff, 0x04, 0xcd, 0x49, 0x11, 0x51, 0x9b, 0x21,
	0x16, 0x79, 0x32, 0x81, 0xc6, 0x86, 0xe1, 0xe5, 0x84, 0xe2, 0xf6, 0x0f, 0x70, 0x90, 0x59, 0x4b,
	0x15, 0xcb, 0xd7, 0x6f, 0x8b, 0x9a, 0x36, 0x6a, 0x34, 0x9c, 0x46, 0x71, 0xa2, 0xa7, 0x79, 0x65,
	0x39, 0xbd, 0xf5, 0x67, 0x09, 0x2e, 0x3e, 0xc3, 0xc1, 0x3f, 0xc2, 0xe5, 0xab, 0xd7, 0xf4, 0xc2,
	0xd1, 0x45, 0x41, 0x10, 0x79, 0x7c, 0x6b, 0xe8, 0x58, 0x67, 0xd9, 0x96, 0x3a, 0x5c, 0x27, 0x8d,
	0xbd, 0xb4, 0xd5, 0xd5, 0xc2, 0x16, 0xdd, 0xc6, 0xc4, 0x33, 0xe2, 0x77, 0x7f, 0xef, 0x03, 0x92,
	0xbf, 0x0f, 0x9f, 0x8d, 0x10, 0x1f, 0xc1, 0xfe, 0x90, 0xb8, 0x8e, 0x8d, 0xbe, 0xc0, 0x08, 0x8e,
	0x99, 0xe3, 0x2a, 0xca, 0x86, 0xd4, 0xe5, 0x1e, 0x45, 0x25, 0x7c, 0x06, 0xe5, 0x2e, 0xb1, 0x95,
	0x47, 0xee, 0x5d, 0x4e, 0x6c, 0xb4, 0x87, 0xcf, 0xa1, 0x62, 0x00, 0x8b, 0xf7, 0xfb, 0x9c, 0xa9,
	0x1e, 0x25, 0x36, 0x15, 0xe8, 0x0d, 0xbe, 0x84, 0xf3, 0x14, 0x16, 0x94, 0x48, 0x2e, 0x94, 0xef,
	0xdc, 0x32, 0x22, 0x07, 0x82, 0xa2, 0x2f, 0xf1, 0x3b, 0xb8, 0x72, 0x58, 0xaa, 0xa0, 0x28, 0xb3,
	0xb9, 0xf0, 0xa9, 0x50, 0x52, 0x10, 0xe6, 0x13, 0x4b, 0x3a, 0x9c, 0xa1, 0x7d, 0xfc, 0x35, 0x34,
	0x0b, 0x86, 0xc5, 0xd9, 0x8d, 0x73, 0xfb, 0x2c, 0x7e, 0x80, 0x9b, 0x50, 0x1f, 0x30, 0x7f, 0xe0,
	0x79, 0x5c, 0x48, 0x6a, 0x2b, 0x39, 0xda, 0xf8, 0x39, 0x2c, 0xfc, 0x78, 0x82, 0x7b, 0xdc, 0x27,
	0xae, 0x92, 0x23, 0xc7, 0x46, 0x5f, 0x61, 0x0c, 0xa7, 0xf6, 0xc0, 0x73, 0x1d, 0x8b, 0x48, 0x9a,
	0x61, 0x47, 0x46, 0x26, 0x37, 0xd0, 0xa7, 0x4c, 0x2a, 0x8f, 0xbb, 0x8e, 0x75, 0xaf, 0x6e, 0x88,
	0xe3, 0x1a, 0xa3, 0x80, 0xeb, 0x80, 0xfb, 0x43, 0xcb, 0x52, 0x82, 0x92, 0xcc, 0x88, 0xeb, 0x58,
	0x12, 0x95, 0x4d, 0x6d, 0x5e, 0x8f, 0x30, 0xc9, 0xfb, 0x2f, 0x42, 0xc7, 0xb8, 0x0a, 0x67, 0x03,
	0xf6, 0x33, 0xe3, 0x77, 0xcc, 0xb8, 0x92, 0xf7, 0x1e, 0x45, 0x27, 0xc6, 0xae, 0x24, 0xe2, 0x96,
	0x4a, 0x65, 0xf5, 0x88, 0xc3, 0x14, 0xe3, 0x52, 0xdd, 0xf0, 0x01, 0xb3, 0xd1, 0x29, 0xae, 0x01,
	0xea, 0x13, 0xe1, 0xf7, 0x52, 0xa7, 0x8a, 0x0a, 0xc1, 0x05, 0x3a, 0x2b, 0xfa, 0x2e, 0x47, 0x79,
	0xc9, 0xc8, 0x94, 0x45, 0x47, 0x9e, 0x23, 0xa8, 0x9d, 0x25, 0xb1, 0xb8, 0x4d, 0x51, 0xc5, 0x94,
	0xb0, 0x39, 0xaa, 0x21, 0x15, 0xbe, 0xc3, 0xd9, 0xd6, 0x0f, 0xc6, 0x0d, 0xa8, 0x99, 0x6e, 0x64,
	0x63, 0x51, 0x74, 0x24, 0x29, 0x33, 0x14, 0x54, 0x35, 0xc5, 0xa5, 0x03, 0xea, 0x11, 0xc6, 0xa8,
	0x5b, 0x0c, 0xae, 0x56, 0xdc, 0x10, 0xd4, 0xf7, 0x38, 0xf3, 0xe9, 0xa6, 0xb3, 0xe7, 0xf8, 0x04,
	0x8e, 0xd2, 0xc8, 0x9d, 0x4f, 0x25, 0xaa, 0x1b, 0xe7, 0x8e, 0xeb, 0xd2, 0x5b, 0xe2, 0xaa, 0x3b,
	0xe1, 0x48, 0x6a, 0xd0, 0x8b, 0x14, 0xcd, 0x47, 0xb7, 0x41, 0x1b, 0x18, 0xc3, 0x89, 0x29, 0x3a,
	0xc5, 0x89, 0xa4, 0x36, 0xfa, 0xa7, 0x84, 0x2f, 0xa1, 0x56, 0x30, 0xb9, 0xec, 0x51, 0x61, 0x7a,
	0xe9, 0x73, 0x86, 0xfe, 0x2d, 0x75, 0x27, 0xd0, 0x8a, 0xe2, 0x59, 0xfb, 0xe1, 0x69, 0xa5, 0xe3,
	0x85, 0x9e, 0xce, 0x74, 0xdc, 0xfe, 0x18, 0x8c, 0xe3, 0xf9, 0xa4, 0x78, 0xd1, 0x66, 0xf9, 0x76,
	0xf1, 0xce, 0x92, 0xf0, 0x82, 0xc9, 0x6f, 0xc1, 0x4c, 0xff, 0xf2, 0xed, 0x6c, 0xbe, 0x7e, 0x78,
	0x1c, 0x9b, 0x9d, 0xd6, 0xd9, 0xb9, 0xde, 0xc9, 0xae, 0x67, 0xeb, 0x3c, 0xe9, 0x98, 0xeb, 0xe3,
	0x6c, 0xd5, 0x7f, 0xf8, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x1b, 0x4a, 0x5e, 0xd5, 0x0b, 0x06, 0x00,
	0x00,
}
