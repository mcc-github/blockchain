/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	pb "github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/factory"
	util2 "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/transientstore"
	privdatacommon "github.com/mcc-github/blockchain/gossip/privdata/common"
	"github.com/mcc-github/blockchain/gossip/privdata/mocks"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/mcc-github/blockchain/protos/common"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protos/peer"
	transientstore2 "github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	viper.Set("peer.gossip.pvtData.pullRetryThreshold", time.Second*3)
	factory.InitFactories(nil)
}



type CollectionCriteria struct {
	Channel    string
	TxId       string
	Collection string
	Namespace  string
}

func fromCollectionCriteria(criteria common.CollectionCriteria) CollectionCriteria {
	return CollectionCriteria{
		TxId:       criteria.TxId,
		Collection: criteria.Collection,
		Namespace:  criteria.Namespace,
		Channel:    criteria.Channel,
	}
}

type persistCall struct {
	*mock.Call
	store *mockTransientStore
}

func (pc *persistCall) expectRWSet(namespace string, collection string, rws []byte) *persistCall {
	if pc.store.persists == nil {
		pc.store.persists = make(map[rwsTriplet]struct{})
	}
	pc.store.persists[rwsTriplet{
		namespace:  namespace,
		collection: collection,
		rwset:      hex.EncodeToString(rws),
	}] = struct{}{}
	return pc
}

type mockTransientStore struct {
	t *testing.T
	mock.Mock
	persists      map[rwsTriplet]struct{}
	lastReqTxID   string
	lastReqFilter map[string]ledger.PvtCollFilter
}

func (store *mockTransientStore) On(methodName string, arguments ...interface{}) *persistCall {
	return &persistCall{
		store: store,
		Call:  store.Mock.On(methodName, arguments...),
	}
}

func (store *mockTransientStore) PurgeByTxids(txids []string) error {
	args := store.Called(txids)
	return args.Error(0)
}

func (store *mockTransientStore) Persist(txid string, blockHeight uint64, res *rwset.TxPvtReadWriteSet) error {
	key := rwsTriplet{
		namespace:  res.NsPvtRwset[0].Namespace,
		collection: res.NsPvtRwset[0].CollectionPvtRwset[0].CollectionName,
		rwset:      hex.EncodeToString(res.NsPvtRwset[0].CollectionPvtRwset[0].Rwset)}
	if _, exists := store.persists[key]; !exists {
		store.t.Fatal("Shouldn't have persisted", res)
	}
	delete(store.persists, key)
	store.Called(txid, blockHeight, res)
	return nil
}

func (store *mockTransientStore) PersistWithConfig(txid string, blockHeight uint64, privateSimulationResultsWithConfig *transientstore2.TxPvtReadWriteSetWithConfigInfo) error {
	res := privateSimulationResultsWithConfig.PvtRwset
	key := rwsTriplet{
		namespace:  res.NsPvtRwset[0].Namespace,
		collection: res.NsPvtRwset[0].CollectionPvtRwset[0].CollectionName,
		rwset:      hex.EncodeToString(res.NsPvtRwset[0].CollectionPvtRwset[0].Rwset)}
	if _, exists := store.persists[key]; !exists {
		store.t.Fatal("Shouldn't have persisted", res)
	}
	delete(store.persists, key)
	store.Called(txid, blockHeight, privateSimulationResultsWithConfig)
	return nil
}

func (store *mockTransientStore) PurgeByHeight(maxBlockNumToRetain uint64) error {
	return store.Called(maxBlockNumToRetain).Error(0)
}

func (store *mockTransientStore) GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error) {
	store.lastReqTxID = txid
	store.lastReqFilter = filter
	args := store.Called(txid, filter)
	if args.Get(1) == nil {
		return args.Get(0).(transientstore.RWSetScanner), nil
	}
	return nil, args.Get(1).(error)
}

type mockRWSetScanner struct {
	err     error
	results []*transientstore.EndorserPvtSimulationResultsWithConfig
}

func (scanner *mockRWSetScanner) withRWSet(ns string, col string) *mockRWSetScanner {
	scanner.results = append(scanner.results, &transientstore.EndorserPvtSimulationResultsWithConfig{
		PvtSimulationResultsWithConfig: &transientstore2.TxPvtReadWriteSetWithConfigInfo{
			PvtRwset: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: ns,
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: col,
								Rwset:          []byte("rws-pre-image"),
							},
						},
					},
				},
			},
			CollectionConfigs: map[string]*common.CollectionConfigPackage{
				ns: {
					Config: []*common.CollectionConfig{
						{
							Payload: &common.CollectionConfig_StaticCollectionConfig{
								StaticCollectionConfig: &common.StaticCollectionConfig{
									Name: col,
								},
							},
						},
					},
				},
			},
		},
	})
	return scanner
}

func (scanner *mockRWSetScanner) Next() (*transientstore.EndorserPvtSimulationResults, error) {
	panic("should not be used")
}

func (scanner *mockRWSetScanner) NextWithConfig() (*transientstore.EndorserPvtSimulationResultsWithConfig, error) {
	if scanner.err != nil {
		return nil, scanner.err
	}
	var res *transientstore.EndorserPvtSimulationResultsWithConfig
	if len(scanner.results) == 0 {
		return nil, nil
	}
	res, scanner.results = scanner.results[len(scanner.results)-1], scanner.results[:len(scanner.results)-1]
	return res, nil
}

func (*mockRWSetScanner) Close() {
}

type committerMock struct {
	mock.Mock
}

func (mock *committerMock) GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error) {
	args := mock.Called()
	return args.Get(0).(ledger.MissingPvtDataTracker), args.Error(1)
}

