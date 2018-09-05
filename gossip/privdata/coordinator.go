/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	util2 "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/committer"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/transientstore"
	privdatacommon "github.com/mcc-github/blockchain/gossip/privdata/common"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protos/peer"
	transientstore2 "github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	pullRetrySleepInterval           = time.Second
	transientBlockRetentionConfigKey = "peer.gossip.pvtData.transientstoreMaxBlockRetention"
	transientBlockRetentionDefault   = 1000
)

var logger = util.GetLogger(util.LoggingPrivModule, "")


type TransientStore interface {
	
	
	PersistWithConfig(txid string, blockHeight uint64, privateSimulationResultsWithConfig *transientstore2.TxPvtReadWriteSetWithConfigInfo) error

	
	Persist(txid string, blockHeight uint64, privateSimulationResults *rwset.TxPvtReadWriteSet) error
	
	
	GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error)

	
	
	PurgeByTxids(txids []string) error

	
	
	
	
	
	
	PurgeByHeight(maxBlockNumToRetain uint64) error
}




type Coordinator interface {
	
	
	StoreBlock(block *common.Block, data util.PvtDataCollections) error

	
	StorePvtData(txid string, privData *transientstore2.TxPvtReadWriteSetWithConfigInfo, blckHeight uint64) error

	
	
	
	
	GetPvtDataAndBlockByNum(seqNum uint64, peerAuth common.SignedData) (*common.Block, util.PvtDataCollections, error)

	
	LedgerHeight() (uint64, error)

	
	Close()
}

type dig2sources map[privdatacommon.DigKey][]*peer.Endorsement

func (d2s dig2sources) keys() []privdatacommon.DigKey {
	res := make([]privdatacommon.DigKey, 0, len(d2s))
	for dig := range d2s {
		res = append(res, dig)
	}
	return res
}



type Fetcher interface {
	fetch(dig2src dig2sources) (*privdatacommon.FetchedPvtDataContainer, error)
}



type Support struct {
	privdata.CollectionStore
	txvalidator.Validator
	committer.Committer
	TransientStore
	Fetcher
}

type coordinator struct {
	selfSignedData common.SignedData
	Support
	transientBlockRetention uint64
}


func NewCoordinator(support Support, selfSignedData common.SignedData) Coordinator {
	transientBlockRetention := uint64(viper.GetInt(transientBlockRetentionConfigKey))
	if transientBlockRetention == 0 {
		logger.Warning("Configuration key", transientBlockRetentionConfigKey, "isn't set, defaulting to", transientBlockRetentionDefault)
		transientBlockRetention = transientBlockRetentionDefault
	}
	return &coordinator{Support: support, selfSignedData: selfSignedData, transientBlockRetention: transientBlockRetention}
}


func (c *coordinator) StorePvtData(txID string, privData *transientstore2.TxPvtReadWriteSetWithConfigInfo, blkHeight uint64) error {
	return c.TransientStore.PersistWithConfig(txID, blkHeight, privData)
}


