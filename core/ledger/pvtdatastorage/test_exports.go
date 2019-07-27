/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/stretchr/testify/assert"
)

func pvtDataConf() *PrivateDataConfig {
	return &PrivateDataConfig{
		PrivateDataConfig: &ledger.PrivateDataConfig{
			BatchesInterval: 1000,
			MaxBatchSize:    5000,
			PurgeInterval:   2,
		},
		StorePath: "",
	}
}


type StoreEnv struct {
	t                 testing.TB
	TestStoreProvider Provider
	TestStore         Store
	ledgerid          string
	btlPolicy         pvtdatapolicy.BTLPolicy
	conf              *PrivateDataConfig
}


func NewTestStoreEnv(
	t *testing.T,
	ledgerid string,
	btlPolicy pvtdatapolicy.BTLPolicy,
	conf *PrivateDataConfig) *StoreEnv {

	storeDir, err := ioutil.TempDir("", "pdstore")
	if err != nil {
		t.Fatalf("Failed to create private data storage directory: %s", err)
	}
	assert := assert.New(t)
	conf.StorePath = storeDir
	testStoreProvider := NewProvider(conf)
	testStore, err := testStoreProvider.OpenStore(ledgerid)
	testStore.Init(btlPolicy)
	assert.NoError(err)
	return &StoreEnv{t, testStoreProvider, testStore, ledgerid, btlPolicy, conf}
}


func (env *StoreEnv) CloseAndReopen() {
	var err error
	env.TestStoreProvider.Close()
	env.TestStoreProvider = NewProvider(env.conf)
	env.TestStore, err = env.TestStoreProvider.OpenStore(env.ledgerid)
	env.TestStore.Init(env.btlPolicy)
	assert.NoError(env.t, err)
}


func (env *StoreEnv) Cleanup() {
	os.RemoveAll(env.conf.StorePath)
}
