

package mocks

import (
	ledger "github.com/mcc-github/blockchain/core/ledger"
	mock "github.com/stretchr/testify/mock"
)


type MissingPvtDataTracker struct {
	mock.Mock
}


func (_m *MissingPvtDataTracker) GetMissingPvtDataInfoForMostRecentBlocks(maxBlocks int) (ledger.MissingPvtDataInfo, error) {
	ret := _m.Called(maxBlocks)

	var r0 ledger.MissingPvtDataInfo
	if rf, ok := ret.Get(0).(func(int) ledger.MissingPvtDataInfo); ok {
		r0 = rf(maxBlocks)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.MissingPvtDataInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(maxBlocks)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
