/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


package mocks

import endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
import mock "github.com/stretchr/testify/mock"


type PluginFactory struct {
	mock.Mock
}


func (_m *PluginFactory) New() endorsement.Plugin {
	ret := _m.Called()

	var r0 endorsement.Plugin
	if rf, ok := ret.Get(0).(func() endorsement.Plugin); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(endorsement.Plugin)
		}
	}

	return r0
}
