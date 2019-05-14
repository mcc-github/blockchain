

package mocks

import api "github.com/mcc-github/blockchain/gossip/api"
import mock "github.com/stretchr/testify/mock"


type SecurityAdvisor struct {
	mock.Mock
}


func (_m *SecurityAdvisor) OrgByPeerIdentity(_a0 api.PeerIdentityType) api.OrgIdentityType {
	ret := _m.Called(_a0)

	var r0 api.OrgIdentityType
	if rf, ok := ret.Get(0).(func(api.PeerIdentityType) api.OrgIdentityType); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.OrgIdentityType)
		}
	}

	return r0
}
