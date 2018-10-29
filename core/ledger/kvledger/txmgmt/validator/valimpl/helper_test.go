/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valimpl

import (
	"fmt"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/flogging/floggingtest"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/internal"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	lutils "github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	putils "github.com/mcc-github/blockchain/protos/utils"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flogging.ActivateSpec("internal=debug")
	viper.Set("peer.fileSystemPath", "/tmp/blockchain/ledgertests/kvledger/txmgmt/validator/internal")
	os.Exit(m.Run())
}

func TestValidateAndPreparePvtBatch(t *testing.T) {
	testDBEnv := &privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	testDB := testDBEnv.GetDBHandle("emptydb")

	pubSimulationResults := [][]byte{}
	pvtDataMap := make(map[uint64]*ledger.TxPvtData)

	txids := []string{"tx1", "tx2", "tx3"}

	
	
	

	
	
	tx1SimulationResults := testutilSampleTxSimulationResults(t, "key1")
	res, err := tx1SimulationResults.GetPubSimulationBytes()
	assert.NoError(t, err)

	
	pubSimulationResults = append(pubSimulationResults, res)

	
	tx1PvtData := &ledger.TxPvtData{SeqInBlock: 0, WriteSet: tx1SimulationResults.PvtSimulationResults}
	pvtDataMap[uint64(0)] = tx1PvtData

	
	
	tx2SimulationResults := testutilSampleTxSimulationResults(t, "key2")
	res, err = tx2SimulationResults.GetPubSimulationBytes()
	assert.NoError(t, err)

	
	pubSimulationResults = append(pubSimulationResults, res)

	
	

	
	
	tx3SimulationResults := testutilSampleTxSimulationResults(t, "key3")
	res, err = tx3SimulationResults.GetPubSimulationBytes()
	assert.NoError(t, err)

	
	pubSimulationResults = append(pubSimulationResults, res)

	
	tx3PvtData := &ledger.TxPvtData{SeqInBlock: 2, WriteSet: tx3SimulationResults.PvtSimulationResults}
	pvtDataMap[uint64(2)] = tx3PvtData

	
	block := testutil.ConstructBlockWithTxid(t, 10, testutil.ConstructRandomBytes(t, 32), pubSimulationResults, txids, false)

	
	expectedPerProcessedBlock := &internal.Block{Num: 10}
	tx1TxRWSet, err := rwsetutil.TxRwSetFromProtoMsg(tx1SimulationResults.PubSimulationResults)
	assert.NoError(t, err)
	expectedPerProcessedBlock.Txs = append(expectedPerProcessedBlock.Txs, &internal.Transaction{IndexInBlock: 0, ID: "tx1", RWSet: tx1TxRWSet})

	tx2TxRWSet, err := rwsetutil.TxRwSetFromProtoMsg(tx2SimulationResults.PubSimulationResults)
	assert.NoError(t, err)
	expectedPerProcessedBlock.Txs = append(expectedPerProcessedBlock.Txs, &internal.Transaction{IndexInBlock: 1, ID: "tx2", RWSet: tx2TxRWSet})

	tx3TxRWSet, err := rwsetutil.TxRwSetFromProtoMsg(tx3SimulationResults.PubSimulationResults)
	assert.NoError(t, err)
	expectedPerProcessedBlock.Txs = append(expectedPerProcessedBlock.Txs, &internal.Transaction{IndexInBlock: 2, ID: "tx3", RWSet: tx3TxRWSet})
	alwaysValidKVFunc := func(key string, value []byte) error {
		return nil
	}
	actualPreProcessedBlock, err := preprocessProtoBlock(nil, alwaysValidKVFunc, block, false)
	assert.NoError(t, err)
	assert.Equal(t, expectedPerProcessedBlock, actualPreProcessedBlock)

	
	
	

	
	mvccValidatedBlock := actualPreProcessedBlock
	mvccValidatedBlock.Txs[0].ValidationCode = peer.TxValidationCode_VALID
	mvccValidatedBlock.Txs[1].ValidationCode = peer.TxValidationCode_VALID
	mvccValidatedBlock.Txs[2].ValidationCode = peer.TxValidationCode_INVALID_OTHER_REASON

	
	expectedPvtUpdates := privacyenabledstate.NewPvtUpdateBatch()
	tx1TxPvtRWSet, err := rwsetutil.TxPvtRwSetFromProtoMsg(tx1SimulationResults.PvtSimulationResults)
	assert.NoError(t, err)
	addPvtRWSetToPvtUpdateBatch(tx1TxPvtRWSet, expectedPvtUpdates, version.NewHeight(uint64(10), uint64(0)))

	actualPvtUpdates, err := validateAndPreparePvtBatch(mvccValidatedBlock, testDB, nil, pvtDataMap)
	assert.NoError(t, err)
	assert.Equal(t, expectedPvtUpdates, actualPvtUpdates)

	expectedtxsFilter := []uint8{uint8(peer.TxValidationCode_VALID), uint8(peer.TxValidationCode_VALID), uint8(peer.TxValidationCode_INVALID_OTHER_REASON)}

	postprocessProtoBlock(block, mvccValidatedBlock)
	assert.Equal(t, expectedtxsFilter, block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
}

func TestPreprocessProtoBlock(t *testing.T) {
	allwaysValidKVfunc := func(key string, value []byte) error {
		return nil
	}
	
	
	gb := testutil.ConstructTestBlock(t, 10, 1, 1)
	_, err := preprocessProtoBlock(nil, allwaysValidKVfunc, gb, false)
	assert.NoError(t, err)
	
	gb = testutil.ConstructTestBlock(t, 11, 1, 1)
	gb.Data = &common.BlockData{Data: [][]byte{{123}}}
	gb.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] =
		lutils.NewTxValidationFlagsSetValue(len(gb.Data.Data), peer.TxValidationCode_VALID)
	_, err = preprocessProtoBlock(nil, allwaysValidKVfunc, gb, false)
	assert.Error(t, err)
	t.Log(err)
	
	gb = testutil.ConstructTestBlock(t, 12, 1, 1)
	envBytes, _ := putils.GetBytesEnvelope(&common.Envelope{Payload: []byte{123}})
	gb.Data = &common.BlockData{Data: [][]byte{envBytes}}
	_, err = preprocessProtoBlock(nil, allwaysValidKVfunc, gb, false)
	assert.Error(t, err)
	t.Log(err)
	
	gb = testutil.ConstructTestBlock(t, 13, 1, 1)
	payloadBytes, _ := putils.GetBytesPayload(&common.Payload{
		Header: &common.Header{ChannelHeader: []byte{123}},
	})
	envBytes, _ = putils.GetBytesEnvelope(&common.Envelope{Payload: payloadBytes})
	gb.Data = &common.BlockData{Data: [][]byte{envBytes}}
	_, err = preprocessProtoBlock(nil, allwaysValidKVfunc, gb, false)
	assert.Error(t, err)
	t.Log(err)

	
	gb = testutil.ConstructTestBlock(t, 14, 1, 1)
	payloadBytes, _ = putils.GetBytesPayload(&common.Payload{
		Header: &common.Header{ChannelHeader: []byte{123}},
	})
	envBytes, _ = putils.GetBytesEnvelope(&common.Envelope{Payload: payloadBytes})
	gb.Data = &common.BlockData{Data: [][]byte{envBytes}}
	flags := lutils.NewTxValidationFlags(len(gb.Data.Data))
	flags.SetFlag(0, peer.TxValidationCode_BAD_CHANNEL_HEADER)
	gb.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = flags
	_, err = preprocessProtoBlock(nil, allwaysValidKVfunc, gb, false)
	assert.NoError(t, err) 

	
	var blockNum uint64 = 15
	txid := "testtxid1234"
	gb = testutil.ConstructBlockWithTxid(t, blockNum, []byte{123},
		[][]byte{{123}}, []string{txid}, false)
	flags = lutils.NewTxValidationFlags(len(gb.Data.Data))
	flags.SetFlag(0, peer.TxValidationCode_BAD_HEADER_EXTENSION)
	gb.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = flags

	
	oldLogger := logger
	defer func() { logger = oldLogger }()
	l, recorder := floggingtest.NewTestLogger(t)
	logger = l

	_, err = preprocessProtoBlock(nil, allwaysValidKVfunc, gb, false)
	assert.NoError(t, err)
	expected := fmt.Sprintf(
		"Channel [%s]: Block [%d] Transaction index [%d] TxId [%s] marked as invalid by committer. Reason code [%s]",
		util.GetTestChainID(), blockNum, 0, txid, peer.TxValidationCode_BAD_HEADER_EXTENSION,
	)
	assert.NotEmpty(t, recorder.MessagesContaining(expected))
}

