/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorsement

import (
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/protos/common"
	. "github.com/mcc-github/blockchain/protos/discovery"
	"github.com/pkg/errors"
)

func principalsFromCollectionConfig(ccp *common.CollectionConfigPackage) (principalSetsByCollectionName, error) {
	principalSetsByCollections := make(principalSetsByCollectionName)
	if ccp == nil {
		return principalSetsByCollections, nil
	}
	for _, colConfig := range ccp.Config {
		staticCol := colConfig.GetStaticCollectionConfig()
		if staticCol == nil {
			
			
			return nil, errors.Errorf("expected a static collection but got %v instead", colConfig)
		}
		if staticCol.MemberOrgsPolicy == nil {
			return nil, errors.Errorf("MemberOrgsPolicy of %s is nil", staticCol.Name)
		}
		pol := staticCol.MemberOrgsPolicy.GetSignaturePolicy()
		if pol == nil {
			return nil, errors.Errorf("policy of %s is nil", staticCol.Name)
		}
		var principals policies.PrincipalSet
		
		for _, principal := range pol.Identities {
			principals = append(principals, principal)
		}
		principalSetsByCollections[staticCol.Name] = principals
	}
	return principalSetsByCollections, nil
}

type principalSetsByCollectionName map[string]policies.PrincipalSet



func (psbc principalSetsByCollectionName) toIdentityFilter(channel string, evaluator principalEvaluator, cc *ChaincodeCall) (identityFilter, error) {
	var principalSets policies.PrincipalSets
	for _, col := range cc.CollectionNames {
		
		
		
		principalSet, exists := psbc[col]
		if !exists {
			return nil, errors.Errorf("collection %s doesn't exist in collection config for chaincode %s", col, cc.Name)
		}
		principalSets = append(principalSets, principalSet)
	}
	return filterForPrincipalSets(channel, evaluator, principalSets), nil
}


func filterForPrincipalSets(channel string, evaluator principalEvaluator, sets policies.PrincipalSets) identityFilter {
	return func(identity api.PeerIdentityType) bool {
		
		
		for _, principalSet := range sets {
			if !isIdentityAuthorizedByPrincipalSet(channel, evaluator, principalSet, identity) {
				return false
			}
		}
		return true
	}
}


func isIdentityAuthorizedByPrincipalSet(channel string, evaluator principalEvaluator, principalSet policies.PrincipalSet, identity api.PeerIdentityType) bool {
	
	
	for _, principal := range principalSet {
		err := evaluator.SatisfiesPrincipal(channel, identity, principal)
		if err != nil {
			continue
		}
		
		
		return true
	}
	return false
}
