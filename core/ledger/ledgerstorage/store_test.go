/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerstorage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage/fsblkstorage"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	btltestutil "github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/testutil"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatastorage"
	lutil "github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/stretchr/testify/assert"
)

var metricsProvider = &disabled.Provider{}

func TestMain(m *testing.M) {
	flogging.ActivateSpec("ledgerstorage,pvtdatastorage=debug")
	os.Exit(m.Run())
}

func TestStore(t *testing.T) {
	storeDir, err := ioutil.TempDir("", "lstore")
	if err != nil {
		t.Fatalf("Failed to create ledger storage directory: %s", err)
	}
	defer os.RemoveAll(storeDir)
	conf := buildPrivateDataConfig(storeDir)
	blockStoreDir := filepath.Join(storeDir, "chains")
	provider, err := NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	defer provider.Close()
	store, err := provider.Open("testLedger")
	store.Init(btlPolicyForSampleData())
	defer store.Shutdown()

	assert.NoError(t, err)
	sampleData := sampleDataWithPvtdataForSelectiveTx(t)
	for _, sampleDatum := range sampleData {
		assert.NoError(t, store.CommitWithPvtData(sampleDatum))
	}

	
	pvtdata, err := store.GetPvtDataByNum(1, nil)
	assert.NoError(t, err)
	assert.Nil(t, pvtdata)

	
	pvtdata, err = store.GetPvtDataByNum(4, nil)
	assert.NoError(t, err)
	assert.Nil(t, pvtdata)

	
	
	
	pvtdata, err = store.GetPvtDataByNum(2, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(pvtdata))
	assert.Equal(t, uint64(3), pvtdata[0].SeqInBlock)
	assert.Equal(t, uint64(5), pvtdata[1].SeqInBlock)
	assert.Equal(t, uint64(6), pvtdata[2].SeqInBlock)

	
	pvtdata, err = store.GetPvtDataByNum(3, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(pvtdata))
	assert.Equal(t, uint64(4), pvtdata[0].SeqInBlock)
	assert.Equal(t, uint64(6), pvtdata[1].SeqInBlock)

	blockAndPvtdata, err := store.GetPvtDataAndBlockByNum(2, nil)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(sampleData[2].Block, blockAndPvtdata.Block))

	blockAndPvtdata, err = store.GetPvtDataAndBlockByNum(3, nil)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(sampleData[3].Block, blockAndPvtdata.Block))

	
	filter := ledger.NewPvtNsCollFilter()
	filter.Add("ns-1", "coll-1")
	blockAndPvtdata, err = store.GetPvtDataAndBlockByNum(3, filter)
	assert.NoError(t, err)
	assert.Equal(t, sampleData[3].Block, blockAndPvtdata.Block)
	
	assert.Equal(t, 2, len(blockAndPvtdata.PvtData))
	
	assert.Equal(t, 1, len(blockAndPvtdata.PvtData[4].WriteSet.NsPvtRwset))
	assert.Equal(t, 1, len(blockAndPvtdata.PvtData[6].WriteSet.NsPvtRwset))
	
	assert.Nil(t, blockAndPvtdata.PvtData[2])

	
	
	
	expectedMissingDataInfo := make(ledger.MissingPvtDataInfo)
	expectedMissingDataInfo.Add(5, 4, "ns-4", "coll-4")
	expectedMissingDataInfo.Add(5, 5, "ns-5", "coll-5")
	missingDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedMissingDataInfo, missingDataInfo)
}

