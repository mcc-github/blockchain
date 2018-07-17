



package gossip

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common2 "github.com/mcc-github/blockchain/protos/common"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type PullMsgType int32

const (
	PullMsgType_UNDEFINED    PullMsgType = 0
	PullMsgType_BLOCK_MSG    PullMsgType = 1
	PullMsgType_IDENTITY_MSG PullMsgType = 2
)

var PullMsgType_name = map[int32]string{
	0: "UNDEFINED",
	1: "BLOCK_MSG",
	2: "IDENTITY_MSG",
}
var PullMsgType_value = map[string]int32{
	"UNDEFINED":    0,
	"BLOCK_MSG":    1,
	"IDENTITY_MSG": 2,
}

func (x PullMsgType) String() string {
	return proto.EnumName(PullMsgType_name, int32(x))
}
func (PullMsgType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type GossipMessage_Tag int32

const (
	GossipMessage_UNDEFINED    GossipMessage_Tag = 0
	GossipMessage_EMPTY        GossipMessage_Tag = 1
	GossipMessage_ORG_ONLY     GossipMessage_Tag = 2
	GossipMessage_CHAN_ONLY    GossipMessage_Tag = 3
	GossipMessage_CHAN_AND_ORG GossipMessage_Tag = 4
	GossipMessage_CHAN_OR_ORG  GossipMessage_Tag = 5
)

var GossipMessage_Tag_name = map[int32]string{
	0: "UNDEFINED",
	1: "EMPTY",
	2: "ORG_ONLY",
	3: "CHAN_ONLY",
	4: "CHAN_AND_ORG",
	5: "CHAN_OR_ORG",
}
var GossipMessage_Tag_value = map[string]int32{
	"UNDEFINED":    0,
	"EMPTY":        1,
	"ORG_ONLY":     2,
	"CHAN_ONLY":    3,
	"CHAN_AND_ORG": 4,
	"CHAN_OR_ORG":  5,
}

func (x GossipMessage_Tag) String() string {
	return proto.EnumName(GossipMessage_Tag_name, int32(x))
}
func (GossipMessage_Tag) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }





type Envelope struct {
	Payload        []byte          `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Signature      []byte          `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	SecretEnvelope *SecretEnvelope `protobuf:"bytes,3,opt,name=secret_envelope,json=secretEnvelope" json:"secret_envelope,omitempty"`
}

func (m *Envelope) Reset()                    { *m = Envelope{} }
func (m *Envelope) String() string            { return proto.CompactTextString(m) }
func (*Envelope) ProtoMessage()               {}
func (*Envelope) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Envelope) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Envelope) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Envelope) GetSecretEnvelope() *SecretEnvelope {
	if m != nil {
		return m.SecretEnvelope
	}
	return nil
}






type SecretEnvelope struct {
	Payload   []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *SecretEnvelope) Reset()                    { *m = SecretEnvelope{} }
func (m *SecretEnvelope) String() string            { return proto.CompactTextString(m) }
func (*SecretEnvelope) ProtoMessage()               {}
func (*SecretEnvelope) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SecretEnvelope) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *SecretEnvelope) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}




type Secret struct {
	
	
	Content isSecret_Content `protobuf_oneof:"content"`
}

func (m *Secret) Reset()                    { *m = Secret{} }
func (m *Secret) String() string            { return proto.CompactTextString(m) }
func (*Secret) ProtoMessage()               {}
func (*Secret) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isSecret_Content interface{ isSecret_Content() }

type Secret_InternalEndpoint struct {
	InternalEndpoint string `protobuf:"bytes,1,opt,name=internalEndpoint,oneof"`
}

func (*Secret_InternalEndpoint) isSecret_Content() {}

func (m *Secret) GetContent() isSecret_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *Secret) GetInternalEndpoint() string {
	if x, ok := m.GetContent().(*Secret_InternalEndpoint); ok {
		return x.InternalEndpoint
	}
	return ""
}


func (*Secret) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Secret_OneofMarshaler, _Secret_OneofUnmarshaler, _Secret_OneofSizer, []interface{}{
		(*Secret_InternalEndpoint)(nil),
	}
}

func _Secret_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Secret)
	
	switch x := m.Content.(type) {
	case *Secret_InternalEndpoint:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.InternalEndpoint)
	case nil:
	default:
		return fmt.Errorf("Secret.Content has unexpected type %T", x)
	}
	return nil
}

func _Secret_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Secret)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Content = &Secret_InternalEndpoint{x}
		return true, err
	default:
		return false, nil
	}
}

func _Secret_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Secret)
	
	switch x := m.Content.(type) {
	case *Secret_InternalEndpoint:
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.InternalEndpoint)))
		n += len(x.InternalEndpoint)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type GossipMessage struct {
	
	
	Nonce uint64 `protobuf:"varint,1,opt,name=nonce" json:"nonce,omitempty"`
	
	
	
	Channel []byte `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	
	
	Tag GossipMessage_Tag `protobuf:"varint,3,opt,name=tag,enum=gossip.GossipMessage_Tag" json:"tag,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Content isGossipMessage_Content `protobuf_oneof:"content"`
}

func (m *GossipMessage) Reset()                    { *m = GossipMessage{} }
func (m *GossipMessage) String() string            { return proto.CompactTextString(m) }
func (*GossipMessage) ProtoMessage()               {}
func (*GossipMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type isGossipMessage_Content interface{ isGossipMessage_Content() }

type GossipMessage_AliveMsg struct {
	AliveMsg *AliveMessage `protobuf:"bytes,5,opt,name=alive_msg,json=aliveMsg,oneof"`
}
type GossipMessage_MemReq struct {
	MemReq *MembershipRequest `protobuf:"bytes,6,opt,name=mem_req,json=memReq,oneof"`
}
type GossipMessage_MemRes struct {
	MemRes *MembershipResponse `protobuf:"bytes,7,opt,name=mem_res,json=memRes,oneof"`
}
type GossipMessage_DataMsg struct {
	DataMsg *DataMessage `protobuf:"bytes,8,opt,name=data_msg,json=dataMsg,oneof"`
}
type GossipMessage_Hello struct {
	Hello *GossipHello `protobuf:"bytes,9,opt,name=hello,oneof"`
}
type GossipMessage_DataDig struct {
	DataDig *DataDigest `protobuf:"bytes,10,opt,name=data_dig,json=dataDig,oneof"`
}
type GossipMessage_DataReq struct {
	DataReq *DataRequest `protobuf:"bytes,11,opt,name=data_req,json=dataReq,oneof"`
}
type GossipMessage_DataUpdate struct {
	DataUpdate *DataUpdate `protobuf:"bytes,12,opt,name=data_update,json=dataUpdate,oneof"`
}
type GossipMessage_Empty struct {
	Empty *Empty `protobuf:"bytes,13,opt,name=empty,oneof"`
}
type GossipMessage_Conn struct {
	Conn *ConnEstablish `protobuf:"bytes,14,opt,name=conn,oneof"`
}
type GossipMessage_StateInfo struct {
	StateInfo *StateInfo `protobuf:"bytes,15,opt,name=state_info,json=stateInfo,oneof"`
}
type GossipMessage_StateSnapshot struct {
	StateSnapshot *StateInfoSnapshot `protobuf:"bytes,16,opt,name=state_snapshot,json=stateSnapshot,oneof"`
}
type GossipMessage_StateInfoPullReq struct {
	StateInfoPullReq *StateInfoPullRequest `protobuf:"bytes,17,opt,name=state_info_pull_req,json=stateInfoPullReq,oneof"`
}
type GossipMessage_StateRequest struct {
	StateRequest *RemoteStateRequest `protobuf:"bytes,18,opt,name=state_request,json=stateRequest,oneof"`
}
type GossipMessage_StateResponse struct {
	StateResponse *RemoteStateResponse `protobuf:"bytes,19,opt,name=state_response,json=stateResponse,oneof"`
}
type GossipMessage_LeadershipMsg struct {
	LeadershipMsg *LeadershipMessage `protobuf:"bytes,20,opt,name=leadership_msg,json=leadershipMsg,oneof"`
}
type GossipMessage_PeerIdentity struct {
	PeerIdentity *PeerIdentity `protobuf:"bytes,21,opt,name=peer_identity,json=peerIdentity,oneof"`
}
type GossipMessage_Ack struct {
	Ack *Acknowledgement `protobuf:"bytes,22,opt,name=ack,oneof"`
}
type GossipMessage_PrivateReq struct {
	PrivateReq *RemotePvtDataRequest `protobuf:"bytes,23,opt,name=privateReq,oneof"`
}
type GossipMessage_PrivateRes struct {
	PrivateRes *RemotePvtDataResponse `protobuf:"bytes,24,opt,name=privateRes,oneof"`
}
type GossipMessage_PrivateData struct {
	PrivateData *PrivateDataMessage `protobuf:"bytes,25,opt,name=private_data,json=privateData,oneof"`
}

func (*GossipMessage_AliveMsg) isGossipMessage_Content()         {}
func (*GossipMessage_MemReq) isGossipMessage_Content()           {}
func (*GossipMessage_MemRes) isGossipMessage_Content()           {}
func (*GossipMessage_DataMsg) isGossipMessage_Content()          {}
func (*GossipMessage_Hello) isGossipMessage_Content()            {}
func (*GossipMessage_DataDig) isGossipMessage_Content()          {}
func (*GossipMessage_DataReq) isGossipMessage_Content()          {}
func (*GossipMessage_DataUpdate) isGossipMessage_Content()       {}
func (*GossipMessage_Empty) isGossipMessage_Content()            {}
func (*GossipMessage_Conn) isGossipMessage_Content()             {}
func (*GossipMessage_StateInfo) isGossipMessage_Content()        {}
func (*GossipMessage_StateSnapshot) isGossipMessage_Content()    {}
func (*GossipMessage_StateInfoPullReq) isGossipMessage_Content() {}
func (*GossipMessage_StateRequest) isGossipMessage_Content()     {}
func (*GossipMessage_StateResponse) isGossipMessage_Content()    {}
func (*GossipMessage_LeadershipMsg) isGossipMessage_Content()    {}
func (*GossipMessage_PeerIdentity) isGossipMessage_Content()     {}
func (*GossipMessage_Ack) isGossipMessage_Content()              {}
func (*GossipMessage_PrivateReq) isGossipMessage_Content()       {}
func (*GossipMessage_PrivateRes) isGossipMessage_Content()       {}
func (*GossipMessage_PrivateData) isGossipMessage_Content()      {}

func (m *GossipMessage) GetContent() isGossipMessage_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *GossipMessage) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *GossipMessage) GetChannel() []byte {
	if m != nil {
		return m.Channel
	}
	return nil
}

