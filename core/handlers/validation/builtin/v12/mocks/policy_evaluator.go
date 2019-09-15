

package mocks

import (
	protoutil "github.com/mcc-github/blockchain/protoutil"
	mock "github.com/stretchr/testify/mock"
)


type PolicyEvaluator struct {
	mock.Mock
}


func (_m *PolicyEvaluator) Evaluate(policyBytes []byte, signatureSet []*protoutil.SignedData) error {
	ret := _m.Called(policyBytes, signatureSet)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []*protoutil.SignedData) error); ok {
		r0 = rf(policyBytes, signatureSet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
