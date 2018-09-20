/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statecouchdb

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb/commontests"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	ledgertestutil "github.com/mcc-github/blockchain/core/ledger/testutil"
	"github.com/mcc-github/blockchain/integration/runner"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	flogging.SetModuleLevel("statecouchdb", "debug")
	flogging.SetModuleLevel("couchdb", "debug")
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	
	ledgertestutil.SetupCoreYAMLConfig()
	viper.Set("peer.fileSystemPath", "/tmp/blockchain/ledgertests/kvledger/txmgmt/statedb/statecouchdb")

	
	couchAddress, cleanup := couchDBSetup()
	defer cleanup()
	viper.Set("ledger.state.stateDatabase", "CouchDB")
	defer viper.Set("ledger.state.stateDatabase", "goleveldb")

	viper.Set("ledger.state.couchDBConfig.couchDBAddress", couchAddress)
	
	
	viper.Set("ledger.state.couchDBConfig.username", "")
	viper.Set("ledger.state.couchDBConfig.password", "")
	viper.Set("ledger.state.couchDBConfig.maxRetries", 3)
	viper.Set("ledger.state.couchDBConfig.maxRetriesOnStartup", 10)
	viper.Set("ledger.state.couchDBConfig.requestTimeout", time.Second*35)
	flogging.SetModuleLevel("statecouchdb", "debug")
	
	return m.Run()
}

func couchDBSetup() (addr string, cleanup func()) {
	externalCouch, set := os.LookupEnv("COUCHDB_ADDR")
	if set {
		return externalCouch, func() {}
	}

	couchDB := &runner.CouchDB{}
	if err := couchDB.Start(); err != nil {
		err := fmt.Errorf("failed to start couchDB: %s", err)
		panic(err)
	}
	return couchDB.Address(), func() { couchDB.Stop() }
}

func TestBasicRW(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testbasicrw_")
	env.Cleanup("testbasicrw_ns")
	env.Cleanup("testbasicrw_ns1")
	env.Cleanup("testbasicrw_ns2")
	defer env.Cleanup("testbasicrw_")
	defer env.Cleanup("testbasicrw_ns")
	defer env.Cleanup("testbasicrw_ns1")
	defer env.Cleanup("testbasicrw_ns2")
	commontests.TestBasicRW(t, env.DBProvider)

}

func TestMultiDBBasicRW(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testmultidbbasicrw_")
	env.Cleanup("testmultidbbasicrw_ns1")
	env.Cleanup("testmultidbbasicrw2_")
	env.Cleanup("testmultidbbasicrw2_ns1")
	defer env.Cleanup("testmultidbbasicrw_")
	defer env.Cleanup("testmultidbbasicrw_ns1")
	defer env.Cleanup("testmultidbbasicrw2_")
	defer env.Cleanup("testmultidbbasicrw2_ns1")
	commontests.TestMultiDBBasicRW(t, env.DBProvider)

}

func TestDeletes(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testdeletes_")
	env.Cleanup("testdeletes_ns")
	defer env.Cleanup("testdeletes_")
	defer env.Cleanup("testdeletes_ns")
	commontests.TestDeletes(t, env.DBProvider)
}

func TestIterator(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testiterator_")
	env.Cleanup("testiterator_ns1")
	env.Cleanup("testiterator_ns2")
	env.Cleanup("testiterator_ns3")
	defer env.Cleanup("testiterator_")
	defer env.Cleanup("testiterator_ns1")
	defer env.Cleanup("testiterator_ns2")
	defer env.Cleanup("testiterator_ns3")
	commontests.TestIterator(t, env.DBProvider)
}



func TestQuery(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testquery_")
	env.Cleanup("testquery_ns1")
	env.Cleanup("testquery_ns2")
	env.Cleanup("testquery_ns3")
	defer env.Cleanup("testquery_")
	defer env.Cleanup("testquery_ns1")
	defer env.Cleanup("testquery_ns2")
	defer env.Cleanup("testquery_ns3")
	commontests.TestQuery(t, env.DBProvider)
}

func TestGetStateMultipleKeys(t *testing.T) {

	env := NewTestVDBEnv(t)
	env.Cleanup("testgetmultiplekeys_")
	env.Cleanup("testgetmultiplekeys_ns1")
	env.Cleanup("testgetmultiplekeys_ns2")
	defer env.Cleanup("testgetmultiplekeys_")
	defer env.Cleanup("testgetmultiplekeys_ns1")
	defer env.Cleanup("testgetmultiplekeys_ns2")
	commontests.TestGetStateMultipleKeys(t, env.DBProvider)
}

func TestGetVersion(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testgetversion_")
	env.Cleanup("testgetversion_ns")
	env.Cleanup("testgetversion_ns2")
	defer env.Cleanup("testgetversion_")
	defer env.Cleanup("testgetversion_ns")
	defer env.Cleanup("testgetversion_ns2")
	commontests.TestGetVersion(t, env.DBProvider)
}

