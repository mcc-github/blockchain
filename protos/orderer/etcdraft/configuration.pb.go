


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
	return fileDescriptor_configuration_ba740eac6f67e7c7, []int{0}
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
	return fileDescriptor_configuration_ba740eac6f67e7c7, []int{1}
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
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Options) Reset()         { *m = Options{} }
func (m *Options) String() string { return proto.CompactTextString(m) }
func (*Options) ProtoMessage()    {}
func (*Options) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_ba740eac6f67e7c7, []int{2}
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
	return fileDescriptor_configuration_ba740eac6f67e7c7, []int{3}
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
	proto.RegisterFile("orderer/etcdraft/configuration.proto", fileDescriptor_configuration_ba740eac6f67e7c7)
}

var fileDescriptor_configuration_ba740eac6f67e7c7 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0xc1, 0x6e, 0x13, 0x31,
	0x10, 0x86, 0xb5, 0x4d, 0x4a, 0x9b, 0x69, 0x96, 0xb4, 0xe6, 0x12, 0x21, 0x0e, 0x51, 0x80, 0x12,
	0x40, 0xda, 0x95, 0x5a, 0x21, 0x21, 0x8e, 0x14, 0x90, 0x72, 0x88, 0x40, 0xa6, 0x27, 0x2e, 0x2b,
	0xc7, 0x99, 0x6c, 0xac, 0x6c, 0xd6, 0x91, 0x3d, 0x89, 0x92, 0x5e, 0x79, 0x40, 0xde, 0x85, 0x27,
	0x40, 0xb6, 0x37, 0xd9, 0x10, 0xf5, 0x66, 0xfd, 0xf3, 0xcd, 0xec, 0xef, 0xfd, 0xc7, 0xf0, 0x4a,
	0x9b, 0x09, 0x1a, 0x34, 0x29, 0x92, 0x9c, 0x18, 0x31, 0xa5, 0x54, 0xea, 0x72, 0xaa, 0xf2, 0x95,
	0x11, 0xa4, 0x74, 0x99, 0x2c, 0x8d, 0x26, 0xcd, 0xce, 0x77, 0xd5, 0x7e, 0x01, 0xe7, 0x23, 0x24,
	0x31, 0x11, 0x24, 0xd8, 0x2d, 0x80, 0xd4, 0xa5, 0xc5, 0x92, 0xd0, 0xd8, 0x6e, 0xd4, 0x6b, 0x0c,
	0x2e, 0x6e, 0x9e, 0x25, 0x3b, 0x34, 0xb9, 0xdb, 0xd5, 0xf8, 0x01, 0xc6, 0xde, 0xc3, 0x99, 0x5e,
	0xba, 0xd1, 0xb6, 0x7b, 0xd2, 0x8b, 0x06, 0x17, 0x37, 0x57, 0x75, 0xc7, 0xf7, 0x50, 0xe0, 0x3b,
	0xa2, 0xff, 0x3b, 0x82, 0xd6, 0x7e, 0x0c, 0x63, 0xd0, 0x9c, 0x69, 0x4b, 0xdd, 0xa8, 0x17, 0x0d,
	0x5a, 0xdc, 0x9f, 0x9d, 0xb6, 0xd4, 0x86, 0xfc, 0xac, 0x98, 0xfb, 0x33, 0xbb, 0x86, 0x8e, 0x2c,
	0x14, 0x96, 0x94, 0x51, 0x61, 0x33, 0x89, 0x86, 0xba, 0x8d, 0x5e, 0x34, 0x68, 0xf3, 0x38, 0xc8,
	0xf7, 0x85, 0xbd, 0xc3, 0xc0, 0x59, 0x34, 0x6b, 0x34, 0x35, 0xd7, 0x0c, 0x5c, 0x90, 0x2b, 0xae,
	0xff, 0x27, 0x82, 0xb3, 0xca, 0x1a, 0x7b, 0x09, 0x31, 0x29, 0x39, 0xcf, 0x94, 0x73, 0xb4, 0x16,
	0x85, 0x37, 0xd3, 0xe4, 0x6d, 0x27, 0x0e, 0x2b, 0xcd, 0x41, 0x58, 0xa0, 0x74, 0x1d, 0x99, 0x2b,
	0x54, 0xee, 0xda, 0x3b, 0xf1, 0x5e, 0xc9, 0x39, 0x7b, 0x0d, 0x4f, 0x67, 0x28, 0x0c, 0x8d, 0x51,
	0x50, 0xa0, 0x1a, 0x9e, 0x8a, 0xf7, 0xaa, 0xc7, 0xde, 0xc1, 0xd5, 0x42, 0x6c, 0x32, 0x55, 0x4e,
	0x0b, 0x95, 0xcf, 0x28, 0x5b, 0xd8, 0xdc, 0x7a, 0x9b, 0x31, 0xef, 0x2c, 0xc4, 0x66, 0x58, 0xe9,
	0x23, 0x9b, 0x5b, 0xf6, 0x06, 0x2e, 0x1d, 0x6b, 0xd5, 0x03, 0x66, 0x4b, 0x34, 0x8e, 0xed, 0x9e,
	0x7a, 0x7f, 0xf1, 0x42, 0x6c, 0x7e, 0xaa, 0x07, 0xfc, 0x81, 0x66, 0x64, 0xf3, 0xfe, 0xdf, 0x08,
	0xda, 0x5c, 0x4c, 0x69, 0x1f, 0xe5, 0xb7, 0x47, 0xa2, 0xbc, 0xae, 0x83, 0x39, 0x64, 0xeb, 0x5c,
	0xed, 0xd7, 0x92, 0xcc, 0xf6, 0xbf, 0x74, 0x07, 0xd0, 0x29, 0x71, 0x43, 0x7b, 0x64, 0xf8, 0xc5,
	0xdf, 0xbd, 0xc9, 0x8f, 0x65, 0xf6, 0x02, 0x5a, 0x6e, 0xf4, 0xb0, 0x9c, 0xe0, 0xc6, 0xdf, 0xbc,
	0xc9, 0x6b, 0xe1, 0x39, 0x87, 0xce, 0xd1, 0x67, 0xd8, 0x25, 0x34, 0xe6, 0xb8, 0xad, 0xfe, 0xb7,
	0x3b, 0xb2, 0xb7, 0x70, 0xba, 0x16, 0xc5, 0x0a, 0xab, 0x45, 0x7a, 0x74, 0xf5, 0x02, 0xf1, 0xe9,
	0xe4, 0x63, 0xf4, 0x39, 0x87, 0x44, 0x9b, 0x3c, 0x99, 0x6d, 0x97, 0x68, 0x0a, 0x9c, 0xe4, 0x68,
	0x92, 0xa9, 0x18, 0x1b, 0x25, 0xc3, 0x92, 0xdb, 0xa4, 0x7a, 0x0a, 0xfb, 0x31, 0xbf, 0x3e, 0xe4,
	0x8a, 0x66, 0xab, 0x71, 0x22, 0xf5, 0x22, 0x3d, 0x68, 0x4b, 0x43, 0x5b, 0x1a, 0xda, 0xd2, 0xe3,
	0x17, 0x34, 0x7e, 0xe2, 0x0b, 0xb7, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x65, 0x55, 0xbd, 0x7e,
	0x5c, 0x03, 0x00, 0x00,
}
