/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"github.com/mcc-github/blockchain/gossip/common"
)





type SecurityAdvisor interface {
	
	
	
	
	
	OrgByPeerIdentity(PeerIdentityType) OrgIdentityType
}



type ChannelNotifier interface {
	JoinChannel(joinMsg JoinChannelMessage, chainID common.ChainID)
}




type JoinChannelMessage interface {

	
	
	SequenceNumber() uint64

	
	Members() []OrgIdentityType

	
	AnchorPeersOf(org OrgIdentityType) []AnchorPeer
}


type AnchorPeer struct {
	Host string 
	Port int    
}


type OrgIdentityType []byte