func (mock *committerMock) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	args := mock.Called()
	return args.Get(0).(ledger.ConfigHistoryRetriever), args.Error(1)
}

func (mock *committerMock) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	args := mock.Called(blockNum, filter)
	return args.Get(0).([]*ledger.TxPvtData), args.Error(1)
}

func (mock *committerMock) CommitWithPvtData(blockAndPvtData *ledger.BlockAndPvtData) error {
	args := mock.Called(blockAndPvtData)
	return args.Error(0)
}

func (mock *committerMock) CommitPvtDataOfOldBlocks(blockPvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error) {
	args := mock.Called(blockPvtData)
	return args.Get(0).([]*ledger.PvtdataHashMismatch), args.Error(1)
}

func (mock *committerMock) GetPvtDataAndBlockByNum(seqNum uint64) (*ledger.BlockAndPvtData, error) {
	args := mock.Called(seqNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ledger.BlockAndPvtData), args.Error(1)
}

func (mock *committerMock) Commit(block *common.Block) error {
	args := mock.Called(block)
	return args.Error(0)
}

func (mock *committerMock) LedgerHeight() (uint64, error) {
	args := mock.Called()
	if args.Get(0) == nil {
		return uint64(0), args.Error(1)
	}
	return args.Get(0).(uint64), args.Error(1)
}

func (mock *committerMock) GetBlocks(blockSeqs []uint64) []*common.Block {
	args := mock.Called(blockSeqs)
	seqs := args.Get(0)
	if seqs == nil {
		return nil
	}
	return seqs.([]*common.Block)
}

func (mock *committerMock) Close() {
	mock.Called()
}

type validatorMock struct {
	err error
}

func (v *validatorMock) Validate(block *common.Block) error {
	if v.err != nil {
		return v.err
	}
	return nil
}

type digests []privdatacommon.DigKey

func (d digests) Equal(other digests) bool {
	flatten := func(d digests) map[privdatacommon.DigKey]struct{} {
		m := map[privdatacommon.DigKey]struct{}{}
		for _, dig := range d {
			m[dig] = struct{}{}
		}
		return m
	}
	return reflect.DeepEqual(flatten(d), flatten(other))
}

type fetchCall struct {
	fetcher *fetcherMock
	*mock.Call
}

func (fc *fetchCall) expectingEndorsers(orgs ...string) *fetchCall {
	if fc.fetcher.expectedEndorsers == nil {
		fc.fetcher.expectedEndorsers = make(map[string]struct{})
	}
	for _, org := range orgs {
		sID := &msp.SerializedIdentity{Mspid: org, IdBytes: []byte(fmt.Sprintf("p0%s", org))}
		b, _ := pb.Marshal(sID)
		fc.fetcher.expectedEndorsers[string(b)] = struct{}{}
	}

	return fc
}

func (fc *fetchCall) expectingDigests(digests []privdatacommon.DigKey) *fetchCall {
	fc.fetcher.expectedDigests = digests
	return fc
}

func (fc *fetchCall) Return(returnArguments ...interface{}) *mock.Call {

	return fc.Call.Return(returnArguments...)
}

type fetcherMock struct {
	t *testing.T
	mock.Mock
	expectedDigests   []privdatacommon.DigKey
	expectedEndorsers map[string]struct{}
}

func (f *fetcherMock) On(methodName string, arguments ...interface{}) *fetchCall {
	return &fetchCall{
		fetcher: f,
		Call:    f.Mock.On(methodName, arguments...),
	}
}

func (f *fetcherMock) fetch(dig2src dig2sources) (*privdatacommon.FetchedPvtDataContainer, error) {
	for _, endorsements := range dig2src {
		for _, endorsement := range endorsements {
			_, exists := f.expectedEndorsers[string(endorsement.Endorser)]
			if !exists {
				f.t.Fatalf("Encountered a non-expected endorser: %s", string(endorsement.Endorser))
			}
			
			delete(f.expectedEndorsers, string(endorsement.Endorser))
		}
	}
	assert.True(f.t, digests(dig2src.keys()).Equal(digests(f.expectedDigests)))
	assert.Empty(f.t, f.expectedEndorsers)
	args := f.Called(dig2src)
	if args.Get(1) == nil {
		return args.Get(0).(*privdatacommon.FetchedPvtDataContainer), nil
	}
	return nil, args.Get(1).(error)
}

func createcollectionStore(expectedSignedData common.SignedData) *collectionStore {
	return &collectionStore{
		expectedSignedData: expectedSignedData,
		policies:           make(map[collectionAccessPolicy]CollectionCriteria),
		store:              make(map[CollectionCriteria]collectionAccessPolicy),
	}
}

type collectionStore struct {
	expectedSignedData common.SignedData
	acceptsAll         bool
	acceptsNone        bool
	lenient            bool
	store              map[CollectionCriteria]collectionAccessPolicy
	policies           map[collectionAccessPolicy]CollectionCriteria
}

func (cs *collectionStore) thatAcceptsAll() *collectionStore {
	cs.acceptsAll = true
	return cs
}

func (cs *collectionStore) thatAcceptsNone() *collectionStore {
	cs.acceptsNone = true
	return cs
}

func (cs *collectionStore) andIsLenient() *collectionStore {
	cs.lenient = true
	return cs
}

func (cs *collectionStore) thatAccepts(cc CollectionCriteria) *collectionStore {
	sp := collectionAccessPolicy{
		cs: cs,
		n:  util.RandomUInt64(),
	}
	cs.store[cc] = sp
	cs.policies[sp] = cc
	return cs
}

