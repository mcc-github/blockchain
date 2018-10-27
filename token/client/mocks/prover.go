
package mocks

import blockchaintoken "github.com/mcc-github/blockchain/token"
import mock "github.com/stretchr/testify/mock"
import token "github.com/mcc-github/blockchain/protos/token"


type Prover struct {
	mock.Mock
}


func (_m *Prover) RequestImport(tokensToIssue []*token.TokenToIssue, signingIdentity blockchaintoken.SigningIdentity) ([]byte, error) {
	ret := _m.Called(tokensToIssue, signingIdentity)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]*token.TokenToIssue, blockchaintoken.SigningIdentity) []byte); ok {
		r0 = rf(tokensToIssue, signingIdentity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*token.TokenToIssue, blockchaintoken.SigningIdentity) error); ok {
		r1 = rf(tokensToIssue, signingIdentity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *Prover) RequestTransfer(tokenIDs [][]byte, shares []*token.RecipientTransferShare, signingIdentity blockchaintoken.SigningIdentity) ([]byte, error) {
	ret := _m.Called(tokenIDs, shares, signingIdentity)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([][]byte, []*token.RecipientTransferShare, blockchaintoken.SigningIdentity) []byte); ok {
		r0 = rf(tokenIDs, shares, signingIdentity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([][]byte, []*token.RecipientTransferShare, blockchaintoken.SigningIdentity) error); ok {
		r1 = rf(tokenIDs, shares, signingIdentity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
