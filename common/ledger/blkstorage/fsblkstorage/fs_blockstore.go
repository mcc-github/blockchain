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

package fsblkstorage

import (
	"time"

	"github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)


type fsBlockStore struct {
	id      string
	conf    *Conf
	fileMgr *blockfileMgr
	stats   *ledgerStats
}


func newFsBlockStore(id string, conf *Conf, indexConfig *blkstorage.IndexConfig,
	dbHandle *leveldbhelper.DBHandle, stats *stats) *fsBlockStore {
	fileMgr := newBlockfileMgr(id, conf, indexConfig, dbHandle)

	
	ledgerStats := stats.ledgerStats(id)
	info := fileMgr.getBlockchainInfo()
	ledgerStats.updateBlockchainHeight(info.Height)

	return &fsBlockStore{id, conf, fileMgr, ledgerStats}
}


func (store *fsBlockStore) AddBlock(block *common.Block) error {
	
	startBlockCommit := time.Now()
	result := store.fileMgr.addBlock(block)
	elapsedBlockCommit := time.Since(startBlockCommit)

	store.updateBlockStats(block.Header.Number, elapsedBlockCommit)

	return result
}


func (store *fsBlockStore) GetBlockchainInfo() (*common.BlockchainInfo, error) {
	return store.fileMgr.getBlockchainInfo(), nil
}


func (store *fsBlockStore) RetrieveBlocks(startNum uint64) (ledger.ResultsIterator, error) {
	return store.fileMgr.retrieveBlocks(startNum)
}


func (store *fsBlockStore) RetrieveBlockByHash(blockHash []byte) (*common.Block, error) {
	return store.fileMgr.retrieveBlockByHash(blockHash)
}


func (store *fsBlockStore) RetrieveBlockByNumber(blockNum uint64) (*common.Block, error) {
	return store.fileMgr.retrieveBlockByNumber(blockNum)
}


func (store *fsBlockStore) RetrieveTxByID(txID string) (*common.Envelope, error) {
	return store.fileMgr.retrieveTransactionByID(txID)
}


func (store *fsBlockStore) RetrieveTxByBlockNumTranNum(blockNum uint64, tranNum uint64) (*common.Envelope, error) {
	return store.fileMgr.retrieveTransactionByBlockNumTranNum(blockNum, tranNum)
}

func (store *fsBlockStore) RetrieveBlockByTxID(txID string) (*common.Block, error) {
	return store.fileMgr.retrieveBlockByTxID(txID)
}

func (store *fsBlockStore) RetrieveTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	return store.fileMgr.retrieveTxValidationCodeByTxID(txID)
}


func (store *fsBlockStore) Shutdown() {
	logger.Debugf("closing fs blockStore:%s", store.id)
	store.fileMgr.close()
}

func (store *fsBlockStore) updateBlockStats(blockNum uint64, blockstorageCommitTime time.Duration) {
	store.stats.updateBlockchainHeight(blockNum + 1)
	store.stats.updateBlockstorageCommitTime(blockstorageCommitTime)
}