func (cs *collectionStore) RetrieveCollectionAccessPolicy(cc common.CollectionCriteria) (privdata.CollectionAccessPolicy, error) {
	if sp, exists := cs.store[fromCollectionCriteria(cc)]; exists {
		return &sp, nil
	}
	if cs.acceptsAll || cs.acceptsNone || cs.lenient {
		return &collectionAccessPolicy{
			cs: cs,
			n:  util.RandomUInt64(),
		}, nil
	}
	return nil, privdata.NoSuchCollectionError{}
}

func (cs *collectionStore) RetrieveCollection(common.CollectionCriteria) (privdata.Collection, error) {
	panic("implement me")
}

func (cs *collectionStore) HasReadAccess(cc common.CollectionCriteria, sp *peer.SignedProposal, qe ledger.QueryExecutor) (bool, error) {
	panic("implement me")
}

func (cs *collectionStore) RetrieveCollectionConfigPackage(cc common.CollectionCriteria) (*common.CollectionConfigPackage, error) {
	return &common.CollectionConfigPackage{
		Config: []*common.CollectionConfig{
			{
				Payload: &common.CollectionConfig_StaticCollectionConfig{
					StaticCollectionConfig: &common.StaticCollectionConfig{
						Name:              cc.Collection,
						MaximumPeerCount:  1,
						RequiredPeerCount: 1,
					},
				},
			},
		},
	}, nil
}

func (cs *collectionStore) RetrieveCollectionPersistenceConfigs(cc common.CollectionCriteria) (privdata.CollectionPersistenceConfigs, error) {
	panic("implement me")
}

func (cs *collectionStore) AccessFilter(channelName string, collectionPolicyConfig *common.CollectionPolicyConfig) (privdata.Filter, error) {
	panic("implement me")
}

type collectionAccessPolicy struct {
	cs *collectionStore
	n  uint64
}

func (cap *collectionAccessPolicy) MemberOrgs() []string {
	return []string{"org0", "org1"}
}

func (cap *collectionAccessPolicy) RequiredPeerCount() int {
	return 1
}

func (cap *collectionAccessPolicy) MaximumPeerCount() int {
	return 2
}

func (cap *collectionAccessPolicy) IsMemberOnlyRead() bool {
	return false
}

func (cap *collectionAccessPolicy) AccessFilter() privdata.Filter {
	return func(sd common.SignedData) bool {
		that, _ := asn1.Marshal(sd)
		this, _ := asn1.Marshal(cap.cs.expectedSignedData)
		if hex.EncodeToString(that) != hex.EncodeToString(this) {
			panic(fmt.Errorf("self signed data passed isn't equal to expected:%v, %v", sd, cap.cs.expectedSignedData))
		}

		if cap.cs.acceptsNone {
			return false
		} else if cap.cs.acceptsAll {
			return true
		}

		_, exists := cap.cs.policies[*cap]
		return exists
	}
}

func TestPvtDataCollections_FailOnEmptyPayload(t *testing.T) {
	collection := &util.PvtDataCollections{
		&ledger.TxPvtData{
			SeqInBlock: uint64(1),
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: "ns1",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "secretCollection",
								Rwset:          []byte{1, 2, 3, 4, 5, 6, 7},
							},
						},
					},
				},
			},
		},

		nil,
	}

	_, err := collection.Marshal()
	assertion := assert.New(t)
	assertion.Error(err, "Expected to fail since second item has nil payload")
	assertion.Equal("Mallformed private data payload, rwset index 1 is nil", fmt.Sprintf("%s", err))
}

func TestPvtDataCollections_FailMarshalingWriteSet(t *testing.T) {
	collection := &util.PvtDataCollections{
		&ledger.TxPvtData{
			SeqInBlock: uint64(1),
			WriteSet:   nil,
		},
	}

	_, err := collection.Marshal()
	assertion := assert.New(t)
	assertion.Error(err, "Expected to fail since first item has nil writeset")
	assertion.Contains(fmt.Sprintf("%s", err), "Could not marshal private rwset index 0")
}

func TestPvtDataCollections_Marshal(t *testing.T) {
	collection := &util.PvtDataCollections{
		&ledger.TxPvtData{
			SeqInBlock: uint64(1),
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: "ns1",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "secretCollection",
								Rwset:          []byte{1, 2, 3, 4, 5, 6, 7},
							},
						},
					},
				},
			},
		},

		&ledger.TxPvtData{
			SeqInBlock: uint64(2),
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: "ns1",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "secretCollection",
								Rwset:          []byte{42, 42, 42, 42, 42, 42, 42},
							},
						},
					},
					{
						Namespace: "ns2",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "otherCollection",
								Rwset:          []byte{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
							},
						},
					},
				},
			},
		},
	}

	bytes, err := collection.Marshal()

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotNil(bytes)
	assertion.Equal(2, len(bytes))
}

func TestPvtDataCollections_Unmarshal(t *testing.T) {
	collection := util.PvtDataCollections{
		&ledger.TxPvtData{
			SeqInBlock: uint64(1),
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: "ns1",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "secretCollection",
								Rwset:          []byte{1, 2, 3, 4, 5, 6, 7},
							},
						},
					},
				},
			},
		},
	}

	bytes, err := collection.Marshal()

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotNil(bytes)
	assertion.Equal(1, len(bytes))

	var newCol util.PvtDataCollections

	err = newCol.Unmarshal(bytes)
	assertion.NoError(err)
	assertion.Equal(1, len(newCol))
	assertion.Equal(newCol[0].SeqInBlock, collection[0].SeqInBlock)
	assertion.True(pb.Equal(newCol[0].WriteSet, collection[0].WriteSet))
}

