/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/pkg/errors"
)


func UpgradeDataFormat(rootFSPath string) error {
	fileLockPath := fileLockPath(rootFSPath)
	fileLock := leveldbhelper.NewFileLock(fileLockPath)
	if err := fileLock.Lock(); err != nil {
		return errors.Wrap(err, "as another peer node command is executing,"+
			" wait for that command to complete its execution or terminate it before retrying")
	}
	defer fileLock.Unlock()

	
	
	return UpgradeIDStoreFormat(rootFSPath)
}


func UpgradeIDStoreFormat(rootFSPath string) error {
	logger.Debugf("Attempting to upgrade idStore data format to current format %s", string(idStoreFormatVersion))

	dbPath := LedgerProviderPath(rootFSPath)
	db := leveldbhelper.CreateDB(&leveldbhelper.Conf{DBPath: dbPath})
	db.Open()
	defer db.Close()

	idStore := &idStore{db, dbPath}
	return idStore.upgradeFormat()
}
