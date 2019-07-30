/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"bytes"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/ledgerstorage"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protoutil"
)



func constructValidAndInvalidPvtData(blocksPvtData []*ledger.BlockPvtData, blockStore *ledgerstorage.Store) (
	map[uint64][]*ledger.TxPvtData, []*ledger.PvtdataHashMismatch, error,
) {
	
	
	
	
	validPvtData := make(map[uint64][]*ledger.TxPvtData)
	var invalidPvtData []*ledger.PvtdataHashMismatch

	for _, blockPvtData := range blocksPvtData {
		validData, invalidData, err := findValidAndInvalidBlockPvtData(blockPvtData, blockStore)
		if err != nil {
			return nil, nil, err
		}
		if len(validData) > 0 {
			validPvtData[blockPvtData.BlockNum] = validData
		}
		invalidPvtData = append(invalidPvtData, invalidData...)
	} 
	return validPvtData, invalidPvtData, nil
}

func findValidAndInvalidBlockPvtData(blockPvtData *ledger.BlockPvtData, blockStore *ledgerstorage.Store) (
	[]*ledger.TxPvtData, []*ledger.PvtdataHashMismatch, error,
) {
	var validPvtData []*ledger.TxPvtData
	var invalidPvtData []*ledger.PvtdataHashMismatch
	for _, txPvtData := range blockPvtData.WriteSets {
		
		logger.Debugf("Retrieving rwset of blockNum:[%d], txNum:[%d]", blockPvtData.BlockNum, txPvtData.SeqInBlock)
		txRWSet, err := retrieveRwsetForTx(blockPvtData.BlockNum, txPvtData.SeqInBlock, blockStore)
		if err != nil {
			return nil, nil, err
		}

		
		logger.Debugf("Constructing valid and invalid pvtData using rwset of blockNum:[%d], txNum:[%d]",
			blockPvtData.BlockNum, txPvtData.SeqInBlock)
		validData, invalidData := findValidAndInvalidTxPvtData(txPvtData, txRWSet, blockPvtData.BlockNum)

		
		
		if validData != nil {
			validPvtData = append(validPvtData, validData)
		}
		invalidPvtData = append(invalidPvtData, invalidData...)
	} 
	return validPvtData, invalidPvtData, nil
}

func retrieveRwsetForTx(blkNum uint64, txNum uint64, blockStore *ledgerstorage.Store) (*rwsetutil.TxRwSet, error) {
	
	
	txEnvelope, err := blockStore.RetrieveTxByBlockNumTranNum(blkNum, txNum)
	if err != nil {
		return nil, err
	}
	
	responsePayload, err := protoutil.GetActionFromEnvelopeMsg(txEnvelope)
	if err != nil {
		return nil, err
	}
	txRWSet := &rwsetutil.TxRwSet{}
	if err := txRWSet.FromProtoBytes(responsePayload.Results); err != nil {
		return nil, err
	}
	return txRWSet, nil
}

func findValidAndInvalidTxPvtData(txPvtData *ledger.TxPvtData, txRWSet *rwsetutil.TxRwSet, blkNum uint64) (
	*ledger.TxPvtData, []*ledger.PvtdataHashMismatch,
) {
	var invalidPvtData []*ledger.PvtdataHashMismatch
	var toDeleteNsColl []*nsColl
	
	
	for _, nsRwset := range txPvtData.WriteSet.NsPvtRwset {
		txNum := txPvtData.SeqInBlock
		invalidData, invalidNsColl := findInvalidNsPvtData(nsRwset, txRWSet, blkNum, txNum)
		invalidPvtData = append(invalidPvtData, invalidData...)
		toDeleteNsColl = append(toDeleteNsColl, invalidNsColl...)
	}
	for _, nsColl := range toDeleteNsColl {
		removeCollFromTxPvtReadWriteSet(txPvtData.WriteSet, nsColl.ns, nsColl.coll)
	}
	if len(txPvtData.WriteSet.NsPvtRwset) == 0 {
		
		
		return nil, invalidPvtData
	}
	return txPvtData, invalidPvtData
}



func removeCollFromTxPvtReadWriteSet(p *rwset.TxPvtReadWriteSet, ns, coll string) {
	for i := 0; i < len(p.NsPvtRwset); i++ {
		n := p.NsPvtRwset[i]
		if n.Namespace != ns {
			continue
		}
		removeCollFromNsPvtWriteSet(n, coll)
		if len(n.CollectionPvtRwset) == 0 {
			p.NsPvtRwset = append(p.NsPvtRwset[:i], p.NsPvtRwset[i+1:]...)
		}
		return
	}
}

func removeCollFromNsPvtWriteSet(n *rwset.NsPvtReadWriteSet, collName string) {
	for i := 0; i < len(n.CollectionPvtRwset); i++ {
		c := n.CollectionPvtRwset[i]
		if c.CollectionName != collName {
			continue
		}
		n.CollectionPvtRwset = append(n.CollectionPvtRwset[:i], n.CollectionPvtRwset[i+1:]...)
		return
	}
}

type nsColl struct {
	ns, coll string
}

func findInvalidNsPvtData(nsRwset *rwset.NsPvtReadWriteSet, txRWSet *rwsetutil.TxRwSet, blkNum, txNum uint64) (
	[]*ledger.PvtdataHashMismatch, []*nsColl,
) {
	var invalidPvtData []*ledger.PvtdataHashMismatch
	var invalidNsColl []*nsColl

	ns := nsRwset.Namespace
	for _, collPvtRwset := range nsRwset.CollectionPvtRwset {
		coll := collPvtRwset.CollectionName
		rwsetHash := txRWSet.GetPvtDataHash(ns, coll)
		if rwsetHash == nil {
			logger.Warningf("namespace: %s collection: %s was not accessed by txNum %d in BlkNum %d. "+
				"Unnecessary pvtdata has been passed", ns, coll, txNum, blkNum)
			invalidNsColl = append(invalidNsColl, &nsColl{ns, coll})
			continue
		}

		if !bytes.Equal(util.ComputeSHA256(collPvtRwset.Rwset), rwsetHash) {
			invalidPvtData = append(invalidPvtData, &ledger.PvtdataHashMismatch{
				BlockNum:     blkNum,
				TxNum:        txNum,
				Namespace:    ns,
				Collection:   coll,
				ExpectedHash: rwsetHash})
			invalidNsColl = append(invalidNsColl, &nsColl{ns, coll})
		}
	}
	return invalidPvtData, invalidNsColl
}
