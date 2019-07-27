/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package committer

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
)

var logger = flogging.MustGetLogger("committer")








type PeerLedgerSupport interface {
	GetPvtDataAndBlockByNum(blockNum uint64, filter ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error)

	GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)

	CommitWithPvtData(blockAndPvtdata *ledger.BlockAndPvtData) error

	CommitPvtDataOfOldBlocks(blockPvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error)

	GetBlockchainInfo() (*common.BlockchainInfo, error)

	GetBlockByNumber(blockNumber uint64) (*common.Block, error)

	GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error)

	GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error)

	Close()
}




type LedgerCommitter struct {
	PeerLedgerSupport
}



func NewLedgerCommitter(ledger PeerLedgerSupport) *LedgerCommitter {
	return &LedgerCommitter{PeerLedgerSupport: ledger}
}


func (lc *LedgerCommitter) CommitWithPvtData(blockAndPvtData *ledger.BlockAndPvtData) error {
	
	if err := lc.PeerLedgerSupport.CommitWithPvtData(blockAndPvtData); err != nil {
		return err
	}

	return nil
}


func (lc *LedgerCommitter) GetPvtDataAndBlockByNum(seqNum uint64) (*ledger.BlockAndPvtData, error) {
	return lc.PeerLedgerSupport.GetPvtDataAndBlockByNum(seqNum, nil)
}


func (lc *LedgerCommitter) LedgerHeight() (uint64, error) {
	info, err := lc.GetBlockchainInfo()
	if err != nil {
		logger.Errorf("Cannot get blockchain info, %s", info)
		return 0, err
	}

	return info.Height, nil
}


func (lc *LedgerCommitter) GetBlocks(blockSeqs []uint64) []*common.Block {
	var blocks []*common.Block

	for _, seqNum := range blockSeqs {
		if blck, err := lc.GetBlockByNumber(seqNum); err != nil {
			logger.Errorf("Not able to acquire block num %d, from the ledger skipping...", seqNum)
			continue
		} else {
			logger.Debug("Appending next block with seqNum = ", seqNum, " to the resulting set")
			blocks = append(blocks, blck)
		}
	}

	return blocks
}
