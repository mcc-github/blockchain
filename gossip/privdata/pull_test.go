/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"bytes"
	"crypto/rand"
	"sync"
	"testing"

	pb "github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/filter"
	"github.com/mcc-github/blockchain/gossip/privdata/mocks"
	"github.com/mcc-github/blockchain/gossip/util"
	fcommon "github.com/mcc-github/blockchain/protos/common"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	logging.SetLevel(logging.DEBUG, util.LoggingPrivModule)
	policy2Filter = make(map[privdata.CollectionAccessPolicy]privdata.Filter)
}



func protoMatcher(pvds ...*proto.PvtDataDigest) func([]*proto.PvtDataDigest) bool {
	return func(ipvds []*proto.PvtDataDigest) bool {
		if len(pvds) != len(ipvds) {
			return false
		}

		for i, pvd := range pvds {
			if !pb.Equal(pvd, ipvds[i]) {
				return false
			}
		}

		return true
	}
}

var policyLock sync.Mutex
var policy2Filter map[privdata.CollectionAccessPolicy]privdata.Filter

type mockCollectionStore struct {
	m map[string]*mockCollectionAccess
}

func newCollectionStore() *mockCollectionStore {
	return &mockCollectionStore{
		m: make(map[string]*mockCollectionAccess),
	}
}

func (cs *mockCollectionStore) withPolicy(collection string, btl uint64) *mockCollectionAccess {
	coll := &mockCollectionAccess{cs: cs, btl: btl}
	cs.m[collection] = coll
	return coll
}

func (cs mockCollectionStore) RetrieveCollectionAccessPolicy(cc fcommon.CollectionCriteria) (privdata.CollectionAccessPolicy, error) {
	return cs.m[cc.Collection], nil
}

func (cs mockCollectionStore) RetrieveCollection(fcommon.CollectionCriteria) (privdata.Collection, error) {
	panic("implement me")
}

func (cs mockCollectionStore) RetrieveCollectionConfigPackage(fcommon.CollectionCriteria) (*fcommon.CollectionConfigPackage, error) {
	panic("implement me")
}

func (cs mockCollectionStore) RetrieveCollectionPersistenceConfigs(cc fcommon.CollectionCriteria) (privdata.CollectionPersistenceConfigs, error) {
	return cs.m[cc.Collection], nil
}

type mockCollectionAccess struct {
	cs  *mockCollectionStore
	btl uint64
}

func (mc *mockCollectionAccess) BlockToLive() uint64 {
	return mc.btl
}

func (mc *mockCollectionAccess) thatMapsTo(peers ...string) *mockCollectionStore {
	policyLock.Lock()
	defer policyLock.Unlock()
	policy2Filter[mc] = func(sd fcommon.SignedData) bool {
		for _, peer := range peers {
			if bytes.Equal(sd.Identity, []byte(peer)) {
				return true
			}
		}
		return false
	}
	return mc.cs
}

func (mc *mockCollectionAccess) MemberOrgs() []string {
	return nil
}

func (mc *mockCollectionAccess) AccessFilter() privdata.Filter {
	policyLock.Lock()
	defer policyLock.Unlock()
	return policy2Filter[mc]
}

func (mc *mockCollectionAccess) RequiredPeerCount() int {
	return 0
}

func (mc *mockCollectionAccess) MaximumPeerCount() int {
	return 0
}

type dataRetrieverMock struct {
	mock.Mock
}

func (dr *dataRetrieverMock) CollectionRWSet(dig []*proto.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, error) {
	args := dr.Called(dig, blockNum)
	return args.Get(0).(Dig2PvtRWSetWithConfig), args.Error(1)
}

type receivedMsg struct {
	responseChan chan proto.ReceivedMessage
	*comm.RemotePeer
	*proto.SignedGossipMessage
}

func (msg *receivedMsg) Ack(_ error) {

}

func (msg *receivedMsg) Respond(message *proto.GossipMessage) {
	m, _ := message.NoopSign()
	msg.responseChan <- &receivedMsg{SignedGossipMessage: m, RemotePeer: &comm.RemotePeer{}}
}

func (msg *receivedMsg) GetGossipMessage() *proto.SignedGossipMessage {
	return msg.SignedGossipMessage
}

