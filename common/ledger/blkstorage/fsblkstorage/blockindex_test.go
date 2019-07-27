/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fsblkstorage

import (
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

type noopIndex struct {
}

func (i *noopIndex) getLastBlockIndexed() (uint64, error) {
	return 0, nil
}
func (i *noopIndex) indexBlock(blockIdxInfo *blockIdxInfo) error {
	return nil
}
func (i *noopIndex) getBlockLocByHash(blockHash []byte) (*fileLocPointer, error) {
	return nil, nil
}
func (i *noopIndex) getBlockLocByBlockNum(blockNum uint64) (*fileLocPointer, error) {
	return nil, nil
}
func (i *noopIndex) getTxLoc(txID string) (*fileLocPointer, error) {
	return nil, nil
}
func (i *noopIndex) getTXLocByBlockNumTranNum(blockNum uint64, tranNum uint64) (*fileLocPointer, error) {
	return nil, nil
}

func (i *noopIndex) getBlockLocByTxID(txID string) (*fileLocPointer, error) {
	return nil, nil
}

func (i *noopIndex) getTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	return peer.TxValidationCode(-1), nil
}

func TestBlockIndexSync(t *testing.T) {
	testBlockIndexSync(t, 10, 5, false)
	testBlockIndexSync(t, 10, 5, true)
	testBlockIndexSync(t, 10, 0, true)
	testBlockIndexSync(t, 10, 10, true)
}

func testBlockIndexSync(t *testing.T, numBlocks int, numBlocksToIndex int, syncByRestart bool) {
	testName := fmt.Sprintf("%v/%v/%v", numBlocks, numBlocksToIndex, syncByRestart)
	t.Run(testName, func(t *testing.T) {
		env := newTestEnv(t, NewConf(testPath(), 0))
		defer env.Cleanup()
		ledgerid := "testledger"
		blkfileMgrWrapper := newTestBlockfileWrapper(env, ledgerid)
		defer blkfileMgrWrapper.close()
		blkfileMgr := blkfileMgrWrapper.blockfileMgr
		origIndex := blkfileMgr.index
		
		blocks := testutil.ConstructTestBlocks(t, numBlocks)
		
		blkfileMgrWrapper.addBlocks(blocks[:numBlocksToIndex])

		
		blkfileMgr.index = &noopIndex{}
		blkfileMgrWrapper.addBlocks(blocks[numBlocksToIndex:])

		
		blkfileMgr.index = origIndex
		
		for i := 0; i < numBlocksToIndex; i++ {
			block, err := blkfileMgr.retrieveBlockByNumber(uint64(i))
			assert.NoError(t, err, "block [%d] should have been present in the index", i)
			assert.Equal(t, blocks[i], block)
		}

		
		for i := numBlocksToIndex + 1; i <= numBlocks; i++ {
			_, err := blkfileMgr.retrieveBlockByNumber(uint64(i))
			assert.Exactly(t, blkstorage.ErrNotFoundInIndex, err)
		}

		
		if syncByRestart {
			blkfileMgrWrapper.close()
			blkfileMgrWrapper = newTestBlockfileWrapper(env, ledgerid)
			defer blkfileMgrWrapper.close()
			blkfileMgr = blkfileMgrWrapper.blockfileMgr
		} else {
			blkfileMgr.syncIndex()
		}

		
		for i := numBlocksToIndex; i < numBlocks; i++ {
			block, err := blkfileMgr.retrieveBlockByNumber(uint64(i))
			assert.NoError(t, err, "block [%d] should have been present in the index", i)
			assert.Equal(t, blocks[i], block)
		}
	})
}

func TestBlockIndexSelectiveIndexingWrongConfig(t *testing.T) {
	testBlockIndexSelectiveIndexingWrongConfig(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockTxID})
	testBlockIndexSelectiveIndexingWrongConfig(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrTxValidationCode})
	testBlockIndexSelectiveIndexingWrongConfig(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockTxID, blkstorage.IndexableAttrBlockNum})
	testBlockIndexSelectiveIndexingWrongConfig(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrTxValidationCode, blkstorage.IndexableAttrBlockNumTranNum})

}

func testBlockIndexSelectiveIndexingWrongConfig(t *testing.T, indexItems []blkstorage.IndexableAttr) {
	var testName string
	for _, s := range indexItems {
		testName = testName + string(s)
	}
	t.Run(testName, func(t *testing.T) {
		env := newTestEnvSelectiveIndexing(t, NewConf(testPath(), 0), indexItems, &disabled.Provider{})
		defer env.Cleanup()

		assert.Panics(t, func() {
			env.provider.OpenBlockStore("test-ledger")
		}, "A panic is expected when index is opened with wrong configs")
	})
}

