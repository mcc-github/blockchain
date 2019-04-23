

package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"


type BlockRetriever struct {
	mock.Mock
}


func (_m *BlockRetriever) Block(number uint64) *common.Block {
	ret := _m.Called(number)

	var r0 *common.Block
	if rf, ok := ret.Get(0).(func(uint64) *common.Block); ok {
		r0 = rf(number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.Block)
		}
	}

	return r0
}
