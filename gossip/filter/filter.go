/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package filter

import (
	"math/rand"

	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/util"
)

func init() { 
	rand.Seed(int64(util.RandomUInt64()))
}




type RoutingFilter func(discovery.NetworkMember) bool


var SelectNonePolicy = func(discovery.NetworkMember) bool {
	return false
}


var SelectAllPolicy = func(discovery.NetworkMember) bool {
	return true
}


func CombineRoutingFilters(filters ...RoutingFilter) RoutingFilter {
	return func(member discovery.NetworkMember) bool {
		for _, filter := range filters {
			if !filter(member) {
				return false
			}
		}
		return true
	}
}


func SelectPeers(k int, peerPool []discovery.NetworkMember, filter RoutingFilter) []*comm.RemotePeer {
	var res []*comm.RemotePeer
	
	for _, index := range rand.Perm(len(peerPool)) {
		
		if len(res) == k {
			break
		}
		peer := peerPool[index]
		
		if filter(peer) {
			p := &comm.RemotePeer{PKIID: peer.PKIid, Endpoint: peer.PreferredEndpoint()}
			res = append(res, p)
		}
	}
	return res
}


func First(peerPool []discovery.NetworkMember, filter RoutingFilter) *comm.RemotePeer {
	for _, p := range peerPool {
		if filter(p) {
			return &comm.RemotePeer{PKIID: p.PKIid, Endpoint: p.PreferredEndpoint()}
		}
	}
	return nil
}


func AnyMatch(peerPool []discovery.NetworkMember, filters ...RoutingFilter) []discovery.NetworkMember {
	var res []discovery.NetworkMember
	for _, peer := range peerPool {
		for _, matches := range filters {
			if matches(peer) {
				res = append(res, peer)
				break
			}
		}
	}
	return res
}
