

package mocks

import cluster "github.com/mcc-github/blockchain/orderer/common/cluster"
import mock "github.com/stretchr/testify/mock"


type Communicator struct {
	mock.Mock
}


func (_m *Communicator) Configure(channel string, members []cluster.RemoteNode) {
	_m.Called(channel, members)
}


func (_m *Communicator) Remote(channel string, id uint64) (*cluster.RemoteContext, error) {
	ret := _m.Called(channel, id)

	var r0 *cluster.RemoteContext
	if rf, ok := ret.Get(0).(func(string, uint64) *cluster.RemoteContext); ok {
		r0 = rf(channel, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cluster.RemoteContext)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64) error); ok {
		r1 = rf(channel, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Communicator) Shutdown() {
	_m.Called()
}
