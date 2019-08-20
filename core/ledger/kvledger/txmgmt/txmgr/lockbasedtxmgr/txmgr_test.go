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
	"strings"
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	btltestutil "github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/testutil"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/stretchr/testify/assert"
)

func TestTxSimulatorWithNoExistingData(t *testing.T) {
	
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxsimulatorwithnoexistingdata"
		testEnv.init(t, testLedgerID, nil)
		testTxSimulatorWithNoExistingData(t, testEnv)
	}
}

func testTxSimulatorWithNoExistingData(t *testing.T, env testEnv) {
	txMgr := env.getTxMgr()
	s, _ := txMgr.NewTxSimulator("test_txid")
	value, err := s.GetState("ns1", "key1")
	assert.NoErrorf(t, err, "Error in GetState(): %s", err)
	assert.Nil(t, value)

	s.SetState("ns1", "key1", []byte("value1"))
	s.SetState("ns1", "key2", []byte("value2"))
	s.SetState("ns2", "key3", []byte("value3"))
	s.SetState("ns2", "key4", []byte("value4"))

	value, _ = s.GetState("ns2", "key3")
	assert.Nil(t, value)

	simulationResults, err := s.GetTxSimulationResults()
	assert.NoError(t, err)
	assert.Nil(t, simulationResults.PvtSimulationResults)
}

func TestTxSimulatorGetResults(t *testing.T) {
	testEnv := testEnvsMap[levelDBtestEnvName]
	testEnv.init(t, "testLedger", nil)
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
	assert.Len(t, simulationResults1.PubSimulationResults.NsRwset, 1)
	
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
	assert.Equal(t, []byte("value1"), value)
	s2.SetState("ns1", "key1", []byte("value1_1"))
	s2.DeleteState("ns2", "key3")
	value, _ = s2.GetState("ns1", "key1")
	assert.Equal(t, []byte("value1"), value)
	s2.Done()
	
	txRWSet2, _ := s2.GetTxSimulationResults()
	txMgrHelper.validateAndCommitRWSet(txRWSet2.PubSimulationResults)

	
	s3, _ := txMgr.NewTxSimulator("test_tx3")
	value, _ = s3.GetState("ns1", "key1")
	assert.Equal(t, []byte("value1_1"), value)
	value, _ = s3.GetState("ns2", "key3")
	assert.Nil(t, value)
	s3.Done()

	
	vv, _ := env.getVDB().GetState("ns1", "key1")
	assert.Equal(t, version.NewHeight(2, 0), vv.Version)
	vv, _ = env.getVDB().GetState("ns1", "key2")
	assert.Equal(t, version.NewHeight(1, 0), vv.Version)
}

func TestTxValidation(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxvalidation"
		testEnv.init(t, testLedgerID, nil)
		testTxValidation(t, testEnv)
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
	assert.Equal(t, []byte("value1"), value)

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

		testLedgerID = "testiterator.2"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 1, 11)

		testLedgerID = "testiterator.3"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 0, 0)

		testLedgerID = "testiterator.4"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 5, 0)

		testLedgerID = "testiterator.5"
		testEnv.init(t, testLedgerID, nil)
		testIterator(t, testEnv, 10, 0, 5)
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
		assert.Equal(t, createTestKey(keyNum), k)
		assert.Equal(t, createTestValue(keyNum), v)
		count++
	}
	assert.Equal(t, expectedCount, count)
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
		assert.NoError(t, err, "An unexpected error was thrown during iterator Next()")
		vkv := queryResult.(*queryresult.KV)
		key := vkv.Key
		assert.Equal(t, expectedKey, key)
	}
	queryResult, err := itr.Next()
	assert.NoError(t, err, "An unexpected error was thrown during iterator Next()")
	assert.Nil(t, queryResult)
}

func TestIteratorWithDeletes(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testiteratorwithdeletes"
		testEnv.init(t, testLedgerID, nil)
		testIteratorWithDeletes(t, testEnv)
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
	assert.Equal(t, createTestKey(3), kv.(*queryresult.KV).Key)
	kv, _ = itr.Next()
	assert.Equal(t, createTestKey(5), kv.(*queryresult.KV).Key)
}

