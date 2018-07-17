/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package couchdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	logging "github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	ledgertestutil "github.com/mcc-github/blockchain/core/ledger/testutil"
	"github.com/mcc-github/blockchain/integration/runner"
	"github.com/spf13/viper"
)

const badConnectURL = "couchdb:5990"
const badParseConnectURL = "http://host.com|5432"
const updateDocumentConflictError = "conflict"
const updateDocumentConflictReason = "Document update conflict."

var couchDBDef *CouchDBDef

func cleanup(database string) error {
	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)

	if err != nil {
		fmt.Println("Unexpected error", err)
		return err
	}
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}
	
	db.DropDatabase()
	return nil
}

type Asset struct {
	ID        string `json:"_id"`
	Rev       string `json:"_rev"`
	AssetName string `json:"asset_name"`
	Color     string `json:"color"`
	Size      string `json:"size"`
	Owner     string `json:"owner"`
}

var assetJSON = []byte(`{"asset_name":"marble1","color":"blue","size":"35","owner":"jerry"}`)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	
	ledgertestutil.SetupCoreYAMLConfig()

	
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

	
	logging.SetModuleLevel("couchdb", "Debug")

	viper.Set("logging.peer", "debug")

	
	couchDBDef = GetCouchDBDefinition()

	
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

func TestDBConnectionDef(t *testing.T) {

	
	_, err := CreateConnectionDefinition(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create database connection definition"))

}

func TestDBBadConnectionDef(t *testing.T) {

	
	_, err := CreateConnectionDefinition(badParseConnectURL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertError(t, err, fmt.Sprintf("Did not receive error when trying to create database connection definition with a bad hostname"))

}

func TestEncodePathElement(t *testing.T) {

	encodedString := encodePathElement("testelement")
	testutil.AssertEquals(t, encodedString, "testelement")

	encodedString = encodePathElement("test element")
	testutil.AssertEquals(t, encodedString, "test%20element")

	encodedString = encodePathElement("/test element")
	testutil.AssertEquals(t, encodedString, "%2Ftest%20element")

	encodedString = encodePathElement("/test element:")
	testutil.AssertEquals(t, encodedString, "%2Ftest%20element:")

	encodedString = encodePathElement("/test+ element:")
	testutil.AssertEquals(t, encodedString, "%2Ftest%2B%20element:")

}

func TestBadCouchDBInstance(t *testing.T) {

	
	badConnectDef := CouchConnectionDef{URL: badParseConnectURL, Username: "", Password: "",
		MaxRetries: 3, MaxRetriesOnStartup: 10, RequestTimeout: time.Second * 30}

	client := &http.Client{}

	
	badCouchDBInstance := CouchInstance{badConnectDef, client}

	
	badDB := CouchDatabase{&badCouchDBInstance, "baddb", 1}

	
	_, err := CreateCouchDatabase(&badCouchDBInstance, "baddbtest")
	testutil.AssertError(t, err, "Error should have been thrown with CreateCouchDatabase and invalid connection")

	
	err = CreateSystemDatabasesIfNotExist(&badCouchDBInstance)
	testutil.AssertError(t, err, "Error should have been thrown with CreateSystemDatabasesIfNotExist and invalid connection")

	
	err = badDB.CreateDatabaseIfNotExist()
	testutil.AssertError(t, err, "Error should have been thrown with CreateDatabaseIfNotExist and invalid connection")

	
	_, _, err = badDB.GetDatabaseInfo()
	testutil.AssertError(t, err, "Error should have been thrown with GetDatabaseInfo and invalid connection")

	
	_, _, err = badCouchDBInstance.VerifyCouchConfig()
	testutil.AssertError(t, err, "Error should have been thrown with VerifyCouchConfig and invalid connection")

	
	_, err = badDB.EnsureFullCommit()
	testutil.AssertError(t, err, "Error should have been thrown with EnsureFullCommit and invalid connection")

	
	_, err = badDB.DropDatabase()
	testutil.AssertError(t, err, "Error should have been thrown with DropDatabase and invalid connection")

	
	_, _, err = badDB.ReadDoc("1")
	testutil.AssertError(t, err, "Error should have been thrown with ReadDoc and invalid connection")

	
	_, err = badDB.SaveDoc("1", "1", nil)
	testutil.AssertError(t, err, "Error should have been thrown with SaveDoc and invalid connection")

	
	err = badDB.DeleteDoc("1", "1")
	testutil.AssertError(t, err, "Error should have been thrown with DeleteDoc and invalid connection")

	
	_, err = badDB.ReadDocRange("1", "2", 1000, 0)
	testutil.AssertError(t, err, "Error should have been thrown with ReadDocRange and invalid connection")

	
	_, err = badDB.QueryDocuments("1")
	testutil.AssertError(t, err, "Error should have been thrown with QueryDocuments and invalid connection")

	
	_, err = badDB.BatchRetrieveDocumentMetadata(nil)
	testutil.AssertError(t, err, "Error should have been thrown with BatchRetrieveDocumentMetadata and invalid connection")

	
	_, err = badDB.BatchUpdateDocuments(nil)
	testutil.AssertError(t, err, "Error should have been thrown with BatchUpdateDocuments and invalid connection")

	
	_, err = badDB.ListIndex()
	testutil.AssertError(t, err, "Error should have been thrown with ListIndex and invalid connection")

	
	_, err = badDB.CreateIndex("")
	testutil.AssertError(t, err, "Error should have been thrown with CreateIndex and invalid connection")

	
	err = badDB.DeleteIndex("", "")
	testutil.AssertError(t, err, "Error should have been thrown with DeleteIndex and invalid connection")

}

func TestDBCreateSaveWithoutRevision(t *testing.T) {

	database := "testdbcreatesavewithoutrevision"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	_, saveerr := db.SaveDoc("2", "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

}

func TestDBCreateEnsureFullCommit(t *testing.T) {

	database := "testdbensurefullcommit"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	_, saveerr := db.SaveDoc("2", "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, commiterr := db.EnsureFullCommit()
	testutil.AssertNoError(t, commiterr, fmt.Sprintf("Error when trying to ensure a full commit"))
}

func TestDBBadDatabaseName(t *testing.T) {

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	_, dberr := CreateCouchDatabase(couchInstance, "testDB")
	testutil.AssertError(t, dberr, "Error should have been thrown for an invalid db name")

	
	couchInstance, err = CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	_, dberr = CreateCouchDatabase(couchInstance, "test132")
	testutil.AssertNoError(t, dberr, fmt.Sprintf("Error when testing a valid database name"))

	
	couchInstance, err = CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	_, dberr = CreateCouchDatabase(couchInstance, "test1234~!@#$%^&*()[]{}.")
	testutil.AssertError(t, dberr, "Error should have been thrown for an invalid db name")

	
	couchInstance, err = CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	_, dberr = CreateCouchDatabase(couchInstance, "a12345678901234567890123456789012345678901234"+
		"56789012345678901234567890123456789012345678901234567890123456789012345678901234567890"+
		"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456"+
		"78901234567890123456789012345678901234567890")
	testutil.AssertError(t, dberr, fmt.Sprintf("Error should have been thrown for invalid database name"))

}

func TestDBBadConnection(t *testing.T) {

	
	
	_, err := CreateCouchInstance(badConnectURL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, 3, couchDBDef.RequestTimeout)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for a bad connection"))
}

func TestBadDBCredentials(t *testing.T) {

	database := "testdbbadcredentials"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	_, err = CreateCouchInstance(couchDBDef.URL, "fred", "fred",
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for bad credentials"))

}

func TestDBCreateDatabaseAndPersist(t *testing.T) {

	
	testDBCreateDatabaseAndPersist(t, couchDBDef.MaxRetries)

	
	testDBCreateDatabaseAndPersist(t, 0)

	
	testBatchBatchOperations(t, couchDBDef.MaxRetries)

	
	testBatchBatchOperations(t, 0)

}

func testDBCreateDatabaseAndPersist(t *testing.T, maxRetries int) {

	database := "testdbcreatedatabaseandpersist"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		maxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	dbResp, _, errdb := db.GetDatabaseInfo()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to retrieve database information"))
	testutil.AssertEquals(t, dbResp.DbName, database)

	
	_, saveerr := db.SaveDoc("idWith/slash", "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	dbGetResp, _, geterr := db.ReadDoc("idWith/slash")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	assetResp := &Asset{}
	geterr = json.Unmarshal(dbGetResp.JSONValue, &assetResp)
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	testutil.AssertEquals(t, assetResp.Owner, "jerry")

	
	_, saveerr = db.SaveDoc("1", "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	dbGetResp, _, geterr = db.ReadDoc("1")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	assetResp = &Asset{}
	geterr = json.Unmarshal(dbGetResp.JSONValue, &assetResp)
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	testutil.AssertEquals(t, assetResp.Owner, "jerry")

	
	assetResp.Owner = "bob"

	
	assetDocUpdated, _ := json.Marshal(assetResp)

	
	_, saveerr = db.SaveDoc("1", "", &CouchDoc{JSONValue: assetDocUpdated, Attachments: nil})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save the updated document"))

	
	dbGetResp, _, geterr = db.ReadDoc("1")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	assetResp = &Asset{}
	json.Unmarshal(dbGetResp.JSONValue, &assetResp)

	
	testutil.AssertEquals(t, assetResp.Owner, "bob")

	testBytes2 := []byte(`test attachment 2`)

	attachment2 := &AttachmentInfo{}
	attachment2.AttachmentBytes = testBytes2
	attachment2.ContentType = "application/octet-stream"
	attachment2.Name = "data"
	attachments2 := []*AttachmentInfo{}
	attachments2 = append(attachments2, attachment2)

	
	_, saveerr = db.SaveDoc("2", "", &CouchDoc{JSONValue: nil, Attachments: attachments2})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	dbGetResp, _, geterr = db.ReadDoc("2")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	testattach := dbGetResp.Attachments[0].AttachmentBytes
	testutil.AssertEquals(t, testattach, testBytes2)

	testBytes3 := []byte{}

	attachment3 := &AttachmentInfo{}
	attachment3.AttachmentBytes = testBytes3
	attachment3.ContentType = "application/octet-stream"
	attachment3.Name = "data"
	attachments3 := []*AttachmentInfo{}
	attachments3 = append(attachments3, attachment3)

	
	_, saveerr = db.SaveDoc("3", "", &CouchDoc{JSONValue: nil, Attachments: attachments3})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	dbGetResp, _, geterr = db.ReadDoc("3")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	testattach = dbGetResp.Attachments[0].AttachmentBytes
	testutil.AssertEquals(t, testattach, testBytes3)

	testBytes4a := []byte(`test attachment 4a`)
	attachment4a := &AttachmentInfo{}
	attachment4a.AttachmentBytes = testBytes4a
	attachment4a.ContentType = "application/octet-stream"
	attachment4a.Name = "data1"

	testBytes4b := []byte(`test attachment 4b`)
	attachment4b := &AttachmentInfo{}
	attachment4b.AttachmentBytes = testBytes4b
	attachment4b.ContentType = "application/octet-stream"
	attachment4b.Name = "data2"

	attachments4 := []*AttachmentInfo{}
	attachments4 = append(attachments4, attachment4a)
	attachments4 = append(attachments4, attachment4b)

	
	_, saveerr = db.SaveDoc("4", "", &CouchDoc{JSONValue: assetJSON, Attachments: attachments4})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save the updated document"))

	
	dbGetResp, _, geterr = db.ReadDoc("4")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	for _, attach4 := range dbGetResp.Attachments {

		currentName := attach4.Name
		if currentName == "data1" {
			testutil.AssertEquals(t, attach4.AttachmentBytes, testBytes4a)
		}
		if currentName == "data2" {
			testutil.AssertEquals(t, attach4.AttachmentBytes, testBytes4b)
		}

	}

	testBytes5a := []byte(`test attachment 5a`)
	attachment5a := &AttachmentInfo{}
	attachment5a.AttachmentBytes = testBytes5a
	attachment5a.ContentType = "application/octet-stream"
	attachment5a.Name = "data1"

	testBytes5b := []byte{}
	attachment5b := &AttachmentInfo{}
	attachment5b.AttachmentBytes = testBytes5b
	attachment5b.ContentType = "application/octet-stream"
	attachment5b.Name = "data2"

	attachments5 := []*AttachmentInfo{}
	attachments5 = append(attachments5, attachment5a)
	attachments5 = append(attachments5, attachment5b)

	
	_, saveerr = db.SaveDoc("5", "", &CouchDoc{JSONValue: assetJSON, Attachments: attachments5})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save the updated document"))

	
	dbGetResp, _, geterr = db.ReadDoc("5")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	for _, attach5 := range dbGetResp.Attachments {

		currentName := attach5.Name
		if currentName == "data1" {
			testutil.AssertEquals(t, attach5.AttachmentBytes, testBytes5a)
		}
		if currentName == "data2" {
			testutil.AssertEquals(t, attach5.AttachmentBytes, testBytes5b)
		}

	}

	
	_, saveerr = db.SaveDoc(string([]byte{0xff, 0xfe, 0xfd}), "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertError(t, saveerr, fmt.Sprintf("Error should have been thrown when saving a document with an invalid ID"))

	
	_, _, readerr := db.ReadDoc(string([]byte{0xff, 0xfe, 0xfd}))
	testutil.AssertError(t, readerr, fmt.Sprintf("Error should have been thrown when reading a document with an invalid ID"))

	
	_, errdbdrop := db.DropDatabase()
	testutil.AssertNoError(t, errdbdrop, fmt.Sprintf("Error dropping database"))

	
	_, _, errdbinfo := db.GetDatabaseInfo()
	testutil.AssertError(t, errdbinfo, fmt.Sprintf("Error should have been thrown for missing database"))

	
	_, saveerr = db.SaveDoc("6", "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertError(t, saveerr, fmt.Sprintf("Error should have been thrown while attempting to save to a deleted database"))

	
	_, _, geterr = db.ReadDoc("6")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error should not have been thrown for a missing database, nil value is returned"))

}

func TestDBRequestTimeout(t *testing.T) {

	database := "testdbrequesttimeout"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	impossibleTimeout := time.Microsecond * 1

	
	
	_, err = CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, 3, impossibleTimeout)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thown while trying to create a couchdb instance with a connection timeout"))

	
	_, err = CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		-1, 3, couchDBDef.RequestTimeout)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown while attempting to create a database"))

}

func TestDBTimeoutConflictRetry(t *testing.T) {

	database := "testdbtimeoutretry"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, 3, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	dbResp, _, errdb := db.GetDatabaseInfo()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to retrieve database information"))
	testutil.AssertEquals(t, dbResp.DbName, database)

	
	_, saveerr := db.SaveDoc("1", "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, _, geterr := db.ReadDoc("1")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	_, saveerr = db.SaveDoc("1", "1-11111111111111111111111111111111", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document with a revision conflict"))

	
	deleteerr := db.DeleteDoc("1", "1-11111111111111111111111111111111")
	testutil.AssertNoError(t, deleteerr, fmt.Sprintf("Error when trying to delete a document with a revision conflict"))

}

func TestDBBadNumberOfRetries(t *testing.T) {

	database := "testdbbadretries"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	_, err = CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		-1, 3, couchDBDef.RequestTimeout)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown while attempting to create a database"))

}

func TestDBBadJSON(t *testing.T) {

	database := "testdbbadjson"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	dbResp, _, errdb := db.GetDatabaseInfo()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to retrieve database information"))
	testutil.AssertEquals(t, dbResp.DbName, database)

	badJSON := []byte(`{"asset_name"}`)

	
	_, saveerr := db.SaveDoc("1", "", &CouchDoc{JSONValue: badJSON, Attachments: nil})
	testutil.AssertError(t, saveerr, fmt.Sprintf("Error should have been thrown for a bad JSON"))

}

func TestPrefixScan(t *testing.T) {

	database := "testprefixscan"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	dbResp, _, errdb := db.GetDatabaseInfo()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to retrieve database information"))
	testutil.AssertEquals(t, dbResp.DbName, database)

	
	for i := 0; i < 20; i++ {
		id1 := string(0) + string(i) + string(0)
		id2 := string(0) + string(i) + string(1)
		id3 := string(0) + string(i) + string(utf8.MaxRune-1)
		_, saveerr := db.SaveDoc(id1, "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
		testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))
		_, saveerr = db.SaveDoc(id2, "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
		testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))
		_, saveerr = db.SaveDoc(id3, "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
		testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	}
	startKey := string(0) + string(10)
	endKey := startKey + string(utf8.MaxRune)
	_, _, geterr := db.ReadDoc(endKey)
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to get lastkey"))

	resultsPtr, geterr := db.ReadDocRange(startKey, endKey, 1000, 0)
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to perform a range scan"))
	testutil.AssertNotNil(t, resultsPtr)
	results := *resultsPtr
	testutil.AssertEquals(t, len(results), 3)
	testutil.AssertEquals(t, results[0].ID, string(0)+string(10)+string(0))
	testutil.AssertEquals(t, results[1].ID, string(0)+string(10)+string(1))
	testutil.AssertEquals(t, results[2].ID, string(0)+string(10)+string(utf8.MaxRune-1))

	
	_, errdbdrop := db.DropDatabase()
	testutil.AssertNoError(t, errdbdrop, fmt.Sprintf("Error dropping database"))

	
	_, _, errdbinfo := db.GetDatabaseInfo()
	testutil.AssertError(t, errdbinfo, fmt.Sprintf("Error should have been thrown for missing database"))

}

func TestDBSaveAttachment(t *testing.T) {

	database := "testdbsaveattachment"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	byteText := []byte(`This is a test document.  This is only a test`)

	attachment := &AttachmentInfo{}
	attachment.AttachmentBytes = byteText
	attachment.ContentType = "text/plain"
	attachment.Name = "valueBytes"

	attachments := []*AttachmentInfo{}
	attachments = append(attachments, attachment)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	_, saveerr := db.SaveDoc("10", "", &CouchDoc{JSONValue: nil, Attachments: attachments})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	couchDoc, _, geterr2 := db.ReadDoc("10")
	testutil.AssertNoError(t, geterr2, fmt.Sprintf("Error when trying to retrieve a document with attachment"))
	testutil.AssertNotNil(t, couchDoc.Attachments)
	testutil.AssertEquals(t, couchDoc.Attachments[0].AttachmentBytes, byteText)

}

func TestDBDeleteDocument(t *testing.T) {

	database := "testdbdeletedocument"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	_, saveerr := db.SaveDoc("2", "", &CouchDoc{JSONValue: assetJSON, Attachments: nil})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, _, readErr := db.ReadDoc("2")
	testutil.AssertNoError(t, readErr, fmt.Sprintf("Error when trying to retrieve a document with attachment"))

	
	deleteErr := db.DeleteDoc("2", "")
	testutil.AssertNoError(t, deleteErr, fmt.Sprintf("Error when trying to delete a document"))

	
	readValue, _, _ := db.ReadDoc("2")
	testutil.AssertNil(t, readValue)

}

func TestDBDeleteNonExistingDocument(t *testing.T) {

	database := "testdbdeletenonexistingdocument"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	deleteErr := db.DeleteDoc("2", "")
	testutil.AssertNoError(t, deleteErr, fmt.Sprintf("Error when trying to delete a non existing document"))
}

func TestCouchDBVersion(t *testing.T) {

	err := checkCouchDBVersion("2.0.0")
	testutil.AssertNoError(t, err, fmt.Sprintf("Error should not have been thrown for valid version"))

	err = checkCouchDBVersion("4.5.0")
	testutil.AssertNoError(t, err, fmt.Sprintf("Error should not have been thrown for valid version"))

	err = checkCouchDBVersion("1.6.5.4")
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for invalid version"))

	err = checkCouchDBVersion("0.0.0.0")
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for invalid version"))

}

func TestIndexOperations(t *testing.T) {

	database := "testindexoperations"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	byteJSON1 := []byte(`{"_id":"1", "asset_name":"marble1","color":"blue","size":1,"owner":"jerry"}`)
	byteJSON2 := []byte(`{"_id":"2", "asset_name":"marble2","color":"red","size":2,"owner":"tom"}`)
	byteJSON3 := []byte(`{"_id":"3", "asset_name":"marble3","color":"green","size":3,"owner":"jerry"}`)
	byteJSON4 := []byte(`{"_id":"4", "asset_name":"marble4","color":"purple","size":4,"owner":"tom"}`)
	byteJSON5 := []byte(`{"_id":"5", "asset_name":"marble5","color":"blue","size":5,"owner":"jerry"}`)
	byteJSON6 := []byte(`{"_id":"6", "asset_name":"marble6","color":"white","size":6,"owner":"tom"}`)
	byteJSON7 := []byte(`{"_id":"7", "asset_name":"marble7","color":"white","size":7,"owner":"tom"}`)
	byteJSON8 := []byte(`{"_id":"8", "asset_name":"marble8","color":"white","size":8,"owner":"tom"}`)
	byteJSON9 := []byte(`{"_id":"9", "asset_name":"marble9","color":"white","size":9,"owner":"tom"}`)
	byteJSON10 := []byte(`{"_id":"10", "asset_name":"marble10","color":"white","size":10,"owner":"tom"}`)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	batchUpdateDocs := []*CouchDoc{}

	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON1, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON2, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON3, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON4, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON5, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON6, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON7, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON8, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON9, Attachments: nil})
	batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: byteJSON10, Attachments: nil})

	_, err = db.BatchUpdateDocuments(batchUpdateDocs)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error adding batch of documents"))

	
	indexDefSize := `{"index":{"fields":[{"size":"desc"}]},"ddoc":"indexSizeSortDoc", "name":"indexSizeSortName","type":"json"}`

	
	_, err = db.CreateIndex(indexDefSize)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while creating an index"))

	
	
	time.Sleep(100 * time.Millisecond)
	listResult, err := db.ListIndex()
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while retrieving indexes"))

	
	testutil.AssertEquals(t, len(listResult), 1)

	
	for _, elem := range listResult {
		testutil.AssertEquals(t, elem.DesignDocument, "indexSizeSortDoc")
		testutil.AssertEquals(t, elem.Name, "indexSizeSortName")
		
		testutil.AssertEquals(t, strings.Contains(elem.Definition, `"fields":[{"size":"desc"}]`), true)
	}

	
	indexDefColor := `{"index":{"fields":[{"color":"desc"}]}}`

	
	_, err = db.CreateIndex(indexDefColor)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while creating an index"))

	
	
	time.Sleep(100 * time.Millisecond)
	listResult, err = db.ListIndex()
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while retrieving indexes"))

	
	testutil.AssertEquals(t, len(listResult), 2)

	
	err = db.DeleteIndex("indexSizeSortDoc", "indexSizeSortName")
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while deleting an index"))

	
	
	time.Sleep(100 * time.Millisecond)
	listResult, err = db.ListIndex()
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while retrieving indexes"))

	
	testutil.AssertEquals(t, len(listResult), 1)

	
	for _, elem := range listResult {
		err = db.DeleteIndex(elem.DesignDocument, elem.Name)
		testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while deleting an index"))
	}

	
	
	time.Sleep(100 * time.Millisecond)
	listResult, err = db.ListIndex()
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while retrieving indexes"))
	testutil.AssertEquals(t, len(listResult), 0)

	
	queryString := `{"selector":{"size": {"$gt": 0}},"fields": ["_id", "_rev", "owner", "asset_name", "color", "size"], "sort":[{"size":"desc"}], "limit": 10,"skip": 0}`

	
	_, err = db.QueryDocuments(queryString)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have thrown while querying without a valid index"))

	
	_, err = db.CreateIndex(indexDefSize)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while creating an index"))

	
	time.Sleep(100 * time.Millisecond)

	
	_, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while querying with an index"))

	
	indexDefSize = `{"index":{"fields":[{"data.size":"desc"},{"data.owner":"desc"}]},"ddoc":"indexSizeOwnerSortDoc", "name":"indexSizeOwnerSortName","type":"json"}`

	
	dbResp, err := db.CreateIndex(indexDefSize)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while creating an index"))

	
	testutil.AssertEquals(t, dbResp.Result, "created")

	
	time.Sleep(100 * time.Millisecond)

	
	dbResp, err = db.CreateIndex(indexDefSize)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while creating an index"))

	
	testutil.AssertEquals(t, dbResp.Result, "exists")

	
	
	time.Sleep(100 * time.Millisecond)
	listResult, err = db.ListIndex()
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while retrieving indexes"))

	
	testutil.AssertEquals(t, len(listResult), 2)

	
	indexDefSize = `{"index"{"fields":[{"data.size":"desc"},{"data.owner":"desc"}]},"ddoc":"indexSizeOwnerSortDoc", "name":"indexSizeOwnerSortName","type":"json"}`

	
	_, err = db.CreateIndex(indexDefSize)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for an invalid index JSON"))

	
	indexDefSize = `{"index":{"fields2":[{"data.size":"desc"},{"data.owner":"desc"}]},"ddoc":"indexSizeOwnerSortDoc", "name":"indexSizeOwnerSortName","type":"json"}`

	
	_, err = db.CreateIndex(indexDefSize)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for an invalid index definition"))

}

