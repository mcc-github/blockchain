/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package couchdb

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/util"
)


func TestCreateCouchDBConnectionAndDB(t *testing.T) {

	database := "testcreatecouchdbconnectionanddb"
	cleanup(database)
	defer cleanup(database)
	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to CreateCouchInstance"))

	_, err = CreateCouchDatabase(couchInstance, database)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to CreateCouchDatabase"))

}


func TestCreateCouchDBSystemDBs(t *testing.T) {

	database := "testcreatecouchdbsystemdb"
	cleanup(database)
	defer cleanup(database)

	
	couchInstance, err := CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)

	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to CreateCouchInstance"))

	err = CreateSystemDatabasesIfNotExist(couchInstance)
	testutil.AssertNoError(t, err, fmt.Sprintf("Error when trying to create system databases"))

	db := CouchDatabase{CouchInstance: couchInstance, DBName: "_users"}

	
	dbResp, _, errdb := db.GetDatabaseInfo()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to retrieve _users database information"))
	testutil.AssertEquals(t, dbResp.DbName, "_users")

	db = CouchDatabase{CouchInstance: couchInstance, DBName: "_replicator"}

	
	dbResp, _, errdb = db.GetDatabaseInfo()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to retrieve _replicator database information"))
	testutil.AssertEquals(t, dbResp.DbName, "_replicator")

	db = CouchDatabase{CouchInstance: couchInstance, DBName: "_global_changes"}

	
	dbResp, _, errdb = db.GetDatabaseInfo()
	testutil.AssertNoError(t, errdb, fmt.Sprintf("Error when trying to retrieve _global_changes database information"))
	testutil.AssertEquals(t, dbResp.DbName, "_global_changes")

}
func TestDatabaseMapping(t *testing.T) {
	
	_, err := mapAndValidateDatabaseName("testDB")
	testutil.AssertError(t, err, "Error expected because the name contains capital letters")

	
	_, err = mapAndValidateDatabaseName("test1234/1")
	testutil.AssertError(t, err, "Error expected because the name contains illegal chars")

	
	_, err = mapAndValidateDatabaseName("5test1234")
	testutil.AssertError(t, err, "Error expected because the name starts with a number")

	
	_, err = mapAndValidateDatabaseName("")
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for an invalid name"))

	_, err = mapAndValidateDatabaseName("a12345678901234567890123456789012345678901234" +
		"56789012345678901234567890123456789012345678901234567890123456789012345678901234567890" +
		"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456" +
		"78901234567890123456789012345678901234567890")
	testutil.AssertError(t, err, fmt.Sprintf("Error should have been thrown for an invalid name"))

	transformedName, err := mapAndValidateDatabaseName("test.my.db-1")
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, transformedName, "test$my$db-1")
}

func TestConstructMetadataDBName(t *testing.T) {
	
	chainName := "tob2g.y-z0f.qwp-rq5g4-ogid5g6oucyryg9sc16mz0t4vuake5q557esz7sn493nf0ghch0xih6dwuirokyoi4jvs67gh6r5v6mhz3-292un2-9egdcs88cstg3f7xa9m1i8v4gj0t3jedsm-woh3kgiqehwej6h93hdy5tr4v.1qmmqjzz0ox62k.507sh3fkw3-mfqh.ukfvxlm5szfbwtpfkd1r4j.cy8oft5obvwqpzjxb27xuw6"

	truncatedChainName := "tob2g.y-z0f.qwp-rq5g4-ogid5g6oucyryg9sc16mz0t4vuak"
	testutil.AssertEquals(t, len(truncatedChainName), chainNameAllowedLength)

	
	
	hash := hex.EncodeToString(util.ComputeSHA256([]byte(chainName)))
	expectedDBName := truncatedChainName + "(" + hash + ")" + "_"
	expectedDBNameLength := 117

	constructedDBName := ConstructMetadataDBName(chainName)
	testutil.AssertEquals(t, len(constructedDBName), expectedDBNameLength)
	testutil.AssertEquals(t, constructedDBName, expectedDBName)
}

