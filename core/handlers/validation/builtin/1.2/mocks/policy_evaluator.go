
package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"


type PolicyEvaluator struct {
	mock.Mock
}


func (_m *PolicyEvaluator) Evaluate(policyBytes []byte, signatureSet []*common.SignedData) error {
	ret := _m.Called(policyBytes, signatureSet)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []*common.SignedData) error); ok {
		r0 = rf(policyBytes, signatureSet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
