
package mocks

import api "github.com/mcc-github/blockchain/gossip/api"
import common "github.com/mcc-github/blockchain/gossip/common"

import gossipdiscovery "github.com/mcc-github/blockchain/gossip/discovery"
import mock "github.com/stretchr/testify/mock"


type GossipSupport struct {
	mock.Mock
}


func (_m *GossipSupport) ChannelExists(channel string) bool {
	ret := _m.Called(channel)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(channel)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}


func (_m *GossipSupport) IdentityInfo() api.PeerIdentitySet {
	ret := _m.Called()

	var r0 api.PeerIdentitySet
	if rf, ok := ret.Get(0).(func() api.PeerIdentitySet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.PeerIdentitySet)
		}
	}

	return r0
}


func (_m *GossipSupport) Peers() gossipdiscovery.Members {
	ret := _m.Called()

	var r0 gossipdiscovery.Members
	if rf, ok := ret.Get(0).(func() gossipdiscovery.Members); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gossipdiscovery.Members)
		}
	}

	return r0
}


func (_m *GossipSupport) PeersOfChannel(_a0 common.ChainID) gossipdiscovery.Members {
	ret := _m.Called(_a0)

	var r0 gossipdiscovery.Members
	if rf, ok := ret.Get(0).(func(common.ChainID) gossipdiscovery.Members); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gossipdiscovery.Members)
		}
	}

	return r0
}
