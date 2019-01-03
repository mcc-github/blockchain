

package mocks

import common "github.com/mcc-github/blockchain/protos/common"

import mock "github.com/stretchr/testify/mock"


type InactiveChainRegistry struct {
	mock.Mock
}


func (_m *InactiveChainRegistry) TrackChain(chainName string, genesisBlock *common.Block, createChainCallback func()) {
	_m.Called(chainName, genesisBlock, createChainCallback)
}
