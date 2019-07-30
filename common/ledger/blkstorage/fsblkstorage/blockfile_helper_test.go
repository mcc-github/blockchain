/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fsblkstorage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/stretchr/testify/assert"
)

func TestConstructCheckpointInfoFromBlockFiles(t *testing.T) {
	testPath := "/tmp/tests/blockchain/common/ledger/blkstorage/fsblkstorage"
	ledgerid := "testLedger"
	conf := NewConf(testPath, 0)
	blkStoreDir := conf.getLedgerBlockDir(ledgerid)
	env := newTestEnv(t, conf)
	util.CreateDirIfMissing(blkStoreDir)
	defer env.Cleanup()

	
	cpInfo, err := constructCheckpointInfoFromBlockFiles(blkStoreDir)
	assert.NoError(t, err)
	assert.Equal(t, &checkpointInfo{isChainEmpty: true, lastBlockNumber: 0, latestFileChunksize: 0, latestFileChunkSuffixNum: 0}, cpInfo)

	w := newTestBlockfileWrapper(env, ledgerid)
	defer w.close()
	blockfileMgr := w.blockfileMgr
	bg, gb := testutil.NewBlockGenerator(t, ledgerid, false)

	
	blockfileMgr.addBlock(gb)
	for _, blk := range bg.NextTestBlocks(3) {
		blockfileMgr.addBlock(blk)
	}
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)

	
	blockfileMgr.moveToNextFile()
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)

	
	for _, blk := range bg.NextTestBlocks(3) {
		blockfileMgr.addBlock(blk)
	}
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)

	
	lastTestBlk := bg.NextTestBlocks(1)[0]
	blockBytes, _, err := serializeBlock(lastTestBlk)
	assert.NoError(t, err)
	partialByte := append(proto.EncodeVarint(uint64(len(blockBytes))), blockBytes[len(blockBytes)/2:]...)
	blockfileMgr.currentFileWriter.append(partialByte, true)
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)

	
	cpInfoBeforeClose := blockfileMgr.cpInfo
	w.close()
	env.provider.Close()
	indexFolder := conf.getIndexDir()
	assert.NoError(t, os.RemoveAll(indexFolder))

	env = newTestEnv(t, conf)
	w = newTestBlockfileWrapper(env, ledgerid)
	blockfileMgr = w.blockfileMgr
	assert.Equal(t, cpInfoBeforeClose, blockfileMgr.cpInfo)

	lastBlkIndexed, err := blockfileMgr.index.getLastBlockIndexed()
	assert.NoError(t, err)
	assert.Equal(t, uint64(6), lastBlkIndexed)

	
	assert.NoError(t, blockfileMgr.addBlock(lastTestBlk))
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)
}

func TestBinarySearchBlockFileNum(t *testing.T) {
	blockStoreRootDir := testPath()
	blocks := testutil.ConstructTestBlocks(t, 100)
	maxFileSie := int(0.1 * float64(testutilEstimateTotalSizeOnDisk(t, blocks)))
	env := newTestEnv(t, NewConf(blockStoreRootDir, maxFileSie))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	blkfileMgr := blkfileMgrWrapper.blockfileMgr

	blkfileMgrWrapper.addBlocks(blocks)

	ledgerDir := (&Conf{blockStorageDir: blockStoreRootDir}).getLedgerBlockDir("testLedger")
	files, err := ioutil.ReadDir(ledgerDir)
	assert.NoError(t, err)
	assert.Len(t, files, 11)

	for i := uint64(0); i < 100; i++ {
		fileNum, err := binarySearchFileNumForBlock(ledgerDir, i)
		assert.NoError(t, err)
		locFromIndex, err := blkfileMgr.index.getBlockLocByBlockNum(i)
		assert.NoError(t, err)
		expectedFileNum := locFromIndex.fileSuffixNum
		assert.Equal(t, expectedFileNum, fileNum)
	}
}

func checkCPInfoFromFile(t *testing.T, blkStoreDir string, expectedCPInfo *checkpointInfo) {
	cpInfo, err := constructCheckpointInfoFromBlockFiles(blkStoreDir)
	assert.NoError(t, err)
	assert.Equal(t, expectedCPInfo, cpInfo)
}
