/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package policies

import (
	"fmt"

	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)






func remap(sp *cb.SignaturePolicy, idRemap map[int]int) *cb.SignaturePolicy {
	switch t := sp.Type.(type) {
	case *cb.SignaturePolicy_NOutOf_:
		rules := []*cb.SignaturePolicy{}
		for _, rule := range t.NOutOf.Rules {
			
			
			rules = append(rules, remap(rule, idRemap))
		}

		return &cb.SignaturePolicy{
			Type: &cb.SignaturePolicy_NOutOf_{
				NOutOf: &cb.SignaturePolicy_NOutOf{
					N:     t.NOutOf.N,
					Rules: rules,
				},
			},
		}
	case *cb.SignaturePolicy_SignedBy:
		
		
		
		newID, in := idRemap[int(t.SignedBy)]
		if !in {
			panic("programming error")
		}

		return &cb.SignaturePolicy{
			Type: &cb.SignaturePolicy_SignedBy{
				SignedBy: int32(newID),
			},
		}
	default:
		panic(fmt.Sprintf("invalid policy type %T", t))
	}
}




func merge(this *cb.SignaturePolicyEnvelope, that *cb.SignaturePolicyEnvelope) {
	
	IDs := this.Identities
	idMap := map[string]int{}
	for i, id := range this.Identities {
		str := id.PrincipalClassification.String() + string(id.Principal)
		idMap[str] = i
	}

	
	
	
	
	
	
	idRemap := map[int]int{}
	for i, id := range that.Identities {
		str := id.PrincipalClassification.String() + string(id.Principal)
		if j, in := idMap[str]; in {
			idRemap[i] = j
		} else {
			idRemap[i] = len(IDs)
			idMap[str] = len(IDs)
			IDs = append(IDs, id)
		}
	}

	this.Identities = IDs

	newEntry := remap(that.Rule, idRemap)

	existingRules := this.Rule.Type.(*cb.SignaturePolicy_NOutOf_).NOutOf.Rules
	this.Rule.Type.(*cb.SignaturePolicy_NOutOf_).NOutOf.Rules = append(existingRules, newEntry)
}



func (p *ImplicitMetaPolicy) Convert() (*cb.SignaturePolicyEnvelope, error) {
	converted := &cb.SignaturePolicyEnvelope{
		Version: 0,
		Rule: &cb.SignaturePolicy{
			Type: &cb.SignaturePolicy_NOutOf_{
				NOutOf: &cb.SignaturePolicy_NOutOf{
					N: int32(p.Threshold),
				},
			},
		},
	}

	
	
	
	
	for i, subPolicy := range p.SubPolicies {
		convertibleSubpolicy, ok := subPolicy.(Converter)
		if !ok {
			return nil, errors.Errorf("subpolicy number %d type %T of policy %s is not convertible", i, subPolicy, p.SubPolicyName)
		}

		spe, err := convertibleSubpolicy.Convert()
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to convert subpolicy number %d of policy %s", i, p.SubPolicyName)
		}

		merge(converted, spe)
	}

	return converted, nil
}
