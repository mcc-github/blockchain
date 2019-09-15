

package mocks

import (
	ledger "github.com/mcc-github/blockchain/core/ledger"
	mock "github.com/stretchr/testify/mock"
)


type LifecycleResources struct {
	mock.Mock
}


func (_m *LifecycleResources) ValidationInfo(channelID string, chaincodeName string, qe ledger.SimpleQueryExecutor) (string, []byte, error, error) {
	ret := _m.Called(channelID, chaincodeName, qe)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, ledger.SimpleQueryExecutor) string); ok {
		r0 = rf(channelID, chaincodeName, qe)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 []byte
	if rf, ok := ret.Get(1).(func(string, string, ledger.SimpleQueryExecutor) []byte); ok {
		r1 = rf(channelID, chaincodeName, qe)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, ledger.SimpleQueryExecutor) error); ok {
		r2 = rf(channelID, chaincodeName, qe)
	} else {
		r2 = ret.Error(2)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(string, string, ledger.SimpleQueryExecutor) error); ok {
		r3 = rf(channelID, chaincodeName, qe)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}
