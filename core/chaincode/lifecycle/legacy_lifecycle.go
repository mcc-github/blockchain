/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"strings"

	"github.com/mcc-github/blockchain/common/util"
	corechaincode "github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"

	"github.com/pkg/errors"
)


type LegacyLifecycle interface {
	corechaincode.Lifecycle
}


type ChaincodeInfoCache interface {
	ChaincodeInfo(channelID, chaincodeName string) (definition *LocalChaincodeInfo, err error)
}




type LegacyDefinition struct {
	Name                string
	Version             string
	HashField           []byte
	EndorsementPlugin   string
	ValidationPlugin    string
	ValidationParameter []byte
	RequiresInitField   bool
}


func (ld *LegacyDefinition) CCName() string {
	return ld.Name
}





func (ld *LegacyDefinition) Hash() []byte {
	return util.ComputeSHA256([]byte(ld.Name + ":" + ld.Version))
}


func (ld *LegacyDefinition) CCVersion() string {
	return ld.Version
}





func (ld *LegacyDefinition) Validation() (string, []byte) {
	return ld.ValidationPlugin, ld.ValidationParameter
}



func (ld *LegacyDefinition) Endorsement() string {
	return ld.EndorsementPlugin
}


func (ld *LegacyDefinition) RequiresInit() bool {
	return ld.RequiresInitField
}

type ChaincodeEndorsementInfo struct {
	Resources  *Resources
	Cache      ChaincodeInfoCache
	LegacyImpl LegacyLifecycle
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
	chaincodeInfo, ok, err := cei.CachedChaincodeInfo(channelID, chaincodeName, qe)
	if err != nil {
		return nil, err
	}
	if !ok {
		return cei.LegacyImpl.ChaincodeDefinition(channelID, chaincodeName, qe)
	}

	chaincodeDefinition := chaincodeInfo.Definition

	return &LegacyDefinition{
		Name:                chaincodeName,
		Version:             chaincodeDefinition.EndorsementInfo.Version,
		EndorsementPlugin:   chaincodeDefinition.EndorsementInfo.EndorsementPlugin,
		RequiresInitField:   chaincodeDefinition.EndorsementInfo.InitRequired,
		ValidationPlugin:    chaincodeDefinition.ValidationInfo.ValidationPlugin,
		ValidationParameter: chaincodeDefinition.ValidationInfo.ValidationParameter,
	}, nil
}


func (cei *ChaincodeEndorsementInfo) ChaincodeContainerInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (*ccprovider.ChaincodeContainerInfo, error) {
	chaincodeInfo, ok, err := cei.CachedChaincodeInfo(channelID, chaincodeName, qe)
	if err != nil {
		return nil, err
	}
	if !ok {
		return cei.LegacyImpl.ChaincodeContainerInfo(channelID, chaincodeName, qe)
	}

	return &ccprovider.ChaincodeContainerInfo{
		Name:      chaincodeName,
		Version:   chaincodeInfo.Definition.EndorsementInfo.Version,
		Path:      chaincodeInfo.InstallInfo.Path,
		Type:      strings.ToUpper(chaincodeInfo.InstallInfo.Type),
		PackageID: chaincodeInfo.InstallInfo.PackageID,
	}, nil
}
