

package mocks

import chaincode "github.com/mcc-github/blockchain/common/chaincode"
import mock "github.com/stretchr/testify/mock"


type LifecycleChangeListener struct {
	mock.Mock
}


func (_m *LifecycleChangeListener) HandleMetadataUpdate(channel string, chaincodes chaincode.MetadataSet) {
	_m.Called(channel, chaincodes)
}
