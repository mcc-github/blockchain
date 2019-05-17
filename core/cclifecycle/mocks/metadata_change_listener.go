

package mocks

import chaincode "github.com/mcc-github/blockchain/common/chaincode"
import mock "github.com/stretchr/testify/mock"


type MetadataChangeListener struct {
	mock.Mock
}


func (_m *MetadataChangeListener) HandleMetadataUpdate(channel string, chaincodes chaincode.MetadataSet) {
	_m.Called(channel, chaincodes)
}
