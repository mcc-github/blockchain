/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package fsblkstorage

import (
	"os"

	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

type rollbackMgr struct {
	ledgerID       string
	ledgerDir      string
	indexDir       string
	dbProvider     *leveldbhelper.Provider
	indexStore     *blockIndex
	targetBlockNum uint64
}


func Rollback(blockStorageDir, ledgerID string, targetBlockNum uint64, indexConfig *blkstorage.IndexConfig) error {
	r, err := newRollbackMgr(blockStorageDir, ledgerID, indexConfig, targetBlockNum)
	if err != nil {
		return err
	}
	defer r.dbProvider.Close()

	if err := recordHeightIfGreaterThanPreviousRecording(r.ledgerDir); err != nil {
		return err
	}

	logger.Infof("Rolling back block index to block number [%d]", targetBlockNum)
	if err := r.rollbackBlockIndex(); err != nil {
		return err
	}

	logger.Infof("Rolling back block files to block number [%d]", targetBlockNum)
	if err := r.rollbackBlockFiles(); err != nil {
		return err
	}

	return nil
}

func newRollbackMgr(blockStorageDir, ledgerID string, indexConfig *blkstorage.IndexConfig, targetBlockNum uint64) (*rollbackMgr, error) {
	r := &rollbackMgr{}

	r.ledgerID = ledgerID
	conf := &Conf{blockStorageDir: blockStorageDir}
	r.ledgerDir = conf.getLedgerBlockDir(ledgerID)
	r.targetBlockNum = targetBlockNum

	r.indexDir = conf.getIndexDir()
	r.dbProvider = leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: r.indexDir})
	var err error
	indexDB := r.dbProvider.GetDBHandle(ledgerID)
	r.indexStore, err = newBlockIndex(indexConfig, indexDB)
	return r, err
}

func (r *rollbackMgr) rollbackBlockIndex() error {
	lastBlockNumber, err := r.indexStore.getLastBlockIndexed()
	if err == errIndexEmpty {
		return nil
	}
	if err != nil {
		return err
	}

	
	
	
	
	
	
	
	batchLimit := uint64(10)

	
	
	start, end := r.targetBlockNum+1, lastBlockNumber
	for end >= start {
		if end-start >= batchLimit {
			start = end - batchLimit + 1 
		}
		logger.Infof("Deleting index associated with block number [%d] to [%d]", start, end)
		if err := r.deleteIndexEntriesRange(start, end); err != nil {
			return err
		}
		start, end = r.targetBlockNum+1, start-1
	}
	return nil
}

func (r *rollbackMgr) deleteIndexEntriesRange(startBlkNum, endBlkNum uint64) error {
	
	
	
	
	
	batch := leveldbhelper.NewUpdateBatch()
	lp, err := r.indexStore.getBlockLocByBlockNum(startBlkNum)
	if err != nil {
		return err
	}
	stream, err := newBlockStream(r.ledgerDir, lp.fileSuffixNum, int64(lp.offset), -1)
	defer stream.close()

	numberOfBlocksToRetrieve := endBlkNum - startBlkNum + 1
	for numberOfBlocksToRetrieve > 0 {
		blockBytes, placementInfo, err := stream.nextBlockBytesAndPlacementInfo()
		if err != nil {
			return err
		}

		blockInfo, err := extractSerializedBlockInfo(blockBytes)
		if err != nil {
			return err
		}

		err = populateBlockInfoWithDuplicateTxids(blockInfo, placementInfo, r.indexStore)
		if err != nil {
			return err
		}
		addIndexEntriesToBeDeleted(batch, blockInfo, r.indexStore)
		numberOfBlocksToRetrieve--
	}

	batch.Put(indexCheckpointKey, encodeBlockNum(startBlkNum-1))
	return r.indexStore.db.WriteBatch(batch, true)
}

func populateBlockInfoWithDuplicateTxids(blockInfo *serializedBlockInfo, placementInfo *blockPlacementInfo, indexStore *blockIndex) error {
	
	if !indexStore.isAttributeIndexed(blkstorage.IndexableAttrBlockTxID) {
		return nil
	}

	for _, txOffset := range blockInfo.txOffsets {
		blockLoc, err := indexStore.getBlockLocByTxID(txOffset.txID)
		
		
		
		
		if err == blkstorage.ErrNotFoundInIndex {
			logger.Warnf("TxID [%s] not found in index... skipping this TxID", txOffset.txID)
			continue
		}
		if err != nil {
			return err
		}
		
		if blockLoc.fileSuffixNum != placementInfo.fileNum || blockLoc.offset != int(placementInfo.blockStartOffset) {
			txOffset.isDuplicate = true
		}
	}
	return nil
}

