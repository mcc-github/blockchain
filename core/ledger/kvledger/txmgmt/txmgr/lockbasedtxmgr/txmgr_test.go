/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lockbasedtxmgr

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	btltestutil "github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/testutil"
	ledgertestutil "github.com/mcc-github/blockchain/core/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ledgertestutil.SetupCoreYAMLConfig()
	flogging.SetModuleLevel("lockbasedtxmgr", "debug")
	flogging.SetModuleLevel("statevalidator", "debug")
	flogging.SetModuleLevel("statebasedval", "debug")
	flogging.SetModuleLevel("statecouchdb", "debug")
	flogging.SetModuleLevel("valimpl", "debug")
	flogging.SetModuleLevel("pvtstatepurgemgmt", "debug")
	flogging.SetModuleLevel("valinternal", "debug")

	viper.Set("peer.fileSystemPath", "/tmp/blockchain/ledgertests/kvledger/txmgmt/txmgr/lockbasedtxmgr")
	os.Exit(m.Run())
}

func TestTxSimulatorWithNoExistingData(t *testing.T) {
	
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxsimulatorwithnoexistingdata"
		testEnv.init(t, testLedgerID, nil)
		testTxSimulatorWithNoExistingData(t, testEnv)
		testEnv.cleanup()
	}
}

func testTxSimulatorWithNoExistingData(t *testing.T, env testEnv) {
	txMgr := env.getTxMgr()
	s, _ := txMgr.NewTxSimulator("test_txid")
	value, err := s.GetState("ns1", "key1")
	testutil.AssertNoError(t, err, fmt.Sprintf("Error in GetState(): %s", err))
	testutil.AssertNil(t, value)

	s.SetState("ns1", "key1", []byte("value1"))
	s.SetState("ns1", "key2", []byte("value2"))
	s.SetState("ns2", "key3", []byte("value3"))
	s.SetState("ns2", "key4", []byte("value4"))

	value, _ = s.GetState("ns2", "key3")
	testutil.AssertNil(t, value)

	simulationResults, err := s.GetTxSimulationResults()
	assert.NoError(t, err)
	assert.Nil(t, simulationResults.PvtSimulationResults)
}

func TestTxSimulatorGetResults(t *testing.T) {
	testEnv := testEnvsMap[levelDBtestEnvName]
	testEnv.init(t, "testLedger", nil)
	defer testEnv.cleanup()
	txMgr := testEnv.getTxMgr()
	populateCollConfigForTest(t, txMgr.(*LockBasedTxMgr),
		[]collConfigkey{
			{"ns1", "coll1"},
			{"ns1", "coll3"},
			{"ns2", "coll2"},
			{"ns3", "coll3"},
		},
		version.NewHeight(1, 1),
	)

	var err error

	
	simulator, _ := testEnv.getTxMgr().NewTxSimulator("test_txid1")
	simulator.GetState("ns1", "key1")
	_, err = simulator.GetPrivateData("ns1", "coll1", "key1")
	assert.NoError(t, err)
	simulator.SetState("ns1", "key1", []byte("value1"))
	
	simulationResults1, err := simulator.GetTxSimulationResults()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(simulationResults1.PubSimulationResults.NsRwset))
	
	buff1 := new(bytes.Buffer)
	assert.NoError(t, gob.NewEncoder(buff1).Encode(simulationResults1))
	frozenSimulationResults1 := &ledger.TxSimulationResults{}
	assert.NoError(t, gob.NewDecoder(buff1).Decode(&frozenSimulationResults1))

	
	simulator.GetState("ns2", "key2")
	simulator.GetPrivateData("ns2", "coll2", "key2")
	simulator.SetState("ns2", "key2", []byte("value2"))
	
	_, err = simulator.GetTxSimulationResults()
	assert.Error(t, err) 
	
	assert.Equal(t, frozenSimulationResults1, simulationResults1)

	
	simulator.Done()
	_, err = simulator.GetState("ns3", "key3")
	assert.Errorf(t, err, "An error is expected when using simulator to get/set data after calling `Done` function()")
	err = simulator.SetState("ns3", "key3", []byte("value3"))
	assert.Errorf(t, err, "An error is expected when using simulator to get/set data after calling `Done` function()")
	_, err = simulator.GetPrivateData("ns3", "coll3", "key3")
	assert.Errorf(t, err, "An error is expected when using simulator to get/set data after calling `Done` function()")
	err = simulator.SetPrivateData("ns3", "coll3", "key3", []byte("value3"))
	assert.Errorf(t, err, "An error is expected when using simulator to get/set data after calling `Done` function()")
}

func TestTxSimulatorWithExistingData(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Run(testEnv.getName(), func(t *testing.T) {
			testLedgerID := "testtxsimulatorwithexistingdata"
			testEnv.init(t, testLedgerID, nil)
			testTxSimulatorWithExistingData(t, testEnv)
			testEnv.cleanup()
		})
	}
}

func testTxSimulatorWithExistingData(t *testing.T, env testEnv) {
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)
	
	s1, _ := txMgr.NewTxSimulator("test_tx1")
	s1.SetState("ns1", "key1", []byte("value1"))
	s1.SetState("ns1", "key2", []byte("value2"))
	s1.SetState("ns2", "key3", []byte("value3"))
	s1.SetState("ns2", "key4", []byte("value4"))
	s1.Done()
	
	txRWSet1, _ := s1.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet1.PubSimulationResults)

	
	s2, _ := txMgr.NewTxSimulator("test_tx2")
	value, _ := s2.GetState("ns1", "key1")
	testutil.AssertEquals(t, value, []byte("value1"))
	s2.SetState("ns1", "key1", []byte("value1_1"))
	s2.DeleteState("ns2", "key3")
	value, _ = s2.GetState("ns1", "key1")
	testutil.AssertEquals(t, value, []byte("value1"))
	s2.Done()
	
	txRWSet2, _ := s2.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet2.PubSimulationResults)

	
	s3, _ := txMgr.NewTxSimulator("test_tx3")
	value, _ = s3.GetState("ns1", "key1")
	testutil.AssertEquals(t, value, []byte("value1_1"))
	value, _ = s3.GetState("ns2", "key3")
	testutil.AssertEquals(t, value, nil)
	s3.Done()

	
	vv, _ := env.getVDB().GetState("ns1", "key1")
	testutil.AssertEquals(t, vv.Version, version.NewHeight(2, 0))
	vv, _ = env.getVDB().GetState("ns1", "key2")
	testutil.AssertEquals(t, vv.Version, version.NewHeight(1, 0))
}

