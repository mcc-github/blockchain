/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoext

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/discovery"
)




func SignedRequestToRequest(sr *discovery.SignedRequest) (*discovery.Request, error) {
	req := &discovery.Request{}
	return req, proto.Unmarshal(sr.Payload, req)
}
