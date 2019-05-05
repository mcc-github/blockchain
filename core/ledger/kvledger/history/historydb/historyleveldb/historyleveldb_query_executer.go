/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package historyleveldb

import (
	"bytes"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/history/historydb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	protoutil "github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)


type LevelHistoryDBQueryExecutor struct {
	historyDB  *historyDB
	blockStore blkstorage.BlockStore
}


func (q *LevelHistoryDBQueryExecutor) GetHistoryForKey(namespace string, key string) (commonledger.ResultsIterator, error) {
	var compositeStartKey []byte
	var compositeEndKey []byte
	compositeStartKey = historydb.ConstructPartialCompositeHistoryKey(namespace, key, false)
	compositeEndKey = historydb.ConstructPartialCompositeHistoryKey(namespace, key, true)

	
	dbItr := q.historyDB.db.GetIterator(compositeStartKey, compositeEndKey)
	return newHistoryScanner(compositeStartKey, namespace, key, dbItr, q.blockStore), nil
}


type historyScanner struct {
	compositePartialKey []byte 
	namespace           string
	key                 string
	dbItr               iterator.Iterator
	blockStore          blkstorage.BlockStore
}

func newHistoryScanner(compositePartialKey []byte, namespace string, key string,
	dbItr iterator.Iterator, blockStore blkstorage.BlockStore) *historyScanner {
	return &historyScanner{compositePartialKey, namespace, key, dbItr, blockStore}
}

func (scanner *historyScanner) Next() (commonledger.QueryResult, error) {
	for {
		if !scanner.dbItr.Next() {
			return nil, nil
		}
		historyKey := scanner.dbItr.Key() 

		
		_, blockNumTranNumBytes := historydb.SplitCompositeHistoryKey(historyKey, scanner.compositePartialKey)

		
		
		
		
		
		
		
		
		
		
		
		
		
		if bytes.Contains(blockNumTranNumBytes[:len(blockNumTranNumBytes)-1], historydb.CompositeKeySep) {
			logger.Debugf("Some other key [%#v] found in the range while scanning history for key [%#v]. Skipping...",
				historyKey, scanner.key)
			continue
		}
		blockNum, bytesConsumed := util.DecodeOrderPreservingVarUint64(blockNumTranNumBytes[0:])
		tranNum, _ := util.DecodeOrderPreservingVarUint64(blockNumTranNumBytes[bytesConsumed:])
		logger.Debugf("Found history record for namespace:%s key:%s at blockNumTranNum %v:%v\n",
			scanner.namespace, scanner.key, blockNum, tranNum)

		
		tranEnvelope, err := scanner.blockStore.RetrieveTxByBlockNumTranNum(blockNum, tranNum)
		if err != nil {
			return nil, err
		}

		
		queryResult, err := getKeyModificationFromTran(tranEnvelope, scanner.namespace, scanner.key)
		if err != nil {
			return nil, err
		}
		logger.Debugf("Found historic key value for namespace:%s key:%s from transaction %s\n",
			scanner.namespace, scanner.key, queryResult.(*queryresult.KeyModification).TxId)
		return queryResult, nil
	}
}

func (scanner *historyScanner) Close() {
	scanner.dbItr.Release()
}


func getKeyModificationFromTran(tranEnvelope *common.Envelope, namespace string, key string) (commonledger.QueryResult, error) {
	logger.Debugf("Entering getKeyModificationFromTran()\n", namespace, key)

	
	payload, err := protoutil.GetPayload(tranEnvelope)
	if err != nil {
		return nil, err
	}

	tx, err := protoutil.GetTransaction(payload.Data)
	if err != nil {
		return nil, err
	}

	_, respPayload, err := protoutil.GetPayloads(tx.Actions[0])
	if err != nil {
		return nil, err
	}

	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, err
	}

	txID := chdr.TxId
	timestamp := chdr.Timestamp

	txRWSet := &rwsetutil.TxRwSet{}

	
	
	if err = txRWSet.FromProtoBytes(respPayload.Results); err != nil {
		return nil, err
	}

	
	for _, nsRWSet := range txRWSet.NsRwSets {
		if nsRWSet.NameSpace == namespace {
			
			for _, kvWrite := range nsRWSet.KvRwSet.Writes {
				if kvWrite.Key == key {
					return &queryresult.KeyModification{TxId: txID, Value: kvWrite.Value,
						Timestamp: timestamp, IsDelete: kvWrite.IsDelete}, nil
				}
			} 
			return nil, errors.New("key not found in namespace's writeset")
		} 
	} 
	return nil, errors.New("namespace not found in transaction's ReadWriteSets")

}
