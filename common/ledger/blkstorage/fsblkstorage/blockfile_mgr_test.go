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
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	ledgerutil "github.com/mcc-github/blockchain/core/ledger/util"

	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	putil "github.com/mcc-github/blockchain/protos/utils"
)

func TestBlockfileMgrBlockReadWrite(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	defer blkfileMgrWrapper.close()
	blocks := testutil.ConstructTestBlocks(t, 10)
	blkfileMgrWrapper.addBlocks(blocks)
	blkfileMgrWrapper.testGetBlockByHash(blocks)
	blkfileMgrWrapper.testGetBlockByNumber(blocks, 0)
}

func TestBlockfileMgrCrashDuringWriting(t *testing.T) {
	testBlockfileMgrCrashDuringWriting(t, 10, 2, 1000, 10)
	testBlockfileMgrCrashDuringWriting(t, 10, 2, 1000, 1)
	testBlockfileMgrCrashDuringWriting(t, 10, 2, 1000, 0)
	testBlockfileMgrCrashDuringWriting(t, 0, 0, 1000, 10)
	testBlockfileMgrCrashDuringWriting(t, 0, 5, 1000, 10)
}

func testBlockfileMgrCrashDuringWriting(t *testing.T, numBlocksBeforeCheckpoint int,
	numBlocksAfterCheckpoint int, numLastBlockBytes int, numPartialBytesToWrite int) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	ledgerid := "testLedger"
	blkfileMgrWrapper := newTestBlockfileWrapper(env, ledgerid)
	bg, gb := testutil.NewBlockGenerator(t, ledgerid, false)

	
	totalBlocks := numBlocksBeforeCheckpoint + numBlocksAfterCheckpoint
	allBlocks := []*common.Block{gb}
	allBlocks = append(allBlocks, bg.NextTestBlocks(totalBlocks+1)...)

	
	blocksBeforeCP := []*common.Block{}
	blocksAfterCP := []*common.Block{}
	if numBlocksBeforeCheckpoint != 0 {
		blocksBeforeCP = allBlocks[0:numBlocksBeforeCheckpoint]
	}
	if numBlocksAfterCheckpoint != 0 {
		blocksAfterCP = allBlocks[numBlocksBeforeCheckpoint : numBlocksBeforeCheckpoint+numBlocksAfterCheckpoint]
	}
	blocksAfterRestart := allBlocks[numBlocksBeforeCheckpoint+numBlocksAfterCheckpoint:]

	
	blkfileMgrWrapper.addBlocks(blocksBeforeCP)
	currentCPInfo := blkfileMgrWrapper.blockfileMgr.cpInfo
	cpInfo1 := &checkpointInfo{
		currentCPInfo.latestFileChunkSuffixNum,
		currentCPInfo.latestFileChunksize,
		currentCPInfo.isChainEmpty,
		currentCPInfo.lastBlockNumber}

	
	blkfileMgrWrapper.addBlocks(blocksAfterCP)
	cpInfo2 := blkfileMgrWrapper.blockfileMgr.cpInfo

	
	lastBlockBytes := []byte{}
	encodedLen := proto.EncodeVarint(uint64(numLastBlockBytes))
	randomBytes := testutil.ConstructRandomBytes(t, numLastBlockBytes)
	lastBlockBytes = append(lastBlockBytes, encodedLen...)
	lastBlockBytes = append(lastBlockBytes, randomBytes...)
	partialBytes := lastBlockBytes[:numPartialBytesToWrite]
	blkfileMgrWrapper.blockfileMgr.currentFileWriter.append(partialBytes, true)
	blkfileMgrWrapper.blockfileMgr.saveCurrentInfo(cpInfo1, true)
	blkfileMgrWrapper.close()

	
	blkfileMgrWrapper = newTestBlockfileWrapper(env, ledgerid)
	defer blkfileMgrWrapper.close()
	cpInfo3 := blkfileMgrWrapper.blockfileMgr.cpInfo
	testutil.AssertEquals(t, cpInfo3, cpInfo2)

	
	blkfileMgrWrapper.addBlocks(blocksAfterRestart)
	testBlockfileMgrBlockIterator(t, blkfileMgrWrapper.blockfileMgr, 0, len(allBlocks)-1, allBlocks)
}

func TestBlockfileMgrBlockIterator(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	defer blkfileMgrWrapper.close()
	blocks := testutil.ConstructTestBlocks(t, 10)
	blkfileMgrWrapper.addBlocks(blocks)
	testBlockfileMgrBlockIterator(t, blkfileMgrWrapper.blockfileMgr, 0, 7, blocks[0:8])
}

