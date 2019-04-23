


package token 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 


type IssueRequest struct {
	
	
	Credential []byte `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	
	TokensToIssue        []*Token `protobuf:"bytes,2,rep,name=tokens_to_issue,json=tokensToIssue,proto3" json:"tokens_to_issue,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IssueRequest) Reset()         { *m = IssueRequest{} }
func (m *IssueRequest) String() string { return proto.CompactTextString(m) }
func (*IssueRequest) ProtoMessage()    {}
func (*IssueRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{0}
}
func (m *IssueRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IssueRequest.Unmarshal(m, b)
}
func (m *IssueRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IssueRequest.Marshal(b, m, deterministic)
}
func (dst *IssueRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IssueRequest.Merge(dst, src)
}
func (m *IssueRequest) XXX_Size() int {
	return xxx_messageInfo_IssueRequest.Size(m)
}
func (m *IssueRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IssueRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IssueRequest proto.InternalMessageInfo

func (m *IssueRequest) GetCredential() []byte {
	if m != nil {
		return m.Credential
	}
	return nil
}

func (m *IssueRequest) GetTokensToIssue() []*Token {
	if m != nil {
		return m.TokensToIssue
	}
	return nil
}


type RecipientShare struct {
	
	Recipient *TokenOwner `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	
	
	
	
	Quantity             string   `protobuf:"bytes,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RecipientShare) Reset()         { *m = RecipientShare{} }
func (m *RecipientShare) String() string { return proto.CompactTextString(m) }
func (*RecipientShare) ProtoMessage()    {}
func (*RecipientShare) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{1}
}
func (m *RecipientShare) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecipientShare.Unmarshal(m, b)
}
func (m *RecipientShare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecipientShare.Marshal(b, m, deterministic)
}
func (dst *RecipientShare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecipientShare.Merge(dst, src)
}
func (m *RecipientShare) XXX_Size() int {
	return xxx_messageInfo_RecipientShare.Size(m)
}
func (m *RecipientShare) XXX_DiscardUnknown() {
	xxx_messageInfo_RecipientShare.DiscardUnknown(m)
}

var xxx_messageInfo_RecipientShare proto.InternalMessageInfo

func (m *RecipientShare) GetRecipient() *TokenOwner {
	if m != nil {
		return m.Recipient
	}
	return nil
}

func (m *RecipientShare) GetQuantity() string {
	if m != nil {
		return m.Quantity
	}
	return ""
}


type TokenTransactions struct {
	Txs                  []*TokenTransaction `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TokenTransactions) Reset()         { *m = TokenTransactions{} }
func (m *TokenTransactions) String() string { return proto.CompactTextString(m) }
func (*TokenTransactions) ProtoMessage()    {}
func (*TokenTransactions) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{2}
}
func (m *TokenTransactions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenTransactions.Unmarshal(m, b)
}
func (m *TokenTransactions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenTransactions.Marshal(b, m, deterministic)
}
func (dst *TokenTransactions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenTransactions.Merge(dst, src)
}
func (m *TokenTransactions) XXX_Size() int {
	return xxx_messageInfo_TokenTransactions.Size(m)
}
func (m *TokenTransactions) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenTransactions.DiscardUnknown(m)
}

var xxx_messageInfo_TokenTransactions proto.InternalMessageInfo

func (m *TokenTransactions) GetTxs() []*TokenTransaction {
	if m != nil {
		return m.Txs
	}
	return nil
}


type TransferRequest struct {
	
	
	Credential []byte `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	
	TokenIds []*TokenId `protobuf:"bytes,2,rep,name=token_ids,json=tokenIds,proto3" json:"token_ids,omitempty"`
	
	Shares               []*RecipientShare `protobuf:"bytes,3,rep,name=shares,proto3" json:"shares,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *TransferRequest) Reset()         { *m = TransferRequest{} }
func (m *TransferRequest) String() string { return proto.CompactTextString(m) }
func (*TransferRequest) ProtoMessage()    {}
func (*TransferRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{3}
}
func (m *TransferRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransferRequest.Unmarshal(m, b)
}
func (m *TransferRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransferRequest.Marshal(b, m, deterministic)
}
func (dst *TransferRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferRequest.Merge(dst, src)
}
func (m *TransferRequest) XXX_Size() int {
	return xxx_messageInfo_TransferRequest.Size(m)
}
func (m *TransferRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TransferRequest proto.InternalMessageInfo

func (m *TransferRequest) GetCredential() []byte {
	if m != nil {
		return m.Credential
	}
	return nil
}

func (m *TransferRequest) GetTokenIds() []*TokenId {
	if m != nil {
		return m.TokenIds
	}
	return nil
}

func (m *TransferRequest) GetShares() []*RecipientShare {
	if m != nil {
		return m.Shares
	}
	return nil
}


type RedeemRequest struct {
	
	
	Credential []byte `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	
	TokenIds []*TokenId `protobuf:"bytes,2,rep,name=token_ids,json=tokenIds,proto3" json:"token_ids,omitempty"`
	
	
	
	
	Quantity             string   `protobuf:"bytes,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RedeemRequest) Reset()         { *m = RedeemRequest{} }
func (m *RedeemRequest) String() string { return proto.CompactTextString(m) }
func (*RedeemRequest) ProtoMessage()    {}
func (*RedeemRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{4}
}
func (m *RedeemRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RedeemRequest.Unmarshal(m, b)
}
func (m *RedeemRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RedeemRequest.Marshal(b, m, deterministic)
}
func (dst *RedeemRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RedeemRequest.Merge(dst, src)
}
func (m *RedeemRequest) XXX_Size() int {
	return xxx_messageInfo_RedeemRequest.Size(m)
}
func (m *RedeemRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RedeemRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RedeemRequest proto.InternalMessageInfo

func (m *RedeemRequest) GetCredential() []byte {
	if m != nil {
		return m.Credential
	}
	return nil
}

func (m *RedeemRequest) GetTokenIds() []*TokenId {
	if m != nil {
		return m.TokenIds
	}
	return nil
}

func (m *RedeemRequest) GetQuantity() string {
	if m != nil {
		return m.Quantity
	}
	return ""
}


type UnspentToken struct {
	
	Id *TokenId `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	
	
	Quantity             string   `protobuf:"bytes,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnspentToken) Reset()         { *m = UnspentToken{} }
func (m *UnspentToken) String() string { return proto.CompactTextString(m) }
func (*UnspentToken) ProtoMessage()    {}
func (*UnspentToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{5}
}
func (m *UnspentToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnspentToken.Unmarshal(m, b)
}
func (m *UnspentToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnspentToken.Marshal(b, m, deterministic)
}
func (dst *UnspentToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnspentToken.Merge(dst, src)
}
func (m *UnspentToken) XXX_Size() int {
	return xxx_messageInfo_UnspentToken.Size(m)
}
func (m *UnspentToken) XXX_DiscardUnknown() {
	xxx_messageInfo_UnspentToken.DiscardUnknown(m)
}

var xxx_messageInfo_UnspentToken proto.InternalMessageInfo

func (m *UnspentToken) GetId() *TokenId {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *UnspentToken) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *UnspentToken) GetQuantity() string {
	if m != nil {
		return m.Quantity
	}
	return ""
}


type UnspentTokens struct {
	
	Tokens               []*UnspentToken `protobuf:"bytes,1,rep,name=tokens,proto3" json:"tokens,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UnspentTokens) Reset()         { *m = UnspentTokens{} }
func (m *UnspentTokens) String() string { return proto.CompactTextString(m) }
func (*UnspentTokens) ProtoMessage()    {}
func (*UnspentTokens) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{6}
}
func (m *UnspentTokens) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnspentTokens.Unmarshal(m, b)
}
func (m *UnspentTokens) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnspentTokens.Marshal(b, m, deterministic)
}
func (dst *UnspentTokens) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnspentTokens.Merge(dst, src)
}
func (m *UnspentTokens) XXX_Size() int {
	return xxx_messageInfo_UnspentTokens.Size(m)
}
func (m *UnspentTokens) XXX_DiscardUnknown() {
	xxx_messageInfo_UnspentTokens.DiscardUnknown(m)
}

