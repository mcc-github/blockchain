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

package statebasedval

import (
	"fmt"
	"os"
	"testing"

	"github.com/mcc-github/blockchain-protos-go/ledger/rwset/kvrwset"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/internal"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/stretchr/testify/assert"
)

type keyValue struct {
	namespace  string
	collection string
	key        string
	keyHash    []byte
	value      []byte
	version    *version.Height
}

func TestMain(m *testing.M) {
	flogging.ActivateSpec("statevalidator,statebasedval,statecouchdb=debug")
	os.Exit(m.Run())
}

func TestValidatorBulkLoadingOfCache(t *testing.T) {
	testDBEnv := privacyenabledstate.CouchDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("testdb")

	validator := NewValidator(db)

	
	batch := privacyenabledstate.NewUpdateBatch()

	
	pubKV1 := keyValue{namespace: "ns1", key: "key1", value: []byte("value1"), version: version.NewHeight(1, 0)}
	pubKV2 := keyValue{namespace: "ns1", key: "key2", value: []byte("value2"), version: version.NewHeight(1, 1)}

	
	hashedKV1 := keyValue{namespace: "ns2", collection: "col1", key: "hashedPvtKey1",
		keyHash: util.ComputeStringHash("hashedPvtKey1"), value: []byte("value1"),
		version: version.NewHeight(1, 2)}
	hashedKV2 := keyValue{namespace: "ns2", collection: "col2", key: "hashedPvtKey2",
		keyHash: util.ComputeStringHash("hashedPvtKey2"), value: []byte("value2"),
		version: version.NewHeight(1, 3)}

	
	batch.PubUpdates.Put(pubKV1.namespace, pubKV1.key, pubKV1.value, pubKV1.version)
	batch.PubUpdates.Put(pubKV2.namespace, pubKV2.key, pubKV2.value, pubKV2.version)
	batch.HashUpdates.Put(hashedKV1.namespace, hashedKV1.collection, hashedKV1.keyHash, hashedKV1.value, hashedKV1.version)
	batch.HashUpdates.Put(hashedKV2.namespace, hashedKV2.collection, hashedKV2.keyHash, hashedKV2.value, hashedKV2.version)

	db.ApplyPrivacyAwareUpdates(batch, version.NewHeight(1, 4))

	
	
	rwsetBuilder1 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder1.AddToReadSet(pubKV1.namespace, pubKV1.key, pubKV1.version)
	rwsetBuilder1.AddToReadSet(pubKV2.namespace, pubKV2.key, pubKV2.version)
	rwsetBuilder1.AddToHashedReadSet(hashedKV1.namespace, hashedKV1.collection, hashedKV1.key, hashedKV1.version)
	rwsetBuilder1.AddToHashedReadSet(hashedKV2.namespace, hashedKV2.collection, hashedKV2.key, hashedKV2.version)

	
	rwsetBuilder2 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder2.AddToReadSet("ns3", "key1", nil)
	rwsetBuilder2.AddToHashedReadSet("ns3", "col1", "hashedPvtKey1", nil)

	
	transRWSets := getTestPubSimulationRWSet(t, rwsetBuilder1, rwsetBuilder2)
	var trans []*internal.Transaction
	for i, tranRWSet := range transRWSets {
		tx := &internal.Transaction{
			ID:             fmt.Sprintf("txid-%d", i),
			IndexInBlock:   i,
			ValidationCode: peer.TxValidationCode_VALID,
			RWSet:          tranRWSet,
		}
		trans = append(trans, tx)
	}
	block := &internal.Block{Num: 1, Txs: trans}

	if validator.db.IsBulkOptimizable() {

		commonStorageDB := validator.db.(*privacyenabledstate.CommonStorageDB)
		bulkOptimizable, _ := commonStorageDB.VersionedDB.(statedb.BulkOptimizable)

		
		validator.db.ClearCachedVersions()

		validator.preLoadCommittedVersionOfRSet(block)

		
		version, keyFound := bulkOptimizable.GetCachedVersion(pubKV1.namespace, pubKV1.key)
		assert.True(t, keyFound)
		assert.Equal(t, pubKV1.version, version)

		
		version, keyFound = bulkOptimizable.GetCachedVersion(pubKV2.namespace, pubKV2.key)
		assert.True(t, keyFound)
		assert.Equal(t, pubKV2.version, version)

		
		
		version, keyFound = bulkOptimizable.GetCachedVersion("ns3", "key1")
		assert.True(t, keyFound)
		assert.Nil(t, version)

		
		version, keyFound = bulkOptimizable.GetCachedVersion("ns4", "key1")
		assert.False(t, keyFound)
		assert.Nil(t, version)

		
		version, keyFound = validator.db.GetCachedKeyHashVersion(hashedKV1.namespace,
			hashedKV1.collection, hashedKV1.keyHash)
		assert.True(t, keyFound)
		assert.Equal(t, hashedKV1.version, version)

		
		version, keyFound = validator.db.GetCachedKeyHashVersion(hashedKV2.namespace,
			hashedKV2.collection, hashedKV2.keyHash)
		assert.True(t, keyFound)
		assert.Equal(t, hashedKV2.version, version)

		
		
		version, keyFound = validator.db.GetCachedKeyHashVersion("ns3", "col1", util.ComputeStringHash("hashedPvtKey1"))
		assert.True(t, keyFound)
		assert.Nil(t, version)

		
		version, keyFound = validator.db.GetCachedKeyHashVersion("ns4", "col1", util.ComputeStringHash("key1"))
		assert.False(t, keyFound)
		assert.Nil(t, version)

		
		validator.db.ClearCachedVersions()

		
		version, keyFound = bulkOptimizable.GetCachedVersion(pubKV1.namespace, pubKV1.key)
		assert.False(t, keyFound)
		assert.Nil(t, version)

		
		version, keyFound = validator.db.GetCachedKeyHashVersion("ns3", "col1", util.ComputeStringHash("hashedPvtKey1"))
		assert.False(t, keyFound)
		assert.Nil(t, version)
	}
}

