


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
	
	Type string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	
	Quantity             uint64   `protobuf:"varint,3,opt,name=quantity" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenToIssue) Reset()         { *m = TokenToIssue{} }
func (m *TokenToIssue) String() string { return proto.CompactTextString(m) }
func (*TokenToIssue) ProtoMessage()    {}
func (*TokenToIssue) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_886d64b71f8f88fa, []int{0}
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


type ImportRequest struct {
	
	
	Credential []byte `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	
	TokensToIssue        []*TokenToIssue `protobuf:"bytes,2,rep,name=tokens_to_issue,json=tokensToIssue" json:"tokens_to_issue,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ImportRequest) Reset()         { *m = ImportRequest{} }
func (m *ImportRequest) String() string { return proto.CompactTextString(m) }
func (*ImportRequest) ProtoMessage()    {}
func (*ImportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_886d64b71f8f88fa, []int{1}
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


type Header struct {
	
	
	Timestamp *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId" json:"channel_id,omitempty"`
	
	
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
	return fileDescriptor_prover_886d64b71f8f88fa, []int{2}
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
	
	Header *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	
	
	
	
	Payload              isCommand_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_886d64b71f8f88fa, []int{3}
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

type isCommand_Payload interface {
	isCommand_Payload()
}

type Command_ImportRequest struct {
	ImportRequest *ImportRequest `protobuf:"bytes,2,opt,name=import_request,json=importRequest,oneof"`
}

func (*Command_ImportRequest) isCommand_Payload() {}

func (m *Command) GetPayload() isCommand_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Command) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Command) GetImportRequest() *ImportRequest {
	if x, ok := m.GetPayload().(*Command_ImportRequest); ok {
		return x.ImportRequest
	}
	return nil
}


func (*Command) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Command_OneofMarshaler, _Command_OneofUnmarshaler, _Command_OneofSizer, []interface{}{
		(*Command_ImportRequest)(nil),
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
	return fileDescriptor_prover_886d64b71f8f88fa, []int{4}
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
	
	
	Timestamp *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	
	
	
	
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
	return fileDescriptor_prover_886d64b71f8f88fa, []int{5}
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
	
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
	
	Payload              []byte   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_886d64b71f8f88fa, []int{6}
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
	
	Header *CommandResponseHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	
	
	
	
	
	Payload              isCommandResponse_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *CommandResponse) Reset()         { *m = CommandResponse{} }
