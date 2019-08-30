/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoext

import "github.com/mcc-github/blockchain-protos-go/discovery"


type QueryType uint8

const (
	InvalidQueryType QueryType = iota
	ConfigQueryType
	PeerMembershipQueryType
	ChaincodeQueryType
	LocalMembershipQueryType
)


func GetQueryType(q *discovery.Query) QueryType {
	switch {
	case q.GetCcQuery() != nil:
		return ChaincodeQueryType
	case q.GetConfigQuery() != nil:
		return ConfigQueryType
	case q.GetPeerQuery() != nil:
		return PeerMembershipQueryType
	case q.GetLocalPeers() != nil:
		return LocalMembershipQueryType
	default:
		return InvalidQueryType
	}
}
