
package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/filter"
	"github.com/mcc-github/blockchain/gossip/gossip"
	"github.com/mcc-github/blockchain/gossip/protoext"
	gossipa "github.com/mcc-github/blockchain/protos/gossip"
)

type Gossip struct {
	AcceptStub        func(common.MessageAcceptor, bool) (<-chan *gossipa.GossipMessage, <-chan protoext.ReceivedMessage)
	acceptMutex       sync.RWMutex
	acceptArgsForCall []struct {
		arg1 common.MessageAcceptor
		arg2 bool
	}
	acceptReturns struct {
		result1 <-chan *gossipa.GossipMessage
		result2 <-chan protoext.ReceivedMessage
	}
	acceptReturnsOnCall map[int]struct {
		result1 <-chan *gossipa.GossipMessage
		result2 <-chan protoext.ReceivedMessage
	}
	GossipStub        func(*gossipa.GossipMessage)
	gossipMutex       sync.RWMutex
	gossipArgsForCall []struct {
		arg1 *gossipa.GossipMessage
	}
	IdentityInfoStub        func() api.PeerIdentitySet
	identityInfoMutex       sync.RWMutex
	identityInfoArgsForCall []struct {
	}
	identityInfoReturns struct {
		result1 api.PeerIdentitySet
	}
	identityInfoReturnsOnCall map[int]struct {
		result1 api.PeerIdentitySet
	}
	IsInMyOrgStub        func(discovery.NetworkMember) bool
	isInMyOrgMutex       sync.RWMutex
	isInMyOrgArgsForCall []struct {
		arg1 discovery.NetworkMember
	}
	isInMyOrgReturns struct {
		result1 bool
	}
	isInMyOrgReturnsOnCall map[int]struct {
		result1 bool
	}
	JoinChanStub        func(api.JoinChannelMessage, common.ChainID)
	joinChanMutex       sync.RWMutex
	joinChanArgsForCall []struct {
		arg1 api.JoinChannelMessage
		arg2 common.ChainID
	}
	LeaveChanStub        func(common.ChainID)
	leaveChanMutex       sync.RWMutex
	leaveChanArgsForCall []struct {
		arg1 common.ChainID
	}
	PeerFilterStub        func(common.ChainID, api.SubChannelSelectionCriteria) (filter.RoutingFilter, error)
	peerFilterMutex       sync.RWMutex
	peerFilterArgsForCall []struct {
		arg1 common.ChainID
		arg2 api.SubChannelSelectionCriteria
	}
	peerFilterReturns struct {
		result1 filter.RoutingFilter
		result2 error
	}
	peerFilterReturnsOnCall map[int]struct {
		result1 filter.RoutingFilter
		result2 error
	}
	PeersStub        func() []discovery.NetworkMember
	peersMutex       sync.RWMutex
	peersArgsForCall []struct {
	}
	peersReturns struct {
		result1 []discovery.NetworkMember
	}
	peersReturnsOnCall map[int]struct {
		result1 []discovery.NetworkMember
	}
	PeersOfChannelStub        func(common.ChainID) []discovery.NetworkMember
	peersOfChannelMutex       sync.RWMutex
	peersOfChannelArgsForCall []struct {
		arg1 common.ChainID
	}
	peersOfChannelReturns struct {
		result1 []discovery.NetworkMember
	}
	peersOfChannelReturnsOnCall map[int]struct {
		result1 []discovery.NetworkMember
	}
	SelfChannelInfoStub        func(common.ChainID) *protoext.SignedGossipMessage
	selfChannelInfoMutex       sync.RWMutex
	selfChannelInfoArgsForCall []struct {
		arg1 common.ChainID
	}
	selfChannelInfoReturns struct {
		result1 *protoext.SignedGossipMessage
	}
	selfChannelInfoReturnsOnCall map[int]struct {
		result1 *protoext.SignedGossipMessage
	}
	SelfMembershipInfoStub        func() discovery.NetworkMember
	selfMembershipInfoMutex       sync.RWMutex
	selfMembershipInfoArgsForCall []struct {
	}
	selfMembershipInfoReturns struct {
		result1 discovery.NetworkMember
	}
	selfMembershipInfoReturnsOnCall map[int]struct {
		result1 discovery.NetworkMember
	}
	SendStub        func(*gossipa.GossipMessage, ...*comm.RemotePeer)
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		arg1 *gossipa.GossipMessage
		arg2 []*comm.RemotePeer
	}
	SendByCriteriaStub        func(*protoext.SignedGossipMessage, gossip.SendCriteria) error
	sendByCriteriaMutex       sync.RWMutex
	sendByCriteriaArgsForCall []struct {
		arg1 *protoext.SignedGossipMessage
		arg2 gossip.SendCriteria
	}
	sendByCriteriaReturns struct {
		result1 error
	}
	sendByCriteriaReturnsOnCall map[int]struct {
		result1 error
	}
	StopStub        func()
	stopMutex       sync.RWMutex
	stopArgsForCall []struct {
	}
	SuspectPeersStub        func(api.PeerSuspector)
	suspectPeersMutex       sync.RWMutex
	suspectPeersArgsForCall []struct {
		arg1 api.PeerSuspector
	}
	UpdateChaincodesStub        func([]*gossipa.Chaincode, common.ChainID)
	updateChaincodesMutex       sync.RWMutex
	updateChaincodesArgsForCall []struct {
		arg1 []*gossipa.Chaincode
		arg2 common.ChainID
	}
	UpdateLedgerHeightStub        func(uint64, common.ChainID)
	updateLedgerHeightMutex       sync.RWMutex
	updateLedgerHeightArgsForCall []struct {
		arg1 uint64
		arg2 common.ChainID
	}
	UpdateMetadataStub        func([]byte)
	updateMetadataMutex       sync.RWMutex
	updateMetadataArgsForCall []struct {
		arg1 []byte
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Gossip) Accept(arg1 common.MessageAcceptor, arg2 bool) (<-chan *gossipa.GossipMessage, <-chan protoext.ReceivedMessage) {
	fake.acceptMutex.Lock()
	ret, specificReturn := fake.acceptReturnsOnCall[len(fake.acceptArgsForCall)]
	fake.acceptArgsForCall = append(fake.acceptArgsForCall, struct {
		arg1 common.MessageAcceptor
		arg2 bool
	}{arg1, arg2})
	fake.recordInvocation("Accept", []interface{}{arg1, arg2})
	fake.acceptMutex.Unlock()
	if fake.AcceptStub != nil {
		return fake.AcceptStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.acceptReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Gossip) AcceptCallCount() int {
	fake.acceptMutex.RLock()
	defer fake.acceptMutex.RUnlock()
	return len(fake.acceptArgsForCall)
}

func (fake *Gossip) AcceptCalls(stub func(common.MessageAcceptor, bool) (<-chan *gossipa.GossipMessage, <-chan protoext.ReceivedMessage)) {
	fake.acceptMutex.Lock()
	defer fake.acceptMutex.Unlock()
	fake.AcceptStub = stub
}

func (fake *Gossip) AcceptArgsForCall(i int) (common.MessageAcceptor, bool) {
	fake.acceptMutex.RLock()
	defer fake.acceptMutex.RUnlock()
	argsForCall := fake.acceptArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Gossip) AcceptReturns(result1 <-chan *gossipa.GossipMessage, result2 <-chan protoext.ReceivedMessage) {
	fake.acceptMutex.Lock()
	defer fake.acceptMutex.Unlock()
	fake.AcceptStub = nil
	fake.acceptReturns = struct {
		result1 <-chan *gossipa.GossipMessage
		result2 <-chan protoext.ReceivedMessage
	}{result1, result2}
}

func (fake *Gossip) AcceptReturnsOnCall(i int, result1 <-chan *gossipa.GossipMessage, result2 <-chan protoext.ReceivedMessage) {
	fake.acceptMutex.Lock()
	defer fake.acceptMutex.Unlock()
	fake.AcceptStub = nil
	if fake.acceptReturnsOnCall == nil {
		fake.acceptReturnsOnCall = make(map[int]struct {
			result1 <-chan *gossipa.GossipMessage
			result2 <-chan protoext.ReceivedMessage
		})
	}
	fake.acceptReturnsOnCall[i] = struct {
		result1 <-chan *gossipa.GossipMessage
		result2 <-chan protoext.ReceivedMessage
	}{result1, result2}
}

func (fake *Gossip) Gossip(arg1 *gossipa.GossipMessage) {
	fake.gossipMutex.Lock()
	fake.gossipArgsForCall = append(fake.gossipArgsForCall, struct {
		arg1 *gossipa.GossipMessage
	}{arg1})
	fake.recordInvocation("Gossip", []interface{}{arg1})
	fake.gossipMutex.Unlock()
	if fake.GossipStub != nil {
		fake.GossipStub(arg1)
	}
}

func (fake *Gossip) GossipCallCount() int {
	fake.gossipMutex.RLock()
	defer fake.gossipMutex.RUnlock()
	return len(fake.gossipArgsForCall)
}

func (fake *Gossip) GossipCalls(stub func(*gossipa.GossipMessage)) {
	fake.gossipMutex.Lock()
	defer fake.gossipMutex.Unlock()
	fake.GossipStub = stub
}

func (fake *Gossip) GossipArgsForCall(i int) *gossipa.GossipMessage {
	fake.gossipMutex.RLock()
	defer fake.gossipMutex.RUnlock()
	argsForCall := fake.gossipArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gossip) IdentityInfo() api.PeerIdentitySet {
	fake.identityInfoMutex.Lock()
	ret, specificReturn := fake.identityInfoReturnsOnCall[len(fake.identityInfoArgsForCall)]
	fake.identityInfoArgsForCall = append(fake.identityInfoArgsForCall, struct {
	}{})
	fake.recordInvocation("IdentityInfo", []interface{}{})
	fake.identityInfoMutex.Unlock()
	if fake.IdentityInfoStub != nil {
		return fake.IdentityInfoStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.identityInfoReturns
	return fakeReturns.result1
}

func (fake *Gossip) IdentityInfoCallCount() int {
	fake.identityInfoMutex.RLock()
	defer fake.identityInfoMutex.RUnlock()
	return len(fake.identityInfoArgsForCall)
}

func (fake *Gossip) IdentityInfoCalls(stub func() api.PeerIdentitySet) {
	fake.identityInfoMutex.Lock()
	defer fake.identityInfoMutex.Unlock()
	fake.IdentityInfoStub = stub
}

func (fake *Gossip) IdentityInfoReturns(result1 api.PeerIdentitySet) {
	fake.identityInfoMutex.Lock()
	defer fake.identityInfoMutex.Unlock()
	fake.IdentityInfoStub = nil
	fake.identityInfoReturns = struct {
		result1 api.PeerIdentitySet
	}{result1}
}

func (fake *Gossip) IdentityInfoReturnsOnCall(i int, result1 api.PeerIdentitySet) {
	fake.identityInfoMutex.Lock()
	defer fake.identityInfoMutex.Unlock()
	fake.IdentityInfoStub = nil
	if fake.identityInfoReturnsOnCall == nil {
		fake.identityInfoReturnsOnCall = make(map[int]struct {
			result1 api.PeerIdentitySet
		})
	}
	fake.identityInfoReturnsOnCall[i] = struct {
		result1 api.PeerIdentitySet
	}{result1}
}

func (fake *Gossip) IsInMyOrg(arg1 discovery.NetworkMember) bool {
	fake.isInMyOrgMutex.Lock()
	ret, specificReturn := fake.isInMyOrgReturnsOnCall[len(fake.isInMyOrgArgsForCall)]
	fake.isInMyOrgArgsForCall = append(fake.isInMyOrgArgsForCall, struct {
		arg1 discovery.NetworkMember
	}{arg1})
	fake.recordInvocation("IsInMyOrg", []interface{}{arg1})
	fake.isInMyOrgMutex.Unlock()
	if fake.IsInMyOrgStub != nil {
		return fake.IsInMyOrgStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.isInMyOrgReturns
	return fakeReturns.result1
}

func (fake *Gossip) IsInMyOrgCallCount() int {
	fake.isInMyOrgMutex.RLock()
	defer fake.isInMyOrgMutex.RUnlock()
	return len(fake.isInMyOrgArgsForCall)
}

func (fake *Gossip) IsInMyOrgCalls(stub func(discovery.NetworkMember) bool) {
	fake.isInMyOrgMutex.Lock()
	defer fake.isInMyOrgMutex.Unlock()
	fake.IsInMyOrgStub = stub
}

func (fake *Gossip) IsInMyOrgArgsForCall(i int) discovery.NetworkMember {
	fake.isInMyOrgMutex.RLock()
	defer fake.isInMyOrgMutex.RUnlock()
	argsForCall := fake.isInMyOrgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gossip) IsInMyOrgReturns(result1 bool) {
	fake.isInMyOrgMutex.Lock()
	defer fake.isInMyOrgMutex.Unlock()
	fake.IsInMyOrgStub = nil
	fake.isInMyOrgReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Gossip) IsInMyOrgReturnsOnCall(i int, result1 bool) {
	fake.isInMyOrgMutex.Lock()
	defer fake.isInMyOrgMutex.Unlock()
	fake.IsInMyOrgStub = nil
	if fake.isInMyOrgReturnsOnCall == nil {
		fake.isInMyOrgReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isInMyOrgReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Gossip) JoinChan(arg1 api.JoinChannelMessage, arg2 common.ChainID) {
	fake.joinChanMutex.Lock()
	fake.joinChanArgsForCall = append(fake.joinChanArgsForCall, struct {
		arg1 api.JoinChannelMessage
		arg2 common.ChainID
	}{arg1, arg2})
	fake.recordInvocation("JoinChan", []interface{}{arg1, arg2})
	fake.joinChanMutex.Unlock()
	if fake.JoinChanStub != nil {
		fake.JoinChanStub(arg1, arg2)
	}
}

func (fake *Gossip) JoinChanCallCount() int {
	fake.joinChanMutex.RLock()
	defer fake.joinChanMutex.RUnlock()
	return len(fake.joinChanArgsForCall)
}

func (fake *Gossip) JoinChanCalls(stub func(api.JoinChannelMessage, common.ChainID)) {
	fake.joinChanMutex.Lock()
	defer fake.joinChanMutex.Unlock()
	fake.JoinChanStub = stub
}

func (fake *Gossip) JoinChanArgsForCall(i int) (api.JoinChannelMessage, common.ChainID) {
	fake.joinChanMutex.RLock()
	defer fake.joinChanMutex.RUnlock()
	argsForCall := fake.joinChanArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Gossip) LeaveChan(arg1 common.ChainID) {
	fake.leaveChanMutex.Lock()
	fake.leaveChanArgsForCall = append(fake.leaveChanArgsForCall, struct {
		arg1 common.ChainID
	}{arg1})
	fake.recordInvocation("LeaveChan", []interface{}{arg1})
	fake.leaveChanMutex.Unlock()
	if fake.LeaveChanStub != nil {
		fake.LeaveChanStub(arg1)
	}
}

func (fake *Gossip) LeaveChanCallCount() int {
	fake.leaveChanMutex.RLock()
	defer fake.leaveChanMutex.RUnlock()
	return len(fake.leaveChanArgsForCall)
}

func (fake *Gossip) LeaveChanCalls(stub func(common.ChainID)) {
	fake.leaveChanMutex.Lock()
	defer fake.leaveChanMutex.Unlock()
	fake.LeaveChanStub = stub
}

func (fake *Gossip) LeaveChanArgsForCall(i int) common.ChainID {
	fake.leaveChanMutex.RLock()
	defer fake.leaveChanMutex.RUnlock()
	argsForCall := fake.leaveChanArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gossip) PeerFilter(arg1 common.ChainID, arg2 api.SubChannelSelectionCriteria) (filter.RoutingFilter, error) {
	fake.peerFilterMutex.Lock()
	ret, specificReturn := fake.peerFilterReturnsOnCall[len(fake.peerFilterArgsForCall)]
	fake.peerFilterArgsForCall = append(fake.peerFilterArgsForCall, struct {
		arg1 common.ChainID
		arg2 api.SubChannelSelectionCriteria
	}{arg1, arg2})
	fake.recordInvocation("PeerFilter", []interface{}{arg1, arg2})
	fake.peerFilterMutex.Unlock()
	if fake.PeerFilterStub != nil {
		return fake.PeerFilterStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.peerFilterReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Gossip) PeerFilterCallCount() int {
	fake.peerFilterMutex.RLock()
	defer fake.peerFilterMutex.RUnlock()
	return len(fake.peerFilterArgsForCall)
}

func (fake *Gossip) PeerFilterCalls(stub func(common.ChainID, api.SubChannelSelectionCriteria) (filter.RoutingFilter, error)) {
	fake.peerFilterMutex.Lock()
	defer fake.peerFilterMutex.Unlock()
	fake.PeerFilterStub = stub
}

func (fake *Gossip) PeerFilterArgsForCall(i int) (common.ChainID, api.SubChannelSelectionCriteria) {
	fake.peerFilterMutex.RLock()
	defer fake.peerFilterMutex.RUnlock()
	argsForCall := fake.peerFilterArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Gossip) PeerFilterReturns(result1 filter.RoutingFilter, result2 error) {
	fake.peerFilterMutex.Lock()
	defer fake.peerFilterMutex.Unlock()
	fake.PeerFilterStub = nil
	fake.peerFilterReturns = struct {
		result1 filter.RoutingFilter
		result2 error
	}{result1, result2}
}

func (fake *Gossip) PeerFilterReturnsOnCall(i int, result1 filter.RoutingFilter, result2 error) {
	fake.peerFilterMutex.Lock()
	defer fake.peerFilterMutex.Unlock()
	fake.PeerFilterStub = nil
	if fake.peerFilterReturnsOnCall == nil {
		fake.peerFilterReturnsOnCall = make(map[int]struct {
			result1 filter.RoutingFilter
			result2 error
		})
	}
	fake.peerFilterReturnsOnCall[i] = struct {
		result1 filter.RoutingFilter
		result2 error
	}{result1, result2}
}

func (fake *Gossip) Peers() []discovery.NetworkMember {
	fake.peersMutex.Lock()
	ret, specificReturn := fake.peersReturnsOnCall[len(fake.peersArgsForCall)]
	fake.peersArgsForCall = append(fake.peersArgsForCall, struct {
	}{})
	fake.recordInvocation("Peers", []interface{}{})
	fake.peersMutex.Unlock()
	if fake.PeersStub != nil {
		return fake.PeersStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.peersReturns
	return fakeReturns.result1
}

func (fake *Gossip) PeersCallCount() int {
	fake.peersMutex.RLock()
	defer fake.peersMutex.RUnlock()
	return len(fake.peersArgsForCall)
}

func (fake *Gossip) PeersCalls(stub func() []discovery.NetworkMember) {
	fake.peersMutex.Lock()
	defer fake.peersMutex.Unlock()
	fake.PeersStub = stub
}

func (fake *Gossip) PeersReturns(result1 []discovery.NetworkMember) {
	fake.peersMutex.Lock()
	defer fake.peersMutex.Unlock()
	fake.PeersStub = nil
	fake.peersReturns = struct {
		result1 []discovery.NetworkMember
	}{result1}
}

func (fake *Gossip) PeersReturnsOnCall(i int, result1 []discovery.NetworkMember) {
	fake.peersMutex.Lock()
	defer fake.peersMutex.Unlock()
	fake.PeersStub = nil
	if fake.peersReturnsOnCall == nil {
		fake.peersReturnsOnCall = make(map[int]struct {
			result1 []discovery.NetworkMember
		})
	}
	fake.peersReturnsOnCall[i] = struct {
		result1 []discovery.NetworkMember
	}{result1}
}

func (fake *Gossip) PeersOfChannel(arg1 common.ChainID) []discovery.NetworkMember {
	fake.peersOfChannelMutex.Lock()
	ret, specificReturn := fake.peersOfChannelReturnsOnCall[len(fake.peersOfChannelArgsForCall)]
	fake.peersOfChannelArgsForCall = append(fake.peersOfChannelArgsForCall, struct {
		arg1 common.ChainID
	}{arg1})
	fake.recordInvocation("PeersOfChannel", []interface{}{arg1})
	fake.peersOfChannelMutex.Unlock()
	if fake.PeersOfChannelStub != nil {
		return fake.PeersOfChannelStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.peersOfChannelReturns
	return fakeReturns.result1
}

func (fake *Gossip) PeersOfChannelCallCount() int {
	fake.peersOfChannelMutex.RLock()
	defer fake.peersOfChannelMutex.RUnlock()
	return len(fake.peersOfChannelArgsForCall)
}

func (fake *Gossip) PeersOfChannelCalls(stub func(common.ChainID) []discovery.NetworkMember) {
	fake.peersOfChannelMutex.Lock()
	defer fake.peersOfChannelMutex.Unlock()
	fake.PeersOfChannelStub = stub
}

func (fake *Gossip) PeersOfChannelArgsForCall(i int) common.ChainID {
	fake.peersOfChannelMutex.RLock()
	defer fake.peersOfChannelMutex.RUnlock()
	argsForCall := fake.peersOfChannelArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gossip) PeersOfChannelReturns(result1 []discovery.NetworkMember) {
	fake.peersOfChannelMutex.Lock()
	defer fake.peersOfChannelMutex.Unlock()
	fake.PeersOfChannelStub = nil
	fake.peersOfChannelReturns = struct {
		result1 []discovery.NetworkMember
	}{result1}
}

func (fake *Gossip) PeersOfChannelReturnsOnCall(i int, result1 []discovery.NetworkMember) {
	fake.peersOfChannelMutex.Lock()
	defer fake.peersOfChannelMutex.Unlock()
	fake.PeersOfChannelStub = nil
	if fake.peersOfChannelReturnsOnCall == nil {
		fake.peersOfChannelReturnsOnCall = make(map[int]struct {
			result1 []discovery.NetworkMember
		})
	}
	fake.peersOfChannelReturnsOnCall[i] = struct {
		result1 []discovery.NetworkMember
	}{result1}
}

func (fake *Gossip) SelfChannelInfo(arg1 common.ChainID) *protoext.SignedGossipMessage {
	fake.selfChannelInfoMutex.Lock()
	ret, specificReturn := fake.selfChannelInfoReturnsOnCall[len(fake.selfChannelInfoArgsForCall)]
	fake.selfChannelInfoArgsForCall = append(fake.selfChannelInfoArgsForCall, struct {
		arg1 common.ChainID
	}{arg1})
	fake.recordInvocation("SelfChannelInfo", []interface{}{arg1})
	fake.selfChannelInfoMutex.Unlock()
	if fake.SelfChannelInfoStub != nil {
		return fake.SelfChannelInfoStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.selfChannelInfoReturns
	return fakeReturns.result1
}

func (fake *Gossip) SelfChannelInfoCallCount() int {
	fake.selfChannelInfoMutex.RLock()
	defer fake.selfChannelInfoMutex.RUnlock()
	return len(fake.selfChannelInfoArgsForCall)
}

func (fake *Gossip) SelfChannelInfoCalls(stub func(common.ChainID) *protoext.SignedGossipMessage) {
	fake.selfChannelInfoMutex.Lock()
	defer fake.selfChannelInfoMutex.Unlock()
	fake.SelfChannelInfoStub = stub
}

func (fake *Gossip) SelfChannelInfoArgsForCall(i int) common.ChainID {
	fake.selfChannelInfoMutex.RLock()
	defer fake.selfChannelInfoMutex.RUnlock()
	argsForCall := fake.selfChannelInfoArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gossip) SelfChannelInfoReturns(result1 *protoext.SignedGossipMessage) {
	fake.selfChannelInfoMutex.Lock()
	defer fake.selfChannelInfoMutex.Unlock()
	fake.SelfChannelInfoStub = nil
	fake.selfChannelInfoReturns = struct {
		result1 *protoext.SignedGossipMessage
	}{result1}
}

func (fake *Gossip) SelfChannelInfoReturnsOnCall(i int, result1 *protoext.SignedGossipMessage) {
	fake.selfChannelInfoMutex.Lock()
	defer fake.selfChannelInfoMutex.Unlock()
	fake.SelfChannelInfoStub = nil
	if fake.selfChannelInfoReturnsOnCall == nil {
		fake.selfChannelInfoReturnsOnCall = make(map[int]struct {
			result1 *protoext.SignedGossipMessage
		})
	}
	fake.selfChannelInfoReturnsOnCall[i] = struct {
		result1 *protoext.SignedGossipMessage
	}{result1}
}

func (fake *Gossip) SelfMembershipInfo() discovery.NetworkMember {
	fake.selfMembershipInfoMutex.Lock()
	ret, specificReturn := fake.selfMembershipInfoReturnsOnCall[len(fake.selfMembershipInfoArgsForCall)]
	fake.selfMembershipInfoArgsForCall = append(fake.selfMembershipInfoArgsForCall, struct {
	}{})
	fake.recordInvocation("SelfMembershipInfo", []interface{}{})
	fake.selfMembershipInfoMutex.Unlock()
	if fake.SelfMembershipInfoStub != nil {
		return fake.SelfMembershipInfoStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.selfMembershipInfoReturns
	return fakeReturns.result1
}

func (fake *Gossip) SelfMembershipInfoCallCount() int {
	fake.selfMembershipInfoMutex.RLock()
	defer fake.selfMembershipInfoMutex.RUnlock()
	return len(fake.selfMembershipInfoArgsForCall)
}

func (fake *Gossip) SelfMembershipInfoCalls(stub func() discovery.NetworkMember) {
	fake.selfMembershipInfoMutex.Lock()
	defer fake.selfMembershipInfoMutex.Unlock()
	fake.SelfMembershipInfoStub = stub
}

func (fake *Gossip) SelfMembershipInfoReturns(result1 discovery.NetworkMember) {
	fake.selfMembershipInfoMutex.Lock()
	defer fake.selfMembershipInfoMutex.Unlock()
	fake.SelfMembershipInfoStub = nil
	fake.selfMembershipInfoReturns = struct {
		result1 discovery.NetworkMember
	}{result1}
}

func (fake *Gossip) SelfMembershipInfoReturnsOnCall(i int, result1 discovery.NetworkMember) {
	fake.selfMembershipInfoMutex.Lock()
	defer fake.selfMembershipInfoMutex.Unlock()
	fake.SelfMembershipInfoStub = nil
	if fake.selfMembershipInfoReturnsOnCall == nil {
		fake.selfMembershipInfoReturnsOnCall = make(map[int]struct {
			result1 discovery.NetworkMember
		})
	}
	fake.selfMembershipInfoReturnsOnCall[i] = struct {
		result1 discovery.NetworkMember
	}{result1}
}

func (fake *Gossip) Send(arg1 *gossipa.GossipMessage, arg2 ...*comm.RemotePeer) {
	fake.sendMutex.Lock()
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		arg1 *gossipa.GossipMessage
		arg2 []*comm.RemotePeer
	}{arg1, arg2})
	fake.recordInvocation("Send", []interface{}{arg1, arg2})
	fake.sendMutex.Unlock()
	if fake.SendStub != nil {
		fake.SendStub(arg1, arg2...)
	}
}

func (fake *Gossip) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *Gossip) SendCalls(stub func(*gossipa.GossipMessage, ...*comm.RemotePeer)) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = stub
}

func (fake *Gossip) SendArgsForCall(i int) (*gossipa.GossipMessage, []*comm.RemotePeer) {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	argsForCall := fake.sendArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Gossip) SendByCriteria(arg1 *protoext.SignedGossipMessage, arg2 gossip.SendCriteria) error {
	fake.sendByCriteriaMutex.Lock()
	ret, specificReturn := fake.sendByCriteriaReturnsOnCall[len(fake.sendByCriteriaArgsForCall)]
	fake.sendByCriteriaArgsForCall = append(fake.sendByCriteriaArgsForCall, struct {
		arg1 *protoext.SignedGossipMessage
		arg2 gossip.SendCriteria
	}{arg1, arg2})
	fake.recordInvocation("SendByCriteria", []interface{}{arg1, arg2})
	fake.sendByCriteriaMutex.Unlock()
	if fake.SendByCriteriaStub != nil {
		return fake.SendByCriteriaStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sendByCriteriaReturns
	return fakeReturns.result1
}

func (fake *Gossip) SendByCriteriaCallCount() int {
	fake.sendByCriteriaMutex.RLock()
	defer fake.sendByCriteriaMutex.RUnlock()
	return len(fake.sendByCriteriaArgsForCall)
}

func (fake *Gossip) SendByCriteriaCalls(stub func(*protoext.SignedGossipMessage, gossip.SendCriteria) error) {
	fake.sendByCriteriaMutex.Lock()
	defer fake.sendByCriteriaMutex.Unlock()
	fake.SendByCriteriaStub = stub
}

func (fake *Gossip) SendByCriteriaArgsForCall(i int) (*protoext.SignedGossipMessage, gossip.SendCriteria) {
	fake.sendByCriteriaMutex.RLock()
	defer fake.sendByCriteriaMutex.RUnlock()
	argsForCall := fake.sendByCriteriaArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Gossip) SendByCriteriaReturns(result1 error) {
	fake.sendByCriteriaMutex.Lock()
	defer fake.sendByCriteriaMutex.Unlock()
	fake.SendByCriteriaStub = nil
	fake.sendByCriteriaReturns = struct {
		result1 error
	}{result1}
}

func (fake *Gossip) SendByCriteriaReturnsOnCall(i int, result1 error) {
	fake.sendByCriteriaMutex.Lock()
	defer fake.sendByCriteriaMutex.Unlock()
	fake.SendByCriteriaStub = nil
	if fake.sendByCriteriaReturnsOnCall == nil {
		fake.sendByCriteriaReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendByCriteriaReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Gossip) Stop() {
	fake.stopMutex.Lock()
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct {
	}{})
	fake.recordInvocation("Stop", []interface{}{})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		fake.StopStub()
	}
}

func (fake *Gossip) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *Gossip) StopCalls(stub func()) {
	fake.stopMutex.Lock()
	defer fake.stopMutex.Unlock()
	fake.StopStub = stub
}

func (fake *Gossip) SuspectPeers(arg1 api.PeerSuspector) {
	fake.suspectPeersMutex.Lock()
	fake.suspectPeersArgsForCall = append(fake.suspectPeersArgsForCall, struct {
		arg1 api.PeerSuspector
	}{arg1})
	fake.recordInvocation("SuspectPeers", []interface{}{arg1})
	fake.suspectPeersMutex.Unlock()
	if fake.SuspectPeersStub != nil {
		fake.SuspectPeersStub(arg1)
	}
}

func (fake *Gossip) SuspectPeersCallCount() int {
	fake.suspectPeersMutex.RLock()
	defer fake.suspectPeersMutex.RUnlock()
	return len(fake.suspectPeersArgsForCall)
}

func (fake *Gossip) SuspectPeersCalls(stub func(api.PeerSuspector)) {
	fake.suspectPeersMutex.Lock()
	defer fake.suspectPeersMutex.Unlock()
	fake.SuspectPeersStub = stub
}

func (fake *Gossip) SuspectPeersArgsForCall(i int) api.PeerSuspector {
	fake.suspectPeersMutex.RLock()
	defer fake.suspectPeersMutex.RUnlock()
	argsForCall := fake.suspectPeersArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gossip) UpdateChaincodes(arg1 []*gossipa.Chaincode, arg2 common.ChainID) {
	var arg1Copy []*gossipa.Chaincode
	if arg1 != nil {
		arg1Copy = make([]*gossipa.Chaincode, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.updateChaincodesMutex.Lock()
	fake.updateChaincodesArgsForCall = append(fake.updateChaincodesArgsForCall, struct {
		arg1 []*gossipa.Chaincode
		arg2 common.ChainID
	}{arg1Copy, arg2})
	fake.recordInvocation("UpdateChaincodes", []interface{}{arg1Copy, arg2})
	fake.updateChaincodesMutex.Unlock()
	if fake.UpdateChaincodesStub != nil {
		fake.UpdateChaincodesStub(arg1, arg2)
	}
}

func (fake *Gossip) UpdateChaincodesCallCount() int {
	fake.updateChaincodesMutex.RLock()
	defer fake.updateChaincodesMutex.RUnlock()
	return len(fake.updateChaincodesArgsForCall)
}

func (fake *Gossip) UpdateChaincodesCalls(stub func([]*gossipa.Chaincode, common.ChainID)) {
	fake.updateChaincodesMutex.Lock()
	defer fake.updateChaincodesMutex.Unlock()
	fake.UpdateChaincodesStub = stub
}

func (fake *Gossip) UpdateChaincodesArgsForCall(i int) ([]*gossipa.Chaincode, common.ChainID) {
	fake.updateChaincodesMutex.RLock()
	defer fake.updateChaincodesMutex.RUnlock()
	argsForCall := fake.updateChaincodesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Gossip) UpdateLedgerHeight(arg1 uint64, arg2 common.ChainID) {
	fake.updateLedgerHeightMutex.Lock()
	fake.updateLedgerHeightArgsForCall = append(fake.updateLedgerHeightArgsForCall, struct {
		arg1 uint64
		arg2 common.ChainID
	}{arg1, arg2})
	fake.recordInvocation("UpdateLedgerHeight", []interface{}{arg1, arg2})
	fake.updateLedgerHeightMutex.Unlock()
	if fake.UpdateLedgerHeightStub != nil {
		fake.UpdateLedgerHeightStub(arg1, arg2)
	}
}

func (fake *Gossip) UpdateLedgerHeightCallCount() int {
	fake.updateLedgerHeightMutex.RLock()
	defer fake.updateLedgerHeightMutex.RUnlock()
	return len(fake.updateLedgerHeightArgsForCall)
}

func (fake *Gossip) UpdateLedgerHeightCalls(stub func(uint64, common.ChainID)) {
	fake.updateLedgerHeightMutex.Lock()
	defer fake.updateLedgerHeightMutex.Unlock()
	fake.UpdateLedgerHeightStub = stub
}

func (fake *Gossip) UpdateLedgerHeightArgsForCall(i int) (uint64, common.ChainID) {
	fake.updateLedgerHeightMutex.RLock()
	defer fake.updateLedgerHeightMutex.RUnlock()
	argsForCall := fake.updateLedgerHeightArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Gossip) UpdateMetadata(arg1 []byte) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.updateMetadataMutex.Lock()
	fake.updateMetadataArgsForCall = append(fake.updateMetadataArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("UpdateMetadata", []interface{}{arg1Copy})
	fake.updateMetadataMutex.Unlock()
	if fake.UpdateMetadataStub != nil {
		fake.UpdateMetadataStub(arg1)
	}
}

func (fake *Gossip) UpdateMetadataCallCount() int {
	fake.updateMetadataMutex.RLock()
	defer fake.updateMetadataMutex.RUnlock()
	return len(fake.updateMetadataArgsForCall)
}

func (fake *Gossip) UpdateMetadataCalls(stub func([]byte)) {
	fake.updateMetadataMutex.Lock()
	defer fake.updateMetadataMutex.Unlock()
	fake.UpdateMetadataStub = stub
}

func (fake *Gossip) UpdateMetadataArgsForCall(i int) []byte {
	fake.updateMetadataMutex.RLock()
	defer fake.updateMetadataMutex.RUnlock()
	argsForCall := fake.updateMetadataArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gossip) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.acceptMutex.RLock()
	defer fake.acceptMutex.RUnlock()
	fake.gossipMutex.RLock()
	defer fake.gossipMutex.RUnlock()
	fake.identityInfoMutex.RLock()
	defer fake.identityInfoMutex.RUnlock()
	fake.isInMyOrgMutex.RLock()
	defer fake.isInMyOrgMutex.RUnlock()
	fake.joinChanMutex.RLock()
	defer fake.joinChanMutex.RUnlock()
	fake.leaveChanMutex.RLock()
	defer fake.leaveChanMutex.RUnlock()
	fake.peerFilterMutex.RLock()
	defer fake.peerFilterMutex.RUnlock()
	fake.peersMutex.RLock()
	defer fake.peersMutex.RUnlock()
	fake.peersOfChannelMutex.RLock()
	defer fake.peersOfChannelMutex.RUnlock()
	fake.selfChannelInfoMutex.RLock()
	defer fake.selfChannelInfoMutex.RUnlock()
	fake.selfMembershipInfoMutex.RLock()
	defer fake.selfMembershipInfoMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	fake.sendByCriteriaMutex.RLock()
	defer fake.sendByCriteriaMutex.RUnlock()
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	fake.suspectPeersMutex.RLock()
	defer fake.suspectPeersMutex.RUnlock()
	fake.updateChaincodesMutex.RLock()
	defer fake.updateChaincodesMutex.RUnlock()
	fake.updateLedgerHeightMutex.RLock()
	defer fake.updateLedgerHeightMutex.RUnlock()
	fake.updateMetadataMutex.RLock()
	defer fake.updateMetadataMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Gossip) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ gossip.Gossip = new(Gossip)