func (c *coordinator) StoreBlock(block *common.Block, privateDataSets util.PvtDataCollections) error {
	if block.Data == nil {
		return errors.New("Block data is empty")
	}
	if block.Header == nil {
		return errors.New("Block header is nil")
	}
	logger.Infof("Received block [%d]", block.Header.Number)

	logger.Debugf("Validating block [%d]", block.Header.Number)
	err := c.Validator.Validate(block)
	if err != nil {
		logger.Errorf("Validation failed: %+v", err)
		return err
	}

	blockAndPvtData := &ledger.BlockAndPvtData{
		Block:        block,
		BlockPvtData: make(map[uint64]*ledger.TxPvtData),
	}

	ownedRWsets, err := computeOwnedRWsets(block, privateDataSets)
	if err != nil {
		logger.Warning("Failed computing owned RWSets", err)
		return err
	}

	privateInfo, err := c.listMissingPrivateData(block, ownedRWsets)
	if err != nil {
		logger.Warning(err)
		return err
	}

	retryThresh := viper.GetDuration("peer.gossip.pvtData.pullRetryThreshold")
	var bFetchFromPeers bool 
	if len(privateInfo.missingKeys) == 0 {
		logger.Debug("No missing collection private write sets to fetch from remote peers")
	} else {
		bFetchFromPeers = true
		logger.Debug("Could not find all collection private write sets in local peer transient store.")
		logger.Debug("Fetching", len(privateInfo.missingKeys), "collection private write sets from remote peers for a maximum duration of", retryThresh)
	}
	start := time.Now()
	limit := start.Add(retryThresh)
	for len(privateInfo.missingKeys) > 0 && time.Now().Before(limit) {
		c.fetchFromPeers(block.Header.Number, ownedRWsets, privateInfo)
		
		
		if len(privateInfo.missingKeys) == 0 {
			break
		}
		time.Sleep(pullRetrySleepInterval)
	}

	
	if bFetchFromPeers {
		if len(privateInfo.missingKeys) == 0 {
			logger.Debug("Fetched all missing collection private write sets from remote peers")
		} else {
			logger.Warning("Could not fetch all missing collection private write sets from remote peers. Will commit block with missing private write sets:", privateInfo.missingKeys)
		}
	}

	
	for seqInBlock, nsRWS := range ownedRWsets.bySeqsInBlock() {
		rwsets := nsRWS.toRWSet()
		logger.Debugf("Added %d namespace private write sets for block [%d], tran [%d]", len(rwsets.NsPvtRwset), block.Header.Number, seqInBlock)
		blockAndPvtData.BlockPvtData[seqInBlock] = &ledger.TxPvtData{
			SeqInBlock: seqInBlock,
			WriteSet:   rwsets,
		}
	}

	
	blockAndPvtData.Missing = &ledger.MissingPrivateDataList{}
	for missingRWS := range privateInfo.missingKeys {
		blockAndPvtData.Missing.Add(missingRWS.txID, missingRWS.seqInBlock, missingRWS.namespace, missingRWS.collection, true)
	}

	
	for _, missingRWS := range privateInfo.missingRWSButIneligible {
		blockAndPvtData.Missing.Add(missingRWS.txID, missingRWS.seqInBlock, missingRWS.namespace, missingRWS.collection, false)
	}

	
	err = c.CommitWithPvtData(blockAndPvtData)
	if err != nil {
		return errors.Wrap(err, "commit failed")
	}

	if len(blockAndPvtData.BlockPvtData) > 0 {
		
		if err := c.PurgeByTxids(privateInfo.txns); err != nil {
			logger.Error("Purging transactions", privateInfo.txns, "failed:", err)
		}
	}

	seq := block.Header.Number
	if seq%c.transientBlockRetention == 0 && seq > c.transientBlockRetention {
		err := c.PurgeByHeight(seq - c.transientBlockRetention)
		if err != nil {
			logger.Error("Failed purging data from transient store at block", seq, ":", err)
		}
	}

	return nil
}

func (c *coordinator) fetchFromPeers(blockSeq uint64, ownedRWsets map[rwSetKey][]byte, privateInfo *privateDataInfo) {
	dig2src := make(map[privdatacommon.DigKey][]*peer.Endorsement)
	privateInfo.missingKeys.foreach(func(k rwSetKey) {
		logger.Debug("Fetching", k, "from peers")
		dig := privdatacommon.DigKey{
			TxId:       k.txID,
			SeqInBlock: k.seqInBlock,
			Collection: k.collection,
			Namespace:  k.namespace,
			BlockSeq:   blockSeq,
		}
		dig2src[dig] = privateInfo.sources[k]
	})
	fetchedData, err := c.fetch(dig2src)
	if err != nil {
		logger.Warning("Failed fetching private data for block", blockSeq, "from peers:", err)
		return
	}

	
	for _, element := range fetchedData.AvailableElements {
		dig := element.Digest
		for _, rws := range element.Payload {
			hash := hex.EncodeToString(util2.ComputeSHA256(rws))
			key := rwSetKey{
				txID:       dig.TxId,
				namespace:  dig.Namespace,
				collection: dig.Collection,
				seqInBlock: dig.SeqInBlock,
				hash:       hash,
			}
			if _, isMissing := privateInfo.missingKeys[key]; !isMissing {
				logger.Debug("Ignoring", key, "because it wasn't found in the block")
				continue
			}
			ownedRWsets[key] = rws
			delete(privateInfo.missingKeys, key)
			
			
			c.TransientStore.Persist(dig.TxId, blockSeq, key.toTxPvtReadWriteSet(rws))
			logger.Debug("Fetched", key)
		}
	}
	
	for _, dig := range fetchedData.PurgedElements {
		
		for missingPvtRWKey := range privateInfo.missingKeys {
			if missingPvtRWKey.namespace == dig.Namespace &&
				missingPvtRWKey.collection == dig.Collection &&
				missingPvtRWKey.txID == dig.TxId {
				delete(privateInfo.missingKeys, missingPvtRWKey)
				logger.Warningf("Missing key because was purged or will soon be purged, "+
					"continue block commit without [%+v] in private rwset", missingPvtRWKey)
			}
		}
	}
}

