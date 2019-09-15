

package mocks

import (
	common "github.com/mcc-github/blockchain-protos-go/common"
	mock "github.com/stretchr/testify/mock"

	policies "github.com/mcc-github/blockchain/common/policies"
)


type SignaturePolicyProvider struct {
	mock.Mock
}


func (_m *SignaturePolicyProvider) NewPolicy(signaturePolicy *common.SignaturePolicyEnvelope) (policies.Policy, error) {
	ret := _m.Called(signaturePolicy)

	var r0 policies.Policy
	if rf, ok := ret.Get(0).(func(*common.SignaturePolicyEnvelope) policies.Policy); ok {
		r0 = rf(signaturePolicy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(policies.Policy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*common.SignaturePolicyEnvelope) error); ok {
		r1 = rf(signaturePolicy)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