func TestValidator(t *testing.T) {
	testDBEnv := privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("TestDB")

	
	batch := privacyenabledstate.NewUpdateBatch()
	batch.PubUpdates.Put("ns1", "key1", []byte("value1"), version.NewHeight(1, 0))
	batch.PubUpdates.Put("ns1", "key2", []byte("value2"), version.NewHeight(1, 1))
	batch.PubUpdates.Put("ns1", "key3", []byte("value3"), version.NewHeight(1, 2))
	batch.PubUpdates.Put("ns1", "key4", []byte("value4"), version.NewHeight(1, 3))
	batch.PubUpdates.Put("ns1", "key5", []byte("value5"), version.NewHeight(1, 4))
	db.ApplyPrivacyAwareUpdates(batch, version.NewHeight(1, 4))

	validator := NewValidator(db)

	
	rwsetBuilder1 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder1.AddToReadSet("ns1", "key1", version.NewHeight(1, 0))
	rwsetBuilder1.AddToReadSet("ns2", "key2", nil)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder1), []int{})

	
	rwsetBuilder2 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder2.AddToReadSet("ns1", "key1", version.NewHeight(1, 1))
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder2), []int{0})

	
	rwsetBuilder3 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder3.AddToReadSet("ns1", "key1", nil)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder3), []int{0})

	
	rwsetBuilder4 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder4.AddToReadSet("ns1", "key1", version.NewHeight(1, 0))
	rwsetBuilder4.AddToWriteSet("ns1", "key1", []byte("value1_new"))

	rwsetBuilder5 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder5.AddToReadSet("ns1", "key1", version.NewHeight(1, 0))
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder4, rwsetBuilder5), []int{1})
}