func (m *GossipMessage) GetTag() GossipMessage_Tag {
	if m != nil {
		return m.Tag
	}
	return GossipMessage_UNDEFINED
}

func (m *GossipMessage) GetAliveMsg() *AliveMessage {
	if x, ok := m.GetContent().(*GossipMessage_AliveMsg); ok {
		return x.AliveMsg
	}
	return nil
}

func (m *GossipMessage) GetMemReq() *MembershipRequest {
	if x, ok := m.GetContent().(*GossipMessage_MemReq); ok {
		return x.MemReq
	}
	return nil
}

func (m *GossipMessage) GetMemRes() *MembershipResponse {
	if x, ok := m.GetContent().(*GossipMessage_MemRes); ok {
		return x.MemRes
	}
	return nil
}

func (m *GossipMessage) GetDataMsg() *DataMessage {
	if x, ok := m.GetContent().(*GossipMessage_DataMsg); ok {
		return x.DataMsg
	}
	return nil
}

func (m *GossipMessage) GetHello() *GossipHello {
	if x, ok := m.GetContent().(*GossipMessage_Hello); ok {
		return x.Hello
	}
	return nil
}

func (m *GossipMessage) GetDataDig() *DataDigest {
	if x, ok := m.GetContent().(*GossipMessage_DataDig); ok {
		return x.DataDig
	}
	return nil
}

func (m *GossipMessage) GetDataReq() *DataRequest {
	if x, ok := m.GetContent().(*GossipMessage_DataReq); ok {
		return x.DataReq
	}
	return nil
}

func (m *GossipMessage) GetDataUpdate() *DataUpdate {
	if x, ok := m.GetContent().(*GossipMessage_DataUpdate); ok {
		return x.DataUpdate
	}
	return nil
}

func (m *GossipMessage) GetEmpty() *Empty {
	if x, ok := m.GetContent().(*GossipMessage_Empty); ok {
		return x.Empty
	}
	return nil
}

func (m *GossipMessage) GetConn() *ConnEstablish {
	if x, ok := m.GetContent().(*GossipMessage_Conn); ok {
		return x.Conn
	}
	return nil
}

func (m *GossipMessage) GetStateInfo() *StateInfo {
	if x, ok := m.GetContent().(*GossipMessage_StateInfo); ok {
		return x.StateInfo
	}
	return nil
}

func (m *GossipMessage) GetStateSnapshot() *StateInfoSnapshot {
	if x, ok := m.GetContent().(*GossipMessage_StateSnapshot); ok {
		return x.StateSnapshot
	}
	return nil
}

func (m *GossipMessage) GetStateInfoPullReq() *StateInfoPullRequest {
	if x, ok := m.GetContent().(*GossipMessage_StateInfoPullReq); ok {
		return x.StateInfoPullReq
	}
	return nil
}

func (m *GossipMessage) GetStateRequest() *RemoteStateRequest {
	if x, ok := m.GetContent().(*GossipMessage_StateRequest); ok {
		return x.StateRequest
	}
	return nil
}

func (m *GossipMessage) GetStateResponse() *RemoteStateResponse {
	if x, ok := m.GetContent().(*GossipMessage_StateResponse); ok {
		return x.StateResponse
	}
	return nil
}

func (m *GossipMessage) GetLeadershipMsg() *LeadershipMessage {
	if x, ok := m.GetContent().(*GossipMessage_LeadershipMsg); ok {
		return x.LeadershipMsg
	}
	return nil
}

func (m *GossipMessage) GetPeerIdentity() *PeerIdentity {
	if x, ok := m.GetContent().(*GossipMessage_PeerIdentity); ok {
		return x.PeerIdentity
	}
	return nil
}

func (m *GossipMessage) GetAck() *Acknowledgement {
	if x, ok := m.GetContent().(*GossipMessage_Ack); ok {
		return x.Ack
	}
	return nil
}

func (m *GossipMessage) GetPrivateReq() *RemotePvtDataRequest {
	if x, ok := m.GetContent().(*GossipMessage_PrivateReq); ok {
		return x.PrivateReq
	}
	return nil
}

func (m *GossipMessage) GetPrivateRes() *RemotePvtDataResponse {
	if x, ok := m.GetContent().(*GossipMessage_PrivateRes); ok {
		return x.PrivateRes
	}
	return nil
}

func (m *GossipMessage) GetPrivateData() *PrivateDataMessage {
	if x, ok := m.GetContent().(*GossipMessage_PrivateData); ok {
		return x.PrivateData
	}
	return nil
}


func (*GossipMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _GossipMessage_OneofMarshaler, _GossipMessage_OneofUnmarshaler, _GossipMessage_OneofSizer, []interface{}{
		(*GossipMessage_AliveMsg)(nil),
		(*GossipMessage_MemReq)(nil),
		(*GossipMessage_MemRes)(nil),
		(*GossipMessage_DataMsg)(nil),
		(*GossipMessage_Hello)(nil),
		(*GossipMessage_DataDig)(nil),
		(*GossipMessage_DataReq)(nil),
		(*GossipMessage_DataUpdate)(nil),
		(*GossipMessage_Empty)(nil),
		(*GossipMessage_Conn)(nil),
		(*GossipMessage_StateInfo)(nil),
		(*GossipMessage_StateSnapshot)(nil),
		(*GossipMessage_StateInfoPullReq)(nil),
		(*GossipMessage_StateRequest)(nil),
		(*GossipMessage_StateResponse)(nil),
		(*GossipMessage_LeadershipMsg)(nil),
		(*GossipMessage_PeerIdentity)(nil),
		(*GossipMessage_Ack)(nil),
		(*GossipMessage_PrivateReq)(nil),
		(*GossipMessage_PrivateRes)(nil),
		(*GossipMessage_PrivateData)(nil),
	}
}

func _GossipMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*GossipMessage)
	
	switch x := m.Content.(type) {
	case *GossipMessage_AliveMsg:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.AliveMsg); err != nil {
			return err
		}
	case *GossipMessage_MemReq:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MemReq); err != nil {
			return err
		}
	case *GossipMessage_MemRes:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MemRes); err != nil {
			return err
		}
	case *GossipMessage_DataMsg:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DataMsg); err != nil {
			return err
		}
	case *GossipMessage_Hello:
		b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Hello); err != nil {
			return err
		}
	case *GossipMessage_DataDig:
		b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DataDig); err != nil {
			return err
		}
	case *GossipMessage_DataReq:
		b.EncodeVarint(11<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DataReq); err != nil {
			return err
		}
	case *GossipMessage_DataUpdate:
		b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DataUpdate); err != nil {
			return err
		}
	case *GossipMessage_Empty:
		b.EncodeVarint(13<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Empty); err != nil {
			return err
		}
	case *GossipMessage_Conn:
		b.EncodeVarint(14<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Conn); err != nil {
			return err
		}
	case *GossipMessage_StateInfo:
		b.EncodeVarint(15<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StateInfo); err != nil {
			return err
		}
	case *GossipMessage_StateSnapshot:
		b.EncodeVarint(16<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StateSnapshot); err != nil {
			return err
		}
	case *GossipMessage_StateInfoPullReq:
		b.EncodeVarint(17<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StateInfoPullReq); err != nil {
			return err
		}
	case *GossipMessage_StateRequest:
		b.EncodeVarint(18<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StateRequest); err != nil {
			return err
		}
	case *GossipMessage_StateResponse:
		b.EncodeVarint(19<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StateResponse); err != nil {
			return err
		}
	case *GossipMessage_LeadershipMsg:
		b.EncodeVarint(20<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LeadershipMsg); err != nil {
			return err
		}
	case *GossipMessage_PeerIdentity:
		b.EncodeVarint(21<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PeerIdentity); err != nil {
			return err
		}
	case *GossipMessage_Ack:
		b.EncodeVarint(22<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ack); err != nil {
			return err
		}
	case *GossipMessage_PrivateReq:
		b.EncodeVarint(23<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PrivateReq); err != nil {
			return err
		}
	case *GossipMessage_PrivateRes:
		b.EncodeVarint(24<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PrivateRes); err != nil {
			return err
		}
	case *GossipMessage_PrivateData:
		b.EncodeVarint(25<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PrivateData); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("GossipMessage.Content has unexpected type %T", x)
	}
	return nil
}

