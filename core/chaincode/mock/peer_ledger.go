
package mock

import (
	"sync"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)

type PeerLedger struct {
	GetBlockchainInfoStub        func() (*common.BlockchainInfo, error)
	getBlockchainInfoMutex       sync.RWMutex
	getBlockchainInfoArgsForCall []struct{}
	getBlockchainInfoReturns     struct {
		result1 *common.BlockchainInfo
		result2 error
	}
	getBlockchainInfoReturnsOnCall map[int]struct {
		result1 *common.BlockchainInfo
		result2 error
	}
	GetBlockByNumberStub        func(blockNumber uint64) (*common.Block, error)
	getBlockByNumberMutex       sync.RWMutex
	getBlockByNumberArgsForCall []struct {
		blockNumber uint64
	}
	getBlockByNumberReturns struct {
		result1 *common.Block
		result2 error
	}
	getBlockByNumberReturnsOnCall map[int]struct {
		result1 *common.Block
		result2 error
	}
	GetBlocksIteratorStub        func(startBlockNumber uint64) (commonledger.ResultsIterator, error)
	getBlocksIteratorMutex       sync.RWMutex
	getBlocksIteratorArgsForCall []struct {
		startBlockNumber uint64
	}
	getBlocksIteratorReturns struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	getBlocksIteratorReturnsOnCall map[int]struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	CloseStub                     func()
	closeMutex                    sync.RWMutex
	closeArgsForCall              []struct{}
	GetTransactionByIDStub        func(txID string) (*peer.ProcessedTransaction, error)
	getTransactionByIDMutex       sync.RWMutex
	getTransactionByIDArgsForCall []struct {
		txID string
	}
	getTransactionByIDReturns struct {
		result1 *peer.ProcessedTransaction
		result2 error
	}
	getTransactionByIDReturnsOnCall map[int]struct {
		result1 *peer.ProcessedTransaction
		result2 error
	}
	GetBlockByHashStub        func(blockHash []byte) (*common.Block, error)
	getBlockByHashMutex       sync.RWMutex
	getBlockByHashArgsForCall []struct {
		blockHash []byte
	}
	getBlockByHashReturns struct {
		result1 *common.Block
		result2 error
	}
	getBlockByHashReturnsOnCall map[int]struct {
		result1 *common.Block
		result2 error
	}
	GetBlockByTxIDStub        func(txID string) (*common.Block, error)
	getBlockByTxIDMutex       sync.RWMutex
	getBlockByTxIDArgsForCall []struct {
		txID string
	}
	getBlockByTxIDReturns struct {
		result1 *common.Block
		result2 error
	}
	getBlockByTxIDReturnsOnCall map[int]struct {
		result1 *common.Block
		result2 error
	}
	GetTxValidationCodeByTxIDStub        func(txID string) (peer.TxValidationCode, error)
	getTxValidationCodeByTxIDMutex       sync.RWMutex
	getTxValidationCodeByTxIDArgsForCall []struct {
		txID string
	}
	getTxValidationCodeByTxIDReturns struct {
		result1 peer.TxValidationCode
		result2 error
	}
	getTxValidationCodeByTxIDReturnsOnCall map[int]struct {
		result1 peer.TxValidationCode
		result2 error
	}
	NewTxSimulatorStub        func(txid string) (ledger.TxSimulator, error)
	newTxSimulatorMutex       sync.RWMutex
	newTxSimulatorArgsForCall []struct {
		txid string
	}
	newTxSimulatorReturns struct {
		result1 ledger.TxSimulator
		result2 error
	}
	newTxSimulatorReturnsOnCall map[int]struct {
		result1 ledger.TxSimulator
		result2 error
	}
	NewQueryExecutorStub        func() (ledger.QueryExecutor, error)
	newQueryExecutorMutex       sync.RWMutex
	newQueryExecutorArgsForCall []struct{}
	newQueryExecutorReturns     struct {
		result1 ledger.QueryExecutor
		result2 error
	}
	newQueryExecutorReturnsOnCall map[int]struct {
		result1 ledger.QueryExecutor
		result2 error
	}
	NewHistoryQueryExecutorStub        func() (ledger.HistoryQueryExecutor, error)
	newHistoryQueryExecutorMutex       sync.RWMutex
	newHistoryQueryExecutorArgsForCall []struct{}
	newHistoryQueryExecutorReturns     struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}
	newHistoryQueryExecutorReturnsOnCall map[int]struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}
	GetPvtDataAndBlockByNumStub        func(blockNum uint64, filter ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error)
	getPvtDataAndBlockByNumMutex       sync.RWMutex
	getPvtDataAndBlockByNumArgsForCall []struct {
		blockNum uint64
		filter   ledger.PvtNsCollFilter
	}
	getPvtDataAndBlockByNumReturns struct {
		result1 *ledger.BlockAndPvtData
		result2 error
	}
	getPvtDataAndBlockByNumReturnsOnCall map[int]struct {
		result1 *ledger.BlockAndPvtData
		result2 error
	}
	GetPvtDataByNumStub        func(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)
	getPvtDataByNumMutex       sync.RWMutex
	getPvtDataByNumArgsForCall []struct {
		blockNum uint64
		filter   ledger.PvtNsCollFilter
	}
	getPvtDataByNumReturns struct {
		result1 []*ledger.TxPvtData
		result2 error
	}
	getPvtDataByNumReturnsOnCall map[int]struct {
		result1 []*ledger.TxPvtData
		result2 error
	}
	CommitWithPvtDataStub        func(blockAndPvtdata *ledger.BlockAndPvtData) error
	commitWithPvtDataMutex       sync.RWMutex
	commitWithPvtDataArgsForCall []struct {
		blockAndPvtdata *ledger.BlockAndPvtData
	}
	commitWithPvtDataReturns struct {
		result1 error
	}
	commitWithPvtDataReturnsOnCall map[int]struct {
		result1 error
	}
	PurgePrivateDataStub        func(maxBlockNumToRetain uint64) error
	purgePrivateDataMutex       sync.RWMutex
	purgePrivateDataArgsForCall []struct {
		maxBlockNumToRetain uint64
	}
	purgePrivateDataReturns struct {
		result1 error
	}
	purgePrivateDataReturnsOnCall map[int]struct {
		result1 error
	}
	PrivateDataMinBlockNumStub        func() (uint64, error)
	privateDataMinBlockNumMutex       sync.RWMutex
	privateDataMinBlockNumArgsForCall []struct{}
	privateDataMinBlockNumReturns     struct {
		result1 uint64
		result2 error
	}
	privateDataMinBlockNumReturnsOnCall map[int]struct {
		result1 uint64
		result2 error
	}
	PruneStub        func(policy commonledger.PrunePolicy) error
	pruneMutex       sync.RWMutex
	pruneArgsForCall []struct {
		policy commonledger.PrunePolicy
	}
	pruneReturns struct {
		result1 error
	}
	pruneReturnsOnCall map[int]struct {
		result1 error
	}
	GetConfigHistoryRetrieverStub        func() (ledger.ConfigHistoryRetriever, error)
	getConfigHistoryRetrieverMutex       sync.RWMutex
	getConfigHistoryRetrieverArgsForCall []struct{}
	getConfigHistoryRetrieverReturns     struct {
		result1 ledger.ConfigHistoryRetriever
		result2 error
	}
	getConfigHistoryRetrieverReturnsOnCall map[int]struct {
		result1 ledger.ConfigHistoryRetriever
		result2 error
	}
	CommitPvtDataStub        func(blockPvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error)
	commitPvtDataMutex       sync.RWMutex
	commitPvtDataArgsForCall []struct {
		blockPvtData []*ledger.BlockPvtData
	}
	commitPvtDataReturns struct {
		result1 []*ledger.PvtdataHashMismatch
		result2 error
	}
	commitPvtDataReturnsOnCall map[int]struct {
		result1 []*ledger.PvtdataHashMismatch
		result2 error
	}
	GetMissingPvtDataTrackerStub        func() (ledger.MissingPvtDataTracker, error)
	getMissingPvtDataTrackerMutex       sync.RWMutex
	getMissingPvtDataTrackerArgsForCall []struct{}
	getMissingPvtDataTrackerReturns     struct {
		result1 ledger.MissingPvtDataTracker
		result2 error
	}
	getMissingPvtDataTrackerReturnsOnCall map[int]struct {
		result1 ledger.MissingPvtDataTracker
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PeerLedger) GetBlockchainInfo() (*common.BlockchainInfo, error) {
	fake.getBlockchainInfoMutex.Lock()
	ret, specificReturn := fake.getBlockchainInfoReturnsOnCall[len(fake.getBlockchainInfoArgsForCall)]
	fake.getBlockchainInfoArgsForCall = append(fake.getBlockchainInfoArgsForCall, struct{}{})
	fake.recordInvocation("GetBlockchainInfo", []interface{}{})
	fake.getBlockchainInfoMutex.Unlock()
	if fake.GetBlockchainInfoStub != nil {
		return fake.GetBlockchainInfoStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getBlockchainInfoReturns.result1, fake.getBlockchainInfoReturns.result2
}

func (fake *PeerLedger) GetBlockchainInfoCallCount() int {
	fake.getBlockchainInfoMutex.RLock()
	defer fake.getBlockchainInfoMutex.RUnlock()
	return len(fake.getBlockchainInfoArgsForCall)
}

func (fake *PeerLedger) GetBlockchainInfoReturns(result1 *common.BlockchainInfo, result2 error) {
	fake.GetBlockchainInfoStub = nil
	fake.getBlockchainInfoReturns = struct {
		result1 *common.BlockchainInfo
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockchainInfoReturnsOnCall(i int, result1 *common.BlockchainInfo, result2 error) {
	fake.GetBlockchainInfoStub = nil
	if fake.getBlockchainInfoReturnsOnCall == nil {
		fake.getBlockchainInfoReturnsOnCall = make(map[int]struct {
			result1 *common.BlockchainInfo
			result2 error
		})
	}
	fake.getBlockchainInfoReturnsOnCall[i] = struct {
		result1 *common.BlockchainInfo
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByNumber(blockNumber uint64) (*common.Block, error) {
	fake.getBlockByNumberMutex.Lock()
	ret, specificReturn := fake.getBlockByNumberReturnsOnCall[len(fake.getBlockByNumberArgsForCall)]
	fake.getBlockByNumberArgsForCall = append(fake.getBlockByNumberArgsForCall, struct {
		blockNumber uint64
	}{blockNumber})
	fake.recordInvocation("GetBlockByNumber", []interface{}{blockNumber})
	fake.getBlockByNumberMutex.Unlock()
	if fake.GetBlockByNumberStub != nil {
		return fake.GetBlockByNumberStub(blockNumber)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getBlockByNumberReturns.result1, fake.getBlockByNumberReturns.result2
}

func (fake *PeerLedger) GetBlockByNumberCallCount() int {
	fake.getBlockByNumberMutex.RLock()
	defer fake.getBlockByNumberMutex.RUnlock()
	return len(fake.getBlockByNumberArgsForCall)
}

func (fake *PeerLedger) GetBlockByNumberArgsForCall(i int) uint64 {
	fake.getBlockByNumberMutex.RLock()
	defer fake.getBlockByNumberMutex.RUnlock()
	return fake.getBlockByNumberArgsForCall[i].blockNumber
}

func (fake *PeerLedger) GetBlockByNumberReturns(result1 *common.Block, result2 error) {
	fake.GetBlockByNumberStub = nil
	fake.getBlockByNumberReturns = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByNumberReturnsOnCall(i int, result1 *common.Block, result2 error) {
	fake.GetBlockByNumberStub = nil
	if fake.getBlockByNumberReturnsOnCall == nil {
		fake.getBlockByNumberReturnsOnCall = make(map[int]struct {
			result1 *common.Block
			result2 error
		})
	}
	fake.getBlockByNumberReturnsOnCall[i] = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlocksIterator(startBlockNumber uint64) (commonledger.ResultsIterator, error) {
	fake.getBlocksIteratorMutex.Lock()
	ret, specificReturn := fake.getBlocksIteratorReturnsOnCall[len(fake.getBlocksIteratorArgsForCall)]
	fake.getBlocksIteratorArgsForCall = append(fake.getBlocksIteratorArgsForCall, struct {
		startBlockNumber uint64
	}{startBlockNumber})
	fake.recordInvocation("GetBlocksIterator", []interface{}{startBlockNumber})
	fake.getBlocksIteratorMutex.Unlock()
	if fake.GetBlocksIteratorStub != nil {
		return fake.GetBlocksIteratorStub(startBlockNumber)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getBlocksIteratorReturns.result1, fake.getBlocksIteratorReturns.result2
}

func (fake *PeerLedger) GetBlocksIteratorCallCount() int {
	fake.getBlocksIteratorMutex.RLock()
	defer fake.getBlocksIteratorMutex.RUnlock()
	return len(fake.getBlocksIteratorArgsForCall)
}

func (fake *PeerLedger) GetBlocksIteratorArgsForCall(i int) uint64 {
	fake.getBlocksIteratorMutex.RLock()
	defer fake.getBlocksIteratorMutex.RUnlock()
	return fake.getBlocksIteratorArgsForCall[i].startBlockNumber
}

func (fake *PeerLedger) GetBlocksIteratorReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.GetBlocksIteratorStub = nil
	fake.getBlocksIteratorReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlocksIteratorReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
	fake.GetBlocksIteratorStub = nil
	if fake.getBlocksIteratorReturnsOnCall == nil {
		fake.getBlocksIteratorReturnsOnCall = make(map[int]struct {
			result1 commonledger.ResultsIterator
			result2 error
		})
	}
	fake.getBlocksIteratorReturnsOnCall[i] = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *PeerLedger) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *PeerLedger) GetTransactionByID(txID string) (*peer.ProcessedTransaction, error) {
	fake.getTransactionByIDMutex.Lock()
	ret, specificReturn := fake.getTransactionByIDReturnsOnCall[len(fake.getTransactionByIDArgsForCall)]
	fake.getTransactionByIDArgsForCall = append(fake.getTransactionByIDArgsForCall, struct {
		txID string
	}{txID})
	fake.recordInvocation("GetTransactionByID", []interface{}{txID})
	fake.getTransactionByIDMutex.Unlock()
	if fake.GetTransactionByIDStub != nil {
		return fake.GetTransactionByIDStub(txID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTransactionByIDReturns.result1, fake.getTransactionByIDReturns.result2
}

func (fake *PeerLedger) GetTransactionByIDCallCount() int {
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	return len(fake.getTransactionByIDArgsForCall)
}

func (fake *PeerLedger) GetTransactionByIDArgsForCall(i int) string {
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	return fake.getTransactionByIDArgsForCall[i].txID
}

func (fake *PeerLedger) GetTransactionByIDReturns(result1 *peer.ProcessedTransaction, result2 error) {
	fake.GetTransactionByIDStub = nil
	fake.getTransactionByIDReturns = struct {
		result1 *peer.ProcessedTransaction
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTransactionByIDReturnsOnCall(i int, result1 *peer.ProcessedTransaction, result2 error) {
	fake.GetTransactionByIDStub = nil
	if fake.getTransactionByIDReturnsOnCall == nil {
		fake.getTransactionByIDReturnsOnCall = make(map[int]struct {
			result1 *peer.ProcessedTransaction
			result2 error
		})
	}
	fake.getTransactionByIDReturnsOnCall[i] = struct {
		result1 *peer.ProcessedTransaction
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByHash(blockHash []byte) (*common.Block, error) {
	var blockHashCopy []byte
	if blockHash != nil {
		blockHashCopy = make([]byte, len(blockHash))
		copy(blockHashCopy, blockHash)
	}
	fake.getBlockByHashMutex.Lock()
	ret, specificReturn := fake.getBlockByHashReturnsOnCall[len(fake.getBlockByHashArgsForCall)]
	fake.getBlockByHashArgsForCall = append(fake.getBlockByHashArgsForCall, struct {
		blockHash []byte
	}{blockHashCopy})
	fake.recordInvocation("GetBlockByHash", []interface{}{blockHashCopy})
	fake.getBlockByHashMutex.Unlock()
	if fake.GetBlockByHashStub != nil {
		return fake.GetBlockByHashStub(blockHash)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getBlockByHashReturns.result1, fake.getBlockByHashReturns.result2
}

func (fake *PeerLedger) GetBlockByHashCallCount() int {
	fake.getBlockByHashMutex.RLock()
	defer fake.getBlockByHashMutex.RUnlock()
	return len(fake.getBlockByHashArgsForCall)
}

func (fake *PeerLedger) GetBlockByHashArgsForCall(i int) []byte {
	fake.getBlockByHashMutex.RLock()
	defer fake.getBlockByHashMutex.RUnlock()
	return fake.getBlockByHashArgsForCall[i].blockHash
}

func (fake *PeerLedger) GetBlockByHashReturns(result1 *common.Block, result2 error) {
	fake.GetBlockByHashStub = nil
	fake.getBlockByHashReturns = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByHashReturnsOnCall(i int, result1 *common.Block, result2 error) {
	fake.GetBlockByHashStub = nil
	if fake.getBlockByHashReturnsOnCall == nil {
		fake.getBlockByHashReturnsOnCall = make(map[int]struct {
			result1 *common.Block
			result2 error
		})
	}
	fake.getBlockByHashReturnsOnCall[i] = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByTxID(txID string) (*common.Block, error) {
	fake.getBlockByTxIDMutex.Lock()
	ret, specificReturn := fake.getBlockByTxIDReturnsOnCall[len(fake.getBlockByTxIDArgsForCall)]
	fake.getBlockByTxIDArgsForCall = append(fake.getBlockByTxIDArgsForCall, struct {
		txID string
	}{txID})
	fake.recordInvocation("GetBlockByTxID", []interface{}{txID})
	fake.getBlockByTxIDMutex.Unlock()
	if fake.GetBlockByTxIDStub != nil {
		return fake.GetBlockByTxIDStub(txID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getBlockByTxIDReturns.result1, fake.getBlockByTxIDReturns.result2
}

func (fake *PeerLedger) GetBlockByTxIDCallCount() int {
	fake.getBlockByTxIDMutex.RLock()
	defer fake.getBlockByTxIDMutex.RUnlock()
	return len(fake.getBlockByTxIDArgsForCall)
}

func (fake *PeerLedger) GetBlockByTxIDArgsForCall(i int) string {
	fake.getBlockByTxIDMutex.RLock()
	defer fake.getBlockByTxIDMutex.RUnlock()
	return fake.getBlockByTxIDArgsForCall[i].txID
}

func (fake *PeerLedger) GetBlockByTxIDReturns(result1 *common.Block, result2 error) {
	fake.GetBlockByTxIDStub = nil
	fake.getBlockByTxIDReturns = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByTxIDReturnsOnCall(i int, result1 *common.Block, result2 error) {
	fake.GetBlockByTxIDStub = nil
	if fake.getBlockByTxIDReturnsOnCall == nil {
		fake.getBlockByTxIDReturnsOnCall = make(map[int]struct {
			result1 *common.Block
			result2 error
		})
	}
	fake.getBlockByTxIDReturnsOnCall[i] = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	fake.getTxValidationCodeByTxIDMutex.Lock()
	ret, specificReturn := fake.getTxValidationCodeByTxIDReturnsOnCall[len(fake.getTxValidationCodeByTxIDArgsForCall)]
	fake.getTxValidationCodeByTxIDArgsForCall = append(fake.getTxValidationCodeByTxIDArgsForCall, struct {
		txID string
	}{txID})
	fake.recordInvocation("GetTxValidationCodeByTxID", []interface{}{txID})
	fake.getTxValidationCodeByTxIDMutex.Unlock()
	if fake.GetTxValidationCodeByTxIDStub != nil {
		return fake.GetTxValidationCodeByTxIDStub(txID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTxValidationCodeByTxIDReturns.result1, fake.getTxValidationCodeByTxIDReturns.result2
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDCallCount() int {
	fake.getTxValidationCodeByTxIDMutex.RLock()
	defer fake.getTxValidationCodeByTxIDMutex.RUnlock()
	return len(fake.getTxValidationCodeByTxIDArgsForCall)
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDArgsForCall(i int) string {
	fake.getTxValidationCodeByTxIDMutex.RLock()
	defer fake.getTxValidationCodeByTxIDMutex.RUnlock()
	return fake.getTxValidationCodeByTxIDArgsForCall[i].txID
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDReturns(result1 peer.TxValidationCode, result2 error) {
	fake.GetTxValidationCodeByTxIDStub = nil
	fake.getTxValidationCodeByTxIDReturns = struct {
		result1 peer.TxValidationCode
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDReturnsOnCall(i int, result1 peer.TxValidationCode, result2 error) {
	fake.GetTxValidationCodeByTxIDStub = nil
	if fake.getTxValidationCodeByTxIDReturnsOnCall == nil {
		fake.getTxValidationCodeByTxIDReturnsOnCall = make(map[int]struct {
			result1 peer.TxValidationCode
			result2 error
		})
	}
	fake.getTxValidationCodeByTxIDReturnsOnCall[i] = struct {
		result1 peer.TxValidationCode
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewTxSimulator(txid string) (ledger.TxSimulator, error) {
	fake.newTxSimulatorMutex.Lock()
	ret, specificReturn := fake.newTxSimulatorReturnsOnCall[len(fake.newTxSimulatorArgsForCall)]
	fake.newTxSimulatorArgsForCall = append(fake.newTxSimulatorArgsForCall, struct {
		txid string
	}{txid})
	fake.recordInvocation("NewTxSimulator", []interface{}{txid})
	fake.newTxSimulatorMutex.Unlock()
	if fake.NewTxSimulatorStub != nil {
		return fake.NewTxSimulatorStub(txid)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newTxSimulatorReturns.result1, fake.newTxSimulatorReturns.result2
}

func (fake *PeerLedger) NewTxSimulatorCallCount() int {
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	return len(fake.newTxSimulatorArgsForCall)
}

func (fake *PeerLedger) NewTxSimulatorArgsForCall(i int) string {
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	return fake.newTxSimulatorArgsForCall[i].txid
}

func (fake *PeerLedger) NewTxSimulatorReturns(result1 ledger.TxSimulator, result2 error) {
	fake.NewTxSimulatorStub = nil
	fake.newTxSimulatorReturns = struct {
		result1 ledger.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewTxSimulatorReturnsOnCall(i int, result1 ledger.TxSimulator, result2 error) {
	fake.NewTxSimulatorStub = nil
	if fake.newTxSimulatorReturnsOnCall == nil {
		fake.newTxSimulatorReturnsOnCall = make(map[int]struct {
			result1 ledger.TxSimulator
			result2 error
		})
	}
	fake.newTxSimulatorReturnsOnCall[i] = struct {
		result1 ledger.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewQueryExecutor() (ledger.QueryExecutor, error) {
	fake.newQueryExecutorMutex.Lock()
	ret, specificReturn := fake.newQueryExecutorReturnsOnCall[len(fake.newQueryExecutorArgsForCall)]
	fake.newQueryExecutorArgsForCall = append(fake.newQueryExecutorArgsForCall, struct{}{})
	fake.recordInvocation("NewQueryExecutor", []interface{}{})
	fake.newQueryExecutorMutex.Unlock()
	if fake.NewQueryExecutorStub != nil {
		return fake.NewQueryExecutorStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newQueryExecutorReturns.result1, fake.newQueryExecutorReturns.result2
}

func (fake *PeerLedger) NewQueryExecutorCallCount() int {
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	return len(fake.newQueryExecutorArgsForCall)
}

func (fake *PeerLedger) NewQueryExecutorReturns(result1 ledger.QueryExecutor, result2 error) {
	fake.NewQueryExecutorStub = nil
	fake.newQueryExecutorReturns = struct {
		result1 ledger.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewQueryExecutorReturnsOnCall(i int, result1 ledger.QueryExecutor, result2 error) {
	fake.NewQueryExecutorStub = nil
	if fake.newQueryExecutorReturnsOnCall == nil {
		fake.newQueryExecutorReturnsOnCall = make(map[int]struct {
			result1 ledger.QueryExecutor
			result2 error
		})
	}
	fake.newQueryExecutorReturnsOnCall[i] = struct {
		result1 ledger.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewHistoryQueryExecutor() (ledger.HistoryQueryExecutor, error) {
	fake.newHistoryQueryExecutorMutex.Lock()
	ret, specificReturn := fake.newHistoryQueryExecutorReturnsOnCall[len(fake.newHistoryQueryExecutorArgsForCall)]
	fake.newHistoryQueryExecutorArgsForCall = append(fake.newHistoryQueryExecutorArgsForCall, struct{}{})
	fake.recordInvocation("NewHistoryQueryExecutor", []interface{}{})
	fake.newHistoryQueryExecutorMutex.Unlock()
	if fake.NewHistoryQueryExecutorStub != nil {
		return fake.NewHistoryQueryExecutorStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newHistoryQueryExecutorReturns.result1, fake.newHistoryQueryExecutorReturns.result2
}

func (fake *PeerLedger) NewHistoryQueryExecutorCallCount() int {
	fake.newHistoryQueryExecutorMutex.RLock()
	defer fake.newHistoryQueryExecutorMutex.RUnlock()
	return len(fake.newHistoryQueryExecutorArgsForCall)
}

func (fake *PeerLedger) NewHistoryQueryExecutorReturns(result1 ledger.HistoryQueryExecutor, result2 error) {
	fake.NewHistoryQueryExecutorStub = nil
	fake.newHistoryQueryExecutorReturns = struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewHistoryQueryExecutorReturnsOnCall(i int, result1 ledger.HistoryQueryExecutor, result2 error) {
	fake.NewHistoryQueryExecutorStub = nil
	if fake.newHistoryQueryExecutorReturnsOnCall == nil {
		fake.newHistoryQueryExecutorReturnsOnCall = make(map[int]struct {
			result1 ledger.HistoryQueryExecutor
			result2 error
		})
	}
	fake.newHistoryQueryExecutorReturnsOnCall[i] = struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataAndBlockByNum(blockNum uint64, filter ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	ret, specificReturn := fake.getPvtDataAndBlockByNumReturnsOnCall[len(fake.getPvtDataAndBlockByNumArgsForCall)]
	fake.getPvtDataAndBlockByNumArgsForCall = append(fake.getPvtDataAndBlockByNumArgsForCall, struct {
		blockNum uint64
		filter   ledger.PvtNsCollFilter
	}{blockNum, filter})
	fake.recordInvocation("GetPvtDataAndBlockByNum", []interface{}{blockNum, filter})
	fake.getPvtDataAndBlockByNumMutex.Unlock()
	if fake.GetPvtDataAndBlockByNumStub != nil {
		return fake.GetPvtDataAndBlockByNumStub(blockNum, filter)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPvtDataAndBlockByNumReturns.result1, fake.getPvtDataAndBlockByNumReturns.result2
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumCallCount() int {
	fake.getPvtDataAndBlockByNumMutex.RLock()
	defer fake.getPvtDataAndBlockByNumMutex.RUnlock()
	return len(fake.getPvtDataAndBlockByNumArgsForCall)
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumArgsForCall(i int) (uint64, ledger.PvtNsCollFilter) {
	fake.getPvtDataAndBlockByNumMutex.RLock()
	defer fake.getPvtDataAndBlockByNumMutex.RUnlock()
	return fake.getPvtDataAndBlockByNumArgsForCall[i].blockNum, fake.getPvtDataAndBlockByNumArgsForCall[i].filter
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumReturns(result1 *ledger.BlockAndPvtData, result2 error) {
	fake.GetPvtDataAndBlockByNumStub = nil
	fake.getPvtDataAndBlockByNumReturns = struct {
		result1 *ledger.BlockAndPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumReturnsOnCall(i int, result1 *ledger.BlockAndPvtData, result2 error) {
	fake.GetPvtDataAndBlockByNumStub = nil
	if fake.getPvtDataAndBlockByNumReturnsOnCall == nil {
		fake.getPvtDataAndBlockByNumReturnsOnCall = make(map[int]struct {
			result1 *ledger.BlockAndPvtData
			result2 error
		})
	}
	fake.getPvtDataAndBlockByNumReturnsOnCall[i] = struct {
		result1 *ledger.BlockAndPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	fake.getPvtDataByNumMutex.Lock()
	ret, specificReturn := fake.getPvtDataByNumReturnsOnCall[len(fake.getPvtDataByNumArgsForCall)]
	fake.getPvtDataByNumArgsForCall = append(fake.getPvtDataByNumArgsForCall, struct {
		blockNum uint64
		filter   ledger.PvtNsCollFilter
	}{blockNum, filter})
	fake.recordInvocation("GetPvtDataByNum", []interface{}{blockNum, filter})
	fake.getPvtDataByNumMutex.Unlock()
	if fake.GetPvtDataByNumStub != nil {
		return fake.GetPvtDataByNumStub(blockNum, filter)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPvtDataByNumReturns.result1, fake.getPvtDataByNumReturns.result2
}

func (fake *PeerLedger) GetPvtDataByNumCallCount() int {
	fake.getPvtDataByNumMutex.RLock()
	defer fake.getPvtDataByNumMutex.RUnlock()
	return len(fake.getPvtDataByNumArgsForCall)
}

func (fake *PeerLedger) GetPvtDataByNumArgsForCall(i int) (uint64, ledger.PvtNsCollFilter) {
	fake.getPvtDataByNumMutex.RLock()
	defer fake.getPvtDataByNumMutex.RUnlock()
	return fake.getPvtDataByNumArgsForCall[i].blockNum, fake.getPvtDataByNumArgsForCall[i].filter
}

func (fake *PeerLedger) GetPvtDataByNumReturns(result1 []*ledger.TxPvtData, result2 error) {
	fake.GetPvtDataByNumStub = nil
	fake.getPvtDataByNumReturns = struct {
		result1 []*ledger.TxPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataByNumReturnsOnCall(i int, result1 []*ledger.TxPvtData, result2 error) {
	fake.GetPvtDataByNumStub = nil
	if fake.getPvtDataByNumReturnsOnCall == nil {
		fake.getPvtDataByNumReturnsOnCall = make(map[int]struct {
			result1 []*ledger.TxPvtData
			result2 error
		})
	}
	fake.getPvtDataByNumReturnsOnCall[i] = struct {
		result1 []*ledger.TxPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) CommitWithPvtData(blockAndPvtdata *ledger.BlockAndPvtData) error {
	fake.commitWithPvtDataMutex.Lock()
	ret, specificReturn := fake.commitWithPvtDataReturnsOnCall[len(fake.commitWithPvtDataArgsForCall)]
	fake.commitWithPvtDataArgsForCall = append(fake.commitWithPvtDataArgsForCall, struct {
		blockAndPvtdata *ledger.BlockAndPvtData
	}{blockAndPvtdata})
	fake.recordInvocation("CommitWithPvtData", []interface{}{blockAndPvtdata})
	fake.commitWithPvtDataMutex.Unlock()
	if fake.CommitWithPvtDataStub != nil {
		return fake.CommitWithPvtDataStub(blockAndPvtdata)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.commitWithPvtDataReturns.result1
}

func (fake *PeerLedger) CommitWithPvtDataCallCount() int {
	fake.commitWithPvtDataMutex.RLock()
	defer fake.commitWithPvtDataMutex.RUnlock()
	return len(fake.commitWithPvtDataArgsForCall)
}

func (fake *PeerLedger) CommitWithPvtDataArgsForCall(i int) *ledger.BlockAndPvtData {
	fake.commitWithPvtDataMutex.RLock()
	defer fake.commitWithPvtDataMutex.RUnlock()
	return fake.commitWithPvtDataArgsForCall[i].blockAndPvtdata
}

func (fake *PeerLedger) CommitWithPvtDataReturns(result1 error) {
	fake.CommitWithPvtDataStub = nil
	fake.commitWithPvtDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) CommitWithPvtDataReturnsOnCall(i int, result1 error) {
	fake.CommitWithPvtDataStub = nil
	if fake.commitWithPvtDataReturnsOnCall == nil {
		fake.commitWithPvtDataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.commitWithPvtDataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) PurgePrivateData(maxBlockNumToRetain uint64) error {
	fake.purgePrivateDataMutex.Lock()
	ret, specificReturn := fake.purgePrivateDataReturnsOnCall[len(fake.purgePrivateDataArgsForCall)]
	fake.purgePrivateDataArgsForCall = append(fake.purgePrivateDataArgsForCall, struct {
		maxBlockNumToRetain uint64
	}{maxBlockNumToRetain})
	fake.recordInvocation("PurgePrivateData", []interface{}{maxBlockNumToRetain})
	fake.purgePrivateDataMutex.Unlock()
	if fake.PurgePrivateDataStub != nil {
		return fake.PurgePrivateDataStub(maxBlockNumToRetain)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.purgePrivateDataReturns.result1
}

func (fake *PeerLedger) PurgePrivateDataCallCount() int {
	fake.purgePrivateDataMutex.RLock()
	defer fake.purgePrivateDataMutex.RUnlock()
	return len(fake.purgePrivateDataArgsForCall)
}

func (fake *PeerLedger) PurgePrivateDataArgsForCall(i int) uint64 {
	fake.purgePrivateDataMutex.RLock()
	defer fake.purgePrivateDataMutex.RUnlock()
	return fake.purgePrivateDataArgsForCall[i].maxBlockNumToRetain
}

func (fake *PeerLedger) PurgePrivateDataReturns(result1 error) {
	fake.PurgePrivateDataStub = nil
	fake.purgePrivateDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) PurgePrivateDataReturnsOnCall(i int, result1 error) {
	fake.PurgePrivateDataStub = nil
	if fake.purgePrivateDataReturnsOnCall == nil {
		fake.purgePrivateDataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.purgePrivateDataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) PrivateDataMinBlockNum() (uint64, error) {
	fake.privateDataMinBlockNumMutex.Lock()
	ret, specificReturn := fake.privateDataMinBlockNumReturnsOnCall[len(fake.privateDataMinBlockNumArgsForCall)]
	fake.privateDataMinBlockNumArgsForCall = append(fake.privateDataMinBlockNumArgsForCall, struct{}{})
	fake.recordInvocation("PrivateDataMinBlockNum", []interface{}{})
	fake.privateDataMinBlockNumMutex.Unlock()
	if fake.PrivateDataMinBlockNumStub != nil {
		return fake.PrivateDataMinBlockNumStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.privateDataMinBlockNumReturns.result1, fake.privateDataMinBlockNumReturns.result2
}

func (fake *PeerLedger) PrivateDataMinBlockNumCallCount() int {
	fake.privateDataMinBlockNumMutex.RLock()
	defer fake.privateDataMinBlockNumMutex.RUnlock()
	return len(fake.privateDataMinBlockNumArgsForCall)
}

func (fake *PeerLedger) PrivateDataMinBlockNumReturns(result1 uint64, result2 error) {
	fake.PrivateDataMinBlockNumStub = nil
	fake.privateDataMinBlockNumReturns = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) PrivateDataMinBlockNumReturnsOnCall(i int, result1 uint64, result2 error) {
	fake.PrivateDataMinBlockNumStub = nil
	if fake.privateDataMinBlockNumReturnsOnCall == nil {
		fake.privateDataMinBlockNumReturnsOnCall = make(map[int]struct {
			result1 uint64
			result2 error
		})
	}
	fake.privateDataMinBlockNumReturnsOnCall[i] = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) Prune(policy commonledger.PrunePolicy) error {
	fake.pruneMutex.Lock()
	ret, specificReturn := fake.pruneReturnsOnCall[len(fake.pruneArgsForCall)]
	fake.pruneArgsForCall = append(fake.pruneArgsForCall, struct {
		policy commonledger.PrunePolicy
	}{policy})
	fake.recordInvocation("Prune", []interface{}{policy})
	fake.pruneMutex.Unlock()
	if fake.PruneStub != nil {
		return fake.PruneStub(policy)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.pruneReturns.result1
}

func (fake *PeerLedger) PruneCallCount() int {
	fake.pruneMutex.RLock()
	defer fake.pruneMutex.RUnlock()
	return len(fake.pruneArgsForCall)
}

func (fake *PeerLedger) PruneArgsForCall(i int) commonledger.PrunePolicy {
	fake.pruneMutex.RLock()
	defer fake.pruneMutex.RUnlock()
	return fake.pruneArgsForCall[i].policy
}

func (fake *PeerLedger) PruneReturns(result1 error) {
	fake.PruneStub = nil
	fake.pruneReturns = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) PruneReturnsOnCall(i int, result1 error) {
	fake.PruneStub = nil
	if fake.pruneReturnsOnCall == nil {
		fake.pruneReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.pruneReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	fake.getConfigHistoryRetrieverMutex.Lock()
	ret, specificReturn := fake.getConfigHistoryRetrieverReturnsOnCall[len(fake.getConfigHistoryRetrieverArgsForCall)]
	fake.getConfigHistoryRetrieverArgsForCall = append(fake.getConfigHistoryRetrieverArgsForCall, struct{}{})
	fake.recordInvocation("GetConfigHistoryRetriever", []interface{}{})
	fake.getConfigHistoryRetrieverMutex.Unlock()
	if fake.GetConfigHistoryRetrieverStub != nil {
		return fake.GetConfigHistoryRetrieverStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getConfigHistoryRetrieverReturns.result1, fake.getConfigHistoryRetrieverReturns.result2
}

func (fake *PeerLedger) GetConfigHistoryRetrieverCallCount() int {
	fake.getConfigHistoryRetrieverMutex.RLock()
	defer fake.getConfigHistoryRetrieverMutex.RUnlock()
	return len(fake.getConfigHistoryRetrieverArgsForCall)
}

func (fake *PeerLedger) GetConfigHistoryRetrieverReturns(result1 ledger.ConfigHistoryRetriever, result2 error) {
	fake.GetConfigHistoryRetrieverStub = nil
	fake.getConfigHistoryRetrieverReturns = struct {
		result1 ledger.ConfigHistoryRetriever
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetConfigHistoryRetrieverReturnsOnCall(i int, result1 ledger.ConfigHistoryRetriever, result2 error) {
	fake.GetConfigHistoryRetrieverStub = nil
	if fake.getConfigHistoryRetrieverReturnsOnCall == nil {
		fake.getConfigHistoryRetrieverReturnsOnCall = make(map[int]struct {
			result1 ledger.ConfigHistoryRetriever
			result2 error
		})
	}
	fake.getConfigHistoryRetrieverReturnsOnCall[i] = struct {
		result1 ledger.ConfigHistoryRetriever
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) CommitPvtData(blockPvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error) {
	var blockPvtDataCopy []*ledger.BlockPvtData
	if blockPvtData != nil {
		blockPvtDataCopy = make([]*ledger.BlockPvtData, len(blockPvtData))
		copy(blockPvtDataCopy, blockPvtData)
	}
	fake.commitPvtDataMutex.Lock()
	ret, specificReturn := fake.commitPvtDataReturnsOnCall[len(fake.commitPvtDataArgsForCall)]
	fake.commitPvtDataArgsForCall = append(fake.commitPvtDataArgsForCall, struct {
		blockPvtData []*ledger.BlockPvtData
	}{blockPvtDataCopy})
	fake.recordInvocation("CommitPvtData", []interface{}{blockPvtDataCopy})
	fake.commitPvtDataMutex.Unlock()
	if fake.CommitPvtDataStub != nil {
		return fake.CommitPvtDataStub(blockPvtData)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.commitPvtDataReturns.result1, fake.commitPvtDataReturns.result2
}

func (fake *PeerLedger) CommitPvtDataCallCount() int {
	fake.commitPvtDataMutex.RLock()
	defer fake.commitPvtDataMutex.RUnlock()
	return len(fake.commitPvtDataArgsForCall)
}

func (fake *PeerLedger) CommitPvtDataArgsForCall(i int) []*ledger.BlockPvtData {
	fake.commitPvtDataMutex.RLock()
	defer fake.commitPvtDataMutex.RUnlock()
	return fake.commitPvtDataArgsForCall[i].blockPvtData
}

func (fake *PeerLedger) CommitPvtDataReturns(result1 []*ledger.PvtdataHashMismatch, result2 error) {
	fake.CommitPvtDataStub = nil
	fake.commitPvtDataReturns = struct {
		result1 []*ledger.PvtdataHashMismatch
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) CommitPvtDataReturnsOnCall(i int, result1 []*ledger.PvtdataHashMismatch, result2 error) {
	fake.CommitPvtDataStub = nil
	if fake.commitPvtDataReturnsOnCall == nil {
		fake.commitPvtDataReturnsOnCall = make(map[int]struct {
			result1 []*ledger.PvtdataHashMismatch
			result2 error
		})
	}
	fake.commitPvtDataReturnsOnCall[i] = struct {
		result1 []*ledger.PvtdataHashMismatch
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error) {
	fake.getMissingPvtDataTrackerMutex.Lock()
	ret, specificReturn := fake.getMissingPvtDataTrackerReturnsOnCall[len(fake.getMissingPvtDataTrackerArgsForCall)]
	fake.getMissingPvtDataTrackerArgsForCall = append(fake.getMissingPvtDataTrackerArgsForCall, struct{}{})
	fake.recordInvocation("GetMissingPvtDataTracker", []interface{}{})
	fake.getMissingPvtDataTrackerMutex.Unlock()
	if fake.GetMissingPvtDataTrackerStub != nil {
		return fake.GetMissingPvtDataTrackerStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMissingPvtDataTrackerReturns.result1, fake.getMissingPvtDataTrackerReturns.result2
}

func (fake *PeerLedger) GetMissingPvtDataTrackerCallCount() int {
	fake.getMissingPvtDataTrackerMutex.RLock()
	defer fake.getMissingPvtDataTrackerMutex.RUnlock()
	return len(fake.getMissingPvtDataTrackerArgsForCall)
}

func (fake *PeerLedger) GetMissingPvtDataTrackerReturns(result1 ledger.MissingPvtDataTracker, result2 error) {
	fake.GetMissingPvtDataTrackerStub = nil
	fake.getMissingPvtDataTrackerReturns = struct {
		result1 ledger.MissingPvtDataTracker
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetMissingPvtDataTrackerReturnsOnCall(i int, result1 ledger.MissingPvtDataTracker, result2 error) {
	fake.GetMissingPvtDataTrackerStub = nil
	if fake.getMissingPvtDataTrackerReturnsOnCall == nil {
		fake.getMissingPvtDataTrackerReturnsOnCall = make(map[int]struct {
			result1 ledger.MissingPvtDataTracker
			result2 error
		})
	}
	fake.getMissingPvtDataTrackerReturnsOnCall[i] = struct {
		result1 ledger.MissingPvtDataTracker
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getBlockchainInfoMutex.RLock()
	defer fake.getBlockchainInfoMutex.RUnlock()
	fake.getBlockByNumberMutex.RLock()
	defer fake.getBlockByNumberMutex.RUnlock()
	fake.getBlocksIteratorMutex.RLock()
	defer fake.getBlocksIteratorMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	fake.getBlockByHashMutex.RLock()
	defer fake.getBlockByHashMutex.RUnlock()
	fake.getBlockByTxIDMutex.RLock()
	defer fake.getBlockByTxIDMutex.RUnlock()
	fake.getTxValidationCodeByTxIDMutex.RLock()
	defer fake.getTxValidationCodeByTxIDMutex.RUnlock()
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	fake.newHistoryQueryExecutorMutex.RLock()
	defer fake.newHistoryQueryExecutorMutex.RUnlock()
	fake.getPvtDataAndBlockByNumMutex.RLock()
	defer fake.getPvtDataAndBlockByNumMutex.RUnlock()
	fake.getPvtDataByNumMutex.RLock()
	defer fake.getPvtDataByNumMutex.RUnlock()
	fake.commitWithPvtDataMutex.RLock()
	defer fake.commitWithPvtDataMutex.RUnlock()
	fake.purgePrivateDataMutex.RLock()
	defer fake.purgePrivateDataMutex.RUnlock()
	fake.privateDataMinBlockNumMutex.RLock()
	defer fake.privateDataMinBlockNumMutex.RUnlock()
	fake.pruneMutex.RLock()
	defer fake.pruneMutex.RUnlock()
	fake.getConfigHistoryRetrieverMutex.RLock()
	defer fake.getConfigHistoryRetrieverMutex.RUnlock()
	fake.commitPvtDataMutex.RLock()
	defer fake.commitPvtDataMutex.RUnlock()
	fake.getMissingPvtDataTrackerMutex.RLock()
	defer fake.getMissingPvtDataTrackerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PeerLedger) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
