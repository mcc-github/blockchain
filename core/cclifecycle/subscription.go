/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cc

import (
	"bytes"
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
)



type Subscription struct {
	sync.Mutex
	lc             *Lifecycle
	channel        string
	queryCreator   QueryCreator
	pendingUpdates []*cceventmgmt.ChaincodeDefinition
}

type depCCsRetriever func(Query, ChaincodePredicate, bool, ...string) (chaincode.MetadataSet, error)



func (sub *Subscription) HandleChaincodeDeploy(chaincodeDefinition *cceventmgmt.ChaincodeDefinition, dbArtifactsTar []byte) error {
	Logger.Debug("Channel", sub.channel, "got a new deployment:", chaincodeDefinition)
	sub.Lock()
	defer sub.Unlock()
	sub.pendingUpdates = append(sub.pendingUpdates, chaincodeDefinition)
	return nil
}

func (sub *Subscription) processPendingUpdate(ccDef *cceventmgmt.ChaincodeDefinition) {
	query, err := sub.queryCreator.NewQuery()
	if err != nil {
		Logger.Errorf("Failed creating a new query for channel %s: %v", sub.channel, err)
		return
	}
	installedCC := []chaincode.InstalledChaincode{{
		Name:    ccDef.Name,
		Version: ccDef.Version,
		Hash:    ccDef.Hash,
	}}
	ccs, err := queryChaincodeDefinitions(query, installedCC, DeployedChaincodes)
	if err != nil {
		Logger.Errorf("Query for channel %s for %v failed with error %v", sub.channel, ccDef, err)
		return
	}
	Logger.Debug("Updating channel", sub.channel, "with", ccs.AsChaincodes())
	sub.lc.updateState(sub.channel, ccs)
	sub.lc.fireChangeListeners(sub.channel)
}



func (sub *Subscription) ChaincodeDeployDone(succeeded bool) {
	
	
	
	
	sub.Lock()
	go func() {
		defer func() {
			sub.pendingUpdates = nil
			sub.Unlock()
		}()
		
		if !succeeded {
			Logger.Errorf("Chaincode deploy for updates %v failed", sub.pendingUpdates)
			return
		}
		for _, update := range sub.pendingUpdates {
			sub.processPendingUpdate(update)
		}
	}()
}

func queryChaincodeDefinitions(query Query, ccs []chaincode.InstalledChaincode, deployedCCs depCCsRetriever) (chaincode.MetadataSet, error) {
	
	installedCCsToIDs := make(map[nameVersion][]byte)
	
	for _, cc := range ccs {
		Logger.Debug("Chaincode", cc, "'s version is", cc.Version, "and Id is", cc.Hash)
		installedCCsToIDs[installedCCToNameVersion(cc)] = cc.Hash
	}

	filter := func(cc chaincode.Metadata) bool {
		installedID, exists := installedCCsToIDs[deployedCCToNameVersion(cc)]
		if !exists {
			Logger.Debug("Chaincode", cc, "is instantiated but a different version is installed")
			return false
		}
		if !bytes.Equal(installedID, cc.Id) {
			Logger.Debug("ID of chaincode", cc, "on filesystem doesn't match ID in ledger")
			return false
		}
		return true
	}

	return deployedCCs(query, filter, false, names(ccs)...)
}
