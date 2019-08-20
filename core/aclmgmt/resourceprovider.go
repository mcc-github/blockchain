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
	"github.com/mcc-github/blockchain/protoutil"
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
	Evaluate(polName string, id []*protoutil.SignedData) error
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

func (pe *policyEvaluatorImpl) Evaluate(polName string, sd []*protoutil.SignedData) error {
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

	
	var sd []*protoutil.SignedData
	switch idinfo := idinfo.(type) {
	case *pb.SignedProposal:
		signedProp := idinfo
		proposal, err := protoutil.UnmarshalProposal(signedProp.ProposalBytes)
		if err != nil {
			return fmt.Errorf("Failing extracting proposal during check policy with policy [%s]: [%s]", polName, err)
		}

		header, err := protoutil.UnmarshalHeader(proposal.Header)
		if err != nil {
			return fmt.Errorf("Failing extracting header during check policy [%s]: [%s]", polName, err)
		}

		shdr, err := protoutil.UnmarshalSignatureHeader(header.SignatureHeader)
		if err != nil {
			return fmt.Errorf("Invalid Proposal's SignatureHeader during check policy [%s]: [%s]", polName, err)
		}

		sd = []*protoutil.SignedData{{
			Data:      signedProp.ProposalBytes,
			Identity:  shdr.Creator,
			Signature: signedProp.Signature,
		}}

	case *common.Envelope:
		var err error
		sd, err = protoutil.EnvelopeAsSignedData(idinfo)
		if err != nil {
			return err
		}

	default:
		return InvalidIdInfo(polName)
	}

	err := rp.pEvaluator.Evaluate(polName, sd)
	if err != nil {
		return fmt.Errorf("failed evaluating policy on signed data during check policy [%s]: [%s]", polName, err)
	}

	return nil
}




type ResourceGetter func(channelID string) channelconfig.Resources


type resourceProvider struct {
	
	resGetter ResourceGetter

	
	defaultProvider defaultACLProvider
}


func newResourceProvider(rg ResourceGetter, defprov defaultACLProvider) *resourceProvider {
	return &resourceProvider{rg, defprov}
}

func (rp *resourceProvider) enforceDefaultBehavior(resName string, channelID string, idinfo interface{}) bool {
	
	
	return rp.defaultProvider.IsPtypePolicy(resName)
}


func (rp *resourceProvider) CheckACL(resName string, channelID string, idinfo interface{}) error {
	if !rp.enforceDefaultBehavior(resName, channelID, idinfo) {
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
	}

	return rp.defaultProvider.CheckACL(resName, channelID, idinfo)
}
