/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package filter

import (
	"context"

	"github.com/mcc-github/blockchain/core/handlers/auth"
	"github.com/mcc-github/blockchain/protos/peer"
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
