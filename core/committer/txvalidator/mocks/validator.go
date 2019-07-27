

package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"


type Validator struct {
	mock.Mock
}


func (_m *Validator) Validate(block *common.Block) error {
	ret := _m.Called(block)

	var r0 error
	if rf, ok := ret.Get(0).(func(*common.Block) error); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
