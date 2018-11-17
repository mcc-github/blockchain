/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)

var chaincodeLogger = flogging.MustGetLogger("chaincode")



type ACLProvider interface {
	CheckACL(resName string, channelID string, idinfo interface{}) error
}


type Registry interface {
	Register(*Handler) error
	Ready(cname string)
	Failed(cname string, err error)
	Deregister(cname string) error
}


type Invoker interface {
	Invoke(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, spec *pb.ChaincodeInput) (*pb.ChaincodeMessage, error)
}


type SystemCCProvider interface {
	IsSysCC(name string) bool
	IsSysCCAndNotInvokableCC2CC(name string) bool
}


type TransactionRegistry interface {
	Add(channelID, txID string) bool
	Remove(channelID, txID string)
}


type ContextRegistry interface {
	Create(txParams *ccprovider.TransactionParams) (*TransactionContext, error)
	Get(chainID, txID string) *TransactionContext
	Delete(chainID, txID string)
	Close()
}


type InstantiationPolicyChecker interface {
	CheckInstantiationPolicy(name, version string, cd *ccprovider.ChaincodeData) error
}


type CheckInstantiationPolicyFunc func(name, version string, cd *ccprovider.ChaincodeData) error

func (c CheckInstantiationPolicyFunc) CheckInstantiationPolicy(name, version string, cd *ccprovider.ChaincodeData) error {
	return c(name, version, cd)
}



type QueryResponseBuilder interface {
	BuildQueryResponse(txContext *TransactionContext, iter commonledger.ResultsIterator,
		iterID string, isPaginated bool, totalReturnLimit int32) (*pb.QueryResponse, error)
}



type ChaincodeDefinitionGetter interface {
	ChaincodeDefinition(chaincodeName string, txSim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error)
}


type LedgerGetter interface {
	GetLedger(cid string) ledger.PeerLedger
}


type UUIDGenerator interface {
	New() string
}
type UUIDGeneratorFunc func() string

func (u UUIDGeneratorFunc) New() string { return u() }


type ApplicationConfigRetriever interface {
	
	
	GetApplicationConfig(cid string) (channelconfig.Application, bool)
}


type Handler struct {
	
	Keepalive time.Duration
	
	SystemCCVersion string
	
	
	DefinitionGetter ChaincodeDefinitionGetter
	
	Invoker Invoker
	
	Registry Registry
	
	ACLProvider ACLProvider
	
	
	TXContexts ContextRegistry
	
	ActiveTransactions TransactionRegistry
	
	SystemCCProvider SystemCCProvider
	
	InstantiationPolicyChecker InstantiationPolicyChecker
	
	QueryResponseBuilder QueryResponseBuilder
	
	LedgerGetter LedgerGetter
	
	UUIDGenerator UUIDGenerator
	
	AppConfig ApplicationConfigRetriever

	
	
	state State
	
	chaincodeID *pb.ChaincodeID
	
	
	ccInstance *sysccprovider.ChaincodeInstance

	
	serialLock sync.Mutex
	
	
	chatStream ccintf.ChaincodeStream
	
	errChan chan error
	
	Metrics *Metrics
}


func (h *Handler) handleMessage(msg *pb.ChaincodeMessage) error {
	chaincodeLogger.Debugf("[%s] Fabric side handling ChaincodeMessage of type: %s in state %s", shorttxid(msg.Txid), msg.Type, h.state)

	if msg.Type == pb.ChaincodeMessage_KEEPALIVE {
		return nil
	}

	switch h.state {
	case Created:
		return h.handleMessageCreatedState(msg)
	case Ready:
		return h.handleMessageReadyState(msg)
	default:
		return errors.Errorf("handle message: invalid state %s for transaction %s", h.state, msg.Txid)
	}
}

func (h *Handler) handleMessageCreatedState(msg *pb.ChaincodeMessage) error {
	switch msg.Type {
	case pb.ChaincodeMessage_REGISTER:
		h.HandleRegister(msg)
	default:
		return fmt.Errorf("[%s] Fabric side handler cannot handle message (%s) while in created state", msg.Txid, msg.Type)
	}
	return nil
}

