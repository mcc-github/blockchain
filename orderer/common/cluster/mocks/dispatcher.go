
package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import orderer "github.com/mcc-github/blockchain/protos/orderer"


type Dispatcher struct {
	mock.Mock
}


func (_m *Dispatcher) DispatchStep(ctx context.Context, request *orderer.StepRequest) (*orderer.StepResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *orderer.StepResponse
	if rf, ok := ret.Get(0).(func(context.Context, *orderer.StepRequest) *orderer.StepResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*orderer.StepResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *orderer.StepRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Dispatcher) DispatchSubmit(ctx context.Context, request *orderer.SubmitRequest) (*orderer.SubmitResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *orderer.SubmitResponse
	if rf, ok := ret.Get(0).(func(context.Context, *orderer.SubmitRequest) *orderer.SubmitResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*orderer.SubmitResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *orderer.SubmitRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
