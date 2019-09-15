

package mocks

import (
	common "github.com/mcc-github/blockchain-protos-go/common"
	mock "github.com/stretchr/testify/mock"

	protoutil "github.com/mcc-github/blockchain/protoutil"
)


type BlockVerifier struct {
	mock.Mock
}


func (_m *BlockVerifier) VerifyBlockSignature(sd []*protoutil.SignedData, config *common.ConfigEnvelope) error {
	ret := _m.Called(sd, config)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*protoutil.SignedData, *common.ConfigEnvelope) error); ok {
		r0 = rf(sd, config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
