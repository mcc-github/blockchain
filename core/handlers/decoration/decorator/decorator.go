/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package decorator

import (
	"github.com/mcc-github/blockchain/core/handlers/decoration"
	"github.com/mcc-github/blockchain/protos/peer"
)


func NewDecorator() decoration.Decorator {
	return &decorator{}
}

type decorator struct {
}


func (d *decorator) Decorate(proposal *peer.Proposal, input *peer.ChaincodeInput) *peer.ChaincodeInput {
	return input
}
