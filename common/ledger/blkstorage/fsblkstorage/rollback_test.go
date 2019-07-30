/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package fsblkstorage

import (
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRollback(t *testing.T) {
	path := testPath()
	blocks := testutil.ConstructTestBlocks(t, 1000)
	
	
	
	maxFileSize := int(0.2 * float64(testutilEstimateTotalSizeOnDisk(t, blocks)))
	
	env := newTestEnv(t, NewConf(path, maxFileSize))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")

	
	blkfileMgrWrapper.addBlocks(blocks)

	
	expectedBlockchainInfo := &common.BlockchainInfo{
		Height:            1000,
		CurrentBlockHash:  protoutil.BlockHeaderHash(blocks[999].Header),
		PreviousBlockHash: protoutil.BlockHeaderHash(blocks[998].Header),
	}
	actualBlockchainInfo := blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	assert.Equal(t, expectedBlockchainInfo, actualBlockchainInfo)

	
	expectedCheckpointInfoLastBlockNumber := uint64(999)
	expectedCheckpointInfoIsChainEmpty := false
	actualCheckpointInfo, err := blkfileMgrWrapper.blockfileMgr.loadCurrentInfo()
	assert.NoError(t, err)
	assert.Equal(t, expectedCheckpointInfoLastBlockNumber, actualCheckpointInfo.lastBlockNumber)
	assert.Equal(t, expectedCheckpointInfoIsChainEmpty, actualCheckpointInfo.isChainEmpty)
	assert.True(t, actualCheckpointInfo.latestFileChunkSuffixNum >= 5)
	
	
	

	
	blkfileMgrWrapper.testGetBlockByNumber(blocks, 0, nil)
	blkfileMgrWrapper.testGetBlockByHash(blocks, nil)
	blkfileMgrWrapper.testGetBlockByTxID(blocks, nil)

	
	env.provider.Close()
	blkfileMgrWrapper.close()

	
	ledgerDir := (&Conf{blockStorageDir: path}).getLedgerBlockDir("testLedger")
	_, _, numBlocksInLastFile, err := scanForLastCompleteBlock(ledgerDir, actualCheckpointInfo.latestFileChunkSuffixNum, 0)
	assert.NoError(t, err)
	lastFileSuffixNum := actualCheckpointInfo.latestFileChunkSuffixNum
	lastBlockNumberInLastFile := uint64(999)
	middleBlockNumberInLastFile := uint64(999 - (numBlocksInLastFile / 2))
	firstBlockNumberInLastFile := uint64(999 - numBlocksInLastFile + 1)

	
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	err = Rollback(path, "testLedger", lastBlockNumberInLastFile-uint64(1), indexConfig)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, lastBlockNumberInLastFile-uint64(1), lastFileSuffixNum, indexConfig)

	
	err = Rollback(path, "testLedger", middleBlockNumberInLastFile, indexConfig)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, middleBlockNumberInLastFile, lastFileSuffixNum, indexConfig)

	
	err = Rollback(path, "testLedger", firstBlockNumberInLastFile, indexConfig)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, firstBlockNumberInLastFile, lastFileSuffixNum, indexConfig)

	
	err = Rollback(path, "testLedger", firstBlockNumberInLastFile-1, indexConfig)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, firstBlockNumberInLastFile-1, lastFileSuffixNum-1, indexConfig)

	
	blockBytes, _, numBlocks, err := scanForLastCompleteBlock(ledgerDir, lastFileSuffixNum/2, 0)
	assert.NoError(t, err)
	blockInfo, err := extractSerializedBlockInfo(blockBytes)
	assert.NoError(t, err)
	middleBlockNumberInMiddleFile := blockInfo.blockHeader.Number - uint64(numBlocks/2)

	
	err = Rollback(path, "testLedger", middleBlockNumberInMiddleFile, indexConfig)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, middleBlockNumberInMiddleFile, lastFileSuffixNum/2, indexConfig)

	
	err = Rollback(path, "testLedger", 5, indexConfig)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, 5, 0, indexConfig)

	
	err = Rollback(path, "testLedger", 1, indexConfig)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, 1, 0, indexConfig)
}



