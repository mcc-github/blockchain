

package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"
import peer "github.com/mcc-github/blockchain/protos/peer"


type Dispatcher struct {
	mock.Mock
}


func (_m *Dispatcher) Dispatch(seq int, payload *common.Payload, envBytes []byte, block *common.Block) (error, peer.TxValidationCode) {
	ret := _m.Called(seq, payload, envBytes, block)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, *common.Payload, []byte, *common.Block) error); ok {
		r0 = rf(seq, payload, envBytes, block)
	} else {
		r0 = ret.Error(0)
	}

	var r1 peer.TxValidationCode
	if rf, ok := ret.Get(1).(func(int, *common.Payload, []byte, *common.Block) peer.TxValidationCode); ok {
		r1 = rf(seq, payload, envBytes, block)
	} else {
		r1 = ret.Get(1).(peer.TxValidationCode)
	}

	return r0, r1
}
