/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"testing"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger"
	"github.com/stretchr/testify/assert"
)

func TestGetMissingPvtDataAfterRollback(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	h := env.newTestHelperCreateLgr("ledger1", t)

	collConf := []*collConf{{name: "coll1", btl: 5}}

	
	h.simulateDeployTx("cc1", collConf)
	h.cutBlockAndCommitLegacy()

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key1", "value1")
	})
	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key2", "value2")
	})

	h.causeMissingPvtData(0)
	blk2 := h.cutBlockAndCommitLegacy()

	h.verifyPvtState("cc1", "coll1", "key2", "value2") 
	h.simulateDataTx("", func(s *simulator) {
		h.assertError(s.GetPrivateData("cc1", "coll1", "key1")) 
	})

	
	h.verifyBlockAndPvtDataSameAs(2, blk2)
	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	expectedMissingPvtDataInfo.Add(2, 0, "cc1", "coll1")
	h.verifyMissingPvtDataSameAs(2, expectedMissingPvtDataInfo)

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key3", "value2")
	})
	blk3 := h.cutBlockAndCommitLegacy()

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key3", "value2")
	})
	blk4 := h.cutBlockAndCommitLegacy()

	
	h.verifyMissingPvtDataSameAs(5, expectedMissingPvtDataInfo)

	
	h.verifyLedgerHeight(5)
	env.closeLedgerMgmt()
	err := kvledger.RollbackKVLedger(env.initializer.Config.RootFSPath, "ledger1", 2)
	assert.NoError(t, err)
	env.initLedgerMgmt()

	h = env.newTestHelperOpenLgr("ledger1", t)
	h.verifyLedgerHeight(3)

	
	h.verifyBlockAndPvtDataSameAs(2, blk2)
	
	
	h.verifyMissingPvtDataSameAs(5, nil)

	
	assert.NoError(t, h.lgr.CommitLegacy(blk3, &ledger.CommitOptions{}))
	
	
	h.verifyMissingPvtDataSameAs(5, nil)

	
	assert.NoError(t, h.lgr.CommitLegacy(blk4, &ledger.CommitOptions{}))
	
	
	h.verifyMissingPvtDataSameAs(5, expectedMissingPvtDataInfo)
}
