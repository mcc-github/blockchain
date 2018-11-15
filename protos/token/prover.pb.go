


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


type TokenToIssue struct {
	
	Recipient []byte `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	
	Quantity             uint64   `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenToIssue) Reset()         { *m = TokenToIssue{} }
func (m *TokenToIssue) String() string { return proto.CompactTextString(m) }
func (*TokenToIssue) ProtoMessage()    {}
func (*TokenToIssue) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_8c62b4bd164d3545, []int{0}
}
func (m *TokenToIssue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenToIssue.Unmarshal(m, b)
}
func (m *TokenToIssue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenToIssue.Marshal(b, m, deterministic)
}
func (dst *TokenToIssue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenToIssue.Merge(dst, src)
}
func (m *TokenToIssue) XXX_Size() int {
	return xxx_messageInfo_TokenToIssue.Size(m)
}
func (m *TokenToIssue) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenToIssue.DiscardUnknown(m)
}

var xxx_messageInfo_TokenToIssue proto.InternalMessageInfo

func (m *TokenToIssue) GetRecipient() []byte {
	if m != nil {
		return m.Recipient
	}
	return nil
}

func (m *TokenToIssue) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *TokenToIssue) GetQuantity() uint64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}


type RecipientTransferShare struct {
	
	Recipient []byte `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	
	Quantity             uint64   `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RecipientTransferShare) Reset()         { *m = RecipientTransferShare{} }
func (m *RecipientTransferShare) String() string { return proto.CompactTextString(m) }
func (*RecipientTransferShare) ProtoMessage()    {}
func (*RecipientTransferShare) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_8c62b4bd164d3545, []int{1}
}
func (m *RecipientTransferShare) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecipientTransferShare.Unmarshal(m, b)
}
func (m *RecipientTransferShare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecipientTransferShare.Marshal(b, m, deterministic)
}
func (dst *RecipientTransferShare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecipientTransferShare.Merge(dst, src)
}
func (m *RecipientTransferShare) XXX_Size() int {
	return xxx_messageInfo_RecipientTransferShare.Size(m)
}
func (m *RecipientTransferShare) XXX_DiscardUnknown() {
	xxx_messageInfo_RecipientTransferShare.DiscardUnknown(m)
}

var xxx_messageInfo_RecipientTransferShare proto.InternalMessageInfo

func (m *RecipientTransferShare) GetRecipient() []byte {
	if m != nil {
		return m.Recipient
	}
	return nil
}

func (m *RecipientTransferShare) GetQuantity() uint64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}