func TestTxValidation(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxvalidation"
		testEnv.init(t, testLedgerID, nil)
		testTxValidation(t, testEnv)
		testEnv.cleanup()
	}
}

func testTxValidation(t *testing.T, env testEnv) {
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)
	
	s1, _ := txMgr.NewTxSimulator("test_tx1")
	s1.SetState("ns1", "key1", []byte("value1"))
	s1.SetState("ns1", "key2", []byte("value2"))
	s1.SetState("ns2", "key3", []byte("value3"))
	s1.SetState("ns2", "key4", []byte("value4"))
	s1.Done()
	
	txRWSet1, _ := s1.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet1.PubSimulationResults)

	
	
	s2, _ := txMgr.NewTxSimulator("test_tx2")
	value, _ := s2.GetState("ns1", "key1")
	testutil.AssertEquals(t, value, []byte("value1"))

	s2.SetState("ns1", "key1", []byte("value1_2"))
	s2.DeleteState("ns2", "key3")
	s2.Done()

	
	
	s3, _ := txMgr.NewTxSimulator("test_tx3")
	s3.GetState("ns1", "key1")
	s3.SetState("ns1", "key1", []byte("value1_3"))
	s3.Done()

	
	
	s4, _ := txMgr.NewTxSimulator("test_tx4")
	s4.GetState("ns2", "key3")
	s4.DeleteState("ns2", "key3")
	s4.Done()

	
	
	s5, _ := txMgr.NewTxSimulator("test_tx5")
	s5.SetState("ns1", "key1", []byte("new_value"))
	s5.GetState("ns1", "key1")
	s5.Done()

	
	
	s6, _ := txMgr.NewTxSimulator("test_tx6")
	s6.SetState("ns1", "new_key", []byte("new_value"))
	s6.Done()

	
	
	
	
	
	

	
	txRWSet2, _ := s2.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet2.PubSimulationResults)

	
	txRWSet3, _ := s3.GetTxSimulationResults()
	txMgrHelper.checkRWsetInvalid(txRWSet3.PubSimulationResults)

	txRWSet4, _ := s4.GetTxSimulationResults()
	txMgrHelper.checkRWsetInvalid(txRWSet4.PubSimulationResults)

	txRWSet5, _ := s5.GetTxSimulationResults()
	txMgrHelper.checkRWsetInvalid(txRWSet5.PubSimulationResults)

	
	txRWSet6, _ := s6.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet6.PubSimulationResults)
}

func TestTxPhantomValidation(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxphantomvalidation"
		testEnv.init(t, testLedgerID, nil)
		testTxPhantomValidation(t, testEnv)
		testEnv.cleanup()
	}
}

func testTxPhantomValidation(t *testing.T, env testEnv) {
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)
	
	s1, _ := txMgr.NewTxSimulator("test_tx1")
	s1.SetState("ns", "key1", []byte("value1"))
	s1.SetState("ns", "key2", []byte("value2"))
	s1.SetState("ns", "key3", []byte("value3"))
	s1.SetState("ns", "key4", []byte("value4"))
	s1.SetState("ns", "key5", []byte("value5"))
	s1.SetState("ns", "key6", []byte("value6"))
	
	txRWSet1, _ := s1.GetTxSimulationResults()
	s1.Done() 
	txMgrHelper.validateAndCommitRWSet(txRWSet1.PubSimulationResults)

	
	s2, _ := txMgr.NewTxSimulator("test_tx2")
	itr2, _ := s2.GetStateRangeScanIterator("ns", "key2", "key5")
	for {
		if result, _ := itr2.Next(); result == nil {
			break
		}
	}
	s2.DeleteState("ns", "key3")
	txRWSet2, _ := s2.GetTxSimulationResults()
	s2.Done()

	
	s3, _ := txMgr.NewTxSimulator("test_tx3")
	itr3, _ := s3.GetStateRangeScanIterator("ns", "key2", "key5")
	for {
		if result, _ := itr3.Next(); result == nil {
			break
		}
	}
	s3.SetState("ns", "key3", []byte("value3_new"))
	txRWSet3, _ := s3.GetTxSimulationResults()
	s3.Done()
	
	s4, _ := txMgr.NewTxSimulator("test_tx4")
	itr4, _ := s4.GetStateRangeScanIterator("ns", "key4", "key6")
	for {
		if result, _ := itr4.Next(); result == nil {
			break
		}
	}
	s4.SetState("ns", "key3", []byte("value3_new"))
	txRWSet4, _ := s4.GetTxSimulationResults()
	s4.Done()

	
	txMgrHelper.validateAndCommitRWSet(txRWSet2.PubSimulationResults)
	
	txMgrHelper.checkRWsetInvalid(txRWSet3.PubSimulationResults)
	
	txMgrHelper.validateAndCommitRWSet(txRWSet4.PubSimulationResults)
}

func TestIterator(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())

		testLedgerID := "testiterator.1"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 2, 7)
		testEnv.cleanup()

		testLedgerID = "testiterator.2"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 1, 11)
		testEnv.cleanup()

		testLedgerID = "testiterator.3"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 0, 0)
		testEnv.cleanup()

		testLedgerID = "testiterator.4"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 5, 0)
		testEnv.cleanup()

		testLedgerID = "testiterator.5"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 0, 5)
		testEnv.cleanup()
	}
}

