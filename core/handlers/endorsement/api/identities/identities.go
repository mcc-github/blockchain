/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorsement

import (
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	"github.com/mcc-github/blockchain/protos/peer"
)


type SigningIdentity interface {
	
	
	Serialize() ([]byte, error)

	
	Sign([]byte) ([]byte, error)
}


type SigningIdentityFetcher interface {
	endorsement.Dependency
	
	SigningIdentityForRequest(*peer.SignedProposal) (SigningIdentity, error)
}
