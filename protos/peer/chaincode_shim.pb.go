


package peer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ChaincodeMessage_Type int32

const (
	ChaincodeMessage_UNDEFINED           ChaincodeMessage_Type = 0
	ChaincodeMessage_REGISTER            ChaincodeMessage_Type = 1
	ChaincodeMessage_REGISTERED          ChaincodeMessage_Type = 2
	ChaincodeMessage_INIT                ChaincodeMessage_Type = 3
	ChaincodeMessage_READY               ChaincodeMessage_Type = 4
	ChaincodeMessage_TRANSACTION         ChaincodeMessage_Type = 5
	ChaincodeMessage_COMPLETED           ChaincodeMessage_Type = 6
	ChaincodeMessage_ERROR               ChaincodeMessage_Type = 7
	ChaincodeMessage_GET_STATE           ChaincodeMessage_Type = 8
	ChaincodeMessage_PUT_STATE           ChaincodeMessage_Type = 9
	ChaincodeMessage_DEL_STATE           ChaincodeMessage_Type = 10
	ChaincodeMessage_INVOKE_CHAINCODE    ChaincodeMessage_Type = 11
	ChaincodeMessage_RESPONSE            ChaincodeMessage_Type = 13
	ChaincodeMessage_GET_STATE_BY_RANGE  ChaincodeMessage_Type = 14
	ChaincodeMessage_GET_QUERY_RESULT    ChaincodeMessage_Type = 15
	ChaincodeMessage_QUERY_STATE_NEXT    ChaincodeMessage_Type = 16
	ChaincodeMessage_QUERY_STATE_CLOSE   ChaincodeMessage_Type = 17
	ChaincodeMessage_KEEPALIVE           ChaincodeMessage_Type = 18
	ChaincodeMessage_GET_HISTORY_FOR_KEY ChaincodeMessage_Type = 19
)

var ChaincodeMessage_Type_name = map[int32]string{
	0:  "UNDEFINED",
	1:  "REGISTER",
	2:  "REGISTERED",
	3:  "INIT",
	4:  "READY",
	5:  "TRANSACTION",
	6:  "COMPLETED",
	7:  "ERROR",
	8:  "GET_STATE",
	9:  "PUT_STATE",
	10: "DEL_STATE",
	11: "INVOKE_CHAINCODE",
	13: "RESPONSE",
	14: "GET_STATE_BY_RANGE",
	15: "GET_QUERY_RESULT",
	16: "QUERY_STATE_NEXT",
	17: "QUERY_STATE_CLOSE",
	18: "KEEPALIVE",
	19: "GET_HISTORY_FOR_KEY",
}
var ChaincodeMessage_Type_value = map[string]int32{
	"UNDEFINED":           0,
	"REGISTER":            1,
	"REGISTERED":          2,
	"INIT":                3,
	"READY":               4,
	"TRANSACTION":         5,
	"COMPLETED":           6,
	"ERROR":               7,
	"GET_STATE":           8,
	"PUT_STATE":           9,
	"DEL_STATE":           10,
	"INVOKE_CHAINCODE":    11,
	"RESPONSE":            13,
	"GET_STATE_BY_RANGE":  14,
	"GET_QUERY_RESULT":    15,
	"QUERY_STATE_NEXT":    16,
	"QUERY_STATE_CLOSE":   17,
	"KEEPALIVE":           18,
	"GET_HISTORY_FOR_KEY": 19,
}

func (x ChaincodeMessage_Type) String() string {
	return proto.EnumName(ChaincodeMessage_Type_name, int32(x))
}
func (ChaincodeMessage_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{0, 0} }