type rwsTriplet struct {
	namespace  string
	collection string
	rwset      string
}

func flattenTxPvtDataMap(pd ledger.TxPvtDataMap) map[uint64]map[rwsTriplet]struct{} {
	m := make(map[uint64]map[rwsTriplet]struct{})
	for seqInBlock, namespaces := range pd {
		triplets := make(map[rwsTriplet]struct{})
		for _, namespace := range namespaces.WriteSet.NsPvtRwset {
			for _, col := range namespace.CollectionPvtRwset {
				triplets[rwsTriplet{
					namespace:  namespace.Namespace,
					collection: col.CollectionName,
					rwset:      hex.EncodeToString(col.Rwset),
				}] = struct{}{}
			}
		}
		m[seqInBlock] = triplets
	}
	return m
}

var expectedCommittedPrivateData1 = map[uint64]*ledger.TxPvtData{
	0: {SeqInBlock: 0, WriteSet: &rwset.TxPvtReadWriteSet{
		DataModel: rwset.TxReadWriteSet_KV,
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{
			{
				Namespace: "ns1",
				CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
					{
						CollectionName: "c1",
						Rwset:          []byte("rws-pre-image"),
					},
					{
						CollectionName: "c2",
						Rwset:          []byte("rws-pre-image"),
					},
				},
			},
		},
	}},
	1: {SeqInBlock: 1, WriteSet: &rwset.TxPvtReadWriteSet{
		DataModel: rwset.TxReadWriteSet_KV,
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{
			{
				Namespace: "ns2",
				CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
					{
						CollectionName: "c1",
						Rwset:          []byte("rws-pre-image"),
					},
				},
			},
		},
	}},
}

var expectedCommittedPrivateData2 = map[uint64]*ledger.TxPvtData{
	0: {SeqInBlock: 0, WriteSet: &rwset.TxPvtReadWriteSet{
		DataModel: rwset.TxReadWriteSet_KV,
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{
			{
				Namespace: "ns3",
				CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
					{
						CollectionName: "c3",
						Rwset:          []byte("rws-pre-image"),
					},
				},
			},
		},
	}},
}

var expectedCommittedPrivateData3 = map[uint64]*ledger.TxPvtData{}

func TestCoordinatorStoreInvalidBlock(t *testing.T) {
	peerSelfSignedData := common.SignedData{
		Identity:  []byte{0, 1, 2},
		Signature: []byte{3, 4, 5},
		Data:      []byte{6, 7, 8},
	}
	hash := util2.ComputeSHA256([]byte("rws-pre-image"))
	committer := &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		t.Fatal("Shouldn't have committed")
	}).Return(nil)
	cs := createcollectionStore(peerSelfSignedData).thatAcceptsAll()
	purgedTxns := make(map[string]struct{})
	store := &mockTransientStore{t: t}
	assertPurged := func(txns ...string) {
		for _, txn := range txns {
			_, exists := purgedTxns[txn]
			assert.True(t, exists)
			delete(purgedTxns, txn)
		}
		assert.Len(t, purgedTxns, 0)
	}
	store.On("PurgeByTxids", mock.Anything).Run(func(args mock.Arguments) {
		for _, txn := range args.Get(0).([]string) {
			purgedTxns[txn] = struct{}{}
		}
	}).Return(nil)
	fetcher := &fetcherMock{t: t}
	pdFactory := &pvtDataFactory{}
	bf := &blockFactory{
		channelID: "test",
	}

	block := bf.withoutMetadata().create()
	
	pvtData := pdFactory.create()
	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	err := coordinator.StoreBlock(block, pvtData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Block.Metadata is nil or Block.Metadata lacks a Tx filter bitmap")

	
	block = bf.create()
	pvtData = pdFactory.create()
	coordinator = NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{fmt.Errorf("failed validating block")},
	}, peerSelfSignedData)
	err = coordinator.StoreBlock(block, pvtData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed validating block")

	
	block = bf.withMetadataSize(100).create()
	pvtData = pdFactory.create()
	coordinator = NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	err = coordinator.StoreBlock(block, pvtData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Block data size")
	assert.Contains(t, err.Error(), "is different from Tx filter size")

	
	
	
	var commitHappened bool
	assertCommitHappened := func() {
		assert.True(t, commitHappened)
		commitHappened = false
	}
	committer = &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		privateDataPassed2Ledger := args.Get(0).(*ledger.BlockAndPvtData).PvtData
		commitHappened = true
		
		assert.Len(t, privateDataPassed2Ledger, 1)
		assert.Equal(t, 0, int(privateDataPassed2Ledger[0].SeqInBlock))
		
		assert.Len(t, privateDataPassed2Ledger[0].WriteSet.NsPvtRwset, 1)
		assert.Equal(t, "ns1", privateDataPassed2Ledger[0].WriteSet.NsPvtRwset[0].Namespace)
		assert.Len(t, privateDataPassed2Ledger[0].WriteSet.NsPvtRwset[0].CollectionPvtRwset, 2)
	}).Return(nil)
	block = bf.withInvalidTxns(1).AddTxn("tx1", "ns1", hash, "c1", "c2").AddTxn("tx2", "ns2", hash, "c1").create()
	pvtData = pdFactory.addRWSet().addNSRWSet("ns1", "c1", "c2").create()
	coordinator = NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	err = coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()
	
	
	
	assertPurged("tx1", "tx2")

	
	block.Header = nil
	err = coordinator.StoreBlock(block, pvtData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Block header is nil")

	
	block.Data = nil
	err = coordinator.StoreBlock(block, pvtData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Block data is empty")
}

func TestCoordinatorToFilterOutPvtRWSetsWithWrongHash(t *testing.T) {
	
	peerSelfSignedData := common.SignedData{
		Identity:  []byte{0, 1, 2},
		Signature: []byte{3, 4, 5},
		Data:      []byte{6, 7, 8},
	}

	expectedPvtData := map[uint64]*ledger.TxPvtData{
		0: {SeqInBlock: 0, WriteSet: &rwset.TxPvtReadWriteSet{
			DataModel: rwset.TxReadWriteSet_KV,
			NsPvtRwset: []*rwset.NsPvtReadWriteSet{
				{
					Namespace: "ns1",
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
						{
							CollectionName: "c1",
							Rwset:          []byte("rws-original"),
						},
					},
				},
			},
		}},
	}

	cs := createcollectionStore(peerSelfSignedData).thatAcceptsAll()
	committer := &committerMock{}
	store := &mockTransientStore{t: t}

	fetcher := &fetcherMock{t: t}

	var commitHappened bool

	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		privateDataPassed2Ledger := args.Get(0).(*ledger.BlockAndPvtData).PvtData
		assert.True(t, reflect.DeepEqual(flattenTxPvtDataMap(privateDataPassed2Ledger),
			flattenTxPvtDataMap(expectedPvtData)))
		commitHappened = true
	}).Return(nil)

	hash := util2.ComputeSHA256([]byte("rws-original"))
	bf := &blockFactory{
		channelID: "test",
	}

	block := bf.AddTxnWithEndorsement("tx1", "ns1", hash, "org1", true, "c1").create()
	store.On("GetTxPvtRWSetByTxid", "tx1", mock.Anything).Return((&mockRWSetScanner{}).withRWSet("ns1", "c1"), nil)

	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)

	fetcher.On("fetch", mock.Anything).expectingDigests([]privdatacommon.DigKey{
		{
			TxId: "tx1", Namespace: "ns1", Collection: "c1", BlockSeq: 1,
		},
	}).expectingEndorsers("org1").Return(&privdatacommon.FetchedPvtDataContainer{
		AvailableElements: []*proto.PvtDataElement{
			{
				Digest: &proto.PvtDataDigest{
					BlockSeq:   1,
					Collection: "c1",
					Namespace:  "ns1",
					TxId:       "tx1",
				},
				Payload: [][]byte{[]byte("rws-original")},
			},
		},
	}, nil)
	store.On("Persist", mock.Anything, uint64(1), mock.Anything).
		expectRWSet("ns1", "c1", []byte("rws-original")).Return(nil)

	purgedTxns := make(map[string]struct{})
	store.On("PurgeByTxids", mock.Anything).Run(func(args mock.Arguments) {
		for _, txn := range args.Get(0).([]string) {
			purgedTxns[txn] = struct{}{}
		}
	}).Return(nil)

	coordinator.StoreBlock(block, nil)
	
	assert.True(t, commitHappened)

	
	_, exists := purgedTxns["tx1"]
	assert.True(t, exists)
}

