


package discovery

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	gossip "github.com/mcc-github/blockchain/protos/gossip"
	msp "github.com/mcc-github/blockchain/protos/msp"
	grpc "google.golang.org/grpc"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 






type SignedRequest struct {
	Payload              []byte   `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedRequest) Reset()         { *m = SignedRequest{} }
func (m *SignedRequest) String() string { return proto.CompactTextString(m) }
func (*SignedRequest) ProtoMessage()    {}
func (*SignedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{0}
}

func (m *SignedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedRequest.Unmarshal(m, b)
}
func (m *SignedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedRequest.Marshal(b, m, deterministic)
}
func (m *SignedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedRequest.Merge(m, src)
}
func (m *SignedRequest) XXX_Size() int {
	return xxx_messageInfo_SignedRequest.Size(m)
}
func (m *SignedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SignedRequest proto.InternalMessageInfo

func (m *SignedRequest) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *SignedRequest) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}



type Request struct {
	
	
	Authentication *AuthInfo `protobuf:"bytes,1,opt,name=authentication,proto3" json:"authentication,omitempty"`
	
	Queries              []*Query `protobuf:"bytes,2,rep,name=queries,proto3" json:"queries,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{1}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetAuthentication() *AuthInfo {
	if m != nil {
		return m.Authentication
	}
	return nil
}

func (m *Request) GetQueries() []*Query {
	if m != nil {
		return m.Queries
	}
	return nil
}

type Response struct {
	
	Results              []*QueryResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{2}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetResults() []*QueryResult {
	if m != nil {
		return m.Results
	}
	return nil
}



type AuthInfo struct {
	
	
	
	ClientIdentity []byte `protobuf:"bytes,1,opt,name=client_identity,json=clientIdentity,proto3" json:"client_identity,omitempty"`
	
	
	
	
	
	
	
	ClientTlsCertHash    []byte   `protobuf:"bytes,2,opt,name=client_tls_cert_hash,json=clientTlsCertHash,proto3" json:"client_tls_cert_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthInfo) Reset()         { *m = AuthInfo{} }
func (m *AuthInfo) String() string { return proto.CompactTextString(m) }
func (*AuthInfo) ProtoMessage()    {}
func (*AuthInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{3}
}

func (m *AuthInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthInfo.Unmarshal(m, b)
}
func (m *AuthInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthInfo.Marshal(b, m, deterministic)
}
func (m *AuthInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthInfo.Merge(m, src)
}
func (m *AuthInfo) XXX_Size() int {
	return xxx_messageInfo_AuthInfo.Size(m)
}
func (m *AuthInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AuthInfo proto.InternalMessageInfo

func (m *AuthInfo) GetClientIdentity() []byte {
	if m != nil {
		return m.ClientIdentity
	}
	return nil
}

func (m *AuthInfo) GetClientTlsCertHash() []byte {
	if m != nil {
		return m.ClientTlsCertHash
	}
	return nil
}


type Query struct {
	Channel string `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	
	
	
	
	
	Query                isQuery_Query `protobuf_oneof:"query"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{4}
}

func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (m *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(m, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

type isQuery_Query interface {
	isQuery_Query()
}

type Query_ConfigQuery struct {
	ConfigQuery *ConfigQuery `protobuf:"bytes,2,opt,name=config_query,json=configQuery,proto3,oneof"`
}

type Query_PeerQuery struct {
	PeerQuery *PeerMembershipQuery `protobuf:"bytes,3,opt,name=peer_query,json=peerQuery,proto3,oneof"`
}

type Query_CcQuery struct {
	CcQuery *ChaincodeQuery `protobuf:"bytes,4,opt,name=cc_query,json=ccQuery,proto3,oneof"`
}

type Query_LocalPeers struct {
	LocalPeers *LocalPeerQuery `protobuf:"bytes,5,opt,name=local_peers,json=localPeers,proto3,oneof"`
}

func (*Query_ConfigQuery) isQuery_Query() {}

func (*Query_PeerQuery) isQuery_Query() {}

func (*Query_CcQuery) isQuery_Query() {}

func (*Query_LocalPeers) isQuery_Query() {}

func (m *Query) GetQuery() isQuery_Query {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *Query) GetConfigQuery() *ConfigQuery {
	if x, ok := m.GetQuery().(*Query_ConfigQuery); ok {
		return x.ConfigQuery
	}
	return nil
}

func (m *Query) GetPeerQuery() *PeerMembershipQuery {
	if x, ok := m.GetQuery().(*Query_PeerQuery); ok {
		return x.PeerQuery
	}
	return nil
}

func (m *Query) GetCcQuery() *ChaincodeQuery {
	if x, ok := m.GetQuery().(*Query_CcQuery); ok {
		return x.CcQuery
	}
	return nil
}

func (m *Query) GetLocalPeers() *LocalPeerQuery {
	if x, ok := m.GetQuery().(*Query_LocalPeers); ok {
		return x.LocalPeers
	}
	return nil
}


func (*Query) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Query_ConfigQuery)(nil),
		(*Query_PeerQuery)(nil),
		(*Query_CcQuery)(nil),
		(*Query_LocalPeers)(nil),
	}
}





type QueryResult struct {
	
	
	
	
	
	Result               isQueryResult_Result `protobuf_oneof:"result"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *QueryResult) Reset()         { *m = QueryResult{} }
func (m *QueryResult) String() string { return proto.CompactTextString(m) }
func (*QueryResult) ProtoMessage()    {}
func (*QueryResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{5}
}

func (m *QueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResult.Unmarshal(m, b)
}
func (m *QueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResult.Marshal(b, m, deterministic)
}
func (m *QueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResult.Merge(m, src)
}
func (m *QueryResult) XXX_Size() int {
	return xxx_messageInfo_QueryResult.Size(m)
}
func (m *QueryResult) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResult proto.InternalMessageInfo

type isQueryResult_Result interface {
	isQueryResult_Result()
}

type QueryResult_Error struct {
	Error *Error `protobuf:"bytes,1,opt,name=error,proto3,oneof"`
}

type QueryResult_ConfigResult struct {
	ConfigResult *ConfigResult `protobuf:"bytes,2,opt,name=config_result,json=configResult,proto3,oneof"`
}

type QueryResult_CcQueryRes struct {
	CcQueryRes *ChaincodeQueryResult `protobuf:"bytes,3,opt,name=cc_query_res,json=ccQueryRes,proto3,oneof"`
}

type QueryResult_Members struct {
	Members *PeerMembershipResult `protobuf:"bytes,4,opt,name=members,proto3,oneof"`
}

func (*QueryResult_Error) isQueryResult_Result() {}

func (*QueryResult_ConfigResult) isQueryResult_Result() {}

func (*QueryResult_CcQueryRes) isQueryResult_Result() {}

func (*QueryResult_Members) isQueryResult_Result() {}

func (m *QueryResult) GetResult() isQueryResult_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *QueryResult) GetError() *Error {
	if x, ok := m.GetResult().(*QueryResult_Error); ok {
		return x.Error
	}
	return nil
}

