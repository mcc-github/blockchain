/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package policy

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)




type Policy interface {
	policies.Policy
}




type ChannelPolicyManagerGetter interface {
	policies.ChannelPolicyManagerGetter
}




type PolicyManager interface {
	policies.Manager
}




type IdentityDeserializer interface {
	msp.IdentityDeserializer
}




type Identity interface {
	msp.Identity
}





type SignaturePolicyProvider interface {
	
	NewPolicy(signaturePolicy *common.SignaturePolicyEnvelope) (policies.Policy, error)
}


type ChannelPolicyReferenceProvider interface {
	
	NewPolicy(channelConfigPolicyReference string) (policies.Policy, error)
}

type ApplicationPolicyEvaluator struct {
	signaturePolicyProvider        SignaturePolicyProvider
	channelPolicyReferenceProvider ChannelPolicyReferenceProvider
}


type Manager interface {
	
	GetPolicy(id string) (policies.Policy, bool)
}

type ChannelPolicyReferenceProviderImpl struct {
	Manager
}

func (c *ChannelPolicyReferenceProviderImpl) NewPolicy(channelConfigPolicyReference string) (policies.Policy, error) {
	p, ok := c.GetPolicy(channelConfigPolicyReference)
	if !ok {
		return nil, errors.Errorf("failed to retrieve policy for reference %s", channelConfigPolicyReference)
	}

	return p, nil
}



type dynamicPolicyManager struct {
	channelPolicyManagerGetter policies.ChannelPolicyManagerGetter
	channelID                  string
}

func (d *dynamicPolicyManager) GetPolicy(id string) (policies.Policy, bool) {
	mgr := d.channelPolicyManagerGetter.Manager(d.channelID)
	if mgr == nil {
		
		
		
		
		panic("programming error")
	}

	return mgr.GetPolicy(id)
}


func New(deserializer msp.IdentityDeserializer, channel string, channelPolicyManagerGetter policies.ChannelPolicyManagerGetter) (*ApplicationPolicyEvaluator, error) {
	mgr := channelPolicyManagerGetter.Manager(channel)
	if mgr == nil {
		return nil, errors.Errorf("failed to retrieve policy manager for channel %s", channel)
	}

	return &ApplicationPolicyEvaluator{
		signaturePolicyProvider: &cauthdsl.EnvelopeBasedPolicyProvider{Deserializer: deserializer},
		channelPolicyReferenceProvider: &ChannelPolicyReferenceProviderImpl{Manager: &dynamicPolicyManager{
			channelID:                  channel,
			channelPolicyManagerGetter: channelPolicyManagerGetter,
		}},
	}, nil
}

func (a *ApplicationPolicyEvaluator) evaluateSignaturePolicy(signaturePolicy *common.SignaturePolicyEnvelope, signatureSet []*protoutil.SignedData) error {
	p, err := a.signaturePolicyProvider.NewPolicy(signaturePolicy)
	if err != nil {
		return errors.WithMessage(err, "could not create evaluator for signature policy")
	}

	return p.Evaluate(signatureSet)
}

func (a *ApplicationPolicyEvaluator) evaluateChannelConfigPolicyReference(channelConfigPolicyReference string, signatureSet []*protoutil.SignedData) error {
	p, err := a.channelPolicyReferenceProvider.NewPolicy(channelConfigPolicyReference)
	if err != nil {
		return errors.WithMessage(err, "could not create evaluator for channel reference policy")
	}

	return p.Evaluate(signatureSet)
}

func (a *ApplicationPolicyEvaluator) Evaluate(policyBytes []byte, signatureSet []*protoutil.SignedData) error {
	p := &peer.ApplicationPolicy{}
	err := proto.Unmarshal(policyBytes, p)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal ApplicationPolicy bytes")
	}

	switch policy := p.Type.(type) {
	case *peer.ApplicationPolicy_SignaturePolicy:
		return a.evaluateSignaturePolicy(policy.SignaturePolicy, signatureSet)
	case *peer.ApplicationPolicy_ChannelConfigPolicyReference:
		return a.evaluateChannelConfigPolicyReference(policy.ChannelConfigPolicyReference, signatureSet)
	default:
		return errors.Errorf("unsupported policy type %T", policy)
	}
}