func TestBlockIndexSelectiveIndexing(t *testing.T) {
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{})
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockHash})
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockNum})
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrTxID})
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockNumTranNum})
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockHash, blkstorage.IndexableAttrBlockNum})
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrTxID, blkstorage.IndexableAttrBlockNumTranNum})
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrTxID, blkstorage.IndexableAttrBlockTxID})
	testBlockIndexSelectiveIndexing(t, []blkstorage.IndexableAttr{blkstorage.IndexableAttrTxID, blkstorage.IndexableAttrTxValidationCode})
}

func testBlockIndexSelectiveIndexing(t *testing.T, indexItems []blkstorage.IndexableAttr) {
	var testName string
	for _, s := range indexItems {
		testName = testName + string(s)
	}
	t.Run(testName, func(t *testing.T) {
		env := newTestEnvSelectiveIndexing(t, NewConf(testPath(), 0), indexItems, &disabled.Provider{})
		defer env.Cleanup()
		blkfileMgrWrapper := newTestBlockfileWrapper(env, "testledger")
		defer blkfileMgrWrapper.close()

		blocks := testutil.ConstructTestBlocks(t, 3)
		
		blkfileMgrWrapper.addBlocks(blocks)
		blockfileMgr := blkfileMgrWrapper.blockfileMgr

		
		
		block, err := blockfileMgr.retrieveBlockByHash(protoutil.BlockHeaderHash(blocks[0].Header))
		if containsAttr(indexItems, blkstorage.IndexableAttrBlockHash) {
			assert.NoError(t, err, "Error while retrieving block by hash")
			assert.Equal(t, blocks[0], block)
		} else {
			assert.Exactly(t, blkstorage.ErrAttrNotIndexed, err)
		}

		
		block, err = blockfileMgr.retrieveBlockByNumber(0)
		if containsAttr(indexItems, blkstorage.IndexableAttrBlockNum) {
			assert.NoError(t, err, "Error while retrieving block by number")
			assert.Equal(t, blocks[0], block)
		} else {
			assert.Exactly(t, blkstorage.ErrAttrNotIndexed, err)
		}

		
		txid, err := protoutil.GetOrComputeTxIDFromEnvelope(blocks[0].Data.Data[0])
		assert.NoError(t, err)
		txEnvelope, err := blockfileMgr.retrieveTransactionByID(txid)
		if containsAttr(indexItems, blkstorage.IndexableAttrTxID) {
			assert.NoError(t, err, "Error while retrieving tx by id")
			txEnvelopeBytes := blocks[0].Data.Data[0]
			txEnvelopeOrig, err := protoutil.GetEnvelopeFromBlock(txEnvelopeBytes)
			assert.NoError(t, err)
			assert.Equal(t, txEnvelopeOrig, txEnvelope)
		} else {
			assert.Exactly(t, blkstorage.ErrAttrNotIndexed, err)
		}

		
		txEnvelope2, err := blockfileMgr.retrieveTransactionByBlockNumTranNum(0, 0)
		if containsAttr(indexItems, blkstorage.IndexableAttrBlockNumTranNum) {
			assert.NoError(t, err, "Error while retrieving tx by blockNum and tranNum")
			txEnvelopeBytes2 := blocks[0].Data.Data[0]
			txEnvelopeOrig2, err2 := protoutil.GetEnvelopeFromBlock(txEnvelopeBytes2)
			assert.NoError(t, err2)
			assert.Equal(t, txEnvelopeOrig2, txEnvelope2)
		} else {
			assert.Exactly(t, blkstorage.ErrAttrNotIndexed, err)
		}

		
		txid, err = protoutil.GetOrComputeTxIDFromEnvelope(blocks[0].Data.Data[0])
		assert.NoError(t, err)
		block, err = blockfileMgr.retrieveBlockByTxID(txid)
		if containsAttr(indexItems, blkstorage.IndexableAttrBlockTxID) {
			assert.NoError(t, err, "Error while retrieving block by txID")
			assert.Equal(t, block, blocks[0])
		} else {
			assert.Exactly(t, blkstorage.ErrAttrNotIndexed, err)
		}

		for _, block := range blocks {
			flags := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

			for idx, d := range block.Data.Data {
				txid, err = protoutil.GetOrComputeTxIDFromEnvelope(d)
				assert.NoError(t, err)

				reason, err := blockfileMgr.retrieveTxValidationCodeByTxID(txid)

				if containsAttr(indexItems, blkstorage.IndexableAttrTxValidationCode) {
					assert.NoError(t, err, "Error while retrieving tx validation code by txID")

					reasonFromFlags := flags.Flag(idx)

					assert.Equal(t, reasonFromFlags, reason)
				} else {
					assert.Exactly(t, blkstorage.ErrAttrNotIndexed, err)
				}
			}
		}
	})
}

func containsAttr(indexItems []blkstorage.IndexableAttr, attr blkstorage.IndexableAttr) bool {
	for _, element := range indexItems {
		if element == attr {
			return true
		}
	}
	return false
}
