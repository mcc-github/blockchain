

package mocks

import ledger "github.com/mcc-github/blockchain/core/ledger"
import mock "github.com/stretchr/testify/mock"


type ConfigHistoryRetriever struct {
	mock.Mock
}


func (_m *ConfigHistoryRetriever) CollectionConfigAt(blockNum uint64, chaincodeName string) (*ledger.CollectionConfigInfo, error) {
	ret := _m.Called(blockNum, chaincodeName)

	var r0 *ledger.CollectionConfigInfo
	if rf, ok := ret.Get(0).(func(uint64, string) *ledger.CollectionConfigInfo); ok {
		r0 = rf(blockNum, chaincodeName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ledger.CollectionConfigInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, string) error); ok {
		r1 = rf(blockNum, chaincodeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *ConfigHistoryRetriever) MostRecentCollectionConfigBelow(blockNum uint64, chaincodeName string) (*ledger.CollectionConfigInfo, error) {
	ret := _m.Called(blockNum, chaincodeName)

	var r0 *ledger.CollectionConfigInfo
	if rf, ok := ret.Get(0).(func(uint64, string) *ledger.CollectionConfigInfo); ok {
		r0 = rf(blockNum, chaincodeName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ledger.CollectionConfigInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, string) error); ok {
		r1 = rf(blockNum, chaincodeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
