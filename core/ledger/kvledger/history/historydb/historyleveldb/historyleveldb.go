/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package historyleveldb

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/history/historydb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	putils "github.com/mcc-github/blockchain/protos/utils"
)

var logger = flogging.MustGetLogger("historyleveldb")

var savePointKey = []byte{0x00}
var emptyValue = []byte{}


type HistoryDBProvider struct {
	dbProvider *leveldbhelper.Provider
}


func NewHistoryDBProvider() *HistoryDBProvider {
	dbPath := ledgerconfig.GetHistoryLevelDBPath()
	logger.Debugf("constructing HistoryDBProvider dbPath=%s", dbPath)
	dbProvider := leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: dbPath})
	return &HistoryDBProvider{dbProvider}
}


func (provider *HistoryDBProvider) GetDBHandle(dbName string) (historydb.HistoryDB, error) {
	return newHistoryDB(provider.dbProvider.GetDBHandle(dbName), dbName), nil
}


func (provider *HistoryDBProvider) Close() {
	provider.dbProvider.Close()
}


type historyDB struct {
	db     *leveldbhelper.DBHandle
	dbName string
}


func newHistoryDB(db *leveldbhelper.DBHandle, dbName string) *historyDB {
	return &historyDB{db, dbName}
}


func (historyDB *historyDB) Open() error {
	
	return nil
}


func (historyDB *historyDB) Close() {
	
}


func (historyDB *historyDB) Commit(block *common.Block) error {

	blockNo := block.Header.Number
	
	var tranNo uint64

	dbBatch := leveldbhelper.NewUpdateBatch()

	logger.Debugf("Channel [%s]: Updating history database for blockNo [%v] with [%d] transactions",
		historyDB.dbName, blockNo, len(block.Data.Data))

	
	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

	
	for _, envBytes := range block.Data.Data {

		
		if txsFilter.IsInvalid(int(tranNo)) {
			logger.Debugf("Channel [%s]: Skipping history write for invalid transaction number %d",
				historyDB.dbName, tranNo)
			tranNo++
			continue
		}

		env, err := putils.GetEnvelopeFromBlock(envBytes)
		if err != nil {
			return err
		}

		payload, err := putils.GetPayload(env)
		if err != nil {
			return err
		}

		chdr, err := putils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			return err
		}

		if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

			
			respPayload, err := putils.GetActionFromEnvelope(envBytes)
			if err != nil {
				return err
			}

			
			txRWSet := &rwsetutil.TxRwSet{}

			
			
			if err = txRWSet.FromProtoBytes(respPayload.Results); err != nil {
				return err
			}
			
			
			for _, nsRWSet := range txRWSet.NsRwSets {
				ns := nsRWSet.NameSpace

				for _, kvWrite := range nsRWSet.KvRwSet.Writes {
					writeKey := kvWrite.Key

					
					compositeHistoryKey := historydb.ConstructCompositeHistoryKey(ns, writeKey, blockNo, tranNo)

					
					dbBatch.Put(compositeHistoryKey, emptyValue)
				}
			}

		} else {
			logger.Debugf("Skipping transaction [%d] since it is not an endorsement transaction\n", tranNo)
		}
		tranNo++
	}

	
	height := version.NewHeight(blockNo, tranNo)
	dbBatch.Put(savePointKey, height.ToBytes())

	
	
	if err := historyDB.db.WriteBatch(dbBatch, true); err != nil {
		return err
	}

	logger.Debugf("Channel [%s]: Updates committed to history database for blockNo [%v]", historyDB.dbName, blockNo)
	return nil
}


func (historyDB *historyDB) NewHistoryQueryExecutor(blockStore blkstorage.BlockStore) (ledger.HistoryQueryExecutor, error) {
	return &LevelHistoryDBQueryExecutor{historyDB, blockStore}, nil
}


func (historyDB *historyDB) GetLastSavepoint() (*version.Height, error) {
	versionBytes, err := historyDB.db.Get(savePointKey)
	if err != nil || versionBytes == nil {
		return nil, err
	}
	height, _ := version.NewHeightFromBytes(versionBytes)
	return height, nil
}


func (historyDB *historyDB) ShouldRecover(lastAvailableBlock uint64) (bool, uint64, error) {
	if !ledgerconfig.IsHistoryDBEnabled() {
		return false, 0, nil
	}
	savepoint, err := historyDB.GetLastSavepoint()
	if err != nil {
		return false, 0, err
	}
	if savepoint == nil {
		return true, 0, nil
	}
	return savepoint.BlockNum != lastAvailableBlock, savepoint.BlockNum + 1, nil
}


func (historyDB *historyDB) CommitLostBlock(blockAndPvtdata *ledger.BlockAndPvtData) error {
	block := blockAndPvtdata.Block
	if err := historyDB.Commit(block); err != nil {
		return err
	}
	return nil
}
