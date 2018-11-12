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
	vsccErrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	common2 "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/mcc-github/blockchain/protos/common"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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


type GossipAdapter interface {
	
	Send(msg *proto.GossipMessage, peers ...*comm.RemotePeer)

	
	
	
	
	Accept(acceptor common2.MessageAcceptor, passThrough bool) (<-chan *proto.GossipMessage, <-chan proto.ReceivedMessage)

	
	
	UpdateLedgerHeight(height uint64, chainID common2.ChainID)

	
	
	PeersOfChannel(common2.ChainID) []discovery.NetworkMember
}



type MCSAdapter interface {
	
	
	
	VerifyBlock(chainID common2.ChainID, seqNum uint64, signedBlock []byte) error

	
	
	
	
	VerifyByChannel(chainID common2.ChainID, peerIdentity api.PeerIdentityType, signature, message []byte) error
}


type ledgerResources interface {
	
	
	StoreBlock(block *common.Block, data util.PvtDataCollections) error

	
	StorePvtData(txid string, privData *transientstore.TxPvtReadWriteSetWithConfigInfo, blckHeight uint64) error

	
	
	
	
	GetPvtDataAndBlockByNum(seqNum uint64, peerAuthInfo common.SignedData) (*common.Block, util.PvtDataCollections, error)

	
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

	
	gossipChan <-chan *proto.GossipMessage

	commChan <-chan proto.ReceivedMessage

	
	payloads PayloadsBuffer

	ledger ledgerResources

	stateResponseCh chan proto.ReceivedMessage

	stateRequestCh chan proto.ReceivedMessage

	stopCh chan struct{}

	done sync.WaitGroup

	once sync.Once

	stateTransferActive int32
}

var logger = util.GetLogger(util.StateLogger, "")



