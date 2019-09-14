


package pvtdatastorage

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 

type ExpiryData struct {
	Map                  map[string]*Collections `protobuf:"bytes,1,rep,name=map,proto3" json:"map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ExpiryData) Reset()         { *m = ExpiryData{} }
func (m *ExpiryData) String() string { return proto.CompactTextString(m) }
func (*ExpiryData) ProtoMessage()    {}
func (*ExpiryData) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f0cbd2d16bac879, []int{0}
}

func (m *ExpiryData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExpiryData.Unmarshal(m, b)
}
func (m *ExpiryData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExpiryData.Marshal(b, m, deterministic)
}
func (m *ExpiryData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExpiryData.Merge(m, src)
}
func (m *ExpiryData) XXX_Size() int {
	return xxx_messageInfo_ExpiryData.Size(m)
}
func (m *ExpiryData) XXX_DiscardUnknown() {
	xxx_messageInfo_ExpiryData.DiscardUnknown(m)
}

var xxx_messageInfo_ExpiryData proto.InternalMessageInfo

func (m *ExpiryData) GetMap() map[string]*Collections {
	if m != nil {
		return m.Map
	}
	return nil
}

type Collections struct {
	
	
	Map map[string]*TxNums `protobuf:"bytes,1,rep,name=map,proto3" json:"map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	
	
	MissingDataMap       map[string]bool `protobuf:"bytes,2,rep,name=missingDataMap,proto3" json:"missingDataMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Collections) Reset()         { *m = Collections{} }
func (m *Collections) String() string { return proto.CompactTextString(m) }
func (*Collections) ProtoMessage()    {}
func (*Collections) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f0cbd2d16bac879, []int{1}
}

func (m *Collections) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Collections.Unmarshal(m, b)
}
func (m *Collections) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Collections.Marshal(b, m, deterministic)
}
func (m *Collections) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Collections.Merge(m, src)
}
func (m *Collections) XXX_Size() int {
	return xxx_messageInfo_Collections.Size(m)
}
func (m *Collections) XXX_DiscardUnknown() {
	xxx_messageInfo_Collections.DiscardUnknown(m)
}

var xxx_messageInfo_Collections proto.InternalMessageInfo

func (m *Collections) GetMap() map[string]*TxNums {
	if m != nil {
		return m.Map
	}
	return nil
}

func (m *Collections) GetMissingDataMap() map[string]bool {
	if m != nil {
		return m.MissingDataMap
	}
	return nil
}

type TxNums struct {
	List                 []uint64 `protobuf:"varint,1,rep,packed,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxNums) Reset()         { *m = TxNums{} }
func (m *TxNums) String() string { return proto.CompactTextString(m) }
func (*TxNums) ProtoMessage()    {}
func (*TxNums) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f0cbd2d16bac879, []int{2}
}

func (m *TxNums) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxNums.Unmarshal(m, b)
}
func (m *TxNums) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxNums.Marshal(b, m, deterministic)
}
func (m *TxNums) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxNums.Merge(m, src)
}
func (m *TxNums) XXX_Size() int {
	return xxx_messageInfo_TxNums.Size(m)
}
func (m *TxNums) XXX_DiscardUnknown() {
	xxx_messageInfo_TxNums.DiscardUnknown(m)
}

var xxx_messageInfo_TxNums proto.InternalMessageInfo

func (m *TxNums) GetList() []uint64 {
	if m != nil {
		return m.List
	}
	return nil
}

