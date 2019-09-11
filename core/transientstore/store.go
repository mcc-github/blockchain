/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transientstore

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/mcc-github/blockchain-protos-go/transientstore"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var logger = flogging.MustGetLogger("transientstore")

var (
	emptyValue = []byte{}
	nilByte    = byte('\x00')
	
	ErrStoreEmpty = errors.New("Transient store is empty")
)






type StoreProvider interface {
	OpenStore(ledgerID string) (Store, error)
	Close()
}


type RWSetScanner interface {
	
	
	
	Next() (*EndorserPvtSimulationResults, error)
	
	Close()
}




type Store interface {
	
	
	Persist(txid string, blockHeight uint64, privateSimulationResultsWithConfig *transientstore.TxPvtReadWriteSetWithConfigInfo) error
	
	
	GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (RWSetScanner, error)
	
	
	PurgeByTxids(txids []string) error
	
	
	
	
	
	
	PurgeByHeight(maxBlockNumToRetain uint64) error
	
	GetMinTransientBlkHt() (uint64, error)
	Shutdown()
}


type EndorserPvtSimulationResults struct {
	ReceivedAtBlockHeight          uint64
	PvtSimulationResultsWithConfig *transientstore.TxPvtReadWriteSetWithConfigInfo
}








type storeProvider struct {
	dbProvider *leveldbhelper.Provider
}


type store struct {
	db       *leveldbhelper.DBHandle
	ledgerID string
}


type RwsetScanner struct {
	txid   string
	dbItr  iterator.Iterator
	filter ledger.PvtNsCollFilter
}


func NewStoreProvider(path string) (StoreProvider, error) {
	dbProvider, err := leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: path})
	if err != nil {
		return nil, err
	}
	return &storeProvider{dbProvider: dbProvider}, nil
}


func (provider *storeProvider) OpenStore(ledgerID string) (Store, error) {
	dbHandle := provider.dbProvider.GetDBHandle(ledgerID)
	return &store{db: dbHandle, ledgerID: ledgerID}, nil
}


func (provider *storeProvider) Close() {
	provider.dbProvider.Close()
}



func (s *store) Persist(txid string, blockHeight uint64,
	privateSimulationResultsWithConfig *transientstore.TxPvtReadWriteSetWithConfigInfo) error {

	logger.Debugf("Persisting private data to transient store for txid [%s] at block height [%d]", txid, blockHeight)

	dbBatch := leveldbhelper.NewUpdateBatch()

	
	
	
	uuid := util.GenerateUUID()
	compositeKeyPvtRWSet := createCompositeKeyForPvtRWSet(txid, uuid, blockHeight)
	privateSimulationResultsWithConfigBytes, err := proto.Marshal(privateSimulationResultsWithConfig)
	if err != nil {
		return err
	}

	
	
	
	
	
	value := append([]byte{nilByte}, privateSimulationResultsWithConfigBytes...)
	dbBatch.Put(compositeKeyPvtRWSet, value)

	

	
	
	
	
	
	compositeKeyPurgeIndexByHeight := createCompositeKeyForPurgeIndexByHeight(blockHeight, txid, uuid)
	dbBatch.Put(compositeKeyPurgeIndexByHeight, emptyValue)

	
	
	
	
	
	
	
	
	
	
	compositeKeyPurgeIndexByTxid := createCompositeKeyForPurgeIndexByTxid(txid, uuid, blockHeight)
	dbBatch.Put(compositeKeyPurgeIndexByTxid, emptyValue)

	return s.db.WriteBatch(dbBatch, true)
}



func (s *store) GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (RWSetScanner, error) {

	logger.Debugf("Getting private data from transient store for transaction %s", txid)

	
	startKey := createTxidRangeStartKey(txid)
	endKey := createTxidRangeEndKey(txid)

	iter := s.db.GetIterator(startKey, endKey)
	return &RwsetScanner{txid, iter, filter}, nil
}




