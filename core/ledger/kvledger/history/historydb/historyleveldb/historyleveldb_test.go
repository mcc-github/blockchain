/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package historyleveldb

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"testing"

	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	util2 "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/history/historydb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/history/historydb/historyleveldb/fakes"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/mcc-github/blockchain/protos/peer"
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

	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewHistoryQueryExecutor")

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
		if count != 4 {
			expectedValue := []byte("value" + strconv.Itoa(count))
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

	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewHistoryQueryExecutor")

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

	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewHistoryQueryExecutor")
	testutilVerifyResults(t, qhistory, "ns1", "key", []string{"value1", "value2"})
	testutilVerifyResults(t, qhistory, "ns1", key1, []string{"dummyVal1"})
	testutilVerifyResults(t, qhistory, "ns1", key2, []string{"dummyVal2"})
	testutilVerifyResults(t, qhistory, "ns1", key3, []string{"dummyVal3"})
	testutilVerifyResults(t, qhistory, "ns1", key4, []string{"dummyVal4"})
	testutilVerifyResults(t, qhistory, "ns1", key5, []string{"dummyVal5"})

	
	testutilCheckKeyInRange(t, qhistory, "ns1", "key", key1, 0)
	testutilCheckKeyInRange(t, qhistory, "ns1", "key", key2, 1)
	testutilCheckKeyInRange(t, qhistory, "ns1", "key", key3, 1)
	testutilCheckKeyInRange(t, qhistory, "ns1", "key", key4, 1)
	testutilCheckKeyInRange(t, qhistory, "ns1", "key", key5, 1)

	
	origLogger := logger
	fakeLogger := &fakes.HistoryleveldbLogger{HistorydbLogger: &fakes.HistorydbLogger{}}
	logger = fakeLogger
	defer func() {
		logger = origLogger
	}()
	exhaustQueryResultsForKey(t, qhistory, "ns1", "key")
	assert.Equalf(t, 4, fakeLogger.WarnfCallCount(), "Get expected number of warning messages")
	testutilVerifyWarningMsg(t, fakeLogger, "ns1", key2,
		"Some other key [%#v] found in the range while scanning history for key [%#v]. Skipping (decoding error: %s)",
		"decoded size from DecodeVarint is invalid, expected <=8, but got 21")
	testutilVerifyWarningMsg(t, fakeLogger, "ns1", key3,
		"Some other key [%#v] found in the range while scanning history for key [%#v]. Skipping (decoding error: %s)",
		"decoded size from DecodeVarint is invalid, expected <=8, but got 12")
	testutilVerifyWarningMsg(t, fakeLogger, "ns1", key4,
		"Some other key [%#v] found in the range while scanning history for key [%#v]. Skipping (decoding error: %s)",
		"number of decoded bytes (4) is not equal to the length of blockNumTranNumBytes (6)")
	testutilVerifyWarningMsg(t, fakeLogger, "ns1", key5,
		"Some other clashing key [%#v] found in the range while scanning history for key [%#v]. Skipping (cannot find block:tx)",
		"")
}



func TestNoKeyFoundForFalseKey(t *testing.T) {
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

	
	
	
	
	
	otherKey := "key\x00\x04\x00"
	txid := util2.GenerateUUID()
	simulator, _ := env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", otherKey, []byte("otherValue"))
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimResBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimResBytes})
	err = store1.AddBlock(block1)
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block1)
	assert.NoError(t, err)

	
	for i := 2; i <= 256; i++ {
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

	
	txid = util2.GenerateUUID()
	simulator, _ = env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", "key2", []byte("key2Value"))
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimResBytes, _ = simRes.GetPubSimulationBytes()
	block257 := bg.NextBlock([][]byte{pubSimResBytes})
	err = store1.AddBlock(block257)
	assert.NoError(t, err)
	err = env.testHistoryDB.Commit(block257)
	assert.NoError(t, err)

	
	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewHistoryQueryExecutor")

	
	origLogger := logger
	fakeLogger := &fakes.HistoryleveldbLogger{HistorydbLogger: &fakes.HistorydbLogger{}}
	logger = fakeLogger
	defer func() {
		logger = origLogger
	}()
	exhaustQueryResultsForKey(t, qhistory, "ns1", "key")
	assert.Equalf(t, 1, fakeLogger.WarnfCallCount(), "Get expected number of warning messages")
	msg, _ := fakeLogger.WarnfArgsForCall(0)
	
	
	
	assert.Equal(t, "Some other key [%#v] found in the range while scanning history for key [%#v]. Skipping (namespace or key not found)", msg)
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

	
	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(store1)
	assert.NoError(t, err, "Error upon NewHistoryQueryExecutor")
	itr, err := qhistory.GetHistoryForKey("ns1", "key")
	assert.NoError(t, err, "Error upon GetHistoryForKey()")

	
	numEntries := 0
	valueInBlock256 := "unknown"
	for {
		kmod, _ := itr.Next()
		if kmod == nil {
			break
		}
		numEntries++
		retrievedValue := string(kmod.(*queryresult.KeyModification).Value)
		if numEntries == 256 {
			valueInBlock256 = retrievedValue
		}
	}
	assert.Equal(t, 256, numEntries)
	assert.Equal(t, "value256", valueInBlock256)
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


func testutilCheckKeyInRange(t *testing.T, hqe ledger.HistoryQueryExecutor, ns, desiredKey, falseKey string, expectedMatchCount int) {
	itr, err := hqe.GetHistoryForKey(ns, desiredKey)
	assert.NoError(t, err, "Error upon GetHistoryForKey()")
	scanner := itr.(*historyScanner)
	compositePartialKey := historydb.ConstructPartialCompositeHistoryKey(ns, falseKey, false)
	count := 0
	for {
		if !scanner.dbItr.Next() {
			break
		}
		historyKey := scanner.dbItr.Key()
		if bytes.Contains(historyKey, compositePartialKey) {
			count++
		}
	}
	assert.Equal(t, expectedMatchCount, count)
}



func exhaustQueryResultsForKey(t *testing.T, hqe ledger.HistoryQueryExecutor, ns, key string) {
	itr, err := hqe.GetHistoryForKey(ns, key)
	assert.NoError(t, err, "Error upon GetHistoryForKey()")
	for {
		kmod, _ := itr.Next()
		if kmod == nil {
			break
		}
	}
}




func testutilVerifyWarningMsg(t *testing.T, fakeLogger *fakes.HistoryleveldbLogger, ns, falseKey, warnMsg, decodeErrMsg string) {
	
	compositePartialKey := historydb.ConstructPartialCompositeHistoryKey(ns, falseKey, false)
	for i := 0; i < fakeLogger.WarnfCallCount(); i++ {
		msg, args := fakeLogger.WarnfArgsForCall(i)

		if !bytes.Contains(args[0].([]byte), compositePartialKey) {
			
			continue
		}

		assert.Equal(t, warnMsg, msg)

		if decodeErrMsg == "" {
			assert.Equal(t, 2, len(args))
		} else {
			assert.Equal(t, 3, len(args))
			assert.Equal(t, decodeErrMsg, fmt.Sprintf("%s", args[2]))
		}
		return
	}
	assert.Fail(t, fmt.Sprintf("did not find false key %s in range", falseKey))
}
