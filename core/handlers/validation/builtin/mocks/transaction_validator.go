
package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import errors "github.com/mcc-github/blockchain/common/errors"
import mock "github.com/stretchr/testify/mock"


type TransactionValidator struct {
	mock.Mock
}


func (_m *TransactionValidator) Validate(block *common.Block, namespace string, txPosition int, actionPosition int, policy []byte) errors.TxValidationError {
	ret := _m.Called(block, namespace, txPosition, actionPosition, policy)

	var r0 errors.TxValidationError
	if rf, ok := ret.Get(0).(func(*common.Block, string, int, int, []byte) errors.TxValidationError); ok {
		r0 = rf(block, namespace, txPosition, actionPosition, policy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(errors.TxValidationError)
		}
	}

	return r0
}
