/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package msp

import (
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/stretchr/testify/assert"
)

func TestInvalidAdminNodeOU(t *testing.T) {
	
	
	
	thisMSP, err := getLocalMSPWithVersionAndError(t, "testdata/nodeous1", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)
	assert.Error(t, err)

	
	thisMSP, err = getLocalMSPWithVersionAndError(t, "testdata/nodeous1", MSPv1_0)
	assert.False(t, thisMSP.(*bccspmsp).ouEnforcement)
	assert.NoError(t, err)
}

func TestInvalidSigningIdentityNodeOU(t *testing.T) {
	
	
	
	thisMSP := getLocalMSPWithVersion(t, "testdata/nodeous2", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)

	id, err := thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	err = id.Validate()
	assert.Error(t, err)

	
	thisMSP, err = getLocalMSPWithVersionAndError(t, "testdata/nodeous1", MSPv1_0)
	assert.False(t, thisMSP.(*bccspmsp).ouEnforcement)
	assert.NoError(t, err)

	id, err = thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	err = id.Validate()
	assert.NoError(t, err)
}

func TestValidMSPWithNodeOU(t *testing.T) {
	
	
	thisMSP := getLocalMSPWithVersion(t, "testdata/nodeous3", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)

	id, err := thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	err = id.Validate()
	assert.NoError(t, err)

	
	thisMSP = getLocalMSPWithVersion(t, "testdata/nodeous3", MSPv1_0)
	assert.False(t, thisMSP.(*bccspmsp).ouEnforcement)

	id, err = thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	err = id.Validate()
	assert.NoError(t, err)
}

func TestValidMSPWithNodeOUAndOrganizationalUnits(t *testing.T) {
	
	
	thisMSP := getLocalMSPWithVersion(t, "testdata/nodeous6", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)

	id, err := thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	err = id.Validate()
	assert.NoError(t, err)

	
	thisMSP = getLocalMSPWithVersion(t, "testdata/nodeous6", MSPv1_0)
	assert.False(t, thisMSP.(*bccspmsp).ouEnforcement)

	id, err = thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	err = id.Validate()
	assert.NoError(t, err)
}

func TestInvalidMSPWithNodeOUAndOrganizationalUnits(t *testing.T) {
	
	
	
	
	thisMSP, err := getLocalMSPWithVersionAndError(t, "testdata/nodeous7", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "could not validate identity's OUs: none of the identity's organizational units")
	}

	
	thisMSP, err = getLocalMSPWithVersionAndError(t, "testdata/nodeous7", MSPv1_0)
	assert.False(t, thisMSP.(*bccspmsp).ouEnforcement)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "could not validate identity's OUs: none of the identity's organizational units")
	}
}

func TestInvalidAdminOU(t *testing.T) {
	
	
	thisMSP, err := getLocalMSPWithVersionAndError(t, "testdata/nodeous4", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "admin 0 is invalid: The identity is not valid under this MSP [SampleOrg]: could not validate identity's OUs: certifiersIdentifier does not match")

	
	thisMSP, err = getLocalMSPWithVersionAndError(t, "testdata/nodeous4", MSPv1_0)
	assert.False(t, thisMSP.(*bccspmsp).ouEnforcement)
	assert.NoError(t, err)
}

func TestInvalidAdminOUNotAClient(t *testing.T) {
	
	
	thisMSP, err := getLocalMSPWithVersionAndError(t, "testdata/nodeous8", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The identity does not contain OU [CLIENT]")

	
	thisMSP, err = getLocalMSPWithVersionAndError(t, "testdata/nodeous8", MSPv1_0)
	assert.False(t, thisMSP.(*bccspmsp).ouEnforcement)
	assert.NoError(t, err)
}

func TestSatisfiesPrincipalPeer(t *testing.T) {
	
	
	thisMSP := getLocalMSPWithVersion(t, "testdata/nodeous3", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)

	
	id, err := thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	err = id.Validate()
	assert.NoError(t, err)

	assert.True(t, t.Run("Check that id is a peer", func(t *testing.T) {
		
		mspID, err := thisMSP.GetIdentifier()
		assert.NoError(t, err)
		principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_PEER, MspIdentifier: mspID})
		assert.NoError(t, err)
		principal := &msp.MSPPrincipal{
			PrincipalClassification: msp.MSPPrincipal_ROLE,
			Principal:               principalBytes}
		err = id.SatisfiesPrincipal(principal)
		assert.NoError(t, err)
	}))

	assert.True(t, t.Run("Check that id is not a client", func(t *testing.T) {
		
		mspID, err := thisMSP.GetIdentifier()
		assert.NoError(t, err)
		principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_CLIENT, MspIdentifier: mspID})
		assert.NoError(t, err)
		principal := &msp.MSPPrincipal{
			PrincipalClassification: msp.MSPPrincipal_ROLE,
			Principal:               principalBytes}
		err = id.SatisfiesPrincipal(principal)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "The identity is not a [CLIENT] under this MSP [SampleOrg]")
	}))
}

