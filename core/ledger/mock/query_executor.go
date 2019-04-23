
package mock

import (
	"sync"

	ledgera "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger"
)

type QueryExecutor struct {
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
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *QueryExecutor) Done() {
	fake.doneMutex.Lock()
	fake.doneArgsForCall = append(fake.doneArgsForCall, struct {
	}{})
	fake.recordInvocation("Done", []interface{}{})
	fake.doneMutex.Unlock()
	if fake.DoneStub != nil {
		fake.DoneStub()
	}
}

func (fake *QueryExecutor) DoneCallCount() int {
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
	return len(fake.doneArgsForCall)
}

func (fake *QueryExecutor) DoneCalls(stub func()) {
	fake.doneMutex.Lock()
	defer fake.doneMutex.Unlock()
	fake.DoneStub = stub
}

func (fake *QueryExecutor) ExecuteQuery(arg1 string, arg2 string) (ledgera.ResultsIterator, error) {
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

func (fake *QueryExecutor) ExecuteQueryCallCount() int {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	return len(fake.executeQueryArgsForCall)
}

func (fake *QueryExecutor) ExecuteQueryCalls(stub func(string, string) (ledgera.ResultsIterator, error)) {
	fake.executeQueryMutex.Lock()
	defer fake.executeQueryMutex.Unlock()
	fake.ExecuteQueryStub = stub
}

func (fake *QueryExecutor) ExecuteQueryArgsForCall(i int) (string, string) {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	argsForCall := fake.executeQueryArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *QueryExecutor) ExecuteQueryReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.executeQueryMutex.Lock()
	defer fake.executeQueryMutex.Unlock()
	fake.ExecuteQueryStub = nil
	fake.executeQueryReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) ExecuteQueryReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
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

