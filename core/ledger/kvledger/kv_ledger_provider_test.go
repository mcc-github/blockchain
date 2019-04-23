/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage/fsblkstorage"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/common/util"
	lgr "github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/mock"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLedgerProvider(t *testing.T) {
	conf, cleanup := testConfig(t)
	defer cleanup()
	
	_ = createTestEnv(t, conf.RootFSPath)
	provider := testutilNewProvider(conf, t)
	numLedgers := 10
	existingLedgerIDs, err := provider.List()
	assert.NoError(t, err)
	assert.Len(t, existingLedgerIDs, 0)
	genesisBlocks := make([]*common.Block, numLedgers)
	for i := 0; i < numLedgers; i++ {
		genesisBlock, _ := configtxtest.MakeGenesisBlock(constructTestLedgerID(i))
		genesisBlocks[i] = genesisBlock
		provider.Create(genesisBlock)
	}
	existingLedgerIDs, err = provider.List()
	assert.NoError(t, err)
	assert.Len(t, existingLedgerIDs, numLedgers)

	provider.Close()

	provider = testutilNewProvider(conf, t)
	defer provider.Close()
	ledgerIds, _ := provider.List()
	assert.Len(t, ledgerIds, numLedgers)
	t.Logf("ledgerIDs=%#v", ledgerIds)
	for i := 0; i < numLedgers; i++ {
		assert.Equal(t, constructTestLedgerID(i), ledgerIds[i])
	}
	for i := 0; i < numLedgers; i++ {
		ledgerid := constructTestLedgerID(i)
		status, _ := provider.Exists(ledgerid)
		assert.True(t, status)
		ledger, err := provider.Open(ledgerid)
		assert.NoError(t, err)
		bcInfo, err := ledger.GetBlockchainInfo()
		ledger.Close()
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), bcInfo.Height)

		
		s := provider.(*Provider).idStore
		gbBytesInProviderStore, err := s.db.Get(s.encodeLedgerKey(ledgerid))
		assert.NoError(t, err)
		gb := &common.Block{}
		assert.NoError(t, proto.Unmarshal(gbBytesInProviderStore, gb))
		assert.True(t, proto.Equal(gb, genesisBlocks[i]), "proto messages are not equal")
	}
	gb, _ := configtxtest.MakeGenesisBlock(constructTestLedgerID(2))
	_, err = provider.Create(gb)
	assert.Equal(t, ErrLedgerIDExists, err)

	status, err := provider.Exists(constructTestLedgerID(numLedgers))
	assert.NoError(t, err, "Failed to check for ledger existence")
	assert.Equal(t, status, false)

	_, err = provider.Open(constructTestLedgerID(numLedgers))
	assert.Equal(t, ErrNonExistingLedgerID, err)
}

func TestRecovery(t *testing.T) {
	conf, cleanup := testConfig(t)
	defer cleanup()
	
	_ = createTestEnv(t, conf.RootFSPath)
	provider := testutilNewProvider(conf, t)

	
	genesisBlock, _ := configtxtest.MakeGenesisBlock(constructTestLedgerID(1))
	ledger, err := provider.(*Provider).openInternal(constructTestLedgerID(1))
	ledger.CommitWithPvtData(&lgr.BlockAndPvtData{Block: genesisBlock})
	ledger.Close()

	
	
	provider.(*Provider).idStore.setUnderConstructionFlag(constructTestLedgerID(1))
	provider.Close()

	
	provider = testutilNewProvider(conf, t)
	
	flag, err := provider.(*Provider).idStore.getUnderConstructionFlag()
	assert.NoError(t, err, "Failed to read the underconstruction flag")
	assert.Equal(t, "", flag)
	ledger, err = provider.Open(constructTestLedgerID(1))
	assert.NoError(t, err, "Failed to open the ledger")
	ledger.Close()

	
	
	provider.(*Provider).idStore.setUnderConstructionFlag(constructTestLedgerID(2))
	provider.Close()

	
	provider = testutilNewProvider(conf, t)
	assert.NoError(t, err, "Provider failed to recover an underConstructionLedger")
	flag, err = provider.(*Provider).idStore.getUnderConstructionFlag()
	assert.NoError(t, err, "Failed to read the underconstruction flag")
	assert.Equal(t, "", flag)

}

