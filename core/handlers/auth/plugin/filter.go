/*
Copyright IBM Corp, SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/mcc-github/blockchain/core/handlers/auth"
	"github.com/mcc-github/blockchain/protos/peer"
	"golang.org/x/net/context"
)


func NewFilter() auth.Filter {
	return &filter{}
}

type filter struct {
	next peer.EndorserServer
}


func (f *filter) Init(next peer.EndorserServer) {
	f.next = next
}


func (f *filter) ProcessProposal(ctx context.Context, signedProp *peer.SignedProposal) (*peer.ProposalResponse, error) {
	return f.next.ProcessProposal(ctx, signedProp)
}

func main() {
}
