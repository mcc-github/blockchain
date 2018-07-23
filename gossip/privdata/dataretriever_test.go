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
	"github.com/mcc-github/blockchain/gossip/privdata/mocks"
	"github.com/mcc-github/blockchain/protos/common"
	gossip2 "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	transientstore2 "github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestNewDataRetriever_GetDataFromTransientStore(t *testing.T) {
	t.Parallel()
	dataStore := &mocks.DataStore{}

	rwSetScanner := &mocks.RWSetScanner{}
	namespace := "testChaincodeName1"
	collectionName := "testCollectionName"

	rwSetScanner.On("Close")
	rwSetScanner.On("NextWithConfig").Return(&transientstore.EndorserPvtSimulationResultsWithConfig{
		ReceivedAtBlockHeight:          2,
		PvtSimulationResultsWithConfig: nil,
	}, nil).Once().On("NextWithConfig").Return(&transientstore.EndorserPvtSimulationResultsWithConfig{
		ReceivedAtBlockHeight: 2,
		PvtSimulationResultsWithConfig: &transientstore2.TxPvtReadWriteSetWithConfigInfo{
			PvtRwset: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					pvtReadWriteSet(namespace, collectionName, []byte{1, 2}),
					pvtReadWriteSet(namespace, collectionName, []byte{3, 4}),
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

	
	
	rwSets, err := retriever.CollectionRWSet([]*gossip2.PvtDataDigest{{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   2,
		TxId:       "testTxID",
		SeqInBlock: 1,
	}}, 2)

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotEmpty(rwSets)
	dig2pvtRWSet := rwSets[DigKey{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   2,
		TxId:       "testTxID",
		SeqInBlock: 1,
	}]
	assertion.NotNil(dig2pvtRWSet)
	pvtRWSets := dig2pvtRWSet.RWSet
	assertion.Equal(2, len(pvtRWSets))

	var mergedRWSet []byte
	for _, rws := range pvtRWSets {
		mergedRWSet = append(mergedRWSet, rws...)
	}

	assertion.Equal([]byte{1, 2, 3, 4}, mergedRWSet)
}


func TestNewDataRetriever_GetDataFromLedger(t *testing.T) {
	t.Parallel()
	dataStore := &mocks.DataStore{}

	namespace := "testChaincodeName1"
	collectionName := "testCollectionName"

	result := []*ledger.TxPvtData{{
		WriteSet: &rwset.TxPvtReadWriteSet{
			DataModel: rwset.TxReadWriteSet_KV,
			NsPvtRwset: []*rwset.NsPvtReadWriteSet{
				pvtReadWriteSet(namespace, collectionName, []byte{1, 2}),
				pvtReadWriteSet(namespace, collectionName, []byte{3, 4}),
			},
		},
		SeqInBlock: 1,
	}}

	dataStore.On("LedgerHeight").Return(uint64(10), nil)
	dataStore.On("GetPvtDataByNum", uint64(5), mock.Anything).Return(result, nil)

	historyRetreiver := &mocks.ConfigHistoryRetriever{}
	historyRetreiver.On("MostRecentCollectionConfigBelow", mock.Anything, namespace).Return(newCollectionConfig(collectionName), nil)
	dataStore.On("GetConfigHistoryRetriever").Return(historyRetreiver, nil)

	retriever := NewDataRetriever(dataStore)

	
	
	rwSets, err := retriever.CollectionRWSet([]*gossip2.PvtDataDigest{{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   uint64(5),
		TxId:       "testTxID",
		SeqInBlock: 1,
	}}, uint64(5))

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotEmpty(rwSets)
	pvtRWSet := rwSets[DigKey{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   5,
		TxId:       "testTxID",
		SeqInBlock: 1,
	}]
	assertion.NotEmpty(pvtRWSet)
	assertion.Equal(2, len(pvtRWSet.RWSet))

	var mergedRWSet []byte
	for _, rws := range pvtRWSet.RWSet {
		mergedRWSet = append(mergedRWSet, rws...)
	}

	assertion.Equal([]byte{1, 2, 3, 4}, mergedRWSet)
}

func TestNewDataRetriever_FailGetPvtDataFromLedger(t *testing.T) {
	t.Parallel()
	dataStore := &mocks.DataStore{}

	namespace := "testChaincodeName1"
	collectionName := "testCollectionName"

	dataStore.On("LedgerHeight").Return(uint64(10), nil)
	dataStore.On("GetPvtDataByNum", uint64(5), mock.Anything).
		Return(nil, errors.New("failing retrieving private data"))

	retriever := NewDataRetriever(dataStore)

	
	
	rwSets, err := retriever.CollectionRWSet([]*gossip2.PvtDataDigest{{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   uint64(5),
		TxId:       "testTxID",
		SeqInBlock: 1,
	}}, uint64(5))

	assertion := assert.New(t)
	assertion.Error(err)
	assertion.Empty(rwSets)
}

func TestNewDataRetriever_GetOnlyRelevantPvtData(t *testing.T) {
	t.Parallel()
	dataStore := &mocks.DataStore{}

	namespace := "testChaincodeName1"
	collectionName := "testCollectionName"

	result := []*ledger.TxPvtData{{
		WriteSet: &rwset.TxPvtReadWriteSet{
			DataModel: rwset.TxReadWriteSet_KV,
			NsPvtRwset: []*rwset.NsPvtReadWriteSet{
				pvtReadWriteSet(namespace, collectionName, []byte{1}),
				pvtReadWriteSet(namespace, collectionName, []byte{2}),
				pvtReadWriteSet("invalidNamespace", collectionName, []byte{0, 0}),
				pvtReadWriteSet(namespace, "invalidCollectionName", []byte{0, 0}),
			},
		},
		SeqInBlock: 1,
	}}

	dataStore.On("LedgerHeight").Return(uint64(10), nil)
	dataStore.On("GetPvtDataByNum", uint64(5), mock.Anything).Return(result, nil)
	historyRetreiver := &mocks.ConfigHistoryRetriever{}
	historyRetreiver.On("MostRecentCollectionConfigBelow", mock.Anything, namespace).Return(newCollectionConfig(collectionName), nil)
	dataStore.On("GetConfigHistoryRetriever").Return(historyRetreiver, nil)

	retriever := NewDataRetriever(dataStore)

	
	
	rwSets, err := retriever.CollectionRWSet([]*gossip2.PvtDataDigest{{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   uint64(5),
		TxId:       "testTxID",
		SeqInBlock: 1,
	}}, 5)

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotEmpty(rwSets)
	pvtRWSet := rwSets[DigKey{
		Namespace:  namespace,
		Collection: collectionName,
		BlockSeq:   5,
		TxId:       "testTxID",
		SeqInBlock: 1,
	}]
	assertion.NotEmpty(pvtRWSet)
	assertion.Equal(2, len(pvtRWSet.RWSet))

	var mergedRWSet []byte
	for _, rws := range pvtRWSet.RWSet {
		mergedRWSet = append(mergedRWSet, rws...)
	}

	assertion.Equal([]byte{1, 2}, mergedRWSet)
}

func TestNewDataRetriever_GetMultipleDigests(t *testing.T) {
	t.Parallel()
	dataStore := &mocks.DataStore{}

	ns1, ns2 := "testChaincodeName1", "testChaincodeName2"
	col1, col2 := "testCollectionName1", "testCollectionName2"

	result := []*ledger.TxPvtData{
		{
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					pvtReadWriteSet(ns1, col1, []byte{1}),
					pvtReadWriteSet(ns1, col1, []byte{2}),
					pvtReadWriteSet("invalidNamespace", col1, []byte{0, 0}),
					pvtReadWriteSet(ns1, "invalidCollectionName", []byte{0, 0}),
				},
			},
			SeqInBlock: 1,
		},
		{
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					pvtReadWriteSet(ns2, col2, []byte{3}),
					pvtReadWriteSet(ns2, col2, []byte{4}),
					pvtReadWriteSet("invalidNamespace", col2, []byte{0, 0}),
					pvtReadWriteSet(ns2, "invalidCollectionName", []byte{0, 0}),
				},
			},
			SeqInBlock: 2,
		},
		{
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					pvtReadWriteSet(ns1, col1, []byte{5}),
					pvtReadWriteSet(ns2, col2, []byte{6}),
					pvtReadWriteSet("invalidNamespace", col2, []byte{0, 0}),
					pvtReadWriteSet(ns2, "invalidCollectionName", []byte{0, 0}),
				},
			},
			SeqInBlock: 3,
		},
	}

	dataStore.On("LedgerHeight").Return(uint64(10), nil)
	dataStore.On("GetPvtDataByNum", uint64(5), mock.Anything).Return(result, nil)
	historyRetreiver := &mocks.ConfigHistoryRetriever{}
	historyRetreiver.On("MostRecentCollectionConfigBelow", mock.Anything, ns1).Return(newCollectionConfig(col1), nil)
	historyRetreiver.On("MostRecentCollectionConfigBelow", mock.Anything, ns2).Return(newCollectionConfig(col2), nil)
	dataStore.On("GetConfigHistoryRetriever").Return(historyRetreiver, nil)

	retriever := NewDataRetriever(dataStore)

	
	
	rwSets, err := retriever.CollectionRWSet([]*gossip2.PvtDataDigest{{
		Namespace:  ns1,
		Collection: col1,
		BlockSeq:   uint64(5),
		TxId:       "testTxID",
		SeqInBlock: 1,
	}, {
		Namespace:  ns2,
		Collection: col2,
		BlockSeq:   uint64(5),
		TxId:       "testTxID",
		SeqInBlock: 2,
	}}, 5)

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotEmpty(rwSets)
	assertion.Equal(2, len(rwSets))

	pvtRWSet := rwSets[DigKey{
		Namespace:  ns1,
		Collection: col1,
		BlockSeq:   5,
		TxId:       "testTxID",
		SeqInBlock: 1,
	}]
	assertion.NotEmpty(pvtRWSet)
	assertion.Equal(2, len(pvtRWSet.RWSet))

	var mergedRWSet []byte
	for _, rws := range pvtRWSet.RWSet {
		mergedRWSet = append(mergedRWSet, rws...)
	}

	pvtRWSet = rwSets[DigKey{
		Namespace:  ns2,
		Collection: col2,
		BlockSeq:   5,
		TxId:       "testTxID",
		SeqInBlock: 2,
	}]
	assertion.NotEmpty(pvtRWSet)
	assertion.Equal(2, len(pvtRWSet.RWSet))
	for _, rws := range pvtRWSet.RWSet {
		mergedRWSet = append(mergedRWSet, rws...)
	}

	assertion.Equal([]byte{1, 2, 3, 4}, mergedRWSet)
}

func newCollectionConfig(collectionName string) *ledger.CollectionConfigInfo {
	return &ledger.CollectionConfigInfo{
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
	}
}

func pvtReadWriteSet(ns string, collectionName string, data []byte) *rwset.NsPvtReadWriteSet {
	return &rwset.NsPvtReadWriteSet{
		Namespace: ns,
		CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{{
			CollectionName: collectionName,
			Rwset:          data,
		}},
	}
}