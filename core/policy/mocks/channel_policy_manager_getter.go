

package mocks

import (
	policies "github.com/mcc-github/blockchain/common/policies"
	mock "github.com/stretchr/testify/mock"
)


type ChannelPolicyManagerGetter struct {
	mock.Mock
}


func (_m *ChannelPolicyManagerGetter) Manager(channelID string) policies.Manager {
	ret := _m.Called(channelID)

	var r0 policies.Manager
	if rf, ok := ret.Get(0).(func(string) policies.Manager); ok {
		r0 = rf(channelID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(policies.Manager)
		}
	}

	return r0
}
