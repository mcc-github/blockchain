

package mocks

import common "github.com/mcc-github/blockchain-protos-go/common"
import etcdraft "github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
import mock "github.com/stretchr/testify/mock"


type InactiveChainRegistry struct {
	mock.Mock
}


func (_m *InactiveChainRegistry) TrackChain(chainName string, genesisBlock *common.Block, createChain etcdraft.CreateChainCallback) {
	_m.Called(chainName, genesisBlock, createChain)
}