func (m *QueryResult) GetConfigResult() *ConfigResult {
	if x, ok := m.GetResult().(*QueryResult_ConfigResult); ok {
		return x.ConfigResult
	}
	return nil
}

func (m *QueryResult) GetCcQueryRes() *ChaincodeQueryResult {
	if x, ok := m.GetResult().(*QueryResult_CcQueryRes); ok {
		return x.CcQueryRes
	}
	return nil
}

func (m *QueryResult) GetMembers() *PeerMembershipResult {
	if x, ok := m.GetResult().(*QueryResult_Members); ok {
		return x.Members
	}
	return nil
}


func (*QueryResult) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*QueryResult_Error)(nil),
		(*QueryResult_ConfigResult)(nil),
		(*QueryResult_CcQueryRes)(nil),
		(*QueryResult_Members)(nil),
	}
}


type ConfigQuery struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigQuery) Reset()         { *m = ConfigQuery{} }
func (m *ConfigQuery) String() string { return proto.CompactTextString(m) }
func (*ConfigQuery) ProtoMessage()    {}
func (*ConfigQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{6}
}

func (m *ConfigQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigQuery.Unmarshal(m, b)
}
func (m *ConfigQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigQuery.Marshal(b, m, deterministic)
}
func (m *ConfigQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigQuery.Merge(m, src)
}
func (m *ConfigQuery) XXX_Size() int {
	return xxx_messageInfo_ConfigQuery.Size(m)
}
func (m *ConfigQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigQuery.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigQuery proto.InternalMessageInfo

type ConfigResult struct {
	
	Msps map[string]*msp.FabricMSPConfig `protobuf:"bytes,1,rep,name=msps,proto3" json:"msps,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	
	Orderers             map[string]*Endpoints `protobuf:"bytes,2,rep,name=orderers,proto3" json:"orderers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ConfigResult) Reset()         { *m = ConfigResult{} }
func (m *ConfigResult) String() string { return proto.CompactTextString(m) }
func (*ConfigResult) ProtoMessage()    {}
func (*ConfigResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{7}
}

func (m *ConfigResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigResult.Unmarshal(m, b)
}
func (m *ConfigResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigResult.Marshal(b, m, deterministic)
}
func (m *ConfigResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigResult.Merge(m, src)
}
func (m *ConfigResult) XXX_Size() int {
	return xxx_messageInfo_ConfigResult.Size(m)
}
func (m *ConfigResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigResult.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigResult proto.InternalMessageInfo

func (m *ConfigResult) GetMsps() map[string]*msp.FabricMSPConfig {
	if m != nil {
		return m.Msps
	}
	return nil
}

func (m *ConfigResult) GetOrderers() map[string]*Endpoints {
	if m != nil {
		return m.Orderers
	}
	return nil
}






type PeerMembershipQuery struct {
	Filter               *ChaincodeInterest `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *PeerMembershipQuery) Reset()         { *m = PeerMembershipQuery{} }
func (m *PeerMembershipQuery) String() string { return proto.CompactTextString(m) }
func (*PeerMembershipQuery) ProtoMessage()    {}
func (*PeerMembershipQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{8}
}

func (m *PeerMembershipQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerMembershipQuery.Unmarshal(m, b)
}
func (m *PeerMembershipQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerMembershipQuery.Marshal(b, m, deterministic)
}
func (m *PeerMembershipQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerMembershipQuery.Merge(m, src)
}
func (m *PeerMembershipQuery) XXX_Size() int {
	return xxx_messageInfo_PeerMembershipQuery.Size(m)
}
func (m *PeerMembershipQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerMembershipQuery.DiscardUnknown(m)
}

var xxx_messageInfo_PeerMembershipQuery proto.InternalMessageInfo

func (m *PeerMembershipQuery) GetFilter() *ChaincodeInterest {
	if m != nil {
		return m.Filter
	}
	return nil
}


type PeerMembershipResult struct {
	PeersByOrg           map[string]*Peers `protobuf:"bytes,1,rep,name=peers_by_org,json=peersByOrg,proto3" json:"peers_by_org,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *PeerMembershipResult) Reset()         { *m = PeerMembershipResult{} }
func (m *PeerMembershipResult) String() string { return proto.CompactTextString(m) }
func (*PeerMembershipResult) ProtoMessage()    {}
func (*PeerMembershipResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{9}
}

func (m *PeerMembershipResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerMembershipResult.Unmarshal(m, b)
}
func (m *PeerMembershipResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerMembershipResult.Marshal(b, m, deterministic)
}
func (m *PeerMembershipResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerMembershipResult.Merge(m, src)
}
func (m *PeerMembershipResult) XXX_Size() int {
	return xxx_messageInfo_PeerMembershipResult.Size(m)
}
func (m *PeerMembershipResult) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerMembershipResult.DiscardUnknown(m)
}

var xxx_messageInfo_PeerMembershipResult proto.InternalMessageInfo

func (m *PeerMembershipResult) GetPeersByOrg() map[string]*Peers {
	if m != nil {
		return m.PeersByOrg
	}
	return nil
}





type ChaincodeQuery struct {
	Interests            []*ChaincodeInterest `protobuf:"bytes,1,rep,name=interests,proto3" json:"interests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ChaincodeQuery) Reset()         { *m = ChaincodeQuery{} }
func (m *ChaincodeQuery) String() string { return proto.CompactTextString(m) }
func (*ChaincodeQuery) ProtoMessage()    {}
func (*ChaincodeQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{10}
}

func (m *ChaincodeQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeQuery.Unmarshal(m, b)
}
func (m *ChaincodeQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeQuery.Marshal(b, m, deterministic)
}
func (m *ChaincodeQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeQuery.Merge(m, src)
}
func (m *ChaincodeQuery) XXX_Size() int {
	return xxx_messageInfo_ChaincodeQuery.Size(m)
}
func (m *ChaincodeQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeQuery.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeQuery proto.InternalMessageInfo

func (m *ChaincodeQuery) GetInterests() []*ChaincodeInterest {
	if m != nil {
		return m.Interests
	}
	return nil
}




type ChaincodeInterest struct {
	Chaincodes           []*ChaincodeCall `protobuf:"bytes,1,rep,name=chaincodes,proto3" json:"chaincodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ChaincodeInterest) Reset()         { *m = ChaincodeInterest{} }
func (m *ChaincodeInterest) String() string { return proto.CompactTextString(m) }
func (*ChaincodeInterest) ProtoMessage()    {}
func (*ChaincodeInterest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{11}
}

func (m *ChaincodeInterest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeInterest.Unmarshal(m, b)
}
func (m *ChaincodeInterest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeInterest.Marshal(b, m, deterministic)
}
func (m *ChaincodeInterest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeInterest.Merge(m, src)
}
func (m *ChaincodeInterest) XXX_Size() int {
	return xxx_messageInfo_ChaincodeInterest.Size(m)
}
func (m *ChaincodeInterest) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeInterest.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeInterest proto.InternalMessageInfo

func (m *ChaincodeInterest) GetChaincodes() []*ChaincodeCall {
	if m != nil {
		return m.Chaincodes
	}
	return nil
}



type ChaincodeCall struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	CollectionNames      []string `protobuf:"bytes,2,rep,name=collection_names,json=collectionNames,proto3" json:"collection_names,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeCall) Reset()         { *m = ChaincodeCall{} }
func (m *ChaincodeCall) String() string { return proto.CompactTextString(m) }
func (*ChaincodeCall) ProtoMessage()    {}
func (*ChaincodeCall) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{12}
}

func (m *ChaincodeCall) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeCall.Unmarshal(m, b)
}
func (m *ChaincodeCall) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeCall.Marshal(b, m, deterministic)
}
func (m *ChaincodeCall) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeCall.Merge(m, src)
}
func (m *ChaincodeCall) XXX_Size() int {
	return xxx_messageInfo_ChaincodeCall.Size(m)
}
func (m *ChaincodeCall) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeCall.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeCall proto.InternalMessageInfo

func (m *ChaincodeCall) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ChaincodeCall) GetCollectionNames() []string {
	if m != nil {
		return m.CollectionNames
	}
	return nil
}



type ChaincodeQueryResult struct {
	Content              []*EndorsementDescriptor `protobuf:"bytes,1,rep,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ChaincodeQueryResult) Reset()         { *m = ChaincodeQueryResult{} }
func (m *ChaincodeQueryResult) String() string { return proto.CompactTextString(m) }
func (*ChaincodeQueryResult) ProtoMessage()    {}
func (*ChaincodeQueryResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{13}
}

func (m *ChaincodeQueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeQueryResult.Unmarshal(m, b)
}
func (m *ChaincodeQueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeQueryResult.Marshal(b, m, deterministic)
}
func (m *ChaincodeQueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeQueryResult.Merge(m, src)
}
func (m *ChaincodeQueryResult) XXX_Size() int {
	return xxx_messageInfo_ChaincodeQueryResult.Size(m)
}
func (m *ChaincodeQueryResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeQueryResult.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeQueryResult proto.InternalMessageInfo

func (m *ChaincodeQueryResult) GetContent() []*EndorsementDescriptor {
	if m != nil {
		return m.Content
	}
	return nil
}


type LocalPeerQuery struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LocalPeerQuery) Reset()         { *m = LocalPeerQuery{} }
func (m *LocalPeerQuery) String() string { return proto.CompactTextString(m) }
func (*LocalPeerQuery) ProtoMessage()    {}
func (*LocalPeerQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{14}
}

func (m *LocalPeerQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LocalPeerQuery.Unmarshal(m, b)
}
func (m *LocalPeerQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LocalPeerQuery.Marshal(b, m, deterministic)
}
func (m *LocalPeerQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LocalPeerQuery.Merge(m, src)
}
func (m *LocalPeerQuery) XXX_Size() int {
	return xxx_messageInfo_LocalPeerQuery.Size(m)
}
func (m *LocalPeerQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_LocalPeerQuery.DiscardUnknown(m)
}

var xxx_messageInfo_LocalPeerQuery proto.InternalMessageInfo













type EndorsementDescriptor struct {
	Chaincode string `protobuf:"bytes,1,opt,name=chaincode,proto3" json:"chaincode,omitempty"`
	
	EndorsersByGroups map[string]*Peers `protobuf:"bytes,2,rep,name=endorsers_by_groups,json=endorsersByGroups,proto3" json:"endorsers_by_groups,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	
	
	
	Layouts              []*Layout `protobuf:"bytes,3,rep,name=layouts,proto3" json:"layouts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *EndorsementDescriptor) Reset()         { *m = EndorsementDescriptor{} }
func (m *EndorsementDescriptor) String() string { return proto.CompactTextString(m) }
func (*EndorsementDescriptor) ProtoMessage()    {}
func (*EndorsementDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{15}
}

func (m *EndorsementDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EndorsementDescriptor.Unmarshal(m, b)
}
func (m *EndorsementDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EndorsementDescriptor.Marshal(b, m, deterministic)
}
func (m *EndorsementDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EndorsementDescriptor.Merge(m, src)
}
func (m *EndorsementDescriptor) XXX_Size() int {
	return xxx_messageInfo_EndorsementDescriptor.Size(m)
}
func (m *EndorsementDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_EndorsementDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_EndorsementDescriptor proto.InternalMessageInfo

func (m *EndorsementDescriptor) GetChaincode() string {
	if m != nil {
		return m.Chaincode
	}
	return ""
}

func (m *EndorsementDescriptor) GetEndorsersByGroups() map[string]*Peers {
	if m != nil {
		return m.EndorsersByGroups
	}
	return nil
}

func (m *EndorsementDescriptor) GetLayouts() []*Layout {
	if m != nil {
		return m.Layouts
	}
	return nil
}



type Layout struct {
	
	
	QuantitiesByGroup    map[string]uint32 `protobuf:"bytes,1,rep,name=quantities_by_group,json=quantitiesByGroup,proto3" json:"quantities_by_group,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Layout) Reset()         { *m = Layout{} }
func (m *Layout) String() string { return proto.CompactTextString(m) }
func (*Layout) ProtoMessage()    {}
func (*Layout) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{16}
}

func (m *Layout) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Layout.Unmarshal(m, b)
}
func (m *Layout) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Layout.Marshal(b, m, deterministic)
}
func (m *Layout) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Layout.Merge(m, src)
}
func (m *Layout) XXX_Size() int {
	return xxx_messageInfo_Layout.Size(m)
}
func (m *Layout) XXX_DiscardUnknown() {
	xxx_messageInfo_Layout.DiscardUnknown(m)
}

var xxx_messageInfo_Layout proto.InternalMessageInfo

func (m *Layout) GetQuantitiesByGroup() map[string]uint32 {
	if m != nil {
		return m.QuantitiesByGroup
	}
	return nil
}


type Peers struct {
	Peers                []*Peer  `protobuf:"bytes,1,rep,name=peers,proto3" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peers) Reset()         { *m = Peers{} }
func (m *Peers) String() string { return proto.CompactTextString(m) }
func (*Peers) ProtoMessage()    {}
func (*Peers) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{17}
}

func (m *Peers) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peers.Unmarshal(m, b)
}
func (m *Peers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peers.Marshal(b, m, deterministic)
}
func (m *Peers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peers.Merge(m, src)
}
func (m *Peers) XXX_Size() int {
	return xxx_messageInfo_Peers.Size(m)
}
func (m *Peers) XXX_DiscardUnknown() {
	xxx_messageInfo_Peers.DiscardUnknown(m)
}

var xxx_messageInfo_Peers proto.InternalMessageInfo

func (m *Peers) GetPeers() []*Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}



type Peer struct {
	
	StateInfo *gossip.Envelope `protobuf:"bytes,1,opt,name=state_info,json=stateInfo,proto3" json:"state_info,omitempty"`
	
	MembershipInfo *gossip.Envelope `protobuf:"bytes,2,opt,name=membership_info,json=membershipInfo,proto3" json:"membership_info,omitempty"`
	
	Identity             []byte   `protobuf:"bytes,3,opt,name=identity,proto3" json:"identity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peer) Reset()         { *m = Peer{} }
func (m *Peer) String() string { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()    {}
func (*Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{18}
}

func (m *Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer.Unmarshal(m, b)
}
func (m *Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer.Marshal(b, m, deterministic)
}
func (m *Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer.Merge(m, src)
}
func (m *Peer) XXX_Size() int {
	return xxx_messageInfo_Peer.Size(m)
}
func (m *Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_Peer proto.InternalMessageInfo

func (m *Peer) GetStateInfo() *gossip.Envelope {
	if m != nil {
		return m.StateInfo
	}
	return nil
}

func (m *Peer) GetMembershipInfo() *gossip.Envelope {
	if m != nil {
		return m.MembershipInfo
	}
	return nil
}

func (m *Peer) GetIdentity() []byte {
	if m != nil {
		return m.Identity
	}
	return nil
}


type Error struct {
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{19}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}


type Endpoints struct {
	Endpoint             []*Endpoint `protobuf:"bytes,1,rep,name=endpoint,proto3" json:"endpoint,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Endpoints) Reset()         { *m = Endpoints{} }
func (m *Endpoints) String() string { return proto.CompactTextString(m) }
func (*Endpoints) ProtoMessage()    {}
func (*Endpoints) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{20}
}

func (m *Endpoints) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoints.Unmarshal(m, b)
}
func (m *Endpoints) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoints.Marshal(b, m, deterministic)
}
func (m *Endpoints) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoints.Merge(m, src)
}
func (m *Endpoints) XXX_Size() int {
	return xxx_messageInfo_Endpoints.Size(m)
}
func (m *Endpoints) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoints.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoints proto.InternalMessageInfo

