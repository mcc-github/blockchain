


package token

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
	return fileDescriptor_456ae20c2189a151, []int{0}
}

func (m *TokenToIssue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenToIssue.Unmarshal(m, b)
}
func (m *TokenToIssue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenToIssue.Marshal(b, m, deterministic)
}
func (m *TokenToIssue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenToIssue.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{1}
}

func (m *RecipientTransferShare) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecipientTransferShare.Unmarshal(m, b)
}
func (m *RecipientTransferShare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecipientTransferShare.Marshal(b, m, deterministic)
}
func (m *RecipientTransferShare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecipientTransferShare.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{2}
}

func (m *TokenOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenOutput.Unmarshal(m, b)
}
func (m *TokenOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenOutput.Marshal(b, m, deterministic)
}
func (m *TokenOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenOutput.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{3}
}

func (m *UnspentTokens) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnspentTokens.Unmarshal(m, b)
}
func (m *UnspentTokens) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnspentTokens.Marshal(b, m, deterministic)
}
func (m *UnspentTokens) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnspentTokens.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{4}
}

func (m *ListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRequest.Unmarshal(m, b)
}
func (m *ListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRequest.Marshal(b, m, deterministic)
}
func (m *ListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRequest.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{5}
}

func (m *ImportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportRequest.Unmarshal(m, b)
}
func (m *ImportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportRequest.Marshal(b, m, deterministic)
}
func (m *ImportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportRequest.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{6}
}

func (m *TransferRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransferRequest.Unmarshal(m, b)
}
func (m *TransferRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransferRequest.Marshal(b, m, deterministic)
}
func (m *TransferRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferRequest.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{7}
}

func (m *RedeemRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RedeemRequest.Unmarshal(m, b)
}
func (m *RedeemRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RedeemRequest.Marshal(b, m, deterministic)
}
func (m *RedeemRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RedeemRequest.Merge(m, src)
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


type AllowanceRecipientShare struct {
	
	Recipient []byte `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	
	Quantity             uint64   `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AllowanceRecipientShare) Reset()         { *m = AllowanceRecipientShare{} }
func (m *AllowanceRecipientShare) String() string { return proto.CompactTextString(m) }
func (*AllowanceRecipientShare) ProtoMessage()    {}
func (*AllowanceRecipientShare) Descriptor() ([]byte, []int) {
	return fileDescriptor_456ae20c2189a151, []int{8}
}

func (m *AllowanceRecipientShare) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AllowanceRecipientShare.Unmarshal(m, b)
}
func (m *AllowanceRecipientShare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AllowanceRecipientShare.Marshal(b, m, deterministic)
}
func (m *AllowanceRecipientShare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllowanceRecipientShare.Merge(m, src)
}
func (m *AllowanceRecipientShare) XXX_Size() int {
	return xxx_messageInfo_AllowanceRecipientShare.Size(m)
}
func (m *AllowanceRecipientShare) XXX_DiscardUnknown() {
	xxx_messageInfo_AllowanceRecipientShare.DiscardUnknown(m)
}

var xxx_messageInfo_AllowanceRecipientShare proto.InternalMessageInfo

func (m *AllowanceRecipientShare) GetRecipient() []byte {
	if m != nil {
		return m.Recipient
	}
	return nil
}

func (m *AllowanceRecipientShare) GetQuantity() uint64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}