type CollElgInfo struct {
	NsCollMap            map[string]*CollNames `protobuf:"bytes,1,rep,name=nsCollMap,proto3" json:"nsCollMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *CollElgInfo) Reset()         { *m = CollElgInfo{} }
func (m *CollElgInfo) String() string { return proto.CompactTextString(m) }
func (*CollElgInfo) ProtoMessage()    {}
func (*CollElgInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f0cbd2d16bac879, []int{3}
}

func (m *CollElgInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollElgInfo.Unmarshal(m, b)
}
func (m *CollElgInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollElgInfo.Marshal(b, m, deterministic)
}
func (m *CollElgInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollElgInfo.Merge(m, src)
}
func (m *CollElgInfo) XXX_Size() int {
	return xxx_messageInfo_CollElgInfo.Size(m)
}
func (m *CollElgInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_CollElgInfo.DiscardUnknown(m)
}

var xxx_messageInfo_CollElgInfo proto.InternalMessageInfo

func (m *CollElgInfo) GetNsCollMap() map[string]*CollNames {
	if m != nil {
		return m.NsCollMap
	}
	return nil
}

type CollNames struct {
	Entries              []string `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CollNames) Reset()         { *m = CollNames{} }
func (m *CollNames) String() string { return proto.CompactTextString(m) }
func (*CollNames) ProtoMessage()    {}
func (*CollNames) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f0cbd2d16bac879, []int{4}
}

func (m *CollNames) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollNames.Unmarshal(m, b)
}
func (m *CollNames) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollNames.Marshal(b, m, deterministic)
}
func (m *CollNames) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollNames.Merge(m, src)
}
func (m *CollNames) XXX_Size() int {
	return xxx_messageInfo_CollNames.Size(m)
}
func (m *CollNames) XXX_DiscardUnknown() {
	xxx_messageInfo_CollNames.DiscardUnknown(m)
}

var xxx_messageInfo_CollNames proto.InternalMessageInfo

func (m *CollNames) GetEntries() []string {
	if m != nil {
		return m.Entries
	}
	return nil
}

func init() {
	proto.RegisterType((*ExpiryData)(nil), "pvtdatastorage.ExpiryData")
	proto.RegisterMapType((map[string]*Collections)(nil), "pvtdatastorage.ExpiryData.MapEntry")
	proto.RegisterType((*Collections)(nil), "pvtdatastorage.Collections")
	proto.RegisterMapType((map[string]*TxNums)(nil), "pvtdatastorage.Collections.MapEntry")
	proto.RegisterMapType((map[string]bool)(nil), "pvtdatastorage.Collections.MissingDataMapEntry")
	proto.RegisterType((*TxNums)(nil), "pvtdatastorage.TxNums")
	proto.RegisterType((*CollElgInfo)(nil), "pvtdatastorage.CollElgInfo")
	proto.RegisterMapType((map[string]*CollNames)(nil), "pvtdatastorage.CollElgInfo.NsCollMapEntry")
	proto.RegisterType((*CollNames)(nil), "pvtdatastorage.CollNames")
}

func init() { proto.RegisterFile("persistent_msgs.proto", fileDescriptor_0f0cbd2d16bac879) }

