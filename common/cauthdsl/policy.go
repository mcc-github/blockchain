/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cauthdsl

import (
	"errors"
	"fmt"

	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/msp"
)

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
		evaluator:    compiled,
		deserializer: pr.deserializer,
	}, sigPolicy, nil

}

type policy struct {
	evaluator    func([]*cb.SignedData, []bool) bool
	deserializer msp.IdentityDeserializer
}


func (p *policy) Evaluate(signatureSet []*cb.SignedData) error {
	if p == nil {
		return fmt.Errorf("No such policy")
	}

	ok := p.evaluator(deduplicate(signatureSet, p.deserializer), make([]bool, len(signatureSet)))
	if !ok {
		return errors.New("signature set did not satisfy policy")
	}
	return nil
}
