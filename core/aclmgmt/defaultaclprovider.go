/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package aclmgmt

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
)

const (
	CHANNELREADERS = policies.ChannelApplicationReaders
	CHANNELWRITERS = policies.ChannelApplicationWriters
)

type defaultACLProvider interface {
	ACLProvider
	IsPtypePolicy(resName string) bool
}



type defaultACLProviderImpl struct {
	policyChecker policy.PolicyChecker

	
	pResourcePolicyMap map[string]string

	
	cResourcePolicyMap map[string]string
}

func newDefaultACLProvider(policyChecker policy.PolicyChecker) defaultACLProvider {
	d := &defaultACLProviderImpl{
		policyChecker:      policyChecker,
		pResourcePolicyMap: map[string]string{},
		cResourcePolicyMap: map[string]string{},
	}

	
	d.pResourcePolicyMap[resources.Lifecycle_InstallChaincode] = mgmt.Admins
	d.pResourcePolicyMap[resources.Lifecycle_QueryInstalledChaincode] = mgmt.Admins
	d.pResourcePolicyMap[resources.Lifecycle_GetInstalledChaincodePackage] = mgmt.Admins
	d.pResourcePolicyMap[resources.Lifecycle_QueryInstalledChaincodes] = mgmt.Admins
	d.pResourcePolicyMap[resources.Lifecycle_ApproveChaincodeDefinitionForMyOrg] = mgmt.Admins

	d.cResourcePolicyMap[resources.Lifecycle_CommitChaincodeDefinition] = CHANNELWRITERS
	d.cResourcePolicyMap[resources.Lifecycle_QueryChaincodeDefinition] = CHANNELWRITERS
	d.cResourcePolicyMap[resources.Lifecycle_QueryChaincodeDefinitions] = CHANNELWRITERS
	d.cResourcePolicyMap[resources.Lifecycle_CheckCommitReadiness] = CHANNELWRITERS

	
	
	d.pResourcePolicyMap[resources.Lscc_Install] = mgmt.Admins
	d.pResourcePolicyMap[resources.Lscc_GetInstalledChaincodes] = mgmt.Admins

	
	d.cResourcePolicyMap[resources.Lscc_Deploy] = ""  
	d.cResourcePolicyMap[resources.Lscc_Upgrade] = "" 
	d.cResourcePolicyMap[resources.Lscc_ChaincodeExists] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Lscc_GetDeploymentSpec] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Lscc_GetChaincodeData] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Lscc_GetInstantiatedChaincodes] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Lscc_GetCollectionsConfig] = CHANNELREADERS

	
	

	
	d.cResourcePolicyMap[resources.Qscc_GetChainInfo] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Qscc_GetBlockByNumber] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Qscc_GetBlockByHash] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Qscc_GetTransactionByID] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Qscc_GetBlockByTxID] = CHANNELREADERS

	
	
	d.pResourcePolicyMap[resources.Cscc_JoinChain] = mgmt.Admins
	d.pResourcePolicyMap[resources.Cscc_GetChannels] = mgmt.Members

	
	d.cResourcePolicyMap[resources.Cscc_GetConfigBlock] = CHANNELREADERS

	
	
	d.cResourcePolicyMap[resources.Peer_Propose] = CHANNELWRITERS
	d.cResourcePolicyMap[resources.Peer_ChaincodeToChaincode] = CHANNELWRITERS

	
	d.cResourcePolicyMap[resources.Event_Block] = CHANNELREADERS
	d.cResourcePolicyMap[resources.Event_FilteredBlock] = CHANNELREADERS

	return d
}

func (d *defaultACLProviderImpl) IsPtypePolicy(resName string) bool {
	_, ok := d.pResourcePolicyMap[resName]
	return ok
}


func (d *defaultACLProviderImpl) CheckACL(resName string, channelID string, idinfo interface{}) error {
	
	policy := d.pResourcePolicyMap[resName]
	if policy != "" {
		channelID = ""
	} else {
		policy = d.cResourcePolicyMap[resName]
		if policy == "" {
			aclLogger.Errorf("Unmapped policy for %s", resName)
			return fmt.Errorf("Unmapped policy for %s", resName)
		}
	}

	switch typedData := idinfo.(type) {
	case *pb.SignedProposal:
		return d.policyChecker.CheckPolicy(channelID, policy, typedData)
	case *common.Envelope:
		sd, err := protoutil.EnvelopeAsSignedData(typedData)
		if err != nil {
			return err
		}
		return d.policyChecker.CheckPolicyBySignedData(channelID, policy, sd)
	case []*protoutil.SignedData:
		return d.policyChecker.CheckPolicyBySignedData(channelID, policy, typedData)
	default:
		aclLogger.Errorf("Unmapped id on checkACL %s", resName)
		return fmt.Errorf("Unknown id on checkACL %s", resName)
	}
}