func testIterator(t *testing.T, env testEnv, numKeys int, startKeyNum int, endKeyNum int) {
	cID := "cid"
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)
	s, _ := txMgr.NewTxSimulator("test_tx1")
	for i := 1; i <= numKeys; i++ {
		k := createTestKey(i)
		v := createTestValue(i)
		t.Logf("Adding k=[%s], v=[%s]", k, v)
		s.SetState(cID, k, v)
	}
	s.Done()
	
	txRWSet, _ := s.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet.PubSimulationResults)

	var startKey string
	var endKey string
	var begin int
	var end int

	if startKeyNum != 0 {
		begin = startKeyNum
		startKey = createTestKey(startKeyNum)
	} else {
		begin = 1 
		startKey = ""
	}

	if endKeyNum != 0 {
		endKey = createTestKey(endKeyNum)
		end = endKeyNum
	} else {
		endKey = ""
		end = numKeys + 1 
	}

	expectedCount := end - begin

	queryExecuter, _ := txMgr.NewQueryExecutor("test_tx2")
	itr, _ := queryExecuter.GetStateRangeScanIterator(cID, startKey, endKey)
	count := 0
	for {
		kv, _ := itr.Next()
		if kv == nil {
			break
		}
		keyNum := begin + count
		k := kv.(*queryresult.KV).Key
		v := kv.(*queryresult.KV).Value
		t.Logf("Retrieved k=%s, v=%s at count=%d start=%s end=%s", k, v, count, startKey, endKey)
		testutil.AssertEquals(t, k, createTestKey(keyNum))
		testutil.AssertEquals(t, v, createTestValue(keyNum))
		count++
	}
	testutil.AssertEquals(t, count, expectedCount)
}

func TestIteratorPaging(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())

		
		testLedgerID := "testiterator.1"
		testEnv.init(t, testLedgerID, nil)
		testIteratorPagingInit(t, testEnv, 10)
		returnKeys := []string{"key_002", "key_003"}
		nextStartKey := testIteratorPaging(t, testEnv, 10, "key_002", "key_007", int32(2), returnKeys)
		returnKeys = []string{"key_004", "key_005"}
		nextStartKey = testIteratorPaging(t, testEnv, 10, nextStartKey, "key_007", int32(2), returnKeys)
		returnKeys = []string{"key_006"}
		testIteratorPaging(t, testEnv, 10, nextStartKey, "key_007", int32(2), returnKeys)
		testEnv.cleanup()
	}
}

func testIteratorPagingInit(t *testing.T, env testEnv, numKeys int) {
	cID := "cid"
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)
	s, _ := txMgr.NewTxSimulator("test_tx1")
	for i := 1; i <= numKeys; i++ {
		k := createTestKey(i)
		v := createTestValue(i)
		t.Logf("Adding k=[%s], v=[%s]", k, v)
		s.SetState(cID, k, v)
	}
	s.Done()
	
	txRWSet, _ := s.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet.PubSimulationResults)
}

func testIteratorPaging(t *testing.T, env testEnv, numKeys int, startKey, endKey string,
	limit int32, expectedKeys []string) string {
	cID := "cid"
	txMgr := env.getTxMgr()

	queryOptions := make(map[string]interface{})
	if limit != 0 {
		queryOptions["limit"] = limit
	}

	queryExecuter, _ := txMgr.NewQueryExecutor("test_tx2")
	itr, _ := queryExecuter.GetStateRangeScanIteratorWithMetadata(cID, startKey, endKey, queryOptions)

	
	testItrWithoutClose(t, itr, expectedKeys)

	returnBookmark := ""
	if limit > 0 {
		returnBookmark = itr.GetBookmarkAndClose()
	}

	return returnBookmark
}


func testItrWithoutClose(t *testing.T, itr ledger.QueryResultsIterator, expectedKeys []string) {
	for _, expectedKey := range expectedKeys {
		queryResult, err := itr.Next()
		testutil.AssertNoError(t, err, "An unexpected error was thrown during iterator Next()")
		vkv := queryResult.(*queryresult.KV)
		key := vkv.Key
		testutil.AssertEquals(t, key, expectedKey)
	}
	queryResult, err := itr.Next()
	testutil.AssertNoError(t, err, "An unexpected error was thrown during iterator Next()")
	testutil.AssertNil(t, queryResult)
}

func TestIteratorWithDeletes(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testiteratorwithdeletes"
		testEnv.init(t, testLedgerID, nil)
		testIteratorWithDeletes(t, testEnv)
		testEnv.cleanup()
	}
}

func testIteratorWithDeletes(t *testing.T, env testEnv) {
	cID := "cid"
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)
	s, _ := txMgr.NewTxSimulator("test_tx1")
	for i := 1; i <= 10; i++ {
		k := createTestKey(i)
		v := createTestValue(i)
		t.Logf("Adding k=[%s], v=[%s]", k, v)
		s.SetState(cID, k, v)
	}
	s.Done()
	
	txRWSet1, _ := s.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet1.PubSimulationResults)

	s, _ = txMgr.NewTxSimulator("test_tx2")
	s.DeleteState(cID, createTestKey(4))
	s.Done()
	
	txRWSet2, _ := s.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet2.PubSimulationResults)

	queryExecuter, _ := txMgr.NewQueryExecutor("test_tx3")
	itr, _ := queryExecuter.GetStateRangeScanIterator(cID, createTestKey(3), createTestKey(6))
	defer itr.Close()
	kv, _ := itr.Next()
	testutil.AssertEquals(t, kv.(*queryresult.KV).Key, createTestKey(3))
	kv, _ = itr.Next()
	testutil.AssertEquals(t, kv.(*queryresult.KV).Key, createTestKey(5))
}

func TestTxValidationWithItr(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxvalidationwithitr"
		testEnv.init(t, testLedgerID, nil)
		testTxValidationWithItr(t, testEnv)
		testEnv.cleanup()
	}
}

