/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transientstore

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/mcc-github/blockchain-protos-go/transientstore"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	commonutil "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tempdir, err := ioutil.TempDir("", "ts")
	if err != nil {
		panic(err)
	}

	rc := m.Run()

	os.RemoveAll(tempdir)
	os.Exit(rc)
}

func testStore(t *testing.T) (s Store, cleanup func()) {
	tempdir, err := ioutil.TempDir("", "ts")
	if err != nil {
		t.Fatalf("Failed to create test directory: %s", err)
	}

	cleanup = func() {
		os.RemoveAll(tempdir)
	}

	sp, err := NewStoreProvider(tempdir)
	if err != nil {
		t.Fatalf("Failed to open test store: %s", err)
	}
	s, err = sp.OpenStore("TestStore")
	require.NoError(t, err)
	return s, cleanup
}

func TestPurgeIndexKeyCodingEncoding(t *testing.T) {
	assert := assert.New(t)
	blkHts := []uint64{0, 10, 20000}
	txids := []string{"txid", ""}
	uuids := []string{"uuid", ""}
	for _, blkHt := range blkHts {
		for _, txid := range txids {
			for _, uuid := range uuids {
				testCase := fmt.Sprintf("blkHt=%d,txid=%s,uuid=%s", blkHt, txid, uuid)
				t.Run(testCase, func(t *testing.T) {
					t.Logf("Running test case [%s]", testCase)
					purgeIndexKey := createCompositeKeyForPurgeIndexByHeight(blkHt, txid, uuid)
					txid1, uuid1, blkHt1, err := splitCompositeKeyOfPurgeIndexByHeight(purgeIndexKey)
					assert.NoError(err)
					assert.Equal(txid, txid1)
					assert.Equal(uuid, uuid1)
					assert.Equal(blkHt, blkHt1)
				})
			}
		}
	}
}

func TestRWSetKeyCodingEncoding(t *testing.T) {
	assert := assert.New(t)
	blkHts := []uint64{0, 10, 20000}
	txids := []string{"txid", ""}
	uuids := []string{"uuid", ""}
	for _, blkHt := range blkHts {
		for _, txid := range txids {
			for _, uuid := range uuids {
				testCase := fmt.Sprintf("blkHt=%d,txid=%s,uuid=%s", blkHt, txid, uuid)
				t.Run(testCase, func(t *testing.T) {
					t.Logf("Running test case [%s]", testCase)
					rwsetKey := createCompositeKeyForPvtRWSet(txid, uuid, blkHt)
					uuid1, blkHt1, err := splitCompositeKeyOfPvtRWSet(rwsetKey)
					assert.NoError(err)
					assert.Equal(uuid, uuid1)
					assert.Equal(blkHt, blkHt1)
				})
			}
		}
	}
}

func TestTransientStorePersistAndRetrieve(t *testing.T) {
	testStore, cleanup := testStore(t)
	defer cleanup()
	assert := assert.New(t)
	txid := "txid-1"
	samplePvtRWSetWithConfig := samplePvtDataWithConfigInfo(t)

	
	var endorsersResults []*EndorserPvtSimulationResults

	
	endorser0SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          10,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser0SimulationResults)

	
	endorser1SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          10,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser1SimulationResults)

	
	var err error
	for i := 0; i < len(endorsersResults); i++ {
		err = testStore.Persist(txid, endorsersResults[i].ReceivedAtBlockHeight,
			endorsersResults[i].PvtSimulationResultsWithConfig)
		assert.NoError(err)
	}

	
	iter, err := testStore.GetTxPvtRWSetByTxid(txid, nil)
	assert.NoError(err)

	var actualEndorsersResults []*EndorserPvtSimulationResults
	for {
		result, err := iter.Next()
		assert.NoError(err)
		if result == nil {
			break
		}
		actualEndorsersResults = append(actualEndorsersResults, result)
	}
	iter.Close()
	sortResults(endorsersResults)
	sortResults(actualEndorsersResults)
	assert.Equal(endorsersResults, actualEndorsersResults)
}

