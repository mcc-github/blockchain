


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
	return fileDescriptor_configuration_65c7ed1a646136e8, []int{0}
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
	return fileDescriptor_configuration_65c7ed1a646136e8, []int{1}
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
	return fileDescriptor_configuration_65c7ed1a646136e8, []int{2}
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

func init() {
	proto.RegisterType((*Metadata)(nil), "etcdraft.Metadata")
	proto.RegisterType((*Consenter)(nil), "etcdraft.Consenter")
	proto.RegisterType((*Options)(nil), "etcdraft.Options")
}

func init() {
	proto.RegisterFile("orderer/etcdraft/configuration.proto", fileDescriptor_configuration_65c7ed1a646136e8)
}

var fileDescriptor_configuration_65c7ed1a646136e8 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0xcd, 0x6e, 0xdb, 0x30,
	0x10, 0x84, 0xa1, 0xda, 0x6d, 0x12, 0xc6, 0x6a, 0x1a, 0xf6, 0xa2, 0xa3, 0xe1, 0xfe, 0x19, 0x2d,
	0x40, 0x01, 0x09, 0xfa, 0x02, 0xcd, 0x29, 0x07, 0xa3, 0x85, 0x9a, 0x53, 0x2f, 0x02, 0x4d, 0xaf,
	0x28, 0x22, 0x94, 0x28, 0x2c, 0x37, 0x81, 0x9b, 0x6b, 0x1f, 0xb0, 0xaf, 0x54, 0x90, 0x94, 0x1c,
	0xa3, 0xb7, 0xc5, 0xcc, 0x37, 0x8b, 0x01, 0x76, 0xd9, 0x7b, 0x87, 0x3b, 0x40, 0xc0, 0x12, 0x48,
	0xed, 0x50, 0x36, 0x54, 0x2a, 0xd7, 0x37, 0x46, 0x3f, 0xa0, 0x24, 0xe3, 0x7a, 0x31, 0xa0, 0x23,
	0xc7, 0x4f, 0x27, 0x77, 0x65, 0xd9, 0xe9, 0x06, 0x48, 0xee, 0x24, 0x49, 0x7e, 0xcd, 0x98, 0x72,
	0xbd, 0x87, 0x9e, 0x00, 0x7d, 0x91, 0x2d, 0x67, 0xeb, 0xf3, 0xab, 0xb7, 0x62, 0x42, 0xc5, 0xcd,
	0xe4, 0x55, 0x47, 0x18, 0xff, 0xc2, 0x4e, 0xdc, 0x10, 0x56, 0xfb, 0xe2, 0xc5, 0x32, 0x5b, 0x9f,
	0x5f, 0x5d, 0x3e, 0x27, 0xbe, 0x27, 0xa3, 0x9a, 0x88, 0xd5, 0x9f, 0x8c, 0x9d, 0x1d, 0xd6, 0x70,
	0xce, 0xe6, 0xad, 0xf3, 0x54, 0x64, 0xcb, 0x6c, 0x7d, 0x56, 0xc5, 0x39, 0x68, 0x83, 0x43, 0x8a,
	0xbb, 0xf2, 0x2a, 0xce, 0xfc, 0x23, 0xbb, 0x50, 0xd6, 0x40, 0x4f, 0x35, 0x59, 0x5f, 0x2b, 0x40,
	0x2a, 0x66, 0xcb, 0x6c, 0xbd, 0xa8, 0xf2, 0x24, 0xdf, 0x59, 0x7f, 0x03, 0x89, 0xf3, 0x80, 0x8f,
	0x80, 0xcf, 0xdc, 0x3c, 0x71, 0x49, 0x1e, 0xb9, 0xd5, 0xdf, 0x8c, 0x9d, 0x8c, 0xd5, 0xf8, 0x3b,
	0x96, 0x93, 0x51, 0xf7, 0xb5, 0x09, 0x8d, 0x1e, 0xa5, 0x8d, 0x65, 0xe6, 0xd5, 0x22, 0x88, 0xb7,
	0xa3, 0x16, 0x20, 0xb0, 0xa0, 0x42, 0xa2, 0x0e, 0xc6, 0xd8, 0x6e, 0x31, 0x89, 0x77, 0x46, 0xdd,
	0xf3, 0x0f, 0xec, 0x75, 0x0b, 0x12, 0x69, 0x0b, 0x92, 0x12, 0x35, 0x8b, 0x54, 0x7e, 0x50, 0x23,
	0xf6, 0x99, 0x5d, 0x76, 0x72, 0x5f, 0x9b, 0xbe, 0xb1, 0x46, 0xb7, 0x54, 0x77, 0x5e, 0xfb, 0x58,
	0x33, 0xaf, 0x2e, 0x3a, 0xb9, 0xbf, 0x1d, 0xf5, 0x8d, 0xd7, 0x9e, 0x7f, 0x62, 0x6f, 0x02, 0xeb,
	0xcd, 0x13, 0xd4, 0x03, 0x60, 0x60, 0x8b, 0x97, 0xb1, 0x5f, 0xde, 0xc9, 0xfd, 0x4f, 0xf3, 0x04,
	0x3f, 0x00, 0x37, 0x5e, 0x7f, 0xd3, 0x4c, 0x38, 0xd4, 0xa2, 0xfd, 0x3d, 0x00, 0x5a, 0xd8, 0x69,
	0x40, 0xd1, 0xc8, 0x2d, 0x1a, 0x95, 0xee, 0xed, 0xc5, 0xf8, 0x15, 0x87, 0xd3, 0xfc, 0xfa, 0xaa,
	0x0d, 0xb5, 0x0f, 0x5b, 0xa1, 0x5c, 0x57, 0x1e, 0xc5, 0xca, 0x14, 0x2b, 0x53, 0xac, 0xfc, 0xff,
	0x99, 0xb6, 0xaf, 0xa2, 0x71, 0xfd, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xfb, 0xd4, 0x3d, 0xca, 0x67,
	0x02, 0x00, 0x00,
}
