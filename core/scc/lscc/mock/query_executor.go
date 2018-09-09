
package mock

import (
	"sync"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger"
)

type QueryExecutor struct {
	GetStateStub        func(namespace string, key string) ([]byte, error)
	getStateMutex       sync.RWMutex
	getStateArgsForCall []struct {
		namespace string
		key       string
	}
	getStateReturns struct {
		result1 []byte
		result2 error
	}
	getStateReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetStateRangeScanIteratorStub        func(namespace string, startKey string, endKey string) (commonledger.ResultsIterator, error)
	getStateRangeScanIteratorMutex       sync.RWMutex
	getStateRangeScanIteratorArgsForCall []struct {
		namespace string
		startKey  string
		endKey    string
	}
	getStateRangeScanIteratorReturns struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	getStateRangeScanIteratorReturnsOnCall map[int]struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	GetStateMetadataStub        func(namespace, key string) (map[string][]byte, error)
	getStateMetadataMutex       sync.RWMutex
	getStateMetadataArgsForCall []struct {
		namespace string
		key       string
	}
	getStateMetadataReturns struct {
		result1 map[string][]byte
		result2 error
	}
	getStateMetadataReturnsOnCall map[int]struct {
		result1 map[string][]byte
		result2 error
	}
	GetStateMultipleKeysStub        func(namespace string, keys []string) ([][]byte, error)
	getStateMultipleKeysMutex       sync.RWMutex
	getStateMultipleKeysArgsForCall []struct {
		namespace string
		keys      []string
	}
	getStateMultipleKeysReturns struct {
		result1 [][]byte
		result2 error
	}
	getStateMultipleKeysReturnsOnCall map[int]struct {
		result1 [][]byte
		result2 error
	}
	GetStateRangeScanIteratorWithMetadataStub        func(namespace string, startKey, endKey string, metadata map[string]interface{}) (ledger.QueryResultsIterator, error)
	getStateRangeScanIteratorWithMetadataMutex       sync.RWMutex
	getStateRangeScanIteratorWithMetadataArgsForCall []struct {
		namespace string
		startKey  string
		endKey    string
		metadata  map[string]interface{}
	}
	getStateRangeScanIteratorWithMetadataReturns struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}
	getStateRangeScanIteratorWithMetadataReturnsOnCall map[int]struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}
	ExecuteQueryStub        func(namespace, query string) (commonledger.ResultsIterator, error)
	executeQueryMutex       sync.RWMutex
	executeQueryArgsForCall []struct {
		namespace string
		query     string
	}
	executeQueryReturns struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	executeQueryReturnsOnCall map[int]struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	ExecuteQueryWithMetadataStub        func(namespace, query string, metadata map[string]interface{}) (ledger.QueryResultsIterator, error)
	executeQueryWithMetadataMutex       sync.RWMutex
	executeQueryWithMetadataArgsForCall []struct {
		namespace string
		query     string
		metadata  map[string]interface{}
	}
	executeQueryWithMetadataReturns struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}
	executeQueryWithMetadataReturnsOnCall map[int]struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}
	GetPrivateDataStub        func(namespace, collection, key string) ([]byte, error)
	getPrivateDataMutex       sync.RWMutex
	getPrivateDataArgsForCall []struct {
		namespace  string
		collection string
		key        string
	}
	getPrivateDataReturns struct {
		result1 []byte
		result2 error
	}
	getPrivateDataReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetPrivateDataMetadataStub        func(namespace, collection, key string) (map[string][]byte, error)
	getPrivateDataMetadataMutex       sync.RWMutex
	getPrivateDataMetadataArgsForCall []struct {
		namespace  string
		collection string
		key        string
	}
	getPrivateDataMetadataReturns struct {
		result1 map[string][]byte
		result2 error
	}
	getPrivateDataMetadataReturnsOnCall map[int]struct {
		result1 map[string][]byte
		result2 error
	}
	GetPrivateDataMultipleKeysStub        func(namespace, collection string, keys []string) ([][]byte, error)
	getPrivateDataMultipleKeysMutex       sync.RWMutex
	getPrivateDataMultipleKeysArgsForCall []struct {
		namespace  string
		collection string
		keys       []string
	}
	getPrivateDataMultipleKeysReturns struct {
		result1 [][]byte
		result2 error
	}
	getPrivateDataMultipleKeysReturnsOnCall map[int]struct {
		result1 [][]byte
		result2 error
	}
	GetPrivateDataRangeScanIteratorStub        func(namespace, collection, startKey, endKey string) (commonledger.ResultsIterator, error)
	getPrivateDataRangeScanIteratorMutex       sync.RWMutex
	getPrivateDataRangeScanIteratorArgsForCall []struct {
		namespace  string
		collection string
		startKey   string
		endKey     string
	}
	getPrivateDataRangeScanIteratorReturns struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	getPrivateDataRangeScanIteratorReturnsOnCall map[int]struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	ExecuteQueryOnPrivateDataStub        func(namespace, collection, query string) (commonledger.ResultsIterator, error)
	executeQueryOnPrivateDataMutex       sync.RWMutex
	executeQueryOnPrivateDataArgsForCall []struct {
		namespace  string
		collection string
		query      string
	}
	executeQueryOnPrivateDataReturns struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	executeQueryOnPrivateDataReturnsOnCall map[int]struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	DoneStub         func()
	doneMutex        sync.RWMutex
	doneArgsForCall  []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *QueryExecutor) GetState(namespace string, key string) ([]byte, error) {
	fake.getStateMutex.Lock()
	ret, specificReturn := fake.getStateReturnsOnCall[len(fake.getStateArgsForCall)]
	fake.getStateArgsForCall = append(fake.getStateArgsForCall, struct {
		namespace string
		key       string
	}{namespace, key})
	fake.recordInvocation("GetState", []interface{}{namespace, key})
	fake.getStateMutex.Unlock()
	if fake.GetStateStub != nil {
		return fake.GetStateStub(namespace, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateReturns.result1, fake.getStateReturns.result2
}

func (fake *QueryExecutor) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *QueryExecutor) GetStateArgsForCall(i int) (string, string) {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return fake.getStateArgsForCall[i].namespace, fake.getStateArgsForCall[i].key
}

func (fake *QueryExecutor) GetStateReturns(result1 []byte, result2 error) {
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *QueryExecutor) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (commonledger.ResultsIterator, error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	ret, specificReturn := fake.getStateRangeScanIteratorReturnsOnCall[len(fake.getStateRangeScanIteratorArgsForCall)]
	fake.getStateRangeScanIteratorArgsForCall = append(fake.getStateRangeScanIteratorArgsForCall, struct {
		namespace string
		startKey  string
		endKey    string
	}{namespace, startKey, endKey})
	fake.recordInvocation("GetStateRangeScanIterator", []interface{}{namespace, startKey, endKey})
	fake.getStateRangeScanIteratorMutex.Unlock()
	if fake.GetStateRangeScanIteratorStub != nil {
		return fake.GetStateRangeScanIteratorStub(namespace, startKey, endKey)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateRangeScanIteratorReturns.result1, fake.getStateRangeScanIteratorReturns.result2
}

func (fake *QueryExecutor) GetStateRangeScanIteratorCallCount() int {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorArgsForCall)
}

func (fake *QueryExecutor) GetStateRangeScanIteratorArgsForCall(i int) (string, string, string) {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	return fake.getStateRangeScanIteratorArgsForCall[i].namespace, fake.getStateRangeScanIteratorArgsForCall[i].startKey, fake.getStateRangeScanIteratorArgsForCall[i].endKey
}

func (fake *QueryExecutor) GetStateRangeScanIteratorReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.GetStateRangeScanIteratorStub = nil
	fake.getStateRangeScanIteratorReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateRangeScanIteratorReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
	fake.GetStateRangeScanIteratorStub = nil
	if fake.getStateRangeScanIteratorReturnsOnCall == nil {
		fake.getStateRangeScanIteratorReturnsOnCall = make(map[int]struct {
			result1 commonledger.ResultsIterator
			result2 error
		})
	}
	fake.getStateRangeScanIteratorReturnsOnCall[i] = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateMetadata(namespace string, key string) (map[string][]byte, error) {
	fake.getStateMetadataMutex.Lock()
	ret, specificReturn := fake.getStateMetadataReturnsOnCall[len(fake.getStateMetadataArgsForCall)]
	fake.getStateMetadataArgsForCall = append(fake.getStateMetadataArgsForCall, struct {
		namespace string
		key       string
	}{namespace, key})
	fake.recordInvocation("GetStateMetadata", []interface{}{namespace, key})
	fake.getStateMetadataMutex.Unlock()
	if fake.GetStateMetadataStub != nil {
		return fake.GetStateMetadataStub(namespace, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateMetadataReturns.result1, fake.getStateMetadataReturns.result2
}

func (fake *QueryExecutor) GetStateMetadataCallCount() int {
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	return len(fake.getStateMetadataArgsForCall)
}

func (fake *QueryExecutor) GetStateMetadataArgsForCall(i int) (string, string) {
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	return fake.getStateMetadataArgsForCall[i].namespace, fake.getStateMetadataArgsForCall[i].key
}

func (fake *QueryExecutor) GetStateMetadataReturns(result1 map[string][]byte, result2 error) {
	fake.GetStateMetadataStub = nil
	fake.getStateMetadataReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateMetadataReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
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

func (fake *QueryExecutor) GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error) {
	var keysCopy []string
	if keys != nil {
		keysCopy = make([]string, len(keys))
		copy(keysCopy, keys)
	}
	fake.getStateMultipleKeysMutex.Lock()
	ret, specificReturn := fake.getStateMultipleKeysReturnsOnCall[len(fake.getStateMultipleKeysArgsForCall)]
	fake.getStateMultipleKeysArgsForCall = append(fake.getStateMultipleKeysArgsForCall, struct {
		namespace string
		keys      []string
	}{namespace, keysCopy})
	fake.recordInvocation("GetStateMultipleKeys", []interface{}{namespace, keysCopy})
	fake.getStateMultipleKeysMutex.Unlock()
	if fake.GetStateMultipleKeysStub != nil {
		return fake.GetStateMultipleKeysStub(namespace, keys)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateMultipleKeysReturns.result1, fake.getStateMultipleKeysReturns.result2
}

func (fake *QueryExecutor) GetStateMultipleKeysCallCount() int {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	return len(fake.getStateMultipleKeysArgsForCall)
}

func (fake *QueryExecutor) GetStateMultipleKeysArgsForCall(i int) (string, []string) {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	return fake.getStateMultipleKeysArgsForCall[i].namespace, fake.getStateMultipleKeysArgsForCall[i].keys
}

func (fake *QueryExecutor) GetStateMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.GetStateMultipleKeysStub = nil
	fake.getStateMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
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

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (ledger.QueryResultsIterator, error) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	ret, specificReturn := fake.getStateRangeScanIteratorWithMetadataReturnsOnCall[len(fake.getStateRangeScanIteratorWithMetadataArgsForCall)]
	fake.getStateRangeScanIteratorWithMetadataArgsForCall = append(fake.getStateRangeScanIteratorWithMetadataArgsForCall, struct {
		namespace string
		startKey  string
		endKey    string
		metadata  map[string]interface{}
	}{namespace, startKey, endKey, metadata})
	fake.recordInvocation("GetStateRangeScanIteratorWithMetadata", []interface{}{namespace, startKey, endKey, metadata})
	fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	if fake.GetStateRangeScanIteratorWithMetadataStub != nil {
		return fake.GetStateRangeScanIteratorWithMetadataStub(namespace, startKey, endKey, metadata)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateRangeScanIteratorWithMetadataReturns.result1, fake.getStateRangeScanIteratorWithMetadataReturns.result2
}

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataCallCount() int {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorWithMetadataArgsForCall)
}

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataArgsForCall(i int) (string, string, string, map[string]interface{}) {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	return fake.getStateRangeScanIteratorWithMetadataArgsForCall[i].namespace, fake.getStateRangeScanIteratorWithMetadataArgsForCall[i].startKey, fake.getStateRangeScanIteratorWithMetadataArgsForCall[i].endKey, fake.getStateRangeScanIteratorWithMetadataArgsForCall[i].metadata
}

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataReturns(result1 ledger.QueryResultsIterator, result2 error) {
	fake.GetStateRangeScanIteratorWithMetadataStub = nil
	fake.getStateRangeScanIteratorWithMetadataReturns = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetStateRangeScanIteratorWithMetadataReturnsOnCall(i int, result1 ledger.QueryResultsIterator, result2 error) {
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

func (fake *QueryExecutor) ExecuteQuery(namespace string, query string) (commonledger.ResultsIterator, error) {
	fake.executeQueryMutex.Lock()
	ret, specificReturn := fake.executeQueryReturnsOnCall[len(fake.executeQueryArgsForCall)]
	fake.executeQueryArgsForCall = append(fake.executeQueryArgsForCall, struct {
		namespace string
		query     string
	}{namespace, query})
	fake.recordInvocation("ExecuteQuery", []interface{}{namespace, query})
	fake.executeQueryMutex.Unlock()
	if fake.ExecuteQueryStub != nil {
		return fake.ExecuteQueryStub(namespace, query)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.executeQueryReturns.result1, fake.executeQueryReturns.result2
}

func (fake *QueryExecutor) ExecuteQueryCallCount() int {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	return len(fake.executeQueryArgsForCall)
}

func (fake *QueryExecutor) ExecuteQueryArgsForCall(i int) (string, string) {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	return fake.executeQueryArgsForCall[i].namespace, fake.executeQueryArgsForCall[i].query
}

func (fake *QueryExecutor) ExecuteQueryReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.ExecuteQueryStub = nil
	fake.executeQueryReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) ExecuteQueryReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
	fake.ExecuteQueryStub = nil
	if fake.executeQueryReturnsOnCall == nil {
		fake.executeQueryReturnsOnCall = make(map[int]struct {
			result1 commonledger.ResultsIterator
			result2 error
		})
	}
	fake.executeQueryReturnsOnCall[i] = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) ExecuteQueryWithMetadata(namespace string, query string, metadata map[string]interface{}) (ledger.QueryResultsIterator, error) {
	fake.executeQueryWithMetadataMutex.Lock()
	ret, specificReturn := fake.executeQueryWithMetadataReturnsOnCall[len(fake.executeQueryWithMetadataArgsForCall)]
	fake.executeQueryWithMetadataArgsForCall = append(fake.executeQueryWithMetadataArgsForCall, struct {
		namespace string
		query     string
		metadata  map[string]interface{}
	}{namespace, query, metadata})
	fake.recordInvocation("ExecuteQueryWithMetadata", []interface{}{namespace, query, metadata})
	fake.executeQueryWithMetadataMutex.Unlock()
	if fake.ExecuteQueryWithMetadataStub != nil {
		return fake.ExecuteQueryWithMetadataStub(namespace, query, metadata)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.executeQueryWithMetadataReturns.result1, fake.executeQueryWithMetadataReturns.result2
}

func (fake *QueryExecutor) ExecuteQueryWithMetadataCallCount() int {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	return len(fake.executeQueryWithMetadataArgsForCall)
}

func (fake *QueryExecutor) ExecuteQueryWithMetadataArgsForCall(i int) (string, string, map[string]interface{}) {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	return fake.executeQueryWithMetadataArgsForCall[i].namespace, fake.executeQueryWithMetadataArgsForCall[i].query, fake.executeQueryWithMetadataArgsForCall[i].metadata
}

func (fake *QueryExecutor) ExecuteQueryWithMetadataReturns(result1 ledger.QueryResultsIterator, result2 error) {
	fake.ExecuteQueryWithMetadataStub = nil
	fake.executeQueryWithMetadataReturns = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) ExecuteQueryWithMetadataReturnsOnCall(i int, result1 ledger.QueryResultsIterator, result2 error) {
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

func (fake *QueryExecutor) GetPrivateData(namespace string, collection string, key string) ([]byte, error) {
	fake.getPrivateDataMutex.Lock()
	ret, specificReturn := fake.getPrivateDataReturnsOnCall[len(fake.getPrivateDataArgsForCall)]
	fake.getPrivateDataArgsForCall = append(fake.getPrivateDataArgsForCall, struct {
		namespace  string
		collection string
		key        string
	}{namespace, collection, key})
	fake.recordInvocation("GetPrivateData", []interface{}{namespace, collection, key})
	fake.getPrivateDataMutex.Unlock()
	if fake.GetPrivateDataStub != nil {
		return fake.GetPrivateDataStub(namespace, collection, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataReturns.result1, fake.getPrivateDataReturns.result2
}

func (fake *QueryExecutor) GetPrivateDataCallCount() int {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	return len(fake.getPrivateDataArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	return fake.getPrivateDataArgsForCall[i].namespace, fake.getPrivateDataArgsForCall[i].collection, fake.getPrivateDataArgsForCall[i].key
}

func (fake *QueryExecutor) GetPrivateDataReturns(result1 []byte, result2 error) {
	fake.GetPrivateDataStub = nil
	fake.getPrivateDataReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *QueryExecutor) GetPrivateDataMetadata(namespace string, collection string, key string) (map[string][]byte, error) {
	fake.getPrivateDataMetadataMutex.Lock()
	ret, specificReturn := fake.getPrivateDataMetadataReturnsOnCall[len(fake.getPrivateDataMetadataArgsForCall)]
	fake.getPrivateDataMetadataArgsForCall = append(fake.getPrivateDataMetadataArgsForCall, struct {
		namespace  string
		collection string
		key        string
	}{namespace, collection, key})
	fake.recordInvocation("GetPrivateDataMetadata", []interface{}{namespace, collection, key})
	fake.getPrivateDataMetadataMutex.Unlock()
	if fake.GetPrivateDataMetadataStub != nil {
		return fake.GetPrivateDataMetadataStub(namespace, collection, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataMetadataReturns.result1, fake.getPrivateDataMetadataReturns.result2
}

func (fake *QueryExecutor) GetPrivateDataMetadataCallCount() int {
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	return len(fake.getPrivateDataMetadataArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataMetadataArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	return fake.getPrivateDataMetadataArgsForCall[i].namespace, fake.getPrivateDataMetadataArgsForCall[i].collection, fake.getPrivateDataMetadataArgsForCall[i].key
}

func (fake *QueryExecutor) GetPrivateDataMetadataReturns(result1 map[string][]byte, result2 error) {
	fake.GetPrivateDataMetadataStub = nil
	fake.getPrivateDataMetadataReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataMetadataReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
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

func (fake *QueryExecutor) GetPrivateDataMultipleKeys(namespace string, collection string, keys []string) ([][]byte, error) {
	var keysCopy []string
	if keys != nil {
		keysCopy = make([]string, len(keys))
		copy(keysCopy, keys)
	}
	fake.getPrivateDataMultipleKeysMutex.Lock()
	ret, specificReturn := fake.getPrivateDataMultipleKeysReturnsOnCall[len(fake.getPrivateDataMultipleKeysArgsForCall)]
	fake.getPrivateDataMultipleKeysArgsForCall = append(fake.getPrivateDataMultipleKeysArgsForCall, struct {
		namespace  string
		collection string
		keys       []string
	}{namespace, collection, keysCopy})
	fake.recordInvocation("GetPrivateDataMultipleKeys", []interface{}{namespace, collection, keysCopy})
	fake.getPrivateDataMultipleKeysMutex.Unlock()
	if fake.GetPrivateDataMultipleKeysStub != nil {
		return fake.GetPrivateDataMultipleKeysStub(namespace, collection, keys)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataMultipleKeysReturns.result1, fake.getPrivateDataMultipleKeysReturns.result2
}

func (fake *QueryExecutor) GetPrivateDataMultipleKeysCallCount() int {
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	return len(fake.getPrivateDataMultipleKeysArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataMultipleKeysArgsForCall(i int) (string, string, []string) {
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	return fake.getPrivateDataMultipleKeysArgsForCall[i].namespace, fake.getPrivateDataMultipleKeysArgsForCall[i].collection, fake.getPrivateDataMultipleKeysArgsForCall[i].keys
}

func (fake *QueryExecutor) GetPrivateDataMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.GetPrivateDataMultipleKeysStub = nil
	fake.getPrivateDataMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
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

func (fake *QueryExecutor) GetPrivateDataRangeScanIterator(namespace string, collection string, startKey string, endKey string) (commonledger.ResultsIterator, error) {
	fake.getPrivateDataRangeScanIteratorMutex.Lock()
	ret, specificReturn := fake.getPrivateDataRangeScanIteratorReturnsOnCall[len(fake.getPrivateDataRangeScanIteratorArgsForCall)]
	fake.getPrivateDataRangeScanIteratorArgsForCall = append(fake.getPrivateDataRangeScanIteratorArgsForCall, struct {
		namespace  string
		collection string
		startKey   string
		endKey     string
	}{namespace, collection, startKey, endKey})
	fake.recordInvocation("GetPrivateDataRangeScanIterator", []interface{}{namespace, collection, startKey, endKey})
	fake.getPrivateDataRangeScanIteratorMutex.Unlock()
	if fake.GetPrivateDataRangeScanIteratorStub != nil {
		return fake.GetPrivateDataRangeScanIteratorStub(namespace, collection, startKey, endKey)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataRangeScanIteratorReturns.result1, fake.getPrivateDataRangeScanIteratorReturns.result2
}

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorCallCount() int {
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	return len(fake.getPrivateDataRangeScanIteratorArgsForCall)
}

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorArgsForCall(i int) (string, string, string, string) {
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	return fake.getPrivateDataRangeScanIteratorArgsForCall[i].namespace, fake.getPrivateDataRangeScanIteratorArgsForCall[i].collection, fake.getPrivateDataRangeScanIteratorArgsForCall[i].startKey, fake.getPrivateDataRangeScanIteratorArgsForCall[i].endKey
}

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.GetPrivateDataRangeScanIteratorStub = nil
	fake.getPrivateDataRangeScanIteratorReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) GetPrivateDataRangeScanIteratorReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
	fake.GetPrivateDataRangeScanIteratorStub = nil
	if fake.getPrivateDataRangeScanIteratorReturnsOnCall == nil {
		fake.getPrivateDataRangeScanIteratorReturnsOnCall = make(map[int]struct {
			result1 commonledger.ResultsIterator
			result2 error
		})
	}
	fake.getPrivateDataRangeScanIteratorReturnsOnCall[i] = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateData(namespace string, collection string, query string) (commonledger.ResultsIterator, error) {
	fake.executeQueryOnPrivateDataMutex.Lock()
	ret, specificReturn := fake.executeQueryOnPrivateDataReturnsOnCall[len(fake.executeQueryOnPrivateDataArgsForCall)]
	fake.executeQueryOnPrivateDataArgsForCall = append(fake.executeQueryOnPrivateDataArgsForCall, struct {
		namespace  string
		collection string
		query      string
	}{namespace, collection, query})
	fake.recordInvocation("ExecuteQueryOnPrivateData", []interface{}{namespace, collection, query})
	fake.executeQueryOnPrivateDataMutex.Unlock()
	if fake.ExecuteQueryOnPrivateDataStub != nil {
		return fake.ExecuteQueryOnPrivateDataStub(namespace, collection, query)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.executeQueryOnPrivateDataReturns.result1, fake.executeQueryOnPrivateDataReturns.result2
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataCallCount() int {
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	return len(fake.executeQueryOnPrivateDataArgsForCall)
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataArgsForCall(i int) (string, string, string) {
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	return fake.executeQueryOnPrivateDataArgsForCall[i].namespace, fake.executeQueryOnPrivateDataArgsForCall[i].collection, fake.executeQueryOnPrivateDataArgsForCall[i].query
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.ExecuteQueryOnPrivateDataStub = nil
	fake.executeQueryOnPrivateDataReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) ExecuteQueryOnPrivateDataReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
	fake.ExecuteQueryOnPrivateDataStub = nil
	if fake.executeQueryOnPrivateDataReturnsOnCall == nil {
		fake.executeQueryOnPrivateDataReturnsOnCall = make(map[int]struct {
			result1 commonledger.ResultsIterator
			result2 error
		})
	}
	fake.executeQueryOnPrivateDataReturnsOnCall[i] = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *QueryExecutor) Done() {
	fake.doneMutex.Lock()
	fake.doneArgsForCall = append(fake.doneArgsForCall, struct{}{})
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

func (fake *QueryExecutor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
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
