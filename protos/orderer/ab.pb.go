


package orderer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "github.com/mcc-github/blockchain/protos/common"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 





type SeekInfo_SeekBehavior int32

const (
	SeekInfo_BLOCK_UNTIL_READY SeekInfo_SeekBehavior = 0
	SeekInfo_FAIL_IF_NOT_READY SeekInfo_SeekBehavior = 1
)

var SeekInfo_SeekBehavior_name = map[int32]string{
	0: "BLOCK_UNTIL_READY",
	1: "FAIL_IF_NOT_READY",
}
var SeekInfo_SeekBehavior_value = map[string]int32{
	"BLOCK_UNTIL_READY": 0,
	"FAIL_IF_NOT_READY": 1,
}

func (x SeekInfo_SeekBehavior) String() string {
	return proto.EnumName(SeekInfo_SeekBehavior_name, int32(x))
}
func (SeekInfo_SeekBehavior) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{5, 0}
}









type SeekInfo_SeekErrorResponse int32

const (
	SeekInfo_STRICT      SeekInfo_SeekErrorResponse = 0
	SeekInfo_BEST_EFFORT SeekInfo_SeekErrorResponse = 1
)

var SeekInfo_SeekErrorResponse_name = map[int32]string{
	0: "STRICT",
	1: "BEST_EFFORT",
}
var SeekInfo_SeekErrorResponse_value = map[string]int32{
	"STRICT":      0,
	"BEST_EFFORT": 1,
}

func (x SeekInfo_SeekErrorResponse) String() string {
	return proto.EnumName(SeekInfo_SeekErrorResponse_name, int32(x))
}
func (SeekInfo_SeekErrorResponse) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{5, 1}
}

type BroadcastResponse struct {
	
	Status common.Status `protobuf:"varint,1,opt,name=status,proto3,enum=common.Status" json:"status,omitempty"`
	
	Info                 string   `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BroadcastResponse) Reset()         { *m = BroadcastResponse{} }
func (m *BroadcastResponse) String() string { return proto.CompactTextString(m) }
func (*BroadcastResponse) ProtoMessage()    {}
func (*BroadcastResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{0}
}
func (m *BroadcastResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BroadcastResponse.Unmarshal(m, b)
}
func (m *BroadcastResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BroadcastResponse.Marshal(b, m, deterministic)
}
func (dst *BroadcastResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BroadcastResponse.Merge(dst, src)
}
func (m *BroadcastResponse) XXX_Size() int {
	return xxx_messageInfo_BroadcastResponse.Size(m)
}
func (m *BroadcastResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BroadcastResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BroadcastResponse proto.InternalMessageInfo

func (m *BroadcastResponse) GetStatus() common.Status {
	if m != nil {
		return m.Status
	}
	return common.Status_UNKNOWN
}

func (m *BroadcastResponse) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

type SeekNewest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SeekNewest) Reset()         { *m = SeekNewest{} }
func (m *SeekNewest) String() string { return proto.CompactTextString(m) }
func (*SeekNewest) ProtoMessage()    {}
func (*SeekNewest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{1}
}
func (m *SeekNewest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SeekNewest.Unmarshal(m, b)
}
func (m *SeekNewest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SeekNewest.Marshal(b, m, deterministic)
}
func (dst *SeekNewest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SeekNewest.Merge(dst, src)
}
func (m *SeekNewest) XXX_Size() int {
	return xxx_messageInfo_SeekNewest.Size(m)
}
func (m *SeekNewest) XXX_DiscardUnknown() {
	xxx_messageInfo_SeekNewest.DiscardUnknown(m)
}

var xxx_messageInfo_SeekNewest proto.InternalMessageInfo

type SeekOldest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SeekOldest) Reset()         { *m = SeekOldest{} }
func (m *SeekOldest) String() string { return proto.CompactTextString(m) }
func (*SeekOldest) ProtoMessage()    {}
func (*SeekOldest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{2}
}
func (m *SeekOldest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SeekOldest.Unmarshal(m, b)
}
func (m *SeekOldest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SeekOldest.Marshal(b, m, deterministic)
}
func (dst *SeekOldest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SeekOldest.Merge(dst, src)
}
func (m *SeekOldest) XXX_Size() int {
	return xxx_messageInfo_SeekOldest.Size(m)
}
func (m *SeekOldest) XXX_DiscardUnknown() {
	xxx_messageInfo_SeekOldest.DiscardUnknown(m)
}

var xxx_messageInfo_SeekOldest proto.InternalMessageInfo

type SeekSpecified struct {
	Number               uint64   `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SeekSpecified) Reset()         { *m = SeekSpecified{} }
