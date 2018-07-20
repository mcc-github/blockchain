/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgermgmt

import (
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/ledger/customtx"
	"github.com/mcc-github/blockchain/core/ledger/kvledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
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


func Initialize(customTxProcessors customtx.Processors, pr *platforms.Registry) {
	once.Do(func() {
		initialize(customTxProcessors, nil, pr)
	})
}

func initialize(customTxProcessors customtx.Processors, statelisteners []ledger.StateListener, pr *platforms.Registry) {
	logger.Info("Initializing ledger mgmt")
	lock.Lock()
	defer lock.Unlock()
	initialized = true
	openedLedgers = make(map[string]ledger.PeerLedger)
	customtx.Initialize(customTxProcessors)
	cceventmgmt.Initialize(pr)
	finalStateListeners := addListenerForCCEventsHandler(statelisteners)
	provider, err := kvledger.NewProvider()
	if err != nil {
		panic(errors.WithMessage(err, "Error in instantiating ledger provider"))
	}
	provider.Initialize(finalStateListeners)
	ledgerProvider = provider
	logger.Info("ledger mgmt initialized")
}




func CreateLedger(genesisBlock *common.Block) (ledger.PeerLedger, error) {
	lock.Lock()
	defer lock.Unlock()
	if !initialized {
		return nil, ErrLedgerMgmtNotInitialized
	}
	id, err := utils.GetChainIDFromBlock(genesisBlock)
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



func addListenerForCCEventsHandler(stateListeners []ledger.StateListener) []ledger.StateListener {
	return append(stateListeners, &cceventmgmt.KVLedgerLSCCStateListener{})
}
