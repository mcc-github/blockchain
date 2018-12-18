


package common 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 


type Status int32

const (
	Status_UNKNOWN                  Status = 0
	Status_SUCCESS                  Status = 200
	Status_BAD_REQUEST              Status = 400
	Status_FORBIDDEN                Status = 403
	Status_NOT_FOUND                Status = 404
	Status_REQUEST_ENTITY_TOO_LARGE Status = 413
	Status_INTERNAL_SERVER_ERROR    Status = 500
	Status_NOT_IMPLEMENTED          Status = 501
	Status_SERVICE_UNAVAILABLE      Status = 503
)

var Status_name = map[int32]string{
	0:   "UNKNOWN",
	200: "SUCCESS",
	400: "BAD_REQUEST",
	403: "FORBIDDEN",
	404: "NOT_FOUND",
	413: "REQUEST_ENTITY_TOO_LARGE",
	500: "INTERNAL_SERVER_ERROR",
	501: "NOT_IMPLEMENTED",
	503: "SERVICE_UNAVAILABLE",
}
var Status_value = map[string]int32{
	"UNKNOWN":                  0,
	"SUCCESS":                  200,
	"BAD_REQUEST":              400,
	"FORBIDDEN":                403,
	"NOT_FOUND":                404,
	"REQUEST_ENTITY_TOO_LARGE": 413,
	"INTERNAL_SERVER_ERROR":    500,
	"NOT_IMPLEMENTED":          501,
	"SERVICE_UNAVAILABLE":      503,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}
func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{0}
}

type HeaderType int32

const (
	HeaderType_MESSAGE                    HeaderType = 0
	HeaderType_CONFIG                     HeaderType = 1
	HeaderType_CONFIG_UPDATE              HeaderType = 2
	HeaderType_ENDORSER_TRANSACTION       HeaderType = 3
	HeaderType_ORDERER_TRANSACTION        HeaderType = 4
	HeaderType_DELIVER_SEEK_INFO          HeaderType = 5
	HeaderType_CHAINCODE_PACKAGE          HeaderType = 6
	HeaderType_PEER_ADMIN_OPERATION       HeaderType = 8
	HeaderType_TOKEN_TRANSACTION          HeaderType = 9
	HeaderType_TOKEN_ENDORSER_TRANSACTION HeaderType = 10
)

var HeaderType_name = map[int32]string{
	0:  "MESSAGE",
	1:  "CONFIG",
	2:  "CONFIG_UPDATE",
	3:  "ENDORSER_TRANSACTION",
	4:  "ORDERER_TRANSACTION",
	5:  "DELIVER_SEEK_INFO",
	6:  "CHAINCODE_PACKAGE",
	8:  "PEER_ADMIN_OPERATION",
	9:  "TOKEN_TRANSACTION",
	10: "TOKEN_ENDORSER_TRANSACTION",
}
var HeaderType_value = map[string]int32{
	"MESSAGE":                    0,
	"CONFIG":                     1,
	"CONFIG_UPDATE":              2,
	"ENDORSER_TRANSACTION":       3,
	"ORDERER_TRANSACTION":        4,
	"DELIVER_SEEK_INFO":          5,
	"CHAINCODE_PACKAGE":          6,
	"PEER_ADMIN_OPERATION":       8,
	"TOKEN_TRANSACTION":          9,
	"TOKEN_ENDORSER_TRANSACTION": 10,
}

func (x HeaderType) String() string {
	return proto.EnumName(HeaderType_name, int32(x))
}
func (HeaderType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{1}
}


type BlockMetadataIndex int32

const (
	BlockMetadataIndex_SIGNATURES          BlockMetadataIndex = 0
	BlockMetadataIndex_LAST_CONFIG         BlockMetadataIndex = 1
	BlockMetadataIndex_TRANSACTIONS_FILTER BlockMetadataIndex = 2
	BlockMetadataIndex_ORDERER             BlockMetadataIndex = 3
)

var BlockMetadataIndex_name = map[int32]string{
	0: "SIGNATURES",
	1: "LAST_CONFIG",
	2: "TRANSACTIONS_FILTER",
	3: "ORDERER",
}
var BlockMetadataIndex_value = map[string]int32{
	"SIGNATURES":          0,
	"LAST_CONFIG":         1,
	"TRANSACTIONS_FILTER": 2,
	"ORDERER":             3,
}

