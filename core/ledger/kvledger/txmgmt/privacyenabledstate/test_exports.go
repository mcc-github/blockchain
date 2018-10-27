/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privacyenabledstate

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb/statecouchdb"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/integration/runner"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


type TestEnv interface {
	Init(t testing.TB)
	GetDBHandle(id string) DB
	GetName() string
	Cleanup()
}




var testEnvs = []TestEnv{&LevelDBCommonStorageTestEnv{}, &CouchDBCommonStorageTestEnv{}}




type LevelDBCommonStorageTestEnv struct {
	t                 testing.TB
	provider          DBProvider
	bookkeeperTestEnv *bookkeeping.TestEnv
}


func (env *LevelDBCommonStorageTestEnv) Init(t testing.TB) {
	viper.Set("ledger.state.stateDatabase", "")
	removeDBPath(t)
	env.bookkeeperTestEnv = bookkeeping.NewTestEnv(t)
	dbProvider, err := NewCommonStorageDBProvider(env.bookkeeperTestEnv.TestProvider)
	assert.NoError(t, err)
	env.t = t
	env.provider = dbProvider
}


func (env *LevelDBCommonStorageTestEnv) GetDBHandle(id string) DB {
	db, err := env.provider.GetDBHandle(id)
	assert.NoError(env.t, err)
	return db
}


func (env *LevelDBCommonStorageTestEnv) GetName() string {
	return "levelDBCommonStorageTestEnv"
}


func (env *LevelDBCommonStorageTestEnv) Cleanup() {
	env.provider.Close()
	env.bookkeeperTestEnv.Cleanup()
	removeDBPath(env.t)
}




type CouchDBCommonStorageTestEnv struct {
	t                 testing.TB
	provider          DBProvider
	bookkeeperTestEnv *bookkeeping.TestEnv
	couchCleanup      func()
}

func (env *CouchDBCommonStorageTestEnv) setupCouch() string {
	externalCouch, set := os.LookupEnv("COUCHDB_ADDR")
	if set {
		env.couchCleanup = func() {}
		return externalCouch
	}

	couchDB := &runner.CouchDB{}
	if err := couchDB.Start(); err != nil {
		err := fmt.Errorf("failed to start couchDB: %s", err)
		panic(err)
	}
	env.couchCleanup = func() { couchDB.Stop() }
	return couchDB.Address()
}


func (env *CouchDBCommonStorageTestEnv) Init(t testing.TB) {
	viper.Set("ledger.state.stateDatabase", "CouchDB")
	couchAddr := env.setupCouch()
	viper.Set("ledger.state.couchDBConfig.couchDBAddress", couchAddr)
	
	
	viper.Set("ledger.state.couchDBConfig.username", "")
	viper.Set("ledger.state.couchDBConfig.password", "")
	viper.Set("ledger.state.couchDBConfig.maxRetries", 3)
	viper.Set("ledger.state.couchDBConfig.maxRetriesOnStartup", 20)
	viper.Set("ledger.state.couchDBConfig.requestTimeout", time.Second*35)

	env.bookkeeperTestEnv = bookkeeping.NewTestEnv(t)
	dbProvider, err := NewCommonStorageDBProvider(env.bookkeeperTestEnv.TestProvider)
	assert.NoError(t, err)
	env.t = t
	env.provider = dbProvider
}


func (env *CouchDBCommonStorageTestEnv) GetDBHandle(id string) DB {
	db, err := env.provider.GetDBHandle(id)
	assert.NoError(env.t, err)
	return db
}


func (env *CouchDBCommonStorageTestEnv) GetName() string {
	return "couchDBCommonStorageTestEnv"
}


func (env *CouchDBCommonStorageTestEnv) Cleanup() {
	csdbProvider, _ := env.provider.(*CommonStorageDBProvider)
	statecouchdb.CleanupDB(env.t, csdbProvider.VersionedDBProvider)

	env.bookkeeperTestEnv.Cleanup()
	env.provider.Close()
	env.couchCleanup()
}

func removeDBPath(t testing.TB) {
	dbPath := ledgerconfig.GetStateLevelDBPath()
	if err := os.RemoveAll(dbPath); err != nil {
		t.Fatalf("Err: %s", err)
		t.FailNow()
	}
}
