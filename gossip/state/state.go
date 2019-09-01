/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package state

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	proto "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/mcc-github/blockchain-protos-go/transientstore"
	vsccErrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	common2 "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)




type GossipStateProvider interface {
	AddPayload(payload *proto.Payload) error

	
	Stop()
}

const (
	defAntiEntropyInterval             = 10 * time.Second
	defAntiEntropyStateResponseTimeout = 3 * time.Second
	defAntiEntropyBatchSize            = 10

	defChannelBufferSize     = 100
	defAntiEntropyMaxRetries = 3

	defMaxBlockDistance = 100

	blocking    = true
	nonBlocking = false

	enqueueRetryInterval = time.Millisecond * 100
)


type Configuration struct {
	AntiEntropyInterval             time.Duration
	AntiEntropyStateResponseTimeout time.Duration
	AntiEntropyBatchSize            uint64
	MaxBlockDistance                int
	AntiEntropyMaxRetries           int
	ChannelBufferSize               int
	EnableStateTransfer             bool
}


type GossipAdapter interface {
	
	Send(msg *proto.GossipMessage, peers ...*comm.RemotePeer)

	
	
	
	
	Accept(acceptor common2.MessageAcceptor, passThrough bool) (<-chan *proto.GossipMessage, <-chan protoext.ReceivedMessage)

	
	
	UpdateLedgerHeight(height uint64, channelID common2.ChannelID)

	
	
	PeersOfChannel(common2.ChannelID) []discovery.NetworkMember
}



type MCSAdapter interface {
	
	
	
	VerifyBlock(channelID common2.ChannelID, seqNum uint64, signedBlock *common.Block) error

	
	
	
	
	VerifyByChannel(channelID common2.ChannelID, peerIdentity api.PeerIdentityType, signature, message []byte) error
}


type ledgerResources interface {
	
	
	StoreBlock(block *common.Block, data util.PvtDataCollections) error

	
	StorePvtData(txid string, privData *transientstore.TxPvtReadWriteSetWithConfigInfo, blckHeight uint64) error

	
	
	
	
	GetPvtDataAndBlockByNum(seqNum uint64, peerAuthInfo protoutil.SignedData) (*common.Block, util.PvtDataCollections, error)

	
	LedgerHeight() (uint64, error)

	
	Close()
}



type ServicesMediator struct {
	GossipAdapter
	MCSAdapter
}




type GossipStateProviderImpl struct {
	
	chainID string

	mediator *ServicesMediator

	
	payloads PayloadsBuffer

	ledger ledgerResources

	stateResponseCh chan protoext.ReceivedMessage

	stateRequestCh chan protoext.ReceivedMessage

	stopCh chan struct{}

	once sync.Once

	stateTransferActive int32

	stateMetrics *metrics.StateMetrics

	requestValidator *stateRequestValidator

	blockingMode bool

	config *StateConfig
}

var logger = util.GetLogger(util.StateLogger, "")


type stateRequestValidator struct {
}


func (v *stateRequestValidator) validate(request *proto.RemoteStateRequest, batchSize uint64) error {
	if request.StartSeqNum > request.EndSeqNum {
		return errors.Errorf("Invalid sequence interval [%d...%d).", request.StartSeqNum, request.EndSeqNum)
	}

	if request.EndSeqNum > batchSize+request.StartSeqNum {
		return errors.Errorf("Requesting blocks range [%d-%d) greater than configured allowed"+
			" (%d) batching size for anti-entropy.", request.StartSeqNum, request.EndSeqNum, batchSize)
	}
	return nil
}



