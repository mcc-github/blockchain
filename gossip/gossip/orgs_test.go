/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	cb "github.com/mcc-github/blockchain-protos-go/common"
	proto "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/gossip/api"
	gcomm "github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/gossip/algo"
	"github.com/mcc-github/blockchain/gossip/gossip/channel"
	"github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	util.SetupTestLogging()
	factory.InitFactories(nil)
}

type configurableCryptoService struct {
	sync.RWMutex
	m map[string]api.OrgIdentityType
}

func (c *configurableCryptoService) Expiration(peerIdentity api.PeerIdentityType) (time.Time, error) {
	return time.Now().Add(time.Hour), nil
}

func (c *configurableCryptoService) putInOrg(port int, org string) {
	identity := fmt.Sprintf("127.0.0.1:%d", port)
	c.Lock()
	c.m[identity] = api.OrgIdentityType(org)
	c.Unlock()
}



func (c *configurableCryptoService) OrgByPeerIdentity(identity api.PeerIdentityType) api.OrgIdentityType {
	c.RLock()
	org := c.m[string(identity)]
	c.RUnlock()
	return org
}



func (c *configurableCryptoService) VerifyByChannel(_ common.ChannelID, identity api.PeerIdentityType, _, _ []byte) error {
	return nil
}

func (*configurableCryptoService) ValidateIdentity(peerIdentity api.PeerIdentityType) error {
	return nil
}


func (*configurableCryptoService) GetPKIidOfCert(peerIdentity api.PeerIdentityType) common.PKIidType {
	return common.PKIidType(peerIdentity)
}



func (*configurableCryptoService) VerifyBlock(channelID common.ChannelID, seqNum uint64, signedBlock *cb.Block) error {
	return nil
}



func (*configurableCryptoService) Sign(msg []byte) ([]byte, error) {
	sig := make([]byte, len(msg))
	copy(sig, msg)
	return sig, nil
}




func (*configurableCryptoService) Verify(peerIdentity api.PeerIdentityType, signature, message []byte) error {
	equal := bytes.Equal(signature, message)
	if !equal {
		return fmt.Errorf("Wrong signature:%v, %v", signature, message)
	}
	return nil
}

func newGossipInstanceWithGRPCWithExternalEndpoint(id int, port int, gRPCServer *comm.GRPCServer,
	certs *common.TLSCertificates, secureDialOpts api.PeerSecureDialOpts, mcs *configurableCryptoService,
	externalEndpoint string, boot ...int) *gossipGRPC {
	conf := &Config{
		BootstrapPeers:               bootPeersWithPorts(boot...),
		ID:                           fmt.Sprintf("p%d", id),
		MaxBlockCountToStore:         100,
		MaxPropagationBurstLatency:   time.Duration(500) * time.Millisecond,
		MaxPropagationBurstSize:      20,
		PropagateIterations:          1,
		PropagatePeerNum:             3,
		PullInterval:                 time.Duration(2) * time.Second,
		PullPeerNum:                  5,
		InternalEndpoint:             fmt.Sprintf("127.0.0.1:%d", port),
		ExternalEndpoint:             externalEndpoint,
		PublishCertPeriod:            time.Duration(4) * time.Second,
		PublishStateInfoInterval:     time.Duration(1) * time.Second,
		RequestStateInfoInterval:     time.Duration(1) * time.Second,
		TimeForMembershipTracker:     5 * time.Second,
		TLSCerts:                     certs,
		DigestWaitTime:               algo.DefDigestWaitTime,
		RequestWaitTime:              algo.DefRequestWaitTime,
		ResponseWaitTime:             algo.DefResponseWaitTime,
		DialTimeout:                  gcomm.DefDialTimeout,
		ConnTimeout:                  gcomm.DefConnTimeout,
		RecvBuffSize:                 gcomm.DefRecvBuffSize,
		SendBuffSize:                 gcomm.DefSendBuffSize,
		MsgExpirationTimeout:         channel.DefMsgExpirationTimeout,
		AliveTimeInterval:            discoveryConfig.AliveTimeInterval,
		AliveExpirationTimeout:       discoveryConfig.AliveExpirationTimeout,
		AliveExpirationCheckInterval: discoveryConfig.AliveExpirationCheckInterval,
		ReconnectInterval:            discoveryConfig.ReconnectInterval,
	}
	selfID := api.PeerIdentityType(conf.InternalEndpoint)
	g := New(conf, gRPCServer.Server(), mcs, mcs, selfID,
		secureDialOpts, metrics.NewGossipMetrics(&disabled.Provider{}))
	go func() {
		gRPCServer.Start()
	}()
	return &gossipGRPC{GossipImpl: g, grpc: gRPCServer}
}

