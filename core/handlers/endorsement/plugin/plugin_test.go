/*
Copyright Hitachi, Ltd. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main_test

import (
	"errors"
	"testing"

	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/core/endorser/mocks"
	mocks2 "github.com/mcc-github/blockchain/core/handlers/endorsement/builtin/mocks"
	plgn "github.com/mcc-github/blockchain/core/handlers/endorsement/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEndorsementPlugin(t *testing.T) {
	factory := plgn.NewPluginFactory()
	plugin := factory.New()
	dependency := &struct{}{}
	err := plugin.Init(dependency)
	assert.EqualError(t, err, "could not find SigningIdentityFetcher in dependencies")

	sif := &mocks.SigningIdentityFetcher{}
	err1 := plugin.Init(sif)
	assert.NoError(t, err1)

	
	
	sif.On("SigningIdentityForRequest", mock.Anything).Return(nil, errors.New("signingIdentityForRequestReturnsError")).Once()
	endorsement2, prepBytes2, err2 := plugin.Endorse(nil, nil)
	assert.Nil(t, endorsement2)
	assert.Nil(t, prepBytes2)
	assert.EqualError(t, err2, "failed fetching signing identity: signingIdentityForRequestReturnsError")

	
	sid := &mocks2.SigningIdentity{}
	sid.On("Serialize").Return(nil, errors.New("serializeReturnsError")).Once()
	sif.On("SigningIdentityForRequest", mock.Anything).Return(sid, nil).Once()
	endorsement3, prepBytes3, err3 := plugin.Endorse(nil, nil)
	assert.Nil(t, endorsement3)
	assert.Nil(t, prepBytes3)
	assert.EqualError(t, err3, "could not serialize the signing identity: serializeReturnsError")

	
	sid.On("Serialize").Return([]byte("Endorser4"), nil).Once()
	sid.On("Sign", mock.Anything).Return(nil, errors.New("signReturnsError")).Once()
	sif.On("SigningIdentityForRequest", mock.Anything).Return(sid, nil).Once()
	endorsement4, prepBytes4, err4 := plugin.Endorse([]byte("prpBytes4"), nil)
	assert.Nil(t, endorsement4)
	assert.Nil(t, prepBytes4)
	assert.EqualError(t, err4, "could not sign the proposal response payload: signReturnsError")

	
	sid.On("Serialize").Return([]byte("Endorser5"), nil).Once()
	sid.On("Sign", mock.Anything).Return([]byte("Signature5"), nil).Once()
	sif.On("SigningIdentityForRequest", mock.Anything).Return(sid, nil).Once()
	endorsement5, prpBytes5, err5 := plugin.Endorse([]byte("prpBytes5"), nil)
	expected5 := &peer.Endorsement{Signature: []byte("Signature5"), Endorser: []byte("Endorser5")}
	assert.NoError(t, err5)
	assert.Equal(t, expected5, endorsement5)
	assert.Equal(t, []byte("prpBytes5"), prpBytes5)
}
