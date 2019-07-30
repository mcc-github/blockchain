/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"testing"
)

func TestRebuildComponents(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()

	h1, h2 := env.newTestHelperCreateLgr("ledger1", t), env.newTestHelperCreateLgr("ledger2", t)
	dataHelper := newSampleDataHelper(t)

	dataHelper.populateLedger(h1)
	dataHelper.populateLedger(h2)

	dataHelper.verifyLedgerContent(h1)
	dataHelper.verifyLedgerContent(h2)

	t.Run("rebuild only statedb",
		func(t *testing.T) {
			env.closeAllLedgersAndDrop(rebuildableStatedb)
			h1, h2 := env.newTestHelperOpenLgr("ledger1", t), env.newTestHelperOpenLgr("ledger2", t)
			dataHelper.verifyLedgerContent(h1)
			dataHelper.verifyLedgerContent(h2)
		},
	)

	t.Run("rebuild statedb and config history",
		func(t *testing.T) {
			env.closeAllLedgersAndDrop(rebuildableStatedb | rebuildableConfigHistory)
			h1, h2 := env.newTestHelperOpenLgr("ledger1", t), env.newTestHelperOpenLgr("ledger2", t)
			dataHelper.verifyLedgerContent(h1)
			dataHelper.verifyLedgerContent(h2)
		},
	)

	t.Run("rebuild statedb and block index",
		func(t *testing.T) {
			env.closeAllLedgersAndDrop(rebuildableStatedb | rebuildableBlockIndex)
			h1, h2 := env.newTestHelperOpenLgr("ledger1", t), env.newTestHelperOpenLgr("ledger2", t)
			dataHelper.verifyLedgerContent(h1)
			dataHelper.verifyLedgerContent(h2)
		},
	)

	t.Run("rebuild statedb and historydb",
		func(t *testing.T) {
			env.closeAllLedgersAndDrop(rebuildableStatedb | rebuildableHistoryDB)
			h1, h2 := env.newTestHelperOpenLgr("ledger1", t), env.newTestHelperOpenLgr("ledger2", t)
			dataHelper.verifyLedgerContent(h1)
			dataHelper.verifyLedgerContent(h2)
		},
	)
}

func TestRebuildComponentsWithBTL(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	h := env.newTestHelperCreateLgr("ledger1", t)
	collConf := []*collConf{{name: "coll1", btl: 0}, {name: "coll2", btl: 1}}

	
	h.simulateDeployTx("cc1", collConf)
	h.cutBlockAndCommitLegacy()

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key1", "value1") 
		s.setPvtdata("cc1", "coll2", "key2", "value2") 
	})
	blk2 := h.cutBlockAndCommitLegacy()

	
	h.verifyPvtState("cc1", "coll1", "key1", "value1") 
	h.verifyPvtState("cc1", "coll2", "key2", "value2") 
	h.verifyBlockAndPvtDataSameAs(2, blk2)             

	

	for i := 0; i < 2; i++ {
		h.simulateDataTx("", func(s *simulator) {
			s.setPvtdata("cc1", "coll1", "someOtherKey", "someOtherVal")
			s.setPvtdata("cc1", "coll2", "someOtherKey", "someOtherVal")
		})
		h.cutBlockAndCommitLegacy()
	}

	
	h.verifyPvtState("cc1", "coll1", "key1", "value1")                  
	h.verifyPvtState("cc1", "coll2", "key2", "")                        
	h.verifyBlockAndPvtData(2, nil, func(r *retrievedBlockAndPvtdata) { 
		r.pvtdataShouldContain(0, "cc1", "coll1", "key1", "value1") 
		r.pvtdataShouldNotContain("cc1", "coll2")                   
	})

	
	env.closeAllLedgersAndDrop(rebuildableStatedb | rebuildableBookkeeper)

	h = env.newTestHelperOpenLgr("ledger1", t)
	h.verifyPvtState("cc1", "coll1", "key1", "value1")                  
	h.verifyPvtState("cc1", "coll2", "key2", "")                        
	h.verifyBlockAndPvtData(2, nil, func(r *retrievedBlockAndPvtdata) { 
		r.pvtdataShouldContain(0, "cc1", "coll1", "key1", "value1") 
		r.pvtdataShouldNotContain("cc1", "coll2")                   
	})

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key3", "value1") 
		s.setPvtdata("cc1", "coll2", "key4", "value2") 
	})
	h.cutBlockAndCommitLegacy()

	
	for i := 0; i < 2; i++ {
		h.simulateDataTx("", func(s *simulator) {
			s.setPvtdata("cc1", "coll1", "someOtherKey", "someOtherVal")
			s.setPvtdata("cc1", "coll2", "someOtherKey", "someOtherVal")
		})
		h.cutBlockAndCommitLegacy()
	}

	
	h.verifyPvtState("cc1", "coll1", "key3", "value1")                  
	h.verifyPvtState("cc1", "coll2", "key4", "")                        
	h.verifyBlockAndPvtData(5, nil, func(r *retrievedBlockAndPvtdata) { 
		r.pvtdataShouldContain(0, "cc1", "coll1", "key3", "value1") 
		r.pvtdataShouldNotContain("cc1", "coll2")                   
	})
}