func NewGossipStateProvider(
	chainID string,
	services *ServicesMediator,
	ledger ledgerResources,
	stateMetrics *metrics.StateMetrics,
	blockingMode bool,
	config *StateConfig,
) GossipStateProvider {

	gossipChan, _ := services.Accept(func(message interface{}) bool {
		
		return protoext.IsDataMsg(message.(*proto.GossipMessage)) &&
			bytes.Equal(message.(*proto.GossipMessage).Channel, []byte(chainID))
	}, false)

	remoteStateMsgFilter := func(message interface{}) bool {
		receivedMsg := message.(protoext.ReceivedMessage)
		msg := receivedMsg.GetGossipMessage()
		if !(protoext.IsRemoteStateMessage(msg.GossipMessage) || msg.GetPrivateData() != nil) {
			return false
		}
		
		if !bytes.Equal(msg.Channel, []byte(chainID)) {
			return false
		}
		connInfo := receivedMsg.GetConnectionInfo()
		authErr := services.VerifyByChannel(msg.Channel, connInfo.Identity, connInfo.Auth.Signature, connInfo.Auth.SignedData)
		if authErr != nil {
			logger.Warning("Got unauthorized request from", string(connInfo.Identity))
			return false
		}
		return true
	}

	
	_, commChan := services.Accept(remoteStateMsgFilter, true)

	height, err := ledger.LedgerHeight()
	if height == 0 {
		
		
		logger.Panic("Committer height cannot be zero, ledger should include at least one block (genesis).")
	}

	if err != nil {
		logger.Error("Could not read ledger info to obtain current ledger height due to: ", errors.WithStack(err))
		
		
		return nil
	}

	s := &GossipStateProviderImpl{
		
		mediator: services,
		
		chainID: chainID,
		
		payloads: &metricsBuffer{
			PayloadsBuffer: NewPayloadsBuffer(height),
			sizeMetrics:    stateMetrics.PayloadBufferSize,
			chainID:        chainID,
		},
		ledger:              ledger,
		stateResponseCh:     make(chan protoext.ReceivedMessage, config.StateChannelSize),
		stateRequestCh:      make(chan protoext.ReceivedMessage, config.StateChannelSize),
		stopCh:              make(chan struct{}),
		stateTransferActive: 0,
		once:                sync.Once{},
		stateMetrics:        stateMetrics,
		requestValidator:    &stateRequestValidator{},
		blockingMode:        blockingMode,
		config:              config,
	}

	logger.Infof("Updating metadata information, "+
		"current ledger sequence is at = %d, next expected block is = %d", height-1, s.payloads.Next())
	logger.Debug("Updating gossip ledger height to", height)
	services.UpdateLedgerHeight(height, common2.ChannelID(s.chainID))

	
	go s.receiveAndQueueGossipMessages(gossipChan)
	go s.receiveAndDispatchDirectMessages(commChan)
	
	go s.deliverPayloads()
	if s.config.StateEnabled {
		
		go s.antiEntropy()
	}
	
	go s.processStateRequests()

	return s
}

func (s *GossipStateProviderImpl) receiveAndQueueGossipMessages(ch <-chan *proto.GossipMessage) {
	for msg := range ch {
		logger.Debug("Received new message via gossip channel")
		go func(msg *proto.GossipMessage) {
			if !bytes.Equal(msg.Channel, []byte(s.chainID)) {
				logger.Warning("Received enqueue for channel",
					string(msg.Channel), "while expecting channel", s.chainID, "ignoring enqueue")
				return
			}

			dataMsg := msg.GetDataMsg()
			if dataMsg != nil {
				if err := s.addPayload(dataMsg.GetPayload(), nonBlocking); err != nil {
					logger.Warningf("Block [%d] received from gossip wasn't added to payload buffer: %v", dataMsg.Payload.SeqNum, err)
					return
				}

			} else {
				logger.Debug("Gossip message received is not of data message type, usually this should not happen.")
			}
		}(msg)
	}
}

func (s *GossipStateProviderImpl) receiveAndDispatchDirectMessages(ch <-chan protoext.ReceivedMessage) {
	for msg := range ch {
		logger.Debug("Dispatching a message", msg)
		go func(msg protoext.ReceivedMessage) {
			gm := msg.GetGossipMessage()
			
			if protoext.IsRemoteStateMessage(gm.GossipMessage) {
				logger.Debug("Handling direct state transfer message")
				
				s.directMessage(msg)
			} else if gm.GetPrivateData() != nil {
				logger.Debug("Handling private data collection message")
				
				s.privateDataMessage(msg)
			}
		}(msg)
	}
}

