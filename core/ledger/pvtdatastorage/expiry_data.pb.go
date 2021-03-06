


package pvtdatastorage 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type ExpiryData struct {
	Map                  map[string]*Collections `protobuf:"bytes,1,rep,name=map" json:"map,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ExpiryData) Reset()         { *m = ExpiryData{} }
func (m *ExpiryData) String() string { return proto.CompactTextString(m) }
func (*ExpiryData) ProtoMessage()    {}
func (*ExpiryData) Descriptor() ([]byte, []int) {
	return fileDescriptor_expiry_data_bdde9920f0b783de, []int{0}
}
func (m *ExpiryData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExpiryData.Unmarshal(m, b)
}
func (m *ExpiryData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExpiryData.Marshal(b, m, deterministic)
}
func (dst *ExpiryData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExpiryData.Merge(dst, src)
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
	
	
	Map map[string]*TxNums `protobuf:"bytes,1,rep,name=map" json:"map,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	
	
	MissingDataMap       map[string]bool `protobuf:"bytes,2,rep,name=missingDataMap" json:"missingDataMap,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Collections) Reset()         { *m = Collections{} }
func (m *Collections) String() string { return proto.CompactTextString(m) }
func (*Collections) ProtoMessage()    {}
func (*Collections) Descriptor() ([]byte, []int) {
	return fileDescriptor_expiry_data_bdde9920f0b783de, []int{1}
}
func (m *Collections) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Collections.Unmarshal(m, b)
}
func (m *Collections) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Collections.Marshal(b, m, deterministic)
}
func (dst *Collections) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Collections.Merge(dst, src)
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
	List                 []uint64 `protobuf:"varint,1,rep,packed,name=list" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxNums) Reset()         { *m = TxNums{} }
func (m *TxNums) String() string { return proto.CompactTextString(m) }
func (*TxNums) ProtoMessage()    {}
func (*TxNums) Descriptor() ([]byte, []int) {
	return fileDescriptor_expiry_data_bdde9920f0b783de, []int{2}
}
func (m *TxNums) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxNums.Unmarshal(m, b)
}
func (m *TxNums) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxNums.Marshal(b, m, deterministic)
}
func (dst *TxNums) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxNums.Merge(dst, src)
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

func init() {
	proto.RegisterType((*ExpiryData)(nil), "pvtdatastorage.ExpiryData")
	proto.RegisterMapType((map[string]*Collections)(nil), "pvtdatastorage.ExpiryData.MapEntry")
	proto.RegisterType((*Collections)(nil), "pvtdatastorage.Collections")
	proto.RegisterMapType((map[string]*TxNums)(nil), "pvtdatastorage.Collections.MapEntry")
	proto.RegisterMapType((map[string]bool)(nil), "pvtdatastorage.Collections.MissingDataMapEntry")
	proto.RegisterType((*TxNums)(nil), "pvtdatastorage.TxNums")
}

func init() { proto.RegisterFile("expiry_data.proto", fileDescriptor_expiry_data_bdde9920f0b783de) }

var fileDescriptor_expiry_data_bdde9920f0b783de = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcf, 0x4a, 0xf3, 0x40,
	0x14, 0xc5, 0x99, 0xa4, 0x5f, 0xe9, 0x77, 0x03, 0x45, 0x47, 0x91, 0x10, 0x5d, 0x84, 0xea, 0x22,
	0x0b, 0x49, 0xb0, 0xa2, 0x94, 0xee, 0xfc, 0xd3, 0x65, 0xbb, 0x88, 0x82, 0xe0, 0x46, 0x26, 0xe9,
	0x98, 0x0e, 0x26, 0x99, 0x61, 0x32, 0x29, 0xcd, 0x9b, 0xf8, 0x1a, 0xbe, 0xa1, 0x24, 0x51, 0x4c,
	0x42, 0xc9, 0xee, 0xce, 0x9d, 0x33, 0xbf, 0x7b, 0xce, 0x65, 0xe0, 0x90, 0xee, 0x04, 0x93, 0xc5,
	0xdb, 0x9a, 0x28, 0xe2, 0x0a, 0xc9, 0x15, 0xc7, 0x63, 0xb1, 0x55, 0xe5, 0x31, 0x53, 0x5c, 0x92,
	0x88, 0x4e, 0x3e, 0x11, 0xc0, 0xa2, 0x52, 0x3d, 0x12, 0x45, 0xf0, 0x0d, 0xe8, 0x09, 0x11, 0x26,
	0xb2, 0x75, 0xc7, 0x98, 0x9e, 0xbb, 0x6d, 0xb1, 0xfb, 0x27, 0x74, 0x97, 0x44, 0x2c, 0x52, 0x25,
	0x0b, 0xbf, 0xd4, 0x5b, 0x4f, 0x30, 0xfa, 0x6d, 0xe0, 0x03, 0xd0, 0x3f, 0x68, 0x61, 0x22, 0x1b,
	0x39, 0xff, 0xfd, 0xb2, 0xc4, 0x57, 0xf0, 0x6f, 0x4b, 0xe2, 0x9c, 0x9a, 0x9a, 0x8d, 0x1c, 0x63,
	0x7a, 0xda, 0xc5, 0x3e, 0xf0, 0x38, 0xa6, 0xa1, 0x62, 0x3c, 0xcd, 0xfc, 0x5a, 0x39, 0xd7, 0x66,
	0x68, 0xf2, 0xa5, 0x81, 0xd1, 0xb8, 0xc2, 0xb7, 0x4d, 0x6f, 0x17, 0x3d, 0x90, 0xb6, 0x39, 0xfc,
	0x02, 0xe3, 0x84, 0x65, 0x19, 0x4b, 0xa3, 0xd2, 0xf9, 0x92, 0x08, 0x53, 0xab, 0x10, 0x5e, 0x2f,
	0xa2, 0xf5, 0xa2, 0xa6, 0x75, 0x30, 0xd6, 0xaa, 0x37, 0xf5, 0x65, 0x3b, 0xf5, 0x49, 0x77, 0xda,
	0xf3, 0x6e, 0x95, 0x27, 0xcd, 0xc0, 0xd6, 0x1d, 0x1c, 0xed, 0x19, 0xbb, 0x07, 0x7d, 0xdc, 0x44,
	0x8f, 0x9a, 0x3b, 0x3b, 0x83, 0x61, 0xcd, 0xc5, 0x18, 0x06, 0x31, 0xcb, 0x54, 0xb5, 0xae, 0x81,
	0x5f, 0xd5, 0xf7, 0xf3, 0xd7, 0x59, 0xc4, 0xd4, 0x26, 0x0f, 0xdc, 0x90, 0x27, 0xde, 0xa6, 0x10,
	0x54, 0xc6, 0x74, 0x1d, 0x51, 0xe9, 0xbd, 0x93, 0x40, 0xb2, 0xd0, 0x0b, 0xb9, 0xa4, 0xde, 0x4f,
	0xab, 0x6d, 0x37, 0x18, 0x56, 0xff, 0xe7, 0xfa, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x3e, 0x0e, 0xe9,
	0x91, 0x54, 0x02, 0x00, 0x00,
}
