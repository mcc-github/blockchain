/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"runtime/debug"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/deliver"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("common.deliverevents")



type PolicyCheckerProvider func(resourceName string) deliver.PolicyCheckerFunc


type DeliverServer struct {
	DeliverHandler          *deliver.Handler
	PolicyCheckerProvider   PolicyCheckerProvider
	CollectionPolicyChecker CollectionPolicyChecker
	IdentityDeserializerMgr IdentityDeserializerManager
}


type Chain interface {
	deliver.Chain
	Ledger() ledger.PeerLedger
}


type CollectionPolicyChecker interface {
	CheckCollectionPolicy(blockNum uint64, ccName string, collName string, cfgHistoryRetriever ledger.ConfigHistoryRetriever, deserializer msp.IdentityDeserializer, signedData *protoutil.SignedData) (bool, error)
}


type IdentityDeserializerManager interface {
	
	
	Deserializer(channel string) (msp.IdentityDeserializer, error)
}


type blockResponseSender struct {
	peer.Deliver_DeliverServer
}


func (brs *blockResponseSender) SendStatusResponse(status common.Status) error {
	reply := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_Status{Status: status},
	}
	return brs.Send(reply)
}


func (brs *blockResponseSender) SendBlockResponse(
	block *common.Block,
	channelID string,
	chain deliver.Chain,
	signedData *protoutil.SignedData,
) error {
	
	response := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_Block{Block: block},
	}
	return brs.Send(response)
}

func (brs *blockResponseSender) DataType() string {
	return "block"
}


type filteredBlockResponseSender struct {
	peer.Deliver_DeliverFilteredServer
}


func (fbrs *filteredBlockResponseSender) SendStatusResponse(status common.Status) error {
	response := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_Status{Status: status},
	}
	return fbrs.Send(response)
}



func (fbrs *filteredBlockResponseSender) IsFiltered() bool {
	return true
}


func (fbrs *filteredBlockResponseSender) SendBlockResponse(
	block *common.Block,
	channelID string,
	chain deliver.Chain,
	signedData *protoutil.SignedData,
) error {
	
	b := blockEvent(*block)
	filteredBlock, err := b.toFilteredBlock()
	if err != nil {
		logger.Warningf("Failed to generate filtered block due to: %s", err)
		return fbrs.SendStatusResponse(common.Status_BAD_REQUEST)
	}
	response := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_FilteredBlock{FilteredBlock: filteredBlock},
	}
	return fbrs.Send(response)
}

func (fbrs *filteredBlockResponseSender) DataType() string {
	return "filtered_block"
}


type blockAndPrivateDataResponseSender struct {
	peer.Deliver_DeliverWithPrivateDataServer
	CollectionPolicyChecker
	IdentityDeserializerManager
}


func (bprs *blockAndPrivateDataResponseSender) SendStatusResponse(status common.Status) error {
	reply := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_Status{Status: status},
	}
	return bprs.Send(reply)
}


func (bprs *blockAndPrivateDataResponseSender) SendBlockResponse(
	block *common.Block,
	channelID string,
	chain deliver.Chain,
	signedData *protoutil.SignedData,
) error {
	pvtData, err := bprs.getPrivateData(block, chain, channelID, signedData)
	if err != nil {
		return err
	}

	blockAndPvtData := &peer.BlockAndPrivateData{
		Block:          block,
		PrivateDataMap: pvtData,
	}
	response := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_BlockAndPrivateData{BlockAndPrivateData: blockAndPvtData},
	}
	return bprs.Send(response)
}

func (bprs *blockAndPrivateDataResponseSender) DataType() string {
	return "block_and_pvtdata"
}


func (bprs *blockAndPrivateDataResponseSender) getPrivateData(
	block *common.Block,
	chain deliver.Chain,
	channelID string,
	signedData *protoutil.SignedData,
) (map[uint64]*rwset.TxPvtReadWriteSet, error) {

	channel, ok := chain.(Chain)
	if !ok {
		return nil, errors.New("wrong chain type")
	}

	pvtData, err := channel.Ledger().GetPvtDataByNum(block.Header.Number, nil)
	if err != nil {
		logger.Errorf("Error getting private data by block number %d on channel %s", block.Header.Number, channelID)
		return nil, errors.Wrapf(err, "error getting private data by block number %d", block.Header.Number)
	}

	seqs2Namespaces := aggregatedCollections(make(map[seqAndDataModel]map[string][]*rwset.CollectionPvtReadWriteSet))

	configHistoryRetriever, err := channel.Ledger().GetConfigHistoryRetriever()
	if err != nil {
		return nil, err
	}

	identityDeserializer, err := bprs.IdentityDeserializerManager.Deserializer(channelID)
	if err != nil {
		return nil, err
	}

	
	for _, item := range pvtData {
		logger.Debugf("Got private data for block number %d, tx sequence %d", block.Header.Number, item.SeqInBlock)
		if item.WriteSet == nil {
			continue
		}
		for _, ns := range item.WriteSet.NsPvtRwset {
			for _, col := range ns.CollectionPvtRwset {
				logger.Debugf("Checking policy for namespace %s, collection %s", ns.Namespace, col.CollectionName)

				eligible, err := bprs.CollectionPolicyChecker.CheckCollectionPolicy(block.Header.Number,
					ns.Namespace, col.CollectionName, configHistoryRetriever, identityDeserializer, signedData)
				if err != nil {
					return nil, err
				}

				if eligible {
					logger.Debugf("Adding private data for namespace %s, collection %s", ns.Namespace, col.CollectionName)
					seqs2Namespaces.addCollection(item.SeqInBlock, item.WriteSet.DataModel, ns.Namespace, col)
				}
			}
		}
	}

	return seqs2Namespaces.asPrivateDataMap(), nil
}


