


package testpb

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Echo struct {
	Payload              []byte   `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Echo) Reset()         { *m = Echo{} }
func (m *Echo) String() string { return proto.CompactTextString(m) }
func (*Echo) ProtoMessage()    {}
func (*Echo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{1}
}

func (m *Echo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Echo.Unmarshal(m, b)
}
func (m *Echo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Echo.Marshal(b, m, deterministic)
}
func (m *Echo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Echo.Merge(m, src)
}
func (m *Echo) XXX_Size() int {
	return xxx_messageInfo_Echo.Size(m)
}
func (m *Echo) XXX_DiscardUnknown() {
	xxx_messageInfo_Echo.DiscardUnknown(m)
}

var xxx_messageInfo_Echo proto.InternalMessageInfo

func (m *Echo) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "Empty")
	proto.RegisterType((*Echo)(nil), "Echo")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x31, 0x4f, 0x87, 0x30,
	0x14, 0xc4, 0xd3, 0xc4, 0x3f, 0xe8, 0x83, 0xa9, 0x13, 0xc1, 0x05, 0x59, 0x24, 0x0e, 0x2d, 0xc1,
	0x6f, 0xa0, 0x61, 0x75, 0x50, 0x27, 0xb7, 0xb6, 0x3c, 0x81, 0xa4, 0x4d, 0x9b, 0x52, 0x4d, 0xf8,
	0xf6, 0xa6, 0x45, 0x16, 0x27, 0xa7, 0xfb, 0x5d, 0x7a, 0xbd, 0xbc, 0x1c, 0x40, 0xc0, 0x2d, 0x30,
	0xe7, 0x6d, 0xb0, 0x6d, 0x0e, 0x97, 0xd1, 0xb8, 0xb0, 0xb7, 0x0d, 0x5c, 0x8d, 0x6a, 0xb1, 0xb4,
	0x82, 0xdc, 0x89, 0x5d, 0x5b, 0x31, 0x55, 0xa4, 0x21, 0x5d, 0xf9, 0x7a, 0xda, 0xe1, 0x01, 0x8a,
	0x77, 0xdc, 0xc2, 0x1b, 0xfa, 0xef, 0x55, 0x21, 0xbd, 0x85, 0x9b, 0xf4, 0xf3, 0x59, 0x68, 0x4d,
	0x33, 0x96, 0xb8, 0xfe, 0xd5, 0xe1, 0x05, 0xca, 0x04, 0xff, 0x09, 0xd3, 0x3b, 0x28, 0x8e, 0x70,
	0xf0, 0x28, 0xcc, 0xdf, 0xe7, 0x8e, 0xf4, 0x64, 0xb8, 0x87, 0x22, 0x5e, 0x77, 0xd6, 0x55, 0x70,
	0x1d, 0x6d, 0x6a, 0xbb, 0xb0, 0x88, 0xf5, 0x21, 0x4f, 0xfd, 0x07, 0x9b, 0xd7, 0xb0, 0x7c, 0x49,
	0xa6, 0xac, 0xe1, 0xcb, 0xee, 0xd0, 0x6b, 0x9c, 0x66, 0xf4, 0xfc, 0x53, 0x48, 0xbf, 0x2a, 0xae,
	0xac, 0x47, 0xae, 0xac, 0x31, 0x3c, 0xae, 0xe0, 0xa4, 0xcc, 0xd2, 0x10, 0x8f, 0x3f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x34, 0x99, 0xfa, 0xcc, 0x16, 0x01, 0x00, 0x00,
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type TestServiceClient interface {
	EmptyCall(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type testServiceClient struct {
	cc *grpc.ClientConn
}

func NewTestServiceClient(cc *grpc.ClientConn) TestServiceClient {
	return &testServiceClient{cc}
}

func (c *testServiceClient) EmptyCall(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/TestService/EmptyCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}


type TestServiceServer interface {
	EmptyCall(context.Context, *Empty) (*Empty, error)
}

func RegisterTestServiceServer(s *grpc.Server, srv TestServiceServer) {
	s.RegisterService(&_TestService_serviceDesc, srv)
}

func _TestService_EmptyCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).EmptyCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/TestService/EmptyCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).EmptyCall(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _TestService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "TestService",
	HandlerType: (*TestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EmptyCall",
			Handler:    _TestService_EmptyCall_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test.proto",
}




type EmptyServiceClient interface {
	EmptyCall(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	EmptyStream(ctx context.Context, opts ...grpc.CallOption) (EmptyService_EmptyStreamClient, error)
}

type emptyServiceClient struct {
	cc *grpc.ClientConn
}

func NewEmptyServiceClient(cc *grpc.ClientConn) EmptyServiceClient {
	return &emptyServiceClient{cc}
}

func (c *emptyServiceClient) EmptyCall(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/EmptyService/EmptyCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emptyServiceClient) EmptyStream(ctx context.Context, opts ...grpc.CallOption) (EmptyService_EmptyStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_EmptyService_serviceDesc.Streams[0], "/EmptyService/EmptyStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &emptyServiceEmptyStreamClient{stream}
	return x, nil
}

type EmptyService_EmptyStreamClient interface {
	Send(*Empty) error
	Recv() (*Empty, error)
	grpc.ClientStream
}

type emptyServiceEmptyStreamClient struct {
	grpc.ClientStream
}

func (x *emptyServiceEmptyStreamClient) Send(m *Empty) error {
	return x.ClientStream.SendMsg(m)
}

func (x *emptyServiceEmptyStreamClient) Recv() (*Empty, error) {
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}


type EmptyServiceServer interface {
	EmptyCall(context.Context, *Empty) (*Empty, error)
	EmptyStream(EmptyService_EmptyStreamServer) error
}

func RegisterEmptyServiceServer(s *grpc.Server, srv EmptyServiceServer) {
	s.RegisterService(&_EmptyService_serviceDesc, srv)
}

func _EmptyService_EmptyCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmptyServiceServer).EmptyCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/EmptyService/EmptyCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmptyServiceServer).EmptyCall(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmptyService_EmptyStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EmptyServiceServer).EmptyStream(&emptyServiceEmptyStreamServer{stream})
}

type EmptyService_EmptyStreamServer interface {
	Send(*Empty) error
	Recv() (*Empty, error)
	grpc.ServerStream
}

type emptyServiceEmptyStreamServer struct {
	grpc.ServerStream
}

func (x *emptyServiceEmptyStreamServer) Send(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *emptyServiceEmptyStreamServer) Recv() (*Empty, error) {
	m := new(Empty)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _EmptyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "EmptyService",
	HandlerType: (*EmptyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EmptyCall",
			Handler:    _EmptyService_EmptyCall_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EmptyStream",
			Handler:       _EmptyService_EmptyStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "test.proto",
}




type EchoServiceClient interface {
	EchoCall(ctx context.Context, in *Echo, opts ...grpc.CallOption) (*Echo, error)
}

type echoServiceClient struct {
	cc *grpc.ClientConn
}

func NewEchoServiceClient(cc *grpc.ClientConn) EchoServiceClient {
	return &echoServiceClient{cc}
}

func (c *echoServiceClient) EchoCall(ctx context.Context, in *Echo, opts ...grpc.CallOption) (*Echo, error) {
	out := new(Echo)
	err := c.cc.Invoke(ctx, "/EchoService/EchoCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}


type EchoServiceServer interface {
	EchoCall(context.Context, *Echo) (*Echo, error)
}

func RegisterEchoServiceServer(s *grpc.Server, srv EchoServiceServer) {
	s.RegisterService(&_EchoService_serviceDesc, srv)
}

func _EchoService_EchoCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Echo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).EchoCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/EchoService/EchoCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).EchoCall(ctx, req.(*Echo))
	}
	return interceptor(ctx, in, info, handler)
}

var _EchoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "EchoService",
	HandlerType: (*EchoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EchoCall",
			Handler:    _EchoService_EchoCall_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test.proto",
}
