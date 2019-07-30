
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
)

type TxMgr struct {
	CommitStub        func() error
	commitMutex       sync.RWMutex
	commitArgsForCall []struct {
	}
	commitReturns struct {
		result1 error
	}
	commitReturnsOnCall map[int]struct {
		result1 error
	}
	CommitLostBlockStub        func(*ledger.BlockAndPvtData) error
	commitLostBlockMutex       sync.RWMutex
	commitLostBlockArgsForCall []struct {
		arg1 *ledger.BlockAndPvtData
	}
	commitLostBlockReturns struct {
		result1 error
	}
	commitLostBlockReturnsOnCall map[int]struct {
		result1 error
	}
	GetLastSavepointStub        func() (*version.Height, error)
	getLastSavepointMutex       sync.RWMutex
	getLastSavepointArgsForCall []struct {
	}
	getLastSavepointReturns struct {
		result1 *version.Height
		result2 error
	}
	getLastSavepointReturnsOnCall map[int]struct {
		result1 *version.Height
		result2 error
	}
	NewQueryExecutorStub        func(string) (ledger.QueryExecutor, error)
	newQueryExecutorMutex       sync.RWMutex
	newQueryExecutorArgsForCall []struct {
		arg1 string
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
	RemoveStaleAndCommitPvtDataOfOldBlocksStub        func(map[uint64][]*ledger.TxPvtData) error
	removeStaleAndCommitPvtDataOfOldBlocksMutex       sync.RWMutex
	removeStaleAndCommitPvtDataOfOldBlocksArgsForCall []struct {
		arg1 map[uint64][]*ledger.TxPvtData
	}
	removeStaleAndCommitPvtDataOfOldBlocksReturns struct {
		result1 error
	}
	removeStaleAndCommitPvtDataOfOldBlocksReturnsOnCall map[int]struct {
		result1 error
	}
	RollbackStub        func()
	rollbackMutex       sync.RWMutex
	rollbackArgsForCall []struct {
	}
	ShouldRecoverStub        func(uint64) (bool, uint64, error)
	shouldRecoverMutex       sync.RWMutex
	shouldRecoverArgsForCall []struct {
		arg1 uint64
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
	ShutdownStub        func()
	shutdownMutex       sync.RWMutex
	shutdownArgsForCall []struct {
	}
	ValidateAndPrepareStub        func(*ledger.BlockAndPvtData, bool) ([]*txmgr.TxStatInfo, []byte, error)
	validateAndPrepareMutex       sync.RWMutex
	validateAndPrepareArgsForCall []struct {
		arg1 *ledger.BlockAndPvtData
		arg2 bool
	}
	validateAndPrepareReturns struct {
		result1 []*txmgr.TxStatInfo
		result2 []byte
		result3 error
	}
	validateAndPrepareReturnsOnCall map[int]struct {
		result1 []*txmgr.TxStatInfo
		result2 []byte
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TxMgr) Commit() error {
	fake.commitMutex.Lock()
	ret, specificReturn := fake.commitReturnsOnCall[len(fake.commitArgsForCall)]
	fake.commitArgsForCall = append(fake.commitArgsForCall, struct {
	}{})
	fake.recordInvocation("Commit", []interface{}{})
	fake.commitMutex.Unlock()
	if fake.CommitStub != nil {
		return fake.CommitStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.commitReturns
	return fakeReturns.result1
}

func (fake *TxMgr) CommitCallCount() int {
	fake.commitMutex.RLock()
	defer fake.commitMutex.RUnlock()
	return len(fake.commitArgsForCall)
}

func (fake *TxMgr) CommitCalls(stub func() error) {
	fake.commitMutex.Lock()
	defer fake.commitMutex.Unlock()
	fake.CommitStub = stub
}

func (fake *TxMgr) CommitReturns(result1 error) {
	fake.commitMutex.Lock()
	defer fake.commitMutex.Unlock()
	fake.CommitStub = nil
	fake.commitReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) CommitReturnsOnCall(i int, result1 error) {
	fake.commitMutex.Lock()
	defer fake.commitMutex.Unlock()
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

func (fake *TxMgr) CommitLostBlock(arg1 *ledger.BlockAndPvtData) error {
	fake.commitLostBlockMutex.Lock()
	ret, specificReturn := fake.commitLostBlockReturnsOnCall[len(fake.commitLostBlockArgsForCall)]
	fake.commitLostBlockArgsForCall = append(fake.commitLostBlockArgsForCall, struct {
		arg1 *ledger.BlockAndPvtData
	}{arg1})
	fake.recordInvocation("CommitLostBlock", []interface{}{arg1})
	fake.commitLostBlockMutex.Unlock()
	if fake.CommitLostBlockStub != nil {
		return fake.CommitLostBlockStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.commitLostBlockReturns
	return fakeReturns.result1
}

func (fake *TxMgr) CommitLostBlockCallCount() int {
	fake.commitLostBlockMutex.RLock()
	defer fake.commitLostBlockMutex.RUnlock()
	return len(fake.commitLostBlockArgsForCall)
}

func (fake *TxMgr) CommitLostBlockCalls(stub func(*ledger.BlockAndPvtData) error) {
	fake.commitLostBlockMutex.Lock()
	defer fake.commitLostBlockMutex.Unlock()
	fake.CommitLostBlockStub = stub
}

func (fake *TxMgr) CommitLostBlockArgsForCall(i int) *ledger.BlockAndPvtData {
	fake.commitLostBlockMutex.RLock()
	defer fake.commitLostBlockMutex.RUnlock()
	argsForCall := fake.commitLostBlockArgsForCall[i]
	return argsForCall.arg1
}

func (fake *TxMgr) CommitLostBlockReturns(result1 error) {
	fake.commitLostBlockMutex.Lock()
	defer fake.commitLostBlockMutex.Unlock()
	fake.CommitLostBlockStub = nil
	fake.commitLostBlockReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) CommitLostBlockReturnsOnCall(i int, result1 error) {
	fake.commitLostBlockMutex.Lock()
	defer fake.commitLostBlockMutex.Unlock()
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

func (fake *TxMgr) GetLastSavepoint() (*version.Height, error) {
	fake.getLastSavepointMutex.Lock()
	ret, specificReturn := fake.getLastSavepointReturnsOnCall[len(fake.getLastSavepointArgsForCall)]
	fake.getLastSavepointArgsForCall = append(fake.getLastSavepointArgsForCall, struct {
	}{})
	fake.recordInvocation("GetLastSavepoint", []interface{}{})
	fake.getLastSavepointMutex.Unlock()
	if fake.GetLastSavepointStub != nil {
		return fake.GetLastSavepointStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getLastSavepointReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxMgr) GetLastSavepointCallCount() int {
	fake.getLastSavepointMutex.RLock()
	defer fake.getLastSavepointMutex.RUnlock()
	return len(fake.getLastSavepointArgsForCall)
}

func (fake *TxMgr) GetLastSavepointCalls(stub func() (*version.Height, error)) {
	fake.getLastSavepointMutex.Lock()
	defer fake.getLastSavepointMutex.Unlock()
	fake.GetLastSavepointStub = stub
}

func (fake *TxMgr) GetLastSavepointReturns(result1 *version.Height, result2 error) {
	fake.getLastSavepointMutex.Lock()
	defer fake.getLastSavepointMutex.Unlock()
	fake.GetLastSavepointStub = nil
	fake.getLastSavepointReturns = struct {
		result1 *version.Height
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) GetLastSavepointReturnsOnCall(i int, result1 *version.Height, result2 error) {
	fake.getLastSavepointMutex.Lock()
	defer fake.getLastSavepointMutex.Unlock()
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

func (fake *TxMgr) NewQueryExecutor(arg1 string) (ledger.QueryExecutor, error) {
	fake.newQueryExecutorMutex.Lock()
	ret, specificReturn := fake.newQueryExecutorReturnsOnCall[len(fake.newQueryExecutorArgsForCall)]
	fake.newQueryExecutorArgsForCall = append(fake.newQueryExecutorArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("NewQueryExecutor", []interface{}{arg1})
	fake.newQueryExecutorMutex.Unlock()
	if fake.NewQueryExecutorStub != nil {
		return fake.NewQueryExecutorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newQueryExecutorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxMgr) NewQueryExecutorCallCount() int {
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	return len(fake.newQueryExecutorArgsForCall)
}

func (fake *TxMgr) NewQueryExecutorCalls(stub func(string) (ledger.QueryExecutor, error)) {
	fake.newQueryExecutorMutex.Lock()
	defer fake.newQueryExecutorMutex.Unlock()
	fake.NewQueryExecutorStub = stub
}

func (fake *TxMgr) NewQueryExecutorArgsForCall(i int) string {
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	argsForCall := fake.newQueryExecutorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *TxMgr) NewQueryExecutorReturns(result1 ledger.QueryExecutor, result2 error) {
	fake.newQueryExecutorMutex.Lock()
	defer fake.newQueryExecutorMutex.Unlock()
	fake.NewQueryExecutorStub = nil
	fake.newQueryExecutorReturns = struct {
		result1 ledger.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) NewQueryExecutorReturnsOnCall(i int, result1 ledger.QueryExecutor, result2 error) {
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

func (fake *TxMgr) NewTxSimulator(arg1 string) (ledger.TxSimulator, error) {
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

func (fake *TxMgr) NewTxSimulatorCallCount() int {
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	return len(fake.newTxSimulatorArgsForCall)
}

func (fake *TxMgr) NewTxSimulatorCalls(stub func(string) (ledger.TxSimulator, error)) {
	fake.newTxSimulatorMutex.Lock()
	defer fake.newTxSimulatorMutex.Unlock()
	fake.NewTxSimulatorStub = stub
}

func (fake *TxMgr) NewTxSimulatorArgsForCall(i int) string {
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	argsForCall := fake.newTxSimulatorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *TxMgr) NewTxSimulatorReturns(result1 ledger.TxSimulator, result2 error) {
	fake.newTxSimulatorMutex.Lock()
	defer fake.newTxSimulatorMutex.Unlock()
	fake.NewTxSimulatorStub = nil
	fake.newTxSimulatorReturns = struct {
		result1 ledger.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *TxMgr) NewTxSimulatorReturnsOnCall(i int, result1 ledger.TxSimulator, result2 error) {
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

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocks(arg1 map[uint64][]*ledger.TxPvtData) error {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Lock()
	ret, specificReturn := fake.removeStaleAndCommitPvtDataOfOldBlocksReturnsOnCall[len(fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall)]
	fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall = append(fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall, struct {
		arg1 map[uint64][]*ledger.TxPvtData
	}{arg1})
	fake.recordInvocation("RemoveStaleAndCommitPvtDataOfOldBlocks", []interface{}{arg1})
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Unlock()
	if fake.RemoveStaleAndCommitPvtDataOfOldBlocksStub != nil {
		return fake.RemoveStaleAndCommitPvtDataOfOldBlocksStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.removeStaleAndCommitPvtDataOfOldBlocksReturns
	return fakeReturns.result1
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksCallCount() int {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RLock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RUnlock()
	return len(fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall)
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksCalls(stub func(map[uint64][]*ledger.TxPvtData) error) {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Lock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Unlock()
	fake.RemoveStaleAndCommitPvtDataOfOldBlocksStub = stub
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksArgsForCall(i int) map[uint64][]*ledger.TxPvtData {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RLock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RUnlock()
	argsForCall := fake.removeStaleAndCommitPvtDataOfOldBlocksArgsForCall[i]
	return argsForCall.arg1
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksReturns(result1 error) {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Lock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Unlock()
	fake.RemoveStaleAndCommitPvtDataOfOldBlocksStub = nil
	fake.removeStaleAndCommitPvtDataOfOldBlocksReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxMgr) RemoveStaleAndCommitPvtDataOfOldBlocksReturnsOnCall(i int, result1 error) {
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Lock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.Unlock()
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

func (fake *TxMgr) Rollback() {
	fake.rollbackMutex.Lock()
	fake.rollbackArgsForCall = append(fake.rollbackArgsForCall, struct {
	}{})
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

func (fake *TxMgr) RollbackCalls(stub func()) {
	fake.rollbackMutex.Lock()
	defer fake.rollbackMutex.Unlock()
	fake.RollbackStub = stub
}

func (fake *TxMgr) ShouldRecover(arg1 uint64) (bool, uint64, error) {
	fake.shouldRecoverMutex.Lock()
	ret, specificReturn := fake.shouldRecoverReturnsOnCall[len(fake.shouldRecoverArgsForCall)]
	fake.shouldRecoverArgsForCall = append(fake.shouldRecoverArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("ShouldRecover", []interface{}{arg1})
	fake.shouldRecoverMutex.Unlock()
	if fake.ShouldRecoverStub != nil {
		return fake.ShouldRecoverStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.shouldRecoverReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *TxMgr) ShouldRecoverCallCount() int {
	fake.shouldRecoverMutex.RLock()
	defer fake.shouldRecoverMutex.RUnlock()
	return len(fake.shouldRecoverArgsForCall)
}

func (fake *TxMgr) ShouldRecoverCalls(stub func(uint64) (bool, uint64, error)) {
	fake.shouldRecoverMutex.Lock()
	defer fake.shouldRecoverMutex.Unlock()
	fake.ShouldRecoverStub = stub
}

func (fake *TxMgr) ShouldRecoverArgsForCall(i int) uint64 {
	fake.shouldRecoverMutex.RLock()
	defer fake.shouldRecoverMutex.RUnlock()
	argsForCall := fake.shouldRecoverArgsForCall[i]
	return argsForCall.arg1
}

func (fake *TxMgr) ShouldRecoverReturns(result1 bool, result2 uint64, result3 error) {
	fake.shouldRecoverMutex.Lock()
	defer fake.shouldRecoverMutex.Unlock()
	fake.ShouldRecoverStub = nil
	fake.shouldRecoverReturns = struct {
		result1 bool
		result2 uint64
		result3 error
	}{result1, result2, result3}
}

func (fake *TxMgr) ShouldRecoverReturnsOnCall(i int, result1 bool, result2 uint64, result3 error) {
	fake.shouldRecoverMutex.Lock()
	defer fake.shouldRecoverMutex.Unlock()
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

func (fake *TxMgr) Shutdown() {
	fake.shutdownMutex.Lock()
	fake.shutdownArgsForCall = append(fake.shutdownArgsForCall, struct {
	}{})
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

func (fake *TxMgr) ShutdownCalls(stub func()) {
	fake.shutdownMutex.Lock()
	defer fake.shutdownMutex.Unlock()
	fake.ShutdownStub = stub
}

func (fake *TxMgr) ValidateAndPrepare(arg1 *ledger.BlockAndPvtData, arg2 bool) ([]*txmgr.TxStatInfo, []byte, error) {
	fake.validateAndPrepareMutex.Lock()
	ret, specificReturn := fake.validateAndPrepareReturnsOnCall[len(fake.validateAndPrepareArgsForCall)]
	fake.validateAndPrepareArgsForCall = append(fake.validateAndPrepareArgsForCall, struct {
		arg1 *ledger.BlockAndPvtData
		arg2 bool
	}{arg1, arg2})
	fake.recordInvocation("ValidateAndPrepare", []interface{}{arg1, arg2})
	fake.validateAndPrepareMutex.Unlock()
	if fake.ValidateAndPrepareStub != nil {
		return fake.ValidateAndPrepareStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.validateAndPrepareReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *TxMgr) ValidateAndPrepareCallCount() int {
	fake.validateAndPrepareMutex.RLock()
	defer fake.validateAndPrepareMutex.RUnlock()
	return len(fake.validateAndPrepareArgsForCall)
}

func (fake *TxMgr) ValidateAndPrepareCalls(stub func(*ledger.BlockAndPvtData, bool) ([]*txmgr.TxStatInfo, []byte, error)) {
	fake.validateAndPrepareMutex.Lock()
	defer fake.validateAndPrepareMutex.Unlock()
	fake.ValidateAndPrepareStub = stub
}

func (fake *TxMgr) ValidateAndPrepareArgsForCall(i int) (*ledger.BlockAndPvtData, bool) {
	fake.validateAndPrepareMutex.RLock()
	defer fake.validateAndPrepareMutex.RUnlock()
	argsForCall := fake.validateAndPrepareArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *TxMgr) ValidateAndPrepareReturns(result1 []*txmgr.TxStatInfo, result2 []byte, result3 error) {
	fake.validateAndPrepareMutex.Lock()
	defer fake.validateAndPrepareMutex.Unlock()
	fake.ValidateAndPrepareStub = nil
	fake.validateAndPrepareReturns = struct {
		result1 []*txmgr.TxStatInfo
		result2 []byte
		result3 error
	}{result1, result2, result3}
}

func (fake *TxMgr) ValidateAndPrepareReturnsOnCall(i int, result1 []*txmgr.TxStatInfo, result2 []byte, result3 error) {
	fake.validateAndPrepareMutex.Lock()
	defer fake.validateAndPrepareMutex.Unlock()
	fake.ValidateAndPrepareStub = nil
	if fake.validateAndPrepareReturnsOnCall == nil {
		fake.validateAndPrepareReturnsOnCall = make(map[int]struct {
			result1 []*txmgr.TxStatInfo
			result2 []byte
			result3 error
		})
	}
	fake.validateAndPrepareReturnsOnCall[i] = struct {
		result1 []*txmgr.TxStatInfo
		result2 []byte
		result3 error
	}{result1, result2, result3}
}

func (fake *TxMgr) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.commitMutex.RLock()
	defer fake.commitMutex.RUnlock()
	fake.commitLostBlockMutex.RLock()
	defer fake.commitLostBlockMutex.RUnlock()
	fake.getLastSavepointMutex.RLock()
	defer fake.getLastSavepointMutex.RUnlock()
	fake.newQueryExecutorMutex.RLock()
	defer fake.newQueryExecutorMutex.RUnlock()
	fake.newTxSimulatorMutex.RLock()
	defer fake.newTxSimulatorMutex.RUnlock()
	fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RLock()
	defer fake.removeStaleAndCommitPvtDataOfOldBlocksMutex.RUnlock()
	fake.rollbackMutex.RLock()
	defer fake.rollbackMutex.RUnlock()
	fake.shouldRecoverMutex.RLock()
	defer fake.shouldRecoverMutex.RUnlock()
	fake.shutdownMutex.RLock()
	defer fake.shutdownMutex.RUnlock()
	fake.validateAndPrepareMutex.RLock()
	defer fake.validateAndPrepareMutex.RUnlock()
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
