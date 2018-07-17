
package mocks

import mock "github.com/stretchr/testify/mock"


type Query struct {
	mock.Mock
}


func (_m *Query) Done() {
	_m.Called()
}


func (_m *Query) GetState(namespace string, key string) ([]byte, error) {
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
