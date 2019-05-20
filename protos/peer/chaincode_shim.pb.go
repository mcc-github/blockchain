


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type ChaincodeMessage_Type int32

const (
	ChaincodeMessage_UNDEFINED             ChaincodeMessage_Type = 0
	ChaincodeMessage_REGISTER              ChaincodeMessage_Type = 1
	ChaincodeMessage_REGISTERED            ChaincodeMessage_Type = 2
	ChaincodeMessage_INIT                  ChaincodeMessage_Type = 3
	ChaincodeMessage_READY                 ChaincodeMessage_Type = 4
	ChaincodeMessage_TRANSACTION           ChaincodeMessage_Type = 5
	ChaincodeMessage_COMPLETED             ChaincodeMessage_Type = 6
	ChaincodeMessage_ERROR                 ChaincodeMessage_Type = 7
	ChaincodeMessage_GET_STATE             ChaincodeMessage_Type = 8
	ChaincodeMessage_PUT_STATE             ChaincodeMessage_Type = 9
	ChaincodeMessage_DEL_STATE             ChaincodeMessage_Type = 10
	ChaincodeMessage_INVOKE_CHAINCODE      ChaincodeMessage_Type = 11
	ChaincodeMessage_RESPONSE              ChaincodeMessage_Type = 13
	ChaincodeMessage_GET_STATE_BY_RANGE    ChaincodeMessage_Type = 14
	ChaincodeMessage_GET_QUERY_RESULT      ChaincodeMessage_Type = 15
	ChaincodeMessage_QUERY_STATE_NEXT      ChaincodeMessage_Type = 16
	ChaincodeMessage_QUERY_STATE_CLOSE     ChaincodeMessage_Type = 17
	ChaincodeMessage_KEEPALIVE             ChaincodeMessage_Type = 18
	ChaincodeMessage_GET_HISTORY_FOR_KEY   ChaincodeMessage_Type = 19
	ChaincodeMessage_GET_STATE_METADATA    ChaincodeMessage_Type = 20
	ChaincodeMessage_PUT_STATE_METADATA    ChaincodeMessage_Type = 21
	ChaincodeMessage_GET_PRIVATE_DATA_HASH ChaincodeMessage_Type = 22
)

var ChaincodeMessage_Type_name = map[int32]string{
	0:  "UNDEFINED",
	1:  "REGISTER",
	2:  "REGISTERED",
	3:  "INIT",
	4:  "READY",
	5:  "TRANSACTION",
	6:  "COMPLETED",
	7:  "ERROR",
	8:  "GET_STATE",
	9:  "PUT_STATE",
	10: "DEL_STATE",
	11: "INVOKE_CHAINCODE",
	13: "RESPONSE",
	14: "GET_STATE_BY_RANGE",
	15: "GET_QUERY_RESULT",
	16: "QUERY_STATE_NEXT",
	17: "QUERY_STATE_CLOSE",
	18: "KEEPALIVE",
	19: "GET_HISTORY_FOR_KEY",
	20: "GET_STATE_METADATA",
	21: "PUT_STATE_METADATA",
	22: "GET_PRIVATE_DATA_HASH",
}
var ChaincodeMessage_Type_value = map[string]int32{
	"UNDEFINED":             0,
	"REGISTER":              1,
	"REGISTERED":            2,
	"INIT":                  3,
	"READY":                 4,
	"TRANSACTION":           5,
	"COMPLETED":             6,
	"ERROR":                 7,
	"GET_STATE":             8,
	"PUT_STATE":             9,
	"DEL_STATE":             10,
	"INVOKE_CHAINCODE":      11,
	"RESPONSE":              13,
	"GET_STATE_BY_RANGE":    14,
	"GET_QUERY_RESULT":      15,
	"QUERY_STATE_NEXT":      16,
	"QUERY_STATE_CLOSE":     17,
	"KEEPALIVE":             18,
	"GET_HISTORY_FOR_KEY":   19,
	"GET_STATE_METADATA":    20,
	"PUT_STATE_METADATA":    21,
	"GET_PRIVATE_DATA_HASH": 22,
}

func (x ChaincodeMessage_Type) String() string {
	return proto.EnumName(ChaincodeMessage_Type_name, int32(x))
}
func (ChaincodeMessage_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{0, 0}
}

