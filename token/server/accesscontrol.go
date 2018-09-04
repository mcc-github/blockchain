/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package server

import (
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/pkg/errors"
)



type SignedDataPolicyChecker interface {
	
	
	CheckPolicyBySignedData(channelID, policyName string, sd []*common.SignedData) error
}


type PolicyBasedAccessControl struct {
	SignedDataPolicyChecker SignedDataPolicyChecker
}

func (ac *PolicyBasedAccessControl) Check(sc *token.SignedCommand, c *token.Command) error {
	switch t := c.GetPayload().(type) {

	case *token.Command_ImportRequest:
		return ac.SignedDataPolicyChecker.CheckPolicyBySignedData(
			c.Header.ChannelId,
			policies.ChannelApplicationWriters,
			[]*common.SignedData{{
				Identity:  c.Header.Creator,
				Data:      sc.Command,
				Signature: sc.Signature,
			}},
		)

	default:
		return errors.Errorf("command type not recognized: %T", t)
	}
}
