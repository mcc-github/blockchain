/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package historyleveldb

import (
	"os"
	"strconv"
	"testing"

	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	util2 "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	viper.Set("peer.fileSystemPath", "/tmp/blockchain/ledgertests/kvledger/history/historydb/historyleveldb")
	flogging.SetModuleLevel("leveldbhelper", "debug")
	flogging.SetModuleLevel("historyleveldb", "debug")
	os.Exit(m.Run())
}


func TestSavepoint(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()

	
	savepoint, err := env.testHistoryDB.GetLastSavepoint()
	testutil.AssertNoError(t, err, "Error upon historyDatabase.GetLastSavepoint()")
	testutil.AssertNil(t, savepoint)

	
	status, blockNum, err := env.testHistoryDB.ShouldRecover(0)
	testutil.AssertNoError(t, err, "Error upon historyDatabase.ShouldRecover()")
	testutil.AssertEquals(t, status, true)
	testutil.AssertEquals(t, blockNum, uint64(0))

	bg, gb := testutil.NewBlockGenerator(t, "testLedger", false)
	testutil.AssertNoError(t, env.testHistoryDB.Commit(gb), "")
	
	savepoint, err = env.testHistoryDB.GetLastSavepoint()
	testutil.AssertNoError(t, err, "Error upon historyDatabase.GetLastSavepoint()")
	testutil.AssertEquals(t, savepoint.BlockNum, uint64(0))

	
	txid := util2.GenerateUUID()
	simulator, _ := env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value1"))
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimResBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimResBytes})
	testutil.AssertNoError(t, env.testHistoryDB.Commit(block1), "")
	savepoint, err = env.testHistoryDB.GetLastSavepoint()
	testutil.AssertNoError(t, err, "Error upon historyDatabase.GetLastSavepoint()")
	testutil.AssertEquals(t, savepoint.BlockNum, uint64(1))

	
	status, blockNum, err = env.testHistoryDB.ShouldRecover(1)
	testutil.AssertNoError(t, err, "Error upon historyDatabase.ShouldRecover()")
	testutil.AssertEquals(t, status, false)
	testutil.AssertEquals(t, blockNum, uint64(2))

	
	txid = util2.GenerateUUID()
	simulator, _ = env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value2"))
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimResBytes, _ = simRes.GetPubSimulationBytes()
	block2 := bg.NextBlock([][]byte{pubSimResBytes})

	
	env.testHistoryDB.CommitLostBlock(&ledger.BlockAndPvtData{Block: block2})
	savepoint, err = env.testHistoryDB.GetLastSavepoint()
	testutil.AssertNoError(t, err, "Error upon historyDatabase.GetLastSavepoint()")
	testutil.AssertEquals(t, savepoint.BlockNum, uint64(2))

	
	status, blockNum, err = env.testHistoryDB.ShouldRecover(10)
	testutil.AssertNoError(t, err, "Error upon historyDatabase.ShouldRecover()")
	testutil.AssertEquals(t, status, true)
	testutil.AssertEquals(t, blockNum, uint64(3))
}

func TestHistory(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	provider := env.testBlockStorageEnv.provider
	ledger1id := "ledger1"
	store1, err := provider.OpenBlockStore(ledger1id)
	testutil.AssertNoError(t, err, "Error upon provider.OpenBlockStore()")
	defer store1.Shutdown()

	bg, gb := testutil.NewBlockGenerator(t, ledger1id, false)
	testutil.AssertNoError(t, store1.AddBlock(gb), "")
	testutil.AssertNoError(t, env.testHistoryDB.Commit(gb), "")

	
	txid := util2.GenerateUUID()
	simulator, _ := env.txmgr.NewTxSimulator(txid)
	value1 := []byte("value1")
	simulator.SetState("ns1", "key7", value1)
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimResBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimResBytes})
	err = store1.AddBlock(block1)
	testutil.AssertNoError(t, err, "")
	err = env.testHistoryDB.Commit(block1)
	testutil.AssertNoError(t, err, "")

	
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
	testutil.AssertNoError(t, err, "")
	err = env.testHistoryDB.Commit(block2)
	testutil.AssertNoError(t, err, "")

	
	txid = util2.GenerateUUID()
	simulator, _ = env.txmgr.NewTxSimulator(txid)
	simulator.DeleteState("ns1", "key7")
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimResBytes, _ = simRes.GetPubSimulationBytes()
	block3 := bg.NextBlock([][]byte{pubSimResBytes})
	err = store1.AddBlock(block3)
	testutil.AssertNoError(t, err, "")
	err = env.testHistoryDB.Commit(block3)
	testutil.AssertNoError(t, err, "")
	t.Logf("Inserted all 3 blocks")

	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(store1)
	testutil.AssertNoError(t, err, "Error upon NewHistoryQueryExecutor")

	itr, err2 := qhistory.GetHistoryForKey("ns1", "key7")
	testutil.AssertNoError(t, err2, "Error upon GetHistoryForKey()")

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
			testutil.AssertEquals(t, retrievedValue, expectedValue)
			testutil.AssertNotEquals(t, retrievedTimestamp, nil)
			testutil.AssertEquals(t, retrievedIsDelete, false)
		} else {
			testutil.AssertEquals(t, retrievedValue, nil)
			testutil.AssertNotEquals(t, retrievedTimestamp, nil)
			testutil.AssertEquals(t, retrievedIsDelete, true)
		}
	}
	testutil.AssertEquals(t, count, 4)
}

