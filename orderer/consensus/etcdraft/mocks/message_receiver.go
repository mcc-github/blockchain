

package mocks

import (
	orderer "github.com/mcc-github/blockchain-protos-go/orderer"
	mock "github.com/stretchr/testify/mock"
)


type MessageReceiver struct {
	mock.Mock
}


func (_m *MessageReceiver) Consensus(req *orderer.ConsensusRequest, sender uint64) error {
	ret := _m.Called(req, sender)

	var r0 error
	if rf, ok := ret.Get(0).(func(*orderer.ConsensusRequest, uint64) error); ok {
		r0 = rf(req, sender)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *MessageReceiver) Submit(req *orderer.SubmitRequest, sender uint64) error {
	ret := _m.Called(req, sender)

	var r0 error
	if rf, ok := ret.Get(0).(func(*orderer.SubmitRequest, uint64) error); ok {
		r0 = rf(req, sender)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
