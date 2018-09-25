/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privacyenabledstate

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
)

func TestMain(m *testing.M) {
	viper.Set("peer.fileSystemPath", "/tmp/blockchain/ledgertests/kvledger/txmgmt/privacyenabledstate")
	
	viper.Set("ledger.state.couchDBConfig.autoWarmIndexes", false)
	os.Exit(m.Run())
}

func TestBatch(t *testing.T) {
	batch := UpdateMap(make(map[string]nsBatch))
	v := version.NewHeight(1, 1)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				batch.Put(fmt.Sprintf("ns-%d", i), fmt.Sprintf("collection-%d", j), fmt.Sprintf("key-%d", k),
					[]byte(fmt.Sprintf("value-%d-%d-%d", i, j, k)), v)
			}
		}
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				vv := batch.Get(fmt.Sprintf("ns-%d", i), fmt.Sprintf("collection-%d", j), fmt.Sprintf("key-%d", k))
				assert.NotNil(t, vv)
				assert.Equal(t,
					&statedb.VersionedValue{Value: []byte(fmt.Sprintf("value-%d-%d-%d", i, j, k)), Version: v},
					vv)
			}
		}
	}
	assert.Nil(t, batch.Get("ns-1", "collection-1", "key-5"))
	assert.Nil(t, batch.Get("ns-1", "collection-5", "key-1"))
	assert.Nil(t, batch.Get("ns-5", "collection-1", "key-1"))
}

func TestHashBatchContains(t *testing.T) {
	batch := NewHashedUpdateBatch()
	batch.Put("ns1", "coll1", []byte("key1"), []byte("val1"), version.NewHeight(1, 1))
	assert.True(t, batch.Contains("ns1", "coll1", []byte("key1")))
	assert.False(t, batch.Contains("ns1", "coll1", []byte("key2")))
	assert.False(t, batch.Contains("ns1", "coll2", []byte("key1")))
	assert.False(t, batch.Contains("ns2", "coll1", []byte("key1")))

	batch.Delete("ns1", "coll1", []byte("deleteKey"), version.NewHeight(1, 1))
	assert.True(t, batch.Contains("ns1", "coll1", []byte("deleteKey")))
	assert.False(t, batch.Contains("ns1", "coll1", []byte("deleteKey1")))
	assert.False(t, batch.Contains("ns1", "coll2", []byte("deleteKey")))
	assert.False(t, batch.Contains("ns2", "coll1", []byte("deleteKey")))
}

func TestDB(t *testing.T) {
	for _, env := range testEnvs {
		t.Run(env.GetName(), func(t *testing.T) {
			testDB(t, env)
		})
	}
}

func testDB(t *testing.T, env TestEnv) {
	env.Init(t)
	defer env.Cleanup()
	db := env.GetDBHandle("test-ledger-id")

	updates := NewUpdateBatch()

	updates.PubUpdates.Put("ns1", "key1", []byte("value1"), version.NewHeight(1, 1))
	updates.PubUpdates.Put("ns1", "key2", []byte("value2"), version.NewHeight(1, 2))
	updates.PubUpdates.Put("ns2", "key3", []byte("value3"), version.NewHeight(1, 3))

	putPvtUpdates(t, updates, "ns1", "coll1", "key1", []byte("pvt_value1"), version.NewHeight(1, 4))
	putPvtUpdates(t, updates, "ns1", "coll1", "key2", []byte("pvt_value2"), version.NewHeight(1, 5))
	putPvtUpdates(t, updates, "ns2", "coll1", "key3", []byte("pvt_value3"), version.NewHeight(1, 6))
	db.ApplyPrivacyAwareUpdates(updates, version.NewHeight(2, 6))
	commonStorageDB := db.(*CommonStorageDB)
	bulkOptimizable, ok := commonStorageDB.VersionedDB.(statedb.BulkOptimizable)
	if ok {
		bulkOptimizable.ClearCachedVersions()
	}

	vv, err := db.GetState("ns1", "key1")
	assert.NoError(t, err)
	assert.Equal(t, &statedb.VersionedValue{Value: []byte("value1"), Version: version.NewHeight(1, 1)}, vv)

	vv, err = db.GetPrivateData("ns1", "coll1", "key1")
	assert.NoError(t, err)
	assert.Equal(t, &statedb.VersionedValue{Value: []byte("pvt_value1"), Version: version.NewHeight(1, 4)}, vv)

	vv, err = db.GetValueHash("ns1", "coll1", util.ComputeStringHash("key1"))
	assert.NoError(t, err)
	assert.Equal(t, &statedb.VersionedValue{Value: util.ComputeStringHash("pvt_value1"), Version: version.NewHeight(1, 4)}, vv)

	committedVersion, err := db.GetKeyHashVersion("ns1", "coll1", util.ComputeStringHash("key1"))
	assert.NoError(t, err)
	assert.Equal(t, version.NewHeight(1, 4), committedVersion)

	updates = NewUpdateBatch()
	updates.PubUpdates.Delete("ns1", "key1", version.NewHeight(2, 7))
	deletePvtUpdates(t, updates, "ns1", "coll1", "key1", version.NewHeight(2, 7))
	db.ApplyPrivacyAwareUpdates(updates, version.NewHeight(2, 7))

	vv, err = db.GetState("ns1", "key1")
	assert.NoError(t, err)
	assert.Nil(t, vv)

	vv, err = db.GetPrivateData("ns1", "coll1", "key1")
	assert.NoError(t, err)
	assert.Nil(t, vv)

	vv, err = db.GetValueHash("ns1", "coll1", util.ComputeStringHash("key1"))
	assert.Nil(t, vv)
}

