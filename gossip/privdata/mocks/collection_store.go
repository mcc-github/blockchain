
package mocks

import common "github.com/mcc-github/blockchain/protos/common"
import mock "github.com/stretchr/testify/mock"
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


func (_m *CollectionStore) RetrieveCollection(_a0 common.CollectionCriteria) (privdata.Collection, error) {
	ret := _m.Called(_a0)

	var r0 privdata.Collection
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria) privdata.Collection); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privdata.Collection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *CollectionStore) RetrieveCollectionAccessPolicy(_a0 common.CollectionCriteria) (privdata.CollectionAccessPolicy, error) {
	ret := _m.Called(_a0)

	var r0 privdata.CollectionAccessPolicy
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria) privdata.CollectionAccessPolicy); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privdata.CollectionAccessPolicy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *CollectionStore) RetrieveCollectionConfigPackage(_a0 common.CollectionCriteria) (*common.CollectionConfigPackage, error) {
	ret := _m.Called(_a0)

	var r0 *common.CollectionConfigPackage
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria) *common.CollectionConfigPackage); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.CollectionConfigPackage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *CollectionStore) RetrieveCollectionPersistenceConfigs(cc common.CollectionCriteria) (privdata.CollectionPersistenceConfigs, error) {
	ret := _m.Called(cc)

	var r0 privdata.CollectionPersistenceConfigs
	if rf, ok := ret.Get(0).(func(common.CollectionCriteria) privdata.CollectionPersistenceConfigs); ok {
		r0 = rf(cc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privdata.CollectionPersistenceConfigs)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.CollectionCriteria) error); ok {
		r1 = rf(cc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
