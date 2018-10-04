/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerconfig

import (
	"path/filepath"

	"github.com/mcc-github/blockchain/core/config"
	"github.com/spf13/viper"
)


func IsCouchDBEnabled() bool {
	stateDatabase := viper.GetString("ledger.state.stateDatabase")
	if stateDatabase == "CouchDB" {
		return true
	}
	return false
}

const confPeerFileSystemPath = "peer.fileSystemPath"
const confLedgersData = "ledgersData"
const confLedgerProvider = "ledgerProvider"
const confStateleveldb = "stateLeveldb"
const confHistoryLeveldb = "historyLeveldb"
const confBookkeeper = "bookkeeper"
const confConfigHistory = "configHistory"
const confChains = "chains"
const confPvtdataStore = "pvtdataStore"
const confTotalQueryLimit = "ledger.state.totalQueryLimit"
const confInternalQueryLimit = "ledger.state.couchDBConfig.internalQueryLimit"
const confEnableHistoryDatabase = "ledger.history.enableHistoryDatabase"
const confMaxBatchSize = "ledger.state.couchDBConfig.maxBatchUpdateSize"
const confAutoWarmIndexes = "ledger.state.couchDBConfig.autoWarmIndexes"
const confWarmIndexesAfterNBlocks = "ledger.state.couchDBConfig.warmIndexesAfterNBlocks"

var confCollElgProcMaxDbBatchSize = &conf{"ledger.pvtdataStore.collElgProcMaxDbBatchSize", 5000}
var confCollElgProcDbBatchesInterval = &conf{"ledger.pvtdataStore.collElgProcDbBatchesInterval", 1000}



func GetRootPath() string {
	sysPath := config.GetPath(confPeerFileSystemPath)
	return filepath.Join(sysPath, confLedgersData)
}


func GetLedgerProviderPath() string {
	return filepath.Join(GetRootPath(), confLedgerProvider)
}


func GetStateLevelDBPath() string {
	return filepath.Join(GetRootPath(), confStateleveldb)
}


func GetHistoryLevelDBPath() string {
	return filepath.Join(GetRootPath(), confHistoryLeveldb)
}


func GetBlockStorePath() string {
	return filepath.Join(GetRootPath(), confChains)
}


func GetPvtdataStorePath() string {
	return filepath.Join(GetRootPath(), confPvtdataStore)
}


func GetInternalBookkeeperPath() string {
	return filepath.Join(GetRootPath(), confBookkeeper)
}


func GetConfigHistoryPath() string {
	return filepath.Join(GetRootPath(), confConfigHistory)
}


func GetMaxBlockfileSize() int {
	return 64 * 1024 * 1024
}


func GetTotalQueryLimit() int {
	totalQueryLimit := viper.GetInt(confTotalQueryLimit)
	
	if !viper.IsSet(confTotalQueryLimit) {
		totalQueryLimit = 10000
	}
	return totalQueryLimit
}


func GetInternalQueryLimit() int {
	internalQueryLimit := viper.GetInt(confInternalQueryLimit)
	
	if !viper.IsSet(confInternalQueryLimit) {
		internalQueryLimit = 1000
	}
	return internalQueryLimit
}


func GetMaxBatchUpdateSize() int {
	maxBatchUpdateSize := viper.GetInt(confMaxBatchSize)
	
	if !viper.IsSet(confMaxBatchSize) {
		maxBatchUpdateSize = 500
	}
	return maxBatchUpdateSize
}



func GetPvtdataStorePurgeInterval() uint64 {
	purgeInterval := viper.GetInt("ledger.pvtdataStore.purgeInterval")
	if purgeInterval <= 0 {
		purgeInterval = 100
	}
	return uint64(purgeInterval)
}



func GetPvtdataStoreCollElgProcMaxDbBatchSize() int {
	collElgProcMaxDbBatchSize := viper.GetInt(confCollElgProcMaxDbBatchSize.Name)
	if collElgProcMaxDbBatchSize <= 0 {
		collElgProcMaxDbBatchSize = confCollElgProcMaxDbBatchSize.DefaultVal
	}
	return collElgProcMaxDbBatchSize
}



func GetPvtdataStoreCollElgProcDbBatchesInterval() int {
	collElgProcDbBatchesInterval := viper.GetInt(confCollElgProcDbBatchesInterval.Name)
	if collElgProcDbBatchesInterval <= 0 {
		collElgProcDbBatchesInterval = confCollElgProcDbBatchesInterval.DefaultVal
	}
	return collElgProcDbBatchesInterval
}


func IsHistoryDBEnabled() bool {
	return viper.GetBool(confEnableHistoryDatabase)
}



func IsQueryReadsHashingEnabled() bool {
	return true
}




func GetMaxDegreeQueryReadsHashing() uint32 {
	return 50
}


func IsAutoWarmIndexesEnabled() bool {
	
	if viper.IsSet(confAutoWarmIndexes) {
		return viper.GetBool(confAutoWarmIndexes)
	}
	return true

}


func GetWarmIndexesAfterNBlocks() int {
	warmAfterNBlocks := viper.GetInt(confWarmIndexesAfterNBlocks)
	
	if !viper.IsSet(confWarmIndexesAfterNBlocks) {
		warmAfterNBlocks = 1
	}
	return warmAfterNBlocks
}

type conf struct {
	Name       string
	DefaultVal int
}
