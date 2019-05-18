/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"sync"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/ledger/confighistory"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/history/historydb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr/lockbasedtxmgr"
	"github.com/mcc-github/blockchain/core/ledger/ledgerstorage"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("kvledger")



type kvLedger struct {
	ledgerID               string
	blockStore             *ledgerstorage.Store
	txtmgmt                txmgr.TxMgr
	historyDB              historydb.HistoryDB
	configHistoryRetriever ledger.ConfigHistoryRetriever
	blockAPIsRWLock        *sync.RWMutex
	stats                  *ledgerStats
}


func newKVLedger(
	ledgerID string,
	blockStore *ledgerstorage.Store,
	versionedDB privacyenabledstate.DB,
	historyDB historydb.HistoryDB,
	configHistoryMgr confighistory.Mgr,
	stateListeners []ledger.StateListener,
	bookkeeperProvider bookkeeping.Provider,
	ccInfoProvider ledger.DeployedChaincodeInfoProvider,
	ccLifecycleEventProvider ledger.ChaincodeLifecycleEventProvider,
	stats *ledgerStats,
) (*kvLedger, error) {
	logger.Debugf("Creating KVLedger ledgerID=%s: ", ledgerID)
	
	
	l := &kvLedger{ledgerID: ledgerID, blockStore: blockStore, historyDB: historyDB, blockAPIsRWLock: &sync.RWMutex{}}

	btlPolicy := pvtdatapolicy.ConstructBTLPolicy(&collectionInfoRetriever{ledgerID, l, ccInfoProvider})
	if err := l.initTxMgr(versionedDB, stateListeners, btlPolicy, bookkeeperProvider, ccInfoProvider); err != nil {
		return nil, err
	}

	l.initBlockStore(btlPolicy)

	
	
	
	ccEventListener := versionedDB.GetChaincodeEventListener()
	logger.Debugf("Register state db for chaincode lifecycle events: %t", ccEventListener != nil)
	if ccEventListener != nil {
		cceventmgmt.GetMgr().Register(ledgerID, ccEventListener)
		ccLifecycleEventProvider.RegisterListener(l.ledgerID, &ccEventListenerAdaptor{ccEventListener})
	}

	
	if err := l.recoverDBs(); err != nil {
		panic(errors.WithMessage(err, "error during state DB recovery"))
	}
	l.configHistoryRetriever = configHistoryMgr.GetRetriever(ledgerID, l)

	info, err := l.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}
	
	stats.updateBlockchainHeight(info.Height)
	l.stats = stats
	return l, nil
}

func (l *kvLedger) initTxMgr(versionedDB privacyenabledstate.DB, stateListeners []ledger.StateListener,
	btlPolicy pvtdatapolicy.BTLPolicy, bookkeeperProvider bookkeeping.Provider, ccInfoProvider ledger.DeployedChaincodeInfoProvider) error {
	var err error
	txmgr, err := lockbasedtxmgr.NewLockBasedTxMgr(l.ledgerID, versionedDB, stateListeners, btlPolicy, bookkeeperProvider, ccInfoProvider)
	if err != nil {
		return err
	}
	l.txtmgmt = txmgr
	
	
	qe, err := txmgr.NewQueryExecutorNoCollChecks()
	if err != nil {
		return err
	}
	defer qe.Done()
	for _, sl := range stateListeners {
		if err := sl.Initialize(l.ledgerID, qe); err != nil {
			return err
		}
	}
	return err
}

func (l *kvLedger) initBlockStore(btlPolicy pvtdatapolicy.BTLPolicy) {
	l.blockStore.Init(btlPolicy)
}



func (l *kvLedger) recoverDBs() error {
	logger.Debugf("Entering recoverDB()")
	if err := l.syncStateAndHistoryDBWithBlockstore(); err != nil {
		return err
	}
	if err := l.syncStateDBWithPvtdatastore(); err != nil {
		return err
	}
	return nil
}