func (x BlockMetadataIndex) String() string {
	return proto.EnumName(BlockMetadataIndex_name, int32(x))
}
func (BlockMetadataIndex) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{2}
}


type LastConfig struct {
	Index                uint64   `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LastConfig) Reset()         { *m = LastConfig{} }
func (m *LastConfig) String() string { return proto.CompactTextString(m) }
func (*LastConfig) ProtoMessage()    {}
func (*LastConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{0}
}
func (m *LastConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LastConfig.Unmarshal(m, b)
}
func (m *LastConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LastConfig.Marshal(b, m, deterministic)
}
func (dst *LastConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LastConfig.Merge(dst, src)
}
func (m *LastConfig) XXX_Size() int {
	return xxx_messageInfo_LastConfig.Size(m)
}
func (m *LastConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_LastConfig.DiscardUnknown(m)
}

var xxx_messageInfo_LastConfig proto.InternalMessageInfo

func (m *LastConfig) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}


type Metadata struct {
	Value                []byte               `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Signatures           []*MetadataSignature `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{1}
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

func (m *Metadata) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Metadata) GetSignatures() []*MetadataSignature {
	if m != nil {
		return m.Signatures
	}
	return nil
}

type MetadataSignature struct {
	SignatureHeader      []byte   `protobuf:"bytes,1,opt,name=signature_header,json=signatureHeader,proto3" json:"signature_header,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetadataSignature) Reset()         { *m = MetadataSignature{} }
func (m *MetadataSignature) String() string { return proto.CompactTextString(m) }
func (*MetadataSignature) ProtoMessage()    {}
func (*MetadataSignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{2}
}
func (m *MetadataSignature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataSignature.Unmarshal(m, b)
}
func (m *MetadataSignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataSignature.Marshal(b, m, deterministic)
}
func (dst *MetadataSignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataSignature.Merge(dst, src)
}
func (m *MetadataSignature) XXX_Size() int {
	return xxx_messageInfo_MetadataSignature.Size(m)
}
func (m *MetadataSignature) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataSignature.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataSignature proto.InternalMessageInfo

func (m *MetadataSignature) GetSignatureHeader() []byte {
	if m != nil {
		return m.SignatureHeader
	}
	return nil
}

func (m *MetadataSignature) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type Header struct {
	ChannelHeader        []byte   `protobuf:"bytes,1,opt,name=channel_header,json=channelHeader,proto3" json:"channel_header,omitempty"`
	SignatureHeader      []byte   `protobuf:"bytes,2,opt,name=signature_header,json=signatureHeader,proto3" json:"signature_header,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{3}
}
func (m *Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Header.Unmarshal(m, b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Header.Marshal(b, m, deterministic)
}
func (dst *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(dst, src)
}
func (m *Header) XXX_Size() int {
	return xxx_messageInfo_Header.Size(m)
}
func (m *Header) XXX_DiscardUnknown() {
	xxx_messageInfo_Header.DiscardUnknown(m)
}

var xxx_messageInfo_Header proto.InternalMessageInfo

func (m *Header) GetChannelHeader() []byte {
	if m != nil {
		return m.ChannelHeader
	}
	return nil
}

func (m *Header) GetSignatureHeader() []byte {
	if m != nil {
		return m.SignatureHeader
	}
	return nil
}


type ChannelHeader struct {
	Type int32 `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	
	Version int32 `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	
	
	Timestamp *timestamp.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	
	ChannelId string `protobuf:"bytes,4,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	
	
	
	
	
	
	TxId string `protobuf:"bytes,5,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	
	
	
	
	
	
	
	Epoch uint64 `protobuf:"varint,6,opt,name=epoch,proto3" json:"epoch,omitempty"`
	
	Extension []byte `protobuf:"bytes,7,opt,name=extension,proto3" json:"extension,omitempty"`
	
	
	TlsCertHash          []byte   `protobuf:"bytes,8,opt,name=tls_cert_hash,json=tlsCertHash,proto3" json:"tls_cert_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChannelHeader) Reset()         { *m = ChannelHeader{} }
func (m *ChannelHeader) String() string { return proto.CompactTextString(m) }
func (*ChannelHeader) ProtoMessage()    {}
func (*ChannelHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{4}
}
func (m *ChannelHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChannelHeader.Unmarshal(m, b)
}
func (m *ChannelHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChannelHeader.Marshal(b, m, deterministic)
}
func (dst *ChannelHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChannelHeader.Merge(dst, src)
}
func (m *ChannelHeader) XXX_Size() int {
	return xxx_messageInfo_ChannelHeader.Size(m)
}
func (m *ChannelHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_ChannelHeader.DiscardUnknown(m)
}

var xxx_messageInfo_ChannelHeader proto.InternalMessageInfo

func (m *ChannelHeader) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *ChannelHeader) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ChannelHeader) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *ChannelHeader) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *ChannelHeader) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *ChannelHeader) GetEpoch() uint64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func (m *ChannelHeader) GetExtension() []byte {
	if m != nil {
		return m.Extension
	}
	return nil
}

