

package mocks

import coreledger "github.com/mcc-github/blockchain/core/ledger"
import ledger "github.com/mcc-github/blockchain/common/ledger"
import mock "github.com/stretchr/testify/mock"


type QueryExecutor struct {
	mock.Mock
}


func (_m *QueryExecutor) Done() {
	_m.Called()
}


func (_m *QueryExecutor) ExecuteQuery(namespace string, query string) (ledger.ResultsIterator, error) {
	ret := _m.Called(namespace, query)

	var r0 ledger.ResultsIterator
	if rf, ok := ret.Get(0).(func(string, string) ledger.ResultsIterator); ok {
		r0 = rf(namespace, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.ResultsIterator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(namespace, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) ExecuteQueryOnPrivateData(namespace string, collection string, query string) (ledger.ResultsIterator, error) {
	ret := _m.Called(namespace, collection, query)

	var r0 ledger.ResultsIterator
	if rf, ok := ret.Get(0).(func(string, string, string) ledger.ResultsIterator); ok {
		r0 = rf(namespace, collection, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.ResultsIterator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(namespace, collection, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) ExecuteQueryWithMetadata(namespace string, query string, metadata map[string]interface{}) (coreledger.QueryResultsIterator, error) {
	ret := _m.Called(namespace, query, metadata)

	var r0 coreledger.QueryResultsIterator
	if rf, ok := ret.Get(0).(func(string, string, map[string]interface{}) coreledger.QueryResultsIterator); ok {
		r0 = rf(namespace, query, metadata)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreledger.QueryResultsIterator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, map[string]interface{}) error); ok {
		r1 = rf(namespace, query, metadata)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetPrivateData(namespace string, collection string, key string) ([]byte, error) {
	ret := _m.Called(namespace, collection, key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string, string) []byte); ok {
		r0 = rf(namespace, collection, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(namespace, collection, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetPrivateDataHash(namespace string, collection string, key string) ([]byte, error) {
	ret := _m.Called(namespace, collection, key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string, string) []byte); ok {
		r0 = rf(namespace, collection, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(namespace, collection, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetPrivateDataMetadata(namespace string, collection string, key string) (map[string][]byte, error) {
	ret := _m.Called(namespace, collection, key)

	var r0 map[string][]byte
	if rf, ok := ret.Get(0).(func(string, string, string) map[string][]byte); ok {
		r0 = rf(namespace, collection, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(namespace, collection, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetPrivateDataMetadataByHash(namespace string, collection string, keyhash []byte) (map[string][]byte, error) {
	ret := _m.Called(namespace, collection, keyhash)

	var r0 map[string][]byte
	if rf, ok := ret.Get(0).(func(string, string, []byte) map[string][]byte); ok {
		r0 = rf(namespace, collection, keyhash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, []byte) error); ok {
		r1 = rf(namespace, collection, keyhash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetPrivateDataMultipleKeys(namespace string, collection string, keys []string) ([][]byte, error) {
	ret := _m.Called(namespace, collection, keys)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func(string, string, []string) [][]byte); ok {
		r0 = rf(namespace, collection, keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, []string) error); ok {
		r1 = rf(namespace, collection, keys)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetPrivateDataRangeScanIterator(namespace string, collection string, startKey string, endKey string) (ledger.ResultsIterator, error) {
	ret := _m.Called(namespace, collection, startKey, endKey)

	var r0 ledger.ResultsIterator
	if rf, ok := ret.Get(0).(func(string, string, string, string) ledger.ResultsIterator); ok {
		r0 = rf(namespace, collection, startKey, endKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.ResultsIterator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(namespace, collection, startKey, endKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetState(namespace string, key string) ([]byte, error) {
	ret := _m.Called(namespace, key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string) []byte); ok {
		r0 = rf(namespace, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(namespace, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetStateMetadata(namespace string, key string) (map[string][]byte, error) {
	ret := _m.Called(namespace, key)

	var r0 map[string][]byte
	if rf, ok := ret.Get(0).(func(string, string) map[string][]byte); ok {
		r0 = rf(namespace, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(namespace, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error) {
	ret := _m.Called(namespace, keys)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func(string, []string) [][]byte); ok {
		r0 = rf(namespace, keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []string) error); ok {
		r1 = rf(namespace, keys)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (ledger.ResultsIterator, error) {
	ret := _m.Called(namespace, startKey, endKey)

	var r0 ledger.ResultsIterator
	if rf, ok := ret.Get(0).(func(string, string, string) ledger.ResultsIterator); ok {
		r0 = rf(namespace, startKey, endKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.ResultsIterator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(namespace, startKey, endKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *QueryExecutor) GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (coreledger.QueryResultsIterator, error) {
	ret := _m.Called(namespace, startKey, endKey, metadata)

	var r0 coreledger.QueryResultsIterator
	if rf, ok := ret.Get(0).(func(string, string, string, map[string]interface{}) coreledger.QueryResultsIterator); ok {
		r0 = rf(namespace, startKey, endKey, metadata)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreledger.QueryResultsIterator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, map[string]interface{}) error); ok {
		r1 = rf(namespace, startKey, endKey, metadata)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
