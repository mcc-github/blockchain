/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valimpl

import (
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/valinternal"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	lgrutil "github.com/mcc-github/blockchain/core/ledger/util"
	lutils "github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	putils "github.com/mcc-github/blockchain/protos/utils"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)

func getPubAndPvtSimulationResults(t *testing.T, key string) *ledger.TxSimulationResults {
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

func TestValidateAndPreparePvtBatch(t *testing.T) {
	pubSimulationResults := [][]byte{}
	pvtDataMap := make(map[uint64]*ledger.TxPvtData)

	txids := []string{"tx1", "tx2", "tx3"}

	
	
	

	
	
	tx1SimulationResults := getPubAndPvtSimulationResults(t, "key1")
	res, err := tx1SimulationResults.GetPubSimulationBytes()
	assert.NoError(t, err)

	
	pubSimulationResults = append(pubSimulationResults, res)

	
	tx1PvtData := &ledger.TxPvtData{0, tx1SimulationResults.PvtSimulationResults}
	pvtDataMap[uint64(0)] = tx1PvtData

	
	
	tx2SimulationResults := getPubAndPvtSimulationResults(t, "key2")
	res, err = tx2SimulationResults.GetPubSimulationBytes()
	assert.NoError(t, err)

	
	pubSimulationResults = append(pubSimulationResults, res)

	
	

	
	
	tx3SimulationResults := getPubAndPvtSimulationResults(t, "key3")
	res, err = tx3SimulationResults.GetPubSimulationBytes()
	assert.NoError(t, err)

	
	pubSimulationResults = append(pubSimulationResults, res)

	
	tx3PvtData := &ledger.TxPvtData{2, tx3SimulationResults.PvtSimulationResults}
	pvtDataMap[uint64(2)] = tx3PvtData

	
	block := testutil.ConstructBlockWithTxid(t, 10, testutil.ConstructRandomBytes(t, 32), pubSimulationResults, txids, false)

	
	expectedPerProcessedBlock := &valinternal.Block{Num: 10}
	tx1TxRWSet, err := rwsetutil.TxRwSetFromProtoMsg(tx1SimulationResults.PubSimulationResults)
	assert.NoError(t, err)
	expectedPerProcessedBlock.Txs = append(expectedPerProcessedBlock.Txs, &valinternal.Transaction{IndexInBlock: 0, ID: "tx1", RWSet: tx1TxRWSet})

	tx2TxRWSet, err := rwsetutil.TxRwSetFromProtoMsg(tx2SimulationResults.PubSimulationResults)
	assert.NoError(t, err)
	expectedPerProcessedBlock.Txs = append(expectedPerProcessedBlock.Txs, &valinternal.Transaction{IndexInBlock: 1, ID: "tx2", RWSet: tx2TxRWSet})

	tx3TxRWSet, err := rwsetutil.TxRwSetFromProtoMsg(tx3SimulationResults.PubSimulationResults)
	assert.NoError(t, err)
	expectedPerProcessedBlock.Txs = append(expectedPerProcessedBlock.Txs, &valinternal.Transaction{IndexInBlock: 2, ID: "tx3", RWSet: tx3TxRWSet})
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

	actualPvtUpdates, err := validateAndPreparePvtBatch(mvccValidatedBlock, pvtDataMap)
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
	
	backend := logging.NewMemoryBackend(1)
	logging.SetBackend(backend)
	_, err = preprocessProtoBlock(nil, allwaysValidKVfunc, gb, false)
	assert.NoError(t, err)
	expected := fmt.Sprintf("Channel [%s]: Block [%d] Transaction index [%d] TxId [%s]"+
		" marked as invalid by committer. Reason code [%s]",
		util.GetTestChainID(), blockNum, 0, txid, peer.TxValidationCode_BAD_HEADER_EXTENSION.String())
	t.Log(expected)
	assert.Equal(t, expected, memoryRecordN(backend, 0).Message())
	
	t.Log(memoryRecordN(backend, 0).Message())

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
	txfilter := lgrutil.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assert.True(t, txfilter.IsValid(0))
	assert.True(t, txfilter.IsValid(1)) 

	internalBlock, err := preprocessProtoBlock(nil, kvValidationFunc, block, false)
	assert.NoError(t, err)
	assert.False(t, txfilter.IsValid(0)) 
	assert.True(t, txfilter.IsValid(1))  
	assert.Len(t, internalBlock.Txs, 1)
	assert.Equal(t, internalBlock.Txs[0].IndexInBlock, 1)
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
