/*
Copyright IBM Corp, SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/core/handlers/decoration"
)


func NewDecorator() decoration.Decorator {
	return &decorator{}
}

type decorator struct {
}


func (d *decorator) Decorate(proposal *peer.Proposal, input *peer.ChaincodeInput) *peer.ChaincodeInput {
	return input
}

func main() {
}