func TestConstructedNamespaceDBName(t *testing.T) {
	

	
	chainName := "tob2g.y-z0f.qwp-rq5g4-ogid5g6oucyryg9sc16mz0t4vuake5q557esz7sn493nf0ghch0xih6dwuirokyoi4jvs67gh6r5v6mhz3-292un2-9egdcs88cstg3f7xa9m1i8v4gj0t3jedsm-woh3kgiqehwej6h93hdy5tr4v.1qmmqjzz0ox62k.507sh3fkw3-mfqh.ukfvxlm5szfbwtpfkd1r4j.cy8oft5obvwqpzjxb27xuw6"

	
	ns := "wMCnSXiV9YoIqNQyNvFVTdM8XnUtvrOFFIWsKelmP5NEszmNLl8YhtOKbFu3P_NgwgsYF8PsfwjYCD8f1XRpANQLoErDHwLlweryqXeJ6vzT2x0pS_GwSx0m6tBI0zOmHQOq_2De8A87x6zUOPwufC2T6dkidFxiuq8Sey2-5vUo_iNKCij3WTeCnKx78PUIg_U1gp4_0KTvYVtRBRvH0kz5usizBxPaiFu3TPhB9XLviScvdUVSbSYJ0Z"
	
	
	coll := "pvWjtfSTXVK8WJus5s6zWoMIciXd7qHRZIusF9SkOS6m8XuHCiJDE9cCRuVerq22Na8qBL2ywDGFpVMIuzfyEXLjeJb0mMuH4cwewT6r1INOTOSYwrikwOLlT_fl0V1L7IQEwUBB8WCvRqSdj6j5-E5aGul_pv_0UeCdwWiyA_GrZmP7ocLzfj2vP8btigrajqdH-irLO2ydEjQUAvf8fiuxru9la402KmKRy457GgI98UHoUdqV3f3FCdR"

	truncatedChainName := "tob2g.y-z0f.qwp-rq5g4-ogid5g6oucyryg9sc16mz0t4vuak"
	truncatedEscapedNs := "w$m$cn$s$xi$v9$yo$iq$n$qy$nv$f$v$td$m8$xn$utvr$o$f"
	truncatedEscapedColl := "pv$wjtf$s$t$x$v$k8$w$jus5s6z$wo$m$ici$xd7q$h$r$z$i"
	testutil.AssertEquals(t, len(truncatedChainName), chainNameAllowedLength)
	testutil.AssertEquals(t, len(truncatedEscapedNs), namespaceNameAllowedLength)
	testutil.AssertEquals(t, len(truncatedEscapedColl), collectionNameAllowedLength)

	untruncatedDBName := chainName + "_" + ns + "$$" + coll
	hash := hex.EncodeToString(util.ComputeSHA256([]byte(untruncatedDBName)))
	expectedDBName := truncatedChainName + "_" + truncatedEscapedNs + "$$" + truncatedEscapedColl + "(" + hash + ")"
	
	
	
	
	expectedDBNameLength := 219

	namespace := ns + "$$" + coll
	constructedDBName := ConstructNamespaceDBName(chainName, namespace)
	testutil.AssertEquals(t, len(constructedDBName), expectedDBNameLength)
	testutil.AssertEquals(t, constructedDBName, expectedDBName)

	

	untruncatedDBName = chainName + "_" + ns
	hash = hex.EncodeToString(util.ComputeSHA256([]byte(untruncatedDBName)))
	expectedDBName = truncatedChainName + "_" + truncatedEscapedNs + "(" + hash + ")"
	
	
	
	expectedDBNameLength = 167

	namespace = ns
	constructedDBName = ConstructNamespaceDBName(chainName, namespace)
	testutil.AssertEquals(t, len(constructedDBName), expectedDBNameLength)
	testutil.AssertEquals(t, constructedDBName, expectedDBName)
}