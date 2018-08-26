
package mocks

import mock "github.com/stretchr/testify/mock"
import orderer "github.com/mcc-github/blockchain/protos/orderer"


type Handler struct {
	mock.Mock
}


func (_m *Handler) OnStep(channel string, sender uint64, req *orderer.StepRequest) (*orderer.StepResponse, error) {
	ret := _m.Called(channel, sender, req)

	var r0 *orderer.StepResponse
	if rf, ok := ret.Get(0).(func(string, uint64, *orderer.StepRequest) *orderer.StepResponse); ok {
		r0 = rf(channel, sender, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*orderer.StepResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64, *orderer.StepRequest) error); ok {
		r1 = rf(channel, sender, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Handler) OnSubmit(channel string, sender uint64, req *orderer.SubmitRequest) (*orderer.SubmitResponse, error) {
	ret := _m.Called(channel, sender, req)

	var r0 *orderer.SubmitResponse
	if rf, ok := ret.Get(0).(func(string, uint64, *orderer.SubmitRequest) *orderer.SubmitResponse); ok {
		r0 = rf(channel, sender, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*orderer.SubmitResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64, *orderer.SubmitRequest) error); ok {
		r1 = rf(channel, sender, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