func TestPhantomValidation(t *testing.T) {
	testDBEnv := privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("TestDB")

	
	batch := privacyenabledstate.NewUpdateBatch()
	batch.PubUpdates.Put("ns1", "key1", []byte("value1"), version.NewHeight(1, 0))
	batch.PubUpdates.Put("ns1", "key2", []byte("value2"), version.NewHeight(1, 1))
	batch.PubUpdates.Put("ns1", "key3", []byte("value3"), version.NewHeight(1, 2))
	batch.PubUpdates.Put("ns1", "key4", []byte("value4"), version.NewHeight(1, 3))
	batch.PubUpdates.Put("ns1", "key5", []byte("value5"), version.NewHeight(1, 4))
	db.ApplyPrivacyAwareUpdates(batch, version.NewHeight(1, 4))

	validator := NewValidator(db)

	
	rwsetBuilder1 := rwsetutil.NewRWSetBuilder()
	rqi1 := &kvrwset.RangeQueryInfo{StartKey: "key2", EndKey: "key4", ItrExhausted: true}
	rwsetutil.SetRawReads(rqi1, []*kvrwset.KVRead{
		rwsetutil.NewKVRead("key2", version.NewHeight(1, 1)),
		rwsetutil.NewKVRead("key3", version.NewHeight(1, 2))})
	rwsetBuilder1.AddToRangeQuerySet("ns1", rqi1)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder1), []int{})

	
	rwsetBuilder2 := rwsetutil.NewRWSetBuilder()
	rqi2 := &kvrwset.RangeQueryInfo{StartKey: "key2", EndKey: "key4", ItrExhausted: false}
	rwsetutil.SetRawReads(rqi2, []*kvrwset.KVRead{
		rwsetutil.NewKVRead("key2", version.NewHeight(1, 1)),
		rwsetutil.NewKVRead("key3", version.NewHeight(1, 2)),
		rwsetutil.NewKVRead("key4", version.NewHeight(1, 2))})
	rwsetBuilder2.AddToRangeQuerySet("ns1", rqi2)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder2), []int{0})

	
	rwsetBuilder3 := rwsetutil.NewRWSetBuilder()
	rqi3 := &kvrwset.RangeQueryInfo{StartKey: "key2", EndKey: "key4", ItrExhausted: false}
	rwsetutil.SetRawReads(rqi3, []*kvrwset.KVRead{
		rwsetutil.NewKVRead("key2", version.NewHeight(1, 1)),
		rwsetutil.NewKVRead("key4", version.NewHeight(1, 3))})
	rwsetBuilder3.AddToRangeQuerySet("ns1", rqi3)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder3), []int{0})

	
	rwsetBuilder4 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder4.AddToWriteSet("ns1", "key3", nil)
	rwsetBuilder5 := rwsetutil.NewRWSetBuilder()
	rqi5 := &kvrwset.RangeQueryInfo{StartKey: "key2", EndKey: "key4", ItrExhausted: false}
	rwsetutil.SetRawReads(rqi5, []*kvrwset.KVRead{
		rwsetutil.NewKVRead("key2", version.NewHeight(1, 1)),
		rwsetutil.NewKVRead("key3", version.NewHeight(1, 2)),
		rwsetutil.NewKVRead("key4", version.NewHeight(1, 3))})
	rwsetBuilder5.AddToRangeQuerySet("ns1", rqi5)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder4, rwsetBuilder5), []int{1})

	
	rwsetBuilder6 := rwsetutil.NewRWSetBuilder()
	rwsetBuilder6.AddToWriteSet("ns1", "key2_1", []byte("value2_1"))

	rwsetBuilder7 := rwsetutil.NewRWSetBuilder()
	rqi7 := &kvrwset.RangeQueryInfo{StartKey: "key2", EndKey: "key4", ItrExhausted: false}
	rwsetutil.SetRawReads(rqi7, []*kvrwset.KVRead{
		rwsetutil.NewKVRead("key2", version.NewHeight(1, 1)),
		rwsetutil.NewKVRead("key3", version.NewHeight(1, 2)),
		rwsetutil.NewKVRead("key4", version.NewHeight(1, 3))})
	rwsetBuilder7.AddToRangeQuerySet("ns1", rqi7)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder6, rwsetBuilder7), []int{1})
}