func TestGetStateMultipleKeys(t *testing.T) {
	for _, env := range testEnvs {
		t.Run(env.GetName(), func(t *testing.T) {
			testGetStateMultipleKeys(t, env)
		})
	}
}

func testGetStateMultipleKeys(t *testing.T, env TestEnv) {
	env.Init(t)
	defer env.Cleanup()
	db := env.GetDBHandle("test-ledger-id")

	updates := NewUpdateBatch()

	updates.PubUpdates.Put("ns1", "key1", []byte("value1"), version.NewHeight(1, 1))
	updates.PubUpdates.Put("ns1", "key2", []byte("value2"), version.NewHeight(1, 2))
	updates.PubUpdates.Put("ns1", "key3", []byte("value3"), version.NewHeight(1, 3))

	putPvtUpdates(t, updates, "ns1", "coll1", "key1", []byte("pvt_value1"), version.NewHeight(1, 4))
	putPvtUpdates(t, updates, "ns1", "coll1", "key2", []byte("pvt_value2"), version.NewHeight(1, 5))
	putPvtUpdates(t, updates, "ns1", "coll1", "key3", []byte("pvt_value3"), version.NewHeight(1, 6))
	db.ApplyPrivacyAwareUpdates(updates, version.NewHeight(2, 6))

	versionedVals, err := db.GetStateMultipleKeys("ns1", []string{"key1", "key3"})
	assert.NoError(t, err)
	assert.Equal(t,
		[]*statedb.VersionedValue{
			{Value: []byte("value1"), Version: version.NewHeight(1, 1)},
			{Value: []byte("value3"), Version: version.NewHeight(1, 3)},
		},
		versionedVals)

	pvtVersionedVals, err := db.GetPrivateDataMultipleKeys("ns1", "coll1", []string{"key1", "key3"})
	assert.NoError(t, err)
	assert.Equal(t,
		[]*statedb.VersionedValue{
			{Value: []byte("pvt_value1"), Version: version.NewHeight(1, 4)},
			{Value: []byte("pvt_value3"), Version: version.NewHeight(1, 6)},
		},
		pvtVersionedVals)
}

func TestGetStateRangeScanIterator(t *testing.T) {
	for _, env := range testEnvs {
		t.Run(env.GetName(), func(t *testing.T) {
			testGetStateRangeScanIterator(t, env)
		})
	}
}