func (m *SeekSpecified) String() string { return proto.CompactTextString(m) }
func (*SeekSpecified) ProtoMessage()    {}
func (*SeekSpecified) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{3}
}
func (m *SeekSpecified) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SeekSpecified.Unmarshal(m, b)
}
func (m *SeekSpecified) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SeekSpecified.Marshal(b, m, deterministic)
}
func (dst *SeekSpecified) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SeekSpecified.Merge(dst, src)
}
func (m *SeekSpecified) XXX_Size() int {
	return xxx_messageInfo_SeekSpecified.Size(m)
}
func (m *SeekSpecified) XXX_DiscardUnknown() {
	xxx_messageInfo_SeekSpecified.DiscardUnknown(m)
}

var xxx_messageInfo_SeekSpecified proto.InternalMessageInfo

func (m *SeekSpecified) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

type SeekPosition struct {
	
	
	
	
	Type                 isSeekPosition_Type `protobuf_oneof:"Type"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *SeekPosition) Reset()         { *m = SeekPosition{} }
func (m *SeekPosition) String() string { return proto.CompactTextString(m) }
func (*SeekPosition) ProtoMessage()    {}
func (*SeekPosition) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{4}
}
func (m *SeekPosition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SeekPosition.Unmarshal(m, b)
}
func (m *SeekPosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SeekPosition.Marshal(b, m, deterministic)
}
func (dst *SeekPosition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SeekPosition.Merge(dst, src)
}
func (m *SeekPosition) XXX_Size() int {
	return xxx_messageInfo_SeekPosition.Size(m)
}
func (m *SeekPosition) XXX_DiscardUnknown() {
	xxx_messageInfo_SeekPosition.DiscardUnknown(m)
}

var xxx_messageInfo_SeekPosition proto.InternalMessageInfo

type isSeekPosition_Type interface {
	isSeekPosition_Type()
}

type SeekPosition_Newest struct {
	Newest *SeekNewest `protobuf:"bytes,1,opt,name=newest,proto3,oneof"`
}

type SeekPosition_Oldest struct {
	Oldest *SeekOldest `protobuf:"bytes,2,opt,name=oldest,proto3,oneof"`
}

type SeekPosition_Specified struct {
	Specified *SeekSpecified `protobuf:"bytes,3,opt,name=specified,proto3,oneof"`
}

func (*SeekPosition_Newest) isSeekPosition_Type() {}

func (*SeekPosition_Oldest) isSeekPosition_Type() {}

func (*SeekPosition_Specified) isSeekPosition_Type() {}

func (m *SeekPosition) GetType() isSeekPosition_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *SeekPosition) GetNewest() *SeekNewest {
	if x, ok := m.GetType().(*SeekPosition_Newest); ok {
		return x.Newest
	}
	return nil
}

func (m *SeekPosition) GetOldest() *SeekOldest {
	if x, ok := m.GetType().(*SeekPosition_Oldest); ok {
		return x.Oldest
	}
	return nil
}

func (m *SeekPosition) GetSpecified() *SeekSpecified {
	if x, ok := m.GetType().(*SeekPosition_Specified); ok {
		return x.Specified
	}
	return nil
}


func (*SeekPosition) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _SeekPosition_OneofMarshaler, _SeekPosition_OneofUnmarshaler, _SeekPosition_OneofSizer, []interface{}{
		(*SeekPosition_Newest)(nil),
		(*SeekPosition_Oldest)(nil),
		(*SeekPosition_Specified)(nil),
	}
}

func _SeekPosition_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*SeekPosition)
	
	switch x := m.Type.(type) {
	case *SeekPosition_Newest:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Newest); err != nil {
			return err
		}
	case *SeekPosition_Oldest:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Oldest); err != nil {
			return err
		}
	case *SeekPosition_Specified:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Specified); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("SeekPosition.Type has unexpected type %T", x)
	}
	return nil
}

func _SeekPosition_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*SeekPosition)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SeekNewest)
		err := b.DecodeMessage(msg)
		m.Type = &SeekPosition_Newest{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SeekOldest)
		err := b.DecodeMessage(msg)
		m.Type = &SeekPosition_Oldest{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SeekSpecified)
		err := b.DecodeMessage(msg)
		m.Type = &SeekPosition_Specified{msg}
		return true, err
	default:
		return false, nil
	}
}

func _SeekPosition_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*SeekPosition)
	
	switch x := m.Type.(type) {
	case *SeekPosition_Newest:
		s := proto.Size(x.Newest)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *SeekPosition_Oldest:
		s := proto.Size(x.Oldest)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *SeekPosition_Specified:
		s := proto.Size(x.Specified)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}





type SeekInfo struct {
	Start                *SeekPosition              `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	Stop                 *SeekPosition              `protobuf:"bytes,2,opt,name=stop,proto3" json:"stop,omitempty"`
	Behavior             SeekInfo_SeekBehavior      `protobuf:"varint,3,opt,name=behavior,proto3,enum=orderer.SeekInfo_SeekBehavior" json:"behavior,omitempty"`
	ErrorResponse        SeekInfo_SeekErrorResponse `protobuf:"varint,4,opt,name=error_response,json=errorResponse,proto3,enum=orderer.SeekInfo_SeekErrorResponse" json:"error_response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *SeekInfo) Reset()         { *m = SeekInfo{} }
func (m *SeekInfo) String() string { return proto.CompactTextString(m) }
func (*SeekInfo) ProtoMessage()    {}
func (*SeekInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{5}
}
func (m *SeekInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SeekInfo.Unmarshal(m, b)
}
func (m *SeekInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SeekInfo.Marshal(b, m, deterministic)
}
func (dst *SeekInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SeekInfo.Merge(dst, src)
}
func (m *SeekInfo) XXX_Size() int {
	return xxx_messageInfo_SeekInfo.Size(m)
}
func (m *SeekInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SeekInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SeekInfo proto.InternalMessageInfo

func (m *SeekInfo) GetStart() *SeekPosition {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *SeekInfo) GetStop() *SeekPosition {
	if m != nil {
		return m.Stop
	}
	return nil
}

func (m *SeekInfo) GetBehavior() SeekInfo_SeekBehavior {
	if m != nil {
		return m.Behavior
	}
	return SeekInfo_BLOCK_UNTIL_READY
}

func (m *SeekInfo) GetErrorResponse() SeekInfo_SeekErrorResponse {
	if m != nil {
		return m.ErrorResponse
	}
	return SeekInfo_STRICT
}

type DeliverResponse struct {
	
	
	
	Type                 isDeliverResponse_Type `protobuf_oneof:"Type"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *DeliverResponse) Reset()         { *m = DeliverResponse{} }
