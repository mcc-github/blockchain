
package mocks

import mock "github.com/stretchr/testify/mock"
import validation "github.com/mcc-github/blockchain/core/handlers/validation/api"


type PluginFactory struct {
	mock.Mock
}


func (_m *PluginFactory) New() validation.Plugin {
	ret := _m.Called()

	var r0 validation.Plugin
	if rf, ok := ret.Get(0).(func() validation.Plugin); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(validation.Plugin)
		}
	}

	return r0
}
