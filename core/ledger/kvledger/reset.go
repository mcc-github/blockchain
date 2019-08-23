/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"github.com/mcc-github/blockchain/common/ledger/blkstorage/fsblkstorage"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/pkg/errors"
)


func ResetAllKVLedgers(rootFSPath string) error {
	fileLockPath := fileLockPath(rootFSPath)
	fileLock := leveldbhelper.NewFileLock(fileLockPath)
	if err := fileLock.Lock(); err != nil {
		return errors.Wrap(err, "as another peer node command is executing,"+
			" wait for that command to complete its execution or terminate it before retrying")
	}
	defer fileLock.Unlock()

	logger.Info("Resetting all channel ledgers to genesis block")
	logger.Infof("Ledger data folder from config = [%s]", rootFSPath)
	if err := dropDBs(rootFSPath); err != nil {
		return err
	}
	if err := resetBlockStorage(rootFSPath); err != nil {
		return err
	}
	logger.Info("All channel ledgers have been successfully reset to the genesis block")
	return nil
}


func LoadPreResetHeight(rootFSPath string) (map[string]uint64, error) {
	blockstorePath := BlockStorePath(rootFSPath)
	logger.Infof("Loading prereset height from path [%s]", blockstorePath)
	return fsblkstorage.LoadPreResetHeight(blockstorePath)
}


func ClearPreResetHeight(rootFSPath string) error {
	blockstorePath := BlockStorePath(rootFSPath)
	logger.Infof("Clearing off prereset height files from path [%s]", blockstorePath)
	return fsblkstorage.ClearPreResetHeight(blockstorePath)
}

func resetBlockStorage(rootFSPath string) error {
	blockstorePath := BlockStorePath(rootFSPath)
	logger.Infof("Resetting BlockStore to genesis block at location [%s]", blockstorePath)
	return fsblkstorage.ResetBlockStore(blockstorePath)
}
