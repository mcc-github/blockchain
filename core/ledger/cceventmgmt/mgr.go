/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cceventmgmt

import (
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
)

var logger = flogging.MustGetLogger("cceventmgmt")

var mgr *Mgr


func Initialize(ccInfoProvider ChaincodeInfoProvider) {
	initialize(ccInfoProvider)
}

func initialize(ccInfoProvider ChaincodeInfoProvider) {
	mgr = newMgr(ccInfoProvider)
}


func GetMgr() *Mgr {
	return mgr
}


type Mgr struct {
	
	
	
	
	rwlock               sync.RWMutex
	infoProvider         ChaincodeInfoProvider
	ccLifecycleListeners map[string][]ChaincodeLifecycleEventListener
	callbackStatus       *callbackStatus
}

func newMgr(chaincodeInfoProvider ChaincodeInfoProvider) *Mgr {
	return &Mgr{
		infoProvider:         chaincodeInfoProvider,
		ccLifecycleListeners: make(map[string][]ChaincodeLifecycleEventListener),
		callbackStatus:       newCallbackStatus()}
}



func (m *Mgr) Register(ledgerid string, l ChaincodeLifecycleEventListener) {
	
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	m.ccLifecycleListeners[ledgerid] = append(m.ccLifecycleListeners[ledgerid], l)
}









func (m *Mgr) HandleChaincodeDeploy(chainid string, chaincodeDefinitions []*ChaincodeDefinition) error {
	logger.Debugf("Channel [%s]: Handling chaincode deploy event for chaincode [%s]", chainid, chaincodeDefinitions)
	
	m.rwlock.RLock()
	for _, chaincodeDefinition := range chaincodeDefinitions {
		installed, dbArtifacts, err := m.infoProvider.RetrieveChaincodeArtifacts(chaincodeDefinition)
		if err != nil {
			return err
		}
		if !installed {
			logger.Infof("Channel [%s]: Chaincode [%s] is not installed hence no need to create chaincode artifacts for endorsement",
				chainid, chaincodeDefinition)
			continue
		}
		m.callbackStatus.setDeployPending(chainid)
		if err := m.invokeHandler(chainid, chaincodeDefinition, dbArtifacts); err != nil {
			logger.Warningf("Channel [%s]: Error while invoking a listener for handling chaincode install event: %s", chainid, err)
			return err
		}
		logger.Debugf("Channel [%s]: Handled chaincode deploy event for chaincode [%s]", chainid, chaincodeDefinitions)
	}
	return nil
}


func (m *Mgr) ChaincodeDeployDone(chainid string) {
	
	defer m.rwlock.RUnlock()
	if m.callbackStatus.isDeployPending(chainid) {
		m.invokeDoneOnHandlers(chainid, true)
		m.callbackStatus.unsetDeployPending(chainid)
	}
}


func (m *Mgr) HandleChaincodeInstall(chaincodeDefinition *ChaincodeDefinition, dbArtifacts []byte) error {
	logger.Debugf("HandleChaincodeInstall() - chaincodeDefinition=%#v", chaincodeDefinition)
	
	m.rwlock.Lock()
	for chainid := range m.ccLifecycleListeners {
		logger.Debugf("Channel [%s]: Handling chaincode install event for chaincode [%s]", chainid, chaincodeDefinition)
		var deployedCCInfo *ledger.DeployedChaincodeInfo
		var err error
		if deployedCCInfo, err = m.infoProvider.GetDeployedChaincodeInfo(chainid, chaincodeDefinition); err != nil {
			logger.Warningf("Channel [%s]: Error while getting the deployment status of chaincode: %s", chainid, err)
			return err
		}
		if deployedCCInfo == nil {
			logger.Debugf("Channel [%s]: Chaincode [%s] is not deployed on channel hence not creating chaincode artifacts.",
				chainid, chaincodeDefinition)
			continue
		}
		if !deployedCCInfo.IsLegacy {
			
			
			logger.Debugf("Channel [%s]: Chaincode [%s] is already defined in new lifecycle hence not creating chaincode artifacts.",
				chainid, chaincodeDefinition)
			continue
		}
		m.callbackStatus.setInstallPending(chainid)
		chaincodeDefinition.CollectionConfigs = deployedCCInfo.ExplicitCollectionConfigPkg
		if err := m.invokeHandler(chainid, chaincodeDefinition, dbArtifacts); err != nil {
			logger.Warningf("Channel [%s]: Error while invoking a listener for handling chaincode install event: %s", chainid, err)
			return err
		}
		logger.Debugf("Channel [%s]: Handled chaincode install event for chaincode [%s]", chainid, chaincodeDefinition)
	}
	return nil
}


func (m *Mgr) ChaincodeInstallDone(succeeded bool) {
	
	defer m.rwlock.Unlock()
	for chainid := range m.callbackStatus.installPending {
		m.invokeDoneOnHandlers(chainid, succeeded)
		m.callbackStatus.unsetInstallPending(chainid)
	}
}

func (m *Mgr) invokeHandler(chainid string, chaincodeDefinition *ChaincodeDefinition, dbArtifactsTar []byte) error {
	listeners := m.ccLifecycleListeners[chainid]
	for _, listener := range listeners {
		if err := listener.HandleChaincodeDeploy(chaincodeDefinition, dbArtifactsTar); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mgr) invokeDoneOnHandlers(chainid string, succeeded bool) {
	listeners := m.ccLifecycleListeners[chainid]
	for _, listener := range listeners {
		listener.ChaincodeDeployDone(succeeded)
	}
}

type callbackStatus struct {
	l              sync.Mutex
	deployPending  map[string]bool
	installPending map[string]bool
}

func newCallbackStatus() *callbackStatus {
	return &callbackStatus{
		deployPending:  make(map[string]bool),
		installPending: make(map[string]bool)}
}

func (s *callbackStatus) setDeployPending(channelID string) {
	s.l.Lock()
	defer s.l.Unlock()
	s.deployPending[channelID] = true
}

func (s *callbackStatus) unsetDeployPending(channelID string) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.deployPending, channelID)
}

func (s *callbackStatus) isDeployPending(channelID string) bool {
	s.l.Lock()
	defer s.l.Unlock()
	return s.deployPending[channelID]
}

func (s *callbackStatus) setInstallPending(channelID string) {
	s.l.Lock()
	defer s.l.Unlock()
	s.installPending[channelID] = true
}

func (s *callbackStatus) unsetInstallPending(channelID string) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.installPending, channelID)
}