var xxx_messageInfo_UnspentTokens proto.InternalMessageInfo

func (m *UnspentTokens) GetTokens() []*UnspentToken {
	if m != nil {
		return m.Tokens
	}
	return nil
}


type ListRequest struct {
	
	Credential           []byte   `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRequest) Reset()         { *m = ListRequest{} }
func (m *ListRequest) String() string { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()    {}
func (*ListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{7}
}
func (m *ListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRequest.Unmarshal(m, b)
}
func (m *ListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRequest.Marshal(b, m, deterministic)
}
func (dst *ListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRequest.Merge(dst, src)
}
func (m *ListRequest) XXX_Size() int {
	return xxx_messageInfo_ListRequest.Size(m)
}
func (m *ListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRequest proto.InternalMessageInfo

func (m *ListRequest) GetCredential() []byte {
	if m != nil {
		return m.Credential
	}
	return nil
}




type TokenOperationRequest struct {
	
	
	Credential []byte `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	
	Operations []*TokenOperation `protobuf:"bytes,2,rep,name=operations,proto3" json:"operations,omitempty"`
	
	
	TokenIds             []*TokenId `protobuf:"bytes,3,rep,name=token_ids,json=tokenIds,proto3" json:"token_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *TokenOperationRequest) Reset()         { *m = TokenOperationRequest{} }
func (m *TokenOperationRequest) String() string { return proto.CompactTextString(m) }
func (*TokenOperationRequest) ProtoMessage()    {}
func (*TokenOperationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{8}
}
func (m *TokenOperationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenOperationRequest.Unmarshal(m, b)
}
func (m *TokenOperationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenOperationRequest.Marshal(b, m, deterministic)
}
func (dst *TokenOperationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenOperationRequest.Merge(dst, src)
}
func (m *TokenOperationRequest) XXX_Size() int {
	return xxx_messageInfo_TokenOperationRequest.Size(m)
}
func (m *TokenOperationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenOperationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TokenOperationRequest proto.InternalMessageInfo

func (m *TokenOperationRequest) GetCredential() []byte {
	if m != nil {
		return m.Credential
	}
	return nil
}

func (m *TokenOperationRequest) GetOperations() []*TokenOperation {
	if m != nil {
		return m.Operations
	}
	return nil
}

func (m *TokenOperationRequest) GetTokenIds() []*TokenId {
	if m != nil {
		return m.TokenIds
	}
	return nil
}


type Header struct {
	
	
	Timestamp *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	
	
	Nonce []byte `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	
	
	Creator []byte `protobuf:"bytes,4,opt,name=creator,proto3" json:"creator,omitempty"`
	
	
	TlsCertHash          []byte   `protobuf:"bytes,5,opt,name=tls_cert_hash,json=tlsCertHash,proto3" json:"tls_cert_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{9}
}
func (m *Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Header.Unmarshal(m, b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Header.Marshal(b, m, deterministic)
}
func (dst *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(dst, src)
}
func (m *Header) XXX_Size() int {
	return xxx_messageInfo_Header.Size(m)
}
func (m *Header) XXX_DiscardUnknown() {
	xxx_messageInfo_Header.DiscardUnknown(m)
}

var xxx_messageInfo_Header proto.InternalMessageInfo

func (m *Header) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Header) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *Header) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

func (m *Header) GetCreator() []byte {
	if m != nil {
		return m.Creator
	}
	return nil
}

func (m *Header) GetTlsCertHash() []byte {
	if m != nil {
		return m.TlsCertHash
	}
	return nil
}


type Command struct {
	
	Header *Header `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	
	
	
	
	
	
	
	
	Payload              isCommand_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{10}
}
func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (dst *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(dst, src)
}
func (m *Command) XXX_Size() int {
	return xxx_messageInfo_Command.Size(m)
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type isCommand_Payload interface {
	isCommand_Payload()
}

type Command_IssueRequest struct {
	IssueRequest *IssueRequest `protobuf:"bytes,2,opt,name=issue_request,json=issueRequest,proto3,oneof"`
}

type Command_TransferRequest struct {
	TransferRequest *TransferRequest `protobuf:"bytes,3,opt,name=transfer_request,json=transferRequest,proto3,oneof"`
}

type Command_ListRequest struct {
	ListRequest *ListRequest `protobuf:"bytes,4,opt,name=list_request,json=listRequest,proto3,oneof"`
}

type Command_RedeemRequest struct {
	RedeemRequest *RedeemRequest `protobuf:"bytes,5,opt,name=redeem_request,json=redeemRequest,proto3,oneof"`
}

type Command_TokenOperationRequest struct {
	TokenOperationRequest *TokenOperationRequest `protobuf:"bytes,6,opt,name=token_operation_request,json=tokenOperationRequest,proto3,oneof"`
}

func (*Command_IssueRequest) isCommand_Payload() {}

func (*Command_TransferRequest) isCommand_Payload() {}

func (*Command_ListRequest) isCommand_Payload() {}

func (*Command_RedeemRequest) isCommand_Payload() {}

func (*Command_TokenOperationRequest) isCommand_Payload() {}

func (m *Command) GetPayload() isCommand_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Command) GetIssueRequest() *IssueRequest {
	if x, ok := m.GetPayload().(*Command_IssueRequest); ok {
		return x.IssueRequest
	}
	return nil
}

