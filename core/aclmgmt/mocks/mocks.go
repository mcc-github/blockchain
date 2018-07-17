/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"testing"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/stretchr/testify/mock"
)

type MockACLProvider struct {
	
	
	mock *mock.Mock
}


func (m *MockACLProvider) Reset() {
	m.mock = &mock.Mock{}
}

func (m *MockACLProvider) CheckACL(resName string, channelID string, idinfo interface{}) error {
	args := m.mock.Called(resName, channelID, idinfo)
	return args.Error(0)
}

func (m *MockACLProvider) GenerateSimulationResults(txEnvelop *common.Envelope, simulator ledger.TxSimulator, initializingLedger bool) error {
	return nil
}


func (m *MockACLProvider) On(methodName string, arguments ...interface{}) *mock.Call {
	return m.mock.On(methodName, arguments...)
}


func (m *MockACLProvider) AssertExpectations(t *testing.T) {
	m.mock.AssertExpectations(t)
}
