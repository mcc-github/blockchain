/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package policies

import (
	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/protoutil"
)


func ImplicitMetaPolicyWithSubPolicy(subPolicyName string, rule cb.ImplicitMetaPolicy_Rule) *cb.ConfigPolicy {
	return &cb.ConfigPolicy{
		Policy: &cb.Policy{
			Type: int32(cb.Policy_IMPLICIT_META),
			Value: protoutil.MarshalOrPanic(&cb.ImplicitMetaPolicy{
				Rule:      rule,
				SubPolicy: subPolicyName,
			}),
		},
	}
}


func TemplateImplicitMetaPolicyWithSubPolicy(path []string, policyName string, subPolicyName string, rule cb.ImplicitMetaPolicy_Rule) *cb.ConfigGroup {
	root := protoutil.NewConfigGroup()
	group := root
	for _, element := range path {
		group.Groups[element] = protoutil.NewConfigGroup()
		group = group.Groups[element]
	}

	group.Policies[policyName] = ImplicitMetaPolicyWithSubPolicy(subPolicyName, rule)
	return root
}



func TemplateImplicitMetaPolicy(path []string, policyName string, rule cb.ImplicitMetaPolicy_Rule) *cb.ConfigGroup {
	return TemplateImplicitMetaPolicyWithSubPolicy(path, policyName, policyName, rule)
}


func TemplateImplicitMetaAnyPolicy(path []string, policyName string) *cb.ConfigGroup {
	return TemplateImplicitMetaPolicy(path, policyName, cb.ImplicitMetaPolicy_ANY)
}


func TemplateImplicitMetaAllPolicy(path []string, policyName string) *cb.ConfigGroup {
	return TemplateImplicitMetaPolicy(path, policyName, cb.ImplicitMetaPolicy_ALL)
}


func TemplateImplicitMetaMajorityPolicy(path []string, policyName string) *cb.ConfigGroup {
	return TemplateImplicitMetaPolicy(path, policyName, cb.ImplicitMetaPolicy_MAJORITY)
}
