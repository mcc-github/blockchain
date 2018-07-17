/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"errors"
	"fmt"

	. "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	. "github.com/mcc-github/blockchain/core/handlers/endorsement/api/identities"
	"github.com/mcc-github/blockchain/protos/peer"
)







type DefaultEndorsementFactory struct {
}


func (*DefaultEndorsementFactory) New() Plugin {
	return &DefaultEndorsement{}
}


type DefaultEndorsement struct {
	SigningIdentityFetcher
}






func (e *DefaultEndorsement) Endorse(prpBytes []byte, sp *peer.SignedProposal) (*peer.Endorsement, []byte, error) {
	signer, err := e.SigningIdentityForRequest(sp)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("failed fetching signing identity: %v", err))
	}
	
	identityBytes, err := signer.Serialize()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("could not serialize the signing identity: %v", err))
	}

	
	signature, err := signer.Sign(append(prpBytes, identityBytes...))
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("could not sign the proposal response payload: %v", err))
	}
	endorsement := &peer.Endorsement{Signature: signature, Endorser: identityBytes}
	return endorsement, prpBytes, nil
}


func (e *DefaultEndorsement) Init(dependencies ...Dependency) error {
	for _, dep := range dependencies {
		sIDFetcher, isSigningIdentityFetcher := dep.(SigningIdentityFetcher)
		if !isSigningIdentityFetcher {
			continue
		}
		e.SigningIdentityFetcher = sIDFetcher
		return nil
	}
	return errors.New("could not find SigningIdentityFetcher in dependencies")
}


func NewPluginFactory() PluginFactory {
	return &DefaultEndorsementFactory{}
}
