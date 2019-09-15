

package mocks

import (
	common "github.com/mcc-github/blockchain-protos-go/common"
	blockledger "github.com/mcc-github/blockchain/common/ledger/blockledger"

	mock "github.com/stretchr/testify/mock"

	orderer "github.com/mcc-github/blockchain-protos-go/orderer"
)


type ReadWriter struct {
	mock.Mock
}


func (_m *ReadWriter) Append(block *common.Block) error {
	ret := _m.Called(block)

	var r0 error
	if rf, ok := ret.Get(0).(func(*common.Block) error); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *ReadWriter) Height() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}


func (_m *ReadWriter) Iterator(startType *orderer.SeekPosition) (blockledger.Iterator, uint64) {
	ret := _m.Called(startType)

	var r0 blockledger.Iterator
	if rf, ok := ret.Get(0).(func(*orderer.SeekPosition) blockledger.Iterator); ok {
		r0 = rf(startType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(blockledger.Iterator)
		}
	}

	var r1 uint64
	if rf, ok := ret.Get(1).(func(*orderer.SeekPosition) uint64); ok {
		r1 = rf(startType)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	return r0, r1
}