func TestMultipleOrgEndpointLeakage(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	
	
	
	
	cs := &configurableCryptoService{m: make(map[string]api.OrgIdentityType)}
	peersInOrg := 5
	orgA := "orgA"
	orgB := "orgB"
	channel := common.ChannelID("TEST")
	orgs := []string{orgA, orgB}
	peers := []*gossipGRPC{}

	expectedMembershipSize := map[string]int{}
	peersWithExternalEndpoints := make(map[string]struct{})

	shouldAKnowB := func(a common.PKIidType, b common.PKIidType) bool {
		orgOfPeerA := cs.OrgByPeerIdentity(api.PeerIdentityType(a))
		orgOfPeerB := cs.OrgByPeerIdentity(api.PeerIdentityType(b))
		_, aHasExternalEndpoint := peersWithExternalEndpoints[string(a)]
		_, bHasExternalEndpoint := peersWithExternalEndpoints[string(b)]
		bothHaveExternalEndpoints := aHasExternalEndpoint && bHasExternalEndpoint
		return bytes.Equal(orgOfPeerA, orgOfPeerB) || bothHaveExternalEndpoints
	}

	shouldKnowInternalEndpoint := func(a common.PKIidType, b common.PKIidType) bool {
		orgOfPeerA := cs.OrgByPeerIdentity(api.PeerIdentityType(a))
		orgOfPeerB := cs.OrgByPeerIdentity(api.PeerIdentityType(b))
		return bytes.Equal(orgOfPeerA, orgOfPeerB)
	}

	var ports []int
	var grpcs []*comm.GRPCServer
	var certs []*common.TLSCertificates
	var secDialOpts []api.PeerSecureDialOpts

	for range orgs {
		for i := 0; i < peersInOrg; i++ {
			port, grpc, cert, secDialOpt, _ := util.CreateGRPCLayer()
			ports = append(ports, port)
			grpcs = append(grpcs, grpc)
			certs = append(certs, cert)
			secDialOpts = append(secDialOpts, secDialOpt)
		}
	}

	for orgIndex, org := range orgs {
		for i := 0; i < peersInOrg; i++ {
			id := orgIndex*peersInOrg + i
			endpoint := fmt.Sprintf("127.0.0.1:%d", ports[id])
			cs.putInOrg(ports[id], org)
			expectedMembershipSize[endpoint] = peersInOrg - 1 
			externalEndpoint := ""
			if i < 2 {
				
				
				
				externalEndpoint = endpoint
				peersWithExternalEndpoints[externalEndpoint] = struct{}{}
				expectedMembershipSize[endpoint] += 2 
			}
			peer := newGossipInstanceWithGRPCWithExternalEndpoint(id, ports[id], grpcs[id], certs[id], secDialOpts[id],
				cs, externalEndpoint)
			peers = append(peers, peer)
		}
	}

	jcm := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			orgA: {
				{Host: "127.0.0.1", Port: ports[1]},
			},
			orgB: {
				{Host: "127.0.0.1", Port: ports[5]},
			},
		},
	}

	for _, p := range peers {
		p.JoinChan(jcm, channel)
		p.UpdateLedgerHeight(1, channel)
	}

	membershipCheck := func() bool {
		for _, p := range peers {
			peerNetMember := p.GossipImpl.selfNetworkMember()
			pkiID := peerNetMember.PKIid
			peersKnown := p.Peers()
			peersToKnow := expectedMembershipSize[string(pkiID)]
			if peersToKnow != len(peersKnown) {
				t.Logf("peer %#v doesn't know the needed amount of peers, extected %#v, actual %#v", peerNetMember.Endpoint, peersToKnow, len(peersKnown))
				return false
			}
			for _, knownPeer := range peersKnown {
				if !shouldAKnowB(pkiID, knownPeer.PKIid) {
					assert.Fail(t, fmt.Sprintf("peer %#v doesn't know %#v", peerNetMember.Endpoint, knownPeer.Endpoint))
					return false
				}
				internalEndpointLen := len(knownPeer.InternalEndpoint)
				if shouldKnowInternalEndpoint(pkiID, knownPeer.PKIid) {
					if internalEndpointLen == 0 {
						t.Logf("peer: %v doesn't know internal endpoint of %v", peerNetMember.InternalEndpoint, string(knownPeer.PKIid))
						return false
					}
				} else {
					if internalEndpointLen != 0 {
						assert.Fail(t, fmt.Sprintf("peer: %v knows internal endpoint of %v (%#v)", peerNetMember.InternalEndpoint, string(knownPeer.PKIid), knownPeer.InternalEndpoint))
						return false
					}
				}
			}
		}
		return true
	}

	waitUntilOrFail(t, membershipCheck, "waiting for all instances to form membership view")

	for _, p := range peers {
		p.Stop()
	}
}

