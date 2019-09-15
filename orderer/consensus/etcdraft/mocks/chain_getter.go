

package mocks

import (
	multichannel "github.com/mcc-github/blockchain/orderer/common/multichannel"
	mock "github.com/stretchr/testify/mock"
)


type ChainGetter struct {
	mock.Mock
}


func (_m *ChainGetter) GetChain(chainID string) *multichannel.ChainSupport {
	ret := _m.Called(chainID)

	var r0 *multichannel.ChainSupport
	if rf, ok := ret.Get(0).(func(string) *multichannel.ChainSupport); ok {
		r0 = rf(chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*multichannel.ChainSupport)
		}
	}

	return r0
}