func (msg *receivedMsg) GetSourceEnvelope() *proto.Envelope {
	panic("implement me")
}

func (msg *receivedMsg) GetConnectionInfo() *proto.ConnectionInfo {
	return &proto.ConnectionInfo{
		Identity: api.PeerIdentityType(msg.RemotePeer.PKIID),
		Auth: &proto.AuthInfo{
			SignedData: []byte{},
			Signature:  []byte{},
		},
	}
}

type mockGossip struct {
	mock.Mock
	msgChan chan proto.ReceivedMessage
	id      *comm.RemotePeer
	network *gossipNetwork
}

func newMockGossip(id *comm.RemotePeer) *mockGossip {
	return &mockGossip{
		msgChan: make(chan proto.ReceivedMessage),
		id:      id,
	}
}

func (g *mockGossip) PeerFilter(channel common.ChainID, messagePredicate api.SubChannelSelectionCriteria) (filter.RoutingFilter, error) {
	for _, call := range g.Mock.ExpectedCalls {
		if call.Method == "PeerFilter" {
			args := g.Called(channel, messagePredicate)
			if args.Get(1) != nil {
				return nil, args.Get(1).(error)
			}
			return args.Get(0).(filter.RoutingFilter), nil
		}
	}
	return func(member discovery.NetworkMember) bool {
		return messagePredicate(api.PeerSignature{
			PeerIdentity: api.PeerIdentityType(member.PKIid),
		})
	}, nil
}

func (g *mockGossip) Send(msg *proto.GossipMessage, peers ...*comm.RemotePeer) {
	sMsg, _ := msg.NoopSign()
	for _, peer := range g.network.peers {
		if bytes.Equal(peer.id.PKIID, peers[0].PKIID) {
			peer.msgChan <- &receivedMsg{
				RemotePeer:          g.id,
				SignedGossipMessage: sMsg,
				responseChan:        g.msgChan,
			}
			return
		}
	}
}

func (g *mockGossip) PeersOfChannel(common.ChainID) []discovery.NetworkMember {
	return g.Called().Get(0).([]discovery.NetworkMember)
}

func (g *mockGossip) Accept(acceptor common.MessageAcceptor, passThrough bool) (<-chan *proto.GossipMessage, <-chan proto.ReceivedMessage) {
	return nil, g.msgChan
}

type peerData struct {
	id           string
	ledgerHeight uint64
}

func membership(knownPeers ...peerData) []discovery.NetworkMember {
	var peers []discovery.NetworkMember
	for _, peer := range knownPeers {
		peers = append(peers, discovery.NetworkMember{
			Endpoint: peer.id,
			PKIid:    common.PKIidType(peer.id),
			Properties: &proto.Properties{
				LedgerHeight: peer.ledgerHeight,
			},
		})
	}
	return peers
}

type gossipNetwork struct {
	peers []*mockGossip
}

func (gn *gossipNetwork) newPuller(id string, ps privdata.CollectionStore, factory CollectionAccessFactory, knownMembers ...discovery.NetworkMember) *puller {
	g := newMockGossip(&comm.RemotePeer{PKIID: common.PKIidType(id), Endpoint: id})
	g.network = gn
	g.On("PeersOfChannel", mock.Anything).Return(knownMembers)

	p := NewPuller(ps, g, &dataRetrieverMock{}, factory, "A")
	gn.peers = append(gn.peers, g)
	return p
}

func newPRWSet() []util.PrivateRWSet {
	b1 := make([]byte, 10)
	b2 := make([]byte, 10)
	rand.Read(b1)
	rand.Read(b2)
	return []util.PrivateRWSet{util.PrivateRWSet(b1), util.PrivateRWSet(b2)}
}

