
package mock

import (
	sync "sync"

	ledger "github.com/mcc-github/blockchain/common/ledger"
	ledgera "github.com/mcc-github/blockchain/core/ledger"
	common "github.com/mcc-github/blockchain/protos/common"
	peer "github.com/mcc-github/blockchain/protos/peer"
)

type PeerLedger struct {
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	CommitPvtDataOfOldBlocksStub        func([]*ledgera.BlockPvtData) ([]*ledgera.PvtdataHashMismatch, error)
	commitPvtDataOfOldBlocksMutex       sync.RWMutex
	commitPvtDataOfOldBlocksArgsForCall []struct {
		arg1 []*ledgera.BlockPvtData
	}
	commitPvtDataOfOldBlocksReturns struct {
		result1 []*ledgera.PvtdataHashMismatch
		result2 error
	}
	commitPvtDataOfOldBlocksReturnsOnCall map[int]struct {
		result1 []*ledgera.PvtdataHashMismatch
		result2 error
	}
	CommitWithPvtDataStub        func(*ledgera.BlockAndPvtData) error
	commitWithPvtDataMutex       sync.RWMutex
	commitWithPvtDataArgsForCall []struct {
		arg1 *ledgera.BlockAndPvtData
	}
	commitWithPvtDataReturns struct {
		result1 error
	}
	commitWithPvtDataReturnsOnCall map[int]struct {
		result1 error
	}
	GetBlockByHashStub        func([]byte) (*common.Block, error)
	getBlockByHashMutex       sync.RWMutex
	getBlockByHashArgsForCall []struct {
		arg1 []byte
	}
	getBlockByHashReturns struct {
		result1 *common.Block
		result2 error
	}
	getBlockByHashReturnsOnCall map[int]struct {
		result1 *common.Block
		result2 error
	}
	GetBlockByNumberStub        func(uint64) (*common.Block, error)
	getBlockByNumberMutex       sync.RWMutex
	getBlockByNumberArgsForCall []struct {
		arg1 uint64
	}
	getBlockByNumberReturns struct {
		result1 *common.Block
		result2 error
	}
	getBlockByNumberReturnsOnCall map[int]struct {
		result1 *common.Block
		result2 error
	}
	GetBlockByTxIDStub        func(string) (*common.Block, error)
	getBlockByTxIDMutex       sync.RWMutex
	getBlockByTxIDArgsForCall []struct {
		arg1 string
	}
	getBlockByTxIDReturns struct {
		result1 *common.Block
		result2 error
	}
	getBlockByTxIDReturnsOnCall map[int]struct {
		result1 *common.Block
		result2 error
	}
	GetBlockchainInfoStub        func() (*common.BlockchainInfo, error)
	getBlockchainInfoMutex       sync.RWMutex
	getBlockchainInfoArgsForCall []struct {
	}
	getBlockchainInfoReturns struct {
		result1 *common.BlockchainInfo
		result2 error
	}
	getBlockchainInfoReturnsOnCall map[int]struct {
		result1 *common.BlockchainInfo
		result2 error
	}
	GetBlocksIteratorStub        func(uint64) (ledger.ResultsIterator, error)
	getBlocksIteratorMutex       sync.RWMutex
	getBlocksIteratorArgsForCall []struct {
		arg1 uint64
	}
	getBlocksIteratorReturns struct {
		result1 ledger.ResultsIterator
		result2 error
	}
	getBlocksIteratorReturnsOnCall map[int]struct {
		result1 ledger.ResultsIterator
		result2 error
	}
	GetConfigHistoryRetrieverStub        func() (ledgera.ConfigHistoryRetriever, error)
	getConfigHistoryRetrieverMutex       sync.RWMutex
	getConfigHistoryRetrieverArgsForCall []struct {
	}
	getConfigHistoryRetrieverReturns struct {
		result1 ledgera.ConfigHistoryRetriever
		result2 error
	}
	getConfigHistoryRetrieverReturnsOnCall map[int]struct {
		result1 ledgera.ConfigHistoryRetriever
		result2 error
	}
	GetMissingPvtDataTrackerStub        func() (ledgera.MissingPvtDataTracker, error)
	getMissingPvtDataTrackerMutex       sync.RWMutex
	getMissingPvtDataTrackerArgsForCall []struct {
	}
	getMissingPvtDataTrackerReturns struct {
		result1 ledgera.MissingPvtDataTracker
		result2 error
	}
	getMissingPvtDataTrackerReturnsOnCall map[int]struct {
		result1 ledgera.MissingPvtDataTracker
		result2 error
	}
	GetPvtDataAndBlockByNumStub        func(uint64, ledgera.PvtNsCollFilter) (*ledgera.BlockAndPvtData, error)
	getPvtDataAndBlockByNumMutex       sync.RWMutex
	getPvtDataAndBlockByNumArgsForCall []struct {
		arg1 uint64
		arg2 ledgera.PvtNsCollFilter
	}
	getPvtDataAndBlockByNumReturns struct {
		result1 *ledgera.BlockAndPvtData
		result2 error
	}
	getPvtDataAndBlockByNumReturnsOnCall map[int]struct {
		result1 *ledgera.BlockAndPvtData
		result2 error
	}
	GetPvtDataByNumStub        func(uint64, ledgera.PvtNsCollFilter) ([]*ledgera.TxPvtData, error)
	getPvtDataByNumMutex       sync.RWMutex
	getPvtDataByNumArgsForCall []struct {
		arg1 uint64
		arg2 ledgera.PvtNsCollFilter
	}
	getPvtDataByNumReturns struct {
		result1 []*ledgera.TxPvtData
		result2 error
	}
	getPvtDataByNumReturnsOnCall map[int]struct {
		result1 []*ledgera.TxPvtData
		result2 error
	}
	GetTransactionByIDStub        func(string) (*peer.ProcessedTransaction, error)
	getTransactionByIDMutex       sync.RWMutex
	getTransactionByIDArgsForCall []struct {
		arg1 string
	}
	getTransactionByIDReturns struct {
		result1 *peer.ProcessedTransaction
		result2 error
	}
	getTransactionByIDReturnsOnCall map[int]struct {
		result1 *peer.ProcessedTransaction
		result2 error
	}
	GetTxValidationCodeByTxIDStub        func(string) (peer.TxValidationCode, error)
	getTxValidationCodeByTxIDMutex       sync.RWMutex
	getTxValidationCodeByTxIDArgsForCall []struct {
		arg1 string
	}
	getTxValidationCodeByTxIDReturns struct {
		result1 peer.TxValidationCode
		result2 error
	}
	getTxValidationCodeByTxIDReturnsOnCall map[int]struct {
		result1 peer.TxValidationCode
		result2 error
	}
	NewHistoryQueryExecutorStub        func() (ledgera.HistoryQueryExecutor, error)
	newHistoryQueryExecutorMutex       sync.RWMutex
	newHistoryQueryExecutorArgsForCall []struct {
	}
	newHistoryQueryExecutorReturns struct {
		result1 ledgera.HistoryQueryExecutor
		result2 error
	}
	newHistoryQueryExecutorReturnsOnCall map[int]struct {
		result1 ledgera.HistoryQueryExecutor
		result2 error
	}
	NewQueryExecutorStub        func() (ledgera.QueryExecutor, error)
	newQueryExecutorMutex       sync.RWMutex
	newQueryExecutorArgsForCall []struct {
	}
	newQueryExecutorReturns struct {
		result1 ledgera.QueryExecutor
		result2 error
	}
	newQueryExecutorReturnsOnCall map[int]struct {
		result1 ledgera.QueryExecutor
		result2 error
	}
	NewTxSimulatorStub        func(string) (ledgera.TxSimulator, error)
	newTxSimulatorMutex       sync.RWMutex
	newTxSimulatorArgsForCall []struct {
		arg1 string
	}
	newTxSimulatorReturns struct {
		result1 ledgera.TxSimulator
		result2 error
	}
	newTxSimulatorReturnsOnCall map[int]struct {
		result1 ledgera.TxSimulator
		result2 error
	}
	PrivateDataMinBlockNumStub        func() (uint64, error)
	privateDataMinBlockNumMutex       sync.RWMutex
	privateDataMinBlockNumArgsForCall []struct {
	}
	privateDataMinBlockNumReturns struct {
		result1 uint64
		result2 error
	}
	privateDataMinBlockNumReturnsOnCall map[int]struct {
		result1 uint64
		result2 error
	}
	PruneStub        func(ledger.PrunePolicy) error
	pruneMutex       sync.RWMutex
	pruneArgsForCall []struct {
		arg1 ledger.PrunePolicy
	}
	pruneReturns struct {
		result1 error
	}
	pruneReturnsOnCall map[int]struct {
		result1 error
	}
	PurgePrivateDataStub        func(uint64) error
	purgePrivateDataMutex       sync.RWMutex
	purgePrivateDataArgsForCall []struct {
		arg1 uint64
	}
	purgePrivateDataReturns struct {
		result1 error
	}
	purgePrivateDataReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PeerLedger) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
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

