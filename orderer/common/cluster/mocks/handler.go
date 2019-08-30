

package mocks

import mock "github.com/stretchr/testify/mock"
import orderer "github.com/mcc-github/blockchain-protos-go/orderer"


type Handler struct {
	mock.Mock
}


func (_m *Handler) OnConsensus(channel string, sender uint64, req *orderer.ConsensusRequest) error {
	ret := _m.Called(channel, sender, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uint64, *orderer.ConsensusRequest) error); ok {
		r0 = rf(channel, sender, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *Handler) OnSubmit(channel string, sender uint64, req *orderer.SubmitRequest) error {
	ret := _m.Called(channel, sender, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uint64, *orderer.SubmitRequest) error); ok {
		r0 = rf(channel, sender, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
