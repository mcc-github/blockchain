/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	gossip2 "github.com/mcc-github/blockchain/gossip/gossip"
	"github.com/mcc-github/blockchain/protos/gossip"
)



type DiscoverySupport struct {
	gossip2.Gossip
}


func NewDiscoverySupport(g gossip2.Gossip) *DiscoverySupport {
	return &DiscoverySupport{g}
}


func (s *DiscoverySupport) ChannelExists(channel string) bool {
	return s.SelfChannelInfo(common.ChainID(channel)) != nil
}



func (s *DiscoverySupport) PeersOfChannel(chain common.ChainID) discovery.Members {
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
