/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package policies

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/protoutil"
	"go.uber.org/zap/zapcore"
)

type ImplicitMetaPolicy struct {
	Threshold   int
	SubPolicies []Policy

	
	managers      map[string]*ManagerImpl
	SubPolicyName string
}


func NewImplicitMetaPolicy(data []byte, managers map[string]*ManagerImpl) (*ImplicitMetaPolicy, error) {
	definition := &cb.ImplicitMetaPolicy{}
	if err := proto.Unmarshal(data, definition); err != nil {
		return nil, fmt.Errorf("Error unmarshaling to ImplicitMetaPolicy: %s", err)
	}

	subPolicies := make([]Policy, len(managers))

	i := 0
	for _, manager := range managers {
		subPolicies[i], _ = manager.GetPolicy(definition.SubPolicy)
		i++
	}

	var threshold int

	switch definition.Rule {
	case cb.ImplicitMetaPolicy_ANY:
		threshold = 1
	case cb.ImplicitMetaPolicy_ALL:
		threshold = len(subPolicies)
	case cb.ImplicitMetaPolicy_MAJORITY:
		threshold = len(subPolicies)/2 + 1
	}

	
	if len(subPolicies) == 0 {
		threshold = 0
	}

	return &ImplicitMetaPolicy{
		SubPolicies:   subPolicies,
		Threshold:     threshold,
		managers:      managers,
		SubPolicyName: definition.SubPolicy,
	}, nil
}


func (imp *ImplicitMetaPolicy) Evaluate(signatureSet []*protoutil.SignedData) error {
	logger.Debugf("This is an implicit meta policy, it will trigger other policy evaluations, whose failures may be benign")
	remaining := imp.Threshold

	defer func() {
		if remaining != 0 {
			
			if logger.IsEnabledFor(zapcore.DebugLevel) {
				var b bytes.Buffer
				b.WriteString(fmt.Sprintf("Evaluation Failed: Only %d policies were satisfied, but needed %d of [ ", imp.Threshold-remaining, imp.Threshold))
				for m := range imp.managers {
					b.WriteString(m)
					b.WriteString("/")
					b.WriteString(imp.SubPolicyName)
					b.WriteString(" ")
				}
				b.WriteString("]")
				logger.Debugf(b.String())
			}
		}
	}()

	for _, policy := range imp.SubPolicies {
		if policy.Evaluate(signatureSet) == nil {
			remaining--
			if remaining == 0 {
				return nil
			}
		}
	}
	if remaining == 0 {
		return nil
	}
	return fmt.Errorf("implicit policy evaluation failed - %d sub-policies were satisfied, but this policy requires %d of the '%s' sub-policies to be satisfied", (imp.Threshold - remaining), imp.Threshold, imp.SubPolicyName)
}