func TestPullerFromOnly1Peer(t *testing.T) {
	t.Parallel()
	
	
	
	gn := &gossipNetwork{}
	policyStore := newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p2")
	factoryMock1 := &collectionAccessFactoryMock{}
	policyMock1 := &collectionAccessPolicyMock{}
	policyMock1.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p2"))
	}, []string{"org1", "org2"})
	factoryMock1.On("AccessPolicy", mock.Anything, mock.Anything).Return(policyMock1, nil)
	p1 := gn.newPuller("p1", policyStore, factoryMock1, membership(peerData{"p2", uint64(1)}, peerData{"p3", uint64(1)})...)

	p2TransientStore := &util.PrivateRWSetWithConfig{
		RWSet: newPRWSet(),
		CollectionConfig: &fcommon.CollectionConfig{
			Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &fcommon.StaticCollectionConfig{
					Name: "col1",
				},
			},
		},
	}
	policyStore = newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p1")
	factoryMock2 := &collectionAccessFactoryMock{}
	policyMock2 := &collectionAccessPolicyMock{}
	policyMock2.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p1"))
	}, []string{"org1", "org2"})
	factoryMock2.On("AccessPolicy", mock.Anything, mock.Anything).Return(policyMock2, nil)

	p2 := gn.newPuller("p2", policyStore, factoryMock2)
	dig := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col1",
		Namespace:  "ns1",
	}

	store := Dig2PvtRWSetWithConfig{
		DigKey{
			TxId:       "txID1",
			Collection: "col1",
			Namespace:  "ns1",
		}: p2TransientStore,
	}

	p2.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), uint64(0)).Return(store, nil)

	factoryMock3 := &collectionAccessFactoryMock{}
	policyMock3 := &collectionAccessPolicyMock{}
	policyMock3.Setup(1, 2, func(data fcommon.SignedData) bool {
		return false
	}, []string{"org1", "org2"})
	factoryMock3.On("AccessPolicy", mock.Anything, mock.Anything).Return(policyMock3, nil)

	p3 := gn.newPuller("p3", newCollectionStore(), factoryMock3)
	p3.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), uint64(0)).Run(func(_ mock.Arguments) {
		t.Fatal("p3 shouldn't have been selected for pull")
	})

	dasf := &digestsAndSourceFactory{}

	fetchedMessages, err := p1.fetch(dasf.mapDigest(toDigKey(dig)).toSources().create(), uint64(1))
	rws1 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[0].Payload[0])
	rws2 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[0].Payload[1])
	fetched := []util.PrivateRWSet{rws1, rws2}
	assert.NoError(t, err)
	assert.Equal(t, p2TransientStore.RWSet, fetched)
}

func TestPullerDataNotAvailable(t *testing.T) {
	t.Parallel()
	
	
	gn := &gossipNetwork{}
	policyStore := newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p2")
	factoryMock := &collectionAccessFactoryMock{}
	factoryMock.On("AccessPolicy", mock.Anything, mock.Anything).Return(&collectionAccessPolicyMock{}, nil)

	p1 := gn.newPuller("p1", policyStore, factoryMock, membership(peerData{"p2", uint64(1)}, peerData{"p3", uint64(1)})...)

	policyStore = newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p1")
	p2 := gn.newPuller("p2", policyStore, factoryMock)
	dig := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col1",
		Namespace:  "ns1",
	}

	store := Dig2PvtRWSetWithConfig{
		DigKey{
			TxId:       "txID1",
			Collection: "col1",
			Namespace:  "ns1",
		}: &util.PrivateRWSetWithConfig{
			RWSet: []util.PrivateRWSet{},
		},
	}

	p2.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), mock.Anything).Return(store, nil)

	p3 := gn.newPuller("p3", newCollectionStore(), factoryMock)
	p3.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), mock.Anything).Run(func(_ mock.Arguments) {
		t.Fatal("p3 shouldn't have been selected for pull")
	})

	dasf := &digestsAndSourceFactory{}
	fetchedMessages, err := p1.fetch(dasf.mapDigest(toDigKey(dig)).toSources().create(), uint64(1))
	assert.Empty(t, fetchedMessages.AvailableElemenets)
	assert.NoError(t, err)
}

func TestPullerNoPeersKnown(t *testing.T) {
	t.Parallel()
	
	gn := &gossipNetwork{}
	policyStore := newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p2", "p3")
	factoryMock := &collectionAccessFactoryMock{}
	factoryMock.On("AccessPolicy", mock.Anything, mock.Anything).Return(&collectionAccessPolicyMock{}, nil)

	p1 := gn.newPuller("p1", policyStore, factoryMock)
	dasf := &digestsAndSourceFactory{}
	d2s := dasf.mapDigest(&DigKey{Collection: "col1", TxId: "txID1", Namespace: "ns1"}).toSources().create()
	fetchedMessages, err := p1.fetch(d2s, uint64(1))
	assert.Empty(t, fetchedMessages)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Empty membership")
}

