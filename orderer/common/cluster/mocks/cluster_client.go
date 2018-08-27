
package mocks

import context "context"
import grpc "google.golang.org/grpc"
import mock "github.com/stretchr/testify/mock"
import orderer "github.com/mcc-github/blockchain/protos/orderer"


type ClusterClient struct {
	mock.Mock
}


func (_m *ClusterClient) Step(ctx context.Context, in *orderer.StepRequest, opts ...grpc.CallOption) (*orderer.StepResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *orderer.StepResponse
	if rf, ok := ret.Get(0).(func(context.Context, *orderer.StepRequest, ...grpc.CallOption) *orderer.StepResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*orderer.StepResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *orderer.StepRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *ClusterClient) Submit(ctx context.Context, opts ...grpc.CallOption) (orderer.Cluster_SubmitClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 orderer.Cluster_SubmitClient
	if rf, ok := ret.Get(0).(func(context.Context, ...grpc.CallOption) orderer.Cluster_SubmitClient); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(orderer.Cluster_SubmitClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
