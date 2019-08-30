/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorsement

import (
	"github.com/mcc-github/blockchain-protos-go/peer"
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
)


type SigningIdentity interface {
	
	
	Serialize() ([]byte, error)

	
	Sign([]byte) ([]byte, error)
}


type SigningIdentityFetcher interface {
	endorsement.Dependency
	
	SigningIdentityForRequest(*peer.SignedProposal) (SigningIdentity, error)
}
