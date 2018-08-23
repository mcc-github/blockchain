/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import (
	"github.com/golang/protobuf/proto"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/peer"
)


type PeerLedgerProvider interface {
	Initialize(statelisteners []StateListener)
	
	
	
	Create(genesisBlock *common.Block) (PeerLedger, error)
	
	Open(ledgerID string) (PeerLedger, error)
	
	Exists(ledgerID string) (bool, error)
	
	List() ([]string, error)
	
	Close()
}



type PeerLedger interface {
	commonledger.Ledger
	
	GetTransactionByID(txID string) (*peer.ProcessedTransaction, error)
	
	GetBlockByHash(blockHash []byte) (*common.Block, error)
	
	GetBlockByTxID(txID string) (*common.Block, error)
	
	GetTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error)
	
	
	
	NewTxSimulator(txid string) (TxSimulator, error)
	
	
	
	NewQueryExecutor() (QueryExecutor, error)
	
	
	
	NewHistoryQueryExecutor() (HistoryQueryExecutor, error)
	
	
	
	GetPvtDataAndBlockByNum(blockNum uint64, filter PvtNsCollFilter) (*BlockAndPvtData, error)
	
	
	
	GetPvtDataByNum(blockNum uint64, filter PvtNsCollFilter) ([]*TxPvtData, error)
	
	CommitWithPvtData(blockAndPvtdata *BlockAndPvtData) error
	
	
	
	PurgePrivateData(maxBlockNumToRetain uint64) error
	
	PrivateDataMinBlockNum() (uint64, error)
	
	Prune(policy commonledger.PrunePolicy) error
	
	GetConfigHistoryRetriever() (ConfigHistoryRetriever, error)
	
	
	
	
	CommitPvtData(blockPvtData []*BlockPvtData) ([]*PvtdataHashMismatch, error)
	
	GetMissingPvtDataTracker() (MissingPvtDataTracker, error)
}



type ValidatedLedger interface {
	commonledger.Ledger
}






type QueryExecutor interface {
	
	GetState(namespace string, key string) ([]byte, error)
	
	GetStateMetadata(namespace, key string) (map[string][]byte, error)
	
	GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error)
	
	
	
	
	
	GetStateRangeScanIterator(namespace string, startKey string, endKey string) (commonledger.ResultsIterator, error)
	
	
	
	
	ExecuteQuery(namespace, query string) (commonledger.ResultsIterator, error)
	
	GetPrivateData(namespace, collection, key string) ([]byte, error)
	
	GetPrivateDataMetadata(namespace, collection, key string) (map[string][]byte, error)
	
	GetPrivateDataMultipleKeys(namespace, collection string, keys []string) ([][]byte, error)
	
	
	
	
	
	GetPrivateDataRangeScanIterator(namespace, collection, startKey, endKey string) (commonledger.ResultsIterator, error)
	
	
	
	
	ExecuteQueryOnPrivateData(namespace, collection, query string) (commonledger.ResultsIterator, error)
	
	Done()
}


type HistoryQueryExecutor interface {
	
	
	GetHistoryForKey(namespace string, key string) (commonledger.ResultsIterator, error)
}



type TxSimulator interface {
	QueryExecutor
	
	SetState(namespace string, key string, value []byte) error
	
	DeleteState(namespace string, key string) error
	
	SetStateMultipleKeys(namespace string, kvs map[string][]byte) error
	
	SetStateMetadata(namespace, key string, metadata map[string][]byte) error
	
	DeleteStateMetadata(namespace, key string) error
	
	ExecuteUpdate(query string) error
	
	SetPrivateData(namespace, collection, key string, value []byte) error
	
	SetPrivateDataMultipleKeys(namespace, collection string, kvs map[string][]byte) error
	
	DeletePrivateData(namespace, collection, key string) error
	
	SetPrivateDataMetadata(namespace, collection, key string, metadata map[string][]byte) error
	
	DeletePrivateDataMetadata(namespace, collection, key string) error
	
	
	
	
	
	
	
	
	
	
	GetTxSimulationResults() (*TxSimulationResults, error)
}


