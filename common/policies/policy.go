/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package policies

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

const (
	
	PathSeparator = "/"

	
	ChannelPrefix = "Channel"

	
	ApplicationPrefix = "Application"

	
	OrdererPrefix = "Orderer"

	
	ChannelReaders = PathSeparator + ChannelPrefix + PathSeparator + "Readers"

	
	ChannelWriters = PathSeparator + ChannelPrefix + PathSeparator + "Writers"

	
	ChannelApplicationReaders = PathSeparator + ChannelPrefix + PathSeparator + ApplicationPrefix + PathSeparator + "Readers"

	
	ChannelApplicationWriters = PathSeparator + ChannelPrefix + PathSeparator + ApplicationPrefix + PathSeparator + "Writers"

	
	ChannelApplicationAdmins = PathSeparator + ChannelPrefix + PathSeparator + ApplicationPrefix + PathSeparator + "Admins"

	
	BlockValidation = PathSeparator + ChannelPrefix + PathSeparator + OrdererPrefix + PathSeparator + "BlockValidation"

	
	ChannelOrdererAdmins = PathSeparator + ChannelPrefix + PathSeparator + OrdererPrefix + PathSeparator + "Admins"

	
	ChannelOrdererWriters = PathSeparator + ChannelPrefix + PathSeparator + OrdererPrefix + PathSeparator + "Writers"

	
	ChannelOrdererReaders = PathSeparator + ChannelPrefix + PathSeparator + OrdererPrefix + PathSeparator + "Readers"
)

var logger = flogging.MustGetLogger("policies")


type PrincipalSet []*msp.MSPPrincipal


type PrincipalSets []PrincipalSet


func (psSets PrincipalSets) ContainingOnly(f func(*msp.MSPPrincipal) bool) PrincipalSets {
	var res PrincipalSets
	for _, set := range psSets {
		if !set.ContainingOnly(f) {
			continue
		}
		res = append(res, set)
	}
	return res
}



func (ps PrincipalSet) ContainingOnly(f func(*msp.MSPPrincipal) bool) bool {
	for _, principal := range ps {
		if !f(principal) {
			return false
		}
	}
	return true
}


func (ps PrincipalSet) UniqueSet() map[*msp.MSPPrincipal]int {
	
	histogram := make(map[struct {
		cls       int32
		principal string
	}]int)
	
	for _, principal := range ps {
		key := struct {
			cls       int32
			principal string
		}{
			cls:       int32(principal.PrincipalClassification),
			principal: string(principal.Principal),
		}
		histogram[key]++
	}
	
	res := make(map[*msp.MSPPrincipal]int)
	for principal, count := range histogram {
		res[&msp.MSPPrincipal{
			PrincipalClassification: msp.MSPPrincipal_Classification(principal.cls),
			Principal:               []byte(principal.principal),
		}] = count
	}
	return res
}


type Policy interface {
	
	Evaluate(signatureSet []*protoutil.SignedData) error
}


type InquireablePolicy interface {
	
	
	SatisfiedBy() []PrincipalSet
}


type Manager interface {
	
	GetPolicy(id string) (Policy, bool)

	
	Manager(path []string) (Manager, bool)
}


type Provider interface {
	
	NewPolicy(data []byte) (Policy, proto.Message, error)
}



type ChannelPolicyManagerGetter interface {
	
	
	Manager(channelID string) (Manager, bool)
}



type ManagerImpl struct {
	path     string 
	policies map[string]Policy
	managers map[string]*ManagerImpl
}


