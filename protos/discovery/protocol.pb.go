



package discovery

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import gossip "github.com/mcc-github/blockchain/protos/gossip"
import msp "github.com/mcc-github/blockchain/protos/msp"
import _ "github.com/mcc-github/blockchain/protos/msp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 






type SignedRequest struct {
	Payload   []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *SignedRequest) Reset()                    { *m = SignedRequest{} }
func (m *SignedRequest) String() string            { return proto.CompactTextString(m) }
func (*SignedRequest) ProtoMessage()               {}
func (*SignedRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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
	
	Queries []*Query `protobuf:"bytes,2,rep,name=queries" json:"queries,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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
	
	Results []*QueryResult `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Response) GetResults() []*QueryResult {
	if m != nil {
		return m.Results
	}
	return nil
}



type AuthInfo struct {
	
	
	
	ClientIdentity []byte `protobuf:"bytes,1,opt,name=client_identity,json=clientIdentity,proto3" json:"client_identity,omitempty"`
	
	
	
	
	
	
	
	ClientTlsCertHash []byte `protobuf:"bytes,2,opt,name=client_tls_cert_hash,json=clientTlsCertHash,proto3" json:"client_tls_cert_hash,omitempty"`
}

func (m *AuthInfo) Reset()                    { *m = AuthInfo{} }
func (m *AuthInfo) String() string            { return proto.CompactTextString(m) }
func (*AuthInfo) ProtoMessage()               {}
func (*AuthInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

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
	
	
	
	
	
	Query isQuery_Query `protobuf_oneof:"query"`
}

func (m *Query) Reset()                    { *m = Query{} }
func (m *Query) String() string            { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()               {}
func (*Query) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type isQuery_Query interface{ isQuery_Query() }

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
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Query_PeerQuery:
		s := proto.Size(x.PeerQuery)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Query_CcQuery:
		s := proto.Size(x.CcQuery)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Query_LocalPeers:
		s := proto.Size(x.LocalPeers)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}





type QueryResult struct {
	
	
	
	
	
	Result isQueryResult_Result `protobuf_oneof:"result"`
}

func (m *QueryResult) Reset()                    { *m = QueryResult{} }
func (m *QueryResult) String() string            { return proto.CompactTextString(m) }
func (*QueryResult) ProtoMessage()               {}
func (*QueryResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type isQueryResult_Result interface{ isQueryResult_Result() }

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
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *QueryResult_ConfigResult:
		s := proto.Size(x.ConfigResult)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *QueryResult_CcQueryRes:
		s := proto.Size(x.CcQueryRes)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *QueryResult_Members:
		s := proto.Size(x.Members)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}


type ConfigQuery struct {
}

