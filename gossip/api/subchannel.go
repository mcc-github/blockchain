/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import "github.com/mcc-github/blockchain/gossip/common"



type RoutingFilter func(peerIdentity PeerIdentityType) bool



type SubChannelSelectionCriteria func(signature PeerSignature) bool




type RoutingFilterFactory interface {
	
	Peers(common.ChannelID, SubChannelSelectionCriteria) RoutingFilter
}
