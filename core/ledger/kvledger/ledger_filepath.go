/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"path/filepath"
)

func fileLockPath(rootFSPath string) string {
	return filepath.Join(rootFSPath, "fileLock")
}

func ledgerProviderPath(rootFSPath string) string {
	return filepath.Join(rootFSPath, "ledgerProvider")
}


func BlockStorePath(rootFSPath string) string {
	return filepath.Join(rootFSPath, "chains")
}


func PvtDataStorePath(rootFSPath string) string {
	return filepath.Join(rootFSPath, "pvtdataStore")
}


func StateDBPath(rootFSPath string) string {
	return filepath.Join(rootFSPath, "stateLeveldb")
}


func HistoryDBPath(rootFSPath string) string {
	return filepath.Join(rootFSPath, "historyLeveldb")
}


func ConfigHistoryDBPath(rootFSPath string) string {
	return filepath.Join(rootFSPath, "configHistory")
}


func BookkeeperDBPath(rootFSPath string) string {
	return filepath.Join(rootFSPath, "bookkeeper")
}