func (h *Handler) handleMessageReadyState(msg *pb.ChaincodeMessage) error {
	switch msg.Type {
	case pb.ChaincodeMessage_COMPLETED, pb.ChaincodeMessage_ERROR:
		h.Notify(msg)

	case pb.ChaincodeMessage_PUT_STATE:
		go h.HandleTransaction(msg, h.HandlePutState)
	case pb.ChaincodeMessage_DEL_STATE:
		go h.HandleTransaction(msg, h.HandleDelState)
	case pb.ChaincodeMessage_INVOKE_CHAINCODE:
		go h.HandleTransaction(msg, h.HandleInvokeChaincode)

	case pb.ChaincodeMessage_GET_STATE:
		go h.HandleTransaction(msg, h.HandleGetState)
	case pb.ChaincodeMessage_GET_STATE_BY_RANGE:
		go h.HandleTransaction(msg, h.HandleGetStateByRange)
	case pb.ChaincodeMessage_GET_QUERY_RESULT:
		go h.HandleTransaction(msg, h.HandleGetQueryResult)
	case pb.ChaincodeMessage_GET_HISTORY_FOR_KEY:
		go h.HandleTransaction(msg, h.HandleGetHistoryForKey)
	case pb.ChaincodeMessage_QUERY_STATE_NEXT:
		go h.HandleTransaction(msg, h.HandleQueryStateNext)
	case pb.ChaincodeMessage_QUERY_STATE_CLOSE:
		go h.HandleTransaction(msg, h.HandleQueryStateClose)

	case pb.ChaincodeMessage_GET_STATE_METADATA:
		go h.HandleTransaction(msg, h.HandleGetStateMetadata)
	case pb.ChaincodeMessage_PUT_STATE_METADATA:
		go h.HandleTransaction(msg, h.HandlePutStateMetadata)
	default:
		return fmt.Errorf("[%s] Fabric side handler cannot handle message (%s) while in ready state", msg.Txid, msg.Type)
	}

	return nil
}

type MessageHandler interface {
	Handle(*pb.ChaincodeMessage, *TransactionContext) (*pb.ChaincodeMessage, error)
}

type handleFunc func(*pb.ChaincodeMessage, *TransactionContext) (*pb.ChaincodeMessage, error)





func (h *Handler) HandleTransaction(msg *pb.ChaincodeMessage, delegate handleFunc) {
	chaincodeLogger.Debugf("[%s] handling %s from chaincode", shorttxid(msg.Txid), msg.Type.String())
	if !h.registerTxid(msg) {
		return
	}

	startTime := time.Now()
	var txContext *TransactionContext
	var err error
	if msg.Type == pb.ChaincodeMessage_INVOKE_CHAINCODE {
		txContext, err = h.getTxContextForInvoke(msg.ChannelId, msg.Txid, msg.Payload, "")
	} else {
		txContext, err = h.isValidTxSim(msg.ChannelId, msg.Txid, "no ledger context")
	}

	chaincodeName := h.chaincodeID.Name + ":" + h.chaincodeID.Version
	h.Metrics.ShimRequestsReceived.With(
		"type", msg.Type.String(),
		"channel", msg.ChannelId,
		"chaincode", chaincodeName,
	).Add(1)

	var resp *pb.ChaincodeMessage
	if err == nil {
		resp, err = delegate(msg, txContext)
	}

	if err != nil {
		err = errors.Wrapf(err, "%s failed: transaction ID: %s", msg.Type, msg.Txid)
		chaincodeLogger.Errorf("[%s] Failed to handle %s. error: %+v", shorttxid(msg.Txid), msg.Type, err)
		resp = &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Payload: []byte(err.Error()), Txid: msg.Txid, ChannelId: msg.ChannelId}
	}

	chaincodeLogger.Debugf("[%s] Completed %s. Sending %s", shorttxid(msg.Txid), msg.Type, resp.Type)
	h.ActiveTransactions.Remove(msg.ChannelId, msg.Txid)
	h.serialSendAsync(resp)
	duration := time.Since(startTime)
	h.Metrics.ShimRequestDuration.With(
		"type", msg.Type.String(),
		"channel", msg.ChannelId,
		"chaincode", chaincodeName,
		"success", strconv.FormatBool(resp.Type != pb.ChaincodeMessage_ERROR),
	).Observe(float64(duration) / float64(time.Second))
	h.Metrics.ShimRequestsCompleted.With(
		"type", msg.Type.String(),
		"channel", msg.ChannelId,
		"chaincode", chaincodeName,
		"success", strconv.FormatBool(resp.Type != pb.ChaincodeMessage_ERROR),
	).Add(1)
}

func shorttxid(txid string) string {
	if len(txid) < 8 {
		return txid
	}
	return txid[0:8]
}



