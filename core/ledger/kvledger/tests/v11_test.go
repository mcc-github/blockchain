/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"fmt"
	"os"
	"testing"

	protopeer "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestV11(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	
	testutil.CopyDir("testdata/v11/sample_ledgers/ledgersData", env.initializer.Config.RootFSPath, true)

	blkIndexPath := fmt.Sprintf("%s/chains/index", env.initializer.Config.RootFSPath)
	historyDBPath := fmt.Sprintf("%s/historyLeveldb", env.initializer.Config.RootFSPath)
	idStorePath := kvledger.LedgerProviderPath(env.initializer.Config.RootFSPath)

	
	assert.PanicsWithValue(
		t,
		fmt.Sprintf("Error in instantiating ledger provider: unexpected format. db path = [%s], data format = [], expected format = [2.0]",
			idStorePath),
		func() { env.initLedgerMgmt() },
	)

	
	require.NoError(t, kvledger.UpgradeIDStoreFormat(env.initializer.Config.RootFSPath))
	
	require.PanicsWithValue(
		t,
		fmt.Sprintf("Error in instantiating ledger provider: unexpected format. db path = [%s], data format = [], expected format = [2.0]",
			blkIndexPath),
		func() { env.initLedgerMgmt() },
	)

	
	require.NoError(t, os.RemoveAll(blkIndexPath))
	require.PanicsWithValue(
		t,
		fmt.Sprintf("Error in instantiating ledger provider: unexpected format. db path = [%s], data format = [], expected format = [2.0]",
			historyDBPath),
		func() { env.initLedgerMgmt() },
	)

	
	require.NoError(t, os.RemoveAll(historyDBPath))
	env.initLedgerMgmt()

	h1, h2 := env.newTestHelperOpenLgr("ledger1", t), env.newTestHelperOpenLgr("ledger2", t)
	dataHelper := &v11SampleDataHelper{}

	dataHelper.verifyBeforeStateRebuild(h1)
	dataHelper.verifyBeforeStateRebuild(h2)

	env.closeAllLedgersAndDrop(rebuildableStatedb)
	h1, h2 = env.newTestHelperOpenLgr("ledger1", t), env.newTestHelperOpenLgr("ledger2", t)
	dataHelper.verifyAfterStateRebuild(h1)
	dataHelper.verifyAfterStateRebuild(h2)

	env.closeAllLedgersAndDrop(rebuildableStatedb + rebuildableBlockIndex + rebuildableConfigHistory + rebuildableHistoryDB)
	h1, h2 = env.newTestHelperOpenLgr("ledger1", t), env.newTestHelperOpenLgr("ledger2", t)
	dataHelper.verifyAfterStateRebuild(h1)
	dataHelper.verifyAfterStateRebuild(h2)
}





type v11SampleDataHelper struct {
}

func (d *v11SampleDataHelper) verifyBeforeStateRebuild(h *testhelper) {
	dataHelper := &v11SampleDataHelper{}
	dataHelper.verifyState(h)
	dataHelper.verifyBlockAndPvtdata(h)
	dataHelper.verifyGetTransactionByID(h)
	dataHelper.verifyConfigHistoryDoesNotExist(h)
	dataHelper.verifyHistory(h)
}

func (d *v11SampleDataHelper) verifyAfterStateRebuild(h *testhelper) {
	dataHelper := &v11SampleDataHelper{}
	dataHelper.verifyState(h)
	dataHelper.verifyBlockAndPvtdata(h)
	dataHelper.verifyGetTransactionByID(h)
	dataHelper.verifyConfigHistory(h)
	dataHelper.verifyHistory(h)
}

