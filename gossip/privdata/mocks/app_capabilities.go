

package mocks

import mock "github.com/stretchr/testify/mock"


type AppCapabilities struct {
	mock.Mock
}


func (_m *AppCapabilities) StorePvtDataOfInvalidTx() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
