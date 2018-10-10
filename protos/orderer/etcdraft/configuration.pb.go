


package etcdraft 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 



type Metadata struct {
	Consenters           []*Consenter `protobuf:"bytes,1,rep,name=consenters" json:"consenters,omitempty"`
	Options              *Options     `protobuf:"bytes,2,opt,name=options" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_5f4f73f0206bc0bb, []int{0}
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
	Host                 string   `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
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
	return fileDescriptor_configuration_5f4f73f0206bc0bb, []int{1}
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
	TickInterval         uint32   `protobuf:"varint,1,opt,name=tick_interval,json=tickInterval" json:"tick_interval,omitempty"`
	ElectionTick         uint32   `protobuf:"varint,2,opt,name=election_tick,json=electionTick" json:"election_tick,omitempty"`
	HeartbeatTick        uint32   `protobuf:"varint,3,opt,name=heartbeat_tick,json=heartbeatTick" json:"heartbeat_tick,omitempty"`
	MaxInflightMsgs      uint32   `protobuf:"varint,4,opt,name=max_inflight_msgs,json=maxInflightMsgs" json:"max_inflight_msgs,omitempty"`
	MaxSizePerMsg        uint64   `protobuf:"varint,5,opt,name=max_size_per_msg,json=maxSizePerMsg" json:"max_size_per_msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Options) Reset()         { *m = Options{} }
func (m *Options) String() string { return proto.CompactTextString(m) }
func (*Options) ProtoMessage()    {}
func (*Options) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_5f4f73f0206bc0bb, []int{2}
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

func (m *Options) GetTickInterval() uint32 {
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

func init() {
	proto.RegisterType((*Metadata)(nil), "etcdraft.Metadata")
	proto.RegisterType((*Consenter)(nil), "etcdraft.Consenter")
	proto.RegisterType((*Options)(nil), "etcdraft.Options")
}

func init() {
	proto.RegisterFile("orderer/etcdraft/configuration.proto", fileDescriptor_configuration_5f4f73f0206bc0bb)
}

var fileDescriptor_configuration_5f4f73f0206bc0bb = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0x4f, 0x6f, 0xd4, 0x30,
	0x10, 0xc5, 0x15, 0x76, 0xa1, 0xad, 0xbb, 0xa1, 0xd4, 0x5c, 0x72, 0x8c, 0x96, 0x7f, 0x11, 0x48,
	0x89, 0xd4, 0x8a, 0x2f, 0x40, 0x4f, 0x3d, 0xac, 0x40, 0xa1, 0x27, 0x2e, 0x91, 0xd7, 0x3b, 0x71,
	0xac, 0x3a, 0x71, 0x34, 0x9e, 0x56, 0x4b, 0xaf, 0x7c, 0x40, 0xbe, 0x12, 0xb2, 0x9d, 0x6c, 0x57,
	0xdc, 0x46, 0xef, 0xfd, 0xde, 0xe8, 0x49, 0x33, 0xec, 0xbd, 0xc5, 0x1d, 0x20, 0x60, 0x05, 0x24,
	0x77, 0x28, 0x5a, 0xaa, 0xa4, 0x1d, 0x5a, 0xad, 0x1e, 0x50, 0x90, 0xb6, 0x43, 0x39, 0xa2, 0x25,
	0xcb, 0x4f, 0x67, 0x77, 0x6d, 0xd8, 0xe9, 0x06, 0x48, 0xec, 0x04, 0x09, 0x7e, 0xcd, 0x98, 0xb4,
	0x83, 0x83, 0x81, 0x00, 0x5d, 0x96, 0xe4, 0x8b, 0xe2, 0xfc, 0xea, 0x6d, 0x39, 0xa3, 0xe5, 0xcd,
	0xec, 0xd5, 0x47, 0x18, 0xff, 0xc2, 0x4e, 0xec, 0xe8, 0x57, 0xbb, 0xec, 0x45, 0x9e, 0x14, 0xe7,
	0x57, 0x97, 0xcf, 0x89, 0xef, 0xd1, 0xa8, 0x67, 0x62, 0xfd, 0x27, 0x61, 0x67, 0x87, 0x35, 0x9c,
	0xb3, 0x65, 0x67, 0x1d, 0x65, 0x49, 0x9e, 0x14, 0x67, 0x75, 0x98, 0xbd, 0x36, 0x5a, 0xa4, 0xb0,
	0x2b, 0xad, 0xc3, 0xcc, 0x3f, 0xb2, 0x0b, 0x69, 0x34, 0x0c, 0xd4, 0x90, 0x71, 0x8d, 0x04, 0xa4,
	0x6c, 0x91, 0x27, 0xc5, 0xaa, 0x4e, 0xa3, 0x7c, 0x67, 0xdc, 0x0d, 0x44, 0xce, 0x01, 0x3e, 0x02,
	0x3e, 0x73, 0xcb, 0xc8, 0x45, 0x79, 0xe2, 0xd6, 0x7f, 0x13, 0x76, 0x32, 0x55, 0xe3, 0xef, 0x58,
	0x4a, 0x5a, 0xde, 0x37, 0xda, 0x37, 0x7a, 0x14, 0x26, 0x94, 0x49, 0xeb, 0x95, 0x17, 0x6f, 0x27,
	0xcd, 0x43, 0x60, 0x40, 0xfa, 0x44, 0xe3, 0x8d, 0xa9, 0xdd, 0x6a, 0x16, 0xef, 0xb4, 0xbc, 0xe7,
	0x1f, 0xd8, 0xeb, 0x0e, 0x04, 0xd2, 0x16, 0x04, 0x45, 0x6a, 0x11, 0xa8, 0xf4, 0xa0, 0x06, 0xec,
	0x33, 0xbb, 0xec, 0xc5, 0xbe, 0xd1, 0x43, 0x6b, 0xb4, 0xea, 0xa8, 0xe9, 0x9d, 0x72, 0xa1, 0x66,
	0x5a, 0x5f, 0xf4, 0x62, 0x7f, 0x3b, 0xe9, 0x1b, 0xa7, 0x1c, 0xff, 0xc4, 0xde, 0x78, 0xd6, 0xe9,
	0x27, 0x68, 0x46, 0x40, 0xcf, 0x66, 0x2f, 0xf3, 0xa4, 0x58, 0xd6, 0x69, 0x2f, 0xf6, 0x3f, 0xf5,
	0x13, 0xfc, 0x00, 0xdc, 0x38, 0xf5, 0x4d, 0xb1, 0xd2, 0xa2, 0x2a, 0xbb, 0xdf, 0x23, 0xa0, 0x81,
	0x9d, 0x02, 0x2c, 0x5b, 0xb1, 0x45, 0x2d, 0xe3, 0xbd, 0x5d, 0x39, 0x7d, 0xc5, 0xe1, 0x34, 0xbf,
	0xbe, 0x2a, 0x4d, 0xdd, 0xc3, 0xb6, 0x94, 0xb6, 0xaf, 0x8e, 0x62, 0x55, 0x8c, 0x55, 0x31, 0x56,
	0xfd, 0xff, 0x4c, 0xdb, 0x57, 0xc1, 0xb8, 0xfe, 0x17, 0x00, 0x00, 0xff, 0xff, 0xc6, 0xd9, 0xc9,
	0x01, 0x67, 0x02, 0x00, 0x00,
}
