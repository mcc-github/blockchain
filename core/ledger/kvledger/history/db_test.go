/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package history

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/ledger/queryresult"
	"github.com/mcc-github/blockchain-protos-go/peer"
	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	util2 "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flogging.ActivateSpec("leveldbhelper,historyleveldb=debug")
	os.Exit(m.Run())
}


func TestSavepoint(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()

	
	savepoint, err := env.testHistoryDB.GetLastSavepoint()
	assert.NoError(t, err, "Error upon historyDatabase.GetLastSavepoint()")
	assert.Nil(t, savepoint)

	
	status, blockNum, err := env.testHistoryDB.ShouldRecover(0)
	assert.NoError(t, err, "Error upon historyDatabase.ShouldRecover()")
	assert.True(t, status)
	assert.Equal(t, uint64(0), blockNum)

	bg, gb := testutil.NewBlockGenerator(t, "testLedger", false)
	assert.NoError(t, env.testHistoryDB.Commit(gb))
	
	savepoint, err = env.testHistoryDB.GetLastSavepoint()
	assert.NoError(t, err, "Error upon historyDatabase.GetLastSavepoint()")
	assert.Equal(t, uint64(0), savepoint.BlockNum)

	
	txid := util2.GenerateUUID()
	simulator, _ := env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value1"))
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimResBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimResBytes})
	assert.NoError(t, env.testHistoryDB.Commit(block1))
	savepoint, err = env.testHistoryDB.GetLastSavepoint()
	assert.NoError(t, err, "Error upon historyDatabase.GetLastSavepoint()")
	assert.Equal(t, uint64(1), savepoint.BlockNum)

	
	status, blockNum, err = env.testHistoryDB.ShouldRecover(1)
	assert.NoError(t, err, "Error upon historyDatabase.ShouldRecover()")
	assert.False(t, status)
	assert.Equal(t, uint64(2), blockNum)

	
	txid = util2.GenerateUUID()
	simulator, _ = env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value2"))
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimResBytes, _ = simRes.GetPubSimulationBytes()
	block2 := bg.NextBlock([][]byte{pubSimResBytes})

	
	env.testHistoryDB.CommitLostBlock(&ledger.BlockAndPvtData{Block: block2})
	savepoint, err = env.testHistoryDB.GetLastSavepoint()
	assert.NoError(t, err, "Error upon historyDatabase.GetLastSavepoint()")
	assert.Equal(t, uint64(2), savepoint.BlockNum)

	
	status, blockNum, err = env.testHistoryDB.ShouldRecover(10)
	assert.NoError(t, err, "Error upon historyDatabase.ShouldRecover()")
	assert.True(t, status)
	assert.Equal(t, uint64(3), blockNum)
}

func TestHistory(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	provider := env.testBlockStorageEnv.provider
	ledger1id := "ledger1"
	store1, err := provider.OpenBlockStore(ledger1id)
	assert.NoError(t, err, "Error upon provider.OpenBlockStore()")
	defer store1.Shutdown()
	assert.Equal(t, "history", env.testHistoryDB.Name())

	bg, gb := testutil.NewBlockGenerator(t, ledger1id, false)
	assert.NoError(t, store1.AddBlock(gb))
	assert.NoError(t, env.testHistoryDB.Commit(gb))

	
	txid := util2.GenerateUUID()
	simulator, _ := env.txmgr.NewTxSimulator(txid)
	value1 := []byte("value1")
	simulator.SetState("ns1", "key7", value1)
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimResBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimResBytes})
	err = store1.AddBlock(block1)
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block1)
	assert.NoError(t, err)

	
	simulationResults := [][]byte{}
	txid = util2.GenerateUUID()
	simulator, _ = env.txmgr.NewTxSimulator(txid)
	value2 := []byte("value2")
	simulator.SetState("ns1", "key7", value2)
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimResBytes, _ = simRes.GetPubSimulationBytes()
	simulationResults = append(simulationResults, pubSimResBytes)
	
	txid2 := util2.GenerateUUID()
	simulator2, _ := env.txmgr.NewTxSimulator(txid2)
	value3 := []byte("value3")
	simulator2.SetState("ns1", "key7", value3)
	simulator2.Done()
	simRes2, _ := simulator2.GetTxSimulationResults()
	pubSimResBytes2, _ := simRes2.GetPubSimulationBytes()
	simulationResults = append(simulationResults, pubSimResBytes2)
	block2 := bg.NextBlock(simulationResults)
	err = store1.AddBlock(block2)
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block2)
	assert.NoError(t, err)

	
	txid = util2.GenerateUUID()
	simulator, _ = env.txmgr.NewTxSimulator(txid)
	simulator.DeleteState("ns1", "key7")
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimResBytes, _ = simRes.GetPubSimulationBytes()
	block3 := bg.NextBlock([][]byte{pubSimResBytes})
	err = store1.AddBlock(block3)
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block3)
	assert.NoError(t, err)
	t.Logf("Inserted all 3 blocks")

	qhistory, err := env.testHistoryDB.NewQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewQueryExecutor")

	itr, err2 := qhistory.GetHistoryForKey("ns1", "key7")
	assert.NoError(t, err2, "Error upon GetHistoryForKey()")

	count := 0
	for {
		
		kmod, _ := itr.Next()
		if kmod == nil {
			break
		}
		txid = kmod.(*queryresult.KeyModification).TxId
		retrievedValue := kmod.(*queryresult.KeyModification).Value
		retrievedTimestamp := kmod.(*queryresult.KeyModification).Timestamp
		retrievedIsDelete := kmod.(*queryresult.KeyModification).IsDelete
		t.Logf("Retrieved history record for key=key7 at TxId=%s with value %v and timestamp %v",
			txid, retrievedValue, retrievedTimestamp)
		count++
		if count != 1 {
			
			expectedValue := []byte("value" + strconv.Itoa(5-count))
			assert.Equal(t, expectedValue, retrievedValue)
			assert.NotNil(t, retrievedTimestamp)
			assert.False(t, retrievedIsDelete)
		} else {
			
			assert.Equal(t, []uint8(nil), retrievedValue)
			assert.NotNil(t, retrievedTimestamp)
			assert.True(t, retrievedIsDelete)
		}
	}
	assert.Equal(t, 4, count)
}