type TxPvtData struct {
	SeqInBlock uint64
	WriteSet   *rwset.TxPvtReadWriteSet
}




type MissingPrivateData struct {
	TxId       string
	SeqInBlock int
	Namespace  string
	Collection string
	Eligible   bool
}



type BlockAndPvtData struct {
	Block        *common.Block
	BlockPvtData map[uint64]*TxPvtData
	Missing      []*MissingPrivateData
}


type BlockPvtData struct {
	BlockNum  uint64
	WriteSets map[uint64]*TxPvtData
}


type PvtCollFilter map[string]bool


type PvtNsCollFilter map[string]PvtCollFilter


func NewPvtNsCollFilter() PvtNsCollFilter {
	return make(map[string]PvtCollFilter)
}


func (pvtdata *TxPvtData) Has(ns string, coll string) bool {
	if pvtdata.WriteSet == nil {
		return false
	}
	for _, nsdata := range pvtdata.WriteSet.NsPvtRwset {
		if nsdata.Namespace == ns {
			for _, colldata := range nsdata.CollectionPvtRwset {
				if colldata.CollectionName == coll {
					return true
				}
			}
		}
	}
	return false
}


func (filter PvtNsCollFilter) Add(ns string, coll string) {
	collFilter, ok := filter[ns]
	if !ok {
		collFilter = make(map[string]bool)
		filter[ns] = collFilter
	}
	collFilter[coll] = true
}


func (filter PvtNsCollFilter) Has(ns string, coll string) bool {
	collFilter, ok := filter[ns]
	if !ok {
		return false
	}
	return collFilter[coll]
}


type TxSimulationResults struct {
	PubSimulationResults *rwset.TxReadWriteSet
	PvtSimulationResults *rwset.TxPvtReadWriteSet
}


func (txSim *TxSimulationResults) GetPubSimulationBytes() ([]byte, error) {
	return proto.Marshal(txSim.PubSimulationResults)
}


func (txSim *TxSimulationResults) GetPvtSimulationBytes() ([]byte, error) {
	if !txSim.ContainsPvtWrites() {
		return nil, nil
	}
	return proto.Marshal(txSim.PvtSimulationResults)
}


func (txSim *TxSimulationResults) ContainsPvtWrites() bool {
	return txSim.PvtSimulationResults != nil
}














type StateListener interface {
	InterestedInNamespaces() []string
	HandleStateUpdates(ledgerID string, stateUpdates StateUpdates, committingBlockNum uint64) error
	StateCommitDone(channelID string)
}


type StateUpdates map[string]interface{}


type ConfigHistoryRetriever interface {
	CollectionConfigAt(blockNum uint64, chaincodeName string) (*CollectionConfigInfo, error)
	MostRecentCollectionConfigBelow(blockNum uint64, chaincodeName string) (*CollectionConfigInfo, error)
}


type MissingPvtDataTracker interface {
	GetMissingPvtDataInfoForMostRecentBlocks(maxBlocks int) (MissingPvtDataInfo, error)
}


type MissingPvtDataInfo map[uint64]MissingBlockPvtdataInfo


type MissingBlockPvtdataInfo map[uint64][]*MissingCollectionPvtDataInfo


type MissingCollectionPvtDataInfo struct {
	ChaincodeName, CollectionName string
}


type CollectionConfigInfo struct {
	CollectionConfig   *common.CollectionConfigPackage
	CommittingBlockNum uint64
}




type ErrCollectionConfigNotYetAvailable struct {
	MaxBlockNumCommitted uint64
	Msg                  string
}

func (e *ErrCollectionConfigNotYetAvailable) Error() string {
	return e.Msg
}


type NotFoundInIndexErr string

func (NotFoundInIndexErr) Error() string {
	return "Entry not found in index"
}




type PvtdataHashMismatch struct {
	BlockNum, TxNum               uint64
	ChaincodeName, CollectionName string
	ExpectedHash                  []byte
}