type transactionActions []*peer.TransactionAction



type blockEvent common.Block


func (s *DeliverServer) DeliverFiltered(srv peer.Deliver_DeliverFilteredServer) error {
	logger.Debugf("Starting new DeliverFiltered handler")
	defer dumpStacktraceOnPanic()
	
	deliverServer := &deliver.Server{
		Receiver:      srv,
		PolicyChecker: s.PolicyCheckerProvider(resources.Event_FilteredBlock),
		ResponseSender: &filteredBlockResponseSender{
			Deliver_DeliverFilteredServer: srv,
		},
	}
	return s.DeliverHandler.Handle(srv.Context(), deliverServer)
}


func (s *DeliverServer) Deliver(srv peer.Deliver_DeliverServer) (err error) {
	logger.Debugf("Starting new Deliver handler")
	defer dumpStacktraceOnPanic()
	
	deliverServer := &deliver.Server{
		PolicyChecker: s.PolicyCheckerProvider(resources.Event_Block),
		Receiver:      srv,
		ResponseSender: &blockResponseSender{
			Deliver_DeliverServer: srv,
		},
	}
	return s.DeliverHandler.Handle(srv.Context(), deliverServer)
}


func (s *DeliverServer) DeliverWithPrivateData(srv peer.Deliver_DeliverWithPrivateDataServer) (err error) {
	logger.Debug("Starting new DeliverWithPrivateData handler")
	defer dumpStacktraceOnPanic()
	if s.CollectionPolicyChecker == nil {
		s.CollectionPolicyChecker = &collPolicyChecker{}
	}
	if s.IdentityDeserializerMgr == nil {
		s.IdentityDeserializerMgr = &identityDeserializerMgr{}
	}
	
	deliverServer := &deliver.Server{
		PolicyChecker: s.PolicyCheckerProvider(resources.Event_Block),
		Receiver:      srv,
		ResponseSender: &blockAndPrivateDataResponseSender{
			Deliver_DeliverWithPrivateDataServer: srv,
			CollectionPolicyChecker:              s.CollectionPolicyChecker,
			IdentityDeserializerManager:          s.IdentityDeserializerMgr,
		},
	}
	err = s.DeliverHandler.Handle(srv.Context(), deliverServer)
	return err
}

func (block *blockEvent) toFilteredBlock() (*peer.FilteredBlock, error) {
	filteredBlock := &peer.FilteredBlock{
		Number: block.Header.Number,
	}

	txsFltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	for txIndex, ebytes := range block.Data.Data {
		var env *common.Envelope
		var err error

		if ebytes == nil {
			logger.Debugf("got nil data bytes for tx index %d, block num %d", txIndex, block.Header.Number)
			continue
		}

		env, err = protoutil.GetEnvelopeFromBlock(ebytes)
		if err != nil {
			logger.Errorf("error getting tx from block, %s", err)
			continue
		}

		
		payload, err := protoutil.UnmarshalPayload(env.Payload)
		if err != nil {
			return nil, errors.WithMessage(err, "could not extract payload from envelope")
		}

		if payload.Header == nil {
			logger.Debugf("transaction payload header is nil, %d, block num %d", txIndex, block.Header.Number)
			continue
		}
		chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			return nil, err
		}

		filteredBlock.ChannelId = chdr.ChannelId

		filteredTransaction := &peer.FilteredTransaction{
			Txid:             chdr.TxId,
			Type:             common.HeaderType(chdr.Type),
			TxValidationCode: txsFltr.Flag(txIndex),
		}

		if filteredTransaction.Type == common.HeaderType_ENDORSER_TRANSACTION {
			tx, err := protoutil.UnmarshalTransaction(payload.Data)
			if err != nil {
				return nil, errors.WithMessage(err, "error unmarshal transaction payload for block event")
			}

			filteredTransaction.Data, err = transactionActions(tx.Actions).toFilteredActions()
			if err != nil {
				logger.Errorf(err.Error())
				return nil, err
			}
		}

		filteredBlock.FilteredTransactions = append(filteredBlock.FilteredTransactions, filteredTransaction)
	}

	return filteredBlock, nil
}