func TestPullPeerFilterError(t *testing.T) {
	t.Parallel()
	
	gn := &gossipNetwork{}
	policyStore := newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p2")
	factoryMock := &collectionAccessFactoryMock{}
	factoryMock.On("AccessPolicy", mock.Anything, mock.Anything).Return(&collectionAccessPolicyMock{}, nil)

	p1 := gn.newPuller("p1", policyStore, factoryMock)
	gn.peers[0].On("PeerFilter", mock.Anything, mock.Anything).Return(nil, errors.New("Failed obtaining filter"))
	dasf := &digestsAndSourceFactory{}
	d2s := dasf.mapDigest(&DigKey{Collection: "col1", TxId: "txID1", Namespace: "ns1"}).toSources().create()
	fetchedMessages, err := p1.fetch(d2s, uint64(1))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed obtaining filter")
	assert.Empty(t, fetchedMessages)
}

func TestPullerPeerNotEligible(t *testing.T) {
	t.Parallel()
	
	
	gn := &gossipNetwork{}
	policyStore := newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p2", "p3")
	factoryMock1 := &collectionAccessFactoryMock{}
	accessPolicyMock1 := &collectionAccessPolicyMock{}
	accessPolicyMock1.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p2")) || bytes.Equal(data.Identity, []byte("p3"))
	}, []string{"org1", "org2"})
	factoryMock1.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock1, nil)

	p1 := gn.newPuller("p1", policyStore, factoryMock1, membership(peerData{"p2", uint64(1)}, peerData{"p3", uint64(1)})...)

	policyStore = newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p2")
	factoryMock2 := &collectionAccessFactoryMock{}
	accessPolicyMock2 := &collectionAccessPolicyMock{}
	accessPolicyMock2.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p2"))
	}, []string{"org1", "org2"})
	factoryMock2.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock2, nil)

	p2 := gn.newPuller("p2", policyStore, factoryMock2)

	dig := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col1",
		Namespace:  "ns1",
	}

	store := Dig2PvtRWSetWithConfig{
		DigKey{
			TxId:       "txID1",
			Collection: "col1",
			Namespace:  "ns1",
		}: &util.PrivateRWSetWithConfig{
			RWSet: newPRWSet(),
			CollectionConfig: &fcommon.CollectionConfig{
				Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
					StaticCollectionConfig: &fcommon.StaticCollectionConfig{
						Name: "col1",
					},
				},
			},
		},
	}

	p2.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), mock.Anything).Return(store, nil)

	policyStore = newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p3")
	factoryMock3 := &collectionAccessFactoryMock{}
	accessPolicyMock3 := &collectionAccessPolicyMock{}
	accessPolicyMock3.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p3"))
	}, []string{"org1", "org2"})
	factoryMock3.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock1, nil)

	p3 := gn.newPuller("p3", policyStore, factoryMock3)
	p3.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), mock.Anything).Return(store, nil)
	dasf := &digestsAndSourceFactory{}
	d2s := dasf.mapDigest(&DigKey{Collection: "col1", TxId: "txID1", Namespace: "ns1"}).toSources().create()
	fetchedMessages, err := p1.fetch(d2s, uint64(1))
	assert.Empty(t, fetchedMessages.AvailableElemenets)
	assert.NoError(t, err)
}

