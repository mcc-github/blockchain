

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	transientstore "github.com/mcc-github/blockchain/core/transientstore"
)


type RWSetScanner struct {
	mock.Mock
}


func (_m *RWSetScanner) Close() {
	_m.Called()
}


func (_m *RWSetScanner) Next() (*transientstore.EndorserPvtSimulationResults, error) {
	ret := _m.Called()

	var r0 *transientstore.EndorserPvtSimulationResults
	if rf, ok := ret.Get(0).(func() *transientstore.EndorserPvtSimulationResults); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transientstore.EndorserPvtSimulationResults)
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