func TestSmallBatchSize(t *testing.T) {
	viper.Set("ledger.state.couchDBConfig.maxBatchUpdateSize", 2)
	env := NewTestVDBEnv(t)
	env.Cleanup("testsmallbatchsize_")
	env.Cleanup("testsmallbatchsize_ns1")
	defer env.Cleanup("testsmallbatchsize_")
	defer env.Cleanup("testsmallbatchsize_ns1")
	defer viper.Set("ledger.state.couchDBConfig.maxBatchUpdateSize", 1000)
	commontests.TestSmallBatchSize(t, env.DBProvider)
}

func TestBatchRetry(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testbatchretry_")
	env.Cleanup("testbatchretry_ns")
	env.Cleanup("testbatchretry_ns1")
	defer env.Cleanup("testbatchretry_")
	defer env.Cleanup("testbatchretry_ns")
	defer env.Cleanup("testbatchretry_ns1")
	commontests.TestBatchWithIndividualRetry(t, env.DBProvider)
}

func TestValueAndMetadataWrites(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testvalueandmetadata_")
	env.Cleanup("testvalueandmetadata_ns1")
	env.Cleanup("testvalueandmetadata_ns2")
	defer env.Cleanup("testvalueandmetadata_")
	defer env.Cleanup("testvalueandmetadata_ns1")
	defer env.Cleanup("testvalueandmetadata_ns2")
	commontests.TestValueAndMetadataWrites(t, env.DBProvider)
}

func TestPaginatedRangeQuery(t *testing.T) {
	env := NewTestVDBEnv(t)
	env.Cleanup("testpaginatedrangequery_")
	env.Cleanup("testpaginatedrangequery_ns1")
	defer env.Cleanup("testpaginatedrangequery_")
	defer env.Cleanup("testpaginatedrangequery_ns1")
	commontests.TestPaginatedRangeQuery(t, env.DBProvider)
}


func TestUtilityFunctions(t *testing.T) {

	env := NewTestVDBEnv(t)
	env.Cleanup("testutilityfunctions_")
	defer env.Cleanup("testutilityfunctions_")

	db, err := env.DBProvider.GetDBHandle("testutilityfunctions")
	testutil.AssertNoError(t, err, "")

	
	byteKeySupported := db.BytesKeySuppoted()
	testutil.AssertEquals(t, byteKeySupported, false)

	
	err = db.ValidateKeyValue("testKey", []byte("Some random bytes"))
	testutil.AssertNil(t, err)

	
	err = db.ValidateKeyValue(string([]byte{0xff, 0xfe, 0xfd}), []byte("Some random bytes"))
	testutil.AssertError(t, err, "ValidateKey should have thrown an error for an invalid utf-8 string")

	reservedFields := []string{"~version", "_id", "_test"}

	
	
	for _, reservedField := range reservedFields {
		testVal := fmt.Sprintf(`{"%s":"dummyVal"}`, reservedField)
		err = db.ValidateKeyValue("testKey", []byte(testVal))
		testutil.AssertError(t, err, fmt.Sprintf(
			"ValidateKey should have thrown an error for a json value %s, as contains one of the reserved fields", testVal))
	}

	
	
	for _, reservedField := range reservedFields {
		testVal := fmt.Sprintf(`{"data.%s":"dummyVal"}`, reservedField)
		err = db.ValidateKeyValue("testKey", []byte(testVal))
		testutil.AssertNoError(t, err, fmt.Sprintf(
			"ValidateKey should not have thrown an error the json value %s since the reserved field was not at the top level", testVal))
	}

	
	err = db.ValidateKeyValue("_testKey", []byte("testValue"))
	testutil.AssertError(t, err, "ValidateKey should have thrown an error for a key that begins with an underscore")

}


