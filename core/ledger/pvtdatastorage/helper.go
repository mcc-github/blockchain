/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"math"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/willf/bitset"
)

func prepareStoreEntries(blockNum uint64, pvtData []*ledger.TxPvtData, btlPolicy pvtdatapolicy.BTLPolicy,
	missingPvtData ledger.TxMissingPvtDataMap) (*storeEntries, error) {
	dataEntries := prepareDataEntries(blockNum, pvtData)

	missingDataEntries := prepareMissingDataEntries(blockNum, missingPvtData)

	expiryEntries, err := prepareExpiryEntries(blockNum, dataEntries, missingDataEntries, btlPolicy)
	if err != nil {
		return nil, err
	}

	return &storeEntries{
		dataEntries:        dataEntries,
		expiryEntries:      expiryEntries,
		missingDataEntries: missingDataEntries}, nil
}

func prepareDataEntries(blockNum uint64, pvtData []*ledger.TxPvtData) []*dataEntry {
	var dataEntries []*dataEntry
	for _, txPvtdata := range pvtData {
		for _, nsPvtdata := range txPvtdata.WriteSet.NsPvtRwset {
			for _, collPvtdata := range nsPvtdata.CollectionPvtRwset {
				txnum := txPvtdata.SeqInBlock
				ns := nsPvtdata.Namespace
				coll := collPvtdata.CollectionName
				dataKey := &dataKey{nsCollBlk{ns, coll, blockNum}, txnum}
				dataEntries = append(dataEntries, &dataEntry{key: dataKey, value: collPvtdata})
			}
		}
	}
	return dataEntries
}

func prepareMissingDataEntries(committingBlk uint64, missingPvtData ledger.TxMissingPvtDataMap) map[missingDataKey]*bitset.BitSet {
	missingDataEntries := make(map[missingDataKey]*bitset.BitSet)

	for txNum, missingData := range missingPvtData {
		for _, nsColl := range missingData {
			key := missingDataKey{nsCollBlk{nsColl.Namespace, nsColl.Collection, committingBlk},
				nsColl.IsEligible}

			if _, ok := missingDataEntries[key]; !ok {
				missingDataEntries[key] = &bitset.BitSet{}
			}
			bitmap := missingDataEntries[key]

			bitmap.Set(uint(txNum))
		}
	}

	return missingDataEntries
}



func prepareExpiryEntries(committingBlk uint64, dataEntries []*dataEntry, missingDataEntries map[missingDataKey]*bitset.BitSet,
	btlPolicy pvtdatapolicy.BTLPolicy) ([]*expiryEntry, error) {

	var expiryEntries []*expiryEntry
	mapByExpiringBlk := make(map[uint64]*ExpiryData)

	
	for _, dataEntry := range dataEntries {
		prepareExpiryEntriesForPresentData(mapByExpiringBlk, dataEntry.key, btlPolicy)
	}

	
	for missingDataKey := range missingDataEntries {
		prepareExpiryEntriesForMissingData(mapByExpiringBlk, &missingDataKey, btlPolicy)
	}

	for expiryBlk, expiryData := range mapByExpiringBlk {
		expiryKey := &expiryKey{expiringBlk: expiryBlk, committingBlk: committingBlk}
		expiryEntries = append(expiryEntries, &expiryEntry{key: expiryKey, value: expiryData})
	}

	return expiryEntries, nil
}


func prepareExpiryEntriesForPresentData(mapByExpiringBlk map[uint64]*ExpiryData, dataKey *dataKey, btlPolicy pvtdatapolicy.BTLPolicy) error {
	expiringBlk, err := btlPolicy.GetExpiringBlock(dataKey.ns, dataKey.coll, dataKey.blkNum)
	if err != nil {
		return err
	}
	if neverExpires(expiringBlk) {
		return nil
	}

	expiryData := getOrCreateExpiryData(mapByExpiringBlk, expiringBlk)

	expiryData.addPresentData(dataKey.ns, dataKey.coll, dataKey.txNum)
	return nil
}