func (fake *QueryExecutor) ExecuteQueryOnPrivateData(arg1 string, arg2 string, arg3 string) (ledgera.ResultsIterator, error) {
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

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataCallCount() int {
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	return len(fake.executeQueryOnPrivateDataArgsForCall)
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataCalls(stub func(string, string, string) (ledgera.ResultsIterator, error)) {
	fake.executeQueryOnPrivateDataMutex.Lock()
	defer fake.executeQueryOnPrivateDataMutex.Unlock()
	fake.ExecuteQueryOnPrivateDataStub = stub
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataArgsForCall(i int) (string, string, string) {
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	argsForCall := fake.executeQueryOnPrivateDataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.executeQueryOnPrivateDataMutex.Lock()
	defer fake.executeQueryOnPrivateDataMutex.Unlock()
	fake.ExecuteQueryOnPrivateDataStub = nil
	fake.executeQueryOnPrivateDataReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
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

func (fake *QueryExecutor) ExecuteQueryWithMetadata(arg1 string, arg2 string, arg3 map[string]interface{}) (ledger.QueryResultsIterator, error) {
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

func (fake *QueryExecutor) ExecuteQueryWithMetadataCallCount() int {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	return len(fake.executeQueryWithMetadataArgsForCall)
}

func (fake *QueryExecutor) ExecuteQueryWithMetadataCalls(stub func(string, string, map[string]interface{}) (ledger.QueryResultsIterator, error)) {
	fake.executeQueryWithMetadataMutex.Lock()
	defer fake.executeQueryWithMetadataMutex.Unlock()
	fake.ExecuteQueryWithMetadataStub = stub
}

func (fake *QueryExecutor) ExecuteQueryWithMetadataArgsForCall(i int) (string, string, map[string]interface{}) {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	argsForCall := fake.executeQueryWithMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *QueryExecutor) ExecuteQueryWithMetadataReturns(result1 ledger.QueryResultsIterator, result2 error) {
	fake.executeQueryWithMetadataMutex.Lock()
	defer fake.executeQueryWithMetadataMutex.Unlock()
	fake.ExecuteQueryWithMetadataStub = nil
	fake.executeQueryWithMetadataReturns = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) ExecuteQueryWithMetadataReturnsOnCall(i int, result1 ledger.QueryResultsIterator, result2 error) {
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

func (fake *QueryExecutor) GetPrivateData(arg1 string, arg2 string, arg3 string) ([]byte, error) {
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

func (fake *QueryExecutor) GetPrivateDataCallCount() int {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	return len(fake.getPrivateDataArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataCalls(stub func(string, string, string) ([]byte, error)) {
	fake.getPrivateDataMutex.Lock()
	defer fake.getPrivateDataMutex.Unlock()
	fake.GetPrivateDataStub = stub
}

func (fake *QueryExecutor) GetPrivateDataArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	argsForCall := fake.getPrivateDataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *QueryExecutor) GetPrivateDataReturns(result1 []byte, result2 error) {
	fake.getPrivateDataMutex.Lock()
	defer fake.getPrivateDataMutex.Unlock()
	fake.GetPrivateDataStub = nil
	fake.getPrivateDataReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *QueryExecutor) GetPrivateDataHash(arg1 string, arg2 string, arg3 string) ([]byte, error) {
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

func (fake *QueryExecutor) GetPrivateDataHashCallCount() int {
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	return len(fake.getPrivateDataHashArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataHashCalls(stub func(string, string, string) ([]byte, error)) {
	fake.getPrivateDataHashMutex.Lock()
	defer fake.getPrivateDataHashMutex.Unlock()
	fake.GetPrivateDataHashStub = stub
}

func (fake *QueryExecutor) GetPrivateDataHashArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	argsForCall := fake.getPrivateDataHashArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *QueryExecutor) GetPrivateDataHashReturns(result1 []byte, result2 error) {
	fake.getPrivateDataHashMutex.Lock()
	defer fake.getPrivateDataHashMutex.Unlock()
	fake.GetPrivateDataHashStub = nil
	fake.getPrivateDataHashReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataHashReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *QueryExecutor) GetPrivateDataMetadata(arg1 string, arg2 string, arg3 string) (map[string][]byte, error) {
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

func (fake *QueryExecutor) GetPrivateDataMetadataCallCount() int {
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	return len(fake.getPrivateDataMetadataArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataMetadataCalls(stub func(string, string, string) (map[string][]byte, error)) {
	fake.getPrivateDataMetadataMutex.Lock()
	defer fake.getPrivateDataMetadataMutex.Unlock()
	fake.GetPrivateDataMetadataStub = stub
}

func (fake *QueryExecutor) GetPrivateDataMetadataArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	argsForCall := fake.getPrivateDataMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *QueryExecutor) GetPrivateDataMetadataReturns(result1 map[string][]byte, result2 error) {
	fake.getPrivateDataMetadataMutex.Lock()
	defer fake.getPrivateDataMetadataMutex.Unlock()
	fake.GetPrivateDataMetadataStub = nil
	fake.getPrivateDataMetadataReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataMetadataReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
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

func (fake *QueryExecutor) GetPrivateDataMetadataByHash(arg1 string, arg2 string, arg3 []byte) (map[string][]byte, error) {
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

func (fake *QueryExecutor) GetPrivateDataMetadataByHashCallCount() int {
	fake.getPrivateDataMetadataByHashMutex.RLock()
	defer fake.getPrivateDataMetadataByHashMutex.RUnlock()
	return len(fake.getPrivateDataMetadataByHashArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataMetadataByHashCalls(stub func(string, string, []byte) (map[string][]byte, error)) {
	fake.getPrivateDataMetadataByHashMutex.Lock()
	defer fake.getPrivateDataMetadataByHashMutex.Unlock()
	fake.GetPrivateDataMetadataByHashStub = stub
}

func (fake *QueryExecutor) GetPrivateDataMetadataByHashArgsForCall(i int) (string, string, []byte) {
	fake.getPrivateDataMetadataByHashMutex.RLock()
	defer fake.getPrivateDataMetadataByHashMutex.RUnlock()
	argsForCall := fake.getPrivateDataMetadataByHashArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *QueryExecutor) GetPrivateDataMetadataByHashReturns(result1 map[string][]byte, result2 error) {
	fake.getPrivateDataMetadataByHashMutex.Lock()
	defer fake.getPrivateDataMetadataByHashMutex.Unlock()
	fake.GetPrivateDataMetadataByHashStub = nil
	fake.getPrivateDataMetadataByHashReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataMetadataByHashReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
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

func (fake *QueryExecutor) GetPrivateDataMultipleKeys(arg1 string, arg2 string, arg3 []string) ([][]byte, error) {
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

func (fake *QueryExecutor) GetPrivateDataMultipleKeysCallCount() int {
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	return len(fake.getPrivateDataMultipleKeysArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataMultipleKeysCalls(stub func(string, string, []string) ([][]byte, error)) {
	fake.getPrivateDataMultipleKeysMutex.Lock()
	defer fake.getPrivateDataMultipleKeysMutex.Unlock()
	fake.GetPrivateDataMultipleKeysStub = stub
}

func (fake *QueryExecutor) GetPrivateDataMultipleKeysArgsForCall(i int) (string, string, []string) {
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	argsForCall := fake.getPrivateDataMultipleKeysArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *QueryExecutor) GetPrivateDataMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.getPrivateDataMultipleKeysMutex.Lock()
	defer fake.getPrivateDataMultipleKeysMutex.Unlock()
	fake.GetPrivateDataMultipleKeysStub = nil
	fake.getPrivateDataMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
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

func (fake *QueryExecutor) GetPrivateDataRangeScanIterator(arg1 string, arg2 string, arg3 string, arg4 string) (ledgera.ResultsIterator, error) {
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

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorCallCount() int {
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	return len(fake.getPrivateDataRangeScanIteratorArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorCalls(stub func(string, string, string, string) (ledgera.ResultsIterator, error)) {
	fake.getPrivateDataRangeScanIteratorMutex.Lock()
	defer fake.getPrivateDataRangeScanIteratorMutex.Unlock()
	fake.GetPrivateDataRangeScanIteratorStub = stub
}

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorArgsForCall(i int) (string, string, string, string) {
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	argsForCall := fake.getPrivateDataRangeScanIteratorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.getPrivateDataRangeScanIteratorMutex.Lock()
	defer fake.getPrivateDataRangeScanIteratorMutex.Unlock()
	fake.GetPrivateDataRangeScanIteratorStub = nil
	fake.getPrivateDataRangeScanIteratorReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
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

func (fake *QueryExecutor) GetState(arg1 string, arg2 string) ([]byte, error) {
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

func (fake *QueryExecutor) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *QueryExecutor) GetStateCalls(stub func(string, string) ([]byte, error)) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = stub
}

func (fake *QueryExecutor) GetStateArgsForCall(i int) (string, string) {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	argsForCall := fake.getStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *QueryExecutor) GetStateReturns(result1 []byte, result2 error) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *QueryExecutor) GetStateMetadata(arg1 string, arg2 string) (map[string][]byte, error) {
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

func (fake *QueryExecutor) GetStateMetadataCallCount() int {
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	return len(fake.getStateMetadataArgsForCall)
}

func (fake *QueryExecutor) GetStateMetadataCalls(stub func(string, string) (map[string][]byte, error)) {
	fake.getStateMetadataMutex.Lock()
	defer fake.getStateMetadataMutex.Unlock()
	fake.GetStateMetadataStub = stub
}

func (fake *QueryExecutor) GetStateMetadataArgsForCall(i int) (string, string) {
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	argsForCall := fake.getStateMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *QueryExecutor) GetStateMetadataReturns(result1 map[string][]byte, result2 error) {
	fake.getStateMetadataMutex.Lock()
	defer fake.getStateMetadataMutex.Unlock()
	fake.GetStateMetadataStub = nil
	fake.getStateMetadataReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateMetadataReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
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

func (fake *QueryExecutor) GetStateMultipleKeys(arg1 string, arg2 []string) ([][]byte, error) {
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

func (fake *QueryExecutor) GetStateMultipleKeysCallCount() int {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	return len(fake.getStateMultipleKeysArgsForCall)
}

func (fake *QueryExecutor) GetStateMultipleKeysCalls(stub func(string, []string) ([][]byte, error)) {
	fake.getStateMultipleKeysMutex.Lock()
	defer fake.getStateMultipleKeysMutex.Unlock()
	fake.GetStateMultipleKeysStub = stub
}

func (fake *QueryExecutor) GetStateMultipleKeysArgsForCall(i int) (string, []string) {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	argsForCall := fake.getStateMultipleKeysArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *QueryExecutor) GetStateMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.getStateMultipleKeysMutex.Lock()
	defer fake.getStateMultipleKeysMutex.Unlock()
	fake.GetStateMultipleKeysStub = nil
	fake.getStateMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
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

func (fake *QueryExecutor) GetStateRangeScanIterator(arg1 string, arg2 string, arg3 string) (ledgera.ResultsIterator, error) {
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

func (fake *QueryExecutor) GetStateRangeScanIteratorCallCount() int {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorArgsForCall)
}

func (fake *QueryExecutor) GetStateRangeScanIteratorCalls(stub func(string, string, string) (ledgera.ResultsIterator, error)) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = stub
}

func (fake *QueryExecutor) GetStateRangeScanIteratorArgsForCall(i int) (string, string, string) {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	argsForCall := fake.getStateRangeScanIteratorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *QueryExecutor) GetStateRangeScanIteratorReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = nil
	fake.getStateRangeScanIteratorReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateRangeScanIteratorReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
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

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadata(arg1 string, arg2 string, arg3 string, arg4 map[string]interface{}) (ledger.QueryResultsIterator, error) {
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

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataCallCount() int {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorWithMetadataArgsForCall)
}

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataCalls(stub func(string, string, string, map[string]interface{}) (ledger.QueryResultsIterator, error)) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	fake.GetStateRangeScanIteratorWithMetadataStub = stub
}

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataArgsForCall(i int) (string, string, string, map[string]interface{}) {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	argsForCall := fake.getStateRangeScanIteratorWithMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataReturns(result1 ledger.QueryResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	fake.GetStateRangeScanIteratorWithMetadataStub = nil
	fake.getStateRangeScanIteratorWithMetadataReturns = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataReturnsOnCall(i int, result1 ledger.QueryResultsIterator, result2 error) {
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

func (fake *QueryExecutor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
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
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *QueryExecutor) recordInvocation(key string, args []interface{}) {
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

var _ ledger.QueryExecutor = new(QueryExecutor)
