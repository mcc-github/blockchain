


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import common "github.com/mcc-github/blockchain/protos/common"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type ServerStatus_StatusCode int32

const (
	ServerStatus_UNDEFINED ServerStatus_StatusCode = 0
	ServerStatus_STARTED   ServerStatus_StatusCode = 1
	ServerStatus_STOPPED   ServerStatus_StatusCode = 2
	ServerStatus_PAUSED    ServerStatus_StatusCode = 3
	ServerStatus_ERROR     ServerStatus_StatusCode = 4
	ServerStatus_UNKNOWN   ServerStatus_StatusCode = 5
)

var ServerStatus_StatusCode_name = map[int32]string{
	0: "UNDEFINED",
	1: "STARTED",
	2: "STOPPED",
	3: "PAUSED",
	4: "ERROR",
	5: "UNKNOWN",
}
var ServerStatus_StatusCode_value = map[string]int32{
	"UNDEFINED": 0,
	"STARTED":   1,
	"STOPPED":   2,
	"PAUSED":    3,
	"ERROR":     4,
	"UNKNOWN":   5,
}

func (x ServerStatus_StatusCode) String() string {
	return proto.EnumName(ServerStatus_StatusCode_name, int32(x))
}
func (ServerStatus_StatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_admin_8752d19060ce13e6, []int{0, 0}
}

