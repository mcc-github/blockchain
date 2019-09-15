

package mocks

import (
	plugin "github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	mock "github.com/stretchr/testify/mock"

	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
)


type Mapper struct {
	mock.Mock
}


func (_m *Mapper) FactoryByName(name plugin.Name) validation.PluginFactory {
	ret := _m.Called(name)

	var r0 validation.PluginFactory
	if rf, ok := ret.Get(0).(func(plugin.Name) validation.PluginFactory); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(validation.PluginFactory)
		}
	}

	return r0
}
