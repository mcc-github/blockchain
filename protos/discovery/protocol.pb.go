


package discovery 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import gossip "github.com/mcc-github/blockchain/protos/gossip"
import msp "github.com/mcc-github/blockchain/protos/msp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 






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
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{0}
}
func (m *SignedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedRequest.Unmarshal(m, b)
}
func (m *SignedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedRequest.Marshal(b, m, deterministic)
}
func (dst *SignedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedRequest.Merge(dst, src)
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
	
	
	Authentication *AuthInfo `protobuf:"bytes,1,opt,name=authentication" json:"authentication,omitempty"`
	
	Queries              []*Query `protobuf:"bytes,2,rep,name=queries" json:"queries,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{1}
}
func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (dst *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(dst, src)
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
	
	Results              []*QueryResult `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{2}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
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
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{3}
}
func (m *AuthInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthInfo.Unmarshal(m, b)
}
func (m *AuthInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthInfo.Marshal(b, m, deterministic)
}
func (dst *AuthInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthInfo.Merge(dst, src)
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
	Channel string `protobuf:"bytes,1,opt,name=channel" json:"channel,omitempty"`
	
	
	
	
	
	Query                isQuery_Query `protobuf_oneof:"query"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{4}
}
func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (dst *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(dst, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

type isQuery_Query interface {
	isQuery_Query()
}

type Query_ConfigQuery struct {
	ConfigQuery *ConfigQuery `protobuf:"bytes,2,opt,name=config_query,json=configQuery,oneof"`
}
type Query_PeerQuery struct {
	PeerQuery *PeerMembershipQuery `protobuf:"bytes,3,opt,name=peer_query,json=peerQuery,oneof"`
}
type Query_CcQuery struct {
	CcQuery *ChaincodeQuery `protobuf:"bytes,4,opt,name=cc_query,json=ccQuery,oneof"`
}
type Query_LocalPeers struct {
	LocalPeers *LocalPeerQuery `protobuf:"bytes,5,opt,name=local_peers,json=localPeers,oneof"`
}

func (*Query_ConfigQuery) isQuery_Query() {}
func (*Query_PeerQuery) isQuery_Query()   {}
func (*Query_CcQuery) isQuery_Query()     {}
func (*Query_LocalPeers) isQuery_Query()  {}

func (m *Query) GetQuery() isQuery_Query {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *Query) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
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


func (*Query) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Query_OneofMarshaler, _Query_OneofUnmarshaler, _Query_OneofSizer, []interface{}{
		(*Query_ConfigQuery)(nil),
		(*Query_PeerQuery)(nil),
		(*Query_CcQuery)(nil),
		(*Query_LocalPeers)(nil),
	}
}

func _Query_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Query)
	
	switch x := m.Query.(type) {
	case *Query_ConfigQuery:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ConfigQuery); err != nil {
			return err
		}
	case *Query_PeerQuery:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PeerQuery); err != nil {
			return err
		}
	case *Query_CcQuery:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CcQuery); err != nil {
			return err
		}
	case *Query_LocalPeers:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LocalPeers); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Query.Query has unexpected type %T", x)
	}
	return nil
}

func _Query_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Query)
	switch tag {
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ConfigQuery)
		err := b.DecodeMessage(msg)
		m.Query = &Query_ConfigQuery{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PeerMembershipQuery)
		err := b.DecodeMessage(msg)
		m.Query = &Query_PeerQuery{msg}
		return true, err
	case 4: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeQuery)
		err := b.DecodeMessage(msg)
		m.Query = &Query_CcQuery{msg}
		return true, err
	case 5: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LocalPeerQuery)
		err := b.DecodeMessage(msg)
		m.Query = &Query_LocalPeers{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Query_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Query)
	
	switch x := m.Query.(type) {
	case *Query_ConfigQuery:
		s := proto.Size(x.ConfigQuery)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Query_PeerQuery:
		s := proto.Size(x.PeerQuery)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Query_CcQuery:
		s := proto.Size(x.CcQuery)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Query_LocalPeers:
		s := proto.Size(x.LocalPeers)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
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
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{5}
}
func (m *QueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResult.Unmarshal(m, b)
}
func (m *QueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResult.Marshal(b, m, deterministic)
}
func (dst *QueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResult.Merge(dst, src)
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
	Error *Error `protobuf:"bytes,1,opt,name=error,oneof"`
}
type QueryResult_ConfigResult struct {
	ConfigResult *ConfigResult `protobuf:"bytes,2,opt,name=config_result,json=configResult,oneof"`
}
type QueryResult_CcQueryRes struct {
	CcQueryRes *ChaincodeQueryResult `protobuf:"bytes,3,opt,name=cc_query_res,json=ccQueryRes,oneof"`
}
type QueryResult_Members struct {
	Members *PeerMembershipResult `protobuf:"bytes,4,opt,name=members,oneof"`
}

func (*QueryResult_Error) isQueryResult_Result()        {}
func (*QueryResult_ConfigResult) isQueryResult_Result() {}
func (*QueryResult_CcQueryRes) isQueryResult_Result()   {}
func (*QueryResult_Members) isQueryResult_Result()      {}

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


func (*QueryResult) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _QueryResult_OneofMarshaler, _QueryResult_OneofUnmarshaler, _QueryResult_OneofSizer, []interface{}{
		(*QueryResult_Error)(nil),
		(*QueryResult_ConfigResult)(nil),
		(*QueryResult_CcQueryRes)(nil),
		(*QueryResult_Members)(nil),
	}
}

func _QueryResult_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*QueryResult)
	
	switch x := m.Result.(type) {
	case *QueryResult_Error:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Error); err != nil {
			return err
		}
	case *QueryResult_ConfigResult:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ConfigResult); err != nil {
			return err
		}
	case *QueryResult_CcQueryRes:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CcQueryRes); err != nil {
			return err
		}
	case *QueryResult_Members:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Members); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("QueryResult.Result has unexpected type %T", x)
	}
	return nil
}

func _QueryResult_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*QueryResult)
	switch tag {
	case 1: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Error)
		err := b.DecodeMessage(msg)
		m.Result = &QueryResult_Error{msg}
		return true, err
	case 2: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ConfigResult)
		err := b.DecodeMessage(msg)
		m.Result = &QueryResult_ConfigResult{msg}
		return true, err
	case 3: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeQueryResult)
		err := b.DecodeMessage(msg)
		m.Result = &QueryResult_CcQueryRes{msg}
		return true, err
	case 4: 
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PeerMembershipResult)
		err := b.DecodeMessage(msg)
		m.Result = &QueryResult_Members{msg}
		return true, err
	default:
		return false, nil
	}
}

func _QueryResult_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*QueryResult)
	
	switch x := m.Result.(type) {
	case *QueryResult_Error:
		s := proto.Size(x.Error)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *QueryResult_ConfigResult:
		s := proto.Size(x.ConfigResult)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *QueryResult_CcQueryRes:
		s := proto.Size(x.CcQueryRes)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case *QueryResult_Members:
		s := proto.Size(x.Members)
		n += 1 
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
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
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{6}
}
func (m *ConfigQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigQuery.Unmarshal(m, b)
}
func (m *ConfigQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigQuery.Marshal(b, m, deterministic)
}
func (dst *ConfigQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigQuery.Merge(dst, src)
}
func (m *ConfigQuery) XXX_Size() int {
	return xxx_messageInfo_ConfigQuery.Size(m)
}
func (m *ConfigQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigQuery.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigQuery proto.InternalMessageInfo

type ConfigResult struct {
	
	Msps map[string]*msp.FabricMSPConfig `protobuf:"bytes,1,rep,name=msps" json:"msps,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	
	Orderers             map[string]*Endpoints `protobuf:"bytes,2,rep,name=orderers" json:"orderers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ConfigResult) Reset()         { *m = ConfigResult{} }
func (m *ConfigResult) String() string { return proto.CompactTextString(m) }
func (*ConfigResult) ProtoMessage()    {}
func (*ConfigResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{7}
}
func (m *ConfigResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigResult.Unmarshal(m, b)
}
func (m *ConfigResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigResult.Marshal(b, m, deterministic)
}
func (dst *ConfigResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigResult.Merge(dst, src)
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
	Filter               *ChaincodeInterest `protobuf:"bytes,1,opt,name=filter" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *PeerMembershipQuery) Reset()         { *m = PeerMembershipQuery{} }