type ChaincodeMessage struct {
	Type      ChaincodeMessage_Type       `protobuf:"varint,1,opt,name=type,enum=protos.ChaincodeMessage_Type" json:"type,omitempty"`
	Timestamp *google_protobuf1.Timestamp `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp,omitempty"`
	Payload   []byte                      `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	Txid      string                      `protobuf:"bytes,4,opt,name=txid" json:"txid,omitempty"`
	Proposal  *SignedProposal             `protobuf:"bytes,5,opt,name=proposal" json:"proposal,omitempty"`
	
	
	
	ChaincodeEvent *ChaincodeEvent `protobuf:"bytes,6,opt,name=chaincode_event,json=chaincodeEvent" json:"chaincode_event,omitempty"`
	
	ChannelId string `protobuf:"bytes,7,opt,name=channel_id,json=channelId" json:"channel_id,omitempty"`
}

func (m *ChaincodeMessage) Reset()                    { *m = ChaincodeMessage{} }
func (m *ChaincodeMessage) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeMessage) ProtoMessage()               {}
func (*ChaincodeMessage) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *ChaincodeMessage) GetType() ChaincodeMessage_Type {
	if m != nil {
		return m.Type
	}
	return ChaincodeMessage_UNDEFINED
}

func (m *ChaincodeMessage) GetTimestamp() *google_protobuf1.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *ChaincodeMessage) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *ChaincodeMessage) GetTxid() string {
	if m != nil {
		return m.Txid
	}
	return ""
}

func (m *ChaincodeMessage) GetProposal() *SignedProposal {
	if m != nil {
		return m.Proposal
	}
	return nil
}

func (m *ChaincodeMessage) GetChaincodeEvent() *ChaincodeEvent {
	if m != nil {
		return m.ChaincodeEvent
	}
	return nil
}

func (m *ChaincodeMessage) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

type GetState struct {
	Key        string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Collection string `protobuf:"bytes,2,opt,name=collection" json:"collection,omitempty"`
}

func (m *GetState) Reset()                    { *m = GetState{} }
func (m *GetState) String() string            { return proto.CompactTextString(m) }
func (*GetState) ProtoMessage()               {}
func (*GetState) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *GetState) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *GetState) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

type PutState struct {
	Key        string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value      []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Collection string `protobuf:"bytes,3,opt,name=collection" json:"collection,omitempty"`
}

func (m *PutState) Reset()                    { *m = PutState{} }
func (m *PutState) String() string            { return proto.CompactTextString(m) }
func (*PutState) ProtoMessage()               {}
func (*PutState) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

func (m *PutState) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *PutState) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *PutState) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

type DelState struct {
	Key        string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Collection string `protobuf:"bytes,2,opt,name=collection" json:"collection,omitempty"`
}

func (m *DelState) Reset()                    { *m = DelState{} }
func (m *DelState) String() string            { return proto.CompactTextString(m) }
func (*DelState) ProtoMessage()               {}
func (*DelState) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *DelState) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *DelState) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

type GetStateByRange struct {
	StartKey   string `protobuf:"bytes,1,opt,name=startKey" json:"startKey,omitempty"`
	EndKey     string `protobuf:"bytes,2,opt,name=endKey" json:"endKey,omitempty"`
	Collection string `protobuf:"bytes,3,opt,name=collection" json:"collection,omitempty"`
}

func (m *GetStateByRange) Reset()                    { *m = GetStateByRange{} }
func (m *GetStateByRange) String() string            { return proto.CompactTextString(m) }
func (*GetStateByRange) ProtoMessage()               {}
func (*GetStateByRange) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *GetStateByRange) GetStartKey() string {
	if m != nil {
		return m.StartKey
	}
	return ""
}

func (m *GetStateByRange) GetEndKey() string {
	if m != nil {
		return m.EndKey
	}
	return ""
}

func (m *GetStateByRange) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

type GetQueryResult struct {
	Query      string `protobuf:"bytes,1,opt,name=query" json:"query,omitempty"`
	Collection string `protobuf:"bytes,2,opt,name=collection" json:"collection,omitempty"`
}

func (m *GetQueryResult) Reset()                    { *m = GetQueryResult{} }
func (m *GetQueryResult) String() string            { return proto.CompactTextString(m) }
func (*GetQueryResult) ProtoMessage()               {}
func (*GetQueryResult) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *GetQueryResult) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

func (m *GetQueryResult) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

type GetHistoryForKey struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *GetHistoryForKey) Reset()                    { *m = GetHistoryForKey{} }
func (m *GetHistoryForKey) String() string            { return proto.CompactTextString(m) }
func (*GetHistoryForKey) ProtoMessage()               {}
func (*GetHistoryForKey) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

func (m *GetHistoryForKey) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type QueryStateNext struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *QueryStateNext) Reset()                    { *m = QueryStateNext{} }
func (m *QueryStateNext) String() string            { return proto.CompactTextString(m) }
func (*QueryStateNext) ProtoMessage()               {}
func (*QueryStateNext) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{7} }

func (m *QueryStateNext) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type QueryStateClose struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *QueryStateClose) Reset()                    { *m = QueryStateClose{} }
func (m *QueryStateClose) String() string            { return proto.CompactTextString(m) }
func (*QueryStateClose) ProtoMessage()               {}
func (*QueryStateClose) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{8} }

func (m *QueryStateClose) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type QueryResultBytes struct {
	ResultBytes []byte `protobuf:"bytes,1,opt,name=resultBytes,proto3" json:"resultBytes,omitempty"`
}

func (m *QueryResultBytes) Reset()                    { *m = QueryResultBytes{} }
func (m *QueryResultBytes) String() string            { return proto.CompactTextString(m) }
func (*QueryResultBytes) ProtoMessage()               {}
func (*QueryResultBytes) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{9} }

func (m *QueryResultBytes) GetResultBytes() []byte {
	if m != nil {
		return m.ResultBytes
	}
	return nil
}

type QueryResponse struct {
	Results []*QueryResultBytes `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
	HasMore bool                `protobuf:"varint,2,opt,name=has_more,json=hasMore" json:"has_more,omitempty"`
	Id      string              `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *QueryResponse) Reset()                    { *m = QueryResponse{} }
func (m *QueryResponse) String() string            { return proto.CompactTextString(m) }
func (*QueryResponse) ProtoMessage()               {}
func (*QueryResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{10} }

func (m *QueryResponse) GetResults() []*QueryResultBytes {
	if m != nil {
		return m.Results
	}
	return nil
}

func (m *QueryResponse) GetHasMore() bool {
	if m != nil {
		return m.HasMore
	}
	return false
}

func (m *QueryResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*ChaincodeMessage)(nil), "protos.ChaincodeMessage")
	proto.RegisterType((*GetState)(nil), "protos.GetState")
	proto.RegisterType((*PutState)(nil), "protos.PutState")
	proto.RegisterType((*DelState)(nil), "protos.DelState")
	proto.RegisterType((*GetStateByRange)(nil), "protos.GetStateByRange")
	proto.RegisterType((*GetQueryResult)(nil), "protos.GetQueryResult")
	proto.RegisterType((*GetHistoryForKey)(nil), "protos.GetHistoryForKey")
	proto.RegisterType((*QueryStateNext)(nil), "protos.QueryStateNext")
	proto.RegisterType((*QueryStateClose)(nil), "protos.QueryStateClose")
	proto.RegisterType((*QueryResultBytes)(nil), "protos.QueryResultBytes")
	proto.RegisterType((*QueryResponse)(nil), "protos.QueryResponse")
	proto.RegisterEnum("protos.ChaincodeMessage_Type", ChaincodeMessage_Type_name, ChaincodeMessage_Type_value)
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4



type ChaincodeSupportClient interface {
	Register(ctx context.Context, opts ...grpc.CallOption) (ChaincodeSupport_RegisterClient, error)
}

type chaincodeSupportClient struct {
	cc *grpc.ClientConn
}

func NewChaincodeSupportClient(cc *grpc.ClientConn) ChaincodeSupportClient {
	return &chaincodeSupportClient{cc}
}

func (c *chaincodeSupportClient) Register(ctx context.Context, opts ...grpc.CallOption) (ChaincodeSupport_RegisterClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ChaincodeSupport_serviceDesc.Streams[0], c.cc, "/protos.ChaincodeSupport/Register", opts...)
	if err != nil {
		return nil, err
	}
	x := &chaincodeSupportRegisterClient{stream}
	return x, nil
}

type ChaincodeSupport_RegisterClient interface {
	Send(*ChaincodeMessage) error
	Recv() (*ChaincodeMessage, error)
	grpc.ClientStream
}

type chaincodeSupportRegisterClient struct {
	grpc.ClientStream
}

func (x *chaincodeSupportRegisterClient) Send(m *ChaincodeMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chaincodeSupportRegisterClient) Recv() (*ChaincodeMessage, error) {
	m := new(ChaincodeMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}



type ChaincodeSupportServer interface {
	Register(ChaincodeSupport_RegisterServer) error
}

func RegisterChaincodeSupportServer(s *grpc.Server, srv ChaincodeSupportServer) {
	s.RegisterService(&_ChaincodeSupport_serviceDesc, srv)
}

func _ChaincodeSupport_Register_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChaincodeSupportServer).Register(&chaincodeSupportRegisterServer{stream})
}

type ChaincodeSupport_RegisterServer interface {
	Send(*ChaincodeMessage) error
	Recv() (*ChaincodeMessage, error)
	grpc.ServerStream
}

type chaincodeSupportRegisterServer struct {
	grpc.ServerStream
}

func (x *chaincodeSupportRegisterServer) Send(m *ChaincodeMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chaincodeSupportRegisterServer) Recv() (*ChaincodeMessage, error) {
	m := new(ChaincodeMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ChaincodeSupport_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.ChaincodeSupport",
	HandlerType: (*ChaincodeSupportServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Register",
			Handler:       _ChaincodeSupport_Register_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "peer/chaincode_shim.proto",
}

func init() { proto.RegisterFile("peer/chaincode_shim.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x95, 0x51, 0x6f, 0xe2, 0x46,
	0x10, 0xc7, 0x8f, 0x00, 0xc1, 0x0c, 0x09, 0xec, 0x6d, 0xae, 0xa9, 0x0f, 0xe9, 0x5a, 0x8a, 0xfa,
	0x40, 0x5f, 0xa0, 0xa5, 0x7d, 0xe8, 0xc3, 0x49, 0x15, 0xc1, 0x1b, 0x62, 0x85, 0xd8, 0xdc, 0xda,
	0x39, 0x1d, 0x7d, 0xb1, 0x1c, 0xbc, 0x67, 0xac, 0x1a, 0xaf, 0x6b, 0x2f, 0xa7, 0xf3, 0x67, 0xe8,
	0x07, 0xeb, 0xd7, 0xaa, 0xd6, 0xc6, 0x84, 0x23, 0x8a, 0x4e, 0xea, 0x13, 0xfe, 0xcf, 0xfc, 0x66,
	0xe6, 0x3f, 0xd6, 0xb2, 0x86, 0xd7, 0x31, 0x63, 0xc9, 0x68, 0xb5, 0x76, 0x83, 0x68, 0xc5, 0x3d,
	0xe6, 0xa4, 0xeb, 0x60, 0x33, 0x8c, 0x13, 0x2e, 0x38, 0x3e, 0xcd, 0x7f, 0xd2, 0x6e, 0xf7, 0x08,
	0x61, 0x9f, 0x58, 0x24, 0x0a, 0xa6, 0x7b, 0x91, 0xe7, 0xe2, 0x84, 0xc7, 0x3c, 0x75, 0xc3, 0x5d,
	0xf0, 0x7b, 0x9f, 0x73, 0x3f, 0x64, 0xa3, 0x5c, 0x3d, 0x6c, 0x3f, 0x8e, 0x44, 0xb0, 0x61, 0xa9,
	0x70, 0x37, 0x71, 0x01, 0xf4, 0xff, 0xa9, 0x03, 0x9a, 0x96, 0xfd, 0xee, 0x58, 0x9a, 0xba, 0x3e,
	0xc3, 0xbf, 0x40, 0x4d, 0x64, 0x31, 0x53, 0x2b, 0xbd, 0xca, 0xa0, 0x3d, 0x7e, 0x53, 0xa0, 0xe9,
	0xf0, 0x98, 0x1b, 0xda, 0x59, 0xcc, 0x68, 0x8e, 0xe2, 0xdf, 0xa1, 0xb9, 0x6f, 0xad, 0x9e, 0xf4,
	0x2a, 0x83, 0xd6, 0xb8, 0x3b, 0x2c, 0x86, 0x0f, 0xcb, 0xe1, 0x43, 0xbb, 0x24, 0xe8, 0x23, 0x8c,
	0x55, 0x68, 0xc4, 0x6e, 0x16, 0x72, 0xd7, 0x53, 0xab, 0xbd, 0xca, 0xe0, 0x8c, 0x96, 0x12, 0x63,
	0xa8, 0x89, 0xcf, 0x81, 0xa7, 0xd6, 0x7a, 0x95, 0x41, 0x93, 0xe6, 0xcf, 0x78, 0x0c, 0x4a, 0xb9,
	0xa2, 0x5a, 0xcf, 0xc7, 0x5c, 0x96, 0xf6, 0xac, 0xc0, 0x8f, 0x98, 0xb7, 0xd8, 0x65, 0xe9, 0x9e,
	0xc3, 0x7f, 0x40, 0xe7, 0xe8, 0x95, 0xa9, 0xa7, 0x5f, 0x96, 0xee, 0x37, 0x23, 0x32, 0x4b, 0xdb,
	0xab, 0x2f, 0x34, 0x7e, 0x03, 0xb0, 0x5a, 0xbb, 0x51, 0xc4, 0x42, 0x27, 0xf0, 0xd4, 0x46, 0x6e,
	0xa7, 0xb9, 0x8b, 0xe8, 0x5e, 0xff, 0xdf, 0x13, 0xa8, 0xc9, 0x57, 0x81, 0xcf, 0xa1, 0x79, 0x6f,
	0x68, 0xe4, 0x5a, 0x37, 0x88, 0x86, 0x5e, 0xe0, 0x33, 0x50, 0x28, 0x99, 0xe9, 0x96, 0x4d, 0x28,
	0xaa, 0xe0, 0x36, 0x40, 0xa9, 0x88, 0x86, 0x4e, 0xb0, 0x02, 0x35, 0xdd, 0xd0, 0x6d, 0x54, 0xc5,
	0x4d, 0xa8, 0x53, 0x32, 0xd1, 0x96, 0xa8, 0x86, 0x3b, 0xd0, 0xb2, 0xe9, 0xc4, 0xb0, 0x26, 0x53,
	0x5b, 0x37, 0x0d, 0x54, 0x97, 0x2d, 0xa7, 0xe6, 0xdd, 0x62, 0x4e, 0x6c, 0xa2, 0xa1, 0x53, 0x89,
	0x12, 0x4a, 0x4d, 0x8a, 0x1a, 0x32, 0x33, 0x23, 0xb6, 0x63, 0xd9, 0x13, 0x9b, 0x20, 0x45, 0xca,
	0xc5, 0x7d, 0x29, 0x9b, 0x52, 0x6a, 0x64, 0xbe, 0x93, 0x80, 0x5f, 0x01, 0xd2, 0x8d, 0xf7, 0xe6,
	0x2d, 0x71, 0xa6, 0x37, 0x13, 0xdd, 0x98, 0x9a, 0x1a, 0x41, 0xad, 0xc2, 0xa0, 0xb5, 0x30, 0x0d,
	0x8b, 0xa0, 0x73, 0x7c, 0x09, 0x78, 0xdf, 0xd0, 0xb9, 0x5a, 0x3a, 0x74, 0x62, 0xcc, 0x08, 0x6a,
	0xcb, 0x5a, 0x19, 0x7f, 0x77, 0x4f, 0xe8, 0xd2, 0xa1, 0xc4, 0xba, 0x9f, 0xdb, 0xa8, 0x23, 0xa3,
	0x45, 0xa4, 0xe0, 0x0d, 0xf2, 0xc1, 0x46, 0x08, 0x7f, 0x03, 0x2f, 0x0f, 0xa3, 0xd3, 0xb9, 0x69,
	0x11, 0xf4, 0x52, 0xba, 0xb9, 0x25, 0x64, 0x31, 0x99, 0xeb, 0xef, 0x09, 0xc2, 0xf8, 0x5b, 0xb8,
	0x90, 0x1d, 0x6f, 0x74, 0xcb, 0x36, 0xe9, 0xd2, 0xb9, 0x36, 0xa9, 0x73, 0x4b, 0x96, 0xe8, 0xa2,
	0xff, 0x16, 0x94, 0x19, 0x13, 0x96, 0x70, 0x05, 0xc3, 0x08, 0xaa, 0x7f, 0xb1, 0x2c, 0x3f, 0x83,
	0x4d, 0x2a, 0x1f, 0xf1, 0x77, 0x00, 0x2b, 0x1e, 0x86, 0x6c, 0x25, 0x02, 0x1e, 0xe5, 0x87, 0xac,
	0x49, 0x0f, 0x22, 0x7d, 0x0a, 0xca, 0x62, 0xfb, 0x6c, 0xf5, 0x2b, 0xa8, 0x7f, 0x72, 0xc3, 0x2d,
	0xcb, 0x0b, 0xcf, 0x68, 0x21, 0x8e, 0x7a, 0x56, 0x9f, 0xf4, 0x7c, 0x0b, 0x8a, 0xc6, 0xc2, 0xff,
	0xeb, 0x88, 0x41, 0xa7, 0xdc, 0xe7, 0x2a, 0xa3, 0x6e, 0xe4, 0x33, 0xdc, 0x05, 0x25, 0x15, 0x6e,
	0x22, 0x6e, 0xf7, 0x9d, 0xf6, 0x1a, 0x5f, 0xc2, 0x29, 0x8b, 0x3c, 0x99, 0x29, 0x5a, 0xed, 0xd4,
	0x57, 0x4d, 0x5e, 0x43, 0x7b, 0xc6, 0xc4, 0xbb, 0x2d, 0x4b, 0x32, 0xca, 0xd2, 0x6d, 0x28, 0xe4,
	0xb2, 0x7f, 0x4b, 0xb9, 0x1b, 0x51, 0x88, 0xaf, 0xda, 0xfd, 0x11, 0xd0, 0x8c, 0x89, 0x9b, 0x20,
	0x15, 0x3c, 0xc9, 0xae, 0x79, 0x22, 0x67, 0x3f, 0x59, 0xba, 0xdf, 0x83, 0x76, 0x3e, 0x2a, 0x5f,
	0xcb, 0x60, 0x9f, 0x05, 0x6e, 0xc3, 0x49, 0xe0, 0xed, 0x90, 0x93, 0xc0, 0xeb, 0xff, 0x00, 0x9d,
	0x47, 0x62, 0x1a, 0xf2, 0x94, 0x3d, 0x41, 0x7e, 0x03, 0x74, 0xe0, 0xf7, 0x2a, 0x13, 0x2c, 0xc5,
	0x3d, 0x68, 0x25, 0x8f, 0x32, 0x87, 0xcf, 0xe8, 0x61, 0xa8, 0x1f, 0xc1, 0x79, 0x59, 0x15, 0xf3,
	0x28, 0x65, 0x78, 0x0c, 0x8d, 0x22, 0x2f, 0xf1, 0xea, 0xa0, 0x35, 0x56, 0xcb, 0xbf, 0xf4, 0x71,
	0x77, 0x5a, 0x82, 0xf8, 0x35, 0x28, 0x6b, 0x37, 0x75, 0x36, 0x3c, 0x29, 0xce, 0x82, 0x42, 0x1b,
	0x6b, 0x37, 0xbd, 0xe3, 0x49, 0xe9, 0xb2, 0x5a, 0xba, 0x1c, 0x7f, 0x38, 0xb8, 0x1c, 0xad, 0x6d,
	0x1c, 0xf3, 0x44, 0x60, 0x0d, 0x14, 0xca, 0xfc, 0x20, 0x15, 0x2c, 0xc1, 0xea, 0x73, 0x57, 0x63,
	0xf7, 0xd9, 0x4c, 0xff, 0xc5, 0xa0, 0xf2, 0x73, 0xe5, 0xca, 0x84, 0x3e, 0x4f, 0xfc, 0xe1, 0x3a,
	0x8b, 0x59, 0x12, 0x32, 0xcf, 0x67, 0xc9, 0xf0, 0xa3, 0xfb, 0x90, 0x04, 0xab, 0xb2, 0x4e, 0xde,
	0xe6, 0x7f, 0xfe, 0xe4, 0x07, 0x62, 0xbd, 0x7d, 0x18, 0xae, 0xf8, 0x66, 0x74, 0x80, 0x8e, 0x0a,
	0xb4, 0xb8, 0xd5, 0xd3, 0x91, 0x44, 0x1f, 0x8a, 0x4f, 0xc4, 0xaf, 0xff, 0x05, 0x00, 0x00, 0xff,
	0xff, 0x18, 0x6a, 0x79, 0xe7, 0x46, 0x06, 0x00, 0x00,
}