func testBlockfileMgrBlockIterator(t *testing.T, blockfileMgr *blockfileMgr,
	firstBlockNum int, lastBlockNum int, expectedBlocks []*common.Block) {
	itr, err := blockfileMgr.retrieveBlocks(uint64(firstBlockNum))
	defer itr.Close()
	testutil.AssertNoError(t, err, "Error while getting blocks iterator")
	numBlocksItrated := 0
	for {
		block, err := itr.Next()
		testutil.AssertNoError(t, err, fmt.Sprintf("Error while getting block number [%d] from iterator", numBlocksItrated))
		testutil.AssertEquals(t, block, expectedBlocks[numBlocksItrated])
		numBlocksItrated++
		if numBlocksItrated == lastBlockNum-firstBlockNum+1 {
			break
		}
	}
	testutil.AssertEquals(t, numBlocksItrated, lastBlockNum-firstBlockNum+1)
}

func TestBlockfileMgrBlockchainInfo(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	defer blkfileMgrWrapper.close()

	bcInfo := blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	testutil.AssertEquals(t, bcInfo, &common.BlockchainInfo{Height: 0, CurrentBlockHash: nil, PreviousBlockHash: nil})

	blocks := testutil.ConstructTestBlocks(t, 10)
	blkfileMgrWrapper.addBlocks(blocks)
	bcInfo = blkfileMgrWrapper.blockfileMgr.getBlockchainInfo()
	testutil.AssertEquals(t, bcInfo.Height, uint64(10))
}

func TestBlockfileMgrGetTxById(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	defer blkfileMgrWrapper.close()
	blocks := testutil.ConstructTestBlocks(t, 2)
	blkfileMgrWrapper.addBlocks(blocks)
	for _, blk := range blocks {
		for j, txEnvelopeBytes := range blk.Data.Data {
			
			txID, err := extractTxID(blk.Data.Data[j])
			testutil.AssertNoError(t, err, "")
			txEnvelopeFromFileMgr, err := blkfileMgrWrapper.blockfileMgr.retrieveTransactionByID(txID)
			testutil.AssertNoError(t, err, "Error while retrieving tx from blkfileMgr")
			txEnvelope, err := putil.GetEnvelopeFromBlock(txEnvelopeBytes)
			testutil.AssertNoError(t, err, "Error while unmarshalling tx")
			testutil.AssertEquals(t, txEnvelopeFromFileMgr, txEnvelope)
		}
	}
}



func TestBlockfileMgrGetTxByIdDuplicateTxid(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkStore, err := env.provider.OpenBlockStore("testLedger")
	testutil.AssertNoError(env.t, err, "")
	blkFileMgr := blkStore.(*fsBlockStore).fileMgr
	bg, gb := testutil.NewBlockGenerator(t, "testLedger", false)
	testutil.AssertNoError(t, blkFileMgr.addBlock(gb), "")

	block1 := bg.NextBlockWithTxid(
		[][]byte{
			[]byte("tx with id=txid-1"),
			[]byte("tx with id=txid-2"),
			[]byte("another tx with existing id=txid-1"),
		},
		[]string{"txid-1", "txid-2", "txid-1"},
	)
	txValidationFlags := ledgerutil.NewTxValidationFlags(3)
	txValidationFlags.SetFlag(0, peer.TxValidationCode_VALID)
	txValidationFlags.SetFlag(1, peer.TxValidationCode_INVALID_OTHER_REASON)
	txValidationFlags.SetFlag(2, peer.TxValidationCode_DUPLICATE_TXID)
	block1.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txValidationFlags
	testutil.AssertNoError(t, blkFileMgr.addBlock(block1), "")

	block2 := bg.NextBlockWithTxid(
		[][]byte{
			[]byte("tx with id=txid-3"),
			[]byte("yet another tx with existing id=txid-1"),
		},
		[]string{"txid-3", "txid-1"},
	)
	txValidationFlags = ledgerutil.NewTxValidationFlags(2)
	txValidationFlags.SetFlag(0, peer.TxValidationCode_VALID)
	txValidationFlags.SetFlag(1, peer.TxValidationCode_DUPLICATE_TXID)
	block2.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txValidationFlags
	testutil.AssertNoError(t, blkFileMgr.addBlock(block2), "")

	txenvp1, err := putil.GetEnvelopeFromBlock(block1.Data.Data[0])
	testutil.AssertNoError(t, err, "")
	txenvp2, err := putil.GetEnvelopeFromBlock(block1.Data.Data[1])
	testutil.AssertNoError(t, err, "")
	txenvp3, err := putil.GetEnvelopeFromBlock(block2.Data.Data[0])
	testutil.AssertNoError(t, err, "")

	indexedTxenvp, _ := blkFileMgr.retrieveTransactionByID("txid-1")
	testutil.AssertEquals(t, indexedTxenvp, txenvp1)
	indexedTxenvp, _ = blkFileMgr.retrieveTransactionByID("txid-2")
	testutil.AssertEquals(t, indexedTxenvp, txenvp2)
	indexedTxenvp, _ = blkFileMgr.retrieveTransactionByID("txid-3")
	testutil.AssertEquals(t, indexedTxenvp, txenvp3)

	blk, _ := blkFileMgr.retrieveBlockByTxID("txid-1")
	testutil.AssertEquals(t, blk, block1)
	blk, _ = blkFileMgr.retrieveBlockByTxID("txid-2")
	testutil.AssertEquals(t, blk, block1)
	blk, _ = blkFileMgr.retrieveBlockByTxID("txid-3")
	testutil.AssertEquals(t, blk, block2)

	validationCode, _ := blkFileMgr.retrieveTxValidationCodeByTxID("txid-1")
	testutil.AssertEquals(t, validationCode, peer.TxValidationCode_VALID)
	validationCode, _ = blkFileMgr.retrieveTxValidationCodeByTxID("txid-2")
	testutil.AssertEquals(t, validationCode, peer.TxValidationCode_INVALID_OTHER_REASON)
	validationCode, _ = blkFileMgr.retrieveTxValidationCodeByTxID("txid-3")
	testutil.AssertEquals(t, validationCode, peer.TxValidationCode_VALID)
}