func (m *PeerMembershipQuery) String() string { return proto.CompactTextString(m) }
func (*PeerMembershipQuery) ProtoMessage()    {}
func (*PeerMembershipQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{8}
}
func (m *PeerMembershipQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerMembershipQuery.Unmarshal(m, b)
}
func (m *PeerMembershipQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerMembershipQuery.Marshal(b, m, deterministic)
}
func (dst *PeerMembershipQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerMembershipQuery.Merge(dst, src)
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
	PeersByOrg           map[string]*Peers `protobuf:"bytes,1,rep,name=peers_by_org,json=peersByOrg" json:"peers_by_org,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *PeerMembershipResult) Reset()         { *m = PeerMembershipResult{} }
func (m *PeerMembershipResult) String() string { return proto.CompactTextString(m) }
func (*PeerMembershipResult) ProtoMessage()    {}
func (*PeerMembershipResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{9}
}
func (m *PeerMembershipResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerMembershipResult.Unmarshal(m, b)
}
func (m *PeerMembershipResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerMembershipResult.Marshal(b, m, deterministic)
}
func (dst *PeerMembershipResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerMembershipResult.Merge(dst, src)
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
	Interests            []*ChaincodeInterest `protobuf:"bytes,1,rep,name=interests" json:"interests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ChaincodeQuery) Reset()         { *m = ChaincodeQuery{} }
func (m *ChaincodeQuery) String() string { return proto.CompactTextString(m) }
func (*ChaincodeQuery) ProtoMessage()    {}
func (*ChaincodeQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{10}
}
func (m *ChaincodeQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeQuery.Unmarshal(m, b)
}
func (m *ChaincodeQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeQuery.Marshal(b, m, deterministic)
}
func (dst *ChaincodeQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeQuery.Merge(dst, src)
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
	Chaincodes           []*ChaincodeCall `protobuf:"bytes,1,rep,name=chaincodes" json:"chaincodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ChaincodeInterest) Reset()         { *m = ChaincodeInterest{} }
func (m *ChaincodeInterest) String() string { return proto.CompactTextString(m) }
func (*ChaincodeInterest) ProtoMessage()    {}
func (*ChaincodeInterest) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{11}
}
func (m *ChaincodeInterest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeInterest.Unmarshal(m, b)
}
func (m *ChaincodeInterest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeInterest.Marshal(b, m, deterministic)
}
func (dst *ChaincodeInterest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeInterest.Merge(dst, src)
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
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	CollectionNames      []string `protobuf:"bytes,2,rep,name=collection_names,json=collectionNames" json:"collection_names,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeCall) Reset()         { *m = ChaincodeCall{} }
func (m *ChaincodeCall) String() string { return proto.CompactTextString(m) }
func (*ChaincodeCall) ProtoMessage()    {}
func (*ChaincodeCall) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{12}
}
func (m *ChaincodeCall) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeCall.Unmarshal(m, b)
}
func (m *ChaincodeCall) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeCall.Marshal(b, m, deterministic)
}
func (dst *ChaincodeCall) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeCall.Merge(dst, src)
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
	Content              []*EndorsementDescriptor `protobuf:"bytes,1,rep,name=content" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ChaincodeQueryResult) Reset()         { *m = ChaincodeQueryResult{} }
func (m *ChaincodeQueryResult) String() string { return proto.CompactTextString(m) }
func (*ChaincodeQueryResult) ProtoMessage()    {}
func (*ChaincodeQueryResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{13}
}
func (m *ChaincodeQueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeQueryResult.Unmarshal(m, b)
}
func (m *ChaincodeQueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeQueryResult.Marshal(b, m, deterministic)
}
func (dst *ChaincodeQueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeQueryResult.Merge(dst, src)
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
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{14}
}
func (m *LocalPeerQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LocalPeerQuery.Unmarshal(m, b)
}
func (m *LocalPeerQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LocalPeerQuery.Marshal(b, m, deterministic)
}
func (dst *LocalPeerQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LocalPeerQuery.Merge(dst, src)
}
func (m *LocalPeerQuery) XXX_Size() int {
	return xxx_messageInfo_LocalPeerQuery.Size(m)
}
func (m *LocalPeerQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_LocalPeerQuery.DiscardUnknown(m)
}

var xxx_messageInfo_LocalPeerQuery proto.InternalMessageInfo













type EndorsementDescriptor struct {
	Chaincode string `protobuf:"bytes,1,opt,name=chaincode" json:"chaincode,omitempty"`
	
	EndorsersByGroups map[string]*Peers `protobuf:"bytes,2,rep,name=endorsers_by_groups,json=endorsersByGroups" json:"endorsers_by_groups,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	
	
	
	Layouts              []*Layout `protobuf:"bytes,3,rep,name=layouts" json:"layouts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *EndorsementDescriptor) Reset()         { *m = EndorsementDescriptor{} }
func (m *EndorsementDescriptor) String() string { return proto.CompactTextString(m) }
func (*EndorsementDescriptor) ProtoMessage()    {}
func (*EndorsementDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{15}
}
func (m *EndorsementDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EndorsementDescriptor.Unmarshal(m, b)
}
func (m *EndorsementDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EndorsementDescriptor.Marshal(b, m, deterministic)
}
func (dst *EndorsementDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EndorsementDescriptor.Merge(dst, src)
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
	
	
	QuantitiesByGroup    map[string]uint32 `protobuf:"bytes,1,rep,name=quantities_by_group,json=quantitiesByGroup" json:"quantities_by_group,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Layout) Reset()         { *m = Layout{} }
