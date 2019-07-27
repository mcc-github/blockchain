


package kvrwset

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 



type KVRWSet struct {
	Reads                []*KVRead          `protobuf:"bytes,1,rep,name=reads,proto3" json:"reads,omitempty"`
	RangeQueriesInfo     []*RangeQueryInfo  `protobuf:"bytes,2,rep,name=range_queries_info,json=rangeQueriesInfo,proto3" json:"range_queries_info,omitempty"`
	Writes               []*KVWrite         `protobuf:"bytes,3,rep,name=writes,proto3" json:"writes,omitempty"`
	MetadataWrites       []*KVMetadataWrite `protobuf:"bytes,4,rep,name=metadata_writes,json=metadataWrites,proto3" json:"metadata_writes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *KVRWSet) Reset()         { *m = KVRWSet{} }
func (m *KVRWSet) String() string { return proto.CompactTextString(m) }
func (*KVRWSet) ProtoMessage()    {}
func (*KVRWSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{0}
}

func (m *KVRWSet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVRWSet.Unmarshal(m, b)
}
func (m *KVRWSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVRWSet.Marshal(b, m, deterministic)
}
func (m *KVRWSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVRWSet.Merge(m, src)
}
func (m *KVRWSet) XXX_Size() int {
	return xxx_messageInfo_KVRWSet.Size(m)
}
func (m *KVRWSet) XXX_DiscardUnknown() {
	xxx_messageInfo_KVRWSet.DiscardUnknown(m)
}

var xxx_messageInfo_KVRWSet proto.InternalMessageInfo

func (m *KVRWSet) GetReads() []*KVRead {
	if m != nil {
		return m.Reads
	}
	return nil
}

func (m *KVRWSet) GetRangeQueriesInfo() []*RangeQueryInfo {
	if m != nil {
		return m.RangeQueriesInfo
	}
	return nil
}

func (m *KVRWSet) GetWrites() []*KVWrite {
	if m != nil {
		return m.Writes
	}
	return nil
}

func (m *KVRWSet) GetMetadataWrites() []*KVMetadataWrite {
	if m != nil {
		return m.MetadataWrites
	}
	return nil
}


type HashedRWSet struct {
	HashedReads          []*KVReadHash          `protobuf:"bytes,1,rep,name=hashed_reads,json=hashedReads,proto3" json:"hashed_reads,omitempty"`
	HashedWrites         []*KVWriteHash         `protobuf:"bytes,2,rep,name=hashed_writes,json=hashedWrites,proto3" json:"hashed_writes,omitempty"`
	MetadataWrites       []*KVMetadataWriteHash `protobuf:"bytes,3,rep,name=metadata_writes,json=metadataWrites,proto3" json:"metadata_writes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *HashedRWSet) Reset()         { *m = HashedRWSet{} }
func (m *HashedRWSet) String() string { return proto.CompactTextString(m) }
func (*HashedRWSet) ProtoMessage()    {}
func (*HashedRWSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{1}
}

func (m *HashedRWSet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HashedRWSet.Unmarshal(m, b)
}
func (m *HashedRWSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HashedRWSet.Marshal(b, m, deterministic)
}
func (m *HashedRWSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HashedRWSet.Merge(m, src)
}
func (m *HashedRWSet) XXX_Size() int {
	return xxx_messageInfo_HashedRWSet.Size(m)
}
func (m *HashedRWSet) XXX_DiscardUnknown() {
	xxx_messageInfo_HashedRWSet.DiscardUnknown(m)
}

var xxx_messageInfo_HashedRWSet proto.InternalMessageInfo

func (m *HashedRWSet) GetHashedReads() []*KVReadHash {
	if m != nil {
		return m.HashedReads
	}
	return nil
}

func (m *HashedRWSet) GetHashedWrites() []*KVWriteHash {
	if m != nil {
		return m.HashedWrites
	}
	return nil
}

func (m *HashedRWSet) GetMetadataWrites() []*KVMetadataWriteHash {
	if m != nil {
		return m.MetadataWrites
	}
	return nil
}



type KVRead struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Version              *Version `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KVRead) Reset()         { *m = KVRead{} }
func (m *KVRead) String() string { return proto.CompactTextString(m) }
func (*KVRead) ProtoMessage()    {}
func (*KVRead) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{2}
}

func (m *KVRead) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVRead.Unmarshal(m, b)
}
func (m *KVRead) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVRead.Marshal(b, m, deterministic)
}
func (m *KVRead) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVRead.Merge(m, src)
}
func (m *KVRead) XXX_Size() int {
	return xxx_messageInfo_KVRead.Size(m)
}
func (m *KVRead) XXX_DiscardUnknown() {
	xxx_messageInfo_KVRead.DiscardUnknown(m)
}

var xxx_messageInfo_KVRead proto.InternalMessageInfo

func (m *KVRead) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KVRead) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}


type KVWrite struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	IsDelete             bool     `protobuf:"varint,2,opt,name=is_delete,json=isDelete,proto3" json:"is_delete,omitempty"`
	Value                []byte   `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KVWrite) Reset()         { *m = KVWrite{} }