func _GossipMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*GossipMessage)
	switch tag {
	case 5: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(AliveMessage)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_AliveMsg{msg}
		return true, err
	case 6: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MembershipRequest)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_MemReq{msg}
		return true, err
	case 7: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MembershipResponse)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_MemRes{msg}
		return true, err
	case 8: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DataMessage)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_DataMsg{msg}
		return true, err
	case 9: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GossipHello)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_Hello{msg}
		return true, err
	case 10: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DataDigest)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_DataDig{msg}
		return true, err
	case 11: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DataRequest)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_DataReq{msg}
		return true, err
	case 12: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DataUpdate)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_DataUpdate{msg}
		return true, err
	case 13: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Empty)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_Empty{msg}
		return true, err
	case 14: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ConnEstablish)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_Conn{msg}
		return true, err
	case 15: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(StateInfo)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_StateInfo{msg}
		return true, err
	case 16: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(StateInfoSnapshot)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_StateSnapshot{msg}
		return true, err
	case 17: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(StateInfoPullRequest)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_StateInfoPullReq{msg}
		return true, err
	case 18: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RemoteStateRequest)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_StateRequest{msg}
		return true, err
	case 19: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RemoteStateResponse)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_StateResponse{msg}
		return true, err
	case 20: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LeadershipMessage)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_LeadershipMsg{msg}
		return true, err
	case 21: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PeerIdentity)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_PeerIdentity{msg}
		return true, err
	case 22: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Acknowledgement)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_Ack{msg}
		return true, err
	case 23: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RemotePvtDataRequest)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_PrivateReq{msg}
		return true, err
	case 24: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RemotePvtDataResponse)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_PrivateRes{msg}
		return true, err
	case 25: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PrivateDataMessage)
		err := b.DecodeMessage(msg)
		m.Content = &GossipMessage_PrivateData{msg}
		return true, err
	default:
		return false, nil
	}
}

