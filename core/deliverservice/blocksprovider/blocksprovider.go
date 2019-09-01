/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blocksprovider

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	gossip_proto "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/gossip/api"
	gossipcommon "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/util"
)



type LedgerInfo interface {
	
	LedgerHeight() (uint64, error)
}



type GossipServiceAdapter interface {
	
	PeersOfChannel(gossipcommon.ChannelID) []discovery.NetworkMember

	
	AddPayload(chainID string, payload *gossip_proto.Payload) error

	
	Gossip(msg *gossip_proto.GossipMessage)
}



type BlocksProvider interface {
	
	DeliverBlocks()

	
	UpdateOrderingEndpoints(endpoints []string)

	
	Stop()
}






type BlocksDeliverer interface {
	
	Recv() (*orderer.DeliverResponse, error)

	
	Send(*common.Envelope) error
}

type streamClient interface {
	BlocksDeliverer

	
	UpdateEndpoints(endpoints []string)

	
	GetEndpoints() []string

	
	Close()

	
	Disconnect()
}


type blocksProviderImpl struct {
	chainID string

	client streamClient

	gossip GossipServiceAdapter

	mcs api.MessageCryptoService

	done int32

	wrongStatusThreshold int
}

const wrongStatusThreshold = 10

var maxRetryDelay = time.Second * 10
var logger = flogging.MustGetLogger("blocksProvider")


func NewBlocksProvider(chainID string, client streamClient, gossip GossipServiceAdapter, mcs api.MessageCryptoService) BlocksProvider {
	return &blocksProviderImpl{
		chainID:              chainID,
		client:               client,
		gossip:               gossip,
		mcs:                  mcs,
		wrongStatusThreshold: wrongStatusThreshold,
	}
}



func (b *blocksProviderImpl) DeliverBlocks() {
	errorStatusCounter := 0
	statusCounter := 0
	defer b.client.Close()
	for !b.isDone() {
		msg, err := b.client.Recv()
		if err != nil {
			logger.Warningf("[%s] Receive error: %s", b.chainID, err.Error())
			return
		}
		switch t := msg.Type.(type) {
		case *orderer.DeliverResponse_Status:
			if t.Status == common.Status_SUCCESS {
				logger.Warningf("[%s] ERROR! Received success for a seek that should never complete", b.chainID)
				return
			}
			if t.Status == common.Status_BAD_REQUEST || t.Status == common.Status_FORBIDDEN {
				logger.Errorf("[%s] Got error %v", b.chainID, t)
				errorStatusCounter++
				if errorStatusCounter > b.wrongStatusThreshold {
					logger.Criticalf("[%s] Wrong statuses threshold passed, stopping block provider", b.chainID)
					return
				}
			} else {
				errorStatusCounter = 0
				logger.Warningf("[%s] Got error %v", b.chainID, t)
			}
			maxDelay := float64(maxRetryDelay)
			currDelay := float64(time.Duration(math.Pow(2, float64(statusCounter))) * 100 * time.Millisecond)
			time.Sleep(time.Duration(math.Min(maxDelay, currDelay)))
			if currDelay < maxDelay {
				statusCounter++
			}
			b.client.Disconnect()
			continue
		case *orderer.DeliverResponse_Block:
			errorStatusCounter = 0
			statusCounter = 0
			blockNum := t.Block.Header.Number

			if err := b.mcs.VerifyBlock(gossipcommon.ChannelID(b.chainID), blockNum, t.Block); err != nil {
				logger.Errorf("[%s] Error verifying block with sequnce number %d, due to %s", b.chainID, blockNum, err)
				continue
			}

			marshaledBlock, err := proto.Marshal(t.Block)
			if err != nil {
				logger.Errorf("[%s] Error serializing block with sequence number %d, due to %s", b.chainID, blockNum, err)
				continue
			}

			numberOfPeers := len(b.gossip.PeersOfChannel(gossipcommon.ChannelID(b.chainID)))
			
			payload := createPayload(blockNum, marshaledBlock)
			
			gossipMsg := createGossipMsg(b.chainID, payload)

			logger.Debugf("[%s] Adding payload to local buffer, blockNum = [%d]", b.chainID, blockNum)
			
			if err := b.gossip.AddPayload(b.chainID, payload); err != nil {
				logger.Warningf("Block [%d] received from ordering service wasn't added to payload buffer: %v", blockNum, err)
			}

			
			logger.Debugf("[%s] Gossiping block [%d], peers number [%d]", b.chainID, blockNum, numberOfPeers)
			if !b.isDone() {
				b.gossip.Gossip(gossipMsg)
			}
		default:
			logger.Warningf("[%s] Received unknown: %v", b.chainID, t)
			return
		}
	}
}


func (b *blocksProviderImpl) Stop() {
	atomic.StoreInt32(&b.done, 1)
	b.client.Close()
}


func (b *blocksProviderImpl) UpdateOrderingEndpoints(endpoints []string) {
	if !b.areEndpointsUpdated(endpoints) {
		
		return
	}
	
	logger.Debug("Updating endpoint, to %s", endpoints)
	b.client.UpdateEndpoints(endpoints)
	logger.Debug("Disconnecting so endpoints update will take effect")
	
	
	b.client.Disconnect()
}

func (b *blocksProviderImpl) areEndpointsUpdated(endpoints []string) bool {
	if len(endpoints) != len(b.client.GetEndpoints()) {
		return true
	}
	
	for _, endpoint := range endpoints {
		if !util.Contains(endpoint, b.client.GetEndpoints()) {
			
			return true
		}
	}
	
	return false
}


func (b *blocksProviderImpl) isDone() bool {
	return atomic.LoadInt32(&b.done) == 1
}

func createGossipMsg(chainID string, payload *gossip_proto.Payload) *gossip_proto.GossipMessage {
	gossipMsg := &gossip_proto.GossipMessage{
		Nonce:   0,
		Tag:     gossip_proto.GossipMessage_CHAN_AND_ORG,
		Channel: []byte(chainID),
		Content: &gossip_proto.GossipMessage_DataMsg{
			DataMsg: &gossip_proto.DataMessage{
				Payload: payload,
			},
		},
	}
	return gossipMsg
}

func createPayload(seqNum uint64, marshaledBlock []byte) *gossip_proto.Payload {
	return &gossip_proto.Payload{
		Data:   marshaledBlock,
		SeqNum: seqNum,
	}
}