func (fake *PeerLedger) CloseCalls(stub func()) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocks(arg1 []*ledgera.BlockPvtData) ([]*ledgera.PvtdataHashMismatch, error) {
	var arg1Copy []*ledgera.BlockPvtData
	if arg1 != nil {
		arg1Copy = make([]*ledgera.BlockPvtData, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.commitPvtDataOfOldBlocksMutex.Lock()
	ret, specificReturn := fake.commitPvtDataOfOldBlocksReturnsOnCall[len(fake.commitPvtDataOfOldBlocksArgsForCall)]
	fake.commitPvtDataOfOldBlocksArgsForCall = append(fake.commitPvtDataOfOldBlocksArgsForCall, struct {
		arg1 []*ledgera.BlockPvtData
	}{arg1Copy})
	fake.recordInvocation("CommitPvtDataOfOldBlocks", []interface{}{arg1Copy})
	fake.commitPvtDataOfOldBlocksMutex.Unlock()
	if fake.CommitPvtDataOfOldBlocksStub != nil {
		return fake.CommitPvtDataOfOldBlocksStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.commitPvtDataOfOldBlocksReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocksCallCount() int {
	fake.commitPvtDataOfOldBlocksMutex.RLock()
	defer fake.commitPvtDataOfOldBlocksMutex.RUnlock()
	return len(fake.commitPvtDataOfOldBlocksArgsForCall)
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocksCalls(stub func([]*ledgera.BlockPvtData) ([]*ledgera.PvtdataHashMismatch, error)) {
	fake.commitPvtDataOfOldBlocksMutex.Lock()
	defer fake.commitPvtDataOfOldBlocksMutex.Unlock()
	fake.CommitPvtDataOfOldBlocksStub = stub
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocksArgsForCall(i int) []*ledgera.BlockPvtData {
	fake.commitPvtDataOfOldBlocksMutex.RLock()
	defer fake.commitPvtDataOfOldBlocksMutex.RUnlock()
	argsForCall := fake.commitPvtDataOfOldBlocksArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocksReturns(result1 []*ledgera.PvtdataHashMismatch, result2 error) {
	fake.commitPvtDataOfOldBlocksMutex.Lock()
	defer fake.commitPvtDataOfOldBlocksMutex.Unlock()
	fake.CommitPvtDataOfOldBlocksStub = nil
	fake.commitPvtDataOfOldBlocksReturns = struct {
		result1 []*ledgera.PvtdataHashMismatch
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocksReturnsOnCall(i int, result1 []*ledgera.PvtdataHashMismatch, result2 error) {
	fake.commitPvtDataOfOldBlocksMutex.Lock()
	defer fake.commitPvtDataOfOldBlocksMutex.Unlock()
	fake.CommitPvtDataOfOldBlocksStub = nil
	if fake.commitPvtDataOfOldBlocksReturnsOnCall == nil {
		fake.commitPvtDataOfOldBlocksReturnsOnCall = make(map[int]struct {
			result1 []*ledgera.PvtdataHashMismatch
			result2 error
		})
	}
	fake.commitPvtDataOfOldBlocksReturnsOnCall[i] = struct {
		result1 []*ledgera.PvtdataHashMismatch
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) CommitWithPvtData(arg1 *ledgera.BlockAndPvtData) error {
	fake.commitWithPvtDataMutex.Lock()
	ret, specificReturn := fake.commitWithPvtDataReturnsOnCall[len(fake.commitWithPvtDataArgsForCall)]
	fake.commitWithPvtDataArgsForCall = append(fake.commitWithPvtDataArgsForCall, struct {
		arg1 *ledgera.BlockAndPvtData
	}{arg1})
	fake.recordInvocation("CommitWithPvtData", []interface{}{arg1})
	fake.commitWithPvtDataMutex.Unlock()
	if fake.CommitWithPvtDataStub != nil {
		return fake.CommitWithPvtDataStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.commitWithPvtDataReturns
	return fakeReturns.result1
}

func (fake *PeerLedger) CommitWithPvtDataCallCount() int {
	fake.commitWithPvtDataMutex.RLock()
	defer fake.commitWithPvtDataMutex.RUnlock()
	return len(fake.commitWithPvtDataArgsForCall)
}

func (fake *PeerLedger) CommitWithPvtDataCalls(stub func(*ledgera.BlockAndPvtData) error) {
	fake.commitWithPvtDataMutex.Lock()
	defer fake.commitWithPvtDataMutex.Unlock()
	fake.CommitWithPvtDataStub = stub
}

func (fake *PeerLedger) CommitWithPvtDataArgsForCall(i int) *ledgera.BlockAndPvtData {
	fake.commitWithPvtDataMutex.RLock()
	defer fake.commitWithPvtDataMutex.RUnlock()
	argsForCall := fake.commitWithPvtDataArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) CommitWithPvtDataReturns(result1 error) {
	fake.commitWithPvtDataMutex.Lock()
	defer fake.commitWithPvtDataMutex.Unlock()
	fake.CommitWithPvtDataStub = nil
	fake.commitWithPvtDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) CommitWithPvtDataReturnsOnCall(i int, result1 error) {
	fake.commitWithPvtDataMutex.Lock()
	defer fake.commitWithPvtDataMutex.Unlock()
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

func (fake *PeerLedger) GetBlockByHash(arg1 []byte) (*common.Block, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.getBlockByHashMutex.Lock()
	ret, specificReturn := fake.getBlockByHashReturnsOnCall[len(fake.getBlockByHashArgsForCall)]
	fake.getBlockByHashArgsForCall = append(fake.getBlockByHashArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("GetBlockByHash", []interface{}{arg1Copy})
	fake.getBlockByHashMutex.Unlock()
	if fake.GetBlockByHashStub != nil {
		return fake.GetBlockByHashStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBlockByHashReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetBlockByHashCallCount() int {
	fake.getBlockByHashMutex.RLock()
	defer fake.getBlockByHashMutex.RUnlock()
	return len(fake.getBlockByHashArgsForCall)
}

func (fake *PeerLedger) GetBlockByHashCalls(stub func([]byte) (*common.Block, error)) {
	fake.getBlockByHashMutex.Lock()
	defer fake.getBlockByHashMutex.Unlock()
	fake.GetBlockByHashStub = stub
}

func (fake *PeerLedger) GetBlockByHashArgsForCall(i int) []byte {
	fake.getBlockByHashMutex.RLock()
	defer fake.getBlockByHashMutex.RUnlock()
	argsForCall := fake.getBlockByHashArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) GetBlockByHashReturns(result1 *common.Block, result2 error) {
	fake.getBlockByHashMutex.Lock()
	defer fake.getBlockByHashMutex.Unlock()
	fake.GetBlockByHashStub = nil
	fake.getBlockByHashReturns = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByHashReturnsOnCall(i int, result1 *common.Block, result2 error) {
	fake.getBlockByHashMutex.Lock()
	defer fake.getBlockByHashMutex.Unlock()
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

func (fake *PeerLedger) GetBlockByNumber(arg1 uint64) (*common.Block, error) {
	fake.getBlockByNumberMutex.Lock()
	ret, specificReturn := fake.getBlockByNumberReturnsOnCall[len(fake.getBlockByNumberArgsForCall)]
	fake.getBlockByNumberArgsForCall = append(fake.getBlockByNumberArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("GetBlockByNumber", []interface{}{arg1})
	fake.getBlockByNumberMutex.Unlock()
	if fake.GetBlockByNumberStub != nil {
		return fake.GetBlockByNumberStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBlockByNumberReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetBlockByNumberCallCount() int {
	fake.getBlockByNumberMutex.RLock()
	defer fake.getBlockByNumberMutex.RUnlock()
	return len(fake.getBlockByNumberArgsForCall)
}

func (fake *PeerLedger) GetBlockByNumberCalls(stub func(uint64) (*common.Block, error)) {
	fake.getBlockByNumberMutex.Lock()
	defer fake.getBlockByNumberMutex.Unlock()
	fake.GetBlockByNumberStub = stub
}

func (fake *PeerLedger) GetBlockByNumberArgsForCall(i int) uint64 {
	fake.getBlockByNumberMutex.RLock()
	defer fake.getBlockByNumberMutex.RUnlock()
	argsForCall := fake.getBlockByNumberArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) GetBlockByNumberReturns(result1 *common.Block, result2 error) {
	fake.getBlockByNumberMutex.Lock()
	defer fake.getBlockByNumberMutex.Unlock()
	fake.GetBlockByNumberStub = nil
	fake.getBlockByNumberReturns = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByNumberReturnsOnCall(i int, result1 *common.Block, result2 error) {
	fake.getBlockByNumberMutex.Lock()
	defer fake.getBlockByNumberMutex.Unlock()
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

func (fake *PeerLedger) GetBlockByTxID(arg1 string) (*common.Block, error) {
	fake.getBlockByTxIDMutex.Lock()
	ret, specificReturn := fake.getBlockByTxIDReturnsOnCall[len(fake.getBlockByTxIDArgsForCall)]
	fake.getBlockByTxIDArgsForCall = append(fake.getBlockByTxIDArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetBlockByTxID", []interface{}{arg1})
	fake.getBlockByTxIDMutex.Unlock()
	if fake.GetBlockByTxIDStub != nil {
		return fake.GetBlockByTxIDStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBlockByTxIDReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetBlockByTxIDCallCount() int {
	fake.getBlockByTxIDMutex.RLock()
	defer fake.getBlockByTxIDMutex.RUnlock()
	return len(fake.getBlockByTxIDArgsForCall)
}

func (fake *PeerLedger) GetBlockByTxIDCalls(stub func(string) (*common.Block, error)) {
	fake.getBlockByTxIDMutex.Lock()
	defer fake.getBlockByTxIDMutex.Unlock()
	fake.GetBlockByTxIDStub = stub
}

func (fake *PeerLedger) GetBlockByTxIDArgsForCall(i int) string {
	fake.getBlockByTxIDMutex.RLock()
	defer fake.getBlockByTxIDMutex.RUnlock()
	argsForCall := fake.getBlockByTxIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) GetBlockByTxIDReturns(result1 *common.Block, result2 error) {
	fake.getBlockByTxIDMutex.Lock()
	defer fake.getBlockByTxIDMutex.Unlock()
	fake.GetBlockByTxIDStub = nil
	fake.getBlockByTxIDReturns = struct {
		result1 *common.Block
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockByTxIDReturnsOnCall(i int, result1 *common.Block, result2 error) {
	fake.getBlockByTxIDMutex.Lock()
	defer fake.getBlockByTxIDMutex.Unlock()
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

func (fake *PeerLedger) GetBlockchainInfo() (*common.BlockchainInfo, error) {
	fake.getBlockchainInfoMutex.Lock()
	ret, specificReturn := fake.getBlockchainInfoReturnsOnCall[len(fake.getBlockchainInfoArgsForCall)]
	fake.getBlockchainInfoArgsForCall = append(fake.getBlockchainInfoArgsForCall, struct {
	}{})
	fake.recordInvocation("GetBlockchainInfo", []interface{}{})
	fake.getBlockchainInfoMutex.Unlock()
	if fake.GetBlockchainInfoStub != nil {
		return fake.GetBlockchainInfoStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBlockchainInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetBlockchainInfoCallCount() int {
	fake.getBlockchainInfoMutex.RLock()
	defer fake.getBlockchainInfoMutex.RUnlock()
	return len(fake.getBlockchainInfoArgsForCall)
}

func (fake *PeerLedger) GetBlockchainInfoCalls(stub func() (*common.BlockchainInfo, error)) {
	fake.getBlockchainInfoMutex.Lock()
	defer fake.getBlockchainInfoMutex.Unlock()
	fake.GetBlockchainInfoStub = stub
}

func (fake *PeerLedger) GetBlockchainInfoReturns(result1 *common.BlockchainInfo, result2 error) {
	fake.getBlockchainInfoMutex.Lock()
	defer fake.getBlockchainInfoMutex.Unlock()
	fake.GetBlockchainInfoStub = nil
	fake.getBlockchainInfoReturns = struct {
		result1 *common.BlockchainInfo
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlockchainInfoReturnsOnCall(i int, result1 *common.BlockchainInfo, result2 error) {
	fake.getBlockchainInfoMutex.Lock()
	defer fake.getBlockchainInfoMutex.Unlock()
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

func (fake *PeerLedger) GetBlocksIterator(arg1 uint64) (ledger.ResultsIterator, error) {
	fake.getBlocksIteratorMutex.Lock()
	ret, specificReturn := fake.getBlocksIteratorReturnsOnCall[len(fake.getBlocksIteratorArgsForCall)]
	fake.getBlocksIteratorArgsForCall = append(fake.getBlocksIteratorArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("GetBlocksIterator", []interface{}{arg1})
	fake.getBlocksIteratorMutex.Unlock()
	if fake.GetBlocksIteratorStub != nil {
		return fake.GetBlocksIteratorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBlocksIteratorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetBlocksIteratorCallCount() int {
	fake.getBlocksIteratorMutex.RLock()
	defer fake.getBlocksIteratorMutex.RUnlock()
	return len(fake.getBlocksIteratorArgsForCall)
}

func (fake *PeerLedger) GetBlocksIteratorCalls(stub func(uint64) (ledger.ResultsIterator, error)) {
	fake.getBlocksIteratorMutex.Lock()
	defer fake.getBlocksIteratorMutex.Unlock()
	fake.GetBlocksIteratorStub = stub
}

func (fake *PeerLedger) GetBlocksIteratorArgsForCall(i int) uint64 {
	fake.getBlocksIteratorMutex.RLock()
	defer fake.getBlocksIteratorMutex.RUnlock()
	argsForCall := fake.getBlocksIteratorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) GetBlocksIteratorReturns(result1 ledger.ResultsIterator, result2 error) {
	fake.getBlocksIteratorMutex.Lock()
	defer fake.getBlocksIteratorMutex.Unlock()
	fake.GetBlocksIteratorStub = nil
	fake.getBlocksIteratorReturns = struct {
		result1 ledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlocksIteratorReturnsOnCall(i int, result1 ledger.ResultsIterator, result2 error) {
	fake.getBlocksIteratorMutex.Lock()
	defer fake.getBlocksIteratorMutex.Unlock()
	fake.GetBlocksIteratorStub = nil
	if fake.getBlocksIteratorReturnsOnCall == nil {
		fake.getBlocksIteratorReturnsOnCall = make(map[int]struct {
			result1 ledger.ResultsIterator
			result2 error
		})
	}
	fake.getBlocksIteratorReturnsOnCall[i] = struct {
		result1 ledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetConfigHistoryRetriever() (ledgera.ConfigHistoryRetriever, error) {
	fake.getConfigHistoryRetrieverMutex.Lock()
	ret, specificReturn := fake.getConfigHistoryRetrieverReturnsOnCall[len(fake.getConfigHistoryRetrieverArgsForCall)]
	fake.getConfigHistoryRetrieverArgsForCall = append(fake.getConfigHistoryRetrieverArgsForCall, struct {
	}{})
	fake.recordInvocation("GetConfigHistoryRetriever", []interface{}{})
	fake.getConfigHistoryRetrieverMutex.Unlock()
	if fake.GetConfigHistoryRetrieverStub != nil {
		return fake.GetConfigHistoryRetrieverStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getConfigHistoryRetrieverReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetConfigHistoryRetrieverCallCount() int {
	fake.getConfigHistoryRetrieverMutex.RLock()
	defer fake.getConfigHistoryRetrieverMutex.RUnlock()
	return len(fake.getConfigHistoryRetrieverArgsForCall)
}

func (fake *PeerLedger) GetConfigHistoryRetrieverCalls(stub func() (ledgera.ConfigHistoryRetriever, error)) {
	fake.getConfigHistoryRetrieverMutex.Lock()
	defer fake.getConfigHistoryRetrieverMutex.Unlock()
	fake.GetConfigHistoryRetrieverStub = stub
}

func (fake *PeerLedger) GetConfigHistoryRetrieverReturns(result1 ledgera.ConfigHistoryRetriever, result2 error) {
	fake.getConfigHistoryRetrieverMutex.Lock()
	defer fake.getConfigHistoryRetrieverMutex.Unlock()
	fake.GetConfigHistoryRetrieverStub = nil
	fake.getConfigHistoryRetrieverReturns = struct {
		result1 ledgera.ConfigHistoryRetriever
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetConfigHistoryRetrieverReturnsOnCall(i int, result1 ledgera.ConfigHistoryRetriever, result2 error) {
	fake.getConfigHistoryRetrieverMutex.Lock()
	defer fake.getConfigHistoryRetrieverMutex.Unlock()
	fake.GetConfigHistoryRetrieverStub = nil
	if fake.getConfigHistoryRetrieverReturnsOnCall == nil {
		fake.getConfigHistoryRetrieverReturnsOnCall = make(map[int]struct {
			result1 ledgera.ConfigHistoryRetriever
			result2 error
		})
	}
	fake.getConfigHistoryRetrieverReturnsOnCall[i] = struct {
		result1 ledgera.ConfigHistoryRetriever
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetMissingPvtDataTracker() (ledgera.MissingPvtDataTracker, error) {
	fake.getMissingPvtDataTrackerMutex.Lock()
	ret, specificReturn := fake.getMissingPvtDataTrackerReturnsOnCall[len(fake.getMissingPvtDataTrackerArgsForCall)]
	fake.getMissingPvtDataTrackerArgsForCall = append(fake.getMissingPvtDataTrackerArgsForCall, struct {
	}{})
	fake.recordInvocation("GetMissingPvtDataTracker", []interface{}{})
	fake.getMissingPvtDataTrackerMutex.Unlock()
	if fake.GetMissingPvtDataTrackerStub != nil {
		return fake.GetMissingPvtDataTrackerStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getMissingPvtDataTrackerReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetMissingPvtDataTrackerCallCount() int {
	fake.getMissingPvtDataTrackerMutex.RLock()
	defer fake.getMissingPvtDataTrackerMutex.RUnlock()
	return len(fake.getMissingPvtDataTrackerArgsForCall)
}

func (fake *PeerLedger) GetMissingPvtDataTrackerCalls(stub func() (ledgera.MissingPvtDataTracker, error)) {
	fake.getMissingPvtDataTrackerMutex.Lock()
	defer fake.getMissingPvtDataTrackerMutex.Unlock()
	fake.GetMissingPvtDataTrackerStub = stub
}

func (fake *PeerLedger) GetMissingPvtDataTrackerReturns(result1 ledgera.MissingPvtDataTracker, result2 error) {
	fake.getMissingPvtDataTrackerMutex.Lock()
	defer fake.getMissingPvtDataTrackerMutex.Unlock()
	fake.GetMissingPvtDataTrackerStub = nil
	fake.getMissingPvtDataTrackerReturns = struct {
		result1 ledgera.MissingPvtDataTracker
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetMissingPvtDataTrackerReturnsOnCall(i int, result1 ledgera.MissingPvtDataTracker, result2 error) {
	fake.getMissingPvtDataTrackerMutex.Lock()
	defer fake.getMissingPvtDataTrackerMutex.Unlock()
	fake.GetMissingPvtDataTrackerStub = nil
	if fake.getMissingPvtDataTrackerReturnsOnCall == nil {
		fake.getMissingPvtDataTrackerReturnsOnCall = make(map[int]struct {
			result1 ledgera.MissingPvtDataTracker
			result2 error
		})
	}
	fake.getMissingPvtDataTrackerReturnsOnCall[i] = struct {
		result1 ledgera.MissingPvtDataTracker
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataAndBlockByNum(arg1 uint64, arg2 ledgera.PvtNsCollFilter) (*ledgera.BlockAndPvtData, error) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	ret, specificReturn := fake.getPvtDataAndBlockByNumReturnsOnCall[len(fake.getPvtDataAndBlockByNumArgsForCall)]
	fake.getPvtDataAndBlockByNumArgsForCall = append(fake.getPvtDataAndBlockByNumArgsForCall, struct {
		arg1 uint64
		arg2 ledgera.PvtNsCollFilter
	}{arg1, arg2})
	fake.recordInvocation("GetPvtDataAndBlockByNum", []interface{}{arg1, arg2})
	fake.getPvtDataAndBlockByNumMutex.Unlock()
	if fake.GetPvtDataAndBlockByNumStub != nil {
		return fake.GetPvtDataAndBlockByNumStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPvtDataAndBlockByNumReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumCallCount() int {
	fake.getPvtDataAndBlockByNumMutex.RLock()
	defer fake.getPvtDataAndBlockByNumMutex.RUnlock()
	return len(fake.getPvtDataAndBlockByNumArgsForCall)
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumCalls(stub func(uint64, ledgera.PvtNsCollFilter) (*ledgera.BlockAndPvtData, error)) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	defer fake.getPvtDataAndBlockByNumMutex.Unlock()
	fake.GetPvtDataAndBlockByNumStub = stub
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumArgsForCall(i int) (uint64, ledgera.PvtNsCollFilter) {
	fake.getPvtDataAndBlockByNumMutex.RLock()
	defer fake.getPvtDataAndBlockByNumMutex.RUnlock()
	argsForCall := fake.getPvtDataAndBlockByNumArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumReturns(result1 *ledgera.BlockAndPvtData, result2 error) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	defer fake.getPvtDataAndBlockByNumMutex.Unlock()
	fake.GetPvtDataAndBlockByNumStub = nil
	fake.getPvtDataAndBlockByNumReturns = struct {
		result1 *ledgera.BlockAndPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumReturnsOnCall(i int, result1 *ledgera.BlockAndPvtData, result2 error) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	defer fake.getPvtDataAndBlockByNumMutex.Unlock()
	fake.GetPvtDataAndBlockByNumStub = nil
	if fake.getPvtDataAndBlockByNumReturnsOnCall == nil {
		fake.getPvtDataAndBlockByNumReturnsOnCall = make(map[int]struct {
			result1 *ledgera.BlockAndPvtData
			result2 error
		})
	}
	fake.getPvtDataAndBlockByNumReturnsOnCall[i] = struct {
		result1 *ledgera.BlockAndPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataByNum(arg1 uint64, arg2 ledgera.PvtNsCollFilter) ([]*ledgera.TxPvtData, error) {
	fake.getPvtDataByNumMutex.Lock()
	ret, specificReturn := fake.getPvtDataByNumReturnsOnCall[len(fake.getPvtDataByNumArgsForCall)]
	fake.getPvtDataByNumArgsForCall = append(fake.getPvtDataByNumArgsForCall, struct {
		arg1 uint64
		arg2 ledgera.PvtNsCollFilter
	}{arg1, arg2})
	fake.recordInvocation("GetPvtDataByNum", []interface{}{arg1, arg2})
	fake.getPvtDataByNumMutex.Unlock()
	if fake.GetPvtDataByNumStub != nil {
		return fake.GetPvtDataByNumStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPvtDataByNumReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetPvtDataByNumCallCount() int {
	fake.getPvtDataByNumMutex.RLock()
	defer fake.getPvtDataByNumMutex.RUnlock()
	return len(fake.getPvtDataByNumArgsForCall)
}

func (fake *PeerLedger) GetPvtDataByNumCalls(stub func(uint64, ledgera.PvtNsCollFilter) ([]*ledgera.TxPvtData, error)) {
	fake.getPvtDataByNumMutex.Lock()
	defer fake.getPvtDataByNumMutex.Unlock()
	fake.GetPvtDataByNumStub = stub
}

func (fake *PeerLedger) GetPvtDataByNumArgsForCall(i int) (uint64, ledgera.PvtNsCollFilter) {
	fake.getPvtDataByNumMutex.RLock()
	defer fake.getPvtDataByNumMutex.RUnlock()
	argsForCall := fake.getPvtDataByNumArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PeerLedger) GetPvtDataByNumReturns(result1 []*ledgera.TxPvtData, result2 error) {
	fake.getPvtDataByNumMutex.Lock()
	defer fake.getPvtDataByNumMutex.Unlock()
	fake.GetPvtDataByNumStub = nil
	fake.getPvtDataByNumReturns = struct {
		result1 []*ledgera.TxPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataByNumReturnsOnCall(i int, result1 []*ledgera.TxPvtData, result2 error) {
	fake.getPvtDataByNumMutex.Lock()
	defer fake.getPvtDataByNumMutex.Unlock()
	fake.GetPvtDataByNumStub = nil
	if fake.getPvtDataByNumReturnsOnCall == nil {
		fake.getPvtDataByNumReturnsOnCall = make(map[int]struct {
			result1 []*ledgera.TxPvtData
			result2 error
		})
	}
	fake.getPvtDataByNumReturnsOnCall[i] = struct {
		result1 []*ledgera.TxPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTransactionByID(arg1 string) (*peer.ProcessedTransaction, error) {
	fake.getTransactionByIDMutex.Lock()
	ret, specificReturn := fake.getTransactionByIDReturnsOnCall[len(fake.getTransactionByIDArgsForCall)]
	fake.getTransactionByIDArgsForCall = append(fake.getTransactionByIDArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetTransactionByID", []interface{}{arg1})
	fake.getTransactionByIDMutex.Unlock()
	if fake.GetTransactionByIDStub != nil {
		return fake.GetTransactionByIDStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getTransactionByIDReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetTransactionByIDCallCount() int {
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	return len(fake.getTransactionByIDArgsForCall)
}

func (fake *PeerLedger) GetTransactionByIDCalls(stub func(string) (*peer.ProcessedTransaction, error)) {
	fake.getTransactionByIDMutex.Lock()
	defer fake.getTransactionByIDMutex.Unlock()
	fake.GetTransactionByIDStub = stub
}

func (fake *PeerLedger) GetTransactionByIDArgsForCall(i int) string {
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	argsForCall := fake.getTransactionByIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) GetTransactionByIDReturns(result1 *peer.ProcessedTransaction, result2 error) {
	fake.getTransactionByIDMutex.Lock()
	defer fake.getTransactionByIDMutex.Unlock()
	fake.GetTransactionByIDStub = nil
	fake.getTransactionByIDReturns = struct {
		result1 *peer.ProcessedTransaction
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTransactionByIDReturnsOnCall(i int, result1 *peer.ProcessedTransaction, result2 error) {
	fake.getTransactionByIDMutex.Lock()
	defer fake.getTransactionByIDMutex.Unlock()
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

func (fake *PeerLedger) GetTxValidationCodeByTxID(arg1 string) (peer.TxValidationCode, error) {
	fake.getTxValidationCodeByTxIDMutex.Lock()
	ret, specificReturn := fake.getTxValidationCodeByTxIDReturnsOnCall[len(fake.getTxValidationCodeByTxIDArgsForCall)]
	fake.getTxValidationCodeByTxIDArgsForCall = append(fake.getTxValidationCodeByTxIDArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetTxValidationCodeByTxID", []interface{}{arg1})
	fake.getTxValidationCodeByTxIDMutex.Unlock()
	if fake.GetTxValidationCodeByTxIDStub != nil {
		return fake.GetTxValidationCodeByTxIDStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getTxValidationCodeByTxIDReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDCallCount() int {
	fake.getTxValidationCodeByTxIDMutex.RLock()
	defer fake.getTxValidationCodeByTxIDMutex.RUnlock()
	return len(fake.getTxValidationCodeByTxIDArgsForCall)
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDCalls(stub func(string) (peer.TxValidationCode, error)) {
	fake.getTxValidationCodeByTxIDMutex.Lock()
	defer fake.getTxValidationCodeByTxIDMutex.Unlock()
	fake.GetTxValidationCodeByTxIDStub = stub
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDArgsForCall(i int) string {
	fake.getTxValidationCodeByTxIDMutex.RLock()
	defer fake.getTxValidationCodeByTxIDMutex.RUnlock()
	argsForCall := fake.getTxValidationCodeByTxIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDReturns(result1 peer.TxValidationCode, result2 error) {
	fake.getTxValidationCodeByTxIDMutex.Lock()
	defer fake.getTxValidationCodeByTxIDMutex.Unlock()
	fake.GetTxValidationCodeByTxIDStub = nil
	fake.getTxValidationCodeByTxIDReturns = struct {
		result1 peer.TxValidationCode
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDReturnsOnCall(i int, result1 peer.TxValidationCode, result2 error) {
	fake.getTxValidationCodeByTxIDMutex.Lock()
	defer fake.getTxValidationCodeByTxIDMutex.Unlock()
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

func (fake *PeerLedger) NewHistoryQueryExecutor() (ledgera.HistoryQueryExecutor, error) {
	fake.newHistoryQueryExecutorMutex.Lock()
	ret, specificReturn := fake.newHistoryQueryExecutorReturnsOnCall[len(fake.newHistoryQueryExecutorArgsForCall)]
	fake.newHistoryQueryExecutorArgsForCall = append(fake.newHistoryQueryExecutorArgsForCall, struct {
	}{})
	fake.recordInvocation("NewHistoryQueryExecutor", []interface{}{})
	fake.newHistoryQueryExecutorMutex.Unlock()
	if fake.NewHistoryQueryExecutorStub != nil {
		return fake.NewHistoryQueryExecutorStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newHistoryQueryExecutorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) NewHistoryQueryExecutorCallCount() int {
	fake.newHistoryQueryExecutorMutex.RLock()
	defer fake.newHistoryQueryExecutorMutex.RUnlock()
	return len(fake.newHistoryQueryExecutorArgsForCall)
}

func (fake *PeerLedger) NewHistoryQueryExecutorCalls(stub func() (ledgera.HistoryQueryExecutor, error)) {
	fake.newHistoryQueryExecutorMutex.Lock()
	defer fake.newHistoryQueryExecutorMutex.Unlock()
	fake.NewHistoryQueryExecutorStub = stub
}

func (fake *PeerLedger) NewHistoryQueryExecutorReturns(result1 ledgera.HistoryQueryExecutor, result2 error) {
	fake.newHistoryQueryExecutorMutex.Lock()
	defer fake.newHistoryQueryExecutorMutex.Unlock()
	fake.NewHistoryQueryExecutorStub = nil
	fake.newHistoryQueryExecutorReturns = struct {
		result1 ledgera.HistoryQueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewHistoryQueryExecutorReturnsOnCall(i int, result1 ledgera.HistoryQueryExecutor, result2 error) {
	fake.newHistoryQueryExecutorMutex.Lock()
	defer fake.newHistoryQueryExecutorMutex.Unlock()
	fake.NewHistoryQueryExecutorStub = nil
	if fake.newHistoryQueryExecutorReturnsOnCall == nil {
		fake.newHistoryQueryExecutorReturnsOnCall = make(map[int]struct {
			result1 ledgera.HistoryQueryExecutor
			result2 error
		})
	}
	fake.newHistoryQueryExecutorReturnsOnCall[i] = struct {
		result1 ledgera.HistoryQueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewQueryExecutor() (ledgera.QueryExecutor, error) {
	fake.newQueryExecutorMutex.Lock()
	ret, specificReturn := fake.newQueryExecutorReturnsOnCall[len(fake.newQueryExecutorArgsForCall)]
	fake.newQueryExecutorArgsForCall = append(fake.newQueryExecutorArgsForCall, struct {
	}{})
	fake.recordInvocation("NewQueryExecutor", []interface{}{})
	fake.newQueryExecutorMutex.Unlock()
	if fake.NewQueryExecutorStub != nil {
		return fake.NewQueryExecutorStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newQueryExecutorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) NewQueryExecutorCallCount() int {
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	return len(fake.newQueryExecutorArgsForCall)
}

func (fake *PeerLedger) NewQueryExecutorCalls(stub func() (ledgera.QueryExecutor, error)) {
	fake.newQueryExecutorMutex.Lock()
	defer fake.newQueryExecutorMutex.Unlock()
	fake.NewQueryExecutorStub = stub
}

func (fake *PeerLedger) NewQueryExecutorReturns(result1 ledgera.QueryExecutor, result2 error) {
	fake.newQueryExecutorMutex.Lock()
	defer fake.newQueryExecutorMutex.Unlock()
	fake.NewQueryExecutorStub = nil
	fake.newQueryExecutorReturns = struct {
		result1 ledgera.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewQueryExecutorReturnsOnCall(i int, result1 ledgera.QueryExecutor, result2 error) {
	fake.newQueryExecutorMutex.Lock()
	defer fake.newQueryExecutorMutex.Unlock()
	fake.NewQueryExecutorStub = nil
	if fake.newQueryExecutorReturnsOnCall == nil {
		fake.newQueryExecutorReturnsOnCall = make(map[int]struct {
			result1 ledgera.QueryExecutor
			result2 error
		})
	}
	fake.newQueryExecutorReturnsOnCall[i] = struct {
		result1 ledgera.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewTxSimulator(arg1 string) (ledgera.TxSimulator, error) {
	fake.newTxSimulatorMutex.Lock()
	ret, specificReturn := fake.newTxSimulatorReturnsOnCall[len(fake.newTxSimulatorArgsForCall)]
	fake.newTxSimulatorArgsForCall = append(fake.newTxSimulatorArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("NewTxSimulator", []interface{}{arg1})
	fake.newTxSimulatorMutex.Unlock()
	if fake.NewTxSimulatorStub != nil {
		return fake.NewTxSimulatorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newTxSimulatorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) NewTxSimulatorCallCount() int {
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	return len(fake.newTxSimulatorArgsForCall)
}

func (fake *PeerLedger) NewTxSimulatorCalls(stub func(string) (ledgera.TxSimulator, error)) {
	fake.newTxSimulatorMutex.Lock()
	defer fake.newTxSimulatorMutex.Unlock()
	fake.NewTxSimulatorStub = stub
}

func (fake *PeerLedger) NewTxSimulatorArgsForCall(i int) string {
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	argsForCall := fake.newTxSimulatorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) NewTxSimulatorReturns(result1 ledgera.TxSimulator, result2 error) {
	fake.newTxSimulatorMutex.Lock()
	defer fake.newTxSimulatorMutex.Unlock()
	fake.NewTxSimulatorStub = nil
	fake.newTxSimulatorReturns = struct {
		result1 ledgera.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewTxSimulatorReturnsOnCall(i int, result1 ledgera.TxSimulator, result2 error) {
	fake.newTxSimulatorMutex.Lock()
	defer fake.newTxSimulatorMutex.Unlock()
	fake.NewTxSimulatorStub = nil
	if fake.newTxSimulatorReturnsOnCall == nil {
		fake.newTxSimulatorReturnsOnCall = make(map[int]struct {
			result1 ledgera.TxSimulator
			result2 error
		})
	}
	fake.newTxSimulatorReturnsOnCall[i] = struct {
		result1 ledgera.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) PrivateDataMinBlockNum() (uint64, error) {
	fake.privateDataMinBlockNumMutex.Lock()
	ret, specificReturn := fake.privateDataMinBlockNumReturnsOnCall[len(fake.privateDataMinBlockNumArgsForCall)]
	fake.privateDataMinBlockNumArgsForCall = append(fake.privateDataMinBlockNumArgsForCall, struct {
	}{})
	fake.recordInvocation("PrivateDataMinBlockNum", []interface{}{})
	fake.privateDataMinBlockNumMutex.Unlock()
	if fake.PrivateDataMinBlockNumStub != nil {
		return fake.PrivateDataMinBlockNumStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.privateDataMinBlockNumReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) PrivateDataMinBlockNumCallCount() int {
	fake.privateDataMinBlockNumMutex.RLock()
	defer fake.privateDataMinBlockNumMutex.RUnlock()
	return len(fake.privateDataMinBlockNumArgsForCall)
}

func (fake *PeerLedger) PrivateDataMinBlockNumCalls(stub func() (uint64, error)) {
	fake.privateDataMinBlockNumMutex.Lock()
	defer fake.privateDataMinBlockNumMutex.Unlock()
	fake.PrivateDataMinBlockNumStub = stub
}

func (fake *PeerLedger) PrivateDataMinBlockNumReturns(result1 uint64, result2 error) {
	fake.privateDataMinBlockNumMutex.Lock()
	defer fake.privateDataMinBlockNumMutex.Unlock()
	fake.PrivateDataMinBlockNumStub = nil
	fake.privateDataMinBlockNumReturns = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) PrivateDataMinBlockNumReturnsOnCall(i int, result1 uint64, result2 error) {
	fake.privateDataMinBlockNumMutex.Lock()
	defer fake.privateDataMinBlockNumMutex.Unlock()
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

func (fake *PeerLedger) Prune(arg1 ledger.PrunePolicy) error {
	fake.pruneMutex.Lock()
	ret, specificReturn := fake.pruneReturnsOnCall[len(fake.pruneArgsForCall)]
	fake.pruneArgsForCall = append(fake.pruneArgsForCall, struct {
		arg1 ledger.PrunePolicy
	}{arg1})
	fake.recordInvocation("Prune", []interface{}{arg1})
	fake.pruneMutex.Unlock()
	if fake.PruneStub != nil {
		return fake.PruneStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.pruneReturns
	return fakeReturns.result1
}

func (fake *PeerLedger) PruneCallCount() int {
	fake.pruneMutex.RLock()
	defer fake.pruneMutex.RUnlock()
	return len(fake.pruneArgsForCall)
}

func (fake *PeerLedger) PruneCalls(stub func(ledger.PrunePolicy) error) {
	fake.pruneMutex.Lock()
	defer fake.pruneMutex.Unlock()
	fake.PruneStub = stub
}

func (fake *PeerLedger) PruneArgsForCall(i int) ledger.PrunePolicy {
	fake.pruneMutex.RLock()
	defer fake.pruneMutex.RUnlock()
	argsForCall := fake.pruneArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) PruneReturns(result1 error) {
	fake.pruneMutex.Lock()
	defer fake.pruneMutex.Unlock()
	fake.PruneStub = nil
	fake.pruneReturns = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) PruneReturnsOnCall(i int, result1 error) {
	fake.pruneMutex.Lock()
	defer fake.pruneMutex.Unlock()
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

func (fake *PeerLedger) PurgePrivateData(arg1 uint64) error {
	fake.purgePrivateDataMutex.Lock()
	ret, specificReturn := fake.purgePrivateDataReturnsOnCall[len(fake.purgePrivateDataArgsForCall)]
	fake.purgePrivateDataArgsForCall = append(fake.purgePrivateDataArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("PurgePrivateData", []interface{}{arg1})
	fake.purgePrivateDataMutex.Unlock()
	if fake.PurgePrivateDataStub != nil {
		return fake.PurgePrivateDataStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.purgePrivateDataReturns
	return fakeReturns.result1
}

func (fake *PeerLedger) PurgePrivateDataCallCount() int {
	fake.purgePrivateDataMutex.RLock()
	defer fake.purgePrivateDataMutex.RUnlock()
	return len(fake.purgePrivateDataArgsForCall)
}

func (fake *PeerLedger) PurgePrivateDataCalls(stub func(uint64) error) {
	fake.purgePrivateDataMutex.Lock()
	defer fake.purgePrivateDataMutex.Unlock()
	fake.PurgePrivateDataStub = stub
}

func (fake *PeerLedger) PurgePrivateDataArgsForCall(i int) uint64 {
	fake.purgePrivateDataMutex.RLock()
	defer fake.purgePrivateDataMutex.RUnlock()
	argsForCall := fake.purgePrivateDataArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) PurgePrivateDataReturns(result1 error) {
	fake.purgePrivateDataMutex.Lock()
	defer fake.purgePrivateDataMutex.Unlock()
	fake.PurgePrivateDataStub = nil
	fake.purgePrivateDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) PurgePrivateDataReturnsOnCall(i int, result1 error) {
	fake.purgePrivateDataMutex.Lock()
	defer fake.purgePrivateDataMutex.Unlock()
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

func (fake *PeerLedger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.commitPvtDataOfOldBlocksMutex.RLock()
	defer fake.commitPvtDataOfOldBlocksMutex.RUnlock()
	fake.commitWithPvtDataMutex.RLock()
	defer fake.commitWithPvtDataMutex.RUnlock()
	fake.getBlockByHashMutex.RLock()
	defer fake.getBlockByHashMutex.RUnlock()
	fake.getBlockByNumberMutex.RLock()
	defer fake.getBlockByNumberMutex.RUnlock()
	fake.getBlockByTxIDMutex.RLock()
	defer fake.getBlockByTxIDMutex.RUnlock()
	fake.getBlockchainInfoMutex.RLock()
	defer fake.getBlockchainInfoMutex.RUnlock()
	fake.getBlocksIteratorMutex.RLock()
	defer fake.getBlocksIteratorMutex.RUnlock()
	fake.getConfigHistoryRetrieverMutex.RLock()
	defer fake.getConfigHistoryRetrieverMutex.RUnlock()
	fake.getMissingPvtDataTrackerMutex.RLock()
	defer fake.getMissingPvtDataTrackerMutex.RUnlock()
	fake.getPvtDataAndBlockByNumMutex.RLock()
	defer fake.getPvtDataAndBlockByNumMutex.RUnlock()
	fake.getPvtDataByNumMutex.RLock()
	defer fake.getPvtDataByNumMutex.RUnlock()
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	fake.getTxValidationCodeByTxIDMutex.RLock()
	defer fake.getTxValidationCodeByTxIDMutex.RUnlock()
	fake.newHistoryQueryExecutorMutex.RLock()
	defer fake.newHistoryQueryExecutorMutex.RUnlock()
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	fake.privateDataMinBlockNumMutex.RLock()
	defer fake.privateDataMinBlockNumMutex.RUnlock()
	fake.pruneMutex.RLock()
	defer fake.pruneMutex.RUnlock()
	fake.purgePrivateDataMutex.RLock()
	defer fake.purgePrivateDataMutex.RUnlock()
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
