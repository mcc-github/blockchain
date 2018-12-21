

package mocks

import cluster "github.com/mcc-github/blockchain/orderer/common/cluster"
import mock "github.com/stretchr/testify/mock"


type ChannelLister struct {
	mock.Mock
}


func (_m *ChannelLister) Channels() []cluster.ChannelGenesisBlock {
	ret := _m.Called()

	var r0 []cluster.ChannelGenesisBlock
	if rf, ok := ret.Get(0).(func() []cluster.ChannelGenesisBlock); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]cluster.ChannelGenesisBlock)
		}
	}

	return r0
}


func (_m *ChannelLister) Close() {
	_m.Called()
}