func testGetStateRangeScanIterator(t *testing.T, env TestEnv) {
	env.Init(t)
	defer env.Cleanup()
	db := env.GetDBHandle("test-ledger-id")

	updates := NewUpdateBatch()

	updates.PubUpdates.Put("ns1", "key1", []byte("value1"), version.NewHeight(1, 1))
	updates.PubUpdates.Put("ns1", "key2", []byte("value2"), version.NewHeight(1, 2))
	updates.PubUpdates.Put("ns1", "key3", []byte("value3"), version.NewHeight(1, 3))
	updates.PubUpdates.Put("ns1", "key4", []byte("value4"), version.NewHeight(1, 4))
	updates.PubUpdates.Put("ns2", "key5", []byte("value5"), version.NewHeight(1, 5))
	updates.PubUpdates.Put("ns2", "key6", []byte("value6"), version.NewHeight(1, 6))
	updates.PubUpdates.Put("ns3", "key7", []byte("value7"), version.NewHeight(1, 7))

	putPvtUpdates(t, updates, "ns1", "coll1", "key1", []byte("pvt_value1"), version.NewHeight(1, 1))
	putPvtUpdates(t, updates, "ns1", "coll1", "key2", []byte("pvt_value2"), version.NewHeight(1, 2))
	putPvtUpdates(t, updates, "ns1", "coll1", "key3", []byte("pvt_value3"), version.NewHeight(1, 3))
	putPvtUpdates(t, updates, "ns1", "coll1", "key4", []byte("pvt_value4"), version.NewHeight(1, 4))
	putPvtUpdates(t, updates, "ns2", "coll1", "key5", []byte("pvt_value5"), version.NewHeight(1, 5))
	putPvtUpdates(t, updates, "ns2", "coll1", "key6", []byte("pvt_value6"), version.NewHeight(1, 6))
	putPvtUpdates(t, updates, "ns3", "coll1", "key7", []byte("pvt_value7"), version.NewHeight(1, 7))
	db.ApplyPrivacyAwareUpdates(updates, version.NewHeight(2, 7))

	itr1, _ := db.GetStateRangeScanIterator("ns1", "key1", "")
	testItr(t, itr1, []string{"key1", "key2", "key3", "key4"})

	itr2, _ := db.GetStateRangeScanIterator("ns1", "key2", "key3")
	testItr(t, itr2, []string{"key2"})

	itr3, _ := db.GetStateRangeScanIterator("ns1", "", "")
	testItr(t, itr3, []string{"key1", "key2", "key3", "key4"})

	itr4, _ := db.GetStateRangeScanIterator("ns2", "", "")
	testItr(t, itr4, []string{"key5", "key6"})

	pvtItr1, _ := db.GetPrivateDataRangeScanIterator("ns1", "coll1", "key1", "")
	testItr(t, pvtItr1, []string{"key1", "key2", "key3", "key4"})

	pvtItr2, _ := db.GetPrivateDataRangeScanIterator("ns1", "coll1", "key2", "key3")
	testItr(t, pvtItr2, []string{"key2"})

	pvtItr3, _ := db.GetPrivateDataRangeScanIterator("ns1", "coll1", "", "")
	testItr(t, pvtItr3, []string{"key1", "key2", "key3", "key4"})

	pvtItr4, _ := db.GetPrivateDataRangeScanIterator("ns2", "coll1", "", "")
	testItr(t, pvtItr4, []string{"key5", "key6"})
}

func TestQueryOnCouchDB(t *testing.T) {
	for _, env := range testEnvs {
		_, ok := env.(*CouchDBCommonStorageTestEnv)
		if !ok {
			continue
		}
		t.Run(env.GetName(), func(t *testing.T) {
			testQueryOnCouchDB(t, env)
		})
	}
}

