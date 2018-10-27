

package mocks

import cluster "github.com/mcc-github/blockchain/orderer/common/cluster"

import mock "github.com/stretchr/testify/mock"


type Configurator struct {
	mock.Mock
}


func (_m *Configurator) Configure(channel string, newNodes []cluster.RemoteNode) {
	_m.Called(channel, newNodes)
}