func (m *CommandResponse) String() string { return proto.CompactTextString(m) }
func (*CommandResponse) ProtoMessage()    {}
func (*CommandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_prover_886d64b71f8f88fa, []int{7}
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

type isCommandResponse_Payload interface {
	isCommandResponse_Payload()
}

type CommandResponse_Err struct {
	Err *Error `protobuf:"bytes,2,opt,name=err,oneof"`
}
type CommandResponse_TokenTransaction struct {
	TokenTransaction *TokenTransaction `protobuf:"bytes,3,opt,name=token_transaction,json=tokenTransaction,oneof"`
}

func (*CommandResponse_Err) isCommandResponse_Payload()              {}
func (*CommandResponse_TokenTransaction) isCommandResponse_Payload() {}

func (m *CommandResponse) GetPayload() isCommandResponse_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *CommandResponse) GetHeader() *CommandResponseHeader {
	if m != nil {
		return m.Header
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
	return fileDescriptor_prover_886d64b71f8f88fa, []int{8}
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
	proto.RegisterType((*ImportRequest)(nil), "protos.ImportRequest")
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
	err := grpc.Invoke(ctx, "/protos.Prover/ProcessCommand", in, out, c.cc, opts...)
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

func init() { proto.RegisterFile("token/prover.proto", fileDescriptor_prover_886d64b71f8f88fa) }

var fileDescriptor_prover_886d64b71f8f88fa = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x51, 0x6b, 0xdb, 0x30,
	0x10, 0x4e, 0x9a, 0x36, 0x6d, 0x2e, 0x49, 0xbb, 0x8a, 0x96, 0x85, 0xb0, 0x6e, 0xa9, 0x07, 0x23,
	0xec, 0xc1, 0x81, 0x8c, 0xc1, 0x60, 0x63, 0x8c, 0x8e, 0xb1, 0xf4, 0xad, 0xd5, 0xfa, 0x34, 0x06,
	0x41, 0xb1, 0xaf, 0xb6, 0x58, 0x2c, 0xb9, 0x92, 0x32, 0x08, 0xec, 0x37, 0x0c, 0xf6, 0x6f, 0xf6,
	0xf3, 0x86, 0x25, 0xd9, 0x71, 0x4a, 0xd9, 0xcb, 0x9e, 0xac, 0xbb, 0xfb, 0xa4, 0xfb, 0xee, 0xfb,
	0x64, 0x01, 0x31, 0xf2, 0x3b, 0x8a, 0x49, 0xae, 0xe4, 0x0f, 0x54, 0x61, 0xae, 0xa4, 0x91, 0xa4,
	0x6d, 0x3f, 0x7a, 0xf8, 0x2c, 0x91, 0x32, 0x59, 0xe2, 0xc4, 0x86, 0x8b, 0xd5, 0xed, 0xc4, 0xf0,
	0x0c, 0xb5, 0x61, 0x59, 0xee, 0x80, 0xc3, 0xc7, 0x6e, 0xb3, 0x51, 0x4c, 0x68, 0x16, 0x19, 0x2e,
	0x85, 0x2b, 0x04, 0xdf, 0xa0, 0x77, 0x53, 0x94, 0x6e, 0xe4, 0xa5, 0xd6, 0x2b, 0x24, 0x4f, 0xa0,
	0xa3, 0x30, 0xe2, 0x39, 0x47, 0x61, 0x06, 0xcd, 0x51, 0x73, 0xdc, 0xa3, 0x9b, 0x04, 0x21, 0xb0,
	0x6b, 0xd6, 0x39, 0x0e, 0x76, 0x46, 0xcd, 0x71, 0x87, 0xda, 0x35, 0x19, 0xc2, 0xc1, 0xdd, 0x8a,
	0x09, 0xc3, 0xcd, 0x7a, 0xd0, 0x1a, 0x35, 0xc7, 0xbb, 0xb4, 0x8a, 0x83, 0x0c, 0xfa, 0x97, 0x59,
	0x2e, 0x95, 0xa1, 0x78, 0xb7, 0x42, 0x6d, 0xc8, 0x53, 0x80, 0x48, 0x61, 0x8c, 0xc2, 0x70, 0xb6,
	0xf4, 0xe7, 0xd7, 0x32, 0xe4, 0x1d, 0x1c, 0x59, 0xa6, 0x7a, 0x6e, 0xe4, 0x9c, 0x17, 0x8c, 0x06,
	0x3b, 0xa3, 0xd6, 0xb8, 0x3b, 0x3d, 0x71, 0x7c, 0x75, 0x58, 0x67, 0x4b, 0xfb, 0x0e, 0xec, 0xc3,
	0xe0, 0x77, 0x13, 0xda, 0x33, 0x64, 0x31, 0x2a, 0xf2, 0x06, 0x3a, 0x95, 0x06, 0xb6, 0x4f, 0x77,
	0x3a, 0x0c, 0x9d, 0x4a, 0x61, 0xa9, 0x52, 0x78, 0x53, 0x22, 0xe8, 0x06, 0x4c, 0xce, 0x00, 0xa2,
	0x94, 0x09, 0x81, 0xcb, 0x39, 0x8f, 0xfd, 0xa4, 0x1d, 0x9f, 0xb9, 0x8c, 0xc9, 0x09, 0xec, 0x09,
	0x29, 0x22, 0xb4, 0xb3, 0xf6, 0xa8, 0x0b, 0xc8, 0x00, 0xf6, 0x23, 0x85, 0xcc, 0x48, 0x35, 0xd8,
	0xb5, 0xf9, 0x32, 0x0c, 0x7e, 0xc2, 0xfe, 0x47, 0x99, 0x65, 0x4c, 0xc4, 0xe4, 0x05, 0xb4, 0x53,
	0xcb, 0xce, 0x13, 0x3a, 0x2c, 0x67, 0x72, 0x9c, 0xa9, 0xaf, 0x92, 0xf7, 0x70, 0xc8, 0xad, 0x6a,
	0x73, 0xe5, 0x64, 0xb3, 0x2c, 0xba, 0xd3, 0xd3, 0x12, 0xbf, 0xa5, 0xe9, 0xac, 0x41, 0xfb, 0xbc,
	0x9e, 0xb8, 0xe8, 0xc0, 0x7e, 0xce, 0xd6, 0x4b, 0xc9, 0xe2, 0xe0, 0x33, 0xf4, 0xbf, 0xf0, 0x44,
	0x60, 0x5c, 0x72, 0x28, 0x88, 0xba, 0xa5, 0x57, 0xbf, 0x0c, 0x0b, 0xe7, 0x35, 0x4f, 0x04, 0x33,
	0x2b, 0xe5, 0x0c, 0xee, 0xd1, 0x4d, 0x22, 0xf8, 0xd5, 0x84, 0x53, 0x7f, 0x06, 0x45, 0x9d, 0x4b,
	0xa1, 0xf1, 0xbf, 0x95, 0x3e, 0x87, 0x9e, 0x6f, 0x3e, 0x4f, 0x99, 0x4e, 0x7d, 0xd3, 0xae, 0xcf,
	0xcd, 0x98, 0x4e, 0xeb, 0xba, 0xb6, 0xb6, 0x75, 0x7d, 0x0b, 0x7b, 0x9f, 0x94, 0x92, 0xaa, 0x80,
	0x64, 0xa8, 0x35, 0x4b, 0xd0, 0x76, 0xef, 0xd0, 0x32, 0x2c, 0x2a, 0x5e, 0x07, 0x7f, 0x74, 0x25,
	0xcb, 0x9f, 0x26, 0x1c, 0xdd, 0x9b, 0x86, 0xbc, 0xbe, 0xe7, 0xce, 0x59, 0xa9, 0xf6, 0x83, 0x63,
	0x57, 0x66, 0x9d, 0x43, 0x0b, 0x95, 0xf2, 0x0e, 0xf5, 0xcb, 0x3d, 0x96, 0xda, 0xac, 0x41, 0x8b,
	0x1a, 0xf9, 0x00, 0xc7, 0xf6, 0x9e, 0xce, 0x6b, 0xbf, 0x9f, 0x1d, 0xa7, 0x3b, 0x3d, 0xf6, 0xf7,
	0x79, 0x53, 0x98, 0x35, 0xe8, 0x23, 0x73, 0x2f, 0x57, 0x77, 0xf4, 0x1a, 0x4e, 0xb7, 0x1c, 0xad,
	0xf8, 0x0f, 0xe1, 0x40, 0xf9, 0xb5, 0xb7, 0xb6, 0x8a, 0xff, 0xed, 0xed, 0x94, 0x42, 0xfb, 0xca,
	0xbe, 0x2a, 0x64, 0x06, 0x87, 0x57, 0x4a, 0x46, 0xa8, 0x75, 0x79, 0x5f, 0xaa, 0x3b, 0xb7, 0xd5,
	0x74, 0x78, 0xf6, 0x60, 0xba, 0xe4, 0x12, 0x34, 0x2e, 0xae, 0xe1, 0xb9, 0x54, 0x49, 0x98, 0xae,
	0x73, 0x54, 0x4b, 0x8c, 0x13, 0x54, 0xe1, 0x2d, 0x5b, 0x28, 0x1e, 0x95, 0x1b, 0xed, 0x8c, 0x5f,
	0x5f, 0x26, 0xdc, 0xa4, 0xab, 0x45, 0x18, 0xc9, 0x6c, 0x52, 0xc3, 0x4e, 0x1c, 0xd6, 0xbd, 0x67,
	0x7a, 0x62, 0xb1, 0x0b, 0xf7, 0xd8, 0xbd, 0xfa, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xeb, 0xc4, 0xe5,
	0x53, 0x09, 0x05, 0x00, 0x00,
}
