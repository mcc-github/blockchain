
package mocks

import cc "github.com/mcc-github/blockchain/core/cclifecycle"
import mock "github.com/stretchr/testify/mock"


type QueryCreator struct {
	mock.Mock
}


func (_m *QueryCreator) NewQuery() (cc.Query, error) {
	ret := _m.Called()

	var r0 cc.Query
	if rf, ok := ret.Get(0).(func() cc.Query); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cc.Query)
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