func TestPullerDifferentPeersDifferentCollections(t *testing.T) {
	t.Parallel()
	
	
	gn := &gossipNetwork{}
	factoryMock1 := &collectionAccessFactoryMock{}
	accessPolicyMock1 := &collectionAccessPolicyMock{}
	accessPolicyMock1.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p2")) || bytes.Equal(data.Identity, []byte("p3"))
	}, []string{"org1", "org2"})
	factoryMock1.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock1, nil)

	policyStore := newCollectionStore().withPolicy("col2", uint64(100)).thatMapsTo("p2").withPolicy("col3", uint64(100)).thatMapsTo("p3")
	p1 := gn.newPuller("p1", policyStore, factoryMock1, membership(peerData{"p2", uint64(1)}, peerData{"p3", uint64(1)})...)

	p2TransientStore := &util.PrivateRWSetWithConfig{
		RWSet: newPRWSet(),
		CollectionConfig: &fcommon.CollectionConfig{
			Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &fcommon.StaticCollectionConfig{
					Name: "col2",
				},
			},
		},
	}

	policyStore = newCollectionStore().withPolicy("col2", uint64(100)).thatMapsTo("p1")
	factoryMock2 := &collectionAccessFactoryMock{}
	accessPolicyMock2 := &collectionAccessPolicyMock{}
	accessPolicyMock2.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p1"))
	}, []string{"org1", "org2"})
	factoryMock2.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock2, nil)

	p2 := gn.newPuller("p2", policyStore, factoryMock2)
	dig1 := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col2",
		Namespace:  "ns1",
	}

	store1 := Dig2PvtRWSetWithConfig{
		DigKey{
			TxId:       "txID1",
			Collection: "col2",
			Namespace:  "ns1",
		}: p2TransientStore,
	}

	p2.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig1)), mock.Anything).Return(store1, nil)

	p3TransientStore := &util.PrivateRWSetWithConfig{
		RWSet: newPRWSet(),
		CollectionConfig: &fcommon.CollectionConfig{
			Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &fcommon.StaticCollectionConfig{
					Name: "col3",
				},
			},
		},
	}

	store2 := Dig2PvtRWSetWithConfig{
		DigKey{
			TxId:       "txID1",
			Collection: "col3",
			Namespace:  "ns1",
		}: p3TransientStore,
	}
	policyStore = newCollectionStore().withPolicy("col3", uint64(100)).thatMapsTo("p1")
	factoryMock3 := &collectionAccessFactoryMock{}
	accessPolicyMock3 := &collectionAccessPolicyMock{}
	accessPolicyMock3.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p1"))
	}, []string{"org1", "org2"})
	factoryMock3.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock3, nil)

	p3 := gn.newPuller("p3", policyStore, factoryMock3)
	dig2 := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col3",
		Namespace:  "ns1",
	}

	p3.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig2)), mock.Anything).Return(store2, nil)

	dasf := &digestsAndSourceFactory{}
	fetchedMessages, err := p1.fetch(dasf.mapDigest(toDigKey(dig1)).toSources().mapDigest(toDigKey(dig2)).toSources().create(), uint64(1))
	assert.NoError(t, err)
	rws1 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[0].Payload[0])
	rws2 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[0].Payload[1])
	rws3 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[1].Payload[0])
	rws4 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[1].Payload[1])
	fetched := []util.PrivateRWSet{rws1, rws2, rws3, rws4}
	assert.Contains(t, fetched, p2TransientStore.RWSet[0])
	assert.Contains(t, fetched, p2TransientStore.RWSet[1])
	assert.Contains(t, fetched, p3TransientStore.RWSet[0])
	assert.Contains(t, fetched, p3TransientStore.RWSet[1])
}

