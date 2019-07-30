/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lockbasedtxmgr

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/mock"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	btltestutil "github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/testutil"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flogging.ActivateSpec(
		"lockbasedtxmgr,statevalidator,statebasedval,statecouchdb,valimpl,pvtstatepurgemgmt,valinternal=debug",
	)
	exitCode := m.Run()
	for _, testEnv := range testEnvs {
		testEnv.cleanup()
	}
	os.Exit(exitCode)
}

type testEnv interface {
	cleanup()
	getName() string
	getTxMgr() txmgr.TxMgr
	getVDB() privacyenabledstate.DB
	init(t *testing.T, testLedgerID string, btlPolicy pvtdatapolicy.BTLPolicy)
}

const (
	levelDBtestEnvName = "levelDB_LockBasedTxMgr"
	couchDBtestEnvName = "couchDB_LockBasedTxMgr"
)



var testEnvs = []testEnv{
	&lockBasedEnv{name: levelDBtestEnvName, testDBEnv: &privacyenabledstate.LevelDBCommonStorageTestEnv{}},
	&lockBasedEnv{name: couchDBtestEnvName, testDBEnv: &privacyenabledstate.CouchDBCommonStorageTestEnv{}},
}

var testEnvsMap = map[string]testEnv{
	levelDBtestEnvName: testEnvs[0],
	couchDBtestEnvName: testEnvs[1],
}



type lockBasedEnv struct {
	dbInitialized      bool
	name               string
	t                  testing.TB
	testBookkeepingEnv *bookkeeping.TestEnv
	testDB             privacyenabledstate.DB
	testDBEnv          privacyenabledstate.TestEnv
	txmgr              txmgr.TxMgr
}

func (env *lockBasedEnv) getName() string {
	return env.name
}

func (env *lockBasedEnv) init(t *testing.T, testLedgerID string, btlPolicy pvtdatapolicy.BTLPolicy) {
	var err error
	env.t = t
	if env.dbInitialized == false {
		env.testDBEnv.Init(t)
		env.dbInitialized = true
	}
	env.testDB = env.testDBEnv.GetDBHandle(testLedgerID)
	assert.NoError(t, err)
	if btlPolicy == nil {
		btlPolicy = btltestutil.SampleBTLPolicy(
			map[[2]string]uint64{},
		)
	}
	env.testBookkeepingEnv = bookkeeping.NewTestEnv(t)
	env.txmgr, err = NewLockBasedTxMgr(
		testLedgerID, env.testDB, nil,
		btlPolicy, env.testBookkeepingEnv.TestProvider,
		&mock.DeployedChaincodeInfoProvider{},
		nil,
	)
	assert.NoError(t, err)

}

func (env *lockBasedEnv) getTxMgr() txmgr.TxMgr {
	return env.txmgr
}

func (env *lockBasedEnv) getVDB() privacyenabledstate.DB {
	return env.testDB
}

func (env *lockBasedEnv) cleanup() {
	if env.dbInitialized {
		env.txmgr.Shutdown()
		env.testDBEnv.Cleanup()
		env.testBookkeepingEnv.Cleanup()
	}
}



type txMgrTestHelper struct {
	t     *testing.T
	txMgr txmgr.TxMgr
	bg    *testutil.BlockGenerator
}

func newTxMgrTestHelper(t *testing.T, txMgr txmgr.TxMgr) *txMgrTestHelper {
	bg, _ := testutil.NewBlockGenerator(t, "testLedger", false)
	return &txMgrTestHelper{t, txMgr, bg}
}

func (h *txMgrTestHelper) validateAndCommitRWSet(txRWSet *rwset.TxReadWriteSet) {
	rwSetBytes, _ := proto.Marshal(txRWSet)
	block := h.bg.NextBlock([][]byte{rwSetBytes})
	_, _, err := h.txMgr.ValidateAndPrepare(&ledger.BlockAndPvtData{Block: block, PvtData: nil}, true)
	assert.NoError(h.t, err)
	txsFltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	invalidTxNum := 0
	for i := 0; i < len(block.Data.Data); i++ {
		if txsFltr.IsInvalid(i) {
			invalidTxNum++
		}
	}
	assert.Equal(h.t, 0, invalidTxNum)
	err = h.txMgr.Commit()
	assert.NoError(h.t, err)
}

func (h *txMgrTestHelper) checkRWsetInvalid(txRWSet *rwset.TxReadWriteSet) {
	rwSetBytes, _ := proto.Marshal(txRWSet)
	block := h.bg.NextBlock([][]byte{rwSetBytes})
	_, _, err := h.txMgr.ValidateAndPrepare(&ledger.BlockAndPvtData{Block: block, PvtData: nil}, true)
	assert.NoError(h.t, err)
	txsFltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	invalidTxNum := 0
	for i := 0; i < len(block.Data.Data); i++ {
		if txsFltr.IsInvalid(i) {
			invalidTxNum++
		}
	}
	assert.Equal(h.t, 1, invalidTxNum)
}

func populateCollConfigForTest(t *testing.T, txMgr *LockBasedTxMgr, nsColls []collConfigkey, ht *version.Height) {
	m := map[string]*common.CollectionConfigPackage{}
	for _, nsColl := range nsColls {
		ns, coll := nsColl.ns, nsColl.coll
		pkg, ok := m[ns]
		if !ok {
			pkg = &common.CollectionConfigPackage{}
			m[ns] = pkg
		}
		sCollConfig := &common.CollectionConfig_StaticCollectionConfig{
			StaticCollectionConfig: &common.StaticCollectionConfig{
				Name: coll,
			},
		}
		pkg.Config = append(pkg.Config, &common.CollectionConfig{Payload: sCollConfig})
	}
	ccInfoProvider := &mock.DeployedChaincodeInfoProvider{}
	ccInfoProvider.ChaincodeInfoStub = func(channelName, ccName string, qe ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error) {
		fmt.Printf("retrieveing info for [%s] from [%s]\n", ccName, m)
		return &ledger.DeployedChaincodeInfo{Name: ccName, ExplicitCollectionConfigPkg: m[ccName]}, nil
	}
	txMgr.ccInfoProvider = ccInfoProvider
}

func testutilPopulateDB(
	t *testing.T, txMgr *LockBasedTxMgr, ns string,
	data []*queryresult.KV, pvtdataHashes []*testutilPvtdata,
	version *version.Height,
) {
	updates := privacyenabledstate.NewUpdateBatch()
	for _, kv := range data {
		updates.PubUpdates.Put(ns, kv.Key, kv.Value, version)
	}
	for _, p := range pvtdataHashes {
		updates.HashUpdates.Put(ns, p.coll, util.ComputeStringHash(p.key), util.ComputeHash(p.value), version)
	}
	txMgr.db.ApplyPrivacyAwareUpdates(updates, version)
}

type testutilPvtdata struct {
	coll, key string
	value     []byte
}
