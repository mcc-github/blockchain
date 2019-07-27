/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/token/ledger"
	"github.com/pkg/errors"
)



type PeerLedgerManager struct {
	Peer *peer.Peer
}

func (p *PeerLedgerManager) GetLedgerReader(channel string) (ledger.LedgerReader, error) {
	l := p.Peer.GetLedger(channel)
	if l == nil {
		return nil, errors.Errorf("ledger not found for channel %s", channel)
	}

	return l.NewQueryExecutor()
}
