

package mocks

import discovery "github.com/mcc-github/blockchain/discovery/client"
import mock "github.com/stretchr/testify/mock"
import protosdiscovery "github.com/mcc-github/blockchain-protos-go/discovery"


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


func (_m *ServiceResponse) Raw() *protosdiscovery.Response {
	ret := _m.Called()

	var r0 *protosdiscovery.Response
	if rf, ok := ret.Get(0).(func() *protosdiscovery.Response); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*protosdiscovery.Response)
		}
	}

	return r0
}