var fileDescriptor_0f0cbd2d16bac879 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x93, 0xdf, 0x6a, 0xdb, 0x30,
	0x14, 0xc6, 0xb1, 0x93, 0x65, 0xf1, 0x09, 0x84, 0xa1, 0xfd, 0xc1, 0xf3, 0x76, 0x11, 0xb2, 0x0d,
	0xc2, 0x18, 0x36, 0xcb, 0xd8, 0x08, 0xb9, 0x5b, 0xdb, 0x40, 0x7b, 0x11, 0x5f, 0xb8, 0x85, 0x40,
	0x6f, 0x8a, 0xe2, 0x28, 0x8e, 0xa8, 0x6d, 0x09, 0x49, 0x09, 0xf1, 0x9b, 0xf4, 0x31, 0xda, 0x37,
	0x2c, 0xb6, 0xf3, 0xc7, 0x0a, 0x26, 0x77, 0xf2, 0xd1, 0x77, 0x7e, 0xe7, 0x3b, 0x1f, 0x16, 0x7c,
	0xe4, 0x44, 0x48, 0x2a, 0x15, 0x49, 0xd5, 0x43, 0x22, 0x23, 0xe9, 0x72, 0xc1, 0x14, 0x43, 0x5d,
	0xbe, 0x51, 0x0b, 0xac, 0xb0, 0x54, 0x4c, 0xe0, 0x88, 0xf4, 0x9f, 0x0c, 0x80, 0xc9, 0x96, 0x53,
	0x91, 0x5d, 0x61, 0x85, 0xd1, 0x5f, 0x68, 0x24, 0x98, 0xdb, 0x46, 0xaf, 0x31, 0xe8, 0x0c, 0xbf,
	0xb9, 0xba, 0xd8, 0x3d, 0x0a, 0xdd, 0x29, 0xe6, 0x93, 0x54, 0x89, 0x2c, 0xc8, 0xf5, 0xce, 0x2d,
	0xb4, 0xf7, 0x05, 0xf4, 0x0e, 0x1a, 0x8f, 0x24, 0xb3, 0x8d, 0x9e, 0x31, 0xb0, 0x82, 0xfc, 0x88,
	0x7e, 0xc3, 0x9b, 0x0d, 0x8e, 0xd7, 0xc4, 0x36, 0x7b, 0xc6, 0xa0, 0x33, 0xfc, 0x72, 0x8a, 0xbd,
	0x64, 0x71, 0x4c, 0x42, 0x45, 0x59, 0x2a, 0x83, 0x52, 0x39, 0x36, 0x47, 0x46, 0xff, 0xc5, 0x84,
	0x4e, 0xe5, 0x0a, 0xfd, 0xab, 0x7a, 0xfb, 0x7e, 0x06, 0xa2, 0x9b, 0x43, 0x33, 0xe8, 0x26, 0x54,
	0x4a, 0x9a, 0x46, 0xb9, 0xf3, 0x29, 0xe6, 0xb6, 0x59, 0x20, 0xbc, 0xb3, 0x08, 0xad, 0xa3, 0xa4,
	0x9d, 0x60, 0x1c, 0xff, 0xec, 0xd6, 0xbf, 0xf4, 0xad, 0x3f, 0x9d, 0x4e, 0xbb, 0xdb, 0xfa, 0xeb,
	0xa4, 0xba, 0xb0, 0xf3, 0x1f, 0xde, 0xd7, 0x8c, 0xad, 0x41, 0x7f, 0xa8, 0xa2, 0xdb, 0xd5, 0xcc,
	0xbe, 0x42, 0xab, 0xe4, 0x22, 0x04, 0xcd, 0x98, 0x4a, 0x55, 0xc4, 0xd5, 0x0c, 0x8a, 0x73, 0xff,
	0xd9, 0x28, 0x13, 0x9d, 0xc4, 0xd1, 0x4d, 0xba, 0x64, 0xe8, 0x1a, 0xac, 0x54, 0xe6, 0x85, 0xe9,
	0x21, 0xd7, 0x9f, 0x75, 0xa1, 0xec, 0xf4, 0xae, 0xbf, 0x17, 0x97, 0x79, 0x1c, 0x9b, 0x9d, 0x19,
	0x74, 0xf5, 0xcb, 0x1a, 0xd7, 0x9e, 0x1e, 0xc8, 0xe7, 0xba, 0x49, 0x3e, 0x4e, 0x88, 0xf6, 0x13,
	0xfc, 0x00, 0xeb, 0x50, 0x47, 0x36, 0xbc, 0x25, 0xa9, 0x12, 0x94, 0xc8, 0xc2, 0xad, 0x15, 0xec,
	0x3f, 0x2f, 0xc6, 0xf7, 0xa3, 0x88, 0xaa, 0xd5, 0x7a, 0xee, 0x86, 0x2c, 0xf1, 0x56, 0x19, 0x27,
	0x22, 0x26, 0x8b, 0x88, 0x08, 0x6f, 0x89, 0xe7, 0x82, 0x86, 0x5e, 0xc8, 0x04, 0xf1, 0x76, 0x25,
	0x7d, 0xee, 0xbc, 0x55, 0xbc, 0x8c, 0x3f, 0xaf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb7, 0x3e, 0x75,
	0xee, 0x32, 0x03, 0x00, 0x00,
}