type ChaincodeMessage struct {
	Type      ChaincodeMessage_Type `protobuf:"varint,1,opt,name=type,proto3,enum=protos.ChaincodeMessage_Type" json:"type,omitempty"`
	Timestamp *timestamp.Timestamp  `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Payload   []byte                `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	Txid      string                `protobuf:"bytes,4,opt,name=txid,proto3" json:"txid,omitempty"`
	Proposal  *SignedProposal       `protobuf:"bytes,5,opt,name=proposal,proto3" json:"proposal,omitempty"`
	
	
	
	ChaincodeEvent *ChaincodeEvent `protobuf:"bytes,6,opt,name=chaincode_event,json=chaincodeEvent,proto3" json:"chaincode_event,omitempty"`
	
	ChannelId            string   `protobuf:"bytes,7,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeMessage) Reset()         { *m = ChaincodeMessage{} }
func (m *ChaincodeMessage) String() string { return proto.CompactTextString(m) }
func (*ChaincodeMessage) ProtoMessage()    {}
func (*ChaincodeMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{0}
}
func (m *ChaincodeMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeMessage.Unmarshal(m, b)
}
func (m *ChaincodeMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeMessage.Marshal(b, m, deterministic)
}
func (dst *ChaincodeMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeMessage.Merge(dst, src)
}
func (m *ChaincodeMessage) XXX_Size() int {
	return xxx_messageInfo_ChaincodeMessage.Size(m)
}
func (m *ChaincodeMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeMessage proto.InternalMessageInfo

func (m *ChaincodeMessage) GetType() ChaincodeMessage_Type {
	if m != nil {
		return m.Type
	}
	return ChaincodeMessage_UNDEFINED
}

func (m *ChaincodeMessage) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *ChaincodeMessage) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *ChaincodeMessage) GetTxid() string {
	if m != nil {
		return m.Txid
	}
	return ""
}

func (m *ChaincodeMessage) GetProposal() *SignedProposal {
	if m != nil {
		return m.Proposal
	}
	return nil
}

func (m *ChaincodeMessage) GetChaincodeEvent() *ChaincodeEvent {
	if m != nil {
		return m.ChaincodeEvent
	}
	return nil
}

func (m *ChaincodeMessage) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}




type GetState struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Collection           string   `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetState) Reset()         { *m = GetState{} }
func (m *GetState) String() string { return proto.CompactTextString(m) }
func (*GetState) ProtoMessage()    {}
func (*GetState) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{1}
}
func (m *GetState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetState.Unmarshal(m, b)
}
func (m *GetState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetState.Marshal(b, m, deterministic)
}
func (dst *GetState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetState.Merge(dst, src)
}
func (m *GetState) XXX_Size() int {
	return xxx_messageInfo_GetState.Size(m)
}
func (m *GetState) XXX_DiscardUnknown() {
	xxx_messageInfo_GetState.DiscardUnknown(m)
}

var xxx_messageInfo_GetState proto.InternalMessageInfo

func (m *GetState) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *GetState) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

type GetStateMetadata struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Collection           string   `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetStateMetadata) Reset()         { *m = GetStateMetadata{} }
func (m *GetStateMetadata) String() string { return proto.CompactTextString(m) }
func (*GetStateMetadata) ProtoMessage()    {}
func (*GetStateMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{2}
}
func (m *GetStateMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetStateMetadata.Unmarshal(m, b)
}
func (m *GetStateMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetStateMetadata.Marshal(b, m, deterministic)
}
func (dst *GetStateMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetStateMetadata.Merge(dst, src)
}
func (m *GetStateMetadata) XXX_Size() int {
	return xxx_messageInfo_GetStateMetadata.Size(m)
}
func (m *GetStateMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_GetStateMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_GetStateMetadata proto.InternalMessageInfo

func (m *GetStateMetadata) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *GetStateMetadata) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}





type PutState struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Collection           string   `protobuf:"bytes,3,opt,name=collection,proto3" json:"collection,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PutState) Reset()         { *m = PutState{} }
func (m *PutState) String() string { return proto.CompactTextString(m) }
func (*PutState) ProtoMessage()    {}
func (*PutState) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{3}
}
func (m *PutState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PutState.Unmarshal(m, b)
}
func (m *PutState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PutState.Marshal(b, m, deterministic)
}
func (dst *PutState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PutState.Merge(dst, src)
}
func (m *PutState) XXX_Size() int {
	return xxx_messageInfo_PutState.Size(m)
}
func (m *PutState) XXX_DiscardUnknown() {
	xxx_messageInfo_PutState.DiscardUnknown(m)
}

var xxx_messageInfo_PutState proto.InternalMessageInfo

func (m *PutState) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *PutState) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *PutState) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

type PutStateMetadata struct {
	Key                  string         `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Collection           string         `protobuf:"bytes,3,opt,name=collection,proto3" json:"collection,omitempty"`
	Metadata             *StateMetadata `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *PutStateMetadata) Reset()         { *m = PutStateMetadata{} }
func (m *PutStateMetadata) String() string { return proto.CompactTextString(m) }
func (*PutStateMetadata) ProtoMessage()    {}
func (*PutStateMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{4}
}
func (m *PutStateMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PutStateMetadata.Unmarshal(m, b)
}
func (m *PutStateMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PutStateMetadata.Marshal(b, m, deterministic)
}
func (dst *PutStateMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PutStateMetadata.Merge(dst, src)
}
func (m *PutStateMetadata) XXX_Size() int {
	return xxx_messageInfo_PutStateMetadata.Size(m)
}
func (m *PutStateMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_PutStateMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_PutStateMetadata proto.InternalMessageInfo

func (m *PutStateMetadata) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *PutStateMetadata) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

func (m *PutStateMetadata) GetMetadata() *StateMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}





