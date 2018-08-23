/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	btltestutil "github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flogging.SetModuleLevel("pvtdatastorage", "debug")
	viper.Set("peer.fileSystemPath", "/tmp/blockchain/core/ledger/pvtdatastorage")
	os.Exit(m.Run())
}

func TestEmptyStore(t *testing.T) {
	env := NewTestStoreEnv(t, "TestEmptyStore", nil)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore
	testEmpty(true, assert, store)
	testPendingBatch(false, assert, store)
}

func TestStoreBasicCommitAndRetrieval(t *testing.T) {
	cs := btltestutil.NewMockCollectionStore()
	cs.SetBTL("ns-1", "coll-1", 0)
	cs.SetBTL("ns-1", "coll-2", 0)
	cs.SetBTL("ns-2", "coll-1", 0)
	cs.SetBTL("ns-2", "coll-2", 0)
	btlPolicy := pvtdatapolicy.ConstructBTLPolicy(cs)
	env := NewTestStoreEnv(t, "TestStoreBasicCommitAndRetrieval", btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore
	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}

	
	assert.NoError(store.Prepare(0, nil, nil))
	assert.NoError(store.Commit())

	
	assert.NoError(store.Prepare(1, testData, nil))
	assert.NoError(store.Commit())

	
	assert.NoError(store.Prepare(2, testData, nil))
	assert.NoError(store.Rollback())

	
	var nilFilter ledger.PvtNsCollFilter
	retrievedData, err := store.GetPvtDataByBlockNum(0, nilFilter)
	assert.NoError(err)
	assert.Nil(retrievedData)

	
	retrievedData, err = store.GetPvtDataByBlockNum(1, nilFilter)
	assert.NoError(err)
	for i, data := range retrievedData {
		assert.Equal(data.SeqInBlock, testData[i].SeqInBlock)
		assert.True(proto.Equal(data.WriteSet, testData[i].WriteSet))
	}

	
	filter := ledger.NewPvtNsCollFilter()
	filter.Add("ns-1", "coll-1")
	filter.Add("ns-2", "coll-2")
	retrievedData, err = store.GetPvtDataByBlockNum(1, filter)
	expectedRetrievedData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-2:coll-2"}),
	}
	for i, data := range retrievedData {
		assert.Equal(data.SeqInBlock, expectedRetrievedData[i].SeqInBlock)
		assert.True(proto.Equal(data.WriteSet, expectedRetrievedData[i].WriteSet))
	}

	
	retrievedData, err = store.GetPvtDataByBlockNum(2, nilFilter)
	_, ok := err.(*ErrOutOfRange)
	assert.True(ok)
	assert.Nil(retrievedData)
}

