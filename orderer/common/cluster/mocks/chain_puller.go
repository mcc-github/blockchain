
package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"


type ChainPuller struct {
	mock.Mock
}


func (_m *ChainPuller) Close() {
	_m.Called()
}


func (_m *ChainPuller) HeightsByEndpoints() map[string]uint64 {
	ret := _m.Called()

	var r0 map[string]uint64
	if rf, ok := ret.Get(0).(func() map[string]uint64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]uint64)
		}
	}

	return r0
}


func (_m *ChainPuller) PullBlock(seq uint64) *common.Block {
	ret := _m.Called(seq)

	var r0 *common.Block
	if rf, ok := ret.Get(0).(func(uint64) *common.Block); ok {
		r0 = rf(seq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.Block)
		}
	}

	return r0
}
