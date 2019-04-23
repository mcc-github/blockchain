


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
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
	TxValidationCode_INVALID_CHAINCODE            TxValidationCode = 25
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
	25:  "INVALID_CHAINCODE",
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
	"INVALID_CHAINCODE":            25,
	"NOT_VALIDATED":                254,
	"INVALID_OTHER_REASON":         255,
}

func (x TxValidationCode) String() string {
	return proto.EnumName(TxValidationCode_name, int32(x))
}
func (TxValidationCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_transaction_6a3cf269faade5fe, []int{0}
}


type MetaDataKeys int32

const (
	MetaDataKeys_VALIDATION_PARAMETER    MetaDataKeys = 0
	MetaDataKeys_VALIDATION_PARAMETER_V2 MetaDataKeys = 1
)

var MetaDataKeys_name = map[int32]string{
	0: "VALIDATION_PARAMETER",
	1: "VALIDATION_PARAMETER_V2",
}
var MetaDataKeys_value = map[string]int32{
	"VALIDATION_PARAMETER":    0,
	"VALIDATION_PARAMETER_V2": 1,
}

func (x MetaDataKeys) String() string {
	return proto.EnumName(MetaDataKeys_name, int32(x))
}
func (MetaDataKeys) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_transaction_6a3cf269faade5fe, []int{1}
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
	return fileDescriptor_transaction_6a3cf269faade5fe, []int{0}
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
	
	TransactionEnvelope *common.Envelope `protobuf:"bytes,1,opt,name=transactionEnvelope,proto3" json:"transactionEnvelope,omitempty"`
	
	ValidationCode       int32    `protobuf:"varint,2,opt,name=validationCode,proto3" json:"validationCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProcessedTransaction) Reset()         { *m = ProcessedTransaction{} }
func (m *ProcessedTransaction) String() string { return proto.CompactTextString(m) }
func (*ProcessedTransaction) ProtoMessage()    {}
func (*ProcessedTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_6a3cf269faade5fe, []int{1}
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
	
	
	Actions              []*TransactionAction `protobuf:"bytes,1,rep,name=actions,proto3" json:"actions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_6a3cf269faade5fe, []int{2}
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
	return fileDescriptor_transaction_6a3cf269faade5fe, []int{3}
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
	
	Action               *ChaincodeEndorsedAction `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ChaincodeActionPayload) Reset()         { *m = ChaincodeActionPayload{} }
func (m *ChaincodeActionPayload) String() string { return proto.CompactTextString(m) }
func (*ChaincodeActionPayload) ProtoMessage()    {}
func (*ChaincodeActionPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_6a3cf269faade5fe, []int{4}
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
	
	
	Endorsements         []*Endorsement `protobuf:"bytes,2,rep,name=endorsements,proto3" json:"endorsements,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ChaincodeEndorsedAction) Reset()         { *m = ChaincodeEndorsedAction{} }
func (m *ChaincodeEndorsedAction) String() string { return proto.CompactTextString(m) }
func (*ChaincodeEndorsedAction) ProtoMessage()    {}
func (*ChaincodeEndorsedAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_6a3cf269faade5fe, []int{5}
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
	proto.RegisterEnum("protos.MetaDataKeys", MetaDataKeys_name, MetaDataKeys_value)
}

func init() { proto.RegisterFile("peer/transaction.proto", fileDescriptor_transaction_6a3cf269faade5fe) }