func TestPullerRetries(t *testing.T) {
	t.Parallel()
	
	
	
	gn := &gossipNetwork{}
	factoryMock1 := &collectionAccessFactoryMock{}
	accessPolicyMock1 := &collectionAccessPolicyMock{}
	accessPolicyMock1.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p2")) || bytes.Equal(data.Identity, []byte("p3")) ||
			bytes.Equal(data.Identity, []byte("p4")) ||
			bytes.Equal(data.Identity, []byte("p5"))
	}, []string{"org1", "org2"})
	factoryMock1.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock1, nil)

	
	policyStore := newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p2", "p3", "p4", "p5")
	p1 := gn.newPuller("p1", policyStore, factoryMock1, membership(peerData{"p2", uint64(1)},
		peerData{"p3", uint64(1)}, peerData{"p4", uint64(1)}, peerData{"p5", uint64(1)})...)

	
	transientStore := &util.PrivateRWSetWithConfig{
		RWSet: newPRWSet(),
		CollectionConfig: &fcommon.CollectionConfig{
			Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &fcommon.StaticCollectionConfig{
					Name: "col1",
				},
			},
		},
	}

	dig := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col1",
		Namespace:  "ns1",
	}

	store := Dig2PvtRWSetWithConfig{
		DigKey{
			TxId:       "txID1",
			Collection: "col1",
			Namespace:  "ns1",
		}: transientStore,
	}

	
	policyStore = newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p2")
	factoryMock2 := &collectionAccessFactoryMock{}
	accessPolicyMock2 := &collectionAccessPolicyMock{}
	accessPolicyMock2.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p2"))
	}, []string{"org1", "org2"})
	factoryMock2.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock2, nil)

	p2 := gn.newPuller("p2", policyStore, factoryMock2)
	p2.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), uint64(0)).Return(store, nil)

	
	policyStore = newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p1")
	factoryMock3 := &collectionAccessFactoryMock{}
	accessPolicyMock3 := &collectionAccessPolicyMock{}
	accessPolicyMock3.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p1"))
	}, []string{"org1", "org2"})
	factoryMock3.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock3, nil)

	p3 := gn.newPuller("p3", policyStore, factoryMock3)
	p3.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), uint64(0)).Return(store, nil)

	
	policyStore = newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p4")
	factoryMock4 := &collectionAccessFactoryMock{}
	accessPolicyMock4 := &collectionAccessPolicyMock{}
	accessPolicyMock4.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p4"))
	}, []string{"org1", "org2"})
	factoryMock4.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock4, nil)

	p4 := gn.newPuller("p4", policyStore, factoryMock4)
	p4.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), uint64(0)).Return(store, nil)

	
	policyStore = newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p5")
	factoryMock5 := &collectionAccessFactoryMock{}
	accessPolicyMock5 := &collectionAccessPolicyMock{}
	accessPolicyMock5.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p5"))
	}, []string{"org1", "org2"})
	factoryMock5.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock5, nil)

	p5 := gn.newPuller("p5", policyStore, factoryMock5)
	p5.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig)), uint64(0)).Return(store, nil)

	
	dasf := &digestsAndSourceFactory{}
	fetchedMessages, err := p1.fetch(dasf.mapDigest(toDigKey(dig)).toSources().create(), uint64(1))
	assert.NoError(t, err)
	rws1 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[0].Payload[0])
	rws2 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[0].Payload[1])
	fetched := []util.PrivateRWSet{rws1, rws2}
	assert.NoError(t, err)
	assert.Equal(t, transientStore.RWSet, fetched)
}

func TestPullerPreferEndorsers(t *testing.T) {
	t.Parallel()
	
	
	
	
	gn := &gossipNetwork{}
	factoryMock := &collectionAccessFactoryMock{}
	accessPolicyMock2 := &collectionAccessPolicyMock{}
	accessPolicyMock2.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p2")) || bytes.Equal(data.Identity, []byte("p1"))
	}, []string{"org1", "org2"})
	factoryMock.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock2, nil)

	policyStore := newCollectionStore().
		withPolicy("col1", uint64(100)).
		thatMapsTo("p1", "p2", "p3", "p4", "p5").
		withPolicy("col2", uint64(100)).
		thatMapsTo("p1", "p2")
	p1 := gn.newPuller("p1", policyStore, factoryMock, membership(peerData{"p2", uint64(1)},
		peerData{"p3", uint64(1)}, peerData{"p4", uint64(1)}, peerData{"p5", uint64(1)})...)

	p3TransientStore := &util.PrivateRWSetWithConfig{
		RWSet: newPRWSet(),
		CollectionConfig: &fcommon.CollectionConfig{
			Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &fcommon.StaticCollectionConfig{
					Name: "col2",
				},
			},
		},
	}

	p2TransientStore := &util.PrivateRWSetWithConfig{
		RWSet: newPRWSet(),
		CollectionConfig: &fcommon.CollectionConfig{
			Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &fcommon.StaticCollectionConfig{
					Name: "col2",
				},
			},
		},
	}

	p2 := gn.newPuller("p2", policyStore, factoryMock)
	p3 := gn.newPuller("p3", policyStore, factoryMock)
	gn.newPuller("p4", policyStore, factoryMock)
	gn.newPuller("p5", policyStore, factoryMock)

	dig1 := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col1",
		Namespace:  "ns1",
	}

	dig2 := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col2",
		Namespace:  "ns1",
	}

	store := Dig2PvtRWSetWithConfig{
		DigKey{
			TxId:       "txID1",
			Collection: "col1",
			Namespace:  "ns1",
		}: p3TransientStore,
		DigKey{
			TxId:       "txID1",
			Collection: "col2",
			Namespace:  "ns1",
		}: p2TransientStore,
	}

	
	
	p2.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig2)), uint64(0)).Return(store, nil)

	
	
	p3.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig1)), uint64(0)).Return(store, nil)

	dasf := &digestsAndSourceFactory{}
	d2s := dasf.mapDigest(toDigKey(dig1)).toSources("p3").mapDigest(toDigKey(dig2)).toSources().create()
	fetchedMessages, err := p1.fetch(d2s, uint64(1))
	assert.NoError(t, err)
	rws1 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[0].Payload[0])
	rws2 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[0].Payload[1])
	rws3 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[1].Payload[0])
	rws4 := util.PrivateRWSet(fetchedMessages.AvailableElemenets[1].Payload[1])
	fetched := []util.PrivateRWSet{rws1, rws2, rws3, rws4}
	assert.Contains(t, fetched, p3TransientStore.RWSet[0])
	assert.Contains(t, fetched, p3TransientStore.RWSet[1])
	assert.Contains(t, fetched, p2TransientStore.RWSet[0])
	assert.Contains(t, fetched, p2TransientStore.RWSet[1])
}

