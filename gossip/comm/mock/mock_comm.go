/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mock

import (
	"time"

	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/util"
	proto "github.com/mcc-github/blockchain/protos/gossip"
)


type socketMock struct {
	
	endpoint string

	
	socket chan interface{}
}


type packetMock struct {
	
	src *socketMock

	
	dst *socketMock

	msg interface{}
}

type channelMock struct {
	accept common.MessageAcceptor

	channel chan proto.ReceivedMessage
}

type commMock struct {
	id string

	members map[string]*socketMock

	acceptors []*channelMock

	deadChannel chan common.PKIidType

	done chan struct{}
}

var logger = util.GetLogger(util.CommMockLogger, "")


func NewCommMock(id string, members map[string]*socketMock) comm.Comm {
	res := &commMock{
		id: id,

		members: members,

		acceptors: make([]*channelMock, 0),

		done: make(chan struct{}),

		deadChannel: make(chan common.PKIidType),
	}
	
	go res.start()

	return res
}


func (packet *packetMock) Respond(msg *proto.GossipMessage) {
	sMsg, _ := msg.NoopSign()
	packet.src.socket <- &packetMock{
		src: packet.dst,
		dst: packet.src,
		msg: sMsg,
	}
}


func (packet *packetMock) Ack(err error) {

}



func (packet *packetMock) GetSourceEnvelope() *proto.Envelope {
	return nil
}


func (packet *packetMock) GetGossipMessage() *proto.SignedGossipMessage {
	return packet.msg.(*proto.SignedGossipMessage)
}



func (packet *packetMock) GetConnectionInfo() *proto.ConnectionInfo {
	return nil
}

func (mock *commMock) start() {
	logger.Debug("Starting communication mock module...")
	for {
		select {
		case <-mock.done:
			{
				
				logger.Debug("Exiting...")
				return
			}
		case msg := <-mock.members[mock.id].socket:
			{
				logger.Debug("Got new message", msg)
				packet := msg.(*packetMock)
				for _, channel := range mock.acceptors {
					
					
					
					if channel.accept(packet) {
						channel.channel <- packet
					}
				}
			}
		}
	}
}


func (mock *commMock) GetPKIid() common.PKIidType {
	return common.PKIidType(mock.id)
}


func (mock *commMock) Send(msg *proto.SignedGossipMessage, peers ...*comm.RemotePeer) {
	for _, peer := range peers {
		logger.Debug("Sending message to peer ", peer.Endpoint, "from ", mock.id)
		mock.members[peer.Endpoint].socket <- &packetMock{
			src: mock.members[mock.id],
			dst: mock.members[peer.Endpoint],
			msg: msg,
		}
	}
}

func (mock *commMock) SendWithAck(_ *proto.SignedGossipMessage, _ time.Duration, _ int, _ ...*comm.RemotePeer) comm.AggregatedSendResult {
	panic("not implemented")
}



func (mock *commMock) Probe(peer *comm.RemotePeer) error {
	return nil
}



func (mock *commMock) Handshake(peer *comm.RemotePeer) (api.PeerIdentityType, error) {
	return nil, nil
}



func (mock *commMock) Accept(accept common.MessageAcceptor) <-chan proto.ReceivedMessage {
	ch := make(chan proto.ReceivedMessage)
	mock.acceptors = append(mock.acceptors, &channelMock{accept, ch})
	return ch
}


func (mock *commMock) PresumedDead() <-chan common.PKIidType {
	return mock.deadChannel
}


func (mock *commMock) CloseConn(peer *comm.RemotePeer) {
	
}


func (mock *commMock) Stop() {
	logger.Debug("Stopping communication module, closing all accepting channels.")
	for _, accept := range mock.acceptors {
		close(accept.channel)
	}
	logger.Debug("[XXX]: Sending done signal to close the module.")
	mock.done <- struct{}{}
}