func NewManagerImpl(path string, providers map[int32]Provider, root *cb.ConfigGroup) (*ManagerImpl, error) {
	var err error
	_, ok := providers[int32(cb.Policy_IMPLICIT_META)]
	if ok {
		logger.Panicf("ImplicitMetaPolicy type must be provider by the policy manager")
	}

	managers := make(map[string]*ManagerImpl)

	for groupName, group := range root.Groups {
		managers[groupName], err = NewManagerImpl(path+PathSeparator+groupName, providers, group)
		if err != nil {
			return nil, err
		}
	}

	policies := make(map[string]Policy)
	for policyName, configPolicy := range root.Policies {
		policy := configPolicy.Policy
		if policy == nil {
			return nil, fmt.Errorf("policy %s at path %s was nil", policyName, path)
		}

		var cPolicy Policy

		if policy.Type == int32(cb.Policy_IMPLICIT_META) {
			imp, err := newImplicitMetaPolicy(policy.Value, managers)
			if err != nil {
				return nil, errors.Wrapf(err, "implicit policy %s at path %s did not compile", policyName, path)
			}
			cPolicy = imp
		} else {
			provider, ok := providers[int32(policy.Type)]
			if !ok {
				return nil, fmt.Errorf("policy %s at path %s has unknown policy type: %v", policyName, path, policy.Type)
			}

			var err error
			cPolicy, _, err = provider.NewPolicy(policy.Value)
			if err != nil {
				return nil, errors.Wrapf(err, "policy %s at path %s did not compile", policyName, path)
			}
		}

		policies[policyName] = cPolicy

		logger.Debugf("Proposed new policy %s for %s", policyName, path)
	}

	for groupName, manager := range managers {
		for policyName, policy := range manager.policies {
			policies[groupName+PathSeparator+policyName] = policy
		}
	}

	return &ManagerImpl{
		path:     path,
		policies: policies,
		managers: managers,
	}, nil
}

type rejectPolicy string

func (rp rejectPolicy) Evaluate(signedData []*protoutil.SignedData) error {
	return fmt.Errorf("No such policy: '%s'", rp)
}


func (pm *ManagerImpl) Manager(path []string) (Manager, bool) {
	logger.Debugf("Manager %s looking up path %v", pm.path, path)
	for manager := range pm.managers {
		logger.Debugf("Manager %s has managers %s", pm.path, manager)
	}
	if len(path) == 0 {
		return pm, true
	}

	m, ok := pm.managers[path[0]]
	if !ok {
		return nil, false
	}

	return m.Manager(path[1:])
}

type policyLogger struct {
	policy     Policy
	policyName string
}

func (pl *policyLogger) Evaluate(signatureSet []*protoutil.SignedData) error {
	if logger.IsEnabledFor(zapcore.DebugLevel) {
		logger.Debugf("== Evaluating %T Policy %s ==", pl.policy, pl.policyName)
		defer logger.Debugf("== Done Evaluating %T Policy %s", pl.policy, pl.policyName)
	}

	err := pl.policy.Evaluate(signatureSet)
	if err != nil {
		logger.Debugf("Signature set did not satisfy policy %s", pl.policyName)
	} else {
		logger.Debugf("Signature set satisfies policy %s", pl.policyName)
	}
	return err
}


func (pm *ManagerImpl) GetPolicy(id string) (Policy, bool) {
	if id == "" {
		logger.Errorf("Returning dummy reject all policy because no policy ID supplied")
		return rejectPolicy(id), false
	}
	var relpath string

	if strings.HasPrefix(id, PathSeparator) {
		if !strings.HasPrefix(id, PathSeparator+pm.path) {
			logger.Debugf("Requested absolute policy %s from %s, returning rejectAll", id, pm.path)
			return rejectPolicy(id), false
		}
		
		relpath = id[1+len(pm.path)+1:]
	} else {
		relpath = id
	}

	policy, ok := pm.policies[relpath]
	if !ok {
		logger.Debugf("Returning dummy reject all policy because %s could not be found in %s/%s", id, pm.path, relpath)
		return rejectPolicy(relpath), false
	}

	return &policyLogger{
		policy:     policy,
		policyName: PathSeparator + pm.path + PathSeparator + relpath,
	}, true
}
