/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	protosgossip "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/gossip/privdata/common"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/pkg/errors"
)





type DataStore interface {
	
	
	GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error)

	
	
	
	GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)

	
	GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error)

	
	LedgerHeight() (uint64, error)
}




type RWSetScanner interface {
	transientstore.RWSetScanner
}


type StorageDataRetriever interface {
	
	
	CollectionRWSet(dig []*protosgossip.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, bool, error)
}

type dataRetriever struct {
	store DataStore
}



func NewDataRetriever(store DataStore) StorageDataRetriever {
	return &dataRetriever{store: store}
}



func (dr *dataRetriever) CollectionRWSet(digests []*protosgossip.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, bool, error) {
	height, err := dr.store.LedgerHeight()
	if err != nil {
		
		return nil, false, errors.Wrap(err, "wasn't able to read ledger height")
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

		return results, false, nil
	}
	
	results, err := dr.fromLedger(digests, blockNum)
	return results, true, err
}

func (dr *dataRetriever) fromLedger(digests []*protosgossip.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, error) {
	filter := make(map[string]ledger.PvtCollFilter)
	for _, dig := range digests {
		if _, ok := filter[dig.Namespace]; !ok {
			filter[dig.Namespace] = make(ledger.PvtCollFilter)
		}
		filter[dig.Namespace][dig.Collection] = true
	}

	pvtData, err := dr.store.GetPvtDataByNum(blockNum, filter)
	if err != nil {
		return nil, errors.Errorf("wasn't able to obtain private data, block sequence number %d, due to %s", blockNum, err)
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
			return nil, errors.Errorf("cannot obtain configuration history retriever, for collection <%s>"+
				" txID <%s> block sequence number <%d> due to <%s>", dig.Collection, dig.TxId, dig.BlockSeq, err)
		}

		configInfo, err := confHistoryRetriever.MostRecentCollectionConfigBelow(dig.BlockSeq, dig.Namespace)
		if err != nil {
			return nil, errors.Errorf("cannot find recent collection config update below block sequence = %d,"+
				" collection name = <%s> for chaincode <%s>", dig.BlockSeq, dig.Collection, dig.Namespace)
		}

		if configInfo == nil {
			return nil, errors.Errorf("no collection config update below block sequence = <%d>"+
				" collection name = <%s> for chaincode <%s> is available ", dig.BlockSeq, dig.Collection, dig.Namespace)
		}
		configs := extractCollectionConfig(configInfo.CollectionConfig, dig.Collection)
		if configs == nil {
			return nil, errors.Errorf("no collection config was found for collection <%s>"+
				" namespace <%s> txID <%s>", dig.Collection, dig.Namespace, dig.TxId)
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

func (dr *dataRetriever) fromTransientStore(dig *protosgossip.PvtDataDigest, filter map[string]ledger.PvtCollFilter) (*util.PrivateRWSetWithConfig, error) {
	results := &util.PrivateRWSetWithConfig{}
	it, err := dr.store.GetTxPvtRWSetByTxid(dig.TxId, filter)
	if err != nil {
		return nil, errors.Errorf("was not able to retrieve private data from transient store, namespace <%s>"+
			", collection name %s, txID <%s>, due to <%s>", dig.Namespace, dig.Collection, dig.TxId, err)
	}
	defer it.Close()

	maxEndorsedAt := uint64(0)
	for {
		res, err := it.Next()
		if err != nil {
			return nil, errors.Errorf("error getting next element out of private data iterator, namespace <%s>"+
				", collection name <%s>, txID <%s>, due to <%s>", dig.Namespace, dig.Collection, dig.TxId, err)
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
			logger.Debug("Received private data namespace ", nsws.Namespace, " instead of ", namespace, " skipping...")
			continue
		}
		for _, col := range nsws.CollectionPvtRwset {
			
			if col.CollectionName != collectionName {
				logger.Debug("Received private data collection ", col.CollectionName, " instead of ", collectionName, " skipping...")
				continue
			}
			
			pRWsets = append(pRWsets, col.Rwset)
		}
	}

	return pRWsets
}
