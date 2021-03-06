/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valinternal

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/storageutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Printf("could not create temp dir %s", err)
		os.Exit(-1)
		return
	}
	flogging.SetModuleLevel("valinternal", "debug")
	viper.Set("peer.fileSystemPath", tempDir)
	os.Exit(m.Run())
}

func TestTxOpsPreparationValueUpdate(t *testing.T) {
	testDBEnv := privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("TestDB")

	ck1, ck2, ck3 :=
		compositeKey{ns: "ns1", key: "key1"},
		compositeKey{ns: "ns1", key: "key2"},
		compositeKey{ns: "ns1", key: "key3"}

	updateBatch := privacyenabledstate.NewUpdateBatch()
	updateBatch.PubUpdates.Put(ck1.ns, ck1.key, []byte("value1"), version.NewHeight(1, 1)) 
	updateBatch.PubUpdates.PutValAndMetadata(                                              
		ck2.ns, ck2.key,
		[]byte("value2"),
		testutilSerializedMetadata(t, map[string][]byte{"metadata2": []byte("metadata2")}),
		version.NewHeight(1, 2))

	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(1, 2)) 
	precedingUpdates := NewPubAndHashUpdates()

	rwset := testutilBuildRwset( 
		t,
		map[compositeKey][]byte{
			ck1: []byte("value1_new"),
			ck2: []byte("value2_new"),
			ck3: []byte("value3_new"),
		},
		nil,
	)

	txOps, err := prepareTxOps(rwset, version.NewHeight(1, 2), precedingUpdates, db)
	assert.NoError(t, err)
	assert.Len(t, txOps, 3)

	ck1ExpectedKeyOps := &keyOps{ 
		flag:  upsertVal,
		value: []byte("value1_new"),
	}

	ck2ExpectedKeyOps := &keyOps{ 
		flag:     upsertVal,
		value:    []byte("value2_new"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata2": []byte("metadata2")}),
	}

	ck3ExpectedKeyOps := &keyOps{ 
		flag:  upsertVal,
		value: []byte("value3_new"),
	}

	assert.Equal(t, ck1ExpectedKeyOps, txOps[ck1])
	assert.Equal(t, ck2ExpectedKeyOps, txOps[ck2])
	assert.Equal(t, ck3ExpectedKeyOps, txOps[ck3])
}

func TestTxOpsPreparationMetadataUpdates(t *testing.T) {
	testDBEnv := privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("TestDB")

	ck1, ck2, ck3 :=
		compositeKey{ns: "ns1", key: "key1"},
		compositeKey{ns: "ns1", key: "key2"},
		compositeKey{ns: "ns1", key: "key3"}

	updateBatch := privacyenabledstate.NewUpdateBatch()
	updateBatch.PubUpdates.Put(ck1.ns, ck1.key, []byte("value1"), version.NewHeight(1, 1)) 
	updateBatch.PubUpdates.PutValAndMetadata(                                              
		ck2.ns, ck2.key,
		[]byte("value2"),
		testutilSerializedMetadata(t, map[string][]byte{"metadata2": []byte("metadata2")}),
		version.NewHeight(1, 2))

	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(1, 2)) 
	precedingUpdates := NewPubAndHashUpdates()

	rwset := testutilBuildRwset( 
		t,
		nil,
		map[compositeKey]map[string][]byte{
			ck1: {"metadata1": []byte("metadata1_new")},
			ck2: {"metadata2": []byte("metadata2_new")},
			ck3: {"metadata3": []byte("metadata3_new")},
		},
	)

	txOps, err := prepareTxOps(rwset, version.NewHeight(1, 2), precedingUpdates, db)
	assert.NoError(t, err)
	assert.Len(t, txOps, 2) 

	ck1ExpectedKeyOps := &keyOps{ 
		flag:     metadataUpdate,
		value:    []byte("value1"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata1": []byte("metadata1_new")}),
	}

	ck2ExpectedKeyOps := &keyOps{ 
		flag:     metadataUpdate,
		value:    []byte("value2"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata2": []byte("metadata2_new")}),
	}

	assert.Equal(t, ck1ExpectedKeyOps, txOps[ck1])
	assert.Equal(t, ck2ExpectedKeyOps, txOps[ck2])
}

func TestTxOpsPreparationMetadataDelete(t *testing.T) {
	testDBEnv := privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("TestDB")

	ck1, ck2, ck3 :=
		compositeKey{ns: "ns1", key: "key1"},
		compositeKey{ns: "ns1", key: "key2"},
		compositeKey{ns: "ns1", key: "key3"}

	updateBatch := privacyenabledstate.NewUpdateBatch()
	updateBatch.PubUpdates.Put(ck1.ns, ck1.key, []byte("value1"), version.NewHeight(1, 1)) 
	updateBatch.PubUpdates.PutValAndMetadata(                                              
		ck2.ns, ck2.key,
		[]byte("value2"),
		testutilSerializedMetadata(t, map[string][]byte{"metadata2": []byte("metadata2")}),
		version.NewHeight(1, 2))

	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(1, 2)) 
	precedingUpdates := NewPubAndHashUpdates()

	rwset := testutilBuildRwset( 
		t,
		nil,
		map[compositeKey]map[string][]byte{
			ck1: {},
			ck2: {},
			ck3: {},
		},
	)

	txOps, err := prepareTxOps(rwset, version.NewHeight(1, 2), precedingUpdates, db)
	assert.NoError(t, err)
	assert.Len(t, txOps, 2) 

	ck1ExpectedKeyOps := &keyOps{ 
		flag:  metadataDelete,
		value: []byte("value1"),
	}

	ck2ExpectedKeyOps := &keyOps{ 
		flag:  metadataDelete,
		value: []byte("value2"),
	}

	assert.Equal(t, ck1ExpectedKeyOps, txOps[ck1])
	assert.Equal(t, ck2ExpectedKeyOps, txOps[ck2])
}