func testQueryOnCouchDB(t *testing.T, env TestEnv) {
	env.Init(t)
	defer env.Cleanup()
	db := env.GetDBHandle("test-ledger-id")
	updates := NewUpdateBatch()

	jsonValues := []string{
		`{"asset_name": "marble1", "color": "blue", "size": 1, "owner": "tom"}`,
		`{"asset_name": "marble2","color": "blue","size": 2,"owner": "jerry"}`,
		`{"asset_name": "marble3","color": "blue","size": 3,"owner": "fred"}`,
		`{"asset_name": "marble4","color": "blue","size": 4,"owner": "martha"}`,
		`{"asset_name": "marble5","color": "blue","size": 5,"owner": "fred"}`,
		`{"asset_name": "marble6","color": "blue","size": 6,"owner": "elaine"}`,
		`{"asset_name": "marble7","color": "blue","size": 7,"owner": "fred"}`,
		`{"asset_name": "marble8","color": "blue","size": 8,"owner": "elaine"}`,
		`{"asset_name": "marble9","color": "green","size": 9,"owner": "fred"}`,
		`{"asset_name": "marble10","color": "green","size": 10,"owner": "mary"}`,
		`{"asset_name": "marble11","color": "cyan","size": 1000007,"owner": "joe"}`,
	}

	for i, jsonValue := range jsonValues {
		updates.PubUpdates.Put("ns1", testKey(i), []byte(jsonValue), version.NewHeight(1, uint64(i)))
		updates.PubUpdates.Put("ns2", testKey(i), []byte(jsonValue), version.NewHeight(1, uint64(i)))
		putPvtUpdates(t, updates, "ns1", "coll1", testKey(i), []byte(jsonValue), version.NewHeight(1, uint64(i)))
		putPvtUpdates(t, updates, "ns2", "coll1", testKey(i), []byte(jsonValue), version.NewHeight(1, uint64(i)))
	}
	db.ApplyPrivacyAwareUpdates(updates, version.NewHeight(1, 11))

	
	itr, err := db.ExecuteQuery("ns1", `{"selector":{"owner":"jerry"}}`)
	assert.NoError(t, err)
	testQueryItr(t, itr, []string{testKey(1)}, []string{"jerry"})

	
	itr, err = db.ExecuteQuery("ns2", `{"selector":{"owner":"jerry"}}`)
	assert.NoError(t, err)
	testQueryItr(t, itr, []string{testKey(1)}, []string{"jerry"})

	
	itr, err = db.ExecuteQueryOnPrivateData("ns1", "coll1", `{"selector":{"owner":"jerry"}}`)
	assert.NoError(t, err)
	testQueryItr(t, itr, []string{testKey(1)}, []string{"jerry"})

	
	itr, err = db.ExecuteQueryOnPrivateData("ns2", "coll1", `{"selector":{"owner":"jerry"}}`)
	assert.NoError(t, err)
	testQueryItr(t, itr, []string{testKey(1)}, []string{"jerry"})

	
	itr, err = db.ExecuteQueryOnPrivateData("ns1", "coll1", "this is an invalid query string")
	assert.Error(t, err, "Should have received an error for invalid query string")

	
	itr, err = db.ExecuteQueryOnPrivateData("ns1", "coll1", `{"selector":{"owner":"not_a_valid_name"}}`)
	assert.NoError(t, err)
	testQueryItr(t, itr, []string{}, []string{})

	
	itr, err = db.ExecuteQueryOnPrivateData("ns1", "coll1", `{"selector":{"color":"green","$or":[{"owner":"fred"},{"owner":"mary"}]}}`)
	assert.NoError(t, err)
	testQueryItr(t, itr, []string{testKey(8), testKey(9)}, []string{"green"}, []string{"green"})

	
	
	itr, err = db.ExecuteQueryOnPrivateData("ns2", "coll1", `{"selector":{"$and":[{"size":{"$eq": 1000007}}]}}`)
	assert.NoError(t, err)
	testQueryItr(t, itr, []string{testKey(10)}, []string{"joe", "1000007"})
}

func TestLongDBNameOnCouchDB(t *testing.T) {
	for _, env := range testEnvs {
		_, ok := env.(*CouchDBCommonStorageTestEnv)
		if !ok {
			continue
		}
		t.Run(env.GetName(), func(t *testing.T) {
			testLongDBNameOnCouchDB(t, env)
		})
	}
}

