

package mock

import ccprovider "github.com/mcc-github/blockchain/core/common/ccprovider"

import ledger "github.com/mcc-github/blockchain/core/ledger"
import mock "github.com/stretchr/testify/mock"


type Lifecycle struct {
	mock.Mock
}


func (_m *Lifecycle) ChaincodeContainerInfo(channelID string, chaincodeName string, qe ledger.SimpleQueryExecutor) (*ccprovider.ChaincodeContainerInfo, error) {
	ret := _m.Called(channelID, chaincodeName, qe)

	var r0 *ccprovider.ChaincodeContainerInfo
	if rf, ok := ret.Get(0).(func(string, string, ledger.SimpleQueryExecutor) *ccprovider.ChaincodeContainerInfo); ok {
		r0 = rf(channelID, chaincodeName, qe)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ccprovider.ChaincodeContainerInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, ledger.SimpleQueryExecutor) error); ok {
		r1 = rf(channelID, chaincodeName, qe)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Lifecycle) ChaincodeDefinition(channelID string, chaincodeName string, qe ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	ret := _m.Called(channelID, chaincodeName, qe)

	var r0 ccprovider.ChaincodeDefinition
	if rf, ok := ret.Get(0).(func(string, string, ledger.SimpleQueryExecutor) ccprovider.ChaincodeDefinition); ok {
		r0 = rf(channelID, chaincodeName, qe)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ccprovider.ChaincodeDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, ledger.SimpleQueryExecutor) error); ok {
		r1 = rf(channelID, chaincodeName, qe)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
