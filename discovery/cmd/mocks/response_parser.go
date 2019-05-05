

package mocks

import discovery "github.com/mcc-github/blockchain/discovery/cmd"
import mock "github.com/stretchr/testify/mock"


type ResponseParser struct {
	mock.Mock
}


func (_m *ResponseParser) ParseResponse(channel string, response discovery.ServiceResponse) error {
	ret := _m.Called(channel, response)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, discovery.ServiceResponse) error); ok {
		r0 = rf(channel, response)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