func (m *KVWrite) String() string { return proto.CompactTextString(m) }
func (*KVWrite) ProtoMessage()    {}
func (*KVWrite) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{3}
}

func (m *KVWrite) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVWrite.Unmarshal(m, b)
}
func (m *KVWrite) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVWrite.Marshal(b, m, deterministic)
}
func (m *KVWrite) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVWrite.Merge(m, src)
}
func (m *KVWrite) XXX_Size() int {
	return xxx_messageInfo_KVWrite.Size(m)
}
func (m *KVWrite) XXX_DiscardUnknown() {
	xxx_messageInfo_KVWrite.DiscardUnknown(m)
}

var xxx_messageInfo_KVWrite proto.InternalMessageInfo

func (m *KVWrite) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KVWrite) GetIsDelete() bool {
	if m != nil {
		return m.IsDelete
	}
	return false
}

func (m *KVWrite) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}


type KVMetadataWrite struct {
	Key                  string             `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Entries              []*KVMetadataEntry `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *KVMetadataWrite) Reset()         { *m = KVMetadataWrite{} }
func (m *KVMetadataWrite) String() string { return proto.CompactTextString(m) }
func (*KVMetadataWrite) ProtoMessage()    {}
func (*KVMetadataWrite) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{4}
}

func (m *KVMetadataWrite) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVMetadataWrite.Unmarshal(m, b)
}
func (m *KVMetadataWrite) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVMetadataWrite.Marshal(b, m, deterministic)
}
func (m *KVMetadataWrite) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVMetadataWrite.Merge(m, src)
}
func (m *KVMetadataWrite) XXX_Size() int {
	return xxx_messageInfo_KVMetadataWrite.Size(m)
}
func (m *KVMetadataWrite) XXX_DiscardUnknown() {
	xxx_messageInfo_KVMetadataWrite.DiscardUnknown(m)
}

var xxx_messageInfo_KVMetadataWrite proto.InternalMessageInfo

func (m *KVMetadataWrite) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KVMetadataWrite) GetEntries() []*KVMetadataEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}




