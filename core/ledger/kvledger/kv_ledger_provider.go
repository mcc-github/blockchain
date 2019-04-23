/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/confighistory"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/history/historydb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/history/historydb/historyleveldb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/ledgerstorage"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	
	ErrLedgerIDExists = errors.New("LedgerID already exists")
	
	ErrNonExistingLedgerID = errors.New("LedgerID does not exist")
	
	ErrLedgerNotOpened = errors.New("ledger is not opened yet")

	underConstructionLedgerKey = []byte("underConstructionLedgerKey")
	ledgerKeyPrefix            = []byte("l")
)


type Provider struct {
	idStore             *idStore
	ledgerStoreProvider *ledgerstorage.Provider
	vdbProvider         privacyenabledstate.DBProvider
	historydbProvider   historydb.HistoryDBProvider
	configHistoryMgr    confighistory.Mgr
	stateListeners      []ledger.StateListener
	bookkeepingProvider bookkeeping.Provider
	initializer         *ledger.Initializer
	collElgNotifier     *collElgNotifier
	stats               *stats
}



func NewProvider() (ledger.PeerLedgerProvider, error) {
	logger.Info("Initializing ledger provider")

	logger.Info("ledger provider Initialized")
	p := &Provider{}
	return p, nil
}


func (p *Provider) Initialize(initializer *ledger.Initializer) error {
	var err error

	p.initializer = initializer
	
	idStore := openIDStore(filepath.Join(p.initializer.Config.RootFSPath, "ledgerProvider"))
	
	ledgerStoreProvider := ledgerstorage.NewProvider(
		p.initializer.Config.RootFSPath,
		p.initializer.Config.PrivateData,
	)
	
	historydbProvider := historyleveldb.NewHistoryDBProvider(
		filepath.Join(p.initializer.Config.RootFSPath, "historyLeveldb"),
	)
	
	configHistoryMgr := confighistory.NewMgr(
		filepath.Join(p.initializer.Config.RootFSPath, "configHistory"),
		initializer.DeployedChaincodeInfoProvider,
	)
	
	collElgNotifier := &collElgNotifier{
		initializer.DeployedChaincodeInfoProvider,
		initializer.MembershipInfoProvider,
		make(map[string]collElgListener),
	}
	
	stateListeners := initializer.StateListeners
	stateListeners = append(stateListeners, collElgNotifier)
	stateListeners = append(stateListeners, configHistoryMgr)

	p.idStore = idStore
	p.ledgerStoreProvider = ledgerStoreProvider
	p.historydbProvider = historydbProvider
	p.configHistoryMgr = configHistoryMgr
	p.stateListeners = stateListeners
	p.collElgNotifier = collElgNotifier
	p.bookkeepingProvider = bookkeeping.NewProvider(
		filepath.Join(p.initializer.Config.RootFSPath, "bookkeeper"),
	)
	p.vdbProvider, err = privacyenabledstate.NewCommonStorageDBProvider(
		p.bookkeepingProvider,
		initializer.MetricsProvider,
		initializer.HealthCheckRegistry,
		initializer.Config.StateDB,
	)
	if err != nil {
		return err
	}
	p.stats = newStats(initializer.MetricsProvider)
	p.recoverUnderConstructionLedger()
	return nil
}






func (p *Provider) Create(genesisBlock *common.Block) (ledger.PeerLedger, error) {
	ledgerID, err := protoutil.GetChainIDFromBlock(genesisBlock)
	if err != nil {
		return nil, err
	}
	exists, err := p.idStore.ledgerIDExists(ledgerID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrLedgerIDExists
	}
	if err = p.idStore.setUnderConstructionFlag(ledgerID); err != nil {
		return nil, err
	}
	lgr, err := p.openInternal(ledgerID)
	if err != nil {
		logger.Errorf("Error opening a new empty ledger. Unsetting under construction flag. Error: %+v", err)
		panicOnErr(p.runCleanup(ledgerID), "Error running cleanup for ledger id [%s]", ledgerID)
		panicOnErr(p.idStore.unsetUnderConstructionFlag(), "Error while unsetting under construction flag")
		return nil, err
	}
	if err := lgr.CommitWithPvtData(&ledger.BlockAndPvtData{
		Block: genesisBlock,
	}); err != nil {
		lgr.Close()
		return nil, err
	}
	panicOnErr(p.idStore.createLedgerID(ledgerID, genesisBlock), "Error while marking ledger as created")
	return lgr, nil
}


func (p *Provider) Open(ledgerID string) (ledger.PeerLedger, error) {
	logger.Debugf("Open() opening kvledger: %s", ledgerID)
	
	exists, err := p.idStore.ledgerIDExists(ledgerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNonExistingLedgerID
	}
	return p.openInternal(ledgerID)
}

func (p *Provider) openInternal(ledgerID string) (ledger.PeerLedger, error) {
	
	blockStore, err := p.ledgerStoreProvider.Open(ledgerID)
	if err != nil {
		return nil, err
	}
	p.collElgNotifier.registerListener(ledgerID, blockStore)

	
	vDB, err := p.vdbProvider.GetDBHandle(ledgerID)
	if err != nil {
		return nil, err
	}

	
	historyDB, err := p.historydbProvider.GetDBHandle(ledgerID)
	if err != nil {
		return nil, err
	}

	
	
	l, err := newKVLedger(
		ledgerID, blockStore, vDB, historyDB, p.configHistoryMgr,
		p.stateListeners, p.bookkeepingProvider,
		p.initializer.DeployedChaincodeInfoProvider,
		p.stats.ledgerStats(ledgerID),
	)
	if err != nil {
		return nil, err
	}
	return l, nil
}


