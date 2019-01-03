

package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"


type ChainReplicator struct {
	mock.Mock
}


func (_m *ChainReplicator) ReplicateChains(lastConfigBlock *common.Block, chains []string) []string {
	ret := _m.Called(lastConfigBlock, chains)

	var r0 []string
	if rf, ok := ret.Get(0).(func(*common.Block, []string) []string); ok {
		r0 = rf(lastConfigBlock, chains)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