func (m *ConfigQuery) Reset()                    { *m = ConfigQuery{} }
func (m *ConfigQuery) String() string            { return proto.CompactTextString(m) }
func (*ConfigQuery) ProtoMessage()               {}
func (*ConfigQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type ConfigResult struct {
	
	Msps map[string]*msp.FabricMSPConfig `protobuf:"bytes,1,rep,name=msps" json:"msps,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	
	Orderers map[string]*Endpoints `protobuf:"bytes,2,rep,name=orderers" json:"orderers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ConfigResult) Reset()                    { *m = ConfigResult{} }
func (m *ConfigResult) String() string            { return proto.CompactTextString(m) }
func (*ConfigResult) ProtoMessage()               {}
func (*ConfigResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

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
}

func (m *PeerMembershipQuery) Reset()                    { *m = PeerMembershipQuery{} }
func (m *PeerMembershipQuery) String() string            { return proto.CompactTextString(m) }
func (*PeerMembershipQuery) ProtoMessage()               {}
func (*PeerMembershipQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }


type PeerMembershipResult struct {
	PeersByOrg map[string]*Peers `protobuf:"bytes,1,rep,name=peers_by_org,json=peersByOrg" json:"peers_by_org,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *PeerMembershipResult) Reset()                    { *m = PeerMembershipResult{} }
func (m *PeerMembershipResult) String() string            { return proto.CompactTextString(m) }
func (*PeerMembershipResult) ProtoMessage()               {}
func (*PeerMembershipResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *PeerMembershipResult) GetPeersByOrg() map[string]*Peers {
	if m != nil {
		return m.PeersByOrg
	}
	return nil
}





type ChaincodeQuery struct {
	Interests []*ChaincodeInterest `protobuf:"bytes,1,rep,name=interests" json:"interests,omitempty"`
}

func (m *ChaincodeQuery) Reset()                    { *m = ChaincodeQuery{} }
func (m *ChaincodeQuery) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeQuery) ProtoMessage()               {}
func (*ChaincodeQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *ChaincodeQuery) GetInterests() []*ChaincodeInterest {
	if m != nil {
		return m.Interests
	}
	return nil
}




type ChaincodeInterest struct {
	Chaincodes []*ChaincodeCall `protobuf:"bytes,1,rep,name=chaincodes" json:"chaincodes,omitempty"`
}

func (m *ChaincodeInterest) Reset()                    { *m = ChaincodeInterest{} }
func (m *ChaincodeInterest) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeInterest) ProtoMessage()               {}
func (*ChaincodeInterest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *ChaincodeInterest) GetChaincodes() []*ChaincodeCall {
	if m != nil {
		return m.Chaincodes
	}
	return nil
}



type ChaincodeCall struct {
	Name            string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	CollectionNames []string `protobuf:"bytes,2,rep,name=collection_names,json=collectionNames" json:"collection_names,omitempty"`
}

func (m *ChaincodeCall) Reset()                    { *m = ChaincodeCall{} }
func (m *ChaincodeCall) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeCall) ProtoMessage()               {}
func (*ChaincodeCall) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

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
	Content []*EndorsementDescriptor `protobuf:"bytes,1,rep,name=content" json:"content,omitempty"`
}

func (m *ChaincodeQueryResult) Reset()                    { *m = ChaincodeQueryResult{} }
func (m *ChaincodeQueryResult) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeQueryResult) ProtoMessage()               {}
func (*ChaincodeQueryResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *ChaincodeQueryResult) GetContent() []*EndorsementDescriptor {
	if m != nil {
		return m.Content
	}
	return nil
}


type LocalPeerQuery struct {
}

func (m *LocalPeerQuery) Reset()                    { *m = LocalPeerQuery{} }
func (m *LocalPeerQuery) String() string            { return proto.CompactTextString(m) }
func (*LocalPeerQuery) ProtoMessage()               {}
func (*LocalPeerQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }













type EndorsementDescriptor struct {
	Chaincode string `protobuf:"bytes,1,opt,name=chaincode" json:"chaincode,omitempty"`
	
	EndorsersByGroups map[string]*Peers `protobuf:"bytes,2,rep,name=endorsers_by_groups,json=endorsersByGroups" json:"endorsers_by_groups,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	
	
	
	Layouts []*Layout `protobuf:"bytes,3,rep,name=layouts" json:"layouts,omitempty"`
}

func (m *EndorsementDescriptor) Reset()                    { *m = EndorsementDescriptor{} }
func (m *EndorsementDescriptor) String() string            { return proto.CompactTextString(m) }
func (*EndorsementDescriptor) ProtoMessage()               {}
func (*EndorsementDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

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
	
	
	QuantitiesByGroup map[string]uint32 `protobuf:"bytes,1,rep,name=quantities_by_group,json=quantitiesByGroup" json:"quantities_by_group,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
}

func (m *Layout) Reset()                    { *m = Layout{} }
func (m *Layout) String() string            { return proto.CompactTextString(m) }
func (*Layout) ProtoMessage()               {}
func (*Layout) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *Layout) GetQuantitiesByGroup() map[string]uint32 {
	if m != nil {
		return m.QuantitiesByGroup
	}
	return nil
}


type Peers struct {
	Peers []*Peer `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty"`
}

func (m *Peers) Reset()                    { *m = Peers{} }
func (m *Peers) String() string            { return proto.CompactTextString(m) }
func (*Peers) ProtoMessage()               {}
func (*Peers) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *Peers) GetPeers() []*Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}



type Peer struct {
	
	StateInfo *gossip.Envelope `protobuf:"bytes,1,opt,name=state_info,json=stateInfo" json:"state_info,omitempty"`
	
	MembershipInfo *gossip.Envelope `protobuf:"bytes,2,opt,name=membership_info,json=membershipInfo" json:"membership_info,omitempty"`
	
	Identity []byte `protobuf:"bytes,3,opt,name=identity,proto3" json:"identity,omitempty"`
}

