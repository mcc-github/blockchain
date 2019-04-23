

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import orderer "github.com/mcc-github/blockchain/protos/orderer"


type Dispatcher struct {
	mock.Mock
}


func (_m *Dispatcher) DispatchConsensus(ctx context.Context, request *orderer.ConsensusRequest) error {
	ret := _m.Called(ctx, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *orderer.ConsensusRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *Dispatcher) DispatchSubmit(ctx context.Context, request *orderer.SubmitRequest) error {
	ret := _m.Called(ctx, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *orderer.SubmitRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