func TestRollbackWithOnlyBlockIndexAttributes(t *testing.T) {
	path := testPath()
	blocks := testutil.ConstructTestBlocks(t, 100)
	
	
	
	maxFileSize := int(0.2 * float64(testutilEstimateTotalSizeOnDisk(t, blocks)))
	

	onlyBlockNumIndex := []blkstorage.IndexableAttr{
		blkstorage.IndexableAttrBlockNum,
	}

	env := newTestEnvSelectiveIndexing(t, NewConf(path, maxFileSize), onlyBlockNumIndex, &disabled.Provider{})
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")

	
	blkfileMgrWrapper.addBlocks(blocks)

	
	expectedBlockchainInfo := &common.BlockchainInfo{
		Height:            100,
		CurrentBlockHash:  protoutil.BlockHeaderHash(blocks[99].Header),
		PreviousBlockHash: protoutil.BlockHeaderHash(blocks[98].Header),
	}
	actualBlockchainInfo := blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	assert.Equal(t, expectedBlockchainInfo, actualBlockchainInfo)

	
	expectedCheckpointInfoLastBlockNumber := uint64(99)
	expectedCheckpointInfoIsChainEmpty := false
	actualCheckpointInfo, err := blkfileMgrWrapper.blockfileMgr.loadCurrentInfo()
	assert.NoError(t, err)
	assert.Equal(t, expectedCheckpointInfoLastBlockNumber, actualCheckpointInfo.lastBlockNumber)
	assert.Equal(t, expectedCheckpointInfoIsChainEmpty, actualCheckpointInfo.isChainEmpty)
	assert.True(t, actualCheckpointInfo.latestFileChunkSuffixNum >= 5)
	
	
	

	
	env.provider.Close()
	blkfileMgrWrapper.close()

	
	onlyBlockNumIndexCfg := &blkstorage.IndexConfig{
		AttrsToIndex: onlyBlockNumIndex,
	}
	err = Rollback(path, "testLedger", 2, onlyBlockNumIndexCfg)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, 2, 0, onlyBlockNumIndexCfg)
}

func TestRollbackWithNoIndexDir(t *testing.T) {
	path := testPath()
	blocks := testutil.ConstructTestBlocks(t, 100)
	
	
	
	maxFileSize := int(0.2 * float64(testutilEstimateTotalSizeOnDisk(t, blocks)))
	
	conf := NewConf(path, maxFileSize)
	env := newTestEnv(t, conf)
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")

	
	blkfileMgrWrapper.addBlocks(blocks)

	
	expectedBlockchainInfo := &common.BlockchainInfo{
		Height:            100,
		CurrentBlockHash:  protoutil.BlockHeaderHash(blocks[99].Header),
		PreviousBlockHash: protoutil.BlockHeaderHash(blocks[98].Header),
	}
	actualBlockchainInfo := blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	assert.Equal(t, expectedBlockchainInfo, actualBlockchainInfo)

	
	expectedCheckpointInfoLastBlockNumber := uint64(99)
	expectedCheckpointInfoIsChainEmpty := false
	actualCheckpointInfo, err := blkfileMgrWrapper.blockfileMgr.loadCurrentInfo()
	assert.NoError(t, err)
	assert.Equal(t, expectedCheckpointInfoLastBlockNumber, actualCheckpointInfo.lastBlockNumber)
	assert.Equal(t, expectedCheckpointInfoIsChainEmpty, actualCheckpointInfo.isChainEmpty)
	assert.True(t, actualCheckpointInfo.latestFileChunkSuffixNum >= 5)
	
	
	

	
	env.provider.Close()
	blkfileMgrWrapper.close()

	
	indexDir := conf.getIndexDir()
	err = os.RemoveAll(indexDir)
	assert.NoError(t, err)

	
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	err = Rollback(path, "testLedger", 2, indexConfig)
	assert.NoError(t, err)
	assertBlockStoreRollback(t, path, "testLedger", maxFileSize, blocks, 2, 0, indexConfig)
}

