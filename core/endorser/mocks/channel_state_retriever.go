
package mocks

import endorser "github.com/mcc-github/blockchain/core/endorser"
import mock "github.com/stretchr/testify/mock"


type ChannelStateRetriever struct {
	mock.Mock
}


func (_m *ChannelStateRetriever) NewQueryCreator(channel string) (endorser.QueryCreator, error) {
	ret := _m.Called(channel)

	var r0 endorser.QueryCreator
	if rf, ok := ret.Get(0).(func(string) endorser.QueryCreator); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(endorser.QueryCreator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}