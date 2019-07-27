
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
)

type TxMgr struct {
	NewQueryExecutorStub        func(txid string) (ledger.QueryExecutor, error)
	newQueryExecutorMutex       sync.RWMutex
	newQueryExecutorArgsForCall []struct {
		txid string
	}
	newQueryExecutorReturns struct {
		result1 ledger.QueryExecutor
		result2 error
	}
	newQueryExecutorReturnsOnCall map[int]struct {
		result1 ledger.QueryExecutor
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
	ValidateAndPrepareStub        func(blockAndPvtdata *ledger.BlockAndPvtData, doMVCCValidation bool) ([]*txmgr.TxStatInfo, error)
	validateAndPrepareMutex       sync.RWMutex
	validateAndPrepareArgsForCall []struct {
		blockAndPvtdata  *ledger.BlockAndPvtData
		doMVCCValidation bool
	}
	validateAndPrepareReturns struct {
		result1 []*txmgr.TxStatInfo
		result2 error
	}
	validateAndPrepareReturnsOnCall map[int]struct {
		result1 []*txmgr.TxStatInfo
		result2 error
	}
	RemoveStaleAndCommitPvtDataOfOldBlocksStub        func(blocksPvtData map[uint64][]*ledger.TxPvtData) error
	removeStaleAndCommitPvtDataOfOldBlocksMutex       sync.RWMutex
	removeStaleAndCommitPvtDataOfOldBlocksArgsForCall []struct {
		blocksPvtData map[uint64][]*ledger.TxPvtData
	}
	removeStaleAndCommitPvtDataOfOldBlocksReturns struct {
		result1 error
	}
	removeStaleAndCommitPvtDataOfOldBlocksReturnsOnCall map[int]struct {
		result1 error
	}
	GetLastSavepointStub        func() (*version.Height, error)
	getLastSavepointMutex       sync.RWMutex
	getLastSavepointArgsForCall []struct{}
	getLastSavepointReturns     struct {
		result1 *version.Height
		result2 error
	}
	getLastSavepointReturnsOnCall map[int]struct {
		result1 *version.Height
		result2 error
	}
	ShouldRecoverStub        func(lastAvailableBlock uint64) (bool, uint64, error)
	shouldRecoverMutex       sync.RWMutex
	shouldRecoverArgsForCall []struct {
		lastAvailableBlock uint64
	}
	shouldRecoverReturns struct {
		result1 bool
		result2 uint64
		result3 error
	}
	shouldRecoverReturnsOnCall map[int]struct {
		result1 bool
		result2 uint64
		result3 error
	}
	CommitLostBlockStub        func(blockAndPvtdata *ledger.BlockAndPvtData) error
	commitLostBlockMutex       sync.RWMutex
	commitLostBlockArgsForCall []struct {
		blockAndPvtdata *ledger.BlockAndPvtData
	}
	commitLostBlockReturns struct {
		result1 error
	}
	commitLostBlockReturnsOnCall map[int]struct {
		result1 error
	}
	CommitStub        func() error
	commitMutex       sync.RWMutex
	commitArgsForCall []struct{}
	commitReturns     struct {
		result1 error
	}
	commitReturnsOnCall map[int]struct {
		result1 error
	}
	RollbackStub        func()
	rollbackMutex       sync.RWMutex
	rollbackArgsForCall []struct{}
	ShutdownStub        func()
	shutdownMutex       sync.RWMutex
	shutdownArgsForCall []struct{}
	invocations         map[string][][]interface{}
	invocationsMutex    sync.RWMutex
}

func (fake *TxMgr) NewQueryExecutor(txid string) (ledger.QueryExecutor, error) {
	fake.newQueryExecutorMutex.Lock()
	ret, specificReturn := fake.newQueryExecutorReturnsOnCall[len(fake.newQueryExecutorArgsForCall)]
	fake.newQueryExecutorArgsForCall = append(fake.newQueryExecutorArgsForCall, struct {
		txid string
	}{txid})
	fake.recordInvocation("NewQueryExecutor", []interface{}{txid})
	fake.newQueryExecutorMutex.Unlock()
	if fake.NewQueryExecutorStub != nil {
		return fake.NewQueryExecutorStub(txid)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newQueryExecutorReturns.result1, fake.newQueryExecutorReturns.result2
}

func (fake *TxMgr) NewQueryExecutorCallCount() int {
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	return len(fake.newQueryExecutorArgsForCall)
}

func (fake *TxMgr) NewQueryExecutorArgsForCall(i int) string {
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	return fake.newQueryExecutorArgsForCall[i].txid
}

func (fake *TxMgr) NewQueryExecutorReturns(result1 ledger.QueryExecutor, result2 error) {
	fake.NewQueryExecutorStub = nil
	fake.newQueryExecutorReturns = struct {
		result1 ledger.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) NewQueryExecutorReturnsOnCall(i int, result1 ledger.QueryExecutor, result2 error) {
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

func (fake *TxMgr) NewTxSimulator(txid string) (ledger.TxSimulator, error) {
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

func (fake *TxMgr) NewTxSimulatorCallCount() int {
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	return len(fake.newTxSimulatorArgsForCall)
}

func (fake *TxMgr) NewTxSimulatorArgsForCall(i int) string {
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	return fake.newTxSimulatorArgsForCall[i].txid
}

func (fake *TxMgr) NewTxSimulatorReturns(result1 ledger.TxSimulator, result2 error) {
	fake.NewTxSimulatorStub = nil
	fake.newTxSimulatorReturns = struct {
		result1 ledger.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) NewTxSimulatorReturnsOnCall(i int, result1 ledger.TxSimulator, result2 error) {
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

func (fake *TxMgr) ValidateAndPrepare(blockAndPvtdata *ledger.BlockAndPvtData, doMVCCValidation bool) ([]*txmgr.TxStatInfo, error) {
	fake.validateAndPrepareMutex.Lock()
	ret, specificReturn := fake.validateAndPrepareReturnsOnCall[len(fake.validateAndPrepareArgsForCall)]
	fake.validateAndPrepareArgsForCall = append(fake.validateAndPrepareArgsForCall, struct {
		blockAndPvtdata  *ledger.BlockAndPvtData
		doMVCCValidation bool
	}{blockAndPvtdata, doMVCCValidation})
	fake.recordInvocation("ValidateAndPrepare", []interface{}{blockAndPvtdata, doMVCCValidation})
	fake.validateAndPrepareMutex.Unlock()
	if fake.ValidateAndPrepareStub != nil {
		return fake.ValidateAndPrepareStub(blockAndPvtdata, doMVCCValidation)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.validateAndPrepareReturns.result1, fake.validateAndPrepareReturns.result2
}

func (fake *TxMgr) ValidateAndPrepareCallCount() int {
	fake.validateAndPrepareMutex.RLock()
	defer fake.validateAndPrepareMutex.RUnlock()
	return len(fake.validateAndPrepareArgsForCall)
}

func (fake *TxMgr) ValidateAndPrepareArgsForCall(i int) (*ledger.BlockAndPvtData, bool) {
	fake.validateAndPrepareMutex.RLock()
	defer fake.validateAndPrepareMutex.RUnlock()
	return fake.validateAndPrepareArgsForCall[i].blockAndPvtdata, fake.validateAndPrepareArgsForCall[i].doMVCCValidation
}

func (fake *TxMgr) ValidateAndPrepareReturns(result1 []*txmgr.TxStatInfo, result2 error) {
	fake.ValidateAndPrepareStub = nil
	fake.validateAndPrepareReturns = struct {
		result1 []*txmgr.TxStatInfo
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) ValidateAndPrepareReturnsOnCall(i int, result1 []*txmgr.TxStatInfo, result2 error) {
	fake.ValidateAndPrepareStub = nil
	if fake.validateAndPrepareReturnsOnCall == nil {
		fake.validateAndPrepareReturnsOnCall = make(map[int]struct {
			result1 []*txmgr.TxStatInfo
			result2 error
		})
	}
	fake.validateAndPrepareReturnsOnCall[i] = struct {
		result1 []*txmgr.TxStatInfo
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocks(blocksPvtData map[uint64][]*ledger.TxPvtData) error {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Lock()
	ret, specificReturn := fake.removeStaleAndCommitPvtDataOfOldBlocksReturnsOnCall[len(fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall)]
	fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall = append(fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall, struct {
		blocksPvtData map[uint64][]*ledger.TxPvtData
	}{blocksPvtData})
	fake.recordInvocation("RemoveStaleAndCommitPvtDataOfOldBlocks", []interface{}{blocksPvtData})
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Unlock()
	if fake.RemoveStaleAndCommitPvtDataOfOldBlocksStub != nil {
		return fake.RemoveStaleAndCommitPvtDataOfOldBlocksStub(blocksPvtData)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.removeStaleAndCommitPvtDataOfOldBlocksReturns.result1
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksCallCount() int {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RLock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RUnlock()
	return len(fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall)
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksArgsForCall(i int) map[uint64][]*ledger.TxPvtData {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RLock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RUnlock()
	return fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall[i].blocksPvtData
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksReturns(result1 error) {
	fake.RemoveStaleAndCommitPvtDataOfOldBlocksStub = nil
	fake.removeStaleAndCommitPvtDataOfOldBlocksReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksReturnsOnCall(i int, result1 error) {
	fake.RemoveStaleAndCommitPvtDataOfOldBlocksStub = nil
	if fake.removeStaleAndCommitPvtDataOfOldBlocksReturnsOnCall == nil {
		fake.removeStaleAndCommitPvtDataOfOldBlocksReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeStaleAndCommitPvtDataOfOldBlocksReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) GetLastSavepoint() (*version.Height, error) {
	fake.getLastSavepointMutex.Lock()
	ret, specificReturn := fake.getLastSavepointReturnsOnCall[len(fake.getLastSavepointArgsForCall)]
	fake.getLastSavepointArgsForCall = append(fake.getLastSavepointArgsForCall, struct{}{})
	fake.recordInvocation("GetLastSavepoint", []interface{}{})
	fake.getLastSavepointMutex.Unlock()
	if fake.GetLastSavepointStub != nil {
		return fake.GetLastSavepointStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getLastSavepointReturns.result1, fake.getLastSavepointReturns.result2
}

func (fake *TxMgr) GetLastSavepointCallCount() int {
	fake.getLastSavepointMutex.RLock()
	defer fake.getLastSavepointMutex.RUnlock()
	return len(fake.getLastSavepointArgsForCall)
}

func (fake *TxMgr) GetLastSavepointReturns(result1 *version.Height, result2 error) {
	fake.GetLastSavepointStub = nil
	fake.getLastSavepointReturns = struct {
		result1 *version.Height
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) GetLastSavepointReturnsOnCall(i int, result1 *version.Height, result2 error) {
	fake.GetLastSavepointStub = nil
	if fake.getLastSavepointReturnsOnCall == nil {
		fake.getLastSavepointReturnsOnCall = make(map[int]struct {
			result1 *version.Height
			result2 error
		})
	}
	fake.getLastSavepointReturnsOnCall[i] = struct {
		result1 *version.Height
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) ShouldRecover(lastAvailableBlock uint64) (bool, uint64, error) {
	fake.shouldRecoverMutex.Lock()
	ret, specificReturn := fake.shouldRecoverReturnsOnCall[len(fake.shouldRecoverArgsForCall)]
	fake.shouldRecoverArgsForCall = append(fake.shouldRecoverArgsForCall, struct {
		lastAvailableBlock uint64
	}{lastAvailableBlock})
	fake.recordInvocation("ShouldRecover", []interface{}{lastAvailableBlock})
	fake.shouldRecoverMutex.Unlock()
	if fake.ShouldRecoverStub != nil {
		return fake.ShouldRecoverStub(lastAvailableBlock)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.shouldRecoverReturns.result1, fake.shouldRecoverReturns.result2, fake.shouldRecoverReturns.result3
}

func (fake *TxMgr) ShouldRecoverCallCount() int {
	fake.shouldRecoverMutex.RLock()
	defer fake.shouldRecoverMutex.RUnlock()
	return len(fake.shouldRecoverArgsForCall)
}

func (fake *TxMgr) ShouldRecoverArgsForCall(i int) uint64 {
	fake.shouldRecoverMutex.RLock()
	defer fake.shouldRecoverMutex.RUnlock()
	return fake.shouldRecoverArgsForCall[i].lastAvailableBlock
}

func (fake *TxMgr) ShouldRecoverReturns(result1 bool, result2 uint64, result3 error) {
	fake.ShouldRecoverStub = nil
	fake.shouldRecoverReturns = struct {
		result1 bool
		result2 uint64
		result3 error
	}{result1, result2, result3}
}

func (fake *TxMgr) ShouldRecoverReturnsOnCall(i int, result1 bool, result2 uint64, result3 error) {
	fake.ShouldRecoverStub = nil
	if fake.shouldRecoverReturnsOnCall == nil {
		fake.shouldRecoverReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 uint64
			result3 error
		})
	}
	fake.shouldRecoverReturnsOnCall[i] = struct {
		result1 bool
		result2 uint64
		result3 error
	}{result1, result2, result3}
}

func (fake *TxMgr) CommitLostBlock(blockAndPvtdata *ledger.BlockAndPvtData) error {
	fake.commitLostBlockMutex.Lock()
	ret, specificReturn := fake.commitLostBlockReturnsOnCall[len(fake.commitLostBlockArgsForCall)]
	fake.commitLostBlockArgsForCall = append(fake.commitLostBlockArgsForCall, struct {
		blockAndPvtdata *ledger.BlockAndPvtData
	}{blockAndPvtdata})
	fake.recordInvocation("CommitLostBlock", []interface{}{blockAndPvtdata})
	fake.commitLostBlockMutex.Unlock()
	if fake.CommitLostBlockStub != nil {
		return fake.CommitLostBlockStub(blockAndPvtdata)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.commitLostBlockReturns.result1
}

func (fake *TxMgr) CommitLostBlockCallCount() int {
	fake.commitLostBlockMutex.RLock()
	defer fake.commitLostBlockMutex.RUnlock()
	return len(fake.commitLostBlockArgsForCall)
}

func (fake *TxMgr) CommitLostBlockArgsForCall(i int) *ledger.BlockAndPvtData {
	fake.commitLostBlockMutex.RLock()
	defer fake.commitLostBlockMutex.RUnlock()
	return fake.commitLostBlockArgsForCall[i].blockAndPvtdata
}

func (fake *TxMgr) CommitLostBlockReturns(result1 error) {
	fake.CommitLostBlockStub = nil
	fake.commitLostBlockReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) CommitLostBlockReturnsOnCall(i int, result1 error) {
	fake.CommitLostBlockStub = nil
	if fake.commitLostBlockReturnsOnCall == nil {
		fake.commitLostBlockReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.commitLostBlockReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) Commit() error {
	fake.commitMutex.Lock()
	ret, specificReturn := fake.commitReturnsOnCall[len(fake.commitArgsForCall)]
	fake.commitArgsForCall = append(fake.commitArgsForCall, struct{}{})
	fake.recordInvocation("Commit", []interface{}{})
	fake.commitMutex.Unlock()
	if fake.CommitStub != nil {
		return fake.CommitStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.commitReturns.result1
}

func (fake *TxMgr) CommitCallCount() int {
	fake.commitMutex.RLock()
	defer fake.commitMutex.RUnlock()
	return len(fake.commitArgsForCall)
}

func (fake *TxMgr) CommitReturns(result1 error) {
	fake.CommitStub = nil
	fake.commitReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) CommitReturnsOnCall(i int, result1 error) {
	fake.CommitStub = nil
	if fake.commitReturnsOnCall == nil {
		fake.commitReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.commitReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) Rollback() {
	fake.rollbackMutex.Lock()
	fake.rollbackArgsForCall = append(fake.rollbackArgsForCall, struct{}{})
	fake.recordInvocation("Rollback", []interface{}{})
	fake.rollbackMutex.Unlock()
	if fake.RollbackStub != nil {
		fake.RollbackStub()
	}
}

func (fake *TxMgr) RollbackCallCount() int {
	fake.rollbackMutex.RLock()
	defer fake.rollbackMutex.RUnlock()
	return len(fake.rollbackArgsForCall)
}

func (fake *TxMgr) Shutdown() {
	fake.shutdownMutex.Lock()
	fake.shutdownArgsForCall = append(fake.shutdownArgsForCall, struct{}{})
	fake.recordInvocation("Shutdown", []interface{}{})
	fake.shutdownMutex.Unlock()
	if fake.ShutdownStub != nil {
		fake.ShutdownStub()
	}
}

func (fake *TxMgr) ShutdownCallCount() int {
	fake.shutdownMutex.RLock()
	defer fake.shutdownMutex.RUnlock()
	return len(fake.shutdownArgsForCall)
}

func (fake *TxMgr) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	fake.validateAndPrepareMutex.RLock()
	defer fake.validateAndPrepareMutex.RUnlock()
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RLock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RUnlock()
	fake.getLastSavepointMutex.RLock()
	defer fake.getLastSavepointMutex.RUnlock()
	fake.shouldRecoverMutex.RLock()
	defer fake.shouldRecoverMutex.RUnlock()
	fake.commitLostBlockMutex.RLock()
	defer fake.commitLostBlockMutex.RUnlock()
	fake.commitMutex.RLock()
	defer fake.commitMutex.RUnlock()
	fake.rollbackMutex.RLock()
	defer fake.rollbackMutex.RUnlock()
	fake.shutdownMutex.RLock()
	defer fake.shutdownMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TxMgr) recordInvocation(key string, args []interface{}) {
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
