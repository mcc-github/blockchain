
package mocks

import mock "github.com/stretchr/testify/mock"
import token "github.com/mcc-github/blockchain/protos/token"


type Loader struct {
	mock.Mock
}


func (_m *Loader) Shares(s string) ([]*token.RecipientShare, error) {
	ret := _m.Called(s)

	var r0 []*token.RecipientShare
	if rf, ok := ret.Get(0).(func(string) []*token.RecipientShare); ok {
		r0 = rf(s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*token.RecipientShare)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(s)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Loader) TokenIDs(s string) ([]*token.TokenId, error) {
	ret := _m.Called(s)

	var r0 []*token.TokenId
	if rf, ok := ret.Get(0).(func(string) []*token.TokenId); ok {
		r0 = rf(s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*token.TokenId)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(s)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Loader) TokenOwner(s string) (*token.TokenOwner, error) {
	ret := _m.Called(s)

	var r0 *token.TokenOwner
	if rf, ok := ret.Get(0).(func(string) *token.TokenOwner); ok {
		r0 = rf(s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token.TokenOwner)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(s)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
