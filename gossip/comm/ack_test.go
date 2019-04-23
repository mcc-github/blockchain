/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"errors"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/stretchr/testify/assert"
)

func TestInterceptAcks(t *testing.T) {
	pubsub := util.NewPubSub()
	pkiID := common.PKIidType("pkiID")
	msgs := make(chan *protoext.SignedGossipMessage, 1)
	handlerFunc := func(message *protoext.SignedGossipMessage) {
		msgs <- message
	}
	wrappedHandler := interceptAcks(handlerFunc, pkiID, pubsub)
	ack := &protoext.SignedGossipMessage{
		GossipMessage: &proto.GossipMessage{
			Nonce: 1,
			Content: &proto.GossipMessage_Ack{
				Ack: &proto.Acknowledgement{},
			},
		},
	}
	sub := pubsub.Subscribe(topicForAck(1, pkiID), time.Second)
	wrappedHandler(ack)
	
	assert.Len(t, msgs, 0)
	_, err := sub.Listen()
	
	assert.NoError(t, err)

	
	notAck := &protoext.SignedGossipMessage{
		GossipMessage: &proto.GossipMessage{
			Nonce: 2,
			Content: &proto.GossipMessage_DataMsg{
				DataMsg: &proto.DataMessage{},
			},
		},
	}
	sub = pubsub.Subscribe(topicForAck(2, pkiID), time.Second)
	wrappedHandler(notAck)
	
	assert.Len(t, msgs, 1)
	_, err = sub.Listen()
	
	assert.Error(t, err)
}

func TestAck(t *testing.T) {
	t.Parallel()

	comm1, _ := newCommInstance(t, naiveSec)
	comm2, port2 := newCommInstance(t, naiveSec)
	defer comm2.Stop()
	comm3, port3 := newCommInstance(t, naiveSec)
	defer comm3.Stop()
	comm4, port4 := newCommInstance(t, naiveSec)
	defer comm4.Stop()

	acceptData := func(o interface{}) bool {
		m := o.(protoext.ReceivedMessage).GetGossipMessage()
		return protoext.IsDataMsg(m.GossipMessage)
	}

	ack := func(c <-chan protoext.ReceivedMessage) {
		msg := <-c
		msg.Ack(nil)
	}

	nack := func(c <-chan protoext.ReceivedMessage) {
		msg := <-c
		msg.Ack(errors.New("Failed processing message because reasons"))
	}

	
	inc2 := comm2.Accept(acceptData)
	inc3 := comm3.Accept(acceptData)

	
	go ack(inc2)
	go ack(inc3)
	res := comm1.SendWithAck(createGossipMsg(), time.Second*3, 2, remotePeer(port2), remotePeer(port3))
	assert.Len(t, res, 2)
	assert.Empty(t, res[0].Error())
	assert.Empty(t, res[1].Error())

	
	t1 := time.Now()
	go ack(inc2)
	go ack(inc3)
	res = comm1.SendWithAck(createGossipMsg(), time.Second*10, 2, remotePeer(port2), remotePeer(port3), remotePeer(port4))
	elapsed := time.Since(t1)
	assert.Len(t, res, 2)
	assert.Empty(t, res[0].Error())
	assert.Empty(t, res[1].Error())
	
	assert.True(t, elapsed < time.Second*5)

	
	go ack(inc2)
	go nack(inc3)
	res = comm1.SendWithAck(createGossipMsg(), time.Second*10, 2, remotePeer(port2), remotePeer(port3), remotePeer(port4))
	assert.Len(t, res, 3)
	assert.Contains(t, []string{res[0].Error(), res[1].Error(), res[2].Error()}, "Failed processing message because reasons")
	assert.Contains(t, []string{res[0].Error(), res[1].Error(), res[2].Error()}, "timed out")

	
	res = comm1.SendWithAck(createGossipMsg(), time.Second*3, 2, remotePeer(port2), remotePeer(port3))
	assert.Len(t, res, 2)
	assert.Contains(t, res[0].Error(), "timed out")
	assert.Contains(t, res[1].Error(), "timed out")
	
	<-inc2
	<-inc3

	
	go ack(inc2)
	go nack(inc3)
	res = comm1.SendWithAck(createGossipMsg(), time.Second*3, 2, remotePeer(port2), remotePeer(port3), remotePeer(port4))
	assert.Len(t, res, 3)
	assert.Contains(t, []string{res[0].Error(), res[1].Error(), res[2].Error()}, "") 
	assert.Contains(t, []string{res[0].Error(), res[1].Error(), res[2].Error()}, "Failed processing message because reasons")
	assert.Contains(t, []string{res[0].Error(), res[1].Error(), res[2].Error()}, "timed out")
	assert.Contains(t, res.String(), "\"Failed processing message because reasons\":1")
	assert.Contains(t, res.String(), "\"timed out\":1")
	assert.Contains(t, res.String(), "\"successes\":1")
	assert.Equal(t, 2, res.NackCount())
	assert.Equal(t, 1, res.AckCount())

	
	res = comm1.SendWithAck(createGossipMsg(), time.Second*3, 1)
	assert.Len(t, res, 0)

	
	comm1.Stop()
	res = comm1.SendWithAck(createGossipMsg(), time.Second*3, 1, remotePeer(port2), remotePeer(port3), remotePeer(port4))
	assert.Len(t, res, 3)
	assert.Contains(t, res[0].Error(), "comm is stopping")
	assert.Contains(t, res[1].Error(), "comm is stopping")
	assert.Contains(t, res[2].Error(), "comm is stopping")
}