func (d *v11SampleDataHelper) verifyState(h *testhelper) {
	lgrid := h.lgrid
	h.verifyPubState("cc1", "key1", d.sampleVal("value13", lgrid))
	h.verifyPubState("cc1", "key2", "")
	h.verifyPvtState("cc1", "coll1", "key3", d.sampleVal("value14", lgrid))
	h.verifyPvtState("cc1", "coll1", "key4", "")
	h.verifyPvtState("cc1", "coll2", "key3", d.sampleVal("value09", lgrid))
	h.verifyPvtState("cc1", "coll2", "key4", d.sampleVal("value10", lgrid))

	h.verifyPubState("cc2", "key1", d.sampleVal("value03", lgrid))
	h.verifyPubState("cc2", "key2", d.sampleVal("value04", lgrid))
	h.verifyPvtState("cc2", "coll1", "key3", d.sampleVal("value07", lgrid))
	h.verifyPvtState("cc2", "coll1", "key4", d.sampleVal("value08", lgrid))
	h.verifyPvtState("cc2", "coll2", "key3", d.sampleVal("value11", lgrid))
	h.verifyPvtState("cc2", "coll2", "key4", d.sampleVal("value12", lgrid))
}

func (d *v11SampleDataHelper) verifyHistory(h *testhelper) {
	lgrid := h.lgrid
	expectedHistoryCC1Key1 := []string{
		d.sampleVal("value13", lgrid),
		d.sampleVal("value01", lgrid),
	}
	h.verifyHistory("cc1", "key1", expectedHistoryCC1Key1)
}

func (d *v11SampleDataHelper) verifyConfigHistory(h *testhelper) {
	lgrid := h.lgrid
	h.verifyMostRecentCollectionConfigBelow(10, "cc1",
		&expectedCollConfInfo{5, d.sampleCollConf2(lgrid, "cc1")})

	h.verifyMostRecentCollectionConfigBelow(5, "cc1",
		&expectedCollConfInfo{3, d.sampleCollConf1(lgrid, "cc1")})

	h.verifyMostRecentCollectionConfigBelow(10, "cc2",
		&expectedCollConfInfo{5, d.sampleCollConf2(lgrid, "cc2")})

	h.verifyMostRecentCollectionConfigBelow(5, "cc2",
		&expectedCollConfInfo{3, d.sampleCollConf1(lgrid, "cc2")})
}

func (d *v11SampleDataHelper) verifyConfigHistoryDoesNotExist(h *testhelper) {
	h.verifyMostRecentCollectionConfigBelow(10, "cc1", nil)
	h.verifyMostRecentCollectionConfigBelow(10, "cc2", nil)
}

func (d *v11SampleDataHelper) verifyBlockAndPvtdata(h *testhelper) {
	lgrid := h.lgrid
	h.verifyBlockAndPvtData(2, nil, func(r *retrievedBlockAndPvtdata) {
		r.hasNumTx(2)
		r.hasNoPvtdata()
	})

	h.verifyBlockAndPvtData(4, nil, func(r *retrievedBlockAndPvtdata) {
		r.hasNumTx(2)
		r.pvtdataShouldContain(0, "cc1", "coll1", "key3", d.sampleVal("value05", lgrid))
		r.pvtdataShouldContain(1, "cc2", "coll1", "key3", d.sampleVal("value07", lgrid))
	})
}

func (d *v11SampleDataHelper) verifyGetTransactionByID(h *testhelper) {
	h.verifyTxValidationCode("txid7", protopeer.TxValidationCode_VALID)
	h.verifyTxValidationCode("txid8", protopeer.TxValidationCode_MVCC_READ_CONFLICT)
}

func (d *v11SampleDataHelper) sampleVal(val, ledgerid string) string {
	return fmt.Sprintf("%s:%s", val, ledgerid)
}

func (d *v11SampleDataHelper) sampleCollConf1(ledgerid, ccName string) []*collConf {
	return []*collConf{
		{name: "coll1", members: []string{"org1", "org2"}},
		{name: ledgerid, members: []string{"org1", "org2"}},
		{name: ccName, members: []string{"org1", "org2"}},
	}
}

func (d *v11SampleDataHelper) sampleCollConf2(ledgerid string, ccName string) []*collConf {
	return []*collConf{
		{name: "coll1", members: []string{"org1", "org2"}},
		{name: "coll2", members: []string{"org1", "org2"}},
		{name: ledgerid, members: []string{"org1", "org2"}},
		{name: ccName, members: []string{"org1", "org2"}},
	}
}