func (c *coordinator) fetchMissingFromTransientStore(missing rwSetKeysByTxIDs, ownedRWsets map[rwSetKey][]byte) {
	
	for txAndSeq, filter := range missing.FiltersByTxIDs() {
		c.fetchFromTransientStore(txAndSeq, filter, ownedRWsets)
	}
}

func (c *coordinator) fetchFromTransientStore(txAndSeq txAndSeqInBlock, filter ledger.PvtNsCollFilter, ownedRWsets map[rwSetKey][]byte) {
	iterator, err := c.TransientStore.GetTxPvtRWSetByTxid(txAndSeq.txID, filter)
	if err != nil {
		logger.Warning("Failed obtaining iterator from transient store:", err)
		return
	}
	defer iterator.Close()
	for {
		res, err := iterator.NextWithConfig()
		if err != nil {
			logger.Error("Failed iterating:", err)
			break
		}
		if res == nil {
			
			break
		}
		if res.PvtSimulationResultsWithConfig == nil {
			logger.Warning("Resultset's PvtSimulationResultsWithConfig for", txAndSeq.txID, "is nil, skipping")
			continue
		}
		simRes := res.PvtSimulationResultsWithConfig
		if simRes.PvtRwset == nil {
			logger.Warning("The PvtRwset of PvtSimulationResultsWithConfig for", txAndSeq.txID, "is nil, skipping")
			continue
		}
		for _, ns := range simRes.PvtRwset.NsPvtRwset {
			for _, col := range ns.CollectionPvtRwset {
				key := rwSetKey{
					txID:       txAndSeq.txID,
					seqInBlock: txAndSeq.seqInBlock,
					collection: col.CollectionName,
					namespace:  ns.Namespace,
					hash:       hex.EncodeToString(util2.ComputeSHA256(col.Rwset)),
				}
				
				ownedRWsets[key] = col.Rwset
			} 
		} 
	} 
}


func computeOwnedRWsets(block *common.Block, blockPvtData util.PvtDataCollections) (rwsetByKeys, error) {
	lastBlockSeq := len(block.Data.Data) - 1

	ownedRWsets := make(map[rwSetKey][]byte)
	for _, txPvtData := range blockPvtData {
		if lastBlockSeq < int(txPvtData.SeqInBlock) {
			logger.Warningf("Claimed SeqInBlock %d but block has only %d transactions", txPvtData.SeqInBlock, len(block.Data.Data))
			continue
		}
		env, err := utils.GetEnvelopeFromBlock(block.Data.Data[txPvtData.SeqInBlock])
		if err != nil {
			return nil, err
		}
		payload, err := utils.GetPayload(env)
		if err != nil {
			return nil, err
		}

		chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			return nil, err
		}
		for _, ns := range txPvtData.WriteSet.NsPvtRwset {
			for _, col := range ns.CollectionPvtRwset {
				computedHash := hex.EncodeToString(util2.ComputeSHA256(col.Rwset))
				ownedRWsets[rwSetKey{
					txID:       chdr.TxId,
					seqInBlock: txPvtData.SeqInBlock,
					collection: col.CollectionName,
					namespace:  ns.Namespace,
					hash:       computedHash,
				}] = col.Rwset
			} 
		} 
	} 
	return ownedRWsets, nil
}

