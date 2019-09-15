

package mocks

import (
	cluster "github.com/mcc-github/blockchain/orderer/common/cluster"
	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)


type SecureDialer struct {
	mock.Mock
}


func (_m *SecureDialer) Dial(address string, verifyFunc cluster.RemoteVerifier) (*grpc.ClientConn, error) {
	ret := _m.Called(address, verifyFunc)

	var r0 *grpc.ClientConn
	if rf, ok := ret.Get(0).(func(string, cluster.RemoteVerifier) *grpc.ClientConn); ok {
		r0 = rf(address, verifyFunc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.ClientConn)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, cluster.RemoteVerifier) error); ok {
		r1 = rf(address, verifyFunc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
