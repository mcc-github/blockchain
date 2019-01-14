/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lscc

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/pkg/errors"
)

const (
	lsccNamespace = "lscc"
)


type DeployedCCInfoProvider struct {
}


func (p *DeployedCCInfoProvider) Namespaces() []string {
	return []string{lsccNamespace}
}


func (p *DeployedCCInfoProvider) UpdatedChaincodes(stateUpdates map[string][]*kvrwset.KVWrite) ([]*ledger.ChaincodeLifecycleInfo, error) {
	lsccUpdates := stateUpdates[lsccNamespace]
	lifecycleInfo := []*ledger.ChaincodeLifecycleInfo{}
	updatedCCNames := map[string]bool{}

	for _, kvWrite := range lsccUpdates {
		if kvWrite.IsDelete {
			
			continue
		}
		
		
		
		if privdata.IsCollectionConfigKey(kvWrite.Key) {
			ccname := privdata.GetCCNameFromCollectionConfigKey(kvWrite.Key)
			updatedCCNames[ccname] = true
			continue
		}
		updatedCCNames[kvWrite.Key] = true
	}

	for updatedCCNames := range updatedCCNames {
		lifecycleInfo = append(lifecycleInfo, &ledger.ChaincodeLifecycleInfo{Name: updatedCCNames})
	}
	return lifecycleInfo, nil
}


func (p *DeployedCCInfoProvider) ChaincodeInfo(channelName, chaincodeName string, qe ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error) {
	chaincodeDataBytes, err := qe.GetState(lsccNamespace, chaincodeName)
	if err != nil || chaincodeDataBytes == nil {
		return nil, err
	}
	chaincodeData := &ccprovider.ChaincodeData{}
	if err := proto.Unmarshal(chaincodeDataBytes, chaincodeData); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling chaincode state data")
	}
	collConfigPkg, err := fetchCollConfigPkg(chaincodeName, qe)
	if err != nil {
		return nil, err
	}
	return &ledger.DeployedChaincodeInfo{
		Name:                chaincodeName,
		Hash:                chaincodeData.Id,
		Version:             chaincodeData.Version,
		CollectionConfigPkg: collConfigPkg,
	}, nil
}


func (p *DeployedCCInfoProvider) CollectionInfo(channelName, chaincodeName, collectionName string, qe ledger.SimpleQueryExecutor) (*common.StaticCollectionConfig, error) {
	collConfigPkg, err := fetchCollConfigPkg(chaincodeName, qe)
	if err != nil || collConfigPkg == nil {
		return nil, err
	}
	for _, conf := range collConfigPkg.Config {
		staticCollConfig := conf.GetStaticCollectionConfig()
		if staticCollConfig != nil && staticCollConfig.Name == collectionName {
			return staticCollConfig, nil
		}
	}
	return nil, nil
}

func fetchCollConfigPkg(chaincodeName string, qe ledger.SimpleQueryExecutor) (*common.CollectionConfigPackage, error) {
	collKey := privdata.BuildCollectionKVSKey(chaincodeName)
	collectionConfigPkgBytes, err := qe.GetState(lsccNamespace, collKey)
	if err != nil || collectionConfigPkgBytes == nil {
		return nil, err
	}
	collectionConfigPkg := &common.CollectionConfigPackage{}
	if err := proto.Unmarshal(collectionConfigPkgBytes, collectionConfigPkg); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling chaincode collection config pkg")
	}
	return collectionConfigPkg, nil
}
