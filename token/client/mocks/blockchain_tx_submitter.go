
package mocks

import mock "github.com/stretchr/testify/mock"


type FabricTxSubmitter struct {
	mock.Mock
}


func (_m *FabricTxSubmitter) Submit(tx []byte) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
