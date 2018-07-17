/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"errors"
	"testing"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/protos/common"
	gossip2 "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	transientstore2 "github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedHistoryRetreiver struct {
	mock.Mock
}

func (mock *mockedHistoryRetreiver) CollectionConfigAt(blockNum uint64, chaincodeName string) (*ledger.CollectionConfigInfo, error) {
	args := mock.Called(blockNum, chaincodeName)
	return args.Get(0).(*ledger.CollectionConfigInfo), args.Error(1)
}

func (mock *mockedHistoryRetreiver) MostRecentCollectionConfigBelow(blockNum uint64, chaincodeName string) (*ledger.CollectionConfigInfo, error) {
	args := mock.Called(blockNum, chaincodeName)
	return args.Get(0).(*ledger.CollectionConfigInfo), args.Error(1)
}

type mockedDataStore struct {
	mock.Mock
}

func (ds *mockedDataStore) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	args := ds.Called()
	return args.Get(0).(ledger.ConfigHistoryRetriever), args.Error(1)
}

func (ds *mockedDataStore) GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error) {
	args := ds.Called(txid, filter)
	return args.Get(0).(transientstore.RWSetScanner), args.Error(1)
}

func (ds *mockedDataStore) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	args := ds.Called(blockNum, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ledger.TxPvtData), args.Error(1)
}

func (ds *mockedDataStore) LedgerHeight() (uint64, error) {
	args := ds.Called()
	return args.Get(0).(uint64), args.Error(1)
}

type mockedRWSetScanner struct {
	mock.Mock
}

func (mock *mockedRWSetScanner) Close() {

}

func (mock *mockedRWSetScanner) Next() (*transientstore.EndorserPvtSimulationResults, error) {
	args := mock.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transientstore.EndorserPvtSimulationResults), args.Error(1)
}

func (mock *mockedRWSetScanner) NextWithConfig() (*transientstore.EndorserPvtSimulationResultsWithConfig, error) {
	args := mock.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transientstore.EndorserPvtSimulationResultsWithConfig), args.Error(1)
}


func TestNewDataRetriever_GetDataFromTransientStore(t *testing.T) {
	t.Parallel()
	dataStore := &mockedDataStore{}

	rwSetScanner := &mockedRWSetScanner{}
	namespace := "testChaincodeName1"
	collectionName := "testCollectionName"

	rwSetScanner.On("NextWithConfig").Return(&transientstore.EndorserPvtSimulationResultsWithConfig{
		ReceivedAtBlockHeight:          2,
		PvtSimulationResultsWithConfig: nil,
	}, nil).Once().On("NextWithConfig").Return(&transientstore.EndorserPvtSimulationResultsWithConfig{
		ReceivedAtBlockHeight: 2,
		PvtSimulationResultsWithConfig: &transientstore2.TxPvtReadWriteSetWithConfigInfo{
			PvtRwset: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: namespace,
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: collectionName,
								Rwset:          []byte{1, 2},
							},
						},
					},
					{
						Namespace: namespace,
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: collectionName,
								Rwset:          []byte{3, 4},
							},
						},
					},
				},
			},
			CollectionConfigs: map[string]*common.CollectionConfigPackage{
				namespace: {
					Config: []*common.CollectionConfig{
						{
							Payload: &common.CollectionConfig_StaticCollectionConfig{
								StaticCollectionConfig: &common.StaticCollectionConfig{
									Name: collectionName,
								},
							},
						},
					},
				},
			},
		},
	}, nil).
		Once(). 
		On("NextWithConfig").Return(nil, nil)

	dataStore.On("LedgerHeight").Return(uint64(1), nil)
	dataStore.On("GetTxPvtRWSetByTxid", "testTxID", mock.Anything).Return(rwSetScanner, nil)

	retriever := NewDataRetriever(dataStore)

	
	
	rwSets, err := retriever.CollectionRWSet(&gossip2.PvtDataDigest{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   2,
		TxId:       "testTxID",
		SeqInBlock: 1,
	})

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotNil(rwSets)
	assertion.NotEmpty(rwSets.RWSet)
	assertion.Equal(2, len(rwSets.RWSet))

	var mergedRWSet []byte
	for _, rws := range rwSets.RWSet {
		mergedRWSet = append(mergedRWSet, rws...)
	}

	assertion.Equal([]byte{1, 2, 3, 4}, mergedRWSet)
}