func TestConfidentiality(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	peersInOrg := 3
	externalEndpointsInOrg := 2

	
	
	
	
	peersWithExternalEndpoints := map[string]struct{}{}

	orgs := []string{"A", "B", "C", "D"}
	channels := []string{"C0", "C1", "C2", "C3"}
	isOrgInChan := func(org string, channel string) bool {
		switch org {
		case "A":
			return channel == "C0" || channel == "C1"
		case "B":
			return channel == "C0" || channel == "C2" || channel == "C3"
		case "C":
			return channel == "C1" || channel == "C2"
		case "D":
			return channel == "C3"
		}

		return false
	}

	var ports []int
	var grpcs []*comm.GRPCServer
	var certs []*common.TLSCertificates
	var secDialOpts []api.PeerSecureDialOpts

	for range orgs {
		for j := 0; j < peersInOrg; j++ {
			port, grpc, cert, secDialOpt, _ := util.CreateGRPCLayer()
			ports = append(ports, port)
			grpcs = append(grpcs, grpc)
			certs = append(certs, cert)
			secDialOpts = append(secDialOpts, secDialOpt)
		}
	}

	
	cs := &configurableCryptoService{m: make(map[string]api.OrgIdentityType)}
	for i, org := range orgs {
		for j := 0; j < peersInOrg; j++ {
			port := ports[i*peersInOrg+j]
			cs.putInOrg(port, org)
		}
	}

	var peers []*gossipGRPC
	orgs2Peers := map[string][]*gossipGRPC{
		"A": {},
		"B": {},
		"C": {},
		"D": {},
	}

	anchorPeersByOrg := map[string]api.AnchorPeer{}

	for i, org := range orgs {
		for j := 0; j < peersInOrg; j++ {
			id := i*peersInOrg + j
			endpoint := fmt.Sprintf("127.0.0.1:%d", ports[id])
			externalEndpoint := ""
			if j < externalEndpointsInOrg { 
				externalEndpoint = endpoint
				peersWithExternalEndpoints[string(endpoint)] = struct{}{}
			}
			peer := newGossipInstanceWithGRPCWithExternalEndpoint(id, ports[id], grpcs[id], certs[id], secDialOpts[id],
				cs, externalEndpoint)
			peers = append(peers, peer)
			orgs2Peers[org] = append(orgs2Peers[org], peer)
			t.Log(endpoint, "id:", id, "externalEndpoint:", externalEndpoint)
			
			if j == 0 {
				anchorPeersByOrg[org] = api.AnchorPeer{
					Host: "127.0.0.1",
					Port: ports[id],
				}
			}
		}
	}

	msgs2Inspect := make(chan *msg, 3000)
	defer close(msgs2Inspect)
	go inspectMsgs(t, msgs2Inspect, cs, peersWithExternalEndpoints)
	finished := int32(0)
	var wg sync.WaitGroup

	msgSelector := func(o interface{}) bool {
		msg := o.(protoext.ReceivedMessage).GetGossipMessage()
		identitiesPull := protoext.IsPullMsg(msg.GossipMessage) && protoext.GetPullMsgType(msg.GossipMessage) == proto.PullMsgType_IDENTITY_MSG
		return protoext.IsAliveMsg(msg.GossipMessage) || protoext.IsStateInfoMsg(msg.GossipMessage) || protoext.IsStateInfoSnapshot(msg.GossipMessage) || msg.GetMemRes() != nil || identitiesPull
	}
	
	
	for _, p := range peers {
		wg.Add(1)
		_, msgs := p.Accept(msgSelector, true)
		peerNetMember := p.GossipImpl.selfNetworkMember()
		targetORg := string(cs.OrgByPeerIdentity(api.PeerIdentityType(peerNetMember.InternalEndpoint)))
		go func(targetOrg string, msgs <-chan protoext.ReceivedMessage) {
			defer wg.Done()
			for receivedMsg := range msgs {
				m := &msg{
					src:           string(cs.OrgByPeerIdentity(receivedMsg.GetConnectionInfo().Identity)),
					dst:           targetORg,
					GossipMessage: receivedMsg.GetGossipMessage().GossipMessage,
				}
				if atomic.LoadInt32(&finished) == int32(1) {
					return
				}
				msgs2Inspect <- m
			}
		}(targetORg, msgs)
	}

	
	joinChanMsgsByChan := map[string]*joinChanMsg{}
	for _, ch := range channels {
		jcm := &joinChanMsg{members2AnchorPeers: map[string][]api.AnchorPeer{}}
		for _, org := range orgs {
			if isOrgInChan(org, ch) {
				jcm.members2AnchorPeers[org] = append(jcm.members2AnchorPeers[org], anchorPeersByOrg[org])
			}
		}
		joinChanMsgsByChan[ch] = jcm
	}

	
	for org, peers := range orgs2Peers {
		for _, ch := range channels {
			if isOrgInChan(org, ch) {
				for _, p := range peers {
					p.JoinChan(joinChanMsgsByChan[ch], common.ChannelID(ch))
					p.UpdateLedgerHeight(1, common.ChannelID(ch))
					go func(p *gossipGRPC, ch string) {
						for i := 0; i < 5; i++ {
							time.Sleep(time.Second)
							p.UpdateLedgerHeight(1, common.ChannelID(ch))
						}
					}(p, ch)
				}
			}
		}
	}

	
	time.Sleep(time.Second * 7)

	assertMembership := func() bool {
		for _, org := range orgs {
			for i, p := range orgs2Peers[org] {
				members := p.Peers()
				expMemberSize := expectedMembershipSize(peersInOrg, externalEndpointsInOrg, org, i < externalEndpointsInOrg)
				peerNetMember := p.GossipImpl.selfNetworkMember()
				membersCount := len(members)
				if membersCount < expMemberSize {
					return false
				}
				
				assert.True(t, membersCount <= expMemberSize, "%s knows too much (%d > %d) peers: %v",
					membersCount, expMemberSize, peerNetMember.PKIid, members)
			}
		}
		return true
	}

	waitUntilOrFail(t, assertMembership, "waiting for all instances to form unified membership view")
	stopPeers(peers)
	wg.Wait()
	atomic.StoreInt32(&finished, int32(1))
}