func (m *DeliverResponse) String() string { return proto.CompactTextString(m) }
func (*DeliverResponse) ProtoMessage()    {}
func (*DeliverResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab_c21cb6da4aa42ea1, []int{6}
}
func (m *DeliverResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeliverResponse.Unmarshal(m, b)
}
func (m *DeliverResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeliverResponse.Marshal(b, m, deterministic)
}
func (dst *DeliverResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeliverResponse.Merge(dst, src)
}
func (m *DeliverResponse) XXX_Size() int {
	return xxx_messageInfo_DeliverResponse.Size(m)
}
func (m *DeliverResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeliverResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeliverResponse proto.InternalMessageInfo

type isDeliverResponse_Type interface {
	isDeliverResponse_Type()
}

type DeliverResponse_Status struct {
	Status common.Status `protobuf:"varint,1,opt,name=status,proto3,enum=common.Status,oneof"`
}

type DeliverResponse_Block struct {
	Block *common.Block `protobuf:"bytes,2,opt,name=block,proto3,oneof"`
}

func (*DeliverResponse_Status) isDeliverResponse_Type() {}

func (*DeliverResponse_Block) isDeliverResponse_Type() {}

func (m *DeliverResponse) GetType() isDeliverResponse_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *DeliverResponse) GetStatus() common.Status {
	if x, ok := m.GetType().(*DeliverResponse_Status); ok {
		return x.Status
	}
	return common.Status_UNKNOWN
}

func (m *DeliverResponse) GetBlock() *common.Block {
	if x, ok := m.GetType().(*DeliverResponse_Block); ok {
		return x.Block
	}
	return nil
}


func (*DeliverResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _DeliverResponse_OneofMarshaler, _DeliverResponse_OneofUnmarshaler, _DeliverResponse_OneofSizer, []interface{}{
		(*DeliverResponse_Status)(nil),
		(*DeliverResponse_Block)(nil),
	}
}