func (m *Endpoints) GetEndpoint() []*Endpoint {
	if m != nil {
		return m.Endpoint
	}
	return nil
}


type Endpoint struct {
	Host                 string   `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}
func (*Endpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce69bf33982206ff, []int{21}
}

func (m *Endpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoint.Unmarshal(m, b)
}
func (m *Endpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoint.Marshal(b, m, deterministic)
}
func (m *Endpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoint.Merge(m, src)
}
func (m *Endpoint) XXX_Size() int {
	return xxx_messageInfo_Endpoint.Size(m)
}
func (m *Endpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoint.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoint proto.InternalMessageInfo

func (m *Endpoint) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Endpoint) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func init() {
	proto.RegisterType((*SignedRequest)(nil), "discovery.SignedRequest")
	proto.RegisterType((*Request)(nil), "discovery.Request")
	proto.RegisterType((*Response)(nil), "discovery.Response")
	proto.RegisterType((*AuthInfo)(nil), "discovery.AuthInfo")
	proto.RegisterType((*Query)(nil), "discovery.Query")
	proto.RegisterType((*QueryResult)(nil), "discovery.QueryResult")
	proto.RegisterType((*ConfigQuery)(nil), "discovery.ConfigQuery")
	proto.RegisterType((*ConfigResult)(nil), "discovery.ConfigResult")
	proto.RegisterMapType((map[string]*msp.FabricMSPConfig)(nil), "discovery.ConfigResult.MspsEntry")
	proto.RegisterMapType((map[string]*Endpoints)(nil), "discovery.ConfigResult.OrderersEntry")
	proto.RegisterType((*PeerMembershipQuery)(nil), "discovery.PeerMembershipQuery")
	proto.RegisterType((*PeerMembershipResult)(nil), "discovery.PeerMembershipResult")
	proto.RegisterMapType((map[string]*Peers)(nil), "discovery.PeerMembershipResult.PeersByOrgEntry")
	proto.RegisterType((*ChaincodeQuery)(nil), "discovery.ChaincodeQuery")
	proto.RegisterType((*ChaincodeInterest)(nil), "discovery.ChaincodeInterest")
	proto.RegisterType((*ChaincodeCall)(nil), "discovery.ChaincodeCall")
	proto.RegisterType((*ChaincodeQueryResult)(nil), "discovery.ChaincodeQueryResult")
	proto.RegisterType((*LocalPeerQuery)(nil), "discovery.LocalPeerQuery")
	proto.RegisterType((*EndorsementDescriptor)(nil), "discovery.EndorsementDescriptor")
	proto.RegisterMapType((map[string]*Peers)(nil), "discovery.EndorsementDescriptor.EndorsersByGroupsEntry")
	proto.RegisterType((*Layout)(nil), "discovery.Layout")
	proto.RegisterMapType((map[string]uint32)(nil), "discovery.Layout.QuantitiesByGroupEntry")
	proto.RegisterType((*Peers)(nil), "discovery.Peers")
	proto.RegisterType((*Peer)(nil), "discovery.Peer")
	proto.RegisterType((*Error)(nil), "discovery.Error")
	proto.RegisterType((*Endpoints)(nil), "discovery.Endpoints")
	proto.RegisterType((*Endpoint)(nil), "discovery.Endpoint")
}

func init() { proto.RegisterFile("discovery/protocol.proto", fileDescriptor_ce69bf33982206ff) }

var fileDescriptor_ce69bf33982206ff = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0x5b, 0x6f, 0x1b, 0x45,
	0x14, 0x8e, 0x9d, 0x38, 0xb6, 0x8f, 0xe3, 0x5c, 0x26, 0xa6, 0x18, 0xab, 0xa2, 0xe9, 0x4a, 0x85,
	0x50, 0xa4, 0x75, 0x15, 0x6e, 0x6d, 0x12, 0x81, 0x9a, 0x4b, 0x9b, 0x88, 0xa6, 0x49, 0xa6, 0x08,
	0x21, 0x5e, 0xac, 0xcd, 0xfa, 0xc4, 0x5e, 0xb1, 0xde, 0xd9, 0xcc, 0xcc, 0x46, 0xda, 0x67, 0xde,
	0xf9, 0x09, 0xbc, 0xf0, 0x82, 0xf8, 0x09, 0xfc, 0x3a, 0xb4, 0x73, 0x59, 0xaf, 0x9d, 0x0d, 0x45,
	0xe2, 0x6d, 0xe6, 0x9c, 0xf3, 0x7d, 0xe7, 0x3a, 0x17, 0xe8, 0x0e, 0x03, 0xe1, 0xb3, 0x5b, 0xe4,
	0x69, 0x3f, 0xe6, 0x4c, 0x32, 0x9f, 0x85, 0xae, 0x5a, 0x90, 0x66, 0xae, 0xe9, 0x75, 0x46, 0x4c,
	0x88, 0x20, 0xee, 0x4f, 0x50, 0x08, 0x6f, 0x84, 0xda, 0xa0, 0xd7, 0x99, 0x88, 0xb8, 0x3f, 0x11,
	0xf1, 0xc0, 0x67, 0xd1, 0x75, 0x30, 0xd2, 0x52, 0xe7, 0x35, 0xb4, 0xdf, 0x05, 0xa3, 0x08, 0x87,
	0x14, 0x6f, 0x12, 0x14, 0x92, 0x74, 0xa1, 0x1e, 0x7b, 0x69, 0xc8, 0xbc, 0x61, 0xb7, 0xb2, 0x55,
	0xd9, 0x5e, 0xa1, 0x76, 0x4b, 0x1e, 0x42, 0x53, 0x04, 0xa3, 0xc8, 0x93, 0x09, 0xc7, 0x6e, 0x55,
	0xe9, 0xa6, 0x02, 0x87, 0x43, 0xdd, 0x52, 0xec, 0xc1, 0xaa, 0x97, 0xc8, 0x31, 0x46, 0x32, 0xf0,
	0x3d, 0x19, 0xb0, 0x48, 0x31, 0xb5, 0x76, 0x36, 0xdd, 0x3c, 0x46, 0xf7, 0x65, 0x22, 0xc7, 0xa7,
	0xd1, 0x35, 0xa3, 0x73, 0xa6, 0xe4, 0x29, 0xd4, 0x6f, 0x12, 0xe4, 0x01, 0x8a, 0x6e, 0x75, 0x6b,
	0x71, 0xbb, 0xb5, 0xb3, 0x5e, 0x40, 0x5d, 0x26, 0xc8, 0x53, 0x6a, 0x0d, 0x9c, 0x7d, 0x68, 0x50,
	0x14, 0x31, 0x8b, 0x04, 0x92, 0x67, 0x50, 0xe7, 0x28, 0x92, 0x50, 0x8a, 0x6e, 0x45, 0xe1, 0x1e,
	0xdc, 0xc1, 0x29, 0x35, 0xb5, 0x66, 0xce, 0x10, 0x1a, 0x36, 0x0a, 0xf2, 0x29, 0xac, 0xf9, 0x61,
	0x80, 0x91, 0x1c, 0x04, 0xc3, 0x2c, 0x18, 0x99, 0x9a, 0xec, 0x57, 0xb5, 0xf8, 0xd4, 0x48, 0x49,
	0x1f, 0x3a, 0xc6, 0x50, 0x86, 0x62, 0xe0, 0x23, 0x97, 0x83, 0xb1, 0x27, 0xc6, 0xa6, 0x1e, 0x1b,
	0x5a, 0xf7, 0x43, 0x28, 0x0e, 0x91, 0xcb, 0x13, 0x4f, 0x8c, 0x9d, 0xdf, 0xab, 0x50, 0x53, 0xee,
	0xb3, 0xca, 0xfa, 0x63, 0x2f, 0x8a, 0x30, 0x54, 0xdc, 0x4d, 0x6a, 0xb7, 0x64, 0x0f, 0x56, 0x74,
	0x53, 0x06, 0x59, 0x66, 0xa9, 0x22, 0x9b, 0x4d, 0xe0, 0x50, 0xa9, 0x15, 0xcf, 0xc9, 0x02, 0x6d,
	0xf9, 0xd3, 0x2d, 0xf9, 0x0e, 0x20, 0x46, 0xe4, 0x06, 0xba, 0xa8, 0xa0, 0x1f, 0x17, 0xa0, 0x17,
	0x88, 0xfc, 0x0c, 0x27, 0x57, 0xc8, 0xc5, 0x38, 0x88, 0x2d, 0x45, 0x33, 0xc3, 0x68, 0x82, 0xaf,
	0xa1, 0xe1, 0xfb, 0x06, 0xbe, 0xa4, 0xe0, 0x1f, 0x15, 0x3d, 0x8f, 0xbd, 0x20, 0xf2, 0xd9, 0x10,
	0x2d, 0xb2, 0xee, 0xfb, 0x1a, 0xb7, 0x0f, 0xad, 0x90, 0xf9, 0x5e, 0x38, 0xc8, 0xa8, 0x44, 0xb7,
	0x76, 0x07, 0xfa, 0x26, 0xd3, 0x5e, 0x58, 0x3f, 0x27, 0x0b, 0x14, 0x42, 0x2b, 0x11, 0x07, 0x75,
	0xa8, 0x29, 0x97, 0xce, 0xaf, 0x55, 0x68, 0x15, 0xfa, 0x43, 0xb6, 0xa1, 0x86, 0x9c, 0x33, 0x6e,
	0x86, 0xa6, 0xd8, 0xfe, 0xe3, 0x4c, 0x7e, 0xb2, 0x40, 0xb5, 0x01, 0xf9, 0x16, 0xda, 0xa6, 0x6c,
	0xba, 0xa5, 0xa6, 0x6e, 0x1f, 0xde, 0xa9, 0x9b, 0x66, 0x3e, 0x59, 0xa0, 0xa6, 0xcc, 0xc6, 0xd3,
	0x21, 0xac, 0xd8, 0xc4, 0x33, 0x06, 0x53, 0xbb, 0x47, 0xf7, 0x26, 0x9f, 0xd3, 0x80, 0x29, 0x01,
	0x45, 0x41, 0xf6, 0xa0, 0x3e, 0xd1, 0xd5, 0x35, 0xc5, 0x7b, 0x74, 0x6f, 0xed, 0x73, 0xbc, 0x45,
	0x1c, 0x34, 0x60, 0x59, 0x87, 0xee, 0xb4, 0xa1, 0x55, 0xe8, 0xb1, 0xf3, 0x57, 0x15, 0x56, 0x8a,
	0xb1, 0x93, 0xaf, 0x60, 0x69, 0x22, 0x62, 0x3b, 0xdb, 0x8f, 0xef, 0x49, 0xd1, 0x3d, 0x13, 0xb1,
	0x38, 0x8e, 0x24, 0x4f, 0xa9, 0x32, 0x27, 0x2f, 0xa1, 0xc1, 0xf8, 0x10, 0x79, 0x16, 0x9e, 0x3e,
	0x4e, 0x4f, 0xee, 0x83, 0x9e, 0x1b, 0x3b, 0x0d, 0xcf, 0x61, 0xbd, 0x33, 0x68, 0xe6, 0xac, 0x64,
	0x1d, 0x16, 0x7f, 0xc1, 0xd4, 0xcc, 0x6f, 0xb6, 0x24, 0x4f, 0xa1, 0x76, 0xeb, 0x85, 0x09, 0x9a,
	0xe2, 0x77, 0xdc, 0x89, 0x88, 0xdd, 0x57, 0xde, 0x15, 0x0f, 0xfc, 0xb3, 0x77, 0x17, 0xc6, 0x83,
	0x36, 0xd9, 0xad, 0x3e, 0xaf, 0xf4, 0x2e, 0xa1, 0x3d, 0xe3, 0xe9, 0xbf, 0x50, 0x16, 0x26, 0x20,
	0x1a, 0xc6, 0x2c, 0x88, 0xa4, 0x28, 0x50, 0x3a, 0xdf, 0xc3, 0x66, 0xc9, 0x90, 0x93, 0x2f, 0x61,
	0xf9, 0x3a, 0x08, 0x25, 0xda, 0x49, 0x7a, 0x58, 0xd6, 0xd8, 0xd3, 0x48, 0x22, 0x47, 0x21, 0xa9,
	0xb1, 0x75, 0xfe, 0xae, 0x40, 0xa7, 0xac, 0x6d, 0xe4, 0x12, 0x56, 0xd4, 0xa0, 0x0f, 0xae, 0xd2,
	0x01, 0xe3, 0x23, 0xd3, 0x89, 0xfe, 0x7b, 0xba, 0xed, 0xea, 0x69, 0x4f, 0xcf, 0xf9, 0x48, 0x17,
	0x56, 0x1d, 0x56, 0x2d, 0xe8, 0x9d, 0xc3, 0xda, 0x9c, 0xba, 0xa4, 0x1a, 0x9f, 0xcc, 0x56, 0x63,
	0x7d, 0xce, 0xe1, 0x4c, 0x25, 0xde, 0xc0, 0xea, 0xec, 0xc8, 0x92, 0x5d, 0x68, 0x06, 0x26, 0x45,
	0x3b, 0x3c, 0xff, 0x5e, 0x87, 0xa9, 0xb9, 0x73, 0x06, 0x1b, 0x77, 0xf4, 0xe4, 0x39, 0x80, 0x6f,
	0x85, 0x96, 0xb1, 0x5b, 0xc6, 0x78, 0xe8, 0x85, 0x21, 0x2d, 0xd8, 0x3a, 0x6f, 0xa1, 0x3d, 0xa3,
	0x24, 0x04, 0x96, 0x22, 0x6f, 0x82, 0x26, 0x59, 0xb5, 0x26, 0x9f, 0xc1, 0xba, 0xcf, 0xc2, 0x10,
	0xfd, 0xec, 0x31, 0x18, 0x64, 0x22, 0x3d, 0xb8, 0x4d, 0xba, 0x36, 0x95, 0xbf, 0xcd, 0xc4, 0x0e,
	0x85, 0x4e, 0xd9, 0xf9, 0x24, 0xbb, 0x50, 0xf7, 0x59, 0x24, 0x31, 0x92, 0x26, 0xbc, 0xad, 0xd9,
	0x01, 0x62, 0x5c, 0xe0, 0x04, 0x23, 0x79, 0x84, 0xc2, 0xe7, 0x41, 0x2c, 0x19, 0xa7, 0x16, 0xe0,
	0xac, 0xc3, 0xea, 0xec, 0xad, 0xe5, 0xfc, 0x51, 0x85, 0x0f, 0x4a, 0x41, 0xd9, 0x7b, 0x98, 0x67,
	0x67, 0x72, 0x98, 0x0a, 0xc8, 0x08, 0x36, 0x51, 0xc3, 0xf4, 0xc8, 0x8c, 0x38, 0x4b, 0x62, 0x7b,
	0x08, 0xbf, 0x79, 0x5f, 0x44, 0x56, 0x9a, 0xcd, 0xc6, 0x6b, 0x85, 0xd4, 0xd3, 0xb3, 0x81, 0xf3,
	0x72, 0xf2, 0x39, 0xd4, 0x43, 0x2f, 0x65, 0x89, 0xcc, 0x2e, 0xb0, 0x8c, 0x7c, 0xa3, 0x78, 0x05,
	0x2b, 0x0d, 0xb5, 0x16, 0xbd, 0x1f, 0xe1, 0x41, 0x39, 0xf3, 0xff, 0x1c, 0xbc, 0x3f, 0x2b, 0xb0,
	0xac, 0x7d, 0x91, 0x9f, 0x60, 0xf3, 0x26, 0xf1, 0xb2, 0xd7, 0x32, 0xc0, 0x69, 0xe6, 0xa6, 0x15,
	0xdb, 0x77, 0x62, 0x73, 0x2f, 0x73, 0x63, 0x13, 0x90, 0xc9, 0xf4, 0x66, 0x5e, 0xde, 0x3b, 0x82,
	0x07, 0xe5, 0xc6, 0x25, 0xc1, 0x77, 0x8a, 0xc1, 0xb7, 0x8b, 0xa1, 0xba, 0x50, 0x53, 0xe1, 0x93,
	0x27, 0x50, 0xd3, 0x2f, 0x97, 0x0e, 0x6d, 0x6d, 0x2e, 0x3f, 0xaa, 0xb5, 0xce, 0x6f, 0x15, 0x58,
	0xca, 0xf6, 0xa4, 0x0f, 0x20, 0xa4, 0x27, 0x71, 0x10, 0x44, 0xd7, 0x2c, 0x7f, 0x9d, 0xf4, 0x5f,
	0xcb, 0x3d, 0x8e, 0x6e, 0x31, 0x64, 0x31, 0xd2, 0xa6, 0xb2, 0x51, 0x9f, 0x8a, 0x17, 0xb0, 0x36,
	0xc9, 0xaf, 0x03, 0x8d, 0xaa, 0xde, 0x83, 0x5a, 0x9d, 0x1a, 0x2a, 0x68, 0x0f, 0x1a, 0xf9, 0x47,
	0x64, 0x51, 0x7d, 0x2d, 0xf2, 0xbd, 0xf3, 0x18, 0x6a, 0xea, 0x21, 0x54, 0x1f, 0x8a, 0x7c, 0xd0,
	0xf5, 0x87, 0xc2, 0x8c, 0xf1, 0x3e, 0x34, 0xf3, 0x9b, 0x92, 0xf4, 0xa1, 0x81, 0x66, 0x63, 0x52,
	0xdd, 0x2c, 0xb9, 0x51, 0x69, 0x6e, 0xe4, 0xec, 0x40, 0xc3, 0x4a, 0xb3, 0x33, 0x3a, 0x66, 0xc2,
	0x3a, 0x50, 0xeb, 0x4c, 0x16, 0x33, 0x2e, 0x4d, 0x69, 0xd5, 0x7a, 0xe7, 0x15, 0x34, 0x8f, 0x2c,
	0x27, 0x79, 0x01, 0x0d, 0xbb, 0x21, 0xc5, 0xbb, 0x61, 0xe6, 0xa7, 0xd9, 0x2b, 0x46, 0x61, 0xbf,
	0x71, 0x07, 0xcf, 0x7e, 0x76, 0x47, 0x81, 0x1c, 0x27, 0x57, 0xae, 0xcf, 0x26, 0xfd, 0x71, 0x1a,
	0x23, 0x0f, 0x71, 0x38, 0x42, 0xde, 0xbf, 0x56, 0x6f, 0x8a, 0xfe, 0xf6, 0x8a, 0x7e, 0x0e, 0xbd,
	0x5a, 0x56, 0x92, 0x2f, 0xfe, 0x09, 0x00, 0x00, 0xff, 0xff, 0x81, 0x4a, 0xe4, 0xb3, 0x1b, 0x0b,
	0x00, 0x00,
}


var _ context.Context
var _ grpc.ClientConn



const _ = grpc.SupportPackageIsVersion4




type DiscoveryClient interface {
	
	Discover(ctx context.Context, in *SignedRequest, opts ...grpc.CallOption) (*Response, error)
}

type discoveryClient struct {
	cc *grpc.ClientConn
}

func NewDiscoveryClient(cc *grpc.ClientConn) DiscoveryClient {
	return &discoveryClient{cc}
}

func (c *discoveryClient) Discover(ctx context.Context, in *SignedRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/discovery.Discovery/Discover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}


type DiscoveryServer interface {
	
	Discover(context.Context, *SignedRequest) (*Response, error)
}

func RegisterDiscoveryServer(s *grpc.Server, srv DiscoveryServer) {
	s.RegisterService(&_Discovery_serviceDesc, srv)
}

func _Discovery_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoveryServer).Discover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discovery.Discovery/Discover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoveryServer).Discover(ctx, req.(*SignedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Discovery_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discovery.Discovery",
	HandlerType: (*DiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Discover",
			Handler:    _Discovery_Discover_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "discovery/protocol.proto",
}