func TestHistoryForInvalidTran(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	provider := env.testBlockStorageEnv.provider
	ledger1id := "ledger1"
	store1, err := provider.OpenBlockStore(ledger1id)
	assert.NoError(t, err, "Error upon provider.OpenBlockStore()")
	defer store1.Shutdown()

	bg, gb := testutil.NewBlockGenerator(t, ledger1id, false)
	assert.NoError(t, store1.AddBlock(gb))
	assert.NoError(t, env.testHistoryDB.Commit(gb))

	
	txid := util2.GenerateUUID()
	simulator, _ := env.txmgr.NewTxSimulator(txid)
	value1 := []byte("value1")
	simulator.SetState("ns1", "key7", value1)
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimResBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimResBytes})

	
	txsFilter := util.TxValidationFlags(block1.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	txsFilter.SetFlag(0, peer.TxValidationCode_INVALID_OTHER_REASON)
	block1.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter

	err = store1.AddBlock(block1)
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block1)
	assert.NoError(t, err)

	qhistory, err := env.testHistoryDB.NewQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewQueryExecutor")

	itr, err2 := qhistory.GetHistoryForKey("ns1", "key7")
	assert.NoError(t, err2, "Error upon GetHistoryForKey()")

	
	kmod, _ := itr.Next()
	assert.Nil(t, kmod)
}



func TestGenesisBlockNoError(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	block, err := configtxtest.MakeGenesisBlock("test_chainid")
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block)
	assert.NoError(t, err)
}



func TestHistoryWithKeyContainingNilBytes(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	provider := env.testBlockStorageEnv.provider
	ledger1id := "ledger1"
	store1, err := provider.OpenBlockStore(ledger1id)
	assert.NoError(t, err, "Error upon provider.OpenBlockStore()")
	defer store1.Shutdown()

	bg, gb := testutil.NewBlockGenerator(t, ledger1id, false)
	assert.NoError(t, store1.AddBlock(gb))
	assert.NoError(t, env.testHistoryDB.Commit(gb))

	
	txid := util2.GenerateUUID()
	simulator, _ := env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", "key", []byte("value1")) 
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimResBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimResBytes})
	err = store1.AddBlock(block1)
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block1)
	assert.NoError(t, err)

	
	simulationResults := [][]byte{}
	txid = util2.GenerateUUID()
	simulator, _ = env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", "key", []byte("value2")) 
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimResBytes, _ = simRes.GetPubSimulationBytes()
	simulationResults = append(simulationResults, pubSimResBytes)

	
	txid2 := util2.GenerateUUID()
	simulator2, _ := env.txmgr.NewTxSimulator(txid2)

	
	key1 := "\x00key\x00\x01\x01\x15"
	simulator2.SetState("ns1", key1, []byte("dummyVal1"))

	
	
	
	key2 := "key\x00\x01\x01\x15" 
	simulator2.SetState("ns1", key2, []byte("dummyVal2"))

	
	
	key3 := "key\x00\x0c" 
	simulator2.SetState("ns1", key3, []byte("dummyVal3"))

	
	
	key4 := "key\x00\x01"
	simulator2.SetState("ns1", key4, []byte("dummyVal4"))

	
	
	key5 := "key\x00\x04\x01"
	simulator2.SetState("ns1", key5, []byte("dummyVal5"))

	
	simulator2.Done()
	simRes2, _ := simulator2.GetTxSimulationResults()
	pubSimResBytes2, _ := simRes2.GetPubSimulationBytes()
	simulationResults = append(simulationResults, pubSimResBytes2)
	block2 := bg.NextBlock(simulationResults)
	err = store1.AddBlock(block2)
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block2)
	assert.NoError(t, err)

	qhistory, err := env.testHistoryDB.NewQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewQueryExecutor")

	
	testutilVerifyResults(t, qhistory, "ns1", "key", []string{"value2", "value1"})
	testutilVerifyResults(t, qhistory, "ns1", key1, []string{"dummyVal1"})
	testutilVerifyResults(t, qhistory, "ns1", key2, []string{"dummyVal2"})
	testutilVerifyResults(t, qhistory, "ns1", key3, []string{"dummyVal3"})
	testutilVerifyResults(t, qhistory, "ns1", key4, []string{"dummyVal4"})
	testutilVerifyResults(t, qhistory, "ns1", key5, []string{"dummyVal5"})

	
	testutilCheckKeyNotInRange(t, qhistory, "ns1", "key", key1)
	testutilCheckKeyNotInRange(t, qhistory, "ns1", "key", key2)
	testutilCheckKeyNotInRange(t, qhistory, "ns1", "key", key3)
	testutilCheckKeyNotInRange(t, qhistory, "ns1", "key", key4)
	testutilCheckKeyNotInRange(t, qhistory, "ns1", "key", key5)
}