func (m *ChannelHeader) GetTlsCertHash() []byte {
	if m != nil {
		return m.TlsCertHash
	}
	return nil
}

type SignatureHeader struct {
	
	Creator []byte `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	
	Nonce                []byte   `protobuf:"bytes,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignatureHeader) Reset()         { *m = SignatureHeader{} }
func (m *SignatureHeader) String() string { return proto.CompactTextString(m) }
func (*SignatureHeader) ProtoMessage()    {}
func (*SignatureHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{5}
}
func (m *SignatureHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignatureHeader.Unmarshal(m, b)
}
func (m *SignatureHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignatureHeader.Marshal(b, m, deterministic)
}
func (dst *SignatureHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignatureHeader.Merge(dst, src)
}
func (m *SignatureHeader) XXX_Size() int {
	return xxx_messageInfo_SignatureHeader.Size(m)
}
func (m *SignatureHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_SignatureHeader.DiscardUnknown(m)
}

var xxx_messageInfo_SignatureHeader proto.InternalMessageInfo

func (m *SignatureHeader) GetCreator() []byte {
	if m != nil {
		return m.Creator
	}
	return nil
}

func (m *SignatureHeader) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}


type Payload struct {
	
	Header *Header `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Payload) Reset()         { *m = Payload{} }
func (m *Payload) String() string { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()    {}
func (*Payload) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{6}
}
func (m *Payload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Payload.Unmarshal(m, b)
}
func (m *Payload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Payload.Marshal(b, m, deterministic)
}
func (dst *Payload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Payload.Merge(dst, src)
}
func (m *Payload) XXX_Size() int {
	return xxx_messageInfo_Payload.Size(m)
}
func (m *Payload) XXX_DiscardUnknown() {
	xxx_messageInfo_Payload.DiscardUnknown(m)
}

var xxx_messageInfo_Payload proto.InternalMessageInfo

func (m *Payload) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Payload) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}


type Envelope struct {
	
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Envelope) Reset()         { *m = Envelope{} }
func (m *Envelope) String() string { return proto.CompactTextString(m) }
func (*Envelope) ProtoMessage()    {}
func (*Envelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{7}
}
func (m *Envelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Envelope.Unmarshal(m, b)
}
func (m *Envelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Envelope.Marshal(b, m, deterministic)
}
func (dst *Envelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Envelope.Merge(dst, src)
}
func (m *Envelope) XXX_Size() int {
	return xxx_messageInfo_Envelope.Size(m)
}
func (m *Envelope) XXX_DiscardUnknown() {
	xxx_messageInfo_Envelope.DiscardUnknown(m)
}

var xxx_messageInfo_Envelope proto.InternalMessageInfo

func (m *Envelope) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Envelope) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}





type Block struct {
	Header               *BlockHeader   `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Data                 *BlockData     `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Metadata             *BlockMetadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{8}
}
func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (dst *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(dst, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetHeader() *BlockHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Block) GetData() *BlockData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Block) GetMetadata() *BlockMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}




type BlockHeader struct {
	Number               uint64   `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	PreviousHash         []byte   `protobuf:"bytes,2,opt,name=previous_hash,json=previousHash,proto3" json:"previous_hash,omitempty"`
	DataHash             []byte   `protobuf:"bytes,3,opt,name=data_hash,json=dataHash,proto3" json:"data_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockHeader) Reset()         { *m = BlockHeader{} }
func (m *BlockHeader) String() string { return proto.CompactTextString(m) }
func (*BlockHeader) ProtoMessage()    {}
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{9}
}
func (m *BlockHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockHeader.Unmarshal(m, b)
}
func (m *BlockHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockHeader.Marshal(b, m, deterministic)
}
func (dst *BlockHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHeader.Merge(dst, src)
}
func (m *BlockHeader) XXX_Size() int {
	return xxx_messageInfo_BlockHeader.Size(m)
}
func (m *BlockHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHeader.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHeader proto.InternalMessageInfo

func (m *BlockHeader) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *BlockHeader) GetPreviousHash() []byte {
	if m != nil {
		return m.PreviousHash
	}
	return nil
}

func (m *BlockHeader) GetDataHash() []byte {
	if m != nil {
		return m.DataHash
	}
	return nil
}

type BlockData struct {
	Data                 [][]byte `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockData) Reset()         { *m = BlockData{} }