func NewGossipStateProvider(chainID string, services *ServicesMediator, ledger ledgerResources) GossipStateProvider {

	gossipChan, _ := services.Accept(func(message interface{}) bool {
		
		return message.(*proto.GossipMessage).IsDataMsg() &&
			bytes.Equal(message.(*proto.GossipMessage).Channel, []byte(chainID))
	}, false)

	remoteStateMsgFilter := func(message interface{}) bool {
		receivedMsg := message.(proto.ReceivedMessage)
		msg := receivedMsg.GetGossipMessage()
		if !(msg.IsRemoteStateMessage() || msg.GetPrivateData() != nil) {
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

		
		gossipChan: gossipChan,

		
		commChan: commChan,

		
		payloads: NewPayloadsBuffer(height),

		ledger: ledger,

		stateResponseCh: make(chan proto.ReceivedMessage, defChannelBufferSize),

		stateRequestCh: make(chan proto.ReceivedMessage, defChannelBufferSize),

		stopCh: make(chan struct{}, 1),

		stateTransferActive: 0,

		once: sync.Once{},
	}

	logger.Infof("Updating metadata information, "+
		"current ledger sequence is at = %d, next expected block is = %d", height-1, s.payloads.Next())
	logger.Debug("Updating gossip ledger height to", height)
	services.UpdateLedgerHeight(height, common2.ChainID(s.chainID))

	s.done.Add(4)

	
	go s.listen()
	
	go s.deliverPayloads()
	
	go s.antiEntropy()
	
	go s.processStateRequests()

	return s
}

func (s *GossipStateProviderImpl) listen() {
	defer s.done.Done()

	for {
		select {
		case msg := <-s.gossipChan:
			logger.Debug("Received new message via gossip channel")
			go s.queueNewMessage(msg)
		case msg := <-s.commChan:
			logger.Debug("Dispatching a message", msg)
			go s.dispatch(msg)
		case <-s.stopCh:
			s.stopCh <- struct{}{}
			logger.Debug("Stop listening for new messages")
			return
		}
	}
}
func (s *GossipStateProviderImpl) dispatch(msg proto.ReceivedMessage) {
	
	if msg.GetGossipMessage().IsRemoteStateMessage() {
		logger.Debug("Handling direct state transfer message")
		
		s.directMessage(msg)
	} else if msg.GetGossipMessage().GetPrivateData() != nil {
		logger.Debug("Handling private data collection message")
		
		s.privateDataMessage(msg)
	}

}
func (s *GossipStateProviderImpl) privateDataMessage(msg proto.ReceivedMessage) {
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

func (s *GossipStateProviderImpl) directMessage(msg proto.ReceivedMessage) {
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
		if len(s.stateRequestCh) < defChannelBufferSize {
			
			
			s.stateRequestCh <- msg
		}
	} else if incoming.GetStateResponse() != nil {
		
		
		if atomic.LoadInt32(&s.stateTransferActive) == 1 {
			
			s.stateResponseCh <- msg
		}
	}
}

func (s *GossipStateProviderImpl) processStateRequests() {
	defer s.done.Done()

	for {
		select {
		case msg := <-s.stateRequestCh:
			s.handleStateRequest(msg)
		case <-s.stopCh:
			s.stopCh <- struct{}{}
			return
		}
	}
}



func (s *GossipStateProviderImpl) handleStateRequest(msg proto.ReceivedMessage) {
	if msg == nil {
		return
	}
	request := msg.GetGossipMessage().GetStateRequest()

	batchSize := request.EndSeqNum - request.StartSeqNum
	if batchSize > defAntiEntropyBatchSize {
		logger.Errorf("Requesting blocks batchSize size (%d) greater than configured allowed"+
			" (%d) batching for anti-entropy. Ignoring request...", batchSize, defAntiEntropyBatchSize)
		return
	}

	if request.StartSeqNum > request.EndSeqNum {
		logger.Errorf("Invalid sequence interval [%d...%d], ignoring request...", request.StartSeqNum, request.EndSeqNum)
		return
	}

	currentHeight, err := s.ledger.LedgerHeight()
	if err != nil {
		logger.Errorf("Cannot access to current ledger height, due to %+v", errors.WithStack(err))
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
		peerAuthInfo := common.SignedData{
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

func (s *GossipStateProviderImpl) handleStateResponse(msg proto.ReceivedMessage) (uint64, error) {
	max := uint64(0)
	
	response := msg.GetGossipMessage().GetStateResponse()
	
	if len(response.GetPayloads()) == 0 {
		return uint64(0), errors.New("Received state transfer response without payload")
	}
	for _, payload := range response.GetPayloads() {
		logger.Debugf("Received payload with sequence number %d.", payload.SeqNum)
		if err := s.mediator.VerifyBlock(common2.ChainID(s.chainID), payload.SeqNum, payload.Data); err != nil {
			err = errors.WithStack(err)
			logger.Warningf("Error verifying block with sequence number %d, due to %+v", payload.SeqNum, err)
			return uint64(0), err
		}
		if max < payload.SeqNum {
			max = payload.SeqNum
		}

		err := s.addPayload(payload, blocking)
		if err != nil {
			logger.Warningf("Block [%d] received from block transfer wasn't added to payload buffer: %v", payload.SeqNum, err)
		}
	}
	return max, nil
}


func (s *GossipStateProviderImpl) Stop() {
	
	
	s.once.Do(func() {
		s.stopCh <- struct{}{}
		
		s.done.Wait()
		
		s.ledger.Close()
		close(s.stateRequestCh)
		close(s.stateResponseCh)
		close(s.stopCh)
	})
}


func (s *GossipStateProviderImpl) queueNewMessage(msg *proto.GossipMessage) {
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
}

func (s *GossipStateProviderImpl) deliverPayloads() {
	defer s.done.Done()

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
			s.stopCh <- struct{}{}
			logger.Debug("State provider has been stopped, finishing to push new blocks.")
			return
		}
	}
}

func (s *GossipStateProviderImpl) antiEntropy() {
	defer s.done.Done()
	defer logger.Debug("State Provider stopped, stopping anti entropy procedure.")

	for {
		select {
		case <-s.stopCh:
			s.stopCh <- struct{}{}
			return
		case <-time.After(defAntiEntropyInterval):
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
	for _, p := range s.mediator.PeersOfChannel(common2.ChainID(s.chainID)) {
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
		next := min(end, prev+defAntiEntropyBatchSize)

		gossipMsg := s.stateRequestMessage(prev, next)

		responseReceived := false
		tryCounts := 0

		for !responseReceived {
			if tryCounts > defAntiEntropyMaxRetries {
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
			case msg := <-s.stateResponseCh:
				if msg.GetGossipMessage().Nonce != gossipMsg.Nonce {
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
			case <-time.After(defAntiEntropyStateResponseTimeout):
			case <-s.stopCh:
				s.stopCh <- struct{}{}
				return
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

	for _, member := range s.mediator.PeersOfChannel(common2.ChainID(s.chainID)) {
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
	blockingMode := blocking
	if viper.GetBool("peer.gossip.nonBlockingCommitMode") {
		blockingMode = false
	}
	return s.addPayload(payload, blockingMode)
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

	if !blockingMode && payload.SeqNum-height >= defMaxBlockDistance {
		return errors.Errorf("Ledger height is at %d, cannot enqueue block with sequence of %d", height, payload.SeqNum)
	}

	for blockingMode && s.payloads.Size() > defMaxBlockDistance*2 {
		time.Sleep(enqueueRetryInterval)
	}

	s.payloads.Push(payload)
	return nil
}

func (s *GossipStateProviderImpl) commitBlock(block *common.Block, pvtData util.PvtDataCollections) error {

	
	if err := s.ledger.StoreBlock(block, pvtData); err != nil {
		logger.Errorf("Got error while committing(%+v)", errors.WithStack(err))
		return err
	}

	
	s.mediator.UpdateLedgerHeight(block.Header.Number+1, common2.ChainID(s.chainID))
	logger.Debugf("[%s] Committed block [%d] with %d transaction(s)",
		s.chainID, block.Header.Number, len(block.Data.Data))

	return nil
}

func min(a uint64, b uint64) uint64 {
	return b ^ ((a ^ b) & (-(uint64(a-b) >> 63)))
}