func testTxValidationWithItr(t *testing.T, env testEnv) {
	cID := "cid"
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)

	
	s1, _ := txMgr.NewTxSimulator("test_tx1")
	for i := 1; i <= 10; i++ {
		k := createTestKey(i)
		v := createTestValue(i)
		t.Logf("Adding k=[%s], v=[%s]", k, v)
		s1.SetState(cID, k, v)
	}
	s1.Done()
	
	txRWSet1, _ := s1.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet1.PubSimulationResults)

	
	s2, _ := txMgr.NewTxSimulator("test_tx2")
	itr, _ := s2.GetStateRangeScanIterator(cID, createTestKey(1), createTestKey(5))
	
	itr.Next()
	itr.Next()
	itr.Close()
	s2.Done()

	
	s3, _ := txMgr.NewTxSimulator("test_tx3")
	itr, _ = s3.GetStateRangeScanIterator(cID, createTestKey(4), createTestKey(6))
	
	itr.Next()
	itr.Next()
	itr.Close()
	s3.Done()

	
	s4, _ := txMgr.NewTxSimulator("test_tx4")
	s4.DeleteState(cID, createTestKey(5))
	s4.Done()

	
	txRWSet4, _ := s4.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet4.PubSimulationResults)

	
	txRWSet3, _ := s3.GetTxSimulationResults()
	txMgrHelper.checkRWsetInvalid(txRWSet3.PubSimulationResults)

	
	txRWSet2, _ := s2.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet2.PubSimulationResults)

}

func TestGetSetMultipeKeys(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testgetsetmultipekeys"
		testEnv.init(t, testLedgerID, nil)
		testGetSetMultipeKeys(t, testEnv)
		testEnv.cleanup()
	}
}

func testGetSetMultipeKeys(t *testing.T, env testEnv) {
	cID := "cid"
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)
	
	s1, _ := txMgr.NewTxSimulator("test_tx1")
	multipleKeyMap := make(map[string][]byte)
	for i := 1; i <= 10; i++ {
		k := createTestKey(i)
		v := createTestValue(i)
		multipleKeyMap[k] = v
	}
	s1.SetStateMultipleKeys(cID, multipleKeyMap)
	s1.Done()
	
	txRWSet, _ := s1.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet.PubSimulationResults)
	qe, _ := txMgr.NewQueryExecutor("test_tx2")
	defer qe.Done()
	multipleKeys := []string{}
	for k := range multipleKeyMap {
		multipleKeys = append(multipleKeys, k)
	}
	values, _ := qe.GetStateMultipleKeys(cID, multipleKeys)
	testutil.AssertEquals(t, len(values), 10)
	for i, v := range values {
		testutil.AssertEquals(t, v, multipleKeyMap[multipleKeys[i]])
	}

	s2, _ := txMgr.NewTxSimulator("test_tx3")
	defer s2.Done()
	values, _ = s2.GetStateMultipleKeys(cID, multipleKeys[5:7])
	testutil.AssertEquals(t, len(values), 2)
	for i, v := range values {
		testutil.AssertEquals(t, v, multipleKeyMap[multipleKeys[i+5]])
	}
}

func createTestKey(i int) string {
	if i == 0 {
		return ""
	}
	return fmt.Sprintf("key_%03d", i)
}

func createTestValue(i int) []byte {
	return []byte(fmt.Sprintf("value_%03d", i))
}


func TestExecuteQuery(t *testing.T) {

	for _, testEnv := range testEnvs {
		
		if testEnv.getName() == couchDBtestEnvName {
			t.Logf("Running test for TestEnv = %s", testEnv.getName())
			testLedgerID := "testexecutequery"
			testEnv.init(t, testLedgerID, nil)
			testExecuteQuery(t, testEnv)
			testEnv.cleanup()
		}
	}
}

func testExecuteQuery(t *testing.T, env testEnv) {

	type Asset struct {
		ID        string `json:"_id"`
		Rev       string `json:"_rev"`
		AssetName string `json:"asset_name"`
		Color     string `json:"color"`
		Size      string `json:"size"`
		Owner     string `json:"owner"`
	}

	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)

	s1, _ := txMgr.NewTxSimulator("test_tx1")

	s1.SetState("ns1", "key1", []byte("value1"))
	s1.SetState("ns1", "key2", []byte("value2"))
	s1.SetState("ns1", "key3", []byte("value3"))
	s1.SetState("ns1", "key4", []byte("value4"))
	s1.SetState("ns1", "key5", []byte("value5"))
	s1.SetState("ns1", "key6", []byte("value6"))
	s1.SetState("ns1", "key7", []byte("value7"))
	s1.SetState("ns1", "key8", []byte("value8"))

	s1.SetState("ns1", "key9", []byte(`{"asset_name":"marble1","color":"red","size":"25","owner":"jerry"}`))
	s1.SetState("ns1", "key10", []byte(`{"asset_name":"marble2","color":"blue","size":"10","owner":"bob"}`))
	s1.SetState("ns1", "key11", []byte(`{"asset_name":"marble3","color":"blue","size":"35","owner":"jerry"}`))
	s1.SetState("ns1", "key12", []byte(`{"asset_name":"marble4","color":"green","size":"15","owner":"bob"}`))
	s1.SetState("ns1", "key13", []byte(`{"asset_name":"marble5","color":"red","size":"35","owner":"jerry"}`))
	s1.SetState("ns1", "key14", []byte(`{"asset_name":"marble6","color":"blue","size":"25","owner":"bob"}`))

	s1.Done()

	
	txRWSet, _ := s1.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet.PubSimulationResults)

	queryExecuter, _ := txMgr.NewQueryExecutor("test_tx2")
	queryString := "{\"selector\":{\"owner\": {\"$eq\": \"bob\"}},\"limit\": 10,\"skip\": 0}"

	itr, err := queryExecuter.ExecuteQuery("ns1", queryString)
	testutil.AssertNoError(t, err, "Error upon ExecuteQuery()")
	counter := 0
	for {
		queryRecord, _ := itr.Next()
		if queryRecord == nil {
			break
		}
		
		assetResp := &Asset{}
		json.Unmarshal(queryRecord.(*queryresult.KV).Value, &assetResp)
		
		testutil.AssertEquals(t, assetResp.Owner, "bob")
		counter++
	}
	
	testutil.AssertEquals(t, counter, 3)
}


