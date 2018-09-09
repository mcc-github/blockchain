/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lockbasedtxmgr

import (
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/storageutil"
)


type lockBasedQueryExecutor struct {
	helper *queryHelper
	txid   string
}

func newQueryExecutor(txmgr *LockBasedTxMgr, txid string) *lockBasedQueryExecutor {
	helper := newQueryHelper(txmgr, nil)
	logger.Debugf("constructing new query executor txid = [%s]", txid)
	return &lockBasedQueryExecutor{helper, txid}
}


func (q *lockBasedQueryExecutor) GetState(ns string, key string) (val []byte, err error) {
	val, _, err = q.helper.getState(ns, key)
	return
}


func (q *lockBasedQueryExecutor) GetStateMetadata(namespace, key string) (map[string][]byte, error) {
	_, metadataBytes, err := q.helper.getState(namespace, key)
	if err != nil {
		return nil, err
	}
	metadata, err := storageutil.DeserializeMetadata(metadataBytes)
	if err != nil {
		return nil, err
	}
	return metadata, err
}


func (q *lockBasedQueryExecutor) GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error) {
	return q.helper.getStateMultipleKeys(namespace, keys)
}





func (q *lockBasedQueryExecutor) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (commonledger.ResultsIterator, error) {
	return q.helper.getStateRangeScanIterator(namespace, startKey, endKey)
}






func (q *lockBasedQueryExecutor) GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (ledger.QueryResultsIterator, error) {
	return q.helper.getStateRangeScanIteratorWithMetadata(namespace, startKey, endKey, metadata)
}


func (q *lockBasedQueryExecutor) ExecuteQuery(namespace, query string) (commonledger.ResultsIterator, error) {
	return q.helper.executeQuery(namespace, query)
}


func (q *lockBasedQueryExecutor) ExecuteQueryWithMetadata(namespace, query string, metadata map[string]interface{}) (ledger.QueryResultsIterator, error) {
	return q.helper.executeQueryWithMetadata(namespace, query, metadata)
}


func (q *lockBasedQueryExecutor) GetPrivateData(namespace, collection, key string) ([]byte, error) {
	return q.helper.getPrivateData(namespace, collection, key)
}


func (q *lockBasedQueryExecutor) GetPrivateDataMetadata(namespace, collection, key string) (map[string][]byte, error) {
	_, metadataBytes, err := q.helper.getPrivateDataValueHash(namespace, collection, key)
	if err != nil {
		return nil, err
	}
	metadata, err := storageutil.DeserializeMetadata(metadataBytes)
	if err != nil {
		return nil, err
	}
	return metadata, err
}


func (q *lockBasedQueryExecutor) GetPrivateDataMultipleKeys(namespace, collection string, keys []string) ([][]byte, error) {
	return q.helper.getPrivateDataMultipleKeys(namespace, collection, keys)
}


func (q *lockBasedQueryExecutor) GetPrivateDataRangeScanIterator(namespace, collection, startKey, endKey string) (commonledger.ResultsIterator, error) {
	return q.helper.getPrivateDataRangeScanIterator(namespace, collection, startKey, endKey)
}


func (q *lockBasedQueryExecutor) ExecuteQueryOnPrivateData(namespace, collection, query string) (commonledger.ResultsIterator, error) {
	return q.helper.executeQueryOnPrivateData(namespace, collection, query)
}


func (q *lockBasedQueryExecutor) Done() {
	logger.Debugf("Done with transaction simulation / query execution [%s]", q.txid)
	q.helper.done()
}