func TestTransientStorePersistAndRetrieveBothOldAndNewProto(t *testing.T) {
	testStore, cleanup := testStore(t)
	defer cleanup()
	assert := assert.New(t)
	txid := "txid-1"
	var receivedAtBlockHeight uint64 = 10
	var err error

	
	samplePvtRWSet := samplePvtData(t)
	err = testStore.(*store).persistOldProto(txid, receivedAtBlockHeight, samplePvtRWSet)
	assert.NoError(err)

	
	samplePvtRWSetWithConfig := samplePvtDataWithConfigInfo(t)
	err = testStore.Persist(txid, receivedAtBlockHeight, samplePvtRWSetWithConfig)
	assert.NoError(err)

	
	var expectedEndorsersResults []*EndorserPvtSimulationResults

	pvtRWSetWithConfigInfo := &transientstore.TxPvtReadWriteSetWithConfigInfo{
		PvtRwset: samplePvtRWSet,
	}

	endorser0SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          receivedAtBlockHeight,
		PvtSimulationResultsWithConfig: pvtRWSetWithConfigInfo,
	}
	expectedEndorsersResults = append(expectedEndorsersResults, endorser0SimulationResults)

	endorser1SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          receivedAtBlockHeight,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	expectedEndorsersResults = append(expectedEndorsersResults, endorser1SimulationResults)

	
	iter, err := testStore.GetTxPvtRWSetByTxid(txid, nil)
	assert.NoError(err)

	var actualEndorsersResults []*EndorserPvtSimulationResults
	for {
		result, err := iter.Next()
		assert.NoError(err)
		if result == nil {
			break
		}
		actualEndorsersResults = append(actualEndorsersResults, result)
	}
	iter.Close()
	sortResults(expectedEndorsersResults)
	sortResults(actualEndorsersResults)
	assert.Equal(expectedEndorsersResults, actualEndorsersResults)
}

func TestTransientStorePurgeByTxids(t *testing.T) {
	testStore, cleanup := testStore(t)
	defer cleanup()
	assert := assert.New(t)

	var txids []string
	var endorsersResults []*EndorserPvtSimulationResults

	samplePvtRWSetWithConfig := samplePvtDataWithConfigInfo(t)

	
	txids = append(txids, "txid-1")
	endorser0SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          10,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser0SimulationResults)

	txids = append(txids, "txid-1")
	endorser1SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          11,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser1SimulationResults)

	
	txids = append(txids, "txid-2")
	endorser2SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          11,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser2SimulationResults)

	
	txids = append(txids, "txid-3")
	endorser3SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          12,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser3SimulationResults)

	txids = append(txids, "txid-3")
	endorser4SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          12,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser4SimulationResults)

	txids = append(txids, "txid-3")
	endorser5SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          13,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser5SimulationResults)

	var err error
	for i := 0; i < len(txids); i++ {
		err = testStore.Persist(txids[i], endorsersResults[i].ReceivedAtBlockHeight,
			endorsersResults[i].PvtSimulationResultsWithConfig)
		assert.NoError(err)
	}

	
	iter, err := testStore.GetTxPvtRWSetByTxid("txid-2", nil)
	assert.NoError(err)

	
	var expectedEndorsersResults []*EndorserPvtSimulationResults
	expectedEndorsersResults = append(expectedEndorsersResults, endorser2SimulationResults)

	
	var actualEndorsersResults []*EndorserPvtSimulationResults
	for true {
		result, err := iter.Next()
		assert.NoError(err)
		if result == nil {
			break
		}
		actualEndorsersResults = append(actualEndorsersResults, result)
	}
	iter.Close()

	
	
	sortResults(expectedEndorsersResults)
	sortResults(actualEndorsersResults)

	assert.Equal(len(expectedEndorsersResults), len(actualEndorsersResults))
	for i, expected := range expectedEndorsersResults {
		assert.Equal(expected.ReceivedAtBlockHeight, actualEndorsersResults[i].ReceivedAtBlockHeight)
		assert.True(proto.Equal(expected.PvtSimulationResultsWithConfig, actualEndorsersResults[i].PvtSimulationResultsWithConfig))
	}

	
	toRemoveTxids := []string{"txid-2", "txid-3"}
	err = testStore.PurgeByTxids(toRemoveTxids)
	assert.NoError(err)

	for _, txid := range toRemoveTxids {

		
		var expectedEndorsersResults *EndorserPvtSimulationResults
		expectedEndorsersResults = nil
		iter, err = testStore.GetTxPvtRWSetByTxid(txid, nil)
		assert.NoError(err)
		
		result, err := iter.Next()
		assert.NoError(err)
		assert.Equal(expectedEndorsersResults, result)
	}

	
	iter, err = testStore.GetTxPvtRWSetByTxid("txid-1", nil)
	assert.NoError(err)

	
	expectedEndorsersResults = nil
	expectedEndorsersResults = append(expectedEndorsersResults, endorser0SimulationResults)
	expectedEndorsersResults = append(expectedEndorsersResults, endorser1SimulationResults)

	
	actualEndorsersResults = nil
	for true {
		result, err := iter.Next()
		assert.NoError(err)
		if result == nil {
			break
		}
		actualEndorsersResults = append(actualEndorsersResults, result)
	}
	iter.Close()

	
	
	sortResults(expectedEndorsersResults)
	sortResults(actualEndorsersResults)

	assert.Equal(len(expectedEndorsersResults), len(actualEndorsersResults))
	for i, expected := range expectedEndorsersResults {
		assert.Equal(expected.ReceivedAtBlockHeight, actualEndorsersResults[i].ReceivedAtBlockHeight)
		assert.True(proto.Equal(expected.PvtSimulationResultsWithConfig, actualEndorsersResults[i].PvtSimulationResultsWithConfig))
	}

	toRemoveTxids = []string{"txid-1"}
	err = testStore.PurgeByTxids(toRemoveTxids)
	assert.NoError(err)

	for _, txid := range toRemoveTxids {

		
		var expectedEndorsersResults *EndorserPvtSimulationResults
		expectedEndorsersResults = nil
		iter, err = testStore.GetTxPvtRWSetByTxid(txid, nil)
		assert.NoError(err)
		
		result, err := iter.Next()
		assert.NoError(err)
		assert.Equal(expectedEndorsersResults, result)
	}

	
	_, err = testStore.GetMinTransientBlkHt()
	assert.Equal(err, ErrStoreEmpty)
}

