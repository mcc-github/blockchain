/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"errors"
	"fmt"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/gossip/privdata/common"
	"github.com/mcc-github/blockchain/gossip/util"
	gossip2 "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
)


type StorageDataRetriever interface {
	
	
	CollectionRWSet(dig []*gossip2.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, error)
}







type DataStore interface {
	
	
	GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error)

	
	
	
	GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)

	
	GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error)

	
	LedgerHeight() (uint64, error)
}

type dataRetriever struct {
	store DataStore
}



func NewDataRetriever(store DataStore) StorageDataRetriever {
	return &dataRetriever{store: store}
}



func (dr *dataRetriever) CollectionRWSet(digests []*gossip2.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, error) {
	height, err := dr.store.LedgerHeight()
	if err != nil {
		
		return nil, fmt.Errorf("wasn't able to read ledger height, due to %s", err)
	}
	if height <= blockNum {
		logger.Debug("Current ledger height ", height, "is below requested block sequence number",
			blockNum, "retrieving private data from transient store")
	}

	if height <= blockNum { 
		results := make(Dig2PvtRWSetWithConfig)
		for _, dig := range digests {
			filter := map[string]ledger.PvtCollFilter{
				dig.Namespace: map[string]bool{
					dig.Collection: true,
				},
			}
			pvtRWSet, err := dr.fromTransientStore(dig, filter)
			if err != nil {
				logger.Errorf("couldn't read from transient store private read-write set, "+
					"digest %+v, because of %s", dig, err)
				continue
			}
			results[common.DigKey{
				Namespace:  dig.Namespace,
				Collection: dig.Collection,
				TxId:       dig.TxId,
				BlockSeq:   dig.BlockSeq,
				SeqInBlock: dig.SeqInBlock,
			}] = pvtRWSet
		}

		return results, nil
	}
	
	return dr.fromLedger(digests, blockNum)
}

func (dr *dataRetriever) fromLedger(digests []*gossip2.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, error) {
	filter := make(map[string]ledger.PvtCollFilter)
	for _, dig := range digests {
		if _, ok := filter[dig.Namespace]; !ok {
			filter[dig.Namespace] = make(ledger.PvtCollFilter)
		}
		filter[dig.Namespace][dig.Collection] = true
	}

	pvtData, err := dr.store.GetPvtDataByNum(blockNum, filter)
	if err != nil {
		return nil, errors.New(fmt.Sprint("wasn't able to obtain private data, block sequence number", blockNum, "due to", err))
	}

	results := make(Dig2PvtRWSetWithConfig)
	for _, dig := range digests {
		dig := dig
		pvtRWSetWithConfig := &util.PrivateRWSetWithConfig{}
		for _, data := range pvtData {
			if data.WriteSet == nil {
				logger.Warning("Received nil write set for collection tx in block", data.SeqInBlock, "block number", blockNum)
				continue
			}

			
			
			if !data.Has(dig.Namespace, dig.Collection) || data.SeqInBlock != dig.SeqInBlock {
				continue
			}

			pvtRWSet := dr.extractPvtRWsets(data.WriteSet.NsPvtRwset, dig.Namespace, dig.Collection)
			pvtRWSetWithConfig.RWSet = append(pvtRWSetWithConfig.RWSet, pvtRWSet...)
		}

		confHistoryRetriever, err := dr.store.GetConfigHistoryRetriever()
		if err != nil {
			return nil, errors.New(fmt.Sprint("cannot obtain configuration history retriever, for collection, ", dig.Collection,
				" txID ", dig.TxId, " block sequence number ", dig.BlockSeq, " due to", err))
		}

		configInfo, err := confHistoryRetriever.MostRecentCollectionConfigBelow(dig.BlockSeq, dig.Namespace)
		if err != nil {
			return nil, errors.New(fmt.Sprint("cannot find recent collection config update below block sequence = ", dig.BlockSeq,
				" collection name = ", dig.Collection, " for chaincode ", dig.Namespace))
		}

		if configInfo == nil {
			return nil, errors.New(fmt.Sprint("no collection config update below block sequence = ", dig.BlockSeq,
				" collection name = ", dig.Collection, " for chaincode ", dig.Namespace, " is available "))
		}
		configs := extractCollectionConfig(configInfo.CollectionConfig, dig.Collection)
		if configs == nil {
			return nil, errors.New(fmt.Sprint("no collection config was found for collection ", dig.Collection,
				" namespace ", dig.Namespace, " txID ", dig.TxId))
		}
		pvtRWSetWithConfig.CollectionConfig = configs
		results[common.DigKey{
			Namespace:  dig.Namespace,
			Collection: dig.Collection,
			TxId:       dig.TxId,
			BlockSeq:   dig.BlockSeq,
			SeqInBlock: dig.SeqInBlock,
		}] = pvtRWSetWithConfig
	}

	return results, nil
}

