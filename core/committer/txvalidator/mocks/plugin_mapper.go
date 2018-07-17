
package mocks

import mock "github.com/stretchr/testify/mock"
import txvalidator "github.com/mcc-github/blockchain/core/committer/txvalidator"
import validation "github.com/mcc-github/blockchain/core/handlers/validation/api"


type PluginMapper struct {
	mock.Mock
}


func (_m *PluginMapper) PluginFactoryByName(name txvalidator.PluginName) validation.PluginFactory {
	ret := _m.Called(name)

	var r0 validation.PluginFactory
	if rf, ok := ret.Get(0).(func(txvalidator.PluginName) validation.PluginFactory); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(validation.PluginFactory)
		}
	}

	return r0
}