func _DeliverResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*DeliverResponse)
	
	switch x := m.Type.(type) {
	case *DeliverResponse_Status:
		b.EncodeVarint(1<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.Status))
	case *DeliverResponse_Block:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Block); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("DeliverResponse.Type has unexpected type %T", x)
	}
	return nil
}

func _DeliverResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*DeliverResponse)
	switch tag {
	case 1: 
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Type = &DeliverResponse_Status{common.Status(x)}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(common.Block)
		err := b.DecodeMessage(msg)
		m.Type = &DeliverResponse_Block{msg}
		return true, err
	default:
		return false, nil
	}
}

func _DeliverResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*DeliverResponse)
	
	switch x := m.Type.(type) {
	case *DeliverResponse_Status:
		n += 1 
		n += proto.SizeVarint(uint64(x.Status))
	case *DeliverResponse_Block:
		s := proto.Size(x.Block)
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
	proto.RegisterType((*BroadcastResponse)(nil), "orderer.BroadcastResponse")
	proto.RegisterType((*SeekNewest)(nil), "orderer.SeekNewest")
	proto.RegisterType((*SeekOldest)(nil), "orderer.SeekOldest")
	proto.RegisterType((*SeekSpecified)(nil), "orderer.SeekSpecified")
	proto.RegisterType((*SeekPosition)(nil), "orderer.SeekPosition")
	proto.RegisterType((*SeekInfo)(nil), "orderer.SeekInfo")
	proto.RegisterType((*DeliverResponse)(nil), "orderer.DeliverResponse")
	proto.RegisterEnum("orderer.SeekInfo_SeekBehavior", SeekInfo_SeekBehavior_name, SeekInfo_SeekBehavior_value)
	proto.RegisterEnum("orderer.SeekInfo_SeekErrorResponse", SeekInfo_SeekErrorResponse_name, SeekInfo_SeekErrorResponse_value)
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type AtomicBroadcastClient interface {
	
	Broadcast(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcast_BroadcastClient, error)
	
	Deliver(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcast_DeliverClient, error)
}

type atomicBroadcastClient struct {
	cc *grpc.ClientConn
}

func NewAtomicBroadcastClient(cc *grpc.ClientConn) AtomicBroadcastClient {
	return &atomicBroadcastClient{cc}
}

func (c *atomicBroadcastClient) Broadcast(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcast_BroadcastClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AtomicBroadcast_serviceDesc.Streams[0], "/orderer.AtomicBroadcast/Broadcast", opts...)
	if err != nil {
		return nil, err
	}
	x := &atomicBroadcastBroadcastClient{stream}
	return x, nil
}

type AtomicBroadcast_BroadcastClient interface {
	Send(*common.Envelope) error
	Recv() (*BroadcastResponse, error)
	grpc.ClientStream
}

type atomicBroadcastBroadcastClient struct {
	grpc.ClientStream
}

func (x *atomicBroadcastBroadcastClient) Send(m *common.Envelope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *atomicBroadcastBroadcastClient) Recv() (*BroadcastResponse, error) {
	m := new(BroadcastResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *atomicBroadcastClient) Deliver(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcast_DeliverClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AtomicBroadcast_serviceDesc.Streams[1], "/orderer.AtomicBroadcast/Deliver", opts...)
	if err != nil {
		return nil, err
	}
	x := &atomicBroadcastDeliverClient{stream}
	return x, nil
}

type AtomicBroadcast_DeliverClient interface {
	Send(*common.Envelope) error
	Recv() (*DeliverResponse, error)
	grpc.ClientStream
}

type atomicBroadcastDeliverClient struct {
	grpc.ClientStream
}

func (x *atomicBroadcastDeliverClient) Send(m *common.Envelope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *atomicBroadcastDeliverClient) Recv() (*DeliverResponse, error) {
	m := new(DeliverResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}


type AtomicBroadcastServer interface {
	
	Broadcast(AtomicBroadcast_BroadcastServer) error
	
	Deliver(AtomicBroadcast_DeliverServer) error
}

func RegisterAtomicBroadcastServer(s *grpc.Server, srv AtomicBroadcastServer) {
	s.RegisterService(&_AtomicBroadcast_serviceDesc, srv)
}

func _AtomicBroadcast_Broadcast_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AtomicBroadcastServer).Broadcast(&atomicBroadcastBroadcastServer{stream})
}

type AtomicBroadcast_BroadcastServer interface {
	Send(*BroadcastResponse) error
	Recv() (*common.Envelope, error)
	grpc.ServerStream
}

type atomicBroadcastBroadcastServer struct {
	grpc.ServerStream
}