func (s *GossipStateProviderImpl) privateDataMessage(msg protoext.ReceivedMessage) {
	if !bytes.Equal(msg.GetGossipMessage().Channel, []byte(s.chainID)) {
		logger.Warning("Received state transfer request for channel",
			string(msg.GetGossipMessage().Channel), "while expecting channel", s.chainID, "skipping request...")
		return
	}

	gossipMsg := msg.GetGossipMessage()
	pvtDataMsg := gossipMsg.GetPrivateData()

	if pvtDataMsg.Payload == nil {
		logger.Warning("Malformed private data message, no payload provided")
		return
	}

	collectionName := pvtDataMsg.Payload.CollectionName
	txID := pvtDataMsg.Payload.TxId
	pvtRwSet := pvtDataMsg.Payload.PrivateRwset

	if len(pvtRwSet) == 0 {
		logger.Warning("Malformed private data message, no rwset provided, collection name = ", collectionName)
		return
	}

	txPvtRwSet := &rwset.TxPvtReadWriteSet{
		DataModel: rwset.TxReadWriteSet_KV,
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{{
			Namespace: pvtDataMsg.Payload.Namespace,
			CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{{
				CollectionName: collectionName,
				Rwset:          pvtRwSet,
			}}},
		},
	}

	txPvtRwSetWithConfig := &transientstore.TxPvtReadWriteSetWithConfigInfo{
		PvtRwset: txPvtRwSet,
		CollectionConfigs: map[string]*common.CollectionConfigPackage{
			pvtDataMsg.Payload.Namespace: pvtDataMsg.Payload.CollectionConfigs,
		},
	}

	if err := s.ledger.StorePvtData(txID, txPvtRwSetWithConfig, pvtDataMsg.Payload.PrivateSimHeight); err != nil {
		logger.Errorf("Wasn't able to persist private data for collection %s, due to %s", collectionName, err)
		msg.Ack(err) 
	}

	msg.Ack(nil)
	logger.Debug("Private data for collection", collectionName, "has been stored")
}

func (s *GossipStateProviderImpl) directMessage(msg protoext.ReceivedMessage) {
	logger.Debug("[ENTER] -> directMessage")
	defer logger.Debug("[EXIT] ->  directMessage")

	if msg == nil {
		logger.Error("Got nil message via end-to-end channel, should not happen!")
		return
	}

	if !bytes.Equal(msg.GetGossipMessage().Channel, []byte(s.chainID)) {
		logger.Warning("Received state transfer request for channel",
			string(msg.GetGossipMessage().Channel), "while expecting channel", s.chainID, "skipping request...")
		return
	}

	incoming := msg.GetGossipMessage()

	if incoming.GetStateRequest() != nil {
		if len(s.stateRequestCh) < s.config.StateChannelSize {
			
			
			s.stateRequestCh <- msg
		}
	} else if incoming.GetStateResponse() != nil {
		
		
		if atomic.LoadInt32(&s.stateTransferActive) == 1 {
			
			s.stateResponseCh <- msg
		}
	}
}

func (s *GossipStateProviderImpl) processStateRequests() {
	for {
		msg, stillOpen := <-s.stateRequestCh
		if !stillOpen {
			return
		}
		s.handleStateRequest(msg)
	}
}



