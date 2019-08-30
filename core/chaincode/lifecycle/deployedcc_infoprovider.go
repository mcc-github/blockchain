/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"fmt"
	"regexp"
	"strings"

	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset/kvrwset"
	"github.com/mcc-github/blockchain-protos-go/msp"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/util"
	validationState "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protoutil"

	"github.com/pkg/errors"
)


type LegacyDeployedCCInfoProvider interface {
	ledger.DeployedChaincodeInfoProvider
}

const (
	LifecycleEndorsementPolicyRef = "/Channel/Application/LifecycleEndorsement"
)

var (
	
	LifecycleDefaultEndorsementPolicyBytes = protoutil.MarshalOrPanic(&pb.ApplicationPolicy{
		Type: &pb.ApplicationPolicy_ChannelConfigPolicyReference{
			ChannelConfigPolicyReference: LifecycleEndorsementPolicyRef,
		},
	})
)

type ValidatorCommitter struct {
	Resources                    *Resources
	LegacyDeployedCCInfoProvider LegacyDeployedCCInfoProvider
}


func (vc *ValidatorCommitter) Namespaces() []string {
	return append([]string{LifecycleNamespace}, vc.LegacyDeployedCCInfoProvider.Namespaces()...)
}

var SequenceMatcher = regexp.MustCompile("^" + NamespacesName + "/fields/([^/]+)/Sequence$")


func (vc *ValidatorCommitter) UpdatedChaincodes(stateUpdates map[string][]*kvrwset.KVWrite) ([]*ledger.ChaincodeLifecycleInfo, error) {
	lifecycleInfo := []*ledger.ChaincodeLifecycleInfo{}

	
	lifecycleUpdates := stateUpdates[LifecycleNamespace]

	for _, kvWrite := range lifecycleUpdates {
		matches := SequenceMatcher.FindStringSubmatch(kvWrite.Key)
		if len(matches) != 2 {
			continue
		}
		
		lifecycleInfo = append(lifecycleInfo, &ledger.ChaincodeLifecycleInfo{Name: matches[1]})
	}

	legacyUpdates, err := vc.LegacyDeployedCCInfoProvider.UpdatedChaincodes(stateUpdates)
	if err != nil {
		return nil, errors.WithMessage(err, "error invoking legacy deployed cc info provider")
	}

	return append(lifecycleInfo, legacyUpdates...), nil
}

func (vc *ValidatorCommitter) ChaincodeInfo(channelName, chaincodeName string, qe ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error) {
	exists, definedChaincode, err := vc.Resources.ChaincodeDefinitionIfDefined(chaincodeName, &SimpleQueryExecutorShim{
		Namespace:           LifecycleNamespace,
		SimpleQueryExecutor: qe,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "could not get info about chaincode")
	}

	if !exists {
		return vc.LegacyDeployedCCInfoProvider.ChaincodeInfo(channelName, chaincodeName, qe)
	}

	ic, err := vc.ChaincodeImplicitCollections(channelName)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create implicit collections for channel")
	}

	return &ledger.DeployedChaincodeInfo{
		Name:                        chaincodeName,
		Version:                     definedChaincode.EndorsementInfo.Version,
		Hash:                        util.ComputeSHA256([]byte(chaincodeName + ":" + definedChaincode.EndorsementInfo.Version)),
		ExplicitCollectionConfigPkg: definedChaincode.Collections,
		ImplicitCollections:         ic,
		IsLegacy:                    false,
	}, nil
}

var ImplicitCollectionMatcher = regexp.MustCompile("^" + ImplicitCollectionNameForOrg("(.+)") + "$")