func TestInvalidJSONFields(t *testing.T) {

	env := NewTestVDBEnv(t)
	env.Cleanup("testinvalidfields_")
	defer env.Cleanup("testinvalidfields_")

	db, err := env.DBProvider.GetDBHandle("testinvalidfields")
	testutil.AssertNoError(t, err, "")

	db.Open()
	defer db.Close()

	batch := statedb.NewUpdateBatch()
	jsonValue1 := `{"_id":"key1","asset_name":"marble1","color":"blue","size":1,"owner":"tom"}`
	batch.Put("ns1", "key1", []byte(jsonValue1), version.NewHeight(1, 1))

	savePoint := version.NewHeight(1, 2)
	err = db.ApplyUpdates(batch, savePoint)
	testutil.AssertError(t, err, "Invalid field _id should have thrown an error")

	batch = statedb.NewUpdateBatch()
	jsonValue1 = `{"_rev":"rev1","asset_name":"marble1","color":"blue","size":1,"owner":"tom"}`
	batch.Put("ns1", "key1", []byte(jsonValue1), version.NewHeight(1, 1))

	savePoint = version.NewHeight(1, 2)
	err = db.ApplyUpdates(batch, savePoint)
	testutil.AssertError(t, err, "Invalid field _rev should have thrown an error")

	batch = statedb.NewUpdateBatch()
	jsonValue1 = `{"_deleted":"true","asset_name":"marble1","color":"blue","size":1,"owner":"tom"}`
	batch.Put("ns1", "key1", []byte(jsonValue1), version.NewHeight(1, 1))

	savePoint = version.NewHeight(1, 2)
	err = db.ApplyUpdates(batch, savePoint)
	testutil.AssertError(t, err, "Invalid field _deleted should have thrown an error")

	batch = statedb.NewUpdateBatch()
	jsonValue1 = `{"~version":"v1","asset_name":"marble1","color":"blue","size":1,"owner":"tom"}`
	batch.Put("ns1", "key1", []byte(jsonValue1), version.NewHeight(1, 1))

	savePoint = version.NewHeight(1, 2)
	err = db.ApplyUpdates(batch, savePoint)
	testutil.AssertError(t, err, "Invalid field ~version should have thrown an error")
}

func TestDebugFunctions(t *testing.T) {

	
	
	loadKeys := []*statedb.CompositeKey{}
	
	compositeKey := statedb.CompositeKey{Namespace: "ns", Key: "key3"}
	loadKeys = append(loadKeys, &compositeKey)
	compositeKey = statedb.CompositeKey{Namespace: "ns", Key: "key4"}
	loadKeys = append(loadKeys, &compositeKey)
	testutil.AssertEquals(t, printCompositeKeys(loadKeys), "[ns,key4],[ns,key4]")

}

