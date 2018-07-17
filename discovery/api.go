/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	common2 "github.com/mcc-github/blockchain/protos/common"
	discovery2 "github.com/mcc-github/blockchain/protos/discovery"
)


type AccessControlSupport interface {
	
	
	EligibleForService(channel string, data common2.SignedData) error
}


type ConfigSequenceSupport interface {
	
	ConfigSequence(channel string) uint64
}





type GossipSupport interface {
	
	ChannelExists(channel string) bool

	
	
	PeersOfChannel(common.ChainID) discovery.Members

	
	Peers() discovery.Members

	
	IdentityInfo() api.PeerIdentitySet
}



type EndorsementSupport interface {
	
	PeersForEndorsement(channel common.ChainID, interest *discovery2.ChaincodeInterest) (*discovery2.EndorsementDescriptor, error)
}


type ConfigSupport interface {
	
	Config(channel string) (*discovery2.ConfigResult, error)
}



type Support interface {
	AccessControlSupport
	GossipSupport
	EndorsementSupport
	ConfigSupport
	ConfigSequenceSupport
}
