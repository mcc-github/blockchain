



package messages

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type Duration struct {
	
	
	Seconds int64 `protobuf:"varint,1,opt,name=seconds" json:"seconds,omitempty"`
	
	
	
	
	
	
	Nanos int32 `protobuf:"varint,2,opt,name=nanos" json:"nanos,omitempty"`
}

func (m *Duration) Reset()                    { *m = Duration{} }
func (m *Duration) String() string            { return proto.CompactTextString(m) }
func (*Duration) ProtoMessage()               {}
func (*Duration) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Duration) GetSeconds() int64 {
	if m != nil {
		return m.Seconds
	}
	return 0
}

func (m *Duration) GetNanos() int32 {
	if m != nil {
		return m.Nanos
	}
	return 0
}

type Timestamp struct {
	
	
	
	Seconds int64 `protobuf:"varint,1,opt,name=seconds" json:"seconds,omitempty"`
	
	
	
	
	Nanos int32 `protobuf:"varint,2,opt,name=nanos" json:"nanos,omitempty"`
}

func (m *Timestamp) Reset()                    { *m = Timestamp{} }
func (m *Timestamp) String() string            { return proto.CompactTextString(m) }
func (*Timestamp) ProtoMessage()               {}
func (*Timestamp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Timestamp) GetSeconds() int64 {
	if m != nil {
		return m.Seconds
	}
	return 0
}

func (m *Timestamp) GetNanos() int32 {
	if m != nil {
		return m.Nanos
	}
	return 0
}

type LoadBalanceRequest struct {
	
	
	
	LoadBalanceRequestType isLoadBalanceRequest_LoadBalanceRequestType `protobuf_oneof:"load_balance_request_type"`
}

func (m *LoadBalanceRequest) Reset()                    { *m = LoadBalanceRequest{} }
func (m *LoadBalanceRequest) String() string            { return proto.CompactTextString(m) }
func (*LoadBalanceRequest) ProtoMessage()               {}
func (*LoadBalanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isLoadBalanceRequest_LoadBalanceRequestType interface {
	isLoadBalanceRequest_LoadBalanceRequestType()
}

type LoadBalanceRequest_InitialRequest struct {
	InitialRequest *InitialLoadBalanceRequest `protobuf:"bytes,1,opt,name=initial_request,json=initialRequest,oneof"`
}
type LoadBalanceRequest_ClientStats struct {
	ClientStats *ClientStats `protobuf:"bytes,2,opt,name=client_stats,json=clientStats,oneof"`
}

func (*LoadBalanceRequest_InitialRequest) isLoadBalanceRequest_LoadBalanceRequestType() {}
func (*LoadBalanceRequest_ClientStats) isLoadBalanceRequest_LoadBalanceRequestType()    {}

func (m *LoadBalanceRequest) GetLoadBalanceRequestType() isLoadBalanceRequest_LoadBalanceRequestType {
	if m != nil {
		return m.LoadBalanceRequestType
	}
	return nil
}

func (m *LoadBalanceRequest) GetInitialRequest() *InitialLoadBalanceRequest {
	if x, ok := m.GetLoadBalanceRequestType().(*LoadBalanceRequest_InitialRequest); ok {
		return x.InitialRequest
	}
	return nil
}

func (m *LoadBalanceRequest) GetClientStats() *ClientStats {
	if x, ok := m.GetLoadBalanceRequestType().(*LoadBalanceRequest_ClientStats); ok {
		return x.ClientStats
	}
	return nil
}


func (*LoadBalanceRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LoadBalanceRequest_OneofMarshaler, _LoadBalanceRequest_OneofUnmarshaler, _LoadBalanceRequest_OneofSizer, []interface{}{
		(*LoadBalanceRequest_InitialRequest)(nil),
		(*LoadBalanceRequest_ClientStats)(nil),
	}
}

func _LoadBalanceRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LoadBalanceRequest)
	
	switch x := m.LoadBalanceRequestType.(type) {
	case *LoadBalanceRequest_InitialRequest:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.InitialRequest); err != nil {
			return err
		}
	case *LoadBalanceRequest_ClientStats:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ClientStats); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LoadBalanceRequest.LoadBalanceRequestType has unexpected type %T", x)
	}
	return nil
}

