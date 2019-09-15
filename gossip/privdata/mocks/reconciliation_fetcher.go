

package mocks

import (
	common "github.com/mcc-github/blockchain/gossip/privdata/common"
	mock "github.com/stretchr/testify/mock"
)


type ReconciliationFetcher struct {
	mock.Mock
}


func (_m *ReconciliationFetcher) FetchReconciledItems(dig2collectionConfig common.Dig2CollectionConfig) (*common.FetchedPvtDataContainer, error) {
	ret := _m.Called(dig2collectionConfig)

	var r0 *common.FetchedPvtDataContainer
	if rf, ok := ret.Get(0).(func(common.Dig2CollectionConfig) *common.FetchedPvtDataContainer); ok {
		r0 = rf(dig2collectionConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.FetchedPvtDataContainer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Dig2CollectionConfig) error); ok {
		r1 = rf(dig2collectionConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
