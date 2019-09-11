/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/confighistory"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/history"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/msgs"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/ledgerstorage"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatastorage"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

const idStoreFormatVersion = "2.0"

var (
	
	ErrLedgerIDExists = errors.New("LedgerID already exists")
	
	ErrNonExistingLedgerID = errors.New("LedgerID does not exist")
	
	ErrLedgerNotOpened = errors.New("ledger is not opened yet")
	
	ErrInactiveLedger = errors.New("Ledger is not active")

	underConstructionLedgerKey = []byte("underConstructionLedgerKey")
	
	ledgerKeyPrefix = []byte{'l'}
	
	ledgerKeyStop = []byte{'l' + 1}
	
	metadataKeyPrefix = []byte{'s'}
	
	metadataKeyStop = []byte{'s' + 1}

	
	formatKey = []byte("f")
)


type Provider struct {
	idStore             *idStore
	ledgerStoreProvider *ledgerstorage.Provider
	vdbProvider         privacyenabledstate.DBProvider
	historydbProvider   *history.DBProvider
	configHistoryMgr    confighistory.Mgr
	stateListeners      []ledger.StateListener
	bookkeepingProvider bookkeeping.Provider
	initializer         *ledger.Initializer
	collElgNotifier     *collElgNotifier
	stats               *stats
	fileLock            *leveldbhelper.FileLock
}



func NewProvider(initializer *ledger.Initializer) (pr *Provider, e error) {
	p := &Provider{
		initializer: initializer,
	}

	defer func() {
		if e != nil {
			p.Close()
		}
	}()

	fileLockPath := fileLockPath(initializer.Config.RootFSPath)
	fileLock := leveldbhelper.NewFileLock(fileLockPath)
	if err := fileLock.Lock(); err != nil {
		return nil, errors.Wrap(err, "as another peer node command is executing,"+
			" wait for that command to complete its execution or terminate it before retrying")
	}

	p.fileLock = fileLock
	
	idStore, err := openIDStore(LedgerProviderPath(p.initializer.Config.RootFSPath))
	if err != nil {
		return nil, err
	}
	p.idStore = idStore
	
	privateData := &pvtdatastorage.PrivateDataConfig{
		PrivateDataConfig: initializer.Config.PrivateDataConfig,
		StorePath:         PvtDataStorePath(p.initializer.Config.RootFSPath),
	}

	ledgerStoreProvider, err := ledgerstorage.NewProvider(
		BlockStorePath(p.initializer.Config.RootFSPath),
		privateData,
		p.initializer.MetricsProvider,
	)
	if err != nil {
		return nil, err
	}
	p.ledgerStoreProvider = ledgerStoreProvider
	if initializer.Config.HistoryDBConfig.Enabled {
		
		historydbProvider, err := history.NewDBProvider(
			HistoryDBPath(p.initializer.Config.RootFSPath),
		)
		if err != nil {
			return nil, err
		}
		p.historydbProvider = historydbProvider
	}
	
	configHistoryMgr, err := confighistory.NewMgr(
		ConfigHistoryDBPath(p.initializer.Config.RootFSPath),
		initializer.DeployedChaincodeInfoProvider,
	)
	if err != nil {
		return nil, err
	}
	p.configHistoryMgr = configHistoryMgr
	

	collElgNotifier := &collElgNotifier{
		initializer.DeployedChaincodeInfoProvider,
		initializer.MembershipInfoProvider,
		make(map[string]collElgListener),
	}
	p.collElgNotifier = collElgNotifier
	
	stateListeners := initializer.StateListeners
	stateListeners = append(stateListeners, collElgNotifier)
	stateListeners = append(stateListeners, configHistoryMgr)
	p.stateListeners = stateListeners

	p.bookkeepingProvider, err = bookkeeping.NewProvider(
		BookkeeperDBPath(p.initializer.Config.RootFSPath),
	)
	if err != nil {
		return nil, err
	}
	stateDB := &privacyenabledstate.StateDBConfig{
		StateDBConfig: initializer.Config.StateDBConfig,
		LevelDBPath:   StateDBPath(p.initializer.Config.RootFSPath),
	}
	p.vdbProvider, err = privacyenabledstate.NewCommonStorageDBProvider(
		p.bookkeepingProvider,
		initializer.MetricsProvider,
		initializer.HealthCheckRegistry,
		stateDB,
	)
	if err != nil {
		return nil, err
	}
	p.stats = newStats(initializer.MetricsProvider)
	p.recoverUnderConstructionLedger()
	return p, nil
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
	if err := lgr.CommitLegacy(&ledger.BlockAndPvtData{Block: genesisBlock}, &ledger.CommitOptions{}); err != nil {
		lgr.Close()
		return nil, err
	}
	panicOnErr(p.idStore.createLedgerID(ledgerID, genesisBlock), "Error while marking ledger as created")
	return lgr, nil
}