func (vc *ValidatorCommitter) CollectionInfo(channelName, chaincodeName, collectionName string, qe ledger.SimpleQueryExecutor) (*cb.StaticCollectionConfig, error) {
	exists, definedChaincode, err := vc.Resources.ChaincodeDefinitionIfDefined(chaincodeName, &SimpleQueryExecutorShim{
		Namespace:           LifecycleNamespace,
		SimpleQueryExecutor: qe,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "could not get chaincode")
	}

	if !exists {
		return vc.LegacyDeployedCCInfoProvider.CollectionInfo(channelName, chaincodeName, collectionName, qe)
	}

	matches := ImplicitCollectionMatcher.FindStringSubmatch(collectionName)
	if len(matches) == 2 {
		return GenerateImplicitCollectionForOrg(matches[1]), nil
	}

	if definedChaincode.Collections != nil {
		for _, conf := range definedChaincode.Collections.Config {
			staticCollConfig := conf.GetStaticCollectionConfig()
			if staticCollConfig != nil && staticCollConfig.Name == collectionName {
				return staticCollConfig, nil
			}
		}
	}
	return nil, nil
}



func (vc *ValidatorCommitter) ImplicitCollections(channelName, chaincodeName string, qe ledger.SimpleQueryExecutor) ([]*cb.StaticCollectionConfig, error) {
	exists, _, err := vc.Resources.ChaincodeDefinitionIfDefined(chaincodeName, &SimpleQueryExecutorShim{
		Namespace:           LifecycleNamespace,
		SimpleQueryExecutor: qe,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "could not get info about chaincode")
	}

	if !exists {
		
		return nil, nil
	}

	return vc.ChaincodeImplicitCollections(channelName)
}


func (vc *ValidatorCommitter) ChaincodeImplicitCollections(channelName string) ([]*cb.StaticCollectionConfig, error) {
	channelConfig := vc.Resources.ChannelConfigSource.GetStableChannelConfig(channelName)
	if channelConfig == nil {
		return nil, errors.Errorf("could not get channelconfig for channel %s", channelName)
	}
	ac, ok := channelConfig.ApplicationConfig()
	if !ok {
		return nil, errors.Errorf("could not get application config for channel %s", channelName)
	}

	orgs := ac.Organizations()
	implicitCollections := make([]*cb.StaticCollectionConfig, 0, len(orgs))
	for _, org := range orgs {
		implicitCollections = append(implicitCollections, GenerateImplicitCollectionForOrg(org.MSPID()))
	}

	return implicitCollections, nil
}

func GenerateImplicitCollectionForOrg(mspid string) *cb.StaticCollectionConfig {
	return &cb.StaticCollectionConfig{
		Name: ImplicitCollectionNameForOrg(mspid),
		MemberOrgsPolicy: &cb.CollectionPolicyConfig{
			Payload: &cb.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: cauthdsl.SignedByMspMember(mspid),
			},
		},
	}
}

func ImplicitCollectionNameForOrg(mspid string) string {
	return fmt.Sprintf("_implicit_org_%s", mspid)
}

func OrgFromImplicitCollectionName(name string) string {
	return strings.TrimPrefix(name, "_implicit_org_")
}

func (vc *ValidatorCommitter) ImplicitCollectionEndorsementPolicyAsBytes(channelID, orgMSPID string) (policy []byte, unexpectedErr, validationErr error) {
	channelConfig := vc.Resources.ChannelConfigSource.GetStableChannelConfig(channelID)
	if channelConfig == nil {
		return nil, errors.Errorf("could not get channel config for channel '%s'", channelID), nil
	}

	ac, ok := channelConfig.ApplicationConfig()
	if !ok {
		return nil, errors.Errorf("could not get application config for channel '%s'", channelID), nil
	}

	matchedOrgName := ""
	for orgName, org := range ac.Organizations() {
		if org.MSPID() == orgMSPID {
			matchedOrgName = orgName
			break
		}
	}

	if matchedOrgName == "" {
		return nil, nil, errors.Errorf("no org found in channel with MSPID '%s'", orgMSPID)
	}

	policyName := fmt.Sprintf("/Channel/Application/%s/Endorsement", matchedOrgName)
	if _, ok := channelConfig.PolicyManager().GetPolicy(policyName); ok {
		return protoutil.MarshalOrPanic(&pb.ApplicationPolicy{
			Type: &pb.ApplicationPolicy_ChannelConfigPolicyReference{
				ChannelConfigPolicyReference: policyName,
			},
		}), nil, nil
	}

	
	

	return protoutil.MarshalOrPanic(&pb.ApplicationPolicy{
		Type: &pb.ApplicationPolicy_SignaturePolicy{
			SignaturePolicy: cauthdsl.SignedByAnyMember([]string{orgMSPID}),
		},
	}), nil, nil
}

