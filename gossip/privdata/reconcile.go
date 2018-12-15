/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"encoding/hex"
	"fmt"
	"math"
	"sync"
	"time"

	util2 "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/committer"
	"github.com/mcc-github/blockchain/core/ledger"
	privdatacommon "github.com/mcc-github/blockchain/gossip/privdata/common"
	"github.com/mcc-github/blockchain/protos/common"
	gossip2 "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	reconcileSleepIntervalConfigKey = "peer.gossip.pvtData.reconcileSleepInterval"
	reconcileSleepIntervalDefault   = time.Minute * 1
	reconcileBatchSizeConfigKey     = "peer.gossip.pvtData.reconcileBatchSize"
	reconcileBatchSizeDefault       = 10
	reconciliationEnabledConfigKey  = "peer.gossip.pvtData.reconciliationEnabled"
)



type ReconciliationFetcher interface {
	FetchReconciledItems(dig2collectionConfig privdatacommon.Dig2CollectionConfig) (*privdatacommon.FetchedPvtDataContainer, error)
}







type PvtDataReconciler interface {
	
	Start()
	
	Stop()
}

type Reconciler struct {
	config *ReconcilerConfig
	ReconciliationFetcher
	committer.Committer
	stopChan  chan struct{}
	startOnce sync.Once
	stopOnce  sync.Once
}



type NoOpReconciler struct {
}

func (*NoOpReconciler) Start() {
	
	logger.Debug("Private data reconciliation has been disabled")
}

func (*NoOpReconciler) Stop() {
	
}


type ReconcilerConfig struct {
	sleepInterval time.Duration
	batchSize     int
	IsEnabled     bool
}


func GetReconcilerConfig() *ReconcilerConfig {
	reconcileSleepInterval := viper.GetDuration(reconcileSleepIntervalConfigKey)
	if reconcileSleepInterval == 0 {
		logger.Warning("Configuration key", reconcileSleepIntervalConfigKey, "isn't set, defaulting to", reconcileSleepIntervalDefault)
		reconcileSleepInterval = reconcileSleepIntervalDefault
	}
	reconcileBatchSize := viper.GetInt(reconcileBatchSizeConfigKey)
	if reconcileBatchSize == 0 {
		logger.Warning("Configuration key", reconcileBatchSizeConfigKey, "isn't set, defaulting to", reconcileBatchSizeDefault)
		reconcileBatchSize = reconcileBatchSizeDefault
	}
	isEnabled := viper.GetBool(reconciliationEnabledConfigKey)
	return &ReconcilerConfig{sleepInterval: reconcileSleepInterval, batchSize: reconcileBatchSize, IsEnabled: isEnabled}
}


func NewReconciler(c committer.Committer, fetcher ReconciliationFetcher, config *ReconcilerConfig) *Reconciler {
	logger.Debug("Private data reconciliation is enabled")
	return &Reconciler{
		config:                config,
		Committer:             c,
		ReconciliationFetcher: fetcher,
		stopChan:              make(chan struct{}),
	}
}

func (r *Reconciler) Stop() {
	r.stopOnce.Do(func() {
		close(r.stopChan)
	})
}

func (r *Reconciler) Start() {
	r.startOnce.Do(func() {
		go r.run()
	})
}

func (r *Reconciler) run() {
	for {
		select {
		case <-r.stopChan:
			return
		case <-time.After(r.config.sleepInterval):
			logger.Debug("Start reconcile missing private info")
			if err := r.reconcile(); err != nil {
				logger.Error("Failed to reconcile missing private info, error: ", err.Error())
				break
			}
		}
	}
}


