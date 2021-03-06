

package mocks

import (
	common "github.com/mcc-github/blockchain-protos-go/common"
	errors "github.com/mcc-github/blockchain/common/errors"

	mock "github.com/stretchr/testify/mock"

	peer "github.com/mcc-github/blockchain-protos-go/peer"
)


type StateBasedValidator struct {
	mock.Mock
}


func (_m *StateBasedValidator) PostValidate(cc string, blockNum uint64, txNum uint64, err error) {
	_m.Called(cc, blockNum, txNum, err)
}


func (_m *StateBasedValidator) PreValidate(txNum uint64, block *common.Block) {
	_m.Called(txNum, block)
}


func (_m *StateBasedValidator) Validate(cc string, blockNum uint64, txNum uint64, rwset []byte, prp []byte, ep []byte, endorsements []*peer.Endorsement) errors.TxValidationError {
	ret := _m.Called(cc, blockNum, txNum, rwset, prp, ep, endorsements)

	var r0 errors.TxValidationError
	if rf, ok := ret.Get(0).(func(string, uint64, uint64, []byte, []byte, []byte, []*peer.Endorsement) errors.TxValidationError); ok {
		r0 = rf(cc, blockNum, txNum, rwset, prp, ep, endorsements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(errors.TxValidationError)
		}
	}

	return r0
}