func (m *Layout) String() string { return proto.CompactTextString(m) }
func (*Layout) ProtoMessage()    {}
func (*Layout) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{16}
}
func (m *Layout) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Layout.Unmarshal(m, b)
}
func (m *Layout) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Layout.Marshal(b, m, deterministic)
}
func (dst *Layout) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Layout.Merge(dst, src)
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
	Peers                []*Peer  `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peers) Reset()         { *m = Peers{} }
func (m *Peers) String() string { return proto.CompactTextString(m) }
func (*Peers) ProtoMessage()    {}
func (*Peers) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{17}
}
func (m *Peers) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peers.Unmarshal(m, b)
}
func (m *Peers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peers.Marshal(b, m, deterministic)
}
func (dst *Peers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peers.Merge(dst, src)
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
	
	StateInfo *gossip.Envelope `protobuf:"bytes,1,opt,name=state_info,json=stateInfo" json:"state_info,omitempty"`
	
	MembershipInfo *gossip.Envelope `protobuf:"bytes,2,opt,name=membership_info,json=membershipInfo" json:"membership_info,omitempty"`
	
	Identity             []byte   `protobuf:"bytes,3,opt,name=identity,proto3" json:"identity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peer) Reset()         { *m = Peer{} }
func (m *Peer) String() string { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()    {}
func (*Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{18}
}
func (m *Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer.Unmarshal(m, b)
}
func (m *Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer.Marshal(b, m, deterministic)
}
func (dst *Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer.Merge(dst, src)
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
	Content              string   `protobuf:"bytes,1,opt,name=content" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{19}
}
func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (dst *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(dst, src)
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
	Endpoint             []*Endpoint `protobuf:"bytes,1,rep,name=endpoint" json:"endpoint,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Endpoints) Reset()         { *m = Endpoints{} }
func (m *Endpoints) String() string { return proto.CompactTextString(m) }
func (*Endpoints) ProtoMessage()    {}
func (*Endpoints) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{20}
}
func (m *Endpoints) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoints.Unmarshal(m, b)
}
func (m *Endpoints) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoints.Marshal(b, m, deterministic)
}
func (dst *Endpoints) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoints.Merge(dst, src)
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
	Host                 string   `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}
func (*Endpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_b2f93a2b7b5bdad4, []int{21}
}
func (m *Endpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoint.Unmarshal(m, b)
}
func (m *Endpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoint.Marshal(b, m, deterministic)
}
func (dst *Endpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoint.Merge(dst, src)
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
	err := grpc.Invoke(ctx, "/discovery.Discovery/Discover", in, out, c.cc, opts...)
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

func init() { proto.RegisterFile("discovery/protocol.proto", fileDescriptor_protocol_b2f93a2b7b5bdad4) }

var fileDescriptor_protocol_b2f93a2b7b5bdad4 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0x5b, 0x6f, 0xe3, 0x44,
	0x14, 0x6e, 0xd2, 0xa6, 0x49, 0x4e, 0x9a, 0x5e, 0xa6, 0x61, 0x09, 0xd1, 0x8a, 0xdd, 0xb5, 0xb4,
	0x50, 0x16, 0xc9, 0x59, 0x95, 0xdb, 0xd2, 0x56, 0xa0, 0xed, 0x85, 0x4d, 0xc5, 0x76, 0xdb, 0xce,
	0x22, 0x84, 0x78, 0x89, 0x5c, 0xe7, 0x34, 0xb1, 0x70, 0x3c, 0xee, 0xcc, 0xb8, 0x92, 0x9f, 0x79,
	0xe7, 0x27, 0xf0, 0xc2, 0x0b, 0xe2, 0x27, 0xf0, 0xeb, 0x90, 0xe7, 0xe2, 0x38, 0xa9, 0xcb, 0x22,
	0xf1, 0xe6, 0x39, 0xe7, 0x7c, 0xe7, 0xf2, 0x9d, 0xe3, 0x99, 0x03, 0xdd, 0x51, 0x20, 0x7c, 0x76,
	0x8b, 0x3c, 0xed, 0xc7, 0x9c, 0x49, 0xe6, 0xb3, 0xd0, 0x55, 0x1f, 0xa4, 0x99, 0x6b, 0x7a, 0x9d,
	0x31, 0x13, 0x22, 0x88, 0xfb, 0x53, 0x14, 0xc2, 0x1b, 0xa3, 0x36, 0xe8, 0x75, 0xa6, 0x22, 0xee,
	0x4f, 0x45, 0x3c, 0xf4, 0x59, 0x74, 0x1d, 0x8c, 0x8b, 0xd2, 0x60, 0x84, 0x91, 0x0c, 0x64, 0x80,
	0x42, 0x4b, 0x9d, 0x57, 0xd0, 0x7e, 0x1b, 0x8c, 0x23, 0x1c, 0x51, 0xbc, 0x49, 0x50, 0x48, 0xd2,
	0x85, 0x7a, 0xec, 0xa5, 0x21, 0xf3, 0x46, 0xdd, 0xca, 0xe3, 0xca, 0xce, 0x1a, 0xb5, 0x47, 0xf2,
	0x10, 0x9a, 0x22, 0x18, 0x47, 0x9e, 0x4c, 0x38, 0x76, 0xab, 0x4a, 0x37, 0x13, 0x38, 0x1c, 0xea,
	0xd6, 0xc5, 0x3e, 0xac, 0x7b, 0x89, 0x9c, 0x64, 0x91, 0x7c, 0x4f, 0x06, 0x2c, 0x52, 0x9e, 0x5a,
	0xbb, 0xdb, 0x6e, 0x9e, 0xb9, 0xfb, 0x32, 0x91, 0x93, 0xd3, 0xe8, 0x9a, 0xd1, 0x05, 0x53, 0xf2,
	0x0c, 0xea, 0x37, 0x09, 0xf2, 0x00, 0x45, 0xb7, 0xfa, 0x78, 0x79, 0xa7, 0xb5, 0xbb, 0x59, 0x40,
	0x5d, 0x26, 0xc8, 0x53, 0x6a, 0x0d, 0x9c, 0x03, 0x68, 0x50, 0x14, 0x31, 0x8b, 0x04, 0x92, 0xe7,
	0x50, 0xe7, 0x28, 0x92, 0x50, 0x8a, 0x6e, 0x45, 0xe1, 0x1e, 0xdc, 0xc1, 0x29, 0x35, 0xb5, 0x66,
	0xce, 0x08, 0x1a, 0x36, 0x0b, 0xf2, 0x31, 0x6c, 0xf8, 0x61, 0x80, 0x91, 0x1c, 0x1a, 0x86, 0x52,
	0x53, 0xfd, 0xba, 0x16, 0x9f, 0x1a, 0x29, 0xe9, 0x43, 0xc7, 0x18, 0xca, 0x50, 0x0c, 0x7d, 0xe4,
	0x72, 0x38, 0xf1, 0xc4, 0xc4, 0xf0, 0xb1, 0xa5, 0x75, 0x3f, 0x84, 0xe2, 0x08, 0xb9, 0x1c, 0x78,
	0x62, 0xe2, 0xfc, 0x5e, 0x85, 0x9a, 0x0a, 0x9f, 0x31, 0xeb, 0x4f, 0xbc, 0x28, 0xc2, 0x50, 0xf9,
	0x6e, 0x52, 0x7b, 0x24, 0xfb, 0xb0, 0xa6, 0x5b, 0x35, 0xcc, 0x2a, 0x4b, 0x95, 0xb3, 0xf9, 0x02,
	0x8e, 0x94, 0x5a, 0xf9, 0x19, 0x2c, 0xd1, 0x96, 0x3f, 0x3b, 0x92, 0x6f, 0x01, 0x62, 0x44, 0x6e,
	0xa0, 0xcb, 0x0a, 0xfa, 0x61, 0x01, 0x7a, 0x81, 0xc8, 0xcf, 0x70, 0x7a, 0x85, 0x5c, 0x4c, 0x82,
	0xd8, 0xba, 0x68, 0x66, 0x18, 0xed, 0xe0, 0x4b, 0x68, 0xf8, 0xbe, 0x81, 0xaf, 0x28, 0xf8, 0x07,
	0xc5, 0xc8, 0x13, 0x2f, 0x88, 0x7c, 0x36, 0x42, 0x8b, 0xac, 0xfb, 0xbe, 0xc6, 0x1d, 0x40, 0x2b,
	0x64, 0xbe, 0x17, 0x0e, 0x33, 0x57, 0xa2, 0x5b, 0xbb, 0x03, 0x7d, 0x9d, 0x69, 0x2f, 0x6c, 0x9c,
	0xc1, 0x12, 0x85, 0xd0, 0x4a, 0xc4, 0x61, 0x1d, 0x6a, 0x2a, 0xa4, 0xf3, 0x6b, 0x15, 0x5a, 0x85,
	0xfe, 0x90, 0x1d, 0xa8, 0x21, 0xe7, 0x8c, 0x9b, 0xa1, 0x29, 0xb6, 0xff, 0x24, 0x93, 0x0f, 0x96,
	0xa8, 0x36, 0x20, 0xdf, 0x40, 0xdb, 0xd0, 0xa6, 0x5b, 0x6a, 0x78, 0x7b, 0xff, 0x0e, 0x6f, 0xda,
	0xf3, 0x60, 0x89, 0x1a, 0x9a, 0x4d, 0xa4, 0x23, 0x58, 0xb3, 0x85, 0x67, 0x1e, 0x0c, 0x77, 0x8f,
	0xee, 0x2d, 0x3e, 0x77, 0x03, 0x86, 0x02, 0x8a, 0x82, 0xec, 0x43, 0x7d, 0xaa, 0xd9, 0x35, 0xe4,
	0x3d, 0xba, 0x97, 0xfb, 0x1c, 0x6f, 0x11, 0x87, 0x0d, 0x58, 0xd5, 0xa9, 0x3b, 0x6d, 0x68, 0x15,
	0x7a, 0xec, 0xfc, 0x55, 0x85, 0xb5, 0x62, 0xee, 0xe4, 0x0b, 0x58, 0x99, 0x8a, 0xd8, 0xce, 0xf6,
	0x93, 0x7b, 0x4a, 0x74, 0xcf, 0x44, 0x2c, 0x4e, 0x22, 0xc9, 0x53, 0xaa, 0xcc, 0xc9, 0x4b, 0x68,
	0x30, 0x3e, 0x42, 0x9e, 0xa5, 0xa7, 0x7f, 0xa7, 0xa7, 0xf7, 0x41, 0xcf, 0x8d, 0x9d, 0x86, 0xe7,
	0xb0, 0xde, 0x19, 0x34, 0x73, 0xaf, 0x64, 0x13, 0x96, 0x7f, 0xc1, 0xd4, 0xcc, 0x6f, 0xf6, 0x49,
	0x9e, 0x41, 0xed, 0xd6, 0x0b, 0x13, 0x34, 0xe4, 0x77, 0xdc, 0xa9, 0x88, 0xdd, 0xef, 0xbc, 0x2b,
	0x1e, 0xf8, 0x67, 0x6f, 0x2f, 0x4c, 0x04, 0x6d, 0xb2, 0x57, 0x7d, 0x51, 0xe9, 0x5d, 0x42, 0x7b,
	0x2e, 0xd2, 0x7f, 0x71, 0x59, 0x98, 0x80, 0x68, 0x14, 0xb3, 0x20, 0x92, 0xa2, 0xe0, 0xd2, 0xf9,
	0x1e, 0xb6, 0x4b, 0x86, 0x9c, 0x7c, 0x0e, 0xab, 0xd7, 0x41, 0x28, 0xd1, 0x4e, 0xd2, 0xc3, 0xb2,
	0xc6, 0x9e, 0x46, 0x12, 0x39, 0x0a, 0x49, 0x8d, 0xad, 0xf3, 0x77, 0x05, 0x3a, 0x65, 0x6d, 0x23,
	0x97, 0xb0, 0xa6, 0x06, 0x7d, 0x78, 0x95, 0x0e, 0x19, 0x1f, 0x9b, 0x4e, 0xf4, 0xdf, 0xd1, 0x6d,
	0x57, 0x4f, 0x7b, 0x7a, 0xce, 0xc7, 0x9a, 0x58, 0xf5, 0xb3, 0x6a, 0x41, 0xef, 0x1c, 0x36, 0x16,
	0xd4, 0x25, 0x6c, 0x7c, 0x34, 0xcf, 0xc6, 0xe6, 0x42, 0xc0, 0x39, 0x26, 0x5e, 0xc3, 0xfa, 0xfc,
	0xc8, 0x92, 0x3d, 0x68, 0x06, 0xa6, 0x44, 0x3b, 0x3c, 0xff, 0xce, 0xc3, 0xcc, 0xdc, 0x39, 0x83,
	0xad, 0x3b, 0x7a, 0xf2, 0x02, 0xc0, 0xb7, 0x42, 0xeb, 0xb1, 0x5b, 0xe6, 0xf1, 0xc8, 0x0b, 0x43,
	0x5a, 0xb0, 0x75, 0xde, 0x40, 0x7b, 0x4e, 0x49, 0x08, 0xac, 0x44, 0xde, 0x14, 0x4d, 0xb1, 0xea,
	0x9b, 0x7c, 0x02, 0x9b, 0x3e, 0x0b, 0x43, 0xf4, 0xb3, 0xc7, 0x60, 0x98, 0x89, 0xf4, 0xe0, 0x36,
	0xe9, 0xc6, 0x4c, 0xfe, 0x26, 0x13, 0x3b, 0x14, 0x3a, 0x65, 0xff, 0x27, 0xd9, 0x83, 0xba, 0xcf,
	0x22, 0x89, 0x91, 0x34, 0xe9, 0x3d, 0x9e, 0x1f, 0x20, 0xc6, 0x05, 0x4e, 0x31, 0x92, 0xc7, 0x28,
	0x7c, 0x1e, 0xc4, 0x92, 0x71, 0x6a, 0x01, 0xce, 0x26, 0xac, 0xcf, 0xdf, 0x5a, 0xce, 0x1f, 0x55,
	0x78, 0xaf, 0x14, 0x94, 0xbd, 0x87, 0x79, 0x75, 0xa6, 0x86, 0x99, 0x80, 0x8c, 0x61, 0x1b, 0x35,
	0x4c, 0x8f, 0xcc, 0x98, 0xb3, 0x24, 0xb6, 0x3f, 0xe1, 0x57, 0xef, 0xca, 0xc8, 0x4a, 0xb3, 0xd9,
	0x78, 0xa5, 0x90, 0x7a, 0x7a, 0xb6, 0x70, 0x51, 0x4e, 0x3e, 0x85, 0x7a, 0xe8, 0xa5, 0x2c, 0x91,
	0xd9, 0x05, 0x96, 0x39, 0xdf, 0x2a, 0x5e, 0xc1, 0x4a, 0x43, 0xad, 0x45, 0xef, 0x47, 0x78, 0x50,
	0xee, 0xf9, 0x7f, 0x0e, 0xde, 0x9f, 0x15, 0x58, 0xd5, 0xb1, 0xc8, 0x4f, 0xb0, 0x7d, 0x93, 0x78,
	0x66, 0xcb, 0xc8, 0x2b, 0x37, 0xad, 0xd8, 0xb9, 0x93, 0x9b, 0x7b, 0x99, 0x1b, 0x9b, 0x84, 0x4c,
	0xa5, 0x37, 0x8b, 0xf2, 0xde, 0x31, 0x3c, 0x28, 0x37, 0x2e, 0x49, 0xbe, 0x53, 0x4c, 0xbe, 0x5d,
	0x4c, 0xd5, 0x85, 0x9a, 0x4a, 0x9f, 0x3c, 0x85, 0x9a, 0x7e, 0xb9, 0x74, 0x6a, 0x1b, 0x0b, 0xf5,
	0x51, 0xad, 0x75, 0x7e, 0xab, 0xc0, 0x4a, 0x76, 0x26, 0x7d, 0x00, 0x21, 0x3d, 0x89, 0xc3, 0x20,
	0xba, 0x66, 0xf9, 0xeb, 0xa4, 0x37, 0x30, 0xf7, 0x24, 0xba, 0xc5, 0x90, 0xc5, 0x48, 0x9b, 0xca,
	0x46, 0x2d, 0x15, 0x5f, 0xc3, 0xc6, 0x34, 0xbf, 0x0e, 0x34, 0xaa, 0x7a, 0x0f, 0x6a, 0x7d, 0x66,
	0xa8, 0xa0, 0x3d, 0x68, 0xe4, 0x8b, 0xc8, 0xb2, 0x5a, 0x2d, 0xf2, 0xb3, 0xf3, 0x04, 0x6a, 0xea,
	0x21, 0x54, 0x0b, 0x45, 0x3e, 0xe8, 0x7a, 0xa1, 0x30, 0x63, 0x7c, 0x00, 0xcd, 0xfc, 0xa6, 0x24,
	0x7d, 0x68, 0xa0, 0x39, 0x98, 0x52, 0xb7, 0x4b, 0x6e, 0x54, 0x9a, 0x1b, 0x39, 0xbb, 0xd0, 0xb0,
	0xd2, 0xec, 0x1f, 0x9d, 0x30, 0x61, 0x03, 0xa8, 0xef, 0x4c, 0x16, 0x33, 0x2e, 0x0d, 0xb5, 0xea,
	0x7b, 0x77, 0x00, 0xcd, 0x63, 0xeb, 0x93, 0xec, 0x43, 0xc3, 0x1e, 0x48, 0xf1, 0x6e, 0x98, 0xdb,
	0x34, 0x7b, 0xc5, 0x2c, 0xec, 0x1a, 0xe7, 0x2c, 0x1d, 0x3e, 0xff, 0xd9, 0x1d, 0x07, 0x72, 0x92,
	0x5c, 0xb9, 0x3e, 0x9b, 0xf6, 0x27, 0x69, 0x8c, 0x3c, 0xc4, 0xd1, 0x18, 0x79, 0xff, 0x5a, 0xbd,
	0x2a, 0x7a, 0x1d, 0x16, 0xfd, 0x1c, 0x7c, 0xb5, 0xaa, 0x24, 0x9f, 0xfd, 0x13, 0x00, 0x00, 0xff,
	0xff, 0x32, 0x94, 0x32, 0xb7, 0x33, 0x0b, 0x00, 0x00,
}