func (r *Reconciler) reconcile() error {
	missingPvtDataTracker, err := r.GetMissingPvtDataTracker()
	if err != nil {
		logger.Error("reconciliation error when trying to get missingPvtDataTrcker:", err)
		return err
	}
	if missingPvtDataTracker == nil {
		logger.Error("got nil as MissingPvtDataTracker, exiting...")
		return errors.New("got nil as MissingPvtDataTracker, exiting...")
	}
	totalReconciled, minBlock, maxBlock := 0, uint64(math.MaxUint64), uint64(0)

	for {
		missingPvtDataInfo, err := missingPvtDataTracker.GetMissingPvtDataInfoForMostRecentBlocks(r.config.batchSize)
		if err != nil {
			logger.Error("reconciliation error when trying to get missing pvt data info recent blocks:", err)
			return err
		}
		
		if len(missingPvtDataInfo) == 0 {
			if totalReconciled > 0 {
				logger.Infof("Reconciliation cycle finished successfully. reconciled %d private data keys from blocks range [%d - %d]", totalReconciled, minBlock, maxBlock)
			} else {
				logger.Debug("Reconciliation cycle finished successfully. no items to reconcile")
			}
			return nil
		}

		logger.Debug("got from ledger", len(missingPvtDataInfo), "blocks with missing private data, trying to reconcile...")

		dig2collectionCfg, minB, maxB := r.getDig2CollectionConfig(missingPvtDataInfo)
		fetchedData, err := r.FetchReconciledItems(dig2collectionCfg)
		if err != nil {
			logger.Error("reconciliation error when trying to fetch missing items from different peers:", err)
			return err
		}
		if len(fetchedData.AvailableElements) == 0 {
			logger.Warning("missing private data is not available on other peers")
			return nil
		}

		pvtDataToCommit := r.preparePvtDataToCommit(fetchedData.AvailableElements)
		
		pvtdataHashMismatch, err := r.CommitPvtDataOfOldBlocks(pvtDataToCommit)
		r.logMismatched(pvtdataHashMismatch)
		if err != nil {
			return errors.Wrap(err, "failed to commit private data")
		}
		if minB < minBlock {
			minBlock = minB
		}
		if maxB > maxBlock {
			maxBlock = maxB
		}
		totalReconciled += len(fetchedData.AvailableElements)
	}
}

type collectionConfigKey struct {
	chaincodeName, collectionName string
	blockNum                      uint64
}

func (r *Reconciler) getDig2CollectionConfig(missingPvtDataInfo ledger.MissingPvtDataInfo) (privdatacommon.Dig2CollectionConfig, uint64, uint64) {
	var minBlock, maxBlock uint64
	minBlock = math.MaxUint64
	maxBlock = 0
	collectionConfigCache := make(map[collectionConfigKey]*common.StaticCollectionConfig)
	dig2collectionCfg := make(map[privdatacommon.DigKey]*common.StaticCollectionConfig)
	for blockNum, blockPvtDataInfo := range missingPvtDataInfo {
		if blockNum < minBlock {
			minBlock = blockNum
		}
		if blockNum > maxBlock {
			maxBlock = blockNum
		}
		for seqInBlock, collectionPvtDataInfo := range blockPvtDataInfo {
			for _, pvtDataInfo := range collectionPvtDataInfo {
				collConfigKey := collectionConfigKey{
					chaincodeName:  pvtDataInfo.Namespace,
					collectionName: pvtDataInfo.Collection,
					blockNum:       blockNum,
				}
				if _, exists := collectionConfigCache[collConfigKey]; !exists {
					collectionConfig, err := r.getMostRecentCollectionConfig(pvtDataInfo.Namespace, pvtDataInfo.Collection, blockNum)
					if err != nil {
						logger.Debug(err)
						continue
					}
					collectionConfigCache[collConfigKey] = collectionConfig
				}
				digKey := privdatacommon.DigKey{
					SeqInBlock: seqInBlock,
					Collection: pvtDataInfo.Collection,
					Namespace:  pvtDataInfo.Namespace,
					BlockSeq:   blockNum,
				}
				dig2collectionCfg[digKey] = collectionConfigCache[collConfigKey]
			}
		}
	}
	return dig2collectionCfg, minBlock, maxBlock
}