func TestTransientStorePurgeByHeight(t *testing.T) {
	testStore, cleanup := testStore(t)
	defer cleanup()
	assert := assert.New(t)

	txid := "txid-1"
	samplePvtRWSetWithConfig := samplePvtDataWithConfigInfo(t)

	
	var endorsersResults []*EndorserPvtSimulationResults

	
	endorser0SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          10,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser0SimulationResults)

	
	endorser1SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          11,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser1SimulationResults)

	
	endorser2SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          12,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser2SimulationResults)

	
	endorser3SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          12,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser3SimulationResults)

	
	endorser4SimulationResults := &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight:          13,
		PvtSimulationResultsWithConfig: samplePvtRWSetWithConfig,
	}
	endorsersResults = append(endorsersResults, endorser4SimulationResults)

	
	var err error
	for i := 0; i < 5; i++ {
		err = testStore.Persist(txid, endorsersResults[i].ReceivedAtBlockHeight,
			endorsersResults[i].PvtSimulationResultsWithConfig)
		assert.NoError(err)
	}

	
	minTransientBlkHtToRetain := uint64(12)
	err = testStore.PurgeByHeight(minTransientBlkHtToRetain)
	assert.NoError(err)

	
	iter, err := testStore.GetTxPvtRWSetByTxid(txid, nil)
	assert.NoError(err)

	
	var expectedEndorsersResults []*EndorserPvtSimulationResults
	expectedEndorsersResults = append(expectedEndorsersResults, endorser2SimulationResults) 
	expectedEndorsersResults = append(expectedEndorsersResults, endorser3SimulationResults) 
	expectedEndorsersResults = append(expectedEndorsersResults, endorser4SimulationResults) 

	
	var actualEndorsersResults []*EndorserPvtSimulationResults
	for true {
		result, err := iter.Next()
		assert.NoError(err)
		if result == nil {
			break
		}
		actualEndorsersResults = append(actualEndorsersResults, result)
	}
	iter.Close()

	
	
	sortResults(expectedEndorsersResults)
	sortResults(actualEndorsersResults)

	assert.Equal(len(expectedEndorsersResults), len(actualEndorsersResults))
	for i, expected := range expectedEndorsersResults {
		assert.Equal(expected.ReceivedAtBlockHeight, actualEndorsersResults[i].ReceivedAtBlockHeight)
		assert.True(proto.Equal(expected.PvtSimulationResultsWithConfig, actualEndorsersResults[i].PvtSimulationResultsWithConfig))
	}

	
	var actualMinTransientBlkHt uint64
	actualMinTransientBlkHt, err = testStore.GetMinTransientBlkHt()
	assert.NoError(err)
	assert.Equal(minTransientBlkHtToRetain, actualMinTransientBlkHt)

	
	minTransientBlkHtToRetain = uint64(15)
	err = testStore.PurgeByHeight(minTransientBlkHtToRetain)
	assert.NoError(err)

	
	actualMinTransientBlkHt, err = testStore.GetMinTransientBlkHt()
	assert.Equal(err, ErrStoreEmpty)

	
	minTransientBlkHtToRetain = uint64(15)
	err = testStore.PurgeByHeight(minTransientBlkHtToRetain)
	
	assert.NoError(err)
}

