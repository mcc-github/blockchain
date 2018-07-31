/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cceventmgmt

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/pkg/errors"
)

const (
	lsccNamespace = "lscc"
)


type KVLedgerLSCCStateListener struct {
}




func (listener *KVLedgerLSCCStateListener) HandleStateUpdates(channelName string, stateUpdates ledger.StateUpdates, committingBlockNum uint64) error {
	kvWrites := stateUpdates[lsccNamespace].([]*kvrwset.KVWrite)
	logger.Debugf("Channel [%s]: Handling state updates in LSCC namespace - stateUpdates=%#v", channelName, kvWrites)
	chaincodeDefs := []*ChaincodeDefinition{}
	chaincodesCollConfigs := make(map[string][]byte)

	for _, kvWrite := range kvWrites {
		
		
		
		if privdata.IsCollectionConfigKey(kvWrite.Key) {
			ccname := privdata.GetCCNameFromCollectionConfigKey(kvWrite.Key)
			chaincodesCollConfigs[ccname] = kvWrite.Value
			continue
		}
		
		if kvWrite.IsDelete {
			continue
		}
		
		logger.Infof("Channel [%s]: Handling LSCC state update for chaincode [%s]", channelName, kvWrite.Key)
		chaincodeData := &ccprovider.ChaincodeData{}
		if err := proto.Unmarshal(kvWrite.Value, chaincodeData); err != nil {
			return errors.Wrap(err, "error unmarshalling chaincode state data")
		}
		chaincodeDefs = append(chaincodeDefs, &ChaincodeDefinition{Name: chaincodeData.CCName(), Version: chaincodeData.CCVersion(), Hash: chaincodeData.Hash()})
	}

	for _, chaincodeDef := range chaincodeDefs {
		chaincodeCollConfigs, ok := chaincodesCollConfigs[chaincodeDef.Name]
		if ok {
			chaincodeDef.CollectionConfigs = chaincodeCollConfigs
		}
	}

	return GetMgr().HandleChaincodeDeploy(channelName, chaincodeDefs)
}


func (listener *KVLedgerLSCCStateListener) InterestedInNamespaces() []string {
	return []string{lsccNamespace}
}


func (listener *KVLedgerLSCCStateListener) StateCommitDone(channelName string) {
	GetMgr().ChaincodeDeployDone(channelName)
}
