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

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	btltestutil "github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flogging.ActivateSpec("pvtdatastorage=debug")
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
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 0,
			{"ns-3", "coll-1"}: 0,
			{"ns-4", "coll-1"}: 0,
			{"ns-4", "coll-2"}: 0,
		},
	)

	env := NewTestStoreEnv(t, "TestStoreBasicCommitAndRetrieval", btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore
	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}

	
	blk1MissingData := make(ledger.TxMissingPvtDataMap)

	
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-1", "coll-2", true)
	blk1MissingData.Add(1, "ns-2", "coll-1", true)
	blk1MissingData.Add(1, "ns-2", "coll-2", true)
	
	blk1MissingData.Add(2, "ns-3", "coll-1", true)
	
	blk1MissingData.Add(4, "ns-4", "coll-1", false)
	blk1MissingData.Add(4, "ns-4", "coll-2", false)

	
	blk2MissingData := make(ledger.TxMissingPvtDataMap)
	
	blk2MissingData.Add(1, "ns-1", "coll-1", true)
	blk2MissingData.Add(1, "ns-1", "coll-2", true)
	
	blk2MissingData.Add(3, "ns-1", "coll-1", true)

	
	assert.NoError(store.Prepare(0, nil, nil))
	assert.NoError(store.Commit())

	
	assert.NoError(store.Prepare(1, testData, blk1MissingData))
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

	
	assert.NoError(store.Prepare(2, testData, blk2MissingData))
	assert.NoError(store.Commit())

	
	
	

	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(2, 3, "ns-1", "coll-1")

	missingPvtDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(1)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-2")

	
	expectedMissingPvtDataInfo.Add(1, 2, "ns-3", "coll-1")

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(2)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)
}

