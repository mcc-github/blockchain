/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/pkg/errors"
)



type ACLProvider interface {
	
	
	
	CheckACL(resName string, channelID string, idinfo interface{}) error
}

type ACLResources struct {
	IssueTokens    string
	TransferTokens string
	ListTokens     string
}


type PolicyBasedAccessControl struct {
	ACLProvider  ACLProvider
	ACLResources *ACLResources
}

func (ac *PolicyBasedAccessControl) Check(sc *token.SignedCommand, c *token.Command) error {
	signedData := []*common.SignedData{{
		Identity:  c.Header.Creator,
		Data:      sc.Command,
		Signature: sc.Signature,
	}}

	switch t := c.GetPayload().(type) {

	case *token.Command_ImportRequest:
		return ac.ACLProvider.CheckACL(
			ac.ACLResources.IssueTokens,
			c.Header.ChannelId,
			signedData,
		)
	case *token.Command_ListRequest:
		return ac.ACLProvider.CheckACL(
			ac.ACLResources.ListTokens,
			c.Header.ChannelId,
			signedData,
		)
	case *token.Command_TransferRequest:
		return ac.ACLProvider.CheckACL(
			ac.ACLResources.TransferTokens,
			c.Header.ChannelId,
			signedData,
		)
	case *token.Command_RedeemRequest:
		
		return ac.ACLProvider.CheckACL(
			ac.ACLResources.TransferTokens,
			c.Header.ChannelId,
			signedData,
		)

	case *token.Command_ApproveRequest:
		
		return ac.ACLProvider.CheckACL(
			ac.ACLResources.TransferTokens,
			c.Header.ChannelId,
			signedData,
		)

	case *token.Command_TransferFromRequest:
		
		return ac.ACLProvider.CheckACL(
			ac.ACLResources.TransferTokens,
			c.Header.ChannelId,
			signedData,
		)

	case *token.Command_ExpectationRequest:
		if c.GetExpectationRequest().GetExpectation() == nil {
			return errors.New("ExpectationRequest has nil Expectation")
		}
		plainExpectation := c.GetExpectationRequest().GetExpectation().GetPlainExpectation()
		if plainExpectation == nil {
			return errors.New("ExpectationRequest has nil PlainExpectation")
		}
		return ac.checkExpectation(plainExpectation, signedData, c)
	default:
		return errors.Errorf("command type not recognized: %T", t)
	}
}


func (ac *PolicyBasedAccessControl) checkExpectation(plainExpectation *token.PlainExpectation, signedData []*common.SignedData, c *token.Command) error {
	switch t := plainExpectation.GetPayload().(type) {
	case *token.PlainExpectation_ImportExpectation:
		return ac.ACLProvider.CheckACL(
			ac.ACLResources.IssueTokens,
			c.Header.ChannelId,
			signedData,
		)
	case *token.PlainExpectation_TransferExpectation:
		return ac.ACLProvider.CheckACL(
			ac.ACLResources.TransferTokens,
			c.Header.ChannelId,
			signedData,
		)
	default:
		return errors.Errorf("expectation payload type not recognized: %T", t)
	}
}