func expectedMembershipSize(peersInOrg, externalEndpointsInOrg int, org string, hasExternalEndpoint bool) int {
	
	
	
	
	
	
	
	
	

	m := map[string]func(x, y int) int{
		"A": func(x, y int) int {
			return x + 2*y
		},
		"B": func(x, y int) int {
			return x + 3*y
		},
		"C": func(x, y int) int {
			return x + 2*y
		},
		"D": func(x, y int) int {
			return x + y
		},
	}

	
	
	if !hasExternalEndpoint {
		externalEndpointsInOrg = 0
	}
	
	return m[org](peersInOrg, externalEndpointsInOrg) - 1
}

func extractOrgsFromMsg(msg *proto.GossipMessage, sec api.SecurityAdvisor) []string {
	if protoext.IsAliveMsg(msg) {
		return []string{string(sec.OrgByPeerIdentity(api.PeerIdentityType(msg.GetAliveMsg().Membership.PkiId)))}
	}

	orgs := map[string]struct{}{}

	if protoext.IsPullMsg(msg) {
		if protoext.IsDigestMsg(msg) || protoext.IsDataReq(msg) {
			var digests []string
			if protoext.IsDigestMsg(msg) {
				digests = util.BytesToStrings(msg.GetDataDig().Digests)
			} else {
				digests = util.BytesToStrings(msg.GetDataReq().Digests)
			}

			for _, dig := range digests {
				org := sec.OrgByPeerIdentity(api.PeerIdentityType(dig))
				orgs[string(org)] = struct{}{}
			}
		}

		if protoext.IsDataUpdate(msg) {
			for _, identityMsg := range msg.GetDataUpdate().Data {
				gMsg, _ := protoext.EnvelopeToGossipMessage(identityMsg)
				id := string(gMsg.GetPeerIdentity().Cert)
				org := sec.OrgByPeerIdentity(api.PeerIdentityType(id))
				orgs[string(org)] = struct{}{}
			}
		}
	}

	if msg.GetMemRes() != nil {
		alive := msg.GetMemRes().Alive
		dead := msg.GetMemRes().Dead
		for _, envp := range append(alive, dead...) {
			msg, _ := protoext.EnvelopeToGossipMessage(envp)
			orgs[string(sec.OrgByPeerIdentity(api.PeerIdentityType(msg.GetAliveMsg().Membership.PkiId)))] = struct{}{}
		}
	}

	res := []string{}
	for org := range orgs {
		res = append(res, org)
	}
	return res
}

