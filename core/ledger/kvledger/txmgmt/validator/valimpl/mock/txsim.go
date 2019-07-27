
package mock

import (
	"sync"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger"
)

type TxSimulator struct {
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
	GetPrivateDataHashStub        func(namespace, collection, key string) ([]byte, error)
	getPrivateDataHashMutex       sync.RWMutex
	getPrivateDataHashArgsForCall []struct {
		namespace  string
		collection string
		key        string
	}
	getPrivateDataHashReturns struct {
		result1 []byte
		result2 error
	}
	getPrivateDataHashReturnsOnCall map[int]struct {
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
	GetPrivateDataMetadataByHashStub        func(namespace, collection string, keyhash []byte) (map[string][]byte, error)
	getPrivateDataMetadataByHashMutex       sync.RWMutex
	getPrivateDataMetadataByHashArgsForCall []struct {
		namespace  string
		collection string
		keyhash    []byte
	}
	getPrivateDataMetadataByHashReturns struct {
		result1 map[string][]byte
		result2 error
	}
	getPrivateDataMetadataByHashReturnsOnCall map[int]struct {
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
	DoneStub            func()
	doneMutex           sync.RWMutex
	doneArgsForCall     []struct{}
	SetStateStub        func(namespace string, key string, value []byte) error
	setStateMutex       sync.RWMutex
	setStateArgsForCall []struct {
		namespace string
		key       string
		value     []byte
	}
	setStateReturns struct {
		result1 error
	}
	setStateReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStateStub        func(namespace string, key string) error
	deleteStateMutex       sync.RWMutex
	deleteStateArgsForCall []struct {
		namespace string
		key       string
	}
	deleteStateReturns struct {
		result1 error
	}
	deleteStateReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateMultipleKeysStub        func(namespace string, kvs map[string][]byte) error
	setStateMultipleKeysMutex       sync.RWMutex
	setStateMultipleKeysArgsForCall []struct {
		namespace string
		kvs       map[string][]byte
	}
	setStateMultipleKeysReturns struct {
		result1 error
	}
	setStateMultipleKeysReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateMetadataStub        func(namespace, key string, metadata map[string][]byte) error
	setStateMetadataMutex       sync.RWMutex
	setStateMetadataArgsForCall []struct {
		namespace string
		key       string
		metadata  map[string][]byte
	}
	setStateMetadataReturns struct {
		result1 error
	}
	setStateMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStateMetadataStub        func(namespace, key string) error
	deleteStateMetadataMutex       sync.RWMutex
	deleteStateMetadataArgsForCall []struct {
		namespace string
		key       string
	}
	deleteStateMetadataReturns struct {
		result1 error
	}
	deleteStateMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	ExecuteUpdateStub        func(query string) error
	executeUpdateMutex       sync.RWMutex
	executeUpdateArgsForCall []struct {
		query string
	}
	executeUpdateReturns struct {
		result1 error
	}
	executeUpdateReturnsOnCall map[int]struct {
		result1 error
	}
	SetPrivateDataStub        func(namespace, collection, key string, value []byte) error
	setPrivateDataMutex       sync.RWMutex
	setPrivateDataArgsForCall []struct {
		namespace  string
		collection string
		key        string
		value      []byte
	}
	setPrivateDataReturns struct {
		result1 error
	}
	setPrivateDataReturnsOnCall map[int]struct {
		result1 error
	}
	SetPrivateDataMultipleKeysStub        func(namespace, collection string, kvs map[string][]byte) error
	setPrivateDataMultipleKeysMutex       sync.RWMutex
	setPrivateDataMultipleKeysArgsForCall []struct {
		namespace  string
		collection string
		kvs        map[string][]byte
	}
	setPrivateDataMultipleKeysReturns struct {
		result1 error
	}
	setPrivateDataMultipleKeysReturnsOnCall map[int]struct {
		result1 error
	}
	DeletePrivateDataStub        func(namespace, collection, key string) error
	deletePrivateDataMutex       sync.RWMutex
	deletePrivateDataArgsForCall []struct {
		namespace  string
		collection string
		key        string
	}
	deletePrivateDataReturns struct {
		result1 error
	}
	deletePrivateDataReturnsOnCall map[int]struct {
		result1 error
	}
	SetPrivateDataMetadataStub        func(namespace, collection, key string, metadata map[string][]byte) error
	setPrivateDataMetadataMutex       sync.RWMutex
	setPrivateDataMetadataArgsForCall []struct {
		namespace  string
		collection string
		key        string
		metadata   map[string][]byte
	}
	setPrivateDataMetadataReturns struct {
		result1 error
	}
	setPrivateDataMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	DeletePrivateDataMetadataStub        func(namespace, collection, key string) error
	deletePrivateDataMetadataMutex       sync.RWMutex
	deletePrivateDataMetadataArgsForCall []struct {
		namespace  string
		collection string
		key        string
	}
	deletePrivateDataMetadataReturns struct {
		result1 error
	}
	deletePrivateDataMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	GetTxSimulationResultsStub        func() (*ledger.TxSimulationResults, error)
	getTxSimulationResultsMutex       sync.RWMutex
	getTxSimulationResultsArgsForCall []struct{}
	getTxSimulationResultsReturns     struct {
		result1 *ledger.TxSimulationResults
		result2 error
	}
	getTxSimulationResultsReturnsOnCall map[int]struct {
		result1 *ledger.TxSimulationResults
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TxSimulator) GetState(namespace string, key string) ([]byte, error) {
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

func (fake *TxSimulator) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *TxSimulator) GetStateArgsForCall(i int) (string, string) {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return fake.getStateArgsForCall[i].namespace, fake.getStateArgsForCall[i].key
}

func (fake *TxSimulator) GetStateReturns(result1 []byte, result2 error) {
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *TxSimulator) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (commonledger.ResultsIterator, error) {
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

func (fake *TxSimulator) GetStateRangeScanIteratorCallCount() int {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorArgsForCall)
}

func (fake *TxSimulator) GetStateRangeScanIteratorArgsForCall(i int) (string, string, string) {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	return fake.getStateRangeScanIteratorArgsForCall[i].namespace, fake.getStateRangeScanIteratorArgsForCall[i].startKey, fake.getStateRangeScanIteratorArgsForCall[i].endKey
}

func (fake *TxSimulator) GetStateRangeScanIteratorReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.GetStateRangeScanIteratorStub = nil
	fake.getStateRangeScanIteratorReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateRangeScanIteratorReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
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

func (fake *TxSimulator) GetStateMetadata(namespace string, key string) (map[string][]byte, error) {
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

func (fake *TxSimulator) GetStateMetadataCallCount() int {
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	return len(fake.getStateMetadataArgsForCall)
}

func (fake *TxSimulator) GetStateMetadataArgsForCall(i int) (string, string) {
	fake.getStateMetadataMutex.RLock()
	defer fake.getStateMetadataMutex.RUnlock()
	return fake.getStateMetadataArgsForCall[i].namespace, fake.getStateMetadataArgsForCall[i].key
}

func (fake *TxSimulator) GetStateMetadataReturns(result1 map[string][]byte, result2 error) {
	fake.GetStateMetadataStub = nil
	fake.getStateMetadataReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateMetadataReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
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

func (fake *TxSimulator) GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error) {
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

func (fake *TxSimulator) GetStateMultipleKeysCallCount() int {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	return len(fake.getStateMultipleKeysArgsForCall)
}

func (fake *TxSimulator) GetStateMultipleKeysArgsForCall(i int) (string, []string) {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	return fake.getStateMultipleKeysArgsForCall[i].namespace, fake.getStateMultipleKeysArgsForCall[i].keys
}

func (fake *TxSimulator) GetStateMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.GetStateMultipleKeysStub = nil
	fake.getStateMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
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

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (ledger.QueryResultsIterator, error) {
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

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataCallCount() int {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorWithMetadataArgsForCall)
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataArgsForCall(i int) (string, string, string, map[string]interface{}) {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	return fake.getStateRangeScanIteratorWithMetadataArgsForCall[i].namespace, fake.getStateRangeScanIteratorWithMetadataArgsForCall[i].startKey, fake.getStateRangeScanIteratorWithMetadataArgsForCall[i].endKey, fake.getStateRangeScanIteratorWithMetadataArgsForCall[i].metadata
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataReturns(result1 ledger.QueryResultsIterator, result2 error) {
	fake.GetStateRangeScanIteratorWithMetadataStub = nil
	fake.getStateRangeScanIteratorWithMetadataReturns = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetStateRangeScanIteratorWithMetadataReturnsOnCall(i int, result1 ledger.QueryResultsIterator, result2 error) {
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

func (fake *TxSimulator) ExecuteQuery(namespace string, query string) (commonledger.ResultsIterator, error) {
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

func (fake *TxSimulator) ExecuteQueryCallCount() int {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	return len(fake.executeQueryArgsForCall)
}

func (fake *TxSimulator) ExecuteQueryArgsForCall(i int) (string, string) {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	return fake.executeQueryArgsForCall[i].namespace, fake.executeQueryArgsForCall[i].query
}

func (fake *TxSimulator) ExecuteQueryReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.ExecuteQueryStub = nil
	fake.executeQueryReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteQueryReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
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

func (fake *TxSimulator) ExecuteQueryWithMetadata(namespace string, query string, metadata map[string]interface{}) (ledger.QueryResultsIterator, error) {
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

func (fake *TxSimulator) ExecuteQueryWithMetadataCallCount() int {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	return len(fake.executeQueryWithMetadataArgsForCall)
}

func (fake *TxSimulator) ExecuteQueryWithMetadataArgsForCall(i int) (string, string, map[string]interface{}) {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	return fake.executeQueryWithMetadataArgsForCall[i].namespace, fake.executeQueryWithMetadataArgsForCall[i].query, fake.executeQueryWithMetadataArgsForCall[i].metadata
}

func (fake *TxSimulator) ExecuteQueryWithMetadataReturns(result1 ledger.QueryResultsIterator, result2 error) {
	fake.ExecuteQueryWithMetadataStub = nil
	fake.executeQueryWithMetadataReturns = struct {
		result1 ledger.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteQueryWithMetadataReturnsOnCall(i int, result1 ledger.QueryResultsIterator, result2 error) {
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

func (fake *TxSimulator) GetPrivateData(namespace string, collection string, key string) ([]byte, error) {
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

func (fake *TxSimulator) GetPrivateDataCallCount() int {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	return len(fake.getPrivateDataArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	return fake.getPrivateDataArgsForCall[i].namespace, fake.getPrivateDataArgsForCall[i].collection, fake.getPrivateDataArgsForCall[i].key
}

func (fake *TxSimulator) GetPrivateDataReturns(result1 []byte, result2 error) {
	fake.GetPrivateDataStub = nil
	fake.getPrivateDataReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *TxSimulator) GetPrivateDataHash(namespace string, collection string, key string) ([]byte, error) {
	fake.getPrivateDataHashMutex.Lock()
	ret, specificReturn := fake.getPrivateDataHashReturnsOnCall[len(fake.getPrivateDataHashArgsForCall)]
	fake.getPrivateDataHashArgsForCall = append(fake.getPrivateDataHashArgsForCall, struct {
		namespace  string
		collection string
		key        string
	}{namespace, collection, key})
	fake.recordInvocation("GetPrivateDataHash", []interface{}{namespace, collection, key})
	fake.getPrivateDataHashMutex.Unlock()
	if fake.GetPrivateDataHashStub != nil {
		return fake.GetPrivateDataHashStub(namespace, collection, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataHashReturns.result1, fake.getPrivateDataHashReturns.result2
}

func (fake *TxSimulator) GetPrivateDataHashCallCount() int {
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	return len(fake.getPrivateDataHashArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataHashArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	return fake.getPrivateDataHashArgsForCall[i].namespace, fake.getPrivateDataHashArgsForCall[i].collection, fake.getPrivateDataHashArgsForCall[i].key
}

func (fake *TxSimulator) GetPrivateDataHashReturns(result1 []byte, result2 error) {
	fake.GetPrivateDataHashStub = nil
	fake.getPrivateDataHashReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataHashReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *TxSimulator) GetPrivateDataMetadata(namespace string, collection string, key string) (map[string][]byte, error) {
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

func (fake *TxSimulator) GetPrivateDataMetadataCallCount() int {
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	return len(fake.getPrivateDataMetadataArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataMetadataArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataMetadataMutex.RLock()
	defer fake.getPrivateDataMetadataMutex.RUnlock()
	return fake.getPrivateDataMetadataArgsForCall[i].namespace, fake.getPrivateDataMetadataArgsForCall[i].collection, fake.getPrivateDataMetadataArgsForCall[i].key
}

func (fake *TxSimulator) GetPrivateDataMetadataReturns(result1 map[string][]byte, result2 error) {
	fake.GetPrivateDataMetadataStub = nil
	fake.getPrivateDataMetadataReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMetadataReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
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

func (fake *TxSimulator) GetPrivateDataMetadataByHash(namespace string, collection string, keyhash []byte) (map[string][]byte, error) {
	var keyhashCopy []byte
	if keyhash != nil {
		keyhashCopy = make([]byte, len(keyhash))
		copy(keyhashCopy, keyhash)
	}
	fake.getPrivateDataMetadataByHashMutex.Lock()
	ret, specificReturn := fake.getPrivateDataMetadataByHashReturnsOnCall[len(fake.getPrivateDataMetadataByHashArgsForCall)]
	fake.getPrivateDataMetadataByHashArgsForCall = append(fake.getPrivateDataMetadataByHashArgsForCall, struct {
		namespace  string
		collection string
		keyhash    []byte
	}{namespace, collection, keyhashCopy})
	fake.recordInvocation("GetPrivateDataMetadataByHash", []interface{}{namespace, collection, keyhashCopy})
	fake.getPrivateDataMetadataByHashMutex.Unlock()
	if fake.GetPrivateDataMetadataByHashStub != nil {
		return fake.GetPrivateDataMetadataByHashStub(namespace, collection, keyhash)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataMetadataByHashReturns.result1, fake.getPrivateDataMetadataByHashReturns.result2
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashCallCount() int {
	fake.getPrivateDataMetadataByHashMutex.RLock()
	defer fake.getPrivateDataMetadataByHashMutex.RUnlock()
	return len(fake.getPrivateDataMetadataByHashArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashArgsForCall(i int) (string, string, []byte) {
	fake.getPrivateDataMetadataByHashMutex.RLock()
	defer fake.getPrivateDataMetadataByHashMutex.RUnlock()
	return fake.getPrivateDataMetadataByHashArgsForCall[i].namespace, fake.getPrivateDataMetadataByHashArgsForCall[i].collection, fake.getPrivateDataMetadataByHashArgsForCall[i].keyhash
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashReturns(result1 map[string][]byte, result2 error) {
	fake.GetPrivateDataMetadataByHashStub = nil
	fake.getPrivateDataMetadataByHashReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMetadataByHashReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
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

func (fake *TxSimulator) GetPrivateDataMultipleKeys(namespace string, collection string, keys []string) ([][]byte, error) {
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

func (fake *TxSimulator) GetPrivateDataMultipleKeysCallCount() int {
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	return len(fake.getPrivateDataMultipleKeysArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataMultipleKeysArgsForCall(i int) (string, string, []string) {
	fake.getPrivateDataMultipleKeysMutex.RLock()
	defer fake.getPrivateDataMultipleKeysMutex.RUnlock()
	return fake.getPrivateDataMultipleKeysArgsForCall[i].namespace, fake.getPrivateDataMultipleKeysArgsForCall[i].collection, fake.getPrivateDataMultipleKeysArgsForCall[i].keys
}

func (fake *TxSimulator) GetPrivateDataMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.GetPrivateDataMultipleKeysStub = nil
	fake.getPrivateDataMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
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

func (fake *TxSimulator) GetPrivateDataRangeScanIterator(namespace string, collection string, startKey string, endKey string) (commonledger.ResultsIterator, error) {
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

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorCallCount() int {
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	return len(fake.getPrivateDataRangeScanIteratorArgsForCall)
}

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorArgsForCall(i int) (string, string, string, string) {
	fake.getPrivateDataRangeScanIteratorMutex.RLock()
	defer fake.getPrivateDataRangeScanIteratorMutex.RUnlock()
	return fake.getPrivateDataRangeScanIteratorArgsForCall[i].namespace, fake.getPrivateDataRangeScanIteratorArgsForCall[i].collection, fake.getPrivateDataRangeScanIteratorArgsForCall[i].startKey, fake.getPrivateDataRangeScanIteratorArgsForCall[i].endKey
}

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.GetPrivateDataRangeScanIteratorStub = nil
	fake.getPrivateDataRangeScanIteratorReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetPrivateDataRangeScanIteratorReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
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

func (fake *TxSimulator) ExecuteQueryOnPrivateData(namespace string, collection string, query string) (commonledger.ResultsIterator, error) {
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

func (fake *TxSimulator) ExecuteQueryOnPrivateDataCallCount() int {
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	return len(fake.executeQueryOnPrivateDataArgsForCall)
}

func (fake *TxSimulator) ExecuteQueryOnPrivateDataArgsForCall(i int) (string, string, string) {
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	return fake.executeQueryOnPrivateDataArgsForCall[i].namespace, fake.executeQueryOnPrivateDataArgsForCall[i].collection, fake.executeQueryOnPrivateDataArgsForCall[i].query
}

func (fake *TxSimulator) ExecuteQueryOnPrivateDataReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.ExecuteQueryOnPrivateDataStub = nil
	fake.executeQueryOnPrivateDataReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) ExecuteQueryOnPrivateDataReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
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

func (fake *TxSimulator) Done() {
	fake.doneMutex.Lock()
	fake.doneArgsForCall = append(fake.doneArgsForCall, struct{}{})
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

func (fake *TxSimulator) SetState(namespace string, key string, value []byte) error {
	var valueCopy []byte
	if value != nil {
		valueCopy = make([]byte, len(value))
		copy(valueCopy, value)
	}
	fake.setStateMutex.Lock()
	ret, specificReturn := fake.setStateReturnsOnCall[len(fake.setStateArgsForCall)]
	fake.setStateArgsForCall = append(fake.setStateArgsForCall, struct {
		namespace string
		key       string
		value     []byte
	}{namespace, key, valueCopy})
	fake.recordInvocation("SetState", []interface{}{namespace, key, valueCopy})
	fake.setStateMutex.Unlock()
	if fake.SetStateStub != nil {
		return fake.SetStateStub(namespace, key, value)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setStateReturns.result1
}

func (fake *TxSimulator) SetStateCallCount() int {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return len(fake.setStateArgsForCall)
}

func (fake *TxSimulator) SetStateArgsForCall(i int) (string, string, []byte) {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return fake.setStateArgsForCall[i].namespace, fake.setStateArgsForCall[i].key, fake.setStateArgsForCall[i].value
}

func (fake *TxSimulator) SetStateReturns(result1 error) {
	fake.SetStateStub = nil
	fake.setStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetStateReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) DeleteState(namespace string, key string) error {
	fake.deleteStateMutex.Lock()
	ret, specificReturn := fake.deleteStateReturnsOnCall[len(fake.deleteStateArgsForCall)]
	fake.deleteStateArgsForCall = append(fake.deleteStateArgsForCall, struct {
		namespace string
		key       string
	}{namespace, key})
	fake.recordInvocation("DeleteState", []interface{}{namespace, key})
	fake.deleteStateMutex.Unlock()
	if fake.DeleteStateStub != nil {
		return fake.DeleteStateStub(namespace, key)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteStateReturns.result1
}

func (fake *TxSimulator) DeleteStateCallCount() int {
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	return len(fake.deleteStateArgsForCall)
}

func (fake *TxSimulator) DeleteStateArgsForCall(i int) (string, string) {
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	return fake.deleteStateArgsForCall[i].namespace, fake.deleteStateArgsForCall[i].key
}

func (fake *TxSimulator) DeleteStateReturns(result1 error) {
	fake.DeleteStateStub = nil
	fake.deleteStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeleteStateReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) SetStateMultipleKeys(namespace string, kvs map[string][]byte) error {
	fake.setStateMultipleKeysMutex.Lock()
	ret, specificReturn := fake.setStateMultipleKeysReturnsOnCall[len(fake.setStateMultipleKeysArgsForCall)]
	fake.setStateMultipleKeysArgsForCall = append(fake.setStateMultipleKeysArgsForCall, struct {
		namespace string
		kvs       map[string][]byte
	}{namespace, kvs})
	fake.recordInvocation("SetStateMultipleKeys", []interface{}{namespace, kvs})
	fake.setStateMultipleKeysMutex.Unlock()
	if fake.SetStateMultipleKeysStub != nil {
		return fake.SetStateMultipleKeysStub(namespace, kvs)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setStateMultipleKeysReturns.result1
}

func (fake *TxSimulator) SetStateMultipleKeysCallCount() int {
	fake.setStateMultipleKeysMutex.RLock()
	defer fake.setStateMultipleKeysMutex.RUnlock()
	return len(fake.setStateMultipleKeysArgsForCall)
}

func (fake *TxSimulator) SetStateMultipleKeysArgsForCall(i int) (string, map[string][]byte) {
	fake.setStateMultipleKeysMutex.RLock()
	defer fake.setStateMultipleKeysMutex.RUnlock()
	return fake.setStateMultipleKeysArgsForCall[i].namespace, fake.setStateMultipleKeysArgsForCall[i].kvs
}

func (fake *TxSimulator) SetStateMultipleKeysReturns(result1 error) {
	fake.SetStateMultipleKeysStub = nil
	fake.setStateMultipleKeysReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetStateMultipleKeysReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) SetStateMetadata(namespace string, key string, metadata map[string][]byte) error {
	fake.setStateMetadataMutex.Lock()
	ret, specificReturn := fake.setStateMetadataReturnsOnCall[len(fake.setStateMetadataArgsForCall)]
	fake.setStateMetadataArgsForCall = append(fake.setStateMetadataArgsForCall, struct {
		namespace string
		key       string
		metadata  map[string][]byte
	}{namespace, key, metadata})
	fake.recordInvocation("SetStateMetadata", []interface{}{namespace, key, metadata})
	fake.setStateMetadataMutex.Unlock()
	if fake.SetStateMetadataStub != nil {
		return fake.SetStateMetadataStub(namespace, key, metadata)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setStateMetadataReturns.result1
}

func (fake *TxSimulator) SetStateMetadataCallCount() int {
	fake.setStateMetadataMutex.RLock()
	defer fake.setStateMetadataMutex.RUnlock()
	return len(fake.setStateMetadataArgsForCall)
}

func (fake *TxSimulator) SetStateMetadataArgsForCall(i int) (string, string, map[string][]byte) {
	fake.setStateMetadataMutex.RLock()
	defer fake.setStateMetadataMutex.RUnlock()
	return fake.setStateMetadataArgsForCall[i].namespace, fake.setStateMetadataArgsForCall[i].key, fake.setStateMetadataArgsForCall[i].metadata
}

func (fake *TxSimulator) SetStateMetadataReturns(result1 error) {
	fake.SetStateMetadataStub = nil
	fake.setStateMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetStateMetadataReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) DeleteStateMetadata(namespace string, key string) error {
	fake.deleteStateMetadataMutex.Lock()
	ret, specificReturn := fake.deleteStateMetadataReturnsOnCall[len(fake.deleteStateMetadataArgsForCall)]
	fake.deleteStateMetadataArgsForCall = append(fake.deleteStateMetadataArgsForCall, struct {
		namespace string
		key       string
	}{namespace, key})
	fake.recordInvocation("DeleteStateMetadata", []interface{}{namespace, key})
	fake.deleteStateMetadataMutex.Unlock()
	if fake.DeleteStateMetadataStub != nil {
		return fake.DeleteStateMetadataStub(namespace, key)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteStateMetadataReturns.result1
}

func (fake *TxSimulator) DeleteStateMetadataCallCount() int {
	fake.deleteStateMetadataMutex.RLock()
	defer fake.deleteStateMetadataMutex.RUnlock()
	return len(fake.deleteStateMetadataArgsForCall)
}

func (fake *TxSimulator) DeleteStateMetadataArgsForCall(i int) (string, string) {
	fake.deleteStateMetadataMutex.RLock()
	defer fake.deleteStateMetadataMutex.RUnlock()
	return fake.deleteStateMetadataArgsForCall[i].namespace, fake.deleteStateMetadataArgsForCall[i].key
}

func (fake *TxSimulator) DeleteStateMetadataReturns(result1 error) {
	fake.DeleteStateMetadataStub = nil
	fake.deleteStateMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeleteStateMetadataReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) ExecuteUpdate(query string) error {
	fake.executeUpdateMutex.Lock()
	ret, specificReturn := fake.executeUpdateReturnsOnCall[len(fake.executeUpdateArgsForCall)]
	fake.executeUpdateArgsForCall = append(fake.executeUpdateArgsForCall, struct {
		query string
	}{query})
	fake.recordInvocation("ExecuteUpdate", []interface{}{query})
	fake.executeUpdateMutex.Unlock()
	if fake.ExecuteUpdateStub != nil {
		return fake.ExecuteUpdateStub(query)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.executeUpdateReturns.result1
}

func (fake *TxSimulator) ExecuteUpdateCallCount() int {
	fake.executeUpdateMutex.RLock()
	defer fake.executeUpdateMutex.RUnlock()
	return len(fake.executeUpdateArgsForCall)
}

func (fake *TxSimulator) ExecuteUpdateArgsForCall(i int) string {
	fake.executeUpdateMutex.RLock()
	defer fake.executeUpdateMutex.RUnlock()
	return fake.executeUpdateArgsForCall[i].query
}

func (fake *TxSimulator) ExecuteUpdateReturns(result1 error) {
	fake.ExecuteUpdateStub = nil
	fake.executeUpdateReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) ExecuteUpdateReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) SetPrivateData(namespace string, collection string, key string, value []byte) error {
	var valueCopy []byte
	if value != nil {
		valueCopy = make([]byte, len(value))
		copy(valueCopy, value)
	}
	fake.setPrivateDataMutex.Lock()
	ret, specificReturn := fake.setPrivateDataReturnsOnCall[len(fake.setPrivateDataArgsForCall)]
	fake.setPrivateDataArgsForCall = append(fake.setPrivateDataArgsForCall, struct {
		namespace  string
		collection string
		key        string
		value      []byte
	}{namespace, collection, key, valueCopy})
	fake.recordInvocation("SetPrivateData", []interface{}{namespace, collection, key, valueCopy})
	fake.setPrivateDataMutex.Unlock()
	if fake.SetPrivateDataStub != nil {
		return fake.SetPrivateDataStub(namespace, collection, key, value)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setPrivateDataReturns.result1
}

func (fake *TxSimulator) SetPrivateDataCallCount() int {
	fake.setPrivateDataMutex.RLock()
	defer fake.setPrivateDataMutex.RUnlock()
	return len(fake.setPrivateDataArgsForCall)
}

func (fake *TxSimulator) SetPrivateDataArgsForCall(i int) (string, string, string, []byte) {
	fake.setPrivateDataMutex.RLock()
	defer fake.setPrivateDataMutex.RUnlock()
	return fake.setPrivateDataArgsForCall[i].namespace, fake.setPrivateDataArgsForCall[i].collection, fake.setPrivateDataArgsForCall[i].key, fake.setPrivateDataArgsForCall[i].value
}

func (fake *TxSimulator) SetPrivateDataReturns(result1 error) {
	fake.SetPrivateDataStub = nil
	fake.setPrivateDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetPrivateDataReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) SetPrivateDataMultipleKeys(namespace string, collection string, kvs map[string][]byte) error {
	fake.setPrivateDataMultipleKeysMutex.Lock()
	ret, specificReturn := fake.setPrivateDataMultipleKeysReturnsOnCall[len(fake.setPrivateDataMultipleKeysArgsForCall)]
	fake.setPrivateDataMultipleKeysArgsForCall = append(fake.setPrivateDataMultipleKeysArgsForCall, struct {
		namespace  string
		collection string
		kvs        map[string][]byte
	}{namespace, collection, kvs})
	fake.recordInvocation("SetPrivateDataMultipleKeys", []interface{}{namespace, collection, kvs})
	fake.setPrivateDataMultipleKeysMutex.Unlock()
	if fake.SetPrivateDataMultipleKeysStub != nil {
		return fake.SetPrivateDataMultipleKeysStub(namespace, collection, kvs)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setPrivateDataMultipleKeysReturns.result1
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysCallCount() int {
	fake.setPrivateDataMultipleKeysMutex.RLock()
	defer fake.setPrivateDataMultipleKeysMutex.RUnlock()
	return len(fake.setPrivateDataMultipleKeysArgsForCall)
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysArgsForCall(i int) (string, string, map[string][]byte) {
	fake.setPrivateDataMultipleKeysMutex.RLock()
	defer fake.setPrivateDataMultipleKeysMutex.RUnlock()
	return fake.setPrivateDataMultipleKeysArgsForCall[i].namespace, fake.setPrivateDataMultipleKeysArgsForCall[i].collection, fake.setPrivateDataMultipleKeysArgsForCall[i].kvs
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysReturns(result1 error) {
	fake.SetPrivateDataMultipleKeysStub = nil
	fake.setPrivateDataMultipleKeysReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetPrivateDataMultipleKeysReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) DeletePrivateData(namespace string, collection string, key string) error {
	fake.deletePrivateDataMutex.Lock()
	ret, specificReturn := fake.deletePrivateDataReturnsOnCall[len(fake.deletePrivateDataArgsForCall)]
	fake.deletePrivateDataArgsForCall = append(fake.deletePrivateDataArgsForCall, struct {
		namespace  string
		collection string
		key        string
	}{namespace, collection, key})
	fake.recordInvocation("DeletePrivateData", []interface{}{namespace, collection, key})
	fake.deletePrivateDataMutex.Unlock()
	if fake.DeletePrivateDataStub != nil {
		return fake.DeletePrivateDataStub(namespace, collection, key)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deletePrivateDataReturns.result1
}

func (fake *TxSimulator) DeletePrivateDataCallCount() int {
	fake.deletePrivateDataMutex.RLock()
	defer fake.deletePrivateDataMutex.RUnlock()
	return len(fake.deletePrivateDataArgsForCall)
}

func (fake *TxSimulator) DeletePrivateDataArgsForCall(i int) (string, string, string) {
	fake.deletePrivateDataMutex.RLock()
	defer fake.deletePrivateDataMutex.RUnlock()
	return fake.deletePrivateDataArgsForCall[i].namespace, fake.deletePrivateDataArgsForCall[i].collection, fake.deletePrivateDataArgsForCall[i].key
}

func (fake *TxSimulator) DeletePrivateDataReturns(result1 error) {
	fake.DeletePrivateDataStub = nil
	fake.deletePrivateDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeletePrivateDataReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) SetPrivateDataMetadata(namespace string, collection string, key string, metadata map[string][]byte) error {
	fake.setPrivateDataMetadataMutex.Lock()
	ret, specificReturn := fake.setPrivateDataMetadataReturnsOnCall[len(fake.setPrivateDataMetadataArgsForCall)]
	fake.setPrivateDataMetadataArgsForCall = append(fake.setPrivateDataMetadataArgsForCall, struct {
		namespace  string
		collection string
		key        string
		metadata   map[string][]byte
	}{namespace, collection, key, metadata})
	fake.recordInvocation("SetPrivateDataMetadata", []interface{}{namespace, collection, key, metadata})
	fake.setPrivateDataMetadataMutex.Unlock()
	if fake.SetPrivateDataMetadataStub != nil {
		return fake.SetPrivateDataMetadataStub(namespace, collection, key, metadata)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setPrivateDataMetadataReturns.result1
}

func (fake *TxSimulator) SetPrivateDataMetadataCallCount() int {
	fake.setPrivateDataMetadataMutex.RLock()
	defer fake.setPrivateDataMetadataMutex.RUnlock()
	return len(fake.setPrivateDataMetadataArgsForCall)
}

func (fake *TxSimulator) SetPrivateDataMetadataArgsForCall(i int) (string, string, string, map[string][]byte) {
	fake.setPrivateDataMetadataMutex.RLock()
	defer fake.setPrivateDataMetadataMutex.RUnlock()
	return fake.setPrivateDataMetadataArgsForCall[i].namespace, fake.setPrivateDataMetadataArgsForCall[i].collection, fake.setPrivateDataMetadataArgsForCall[i].key, fake.setPrivateDataMetadataArgsForCall[i].metadata
}

func (fake *TxSimulator) SetPrivateDataMetadataReturns(result1 error) {
	fake.SetPrivateDataMetadataStub = nil
	fake.setPrivateDataMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) SetPrivateDataMetadataReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) DeletePrivateDataMetadata(namespace string, collection string, key string) error {
	fake.deletePrivateDataMetadataMutex.Lock()
	ret, specificReturn := fake.deletePrivateDataMetadataReturnsOnCall[len(fake.deletePrivateDataMetadataArgsForCall)]
	fake.deletePrivateDataMetadataArgsForCall = append(fake.deletePrivateDataMetadataArgsForCall, struct {
		namespace  string
		collection string
		key        string
	}{namespace, collection, key})
	fake.recordInvocation("DeletePrivateDataMetadata", []interface{}{namespace, collection, key})
	fake.deletePrivateDataMetadataMutex.Unlock()
	if fake.DeletePrivateDataMetadataStub != nil {
		return fake.DeletePrivateDataMetadataStub(namespace, collection, key)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deletePrivateDataMetadataReturns.result1
}

func (fake *TxSimulator) DeletePrivateDataMetadataCallCount() int {
	fake.deletePrivateDataMetadataMutex.RLock()
	defer fake.deletePrivateDataMetadataMutex.RUnlock()
	return len(fake.deletePrivateDataMetadataArgsForCall)
}

func (fake *TxSimulator) DeletePrivateDataMetadataArgsForCall(i int) (string, string, string) {
	fake.deletePrivateDataMetadataMutex.RLock()
	defer fake.deletePrivateDataMetadataMutex.RUnlock()
	return fake.deletePrivateDataMetadataArgsForCall[i].namespace, fake.deletePrivateDataMetadataArgsForCall[i].collection, fake.deletePrivateDataMetadataArgsForCall[i].key
}

func (fake *TxSimulator) DeletePrivateDataMetadataReturns(result1 error) {
	fake.DeletePrivateDataMetadataStub = nil
	fake.deletePrivateDataMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *TxSimulator) DeletePrivateDataMetadataReturnsOnCall(i int, result1 error) {
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

func (fake *TxSimulator) GetTxSimulationResults() (*ledger.TxSimulationResults, error) {
	fake.getTxSimulationResultsMutex.Lock()
	ret, specificReturn := fake.getTxSimulationResultsReturnsOnCall[len(fake.getTxSimulationResultsArgsForCall)]
	fake.getTxSimulationResultsArgsForCall = append(fake.getTxSimulationResultsArgsForCall, struct{}{})
	fake.recordInvocation("GetTxSimulationResults", []interface{}{})
	fake.getTxSimulationResultsMutex.Unlock()
	if fake.GetTxSimulationResultsStub != nil {
		return fake.GetTxSimulationResultsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTxSimulationResultsReturns.result1, fake.getTxSimulationResultsReturns.result2
}

func (fake *TxSimulator) GetTxSimulationResultsCallCount() int {
	fake.getTxSimulationResultsMutex.RLock()
	defer fake.getTxSimulationResultsMutex.RUnlock()
	return len(fake.getTxSimulationResultsArgsForCall)
}

func (fake *TxSimulator) GetTxSimulationResultsReturns(result1 *ledger.TxSimulationResults, result2 error) {
	fake.GetTxSimulationResultsStub = nil
	fake.getTxSimulationResultsReturns = struct {
		result1 *ledger.TxSimulationResults
		result2 error
	}{result1, result2}
}

func (fake *TxSimulator) GetTxSimulationResultsReturnsOnCall(i int, result1 *ledger.TxSimulationResults, result2 error) {
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

func (fake *TxSimulator) Invocations() map[string][][]interface{} {
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
	fake.executeQueryOnPrivateDataMutex.RLock()
	defer fake.executeQueryOnPrivateDataMutex.RUnlock()
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	fake.setStateMultipleKeysMutex.RLock()
	defer fake.setStateMultipleKeysMutex.RUnlock()
	fake.setStateMetadataMutex.RLock()
	defer fake.setStateMetadataMutex.RUnlock()
	fake.deleteStateMetadataMutex.RLock()
	defer fake.deleteStateMetadataMutex.RUnlock()
	fake.executeUpdateMutex.RLock()
	defer fake.executeUpdateMutex.RUnlock()
	fake.setPrivateDataMutex.RLock()
	defer fake.setPrivateDataMutex.RUnlock()
	fake.setPrivateDataMultipleKeysMutex.RLock()
	defer fake.setPrivateDataMultipleKeysMutex.RUnlock()
	fake.deletePrivateDataMutex.RLock()
	defer fake.deletePrivateDataMutex.RUnlock()
	fake.setPrivateDataMetadataMutex.RLock()
	defer fake.setPrivateDataMetadataMutex.RUnlock()
	fake.deletePrivateDataMetadataMutex.RLock()
	defer fake.deletePrivateDataMetadataMutex.RUnlock()
	fake.getTxSimulationResultsMutex.RLock()
	defer fake.getTxSimulationResultsMutex.RUnlock()
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
