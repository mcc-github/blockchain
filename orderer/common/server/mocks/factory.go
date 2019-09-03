

package mocks

import blockledger "github.com/mcc-github/blockchain/common/ledger/blockledger"
import mock "github.com/stretchr/testify/mock"


type Factory struct {
	mock.Mock
}


func (_m *Factory) ChannelIDs() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}


func (_m *Factory) Close() {
	_m.Called()
}


func (_m *Factory) GetOrCreate(chainID string) (blockledger.ReadWriter, error) {
	ret := _m.Called(chainID)

	var r0 blockledger.ReadWriter
	if rf, ok := ret.Get(0).(func(string) blockledger.ReadWriter); ok {
		r0 = rf(chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(blockledger.ReadWriter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(chainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
