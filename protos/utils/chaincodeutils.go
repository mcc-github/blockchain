/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)



func UnmarshalChaincodeDeploymentSpec(cdsBytes []byte) (*peer.ChaincodeDeploymentSpec, error) {
	cds := &peer.ChaincodeDeploymentSpec{}
	err := proto.Unmarshal(cdsBytes, cds)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeDeploymentSpec")
	}

	return cds, nil
}
