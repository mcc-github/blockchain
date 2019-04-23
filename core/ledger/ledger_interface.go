/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-lib-go/healthz"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/core/ledger/util/couchdb"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/mcc-github/blockchain/protos/peer"
)


type Initializer struct {
	StateListeners                []StateListener
	DeployedChaincodeInfoProvider DeployedChaincodeInfoProvider
	MembershipInfoProvider        MembershipInfoProvider
	MetricsProvider               metrics.Provider
	HealthCheckRegistry           HealthCheckRegistry
	Config                        *Config
}


type Config struct {
	
	StateDB *StateDB
}


type StateDB struct {
	
	
	StateDatabase string
	
	
	CouchDB *couchdb.Config
}


type PeerLedgerProvider interface {
	Initialize(initializer *Initializer) error
	
	
	
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
	
	GetConfigHistoryRetriever() (ConfigHistoryRetriever, error)
	
	
	
	
	CommitPvtDataOfOldBlocks(blockPvtData []*BlockPvtData) ([]*PvtdataHashMismatch, error)
	
	GetMissingPvtDataTracker() (MissingPvtDataTracker, error)
}



type ValidatedLedger interface {
	commonledger.Ledger
}


type SimpleQueryExecutor interface {
	
	GetState(namespace string, key string) ([]byte, error)
	
	
	
	
	
	GetStateRangeScanIterator(namespace string, startKey string, endKey string) (commonledger.ResultsIterator, error)
	
	
	
	
	GetPrivateDataHash(namespace, collection, key string) ([]byte, error)
}






type QueryExecutor interface {
	SimpleQueryExecutor
	
	GetStateMetadata(namespace, key string) (map[string][]byte, error)
	
	GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error)
	
	
	
	
	
	
	GetStateRangeScanIteratorWithMetadata(namespace string, startKey, endKey string, metadata map[string]interface{}) (QueryResultsIterator, error)
	
	
	
	
	ExecuteQuery(namespace, query string) (commonledger.ResultsIterator, error)
	
	
	
	
	
	ExecuteQueryWithMetadata(namespace, query string, metadata map[string]interface{}) (QueryResultsIterator, error)
	
	GetPrivateData(namespace, collection, key string) ([]byte, error)
	
	GetPrivateDataMetadata(namespace, collection, key string) (map[string][]byte, error)
	
	GetPrivateDataMetadataByHash(namespace, collection string, keyhash []byte) (map[string][]byte, error)
	
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


type QueryResultsIterator interface {
	commonledger.ResultsIterator
	
	GetBookmarkAndClose() string
}


type TxPvtData struct {
	SeqInBlock uint64
	WriteSet   *rwset.TxPvtReadWriteSet
}


type TxPvtDataMap map[uint64]*TxPvtData





type MissingPvtData struct {
	Namespace  string
	Collection string
	IsEligible bool
}



type TxMissingPvtDataMap map[uint64][]*MissingPvtData



type BlockAndPvtData struct {
	Block          *common.Block
	PvtData        TxPvtDataMap
	MissingPvtData TxMissingPvtDataMap
}


type BlockPvtData struct {
	BlockNum  uint64
	WriteSets TxPvtDataMap
}


func (txMissingPvtData TxMissingPvtDataMap) Add(txNum uint64, ns, coll string, isEligible bool) {
	txMissingPvtData[txNum] = append(txMissingPvtData[txNum], &MissingPvtData{ns, coll, isEligible})
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
	Initialize(ledgerID string, qe SimpleQueryExecutor) error
	InterestedInNamespaces() []string
	HandleStateUpdates(trigger *StateUpdateTrigger) error
	StateCommitDone(channelID string)
}


type StateUpdateTrigger struct {
	LedgerID                    string
	StateUpdates                StateUpdates
	CommittingBlockNum          uint64
	CommittedStateQueryExecutor SimpleQueryExecutor
	PostCommitQueryExecutor     SimpleQueryExecutor
}


type StateUpdates map[string]*KVStateUpdates