func (m *Peer) Reset()                    { *m = Peer{} }
func (m *Peer) String() string            { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()               {}
func (*Peer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

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
	Content string `protobuf:"bytes,1,opt,name=content" json:"content,omitempty"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

func (m *Error) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}


type Endpoints struct {
	Endpoint []*Endpoint `protobuf:"bytes,1,rep,name=endpoint" json:"endpoint,omitempty"`
}

func (m *Endpoints) Reset()                    { *m = Endpoints{} }
func (m *Endpoints) String() string            { return proto.CompactTextString(m) }
func (*Endpoints) ProtoMessage()               {}
func (*Endpoints) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

func (m *Endpoints) GetEndpoint() []*Endpoint {
	if m != nil {
		return m.Endpoint
	}
	return nil
}


type Endpoint struct {
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	Port uint32 `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
}

func (m *Endpoint) Reset()                    { *m = Endpoint{} }
func (m *Endpoint) String() string            { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()               {}
func (*Endpoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{21} }

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
	proto.RegisterType((*PeerMembershipQuery)(nil), "discovery.PeerMembershipQuery")
	proto.RegisterType((*PeerMembershipResult)(nil), "discovery.PeerMembershipResult")
	proto.RegisterType((*ChaincodeQuery)(nil), "discovery.ChaincodeQuery")
	proto.RegisterType((*ChaincodeInterest)(nil), "discovery.ChaincodeInterest")
	proto.RegisterType((*ChaincodeCall)(nil), "discovery.ChaincodeCall")
	proto.RegisterType((*ChaincodeQueryResult)(nil), "discovery.ChaincodeQueryResult")
	proto.RegisterType((*LocalPeerQuery)(nil), "discovery.LocalPeerQuery")
	proto.RegisterType((*EndorsementDescriptor)(nil), "discovery.EndorsementDescriptor")
	proto.RegisterType((*Layout)(nil), "discovery.Layout")
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

func init() { proto.RegisterFile("discovery/protocol.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0x5b, 0x6f, 0x1b, 0x45,
	0x14, 0x8e, 0x9d, 0xb8, 0xb6, 0x8f, 0xe3, 0x5c, 0x26, 0x6e, 0x30, 0x56, 0x45, 0xd3, 0x95, 0x0a,
	0xa1, 0x48, 0xeb, 0x2a, 0x08, 0x28, 0x49, 0x04, 0x6a, 0x2e, 0xd4, 0x91, 0x9a, 0x26, 0x99, 0x22,
	0x84, 0x78, 0xb1, 0x36, 0xeb, 0x13, 0x7b, 0xc5, 0x7a, 0x67, 0x33, 0x33, 0x1b, 0x69, 0x9f, 0x79,
	0xe7, 0x27, 0xf0, 0xc2, 0x0b, 0xe2, 0x27, 0xf0, 0xeb, 0xd0, 0xce, 0x65, 0xbd, 0x76, 0x36, 0x2a,
	0x12, 0x6f, 0x33, 0xe7, 0x9c, 0xef, 0x3b, 0x97, 0x39, 0x33, 0x73, 0xa0, 0x3b, 0x0a, 0x84, 0xcf,
	0xee, 0x90, 0xa7, 0xfd, 0x98, 0x33, 0xc9, 0x7c, 0x16, 0xba, 0x6a, 0x41, 0x9a, 0xb9, 0xa6, 0xd7,
	0x19, 0x33, 0x21, 0x82, 0xb8, 0x3f, 0x45, 0x21, 0xbc, 0x31, 0x6a, 0x83, 0x5e, 0x67, 0x2a, 0xe2,
	0xfe, 0x54, 0xc4, 0x43, 0x9f, 0x45, 0x37, 0xc1, 0xb8, 0x28, 0x0d, 0x46, 0x18, 0xc9, 0x40, 0x06,
	0x28, 0xb4, 0xd4, 0x79, 0x03, 0xed, 0xf7, 0xc1, 0x38, 0xc2, 0x11, 0xc5, 0xdb, 0x04, 0x85, 0x24,
	0x5d, 0xa8, 0xc7, 0x5e, 0x1a, 0x32, 0x6f, 0xd4, 0xad, 0xec, 0x54, 0x76, 0x57, 0xa9, 0xdd, 0x92,
	0x27, 0xd0, 0x14, 0xc1, 0x38, 0xf2, 0x64, 0xc2, 0xb1, 0x5b, 0x55, 0xba, 0x99, 0xc0, 0xe1, 0x50,
	0xb7, 0x14, 0x07, 0xb0, 0xe6, 0x25, 0x72, 0x92, 0x79, 0xf2, 0x3d, 0x19, 0xb0, 0x48, 0x31, 0xb5,
	0xf6, 0xb6, 0xdc, 0x3c, 0x72, 0xf7, 0x75, 0x22, 0x27, 0x67, 0xd1, 0x0d, 0xa3, 0x0b, 0xa6, 0xe4,
	0x05, 0xd4, 0x6f, 0x13, 0xe4, 0x01, 0x8a, 0x6e, 0x75, 0x67, 0x79, 0xb7, 0xb5, 0xb7, 0x51, 0x40,
	0x5d, 0x25, 0xc8, 0x53, 0x6a, 0x0d, 0x9c, 0x43, 0x68, 0x50, 0x14, 0x31, 0x8b, 0x04, 0x92, 0x97,
	0x50, 0xe7, 0x28, 0x92, 0x50, 0x8a, 0x6e, 0x45, 0xe1, 0xb6, 0xef, 0xe1, 0x94, 0x9a, 0x5a, 0x33,
	0x67, 0x04, 0x0d, 0x1b, 0x05, 0xf9, 0x0c, 0xd6, 0xfd, 0x30, 0xc0, 0x48, 0x0e, 0x4d, 0x85, 0x52,
	0x93, 0xfd, 0x9a, 0x16, 0x9f, 0x19, 0x29, 0xe9, 0x43, 0xc7, 0x18, 0xca, 0x50, 0x0c, 0x7d, 0xe4,
	0x72, 0x38, 0xf1, 0xc4, 0xc4, 0xd4, 0x63, 0x53, 0xeb, 0x7e, 0x0c, 0xc5, 0x31, 0x72, 0x39, 0xf0,
	0xc4, 0xc4, 0xf9, 0xa3, 0x0a, 0x35, 0xe5, 0x3e, 0xab, 0xac, 0x3f, 0xf1, 0xa2, 0x08, 0x43, 0xc5,
	0xdd, 0xa4, 0x76, 0x4b, 0x0e, 0x60, 0x55, 0x1f, 0xd5, 0x30, 0xcb, 0x2c, 0x55, 0x64, 0xf3, 0x09,
	0x1c, 0x2b, 0xb5, 0xe2, 0x19, 0x2c, 0xd1, 0x96, 0x3f, 0xdb, 0x92, 0xef, 0x01, 0x62, 0x44, 0x6e,
	0xa0, 0xcb, 0x0a, 0xfa, 0x49, 0x01, 0x7a, 0x89, 0xc8, 0xcf, 0x71, 0x7a, 0x8d, 0x5c, 0x4c, 0x82,
	0xd8, 0x52, 0x34, 0x33, 0x8c, 0x26, 0xf8, 0x1a, 0x1a, 0xbe, 0x6f, 0xe0, 0x2b, 0x0a, 0xfe, 0x71,
	0xd1, 0xf3, 0xc4, 0x0b, 0x22, 0x9f, 0x8d, 0xd0, 0x22, 0xeb, 0xbe, 0xaf, 0x71, 0x87, 0xd0, 0x0a,
	0x99, 0xef, 0x85, 0xc3, 0x8c, 0x4a, 0x74, 0x6b, 0xf7, 0xa0, 0x6f, 0x33, 0xed, 0xa5, 0xf5, 0x33,
	0x58, 0xa2, 0x10, 0x5a, 0x89, 0x38, 0xaa, 0x43, 0x4d, 0xb9, 0x74, 0x7e, 0xab, 0x42, 0xab, 0x70,
	0x3e, 0x64, 0x17, 0x6a, 0xc8, 0x39, 0xe3, 0xa6, 0x69, 0x8a, 0xc7, 0x7f, 0x9a, 0xc9, 0x07, 0x4b,
	0x54, 0x1b, 0x90, 0xef, 0xa0, 0x6d, 0xca, 0xa6, 0x8f, 0xd4, 0xd4, 0xed, 0xa3, 0x7b, 0x75, 0xd3,
	0xcc, 0x83, 0x25, 0x6a, 0xca, 0x6c, 0x3c, 0x1d, 0xc3, 0xaa, 0x4d, 0x3c, 0x63, 0x30, 0xb5, 0x7b,
	0xfa, 0x60, 0xf2, 0x39, 0x0d, 0x98, 0x12, 0x50, 0x14, 0xe4, 0x00, 0xea, 0x53, 0x5d, 0x5d, 0x53,
	0xbc, 0xa7, 0x0f, 0xd6, 0x3e, 0xc7, 0x5b, 0xc4, 0x51, 0x03, 0x1e, 0xe9, 0xd0, 0x9d, 0x36, 0xb4,
	0x0a, 0x67, 0xec, 0xfc, 0x5d, 0x85, 0xd5, 0x62, 0xec, 0xe4, 0x2b, 0x58, 0x99, 0x8a, 0xd8, 0xf6,
	0xf6, 0xb3, 0x07, 0x52, 0x74, 0xcf, 0x45, 0x2c, 0x4e, 0x23, 0xc9, 0x53, 0xaa, 0xcc, 0xc9, 0x6b,
	0x68, 0x30, 0x3e, 0x42, 0x9e, 0x85, 0xa7, 0xaf, 0xd3, 0xf3, 0x87, 0xa0, 0x17, 0xc6, 0x4e, 0xc3,
	0x73, 0x58, 0xef, 0x1c, 0x9a, 0x39, 0x2b, 0xd9, 0x80, 0xe5, 0x5f, 0x31, 0x35, 0xfd, 0x9b, 0x2d,
	0xc9, 0x0b, 0xa8, 0xdd, 0x79, 0x61, 0x82, 0xa6, 0xf8, 0x1d, 0x77, 0x2a, 0x62, 0xf7, 0x07, 0xef,
	0x9a, 0x07, 0xfe, 0xf9, 0xfb, 0x4b, 0xe3, 0x41, 0x9b, 0xec, 0x57, 0x5f, 0x55, 0x7a, 0x57, 0xd0,
	0x9e, 0xf3, 0xf4, 0x5f, 0x28, 0x0b, 0x1d, 0x10, 0x8d, 0x62, 0x16, 0x44, 0x52, 0x14, 0x28, 0x9d,
	0xc7, 0xb0, 0x55, 0xd2, 0xe4, 0xce, 0x3f, 0x15, 0xe8, 0x94, 0x1d, 0x00, 0xb9, 0x82, 0x55, 0xd5,
	0xb2, 0xc3, 0xeb, 0x74, 0xc8, 0xf8, 0xd8, 0xd4, 0xb4, 0xff, 0x81, 0x73, 0x73, 0x75, 0xdf, 0xa6,
	0x17, 0x7c, 0xac, 0x4b, 0xa4, 0xae, 0x9d, 0x16, 0xf4, 0x2e, 0x60, 0x7d, 0x41, 0x5d, 0x92, 0xd7,
	0xa7, 0xf3, 0x79, 0x6d, 0x2c, 0x38, 0x9c, 0xcb, 0xe9, 0x2d, 0xac, 0xcd, 0x37, 0x1f, 0xd9, 0x87,
	0x66, 0x10, 0x49, 0xe4, 0x28, 0xf2, 0x27, 0xee, 0x49, 0x59, 0xab, 0x9e, 0x19, 0x23, 0x3a, 0x33,
	0x77, 0xce, 0x61, 0xf3, 0x9e, 0x9e, 0xbc, 0x02, 0xf0, 0xad, 0xd0, 0x32, 0x76, 0xcb, 0x18, 0x8f,
	0xbd, 0x30, 0xa4, 0x05, 0x5b, 0xe7, 0x1d, 0xb4, 0xe7, 0x94, 0x84, 0xc0, 0x4a, 0xe4, 0x4d, 0xd1,
	0x24, 0xab, 0xd6, 0xe4, 0x73, 0xd8, 0xf0, 0x59, 0x18, 0xa2, 0x9f, 0x3d, 0xeb, 0xc3, 0x4c, 0xa4,
	0x5b, 0xb0, 0x49, 0xd7, 0x67, 0xf2, 0x77, 0x99, 0xd8, 0xa1, 0xd0, 0x29, 0xbb, 0x69, 0x64, 0x1f,
	0xea, 0x3e, 0x8b, 0x24, 0x46, 0xd2, 0x84, 0xb7, 0x33, 0xdf, 0x0a, 0x8c, 0x0b, 0x9c, 0x62, 0x24,
	0x4f, 0x50, 0xf8, 0x3c, 0x88, 0x25, 0xe3, 0xd4, 0x02, 0x9c, 0x0d, 0x58, 0x9b, 0x7f, 0x7f, 0x9c,
	0x3f, 0xab, 0xf0, 0xb8, 0x14, 0x94, 0xfd, 0x6c, 0x79, 0x76, 0x26, 0x87, 0x99, 0x80, 0x8c, 0x61,
	0x0b, 0x35, 0x4c, 0xb7, 0xcc, 0x98, 0xb3, 0x24, 0xb6, 0xd7, 0xe9, 0x9b, 0x0f, 0x45, 0x64, 0xa5,
	0x59, 0x6f, 0xbc, 0x51, 0x48, 0xdd, 0x3d, 0x9b, 0xb8, 0x28, 0x27, 0x5f, 0x40, 0x3d, 0xf4, 0x52,
	0x96, 0xc8, 0xec, 0x29, 0xca, 0xc8, 0x37, 0x8b, 0x8f, 0xa9, 0xd2, 0x50, 0x6b, 0xd1, 0xfb, 0x09,
	0xb6, 0xcb, 0x99, 0xff, 0x67, 0xe3, 0xfd, 0x55, 0x81, 0x47, 0xda, 0x17, 0xf9, 0x19, 0xb6, 0x6e,
	0x13, 0xcf, 0xcc, 0x0b, 0x79, 0xe6, 0xe6, 0x28, 0x76, 0xef, 0xc5, 0xe6, 0x5e, 0xe5, 0xc6, 0x26,
	0x20, 0x93, 0xe9, 0xed, 0xa2, 0xbc, 0x77, 0x02, 0xdb, 0xe5, 0xc6, 0x25, 0xc1, 0x77, 0x8a, 0xc1,
	0xb7, 0x8b, 0xa1, 0xba, 0x50, 0x53, 0xe1, 0x93, 0xe7, 0x50, 0xd3, 0x7f, 0x90, 0x0e, 0x6d, 0x7d,
	0x21, 0x3f, 0xaa, 0xb5, 0xce, 0xef, 0x15, 0x58, 0xc9, 0xf6, 0xa4, 0x0f, 0x20, 0xa4, 0x27, 0x71,
	0x18, 0x44, 0x37, 0x2c, 0xff, 0x67, 0xf4, 0x2c, 0xe5, 0x9e, 0x46, 0x77, 0x18, 0xb2, 0x18, 0x69,
	0x53, 0xd9, 0xa8, 0xf1, 0xe0, 0x5b, 0x58, 0x9f, 0xe6, 0xcf, 0x81, 0x46, 0x55, 0x1f, 0x40, 0xad,
	0xcd, 0x0c, 0x15, 0xb4, 0x07, 0x8d, 0x7c, 0xa4, 0x58, 0x56, 0x43, 0x42, 0xbe, 0x77, 0x9e, 0x41,
	0x4d, 0x7d, 0x69, 0x6a, 0x34, 0xc8, 0x1b, 0x5d, 0x8f, 0x06, 0xa6, 0x8d, 0x0f, 0xa1, 0x99, 0xbf,
	0x79, 0xa4, 0x0f, 0x0d, 0x34, 0x1b, 0x93, 0xea, 0x56, 0xc9, 0xdb, 0x48, 0x73, 0x23, 0x67, 0x0f,
	0x1a, 0x56, 0x9a, 0xdd, 0xd1, 0x09, 0x13, 0xd6, 0x81, 0x5a, 0x67, 0xb2, 0x98, 0x71, 0x69, 0x4a,
	0xab, 0xd6, 0x7b, 0x03, 0x68, 0x9e, 0x58, 0x4e, 0x72, 0x00, 0x0d, 0xbb, 0x21, 0xc5, 0xb7, 0x61,
	0x6e, 0x66, 0xec, 0x15, 0xa3, 0xb0, 0x03, 0x99, 0xb3, 0x74, 0xf4, 0xf2, 0x17, 0x77, 0x1c, 0xc8,
	0x49, 0x72, 0xed, 0xfa, 0x6c, 0xda, 0x9f, 0xa4, 0x31, 0xf2, 0x10, 0x47, 0x63, 0xe4, 0xfd, 0x1b,
	0xf5, 0x3f, 0xe8, 0xc1, 0x56, 0xf4, 0x73, 0xf0, 0xf5, 0x23, 0x25, 0xf9, 0xf2, 0xdf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x59, 0x7e, 0x3b, 0x2c, 0xfd, 0x0a, 0x00, 0x00,
}