func (s *GossipStateProviderImpl) handleStateRequest(msg protoext.ReceivedMessage) {
	if msg == nil {
		return
	}
	request := msg.GetGossipMessage().GetStateRequest()

	if err := s.requestValidator.validate(request, s.config.StateBatchSize); err != nil {
		logger.Errorf("State request validation failed, %s. Ignoring request...", err)
		return
	}

	currentHeight, err := s.ledger.LedgerHeight()
	if err != nil {
		logger.Errorf("Cannot access to current ledger height, due to %+v", err)
		return
	}
	if currentHeight < request.EndSeqNum {
		logger.Warningf("Received state request to transfer blocks with sequence numbers higher  [%d...%d] "+
			"than available in ledger (%d)", request.StartSeqNum, request.StartSeqNum, currentHeight)
	}

	endSeqNum := min(currentHeight, request.EndSeqNum)

	response := &proto.RemoteStateResponse{Payloads: make([]*proto.Payload, 0)}
	for seqNum := request.StartSeqNum; seqNum <= endSeqNum; seqNum++ {
		logger.Debug("Reading block ", seqNum, " with private data from the coordinator service")
		connInfo := msg.GetConnectionInfo()
		peerAuthInfo := protoutil.SignedData{
			Data:      connInfo.Auth.SignedData,
			Signature: connInfo.Auth.Signature,
			Identity:  connInfo.Identity,
		}
		block, pvtData, err := s.ledger.GetPvtDataAndBlockByNum(seqNum, peerAuthInfo)

		if err != nil {
			logger.Errorf("cannot read block number %d from ledger, because %+v, skipping...", seqNum, err)
			continue
		}

		if block == nil {
			logger.Errorf("Wasn't able to read block with sequence number %d from ledger, skipping....", seqNum)
			continue
		}

		blockBytes, err := pb.Marshal(block)

		if err != nil {
			logger.Errorf("Could not marshal block: %+v", errors.WithStack(err))
			continue
		}

		var pvtBytes [][]byte
		if pvtData != nil {
			
			pvtBytes, err = pvtData.Marshal()
			if err != nil {
				logger.Errorf("Failed to marshal private rwset for block %d due to %+v", seqNum, errors.WithStack(err))
				continue
			}
		}

		
		response.Payloads = append(response.Payloads, &proto.Payload{
			SeqNum:      seqNum,
			Data:        blockBytes,
			PrivateData: pvtBytes,
		})
	}
	
	msg.Respond(&proto.GossipMessage{
		
		Nonce:   msg.GetGossipMessage().Nonce,
		Tag:     proto.GossipMessage_CHAN_OR_ORG,
		Channel: []byte(s.chainID),
		Content: &proto.GossipMessage_StateResponse{StateResponse: response},
	})
}

func (s *GossipStateProviderImpl) handleStateResponse(msg protoext.ReceivedMessage) (uint64, error) {
	max := uint64(0)
	
	response := msg.GetGossipMessage().GetStateResponse()
	
	if len(response.GetPayloads()) == 0 {
		return uint64(0), errors.New("Received state transfer response without payload")
	}
	for _, payload := range response.GetPayloads() {
		logger.Debugf("Received payload with sequence number %d.", payload.SeqNum)
		block, err := protoutil.UnmarshalBlock(payload.Data)
		if err != nil {
			logger.Warningf("Error unmarshaling payload to block for sequence number %d, due to %+v", payload.SeqNum, err)
			return uint64(0), err
		}

		if err := s.mediator.VerifyBlock(common2.ChannelID(s.chainID), payload.SeqNum, block); err != nil {
			err = errors.WithStack(err)
			logger.Warningf("Error verifying block with sequence number %d, due to %+v", payload.SeqNum, err)
			return uint64(0), err
		}
		if max < payload.SeqNum {
			max = payload.SeqNum
		}

		err = s.addPayload(payload, blocking)
		if err != nil {
			logger.Warningf("Block [%d] received from block transfer wasn't added to payload buffer: %v", payload.SeqNum, err)
		}
	}
	return max, nil
}


func (s *GossipStateProviderImpl) Stop() {
	
	
	s.once.Do(func() {
		close(s.stopCh)
		
		s.ledger.Close()
		close(s.stateRequestCh)
		close(s.stateResponseCh)
	})
}