func ParseName(ccName string) *sysccprovider.ChaincodeInstance {
	ci := &sysccprovider.ChaincodeInstance{}

	z := strings.SplitN(ccName, "/", 2)
	if len(z) == 2 {
		ci.ChainID = z[1]
	}
	z = strings.SplitN(z[0], ":", 2)
	if len(z) == 2 {
		ci.ChaincodeVersion = z[1]
	}
	ci.ChaincodeName = z[0]

	return ci
}

func (h *Handler) ChaincodeName() string {
	if h.ccInstance == nil {
		return ""
	}
	return h.ccInstance.ChaincodeName
}


func (h *Handler) serialSend(msg *pb.ChaincodeMessage) error {
	h.serialLock.Lock()
	defer h.serialLock.Unlock()

	if err := h.chatStream.Send(msg); err != nil {
		err = errors.WithMessage(err, fmt.Sprintf("[%s] error sending %s", shorttxid(msg.Txid), msg.Type))
		chaincodeLogger.Errorf("%+v", err)
		return err
	}

	return nil
}






func (h *Handler) serialSendAsync(msg *pb.ChaincodeMessage) {
	go func() {
		if err := h.serialSend(msg); err != nil {
			
			resp := &pb.ChaincodeMessage{
				Type:      pb.ChaincodeMessage_ERROR,
				Payload:   []byte(err.Error()),
				Txid:      msg.Txid,
				ChannelId: msg.ChannelId,
			}
			h.Notify(resp)

			
			h.errChan <- err
		}
	}()
}


func (h *Handler) checkACL(signedProp *pb.SignedProposal, proposal *pb.Proposal, ccIns *sysccprovider.ChaincodeInstance) error {
	
	
	if h.SystemCCProvider.IsSysCCAndNotInvokableCC2CC(ccIns.ChaincodeName) {
		return errors.Errorf("system chaincode %s cannot be invoked with a cc2cc invocation", ccIns.ChaincodeName)
	}

	
	
	
	
	

	if h.SystemCCProvider.IsSysCC(ccIns.ChaincodeName) {
		
		return nil
	}

	
	if signedProp == nil {
		return errors.Errorf("signed proposal must not be nil from caller [%s]", ccIns.String())
	}

	return h.ACLProvider.CheckACL(resources.Peer_ChaincodeToChaincode, ccIns.ChainID, signedProp)
}

func (h *Handler) deregister() {
	if h.chaincodeID != nil {
		h.Registry.Deregister(h.chaincodeID.Name)
	}
}

func (h *Handler) ProcessStream(stream ccintf.ChaincodeStream) error {
	defer h.deregister()

	h.chatStream = stream
	h.errChan = make(chan error, 1)

	var keepaliveCh <-chan time.Time
	if h.Keepalive != 0 {
		ticker := time.NewTicker(h.Keepalive)
		defer ticker.Stop()
		keepaliveCh = ticker.C
	}

	
	type recvMsg struct {
		msg *pb.ChaincodeMessage
		err error
	}
	msgAvail := make(chan *recvMsg, 1)

	receiveMessage := func() {
		in, err := h.chatStream.Recv()
		msgAvail <- &recvMsg{in, err}
	}

	go receiveMessage()
	for {
		select {
		case rmsg := <-msgAvail:
			switch {
			
			case rmsg.err == io.EOF:
				chaincodeLogger.Debugf("received EOF, ending chaincode support stream: %s", rmsg.err)
				return rmsg.err
			case rmsg.err != nil:
				err := errors.Wrap(rmsg.err, "receive failed")
				chaincodeLogger.Errorf("handling chaincode support stream: %+v", err)
				return err
			case rmsg.msg == nil:
				err := errors.New("received nil message, ending chaincode support stream")
				chaincodeLogger.Debugf("%+v", err)
				return err
			default:
				err := h.handleMessage(rmsg.msg)
				if err != nil {
					err = errors.WithMessage(err, "error handling message, ending stream")
					chaincodeLogger.Errorf("[%s] %+v", shorttxid(rmsg.msg.Txid), err)
					return err
				}

				go receiveMessage()
			}

		case sendErr := <-h.errChan:
			err := errors.Wrapf(sendErr, "received error while sending message, ending chaincode support stream")
			chaincodeLogger.Errorf("%s", err)
			return err
		case <-keepaliveCh:
			
			
			h.serialSendAsync(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_KEEPALIVE})
			continue
		}
	}
}


