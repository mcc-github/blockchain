/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"bytes"
	"errors"
	"testing"
	"time"

	pb "github.com/mcc-github/blockchain-protos-go/common"
	mb "github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

func createCollectionPolicyConfig(accessPolicy *pb.SignaturePolicyEnvelope) *pb.CollectionPolicyConfig {
	cpcSp := &pb.CollectionPolicyConfig_SignaturePolicy{
		SignaturePolicy: accessPolicy,
	}
	cpc := &pb.CollectionPolicyConfig{
		Payload: cpcSp,
	}
	return cpc
}

type mockIdentity struct {
	idBytes []byte
}

func (id *mockIdentity) Anonymous() bool {
	panic("implement me")
}

func (id *mockIdentity) ExpiresAt() time.Time {
	return time.Time{}
}

func (id *mockIdentity) SatisfiesPrincipal(p *mb.MSPPrincipal) error {
	if bytes.Equal(id.idBytes, p.Principal) {
		return nil
	}
	return errors.New("Principals do not match")
}

func (id *mockIdentity) GetIdentifier() *msp.IdentityIdentifier {
	return &msp.IdentityIdentifier{Mspid: "Mock", Id: string(id.idBytes)}
}

func (id *mockIdentity) GetMSPIdentifier() string {
	return string(id.idBytes)
}

func (id *mockIdentity) Validate() error {
	return nil
}

func (id *mockIdentity) GetOrganizationalUnits() []*msp.OUIdentifier {
	return nil
}

func (id *mockIdentity) Verify(msg []byte, sig []byte) error {
	if bytes.Equal(sig, []byte("badsigned")) {
		return errors.New("Invalid signature")
	}
	return nil
}

func (id *mockIdentity) Serialize() ([]byte, error) {
	return id.idBytes, nil
}

type mockDeserializer struct {
	fail error
}

func (md *mockDeserializer) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	if md.fail != nil {
		return nil, md.fail
	}
	return &mockIdentity{idBytes: serializedIdentity}, nil
}

func (md *mockDeserializer) IsWellFormed(_ *mb.SerializedIdentity) error {
	return nil
}

func TestNewSimpleCollectionWithBadConfig(t *testing.T) {
	
	_, err := NewSimpleCollection(nil, &mockDeserializer{})
	assert.Error(t, err)

	
	collectionConfig := &pb.StaticCollectionConfig{
		Name:              "test collection",
		RequiredPeerCount: 1,
		MemberOrgsPolicy:  getBadAccessPolicy([]string{"peer0", "peer1"}, 3),
	}
	_, err = NewSimpleCollection(collectionConfig, &mockDeserializer{})
	assert.Error(t, err)
	assert.EqualError(t, err, "failed constructing policy object out of collection policy config: identity index out of range, requested 3, but identities length is 2")
}

func TestNewSimpleCollectionWithGoodConfig(t *testing.T) {
	
	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	accessPolicy := createCollectionPolicyConfig(policyEnvelope)

	
	collectionConfig := &pb.StaticCollectionConfig{
		Name:              "test collection",
		RequiredPeerCount: 1,
		MemberOrgsPolicy:  accessPolicy,
	}

	
	sc, err := NewSimpleCollection(collectionConfig, &mockDeserializer{})
	assert.NoError(t, err)

	
	assert.True(t, sc.CollectionID() == "test collection")

	
	members := sc.MemberOrgs()
	assert.True(t, members[0] == "signer0")
	assert.True(t, members[1] == "signer1")

	
	assert.True(t, sc.RequiredPeerCount() == 1)
}

func TestSetupWithBadConfig(t *testing.T) {
	
	var sc SimpleCollection
	err := sc.Setup(&pb.StaticCollectionConfig{}, &mockDeserializer{})
	assert.Error(t, err)

	
	collectionConfig := &pb.StaticCollectionConfig{
		Name:              "test collection",
		RequiredPeerCount: 1,
		MemberOrgsPolicy:  getBadAccessPolicy([]string{"peer0", "peer1"}, 3),
	}
	err = sc.Setup(collectionConfig, &mockDeserializer{})
	assert.Error(t, err)
	assert.EqualError(t, err, "failed constructing policy object out of collection policy config: identity index out of range, requested 3, but identities length is 2")
}

func TestSetupGoodConfigCollection(t *testing.T) {
	
	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	accessPolicy := createCollectionPolicyConfig(policyEnvelope)

	
	collectionConfig := &pb.StaticCollectionConfig{
		Name:              "test collection",
		RequiredPeerCount: 1,
		MemberOrgsPolicy:  accessPolicy,
	}

	
	var sc SimpleCollection
	err := sc.Setup(collectionConfig, &mockDeserializer{})
	assert.NoError(t, err)

	
	assert.True(t, sc.CollectionID() == "test collection")

	
	members := sc.MemberOrgs()
	assert.True(t, members[0] == "signer0")
	assert.True(t, members[1] == "signer1")

	
	assert.True(t, sc.RequiredPeerCount() == 1)
}

func TestSimpleCollectionFilter(t *testing.T) {
	
	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	accessPolicy := createCollectionPolicyConfig(policyEnvelope)

	
	collectionConfig := &pb.StaticCollectionConfig{
		Name:              "test collection",
		RequiredPeerCount: 1,
		MemberOrgsPolicy:  accessPolicy,
	}

	
	var sc SimpleCollection
	err := sc.Setup(collectionConfig, &mockDeserializer{})
	assert.NoError(t, err)

	
	var cap CollectionAccessPolicy
	cap = &sc
	accessFilter := cap.AccessFilter()

	
	notMember := protoutil.SignedData{
		Identity:  []byte{1, 2, 3},
		Signature: []byte{},
		Data:      []byte{},
	}
	assert.False(t, accessFilter(notMember))

	
	member := protoutil.SignedData{
		Identity:  signers[0],
		Signature: []byte{},
		Data:      []byte{},
	}
	assert.True(t, accessFilter(member))
}
