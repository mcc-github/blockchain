/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package queryutil

import (
	"fmt"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
)

type itrCombiner struct {
	namespace string
	holders   []*itrHolder
}

func newItrCombiner(namespace string, baseIterators []statedb.ResultsIterator) (*itrCombiner, error) {
	var holders []*itrHolder
	for _, itr := range baseIterators {
		res, err := itr.Next()
		if err != nil {
			for _, holder := range holders {
				holder.itr.Close()
			}
			return nil, err
		}
		if res != nil {
			holders = append(holders, &itrHolder{itr, res.(*statedb.VersionedKV)})
		}
	}
	return &itrCombiner{namespace, holders}, nil
}





func (combiner *itrCombiner) Next() (commonledger.QueryResult, error) {
	logger.Debugf("Iterators position at beginning: %s", combiner.holders)
	if len(combiner.holders) == 0 {
		return nil, nil
	}
	smallestHolderIndex := 0
	for i := 1; i < len(combiner.holders); i++ {
		smallestKey, holderKey := combiner.keyAt(smallestHolderIndex), combiner.keyAt(i)
		switch {
		case holderKey == smallestKey: 
			
			
			removed, err := combiner.moveItrAndRemoveIfExhausted(i)
			if err != nil {
				return nil, err
			}
			if removed { 
				
				i--
			}
		case holderKey < smallestKey:
			smallestHolderIndex = i
		default:
			
		}
	}
	kv := combiner.kvAt(smallestHolderIndex)
	combiner.moveItrAndRemoveIfExhausted(smallestHolderIndex)
	if kv.IsDelete() {
		return combiner.Next()
	}
	logger.Debugf("Key [%s] selected from iterator at index [%d]", kv.Key, smallestHolderIndex)
	logger.Debugf("Iterators position at end: %s", combiner.holders)
	return &queryresult.KV{Namespace: combiner.namespace, Key: kv.Key, Value: kv.Value}, nil
}



func (combiner *itrCombiner) moveItrAndRemoveIfExhausted(i int) (removed bool, err error) {
	holder := combiner.holders[i]
	exhausted, err := holder.moveToNext()
	if err != nil {
		return false, err
	}
	if exhausted {
		combiner.holders[i].itr.Close()
		combiner.holders = append(combiner.holders[:i], combiner.holders[i+1:]...)

	}
	return exhausted, nil
}


func (combiner *itrCombiner) kvAt(i int) *statedb.VersionedKV {
	return combiner.holders[i].kv
}


func (combiner *itrCombiner) keyAt(i int) string {
	return combiner.kvAt(i).Key
}


func (combiner *itrCombiner) Close() {
	for _, holder := range combiner.holders {
		holder.itr.Close()
	}
}


type itrHolder struct {
	itr statedb.ResultsIterator
	kv  *statedb.VersionedKV
}


func (holder *itrHolder) moveToNext() (exhausted bool, err error) {
	var res statedb.QueryResult
	if res, err = holder.itr.Next(); err != nil {
		return false, err
	}
	if res != nil {
		holder.kv = res.(*statedb.VersionedKV)
	}
	return res == nil, nil
}


func (holder *itrHolder) String() string {
	return fmt.Sprintf("{%s}", holder.kv.Key)
}