func TestExpiryDataNotIncluded(t *testing.T) {
	ledgerid := "TestExpiryDataNotIncluded"
	cs := btltestutil.NewMockCollectionStore()
	cs.SetBTL("ns-1", "coll-1", 1)
	cs.SetBTL("ns-1", "coll-2", 0)
	cs.SetBTL("ns-2", "coll-1", 0)
	cs.SetBTL("ns-2", "coll-2", 2)
	btlPolicy := pvtdatapolicy.ConstructBTLPolicy(cs)

	env := NewTestStoreEnv(t, ledgerid, btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore

	
	assert.NoError(store.Prepare(0, nil, nil))
	assert.NoError(store.Commit())

	
	testDataForBlk1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	assert.NoError(store.Prepare(1, testDataForBlk1, nil))
	assert.NoError(store.Commit())

	
	testDataForBlk2 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 5, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	assert.NoError(store.Prepare(2, testDataForBlk2, nil))
	assert.NoError(store.Commit())

	retrievedData, _ := store.GetPvtDataByBlockNum(1, nil)
	
	for i, data := range retrievedData {
		assert.Equal(data.SeqInBlock, testDataForBlk1[i].SeqInBlock)
		assert.True(proto.Equal(data.WriteSet, testDataForBlk1[i].WriteSet))
	}

	
	assert.NoError(store.Prepare(3, nil, nil))
	assert.NoError(store.Commit())

	
	expectedPvtdataFromBlock1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(1, nil)
	testutil.AssertEquals(t, retrievedData, expectedPvtdataFromBlock1)

	
	assert.NoError(store.Prepare(4, nil, nil))
	assert.NoError(store.Commit())

	
	expectedPvtdataFromBlock1 = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-2", "ns-2:coll-1"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-2", "ns-2:coll-1"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(1, nil)
	testutil.AssertEquals(t, retrievedData, expectedPvtdataFromBlock1)

	
	expectedPvtdataFromBlock2 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 5, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(2, nil)
	testutil.AssertEquals(t, retrievedData, expectedPvtdataFromBlock2)
}

func TestStorePurge(t *testing.T) {
	ledgerid := "TestStorePurge"
	viper.Set("ledger.pvtdataStore.purgeInterval", 2)
	cs := btltestutil.NewMockCollectionStore()
	cs.SetBTL("ns-1", "coll-1", 1)
	cs.SetBTL("ns-1", "coll-2", 0)
	cs.SetBTL("ns-2", "coll-1", 0)
	cs.SetBTL("ns-2", "coll-2", 4)
	btlPolicy := pvtdatapolicy.ConstructBTLPolicy(cs)

	env := NewTestStoreEnv(t, ledgerid, btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	s := env.TestStore

	
	assert.NoError(s.Prepare(0, nil, nil))
	assert.NoError(s.Commit())

	
	testDataForBlk1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	assert.NoError(s.Prepare(1, testDataForBlk1, nil))
	assert.NoError(s.Commit())

	
	assert.NoError(s.Prepare(2, nil, nil))
	assert.NoError(s.Commit())
	
	testWaitForPurgerRoutineToFinish(s)
	assert.True(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-1", coll: "coll-1"}))
	assert.True(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-2", coll: "coll-2"}))

	
	assert.NoError(s.Prepare(3, nil, nil))
	assert.NoError(s.Commit())
	
	testWaitForPurgerRoutineToFinish(s)
	assert.True(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-1", coll: "coll-1"}))
	assert.True(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-2", coll: "coll-2"}))

	
	assert.NoError(s.Prepare(4, nil, nil))
	assert.NoError(s.Commit())
	
	
	testWaitForPurgerRoutineToFinish(s)
	assert.False(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-1", coll: "coll-1"}))
	assert.True(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-2", coll: "coll-2"}))

	
	assert.NoError(s.Prepare(5, nil, nil))
	assert.NoError(s.Commit())
	
	testWaitForPurgerRoutineToFinish(s)
	assert.False(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-1", coll: "coll-1"}))
	assert.True(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-2", coll: "coll-2"}))

	
	assert.NoError(s.Prepare(6, nil, nil))
	assert.NoError(s.Commit())
	
	testWaitForPurgerRoutineToFinish(s)
	assert.False(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-1", coll: "coll-1"}))
	assert.False(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-2", coll: "coll-2"}))

	
	assert.True(testDataKeyExists(t, s, &dataKey{blkNum: 1, txNum: 2, ns: "ns-1", coll: "coll-2"}))
}

func TestStoreState(t *testing.T) {
	cs := btltestutil.NewMockCollectionStore()
	cs.SetBTL("ns-1", "coll-1", 0)
	cs.SetBTL("ns-1", "coll-2", 0)
	btlPolicy := pvtdatapolicy.ConstructBTLPolicy(cs)

	env := NewTestStoreEnv(t, "TestStoreState", btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore
	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 0, []string{"ns-1:coll-1", "ns-1:coll-2"}),
	}
	_, ok := store.Prepare(1, testData, nil).(*ErrIllegalArgs)
	assert.True(ok)

	assert.Nil(store.Prepare(0, testData, nil))
	assert.NoError(store.Commit())

	assert.Nil(store.Prepare(1, testData, nil))
	_, ok = store.Prepare(2, testData, nil).(*ErrIllegalCall)
	assert.True(ok)
}

func TestInitLastCommittedBlock(t *testing.T) {
	env := NewTestStoreEnv(t, "TestStoreState", nil)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore
	existingLastBlockNum := uint64(25)
	assert.NoError(store.InitLastCommittedBlock(existingLastBlockNum))

	testEmpty(false, assert, store)
	testPendingBatch(false, assert, store)
	testLastCommittedBlockHeight(existingLastBlockNum+1, assert, store)

	env.CloseAndReopen()
	testEmpty(false, assert, store)
	testPendingBatch(false, assert, store)
	testLastCommittedBlockHeight(existingLastBlockNum+1, assert, store)

	err := store.InitLastCommittedBlock(30)
	_, ok := err.(*ErrIllegalCall)
	assert.True(ok)
}



func testEmpty(expectedEmpty bool, assert *assert.Assertions, store Store) {
	isEmpty, err := store.IsEmpty()
	assert.NoError(err)
	assert.Equal(expectedEmpty, isEmpty)
}

func testPendingBatch(expectedPending bool, assert *assert.Assertions, store Store) {
	hasPendingBatch, err := store.HasPendingBatch()
	assert.NoError(err)
	assert.Equal(expectedPending, hasPendingBatch)
}

func testLastCommittedBlockHeight(expectedBlockHt uint64, assert *assert.Assertions, store Store) {
	blkHt, err := store.LastCommittedBlockHeight()
	assert.NoError(err)
	assert.Equal(expectedBlockHt, blkHt)
}

func testDataKeyExists(t *testing.T, s Store, dataKey *dataKey) bool {
	dataKeyBytes := encodeDataKey(dataKey)
	val, err := s.(*store).db.Get(dataKeyBytes)
	assert.NoError(t, err)
	return len(val) != 0
}

func testWaitForPurgerRoutineToFinish(s Store) {
	time.Sleep(1 * time.Second)
	s.(*store).purgerLock.Lock()
	s.(*store).purgerLock.Unlock()
}

func produceSamplePvtdata(t *testing.T, txNum uint64, nsColls []string) *ledger.TxPvtData {
	builder := rwsetutil.NewRWSetBuilder()
	for _, nsColl := range nsColls {
		nsCollSplit := strings.Split(nsColl, ":")
		ns := nsCollSplit[0]
		coll := nsCollSplit[1]
		builder.AddToPvtAndHashedWriteSet(ns, coll, fmt.Sprintf("key-%s-%s", ns, coll), []byte(fmt.Sprintf("value-%s-%s", ns, coll)))
	}
	simRes, err := builder.GetTxSimulationResults()
	assert.NoError(t, err)
	return &ledger.TxPvtData{SeqInBlock: txNum, WriteSet: simRes.PvtSimulationResults}
}