func TestRichQuery(t *testing.T) {

	byteJSON01 := []byte(`{"asset_name":"marble01","color":"blue","size":1,"owner":"jerry"}`)
	byteJSON02 := []byte(`{"asset_name":"marble02","color":"red","size":2,"owner":"tom"}`)
	byteJSON03 := []byte(`{"asset_name":"marble03","color":"green","size":3,"owner":"jerry"}`)
	byteJSON04 := []byte(`{"asset_name":"marble04","color":"purple","size":4,"owner":"tom"}`)
	byteJSON05 := []byte(`{"asset_name":"marble05","color":"blue","size":5,"owner":"jerry"}`)
	byteJSON06 := []byte(`{"asset_name":"marble06","color":"white","size":6,"owner":"tom"}`)
	byteJSON07 := []byte(`{"asset_name":"marble07","color":"white","size":7,"owner":"tom"}`)
	byteJSON08 := []byte(`{"asset_name":"marble08","color":"white","size":8,"owner":"tom"}`)
	byteJSON09 := []byte(`{"asset_name":"marble09","color":"white","size":9,"owner":"tom"}`)
	byteJSON10 := []byte(`{"asset_name":"marble10","color":"white","size":10,"owner":"tom"}`)
	byteJSON11 := []byte(`{"asset_name":"marble11","color":"green","size":11,"owner":"tom"}`)
	byteJSON12 := []byte(`{"asset_name":"marble12","color":"green","size":12,"owner":"frank"}`)

	attachment1 := &AttachmentInfo{}
	attachment1.AttachmentBytes = []byte(`marble01 - test attachment`)
	attachment1.ContentType = "application/octet-stream"
	attachment1.Name = "data"
	attachments1 := []*AttachmentInfo{}
	attachments1 = append(attachments1, attachment1)

	attachment2 := &AttachmentInfo{}
	attachment2.AttachmentBytes = []byte(`marble02 - test attachment`)
	attachment2.ContentType = "application/octet-stream"
	attachment2.Name = "data"
	attachments2 := []*AttachmentInfo{}
	attachments2 = append(attachments2, attachment2)

	attachment3 := &AttachmentInfo{}
	attachment3.AttachmentBytes = []byte(`marble03 - test attachment`)
	attachment3.ContentType = "application/octet-stream"
	attachment3.Name = "data"
	attachments3 := []*AttachmentInfo{}
	attachments3 = append(attachments3, attachment3)

	attachment4 := &AttachmentInfo{}
	attachment4.AttachmentBytes = []byte(`marble04 - test attachment`)
	attachment4.ContentType = "application/octet-stream"
	attachment4.Name = "data"
	attachments4 := []*AttachmentInfo{}
	attachments4 = append(attachments4, attachment4)

	attachment5 := &AttachmentInfo{}
	attachment5.AttachmentBytes = []byte(`marble05 - test attachment`)
	attachment5.ContentType = "application/octet-stream"
	attachment5.Name = "data"
	attachments5 := []*AttachmentInfo{}
	attachments5 = append(attachments5, attachment5)

	attachment6 := &AttachmentInfo{}
	attachment6.AttachmentBytes = []byte(`marble06 - test attachment`)
	attachment6.ContentType = "application/octet-stream"
	attachment6.Name = "data"
	attachments6 := []*AttachmentInfo{}
	attachments6 = append(attachments6, attachment6)

	attachment7 := &AttachmentInfo{}
	attachment7.AttachmentBytes = []byte(`marble07 - test attachment`)
	attachment7.ContentType = "application/octet-stream"
	attachment7.Name = "data"
	attachments7 := []*AttachmentInfo{}
	attachments7 = append(attachments7, attachment7)

	attachment8 := &AttachmentInfo{}
	attachment8.AttachmentBytes = []byte(`marble08 - test attachment`)
	attachment8.ContentType = "application/octet-stream"
	attachment7.Name = "data"
	attachments8 := []*AttachmentInfo{}
	attachments8 = append(attachments8, attachment8)

	attachment9 := &AttachmentInfo{}
	attachment9.AttachmentBytes = []byte(`marble09 - test attachment`)
	attachment9.ContentType = "application/octet-stream"
	attachment9.Name = "data"
	attachments9 := []*AttachmentInfo{}
	attachments9 = append(attachments9, attachment9)

	attachment10 := &AttachmentInfo{}
	attachment10.AttachmentBytes = []byte(`marble10 - test attachment`)
	attachment10.ContentType = "application/octet-stream"
	attachment10.Name = "data"
	attachments10 := []*AttachmentInfo{}
	attachments10 = append(attachments10, attachment10)

	attachment11 := &AttachmentInfo{}
	attachment11.AttachmentBytes = []byte(`marble11 - test attachment`)
	attachment11.ContentType = "application/octet-stream"
	attachment11.Name = "data"
	attachments11 := []*AttachmentInfo{}
	attachments11 = append(attachments11, attachment11)

	attachment12 := &AttachmentInfo{}
	attachment12.AttachmentBytes = []byte(`marble12 - test attachment`)
	attachment12.ContentType = "application/octet-stream"
	attachment12.Name = "data"
	attachments12 := []*AttachmentInfo{}
	attachments12 = append(attachments12, attachment12)

	database := "testrichquery"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	_, saveerr := db.SaveDoc("marble01", "", &CouchDoc{JSONValue: byteJSON01, Attachments: attachments1})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble02", "", &CouchDoc{JSONValue: byteJSON02, Attachments: attachments2})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble03", "", &CouchDoc{JSONValue: byteJSON03, Attachments: attachments3})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble04", "", &CouchDoc{JSONValue: byteJSON04, Attachments: attachments4})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble05", "", &CouchDoc{JSONValue: byteJSON05, Attachments: attachments5})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble06", "", &CouchDoc{JSONValue: byteJSON06, Attachments: attachments6})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble07", "", &CouchDoc{JSONValue: byteJSON07, Attachments: attachments7})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble08", "", &CouchDoc{JSONValue: byteJSON08, Attachments: attachments8})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble09", "", &CouchDoc{JSONValue: byteJSON09, Attachments: attachments9})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble10", "", &CouchDoc{JSONValue: byteJSON10, Attachments: attachments10})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble11", "", &CouchDoc{JSONValue: byteJSON11, Attachments: attachments11})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	_, saveerr = db.SaveDoc("marble12", "", &CouchDoc{JSONValue: byteJSON12, Attachments: attachments12})
	testutil.AssertNoError(t, saveerr, fmt.Sprintf("Error when trying to save a document"))

	
	queryString := `{"selector":{"owner":}}`

	_, err = db.QueryDocuments(queryString)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for bad json"))

	
	queryString = `{"selector":{"owner":{"$eq":"jerry"}}}`

	queryResult, err := db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 3)

	
	queryString = `{"selector":{"owner":"jerry"}}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 3)

	
	queryString = `{"selector":{"owner":{"$eq":"jerry"}},"fields": ["owner","asset_name","color","size"]}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 3)

	
	queryString = `{"selector":{"$or":[{"owner":{"$eq":"jerry"}},{"owner": {"$eq": "frank"}}]}}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 4)

	
	queryString = `{"selector":{"color":"green","$or":[{"owner":"tom"},{"owner":"frank"}]}}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 2)

	
	queryString = `{"selector":{"$and":[{"size":{"$gte":2}},{"size":{"$lte":5}}]}}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 4)

	
	queryString = `{"selector":{"$and":[{"size":{"$gte":3}},{"size":{"$lte":10}},{"$not":{"size":7}}]}}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 7)

	
	queryString = `{"selector":{"$and":[{"size":{"$gte":2}},{"size":{"$lte":10}},{"$nor":[{"size":3},{"size":5},{"size":7}]}]}}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 6)

	
	queryResult, err = db.ReadDocRange("marble02", "marble06", 10000, 0)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a range query"))

	
	testutil.AssertEquals(t, len(*queryResult), 4)

	
	queryString = `{"selector":{"owner":{"$eq":"tom"}}}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 8)

	
	queryString = `{"selector":{"owner":{"$eq":"tom"}},"limit":2}`

	queryResult, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query"))

	
	testutil.AssertEquals(t, len(*queryResult), 2)

	
	queryString = `{"selector":{"owner":"tom"}, "use_index":["indexOwnerDoc","indexOwner"]}`

	_, err = db.QueryDocuments(queryString)
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for an invalid index"))

	
	indexDefSize := `{"index":{"fields":[{"size":"desc"}]},"ddoc":"indexSizeSortDoc", "name":"indexSizeSortName","type":"json"}`

	
	_, err = db.CreateIndex(indexDefSize)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error thrown while creating an index"))

	
	time.Sleep(100 * time.Millisecond)

	
	queryString = `{"selector":{"size":{"$gt":0}}, "use_index":["indexSizeSortDoc","indexSizeSortName"]}`

	_, err = db.QueryDocuments(queryString)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to execute a query with a valid index"))

	
	queryString = `{"selector":{"owner":{"$eq":"tom"}}, "use_index":"indexSizeSortName"}`

	
	_, err = db.QueryDocuments(queryString)
	testutil.AssertError(t, err, fmt.Sprintf("400 error should have been thrown for a missing index"))
	testutil.AssertEquals(t, strings.Contains(err.Error(), "Status Code:400"), true)

	
	queryString = `{"selector":{"owner":{"$eq":"tom"}}, "use_index":["indexSizeSortDoc","indexSizeSortName"]}`

	
	_, err = db.QueryDocuments(queryString)
	testutil.AssertError(t, err, fmt.Sprintf("500 error should have been thrown for a missing index with design doc specified"))
	testutil.AssertEquals(t, strings.Contains(err.Error(), "Status Code:500"), true)

}

