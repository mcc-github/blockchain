/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/scc"

	"github.com/pkg/errors"
)







type LegacyLifecycle interface {
	ChaincodeDefinition(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error)
}


type ChaincodeInfoCache interface {
	ChaincodeInfo(channelID, chaincodeName string) (definition *LocalChaincodeInfo, err error)
}




type LegacyDefinition struct {
	Version           string
	EndorsementPlugin string
	RequiresInitField bool
	ChaincodeIDField  string
}


func (ld *LegacyDefinition) CCVersion() string {
	return ld.Version
}



func (ld *LegacyDefinition) Endorsement() string {
	return ld.EndorsementPlugin
}


func (ld *LegacyDefinition) RequiresInit() bool {
	return ld.RequiresInitField
}


func (ld *LegacyDefinition) ChaincodeID() string {
	return ld.ChaincodeIDField
}

type ChaincodeEndorsementInfo struct {
	Resources   *Resources
	Cache       ChaincodeInfoCache
	LegacyImpl  LegacyLifecycle
	BuiltinSCCs scc.BuiltinSCCs
}

func (cei *ChaincodeEndorsementInfo) CachedChaincodeInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (*LocalChaincodeInfo, bool, error) {
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


func (cei *ChaincodeEndorsementInfo) ChaincodeDefinition(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	if cei.BuiltinSCCs.IsSysCC(chaincodeName) {
		return &LegacyDefinition{
			Version:           scc.SysCCVersion,
			EndorsementPlugin: "escc",
			RequiresInitField: false,
			ChaincodeIDField:  scc.ChaincodeID(chaincodeName),
		}, nil
	}

	chaincodeInfo, ok, err := cei.CachedChaincodeInfo(channelID, chaincodeName, qe)
	if err != nil {
		return nil, err
	}
	if !ok {
		return cei.LegacyImpl.ChaincodeDefinition(channelID, chaincodeName, qe)
	}

	chaincodeDefinition := chaincodeInfo.Definition

	return &LegacyDefinition{
		Version:           chaincodeDefinition.EndorsementInfo.Version,
		EndorsementPlugin: chaincodeDefinition.EndorsementInfo.EndorsementPlugin,
		RequiresInitField: chaincodeDefinition.EndorsementInfo.InitRequired,

		
		ChaincodeIDField: string(chaincodeInfo.InstallInfo.PackageID),
	}, nil
}
