/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/protos/gossip"
)



type Gossip interface {
	
	IdentityInfo() api.PeerIdentitySet
	
	Peers() []discovery.NetworkMember
	
	
	PeersOfChannel(common.ChannelID) []discovery.NetworkMember
	
	SelfChannelInfo(common.ChannelID) *protoext.SignedGossipMessage
	
	SelfMembershipInfo() discovery.NetworkMember
}



type DiscoverySupport struct {
	Gossip
}


func NewDiscoverySupport(g Gossip) *DiscoverySupport {
	return &DiscoverySupport{g}
}


func (s *DiscoverySupport) ChannelExists(channel string) bool {
	return s.SelfChannelInfo(common.ChannelID(channel)) != nil
}



func (s *DiscoverySupport) PeersOfChannel(chain common.ChannelID) discovery.Members {
	msg := s.SelfChannelInfo(chain)
	if msg == nil {
		return nil
	}
	stateInf := msg.GetStateInfo()
	selfMember := discovery.NetworkMember{
		Properties: stateInf.Properties,
		PKIid:      stateInf.PkiId,
		Envelope:   msg.Envelope,
	}
	return append(s.Gossip.PeersOfChannel(chain), selfMember)
}


func (s *DiscoverySupport) Peers() discovery.Members {
	peers := s.Gossip.Peers()
	peers = append(peers, s.Gossip.SelfMembershipInfo())
	
	return discovery.Members(peers).Filter(discovery.HasExternalEndpoint).Map(sanitizeEnvelope)
}

func sanitizeEnvelope(member discovery.NetworkMember) discovery.NetworkMember {
	
	returnedMember := member
	if returnedMember.Envelope == nil {
		return returnedMember
	}
	returnedMember.Envelope = &gossip.Envelope{
		Payload:   member.Envelope.Payload,
		Signature: member.Envelope.Signature,
	}
	return returnedMember
}