func TestPullerAvoidPullingPurgedData(t *testing.T) {
	
	
	
	

	t.Parallel()
	gn := &gossipNetwork{}
	factoryMock := &collectionAccessFactoryMock{}
	accessPolicyMock2 := &collectionAccessPolicyMock{}
	accessPolicyMock2.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p1"))
	}, []string{"org1", "org2"})
	factoryMock.On("AccessPolicy", mock.Anything, mock.Anything).Return(accessPolicyMock2, nil)

	policyStore := newCollectionStore().withPolicy("col1", uint64(100)).thatMapsTo("p1", "p2", "p3").
		withPolicy("col2", uint64(1000)).thatMapsTo("p1", "p2", "p3")

	
	p1 := gn.newPuller("p1", policyStore, factoryMock, membership(peerData{"p2", uint64(1)},
		peerData{"p3", uint64(111)})...)

	privateData1 := &util.PrivateRWSetWithConfig{
		RWSet: newPRWSet(),
		CollectionConfig: &fcommon.CollectionConfig{
			Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &fcommon.StaticCollectionConfig{
					Name: "col1",
				},
			},
		},
	}
	privateData2 := &util.PrivateRWSetWithConfig{
		RWSet: newPRWSet(),
		CollectionConfig: &fcommon.CollectionConfig{
			Payload: &fcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &fcommon.StaticCollectionConfig{
					Name: "col2",
				},
			},
		},
	}

	p2 := gn.newPuller("p2", policyStore, factoryMock)
	p3 := gn.newPuller("p3", policyStore, factoryMock)

	dig1 := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col1",
		Namespace:  "ns1",
	}

	dig2 := &proto.PvtDataDigest{
		TxId:       "txID1",
		Collection: "col2",
		Namespace:  "ns1",
	}

	store := Dig2PvtRWSetWithConfig{
		DigKey{
			TxId:       "txID1",
			Collection: "col1",
			Namespace:  "ns1",
		}: privateData1,
		DigKey{
			TxId:       "txID1",
			Collection: "col2",
			Namespace:  "ns1",
		}: privateData2,
	}

	p2.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig1)), 0).Return(store, nil)
	p3.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig1)), 0).Return(store, nil).
		Run(
			func(arg mock.Arguments) {
				assert.Fail(t, "we should not fetch private data from peers where it was purged")
			},
		)

	p3.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig2)), uint64(0)).Return(store, nil)
	p2.PrivateDataRetriever.(*dataRetrieverMock).On("CollectionRWSet", mock.MatchedBy(protoMatcher(dig2)), uint64(0)).Return(store, nil).
		Run(
			func(mock.Arguments) {
				assert.Fail(t, "we should not fetch private data of collection2 from peer 2")

			},
		)

	dasf := &digestsAndSourceFactory{}
	d2s := dasf.mapDigest(toDigKey(dig1)).toSources("p3", "p2").mapDigest(toDigKey(dig2)).toSources("p3").create()
	
	fetchedMessages, err := p1.fetch(d2s, uint64(1))

	assert.NoError(t, err)
	assert.Equal(t, 1, len(fetchedMessages.PurgedElements))
	assert.Equal(t, dig1, fetchedMessages.PurgedElements[0])
	p3.PrivateDataRetriever.(*dataRetrieverMock).AssertNumberOfCalls(t, "CollectionRWSet", 1)

}