func (s *GossipStateProviderImpl) deliverPayloads() {
	for {
		select {
		
		case <-s.payloads.Ready():
			logger.Debugf("[%s] Ready to transfer payloads (blocks) to the ledger, next block number is = [%d]", s.chainID, s.payloads.Next())
			
			for payload := s.payloads.Pop(); payload != nil; payload = s.payloads.Pop() {
				rawBlock := &common.Block{}
				if err := pb.Unmarshal(payload.Data, rawBlock); err != nil {
					logger.Errorf("Error getting block with seqNum = %d due to (%+v)...dropping block", payload.SeqNum, errors.WithStack(err))
					continue
				}
				if rawBlock.Data == nil || rawBlock.Header == nil {
					logger.Errorf("Block with claimed sequence %d has no header (%v) or data (%v)",
						payload.SeqNum, rawBlock.Header, rawBlock.Data)
					continue
				}
				logger.Debugf("[%s] Transferring block [%d] with %d transaction(s) to the ledger", s.chainID, payload.SeqNum, len(rawBlock.Data.Data))

				
				var p util.PvtDataCollections
				if payload.PrivateData != nil {
					err := p.Unmarshal(payload.PrivateData)
					if err != nil {
						logger.Errorf("Wasn't able to unmarshal private data for block seqNum = %d due to (%+v)...dropping block", payload.SeqNum, errors.WithStack(err))
						continue
					}
				}
				if err := s.commitBlock(rawBlock, p); err != nil {
					if executionErr, isExecutionErr := err.(*vsccErrors.VSCCExecutionFailureError); isExecutionErr {
						logger.Errorf("Failed executing VSCC due to %v. Aborting chain processing", executionErr)
						return
					}
					logger.Panicf("Cannot commit block to the ledger due to %+v", errors.WithStack(err))
				}
			}
		case <-s.stopCh:
			logger.Debug("State provider has been stopped, finishing to push new blocks.")
			return
		}
	}
}

func (s *GossipStateProviderImpl) antiEntropy() {
	defer logger.Debug("State Provider stopped, stopping anti entropy procedure.")

	for {
		select {
		case <-s.stopCh:
			return
		case <-time.After(s.config.StateCheckInterval):
			ourHeight, err := s.ledger.LedgerHeight()
			if err != nil {
				
				logger.Errorf("Cannot obtain ledger height, due to %+v", errors.WithStack(err))
				continue
			}
			if ourHeight == 0 {
				logger.Error("Ledger reported block height of 0 but this should be impossible")
				continue
			}
			maxHeight := s.maxAvailableLedgerHeight()
			if ourHeight >= maxHeight {
				continue
			}

			s.requestBlocksInRange(uint64(ourHeight), uint64(maxHeight)-1)
		}
	}
}



func (s *GossipStateProviderImpl) maxAvailableLedgerHeight() uint64 {
	max := uint64(0)
	for _, p := range s.mediator.PeersOfChannel(common2.ChannelID(s.chainID)) {
		if p.Properties == nil {
			logger.Debug("Peer", p.PreferredEndpoint(), "doesn't have properties, skipping it")
			continue
		}
		peerHeight := p.Properties.LedgerHeight
		if max < peerHeight {
			max = peerHeight
		}
	}
	return max
}



func (s *GossipStateProviderImpl) requestBlocksInRange(start uint64, end uint64) {
	atomic.StoreInt32(&s.stateTransferActive, 1)
	defer atomic.StoreInt32(&s.stateTransferActive, 0)

	for prev := start; prev <= end; {
		next := min(end, prev+s.config.StateBatchSize)

		gossipMsg := s.stateRequestMessage(prev, next)

		responseReceived := false
		tryCounts := 0

		for !responseReceived {
			if tryCounts > s.config.StateMaxRetries {
				logger.Warningf("Wasn't  able to get blocks in range [%d...%d), after %d retries",
					prev, next, tryCounts)
				return
			}
			
			peer, err := s.selectPeerToRequestFrom(next)
			if err != nil {
				logger.Warningf("Cannot send state request for blocks in range [%d...%d), due to %+v",
					prev, next, errors.WithStack(err))
				return
			}

			logger.Debugf("State transfer, with peer %s, requesting blocks in range [%d...%d), "+
				"for chainID %s", peer.Endpoint, prev, next, s.chainID)

			s.mediator.Send(gossipMsg, peer)
			tryCounts++

			
			select {
			case msg, stillOpen := <-s.stateResponseCh:
				if !stillOpen {
					return
				}
				if msg.GetGossipMessage().Nonce !=
					gossipMsg.Nonce {
					continue
				}
				
				index, err := s.handleStateResponse(msg)
				if err != nil {
					logger.Warningf("Wasn't able to process state response for "+
						"blocks [%d...%d], due to %+v", prev, next, errors.WithStack(err))
					continue
				}
				prev = index + 1
				responseReceived = true
			case <-time.After(s.config.StateResponseTimeout):
			}
		}
	}
}