type KVReadHash struct {
	KeyHash              []byte   `protobuf:"bytes,1,opt,name=key_hash,json=keyHash,proto3" json:"key_hash,omitempty"`
	Version              *Version `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KVReadHash) Reset()         { *m = KVReadHash{} }
func (m *KVReadHash) String() string { return proto.CompactTextString(m) }
func (*KVReadHash) ProtoMessage()    {}
func (*KVReadHash) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{5}
}

func (m *KVReadHash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVReadHash.Unmarshal(m, b)
}
func (m *KVReadHash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVReadHash.Marshal(b, m, deterministic)
}
func (m *KVReadHash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVReadHash.Merge(m, src)
}
func (m *KVReadHash) XXX_Size() int {
	return xxx_messageInfo_KVReadHash.Size(m)
}
func (m *KVReadHash) XXX_DiscardUnknown() {
	xxx_messageInfo_KVReadHash.DiscardUnknown(m)
}

var xxx_messageInfo_KVReadHash proto.InternalMessageInfo

func (m *KVReadHash) GetKeyHash() []byte {
	if m != nil {
		return m.KeyHash
	}
	return nil
}

func (m *KVReadHash) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}


type KVWriteHash struct {
	KeyHash              []byte   `protobuf:"bytes,1,opt,name=key_hash,json=keyHash,proto3" json:"key_hash,omitempty"`
	IsDelete             bool     `protobuf:"varint,2,opt,name=is_delete,json=isDelete,proto3" json:"is_delete,omitempty"`
	ValueHash            []byte   `protobuf:"bytes,3,opt,name=value_hash,json=valueHash,proto3" json:"value_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KVWriteHash) Reset()         { *m = KVWriteHash{} }
func (m *KVWriteHash) String() string { return proto.CompactTextString(m) }
func (*KVWriteHash) ProtoMessage()    {}
func (*KVWriteHash) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{6}
}

func (m *KVWriteHash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVWriteHash.Unmarshal(m, b)
}
func (m *KVWriteHash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVWriteHash.Marshal(b, m, deterministic)
}
func (m *KVWriteHash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVWriteHash.Merge(m, src)
}
func (m *KVWriteHash) XXX_Size() int {
	return xxx_messageInfo_KVWriteHash.Size(m)
}
func (m *KVWriteHash) XXX_DiscardUnknown() {
	xxx_messageInfo_KVWriteHash.DiscardUnknown(m)
}

var xxx_messageInfo_KVWriteHash proto.InternalMessageInfo

func (m *KVWriteHash) GetKeyHash() []byte {
	if m != nil {
		return m.KeyHash
	}
	return nil
}

func (m *KVWriteHash) GetIsDelete() bool {
	if m != nil {
		return m.IsDelete
	}
	return false
}

func (m *KVWriteHash) GetValueHash() []byte {
	if m != nil {
		return m.ValueHash
	}
	return nil
}


type KVMetadataWriteHash struct {
	KeyHash              []byte             `protobuf:"bytes,1,opt,name=key_hash,json=keyHash,proto3" json:"key_hash,omitempty"`
	Entries              []*KVMetadataEntry `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *KVMetadataWriteHash) Reset()         { *m = KVMetadataWriteHash{} }
func (m *KVMetadataWriteHash) String() string { return proto.CompactTextString(m) }
func (*KVMetadataWriteHash) ProtoMessage()    {}
func (*KVMetadataWriteHash) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{7}
}

func (m *KVMetadataWriteHash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVMetadataWriteHash.Unmarshal(m, b)
}
func (m *KVMetadataWriteHash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVMetadataWriteHash.Marshal(b, m, deterministic)
}
func (m *KVMetadataWriteHash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVMetadataWriteHash.Merge(m, src)
}
func (m *KVMetadataWriteHash) XXX_Size() int {
	return xxx_messageInfo_KVMetadataWriteHash.Size(m)
}
func (m *KVMetadataWriteHash) XXX_DiscardUnknown() {
	xxx_messageInfo_KVMetadataWriteHash.DiscardUnknown(m)
}

var xxx_messageInfo_KVMetadataWriteHash proto.InternalMessageInfo

func (m *KVMetadataWriteHash) GetKeyHash() []byte {
	if m != nil {
		return m.KeyHash
	}
	return nil
}

func (m *KVMetadataWriteHash) GetEntries() []*KVMetadataEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}


type KVMetadataEntry struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KVMetadataEntry) Reset()         { *m = KVMetadataEntry{} }
func (m *KVMetadataEntry) String() string { return proto.CompactTextString(m) }
func (*KVMetadataEntry) ProtoMessage()    {}
func (*KVMetadataEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{8}
}

func (m *KVMetadataEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVMetadataEntry.Unmarshal(m, b)
}
func (m *KVMetadataEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVMetadataEntry.Marshal(b, m, deterministic)
}
func (m *KVMetadataEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVMetadataEntry.Merge(m, src)
}
func (m *KVMetadataEntry) XXX_Size() int {
	return xxx_messageInfo_KVMetadataEntry.Size(m)
}
func (m *KVMetadataEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_KVMetadataEntry.DiscardUnknown(m)
}

var xxx_messageInfo_KVMetadataEntry proto.InternalMessageInfo

func (m *KVMetadataEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *KVMetadataEntry) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}





type Version struct {
	BlockNum             uint64   `protobuf:"varint,1,opt,name=block_num,json=blockNum,proto3" json:"block_num,omitempty"`
	TxNum                uint64   `protobuf:"varint,2,opt,name=tx_num,json=txNum,proto3" json:"tx_num,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{9}
}