func TestPreprocessProtoBlockInvalidWriteset(t *testing.T) {
	kvValidationFunc := func(key string, value []byte) error {
		if value[0] == '_' {
			return fmt.Errorf("value [%s] found to be invalid by 'kvValidationFunc for testing'", value)
		}
		return nil
	}

	rwSetBuilder := rwsetutil.NewRWSetBuilder()
	rwSetBuilder.AddToWriteSet("ns", "key", []byte("_invalidValue")) 
	simulation1, err := rwSetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	simulation1Bytes, err := simulation1.GetPubSimulationBytes()
	assert.NoError(t, err)

	rwSetBuilder = rwsetutil.NewRWSetBuilder()
	rwSetBuilder.AddToWriteSet("ns", "key", []byte("validValue")) 
	simulation2, err := rwSetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	simulation2Bytes, err := simulation2.GetPubSimulationBytes()
	assert.NoError(t, err)

	block := testutil.ConstructBlock(t, 1, testutil.ConstructRandomBytes(t, 32),
		[][]byte{simulation1Bytes, simulation2Bytes}, false) 
	txfilter := lutils.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assert.True(t, txfilter.IsValid(0))
	assert.True(t, txfilter.IsValid(1)) 

	internalBlock, err := preprocessProtoBlock(nil, kvValidationFunc, block, false)
	assert.NoError(t, err)
	assert.False(t, txfilter.IsValid(0)) 
	assert.True(t, txfilter.IsValid(1))  
	assert.Len(t, internalBlock.Txs, 1)
	assert.Equal(t, internalBlock.Txs[0].IndexInBlock, 1)
}

