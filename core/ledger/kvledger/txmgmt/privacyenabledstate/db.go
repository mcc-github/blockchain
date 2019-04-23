/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privacyenabledstate

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
)


type DBProvider interface {
	
	GetDBHandle(id string) (DB, error)
	
	Close()
}


type DB interface {
	statedb.VersionedDB
	IsBulkOptimizable() bool
	LoadCommittedVersionsOfPubAndHashedKeys(pubKeys []*statedb.CompositeKey, hashedKeys []*HashedCompositeKey) error
	GetCachedKeyHashVersion(namespace, collection string, keyHash []byte) (*version.Height, bool)
	ClearCachedVersions()
	GetChaincodeEventListener() cceventmgmt.ChaincodeLifecycleEventListener
	GetPrivateData(namespace, collection, key string) (*statedb.VersionedValue, error)
	GetPrivateDataHash(namespace, collection, key string) (*statedb.VersionedValue, error)
	GetValueHash(namespace, collection string, keyHash []byte) (*statedb.VersionedValue, error)
	GetKeyHashVersion(namespace, collection string, keyHash []byte) (*version.Height, error)
	GetPrivateDataMultipleKeys(namespace, collection string, keys []string) ([]*statedb.VersionedValue, error)
	GetPrivateDataRangeScanIterator(namespace, collection, startKey, endKey string) (statedb.ResultsIterator, error)
	GetStateMetadata(namespace, key string) ([]byte, error)
	GetPrivateDataMetadataByHash(namespace, collection string, keyHash []byte) ([]byte, error)
	ExecuteQueryOnPrivateData(namespace, collection, query string) (statedb.ResultsIterator, error)
	ApplyPrivacyAwareUpdates(updates *UpdateBatch, height *version.Height) error
}


type PvtdataCompositeKey struct {
	Namespace      string
	CollectionName string
	Key            string
}


type HashedCompositeKey struct {
	Namespace      string
	CollectionName string
	KeyHash        string
}


type PvtKVWrite struct {
	Key      string
	IsDelete bool
	Value    []byte
	Version  *version.Height
}



type UpdateBatch struct {
	PubUpdates  *PubUpdateBatch
	HashUpdates *HashedUpdateBatch
	PvtUpdates  *PvtUpdateBatch
}


type PubUpdateBatch struct {
	*statedb.UpdateBatch
}


type HashedUpdateBatch struct {
	UpdateMap
}


type PvtUpdateBatch struct {
	UpdateMap
}


type UpdateMap map[string]nsBatch


type nsBatch struct {
	*statedb.UpdateBatch
}


func NewUpdateBatch() *UpdateBatch {
	return &UpdateBatch{NewPubUpdateBatch(), NewHashedUpdateBatch(), NewPvtUpdateBatch()}
}


func NewPubUpdateBatch() *PubUpdateBatch {
	return &PubUpdateBatch{statedb.NewUpdateBatch()}
}


func NewHashedUpdateBatch() *HashedUpdateBatch {
	return &HashedUpdateBatch{make(map[string]nsBatch)}
}


func NewPvtUpdateBatch() *PvtUpdateBatch {
	return &PvtUpdateBatch{make(map[string]nsBatch)}
}


func (b UpdateMap) IsEmpty() bool {
	return len(b) == 0
}


func (b UpdateMap) Put(ns, coll, key string, value []byte, version *version.Height) {
	b.PutValAndMetadata(ns, coll, key, value, nil, version)
}


func (b UpdateMap) PutValAndMetadata(ns, coll, key string, value []byte, metadata []byte, version *version.Height) {
	b.getOrCreateNsBatch(ns).PutValAndMetadata(coll, key, value, metadata, version)
}


func (b UpdateMap) Delete(ns, coll, key string, version *version.Height) {
	b.getOrCreateNsBatch(ns).Delete(coll, key, version)
}


func (b UpdateMap) Get(ns, coll, key string) *statedb.VersionedValue {
	nsPvtBatch, ok := b[ns]
	if !ok {
		return nil
	}
	return nsPvtBatch.Get(coll, key)
}


func (b UpdateMap) Contains(ns, coll, key string) bool {
	nsBatch, ok := b[ns]
	if !ok {
		return false
	}
	return nsBatch.Exists(coll, key)
}

func (nsb nsBatch) GetCollectionNames() []string {
	return nsb.GetUpdatedNamespaces()
}

func (b UpdateMap) getOrCreateNsBatch(ns string) nsBatch {
	batch, ok := b[ns]
	if !ok {
		batch = nsBatch{statedb.NewUpdateBatch()}
		b[ns] = batch
	}
	return batch
}


func (h HashedUpdateBatch) Contains(ns, coll string, keyHash []byte) bool {
	return h.UpdateMap.Contains(ns, coll, string(keyHash))
}


func (h HashedUpdateBatch) Put(ns, coll string, key []byte, value []byte, version *version.Height) {
	h.PutValHashAndMetadata(ns, coll, key, value, nil, version)
}



func (h HashedUpdateBatch) PutValHashAndMetadata(ns, coll string, key []byte, value []byte, metadata []byte, version *version.Height) {
	h.UpdateMap.PutValAndMetadata(ns, coll, string(key), value, metadata, version)
}


func (h HashedUpdateBatch) Delete(ns, coll string, key []byte, version *version.Height) {
	h.UpdateMap.Delete(ns, coll, string(key), version)
}


func (h HashedUpdateBatch) ToCompositeKeyMap() map[HashedCompositeKey]*statedb.VersionedValue {
	m := make(map[HashedCompositeKey]*statedb.VersionedValue)
	for ns, nsBatch := range h.UpdateMap {
		for _, coll := range nsBatch.GetCollectionNames() {
			for key, vv := range nsBatch.GetUpdates(coll) {
				m[HashedCompositeKey{ns, coll, key}] = vv
			}
		}
	}
	return m
}


type PvtdataCompositeKeyMap map[PvtdataCompositeKey]*statedb.VersionedValue


func (p PvtUpdateBatch) ToCompositeKeyMap() PvtdataCompositeKeyMap {
	m := make(PvtdataCompositeKeyMap)
	for ns, nsBatch := range p.UpdateMap {
		for _, coll := range nsBatch.GetCollectionNames() {
			for key, vv := range nsBatch.GetUpdates(coll) {
				m[PvtdataCompositeKey{ns, coll, key}] = vv
			}
		}
	}
	return m
}


func (hck *HashedCompositeKey) String() string {
	return fmt.Sprintf("ns=%s, collection=%s, keyHash=%x", hck.Namespace, hck.CollectionName, hck.KeyHash)
}
