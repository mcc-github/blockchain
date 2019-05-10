/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package support

import (
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/discovery"
	"github.com/mcc-github/blockchain/discovery/support/acl"
)




type GossipSupport interface {
	discovery.GossipSupport
}




type ChannelPolicyManagerGetter interface {
	acl.ChannelPolicyManagerGetter
}



type PolicyManager interface {
	policies.Manager
}


type DiscoverySupport struct {
	discovery.AccessControlSupport
	discovery.GossipSupport
	discovery.EndorsementSupport
	discovery.ConfigSupport
	discovery.ConfigSequenceSupport
}


func NewDiscoverySupport(
	access discovery.AccessControlSupport,
	gossip discovery.GossipSupport,
	endorsement discovery.EndorsementSupport,
	config discovery.ConfigSupport,
	sequence discovery.ConfigSequenceSupport,
) *DiscoverySupport {
	return &DiscoverySupport{
		AccessControlSupport:  access,
		GossipSupport:         gossip,
		EndorsementSupport:    endorsement,
		ConfigSupport:         config,
		ConfigSequenceSupport: sequence,
	}
}
