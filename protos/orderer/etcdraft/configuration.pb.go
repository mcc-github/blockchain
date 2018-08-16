


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
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_d24fbe9bcd75d252, []int{0}
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
	return fileDescriptor_configuration_d24fbe9bcd75d252, []int{1}
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

func init() {
	proto.RegisterType((*Metadata)(nil), "etcdraft.Metadata")
	proto.RegisterType((*Consenter)(nil), "etcdraft.Consenter")
}

func init() {
	proto.RegisterFile("orderer/etcdraft/configuration.proto", fileDescriptor_configuration_d24fbe9bcd75d252)
}

var fileDescriptor_configuration_d24fbe9bcd75d252 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xb1, 0x6b, 0xf3, 0x30,
	0x10, 0xc5, 0xd1, 0x97, 0xf0, 0x91, 0xa8, 0x0d, 0x05, 0x75, 0xf1, 0x68, 0x42, 0x29, 0x9e, 0x24,
	0x68, 0xe8, 0x5c, 0x68, 0xe6, 0x2e, 0xa6, 0x53, 0x97, 0x20, 0xcb, 0x67, 0x59, 0xe0, 0xea, 0xcc,
	0xe9, 0x52, 0xe8, 0xdc, 0x7f, 0xbc, 0xd4, 0x8a, 0xd3, 0xd0, 0xed, 0x78, 0xbf, 0xdf, 0x1b, 0xee,
	0xc9, 0x3b, 0xa4, 0x16, 0x08, 0xc8, 0x00, 0xbb, 0x96, 0x6c, 0xc7, 0xc6, 0x61, 0xec, 0x82, 0x3f,
	0x92, 0xe5, 0x80, 0x51, 0x8f, 0x84, 0x8c, 0x6a, 0x35, 0xd3, 0xed, 0x93, 0x5c, 0xbd, 0x00, 0xdb,
	0xd6, 0xb2, 0x55, 0x3b, 0x29, 0x1d, 0xc6, 0x04, 0x91, 0x81, 0x52, 0x21, 0xca, 0x45, 0x75, 0xf5,
	0x70, 0xab, 0x67, 0x55, 0xef, 0x67, 0x56, 0x5f, 0x68, 0xdb, 0x2f, 0x21, 0xd7, 0x67, 0xa2, 0x94,
	0x5c, 0xf6, 0x98, 0xb8, 0x10, 0xa5, 0xa8, 0xd6, 0xf5, 0x74, 0xff, 0x64, 0x23, 0x12, 0x17, 0xff,
	0x4a, 0x51, 0x6d, 0xea, 0xe9, 0x56, 0xf7, 0xf2, 0xc6, 0x0d, 0x01, 0x22, 0x1f, 0x78, 0x48, 0x07,
	0x07, 0xc4, 0xc5, 0xa2, 0x14, 0xd5, 0x75, 0xbd, 0xc9, 0xf1, 0xeb, 0x90, 0xf6, 0x90, 0xbd, 0x04,
	0xf4, 0x01, 0xf4, 0xeb, 0x2d, 0xb3, 0x97, 0xe3, 0x93, 0xf7, 0xec, 0xa5, 0x46, 0xf2, 0xba, 0xff,
	0x1c, 0x81, 0x06, 0x68, 0x3d, 0x90, 0xee, 0x6c, 0x43, 0xc1, 0xe5, 0x87, 0x93, 0x3e, 0xcd, 0x72,
	0xfe, 0xe6, 0xed, 0xd1, 0x07, 0xee, 0x8f, 0x8d, 0x76, 0xf8, 0x6e, 0x2e, 0x6a, 0x26, 0xd7, 0x4c,
	0xae, 0x99, 0xbf, 0x6b, 0x36, 0xff, 0x27, 0xb0, 0xfb, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xc7, 0xd4,
	0xed, 0xaf, 0x68, 0x01, 0x00, 0x00,
}