func (dr *dataRetriever) fromTransientStore(dig *gossip2.PvtDataDigest, filter map[string]ledger.PvtCollFilter) (*util.PrivateRWSetWithConfig, error) {
	results := &util.PrivateRWSetWithConfig{}
	it, err := dr.store.GetTxPvtRWSetByTxid(dig.TxId, filter)
	if err != nil {
		return nil, errors.New(fmt.Sprint("was not able to retrieve private data from transient store, namespace", dig.Namespace,
			", collection name", dig.Collection, ", txID", dig.TxId, ", due to", err))
	}
	defer it.Close()

	maxEndorsedAt := uint64(0)
	for {
		res, err := it.NextWithConfig()
		if err != nil {
			return nil, errors.New(fmt.Sprint("error getting next element out of private data iterator, namespace", dig.Namespace,
				", collection name", dig.Collection, ", txID", dig.TxId, ", due to", err))
		}
		if res == nil {
			return results, nil
		}
		rws := res.PvtSimulationResultsWithConfig
		if rws == nil {
			logger.Debug("Skipping nil PvtSimulationResultsWithConfig received at block height", res.ReceivedAtBlockHeight)
			continue
		}
		txPvtRWSet := rws.PvtRwset
		if txPvtRWSet == nil {
			logger.Debug("Skipping empty PvtRwset of PvtSimulationResultsWithConfig received at block height", res.ReceivedAtBlockHeight)
			continue
		}

		colConfigs, found := rws.CollectionConfigs[dig.Namespace]
		if !found {
			logger.Error("No collection config was found for chaincode", dig.Namespace, "collection name",
				dig.Namespace, "txID", dig.TxId)
			continue
		}

		configs := extractCollectionConfig(colConfigs, dig.Collection)
		if configs == nil {
			logger.Error("No collection config was found for collection", dig.Collection,
				"namespace", dig.Namespace, "txID", dig.TxId)
			continue
		}

		pvtRWSet := dr.extractPvtRWsets(txPvtRWSet.NsPvtRwset, dig.Namespace, dig.Collection)
		if rws.EndorsedAt >= maxEndorsedAt {
			maxEndorsedAt = rws.EndorsedAt
			results.CollectionConfig = configs
		}
		results.RWSet = append(results.RWSet, pvtRWSet...)
	}
}

func (dr *dataRetriever) extractPvtRWsets(pvtRWSets []*rwset.NsPvtReadWriteSet, namespace string, collectionName string) []util.PrivateRWSet {
	pRWsets := []util.PrivateRWSet{}

	
	for _, nsws := range pvtRWSets {
		
		if nsws.Namespace != namespace {
			logger.Warning("Received private data namespace ", nsws.Namespace, " instead of ", namespace, " skipping...")
			continue
		}
		for _, col := range nsws.CollectionPvtRwset {
			
			if col.CollectionName != collectionName {
				logger.Warning("Received private data collection ", col.CollectionName, " instead of ", collectionName, " skipping...")
				continue
			}
			
			pRWsets = append(pRWsets, col.Rwset)
		}
	}

	return pRWsets
}