func testBatchBatchOperations(t *testing.T, maxRetries int) {

	byteJSON01 := []byte(`{"_id":"marble01","asset_name":"marble01","color":"blue","size":"1","owner":"jerry"}`)
	byteJSON02 := []byte(`{"_id":"marble02","asset_name":"marble02","color":"red","size":"2","owner":"tom"}`)
	byteJSON03 := []byte(`{"_id":"marble03","asset_name":"marble03","color":"green","size":"3","owner":"jerry"}`)
	byteJSON04 := []byte(`{"_id":"marble04","asset_name":"marble04","color":"purple","size":"4","owner":"tom"}`)
	byteJSON05 := []byte(`{"_id":"marble05","asset_name":"marble05","color":"blue","size":"5","owner":"jerry"}`)
	byteJSON06 := []byte(`{"_id":"marble06#$&'()*+,/:;=?@[]","asset_name":"marble06#$&'()*+,/:;=?@[]","color":"blue","size":"6","owner":"jerry"}`)

	attachment1 := &AttachmentInfo{}
	attachment1.AttachmentBytes = []byte(`marble01 - test attachment`)
	attachment1.ContentType = "application/octet-stream"
	attachment1.Name = "data"
	attachments1 := []*AttachmentInfo{}
	attachments1 = append(attachments1, attachment1)

	attachment2 := &AttachmentInfo{}
	attachment2.AttachmentBytes = []byte(`marble02 - test attachment`)
	attachment2.ContentType = "application/octet-stream"
	attachment2.Name = "data"
	attachments2 := []*AttachmentInfo{}
	attachments2 = append(attachments2, attachment2)

	attachment3 := &AttachmentInfo{}
	attachment3.AttachmentBytes = []byte(`marble03 - test attachment`)
	attachment3.ContentType = "application/octet-stream"
	attachment3.Name = "data"
	attachments3 := []*AttachmentInfo{}
	attachments3 = append(attachments3, attachment3)

	attachment4 := &AttachmentInfo{}
	attachment4.AttachmentBytes = []byte(`marble04 - test attachment`)
	attachment4.ContentType = "application/octet-stream"
	attachment4.Name = "data"
	attachments4 := []*AttachmentInfo{}
	attachments4 = append(attachments4, attachment4)

	attachment5 := &AttachmentInfo{}
	attachment5.AttachmentBytes = []byte(`marble05 - test attachment`)
	attachment5.ContentType = "application/octet-stream"
	attachment5.Name = "data"
	attachments5 := []*AttachmentInfo{}
	attachments5 = append(attachments5, attachment5)

	attachment6 := &AttachmentInfo{}
	attachment6.AttachmentBytes = []byte(`marble06#$&'()*+,/:;=?@[] - test attachment`)
	attachment6.ContentType = "application/octet-stream"
	attachment6.Name = "data"
	attachments6 := []*AttachmentInfo{}
	attachments6 = append(attachments6, attachment6)

	database := "testbatch"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	batchUpdateDocs := []*CouchDoc{}

	value1 := &CouchDoc{JSONValue: byteJSON01, Attachments: attachments1}
	value2 := &CouchDoc{JSONValue: byteJSON02, Attachments: attachments2}
	value3 := &CouchDoc{JSONValue: byteJSON03, Attachments: attachments3}
	value4 := &CouchDoc{JSONValue: byteJSON04, Attachments: attachments4}
	value5 := &CouchDoc{JSONValue: byteJSON05, Attachments: attachments5}
	value6 := &CouchDoc{JSONValue: byteJSON06, Attachments: attachments6}

	batchUpdateDocs = append(batchUpdateDocs, value1)
	batchUpdateDocs = append(batchUpdateDocs, value2)
	batchUpdateDocs = append(batchUpdateDocs, value3)
	batchUpdateDocs = append(batchUpdateDocs, value4)
	batchUpdateDocs = append(batchUpdateDocs, value5)
	batchUpdateDocs = append(batchUpdateDocs, value6)

	batchUpdateResp, err := db.BatchUpdateDocuments(batchUpdateDocs)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to update a batch of documents"))

	
	for _, updateDoc := range batchUpdateResp {
		testutil.AssertEquals(t, updateDoc.Ok, true)
	}

	
	
	dbGetResp, _, geterr := db.ReadDoc("marble01")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when attempting read a document"))

	assetResp := &Asset{}
	geterr = json.Unmarshal(dbGetResp.JSONValue, &assetResp)
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))
	
	testutil.AssertEquals(t, assetResp.Owner, "jerry")

	
	
	
	dbGetResp, _, geterr = db.ReadDoc("marble06#$&'()*+,/:;=?@[]")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when attempting read a document"))

	assetResp = &Asset{}
	geterr = json.Unmarshal(dbGetResp.JSONValue, &assetResp)
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))
	
	testutil.AssertEquals(t, assetResp.Owner, "jerry")

	
	
	dbGetResp, _, geterr = db.ReadDoc("marble03")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when attempting read a document"))
	
	attachments := dbGetResp.Attachments
	
	retrievedAttachment := attachments[0]
	
	testutil.AssertEquals(t, attachment3.AttachmentBytes, retrievedAttachment.AttachmentBytes)
	
	
	batchUpdateDocs = []*CouchDoc{}
	batchUpdateDocs = append(batchUpdateDocs, value1)
	batchUpdateDocs = append(batchUpdateDocs, value2)
	batchUpdateResp, err = db.BatchUpdateDocuments(batchUpdateDocs)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to update a batch of documents"))
	
	
	for _, updateDoc := range batchUpdateResp {
		testutil.AssertEquals(t, updateDoc.Ok, false)
		testutil.AssertEquals(t, updateDoc.Error, updateDocumentConflictError)
		testutil.AssertEquals(t, updateDoc.Reason, updateDocumentConflictReason)
	}

	
	

	var keys []string

	keys = append(keys, "marble01")
	keys = append(keys, "marble03")

	batchRevs, err := db.BatchRetrieveDocumentMetadata(keys)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting retrieve revisions"))

	batchUpdateDocs = []*CouchDoc{}

	
	for _, revdoc := range batchRevs {
		if revdoc.ID == "marble01" {
			
			marble01Doc := addRevisionAndDeleteStatus(revdoc.Rev, byteJSON01, false)
			batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: marble01Doc, Attachments: attachments1})
		}

		if revdoc.ID == "marble03" {
			
			marble03Doc := addRevisionAndDeleteStatus(revdoc.Rev, byteJSON03, false)
			batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: marble03Doc, Attachments: attachments3})
		}
	}

	
	batchUpdateResp, err = db.BatchUpdateDocuments(batchUpdateDocs)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to update a batch of documents"))
	
	for _, updateDoc := range batchUpdateResp {
		testutil.AssertEquals(t, updateDoc.Ok, true)
	}

	
	

	keys = []string{}

	keys = append(keys, "marble02")
	keys = append(keys, "marble04")

	batchRevs, err = db.BatchRetrieveDocumentMetadata(keys)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting retrieve revisions"))

	batchUpdateDocs = []*CouchDoc{}

	
	for _, revdoc := range batchRevs {
		if revdoc.ID == "marble02" {
			
			marble02Doc := addRevisionAndDeleteStatus(revdoc.Rev, byteJSON02, true)
			batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: marble02Doc, Attachments: attachments1})
		}
		if revdoc.ID == "marble04" {
			
			marble04Doc := addRevisionAndDeleteStatus(revdoc.Rev, byteJSON04, true)
			batchUpdateDocs = append(batchUpdateDocs, &CouchDoc{JSONValue: marble04Doc, Attachments: attachments3})
		}
	}

	
	batchUpdateResp, err = db.BatchUpdateDocuments(batchUpdateDocs)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when attempting to update a batch of documents"))

	
	for _, updateDoc := range batchUpdateResp {
		testutil.AssertEquals(t, updateDoc.Ok, true)
	}

	
	dbGetResp, _, geterr = db.ReadDoc("marble02")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	testutil.AssertNil(t, dbGetResp)

	
	dbGetResp, _, geterr = db.ReadDoc("marble04")
	testutil.AssertNoError(t, geterr, fmt.Sprintf("Error when trying to retrieve a document"))

	
	testutil.AssertNil(t, dbGetResp)

}


