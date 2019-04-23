/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privacyenabledstate

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb/statecouchdb"
	"github.com/mcc-github/blockchain/core/ledger/mock"
	"github.com/mcc-github/blockchain/core/ledger/util/couchdb"
	"github.com/mcc-github/blockchain/integration/runner"
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
	dbPath            string
}


func (env *LevelDBCommonStorageTestEnv) Init(t testing.TB) {
	dbPath, err := ioutil.TempDir("", "cstestenv")
	if err != nil {
		t.Fatalf("Failed to create level db storage directory: %s", err)
	}
	env.bookkeeperTestEnv = bookkeeping.NewTestEnv(t)
	dbProvider, err := NewCommonStorageDBProvider(
		env.bookkeeperTestEnv.TestProvider,
		&disabled.Provider{},
		&mock.HealthCheckRegistry{},
		&ledger.StateDB{
			LevelDBPath: dbPath,
		},
	)
	assert.NoError(t, err)
	env.t = t
	env.provider = dbProvider
	env.dbPath = dbPath
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
	os.RemoveAll(env.dbPath)
}




type CouchDBCommonStorageTestEnv struct {
	t                 testing.TB
	provider          DBProvider
	bookkeeperTestEnv *bookkeeping.TestEnv
	redoPath          string
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
	redoPath, err := ioutil.TempDir("", "pestate")
	if err != nil {
		t.Fatalf("Failed to create redo log directory: %s", err)
	}
	couchAddress := env.setupCouch()

	stateDBConfig := &ledger.StateDB{
		StateDatabase: "CouchDB",
		CouchDB: &couchdb.Config{
			Address:             couchAddress,
			Username:            "",
			Password:            "",
			MaxRetries:          3,
			MaxRetriesOnStartup: 20,
			RequestTimeout:      35 * time.Second,
			InternalQueryLimit:  1000,
			MaxBatchUpdateSize:  1000,
			RedoLogPath:         redoPath,
		},
	}

	env.bookkeeperTestEnv = bookkeeping.NewTestEnv(t)
	dbProvider, err := NewCommonStorageDBProvider(
		env.bookkeeperTestEnv.TestProvider,
		&disabled.Provider{},
		&mock.HealthCheckRegistry{},
		stateDBConfig,
	)
	assert.NoError(t, err)
	env.t = t
	env.provider = dbProvider
	env.redoPath = redoPath
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
	os.RemoveAll(env.redoPath)
	env.bookkeeperTestEnv.Cleanup()
	env.provider.Close()
	env.couchCleanup()
}
