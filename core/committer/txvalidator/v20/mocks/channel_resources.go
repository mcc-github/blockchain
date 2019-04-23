

package mocks

import channelconfig "github.com/mcc-github/blockchain/common/channelconfig"
import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"
import msp "github.com/mcc-github/blockchain/msp"


type ChannelResources struct {
	mock.Mock
}


func (_m *ChannelResources) Apply(configtx *common.ConfigEnvelope) error {
	ret := _m.Called(configtx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*common.ConfigEnvelope) error); ok {
		r0 = rf(configtx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *ChannelResources) Capabilities() channelconfig.ApplicationCapabilities {
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


func (_m *ChannelResources) GetMSPIDs(cid string) []string {
	ret := _m.Called(cid)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(cid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}


func (_m *ChannelResources) MSPManager() msp.MSPManager {
	ret := _m.Called()

	var r0 msp.MSPManager
	if rf, ok := ret.Get(0).(func() msp.MSPManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(msp.MSPManager)
		}
	}

	return r0
}
