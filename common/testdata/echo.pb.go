


package grpclogging_test

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type Message struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	SequenceNumber       int32    `protobuf:"varint,2,opt,name=sequence_number,json=sequenceNumber,proto3" json:"sequence_number,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_4a6460a7a20c30ff, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Message) GetSequenceNumber() int32 {
	if m != nil {
		return m.SequenceNumber
	}
	return 0
}

func init() {
	proto.RegisterType((*Message)(nil), "Message")
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type EchoServiceClient interface {
	Echo(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	EchoStream(ctx context.Context, opts ...grpc.CallOption) (EchoService_EchoStreamClient, error)
}

type echoServiceClient struct {
	cc *grpc.ClientConn
}

func NewEchoServiceClient(cc *grpc.ClientConn) EchoServiceClient {
	return &echoServiceClient{cc}
}

func (c *echoServiceClient) Echo(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/EchoService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoServiceClient) EchoStream(ctx context.Context, opts ...grpc.CallOption) (EchoService_EchoStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_EchoService_serviceDesc.Streams[0], "/EchoService/EchoStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &echoServiceEchoStreamClient{stream}
	return x, nil
}

type EchoService_EchoStreamClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type echoServiceEchoStreamClient struct {
	grpc.ClientStream
}

func (x *echoServiceEchoStreamClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *echoServiceEchoStreamClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}


type EchoServiceServer interface {
	Echo(context.Context, *Message) (*Message, error)
	EchoStream(EchoService_EchoStreamServer) error
}

func RegisterEchoServiceServer(s *grpc.Server, srv EchoServiceServer) {
	s.RegisterService(&_EchoService_serviceDesc, srv)
}

func _EchoService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/EchoService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Echo(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _EchoService_EchoStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EchoServiceServer).EchoStream(&echoServiceEchoStreamServer{stream})
}

type EchoService_EchoStreamServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type echoServiceEchoStreamServer struct {
	grpc.ServerStream
}

func (x *echoServiceEchoStreamServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *echoServiceEchoStreamServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _EchoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "EchoService",
	HandlerType: (*EchoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EchoService_Echo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EchoStream",
			Handler:       _EchoService_EchoStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "testdata/echo.proto",
}

func init() { proto.RegisterFile("testdata/echo.proto", fileDescriptor_echo_4a6460a7a20c30ff) }

var fileDescriptor_echo_4a6460a7a20c30ff = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0x49, 0x2d, 0x2e,
	0x49, 0x49, 0x2c, 0x49, 0xd4, 0x4f, 0x4d, 0xce, 0xc8, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57,
	0xf2, 0xe1, 0x62, 0xf7, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0x92, 0xe0, 0x62, 0xcf, 0x85,
	0x30, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x60, 0x5c, 0x21, 0x75, 0x2e, 0xfe, 0xe2, 0xd4,
	0xc2, 0xd2, 0xd4, 0xbc, 0xe4, 0xd4, 0xf8, 0xbc, 0xd2, 0xdc, 0xa4, 0xd4, 0x22, 0x09, 0x26, 0x05,
	0x46, 0x0d, 0xd6, 0x20, 0x3e, 0x98, 0xb0, 0x1f, 0x58, 0xd4, 0xc8, 0x9f, 0x8b, 0xdb, 0x35, 0x39,
	0x23, 0x3f, 0x38, 0xb5, 0xa8, 0x2c, 0x33, 0x39, 0x55, 0x48, 0x8a, 0x8b, 0x05, 0xc4, 0x15, 0xe2,
	0xd0, 0x83, 0xda, 0x21, 0x05, 0x67, 0x09, 0xa9, 0x70, 0x71, 0x81, 0x95, 0x96, 0x14, 0xa5, 0x26,
	0xe6, 0x62, 0x53, 0xa1, 0xc1, 0x68, 0xc0, 0xe8, 0x24, 0x14, 0x25, 0x90, 0x5e, 0x54, 0x90, 0x9c,
	0x93, 0x9f, 0x9e, 0x9e, 0x99, 0x97, 0x1e, 0x0f, 0xf2, 0x41, 0x12, 0x1b, 0xd8, 0xe5, 0xc6, 0x80,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x4d, 0xa0, 0x63, 0xf4, 0xd0, 0x00, 0x00, 0x00,
}
