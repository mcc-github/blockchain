
package mocks

import context "context"
import metadata "google.golang.org/grpc/metadata"
import mock "github.com/stretchr/testify/mock"
import orderer "github.com/mcc-github/blockchain/protos/orderer"


type SubmitStream struct {
	mock.Mock
}


func (_m *SubmitStream) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}


func (_m *SubmitStream) Recv() (*orderer.SubmitRequest, error) {
	ret := _m.Called()

	var r0 *orderer.SubmitRequest
	if rf, ok := ret.Get(0).(func() *orderer.SubmitRequest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*orderer.SubmitRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *SubmitStream) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *SubmitStream) Send(response *orderer.SubmitResponse) error {
	ret := _m.Called(response)

	var r0 error
	if rf, ok := ret.Get(0).(func(*orderer.SubmitResponse) error); ok {
		r0 = rf(response)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *SubmitStream) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *SubmitStream) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *SubmitStream) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *SubmitStream) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}
