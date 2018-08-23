/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package confighistory

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("confighistory")

const (
	lsccNamespace = "lscc"
)


type Mgr interface {
	ledger.StateListener
	GetRetriever(ledgerID string, ledgerInfoRetriever LedgerInfoRetriever) ledger.ConfigHistoryRetriever
	Close()
}

type mgr struct {
	dbProvider *dbProvider
}


func NewMgr() Mgr {
	return newMgr(dbPath())
}

func newMgr(dbPath string) Mgr {
	return &mgr{newDBProvider(dbPath)}
}


func (m *mgr) InterestedInNamespaces() []string {
	return []string{lsccNamespace}
}


func (m *mgr) StateCommitDone(ledgerID string) {
	
}





func (m *mgr) HandleStateUpdates(trigger *ledger.StateUpdateTrigger) error {
	batch := prepareDBBatch(trigger.StateUpdates, trigger.CommittingBlockNum)
	dbHandle := m.dbProvider.getDB(trigger.LedgerID)
	return dbHandle.writeBatch(batch, true)
}


func (m *mgr) GetRetriever(ledgerID string, ledgerInfoRetriever LedgerInfoRetriever) ledger.ConfigHistoryRetriever {
	return &retriever{dbHandle: m.dbProvider.getDB(ledgerID), ledgerInfoRetriever: ledgerInfoRetriever}
}


func (m *mgr) Close() {
	m.dbProvider.Close()
}

type retriever struct {
	ledgerInfoRetriever LedgerInfoRetriever
	dbHandle            *db
}


func (r *retriever) MostRecentCollectionConfigBelow(blockNum uint64, chaincodeName string) (*ledger.CollectionConfigInfo, error) {
	compositeKV, err := r.dbHandle.mostRecentEntryBelow(blockNum, lsccNamespace, constructCollectionConfigKey(chaincodeName))
	if err != nil || compositeKV == nil {
		return nil, err
	}
	return compositeKVToCollectionConfig(compositeKV)
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

	compositeKV, err := r.dbHandle.entryAt(blockNum, lsccNamespace, constructCollectionConfigKey(chaincodeName))
	if err != nil || compositeKV == nil {
		return nil, err
	}
	return compositeKVToCollectionConfig(compositeKV)
}

func prepareDBBatch(stateUpdates ledger.StateUpdates, committingBlock uint64) *batch {
	batch := newBatch()
	lsccWrites := stateUpdates[lsccNamespace]
	for _, kv := range lsccWrites.([]*kvrwset.KVWrite) {
		if !privdata.IsCollectionConfigKey(kv.Key) {
			continue
		}
		batch.add(lsccNamespace, kv.Key, committingBlock, kv.Value)
	}
	return batch
}

func compositeKVToCollectionConfig(compositeKV *compositeKV) (*ledger.CollectionConfigInfo, error) {
	conf := &common.CollectionConfigPackage{}
	if err := proto.Unmarshal(compositeKV.value, conf); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling compositeKV to collection config")
	}
	return &ledger.CollectionConfigInfo{CollectionConfig: conf, CommittingBlockNum: compositeKV.blockNum}, nil
}

func constructCollectionConfigKey(chaincodeName string) string {
	return privdata.BuildCollectionKVSKey(chaincodeName)
}

func isCollectionConfigKey(key string) bool {
	return privdata.IsCollectionConfigKey(key)
}

func dbPath() string {
	return ledgerconfig.GetConfigHistoryPath()
}


type LedgerInfoRetriever interface {
	GetBlockchainInfo() (*common.BlockchainInfo, error)
}
