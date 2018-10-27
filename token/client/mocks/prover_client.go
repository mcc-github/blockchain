
package mocks

import context "context"
import grpc "google.golang.org/grpc"
import mock "github.com/stretchr/testify/mock"
import token "github.com/mcc-github/blockchain/protos/token"


type ProverClient struct {
	mock.Mock
}


func (_m *ProverClient) ProcessCommand(ctx context.Context, in *token.SignedCommand, opts ...grpc.CallOption) (*token.SignedCommandResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *token.SignedCommandResponse
	if rf, ok := ret.Get(0).(func(context.Context, *token.SignedCommand, ...grpc.CallOption) *token.SignedCommandResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token.SignedCommandResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *token.SignedCommand, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