func (l *kvLedger) syncStateAndHistoryDBWithBlockstore() error {
	
	info, _ := l.blockStore.GetBlockchainInfo()
	if info.Height == 0 {
		logger.Debug("Block storage is empty.")
		return nil
	}
	lastAvailableBlockNum := info.Height - 1
	recoverables := []recoverable{l.txtmgmt}
	if l.historyDB != nil {
		recoverables = append(recoverables, l.historyDB)
	}
	recoverers := []*recoverer{}
	for _, recoverable := range recoverables {
		recoverFlag, firstBlockNum, err := recoverable.ShouldRecover(lastAvailableBlockNum)
		if err != nil {
			return err
		}
		if recoverFlag {
			recoverers = append(recoverers, &recoverer{firstBlockNum, recoverable})
		}
	}
	if len(recoverers) == 0 {
		return nil
	}
	if len(recoverers) == 1 {
		return l.recommitLostBlocks(recoverers[0].firstBlockNum, lastAvailableBlockNum, recoverers[0].recoverable)
	}

	
	if recoverers[0].firstBlockNum > recoverers[1].firstBlockNum {
		
		recoverers[0], recoverers[1] = recoverers[1], recoverers[0]
	}
	if recoverers[0].firstBlockNum != recoverers[1].firstBlockNum {
		
		if err := l.recommitLostBlocks(recoverers[0].firstBlockNum, recoverers[1].firstBlockNum-1,
			recoverers[0].recoverable); err != nil {
			return err
		}
	}
	
	return l.recommitLostBlocks(recoverers[1].firstBlockNum, lastAvailableBlockNum,
		recoverers[0].recoverable, recoverers[1].recoverable)
}

func (l *kvLedger) syncStateDBWithPvtdatastore() error {
	
	
	
	
	
	
	
	
	blocksPvtData, err := l.blockStore.GetLastUpdatedOldBlocksPvtData()
	if err != nil {
		return err
	}
	if err := l.txtmgmt.RemoveStaleAndCommitPvtDataOfOldBlocks(blocksPvtData); err != nil {
		return err
	}
	if err := l.blockStore.ResetLastUpdatedOldBlocksList(); err != nil {
		return err
	}

	return nil
}



func (l *kvLedger) recommitLostBlocks(firstBlockNum uint64, lastBlockNum uint64, recoverables ...recoverable) error {
	logger.Infof("Recommitting lost blocks - firstBlockNum=%d, lastBlockNum=%d, recoverables=%#v", firstBlockNum, lastBlockNum, recoverables)
	var err error
	var blockAndPvtdata *ledger.BlockAndPvtData
	for blockNumber := firstBlockNum; blockNumber <= lastBlockNum; blockNumber++ {
		if blockAndPvtdata, err = l.GetPvtDataAndBlockByNum(blockNumber, nil); err != nil {
			return err
		}
		for _, r := range recoverables {
			if err := r.CommitLostBlock(blockAndPvtdata); err != nil {
				return err
			}
		}
	}
	logger.Infof("Recommitted lost blocks - firstBlockNum=%d, lastBlockNum=%d, recoverables=%#v", firstBlockNum, lastBlockNum, recoverables)
	return nil
}


func (l *kvLedger) GetTransactionByID(txID string) (*peer.ProcessedTransaction, error) {
	tranEnv, err := l.blockStore.RetrieveTxByID(txID)
	if err != nil {
		return nil, err
	}
	txVResult, err := l.blockStore.RetrieveTxValidationCodeByTxID(txID)
	if err != nil {
		return nil, err
	}
	processedTran := &peer.ProcessedTransaction{TransactionEnvelope: tranEnv, ValidationCode: int32(txVResult)}
	l.blockAPIsRWLock.RLock()
	l.blockAPIsRWLock.RUnlock()
	return processedTran, nil
}


func (l *kvLedger) GetBlockchainInfo() (*common.BlockchainInfo, error) {
	bcInfo, err := l.blockStore.GetBlockchainInfo()
	l.blockAPIsRWLock.RLock()
	defer l.blockAPIsRWLock.RUnlock()
	return bcInfo, err
}



func (l *kvLedger) GetBlockByNumber(blockNumber uint64) (*common.Block, error) {
	block, err := l.blockStore.RetrieveBlockByNumber(blockNumber)
	l.blockAPIsRWLock.RLock()
	l.blockAPIsRWLock.RUnlock()
	return block, err
}




func (l *kvLedger) GetBlocksIterator(startBlockNumber uint64) (commonledger.ResultsIterator, error) {
	blkItr, err := l.blockStore.RetrieveBlocks(startBlockNumber)
	if err != nil {
		return nil, err
	}
	return &blocksItr{l.blockAPIsRWLock, blkItr}, nil
}