type DelState struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Collection           string   `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DelState) Reset()         { *m = DelState{} }
func (m *DelState) String() string { return proto.CompactTextString(m) }
func (*DelState) ProtoMessage()    {}
func (*DelState) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{5}
}
func (m *DelState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DelState.Unmarshal(m, b)
}
func (m *DelState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DelState.Marshal(b, m, deterministic)
}
func (dst *DelState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DelState.Merge(dst, src)
}
func (m *DelState) XXX_Size() int {
	return xxx_messageInfo_DelState.Size(m)
}
func (m *DelState) XXX_DiscardUnknown() {
	xxx_messageInfo_DelState.DiscardUnknown(m)
}

var xxx_messageInfo_DelState proto.InternalMessageInfo

func (m *DelState) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *DelState) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}





type GetStateByRange struct {
	StartKey             string   `protobuf:"bytes,1,opt,name=startKey,proto3" json:"startKey,omitempty"`
	EndKey               string   `protobuf:"bytes,2,opt,name=endKey,proto3" json:"endKey,omitempty"`
	Collection           string   `protobuf:"bytes,3,opt,name=collection,proto3" json:"collection,omitempty"`
	Metadata             []byte   `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetStateByRange) Reset()         { *m = GetStateByRange{} }
func (m *GetStateByRange) String() string { return proto.CompactTextString(m) }
func (*GetStateByRange) ProtoMessage()    {}
func (*GetStateByRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{6}
}
func (m *GetStateByRange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetStateByRange.Unmarshal(m, b)
}
func (m *GetStateByRange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetStateByRange.Marshal(b, m, deterministic)
}
func (dst *GetStateByRange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetStateByRange.Merge(dst, src)
}
func (m *GetStateByRange) XXX_Size() int {
	return xxx_messageInfo_GetStateByRange.Size(m)
}
func (m *GetStateByRange) XXX_DiscardUnknown() {
	xxx_messageInfo_GetStateByRange.DiscardUnknown(m)
}

var xxx_messageInfo_GetStateByRange proto.InternalMessageInfo

func (m *GetStateByRange) GetStartKey() string {
	if m != nil {
		return m.StartKey
	}
	return ""
}

func (m *GetStateByRange) GetEndKey() string {
	if m != nil {
		return m.EndKey
	}
	return ""
}

func (m *GetStateByRange) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

func (m *GetStateByRange) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}





