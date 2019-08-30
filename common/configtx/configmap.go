/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/protoutil"
)

const (
	groupPrefix  = "[Group]  "
	valuePrefix  = "[Value]  "
	policyPrefix = "[Policy] "

	pathSeparator = "/"

	
	hackyFixOrdererCapabilities = "[Value]  /Channel/Orderer/Capabilities"
	hackyFixNewModPolicy        = "Admins"
)



func mapConfig(channelGroup *cb.ConfigGroup, rootGroupKey string) (map[string]comparable, error) {
	result := make(map[string]comparable)
	if channelGroup != nil {
		err := recurseConfig(result, []string{rootGroupKey}, channelGroup)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}


func addToMap(cg comparable, result map[string]comparable) error {
	var fqPath string

	switch {
	case cg.ConfigGroup != nil:
		fqPath = groupPrefix
	case cg.ConfigValue != nil:
		fqPath = valuePrefix
	case cg.ConfigPolicy != nil:
		fqPath = policyPrefix
	}

	if err := validateConfigID(cg.key); err != nil {
		return fmt.Errorf("Illegal characters in key: %s", fqPath)
	}

	if len(cg.path) == 0 {
		fqPath += pathSeparator + cg.key
	} else {
		fqPath += pathSeparator + strings.Join(cg.path, pathSeparator) + pathSeparator + cg.key
	}

	logger.Debugf("Adding to config map: %s", fqPath)

	result[fqPath] = cg

	return nil
}


func recurseConfig(result map[string]comparable, path []string, group *cb.ConfigGroup) error {
	if err := addToMap(comparable{key: path[len(path)-1], path: path[:len(path)-1], ConfigGroup: group}, result); err != nil {
		return err
	}

	for key, group := range group.Groups {
		nextPath := make([]string, len(path)+1)
		copy(nextPath, path)
		nextPath[len(nextPath)-1] = key
		if err := recurseConfig(result, nextPath, group); err != nil {
			return err
		}
	}

	for key, value := range group.Values {
		if err := addToMap(comparable{key: key, path: path, ConfigValue: value}, result); err != nil {
			return err
		}
	}

	for key, policy := range group.Policies {
		if err := addToMap(comparable{key: key, path: path, ConfigPolicy: policy}, result); err != nil {
			return err
		}
	}

	return nil
}



func configMapToConfig(configMap map[string]comparable, rootGroupKey string) (*cb.ConfigGroup, error) {
	rootPath := pathSeparator + rootGroupKey
	return recurseConfigMap(rootPath, configMap)
}



func recurseConfigMap(path string, configMap map[string]comparable) (*cb.ConfigGroup, error) {
	groupPath := groupPrefix + path
	group, ok := configMap[groupPath]
	if !ok {
		return nil, fmt.Errorf("Missing group at path: %s", groupPath)
	}

	if group.ConfigGroup == nil {
		return nil, fmt.Errorf("ConfigGroup not found at group path: %s", groupPath)
	}

	newConfigGroup := protoutil.NewConfigGroup()
	proto.Merge(newConfigGroup, group.ConfigGroup)

	for key := range group.Groups {
		updatedGroup, err := recurseConfigMap(path+pathSeparator+key, configMap)
		if err != nil {
			return nil, err
		}
		newConfigGroup.Groups[key] = updatedGroup
	}

	for key := range group.Values {
		valuePath := valuePrefix + path + pathSeparator + key
		value, ok := configMap[valuePath]
		if !ok {
			return nil, fmt.Errorf("Missing value at path: %s", valuePath)
		}
		if value.ConfigValue == nil {
			return nil, fmt.Errorf("ConfigValue not found at value path: %s", valuePath)
		}
		newConfigGroup.Values[key] = proto.Clone(value.ConfigValue).(*cb.ConfigValue)
	}

	for key := range group.Policies {
		policyPath := policyPrefix + path + pathSeparator + key
		policy, ok := configMap[policyPath]
		if !ok {
			return nil, fmt.Errorf("Missing policy at path: %s", policyPath)
		}
		if policy.ConfigPolicy == nil {
			return nil, fmt.Errorf("ConfigPolicy not found at policy path: %s", policyPath)
		}
		newConfigGroup.Policies[key] = proto.Clone(policy.ConfigPolicy).(*cb.ConfigPolicy)
		logger.Debugf("Setting policy for key %s to %+v", key, group.Policies[key])
	}

	
	
	
	
	
	
	
	if _, ok := configMap[hackyFixOrdererCapabilities]; ok {
		
		if newConfigGroup.ModPolicy == "" {
			logger.Debugf("Performing upgrade of group %s empty mod_policy", groupPath)
			newConfigGroup.ModPolicy = hackyFixNewModPolicy
		}

		for key, value := range newConfigGroup.Values {
			if value.ModPolicy == "" {
				logger.Debugf("Performing upgrade of value %s empty mod_policy", valuePrefix+path+pathSeparator+key)
				value.ModPolicy = hackyFixNewModPolicy
			}
		}

		for key, policy := range newConfigGroup.Policies {
			if policy.ModPolicy == "" {
				logger.Debugf("Performing upgrade of policy %s empty mod_policy", policyPrefix+path+pathSeparator+key)

				policy.ModPolicy = hackyFixNewModPolicy
			}
		}
	}

	return newConfigGroup, nil
}
