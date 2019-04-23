

package mocks

import mock "github.com/stretchr/testify/mock"

import validation "github.com/mcc-github/blockchain/core/handlers/validation/api/state"


type CollectionResources struct {
	mock.Mock
}


func (_m *CollectionResources) CollectionValidationInfo(chaincodeName string, collectionName string, state validation.State) ([]byte, error, error) {
	ret := _m.Called(chaincodeName, collectionName, state)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string, validation.State) []byte); ok {
		r0 = rf(chaincodeName, collectionName, state)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, validation.State) error); ok {
		r1 = rf(chaincodeName, collectionName, state)
	} else {
		r1 = ret.Error(1)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, validation.State) error); ok {
		r2 = rf(chaincodeName, collectionName, state)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