func TestHandleChaincodeDeploy(t *testing.T) {

	env := NewTestVDBEnv(t)
	env.Cleanup("testinit_")
	env.Cleanup("testinit_ns1")
	env.Cleanup("testinit_ns2")
	defer env.Cleanup("testinit_")
	defer env.Cleanup("testinit_ns1")
	defer env.Cleanup("testinit_ns2")

	db, err := env.DBProvider.GetDBHandle("testinit")
	testutil.AssertNoError(t, err, "")
	db.Open()
	defer db.Close()
	batch := statedb.NewUpdateBatch()

	jsonValue1 := "{\"asset_name\": \"marble1\",\"color\": \"blue\",\"size\": 1,\"owner\": \"tom\"}"
	batch.Put("ns1", "key1", []byte(jsonValue1), version.NewHeight(1, 1))
	jsonValue2 := "{\"asset_name\": \"marble2\",\"color\": \"blue\",\"size\": 2,\"owner\": \"jerry\"}"
	batch.Put("ns1", "key2", []byte(jsonValue2), version.NewHeight(1, 2))
	jsonValue3 := "{\"asset_name\": \"marble3\",\"color\": \"blue\",\"size\": 3,\"owner\": \"fred\"}"
	batch.Put("ns1", "key3", []byte(jsonValue3), version.NewHeight(1, 3))
	jsonValue4 := "{\"asset_name\": \"marble4\",\"color\": \"blue\",\"size\": 4,\"owner\": \"martha\"}"
	batch.Put("ns1", "key4", []byte(jsonValue4), version.NewHeight(1, 4))
	jsonValue5 := "{\"asset_name\": \"marble5\",\"color\": \"blue\",\"size\": 5,\"owner\": \"fred\"}"
	batch.Put("ns1", "key5", []byte(jsonValue5), version.NewHeight(1, 5))
	jsonValue6 := "{\"asset_name\": \"marble6\",\"color\": \"blue\",\"size\": 6,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key6", []byte(jsonValue6), version.NewHeight(1, 6))
	jsonValue7 := "{\"asset_name\": \"marble7\",\"color\": \"blue\",\"size\": 7,\"owner\": \"fred\"}"
	batch.Put("ns1", "key7", []byte(jsonValue7), version.NewHeight(1, 7))
	jsonValue8 := "{\"asset_name\": \"marble8\",\"color\": \"blue\",\"size\": 8,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key8", []byte(jsonValue8), version.NewHeight(1, 8))
	jsonValue9 := "{\"asset_name\": \"marble9\",\"color\": \"green\",\"size\": 9,\"owner\": \"fred\"}"
	batch.Put("ns1", "key9", []byte(jsonValue9), version.NewHeight(1, 9))
	jsonValue10 := "{\"asset_name\": \"marble10\",\"color\": \"green\",\"size\": 10,\"owner\": \"mary\"}"
	batch.Put("ns1", "key10", []byte(jsonValue10), version.NewHeight(1, 10))
	jsonValue11 := "{\"asset_name\": \"marble11\",\"color\": \"cyan\",\"size\": 1000007,\"owner\": \"joe\"}"
	batch.Put("ns1", "key11", []byte(jsonValue11), version.NewHeight(1, 11))

	
	batch.Put("ns2", "key1", []byte(jsonValue1), version.NewHeight(1, 12))
	batch.Put("ns2", "key2", []byte(jsonValue2), version.NewHeight(1, 13))
	batch.Put("ns2", "key3", []byte(jsonValue3), version.NewHeight(1, 14))
	batch.Put("ns2", "key4", []byte(jsonValue4), version.NewHeight(1, 15))
	batch.Put("ns2", "key5", []byte(jsonValue5), version.NewHeight(1, 16))
	batch.Put("ns2", "key6", []byte(jsonValue6), version.NewHeight(1, 17))
	batch.Put("ns2", "key7", []byte(jsonValue7), version.NewHeight(1, 18))
	batch.Put("ns2", "key8", []byte(jsonValue8), version.NewHeight(1, 19))
	batch.Put("ns2", "key9", []byte(jsonValue9), version.NewHeight(1, 20))
	batch.Put("ns2", "key10", []byte(jsonValue10), version.NewHeight(1, 21))

	savePoint := version.NewHeight(2, 22)
	db.ApplyUpdates(batch, savePoint)

	
	dbArtifactsTarBytes := testutil.CreateTarBytesForTest(
		[]*testutil.TarFileEntry{
			{Name: "META-INF/statedb/couchdb/indexes/indexColorSortName.json", Body: `{"index":{"fields":[{"color":"desc"}]},"ddoc":"indexColorSortName","name":"indexColorSortName","type":"json"}`},
			{Name: "META-INF/statedb/couchdb/indexes/indexSizeSortName.json", Body: `{"index":{"fields":[{"size":"desc"}]},"ddoc":"indexSizeSortName","name":"indexSizeSortName","type":"json"}`},
			{Name: "META-INF/statedb/couchdb/collections/collectionMarbles/indexes/indexCollMarbles.json", Body: `{"index":{"fields":["docType","owner"]},"ddoc":"indexCollectionMarbles", "name":"indexCollectionMarbles","type":"json"}`},
			{Name: "META-INF/statedb/couchdb/collections/collectionMarblesPrivateDetails/indexes/indexCollPrivDetails.json", Body: `{"index":{"fields":["docType","price"]},"ddoc":"indexPrivateDetails", "name":"indexPrivateDetails","type":"json"}`},
		},
	)

	
	queryString := `{"selector":{"owner":"fred"}}`

	_, err = db.ExecuteQuery("ns1", queryString)
	testutil.AssertNoError(t, err, "")

	
	queryString = `{"selector":{"owner":"fred"}, "sort": [{"size": "desc"}]}`

	_, err = db.ExecuteQuery("ns1", queryString)
	testutil.AssertError(t, err, "Error should have been thrown for a missing index")

	indexCapable, ok := db.(statedb.IndexCapable)

	if !ok {
		t.Fatalf("Couchdb state impl is expected to implement interface `statedb.IndexCapable`")
	}

	fileEntries, errExtract := ccprovider.ExtractFileEntries(dbArtifactsTarBytes, "couchdb")
	testutil.AssertNoError(t, errExtract, "")

	indexCapable.ProcessIndexesForChaincodeDeploy("ns1", fileEntries["META-INF/statedb/couchdb/indexes"])
	
	time.Sleep(100 * time.Millisecond)
	
	queryString = `{"selector":{"owner":"fred"}, "sort": [{"size": "desc"}]}`

	
	_, err = db.ExecuteQuery("ns1", queryString)
	testutil.AssertNoError(t, err, "")

	
	_, err = db.ExecuteQuery("ns2", queryString)
	testutil.AssertError(t, err, "Error should have been thrown for a missing index")

}

func TestTryCastingToJSON(t *testing.T) {
	sampleJSON := []byte(`{"a":"A", "b":"B"}`)
	isJSON, jsonVal := tryCastingToJSON(sampleJSON)
	testutil.AssertEquals(t, isJSON, true)
	testutil.AssertEquals(t, jsonVal["a"], "A")
	testutil.AssertEquals(t, jsonVal["b"], "B")

	sampleNonJSON := []byte(`This is not a json`)
	isJSON, jsonVal = tryCastingToJSON(sampleNonJSON)
	testutil.AssertEquals(t, isJSON, false)
}

