
package mocks

import chaincode "github.com/mcc-github/blockchain/common/chaincode"
import mock "github.com/stretchr/testify/mock"


type LifeCycleChangeListener struct {
	mock.Mock
}


func (_m *LifeCycleChangeListener) LifeCycleChangeListener(channel string, chaincodes chaincode.MetadataSet) {
	_m.Called(channel, chaincodes)
}
