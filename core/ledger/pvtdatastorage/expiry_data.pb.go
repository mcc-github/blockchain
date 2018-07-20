


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
	return fileDescriptor_expiry_data_8c4813946863f3d0, []int{0}
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
	Map                  map[string]*TxNums `protobuf:"bytes,1,rep,name=map" json:"map,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Collections) Reset()         { *m = Collections{} }
func (m *Collections) String() string { return proto.CompactTextString(m) }
func (*Collections) ProtoMessage()    {}
func (*Collections) Descriptor() ([]byte, []int) {
	return fileDescriptor_expiry_data_8c4813946863f3d0, []int{1}
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
	return fileDescriptor_expiry_data_8c4813946863f3d0, []int{2}
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
	proto.RegisterType((*TxNums)(nil), "pvtdatastorage.TxNums")
}

func init() { proto.RegisterFile("expiry_data.proto", fileDescriptor_expiry_data_8c4813946863f3d0) }

var fileDescriptor_expiry_data_8c4813946863f3d0 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xad, 0x28, 0xc8,
	0x2c, 0xaa, 0x8c, 0x4f, 0x49, 0x2c, 0x49, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2b,
	0x28, 0x2b, 0x01, 0x71, 0x8b, 0x4b, 0xf2, 0x8b, 0x12, 0xd3, 0x53, 0x95, 0x66, 0x30, 0x72, 0x71,
	0xb9, 0x82, 0x55, 0xb9, 0x24, 0x96, 0x24, 0x0a, 0x99, 0x72, 0x31, 0xe7, 0x26, 0x16, 0x48, 0x30,
	0x2a, 0x30, 0x6b, 0x70, 0x1b, 0x29, 0xeb, 0xa1, 0x2a, 0xd6, 0x43, 0x28, 0xd4, 0xf3, 0x4d, 0x2c,
	0x70, 0xcd, 0x2b, 0x29, 0xaa, 0x0c, 0x02, 0xa9, 0x97, 0x0a, 0xe6, 0xe2, 0x80, 0x09, 0x08, 0x09,
	0x70, 0x31, 0x67, 0xa7, 0x56, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0x98, 0x42, 0x86,
	0x5c, 0xac, 0x65, 0x89, 0x39, 0xa5, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0xd2, 0xe8,
	0xc6, 0x3a, 0xe7, 0xe7, 0xe4, 0xa4, 0x26, 0x97, 0x64, 0xe6, 0xe7, 0x15, 0x07, 0x41, 0x54, 0x5a,
	0x31, 0x59, 0x30, 0x2a, 0x4d, 0x65, 0xe4, 0xe2, 0x46, 0x92, 0x12, 0x32, 0x43, 0x76, 0x9b, 0x0a,
	0x1e, 0x43, 0xd0, 0x1c, 0xe7, 0x87, 0xd7, 0x71, 0x3a, 0xa8, 0x8e, 0x13, 0x43, 0x37, 0x37, 0xa4,
	0xc2, 0xaf, 0x34, 0x17, 0xc5, 0x5d, 0x32, 0x5c, 0x6c, 0x10, 0x41, 0x21, 0x21, 0x2e, 0x96, 0x9c,
	0xcc, 0xe2, 0x12, 0xb0, 0x93, 0x58, 0x82, 0xc0, 0x6c, 0x27, 0xab, 0x28, 0x8b, 0xf4, 0xcc, 0x92,
	0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0x8c, 0xca, 0x82, 0xd4, 0xa2, 0x9c, 0xd4, 0x94,
	0xf4, 0xd4, 0x22, 0xfd, 0xb4, 0xc4, 0xa4, 0xa2, 0xcc, 0x64, 0xfd, 0xe4, 0xfc, 0xa2, 0x54, 0x7d,
	0xa8, 0x10, 0xaa, 0x5d, 0x49, 0x6c, 0xe0, 0x38, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xb7,
	0xd4, 0xf8, 0x9d, 0xb8, 0x01, 0x00, 0x00,
}
