/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lockbasedtxmgr

import (
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/pvtstatepurgemgmt"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/queryutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/valimpl"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
)

var logger = flogging.MustGetLogger("lockbasedtxmgr")



type LockBasedTxMgr struct {
	ledgerid        string
	db              privacyenabledstate.DB
	pvtdataPurgeMgr *pvtdataPurgeMgr
	validator       validator.Validator
	stateListeners  []ledger.StateListener
	ccInfoProvider  ledger.DeployedChaincodeInfoProvider
	commitRWLock    sync.RWMutex
	current         *current
}

type current struct {
	block     *common.Block
	batch     *privacyenabledstate.UpdateBatch
	listeners []ledger.StateListener
}

func (c *current) blockNum() uint64 {
	return c.block.Header.Number
}

func (c *current) maxTxNumber() uint64 {
	return uint64(len(c.block.Data.Data)) - 1
}


func NewLockBasedTxMgr(ledgerid string, db privacyenabledstate.DB, stateListeners []ledger.StateListener,
	btlPolicy pvtdatapolicy.BTLPolicy, bookkeepingProvider bookkeeping.Provider, ccInfoProvider ledger.DeployedChaincodeInfoProvider) (*LockBasedTxMgr, error) {
	db.Open()
	txmgr := &LockBasedTxMgr{ledgerid: ledgerid, db: db, stateListeners: stateListeners, ccInfoProvider: ccInfoProvider}
	pvtstatePurgeMgr, err := pvtstatepurgemgmt.InstantiatePurgeMgr(ledgerid, db, btlPolicy, bookkeepingProvider)
	if err != nil {
		return nil, err
	}
	txmgr.pvtdataPurgeMgr = &pvtdataPurgeMgr{pvtstatePurgeMgr, false}
	txmgr.validator = valimpl.NewStatebasedValidator(txmgr, db)
	return txmgr, nil
}



func (txmgr *LockBasedTxMgr) GetLastSavepoint() (*version.Height, error) {
	return txmgr.db.GetLatestSavePoint()
}


func (txmgr *LockBasedTxMgr) NewQueryExecutor(txid string) (ledger.QueryExecutor, error) {
	qe := newQueryExecutor(txmgr, txid)
	txmgr.commitRWLock.RLock()
	return qe, nil
}


func (txmgr *LockBasedTxMgr) NewTxSimulator(txid string) (ledger.TxSimulator, error) {
	logger.Debugf("constructing new tx simulator")
	s, err := newLockBasedTxSimulator(txmgr, txid)
	if err != nil {
		return nil, err
	}
	txmgr.commitRWLock.RLock()
	return s, nil
}


func (txmgr *LockBasedTxMgr) ValidateAndPrepare(blockAndPvtdata *ledger.BlockAndPvtData, doMVCCValidation bool) error {
	block := blockAndPvtdata.Block
	logger.Debugf("Waiting for purge mgr to finish the background job of computing expirying keys for the block")
	txmgr.pvtdataPurgeMgr.WaitForPrepareToFinish()
	logger.Debugf("Validating new block with num trans = [%d]", len(block.Data.Data))
	batch, err := txmgr.validator.ValidateAndPrepareBatch(blockAndPvtdata, doMVCCValidation)
	if err != nil {
		txmgr.reset()
		return err
	}
	txmgr.current = &current{block: block, batch: batch}
	if err := txmgr.invokeNamespaceListeners(); err != nil {
		txmgr.reset()
		return err
	}
	return nil
}

func (txmgr *LockBasedTxMgr) invokeNamespaceListeners() error {
	for _, listener := range txmgr.stateListeners {
		stateUpdatesForListener := extractStateUpdates(txmgr.current.batch, listener.InterestedInNamespaces())
		if len(stateUpdatesForListener) == 0 {
			continue
		}
		txmgr.current.listeners = append(txmgr.current.listeners, listener)

		committedStateQueryExecuter := &queryutil.QECombiner{
			QueryExecuters: []queryutil.QueryExecuter{txmgr.db}}

		postCommitQueryExecuter := &queryutil.QECombiner{
			QueryExecuters: []queryutil.QueryExecuter{
				&queryutil.UpdateBatchBackedQueryExecuter{UpdateBatch: txmgr.current.batch.PubUpdates.UpdateBatch},
				txmgr.db,
			},
		}

		trigger := &ledger.StateUpdateTrigger{
			LedgerID:                    txmgr.ledgerid,
			StateUpdates:                stateUpdatesForListener,
			CommittingBlockNum:          txmgr.current.blockNum(),
			CommittedStateQueryExecutor: committedStateQueryExecuter,
			PostCommitQueryExecutor:     postCommitQueryExecuter,
		}
		if err := listener.HandleStateUpdates(trigger); err != nil {
			return err
		}
		logger.Debugf("Invoking listener for state changes:%s", listener)
	}
	return nil
}


