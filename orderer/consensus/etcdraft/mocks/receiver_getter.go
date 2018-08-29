
package mocks

import etcdraft "github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
import mock "github.com/stretchr/testify/mock"


type ReceiverGetter struct {
	mock.Mock
}


func (_m *ReceiverGetter) ReceiverByChain(channelID string) etcdraft.MessageReceiver {
	ret := _m.Called(channelID)

	var r0 etcdraft.MessageReceiver
	if rf, ok := ret.Get(0).(func(string) etcdraft.MessageReceiver); ok {
		r0 = rf(channelID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(etcdraft.MessageReceiver)
		}
	}

	return r0
}
