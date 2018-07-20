


package grpc 

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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_test_7060a545c9ec914d, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
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
	return fileDescriptor_test_7060a545c9ec914d, []int{1}
}
func (m *Echo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Echo.Unmarshal(m, b)
}
func (m *Echo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Echo.Marshal(b, m, deterministic)
}
func (dst *Echo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Echo.Merge(dst, src)
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
	err := grpc.Invoke(ctx, "/TestService/EmptyCall", in, out, c.cc, opts...)
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
	err := grpc.Invoke(ctx, "/EchoService/EchoCall", in, out, c.cc, opts...)
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

func init() { proto.RegisterFile("test.proto", fileDescriptor_test_7060a545c9ec914d) }

var fileDescriptor_test_7060a545c9ec914d = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8e, 0x41, 0x8b, 0x83, 0x30,
	0x10, 0x85, 0x11, 0x56, 0xdd, 0x1d, 0xf7, 0xe4, 0x49, 0xdc, 0x8b, 0x78, 0xd9, 0xd2, 0x43, 0x02,
	0x96, 0xd2, 0x7b, 0x8b, 0x7f, 0xa0, 0xed, 0xa9, 0xb7, 0x18, 0xa7, 0x2a, 0x24, 0x24, 0xc4, 0x69,
	0xc1, 0x7f, 0x5f, 0x92, 0xd6, 0xd3, 0x7b, 0x0f, 0xbe, 0x19, 0x3e, 0x00, 0xc2, 0x99, 0x98, 0x75,
	0x86, 0x4c, 0x9d, 0x42, 0xdc, 0x6a, 0x4b, 0x4b, 0x5d, 0xc1, 0x57, 0x2b, 0x47, 0x93, 0x17, 0x90,
	0x5a, 0xb1, 0x28, 0x23, 0xfa, 0x22, 0xaa, 0xa2, 0xcd, 0xef, 0x79, 0x9d, 0xcd, 0x16, 0xb2, 0x2b,
	0xce, 0x74, 0x41, 0xf7, 0x9c, 0x24, 0xe6, 0x7f, 0xf0, 0x13, 0x2e, 0x4f, 0x42, 0xa9, 0x3c, 0x61,
	0xa1, 0x97, 0x9f, 0x6c, 0xfe, 0x21, 0xf3, 0xdf, 0x56, 0xb6, 0x80, 0x6f, 0x3f, 0x03, 0x1a, 0x33,
	0x5f, 0xcb, 0x77, 0x1c, 0x0f, 0xb7, 0xfd, 0x30, 0xd1, 0xf8, 0xe8, 0x98, 0x34, 0x9a, 0x8f, 0x8b,
	0x45, 0xa7, 0xb0, 0x1f, 0xd0, 0xf1, 0xbb, 0xe8, 0xdc, 0x24, 0xb9, 0x34, 0x0e, 0xb9, 0x34, 0x5a,
	0x73, 0x6f, 0xdd, 0x0b, 0x12, 0x7c, 0x70, 0x56, 0x76, 0x49, 0xf0, 0xdf, 0xbd, 0x02, 0x00, 0x00,
	0xff, 0xff, 0xa7, 0xf6, 0xb7, 0xb5, 0xcd, 0x00, 0x00, 0x00,
}
