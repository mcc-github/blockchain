/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validator

import (
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/stretchr/testify/mock"
)


type MockValidator struct {
	mock.Mock
}


func (m *MockValidator) Validate(block *common.Block) error {
	if len(m.ExpectedCalls) == 0 {
		return nil
	}
	return m.Called().Error(0)
}


type MockVsccValidator struct {
}


func (v *MockVsccValidator) VSCCValidateTx(seq int, payload *common.Payload, envBytes []byte, block *common.Block) (error, peer.TxValidationCode) {
	return nil, peer.TxValidationCode_VALID
}
