/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kvledger

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage/fsblkstorage"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/util"
	ledgerproto "github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	putils "github.com/mcc-github/blockchain/protos/utils"
	"github.com/spf13/viper"
)

func TestLedgerProvider(t *testing.T) {
	env := newTestEnv(t)
	defer env.cleanup()
	numLedgers := 10
	provider, _ := NewProvider()
	existingLedgerIDs, err := provider.List()
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, len(existingLedgerIDs), 0)
	for i := 0; i < numLedgers; i++ {
		genesisBlock, _ := configtxtest.MakeGenesisBlock(constructTestLedgerID(i))
		provider.Create(genesisBlock)
	}
	existingLedgerIDs, err = provider.List()
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, len(existingLedgerIDs), numLedgers)

	provider.Close()

	provider, _ = NewProvider()
	defer provider.Close()
	ledgerIds, _ := provider.List()
	testutil.AssertEquals(t, len(ledgerIds), numLedgers)
	t.Logf("ledgerIDs=%#v", ledgerIds)
	for i := 0; i < numLedgers; i++ {
		testutil.AssertEquals(t, ledgerIds[i], constructTestLedgerID(i))
	}
	for i := 0; i < numLedgers; i++ {
		status, _ := provider.Exists(constructTestLedgerID(i))
		testutil.AssertEquals(t, status, true)
		ledger, err := provider.Open(constructTestLedgerID(i))
		testutil.AssertNoError(t, err, "")
		bcInfo, err := ledger.GetBlockchainInfo()
		ledger.Close()
		testutil.AssertNoError(t, err, "")
		testutil.AssertEquals(t, bcInfo.Height, uint64(1))
	}
	gb, _ := configtxtest.MakeGenesisBlock(constructTestLedgerID(2))
	_, err = provider.Create(gb)
	testutil.AssertEquals(t, err, ErrLedgerIDExists)

	status, err := provider.Exists(constructTestLedgerID(numLedgers))
	testutil.AssertNoError(t, err, "Failed to check for ledger existence")
	testutil.AssertEquals(t, false, status)

	_, err = provider.Open(constructTestLedgerID(numLedgers))
	testutil.AssertEquals(t, err, ErrNonExistingLedgerID)
}

func TestRecovery(t *testing.T) {
	env := newTestEnv(t)
	defer env.cleanup()
	provider, _ := NewProvider()

	
	genesisBlock, _ := configtxtest.MakeGenesisBlock(constructTestLedgerID(1))
	ledger, err := provider.(*Provider).openInternal(constructTestLedgerID(1))
	ledger.CommitWithPvtData(&ledgerproto.BlockAndPvtData{Block: genesisBlock})
	ledger.Close()

	
	
	provider.(*Provider).idStore.setUnderConstructionFlag(constructTestLedgerID(1))
	provider.Close()

	
	provider, err = NewProvider()
	testutil.AssertNoError(t, err, "Provider failed to recover an underConstructionLedger")
	
	flag, err := provider.(*Provider).idStore.getUnderConstructionFlag()
	testutil.AssertNoError(t, err, "Failed to read the underconstruction flag")
	testutil.AssertEquals(t, flag, "")
	ledger, err = provider.Open(constructTestLedgerID(1))
	testutil.AssertNoError(t, err, "Failed to open the ledger")
	ledger.Close()

	
	
	provider.(*Provider).idStore.setUnderConstructionFlag(constructTestLedgerID(2))
	provider.Close()

	
	provider, err = NewProvider()
	testutil.AssertNoError(t, err, "Provider failed to recover an underConstructionLedger")
	flag, err = provider.(*Provider).idStore.getUnderConstructionFlag()
	testutil.AssertNoError(t, err, "Failed to read the underconstruction flag")
	testutil.AssertEquals(t, flag, "")

}

func TestMultipleLedgerBasicRW(t *testing.T) {
	env := newTestEnv(t)
	defer env.cleanup()
	numLedgers := 10
	provider, _ := NewProvider()
	ledgers := make([]ledgerproto.PeerLedger, numLedgers)
	for i := 0; i < numLedgers; i++ {
		bg, gb := testutil.NewBlockGenerator(t, constructTestLedgerID(i), false)
		l, err := provider.Create(gb)
		testutil.AssertNoError(t, err, "")
		ledgers[i] = l
		txid := util.GenerateUUID()
		s, _ := l.NewTxSimulator(txid)
		err = s.SetState("ns", "testKey", []byte(fmt.Sprintf("testValue_%d", i)))
		s.Done()
		testutil.AssertNoError(t, err, "")
		res, err := s.GetTxSimulationResults()
		testutil.AssertNoError(t, err, "")
		pubSimBytes, _ := res.GetPubSimulationBytes()
		b := bg.NextBlock([][]byte{pubSimBytes})
		err = l.CommitWithPvtData(&ledgerproto.BlockAndPvtData{Block: b})
		l.Close()
		testutil.AssertNoError(t, err, "")
	}

	provider.Close()

	provider, _ = NewProvider()
	defer provider.Close()
	ledgers = make([]ledgerproto.PeerLedger, numLedgers)
	for i := 0; i < numLedgers; i++ {
		l, err := provider.Open(constructTestLedgerID(i))
		testutil.AssertNoError(t, err, "")
		ledgers[i] = l
	}

	for i, l := range ledgers {
		q, _ := l.NewQueryExecutor()
		val, err := q.GetState("ns", "testKey")
		q.Done()
		testutil.AssertNoError(t, err, "")
		testutil.AssertEquals(t, val, []byte(fmt.Sprintf("testValue_%d", i)))
		l.Close()
	}
}

