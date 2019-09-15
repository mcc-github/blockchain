

package mocks

import (
	cluster "github.com/mcc-github/blockchain/orderer/common/cluster"
	mock "github.com/stretchr/testify/mock"
)


type VerifierRetriever struct {
	mock.Mock
}


func (_m *VerifierRetriever) RetrieveVerifier(channel string) cluster.BlockVerifier {
	ret := _m.Called(channel)

	var r0 cluster.BlockVerifier
	if rf, ok := ret.Get(0).(func(string) cluster.BlockVerifier); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cluster.BlockVerifier)
		}
	}

	return r0
}