func TestTxValidationWithItr(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxvalidationwithitr"
		testEnv.init(t, testLedgerID, nil)
		testTxValidationWithItr(t, testEnv)
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
	assert.Len(t, values, 10)
	for i, v := range values {
		assert.Equal(t, multipleKeyMap[multipleKeys[i]], v)
	}

	s2, _ := txMgr.NewTxSimulator("test_tx3")
	defer s2.Done()
	values, _ = s2.GetStateMultipleKeys(cID, multipleKeys[5:7])
	assert.Len(t, values, 2)
	for i, v := range values {
		assert.Equal(t, multipleKeyMap[multipleKeys[i+5]], v)
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
	assert.NoError(t, err, "Error upon ExecuteQuery()")
	counter := 0
	for {
		queryRecord, _ := itr.Next()
		if queryRecord == nil {
			break
		}
		
		assetResp := &Asset{}
		json.Unmarshal(queryRecord.(*queryresult.KV).Value, &assetResp)
		
		assert.Equal(t, "bob", assetResp.Owner)
		counter++
	}
	
	assert.Equal(t, 3, counter)
}


func TestExecutePaginatedQuery(t *testing.T) {

	for _, testEnv := range testEnvs {
		
		if testEnv.getName() == couchDBtestEnvName {
			t.Logf("Running test for TestEnv = %s", testEnv.getName())
			testLedgerID := "testexecutepaginatedquery"
			testEnv.init(t, testLedgerID, nil)
			testExecutePaginatedQuery(t, testEnv)
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
	assert.NoError(t, err, "Error upon ExecuteQueryWithMetadata()")
	counter := 0
	for {
		queryRecord, _ := itr.Next()
		if queryRecord == nil {
			break
		}
		
		assetResp := &Asset{}
		json.Unmarshal(queryRecord.(*queryresult.KV).Value, &assetResp)
		
		assert.Equal(t, "bob", assetResp.Owner)
		counter++
	}
	
	assert.Equal(t, 2, counter)

	bookmark := itr.GetBookmarkAndClose()

	queryOptions = map[string]interface{}{
		"limit": int32(2),
	}
	if bookmark != "" {
		queryOptions["bookmark"] = bookmark
	}

	itr, err = queryExecuter.ExecuteQueryWithMetadata("ns1", queryString, queryOptions)
	assert.NoError(t, err, "Error upon ExecuteQuery()")
	counter = 0
	for {
		queryRecord, _ := itr.Next()
		if queryRecord == nil {
			break
		}
		
		assetResp := &Asset{}
		json.Unmarshal(queryRecord.(*queryresult.KV).Value, &assetResp)
		
		assert.Equal(t, "bob", assetResp.Owner)
		counter++
	}
	
	assert.Equal(t, 1, counter)
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
			assert.NoError(t, err)
		}
		if testEnv.getName() == couchDBtestEnvName {
			assert.Error(t, err)
		}
	}
}



func TestTxSimulatorUnsupportedTx(t *testing.T) {
	testEnv := testEnvs[0]
	testEnv.init(t, "testtxsimulatorunsupportedtx", nil)
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
	assert.NoError(t, err)
	_, err = simulator.GetPrivateDataRangeScanIterator("ns1", "coll1", "startKey", "endKey")
	_, ok := err.(*txmgr.ErrUnsupportedTransaction)
	assert.True(t, ok)

	simulator, _ = txMgr.NewTxSimulator("txid2")
	_, err = simulator.GetPrivateDataRangeScanIterator("ns1", "coll1", "startKey", "endKey")
	assert.NoError(t, err)
	err = simulator.SetState("ns", "key", []byte("value"))
	_, ok = err.(*txmgr.ErrUnsupportedTransaction)
	assert.True(t, ok)

	queryOptions := map[string]interface{}{
		"limit": int32(2),
	}

	simulator, _ = txMgr.NewTxSimulator("txid3")
	err = simulator.SetState("ns", "key", []byte("value"))
	assert.NoError(t, err)
	_, err = simulator.GetStateRangeScanIteratorWithMetadata("ns1", "startKey", "endKey", queryOptions)
	_, ok = err.(*txmgr.ErrUnsupportedTransaction)
	assert.True(t, ok)

	simulator, _ = txMgr.NewTxSimulator("txid4")
	_, err = simulator.GetStateRangeScanIteratorWithMetadata("ns1", "startKey", "endKey", queryOptions)
	assert.NoError(t, err)
	err = simulator.SetState("ns", "key", []byte("value"))
	_, ok = err.(*txmgr.ErrUnsupportedTransaction)
	assert.True(t, ok)

}


func TestTxSimulatorQueryUnsupportedTx(t *testing.T) {

	for _, testEnv := range testEnvs {
		
		if testEnv.getName() == couchDBtestEnvName {
			t.Logf("Running test for TestEnv = %s", testEnv.getName())
			testLedgerID := "testtxsimulatorunsupportedtxqueries"
			testEnv.init(t, testLedgerID, nil)
			testTxSimulatorQueryUnsupportedTx(t, testEnv)
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
	assert.NoError(t, err)
	_, err = simulator.ExecuteQueryWithMetadata("ns1", queryString, queryOptions)
	_, ok := err.(*txmgr.ErrUnsupportedTransaction)
	assert.True(t, ok)

	simulator, _ = txMgr.NewTxSimulator("txid2")
	_, err = simulator.ExecuteQueryWithMetadata("ns1", queryString, queryOptions)
	assert.NoError(t, err)
	err = simulator.SetState("ns1", "key1", []byte(`{"asset_name":"marble1","color":"red","size":"25","owner":"jerry"}`))
	_, ok = err.(*txmgr.ErrUnsupportedTransaction)
	assert.True(t, ok)

}

func TestConstructUniquePvtData(t *testing.T) {
	v1 := []byte{1}
	
	pvtDataBlk1Tx1 := producePvtdata(t, 1, []string{"ns1:coll1"}, []string{"key1"}, [][]byte{v1})
	
	pvtDataBlk1Tx2 := producePvtdata(t, 2, []string{"ns1:coll1", "ns1:coll2"}, []string{"key2", "key3"}, [][]byte{v1, v1})
	
	pvtDataBlk1Tx3 := producePvtdata(t, 3, []string{"ns1:coll2"}, []string{"key4"}, [][]byte{v1})

	v2 := []byte{2}
	
	pvtDataBlk2Tx1 := producePvtdata(t, 1, []string{"ns1:coll1"}, []string{"key1"}, [][]byte{v2})
	
	pvtDataBlk2Tx2 := producePvtdata(t, 2, []string{"ns1:coll1"}, []string{"key2"}, [][]byte{nil})

	v3 := []byte{3}
	
	pvtDataBlk3Tx1 := producePvtdata(t, 1, []string{"ns1:coll1"}, []string{"key1"}, [][]byte{v3})

	blocksPvtData := map[uint64][]*ledger.TxPvtData{
		1: {
			pvtDataBlk1Tx1,
			pvtDataBlk1Tx2,
			pvtDataBlk1Tx3,
		},
		2: {
			pvtDataBlk2Tx1,
			pvtDataBlk2Tx2,
		},
		3: {
			pvtDataBlk3Tx1,
		},
	}

	hashedCompositeKeyNs1Coll2Key3 := privacyenabledstate.HashedCompositeKey{Namespace: "ns1", CollectionName: "coll2", KeyHash: string(util.ComputeStringHash("key3"))}
	pvtKVWriteNs1Coll2Key3 := &privacyenabledstate.PvtKVWrite{Key: "key3", IsDelete: false, Value: v1, Version: version.NewHeight(1, 2)}

	hashedCompositeKeyNs1Coll2Key4 := privacyenabledstate.HashedCompositeKey{Namespace: "ns1", CollectionName: "coll2", KeyHash: string(util.ComputeStringHash("key4"))}
	pvtKVWriteNs1Coll2Key4 := &privacyenabledstate.PvtKVWrite{Key: "key4", IsDelete: false, Value: v1, Version: version.NewHeight(1, 3)}

	hashedCompositeKeyNs1Coll1Key2 := privacyenabledstate.HashedCompositeKey{Namespace: "ns1", CollectionName: "coll1", KeyHash: string(util.ComputeStringHash("key2"))}
	pvtKVWriteNs1Coll1Key2 := &privacyenabledstate.PvtKVWrite{Key: "key2", IsDelete: true, Value: nil, Version: version.NewHeight(2, 2)}

	hashedCompositeKeyNs1Coll1Key1 := privacyenabledstate.HashedCompositeKey{Namespace: "ns1", CollectionName: "coll1", KeyHash: string(util.ComputeStringHash("key1"))}
	pvtKVWriteNs1Coll1Key1 := &privacyenabledstate.PvtKVWrite{Key: "key1", IsDelete: false, Value: v3, Version: version.NewHeight(3, 1)}

	expectedUniquePvtData := uniquePvtDataMap{
		hashedCompositeKeyNs1Coll2Key3: pvtKVWriteNs1Coll2Key3,
		hashedCompositeKeyNs1Coll2Key4: pvtKVWriteNs1Coll2Key4,
		hashedCompositeKeyNs1Coll1Key2: pvtKVWriteNs1Coll1Key2,
		hashedCompositeKeyNs1Coll1Key1: pvtKVWriteNs1Coll1Key1,
	}

	uniquePvtData, err := constructUniquePvtData(blocksPvtData)
	assert.NoError(t, err)
	assert.Equal(t, expectedUniquePvtData, uniquePvtData)
}

func TestFindAndRemoveStalePvtData(t *testing.T) {
	ledgerid := "TestFindAndRemoveStalePvtData"
	testEnv := testEnvs[0]
	testEnv.init(t, ledgerid, nil)
	db := testEnv.getVDB()

	batch := privacyenabledstate.NewUpdateBatch()
	batch.HashUpdates.Put("ns1", "coll1", util.ComputeStringHash("key1"), util.ComputeStringHash("value_1_1_1"), version.NewHeight(1, 1))
	batch.HashUpdates.Put("ns1", "coll2", util.ComputeStringHash("key2"), util.ComputeStringHash("value_1_2_2"), version.NewHeight(1, 2))
	batch.HashUpdates.Put("ns2", "coll1", util.ComputeStringHash("key2"), util.ComputeStringHash("value_2_1_2"), version.NewHeight(2, 1))
	batch.HashUpdates.Put("ns2", "coll2", util.ComputeStringHash("key3"), util.ComputeStringHash("value_2_2_3"), version.NewHeight(10, 10))

	
	db.ApplyPrivacyAwareUpdates(batch, version.NewHeight(11, 1))

	
	
	
	hashedCompositeKeyNs1Coll1Key1 := privacyenabledstate.HashedCompositeKey{Namespace: "ns1", CollectionName: "coll1", KeyHash: string(util.ComputeStringHash("key1"))}
	pvtKVWriteNs1Coll1Key1 := &privacyenabledstate.PvtKVWrite{Key: "key1", IsDelete: false, Value: []byte("old_value_1_1_1"), Version: version.NewHeight(1, 0)}

	
	hashedCompositeKeyNs2Coll1Key2 := privacyenabledstate.HashedCompositeKey{Namespace: "ns2", CollectionName: "coll1", KeyHash: string(util.ComputeStringHash("key2"))}
	pvtKVWriteNs2Coll1Key2 := &privacyenabledstate.PvtKVWrite{Key: "key2", IsDelete: false, Value: []byte("value_2_1_2"), Version: version.NewHeight(2, 1)}

	
	hashedCompositeKeyNs1Coll3Key3 := privacyenabledstate.HashedCompositeKey{Namespace: "ns1", CollectionName: "coll3", KeyHash: string(util.ComputeStringHash("key3"))}
	pvtKVWriteNs1Coll3Key3 := &privacyenabledstate.PvtKVWrite{Key: "key3", IsDelete: true, Value: nil, Version: version.NewHeight(2, 3)}

	
	
	hashedCompositeKeyNs2Coll2Key3 := privacyenabledstate.HashedCompositeKey{Namespace: "ns2", CollectionName: "coll2", KeyHash: string(util.ComputeStringHash("key3"))}
	pvtKVWriteNs2Coll2Key3 := &privacyenabledstate.PvtKVWrite{Key: "key3", IsDelete: false, Value: []byte("value_2_2_3"), Version: version.NewHeight(9, 9)}

	uniquePvtData := uniquePvtDataMap{
		hashedCompositeKeyNs1Coll1Key1: pvtKVWriteNs1Coll1Key1,
		hashedCompositeKeyNs2Coll1Key2: pvtKVWriteNs2Coll1Key2,
		hashedCompositeKeyNs1Coll3Key3: pvtKVWriteNs1Coll3Key3,
		hashedCompositeKeyNs2Coll2Key3: pvtKVWriteNs2Coll2Key3,
	}

	
	expectedBatch := privacyenabledstate.NewUpdateBatch()
	expectedBatch.PvtUpdates.Put("ns2", "coll1", "key2", []byte("value_2_1_2"), version.NewHeight(2, 1))
	expectedBatch.PvtUpdates.Delete("ns1", "coll3", "key3", version.NewHeight(2, 3))
	expectedBatch.PvtUpdates.Put("ns2", "coll2", "key3", []byte("value_2_2_3"), version.NewHeight(10, 10))

	err := uniquePvtData.findAndRemoveStalePvtData(db)
	assert.NoError(t, err, "uniquePvtData.findAndRemoveStatePvtData resulted in an error")
	batch = uniquePvtData.transformToUpdateBatch()
	assert.Equal(t, expectedBatch.PvtUpdates, batch.PvtUpdates)
}

func producePvtdata(t *testing.T, txNum uint64, nsColls []string, keys []string, values [][]byte) *ledger.TxPvtData {
	builder := rwsetutil.NewRWSetBuilder()
	for index, nsColl := range nsColls {
		nsCollSplit := strings.Split(nsColl, ":")
		ns := nsCollSplit[0]
		coll := nsCollSplit[1]
		key := keys[index]
		value := values[index]
		builder.AddToPvtAndHashedWriteSet(ns, coll, key, value)
	}
	simRes, err := builder.GetTxSimulationResults()
	assert.NoError(t, err)
	return &ledger.TxPvtData{
		SeqInBlock: txNum,
		WriteSet:   simRes.PvtSimulationResults,
	}
}

func TestRemoveStaleAndCommitPvtDataOfOldBlocks(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testValidationAndCommitOfOldPvtData(t, testEnv)
	}
}

func testValidationAndCommitOfOldPvtData(t *testing.T, env testEnv) {
	ledgerid := "testvalidationandcommitofoldpvtdata"
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns1", "coll1"}: 0,
			{"ns1", "coll2"}: 0,
		},
	)
	env.init(t, ledgerid, btlPolicy)
	txMgr := env.getTxMgr()
	populateCollConfigForTest(t, txMgr.(*LockBasedTxMgr),
		[]collConfigkey{
			{"ns1", "coll1"},
			{"ns1", "coll2"},
		},
		version.NewHeight(1, 1),
	)

	db := env.getVDB()
	updateBatch := privacyenabledstate.NewUpdateBatch()
	
	updateBatch.HashUpdates.Put("ns1", "coll1", util.ComputeStringHash("key1"), util.ComputeStringHash("value1"), version.NewHeight(1, 1)) 
	updateBatch.HashUpdates.Put("ns1", "coll1", util.ComputeStringHash("key2"), util.ComputeStringHash("value2"), version.NewHeight(1, 2)) 
	updateBatch.HashUpdates.Put("ns1", "coll2", util.ComputeStringHash("key3"), util.ComputeStringHash("value3"), version.NewHeight(1, 2)) 
	updateBatch.HashUpdates.Put("ns1", "coll2", util.ComputeStringHash("key4"), util.ComputeStringHash("value4"), version.NewHeight(1, 3)) 
	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(1, 2))

	updateBatch = privacyenabledstate.NewUpdateBatch()
	updateBatch.HashUpdates.Put("ns1", "coll1", util.ComputeStringHash("key1"), util.ComputeStringHash("new-value1"), version.NewHeight(2, 1)) 
	updateBatch.HashUpdates.Delete("ns1", "coll1", util.ComputeStringHash("key2"), version.NewHeight(2, 2))                                    
	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(2, 2))

	updateBatch = privacyenabledstate.NewUpdateBatch()
	updateBatch.HashUpdates.Put("ns1", "coll1", util.ComputeStringHash("key1"), util.ComputeStringHash("another-new-value1"), version.NewHeight(3, 1)) 
	updateBatch.HashUpdates.Put("ns1", "coll2", util.ComputeStringHash("key3"), util.ComputeStringHash("value3"), version.NewHeight(3, 2))             
	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(3, 2))

	v1 := []byte("value1")
	
	pvtDataBlk1Tx1 := producePvtdata(t, 1, []string{"ns1:coll1"}, []string{"key1"}, [][]byte{v1})
	
	
	v2 := []byte("value2")
	v3 := []byte("value3")
	pvtDataBlk1Tx2 := producePvtdata(t, 2, []string{"ns1:coll1", "ns1:coll2"}, []string{"key2", "key3"}, [][]byte{v2, v3})
	
	v4 := []byte("value4")
	pvtDataBlk1Tx3 := producePvtdata(t, 3, []string{"ns1:coll2"}, []string{"key4"}, [][]byte{v4})

	nv1 := []byte("new-value1")
	
	pvtDataBlk2Tx1 := producePvtdata(t, 1, []string{"ns1:coll1"}, []string{"key1"}, [][]byte{nv1})
	
	pvtDataBlk2Tx2 := producePvtdata(t, 2, []string{"ns1:coll1"}, []string{"key2"}, [][]byte{nil})

	anv1 := []byte("another-new-value1")
	
	pvtDataBlk3Tx1 := producePvtdata(t, 1, []string{"ns1:coll1"}, []string{"key1"}, [][]byte{anv1})
	
	pvtDataBlk3Tx2 := producePvtdata(t, 2, []string{"ns1:coll2"}, []string{"key3"}, [][]byte{v3})

	blocksPvtData := map[uint64][]*ledger.TxPvtData{
		1: {
			pvtDataBlk1Tx1,
			pvtDataBlk1Tx2,
			pvtDataBlk1Tx3,
		},
		2: {
			pvtDataBlk2Tx1,
			pvtDataBlk2Tx2,
		},
		3: {
			pvtDataBlk3Tx1,
			pvtDataBlk3Tx2,
		},
	}

	err := txMgr.RemoveStaleAndCommitPvtDataOfOldBlocks(blocksPvtData)
	assert.NoError(t, err)

	vv, err := db.GetPrivateData("ns1", "coll1", "key1")
	assert.NoError(t, err)
	assert.Equal(t, anv1, vv.Value) 

	vv, err = db.GetPrivateData("ns1", "coll1", "key2")
	assert.NoError(t, err)
	assert.Equal(t, nil, nil) 

	vv, err = db.GetPrivateData("ns1", "coll2", "key3")
	assert.NoError(t, err)
	assert.Equal(t, v3, vv.Value)
	assert.Equal(t, version.NewHeight(3, 2), vv.Version) 

	vv, err = db.GetPrivateData("ns1", "coll2", "key4")
	assert.NoError(t, err)
	assert.Equal(t, v4, vv.Value)
}