func TestLedgerBackup(t *testing.T) {
	ledgerid := "TestLedger"
	originalPath := "/tmp/blockchain/ledgertests/kvledger1"
	restorePath := "/tmp/blockchain/ledgertests/kvledger2"
	viper.Set("ledger.history.enableHistoryDatabase", true)

	
	env := createTestEnv(t, originalPath)
	provider, _ := NewProvider()
	bg, gb := testutil.NewBlockGenerator(t, ledgerid, false)
	gbHash := gb.Header.Hash()
	ledger, _ := provider.Create(gb)

	txid := util.GenerateUUID()
	simulator, _ := ledger.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value1"))
	simulator.SetState("ns1", "key2", []byte("value2"))
	simulator.SetState("ns1", "key3", []byte("value3"))
	simulator.Done()
	simRes, _ := simulator.GetTxSimulationResults()
	pubSimBytes, _ := simRes.GetPubSimulationBytes()
	block1 := bg.NextBlock([][]byte{pubSimBytes})
	ledger.CommitWithPvtData(&ledgerproto.BlockAndPvtData{Block: block1})

	txid = util.GenerateUUID()
	simulator, _ = ledger.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value4"))
	simulator.SetState("ns1", "key2", []byte("value5"))
	simulator.SetState("ns1", "key3", []byte("value6"))
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimBytes, _ = simRes.GetPubSimulationBytes()
	block2 := bg.NextBlock([][]byte{pubSimBytes})
	ledger.CommitWithPvtData(&ledgerproto.BlockAndPvtData{Block: block2})

	ledger.Close()
	provider.Close()

	
	env = createTestEnv(t, restorePath)

	
	
	testutil.AssertNoError(t, os.RemoveAll(ledgerconfig.GetStateLevelDBPath()), "")
	testutil.AssertNoError(t, os.RemoveAll(ledgerconfig.GetHistoryLevelDBPath()), "")
	testutil.AssertNoError(t, os.RemoveAll(filepath.Join(ledgerconfig.GetBlockStorePath(), fsblkstorage.IndexDir)), "")
	testutil.AssertNoError(t, os.Rename(originalPath, restorePath), "")
	defer env.cleanup()

	
	provider, _ = NewProvider()
	defer provider.Close()

	_, err := provider.Create(gb)
	testutil.AssertEquals(t, err, ErrLedgerIDExists)

	ledger, _ = provider.Open(ledgerid)
	defer ledger.Close()

	block1Hash := block1.Header.Hash()
	block2Hash := block2.Header.Hash()
	bcInfo, _ := ledger.GetBlockchainInfo()
	testutil.AssertEquals(t, bcInfo, &common.BlockchainInfo{
		Height: 3, CurrentBlockHash: block2Hash, PreviousBlockHash: block1Hash})

	b0, _ := ledger.GetBlockByHash(gbHash)
	testutil.AssertEquals(t, b0, gb)

	b1, _ := ledger.GetBlockByHash(block1Hash)
	testutil.AssertEquals(t, b1, block1)

	b2, _ := ledger.GetBlockByHash(block2Hash)
	testutil.AssertEquals(t, b2, block2)

	b0, _ = ledger.GetBlockByNumber(0)
	testutil.AssertEquals(t, b0, gb)

	b1, _ = ledger.GetBlockByNumber(1)
	testutil.AssertEquals(t, b1, block1)

	b2, _ = ledger.GetBlockByNumber(2)
	testutil.AssertEquals(t, b2, block2)

	
	txEnvBytes2 := block1.Data.Data[0]
	txEnv2, err := putils.GetEnvelopeFromBlock(txEnvBytes2)
	testutil.AssertNoError(t, err, "Error upon GetEnvelopeFromBlock")
	payload2, err := putils.GetPayload(txEnv2)
	testutil.AssertNoError(t, err, "Error upon GetPayload")
	chdr, err := putils.UnmarshalChannelHeader(payload2.Header.ChannelHeader)
	testutil.AssertNoError(t, err, "Error upon GetChannelHeaderFromBytes")
	txID2 := chdr.TxId
	processedTran2, err := ledger.GetTransactionByID(txID2)
	testutil.AssertNoError(t, err, "Error upon GetTransactionByID")
	
	retrievedTxEnv2 := processedTran2.TransactionEnvelope
	testutil.AssertEquals(t, retrievedTxEnv2, txEnv2)

	qe, _ := ledger.NewQueryExecutor()
	value1, _ := qe.GetState("ns1", "key1")
	testutil.AssertEquals(t, value1, []byte("value4"))

	hqe, err := ledger.NewHistoryQueryExecutor()
	testutil.AssertNoError(t, err, "")
	itr, err := hqe.GetHistoryForKey("ns1", "key1")
	testutil.AssertNoError(t, err, "")
	defer itr.Close()

	result1, err := itr.Next()
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, result1.(*queryresult.KeyModification).Value, []byte("value1"))
	result2, err := itr.Next()
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, result2.(*queryresult.KeyModification).Value, []byte("value4"))
}

func constructTestLedgerID(i int) string {
	return fmt.Sprintf("ledger_%06d", i)
}
