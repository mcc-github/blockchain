/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorsement

import (
	"github.com/mcc-github/blockchain-protos-go/peer"
)


type Argument interface {
	Dependency
	
	Arg() []byte
}


type Dependency interface {
}


type Plugin interface {
	
	
	
	
	
	Endorse(payload []byte, sp *peer.SignedProposal) (*peer.Endorsement, []byte, error)

	
	Init(dependencies ...Dependency) error
}


type PluginFactory interface {
	New() Plugin
}