func testLongDBNameOnCouchDB(t *testing.T, env TestEnv) {
	env.Init(t)
	defer env.Cleanup()

	
	
	db := env.GetDBHandle("w1coaii9ck3l8red6a5cf3rwbe1b4wvbzcrrfl7samu7px8b9gf-4hft7wrgdmzzjj9ure4cbffucaj78nbj9ej.kvl3bus1iq1qir9xlhb8a1wipuksgs3g621elzy1prr658087exwrhp-y4j55o9cld242v--oeh3br1g7m8d6l8jobn.y42cgjt1.u1ik8qxnv4ohh9kr2w2zc8hqir5u4ev23s7jygrg....s7.ohp-5bcxari8nji")

	updates := NewUpdateBatch()

	
	ns := "wMCnSXiV9YoIqNQyNvFVTdM8XnUtvrOFFIWsKelmP5NEszmNLl8YhtOKbFu3P_NgwgsYF8PsfwjYCD8f1XRpANQLoErDHwLlweryqXeJ6vzT2x0pS_GwSx0m6tBI0zOmHQOq_2De8A87x6zUOPwufC2T6dkidFxiuq8Sey2-5vUo_iNKCij3WTeCnKx78PUIg_U1gp4_0KTvYVtRBRvH0kz5usizBxPaiFu3TPhB9XLviScvdUVSbSYJ0Z"
	coll := "vWjtfSTXVK8WJus5s6zWoMIciXd7qHRZIusF9SkOS6m8XuHCiJDE9cCRuVerq22Na8qBL2ywDGFpVMIuzfyEXLjeJb0mMuH4cwewT6r1INOTOSYwrikwOLlT_fl0V1L7IQEwUBB8WCvRqSdj6j5-E5aGul_pv_0UeCdwWiyA_GrZmP7ocLzfj2vP8btigrajqdH-irLO2ydEjQUAvf8fiuxru9la402KmKRy457GgI98UHoUdqV3f3FCdR"

	updates.PubUpdates.Put(ns, "key1", []byte("value1"), version.NewHeight(1, 1))
	updates.PvtUpdates.Put(ns, coll, "key1", []byte("pvt_value"), version.NewHeight(1, 2))
	updates.HashUpdates.Put(ns, coll, util.ComputeStringHash("key1"), util.ComputeHash([]byte("pvt_value")), version.NewHeight(1, 2))

	db.ApplyPrivacyAwareUpdates(updates, version.NewHeight(2, 6))

	vv, err := db.GetState(ns, "key1")
	assert.NoError(t, err)
	assert.Equal(t, &statedb.VersionedValue{Value: []byte("value1"), Version: version.NewHeight(1, 1)}, vv)

	vv, err = db.GetPrivateData(ns, coll, "key1")
	assert.NoError(t, err)
	assert.Equal(t, &statedb.VersionedValue{Value: []byte("pvt_value"), Version: version.NewHeight(1, 2)}, vv)
}

func testItr(t *testing.T, itr statedb.ResultsIterator, expectedKeys []string) {
	defer itr.Close()
	for _, expectedKey := range expectedKeys {
		queryResult, _ := itr.Next()
		vkv := queryResult.(*statedb.VersionedKV)
		key := vkv.Key
		assert.Equal(t, expectedKey, key)
	}
	last, err := itr.Next()
	assert.NoError(t, err)
	assert.Nil(t, last)
}

func testQueryItr(t *testing.T, itr statedb.ResultsIterator, expectedKeys []string, expectedValStrs ...[]string) {
	defer itr.Close()
	for i, expectedKey := range expectedKeys {
		queryResult, _ := itr.Next()
		vkv := queryResult.(*statedb.VersionedKV)
		key := vkv.Key
		valStr := string(vkv.Value)
		assert.Equal(t, expectedKey, key)
		for _, expectedValStr := range expectedValStrs[i] {
			assert.Contains(t, valStr, expectedValStr)
		}
	}
	last, err := itr.Next()
	assert.NoError(t, err)
	assert.Nil(t, last)
}

func testKey(i int) string {
	return fmt.Sprintf("key%d", i)
}

func TestCompositeKeyMap(t *testing.T) {
	b := NewPvtUpdateBatch()
	b.Put("ns1", "coll1", "key1", []byte("testVal1"), nil)
	b.Delete("ns1", "coll2", "key2", nil)
	b.Put("ns2", "coll1", "key1", []byte("testVal3"), nil)
	b.Put("ns2", "coll2", "key2", []byte("testVal4"), nil)
	m := b.ToCompositeKeyMap()
	assert.Len(t, m, 4)
	vv, ok := m[PvtdataCompositeKey{"ns1", "coll1", "key1"}]
	assert.True(t, ok)
	assert.Equal(t, []byte("testVal1"), vv.Value)
	vv, ok = m[PvtdataCompositeKey{"ns1", "coll2", "key2"}]
	assert.Nil(t, vv.Value)
	assert.True(t, ok)
	_, ok = m[PvtdataCompositeKey{"ns2", "coll1", "key1"}]
	assert.True(t, ok)
	_, ok = m[PvtdataCompositeKey{"ns2", "coll2", "key2"}]
	assert.True(t, ok)
	_, ok = m[PvtdataCompositeKey{"ns2", "coll1", "key8888"}]
	assert.False(t, ok)
}

func TestHandleChainCodeDeployOnCouchDB(t *testing.T) {
	for _, env := range testEnvs {
		_, ok := env.(*CouchDBCommonStorageTestEnv)
		if !ok {
			continue
		}
		t.Run(env.GetName(), func(t *testing.T) {
			testHandleChainCodeDeploy(t, env)
		})
	}
}