func TestCoordinatorStoreBlock(t *testing.T) {
	peerSelfSignedData := common.SignedData{
		Identity:  []byte{0, 1, 2},
		Signature: []byte{3, 4, 5},
		Data:      []byte{6, 7, 8},
	}
	

	cs := createcollectionStore(peerSelfSignedData).thatAcceptsAll()

	var commitHappened bool
	assertCommitHappened := func() {
		assert.True(t, commitHappened)
		commitHappened = false
	}
	committer := &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		privateDataPassed2Ledger := args.Get(0).(*ledger.BlockAndPvtData).PvtData
		assert.True(t, reflect.DeepEqual(flattenTxPvtDataMap(privateDataPassed2Ledger),
			flattenTxPvtDataMap(expectedCommittedPrivateData1)))
		commitHappened = true
	}).Return(nil)
	purgedTxns := make(map[string]struct{})
	store := &mockTransientStore{t: t}
	store.On("PurgeByTxids", mock.Anything).Run(func(args mock.Arguments) {
		for _, txn := range args.Get(0).([]string) {
			purgedTxns[txn] = struct{}{}
		}
	}).Return(nil)
	assertPurged := func(txns ...string) {
		for _, txn := range txns {
			_, exists := purgedTxns[txn]
			assert.True(t, exists)
			delete(purgedTxns, txn)
		}
		assert.Len(t, purgedTxns, 0)
	}

	fetcher := &fetcherMock{t: t}

	hash := util2.ComputeSHA256([]byte("rws-pre-image"))
	pdFactory := &pvtDataFactory{}
	bf := &blockFactory{
		channelID: "test",
	}
	block := bf.AddTxnWithEndorsement("tx1", "ns1", hash, "org1", true, "c1", "c2").
		AddTxnWithEndorsement("tx2", "ns2", hash, "org2", true, "c1").create()

	fmt.Println("Scenario I")
	
	
	
	pvtData := pdFactory.addRWSet().addNSRWSet("ns1", "c1", "c2").addRWSet().addNSRWSet("ns2", "c1").create()
	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	err := coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()
	assertPurged("tx1", "tx2")

	fmt.Println("Scenario II")
	
	
	store.On("GetTxPvtRWSetByTxid", "tx1", mock.Anything).Return((&mockRWSetScanner{}).withRWSet("ns1", "c2"), nil)
	store.On("GetTxPvtRWSetByTxid", "tx2", mock.Anything).Return(&mockRWSetScanner{}, nil)
	pvtData = pdFactory.addRWSet().addNSRWSet("ns1", "c1").addRWSet().addNSRWSet("ns2", "c1").create()
	err = coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()
	assertPurged("tx1", "tx2")
	assert.Equal(t, "tx1", store.lastReqTxID)
	assert.Equal(t, map[string]ledger.PvtCollFilter{
		"ns1": map[string]bool{
			"c2": true,
		},
	}, store.lastReqFilter)

	fmt.Println("Scenario III")
	
	
	
	
	
	fetcher.On("fetch", mock.Anything).expectingDigests([]privdatacommon.DigKey{
		{
			TxId: "tx1", Namespace: "ns1", Collection: "c2", BlockSeq: 1,
		},
		{
			TxId: "tx2", Namespace: "ns2", Collection: "c1", BlockSeq: 1, SeqInBlock: 1,
		},
	}).expectingEndorsers("org1").Return(&privdatacommon.FetchedPvtDataContainer{
		AvailableElements: []*proto.PvtDataElement{
			{
				Digest: &proto.PvtDataDigest{
					BlockSeq:   1,
					Collection: "c2",
					Namespace:  "ns1",
					TxId:       "tx1",
				},
				Payload: [][]byte{[]byte("rws-pre-image")},
			},
			{
				Digest: &proto.PvtDataDigest{
					SeqInBlock: 1,
					BlockSeq:   1,
					Collection: "c1",
					Namespace:  "ns2",
					TxId:       "tx2",
				},
				Payload: [][]byte{[]byte("rws-pre-image")},
			},
		},
	}, nil)
	store.On("Persist", mock.Anything, uint64(1), mock.Anything).
		expectRWSet("ns1", "c2", []byte("rws-pre-image")).
		expectRWSet("ns2", "c1", []byte("rws-pre-image")).Return(nil)
	pvtData = pdFactory.addRWSet().addNSRWSet("ns1", "c1").create()
	err = coordinator.StoreBlock(block, pvtData)
	assertPurged("tx1", "tx2")
	assert.NoError(t, err)
	assertCommitHappened()

	fmt.Println("Scenario IV")
	
	pvtData = pdFactory.addRWSet().addNSRWSet("ns1", "c1", "c2", "c3").
		addRWSet().addNSRWSet("ns2", "c1", "c3").addRWSet().addNSRWSet("ns1", "c4").create()
	err = coordinator.StoreBlock(block, pvtData)
	assertPurged("tx1", "tx2")
	assert.NoError(t, err)
	assertCommitHappened()

	fmt.Println("Scenario V")
	
	
	
	mockCs := &mocks.CollectionStore{}
	mockCs.On("RetrieveCollectionAccessPolicy", mock.Anything).Return(nil, errors.New("test error"))
	coordinator = NewCoordinator(Support{
		CollectionStore: mockCs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	err = coordinator.StoreBlock(block, nil)
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	fmt.Println("Scenario VI")
	
	
	
	block = bf.AddTxn("tx3", "ns3", hash, "c3").create()
	fetcher = &fetcherMock{t: t}
	fetcher.On("fetch", mock.Anything).expectingDigests([]privdatacommon.DigKey{
		{
			TxId: "tx3", Namespace: "ns3", Collection: "c3", BlockSeq: 1,
		},
	}).Return(&privdatacommon.FetchedPvtDataContainer{
		AvailableElements: []*proto.PvtDataElement{
			{
				Digest: &proto.PvtDataDigest{
					BlockSeq:   1,
					Collection: "c3",
					Namespace:  "ns3",
					TxId:       "tx3",
				},
				Payload: [][]byte{[]byte("rws-pre-image")},
			},
		},
	}, nil)
	store = &mockTransientStore{t: t}
	store.On("PurgeByTxids", mock.Anything).Run(func(args mock.Arguments) {
		for _, txn := range args.Get(0).([]string) {
			purgedTxns[txn] = struct{}{}
		}
	}).Return(nil)
	store.On("GetTxPvtRWSetByTxid", "tx3", mock.Anything).Return(&mockRWSetScanner{err: errors.New("uh oh")}, nil)
	store.On("Persist", mock.Anything, uint64(1), mock.Anything).expectRWSet("ns3", "c3", []byte("rws-pre-image"))
	committer = &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		privateDataPassed2Ledger := args.Get(0).(*ledger.BlockAndPvtData).PvtData
		assert.True(t, reflect.DeepEqual(flattenTxPvtDataMap(privateDataPassed2Ledger),
			flattenTxPvtDataMap(expectedCommittedPrivateData2)))
		commitHappened = true
	}).Return(nil)
	coordinator = NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	err = coordinator.StoreBlock(block, nil)
	assertPurged("tx3")
	assert.NoError(t, err)
	assertCommitHappened()

	fmt.Println("Scenario VII")
	
	
	
	
	block = bf.AddTxn("tx3", "ns3", hash, "c3", "c2", "c1").AddTxn("tx1", "ns1", hash, "c1").create()
	cs = createcollectionStore(peerSelfSignedData).thatAccepts(CollectionCriteria{
		TxId:       "tx3",
		Collection: "c3",
		Namespace:  "ns3",
		Channel:    "test",
	})
	store = &mockTransientStore{t: t}
	store.On("PurgeByTxids", mock.Anything).Run(func(args mock.Arguments) {
		for _, txn := range args.Get(0).([]string) {
			purgedTxns[txn] = struct{}{}
		}
	}).Return(nil)
	fetcher = &fetcherMock{t: t}
	committer = &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		privateDataPassed2Ledger := args.Get(0).(*ledger.BlockAndPvtData).PvtData
		assert.True(t, reflect.DeepEqual(flattenTxPvtDataMap(privateDataPassed2Ledger),
			flattenTxPvtDataMap(expectedCommittedPrivateData2)))
		commitHappened = true
	}).Return(nil)
	coordinator = NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)

	pvtData = pdFactory.addRWSet().addNSRWSet("ns3", "c3").create()
	err = coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()
	
	assertPurged("tx3", "tx1")
}

