
package mocks

import errors "github.com/mcc-github/blockchain/common/errors"
import mock "github.com/stretchr/testify/mock"


type TransactionValidator struct {
	mock.Mock
}


func (_m *TransactionValidator) Validate(txData []byte, policy []byte) errors.TxValidationError {
	ret := _m.Called(txData, policy)

	var r0 errors.TxValidationError
	if rf, ok := ret.Get(0).(func([]byte, []byte) errors.TxValidationError); ok {
		r0 = rf(txData, policy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(errors.TxValidationError)
		}
	}

	return r0
}
