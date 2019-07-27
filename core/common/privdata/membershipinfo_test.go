/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"testing"

	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

func TestMembershipInfoProvider(t *testing.T) {
	peerSelfSignedData := protoutil.SignedData{
		Identity:  []byte("peer0"),
		Signature: []byte{1, 2, 3},
		Data:      []byte{4, 5, 6},
	}

	identityDeserializer := func(chainID string) msp.IdentityDeserializer {
		return &mockDeserializer{}
	}

	
	membershipProvider := NewMembershipInfoProvider(peerSelfSignedData, identityDeserializer)
	res, err := membershipProvider.AmMemberOf("test1", getAccessPolicy([]string{"peer0", "peer1"}))
	assert.True(t, res)
	assert.Nil(t, err)

	
	res, err = membershipProvider.AmMemberOf("test1", getAccessPolicy([]string{"peer2", "peer3"}))
	assert.False(t, res)
	assert.Nil(t, err)

	
	res, err = membershipProvider.AmMemberOf("test1", nil)
	assert.False(t, res)
	assert.Nil(t, err)

	
	res, err = membershipProvider.AmMemberOf("test1", getBadAccessPolicy([]string{"signer0"}, 1))
	assert.False(t, res)
	assert.Nil(t, err)
}

func getAccessPolicy(signers []string) *common.CollectionPolicyConfig {
	var data [][]byte
	for _, signer := range signers {
		data = append(data, []byte(signer))
	}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), data)
	return createCollectionPolicyConfig(policyEnvelope)
}

func getBadAccessPolicy(signers []string, badIndex int32) *common.CollectionPolicyConfig {
	var data [][]byte
	for _, signer := range signers {
		data = append(data, []byte(signer))
	}
	
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(badIndex)), data)
	return createCollectionPolicyConfig(policyEnvelope)
}