func (m *Command) GetTransferRequest() *TransferRequest {
	if x, ok := m.GetPayload().(*Command_TransferRequest); ok {
		return x.TransferRequest
	}
	return nil
}

func (m *Command) GetListRequest() *ListRequest {
	if x, ok := m.GetPayload().(*Command_ListRequest); ok {
		return x.ListRequest
	}
	return nil
}

func (m *Command) GetRedeemRequest() *RedeemRequest {
	if x, ok := m.GetPayload().(*Command_RedeemRequest); ok {
		return x.RedeemRequest
	}
	return nil
}

func (m *Command) GetTokenOperationRequest() *TokenOperationRequest {
	if x, ok := m.GetPayload().(*Command_TokenOperationRequest); ok {
		return x.TokenOperationRequest
	}
	return nil
}


func (*Command) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Command_OneofMarshaler, _Command_OneofUnmarshaler, _Command_OneofSizer, []interface{}{
		(*Command_IssueRequest)(nil),
		(*Command_TransferRequest)(nil),
		(*Command_ListRequest)(nil),
		(*Command_RedeemRequest)(nil),
		(*Command_TokenOperationRequest)(nil),
	}
}

func _Command_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Command)
	
	switch x := m.Payload.(type) {
	case *Command_IssueRequest:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.IssueRequest); err != nil {
			return err
		}
	case *Command_TransferRequest:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TransferRequest); err != nil {
			return err
		}
	case *Command_ListRequest:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ListRequest); err != nil {
			return err
		}
	case *Command_RedeemRequest:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.RedeemRequest); err != nil {
			return err
		}
	case *Command_TokenOperationRequest:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TokenOperationRequest); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Command.Payload has unexpected type %T", x)
	}
	return nil
}