func TestTxOpsPreparationMixedUpdates(t *testing.T) {
	testDBEnv := privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("TestDB")

	ck1, ck2, ck3, ck4 :=
		compositeKey{ns: "ns1", key: "key1"},
		compositeKey{ns: "ns1", key: "key2"},
		compositeKey{ns: "ns1", key: "key3"},
		compositeKey{ns: "ns1", key: "key4"}

	updateBatch := privacyenabledstate.NewUpdateBatch()
	updateBatch.PubUpdates.Put(ck1.ns, ck1.key, []byte("value1"), version.NewHeight(1, 1)) 
	updateBatch.PubUpdates.Put(ck2.ns, ck2.key, []byte("value2"), version.NewHeight(1, 2)) 
	updateBatch.PubUpdates.PutValAndMetadata(                                              
		ck3.ns, ck3.key,
		[]byte("value3"),
		testutilSerializedMetadata(t, map[string][]byte{"metadata3": []byte("metadata3")}),
		version.NewHeight(1, 3))
	updateBatch.PubUpdates.PutValAndMetadata( 
		ck4.ns, ck4.key,
		[]byte("value4"),
		testutilSerializedMetadata(t, map[string][]byte{"metadata4": []byte("metadata4")}),
		version.NewHeight(1, 4))

	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(1, 2)) 

	precedingUpdates := NewPubAndHashUpdates()

	rwset := testutilBuildRwset( 
		t,
		map[compositeKey][]byte{
			ck1: []byte("value1_new"),
			ck2: []byte("value2_new"),
			ck4: []byte("value4_new"),
		},
		map[compositeKey]map[string][]byte{
			ck2: {"metadata2": []byte("metadata2_new")},
			ck3: {"metadata3": []byte("metadata3_new")},
		},
	)

	txOps, err := prepareTxOps(rwset, version.NewHeight(1, 2), precedingUpdates, db)
	assert.NoError(t, err)
	assert.Len(t, txOps, 4)

	ck1ExpectedKeyOps := &keyOps{ 
		flag:  upsertVal,
		value: []byte("value1_new"),
	}

	ck2ExpectedKeyOps := &keyOps{ 
		flag:     upsertVal + metadataUpdate,
		value:    []byte("value2_new"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata2": []byte("metadata2_new")}),
	}

	ck3ExpectedKeyOps := &keyOps{ 
		flag:     metadataUpdate,
		value:    []byte("value3"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata3": []byte("metadata3_new")}),
	}

	ck4ExpectedKeyOps := &keyOps{ 
		flag:     upsertVal,
		value:    []byte("value4_new"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata4": []byte("metadata4")}),
	}

	assert.Equal(t, ck1ExpectedKeyOps, txOps[ck1])
	assert.Equal(t, ck2ExpectedKeyOps, txOps[ck2])
	assert.Equal(t, ck3ExpectedKeyOps, txOps[ck3])
	assert.Equal(t, ck4ExpectedKeyOps, txOps[ck4])
}

func TestTxOpsPreparationPvtdataHashes(t *testing.T) {
	testDBEnv := privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	db := testDBEnv.GetDBHandle("TestDB")

	ck1, ck2, ck3, ck4 :=
		compositeKey{ns: "ns1", coll: "coll1", key: "key1"},
		compositeKey{ns: "ns1", coll: "coll1", key: "key2"},
		compositeKey{ns: "ns1", coll: "coll1", key: "key3"},
		compositeKey{ns: "ns1", coll: "coll1", key: "key4"}

	ck1Hash, ck2Hash, ck3Hash, ck4Hash :=
		compositeKey{ns: "ns1", coll: "coll1", key: string(util.ComputeStringHash("key1"))},
		compositeKey{ns: "ns1", coll: "coll1", key: string(util.ComputeStringHash("key2"))},
		compositeKey{ns: "ns1", coll: "coll1", key: string(util.ComputeStringHash("key3"))},
		compositeKey{ns: "ns1", coll: "coll1", key: string(util.ComputeStringHash("key4"))}

	updateBatch := privacyenabledstate.NewUpdateBatch()

	updateBatch.HashUpdates.Put(ck1.ns, ck1.coll, util.ComputeStringHash(ck1.key),
		util.ComputeStringHash("value1"), version.NewHeight(1, 1)) 

	updateBatch.HashUpdates.Put(ck2.ns, ck2.coll, util.ComputeStringHash(ck2.key),
		util.ComputeStringHash("value2"), version.NewHeight(1, 2)) 

	updateBatch.HashUpdates.PutValAndMetadata( 
		ck3.ns, ck3.coll, string(util.ComputeStringHash(ck3.key)),
		util.ComputeStringHash("value3"),
		testutilSerializedMetadata(t, map[string][]byte{"metadata3": []byte("metadata3")}),
		version.NewHeight(1, 3))

	updateBatch.HashUpdates.PutValAndMetadata( 
		ck4.ns, ck4.coll, string(util.ComputeStringHash(ck4.key)),
		util.ComputeStringHash("value4"),
		testutilSerializedMetadata(t, map[string][]byte{"metadata4": []byte("metadata4")}),
		version.NewHeight(1, 4))

	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(1, 2)) 

	precedingUpdates := NewPubAndHashUpdates()
	rwset := testutilBuildRwset( 
		t,
		map[compositeKey][]byte{
			ck1: []byte("value1_new"),
			ck2: []byte("value2_new"),
			ck4: []byte("value4_new"),
		},
		map[compositeKey]map[string][]byte{
			ck2: {"metadata2": []byte("metadata2_new")},
			ck3: {"metadata3": []byte("metadata3_new")},
		},
	)

	txOps, err := prepareTxOps(rwset, version.NewHeight(1, 2), precedingUpdates, db)
	assert.NoError(t, err)
	assert.Len(t, txOps, 4)

	ck1ExpectedKeyOps := &keyOps{ 
		flag:  upsertVal,
		value: util.ComputeStringHash("value1_new"),
	}

	ck2ExpectedKeyOps := &keyOps{ 
		flag:     upsertVal + metadataUpdate,
		value:    util.ComputeStringHash("value2_new"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata2": []byte("metadata2_new")}),
	}

	ck3ExpectedKeyOps := &keyOps{ 
		flag:     metadataUpdate,
		value:    util.ComputeStringHash("value3"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata3": []byte("metadata3_new")}),
	}

	ck4ExpectedKeyOps := &keyOps{ 
		flag:     upsertVal,
		value:    util.ComputeStringHash("value4_new"),
		metadata: testutilSerializedMetadata(t, map[string][]byte{"metadata4": []byte("metadata4")}),
	}

	assert.Equal(t, ck1ExpectedKeyOps, txOps[ck1Hash])
	assert.Equal(t, ck2ExpectedKeyOps, txOps[ck2Hash])
	assert.Equal(t, ck3ExpectedKeyOps, txOps[ck3Hash])
	assert.Equal(t, ck4ExpectedKeyOps, txOps[ck4Hash])
}

func testutilBuildRwset(t *testing.T,
	kvWrites map[compositeKey][]byte,
	metadataWrites map[compositeKey]map[string][]byte) *rwsetutil.TxRwSet {
	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	for kvwrite, val := range kvWrites {
		if kvwrite.coll == "" {
			rwsetBuilder.AddToWriteSet(kvwrite.ns, kvwrite.key, val)
		} else {
			rwsetBuilder.AddToPvtAndHashedWriteSet(kvwrite.ns, kvwrite.coll, kvwrite.key, val)
		}
	}

	for metadataWrite, metadataVal := range metadataWrites {
		if metadataWrite.coll == "" {
			rwsetBuilder.AddToMetadataWriteSet(metadataWrite.ns, metadataWrite.key, metadataVal)
		} else {
			rwsetBuilder.AddToHashedMetadataWriteSet(metadataWrite.ns, metadataWrite.coll, metadataWrite.key, metadataVal)
		}
	}
	return rwsetBuilder.GetTxReadWriteSet()
}

func testutilSerializedMetadata(t *testing.T, metadataMap map[string][]byte) []byte {
	metadataEntries := []*kvrwset.KVMetadataEntry{}
	for metadataK, metadataV := range metadataMap {
		metadataEntries = append(metadataEntries, &kvrwset.KVMetadataEntry{Name: metadataK, Value: metadataV})
	}
	metadataBytes, err := storageutil.SerializeMetadata(metadataEntries)
	assert.NoError(t, err)
	return metadataBytes
}