func (vc *ValidatorCommitter) LifecycleEndorsementPolicyAsBytes(channelID string) ([]byte, error) {
	channelConfig := vc.Resources.ChannelConfigSource.GetStableChannelConfig(channelID)
	if channelConfig == nil {
		return nil, errors.Errorf("could not get channel config for channel '%s'", channelID)
	}

	if _, ok := channelConfig.PolicyManager().GetPolicy(LifecycleEndorsementPolicyRef); ok {
		return LifecycleDefaultEndorsementPolicyBytes, nil
	}

	
	
	ac, ok := channelConfig.ApplicationConfig()
	if !ok {
		return nil, errors.Errorf("could not get application config for channel '%s'", channelID)
	}
	orgs := ac.Organizations()
	mspids := make([]string, 0, len(orgs))
	for _, org := range orgs {
		mspids = append(mspids, org.MSPID())
	}

	return protoutil.MarshalOrPanic(&pb.ApplicationPolicy{
		Type: &pb.ApplicationPolicy_SignaturePolicy{
			SignaturePolicy: cauthdsl.SignedByNOutOfGivenRole(int32(len(mspids)/2+1), msp.MSPRole_MEMBER, mspids),
		},
	}), nil
}







func (vc *ValidatorCommitter) ValidationInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (plugin string, args []byte, unexpectedErr error, validationErr error) {
	
	exists, definedChaincode, err := vc.Resources.ChaincodeDefinitionIfDefined(chaincodeName, &SimpleQueryExecutorShim{
		Namespace:           LifecycleNamespace,
		SimpleQueryExecutor: qe,
	})
	if err != nil {
		return "", nil, errors.WithMessage(err, "could not get chaincode"), nil
	}

	if !exists {
		
		
		
		return "", nil, nil, nil
	}

	if chaincodeName == LifecycleNamespace {
		b, err := vc.LifecycleEndorsementPolicyAsBytes(channelID)
		if err != nil {
			return "", nil, errors.WithMessage(err, "unexpected failure to create lifecycle endorsement policy"), nil
		}
		return "vscc", b, nil, nil
	}

	return definedChaincode.ValidationInfo.ValidationPlugin, definedChaincode.ValidationInfo.ValidationParameter, nil, nil
}


func (vc *ValidatorCommitter) CollectionValidationInfo(channelID, chaincodeName, collectionName string, state validationState.State) (args []byte, unexpectedErr, validationErr error) {
	exists, definedChaincode, err := vc.Resources.ChaincodeDefinitionIfDefined(chaincodeName, &ValidatorStateShim{
		Namespace:      LifecycleNamespace,
		ValidatorState: state,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "could not get chaincode"), nil
	}

	if !exists {
		
		
		
		return nil, nil, nil
	}

	matches := ImplicitCollectionMatcher.FindStringSubmatch(collectionName)
	if len(matches) == 2 {
		return vc.ImplicitCollectionEndorsementPolicyAsBytes(channelID, matches[1])
	}

	if definedChaincode.Collections != nil {
		for _, conf := range definedChaincode.Collections.Config {
			staticCollConfig := conf.GetStaticCollectionConfig()
			if staticCollConfig != nil && staticCollConfig.Name == collectionName {
				
				return definedChaincode.ValidationInfo.ValidationParameter, nil, nil
			}
		}
	}

	return nil, nil, errors.Errorf("no such collection '%s'", collectionName)
}
