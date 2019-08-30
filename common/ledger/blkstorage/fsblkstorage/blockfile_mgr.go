/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fsblkstorage

import (
	"bytes"
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("fsblkstorage")

const (
	blockfilePrefix = "blockfile_"
)

var (
	blkMgrInfoKey = []byte("blkMgrInfo")
)

type blockfileMgr struct {
	rootDir           string
	conf              *Conf
	db                *leveldbhelper.DBHandle
	index             index
	cpInfo            *checkpointInfo
	cpInfoCond        *sync.Cond
	currentFileWriter *blockfileWriter
	bcInfo            atomic.Value
}


func newBlockfileMgr(id string, conf *Conf, indexConfig *blkstorage.IndexConfig, indexStore *leveldbhelper.DBHandle) *blockfileMgr {
	logger.Debugf("newBlockfileMgr() initializing file-based block storage for ledger: %s ", id)
	
	rootDir := conf.getLedgerBlockDir(id)
	_, err := util.CreateDirIfMissing(rootDir)
	if err != nil {
		panic(fmt.Sprintf("Error creating block storage root dir [%s]: %s", rootDir, err))
	}
	
	mgr := &blockfileMgr{rootDir: rootDir, conf: conf, db: indexStore}

	
	
	
	cpInfo, err := mgr.loadCurrentInfo()
	if err != nil {
		panic(fmt.Sprintf("Could not get block file info for current block file from db: %s", err))
	}
	if cpInfo == nil {
		logger.Info(`Getting block information from block storage`)
		if cpInfo, err = constructCheckpointInfoFromBlockFiles(rootDir); err != nil {
			panic(fmt.Sprintf("Could not build checkpoint info from block files: %s", err))
		}
		logger.Debugf("Info constructed by scanning the blocks dir = %s", spew.Sdump(cpInfo))
	} else {
		logger.Debug(`Synching block information from block storage (if needed)`)
		syncCPInfoFromFS(rootDir, cpInfo)
	}
	err = mgr.saveCurrentInfo(cpInfo, true)
	if err != nil {
		panic(fmt.Sprintf("Could not save next block file info to db: %s", err))
	}

	
	
	currentFileWriter, err := newBlockfileWriter(deriveBlockfilePath(rootDir, cpInfo.latestFileChunkSuffixNum))
	if err != nil {
		panic(fmt.Sprintf("Could not open writer to current file: %s", err))
	}
	
	err = currentFileWriter.truncateFile(cpInfo.latestFileChunksize)
	if err != nil {
		panic(fmt.Sprintf("Could not truncate current file to known size in db: %s", err))
	}

	
	if mgr.index, err = newBlockIndex(indexConfig, indexStore); err != nil {
		panic(fmt.Sprintf("error in block index: %s", err))
	}

	
	mgr.cpInfo = cpInfo
	mgr.currentFileWriter = currentFileWriter
	
	
	mgr.cpInfoCond = sync.NewCond(&sync.Mutex{})

	
	bcInfo := &common.BlockchainInfo{
		Height:            0,
		CurrentBlockHash:  nil,
		PreviousBlockHash: nil}

	if !cpInfo.isChainEmpty {
		
		mgr.syncIndex()
		lastBlockHeader, err := mgr.retrieveBlockHeaderByNumber(cpInfo.lastBlockNumber)
		if err != nil {
			panic(fmt.Sprintf("Could not retrieve header of the last block form file: %s", err))
		}
		lastBlockHash := protoutil.BlockHeaderHash(lastBlockHeader)
		previousBlockHash := lastBlockHeader.PreviousHash
		bcInfo = &common.BlockchainInfo{
			Height:            cpInfo.lastBlockNumber + 1,
			CurrentBlockHash:  lastBlockHash,
			PreviousBlockHash: previousBlockHash}
	}
	mgr.bcInfo.Store(bcInfo)
	return mgr
}





func syncCPInfoFromFS(rootDir string, cpInfo *checkpointInfo) {
	logger.Debugf("Starting checkpoint=%s", cpInfo)
	
	filePath := deriveBlockfilePath(rootDir, cpInfo.latestFileChunkSuffixNum)
	exists, size, err := util.FileExists(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error in checking whether file [%s] exists: %s", filePath, err))
	}
	logger.Debugf("status of file [%s]: exists=[%t], size=[%d]", filePath, exists, size)
	
	
	
	if !exists || int(size) == cpInfo.latestFileChunksize {
		
		return
	}
	
	_, endOffsetLastBlock, numBlocks, err := scanForLastCompleteBlock(
		rootDir, cpInfo.latestFileChunkSuffixNum, int64(cpInfo.latestFileChunksize))
	if err != nil {
		panic(fmt.Sprintf("Could not open current file for detecting last block in the file: %s", err))
	}
	cpInfo.latestFileChunksize = int(endOffsetLastBlock)
	if numBlocks == 0 {
		return
	}
	
	if cpInfo.isChainEmpty {
		cpInfo.lastBlockNumber = uint64(numBlocks - 1)
	} else {
		cpInfo.lastBlockNumber += uint64(numBlocks)
	}
	cpInfo.isChainEmpty = false
	logger.Debugf("Checkpoint after updates by scanning the last file segment:%s", cpInfo)
}

func deriveBlockfilePath(rootDir string, suffixNum int) string {
	return rootDir + "/" + blockfilePrefix + fmt.Sprintf("%06d", suffixNum)
}

func (mgr *blockfileMgr) close() {
	mgr.currentFileWriter.close()
}

func (mgr *blockfileMgr) moveToNextFile() {
	cpInfo := &checkpointInfo{
		latestFileChunkSuffixNum: mgr.cpInfo.latestFileChunkSuffixNum + 1,
		latestFileChunksize:      0,
		lastBlockNumber:          mgr.cpInfo.lastBlockNumber}

	nextFileWriter, err := newBlockfileWriter(
		deriveBlockfilePath(mgr.rootDir, cpInfo.latestFileChunkSuffixNum))

	if err != nil {
		panic(fmt.Sprintf("Could not open writer to next file: %s", err))
	}
	mgr.currentFileWriter.close()
	err = mgr.saveCurrentInfo(cpInfo, true)
	if err != nil {
		panic(fmt.Sprintf("Could not save next block file info to db: %s", err))
	}
	mgr.currentFileWriter = nextFileWriter
	mgr.updateCheckpoint(cpInfo)
}

func (mgr *blockfileMgr) addBlock(block *common.Block) error {
	bcInfo := mgr.getBlockchainInfo()
	if block.Header.Number != bcInfo.Height {
		return errors.Errorf(
			"block number should have been %d but was %d",
			mgr.getBlockchainInfo().Height, block.Header.Number,
		)
	}

	
	
	
	
	if !bytes.Equal(block.Header.PreviousHash, bcInfo.CurrentBlockHash) {
		return errors.Errorf(
			"unexpected Previous block hash. Expected PreviousHash = [%x], PreviousHash referred in the latest block= [%x]",
			bcInfo.CurrentBlockHash, block.Header.PreviousHash,
		)
	}
	blockBytes, info, err := serializeBlock(block)
	if err != nil {
		return errors.WithMessage(err, "error serializing block")
	}
	blockHash := protoutil.BlockHeaderHash(block.Header)
	
	txOffsets := info.txOffsets
	currentOffset := mgr.cpInfo.latestFileChunksize

	blockBytesLen := len(blockBytes)
	blockBytesEncodedLen := proto.EncodeVarint(uint64(blockBytesLen))
	totalBytesToAppend := blockBytesLen + len(blockBytesEncodedLen)

	
	
	if currentOffset+totalBytesToAppend > mgr.conf.maxBlockfileSize {
		mgr.moveToNextFile()
		currentOffset = 0
	}
	
	err = mgr.currentFileWriter.append(blockBytesEncodedLen, false)
	if err == nil {
		
		err = mgr.currentFileWriter.append(blockBytes, true)
	}
	if err != nil {
		truncateErr := mgr.currentFileWriter.truncateFile(mgr.cpInfo.latestFileChunksize)
		if truncateErr != nil {
			panic(fmt.Sprintf("Could not truncate current file to known size after an error during block append: %s", err))
		}
		return errors.WithMessage(err, "error appending block to file")
	}

	
	currentCPInfo := mgr.cpInfo
	newCPInfo := &checkpointInfo{
		latestFileChunkSuffixNum: currentCPInfo.latestFileChunkSuffixNum,
		latestFileChunksize:      currentCPInfo.latestFileChunksize + totalBytesToAppend,
		isChainEmpty:             false,
		lastBlockNumber:          block.Header.Number}
	
	if err = mgr.saveCurrentInfo(newCPInfo, false); err != nil {
		truncateErr := mgr.currentFileWriter.truncateFile(currentCPInfo.latestFileChunksize)
		if truncateErr != nil {
			panic(fmt.Sprintf("Error in truncating current file to known size after an error in saving checkpoint info: %s", err))
		}
		return errors.WithMessage(err, "error saving current file info to db")
	}

	
	blockFLP := &fileLocPointer{fileSuffixNum: newCPInfo.latestFileChunkSuffixNum}
	blockFLP.offset = currentOffset
	
	for _, txOffset := range txOffsets {
		txOffset.loc.offset += len(blockBytesEncodedLen)
	}
	
	if err = mgr.index.indexBlock(&blockIdxInfo{
		blockNum: block.Header.Number, blockHash: blockHash,
		flp: blockFLP, txOffsets: txOffsets, metadata: block.Metadata}); err != nil {
		return err
	}

	
	mgr.updateCheckpoint(newCPInfo)
	mgr.updateBlockchainInfo(blockHash, block)
	return nil
}

func (mgr *blockfileMgr) syncIndex() error {
	var lastBlockIndexed uint64
	var indexEmpty bool
	var err error
	
	if lastBlockIndexed, err = mgr.index.getLastBlockIndexed(); err != nil {
		if err != errIndexEmpty {
			return err
		}
		indexEmpty = true
	}

	
	startFileNum := 0
	startOffset := 0
	skipFirstBlock := false
	
	endFileNum := mgr.cpInfo.latestFileChunkSuffixNum
	startingBlockNum := uint64(0)

	
	if !indexEmpty {
		if lastBlockIndexed == mgr.cpInfo.lastBlockNumber {
			logger.Debug("Both the block files and indices are in sync.")
			return nil
		}
		logger.Debugf("Last block indexed [%d], Last block present in block files [%d]", lastBlockIndexed, mgr.cpInfo.lastBlockNumber)
		var flp *fileLocPointer
		if flp, err = mgr.index.getBlockLocByBlockNum(lastBlockIndexed); err != nil {
			return err
		}
		startFileNum = flp.fileSuffixNum
		startOffset = flp.locPointer.offset
		skipFirstBlock = true
		startingBlockNum = lastBlockIndexed + 1
	} else {
		logger.Debugf("No block indexed, Last block present in block files=[%d]", mgr.cpInfo.lastBlockNumber)
	}

	logger.Infof("Start building index from block [%d] to last block [%d]", startingBlockNum, mgr.cpInfo.lastBlockNumber)

	
	var stream *blockStream
	if stream, err = newBlockStream(mgr.rootDir, startFileNum, int64(startOffset), endFileNum); err != nil {
		return err
	}
	var blockBytes []byte
	var blockPlacementInfo *blockPlacementInfo

	if skipFirstBlock {
		if blockBytes, _, err = stream.nextBlockBytesAndPlacementInfo(); err != nil {
			return err
		}
		if blockBytes == nil {
			return errors.Errorf("block bytes for block num = [%d] should not be nil here. The indexes for the block are already present",
				lastBlockIndexed)
		}
	}

	
	
	
	blockIdxInfo := &blockIdxInfo{}
	for {
		if blockBytes, blockPlacementInfo, err = stream.nextBlockBytesAndPlacementInfo(); err != nil {
			return err
		}
		if blockBytes == nil {
			break
		}
		info, err := extractSerializedBlockInfo(blockBytes)
		if err != nil {
			return err
		}

		
		
		numBytesToShift := int(blockPlacementInfo.blockBytesOffset - blockPlacementInfo.blockStartOffset)
		for _, offset := range info.txOffsets {
			offset.loc.offset += numBytesToShift
		}

		
		blockIdxInfo.blockHash = protoutil.BlockHeaderHash(info.blockHeader)
		blockIdxInfo.blockNum = info.blockHeader.Number
		blockIdxInfo.flp = &fileLocPointer{fileSuffixNum: blockPlacementInfo.fileNum,
			locPointer: locPointer{offset: int(blockPlacementInfo.blockStartOffset)}}
		blockIdxInfo.txOffsets = info.txOffsets
		blockIdxInfo.metadata = info.metadata

		logger.Debugf("syncIndex() indexing block [%d]", blockIdxInfo.blockNum)
		if err = mgr.index.indexBlock(blockIdxInfo); err != nil {
			return err
		}
		if blockIdxInfo.blockNum%10000 == 0 {
			logger.Infof("Indexed block number [%d]", blockIdxInfo.blockNum)
		}
	}
	logger.Infof("Finished building index. Last block indexed [%d]", blockIdxInfo.blockNum)
	return nil
}

func (mgr *blockfileMgr) getBlockchainInfo() *common.BlockchainInfo {
	return mgr.bcInfo.Load().(*common.BlockchainInfo)
}

func (mgr *blockfileMgr) updateCheckpoint(cpInfo *checkpointInfo) {
	mgr.cpInfoCond.L.Lock()
	defer mgr.cpInfoCond.L.Unlock()
	mgr.cpInfo = cpInfo
	logger.Debugf("Broadcasting about update checkpointInfo: %s", cpInfo)
	mgr.cpInfoCond.Broadcast()
}

func (mgr *blockfileMgr) updateBlockchainInfo(latestBlockHash []byte, latestBlock *common.Block) {
	currentBCInfo := mgr.getBlockchainInfo()
	newBCInfo := &common.BlockchainInfo{
		Height:            currentBCInfo.Height + 1,
		CurrentBlockHash:  latestBlockHash,
		PreviousBlockHash: latestBlock.Header.PreviousHash}

	mgr.bcInfo.Store(newBCInfo)
}

func (mgr *blockfileMgr) retrieveBlockByHash(blockHash []byte) (*common.Block, error) {
	logger.Debugf("retrieveBlockByHash() - blockHash = [%#v]", blockHash)
	loc, err := mgr.index.getBlockLocByHash(blockHash)
	if err != nil {
		return nil, err
	}
	return mgr.fetchBlock(loc)
}

func (mgr *blockfileMgr) retrieveBlockByNumber(blockNum uint64) (*common.Block, error) {
	logger.Debugf("retrieveBlockByNumber() - blockNum = [%d]", blockNum)

	
	if blockNum == math.MaxUint64 {
		blockNum = mgr.getBlockchainInfo().Height - 1
	}

	loc, err := mgr.index.getBlockLocByBlockNum(blockNum)
	if err != nil {
		return nil, err
	}
	return mgr.fetchBlock(loc)
}

func (mgr *blockfileMgr) retrieveBlockByTxID(txID string) (*common.Block, error) {
	logger.Debugf("retrieveBlockByTxID() - txID = [%s]", txID)

	loc, err := mgr.index.getBlockLocByTxID(txID)

	if err != nil {
		return nil, err
	}
	return mgr.fetchBlock(loc)
}

func (mgr *blockfileMgr) retrieveTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	logger.Debugf("retrieveTxValidationCodeByTxID() - txID = [%s]", txID)
	return mgr.index.getTxValidationCodeByTxID(txID)
}

