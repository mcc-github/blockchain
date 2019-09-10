/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/scc"

	"github.com/pkg/errors"
)




type Lifecycle interface {
	ChaincodeEndorsementInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (*ChaincodeEndorsementInfo, error)
}


type ChaincodeInfoCache interface {
	ChaincodeInfo(channelID, chaincodeName string) (definition *LocalChaincodeInfo, err error)
}


type ChaincodeEndorsementInfo struct {
	
	Version string

	
	
	EnforceInit bool

	
	ChaincodeID string

	
	EndorsementPlugin string
}

type ChaincodeEndorsementInfoSource struct {
	Resources   *Resources
	Cache       ChaincodeInfoCache
	LegacyImpl  Lifecycle
	BuiltinSCCs scc.BuiltinSCCs
}

func (cei *ChaincodeEndorsementInfoSource) CachedChaincodeInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (*LocalChaincodeInfo, bool, error) {
	var qes ReadableState = &SimpleQueryExecutorShim{
		Namespace:           LifecycleNamespace,
		SimpleQueryExecutor: qe,
	}

	if qe == nil {
		
		
		
		
		qes = &DummyQueryExecutorShim{}
	}

	currentSequence, err := cei.Resources.Serializer.DeserializeFieldAsInt64(NamespacesName, chaincodeName, "Sequence", qes)
	if err != nil {
		return nil, false, errors.WithMessagef(err, "could not get current sequence for chaincode '%s' on channel '%s'", chaincodeName, channelID)
	}

	
	if currentSequence == 0 {
		return nil, false, nil
	}

	chaincodeInfo, err := cei.Cache.ChaincodeInfo(channelID, chaincodeName)
	if err != nil {
		return nil, false, errors.WithMessage(err, "could not get approved chaincode info from cache")
	}

	if chaincodeInfo.Definition.Sequence != currentSequence {
		
		
		
		
		return nil, false, errors.Errorf("chaincode cache at sequence %d but current sequence is %d, chaincode definition for '%s' changed during invoke", chaincodeInfo.Definition.Sequence, currentSequence, chaincodeName)
	}

	if !chaincodeInfo.Approved {
		return nil, false, errors.Errorf("chaincode definition for '%s' at sequence %d on channel '%s' has not yet been approved by this org", chaincodeName, currentSequence, channelID)
	}

	if chaincodeInfo.InstallInfo == nil {
		return nil, false, errors.Errorf("chaincode definition for '%s' exists, but chaincode is not installed", chaincodeName)
	}

	return chaincodeInfo, true, nil

}



func (cei *ChaincodeEndorsementInfoSource) ChaincodeEndorsementInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (*ChaincodeEndorsementInfo, error) {
	if cei.BuiltinSCCs.IsSysCC(chaincodeName) {
		return &ChaincodeEndorsementInfo{
			Version:           scc.SysCCVersion,
			EndorsementPlugin: "escc",
			EnforceInit:       false,
			ChaincodeID:       scc.ChaincodeID(chaincodeName),
		}, nil
	}

	chaincodeInfo, ok, err := cei.CachedChaincodeInfo(channelID, chaincodeName, qe)
	if err != nil {
		return nil, err
	}
	if !ok {
		return cei.LegacyImpl.ChaincodeEndorsementInfo(channelID, chaincodeName, qe)
	}

	return &ChaincodeEndorsementInfo{
		Version:           chaincodeInfo.Definition.EndorsementInfo.Version,
		EnforceInit:       chaincodeInfo.Definition.EndorsementInfo.InitRequired,
		EndorsementPlugin: chaincodeInfo.Definition.EndorsementInfo.EndorsementPlugin,
		ChaincodeID:       chaincodeInfo.InstallInfo.PackageID, 
	}, nil
}
