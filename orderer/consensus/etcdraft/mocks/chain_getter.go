
package mocks

import mock "github.com/stretchr/testify/mock"
import multichannel "github.com/mcc-github/blockchain/orderer/common/multichannel"


type ChainGetter struct {
	mock.Mock
}


func (_m *ChainGetter) GetChain(chainID string) (*multichannel.ChainSupport, bool) {
	ret := _m.Called(chainID)

	var r0 *multichannel.ChainSupport
	if rf, ok := ret.Get(0).(func(string) *multichannel.ChainSupport); ok {
		r0 = rf(chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*multichannel.ChainSupport)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(chainID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}
