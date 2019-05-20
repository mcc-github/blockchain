/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	"github.com/mcc-github/blockchain/protos/peer"
)

type NoOpEndorser struct {
}

func (*NoOpEndorser) Endorse(payload []byte, sp *peer.SignedProposal) (*peer.Endorsement, []byte, error) {
	return nil, payload, nil
}

func (*NoOpEndorser) Init(dependencies ...endorsement.Dependency) error {
	return nil
}

type NoOpEndorserFactory struct {
}

func (*NoOpEndorserFactory) New() endorsement.Plugin {
	return &NoOpEndorser{}
}


func NewPluginFactory() endorsement.PluginFactory {
	return &NoOpEndorserFactory{}
}
