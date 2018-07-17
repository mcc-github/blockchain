/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	"bytes"

	cb "github.com/mcc-github/blockchain/protos/common"
)

type comparable struct {
	*cb.ConfigGroup
	*cb.ConfigValue
	*cb.ConfigPolicy
	key  string
	path []string
}

func (cg comparable) equals(other comparable) bool {
	switch {
	case cg.ConfigGroup != nil:
		if other.ConfigGroup == nil {
			return false
		}
		return equalConfigGroup(cg.ConfigGroup, other.ConfigGroup)
	case cg.ConfigValue != nil:
		if other.ConfigValue == nil {
			return false
		}
		return equalConfigValues(cg.ConfigValue, other.ConfigValue)
	case cg.ConfigPolicy != nil:
		if other.ConfigPolicy == nil {
			return false
		}
		return equalConfigPolicies(cg.ConfigPolicy, other.ConfigPolicy)
	}

	
	return false
}

func (cg comparable) version() uint64 {
	switch {
	case cg.ConfigGroup != nil:
		return cg.ConfigGroup.Version
	case cg.ConfigValue != nil:
		return cg.ConfigValue.Version
	case cg.ConfigPolicy != nil:
		return cg.ConfigPolicy.Version
	}

	
	return 0
}

func (cg comparable) modPolicy() string {
	switch {
	case cg.ConfigGroup != nil:
		return cg.ConfigGroup.ModPolicy
	case cg.ConfigValue != nil:
		return cg.ConfigValue.ModPolicy
	case cg.ConfigPolicy != nil:
		return cg.ConfigPolicy.ModPolicy
	}

	
	return ""
}

func equalConfigValues(lhs, rhs *cb.ConfigValue) bool {
	return lhs.Version == rhs.Version &&
		lhs.ModPolicy == rhs.ModPolicy &&
		bytes.Equal(lhs.Value, rhs.Value)
}

func equalConfigPolicies(lhs, rhs *cb.ConfigPolicy) bool {
	if lhs.Version != rhs.Version ||
		lhs.ModPolicy != rhs.ModPolicy {
		return false
	}

	if lhs.Policy == nil || rhs.Policy == nil {
		return lhs.Policy == rhs.Policy
	}

	return lhs.Policy.Type == rhs.Policy.Type &&
		bytes.Equal(lhs.Policy.Value, rhs.Policy.Value)
}




func subsetOfGroups(inner, outer map[string]*cb.ConfigGroup) bool {
	
	if len(inner) == 0 {
		return true
	}

	
	if len(inner) > len(outer) {
		return false
	}

	
	for key := range inner {
		if _, ok := outer[key]; !ok {
			return false
		}
	}

	return true
}

func subsetOfPolicies(inner, outer map[string]*cb.ConfigPolicy) bool {
	
	if len(inner) == 0 {
		return true
	}

	
	if len(inner) > len(outer) {
		return false
	}

	
	for key := range inner {
		if _, ok := outer[key]; !ok {
			return false
		}
	}

	return true
}

func subsetOfValues(inner, outer map[string]*cb.ConfigValue) bool {
	
	if len(inner) == 0 {
		return true
	}

	
	if len(inner) > len(outer) {
		return false
	}

	
	for key := range inner {
		if _, ok := outer[key]; !ok {
			return false
		}
	}

	return true
}

func equalConfigGroup(lhs, rhs *cb.ConfigGroup) bool {
	if lhs.Version != rhs.Version ||
		lhs.ModPolicy != rhs.ModPolicy {
		return false
	}

	if !subsetOfGroups(lhs.Groups, rhs.Groups) ||
		!subsetOfGroups(rhs.Groups, lhs.Groups) ||
		!subsetOfPolicies(lhs.Policies, rhs.Policies) ||
		!subsetOfPolicies(rhs.Policies, lhs.Policies) ||
		!subsetOfValues(lhs.Values, rhs.Values) ||
		!subsetOfValues(rhs.Values, lhs.Values) {
		return false
	}

	return true
}