func (s *GossipStateProviderImpl) stateRequestMessage(beginSeq uint64, endSeq uint64) *proto.GossipMessage {
	return &proto.GossipMessage{
		Nonce:   util.RandomUInt64(),
		Tag:     proto.GossipMessage_CHAN_OR_ORG,
		Channel: []byte(s.chainID),
		Content: &proto.GossipMessage_StateRequest{
			StateRequest: &proto.RemoteStateRequest{
				StartSeqNum: beginSeq,
				EndSeqNum:   endSeq,
			},
		},
	}
}


func (s *GossipStateProviderImpl) selectPeerToRequestFrom(height uint64) (*comm.RemotePeer, error) {
	
	peers := s.filterPeers(s.hasRequiredHeight(height))

	n := len(peers)
	if n == 0 {
		return nil, errors.New("there are no peers to ask for missing blocks from")
	}

	
	return peers[util.RandomInt(n)], nil
}


func (s *GossipStateProviderImpl) filterPeers(predicate func(peer discovery.NetworkMember) bool) []*comm.RemotePeer {
	var peers []*comm.RemotePeer

	for _, member := range s.mediator.PeersOfChannel(common2.ChannelID(s.chainID)) {
		if predicate(member) {
			peers = append(peers, &comm.RemotePeer{Endpoint: member.PreferredEndpoint(), PKIID: member.PKIid})
		}
	}

	return peers
}



func (s *GossipStateProviderImpl) hasRequiredHeight(height uint64) func(peer discovery.NetworkMember) bool {
	return func(peer discovery.NetworkMember) bool {
		if peer.Properties != nil {
			return peer.Properties.LedgerHeight >= height
		}
		logger.Debug(peer.PreferredEndpoint(), "doesn't have properties")
		return false
	}
}


func (s *GossipStateProviderImpl) AddPayload(payload *proto.Payload) error {
	return s.addPayload(payload, s.blockingMode)
}





func (s *GossipStateProviderImpl) addPayload(payload *proto.Payload, blockingMode bool) error {
	if payload == nil {
		return errors.New("Given payload is nil")
	}
	logger.Debugf("[%s] Adding payload to local buffer, blockNum = [%d]", s.chainID, payload.SeqNum)
	height, err := s.ledger.LedgerHeight()
	if err != nil {
		return errors.Wrap(err, "Failed obtaining ledger height")
	}

	if !blockingMode && payload.SeqNum-height >= uint64(s.config.StateBlockBufferSize) {
		return errors.Errorf("Ledger height is at %d, cannot enqueue block with sequence of %d", height, payload.SeqNum)
	}

	for blockingMode && s.payloads.Size() > s.config.StateBlockBufferSize*2 {
		time.Sleep(enqueueRetryInterval)
	}

	s.payloads.Push(payload)
	logger.Debugf("Blocks payloads buffer size for channel [%s] is %d blocks", s.chainID, s.payloads.Size())
	return nil
}

func (s *GossipStateProviderImpl) commitBlock(block *common.Block, pvtData util.PvtDataCollections) error {

	t1 := time.Now()

	
	if err := s.ledger.StoreBlock(block, pvtData); err != nil {
		logger.Errorf("Got error while committing(%+v)", errors.WithStack(err))
		return err
	}

	sinceT1 := time.Since(t1)
	s.stateMetrics.CommitDuration.With("channel", s.chainID).Observe(sinceT1.Seconds())

	
	s.mediator.UpdateLedgerHeight(block.Header.Number+1, common2.ChannelID(s.chainID))
	logger.Debugf("[%s] Committed block [%d] with %d transaction(s)",
		s.chainID, block.Header.Number, len(block.Data.Data))

	s.stateMetrics.Height.With("channel", s.chainID).Set(float64(block.Header.Number + 1))

	return nil
}

func min(a uint64, b uint64) uint64 {
	return b ^ ((a ^ b) & (-(uint64(a-b) >> 63)))
}