func TestStoreWithExistingBlockchain(t *testing.T) {
	testLedgerid := "test-ledger"
	storeDir, err := ioutil.TempDir("", "lstore")
	if err != nil {
		t.Fatalf("Failed to create ledger storage directory: %s", err)
	}
	defer os.RemoveAll(storeDir)

	
	attrsToIndex := []blkstorage.IndexableAttr{
		blkstorage.IndexableAttrBlockHash,
		blkstorage.IndexableAttrBlockNum,
		blkstorage.IndexableAttrTxID,
		blkstorage.IndexableAttrBlockNumTranNum,
	}
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	blockStoreProvider, err := fsblkstorage.NewProvider(
		fsblkstorage.NewConf(filepath.Join(storeDir, "chains"), maxBlockFileSize),
		indexConfig,
		metricsProvider,
	)
	assert.NoError(t, err)
	blkStore, err := blockStoreProvider.OpenBlockStore(testLedgerid)
	assert.NoError(t, err)
	testBlocks := testutil.ConstructTestBlocks(t, 10)

	existingBlocks := testBlocks[0:9]
	blockToAdd := testBlocks[9:][0]

	
	for _, blk := range existingBlocks {
		assert.NoError(t, blkStore.AddBlock(blk))
	}
	blockStoreProvider.Close()

	
	
	conf := buildPrivateDataConfig(storeDir)
	blockStoreDir := filepath.Join(storeDir, "chains")
	provider, err := NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	defer provider.Close()
	store, err := provider.Open(testLedgerid)
	store.Init(btlPolicyForSampleData())
	defer store.Shutdown()

	
	pvtdataBlockHt, err := store.pvtdataStore.LastCommittedBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(9), pvtdataBlockHt)

	
	pvtdata := samplePvtData(t, []uint64{0})
	assert.NoError(t, store.CommitWithPvtData(&ledger.BlockAndPvtData{Block: blockToAdd, PvtData: pvtdata}))
	pvtdataBlockHt, err = store.pvtdataStore.LastCommittedBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(10), pvtdataBlockHt)
}

func TestCrashAfterPvtdataStoreCommit(t *testing.T) {
	storeDir, err := ioutil.TempDir("", "lstore")
	if err != nil {
		t.Fatalf("Failed to create ledger storage directory: %s", err)
	}
	defer os.RemoveAll(storeDir)
	conf := buildPrivateDataConfig(storeDir)
	blockStoreDir := filepath.Join(storeDir, "chains")
	provider, err := NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	defer provider.Close()
	store, err := provider.Open("testLedger")
	store.Init(btlPolicyForSampleData())
	defer store.Shutdown()
	assert.NoError(t, err)

	sampleData := sampleDataWithPvtdataForAllTxs(t)
	dataBeforeCrash := sampleData[0:3]
	dataAtCrash := sampleData[3]

	for _, sampleDatum := range dataBeforeCrash {
		assert.NoError(t, store.CommitWithPvtData(sampleDatum))
	}
	blockNumAtCrash := dataAtCrash.Block.Header.Number
	var pvtdataAtCrash []*ledger.TxPvtData
	for _, p := range dataAtCrash.PvtData {
		pvtdataAtCrash = append(pvtdataAtCrash, p)
	}
	
	store.pvtdataStore.Commit(blockNumAtCrash, pvtdataAtCrash, nil)
	store.Shutdown()
	provider.Close()
	provider, err = NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	store, err = provider.Open("testLedger")
	assert.NoError(t, err)
	store.Init(btlPolicyForSampleData())

	
	testVerifyPvtData(t, store, blockNumAtCrash, dataAtCrash.PvtData)
	bcInfo, err := store.GetBlockchainInfo()
	assert.NoError(t, err)
	assert.Equal(t, blockNumAtCrash, bcInfo.Height)

	
	
	
	expectedPvtData := dataAtCrash.PvtData
	dataAtCrash.PvtData = make(ledger.TxPvtDataMap)
	dataAtCrash.PvtData[0] = &ledger.TxPvtData{
		SeqInBlock: 0,
		WriteSet: &rwset.TxPvtReadWriteSet{
			NsPvtRwset: []*rwset.NsPvtReadWriteSet{
				{
					Namespace: "ns-1",
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
						{
							CollectionName: "coll-1",
							Rwset:          []byte("pvtdata"),
						},
					},
				},
			},
		},
	}
	assert.NoError(t, store.CommitWithPvtData(dataAtCrash))
	testVerifyPvtData(t, store, blockNumAtCrash, expectedPvtData)
	bcInfo, err = store.GetBlockchainInfo()
	assert.NoError(t, err)
	assert.Equal(t, blockNumAtCrash+1, bcInfo.Height)

}

func testVerifyPvtData(t *testing.T, store *Store, blockNum uint64, expectedPvtData ledger.TxPvtDataMap) {
	pvtdata, err := store.GetPvtDataByNum(blockNum, nil)
	assert.NoError(t, err)
	constructed := constructPvtdataMap(pvtdata)
	assert.Equal(t, len(expectedPvtData), len(constructed))
	for k, v := range expectedPvtData {
		ov, ok := constructed[k]
		assert.True(t, ok)
		assert.Equal(t, v.SeqInBlock, ov.SeqInBlock)
		assert.True(t, proto.Equal(v.WriteSet, ov.WriteSet))
	}
}