func _Command_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Command)
	switch tag {
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(IssueRequest)
		err := b.DecodeMessage(msg)
		m.Payload = &Command_IssueRequest{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TransferRequest)
		err := b.DecodeMessage(msg)
		m.Payload = &Command_TransferRequest{msg}
		return true, err
	case 4: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ListRequest)
		err := b.DecodeMessage(msg)
		m.Payload = &Command_ListRequest{msg}
		return true, err
	case 5: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RedeemRequest)
		err := b.DecodeMessage(msg)
		m.Payload = &Command_RedeemRequest{msg}
		return true, err
	case 6: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TokenOperationRequest)
		err := b.DecodeMessage(msg)
		m.Payload = &Command_TokenOperationRequest{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Command_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Command)
	
	switch x := m.Payload.(type) {
	case *Command_IssueRequest:
		s := proto.Size(x.IssueRequest)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Command_TransferRequest:
		s := proto.Size(x.TransferRequest)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Command_ListRequest:
		s := proto.Size(x.ListRequest)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Command_RedeemRequest:
		s := proto.Size(x.RedeemRequest)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Command_TokenOperationRequest:
		s := proto.Size(x.TokenOperationRequest)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type SignedCommand struct {
	
	Command []byte `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedCommand) Reset()         { *m = SignedCommand{} }
func (m *SignedCommand) String() string { return proto.CompactTextString(m) }
func (*SignedCommand) ProtoMessage()    {}
func (*SignedCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{11}
}
func (m *SignedCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedCommand.Unmarshal(m, b)
}
func (m *SignedCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedCommand.Marshal(b, m, deterministic)
}
func (dst *SignedCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedCommand.Merge(dst, src)
}
func (m *SignedCommand) XXX_Size() int {
	return xxx_messageInfo_SignedCommand.Size(m)
}
func (m *SignedCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedCommand.DiscardUnknown(m)
}

var xxx_messageInfo_SignedCommand proto.InternalMessageInfo

func (m *SignedCommand) GetCommand() []byte {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *SignedCommand) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type CommandResponseHeader struct {
	
	
	Timestamp *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	
	
	
	
	CommandHash []byte `protobuf:"bytes,2,opt,name=command_hash,json=commandHash,proto3" json:"command_hash,omitempty"`
	
	Creator              []byte   `protobuf:"bytes,3,opt,name=creator,proto3" json:"creator,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandResponseHeader) Reset()         { *m = CommandResponseHeader{} }
func (m *CommandResponseHeader) String() string { return proto.CompactTextString(m) }
func (*CommandResponseHeader) ProtoMessage()    {}
func (*CommandResponseHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{12}
}
func (m *CommandResponseHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandResponseHeader.Unmarshal(m, b)
}
func (m *CommandResponseHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandResponseHeader.Marshal(b, m, deterministic)
}
func (dst *CommandResponseHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandResponseHeader.Merge(dst, src)
}
func (m *CommandResponseHeader) XXX_Size() int {
	return xxx_messageInfo_CommandResponseHeader.Size(m)
}
func (m *CommandResponseHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandResponseHeader.DiscardUnknown(m)
}

var xxx_messageInfo_CommandResponseHeader proto.InternalMessageInfo

func (m *CommandResponseHeader) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *CommandResponseHeader) GetCommandHash() []byte {
	if m != nil {
		return m.CommandHash
	}
	return nil
}

func (m *CommandResponseHeader) GetCreator() []byte {
	if m != nil {
		return m.Creator
	}
	return nil
}


type Error struct {
	
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	
	Payload              []byte   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{13}
}
func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (dst *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(dst, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Error) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}


type CommandResponse struct {
	
	Header *CommandResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	
	
	
	
	
	
	
	Payload              isCommandResponse_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *CommandResponse) Reset()         { *m = CommandResponse{} }
func (m *CommandResponse) String() string { return proto.CompactTextString(m) }
func (*CommandResponse) ProtoMessage()    {}
func (*CommandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{14}
}
func (m *CommandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandResponse.Unmarshal(m, b)
}
func (m *CommandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandResponse.Marshal(b, m, deterministic)
}
func (dst *CommandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandResponse.Merge(dst, src)
}
func (m *CommandResponse) XXX_Size() int {
	return xxx_messageInfo_CommandResponse.Size(m)
}
func (m *CommandResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CommandResponse proto.InternalMessageInfo

func (m *CommandResponse) GetHeader() *CommandResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type isCommandResponse_Payload interface {
	isCommandResponse_Payload()
}

type CommandResponse_Err struct {
	Err *Error `protobuf:"bytes,2,opt,name=err,proto3,oneof"`
}

type CommandResponse_TokenTransaction struct {
	TokenTransaction *TokenTransaction `protobuf:"bytes,3,opt,name=token_transaction,json=tokenTransaction,proto3,oneof"`
}

type CommandResponse_UnspentTokens struct {
	UnspentTokens *UnspentTokens `protobuf:"bytes,4,opt,name=unspent_tokens,json=unspentTokens,proto3,oneof"`
}

type CommandResponse_TokenTransactions struct {
	TokenTransactions *TokenTransactions `protobuf:"bytes,5,opt,name=token_transactions,json=tokenTransactions,proto3,oneof"`
}

func (*CommandResponse_Err) isCommandResponse_Payload() {}

func (*CommandResponse_TokenTransaction) isCommandResponse_Payload() {}

func (*CommandResponse_UnspentTokens) isCommandResponse_Payload() {}

func (*CommandResponse_TokenTransactions) isCommandResponse_Payload() {}

func (m *CommandResponse) GetPayload() isCommandResponse_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *CommandResponse) GetErr() *Error {
	if x, ok := m.GetPayload().(*CommandResponse_Err); ok {
		return x.Err
	}
	return nil
}

func (m *CommandResponse) GetTokenTransaction() *TokenTransaction {
	if x, ok := m.GetPayload().(*CommandResponse_TokenTransaction); ok {
		return x.TokenTransaction
	}
	return nil
}

func (m *CommandResponse) GetUnspentTokens() *UnspentTokens {
	if x, ok := m.GetPayload().(*CommandResponse_UnspentTokens); ok {
		return x.UnspentTokens
	}
	return nil
}

func (m *CommandResponse) GetTokenTransactions() *TokenTransactions {
	if x, ok := m.GetPayload().(*CommandResponse_TokenTransactions); ok {
		return x.TokenTransactions
	}
	return nil
}


func (*CommandResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CommandResponse_OneofMarshaler, _CommandResponse_OneofUnmarshaler, _CommandResponse_OneofSizer, []interface{}{
		(*CommandResponse_Err)(nil),
		(*CommandResponse_TokenTransaction)(nil),
		(*CommandResponse_UnspentTokens)(nil),
		(*CommandResponse_TokenTransactions)(nil),
	}
}

func _CommandResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*CommandResponse)
	
	switch x := m.Payload.(type) {
	case *CommandResponse_Err:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Err); err != nil {
			return err
		}
	case *CommandResponse_TokenTransaction:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TokenTransaction); err != nil {
			return err
		}
	case *CommandResponse_UnspentTokens:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.UnspentTokens); err != nil {
			return err
		}
	case *CommandResponse_TokenTransactions:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TokenTransactions); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("CommandResponse.Payload has unexpected type %T", x)
	}
	return nil
}