func TestMultipleLedgerBasicRW(t *testing.T) {
	conf, cleanup := testConfig(t)
	defer cleanup()
	
	_ = createTestEnv(t, conf.RootFSPath)
	provider := testutilNewProvider(conf, t)
	numLedgers := 10
	ledgers := make([]lgr.PeerLedger, numLedgers)
	for i := 0; i < numLedgers; i++ {
		bg, gb := testutil.NewBlockGenerator(t, constructTestLedgerID(i), false)
		l, err := provider.Create(gb)
		assert.NoError(t, err)
		ledgers[i] = l
		txid := util.GenerateUUID()
		s, _ := l.NewTxSimulator(txid)
		err = s.SetState("ns", "testKey", []byte(fmt.Sprintf("testValue_%d", i)))
		s.Done()
		assert.NoError(t, err)
		res, err := s.GetTxSimulationResults()
		assert.NoError(t, err)
		pubSimBytes, _ := res.GetPubSimulationBytes()
		b := bg.NextBlock([][]byte{pubSimBytes})
		err = l.CommitWithPvtData(&lgr.BlockAndPvtData{Block: b})
		l.Close()
		assert.NoError(t, err)
	}

	provider.Close()

	provider = testutilNewProvider(conf, t)
	defer provider.Close()
	ledgers = make([]lgr.PeerLedger, numLedgers)
	for i := 0; i < numLedgers; i++ {
		l, err := provider.Open(constructTestLedgerID(i))
		assert.NoError(t, err)
		ledgers[i] = l
	}

	for i, l := range ledgers {
		q, _ := l.NewQueryExecutor()
		val, err := q.GetState("ns", "testKey")
		q.Done()
		assert.NoError(t, err)
		assert.Equal(t, []byte(fmt.Sprintf("testValue_%d", i)), val)
		l.Close()
	}
}

func TestLedgerBackup(t *testing.T) {
	ledgerid := "TestLedger"
	originalPath := "/tmp/blockchain/ledgertests/kvledger1"
	restorePath := "/tmp/blockchain/ledgertests/kvledger2"
	viper.Set("ledger.history.enableHistoryDatabase", true)

	
	env := createTestEnv(t, originalPath)
	origConf := &lgr.Config{
		RootFSPath: originalPath,
		StateDB: &lgr.StateDB{
			LevelDBPath: filepath.Join(originalPath, "stateLeveldb"),
		},
		PrivateData: &lgr.PrivateData{
			StorePath:       filepath.Join(originalPath, "pvtdataStore"),
			MaxBatchSize:    5000,
			BatchesInterval: 1000,
			PurgeInterval:   100,
		},
	}
	provider := testutilNewProvider(origConf, t)
	bg, gb := testutil.NewBlockGenerator(t, ledgerid, false)
	gbHash := protoutil.BlockHeaderHash(gb.Header)
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
	ledger.CommitWithPvtData(&lgr.BlockAndPvtData{Block: block1})

	txid = util.GenerateUUID()
	simulator, _ = ledger.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value4"))
	simulator.SetState("ns1", "key2", []byte("value5"))
	simulator.SetState("ns1", "key3", []byte("value6"))
	simulator.Done()
	simRes, _ = simulator.GetTxSimulationResults()
	pubSimBytes, _ = simRes.GetPubSimulationBytes()
	block2 := bg.NextBlock([][]byte{pubSimBytes})
	ledger.CommitWithPvtData(&lgr.BlockAndPvtData{Block: block2})

	ledger.Close()
	provider.Close()

	
	env = createTestEnv(t, restorePath)

	
	
	assert.NoError(t, os.RemoveAll(filepath.Join(originalPath, "stateLeveldb")))
	assert.NoError(t, os.RemoveAll(filepath.Join(originalPath, "historyLeveldb")))
	assert.NoError(t, os.RemoveAll(filepath.Join(originalPath, "chains", fsblkstorage.IndexDir)))
	assert.NoError(t, os.Rename(originalPath, restorePath))
	defer env.cleanup()

	
	restoreConf := &lgr.Config{
		RootFSPath: restorePath,
		StateDB: &lgr.StateDB{
			LevelDBPath: filepath.Join(restorePath, "stateLeveldb"),
		},
		PrivateData: &lgr.PrivateData{
			StorePath:       filepath.Join(restorePath, "pvtdataStore"),
			MaxBatchSize:    5000,
			BatchesInterval: 1000,
			PurgeInterval:   100,
		},
	}
	provider = testutilNewProvider(restoreConf, t)
	defer provider.Close()

	_, err := provider.Create(gb)
	assert.Equal(t, ErrLedgerIDExists, err)

	ledger, _ = provider.Open(ledgerid)
	defer ledger.Close()

	block1Hash := protoutil.BlockHeaderHash(block1.Header)
	block2Hash := protoutil.BlockHeaderHash(block2.Header)
	bcInfo, _ := ledger.GetBlockchainInfo()
	assert.Equal(t, &common.BlockchainInfo{
		Height: 3, CurrentBlockHash: block2Hash, PreviousBlockHash: block1Hash,
	}, bcInfo)

	b0, _ := ledger.GetBlockByHash(gbHash)
	assert.True(t, proto.Equal(b0, gb), "proto messages are not equal")

	b1, _ := ledger.GetBlockByHash(block1Hash)
	assert.True(t, proto.Equal(b1, block1), "proto messages are not equal")

	b2, _ := ledger.GetBlockByHash(block2Hash)
	assert.True(t, proto.Equal(b2, block2), "proto messages are not equal")

	b0, _ = ledger.GetBlockByNumber(0)
	assert.True(t, proto.Equal(b0, gb), "proto messages are not equal")

	b1, _ = ledger.GetBlockByNumber(1)
	assert.True(t, proto.Equal(b1, block1), "proto messages are not equal")

	b2, _ = ledger.GetBlockByNumber(2)
	assert.True(t, proto.Equal(b2, block2), "proto messages are not equal")

	
	txEnvBytes2 := block1.Data.Data[0]
	txEnv2, err := protoutil.GetEnvelopeFromBlock(txEnvBytes2)
	assert.NoError(t, err, "Error upon GetEnvelopeFromBlock")
	payload2, err := protoutil.GetPayload(txEnv2)
	assert.NoError(t, err, "Error upon GetPayload")
	chdr, err := protoutil.UnmarshalChannelHeader(payload2.Header.ChannelHeader)
	assert.NoError(t, err, "Error upon GetChannelHeaderFromBytes")
	txID2 := chdr.TxId
	processedTran2, err := ledger.GetTransactionByID(txID2)
	assert.NoError(t, err, "Error upon GetTransactionByID")
	
	retrievedTxEnv2 := processedTran2.TransactionEnvelope
	assert.Equal(t, txEnv2, retrievedTxEnv2)

	qe, _ := ledger.NewQueryExecutor()
	value1, _ := qe.GetState("ns1", "key1")
	assert.Equal(t, []byte("value4"), value1)

	hqe, err := ledger.NewHistoryQueryExecutor()
	assert.NoError(t, err)
	itr, err := hqe.GetHistoryForKey("ns1", "key1")
	assert.NoError(t, err)
	defer itr.Close()

	result1, err := itr.Next()
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), result1.(*queryresult.KeyModification).Value)
	result2, err := itr.Next()
	assert.NoError(t, err)
	assert.Equal(t, []byte("value4"), result2.(*queryresult.KeyModification).Value)
}