func TestExecutePaginatedQuery(t *testing.T) {

	for _, testEnv := range testEnvs {
		
		if testEnv.getName() == couchDBtestEnvName {
			t.Logf("Running test for TestEnv = %s", testEnv.getName())
			testLedgerID := "testexecutepaginatedquery"
			testEnv.init(t, testLedgerID, nil)
			testExecutePaginatedQuery(t, testEnv)
			testEnv.cleanup()
		}
	}
}

func testExecutePaginatedQuery(t *testing.T, env testEnv) {

	type Asset struct {
		ID        string `json:"_id"`
		Rev       string `json:"_rev"`
		AssetName string `json:"asset_name"`
		Color     string `json:"color"`
		Size      string `json:"size"`
		Owner     string `json:"owner"`
	}

	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)

	s1, _ := txMgr.NewTxSimulator("test_tx1")

	s1.SetState("ns1", "key1", []byte(`{"asset_name":"marble1","color":"red","size":"25","owner":"jerry"}`))
	s1.SetState("ns1", "key2", []byte(`{"asset_name":"marble2","color":"blue","size":"10","owner":"bob"}`))
	s1.SetState("ns1", "key3", []byte(`{"asset_name":"marble3","color":"blue","size":"35","owner":"jerry"}`))
	s1.SetState("ns1", "key4", []byte(`{"asset_name":"marble4","color":"green","size":"15","owner":"bob"}`))
	s1.SetState("ns1", "key5", []byte(`{"asset_name":"marble5","color":"red","size":"35","owner":"jerry"}`))
	s1.SetState("ns1", "key6", []byte(`{"asset_name":"marble6","color":"blue","size":"25","owner":"bob"}`))

	s1.Done()

	
	txRWSet, _ := s1.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet.PubSimulationResults)

	queryExecuter, _ := txMgr.NewQueryExecutor("test_tx2")
	queryString := `{"selector":{"owner":{"$eq":"bob"}}}`

	queryOptions := map[string]interface{}{
		"limit": int32(2),
	}

	itr, err := queryExecuter.ExecuteQueryWithMetadata("ns1", queryString, queryOptions)
	testutil.AssertNoError(t, err, "Error upon ExecuteQueryWithMetadata()")
	counter := 0
	for {
		queryRecord, _ := itr.Next()
		if queryRecord == nil {
			break
		}
		
		assetResp := &Asset{}
		json.Unmarshal(queryRecord.(*queryresult.KV).Value, &assetResp)
		
		testutil.AssertEquals(t, assetResp.Owner, "bob")
		counter++
	}
	
	testutil.AssertEquals(t, counter, 2)

	bookmark := itr.GetBookmarkAndClose()

	queryOptions = map[string]interface{}{
		"limit": int32(2),
	}
	if bookmark != "" {
		queryOptions["bookmark"] = bookmark
	}

	itr, err = queryExecuter.ExecuteQueryWithMetadata("ns1", queryString, queryOptions)
	testutil.AssertNoError(t, err, "Error upon ExecuteQuery()")
	counter = 0
	for {
		queryRecord, _ := itr.Next()
		if queryRecord == nil {
			break
		}
		
		assetResp := &Asset{}
		json.Unmarshal(queryRecord.(*queryresult.KV).Value, &assetResp)
		
		testutil.AssertEquals(t, assetResp.Owner, "bob")
		counter++
	}
	
	testutil.AssertEquals(t, counter, 1)
}

func TestValidateKey(t *testing.T) {
	nonUTF8Key := string([]byte{0xff, 0xff})
	dummyValue := []byte("dummyValue")
	for _, testEnv := range testEnvs {
		testLedgerID := "test.validate.key"
		testEnv.init(t, testLedgerID, nil)
		txSimulator, _ := testEnv.getTxMgr().NewTxSimulator("test_tx1")
		err := txSimulator.SetState("ns1", nonUTF8Key, dummyValue)
		if testEnv.getName() == levelDBtestEnvName {
			testutil.AssertNoError(t, err, "")
		}
		if testEnv.getName() == couchDBtestEnvName {
			testutil.AssertError(t, err, "")
		}
		testEnv.cleanup()
	}
}



func TestTxSimulatorUnsupportedTx(t *testing.T) {
	testEnv := testEnvs[0]
	testEnv.init(t, "testtxsimulatorunsupportedtx", nil)
	defer testEnv.cleanup()
	txMgr := testEnv.getTxMgr()
	populateCollConfigForTest(t, txMgr.(*LockBasedTxMgr),
		[]collConfigkey{
			{"ns1", "coll1"},
			{"ns1", "coll2"},
			{"ns1", "coll3"},
			{"ns1", "coll4"},
		},
		version.NewHeight(1, 1))

	simulator, _ := txMgr.NewTxSimulator("txid1")
	err := simulator.SetState("ns", "key", []byte("value"))
	testutil.AssertNoError(t, err, "")
	_, err = simulator.GetPrivateDataRangeScanIterator("ns1", "coll1", "startKey", "endKey")
	_, ok := err.(*txmgr.ErrUnsupportedTransaction)
	testutil.AssertEquals(t, ok, true)

	simulator, _ = txMgr.NewTxSimulator("txid2")
	_, err = simulator.GetPrivateDataRangeScanIterator("ns1", "coll1", "startKey", "endKey")
	testutil.AssertNoError(t, err, "")
	err = simulator.SetState("ns", "key", []byte("value"))
	_, ok = err.(*txmgr.ErrUnsupportedTransaction)
	testutil.AssertEquals(t, ok, true)

	queryOptions := map[string]interface{}{
		"limit": int32(2),
	}

	simulator, _ = txMgr.NewTxSimulator("txid3")
	err = simulator.SetState("ns", "key", []byte("value"))
	testutil.AssertNoError(t, err, "")
	_, err = simulator.GetStateRangeScanIteratorWithMetadata("ns1", "startKey", "endKey", queryOptions)
	_, ok = err.(*txmgr.ErrUnsupportedTransaction)
	testutil.AssertEquals(t, ok, true)

	simulator, _ = txMgr.NewTxSimulator("txid4")
	_, err = simulator.GetStateRangeScanIteratorWithMetadata("ns1", "startKey", "endKey", queryOptions)
	testutil.AssertNoError(t, err, "")
	err = simulator.SetState("ns", "key", []byte("value"))
	_, ok = err.(*txmgr.ErrUnsupportedTransaction)
	testutil.AssertEquals(t, ok, true)

}


