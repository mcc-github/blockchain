/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cceventmgmt

import (
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
)


type KVLedgerLSCCStateListener struct {
	DeployedChaincodeInfoProvider ledger.DeployedChaincodeInfoProvider
}




func (listener *KVLedgerLSCCStateListener) HandleStateUpdates(trigger *ledger.StateUpdateTrigger) error {
	channelName, kvWrites, postCommitQE, deployCCInfoProvider :=
		trigger.LedgerID, convertToKVWrites(trigger.StateUpdates), trigger.PostCommitQueryExecutor, listener.DeployedChaincodeInfoProvider

	logger.Debugf("Channel [%s]: Handling state updates in LSCC namespace - stateUpdates=%#v", channelName, kvWrites)
	updatedChaincodes, err := deployCCInfoProvider.UpdatedChaincodes(kvWrites)
	if err != nil {
		return err
	}
	chaincodeDefs := []*ChaincodeDefinition{}
	for _, updatedChaincode := range updatedChaincodes {
		logger.Infof("Channel [%s]: Handling deploy or update of chaincode [%s]", channelName, updatedChaincode.Name)
		if updatedChaincode.Deleted {
			
			continue
		}
		deployedCCInfo, err := deployCCInfoProvider.ChaincodeInfo(channelName, updatedChaincode.Name, postCommitQE)
		if err != nil {
			return err
		}
		chaincodeDefs = append(chaincodeDefs, &ChaincodeDefinition{
			Name:              deployedCCInfo.Name,
			Hash:              deployedCCInfo.Hash,
			Version:           deployedCCInfo.Version,
			CollectionConfigs: deployedCCInfo.CollectionConfigPkg,
		})
	}
	return GetMgr().HandleChaincodeDeploy(channelName, chaincodeDefs)
}


func (listener *KVLedgerLSCCStateListener) InterestedInNamespaces() []string {
	return listener.DeployedChaincodeInfoProvider.Namespaces()
}


func (listener *KVLedgerLSCCStateListener) StateCommitDone(channelName string) {
	GetMgr().ChaincodeDeployDone(channelName)
}

func convertToKVWrites(stateUpdates ledger.StateUpdates) map[string][]*kvrwset.KVWrite {
	m := map[string][]*kvrwset.KVWrite{}
	for ns, updates := range stateUpdates {
		m[ns] = updates.([]*kvrwset.KVWrite)
	}
	return m
}