func TestBlockfileMgrGetTxByBlockNumTranNum(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	defer blkfileMgrWrapper.close()
	blocks := testutil.ConstructTestBlocks(t, 10)
	blkfileMgrWrapper.addBlocks(blocks)
	for blockIndex, blk := range blocks {
		for tranIndex, txEnvelopeBytes := range blk.Data.Data {
			
			txEnvelopeFromFileMgr, err := blkfileMgrWrapper.blockfileMgr.retrieveTransactionByBlockNumTranNum(uint64(blockIndex), uint64(tranIndex))
			testutil.AssertNoError(t, err, "Error while retrieving tx from blkfileMgr")
			txEnvelope, err := putil.GetEnvelopeFromBlock(txEnvelopeBytes)
			testutil.AssertNoError(t, err, "Error while unmarshalling tx")
			testutil.AssertEquals(t, txEnvelopeFromFileMgr, txEnvelope)
		}
	}
}

func TestBlockfileMgrRestart(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	ledgerid := "testLedger"
	blkfileMgrWrapper := newTestBlockfileWrapper(env, ledgerid)
	blocks := testutil.ConstructTestBlocks(t, 10)
	blkfileMgrWrapper.addBlocks(blocks)
	expectedHeight := uint64(10)
	testutil.AssertEquals(t, blkfileMgrWrapper.blockfileMgr.getBlockchainInfo().Height, expectedHeight)
	blkfileMgrWrapper.close()

	blkfileMgrWrapper = newTestBlockfileWrapper(env, ledgerid)
	defer blkfileMgrWrapper.close()
	testutil.AssertEquals(t, int(blkfileMgrWrapper.blockfileMgr.cpInfo.lastBlockNumber), 9)
	blkfileMgrWrapper.testGetBlockByHash(blocks)
	testutil.AssertEquals(t, blkfileMgrWrapper.blockfileMgr.getBlockchainInfo().Height, expectedHeight)
}

func TestBlockfileMgrFileRolling(t *testing.T) {
	blocks := testutil.ConstructTestBlocks(t, 200)
	size := 0
	for _, block := range blocks[:100] {
		by, _, err := serializeBlock(block)
		testutil.AssertNoError(t, err, "Error while serializing block")
		blockBytesSize := len(by)
		encodedLen := proto.EncodeVarint(uint64(blockBytesSize))
		size += blockBytesSize + len(encodedLen)
	}

	maxFileSie := int(0.75 * float64(size))
	env := newTestEnv(t, NewConf(testPath(), maxFileSie))
	defer env.Cleanup()
	ledgerid := "testLedger"
	blkfileMgrWrapper := newTestBlockfileWrapper(env, ledgerid)
	blkfileMgrWrapper.addBlocks(blocks[:100])
	testutil.AssertEquals(t, blkfileMgrWrapper.blockfileMgr.cpInfo.latestFileChunkSuffixNum, 1)
	blkfileMgrWrapper.testGetBlockByHash(blocks[:100])
	blkfileMgrWrapper.close()

	blkfileMgrWrapper = newTestBlockfileWrapper(env, ledgerid)
	defer blkfileMgrWrapper.close()
	blkfileMgrWrapper.addBlocks(blocks[100:])
	testutil.AssertEquals(t, blkfileMgrWrapper.blockfileMgr.cpInfo.latestFileChunkSuffixNum, 2)
	blkfileMgrWrapper.testGetBlockByHash(blocks[100:])
}

func TestBlockfileMgrGetBlockByTxID(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	defer blkfileMgrWrapper.close()
	blocks := testutil.ConstructTestBlocks(t, 10)
	blkfileMgrWrapper.addBlocks(blocks)
	for _, blk := range blocks {
		for j := range blk.Data.Data {
			
			txID, err := extractTxID(blk.Data.Data[j])
			testutil.AssertNoError(t, err, "")

			blockFromFileMgr, err := blkfileMgrWrapper.blockfileMgr.retrieveBlockByTxID(txID)
			testutil.AssertNoError(t, err, "Error while retrieving block from blkfileMgr")
			testutil.AssertEquals(t, blockFromFileMgr, blk)
		}
	}
}
