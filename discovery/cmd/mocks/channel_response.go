
package mocks

import client "github.com/mcc-github/blockchain/discovery/client"
import discovery "github.com/mcc-github/blockchain/protos/discovery"
import mock "github.com/stretchr/testify/mock"


type ChannelResponse struct {
	mock.Mock
}


func (_m *ChannelResponse) Config() (*discovery.ConfigResult, error) {
	ret := _m.Called()

	var r0 *discovery.ConfigResult
	if rf, ok := ret.Get(0).(func() *discovery.ConfigResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discovery.ConfigResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *ChannelResponse) Endorsers(invocationChain client.InvocationChain, ps client.PrioritySelector, ef client.ExclusionFilter) (client.Endorsers, error) {
	ret := _m.Called(invocationChain, ps, ef)

	var r0 client.Endorsers
	if rf, ok := ret.Get(0).(func(client.InvocationChain, client.PrioritySelector, client.ExclusionFilter) client.Endorsers); ok {
		r0 = rf(invocationChain, ps, ef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.Endorsers)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(client.InvocationChain, client.PrioritySelector, client.ExclusionFilter) error); ok {
		r1 = rf(invocationChain, ps, ef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *ChannelResponse) Peers(invocationChain ...*discovery.ChaincodeCall) ([]*client.Peer, error) {
	_va := make([]interface{}, len(invocationChain))
	for _i := range invocationChain {
		_va[_i] = invocationChain[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*client.Peer
	if rf, ok := ret.Get(0).(func(...*discovery.ChaincodeCall) []*client.Peer); ok {
		r0 = rf(invocationChain...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*client.Peer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...*discovery.ChaincodeCall) error); ok {
		r1 = rf(invocationChain...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
