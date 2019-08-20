/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"testing"

	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/msp"
	mspprotos "github.com/mcc-github/blockchain/protos/msp"
	"github.com/stretchr/testify/assert"
)

func TestMSPConfigManager(t *testing.T) {
	mspDir, err := configtest.GetDevMspDir()
	assert.NoError(t, err)
	conf, err := msp.GetLocalMspConfig(mspDir, nil, "SampleOrg")
	assert.NoError(t, err)

	

	mspVers := []msp.MSPVersion{msp.MSPv1_0, msp.MSPv1_1}

	for _, ver := range mspVers {
		mspCH := NewMSPConfigHandler(ver, factory.DefaultBCCSP)

		_, err = mspCH.ProposeMSP(conf)
		assert.NoError(t, err)

		mgr, err := mspCH.CreateMSPManager()
		assert.NoError(t, err)
		assert.NotNil(t, mgr)

		msps, err := mgr.GetMSPs()
		assert.NoError(t, err)

		if msps == nil || len(msps) == 0 {
			t.Fatalf("There are no MSPS in the manager")
		}

		for _, mspInst := range msps {
			assert.Equal(t, mspInst.GetVersion(), msp.MSPVersion(ver))
		}
	}
}

func TestMSPConfigFailure(t *testing.T) {
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	mspCH := NewMSPConfigHandler(msp.MSPv1_0, cryptoProvider)

	
	t.Run("Bad proto", func(t *testing.T) {
		_, err := mspCH.ProposeMSP(&mspprotos.MSPConfig{Config: []byte("BARF!")})
		assert.Error(t, err)
	})

	t.Run("Bad MSP Type", func(t *testing.T) {
		_, err := mspCH.ProposeMSP(&mspprotos.MSPConfig{Type: int32(10)})
		assert.Error(t, err)
	})
}
