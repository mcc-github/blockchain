/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cauthdsl

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	mspp "github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
)

type Identity interface {
	
	
	
	
	SatisfiesPrincipal(principal *mspp.MSPPrincipal) error

	
	GetIdentifier() *msp.IdentityIdentifier
}

type IdentityAndSignature interface {
	
	Identity() (Identity, error)

	
	Verify() error
}

type deserializeAndVerify struct {
	signedData           *protoutil.SignedData
	deserializer         msp.IdentityDeserializer
	deserializedIdentity msp.Identity
}

func (d *deserializeAndVerify) Identity() (Identity, error) {
	deserializedIdentity, err := d.deserializer.DeserializeIdentity(d.signedData.Identity)
	if err != nil {
		return nil, err
	}

	d.deserializedIdentity = deserializedIdentity
	return deserializedIdentity, nil
}

func (d *deserializeAndVerify) Verify() error {
	if d.deserializedIdentity == nil {
		cauthdslLogger.Panicf("programming error, Identity must be called prior to Verify")
	}
	return d.deserializedIdentity.Verify(d.signedData.Data, d.signedData.Signature)
}

type provider struct {
	deserializer msp.IdentityDeserializer
}


func NewPolicyProvider(deserializer msp.IdentityDeserializer) policies.Provider {
	return &provider{
		deserializer: deserializer,
	}
}


func (pr *provider) NewPolicy(data []byte) (policies.Policy, proto.Message, error) {
	sigPolicy := &cb.SignaturePolicyEnvelope{}
	if err := proto.Unmarshal(data, sigPolicy); err != nil {
		return nil, nil, fmt.Errorf("Error unmarshaling to SignaturePolicy: %s", err)
	}

	if sigPolicy.Version != 0 {
		return nil, nil, fmt.Errorf("This evaluator only understands messages of version 0, but version was %d", sigPolicy.Version)
	}

	compiled, err := compile(sigPolicy.Rule, sigPolicy.Identities, pr.deserializer)
	if err != nil {
		return nil, nil, err
	}

	return &policy{
		evaluator:               compiled,
		deserializer:            pr.deserializer,
		signaturePolicyEnvelope: sigPolicy,
	}, sigPolicy, nil

}


type EnvelopeBasedPolicyProvider struct {
	Deserializer msp.IdentityDeserializer
}


func (pp *EnvelopeBasedPolicyProvider) NewPolicy(sigPolicy *cb.SignaturePolicyEnvelope) (policies.Policy, error) {
	if sigPolicy == nil {
		return nil, errors.New("invalid arguments")
	}

	compiled, err := compile(sigPolicy.Rule, sigPolicy.Identities, pp.Deserializer)
	if err != nil {
		return nil, err
	}

	return &policy{
		evaluator:               compiled,
		deserializer:            pp.Deserializer,
		signaturePolicyEnvelope: sigPolicy,
	}, nil
}

type policy struct {
	signaturePolicyEnvelope *cb.SignaturePolicyEnvelope
	evaluator               func([]IdentityAndSignature, []bool) bool
	deserializer            msp.IdentityDeserializer
}


func (p *policy) Evaluate(signatureSet []*protoutil.SignedData) error {
	if p == nil {
		return fmt.Errorf("No such policy")
	}
	idAndS := make([]IdentityAndSignature, len(signatureSet))
	for i, sd := range signatureSet {
		idAndS[i] = &deserializeAndVerify{
			signedData:   sd,
			deserializer: p.deserializer,
		}
	}

	ok := p.evaluator(deduplicate(idAndS), make([]bool, len(signatureSet)))
	if !ok {
		return errors.New("signature set did not satisfy policy")
	}
	return nil
}

func (p *policy) Convert() (*cb.SignaturePolicyEnvelope, error) {
	if p.signaturePolicyEnvelope == nil {
		return nil, errors.New("nil policy field")
	}

	return p.signaturePolicyEnvelope, nil
}
