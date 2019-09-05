
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	peera "github.com/mcc-github/blockchain-protos-go/peer"
	ledgera "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger"
)

type PeerLedger struct {
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	CommitLegacyStub        func(*ledger.BlockAndPvtData, *ledger.CommitOptions) error
	commitLegacyMutex       sync.RWMutex
	commitLegacyArgsForCall []struct {
		arg1 *ledger.BlockAndPvtData
		arg2 *ledger.CommitOptions
	}
	commitLegacyReturns struct {
		result1 error
	}
	commitLegacyReturnsOnCall map[int]struct {
		result1 error
	}
	CommitPvtDataOfOldBlocksStub        func([]*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error)
	commitPvtDataOfOldBlocksMutex       sync.RWMutex
	commitPvtDataOfOldBlocksArgsForCall []struct {
		arg1 []*ledger.BlockPvtData
	}
	commitPvtDataOfOldBlocksReturns struct {
		result1 []*ledger.PvtdataHashMismatch
		result2 error
	}
	commitPvtDataOfOldBlocksReturnsOnCall map[int]struct {
		result1 []*ledger.PvtdataHashMismatch
		result2 error
	}
	DoesPvtDataInfoExistStub        func(uint64) (bool, error)
	doesPvtDataInfoExistMutex       sync.RWMutex
	doesPvtDataInfoExistArgsForCall []struct {
		arg1 uint64
	}
	doesPvtDataInfoExistReturns struct {
		result1 bool
		result2 error
	}
	doesPvtDataInfoExistReturnsOnCall map[int]struct {
		result1 bool
		result2 error
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
	GetBlocksIteratorStub        func(uint64) (ledgera.ResultsIterator, error)
	getBlocksIteratorMutex       sync.RWMutex
	getBlocksIteratorArgsForCall []struct {
		arg1 uint64
	}
	getBlocksIteratorReturns struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	getBlocksIteratorReturnsOnCall map[int]struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	GetConfigHistoryRetrieverStub        func() (ledger.ConfigHistoryRetriever, error)
	getConfigHistoryRetrieverMutex       sync.RWMutex
	getConfigHistoryRetrieverArgsForCall []struct {
	}
	getConfigHistoryRetrieverReturns struct {
		result1 ledger.ConfigHistoryRetriever
		result2 error
	}
	getConfigHistoryRetrieverReturnsOnCall map[int]struct {
		result1 ledger.ConfigHistoryRetriever
		result2 error
	}
	GetMissingPvtDataTrackerStub        func() (ledger.MissingPvtDataTracker, error)
	getMissingPvtDataTrackerMutex       sync.RWMutex
	getMissingPvtDataTrackerArgsForCall []struct {
	}
	getMissingPvtDataTrackerReturns struct {
		result1 ledger.MissingPvtDataTracker
		result2 error
	}
	getMissingPvtDataTrackerReturnsOnCall map[int]struct {
		result1 ledger.MissingPvtDataTracker
		result2 error
	}
	GetPvtDataAndBlockByNumStub        func(uint64, ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error)
	getPvtDataAndBlockByNumMutex       sync.RWMutex
	getPvtDataAndBlockByNumArgsForCall []struct {
		arg1 uint64
		arg2 ledger.PvtNsCollFilter
	}
	getPvtDataAndBlockByNumReturns struct {
		result1 *ledger.BlockAndPvtData
		result2 error
	}
	getPvtDataAndBlockByNumReturnsOnCall map[int]struct {
		result1 *ledger.BlockAndPvtData
		result2 error
	}
	GetPvtDataByNumStub        func(uint64, ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)
	getPvtDataByNumMutex       sync.RWMutex
	getPvtDataByNumArgsForCall []struct {
		arg1 uint64
		arg2 ledger.PvtNsCollFilter
	}
	getPvtDataByNumReturns struct {
		result1 []*ledger.TxPvtData
		result2 error
	}
	getPvtDataByNumReturnsOnCall map[int]struct {
		result1 []*ledger.TxPvtData
		result2 error
	}
	GetTransactionByIDStub        func(string) (*peera.ProcessedTransaction, error)
	getTransactionByIDMutex       sync.RWMutex
	getTransactionByIDArgsForCall []struct {
		arg1 string
	}
	getTransactionByIDReturns struct {
		result1 *peera.ProcessedTransaction
		result2 error
	}
	getTransactionByIDReturnsOnCall map[int]struct {
		result1 *peera.ProcessedTransaction
		result2 error
	}
	GetTxValidationCodeByTxIDStub        func(string) (peera.TxValidationCode, error)
	getTxValidationCodeByTxIDMutex       sync.RWMutex
	getTxValidationCodeByTxIDArgsForCall []struct {
		arg1 string
	}
	getTxValidationCodeByTxIDReturns struct {
		result1 peera.TxValidationCode
		result2 error
	}
	getTxValidationCodeByTxIDReturnsOnCall map[int]struct {
		result1 peera.TxValidationCode
		result2 error
	}
	NewHistoryQueryExecutorStub        func() (ledger.HistoryQueryExecutor, error)
	newHistoryQueryExecutorMutex       sync.RWMutex
	newHistoryQueryExecutorArgsForCall []struct {
	}
	newHistoryQueryExecutorReturns struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}
	newHistoryQueryExecutorReturnsOnCall map[int]struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}
	NewQueryExecutorStub        func() (ledger.QueryExecutor, error)
	newQueryExecutorMutex       sync.RWMutex
	newQueryExecutorArgsForCall []struct {
	}
	newQueryExecutorReturns struct {
		result1 ledger.QueryExecutor
		result2 error
	}
	newQueryExecutorReturnsOnCall map[int]struct {
		result1 ledger.QueryExecutor
		result2 error
	}
	NewTxSimulatorStub        func(string) (ledger.TxSimulator, error)
	newTxSimulatorMutex       sync.RWMutex
	newTxSimulatorArgsForCall []struct {
		arg1 string
	}
	newTxSimulatorReturns struct {
		result1 ledger.TxSimulator
		result2 error
	}
	newTxSimulatorReturnsOnCall map[int]struct {
		result1 ledger.TxSimulator
		result2 error
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

func (fake *PeerLedger) CommitLegacy(arg1 *ledger.BlockAndPvtData, arg2 *ledger.CommitOptions) error {
	fake.commitLegacyMutex.Lock()
	ret, specificReturn := fake.commitLegacyReturnsOnCall[len(fake.commitLegacyArgsForCall)]
	fake.commitLegacyArgsForCall = append(fake.commitLegacyArgsForCall, struct {
		arg1 *ledger.BlockAndPvtData
		arg2 *ledger.CommitOptions
	}{arg1, arg2})
	fake.recordInvocation("CommitLegacy", []interface{}{arg1, arg2})
	fake.commitLegacyMutex.Unlock()
	if fake.CommitLegacyStub != nil {
		return fake.CommitLegacyStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.commitLegacyReturns
	return fakeReturns.result1
}

func (fake *PeerLedger) CommitLegacyCallCount() int {
	fake.commitLegacyMutex.RLock()
	defer fake.commitLegacyMutex.RUnlock()
	return len(fake.commitLegacyArgsForCall)
}

func (fake *PeerLedger) CommitLegacyCalls(stub func(*ledger.BlockAndPvtData, *ledger.CommitOptions) error) {
	fake.commitLegacyMutex.Lock()
	defer fake.commitLegacyMutex.Unlock()
	fake.CommitLegacyStub = stub
}

func (fake *PeerLedger) CommitLegacyArgsForCall(i int) (*ledger.BlockAndPvtData, *ledger.CommitOptions) {
	fake.commitLegacyMutex.RLock()
	defer fake.commitLegacyMutex.RUnlock()
	argsForCall := fake.commitLegacyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PeerLedger) CommitLegacyReturns(result1 error) {
	fake.commitLegacyMutex.Lock()
	defer fake.commitLegacyMutex.Unlock()
	fake.CommitLegacyStub = nil
	fake.commitLegacyReturns = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) CommitLegacyReturnsOnCall(i int, result1 error) {
	fake.commitLegacyMutex.Lock()
	defer fake.commitLegacyMutex.Unlock()
	fake.CommitLegacyStub = nil
	if fake.commitLegacyReturnsOnCall == nil {
		fake.commitLegacyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.commitLegacyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocks(arg1 []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error) {
	var arg1Copy []*ledger.BlockPvtData
	if arg1 != nil {
		arg1Copy = make([]*ledger.BlockPvtData, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.commitPvtDataOfOldBlocksMutex.Lock()
	ret, specificReturn := fake.commitPvtDataOfOldBlocksReturnsOnCall[len(fake.commitPvtDataOfOldBlocksArgsForCall)]
	fake.commitPvtDataOfOldBlocksArgsForCall = append(fake.commitPvtDataOfOldBlocksArgsForCall, struct {
		arg1 []*ledger.BlockPvtData
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

func (fake *PeerLedger) CommitPvtDataOfOldBlocksCalls(stub func([]*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error)) {
	fake.commitPvtDataOfOldBlocksMutex.Lock()
	defer fake.commitPvtDataOfOldBlocksMutex.Unlock()
	fake.CommitPvtDataOfOldBlocksStub = stub
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocksArgsForCall(i int) []*ledger.BlockPvtData {
	fake.commitPvtDataOfOldBlocksMutex.RLock()
	defer fake.commitPvtDataOfOldBlocksMutex.RUnlock()
	argsForCall := fake.commitPvtDataOfOldBlocksArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocksReturns(result1 []*ledger.PvtdataHashMismatch, result2 error) {
	fake.commitPvtDataOfOldBlocksMutex.Lock()
	defer fake.commitPvtDataOfOldBlocksMutex.Unlock()
	fake.CommitPvtDataOfOldBlocksStub = nil
	fake.commitPvtDataOfOldBlocksReturns = struct {
		result1 []*ledger.PvtdataHashMismatch
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) CommitPvtDataOfOldBlocksReturnsOnCall(i int, result1 []*ledger.PvtdataHashMismatch, result2 error) {
	fake.commitPvtDataOfOldBlocksMutex.Lock()
	defer fake.commitPvtDataOfOldBlocksMutex.Unlock()
	fake.CommitPvtDataOfOldBlocksStub = nil
	if fake.commitPvtDataOfOldBlocksReturnsOnCall == nil {
		fake.commitPvtDataOfOldBlocksReturnsOnCall = make(map[int]struct {
			result1 []*ledger.PvtdataHashMismatch
			result2 error
		})
	}
	fake.commitPvtDataOfOldBlocksReturnsOnCall[i] = struct {
		result1 []*ledger.PvtdataHashMismatch
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) DoesPvtDataInfoExist(arg1 uint64) (bool, error) {
	fake.doesPvtDataInfoExistMutex.Lock()
	ret, specificReturn := fake.doesPvtDataInfoExistReturnsOnCall[len(fake.doesPvtDataInfoExistArgsForCall)]
	fake.doesPvtDataInfoExistArgsForCall = append(fake.doesPvtDataInfoExistArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("DoesPvtDataInfoExist", []interface{}{arg1})
	fake.doesPvtDataInfoExistMutex.Unlock()
	if fake.DoesPvtDataInfoExistStub != nil {
		return fake.DoesPvtDataInfoExistStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.doesPvtDataInfoExistReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerLedger) DoesPvtDataInfoExistCallCount() int {
	fake.doesPvtDataInfoExistMutex.RLock()
	defer fake.doesPvtDataInfoExistMutex.RUnlock()
	return len(fake.doesPvtDataInfoExistArgsForCall)
}

func (fake *PeerLedger) DoesPvtDataInfoExistCalls(stub func(uint64) (bool, error)) {
	fake.doesPvtDataInfoExistMutex.Lock()
	defer fake.doesPvtDataInfoExistMutex.Unlock()
	fake.DoesPvtDataInfoExistStub = stub
}

func (fake *PeerLedger) DoesPvtDataInfoExistArgsForCall(i int) uint64 {
	fake.doesPvtDataInfoExistMutex.RLock()
	defer fake.doesPvtDataInfoExistMutex.RUnlock()
	argsForCall := fake.doesPvtDataInfoExistArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerLedger) DoesPvtDataInfoExistReturns(result1 bool, result2 error) {
	fake.doesPvtDataInfoExistMutex.Lock()
	defer fake.doesPvtDataInfoExistMutex.Unlock()
	fake.DoesPvtDataInfoExistStub = nil
	fake.doesPvtDataInfoExistReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) DoesPvtDataInfoExistReturnsOnCall(i int, result1 bool, result2 error) {
	fake.doesPvtDataInfoExistMutex.Lock()
	defer fake.doesPvtDataInfoExistMutex.Unlock()
	fake.DoesPvtDataInfoExistStub = nil
	if fake.doesPvtDataInfoExistReturnsOnCall == nil {
		fake.doesPvtDataInfoExistReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.doesPvtDataInfoExistReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
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

func (fake *PeerLedger) GetBlocksIterator(arg1 uint64) (ledgera.ResultsIterator, error) {
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

func (fake *PeerLedger) GetBlocksIteratorCalls(stub func(uint64) (ledgera.ResultsIterator, error)) {
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

func (fake *PeerLedger) GetBlocksIteratorReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.getBlocksIteratorMutex.Lock()
	defer fake.getBlocksIteratorMutex.Unlock()
	fake.GetBlocksIteratorStub = nil
	fake.getBlocksIteratorReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetBlocksIteratorReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
	fake.getBlocksIteratorMutex.Lock()
	defer fake.getBlocksIteratorMutex.Unlock()
	fake.GetBlocksIteratorStub = nil
	if fake.getBlocksIteratorReturnsOnCall == nil {
		fake.getBlocksIteratorReturnsOnCall = make(map[int]struct {
			result1 ledgera.ResultsIterator
			result2 error
		})
	}
	fake.getBlocksIteratorReturnsOnCall[i] = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
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

func (fake *PeerLedger) GetConfigHistoryRetrieverCalls(stub func() (ledger.ConfigHistoryRetriever, error)) {
	fake.getConfigHistoryRetrieverMutex.Lock()
	defer fake.getConfigHistoryRetrieverMutex.Unlock()
	fake.GetConfigHistoryRetrieverStub = stub
}

func (fake *PeerLedger) GetConfigHistoryRetrieverReturns(result1 ledger.ConfigHistoryRetriever, result2 error) {
	fake.getConfigHistoryRetrieverMutex.Lock()
	defer fake.getConfigHistoryRetrieverMutex.Unlock()
	fake.GetConfigHistoryRetrieverStub = nil
	fake.getConfigHistoryRetrieverReturns = struct {
		result1 ledger.ConfigHistoryRetriever
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetConfigHistoryRetrieverReturnsOnCall(i int, result1 ledger.ConfigHistoryRetriever, result2 error) {
	fake.getConfigHistoryRetrieverMutex.Lock()
	defer fake.getConfigHistoryRetrieverMutex.Unlock()
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

func (fake *PeerLedger) GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error) {
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

func (fake *PeerLedger) GetMissingPvtDataTrackerCalls(stub func() (ledger.MissingPvtDataTracker, error)) {
	fake.getMissingPvtDataTrackerMutex.Lock()
	defer fake.getMissingPvtDataTrackerMutex.Unlock()
	fake.GetMissingPvtDataTrackerStub = stub
}

func (fake *PeerLedger) GetMissingPvtDataTrackerReturns(result1 ledger.MissingPvtDataTracker, result2 error) {
	fake.getMissingPvtDataTrackerMutex.Lock()
	defer fake.getMissingPvtDataTrackerMutex.Unlock()
	fake.GetMissingPvtDataTrackerStub = nil
	fake.getMissingPvtDataTrackerReturns = struct {
		result1 ledger.MissingPvtDataTracker
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetMissingPvtDataTrackerReturnsOnCall(i int, result1 ledger.MissingPvtDataTracker, result2 error) {
	fake.getMissingPvtDataTrackerMutex.Lock()
	defer fake.getMissingPvtDataTrackerMutex.Unlock()
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

func (fake *PeerLedger) GetPvtDataAndBlockByNum(arg1 uint64, arg2 ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	ret, specificReturn := fake.getPvtDataAndBlockByNumReturnsOnCall[len(fake.getPvtDataAndBlockByNumArgsForCall)]
	fake.getPvtDataAndBlockByNumArgsForCall = append(fake.getPvtDataAndBlockByNumArgsForCall, struct {
		arg1 uint64
		arg2 ledger.PvtNsCollFilter
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

func (fake *PeerLedger) GetPvtDataAndBlockByNumCalls(stub func(uint64, ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error)) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	defer fake.getPvtDataAndBlockByNumMutex.Unlock()
	fake.GetPvtDataAndBlockByNumStub = stub
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumArgsForCall(i int) (uint64, ledger.PvtNsCollFilter) {
	fake.getPvtDataAndBlockByNumMutex.RLock()
	defer fake.getPvtDataAndBlockByNumMutex.RUnlock()
	argsForCall := fake.getPvtDataAndBlockByNumArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumReturns(result1 *ledger.BlockAndPvtData, result2 error) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	defer fake.getPvtDataAndBlockByNumMutex.Unlock()
	fake.GetPvtDataAndBlockByNumStub = nil
	fake.getPvtDataAndBlockByNumReturns = struct {
		result1 *ledger.BlockAndPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataAndBlockByNumReturnsOnCall(i int, result1 *ledger.BlockAndPvtData, result2 error) {
	fake.getPvtDataAndBlockByNumMutex.Lock()
	defer fake.getPvtDataAndBlockByNumMutex.Unlock()
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

func (fake *PeerLedger) GetPvtDataByNum(arg1 uint64, arg2 ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	fake.getPvtDataByNumMutex.Lock()
	ret, specificReturn := fake.getPvtDataByNumReturnsOnCall[len(fake.getPvtDataByNumArgsForCall)]
	fake.getPvtDataByNumArgsForCall = append(fake.getPvtDataByNumArgsForCall, struct {
		arg1 uint64
		arg2 ledger.PvtNsCollFilter
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

func (fake *PeerLedger) GetPvtDataByNumCalls(stub func(uint64, ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)) {
	fake.getPvtDataByNumMutex.Lock()
	defer fake.getPvtDataByNumMutex.Unlock()
	fake.GetPvtDataByNumStub = stub
}

func (fake *PeerLedger) GetPvtDataByNumArgsForCall(i int) (uint64, ledger.PvtNsCollFilter) {
	fake.getPvtDataByNumMutex.RLock()
	defer fake.getPvtDataByNumMutex.RUnlock()
	argsForCall := fake.getPvtDataByNumArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PeerLedger) GetPvtDataByNumReturns(result1 []*ledger.TxPvtData, result2 error) {
	fake.getPvtDataByNumMutex.Lock()
	defer fake.getPvtDataByNumMutex.Unlock()
	fake.GetPvtDataByNumStub = nil
	fake.getPvtDataByNumReturns = struct {
		result1 []*ledger.TxPvtData
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetPvtDataByNumReturnsOnCall(i int, result1 []*ledger.TxPvtData, result2 error) {
	fake.getPvtDataByNumMutex.Lock()
	defer fake.getPvtDataByNumMutex.Unlock()
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

func (fake *PeerLedger) GetTransactionByID(arg1 string) (*peera.ProcessedTransaction, error) {
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

func (fake *PeerLedger) GetTransactionByIDCalls(stub func(string) (*peera.ProcessedTransaction, error)) {
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

func (fake *PeerLedger) GetTransactionByIDReturns(result1 *peera.ProcessedTransaction, result2 error) {
	fake.getTransactionByIDMutex.Lock()
	defer fake.getTransactionByIDMutex.Unlock()
	fake.GetTransactionByIDStub = nil
	fake.getTransactionByIDReturns = struct {
		result1 *peera.ProcessedTransaction
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTransactionByIDReturnsOnCall(i int, result1 *peera.ProcessedTransaction, result2 error) {
	fake.getTransactionByIDMutex.Lock()
	defer fake.getTransactionByIDMutex.Unlock()
	fake.GetTransactionByIDStub = nil
	if fake.getTransactionByIDReturnsOnCall == nil {
		fake.getTransactionByIDReturnsOnCall = make(map[int]struct {
			result1 *peera.ProcessedTransaction
			result2 error
		})
	}
	fake.getTransactionByIDReturnsOnCall[i] = struct {
		result1 *peera.ProcessedTransaction
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTxValidationCodeByTxID(arg1 string) (peera.TxValidationCode, error) {
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

func (fake *PeerLedger) GetTxValidationCodeByTxIDCalls(stub func(string) (peera.TxValidationCode, error)) {
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

func (fake *PeerLedger) GetTxValidationCodeByTxIDReturns(result1 peera.TxValidationCode, result2 error) {
	fake.getTxValidationCodeByTxIDMutex.Lock()
	defer fake.getTxValidationCodeByTxIDMutex.Unlock()
	fake.GetTxValidationCodeByTxIDStub = nil
	fake.getTxValidationCodeByTxIDReturns = struct {
		result1 peera.TxValidationCode
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) GetTxValidationCodeByTxIDReturnsOnCall(i int, result1 peera.TxValidationCode, result2 error) {
	fake.getTxValidationCodeByTxIDMutex.Lock()
	defer fake.getTxValidationCodeByTxIDMutex.Unlock()
	fake.GetTxValidationCodeByTxIDStub = nil
	if fake.getTxValidationCodeByTxIDReturnsOnCall == nil {
		fake.getTxValidationCodeByTxIDReturnsOnCall = make(map[int]struct {
			result1 peera.TxValidationCode
			result2 error
		})
	}
	fake.getTxValidationCodeByTxIDReturnsOnCall[i] = struct {
		result1 peera.TxValidationCode
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewHistoryQueryExecutor() (ledger.HistoryQueryExecutor, error) {
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

func (fake *PeerLedger) NewHistoryQueryExecutorCalls(stub func() (ledger.HistoryQueryExecutor, error)) {
	fake.newHistoryQueryExecutorMutex.Lock()
	defer fake.newHistoryQueryExecutorMutex.Unlock()
	fake.NewHistoryQueryExecutorStub = stub
}

func (fake *PeerLedger) NewHistoryQueryExecutorReturns(result1 ledger.HistoryQueryExecutor, result2 error) {
	fake.newHistoryQueryExecutorMutex.Lock()
	defer fake.newHistoryQueryExecutorMutex.Unlock()
	fake.NewHistoryQueryExecutorStub = nil
	fake.newHistoryQueryExecutorReturns = struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewHistoryQueryExecutorReturnsOnCall(i int, result1 ledger.HistoryQueryExecutor, result2 error) {
	fake.newHistoryQueryExecutorMutex.Lock()
	defer fake.newHistoryQueryExecutorMutex.Unlock()
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

func (fake *PeerLedger) NewQueryExecutor() (ledger.QueryExecutor, error) {
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

func (fake *PeerLedger) NewQueryExecutorCalls(stub func() (ledger.QueryExecutor, error)) {
	fake.newQueryExecutorMutex.Lock()
	defer fake.newQueryExecutorMutex.Unlock()
	fake.NewQueryExecutorStub = stub
}

func (fake *PeerLedger) NewQueryExecutorReturns(result1 ledger.QueryExecutor, result2 error) {
	fake.newQueryExecutorMutex.Lock()
	defer fake.newQueryExecutorMutex.Unlock()
	fake.NewQueryExecutorStub = nil
	fake.newQueryExecutorReturns = struct {
		result1 ledger.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewQueryExecutorReturnsOnCall(i int, result1 ledger.QueryExecutor, result2 error) {
	fake.newQueryExecutorMutex.Lock()
	defer fake.newQueryExecutorMutex.Unlock()
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

func (fake *PeerLedger) NewTxSimulator(arg1 string) (ledger.TxSimulator, error) {
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

func (fake *PeerLedger) NewTxSimulatorCalls(stub func(string) (ledger.TxSimulator, error)) {
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

func (fake *PeerLedger) NewTxSimulatorReturns(result1 ledger.TxSimulator, result2 error) {
	fake.newTxSimulatorMutex.Lock()
	defer fake.newTxSimulatorMutex.Unlock()
	fake.NewTxSimulatorStub = nil
	fake.newTxSimulatorReturns = struct {
		result1 ledger.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *PeerLedger) NewTxSimulatorReturnsOnCall(i int, result1 ledger.TxSimulator, result2 error) {
	fake.newTxSimulatorMutex.Lock()
	defer fake.newTxSimulatorMutex.Unlock()
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

func (fake *PeerLedger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.commitLegacyMutex.RLock()
	defer fake.commitLegacyMutex.RUnlock()
	fake.commitPvtDataOfOldBlocksMutex.RLock()
	defer fake.commitPvtDataOfOldBlocksMutex.RUnlock()
	fake.doesPvtDataInfoExistMutex.RLock()
	defer fake.doesPvtDataInfoExistMutex.RUnlock()
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