func (x *atomicBroadcastBroadcastServer) Send(m *BroadcastResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *atomicBroadcastBroadcastServer) Recv() (*common.Envelope, error) {
	m := new(common.Envelope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _AtomicBroadcast_Deliver_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AtomicBroadcastServer).Deliver(&atomicBroadcastDeliverServer{stream})
}

type AtomicBroadcast_DeliverServer interface {
	Send(*DeliverResponse) error
	Recv() (*common.Envelope, error)
	grpc.ServerStream
}

type atomicBroadcastDeliverServer struct {
	grpc.ServerStream
}

func (x *atomicBroadcastDeliverServer) Send(m *DeliverResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *atomicBroadcastDeliverServer) Recv() (*common.Envelope, error) {
	m := new(common.Envelope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AtomicBroadcast_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orderer.AtomicBroadcast",
	HandlerType: (*AtomicBroadcastServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Broadcast",
			Handler:       _AtomicBroadcast_Broadcast_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Deliver",
			Handler:       _AtomicBroadcast_Deliver_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "orderer/ab.proto",
}

func init() { proto.RegisterFile("orderer/ab.proto", fileDescriptor_ab_c21cb6da4aa42ea1) }

var fileDescriptor_ab_c21cb6da4aa42ea1 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x93, 0x51, 0x6f, 0xda, 0x30,
	0x10, 0xc7, 0x49, 0x47, 0xd3, 0xf6, 0xda, 0x52, 0x70, 0xd5, 0x0a, 0xf1, 0x30, 0x55, 0x99, 0xba,
	0x31, 0x6d, 0x4b, 0x2a, 0x26, 0x4d, 0xda, 0xb4, 0x3d, 0x90, 0x16, 0x04, 0x1b, 0x82, 0xc9, 0xa4,
	0x0f, 0xdb, 0x4b, 0x94, 0x84, 0x03, 0xb2, 0x42, 0x1c, 0x39, 0x81, 0xa9, 0x9f, 0x61, 0x5f, 0x65,
	0x9f, 0x68, 0x9f, 0x66, 0xb2, 0xe3, 0x40, 0x69, 0xbb, 0x3d, 0x25, 0x77, 0xfe, 0xfd, 0xef, 0xfe,
	0xf6, 0xd9, 0x50, 0x66, 0x7c, 0x84, 0x1c, 0xb9, 0xe5, 0xf9, 0x66, 0xcc, 0x59, 0xca, 0xc8, 0x8e,
	0xca, 0xd4, 0x8e, 0x03, 0x36, 0x9f, 0xb3, 0xc8, 0xca, 0x3e, 0xd9, 0xaa, 0x31, 0x80, 0x8a, 0xcd,
	0x99, 0x37, 0x0a, 0xbc, 0x24, 0xa5, 0x98, 0xc4, 0x2c, 0x4a, 0x90, 0x3c, 0x07, 0x3d, 0x49, 0xbd,
	0x74, 0x91, 0x54, 0xb5, 0x33, 0xad, 0x5e, 0x6a, 0x94, 0x4c, 0xa5, 0x19, 0xca, 0x2c, 0x55, 0xab,
	0x84, 0x40, 0x31, 0x8c, 0xc6, 0xac, 0xba, 0x75, 0xa6, 0xd5, 0xf7, 0xa8, 0xfc, 0x37, 0x0e, 0x00,
	0x86, 0x88, 0x37, 0x7d, 0xfc, 0x89, 0x49, 0x9a, 0x47, 0x83, 0xd9, 0x48, 0x44, 0x2f, 0xe0, 0x50,
	0x44, 0xc3, 0x18, 0x83, 0x70, 0x1c, 0xe2, 0x88, 0x9c, 0x82, 0x1e, 0x2d, 0xe6, 0x3e, 0x72, 0xd9,
	0xa8, 0x48, 0x55, 0x64, 0xfc, 0xd6, 0xe0, 0x40, 0x90, 0x5f, 0x59, 0x12, 0xa6, 0x21, 0x8b, 0xc8,
	0x1b, 0xd0, 0x23, 0x59, 0x51, 0x82, 0xfb, 0x8d, 0x63, 0x53, 0xed, 0xca, 0x5c, 0x37, 0xeb, 0x14,
	0xa8, 0x82, 0x04, 0xce, 0x64, 0x4b, 0x69, 0xed, 0x3e, 0x9e, 0xb9, 0x11, 0x78, 0x06, 0x91, 0x77,
	0xb0, 0x97, 0xe4, 0x9e, 0xaa, 0x4f, 0xa4, 0xe2, 0x74, 0x43, 0xb1, 0x72, 0xdc, 0x29, 0xd0, 0x35,
	0x6a, 0xeb, 0x50, 0x74, 0x6e, 0x63, 0x34, 0xfe, 0x6c, 0xc1, 0xae, 0xc0, 0xba, 0xd1, 0x98, 0x91,
	0x57, 0xb0, 0x9d, 0xa4, 0x1e, 0xcf, 0x9d, 0x9e, 0x6c, 0x14, 0xca, 0x37, 0x44, 0x33, 0x86, 0xbc,
	0x84, 0x62, 0x92, 0xb2, 0x58, 0xd9, 0xfc, 0x07, 0x2b, 0x11, 0xf2, 0x01, 0x76, 0x7d, 0x9c, 0x7a,
	0xcb, 0x90, 0x71, 0xe9, 0xb1, 0xd4, 0x78, 0xba, 0x81, 0x8b, 0xe6, 0xf2, 0xc7, 0x56, 0x14, 0x5d,
	0xf1, 0xe4, 0x33, 0x94, 0x90, 0x73, 0xc6, 0x5d, 0xae, 0x46, 0x5c, 0x2d, 0xca, 0x0a, 0xcf, 0x1e,
	0xaf, 0xd0, 0x12, 0x6c, 0x7e, 0x1b, 0xe8, 0x21, 0xde, 0x0d, 0x8d, 0x8f, 0xd9, 0x68, 0xf2, 0x2e,
	0xe4, 0x04, 0x2a, 0x76, 0x6f, 0x70, 0xf9, 0xc5, 0xbd, 0xee, 0x3b, 0xdd, 0x9e, 0x4b, 0x5b, 0xcd,
	0xab, 0x6f, 0xe5, 0x82, 0x48, 0xb7, 0x9b, 0xdd, 0x9e, 0xdb, 0x6d, 0xbb, 0xfd, 0x81, 0xa3, 0xd2,
	0x9a, 0x71, 0x01, 0x95, 0x07, 0x1d, 0x08, 0x80, 0x3e, 0x74, 0x68, 0xf7, 0xd2, 0x29, 0x17, 0xc8,
	0x11, 0xec, 0xdb, 0xad, 0xa1, 0xe3, 0xb6, 0xda, 0xed, 0x01, 0x75, 0xca, 0x9a, 0xf1, 0x03, 0x8e,
	0xae, 0x70, 0x16, 0x2e, 0x71, 0xcd, 0xd7, 0xff, 0x7f, 0x3f, 0xc5, 0x64, 0xd5, 0x0d, 0x3d, 0x87,
	0x6d, 0x7f, 0xc6, 0x82, 0x1b, 0x75, 0xc0, 0x87, 0x39, 0x68, 0x8b, 0x64, 0xa7, 0x40, 0xb3, 0xd5,
	0x7c, 0x90, 0x8d, 0x5f, 0x1a, 0x1c, 0x35, 0x53, 0x36, 0x0f, 0x83, 0xd5, 0xa3, 0x20, 0x9f, 0x60,
	0x6f, 0x1d, 0x94, 0xf3, 0x02, 0xad, 0x68, 0x89, 0x33, 0x16, 0x63, 0xad, 0xb6, 0x3a, 0xc2, 0x07,
	0xef, 0xa8, 0xae, 0x5d, 0x68, 0xe4, 0x3d, 0xec, 0x28, 0xfb, 0x8f, 0x88, 0xab, 0x2b, 0xf1, 0xbd,
	0x2d, 0x0a, 0xa9, 0x7d, 0x0d, 0xe7, 0x8c, 0x4f, 0xcc, 0xe9, 0x6d, 0x8c, 0x7c, 0x86, 0xa3, 0x09,
	0x72, 0x73, 0xec, 0xf9, 0x3c, 0x0c, 0xb2, 0xb7, 0x9b, 0xe4, 0xe2, 0xef, 0xaf, 0x27, 0x61, 0x3a,
	0x5d, 0xf8, 0xa2, 0xbc, 0x75, 0x87, 0xb6, 0x32, 0xda, 0xca, 0x68, 0x4b, 0xd1, 0xbe, 0x2e, 0xe3,
	0xb7, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xf1, 0x8d, 0x09, 0x12, 0x2b, 0x04, 0x00, 0x00,
}