func TestTxSimulatorQueryUnsupportedTx(t *testing.T) {

	for _, testEnv := range testEnvs {
		
		if testEnv.getName() == couchDBtestEnvName {
			t.Logf("Running test for TestEnv = %s", testEnv.getName())
			testLedgerID := "testtxsimulatorunsupportedtxqueries"
			testEnv.init(t, testLedgerID, nil)
			testTxSimulatorQueryUnsupportedTx(t, testEnv)
			testEnv.cleanup()
		}
	}
}

func testTxSimulatorQueryUnsupportedTx(t *testing.T, env testEnv) {

	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)

	s1, _ := txMgr.NewTxSimulator("test_tx1")

	s1.SetState("ns1", "key1", []byte(`{"asset_name":"marble1","color":"red","size":"25","owner":"jerry"}`))

	s1.Done()

	
	txRWSet, _ := s1.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet.PubSimulationResults)

	queryString := `{"selector":{"owner":{"$eq":"bob"}}}`
	queryOptions := map[string]interface{}{
		"limit": int32(2),
	}

	simulator, _ := txMgr.NewTxSimulator("txid1")
	err := simulator.SetState("ns1", "key1", []byte(`{"asset_name":"marble1","color":"red","size":"25","owner":"jerry"}`))
	testutil.AssertNoError(t, err, "")
	_, err = simulator.ExecuteQueryWithMetadata("ns1", queryString, queryOptions)
	_, ok := err.(*txmgr.ErrUnsupportedTransaction)
	testutil.AssertEquals(t, ok, true)

	simulator, _ = txMgr.NewTxSimulator("txid2")
	_, err = simulator.ExecuteQueryWithMetadata("ns1", queryString, queryOptions)
	testutil.AssertNoError(t, err, "")
	err = simulator.SetState("ns1", "key1", []byte(`{"asset_name":"marble1","color":"red","size":"25","owner":"jerry"}`))
	_, ok = err.(*txmgr.ErrUnsupportedTransaction)
	testutil.AssertEquals(t, ok, true)

}

func TestTxSimulatorMissingPvtdata(t *testing.T) {
	testEnv := testEnvs[0]
	testEnv.init(t, "TestTxSimulatorUnsupportedTxQueries", nil)
	defer testEnv.cleanup()

	txMgr := testEnv.getTxMgr()
	populateCollConfigForTest(t, txMgr.(*LockBasedTxMgr),
		[]collConfigkey{
			{"ns1", "coll1"},
			{"ns1", "coll2"},
			{"ns1", "coll3"},
			{"ns1", "coll4"},
		},
		version.NewHeight(1, 1),
	)

	db := testEnv.getVDB()
	updateBatch := privacyenabledstate.NewUpdateBatch()
	updateBatch.HashUpdates.Put("ns1", "coll1", util.ComputeStringHash("key1"), util.ComputeStringHash("value1"), version.NewHeight(1, 1))
	updateBatch.PvtUpdates.Put("ns1", "coll1", "key1", []byte("value1"), version.NewHeight(1, 1))
	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(1, 1))

	simulator, _ := txMgr.NewTxSimulator("testTxid1")
	val, _ := simulator.GetPrivateData("ns1", "coll1", "key1")
	testutil.AssertEquals(t, val, []byte("value1"))
	simulator.Done()

	updateBatch = privacyenabledstate.NewUpdateBatch()
	updateBatch.HashUpdates.Put("ns1", "coll1", util.ComputeStringHash("key1"), util.ComputeStringHash("value1"), version.NewHeight(2, 1))
	updateBatch.HashUpdates.Put("ns1", "coll2", util.ComputeStringHash("key2"), util.ComputeStringHash("value2"), version.NewHeight(2, 1))
	updateBatch.HashUpdates.Put("ns1", "coll3", util.ComputeStringHash("key3"), util.ComputeStringHash("value3"), version.NewHeight(2, 1))
	updateBatch.PvtUpdates.Put("ns1", "coll3", "key3", []byte("value3"), version.NewHeight(2, 1))
	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(2, 1))
	simulator, _ = txMgr.NewTxSimulator("testTxid2")
	defer simulator.Done()

	val, err := simulator.GetPrivateData("ns1", "coll1", "key1")
	_, ok := err.(*txmgr.ErrPvtdataNotAvailable)
	testutil.AssertEquals(t, ok, true)

	val, err = simulator.GetPrivateData("ns1", "coll2", "key2")
	_, ok = err.(*txmgr.ErrPvtdataNotAvailable)
	testutil.AssertEquals(t, ok, true)

	val, err = simulator.GetPrivateData("ns1", "coll3", "key3")
	testutil.AssertEquals(t, val, []byte("value3"))

	val, err = simulator.GetPrivateData("ns1", "coll4", "key4")
	testutil.AssertNoError(t, err, "")
	testutil.AssertNil(t, val)
}