func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (m *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(m, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetBlockNum() uint64 {
	if m != nil {
		return m.BlockNum
	}
	return 0
}

func (m *Version) GetTxNum() uint64 {
	if m != nil {
		return m.TxNum
	}
	return 0
}







type RangeQueryInfo struct {
	StartKey     string `protobuf:"bytes,1,opt,name=start_key,json=startKey,proto3" json:"start_key,omitempty"`
	EndKey       string `protobuf:"bytes,2,opt,name=end_key,json=endKey,proto3" json:"end_key,omitempty"`
	ItrExhausted bool   `protobuf:"varint,3,opt,name=itr_exhausted,json=itrExhausted,proto3" json:"itr_exhausted,omitempty"`
	
	
	
	ReadsInfo            isRangeQueryInfo_ReadsInfo `protobuf_oneof:"reads_info"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *RangeQueryInfo) Reset()         { *m = RangeQueryInfo{} }
func (m *RangeQueryInfo) String() string { return proto.CompactTextString(m) }
func (*RangeQueryInfo) ProtoMessage()    {}
func (*RangeQueryInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{10}
}

func (m *RangeQueryInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RangeQueryInfo.Unmarshal(m, b)
}
func (m *RangeQueryInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RangeQueryInfo.Marshal(b, m, deterministic)
}
func (m *RangeQueryInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RangeQueryInfo.Merge(m, src)
}
func (m *RangeQueryInfo) XXX_Size() int {
	return xxx_messageInfo_RangeQueryInfo.Size(m)
}
func (m *RangeQueryInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RangeQueryInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RangeQueryInfo proto.InternalMessageInfo

func (m *RangeQueryInfo) GetStartKey() string {
	if m != nil {
		return m.StartKey
	}
	return ""
}

func (m *RangeQueryInfo) GetEndKey() string {
	if m != nil {
		return m.EndKey
	}
	return ""
}

func (m *RangeQueryInfo) GetItrExhausted() bool {
	if m != nil {
		return m.ItrExhausted
	}
	return false
}

type isRangeQueryInfo_ReadsInfo interface {
	isRangeQueryInfo_ReadsInfo()
}

type RangeQueryInfo_RawReads struct {
	RawReads *QueryReads `protobuf:"bytes,4,opt,name=raw_reads,json=rawReads,proto3,oneof"`
}

type RangeQueryInfo_ReadsMerkleHashes struct {
	ReadsMerkleHashes *QueryReadsMerkleSummary `protobuf:"bytes,5,opt,name=reads_merkle_hashes,json=readsMerkleHashes,proto3,oneof"`
}

func (*RangeQueryInfo_RawReads) isRangeQueryInfo_ReadsInfo() {}

func (*RangeQueryInfo_ReadsMerkleHashes) isRangeQueryInfo_ReadsInfo() {}

func (m *RangeQueryInfo) GetReadsInfo() isRangeQueryInfo_ReadsInfo {
	if m != nil {
		return m.ReadsInfo
	}
	return nil
}

func (m *RangeQueryInfo) GetRawReads() *QueryReads {
	if x, ok := m.GetReadsInfo().(*RangeQueryInfo_RawReads); ok {
		return x.RawReads
	}
	return nil
}

func (m *RangeQueryInfo) GetReadsMerkleHashes() *QueryReadsMerkleSummary {
	if x, ok := m.GetReadsInfo().(*RangeQueryInfo_ReadsMerkleHashes); ok {
		return x.ReadsMerkleHashes
	}
	return nil
}


func (*RangeQueryInfo) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*RangeQueryInfo_RawReads)(nil),
		(*RangeQueryInfo_ReadsMerkleHashes)(nil),
	}
}


type QueryReads struct {
	KvReads              []*KVRead `protobuf:"bytes,1,rep,name=kv_reads,json=kvReads,proto3" json:"kv_reads,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *QueryReads) Reset()         { *m = QueryReads{} }
func (m *QueryReads) String() string { return proto.CompactTextString(m) }
func (*QueryReads) ProtoMessage()    {}
func (*QueryReads) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{11}
}

func (m *QueryReads) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryReads.Unmarshal(m, b)
}
func (m *QueryReads) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryReads.Marshal(b, m, deterministic)
}
func (m *QueryReads) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryReads.Merge(m, src)
}
func (m *QueryReads) XXX_Size() int {
	return xxx_messageInfo_QueryReads.Size(m)
}
func (m *QueryReads) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryReads.DiscardUnknown(m)
}

