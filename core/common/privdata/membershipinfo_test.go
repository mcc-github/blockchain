/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"testing"

	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/stretchr/testify/assert"
)

func TestMembershipInfoProvider(t *testing.T) {
	
	peerSelfSignedData := common.SignedData{
		Identity:  []byte("peer0"),
		Signature: []byte{1, 2, 3},
		Data:      []byte{4, 5, 6},
	}

	collectionStore := NewSimpleCollectionStore(&mockStoreSupport{})

	
	membershipProvider := NewMembershipInfoProvider("test1", peerSelfSignedData, collectionStore)
	res, err := membershipProvider.AmMemberOf(getAccessPolicy([]string{"peer0", "peer1"}))
	assert.True(t, res)
	assert.Nil(t, err)

	
	res, err = membershipProvider.AmMemberOf(getAccessPolicy([]string{"peer2", "peer3"}))
	assert.False(t, res)
	assert.Nil(t, err)

	
	res, err = membershipProvider.AmMemberOf(nil)
	assert.False(t, res)
	assert.Error(t, err)
	assert.Equal(t, "Collection config policy is nil", err.Error())
}

func getAccessPolicy(signers []string) *common.CollectionPolicyConfig {
	var data [][]byte
	for _, signer := range signers {
		data = append(data, []byte(signer))
	}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), data)
	return createCollectionPolicyConfig(policyEnvelope)
}
