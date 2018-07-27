/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rwsetutil

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
)

var logger = flogging.MustGetLogger("rwsetutil")


type RWSetBuilder struct {
	pubRwBuilderMap map[string]*nsPubRwBuilder
	pvtRwBuilderMap map[string]*nsPvtRwBuilder
}

type nsPubRwBuilder struct {
	namespace         string
	readMap           map[string]*kvrwset.KVRead 
	writeMap          map[string]*kvrwset.KVWrite
	metadataWriteMap  map[string]*kvrwset.KVMetadataWrite
	rangeQueriesMap   map[rangeQueryKey]*kvrwset.RangeQueryInfo 
	rangeQueriesKeys  []rangeQueryKey
	collHashRwBuilder map[string]*collHashRwBuilder
}

type collHashRwBuilder struct {
	collName         string
	readMap          map[string]*kvrwset.KVReadHash
	writeMap         map[string]*kvrwset.KVWriteHash
	metadataWriteMap map[string]*kvrwset.KVMetadataWriteHash
	pvtDataHash      []byte
}

type nsPvtRwBuilder struct {
	namespace         string
	collPvtRwBuilders map[string]*collPvtRwBuilder
}

type collPvtRwBuilder struct {
	collectionName string
	writeMap       map[string]*kvrwset.KVWrite
}

type rangeQueryKey struct {
	startKey     string
	endKey       string
	itrExhausted bool
}


func NewRWSetBuilder() *RWSetBuilder {
	return &RWSetBuilder{make(map[string]*nsPubRwBuilder), make(map[string]*nsPvtRwBuilder)}
}


func (b *RWSetBuilder) AddToReadSet(ns string, key string, version *version.Height) {
	nsPubRwBuilder := b.getOrCreateNsPubRwBuilder(ns)
	nsPubRwBuilder.readMap[key] = NewKVRead(key, version)
}


func (b *RWSetBuilder) AddToWriteSet(ns string, key string, value []byte) {
	nsPubRwBuilder := b.getOrCreateNsPubRwBuilder(ns)
	nsPubRwBuilder.writeMap[key] = newKVWrite(key, value)
}



func (b *RWSetBuilder) AddToMetadataWriteSet(ns, key string, metadata map[string][]byte) {
	b.getOrCreateNsPubRwBuilder(ns).
		metadataWriteMap[key] = mapToMetadataWrite(key, metadata)
}


func (b *RWSetBuilder) AddToRangeQuerySet(ns string, rqi *kvrwset.RangeQueryInfo) {
	nsPubRwBuilder := b.getOrCreateNsPubRwBuilder(ns)
	key := rangeQueryKey{rqi.StartKey, rqi.EndKey, rqi.ItrExhausted}
	_, ok := nsPubRwBuilder.rangeQueriesMap[key]
	if !ok {
		nsPubRwBuilder.rangeQueriesMap[key] = rqi
		nsPubRwBuilder.rangeQueriesKeys = append(nsPubRwBuilder.rangeQueriesKeys, key)
	}
}


func (b *RWSetBuilder) AddToHashedReadSet(ns string, coll string, key string, version *version.Height) {
	kvReadHash := newPvtKVReadHash(key, version)
	b.getOrCreateCollHashedRwBuilder(ns, coll).readMap[key] = kvReadHash
}


func (b *RWSetBuilder) AddToPvtAndHashedWriteSet(ns string, coll string, key string, value []byte) {
	kvWrite, kvWriteHash := newPvtKVWriteAndHash(key, value)
	b.getOrCreateCollPvtRwBuilder(ns, coll).writeMap[key] = kvWrite
	b.getOrCreateCollHashedRwBuilder(ns, coll).writeMap[key] = kvWriteHash
}


func (b *RWSetBuilder) AddToHashedMetadataWriteSet(ns, coll, key string, metadata map[string][]byte) {
	b.getOrCreateCollHashedRwBuilder(ns, coll).
		metadataWriteMap[key] = mapToMetadataWriteHash(key, metadata)
}



func (b *RWSetBuilder) GetTxSimulationResults() (*ledger.TxSimulationResults, error) {
	pvtData := b.getTxPvtReadWriteSet()
	var err error

	var pubDataProto *rwset.TxReadWriteSet
	var pvtDataProto *rwset.TxPvtReadWriteSet

	
	if pvtData != nil {
		if pvtDataProto, err = pvtData.toProtoMsg(); err != nil {
			return nil, err
		}
		for _, ns := range pvtDataProto.NsPvtRwset {
			for _, coll := range ns.CollectionPvtRwset {
				b.setPvtCollectionHash(ns.Namespace, coll.CollectionName, coll.Rwset)
			}
		}
	}
	
	pubSet := b.GetTxReadWriteSet()
	if pubSet != nil {
		if pubDataProto, err = b.GetTxReadWriteSet().toProtoMsg(); err != nil {
			return nil, err
		}
	}
	return &ledger.TxSimulationResults{
		PubSimulationResults: pubDataProto,
		PvtSimulationResults: pvtDataProto,
	}, nil
}

