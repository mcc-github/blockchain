/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/protoext"
)



type Comm interface {

	
	GetPKIid() common.PKIidType

	
	Send(msg *protoext.SignedGossipMessage, peers ...*RemotePeer)

	
	SendWithAck(msg *protoext.SignedGossipMessage, timeout time.Duration, minAck int, peers ...*RemotePeer) AggregatedSendResult

	
	
	Probe(peer *RemotePeer) error

	
	
	Handshake(peer *RemotePeer) (api.PeerIdentityType, error)

	
	
	Accept(common.MessageAcceptor) <-chan protoext.ReceivedMessage

	
	PresumedDead() <-chan common.PKIidType

	
	CloseConn(peer *RemotePeer)

	
	Stop()
}


type RemotePeer struct {
	Endpoint string
	PKIID    common.PKIidType
}


type SendResult struct {
	error
	RemotePeer
}



func (sr SendResult) Error() string {
	if sr.error != nil {
		return sr.error.Error()
	}
	return ""
}


type AggregatedSendResult []SendResult


func (ar AggregatedSendResult) AckCount() int {
	c := 0
	for _, ack := range ar {
		if ack.error == nil {
			c++
		}
	}
	return c
}


func (ar AggregatedSendResult) NackCount() int {
	return len(ar) - ar.AckCount()
}



func (ar AggregatedSendResult) String() string {
	errMap := map[string]int{}
	for _, ack := range ar {
		if ack.error == nil {
			continue
		}
		errMap[ack.Error()]++
	}

	ackCount := ar.AckCount()
	output := map[string]interface{}{}
	if ackCount > 0 {
		output["successes"] = ackCount
	}
	if ackCount < len(ar) {
		output["failures"] = errMap
	}
	b, _ := json.Marshal(output)
	return string(b)
}


func (p *RemotePeer) String() string {
	return fmt.Sprintf("%s, PKIid:%v", p.Endpoint, p.PKIID)
}
