/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package committer

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/events/producer"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
)






var logger *logging.Logger 

func init() {
	logger = flogging.MustGetLogger("committer")
}



type PeerLedgerSupport interface {
	GetPvtDataAndBlockByNum(blockNum uint64, filter ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error)

	GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)

	CommitWithPvtData(blockAndPvtdata *ledger.BlockAndPvtData) error

	GetBlockchainInfo() (*common.BlockchainInfo, error)

	GetBlockByNumber(blockNumber uint64) (*common.Block, error)

	GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error)

	Close()
}




type LedgerCommitter struct {
	PeerLedgerSupport
	eventer ConfigBlockEventer
}



type ConfigBlockEventer func(block *common.Block) error



func NewLedgerCommitter(ledger PeerLedgerSupport) *LedgerCommitter {
	return NewLedgerCommitterReactive(ledger, func(_ *common.Block) error { return nil })
}




func NewLedgerCommitterReactive(ledger PeerLedgerSupport, eventer ConfigBlockEventer) *LedgerCommitter {
	return &LedgerCommitter{PeerLedgerSupport: ledger, eventer: eventer}
}



func (lc *LedgerCommitter) preCommit(block *common.Block) error {
	
	if utils.IsConfigBlock(block) {
		logger.Debug("Received configuration update, calling CSCC ConfigUpdate")
		if err := lc.eventer(block); err != nil {
			return errors.WithMessage(err, "could not update CSCC with new configuration update")
		}
	}
	return nil
}


func (lc *LedgerCommitter) CommitWithPvtData(blockAndPvtData *ledger.BlockAndPvtData) error {
	
	
	if err := lc.preCommit(blockAndPvtData.Block); err != nil {
		return err
	}

	
	if err := lc.PeerLedgerSupport.CommitWithPvtData(blockAndPvtData); err != nil {
		return err
	}

	
	lc.postCommit(blockAndPvtData.Block)

	return nil
}


func (lc *LedgerCommitter) GetPvtDataAndBlockByNum(seqNum uint64) (*ledger.BlockAndPvtData, error) {
	return lc.PeerLedgerSupport.GetPvtDataAndBlockByNum(seqNum, nil)
}


func (lc *LedgerCommitter) postCommit(block *common.Block) {
	
	bevent, fbevent, channelID, err := producer.CreateBlockEvents(block)
	if err != nil {
		logger.Errorf("Channel [%s] Error processing block events for block number [%d]: %+v", channelID, block.Header.Number, err)
	} else {
		if err := producer.Send(bevent); err != nil {
			logger.Errorf("Channel [%s] Error sending block event for block number [%d]: %+v", channelID, block.Header.Number, err)
		}
		if err := producer.Send(fbevent); err != nil {
			logger.Errorf("Channel [%s] Error sending filtered block event for block number [%d]: %+v", channelID, block.Header.Number, err)
		}
	}
}


func (lc *LedgerCommitter) LedgerHeight() (uint64, error) {
	var info *common.BlockchainInfo
	var err error
	if info, err = lc.GetBlockchainInfo(); err != nil {
		logger.Errorf("Cannot get blockchain info, %s", info)
		return uint64(0), err
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