func TestCommitPvtDataOfOldBlocks(t *testing.T) {
	viper.Set("ledger.pvtdataStore.purgeInterval", 2)
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 3,
			{"ns-1", "coll-2"}: 1,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 1,
			{"ns-3", "coll-1"}: 0,
			{"ns-3", "coll-2"}: 3,
			{"ns-4", "coll-1"}: 0,
			{"ns-4", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, "TestCommitPvtDataOfOldBlocks", btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore

	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}

	
	blk1MissingData := make(ledger.TxMissingPvtDataMap)

	
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-1", "coll-2", true)
	blk1MissingData.Add(1, "ns-2", "coll-1", true)
	blk1MissingData.Add(1, "ns-2", "coll-2", true)
	
	blk1MissingData.Add(2, "ns-1", "coll-1", true)
	blk1MissingData.Add(2, "ns-1", "coll-2", true)
	blk1MissingData.Add(2, "ns-3", "coll-1", true)
	blk1MissingData.Add(2, "ns-3", "coll-2", true)

	
	blk2MissingData := make(ledger.TxMissingPvtDataMap)
	
	blk2MissingData.Add(1, "ns-1", "coll-1", true)
	blk2MissingData.Add(1, "ns-1", "coll-2", true)
	
	blk2MissingData.Add(3, "ns-1", "coll-1", true)

	
	assert.NoError(store.Prepare(0, nil, nil))
	assert.NoError(store.Commit())

	
	assert.NoError(store.Prepare(1, testData, blk1MissingData))
	assert.NoError(store.Commit())

	
	assert.NoError(store.Prepare(2, nil, blk2MissingData))
	assert.NoError(store.Commit())

	
	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-2")

	
	expectedMissingPvtDataInfo.Add(1, 2, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 2, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 2, "ns-3", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 2, "ns-3", "coll-2")

	
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")
	
	expectedMissingPvtDataInfo.Add(2, 3, "ns-1", "coll-1")

	missingPvtDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(2)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	
	oldBlocksPvtData := make(map[uint64][]*ledger.TxPvtData)
	oldBlocksPvtData[1] = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 1, []string{"ns-1:coll-1", "ns-2:coll-1"}),
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-3:coll-1"}),
	}
	oldBlocksPvtData[2] = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-1"}),
	}

	err = store.CommitPvtDataOfOldBlocks(oldBlocksPvtData)
	assert.NoError(err)

	
	ns1Coll1Blk1Tx1 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-1", blkNum: 1}, txNum: 1}
	ns2Coll1Blk1Tx1 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-2", coll: "coll-1", blkNum: 1}, txNum: 1}
	ns1Coll1Blk1Tx2 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-1", blkNum: 1}, txNum: 2}
	ns3Coll1Blk1Tx2 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-3", coll: "coll-1", blkNum: 1}, txNum: 2}
	ns1Coll1Blk2Tx3 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-1", blkNum: 2}, txNum: 3}

	assert.True(testDataKeyExists(t, store, ns1Coll1Blk1Tx1))
	assert.True(testDataKeyExists(t, store, ns2Coll1Blk1Tx1))
	assert.True(testDataKeyExists(t, store, ns1Coll1Blk1Tx2))
	assert.True(testDataKeyExists(t, store, ns3Coll1Blk1Tx2))
	assert.True(testDataKeyExists(t, store, ns1Coll1Blk2Tx3))

	
	var nilFilter ledger.PvtNsCollFilter
	retrievedData, err := store.GetPvtDataByBlockNum(2, nilFilter)
	assert.NoError(err)
	for i, data := range retrievedData {
		assert.Equal(data.SeqInBlock, oldBlocksPvtData[2][i].SeqInBlock)
		assert.True(proto.Equal(data.WriteSet, oldBlocksPvtData[2][i].WriteSet))
	}

	expectedMissingPvtDataInfo = make(ledger.MissingPvtDataInfo)
	
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-2")

	
	expectedMissingPvtDataInfo.Add(1, 2, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 2, "ns-3", "coll-2")

	
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(2)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	expectedBlockList := []uint64{1, 2}
	blockList, err := store.GetLastUpdatedOldBlocksList()
	assert.NoError(err)
	assert.Equal(expectedBlockList, blockList)

	err = store.ResetLastUpdatedOldBlocksList()
	assert.NoError(err)

	expectedBlockList = []uint64(nil)
	blockList, err = store.GetLastUpdatedOldBlocksList()
	assert.NoError(err)
	assert.Nil(blockList)

	
	assert.NoError(store.Prepare(3, nil, nil))
	assert.NoError(store.Commit())

	
	
	oldBlocksPvtData = make(map[uint64][]*ledger.TxPvtData)
	oldBlocksPvtData[1] = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 1, []string{"ns-1:coll-2"}), 
		
		produceSamplePvtdata(t, 2, []string{"ns-3:coll-2"}), 
	}

	err = store.CommitPvtDataOfOldBlocks(oldBlocksPvtData)
	assert.NoError(err)

	ns1Coll2Blk1Tx1 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-2", blkNum: 1}, txNum: 1}
	ns2Coll2Blk1Tx1 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-2", coll: "coll-2", blkNum: 1}, txNum: 1}
	ns1Coll2Blk1Tx2 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-2", blkNum: 1}, txNum: 2}
	ns3Coll2Blk1Tx2 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-3", coll: "coll-2", blkNum: 1}, txNum: 2}

	
	
	
	assert.True(testDataKeyExists(t, store, ns1Coll2Blk1Tx1))  
	assert.False(testDataKeyExists(t, store, ns2Coll2Blk1Tx1)) 
	assert.False(testDataKeyExists(t, store, ns1Coll2Blk1Tx2)) 
	assert.True(testDataKeyExists(t, store, ns3Coll2Blk1Tx2))  

	err = store.ResetLastUpdatedOldBlocksList()
	assert.NoError(err)

	
	assert.NoError(store.Prepare(4, nil, nil))
	assert.NoError(store.Commit())

	testWaitForPurgerRoutineToFinish(store)

	
	
	oldBlocksPvtData = make(map[uint64][]*ledger.TxPvtData)
	oldBlocksPvtData[1] = []*ledger.TxPvtData{
		
		
		produceSamplePvtdata(t, 1, []string{"ns-2:coll-2"}),
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-2"}),
	}

	err = store.CommitPvtDataOfOldBlocks(oldBlocksPvtData)
	assert.NoError(err)

	ns1Coll2Blk1Tx1 = &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-2", blkNum: 1}, txNum: 1}
	ns2Coll2Blk1Tx1 = &dataKey{nsCollBlk: nsCollBlk{ns: "ns-2", coll: "coll-2", blkNum: 1}, txNum: 1}
	ns1Coll2Blk1Tx2 = &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-2", blkNum: 1}, txNum: 2}
	ns3Coll2Blk1Tx2 = &dataKey{nsCollBlk: nsCollBlk{ns: "ns-3", coll: "coll-2", blkNum: 1}, txNum: 2}

	assert.False(testDataKeyExists(t, store, ns1Coll2Blk1Tx1)) 
	assert.False(testDataKeyExists(t, store, ns2Coll2Blk1Tx1)) 
	assert.False(testDataKeyExists(t, store, ns1Coll2Blk1Tx2)) 
	assert.True(testDataKeyExists(t, store, ns3Coll2Blk1Tx2))  
}