type readWriteSets []readWriteSet

func (s readWriteSets) toRWSet() *rwset.TxPvtReadWriteSet {
	namespaces := make(map[string]*rwset.NsPvtReadWriteSet)
	dataModel := rwset.TxReadWriteSet_KV
	for _, rws := range s {
		if _, exists := namespaces[rws.namespace]; !exists {
			namespaces[rws.namespace] = &rwset.NsPvtReadWriteSet{
				Namespace: rws.namespace,
			}
		}
		col := &rwset.CollectionPvtReadWriteSet{
			CollectionName: rws.collection,
			Rwset:          rws.rws,
		}
		namespaces[rws.namespace].CollectionPvtRwset = append(namespaces[rws.namespace].CollectionPvtRwset, col)
	}

	var namespaceSlice []*rwset.NsPvtReadWriteSet
	for _, nsRWset := range namespaces {
		namespaceSlice = append(namespaceSlice, nsRWset)
	}

	return &rwset.TxPvtReadWriteSet{
		DataModel:  dataModel,
		NsPvtRwset: namespaceSlice,
	}
}

type readWriteSet struct {
	rwSetKey
	rws []byte
}

type rwsetByKeys map[rwSetKey][]byte

func (s rwsetByKeys) bySeqsInBlock() map[uint64]readWriteSets {
	res := make(map[uint64]readWriteSets)
	for k, rws := range s {
		res[k.seqInBlock] = append(res[k.seqInBlock], readWriteSet{
			rws:      rws,
			rwSetKey: k,
		})
	}
	return res
}

type rwsetKeys map[rwSetKey]struct{}


func (s rwsetKeys) String() string {
	var buffer bytes.Buffer
	for k := range s {
		buffer.WriteString(fmt.Sprintf("%s\n", k.String()))
	}
	return buffer.String()
}


func (s rwsetKeys) foreach(f func(key rwSetKey)) {
	for k := range s {
		f(k)
	}
}


func (s rwsetKeys) exclude(exists func(key rwSetKey) bool) {
	for k := range s {
		if exists(k) {
			delete(s, k)
		}
	}
}

type txAndSeqInBlock struct {
	txID       string
	seqInBlock uint64
}

type rwSetKeysByTxIDs map[txAndSeqInBlock][]rwSetKey

func (s rwSetKeysByTxIDs) flatten() rwsetKeys {
	m := make(map[rwSetKey]struct{})
	for _, keys := range s {
		for _, k := range keys {
			m[k] = struct{}{}
		}
	}
	return m
}

func (s rwSetKeysByTxIDs) FiltersByTxIDs() map[txAndSeqInBlock]ledger.PvtNsCollFilter {
	filters := make(map[txAndSeqInBlock]ledger.PvtNsCollFilter)
	for txAndSeq, rwsKeys := range s {
		filter := ledger.NewPvtNsCollFilter()
		for _, rwskey := range rwsKeys {
			filter.Add(rwskey.namespace, rwskey.collection)
		}
		filters[txAndSeq] = filter
	}

	return filters
}

type rwSetKey struct {
	txID       string
	seqInBlock uint64
	namespace  string
	collection string
	hash       string
}


func (k *rwSetKey) String() string {
	return fmt.Sprintf("txID: %s, seq: %d, namespace: %s, collection: %s, hash: %s", k.txID, k.seqInBlock, k.namespace, k.collection, k.hash)
}

func (k *rwSetKey) toTxPvtReadWriteSet(rws []byte) *rwset.TxPvtReadWriteSet {
	return &rwset.TxPvtReadWriteSet{
		DataModel: rwset.TxReadWriteSet_KV,
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{
			{
				Namespace: k.namespace,
				CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
					{
						CollectionName: k.collection,
						Rwset:          rws,
					},
				},
			},
		},
	}
}

type txns []string
type blockData [][]byte
type blockConsumer func(seqInBlock uint64, chdr *common.ChannelHeader, txRWSet *rwsetutil.TxRwSet, endorsers []*peer.Endorsement)