func addRevisionAndDeleteStatus(revision string, value []byte, deleted bool) []byte {

	
	jsonMap := make(map[string]interface{})

	json.Unmarshal(value, &jsonMap)

	
	if revision != "" {
		jsonMap["_rev"] = revision
	}

	
	if deleted {
		jsonMap["_deleted"] = true
	}
	
	returnJSON, _ := json.Marshal(jsonMap)

	return returnJSON

}

func TestDatabaseSecuritySettings(t *testing.T) {

	database := "testdbsecuritysettings"
	err := cleanup(database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to cleanup  Error: %s", err))
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create couch instance"))
	db := CouchDatabase{CouchInstance: couchInstance, DBName: database}

	
	errdb := db.CreateDatabaseIfNotExist()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to create database"))

	
	securityPermissions := &DatabaseSecurity{}
	securityPermissions.Admins.Names = append(securityPermissions.Admins.Names, "admin")
	securityPermissions.Members.Names = append(securityPermissions.Members.Names, "admin")

	
	err = db.ApplyDatabaseSecurity(securityPermissions)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to apply database security"))

	
	databaseSecurity, err := db.GetDatabaseSecurity()
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when retrieving database security"))

	
	testutil.AssertEquals(t, databaseSecurity.Admins.Names[0], "admin")

	
	testutil.AssertEquals(t, databaseSecurity.Members.Names[0], "admin")

	
	securityPermissions = &DatabaseSecurity{}

	
	err = db.ApplyDatabaseSecurity(securityPermissions)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to apply database security"))

	
	databaseSecurity, err = db.GetDatabaseSecurity()
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when retrieving database security"))

	
	testutil.AssertEquals(t, len(databaseSecurity.Admins.Names), 0)

	
	testutil.AssertEquals(t, len(databaseSecurity.Members.Names), 0)

}
