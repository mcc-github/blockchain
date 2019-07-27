/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/chaincode/shim/ext/statebased"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/stretchr/testify/assert"
)

func TestAddOrg(t *testing.T) {
	
	ep, err := statebased.NewStateEP(nil)
	assert.NoError(t, err)
	err = ep.AddOrgs(statebased.RoleTypePeer, "Org1")
	assert.NoError(t, err)

	
	err = ep.AddOrgs("unknown", "Org1")
	assert.Equal(t, &statebased.RoleTypeDoesNotExistError{RoleType: statebased.RoleType("unknown")}, err)
	assert.EqualError(t, err, "role type unknown does not exist")

	epBytes, err := ep.Policy()
	assert.NoError(t, err)
	expectedEP := signedByMspPeer("Org1", t)
	expectedEPBytes, err := proto.Marshal(expectedEP)
	assert.NoError(t, err)
	assert.Equal(t, expectedEPBytes, epBytes)
}

func TestListOrgs(t *testing.T) {
	expectedEP := signedByMspPeer("Org1", t)
	expectedEPBytes, err := proto.Marshal(expectedEP)
	assert.NoError(t, err)

	
	ep, err := statebased.NewStateEP(expectedEPBytes)
	orgs := ep.ListOrgs()
	assert.Equal(t, []string{"Org1"}, orgs)
}

func TestDelAddOrg(t *testing.T) {
	expectedEP := signedByMspPeer("Org1", t)
	expectedEPBytes, err := proto.Marshal(expectedEP)
	assert.NoError(t, err)
	ep, err := statebased.NewStateEP(expectedEPBytes)

	
	orgs := ep.ListOrgs()
	assert.Equal(t, []string{"Org1"}, orgs)

	
	ep.AddOrgs(statebased.RoleTypePeer, "Org2")
	ep.DelOrgs("Org1")

	
	epBytes, err := ep.Policy()
	assert.NoError(t, err)
	expectedEP = signedByMspPeer("Org2", t)
	expectedEPBytes, err = proto.Marshal(expectedEP)
	assert.NoError(t, err)
	assert.Equal(t, expectedEPBytes, epBytes)
}



func signedByMspPeer(mspId string, t *testing.T) *common.SignaturePolicyEnvelope {
	
	principal, err := proto.Marshal(
		&msp.MSPRole{
			Role:          msp.MSPRole_PEER,
			MspIdentifier: mspId,
		},
	)
	if err != nil {
		t.Fatalf("failed to marshal principal: %s", err)
	}

	
	p := &common.SignaturePolicyEnvelope{
		Version: 0,
		Rule: &common.SignaturePolicy{
			Type: &common.SignaturePolicy_NOutOf_{
				NOutOf: &common.SignaturePolicy_NOutOf{
					N: 1,
					Rules: []*common.SignaturePolicy{
						{
							Type: &common.SignaturePolicy_SignedBy{
								SignedBy: 0,
							},
						},
					},
				},
			},
		},
		Identities: []*msp.MSPPrincipal{
			{
				PrincipalClassification: msp.MSPPrincipal_ROLE,
				Principal:               principal,
			},
		},
	}

	return p
}