type TokenOutput struct {
	
	Id []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	
	Quantity             uint64   `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenOutput) Reset()         { *m = TokenOutput{} }
func (m *TokenOutput) String() string { return proto.CompactTextString(m) }
func (*TokenOutput) ProtoMessage()    {}
func (*TokenOutput) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_8c62b4bd164d3545, []int{2}
}
func (m *TokenOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenOutput.Unmarshal(m, b)
}
func (m *TokenOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenOutput.Marshal(b, m, deterministic)
}
func (dst *TokenOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenOutput.Merge(dst, src)
}
func (m *TokenOutput) XXX_Size() int {
	return xxx_messageInfo_TokenOutput.Size(m)
}
func (m *TokenOutput) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenOutput.DiscardUnknown(m)
}

var xxx_messageInfo_TokenOutput proto.InternalMessageInfo

func (m *TokenOutput) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *TokenOutput) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *TokenOutput) GetQuantity() uint64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}


type UnspentTokens struct {
	Tokens               []*TokenOutput `protobuf:"bytes,1,rep,name=tokens,proto3" json:"tokens,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *UnspentTokens) Reset()         { *m = UnspentTokens{} }
func (m *UnspentTokens) String() string { return proto.CompactTextString(m) }
func (*UnspentTokens) ProtoMessage()    {}
func (*UnspentTokens) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_8c62b4bd164d3545, []int{3}
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

func (m *UnspentTokens) GetTokens() []*TokenOutput {
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
	return fileDescriptor_prover_8c62b4bd164d3545, []int{4}
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


type ImportRequest struct {
	
	
	Credential []byte `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	
	TokensToIssue        []*TokenToIssue `protobuf:"bytes,2,rep,name=tokens_to_issue,json=tokensToIssue,proto3" json:"tokens_to_issue,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ImportRequest) Reset()         { *m = ImportRequest{} }
func (m *ImportRequest) String() string { return proto.CompactTextString(m) }
func (*ImportRequest) ProtoMessage()    {}
func (*ImportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_8c62b4bd164d3545, []int{5}
}
func (m *ImportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportRequest.Unmarshal(m, b)
}
func (m *ImportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportRequest.Marshal(b, m, deterministic)
}
func (dst *ImportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportRequest.Merge(dst, src)
}
func (m *ImportRequest) XXX_Size() int {
	return xxx_messageInfo_ImportRequest.Size(m)
}
func (m *ImportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ImportRequest proto.InternalMessageInfo

func (m *ImportRequest) GetCredential() []byte {
	if m != nil {
		return m.Credential
	}
	return nil
}

func (m *ImportRequest) GetTokensToIssue() []*TokenToIssue {
	if m != nil {
		return m.TokensToIssue
	}
	return nil
}


type TransferRequest struct {
	Credential           []byte                    `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	TokenIds             [][]byte                  `protobuf:"bytes,2,rep,name=token_ids,json=tokenIds,proto3" json:"token_ids,omitempty"`
	Shares               []*RecipientTransferShare `protobuf:"bytes,3,rep,name=shares,proto3" json:"shares,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *TransferRequest) Reset()         { *m = TransferRequest{} }
func (m *TransferRequest) String() string { return proto.CompactTextString(m) }
func (*TransferRequest) ProtoMessage()    {}
func (*TransferRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_8c62b4bd164d3545, []int{6}
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

func (m *TransferRequest) GetTokenIds() [][]byte {
	if m != nil {
		return m.TokenIds
	}
	return nil
}

func (m *TransferRequest) GetShares() []*RecipientTransferShare {
	if m != nil {
		return m.Shares
	}
	return nil
}


type RedeemRequest struct {
	
	
	Credential []byte `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	
	TokenIds [][]byte `protobuf:"bytes,2,rep,name=token_ids,json=tokenIds,proto3" json:"token_ids,omitempty"`
	
	QuantityToRedeem     uint64   `protobuf:"varint,3,opt,name=quantity_to_redeem,json=quantityToRedeem,proto3" json:"quantity_to_redeem,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RedeemRequest) Reset()         { *m = RedeemRequest{} }
func (m *RedeemRequest) String() string { return proto.CompactTextString(m) }
func (*RedeemRequest) ProtoMessage()    {}
func (*RedeemRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_8c62b4bd164d3545, []int{7}
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

func (m *RedeemRequest) GetTokenIds() [][]byte {
	if m != nil {
		return m.TokenIds
	}
	return nil
}

func (m *RedeemRequest) GetQuantityToRedeem() uint64 {
	if m != nil {
		return m.QuantityToRedeem
	}
	return 0
}


type Header struct {
	
	
	Timestamp *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	
	
	Nonce []byte `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	
	
	Creator              []byte   `protobuf:"bytes,4,opt,name=creator,proto3" json:"creator,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_8c62b4bd164d3545, []int{8}
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
	return fileDescriptor_prover_8c62b4bd164d3545, []int{9}
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

type Command_ImportRequest struct {
	ImportRequest *ImportRequest `protobuf:"bytes,2,opt,name=import_request,json=importRequest,proto3,oneof"`
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

func (*Command_ImportRequest) isCommand_Payload() {}

func (*Command_TransferRequest) isCommand_Payload() {}

func (*Command_ListRequest) isCommand_Payload() {}

func (*Command_RedeemRequest) isCommand_Payload() {}

func (m *Command) GetPayload() isCommand_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Command) GetImportRequest() *ImportRequest {
	if x, ok := m.GetPayload().(*Command_ImportRequest); ok {
		return x.ImportRequest
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


func (*Command) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Command_OneofMarshaler, _Command_OneofUnmarshaler, _Command_OneofSizer, []interface{}{
		(*Command_ImportRequest)(nil),
		(*Command_TransferRequest)(nil),
		(*Command_ListRequest)(nil),
		(*Command_RedeemRequest)(nil),
	}
}

func _Command_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Command)
	
	switch x := m.Payload.(type) {
	case *Command_ImportRequest:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ImportRequest); err != nil {
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
		msg := new(ImportRequest)
		err := b.DecodeMessage(msg)
		m.Payload = &Command_ImportRequest{msg}
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
	default:
		return false, nil
	}
}

func _Command_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Command)
	
	switch x := m.Payload.(type) {
	case *Command_ImportRequest:
		s := proto.Size(x.ImportRequest)
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
	return fileDescriptor_prover_8c62b4bd164d3545, []int{10}
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
	return fileDescriptor_prover_8c62b4bd164d3545, []int{11}
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
	return fileDescriptor_prover_8c62b4bd164d3545, []int{12}
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
	return fileDescriptor_prover_8c62b4bd164d3545, []int{13}
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

func (*CommandResponse_Err) isCommandResponse_Payload() {}

func (*CommandResponse_TokenTransaction) isCommandResponse_Payload() {}

func (*CommandResponse_UnspentTokens) isCommandResponse_Payload() {}

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


func (*CommandResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CommandResponse_OneofMarshaler, _CommandResponse_OneofUnmarshaler, _CommandResponse_OneofSizer, []interface{}{
		(*CommandResponse_Err)(nil),
		(*CommandResponse_TokenTransaction)(nil),
		(*CommandResponse_UnspentTokens)(nil),
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
	return fileDescriptor_prover_8c62b4bd164d3545, []int{14}
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
	proto.RegisterType((*TokenToIssue)(nil), "protos.TokenToIssue")
	proto.RegisterType((*RecipientTransferShare)(nil), "protos.RecipientTransferShare")
	proto.RegisterType((*TokenOutput)(nil), "protos.TokenOutput")
	proto.RegisterType((*UnspentTokens)(nil), "protos.UnspentTokens")
	proto.RegisterType((*ListRequest)(nil), "protos.ListRequest")
	proto.RegisterType((*ImportRequest)(nil), "protos.ImportRequest")
	proto.RegisterType((*TransferRequest)(nil), "protos.TransferRequest")
	proto.RegisterType((*RedeemRequest)(nil), "protos.RedeemRequest")
	proto.RegisterType((*Header)(nil), "protos.Header")
	proto.RegisterType((*Command)(nil), "protos.Command")
	proto.RegisterType((*SignedCommand)(nil), "protos.SignedCommand")
	proto.RegisterType((*CommandResponseHeader)(nil), "protos.CommandResponseHeader")
	proto.RegisterType((*Error)(nil), "protos.Error")
	proto.RegisterType((*CommandResponse)(nil), "protos.CommandResponse")
	proto.RegisterType((*SignedCommandResponse)(nil), "protos.SignedCommandResponse")
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
	err := c.cc.Invoke(ctx, "/protos.Prover/ProcessCommand", in, out, opts...)
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
		FullMethod: "/protos.Prover/ProcessCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProverServer).ProcessCommand(ctx, req.(*SignedCommand))
	}
	return interceptor(ctx, in, info, handler)
}

var _Prover_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Prover",
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

func init() { proto.RegisterFile("token/prover.proto", fileDescriptor_prover_8c62b4bd164d3545) }

var fileDescriptor_prover_8c62b4bd164d3545 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0x5b, 0x8f, 0xdb, 0x44,
	0x14, 0x8e, 0x93, 0x6c, 0x76, 0x73, 0x72, 0xd9, 0xed, 0xd0, 0xa5, 0x56, 0x60, 0xcb, 0xd6, 0x48,
	0x68, 0xc5, 0xc5, 0x91, 0x16, 0x81, 0x2a, 0x51, 0x55, 0xa8, 0x80, 0xf0, 0x4a, 0x20, 0xda, 0x69,
	0x78, 0x41, 0x48, 0xd1, 0xc4, 0x9e, 0xb5, 0x47, 0xc4, 0x1e, 0x77, 0x66, 0x8c, 0x14, 0x7e, 0x00,
	0x8f, 0x48, 0xfc, 0x1f, 0x7e, 0x18, 0x8f, 0xc8, 0x73, 0x71, 0xec, 0x68, 0x05, 0x45, 0xf0, 0x14,
	0x9f, 0xcb, 0x9c, 0xf3, 0xf9, 0x3b, 0xdf, 0x19, 0x07, 0x90, 0xe2, 0x3f, 0xd1, 0x62, 0x59, 0x0a,
	0xfe, 0x33, 0x15, 0x61, 0x29, 0xb8, 0xe2, 0x68, 0xa4, 0x7f, 0xe4, 0xe2, 0x9d, 0x94, 0xf3, 0x74,
	0x4b, 0x97, 0xda, 0xdc, 0x54, 0xb7, 0x4b, 0xc5, 0x72, 0x2a, 0x15, 0xc9, 0x4b, 0x93, 0xb8, 0x78,
	0x60, 0x0e, 0x2b, 0x41, 0x0a, 0x49, 0x62, 0xc5, 0x78, 0x61, 0x02, 0xc1, 0x8f, 0x30, 0x5d, 0xd5,
	0xa1, 0x15, 0xbf, 0x91, 0xb2, 0xa2, 0xe8, 0x6d, 0x18, 0x0b, 0x1a, 0xb3, 0x92, 0xd1, 0x42, 0xf9,
	0xde, 0xa5, 0x77, 0x35, 0xc5, 0x7b, 0x07, 0x42, 0x30, 0x54, 0xbb, 0x92, 0xfa, 0xfd, 0x4b, 0xef,
	0x6a, 0x8c, 0xf5, 0x33, 0x5a, 0xc0, 0xc9, 0xab, 0x8a, 0x14, 0x8a, 0xa9, 0x9d, 0x3f, 0xb8, 0xf4,
	0xae, 0x86, 0xb8, 0xb1, 0x03, 0x0c, 0x6f, 0x62, 0x77, 0x78, 0x55, 0xf7, 0xbe, 0xa5, 0xe2, 0x65,
	0x46, 0xc4, 0x3f, 0xf5, 0x69, 0xd7, 0xec, 0x1f, 0xd4, 0xfc, 0x16, 0x26, 0x1a, 0xf1, 0x77, 0x95,
	0x2a, 0x2b, 0x85, 0xe6, 0xd0, 0x67, 0x89, 0xad, 0xd0, 0x67, 0xc9, 0xbf, 0x86, 0xf8, 0x04, 0x66,
	0xdf, 0x17, 0xb2, 0xac, 0x01, 0xd6, 0x55, 0x25, 0xfa, 0x00, 0x46, 0x9a, 0x2c, 0xe9, 0x7b, 0x97,
	0x83, 0xab, 0xc9, 0xf5, 0x1b, 0x86, 0x29, 0x19, 0xb6, 0xba, 0x62, 0x9b, 0x12, 0x7c, 0x04, 0x93,
	0x6f, 0x98, 0x54, 0x98, 0xbe, 0xaa, 0xa8, 0x54, 0xe8, 0x21, 0x40, 0x2c, 0x68, 0x42, 0x0b, 0xc5,
	0xc8, 0xd6, 0x82, 0x6a, 0x79, 0x82, 0x1c, 0x66, 0x37, 0x79, 0xc9, 0xc5, 0xeb, 0x1e, 0x40, 0x4f,
	0xe0, 0xd4, 0x74, 0x5a, 0x2b, 0xbe, 0x66, 0xf5, 0x84, 0xfc, 0xbe, 0x46, 0x75, 0xbf, 0x83, 0xca,
	0x4e, 0x0f, 0xcf, 0x4c, 0xb2, 0x35, 0x83, 0x5f, 0x3d, 0x38, 0x75, 0xb4, 0xbf, 0x6e, 0xc7, 0xb7,
	0x60, 0xac, 0x8b, 0xac, 0x59, 0x22, 0x75, 0xaf, 0x29, 0x3e, 0xd1, 0x8e, 0x9b, 0x44, 0xa2, 0x4f,
	0x61, 0x24, 0xeb, 0xf1, 0x49, 0x7f, 0xa0, 0x51, 0x3c, 0x74, 0x28, 0xee, 0x9e, 0x32, 0xb6, 0xd9,
	0xc1, 0x2f, 0x30, 0xc3, 0x34, 0xa1, 0x34, 0xff, 0x5f, 0x50, 0x7c, 0x08, 0xc8, 0x8d, 0xaf, 0xa6,
	0x45, 0xe8, 0xca, 0x76, 0xb0, 0x67, 0x2e, 0xb2, 0xe2, 0xa6, 0x63, 0xf0, 0xbb, 0x07, 0xa3, 0x88,
	0x92, 0x84, 0x0a, 0xf4, 0x18, 0xc6, 0xcd, 0x62, 0xe8, 0xa6, 0x93, 0xeb, 0x45, 0x68, 0x56, 0x27,
	0x74, 0xab, 0x13, 0xae, 0x5c, 0x06, 0xde, 0x27, 0xa3, 0x0b, 0x80, 0x38, 0x23, 0x45, 0x41, 0xb7,
	0x6b, 0x96, 0x58, 0x6d, 0x8d, 0xad, 0xe7, 0x26, 0x41, 0xf7, 0xe1, 0xa8, 0xe0, 0x45, 0x4c, 0x35,
	0x88, 0x29, 0x36, 0x06, 0xf2, 0xe1, 0x38, 0x16, 0x94, 0x28, 0x2e, 0xfc, 0xa1, 0xf6, 0x3b, 0x33,
	0xf8, 0xa3, 0x0f, 0xc7, 0x5f, 0xf0, 0x3c, 0x27, 0x45, 0x82, 0xde, 0x83, 0x51, 0xa6, 0xe1, 0x59,
	0x44, 0x73, 0xc7, 0xa9, 0x01, 0x8d, 0x6d, 0x14, 0x3d, 0x85, 0x39, 0xd3, 0xda, 0x59, 0x0b, 0x43,
	0xa2, 0x86, 0x31, 0xb9, 0x3e, 0x77, 0xf9, 0x1d, 0x65, 0x45, 0x3d, 0x3c, 0x63, 0x1d, 0xa9, 0x7d,
	0x09, 0x67, 0xca, 0x0e, 0xa7, 0xa9, 0x30, 0xd0, 0x15, 0x1e, 0x34, 0x5a, 0xea, 0x6a, 0x25, 0xea,
	0xe1, 0x53, 0x75, 0x20, 0x9f, 0xc7, 0x30, 0xdd, 0x32, 0xb9, 0xc7, 0x30, 0xd4, 0x15, 0x9a, 0x1d,
	0x69, 0x2d, 0x43, 0xd4, 0xc3, 0x93, 0x6d, 0x6b, 0x37, 0x9e, 0xc2, 0xdc, 0x4c, 0xaa, 0x39, 0x7b,
	0xd4, 0xc5, 0xdf, 0x51, 0x48, 0x8d, 0x5f, 0xb4, 0x1d, 0xcf, 0xc6, 0x70, 0x5c, 0x92, 0xdd, 0x96,
	0x93, 0x24, 0xf8, 0x1a, 0x66, 0x2f, 0x59, 0x5a, 0xd0, 0xc4, 0x71, 0x58, 0x33, 0x6d, 0x1e, 0xad,
	0x96, 0x9c, 0x59, 0xdf, 0x33, 0x92, 0xa5, 0x05, 0x51, 0x95, 0x30, 0x77, 0xc2, 0x14, 0xef, 0x1d,
	0xc1, 0x6f, 0x1e, 0x9c, 0xdb, 0x1a, 0x98, 0xca, 0x92, 0x17, 0x92, 0xfe, 0x67, 0xa9, 0x3c, 0x82,
	0xa9, 0x6d, 0xbe, 0xce, 0x88, 0xcc, 0x6c, 0xd3, 0x89, 0xf5, 0x45, 0x44, 0x66, 0x6d, 0x61, 0x0c,
	0xba, 0xc2, 0xf8, 0x0c, 0x8e, 0xbe, 0x12, 0x82, 0x8b, 0x3a, 0x25, 0xa7, 0x52, 0x92, 0x94, 0xea,
	0xee, 0x63, 0xec, 0xcc, 0x3a, 0x62, 0x79, 0xb0, 0xa5, 0x1b, 0x5a, 0xfe, 0xf4, 0xe0, 0xf4, 0xe0,
	0x6d, 0xd0, 0x27, 0x07, 0xea, 0xba, 0x70, 0x6c, 0xdf, 0xf9, 0xda, 0x8d, 0xd8, 0x1e, 0xc1, 0x80,
	0x0a, 0x61, 0x15, 0x36, 0x73, 0x67, 0x34, 0xb4, 0xa8, 0x87, 0xeb, 0x18, 0xfa, 0x1c, 0xee, 0x99,
	0x15, 0x6d, 0x7d, 0x54, 0xac, 0xa0, 0xee, 0xd9, 0x5b, 0x69, 0x1f, 0x88, 0x7a, 0xf8, 0x4c, 0x1d,
	0xf8, 0x6a, 0x45, 0x54, 0xe6, 0xea, 0x5d, 0xdb, 0x1b, 0x77, 0xd8, 0x55, 0x44, 0xe7, 0x62, 0xae,
	0x15, 0x51, 0xb5, 0x1d, 0x6d, 0x45, 0xbc, 0x80, 0xf3, 0x8e, 0x22, 0x9a, 0xf7, 0x5f, 0xc0, 0x89,
	0xb0, 0xcf, 0x56, 0x1a, 0x8d, 0xfd, 0xf7, 0xda, 0xb8, 0xc6, 0x30, 0x7a, 0xae, 0xbf, 0xb5, 0x28,
	0x82, 0xf9, 0x73, 0xc1, 0x63, 0x2a, 0xa5, 0xd3, 0x5b, 0x83, 0xb0, 0xd3, 0x74, 0x71, 0x71, 0xa7,
	0xdb, 0x61, 0x09, 0x7a, 0xcf, 0x5e, 0xc0, 0xbb, 0x5c, 0xa4, 0x61, 0xb6, 0x2b, 0xa9, 0xd8, 0xd2,
	0x24, 0xa5, 0x22, 0xbc, 0x25, 0x1b, 0xc1, 0x62, 0x77, 0x50, 0xf3, 0xf0, 0xc3, 0xfb, 0x29, 0x53,
	0x59, 0xb5, 0x09, 0x63, 0x9e, 0x2f, 0x5b, 0xb9, 0x4b, 0x93, 0x6b, 0xbe, 0xf2, 0x72, 0xa9, 0x73,
	0x37, 0xe6, 0x2f, 0xc0, 0xc7, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0xe5, 0xe8, 0xf0, 0x9e, 0x1f,
	0x08, 0x00, 0x00,
}
