/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func dropDBs(rootFSPath string) error {
	
	
	
	
	
	
	
	
	
	if err := dropStateLevelDB(rootFSPath); err != nil {
		return err
	}
	if err := dropConfigHistoryDB(rootFSPath); err != nil {
		return err
	}
	if err := dropBookkeeperDB(rootFSPath); err != nil {
		return err
	}
	if err := dropHistoryDB(rootFSPath); err != nil {
		return err
	}
	return nil
}

func dropStateLevelDB(rootFSPath string) error {
	stateLeveldbPath := filepath.Join(rootFSPath, "stateLeveldb")
	logger.Infof("Dropping StateLevelDB at location [%s] ...if present", stateLeveldbPath)
	return os.RemoveAll(stateLeveldbPath)
}

func dropConfigHistoryDB(rootFSPath string) error {
	configHistoryDBPath := filepath.Join(rootFSPath, "configHistory")
	logger.Infof("Dropping ConfigHistoryDB at location [%s]", configHistoryDBPath)
	err := os.RemoveAll(configHistoryDBPath)
	return errors.Wrapf(err, "error removing the ConfigHistoryDB located at %s", configHistoryDBPath)
}

func dropBookkeeperDB(rootFSPath string) error {
	bookkeeperDBPath := filepath.Join(rootFSPath, "bookkeeper")
	logger.Infof("Dropping BookkeeperDB at location [%s]", bookkeeperDBPath)
	err := os.RemoveAll(bookkeeperDBPath)
	return errors.Wrapf(err, "error removing the BookkeeperDB located at %s", bookkeeperDBPath)
}

func dropHistoryDB(rootFSPath string) error {
	historyDBPath := filepath.Join(rootFSPath, "historyLeveldb")
	logger.Infof("Dropping HistoryDB at location [%s] ...if present", historyDBPath)
	return os.RemoveAll(historyDBPath)
}
