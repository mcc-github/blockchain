/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import "github.com/mcc-github/blockchain/common/chaincode"



type LegacyMetadataProvider interface {
	Metadata(channel string, cc string, includeCollections bool) *chaincode.Metadata
}



type ChaincodeInfoProvider interface {
	
	
	ChaincodeInfo(channelID, name string) (*LocalChaincodeInfo, error)
}

type MetadataProvider struct {
	ChaincodeInfoProvider  ChaincodeInfoProvider
	LegacyMetadataProvider LegacyMetadataProvider
}

func (mp *MetadataProvider) Metadata(channel string, ccName string, includeCollections bool) *chaincode.Metadata {
	ccInfo, err := mp.ChaincodeInfoProvider.ChaincodeInfo(channel, ccName)
	if err != nil {
		logger.Debugf("chaincode '%s' on channel '%s' not defined in _lifecycle. requesting metadata from lscc", ccName, channel)
		
		return mp.LegacyMetadataProvider.Metadata(channel, ccName, includeCollections)
	}

	ccMetadata := &chaincode.Metadata{
		Name:              ccName,
		Version:           ccInfo.Definition.EndorsementInfo.Version,
		Policy:            ccInfo.Definition.ValidationInfo.ValidationParameter,
		CollectionsConfig: ccInfo.Definition.Collections,
	}

	return ccMetadata
}
