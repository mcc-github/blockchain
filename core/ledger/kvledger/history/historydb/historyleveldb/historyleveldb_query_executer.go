/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package historyleveldb

import (
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

		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		blockNum, tranNum, err := decodeBlockNumTranNum(blockNumTranNumBytes)
		if err != nil {
			logger.Warnf("Some other key [%#v] found in the range while scanning history for key [%#v]. Skipping (decoding error: %s)",
				historyKey, scanner.key, err)
			continue
		}

		logger.Debugf("Found history record for namespace:%s key:%s at blockNumTranNum %v:%v\n",
			scanner.namespace, scanner.key, blockNum, tranNum)

		
		tranEnvelope, err := scanner.blockStore.RetrieveTxByBlockNumTranNum(blockNum, tranNum)
		if err == blkstorage.ErrNotFoundInIndex {
			logger.Warnf("Some other clashing key [%#v] found in the range while scanning history for key [%#v]. Skipping (cannot find block:tx)",
				historyKey, scanner.key)
			continue
		}
		if err != nil {
			return nil, err
		}

		
		queryResult, err := getKeyModificationFromTran(tranEnvelope, scanner.namespace, scanner.key)
		if err != nil {
			return nil, err
		}
		if queryResult == nil {
			
			
			
			
			logger.Warnf("Some other key [%#v] found in the range while scanning history for key [%#v]. Skipping (namespace or key not found)",
				historyKey, scanner.key)
			continue
		}
		logger.Debugf("Found historic key value for namespace:%s key:%s from transaction %s",
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
			logger.Debugf("key [%s] not found in namespace [%s]'s writeset", key, namespace)
			return nil, nil
		} 
	} 
	logger.Debugf("namespace [%s] not found in transaction's ReadWriteSets", namespace)
	return nil, nil
}


func decodeBlockNumTranNum(blockNumTranNumBytes []byte) (uint64, uint64, error) {
	blockNum, blockBytesConsumed, err := util.DecodeOrderPreservingVarUint64(blockNumTranNumBytes)
	if err != nil {
		return 0, 0, err
	}

	tranNum, tranBytesConsumed, err := util.DecodeOrderPreservingVarUint64(blockNumTranNumBytes[blockBytesConsumed:])
	if err != nil {
		return 0, 0, err
	}

	
	if blockBytesConsumed+tranBytesConsumed != len(blockNumTranNumBytes) {
		return 0, 0, errors.Errorf("number of decoded bytes (%d) is not equal to the length of blockNumTranNumBytes (%d)",
			blockBytesConsumed+tranBytesConsumed, len(blockNumTranNumBytes))
	}

	return blockNum, tranNum, nil
}