func createCollectionConfig(collectionName string) *common.CollectionConfig {
	return &common.CollectionConfig{
		Payload: &common.CollectionConfig_StaticCollectionConfig{
			StaticCollectionConfig: &common.StaticCollectionConfig{
				Name:              collectionName,
				MemberOrgsPolicy:  nil,
				RequiredPeerCount: 0,
				MaximumPeerCount:  0,
				BlockToLive:       0,
			},
		},
	}
}

func testHandleChainCodeDeploy(t *testing.T, env TestEnv) {
	env.Init(t)
	defer env.Cleanup()
	db := env.GetDBHandle("test-handle-chaincode-deploy")

	coll1 := createCollectionConfig("collectionMarbles")
	ccp := &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1}}
	ccpBytes, err := proto.Marshal(ccp)
	assert.NoError(t, err)
	assert.NotNil(t, ccpBytes)

	chaincodeDef := &cceventmgmt.ChaincodeDefinition{Name: "ns1", Hash: nil, Version: "", CollectionConfigs: ccpBytes}

	commonStorageDB := db.(*CommonStorageDB)

	
	dbArtifactsTarBytes := testutil.CreateTarBytesForTest(
		[]*testutil.TarFileEntry{
			{Name: "META-INF/statedb/couchdb/indexes/indexColorSortName.json", Body: `{"index":{"fields":[{"color":"desc"}]},"ddoc":"indexColorSortName","name":"indexColorSortName","type":"json"}`},
			{Name: "META-INF/statedb/couchdb/indexes/indexSizeSortName.json", Body: `{"index":{"fields":[{"size":"desc"}]},"ddoc":"indexSizeSortName","name":"indexSizeSortName","type":"json"}`},
			{Name: "META-INF/statedb/couchdb/collections/collectionMarbles/indexes/indexCollMarbles.json", Body: `{"index":{"fields":["docType","owner"]},"ddoc":"indexCollectionMarbles", "name":"indexCollectionMarbles","type":"json"}`},
			{Name: "META-INF/statedb/couchdb/collections/collectionMarblesPrivateDetails/indexes/indexCollPrivDetails.json", Body: `{"index":{"fields":["docType","price"]},"ddoc":"indexPrivateDetails", "name":"indexPrivateDetails","type":"json"}`},
		},
	)

	
	fileEntries, err := ccprovider.ExtractFileEntries(dbArtifactsTarBytes, "couchdb")
	assert.NoError(t, err)

	
	assert.Len(t, fileEntries, 3)

	
	assert.Len(t, fileEntries["META-INF/statedb/couchdb/indexes"], 2)

	
	assert.Len(t, fileEntries["META-INF/statedb/couchdb/collections/collectionMarbles/indexes"], 1)

	
	expectedJSON := []byte(`{"index":{"fields":["docType","owner"]},"ddoc":"indexCollectionMarbles", "name":"indexCollectionMarbles","type":"json"}`)
	actualJSON := fileEntries["META-INF/statedb/couchdb/collections/collectionMarbles/indexes"][0].FileContent
	assert.Equal(t, expectedJSON, actualJSON)

	
	
	err = commonStorageDB.HandleChaincodeDeploy(chaincodeDef, dbArtifactsTarBytes)
	assert.NoError(t, err)

	coll2 := createCollectionConfig("collectionMarblesPrivateDetails")
	ccp = &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1, coll2}}
	ccpBytes, err = proto.Marshal(ccp)
	assert.NoError(t, err)
	assert.NotNil(t, ccpBytes)

	chaincodeDef = &cceventmgmt.ChaincodeDefinition{Name: "ns1", Hash: nil, Version: "", CollectionConfigs: ccpBytes}

	
	
	
	
	err = commonStorageDB.HandleChaincodeDeploy(chaincodeDef, dbArtifactsTarBytes)
	assert.NoError(t, err)

	chaincodeDef = &cceventmgmt.ChaincodeDefinition{Name: "ns1", Hash: nil, Version: "", CollectionConfigs: nil}

	
	
	
	
	err = commonStorageDB.HandleChaincodeDeploy(chaincodeDef, dbArtifactsTarBytes)
	assert.NoError(t, err)

	
	err = commonStorageDB.HandleChaincodeDeploy(chaincodeDef, nil)
	assert.NoError(t, err)

	
	err = commonStorageDB.HandleChaincodeDeploy(chaincodeDef, []byte(`This is a really bad tar file`))
	assert.NoError(t, err, "Error should not have been thrown for a bad tar file")

	
	err = commonStorageDB.HandleChaincodeDeploy(nil, dbArtifactsTarBytes)
	assert.Error(t, err, "Error should have been thrown for a nil chaincodeDefinition")

	
	badSyntaxFileContent := `{"index":{"fields": This is a bad json}`
	dbArtifactsTarBytes = testutil.CreateTarBytesForTest(
		[]*testutil.TarFileEntry{
			{Name: "META-INF/statedb/couchdb/indexes/indexSizeSortName.json", Body: `{"index":{"fields":[{"size":"desc"}]},"ddoc":"indexSizeSortName","name":"indexSizeSortName","type":"json"}`},
			{Name: "META-INF/statedb/couchdb/indexes/badSyntax.json", Body: badSyntaxFileContent},
		},
	)

	
	fileEntries, err = ccprovider.ExtractFileEntries(dbArtifactsTarBytes, "couchdb")
	assert.NoError(t, err)

	
	assert.Len(t, fileEntries, 1)

	err = commonStorageDB.HandleChaincodeDeploy(chaincodeDef, dbArtifactsTarBytes)
	assert.NoError(t, err)

}

