

package mocks

import (
	ledger "github.com/mcc-github/blockchain/core/ledger"
	mock "github.com/stretchr/testify/mock"

	transientstore "github.com/mcc-github/blockchain/core/transientstore"
)


type DataStore struct {
	mock.Mock
}


func (_m *DataStore) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	ret := _m.Called()

	var r0 ledger.ConfigHistoryRetriever
	if rf, ok := ret.Get(0).(func() ledger.ConfigHistoryRetriever); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.ConfigHistoryRetriever)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *DataStore) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	ret := _m.Called(blockNum, filter)

	var r0 []*ledger.TxPvtData
	if rf, ok := ret.Get(0).(func(uint64, ledger.PvtNsCollFilter) []*ledger.TxPvtData); ok {
		r0 = rf(blockNum, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ledger.TxPvtData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, ledger.PvtNsCollFilter) error); ok {
		r1 = rf(blockNum, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *DataStore) GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error) {
	ret := _m.Called(txid, filter)

	var r0 transientstore.RWSetScanner
	if rf, ok := ret.Get(0).(func(string, ledger.PvtNsCollFilter) transientstore.RWSetScanner); ok {
		r0 = rf(txid, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transientstore.RWSetScanner)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ledger.PvtNsCollFilter) error); ok {
		r1 = rf(txid, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *DataStore) LedgerHeight() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