func TestValidateRollbackParams(t *testing.T) {
	path := testPath()
	env := newTestEnv(t, NewConf(path, 1024*24))
	defer env.Cleanup()

	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")

	
	blocks := testutil.ConstructTestBlocks(t, 10)
	blkfileMgrWrapper.addBlocks(blocks)

	
	err := ValidateRollbackParams(path, "testLedger", 5)
	assert.NoError(t, err)

	
	err = ValidateRollbackParams(path, "noLedger", 5)
	assert.Equal(t, "ledgerID [noLedger] does not exist", err.Error())

	err = ValidateRollbackParams(path, "testLedger", 15)
	assert.Equal(t, "target block number [15] should be less than the biggest block number [9]", err.Error())
}

func TestDuplicateTxIDDuringRollback(t *testing.T) {
	path := testPath()
	blocks := testutil.ConstructTestBlocks(t, 4)
	maxFileSize := 1024 * 1024 * 4
	env := newTestEnv(t, NewConf(path, maxFileSize))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	blocks[3].Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER][0] = byte(peer.TxValidationCode_DUPLICATE_TXID)
	testutil.SetTxID(t, blocks[3], 0, "tx0")
	testutil.SetTxID(t, blocks[2], 0, "tx0")

	
	blkfileMgrWrapper.addBlocks(blocks)

	
	expectedBlockchainInfo := &common.BlockchainInfo{
		Height:            4,
		CurrentBlockHash:  protoutil.BlockHeaderHash(blocks[3].Header),
		PreviousBlockHash: protoutil.BlockHeaderHash(blocks[2].Header),
	}
	actualBlockchainInfo := blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	assert.Equal(t, expectedBlockchainInfo, actualBlockchainInfo)

	
	blkfileMgrWrapper.testGetTransactionByTxID("tx0", blocks[2].Data.Data[0], nil)

	
	env.provider.Close()
	blkfileMgrWrapper.close()

	
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	err := Rollback(path, "testLedger", 2, indexConfig)
	assert.NoError(t, err)

	env = newTestEnv(t, NewConf(path, maxFileSize))
	blkfileMgrWrapper = newTestBlockfileWrapper(env, "testLedger")

	
	expectedBlockchainInfo = &common.BlockchainInfo{
		Height:            3,
		CurrentBlockHash:  protoutil.BlockHeaderHash(blocks[2].Header),
		PreviousBlockHash: protoutil.BlockHeaderHash(blocks[1].Header),
	}
	actualBlockchainInfo = blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	assert.Equal(t, expectedBlockchainInfo, actualBlockchainInfo)

	
	blkfileMgrWrapper.testGetTransactionByTxID("tx0", blocks[2].Data.Data[0], nil)
}