func addIndexEntriesToBeDeleted(batch *leveldbhelper.UpdateBatch, blockInfo *serializedBlockInfo, indexStore *blockIndex) error {
	if indexStore.isAttributeIndexed(blkstorage.IndexableAttrBlockHash) {
		batch.Delete(constructBlockHashKey(protoutil.BlockHeaderHash(blockInfo.blockHeader)))
	}

	if indexStore.isAttributeIndexed(blkstorage.IndexableAttrBlockNum) {
		batch.Delete(constructBlockNumKey(blockInfo.blockHeader.Number))
	}

	if indexStore.isAttributeIndexed(blkstorage.IndexableAttrBlockNumTranNum) {
		for txIndex := range blockInfo.txOffsets {
			batch.Delete(constructBlockNumTranNumKey(blockInfo.blockHeader.Number, uint64(txIndex)))
		}
	}

	for _, txOffset := range blockInfo.txOffsets {
		if txOffset.isDuplicate {
			continue
		}

		if indexStore.isAttributeIndexed(blkstorage.IndexableAttrTxID) {
			batch.Delete(constructTxIDKey(txOffset.txID))
		}

		if indexStore.isAttributeIndexed(blkstorage.IndexableAttrBlockTxID) {
			batch.Delete(constructBlockTxIDKey(txOffset.txID))
		}

		if indexStore.isAttributeIndexed(blkstorage.IndexableAttrTxValidationCode) {
			batch.Delete(constructTxValidationCodeIDKey(txOffset.txID))
		}
	}

	return nil
}

func (r *rollbackMgr) rollbackBlockFiles() error {
	logger.Infof("Deleting checkpointInfo")
	if err := r.indexStore.db.Delete(blkMgrInfoKey, true); err != nil {
		return err
	}
	
	targetFileNum, err := binarySearchFileNumForBlock(r.ledgerDir, r.targetBlockNum)
	if err != nil {
		return err
	}
	lastFileNum, err := retrieveLastFileSuffix(r.ledgerDir)
	if err != nil {
		return err
	}

	logger.Infof("Removing all block files with suffixNum in the range [%d] to [%d]",
		targetFileNum+1, lastFileNum)

	for n := lastFileNum; n >= targetFileNum+1; n-- {
		filepath := deriveBlockfilePath(r.ledgerDir, n)
		if err := os.Remove(filepath); err != nil {
			return errors.Wrapf(err, "error removing the block file [%s]", filepath)
		}
	}

	logger.Infof("Truncating block file [%d] to the end boundary of block number [%d]", targetFileNum, r.targetBlockNum)
	endOffset, err := calculateEndOffSet(r.ledgerDir, targetFileNum, r.targetBlockNum)
	if err != nil {
		return err
	}

	filePath := deriveBlockfilePath(r.ledgerDir, targetFileNum)
	if err := os.Truncate(filePath, endOffset); err != nil {
		return errors.Wrapf(err, "error trucating the block file [%s]", filePath)
	}

	return nil
}

func calculateEndOffSet(ledgerDir string, targetBlkFileNum int, blockNum uint64) (int64, error) {
	stream, err := newBlockfileStream(ledgerDir, targetBlkFileNum, 0)
	if err != nil {
		return 0, err
	}
	defer stream.close()
	for {
		blockBytes, err := stream.nextBlockBytes()
		if err != nil {
			return 0, err
		}
		blockInfo, err := extractSerializedBlockInfo(blockBytes)
		if err != nil {
			return 0, err
		}
		if blockInfo.blockHeader.Number == blockNum {
			break
		}
	}
	return stream.currentOffset, nil
}



func ValidateRollbackParams(blockStorageDir, ledgerID string, targetBlockNum uint64) error {
	logger.Infof("Validating the rollback parameters: ledgerID [%s], block number [%d]",
		ledgerID, targetBlockNum)
	conf := &Conf{blockStorageDir: blockStorageDir}
	ledgerDir := conf.getLedgerBlockDir(ledgerID)
	if err := validateLedgerID(ledgerDir, ledgerID); err != nil {
		return err
	}
	if err := validateTargetBlkNum(ledgerDir, targetBlockNum); err != nil {
		return err
	}
	return nil

}

func validateLedgerID(ledgerDir, ledgerID string) error {
	logger.Debugf("Validating the existance of ledgerID [%s]", ledgerID)
	exists, _, err := util.FileExists(ledgerDir)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("ledgerID [%s] does not exist", ledgerID)
	}
	return nil
}

func validateTargetBlkNum(ledgerDir string, targetBlockNum uint64) error {
	logger.Debugf("Validating the given block number [%d] agains the ledger block height", targetBlockNum)
	cpInfo, err := constructCheckpointInfoFromBlockFiles(ledgerDir)
	if err != nil {
		return err
	}
	if cpInfo.lastBlockNumber <= targetBlockNum {
		return errors.Errorf("target block number [%d] should be less than the biggest block number [%d]",
			targetBlockNum, cpInfo.lastBlockNumber)
	}
	return nil
}
