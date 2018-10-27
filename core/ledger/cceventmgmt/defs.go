/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cceventmgmt

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
)


type ChaincodeDefinition struct {
	Name              string
	Hash              []byte
	Version           string
	CollectionConfigs *common.CollectionConfigPackage
}

func (cdef *ChaincodeDefinition) String() string {
	return fmt.Sprintf("Name=%s, Version=%s, Hash=%#v", cdef.Name, cdef.Version, cdef.Hash)
}




type ChaincodeLifecycleEventListener interface {
	
	
	
	
	HandleChaincodeDeploy(chaincodeDefinition *ChaincodeDefinition, dbArtifactsTar []byte) error
	
	
	ChaincodeDeployDone(succeeded bool)
}


type ChaincodeInfoProvider interface {
	
	
	
	GetDeployedChaincodeInfo(chainid string, chaincodeDefinition *ChaincodeDefinition) (*ledger.DeployedChaincodeInfo, error)
	
	
	RetrieveChaincodeArtifacts(chaincodeDefinition *ChaincodeDefinition) (installed bool, dbArtifactsTar []byte, err error)
}
