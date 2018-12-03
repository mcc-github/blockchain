/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statedb

import (
	"fmt"
	"sort"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/util"
)





type VersionedDBProvider interface {
	
	GetDBHandle(id string) (VersionedDB, error)
	
	Close()
}


type VersionedDB interface {
	
	GetState(namespace string, key string) (*VersionedValue, error)
	
	GetVersion(namespace string, key string) (*version.Height, error)
	
	GetStateMultipleKeys(namespace string, keys []string) ([]*VersionedValue, error)
	
	
	
	
	GetStateRangeScanIterator(namespace string, startKey string, endKey string) (ResultsIterator, error)
	
	
	
	
	
	GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (QueryResultsIterator, error)
	
	ExecuteQuery(namespace, query string) (ResultsIterator, error)
	
	
	
	ExecuteQueryWithMetadata(namespace, query string, metadata map[string]interface{}) (QueryResultsIterator, error)
	
	
	
	ApplyUpdates(batch *UpdateBatch, height *version.Height) error
	
	
	GetLatestSavePoint() (*version.Height, error)
	
	
	
	
	
	ValidateKeyValue(key string, value []byte) error
	
	
	BytesKeySupported() bool
	
	Open() error
	
	Close()
}



type BulkOptimizable interface {
	LoadCommittedVersions(keys []*CompositeKey) error
	GetCachedVersion(namespace, key string) (*version.Height, bool)
	ClearCachedVersions()
}



type IndexCapable interface {
	GetDBType() string
	ProcessIndexesForChaincodeDeploy(namespace string, fileEntries []*ccprovider.TarFileEntry) error
}


type CompositeKey struct {
	Namespace string
	Key       string
}


type VersionedValue struct {
	Value    []byte
	Metadata []byte
	Version  *version.Height
}


func (vv *VersionedValue) IsDelete() bool {
	return vv.Value == nil
}


type VersionedKV struct {
	CompositeKey
	VersionedValue
}


type ResultsIterator interface {
	Next() (QueryResult, error)
	Close()
}


type QueryResultsIterator interface {
	ResultsIterator
	GetBookmarkAndClose() string
}


type QueryResult interface{}

type nsUpdates struct {
	m map[string]*VersionedValue
}

func newNsUpdates() *nsUpdates {
	return &nsUpdates{make(map[string]*VersionedValue)}
}


type UpdateBatch struct {
	updates map[string]*nsUpdates
}


func NewUpdateBatch() *UpdateBatch {
	return &UpdateBatch{make(map[string]*nsUpdates)}
}


func (batch *UpdateBatch) Get(ns string, key string) *VersionedValue {
	nsUpdates, ok := batch.updates[ns]
	if !ok {
		return nil
	}
	vv, ok := nsUpdates.m[key]
	if !ok {
		return nil
	}
	return vv
}


func (batch *UpdateBatch) Put(ns string, key string, value []byte, version *version.Height) {
	batch.PutValAndMetadata(ns, key, value, nil, version)
}



func (batch *UpdateBatch) PutValAndMetadata(ns string, key string, value []byte, metadata []byte, version *version.Height) {
	if value == nil {
		panic("Nil value not allowed. Instead call 'Delete' function")
	}
	batch.Update(ns, key, &VersionedValue{value, metadata, version})
}


func (batch *UpdateBatch) Delete(ns string, key string, version *version.Height) {
	batch.Update(ns, key, &VersionedValue{nil, nil, version})
}


func (batch *UpdateBatch) Exists(ns string, key string) bool {
	nsUpdates, ok := batch.updates[ns]
	if !ok {
		return false
	}
	_, ok = nsUpdates.m[key]
	return ok
}


func (batch *UpdateBatch) GetUpdatedNamespaces() []string {
	namespaces := make([]string, len(batch.updates))
	i := 0
	for ns := range batch.updates {
		namespaces[i] = ns
		i++
	}
	return namespaces
}


func (batch *UpdateBatch) Update(ns string, key string, vv *VersionedValue) {
	batch.getOrCreateNsUpdates(ns).m[key] = vv
}


func (batch *UpdateBatch) GetUpdates(ns string) map[string]*VersionedValue {
	nsUpdates, ok := batch.updates[ns]
	if !ok {
		return nil
	}
	return nsUpdates.m
}








func (batch *UpdateBatch) GetRangeScanIterator(ns string, startKey string, endKey string) QueryResultsIterator {
	return newNsIterator(ns, startKey, endKey, batch)
}

func (batch *UpdateBatch) getOrCreateNsUpdates(ns string) *nsUpdates {
	nsUpdates := batch.updates[ns]
	if nsUpdates == nil {
		nsUpdates = newNsUpdates()
		batch.updates[ns] = nsUpdates
	}
	return nsUpdates
}

type nsIterator struct {
	ns         string
	nsUpdates  *nsUpdates
	sortedKeys []string
	nextIndex  int
	lastIndex  int
}

func newNsIterator(ns string, startKey string, endKey string, batch *UpdateBatch) *nsIterator {
	nsUpdates, ok := batch.updates[ns]
	if !ok {
		return &nsIterator{}
	}
	sortedKeys := util.GetSortedKeys(nsUpdates.m)
	var nextIndex int
	var lastIndex int
	if startKey == "" {
		nextIndex = 0
	} else {
		nextIndex = sort.SearchStrings(sortedKeys, startKey)
	}
	if endKey == "" {
		lastIndex = len(sortedKeys)
	} else {
		lastIndex = sort.SearchStrings(sortedKeys, endKey)
	}
	return &nsIterator{ns, nsUpdates, sortedKeys, nextIndex, lastIndex}
}


func (itr *nsIterator) Next() (QueryResult, error) {
	if itr.nextIndex >= itr.lastIndex {
		return nil, nil
	}
	key := itr.sortedKeys[itr.nextIndex]
	vv := itr.nsUpdates.m[key]
	itr.nextIndex++
	return &VersionedKV{CompositeKey{itr.ns, key}, VersionedValue{vv.Value, vv.Metadata, vv.Version}}, nil
}


func (itr *nsIterator) Close() {
	
}


func (itr *nsIterator) GetBookmarkAndClose() string {
	
	return ""
}

const optionLimit = "limit"


func ValidateRangeMetadata(metadata map[string]interface{}) error {
	for key, keyVal := range metadata {
		switch key {

		case optionLimit:
			
			if _, ok := keyVal.(int32); ok {
				continue
			}
			return fmt.Errorf("Invalid entry, \"limit\" must be a int32")

		default:
			return fmt.Errorf("Invalid entry, option %s not recognized", key)
		}
	}
	return nil
}