func (data blockData) forEachTxn(txsFilter txValidationFlags, consumer blockConsumer) txns {
	var txList []string
	for seqInBlock, envBytes := range data {
		env, err := utils.GetEnvelopeFromBlock(envBytes)
		if err != nil {
			logger.Warning("Invalid envelope:", err)
			continue
		}

		payload, err := utils.GetPayload(env)
		if err != nil {
			logger.Warning("Invalid payload:", err)
			continue
		}

		chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			logger.Warning("Invalid channel header:", err)
			continue
		}

		if chdr.Type != int32(common.HeaderType_ENDORSER_TRANSACTION) {
			continue
		}

		txList = append(txList, chdr.TxId)

		if txsFilter[seqInBlock] != uint8(peer.TxValidationCode_VALID) {
			logger.Debug("Skipping Tx", seqInBlock, "because it's invalid. Status is", txsFilter[seqInBlock])
			continue
		}

		respPayload, err := utils.GetActionFromEnvelope(envBytes)
		if err != nil {
			logger.Warning("Failed obtaining action from envelope", err)
			continue
		}

		tx, err := utils.GetTransaction(payload.Data)
		if err != nil {
			logger.Warning("Invalid transaction in payload data for tx ", chdr.TxId, ":", err)
			continue
		}

		ccActionPayload, err := utils.GetChaincodeActionPayload(tx.Actions[0].Payload)
		if err != nil {
			logger.Warning("Invalid chaincode action in payload for tx", chdr.TxId, ":", err)
			continue
		}

		if ccActionPayload.Action == nil {
			logger.Warning("Action in ChaincodeActionPayload for", chdr.TxId, "is nil")
			continue
		}

		txRWSet := &rwsetutil.TxRwSet{}
		if err = txRWSet.FromProtoBytes(respPayload.Results); err != nil {
			logger.Warning("Failed obtaining TxRwSet from ChaincodeAction's results", err)
			continue
		}
		consumer(uint64(seqInBlock), chdr, txRWSet, ccActionPayload.Action.Endorsements)
	}
	return txList
}

func endorsersFromOrgs(ns string, col string, endorsers []*peer.Endorsement, orgs []string) []*peer.Endorsement {
	var res []*peer.Endorsement
	for _, e := range endorsers {
		sID := &msp.SerializedIdentity{}
		err := proto.Unmarshal(e.Endorser, sID)
		if err != nil {
			logger.Warning("Failed unmarshalling endorser:", err)
			continue
		}
		if !util.Contains(sID.Mspid, orgs) {
			logger.Debug(sID.Mspid, "isn't among the collection's orgs:", orgs, "for namespace", ns, ",collection", col)
			continue
		}
		res = append(res, e)
	}
	return res
}

type privateDataInfo struct {
	sources                 map[rwSetKey][]*peer.Endorsement
	missingKeysByTxIDs      rwSetKeysByTxIDs
	missingKeys             rwsetKeys
	txns                    txns
	missingRWSButIneligible []rwSetKey
}


func (c *coordinator) listMissingPrivateData(block *common.Block, ownedRWsets map[rwSetKey][]byte) (*privateDataInfo, error) {
	if block.Metadata == nil || len(block.Metadata.Metadata) <= int(common.BlockMetadataIndex_TRANSACTIONS_FILTER) {
		return nil, errors.New("Block.Metadata is nil or Block.Metadata lacks a Tx filter bitmap")
	}
	txsFilter := txValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	if len(txsFilter) != len(block.Data.Data) {
		return nil, errors.Errorf("Block data size(%d) is different from Tx filter size(%d)", len(block.Data.Data), len(txsFilter))
	}

	sources := make(map[rwSetKey][]*peer.Endorsement)
	privateRWsetsInBlock := make(map[rwSetKey]struct{})
	missing := make(rwSetKeysByTxIDs)
	data := blockData(block.Data.Data)
	bi := &transactionInspector{
		sources:              sources,
		missingKeys:          missing,
		ownedRWsets:          ownedRWsets,
		privateRWsetsInBlock: privateRWsetsInBlock,
		coordinator:          c,
	}
	txList := data.forEachTxn(txsFilter, bi.inspectTransaction)

	privateInfo := &privateDataInfo{
		sources:            sources,
		missingKeysByTxIDs: missing,
		txns:               txList,
		missingRWSButIneligible: bi.missingRWSButIneligible,
	}

	logger.Debug("Retrieving private write sets for", len(privateInfo.missingKeysByTxIDs), "transactions from transient store")

	
	c.fetchMissingFromTransientStore(privateInfo.missingKeysByTxIDs, ownedRWsets)

	
	
	for k := range ownedRWsets {
		if _, exists := privateRWsetsInBlock[k]; !exists {
			logger.Warning("Removed", k.namespace, k.collection, "hash", k.hash, "from the data passed to the ledger")
			delete(ownedRWsets, k)
		}
	}

	privateInfo.missingKeys = privateInfo.missingKeysByTxIDs.flatten()
	
	privateInfo.missingKeys.exclude(func(key rwSetKey) bool {
		_, exists := ownedRWsets[key]
		return exists
	})

	return privateInfo, nil
}