func TestIncrementPvtdataVersionIfNeeded(t *testing.T) {
	testDBEnv := &privacyenabledstate.LevelDBCommonStorageTestEnv{}
	testDBEnv.Init(t)
	defer testDBEnv.Cleanup()
	testDB := testDBEnv.GetDBHandle("testdb")
	updateBatch := privacyenabledstate.NewUpdateBatch()
	
	updateBatch.PvtUpdates.Put("ns", "coll1", "key1", []byte("value1"), version.NewHeight(1, 1))
	updateBatch.PvtUpdates.Put("ns", "coll2", "key2", []byte("value2"), version.NewHeight(1, 2))
	updateBatch.PvtUpdates.Put("ns", "coll3", "key3", []byte("value3"), version.NewHeight(1, 3))
	updateBatch.PvtUpdates.Put("ns", "col4", "key4", []byte("value4"), version.NewHeight(1, 4))
	testDB.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(1, 4))

	
	hashUpdates := privacyenabledstate.NewHashedUpdateBatch()
	hashUpdates.PutValHashAndMetadata("ns", "coll1", lutils.ComputeStringHash("key1"),
		lutils.ComputeStringHash("value1_set_by_tx1"), []byte("metadata1_set_by_tx2"), version.NewHeight(2, 2)) 
	hashUpdates.PutValHashAndMetadata("ns", "coll2", lutils.ComputeStringHash("key2"),
		lutils.ComputeStringHash("value2"), []byte("metadata2_set_by_tx4"), version.NewHeight(2, 4)) 
	hashUpdates.PutValHashAndMetadata("ns", "coll3", lutils.ComputeStringHash("key3"),
		lutils.ComputeStringHash("value3_set_by_tx6"), []byte("metadata3"), version.NewHeight(2, 6)) 
	pubAndHashedUpdatesBatch := &internal.PubAndHashUpdates{HashUpdates: hashUpdates}

	
	pvtUpdateBatch := privacyenabledstate.NewPvtUpdateBatch()
	pvtUpdateBatch.Put("ns", "coll1", "key1", []byte("value1_set_by_tx1"), version.NewHeight(2, 1))
	pvtUpdateBatch.Put("ns", "coll3", "key3", []byte("value3_set_by_tx5"), version.NewHeight(2, 5))
	
	metadataUpdates := metadataUpdates{collKey{"ns", "coll1", "key1"}: true, collKey{"ns", "coll2", "key2"}: true}

	
	err := incrementPvtdataVersionIfNeeded(metadataUpdates, pvtUpdateBatch, pubAndHashedUpdatesBatch, testDB)
	assert.NoError(t, err)

	assert.Equal(t,
		&statedb.VersionedValue{Value: []byte("value1_set_by_tx1"), Version: version.NewHeight(2, 2)}, 
		pvtUpdateBatch.Get("ns", "coll1", "key1"),
	)

	assert.Equal(t,
		&statedb.VersionedValue{Value: []byte("value2"), Version: version.NewHeight(2, 4)}, 
		pvtUpdateBatch.Get("ns", "coll2", "key2"),
	)

	assert.Equal(t,
		&statedb.VersionedValue{Value: []byte("value3_set_by_tx5"), Version: version.NewHeight(2, 5)}, 
		pvtUpdateBatch.Get("ns", "coll3", "key3"),
	)
}