func (p *Provider) Exists(ledgerID string) (bool, error) {
	return p.idStore.ledgerIDExists(ledgerID)
}


func (p *Provider) List() ([]string, error) {
	return p.idStore.getAllLedgerIds()
}


func (p *Provider) Close() {
	p.idStore.close()
	p.ledgerStoreProvider.Close()
	p.vdbProvider.Close()
	p.historydbProvider.Close()
	p.bookkeepingProvider.Close()
	p.configHistoryMgr.Close()
}





func (p *Provider) recoverUnderConstructionLedger() {
	logger.Debugf("Recovering under construction ledger")
	ledgerID, err := p.idStore.getUnderConstructionFlag()
	panicOnErr(err, "Error while checking whether the under construction flag is set")
	if ledgerID == "" {
		logger.Debugf("No under construction ledger found. Quitting recovery")
		return
	}
	logger.Infof("ledger [%s] found as under construction", ledgerID)
	ledger, err := p.openInternal(ledgerID)
	panicOnErr(err, "Error while opening under construction ledger [%s]", ledgerID)
	bcInfo, err := ledger.GetBlockchainInfo()
	panicOnErr(err, "Error while getting blockchain info for the under construction ledger [%s]", ledgerID)
	ledger.Close()

	switch bcInfo.Height {
	case 0:
		logger.Infof("Genesis block was not committed. Hence, the peer ledger not created. unsetting the under construction flag")
		panicOnErr(p.runCleanup(ledgerID), "Error while running cleanup for ledger id [%s]", ledgerID)
		panicOnErr(p.idStore.unsetUnderConstructionFlag(), "Error while unsetting under construction flag")
	case 1:
		logger.Infof("Genesis block was committed. Hence, marking the peer ledger as created")
		genesisBlock, err := ledger.GetBlockByNumber(0)
		panicOnErr(err, "Error while retrieving genesis block from blockchain for ledger [%s]", ledgerID)
		panicOnErr(p.idStore.createLedgerID(ledgerID, genesisBlock), "Error while adding ledgerID [%s] to created list", ledgerID)
	default:
		panic(errors.Errorf(
			"data inconsistency: under construction flag is set for ledger [%s] while the height of the blockchain is [%d]",
			ledgerID, bcInfo.Height))
	}
	return
}



func (p *Provider) runCleanup(ledgerID string) error {
	
	
	
	
	
	return nil
}

func panicOnErr(err error, mgsFormat string, args ...interface{}) {
	if err == nil {
		return
	}
	args = append(args, err)
	panic(fmt.Sprintf(mgsFormat+" Error: %s", args...))
}




type idStore struct {
	db *leveldbhelper.DB
}

func openIDStore(path string) *idStore {
	db := leveldbhelper.CreateDB(&leveldbhelper.Conf{DBPath: path})
	db.Open()
	return &idStore{db}
}

func (s *idStore) setUnderConstructionFlag(ledgerID string) error {
	return s.db.Put(underConstructionLedgerKey, []byte(ledgerID), true)
}

func (s *idStore) unsetUnderConstructionFlag() error {
	return s.db.Delete(underConstructionLedgerKey, true)
}

func (s *idStore) getUnderConstructionFlag() (string, error) {
	val, err := s.db.Get(underConstructionLedgerKey)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (s *idStore) createLedgerID(ledgerID string, gb *common.Block) error {
	key := s.encodeLedgerKey(ledgerID)
	var val []byte
	var err error
	if val, err = s.db.Get(key); err != nil {
		return err
	}
	if val != nil {
		return ErrLedgerIDExists
	}
	if val, err = proto.Marshal(gb); err != nil {
		return err
	}
	batch := &leveldb.Batch{}
	batch.Put(key, val)
	batch.Delete(underConstructionLedgerKey)
	return s.db.WriteBatch(batch, true)
}

func (s *idStore) ledgerIDExists(ledgerID string) (bool, error) {
	key := s.encodeLedgerKey(ledgerID)
	val := []byte{}
	err := error(nil)
	if val, err = s.db.Get(key); err != nil {
		return false, err
	}
	return val != nil, nil
}

func (s *idStore) getAllLedgerIds() ([]string, error) {
	var ids []string
	itr := s.db.GetIterator(nil, nil)
	defer itr.Release()
	itr.First()
	for itr.Valid() {
		if bytes.Equal(itr.Key(), underConstructionLedgerKey) {
			continue
		}
		id := string(s.decodeLedgerID(itr.Key()))
		ids = append(ids, id)
		itr.Next()
	}
	return ids, nil
}

func (s *idStore) close() {
	s.db.Close()
}

func (s *idStore) encodeLedgerKey(ledgerID string) []byte {
	return append(ledgerKeyPrefix, []byte(ledgerID)...)
}

func (s *idStore) decodeLedgerID(key []byte) string {
	return string(key[len(ledgerKeyPrefix):])
}
