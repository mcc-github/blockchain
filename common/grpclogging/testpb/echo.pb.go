


package testpb

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
	Sequence             int32    `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_c06a4742710c8217, []int{0}
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

func (m *Message) GetSequence() int32 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func init() {
	proto.RegisterType((*Message)(nil), "testpb.Message")
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
	err := c.cc.Invoke(ctx, "/testpb.EchoService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoServiceClient) EchoStream(ctx context.Context, opts ...grpc.CallOption) (EchoService_EchoStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_EchoService_serviceDesc.Streams[0], "/testpb.EchoService/EchoStream", opts...)
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
		FullMethod: "/testpb.EchoService/Echo",
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
	ServiceName: "testpb.EchoService",
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
	Metadata: "testpb/echo.proto",
}

func init() { proto.RegisterFile("testpb/echo.proto", fileDescriptor_echo_c06a4742710c8217) }

var fileDescriptor_echo_c06a4742710c8217 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x49, 0x2d, 0x2e,
	0x29, 0x48, 0xd2, 0x4f, 0x4d, 0xce, 0xc8, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83,
	0x08, 0x29, 0xd9, 0x73, 0xb1, 0xfb, 0xa6, 0x16, 0x17, 0x27, 0xa6, 0xa7, 0x0a, 0x49, 0x70, 0xb1,
	0xe7, 0x42, 0x98, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x30, 0xae, 0x90, 0x14, 0x17, 0x47,
	0x71, 0x6a, 0x61, 0x69, 0x6a, 0x5e, 0x72, 0xaa, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x6b, 0x10, 0x9c,
	0x6f, 0x94, 0xcd, 0xc5, 0xed, 0x9a, 0x9c, 0x91, 0x1f, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a,
	0xa4, 0xc1, 0xc5, 0x02, 0xe2, 0x0a, 0xf1, 0xeb, 0x41, 0x2c, 0xd0, 0x83, 0x9a, 0x2e, 0x85, 0x2e,
	0x20, 0x64, 0xc4, 0xc5, 0x05, 0xd6, 0x58, 0x52, 0x94, 0x9a, 0x98, 0x4b, 0x58, 0xbd, 0x06, 0xa3,
	0x01, 0xa3, 0x13, 0x47, 0x14, 0xd4, 0xdd, 0x49, 0x6c, 0x60, 0x6f, 0x18, 0x03, 0x02, 0x00, 0x00,
	0xff, 0xff, 0xee, 0x2c, 0xb0, 0x05, 0xdb, 0x00, 0x00, 0x00,
}