func (r *Reconciler) getMostRecentCollectionConfig(chaincodeName string, collectionName string, blockNum uint64) (*common.StaticCollectionConfig, error) {
	configHistoryRetriever, err := r.GetConfigHistoryRetriever()
	if err != nil {
		return nil, errors.Wrap(err, "configHistoryRetriever is not available")
	}

	configInfo, err := configHistoryRetriever.MostRecentCollectionConfigBelow(blockNum, chaincodeName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot find recent collection config update below block sequence = %d for chaincode %s", blockNum, chaincodeName))
	}
	if configInfo == nil {
		return nil, errors.New(fmt.Sprintf("no collection config update below block sequence = %d for chaincode %s is available", blockNum, chaincodeName))
	}

	collectionConfig := extractCollectionConfig(configInfo.CollectionConfig, collectionName)
	if collectionConfig == nil {
		return nil, errors.New(fmt.Sprintf("no collection config was found for collection %s for chaincode %s", collectionName, chaincodeName))
	}

	staticCollectionConfig, wasCastingSuccessful := collectionConfig.Payload.(*common.CollectionConfig_StaticCollectionConfig)
	if !wasCastingSuccessful {
		return nil, errors.New(fmt.Sprintf("expected collection config of type CollectionConfig_StaticCollectionConfig for collection %s for chaincode %s, while got different config type...", collectionName, chaincodeName))
	}
	return staticCollectionConfig.StaticCollectionConfig, nil
}

func (r *Reconciler) preparePvtDataToCommit(elements []*gossip2.PvtDataElement) []*ledger.BlockPvtData {
	rwSetByBlockByKeys := r.groupRwsetByBlock(elements)

	
	var pvtDataToCommit []*ledger.BlockPvtData

	for blockNum, rwSetKeys := range rwSetByBlockByKeys {
		blockPvtData := &ledger.BlockPvtData{
			BlockNum:  blockNum,
			WriteSets: make(map[uint64]*ledger.TxPvtData),
		}
		for seqInBlock, nsRWS := range rwSetKeys.bySeqsInBlock() {
			rwsets := nsRWS.toRWSet()
			logger.Debugf("Preparing to commit [%d] private write set, missed from transaction index [%d] of block number [%d]", len(rwsets.NsPvtRwset), seqInBlock, blockNum)
			blockPvtData.WriteSets[seqInBlock] = &ledger.TxPvtData{
				SeqInBlock: seqInBlock,
				WriteSet:   rwsets,
			}
		}
		pvtDataToCommit = append(pvtDataToCommit, blockPvtData)
	}
	return pvtDataToCommit
}

func (r *Reconciler) logMismatched(pvtdataMismatched []*ledger.PvtdataHashMismatch) {
	if len(pvtdataMismatched) > 0 {
		for _, hashMismatch := range pvtdataMismatched {
			logger.Warningf("failed to reconcile pvtdata chaincode %s, collection %s, block num %d, tx num %d due to hash mismatch",
				hashMismatch.Namespace, hashMismatch.Collection, hashMismatch.BlockNum, hashMismatch.TxNum)
		}
	}
}


func (r *Reconciler) groupRwsetByBlock(elements []*gossip2.PvtDataElement) map[uint64]rwsetByKeys {
	rwSetByBlockByKeys := make(map[uint64]rwsetByKeys) 

	
	for _, element := range elements {
		dig := element.Digest
		if _, exists := rwSetByBlockByKeys[dig.BlockSeq]; !exists {
			rwSetByBlockByKeys[dig.BlockSeq] = make(map[rwSetKey][]byte)
		}
		for _, rws := range element.Payload {
			hash := hex.EncodeToString(util2.ComputeSHA256(rws))
			key := rwSetKey{
				txID:       dig.TxId,
				namespace:  dig.Namespace,
				collection: dig.Collection,
				seqInBlock: dig.SeqInBlock,
				hash:       hash,
			}
			rwSetByBlockByKeys[dig.BlockSeq][key] = rws
		}
	}
	return rwSetByBlockByKeys
}