func (s *store) PurgeByTxids(txids []string) error {

	logger.Debug("Purging private data from transient store for committed txids")

	dbBatch := leveldbhelper.NewUpdateBatch()

	for _, txid := range txids {
		
		startKey := createPurgeIndexByTxidRangeStartKey(txid)
		endKey := createPurgeIndexByTxidRangeEndKey(txid)

		iter := s.db.GetIterator(startKey, endKey)

		
		
		for iter.Next() {
			

			
			compositeKeyPurgeIndexByTxid := iter.Key()
			
			
			uuid, blockHeight, err := splitCompositeKeyOfPurgeIndexByTxid(compositeKeyPurgeIndexByTxid)
			if err != nil {
				return err
			}
			compositeKeyPvtRWSet := createCompositeKeyForPvtRWSet(txid, uuid, blockHeight)
			dbBatch.Delete(compositeKeyPvtRWSet)

			
			compositeKeyPurgeIndexByHeight := createCompositeKeyForPurgeIndexByHeight(blockHeight, txid, uuid)
			dbBatch.Delete(compositeKeyPurgeIndexByHeight)

			
			dbBatch.Delete(compositeKeyPurgeIndexByTxid)
		}
		iter.Release()
	}
	
	
	return s.db.WriteBatch(dbBatch, true)
}







func (s *store) PurgeByHeight(maxBlockNumToRetain uint64) error {

	logger.Debugf("Purging orphaned private data from transient store received prior to block [%d]", maxBlockNumToRetain)

	
	startKey := createPurgeIndexByHeightRangeStartKey(0)
	endKey := createPurgeIndexByHeightRangeEndKey(maxBlockNumToRetain - 1)
	iter := s.db.GetIterator(startKey, endKey)

	dbBatch := leveldbhelper.NewUpdateBatch()

	
	
	for iter.Next() {
		

		
		compositeKeyPurgeIndexByHeight := iter.Key()
		txid, uuid, blockHeight, err := splitCompositeKeyOfPurgeIndexByHeight(compositeKeyPurgeIndexByHeight)
		if err != nil {
			return err
		}
		logger.Debugf("Purging from transient store private data simulated at block [%d]: txid [%s] uuid [%s]", blockHeight, txid, uuid)

		compositeKeyPvtRWSet := createCompositeKeyForPvtRWSet(txid, uuid, blockHeight)
		dbBatch.Delete(compositeKeyPvtRWSet)

		
		compositeKeyPurgeIndexByTxid := createCompositeKeyForPurgeIndexByTxid(txid, uuid, blockHeight)
		dbBatch.Delete(compositeKeyPurgeIndexByTxid)

		
		dbBatch.Delete(compositeKeyPurgeIndexByHeight)
	}
	iter.Release()

	return s.db.WriteBatch(dbBatch, true)
}


func (s *store) GetMinTransientBlkHt() (uint64, error) {
	
	
	
	
	startKey := createPurgeIndexByHeightRangeStartKey(0)
	iter := s.db.GetIterator(startKey, nil)
	defer iter.Release()
	
	if iter.Next() {
		dbKey := iter.Key()
		_, _, blockHeight, err := splitCompositeKeyOfPurgeIndexByHeight(dbKey)
		return blockHeight, err
	}
	
	
	
	return 0, ErrStoreEmpty
}

func (s *store) Shutdown() {
	
}



func (scanner *RwsetScanner) Next() (*EndorserPvtSimulationResults, error) {
	if !scanner.dbItr.Next() {
		return nil, nil
	}
	dbKey := scanner.dbItr.Key()
	dbVal := scanner.dbItr.Value()
	_, blockHeight, err := splitCompositeKeyOfPvtRWSet(dbKey)
	if err != nil {
		return nil, err
	}

	txPvtRWSet := &rwset.TxPvtReadWriteSet{}
	filteredTxPvtRWSet := &rwset.TxPvtReadWriteSet{}
	txPvtRWSetWithConfig := &transientstore.TxPvtReadWriteSetWithConfigInfo{}

	if dbVal[0] == nilByte {
		
		if err := proto.Unmarshal(dbVal[1:], txPvtRWSetWithConfig); err != nil {
			return nil, err
		}

		filteredTxPvtRWSet = trimPvtWSet(txPvtRWSetWithConfig.GetPvtRwset(), scanner.filter)
		configs, err := trimPvtCollectionConfigs(txPvtRWSetWithConfig.CollectionConfigs, scanner.filter)
		if err != nil {
			return nil, err
		}
		txPvtRWSetWithConfig.CollectionConfigs = configs
	} else {
		
		if err := proto.Unmarshal(dbVal, txPvtRWSet); err != nil {
			return nil, err
		}
		filteredTxPvtRWSet = trimPvtWSet(txPvtRWSet, scanner.filter)
	}

	txPvtRWSetWithConfig.PvtRwset = filteredTxPvtRWSet

	return &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          blockHeight,
		PvtSimulationResultsWithConfig: txPvtRWSetWithConfig,
	}, nil
}


func (scanner *RwsetScanner) Close() {
	scanner.dbItr.Release()
}
