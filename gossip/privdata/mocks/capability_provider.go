

package mocks

import (
	channelconfig "github.com/mcc-github/blockchain/common/channelconfig"
	mock "github.com/stretchr/testify/mock"
)


type CapabilityProvider struct {
	mock.Mock
}


func (_m *CapabilityProvider) Capabilities() channelconfig.ApplicationCapabilities {
	ret := _m.Called()

	var r0 channelconfig.ApplicationCapabilities
	if rf, ok := ret.Get(0).(func() channelconfig.ApplicationCapabilities); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(channelconfig.ApplicationCapabilities)
		}
	}

	return r0
}