func (l *kvLedger) GetBlockByHash(blockHash []byte) (*common.Block, error) {
	block, err := l.blockStore.RetrieveBlockByHash(blockHash)
	l.blockAPIsRWLock.RLock()
	l.blockAPIsRWLock.RUnlock()
	return block, err
}


func (l *kvLedger) GetBlockByTxID(txID string) (*common.Block, error) {
	block, err := l.blockStore.RetrieveBlockByTxID(txID)
	l.blockAPIsRWLock.RLock()
	l.blockAPIsRWLock.RUnlock()
	return block, err
}

func (l *kvLedger) GetTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	txValidationCode, err := l.blockStore.RetrieveTxValidationCodeByTxID(txID)
	l.blockAPIsRWLock.RLock()
	l.blockAPIsRWLock.RUnlock()
	return txValidationCode, err
}


func (l *kvLedger) NewTxSimulator(txid string) (ledger.TxSimulator, error) {
	return l.txtmgmt.NewTxSimulator(txid)
}




func (l *kvLedger) NewQueryExecutor() (ledger.QueryExecutor, error) {
	return l.txtmgmt.NewQueryExecutor(util.GenerateUUID())
}





func (l *kvLedger) NewHistoryQueryExecutor() (ledger.HistoryQueryExecutor, error) {
	if l.historyDB != nil {
		return l.historyDB.NewHistoryQueryExecutor(l.blockStore)
	}
	return nil, nil
}


func (l *kvLedger) CommitWithPvtData(pvtdataAndBlock *ledger.BlockAndPvtData) error {
	var err error
	block := pvtdataAndBlock.Block
	blockNo := pvtdataAndBlock.Block.Header.Number

	startBlockProcessing := time.Now()
	logger.Debugf("[%s] Validating state for block [%d]", l.ledgerID, blockNo)
	txstatsInfo, err := l.txtmgmt.ValidateAndPrepare(pvtdataAndBlock, true)
	if err != nil {
		return err
	}
	elapsedBlockProcessing := time.Since(startBlockProcessing)

	startCommitBlockStorage := time.Now()
	logger.Debugf("[%s] Committing block [%d] to storage", l.ledgerID, blockNo)
	l.blockAPIsRWLock.Lock()
	defer l.blockAPIsRWLock.Unlock()
	if err = l.blockStore.CommitWithPvtData(pvtdataAndBlock); err != nil {
		return err
	}
	elapsedCommitBlockStorage := time.Since(startCommitBlockStorage)

	startCommitState := time.Now()
	logger.Debugf("[%s] Committing block [%d] transactions to state database", l.ledgerID, blockNo)
	if err = l.txtmgmt.Commit(); err != nil {
		panic(errors.WithMessage(err, "error during commit to txmgr"))
	}
	elapsedCommitState := time.Since(startCommitState)

	
	
	if l.historyDB != nil {
		logger.Debugf("[%s] Committing block [%d] transactions to history database", l.ledgerID, blockNo)
		if err := l.historyDB.Commit(block); err != nil {
			panic(errors.WithMessage(err, "Error during commit to history db"))
		}
	}

	elapsedCommitWithPvtData := time.Since(startBlockProcessing)

	logger.Infof("[%s] Committed block [%d] with %d transaction(s) in %dms (state_validation=%dms block_commit=%dms state_commit=%dms)",
		l.ledgerID, block.Header.Number, len(block.Data.Data),
		elapsedCommitWithPvtData/time.Millisecond,
		elapsedBlockProcessing/time.Millisecond,
		elapsedCommitBlockStorage/time.Millisecond,
		elapsedCommitState/time.Millisecond,
	)
	l.updateBlockStats(blockNo,
		elapsedBlockProcessing,
		elapsedCommitBlockStorage,
		elapsedCommitState,
		txstatsInfo,
	)
	return nil
}

func (l *kvLedger) updateBlockStats(
	blockNum uint64,
	blockProcessingTime time.Duration,
	blockstorageCommitTime time.Duration,
	statedbCommitTime time.Duration,
	txstatsInfo []*txmgr.TxStatInfo,
) {
	l.stats.updateBlockchainHeight(blockNum + 1)
	l.stats.updateBlockProcessingTime(blockProcessingTime)
	l.stats.updateBlockstorageCommitTime(blockstorageCommitTime)
	l.stats.updateStatedbCommitTime(statedbCommitTime)
	l.stats.updateTransactionsStats(txstatsInfo)
}



