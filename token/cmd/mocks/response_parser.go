
package mocks

import mock "github.com/stretchr/testify/mock"
import token "github.com/mcc-github/blockchain/token/cmd"


type ResponseParser struct {
	mock.Mock
}


func (_m *ResponseParser) ParseResponse(response token.StubResponse) error {
	ret := _m.Called(response)

	var r0 error
	if rf, ok := ret.Get(0).(func(token.StubResponse) error); ok {
		r0 = rf(response)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