type ServerStatus struct {
	Status               ServerStatus_StatusCode `protobuf:"varint,1,opt,name=status,proto3,enum=protos.ServerStatus_StatusCode" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ServerStatus) Reset()         { *m = ServerStatus{} }
func (m *ServerStatus) String() string { return proto.CompactTextString(m) }
func (*ServerStatus) ProtoMessage()    {}
func (*ServerStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_8752d19060ce13e6, []int{0}
}
func (m *ServerStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerStatus.Unmarshal(m, b)
}
func (m *ServerStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerStatus.Marshal(b, m, deterministic)
}
func (dst *ServerStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerStatus.Merge(dst, src)
}
func (m *ServerStatus) XXX_Size() int {
	return xxx_messageInfo_ServerStatus.Size(m)
}
func (m *ServerStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ServerStatus proto.InternalMessageInfo

func (m *ServerStatus) GetStatus() ServerStatus_StatusCode {
	if m != nil {
		return m.Status
	}
	return ServerStatus_UNDEFINED
}

type LogLevelRequest struct {
	LogModule            string   `protobuf:"bytes,1,opt,name=log_module,json=logModule,proto3" json:"log_module,omitempty"`
	LogLevel             string   `protobuf:"bytes,2,opt,name=log_level,json=logLevel,proto3" json:"log_level,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogLevelRequest) Reset()         { *m = LogLevelRequest{} }
func (m *LogLevelRequest) String() string { return proto.CompactTextString(m) }
func (*LogLevelRequest) ProtoMessage()    {}
func (*LogLevelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_8752d19060ce13e6, []int{1}
}
func (m *LogLevelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogLevelRequest.Unmarshal(m, b)
}
func (m *LogLevelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogLevelRequest.Marshal(b, m, deterministic)
}
func (dst *LogLevelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogLevelRequest.Merge(dst, src)
}
func (m *LogLevelRequest) XXX_Size() int {
	return xxx_messageInfo_LogLevelRequest.Size(m)
}
func (m *LogLevelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogLevelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogLevelRequest proto.InternalMessageInfo

func (m *LogLevelRequest) GetLogModule() string {
	if m != nil {
		return m.LogModule
	}
	return ""
}

func (m *LogLevelRequest) GetLogLevel() string {
	if m != nil {
		return m.LogLevel
	}
	return ""
}

type LogLevelResponse struct {
	LogModule            string   `protobuf:"bytes,1,opt,name=log_module,json=logModule,proto3" json:"log_module,omitempty"`
	LogLevel             string   `protobuf:"bytes,2,opt,name=log_level,json=logLevel,proto3" json:"log_level,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogLevelResponse) Reset()         { *m = LogLevelResponse{} }
func (m *LogLevelResponse) String() string { return proto.CompactTextString(m) }
func (*LogLevelResponse) ProtoMessage()    {}
func (*LogLevelResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_8752d19060ce13e6, []int{2}
}
func (m *LogLevelResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogLevelResponse.Unmarshal(m, b)
}
func (m *LogLevelResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogLevelResponse.Marshal(b, m, deterministic)
}
func (dst *LogLevelResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogLevelResponse.Merge(dst, src)
}
func (m *LogLevelResponse) XXX_Size() int {
	return xxx_messageInfo_LogLevelResponse.Size(m)
}
func (m *LogLevelResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogLevelResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogLevelResponse proto.InternalMessageInfo

func (m *LogLevelResponse) GetLogModule() string {
	if m != nil {
		return m.LogModule
	}
	return ""
}

func (m *LogLevelResponse) GetLogLevel() string {
	if m != nil {
		return m.LogLevel
	}
	return ""
}

type LogSpecResponse struct {
	LogSpec              string   `protobuf:"bytes,1,opt,name=log_spec,json=logSpec,proto3" json:"log_spec,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogSpecResponse) Reset()         { *m = LogSpecResponse{} }
func (m *LogSpecResponse) String() string { return proto.CompactTextString(m) }
func (*LogSpecResponse) ProtoMessage()    {}
func (*LogSpecResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_8752d19060ce13e6, []int{3}
}
func (m *LogSpecResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogSpecResponse.Unmarshal(m, b)
}
func (m *LogSpecResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogSpecResponse.Marshal(b, m, deterministic)
}
func (dst *LogSpecResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogSpecResponse.Merge(dst, src)
}
func (m *LogSpecResponse) XXX_Size() int {
	return xxx_messageInfo_LogSpecResponse.Size(m)
}
func (m *LogSpecResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogSpecResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogSpecResponse proto.InternalMessageInfo

func (m *LogSpecResponse) GetLogSpec() string {
	if m != nil {
		return m.LogSpec
	}
	return ""
}

type AdminOperation struct {
	
	
	Content              isAdminOperation_Content `protobuf_oneof:"content"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *AdminOperation) Reset()         { *m = AdminOperation{} }
func (m *AdminOperation) String() string { return proto.CompactTextString(m) }
func (*AdminOperation) ProtoMessage()    {}
func (*AdminOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_8752d19060ce13e6, []int{4}
}
func (m *AdminOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdminOperation.Unmarshal(m, b)
}
func (m *AdminOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdminOperation.Marshal(b, m, deterministic)
}
func (dst *AdminOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdminOperation.Merge(dst, src)
}
func (m *AdminOperation) XXX_Size() int {
	return xxx_messageInfo_AdminOperation.Size(m)
}
func (m *AdminOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_AdminOperation.DiscardUnknown(m)
}

var xxx_messageInfo_AdminOperation proto.InternalMessageInfo

type isAdminOperation_Content interface {
	isAdminOperation_Content()
}

type AdminOperation_LogReq struct {
	LogReq *LogLevelRequest `protobuf:"bytes,1,opt,name=logReq,proto3,oneof"`
}

func (*AdminOperation_LogReq) isAdminOperation_Content() {}

func (m *AdminOperation) GetContent() isAdminOperation_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *AdminOperation) GetLogReq() *LogLevelRequest {
	if x, ok := m.GetContent().(*AdminOperation_LogReq); ok {
		return x.LogReq
	}
	return nil
}


func (*AdminOperation) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _AdminOperation_OneofMarshaler, _AdminOperation_OneofUnmarshaler, _AdminOperation_OneofSizer, []interface{}{
		(*AdminOperation_LogReq)(nil),
	}
}

func _AdminOperation_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*AdminOperation)
	
	switch x := m.Content.(type) {
	case *AdminOperation_LogReq:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LogReq); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("AdminOperation.Content has unexpected type %T", x)
	}
	return nil
}

func _AdminOperation_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*AdminOperation)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LogLevelRequest)
		err := b.DecodeMessage(msg)
		m.Content = &AdminOperation_LogReq{msg}
		return true, err
	default:
		return false, nil
	}
}

func _AdminOperation_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*AdminOperation)
	
	switch x := m.Content.(type) {
	case *AdminOperation_LogReq:
		s := proto.Size(x.LogReq)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ServerStatus)(nil), "protos.ServerStatus")
	proto.RegisterType((*LogLevelRequest)(nil), "protos.LogLevelRequest")
	proto.RegisterType((*LogLevelResponse)(nil), "protos.LogLevelResponse")
	proto.RegisterType((*LogSpecResponse)(nil), "protos.LogSpecResponse")
	proto.RegisterType((*AdminOperation)(nil), "protos.AdminOperation")
	proto.RegisterEnum("protos.ServerStatus_StatusCode", ServerStatus_StatusCode_name, ServerStatus_StatusCode_value)
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type AdminClient interface {
	GetStatus(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*ServerStatus, error)
	StartServer(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*ServerStatus, error)
	GetModuleLogLevel(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*LogLevelResponse, error)
	SetModuleLogLevel(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*LogLevelResponse, error)
	RevertLogLevels(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*empty.Empty, error)
	GetLogSpec(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*LogSpecResponse, error)
}

type adminClient struct {
	cc *grpc.ClientConn
}

func NewAdminClient(cc *grpc.ClientConn) AdminClient {
	return &adminClient{cc}
}

func (c *adminClient) GetStatus(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*ServerStatus, error) {
	out := new(ServerStatus)
	err := c.cc.Invoke(ctx, "/protos.Admin/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) StartServer(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*ServerStatus, error) {
	out := new(ServerStatus)
	err := c.cc.Invoke(ctx, "/protos.Admin/StartServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) GetModuleLogLevel(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*LogLevelResponse, error) {
	out := new(LogLevelResponse)
	err := c.cc.Invoke(ctx, "/protos.Admin/GetModuleLogLevel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) SetModuleLogLevel(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*LogLevelResponse, error) {
	out := new(LogLevelResponse)
	err := c.cc.Invoke(ctx, "/protos.Admin/SetModuleLogLevel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) RevertLogLevels(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/protos.Admin/RevertLogLevels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) GetLogSpec(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*LogSpecResponse, error) {
	out := new(LogSpecResponse)
	err := c.cc.Invoke(ctx, "/protos.Admin/GetLogSpec", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}


type AdminServer interface {
	GetStatus(context.Context, *common.Envelope) (*ServerStatus, error)
	StartServer(context.Context, *common.Envelope) (*ServerStatus, error)
	GetModuleLogLevel(context.Context, *common.Envelope) (*LogLevelResponse, error)
	SetModuleLogLevel(context.Context, *common.Envelope) (*LogLevelResponse, error)
	RevertLogLevels(context.Context, *common.Envelope) (*empty.Empty, error)
	GetLogSpec(context.Context, *common.Envelope) (*LogSpecResponse, error)
}

func RegisterAdminServer(s *grpc.Server, srv AdminServer) {
	s.RegisterService(&_Admin_serviceDesc, srv)
}

func _Admin_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Envelope)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).GetStatus(ctx, req.(*common.Envelope))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_StartServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Envelope)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).StartServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/StartServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).StartServer(ctx, req.(*common.Envelope))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_GetModuleLogLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Envelope)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).GetModuleLogLevel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/GetModuleLogLevel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).GetModuleLogLevel(ctx, req.(*common.Envelope))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_SetModuleLogLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Envelope)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).SetModuleLogLevel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/SetModuleLogLevel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).SetModuleLogLevel(ctx, req.(*common.Envelope))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_RevertLogLevels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Envelope)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).RevertLogLevels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/RevertLogLevels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).RevertLogLevels(ctx, req.(*common.Envelope))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_GetLogSpec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Envelope)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).GetLogSpec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/GetLogSpec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).GetLogSpec(ctx, req.(*common.Envelope))
	}
	return interceptor(ctx, in, info, handler)
}

var _Admin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _Admin_GetStatus_Handler,
		},
		{
			MethodName: "StartServer",
			Handler:    _Admin_StartServer_Handler,
		},
		{
			MethodName: "GetModuleLogLevel",
			Handler:    _Admin_GetModuleLogLevel_Handler,
		},
		{
			MethodName: "SetModuleLogLevel",
			Handler:    _Admin_SetModuleLogLevel_Handler,
		},
		{
			MethodName: "RevertLogLevels",
			Handler:    _Admin_RevertLogLevels_Handler,
		},
		{
			MethodName: "GetLogSpec",
			Handler:    _Admin_GetLogSpec_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "peer/admin.proto",
}

func init() { proto.RegisterFile("peer/admin.proto", fileDescriptor_admin_8752d19060ce13e6) }

var fileDescriptor_admin_8752d19060ce13e6 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x6d, 0x07, 0x6d, 0xc9, 0xed, 0xd8, 0x82, 0x41, 0x50, 0x3a, 0x21, 0x50, 0x9e, 0x40, 0x42,
	0x89, 0x28, 0x42, 0x13, 0x0f, 0x3c, 0xb4, 0x34, 0x14, 0xc4, 0x96, 0x56, 0xce, 0x2a, 0x04, 0x12,
	0xaa, 0xd2, 0xf4, 0xce, 0xab, 0x70, 0xe2, 0xcc, 0x71, 0x2a, 0xed, 0x77, 0xf8, 0x49, 0x5e, 0x51,
	0xec, 0x84, 0x55, 0xd0, 0x17, 0xb4, 0x27, 0xc7, 0xd7, 0xe7, 0x9c, 0xdc, 0x7b, 0x8f, 0x0e, 0xd8,
	0x19, 0xa2, 0xf4, 0xa2, 0x55, 0xb2, 0x4e, 0xdd, 0x4c, 0x0a, 0x25, 0x48, 0x5b, 0x1f, 0x79, 0xff,
	0x88, 0x09, 0xc1, 0x38, 0x7a, 0xfa, 0xba, 0x2c, 0xce, 0x3d, 0x4c, 0x32, 0x75, 0x65, 0x40, 0xfd,
	0xfb, 0xb1, 0x48, 0x12, 0x91, 0x7a, 0xe6, 0x30, 0x45, 0xe7, 0x67, 0x13, 0xf6, 0x43, 0x94, 0x1b,
	0x94, 0xa1, 0x8a, 0x54, 0x91, 0x93, 0x63, 0x68, 0xe7, 0xfa, 0xab, 0xd7, 0x7c, 0xd6, 0x7c, 0x7e,
	0x30, 0x78, 0x6a, 0x80, 0xb9, 0xbb, 0x8d, 0x72, 0xcd, 0xf1, 0x5e, 0xac, 0x90, 0x56, 0x70, 0xe7,
	0x2b, 0xc0, 0x75, 0x95, 0xdc, 0x05, 0x6b, 0x1e, 0x8c, 0xfd, 0x0f, 0x9f, 0x02, 0x7f, 0x6c, 0x37,
	0x48, 0x17, 0x3a, 0xe1, 0xd9, 0x90, 0x9e, 0xf9, 0x63, 0xbb, 0x69, 0x2e, 0xd3, 0xd9, 0xcc, 0x1f,
	0xdb, 0x7b, 0x04, 0xa0, 0x3d, 0x1b, 0xce, 0x43, 0x7f, 0x6c, 0xdf, 0x22, 0x16, 0xb4, 0x7c, 0x4a,
	0xa7, 0xd4, 0xbe, 0x5d, 0x62, 0xe6, 0xc1, 0xe7, 0x60, 0xfa, 0x25, 0xb0, 0x5b, 0xce, 0x29, 0x1c,
	0x9e, 0x08, 0x76, 0x82, 0x1b, 0xe4, 0x14, 0x2f, 0x0b, 0xcc, 0x15, 0x79, 0x02, 0xc0, 0x05, 0x5b,
	0x24, 0x62, 0x55, 0x70, 0xd4, 0xad, 0x5a, 0xd4, 0xe2, 0x82, 0x9d, 0xea, 0x02, 0x39, 0x82, 0xf2,
	0xb2, 0xe0, 0x25, 0xa5, 0xb7, 0xa7, 0x5f, 0xef, 0xf0, 0x4a, 0xc2, 0x09, 0xc0, 0xbe, 0x96, 0xcb,
	0x33, 0x91, 0xe6, 0x78, 0x23, 0xbd, 0x97, 0xba, 0xbd, 0x30, 0xc3, 0xf8, 0x8f, 0xdc, 0x63, 0x28,
	0x9f, 0x17, 0x79, 0x86, 0x71, 0x25, 0xd6, 0xe1, 0x06, 0xe2, 0x04, 0x70, 0x30, 0x2c, 0xad, 0x9b,
	0x66, 0x28, 0x23, 0xb5, 0x16, 0x29, 0x79, 0x05, 0x6d, 0x2e, 0x18, 0xc5, 0x4b, 0x0d, 0xed, 0x0e,
	0x1e, 0xd5, 0x2b, 0xff, 0x6b, 0xe8, 0x8f, 0x0d, 0x5a, 0x01, 0x47, 0x16, 0x74, 0x62, 0x91, 0x2a,
	0x4c, 0xd5, 0xe0, 0xd7, 0x1e, 0xb4, 0xb4, 0x20, 0x79, 0x03, 0xd6, 0x04, 0x55, 0xe5, 0xa3, 0xed,
	0x56, 0x3e, 0xfb, 0xe9, 0x06, 0xb9, 0xc8, 0xb0, 0xff, 0x60, 0x97, 0x93, 0x4e, 0x83, 0x1c, 0x43,
	0x37, 0x54, 0x91, 0x54, 0xa6, 0xfc, 0x1f, 0xc4, 0x21, 0xdc, 0x9b, 0xa0, 0x32, 0x1b, 0xaa, 0x5b,
	0xdd, 0x41, 0xef, 0xfd, 0x3b, 0x8e, 0xd9, 0x92, 0x91, 0x08, 0x6f, 0x28, 0xf1, 0x0e, 0x0e, 0x29,
	0x6e, 0x50, 0xaa, 0xfa, 0x6d, 0xd7, 0xec, 0x0f, 0x5d, 0x93, 0x0c, 0xb7, 0x4e, 0x86, 0xeb, 0x97,
	0xc9, 0x70, 0x1a, 0xe4, 0x2d, 0xc0, 0x04, 0x55, 0xe5, 0xdf, 0x0e, 0xe6, 0xb6, 0x19, 0xdb, 0x16,
	0x3b, 0x8d, 0xd1, 0x77, 0x70, 0x84, 0x64, 0xee, 0xc5, 0x55, 0x86, 0x92, 0xe3, 0x8a, 0xa1, 0x74,
	0xcf, 0xa3, 0xa5, 0x5c, 0xc7, 0x35, 0xa5, 0xcc, 0xe9, 0x68, 0x5f, 0x9b, 0x33, 0x8b, 0xe2, 0x1f,
	0x11, 0xc3, 0x6f, 0x2f, 0xd8, 0x5a, 0x5d, 0x14, 0xcb, 0xf2, 0x37, 0xde, 0x16, 0xd1, 0x33, 0x44,
	0x13, 0xdc, 0xdc, 0x2b, 0x89, 0x4b, 0x13, 0xea, 0xd7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xaa,
	0xa9, 0xb7, 0x35, 0xef, 0x03, 0x00, 0x00,
}
