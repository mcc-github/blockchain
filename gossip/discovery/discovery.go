/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"fmt"

	"github.com/mcc-github/blockchain/gossip/common"
	proto "github.com/mcc-github/blockchain/protos/gossip"
)


type CryptoService interface {
	
	ValidateAliveMsg(message *proto.SignedGossipMessage) bool

	
	SignMessage(m *proto.GossipMessage, internalEndpoint string) *proto.Envelope
}



type EnvelopeFilter func(message *proto.SignedGossipMessage) *proto.Envelope




type Sieve func(message *proto.SignedGossipMessage) bool










type DisclosurePolicy func(remotePeer *NetworkMember) (Sieve, EnvelopeFilter)


type CommService interface {
	
	Gossip(msg *proto.SignedGossipMessage)

	
	
	SendToPeer(peer *NetworkMember, msg *proto.SignedGossipMessage)

	
	Ping(peer *NetworkMember) bool

	
	Accept() <-chan proto.ReceivedMessage

	
	PresumedDead() <-chan common.PKIidType

	
	CloseConn(peer *NetworkMember)

	
	
	Forward(msg proto.ReceivedMessage)
}


type NetworkMember struct {
	Endpoint         string
	Metadata         []byte
	PKIid            common.PKIidType
	InternalEndpoint string
	Properties       *proto.Properties
	*proto.Envelope
}


func (n *NetworkMember) String() string {
	return fmt.Sprintf("Endpoint: %s, InternalEndpoint: %s, PKI-ID: %v, Metadata: %v", n.Endpoint, n.InternalEndpoint, n.PKIid, n.Metadata)
}




func (n NetworkMember) PreferredEndpoint() string {
	if n.InternalEndpoint != "" {
		return n.InternalEndpoint
	}
	return n.Endpoint
}




type PeerIdentification struct {
	ID      common.PKIidType
	SelfOrg bool
}

type identifier func() (*PeerIdentification, error)


type Discovery interface {

	
	Lookup(PKIID common.PKIidType) *NetworkMember

	
	Self() NetworkMember

	
	UpdateMetadata([]byte)

	
	UpdateEndpoint(string)

	
	Stop()

	
	GetMembership() []NetworkMember

	
	
	InitiateSync(peerNum int)

	
	
	
	
	Connect(member NetworkMember, id identifier)
}


type Members []NetworkMember



func (members Members) ByID() map[string]NetworkMember {
	res := make(map[string]NetworkMember, len(members))
	for _, peer := range members {
		res[string(peer.PKIid)] = peer
	}
	return res
}


func (members Members) Intersect(otherMembers Members) Members {
	var res Members
	m := otherMembers.ByID()
	for _, member := range members {
		if _, exists := m[string(member.PKIid)]; exists {
			res = append(res, member)
		}
	}
	return res
}


func (members Members) Filter(filter func(member NetworkMember) bool) Members {
	var res Members
	for _, member := range members {
		if filter(member) {
			res = append(res, member)
		}
	}
	return res
}


func (members Members) Map(f func(member NetworkMember) NetworkMember) Members {
	var res Members
	for _, m := range members {
		res = append(res, f(m))
	}
	return res
}


func HasExternalEndpoint(member NetworkMember) bool {
	return member.Endpoint != ""
}