func TestNewDataRetriever_GetDataFromLedger(t *testing.T) {
	t.Parallel()
	dataStore := &mockedDataStore{}

	namespace := "testChaincodeName1"
	collectionName := "testCollectionName"

	result := []*ledger.TxPvtData{{
		WriteSet: &rwset.TxPvtReadWriteSet{
			DataModel: rwset.TxReadWriteSet_KV,
			NsPvtRwset: []*rwset.NsPvtReadWriteSet{
				{
					Namespace: namespace,
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{{
						CollectionName: collectionName,
						Rwset:          []byte{1, 2},
					}},
				},
				{
					Namespace: namespace,
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{{
						CollectionName: collectionName,
						Rwset:          []byte{3, 4},
					}},
				},
			},
		},
		SeqInBlock: 1,
	}}

	dataStore.On("LedgerHeight").Return(uint64(10), nil)
	dataStore.On("GetPvtDataByNum", uint64(5), mock.Anything).Return(result, nil)

	historyRetreiver := &mockedHistoryRetreiver{}
	historyRetreiver.On("MostRecentCollectionConfigBelow", mock.Anything, namespace).Return(&ledger.CollectionConfigInfo{
		CollectionConfig: &common.CollectionConfigPackage{
			Config: []*common.CollectionConfig{
				{
					Payload: &common.CollectionConfig_StaticCollectionConfig{
						StaticCollectionConfig: &common.StaticCollectionConfig{
							Name: collectionName,
						},
					},
				},
			},
		},
	}, nil)
	dataStore.On("GetConfigHistoryRetriever").Return(historyRetreiver, nil)

	retriever := NewDataRetriever(dataStore)

	
	
	rwSets, err := retriever.CollectionRWSet(&gossip2.PvtDataDigest{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   uint64(5),
		TxId:       "testTxID",
		SeqInBlock: 1,
	})

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotNil(rwSets)
	assertion.NotEmpty(rwSets)
	assertion.Equal(2, len(rwSets.RWSet))

	var mergedRWSet []byte
	for _, rws := range rwSets.RWSet {
		mergedRWSet = append(mergedRWSet, rws...)
	}

	assertion.Equal([]byte{1, 2, 3, 4}, mergedRWSet)
}

func TestNewDataRetriever_FailGetPvtDataFromLedger(t *testing.T) {
	t.Parallel()
	dataStore := &mockedDataStore{}

	namespace := "testChaincodeName1"
	collectionName := "testCollectionName"

	dataStore.On("LedgerHeight").Return(uint64(10), nil)
	dataStore.On("GetPvtDataByNum", uint64(5), mock.Anything).
		Return(nil, errors.New("failing retrieving private data"))

	retriever := NewDataRetriever(dataStore)

	
	
	rwSets, err := retriever.CollectionRWSet(&gossip2.PvtDataDigest{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   uint64(5),
		TxId:       "testTxID",
		SeqInBlock: 1,
	})

	assertion := assert.New(t)
	assertion.Error(err)
	assertion.Nil(rwSets)
}

func TestNewDataRetriever_GetOnlyRelevantPvtData(t *testing.T) {
	t.Parallel()
	dataStore := &mockedDataStore{}

	namespace := "testChaincodeName1"
	collectionName := "testCollectionName"

	result := []*ledger.TxPvtData{{
		WriteSet: &rwset.TxPvtReadWriteSet{
			DataModel: rwset.TxReadWriteSet_KV,
			NsPvtRwset: []*rwset.NsPvtReadWriteSet{
				{
					Namespace: namespace,
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{{
						CollectionName: collectionName,
						Rwset:          []byte{1},
					}},
				},
				{
					Namespace: namespace,
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{{
						CollectionName: collectionName,
						Rwset:          []byte{2},
					}},
				},
				{
					Namespace: "invalidNamespace",
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{{
						CollectionName: collectionName,
						Rwset:          []byte{0, 0},
					}},
				},
				{
					Namespace: namespace,
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{{
						CollectionName: "invalidCollectionName",
						Rwset:          []byte{0, 0},
					}},
				},
			},
		},
		SeqInBlock: 1,
	}}

	dataStore.On("LedgerHeight").Return(uint64(10), nil)
	dataStore.On("GetPvtDataByNum", uint64(5), mock.Anything).Return(result, nil)
	historyRetreiver := &mockedHistoryRetreiver{}
	historyRetreiver.On("MostRecentCollectionConfigBelow", mock.Anything, namespace).Return(&ledger.CollectionConfigInfo{
		CollectionConfig: &common.CollectionConfigPackage{
			Config: []*common.CollectionConfig{
				{
					Payload: &common.CollectionConfig_StaticCollectionConfig{
						StaticCollectionConfig: &common.StaticCollectionConfig{
							Name: collectionName,
						},
					},
				},
			},
		},
	}, nil)
	dataStore.On("GetConfigHistoryRetriever").Return(historyRetreiver, nil)

	retriever := NewDataRetriever(dataStore)

	
	
	rwSets, err := retriever.CollectionRWSet(&gossip2.PvtDataDigest{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   uint64(5),
		TxId:       "testTxID",
		SeqInBlock: 1,
	})

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotNil(rwSets)
	assertion.NotEmpty(rwSets)
	assertion.Equal(2, len(rwSets.RWSet))

	var mergedRWSet []byte
	for _, rws := range rwSets.RWSet {
		mergedRWSet = append(mergedRWSet, rws...)
	}

	assertion.Equal([]byte{1, 2}, mergedRWSet)

}