func prepareExpiryEntriesForMissingData(mapByExpiringBlk map[uint64]*ExpiryData, missingKey *missingDataKey, btlPolicy pvtdatapolicy.BTLPolicy) error {
	expiringBlk, err := btlPolicy.GetExpiringBlock(missingKey.ns, missingKey.coll, missingKey.blkNum)
	if err != nil {
		return err
	}
	if neverExpires(expiringBlk) {
		return nil
	}

	expiryData := getOrCreateExpiryData(mapByExpiringBlk, expiringBlk)

	expiryData.addMissingData(missingKey.ns, missingKey.coll)
	return nil
}

func getOrCreateExpiryData(mapByExpiringBlk map[uint64]*ExpiryData, expiringBlk uint64) *ExpiryData {
	expiryData, ok := mapByExpiringBlk[expiringBlk]
	if !ok {
		expiryData = newExpiryData()
		mapByExpiringBlk[expiringBlk] = expiryData
	}
	return expiryData
}


func deriveKeys(expiryEntry *expiryEntry) (dataKeys []*dataKey, missingDataKeys []*missingDataKey) {
	for ns, colls := range expiryEntry.value.Map {
		
		for coll, txNums := range colls.Map {
			for _, txNum := range txNums.List {
				dataKeys = append(dataKeys,
					&dataKey{nsCollBlk{ns, coll, expiryEntry.key.committingBlk}, txNum})
			}
		}
		
		for coll := range colls.MissingDataMap {
			
			missingDataKeys = append(missingDataKeys,
				&missingDataKey{nsCollBlk{ns, coll, expiryEntry.key.committingBlk}, true})
			missingDataKeys = append(missingDataKeys,
				&missingDataKey{nsCollBlk{ns, coll, expiryEntry.key.committingBlk}, false})

		}
	}
	return
}

func passesFilter(dataKey *dataKey, filter ledger.PvtNsCollFilter) bool {
	return filter == nil || filter.Has(dataKey.ns, dataKey.coll)
}

func isExpired(key nsCollBlk, btl pvtdatapolicy.BTLPolicy, latestBlkNum uint64) (bool, error) {
	expiringBlk, err := btl.GetExpiringBlock(key.ns, key.coll, key.blkNum)
	if err != nil {
		return false, err
	}

	return latestBlkNum >= expiringBlk, nil
}

func neverExpires(expiringBlkNum uint64) bool {
	return expiringBlkNum == math.MaxUint64
}

type txPvtdataAssembler struct {
	blockNum, txNum uint64
	txWset          *rwset.TxPvtReadWriteSet
	currentNsWSet   *rwset.NsPvtReadWriteSet
	firstCall       bool
}

func newTxPvtdataAssembler(blockNum, txNum uint64) *txPvtdataAssembler {
	return &txPvtdataAssembler{blockNum, txNum, &rwset.TxPvtReadWriteSet{}, nil, true}
}

func (a *txPvtdataAssembler) add(ns string, collPvtWset *rwset.CollectionPvtReadWriteSet) {
	
	if a.firstCall {
		a.currentNsWSet = &rwset.NsPvtReadWriteSet{Namespace: ns}
		a.firstCall = false
	}

	
	if a.currentNsWSet.Namespace != ns {
		a.txWset.NsPvtRwset = append(a.txWset.NsPvtRwset, a.currentNsWSet)
		a.currentNsWSet = &rwset.NsPvtReadWriteSet{Namespace: ns}
	}
	
	a.currentNsWSet.CollectionPvtRwset = append(a.currentNsWSet.CollectionPvtRwset, collPvtWset)
}

func (a *txPvtdataAssembler) done() {
	if a.currentNsWSet != nil {
		a.txWset.NsPvtRwset = append(a.txWset.NsPvtRwset, a.currentNsWSet)
	}
	a.currentNsWSet = nil
}

func (a *txPvtdataAssembler) getTxPvtdata() *ledger.TxPvtData {
	a.done()
	return &ledger.TxPvtData{SeqInBlock: a.txNum, WriteSet: a.txWset}
}