func TestProceedWithoutPrivateData(t *testing.T) {
	
	
	peerSelfSignedData := common.SignedData{
		Identity:  []byte{0, 1, 2},
		Signature: []byte{3, 4, 5},
		Data:      []byte{6, 7, 8},
	}
	cs := createcollectionStore(peerSelfSignedData).thatAcceptsAll()
	var commitHappened bool
	assertCommitHappened := func() {
		assert.True(t, commitHappened)
		commitHappened = false
	}
	committer := &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		blockAndPrivateData := args.Get(0).(*ledger.BlockAndPvtData)
		privateDataPassed2Ledger := blockAndPrivateData.PvtData
		assert.True(t, reflect.DeepEqual(flattenTxPvtDataMap(privateDataPassed2Ledger),
			flattenTxPvtDataMap(expectedCommittedPrivateData2)))
		missingPrivateData := blockAndPrivateData.MissingPvtData
		expectedMissingPvtData := make(ledger.TxMissingPvtDataMap)
		expectedMissingPvtData.Add(0, "ns3", "c2", true)
		assert.Equal(t, expectedMissingPvtData, missingPrivateData)
		commitHappened = true
	}).Return(nil)
	purgedTxns := make(map[string]struct{})
	store := &mockTransientStore{t: t}
	store.On("GetTxPvtRWSetByTxid", "tx1", mock.Anything).Return(&mockRWSetScanner{}, nil)
	store.On("PurgeByTxids", mock.Anything).Run(func(args mock.Arguments) {
		for _, txn := range args.Get(0).([]string) {
			purgedTxns[txn] = struct{}{}
		}
	}).Return(nil)
	assertPurged := func(txns ...string) {
		for _, txn := range txns {
			_, exists := purgedTxns[txn]
			assert.True(t, exists)
			delete(purgedTxns, txn)
		}
		assert.Len(t, purgedTxns, 0)
	}

	fetcher := &fetcherMock{t: t}
	
	fetcher.On("fetch", mock.Anything).expectingDigests([]privdatacommon.DigKey{
		{
			TxId: "tx1", Namespace: "ns3", Collection: "c2", BlockSeq: 1,
		},
	}).Return(&privdatacommon.FetchedPvtDataContainer{
		AvailableElements: []*proto.PvtDataElement{
			{
				Digest: &proto.PvtDataDigest{
					BlockSeq:   1,
					Collection: "c2",
					Namespace:  "ns3",
					TxId:       "tx1",
				},
				Payload: [][]byte{[]byte("wrong pre-image")},
			},
		},
	}, nil)

	hash := util2.ComputeSHA256([]byte("rws-pre-image"))
	pdFactory := &pvtDataFactory{}
	bf := &blockFactory{
		channelID: "test",
	}

	block := bf.AddTxn("tx1", "ns3", hash, "c3", "c2").create()
	pvtData := pdFactory.addRWSet().addNSRWSet("ns3", "c3").create()
	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	err := coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()
	assertPurged("tx1")
}

