

package mocks

import mock "github.com/stretchr/testify/mock"
import policies "github.com/mcc-github/blockchain/common/policies"


type ChannelPolicyManagerGetter struct {
	mock.Mock
}


func (_m *ChannelPolicyManagerGetter) Manager(channelID string) (policies.Manager, bool) {
	ret := _m.Called(channelID)

	var r0 policies.Manager
	if rf, ok := ret.Get(0).(func(string) policies.Manager); ok {
		r0 = rf(channelID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(policies.Manager)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(channelID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}
