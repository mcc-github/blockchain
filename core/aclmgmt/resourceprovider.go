/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package aclmgmt

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
)




type PolicyNotFound string

func (e PolicyNotFound) Error() string {
	return fmt.Sprintf("policy %s not found", string(e))
}


type InvalidIdInfo string

func (e InvalidIdInfo) Error() string {
	return fmt.Sprintf("Invalid id for policy [%s]", string(e))
}




type policyEvaluator interface {
	PolicyRefForAPI(resName string) string
	Evaluate(polName string, id []*common.SignedData) error
}


type policyEvaluatorImpl struct {
	bundle channelconfig.Resources
}

func (pe *policyEvaluatorImpl) PolicyRefForAPI(resName string) string {
	app, exists := pe.bundle.ApplicationConfig()
	if !exists {
		return ""
	}

	pm := app.APIPolicyMapper()
	if pm == nil {
		return ""
	}

	return pm.PolicyRefForAPI(resName)
}

func (pe *policyEvaluatorImpl) Evaluate(polName string, sd []*common.SignedData) error {
	policy, ok := pe.bundle.PolicyManager().GetPolicy(polName)
	if !ok {
		return PolicyNotFound(polName)
	}

	return policy.Evaluate(sd)
}




type aclmgmtPolicyProvider interface {
	
	GetPolicyName(resName string) string

	
	CheckACL(polName string, idinfo interface{}) error
}


type aclmgmtPolicyProviderImpl struct {
	pEvaluator policyEvaluator
}


func (rp *aclmgmtPolicyProviderImpl) GetPolicyName(resName string) string {
	return rp.pEvaluator.PolicyRefForAPI(resName)
}



func (rp *aclmgmtPolicyProviderImpl) CheckACL(polName string, idinfo interface{}) error {
	aclLogger.Debugf("acl check(%s)", polName)

	
	var sd []*common.SignedData
	var err error
	switch idinfo.(type) {
	case *pb.SignedProposal:
		signedProp, _ := idinfo.(*pb.SignedProposal)
		
		proposal, err := utils.GetProposal(signedProp.ProposalBytes)
		if err != nil {
			return fmt.Errorf("Failing extracting proposal during check policy with policy [%s]: [%s]", polName, err)
		}

		header, err := utils.GetHeader(proposal.Header)
		if err != nil {
			return fmt.Errorf("Failing extracting header during check policy [%s]: [%s]", polName, err)
		}

		shdr, err := utils.GetSignatureHeader(header.SignatureHeader)
		if err != nil {
			return fmt.Errorf("Invalid Proposal's SignatureHeader during check policy [%s]: [%s]", polName, err)
		}

		sd = []*common.SignedData{{
			Data:      signedProp.ProposalBytes,
			Identity:  shdr.Creator,
			Signature: signedProp.Signature,
		}}
	case *common.Envelope:
		sd, err = idinfo.(*common.Envelope).AsSignedData()
		if err != nil {
			return err
		}
	default:
		return InvalidIdInfo(polName)
	}

	err = rp.pEvaluator.Evaluate(polName, sd)
	if err != nil {
		return fmt.Errorf("failed evaluating policy on signed data during check policy [%s]: [%s]", polName, err)
	}

	return nil
}




type ResourceGetter func(channelID string) channelconfig.Resources


type resourceProvider struct {
	
	resGetter ResourceGetter

	
	defaultProvider ACLProvider
}


func newResourceProvider(rg ResourceGetter, defprov ACLProvider) *resourceProvider {
	return &resourceProvider{rg, defprov}
}


func (rp *resourceProvider) CheckACL(resName string, channelID string, idinfo interface{}) error {
	resCfg := rp.resGetter(channelID)

	if resCfg != nil {
		pp := &aclmgmtPolicyProviderImpl{&policyEvaluatorImpl{resCfg}}
		policyName := pp.GetPolicyName(resName)
		if policyName != "" {
			aclLogger.Debugf("acl policy %s found in config for resource %s", policyName, resName)
			return pp.CheckACL(policyName, idinfo)
		}
		aclLogger.Debugf("acl policy not found in config for resource %s", resName)
	}

	return rp.defaultProvider.CheckACL(resName, channelID, idinfo)
}