func TestHistoryForInvalidTran(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	provider := env.testBlockStorageEnv.provider
	ledger1id := "ledger1"
	store1, err := provider.OpenBlockStore(ledger1id)
	testutil.AssertNoError(t, err, "Error upon provider.OpenBlockStore()")
	defer store1.Shutdown()

	bg, gb := testutil.NewBlockGenerator(t, ledger1id, false)
	testutil.AssertNoError(t, store1.AddBlock(gb), "")
	testutil.AssertNoError(t, env.testHistoryDB.Commit(gb), "")

	
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
	testutil.AssertNoError(t, err, "")
	err = env.testHistoryDB.Commit(block1)
	testutil.AssertNoError(t, err, "")

	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(store1)
	testutil.AssertNoError(t, err, "Error upon NewHistoryQueryExecutor")

	itr, err2 := qhistory.GetHistoryForKey("ns1", "key7")
	testutil.AssertNoError(t, err2, "Error upon GetHistoryForKey()")

	
	kmod, _ := itr.Next()
	testutil.AssertNil(t, kmod)
}


func TestHistoryDisabled(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	viper.Set("ledger.history.enableHistoryDatabase", "false")
	
	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(nil)
	testutil.AssertNoError(t, err, "Error upon NewHistoryQueryExecutor")
	_, err2 := qhistory.GetHistoryForKey("ns1", "key7")
	testutil.AssertError(t, err2, "Error should have been returned for GetHistoryForKey() when history disabled")
}



func TestGenesisBlockNoError(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	block, err := configtxtest.MakeGenesisBlock("test_chainid")
	testutil.AssertNoError(t, err, "")
	err = env.testHistoryDB.Commit(block)
	testutil.AssertNoError(t, err, "")
}



func TestHistoryWithKeyContainingNilBytes(t *testing.T) {
	env := newTestHistoryEnv(t)
	defer env.cleanup()
	provider := env.testBlockStorageEnv.provider
	ledger1id := "ledger1"
	store1, err := provider.OpenBlockStore(ledger1id)
	testutil.AssertNoError(t, err, "Error upon provider.OpenBlockStore()")
	defer store1.Shutdown()

	bg, gb := testutil.NewBlockGenerator(t, ledger1id, false)
	testutil.AssertNoError(t, store1.AddBlock(gb), "")
	testutil.AssertNoError(t, env.testHistoryDB.Commit(gb), "")

	
	txid := util2.GenerateUUID()
	simulator, _ := env.txmgr.NewTxSimulator(txid)
	simulator.SetState("ns1", "key", []byte("value1")) 
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimResBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimResBytes})
	err = store1.AddBlock(block1)
	testutil.AssertNoError(t, err, "")
	err = env.testHistoryDB.Commit(block1)
	testutil.AssertNoError(t, err, "")

	
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
	
	simulator2.SetState("ns1", "key\x00\x01\x01\x15", []byte("dummyVal1"))
	simulator2.SetState("ns1", "\x00key\x00\x01\x01\x15", []byte("dummyVal2"))
	simulator2.Done()
	simRes2, _ := simulator2.GetTxSimulationResults()
	pubSimResBytes2, _ := simRes2.GetPubSimulationBytes()
	simulationResults = append(simulationResults, pubSimResBytes2)
	block2 := bg.NextBlock(simulationResults)
	err = store1.AddBlock(block2)
	testutil.AssertNoError(t, err, "")
	err = env.testHistoryDB.Commit(block2)
	testutil.AssertNoError(t, err, "")

	qhistory, err := env.testHistoryDB.NewHistoryQueryExecutor(store1)
	testutil.AssertNoError(t, err, "Error upon NewHistoryQueryExecutor")
	testutilVerifyResults(t, qhistory, "ns1", "key", []string{"value1", "value2"})
	testutilVerifyResults(t, qhistory, "ns1", "key\x00\x01\x01\x15", []string{"dummyVal1"})
	testutilVerifyResults(t, qhistory, "ns1", "\x00key\x00\x01\x01\x15", []string{"dummyVal2"})
}

func testutilVerifyResults(t *testing.T, hqe ledger.HistoryQueryExecutor, ns, key string, expectedVals []string) {
	itr, err := hqe.GetHistoryForKey(ns, key)
	testutil.AssertNoError(t, err, "Error upon GetHistoryForKey()")
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
	testutil.AssertEquals(t, retrievedVals, expectedVals)
}
