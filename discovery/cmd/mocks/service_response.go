

package mocks

import (
	blockchain_protos_godiscovery "github.com/mcc-github/blockchain-protos-go/discovery"
	discovery "github.com/mcc-github/blockchain/discovery/client"

	mock "github.com/stretchr/testify/mock"
)


type ServiceResponse struct {
	mock.Mock
}


func (_m *ServiceResponse) ForChannel(_a0 string) discovery.ChannelResponse {
	ret := _m.Called(_a0)

	var r0 discovery.ChannelResponse
	if rf, ok := ret.Get(0).(func(string) discovery.ChannelResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(discovery.ChannelResponse)
		}
	}

	return r0
}


func (_m *ServiceResponse) ForLocal() discovery.LocalResponse {
	ret := _m.Called()

	var r0 discovery.LocalResponse
	if rf, ok := ret.Get(0).(func() discovery.LocalResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(discovery.LocalResponse)
		}
	}

	return r0
}


func (_m *ServiceResponse) Raw() *blockchain_protos_godiscovery.Response {
	ret := _m.Called()

	var r0 *blockchain_protos_godiscovery.Response
	if rf, ok := ret.Get(0).(func() *blockchain_protos_godiscovery.Response); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*blockchain_protos_godiscovery.Response)
		}
	}

	return r0
}