func (h *Handler) sendReady() error {
	chaincodeLogger.Debugf("sending READY for chaincode %+v", h.chaincodeID)
	ccMsg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_READY}

	
	if err := h.serialSend(ccMsg); err != nil {
		chaincodeLogger.Errorf("error sending READY (%s) for chaincode %+v", err, h.chaincodeID)
		return err
	}

	h.state = Ready

	chaincodeLogger.Debugf("Changed to state ready for chaincode %+v", h.chaincodeID)

	return nil
}



func (h *Handler) notifyRegistry(err error) {
	if err == nil {
		err = h.sendReady()
	}

	if err != nil {
		h.Registry.Failed(h.chaincodeID.Name, err)
		chaincodeLogger.Errorf("failed to start %s", h.chaincodeID)
		return
	}

	h.Registry.Ready(h.chaincodeID.Name)
}


func (h *Handler) HandleRegister(msg *pb.ChaincodeMessage) {
	chaincodeLogger.Debugf("Received %s in state %s", msg.Type, h.state)
	chaincodeID := &pb.ChaincodeID{}
	err := proto.Unmarshal(msg.Payload, chaincodeID)
	if err != nil {
		chaincodeLogger.Errorf("Error in received %s, could NOT unmarshal registration info: %s", pb.ChaincodeMessage_REGISTER, err)
		return
	}

	
	h.chaincodeID = chaincodeID
	err = h.Registry.Register(h)
	if err != nil {
		h.notifyRegistry(err)
		return
	}

	
	
	h.ccInstance = ParseName(h.chaincodeID.Name)

	chaincodeLogger.Debugf("Got %s for chaincodeID = %s, sending back %s", pb.ChaincodeMessage_REGISTER, chaincodeID, pb.ChaincodeMessage_REGISTERED)
	if err := h.serialSend(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTERED}); err != nil {
		chaincodeLogger.Errorf("error sending %s: %s", pb.ChaincodeMessage_REGISTERED, err)
		h.notifyRegistry(err)
		return
	}

	h.state = Established

	chaincodeLogger.Debugf("Changed state to established for %+v", h.chaincodeID)

	
	h.notifyRegistry(nil)
}

func (h *Handler) Notify(msg *pb.ChaincodeMessage) {
	tctx := h.TXContexts.Get(msg.ChannelId, msg.Txid)
	if tctx == nil {
		chaincodeLogger.Debugf("notifier Txid:%s, channelID:%s does not exist for handling message %s", msg.Txid, msg.ChannelId, msg.Type)
		return
	}

	chaincodeLogger.Debugf("[%s] notifying Txid:%s, channelID:%s", shorttxid(msg.Txid), msg.Txid, msg.ChannelId)
	tctx.ResponseNotifier <- msg
	tctx.CloseQueryIterators()
}


func (h *Handler) isValidTxSim(channelID string, txid string, fmtStr string, args ...interface{}) (*TransactionContext, error) {
	txContext := h.TXContexts.Get(channelID, txid)
	if txContext == nil || txContext.TXSimulator == nil {
		err := errors.Errorf(fmtStr, args...)
		chaincodeLogger.Errorf("no ledger context: %s %s\n\n %+v", channelID, txid, err)
		return nil, err
	}
	return txContext, nil
}


func (h *Handler) registerTxid(msg *pb.ChaincodeMessage) bool {
	
	if h.ActiveTransactions.Add(msg.ChannelId, msg.Txid) {
		return true
	}

	
	chaincodeName := "unknown"
	if h.chaincodeID != nil {
		chaincodeName = h.chaincodeID.Name
	}
	chaincodeLogger.Errorf("[%s] Another request pending for this CC: %s, Txid: %s, ChannelID: %s. Cannot process.", shorttxid(msg.Txid), chaincodeName, msg.Txid, msg.ChannelId)
	return false
}

func (h *Handler) checkMetadataCap(msg *pb.ChaincodeMessage) error {
	ac, exists := h.AppConfig.GetApplicationConfig(msg.ChannelId)
	if !exists {
		return errors.Errorf("application config does not exist for %s", msg.ChannelId)
	}

	if !ac.Capabilities().KeyLevelEndorsement() {
		return errors.New("key level endorsement is not enabled")
	}
	return nil
}


