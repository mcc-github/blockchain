/*
Copyright IBM Corp, SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package decoration

import (
	"github.com/mcc-github/blockchain/protos/peer"
)


type Decorator interface {
	
	Decorate(proposal *peer.Proposal, input *peer.ChaincodeInput) *peer.ChaincodeInput
}


func Apply(proposal *peer.Proposal, input *peer.ChaincodeInput,
	decorators ...Decorator) *peer.ChaincodeInput {
	for _, decorator := range decorators {
		input = decorator.Decorate(proposal, input)
	}

	return input
}
