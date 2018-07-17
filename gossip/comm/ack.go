/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/util"
	proto "github.com/mcc-github/blockchain/protos/gossip"
)

type sendFunc func(peer *RemotePeer, msg *proto.SignedGossipMessage)
type waitFunc func(*RemotePeer) error

type ackSendOperation struct {
	snd        sendFunc
	waitForAck waitFunc
}

func newAckSendOperation(snd sendFunc, waitForAck waitFunc) *ackSendOperation {
	return &ackSendOperation{
		snd:        snd,
		waitForAck: waitForAck,
	}
}

func (aso *ackSendOperation) send(msg *proto.SignedGossipMessage, minAckNum int, peers ...*RemotePeer) []SendResult {
	successAcks := 0
	results := []SendResult{}

	acks := make(chan SendResult, len(peers))
	
	for _, p := range peers {
		go func(p *RemotePeer) {
			
			aso.snd(p, msg)
			
			err := aso.waitForAck(p)
			acks <- SendResult{
				RemotePeer: *p,
				error:      err,
			}
		}(p)
	}
	for {
		ack := <-acks
		results = append(results, SendResult{
			error:      ack.error,
			RemotePeer: ack.RemotePeer,
		})
		if ack.error == nil {
			successAcks++
		}
		if successAcks == minAckNum || len(results) == len(peers) {
			break
		}
	}
	return results
}

func interceptAcks(nextHandler handler, remotePeerID common.PKIidType, pubSub *util.PubSub) func(*proto.SignedGossipMessage) {
	return func(m *proto.SignedGossipMessage) {
		if m.IsAck() {
			topic := topicForAck(m.Nonce, remotePeerID)
			pubSub.Publish(topic, m.GetAck())
			return
		}
		nextHandler(m)
	}
}