func TestTxSimulatorMissingPvtdata(t *testing.T) {
	testEnv := testEnvs[0]
	testEnv.init(t, "TestTxSimulatorUnsupportedTxQueries", nil)

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

	assert.True(t, testPvtValueEqual(t, txMgr, "ns1", "coll1", "key1", []byte("value1")))

	updateBatch = privacyenabledstate.NewUpdateBatch()
	updateBatch.HashUpdates.Put("ns1", "coll1", util.ComputeStringHash("key1"), util.ComputeStringHash("value1"), version.NewHeight(2, 1))
	updateBatch.HashUpdates.Put("ns1", "coll2", util.ComputeStringHash("key2"), util.ComputeStringHash("value2"), version.NewHeight(2, 1))
	updateBatch.HashUpdates.Put("ns1", "coll3", util.ComputeStringHash("key3"), util.ComputeStringHash("value3"), version.NewHeight(2, 1))
	updateBatch.PvtUpdates.Put("ns1", "coll3", "key3", []byte("value3"), version.NewHeight(2, 1))
	db.ApplyPrivacyAwareUpdates(updateBatch, version.NewHeight(2, 1))

	assert.False(t, testPvtKeyExist(t, txMgr, "ns1", "coll1", "key1"))

	assert.False(t, testPvtKeyExist(t, txMgr, "ns1", "coll2", "key2"))

	assert.True(t, testPvtValueEqual(t, txMgr, "ns1", "coll3", "key3", []byte("value3")))

	assert.True(t, testPvtValueEqual(t, txMgr, "ns1", "coll4", "key4", nil))
}

