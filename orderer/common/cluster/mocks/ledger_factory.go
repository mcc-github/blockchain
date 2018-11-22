
package mocks

import cluster "github.com/mcc-github/blockchain/orderer/common/cluster"
import mock "github.com/stretchr/testify/mock"


type LedgerFactory struct {
	mock.Mock
}


func (_m *LedgerFactory) ChainIDs() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}


func (_m *LedgerFactory) Close() {
	_m.Called()
}


func (_m *LedgerFactory) GetOrCreate(chainID string) (cluster.LedgerWriter, error) {
	ret := _m.Called(chainID)

	var r0 cluster.LedgerWriter
	if rf, ok := ret.Get(0).(func(string) cluster.LedgerWriter); ok {
		r0 = rf(chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cluster.LedgerWriter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(chainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}