func TestProceedWithInEligiblePrivateData(t *testing.T) {
	
	
	peerSelfSignedData := common.SignedData{
		Identity:  []byte{0, 1, 2},
		Signature: []byte{3, 4, 5},
		Data:      []byte{6, 7, 8},
	}

	cs := createcollectionStore(peerSelfSignedData).thatAcceptsNone()

	var commitHappened bool
	assertCommitHappened := func() {
		assert.True(t, commitHappened)
		commitHappened = false
	}
	committer := &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		blockAndPrivateData := args.Get(0).(*ledger.BlockAndPvtData)
		privateDataPassed2Ledger := blockAndPrivateData.PvtData
		assert.True(t, reflect.DeepEqual(flattenTxPvtDataMap(privateDataPassed2Ledger),
			flattenTxPvtDataMap(expectedCommittedPrivateData3)))
		missingPrivateData := blockAndPrivateData.MissingPvtData
		expectedMissingPvtData := make(ledger.TxMissingPvtDataMap)
		expectedMissingPvtData.Add(0, "ns3", "c2", false)
		assert.Equal(t, expectedMissingPvtData, missingPrivateData)
		commitHappened = true
	}).Return(nil)

	hash := util2.ComputeSHA256([]byte("rws-pre-image"))
	bf := &blockFactory{
		channelID: "test",
	}

	block := bf.AddTxn("tx1", "ns3", hash, "c2").create()

	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         nil,
		TransientStore:  nil,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	err := coordinator.StoreBlock(block, nil)
	assert.NoError(t, err)
	assertCommitHappened()
}