type transactionInspector struct {
	*coordinator
	privateRWsetsInBlock    map[rwSetKey]struct{}
	missingKeys             rwSetKeysByTxIDs
	sources                 map[rwSetKey][]*peer.Endorsement
	ownedRWsets             map[rwSetKey][]byte
	missingRWSButIneligible []rwSetKey
}

func (bi *transactionInspector) inspectTransaction(seqInBlock uint64, chdr *common.ChannelHeader, txRWSet *rwsetutil.TxRwSet, endorsers []*peer.Endorsement) {
	for _, ns := range txRWSet.NsRwSets {
		for _, hashedCollection := range ns.CollHashedRwSets {
			if !containsWrites(chdr.TxId, ns.NameSpace, hashedCollection) {
				continue
			}
			
			
			
			
			
			policy := bi.accessPolicyForCollection(chdr, ns.NameSpace, hashedCollection.CollectionName)
			if policy == nil {
				logger.Errorf("Failed to retrieve collection config for channel [%s], chaincode [%s], collection name [%s] for txID [%s]. Skipping.",
					chdr.ChannelId, ns.NameSpace, hashedCollection.CollectionName, chdr.TxId)
				continue
			}

			key := rwSetKey{
				txID:       chdr.TxId,
				seqInBlock: seqInBlock,
				hash:       hex.EncodeToString(hashedCollection.PvtRwSetHash),
				namespace:  ns.NameSpace,
				collection: hashedCollection.CollectionName,
			}

			if !bi.isEligible(policy, ns.NameSpace, hashedCollection.CollectionName) {
				logger.Debugf("Peer is not eligible for collection, channel [%s], chaincode [%s], "+
					"collection name [%s], txID [%s] the policy is [%#v]. Skipping.",
					chdr.ChannelId, ns.NameSpace, hashedCollection.CollectionName, chdr.TxId, policy)
				bi.missingRWSButIneligible = append(bi.missingRWSButIneligible, key)
				continue
			}

			bi.privateRWsetsInBlock[key] = struct{}{}
			if _, exists := bi.ownedRWsets[key]; !exists {
				txAndSeq := txAndSeqInBlock{
					txID:       chdr.TxId,
					seqInBlock: seqInBlock,
				}
				bi.missingKeys[txAndSeq] = append(bi.missingKeys[txAndSeq], key)
				bi.sources[key] = endorsersFromOrgs(ns.NameSpace, hashedCollection.CollectionName, endorsers, policy.MemberOrgs())
			}
		} 
	} 
}



func (c *coordinator) accessPolicyForCollection(chdr *common.ChannelHeader, namespace string, col string) privdata.CollectionAccessPolicy {
	cp := common.CollectionCriteria{
		Channel:    chdr.ChannelId,
		Namespace:  namespace,
		Collection: col,
		TxId:       chdr.TxId,
	}
	sp, err := c.CollectionStore.RetrieveCollectionAccessPolicy(cp)
	if err != nil {
		logger.Warning("Failed obtaining policy for", cp, ":", err, "skipping collection")
		return nil
	}
	return sp
}


