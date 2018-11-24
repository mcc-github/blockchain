


package etcdraft 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



type Metadata struct {
	Consenters           []*Consenter `protobuf:"bytes,1,rep,name=consenters,proto3" json:"consenters,omitempty"`
	Options              *Options     `protobuf:"bytes,2,opt,name=options,proto3" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_7f128e6e8996db94, []int{0}
}
func (m *Metadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metadata.Unmarshal(m, b)
}
func (m *Metadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metadata.Marshal(b, m, deterministic)
}
func (dst *Metadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metadata.Merge(dst, src)
}
func (m *Metadata) XXX_Size() int {
	return xxx_messageInfo_Metadata.Size(m)
}
func (m *Metadata) XXX_DiscardUnknown() {
	xxx_messageInfo_Metadata.DiscardUnknown(m)
}

var xxx_messageInfo_Metadata proto.InternalMessageInfo

func (m *Metadata) GetConsenters() []*Consenter {
	if m != nil {
		return m.Consenters
	}
	return nil
}

func (m *Metadata) GetOptions() *Options {
	if m != nil {
		return m.Options
	}
	return nil
}


type Consenter struct {
	Host                 string   `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	ClientTlsCert        []byte   `protobuf:"bytes,3,opt,name=client_tls_cert,json=clientTlsCert,proto3" json:"client_tls_cert,omitempty"`
	ServerTlsCert        []byte   `protobuf:"bytes,4,opt,name=server_tls_cert,json=serverTlsCert,proto3" json:"server_tls_cert,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Consenter) Reset()         { *m = Consenter{} }
func (m *Consenter) String() string { return proto.CompactTextString(m) }
func (*Consenter) ProtoMessage()    {}
func (*Consenter) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_7f128e6e8996db94, []int{1}
}
func (m *Consenter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Consenter.Unmarshal(m, b)
}
func (m *Consenter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Consenter.Marshal(b, m, deterministic)
}
func (dst *Consenter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Consenter.Merge(dst, src)
}
func (m *Consenter) XXX_Size() int {
	return xxx_messageInfo_Consenter.Size(m)
}
func (m *Consenter) XXX_DiscardUnknown() {
	xxx_messageInfo_Consenter.DiscardUnknown(m)
}

var xxx_messageInfo_Consenter proto.InternalMessageInfo

func (m *Consenter) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Consenter) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Consenter) GetClientTlsCert() []byte {
	if m != nil {
		return m.ClientTlsCert
	}
	return nil
}

func (m *Consenter) GetServerTlsCert() []byte {
	if m != nil {
		return m.ServerTlsCert
	}
	return nil
}



type Options struct {
	TickInterval         uint64   `protobuf:"varint,1,opt,name=tick_interval,json=tickInterval,proto3" json:"tick_interval,omitempty"`
	ElectionTick         uint32   `protobuf:"varint,2,opt,name=election_tick,json=electionTick,proto3" json:"election_tick,omitempty"`
	HeartbeatTick        uint32   `protobuf:"varint,3,opt,name=heartbeat_tick,json=heartbeatTick,proto3" json:"heartbeat_tick,omitempty"`
	MaxInflightMsgs      uint32   `protobuf:"varint,4,opt,name=max_inflight_msgs,json=maxInflightMsgs,proto3" json:"max_inflight_msgs,omitempty"`
	MaxSizePerMsg        uint64   `protobuf:"varint,5,opt,name=max_size_per_msg,json=maxSizePerMsg,proto3" json:"max_size_per_msg,omitempty"`
	SnapshotInterval     uint64   `protobuf:"varint,6,opt,name=snapshot_interval,json=snapshotInterval,proto3" json:"snapshot_interval,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Options) Reset()         { *m = Options{} }
func (m *Options) String() string { return proto.CompactTextString(m) }
func (*Options) ProtoMessage()    {}
func (*Options) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_7f128e6e8996db94, []int{2}
}
func (m *Options) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Options.Unmarshal(m, b)
}
func (m *Options) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Options.Marshal(b, m, deterministic)
}
func (dst *Options) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Options.Merge(dst, src)
}
func (m *Options) XXX_Size() int {
	return xxx_messageInfo_Options.Size(m)
}
func (m *Options) XXX_DiscardUnknown() {
	xxx_messageInfo_Options.DiscardUnknown(m)
}

var xxx_messageInfo_Options proto.InternalMessageInfo

func (m *Options) GetTickInterval() uint64 {
	if m != nil {
		return m.TickInterval
	}
	return 0
}

func (m *Options) GetElectionTick() uint32 {
	if m != nil {
		return m.ElectionTick
	}
	return 0
}

func (m *Options) GetHeartbeatTick() uint32 {
	if m != nil {
		return m.HeartbeatTick
	}
	return 0
}