func memoryRecordN(b *logging.MemoryBackend, n int) *logging.Record {
	node := b.Head()
	for i := 0; i < n; i++ {
		if node == nil {
			break
		}
		node = node.Next()
	}
	if node == nil {
		return nil
	}
	return node.Record
}

func testutilSampleTxSimulationResults(t *testing.T, key string) *ledger.TxSimulationResults {
	rwSetBuilder := rwsetutil.NewRWSetBuilder()
	
	rwSetBuilder.AddToReadSet("ns1", key, version.NewHeight(1, 1))
	rwSetBuilder.AddToReadSet("ns2", key, version.NewHeight(1, 1))
	rwSetBuilder.AddToWriteSet("ns2", key, []byte("ns2-key1-value"))

	
	rwSetBuilder.AddToHashedReadSet("ns1", "coll1", key, version.NewHeight(1, 1))
	rwSetBuilder.AddToHashedReadSet("ns1", "coll2", key, version.NewHeight(1, 1))
	rwSetBuilder.AddToPvtAndHashedWriteSet("ns1", "coll2", key, []byte("pvt-ns1-coll2-key1-value"))

	
	rwSetBuilder.AddToHashedReadSet("ns2", "coll1", key, version.NewHeight(1, 1))
	rwSetBuilder.AddToHashedReadSet("ns2", "coll2", key, version.NewHeight(1, 1))
	rwSetBuilder.AddToPvtAndHashedWriteSet("ns2", "coll2", key, []byte("pvt-ns2-coll2-key1-value"))
	rwSetBuilder.AddToPvtAndHashedWriteSet("ns2", "coll3", key, nil)

	rwSetBuilder.AddToHashedReadSet("ns3", "coll1", key, version.NewHeight(1, 1))

	pubAndPvtSimulationResults, err := rwSetBuilder.GetTxSimulationResults()
	if err != nil {
		t.Fatalf("ConstructSimulationResultsWithPvtData failed while getting simulation results, err %s", err)
	}

	return pubAndPvtSimulationResults
}