func TestPhantomHashBasedValidation(t *testing.T) {
	testDBEnv := privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("TestDB")

	
	batch := privacyenabledstate.NewUpdateBatch()
	batch.PubUpdates.Put("ns1", "key1", []byte("value1"), version.NewHeight(1, 0))
	batch.PubUpdates.Put("ns1", "key2", []byte("value2"), version.NewHeight(1, 1))
	batch.PubUpdates.Put("ns1", "key3", []byte("value3"), version.NewHeight(1, 2))
	batch.PubUpdates.Put("ns1", "key4", []byte("value4"), version.NewHeight(1, 3))
	batch.PubUpdates.Put("ns1", "key5", []byte("value5"), version.NewHeight(1, 4))
	batch.PubUpdates.Put("ns1", "key6", []byte("value6"), version.NewHeight(1, 5))
	batch.PubUpdates.Put("ns1", "key7", []byte("value7"), version.NewHeight(1, 6))
	batch.PubUpdates.Put("ns1", "key8", []byte("value8"), version.NewHeight(1, 7))
	batch.PubUpdates.Put("ns1", "key9", []byte("value9"), version.NewHeight(1, 8))
	db.ApplyPrivacyAwareUpdates(batch, version.NewHeight(1, 8))

	validator := NewValidator(db)

	rwsetBuilder1 := rwsetutil.NewRWSetBuilder()
	rqi1 := &kvrwset.RangeQueryInfo{StartKey: "key2", EndKey: "key9", ItrExhausted: true}
	kvReadsDuringSimulation1 := []*kvrwset.KVRead{
		rwsetutil.NewKVRead("key2", version.NewHeight(1, 1)),
		rwsetutil.NewKVRead("key3", version.NewHeight(1, 2)),
		rwsetutil.NewKVRead("key4", version.NewHeight(1, 3)),
		rwsetutil.NewKVRead("key5", version.NewHeight(1, 4)),
		rwsetutil.NewKVRead("key6", version.NewHeight(1, 5)),
		rwsetutil.NewKVRead("key7", version.NewHeight(1, 6)),
		rwsetutil.NewKVRead("key8", version.NewHeight(1, 7)),
	}
	rwsetutil.SetMerkelSummary(rqi1, buildTestHashResults(t, 2, kvReadsDuringSimulation1))
	rwsetBuilder1.AddToRangeQuerySet("ns1", rqi1)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder1), []int{})

	rwsetBuilder2 := rwsetutil.NewRWSetBuilder()
	rqi2 := &kvrwset.RangeQueryInfo{StartKey: "key1", EndKey: "key9", ItrExhausted: false}
	kvReadsDuringSimulation2 := []*kvrwset.KVRead{
		rwsetutil.NewKVRead("key1", version.NewHeight(1, 0)),
		rwsetutil.NewKVRead("key2", version.NewHeight(1, 1)),
		rwsetutil.NewKVRead("key3", version.NewHeight(1, 1)),
		rwsetutil.NewKVRead("key4", version.NewHeight(1, 3)),
		rwsetutil.NewKVRead("key5", version.NewHeight(1, 4)),
		rwsetutil.NewKVRead("key6", version.NewHeight(1, 5)),
		rwsetutil.NewKVRead("key7", version.NewHeight(1, 6)),
		rwsetutil.NewKVRead("key8", version.NewHeight(1, 7)),
		rwsetutil.NewKVRead("key9", version.NewHeight(1, 8)),
	}
	rwsetutil.SetMerkelSummary(rqi2, buildTestHashResults(t, 2, kvReadsDuringSimulation2))
	rwsetBuilder2.AddToRangeQuerySet("ns1", rqi2)
	checkValidation(t, validator, getTestPubSimulationRWSet(t, rwsetBuilder2), []int{0})
}

func checkValidation(t *testing.T, val *Validator, transRWSets []*rwsetutil.TxRwSet, expectedInvalidTxIndexes []int) {
	var trans []*internal.Transaction
	for i, tranRWSet := range transRWSets {
		tx := &internal.Transaction{
			ID:             fmt.Sprintf("txid-%d", i),
			IndexInBlock:   i,
			ValidationCode: peer.TxValidationCode_VALID,
			RWSet:          tranRWSet,
		}
		trans = append(trans, tx)
	}
	block := &internal.Block{Num: 1, Txs: trans}
	_, err := val.ValidateAndPrepareBatch(block, true)
	assert.NoError(t, err)
	t.Logf("block.Txs[0].ValidationCode = %d", block.Txs[0].ValidationCode)
	var invalidTxs []int
	for _, tx := range block.Txs {
		if tx.ValidationCode != peer.TxValidationCode_VALID {
			invalidTxs = append(invalidTxs, tx.IndexInBlock)
		}
	}
	assert.Equal(t, len(expectedInvalidTxIndexes), len(invalidTxs))
	assert.ElementsMatch(t, invalidTxs, expectedInvalidTxIndexes)
}

func buildTestHashResults(t *testing.T, maxDegree int, kvReads []*kvrwset.KVRead) *kvrwset.QueryReadsMerkleSummary {
	if len(kvReads) <= maxDegree {
		t.Fatal("This method should be called with number of KVReads more than maxDegree; Else, hashing won't be performedrwset")
	}
	helper, _ := rwsetutil.NewRangeQueryResultsHelper(true, uint32(maxDegree))
	for _, kvRead := range kvReads {
		helper.AddResult(kvRead)
	}
	_, h, err := helper.Done()
	assert.NoError(t, err)
	assert.NotNil(t, h)
	return h
}

func getTestPubSimulationRWSet(t *testing.T, builders ...*rwsetutil.RWSetBuilder) []*rwsetutil.TxRwSet {
	var pubRWSets []*rwsetutil.TxRwSet
	for _, b := range builders {
		s, e := b.GetTxSimulationResults()
		assert.NoError(t, e)
		sBytes, err := s.GetPubSimulationBytes()
		assert.NoError(t, err)
		pubRWSet := &rwsetutil.TxRwSet{}
		assert.NoError(t, pubRWSet.FromProtoBytes(sBytes))
		pubRWSets = append(pubRWSets, pubRWSet)
	}
	return pubRWSets
}
