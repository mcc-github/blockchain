/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased_test

import (
	"testing"

	"github.com/mcc-github/blockchain/common/cauthdsl"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/chaincode/shim/ext/statebased"
	"github.com/stretchr/testify/assert"
)

func TestAddOrg(t *testing.T) {
	
	ep, err := statebased.NewStateEP(nil)
	assert.NoError(t, err)
	err = ep.AddOrgs(statebased.RoleTypePeer, "Org1")
	assert.NoError(t, err)
	err = ep.AddOrgs("unknown", "Org1")
	assert.Error(t, err)
	rtdnee, ok := err.(*statebased.RoleTypeDoesNotExistError)
	assert.True(t, ok)
	assert.Equal(t, statebased.RoleType("unknown"), rtdnee.RoleType)
	epBytes, err := ep.Policy()
	assert.NoError(t, err)
	expectedEP := cauthdsl.SignedByMspPeer("Org1")
	expectedEPBytes, err := proto.Marshal(expectedEP)
	assert.NoError(t, err)
	assert.Equal(t, expectedEPBytes, epBytes)
}

func TestListOrgs(t *testing.T) {
	expectedEP := cauthdsl.SignedByMspPeer("Org1")
	expectedEPBytes, err := proto.Marshal(expectedEP)
	assert.NoError(t, err)

	
	ep, err := statebased.NewStateEP(expectedEPBytes)
	orgs := ep.ListOrgs()
	assert.Equal(t, []string{"Org1"}, orgs)
}

func TestDelAddOrg(t *testing.T) {
	expectedEP := cauthdsl.SignedByMspPeer("Org1")
	expectedEPBytes, err := proto.Marshal(expectedEP)
	assert.NoError(t, err)
	ep, err := statebased.NewStateEP(expectedEPBytes)

	
	orgs := ep.ListOrgs()
	assert.Equal(t, []string{"Org1"}, orgs)

	
	ep.AddOrgs(statebased.RoleTypePeer, "Org2")
	ep.DelOrgs("Org1")

	
	epBytes, err := ep.Policy()
	assert.NoError(t, err)
	expectedEP = cauthdsl.SignedByMspPeer("Org2")
	expectedEPBytes, err = proto.Marshal(expectedEP)
	assert.NoError(t, err)
	assert.Equal(t, expectedEPBytes, epBytes)
}