func TestHandleChaincodeDeployErroneousIndexFile(t *testing.T) {
	channelName := "ch1"
	env := NewTestVDBEnv(t)
	env.Cleanup(channelName)
	defer env.Cleanup(channelName)
	db, err := env.DBProvider.GetDBHandle(channelName)
	testutil.AssertNoError(t, err, "")
	db.Open()
	defer db.Close()

	batch := statedb.NewUpdateBatch()
	batch.Put("ns1", "key1", []byte(`{"asset_name": "marble1","color": "blue","size": 1,"owner": "tom"}`), version.NewHeight(1, 1))
	batch.Put("ns1", "key2", []byte(`{"asset_name": "marble2","color": "blue","size": 2,"owner": "jerry"}`), version.NewHeight(1, 2))

	
	badSyntaxFileContent := `{"index":{"fields": This is a bad json}`
	dbArtifactsTarBytes := testutil.CreateTarBytesForTest(
		[]*testutil.TarFileEntry{
			{Name: "META-INF/statedb/couchdb/indexes/indexSizeSortName.json", Body: `{"index":{"fields":[{"size":"desc"}]},"ddoc":"indexSizeSortName","name":"indexSizeSortName","type":"json"}`},
			{Name: "META-INF/statedb/couchdb/indexes/badSyntax.json", Body: badSyntaxFileContent},
		},
	)

	indexCapable, ok := db.(statedb.IndexCapable)
	if !ok {
		t.Fatalf("Couchdb state impl is expected to implement interface `statedb.IndexCapable`")
	}

	fileEntries, errExtract := ccprovider.ExtractFileEntries(dbArtifactsTarBytes, "couchdb")
	testutil.AssertNoError(t, errExtract, "")

	indexCapable.ProcessIndexesForChaincodeDeploy("ns1", fileEntries["META-INF/statedb/couchdb/indexes"])

	
	time.Sleep(100 * time.Millisecond)
	
	_, err = db.ExecuteQuery("ns1", `{"selector":{"owner":"fred"}, "sort": [{"size": "desc"}]}`)
	testutil.AssertNoError(t, err, "")
}

func TestIsBulkOptimizable(t *testing.T) {
	var db statedb.VersionedDB = &VersionedDB{}
	_, ok := db.(statedb.BulkOptimizable)
	if !ok {
		t.Fatal("state couch db is expected to implement interface statedb.BulkOptimizable")
	}
}

func printCompositeKeys(keys []*statedb.CompositeKey) string {

	compositeKeyString := []string{}
	for _, key := range keys {
		compositeKeyString = append(compositeKeyString, "["+key.Namespace+","+key.Key+"]")
	}
	return strings.Join(compositeKeyString, ",")
}


