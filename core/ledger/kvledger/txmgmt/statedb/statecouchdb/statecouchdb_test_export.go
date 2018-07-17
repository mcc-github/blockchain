/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package statecouchdb

import (
	"strings"
	"testing"

	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/util/couchdb"
)


type TestVDBEnv struct {
	t          testing.TB
	DBProvider statedb.VersionedDBProvider
}


func NewTestVDBEnv(t testing.TB) *TestVDBEnv {
	t.Logf("Creating new TestVDBEnv")

	dbProvider, _ := NewVersionedDBProvider()
	testVDBEnv := &TestVDBEnv{t, dbProvider}
	
	return testVDBEnv
}


func (env *TestVDBEnv) Cleanup(dbName string) {
	env.t.Logf("Cleaningup TestVDBEnv")
	env.DBProvider.Close()
	CleanupDB(dbName)
}


func CleanupDB(dbName string) {
	
	couchDBDef := couchdb.GetCouchDBDefinition()
	couchInstance, _ := couchdb.CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password,
		couchDBDef.MaxRetries, couchDBDef.MaxRetriesOnStartup, couchDBDef.RequestTimeout)
	db := couchdb.CouchDatabase{CouchInstance: couchInstance, DBName: strings.ToLower(dbName)}
	
	db.DropDatabase()
}
