

package mocks

import (
	common "github.com/mcc-github/blockchain-protos-go/common"
	mock "github.com/stretchr/testify/mock"

	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
)


type Plugin struct {
	mock.Mock
}


func (_m *Plugin) Init(dependencies ...validation.Dependency) error {
	_va := make([]interface{}, len(dependencies))
	for _i := range dependencies {
		_va[_i] = dependencies[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...validation.Dependency) error); ok {
		r0 = rf(dependencies...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *Plugin) Validate(block *common.Block, namespace string, txPosition int, actionPosition int, contextData ...validation.ContextDatum) error {
	_va := make([]interface{}, len(contextData))
	for _i := range contextData {
		_va[_i] = contextData[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, block, namespace, txPosition, actionPosition)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*common.Block, string, int, int, ...validation.ContextDatum) error); ok {
		r0 = rf(block, namespace, txPosition, actionPosition, contextData...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
