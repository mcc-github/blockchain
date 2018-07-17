/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorsement

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/policies/inquire"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	discoveryprotos "github.com/mcc-github/blockchain/protos/discovery"
	"github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var pkiID2MSPID = map[string]string{
	"p0":  "Org0MSP",
	"p1":  "Org1MSP",
	"p2":  "Org2MSP",
	"p3":  "Org3MSP",
	"p4":  "Org4MSP",
	"p5":  "Org5MSP",
	"p6":  "Org6MSP",
	"p7":  "Org7MSP",
	"p8":  "Org8MSP",
	"p9":  "Org9MSP",
	"p10": "Org10MSP",
	"p11": "Org11MSP",
	"p12": "Org12MSP",
}

func TestPeersForEndorsement(t *testing.T) {
	peerRole := func(pkiID string) *msp.MSPPrincipal {
		return &msp.MSPPrincipal{
			PrincipalClassification: msp.MSPPrincipal_ROLE,
			Principal: utils.MarshalOrPanic(&msp.MSPRole{
				MspIdentifier: pkiID2MSPID[pkiID],
				Role:          msp.MSPRole_PEER,
			}),
		}
	}
	extractPeers := func(desc *discoveryprotos.EndorsementDescriptor) map[string]struct{} {
		res := make(map[string]struct{})
		for _, endorsers := range desc.EndorsersByGroups {
			for _, p := range endorsers.Peers {
				res[string(p.Identity)] = struct{}{}
				assert.Equal(t, string(p.Identity), string(p.MembershipInfo.Payload))
				assert.Equal(t, string(p.Identity), string(p.StateInfo.Payload))
			}
		}
		return res
	}
	cc := "chaincode"
	mf := &metadataFetcher{}
	g := &gossipMock{}
	pf := &policyFetcherMock{}
	ccWithMissingPolicy := "chaincodeWithMissingPolicy"
	channel := common.ChainID("test")
	alivePeers := peerSet{
		newPeer(0),
		newPeer(2),
		newPeer(4),
		newPeer(6),
		newPeer(8),
		newPeer(10),
		newPeer(11),
		newPeer(12),
	}

	identities := identitySet(pkiID2MSPID)

	chanPeers := peerSet{
		newPeer(0).withChaincode(cc, "1.0"),
		newPeer(3).withChaincode(cc, "1.0"),
		newPeer(6).withChaincode(cc, "1.0"),
		newPeer(9).withChaincode(cc, "1.0"),
		newPeer(11).withChaincode(cc, "1.0"),
		newPeer(12).withChaincode(cc, "1.0"),
	}
	g.On("Peers").Return(alivePeers.toMembers())
	g.On("IdentityInfo").Return(identities)

	
	t.Run("PolicyNotFound", func(t *testing.T) {
		pf.On("PolicyByChaincode", ccWithMissingPolicy).Return(nil).Once()
		g.On("PeersOfChannel").Return(chanPeers.toMembers()).Once()
		mf.On("Metadata").Return(&chaincode.Metadata{Name: cc, Version: "1.0"}).Once()
		analyzer := NewEndorsementAnalyzer(g, pf, &principalEvaluatorMock{}, mf)
		desc, err := analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{Chaincodes: []*discoveryprotos.ChaincodeCall{{Name: ccWithMissingPolicy}}})
		assert.Nil(t, desc)
		assert.Equal(t, "policy not found", err.Error())
	})

	t.Run("NotEnoughPeers", func(t *testing.T) {
		
		
		
		
		pb := principalBuilder{}
		policy := pb.newSet().addPrincipal(peerRole("p1")).addPrincipal(peerRole("p6")).
			newSet().addPrincipal(peerRole("p11")).addPrincipal(peerRole("p11")).buildPolicy()
		g.On("PeersOfChannel").Return(chanPeers.toMembers()).Once()
		mf.On("Metadata").Return(&chaincode.Metadata{Name: cc, Version: "1.0"}).Once()
		analyzer := NewEndorsementAnalyzer(g, pf, &principalEvaluatorMock{}, mf)
		pf.On("PolicyByChaincode", cc).Return(policy).Once()
		desc, err := analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{Chaincodes: []*discoveryprotos.ChaincodeCall{{Name: cc}}})
		assert.Nil(t, desc)
		assert.Equal(t, err.Error(), "cannot satisfy any principal combination")
	})

	t.Run("DisjointViews", func(t *testing.T) {
		pb := principalBuilder{}
		
		
		
		
		policy := pb.newSet().addPrincipal(peerRole("p0")).addPrincipal(peerRole("p6")).
			newSet().addPrincipal(peerRole("p10")).addPrincipal(peerRole("p12")).buildPolicy()
		g.On("PeersOfChannel").Return(chanPeers.toMembers()).Once()
		mf.On("Metadata").Return(&chaincode.Metadata{Name: cc, Version: "1.0"}).Once()
		analyzer := NewEndorsementAnalyzer(g, pf, &principalEvaluatorMock{}, mf)
		pf.On("PolicyByChaincode", cc).Return(policy).Once()
		desc, err := analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{Chaincodes: []*discoveryprotos.ChaincodeCall{{Name: cc}}})
		assert.NoError(t, err)
		assert.NotNil(t, desc)
		assert.Len(t, desc.Layouts, 1)
		assert.Len(t, desc.Layouts[0].QuantitiesByGroup, 2)
		assert.Equal(t, map[string]struct{}{
			peerIdentityString("p0"): {},
			peerIdentityString("p6"): {},
		}, extractPeers(desc))
	})

	t.Run("MultipleCombinations", func(t *testing.T) {
		
		
		
		
		pb := principalBuilder{}
		policy := pb.newSet().addPrincipal(peerRole("p0")).addPrincipal(peerRole("p6")).
			newSet().addPrincipal(peerRole("p12")).buildPolicy()
		g.On("PeersOfChannel").Return(chanPeers.toMembers()).Once()
		mf.On("Metadata").Return(&chaincode.Metadata{Name: cc, Version: "1.0"}).Once()
		analyzer := NewEndorsementAnalyzer(g, pf, &principalEvaluatorMock{}, mf)
		pf.On("PolicyByChaincode", cc).Return(policy).Once()
		desc, err := analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{Chaincodes: []*discoveryprotos.ChaincodeCall{{Name: cc}}})
		assert.NoError(t, err)
		assert.NotNil(t, desc)
		assert.Len(t, desc.Layouts, 2)
		assert.Len(t, desc.Layouts[0].QuantitiesByGroup, 2)
		assert.Len(t, desc.Layouts[1].QuantitiesByGroup, 1)
		assert.Equal(t, map[string]struct{}{
			peerIdentityString("p0"):  {},
			peerIdentityString("p6"):  {},
			peerIdentityString("p12"): {},
		}, extractPeers(desc))
	})

	t.Run("WrongVersionInstalled", func(t *testing.T) {
		
		
		mf.On("Metadata").Return(&chaincode.Metadata{
			Name: cc, Version: "1.1",
		}).Once()
		pb := principalBuilder{}
		policy := pb.newSet().addPrincipal(peerRole("p0")).addPrincipal(peerRole("p6")).
			newSet().addPrincipal(peerRole("p12")).buildPolicy()
		g.On("PeersOfChannel").Return(chanPeers.toMembers()).Once()
		pf.On("PolicyByChaincode", cc).Return(policy).Once()
		analyzer := NewEndorsementAnalyzer(g, pf, &principalEvaluatorMock{}, mf)
		desc, err := analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{Chaincodes: []*discoveryprotos.ChaincodeCall{{Name: cc}}})
		assert.Nil(t, desc)
		assert.Equal(t, err.Error(), "chaincode isn't installed on sufficient organizations required by the endorsement policy")

		
		
		chanPeers := peerSet{
			newPeer(0).withChaincode(cc, "1.0"),
			newPeer(3).withChaincode(cc, "1.0"),
			newPeer(6).withChaincode(cc, "1.0"),
			newPeer(9).withChaincode(cc, "1.0"),
			newPeer(12).withChaincode(cc, "1.0"),
		}
		chanPeers[0].Properties.Chaincodes[0].Version = "0.6"
		chanPeers[4].Properties = nil
		g.On("PeersOfChannel").Return(chanPeers.toMembers()).Once()
		pf.On("PolicyByChaincode", cc).Return(policy).Once()
		mf.On("Metadata").Return(&chaincode.Metadata{
			Name: cc, Version: "1.0",
		}).Once()
		desc, err = analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{Chaincodes: []*discoveryprotos.ChaincodeCall{{Name: cc}}})
		assert.Nil(t, desc)
		assert.Equal(t, err.Error(), "chaincode isn't installed on sufficient organizations required by the endorsement policy")
	})

	t.Run("NoChaincodeMetadataFromLedger", func(t *testing.T) {
		
		
		
		pb := principalBuilder{}
		policy := pb.newSet().addPrincipal(peerRole("p0")).addPrincipal(peerRole("p6")).
			newSet().addPrincipal(peerRole("p12")).buildPolicy()
		pf.On("PolicyByChaincode", cc).Return(policy).Once()
		mf.On("Metadata").Return(nil).Once()
		analyzer := NewEndorsementAnalyzer(g, pf, &principalEvaluatorMock{}, mf)
		desc, err := analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{Chaincodes: []*discoveryprotos.ChaincodeCall{{Name: cc}}})
		assert.Nil(t, desc)
		assert.Equal(t, err.Error(), "No metadata was found for chaincode chaincode in channel test")
	})

	t.Run("Collections", func(t *testing.T) {
		
		
		
		
		collectionOrgs := []*msp.MSPPrincipal{
			peerRole("p0"),
			peerRole("p12"),
		}
		col2principals := map[string][]*msp.MSPPrincipal{
			"collection": collectionOrgs,
		}
		mf.On("Metadata").Return(&chaincode.Metadata{
			Name: cc, Version: "1.0", CollectionsConfig: buildCollectionConfig(col2principals),
		}).Once()
		pb := principalBuilder{}
		policy := pb.newSet().addPrincipal(peerRole("p0")).
			addPrincipal(peerRole("p6")).newSet().
			addPrincipal(peerRole("p12")).buildPolicy()
		g.On("PeersOfChannel").Return(chanPeers.toMembers()).Once()
		pf.On("PolicyByChaincode", cc).Return(policy).Once()
		analyzer := NewEndorsementAnalyzer(g, pf, &principalEvaluatorMock{}, mf)
		desc, err := analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{
			Chaincodes: []*discoveryprotos.ChaincodeCall{
				{
					Name:            cc,
					CollectionNames: []string{"collection"},
				},
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, desc)
		assert.Len(t, desc.Layouts, 1)
		assert.Len(t, desc.Layouts[0].QuantitiesByGroup, 1)
		assert.Equal(t, map[string]struct{}{
			peerIdentityString("p12"): {},
		}, extractPeers(desc))
	})

	t.Run("Chaincode2Chaincode", func(t *testing.T) {
		
		
		
		
		
		
		

		chanPeers := peerSet{}
		for _, id := range []int{0, 2, 4, 6, 10, 12} {
			peer := newPeer(id).withChaincode("cc1", "1.0").withChaincode("cc2", "1.0").withChaincode("cc3", "1.0")
			chanPeers = append(chanPeers, peer)
		}

		g.On("PeersOfChannel").Return(chanPeers.toMembers()).Once()

		mf.On("Metadata").Return(&chaincode.Metadata{
			Name: "cc1", Version: "1.0",
		}).Once()
		mf.On("Metadata").Return(&chaincode.Metadata{
			Name: "cc2", Version: "1.0",
		}).Once()
		mf.On("Metadata").Return(&chaincode.Metadata{
			Name: "cc3", Version: "1.0",
		}).Once()

		pb := principalBuilder{}
		cc1policy := pb.newSet().addPrincipal(peerRole("p0")).addPrincipal(peerRole("p2")).
			newSet().addPrincipal(peerRole("p6")).addPrincipal(peerRole("p10")).buildPolicy()

		pf.On("PolicyByChaincode", "cc1").Return(cc1policy).Once()

		cc2policy := pb.newSet().addPrincipal(peerRole("p6")).
			addPrincipal(peerRole("p10")).addPrincipal(peerRole("p12")).buildPolicy()
		pf.On("PolicyByChaincode", "cc2").Return(cc2policy).Once()

		cc3policy := pb.newSet().addPrincipal(peerRole("p4")).
			addPrincipal(peerRole("p12")).buildPolicy()
		pf.On("PolicyByChaincode", "cc3").Return(cc3policy).Once()

		analyzer := NewEndorsementAnalyzer(g, pf, &principalEvaluatorMock{}, mf)
		desc, err := analyzer.PeersForEndorsement(channel, &discoveryprotos.ChaincodeInterest{
			Chaincodes: []*discoveryprotos.ChaincodeCall{
				{
					Name: "cc1",
				},
				{
					Name: "cc2",
				},
				{
					Name: "cc3",
				},
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, desc)
		assert.Len(t, desc.Layouts, 1)
		assert.Len(t, desc.Layouts[0].QuantitiesByGroup, 4)
		assert.Equal(t, map[string]struct{}{
			peerIdentityString("p4"):  {},
			peerIdentityString("p6"):  {},
			peerIdentityString("p10"): {},
			peerIdentityString("p12"): {},
		}, extractPeers(desc))
	})
}

func TestPop(t *testing.T) {
	slice := []inquire.ComparablePrincipalSets{{}, {}}
	assert.Len(t, slice, 2)
	_, slice, err := popComparablePrincipalSets(slice)
	assert.NoError(t, err)
	assert.Len(t, slice, 1)
	_, slice, err = popComparablePrincipalSets(slice)
	assert.Len(t, slice, 0)
	_, slice, err = popComparablePrincipalSets(slice)
	assert.Error(t, err)
	assert.Equal(t, "no principal sets remained after filtering", err.Error())
}

func TestMergePrincipalSetsNilInput(t *testing.T) {
	_, err := mergePrincipalSets(nil)
	assert.Error(t, err)
	assert.Equal(t, "no principal sets remained after filtering", err.Error())
}

func TestComputePrincipalSetsNoPolicies(t *testing.T) {
	

	interest := &discoveryprotos.ChaincodeInterest{
		Chaincodes: []*discoveryprotos.ChaincodeCall{},
	}
	ea := &endorsementAnalyzer{}
	acceptAll := func(policies.PrincipalSet) bool {
		return true
	}
	_, err := ea.computePrincipalSets(common.ChainID("mychannel"), interest, acceptAll)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no principal sets remained after filtering")
}

func TestLoadMetadataAndFiltersCollectionNotPresentInConfig(t *testing.T) {
	interest := &discoveryprotos.ChaincodeInterest{
		Chaincodes: []*discoveryprotos.ChaincodeCall{
			{
				Name:            "mycc",
				CollectionNames: []string{"bar"},
			},
		},
	}

	org1AndOrg2 := []*msp.MSPPrincipal{orgPrincipal("Org1MSP"), orgPrincipal("Org2MSP")}
	col2principals := map[string][]*msp.MSPPrincipal{
		"foo": org1AndOrg2,
	}
	config := buildCollectionConfig(col2principals)

	mdf := &metadataFetcher{}
	mdf.On("Metadata").Return(&chaincode.Metadata{
		Name:              "mycc",
		CollectionsConfig: config,
		Policy:            []byte{1, 2, 3},
	})

	_, err := loadMetadataAndFilters(metadataAndFilterContext{
		identityInfoByID: nil,
		evaluator:        nil,
		chainID:          common.ChainID("mychannel"),
		fetch:            mdf,
		interest:         interest,
	})

	assert.Equal(t, "collection bar doesn't exist in collection config for chaincode mycc", err.Error())
}

func TestLoadMetadataAndFiltersInvalidCollectionData(t *testing.T) {
	interest := &discoveryprotos.ChaincodeInterest{
		Chaincodes: []*discoveryprotos.ChaincodeCall{
			{
				Name:            "mycc",
				CollectionNames: []string{"col1"},
			},
		},
	}
	mdf := &metadataFetcher{}
	mdf.On("Metadata").Return(&chaincode.Metadata{
		Name:              "mycc",
		CollectionsConfig: []byte{1, 2, 3},
		Policy:            []byte{1, 2, 3},
	})

	_, err := loadMetadataAndFilters(metadataAndFilterContext{
		identityInfoByID: nil,
		evaluator:        nil,
		chainID:          common.ChainID("mychannel"),
		fetch:            mdf,
		interest:         interest,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid collection bytes")
}

type peerSet []*peerInfo

func (p peerSet) toMembers() discovery.Members {
	var members discovery.Members
	for _, peer := range p {
		members = append(members, peer.NetworkMember)
	}
	return members
}

func identitySet(pkiID2MSPID map[string]string) api.PeerIdentitySet {
	var res api.PeerIdentitySet
	for pkiID, mspID := range pkiID2MSPID {
		sID := &msp.SerializedIdentity{
			Mspid:   pkiID2MSPID[pkiID],
			IdBytes: []byte(pkiID),
		}
		res = append(res, api.PeerIdentityInfo{
			Identity:     api.PeerIdentityType(utils.MarshalOrPanic(sID)),
			PKIId:        common.PKIidType(pkiID),
			Organization: api.OrgIdentityType(mspID),
		})
	}
	return res
}

type peerInfo struct {
	identity api.PeerIdentityType
	pkiID    common.PKIidType
	discovery.NetworkMember
}

func peerIdentityString(id string) string {
	return string(utils.MarshalOrPanic(&msp.SerializedIdentity{
		Mspid:   pkiID2MSPID[id],
		IdBytes: []byte(id),
	}))
}

func newPeer(i int) *peerInfo {
	p := fmt.Sprintf("p%d", i)
	identity := utils.MarshalOrPanic(&msp.SerializedIdentity{
		Mspid:   pkiID2MSPID[p],
		IdBytes: []byte(p),
	})
	return &peerInfo{
		pkiID:    common.PKIidType(p),
		identity: api.PeerIdentityType(identity),
		NetworkMember: discovery.NetworkMember{
			PKIid:            common.PKIidType(p),
			Endpoint:         p,
			InternalEndpoint: p,
			Envelope: &gossip.Envelope{
				Payload: []byte(identity),
			},
		},
	}
}

func (pi *peerInfo) withChaincode(name, version string) *peerInfo {
	if pi.Properties == nil {
		pi.Properties = &gossip.Properties{}
	}
	pi.Properties.Chaincodes = append(pi.Properties.Chaincodes, &gossip.Chaincode{
		Name:    name,
		Version: version,
	})
	return pi
}

type gossipMock struct {
	mock.Mock
}

func (g *gossipMock) IdentityInfo() api.PeerIdentitySet {
	return g.Called().Get(0).(api.PeerIdentitySet)
}

func (g *gossipMock) PeersOfChannel(_ common.ChainID) discovery.Members {
	members := g.Called().Get(0)
	return members.(discovery.Members)
}

func (g *gossipMock) Peers() discovery.Members {
	members := g.Called().Get(0)
	return members.(discovery.Members)
}

type policyFetcherMock struct {
	mock.Mock
}

func (pf *policyFetcherMock) PolicyByChaincode(channel string, chaincode string) policies.InquireablePolicy {
	arg := pf.Called(chaincode)
	if arg.Get(0) == nil {
		return nil
	}
	return arg.Get(0).(policies.InquireablePolicy)
}

type principalBuilder struct {
	ip inquireablePolicy
}

func (pb *principalBuilder) buildPolicy() inquireablePolicy {
	defer func() {
		pb.ip = nil
	}()
	return pb.ip
}

func (pb *principalBuilder) newSet() *principalBuilder {
	pb.ip = append(pb.ip, make(policies.PrincipalSet, 0))
	return pb
}

func (pb *principalBuilder) addPrincipal(principal *msp.MSPPrincipal) *principalBuilder {
	pb.ip[len(pb.ip)-1] = append(pb.ip[len(pb.ip)-1], principal)
	return pb
}

type inquireablePolicy []policies.PrincipalSet

func (ip inquireablePolicy) SatisfiedBy() []policies.PrincipalSet {
	return ip
}

type principalEvaluatorMock struct {
}

func (pe *principalEvaluatorMock) MSPOfPrincipal(principal *msp.MSPPrincipal) string {
	role := &msp.MSPRole{}
	proto.Unmarshal(principal.Principal, role)
	return role.MspIdentifier
}

func (pe *principalEvaluatorMock) SatisfiesPrincipal(channel string, identity []byte, principal *msp.MSPPrincipal) error {
	peerRole := &msp.MSPRole{}
	if err := proto.Unmarshal(principal.Principal, peerRole); err != nil {
		return err
	}
	sId := &msp.SerializedIdentity{}
	if err := proto.Unmarshal(identity, sId); err != nil {
		return err
	}
	if peerRole.MspIdentifier == sId.Mspid {
		return nil
	}
	return errors.New("not satisfies")
}

type metadataFetcher struct {
	mock.Mock
}

func (mf *metadataFetcher) Metadata(channel string, cc string, _ bool) *chaincode.Metadata {
	arg := mf.Called().Get(0)
	if arg == nil {
		return nil
	}
	return arg.(*chaincode.Metadata)
}
