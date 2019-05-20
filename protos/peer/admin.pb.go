


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
	return fileDescriptor_admin_b797b6f31dc49faf, []int{0, 0}
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
	return fileDescriptor_admin_b797b6f31dc49faf, []int{0}
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
	return fileDescriptor_admin_b797b6f31dc49faf, []int{1}
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
	return fileDescriptor_admin_b797b6f31dc49faf, []int{2}
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

type LogSpecRequest struct {
	LogSpec              string   `protobuf:"bytes,1,opt,name=log_spec,json=logSpec,proto3" json:"log_spec,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogSpecRequest) Reset()         { *m = LogSpecRequest{} }
func (m *LogSpecRequest) String() string { return proto.CompactTextString(m) }
func (*LogSpecRequest) ProtoMessage()    {}
func (*LogSpecRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_b797b6f31dc49faf, []int{3}
}
func (m *LogSpecRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogSpecRequest.Unmarshal(m, b)
}
func (m *LogSpecRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogSpecRequest.Marshal(b, m, deterministic)
}
func (dst *LogSpecRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogSpecRequest.Merge(dst, src)
}
func (m *LogSpecRequest) XXX_Size() int {
	return xxx_messageInfo_LogSpecRequest.Size(m)
}
func (m *LogSpecRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogSpecRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogSpecRequest proto.InternalMessageInfo

func (m *LogSpecRequest) GetLogSpec() string {
	if m != nil {
		return m.LogSpec
	}
	return ""
}

type LogSpecResponse struct {
	LogSpec              string   `protobuf:"bytes,1,opt,name=log_spec,json=logSpec,proto3" json:"log_spec,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogSpecResponse) Reset()         { *m = LogSpecResponse{} }
func (m *LogSpecResponse) String() string { return proto.CompactTextString(m) }
func (*LogSpecResponse) ProtoMessage()    {}
func (*LogSpecResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_b797b6f31dc49faf, []int{4}
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

func (m *LogSpecResponse) GetError() string {
	if m != nil {
		return m.Error
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
	return fileDescriptor_admin_b797b6f31dc49faf, []int{5}
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

type AdminOperation_LogSpecReq struct {
	LogSpecReq *LogSpecRequest `protobuf:"bytes,2,opt,name=logSpecReq,proto3,oneof"`
}

func (*AdminOperation_LogReq) isAdminOperation_Content() {}

func (*AdminOperation_LogSpecReq) isAdminOperation_Content() {}

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

func (m *AdminOperation) GetLogSpecReq() *LogSpecRequest {
	if x, ok := m.GetContent().(*AdminOperation_LogSpecReq); ok {
		return x.LogSpecReq
	}
	return nil
}


func (*AdminOperation) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _AdminOperation_OneofMarshaler, _AdminOperation_OneofUnmarshaler, _AdminOperation_OneofSizer, []interface{}{
		(*AdminOperation_LogReq)(nil),
		(*AdminOperation_LogSpecReq)(nil),
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
	case *AdminOperation_LogSpecReq:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LogSpecReq); err != nil {
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
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LogSpecRequest)
		err := b.DecodeMessage(msg)
		m.Content = &AdminOperation_LogSpecReq{msg}
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
	case *AdminOperation_LogSpecReq:
		s := proto.Size(x.LogSpecReq)
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
	proto.RegisterType((*LogSpecRequest)(nil), "protos.LogSpecRequest")
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
	SetLogSpec(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*LogSpecResponse, error)
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

func (c *adminClient) SetLogSpec(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (*LogSpecResponse, error) {
	out := new(LogSpecResponse)
	err := c.cc.Invoke(ctx, "/protos.Admin/SetLogSpec", in, out, opts...)
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
	SetLogSpec(context.Context, *common.Envelope) (*LogSpecResponse, error)
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

func _Admin_SetLogSpec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Envelope)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).SetLogSpec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/SetLogSpec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).SetLogSpec(ctx, req.(*common.Envelope))
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
		{
			MethodName: "SetLogSpec",
			Handler:    _Admin_SetLogSpec_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "peer/admin.proto",
}

func init() { proto.RegisterFile("peer/admin.proto", fileDescriptor_admin_b797b6f31dc49faf) }

var fileDescriptor_admin_b797b6f31dc49faf = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x6f, 0x6b, 0xd3, 0x40,
	0x18, 0x5f, 0x37, 0xdb, 0x9a, 0xa7, 0xb3, 0xcb, 0xce, 0xb1, 0xd5, 0x0e, 0x51, 0xf2, 0x4a, 0x11,
	0x12, 0xec, 0x90, 0x0a, 0xbe, 0x18, 0xad, 0x8d, 0x9d, 0xd8, 0xa5, 0x25, 0x59, 0x11, 0x05, 0x29,
	0x69, 0xfa, 0x2c, 0x2b, 0x5e, 0x73, 0xd9, 0xe5, 0x5a, 0xd8, 0x37, 0xf0, 0x73, 0x88, 0x1f, 0x54,
	0x72, 0x77, 0x61, 0xc5, 0x55, 0x10, 0xf7, 0xea, 0xee, 0x79, 0xee, 0xf7, 0xfb, 0x3d, 0xff, 0xee,
	0x0e, 0xcc, 0x14, 0x91, 0x3b, 0xe1, 0x6c, 0x31, 0x4f, 0xec, 0x94, 0x33, 0xc1, 0x48, 0x45, 0x2e,
	0x59, 0xf3, 0x38, 0x66, 0x2c, 0xa6, 0xe8, 0x48, 0x73, 0xba, 0xbc, 0x74, 0x70, 0x91, 0x8a, 0x1b,
	0x05, 0x6a, 0x3e, 0x8e, 0xd8, 0x62, 0xc1, 0x12, 0x47, 0x2d, 0xca, 0x69, 0xfd, 0x2c, 0xc1, 0x6e,
	0x80, 0x7c, 0x85, 0x3c, 0x10, 0xa1, 0x58, 0x66, 0xa4, 0x0d, 0x95, 0x4c, 0xee, 0x1a, 0xa5, 0xe7,
	0xa5, 0x17, 0xf5, 0xd6, 0x33, 0x05, 0xcc, 0xec, 0x75, 0x94, 0xad, 0x96, 0xf7, 0x6c, 0x86, 0xbe,
	0x86, 0x5b, 0x5f, 0x00, 0x6e, 0xbd, 0xe4, 0x11, 0x18, 0x63, 0xaf, 0xe7, 0x7e, 0xf8, 0xe8, 0xb9,
	0x3d, 0x73, 0x8b, 0xd4, 0xa0, 0x1a, 0x5c, 0x74, 0xfc, 0x0b, 0xb7, 0x67, 0x96, 0x94, 0x31, 0x1c,
	0x8d, 0xdc, 0x9e, 0xb9, 0x4d, 0x00, 0x2a, 0xa3, 0xce, 0x38, 0x70, 0x7b, 0xe6, 0x0e, 0x31, 0xa0,
	0xec, 0xfa, 0xfe, 0xd0, 0x37, 0x1f, 0xe4, 0x98, 0xb1, 0xf7, 0xc9, 0x1b, 0x7e, 0xf6, 0xcc, 0xb2,
	0x75, 0x0e, 0x7b, 0x03, 0x16, 0x0f, 0x70, 0x85, 0xd4, 0xc7, 0xeb, 0x25, 0x66, 0x82, 0x3c, 0x05,
	0xa0, 0x2c, 0x9e, 0x2c, 0xd8, 0x6c, 0x49, 0x51, 0xa6, 0x6a, 0xf8, 0x06, 0x65, 0xf1, 0xb9, 0x74,
	0x90, 0x63, 0xc8, 0x8d, 0x09, 0xcd, 0x29, 0x8d, 0x6d, 0x79, 0xfa, 0x90, 0x6a, 0x09, 0xcb, 0x03,
	0xf3, 0x56, 0x2e, 0x4b, 0x59, 0x92, 0xe1, 0xbd, 0xf4, 0x5e, 0x41, 0x7d, 0xc0, 0xe2, 0x20, 0xc5,
	0xa8, 0xc8, 0xee, 0x09, 0xe4, 0xa7, 0x93, 0x2c, 0xc5, 0x48, 0x6b, 0x55, 0xa9, 0x42, 0x58, 0x5d,
	0x59, 0x8b, 0x02, 0xeb, 0xd8, 0x7f, 0x47, 0x93, 0x03, 0x28, 0x23, 0xe7, 0x8c, 0xeb, 0x98, 0xca,
	0xb0, 0x7e, 0x94, 0xa0, 0xde, 0xc9, 0xc7, 0x3f, 0x4c, 0x91, 0x87, 0x62, 0xce, 0x12, 0xf2, 0x1a,
	0x2a, 0x94, 0xc5, 0x3e, 0x5e, 0x4b, 0x85, 0x5a, 0xeb, 0xa8, 0x18, 0xdb, 0x1f, 0x8d, 0x3b, 0xdb,
	0xf2, 0x35, 0x90, 0xbc, 0x95, 0x25, 0xeb, 0xb4, 0x65, 0x80, 0x5a, 0xeb, 0x70, 0x8d, 0xb6, 0x56,
	0xd0, 0xd9, 0x96, 0xbf, 0x86, 0xed, 0x1a, 0x50, 0x8d, 0x58, 0x22, 0x30, 0x11, 0xad, 0x5f, 0x3b,
	0x50, 0x96, 0xa9, 0x90, 0x13, 0x30, 0xfa, 0x28, 0xf4, 0x2d, 0x32, 0x6d, 0x7d, 0xcb, 0xdc, 0x64,
	0x85, 0x94, 0xa5, 0xd8, 0x3c, 0xd8, 0x74, 0x8f, 0xc8, 0x1b, 0xa8, 0x05, 0x22, 0xe4, 0x42, 0x39,
	0xff, 0x99, 0x76, 0x0a, 0xfb, 0x7d, 0x14, 0x6a, 0x36, 0x45, 0x81, 0x1b, 0xc8, 0x8d, 0xbb, 0x4d,
	0xd0, 0x2d, 0x3f, 0x85, 0xfd, 0xe0, 0x5e, 0x02, 0xef, 0x60, 0xcf, 0xc7, 0x15, 0x72, 0x51, 0x9c,
	0x6c, 0xaa, 0xf9, 0xd0, 0x56, 0xef, 0xd1, 0x2e, 0xde, 0xa3, 0xed, 0xe6, 0xef, 0x91, 0xb4, 0x01,
	0xfa, 0x28, 0x74, 0x8b, 0x37, 0xf0, 0x8e, 0xee, 0x4c, 0x41, 0x47, 0x6d, 0x03, 0x04, 0xff, 0x43,
	0xec, 0x7e, 0x03, 0x8b, 0xf1, 0xd8, 0xbe, 0xba, 0x49, 0x91, 0x53, 0x9c, 0xc5, 0xc8, 0xed, 0xcb,
	0x70, 0xca, 0xe7, 0x51, 0x41, 0xc8, 0xbf, 0x94, 0xee, 0xae, 0x9c, 0xe4, 0x28, 0x8c, 0xbe, 0x87,
	0x31, 0x7e, 0x7d, 0x19, 0xcf, 0xc5, 0xd5, 0x72, 0x9a, 0x07, 0x71, 0xd6, 0x88, 0x8e, 0x22, 0xaa,
	0x3f, 0x26, 0x73, 0x72, 0xe2, 0x54, 0xfd, 0x3f, 0x27, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x0e,
	0x63, 0x5f, 0xc4, 0x9a, 0x04, 0x00, 0x00,
}