func TestRemoveStaleAndCommitPvtDataOfOldBlocksWithExpiry(t *testing.T) {
	ledgerid := "TestTxSimulatorMissingPvtdataExpiry"
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns", "coll"}: 1,
		},
	)
	testEnv := testEnvs[0]
	testEnv.init(t, ledgerid, btlPolicy)

	txMgr := testEnv.getTxMgr()
	populateCollConfigForTest(t, txMgr.(*LockBasedTxMgr),
		[]collConfigkey{
			{"ns", "coll"},
		},
		version.NewHeight(1, 1),
	)

	bg, _ := testutil.NewBlockGenerator(t, ledgerid, false)

	
	
	blkAndPvtdata := prepareNextBlockForTest(t, txMgr, bg, "txid-1",
		map[string]string{"pubkey1": "pub-value1"}, map[string]string{"pvtkey1": "pvt-value1"}, true)
	_, _, err := txMgr.ValidateAndPrepare(blkAndPvtdata, true)
	assert.NoError(t, err)
	
	assert.NoError(t, txMgr.Commit())

	
	assert.False(t, testPvtKeyExist(t, txMgr, "ns", "coll", "pvtkey1"))

	
	v1 := []byte("pvt-value1")
	pvtDataBlk1Tx1 := producePvtdata(t, 1, []string{"ns:coll"}, []string{"pvtkey1"}, [][]byte{v1})
	blocksPvtData := map[uint64][]*ledger.TxPvtData{
		1: {
			pvtDataBlk1Tx1,
		},
	}
	err = txMgr.RemoveStaleAndCommitPvtDataOfOldBlocks(blocksPvtData)
	assert.NoError(t, err)

	
	assert.True(t, testPvtValueEqual(t, txMgr, "ns", "coll", "pvtkey1", v1))

	
	
	blkAndPvtdata = prepareNextBlockForTest(t, txMgr, bg, "txid-2",
		map[string]string{"pubkey2": "pub-value2"}, map[string]string{"pvtkey2": "pvt-value2"}, true)
	_, _, err = txMgr.ValidateAndPrepare(blkAndPvtdata, true)
	assert.NoError(t, err)
	
	assert.NoError(t, txMgr.Commit())

	
	assert.False(t, testPvtKeyExist(t, txMgr, "ns", "coll", "pvtkey2"))

	blkAndPvtdata = prepareNextBlockForTest(t, txMgr, bg, "txid-3",
		map[string]string{"pubkey3": "pub-value3"}, nil, false)
	_, _, err = txMgr.ValidateAndPrepare(blkAndPvtdata, true)
	assert.NoError(t, err)
	
	assert.NoError(t, txMgr.Commit())

	
	

	
	v2 := []byte("pvt-value2")
	pvtDataBlk2Tx1 := producePvtdata(t, 1, []string{"ns:coll"}, []string{"pvtkey2"}, [][]byte{v2})
	blocksPvtData = map[uint64][]*ledger.TxPvtData{
		2: {
			pvtDataBlk2Tx1,
		},
	}

	err = txMgr.RemoveStaleAndCommitPvtDataOfOldBlocks(blocksPvtData)
	assert.NoError(t, err)

	
	assert.True(t, testPvtValueEqual(t, txMgr, "ns", "coll", "pvtkey2", v2))

	blkAndPvtdata = prepareNextBlockForTest(t, txMgr, bg, "txid-4",
		map[string]string{"pubkey4": "pub-value4"}, nil, false)
	_, _, err = txMgr.ValidateAndPrepare(blkAndPvtdata, true)
	assert.NoError(t, err)
	
	assert.NoError(t, txMgr.Commit())

	assert.True(t, testPvtValueEqual(t, txMgr, "ns", "coll", "pvtkey2", nil))
}

