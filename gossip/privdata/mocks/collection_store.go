

package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import ledger "github.com/mcc-github/blockchain/core/ledger"
import mock "github.com/stretchr/testify/mock"
import peer "github.com/mcc-github/blockchain/protos/peer"
import privdata "github.com/mcc-github/blockchain/core/common/privdata"


type CollectionStore struct {
	mock.Mock
}


func (_m *CollectionStore) AccessFilter(channelName string, collectionPolicyConfig *common.CollectionPolicyConfig) (privdata.Filter, error) {
	ret := _m.Called(channelName, collectionPolicyConfig)

	var r0 privdata.Filter
	if rf, ok := ret.Get(0).(func(string, *common.CollectionPolicyConfig) privdata.Filter); ok {
		r0 = rf(channelName, collectionPolicyConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privdata.Filter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *common.CollectionPolicyConfig) error); ok {
		r1 = rf(channelName, collectionPolicyConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *CollectionStore) HasReadAccess(collCriteria common.CollectionCriteria, signedProposal *peer.SignedProposal, qe ledger.QueryExecutor) (bool, error) {
	ret := _m.Called(collCriteria, signedProposal, qe)

	var r0 bool
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria, *peer.SignedProposal, ledger.QueryExecutor) bool); ok {
		r0 = rf(collCriteria, signedProposal, qe)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria, *peer.SignedProposal, ledger.QueryExecutor) error); ok {
		r1 = rf(collCriteria, signedProposal, qe)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *CollectionStore) RetrieveCollection(collCriteria common.CollectionCriteria) (privdata.Collection, error) {
	ret := _m.Called(collCriteria)

	var r0 privdata.Collection
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria) privdata.Collection); ok {
		r0 = rf(collCriteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privdata.Collection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria) error); ok {
		r1 = rf(collCriteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *CollectionStore) RetrieveCollectionAccessPolicy(collCriteria common.CollectionCriteria) (privdata.CollectionAccessPolicy, error) {
	ret := _m.Called(collCriteria)

	var r0 privdata.CollectionAccessPolicy
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria) privdata.CollectionAccessPolicy); ok {
		r0 = rf(collCriteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privdata.CollectionAccessPolicy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria) error); ok {
		r1 = rf(collCriteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *CollectionStore) RetrieveCollectionConfigPackage(collCriteria common.CollectionCriteria) (*common.CollectionConfigPackage, error) {
	ret := _m.Called(collCriteria)

	var r0 *common.CollectionConfigPackage
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria) *common.CollectionConfigPackage); ok {
		r0 = rf(collCriteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.CollectionConfigPackage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria) error); ok {
		r1 = rf(collCriteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *CollectionStore) RetrieveCollectionPersistenceConfigs(collCriteria common.CollectionCriteria) (privdata.CollectionPersistenceConfigs, error) {
	ret := _m.Called(collCriteria)

	var r0 privdata.CollectionPersistenceConfigs
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria) privdata.CollectionPersistenceConfigs); ok {
		r0 = rf(collCriteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privdata.CollectionPersistenceConfigs)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria) error); ok {
		r1 = rf(collCriteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
