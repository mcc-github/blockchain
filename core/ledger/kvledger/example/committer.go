/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package example

import (
	"github.com/mcc-github/blockchain/core/ledger"

	"github.com/mcc-github/blockchain/protos/common"
)


type Committer struct {
	ledger ledger.PeerLedger
}


func ConstructCommitter(ledger ledger.PeerLedger) *Committer {
	return &Committer{ledger}
}


func (c *Committer) Commit(rawBlock *common.Block) error {
	logger.Debugf("Committer validating the block...")
	if err := c.ledger.CommitWithPvtData(&ledger.BlockAndPvtData{Block: rawBlock}); err != nil {
		return err
	}
	return nil
}