func (c *coordinator) isEligible(ap privdata.CollectionAccessPolicy, namespace string, col string) bool {
	filt := ap.AccessFilter()
	eligible := filt(c.selfSignedData)
	if !eligible {
		logger.Debug("Skipping namespace", namespace, "collection", col, "because we're not eligible for the private data")
	}
	return eligible
}

type seqAndDataModel struct {
	seq       uint64
	dataModel rwset.TxReadWriteSet_DataModel
}



type aggregatedCollections map[seqAndDataModel]map[string][]*rwset.CollectionPvtReadWriteSet

func (ac aggregatedCollections) addCollection(seqInBlock uint64, dm rwset.TxReadWriteSet_DataModel, namespace string, col *rwset.CollectionPvtReadWriteSet) {
	seq := seqAndDataModel{
		dataModel: dm,
		seq:       seqInBlock,
	}
	if _, exists := ac[seq]; !exists {
		ac[seq] = make(map[string][]*rwset.CollectionPvtReadWriteSet)
	}
	ac[seq][namespace] = append(ac[seq][namespace], col)
}

func (ac aggregatedCollections) asPrivateData() []*ledger.TxPvtData {
	var data []*ledger.TxPvtData
	for seq, ns := range ac {
		txPrivateData := &ledger.TxPvtData{
			SeqInBlock: seq.seq,
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: seq.dataModel,
			},
		}
		for namespaceName, cols := range ns {
			txPrivateData.WriteSet.NsPvtRwset = append(txPrivateData.WriteSet.NsPvtRwset, &rwset.NsPvtReadWriteSet{
				Namespace:          namespaceName,
				CollectionPvtRwset: cols,
			})
		}
		data = append(data, txPrivateData)
	}
	return data
}





func (c *coordinator) GetPvtDataAndBlockByNum(seqNum uint64, peerAuthInfo common.SignedData) (*common.Block, util.PvtDataCollections, error) {
	blockAndPvtData, err := c.Committer.GetPvtDataAndBlockByNum(seqNum)
	if err != nil {
		return nil, nil, err
	}

	seqs2Namespaces := aggregatedCollections(make(map[seqAndDataModel]map[string][]*rwset.CollectionPvtReadWriteSet))
	data := blockData(blockAndPvtData.Block.Data.Data)
	data.forEachTxn(make(txValidationFlags, len(data)), func(seqInBlock uint64, chdr *common.ChannelHeader, txRWSet *rwsetutil.TxRwSet, _ []*peer.Endorsement) {
		item, exists := blockAndPvtData.BlockPvtData[seqInBlock]
		if !exists {
			return
		}

		for _, ns := range item.WriteSet.NsPvtRwset {
			for _, col := range ns.CollectionPvtRwset {
				cc := common.CollectionCriteria{
					Channel:    chdr.ChannelId,
					TxId:       chdr.TxId,
					Namespace:  ns.Namespace,
					Collection: col.CollectionName,
				}
				sp, err := c.CollectionStore.RetrieveCollectionAccessPolicy(cc)
				if err != nil {
					logger.Warning("Failed obtaining policy for", cc, ":", err)
					continue
				}
				isAuthorized := sp.AccessFilter()
				if isAuthorized == nil {
					logger.Warning("Failed obtaining filter for", cc)
					continue
				}
				if !isAuthorized(peerAuthInfo) {
					logger.Debug("Skipping", cc, "because peer isn't authorized")
					continue
				}
				seqs2Namespaces.addCollection(seqInBlock, item.WriteSet.DataModel, ns.Namespace, col)
			}
		}
	})

	return blockAndPvtData.Block, seqs2Namespaces.asPrivateData(), nil
}


func containsWrites(txID string, namespace string, colHashedRWSet *rwsetutil.CollHashedRwSet) bool {
	if colHashedRWSet.HashedRwSet == nil {
		logger.Warningf("HashedRWSet of tx %s, namespace %s, collection %s is nil", txID, namespace, colHashedRWSet.CollectionName)
		return false
	}
	if len(colHashedRWSet.HashedRwSet.HashedWrites) == 0 {
		logger.Debugf("HashedRWSet of tx %s, namespace %s, collection %s doesn't contain writes", txID, namespace, colHashedRWSet.CollectionName)
		return false
	}
	return true
}