func TestExpiryDataNotIncluded(t *testing.T) {
	ledgerid := "TestExpiryDataNotIncluded"
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 1,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 2,
			{"ns-3", "coll-1"}: 1,
			{"ns-3", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, ledgerid, btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore

	
	blk1MissingData := make(ledger.TxMissingPvtDataMap)
	
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-1", "coll-2", true)
	
	blk1MissingData.Add(4, "ns-3", "coll-1", false)
	blk1MissingData.Add(4, "ns-3", "coll-2", false)

	
	blk2MissingData := make(ledger.TxMissingPvtDataMap)
	
	blk2MissingData.Add(1, "ns-1", "coll-1", true)
	blk2MissingData.Add(1, "ns-1", "coll-2", true)

	
	assert.NoError(store.Prepare(0, nil, nil))
	assert.NoError(store.Commit())

	
	testDataForBlk1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	assert.NoError(store.Prepare(1, testDataForBlk1, blk1MissingData))
	assert.NoError(store.Commit())

	
	testDataForBlk2 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 5, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	assert.NoError(store.Prepare(2, testDataForBlk2, blk2MissingData))
	assert.NoError(store.Commit())

	retrievedData, _ := store.GetPvtDataByBlockNum(1, nil)
	
	for i, data := range retrievedData {
		assert.Equal(data.SeqInBlock, testDataForBlk1[i].SeqInBlock)
		assert.True(proto.Equal(data.WriteSet, testDataForBlk1[i].WriteSet))
	}

	
	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")

	
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")

	missingPvtDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	
	assert.NoError(store.Prepare(3, nil, nil))
	assert.NoError(store.Commit())

	
	expectedPvtdataFromBlock1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(1, nil)
	assert.Equal(expectedPvtdataFromBlock1, retrievedData)

	
	expectedMissingPvtDataInfo = make(ledger.MissingPvtDataInfo)
	
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")
	
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	
	assert.NoError(store.Prepare(4, nil, nil))
	assert.NoError(store.Commit())

	
	expectedPvtdataFromBlock1 = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-2", "ns-2:coll-1"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-2", "ns-2:coll-1"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(1, nil)
	assert.Equal(expectedPvtdataFromBlock1, retrievedData)

	
	expectedPvtdataFromBlock2 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 5, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(2, nil)
	assert.Equal(expectedPvtdataFromBlock2, retrievedData)

	
	expectedMissingPvtDataInfo = make(ledger.MissingPvtDataInfo)
	
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")

	
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

}

func TestStorePurge(t *testing.T) {
	ledgerid := "TestStorePurge"
	viper.Set("ledger.pvtdataStore.purgeInterval", 2)
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 1,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 4,
			{"ns-3", "coll-1"}: 1,
			{"ns-3", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, ledgerid, btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	s := env.TestStore

	
	assert.NoError(s.Prepare(0, nil, nil))
	assert.NoError(s.Commit())

	
	blk1MissingData := make(ledger.TxMissingPvtDataMap)
	
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-1", "coll-2", true)
	
	blk1MissingData.Add(4, "ns-3", "coll-1", false)
	blk1MissingData.Add(4, "ns-3", "coll-2", false)

	
	testDataForBlk1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	assert.NoError(s.Prepare(1, testDataForBlk1, blk1MissingData))
	assert.NoError(s.Commit())

	
	assert.NoError(s.Prepare(2, nil, nil))
	assert.NoError(s.Commit())
	
	ns1Coll1 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-1", blkNum: 1}, txNum: 2}
	ns2Coll2 := &dataKey{nsCollBlk: nsCollBlk{ns: "ns-2", coll: "coll-2", blkNum: 1}, txNum: 2}

	
	ns1Coll1elgMD := &missingDataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-1", blkNum: 1}, isEligible: true}
	ns1Coll2elgMD := &missingDataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-2", blkNum: 1}, isEligible: true}

	
	ns3Coll1inelgMD := &missingDataKey{nsCollBlk: nsCollBlk{ns: "ns-3", coll: "coll-1", blkNum: 1}, isEligible: false}
	ns3Coll2inelgMD := &missingDataKey{nsCollBlk: nsCollBlk{ns: "ns-3", coll: "coll-2", blkNum: 1}, isEligible: false}

	testWaitForPurgerRoutineToFinish(s)
	assert.True(testDataKeyExists(t, s, ns1Coll1))
	assert.True(testDataKeyExists(t, s, ns2Coll2))

	assert.True(testMissingDataKeyExists(t, s, ns1Coll1elgMD))
	assert.True(testMissingDataKeyExists(t, s, ns1Coll2elgMD))

	assert.True(testMissingDataKeyExists(t, s, ns3Coll1inelgMD))
	assert.True(testMissingDataKeyExists(t, s, ns3Coll2inelgMD))

	
	assert.NoError(s.Prepare(3, nil, nil))
	assert.NoError(s.Commit())
	
	testWaitForPurgerRoutineToFinish(s)
	assert.True(testDataKeyExists(t, s, ns1Coll1))
	assert.True(testDataKeyExists(t, s, ns2Coll2))
	
	assert.True(testMissingDataKeyExists(t, s, ns1Coll1elgMD))
	assert.True(testMissingDataKeyExists(t, s, ns1Coll2elgMD))
	
	assert.True(testMissingDataKeyExists(t, s, ns3Coll1inelgMD))
	assert.True(testMissingDataKeyExists(t, s, ns3Coll2inelgMD))

	
	assert.NoError(s.Prepare(4, nil, nil))
	assert.NoError(s.Commit())
	
	
	testWaitForPurgerRoutineToFinish(s)
	assert.False(testDataKeyExists(t, s, ns1Coll1))
	assert.True(testDataKeyExists(t, s, ns2Coll2))
	
	assert.False(testMissingDataKeyExists(t, s, ns1Coll1elgMD))
	assert.True(testMissingDataKeyExists(t, s, ns1Coll2elgMD))
	
	assert.False(testMissingDataKeyExists(t, s, ns3Coll1inelgMD))
	assert.True(testMissingDataKeyExists(t, s, ns3Coll2inelgMD))

	
	assert.NoError(s.Prepare(5, nil, nil))
	assert.NoError(s.Commit())
	
	testWaitForPurgerRoutineToFinish(s)
	assert.False(testDataKeyExists(t, s, ns1Coll1))
	assert.True(testDataKeyExists(t, s, ns2Coll2))

	
	assert.NoError(s.Prepare(6, nil, nil))
	assert.NoError(s.Commit())
	
	testWaitForPurgerRoutineToFinish(s)
	assert.False(testDataKeyExists(t, s, ns1Coll1))
	assert.False(testDataKeyExists(t, s, ns2Coll2))

	
	assert.True(testDataKeyExists(t, s, &dataKey{nsCollBlk: nsCollBlk{ns: "ns-1", coll: "coll-2", blkNum: 1}, txNum: 2}))
}

func TestStoreState(t *testing.T) {
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
		},
	)
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

func TestCollElgEnabled(t *testing.T) {
	testCollElgEnabled(t)
	defaultValBatchSize := ledgerconfig.GetPvtdataStoreCollElgProcMaxDbBatchSize()
	defaultValInterval := ledgerconfig.GetPvtdataStoreCollElgProcDbBatchesInterval()
	defer func() {
		viper.Set("ledger.pvtdataStore.collElgProcMaxDbBatchSize", defaultValBatchSize)
		viper.Set("ledger.pvtdataStore.collElgProcMaxDbBatchSize", defaultValInterval)
	}()
	viper.Set("ledger.pvtdataStore.collElgProcMaxDbBatchSize", 1)
	viper.Set("ledger.pvtdataStore.collElgProcDbBatchesInterval", 1)
	testCollElgEnabled(t)
}

func testCollElgEnabled(t *testing.T) {
	ledgerid := "TestCollElgEnabled"
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, ledgerid, btlPolicy)
	defer env.Cleanup()
	assert := assert.New(t)
	store := env.TestStore

	

	
	assert.NoError(store.Prepare(0, nil, nil))
	assert.NoError(store.Commit())

	
	blk1MissingData := make(ledger.TxMissingPvtDataMap)
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-2", "coll-1", true)
	blk1MissingData.Add(4, "ns-1", "coll-2", false)
	blk1MissingData.Add(4, "ns-2", "coll-2", false)
	testDataForBlk1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1"}),
	}
	assert.NoError(store.Prepare(1, testDataForBlk1, blk1MissingData))
	assert.NoError(store.Commit())

	
	blk2MissingData := make(ledger.TxMissingPvtDataMap)
	
	blk2MissingData.Add(1, "ns-1", "coll-2", false)
	blk2MissingData.Add(1, "ns-2", "coll-2", false)
	testDataForBlk2 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-1"}),
	}
	assert.NoError(store.Prepare(2, testDataForBlk2, blk2MissingData))
	assert.NoError(store.Commit())

	
	
	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-1")
	missingPvtDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	
	store.ProcessCollsEligibilityEnabled(
		5,
		map[string][]string{
			"ns-1": {"coll-2"},
		},
	)
	testutilWaitForCollElgProcToFinish(store)

	
	
	expectedMissingPvtDataInfo.Add(1, 4, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")
	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	assert.NoError(err)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	
	store.ProcessCollsEligibilityEnabled(6,
		map[string][]string{
			"ns-2": {"coll-2"},
		},
	)
	testutilWaitForCollElgProcToFinish(store)

	
	
	expectedMissingPvtDataInfo.Add(1, 4, "ns-2", "coll-2")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-2", "coll-2")
	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	assert.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)
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

func testMissingDataKeyExists(t *testing.T, s Store, missingDataKey *missingDataKey) bool {
	dataKeyBytes := encodeMissingDataKey(missingDataKey)
	val, err := s.(*store).db.Get(dataKeyBytes)
	assert.NoError(t, err)
	return len(val) != 0
}

func testWaitForPurgerRoutineToFinish(s Store) {
	time.Sleep(1 * time.Second)
	s.(*store).purgerLock.Lock()
	s.(*store).purgerLock.Unlock()
}

func testutilWaitForCollElgProcToFinish(s Store) {
	s.(*store).collElgProcSync.waitForDone()
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