func testPvtKeyExist(t *testing.T, txMgr txmgr.TxMgr, ns, coll, key string) bool {
	simulator, _ := txMgr.NewTxSimulator("tx-tmp")
	defer simulator.Done()
	_, err := simulator.GetPrivateData(ns, coll, key)
	_, ok := err.(*txmgr.ErrPvtdataNotAvailable)
	return !ok
}

func testPvtValueEqual(t *testing.T, txMgr txmgr.TxMgr, ns, coll, key string, value []byte) bool {
	simulator, _ := txMgr.NewTxSimulator("tx-tmp")
	defer simulator.Done()
	pvtValue, err := simulator.GetPrivateData(ns, coll, key)
	assert.NoError(t, err)
	if bytes.Equal(pvtValue, value) {
		return true
	}
	return false
}

func TestDeleteOnCursor(t *testing.T) {
	cID := "cid"
	env := testEnvs[0]
	env.init(t, "TestDeleteOnCursor", nil)

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
		assert.NoError(t, err)
		assert.NotNil(t, kv)
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
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	key := kv.(*queryresult.KV).Key
	assert.Equal(t, "key_005", key)
	itr3.Close()
	s3.Done()
}

func TestTxSimulatorMissingPvtdataExpiry(t *testing.T) {
	ledgerid := "TestTxSimulatorMissingPvtdataExpiry"
	testEnv := testEnvs[0]
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns", "coll"}: 1,
		},
	)
	testEnv.init(t, ledgerid, btlPolicy)

	txMgr := testEnv.getTxMgr()
	populateCollConfigForTest(t, txMgr.(*LockBasedTxMgr), []collConfigkey{{"ns", "coll"}}, version.NewHeight(1, 1))

	bg, _ := testutil.NewBlockGenerator(t, ledgerid, false)

	blkAndPvtdata := prepareNextBlockForTest(t, txMgr, bg, "txid-1",
		map[string]string{"pubkey1": "pub-value1"}, map[string]string{"pvtkey1": "pvt-value1"}, false)
	_, _, err := txMgr.ValidateAndPrepare(blkAndPvtdata, true)
	assert.NoError(t, err)
	assert.NoError(t, txMgr.Commit())

	assert.True(t, testPvtValueEqual(t, txMgr, "ns", "coll", "pvtkey1", []byte("pvt-value1")))

	blkAndPvtdata = prepareNextBlockForTest(t, txMgr, bg, "txid-2",

		map[string]string{"pubkey1": "pub-value2"}, map[string]string{"pvtkey2": "pvt-value2"}, false)
	_, _, err = txMgr.ValidateAndPrepare(blkAndPvtdata, true)
	assert.NoError(t, err)
	assert.NoError(t, txMgr.Commit())

	assert.True(t, testPvtValueEqual(t, txMgr, "ns", "coll", "pvtkey1", []byte("pvt-value1")))

	blkAndPvtdata = prepareNextBlockForTest(t, txMgr, bg, "txid-2",
		map[string]string{"pubkey1": "pub-value3"}, map[string]string{"pvtkey3": "pvt-value3"}, false)
	_, _, err = txMgr.ValidateAndPrepare(blkAndPvtdata, true)
	assert.NoError(t, err)
	assert.NoError(t, txMgr.Commit())

	assert.True(t, testPvtValueEqual(t, txMgr, "ns", "coll", "pvtkey1", nil))
}

