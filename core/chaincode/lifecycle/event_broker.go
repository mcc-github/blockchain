/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"sync"

	ccpersistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/pkg/errors"
)


type EventBroker struct {
	chaincodeStore       ChaincodeStore
	pkgParser            PackageParser
	defineCallbackStatus *sync.Map

	mutex     sync.Mutex
	listeners map[string][]ledger.ChaincodeLifecycleEventListener
}

func NewEventBroker(chaincodeStore ChaincodeStore, pkgParser PackageParser) *EventBroker {
	return &EventBroker{
		chaincodeStore:       chaincodeStore,
		pkgParser:            pkgParser,
		listeners:            make(map[string][]ledger.ChaincodeLifecycleEventListener),
		defineCallbackStatus: &sync.Map{},
	}
}

func (b *EventBroker) RegisterListener(channelID string, listener ledger.ChaincodeLifecycleEventListener) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.listeners[channelID] = append(b.listeners[channelID], listener)
}


func (b *EventBroker) ProcessInstallEvent(localChaincode *LocalChaincode) {
	logger.Debugf("ProcessInstallEvent() - localChaincode = %s", localChaincode.Info)
	dbArtifacts, err := b.loadDBArtifacts(localChaincode.Info.PackageID)
	if err != nil {
		logger.Errorf("Error while loading db artifacts for chaincode package with package ID [%s]: %s",
			localChaincode.Info.PackageID, err)
		return
	}
	for channelID, channelCache := range localChaincode.References {
		listenersInvokedOnChannel := false
		for chaincodeName, cachedChaincode := range channelCache {
			if !isChaincodeInvokable(cachedChaincode) {
				continue
			}
			ccdef := &ledger.ChaincodeDefinition{
				Name:              chaincodeName,
				Version:           cachedChaincode.Definition.EndorsementInfo.Version,
				Hash:              []byte(cachedChaincode.InstallInfo.PackageID),
				CollectionConfigs: cachedChaincode.Definition.Collections,
			}
			b.invokeListeners(channelID, ccdef, dbArtifacts)
			listenersInvokedOnChannel = true
		}
		if listenersInvokedOnChannel {
			
			
			
			
			
			
			
			
			
			b.invokeDoneOnListeners(channelID, true)
		}
	}
	return
}





func (b *EventBroker) ProcessApproveOrDefineEvent(channelID string, chaincodeName string, cachedChaincode *CachedChaincodeDefinition) {
	logger.Debugw("processApproveOrDefineEvent()", "channelID", channelID, "chaincodeName", chaincodeName, "cachedChaincode", cachedChaincode)
	if !isChaincodeInvokable(cachedChaincode) {
		return
	}
	dbArtifacts, err := b.loadDBArtifacts(cachedChaincode.InstallInfo.PackageID)
	if err != nil {
		logger.Errorf("Error while loading db artifacts for chaincode package with package ID [%s]: %s",
			cachedChaincode.InstallInfo.PackageID, err)
		return
	}
	ccdef := &ledger.ChaincodeDefinition{
		Name:              chaincodeName,
		Version:           cachedChaincode.Definition.EndorsementInfo.Version,
		Hash:              []byte(cachedChaincode.InstallInfo.PackageID),
		CollectionConfigs: cachedChaincode.Definition.Collections,
	}
	b.invokeListeners(channelID, ccdef, dbArtifacts)
	b.defineCallbackStatus.Store(channelID, struct{}{})
	return
}



func (b *EventBroker) ApproveOrDefineCommitted(channelID string) {
	_, ok := b.defineCallbackStatus.Load(channelID)
	if !ok {
		return
	}
	b.invokeDoneOnListeners(channelID, true)
	b.defineCallbackStatus.Delete(channelID)
}

func (b *EventBroker) invokeListeners(channelID string, legacyDefinition *ledger.ChaincodeDefinition, dbArtifacts []byte) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	channelListeners := b.listeners[channelID]
	for _, l := range channelListeners {
		if err := l.HandleChaincodeDeploy(legacyDefinition, dbArtifacts); err != nil {
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			logger.Errorf("Error from listener during processing chaincode lifecycle event - %+v", errors.WithStack(err))
		}
	}
}

func (b *EventBroker) invokeDoneOnListeners(channelID string, succeeded bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	channelListeners := b.listeners[channelID]
	for _, l := range channelListeners {
		l.ChaincodeDeployDone(succeeded)
	}
}

func (b *EventBroker) loadDBArtifacts(packageID ccpersistence.PackageID) ([]byte, error) {
	pkgBytes, err := b.chaincodeStore.Load(packageID)
	if err != nil {
		return nil, err
	}
	pkg, err := b.pkgParser.Parse(pkgBytes)
	if err != nil {
		return nil, err
	}
	return pkg.DBArtifacts, nil
}


func isChaincodeInvokable(ccInfo *CachedChaincodeDefinition) bool {
	return ccInfo.Approved && ccInfo.InstallInfo != nil && ccInfo.Definition != nil
}
