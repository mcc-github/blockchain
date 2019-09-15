

package mocks

import (
	transientstore "github.com/mcc-github/blockchain/core/transientstore"
	mock "github.com/stretchr/testify/mock"
)


type TransientStoreRetriever struct {
	mock.Mock
}


func (_m *TransientStoreRetriever) StoreForChannel(channel string) transientstore.Store {
	ret := _m.Called(channel)

	var r0 transientstore.Store
	if rf, ok := ret.Get(0).(func(string) transientstore.Store); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transientstore.Store)
		}
	}

	return r0
}
