/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"github.com/golang/protobuf/proto"
)


type QueryType uint8

const (
	InvalidQueryType QueryType = iota
	ConfigQueryType
	PeerMembershipQueryType
	ChaincodeQueryType
	LocalMembershipQueryType
)


func (q *Query) GetType() QueryType {
	if q.GetCcQuery() != nil {
		return ChaincodeQueryType
	}
	if q.GetConfigQuery() != nil {
		return ConfigQueryType
	}
	if q.GetPeerQuery() != nil {
		return PeerMembershipQueryType
	}
	if q.GetLocalPeers() != nil {
		return LocalMembershipQueryType
	}
	return InvalidQueryType
}




func (sr *SignedRequest) ToRequest() (*Request, error) {
	req := &Request{}
	return req, proto.Unmarshal(sr.Payload, req)
}



func (m *Response) ConfigAt(i int) (*ConfigResult, *Error) {
	r := m.Results[i]
	return r.GetConfigResult(), r.GetError()
}



func (m *Response) MembershipAt(i int) (*PeerMembershipResult, *Error) {
	r := m.Results[i]
	return r.GetMembers(), r.GetError()
}



func (m *Response) EndorsersAt(i int) (*ChaincodeQueryResult, *Error) {
	r := m.Results[i]
	return r.GetCcQueryRes(), r.GetError()
}