func (b *RWSetBuilder) setPvtCollectionHash(ns string, coll string, pvtDataProto []byte) {
	collHashedBuilder := b.getOrCreateCollHashedRwBuilder(ns, coll)
	collHashedBuilder.pvtDataHash = util.ComputeHash(pvtDataProto)
}



func (b *RWSetBuilder) GetTxReadWriteSet() *TxRwSet {
	sortedNsPubBuilders := []*nsPubRwBuilder{}
	util.GetValuesBySortedKeys(&(b.pubRwBuilderMap), &sortedNsPubBuilders)

	var nsPubRwSets []*NsRwSet
	for _, nsPubRwBuilder := range sortedNsPubBuilders {
		nsPubRwSets = append(nsPubRwSets, nsPubRwBuilder.build())
	}
	return &TxRwSet{NsRwSets: nsPubRwSets}
}


func (b *RWSetBuilder) getTxPvtReadWriteSet() *TxPvtRwSet {
	sortedNsPvtBuilders := []*nsPvtRwBuilder{}
	util.GetValuesBySortedKeys(&(b.pvtRwBuilderMap), &sortedNsPvtBuilders)

	var nsPvtRwSets []*NsPvtRwSet
	for _, nsPvtRwBuilder := range sortedNsPvtBuilders {
		nsPvtRwSets = append(nsPvtRwSets, nsPvtRwBuilder.build())
	}
	if len(nsPvtRwSets) == 0 {
		return nil
	}
	return &TxPvtRwSet{NsPvtRwSet: nsPvtRwSets}
}

func (b *nsPubRwBuilder) build() *NsRwSet {
	var readSet []*kvrwset.KVRead
	var writeSet []*kvrwset.KVWrite
	var metadataWriteSet []*kvrwset.KVMetadataWrite
	var rangeQueriesInfo []*kvrwset.RangeQueryInfo
	var collHashedRwSet []*CollHashedRwSet
	
	util.GetValuesBySortedKeys(&(b.readMap), &readSet)
	
	util.GetValuesBySortedKeys(&(b.writeMap), &writeSet)
	util.GetValuesBySortedKeys(&(b.metadataWriteMap), &metadataWriteSet)
	
	for _, key := range b.rangeQueriesKeys {
		rangeQueriesInfo = append(rangeQueriesInfo, b.rangeQueriesMap[key])
	}
	
	sortedCollBuilders := []*collHashRwBuilder{}
	util.GetValuesBySortedKeys(&(b.collHashRwBuilder), &sortedCollBuilders)
	for _, collBuilder := range sortedCollBuilders {
		collHashedRwSet = append(collHashedRwSet, collBuilder.build())
	}
	return &NsRwSet{
		NameSpace: b.namespace,
		KvRwSet: &kvrwset.KVRWSet{
			Reads:            readSet,
			Writes:           writeSet,
			MetadataWrites:   metadataWriteSet,
			RangeQueriesInfo: rangeQueriesInfo,
		},
		CollHashedRwSets: collHashedRwSet,
	}
}

func (b *nsPvtRwBuilder) build() *NsPvtRwSet {
	sortedCollBuilders := []*collPvtRwBuilder{}
	util.GetValuesBySortedKeys(&(b.collPvtRwBuilders), &sortedCollBuilders)

	var collPvtRwSets []*CollPvtRwSet
	for _, collBuilder := range sortedCollBuilders {
		collPvtRwSets = append(collPvtRwSets, collBuilder.build())
	}
	return &NsPvtRwSet{NameSpace: b.namespace, CollPvtRwSets: collPvtRwSets}
}

func (b *collHashRwBuilder) build() *CollHashedRwSet {
	var readSet []*kvrwset.KVReadHash
	var writeSet []*kvrwset.KVWriteHash
	var metadataWriteSet []*kvrwset.KVMetadataWriteHash

	util.GetValuesBySortedKeys(&(b.readMap), &readSet)
	util.GetValuesBySortedKeys(&(b.writeMap), &writeSet)
	util.GetValuesBySortedKeys(&(b.metadataWriteMap), &metadataWriteSet)
	return &CollHashedRwSet{
		CollectionName: b.collName,
		HashedRwSet: &kvrwset.HashedRWSet{
			HashedReads:    readSet,
			HashedWrites:   writeSet,
			MetadataWrites: metadataWriteSet,
		},
		PvtRwSetHash: b.pvtDataHash,
	}
}

