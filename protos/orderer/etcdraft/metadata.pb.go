


package etcdraft

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 




type BlockMetadata struct {
	
	
	ConsenterIds []uint64 `protobuf:"varint,1,rep,packed,name=consenter_ids,json=consenterIds,proto3" json:"consenter_ids,omitempty"`
	
	
	NextConsenterId uint64 `protobuf:"varint,2,opt,name=next_consenter_id,json=nextConsenterId,proto3" json:"next_consenter_id,omitempty"`
	
	RaftIndex            uint64   `protobuf:"varint,3,opt,name=raft_index,json=raftIndex,proto3" json:"raft_index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockMetadata) Reset()         { *m = BlockMetadata{} }
func (m *BlockMetadata) String() string { return proto.CompactTextString(m) }
func (*BlockMetadata) ProtoMessage()    {}
func (*BlockMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d0323e5051228ea, []int{0}
}

func (m *BlockMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockMetadata.Unmarshal(m, b)
}
func (m *BlockMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockMetadata.Marshal(b, m, deterministic)
}
func (m *BlockMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockMetadata.Merge(m, src)
}
func (m *BlockMetadata) XXX_Size() int {
	return xxx_messageInfo_BlockMetadata.Size(m)
}
func (m *BlockMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_BlockMetadata proto.InternalMessageInfo

func (m *BlockMetadata) GetConsenterIds() []uint64 {
	if m != nil {
		return m.ConsenterIds
	}
	return nil
}

func (m *BlockMetadata) GetNextConsenterId() uint64 {
	if m != nil {
		return m.NextConsenterId
	}
	return 0
}

func (m *BlockMetadata) GetRaftIndex() uint64 {
	if m != nil {
		return m.RaftIndex
	}
	return 0
}


type ClusterMetadata struct {
	
	ActiveNodes          []uint64 `protobuf:"varint,1,rep,packed,name=active_nodes,json=activeNodes,proto3" json:"active_nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClusterMetadata) Reset()         { *m = ClusterMetadata{} }
func (m *ClusterMetadata) String() string { return proto.CompactTextString(m) }
func (*ClusterMetadata) ProtoMessage()    {}
func (*ClusterMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d0323e5051228ea, []int{1}
}

func (m *ClusterMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClusterMetadata.Unmarshal(m, b)
}
func (m *ClusterMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClusterMetadata.Marshal(b, m, deterministic)
}
func (m *ClusterMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterMetadata.Merge(m, src)
}
func (m *ClusterMetadata) XXX_Size() int {
	return xxx_messageInfo_ClusterMetadata.Size(m)
}
func (m *ClusterMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterMetadata proto.InternalMessageInfo

func (m *ClusterMetadata) GetActiveNodes() []uint64 {
	if m != nil {
		return m.ActiveNodes
	}
	return nil
}

func init() {
	proto.RegisterType((*BlockMetadata)(nil), "etcdraft.BlockMetadata")
	proto.RegisterType((*ClusterMetadata)(nil), "etcdraft.ClusterMetadata")
}

func init() { proto.RegisterFile("orderer/etcdraft/metadata.proto", fileDescriptor_6d0323e5051228ea) }

var fileDescriptor_6d0323e5051228ea = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x3d, 0x4f, 0x02, 0x41,
	0x10, 0x86, 0x83, 0x10, 0xa3, 0x23, 0x84, 0x78, 0xd5, 0x35, 0x46, 0xc4, 0x86, 0x58, 0xec, 0x16,
	0xea, 0x1f, 0x80, 0x8a, 0x42, 0x0b, 0x4a, 0x9b, 0xcb, 0xde, 0xee, 0x70, 0x6c, 0x3c, 0x76, 0xc8,
	0xec, 0x60, 0xb0, 0xf2, 0xaf, 0x9b, 0xe5, 0x38, 0xbc, 0xd8, 0x3e, 0xef, 0xf3, 0x66, 0x3e, 0xe0,
	0x9e, 0xd8, 0x21, 0x23, 0x6b, 0x14, 0xeb, 0xd8, 0xac, 0x45, 0x6f, 0x51, 0x8c, 0x33, 0x62, 0xd4,
	0x8e, 0x49, 0x28, 0xbb, 0x6a, 0x83, 0xe9, 0x0f, 0x8c, 0xe6, 0x35, 0xd9, 0xcf, 0xb7, 0x93, 0x90,
	0x3d, 0xc2, 0xc8, 0x52, 0x88, 0x18, 0x04, 0xb9, 0xf0, 0x2e, 0xe6, 0xbd, 0x49, 0x7f, 0x36, 0x58,
	0x0d, 0xcf, 0x70, 0xe9, 0x62, 0xf6, 0x04, 0xb7, 0x01, 0x0f, 0x52, 0x74, 0xcd, 0xfc, 0x62, 0xd2,
	0x9b, 0x0d, 0x56, 0xe3, 0x14, 0x2c, 0xfe, 0xe4, 0xec, 0x0e, 0x20, 0x4d, 0x2a, 0x7c, 0x70, 0x78,
	0xc8, 0xfb, 0x47, 0xe9, 0x3a, 0x91, 0x65, 0x02, 0xd3, 0x17, 0x18, 0x2f, 0xea, 0x7d, 0x14, 0xe4,
	0xf3, 0x0a, 0x0f, 0x30, 0x34, 0x56, 0xfc, 0x17, 0x16, 0x81, 0x1c, 0xb6, 0x1b, 0xdc, 0x34, 0xec,
	0x3d, 0xa1, 0x79, 0x05, 0x8a, 0xb8, 0x52, 0x9b, 0xef, 0x1d, 0x72, 0x8d, 0xae, 0x42, 0x56, 0x6b,
	0x53, 0xb2, 0xb7, 0xcd, 0x81, 0x51, 0x9d, 0x3e, 0xa0, 0xda, 0x43, 0x3f, 0x5e, 0x2b, 0x2f, 0x9b,
	0x7d, 0xa9, 0x2c, 0x6d, 0x75, 0xa7, 0xa6, 0x9b, 0x9a, 0x6e, 0x6a, 0xfa, 0xff, 0xe3, 0xca, 0xcb,
	0x63, 0xf0, 0xfc, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xc6, 0xb8, 0xaf, 0x60, 0x53, 0x01, 0x00, 0x00,
}