func TestCoordinatorGetBlocks(t *testing.T) {
	sd := common.SignedData{
		Identity:  []byte{0, 1, 2},
		Signature: []byte{3, 4, 5},
		Data:      []byte{6, 7, 8},
	}
	cs := createcollectionStore(sd).thatAcceptsAll()
	committer := &committerMock{}
	store := &mockTransientStore{t: t}
	fetcher := &fetcherMock{t: t}
	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, sd)

	hash := util2.ComputeSHA256([]byte("rws-pre-image"))
	bf := &blockFactory{
		channelID: "test",
	}
	block := bf.AddTxn("tx1", "ns1", hash, "c1", "c2").AddTxn("tx2", "ns2", hash, "c1").create()

	
	
	cs = createcollectionStore(sd).thatAccepts(CollectionCriteria{
		Namespace:  "ns1",
		Collection: "c2",
		TxId:       "tx1",
		Channel:    "test",
	})
	committer.Mock = mock.Mock{}
	committer.On("GetPvtDataAndBlockByNum", mock.Anything).Return(&ledger.BlockAndPvtData{
		Block:   block,
		PvtData: expectedCommittedPrivateData1,
	}, nil)
	coordinator = NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, sd)
	expectedPrivData := (&pvtDataFactory{}).addRWSet().addNSRWSet("ns1", "c2").create()
	block2, returnedPrivateData, err := coordinator.GetPvtDataAndBlockByNum(1, sd)
	assert.NoError(t, err)
	assert.Equal(t, block, block2)
	assert.Equal(t, expectedPrivData, []*ledger.TxPvtData(returnedPrivateData))

	
	committer.Mock = mock.Mock{}
	committer.On("GetPvtDataAndBlockByNum", mock.Anything).Return(nil, errors.New("uh oh"))
	block2, returnedPrivateData, err = coordinator.GetPvtDataAndBlockByNum(1, sd)
	assert.Nil(t, block2)
	assert.Empty(t, returnedPrivateData)
	assert.Error(t, err)
}

func TestPurgeByHeight(t *testing.T) {
	
	
	peerSelfSignedData := common.SignedData{}
	cs := createcollectionStore(peerSelfSignedData).thatAcceptsAll()

	var purgeHappened bool
	assertPurgeHappened := func() {
		assert.True(t, purgeHappened)
		purgeHappened = false
	}
	committer := &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Return(nil)
	store := &mockTransientStore{t: t}
	store.On("PurgeByHeight", uint64(1000)).Return(nil).Once().Run(func(_ mock.Arguments) {
		purgeHappened = true
	})
	store.On("PurgeByHeight", uint64(2000)).Return(nil).Once().Run(func(_ mock.Arguments) {
		purgeHappened = true
	})
	store.On("PurgeByTxids", mock.Anything).Return(nil)
	fetcher := &fetcherMock{t: t}

	bf := &blockFactory{
		channelID: "test",
	}
	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)

	for i := 0; i <= 3000; i++ {
		block := bf.create()
		block.Header.Number = uint64(i)
		err := coordinator.StoreBlock(block, nil)
		assert.NoError(t, err)
		if i != 2000 && i != 3000 {
			assert.False(t, purgeHappened)
		} else {
			assertPurgeHappened()
		}
	}
}

func TestCoordinatorStorePvtData(t *testing.T) {
	cs := createcollectionStore(common.SignedData{}).thatAcceptsAll()
	committer := &committerMock{}
	store := &mockTransientStore{t: t}
	store.On("PersistWithConfig", mock.Anything, uint64(5), mock.Anything).
		expectRWSet("ns1", "c1", []byte("rws-pre-image")).Return(nil)
	fetcher := &fetcherMock{t: t}
	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, common.SignedData{})
	pvtData := (&pvtDataFactory{}).addRWSet().addNSRWSet("ns1", "c1").create()
	
	err := coordinator.StorePvtData("tx1", &transientstore2.TxPvtReadWriteSetWithConfigInfo{
		PvtRwset:          pvtData[0].WriteSet,
		CollectionConfigs: make(map[string]*common.CollectionConfigPackage),
	}, uint64(5))
	assert.NoError(t, err)
}

func TestContainsWrites(t *testing.T) {
	
	col := &rwsetutil.CollHashedRwSet{
		CollectionName: "col1",
	}
	assert.False(t, containsWrites("tx", "ns", col))

	
	col.HashedRwSet = &kvrwset.HashedRWSet{}
	assert.False(t, containsWrites("tx", "ns", col))

	
	col.HashedRwSet.HashedWrites = append(col.HashedRwSet.HashedWrites, &kvrwset.KVWriteHash{})
	assert.True(t, containsWrites("tx", "ns", col))
}

func TestIgnoreReadOnlyColRWSets(t *testing.T) {
	
	
	
	
	
	
	
	
	peerSelfSignedData := common.SignedData{
		Identity:  []byte{0, 1, 2},
		Signature: []byte{3, 4, 5},
		Data:      []byte{6, 7, 8},
	}
	cs := createcollectionStore(peerSelfSignedData).thatAcceptsAll()
	var commitHappened bool
	assertCommitHappened := func() {
		assert.True(t, commitHappened)
		commitHappened = false
	}
	committer := &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		blockAndPrivateData := args.Get(0).(*ledger.BlockAndPvtData)
		
		assert.Empty(t, blockAndPrivateData.PvtData)
		
		assert.Empty(t, blockAndPrivateData.MissingPvtData)
		commitHappened = true
	}).Return(nil)
	store := &mockTransientStore{t: t}

	fetcher := &fetcherMock{t: t}
	hash := util2.ComputeSHA256([]byte("rws-pre-image"))
	bf := &blockFactory{
		channelID: "test",
	}
	
	block := bf.AddReadOnlyTxn("tx1", "ns3", hash, "c3", "c2").create()
	coordinator := NewCoordinator(Support{
		CollectionStore: cs,
		Committer:       committer,
		Fetcher:         fetcher,
		TransientStore:  store,
		Validator:       &validatorMock{},
	}, peerSelfSignedData)
	
	
	err := coordinator.StoreBlock(block, nil)
	assert.NoError(t, err)
	assertCommitHappened()
}
