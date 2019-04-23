/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoext

import (
	"fmt"

	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/protos/gossip"
)







type ReceivedMessage interface {
	
	Respond(msg *gossip.GossipMessage)

	
	GetGossipMessage() *SignedGossipMessage

	
	
	GetSourceEnvelope() *gossip.Envelope

	
	
	GetConnectionInfo() *ConnectionInfo

	
	
	
	Ack(err error)
}



type ConnectionInfo struct {
	ID       common.PKIidType
	Auth     *AuthInfo
	Identity api.PeerIdentityType
	Endpoint string
}


func (c *ConnectionInfo) String() string {
	return fmt.Sprintf("%s %v", c.Endpoint, c.ID)
}




type AuthInfo struct {
	SignedData []byte
	Signature  []byte
}
