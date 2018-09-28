/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valimpl

import (
	"bytes"
	"fmt"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/customtx"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/internal"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
)




func validateAndPreparePvtBatch(block *internal.Block, db privacyenabledstate.DB,
	pubAndHashUpdates *internal.PubAndHashUpdates, pvtdata map[uint64]*ledger.TxPvtData) (*privacyenabledstate.PvtUpdateBatch, error) {
	pvtUpdates := privacyenabledstate.NewPvtUpdateBatch()
	metadataUpdates := metadataUpdates{}
	for _, tx := range block.Txs {
		if tx.ValidationCode != peer.TxValidationCode_VALID {
			continue
		}
		if !tx.ContainsPvtWrites() {
			continue
		}
		txPvtdata := pvtdata[uint64(tx.IndexInBlock)]
		if txPvtdata == nil {
			continue
		}
		if requiresPvtdataValidation(txPvtdata) {
			if err := validatePvtdata(tx, txPvtdata); err != nil {
				return nil, err
			}
		}
		var pvtRWSet *rwsetutil.TxPvtRwSet
		var err error
		if pvtRWSet, err = rwsetutil.TxPvtRwSetFromProtoMsg(txPvtdata.WriteSet); err != nil {
			return nil, err
		}
		addPvtRWSetToPvtUpdateBatch(pvtRWSet, pvtUpdates, version.NewHeight(block.Num, uint64(tx.IndexInBlock)))
		addEntriesToMetadataUpdates(metadataUpdates, pvtRWSet)
	}
	if err := incrementPvtdataVersionIfNeeded(metadataUpdates, pvtUpdates, pubAndHashUpdates, db); err != nil {
		return nil, err
	}
	return pvtUpdates, nil
}





func requiresPvtdataValidation(tx *ledger.TxPvtData) bool {
	return true
}



func validatePvtdata(tx *internal.Transaction, pvtdata *ledger.TxPvtData) error {
	if pvtdata.WriteSet == nil {
		return nil
	}

	for _, nsPvtdata := range pvtdata.WriteSet.NsPvtRwset {
		for _, collPvtdata := range nsPvtdata.CollectionPvtRwset {
			collPvtdataHash := util.ComputeHash(collPvtdata.Rwset)
			hashInPubdata := tx.RetrieveHash(nsPvtdata.Namespace, collPvtdata.CollectionName)
			if !bytes.Equal(collPvtdataHash, hashInPubdata) {
				return &validator.ErrPvtdataHashMissmatch{
					Msg: fmt.Sprintf(`Hash of pvt data for collection [%s:%s] does not match with the corresponding hash in the public data.
					public hash = [%#v], pvt data hash = [%#v]`, nsPvtdata.Namespace, collPvtdata.CollectionName, hashInPubdata, collPvtdataHash),
				}
			}
		}
	}
	return nil
}



func preprocessProtoBlock(txmgr txmgr.TxMgr, validateKVFunc func(key string, value []byte) error,
	block *common.Block, doMVCCValidation bool) (*internal.Block, error) {
	b := &internal.Block{Num: block.Header.Number}
	
	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	for txIndex, envBytes := range block.Data.Data {
		var env *common.Envelope
		var chdr *common.ChannelHeader
		var payload *common.Payload
		var err error
		if env, err = utils.GetEnvelopeFromBlock(envBytes); err == nil {
			if payload, err = utils.GetPayload(env); err == nil {
				chdr, err = utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
			}
		}
		if txsFilter.IsInvalid(txIndex) {
			
			logger.Warningf("Channel [%s]: Block [%d] Transaction index [%d] TxId [%s]"+
				" marked as invalid by committer. Reason code [%s]",
				chdr.GetChannelId(), block.Header.Number, txIndex, chdr.GetTxId(),
				txsFilter.Flag(txIndex).String())
			continue
		}
		if err != nil {
			return nil, err
		}

		var txRWSet *rwsetutil.TxRwSet
		txType := common.HeaderType(chdr.Type)
		logger.Debugf("txType=%s", txType)
		if txType == common.HeaderType_ENDORSER_TRANSACTION {
			
			respPayload, err := utils.GetActionFromEnvelope(envBytes)
			if err != nil {
				txsFilter.SetFlag(txIndex, peer.TxValidationCode_NIL_TXACTION)
				continue
			}
			txRWSet = &rwsetutil.TxRwSet{}
			if err = txRWSet.FromProtoBytes(respPayload.Results); err != nil {
				txsFilter.SetFlag(txIndex, peer.TxValidationCode_INVALID_OTHER_REASON)
				continue
			}
		} else {
			rwsetProto, err := processNonEndorserTx(env, chdr.TxId, txType, txmgr, !doMVCCValidation)
			if _, ok := err.(*customtx.InvalidTxError); ok {
				txsFilter.SetFlag(txIndex, peer.TxValidationCode_INVALID_OTHER_REASON)
				continue
			}
			if err != nil {
				return nil, err
			}
			if rwsetProto != nil {
				if txRWSet, err = rwsetutil.TxRwSetFromProtoMsg(rwsetProto); err != nil {
					return nil, err
				}
			}
		}
		if txRWSet != nil {
			if err := validateWriteset(txRWSet, validateKVFunc); err != nil {
				logger.Warningf("Channel [%s]: Block [%d] Transaction index [%d] TxId [%s]"+
					" marked as invalid. Reason code [%s]",
					chdr.GetChannelId(), block.Header.Number, txIndex, chdr.GetTxId(), peer.TxValidationCode_INVALID_WRITESET)
				txsFilter.SetFlag(txIndex, peer.TxValidationCode_INVALID_WRITESET)
				continue
			}
			b.Txs = append(b.Txs, &internal.Transaction{IndexInBlock: txIndex, ID: chdr.TxId, RWSet: txRWSet})
		}
	}
	return b, nil
}