func TestAddAfterPvtdataStoreError(t *testing.T) {
	storeDir, err := ioutil.TempDir("", "lstore")
	if err != nil {
		t.Fatalf("Failed to create ledger storage directory: %s", err)
	}
	defer os.RemoveAll(storeDir)
	conf := buildPrivateDataConfig(storeDir)
	blockStoreDir := filepath.Join(storeDir, "chains")
	provider, err := NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	defer provider.Close()
	store, err := provider.Open("testLedger")
	store.Init(btlPolicyForSampleData())
	defer store.Shutdown()
	assert.NoError(t, err)

	sampleData := sampleDataWithPvtdataForAllTxs(t)
	for _, d := range sampleData[0:9] {
		assert.NoError(t, store.CommitWithPvtData(d))
	}
	
	
	assert.Error(t, store.CommitWithPvtData(sampleData[8]))

	
	pvtStoreCommitHt, err := store.pvtdataStore.LastCommittedBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(9), pvtStoreCommitHt)

	
	assert.NoError(t, store.CommitWithPvtData(sampleData[9]))
	pvtStoreCommitHt, err = store.pvtdataStore.LastCommittedBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(10), pvtStoreCommitHt)
}

func TestAddAfterBlkStoreError(t *testing.T) {
	storeDir, err := ioutil.TempDir("", "lstore")
	if err != nil {
		t.Fatalf("Failed to create ledger storage directory: %s", err)
	}
	defer os.RemoveAll(storeDir)
	conf := buildPrivateDataConfig(storeDir)
	blockStoreDir := filepath.Join(storeDir, "chains")
	provider, err := NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	defer provider.Close()
	store, err := provider.Open("testLedger")
	store.Init(btlPolicyForSampleData())
	defer store.Shutdown()
	assert.NoError(t, err)

	sampleData := sampleDataWithPvtdataForAllTxs(t)
	for _, d := range sampleData[0:9] {
		assert.NoError(t, store.CommitWithPvtData(d))
	}
	lastBlkAndPvtData := sampleData[9]
	
	store.BlockStore.AddBlock(lastBlkAndPvtData.Block)
	
	assert.Error(t, store.CommitWithPvtData(lastBlkAndPvtData))
	
	pvtStoreCommitHt, err := store.pvtdataStore.LastCommittedBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(10), pvtStoreCommitHt)
}

func TestPvtStoreAheadOfBlockStore(t *testing.T) {
	storeDir, err := ioutil.TempDir("", "lstore")
	if err != nil {
		t.Fatalf("Failed to create ledger storage directory: %s", err)
	}
	defer os.RemoveAll(storeDir)
	conf := buildPrivateDataConfig(storeDir)
	blockStoreDir := filepath.Join(storeDir, "chains")
	provider, err := NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	defer provider.Close()
	store, err := provider.Open("testLedger")
	store.Init(btlPolicyForSampleData())
	defer store.Shutdown()
	assert.NoError(t, err)

	
	assert.False(t, store.IsPvtStoreAheadOfBlockStore())

	sampleData := sampleDataWithPvtdataForSelectiveTx(t)
	for _, d := range sampleData[0:9] { 
		assert.NoError(t, store.CommitWithPvtData(d))
	}
	assert.False(t, store.IsPvtStoreAheadOfBlockStore())

	
	store.Shutdown()
	provider.Close()
	provider, err = NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	store, err = provider.Open("testLedger")
	assert.NoError(t, err)
	store.Init(btlPolicyForSampleData())

	
	info, err := store.GetBlockchainInfo()
	assert.NoError(t, err)
	assert.Equal(t, uint64(9), info.Height)
	pvtStoreHt, err := store.pvtdataStore.LastCommittedBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(9), pvtStoreHt)
	assert.False(t, store.IsPvtStoreAheadOfBlockStore())

	lastBlkAndPvtData := sampleData[9]
	
	
	validTxPvtData, validTxMissingPvtData := constructPvtDataAndMissingData(lastBlkAndPvtData)
	err = store.pvtdataStore.Commit(lastBlkAndPvtData.Block.Header.Number, validTxPvtData, validTxMissingPvtData)
	assert.NoError(t, err)

	
	store.Shutdown()
	provider.Close()
	provider, err = NewProvider(blockStoreDir, conf, metricsProvider)
	assert.NoError(t, err)
	store, err = provider.Open("testLedger")
	assert.NoError(t, err)
	store.Init(btlPolicyForSampleData())

	
	info, err = store.GetBlockchainInfo()
	assert.NoError(t, err)
	assert.Equal(t, uint64(9), info.Height)
	pvtStoreHt, err = store.pvtdataStore.LastCommittedBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(10), pvtStoreHt)
	assert.True(t, store.IsPvtStoreAheadOfBlockStore())

	
	assert.NoError(t, store.CommitWithPvtData(lastBlkAndPvtData))
	info, err = store.GetBlockchainInfo()
	assert.NoError(t, err)
	assert.Equal(t, uint64(10), info.Height)
	pvtStoreHt, err = store.pvtdataStore.LastCommittedBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(10), pvtStoreHt)
	assert.False(t, store.IsPvtStoreAheadOfBlockStore())
}

