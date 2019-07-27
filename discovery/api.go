/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	discprotos "github.com/mcc-github/blockchain/protos/discovery"
	"github.com/mcc-github/blockchain/protoutil"
)


type AccessControlSupport interface {
	
	
	EligibleForService(channel string, data protoutil.SignedData) error
}


type ConfigSequenceSupport interface {
	
	ConfigSequence(channel string) uint64
}



type GossipSupport interface {
	
	ChannelExists(channel string) bool

	
	
	PeersOfChannel(common.ChannelID) discovery.Members

	
	Peers() discovery.Members

	
	IdentityInfo() api.PeerIdentitySet
}



type EndorsementSupport interface {
	
	PeersForEndorsement(channel common.ChannelID, interest *discprotos.ChaincodeInterest) (*discprotos.EndorsementDescriptor, error)

	
	
	
	
	PeersAuthorizedByCriteria(chainID common.ChannelID, interest *discprotos.ChaincodeInterest) (discovery.Members, error)
}


type ConfigSupport interface {
	
	Config(channel string) (*discprotos.ConfigResult, error)
}



type Support interface {
	AccessControlSupport
	GossipSupport
	EndorsementSupport
	ConfigSupport
	ConfigSequenceSupport
}
