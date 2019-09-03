/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgermgmt

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/ledger/kvledger"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("ledgermgmt")


var ErrLedgerAlreadyOpened = errors.New("ledger already opened")


var ErrLedgerMgmtNotInitialized = errors.New("ledger mgmt should be initialized before using")


type LedgerMgr struct {
	lock           sync.Mutex
	openedLedgers  map[string]ledger.PeerLedger
	ledgerProvider ledger.PeerLedgerProvider
}


type Initializer struct {
	CustomTxProcessors              map[common.HeaderType]ledger.CustomTxProcessor
	StateListeners                  []ledger.StateListener
	DeployedChaincodeInfoProvider   ledger.DeployedChaincodeInfoProvider
	MembershipInfoProvider          ledger.MembershipInfoProvider
	ChaincodeLifecycleEventProvider ledger.ChaincodeLifecycleEventProvider
	MetricsProvider                 metrics.Provider
	HealthCheckRegistry             ledger.HealthCheckRegistry
	Config                          *ledger.Config
}


func NewLedgerMgr(initializer *Initializer) *LedgerMgr {
	logger.Info("Initializing LedgerMgr")
	finalStateListeners := addListenerForCCEventsHandler(
		initializer.DeployedChaincodeInfoProvider,
		initializer.StateListeners,
	)
	provider, err := kvledger.NewProvider(
		&ledger.Initializer{
			StateListeners:                  finalStateListeners,
			DeployedChaincodeInfoProvider:   initializer.DeployedChaincodeInfoProvider,
			MembershipInfoProvider:          initializer.MembershipInfoProvider,
			ChaincodeLifecycleEventProvider: initializer.ChaincodeLifecycleEventProvider,
			MetricsProvider:                 initializer.MetricsProvider,
			HealthCheckRegistry:             initializer.HealthCheckRegistry,
			Config:                          initializer.Config,
			CustomTxProcessors:              initializer.CustomTxProcessors,
		},
	)
	if err != nil {
		panic(fmt.Sprintf("Error in instantiating ledger provider: %s", err))
	}
	ledgerMgr := &LedgerMgr{
		openedLedgers:  make(map[string]ledger.PeerLedger),
		ledgerProvider: provider,
	}
	
	cceventmgmt.Initialize(&chaincodeInfoProviderImpl{
		ledgerMgr,
		initializer.DeployedChaincodeInfoProvider,
	})
	logger.Info("Initialized LedgerMgr")
	return ledgerMgr
}




func (m *LedgerMgr) CreateLedger(id string, genesisBlock *common.Block) (ledger.PeerLedger, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	logger.Infof("Creating ledger [%s] with genesis block", id)
	l, err := m.ledgerProvider.Create(genesisBlock)
	if err != nil {
		return nil, err
	}
	m.openedLedgers[id] = l
	logger.Infof("Created ledger [%s] with genesis block", id)
	return &closableLedger{
		ledgerMgr:  m,
		id:         id,
		PeerLedger: l,
	}, nil
}


func (m *LedgerMgr) OpenLedger(id string) (ledger.PeerLedger, error) {
	logger.Infof("Opening ledger with id = %s", id)
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.openedLedgers[id]
	if ok {
		return nil, ErrLedgerAlreadyOpened
	}
	l, err := m.ledgerProvider.Open(id)
	if err != nil {
		return nil, err
	}
	m.openedLedgers[id] = l
	logger.Infof("Opened ledger with id = %s", id)
	return &closableLedger{
		ledgerMgr:  m,
		id:         id,
		PeerLedger: l,
	}, nil
}


func (m *LedgerMgr) GetLedgerIDs() ([]string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.ledgerProvider.List()
}


func (m *LedgerMgr) Close() {
	logger.Infof("Closing ledger mgmt")
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, l := range m.openedLedgers {
		l.Close()
	}
	m.ledgerProvider.Close()
	m.openedLedgers = nil
	logger.Infof("ledger mgmt closed")
}

func (m *LedgerMgr) getOpenedLedger(ledgerID string) (ledger.PeerLedger, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	l, ok := m.openedLedgers[ledgerID]
	if !ok {
		return nil, errors.Errorf("Ledger not opened [%s]", ledgerID)
	}
	return l, nil
}

func (m *LedgerMgr) closeLedger(ledgerID string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	l, ok := m.openedLedgers[ledgerID]
	if ok {
		l.Close()
		delete(m.openedLedgers, ledgerID)
	}
}


type closableLedger struct {
	ledgerMgr *LedgerMgr
	id        string
	ledger.PeerLedger
}


func (l *closableLedger) Close() {
	l.ledgerMgr.closeLedger(l.id)
}



func addListenerForCCEventsHandler(
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
	stateListeners []ledger.StateListener) []ledger.StateListener {
	return append(stateListeners, &cceventmgmt.KVLedgerLSCCStateListener{DeployedChaincodeInfoProvider: deployedCCInfoProvider})
}


type chaincodeInfoProviderImpl struct {
	ledgerMgr              *LedgerMgr
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider
}


func (p *chaincodeInfoProviderImpl) GetDeployedChaincodeInfo(chainid string,
	chaincodeDefinition *cceventmgmt.ChaincodeDefinition) (*ledger.DeployedChaincodeInfo, error) {
	ledger, err := p.ledgerMgr.getOpenedLedger(chainid)
	if err != nil {
		return nil, err
	}
	qe, err := ledger.NewQueryExecutor()
	if err != nil {
		return nil, err
	}
	defer qe.Done()
	deployedChaincodeInfo, err := p.deployedCCInfoProvider.ChaincodeInfo(chainid, chaincodeDefinition.Name, qe)
	if err != nil || deployedChaincodeInfo == nil {
		return nil, err
	}
	if deployedChaincodeInfo.Version != chaincodeDefinition.Version ||
		!bytes.Equal(deployedChaincodeInfo.Hash, chaincodeDefinition.Hash) {
		
		return nil, nil
	}
	return deployedChaincodeInfo, nil
}


func (p *chaincodeInfoProviderImpl) RetrieveChaincodeArtifacts(chaincodeDefinition *cceventmgmt.ChaincodeDefinition) (installed bool, dbArtifactsTar []byte, err error) {
	return ccprovider.ExtractStatedbArtifactsForChaincode(chaincodeDefinition.Name + ":" + chaincodeDefinition.Version)
}