type GetQueryResult struct {
	Query                string   `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	Collection           string   `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	Metadata             []byte   `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetQueryResult) Reset()         { *m = GetQueryResult{} }
func (m *GetQueryResult) String() string { return proto.CompactTextString(m) }
func (*GetQueryResult) ProtoMessage()    {}
func (*GetQueryResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{7}
}
func (m *GetQueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetQueryResult.Unmarshal(m, b)
}
func (m *GetQueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetQueryResult.Marshal(b, m, deterministic)
}
func (dst *GetQueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetQueryResult.Merge(dst, src)
}
func (m *GetQueryResult) XXX_Size() int {
	return xxx_messageInfo_GetQueryResult.Size(m)
}
func (m *GetQueryResult) XXX_DiscardUnknown() {
	xxx_messageInfo_GetQueryResult.DiscardUnknown(m)
}

var xxx_messageInfo_GetQueryResult proto.InternalMessageInfo

func (m *GetQueryResult) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

func (m *GetQueryResult) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

func (m *GetQueryResult) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}




type QueryMetadata struct {
	PageSize             int32    `protobuf:"varint,1,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
	Bookmark             string   `protobuf:"bytes,2,opt,name=bookmark,proto3" json:"bookmark,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryMetadata) Reset()         { *m = QueryMetadata{} }
func (m *QueryMetadata) String() string { return proto.CompactTextString(m) }
func (*QueryMetadata) ProtoMessage()    {}
func (*QueryMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{8}
}
func (m *QueryMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryMetadata.Unmarshal(m, b)
}
func (m *QueryMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryMetadata.Marshal(b, m, deterministic)
}
func (dst *QueryMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMetadata.Merge(dst, src)
}
func (m *QueryMetadata) XXX_Size() int {
	return xxx_messageInfo_QueryMetadata.Size(m)
}
func (m *QueryMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMetadata proto.InternalMessageInfo

func (m *QueryMetadata) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *QueryMetadata) GetBookmark() string {
	if m != nil {
		return m.Bookmark
	}
	return ""
}



type GetHistoryForKey struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetHistoryForKey) Reset()         { *m = GetHistoryForKey{} }
func (m *GetHistoryForKey) String() string { return proto.CompactTextString(m) }
func (*GetHistoryForKey) ProtoMessage()    {}
func (*GetHistoryForKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{9}
}
func (m *GetHistoryForKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetHistoryForKey.Unmarshal(m, b)
}
func (m *GetHistoryForKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetHistoryForKey.Marshal(b, m, deterministic)
}
func (dst *GetHistoryForKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetHistoryForKey.Merge(dst, src)
}
func (m *GetHistoryForKey) XXX_Size() int {
	return xxx_messageInfo_GetHistoryForKey.Size(m)
}
func (m *GetHistoryForKey) XXX_DiscardUnknown() {
	xxx_messageInfo_GetHistoryForKey.DiscardUnknown(m)
}

var xxx_messageInfo_GetHistoryForKey proto.InternalMessageInfo

func (m *GetHistoryForKey) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type QueryStateNext struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryStateNext) Reset()         { *m = QueryStateNext{} }
func (m *QueryStateNext) String() string { return proto.CompactTextString(m) }
func (*QueryStateNext) ProtoMessage()    {}
func (*QueryStateNext) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{10}
}
func (m *QueryStateNext) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryStateNext.Unmarshal(m, b)
}
func (m *QueryStateNext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryStateNext.Marshal(b, m, deterministic)
}
func (dst *QueryStateNext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryStateNext.Merge(dst, src)
}
func (m *QueryStateNext) XXX_Size() int {
	return xxx_messageInfo_QueryStateNext.Size(m)
}
func (m *QueryStateNext) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryStateNext.DiscardUnknown(m)
}

var xxx_messageInfo_QueryStateNext proto.InternalMessageInfo

func (m *QueryStateNext) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type QueryStateClose struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryStateClose) Reset()         { *m = QueryStateClose{} }
func (m *QueryStateClose) String() string { return proto.CompactTextString(m) }
func (*QueryStateClose) ProtoMessage()    {}
func (*QueryStateClose) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{11}
}
func (m *QueryStateClose) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryStateClose.Unmarshal(m, b)
}
func (m *QueryStateClose) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryStateClose.Marshal(b, m, deterministic)
}
func (dst *QueryStateClose) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryStateClose.Merge(dst, src)
}
func (m *QueryStateClose) XXX_Size() int {
	return xxx_messageInfo_QueryStateClose.Size(m)
}
func (m *QueryStateClose) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryStateClose.DiscardUnknown(m)
}

var xxx_messageInfo_QueryStateClose proto.InternalMessageInfo

func (m *QueryStateClose) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}


type QueryResultBytes struct {
	ResultBytes          []byte   `protobuf:"bytes,1,opt,name=resultBytes,proto3" json:"resultBytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryResultBytes) Reset()         { *m = QueryResultBytes{} }
func (m *QueryResultBytes) String() string { return proto.CompactTextString(m) }
func (*QueryResultBytes) ProtoMessage()    {}
func (*QueryResultBytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{12}
}
func (m *QueryResultBytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResultBytes.Unmarshal(m, b)
}
func (m *QueryResultBytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResultBytes.Marshal(b, m, deterministic)
}
func (dst *QueryResultBytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResultBytes.Merge(dst, src)
}
func (m *QueryResultBytes) XXX_Size() int {
	return xxx_messageInfo_QueryResultBytes.Size(m)
}
func (m *QueryResultBytes) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResultBytes.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResultBytes proto.InternalMessageInfo

func (m *QueryResultBytes) GetResultBytes() []byte {
	if m != nil {
		return m.ResultBytes
	}
	return nil
}






type QueryResponse struct {
	Results              []*QueryResultBytes `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	HasMore              bool                `protobuf:"varint,2,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"`
	Id                   string              `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Metadata             []byte              `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *QueryResponse) Reset()         { *m = QueryResponse{} }
func (m *QueryResponse) String() string { return proto.CompactTextString(m) }
func (*QueryResponse) ProtoMessage()    {}
func (*QueryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{13}
}
func (m *QueryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResponse.Unmarshal(m, b)
}
func (m *QueryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResponse.Marshal(b, m, deterministic)
}
func (dst *QueryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResponse.Merge(dst, src)
}
func (m *QueryResponse) XXX_Size() int {
	return xxx_messageInfo_QueryResponse.Size(m)
}
func (m *QueryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResponse proto.InternalMessageInfo

func (m *QueryResponse) GetResults() []*QueryResultBytes {
	if m != nil {
		return m.Results
	}
	return nil
}

func (m *QueryResponse) GetHasMore() bool {
	if m != nil {
		return m.HasMore
	}
	return false
}

func (m *QueryResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *QueryResponse) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}



type QueryResponseMetadata struct {
	FetchedRecordsCount  int32    `protobuf:"varint,1,opt,name=fetched_records_count,json=fetchedRecordsCount,proto3" json:"fetched_records_count,omitempty"`
	Bookmark             string   `protobuf:"bytes,2,opt,name=bookmark,proto3" json:"bookmark,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryResponseMetadata) Reset()         { *m = QueryResponseMetadata{} }
func (m *QueryResponseMetadata) String() string { return proto.CompactTextString(m) }
func (*QueryResponseMetadata) ProtoMessage()    {}
func (*QueryResponseMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{14}
}
func (m *QueryResponseMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResponseMetadata.Unmarshal(m, b)
}
func (m *QueryResponseMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResponseMetadata.Marshal(b, m, deterministic)
}
func (dst *QueryResponseMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResponseMetadata.Merge(dst, src)
}
func (m *QueryResponseMetadata) XXX_Size() int {
	return xxx_messageInfo_QueryResponseMetadata.Size(m)
}
func (m *QueryResponseMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResponseMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResponseMetadata proto.InternalMessageInfo

func (m *QueryResponseMetadata) GetFetchedRecordsCount() int32 {
	if m != nil {
		return m.FetchedRecordsCount
	}
	return 0
}

func (m *QueryResponseMetadata) GetBookmark() string {
	if m != nil {
		return m.Bookmark
	}
	return ""
}

type StateMetadata struct {
	Metakey              string   `protobuf:"bytes,1,opt,name=metakey,proto3" json:"metakey,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StateMetadata) Reset()         { *m = StateMetadata{} }
func (m *StateMetadata) String() string { return proto.CompactTextString(m) }
func (*StateMetadata) ProtoMessage()    {}
func (*StateMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{15}
}
func (m *StateMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateMetadata.Unmarshal(m, b)
}
func (m *StateMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateMetadata.Marshal(b, m, deterministic)
}
func (dst *StateMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateMetadata.Merge(dst, src)
}
func (m *StateMetadata) XXX_Size() int {
	return xxx_messageInfo_StateMetadata.Size(m)
}
func (m *StateMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_StateMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_StateMetadata proto.InternalMessageInfo

func (m *StateMetadata) GetMetakey() string {
	if m != nil {
		return m.Metakey
	}
	return ""
}

func (m *StateMetadata) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type StateMetadataResult struct {
	Entries              []*StateMetadata `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *StateMetadataResult) Reset()         { *m = StateMetadataResult{} }
func (m *StateMetadataResult) String() string { return proto.CompactTextString(m) }
func (*StateMetadataResult) ProtoMessage()    {}
func (*StateMetadataResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_chaincode_shim_2ef4390ed21a360d, []int{16}
}
func (m *StateMetadataResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateMetadataResult.Unmarshal(m, b)
}
func (m *StateMetadataResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateMetadataResult.Marshal(b, m, deterministic)
}
func (dst *StateMetadataResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateMetadataResult.Merge(dst, src)
}
func (m *StateMetadataResult) XXX_Size() int {
	return xxx_messageInfo_StateMetadataResult.Size(m)
}
func (m *StateMetadataResult) XXX_DiscardUnknown() {
	xxx_messageInfo_StateMetadataResult.DiscardUnknown(m)
}

var xxx_messageInfo_StateMetadataResult proto.InternalMessageInfo

func (m *StateMetadataResult) GetEntries() []*StateMetadata {
	if m != nil {
		return m.Entries
	}
	return nil
}

func init() {
	proto.RegisterType((*ChaincodeMessage)(nil), "protos.ChaincodeMessage")
	proto.RegisterType((*GetState)(nil), "protos.GetState")
	proto.RegisterType((*GetStateMetadata)(nil), "protos.GetStateMetadata")
	proto.RegisterType((*PutState)(nil), "protos.PutState")
	proto.RegisterType((*PutStateMetadata)(nil), "protos.PutStateMetadata")
	proto.RegisterType((*DelState)(nil), "protos.DelState")
	proto.RegisterType((*GetStateByRange)(nil), "protos.GetStateByRange")
	proto.RegisterType((*GetQueryResult)(nil), "protos.GetQueryResult")
	proto.RegisterType((*QueryMetadata)(nil), "protos.QueryMetadata")
	proto.RegisterType((*GetHistoryForKey)(nil), "protos.GetHistoryForKey")
	proto.RegisterType((*QueryStateNext)(nil), "protos.QueryStateNext")
	proto.RegisterType((*QueryStateClose)(nil), "protos.QueryStateClose")
	proto.RegisterType((*QueryResultBytes)(nil), "protos.QueryResultBytes")
	proto.RegisterType((*QueryResponse)(nil), "protos.QueryResponse")
	proto.RegisterType((*QueryResponseMetadata)(nil), "protos.QueryResponseMetadata")
	proto.RegisterType((*StateMetadata)(nil), "protos.StateMetadata")
	proto.RegisterType((*StateMetadataResult)(nil), "protos.StateMetadataResult")
	proto.RegisterEnum("protos.ChaincodeMessage_Type", ChaincodeMessage_Type_name, ChaincodeMessage_Type_value)
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type ChaincodeSupportClient interface {
	Register(ctx context.Context, opts ...grpc.CallOption) (ChaincodeSupport_RegisterClient, error)
}

type chaincodeSupportClient struct {
	cc *grpc.ClientConn
}

func NewChaincodeSupportClient(cc *grpc.ClientConn) ChaincodeSupportClient {
	return &chaincodeSupportClient{cc}
}

func (c *chaincodeSupportClient) Register(ctx context.Context, opts ...grpc.CallOption) (ChaincodeSupport_RegisterClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ChaincodeSupport_serviceDesc.Streams[0], "/protos.ChaincodeSupport/Register", opts...)
	if err != nil {
		return nil, err
	}
	x := &chaincodeSupportRegisterClient{stream}
	return x, nil
}

type ChaincodeSupport_RegisterClient interface {
	Send(*ChaincodeMessage) error
	Recv() (*ChaincodeMessage, error)
	grpc.ClientStream
}

type chaincodeSupportRegisterClient struct {
	grpc.ClientStream
}

func (x *chaincodeSupportRegisterClient) Send(m *ChaincodeMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chaincodeSupportRegisterClient) Recv() (*ChaincodeMessage, error) {
	m := new(ChaincodeMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}


type ChaincodeSupportServer interface {
	Register(ChaincodeSupport_RegisterServer) error
}

func RegisterChaincodeSupportServer(s *grpc.Server, srv ChaincodeSupportServer) {
	s.RegisterService(&_ChaincodeSupport_serviceDesc, srv)
}

func _ChaincodeSupport_Register_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChaincodeSupportServer).Register(&chaincodeSupportRegisterServer{stream})
}

type ChaincodeSupport_RegisterServer interface {
	Send(*ChaincodeMessage) error
	Recv() (*ChaincodeMessage, error)
	grpc.ServerStream
}

type chaincodeSupportRegisterServer struct {
	grpc.ServerStream
}

func (x *chaincodeSupportRegisterServer) Send(m *ChaincodeMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chaincodeSupportRegisterServer) Recv() (*ChaincodeMessage, error) {
	m := new(ChaincodeMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ChaincodeSupport_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.ChaincodeSupport",
	HandlerType: (*ChaincodeSupportServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Register",
			Handler:       _ChaincodeSupport_Register_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "peer/chaincode_shim.proto",
}

func init() {
	proto.RegisterFile("peer/chaincode_shim.proto", fileDescriptor_chaincode_shim_2ef4390ed21a360d)
}

var fileDescriptor_chaincode_shim_2ef4390ed21a360d = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x56, 0xcf, 0x6f, 0xe2, 0x46,
	0x14, 0x2e, 0x01, 0x82, 0x79, 0x24, 0x64, 0x76, 0x12, 0x52, 0x82, 0xb4, 0x2d, 0x45, 0x3d, 0xd0,
	0x0b, 0x74, 0x69, 0x0f, 0x3d, 0x54, 0x5a, 0x11, 0x98, 0x10, 0x94, 0x04, 0xd8, 0xb1, 0x13, 0x35,
	0xbd, 0x58, 0xc6, 0x9e, 0x18, 0x2b, 0xc6, 0xe3, 0xda, 0xc3, 0x76, 0xe9, 0xad, 0xd7, 0x1e, 0xfb,
	0xc7, 0xf5, 0xef, 0xa9, 0xc6, 0xbf, 0x02, 0xa4, 0xd9, 0xad, 0xf6, 0x84, 0xbf, 0xf7, 0xbe, 0xf9,
	0xde, 0xaf, 0x79, 0xd8, 0x70, 0xe6, 0x33, 0x16, 0x74, 0xcd, 0x85, 0xe1, 0x78, 0x26, 0xb7, 0x98,
	0x1e, 0x2e, 0x9c, 0x65, 0xc7, 0x0f, 0xb8, 0xe0, 0x78, 0x3f, 0xfa, 0x09, 0x1b, 0x8d, 0x1d, 0x0a,
	0x7b, 0xcf, 0x3c, 0x11, 0x73, 0x1a, 0xc7, 0x91, 0xcf, 0x0f, 0xb8, 0xcf, 0x43, 0xc3, 0x4d, 0x8c,
	0x5f, 0xdb, 0x9c, 0xdb, 0x2e, 0xeb, 0x46, 0x68, 0xbe, 0x7a, 0xe8, 0x0a, 0x67, 0xc9, 0x42, 0x61,
	0x2c, 0xfd, 0x98, 0xd0, 0xfa, 0xa7, 0x08, 0x68, 0x90, 0xea, 0xdd, 0xb0, 0x30, 0x34, 0x6c, 0x86,
	0xdf, 0x40, 0x41, 0xac, 0x7d, 0x56, 0xcf, 0x35, 0x73, 0xed, 0x6a, 0xef, 0x75, 0x4c, 0x0d, 0x3b,
	0xbb, 0xbc, 0x8e, 0xb6, 0xf6, 0x19, 0x8d, 0xa8, 0xf8, 0x27, 0x28, 0x67, 0xd2, 0xf5, 0xbd, 0x66,
	0xae, 0x5d, 0xe9, 0x35, 0x3a, 0x71, 0xf0, 0x4e, 0x1a, 0xbc, 0xa3, 0xa5, 0x0c, 0xfa, 0x44, 0xc6,
	0x75, 0x28, 0xf9, 0xc6, 0xda, 0xe5, 0x86, 0x55, 0xcf, 0x37, 0x73, 0xed, 0x03, 0x9a, 0x42, 0x8c,
	0xa1, 0x20, 0x3e, 0x38, 0x56, 0xbd, 0xd0, 0xcc, 0xb5, 0xcb, 0x34, 0x7a, 0xc6, 0x3d, 0x50, 0xd2,
	0x12, 0xeb, 0xc5, 0x28, 0xcc, 0x69, 0x9a, 0x9e, 0xea, 0xd8, 0x1e, 0xb3, 0x66, 0x89, 0x97, 0x66,
	0x3c, 0xfc, 0x16, 0x8e, 0x76, 0x5a, 0x56, 0xdf, 0xdf, 0x3e, 0x9a, 0x55, 0x46, 0xa4, 0x97, 0x56,
	0xcd, 0x2d, 0x8c, 0x5f, 0x03, 0x98, 0x0b, 0xc3, 0xf3, 0x98, 0xab, 0x3b, 0x56, 0xbd, 0x14, 0xa5,
	0x53, 0x4e, 0x2c, 0x63, 0xab, 0xf5, 0x77, 0x1e, 0x0a, 0xb2, 0x15, 0xf8, 0x10, 0xca, 0xb7, 0x93,
	0x21, 0xb9, 0x18, 0x4f, 0xc8, 0x10, 0x7d, 0x81, 0x0f, 0x40, 0xa1, 0x64, 0x34, 0x56, 0x35, 0x42,
	0x51, 0x0e, 0x57, 0x01, 0x52, 0x44, 0x86, 0x68, 0x0f, 0x2b, 0x50, 0x18, 0x4f, 0xc6, 0x1a, 0xca,
	0xe3, 0x32, 0x14, 0x29, 0xe9, 0x0f, 0xef, 0x51, 0x01, 0x1f, 0x41, 0x45, 0xa3, 0xfd, 0x89, 0xda,
	0x1f, 0x68, 0xe3, 0xe9, 0x04, 0x15, 0xa5, 0xe4, 0x60, 0x7a, 0x33, 0xbb, 0x26, 0x1a, 0x19, 0xa2,
	0x7d, 0x49, 0x25, 0x94, 0x4e, 0x29, 0x2a, 0x49, 0xcf, 0x88, 0x68, 0xba, 0xaa, 0xf5, 0x35, 0x82,
	0x14, 0x09, 0x67, 0xb7, 0x29, 0x2c, 0x4b, 0x38, 0x24, 0xd7, 0x09, 0x04, 0x7c, 0x02, 0x68, 0x3c,
	0xb9, 0x9b, 0x5e, 0x11, 0x7d, 0x70, 0xd9, 0x1f, 0x4f, 0x06, 0xd3, 0x21, 0x41, 0x95, 0x38, 0x41,
	0x75, 0x36, 0x9d, 0xa8, 0x04, 0x1d, 0xe2, 0x53, 0xc0, 0x99, 0xa0, 0x7e, 0x7e, 0xaf, 0xd3, 0xfe,
	0x64, 0x44, 0x50, 0x55, 0x9e, 0x95, 0xf6, 0x77, 0xb7, 0x84, 0xde, 0xeb, 0x94, 0xa8, 0xb7, 0xd7,
	0x1a, 0x3a, 0x92, 0xd6, 0xd8, 0x12, 0xf3, 0x27, 0xe4, 0x17, 0x0d, 0x21, 0x5c, 0x83, 0x57, 0x9b,
	0xd6, 0xc1, 0xf5, 0x54, 0x25, 0xe8, 0x95, 0xcc, 0xe6, 0x8a, 0x90, 0x59, 0xff, 0x7a, 0x7c, 0x47,
	0x10, 0xc6, 0x5f, 0xc2, 0xb1, 0x54, 0xbc, 0x1c, 0xab, 0xda, 0x94, 0xde, 0xeb, 0x17, 0x53, 0xaa,
	0x5f, 0x91, 0x7b, 0x74, 0xbc, 0x9d, 0xc2, 0x0d, 0xd1, 0xfa, 0xc3, 0xbe, 0xd6, 0x47, 0x27, 0xd2,
	0x9e, 0x15, 0xf7, 0x64, 0xaf, 0xe1, 0x33, 0xa8, 0x49, 0xfe, 0x8c, 0x8e, 0xef, 0xa4, 0x47, 0x5a,
	0xf5, 0xcb, 0xbe, 0x7a, 0x89, 0x4e, 0x5b, 0x3f, 0x83, 0x32, 0x62, 0x42, 0x15, 0x86, 0x60, 0x18,
	0x41, 0xfe, 0x91, 0xad, 0xa3, 0xeb, 0x5c, 0xa6, 0xf2, 0x11, 0x7f, 0x05, 0x60, 0x72, 0xd7, 0x65,
	0xa6, 0x70, 0xb8, 0x17, 0xdd, 0xd7, 0x32, 0xdd, 0xb0, 0xb4, 0x86, 0x80, 0xd2, 0xd3, 0x37, 0x4c,
	0x18, 0x96, 0x21, 0x8c, 0xcf, 0x50, 0xa1, 0xa0, 0xcc, 0x56, 0x2f, 0xe6, 0x70, 0x02, 0xc5, 0xf7,
	0x86, 0xbb, 0x62, 0xd1, 0xc1, 0x03, 0x1a, 0x83, 0x1d, 0xcd, 0xfc, 0x33, 0xcd, 0xdf, 0x01, 0xa5,
	0x9a, 0xff, 0x3b, 0xb3, 0x67, 0x2a, 0xf8, 0x0d, 0x28, 0xcb, 0xe4, 0x74, 0xb4, 0x5e, 0x95, 0x5e,
	0x2d, 0x5b, 0xa3, 0x4d, 0x69, 0x9a, 0xd1, 0x64, 0x43, 0x87, 0xcc, 0xfd, 0xdc, 0x86, 0xfe, 0x99,
	0x83, 0xa3, 0xb4, 0xa3, 0xe7, 0x6b, 0x6a, 0x78, 0x36, 0xc3, 0x0d, 0x50, 0x42, 0x61, 0x04, 0xe2,
	0x2a, 0x93, 0xca, 0x30, 0x3e, 0x85, 0x7d, 0xe6, 0x59, 0xd2, 0x13, 0x6b, 0x25, 0xe8, 0x93, 0x85,
	0x35, 0x76, 0x0a, 0x3b, 0xd8, 0xa8, 0x60, 0x0e, 0xd5, 0x11, 0x13, 0xef, 0x56, 0x2c, 0x58, 0x53,
	0x16, 0xae, 0x5c, 0x21, 0x47, 0xf0, 0x9b, 0x84, 0x49, 0xf8, 0x18, 0x7c, 0xaa, 0x96, 0xad, 0x18,
	0xf9, 0x9d, 0x18, 0x23, 0x38, 0x8c, 0x02, 0x64, 0xb3, 0x69, 0x80, 0xe2, 0x1b, 0x36, 0x53, 0x9d,
	0x3f, 0xe2, 0xff, 0xd3, 0x22, 0xcd, 0xb0, 0xf4, 0xcd, 0x39, 0x7f, 0x5c, 0x1a, 0xc1, 0x63, 0x12,
	0x26, 0xc3, 0xad, 0x6f, 0xa3, 0x1b, 0x78, 0xe9, 0x84, 0x82, 0x07, 0xeb, 0x0b, 0x1e, 0xc8, 0xe2,
	0x9f, 0xb5, 0xbd, 0xd5, 0x84, 0x6a, 0x14, 0x2e, 0xea, 0xeb, 0x84, 0x7d, 0x10, 0xb8, 0x0a, 0x7b,
	0x8e, 0x95, 0x50, 0xf6, 0x1c, 0xab, 0xf5, 0x0d, 0x1c, 0x3d, 0x31, 0x06, 0x2e, 0x0f, 0xd9, 0x33,
	0xca, 0x8f, 0x80, 0x36, 0x9a, 0x72, 0xbe, 0x16, 0x2c, 0xc4, 0x4d, 0xa8, 0x04, 0x4f, 0x30, 0x22,
	0x1f, 0xd0, 0x4d, 0x53, 0xeb, 0xaf, 0x5c, 0x52, 0x2a, 0x65, 0xa1, 0xcf, 0xbd, 0x90, 0xe1, 0x1e,
	0x94, 0x62, 0x82, 0xe4, 0xe7, 0xdb, 0x95, 0x5e, 0x3d, 0xbd, 0x53, 0xbb, 0xf2, 0x34, 0x25, 0xe2,
	0x33, 0x50, 0x16, 0x46, 0xa8, 0x2f, 0x79, 0x10, 0xef, 0x81, 0x42, 0x4b, 0x0b, 0x23, 0xbc, 0xe1,
	0x41, 0x9a, 0x66, 0x3e, 0x4d, 0xf3, 0xa3, 0xa3, 0xb5, 0xa1, 0xb6, 0x95, 0x4b, 0xd6, 0xfe, 0x1e,
	0xd4, 0x1e, 0x98, 0x30, 0x17, 0xcc, 0xd2, 0x03, 0x66, 0xf2, 0xc0, 0x0a, 0x75, 0x93, 0xaf, 0x3c,
	0x91, 0xcc, 0xe2, 0x38, 0x71, 0xd2, 0xd8, 0x37, 0x90, 0xae, 0x8f, 0x8e, 0xe5, 0x2d, 0x1c, 0x6e,
	0xef, 0x5e, 0x1d, 0x4a, 0x32, 0x8b, 0xa7, 0xb9, 0xa4, 0xf0, 0xbf, 0xf7, 0xbb, 0x75, 0x01, 0xc7,
	0xdb, 0x1b, 0x16, 0xdf, 0xc4, 0x2e, 0x94, 0x98, 0x27, 0x02, 0x87, 0xa5, 0xbd, 0x7b, 0x61, 0x1f,
	0x53, 0x56, 0xef, 0x6e, 0xe3, 0xbd, 0xad, 0xae, 0x7c, 0x9f, 0x07, 0x02, 0x9f, 0x83, 0x42, 0x99,
	0xed, 0x84, 0x82, 0x05, 0xb8, 0xfe, 0xd2, 0x5b, 0xbb, 0xf1, 0xa2, 0xa7, 0x9d, 0xfb, 0x3e, 0x77,
	0x3e, 0x85, 0x16, 0x0f, 0xec, 0xce, 0x62, 0xed, 0xb3, 0xc0, 0x65, 0x96, 0xcd, 0x82, 0xce, 0x83,
	0x31, 0x0f, 0x1c, 0x33, 0x3d, 0x25, 0x3f, 0x33, 0x7e, 0xfd, 0xce, 0x76, 0xc4, 0x62, 0x35, 0xef,
	0x98, 0x7c, 0xd9, 0xdd, 0xa0, 0x76, 0x63, 0x6a, 0xfc, 0xb9, 0x11, 0x76, 0x25, 0x75, 0x1e, 0x7f,
	0xbb, 0xfc, 0xf0, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x45, 0xde, 0x24, 0x26, 0xdf, 0x08, 0x00,
	0x00,
}