func (mgr *blockfileMgr) retrieveBlockHeaderByNumber(blockNum uint64) (*common.BlockHeader, error) {
	logger.Debugf("retrieveBlockHeaderByNumber() - blockNum = [%d]", blockNum)
	loc, err := mgr.index.getBlockLocByBlockNum(blockNum)
	if err != nil {
		return nil, err
	}
	blockBytes, err := mgr.fetchBlockBytes(loc)
	if err != nil {
		return nil, err
	}
	info, err := extractSerializedBlockInfo(blockBytes)
	if err != nil {
		return nil, err
	}
	return info.blockHeader, nil
}

func (mgr *blockfileMgr) retrieveBlocks(startNum uint64) (*blocksItr, error) {
	return newBlockItr(mgr, startNum), nil
}

func (mgr *blockfileMgr) retrieveTransactionByID(txID string) (*common.Envelope, error) {
	logger.Debugf("retrieveTransactionByID() - txId = [%s]", txID)
	loc, err := mgr.index.getTxLoc(txID)
	if err != nil {
		return nil, err
	}
	return mgr.fetchTransactionEnvelope(loc)
}

func (mgr *blockfileMgr) retrieveTransactionByBlockNumTranNum(blockNum uint64, tranNum uint64) (*common.Envelope, error) {
	logger.Debugf("retrieveTransactionByBlockNumTranNum() - blockNum = [%d], tranNum = [%d]", blockNum, tranNum)
	loc, err := mgr.index.getTXLocByBlockNumTranNum(blockNum, tranNum)
	if err != nil {
		return nil, err
	}
	return mgr.fetchTransactionEnvelope(loc)
}