func (b *collPvtRwBuilder) build() *CollPvtRwSet {
	var writeSet []*kvrwset.KVWrite
	util.GetValuesBySortedKeys(&(b.writeMap), &writeSet)
	return &CollPvtRwSet{
		CollectionName: b.collectionName,
		KvRwSet: &kvrwset.KVRWSet{
			Writes: writeSet,
		},
	}
}

func (b *RWSetBuilder) getOrCreateNsPubRwBuilder(ns string) *nsPubRwBuilder {
	nsPubRwBuilder, ok := b.pubRwBuilderMap[ns]
	if !ok {
		nsPubRwBuilder = newNsPubRwBuilder(ns)
		b.pubRwBuilderMap[ns] = nsPubRwBuilder
	}
	return nsPubRwBuilder
}

func (b *RWSetBuilder) getOrCreateNsPvtRwBuilder(ns string) *nsPvtRwBuilder {
	nsPvtRwBuilder, ok := b.pvtRwBuilderMap[ns]
	if !ok {
		nsPvtRwBuilder = newNsPvtRwBuilder(ns)
		b.pvtRwBuilderMap[ns] = nsPvtRwBuilder
	}
	return nsPvtRwBuilder
}

func (b *RWSetBuilder) getOrCreateCollHashedRwBuilder(ns string, coll string) *collHashRwBuilder {
	nsPubRwBuilder := b.getOrCreateNsPubRwBuilder(ns)
	collHashRwBuilder, ok := nsPubRwBuilder.collHashRwBuilder[coll]
	if !ok {
		collHashRwBuilder = newCollHashRwBuilder(coll)
		nsPubRwBuilder.collHashRwBuilder[coll] = collHashRwBuilder
	}
	return collHashRwBuilder
}

func (b *RWSetBuilder) getOrCreateCollPvtRwBuilder(ns string, coll string) *collPvtRwBuilder {
	nsPvtRwBuilder := b.getOrCreateNsPvtRwBuilder(ns)
	collPvtRwBuilder, ok := nsPvtRwBuilder.collPvtRwBuilders[coll]
	if !ok {
		collPvtRwBuilder = newCollPvtRwBuilder(coll)
		nsPvtRwBuilder.collPvtRwBuilders[coll] = collPvtRwBuilder
	}
	return collPvtRwBuilder
}

func newNsPubRwBuilder(namespace string) *nsPubRwBuilder {
	return &nsPubRwBuilder{
		namespace,
		make(map[string]*kvrwset.KVRead),
		make(map[string]*kvrwset.KVWrite),
		make(map[string]*kvrwset.KVMetadataWrite),
		make(map[rangeQueryKey]*kvrwset.RangeQueryInfo),
		nil,
		make(map[string]*collHashRwBuilder),
	}
}

func newNsPvtRwBuilder(namespace string) *nsPvtRwBuilder {
	return &nsPvtRwBuilder{namespace, make(map[string]*collPvtRwBuilder)}
}

func newCollHashRwBuilder(collName string) *collHashRwBuilder {
	return &collHashRwBuilder{
		collName,
		make(map[string]*kvrwset.KVReadHash),
		make(map[string]*kvrwset.KVWriteHash),
		make(map[string]*kvrwset.KVMetadataWriteHash),
		nil,
	}
}

func newCollPvtRwBuilder(collName string) *collPvtRwBuilder {
	return &collPvtRwBuilder{collName, make(map[string]*kvrwset.KVWrite)}
}

func mapToMetadataWrite(key string, m map[string][]byte) *kvrwset.KVMetadataWrite {
	proto := &kvrwset.KVMetadataWrite{Key: key}
	names := util.GetSortedKeys(m)
	for _, name := range names {
		proto.Entries = append(proto.Entries,
			&kvrwset.KVMetadataEntry{Name: name, Value: m[name]},
		)
	}
	return proto
}

func mapToMetadataWriteHash(key string, m map[string][]byte) *kvrwset.KVMetadataWriteHash {
	proto := &kvrwset.KVMetadataWriteHash{KeyHash: util.ComputeStringHash(key)}
	names := util.GetSortedKeys(m)
	for _, name := range names {
		proto.Entries = append(proto.Entries,
			&kvrwset.KVMetadataEntry{Name: name, Value: m[name]},
		)
	}
	return proto
}
