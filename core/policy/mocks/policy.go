

package mocks

import mock "github.com/stretchr/testify/mock"

import protoutil "github.com/mcc-github/blockchain/protoutil"


type Policy struct {
	mock.Mock
}


func (_m *Policy) Evaluate(signatureSet []*protoutil.SignedData) error {
	ret := _m.Called(signatureSet)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*protoutil.SignedData) error); ok {
		r0 = rf(signatureSet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
