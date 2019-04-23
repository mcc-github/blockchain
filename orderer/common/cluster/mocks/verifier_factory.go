

package mocks

import cluster "github.com/mcc-github/blockchain/orderer/common/cluster"
import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"


type VerifierFactory struct {
	mock.Mock
}


func (_m *VerifierFactory) VerifierFromConfig(configuration *common.ConfigEnvelope, channel string) (cluster.BlockVerifier, error) {
	ret := _m.Called(configuration, channel)

	var r0 cluster.BlockVerifier
	if rf, ok := ret.Get(0).(func(*common.ConfigEnvelope, string) cluster.BlockVerifier); ok {
		r0 = rf(configuration, channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cluster.BlockVerifier)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*common.ConfigEnvelope, string) error); ok {
		r1 = rf(configuration, channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