func TestSatisfiesPrincipalClient(t *testing.T) {
	
	
	thisMSP := getLocalMSPWithVersion(t, "testdata/nodeous3", MSPv1_1)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)

	
	assert.Equal(t, 1, len(thisMSP.(*bccspmsp).admins))
	id := thisMSP.(*bccspmsp).admins[0]

	err := id.Validate()
	assert.NoError(t, err)

	
	assert.True(t, t.Run("Check that id is a client", func(t *testing.T) {
		mspID, err := thisMSP.GetIdentifier()
		assert.NoError(t, err)
		principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_CLIENT, MspIdentifier: mspID})
		assert.NoError(t, err)
		principal := &msp.MSPPrincipal{
			PrincipalClassification: msp.MSPPrincipal_ROLE,
			Principal:               principalBytes}
		err = id.SatisfiesPrincipal(principal)
		assert.NoError(t, err)
	}))

	assert.True(t, t.Run("Check that id is not a peer", func(t *testing.T) {
		
		mspID, err := thisMSP.GetIdentifier()
		assert.NoError(t, err)
		principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_PEER, MspIdentifier: mspID})
		assert.NoError(t, err)
		principal := &msp.MSPPrincipal{
			PrincipalClassification: msp.MSPPrincipal_ROLE,
			Principal:               principalBytes}
		err = id.SatisfiesPrincipal(principal)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "The identity is not a [PEER] under this MSP [SampleOrg]")
	}))
}

func TestSatisfiesPrincipalAdmin(t *testing.T) {
	
	
	thisMSP := getLocalMSPWithVersion(t, "testdata/nodeouadmin", MSPv1_4_2)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)

	cert, err := readFile("testdata/nodeouadmin/adm/testadmincert.pem")
	assert.NoError(t, err)

	id, _, err := thisMSP.(*bccspmsp).getIdentityFromConf(cert)
	assert.NoError(t, err)

	principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_ADMIN, MspIdentifier: "SampleOrg"})
	assert.NoError(t, err)
	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ROLE,
		Principal:               principalBytes}
	err = id.SatisfiesPrincipal(principal)
	assert.NoError(t, err)
}

func TestLoad142MSPWithInvalidAdminConfiguration(t *testing.T) {
	
	
	conf, err := GetLocalMspConfig("testdata/nodeouadmin2", nil, "SampleOrg")
	assert.NoError(t, err)

	ks, err := sw.NewFileBasedKeyStore(nil, filepath.Join("testdata/nodeouadmin2", "keystore"), true)
	assert.NoError(t, err)
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	thisMSP, err := NewBccspMspWithKeyStore(MSPv1_4_2, ks, cryptoProvider)
	assert.NoError(t, err)

	err = thisMSP.Setup(conf)
	assert.Error(t, err)
	assert.Equal(t, "administrators must be declared when no admin ou classification is set", err.Error())

	
	
	conf, err = GetLocalMspConfig("testdata/nodeouadmin3", nil, "SampleOrg")
	assert.NoError(t, err)

	ks, err = sw.NewFileBasedKeyStore(nil, filepath.Join("testdata/nodeouadmin3", "keystore"), true)
	assert.NoError(t, err)
	thisMSP, err = NewBccspMspWithKeyStore(MSPv1_4_2, ks, cryptoProvider)
	assert.NoError(t, err)

	err = thisMSP.Setup(conf)
	assert.Error(t, err)
	assert.Equal(t, "administrators must be declared when no admin ou classification is set", err.Error())
}

