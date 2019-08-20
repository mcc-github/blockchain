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
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	protoutil "github.com/mcc-github/blockchain/protoutil"
)

var logger historydbLogger = flogging.MustGetLogger("historyleveldb")

var savePointKey = []byte{0x00}
var emptyValue = []byte{}




type historydbLogger interface {
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
}


type HistoryDBProvider struct {
	dbProvider *leveldbhelper.Provider
}


func NewHistoryDBProvider(dbPath string) *HistoryDBProvider {
	logger.Debugf("constructing HistoryDBProvider dbPath=%s", dbPath)
	dbProvider := leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: dbPath})
	return &HistoryDBProvider{
		dbProvider: dbProvider,
	}
}


func (p *HistoryDBProvider) GetDBHandle(dbName string) (historydb.HistoryDB, error) {
	return newHistoryDB(
			p.dbProvider.GetDBHandle(dbName),
			dbName,
		),
		nil
}


func (p *HistoryDBProvider) Close() {
	p.dbProvider.Close()
}


type historyDB struct {
	db     *leveldbhelper.DBHandle
	dbName string
}


func newHistoryDB(db *leveldbhelper.DBHandle, dbName string) *historyDB {
	return &historyDB{
		db:     db,
		dbName: dbName,
	}
}


func (h *historyDB) Commit(block *common.Block) error {

	blockNo := block.Header.Number
	
	var tranNo uint64

	dbBatch := leveldbhelper.NewUpdateBatch()

	logger.Debugf("Channel [%s]: Updating history database for blockNo [%v] with [%d] transactions",
		h.dbName, blockNo, len(block.Data.Data))

	
	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

	
	for _, envBytes := range block.Data.Data {

		
		if txsFilter.IsInvalid(int(tranNo)) {
			logger.Debugf("Channel [%s]: Skipping history write for invalid transaction number %d",
				h.dbName, tranNo)
			tranNo++
			continue
		}

		env, err := protoutil.GetEnvelopeFromBlock(envBytes)
		if err != nil {
			return err
		}

		payload, err := protoutil.UnmarshalPayload(env.Payload)
		if err != nil {
			return err
		}

		chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			return err
		}

		if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

			
			respPayload, err := protoutil.GetActionFromEnvelope(envBytes)
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

	
	
	if err := h.db.WriteBatch(dbBatch, true); err != nil {
		return err
	}

	logger.Debugf("Channel [%s]: Updates committed to history database for blockNo [%v]", h.dbName, blockNo)
	return nil
}


func (h *historyDB) NewHistoryQueryExecutor(blockStore blkstorage.BlockStore) (ledger.HistoryQueryExecutor, error) {
	return &LevelHistoryDBQueryExecutor{h, blockStore}, nil
}


func (h *historyDB) GetLastSavepoint() (*version.Height, error) {
	versionBytes, err := h.db.Get(savePointKey)
	if err != nil || versionBytes == nil {
		return nil, err
	}
	height, _, err := version.NewHeightFromBytes(versionBytes)
	if err != nil {
		return nil, err
	}
	return height, nil
}


func (h *historyDB) ShouldRecover(lastAvailableBlock uint64) (bool, uint64, error) {
	savepoint, err := h.GetLastSavepoint()
	if err != nil {
		return false, 0, err
	}
	if savepoint == nil {
		return true, 0, nil
	}
	return savepoint.BlockNum != lastAvailableBlock, savepoint.BlockNum + 1, nil
}


func (historyDB *historyDB) Name() string {
	return "history"
}


func (h *historyDB) CommitLostBlock(blockAndPvtdata *ledger.BlockAndPvtData) error {
	block := blockAndPvtdata.Block

	
	if block.Header.Number%1000 == 0 {
		logger.Infof("Recommitting block [%d] to history database", block.Header.Number)
	} else {
		logger.Debugf("Recommitting block [%d] to history database", block.Header.Number)
	}

	if err := h.Commit(block); err != nil {
		return err
	}
	return nil
}