func TestHistoryWithBlockNumber256(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	provider := env.testBlockStorageEnv.provider
	ledger1id := "ledger1"
	store1, err := provider.OpenBlockStore(ledger1id)
	assert.NoError(t, err, "Error upon provider.OpenBlockStore()")
	defer store1.Shutdown()

	bg, gb := testutil.NewBlockGenerator(t, ledger1id, false)
	assert.NoError(t, store1.AddBlock(gb))
	assert.NoError(t, env.testHistoryDB.Commit(gb))

	
	for i := 1; i <= 256; i++ {
		txid := util2.GenerateUUID()
		simulator, _ := env.txmgr.NewTxSimulator(txid)
		value := fmt.Sprintf("value%d", i)
		simulator.SetState("ns1", "key", []byte(value))
		simulator.Done()
		simRes, _ := simulator.GetTxSimulationResults()
		pubSimResBytes, _ := simRes.GetPubSimulationBytes()
		block := bg.NextBlock([][]byte{pubSimResBytes})
		err = store1.AddBlock(block)
		assert.NoError(t, err)
		err = env.testHistoryDB.Commit(block)
		assert.NoError(t, err)
	}

	
	qhistory, err := env.testHistoryDB.NewQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewQueryExecutor")

	
	expectedHistoryResults := make([]string, 0)
	for i := 256; i >= 1; i-- {
		expectedHistoryResults = append(expectedHistoryResults, fmt.Sprintf("value%d", i))
	}
	testutilVerifyResults(t, qhistory, "ns1", "key", expectedHistoryResults)
}

func TestName(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	assert.Equal(t, "history", env.testHistoryDB.Name())
}


func testutilVerifyResults(t *testing.T, hqe ledger.HistoryQueryExecutor, ns, key string, expectedVals []string) {
	itr, err := hqe.GetHistoryForKey(ns, key)
	assert.NoError(t, err, "Error upon GetHistoryForKey()")
	retrievedVals := []string{}
	for {
		kmod, _ := itr.Next()
		if kmod == nil {
			break
		}
		txid := kmod.(*queryresult.KeyModification).TxId
		retrievedValue := string(kmod.(*queryresult.KeyModification).Value)
		retrievedVals = append(retrievedVals, retrievedValue)
		t.Logf("Retrieved history record at TxId=%s with value %s", txid, retrievedValue)
	}
	assert.Equal(t, expectedVals, retrievedVals)
}


func testutilCheckKeyNotInRange(t *testing.T, hqe ledger.HistoryQueryExecutor, ns, desiredKey, falseKey string) {
	itr, err := hqe.GetHistoryForKey(ns, desiredKey)
	assert.NoError(t, err, "Error upon GetHistoryForKey()")
	scanner := itr.(*historyScanner)
	rangeScanKeys := constructRangeScan(ns, falseKey)
	for {
		if !scanner.dbItr.Next() {
			break
		}
		historyKey := scanner.dbItr.Key()
		if bytes.Contains(historyKey, rangeScanKeys.startKey) {
			assert.Failf(t, "false key %s should not be returned in range query for key %s", falseKey, desiredKey)
		}
	}
}