var xxx_messageInfo_QueryReads proto.InternalMessageInfo

func (m *QueryReads) GetKvReads() []*KVRead {
	if m != nil {
		return m.KvReads
	}
	return nil
}






type QueryReadsMerkleSummary struct {
	MaxDegree            uint32   `protobuf:"varint,1,opt,name=max_degree,json=maxDegree,proto3" json:"max_degree,omitempty"`
	MaxLevel             uint32   `protobuf:"varint,2,opt,name=max_level,json=maxLevel,proto3" json:"max_level,omitempty"`
	MaxLevelHashes       [][]byte `protobuf:"bytes,3,rep,name=max_level_hashes,json=maxLevelHashes,proto3" json:"max_level_hashes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryReadsMerkleSummary) Reset()         { *m = QueryReadsMerkleSummary{} }
func (m *QueryReadsMerkleSummary) String() string { return proto.CompactTextString(m) }
func (*QueryReadsMerkleSummary) ProtoMessage()    {}
func (*QueryReadsMerkleSummary) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee5d686eab23a142, []int{12}
}

func (m *QueryReadsMerkleSummary) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryReadsMerkleSummary.Unmarshal(m, b)
}
func (m *QueryReadsMerkleSummary) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryReadsMerkleSummary.Marshal(b, m, deterministic)
}
func (m *QueryReadsMerkleSummary) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryReadsMerkleSummary.Merge(m, src)
}
func (m *QueryReadsMerkleSummary) XXX_Size() int {
	return xxx_messageInfo_QueryReadsMerkleSummary.Size(m)
}
func (m *QueryReadsMerkleSummary) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryReadsMerkleSummary.DiscardUnknown(m)
}

var xxx_messageInfo_QueryReadsMerkleSummary proto.InternalMessageInfo

func (m *QueryReadsMerkleSummary) GetMaxDegree() uint32 {
	if m != nil {
		return m.MaxDegree
	}
	return 0
}

func (m *QueryReadsMerkleSummary) GetMaxLevel() uint32 {
	if m != nil {
		return m.MaxLevel
	}
	return 0
}

func (m *QueryReadsMerkleSummary) GetMaxLevelHashes() [][]byte {
	if m != nil {
		return m.MaxLevelHashes
	}
	return nil
}

func init() {
	proto.RegisterType((*KVRWSet)(nil), "kvrwset.KVRWSet")
	proto.RegisterType((*HashedRWSet)(nil), "kvrwset.HashedRWSet")
	proto.RegisterType((*KVRead)(nil), "kvrwset.KVRead")
	proto.RegisterType((*KVWrite)(nil), "kvrwset.KVWrite")
	proto.RegisterType((*KVMetadataWrite)(nil), "kvrwset.KVMetadataWrite")
	proto.RegisterType((*KVReadHash)(nil), "kvrwset.KVReadHash")
	proto.RegisterType((*KVWriteHash)(nil), "kvrwset.KVWriteHash")
	proto.RegisterType((*KVMetadataWriteHash)(nil), "kvrwset.KVMetadataWriteHash")
	proto.RegisterType((*KVMetadataEntry)(nil), "kvrwset.KVMetadataEntry")
	proto.RegisterType((*Version)(nil), "kvrwset.Version")
	proto.RegisterType((*RangeQueryInfo)(nil), "kvrwset.RangeQueryInfo")
	proto.RegisterType((*QueryReads)(nil), "kvrwset.QueryReads")
	proto.RegisterType((*QueryReadsMerkleSummary)(nil), "kvrwset.QueryReadsMerkleSummary")
}

func init() {
	proto.RegisterFile("ledger/rwset/kvrwset/kv_rwset.proto", fileDescriptor_ee5d686eab23a142)
}

var fileDescriptor_ee5d686eab23a142 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0xdf, 0x6b, 0xdb, 0x48,
	0x10, 0x8e, 0x7f, 0xcb, 0x63, 0x3b, 0xf1, 0x6d, 0x72, 0x44, 0xc7, 0xdd, 0x81, 0x51, 0x38, 0x30,
	0x79, 0xb0, 0xc1, 0x07, 0xc7, 0x85, 0xe3, 0x1e, 0x5a, 0xe2, 0x92, 0x92, 0x26, 0xd0, 0x0d, 0x24,
	0xd0, 0x17, 0xb1, 0x8e, 0x26, 0xb6, 0xb0, 0x25, 0xa5, 0xab, 0x95, 0x6d, 0x3d, 0x95, 0xfe, 0x75,
	0xfd, 0x47, 0xfa, 0x87, 0x94, 0x9d, 0x95, 0x63, 0xc5, 0x75, 0x0c, 0xed, 0x93, 0xb5, 0xf3, 0xcd,
	0x37, 0x3b, 0xdf, 0x37, 0xde, 0x5d, 0x38, 0x99, 0xa1, 0x37, 0x46, 0xd9, 0x97, 0x8b, 0x18, 0x55,
	0x7f, 0x3a, 0x5f, 0xfd, 0xba, 0xf4, 0xd1, 0x7b, 0x94, 0x91, 0x8a, 0x58, 0x2d, 0x8b, 0x3b, 0x5f,
	0x0b, 0x50, 0xbb, 0xbc, 0xe5, 0x77, 0x37, 0xa8, 0xd8, 0x5f, 0x50, 0x91, 0x28, 0xbc, 0xd8, 0x2e,
	0x74, 0x4a, 0xdd, 0xc6, 0xe0, 0xa0, 0x97, 0x25, 0xf5, 0x2e, 0x6f, 0x39, 0x0a, 0x8f, 0x1b, 0x94,
	0x0d, 0x81, 0x49, 0x11, 0x8e, 0xd1, 0xfd, 0x98, 0xa0, 0xf4, 0x31, 0x76, 0xfd, 0xf0, 0x21, 0xb2,
	0x8b, 0xc4, 0x39, 0x7e, 0xe2, 0x70, 0x9d, 0xf2, 0x3e, 0x41, 0x99, 0xbe, 0x0d, 0x1f, 0x22, 0xde,
	0x96, 0xab, 0xb5, 0x8f, 0xb1, 0x8e, 0xb0, 0x2e, 0x54, 0x17, 0xd2, 0x57, 0x18, 0xdb, 0x25, 0xa2,
	0xb6, 0x73, 0xdb, 0xdd, 0x69, 0x80, 0x67, 0x38, 0x7b, 0x05, 0x07, 0x01, 0x2a, 0xe1, 0x09, 0x25,
	0xdc, 0x8c, 0x52, 0x26, 0x8a, 0x9d, 0xa3, 0x5c, 0x65, 0x19, 0x86, 0xba, 0x1f, 0xe4, 0x97, 0xb1,
	0xf3, 0xa5, 0x00, 0x8d, 0x0b, 0x11, 0x4f, 0xd0, 0x33, 0x52, 0xff, 0x81, 0xe6, 0x84, 0x96, 0x6e,
	0x5e, 0xf1, 0xe1, 0x86, 0x62, 0xcd, 0xe0, 0x0d, 0x93, 0xc8, 0x49, 0xfb, 0x19, 0xb4, 0x32, 0x5e,
	0xd6, 0x88, 0x91, 0x7d, 0xb4, 0xd9, 0x3b, 0x31, 0xb3, 0x2d, 0x4c, 0x0b, 0x6c, 0xf8, 0xbd, 0x0a,
	0x23, 0xfc, 0x8f, 0x97, 0x54, 0x50, 0x91, 0x4d, 0x25, 0x6f, 0xa0, 0x6a, 0x9a, 0x63, 0x6d, 0x28,
	0x4d, 0x31, 0xb5, 0x0b, 0x9d, 0x42, 0xb7, 0xce, 0xf5, 0x27, 0x3b, 0x85, 0xda, 0x1c, 0x65, 0xec,
	0x47, 0xa1, 0x5d, 0xec, 0x14, 0x9e, 0x79, 0x7a, 0x6b, 0xe2, 0x7c, 0x95, 0xe0, 0x5c, 0xeb, 0xb9,
	0x53, 0xcd, 0x2d, 0x85, 0x7e, 0x87, 0xba, 0x1f, 0xbb, 0x1e, 0xce, 0x50, 0x21, 0x95, 0xb2, 0xb8,
	0xe5, 0xc7, 0xe7, 0xb4, 0x66, 0x47, 0x50, 0x99, 0x8b, 0x59, 0x82, 0x76, 0xa9, 0x53, 0xe8, 0x36,
	0xb9, 0x59, 0x38, 0x77, 0x70, 0xb0, 0xd1, 0xfe, 0x96, 0xba, 0x03, 0xa8, 0x61, 0xa8, 0xf4, 0x5f,
	0x20, 0x33, 0x6e, 0xdb, 0x04, 0x87, 0xa1, 0x92, 0x29, 0x5f, 0x25, 0x3a, 0x37, 0x00, 0xeb, 0x69,
	0xb0, 0xdf, 0xc0, 0x9a, 0x62, 0xea, 0x6a, 0x67, 0xa9, 0x70, 0x93, 0xd7, 0xa6, 0x98, 0x12, 0xf4,
	0x23, 0xea, 0x3d, 0x68, 0xe4, 0x26, 0xb5, 0xab, 0xea, 0x4e, 0x2b, 0xfe, 0x04, 0x20, 0xf5, 0x86,
	0x69, 0xfc, 0xa8, 0x53, 0x44, 0x73, 0x1d, 0x0f, 0x0e, 0xb7, 0x8c, 0x74, 0xd7, 0x6e, 0x3f, 0x63,
	0xd0, 0x7f, 0x79, 0xe7, 0x09, 0x63, 0x0c, 0xca, 0xa1, 0x08, 0x30, 0xb3, 0x9e, 0xbe, 0xd7, 0x63,
	0x2b, 0xe6, 0xc7, 0xf6, 0x3f, 0xd4, 0x32, 0x73, 0xb4, 0xd2, 0xd1, 0x2c, 0xba, 0x9f, 0xba, 0x61,
	0x12, 0x10, 0xb3, 0xcc, 0x2d, 0x0a, 0x5c, 0x27, 0x01, 0xfb, 0x15, 0xaa, 0x6a, 0x49, 0x48, 0x91,
	0x90, 0x8a, 0x5a, 0x5e, 0x27, 0x81, 0xf3, 0xb9, 0x08, 0xfb, 0xcf, 0x4f, 0xba, 0x2e, 0x13, 0x2b,
	0x21, 0x95, 0xbb, 0x9e, 0xbd, 0x45, 0x81, 0x4b, 0x4c, 0xd9, 0xb1, 0xd6, 0xe7, 0x11, 0x54, 0x24,
	0xa8, 0x8a, 0xa1, 0xa7, 0x81, 0x13, 0x68, 0xf9, 0x4a, 0xba, 0xb8, 0x9c, 0x88, 0x24, 0x56, 0xe8,
	0x91, 0x99, 0x16, 0x6f, 0xfa, 0x4a, 0x0e, 0x57, 0x31, 0x36, 0x80, 0xba, 0x14, 0x8b, 0xec, 0xc8,
	0x96, 0x69, 0xc6, 0xeb, 0x23, 0x4b, 0x1d, 0xd0, 0x29, 0xbd, 0xd8, 0xe3, 0x96, 0x14, 0x0b, 0x73,
	0x62, 0x39, 0x1c, 0x52, 0xbe, 0x1b, 0xa0, 0x9c, 0xce, 0xcc, 0xa4, 0x30, 0xb6, 0x2b, 0xc4, 0xee,
	0x6c, 0x61, 0x5f, 0x51, 0xde, 0x4d, 0x12, 0x04, 0x42, 0xa6, 0x17, 0x7b, 0xfc, 0x17, 0xb9, 0x8e,
	0xd2, 0x15, 0x12, 0xbf, 0x6e, 0x02, 0x98, 0x9a, 0xfa, 0xe6, 0x73, 0xfe, 0x05, 0x58, 0xb3, 0xd9,
	0x29, 0x58, 0xfa, 0xae, 0xdd, 0x75, 0x8f, 0xd6, 0xa6, 0x73, 0xca, 0x75, 0x3e, 0xc1, 0xf1, 0x0b,
	0xfb, 0xea, 0x7f, 0x56, 0x20, 0x96, 0xae, 0x87, 0x63, 0x89, 0x66, 0x8e, 0x2d, 0x5e, 0x0f, 0xc4,
	0xf2, 0x9c, 0x02, 0xda, 0x64, 0x0d, 0xcf, 0x70, 0x8e, 0x33, 0x72, 0xb2, 0xc5, 0xad, 0x40, 0x2c,
	0xdf, 0xe9, 0x35, 0xeb, 0x42, 0xfb, 0x09, 0x5c, 0xe9, 0xd5, 0x57, 0x4d, 0x93, 0xef, 0xaf, 0x72,
	0x32, 0x21, 0x11, 0x0c, 0x22, 0x39, 0xee, 0x4d, 0xd2, 0x47, 0x94, 0xe6, 0xd9, 0xe8, 0x3d, 0x88,
	0x91, 0xf4, 0xef, 0xcd, 0x33, 0x11, 0xf7, 0xb2, 0xa0, 0x69, 0x3f, 0x93, 0xf1, 0xe1, 0x6c, 0xec,
	0xab, 0x49, 0x32, 0xea, 0xdd, 0x47, 0x41, 0x3f, 0x47, 0xed, 0x1b, 0x6a, 0xdf, 0x50, 0xfb, 0xdb,
	0x9e, 0xa1, 0x51, 0x95, 0xc0, 0xbf, 0xbf, 0x05, 0x00, 0x00, 0xff, 0xff, 0xe3, 0xe1, 0xb5, 0x07,
	0xa5, 0x06, 0x00, 0x00,
}