func (h *Handler) HandleGetState(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	getState := &pb.GetState{}
	err := proto.Unmarshal(msg.Payload, getState)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	chaincodeName := h.ChaincodeName()
	chaincodeLogger.Debugf("[%s] getting state for chaincode %s, key %s, channel %s", shorttxid(msg.Txid), chaincodeName, getState.Key, txContext.ChainID)

	var res []byte
	if isCollectionSet(getState.Collection) {
		res, err = txContext.TXSimulator.GetPrivateData(chaincodeName, getState.Collection, getState.Key)
	} else {
		res, err = txContext.TXSimulator.GetState(chaincodeName, getState.Key)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if res == nil {
		chaincodeLogger.Debugf("[%s] No state associated with key: %s. Sending %s with an empty payload", shorttxid(msg.Txid), getState.Key, pb.ChaincodeMessage_RESPONSE)
	}

	
	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: res, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}


func (h *Handler) HandleGetStateMetadata(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	err := h.checkMetadataCap(msg)
	if err != nil {
		return nil, err
	}

	getStateMetadata := &pb.GetStateMetadata{}
	err = proto.Unmarshal(msg.Payload, getStateMetadata)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	chaincodeName := h.ChaincodeName()
	chaincodeLogger.Debugf("[%s] getting state metadata for chaincode %s, key %s, channel %s", shorttxid(msg.Txid), chaincodeName, getStateMetadata.Key, txContext.ChainID)

	var metadata map[string][]byte
	if isCollectionSet(getStateMetadata.Collection) {
		metadata, err = txContext.TXSimulator.GetPrivateDataMetadata(chaincodeName, getStateMetadata.Collection, getStateMetadata.Key)
	} else {
		metadata, err = txContext.TXSimulator.GetStateMetadata(chaincodeName, getStateMetadata.Key)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var metadataResult pb.StateMetadataResult
	for metakey := range metadata {
		md := &pb.StateMetadata{Metakey: metakey, Value: metadata[metakey]}
		metadataResult.Entries = append(metadataResult.Entries, md)
	}
	res, err := proto.Marshal(&metadataResult)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	
	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: res, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}


func (h *Handler) HandleGetStateByRange(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	getStateByRange := &pb.GetStateByRange{}
	err := proto.Unmarshal(msg.Payload, getStateByRange)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	metadata, err := getQueryMetadataFromBytes(getStateByRange.Metadata)
	if err != nil {
		return nil, err
	}

	totalReturnLimit := calculateTotalReturnLimit(metadata)

	iterID := h.UUIDGenerator.New()
	chaincodeName := h.ChaincodeName()

	var rangeIter commonledger.ResultsIterator
	var paginationInfo map[string]interface{}

	isPaginated := false

	if isCollectionSet(getStateByRange.Collection) {
		rangeIter, err = txContext.TXSimulator.GetPrivateDataRangeScanIterator(chaincodeName, getStateByRange.Collection,
			getStateByRange.StartKey, getStateByRange.EndKey)
	} else if isMetadataSetForPagination(metadata) {
		paginationInfo, err = createPaginationInfoFromMetadata(metadata, totalReturnLimit, pb.ChaincodeMessage_GET_STATE_BY_RANGE)
		if err != nil {
			return nil, err
		}
		isPaginated = true

		startKey := getStateByRange.StartKey

		if isMetadataSetForPagination(metadata) {
			if metadata.Bookmark != "" {
				startKey = metadata.Bookmark
			}
		}
		rangeIter, err = txContext.TXSimulator.GetStateRangeScanIteratorWithMetadata(chaincodeName,
			startKey, getStateByRange.EndKey, paginationInfo)
	} else {
		rangeIter, err = txContext.TXSimulator.GetStateRangeScanIterator(chaincodeName, getStateByRange.StartKey, getStateByRange.EndKey)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	txContext.InitializeQueryContext(iterID, rangeIter)

	payload, err := h.QueryResponseBuilder.BuildQueryResponse(txContext, rangeIter, iterID, isPaginated, totalReturnLimit)
	if err != nil {
		txContext.CleanupQueryContext(iterID)
		return nil, errors.WithStack(err)
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		txContext.CleanupQueryContext(iterID)
		return nil, errors.Wrap(err, "marshal failed")
	}

	chaincodeLogger.Debugf("Got keys and values. Sending %s", pb.ChaincodeMessage_RESPONSE)
	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: payloadBytes, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}


func (h *Handler) HandleQueryStateNext(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	queryStateNext := &pb.QueryStateNext{}
	err := proto.Unmarshal(msg.Payload, queryStateNext)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	queryIter := txContext.GetQueryIterator(queryStateNext.Id)
	if queryIter == nil {
		return nil, errors.New("query iterator not found")
	}

	totalReturnLimit := calculateTotalReturnLimit(nil)

	payload, err := h.QueryResponseBuilder.BuildQueryResponse(txContext, queryIter, queryStateNext.Id, false, totalReturnLimit)
	if err != nil {
		txContext.CleanupQueryContext(queryStateNext.Id)
		return nil, errors.WithStack(err)
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		txContext.CleanupQueryContext(queryStateNext.Id)
		return nil, errors.Wrap(err, "marshal failed")
	}

	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: payloadBytes, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}


func (h *Handler) HandleQueryStateClose(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	queryStateClose := &pb.QueryStateClose{}
	err := proto.Unmarshal(msg.Payload, queryStateClose)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	iter := txContext.GetQueryIterator(queryStateClose.Id)
	if iter != nil {
		txContext.CleanupQueryContext(queryStateClose.Id)
	}

	payload := &pb.QueryResponse{HasMore: false, Id: queryStateClose.Id}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "marshal failed")
	}

	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: payloadBytes, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}


func (h *Handler) HandleGetQueryResult(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	iterID := h.UUIDGenerator.New()
	chaincodeName := h.ChaincodeName()

	getQueryResult := &pb.GetQueryResult{}
	err := proto.Unmarshal(msg.Payload, getQueryResult)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	metadata, err := getQueryMetadataFromBytes(getQueryResult.Metadata)
	if err != nil {
		return nil, err
	}

	totalReturnLimit := calculateTotalReturnLimit(metadata)
	isPaginated := false

	var executeIter commonledger.ResultsIterator
	var paginationInfo map[string]interface{}

	if isCollectionSet(getQueryResult.Collection) {
		executeIter, err = txContext.TXSimulator.ExecuteQueryOnPrivateData(chaincodeName, getQueryResult.Collection, getQueryResult.Query)
	} else if isMetadataSetForPagination(metadata) {
		paginationInfo, err = createPaginationInfoFromMetadata(metadata, totalReturnLimit, pb.ChaincodeMessage_GET_QUERY_RESULT)
		if err != nil {
			return nil, err
		}
		isPaginated = true
		executeIter, err = txContext.TXSimulator.ExecuteQueryWithMetadata(chaincodeName,
			getQueryResult.Query, paginationInfo)

	} else {
		executeIter, err = txContext.TXSimulator.ExecuteQuery(chaincodeName, getQueryResult.Query)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	txContext.InitializeQueryContext(iterID, executeIter)

	payload, err := h.QueryResponseBuilder.BuildQueryResponse(txContext, executeIter, iterID, isPaginated, totalReturnLimit)
	if err != nil {
		txContext.CleanupQueryContext(iterID)
		return nil, errors.WithStack(err)
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		txContext.CleanupQueryContext(iterID)
		return nil, errors.Wrap(err, "marshal failed")
	}

	chaincodeLogger.Debugf("Got keys and values. Sending %s", pb.ChaincodeMessage_RESPONSE)
	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: payloadBytes, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}


func (h *Handler) HandleGetHistoryForKey(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	iterID := h.UUIDGenerator.New()
	chaincodeName := h.ChaincodeName()

	getHistoryForKey := &pb.GetHistoryForKey{}
	err := proto.Unmarshal(msg.Payload, getHistoryForKey)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	historyIter, err := txContext.HistoryQueryExecutor.GetHistoryForKey(chaincodeName, getHistoryForKey.Key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	totalReturnLimit := calculateTotalReturnLimit(nil)

	txContext.InitializeQueryContext(iterID, historyIter)
	payload, err := h.QueryResponseBuilder.BuildQueryResponse(txContext, historyIter, iterID, false, totalReturnLimit)
	if err != nil {
		txContext.CleanupQueryContext(iterID)
		return nil, errors.WithStack(err)
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		txContext.CleanupQueryContext(iterID)
		return nil, errors.Wrap(err, "marshal failed")
	}

	chaincodeLogger.Debugf("Got keys and values. Sending %s", pb.ChaincodeMessage_RESPONSE)
	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: payloadBytes, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}

func isCollectionSet(collection string) bool {
	return collection != ""
}

func isMetadataSetForPagination(metadata *pb.QueryMetadata) bool {
	if metadata == nil {
		return false
	}

	if metadata.PageSize == 0 && metadata.Bookmark == "" {
		return false
	}

	return true
}

func getQueryMetadataFromBytes(metadataBytes []byte) (*pb.QueryMetadata, error) {
	if metadataBytes != nil {
		metadata := &pb.QueryMetadata{}
		err := proto.Unmarshal(metadataBytes, metadata)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal failed")
		}
		return metadata, nil
	}
	return nil, nil
}

func createPaginationInfoFromMetadata(metadata *pb.QueryMetadata, totalReturnLimit int32, queryType pb.ChaincodeMessage_Type) (map[string]interface{}, error) {
	paginationInfoMap := make(map[string]interface{})

	switch queryType {
	case pb.ChaincodeMessage_GET_QUERY_RESULT:
		paginationInfoMap["bookmark"] = metadata.Bookmark
	case pb.ChaincodeMessage_GET_STATE_BY_RANGE:
		
	default:
		return nil, errors.New("query type must be either GetQueryResult or GetStateByRange")
	}

	paginationInfoMap["limit"] = totalReturnLimit
	return paginationInfoMap, nil
}

func calculateTotalReturnLimit(metadata *pb.QueryMetadata) int32 {
	totalReturnLimit := int32(ledgerconfig.GetTotalQueryLimit())
	if metadata != nil {
		pageSize := int32(metadata.PageSize)
		if pageSize > 0 && pageSize < totalReturnLimit {
			totalReturnLimit = pageSize
		}
	}
	return totalReturnLimit
}

func (h *Handler) getTxContextForInvoke(channelID string, txid string, payload []byte, format string, args ...interface{}) (*TransactionContext, error) {
	
	if channelID != "" {
		return h.isValidTxSim(channelID, txid, "could not get valid transaction")
	}

	chaincodeSpec := &pb.ChaincodeSpec{}
	err := proto.Unmarshal(payload, chaincodeSpec)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	
	
	
	targetInstance := ParseName(chaincodeSpec.ChaincodeId.Name)

	
	
	
	if !h.SystemCCProvider.IsSysCC(targetInstance.ChaincodeName) {
		
		return h.isValidTxSim("", txid, "could not get valid transaction")
	}

	
	
	txContext := h.TXContexts.Get(channelID, txid)
	if txContext == nil {
		return nil, errors.New("failed to get transaction context")
	}

	return txContext, nil
}

func (h *Handler) HandlePutState(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	putState := &pb.PutState{}
	err := proto.Unmarshal(msg.Payload, putState)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	chaincodeName := h.ChaincodeName()
	if isCollectionSet(putState.Collection) {
		err = txContext.TXSimulator.SetPrivateData(chaincodeName, putState.Collection, putState.Key, putState.Value)
	} else {
		err = txContext.TXSimulator.SetState(chaincodeName, putState.Key, putState.Value)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}

func (h *Handler) HandlePutStateMetadata(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	err := h.checkMetadataCap(msg)
	if err != nil {
		return nil, err
	}

	putStateMetadata := &pb.PutStateMetadata{}
	err = proto.Unmarshal(msg.Payload, putStateMetadata)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	metadata := make(map[string][]byte)
	metadata[putStateMetadata.Metadata.Metakey] = putStateMetadata.Metadata.Value

	chaincodeName := h.ChaincodeName()
	if isCollectionSet(putStateMetadata.Collection) {
		err = txContext.TXSimulator.SetPrivateDataMetadata(chaincodeName, putStateMetadata.Collection, putStateMetadata.Key, metadata)
	} else {
		err = txContext.TXSimulator.SetStateMetadata(chaincodeName, putStateMetadata.Key, metadata)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}

func (h *Handler) HandleDelState(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	delState := &pb.DelState{}
	err := proto.Unmarshal(msg.Payload, delState)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	chaincodeName := h.ChaincodeName()
	if isCollectionSet(delState.Collection) {
		err = txContext.TXSimulator.DeletePrivateData(chaincodeName, delState.Collection, delState.Key)
	} else {
		err = txContext.TXSimulator.DeleteState(chaincodeName, delState.Key)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	
	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}


func (h *Handler) HandleInvokeChaincode(msg *pb.ChaincodeMessage, txContext *TransactionContext) (*pb.ChaincodeMessage, error) {
	chaincodeLogger.Debugf("[%s] C-call-C", shorttxid(msg.Txid))

	chaincodeSpec := &pb.ChaincodeSpec{}
	err := proto.Unmarshal(msg.Payload, chaincodeSpec)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	
	
	
	targetInstance := ParseName(chaincodeSpec.ChaincodeId.Name)
	chaincodeSpec.ChaincodeId.Name = targetInstance.ChaincodeName
	if targetInstance.ChainID == "" {
		
		targetInstance.ChainID = txContext.ChainID
	}
	chaincodeLogger.Debugf("[%s] C-call-C %s on channel %s", shorttxid(msg.Txid), targetInstance.ChaincodeName, targetInstance.ChainID)

	err = h.checkACL(txContext.SignedProp, txContext.Proposal, targetInstance)
	if err != nil {
		chaincodeLogger.Errorf(
			"[%s] C-call-C %s on channel %s failed check ACL [%v]: [%s]",
			shorttxid(msg.Txid),
			targetInstance.ChaincodeName,
			targetInstance.ChainID,
			txContext.SignedProp,
			err,
		)
		return nil, errors.WithStack(err)
	}

	
	
	txParams := &ccprovider.TransactionParams{
		TxID:                 msg.Txid,
		ChannelID:            targetInstance.ChainID,
		SignedProp:           txContext.SignedProp,
		Proposal:             txContext.Proposal,
		TXSimulator:          txContext.TXSimulator,
		HistoryQueryExecutor: txContext.HistoryQueryExecutor,
	}

	if targetInstance.ChainID != txContext.ChainID {
		lgr := h.LedgerGetter.GetLedger(targetInstance.ChainID)
		if lgr == nil {
			return nil, errors.Errorf("failed to find ledger for channel: %s", targetInstance.ChainID)
		}

		sim, err := lgr.NewTxSimulator(msg.Txid)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		defer sim.Done()

		hqe, err := lgr.NewHistoryQueryExecutor()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		txParams.TXSimulator = sim
		txParams.HistoryQueryExecutor = hqe
	}

	chaincodeLogger.Debugf("[%s] getting chaincode data for %s on channel %s", shorttxid(msg.Txid), targetInstance.ChaincodeName, targetInstance.ChainID)

	version := h.SystemCCVersion
	if !h.SystemCCProvider.IsSysCC(targetInstance.ChaincodeName) {
		
		cd, err := h.DefinitionGetter.ChaincodeDefinition(targetInstance.ChaincodeName, txParams.TXSimulator)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		version = cd.CCVersion()

		err = h.InstantiationPolicyChecker.CheckInstantiationPolicy(targetInstance.ChaincodeName, version, cd.(*ccprovider.ChaincodeData))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	
	chaincodeLogger.Debugf("[%s] launching chaincode %s on channel %s", shorttxid(msg.Txid), targetInstance.ChaincodeName, targetInstance.ChainID)

	cccid := &ccprovider.CCContext{
		Name:    targetInstance.ChaincodeName,
		Version: version,
	}

	
	responseMessage, err := h.Invoker.Invoke(txParams, cccid, chaincodeSpec.Input)
	if err != nil {
		return nil, errors.Wrap(err, "execute failed")
	}

	
	
	res, err := proto.Marshal(responseMessage)
	if err != nil {
		return nil, errors.Wrap(err, "marshal failed")
	}

	return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: res, Txid: msg.Txid, ChannelId: msg.ChannelId}, nil
}

func (h *Handler) Execute(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, msg *pb.ChaincodeMessage, timeout time.Duration) (*pb.ChaincodeMessage, error) {
	chaincodeLogger.Debugf("Entry")
	defer chaincodeLogger.Debugf("Exit")

	txctx, err := h.TXContexts.Create(txParams)
	if err != nil {
		return nil, err
	}
	defer h.TXContexts.Delete(msg.ChannelId, msg.Txid)

	if err := h.setChaincodeProposal(txParams.SignedProp, txParams.Proposal, msg); err != nil {
		return nil, err
	}

	h.serialSendAsync(msg)

	var ccresp *pb.ChaincodeMessage
	select {
	case ccresp = <-txctx.ResponseNotifier:
		
		
	case <-time.After(timeout):
		err = errors.New("timeout expired while executing transaction")
	}

	return ccresp, err
}

func (h *Handler) setChaincodeProposal(signedProp *pb.SignedProposal, prop *pb.Proposal, msg *pb.ChaincodeMessage) error {
	if prop != nil && signedProp == nil {
		return errors.New("failed getting proposal context. Signed proposal is nil")
	}
	
	
	if prop != nil {
		msg.Proposal = signedProp
	}
	return nil
}

func (h *Handler) State() State { return h.state }
func (h *Handler) Close()       { h.TXContexts.Close() }

type State int

const (
	Created State = iota
	Established
	Ready
)

func (s State) String() string {
	switch s {
	case Created:
		return "created"
	case Established:
		return "established"
	case Ready:
		return "ready"
	default:
		return "UNKNOWN"
	}
}
