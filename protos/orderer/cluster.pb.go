


package orderer

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common "github.com/mcc-github/blockchain/protos/common"
	grpc "google.golang.org/grpc"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 


type StepRequest struct {
	
	
	
	Payload              isStepRequest_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *StepRequest) Reset()         { *m = StepRequest{} }
func (m *StepRequest) String() string { return proto.CompactTextString(m) }
func (*StepRequest) ProtoMessage()    {}
func (*StepRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e3b50707fd3a71f2, []int{0}
}

func (m *StepRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StepRequest.Unmarshal(m, b)
}
func (m *StepRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StepRequest.Marshal(b, m, deterministic)
}
func (m *StepRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StepRequest.Merge(m, src)
}
func (m *StepRequest) XXX_Size() int {
	return xxx_messageInfo_StepRequest.Size(m)
}
func (m *StepRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StepRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StepRequest proto.InternalMessageInfo

type isStepRequest_Payload interface {
	isStepRequest_Payload()
}

type StepRequest_ConsensusRequest struct {
	ConsensusRequest *ConsensusRequest `protobuf:"bytes,1,opt,name=consensus_request,json=consensusRequest,proto3,oneof"`
}

type StepRequest_SubmitRequest struct {
	SubmitRequest *SubmitRequest `protobuf:"bytes,2,opt,name=submit_request,json=submitRequest,proto3,oneof"`
}

func (*StepRequest_ConsensusRequest) isStepRequest_Payload() {}

func (*StepRequest_SubmitRequest) isStepRequest_Payload() {}

func (m *StepRequest) GetPayload() isStepRequest_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *StepRequest) GetConsensusRequest() *ConsensusRequest {
	if x, ok := m.GetPayload().(*StepRequest_ConsensusRequest); ok {
		return x.ConsensusRequest
	}
	return nil
}

func (m *StepRequest) GetSubmitRequest() *SubmitRequest {
	if x, ok := m.GetPayload().(*StepRequest_SubmitRequest); ok {
		return x.SubmitRequest
	}
	return nil
}


func (*StepRequest) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*StepRequest_ConsensusRequest)(nil),
		(*StepRequest_SubmitRequest)(nil),
	}
}


type StepResponse struct {
	
	
	Payload              isStepResponse_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *StepResponse) Reset()         { *m = StepResponse{} }
func (m *StepResponse) String() string { return proto.CompactTextString(m) }
func (*StepResponse) ProtoMessage()    {}
func (*StepResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e3b50707fd3a71f2, []int{1}
}

func (m *StepResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StepResponse.Unmarshal(m, b)
}
func (m *StepResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StepResponse.Marshal(b, m, deterministic)
}
func (m *StepResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StepResponse.Merge(m, src)
}
func (m *StepResponse) XXX_Size() int {
	return xxx_messageInfo_StepResponse.Size(m)
}
func (m *StepResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StepResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StepResponse proto.InternalMessageInfo

type isStepResponse_Payload interface {
	isStepResponse_Payload()
}

type StepResponse_SubmitRes struct {
	SubmitRes *SubmitResponse `protobuf:"bytes,1,opt,name=submit_res,json=submitRes,proto3,oneof"`
}

func (*StepResponse_SubmitRes) isStepResponse_Payload() {}

func (m *StepResponse) GetPayload() isStepResponse_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *StepResponse) GetSubmitRes() *SubmitResponse {
	if x, ok := m.GetPayload().(*StepResponse_SubmitRes); ok {
		return x.SubmitRes
	}
	return nil
}


func (*StepResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*StepResponse_SubmitRes)(nil),
	}
}


type ConsensusRequest struct {
	Channel              string   `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	Payload              []byte   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	Metadata             []byte   `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsensusRequest) Reset()         { *m = ConsensusRequest{} }
func (m *ConsensusRequest) String() string { return proto.CompactTextString(m) }
func (*ConsensusRequest) ProtoMessage()    {}
func (*ConsensusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e3b50707fd3a71f2, []int{2}
}

func (m *ConsensusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsensusRequest.Unmarshal(m, b)
}
func (m *ConsensusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsensusRequest.Marshal(b, m, deterministic)
}
func (m *ConsensusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusRequest.Merge(m, src)
}
func (m *ConsensusRequest) XXX_Size() int {
	return xxx_messageInfo_ConsensusRequest.Size(m)
}
func (m *ConsensusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusRequest proto.InternalMessageInfo

func (m *ConsensusRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *ConsensusRequest) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *ConsensusRequest) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}


type SubmitRequest struct {
	Channel string `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	
	
	
	LastValidationSeq uint64 `protobuf:"varint,2,opt,name=last_validation_seq,json=lastValidationSeq,proto3" json:"last_validation_seq,omitempty"`
	
	
	Payload              *common.Envelope `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SubmitRequest) Reset()         { *m = SubmitRequest{} }
