/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger"
	"github.com/stretchr/testify/assert"
)

func TestRollbackKVLedger(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	
	dataHelper := newSampleDataHelper(t)
	var blockchainsInfo []*common.BlockchainInfo

	h := env.newTestHelperCreateLgr("testLedger", t)
	
	dataHelper.populateLedger(h)
	dataHelper.verifyLedgerContent(h)
	bcInfo, err := h.lgr.GetBlockchainInfo()
	assert.NoError(t, err)
	blockchainsInfo = append(blockchainsInfo, bcInfo)
	env.closeLedgerMgmt()

	
	err = kvledger.RollbackKVLedger(env.initializer.Config.RootFSPath, "noLedger", 0)
	assert.Equal(t, "ledgerID [noLedger] does not exist", err.Error())
	err = kvledger.RollbackKVLedger(env.initializer.Config.RootFSPath, "testLedger", bcInfo.Height)
	expectedErr := fmt.Sprintf("target block number [%d] should be less than the biggest block number [%d]",
		bcInfo.Height, bcInfo.Height-1)
	assert.Equal(t, expectedErr, err.Error())

	
	targetBlockNum := bcInfo.Height - 3
	err = kvledger.RollbackKVLedger(env.initializer.Config.RootFSPath, "testLedger", targetBlockNum)
	assert.NoError(t, err)
	rebuildable := rebuildableStatedb + rebuildableBookkeeper + rebuildableConfigHistory + rebuildableHistoryDB
	env.verifyRebuilableDoesNotExist(rebuildable)
	env.initLedgerMgmt()
	preResetHt, err := kvledger.LoadPreResetHeight(env.initializer.Config.RootFSPath)
	assert.Equal(t, bcInfo.Height, preResetHt["testLedger"])
	t.Logf("preResetHt = %#v", preResetHt)

	h = env.newTestHelperOpenLgr("testLedger", t)
	h.verifyLedgerHeight(targetBlockNum + 1)
	targetBlockNumIndex := targetBlockNum - 1
	for _, b := range dataHelper.submittedData["testLedger"].Blocks[targetBlockNumIndex+1:] {
		
		
		assert.NoError(t, h.lgr.CommitLegacy(b, &ledger.CommitOptions{FetchPvtDataFromLedger: true}))
	}
	actualBcInfo, err := h.lgr.GetBlockchainInfo()
	assert.Equal(t, bcInfo, actualBcInfo)
	dataHelper.verifyLedgerContent(h)
	
}

func TestRollbackKVLedgerWithBTL(t *testing.T) {
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

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "someOtherKey", "someOtherVal")
		s.setPvtdata("cc1", "coll2", "someOtherKey", "someOtherVal")
	})
	h.cutBlockAndCommitLegacy()
	env.closeLedgerMgmt()

	
	err := kvledger.RollbackKVLedger(env.initializer.Config.RootFSPath, "ledger1", 4)
	assert.NoError(t, err)
	rebuildable := rebuildableStatedb | rebuildableBookkeeper | rebuildableConfigHistory | rebuildableHistoryDB
	env.verifyRebuilableDoesNotExist(rebuildable)

	env.initLedgerMgmt()
	h = env.newTestHelperOpenLgr("ledger1", t)
	h.verifyPvtState("cc1", "coll1", "key1", "value1")                  
	h.verifyPvtState("cc1", "coll2", "key2", "")                        
	h.verifyBlockAndPvtData(2, nil, func(r *retrievedBlockAndPvtdata) { 
		r.pvtdataShouldContain(0, "cc1", "coll1", "key1", "value1") 
		r.pvtdataShouldNotContain("cc1", "coll2")                   
	})

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "someOtherKey", "someOtherVal")
		s.setPvtdata("cc1", "coll2", "someOtherKey", "someOtherVal")
	})
	h.cutBlockAndCommitLegacy()

	
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
	h.verifyBlockAndPvtData(6, nil, func(r *retrievedBlockAndPvtdata) { 
		r.pvtdataShouldContain(0, "cc1", "coll1", "key3", "value1") 
		r.pvtdataShouldNotContain("cc1", "coll2")                   
	})
}