var fileDescriptor_transaction_6a3cf269faade5fe = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x55, 0x5d, 0x6f, 0xe2, 0x46,
	0x14, 0x5d, 0xb2, 0x4d, 0xd2, 0x5c, 0x48, 0x32, 0x0c, 0x84, 0x00, 0x8d, 0xda, 0x15, 0x0f, 0xd5,
	0x76, 0x2b, 0x81, 0x94, 0x7d, 0xa8, 0x54, 0xf5, 0x65, 0xb0, 0x27, 0xc1, 0x5a, 0x33, 0x63, 0x8d,
	0x07, 0x42, 0xfa, 0xd0, 0x91, 0x03, 0x53, 0x82, 0x4a, 0x6c, 0x64, 0xb3, 0xab, 0xe6, 0xb5, 0x3f,
	0xa0, 0xfd, 0x21, 0xfd, 0x8f, 0x6d, 0x35, 0xfe, 0x00, 0x92, 0x4d, 0x5f, 0x30, 0x73, 0xce, 0xb9,
	0xf7, 0x9e, 0x7b, 0x2f, 0x8c, 0xa1, 0xb1, 0xd2, 0x3a, 0xee, 0xad, 0xe3, 0x20, 0x4c, 0x82, 0xe9,
	0x7a, 0x11, 0x85, 0xdd, 0x55, 0x1c, 0xad, 0x23, 0x7c, 0x90, 0x3e, 0x92, 0xf6, 0x45, 0xca, 0xaf,
	0xe2, 0x68, 0x15, 0x25, 0xc1, 0x52, 0xc5, 0x3a, 0x59, 0x45, 0x61, 0xa2, 0x33, 0x55, 0xbb, 0x36,
	0x8d, 0x1e, 0x1e, 0xa2, 0xb0, 0x97, 0x3d, 0x32, 0xb0, 0xf3, 0x0b, 0x54, 0xfd, 0xc5, 0x3c, 0xd4,
	0x33, 0xb9, 0xcd, 0x8a, 0xbf, 0x87, 0xea, 0x4e, 0x11, 0x75, 0xf7, 0xb8, 0xd6, 0x49, 0xb3, 0xf4,
	0xa6, 0xf4, 0xb6, 0x22, 0xd0, 0x0e, 0xd1, 0x37, 0x38, 0xbe, 0x80, 0xa3, 0x64, 0x31, 0x0f, 0x83,
	0xf5, 0xc7, 0x58, 0x37, 0xf7, 0x52, 0xd1, 0x16, 0xe8, 0xfc, 0x51, 0x82, 0xba, 0x17, 0x47, 0x53,
	0x9d, 0x24, 0x4f, 0x6b, 0xf4, 0xa1, 0xb6, 0x93, 0x8a, 0x86, 0x9f, 0xf4, 0x32, 0x5a, 0xe9, 0xb4,
	0x4a, 0xf9, 0x12, 0x75, 0x73, 0x93, 0x05, 0x2e, 0x5e, 0x12, 0xe3, 0x6f, 0xe1, 0xe4, 0x53, 0xb0,
	0x5c, 0xcc, 0x02, 0x83, 0x5a, 0xd1, 0x2c, 0xab, 0xbf, 0x2f, 0x9e, 0xa1, 0x9d, 0x3e, 0x94, 0x77,
	0x4b, 0xbf, 0x87, 0xc3, 0xec, 0x9b, 0x69, 0xea, 0xf5, 0xdb, 0xf2, 0x65, 0x2b, 0x1b, 0x46, 0xd2,
	0xdd, 0x51, 0x91, 0xf4, 0x53, 0x14, 0xca, 0x0e, 0x85, 0xea, 0x67, 0x2c, 0x6e, 0xc0, 0xc1, 0xbd,
	0x0e, 0x66, 0x3a, 0xce, 0xa7, 0x93, 0x9f, 0x70, 0x13, 0x0e, 0x57, 0xc1, 0xe3, 0x32, 0x0a, 0x66,
	0xf9, 0x44, 0x8a, 0x63, 0xe7, 0xaf, 0x12, 0x34, 0xac, 0xfb, 0x60, 0x11, 0x4e, 0xa3, 0x99, 0xce,
	0xb2, 0x78, 0x19, 0x85, 0x7f, 0x82, 0xf6, 0xb4, 0x60, 0xd4, 0x66, 0x89, 0x45, 0x9e, 0xac, 0x40,
	0x73, 0xa3, 0xf0, 0x72, 0x41, 0x11, 0xfd, 0x03, 0x1c, 0x64, 0xd6, 0xd2, 0x8a, 0xe5, 0xcb, 0x6f,
	0x8a, 0x9e, 0x36, 0xd5, 0x68, 0x38, 0x8b, 0xe2, 0x44, 0xcf, 0xf2, 0xce, 0x72, 0x79, 0xe7, 0xcf,
	0x12, 0x9c, 0xff, 0x8f, 0x06, 0xff, 0x08, 0xad, 0xcf, 0x7e, 0x4d, 0xcf, 0x1c, 0x9d, 0x17, 0x02,
	0x91, 0xf3, 0x5b, 0x43, 0x15, 0x9d, 0x65, 0x7b, 0xd0, 0xe1, 0x3a, 0x69, 0xee, 0xa5, 0xa3, 0xae,
	0x15, 0xb6, 0xe8, 0x96, 0x13, 0x4f, 0x84, 0xef, 0xfe, 0xde, 0x07, 0x24, 0x7f, 0x1f, 0x3f, 0x59,
	0x21, 0x3e, 0x82, 0xfd, 0x31, 0x71, 0x1d, 0x1b, 0xbd, 0xc2, 0x08, 0x2a, 0xcc, 0x71, 0x15, 0x65,
	0x63, 0xea, 0x72, 0x8f, 0xa2, 0x12, 0x3e, 0x85, 0x72, 0x9f, 0xd8, 0xca, 0x23, 0xb7, 0x2e, 0x27,
	0x36, 0xda, 0xc3, 0x67, 0x50, 0x35, 0x80, 0xc5, 0x87, 0x43, 0xce, 0xd4, 0x80, 0x12, 0x9b, 0x0a,
	0xf4, 0x1a, 0xb7, 0xe0, 0x2c, 0x85, 0x05, 0x25, 0x92, 0x0b, 0xe5, 0x3b, 0xd7, 0x8c, 0xc8, 0x91,
	0xa0, 0xe8, 0x0b, 0xfc, 0x06, 0x2e, 0x1c, 0x96, 0x56, 0x50, 0x94, 0xd9, 0x5c, 0xf8, 0x54, 0x28,
	0x29, 0x08, 0xf3, 0x89, 0x25, 0x1d, 0xce, 0xd0, 0x3e, 0xfe, 0x1a, 0xda, 0x85, 0xc2, 0xe2, 0xec,
	0xca, 0xb9, 0x7e, 0xc2, 0x1f, 0xe0, 0x36, 0x34, 0x46, 0xcc, 0x1f, 0x79, 0x1e, 0x17, 0x92, 0xda,
	0x4a, 0x4e, 0x36, 0x7e, 0x0e, 0x0b, 0x3f, 0x9e, 0xe0, 0x1e, 0xf7, 0x89, 0xab, 0xe4, 0xc4, 0xb1,
	0xd1, 0x97, 0x18, 0xc3, 0x89, 0x3d, 0xf2, 0x5c, 0xc7, 0x22, 0x92, 0x66, 0xd8, 0x91, 0x29, 0x93,
	0x1b, 0x18, 0x52, 0x26, 0x95, 0xc7, 0x5d, 0xc7, 0xba, 0x55, 0x57, 0xc4, 0x71, 0x8d, 0x51, 0xc0,
	0x0d, 0xc0, 0xc3, 0xb1, 0x65, 0x29, 0x41, 0x49, 0x66, 0xc4, 0x75, 0x2c, 0x89, 0xca, 0xa6, 0x37,
	0x6f, 0x40, 0x98, 0xe4, 0xc3, 0x67, 0x54, 0x05, 0xd7, 0xe0, 0x74, 0xc4, 0x3e, 0x30, 0x7e, 0xc3,
	0x8c, 0x2b, 0x79, 0xeb, 0x51, 0x74, 0x6c, 0xec, 0x4a, 0x22, 0xae, 0xa9, 0x54, 0xd6, 0x80, 0x38,
	0x4c, 0x31, 0x2e, 0xd5, 0x15, 0x1f, 0x31, 0x1b, 0x9d, 0xe0, 0x3a, 0xa0, 0x21, 0x11, 0xfe, 0x20,
	0x75, 0xaa, 0xa8, 0x10, 0x5c, 0xa0, 0xd3, 0x62, 0xee, 0x72, 0x92, 0xb7, 0x8c, 0x4c, 0x5b, 0x74,
	0xe2, 0x39, 0x82, 0xda, 0x59, 0x12, 0x8b, 0xdb, 0x14, 0x55, 0x4d, 0x0b, 0x9b, 0xa3, 0x1a, 0x53,
	0xe1, 0x3b, 0x9c, 0x6d, 0xfd, 0x60, 0xdc, 0x84, 0xba, 0x99, 0x46, 0xb6, 0x16, 0x45, 0x27, 0x92,
	0x32, 0x23, 0x41, 0x35, 0xd3, 0x5c, 0xba, 0xa0, 0x01, 0x61, 0x8c, 0xba, 0xc5, 0xe2, 0xea, 0x45,
	0x84, 0xa0, 0xbe, 0xc7, 0x99, 0x4f, 0x37, 0x93, 0x3d, 0xc3, 0xc7, 0x70, 0x94, 0x32, 0x37, 0x3e,
	0x95, 0xa8, 0x61, 0x9c, 0x3b, 0xae, 0x4b, 0xaf, 0x89, 0xab, 0x6e, 0x84, 0x23, 0xa9, 0x41, 0xcf,
	0x53, 0x34, 0x5f, 0xdd, 0x06, 0x6d, 0x1a, 0xf7, 0x9b, 0x85, 0x6e, 0xdc, 0xb7, 0x30, 0x86, 0x63,
	0x33, 0x8b, 0x94, 0x20, 0x92, 0xda, 0xe8, 0x9f, 0x12, 0x6e, 0x41, 0xbd, 0x90, 0x72, 0x39, 0xa0,
	0xc2, 0x8c, 0xd8, 0xe7, 0x0c, 0xfd, 0x5b, 0x7a, 0x47, 0xa1, 0x32, 0xd4, 0xeb, 0xc0, 0x0e, 0xd6,
	0xc1, 0x07, 0xfd, 0x98, 0x18, 0xab, 0x79, 0xa8, 0xe9, 0xda, 0x23, 0x82, 0x0c, 0xa9, 0xa4, 0x02,
	0xbd, 0xc2, 0x5f, 0xc1, 0xf9, 0x4b, 0x8c, 0x1a, 0x5f, 0xa2, 0x52, 0x7f, 0x0a, 0x9d, 0x28, 0x9e,
	0x77, 0xef, 0x1f, 0x57, 0x3a, 0x5e, 0xea, 0xd9, 0x5c, 0xc7, 0xdd, 0x5f, 0x83, 0xbb, 0x78, 0x31,
	0x2d, 0xfe, 0x2f, 0xe6, 0x6a, 0xef, 0xe3, 0x9d, 0x2b, 0xc8, 0x0b, 0xa6, 0xbf, 0x05, 0x73, 0xfd,
	0xf3, 0x77, 0xf3, 0xc5, 0xfa, 0xfe, 0xe3, 0x9d, 0xb9, 0x31, 0x7b, 0x3b, 0xe1, 0xbd, 0x2c, 0xbc,
	0x97, 0x85, 0xf7, 0x4c, 0xf8, 0x5d, 0xf6, 0x9e, 0x78, 0xff, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x87, 0xb2, 0x02, 0x42, 0x48, 0x06, 0x00, 0x00,
}