func TestSatisfiesPrincipalOrderer(t *testing.T) {
	
	
	thisMSP := getLocalMSPWithVersion(t, "testdata/nodeouorderer", MSPv1_4_2)
	assert.True(t, thisMSP.(*bccspmsp).ouEnforcement)

	id, err := thisMSP.(*bccspmsp).GetDefaultSigningIdentity()
	assert.NoError(t, err)

	principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_ORDERER, MspIdentifier: "SampleOrg"})
	assert.NoError(t, err)
	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ROLE,
		Principal:               principalBytes}
	err = id.SatisfiesPrincipal(principal)
	assert.NoError(t, err)
}

func TestLoad142MSPWithInvalidOrdererConfiguration(t *testing.T) {
	
	
	conf, err := GetLocalMspConfig("testdata/nodeouorderer2", nil, "SampleOrg")
	assert.NoError(t, err)

	ks, err := sw.NewFileBasedKeyStore(nil, filepath.Join("testdata/nodeouorderer2", "keystore"), true)
	assert.NoError(t, err)
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	thisMSP, err := NewBccspMspWithKeyStore(MSPv1_4_2, ks, cryptoProvider)
	assert.NoError(t, err)

	err = thisMSP.Setup(conf)
	assert.NoError(t, err)
	id, err := thisMSP.(*bccspmsp).GetDefaultSigningIdentity()
	assert.NoError(t, err)

	principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_ORDERER, MspIdentifier: "SampleOrg"})
	assert.NoError(t, err)
	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ROLE,
		Principal:               principalBytes}
	err = id.SatisfiesPrincipal(principal)
	assert.Error(t, err)
	assert.Equal(t, "The identity is not a [ORDERER] under this MSP [SampleOrg]: cannot test for classification, node ou for type [ORDERER], not defined, msp: [SampleOrg]", err.Error())

	
	
	conf, err = GetLocalMspConfig("testdata/nodeouorderer3", nil, "SampleOrg")
	assert.NoError(t, err)

	ks, err = sw.NewFileBasedKeyStore(nil, filepath.Join("testdata/nodeouorderer3", "keystore"), true)
	assert.NoError(t, err)
	thisMSP, err = NewBccspMspWithKeyStore(MSPv1_4_2, ks, cryptoProvider)
	assert.NoError(t, err)

	err = thisMSP.Setup(conf)
	assert.NoError(t, err)
	id, err = thisMSP.(*bccspmsp).GetDefaultSigningIdentity()
	assert.NoError(t, err)

	principalBytes, err = proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_ORDERER, MspIdentifier: "SampleOrg"})
	assert.NoError(t, err)
	principal = &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ROLE,
		Principal:               principalBytes}
	err = id.SatisfiesPrincipal(principal)
	assert.Error(t, err)
	assert.Equal(t, "The identity is not a [ORDERER] under this MSP [SampleOrg]: cannot test for classification, node ou for type [ORDERER], not defined, msp: [SampleOrg]", err.Error())
}

func TestValidMSPWithNodeOUMissingClassification(t *testing.T) {
	
	
	_, err := getLocalMSPWithVersionAndError(t, "testdata/nodeousbadconf1", MSPv1_3)
	assert.Error(t, err)
	assert.Equal(t, "Failed setting up NodeOUs. ClientOU must be different from nil.", err.Error())

	_, err = getLocalMSPWithVersionAndError(t, "testdata/nodeousbadconf1", MSPv1_4_2)
	assert.Error(t, err)
	assert.Equal(t, "admin 0 is invalid [cannot test for classification, node ou for type [CLIENT], not defined, msp: [SampleOrg],The identity does not contain OU [ADMIN], MSP: [SampleOrg]]", err.Error())

	
	
	_, err = getLocalMSPWithVersionAndError(t, "testdata/nodeousbadconf2", MSPv1_3)
	assert.Error(t, err)
	assert.Equal(t, "Failed setting up NodeOUs. PeerOU must be different from nil.", err.Error())

	_, err = getLocalMSPWithVersionAndError(t, "testdata/nodeousbadconf2", MSPv1_4_2)
	assert.NoError(t, err)
}
