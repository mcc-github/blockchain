/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package internal

import (
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset/kvrwset"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/storageutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
)

func prepareTxOps(rwset *rwsetutil.TxRwSet, txht *version.Height,
	precedingUpdates *PubAndHashUpdates, db privacyenabledstate.DB) (txOps, error) {
	txops := txOps{}
	txops.applyTxRwset(rwset)
	
	for ck, keyop := range txops {
		
		
		if keyop.isDelete() || keyop.isUpsertAndMetadataUpdate() {
			continue
		}

		
		if keyop.isOnlyUpsert() {
			latestMetadata, err := retrieveLatestMetadata(ck.ns, ck.coll, ck.key, precedingUpdates, db)
			if err != nil {
				return nil, err
			}
			keyop.metadata = latestMetadata
			continue
		}

		
		
		latestVal, err := retrieveLatestState(ck.ns, ck.coll, ck.key, precedingUpdates, db)
		if err != nil {
			return nil, err
		}
		if latestVal != nil {
			keyop.value = latestVal.Value
		} else {
			delete(txops, ck)
		}
	}
	
	return txops, nil
}



func (txops txOps) applyTxRwset(rwset *rwsetutil.TxRwSet) error {
	for _, nsRWSet := range rwset.NsRwSets {
		ns := nsRWSet.NameSpace
		for _, kvWrite := range nsRWSet.KvRwSet.Writes {
			txops.applyKVWrite(ns, "", kvWrite)
		}
		for _, kvMetadataWrite := range nsRWSet.KvRwSet.MetadataWrites {
			txops.applyMetadata(ns, "", kvMetadataWrite)
		}

		
		for _, collHashRWset := range nsRWSet.CollHashedRwSets {
			coll := collHashRWset.CollectionName
			for _, hashedWrite := range collHashRWset.HashedRwSet.HashedWrites {
				txops.applyKVWrite(ns, coll,
					&kvrwset.KVWrite{
						Key:      string(hashedWrite.KeyHash),
						Value:    hashedWrite.ValueHash,
						IsDelete: hashedWrite.IsDelete,
					},
				)
			}

			for _, metadataWrite := range collHashRWset.HashedRwSet.MetadataWrites {
				txops.applyMetadata(ns, coll,
					&kvrwset.KVMetadataWrite{
						Key:     string(metadataWrite.KeyHash),
						Entries: metadataWrite.Entries,
					},
				)
			}
		}
	}
	return nil
}


func (txops txOps) applyKVWrite(ns, coll string, kvWrite *kvrwset.KVWrite) {
	if kvWrite.IsDelete {
		txops.delete(compositeKey{ns, coll, kvWrite.Key})
	} else {
		txops.upsert(compositeKey{ns, coll, kvWrite.Key}, kvWrite.Value)
	}
}


func (txops txOps) applyMetadata(ns, coll string, metadataWrite *kvrwset.KVMetadataWrite) error {
	if metadataWrite.Entries == nil {
		txops.metadataDelete(compositeKey{ns, coll, metadataWrite.Key})
	} else {
		metadataBytes, err := storageutil.SerializeMetadata(metadataWrite.Entries)
		if err != nil {
			return err
		}
		txops.metadataUpdate(compositeKey{ns, coll, metadataWrite.Key}, metadataBytes)
	}
	return nil
}






func retrieveLatestState(ns, coll, key string,
	precedingUpdates *PubAndHashUpdates, db privacyenabledstate.DB) (*statedb.VersionedValue, error) {
	var vv *statedb.VersionedValue
	var err error
	if coll == "" {
		vv := precedingUpdates.PubUpdates.Get(ns, key)
		if vv == nil {
			vv, err = db.GetState(ns, key)
		}
		return vv, err
	}

	vv = precedingUpdates.HashUpdates.Get(ns, coll, key)
	if vv == nil {
		vv, err = db.GetValueHash(ns, coll, []byte(key))
	}
	return vv, err
}

func retrieveLatestMetadata(ns, coll, key string,
	precedingUpdates *PubAndHashUpdates, db privacyenabledstate.DB) ([]byte, error) {
	if coll == "" {
		vv := precedingUpdates.PubUpdates.Get(ns, key)
		if vv != nil {
			return vv.Metadata, nil
		}
		return db.GetStateMetadata(ns, key)
	}
	vv := precedingUpdates.HashUpdates.Get(ns, coll, key)
	if vv != nil {
		return vv.Metadata, nil
	}
	return db.GetPrivateDataMetadataByHash(ns, coll, []byte(key))
}

type keyOpsFlag uint8

const (
	upsertVal      keyOpsFlag = 1 
	metadataUpdate            = 2 
	metadataDelete            = 4 
	keyDelete                 = 8 
)

type compositeKey struct {
	ns, coll, key string
}

type txOps map[compositeKey]*keyOps

type keyOps struct {
	flag     keyOpsFlag
	value    []byte
	metadata []byte
}



func (txops txOps) upsert(k compositeKey, val []byte) {
	keyops := txops.getOrCreateKeyEntry(k)
	keyops.flag += upsertVal
	keyops.value = val
}

func (txops txOps) delete(k compositeKey) {
	keyops := txops.getOrCreateKeyEntry(k)
	keyops.flag += keyDelete
}

func (txops txOps) metadataUpdate(k compositeKey, metadata []byte) {
	keyops := txops.getOrCreateKeyEntry(k)
	keyops.flag += metadataUpdate
	keyops.metadata = metadata
}

func (txops txOps) metadataDelete(k compositeKey) {
	keyops := txops.getOrCreateKeyEntry(k)
	keyops.flag += metadataDelete
}

func (txops txOps) getOrCreateKeyEntry(k compositeKey) *keyOps {
	keyops, ok := txops[k]
	if !ok {
		keyops = &keyOps{}
		txops[k] = keyops
	}
	return keyops
}



func (keyops keyOps) isDelete() bool {
	return keyops.flag&(keyDelete) == keyDelete
}

func (keyops keyOps) isUpsertAndMetadataUpdate() bool {
	if keyops.flag&upsertVal == upsertVal {
		return keyops.flag&metadataUpdate == metadataUpdate ||
			keyops.flag&metadataDelete == metadataDelete
	}
	return false
}

func (keyops keyOps) isOnlyUpsert() bool {
	return keyops.flag|upsertVal == upsertVal
}
