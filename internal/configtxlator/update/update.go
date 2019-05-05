/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package update

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
)

func computePoliciesMapUpdate(original, updated map[string]*cb.ConfigPolicy) (readSet, writeSet, sameSet map[string]*cb.ConfigPolicy, updatedMembers bool) {
	readSet = make(map[string]*cb.ConfigPolicy)
	writeSet = make(map[string]*cb.ConfigPolicy)

	
	
	sameSet = make(map[string]*cb.ConfigPolicy)

	for policyName, originalPolicy := range original {
		updatedPolicy, ok := updated[policyName]
		if !ok {
			updatedMembers = true
			continue
		}

		if originalPolicy.ModPolicy == updatedPolicy.ModPolicy && proto.Equal(originalPolicy.Policy, updatedPolicy.Policy) {
			sameSet[policyName] = &cb.ConfigPolicy{
				Version: originalPolicy.Version,
			}
			continue
		}

		writeSet[policyName] = &cb.ConfigPolicy{
			Version:   originalPolicy.Version + 1,
			ModPolicy: updatedPolicy.ModPolicy,
			Policy:    updatedPolicy.Policy,
		}
	}

	for policyName, updatedPolicy := range updated {
		if _, ok := original[policyName]; ok {
			
			continue
		}
		updatedMembers = true
		writeSet[policyName] = &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: updatedPolicy.ModPolicy,
			Policy:    updatedPolicy.Policy,
		}
	}

	return
}

func computeValuesMapUpdate(original, updated map[string]*cb.ConfigValue) (readSet, writeSet, sameSet map[string]*cb.ConfigValue, updatedMembers bool) {
	readSet = make(map[string]*cb.ConfigValue)
	writeSet = make(map[string]*cb.ConfigValue)

	
	
	sameSet = make(map[string]*cb.ConfigValue)

	for valueName, originalValue := range original {
		updatedValue, ok := updated[valueName]
		if !ok {
			updatedMembers = true
			continue
		}

		if originalValue.ModPolicy == updatedValue.ModPolicy && bytes.Equal(originalValue.Value, updatedValue.Value) {
			sameSet[valueName] = &cb.ConfigValue{
				Version: originalValue.Version,
			}
			continue
		}

		writeSet[valueName] = &cb.ConfigValue{
			Version:   originalValue.Version + 1,
			ModPolicy: updatedValue.ModPolicy,
			Value:     updatedValue.Value,
		}
	}

	for valueName, updatedValue := range updated {
		if _, ok := original[valueName]; ok {
			
			continue
		}
		updatedMembers = true
		writeSet[valueName] = &cb.ConfigValue{
			Version:   0,
			ModPolicy: updatedValue.ModPolicy,
			Value:     updatedValue.Value,
		}
	}

	return
}

func computeGroupsMapUpdate(original, updated map[string]*cb.ConfigGroup) (readSet, writeSet, sameSet map[string]*cb.ConfigGroup, updatedMembers bool) {
	readSet = make(map[string]*cb.ConfigGroup)
	writeSet = make(map[string]*cb.ConfigGroup)

	
	
	sameSet = make(map[string]*cb.ConfigGroup)

	for groupName, originalGroup := range original {
		updatedGroup, ok := updated[groupName]
		if !ok {
			updatedMembers = true
			continue
		}

		groupReadSet, groupWriteSet, groupUpdated := computeGroupUpdate(originalGroup, updatedGroup)
		if !groupUpdated {
			sameSet[groupName] = groupReadSet
			continue
		}

		readSet[groupName] = groupReadSet
		writeSet[groupName] = groupWriteSet

	}

	for groupName, updatedGroup := range updated {
		if _, ok := original[groupName]; ok {
			
			continue
		}
		updatedMembers = true
		_, groupWriteSet, _ := computeGroupUpdate(protoutil.NewConfigGroup(), updatedGroup)
		writeSet[groupName] = &cb.ConfigGroup{
			Version:   0,
			ModPolicy: updatedGroup.ModPolicy,
			Policies:  groupWriteSet.Policies,
			Values:    groupWriteSet.Values,
			Groups:    groupWriteSet.Groups,
		}
	}

	return
}

func computeGroupUpdate(original, updated *cb.ConfigGroup) (readSet, writeSet *cb.ConfigGroup, updatedGroup bool) {
	readSetPolicies, writeSetPolicies, sameSetPolicies, policiesMembersUpdated := computePoliciesMapUpdate(original.Policies, updated.Policies)
	readSetValues, writeSetValues, sameSetValues, valuesMembersUpdated := computeValuesMapUpdate(original.Values, updated.Values)
	readSetGroups, writeSetGroups, sameSetGroups, groupsMembersUpdated := computeGroupsMapUpdate(original.Groups, updated.Groups)

	
	if !(policiesMembersUpdated || valuesMembersUpdated || groupsMembersUpdated || original.ModPolicy != updated.ModPolicy) {

		
		if len(readSetPolicies) == 0 &&
			len(writeSetPolicies) == 0 &&
			len(readSetValues) == 0 &&
			len(writeSetValues) == 0 &&
			len(readSetGroups) == 0 &&
			len(writeSetGroups) == 0 {
			return &cb.ConfigGroup{
					Version: original.Version,
				}, &cb.ConfigGroup{
					Version: original.Version,
				}, false
		}

		return &cb.ConfigGroup{
				Version:  original.Version,
				Policies: readSetPolicies,
				Values:   readSetValues,
				Groups:   readSetGroups,
			}, &cb.ConfigGroup{
				Version:  original.Version,
				Policies: writeSetPolicies,
				Values:   writeSetValues,
				Groups:   writeSetGroups,
			}, true
	}

	for k, samePolicy := range sameSetPolicies {
		readSetPolicies[k] = samePolicy
		writeSetPolicies[k] = samePolicy
	}

	for k, sameValue := range sameSetValues {
		readSetValues[k] = sameValue
		writeSetValues[k] = sameValue
	}

	for k, sameGroup := range sameSetGroups {
		readSetGroups[k] = sameGroup
		writeSetGroups[k] = sameGroup
	}

	return &cb.ConfigGroup{
			Version:  original.Version,
			Policies: readSetPolicies,
			Values:   readSetValues,
			Groups:   readSetGroups,
		}, &cb.ConfigGroup{
			Version:   original.Version + 1,
			Policies:  writeSetPolicies,
			Values:    writeSetValues,
			Groups:    writeSetGroups,
			ModPolicy: updated.ModPolicy,
		}, true
}

func Compute(original, updated *cb.Config) (*cb.ConfigUpdate, error) {
	if original.ChannelGroup == nil {
		return nil, fmt.Errorf("no channel group included for original config")
	}

	if updated.ChannelGroup == nil {
		return nil, fmt.Errorf("no channel group included for updated config")
	}

	readSet, writeSet, groupUpdated := computeGroupUpdate(original.ChannelGroup, updated.ChannelGroup)
	if !groupUpdated {
		return nil, fmt.Errorf("no differences detected between original and updated config")
	}
	return &cb.ConfigUpdate{
		ReadSet:  readSet,
		WriteSet: writeSet,
	}, nil
}