func (mgr *blockfileMgr) fetchBlock(lp *fileLocPointer) (*common.Block, error) {
	blockBytes, err := mgr.fetchBlockBytes(lp)
	if err != nil {
		return nil, err
	}
	block, err := deserializeBlock(blockBytes)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (mgr *blockfileMgr) fetchTransactionEnvelope(lp *fileLocPointer) (*common.Envelope, error) {
	logger.Debugf("Entering fetchTransactionEnvelope() %v\n", lp)
	var err error
	var txEnvelopeBytes []byte
	if txEnvelopeBytes, err = mgr.fetchRawBytes(lp); err != nil {
		return nil, err
	}
	_, n := proto.DecodeVarint(txEnvelopeBytes)
	return protoutil.GetEnvelopeFromBlock(txEnvelopeBytes[n:])
}

func (mgr *blockfileMgr) fetchBlockBytes(lp *fileLocPointer) ([]byte, error) {
	stream, err := newBlockfileStream(mgr.rootDir, lp.fileSuffixNum, int64(lp.offset))
	if err != nil {
		return nil, err
	}
	defer stream.close()
	b, err := stream.nextBlockBytes()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (mgr *blockfileMgr) fetchRawBytes(lp *fileLocPointer) ([]byte, error) {
	filePath := deriveBlockfilePath(mgr.rootDir, lp.fileSuffixNum)
	reader, err := newBlockfileReader(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.close()
	b, err := reader.read(lp.offset, lp.bytesLength)
	if err != nil {
		return nil, err
	}
	return b, nil
}


func (mgr *blockfileMgr) loadCurrentInfo() (*checkpointInfo, error) {
	var b []byte
	var err error
	if b, err = mgr.db.Get(blkMgrInfoKey); b == nil || err != nil {
		return nil, err
	}
	i := &checkpointInfo{}
	if err = i.unmarshal(b); err != nil {
		return nil, err
	}
	logger.Debugf("loaded checkpointInfo:%s", i)
	return i, nil
}

func (mgr *blockfileMgr) saveCurrentInfo(i *checkpointInfo, sync bool) error {
	b, err := i.marshal()
	if err != nil {
		return err
	}
	if err = mgr.db.Put(blkMgrInfoKey, b, sync); err != nil {
		return err
	}
	return nil
}



func scanForLastCompleteBlock(rootDir string, fileNum int, startingOffset int64) ([]byte, int64, int, error) {
	
	numBlocks := 0
	var lastBlockBytes []byte
	blockStream, errOpen := newBlockfileStream(rootDir, fileNum, startingOffset)
	if errOpen != nil {
		return nil, 0, 0, errOpen
	}
	defer blockStream.close()
	var errRead error
	var blockBytes []byte
	for {
		blockBytes, errRead = blockStream.nextBlockBytes()
		if blockBytes == nil || errRead != nil {
			break
		}
		lastBlockBytes = blockBytes
		numBlocks++
	}
	if errRead == ErrUnexpectedEndOfBlockfile {
		logger.Debugf(`Error:%s
		The error may happen if a crash has happened during block appending.
		Resetting error to nil and returning current offset as a last complete block's end offset`, errRead)
		errRead = nil
	}
	logger.Debugf("scanForLastCompleteBlock(): last complete block ends at offset=[%d]", blockStream.currentOffset)
	return lastBlockBytes, blockStream.currentOffset, numBlocks, errRead
}


type checkpointInfo struct {
	latestFileChunkSuffixNum int
	latestFileChunksize      int
	isChainEmpty             bool
	lastBlockNumber          uint64
}

func (i *checkpointInfo) marshal() ([]byte, error) {
	buffer := proto.NewBuffer([]byte{})
	var err error
	if err = buffer.EncodeVarint(uint64(i.latestFileChunkSuffixNum)); err != nil {
		return nil, errors.Wrapf(err, "error encoding the latestFileChunkSuffixNum [%d]", i.latestFileChunkSuffixNum)
	}
	if err = buffer.EncodeVarint(uint64(i.latestFileChunksize)); err != nil {
		return nil, errors.Wrapf(err, "error encoding the latestFileChunksize [%d]", i.latestFileChunksize)
	}
	if err = buffer.EncodeVarint(i.lastBlockNumber); err != nil {
		return nil, errors.Wrapf(err, "error encoding the lastBlockNumber [%d]", i.lastBlockNumber)
	}
	var chainEmptyMarker uint64
	if i.isChainEmpty {
		chainEmptyMarker = 1
	}
	if err = buffer.EncodeVarint(chainEmptyMarker); err != nil {
		return nil, errors.Wrapf(err, "error encoding chainEmptyMarker [%d]", chainEmptyMarker)
	}
	return buffer.Bytes(), nil
}

func (i *checkpointInfo) unmarshal(b []byte) error {
	buffer := proto.NewBuffer(b)
	var val uint64
	var chainEmptyMarker uint64
	var err error

	if val, err = buffer.DecodeVarint(); err != nil {
		return err
	}
	i.latestFileChunkSuffixNum = int(val)

	if val, err = buffer.DecodeVarint(); err != nil {
		return err
	}
	i.latestFileChunksize = int(val)

	if val, err = buffer.DecodeVarint(); err != nil {
		return err
	}
	i.lastBlockNumber = val
	if chainEmptyMarker, err = buffer.DecodeVarint(); err != nil {
		return err
	}
	i.isChainEmpty = chainEmptyMarker == 1
	return nil
}

func (i *checkpointInfo) String() string {
	return fmt.Sprintf("latestFileChunkSuffixNum=[%d], latestFileChunksize=[%d], isChainEmpty=[%t], lastBlockNumber=[%d]",
		i.latestFileChunkSuffixNum, i.latestFileChunksize, i.isChainEmpty, i.lastBlockNumber)
}