func _CommandResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*CommandResponse)
	switch tag {
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Error)
		err := b.DecodeMessage(msg)
		m.Payload = &CommandResponse_Err{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TokenTransaction)
		err := b.DecodeMessage(msg)
		m.Payload = &CommandResponse_TokenTransaction{msg}
		return true, err
	case 4: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(UnspentTokens)
		err := b.DecodeMessage(msg)
		m.Payload = &CommandResponse_UnspentTokens{msg}
		return true, err
	case 5: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TokenTransactions)
		err := b.DecodeMessage(msg)
		m.Payload = &CommandResponse_TokenTransactions{msg}
		return true, err
	default:
		return false, nil
	}
}

func _CommandResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*CommandResponse)
	
	switch x := m.Payload.(type) {
	case *CommandResponse_Err:
		s := proto.Size(x.Err)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CommandResponse_TokenTransaction:
		s := proto.Size(x.TokenTransaction)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CommandResponse_UnspentTokens:
		s := proto.Size(x.UnspentTokens)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CommandResponse_TokenTransactions:
		s := proto.Size(x.TokenTransactions)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type SignedCommandResponse struct {
	
	Response []byte `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
	
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedCommandResponse) Reset()         { *m = SignedCommandResponse{} }
func (m *SignedCommandResponse) String() string { return proto.CompactTextString(m) }
func (*SignedCommandResponse) ProtoMessage()    {}
func (*SignedCommandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_2a2fc6c2f642aa49, []int{15}
}
func (m *SignedCommandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedCommandResponse.Unmarshal(m, b)
}
func (m *SignedCommandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedCommandResponse.Marshal(b, m, deterministic)
}
func (dst *SignedCommandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedCommandResponse.Merge(dst, src)
}
func (m *SignedCommandResponse) XXX_Size() int {
	return xxx_messageInfo_SignedCommandResponse.Size(m)
}
func (m *SignedCommandResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedCommandResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SignedCommandResponse proto.InternalMessageInfo

func (m *SignedCommandResponse) GetResponse() []byte {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *SignedCommandResponse) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*IssueRequest)(nil), "token.IssueRequest")
	proto.RegisterType((*RecipientShare)(nil), "token.RecipientShare")
	proto.RegisterType((*TokenTransactions)(nil), "token.TokenTransactions")
	proto.RegisterType((*TransferRequest)(nil), "token.TransferRequest")
	proto.RegisterType((*RedeemRequest)(nil), "token.RedeemRequest")
	proto.RegisterType((*UnspentToken)(nil), "token.UnspentToken")
	proto.RegisterType((*UnspentTokens)(nil), "token.UnspentTokens")
	proto.RegisterType((*ListRequest)(nil), "token.ListRequest")
	proto.RegisterType((*TokenOperationRequest)(nil), "token.TokenOperationRequest")
	proto.RegisterType((*Header)(nil), "token.Header")
	proto.RegisterType((*Command)(nil), "token.Command")
	proto.RegisterType((*SignedCommand)(nil), "token.SignedCommand")
	proto.RegisterType((*CommandResponseHeader)(nil), "token.CommandResponseHeader")
	proto.RegisterType((*Error)(nil), "token.Error")
	proto.RegisterType((*CommandResponse)(nil), "token.CommandResponse")
	proto.RegisterType((*SignedCommandResponse)(nil), "token.SignedCommandResponse")
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type ProverClient interface {
	
	
	
	
	ProcessCommand(ctx context.Context, in *SignedCommand, opts ...grpc.CallOption) (*SignedCommandResponse, error)
}

type proverClient struct {
	cc *grpc.ClientConn
}

func NewProverClient(cc *grpc.ClientConn) ProverClient {
	return &proverClient{cc}
}

func (c *proverClient) ProcessCommand(ctx context.Context, in *SignedCommand, opts ...grpc.CallOption) (*SignedCommandResponse, error) {
	out := new(SignedCommandResponse)
	err := c.cc.Invoke(ctx, "/token.Prover/ProcessCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}


type ProverServer interface {
	
	
	
	
	ProcessCommand(context.Context, *SignedCommand) (*SignedCommandResponse, error)
}

func RegisterProverServer(s *grpc.Server, srv ProverServer) {
	s.RegisterService(&_Prover_serviceDesc, srv)
}

func _Prover_ProcessCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignedCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProverServer).ProcessCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.Prover/ProcessCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProverServer).ProcessCommand(ctx, req.(*SignedCommand))
	}
	return interceptor(ctx, in, info, handler)
}

var _Prover_serviceDesc = grpc.ServiceDesc{
	ServiceName: "token.Prover",
	HandlerType: (*ProverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessCommand",
			Handler:    _Prover_ProcessCommand_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token/prover.proto",
}

func init() { proto.RegisterFile("token/prover.proto", fileDescriptor_prover_2a2fc6c2f642aa49) }

var fileDescriptor_prover_2a2fc6c2f642aa49 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0xef, 0x6e, 0x23, 0x35,
	0x10, 0x4f, 0x9a, 0x26, 0x6d, 0x26, 0x49, 0x7b, 0x35, 0x97, 0xbb, 0x55, 0x54, 0x8e, 0xb2, 0x08,
	0xa9, 0x70, 0xba, 0x8d, 0x54, 0x0e, 0x81, 0xf8, 0xf7, 0xe1, 0x2a, 0x8e, 0x54, 0x42, 0xa2, 0xe7,
	0x2b, 0x7c, 0x40, 0x82, 0x95, 0xbb, 0x3b, 0x4d, 0x2c, 0x36, 0xeb, 0x3d, 0xdb, 0x81, 0xeb, 0x3b,
	0x00, 0x2f, 0xc0, 0x73, 0xf0, 0x0a, 0x3c, 0x17, 0x5a, 0xdb, 0xbb, 0xf1, 0x46, 0xbd, 0xaa, 0x12,
	0xe2, 0x4b, 0x64, 0x8f, 0x7f, 0xf3, 0xf3, 0x78, 0xe6, 0x37, 0x93, 0x05, 0xa2, 0xc5, 0x2f, 0x98,
	0x4f, 0x0b, 0x29, 0x7e, 0x45, 0x19, 0x15, 0x52, 0x68, 0x41, 0xba, 0xc6, 0x36, 0x79, 0x67, 0x2e,
	0xc4, 0x3c, 0xc3, 0xa9, 0x31, 0x5e, 0xae, 0xae, 0xa6, 0x9a, 0x2f, 0x51, 0x69, 0xb6, 0x2c, 0x2c,
	0x6e, 0xf2, 0xc0, 0xfa, 0x8a, 0x02, 0x25, 0xd3, 0x5c, 0xe4, 0xca, 0xd9, 0x1f, 0x5a, 0xbb, 0x96,
	0x2c, 0x57, 0x2c, 0x29, 0x4f, 0xec, 0x41, 0x98, 0xc2, 0xf0, 0x4c, 0xa9, 0x15, 0x52, 0x7c, 0xb5,
	0x42, 0xa5, 0xc9, 0x23, 0x80, 0x44, 0x62, 0x8a, 0xb9, 0xe6, 0x2c, 0x0b, 0xda, 0x47, 0xed, 0xe3,
	0x21, 0xf5, 0x2c, 0xe4, 0x29, 0xec, 0x1b, 0x2a, 0x15, 0x6b, 0x11, 0xf3, 0xd2, 0x33, 0xd8, 0x3a,
	0xea, 0x1c, 0x0f, 0x4e, 0x86, 0x91, 0xb1, 0x47, 0x17, 0xe5, 0x2f, 0x1d, 0x59, 0xd0, 0x85, 0x30,
	0xe4, 0xe1, 0x4f, 0xb0, 0x47, 0x31, 0xe1, 0x05, 0xc7, 0x5c, 0xbf, 0x5c, 0x30, 0x89, 0x64, 0x0a,
	0x7d, 0x59, 0x59, 0xcc, 0x35, 0x83, 0x93, 0x03, 0x9f, 0xe1, 0xbb, 0xdf, 0x72, 0x94, 0x74, 0x8d,
	0x21, 0x13, 0xd8, 0x7d, 0xb5, 0x62, 0xb9, 0xe6, 0xfa, 0x3a, 0xd8, 0x3a, 0x6a, 0x1f, 0xf7, 0x69,
	0xbd, 0x0f, 0xbf, 0x82, 0x03, 0xe3, 0x74, 0xb1, 0x7e, 0x9e, 0x22, 0x1f, 0x40, 0x47, 0xbf, 0x56,
	0x41, 0xdb, 0x44, 0xf7, 0xd0, 0xe7, 0xf6, 0x60, 0xb4, 0xc4, 0x84, 0x7f, 0xb4, 0x61, 0xdf, 0x18,
	0xaf, 0x50, 0xde, 0x35, 0x11, 0x8f, 0xa1, 0x6f, 0x28, 0x63, 0x9e, 0x2a, 0x97, 0x82, 0x3d, 0xff,
	0x92, 0xb3, 0x94, 0xee, 0x6a, 0xbb, 0x50, 0xe4, 0x09, 0xf4, 0x54, 0xf9, 0x6c, 0x15, 0x74, 0x0c,
	0x72, 0xec, 0x90, 0xcd, 0xa4, 0x50, 0x07, 0x0a, 0x5f, 0xc3, 0x88, 0x62, 0x8a, 0xb8, 0xfc, 0x5f,
	0x82, 0xf1, 0x33, 0xd9, 0xd9, 0xc8, 0xe4, 0xcf, 0x30, 0xfc, 0x3e, 0x57, 0x05, 0xe6, 0xda, 0xf8,
	0x91, 0x47, 0xb0, 0xc5, 0x53, 0x57, 0x9f, 0x4d, 0xc6, 0x2d, 0x9e, 0x12, 0x02, 0xdb, 0xfa, 0xba,
	0x40, 0x57, 0x11, 0xb3, 0xbe, 0x95, 0xff, 0x0b, 0x18, 0xf9, 0xfc, 0x8a, 0x3c, 0x86, 0x9e, 0x95,
	0x8a, 0x2b, 0xd4, 0x5b, 0xee, 0x12, 0x1f, 0x45, 0x1d, 0x24, 0x7c, 0x02, 0x83, 0x6f, 0xb9, 0xd2,
	0x77, 0xcc, 0x4a, 0xf8, 0x57, 0x1b, 0xc6, 0x56, 0x4c, 0x55, 0x3b, 0xdc, 0x35, 0x9f, 0x1f, 0x03,
	0xac, 0x5b, 0xc8, 0x25, 0x74, 0xdc, 0x90, 0x67, 0xcd, 0xe8, 0x01, 0x9b, 0x65, 0xe8, 0xdc, 0x5e,
	0x86, 0xf0, 0xef, 0x36, 0xf4, 0x66, 0xc8, 0x52, 0x94, 0xe4, 0x53, 0xe8, 0xd7, 0x8d, 0xec, 0x92,
	0x3d, 0x89, 0x6c, 0xab, 0x47, 0x55, 0xab, 0x47, 0x17, 0x15, 0x82, 0xae, 0xc1, 0xe4, 0x6d, 0x80,
	0x64, 0xc1, 0xf2, 0x1c, 0xb3, 0x98, 0xa7, 0xae, 0x0a, 0x7d, 0x67, 0x39, 0x4b, 0xc9, 0x7d, 0xe8,
	0xe6, 0x22, 0x4f, 0xd0, 0xd4, 0x61, 0x48, 0xed, 0x86, 0x04, 0xb0, 0x93, 0x48, 0x64, 0x5a, 0xc8,
	0x60, 0xdb, 0xd8, 0xab, 0x2d, 0x09, 0x61, 0xa4, 0x33, 0x15, 0x27, 0x28, 0x75, 0xbc, 0x60, 0x6a,
	0x11, 0x74, 0xcd, 0xf9, 0x40, 0x67, 0xea, 0x14, 0xa5, 0x9e, 0x31, 0xb5, 0x08, 0x7f, 0xef, 0xc0,
	0xce, 0xa9, 0x58, 0x2e, 0x59, 0x9e, 0x92, 0xf7, 0xa1, 0xb7, 0x30, 0x4f, 0x70, 0x51, 0x8f, 0xdc,
	0x6b, 0xed, 0xbb, 0xa8, 0x3b, 0x24, 0x9f, 0xc1, 0xc8, 0x8c, 0x8a, 0x58, 0xda, 0xfc, 0x9b, 0x40,
	0xd7, 0xb5, 0xf6, 0x07, 0xd0, 0xac, 0x45, 0x87, 0xdc, 0x1f, 0x48, 0xa7, 0x70, 0x4f, 0xbb, 0xd6,
	0xac, 0xdd, 0x3b, 0xc6, 0xfd, 0x41, 0x95, 0xda, 0x66, 0xe7, 0xce, 0x5a, 0x74, 0x5f, 0x6f, 0x34,
	0xf3, 0x27, 0x30, 0xcc, 0xb8, 0xd2, 0x35, 0xc1, 0xb6, 0x21, 0x20, 0x8e, 0xc0, 0xd3, 0xd4, 0xac,
	0x45, 0x07, 0x99, 0x27, 0xb1, 0x2f, 0x61, 0x4f, 0x9a, 0x4e, 0xac, 0x5d, 0xbb, 0xc6, 0xf5, 0x7e,
	0xdd, 0xc0, 0x5e, 0x9b, 0xce, 0x5a, 0x74, 0x24, 0x1b, 0x7d, 0xfb, 0x03, 0xd8, 0xc1, 0x1b, 0xd7,
	0x22, 0xa9, 0x79, 0x7a, 0x86, 0xe7, 0xf0, 0x66, 0x51, 0xd5, 0x7c, 0x63, 0x7d, 0xd3, 0xc1, 0xb3,
	0x3e, 0xec, 0x14, 0xec, 0x3a, 0x13, 0x2c, 0x0d, 0xbf, 0x81, 0xd1, 0x4b, 0x3e, 0xcf, 0x31, 0xad,
	0x6a, 0x52, 0x56, 0xd7, 0x2e, 0x9d, 0xb0, 0xab, 0x2d, 0x39, 0x84, 0xbe, 0xe2, 0xf3, 0x9c, 0xe9,
	0x95, 0xb4, 0x1d, 0x3b, 0xa4, 0x6b, 0x43, 0xf8, 0x67, 0x1b, 0xc6, 0x8e, 0x83, 0xa2, 0x2a, 0x44,
	0xae, 0xf0, 0x3f, 0xcb, 0xf3, 0x5d, 0x18, 0xba, 0xcb, 0xad, 0x9c, 0xec, 0xa5, 0x03, 0x67, 0x2b,
	0xe5, 0xe4, 0x8b, 0xb1, 0xd3, 0x10, 0x63, 0xf8, 0x39, 0x74, 0xbf, 0x96, 0x52, 0xc8, 0x12, 0xb2,
	0x44, 0xa5, 0xd8, 0x1c, 0xcd, 0xed, 0x7d, 0x5a, 0x6d, 0xcb, 0x13, 0x97, 0x07, 0x47, 0x5d, 0xa7,
	0xe5, 0x9f, 0x2d, 0xd8, 0xdf, 0x78, 0x0d, 0x79, 0xba, 0xa1, 0xd6, 0x2a, 0xf9, 0x37, 0xbe, 0xba,
	0x16, 0xef, 0x11, 0x74, 0x50, 0x4a, 0x27, 0xd9, 0xea, 0x5f, 0xce, 0x04, 0x36, 0x6b, 0xd1, 0xf2,
	0x88, 0x3c, 0x87, 0x03, 0x5b, 0x65, 0xef, 0xef, 0xd5, 0x69, 0xf4, 0x4d, 0xff, 0x3b, 0xb3, 0x16,
	0xbd, 0xa7, 0x37, 0x6c, 0xa5, 0xd8, 0x56, 0x76, 0xec, 0xc5, 0x6e, 0x26, 0x6e, 0x37, 0xc4, 0xd6,
	0x98, 0x9c, 0xa5, 0xd8, 0x56, 0x8d, 0x51, 0x7a, 0xe6, 0xbe, 0x1c, 0xfc, 0x30, 0x94, 0xd3, 0x6b,
	0xf0, 0x86, 0x38, 0x4a, 0x9a, 0x83, 0xcd, 0x40, 0x94, 0xaf, 0xaf, 0x17, 0x30, 0x6e, 0xe8, 0xab,
	0xce, 0xe6, 0x04, 0x76, 0xa5, 0x5b, 0x3b, 0xa1, 0xd5, 0xfb, 0xdb, 0x95, 0x76, 0x72, 0x0e, 0xbd,
	0x73, 0xf3, 0x71, 0x43, 0x9e, 0xc3, 0xde, 0xb9, 0x14, 0x09, 0x2a, 0x55, 0xa9, 0xb7, 0x7a, 0x6b,
	0xe3, 0xce, 0xc9, 0xe1, 0x4d, 0xd6, 0x2a, 0x92, 0xb0, 0xf5, 0xec, 0x05, 0xbc, 0x27, 0xe4, 0x3c,
	0x5a, 0x5c, 0x17, 0x28, 0x33, 0x4c, 0xe7, 0x28, 0xa3, 0x2b, 0x76, 0x29, 0x79, 0x62, 0xf5, 0xa9,
	0xac, 0xfb, 0x8f, 0x1f, 0xce, 0xb9, 0x5e, 0xac, 0x2e, 0xa3, 0x44, 0x2c, 0xa7, 0x1e, 0x76, 0x6a,
	0xb1, 0xf6, 0xab, 0x4a, 0x4d, 0x0d, 0xf6, 0xb2, 0x67, 0x76, 0x1f, 0xfd, 0x1b, 0x00, 0x00, 0xff,
	0xff, 0x93, 0x22, 0x54, 0x97, 0x8e, 0x09, 0x00, 0x00,
}