func TestTxWithPubMetadata(t *testing.T) {
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testLedgerID := "testtxwithpubmetadata"
		testEnv.init(t, testLedgerID, nil)
		testTxWithPubMetadata(t, testEnv)
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
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns", "coll"}: 1000,
		},
	)
	for _, testEnv := range testEnvs {
		t.Logf("Running test for TestEnv = %s", testEnv.getName())
		testEnv.init(t, ledgerid, btlPolicy)
		testTxWithPvtdataMetadata(t, testEnv, ns, coll)
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
	_, _, err := txMgr.ValidateAndPrepare(blkAndPvtdata1, true)
	assert.NoError(t, err)
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
	_, _, err = txMgr.ValidateAndPrepare(blkAndPvtdata2, true)
	assert.NoError(t, err)
	assert.NoError(t, txMgr.Commit())

	
	qe, _ = txMgr.NewQueryExecutor("test_tx4")
	checkPvtdataTestQueryResults(t, qe, ns, coll, key1, value1, updatedMetadata1)
	checkPvtdataTestQueryResults(t, qe, ns, coll, key2, value2, nil)
	qe.Done()
}

func prepareNextBlockForTest(t *testing.T, txMgr txmgr.TxMgr, bg *testutil.BlockGenerator,
	txid string, pubKVs map[string]string, pvtKVs map[string]string, isMissing bool) *ledger.BlockAndPvtData {
	simulator, _ := txMgr.NewTxSimulator(txid)
	
	for k, v := range pubKVs {
		simulator.SetState("ns", k, []byte(v))
	}
	for k, v := range pvtKVs {
		simulator.SetPrivateData("ns", "coll", k, []byte(v))
	}
	simulator.Done()
	if isMissing {
		return prepareNextBlockForTestFromSimulatorWithMissingData(t, bg, simulator, txid, 1, "ns", "coll", true)
	}
	return prepareNextBlockForTestFromSimulator(t, bg, simulator)
}