type ApproveRequest struct {
	
	Credential []byte `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	
	AllowanceShares []*AllowanceRecipientShare `protobuf:"bytes,2,rep,name=allowance_shares,json=allowanceShares,proto3" json:"allowance_shares,omitempty"`
	
	TokenIds             [][]byte `protobuf:"bytes,3,rep,name=token_ids,json=tokenIds,proto3" json:"token_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApproveRequest) Reset()         { *m = ApproveRequest{} }
func (m *ApproveRequest) String() string { return proto.CompactTextString(m) }
func (*ApproveRequest) ProtoMessage()    {}
func (*ApproveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_456ae20c2189a151, []int{9}
}

func (m *ApproveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApproveRequest.Unmarshal(m, b)
}
func (m *ApproveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApproveRequest.Marshal(b, m, deterministic)
}
func (m *ApproveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApproveRequest.Merge(m, src)
}
func (m *ApproveRequest) XXX_Size() int {
	return xxx_messageInfo_ApproveRequest.Size(m)
}
func (m *ApproveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ApproveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ApproveRequest proto.InternalMessageInfo

func (m *ApproveRequest) GetCredential() []byte {
	if m != nil {
		return m.Credential
	}
	return nil
}

func (m *ApproveRequest) GetAllowanceShares() []*AllowanceRecipientShare {
	if m != nil {
		return m.AllowanceShares
	}
	return nil
}

func (m *ApproveRequest) GetTokenIds() [][]byte {
	if m != nil {
		return m.TokenIds
	}
	return nil
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
	return fileDescriptor_456ae20c2189a151, []int{10}
}

func (m *Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Header.Unmarshal(m, b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Header.Marshal(b, m, deterministic)
}
func (m *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{11}
}

func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (m *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(m, src)
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

type Command_ApproveRequest struct {
	ApproveRequest *ApproveRequest `protobuf:"bytes,6,opt,name=approve_request,json=approveRequest,proto3,oneof"`
}

type Command_TransferFromRequest struct {
	TransferFromRequest *TransferRequest `protobuf:"bytes,7,opt,name=transfer_from_request,json=transferFromRequest,proto3,oneof"`
}

func (*Command_ImportRequest) isCommand_Payload() {}

func (*Command_TransferRequest) isCommand_Payload() {}

func (*Command_ListRequest) isCommand_Payload() {}

func (*Command_RedeemRequest) isCommand_Payload() {}

func (*Command_ApproveRequest) isCommand_Payload() {}

func (*Command_TransferFromRequest) isCommand_Payload() {}

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

func (m *Command) GetApproveRequest() *ApproveRequest {
	if x, ok := m.GetPayload().(*Command_ApproveRequest); ok {
		return x.ApproveRequest
	}
	return nil
}

func (m *Command) GetTransferFromRequest() *TransferRequest {
	if x, ok := m.GetPayload().(*Command_TransferFromRequest); ok {
		return x.TransferFromRequest
	}
	return nil
}


func (*Command) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Command_OneofMarshaler, _Command_OneofUnmarshaler, _Command_OneofSizer, []interface{}{
		(*Command_ImportRequest)(nil),
		(*Command_TransferRequest)(nil),
		(*Command_ListRequest)(nil),
		(*Command_RedeemRequest)(nil),
		(*Command_ApproveRequest)(nil),
		(*Command_TransferFromRequest)(nil),
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
	case *Command_ApproveRequest:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ApproveRequest); err != nil {
			return err
		}
	case *Command_TransferFromRequest:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TransferFromRequest); err != nil {
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
	case 6: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ApproveRequest)
		err := b.DecodeMessage(msg)
		m.Payload = &Command_ApproveRequest{msg}
		return true, err
	case 7: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TransferRequest)
		err := b.DecodeMessage(msg)
		m.Payload = &Command_TransferFromRequest{msg}
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
	case *Command_ApproveRequest:
		s := proto.Size(x.ApproveRequest)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Command_TransferFromRequest:
		s := proto.Size(x.TransferFromRequest)
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
	return fileDescriptor_456ae20c2189a151, []int{12}
}

func (m *SignedCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedCommand.Unmarshal(m, b)
}
func (m *SignedCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedCommand.Marshal(b, m, deterministic)
}
func (m *SignedCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedCommand.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{13}
}

func (m *CommandResponseHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandResponseHeader.Unmarshal(m, b)
}
func (m *CommandResponseHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandResponseHeader.Marshal(b, m, deterministic)
}
func (m *CommandResponseHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandResponseHeader.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{14}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{15}
}

func (m *CommandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandResponse.Unmarshal(m, b)
}
func (m *CommandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandResponse.Marshal(b, m, deterministic)
}
func (m *CommandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandResponse.Merge(m, src)
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
	return fileDescriptor_456ae20c2189a151, []int{16}
}

func (m *SignedCommandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedCommandResponse.Unmarshal(m, b)
}
func (m *SignedCommandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedCommandResponse.Marshal(b, m, deterministic)
}
func (m *SignedCommandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedCommandResponse.Merge(m, src)
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
	proto.RegisterType((*AllowanceRecipientShare)(nil), "protos.AllowanceRecipientShare")
	proto.RegisterType((*ApproveRequest)(nil), "protos.ApproveRequest")
	proto.RegisterType((*Header)(nil), "protos.Header")
	proto.RegisterType((*Command)(nil), "protos.Command")
	proto.RegisterType((*SignedCommand)(nil), "protos.SignedCommand")
	proto.RegisterType((*CommandResponseHeader)(nil), "protos.CommandResponseHeader")
	proto.RegisterType((*Error)(nil), "protos.Error")
	proto.RegisterType((*CommandResponse)(nil), "protos.CommandResponse")
	proto.RegisterType((*SignedCommandResponse)(nil), "protos.SignedCommandResponse")
}

func init() { proto.RegisterFile("token/prover.proto", fileDescriptor_456ae20c2189a151) }

var fileDescriptor_456ae20c2189a151 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0x5b, 0x6f, 0x1b, 0x45,
	0x14, 0xf6, 0xda, 0x8e, 0x13, 0x1f, 0xdf, 0xd2, 0x69, 0xd3, 0x58, 0x86, 0xb4, 0xe9, 0x22, 0xa1,
	0x88, 0x8b, 0x2d, 0x05, 0x81, 0x2a, 0x51, 0x55, 0xa4, 0xdc, 0x36, 0x88, 0x8a, 0x76, 0x62, 0x5e,
	0x10, 0x92, 0x35, 0xd9, 0x9d, 0xd8, 0x23, 0x76, 0x77, 0xb6, 0x33, 0xb3, 0x20, 0xf3, 0x03, 0x78,
	0x44, 0xe2, 0x91, 0x3f, 0xc5, 0xef, 0xe1, 0x11, 0xed, 0x5c, 0xd6, 0xbb, 0x56, 0xa0, 0x41, 0xe9,
	0x93, 0x7d, 0xce, 0x9c, 0x39, 0xe7, 0x9b, 0xef, 0x7c, 0x67, 0x66, 0x01, 0x29, 0xfe, 0x13, 0x4d,
	0x67, 0x99, 0xe0, 0x3f, 0x53, 0x31, 0xcd, 0x04, 0x57, 0x1c, 0x75, 0xf4, 0x8f, 0x9c, 0x3c, 0x5c,
	0x72, 0xbe, 0x8c, 0xe9, 0x4c, 0x9b, 0x97, 0xf9, 0xd5, 0x4c, 0xb1, 0x84, 0x4a, 0x45, 0x92, 0xcc,
	0x04, 0x4e, 0x0e, 0xcd, 0x66, 0x25, 0x48, 0x2a, 0x49, 0xa8, 0x18, 0x4f, 0xcd, 0x82, 0xff, 0x23,
	0xf4, 0xe7, 0xc5, 0xd2, 0x9c, 0x9f, 0x4b, 0x99, 0x53, 0xf4, 0x36, 0x74, 0x05, 0x0d, 0x59, 0xc6,
	0x68, 0xaa, 0xc6, 0xde, 0xb1, 0x77, 0xd2, 0xc7, 0x1b, 0x07, 0x42, 0xd0, 0x56, 0xeb, 0x8c, 0x8e,
	0x9b, 0xc7, 0xde, 0x49, 0x17, 0xeb, 0xff, 0x68, 0x02, 0x7b, 0xaf, 0x72, 0x92, 0x2a, 0xa6, 0xd6,
	0xe3, 0xd6, 0xb1, 0x77, 0xd2, 0xc6, 0xa5, 0xed, 0x63, 0xb8, 0x8f, 0xdd, 0xe6, 0x79, 0x51, 0xfb,
	0x8a, 0x8a, 0x8b, 0x15, 0x11, 0xaf, 0xab, 0x53, 0xcd, 0xd9, 0xdc, 0xca, 0xf9, 0x1c, 0x7a, 0x1a,
	0xf1, 0x77, 0xb9, 0xca, 0x72, 0x85, 0x86, 0xd0, 0x64, 0x91, 0xcd, 0xd0, 0x64, 0xd1, 0xff, 0x86,
	0xf8, 0x04, 0x06, 0xdf, 0xa7, 0x32, 0x2b, 0x00, 0x16, 0x59, 0x25, 0x7a, 0x1f, 0x3a, 0x9a, 0x2c,
	0x39, 0xf6, 0x8e, 0x5b, 0x27, 0xbd, 0xd3, 0xbb, 0x86, 0x29, 0x39, 0xad, 0x54, 0xc5, 0x36, 0xc4,
	0xff, 0x10, 0x7a, 0xdf, 0x32, 0xa9, 0x30, 0x7d, 0x95, 0x53, 0xa9, 0xd0, 0x03, 0x80, 0x50, 0xd0,
	0x88, 0xa6, 0x8a, 0x91, 0xd8, 0x82, 0xaa, 0x78, 0xfc, 0x04, 0x06, 0xe7, 0x49, 0xc6, 0xc5, 0x4d,
	0x37, 0xa0, 0x27, 0x30, 0x32, 0x95, 0x16, 0x8a, 0x2f, 0x58, 0xd1, 0xa1, 0x71, 0x53, 0xa3, 0xba,
	0x57, 0x43, 0x65, 0xbb, 0x87, 0x07, 0x26, 0xd8, 0x9a, 0xfe, 0x6f, 0x1e, 0x8c, 0x1c, 0xed, 0x37,
	0xad, 0xf8, 0x16, 0x74, 0x75, 0x92, 0x05, 0x8b, 0xa4, 0xae, 0xd5, 0xc7, 0x7b, 0xda, 0x71, 0x1e,
	0x49, 0xf4, 0x09, 0x74, 0x64, 0xd1, 0x3e, 0x39, 0x6e, 0x69, 0x14, 0x0f, 0x1c, 0x8a, 0xeb, 0xbb,
	0x8c, 0x6d, 0xb4, 0xff, 0x2b, 0x0c, 0x30, 0x8d, 0x28, 0x4d, 0xde, 0x08, 0x8a, 0x0f, 0x00, 0xb9,
	0xf6, 0x15, 0xb4, 0x08, 0x9d, 0xd9, 0x36, 0x76, 0xdf, 0xad, 0xcc, 0xb9, 0xa9, 0xe8, 0x5f, 0xc0,
	0xe1, 0x59, 0x1c, 0xf3, 0x5f, 0x48, 0x1a, 0xd2, 0x12, 0xe6, 0x6d, 0x45, 0xf8, 0xa7, 0x07, 0xc3,
	0xb3, 0x4c, 0xcf, 0xe2, 0x4d, 0x8f, 0xf4, 0x0d, 0xec, 0x13, 0x87, 0x63, 0x61, 0x59, 0x34, 0xbd,
	0x7c, 0xe8, 0x58, 0xfc, 0x17, 0x9c, 0x78, 0x54, 0x6e, 0xd4, 0xb6, 0xac, 0xd3, 0xd3, 0xaa, 0xd3,
	0xe3, 0xff, 0xe1, 0x41, 0x27, 0xa0, 0x24, 0xa2, 0x02, 0x3d, 0x86, 0x6e, 0x79, 0x13, 0x68, 0x48,
	0xbd, 0xd3, 0xc9, 0xd4, 0xdc, 0x15, 0x53, 0x77, 0x57, 0x4c, 0xe7, 0x2e, 0x02, 0x6f, 0x82, 0xd1,
	0x11, 0x40, 0xb8, 0x22, 0x69, 0x4a, 0xe3, 0x05, 0x8b, 0xec, 0x30, 0x75, 0xad, 0xe7, 0x3c, 0x42,
	0xf7, 0x60, 0x27, 0xe5, 0x69, 0x48, 0x35, 0xeb, 0x7d, 0x6c, 0x0c, 0x34, 0x86, 0xdd, 0x50, 0x50,
	0xa2, 0xb8, 0x18, 0xb7, 0xb5, 0xdf, 0x99, 0xfe, 0x5f, 0x2d, 0xd8, 0xfd, 0x9c, 0x27, 0x09, 0x49,
	0x23, 0xf4, 0x2e, 0x74, 0x56, 0x1a, 0x9e, 0x45, 0x34, 0x74, 0xc7, 0x37, 0xa0, 0xb1, 0x5d, 0x45,
	0x4f, 0x61, 0xc8, 0xf4, 0xb0, 0x2c, 0x84, 0xa1, 0x58, 0xc3, 0xe8, 0x9d, 0x1e, 0xb8, 0xf8, 0xda,
	0x28, 0x05, 0x0d, 0x3c, 0x60, 0xb5, 0xd9, 0xfa, 0x02, 0xf6, 0x95, 0x55, 0x63, 0x99, 0xa1, 0xa5,
	0x33, 0x1c, 0x96, 0xc3, 0x53, 0x1f, 0x8e, 0xa0, 0x81, 0x47, 0x6a, 0x6b, 0x5e, 0x1e, 0x43, 0x3f,
	0x66, 0x72, 0x83, 0xa1, 0xad, 0x33, 0x94, 0x97, 0x42, 0x65, 0xfa, 0x83, 0x06, 0xee, 0xc5, 0x95,
	0xcb, 0xe0, 0x29, 0x0c, 0x8d, 0x34, 0xcb, 0xbd, 0x3b, 0x75, 0xfc, 0xb5, 0x91, 0x28, 0xf0, 0x8b,
	0xda, 0x8c, 0x9c, 0xc1, 0x88, 0x18, 0x89, 0x95, 0x09, 0x3a, 0x3a, 0xc1, 0xfd, 0x52, 0x2f, 0x35,
	0x05, 0x06, 0x0d, 0x3c, 0x24, 0x75, 0x4d, 0x3e, 0x87, 0x83, 0x92, 0x82, 0x2b, 0xc1, 0x37, 0x48,
	0x76, 0x5f, 0xc7, 0xc3, 0x5d, 0xb7, 0xef, 0x2b, 0xc1, 0x1d, 0xa2, 0x67, 0x5d, 0xd8, 0xcd, 0xc8,
	0x3a, 0xe6, 0x24, 0xf2, 0xbf, 0x86, 0xc1, 0x05, 0x5b, 0xa6, 0x34, 0x72, 0x5d, 0x2d, 0x7a, 0x6f,
	0xfe, 0x5a, 0xed, 0x3b, 0xb3, 0x98, 0x32, 0xc9, 0x96, 0x29, 0x51, 0xb9, 0x30, 0xd7, 0x72, 0x1f,
	0x6f, 0x1c, 0xfe, 0xef, 0x1e, 0x1c, 0xd8, 0x1c, 0x98, 0xca, 0x8c, 0xa7, 0x92, 0xde, 0x5a, 0xbc,
	0x8f, 0xa0, 0x6f, 0x8b, 0x2f, 0x56, 0x44, 0xae, 0x6c, 0xd1, 0x9e, 0xf5, 0x05, 0x44, 0xae, 0xaa,
	0x52, 0x6d, 0xd5, 0xa5, 0xfa, 0x29, 0xec, 0x7c, 0x29, 0x04, 0x17, 0x45, 0x48, 0x42, 0xa5, 0x24,
	0x4b, 0xaa, 0xab, 0x77, 0xb1, 0x33, 0x8b, 0x15, 0xcb, 0x83, 0x4d, 0x5d, 0xd2, 0xf2, 0xb7, 0x07,
	0xa3, 0xad, 0xd3, 0xa0, 0x8f, 0xb7, 0xf4, 0x7e, 0xe4, 0x58, 0xbf, 0xf6, 0xd8, 0xa5, 0xfc, 0x1f,
	0x41, 0x8b, 0x0a, 0x61, 0x35, 0x3f, 0x70, 0x7b, 0x34, 0xb4, 0xa0, 0x81, 0x8b, 0x35, 0xf4, 0x19,
	0xdc, 0x31, 0xd7, 0x40, 0xe5, 0x5d, 0xb7, 0x12, 0xbf, 0x63, 0x1f, 0x86, 0xcd, 0x42, 0xd0, 0xc0,
	0xfb, 0x6a, 0xcb, 0x57, 0x68, 0x34, 0x37, 0xaf, 0xdf, 0xc2, 0x3e, 0x7a, 0xed, 0xba, 0x46, 0x6b,
	0x6f, 0x63, 0xa1, 0xd1, 0xbc, 0xea, 0xa8, 0x2a, 0xe2, 0x25, 0x1c, 0xd4, 0x14, 0x51, 0x9e, 0x7f,
	0x02, 0x7b, 0xc2, 0xfe, 0xb7, 0xd2, 0x28, 0xed, 0xff, 0xd6, 0xc6, 0x29, 0x86, 0xce, 0x0b, 0xfd,
	0xb9, 0x83, 0x02, 0x18, 0xbe, 0x10, 0x3c, 0xa4, 0x52, 0x3a, 0xbd, 0x95, 0x08, 0x6b, 0x45, 0x27,
	0x47, 0xd7, 0xba, 0x1d, 0x16, 0xbf, 0xf1, 0xec, 0x25, 0xbc, 0xc3, 0xc5, 0x72, 0xba, 0x5a, 0x67,
	0x54, 0xc4, 0x34, 0x5a, 0x52, 0x31, 0xbd, 0x22, 0x97, 0x82, 0x85, 0x6e, 0xa3, 0xe6, 0xe1, 0x87,
	0xf7, 0x96, 0x4c, 0xad, 0xf2, 0xcb, 0x69, 0xc8, 0x93, 0x59, 0x25, 0x76, 0x66, 0x62, 0xcd, 0x87,
	0x96, 0x9c, 0xe9, 0xd8, 0x4b, 0xf3, 0x15, 0xf6, 0xd1, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x36,
	0xe1, 0xb6, 0xa8, 0xa2, 0x09, 0x00, 0x00,
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
