/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package policy

import (
	"errors"
	"fmt"

	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protoutil"
)



type PolicyChecker interface {
	
	
	
	CheckPolicy(channelID, policyName string, signedProp *pb.SignedProposal) error

	
	
	
	CheckPolicyBySignedData(channelID, policyName string, sd []*protoutil.SignedData) error

	
	
	CheckPolicyNoChannel(policyName string, signedProp *pb.SignedProposal) error
}

type policyChecker struct {
	channelPolicyManagerGetter policies.ChannelPolicyManagerGetter
	localMSP                   msp.IdentityDeserializer
	principalGetter            mgmt.MSPPrincipalGetter
}


func NewPolicyChecker(channelPolicyManagerGetter policies.ChannelPolicyManagerGetter, localMSP msp.IdentityDeserializer, principalGetter mgmt.MSPPrincipalGetter) PolicyChecker {
	return &policyChecker{channelPolicyManagerGetter, localMSP, principalGetter}
}



func (p *policyChecker) CheckPolicy(channelID, policyName string, signedProp *pb.SignedProposal) error {
	if channelID == "" {
		return p.CheckPolicyNoChannel(policyName, signedProp)
	}

	if policyName == "" {
		return fmt.Errorf("Invalid policy name during check policy on channel [%s]. Name must be different from nil.", channelID)
	}

	if signedProp == nil {
		return fmt.Errorf("Invalid signed proposal during check policy on channel [%s] with policy [%s]", channelID, policyName)
	}

	
	policyManager := p.channelPolicyManagerGetter.Manager(channelID)
	if policyManager == nil {
		return fmt.Errorf("Failed to get policy manager for channel [%s]", channelID)
	}

	
	proposal, err := protoutil.UnmarshalProposal(signedProp.ProposalBytes)
	if err != nil {
		return fmt.Errorf("Failing extracting proposal during check policy on channel [%s] with policy [%s]: [%s]", channelID, policyName, err)
	}

	header, err := protoutil.UnmarshalHeader(proposal.Header)
	if err != nil {
		return fmt.Errorf("Failing extracting header during check policy on channel [%s] with policy [%s]: [%s]", channelID, policyName, err)
	}

	shdr, err := protoutil.UnmarshalSignatureHeader(header.SignatureHeader)
	if err != nil {
		return fmt.Errorf("Invalid Proposal's SignatureHeader during check policy on channel [%s] with policy [%s]: [%s]", channelID, policyName, err)
	}

	sd := []*protoutil.SignedData{{
		Data:      signedProp.ProposalBytes,
		Identity:  shdr.Creator,
		Signature: signedProp.Signature,
	}}

	return p.CheckPolicyBySignedData(channelID, policyName, sd)
}



func (p *policyChecker) CheckPolicyNoChannel(policyName string, signedProp *pb.SignedProposal) error {
	if policyName == "" {
		return errors.New("Invalid policy name during channelless check policy. Name must be different from nil.")
	}

	if signedProp == nil {
		return fmt.Errorf("Invalid signed proposal during channelless check policy with policy [%s]", policyName)
	}

	proposal, err := protoutil.UnmarshalProposal(signedProp.ProposalBytes)
	if err != nil {
		return fmt.Errorf("Failing extracting proposal during channelless check policy with policy [%s]: [%s]", policyName, err)
	}

	header, err := protoutil.UnmarshalHeader(proposal.Header)
	if err != nil {
		return fmt.Errorf("Failing extracting header during channelless check policy with policy [%s]: [%s]", policyName, err)
	}

	shdr, err := protoutil.UnmarshalSignatureHeader(header.SignatureHeader)
	if err != nil {
		return fmt.Errorf("Invalid Proposal's SignatureHeader during channelless check policy with policy [%s]: [%s]", policyName, err)
	}

	
	id, err := p.localMSP.DeserializeIdentity(shdr.Creator)
	if err != nil {
		return fmt.Errorf("Failed deserializing proposal creator during channelless check policy with policy [%s]: [%s]", policyName, err)
	}

	
	principal, err := p.principalGetter.Get(policyName)
	if err != nil {
		return fmt.Errorf("Failed getting local MSP principal during channelless check policy with policy [%s]: [%s]", policyName, err)
	}

	
	err = id.SatisfiesPrincipal(principal)
	if err != nil {
		return fmt.Errorf("Failed verifying that proposal's creator satisfies local MSP principal during channelless check policy with policy [%s]: [%s]", policyName, err)
	}

	
	return id.Verify(signedProp.ProposalBytes, signedProp.Signature)
}



func (p *policyChecker) CheckPolicyBySignedData(channelID, policyName string, sd []*protoutil.SignedData) error {
	if channelID == "" {
		return errors.New("Invalid channel ID name during check policy on signed data. Name must be different from nil.")
	}

	if policyName == "" {
		return fmt.Errorf("Invalid policy name during check policy on signed data on channel [%s]. Name must be different from nil.", channelID)
	}

	if sd == nil {
		return fmt.Errorf("Invalid signed data during check policy on channel [%s] with policy [%s]", channelID, policyName)
	}

	
	policyManager := p.channelPolicyManagerGetter.Manager(channelID)
	if policyManager == nil {
		return fmt.Errorf("Failed to get policy manager for channel [%s]", channelID)
	}

	
	policy, _ := policyManager.GetPolicy(policyName)

	
	err := policy.Evaluate(sd)
	if err != nil {
		return fmt.Errorf("Failed evaluating policy on signed data during check policy on channel [%s] with policy [%s]: [%s]", channelID, policyName, err)
	}

	return nil
}