type counterDataRetreiver struct {
	numberOfCalls int
	PrivateDataRetriever
}

func (c *counterDataRetreiver) CollectionRWSet(dig []*proto.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, error) {
	c.numberOfCalls += 1
	return c.PrivateDataRetriever.CollectionRWSet(dig, blockNum)
}

func (c *counterDataRetreiver) getNumberOfCalls() int {
	return c.numberOfCalls
}

func TestPullerIntegratedWithDataRetreiver(t *testing.T) {
	t.Parallel()
	gn := &gossipNetwork{}

	ns1, ns2 := "testChaincodeName1", "testChaincodeName2"
	col1, col2 := "testCollectionName1", "testCollectionName2"

	ap := &collectionAccessPolicyMock{}
	ap.Setup(1, 2, func(data fcommon.SignedData) bool {
		return bytes.Equal(data.Identity, []byte("p1"))
	}, []string{"org1", "org2"})

	factoryMock := &collectionAccessFactoryMock{}
	factoryMock.On("AccessPolicy", mock.Anything, mock.Anything).Return(ap, nil)

	policyStore := newCollectionStore().withPolicy(col1, uint64(1000)).thatMapsTo("p1", "p2").
		withPolicy(col2, uint64(1000)).thatMapsTo("p1", "p2")

	p1 := gn.newPuller("p1", policyStore, factoryMock, membership(peerData{"p2", uint64(10)})...)
	p2 := gn.newPuller("p2", policyStore, factoryMock, membership(peerData{"p1", uint64(1)})...)

	dataStore := &mocks.DataStore{}
	result := []*ledger.TxPvtData{
		{
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					pvtReadWriteSet(ns1, col1, []byte{1}),
					pvtReadWriteSet(ns1, col1, []byte{2}),
				},
			},
			SeqInBlock: 1,
		},
		{
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					pvtReadWriteSet(ns2, col2, []byte{3}),
					pvtReadWriteSet(ns2, col2, []byte{4}),
				},
			},
			SeqInBlock: 2,
		},
	}

	dataStore.On("LedgerHeight").Return(uint64(10), nil)
	dataStore.On("GetPvtDataByNum", uint64(5), mock.Anything).Return(result, nil)
	historyRetreiver := &mocks.ConfigHistoryRetriever{}
	historyRetreiver.On("MostRecentCollectionConfigBelow", mock.Anything, ns1).Return(newCollectionConfig(col1), nil)
	historyRetreiver.On("MostRecentCollectionConfigBelow", mock.Anything, ns2).Return(newCollectionConfig(col2), nil)
	dataStore.On("GetConfigHistoryRetriever").Return(historyRetreiver, nil)

	dataRetreiver := &counterDataRetreiver{PrivateDataRetriever: NewDataRetriever(dataStore), numberOfCalls: 0}
	p2.PrivateDataRetriever = dataRetreiver

	dig1 := &DigKey{
		TxId:       "txID1",
		Collection: col1,
		Namespace:  ns1,
		BlockSeq:   5,
		SeqInBlock: 1,
	}

	dig2 := &DigKey{
		TxId:       "txID1",
		Collection: col2,
		Namespace:  ns2,
		BlockSeq:   5,
		SeqInBlock: 2,
	}

	dasf := &digestsAndSourceFactory{}
	d2s := dasf.mapDigest(dig1).toSources("p2").mapDigest(dig2).toSources("p2").create()
	fetchedMessages, err := p1.fetch(d2s, uint64(1))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(fetchedMessages.AvailableElemenets))
	assert.Equal(t, 1, dataRetreiver.getNumberOfCalls())
	assert.Equal(t, 2, len(fetchedMessages.AvailableElemenets[0].Payload))
	assert.Equal(t, 2, len(fetchedMessages.AvailableElemenets[1].Payload))
}

func toDigKey(dig *proto.PvtDataDigest) *DigKey {
	return &DigKey{
		TxId:       dig.TxId,
		BlockSeq:   dig.BlockSeq,
		SeqInBlock: dig.SeqInBlock,
		Namespace:  dig.Namespace,
		Collection: dig.Collection,
	}
}