func prepareNextBlockForTestFromSimulator(t *testing.T, bg *testutil.BlockGenerator, simulator ledger.TxSimulator) *ledger.BlockAndPvtData {
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimBytes, _ := simRes.GetPubSimulationBytes()
	block := bg.NextBlock([][]byte{pubSimBytes})
	return &ledger.BlockAndPvtData{Block: block,
		PvtData: ledger.TxPvtDataMap{0: {SeqInBlock: 0, WriteSet: simRes.PvtSimulationResults}},
	}
}

func prepareNextBlockForTestFromSimulatorWithMissingData(t *testing.T, bg *testutil.BlockGenerator, simulator ledger.TxSimulator,
	txid string, txNum uint64, ns, coll string, isEligible bool) *ledger.BlockAndPvtData {
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimBytes, _ := simRes.GetPubSimulationBytes()
	block := bg.NextBlock([][]byte{pubSimBytes})
	missingData := make(ledger.TxMissingPvtDataMap)
	missingData.Add(txNum, ns, coll, isEligible)
	return &ledger.BlockAndPvtData{Block: block, MissingPvtData: missingData}
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

func TestName(t *testing.T) {
	testEnv := testEnvsMap[levelDBtestEnvName]
	testEnv.init(t, "testLedger", nil)
	defer testEnv.cleanup()
	txMgr := testEnv.getTxMgr()
	assert.Equal(t, "state", txMgr.Name())
}