func inspectMsgs(t *testing.T, msgChan chan *msg, sec api.SecurityAdvisor, peersWithExternalEndpoints map[string]struct{}) {
	for msg := range msgChan {
		
		
		if msg.src == msg.dst {
			continue
		}
		if protoext.IsStateInfoMsg(msg.GossipMessage) || protoext.IsStateInfoSnapshot(msg.GossipMessage) {
			inspectStateInfoMsg(t, msg, peersWithExternalEndpoints)
			continue
		}
		
		
		
		orgs := extractOrgsFromMsg(msg.GossipMessage, sec)
		s := []string{msg.src, msg.dst}
		assert.True(t, isSubset(orgs, s), "%v isn't a subset of %v", orgs, s)

		
		if msg.dst == "D" {
			assert.NotContains(t, "A", orgs)
			assert.NotContains(t, "C", orgs)
		}

		if msg.dst == "A" || msg.dst == "C" {
			assert.NotContains(t, "D", orgs)
		}

		
		
		isIdentityPull := protoext.IsPullMsg(msg.GossipMessage) && protoext.GetPullMsgType(msg.GossipMessage) == proto.PullMsgType_IDENTITY_MSG
		if !(isIdentityPull && protoext.IsDataUpdate(msg.GossipMessage)) {
			continue
		}
		for _, envp := range msg.GetDataUpdate().Data {
			identityMsg, _ := protoext.EnvelopeToGossipMessage(envp)
			pkiID := identityMsg.GetPeerIdentity().PkiId
			_, hasExternalEndpoint := peersWithExternalEndpoints[string(pkiID)]
			assert.True(t, hasExternalEndpoint,
				"Peer %s doesn't have an external endpoint but its identity was gossiped", string(pkiID))
		}
	}
}

func inspectStateInfoMsg(t *testing.T, m *msg, peersWithExternalEndpoints map[string]struct{}) {
	if protoext.IsStateInfoMsg(m.GossipMessage) {
		pkiID := m.GetStateInfo().PkiId
		_, hasExternalEndpoint := peersWithExternalEndpoints[string(pkiID)]
		assert.True(t, hasExternalEndpoint, "peer %s has no external endpoint but crossed an org", string(pkiID))
		return
	}

	for _, envp := range m.GetStateSnapshot().Elements {
		msg, _ := protoext.EnvelopeToGossipMessage(envp)
		pkiID := msg.GetStateInfo().PkiId
		_, hasExternalEndpoint := peersWithExternalEndpoints[string(pkiID)]
		assert.True(t, hasExternalEndpoint, "peer %s has no external endpoint but crossed an org", string(pkiID))
	}
}

type msg struct {
	src string
	dst string
	*proto.GossipMessage
}

func isSubset(a []string, b []string) bool {
	for _, s1 := range a {
		found := false
		for _, s2 := range b {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