func (l *kvLedger) GetMissingPvtDataInfoForMostRecentBlocks(maxBlock int) (ledger.MissingPvtDataInfo, error) {
	return l.blockStore.GetMissingPvtDataInfoForMostRecentBlocks(maxBlock)
}



func (l *kvLedger) GetPvtDataAndBlockByNum(blockNum uint64, filter ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error) {
	blockAndPvtdata, err := l.blockStore.GetPvtDataAndBlockByNum(blockNum, filter)
	l.blockAPIsRWLock.RLock()
	l.blockAPIsRWLock.RUnlock()
	return blockAndPvtdata, err
}



func (l *kvLedger) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	pvtdata, err := l.blockStore.GetPvtDataByNum(blockNum, filter)
	l.blockAPIsRWLock.RLock()
	l.blockAPIsRWLock.RUnlock()
	return pvtdata, err
}

func (l *kvLedger) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	return l.configHistoryRetriever, nil
}

func (l *kvLedger) CommitPvtDataOfOldBlocks(pvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error) {
	logger.Debugf("[%s:] Comparing pvtData of [%d] old blocks against the hashes in transaction's rwset to find valid and invalid data",
		l.ledgerID, len(pvtData))
	validPvtData, hashMismatches, err := ConstructValidAndInvalidPvtData(pvtData, l.blockStore)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%s:] Committing pvtData of [%d] old blocks to the pvtdatastore", l.ledgerID, len(pvtData))
	err = l.blockStore.CommitPvtDataOfOldBlocks(validPvtData)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%s:] Committing pvtData of [%d] old blocks to the stateDB", l.ledgerID, len(pvtData))
	err = l.txtmgmt.RemoveStaleAndCommitPvtDataOfOldBlocks(validPvtData)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%s:] Clearing the bookkeeping information from pvtdatastore", l.ledgerID)
	if err := l.blockStore.ResetLastUpdatedOldBlocksList(); err != nil {
		return nil, err
	}

	return hashMismatches, nil
}

func (l *kvLedger) GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error) {
	return l, nil
}


func (l *kvLedger) Close() {
	l.blockStore.Shutdown()
	l.txtmgmt.Shutdown()
}

type blocksItr struct {
	blockAPIsRWLock *sync.RWMutex
	blocksItr       commonledger.ResultsIterator
}

func (itr *blocksItr) Next() (commonledger.QueryResult, error) {
	block, err := itr.blocksItr.Next()
	if err != nil {
		return nil, err
	}
	itr.blockAPIsRWLock.RLock()
	itr.blockAPIsRWLock.RUnlock()
	return block, nil
}

func (itr *blocksItr) Close() {
	itr.blocksItr.Close()
}

type collectionInfoRetriever struct {
	ledgerID     string
	ledger       ledger.PeerLedger
	infoProvider ledger.DeployedChaincodeInfoProvider
}

func (r *collectionInfoRetriever) CollectionInfo(chaincodeName, collectionName string) (*common.StaticCollectionConfig, error) {
	qe, err := r.ledger.NewQueryExecutor()
	if err != nil {
		return nil, err
	}
	defer qe.Done()
	return r.infoProvider.CollectionInfo(r.ledgerID, chaincodeName, collectionName, qe)
}

type ccEventListenerAdaptor struct {
	legacyEventListener cceventmgmt.ChaincodeLifecycleEventListener
}

func (a *ccEventListenerAdaptor) HandleChaincodeDeploy(chaincodeDefinition *ledger.ChaincodeDefinition, dbArtifactsTar []byte) error {
	return a.legacyEventListener.HandleChaincodeDeploy(&cceventmgmt.ChaincodeDefinition{
		Name:              chaincodeDefinition.Name,
		Hash:              chaincodeDefinition.Hash,
		Version:           chaincodeDefinition.Version,
		CollectionConfigs: chaincodeDefinition.CollectionConfigs,
	},
		dbArtifactsTar,
	)
}

func (a *ccEventListenerAdaptor) ChaincodeDeployDone(succeeded bool) {
	a.legacyEventListener.ChaincodeDeployDone(succeeded)
}