func processNonEndorserTx(txEnv *common.Envelope, txid string, txType common.HeaderType, txmgr txmgr.TxMgr, synchingState bool) (*rwset.TxReadWriteSet, error) {
	logger.Debugf("Performing custom processing for transaction [txid=%s], [txType=%s]", txid, txType)
	processor := customtx.GetProcessor(txType)
	logger.Debugf("Processor for custom tx processing:%#v", processor)
	if processor == nil {
		return nil, nil
	}

	var err error
	var sim ledger.TxSimulator
	var simRes *ledger.TxSimulationResults
	if sim, err = txmgr.NewTxSimulator(txid); err != nil {
		return nil, err
	}
	defer sim.Done()
	if err = processor.GenerateSimulationResults(txEnv, sim, synchingState); err != nil {
		return nil, err
	}
	if simRes, err = sim.GetTxSimulationResults(); err != nil {
		return nil, err
	}
	return simRes.PubSimulationResults, nil
}

func validateWriteset(txRWSet *rwsetutil.TxRwSet, validateKVFunc func(key string, value []byte) error) error {
	for _, nsRwSet := range txRWSet.NsRwSets {
		pubWriteset := nsRwSet.KvRwSet
		if pubWriteset == nil {
			continue
		}
		for _, kvwrite := range pubWriteset.Writes {
			if err := validateKVFunc(kvwrite.Key, kvwrite.Value); err != nil {
				return err
			}
		}
	}
	return nil
}


func postprocessProtoBlock(block *common.Block, validatedBlock *internal.Block) {
	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	for _, tx := range validatedBlock.Txs {
		txsFilter.SetFlag(tx.IndexInBlock, tx.ValidationCode)
	}
	block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
}

func addPvtRWSetToPvtUpdateBatch(pvtRWSet *rwsetutil.TxPvtRwSet, pvtUpdateBatch *privacyenabledstate.PvtUpdateBatch, ver *version.Height) {
	for _, ns := range pvtRWSet.NsPvtRwSet {
		for _, coll := range ns.CollPvtRwSets {
			for _, kvwrite := range coll.KvRwSet.Writes {
				if !kvwrite.IsDelete {
					pvtUpdateBatch.Put(ns.NameSpace, coll.CollectionName, kvwrite.Key, kvwrite.Value, ver)
				} else {
					pvtUpdateBatch.Delete(ns.NameSpace, coll.CollectionName, kvwrite.Key, ver)
				}
			}
		}
	}
}






func incrementPvtdataVersionIfNeeded(
	metadataUpdates metadataUpdates,
	pvtUpdateBatch *privacyenabledstate.PvtUpdateBatch,
	pubAndHashUpdates *internal.PubAndHashUpdates,
	db privacyenabledstate.DB) error {

	for collKey := range metadataUpdates {
		ns, coll, key := collKey.ns, collKey.coll, collKey.key
		keyHash := util.ComputeStringHash(key)
		hashedVal := pubAndHashUpdates.HashUpdates.Get(ns, coll, string(keyHash))
		if hashedVal == nil {
			
			
			
			continue
		}
		latestVal, err := retrieveLatestVal(ns, coll, key, pvtUpdateBatch, db)
		if err != nil {
			return err
		}
		if latestVal == nil || 
			version.AreSame(latestVal.Version, hashedVal.Version) { 
			continue
		}
		
		
		latestValHash := util.ComputeHash(latestVal.Value)
		if bytes.Equal(latestValHash, hashedVal.Value) { 
			
			pvtUpdateBatch.Put(ns, coll, key, latestVal.Value, hashedVal.Version)
		}
	}
	return nil
}

type collKey struct {
	ns, coll, key string
}

type metadataUpdates map[collKey]bool

func addEntriesToMetadataUpdates(metadataUpdates metadataUpdates, pvtRWSet *rwsetutil.TxPvtRwSet) {
	for _, ns := range pvtRWSet.NsPvtRwSet {
		for _, coll := range ns.CollPvtRwSets {
			for _, metadataWrite := range coll.KvRwSet.MetadataWrites {
				ns, coll, key := ns.NameSpace, coll.CollectionName, metadataWrite.Key
				metadataUpdates[collKey{ns, coll, key}] = true
			}
		}
	}
}

func retrieveLatestVal(ns, coll, key string, pvtUpdateBatch *privacyenabledstate.PvtUpdateBatch,
	db privacyenabledstate.DB) (val *statedb.VersionedValue, err error) {
	val = pvtUpdateBatch.Get(ns, coll, key)
	if val == nil {
		val, err = db.GetPrivateData(ns, coll, key)
	}
	return
}