func (txmgr *LockBasedTxMgr) Shutdown() {
	
	
	txmgr.pvtdataPurgeMgr.WaitForPrepareToFinish()
	txmgr.db.Close()
}


func (txmgr *LockBasedTxMgr) Commit() error {
	
	
	
	if !txmgr.pvtdataPurgeMgr.usedOnce {
		txmgr.pvtdataPurgeMgr.PrepareForExpiringKeys(txmgr.current.blockNum())
		txmgr.pvtdataPurgeMgr.usedOnce = true
	}
	defer func() {
		txmgr.clearCache()
		txmgr.pvtdataPurgeMgr.PrepareForExpiringKeys(txmgr.current.blockNum() + 1)
		logger.Debugf("Cleared version cache and launched the background routine for preparing keys to purge with the next block")
		txmgr.reset()
	}()

	logger.Debugf("Committing updates to state database")
	if txmgr.current == nil {
		panic("validateAndPrepare() method should have been called before calling commit()")
	}

	if err := txmgr.pvtdataPurgeMgr.DeleteExpiredAndUpdateBookkeeping(
		txmgr.current.batch.PvtUpdates, txmgr.current.batch.HashUpdates); err != nil {
		return err
	}

	txmgr.commitRWLock.Lock()
	defer txmgr.commitRWLock.Unlock()
	logger.Debugf("Write lock acquired for committing updates to state database")
	commitHeight := version.NewHeight(txmgr.current.blockNum(), txmgr.current.maxTxNumber())
	if err := txmgr.db.ApplyPrivacyAwareUpdates(txmgr.current.batch, commitHeight); err != nil {
		return err
	}
	logger.Debugf("Updates committed to state database")

	
	if err := txmgr.pvtdataPurgeMgr.BlockCommitDone(); err != nil {
		return err
	}
	
	
	txmgr.updateStateListeners()
	return nil
}


func (txmgr *LockBasedTxMgr) Rollback() {
	txmgr.reset()
}


func (txmgr *LockBasedTxMgr) clearCache() {
	if txmgr.db.IsBulkOptimizable() {
		txmgr.db.ClearCachedVersions()
	}
}


func (txmgr *LockBasedTxMgr) ShouldRecover(lastAvailableBlock uint64) (bool, uint64, error) {
	savepoint, err := txmgr.GetLastSavepoint()
	if err != nil {
		return false, 0, err
	}
	if savepoint == nil {
		return true, 0, nil
	}
	return savepoint.BlockNum != lastAvailableBlock, savepoint.BlockNum + 1, nil
}


func (txmgr *LockBasedTxMgr) CommitLostBlock(blockAndPvtdata *ledger.BlockAndPvtData) error {
	block := blockAndPvtdata.Block
	logger.Debugf("Constructing updateSet for the block %d", block.Header.Number)
	if err := txmgr.ValidateAndPrepare(blockAndPvtdata, false); err != nil {
		return err
	}

	
	if block.Header.Number%1000 == 0 {
		logger.Infof("Recommitting block [%d] to state database", block.Header.Number)
	} else {
		logger.Debugf("Recommitting block [%d] to state database", block.Header.Number)
	}

	return txmgr.Commit()
}

func extractStateUpdates(batch *privacyenabledstate.UpdateBatch, namespaces []string) ledger.StateUpdates {
	stateupdates := make(ledger.StateUpdates)
	for _, namespace := range namespaces {
		updatesMap := batch.PubUpdates.GetUpdates(namespace)
		var kvwrites []*kvrwset.KVWrite
		for key, versionedValue := range updatesMap {
			kvwrites = append(kvwrites, &kvrwset.KVWrite{Key: key, IsDelete: versionedValue.Value == nil, Value: versionedValue.Value})
			if len(kvwrites) > 0 {
				stateupdates[namespace] = kvwrites
			}
		}
	}
	return stateupdates
}

func (txmgr *LockBasedTxMgr) updateStateListeners() {
	for _, l := range txmgr.current.listeners {
		l.StateCommitDone(txmgr.ledgerid)
	}
}

func (txmgr *LockBasedTxMgr) reset() {
	txmgr.current = nil
}



type pvtdataPurgeMgr struct {
	pvtstatepurgemgmt.PurgeMgr
	usedOnce bool
}
