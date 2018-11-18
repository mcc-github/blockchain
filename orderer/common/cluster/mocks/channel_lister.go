
package mocks

import mock "github.com/stretchr/testify/mock"


type ChannelLister struct {
	mock.Mock
}


func (_m *ChannelLister) Channels() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}


func (_m *ChannelLister) Close() {
	_m.Called()
}
