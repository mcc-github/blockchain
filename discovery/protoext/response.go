/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoext

import "github.com/mcc-github/blockchain-protos-go/discovery"



func ResponseConfigAt(m *discovery.Response, i int) (*discovery.ConfigResult, *discovery.Error) {
	r := m.Results[i]
	return r.GetConfigResult(), r.GetError()
}



func ResponseMembershipAt(m *discovery.Response, i int) (*discovery.PeerMembershipResult, *discovery.Error) {
	r := m.Results[i]
	return r.GetMembers(), r.GetError()
}



func ResponseEndorsersAt(m *discovery.Response, i int) (*discovery.ChaincodeQueryResult, *discovery.Error) {
	r := m.Results[i]
	return r.GetCcQueryRes(), r.GetError()
}