func TestPaginatedQuery(t *testing.T) {

	env := NewTestVDBEnv(t)
	env.Cleanup("testpaginatedquery_")
	env.Cleanup("testpaginatedquery_ns1")
	defer env.Cleanup("testpaginatedquery_")
	defer env.Cleanup("testpaginatedquery_ns1")

	db, err := env.DBProvider.GetDBHandle("testpaginatedquery")
	testutil.AssertNoError(t, err, "")
	db.Open()
	defer db.Close()

	batch := statedb.NewUpdateBatch()
	jsonValue1 := "{\"asset_name\": \"marble1\",\"color\": \"blue\",\"size\": 1,\"owner\": \"tom\"}"
	batch.Put("ns1", "key1", []byte(jsonValue1), version.NewHeight(1, 1))
	jsonValue2 := "{\"asset_name\": \"marble2\",\"color\": \"red\",\"size\": 2,\"owner\": \"jerry\"}"
	batch.Put("ns1", "key2", []byte(jsonValue2), version.NewHeight(1, 2))
	jsonValue3 := "{\"asset_name\": \"marble3\",\"color\": \"red\",\"size\": 3,\"owner\": \"fred\"}"
	batch.Put("ns1", "key3", []byte(jsonValue3), version.NewHeight(1, 3))
	jsonValue4 := "{\"asset_name\": \"marble4\",\"color\": \"red\",\"size\": 4,\"owner\": \"martha\"}"
	batch.Put("ns1", "key4", []byte(jsonValue4), version.NewHeight(1, 4))
	jsonValue5 := "{\"asset_name\": \"marble5\",\"color\": \"blue\",\"size\": 5,\"owner\": \"fred\"}"
	batch.Put("ns1", "key5", []byte(jsonValue5), version.NewHeight(1, 5))
	jsonValue6 := "{\"asset_name\": \"marble6\",\"color\": \"red\",\"size\": 6,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key6", []byte(jsonValue6), version.NewHeight(1, 6))
	jsonValue7 := "{\"asset_name\": \"marble7\",\"color\": \"blue\",\"size\": 7,\"owner\": \"fred\"}"
	batch.Put("ns1", "key7", []byte(jsonValue7), version.NewHeight(1, 7))
	jsonValue8 := "{\"asset_name\": \"marble8\",\"color\": \"red\",\"size\": 8,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key8", []byte(jsonValue8), version.NewHeight(1, 8))
	jsonValue9 := "{\"asset_name\": \"marble9\",\"color\": \"green\",\"size\": 9,\"owner\": \"fred\"}"
	batch.Put("ns1", "key9", []byte(jsonValue9), version.NewHeight(1, 9))
	jsonValue10 := "{\"asset_name\": \"marble10\",\"color\": \"green\",\"size\": 10,\"owner\": \"mary\"}"
	batch.Put("ns1", "key10", []byte(jsonValue10), version.NewHeight(1, 10))

	jsonValue11 := "{\"asset_name\": \"marble11\",\"color\": \"cyan\",\"size\": 11,\"owner\": \"joe\"}"
	batch.Put("ns1", "key11", []byte(jsonValue11), version.NewHeight(1, 11))
	jsonValue12 := "{\"asset_name\": \"marble12\",\"color\": \"red\",\"size\": 12,\"owner\": \"martha\"}"
	batch.Put("ns1", "key12", []byte(jsonValue12), version.NewHeight(1, 4))
	jsonValue13 := "{\"asset_name\": \"marble13\",\"color\": \"red\",\"size\": 13,\"owner\": \"james\"}"
	batch.Put("ns1", "key13", []byte(jsonValue13), version.NewHeight(1, 4))
	jsonValue14 := "{\"asset_name\": \"marble14\",\"color\": \"red\",\"size\": 14,\"owner\": \"fred\"}"
	batch.Put("ns1", "key14", []byte(jsonValue14), version.NewHeight(1, 4))
	jsonValue15 := "{\"asset_name\": \"marble15\",\"color\": \"red\",\"size\": 15,\"owner\": \"mary\"}"
	batch.Put("ns1", "key15", []byte(jsonValue15), version.NewHeight(1, 4))
	jsonValue16 := "{\"asset_name\": \"marble16\",\"color\": \"red\",\"size\": 16,\"owner\": \"robert\"}"
	batch.Put("ns1", "key16", []byte(jsonValue16), version.NewHeight(1, 4))
	jsonValue17 := "{\"asset_name\": \"marble17\",\"color\": \"red\",\"size\": 17,\"owner\": \"alan\"}"
	batch.Put("ns1", "key17", []byte(jsonValue17), version.NewHeight(1, 4))
	jsonValue18 := "{\"asset_name\": \"marble18\",\"color\": \"red\",\"size\": 18,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key18", []byte(jsonValue18), version.NewHeight(1, 4))
	jsonValue19 := "{\"asset_name\": \"marble19\",\"color\": \"red\",\"size\": 19,\"owner\": \"alan\"}"
	batch.Put("ns1", "key19", []byte(jsonValue19), version.NewHeight(1, 4))
	jsonValue20 := "{\"asset_name\": \"marble20\",\"color\": \"red\",\"size\": 20,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key20", []byte(jsonValue20), version.NewHeight(1, 4))

	jsonValue21 := "{\"asset_name\": \"marble21\",\"color\": \"cyan\",\"size\": 21,\"owner\": \"joe\"}"
	batch.Put("ns1", "key21", []byte(jsonValue21), version.NewHeight(1, 11))
	jsonValue22 := "{\"asset_name\": \"marble22\",\"color\": \"red\",\"size\": 22,\"owner\": \"martha\"}"
	batch.Put("ns1", "key22", []byte(jsonValue22), version.NewHeight(1, 4))
	jsonValue23 := "{\"asset_name\": \"marble23\",\"color\": \"blue\",\"size\": 23,\"owner\": \"james\"}"
	batch.Put("ns1", "key23", []byte(jsonValue23), version.NewHeight(1, 4))
	jsonValue24 := "{\"asset_name\": \"marble24\",\"color\": \"red\",\"size\": 24,\"owner\": \"fred\"}"
	batch.Put("ns1", "key24", []byte(jsonValue24), version.NewHeight(1, 4))
	jsonValue25 := "{\"asset_name\": \"marble25\",\"color\": \"red\",\"size\": 25,\"owner\": \"mary\"}"
	batch.Put("ns1", "key25", []byte(jsonValue25), version.NewHeight(1, 4))
	jsonValue26 := "{\"asset_name\": \"marble26\",\"color\": \"red\",\"size\": 26,\"owner\": \"robert\"}"
	batch.Put("ns1", "key26", []byte(jsonValue26), version.NewHeight(1, 4))
	jsonValue27 := "{\"asset_name\": \"marble27\",\"color\": \"green\",\"size\": 27,\"owner\": \"alan\"}"
	batch.Put("ns1", "key27", []byte(jsonValue27), version.NewHeight(1, 4))
	jsonValue28 := "{\"asset_name\": \"marble28\",\"color\": \"red\",\"size\": 28,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key28", []byte(jsonValue28), version.NewHeight(1, 4))
	jsonValue29 := "{\"asset_name\": \"marble29\",\"color\": \"red\",\"size\": 29,\"owner\": \"alan\"}"
	batch.Put("ns1", "key29", []byte(jsonValue29), version.NewHeight(1, 4))
	jsonValue30 := "{\"asset_name\": \"marble30\",\"color\": \"red\",\"size\": 30,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key30", []byte(jsonValue30), version.NewHeight(1, 4))

	jsonValue31 := "{\"asset_name\": \"marble31\",\"color\": \"cyan\",\"size\": 31,\"owner\": \"joe\"}"
	batch.Put("ns1", "key31", []byte(jsonValue31), version.NewHeight(1, 11))
	jsonValue32 := "{\"asset_name\": \"marble32\",\"color\": \"red\",\"size\": 32,\"owner\": \"martha\"}"
	batch.Put("ns1", "key32", []byte(jsonValue32), version.NewHeight(1, 4))
	jsonValue33 := "{\"asset_name\": \"marble33\",\"color\": \"red\",\"size\": 33,\"owner\": \"james\"}"
	batch.Put("ns1", "key33", []byte(jsonValue33), version.NewHeight(1, 4))
	jsonValue34 := "{\"asset_name\": \"marble34\",\"color\": \"red\",\"size\": 34,\"owner\": \"fred\"}"
	batch.Put("ns1", "key34", []byte(jsonValue34), version.NewHeight(1, 4))
	jsonValue35 := "{\"asset_name\": \"marble35\",\"color\": \"red\",\"size\": 35,\"owner\": \"mary\"}"
	batch.Put("ns1", "key35", []byte(jsonValue35), version.NewHeight(1, 4))
	jsonValue36 := "{\"asset_name\": \"marble36\",\"color\": \"orange\",\"size\": 36,\"owner\": \"robert\"}"
	batch.Put("ns1", "key36", []byte(jsonValue36), version.NewHeight(1, 4))
	jsonValue37 := "{\"asset_name\": \"marble37\",\"color\": \"red\",\"size\": 37,\"owner\": \"alan\"}"
	batch.Put("ns1", "key37", []byte(jsonValue37), version.NewHeight(1, 4))
	jsonValue38 := "{\"asset_name\": \"marble38\",\"color\": \"yellow\",\"size\": 38,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key38", []byte(jsonValue38), version.NewHeight(1, 4))
	jsonValue39 := "{\"asset_name\": \"marble39\",\"color\": \"red\",\"size\": 39,\"owner\": \"alan\"}"
	batch.Put("ns1", "key39", []byte(jsonValue39), version.NewHeight(1, 4))
	jsonValue40 := "{\"asset_name\": \"marble40\",\"color\": \"red\",\"size\": 40,\"owner\": \"elaine\"}"
	batch.Put("ns1", "key40", []byte(jsonValue40), version.NewHeight(1, 4))

	savePoint := version.NewHeight(2, 22)
	db.ApplyUpdates(batch, savePoint)

	
	dbArtifactsTarBytes := testutil.CreateTarBytesForTest(
		[]*testutil.TarFileEntry{
			{Name: "META-INF/statedb/couchdb/indexes/indexSizeSortName.json", Body: `{"index":{"fields":[{"size":"desc"}]},"ddoc":"indexSizeSortName","name":"indexSizeSortName","type":"json"}`},
		},
	)

	
	queryString := `{"selector":{"color":"red"}}`

	_, err = db.ExecuteQuery("ns1", queryString)
	testutil.AssertNoError(t, err, "")

	
	queryString = `{"selector":{"color":"red"}, "sort": [{"size": "asc"}]}`

	indexCapable, ok := db.(statedb.IndexCapable)

	if !ok {
		t.Fatalf("Couchdb state impl is expected to implement interface `statedb.IndexCapable`")
	}

	fileEntries, errExtract := ccprovider.ExtractFileEntries(dbArtifactsTarBytes, "couchdb")
	testutil.AssertNoError(t, errExtract, "")

	indexCapable.ProcessIndexesForChaincodeDeploy("ns1", fileEntries["META-INF/statedb/couchdb/indexes"])
	
	time.Sleep(100 * time.Millisecond)
	
	queryString = `{"selector":{"color":"red"}, "sort": [{"size": "asc"}]}`

	
	_, err = db.ExecuteQuery("ns1", queryString)
	testutil.AssertNoError(t, err, "")

	
	
	returnKeys := []string{"key2", "key3", "key4", "key6", "key8", "key12", "key13", "key14", "key15", "key16"}
	bookmark, err := executeQuery(t, db, "ns1", queryString, "", int32(10), returnKeys)
	testutil.AssertNoError(t, err, "")
	returnKeys = []string{"key17", "key18", "key19", "key20", "key22", "key24", "key25", "key26", "key28", "key29"}
	bookmark, err = executeQuery(t, db, "ns1", queryString, bookmark, int32(10), returnKeys)
	testutil.AssertNoError(t, err, "")

	returnKeys = []string{"key30", "key32", "key33", "key34", "key35", "key37", "key39", "key40"}
	_, err = executeQuery(t, db, "ns1", queryString, bookmark, int32(10), returnKeys)
	testutil.AssertNoError(t, err, "")

	
	
	returnKeys = []string{"key2", "key3", "key4", "key6", "key8", "key12", "key13", "key14", "key15",
		"key16", "key17", "key18", "key19", "key20", "key22", "key24", "key25", "key26", "key28", "key29",
		"key30", "key32", "key33", "key34", "key35", "key37", "key39", "key40"}
	_, err = executeQuery(t, db, "ns1", queryString, "", int32(50), returnKeys)
	testutil.AssertNoError(t, err, "")

	
	viper.Set("ledger.state.couchDBConfig.internalQueryLimit", 50)

	
	
	returnKeys = []string{"key2", "key3", "key4", "key6", "key8", "key12", "key13", "key14", "key15", "key16"}
	bookmark, err = executeQuery(t, db, "ns1", queryString, "", int32(10), returnKeys)
	testutil.AssertNoError(t, err, "")
	returnKeys = []string{"key17", "key18", "key19", "key20", "key22", "key24", "key25", "key26", "key28", "key29"}
	bookmark, err = executeQuery(t, db, "ns1", queryString, bookmark, int32(10), returnKeys)
	testutil.AssertNoError(t, err, "")
	returnKeys = []string{"key30", "key32", "key33", "key34", "key35", "key37", "key39", "key40"}
	_, err = executeQuery(t, db, "ns1", queryString, bookmark, int32(10), returnKeys)
	testutil.AssertNoError(t, err, "")

	
	viper.Set("ledger.state.couchDBConfig.internalQueryLimit", 10)

	
	returnKeys = []string{"key2", "key3", "key4", "key6", "key8", "key12", "key13", "key14", "key15",
		"key16", "key17", "key18", "key19", "key20", "key22", "key24", "key25", "key26", "key28", "key29",
		"key30", "key32", "key33", "key34", "key35", "key37", "key39", "key40"}
	_, err = executeQuery(t, db, "ns1", queryString, "", int32(0), returnKeys)
	testutil.AssertNoError(t, err, "")

	
	viper.Set("ledger.state.couchDBConfig.internalQueryLimit", 5)

	
	returnKeys = []string{"key2", "key3", "key4", "key6", "key8", "key12", "key13", "key14", "key15", "key16"}
	_, err = executeQuery(t, db, "ns1", queryString, "", int32(10), returnKeys)
	testutil.AssertNoError(t, err, "")

	
	viper.Set("ledger.state.couchDBConfig.internalQueryLimit", 1000)
}