type KVStateUpdates struct {
	PublicUpdates   []*kvrwset.KVWrite
	CollHashUpdates map[string][]*kvrwset.KVWriteHash
}


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
	Namespace, Collection string
}


type CollectionConfigInfo struct {
	CollectionConfig   *common.CollectionConfigPackage
	CommittingBlockNum uint64
}


func (missingPvtDataInfo MissingPvtDataInfo) Add(blkNum, txNum uint64, ns, coll string) {
	missingBlockPvtDataInfo, ok := missingPvtDataInfo[blkNum]
	if !ok {
		missingBlockPvtDataInfo = make(MissingBlockPvtdataInfo)
		missingPvtDataInfo[blkNum] = missingBlockPvtDataInfo
	}

	if _, ok := missingBlockPvtDataInfo[txNum]; !ok {
		missingBlockPvtDataInfo[txNum] = []*MissingCollectionPvtDataInfo{}
	}

	missingBlockPvtDataInfo[txNum] = append(missingBlockPvtDataInfo[txNum],
		&MissingCollectionPvtDataInfo{
			Namespace:  ns,
			Collection: coll})
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



type CollConfigNotDefinedError struct {
	Ns string
}

func (e *CollConfigNotDefinedError) Error() string {
	return fmt.Sprintf("collection config not defined for chaincode [%s], pass the collection configuration upon chaincode definition/instantiation", e.Ns)
}



type InvalidCollNameError struct {
	Ns, Coll string
}

func (e *InvalidCollNameError) Error() string {
	return fmt.Sprintf("collection [%s] not defined in the collection config for chaincode [%s]", e.Coll, e.Ns)
}




type PvtdataHashMismatch struct {
	BlockNum, TxNum       uint64
	Namespace, Collection string
	ExpectedHash          []byte
}



type DeployedChaincodeInfoProvider interface {
	
	Namespaces() []string
	
	UpdatedChaincodes(stateUpdates map[string][]*kvrwset.KVWrite) ([]*ChaincodeLifecycleInfo, error)
	
	ChaincodeInfo(channelName, chaincodeName string, qe SimpleQueryExecutor) (*DeployedChaincodeInfo, error)
	
	CollectionInfo(channelName, chaincodeName, collectionName string, qe SimpleQueryExecutor) (*common.StaticCollectionConfig, error)
	
	ImplicitCollections(channelName, chaincodeName string, qe SimpleQueryExecutor) ([]*common.StaticCollectionConfig, error)
}


type DeployedChaincodeInfo struct {
	Name                        string
	Hash                        []byte
	Version                     string
	ExplicitCollectionConfigPkg *common.CollectionConfigPackage
	ImplicitCollections         []*common.StaticCollectionConfig
}


func (dci DeployedChaincodeInfo) AllCollectionsConfigPkg() *common.CollectionConfigPackage {
	var combinedColls []*common.CollectionConfig
	if dci.ExplicitCollectionConfigPkg != nil {
		for _, explicitColl := range dci.ExplicitCollectionConfigPkg.Config {
			combinedColls = append(combinedColls, explicitColl)
		}
	}
	for _, implicitColl := range dci.ImplicitCollections {
		c := &common.CollectionConfig{}
		c.Payload = &common.CollectionConfig_StaticCollectionConfig{StaticCollectionConfig: implicitColl}
		combinedColls = append(combinedColls, c)
	}
	if combinedColls == nil {
		return nil
	}
	return &common.CollectionConfigPackage{
		Config: combinedColls,
	}
}


type ChaincodeLifecycleInfo struct {
	Name    string
	Deleted bool
	Details *ChaincodeLifecycleDetails 
}


type ChaincodeLifecycleDetails struct {
	Updated bool 
	
	HashChanged        bool     
	CollectionsUpdated []string 
	CollectionsRemoved []string 
}



type MembershipInfoProvider interface {
	
	AmMemberOf(channelName string, collectionPolicyConfig *common.CollectionPolicyConfig) (bool, error)
}

type HealthCheckRegistry interface {
	RegisterChecker(string, healthz.HealthChecker) error
}