func constructTestLedgerID(i int) string {
	return fmt.Sprintf("ledger_%06d", i)
}

func testConfig(t *testing.T) (conf *lgr.Config, cleanup func()) {
	path, err := ioutil.TempDir("", "kvledger")
	if err != nil {
		t.Fatalf("Failed to create test ledger directory: %s", err)
	}
	conf = &lgr.Config{
		RootFSPath: path,
		StateDB: &lgr.StateDB{
			LevelDBPath: filepath.Join(path, "stateLeveldb"),
		},
		PrivateData: &lgr.PrivateData{
			StorePath:       filepath.Join(path, "pvtdataStore"),
			MaxBatchSize:    5000,
			BatchesInterval: 1000,
			PurgeInterval:   100,
		},
	}
	cleanup = func() {
		os.RemoveAll(path)
	}

	return conf, cleanup
}

func testutilNewProvider(conf *lgr.Config, t *testing.T) lgr.PeerLedgerProvider {
	provider, err := NewProvider()
	assert.NoError(t, err)
	provider.Initialize(&lgr.Initializer{
		DeployedChaincodeInfoProvider: &mock.DeployedChaincodeInfoProvider{},
		MetricsProvider:               &disabled.Provider{},
		Config:                        conf,
	})
	return provider
}

func testutilNewProviderWithCollectionConfig(
	t *testing.T,
	namespace string,
	btlConfigs map[string]uint64,
	conf *lgr.Config,
) lgr.PeerLedgerProvider {
	provider := testutilNewProvider(conf, t)
	mockCCInfoProvider := provider.(*Provider).initializer.DeployedChaincodeInfoProvider.(*mock.DeployedChaincodeInfoProvider)
	collMap := map[string]*common.StaticCollectionConfig{}
	var collConf []*common.CollectionConfig
	for collName, btl := range btlConfigs {
		staticConf := &common.StaticCollectionConfig{Name: collName, BlockToLive: btl}
		collMap[collName] = staticConf
		collectionConf := &common.CollectionConfig{}
		collectionConf.Payload = &common.CollectionConfig_StaticCollectionConfig{StaticCollectionConfig: staticConf}
		collConf = append(collConf, collectionConf)
	}
	collectionConfPkg := &common.CollectionConfigPackage{Config: collConf}

	mockCCInfoProvider.ChaincodeInfoStub = func(channelName, ccName string, qe lgr.SimpleQueryExecutor) (*lgr.DeployedChaincodeInfo, error) {
		if ccName == namespace {
			return &lgr.DeployedChaincodeInfo{
				Name: namespace, ExplicitCollectionConfigPkg: collectionConfPkg}, nil
		}
		return nil, nil
	}

	mockCCInfoProvider.CollectionInfoStub = func(channelName, ccName, collName string, qe lgr.SimpleQueryExecutor) (*common.StaticCollectionConfig, error) {
		if ccName == namespace {
			return collMap[collName], nil
		}
		return nil, nil
	}
	return provider
}
