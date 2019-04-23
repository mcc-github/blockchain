/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"sync"

	"github.com/mcc-github/blockchain/gossip/protoext"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/pkg/errors"
)


type ReceivedMessageImpl struct {
	*protoext.SignedGossipMessage
	lock     sync.Locker
	conn     *connection
	connInfo *protoext.ConnectionInfo
}



func (m *ReceivedMessageImpl) GetSourceEnvelope() *proto.Envelope {
	return m.Envelope
}


func (m *ReceivedMessageImpl) Respond(msg *proto.GossipMessage) {
	sMsg, err := protoext.NoopSign(msg)
	if err != nil {
		err = errors.WithStack(err)
		m.conn.logger.Errorf("Failed creating SignedGossipMessage: %+v", err)
		return
	}
	m.conn.send(sMsg, func(e error) {}, blockingSend)
}


func (m *ReceivedMessageImpl) GetGossipMessage() *protoext.SignedGossipMessage {
	return m.SignedGossipMessage
}



func (m *ReceivedMessageImpl) GetConnectionInfo() *protoext.ConnectionInfo {
	return m.connInfo
}


func (m *ReceivedMessageImpl) Ack(err error) {
	ackMsg := &proto.GossipMessage{
		Nonce: m.GetGossipMessage().Nonce,
		Content: &proto.GossipMessage_Ack{
			Ack: &proto.Acknowledgement{},
		},
	}
	if err != nil {
		ackMsg.GetAck().Error = err.Error()
	}
	m.Respond(ackMsg)
}
