
package mocks

import cmd "github.com/mcc-github/blockchain/token/cmd"
import mock "github.com/stretchr/testify/mock"
import time "time"
import token "github.com/mcc-github/blockchain/protos/token"


type Stub struct {
	mock.Mock
}


func (_m *Stub) Issue(tokensToIssue []*token.Token, waitTimeout time.Duration) (cmd.StubResponse, error) {
	ret := _m.Called(tokensToIssue, waitTimeout)

	var r0 cmd.StubResponse
	if rf, ok := ret.Get(0).(func([]*token.Token, time.Duration) cmd.StubResponse); ok {
		r0 = rf(tokensToIssue, waitTimeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cmd.StubResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*token.Token, time.Duration) error); ok {
		r1 = rf(tokensToIssue, waitTimeout)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Stub) ListTokens() (cmd.StubResponse, error) {
	ret := _m.Called()

	var r0 cmd.StubResponse
	if rf, ok := ret.Get(0).(func() cmd.StubResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cmd.StubResponse)
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


func (_m *Stub) Redeem(tokenIDs []*token.TokenId, quantity string, waitTimeout time.Duration) (cmd.StubResponse, error) {
	ret := _m.Called(tokenIDs, quantity, waitTimeout)

	var r0 cmd.StubResponse
	if rf, ok := ret.Get(0).(func([]*token.TokenId, string, time.Duration) cmd.StubResponse); ok {
		r0 = rf(tokenIDs, quantity, waitTimeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cmd.StubResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*token.TokenId, string, time.Duration) error); ok {
		r1 = rf(tokenIDs, quantity, waitTimeout)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Stub) Setup(configFilePath string, channel string, mspPath string, mspID string) error {
	ret := _m.Called(configFilePath, channel, mspPath, mspID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) error); ok {
		r0 = rf(configFilePath, channel, mspPath, mspID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *Stub) Transfer(tokenIDs []*token.TokenId, shares []*token.RecipientShare, waitTimeout time.Duration) (cmd.StubResponse, error) {
	ret := _m.Called(tokenIDs, shares, waitTimeout)

	var r0 cmd.StubResponse
	if rf, ok := ret.Get(0).(func([]*token.TokenId, []*token.RecipientShare, time.Duration) cmd.StubResponse); ok {
		r0 = rf(tokenIDs, shares, waitTimeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cmd.StubResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*token.TokenId, []*token.RecipientShare, time.Duration) error); ok {
		r1 = rf(tokenIDs, shares, waitTimeout)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
