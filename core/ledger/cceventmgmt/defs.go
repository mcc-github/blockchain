/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cceventmgmt

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
)


type ChaincodeDefinition struct {
	Name    string
	Hash    []byte
	Version string
}

func (cdef *ChaincodeDefinition) String() string {
	return fmt.Sprintf("Name=%s, Version=%s, Hash=%#v", cdef.Name, cdef.Version, cdef.Hash)
}




type ChaincodeLifecycleEventListener interface {
	
	
	
	
	HandleChaincodeDeploy(chaincodeDefinition *ChaincodeDefinition, dbArtifactsTar []byte) error
	
	
	ChaincodeDeployDone(succeeded bool)
}


type ChaincodeInfoProvider interface {
	
	
	
	IsChaincodeDeployed(chainid string, chaincodeDefinition *ChaincodeDefinition, sccp sysccprovider.SystemChaincodeProvider) (bool, error)
	
	
	RetrieveChaincodeArtifacts(chaincodeDefinition *ChaincodeDefinition) (installed bool, dbArtifactsTar []byte, err error)
}

type chaincodeInfoProviderImpl struct {
	PlatformRegistry *platforms.Registry
}


func (p *chaincodeInfoProviderImpl) IsChaincodeDeployed(chainid string, chaincodeDefinition *ChaincodeDefinition, sccp sysccprovider.SystemChaincodeProvider) (bool, error) {
	return ccprovider.IsChaincodeDeployed(chainid, chaincodeDefinition.Name, chaincodeDefinition.Version, chaincodeDefinition.Hash, sccp)
}


func (p *chaincodeInfoProviderImpl) RetrieveChaincodeArtifacts(chaincodeDefinition *ChaincodeDefinition) (installed bool, dbArtifactsTar []byte, err error) {
	return ccprovider.ExtractStatedbArtifactsForChaincode(chaincodeDefinition.Name, chaincodeDefinition.Version, p.PlatformRegistry)
}
