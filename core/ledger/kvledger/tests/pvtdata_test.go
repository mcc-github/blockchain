/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"testing"
)

func TestMissingCollConfig(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	h := env.newTestHelperCreateLgr("ledger1", t)

	collConf := []*collConf{{name: "coll1", btl: 5}}

	
	h.simulateDeployTx("cc1", nil)
	h.cutBlockAndCommitWithPvtdata()

	
	h.simulateDataTx("", func(s *simulator) {
		h.assertError(s.GetPrivateData("cc1", "coll1", "key"))
		h.assertError(s.SetPrivateData("cc1", "coll1", "key", []byte("value")))
		h.assertError(s.DeletePrivateData("cc1", "coll1", "key"))
	})

	
	h.simulateUpgradeTx("cc1", collConf)
	h.cutBlockAndCommitWithPvtdata()

	
	
	h.simulateDataTx("", func(s *simulator) {
		h.assertNoError(s.GetPrivateData("cc1", "coll1", "key1"))
		h.assertNoError(s.SetPrivateData("cc1", "coll1", "key2", []byte("value")))
		h.assertNoError(s.DeletePrivateData("cc1", "coll1", "key3"))
		h.assertError(s.GetPrivateData("cc1", "coll2", "key"))
		h.assertError(s.SetPrivateData("cc1", "coll2", "key", []byte("value")))
		h.assertError(s.DeletePrivateData("cc1", "coll2", "key"))
	})
}

func TestTxWithMissingPvtdata(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	h := env.newTestHelperCreateLgr("ledger1", t)

	collConf := []*collConf{{name: "coll1", btl: 5}}

	
	h.simulateDeployTx("cc1", collConf)
	h.cutBlockAndCommitWithPvtdata()

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key1", "value1")
	})
	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key2", "value2")
	})
	h.simulatedTrans[0].Pvtws = nil 
	blk2 := h.cutBlockAndCommitWithPvtdata()

	h.verifyPvtState("cc1", "coll1", "key2", "value2") 
	h.simulateDataTx("", func(s *simulator) {
		h.assertError(s.GetPrivateData("cc1", "coll1", "key1")) 
	})

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key1", "newvalue1")
	})
	blk3 := h.cutBlockAndCommitWithPvtdata()
	h.verifyPvtState("cc1", "coll1", "key1", "newvalue1") 
	h.verifyBlockAndPvtDataSameAs(2, blk2)
	h.verifyBlockAndPvtDataSameAs(3, blk3)
}

func TestTxWithWrongPvtdata(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	h := env.newTestHelperCreateLgr("ledger1", t)

	collConf := []*collConf{{name: "coll1", btl: 5}}

	
	h.simulateDeployTx("cc1", collConf)
	h.cutBlockAndCommitWithPvtdata()

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key1", "value1")
	})
	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key2", "value2")
	})
	h.simulatedTrans[0].Pvtws = h.simulatedTrans[1].Pvtws 
	
	h.cutBlockAndCommitExpectError()
	h.verifyPvtState("cc1", "coll1", "key2", "")
}

func TestBTL(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	h := env.newTestHelperCreateLgr("ledger1", t)
	collConf := []*collConf{{name: "coll1", btl: 0}, {name: "coll2", btl: 5}}

	
	h.simulateDeployTx("cc1", collConf)
	h.cutBlockAndCommitWithPvtdata()

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key1", "value1") 
		s.setPvtdata("cc1", "coll2", "key2", "value2") 
	})
	blk2 := h.cutBlockAndCommitWithPvtdata()

	
	for i := 0; i < 5; i++ {
		h.simulateDataTx("", func(s *simulator) {
			s.setPvtdata("cc1", "coll1", "someOtherKey", "someOtherVal")
			s.setPvtdata("cc1", "coll2", "someOtherKey", "someOtherVal")
		})
		h.cutBlockAndCommitWithPvtdata()
	}

	
	h.verifyPvtState("cc1", "coll1", "key1", "value1") 
	h.verifyPvtState("cc1", "coll2", "key2", "value2") 
	h.verifyBlockAndPvtDataSameAs(2, blk2)             

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "someOtherKey", "someOtherVal")
		s.setPvtdata("cc1", "coll2", "someOtherKey", "someOtherVal")
	})
	h.cutBlockAndCommitWithPvtdata()

	
	h.verifyPvtState("cc1", "coll1", "key1", "value1")                  
	h.verifyPvtState("cc1", "coll2", "key2", "")                        
	h.verifyBlockAndPvtData(2, nil, func(r *retrievedBlockAndPvtdata) { 
		r.pvtdataShouldContain(0, "cc1", "coll1", "key1", "value1") 
		r.pvtdataShouldNotContain("cc1", "coll2")                   
	})
}