func (m *BlockData) String() string { return proto.CompactTextString(m) }
func (*BlockData) ProtoMessage()    {}
func (*BlockData) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{10}
}
func (m *BlockData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockData.Unmarshal(m, b)
}
func (m *BlockData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockData.Marshal(b, m, deterministic)
}
func (dst *BlockData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockData.Merge(dst, src)
}
func (m *BlockData) XXX_Size() int {
	return xxx_messageInfo_BlockData.Size(m)
}
func (m *BlockData) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockData.DiscardUnknown(m)
}

var xxx_messageInfo_BlockData proto.InternalMessageInfo

func (m *BlockData) GetData() [][]byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type BlockMetadata struct {
	Metadata             [][]byte `protobuf:"bytes,1,rep,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockMetadata) Reset()         { *m = BlockMetadata{} }
func (m *BlockMetadata) String() string { return proto.CompactTextString(m) }
func (*BlockMetadata) ProtoMessage()    {}
func (*BlockMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_61611b181d9b3dc1, []int{11}
}
func (m *BlockMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockMetadata.Unmarshal(m, b)
}
func (m *BlockMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockMetadata.Marshal(b, m, deterministic)
}
func (dst *BlockMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockMetadata.Merge(dst, src)
}
func (m *BlockMetadata) XXX_Size() int {
	return xxx_messageInfo_BlockMetadata.Size(m)
}
func (m *BlockMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_BlockMetadata proto.InternalMessageInfo

func (m *BlockMetadata) GetMetadata() [][]byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func init() {
	proto.RegisterType((*LastConfig)(nil), "common.LastConfig")
	proto.RegisterType((*Metadata)(nil), "common.Metadata")
	proto.RegisterType((*MetadataSignature)(nil), "common.MetadataSignature")
	proto.RegisterType((*Header)(nil), "common.Header")
	proto.RegisterType((*ChannelHeader)(nil), "common.ChannelHeader")
	proto.RegisterType((*SignatureHeader)(nil), "common.SignatureHeader")
	proto.RegisterType((*Payload)(nil), "common.Payload")
	proto.RegisterType((*Envelope)(nil), "common.Envelope")
	proto.RegisterType((*Block)(nil), "common.Block")
	proto.RegisterType((*BlockHeader)(nil), "common.BlockHeader")
	proto.RegisterType((*BlockData)(nil), "common.BlockData")
	proto.RegisterType((*BlockMetadata)(nil), "common.BlockMetadata")
	proto.RegisterEnum("common.Status", Status_name, Status_value)
	proto.RegisterEnum("common.HeaderType", HeaderType_name, HeaderType_value)
	proto.RegisterEnum("common.BlockMetadataIndex", BlockMetadataIndex_name, BlockMetadataIndex_value)
}

func init() { proto.RegisterFile("common/common.proto", fileDescriptor_common_61611b181d9b3dc1) }

var fileDescriptor_common_61611b181d9b3dc1 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x55, 0x4f, 0x6f, 0xeb, 0xc4,
	0x17, 0x6d, 0xe2, 0xfc, 0xbd, 0x69, 0x5a, 0x77, 0xd2, 0xfe, 0x9e, 0x7f, 0x85, 0xc7, 0xab, 0x0c,
	0x0f, 0x95, 0x56, 0x4a, 0x45, 0xd9, 0xc0, 0xd2, 0xb1, 0xa7, 0xad, 0xd5, 0x74, 0x1c, 0xc6, 0xce,
	0x43, 0xbc, 0x87, 0x64, 0xb9, 0xc9, 0x34, 0x89, 0x48, 0xec, 0xc8, 0x9e, 0x54, 0xed, 0x9a, 0x3d,
	0x42, 0x82, 0x2d, 0xdf, 0x85, 0x25, 0x9f, 0x85, 0x35, 0x88, 0x2d, 0x1a, 0x8f, 0xed, 0x97, 0x94,
	0x4a, 0xac, 0xe2, 0x73, 0xe6, 0xf8, 0xde, 0x33, 0xf7, 0xcc, 0xc4, 0xd0, 0x19, 0x45, 0x8b, 0x45,
	0x14, 0x9e, 0xc9, 0x9f, 0xee, 0x32, 0x8e, 0x78, 0x84, 0x6a, 0x12, 0x1d, 0xbe, 0x9a, 0x44, 0xd1,
	0x64, 0xce, 0xce, 0x52, 0xf6, 0x76, 0x75, 0x77, 0xc6, 0x67, 0x0b, 0x96, 0xf0, 0x60, 0xb1, 0x94,
	0x42, 0x5d, 0x07, 0xe8, 0x07, 0x09, 0x37, 0xa3, 0xf0, 0x6e, 0x36, 0x41, 0xfb, 0x50, 0x9d, 0x85,
	0x63, 0xf6, 0xa0, 0x95, 0x8e, 0x4a, 0xc7, 0x15, 0x2a, 0x81, 0xfe, 0x0e, 0x1a, 0x37, 0x8c, 0x07,
	0xe3, 0x80, 0x07, 0x42, 0x71, 0x1f, 0xcc, 0x57, 0x2c, 0x55, 0x6c, 0x53, 0x09, 0xd0, 0x57, 0x00,
	0xc9, 0x6c, 0x12, 0x06, 0x7c, 0x15, 0xb3, 0x44, 0x2b, 0x1f, 0x29, 0xc7, 0xad, 0xf3, 0xff, 0x77,
	0x33, 0x47, 0xf9, 0xbb, 0x6e, 0xae, 0xa0, 0x6b, 0x62, 0xfd, 0x3b, 0xd8, 0xfb, 0x97, 0x00, 0x7d,
	0x06, 0x6a, 0x21, 0xf1, 0xa7, 0x2c, 0x18, 0xb3, 0x38, 0x6b, 0xb8, 0x5b, 0xf0, 0x57, 0x29, 0x8d,
	0x3e, 0x84, 0x66, 0x41, 0x69, 0xe5, 0x54, 0xf3, 0x9e, 0xd0, 0xdf, 0x42, 0x2d, 0xd3, 0xbd, 0x86,
	0x9d, 0xd1, 0x34, 0x08, 0x43, 0x36, 0xdf, 0x2c, 0xd8, 0xce, 0xd8, 0x4c, 0xf6, 0x5c, 0xe7, 0xf2,
	0xb3, 0x9d, 0xf5, 0x1f, 0xca, 0xd0, 0x36, 0x37, 0x5e, 0x46, 0x50, 0xe1, 0x8f, 0x4b, 0x39, 0x9b,
	0x2a, 0x4d, 0x9f, 0x91, 0x06, 0xf5, 0x7b, 0x16, 0x27, 0xb3, 0x28, 0x4c, 0xeb, 0x54, 0x69, 0x0e,
	0xd1, 0x97, 0xd0, 0x2c, 0xd2, 0xd0, 0x94, 0xa3, 0xd2, 0x71, 0xeb, 0xfc, 0xb0, 0x2b, 0xf3, 0xea,
	0xe6, 0x79, 0x75, 0xbd, 0x5c, 0x41, 0xdf, 0x8b, 0xd1, 0x4b, 0x80, 0x7c, 0x2f, 0xb3, 0xb1, 0x56,
	0x39, 0x2a, 0x1d, 0x37, 0x69, 0x33, 0x63, 0xec, 0x31, 0xea, 0x40, 0x95, 0x3f, 0x88, 0x95, 0x6a,
	0xba, 0x52, 0xe1, 0x0f, 0xf6, 0x58, 0x04, 0xc7, 0x96, 0xd1, 0x68, 0xaa, 0xd5, 0x64, 0xb4, 0x29,
	0x10, 0xd3, 0x63, 0x0f, 0x9c, 0x85, 0xa9, 0xbf, 0xba, 0x9c, 0x5e, 0x41, 0x20, 0x1d, 0xda, 0x7c,
	0x9e, 0xf8, 0x23, 0x16, 0x73, 0x7f, 0x1a, 0x24, 0x53, 0xad, 0x91, 0x2a, 0x5a, 0x7c, 0x9e, 0x98,
	0x2c, 0xe6, 0x57, 0x41, 0x32, 0xd5, 0x0d, 0xd8, 0x75, 0x9f, 0x44, 0xa2, 0x41, 0x7d, 0x14, 0xb3,
	0x80, 0x47, 0xf9, 0x8c, 0x73, 0x28, 0x4c, 0x84, 0x51, 0x38, 0xca, 0x83, 0x92, 0x40, 0xc7, 0x50,
	0x1f, 0x04, 0x8f, 0xf3, 0x28, 0x18, 0xa3, 0x4f, 0xa1, 0xb6, 0x96, 0x4e, 0xeb, 0x7c, 0x27, 0x3f,
	0x44, 0xb2, 0x34, 0xcd, 0x56, 0xc5, 0xa4, 0xc5, 0x89, 0xc9, 0xea, 0xa4, 0xcf, 0x7a, 0x0f, 0x1a,
	0x38, 0xbc, 0x67, 0xf3, 0x48, 0x4e, 0x7d, 0x29, 0x4b, 0xe6, 0x16, 0x32, 0xf8, 0x1f, 0xe7, 0xe5,
	0xc7, 0x12, 0x54, 0x7b, 0xf3, 0x68, 0xf4, 0x3d, 0x3a, 0x7d, 0xe2, 0xa4, 0x93, 0x3b, 0x49, 0x97,
	0x9f, 0xd8, 0x79, 0xbd, 0x66, 0xa7, 0x75, 0xbe, 0xb7, 0x21, 0xb5, 0x02, 0x1e, 0x48, 0x87, 0xe8,
	0x73, 0x68, 0x2c, 0xb2, 0xb3, 0x9e, 0x05, 0x7e, 0xb0, 0x21, 0xcd, 0x2f, 0x02, 0x2d, 0x64, 0xfa,
	0x04, 0x5a, 0x6b, 0x0d, 0xd1, 0xff, 0xa0, 0x16, 0xae, 0x16, 0xb7, 0x99, 0xab, 0x0a, 0xcd, 0x10,
	0xfa, 0x18, 0xda, 0xcb, 0x98, 0xdd, 0xcf, 0xa2, 0x55, 0x22, 0x93, 0x92, 0x3b, 0xdb, 0xce, 0x49,
	0x11, 0x15, 0xfa, 0x00, 0x9a, 0xa2, 0xa6, 0x14, 0x28, 0xa9, 0xa0, 0x21, 0x88, 0x34, 0xc7, 0x57,
	0xd0, 0x2c, 0xec, 0x16, 0xe3, 0x2d, 0x1d, 0x29, 0xc5, 0x78, 0x4f, 0xa1, 0xbd, 0x61, 0x12, 0x1d,
	0xae, 0xed, 0x46, 0x0a, 0x0b, 0x7c, 0xf2, 0x5b, 0x09, 0x6a, 0x2e, 0x0f, 0xf8, 0x2a, 0x41, 0x2d,
	0xa8, 0x0f, 0xc9, 0x35, 0x71, 0xbe, 0x21, 0xea, 0x16, 0xda, 0x86, 0xba, 0x3b, 0x34, 0x4d, 0xec,
	0xba, 0xea, 0xef, 0x25, 0xa4, 0x42, 0xab, 0x67, 0x58, 0x3e, 0xc5, 0x5f, 0x0f, 0xb1, 0xeb, 0xa9,
	0x3f, 0x29, 0x68, 0x07, 0x9a, 0x17, 0x0e, 0xed, 0xd9, 0x96, 0x85, 0x89, 0xfa, 0x73, 0x8a, 0x89,
	0xe3, 0xf9, 0x17, 0xce, 0x90, 0x58, 0xea, 0x2f, 0x0a, 0x7a, 0x09, 0x5a, 0xa6, 0xf6, 0x31, 0xf1,
	0x6c, 0xef, 0x5b, 0xdf, 0x73, 0x1c, 0xbf, 0x6f, 0xd0, 0x4b, 0xac, 0xfe, 0xaa, 0xa0, 0x43, 0x38,
	0xb0, 0x89, 0x87, 0x29, 0x31, 0xfa, 0xbe, 0x8b, 0xe9, 0x1b, 0x4c, 0x7d, 0x4c, 0xa9, 0x43, 0xd5,
	0x3f, 0x15, 0xb4, 0x0f, 0xbb, 0xa2, 0x94, 0x7d, 0x33, 0xe8, 0xe3, 0x1b, 0x4c, 0x3c, 0x6c, 0xa9,
	0x7f, 0x29, 0x48, 0x83, 0x8e, 0x10, 0xda, 0x26, 0xf6, 0x87, 0xc4, 0x78, 0x63, 0xd8, 0x7d, 0xa3,
	0xd7, 0xc7, 0xea, 0xdf, 0xca, 0xc9, 0x1f, 0x25, 0x00, 0x39, 0x75, 0x4f, 0xdc, 0xe3, 0x16, 0xd4,
	0x6f, 0xb0, 0xeb, 0x1a, 0x97, 0x58, 0xdd, 0x42, 0x00, 0x35, 0xd3, 0x21, 0x17, 0xf6, 0xa5, 0x5a,
	0x42, 0x7b, 0xd0, 0x96, 0xcf, 0xfe, 0x70, 0x60, 0x19, 0x1e, 0x56, 0xcb, 0x48, 0x83, 0x7d, 0x4c,
	0x2c, 0x87, 0xba, 0x98, 0xfa, 0x1e, 0x35, 0x88, 0x6b, 0x98, 0x9e, 0xed, 0x10, 0x55, 0x41, 0x2f,
	0xa0, 0xe3, 0x50, 0x0b, 0xd3, 0x27, 0x0b, 0x15, 0x74, 0x00, 0x7b, 0x16, 0xee, 0xdb, 0xc2, 0xb1,
	0x8b, 0xf1, 0xb5, 0x6f, 0x93, 0x0b, 0x47, 0xad, 0x0a, 0xda, 0xbc, 0x32, 0x6c, 0x62, 0x3a, 0x16,
	0xf6, 0x07, 0x86, 0x79, 0x2d, 0xfa, 0xd7, 0x44, 0x83, 0x01, 0xc6, 0xd4, 0x37, 0xac, 0x1b, 0x9b,
	0xf8, 0xce, 0x00, 0x53, 0x23, 0xad, 0xd3, 0x10, 0x2f, 0x78, 0xce, 0x35, 0x26, 0x1b, 0xe5, 0x9b,
	0xe8, 0x23, 0x38, 0x94, 0xf4, 0xb3, 0xbe, 0xe0, 0xe4, 0x1d, 0xa0, 0x8d, 0x70, 0x6d, 0xf1, 0xc7,
	0x8f, 0x76, 0x00, 0x5c, 0xfb, 0x92, 0x18, 0xde, 0x90, 0x62, 0x57, 0xdd, 0x42, 0xbb, 0xd0, 0xea,
	0x1b, 0xae, 0xe7, 0x17, 0x7b, 0x7f, 0x01, 0x9d, 0xb5, 0x3a, 0xae, 0x7f, 0x61, 0xf7, 0x3d, 0x4c,
	0xd5, 0xb2, 0x98, 0x56, 0xb6, 0x4f, 0x55, 0xe9, 0xb9, 0xf0, 0x49, 0x14, 0x4f, 0xba, 0xd3, 0xc7,
	0x25, 0x8b, 0xe7, 0x6c, 0x3c, 0x61, 0x71, 0xf7, 0x2e, 0xb8, 0x8d, 0x67, 0x23, 0xf9, 0x37, 0x97,
	0x64, 0x77, 0xe0, 0xed, 0xe9, 0x64, 0xc6, 0xa7, 0xab, 0x5b, 0x01, 0xcf, 0xd6, 0xc4, 0x67, 0x52,
	0x2c, 0xbf, 0x61, 0x49, 0xf6, 0x9d, 0xbb, 0xad, 0xa5, 0xf0, 0x8b, 0x7f, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x99, 0x03, 0x31, 0xfb, 0xff, 0x06, 0x00, 0x00,
}