func _LoadBalanceRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LoadBalanceRequest)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(InitialLoadBalanceRequest)
		err := b.DecodeMessage(msg)
		m.LoadBalanceRequestType = &LoadBalanceRequest_InitialRequest{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ClientStats)
		err := b.DecodeMessage(msg)
		m.LoadBalanceRequestType = &LoadBalanceRequest_ClientStats{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LoadBalanceRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LoadBalanceRequest)
	
	switch x := m.LoadBalanceRequestType.(type) {
	case *LoadBalanceRequest_InitialRequest:
		s := proto.Size(x.InitialRequest)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LoadBalanceRequest_ClientStats:
		s := proto.Size(x.ClientStats)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type InitialLoadBalanceRequest struct {
	
	
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *InitialLoadBalanceRequest) Reset()                    { *m = InitialLoadBalanceRequest{} }
func (m *InitialLoadBalanceRequest) String() string            { return proto.CompactTextString(m) }
func (*InitialLoadBalanceRequest) ProtoMessage()               {}
func (*InitialLoadBalanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *InitialLoadBalanceRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}



type ClientStats struct {
	
	Timestamp *Timestamp `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	
	NumCallsStarted int64 `protobuf:"varint,2,opt,name=num_calls_started,json=numCallsStarted" json:"num_calls_started,omitempty"`
	
	NumCallsFinished int64 `protobuf:"varint,3,opt,name=num_calls_finished,json=numCallsFinished" json:"num_calls_finished,omitempty"`
	
	
	NumCallsFinishedWithDropForRateLimiting int64 `protobuf:"varint,4,opt,name=num_calls_finished_with_drop_for_rate_limiting,json=numCallsFinishedWithDropForRateLimiting" json:"num_calls_finished_with_drop_for_rate_limiting,omitempty"`
	
	
	NumCallsFinishedWithDropForLoadBalancing int64 `protobuf:"varint,5,opt,name=num_calls_finished_with_drop_for_load_balancing,json=numCallsFinishedWithDropForLoadBalancing" json:"num_calls_finished_with_drop_for_load_balancing,omitempty"`
	
	NumCallsFinishedWithClientFailedToSend int64 `protobuf:"varint,6,opt,name=num_calls_finished_with_client_failed_to_send,json=numCallsFinishedWithClientFailedToSend" json:"num_calls_finished_with_client_failed_to_send,omitempty"`
	
	
	NumCallsFinishedKnownReceived int64 `protobuf:"varint,7,opt,name=num_calls_finished_known_received,json=numCallsFinishedKnownReceived" json:"num_calls_finished_known_received,omitempty"`
}

func (m *ClientStats) Reset()                    { *m = ClientStats{} }
func (m *ClientStats) String() string            { return proto.CompactTextString(m) }
func (*ClientStats) ProtoMessage()               {}
func (*ClientStats) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ClientStats) GetTimestamp() *Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *ClientStats) GetNumCallsStarted() int64 {
	if m != nil {
		return m.NumCallsStarted
	}
	return 0
}

func (m *ClientStats) GetNumCallsFinished() int64 {
	if m != nil {
		return m.NumCallsFinished
	}
	return 0
}

func (m *ClientStats) GetNumCallsFinishedWithDropForRateLimiting() int64 {
	if m != nil {
		return m.NumCallsFinishedWithDropForRateLimiting
	}
	return 0
}

func (m *ClientStats) GetNumCallsFinishedWithDropForLoadBalancing() int64 {
	if m != nil {
		return m.NumCallsFinishedWithDropForLoadBalancing
	}
	return 0
}

func (m *ClientStats) GetNumCallsFinishedWithClientFailedToSend() int64 {
	if m != nil {
		return m.NumCallsFinishedWithClientFailedToSend
	}
	return 0
}

func (m *ClientStats) GetNumCallsFinishedKnownReceived() int64 {
	if m != nil {
		return m.NumCallsFinishedKnownReceived
	}
	return 0
}

type LoadBalanceResponse struct {
	
	
	
	LoadBalanceResponseType isLoadBalanceResponse_LoadBalanceResponseType `protobuf_oneof:"load_balance_response_type"`
}

func (m *LoadBalanceResponse) Reset()                    { *m = LoadBalanceResponse{} }
func (m *LoadBalanceResponse) String() string            { return proto.CompactTextString(m) }
func (*LoadBalanceResponse) ProtoMessage()               {}
func (*LoadBalanceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type isLoadBalanceResponse_LoadBalanceResponseType interface {
	isLoadBalanceResponse_LoadBalanceResponseType()
}

type LoadBalanceResponse_InitialResponse struct {
	InitialResponse *InitialLoadBalanceResponse `protobuf:"bytes,1,opt,name=initial_response,json=initialResponse,oneof"`
}
type LoadBalanceResponse_ServerList struct {
	ServerList *ServerList `protobuf:"bytes,2,opt,name=server_list,json=serverList,oneof"`
}

func (*LoadBalanceResponse_InitialResponse) isLoadBalanceResponse_LoadBalanceResponseType() {}
func (*LoadBalanceResponse_ServerList) isLoadBalanceResponse_LoadBalanceResponseType()      {}

func (m *LoadBalanceResponse) GetLoadBalanceResponseType() isLoadBalanceResponse_LoadBalanceResponseType {
	if m != nil {
		return m.LoadBalanceResponseType
	}
	return nil
}

func (m *LoadBalanceResponse) GetInitialResponse() *InitialLoadBalanceResponse {
	if x, ok := m.GetLoadBalanceResponseType().(*LoadBalanceResponse_InitialResponse); ok {
		return x.InitialResponse
	}
	return nil
}

func (m *LoadBalanceResponse) GetServerList() *ServerList {
	if x, ok := m.GetLoadBalanceResponseType().(*LoadBalanceResponse_ServerList); ok {
		return x.ServerList
	}
	return nil
}


func (*LoadBalanceResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LoadBalanceResponse_OneofMarshaler, _LoadBalanceResponse_OneofUnmarshaler, _LoadBalanceResponse_OneofSizer, []interface{}{
		(*LoadBalanceResponse_InitialResponse)(nil),
		(*LoadBalanceResponse_ServerList)(nil),
	}
}

func _LoadBalanceResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LoadBalanceResponse)
	
	switch x := m.LoadBalanceResponseType.(type) {
	case *LoadBalanceResponse_InitialResponse:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.InitialResponse); err != nil {
			return err
		}
	case *LoadBalanceResponse_ServerList:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ServerList); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LoadBalanceResponse.LoadBalanceResponseType has unexpected type %T", x)
	}
	return nil
}

func _LoadBalanceResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LoadBalanceResponse)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(InitialLoadBalanceResponse)
		err := b.DecodeMessage(msg)
		m.LoadBalanceResponseType = &LoadBalanceResponse_InitialResponse{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ServerList)
		err := b.DecodeMessage(msg)
		m.LoadBalanceResponseType = &LoadBalanceResponse_ServerList{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LoadBalanceResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LoadBalanceResponse)
	
	switch x := m.LoadBalanceResponseType.(type) {
	case *LoadBalanceResponse_InitialResponse:
		s := proto.Size(x.InitialResponse)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LoadBalanceResponse_ServerList:
		s := proto.Size(x.ServerList)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type InitialLoadBalanceResponse struct {
	
	
	
	
	
	LoadBalancerDelegate string `protobuf:"bytes,1,opt,name=load_balancer_delegate,json=loadBalancerDelegate" json:"load_balancer_delegate,omitempty"`
	
	
	
	ClientStatsReportInterval *Duration `protobuf:"bytes,2,opt,name=client_stats_report_interval,json=clientStatsReportInterval" json:"client_stats_report_interval,omitempty"`
}

func (m *InitialLoadBalanceResponse) Reset()                    { *m = InitialLoadBalanceResponse{} }
func (m *InitialLoadBalanceResponse) String() string            { return proto.CompactTextString(m) }
func (*InitialLoadBalanceResponse) ProtoMessage()               {}
func (*InitialLoadBalanceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *InitialLoadBalanceResponse) GetLoadBalancerDelegate() string {
	if m != nil {
		return m.LoadBalancerDelegate
	}
	return ""
}

func (m *InitialLoadBalanceResponse) GetClientStatsReportInterval() *Duration {
	if m != nil {
		return m.ClientStatsReportInterval
	}
	return nil
}

type ServerList struct {
	
	
	
	
	Servers []*Server `protobuf:"bytes,1,rep,name=servers" json:"servers,omitempty"`
}

func (m *ServerList) Reset()                    { *m = ServerList{} }
func (m *ServerList) String() string            { return proto.CompactTextString(m) }
func (*ServerList) ProtoMessage()               {}
func (*ServerList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ServerList) GetServers() []*Server {
	if m != nil {
		return m.Servers
	}
	return nil
}





type Server struct {
	
	
	IpAddress []byte `protobuf:"bytes,1,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	
	Port int32 `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
	
	
	
	
	
	
	LoadBalanceToken string `protobuf:"bytes,3,opt,name=load_balance_token,json=loadBalanceToken" json:"load_balance_token,omitempty"`
	
	
	DropForRateLimiting bool `protobuf:"varint,4,opt,name=drop_for_rate_limiting,json=dropForRateLimiting" json:"drop_for_rate_limiting,omitempty"`
	
	
	DropForLoadBalancing bool `protobuf:"varint,5,opt,name=drop_for_load_balancing,json=dropForLoadBalancing" json:"drop_for_load_balancing,omitempty"`
}

func (m *Server) Reset()                    { *m = Server{} }
func (m *Server) String() string            { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()               {}
func (*Server) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Server) GetIpAddress() []byte {
	if m != nil {
		return m.IpAddress
	}
	return nil
}

func (m *Server) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Server) GetLoadBalanceToken() string {
	if m != nil {
		return m.LoadBalanceToken
	}
	return ""
}

func (m *Server) GetDropForRateLimiting() bool {
	if m != nil {
		return m.DropForRateLimiting
	}
	return false
}

func (m *Server) GetDropForLoadBalancing() bool {
	if m != nil {
		return m.DropForLoadBalancing
	}
	return false
}

func init() {
	proto.RegisterType((*Duration)(nil), "grpc.lb.v1.Duration")
	proto.RegisterType((*Timestamp)(nil), "grpc.lb.v1.Timestamp")
	proto.RegisterType((*LoadBalanceRequest)(nil), "grpc.lb.v1.LoadBalanceRequest")
	proto.RegisterType((*InitialLoadBalanceRequest)(nil), "grpc.lb.v1.InitialLoadBalanceRequest")
	proto.RegisterType((*ClientStats)(nil), "grpc.lb.v1.ClientStats")
	proto.RegisterType((*LoadBalanceResponse)(nil), "grpc.lb.v1.LoadBalanceResponse")
	proto.RegisterType((*InitialLoadBalanceResponse)(nil), "grpc.lb.v1.InitialLoadBalanceResponse")
	proto.RegisterType((*ServerList)(nil), "grpc.lb.v1.ServerList")
	proto.RegisterType((*Server)(nil), "grpc.lb.v1.Server")
}

func init() { proto.RegisterFile("grpc_lb_v1/messages/messages.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xdd, 0x4e, 0x1b, 0x3b,
	0x10, 0x26, 0x27, 0x01, 0x92, 0x09, 0x3a, 0xe4, 0x98, 0x1c, 0x08, 0x14, 0x24, 0xba, 0x52, 0x69,
	0x54, 0xd1, 0x20, 0xa0, 0xbd, 0xe8, 0xcf, 0x45, 0x1b, 0x10, 0x0a, 0x2d, 0x17, 0x95, 0x43, 0x55,
	0xa9, 0x52, 0x65, 0x39, 0xd9, 0x21, 0x58, 0x6c, 0xec, 0xad, 0xed, 0x04, 0xf5, 0x11, 0xfa, 0x28,
	0x7d, 0x8c, 0xaa, 0xcf, 0xd0, 0xf7, 0xa9, 0xd6, 0xbb, 0x9b, 0x5d, 0x20, 0x80, 0x7a, 0x67, 0x8f,
	0xbf, 0xf9, 0xbe, 0xf1, 0xac, 0xbf, 0x59, 0xf0, 0x06, 0x3a, 0xec, 0xb3, 0xa0, 0xc7, 0xc6, 0xbb,
	0x3b, 0x43, 0x34, 0x86, 0x0f, 0xd0, 0x4c, 0x16, 0xad, 0x50, 0x2b, 0xab, 0x08, 0x44, 0x98, 0x56,
	0xd0, 0x6b, 0x8d, 0x77, 0xbd, 0x97, 0x50, 0x3e, 0x1c, 0x69, 0x6e, 0x85, 0x92, 0xa4, 0x01, 0xf3,
	0x06, 0xfb, 0x4a, 0xfa, 0xa6, 0x51, 0xd8, 0x2c, 0x34, 0x8b, 0x34, 0xdd, 0x92, 0x3a, 0xcc, 0x4a,
	0x2e, 0x95, 0x69, 0xfc, 0xb3, 0x59, 0x68, 0xce, 0xd2, 0x78, 0xe3, 0xbd, 0x82, 0xca, 0xa9, 0x18,
	0xa2, 0xb1, 0x7c, 0x18, 0xfe, 0x75, 0xf2, 0xcf, 0x02, 0x90, 0x13, 0xc5, 0xfd, 0x36, 0x0f, 0xb8,
	0xec, 0x23, 0xc5, 0xaf, 0x23, 0x34, 0x96, 0x7c, 0x80, 0x45, 0x21, 0x85, 0x15, 0x3c, 0x60, 0x3a,
	0x0e, 0x39, 0xba, 0xea, 0xde, 0xa3, 0x56, 0x56, 0x75, 0xeb, 0x38, 0x86, 0xdc, 0xcc, 0xef, 0xcc,
	0xd0, 0x7f, 0x93, 0xfc, 0x94, 0xf1, 0x35, 0x2c, 0xf4, 0x03, 0x81, 0xd2, 0x32, 0x63, 0xb9, 0x8d,
	0xab, 0xa8, 0xee, 0xad, 0xe4, 0xe9, 0x0e, 0xdc, 0x79, 0x37, 0x3a, 0xee, 0xcc, 0xd0, 0x6a, 0x3f,
	0xdb, 0xb6, 0x1f, 0xc0, 0x6a, 0xa0, 0xb8, 0xcf, 0x7a, 0xb1, 0x4c, 0x5a, 0x14, 0xb3, 0xdf, 0x42,
	0xf4, 0x76, 0x60, 0xf5, 0xd6, 0x4a, 0x08, 0x81, 0x92, 0xe4, 0x43, 0x74, 0xe5, 0x57, 0xa8, 0x5b,
	0x7b, 0xdf, 0x4b, 0x50, 0xcd, 0x89, 0x91, 0x7d, 0xa8, 0xd8, 0xb4, 0x83, 0xc9, 0x3d, 0xff, 0xcf,
	0x17, 0x36, 0x69, 0x2f, 0xcd, 0x70, 0xe4, 0x09, 0xfc, 0x27, 0x47, 0x43, 0xd6, 0xe7, 0x41, 0x60,
	0xa2, 0x3b, 0x69, 0x8b, 0xbe, 0xbb, 0x55, 0x91, 0x2e, 0xca, 0xd1, 0xf0, 0x20, 0x8a, 0x77, 0xe3,
	0x30, 0xd9, 0x06, 0x92, 0x61, 0xcf, 0x84, 0x14, 0xe6, 0x1c, 0xfd, 0x46, 0xd1, 0x81, 0x6b, 0x29,
	0xf8, 0x28, 0x89, 0x13, 0x06, 0xad, 0x9b, 0x68, 0x76, 0x29, 0xec, 0x39, 0xf3, 0xb5, 0x0a, 0xd9,
	0x99, 0xd2, 0x4c, 0x73, 0x8b, 0x2c, 0x10, 0x43, 0x61, 0x85, 0x1c, 0x34, 0x4a, 0x8e, 0xe9, 0xf1,
	0x75, 0xa6, 0x4f, 0xc2, 0x9e, 0x1f, 0x6a, 0x15, 0x1e, 0x29, 0x4d, 0xb9, 0xc5, 0x93, 0x04, 0x4e,
	0x38, 0xec, 0xdc, 0x2b, 0x90, 0x6b, 0x77, 0xa4, 0x30, 0xeb, 0x14, 0x9a, 0x77, 0x28, 0x64, 0xbd,
	0x8f, 0x24, 0xbe, 0xc0, 0xd3, 0xdb, 0x24, 0x92, 0x67, 0x70, 0xc6, 0x45, 0x80, 0x3e, 0xb3, 0x8a,
	0x19, 0x94, 0x7e, 0x63, 0xce, 0x09, 0x6c, 0x4d, 0x13, 0x88, 0x3f, 0xd5, 0x91, 0xc3, 0x9f, 0xaa,
	0x2e, 0x4a, 0x9f, 0x74, 0xe0, 0xe1, 0x14, 0xfa, 0x0b, 0xa9, 0x2e, 0x25, 0xd3, 0xd8, 0x47, 0x31,
	0x46, 0xbf, 0x31, 0xef, 0x28, 0x37, 0xae, 0x53, 0xbe, 0x8f, 0x50, 0x34, 0x01, 0x79, 0xbf, 0x0a,
	0xb0, 0x74, 0xe5, 0xd9, 0x98, 0x50, 0x49, 0x83, 0xa4, 0x0b, 0xb5, 0xcc, 0x01, 0x71, 0x2c, 0x79,
	0x1a, 0x5b, 0xf7, 0x59, 0x20, 0x46, 0x77, 0x66, 0xe8, 0xe2, 0xc4, 0x03, 0x09, 0xe9, 0x0b, 0xa8,
	0x1a, 0xd4, 0x63, 0xd4, 0x2c, 0x10, 0xc6, 0x26, 0x1e, 0x58, 0xce, 0xf3, 0x75, 0xdd, 0xf1, 0x89,
	0x70, 0x1e, 0x02, 0x33, 0xd9, 0xb5, 0xd7, 0x61, 0xed, 0x9a, 0x03, 0x62, 0xce, 0xd8, 0x02, 0x3f,
	0x0a, 0xb0, 0x76, 0x7b, 0x29, 0xe4, 0x19, 0x2c, 0xe7, 0x93, 0x35, 0xf3, 0x31, 0xc0, 0x01, 0xb7,
	0xa9, 0x2d, 0xea, 0x41, 0x96, 0xa4, 0x0f, 0x93, 0x33, 0xf2, 0x11, 0xd6, 0xf3, 0x96, 0x65, 0x1a,
	0x43, 0xa5, 0x2d, 0x13, 0xd2, 0xa2, 0x1e, 0xf3, 0x20, 0x29, 0xbf, 0x9e, 0x2f, 0x3f, 0x1d, 0x62,
	0x74, 0x35, 0xe7, 0x5e, 0xea, 0xf2, 0x8e, 0x93, 0x34, 0xef, 0x0d, 0x40, 0x76, 0x4b, 0xb2, 0x1d,
	0x0d, 0xac, 0x68, 0x17, 0x0d, 0xac, 0x62, 0xb3, 0xba, 0x47, 0x6e, 0xb6, 0x83, 0xa6, 0x90, 0x77,
	0xa5, 0x72, 0xb1, 0x56, 0xf2, 0x7e, 0x17, 0x60, 0x2e, 0x3e, 0x21, 0x1b, 0x00, 0x22, 0x64, 0xdc,
	0xf7, 0x35, 0x9a, 0x78, 0xe4, 0x2d, 0xd0, 0x8a, 0x08, 0xdf, 0xc6, 0x81, 0xc8, 0xfd, 0x91, 0x76,
	0x32, 0xf3, 0xdc, 0x3a, 0x32, 0xe3, 0x95, 0x4e, 0x5a, 0x75, 0x81, 0xd2, 0x99, 0xb1, 0x42, 0x6b,
	0xb9, 0x46, 0x9c, 0x46, 0x71, 0xb2, 0x0f, 0xcb, 0x77, 0x98, 0xae, 0x4c, 0x97, 0xfc, 0x29, 0x06,
	0x7b, 0x0e, 0x2b, 0x77, 0x19, 0xa9, 0x4c, 0xeb, 0xfe, 0x14, 0xd3, 0xb4, 0xe1, 0x73, 0x39, 0xfd,
	0x47, 0xf4, 0xe6, 0xdc, 0x4f, 0x62, 0xff, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa3, 0x36, 0x86,
	0xa6, 0x4a, 0x06, 0x00, 0x00,
}