func TestTransientStoreRetrievalWithFilter(t *testing.T) {
	testStore, cleanup := testStore(t)
	defer cleanup()

	samplePvtSimResWithConfig := samplePvtDataWithConfigInfo(t)

	testTxid := "testTxid"
	numEntries := 5
	for i := 0; i < numEntries; i++ {
		testStore.Persist(testTxid, uint64(i), samplePvtSimResWithConfig)
	}

	filter := ledger.NewPvtNsCollFilter()
	filter.Add("ns-1", "coll-1")
	filter.Add("ns-2", "coll-2")

	itr, err := testStore.GetTxPvtRWSetByTxid(testTxid, filter)
	assert.NoError(t, err)

	var actualRes []*EndorserPvtSimulationResults
	for {
		res, err := itr.Next()
		if res == nil || err != nil {
			assert.NoError(t, err)
			break
		}
		actualRes = append(actualRes, res)
	}

	
	expectedSimulationRes := samplePvtSimResWithConfig
	expectedSimulationRes.GetPvtRwset().NsPvtRwset[0].CollectionPvtRwset = expectedSimulationRes.GetPvtRwset().NsPvtRwset[0].CollectionPvtRwset[0:1]
	expectedSimulationRes.GetPvtRwset().NsPvtRwset[1].CollectionPvtRwset = expectedSimulationRes.GetPvtRwset().NsPvtRwset[1].CollectionPvtRwset[1:]
	expectedSimulationRes.CollectionConfigs, err = trimPvtCollectionConfigs(expectedSimulationRes.CollectionConfigs, filter)
	assert.NoError(t, err)
	for ns, colName := range map[string]string{"ns-1": "coll-1", "ns-2": "coll-2"} {
		config := expectedSimulationRes.CollectionConfigs[ns]
		assert.NotNil(t, config)
		ns1Config := config.Config
		assert.Equal(t, len(ns1Config), 1)
		ns1ColConfig := ns1Config[0].GetStaticCollectionConfig()
		assert.NotNil(t, ns1ColConfig.Name, colName)
	}

	var expectedRes []*EndorserPvtSimulationResults
	for i := 0; i < numEntries; i++ {
		expectedRes = append(expectedRes, &EndorserPvtSimulationResults{uint64(i), expectedSimulationRes})
	}

	
	
	sortResults(expectedRes)
	sortResults(actualRes)
	assert.Equal(t, len(expectedRes), len(actualRes))
	for i, expected := range expectedRes {
		assert.Equal(t, expected.ReceivedAtBlockHeight, actualRes[i].ReceivedAtBlockHeight)
		assert.True(t, proto.Equal(expected.PvtSimulationResultsWithConfig, actualRes[i].PvtSimulationResultsWithConfig))
	}

}

