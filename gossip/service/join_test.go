/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"sync"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/filter"
	"github.com/mcc-github/blockchain/gossip/gossip"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type secAdvMock struct {
}

func init() {
	util.SetupTestLogging()
}

func (s *secAdvMock) OrgByPeerIdentity(identity api.PeerIdentityType) api.OrgIdentityType {
	return api.OrgIdentityType(identity)
}

type gossipMock struct {
	mock.Mock
}

func (g *gossipMock) SelfChannelInfo(common.ChannelID) *protoext.SignedGossipMessage {
	panic("implement me")
}

func (g *gossipMock) SelfMembershipInfo() discovery.NetworkMember {
	panic("implement me")
}

func (*gossipMock) PeerFilter(channel common.ChannelID, messagePredicate api.SubChannelSelectionCriteria) (filter.RoutingFilter, error) {
	panic("implement me")
}

func (*gossipMock) SuspectPeers(s api.PeerSuspector) {
	panic("implement me")
}

func (*gossipMock) Send(msg *proto.GossipMessage, peers ...*comm.RemotePeer) {
	panic("implement me")
}

func (*gossipMock) Peers() []discovery.NetworkMember {
	panic("implement me")
}

func (*gossipMock) PeersOfChannel(common.ChannelID) []discovery.NetworkMember {
	panic("implement me")
}

func (*gossipMock) UpdateMetadata(metadata []byte) {
	panic("implement me")
}



func (*gossipMock) UpdateLedgerHeight(height uint64, channelID common.ChannelID) {
	panic("implement me")
}



func (*gossipMock) UpdateChaincodes(chaincode []*proto.Chaincode, channelID common.ChannelID) {
	panic("implement me")
}

func (*gossipMock) Gossip(msg *proto.GossipMessage) {
	panic("implement me")
}

func (*gossipMock) Accept(acceptor common.MessageAcceptor, passThrough bool) (<-chan *proto.GossipMessage, <-chan protoext.ReceivedMessage) {
	panic("implement me")
}

func (g *gossipMock) JoinChan(joinMsg api.JoinChannelMessage, channelID common.ChannelID) {
	g.Called(joinMsg, channelID)
}

func (g *gossipMock) LeaveChan(channelID common.ChannelID) {
	panic("implement me")
}

func (g *gossipMock) IdentityInfo() api.PeerIdentitySet {
	panic("implement me")
}

func (*gossipMock) IsInMyOrg(member discovery.NetworkMember) bool {
	panic("implement me")
}

func (*gossipMock) Stop() {
	panic("implement me")
}

func (*gossipMock) SendByCriteria(*protoext.SignedGossipMessage, gossip.SendCriteria) error {
	panic("implement me")
}

type appOrgMock struct {
	id string
}

func (*appOrgMock) Name() string {
	panic("implement me")
}

func (ao *appOrgMock) MSPID() string {
	return ao.id
}

func (ao *appOrgMock) AnchorPeers() []*peer.AnchorPeer {
	return []*peer.AnchorPeer{}
}

type configMock struct {
	orgs2AppOrgs map[string]channelconfig.ApplicationOrg
}

func (c *configMock) OrdererAddresses() []string {
	return []string{"localhost:7050"}
}

func (*configMock) ChannelID() string {
	return "A"
}

func (c *configMock) Organizations() map[string]channelconfig.ApplicationOrg {
	return c.orgs2AppOrgs
}

func (*configMock) Sequence() uint64 {
	return 0
}

func TestJoinChannelConfig(t *testing.T) {
	
	
	

	failChan := make(chan struct{}, 1)
	g1SvcMock := &gossipMock{}
	g1SvcMock.On("JoinChan", mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		failChan <- struct{}{}
	})
	g1 := &GossipService{secAdv: &secAdvMock{}, peerIdentity: api.PeerIdentityType("OrgMSP0"), gossipSvc: g1SvcMock}
	g1.updateAnchors(&configMock{
		orgs2AppOrgs: map[string]channelconfig.ApplicationOrg{
			"Org0": &appOrgMock{id: "Org0"},
		},
	})
	select {
	case <-time.After(time.Second):
	case <-failChan:
		assert.Fail(t, "Joined a badly configured channel")
	}

	succChan := make(chan struct{}, 1)
	g2SvcMock := &gossipMock{}
	g2SvcMock.On("JoinChan", mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		succChan <- struct{}{}
	})
	g2 := &GossipService{secAdv: &secAdvMock{}, peerIdentity: api.PeerIdentityType("Org0"), gossipSvc: g2SvcMock}
	g2.updateAnchors(&configMock{
		orgs2AppOrgs: map[string]channelconfig.ApplicationOrg{
			"Org0": &appOrgMock{id: "Org0"},
		},
	})
	select {
	case <-time.After(time.Second):
		assert.Fail(t, "Didn't join a channel (should have done so within the time period)")
	case <-succChan:

	}
}

func TestJoinChannelNoAnchorPeers(t *testing.T) {
	
	
	

	var joinChanCalled sync.WaitGroup
	joinChanCalled.Add(1)
	gMock := &gossipMock{}
	gMock.On("JoinChan", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		defer joinChanCalled.Done()
		jcm := args.Get(0).(api.JoinChannelMessage)
		channel := args.Get(1).(common.ChannelID)
		assert.Len(t, jcm.Members(), 2)
		assert.Contains(t, jcm.Members(), api.OrgIdentityType("Org0"))
		assert.Contains(t, jcm.Members(), api.OrgIdentityType("Org1"))
		assert.Equal(t, "A", string(channel))
	})

	g := &GossipService{secAdv: &secAdvMock{}, peerIdentity: api.PeerIdentityType("Org0"), gossipSvc: gMock}

	appOrg0 := &appOrgMock{id: "Org0"}
	appOrg1 := &appOrgMock{id: "Org1"}

	
	assert.Empty(t, appOrg0.AnchorPeers())
	assert.Empty(t, appOrg1.AnchorPeers())

	g.updateAnchors(&configMock{
		orgs2AppOrgs: map[string]channelconfig.ApplicationOrg{
			"Org0": appOrg0,
			"Org1": appOrg1,
		},
	})
	joinChanCalled.Wait()
}
