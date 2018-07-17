/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mgmt

import (
	"testing"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/msp"
	"github.com/stretchr/testify/assert"
)

func TestGetManagerForChains(t *testing.T) {
	
	mspMgr1 := GetManagerForChain("test")
	
	if mspMgr1 == nil {
		t.Fatal("mspMgr1 fail")
	}

	
	mspMgr2 := GetManagerForChain("test")
	
	if mspMgr2 != mspMgr1 {
		t.Fatal("mspMgr2 != mspMgr1 fail")
	}
}

func TestGetManagerForChains_usingMSPConfigHandlers(t *testing.T) {
	XXXSetMSPManager("foo", msp.NewMSPManager())
	msp2 := GetManagerForChain("foo")
	
	if msp2 == nil {
		t.FailNow()
	}
}

func TestGetIdentityDeserializer(t *testing.T) {
	XXXSetMSPManager("baz", msp.NewMSPManager())
	ids := GetIdentityDeserializer("baz")
	assert.NotNil(t, ids)
	ids = GetIdentityDeserializer("")
	assert.NotNil(t, ids)
}

func TestGetLocalSigningIdentityOrPanic(t *testing.T) {
	sid := GetLocalSigningIdentityOrPanic()
	assert.NotNil(t, sid)
}

func TestUpdateLocalMspCache(t *testing.T) {
	
	localMsp = nil

	
	firstMsp := GetLocalMSP()
	
	secondMsp := GetLocalMSP()
	
	thirdMsp := GetLocalMSP()

	
	if thirdMsp != secondMsp {
		t.Fatalf("thirdMsp != secondMsp")
	}
	
	if firstMsp != secondMsp {
		t.Fatalf("firstMsp != secondMsp")
	}
}

func TestNewMSPMgmtMgr(t *testing.T) {
	err := LoadMSPSetupForTesting()
	assert.Nil(t, err)

	
	mspMgmtMgr := GetManagerForChain("fake")

	id := GetLocalSigningIdentityOrPanic()
	assert.NotNil(t, id)

	serializedID, err := id.Serialize()
	if err != nil {
		t.Fatalf("Serialize should have succeeded, got err %s", err)
		return
	}

	idBack, err := mspMgmtMgr.DeserializeIdentity(serializedID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "channel doesn't exist")
	assert.Nil(t, idBack, "deserialized identity should have been nil")

	
	mspMgmtMgr = GetManagerForChain(util.GetTestChainID())

	id = GetLocalSigningIdentityOrPanic()
	assert.NotNil(t, id)

	serializedID, err = id.Serialize()
	if err != nil {
		t.Fatalf("Serialize should have succeeded, got err %s", err)
		return
	}

	idBack, err = mspMgmtMgr.DeserializeIdentity(serializedID)
	assert.NoError(t, err)
	assert.NotNil(t, idBack, "deserialized identity should not have been nil")
}

func LoadMSPSetupForTesting() error {
	dir, err := configtest.GetDevMspDir()
	if err != nil {
		return err
	}
	conf, err := msp.GetLocalMspConfig(dir, nil, "SampleOrg")
	if err != nil {
		return err
	}

	err = GetLocalMSP().Setup(conf)
	if err != nil {
		return err
	}

	err = GetManagerForChain(util.GetTestChainID()).Setup([]msp.MSP{GetLocalMSP()})
	if err != nil {
		return err
	}

	return nil
}
