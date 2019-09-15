

package mocks

import (
	policies "github.com/mcc-github/blockchain/common/policies"
	mock "github.com/stretchr/testify/mock"
)


type ChannelPolicyReferenceProvider struct {
	mock.Mock
}


func (_m *ChannelPolicyReferenceProvider) NewPolicy(channelConfigPolicyReference string) (policies.Policy, error) {
	ret := _m.Called(channelConfigPolicyReference)

	var r0 policies.Policy
	if rf, ok := ret.Get(0).(func(string) policies.Policy); ok {
		r0 = rf(channelConfigPolicyReference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(policies.Policy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(channelConfigPolicyReference)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
