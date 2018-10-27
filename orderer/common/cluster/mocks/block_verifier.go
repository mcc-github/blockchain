
package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"


type BlockVerifier struct {
	mock.Mock
}


func (_m *BlockVerifier) VerifyBlockSignature(sd []*common.SignedData) error {
	ret := _m.Called(sd)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*common.SignedData) error); ok {
		r0 = rf(sd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
