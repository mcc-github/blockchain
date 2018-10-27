


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
	return fileDescriptor_prover_405073bc77f17ef9, []int{0}
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
	
	Quantity             uint64   `protobuf:"varint,2,opt,name=quantity" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RecipientTransferShare) Reset()         { *m = RecipientTransferShare{} }
func (m *RecipientTransferShare) String() string { return proto.CompactTextString(m) }
func (*RecipientTransferShare) ProtoMessage()    {}
func (*RecipientTransferShare) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_405073bc77f17ef9, []int{1}
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
	return fileDescriptor_prover_405073bc77f17ef9, []int{2}
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
	TokenIDs             [][]byte                  `protobuf:"bytes,2,rep,name=tokenIDs,proto3" json:"tokenIDs,omitempty"`
	Shares               []*RecipientTransferShare `protobuf:"bytes,3,rep,name=shares" json:"shares,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *TransferRequest) Reset()         { *m = TransferRequest{} }
func (m *TransferRequest) String() string { return proto.CompactTextString(m) }
func (*TransferRequest) ProtoMessage()    {}
func (*TransferRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_405073bc77f17ef9, []int{3}
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

func (m *TransferRequest) GetTokenIDs() [][]byte {
	if m != nil {
		return m.TokenIDs
	}
	return nil
}

func (m *TransferRequest) GetShares() []*RecipientTransferShare {
	if m != nil {
		return m.Shares
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
	return fileDescriptor_prover_405073bc77f17ef9, []int{4}
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
	return fileDescriptor_prover_405073bc77f17ef9, []int{5}
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
	TransferRequest *TransferRequest `protobuf:"bytes,3,opt,name=transfer_request,json=transferRequest,oneof"`
}

func (*Command_ImportRequest) isCommand_Payload()   {}
func (*Command_TransferRequest) isCommand_Payload() {}

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


func (*Command) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Command_OneofMarshaler, _Command_OneofUnmarshaler, _Command_OneofSizer, []interface{}{
		(*Command_ImportRequest)(nil),
		(*Command_TransferRequest)(nil),
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
	return fileDescriptor_prover_405073bc77f17ef9, []int{6}
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
	return fileDescriptor_prover_405073bc77f17ef9, []int{7}
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
	return fileDescriptor_prover_405073bc77f17ef9, []int{8}
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
	return fileDescriptor_prover_405073bc77f17ef9, []int{9}
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

func (*CommandResponse_Err) isCommandResponse_Payload() {}

func (*CommandResponse_TokenTransaction) isCommandResponse_Payload() {}

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


func (*CommandResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CommandResponse_OneofMarshaler, _CommandResponse_OneofUnmarshaler, _CommandResponse_OneofSizer, []interface{}{
		(*CommandResponse_Err)(nil),
		(*CommandResponse_TokenTransaction)(nil),
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
	return fileDescriptor_prover_405073bc77f17ef9, []int{10}
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
	proto.RegisterType((*ImportRequest)(nil), "protos.ImportRequest")
	proto.RegisterType((*TransferRequest)(nil), "protos.TransferRequest")
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

func init() { proto.RegisterFile("token/prover.proto", fileDescriptor_prover_405073bc77f17ef9) }

var fileDescriptor_prover_405073bc77f17ef9 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0x5f, 0x6b, 0xdb, 0x3a,
	0x14, 0xc0, 0xe3, 0xa6, 0x4d, 0x9b, 0x93, 0xa4, 0x69, 0x45, 0x7b, 0x1b, 0xc2, 0x6d, 0x6f, 0xea,
	0x0b, 0x97, 0x70, 0x1f, 0x1c, 0xc8, 0xd8, 0x18, 0x6c, 0x8c, 0xd1, 0x75, 0x2c, 0x79, 0x6b, 0xd5,
	0x3c, 0x8d, 0x41, 0x50, 0xec, 0x53, 0xdb, 0x2c, 0xb6, 0x5c, 0x49, 0x19, 0xe4, 0x03, 0xec, 0x75,
	0xb0, 0x6f, 0xb3, 0xe7, 0x7d, 0xb2, 0x61, 0x59, 0x72, 0xfe, 0x50, 0xb6, 0xc1, 0x9e, 0xe2, 0xf3,
	0x47, 0xe7, 0xdf, 0xef, 0x48, 0x01, 0xa2, 0xf8, 0x47, 0x4c, 0x07, 0x99, 0xe0, 0x9f, 0x50, 0x78,
	0x99, 0xe0, 0x8a, 0x93, 0x9a, 0xfe, 0x91, 0xdd, 0x7f, 0x42, 0xce, 0xc3, 0x39, 0x0e, 0xb4, 0x38,
	0x5b, 0xdc, 0x0f, 0x54, 0x9c, 0xa0, 0x54, 0x2c, 0xc9, 0x0a, 0xc7, 0xee, 0x59, 0x71, 0x58, 0x09,
	0x96, 0x4a, 0xe6, 0xab, 0x98, 0xa7, 0x85, 0xc1, 0xfd, 0x00, 0xcd, 0x49, 0x6e, 0x9a, 0xf0, 0xb1,
	0x94, 0x0b, 0x24, 0x7f, 0x43, 0x5d, 0xa0, 0x1f, 0x67, 0x31, 0xa6, 0xaa, 0xe3, 0xf4, 0x9c, 0x7e,
	0x93, 0xae, 0x14, 0x84, 0xc0, 0xae, 0x5a, 0x66, 0xd8, 0xd9, 0xe9, 0x39, 0xfd, 0x3a, 0xd5, 0xdf,
	0xa4, 0x0b, 0x07, 0x0f, 0x0b, 0x96, 0xaa, 0x58, 0x2d, 0x3b, 0xd5, 0x9e, 0xd3, 0xdf, 0xa5, 0xa5,
	0xec, 0x52, 0xf8, 0x8b, 0xda, 0xc3, 0x93, 0x3c, 0xf7, 0x3d, 0x8a, 0xbb, 0x88, 0x89, 0x5f, 0xe5,
	0x59, 0x8f, 0xb9, 0xb3, 0x15, 0x33, 0x81, 0xd6, 0x38, 0xc9, 0xb8, 0x50, 0x14, 0x1f, 0x16, 0x28,
	0x15, 0xb9, 0x00, 0xf0, 0x05, 0x06, 0x98, 0xaa, 0x98, 0xcd, 0x4d, 0xac, 0x35, 0x0d, 0x79, 0x09,
	0x6d, 0xdd, 0xbd, 0x9c, 0x2a, 0x3e, 0x8d, 0xf3, 0x2e, 0x3b, 0x3b, 0xbd, 0x6a, 0xbf, 0x31, 0x3c,
	0x29, 0x66, 0x20, 0xbd, 0xf5, 0x09, 0xd0, 0x56, 0xe1, 0x6c, 0x44, 0xf7, 0xb3, 0x03, 0x6d, 0x5b,
	0xfa, 0xef, 0x66, 0xec, 0xc2, 0x81, 0x0e, 0x32, 0xbe, 0x96, 0x3a, 0x55, 0x93, 0x96, 0x32, 0x79,
	0x06, 0x35, 0x99, 0x4f, 0x40, 0x76, 0xaa, 0xba, 0x88, 0x0b, 0x5b, 0xc4, 0xe3, 0x83, 0xa2, 0xc6,
	0xdb, 0xfd, 0xea, 0x40, 0x6d, 0x84, 0x2c, 0x40, 0x41, 0x9e, 0x43, 0xbd, 0xe4, 0xab, 0xb3, 0x37,
	0x86, 0x5d, 0xaf, 0xd8, 0x00, 0xcf, 0x6e, 0x80, 0x37, 0xb1, 0x1e, 0x74, 0xe5, 0x4c, 0xce, 0x01,
	0xfc, 0x88, 0xa5, 0x29, 0xce, 0xa7, 0x71, 0x60, 0x28, 0xd6, 0x8d, 0x66, 0x1c, 0x90, 0x13, 0xd8,
	0x4b, 0x79, 0xea, 0xa3, 0xe6, 0xd8, 0xa4, 0x85, 0x40, 0x3a, 0xb0, 0xef, 0x0b, 0x64, 0x8a, 0x8b,
	0xce, 0xae, 0xd6, 0x5b, 0xd1, 0xfd, 0xee, 0xc0, 0xfe, 0x1b, 0x9e, 0x24, 0x2c, 0x0d, 0xc8, 0x7f,
	0x50, 0x8b, 0x74, 0x79, 0xa6, 0xa2, 0x43, 0xdb, 0x57, 0x51, 0x34, 0x35, 0x56, 0xf2, 0x0a, 0x0e,
	0x63, 0x8d, 0x6f, 0x2a, 0x8a, 0x69, 0xea, 0x32, 0x1a, 0xc3, 0x53, 0xeb, 0xbf, 0x01, 0x77, 0x54,
	0xa1, 0xad, 0x78, 0x83, 0xf6, 0x35, 0x1c, 0x29, 0x33, 0xa0, 0x32, 0x42, 0x55, 0x47, 0x38, 0x2b,
	0x71, 0x6e, 0xe2, 0x1a, 0x55, 0x68, 0x5b, 0x6d, 0xaa, 0xae, 0xea, 0xb0, 0x9f, 0xb1, 0xe5, 0x9c,
	0xb3, 0xc0, 0x7d, 0x07, 0xad, 0xbb, 0x38, 0x4c, 0x31, 0xb0, 0x9d, 0xe4, 0xfd, 0x16, 0x9f, 0x06,
	0xad, 0x15, 0xf3, 0xa5, 0x95, 0x71, 0x98, 0x32, 0xb5, 0x10, 0xc5, 0x1d, 0x68, 0xd2, 0x95, 0xc2,
	0xfd, 0xe2, 0xc0, 0xa9, 0x89, 0x41, 0x51, 0x66, 0x3c, 0x95, 0xf8, 0xc7, 0xc0, 0x2e, 0xa1, 0x69,
	0x92, 0x4f, 0x23, 0x26, 0x23, 0x93, 0xb4, 0x61, 0x74, 0x23, 0x26, 0xa3, 0x75, 0x3c, 0xd5, 0x4d,
	0x3c, 0x2f, 0x60, 0xef, 0xad, 0x10, 0x5c, 0xe4, 0x2e, 0x09, 0x4a, 0xc9, 0x42, 0xd4, 0xd9, 0xeb,
	0xd4, 0x8a, 0xb9, 0xc5, 0xcc, 0xc1, 0x84, 0x2e, 0xc7, 0xf2, 0xcd, 0x81, 0xf6, 0x56, 0x37, 0xe4,
	0xe9, 0x16, 0xe3, 0x73, 0x3b, 0xf1, 0x47, 0xdb, 0x2e, 0x91, 0x5f, 0x42, 0x15, 0x85, 0x30, 0x9c,
	0x5b, 0xf6, 0x8c, 0x2e, 0x6d, 0x54, 0xa1, 0xb9, 0x8d, 0xbc, 0x86, 0x63, 0x7d, 0x43, 0xa6, 0x6b,
	0x2f, 0x94, 0xc1, 0x7a, 0x6c, 0xae, 0xe7, 0xca, 0x30, 0xaa, 0xd0, 0x23, 0xb5, 0xa5, 0x5b, 0x27,
	0x7a, 0x0b, 0xa7, 0x1b, 0x44, 0xcb, 0xfa, 0xbb, 0x70, 0x20, 0xcc, 0xb7, 0x41, 0x5b, 0xca, 0x3f,
	0x67, 0x3b, 0xa4, 0x50, 0xbb, 0xd1, 0x0f, 0x2f, 0x19, 0xc1, 0xe1, 0x8d, 0xe0, 0x3e, 0x4a, 0x69,
	0xf7, 0xa5, 0xdc, 0xdc, 0x8d, 0xa4, 0xdd, 0xf3, 0x47, 0xd5, 0xb6, 0x16, 0xb7, 0x72, 0x75, 0x0b,
	0xff, 0x72, 0x11, 0x7a, 0xd1, 0x32, 0x43, 0x31, 0xc7, 0x20, 0x44, 0xe1, 0xdd, 0xb3, 0x99, 0x88,
	0x7d, 0x7b, 0x50, 0xf7, 0xf8, 0xfe, 0xff, 0x30, 0x56, 0xd1, 0x62, 0xe6, 0xf9, 0x3c, 0x19, 0xac,
	0xf9, 0x0e, 0x0a, 0xdf, 0xe2, 0xc9, 0x97, 0x03, 0xed, 0x3b, 0x2b, 0xfe, 0x0f, 0x9e, 0xfc, 0x08,
	0x00, 0x00, 0xff, 0xff, 0xf0, 0x1f, 0x1b, 0x93, 0x2c, 0x06, 0x00, 0x00,
}