func _GossipMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*GossipMessage)
	
	switch x := m.Content.(type) {
	case *GossipMessage_AliveMsg:
		s := proto.Size(x.AliveMsg)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_MemReq:
		s := proto.Size(x.MemReq)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_MemRes:
		s := proto.Size(x.MemRes)
		n += proto.SizeVarint(7<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_DataMsg:
		s := proto.Size(x.DataMsg)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_Hello:
		s := proto.Size(x.Hello)
		n += proto.SizeVarint(9<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_DataDig:
		s := proto.Size(x.DataDig)
		n += proto.SizeVarint(10<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_DataReq:
		s := proto.Size(x.DataReq)
		n += proto.SizeVarint(11<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_DataUpdate:
		s := proto.Size(x.DataUpdate)
		n += proto.SizeVarint(12<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_Empty:
		s := proto.Size(x.Empty)
		n += proto.SizeVarint(13<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_Conn:
		s := proto.Size(x.Conn)
		n += proto.SizeVarint(14<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_StateInfo:
		s := proto.Size(x.StateInfo)
		n += proto.SizeVarint(15<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_StateSnapshot:
		s := proto.Size(x.StateSnapshot)
		n += proto.SizeVarint(16<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_StateInfoPullReq:
		s := proto.Size(x.StateInfoPullReq)
		n += proto.SizeVarint(17<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_StateRequest:
		s := proto.Size(x.StateRequest)
		n += proto.SizeVarint(18<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_StateResponse:
		s := proto.Size(x.StateResponse)
		n += proto.SizeVarint(19<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_LeadershipMsg:
		s := proto.Size(x.LeadershipMsg)
		n += proto.SizeVarint(20<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_PeerIdentity:
		s := proto.Size(x.PeerIdentity)
		n += proto.SizeVarint(21<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_Ack:
		s := proto.Size(x.Ack)
		n += proto.SizeVarint(22<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_PrivateReq:
		s := proto.Size(x.PrivateReq)
		n += proto.SizeVarint(23<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_PrivateRes:
		s := proto.Size(x.PrivateRes)
		n += proto.SizeVarint(24<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GossipMessage_PrivateData:
		s := proto.Size(x.PrivateData)
		n += proto.SizeVarint(25<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}



type StateInfo struct {
	Timestamp *PeerTime `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp,omitempty"`
	PkiId     []byte    `protobuf:"bytes,3,opt,name=pki_id,json=pkiId,proto3" json:"pki_id,omitempty"`
	
	
	
	Channel_MAC []byte      `protobuf:"bytes,4,opt,name=channel_MAC,json=channelMAC,proto3" json:"channel_MAC,omitempty"`
	Properties  *Properties `protobuf:"bytes,5,opt,name=properties" json:"properties,omitempty"`
}

func (m *StateInfo) Reset()                    { *m = StateInfo{} }
func (m *StateInfo) String() string            { return proto.CompactTextString(m) }
func (*StateInfo) ProtoMessage()               {}
func (*StateInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *StateInfo) GetTimestamp() *PeerTime {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *StateInfo) GetPkiId() []byte {
	if m != nil {
		return m.PkiId
	}
	return nil
}

func (m *StateInfo) GetChannel_MAC() []byte {
	if m != nil {
		return m.Channel_MAC
	}
	return nil
}

func (m *StateInfo) GetProperties() *Properties {
	if m != nil {
		return m.Properties
	}
	return nil
}

type Properties struct {
	LedgerHeight uint64       `protobuf:"varint,1,opt,name=ledger_height,json=ledgerHeight" json:"ledger_height,omitempty"`
	LeftChannel  bool         `protobuf:"varint,2,opt,name=left_channel,json=leftChannel" json:"left_channel,omitempty"`
	Chaincodes   []*Chaincode `protobuf:"bytes,3,rep,name=chaincodes" json:"chaincodes,omitempty"`
}

func (m *Properties) Reset()                    { *m = Properties{} }
func (m *Properties) String() string            { return proto.CompactTextString(m) }
func (*Properties) ProtoMessage()               {}
func (*Properties) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Properties) GetLedgerHeight() uint64 {
	if m != nil {
		return m.LedgerHeight
	}
	return 0
}

func (m *Properties) GetLeftChannel() bool {
	if m != nil {
		return m.LeftChannel
	}
	return false
}

func (m *Properties) GetChaincodes() []*Chaincode {
	if m != nil {
		return m.Chaincodes
	}
	return nil
}


type StateInfoSnapshot struct {
	Elements []*Envelope `protobuf:"bytes,1,rep,name=elements" json:"elements,omitempty"`
}

func (m *StateInfoSnapshot) Reset()                    { *m = StateInfoSnapshot{} }
func (m *StateInfoSnapshot) String() string            { return proto.CompactTextString(m) }
func (*StateInfoSnapshot) ProtoMessage()               {}
func (*StateInfoSnapshot) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *StateInfoSnapshot) GetElements() []*Envelope {
	if m != nil {
		return m.Elements
	}
	return nil
}



type StateInfoPullRequest struct {
	
	
	
	Channel_MAC []byte `protobuf:"bytes,1,opt,name=channel_MAC,json=channelMAC,proto3" json:"channel_MAC,omitempty"`
}

func (m *StateInfoPullRequest) Reset()                    { *m = StateInfoPullRequest{} }
func (m *StateInfoPullRequest) String() string            { return proto.CompactTextString(m) }
func (*StateInfoPullRequest) ProtoMessage()               {}
func (*StateInfoPullRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *StateInfoPullRequest) GetChannel_MAC() []byte {
	if m != nil {
		return m.Channel_MAC
	}
	return nil
}




type ConnEstablish struct {
	PkiId       []byte `protobuf:"bytes,1,opt,name=pki_id,json=pkiId,proto3" json:"pki_id,omitempty"`
	Identity    []byte `protobuf:"bytes,2,opt,name=identity,proto3" json:"identity,omitempty"`
	TlsCertHash []byte `protobuf:"bytes,3,opt,name=tls_cert_hash,json=tlsCertHash,proto3" json:"tls_cert_hash,omitempty"`
}

func (m *ConnEstablish) Reset()                    { *m = ConnEstablish{} }
func (m *ConnEstablish) String() string            { return proto.CompactTextString(m) }
func (*ConnEstablish) ProtoMessage()               {}
func (*ConnEstablish) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ConnEstablish) GetPkiId() []byte {
	if m != nil {
		return m.PkiId
	}
	return nil
}

func (m *ConnEstablish) GetIdentity() []byte {
	if m != nil {
		return m.Identity
	}
	return nil
}

func (m *ConnEstablish) GetTlsCertHash() []byte {
	if m != nil {
		return m.TlsCertHash
	}
	return nil
}




type PeerIdentity struct {
	PkiId    []byte `protobuf:"bytes,1,opt,name=pki_id,json=pkiId,proto3" json:"pki_id,omitempty"`
	Cert     []byte `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"`
	Metadata []byte `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (m *PeerIdentity) Reset()                    { *m = PeerIdentity{} }
func (m *PeerIdentity) String() string            { return proto.CompactTextString(m) }
func (*PeerIdentity) ProtoMessage()               {}
func (*PeerIdentity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *PeerIdentity) GetPkiId() []byte {
	if m != nil {
		return m.PkiId
	}
	return nil
}

func (m *PeerIdentity) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *PeerIdentity) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}



type DataRequest struct {
	Nonce   uint64      `protobuf:"varint,1,opt,name=nonce" json:"nonce,omitempty"`
	Digests []string    `protobuf:"bytes,2,rep,name=digests" json:"digests,omitempty"`
	MsgType PullMsgType `protobuf:"varint,3,opt,name=msg_type,json=msgType,enum=gossip.PullMsgType" json:"msg_type,omitempty"`
}

func (m *DataRequest) Reset()                    { *m = DataRequest{} }
func (m *DataRequest) String() string            { return proto.CompactTextString(m) }
func (*DataRequest) ProtoMessage()               {}
func (*DataRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *DataRequest) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *DataRequest) GetDigests() []string {
	if m != nil {
		return m.Digests
	}
	return nil
}

func (m *DataRequest) GetMsgType() PullMsgType {
	if m != nil {
		return m.MsgType
	}
	return PullMsgType_UNDEFINED
}



type GossipHello struct {
	Nonce    uint64      `protobuf:"varint,1,opt,name=nonce" json:"nonce,omitempty"`
	Metadata []byte      `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	MsgType  PullMsgType `protobuf:"varint,3,opt,name=msg_type,json=msgType,enum=gossip.PullMsgType" json:"msg_type,omitempty"`
}

func (m *GossipHello) Reset()                    { *m = GossipHello{} }
func (m *GossipHello) String() string            { return proto.CompactTextString(m) }
func (*GossipHello) ProtoMessage()               {}
func (*GossipHello) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *GossipHello) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *GossipHello) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *GossipHello) GetMsgType() PullMsgType {
	if m != nil {
		return m.MsgType
	}
	return PullMsgType_UNDEFINED
}



type DataUpdate struct {
	Nonce   uint64      `protobuf:"varint,1,opt,name=nonce" json:"nonce,omitempty"`
	Data    []*Envelope `protobuf:"bytes,2,rep,name=data" json:"data,omitempty"`
	MsgType PullMsgType `protobuf:"varint,3,opt,name=msg_type,json=msgType,enum=gossip.PullMsgType" json:"msg_type,omitempty"`
}

func (m *DataUpdate) Reset()                    { *m = DataUpdate{} }
func (m *DataUpdate) String() string            { return proto.CompactTextString(m) }
func (*DataUpdate) ProtoMessage()               {}
func (*DataUpdate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *DataUpdate) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *DataUpdate) GetData() []*Envelope {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *DataUpdate) GetMsgType() PullMsgType {
	if m != nil {
		return m.MsgType
	}
	return PullMsgType_UNDEFINED
}



type DataDigest struct {
	Nonce   uint64      `protobuf:"varint,1,opt,name=nonce" json:"nonce,omitempty"`
	Digests []string    `protobuf:"bytes,2,rep,name=digests" json:"digests,omitempty"`
	MsgType PullMsgType `protobuf:"varint,3,opt,name=msg_type,json=msgType,enum=gossip.PullMsgType" json:"msg_type,omitempty"`
}

func (m *DataDigest) Reset()                    { *m = DataDigest{} }
func (m *DataDigest) String() string            { return proto.CompactTextString(m) }
func (*DataDigest) ProtoMessage()               {}
func (*DataDigest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *DataDigest) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *DataDigest) GetDigests() []string {
	if m != nil {
		return m.Digests
	}
	return nil
}

func (m *DataDigest) GetMsgType() PullMsgType {
	if m != nil {
		return m.MsgType
	}
	return PullMsgType_UNDEFINED
}


type DataMessage struct {
	Payload *Payload `protobuf:"bytes,1,opt,name=payload" json:"payload,omitempty"`
}

func (m *DataMessage) Reset()                    { *m = DataMessage{} }
func (m *DataMessage) String() string            { return proto.CompactTextString(m) }
func (*DataMessage) ProtoMessage()               {}
func (*DataMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *DataMessage) GetPayload() *Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}




type PrivateDataMessage struct {
	Payload *PrivatePayload `protobuf:"bytes,1,opt,name=payload" json:"payload,omitempty"`
}

func (m *PrivateDataMessage) Reset()                    { *m = PrivateDataMessage{} }
func (m *PrivateDataMessage) String() string            { return proto.CompactTextString(m) }
func (*PrivateDataMessage) ProtoMessage()               {}
func (*PrivateDataMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *PrivateDataMessage) GetPayload() *PrivatePayload {
	if m != nil {
		return m.Payload
	}
	return nil
}


type Payload struct {
	SeqNum      uint64   `protobuf:"varint,1,opt,name=seq_num,json=seqNum" json:"seq_num,omitempty"`
	Data        []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	PrivateData [][]byte `protobuf:"bytes,3,rep,name=private_data,json=privateData,proto3" json:"private_data,omitempty"`
}

func (m *Payload) Reset()                    { *m = Payload{} }
func (m *Payload) String() string            { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()               {}
func (*Payload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *Payload) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *Payload) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Payload) GetPrivateData() [][]byte {
	if m != nil {
		return m.PrivateData
	}
	return nil
}




type PrivatePayload struct {
	CollectionName    string                           `protobuf:"bytes,1,opt,name=collection_name,json=collectionName" json:"collection_name,omitempty"`
	Namespace         string                           `protobuf:"bytes,2,opt,name=namespace" json:"namespace,omitempty"`
	TxId              string                           `protobuf:"bytes,3,opt,name=tx_id,json=txId" json:"tx_id,omitempty"`
	PrivateRwset      []byte                           `protobuf:"bytes,4,opt,name=private_rwset,json=privateRwset,proto3" json:"private_rwset,omitempty"`
	PrivateSimHeight  uint64                           `protobuf:"varint,5,opt,name=private_sim_height,json=privateSimHeight" json:"private_sim_height,omitempty"`
	CollectionConfigs *common2.CollectionConfigPackage `protobuf:"bytes,6,opt,name=collection_configs,json=collectionConfigs" json:"collection_configs,omitempty"`
}

func (m *PrivatePayload) Reset()                    { *m = PrivatePayload{} }
func (m *PrivatePayload) String() string            { return proto.CompactTextString(m) }
func (*PrivatePayload) ProtoMessage()               {}
func (*PrivatePayload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *PrivatePayload) GetCollectionName() string {
	if m != nil {
		return m.CollectionName
	}
	return ""
}

func (m *PrivatePayload) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *PrivatePayload) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *PrivatePayload) GetPrivateRwset() []byte {
	if m != nil {
		return m.PrivateRwset
	}
	return nil
}

func (m *PrivatePayload) GetPrivateSimHeight() uint64 {
	if m != nil {
		return m.PrivateSimHeight
	}
	return 0
}

func (m *PrivatePayload) GetCollectionConfigs() *common2.CollectionConfigPackage {
	if m != nil {
		return m.CollectionConfigs
	}
	return nil
}



type AliveMessage struct {
	Membership *Member   `protobuf:"bytes,1,opt,name=membership" json:"membership,omitempty"`
	Timestamp  *PeerTime `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp,omitempty"`
	Identity   []byte    `protobuf:"bytes,4,opt,name=identity,proto3" json:"identity,omitempty"`
}

func (m *AliveMessage) Reset()                    { *m = AliveMessage{} }
func (m *AliveMessage) String() string            { return proto.CompactTextString(m) }
func (*AliveMessage) ProtoMessage()               {}
func (*AliveMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func (m *AliveMessage) GetMembership() *Member {
	if m != nil {
		return m.Membership
	}
	return nil
}

func (m *AliveMessage) GetTimestamp() *PeerTime {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *AliveMessage) GetIdentity() []byte {
	if m != nil {
		return m.Identity
	}
	return nil
}



type LeadershipMessage struct {
	PkiId         []byte    `protobuf:"bytes,1,opt,name=pki_id,json=pkiId,proto3" json:"pki_id,omitempty"`
	Timestamp     *PeerTime `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp,omitempty"`
	IsDeclaration bool      `protobuf:"varint,3,opt,name=is_declaration,json=isDeclaration" json:"is_declaration,omitempty"`
}

func (m *LeadershipMessage) Reset()                    { *m = LeadershipMessage{} }
func (m *LeadershipMessage) String() string            { return proto.CompactTextString(m) }
func (*LeadershipMessage) ProtoMessage()               {}
func (*LeadershipMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

func (m *LeadershipMessage) GetPkiId() []byte {
	if m != nil {
		return m.PkiId
	}
	return nil
}

func (m *LeadershipMessage) GetTimestamp() *PeerTime {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *LeadershipMessage) GetIsDeclaration() bool {
	if m != nil {
		return m.IsDeclaration
	}
	return false
}


type PeerTime struct {
	IncNum uint64 `protobuf:"varint,1,opt,name=inc_num,json=incNum" json:"inc_num,omitempty"`
	SeqNum uint64 `protobuf:"varint,2,opt,name=seq_num,json=seqNum" json:"seq_num,omitempty"`
}

func (m *PeerTime) Reset()                    { *m = PeerTime{} }
func (m *PeerTime) String() string            { return proto.CompactTextString(m) }
func (*PeerTime) ProtoMessage()               {}
func (*PeerTime) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

func (m *PeerTime) GetIncNum() uint64 {
	if m != nil {
		return m.IncNum
	}
	return 0
}

func (m *PeerTime) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}



type MembershipRequest struct {
	SelfInformation *Envelope `protobuf:"bytes,1,opt,name=self_information,json=selfInformation" json:"self_information,omitempty"`
	Known           [][]byte  `protobuf:"bytes,2,rep,name=known,proto3" json:"known,omitempty"`
}

func (m *MembershipRequest) Reset()                    { *m = MembershipRequest{} }
func (m *MembershipRequest) String() string            { return proto.CompactTextString(m) }
func (*MembershipRequest) ProtoMessage()               {}
func (*MembershipRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{21} }

func (m *MembershipRequest) GetSelfInformation() *Envelope {
	if m != nil {
		return m.SelfInformation
	}
	return nil
}

func (m *MembershipRequest) GetKnown() [][]byte {
	if m != nil {
		return m.Known
	}
	return nil
}


type MembershipResponse struct {
	Alive []*Envelope `protobuf:"bytes,1,rep,name=alive" json:"alive,omitempty"`
	Dead  []*Envelope `protobuf:"bytes,2,rep,name=dead" json:"dead,omitempty"`
}

func (m *MembershipResponse) Reset()                    { *m = MembershipResponse{} }
func (m *MembershipResponse) String() string            { return proto.CompactTextString(m) }
func (*MembershipResponse) ProtoMessage()               {}
func (*MembershipResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{22} }

func (m *MembershipResponse) GetAlive() []*Envelope {
	if m != nil {
		return m.Alive
	}
	return nil
}

func (m *MembershipResponse) GetDead() []*Envelope {
	if m != nil {
		return m.Dead
	}
	return nil
}



type Member struct {
	Endpoint string `protobuf:"bytes,1,opt,name=endpoint" json:"endpoint,omitempty"`
	Metadata []byte `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	PkiId    []byte `protobuf:"bytes,3,opt,name=pki_id,json=pkiId,proto3" json:"pki_id,omitempty"`
}

func (m *Member) Reset()                    { *m = Member{} }
func (m *Member) String() string            { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()               {}
func (*Member) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{23} }

func (m *Member) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *Member) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Member) GetPkiId() []byte {
	if m != nil {
		return m.PkiId
	}
	return nil
}


type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{24} }



type RemoteStateRequest struct {
	StartSeqNum uint64 `protobuf:"varint,1,opt,name=start_seq_num,json=startSeqNum" json:"start_seq_num,omitempty"`
	EndSeqNum   uint64 `protobuf:"varint,2,opt,name=end_seq_num,json=endSeqNum" json:"end_seq_num,omitempty"`
}

func (m *RemoteStateRequest) Reset()                    { *m = RemoteStateRequest{} }
func (m *RemoteStateRequest) String() string            { return proto.CompactTextString(m) }
func (*RemoteStateRequest) ProtoMessage()               {}
func (*RemoteStateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{25} }

func (m *RemoteStateRequest) GetStartSeqNum() uint64 {
	if m != nil {
		return m.StartSeqNum
	}
	return 0
}

func (m *RemoteStateRequest) GetEndSeqNum() uint64 {
	if m != nil {
		return m.EndSeqNum
	}
	return 0
}



type RemoteStateResponse struct {
	Payloads []*Payload `protobuf:"bytes,1,rep,name=payloads" json:"payloads,omitempty"`
}

func (m *RemoteStateResponse) Reset()                    { *m = RemoteStateResponse{} }
func (m *RemoteStateResponse) String() string            { return proto.CompactTextString(m) }
func (*RemoteStateResponse) ProtoMessage()               {}
func (*RemoteStateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{26} }

func (m *RemoteStateResponse) GetPayloads() []*Payload {
	if m != nil {
		return m.Payloads
	}
	return nil
}



type RemotePvtDataRequest struct {
	Digests []*PvtDataDigest `protobuf:"bytes,1,rep,name=digests" json:"digests,omitempty"`
}

func (m *RemotePvtDataRequest) Reset()                    { *m = RemotePvtDataRequest{} }
func (m *RemotePvtDataRequest) String() string            { return proto.CompactTextString(m) }
func (*RemotePvtDataRequest) ProtoMessage()               {}
func (*RemotePvtDataRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{27} }

func (m *RemotePvtDataRequest) GetDigests() []*PvtDataDigest {
	if m != nil {
		return m.Digests
	}
	return nil
}


type PvtDataDigest struct {
	TxId       string `protobuf:"bytes,1,opt,name=tx_id,json=txId" json:"tx_id,omitempty"`
	Namespace  string `protobuf:"bytes,2,opt,name=namespace" json:"namespace,omitempty"`
	Collection string `protobuf:"bytes,3,opt,name=collection" json:"collection,omitempty"`
	BlockSeq   uint64 `protobuf:"varint,4,opt,name=block_seq,json=blockSeq" json:"block_seq,omitempty"`
	SeqInBlock uint64 `protobuf:"varint,5,opt,name=seq_in_block,json=seqInBlock" json:"seq_in_block,omitempty"`
}

func (m *PvtDataDigest) Reset()                    { *m = PvtDataDigest{} }
func (m *PvtDataDigest) String() string            { return proto.CompactTextString(m) }
func (*PvtDataDigest) ProtoMessage()               {}
func (*PvtDataDigest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{28} }

func (m *PvtDataDigest) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *PvtDataDigest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *PvtDataDigest) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

func (m *PvtDataDigest) GetBlockSeq() uint64 {
	if m != nil {
		return m.BlockSeq
	}
	return 0
}

func (m *PvtDataDigest) GetSeqInBlock() uint64 {
	if m != nil {
		return m.SeqInBlock
	}
	return 0
}



type RemotePvtDataResponse struct {
	Elements []*PvtDataElement `protobuf:"bytes,1,rep,name=elements" json:"elements,omitempty"`
}

func (m *RemotePvtDataResponse) Reset()                    { *m = RemotePvtDataResponse{} }
func (m *RemotePvtDataResponse) String() string            { return proto.CompactTextString(m) }
func (*RemotePvtDataResponse) ProtoMessage()               {}
func (*RemotePvtDataResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{29} }

func (m *RemotePvtDataResponse) GetElements() []*PvtDataElement {
	if m != nil {
		return m.Elements
	}
	return nil
}

type PvtDataElement struct {
	Digest *PvtDataDigest `protobuf:"bytes,1,opt,name=digest" json:"digest,omitempty"`
	
	Payload [][]byte `protobuf:"bytes,2,rep,name=payload,proto3" json:"payload,omitempty"`
}

func (m *PvtDataElement) Reset()                    { *m = PvtDataElement{} }
func (m *PvtDataElement) String() string            { return proto.CompactTextString(m) }
func (*PvtDataElement) ProtoMessage()               {}
func (*PvtDataElement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{30} }

func (m *PvtDataElement) GetDigest() *PvtDataDigest {
	if m != nil {
		return m.Digest
	}
	return nil
}

func (m *PvtDataElement) GetPayload() [][]byte {
	if m != nil {
		return m.Payload
	}
	return nil
}



type PvtDataPayload struct {
	TxSeqInBlock uint64 `protobuf:"varint,1,opt,name=tx_seq_in_block,json=txSeqInBlock" json:"tx_seq_in_block,omitempty"`
	
	
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *PvtDataPayload) Reset()                    { *m = PvtDataPayload{} }
func (m *PvtDataPayload) String() string            { return proto.CompactTextString(m) }
func (*PvtDataPayload) ProtoMessage()               {}
func (*PvtDataPayload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{31} }

func (m *PvtDataPayload) GetTxSeqInBlock() uint64 {
	if m != nil {
		return m.TxSeqInBlock
	}
	return 0
}

func (m *PvtDataPayload) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type Acknowledgement struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
}

func (m *Acknowledgement) Reset()                    { *m = Acknowledgement{} }
func (m *Acknowledgement) String() string            { return proto.CompactTextString(m) }
func (*Acknowledgement) ProtoMessage()               {}
func (*Acknowledgement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{32} }

func (m *Acknowledgement) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}



type Chaincode struct {
	Name     string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Version  string `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
	Metadata []byte `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (m *Chaincode) Reset()                    { *m = Chaincode{} }
func (m *Chaincode) String() string            { return proto.CompactTextString(m) }
func (*Chaincode) ProtoMessage()               {}
func (*Chaincode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{33} }

func (m *Chaincode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Chaincode) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Chaincode) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func init() {
	proto.RegisterType((*Envelope)(nil), "gossip.Envelope")
	proto.RegisterType((*SecretEnvelope)(nil), "gossip.SecretEnvelope")
	proto.RegisterType((*Secret)(nil), "gossip.Secret")
	proto.RegisterType((*GossipMessage)(nil), "gossip.GossipMessage")
	proto.RegisterType((*StateInfo)(nil), "gossip.StateInfo")
	proto.RegisterType((*Properties)(nil), "gossip.Properties")
	proto.RegisterType((*StateInfoSnapshot)(nil), "gossip.StateInfoSnapshot")
	proto.RegisterType((*StateInfoPullRequest)(nil), "gossip.StateInfoPullRequest")
	proto.RegisterType((*ConnEstablish)(nil), "gossip.ConnEstablish")
	proto.RegisterType((*PeerIdentity)(nil), "gossip.PeerIdentity")
	proto.RegisterType((*DataRequest)(nil), "gossip.DataRequest")
	proto.RegisterType((*GossipHello)(nil), "gossip.GossipHello")
	proto.RegisterType((*DataUpdate)(nil), "gossip.DataUpdate")
	proto.RegisterType((*DataDigest)(nil), "gossip.DataDigest")
	proto.RegisterType((*DataMessage)(nil), "gossip.DataMessage")
	proto.RegisterType((*PrivateDataMessage)(nil), "gossip.PrivateDataMessage")
	proto.RegisterType((*Payload)(nil), "gossip.Payload")
	proto.RegisterType((*PrivatePayload)(nil), "gossip.PrivatePayload")
	proto.RegisterType((*AliveMessage)(nil), "gossip.AliveMessage")
	proto.RegisterType((*LeadershipMessage)(nil), "gossip.LeadershipMessage")
	proto.RegisterType((*PeerTime)(nil), "gossip.PeerTime")
	proto.RegisterType((*MembershipRequest)(nil), "gossip.MembershipRequest")
	proto.RegisterType((*MembershipResponse)(nil), "gossip.MembershipResponse")
	proto.RegisterType((*Member)(nil), "gossip.Member")
	proto.RegisterType((*Empty)(nil), "gossip.Empty")
	proto.RegisterType((*RemoteStateRequest)(nil), "gossip.RemoteStateRequest")
	proto.RegisterType((*RemoteStateResponse)(nil), "gossip.RemoteStateResponse")
	proto.RegisterType((*RemotePvtDataRequest)(nil), "gossip.RemotePvtDataRequest")
	proto.RegisterType((*PvtDataDigest)(nil), "gossip.PvtDataDigest")
	proto.RegisterType((*RemotePvtDataResponse)(nil), "gossip.RemotePvtDataResponse")
	proto.RegisterType((*PvtDataElement)(nil), "gossip.PvtDataElement")
	proto.RegisterType((*PvtDataPayload)(nil), "gossip.PvtDataPayload")
	proto.RegisterType((*Acknowledgement)(nil), "gossip.Acknowledgement")
	proto.RegisterType((*Chaincode)(nil), "gossip.Chaincode")
	proto.RegisterEnum("gossip.PullMsgType", PullMsgType_name, PullMsgType_value)
	proto.RegisterEnum("gossip.GossipMessage_Tag", GossipMessage_Tag_name, GossipMessage_Tag_value)
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4



type GossipClient interface {
	
	GossipStream(ctx context.Context, opts ...grpc.CallOption) (Gossip_GossipStreamClient, error)
	
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type gossipClient struct {
	cc *grpc.ClientConn
}

func NewGossipClient(cc *grpc.ClientConn) GossipClient {
	return &gossipClient{cc}
}

func (c *gossipClient) GossipStream(ctx context.Context, opts ...grpc.CallOption) (Gossip_GossipStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Gossip_serviceDesc.Streams[0], c.cc, "/gossip.Gossip/GossipStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &gossipGossipStreamClient{stream}
	return x, nil
}

type Gossip_GossipStreamClient interface {
	Send(*Envelope) error
	Recv() (*Envelope, error)
	grpc.ClientStream
}

type gossipGossipStreamClient struct {
	grpc.ClientStream
}

func (x *gossipGossipStreamClient) Send(m *Envelope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gossipGossipStreamClient) Recv() (*Envelope, error) {
	m := new(Envelope)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gossipClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/gossip.Gossip/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}



type GossipServer interface {
	
	GossipStream(Gossip_GossipStreamServer) error
	
	Ping(context.Context, *Empty) (*Empty, error)
}

func RegisterGossipServer(s *grpc.Server, srv GossipServer) {
	s.RegisterService(&_Gossip_serviceDesc, srv)
}

func _Gossip_GossipStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GossipServer).GossipStream(&gossipGossipStreamServer{stream})
}

type Gossip_GossipStreamServer interface {
	Send(*Envelope) error
	Recv() (*Envelope, error)
	grpc.ServerStream
}

type gossipGossipStreamServer struct {
	grpc.ServerStream
}

func (x *gossipGossipStreamServer) Send(m *Envelope) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gossipGossipStreamServer) Recv() (*Envelope, error) {
	m := new(Envelope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Gossip_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gossip.Gossip/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Gossip_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gossip.Gossip",
	HandlerType: (*GossipServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Gossip_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GossipStream",
			Handler:       _Gossip_GossipStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "gossip/message.proto",
}

func init() { proto.RegisterFile("gossip/message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x58, 0xdd, 0x6f, 0xe3, 0xc6,
	0x11, 0x17, 0x6d, 0x49, 0x16, 0x47, 0x1f, 0x96, 0xd7, 0xbe, 0x3b, 0xc6, 0x49, 0x13, 0x87, 0xed,
	0x25, 0xd7, 0xfa, 0x62, 0x5f, 0x9d, 0x16, 0x0d, 0x90, 0xb6, 0x07, 0x5b, 0x52, 0x2c, 0x21, 0x27,
	0x9d, 0x4b, 0xfb, 0xd0, 0xba, 0x2f, 0xc4, 0x9a, 0x5c, 0x53, 0xac, 0xc9, 0x25, 0xcd, 0x5d, 0x3b,
	0xf6, 0x63, 0xd1, 0x87, 0x00, 0x7d, 0xe9, 0xdf, 0xd0, 0xa7, 0xfe, 0x9b, 0xc5, 0xee, 0xf2, 0x53,
	0xb2, 0x0f, 0xb8, 0x00, 0x79, 0xd3, 0x7c, 0xef, 0xce, 0xce, 0xfc, 0x66, 0x28, 0xd8, 0xf2, 0x22,
	0xc6, 0xfc, 0x78, 0x3f, 0x24, 0x8c, 0x61, 0x8f, 0xec, 0xc5, 0x49, 0xc4, 0x23, 0xd4, 0x54, 0xdc,
	0xed, 0x67, 0x4e, 0x14, 0x86, 0x11, 0xdd, 0x77, 0xa2, 0x20, 0x20, 0x0e, 0xf7, 0x23, 0xaa, 0x14,
	0xcc, 0x7f, 0x69, 0xd0, 0x1a, 0xd1, 0x5b, 0x12, 0x44, 0x31, 0x41, 0x06, 0xac, 0xc5, 0xf8, 0x3e,
	0x88, 0xb0, 0x6b, 0x68, 0x3b, 0xda, 0x8b, 0x8e, 0x95, 0x91, 0xe8, 0x13, 0xd0, 0x99, 0xef, 0x51,
	0xcc, 0x6f, 0x12, 0x62, 0xac, 0x48, 0x59, 0xc1, 0x40, 0xaf, 0x61, 0x9d, 0x11, 0x27, 0x21, 0xdc,
	0x26, 0xa9, 0x2b, 0x63, 0x75, 0x47, 0x7b, 0xd1, 0x3e, 0x78, 0xba, 0xa7, 0xe2, 0xef, 0x9d, 0x4a,
	0x71, 0x16, 0xc8, 0xea, 0xb1, 0x0a, 0x6d, 0x8e, 0xa1, 0x57, 0xd5, 0xf8, 0xa9, 0x47, 0x31, 0x0f,
	0xa1, 0xa9, 0x3c, 0xa1, 0x97, 0xd0, 0xf7, 0x29, 0x27, 0x09, 0xc5, 0xc1, 0x88, 0xba, 0x71, 0xe4,
	0x53, 0x2e, 0x5d, 0xe9, 0xe3, 0x9a, 0xb5, 0x24, 0x39, 0xd2, 0x61, 0xcd, 0x89, 0x28, 0x27, 0x94,
	0x9b, 0x3f, 0xb6, 0xa1, 0x7b, 0x2c, 0x8f, 0x3d, 0x55, 0xb9, 0x44, 0x5b, 0xd0, 0xa0, 0x11, 0x75,
	0x88, 0xb4, 0xaf, 0x5b, 0x8a, 0x10, 0x47, 0x74, 0xe6, 0x98, 0x52, 0x12, 0xa4, 0xc7, 0xc8, 0x48,
	0xb4, 0x0b, 0xab, 0x1c, 0x7b, 0x32, 0x07, 0xbd, 0x83, 0x8f, 0xb2, 0x1c, 0x54, 0x7c, 0xee, 0x9d,
	0x61, 0xcf, 0x12, 0x5a, 0xe8, 0x6b, 0xd0, 0x71, 0xe0, 0xdf, 0x12, 0x3b, 0x64, 0x9e, 0xd1, 0x90,
	0x69, 0xdb, 0xca, 0x4c, 0x0e, 0x85, 0x20, 0xb5, 0x18, 0xd7, 0xac, 0x96, 0x54, 0x9c, 0x32, 0x0f,
	0xfd, 0x0e, 0xd6, 0x42, 0x12, 0xda, 0x09, 0xb9, 0x36, 0x9a, 0xd2, 0x24, 0x8f, 0x32, 0x25, 0xe1,
	0x05, 0x49, 0xd8, 0xdc, 0x8f, 0x2d, 0x72, 0x7d, 0x43, 0x18, 0x1f, 0xd7, 0xac, 0x66, 0x48, 0x42,
	0x8b, 0x5c, 0xa3, 0xdf, 0x67, 0x56, 0xcc, 0x58, 0x93, 0x56, 0xdb, 0x0f, 0x59, 0xb1, 0x38, 0xa2,
	0x8c, 0xe4, 0x66, 0x0c, 0xbd, 0x82, 0x96, 0x8b, 0x39, 0x96, 0x07, 0x6c, 0x49, 0xbb, 0xcd, 0xcc,
	0x6e, 0x88, 0x39, 0x2e, 0xce, 0xb7, 0x26, 0xd4, 0xc4, 0xf1, 0x76, 0xa1, 0x31, 0x27, 0x41, 0x10,
	0x19, 0x7a, 0x55, 0x5d, 0xa5, 0x60, 0x2c, 0x44, 0xe3, 0x9a, 0xa5, 0x74, 0xd0, 0x7e, 0xea, 0xde,
	0xf5, 0x3d, 0x03, 0xa4, 0x3e, 0x2a, 0xbb, 0x1f, 0xfa, 0x9e, 0xba, 0x85, 0xf4, 0x3e, 0xf4, 0xbd,
	0xfc, 0x3c, 0xe2, 0xf6, 0xed, 0xe5, 0xf3, 0x14, 0xf7, 0x96, 0x16, 0xea, 0xe2, 0x6d, 0x69, 0x71,
	0x13, 0xbb, 0x98, 0x13, 0xa3, 0xb3, 0x1c, 0xe5, 0x9d, 0x94, 0x8c, 0x6b, 0x16, 0xb8, 0x39, 0x85,
	0x9e, 0x43, 0x83, 0x84, 0x31, 0xbf, 0x37, 0xba, 0xd2, 0xa0, 0x9b, 0x19, 0x8c, 0x04, 0x53, 0x5c,
	0x40, 0x4a, 0xd1, 0x2e, 0xd4, 0x9d, 0x88, 0x52, 0xa3, 0x27, 0xb5, 0x9e, 0x64, 0x5a, 0x83, 0x88,
	0xd2, 0x11, 0xe3, 0xf8, 0x22, 0xf0, 0xd9, 0x7c, 0x5c, 0xb3, 0xa4, 0x12, 0x3a, 0x00, 0x60, 0x1c,
	0x73, 0x62, 0xfb, 0xf4, 0x32, 0x32, 0xd6, 0xa5, 0xc9, 0x46, 0xde, 0x26, 0x42, 0x32, 0xa1, 0x97,
	0x22, 0x3b, 0x3a, 0xcb, 0x08, 0x74, 0x04, 0x3d, 0x65, 0xc3, 0x28, 0x8e, 0xd9, 0x3c, 0xe2, 0x46,
	0xbf, 0xfa, 0xe8, 0xb9, 0xdd, 0x69, 0xaa, 0x30, 0xae, 0x59, 0x5d, 0x69, 0x92, 0x31, 0xd0, 0x14,
	0x36, 0x8b, 0xb8, 0x76, 0x7c, 0x13, 0x04, 0x32, 0x7f, 0x1b, 0xd2, 0xd1, 0x27, 0x4b, 0x8e, 0x4e,
	0x6e, 0x82, 0xa0, 0x48, 0x64, 0x9f, 0x2d, 0xf0, 0xd1, 0x21, 0x28, 0xff, 0xc2, 0x89, 0x50, 0x32,
	0x50, 0xb5, 0xa0, 0x2c, 0x12, 0x46, 0x9c, 0x48, 0x77, 0x85, 0x9b, 0x0e, 0x2b, 0xd1, 0x68, 0x98,
	0xdd, 0x2a, 0x49, 0x4b, 0xce, 0xd8, 0x94, 0x3e, 0x3e, 0x7e, 0xd0, 0x47, 0x5e, 0x95, 0x5d, 0x56,
	0x66, 0x88, 0xdc, 0x04, 0x04, 0xbb, 0xaa, 0x78, 0x65, 0x89, 0x6e, 0x55, 0x73, 0xf3, 0x26, 0x97,
	0x16, 0x85, 0xda, 0x2d, 0x4c, 0x44, 0xb9, 0x7e, 0x0b, 0xdd, 0x98, 0x90, 0xc4, 0xf6, 0x5d, 0x42,
	0xb9, 0xcf, 0xef, 0x8d, 0x27, 0xd5, 0x36, 0x3c, 0x21, 0x24, 0x99, 0xa4, 0x32, 0x71, 0x8d, 0xb8,
	0x44, 0x8b, 0x66, 0xc7, 0xce, 0x95, 0xf1, 0x54, 0x9a, 0x3c, 0xcb, 0x3b, 0xd7, 0xb9, 0xa2, 0xd1,
	0x0f, 0x01, 0x71, 0x3d, 0x12, 0x12, 0x2a, 0x2e, 0x2f, 0xb4, 0xd0, 0x9f, 0x01, 0xe2, 0xc4, 0xbf,
	0x55, 0x59, 0x30, 0x9e, 0x55, 0x93, 0xaf, 0xee, 0x7b, 0x72, 0xcb, 0xab, 0x55, 0x5c, 0xb2, 0x40,
	0xaf, 0x4b, 0xf6, 0xcc, 0x30, 0xa4, 0xfd, 0x2f, 0x1e, 0xb1, 0xcf, 0x33, 0x56, 0x32, 0x41, 0xaf,
	0xa1, 0x93, 0x52, 0xb6, 0x28, 0x74, 0xe3, 0xa3, 0xea, 0xb3, 0x9d, 0x28, 0x59, 0xb5, 0xad, 0xdb,
	0x71, 0xc1, 0x35, 0x6d, 0x58, 0x3d, 0xc3, 0x1e, 0xea, 0x82, 0xfe, 0x6e, 0x36, 0x1c, 0x7d, 0x37,
	0x99, 0x8d, 0x86, 0xfd, 0x1a, 0xd2, 0xa1, 0x31, 0x9a, 0x9e, 0x9c, 0x9d, 0xf7, 0x35, 0xd4, 0x81,
	0xd6, 0x5b, 0xeb, 0xd8, 0x7e, 0x3b, 0x7b, 0x73, 0xde, 0x5f, 0x11, 0x7a, 0x83, 0xf1, 0xe1, 0x4c,
	0x91, 0xab, 0xa8, 0x0f, 0x1d, 0x49, 0x1e, 0xce, 0x86, 0xf6, 0x5b, 0xeb, 0xb8, 0x5f, 0x47, 0xeb,
	0xd0, 0x56, 0x0a, 0x96, 0x64, 0x34, 0xca, 0x48, 0xfc, 0x3f, 0x0d, 0xf4, 0xbc, 0x22, 0xd1, 0x1e,
	0xe8, 0xdc, 0x0f, 0x09, 0xe3, 0x38, 0x8c, 0x25, 0xe2, 0xb6, 0x0f, 0xfa, 0xe5, 0x17, 0x3a, 0xf3,
	0x43, 0x62, 0x15, 0x2a, 0xe8, 0x09, 0x34, 0xe3, 0x2b, 0xdf, 0xf6, 0x5d, 0x09, 0xc4, 0x1d, 0xab,
	0x11, 0x5f, 0xf9, 0x13, 0x17, 0x7d, 0x06, 0xed, 0x14, 0xa7, 0xed, 0xe9, 0xe1, 0xc0, 0xa8, 0x4b,
	0x19, 0xa4, 0xac, 0xe9, 0xe1, 0x40, 0x74, 0x68, 0x9c, 0x44, 0x31, 0x49, 0xb8, 0x4f, 0x58, 0x8a,
	0xc8, 0xa8, 0x48, 0x50, 0x26, 0xb1, 0x4a, 0x5a, 0xe6, 0x8f, 0x1a, 0x40, 0x21, 0x42, 0xbf, 0x84,
	0xae, 0x7c, 0xfa, 0xc4, 0x9e, 0x13, 0xdf, 0x9b, 0xf3, 0x74, 0x70, 0x74, 0x14, 0x73, 0x2c, 0x79,
	0xe8, 0x73, 0xe8, 0x04, 0xe4, 0x92, 0xdb, 0xe5, 0x21, 0xd2, 0xb2, 0xda, 0x82, 0x37, 0x48, 0x07,
	0xc9, 0x6f, 0x41, 0x1c, 0xcc, 0xa7, 0x4e, 0xe4, 0x12, 0x66, 0xac, 0xee, 0xac, 0x96, 0xc1, 0x62,
	0x90, 0x49, 0xac, 0x92, 0x92, 0x79, 0x08, 0x1b, 0x4b, 0x68, 0x80, 0x5e, 0x42, 0x8b, 0x04, 0xb2,
	0x10, 0x99, 0xa1, 0x49, 0x2f, 0x79, 0xe6, 0xf2, 0x99, 0x9c, 0x6b, 0x98, 0x7f, 0x80, 0xad, 0x87,
	0x70, 0x60, 0x31, 0x73, 0xda, 0x62, 0xe6, 0xcc, 0x4b, 0xe8, 0x56, 0x40, 0xaf, 0xf4, 0x04, 0x5a,
	0xf9, 0x09, 0xb6, 0xa1, 0x95, 0xb7, 0x9a, 0x1a, 0x9d, 0x39, 0x8d, 0x4c, 0xe8, 0xf2, 0x80, 0xd9,
	0x0e, 0x49, 0xb8, 0x3d, 0xc7, 0x6c, 0x9e, 0x3e, 0x5e, 0x9b, 0x07, 0x6c, 0x40, 0x12, 0x3e, 0xc6,
	0x6c, 0x6e, 0xbe, 0x83, 0x4e, 0xb9, 0x25, 0x1f, 0x0b, 0x83, 0xa0, 0x2e, 0xdc, 0xa4, 0x21, 0xe4,
	0x6f, 0x11, 0x3a, 0x24, 0x1c, 0xcb, 0xda, 0x57, 0x9e, 0x73, 0xda, 0x0c, 0xa1, 0x5d, 0xea, 0xbc,
	0xc7, 0xa7, 0xbe, 0x2b, 0x27, 0x12, 0x33, 0x56, 0x76, 0x56, 0x5f, 0xe8, 0x56, 0x46, 0xa2, 0x3d,
	0x68, 0x85, 0xcc, 0xb3, 0xf9, 0x7d, 0xba, 0xfe, 0xf4, 0x8a, 0xb1, 0x24, 0xb2, 0x38, 0x65, 0xde,
	0xd9, 0x7d, 0x4c, 0xac, 0xb5, 0x50, 0xfd, 0x30, 0x23, 0x68, 0x97, 0xe6, 0xe1, 0x23, 0xe1, 0xca,
	0xe7, 0x5d, 0xa9, 0x9e, 0xf7, 0x83, 0x03, 0xde, 0x01, 0x14, 0xa3, 0xee, 0x91, 0x78, 0xbf, 0x82,
	0x7a, 0x1a, 0xeb, 0xe1, 0x2a, 0xa9, 0xff, 0xa4, 0xc8, 0x81, 0x8a, 0xac, 0x46, 0xf9, 0xcf, 0x9e,
	0xd8, 0x6f, 0xd4, 0x3b, 0x66, 0xdb, 0xdb, 0xaf, 0xab, 0xab, 0x64, 0xfb, 0x60, 0x3d, 0xb7, 0x56,
	0xec, 0x7c, 0xb7, 0x34, 0xbf, 0x03, 0xb4, 0x8c, 0x80, 0xe8, 0xd5, 0xa2, 0x83, 0xa7, 0x0b, 0x70,
	0xb9, 0xe4, 0xe7, 0x1c, 0xd6, 0x52, 0x1e, 0x7a, 0x06, 0x6b, 0x8c, 0x5c, 0xdb, 0xf4, 0x26, 0x4c,
	0xaf, 0xdb, 0x64, 0xe4, 0x7a, 0x76, 0x13, 0x8a, 0xea, 0x2c, 0xbd, 0xaa, 0xca, 0xeb, 0xe7, 0x0b,
	0xe8, 0x2c, 0x3a, 0xbe, 0x53, 0xc5, 0xdf, 0xff, 0xac, 0x40, 0xaf, 0x1a, 0x16, 0x7d, 0x09, 0xeb,
	0xc5, 0x5e, 0x6f, 0x53, 0x1c, 0xaa, 0xcc, 0xea, 0x56, 0xaf, 0x60, 0xcf, 0x70, 0x48, 0xc4, 0xea,
	0x2c, 0xa4, 0x2c, 0xc6, 0x8e, 0x5a, 0x9d, 0x75, 0xab, 0x60, 0xa0, 0x4d, 0x68, 0xf0, 0xbb, 0x0c,
	0x2e, 0x75, 0xab, 0xce, 0xef, 0x26, 0xae, 0x40, 0xb2, 0xec, 0x44, 0xc9, 0x0f, 0x8c, 0xf0, 0x14,
	0x2f, 0xb3, 0x63, 0x5a, 0x82, 0x87, 0x5e, 0x02, 0xca, 0x94, 0x98, 0x1f, 0x66, 0x98, 0xd7, 0x90,
	0xd7, 0xed, 0xa7, 0x92, 0x53, 0x3f, 0x4c, 0x71, 0x6f, 0x06, 0xa8, 0x74, 0x5c, 0x27, 0xa2, 0x97,
	0xbe, 0xc7, 0xd2, 0x35, 0xf6, 0xb3, 0x3d, 0xf5, 0xa1, 0xb2, 0x37, 0xc8, 0x35, 0x06, 0x52, 0xe1,
	0x04, 0x3b, 0x57, 0xd8, 0x23, 0xd6, 0x86, 0xb3, 0x20, 0x60, 0xe6, 0xbf, 0x35, 0xe8, 0x94, 0x17,
	0x65, 0xb4, 0x07, 0x10, 0xe6, 0xfb, 0x6c, 0xfa, 0x64, 0xbd, 0xea, 0xa6, 0x6b, 0x95, 0x34, 0x3e,
	0x78, 0xb0, 0x94, 0xe1, 0xab, 0x5e, 0x85, 0x2f, 0xf3, 0x9f, 0x1a, 0x6c, 0x2c, 0x6d, 0x1c, 0x8f,
	0x01, 0xd4, 0x87, 0x06, 0x7e, 0x0e, 0x3d, 0x9f, 0xd9, 0x2e, 0x71, 0x02, 0x9c, 0x60, 0x91, 0x02,
	0xf9, 0x54, 0x2d, 0xab, 0xeb, 0xb3, 0x61, 0xc1, 0x34, 0xff, 0x08, 0xad, 0xcc, 0x5a, 0x94, 0x9f,
	0x4f, 0x9d, 0x72, 0xf9, 0xf9, 0xd4, 0x11, 0xe5, 0x57, 0xaa, 0xcb, 0x95, 0x72, 0x5d, 0x9a, 0x97,
	0xb0, 0xb1, 0xf4, 0x0d, 0x81, 0xbe, 0x85, 0x3e, 0x23, 0xc1, 0xa5, 0x5c, 0x1e, 0x93, 0x50, 0xc5,
	0xd6, 0xaa, 0x07, 0xce, 0x21, 0x62, 0x5d, 0x68, 0x4e, 0x0a, 0x45, 0xd1, 0xef, 0x62, 0x19, 0xa2,
	0xb2, 0xaf, 0x3b, 0x96, 0x22, 0xcc, 0x0b, 0x40, 0xcb, 0x5f, 0x1d, 0xe8, 0x0b, 0x68, 0xc8, 0x8f,
	0x9c, 0x47, 0xc7, 0x94, 0x12, 0x4b, 0x9c, 0x22, 0xd8, 0x7d, 0x0f, 0x4e, 0x11, 0xec, 0x9a, 0x7f,
	0x85, 0xa6, 0x8a, 0x21, 0xde, 0x8c, 0x54, 0xbe, 0x02, 0xad, 0x9c, 0x7e, 0x2f, 0xc6, 0x3e, 0xbc,
	0x44, 0x98, 0x6b, 0xd0, 0x90, 0x1f, 0x01, 0xe6, 0xdf, 0x00, 0x2d, 0xaf, 0xba, 0x62, 0x88, 0x31,
	0x8e, 0x13, 0x6e, 0x57, 0x5b, 0xbf, 0x2d, 0x99, 0xa7, 0xaa, 0xff, 0x3f, 0x85, 0x36, 0xa1, 0xae,
	0x5d, 0x7d, 0x04, 0x9d, 0x50, 0x57, 0xc9, 0xcd, 0x23, 0xd8, 0x7c, 0x60, 0x01, 0x46, 0xbb, 0xd0,
	0x4a, 0x51, 0x26, 0x1b, 0xe5, 0x4b, 0x70, 0x96, 0x2b, 0x98, 0xc7, 0xb0, 0xf5, 0xd0, 0x52, 0x89,
	0xf6, 0x0b, 0xac, 0x55, 0x3e, 0xf2, 0x8f, 0x96, 0x54, 0x51, 0x21, 0x75, 0x0e, 0xc1, 0xe6, 0x7f,
	0x35, 0xe8, 0x56, 0x44, 0x05, 0x5a, 0x68, 0x25, 0xb4, 0x78, 0x3f, 0xc0, 0x7c, 0x0a, 0x50, 0x74,
	0x6f, 0x8a, 0x32, 0x25, 0x0e, 0xfa, 0x18, 0xf4, 0x8b, 0x20, 0x72, 0xae, 0x44, 0x4e, 0x64, 0x63,
	0xd5, 0xad, 0x96, 0x64, 0x9c, 0x92, 0x6b, 0xb4, 0x03, 0x1d, 0x91, 0x2a, 0x9f, 0xda, 0x92, 0x95,
	0xa2, 0x0b, 0x30, 0x72, 0x3d, 0xa1, 0x47, 0x82, 0x63, 0x7e, 0x0f, 0x4f, 0x1e, 0xdc, 0x80, 0xd1,
	0xc1, 0xd2, 0xf6, 0xf3, 0x74, 0xe1, 0xba, 0x23, 0x25, 0x2e, 0xed, 0x40, 0xe7, 0xd0, 0xab, 0xca,
	0xd0, 0x57, 0xd0, 0x54, 0xd9, 0x48, 0x0b, 0xff, 0x91, 0x94, 0xa5, 0x4a, 0xe5, 0x3f, 0x30, 0x54,
	0xd9, 0xe7, 0xc3, 0xe1, 0x2f, 0xb9, 0xeb, 0x0c, 0xc0, 0x9f, 0xc3, 0x3a, 0xbf, 0xb3, 0x2b, 0xd7,
	0x4b, 0x17, 0x46, 0x7e, 0x77, 0x9a, 0x5f, 0xb0, 0xea, 0xb2, 0xfc, 0x9f, 0x88, 0xf9, 0x25, 0xac,
	0x2f, 0x7c, 0x70, 0x88, 0xa6, 0x23, 0x49, 0x12, 0x25, 0xe9, 0xfb, 0x28, 0xc2, 0x7c, 0x07, 0x7a,
	0xbe, 0x36, 0x8a, 0x09, 0x54, 0x1a, 0x16, 0xf2, 0xb7, 0x88, 0x71, 0x4b, 0x12, 0x26, 0x1e, 0x48,
	0xbd, 0x5f, 0x46, 0xbe, 0x6f, 0x73, 0xfa, 0xcd, 0x9f, 0xa0, 0x5d, 0x9a, 0xc4, 0x8b, 0x1f, 0x07,
	0x5d, 0xd0, 0x8f, 0xde, 0xbc, 0x1d, 0x7c, 0x6f, 0x4f, 0x4f, 0x8f, 0xfb, 0x9a, 0xf8, 0x06, 0x98,
	0x0c, 0x47, 0xb3, 0xb3, 0xc9, 0xd9, 0xb9, 0xe4, 0xac, 0x1c, 0xfc, 0x03, 0x9a, 0x6a, 0x13, 0x42,
	0xdf, 0x40, 0x47, 0xfd, 0x3a, 0xe5, 0x09, 0xc1, 0x21, 0x5a, 0x6a, 0xec, 0xed, 0x25, 0x8e, 0x59,
	0x7b, 0xa1, 0xbd, 0xd2, 0xd0, 0x17, 0x50, 0x3f, 0xf1, 0xa9, 0x87, 0xaa, 0x1f, 0xe9, 0xdb, 0x55,
	0xd2, 0xac, 0x1d, 0x7d, 0xf5, 0xf7, 0x5d, 0xcf, 0xe7, 0xf3, 0x9b, 0x0b, 0x31, 0x69, 0xf6, 0xe7,
	0xf7, 0x31, 0x49, 0xd4, 0x56, 0xbe, 0x7f, 0x89, 0x2f, 0x12, 0xdf, 0xd9, 0x97, 0xff, 0x8b, 0xb1,
	0x7d, 0x65, 0x76, 0xd1, 0x94, 0xe4, 0xd7, 0xff, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x22, 0x0e, 0xc2,
	0x6f, 0x5f, 0x13, 0x00, 0x00,
}