func TestRollbackTxIDMissingFromIndex(t *testing.T) {
	
	path := testPath()
	blocks := testutil.ConstructTestBlocks(t, 4)
	testutil.SetTxID(t, blocks[3], 0, "blk3_tx0")

	env := newTestEnv(t, NewConf(path, 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")

	blkfileMgrWrapper.addBlocks(blocks)
	expectedBlockchainInfo := &common.BlockchainInfo{
		Height:            4,
		CurrentBlockHash:  protoutil.BlockHeaderHash(blocks[3].Header),
		PreviousBlockHash: protoutil.BlockHeaderHash(blocks[2].Header),
	}
	actualBlockchainInfo := blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	assert.Equal(t, expectedBlockchainInfo, actualBlockchainInfo)
	blkfileMgrWrapper.testGetTransactionByTxID("blk3_tx0", blocks[3].Data.Data[0], nil)

	
	testutilDeleteTxFromIndex(t, blkfileMgrWrapper.blockfileMgr.db, 3, 0, "blk3_tx0")
	env.provider.Close()
	blkfileMgrWrapper.close()

	
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	err := Rollback(path, "testLedger", 2, indexConfig)
	assert.NoError(t, err)

	
	env = newTestEnv(t, NewConf(path, 0))
	blkfileMgrWrapper = newTestBlockfileWrapper(env, "testLedger")
	expectedBlockchainInfo = &common.BlockchainInfo{
		Height:            3,
		CurrentBlockHash:  protoutil.BlockHeaderHash(blocks[2].Header),
		PreviousBlockHash: protoutil.BlockHeaderHash(blocks[1].Header),
	}
	actualBlockchainInfo = blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	assert.Equal(t, expectedBlockchainInfo, actualBlockchainInfo)
}

func assertBlockStoreRollback(t *testing.T, path, ledgerID string, maxFileSize int, blocks []*common.Block,
	rollbackedToBlkNum uint64, lastFileSuffixNum int, indexConfig *blkstorage.IndexConfig) {

	env := newTestEnvSelectiveIndexing(t, NewConf(path, maxFileSize), indexConfig.AttrsToIndex, &disabled.Provider{})
	blkfileMgrWrapper := newTestBlockfileWrapper(env, ledgerID)

	
	expectedBlockchainInfo := &common.BlockchainInfo{
		Height:            rollbackedToBlkNum + 1,
		CurrentBlockHash:  protoutil.BlockHeaderHash(blocks[rollbackedToBlkNum].Header),
		PreviousBlockHash: protoutil.BlockHeaderHash(blocks[rollbackedToBlkNum-1].Header),
	}
	actualBlockchainInfo := blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	assert.Equal(t, expectedBlockchainInfo, actualBlockchainInfo)

	
	expectedCheckpointInfoLastBlockNumber := rollbackedToBlkNum
	expectedCheckpointInfoIsChainEmpty := false
	expectedBlockchainInfoLastFileSuffixNum := lastFileSuffixNum
	actualCheckpointInfo, err := blkfileMgrWrapper.blockfileMgr.loadCurrentInfo()
	assert.NoError(t, err)
	assert.Equal(t, expectedCheckpointInfoLastBlockNumber, actualCheckpointInfo.lastBlockNumber)
	assert.Equal(t, expectedCheckpointInfoIsChainEmpty, actualCheckpointInfo.isChainEmpty)
	assert.Equal(t, expectedBlockchainInfoLastFileSuffixNum, actualCheckpointInfo.latestFileChunkSuffixNum)

	
	if blkfileMgrWrapper.blockfileMgr.index.isAttributeIndexed(blkstorage.IndexableAttrBlockNum) {
		blkfileMgrWrapper.testGetBlockByNumber(blocks[:rollbackedToBlkNum+1], 0, nil)
	}
	if blkfileMgrWrapper.blockfileMgr.index.isAttributeIndexed(blkstorage.IndexableAttrBlockHash) {
		blkfileMgrWrapper.testGetBlockByHash(blocks[:rollbackedToBlkNum+1], nil)
	}
	if blkfileMgrWrapper.blockfileMgr.index.isAttributeIndexed(blkstorage.IndexableAttrTxID) {
		blkfileMgrWrapper.testGetBlockByTxID(blocks[:rollbackedToBlkNum+1], nil)
	}

	
	
	expectedErr := errors.New("Entry not found in index")
	if blkfileMgrWrapper.blockfileMgr.index.isAttributeIndexed(blkstorage.IndexableAttrBlockHash) {
		blkfileMgrWrapper.testGetBlockByHash(blocks[rollbackedToBlkNum+1:], expectedErr)
	}
	if blkfileMgrWrapper.blockfileMgr.index.isAttributeIndexed(blkstorage.IndexableAttrTxID) {
		blkfileMgrWrapper.testGetBlockByTxID(blocks[rollbackedToBlkNum+1:], expectedErr)
	}

	
	env.provider.Close()
	blkfileMgrWrapper.close()
}

func testutilDeleteTxFromIndex(t *testing.T, db *leveldbhelper.DBHandle, blkNum, txNum uint64, txID string) {
	indexKeys := [][]byte{
		constructTxIDKey(txID),
		constructBlockNumTranNumKey(blkNum, txNum),
		constructBlockTxIDKey(txID),
		constructTxValidationCodeIDKey(txID),
	}

	batch := leveldbhelper.NewUpdateBatch()
	for _, key := range indexKeys {
		value, err := db.Get(key)
		assert.NoError(t, err)
		assert.NotNil(t, value)
		batch.Delete(key)
	}
	db.WriteBatch(batch, true)

	for _, key := range indexKeys {
		value, err := db.Get(key)
		assert.NoError(t, err)
		assert.Nil(t, value)
	}
}