func (ta transactionActions) toFilteredActions() (*peer.FilteredTransaction_TransactionActions, error) {
	transactionActions := &peer.FilteredTransactionActions{}
	for _, action := range ta {
		chaincodeActionPayload, err := protoutil.UnmarshalChaincodeActionPayload(action.Payload)
		if err != nil {
			return nil, errors.WithMessage(err, "error unmarshal transaction action payload for block event")
		}

		if chaincodeActionPayload.Action == nil {
			logger.Debugf("chaincode action, the payload action is nil, skipping")
			continue
		}
		propRespPayload, err := protoutil.UnmarshalProposalResponsePayload(chaincodeActionPayload.Action.ProposalResponsePayload)
		if err != nil {
			return nil, errors.WithMessage(err, "error unmarshal proposal response payload for block event")
		}

		caPayload, err := protoutil.UnmarshalChaincodeAction(propRespPayload.Extension)
		if err != nil {
			return nil, errors.WithMessage(err, "error unmarshal chaincode action for block event")
		}

		ccEvent, err := protoutil.UnmarshalChaincodeEvents(caPayload.Events)
		if err != nil {
			return nil, errors.WithMessage(err, "error unmarshal chaincode event for block event")
		}

		if ccEvent.GetChaincodeId() != "" {
			filteredAction := &peer.FilteredChaincodeAction{
				ChaincodeEvent: &peer.ChaincodeEvent{
					TxId:        ccEvent.TxId,
					ChaincodeId: ccEvent.ChaincodeId,
					EventName:   ccEvent.EventName,
				},
			}
			transactionActions.ChaincodeActions = append(transactionActions.ChaincodeActions, filteredAction)
		}
	}
	return &peer.FilteredTransaction_TransactionActions{
		TransactionActions: transactionActions,
	}, nil
}

func dumpStacktraceOnPanic() {
	func() {
		if r := recover(); r != nil {
			logger.Criticalf("Deliver client triggered panic: %s\n%s", r, debug.Stack())
		}
		logger.Debugf("Closing Deliver stream")
	}()
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



func (ac aggregatedCollections) asPrivateDataMap() map[uint64]*rwset.TxPvtReadWriteSet {
	var pvtDataMap = make(map[uint64]*rwset.TxPvtReadWriteSet)
	for seq, ns := range ac {
		
		txPvtRWSet := &rwset.TxPvtReadWriteSet{
			DataModel: seq.dataModel,
		}

		for namespaceName, cols := range ns {
			txPvtRWSet.NsPvtRwset = append(txPvtRWSet.NsPvtRwset, &rwset.NsPvtReadWriteSet{
				Namespace:          namespaceName,
				CollectionPvtRwset: cols,
			})
		}

		pvtDataMap[seq.seq] = txPvtRWSet
	}
	return pvtDataMap
}



type identityDeserializerMgr struct {
}

func (*identityDeserializerMgr) Deserializer(channelID string) (msp.IdentityDeserializer, error) {
	id, ok := mgmt.GetDeserializers()[channelID]
	if !ok {
		return nil, errors.Errorf("channel %s not found", channelID)
	}
	return id, nil
}


type collPolicyChecker struct {
}


func (cs *collPolicyChecker) CheckCollectionPolicy(
	blockNum uint64,
	ccName string,
	collName string,
	cfgHistoryRetriever ledger.ConfigHistoryRetriever,
	deserializer msp.IdentityDeserializer,
	signedData *protoutil.SignedData,
) (bool, error) {
	configInfo, err := cfgHistoryRetriever.MostRecentCollectionConfigBelow(blockNum, ccName)
	if err != nil {
		return false, errors.WithMessagef(err, "error getting most recent collection config below block sequence = %d for chaincode %s", blockNum, ccName)
	}

	staticCollConfig := extractStaticCollectionConfig(configInfo.CollectionConfig, collName)
	if staticCollConfig == nil {
		return false, errors.Errorf("no collection config was found for collection %s for chaincode %s", collName, ccName)
	}

	if !staticCollConfig.MemberOnlyRead {
		return true, nil
	}

	
	collAP := &privdata.SimpleCollection{}
	err = collAP.Setup(staticCollConfig, deserializer)
	if err != nil {
		return false, errors.WithMessagef(err, "error setting up collection  %s", staticCollConfig.Name)
	}
	logger.Debugf("got collection access policy")

	collFilter := collAP.AccessFilter()
	if collFilter == nil {
		logger.Debugf("collection %s has no access filter, skipping...", collName)
		return false, nil
	}

	eligible := collFilter(*signedData)
	return eligible, nil
}

func extractStaticCollectionConfig(configPackage *common.CollectionConfigPackage, collectionName string) *common.StaticCollectionConfig {
	for _, config := range configPackage.Config {
		switch cconf := config.Payload.(type) {
		case *common.CollectionConfig_StaticCollectionConfig:
			if cconf.StaticCollectionConfig.Name == collectionName {
				return cconf.StaticCollectionConfig
			}
		default:
			return nil
		}
	}
	return nil
}
