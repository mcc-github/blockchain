/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgermgmt

import (
	"bytes"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/ledger/customtx"
	"github.com/mcc-github/blockchain/core/ledger/kvledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("ledgermgmt")


var ErrLedgerAlreadyOpened = errors.New("ledger already opened")


var ErrLedgerMgmtNotInitialized = errors.New("ledger mgmt should be initialized before using")

var openedLedgers map[string]ledger.PeerLedger
var ledgerProvider ledger.PeerLedgerProvider
var lock sync.Mutex
var initialized bool
var once sync.Once


type Initializer struct {
	CustomTxProcessors            customtx.Processors
	StateListeners                []ledger.StateListener
	PlatformRegistry              *platforms.Registry
	DeployedChaincodeInfoProvider ledger.DeployedChaincodeInfoProvider
	MembershipInfoProvider        ledger.MembershipInfoProvider
	MetricsProvider               metrics.Provider
	HealthCheckRegistry           ledger.HealthCheckRegistry
	Config                        *ledger.Config
}


func Initialize(initializer *Initializer) {
	once.Do(func() {
		initialize(initializer)
	})
}

func initialize(initializer *Initializer) {
	logger.Info("Initializing ledger mgmt")
	lock.Lock()
	defer lock.Unlock()
	initialized = true
	openedLedgers = make(map[string]ledger.PeerLedger)
	customtx.Initialize(initializer.CustomTxProcessors)
	cceventmgmt.Initialize(&chaincodeInfoProviderImpl{
		initializer.PlatformRegistry,
		initializer.DeployedChaincodeInfoProvider,
	})
	finalStateListeners := addListenerForCCEventsHandler(initializer.DeployedChaincodeInfoProvider, initializer.StateListeners)
	provider, err := kvledger.NewProvider()
	if err != nil {
		panic(errors.WithMessage(err, "Error in instantiating ledger provider"))
	}
	err = provider.Initialize(&ledger.Initializer{
		StateListeners:                finalStateListeners,
		DeployedChaincodeInfoProvider: initializer.DeployedChaincodeInfoProvider,
		MembershipInfoProvider:        initializer.MembershipInfoProvider,
		MetricsProvider:               initializer.MetricsProvider,
		HealthCheckRegistry:           initializer.HealthCheckRegistry,
		Config:                        initializer.Config,
	})
	if err != nil {
		panic(errors.WithMessage(err, "Error initializing ledger provider"))
	}
	ledgerProvider = provider
	logger.Info("ledger mgmt initialized")
}




func CreateLedger(genesisBlock *common.Block) (ledger.PeerLedger, error) {
	lock.Lock()
	defer lock.Unlock()
	if !initialized {
		return nil, ErrLedgerMgmtNotInitialized
	}
	id, err := protoutil.GetChainIDFromBlock(genesisBlock)
	if err != nil {
		return nil, err
	}

	logger.Infof("Creating ledger [%s] with genesis block", id)
	l, err := ledgerProvider.Create(genesisBlock)
	if err != nil {
		return nil, err
	}
	l = wrapLedger(id, l)
	openedLedgers[id] = l
	logger.Infof("Created ledger [%s] with genesis block", id)
	return l, nil
}


func OpenLedger(id string) (ledger.PeerLedger, error) {
	logger.Infof("Opening ledger with id = %s", id)
	lock.Lock()
	defer lock.Unlock()
	if !initialized {
		return nil, ErrLedgerMgmtNotInitialized
	}
	l, ok := openedLedgers[id]
	if ok {
		return nil, ErrLedgerAlreadyOpened
	}
	l, err := ledgerProvider.Open(id)
	if err != nil {
		return nil, err
	}
	l = wrapLedger(id, l)
	openedLedgers[id] = l
	logger.Infof("Opened ledger with id = %s", id)
	return l, nil
}


func GetLedgerIDs() ([]string, error) {
	lock.Lock()
	defer lock.Unlock()
	if !initialized {
		return nil, ErrLedgerMgmtNotInitialized
	}
	return ledgerProvider.List()
}


func Close() {
	logger.Infof("Closing ledger mgmt")
	lock.Lock()
	defer lock.Unlock()
	if !initialized {
		return
	}
	for _, l := range openedLedgers {
		l.(*closableLedger).closeWithoutLock()
	}
	ledgerProvider.Close()
	openedLedgers = nil
	logger.Infof("ledger mgmt closed")
}

func wrapLedger(id string, l ledger.PeerLedger) ledger.PeerLedger {
	return &closableLedger{id, l}
}


type closableLedger struct {
	id string
	ledger.PeerLedger
}


func (l *closableLedger) Close() {
	lock.Lock()
	defer lock.Unlock()
	l.closeWithoutLock()
}

func (l *closableLedger) closeWithoutLock() {
	l.PeerLedger.Close()
	delete(openedLedgers, l.id)
}



func addListenerForCCEventsHandler(
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
	stateListeners []ledger.StateListener) []ledger.StateListener {
	return append(stateListeners, &cceventmgmt.KVLedgerLSCCStateListener{DeployedChaincodeInfoProvider: deployedCCInfoProvider})
}


type chaincodeInfoProviderImpl struct {
	pr                     *platforms.Registry
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider
}


func (p *chaincodeInfoProviderImpl) GetDeployedChaincodeInfo(chainid string,
	chaincodeDefinition *cceventmgmt.ChaincodeDefinition) (*ledger.DeployedChaincodeInfo, error) {
	lock.Lock()
	ledger := openedLedgers[chainid]
	lock.Unlock()
	if ledger == nil {
		return nil, errors.Errorf("Ledger not opened [%s]", chainid)
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
	return ccprovider.ExtractStatedbArtifactsForChaincode(chaincodeDefinition.Name, chaincodeDefinition.Version, p.pr)
}
