/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package confighistory

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("confighistory")

const (
	collectionConfigNamespace = "lscc" 
)


type Mgr interface {
	ledger.StateListener
	GetRetriever(ledgerID string, ledgerInfoRetriever LedgerInfoRetriever) ledger.ConfigHistoryRetriever
	Close()
}

type mgr struct {
	ccInfoProvider ledger.DeployedChaincodeInfoProvider
	dbProvider     *dbProvider
}


func NewMgr(ccInfoProvider ledger.DeployedChaincodeInfoProvider) Mgr {
	return newMgr(ccInfoProvider, dbPath())
}

func newMgr(ccInfoProvider ledger.DeployedChaincodeInfoProvider, dbPath string) Mgr {
	return &mgr{ccInfoProvider, newDBProvider(dbPath)}
}

func (m *mgr) Initialize(ledgerID string, qe ledger.SimpleQueryExecutor) error {
	
	return nil
}


func (m *mgr) InterestedInNamespaces() []string {
	return m.ccInfoProvider.Namespaces()
}


func (m *mgr) StateCommitDone(ledgerID string) {
	
}





func (m *mgr) HandleStateUpdates(trigger *ledger.StateUpdateTrigger) error {
	updatedCCs, err := m.ccInfoProvider.UpdatedChaincodes(extractPublicUpdates(trigger.StateUpdates))
	if err != nil {
		return err
	}
	if len(updatedCCs) == 0 {
		logger.Errorf("Config history manager is expected to recieve events only if at least one chaincode is updated stateUpdates = %#v",
			trigger.StateUpdates)
		return nil
	}
	updatedCollConfigs := map[string]*common.CollectionConfigPackage{}
	for _, cc := range updatedCCs {
		ccInfo, err := m.ccInfoProvider.ChaincodeInfo(trigger.LedgerID, cc.Name, trigger.PostCommitQueryExecutor)
		if err != nil {
			return err
		}
		if ccInfo.ExplicitCollectionConfigPkg == nil {
			continue
		}
		updatedCollConfigs[ccInfo.Name] = ccInfo.ExplicitCollectionConfigPkg
	}
	if len(updatedCollConfigs) == 0 {
		return nil
	}
	batch, err := prepareDBBatch(updatedCollConfigs, trigger.CommittingBlockNum)
	if err != nil {
		return err
	}
	dbHandle := m.dbProvider.getDB(trigger.LedgerID)
	return dbHandle.writeBatch(batch, true)
}


func (m *mgr) GetRetriever(ledgerID string, ledgerInfoRetriever LedgerInfoRetriever) ledger.ConfigHistoryRetriever {
	return &retriever{
		ledgerInfoRetriever:    ledgerInfoRetriever,
		ledgerID:               ledgerID,
		deployedCCInfoProvider: m.ccInfoProvider,
		dbHandle:               m.dbProvider.getDB(ledgerID),
	}
}


func (m *mgr) Close() {
	m.dbProvider.Close()
}

type retriever struct {
	ledgerInfoRetriever    LedgerInfoRetriever
	ledgerID               string
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider
	dbHandle               *db
}


func (r *retriever) MostRecentCollectionConfigBelow(blockNum uint64, chaincodeName string) (*ledger.CollectionConfigInfo, error) {
	compositeKV, err := r.dbHandle.mostRecentEntryBelow(blockNum, collectionConfigNamespace, constructCollectionConfigKey(chaincodeName))
	if err != nil {
		return nil, err
	}
	implicitColls, err := r.getImplicitCollection(chaincodeName)
	if err != nil {
		return nil, err
	}

	return constructCollectionConfigInfo(compositeKV, implicitColls)
}


func (r *retriever) CollectionConfigAt(blockNum uint64, chaincodeName string) (*ledger.CollectionConfigInfo, error) {
	info, err := r.ledgerInfoRetriever.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}
	maxCommittedBlockNum := info.Height - 1
	if maxCommittedBlockNum < blockNum {
		return nil, &ledger.ErrCollectionConfigNotYetAvailable{MaxBlockNumCommitted: maxCommittedBlockNum,
			Msg: fmt.Sprintf("The maximum block number committed [%d] is less than the requested block number [%d]", maxCommittedBlockNum, blockNum)}
	}

	compositeKV, err := r.dbHandle.entryAt(blockNum, collectionConfigNamespace, constructCollectionConfigKey(chaincodeName))
	if err != nil {
		return nil, err
	}
	implicitColls, err := r.getImplicitCollection(chaincodeName)
	if err != nil {
		return nil, err
	}
	return constructCollectionConfigInfo(compositeKV, implicitColls)
}

func (r *retriever) getImplicitCollection(chaincodeName string) ([]*common.StaticCollectionConfig, error) {
	qe, err := r.ledgerInfoRetriever.NewQueryExecutor()
	if err != nil {
		return nil, err
	}
	defer qe.Done()
	return r.deployedCCInfoProvider.ImplicitCollections(r.ledgerID, chaincodeName, qe)
}

func prepareDBBatch(chaincodeCollConfigs map[string]*common.CollectionConfigPackage, committingBlockNum uint64) (*batch, error) {
	batch := newBatch()
	for ccName, collConfig := range chaincodeCollConfigs {
		key := constructCollectionConfigKey(ccName)
		var configBytes []byte
		var err error
		if configBytes, err = proto.Marshal(collConfig); err != nil {
			return nil, errors.WithStack(err)
		}
		batch.add(collectionConfigNamespace, key, committingBlockNum, configBytes)
	}
	return batch, nil
}

func compositeKVToCollectionConfig(compositeKV *compositeKV) (*ledger.CollectionConfigInfo, error) {
	conf := &common.CollectionConfigPackage{}
	if err := proto.Unmarshal(compositeKV.value, conf); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling compositeKV to collection config")
	}
	return &ledger.CollectionConfigInfo{
		CollectionConfig:   conf,
		CommittingBlockNum: compositeKV.blockNum,
	}, nil
}

func constructCollectionConfigKey(chaincodeName string) string {
	return chaincodeName + "~collection" 
}

func dbPath() string {
	return ledgerconfig.GetConfigHistoryPath()
}

func extractPublicUpdates(stateUpdates ledger.StateUpdates) map[string][]*kvrwset.KVWrite {
	m := map[string][]*kvrwset.KVWrite{}
	for ns, updates := range stateUpdates {
		m[ns] = updates.PublicUpdates
	}
	return m
}

func constructCollectionConfigInfo(
	compositeKV *compositeKV,
	implicitColls []*common.StaticCollectionConfig,
) (*ledger.CollectionConfigInfo, error) {
	var collConf *ledger.CollectionConfigInfo
	var err error

	if compositeKV == nil && len(implicitColls) == 0 {
		return nil, nil
	}

	collConf = &ledger.CollectionConfigInfo{
		CollectionConfig: &common.CollectionConfigPackage{},
	}
	if compositeKV != nil {
		if collConf, err = compositeKVToCollectionConfig(compositeKV); err != nil {
			return nil, err
		}
	}

	for _, implicitColl := range implicitColls {
		cc := &common.CollectionConfig{}
		cc.Payload = &common.CollectionConfig_StaticCollectionConfig{StaticCollectionConfig: implicitColl}
		collConf.CollectionConfig.Config = append(
			collConf.CollectionConfig.Config,
			cc,
		)
	}
	return collConf, nil
}


type LedgerInfoRetriever interface {
	GetBlockchainInfo() (*common.BlockchainInfo, error)
	NewQueryExecutor() (ledger.QueryExecutor, error)
}
