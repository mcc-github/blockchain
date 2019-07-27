/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"fmt"
	"time"

	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/filter"
	"github.com/mcc-github/blockchain/gossip/protoext"
)



type emittedGossipMessage struct {
	*protoext.SignedGossipMessage
	filter func(id common.PKIidType) bool
}


type SendCriteria struct {
	Timeout    time.Duration        
	MinAck     int                  
	MaxPeers   int                  
	IsEligible filter.RoutingFilter 
	Channel    common.ChannelID     
	
}


func (sc SendCriteria) String() string {
	return fmt.Sprintf("channel: %s, tout: %v, minAck: %d, maxPeers: %d", sc.Channel, sc.Timeout, sc.MinAck, sc.MaxPeers)
}
