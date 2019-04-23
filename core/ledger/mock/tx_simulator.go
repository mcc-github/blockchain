
package mock

import (
	"sync"

	ledgera "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger"
)

type TxSimulator struct {
	DeletePrivateDataStub        func(string, string, string) error
	deletePrivateDataMutex       sync.RWMutex
	deletePrivateDataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	deletePrivateDataReturns struct {
		result1 error
	}
	deletePrivateDataReturnsOnCall map[int]struct {
		result1 error
	}
	DeletePrivateDataMetadataStub        func(string, string, string) error
	deletePrivateDataMetadataMutex       sync.RWMutex
	deletePrivateDataMetadataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	deletePrivateDataMetadataReturns struct {
		result1 error
	}
	deletePrivateDataMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStateStub        func(string, string) error
	deleteStateMutex       sync.RWMutex
	deleteStateArgsForCall []struct {
		arg1 string
		arg2 string
	}
	deleteStateReturns struct {
		result1 error
	}
	deleteStateReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStateMetadataStub        func(string, string) error
	deleteStateMetadataMutex       sync.RWMutex
	deleteStateMetadataArgsForCall []struct {
		arg1 string
		arg2 string
	}
	deleteStateMetadataReturns struct {
		result1 error
	}
	deleteStateMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	DoneStub        func()
	doneMutex       sync.RWMutex
	doneArgsForCall []struct {
	}
	ExecuteQueryStub        func(string, string) (ledgera.ResultsIterator, error)
	executeQueryMutex       sync.RWMutex
	executeQueryArgsForCall []struct {
		arg1 string
		arg2 string
	}
	executeQueryReturns struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	executeQueryReturnsOnCall map[int]struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	ExecuteQueryOnPrivateDataStub        func(string, string, string) (ledgera.ResultsIterator, error)
	executeQueryOnPrivateDataMutex       sync.RWMutex
	executeQueryOnPrivateDataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	executeQueryOnPrivateDataReturns struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	executeQueryOnPrivateDataReturnsOnCall map[int]struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	ExecuteQueryWithMetadataStub        func(string, string, map[string]interface{}) (ledger.QueryResultsIterator, error)
	executeQueryWithMetadataMutex       sync.RWMutex
	executeQueryWithMetadataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 map[string]interface{}
	}
	executeQueryWithMetadataReturns struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}
	executeQueryWithMetadataReturnsOnCall map[int]struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}
	ExecuteUpdateStub        func(string) error
	executeUpdateMutex       sync.RWMutex
	executeUpdateArgsForCall []struct {
		arg1 string
	}
	executeUpdateReturns struct {
		result1 error
	}
	executeUpdateReturnsOnCall map[int]struct {
		result1 error
	}
	GetPrivateDataStub        func(string, string, string) ([]byte, error)
	getPrivateDataMutex       sync.RWMutex
	getPrivateDataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	getPrivateDataReturns struct {
		result1 []byte
		result2 error
	}
	getPrivateDataReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetPrivateDataHashStub        func(string, string, string) ([]byte, error)
	getPrivateDataHashMutex       sync.RWMutex
	getPrivateDataHashArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	getPrivateDataHashReturns struct {
		result1 []byte
		result2 error
	}
	getPrivateDataHashReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetPrivateDataMetadataStub        func(string, string, string) (map[string][]byte, error)
	getPrivateDataMetadataMutex       sync.RWMutex
	getPrivateDataMetadataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	getPrivateDataMetadataReturns struct {
		result1 map[string][]byte
		result2 error
	}
	getPrivateDataMetadataReturnsOnCall map[int]struct {
		result1 map[string][]byte
		result2 error
	}
	GetPrivateDataMetadataByHashStub        func(string, string, []byte) (map[string][]byte, error)
	getPrivateDataMetadataByHashMutex       sync.RWMutex
	getPrivateDataMetadataByHashArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 []byte
	}
	getPrivateDataMetadataByHashReturns struct {
		result1 map[string][]byte
		result2 error
	}
	getPrivateDataMetadataByHashReturnsOnCall map[int]struct {
		result1 map[string][]byte
		result2 error
	}
	GetPrivateDataMultipleKeysStub        func(string, string, []string) ([][]byte, error)
	getPrivateDataMultipleKeysMutex       sync.RWMutex
	getPrivateDataMultipleKeysArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 []string
	}
	getPrivateDataMultipleKeysReturns struct {
		result1 [][]byte
		result2 error
	}
	getPrivateDataMultipleKeysReturnsOnCall map[int]struct {
		result1 [][]byte
		result2 error
	}
	GetPrivateDataRangeScanIteratorStub        func(string, string, string, string) (ledgera.ResultsIterator, error)
	getPrivateDataRangeScanIteratorMutex       sync.RWMutex
	getPrivateDataRangeScanIteratorArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 string
	}
	getPrivateDataRangeScanIteratorReturns struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	getPrivateDataRangeScanIteratorReturnsOnCall map[int]struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	GetStateStub        func(string, string) ([]byte, error)
	getStateMutex       sync.RWMutex
	getStateArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getStateReturns struct {
		result1 []byte
		result2 error
	}
	getStateReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetStateMetadataStub        func(string, string) (map[string][]byte, error)
	getStateMetadataMutex       sync.RWMutex
	getStateMetadataArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getStateMetadataReturns struct {
		result1 map[string][]byte
		result2 error
	}
	getStateMetadataReturnsOnCall map[int]struct {
		result1 map[string][]byte
		result2 error
	}
	GetStateMultipleKeysStub        func(string, []string) ([][]byte, error)
	getStateMultipleKeysMutex       sync.RWMutex
	getStateMultipleKeysArgsForCall []struct {
		arg1 string
		arg2 []string
	}
	getStateMultipleKeysReturns struct {
		result1 [][]byte
		result2 error
	}
	getStateMultipleKeysReturnsOnCall map[int]struct {
		result1 [][]byte
		result2 error
	}
	GetStateRangeScanIteratorStub        func(string, string, string) (ledgera.ResultsIterator, error)
	getStateRangeScanIteratorMutex       sync.RWMutex
	getStateRangeScanIteratorArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	getStateRangeScanIteratorReturns struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	getStateRangeScanIteratorReturnsOnCall map[int]struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	GetStateRangeScanIteratorWithMetadataStub        func(string, string, string, map[string]interface{}) (ledger.QueryResultsIterator, error)
	getStateRangeScanIteratorWithMetadataMutex       sync.RWMutex
	getStateRangeScanIteratorWithMetadataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 map[string]interface{}
	}
	getStateRangeScanIteratorWithMetadataReturns struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}
	getStateRangeScanIteratorWithMetadataReturnsOnCall map[int]struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}
	GetTxSimulationResultsStub        func() (*ledger.TxSimulationResults, error)
	getTxSimulationResultsMutex       sync.RWMutex
	getTxSimulationResultsArgsForCall []struct {
	}
	getTxSimulationResultsReturns struct {
		result1 *ledger.TxSimulationResults
		result2 error
	}
	getTxSimulationResultsReturnsOnCall map[int]struct {
		result1 *ledger.TxSimulationResults
		result2 error
	}
	SetPrivateDataStub        func(string, string, string, []byte) error
	setPrivateDataMutex       sync.RWMutex
	setPrivateDataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 []byte
	}
	setPrivateDataReturns struct {
		result1 error
	}
	setPrivateDataReturnsOnCall map[int]struct {
		result1 error
	}
	SetPrivateDataMetadataStub        func(string, string, string, map[string][]byte) error
	setPrivateDataMetadataMutex       sync.RWMutex
	setPrivateDataMetadataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 map[string][]byte
	}
	setPrivateDataMetadataReturns struct {
		result1 error
	}
	setPrivateDataMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	SetPrivateDataMultipleKeysStub        func(string, string, map[string][]byte) error
	setPrivateDataMultipleKeysMutex       sync.RWMutex
	setPrivateDataMultipleKeysArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 map[string][]byte
	}
	setPrivateDataMultipleKeysReturns struct {
		result1 error
	}
	setPrivateDataMultipleKeysReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateStub        func(string, string, []byte) error
	setStateMutex       sync.RWMutex
	setStateArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 []byte
	}
	setStateReturns struct {
		result1 error
	}
	setStateReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateMetadataStub        func(string, string, map[string][]byte) error
	setStateMetadataMutex       sync.RWMutex
	setStateMetadataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 map[string][]byte
	}
	setStateMetadataReturns struct {
		result1 error
	}
	setStateMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateMultipleKeysStub        func(string, map[string][]byte) error
	setStateMultipleKeysMutex       sync.RWMutex
	setStateMultipleKeysArgsForCall []struct {
		arg1 string
		arg2 map[string][]byte
	}
	setStateMultipleKeysReturns struct {
		result1 error
	}
	setStateMultipleKeysReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TxSimulator) DeletePrivateData(arg1 string, arg2 string, arg3 string) error {
	fake.deletePrivateDataMutex.Lock()
	ret, specificReturn := fake.deletePrivateDataReturnsOnCall[len(fake.deletePrivateDataArgsForCall)]
	fake.deletePrivateDataArgsForCall = append(fake.deletePrivateDataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("DeletePrivateData", []interface{}{arg1, arg2, arg3})
	fake.deletePrivateDataMutex.Unlock()
	if fake.DeletePrivateDataStub != nil {
		return fake.DeletePrivateDataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deletePrivateDataReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) DeletePrivateDataCallCount() int {
	fake.deletePrivateDataMutex.RLock()
	defer fake.deletePrivateDataMutex.RUnlock()
	return len(fake.deletePrivateDataArgsForCall)
}

func (fake *TxSimulator) DeletePrivateDataCalls(stub func(string, string, string) error) {
	fake.deletePrivateDataMutex.Lock()
	defer fake.deletePrivateDataMutex.Unlock()
	fake.DeletePrivateDataStub = stub
}

func (fake *TxSimulator) DeletePrivateDataArgsForCall(i int) (string, string, string) {
	fake.deletePrivateDataMutex.RLock()
	defer fake.deletePrivateDataMutex.RUnlock()
	argsForCall := fake.deletePrivateDataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) DeletePrivateDataReturns(result1 error) {
	fake.deletePrivateDataMutex.Lock()
	defer fake.deletePrivateDataMutex.Unlock()
	fake.DeletePrivateDataStub = nil
	fake.deletePrivateDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeletePrivateDataReturnsOnCall(i int, result1 error) {
	fake.deletePrivateDataMutex.Lock()
	defer fake.deletePrivateDataMutex.Unlock()
	fake.DeletePrivateDataStub = nil
	if fake.deletePrivateDataReturnsOnCall == nil {
		fake.deletePrivateDataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deletePrivateDataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeletePrivateDataMetadata(arg1 string, arg2 string, arg3 string) error {
	fake.deletePrivateDataMetadataMutex.Lock()
	ret, specificReturn := fake.deletePrivateDataMetadataReturnsOnCall[len(fake.deletePrivateDataMetadataArgsForCall)]
	fake.deletePrivateDataMetadataArgsForCall = append(fake.deletePrivateDataMetadataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("DeletePrivateDataMetadata", []interface{}{arg1, arg2, arg3})
	fake.deletePrivateDataMetadataMutex.Unlock()
	if fake.DeletePrivateDataMetadataStub != nil {
		return fake.DeletePrivateDataMetadataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deletePrivateDataMetadataReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) DeletePrivateDataMetadataCallCount() int {
	fake.deletePrivateDataMetadataMutex.RLock()
	defer fake.deletePrivateDataMetadataMutex.RUnlock()
	return len(fake.deletePrivateDataMetadataArgsForCall)
}

func (fake *TxSimulator) DeletePrivateDataMetadataCalls(stub func(string, string, string) error) {
	fake.deletePrivateDataMetadataMutex.Lock()
	defer fake.deletePrivateDataMetadataMutex.Unlock()
	fake.DeletePrivateDataMetadataStub = stub
}

func (fake *TxSimulator) DeletePrivateDataMetadataArgsForCall(i int) (string, string, string) {
	fake.deletePrivateDataMetadataMutex.RLock()
	defer fake.deletePrivateDataMetadataMutex.RUnlock()
	argsForCall := fake.deletePrivateDataMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) DeletePrivateDataMetadataReturns(result1 error) {
	fake.deletePrivateDataMetadataMutex.Lock()
	defer fake.deletePrivateDataMetadataMutex.Unlock()
	fake.DeletePrivateDataMetadataStub = nil
	fake.deletePrivateDataMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeletePrivateDataMetadataReturnsOnCall(i int, result1 error) {
	fake.deletePrivateDataMetadataMutex.Lock()
	defer fake.deletePrivateDataMetadataMutex.Unlock()
	fake.DeletePrivateDataMetadataStub = nil
	if fake.deletePrivateDataMetadataReturnsOnCall == nil {
		fake.deletePrivateDataMetadataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deletePrivateDataMetadataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeleteState(arg1 string, arg2 string) error {
	fake.deleteStateMutex.Lock()
	ret, specificReturn := fake.deleteStateReturnsOnCall[len(fake.deleteStateArgsForCall)]
	fake.deleteStateArgsForCall = append(fake.deleteStateArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("DeleteState", []interface{}{arg1, arg2})
	fake.deleteStateMutex.Unlock()
	if fake.DeleteStateStub != nil {
		return fake.DeleteStateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteStateReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) DeleteStateCallCount() int {
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	return len(fake.deleteStateArgsForCall)
}

func (fake *TxSimulator) DeleteStateCalls(stub func(string, string) error) {
	fake.deleteStateMutex.Lock()
	defer fake.deleteStateMutex.Unlock()
	fake.DeleteStateStub = stub
}

func (fake *TxSimulator) DeleteStateArgsForCall(i int) (string, string) {
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	argsForCall := fake.deleteStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *TxSimulator) DeleteStateReturns(result1 error) {
	fake.deleteStateMutex.Lock()
	defer fake.deleteStateMutex.Unlock()
	fake.DeleteStateStub = nil
	fake.deleteStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeleteStateReturnsOnCall(i int, result1 error) {
	fake.deleteStateMutex.Lock()
	defer fake.deleteStateMutex.Unlock()
	fake.DeleteStateStub = nil
	if fake.deleteStateReturnsOnCall == nil {
		fake.deleteStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeleteStateMetadata(arg1 string, arg2 string) error {
	fake.deleteStateMetadataMutex.Lock()
	ret, specificReturn := fake.deleteStateMetadataReturnsOnCall[len(fake.deleteStateMetadataArgsForCall)]
	fake.deleteStateMetadataArgsForCall = append(fake.deleteStateMetadataArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("DeleteStateMetadata", []interface{}{arg1, arg2})
	fake.deleteStateMetadataMutex.Unlock()
	if fake.DeleteStateMetadataStub != nil {
		return fake.DeleteStateMetadataStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteStateMetadataReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) DeleteStateMetadataCallCount() int {
	fake.deleteStateMetadataMutex.RLock()
	defer fake.deleteStateMetadataMutex.RUnlock()
	return len(fake.deleteStateMetadataArgsForCall)
}

func (fake *TxSimulator) DeleteStateMetadataCalls(stub func(string, string) error) {
	fake.deleteStateMetadataMutex.Lock()
	defer fake.deleteStateMetadataMutex.Unlock()
	fake.DeleteStateMetadataStub = stub
}

func (fake *TxSimulator) DeleteStateMetadataArgsForCall(i int) (string, string) {
	fake.deleteStateMetadataMutex.RLock()
	defer fake.deleteStateMetadataMutex.RUnlock()
	argsForCall := fake.deleteStateMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *TxSimulator) DeleteStateMetadataReturns(result1 error) {
	fake.deleteStateMetadataMutex.Lock()
	defer fake.deleteStateMetadataMutex.Unlock()
	fake.DeleteStateMetadataStub = nil
	fake.deleteStateMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeleteStateMetadataReturnsOnCall(i int, result1 error) {
	fake.deleteStateMetadataMutex.Lock()
	defer fake.deleteStateMetadataMutex.Unlock()
	fake.DeleteStateMetadataStub = nil
	if fake.deleteStateMetadataReturnsOnCall == nil {
		fake.deleteStateMetadataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteStateMetadataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) Done() {
	fake.doneMutex.Lock()
	fake.doneArgsForCall = append(fake.doneArgsForCall, struct {
	}{})
	fake.recordInvocation("Done", []interface{}{})
	fake.doneMutex.Unlock()
	if fake.DoneStub != nil {
		fake.DoneStub()
	}
}

func (fake *TxSimulator) DoneCallCount() int {
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
	return len(fake.doneArgsForCall)
}

func (fake *TxSimulator) DoneCalls(stub func()) {
	fake.doneMutex.Lock()
	defer fake.doneMutex.Unlock()
	fake.DoneStub = stub
}

func (fake *TxSimulator) ExecuteQuery(arg1 string, arg2 string) (ledgera.ResultsIterator, error) {
	fake.executeQueryMutex.Lock()
	ret, specificReturn := fake.executeQueryReturnsOnCall[len(fake.executeQueryArgsForCall)]
	fake.executeQueryArgsForCall = append(fake.executeQueryArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ExecuteQuery", []interface{}{arg1, arg2})
	fake.executeQueryMutex.Unlock()
	if fake.ExecuteQueryStub != nil {
		return fake.ExecuteQueryStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.executeQueryReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) ExecuteQueryCallCount() int {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	return len(fake.executeQueryArgsForCall)
}

func (fake *TxSimulator) ExecuteQueryCalls(stub func(string, string) (ledgera.ResultsIterator, error)) {
	fake.executeQueryMutex.Lock()
	defer fake.executeQueryMutex.Unlock()
	fake.ExecuteQueryStub = stub
}

func (fake *TxSimulator) ExecuteQueryArgsForCall(i int) (string, string) {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	argsForCall := fake.executeQueryArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *TxSimulator) ExecuteQueryReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.executeQueryMutex.Lock()
	defer fake.executeQueryMutex.Unlock()
	fake.ExecuteQueryStub = nil
	fake.executeQueryReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteQueryReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
	fake.executeQueryMutex.Lock()
	defer fake.executeQueryMutex.Unlock()
	fake.ExecuteQueryStub = nil
	if fake.executeQueryReturnsOnCall == nil {
		fake.executeQueryReturnsOnCall = make(map[int]struct {
			result1 ledgera.ResultsIterator
			result2 error
		})
	}
	fake.executeQueryReturnsOnCall[i] = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteQueryOnPrivateData(arg1 string, arg2 string, arg3 string) (ledgera.ResultsIterator, error) {
	fake.executeQueryOnPrivateDataMutex.Lock()
	ret, specificReturn := fake.executeQueryOnPrivateDataReturnsOnCall[len(fake.executeQueryOnPrivateDataArgsForCall)]
	fake.executeQueryOnPrivateDataArgsForCall = append(fake.executeQueryOnPrivateDataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("ExecuteQueryOnPrivateData", []interface{}{arg1, arg2, arg3})
	fake.executeQueryOnPrivateDataMutex.Unlock()
	if fake.ExecuteQueryOnPrivateDataStub != nil {
		return fake.ExecuteQueryOnPrivateDataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.executeQueryOnPrivateDataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) ExecuteQueryOnPrivateDataCallCount() int {
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	return len(fake.executeQueryOnPrivateDataArgsForCall)
}

func (fake *TxSimulator) ExecuteQueryOnPrivateDataCalls(stub func(string, string, string) (ledgera.ResultsIterator, error)) {
	fake.executeQueryOnPrivateDataMutex.Lock()
	defer fake.executeQueryOnPrivateDataMutex.Unlock()
	fake.ExecuteQueryOnPrivateDataStub = stub
}

func (fake *TxSimulator) ExecuteQueryOnPrivateDataArgsForCall(i int) (string, string, string) {
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	argsForCall := fake.executeQueryOnPrivateDataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) ExecuteQueryOnPrivateDataReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.executeQueryOnPrivateDataMutex.Lock()
	defer fake.executeQueryOnPrivateDataMutex.Unlock()
	fake.ExecuteQueryOnPrivateDataStub = nil
	fake.executeQueryOnPrivateDataReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteQueryOnPrivateDataReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
	fake.executeQueryOnPrivateDataMutex.Lock()
	defer fake.executeQueryOnPrivateDataMutex.Unlock()
	fake.ExecuteQueryOnPrivateDataStub = nil
	if fake.executeQueryOnPrivateDataReturnsOnCall == nil {
		fake.executeQueryOnPrivateDataReturnsOnCall = make(map[int]struct {
			result1 ledgera.ResultsIterator
			result2 error
		})
	}
	fake.executeQueryOnPrivateDataReturnsOnCall[i] = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteQueryWithMetadata(arg1 string, arg2 string, arg3 map[string]interface{}) (ledger.QueryResultsIterator, error) {
	fake.executeQueryWithMetadataMutex.Lock()
	ret, specificReturn := fake.executeQueryWithMetadataReturnsOnCall[len(fake.executeQueryWithMetadataArgsForCall)]
	fake.executeQueryWithMetadataArgsForCall = append(fake.executeQueryWithMetadataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 map[string]interface{}
	}{arg1, arg2, arg3})
	fake.recordInvocation("ExecuteQueryWithMetadata", []interface{}{arg1, arg2, arg3})
	fake.executeQueryWithMetadataMutex.Unlock()
	if fake.ExecuteQueryWithMetadataStub != nil {
		return fake.ExecuteQueryWithMetadataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.executeQueryWithMetadataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) ExecuteQueryWithMetadataCallCount() int {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	return len(fake.executeQueryWithMetadataArgsForCall)
}

func (fake *TxSimulator) ExecuteQueryWithMetadataCalls(stub func(string, string, map[string]interface{}) (ledger.QueryResultsIterator, error)) {
	fake.executeQueryWithMetadataMutex.Lock()
	defer fake.executeQueryWithMetadataMutex.Unlock()
	fake.ExecuteQueryWithMetadataStub = stub
}

func (fake *TxSimulator) ExecuteQueryWithMetadataArgsForCall(i int) (string, string, map[string]interface{}) {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	argsForCall := fake.executeQueryWithMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) ExecuteQueryWithMetadataReturns(result1 ledger.QueryResultsIterator, result2 error) {
	fake.executeQueryWithMetadataMutex.Lock()
	defer fake.executeQueryWithMetadataMutex.Unlock()
	fake.ExecuteQueryWithMetadataStub = nil
	fake.executeQueryWithMetadataReturns = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteQueryWithMetadataReturnsOnCall(i int, result1 ledger.QueryResultsIterator, result2 error) {
	fake.executeQueryWithMetadataMutex.Lock()
	defer fake.executeQueryWithMetadataMutex.Unlock()
	fake.ExecuteQueryWithMetadataStub = nil
	if fake.executeQueryWithMetadataReturnsOnCall == nil {
		fake.executeQueryWithMetadataReturnsOnCall = make(map[int]struct {
			result1 ledger.QueryResultsIterator
			result2 error
		})
	}
	fake.executeQueryWithMetadataReturnsOnCall[i] = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteUpdate(arg1 string) error {
	fake.executeUpdateMutex.Lock()
	ret, specificReturn := fake.executeUpdateReturnsOnCall[len(fake.executeUpdateArgsForCall)]
	fake.executeUpdateArgsForCall = append(fake.executeUpdateArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ExecuteUpdate", []interface{}{arg1})
	fake.executeUpdateMutex.Unlock()
	if fake.ExecuteUpdateStub != nil {
		return fake.ExecuteUpdateStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.executeUpdateReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) ExecuteUpdateCallCount() int {
	fake.executeUpdateMutex.RLock()
	defer fake.executeUpdateMutex.RUnlock()
	return len(fake.executeUpdateArgsForCall)
}

func (fake *TxSimulator) ExecuteUpdateCalls(stub func(string) error) {
	fake.executeUpdateMutex.Lock()
	defer fake.executeUpdateMutex.Unlock()
	fake.ExecuteUpdateStub = stub
}

func (fake *TxSimulator) ExecuteUpdateArgsForCall(i int) string {
	fake.executeUpdateMutex.RLock()
	defer fake.executeUpdateMutex.RUnlock()
	argsForCall := fake.executeUpdateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *TxSimulator) ExecuteUpdateReturns(result1 error) {
	fake.executeUpdateMutex.Lock()
	defer fake.executeUpdateMutex.Unlock()
	fake.ExecuteUpdateStub = nil
	fake.executeUpdateReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) ExecuteUpdateReturnsOnCall(i int, result1 error) {
	fake.executeUpdateMutex.Lock()
	defer fake.executeUpdateMutex.Unlock()
	fake.ExecuteUpdateStub = nil
	if fake.executeUpdateReturnsOnCall == nil {
		fake.executeUpdateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.executeUpdateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) GetPrivateData(arg1 string, arg2 string, arg3 string) ([]byte, error) {
	fake.getPrivateDataMutex.Lock()
	ret, specificReturn := fake.getPrivateDataReturnsOnCall[len(fake.getPrivateDataArgsForCall)]
	fake.getPrivateDataArgsForCall = append(fake.getPrivateDataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetPrivateData", []interface{}{arg1, arg2, arg3})
	fake.getPrivateDataMutex.Unlock()
	if fake.GetPrivateDataStub != nil {
		return fake.GetPrivateDataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPrivateDataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetPrivateDataCallCount() int {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	return len(fake.getPrivateDataArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataCalls(stub func(string, string, string) ([]byte, error)) {
	fake.getPrivateDataMutex.Lock()
	defer fake.getPrivateDataMutex.Unlock()
	fake.GetPrivateDataStub = stub
}

func (fake *TxSimulator) GetPrivateDataArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	argsForCall := fake.getPrivateDataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) GetPrivateDataReturns(result1 []byte, result2 error) {
	fake.getPrivateDataMutex.Lock()
	defer fake.getPrivateDataMutex.Unlock()
	fake.GetPrivateDataStub = nil
	fake.getPrivateDataReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getPrivateDataMutex.Lock()
	defer fake.getPrivateDataMutex.Unlock()
	fake.GetPrivateDataStub = nil
	if fake.getPrivateDataReturnsOnCall == nil {
		fake.getPrivateDataReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getPrivateDataReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataHash(arg1 string, arg2 string, arg3 string) ([]byte, error) {
	fake.getPrivateDataHashMutex.Lock()
	ret, specificReturn := fake.getPrivateDataHashReturnsOnCall[len(fake.getPrivateDataHashArgsForCall)]
	fake.getPrivateDataHashArgsForCall = append(fake.getPrivateDataHashArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetPrivateDataHash", []interface{}{arg1, arg2, arg3})
	fake.getPrivateDataHashMutex.Unlock()
	if fake.GetPrivateDataHashStub != nil {
		return fake.GetPrivateDataHashStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPrivateDataHashReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetPrivateDataHashCallCount() int {
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	return len(fake.getPrivateDataHashArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataHashCalls(stub func(string, string, string) ([]byte, error)) {
	fake.getPrivateDataHashMutex.Lock()
	defer fake.getPrivateDataHashMutex.Unlock()
	fake.GetPrivateDataHashStub = stub
}

func (fake *TxSimulator) GetPrivateDataHashArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	argsForCall := fake.getPrivateDataHashArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) GetPrivateDataHashReturns(result1 []byte, result2 error) {
	fake.getPrivateDataHashMutex.Lock()
	defer fake.getPrivateDataHashMutex.Unlock()
	fake.GetPrivateDataHashStub = nil
	fake.getPrivateDataHashReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataHashReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getPrivateDataHashMutex.Lock()
	defer fake.getPrivateDataHashMutex.Unlock()
	fake.GetPrivateDataHashStub = nil
	if fake.getPrivateDataHashReturnsOnCall == nil {
		fake.getPrivateDataHashReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getPrivateDataHashReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMetadata(arg1 string, arg2 string, arg3 string) (map[string][]byte, error) {
	fake.getPrivateDataMetadataMutex.Lock()
	ret, specificReturn := fake.getPrivateDataMetadataReturnsOnCall[len(fake.getPrivateDataMetadataArgsForCall)]
	fake.getPrivateDataMetadataArgsForCall = append(fake.getPrivateDataMetadataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetPrivateDataMetadata", []interface{}{arg1, arg2, arg3})
	fake.getPrivateDataMetadataMutex.Unlock()
	if fake.GetPrivateDataMetadataStub != nil {
		return fake.GetPrivateDataMetadataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPrivateDataMetadataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetPrivateDataMetadataCallCount() int {
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	return len(fake.getPrivateDataMetadataArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataMetadataCalls(stub func(string, string, string) (map[string][]byte, error)) {
	fake.getPrivateDataMetadataMutex.Lock()
	defer fake.getPrivateDataMetadataMutex.Unlock()
	fake.GetPrivateDataMetadataStub = stub
}

func (fake *TxSimulator) GetPrivateDataMetadataArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	argsForCall := fake.getPrivateDataMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) GetPrivateDataMetadataReturns(result1 map[string][]byte, result2 error) {
	fake.getPrivateDataMetadataMutex.Lock()
	defer fake.getPrivateDataMetadataMutex.Unlock()
	fake.GetPrivateDataMetadataStub = nil
	fake.getPrivateDataMetadataReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMetadataReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
	fake.getPrivateDataMetadataMutex.Lock()
	defer fake.getPrivateDataMetadataMutex.Unlock()
	fake.GetPrivateDataMetadataStub = nil
	if fake.getPrivateDataMetadataReturnsOnCall == nil {
		fake.getPrivateDataMetadataReturnsOnCall = make(map[int]struct {
			result1 map[string][]byte
			result2 error
		})
	}
	fake.getPrivateDataMetadataReturnsOnCall[i] = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMetadataByHash(arg1 string, arg2 string, arg3 []byte) (map[string][]byte, error) {
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.getPrivateDataMetadataByHashMutex.Lock()
	ret, specificReturn := fake.getPrivateDataMetadataByHashReturnsOnCall[len(fake.getPrivateDataMetadataByHashArgsForCall)]
	fake.getPrivateDataMetadataByHashArgsForCall = append(fake.getPrivateDataMetadataByHashArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 []byte
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("GetPrivateDataMetadataByHash", []interface{}{arg1, arg2, arg3Copy})
	fake.getPrivateDataMetadataByHashMutex.Unlock()
	if fake.GetPrivateDataMetadataByHashStub != nil {
		return fake.GetPrivateDataMetadataByHashStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPrivateDataMetadataByHashReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashCallCount() int {
	fake.getPrivateDataMetadataByHashMutex.RLock()
	defer fake.getPrivateDataMetadataByHashMutex.RUnlock()
	return len(fake.getPrivateDataMetadataByHashArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashCalls(stub func(string, string, []byte) (map[string][]byte, error)) {
	fake.getPrivateDataMetadataByHashMutex.Lock()
	defer fake.getPrivateDataMetadataByHashMutex.Unlock()
	fake.GetPrivateDataMetadataByHashStub = stub
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashArgsForCall(i int) (string, string, []byte) {
	fake.getPrivateDataMetadataByHashMutex.RLock()
	defer fake.getPrivateDataMetadataByHashMutex.RUnlock()
	argsForCall := fake.getPrivateDataMetadataByHashArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashReturns(result1 map[string][]byte, result2 error) {
	fake.getPrivateDataMetadataByHashMutex.Lock()
	defer fake.getPrivateDataMetadataByHashMutex.Unlock()
	fake.GetPrivateDataMetadataByHashStub = nil
	fake.getPrivateDataMetadataByHashReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
	fake.getPrivateDataMetadataByHashMutex.Lock()
	defer fake.getPrivateDataMetadataByHashMutex.Unlock()
	fake.GetPrivateDataMetadataByHashStub = nil
	if fake.getPrivateDataMetadataByHashReturnsOnCall == nil {
		fake.getPrivateDataMetadataByHashReturnsOnCall = make(map[int]struct {
			result1 map[string][]byte
			result2 error
		})
	}
	fake.getPrivateDataMetadataByHashReturnsOnCall[i] = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMultipleKeys(arg1 string, arg2 string, arg3 []string) ([][]byte, error) {
	var arg3Copy []string
	if arg3 != nil {
		arg3Copy = make([]string, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.getPrivateDataMultipleKeysMutex.Lock()
	ret, specificReturn := fake.getPrivateDataMultipleKeysReturnsOnCall[len(fake.getPrivateDataMultipleKeysArgsForCall)]
	fake.getPrivateDataMultipleKeysArgsForCall = append(fake.getPrivateDataMultipleKeysArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 []string
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("GetPrivateDataMultipleKeys", []interface{}{arg1, arg2, arg3Copy})
	fake.getPrivateDataMultipleKeysMutex.Unlock()
	if fake.GetPrivateDataMultipleKeysStub != nil {
		return fake.GetPrivateDataMultipleKeysStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPrivateDataMultipleKeysReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetPrivateDataMultipleKeysCallCount() int {
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	return len(fake.getPrivateDataMultipleKeysArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataMultipleKeysCalls(stub func(string, string, []string) ([][]byte, error)) {
	fake.getPrivateDataMultipleKeysMutex.Lock()
	defer fake.getPrivateDataMultipleKeysMutex.Unlock()
	fake.GetPrivateDataMultipleKeysStub = stub
}

func (fake *TxSimulator) GetPrivateDataMultipleKeysArgsForCall(i int) (string, string, []string) {
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	argsForCall := fake.getPrivateDataMultipleKeysArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) GetPrivateDataMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.getPrivateDataMultipleKeysMutex.Lock()
	defer fake.getPrivateDataMultipleKeysMutex.Unlock()
	fake.GetPrivateDataMultipleKeysStub = nil
	fake.getPrivateDataMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
	fake.getPrivateDataMultipleKeysMutex.Lock()
	defer fake.getPrivateDataMultipleKeysMutex.Unlock()
	fake.GetPrivateDataMultipleKeysStub = nil
	if fake.getPrivateDataMultipleKeysReturnsOnCall == nil {
		fake.getPrivateDataMultipleKeysReturnsOnCall = make(map[int]struct {
			result1 [][]byte
			result2 error
		})
	}
	fake.getPrivateDataMultipleKeysReturnsOnCall[i] = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataRangeScanIterator(arg1 string, arg2 string, arg3 string, arg4 string) (ledgera.ResultsIterator, error) {
	fake.getPrivateDataRangeScanIteratorMutex.Lock()
	ret, specificReturn := fake.getPrivateDataRangeScanIteratorReturnsOnCall[len(fake.getPrivateDataRangeScanIteratorArgsForCall)]
	fake.getPrivateDataRangeScanIteratorArgsForCall = append(fake.getPrivateDataRangeScanIteratorArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("GetPrivateDataRangeScanIterator", []interface{}{arg1, arg2, arg3, arg4})
	fake.getPrivateDataRangeScanIteratorMutex.Unlock()
	if fake.GetPrivateDataRangeScanIteratorStub != nil {
		return fake.GetPrivateDataRangeScanIteratorStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPrivateDataRangeScanIteratorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorCallCount() int {
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	return len(fake.getPrivateDataRangeScanIteratorArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorCalls(stub func(string, string, string, string) (ledgera.ResultsIterator, error)) {
	fake.getPrivateDataRangeScanIteratorMutex.Lock()
	defer fake.getPrivateDataRangeScanIteratorMutex.Unlock()
	fake.GetPrivateDataRangeScanIteratorStub = stub
}

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorArgsForCall(i int) (string, string, string, string) {
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	argsForCall := fake.getPrivateDataRangeScanIteratorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.getPrivateDataRangeScanIteratorMutex.Lock()
	defer fake.getPrivateDataRangeScanIteratorMutex.Unlock()
	fake.GetPrivateDataRangeScanIteratorStub = nil
	fake.getPrivateDataRangeScanIteratorReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
	fake.getPrivateDataRangeScanIteratorMutex.Lock()
	defer fake.getPrivateDataRangeScanIteratorMutex.Unlock()
	fake.GetPrivateDataRangeScanIteratorStub = nil
	if fake.getPrivateDataRangeScanIteratorReturnsOnCall == nil {
		fake.getPrivateDataRangeScanIteratorReturnsOnCall = make(map[int]struct {
			result1 ledgera.ResultsIterator
			result2 error
		})
	}
	fake.getPrivateDataRangeScanIteratorReturnsOnCall[i] = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetState(arg1 string, arg2 string) ([]byte, error) {
	fake.getStateMutex.Lock()
	ret, specificReturn := fake.getStateReturnsOnCall[len(fake.getStateArgsForCall)]
	fake.getStateArgsForCall = append(fake.getStateArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetState", []interface{}{arg1, arg2})
	fake.getStateMutex.Unlock()
	if fake.GetStateStub != nil {
		return fake.GetStateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *TxSimulator) GetStateCalls(stub func(string, string) ([]byte, error)) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = stub
}

func (fake *TxSimulator) GetStateArgsForCall(i int) (string, string) {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	argsForCall := fake.getStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *TxSimulator) GetStateReturns(result1 []byte, result2 error) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = nil
	if fake.getStateReturnsOnCall == nil {
		fake.getStateReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getStateReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateMetadata(arg1 string, arg2 string) (map[string][]byte, error) {
	fake.getStateMetadataMutex.Lock()
	ret, specificReturn := fake.getStateMetadataReturnsOnCall[len(fake.getStateMetadataArgsForCall)]
	fake.getStateMetadataArgsForCall = append(fake.getStateMetadataArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetStateMetadata", []interface{}{arg1, arg2})
	fake.getStateMetadataMutex.Unlock()
	if fake.GetStateMetadataStub != nil {
		return fake.GetStateMetadataStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateMetadataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetStateMetadataCallCount() int {
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	return len(fake.getStateMetadataArgsForCall)
}

func (fake *TxSimulator) GetStateMetadataCalls(stub func(string, string) (map[string][]byte, error)) {
	fake.getStateMetadataMutex.Lock()
	defer fake.getStateMetadataMutex.Unlock()
	fake.GetStateMetadataStub = stub
}

func (fake *TxSimulator) GetStateMetadataArgsForCall(i int) (string, string) {
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	argsForCall := fake.getStateMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *TxSimulator) GetStateMetadataReturns(result1 map[string][]byte, result2 error) {
	fake.getStateMetadataMutex.Lock()
	defer fake.getStateMetadataMutex.Unlock()
	fake.GetStateMetadataStub = nil
	fake.getStateMetadataReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateMetadataReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
	fake.getStateMetadataMutex.Lock()
	defer fake.getStateMetadataMutex.Unlock()
	fake.GetStateMetadataStub = nil
	if fake.getStateMetadataReturnsOnCall == nil {
		fake.getStateMetadataReturnsOnCall = make(map[int]struct {
			result1 map[string][]byte
			result2 error
		})
	}
	fake.getStateMetadataReturnsOnCall[i] = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateMultipleKeys(arg1 string, arg2 []string) ([][]byte, error) {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.getStateMultipleKeysMutex.Lock()
	ret, specificReturn := fake.getStateMultipleKeysReturnsOnCall[len(fake.getStateMultipleKeysArgsForCall)]
	fake.getStateMultipleKeysArgsForCall = append(fake.getStateMultipleKeysArgsForCall, struct {
		arg1 string
		arg2 []string
	}{arg1, arg2Copy})
	fake.recordInvocation("GetStateMultipleKeys", []interface{}{arg1, arg2Copy})
	fake.getStateMultipleKeysMutex.Unlock()
	if fake.GetStateMultipleKeysStub != nil {
		return fake.GetStateMultipleKeysStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateMultipleKeysReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetStateMultipleKeysCallCount() int {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	return len(fake.getStateMultipleKeysArgsForCall)
}

func (fake *TxSimulator) GetStateMultipleKeysCalls(stub func(string, []string) ([][]byte, error)) {
	fake.getStateMultipleKeysMutex.Lock()
	defer fake.getStateMultipleKeysMutex.Unlock()
	fake.GetStateMultipleKeysStub = stub
}

func (fake *TxSimulator) GetStateMultipleKeysArgsForCall(i int) (string, []string) {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	argsForCall := fake.getStateMultipleKeysArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *TxSimulator) GetStateMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.getStateMultipleKeysMutex.Lock()
	defer fake.getStateMultipleKeysMutex.Unlock()
	fake.GetStateMultipleKeysStub = nil
	fake.getStateMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
	fake.getStateMultipleKeysMutex.Lock()
	defer fake.getStateMultipleKeysMutex.Unlock()
	fake.GetStateMultipleKeysStub = nil
	if fake.getStateMultipleKeysReturnsOnCall == nil {
		fake.getStateMultipleKeysReturnsOnCall = make(map[int]struct {
			result1 [][]byte
			result2 error
		})
	}
	fake.getStateMultipleKeysReturnsOnCall[i] = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateRangeScanIterator(arg1 string, arg2 string, arg3 string) (ledgera.ResultsIterator, error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	ret, specificReturn := fake.getStateRangeScanIteratorReturnsOnCall[len(fake.getStateRangeScanIteratorArgsForCall)]
	fake.getStateRangeScanIteratorArgsForCall = append(fake.getStateRangeScanIteratorArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetStateRangeScanIterator", []interface{}{arg1, arg2, arg3})
	fake.getStateRangeScanIteratorMutex.Unlock()
	if fake.GetStateRangeScanIteratorStub != nil {
		return fake.GetStateRangeScanIteratorStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateRangeScanIteratorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetStateRangeScanIteratorCallCount() int {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorArgsForCall)
}

func (fake *TxSimulator) GetStateRangeScanIteratorCalls(stub func(string, string, string) (ledgera.ResultsIterator, error)) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = stub
}

func (fake *TxSimulator) GetStateRangeScanIteratorArgsForCall(i int) (string, string, string) {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	argsForCall := fake.getStateRangeScanIteratorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) GetStateRangeScanIteratorReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = nil
	fake.getStateRangeScanIteratorReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateRangeScanIteratorReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = nil
	if fake.getStateRangeScanIteratorReturnsOnCall == nil {
		fake.getStateRangeScanIteratorReturnsOnCall = make(map[int]struct {
			result1 ledgera.ResultsIterator
			result2 error
		})
	}
	fake.getStateRangeScanIteratorReturnsOnCall[i] = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadata(arg1 string, arg2 string, arg3 string, arg4 map[string]interface{}) (ledger.QueryResultsIterator, error) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	ret, specificReturn := fake.getStateRangeScanIteratorWithMetadataReturnsOnCall[len(fake.getStateRangeScanIteratorWithMetadataArgsForCall)]
	fake.getStateRangeScanIteratorWithMetadataArgsForCall = append(fake.getStateRangeScanIteratorWithMetadataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 map[string]interface{}
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("GetStateRangeScanIteratorWithMetadata", []interface{}{arg1, arg2, arg3, arg4})
	fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	if fake.GetStateRangeScanIteratorWithMetadataStub != nil {
		return fake.GetStateRangeScanIteratorWithMetadataStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateRangeScanIteratorWithMetadataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataCallCount() int {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorWithMetadataArgsForCall)
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataCalls(stub func(string, string, string, map[string]interface{}) (ledger.QueryResultsIterator, error)) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	fake.GetStateRangeScanIteratorWithMetadataStub = stub
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataArgsForCall(i int) (string, string, string, map[string]interface{}) {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	argsForCall := fake.getStateRangeScanIteratorWithMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataReturns(result1 ledger.QueryResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	fake.GetStateRangeScanIteratorWithMetadataStub = nil
	fake.getStateRangeScanIteratorWithMetadataReturns = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataReturnsOnCall(i int, result1 ledger.QueryResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	fake.GetStateRangeScanIteratorWithMetadataStub = nil
	if fake.getStateRangeScanIteratorWithMetadataReturnsOnCall == nil {
		fake.getStateRangeScanIteratorWithMetadataReturnsOnCall = make(map[int]struct {
			result1 ledger.QueryResultsIterator
			result2 error
		})
	}
	fake.getStateRangeScanIteratorWithMetadataReturnsOnCall[i] = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetTxSimulationResults() (*ledger.TxSimulationResults, error) {
	fake.getTxSimulationResultsMutex.Lock()
	ret, specificReturn := fake.getTxSimulationResultsReturnsOnCall[len(fake.getTxSimulationResultsArgsForCall)]
	fake.getTxSimulationResultsArgsForCall = append(fake.getTxSimulationResultsArgsForCall, struct {
	}{})
	fake.recordInvocation("GetTxSimulationResults", []interface{}{})
	fake.getTxSimulationResultsMutex.Unlock()
	if fake.GetTxSimulationResultsStub != nil {
		return fake.GetTxSimulationResultsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getTxSimulationResultsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TxSimulator) GetTxSimulationResultsCallCount() int {
	fake.getTxSimulationResultsMutex.RLock()
	defer fake.getTxSimulationResultsMutex.RUnlock()
	return len(fake.getTxSimulationResultsArgsForCall)
}

func (fake *TxSimulator) GetTxSimulationResultsCalls(stub func() (*ledger.TxSimulationResults, error)) {
	fake.getTxSimulationResultsMutex.Lock()
	defer fake.getTxSimulationResultsMutex.Unlock()
	fake.GetTxSimulationResultsStub = stub
}

func (fake *TxSimulator) GetTxSimulationResultsReturns(result1 *ledger.TxSimulationResults, result2 error) {
	fake.getTxSimulationResultsMutex.Lock()
	defer fake.getTxSimulationResultsMutex.Unlock()
	fake.GetTxSimulationResultsStub = nil
	fake.getTxSimulationResultsReturns = struct {
		result1 *ledger.TxSimulationResults
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetTxSimulationResultsReturnsOnCall(i int, result1 *ledger.TxSimulationResults, result2 error) {
	fake.getTxSimulationResultsMutex.Lock()
	defer fake.getTxSimulationResultsMutex.Unlock()
	fake.GetTxSimulationResultsStub = nil
	if fake.getTxSimulationResultsReturnsOnCall == nil {
		fake.getTxSimulationResultsReturnsOnCall = make(map[int]struct {
			result1 *ledger.TxSimulationResults
			result2 error
		})
	}
	fake.getTxSimulationResultsReturnsOnCall[i] = struct {
		result1 *ledger.TxSimulationResults
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) SetPrivateData(arg1 string, arg2 string, arg3 string, arg4 []byte) error {
	var arg4Copy []byte
	if arg4 != nil {
		arg4Copy = make([]byte, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.setPrivateDataMutex.Lock()
	ret, specificReturn := fake.setPrivateDataReturnsOnCall[len(fake.setPrivateDataArgsForCall)]
	fake.setPrivateDataArgsForCall = append(fake.setPrivateDataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 []byte
	}{arg1, arg2, arg3, arg4Copy})
	fake.recordInvocation("SetPrivateData", []interface{}{arg1, arg2, arg3, arg4Copy})
	fake.setPrivateDataMutex.Unlock()
	if fake.SetPrivateDataStub != nil {
		return fake.SetPrivateDataStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setPrivateDataReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) SetPrivateDataCallCount() int {
	fake.setPrivateDataMutex.RLock()
	defer fake.setPrivateDataMutex.RUnlock()
	return len(fake.setPrivateDataArgsForCall)
}

func (fake *TxSimulator) SetPrivateDataCalls(stub func(string, string, string, []byte) error) {
	fake.setPrivateDataMutex.Lock()
	defer fake.setPrivateDataMutex.Unlock()
	fake.SetPrivateDataStub = stub
}

func (fake *TxSimulator) SetPrivateDataArgsForCall(i int) (string, string, string, []byte) {
	fake.setPrivateDataMutex.RLock()
	defer fake.setPrivateDataMutex.RUnlock()
	argsForCall := fake.setPrivateDataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *TxSimulator) SetPrivateDataReturns(result1 error) {
	fake.setPrivateDataMutex.Lock()
	defer fake.setPrivateDataMutex.Unlock()
	fake.SetPrivateDataStub = nil
	fake.setPrivateDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetPrivateDataReturnsOnCall(i int, result1 error) {
	fake.setPrivateDataMutex.Lock()
	defer fake.setPrivateDataMutex.Unlock()
	fake.SetPrivateDataStub = nil
	if fake.setPrivateDataReturnsOnCall == nil {
		fake.setPrivateDataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setPrivateDataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetPrivateDataMetadata(arg1 string, arg2 string, arg3 string, arg4 map[string][]byte) error {
	fake.setPrivateDataMetadataMutex.Lock()
	ret, specificReturn := fake.setPrivateDataMetadataReturnsOnCall[len(fake.setPrivateDataMetadataArgsForCall)]
	fake.setPrivateDataMetadataArgsForCall = append(fake.setPrivateDataMetadataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 map[string][]byte
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("SetPrivateDataMetadata", []interface{}{arg1, arg2, arg3, arg4})
	fake.setPrivateDataMetadataMutex.Unlock()
	if fake.SetPrivateDataMetadataStub != nil {
		return fake.SetPrivateDataMetadataStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setPrivateDataMetadataReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) SetPrivateDataMetadataCallCount() int {
	fake.setPrivateDataMetadataMutex.RLock()
	defer fake.setPrivateDataMetadataMutex.RUnlock()
	return len(fake.setPrivateDataMetadataArgsForCall)
}

func (fake *TxSimulator) SetPrivateDataMetadataCalls(stub func(string, string, string, map[string][]byte) error) {
	fake.setPrivateDataMetadataMutex.Lock()
	defer fake.setPrivateDataMetadataMutex.Unlock()
	fake.SetPrivateDataMetadataStub = stub
}

func (fake *TxSimulator) SetPrivateDataMetadataArgsForCall(i int) (string, string, string, map[string][]byte) {
	fake.setPrivateDataMetadataMutex.RLock()
	defer fake.setPrivateDataMetadataMutex.RUnlock()
	argsForCall := fake.setPrivateDataMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *TxSimulator) SetPrivateDataMetadataReturns(result1 error) {
	fake.setPrivateDataMetadataMutex.Lock()
	defer fake.setPrivateDataMetadataMutex.Unlock()
	fake.SetPrivateDataMetadataStub = nil
	fake.setPrivateDataMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetPrivateDataMetadataReturnsOnCall(i int, result1 error) {
	fake.setPrivateDataMetadataMutex.Lock()
	defer fake.setPrivateDataMetadataMutex.Unlock()
	fake.SetPrivateDataMetadataStub = nil
	if fake.setPrivateDataMetadataReturnsOnCall == nil {
		fake.setPrivateDataMetadataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setPrivateDataMetadataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetPrivateDataMultipleKeys(arg1 string, arg2 string, arg3 map[string][]byte) error {
	fake.setPrivateDataMultipleKeysMutex.Lock()
	ret, specificReturn := fake.setPrivateDataMultipleKeysReturnsOnCall[len(fake.setPrivateDataMultipleKeysArgsForCall)]
	fake.setPrivateDataMultipleKeysArgsForCall = append(fake.setPrivateDataMultipleKeysArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 map[string][]byte
	}{arg1, arg2, arg3})
	fake.recordInvocation("SetPrivateDataMultipleKeys", []interface{}{arg1, arg2, arg3})
	fake.setPrivateDataMultipleKeysMutex.Unlock()
	if fake.SetPrivateDataMultipleKeysStub != nil {
		return fake.SetPrivateDataMultipleKeysStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setPrivateDataMultipleKeysReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysCallCount() int {
	fake.setPrivateDataMultipleKeysMutex.RLock()
	defer fake.setPrivateDataMultipleKeysMutex.RUnlock()
	return len(fake.setPrivateDataMultipleKeysArgsForCall)
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysCalls(stub func(string, string, map[string][]byte) error) {
	fake.setPrivateDataMultipleKeysMutex.Lock()
	defer fake.setPrivateDataMultipleKeysMutex.Unlock()
	fake.SetPrivateDataMultipleKeysStub = stub
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysArgsForCall(i int) (string, string, map[string][]byte) {
	fake.setPrivateDataMultipleKeysMutex.RLock()
	defer fake.setPrivateDataMultipleKeysMutex.RUnlock()
	argsForCall := fake.setPrivateDataMultipleKeysArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysReturns(result1 error) {
	fake.setPrivateDataMultipleKeysMutex.Lock()
	defer fake.setPrivateDataMultipleKeysMutex.Unlock()
	fake.SetPrivateDataMultipleKeysStub = nil
	fake.setPrivateDataMultipleKeysReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysReturnsOnCall(i int, result1 error) {
	fake.setPrivateDataMultipleKeysMutex.Lock()
	defer fake.setPrivateDataMultipleKeysMutex.Unlock()
	fake.SetPrivateDataMultipleKeysStub = nil
	if fake.setPrivateDataMultipleKeysReturnsOnCall == nil {
		fake.setPrivateDataMultipleKeysReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setPrivateDataMultipleKeysReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetState(arg1 string, arg2 string, arg3 []byte) error {
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.setStateMutex.Lock()
	ret, specificReturn := fake.setStateReturnsOnCall[len(fake.setStateArgsForCall)]
	fake.setStateArgsForCall = append(fake.setStateArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 []byte
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("SetState", []interface{}{arg1, arg2, arg3Copy})
	fake.setStateMutex.Unlock()
	if fake.SetStateStub != nil {
		return fake.SetStateStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setStateReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) SetStateCallCount() int {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return len(fake.setStateArgsForCall)
}

func (fake *TxSimulator) SetStateCalls(stub func(string, string, []byte) error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
	fake.SetStateStub = stub
}

func (fake *TxSimulator) SetStateArgsForCall(i int) (string, string, []byte) {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	argsForCall := fake.setStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) SetStateReturns(result1 error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
	fake.SetStateStub = nil
	fake.setStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetStateReturnsOnCall(i int, result1 error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
	fake.SetStateStub = nil
	if fake.setStateReturnsOnCall == nil {
		fake.setStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetStateMetadata(arg1 string, arg2 string, arg3 map[string][]byte) error {
	fake.setStateMetadataMutex.Lock()
	ret, specificReturn := fake.setStateMetadataReturnsOnCall[len(fake.setStateMetadataArgsForCall)]
	fake.setStateMetadataArgsForCall = append(fake.setStateMetadataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 map[string][]byte
	}{arg1, arg2, arg3})
	fake.recordInvocation("SetStateMetadata", []interface{}{arg1, arg2, arg3})
	fake.setStateMetadataMutex.Unlock()
	if fake.SetStateMetadataStub != nil {
		return fake.SetStateMetadataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setStateMetadataReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) SetStateMetadataCallCount() int {
	fake.setStateMetadataMutex.RLock()
	defer fake.setStateMetadataMutex.RUnlock()
	return len(fake.setStateMetadataArgsForCall)
}

func (fake *TxSimulator) SetStateMetadataCalls(stub func(string, string, map[string][]byte) error) {
	fake.setStateMetadataMutex.Lock()
	defer fake.setStateMetadataMutex.Unlock()
	fake.SetStateMetadataStub = stub
}

func (fake *TxSimulator) SetStateMetadataArgsForCall(i int) (string, string, map[string][]byte) {
	fake.setStateMetadataMutex.RLock()
	defer fake.setStateMetadataMutex.RUnlock()
	argsForCall := fake.setStateMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TxSimulator) SetStateMetadataReturns(result1 error) {
	fake.setStateMetadataMutex.Lock()
	defer fake.setStateMetadataMutex.Unlock()
	fake.SetStateMetadataStub = nil
	fake.setStateMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetStateMetadataReturnsOnCall(i int, result1 error) {
	fake.setStateMetadataMutex.Lock()
	defer fake.setStateMetadataMutex.Unlock()
	fake.SetStateMetadataStub = nil
	if fake.setStateMetadataReturnsOnCall == nil {
		fake.setStateMetadataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setStateMetadataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetStateMultipleKeys(arg1 string, arg2 map[string][]byte) error {
	fake.setStateMultipleKeysMutex.Lock()
	ret, specificReturn := fake.setStateMultipleKeysReturnsOnCall[len(fake.setStateMultipleKeysArgsForCall)]
	fake.setStateMultipleKeysArgsForCall = append(fake.setStateMultipleKeysArgsForCall, struct {
		arg1 string
		arg2 map[string][]byte
	}{arg1, arg2})
	fake.recordInvocation("SetStateMultipleKeys", []interface{}{arg1, arg2})
	fake.setStateMultipleKeysMutex.Unlock()
	if fake.SetStateMultipleKeysStub != nil {
		return fake.SetStateMultipleKeysStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setStateMultipleKeysReturns
	return fakeReturns.result1
}

func (fake *TxSimulator) SetStateMultipleKeysCallCount() int {
	fake.setStateMultipleKeysMutex.RLock()
	defer fake.setStateMultipleKeysMutex.RUnlock()
	return len(fake.setStateMultipleKeysArgsForCall)
}

func (fake *TxSimulator) SetStateMultipleKeysCalls(stub func(string, map[string][]byte) error) {
	fake.setStateMultipleKeysMutex.Lock()
	defer fake.setStateMultipleKeysMutex.Unlock()
	fake.SetStateMultipleKeysStub = stub
}

func (fake *TxSimulator) SetStateMultipleKeysArgsForCall(i int) (string, map[string][]byte) {
	fake.setStateMultipleKeysMutex.RLock()
	defer fake.setStateMultipleKeysMutex.RUnlock()
	argsForCall := fake.setStateMultipleKeysArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *TxSimulator) SetStateMultipleKeysReturns(result1 error) {
	fake.setStateMultipleKeysMutex.Lock()
	defer fake.setStateMultipleKeysMutex.Unlock()
	fake.SetStateMultipleKeysStub = nil
	fake.setStateMultipleKeysReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetStateMultipleKeysReturnsOnCall(i int, result1 error) {
	fake.setStateMultipleKeysMutex.Lock()
	defer fake.setStateMultipleKeysMutex.Unlock()
	fake.SetStateMultipleKeysStub = nil
	if fake.setStateMultipleKeysReturnsOnCall == nil {
		fake.setStateMultipleKeysReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setStateMultipleKeysReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deletePrivateDataMutex.RLock()
	defer fake.deletePrivateDataMutex.RUnlock()
	fake.deletePrivateDataMetadataMutex.RLock()
	defer fake.deletePrivateDataMetadataMutex.RUnlock()
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	fake.deleteStateMetadataMutex.RLock()
	defer fake.deleteStateMetadataMutex.RUnlock()
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	fake.executeUpdateMutex.RLock()
	defer fake.executeUpdateMutex.RUnlock()
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	fake.getPrivateDataMetadataByHashMutex.RLock()
	defer fake.getPrivateDataMetadataByHashMutex.RUnlock()
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	fake.getTxSimulationResultsMutex.RLock()
	defer fake.getTxSimulationResultsMutex.RUnlock()
	fake.setPrivateDataMutex.RLock()
	defer fake.setPrivateDataMutex.RUnlock()
	fake.setPrivateDataMetadataMutex.RLock()
	defer fake.setPrivateDataMetadataMutex.RUnlock()
	fake.setPrivateDataMultipleKeysMutex.RLock()
	defer fake.setPrivateDataMultipleKeysMutex.RUnlock()
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	fake.setStateMetadataMutex.RLock()
	defer fake.setStateMetadataMutex.RUnlock()
	fake.setStateMultipleKeysMutex.RLock()
	defer fake.setStateMultipleKeysMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TxSimulator) recordInvocation(key string, args []interface{}) {
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

var _ ledger.TxSimulator = new(TxSimulator)