func TestMetadataRetrieval(t *testing.T) {
	for _, env := range testEnvs {
		t.Run(env.GetName(), func(t *testing.T) {
			testMetadataRetrieval(t, env)
		})
	}
}

func testMetadataRetrieval(t *testing.T, env TestEnv) {
	env.Init(t)
	defer env.Cleanup()
	db := env.GetDBHandle("test-ledger-id")

	updates := NewUpdateBatch()
	updates.PubUpdates.PutValAndMetadata("ns1", "key1", []byte("value1"), []byte("metadata1"), version.NewHeight(1, 1))
	updates.PubUpdates.PutValAndMetadata("ns1", "key2", []byte("value2"), nil, version.NewHeight(1, 2))
	updates.PubUpdates.PutValAndMetadata("ns2", "key3", []byte("value3"), nil, version.NewHeight(1, 3))

	putPvtUpdatesWithMetadata(t, updates, "ns1", "coll1", "key1", []byte("pvt_value1"), []byte("metadata1"), version.NewHeight(1, 4))
	putPvtUpdatesWithMetadata(t, updates, "ns1", "coll1", "key2", []byte("pvt_value2"), nil, version.NewHeight(1, 5))
	putPvtUpdatesWithMetadata(t, updates, "ns2", "coll1", "key3", []byte("pvt_value3"), nil, version.NewHeight(1, 6))
	db.ApplyPrivacyAwareUpdates(updates, version.NewHeight(2, 6))

	vm, _ := db.GetStateMetadata("ns1", "key1")
	assert.Equal(t, vm, []byte("metadata1"))
	vm, _ = db.GetStateMetadata("ns1", "key2")
	assert.Nil(t, vm)
	vm, _ = db.GetStateMetadata("ns2", "key3")
	assert.Nil(t, vm)

	vm, _ = db.GetPrivateDataMetadataByHash("ns1", "coll1", util.ComputeStringHash("key1"))
	assert.Equal(t, vm, []byte("metadata1"))
	vm, _ = db.GetPrivateDataMetadataByHash("ns1", "coll1", util.ComputeStringHash("key2"))
	assert.Nil(t, vm)
	vm, _ = db.GetPrivateDataMetadataByHash("ns2", "coll1", util.ComputeStringHash("key3"))
	assert.Nil(t, vm)
}

func putPvtUpdates(t *testing.T, updates *UpdateBatch, ns, coll, key string, value []byte, ver *version.Height) {
	updates.PvtUpdates.Put(ns, coll, key, value, ver)
	updates.HashUpdates.Put(ns, coll, util.ComputeStringHash(key), util.ComputeHash(value), ver)
}

func putPvtUpdatesWithMetadata(t *testing.T, updates *UpdateBatch, ns, coll, key string, value []byte, metadata []byte, ver *version.Height) {
	updates.PvtUpdates.Put(ns, coll, key, value, ver)
	updates.HashUpdates.PutValHashAndMetadata(ns, coll, util.ComputeStringHash(key), util.ComputeHash(value), metadata, ver)
}

func deletePvtUpdates(t *testing.T, updates *UpdateBatch, ns, coll, key string, ver *version.Height) {
	updates.PvtUpdates.Delete(ns, coll, key, ver)
	updates.HashUpdates.Delete(ns, coll, util.ComputeStringHash(key), ver)
}