func TestConstructPvtdataMap(t *testing.T) {
	assert.Nil(t, constructPvtdataMap(nil))
}

func sampleDataWithPvtdataForSelectiveTx(t *testing.T) []*ledger.BlockAndPvtData {
	var blockAndpvtdata []*ledger.BlockAndPvtData
	blocks := testutil.ConstructTestBlocks(t, 10)
	for i := 0; i < 10; i++ {
		blockAndpvtdata = append(blockAndpvtdata, &ledger.BlockAndPvtData{Block: blocks[i]})
	}

	
	blockAndpvtdata[2].PvtData = samplePvtData(t, []uint64{3, 5, 6})
	txFilter := lutil.TxValidationFlags(blockAndpvtdata[2].Block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	txFilter.SetFlag(6, pb.TxValidationCode_INVALID_WRITESET)
	blockAndpvtdata[2].Block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txFilter

	
	blockAndpvtdata[3].PvtData = samplePvtData(t, []uint64{4, 6})

	
	missingData := make(ledger.TxMissingPvtDataMap)
	missingData.Add(4, "ns-4", "coll-4", true)
	missingData.Add(5, "ns-5", "coll-5", true)
	blockAndpvtdata[5].MissingPvtData = missingData
	txFilter = lutil.TxValidationFlags(blockAndpvtdata[5].Block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	txFilter.SetFlag(5, pb.TxValidationCode_INVALID_WRITESET)
	blockAndpvtdata[5].Block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txFilter

	return blockAndpvtdata
}

func sampleDataWithPvtdataForAllTxs(t *testing.T) []*ledger.BlockAndPvtData {
	var blockAndpvtdata []*ledger.BlockAndPvtData
	blocks := testutil.ConstructTestBlocks(t, 10)
	for i := 0; i < 10; i++ {
		blockAndpvtdata = append(blockAndpvtdata,
			&ledger.BlockAndPvtData{
				Block:   blocks[i],
				PvtData: samplePvtData(t, []uint64{uint64(i), uint64(i + 1)}),
			},
		)
	}
	return blockAndpvtdata
}

func samplePvtData(t *testing.T, txNums []uint64) map[uint64]*ledger.TxPvtData {
	pvtWriteSet := &rwset.TxPvtReadWriteSet{DataModel: rwset.TxReadWriteSet_KV}
	pvtWriteSet.NsPvtRwset = []*rwset.NsPvtReadWriteSet{
		{
			Namespace: "ns-1",
			CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
				{
					CollectionName: "coll-1",
					Rwset:          []byte("RandomBytes-PvtRWSet-ns1-coll1"),
				},
				{
					CollectionName: "coll-2",
					Rwset:          []byte("RandomBytes-PvtRWSet-ns1-coll2"),
				},
			},
		},
	}
	var pvtData []*ledger.TxPvtData
	for _, txNum := range txNums {
		pvtData = append(pvtData, &ledger.TxPvtData{SeqInBlock: txNum, WriteSet: pvtWriteSet})
	}
	return constructPvtdataMap(pvtData)
}

func btlPolicyForSampleData() pvtdatapolicy.BTLPolicy {
	return btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
		},
	)
}

func buildPrivateDataConfig(rootFSPath string) *pvtdatastorage.PrivateDataConfig {
	return &pvtdatastorage.PrivateDataConfig{
		PrivateDataConfig: &ledger.PrivateDataConfig{
			PurgeInterval: 1,
		},
		StorePath: filepath.Join(rootFSPath, "pvtdataStore"),
	}
}
