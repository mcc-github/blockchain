/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/blkstorage/fsblkstorage"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/stretchr/testify/assert"
)

func TestResetAllLedgers(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	
	dataHelper := newSampleDataHelper(t)
	var genesisBlocks []*common.Block
	var blockchainsInfo []*common.BlockchainInfo

	
	
	for i := 0; i < 10; i++ {
		h := env.newTestHelperCreateLgr(fmt.Sprintf("ledger-%d", i), t)
		dataHelper.populateLedger(h)
		dataHelper.verifyLedgerContent(h)
		gb, err := h.lgr.GetBlockByNumber(0)
		assert.NoError(t, err)
		genesisBlocks = append(genesisBlocks, gb)
		bcInfo, err := h.lgr.GetBlockchainInfo()
		assert.NoError(t, err)
		blockchainsInfo = append(blockchainsInfo, bcInfo)
	}
	env.closeLedgerMgmt()

	
	rootFSPath := env.initializer.Config.RootFSPath
	err := kvledger.ResetAllKVLedgers(rootFSPath)
	assert.NoError(t, err)
	rebuildable := rebuildableStatedb | rebuildableBookkeeper | rebuildableConfigHistory | rebuildableHistoryDB | rebuildableBlockIndex
	env.verifyRebuilableDoesNotExist(rebuildable)
	env.initLedgerMgmt()
	preResetHt, err := kvledger.LoadPreResetHeight(rootFSPath)
	t.Logf("preResetHt = %#v", preResetHt)
	
	
	
	
	
	for i := 0; i < 10; i++ {
		ledgerID := fmt.Sprintf("ledger-%d", i)
		h := env.newTestHelperOpenLgr(ledgerID, t)
		h.verifyLedgerHeight(1)
		assert.Equal(t, blockchainsInfo[i].Height, preResetHt[ledgerID])
		gb, err := h.lgr.GetBlockByNumber(0)
		assert.NoError(t, err)
		assert.Equal(t, genesisBlocks[i], gb)
		for _, b := range dataHelper.submittedData[ledgerID].Blocks {
			assert.NoError(t, h.lgr.CommitLegacy(b, &ledger.CommitOptions{}))
		}
		bcInfo, err := h.lgr.GetBlockchainInfo()
		assert.NoError(t, err)
		assert.Equal(t, blockchainsInfo[i], bcInfo)
		dataHelper.verifyLedgerContent(h)
	}

	assert.NoError(t, kvledger.ClearPreResetHeight(env.initializer.Config.RootFSPath))
	preResetHt, err = kvledger.LoadPreResetHeight(env.initializer.Config.RootFSPath)
	assert.NoError(t, err)
	assert.Len(t, preResetHt, 0)
}

func TestResetAllLedgersWithBTL(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	h := env.newTestHelperCreateLgr("ledger1", t)
	collConf := []*collConf{{name: "coll1", btl: 0}, {name: "coll2", btl: 1}}

	
	h.simulateDeployTx("cc1", collConf)
	blk1 := h.cutBlockAndCommitLegacy()

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "key1", "value1") 
		s.setPvtdata("cc1", "coll2", "key2", "value2") 
	})
	blk2 := h.cutBlockAndCommitLegacy()

	
	h.verifyPvtState("cc1", "coll1", "key1", "value1") 
	h.verifyPvtState("cc1", "coll2", "key2", "value2") 
	h.verifyBlockAndPvtDataSameAs(2, blk2)             

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "someOtherKey", "someOtherVal")
		s.setPvtdata("cc1", "coll2", "someOtherKey", "someOtherVal")
	})
	blk3 := h.cutBlockAndCommitLegacy()

	
	h.simulateDataTx("", func(s *simulator) {
		s.setPvtdata("cc1", "coll1", "someOtherKey", "someOtherVal")
		s.setPvtdata("cc1", "coll2", "someOtherKey", "someOtherVal")
	})
	blk4 := h.cutBlockAndCommitLegacy()

	
	h.verifyPvtState("cc1", "coll1", "key1", "value1")                  
	h.verifyPvtState("cc1", "coll2", "key2", "")                        
	h.verifyBlockAndPvtData(2, nil, func(r *retrievedBlockAndPvtdata) { 
		r.pvtdataShouldContain(0, "cc1", "coll1", "key1", "value1") 
		r.pvtdataShouldNotContain("cc1", "coll2")                   
	})

	env.closeLedgerMgmt()

	
	err := kvledger.ResetAllKVLedgers(env.initializer.Config.RootFSPath)
	assert.NoError(t, err)
	rebuildable := rebuildableStatedb | rebuildableBookkeeper | rebuildableConfigHistory | rebuildableHistoryDB | rebuildableBlockIndex
	env.verifyRebuilableDoesNotExist(rebuildable)
	env.initLedgerMgmt()

	
	preResetHt, err := kvledger.LoadPreResetHeight(env.initializer.Config.RootFSPath)
	t.Logf("preResetHt = %#v", preResetHt)
	assert.Equal(t, uint64(5), preResetHt["ledger1"])
	h = env.newTestHelperOpenLgr("ledger1", t)
	h.verifyLedgerHeight(1)

	
	assert.NoError(t, h.lgr.CommitLegacy(blk1, &ledger.CommitOptions{}))
	assert.NoError(t, h.lgr.CommitLegacy(blk2, &ledger.CommitOptions{}))
	
	h.verifyPvtState("cc1", "coll1", "key1", "value1") 
	h.verifyPvtState("cc1", "coll2", "key2", "value2") 
	assert.NoError(t, h.lgr.CommitLegacy(blk3, &ledger.CommitOptions{}))
	assert.NoError(t, h.lgr.CommitLegacy(blk4, &ledger.CommitOptions{}))

	
	h.verifyPvtState("cc1", "coll1", "key1", "value1")                  
	h.verifyPvtState("cc1", "coll2", "key2", "")                        
	h.verifyBlockAndPvtData(2, nil, func(r *retrievedBlockAndPvtdata) { 
		r.pvtdataShouldContain(0, "cc1", "coll1", "key1", "value1") 
		r.pvtdataShouldNotContain("cc1", "coll2")                   
	})
}

func TestResetLedgerWithoutDroppingDBs(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	
	dataHelper := newSampleDataHelper(t)

	
	h := env.newTestHelperCreateLgr("ledger-1", t)
	dataHelper.populateLedger(h)
	dataHelper.verifyLedgerContent(h)
	env.closeLedgerMgmt()

	
	blockstorePath := filepath.Join(env.initializer.Config.RootFSPath, "chains")
	err := fsblkstorage.ResetBlockStore(blockstorePath)
	assert.NoError(t, err)
	rebuildable := rebuildableStatedb | rebuildableBookkeeper | rebuildableConfigHistory | rebuildableHistoryDB
	env.verifyRebuilablesExist(rebuildable)
	rebuildable = rebuildableBlockIndex
	env.verifyRebuilableDoesNotExist(rebuildable)
	env.initLedgerMgmt()
	preResetHt, err := kvledger.LoadPreResetHeight(env.initializer.Config.RootFSPath)
	t.Logf("preResetHt = %#v", preResetHt)
	_, err = env.ledgerMgr.OpenLedger("ledger-1")
	assert.Error(t, err)
	
	assert.EqualError(t, err, "the state database [height=9] is ahead of the block store [height=1]. "+
		"This is possible when the state database is not dropped after a ledger reset/rollback. "+
		"The state database can safely be dropped and will be rebuilt up to block store height upon the next peer start.")
}