func TestDeleteOnCursor(t *testing.T) {
	cID := "cid"
	env := testEnvs[0]
	env.init(t, "TestDeleteOnCursor", nil)
	defer env.cleanup()

	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)

	
	s, _ := txMgr.NewTxSimulator("test_tx1")
	for i := 1; i <= 10; i++ {
		k := createTestKey(i)
		v := createTestValue(i)
		t.Logf("Adding k=[%s], v=[%s]", k, v)
		s.SetState(cID, k, v)
	}
	s.Done()
	txRWSet1, _ := s.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet1.PubSimulationResults)

	
	s2, _ := txMgr.NewTxSimulator("test_tx2")
	itr2, _ := s2.GetStateRangeScanIterator(cID, createTestKey(1), createTestKey(5))
	for i := 1; i <= 4; i++ {
		kv, err := itr2.Next()
		testutil.AssertNoError(t, err, "")
		testutil.AssertNotNil(t, kv)
		key := kv.(*queryresult.KV).Key
		s2.DeleteState(cID, key)
	}
	itr2.Close()
	s2.Done()
	txRWSet2, _ := s2.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet2.PubSimulationResults)

	
	s3, _ := txMgr.NewTxSimulator("test_tx3")
	itr3, _ := s3.GetStateRangeScanIterator(cID, createTestKey(1), createTestKey(10))
	kv, err := itr3.Next()
	testutil.AssertNoError(t, err, "")
	testutil.AssertNotNil(t, kv)
	key := kv.(*queryresult.KV).Key
	testutil.AssertEquals(t, key, "key_005")
	itr3.Close()
	s3.Done()
}

func TestTxSimulatorMissingPvtdataExpiry(t *testing.T) {
	ledgerid := "TestTxSimulatorMissingPvtdataExpiry"
	testEnv := testEnvs[0]
	cs := btltestutil.NewMockCollectionStore()
	cs.SetBTL("ns", "coll", 1)
	testEnv.init(t, ledgerid, pvtdatapolicy.ConstructBTLPolicy(cs))
	defer testEnv.cleanup()

	txMgr := testEnv.getTxMgr()
	populateCollConfigForTest(t, txMgr.(*LockBasedTxMgr), []collConfigkey{{"ns", "coll"}}, version.NewHeight(1, 1))

	viper.Set(fmt.Sprintf("ledger.pvtdata.btlpolicy.%s.ns.coll", ledgerid), 1)
	bg, _ := testutil.NewBlockGenerator(t, ledgerid, false)

	blkAndPvtdata := prepareNextBlockForTest(t, txMgr, bg, "txid-1",
		map[string]string{"pubkey1": "pub-value1"}, map[string]string{"pvtkey1": "pvt-value1"})
	testutil.AssertNoError(t, txMgr.ValidateAndPrepare(blkAndPvtdata, true), "")
	testutil.AssertNoError(t, txMgr.Commit(), "")

	simulator, _ := txMgr.NewTxSimulator("tx-tmp")
	pvtval, err := simulator.GetPrivateData("ns", "coll", "pvtkey1")
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, pvtval, []byte("pvt-value1"))
	simulator.Done()

	blkAndPvtdata = prepareNextBlockForTest(t, txMgr, bg, "txid-2",
		map[string]string{"pubkey1": "pub-value2"}, map[string]string{"pvtkey2": "pvt-value2"})
	testutil.AssertNoError(t, txMgr.ValidateAndPrepare(blkAndPvtdata, true), "")
	testutil.AssertNoError(t, txMgr.Commit(), "")

	simulator, _ = txMgr.NewTxSimulator("tx-tmp")
	pvtval, _ = simulator.GetPrivateData("ns", "coll", "pvtkey1")
	testutil.AssertEquals(t, pvtval, []byte("pvt-value1"))
	simulator.Done()

	blkAndPvtdata = prepareNextBlockForTest(t, txMgr, bg, "txid-2",
		map[string]string{"pubkey1": "pub-value3"}, map[string]string{"pvtkey3": "pvt-value3"})
	testutil.AssertNoError(t, txMgr.ValidateAndPrepare(blkAndPvtdata, true), "")
	testutil.AssertNoError(t, txMgr.Commit(), "")

	simulator, _ = txMgr.NewTxSimulator("tx-tmp")
	pvtval, _ = simulator.GetPrivateData("ns", "coll", "pvtkey1")
	testutil.AssertNil(t, pvtval)
	simulator.Done()
}

func TestTxWithPubMetadata(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxwithpubmetadata"
		testEnv.init(t, testLedgerID, nil)
		testTxWithPubMetadata(t, testEnv)
		testEnv.cleanup()
	}
}

func testTxWithPubMetadata(t *testing.T, env testEnv) {
	namespace := "testns"
	txMgr := env.getTxMgr()
	txMgrHelper := newTxMgrTestHelper(t, txMgr)

	
	s1, _ := txMgr.NewTxSimulator("test_tx1")
	key1, value1, metadata1 := "key1", []byte("value1"), map[string][]byte{"entry1": []byte("meatadata1-entry1")}
	key2, value2, metadata2 := "key2", []byte("value2"), map[string][]byte{"entry1": []byte("meatadata2-entry1")}
	key3, metadata3 := "key3", map[string][]byte{"entry1": []byte("meatadata3-entry")}

	s1.SetState(namespace, key1, value1)
	s1.SetStateMetadata(namespace, key1, metadata1)
	s1.SetState(namespace, key2, value2)
	s1.SetStateMetadata(namespace, key2, metadata2)
	s1.SetStateMetadata(namespace, key3, metadata3)
	s1.Done()
	txRWSet1, _ := s1.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet1.PubSimulationResults)

	
	qe, _ := txMgr.NewQueryExecutor("test_tx2")
	checkTestQueryResults(t, qe, namespace, key1, value1, metadata1)
	checkTestQueryResults(t, qe, namespace, key2, value2, metadata2)
	checkTestQueryResults(t, qe, namespace, key3, nil, nil)
	qe.Done()

	
	updatedMetadata1 := map[string][]byte{"entry1": []byte("meatadata1-entry1"), "entry2": []byte("meatadata1-entry2")}
	s2, _ := txMgr.NewTxSimulator("test_tx3")
	s2.SetStateMetadata(namespace, key1, updatedMetadata1)
	s2.DeleteStateMetadata(namespace, key2)
	s2.Done()
	txRWSet2, _ := s2.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet2.PubSimulationResults)

	
	qe, _ = txMgr.NewQueryExecutor("test_tx4")
	checkTestQueryResults(t, qe, namespace, key1, value1, updatedMetadata1)
	checkTestQueryResults(t, qe, namespace, key2, value2, nil)
	qe.Done()
}