func executeQuery(t *testing.T, db statedb.VersionedDB, namespace, query, bookmark string, limit int32, returnKeys []string) (string, error) {

	var itr statedb.ResultsIterator
	var err error

	if limit == int32(0) && bookmark == "" {
		itr, err = db.ExecuteQuery(namespace, query)
		if err != nil {
			return "", err
		}
	} else {
		queryOptions := make(map[string]interface{})
		if bookmark != "" {
			queryOptions["bookmark"] = bookmark
		}
		if limit != 0 {
			queryOptions["limit"] = limit
		}

		itr, err = db.ExecuteQueryWithMetadata(namespace, query, queryOptions)
		if err != nil {
			return "", err
		}
	}

	
	commontests.TestItrWithoutClose(t, itr, returnKeys)

	returnBookmark := ""
	if queryResultItr, ok := itr.(statedb.QueryResultsIterator); ok {
		returnBookmark = queryResultItr.GetBookmarkAndClose()
	}

	return returnBookmark, nil
}


func TestPaginatedQueryValidation(t *testing.T) {

	queryOptions := make(map[string]interface{})
	queryOptions["bookmark"] = "Test1"
	queryOptions["limit"] = int32(10)

	err := validateQueryMetadata(queryOptions)
	testutil.AssertNoError(t, err, "An error was thrown for a valid options")

	queryOptions = make(map[string]interface{})
	queryOptions["bookmark"] = "Test1"
	queryOptions["limit"] = float64(10.2)

	err = validateQueryMetadata(queryOptions)
	testutil.AssertError(t, err, "An should have been thrown for an invalid options")

	queryOptions = make(map[string]interface{})
	queryOptions["bookmark"] = "Test1"
	queryOptions["limit"] = "10"

	err = validateQueryMetadata(queryOptions)
	testutil.AssertError(t, err, "An should have been thrown for an invalid options")

	queryOptions = make(map[string]interface{})
	queryOptions["bookmark"] = int32(10)
	queryOptions["limit"] = "10"

	err = validateQueryMetadata(queryOptions)
	testutil.AssertError(t, err, "An should have been thrown for an invalid options")

	queryOptions = make(map[string]interface{})
	queryOptions["bookmark"] = "Test1"
	queryOptions["limit1"] = int32(10)

	err = validateQueryMetadata(queryOptions)
	testutil.AssertError(t, err, "An should have been thrown for an invalid options")

	queryOptions = make(map[string]interface{})
	queryOptions["bookmark1"] = "Test1"
	queryOptions["limit1"] = int32(10)

	err = validateQueryMetadata(queryOptions)
	testutil.AssertError(t, err, "An should have been thrown for an invalid options")

}
