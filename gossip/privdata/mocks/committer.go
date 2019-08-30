

package mocks

import common "github.com/mcc-github/blockchain-protos-go/common"
import ledger "github.com/mcc-github/blockchain/core/ledger"
import mock "github.com/stretchr/testify/mock"


type Committer struct {
	mock.Mock
}


func (_m *Committer) Close() {
	_m.Called()
}


func (_m *Committer) CommitLegacy(blockAndPvtData *ledger.BlockAndPvtData, commitOpts *ledger.CommitOptions) error {
	ret := _m.Called(blockAndPvtData, commitOpts)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ledger.BlockAndPvtData, *ledger.CommitOptions) error); ok {
		r0 = rf(blockAndPvtData, commitOpts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *Committer) CommitPvtDataOfOldBlocks(blockPvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error) {
	ret := _m.Called(blockPvtData)

	var r0 []*ledger.PvtdataHashMismatch
	if rf, ok := ret.Get(0).(func([]*ledger.BlockPvtData) []*ledger.PvtdataHashMismatch); ok {
		r0 = rf(blockPvtData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ledger.PvtdataHashMismatch)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*ledger.BlockPvtData) error); ok {
		r1 = rf(blockPvtData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Committer) DoesPvtDataInfoExistInLedger(blockNum uint64) (bool, error) {
	ret := _m.Called(blockNum)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64) bool); ok {
		r0 = rf(blockNum)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(blockNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Committer) GetBlocks(blockSeqs []uint64) []*common.Block {
	ret := _m.Called(blockSeqs)

	var r0 []*common.Block
	if rf, ok := ret.Get(0).(func([]uint64) []*common.Block); ok {
		r0 = rf(blockSeqs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*common.Block)
		}
	}

	return r0
}


func (_m *Committer) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	ret := _m.Called()

	var r0 ledger.ConfigHistoryRetriever
	if rf, ok := ret.Get(0).(func() ledger.ConfigHistoryRetriever); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.ConfigHistoryRetriever)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Committer) GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error) {
	ret := _m.Called()

	var r0 ledger.MissingPvtDataTracker
	if rf, ok := ret.Get(0).(func() ledger.MissingPvtDataTracker); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.MissingPvtDataTracker)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Committer) GetPvtDataAndBlockByNum(seqNum uint64) (*ledger.BlockAndPvtData, error) {
	ret := _m.Called(seqNum)

	var r0 *ledger.BlockAndPvtData
	if rf, ok := ret.Get(0).(func(uint64) *ledger.BlockAndPvtData); ok {
		r0 = rf(seqNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ledger.BlockAndPvtData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(seqNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Committer) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	ret := _m.Called(blockNum, filter)

	var r0 []*ledger.TxPvtData
	if rf, ok := ret.Get(0).(func(uint64, ledger.PvtNsCollFilter) []*ledger.TxPvtData); ok {
		r0 = rf(blockNum, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ledger.TxPvtData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, ledger.PvtNsCollFilter) error); ok {
		r1 = rf(blockNum, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Committer) LedgerHeight() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
