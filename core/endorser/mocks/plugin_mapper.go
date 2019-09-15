

package mocks

import (
	endorser "github.com/mcc-github/blockchain/core/endorser"
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"

	mock "github.com/stretchr/testify/mock"
)


type PluginMapper struct {
	mock.Mock
}


func (_m *PluginMapper) PluginFactoryByName(name endorser.PluginName) endorsement.PluginFactory {
	ret := _m.Called(name)

	var r0 endorsement.PluginFactory
	if rf, ok := ret.Get(0).(func(endorser.PluginName) endorsement.PluginFactory); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(endorsement.PluginFactory)
		}
	}

	return r0
}
