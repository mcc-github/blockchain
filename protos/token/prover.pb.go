


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

func (*Command_ImportRequest) isCommand_Payload() {}

func (*Command_TransferRequest) isCommand_Payload() {}

func (*Command_ListRequest) isCommand_Payload() {}

func (*Command_RedeemRequest) isCommand_Payload() {}

func (*Command_ApproveRequest) isCommand_Payload() {}

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


func (*Command) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Command_OneofMarshaler, _Command_OneofUnmarshaler, _Command_OneofSizer, []interface{}{
		(*Command_ImportRequest)(nil),
		(*Command_TransferRequest)(nil),
		(*Command_ListRequest)(nil),
		(*Command_RedeemRequest)(nil),
		(*Command_ApproveRequest)(nil),
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
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xdd, 0x6f, 0xe3, 0x44,
	0x10, 0x8f, 0x93, 0x34, 0x6d, 0x26, 0x5f, 0xbd, 0xe5, 0x7a, 0x8d, 0x02, 0xbd, 0xeb, 0x19, 0x09,
	0x55, 0x7c, 0x24, 0x52, 0x11, 0xe8, 0x24, 0x4e, 0x27, 0x7a, 0x80, 0x70, 0x11, 0x88, 0xbb, 0x6d,
	0x78, 0x41, 0x48, 0xd1, 0xd6, 0xde, 0xc6, 0x2b, 0x6c, 0xaf, 0x6f, 0x77, 0x0d, 0x2a, 0x7f, 0x00,
	0x8f, 0x48, 0x3c, 0xf2, 0x9f, 0xf2, 0x06, 0xf2, 0x7e, 0x38, 0x76, 0x54, 0xa0, 0xa8, 0xf7, 0x94,
	0xcc, 0xec, 0xec, 0xcc, 0x6f, 0x7f, 0xf3, 0x9b, 0x5d, 0x03, 0x52, 0xfc, 0x47, 0x9a, 0x2d, 0x72,
	0xc1, 0x7f, 0xa2, 0x62, 0x9e, 0x0b, 0xae, 0x38, 0xea, 0xe9, 0x1f, 0x39, 0x7b, 0xb4, 0xe6, 0x7c,
	0x9d, 0xd0, 0x85, 0x36, 0x2f, 0x8b, 0xab, 0x85, 0x62, 0x29, 0x95, 0x8a, 0xa4, 0xb9, 0x09, 0x9c,
	0x1d, 0x9a, 0xcd, 0x4a, 0x90, 0x4c, 0x92, 0x50, 0x31, 0x9e, 0x99, 0x05, 0xff, 0x07, 0x18, 0x2e,
	0xcb, 0xa5, 0x25, 0x3f, 0x97, 0xb2, 0xa0, 0xe8, 0x2d, 0xe8, 0x0b, 0x1a, 0xb2, 0x9c, 0xd1, 0x4c,
	0x4d, 0xbd, 0x63, 0xef, 0x64, 0x88, 0x37, 0x0e, 0x84, 0xa0, 0xab, 0xae, 0x73, 0x3a, 0x6d, 0x1f,
	0x7b, 0x27, 0x7d, 0xac, 0xff, 0xa3, 0x19, 0xec, 0xbd, 0x2a, 0x48, 0xa6, 0x98, 0xba, 0x9e, 0x76,
	0x8e, 0xbd, 0x93, 0x2e, 0xae, 0x6c, 0x1f, 0xc3, 0x03, 0xec, 0x36, 0x2f, 0xcb, 0xda, 0x57, 0x54,
	0x5c, 0xc4, 0x44, 0xfc, 0x57, 0x9d, 0x7a, 0xce, 0xf6, 0x56, 0xce, 0x6f, 0x60, 0xa0, 0x11, 0x7f,
	0x5b, 0xa8, 0xbc, 0x50, 0x68, 0x0c, 0x6d, 0x16, 0xd9, 0x0c, 0x6d, 0x16, 0xfd, 0x6f, 0x88, 0x4f,
	0x61, 0xf4, 0x5d, 0x26, 0xf3, 0x12, 0x60, 0x99, 0x55, 0xa2, 0xf7, 0xa0, 0xa7, 0xc9, 0x92, 0x53,
	0xef, 0xb8, 0x73, 0x32, 0x38, 0x7d, 0xc3, 0x30, 0x25, 0xe7, 0xb5, 0xaa, 0xd8, 0x86, 0xf8, 0x1f,
	0xc0, 0xe0, 0x6b, 0x26, 0x15, 0xa6, 0xaf, 0x0a, 0x2a, 0x15, 0x7a, 0x08, 0x10, 0x0a, 0x1a, 0xd1,
	0x4c, 0x31, 0x92, 0x58, 0x50, 0x35, 0x8f, 0x9f, 0xc2, 0xe8, 0x3c, 0xcd, 0xb9, 0xb8, 0xed, 0x06,
	0xf4, 0x14, 0x26, 0xa6, 0xd2, 0x4a, 0xf1, 0x15, 0x2b, 0x3b, 0x34, 0x6d, 0x6b, 0x54, 0xf7, 0x1b,
	0xa8, 0x6c, 0xf7, 0xf0, 0xc8, 0x04, 0x5b, 0xd3, 0xff, 0xd5, 0x83, 0x89, 0xa3, 0xfd, 0xb6, 0x15,
	0xdf, 0x84, 0xbe, 0x4e, 0xb2, 0x62, 0x91, 0xd4, 0xb5, 0x86, 0x78, 0x4f, 0x3b, 0xce, 0x23, 0x89,
	0x3e, 0x86, 0x9e, 0x2c, 0xdb, 0x27, 0xa7, 0x1d, 0x8d, 0xe2, 0xa1, 0x43, 0x71, 0x73, 0x97, 0xb1,
	0x8d, 0xf6, 0x7f, 0x81, 0x11, 0xa6, 0x11, 0xa5, 0xe9, 0x6b, 0x41, 0xf1, 0x3e, 0x20, 0xd7, 0xbe,
	0x92, 0x16, 0xa1, 0x33, 0xdb, 0xc6, 0xee, 0xbb, 0x95, 0x25, 0x37, 0x15, 0xfd, 0x0b, 0x38, 0x3c,
	0x4b, 0x12, 0xfe, 0x33, 0xc9, 0x42, 0x5a, 0xc1, 0xbc, 0xab, 0x08, 0xff, 0xf0, 0x60, 0x7c, 0x96,
	0xeb, 0x59, 0xbc, 0xed, 0x91, 0xbe, 0x82, 0x7d, 0xe2, 0x70, 0xac, 0x2c, 0x8b, 0xa6, 0x97, 0x8f,
	0x1c, 0x8b, 0xff, 0x80, 0x13, 0x4f, 0xaa, 0x8d, 0xda, 0x96, 0x4d, 0x7a, 0x3a, 0x4d, 0x7a, 0xfc,
	0xdf, 0x3d, 0xe8, 0x05, 0x94, 0x44, 0x54, 0xa0, 0x27, 0xd0, 0xaf, 0x6e, 0x02, 0x0d, 0x69, 0x70,
	0x3a, 0x9b, 0x9b, 0xbb, 0x62, 0xee, 0xee, 0x8a, 0xf9, 0xd2, 0x45, 0xe0, 0x4d, 0x30, 0x3a, 0x02,
	0x08, 0x63, 0x92, 0x65, 0x34, 0x59, 0xb1, 0xc8, 0x0e, 0x53, 0xdf, 0x7a, 0xce, 0x23, 0x74, 0x1f,
	0x76, 0x32, 0x9e, 0x85, 0x54, 0xb3, 0x3e, 0xc4, 0xc6, 0x40, 0x53, 0xd8, 0x0d, 0x05, 0x25, 0x8a,
	0x8b, 0x69, 0x57, 0xfb, 0x9d, 0xe9, 0xff, 0xd5, 0x86, 0xdd, 0xcf, 0x78, 0x9a, 0x92, 0x2c, 0x42,
	0xef, 0x40, 0x2f, 0xd6, 0xf0, 0x2c, 0xa2, 0xb1, 0x3b, 0xbe, 0x01, 0x8d, 0xed, 0x2a, 0x7a, 0x06,
	0x63, 0xa6, 0x87, 0x65, 0x25, 0x0c, 0xc5, 0x1a, 0xc6, 0xe0, 0xf4, 0xc0, 0xc5, 0x37, 0x46, 0x29,
	0x68, 0xe1, 0x11, 0x6b, 0xcc, 0xd6, 0xe7, 0xb0, 0xaf, 0xac, 0x1a, 0xab, 0x0c, 0x1d, 0x9d, 0xe1,
	0xb0, 0x1a, 0x9e, 0xe6, 0x70, 0x04, 0x2d, 0x3c, 0x51, 0x5b, 0xf3, 0xf2, 0x04, 0x86, 0x09, 0x93,
	0x1b, 0x0c, 0x5d, 0x9d, 0xa1, 0xba, 0x14, 0x6a, 0xd3, 0x1f, 0xb4, 0xf0, 0x20, 0xa9, 0x5d, 0x06,
	0xcf, 0x60, 0x6c, 0xa4, 0x59, 0xed, 0xdd, 0x69, 0xe2, 0x6f, 0x8c, 0x44, 0x89, 0x5f, 0x34, 0x66,
	0xe4, 0x0c, 0x26, 0xc4, 0x48, 0xac, 0x4a, 0xd0, 0xd3, 0x09, 0x1e, 0x54, 0x7a, 0x69, 0x28, 0x30,
	0x68, 0xe1, 0x31, 0x69, 0x78, 0x9e, 0xf7, 0x61, 0x37, 0x27, 0xd7, 0x09, 0x27, 0x91, 0xff, 0x25,
	0x8c, 0x2e, 0xd8, 0x3a, 0xa3, 0x91, 0x6b, 0x43, 0xd9, 0x2c, 0xf3, 0xd7, 0x8a, 0xd5, 0x99, 0xe5,
	0x58, 0x48, 0xb6, 0xce, 0x88, 0x2a, 0x84, 0xb9, 0x47, 0x87, 0x78, 0xe3, 0xf0, 0x7f, 0xf3, 0xe0,
	0xc0, 0xe6, 0xc0, 0x54, 0xe6, 0x3c, 0x93, 0xf4, 0xce, 0x6a, 0x7b, 0x0c, 0x43, 0x5b, 0x7c, 0x15,
	0x13, 0x19, 0xdb, 0xa2, 0x03, 0xeb, 0x0b, 0x88, 0x8c, 0xeb, 0xda, 0xea, 0x34, 0xb5, 0xf5, 0x09,
	0xec, 0x7c, 0x21, 0x04, 0x17, 0x65, 0x48, 0x4a, 0xa5, 0x24, 0x6b, 0xaa, 0xab, 0xf7, 0xb1, 0x33,
	0xcb, 0x15, 0xcb, 0x83, 0x4d, 0x5d, 0xd1, 0xf2, 0xa7, 0x07, 0x93, 0xad, 0xd3, 0xa0, 0x8f, 0xb6,
	0x04, 0x7a, 0xe4, 0xf8, 0xbe, 0xf1, 0xd8, 0x95, 0x5e, 0x1f, 0x43, 0x87, 0x0a, 0x61, 0x45, 0x3a,
	0x72, 0x7b, 0x34, 0xb4, 0xa0, 0x85, 0xcb, 0x35, 0xf4, 0x29, 0xdc, 0x33, 0x73, 0x5b, 0x7b, 0x88,
	0xad, 0x26, 0xef, 0xd9, 0x9b, 0x7c, 0xb3, 0x10, 0xb4, 0xf0, 0xbe, 0xda, 0xf2, 0x95, 0xa2, 0x2a,
	0xcc, 0x73, 0xb5, 0xb2, 0xaf, 0x54, 0xb7, 0x29, 0xaa, 0xc6, 0x63, 0x56, 0x8a, 0xaa, 0xa8, 0x3b,
	0xea, 0x8a, 0x78, 0x09, 0x07, 0x0d, 0x45, 0x54, 0xe7, 0x9f, 0xc1, 0x9e, 0xb0, 0xff, 0xad, 0x34,
	0x2a, 0xfb, 0xdf, 0xb5, 0x71, 0x8a, 0xa1, 0xf7, 0x42, 0x7f, 0x9f, 0xa0, 0x00, 0xc6, 0x2f, 0x04,
	0x0f, 0xa9, 0x94, 0x4e, 0x6f, 0x15, 0xc2, 0x46, 0xd1, 0xd9, 0xd1, 0x8d, 0x6e, 0x87, 0xc5, 0x6f,
	0x3d, 0x7f, 0x09, 0x6f, 0x73, 0xb1, 0x9e, 0xc7, 0xd7, 0x39, 0x15, 0x09, 0x8d, 0xd6, 0x54, 0xcc,
	0xaf, 0xc8, 0xa5, 0x60, 0xa1, 0xdb, 0xa8, 0x79, 0xf8, 0xfe, 0xdd, 0x35, 0x53, 0x71, 0x71, 0x39,
	0x0f, 0x79, 0xba, 0xa8, 0xc5, 0x2e, 0x4c, 0xac, 0xf9, 0x32, 0x92, 0x0b, 0x1d, 0x7b, 0x69, 0x3e,
	0x9b, 0x3e, 0xfc, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xbf, 0xbe, 0x53, 0x91, 0x53, 0x09, 0x00, 0x00,
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