func (p *Provider) Open(ledgerID string) (ledger.PeerLedger, error) {
	logger.Debugf("Open() opening kvledger: %s", ledgerID)
	
	active, exists, err := p.idStore.ledgerIDActive(ledgerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNonExistingLedgerID
	}
	if !active {
		return nil, ErrInactiveLedger
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

	
	var historyDB *history.DB
	if p.historydbProvider != nil {
		historyDB, err = p.historydbProvider.GetDBHandle(ledgerID)
		if err != nil {
			return nil, err
		}
	}

	
	
	l, err := newKVLedger(
		ledgerID,
		blockStore,
		vDB,
		historyDB,
		p.configHistoryMgr,
		p.stateListeners,
		p.bookkeepingProvider,
		p.initializer.DeployedChaincodeInfoProvider,
		p.initializer.ChaincodeLifecycleEventProvider,
		p.stats.ledgerStats(ledgerID),
		p.initializer.CustomTxProcessors,
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
	return p.idStore.getActiveLedgerIDs()
}


func (p *Provider) Close() {
	if p.idStore != nil {
		p.idStore.close()
	}
	if p.ledgerStoreProvider != nil {
		p.ledgerStoreProvider.Close()
	}
	if p.vdbProvider != nil {
		p.vdbProvider.Close()
	}
	if p.bookkeepingProvider != nil {
		p.bookkeepingProvider.Close()
	}
	if p.configHistoryMgr != nil {
		p.configHistoryMgr.Close()
	}
	if p.historydbProvider != nil {
		p.historydbProvider.Close()
	}
	if p.fileLock != nil {
		p.fileLock.Unlock()
	}
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
	db     *leveldbhelper.DB
	dbPath string
}

func openIDStore(path string) (s *idStore, e error) {
	db := leveldbhelper.CreateDB(&leveldbhelper.Conf{DBPath: path})
	db.Open()
	defer func() {
		if e != nil {
			db.Close()
		}
	}()

	emptyDB, err := db.IsEmpty()
	if err != nil {
		return nil, err
	}

	if emptyDB {
		
		err := db.Put(formatKey, []byte(idStoreFormatVersion), true)
		if err != nil {
			return nil, err
		}
		return &idStore{db, path}, nil
	}

	
	formatVersion, err := db.Get(formatKey)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(formatVersion, []byte(idStoreFormatVersion)) {
		return nil, &leveldbhelper.ErrFormatVersionMismatch{ExpectedFormatVersion: idStoreFormatVersion, DataFormatVersion: string(formatVersion), DBPath: path}
	}
	return &idStore{db, path}, nil
}

func (s *idStore) upgradeFormat() error {
	format, err := s.db.Get(formatKey)
	if err != nil {
		return err
	}
	idStoreFormatBytes := []byte(idStoreFormatVersion)
	if bytes.Equal(format, idStoreFormatBytes) {
		logger.Debug("Format is current, nothing to do")
		return nil
	}
	if format != nil {
		err = &leveldbhelper.ErrFormatVersionMismatch{ExpectedFormatVersion: "", DataFormatVersion: string(format), DBPath: s.dbPath}
		logger.Errorf("Failed to upgrade format [%#v] to new format [%#v]: %s", format, idStoreFormatBytes, err)
		return err
	}

	logger.Infof("The ledgerProvider db format is old, upgrading to the new format %s", idStoreFormatVersion)

	batch := &leveldb.Batch{}
	batch.Put(formatKey, idStoreFormatBytes)

	
	metadata, err := protoutil.Marshal(&msgs.LedgerMetadata{Status: msgs.Status_ACTIVE})
	if err != nil {
		logger.Errorf("Error marshalling ledger metadata: %s", err)
		return errors.Wrapf(err, "error marshalling ledger metadata")
	}
	itr := s.db.GetIterator(ledgerKeyPrefix, ledgerKeyStop)
	defer itr.Release()
	for itr.Error() == nil && itr.Next() {
		id := s.decodeLedgerID(itr.Key(), ledgerKeyPrefix)
		batch.Put(s.encodeLedgerKey(id, metadataKeyPrefix), metadata)
	}
	if err = itr.Error(); err != nil {
		logger.Errorf("Error while upgrading idStore format: %s", err)
		return errors.Wrapf(err, "error while upgrading idStore format")
	}

	return s.db.WriteBatch(batch, true)
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
	gbKey := s.encodeLedgerKey(ledgerID, ledgerKeyPrefix)
	metadataKey := s.encodeLedgerKey(ledgerID, metadataKeyPrefix)
	var val []byte
	var metadata []byte
	var err error
	if val, err = s.db.Get(gbKey); err != nil {
		return err
	}
	if val != nil {
		return ErrLedgerIDExists
	}
	if val, err = proto.Marshal(gb); err != nil {
		return err
	}
	if metadata, err = protoutil.Marshal(&msgs.LedgerMetadata{Status: msgs.Status_ACTIVE}); err != nil {
		return err
	}
	batch := &leveldb.Batch{}
	batch.Put(gbKey, val)
	batch.Put(metadataKey, metadata)
	batch.Delete(underConstructionLedgerKey)
	return s.db.WriteBatch(batch, true)
}

func (s *idStore) ledgerIDExists(ledgerID string) (bool, error) {
	key := s.encodeLedgerKey(ledgerID, ledgerKeyPrefix)
	val := []byte{}
	err := error(nil)
	if val, err = s.db.Get(key); err != nil {
		return false, err
	}
	return val != nil, nil
}


func (s *idStore) ledgerIDActive(ledgerID string) (bool, bool, error) {
	key := s.encodeLedgerKey(ledgerID, metadataKeyPrefix)
	val, err := s.db.Get(key)
	if val == nil || err != nil {
		return false, false, err
	}
	metadata := &msgs.LedgerMetadata{}
	if err = proto.Unmarshal(val, metadata); err != nil {
		logger.Errorf("Error unmarshing ledger metadata: %s", err)
		return false, false, errors.Wrapf(err, "error unmarshing ledger metadata")
	}
	return metadata.Status == msgs.Status_ACTIVE, true, nil
}

func (s *idStore) getActiveLedgerIDs() ([]string, error) {
	var ids []string
	itr := s.db.GetIterator(metadataKeyPrefix, metadataKeyStop)
	defer itr.Release()
	for itr.Error() == nil && itr.Next() {
		metadata := &msgs.LedgerMetadata{}
		if err := proto.Unmarshal(itr.Value(), metadata); err != nil {
			logger.Errorf("Error unmarshing ledger metadata: %s", err)
			return nil, errors.Wrapf(err, "error unmarshing ledger metadata")
		}
		if metadata.Status == msgs.Status_ACTIVE {
			id := s.decodeLedgerID(itr.Key(), metadataKeyPrefix)
			ids = append(ids, id)
		}
	}
	if err := itr.Error(); err != nil {
		logger.Errorf("Error getting ledger ids from idStore: %s", err)
		return nil, errors.Wrapf(err, "error getting ledger ids from idStore")
	}
	return ids, nil
}

func (s *idStore) close() {
	s.db.Close()
}

func (s *idStore) encodeLedgerKey(ledgerID string, prefix []byte) []byte {
	return append(prefix, []byte(ledgerID)...)
}

func (s *idStore) decodeLedgerID(key []byte, prefix []byte) string {
	return string(key[len(prefix):])
}