func sortResults(res []*EndorserPvtSimulationResults) {
	
	
	var sortCondition = func(i, j int) bool {
		if res[i].ReceivedAtBlockHeight == res[j].ReceivedAtBlockHeight {
			resI, _ := proto.Marshal(res[i].PvtSimulationResultsWithConfig)
			resJ, _ := proto.Marshal(res[j].PvtSimulationResultsWithConfig)
			
			return string(util.ComputeHash(resI)) < string(util.ComputeHash(resJ))
		}
		return res[i].ReceivedAtBlockHeight < res[j].ReceivedAtBlockHeight
	}
	sort.SliceStable(res, sortCondition)
}

func samplePvtData(t *testing.T) *rwset.TxPvtReadWriteSet {
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

		{
			Namespace: "ns-2",
			CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
				{
					CollectionName: "coll-1",
					Rwset:          []byte("RandomBytes-PvtRWSet-ns2-coll1"),
				},
				{
					CollectionName: "coll-2",
					Rwset:          []byte("RandomBytes-PvtRWSet-ns2-coll2"),
				},
			},
		},
	}
	return pvtWriteSet
}

func samplePvtDataWithConfigInfo(t *testing.T) *transientstore.TxPvtReadWriteSetWithConfigInfo {
	pvtWriteSet := samplePvtData(t)
	pvtRWSetWithConfigInfo := &transientstore.TxPvtReadWriteSetWithConfigInfo{
		PvtRwset: pvtWriteSet,
		CollectionConfigs: map[string]*common.CollectionConfigPackage{
			"ns-1": {
				Config: []*common.CollectionConfig{
					sampleCollectionConfigPackage("coll-1"),
					sampleCollectionConfigPackage("coll-2"),
				},
			},
			"ns-2": {
				Config: []*common.CollectionConfig{
					sampleCollectionConfigPackage("coll-1"),
					sampleCollectionConfigPackage("coll-2"),
				},
			},
		},
	}
	return pvtRWSetWithConfigInfo
}

func createCollectionConfig(collectionName string, signaturePolicyEnvelope *common.SignaturePolicyEnvelope,
	requiredPeerCount int32, maximumPeerCount int32,
) *common.CollectionConfig {
	signaturePolicy := &common.CollectionPolicyConfig_SignaturePolicy{
		SignaturePolicy: signaturePolicyEnvelope,
	}
	accessPolicy := &common.CollectionPolicyConfig{
		Payload: signaturePolicy,
	}

	return &common.CollectionConfig{
		Payload: &common.CollectionConfig_StaticCollectionConfig{
			StaticCollectionConfig: &common.StaticCollectionConfig{
				Name:              collectionName,
				MemberOrgsPolicy:  accessPolicy,
				RequiredPeerCount: requiredPeerCount,
				MaximumPeerCount:  maximumPeerCount,
			},
		},
	}
}

func sampleCollectionConfigPackage(colName string) *common.CollectionConfig {
	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)

	var requiredPeerCount, maximumPeerCount int32
	requiredPeerCount = 1
	maximumPeerCount = 2

	return createCollectionConfig(colName, policyEnvelope, requiredPeerCount, maximumPeerCount)
}



func (s *store) persistOldProto(txid string, blockHeight uint64,
	privateSimulationResults *rwset.TxPvtReadWriteSet) error {

	logger.Debugf("Persisting private data to transient store for txid [%s] at block height [%d]", txid, blockHeight)

	dbBatch := leveldbhelper.NewUpdateBatch()

	
	
	
	uuid := commonutil.GenerateUUID()
	compositeKeyPvtRWSet := createCompositeKeyForPvtRWSet(txid, uuid, blockHeight)
	privateSimulationResultsBytes, err := proto.Marshal(privateSimulationResults)
	if err != nil {
		return err
	}
	dbBatch.Put(compositeKeyPvtRWSet, privateSimulationResultsBytes)

	

	
	
	
	
	
	compositeKeyPurgeIndexByHeight := createCompositeKeyForPurgeIndexByHeight(blockHeight, txid, uuid)
	dbBatch.Put(compositeKeyPurgeIndexByHeight, emptyValue)

	
	
	
	
	
	
	
	
	
	
	compositeKeyPurgeIndexByTxid := createCompositeKeyForPurgeIndexByTxid(txid, uuid, blockHeight)
	dbBatch.Put(compositeKeyPurgeIndexByTxid, emptyValue)

	return s.db.WriteBatch(dbBatch, true)
}