func (m *SubmitRequest) String() string { return proto.CompactTextString(m) }
func (*SubmitRequest) ProtoMessage()    {}
func (*SubmitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e3b50707fd3a71f2, []int{3}
}

func (m *SubmitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubmitRequest.Unmarshal(m, b)
}
func (m *SubmitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubmitRequest.Marshal(b, m, deterministic)
}
func (m *SubmitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitRequest.Merge(m, src)
}
func (m *SubmitRequest) XXX_Size() int {
	return xxx_messageInfo_SubmitRequest.Size(m)
}
func (m *SubmitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitRequest proto.InternalMessageInfo

func (m *SubmitRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *SubmitRequest) GetLastValidationSeq() uint64 {
	if m != nil {
		return m.LastValidationSeq
	}
	return 0
}

func (m *SubmitRequest) GetPayload() *common.Envelope {
	if m != nil {
		return m.Payload
	}
	return nil
}



type SubmitResponse struct {
	Channel string `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	
	Status common.Status `protobuf:"varint,2,opt,name=status,proto3,enum=common.Status" json:"status,omitempty"`
	
	Info                 string   `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubmitResponse) Reset()         { *m = SubmitResponse{} }
func (m *SubmitResponse) String() string { return proto.CompactTextString(m) }
func (*SubmitResponse) ProtoMessage()    {}
func (*SubmitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e3b50707fd3a71f2, []int{4}
}

func (m *SubmitResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubmitResponse.Unmarshal(m, b)
}
func (m *SubmitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubmitResponse.Marshal(b, m, deterministic)
}
func (m *SubmitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitResponse.Merge(m, src)
}
func (m *SubmitResponse) XXX_Size() int {
	return xxx_messageInfo_SubmitResponse.Size(m)
}
func (m *SubmitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitResponse proto.InternalMessageInfo

func (m *SubmitResponse) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *SubmitResponse) GetStatus() common.Status {
	if m != nil {
		return m.Status
	}
	return common.Status_UNKNOWN
}

func (m *SubmitResponse) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func init() {
	proto.RegisterType((*StepRequest)(nil), "orderer.StepRequest")
	proto.RegisterType((*StepResponse)(nil), "orderer.StepResponse")
	proto.RegisterType((*ConsensusRequest)(nil), "orderer.ConsensusRequest")
	proto.RegisterType((*SubmitRequest)(nil), "orderer.SubmitRequest")
	proto.RegisterType((*SubmitResponse)(nil), "orderer.SubmitResponse")
}

func init() { proto.RegisterFile("orderer/cluster.proto", fileDescriptor_e3b50707fd3a71f2) }

var fileDescriptor_e3b50707fd3a71f2 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x5d, 0x58, 0xb5, 0x90, 0xbb, 0x2d, 0xea, 0x3c, 0x06, 0xa5, 0x4f, 0x28, 0x12, 0x68, 0x42,
	0x28, 0x41, 0xe5, 0x01, 0xde, 0x90, 0x3a, 0x21, 0xf5, 0xd9, 0x11, 0x3c, 0xf0, 0x52, 0x39, 0xc9,
	0x6d, 0x1b, 0x29, 0xb1, 0x53, 0xdb, 0x99, 0xb4, 0x0f, 0xe0, 0x4b, 0xf8, 0x51, 0x14, 0xdb, 0x49,
	0xba, 0x22, 0xf5, 0x29, 0xb9, 0xe7, 0x1c, 0x9f, 0x7b, 0xec, 0x7b, 0xe1, 0x4e, 0xc8, 0x02, 0x25,
	0xca, 0x24, 0xaf, 0x5a, 0xa5, 0x51, 0xc6, 0x8d, 0x14, 0x5a, 0x10, 0xdf, 0xc1, 0xf3, 0xdb, 0x5c,
	0xd4, 0xb5, 0xe0, 0x89, 0xfd, 0x58, 0x36, 0xfa, 0xeb, 0xc1, 0x65, 0xaa, 0xb1, 0xa1, 0xb8, 0x6f,
	0x51, 0x69, 0xb2, 0x82, 0x9b, 0x5c, 0x70, 0x85, 0x5c, 0xb5, 0x6a, 0x2d, 0x2d, 0x38, 0xf3, 0xde,
	0x79, 0xf7, 0x97, 0x8b, 0xb7, 0xb1, 0x73, 0x8a, 0x1f, 0x7a, 0x85, 0x3b, 0xb5, 0x3a, 0xa3, 0xd3,
	0xfc, 0x08, 0x23, 0xdf, 0x21, 0x54, 0x6d, 0x56, 0x97, 0x7a, 0xb0, 0x79, 0x61, 0x6c, 0x5e, 0x0f,
	0x36, 0xa9, 0xa1, 0x47, 0x8f, 0x6b, 0x75, 0x08, 0x2c, 0x03, 0xf0, 0x1b, 0xf6, 0x54, 0x09, 0x56,
	0x44, 0x29, 0x5c, 0xd9, 0x90, 0xaa, 0xe9, 0xda, 0x90, 0x6f, 0x00, 0x83, 0xb7, 0x72, 0xf1, 0xde,
	0xfc, 0xe7, 0x6b, 0xc5, 0xab, 0x33, 0x1a, 0xf4, 0xc6, 0xea, 0xd0, 0x34, 0x83, 0xe9, 0xf1, 0x45,
	0xc8, 0x0c, 0xfc, 0x7c, 0xc7, 0x38, 0xc7, 0xca, 0xb8, 0x06, 0xb4, 0x2f, 0x3b, 0xc6, 0x1d, 0x34,
	0xf7, 0xb8, 0xa2, 0x7d, 0x49, 0xe6, 0xf0, 0xb2, 0x46, 0xcd, 0x0a, 0xa6, 0xd9, 0xec, 0xdc, 0x50,
	0x43, 0x1d, 0xfd, 0xf1, 0xe0, 0xfa, 0xd9, 0x35, 0x4f, 0x74, 0x88, 0xe1, 0xb6, 0x62, 0x4a, 0xaf,
	0x1f, 0x59, 0x55, 0x16, 0x4c, 0x97, 0x82, 0xaf, 0x15, 0xee, 0x4d, 0xb7, 0x09, 0xbd, 0xe9, 0xa8,
	0x5f, 0x03, 0x93, 0xe2, 0x9e, 0x7c, 0x1c, 0x13, 0x9d, 0x9b, 0x17, 0x98, 0xc6, 0x6e, 0xb4, 0x3f,
	0xf8, 0x23, 0x56, 0xa2, 0xc1, 0x21, 0x63, 0xb4, 0x81, 0xf0, 0xf9, 0xab, 0x9c, 0xc8, 0xf1, 0x01,
	0x2e, 0x94, 0x66, 0xba, 0x55, 0xa6, 0x75, 0xb8, 0x08, 0x7b, 0xdb, 0xd4, 0xa0, 0xd4, 0xb1, 0x84,
	0xc0, 0xa4, 0xe4, 0x1b, 0x61, 0x9a, 0x07, 0xd4, 0xfc, 0x2f, 0x96, 0xe0, 0x3f, 0xd8, 0xed, 0x23,
	0x5f, 0x61, 0xd2, 0xcd, 0x8c, 0xbc, 0x1a, 0xe7, 0x32, 0xee, 0xd9, 0xfc, 0xee, 0x08, 0xb5, 0xa9,
	0xee, 0xbd, 0xcf, 0xde, 0xf2, 0x27, 0xbc, 0x17, 0x72, 0x1b, 0xef, 0x9e, 0x1a, 0x94, 0x15, 0x16,
	0x5b, 0x94, 0xf1, 0x86, 0x65, 0xb2, 0xcc, 0xed, 0xca, 0xaa, 0xfe, 0xe4, 0xef, 0x4f, 0xdb, 0x52,
	0xef, 0xda, 0xac, 0x8b, 0x97, 0x1c, 0xa8, 0x13, 0xab, 0x4e, 0xac, 0x3a, 0x71, 0xea, 0xec, 0xc2,
	0xd4, 0x5f, 0xfe, 0x05, 0x00, 0x00, 0xff, 0xff, 0x3d, 0x90, 0xa5, 0x9d, 0x27, 0x03, 0x00, 0x00,
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type ClusterClient interface {
	
	Step(ctx context.Context, opts ...grpc.CallOption) (Cluster_StepClient, error)
}

type clusterClient struct {
	cc *grpc.ClientConn
}

func NewClusterClient(cc *grpc.ClientConn) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) Step(ctx context.Context, opts ...grpc.CallOption) (Cluster_StepClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Cluster_serviceDesc.Streams[0], "/orderer.Cluster/Step", opts...)
	if err != nil {
		return nil, err
	}
	x := &clusterStepClient{stream}
	return x, nil
}

type Cluster_StepClient interface {
	Send(*StepRequest) error
	Recv() (*StepResponse, error)
	grpc.ClientStream
}

type clusterStepClient struct {
	grpc.ClientStream
}

func (x *clusterStepClient) Send(m *StepRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *clusterStepClient) Recv() (*StepResponse, error) {
	m := new(StepResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}


type ClusterServer interface {
	
	Step(Cluster_StepServer) error
}

func RegisterClusterServer(s *grpc.Server, srv ClusterServer) {
	s.RegisterService(&_Cluster_serviceDesc, srv)
}

func _Cluster_Step_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ClusterServer).Step(&clusterStepServer{stream})
}

type Cluster_StepServer interface {
	Send(*StepResponse) error
	Recv() (*StepRequest, error)
	grpc.ServerStream
}

type clusterStepServer struct {
	grpc.ServerStream
}

func (x *clusterStepServer) Send(m *StepResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *clusterStepServer) Recv() (*StepRequest, error) {
	m := new(StepRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Cluster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orderer.Cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Step",
			Handler:       _Cluster_Step_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "orderer/cluster.proto",
}