func (m *Options) GetMaxInflightMsgs() uint32 {
	if m != nil {
		return m.MaxInflightMsgs
	}
	return 0
}

func (m *Options) GetMaxSizePerMsg() uint64 {
	if m != nil {
		return m.MaxSizePerMsg
	}
	return 0
}

func (m *Options) GetSnapshotInterval() uint64 {
	if m != nil {
		return m.SnapshotInterval
	}
	return 0
}




type RaftMetadata struct {
	
	
	Consenters map[uint64]*Consenter `protobuf:"bytes,1,rep,name=consenters,proto3" json:"consenters,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	
	
	NextConsenterID uint64 `protobuf:"varint,2,opt,name=nextConsenterID,proto3" json:"nextConsenterID,omitempty"`
	
	RaftIndex            uint64   `protobuf:"varint,3,opt,name=raftIndex,proto3" json:"raftIndex,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RaftMetadata) Reset()         { *m = RaftMetadata{} }
func (m *RaftMetadata) String() string { return proto.CompactTextString(m) }
func (*RaftMetadata) ProtoMessage()    {}
func (*RaftMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_7f128e6e8996db94, []int{3}
}
func (m *RaftMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RaftMetadata.Unmarshal(m, b)
}
func (m *RaftMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RaftMetadata.Marshal(b, m, deterministic)
}
func (dst *RaftMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RaftMetadata.Merge(dst, src)
}
func (m *RaftMetadata) XXX_Size() int {
	return xxx_messageInfo_RaftMetadata.Size(m)
}
func (m *RaftMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_RaftMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_RaftMetadata proto.InternalMessageInfo

func (m *RaftMetadata) GetConsenters() map[uint64]*Consenter {
	if m != nil {
		return m.Consenters
	}
	return nil
}

func (m *RaftMetadata) GetNextConsenterID() uint64 {
	if m != nil {
		return m.NextConsenterID
	}
	return 0
}

func (m *RaftMetadata) GetRaftIndex() uint64 {
	if m != nil {
		return m.RaftIndex
	}
	return 0
}

func init() {
	proto.RegisterType((*Metadata)(nil), "etcdraft.Metadata")
	proto.RegisterType((*Consenter)(nil), "etcdraft.Consenter")
	proto.RegisterType((*Options)(nil), "etcdraft.Options")
	proto.RegisterType((*RaftMetadata)(nil), "etcdraft.RaftMetadata")
	proto.RegisterMapType((map[uint64]*Consenter)(nil), "etcdraft.RaftMetadata.ConsentersEntry")
}

func init() {
	proto.RegisterFile("orderer/etcdraft/configuration.proto", fileDescriptor_configuration_7f128e6e8996db94)
}

var fileDescriptor_configuration_7f128e6e8996db94 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0x4f, 0x8f, 0xd3, 0x30,
	0x10, 0xc5, 0x95, 0x6d, 0xf7, 0xdf, 0x6c, 0x43, 0x5b, 0x73, 0xa9, 0x10, 0x87, 0xaa, 0xc0, 0x52,
	0x58, 0x29, 0x91, 0x76, 0x85, 0x84, 0x38, 0xb2, 0x80, 0xd4, 0x43, 0x05, 0x32, 0x7b, 0xe2, 0x12,
	0xb9, 0xe9, 0x34, 0xb1, 0x9a, 0xc6, 0x91, 0x3d, 0xad, 0xda, 0xbd, 0xf2, 0x8d, 0x39, 0x73, 0x40,
	0x76, 0x92, 0xa6, 0x54, 0x7b, 0xb3, 0xde, 0xfb, 0xcd, 0xe4, 0xc5, 0x33, 0x86, 0xd7, 0x4a, 0xcf,
	0x51, 0xa3, 0x0e, 0x91, 0xe2, 0xb9, 0x16, 0x0b, 0x0a, 0x63, 0x95, 0x2f, 0x64, 0xb2, 0xd6, 0x82,
	0xa4, 0xca, 0x83, 0x42, 0x2b, 0x52, 0xec, 0xa2, 0x76, 0x47, 0x19, 0x5c, 0x4c, 0x91, 0xc4, 0x5c,
	0x90, 0x60, 0x77, 0x00, 0xb1, 0xca, 0x0d, 0xe6, 0x84, 0xda, 0x0c, 0xbc, 0x61, 0x6b, 0x7c, 0x75,
	0xfb, 0x3c, 0xa8, 0xd1, 0xe0, 0xbe, 0xf6, 0xf8, 0x01, 0xc6, 0x6e, 0xe0, 0x5c, 0x15, 0xb6, 0xb5,
	0x19, 0x9c, 0x0c, 0xbd, 0xf1, 0xd5, 0x6d, 0xbf, 0xa9, 0xf8, 0x5e, 0x1a, 0xbc, 0x26, 0x46, 0xbf,
	0x3d, 0xb8, 0xdc, 0xb7, 0x61, 0x0c, 0xda, 0xa9, 0x32, 0x34, 0xf0, 0x86, 0xde, 0xf8, 0x92, 0xbb,
	0xb3, 0xd5, 0x0a, 0xa5, 0xc9, 0xf5, 0xf2, 0xb9, 0x3b, 0xb3, 0x6b, 0xe8, 0xc6, 0x99, 0xc4, 0x9c,
	0x22, 0xca, 0x4c, 0x14, 0xa3, 0xa6, 0x41, 0x6b, 0xe8, 0x8d, 0x3b, 0xdc, 0x2f, 0xe5, 0x87, 0xcc,
	0xdc, 0x63, 0xc9, 0x19, 0xd4, 0x1b, 0xd4, 0x0d, 0xd7, 0x2e, 0xb9, 0x52, 0xae, 0xb8, 0xd1, 0x5f,
	0x0f, 0xce, 0xab, 0x68, 0xec, 0x15, 0xf8, 0x24, 0xe3, 0x65, 0x24, 0x6d, 0xa2, 0x8d, 0xc8, 0x5c,
	0x98, 0x36, 0xef, 0x58, 0x71, 0x52, 0x69, 0x16, 0xc2, 0x0c, 0x63, 0x5b, 0x11, 0x59, 0xa3, 0x4a,
	0xd7, 0xa9, 0xc5, 0x07, 0x19, 0x2f, 0xd9, 0x1b, 0x78, 0x96, 0xa2, 0xd0, 0x34, 0x43, 0x41, 0x25,
	0xd5, 0x72, 0x94, 0xbf, 0x57, 0x1d, 0xf6, 0x1e, 0xfa, 0x2b, 0xb1, 0x8d, 0x64, 0xbe, 0xc8, 0x64,
	0x92, 0x52, 0xb4, 0x32, 0x89, 0x71, 0x31, 0x7d, 0xde, 0x5d, 0x89, 0xed, 0xa4, 0xd2, 0xa7, 0x26,
	0x31, 0xec, 0x2d, 0xf4, 0x2c, 0x6b, 0xe4, 0x23, 0x46, 0x05, 0x6a, 0xcb, 0x0e, 0x4e, 0x5d, 0x3e,
	0x7f, 0x25, 0xb6, 0x3f, 0xe5, 0x23, 0xfe, 0x40, 0x3d, 0x35, 0x09, 0xbb, 0x81, 0xbe, 0xc9, 0x45,
	0x61, 0x52, 0x45, 0xcd, 0x9f, 0x9c, 0x39, 0xb2, 0x57, 0x1b, 0xf5, 0xdf, 0x8c, 0xfe, 0x78, 0xd0,
	0xe1, 0x62, 0x41, 0xfb, 0xb9, 0x7f, 0x7b, 0x62, 0xee, 0xd7, 0xcd, 0x14, 0x0f, 0xd9, 0x66, 0x09,
	0xcc, 0xd7, 0x9c, 0xf4, 0xee, 0xbf, 0x55, 0x18, 0x43, 0x37, 0xc7, 0x2d, 0xed, 0x91, 0xc9, 0x17,
	0x77, 0x51, 0x6d, 0x7e, 0x2c, 0xb3, 0x97, 0x70, 0x69, 0x5b, 0x4f, 0xf2, 0x39, 0x6e, 0xdd, 0x35,
	0xb5, 0x79, 0x23, 0xbc, 0xe0, 0xd0, 0x3d, 0xfa, 0x0c, 0xeb, 0x41, 0x6b, 0x89, 0xbb, 0x6a, 0x38,
	0xf6, 0xc8, 0xde, 0xc1, 0xe9, 0x46, 0x64, 0x6b, 0xac, 0xb6, 0xee, 0xc9, 0x3d, 0x2d, 0x89, 0x4f,
	0x27, 0x1f, 0xbd, 0xcf, 0x09, 0x04, 0x4a, 0x27, 0x41, 0xba, 0x2b, 0x50, 0x67, 0x38, 0x4f, 0x50,
	0x07, 0x0b, 0x31, 0xd3, 0x32, 0x2e, 0x5f, 0x84, 0x09, 0xaa, 0x77, 0xb3, 0x6f, 0xf3, 0xeb, 0x43,
	0x22, 0x29, 0x5d, 0xcf, 0x82, 0x58, 0xad, 0xc2, 0x83, 0xb2, 0xb0, 0x2c, 0x0b, 0xcb, 0xb2, 0xf0,
	0xf8, 0xb9, 0xcd, 0xce, 0x9c, 0x71, 0xf7, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x4b, 0x2a, 0x56,
	0x89, 0x03, 0x00, 0x00,
}
