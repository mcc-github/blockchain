/*
Copyright IBM Corp, SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"github.com/mcc-github/blockchain/protos/peer"
)



type Filter interface {
	peer.EndorserServer
	
	Init(next peer.EndorserServer)
}



func ChainFilters(endorser peer.EndorserServer, filters ...Filter) peer.EndorserServer {
	if len(filters) == 0 {
		return endorser
	}

	
	for i := 0; i < len(filters)-1; i++ {
		filters[i].Init(filters[i+1])
	}

	
	filters[len(filters)-1].Init(endorser)

	return filters[0]
}