func TestTxWithPvtdataMetadata(t *testing.T) {
	ledgerid, ns, coll := "testtxwithpvtdatametadata", "ns", "coll"
	cs := btltestutil.NewMockCollectionStore()
	cs.SetBTL("ns", "coll", 1000)
	btlPolicy := pvtdatapolicy.ConstructBTLPolicy(cs)
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testEnv.init(t, ledgerid, btlPolicy)
		testTxWithPvtdataMetadata(t, testEnv, ns, coll)
		testEnv.cleanup()
	}
}

func testTxWithPvtdataMetadata(t *testing.T, env testEnv, ns, coll string) {
	ledgerid := "testtxwithpvtdatametadata"
	txMgr := env.getTxMgr()
	bg, _ := testutil.NewBlockGenerator(t, ledgerid, false)

	populateCollConfigForTest(t, txMgr.(*LockBasedTxMgr), []collConfigkey{{"ns", "coll"}}, version.NewHeight(1, 1))

	
	s1, _ := txMgr.NewTxSimulator("test_tx1")
	key1, value1, metadata1 := "key1", []byte("value1"), map[string][]byte{"entry1": []byte("meatadata1-entry1")}
	key2, value2, metadata2 := "key2", []byte("value2"), map[string][]byte{"entry1": []byte("meatadata2-entry1")}
	key3, metadata3 := "key3", map[string][]byte{"entry1": []byte("meatadata3-entry")}
	s1.SetPrivateData(ns, coll, key1, value1)
	s1.SetPrivateDataMetadata(ns, coll, key1, metadata1)
	s1.SetPrivateData(ns, coll, key2, value2)
	s1.SetPrivateDataMetadata(ns, coll, key2, metadata2)
	s1.SetPrivateDataMetadata(ns, coll, key3, metadata3)
	s1.Done()

	blkAndPvtdata1 := prepareNextBlockForTestFromSimulator(t, bg, s1)
	assert.NoError(t, txMgr.ValidateAndPrepare(blkAndPvtdata1, true))
	assert.NoError(t, txMgr.Commit())

	
	qe, _ := txMgr.NewQueryExecutor("test_tx2")
	checkPvtdataTestQueryResults(t, qe, ns, coll, key1, value1, metadata1)
	checkPvtdataTestQueryResults(t, qe, ns, coll, key2, value2, metadata2)
	checkPvtdataTestQueryResults(t, qe, ns, coll, key3, nil, nil)
	qe.Done()

	
	updatedMetadata1 := map[string][]byte{"entry1": []byte("meatadata1-entry1"), "entry2": []byte("meatadata1-entry2")}
	s2, _ := txMgr.NewTxSimulator("test_tx3")
	s2.SetPrivateDataMetadata(ns, coll, key1, updatedMetadata1)
	s2.DeletePrivateDataMetadata(ns, coll, key2)
	s2.Done()

	blkAndPvtdata2 := prepareNextBlockForTestFromSimulator(t, bg, s2)
	assert.NoError(t, txMgr.ValidateAndPrepare(blkAndPvtdata2, true))
	assert.NoError(t, txMgr.Commit())

	
	qe, _ = txMgr.NewQueryExecutor("test_tx4")
	checkPvtdataTestQueryResults(t, qe, ns, coll, key1, value1, updatedMetadata1)
	checkPvtdataTestQueryResults(t, qe, ns, coll, key2, value2, nil)
	qe.Done()
}

func prepareNextBlockForTest(t *testing.T, txMgr txmgr.TxMgr, bg *testutil.BlockGenerator,
	txid string, pubKVs map[string]string, pvtKVs map[string]string) *ledger.BlockAndPvtData {
	simulator, _ := txMgr.NewTxSimulator(txid)
	
	for k, v := range pubKVs {
		simulator.SetState("ns", k, []byte(v))
	}
	for k, v := range pvtKVs {
		simulator.SetPrivateData("ns", "coll", k, []byte(v))
	}
	simulator.Done()
	return prepareNextBlockForTestFromSimulator(t, bg, simulator)
}

func prepareNextBlockForTestFromSimulator(t *testing.T, bg *testutil.BlockGenerator, simulator ledger.TxSimulator) *ledger.BlockAndPvtData {
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimBytes, _ := simRes.GetPubSimulationBytes()
	block := bg.NextBlock([][]byte{pubSimBytes})
	return &ledger.BlockAndPvtData{Block: block,
		BlockPvtData: map[uint64]*ledger.TxPvtData{0: {SeqInBlock: 0, WriteSet: simRes.PvtSimulationResults}},
	}
}

func checkTestQueryResults(t *testing.T, qe ledger.QueryExecutor, ns, key string,
	expectedVal []byte, expectedMetadata map[string][]byte) {
	committedVal, err := qe.GetState(ns, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedVal, committedVal)

	committedMetadata, err := qe.GetStateMetadata(ns, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedMetadata, committedMetadata)
	t.Logf("key=%s, value=%s, metadata=%s", key, committedVal, committedMetadata)
}

func checkPvtdataTestQueryResults(t *testing.T, qe ledger.QueryExecutor, ns, coll, key string,
	expectedVal []byte, expectedMetadata map[string][]byte) {
	committedVal, err := qe.GetPrivateData(ns, coll, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedVal, committedVal)

	committedMetadata, err := qe.GetPrivateDataMetadata(ns, coll, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedMetadata, committedMetadata)
	t.Logf("key=%s, value=%s, metadata=%s", key, committedVal, committedMetadata)
}
