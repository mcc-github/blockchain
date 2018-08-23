/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package queryutil

import (
	"github.com/mcc-github/blockchain/common/flogging"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
)

var logger = flogging.MustGetLogger("util")




type QueryExecuter interface {
	GetState(namespace, key string) (*statedb.VersionedValue, error)
	GetStateRangeScanIterator(namespace, startKey, endKey string) (statedb.ResultsIterator, error)
}




type QECombiner struct {
	QueryExecuters []QueryExecuter 
}


func (c *QECombiner) GetState(namespace string, key string) ([]byte, error) {
	var vv *statedb.VersionedValue
	var val []byte
	var err error
	for _, qe := range c.QueryExecuters {
		if vv, err = qe.GetState(namespace, key); err != nil {
			return nil, err
		}
		if vv != nil {
			if !vv.IsDelete() {
				val = vv.Value
			}
			break
		}
	}
	return val, nil
}


func (c *QECombiner) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (commonledger.ResultsIterator, error) {
	var itrs []statedb.ResultsIterator
	for _, qe := range c.QueryExecuters {
		itr, err := qe.GetStateRangeScanIterator(namespace, startKey, endKey)
		if err != nil {
			for _, itr := range itrs {
				itr.Close()
			}
			return nil, err
		}
		itrs = append(itrs, itr)
	}
	itrCombiner, err := newItrCombiner(namespace, itrs)
	if err != nil {
		return nil, err
	}
	return itrCombiner, nil
}


type UpdateBatchBackedQueryExecuter struct {
	UpdateBatch *statedb.UpdateBatch
}


func (qe *UpdateBatchBackedQueryExecuter) GetState(ns, key string) (*statedb.VersionedValue, error) {
	return qe.UpdateBatch.Get(ns, key), nil
}


func (qe *UpdateBatchBackedQueryExecuter) GetStateRangeScanIterator(namespace, startKey, endKey string) (statedb.ResultsIterator, error) {
	return qe.UpdateBatch.GetRangeScanIterator(namespace, startKey, endKey), nil
}
