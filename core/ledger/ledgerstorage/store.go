/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerstorage

import (
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage/fsblkstorage"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatastorage"
	"github.com/mcc-github/blockchain/protos/common"
)

const maxBlockFileSize = 64 * 1024 * 1024

var logger = flogging.MustGetLogger("ledgerstorage")


type Provider struct {
	blkStoreProvider     blkstorage.BlockStoreProvider
	pvtdataStoreProvider pvtdatastorage.Provider
}


type Store struct {
	blkstorage.BlockStore
	pvtdataStore                pvtdatastorage.Store
	rwlock                      sync.RWMutex
	isPvtstoreAheadOfBlockstore atomic.Value
}

var attrsToIndex = []blkstorage.IndexableAttr{
	blkstorage.IndexableAttrBlockHash,
	blkstorage.IndexableAttrBlockNum,
	blkstorage.IndexableAttrTxID,
	blkstorage.IndexableAttrBlockNumTranNum,
	
	blkstorage.IndexableAttrBlockTxID,
	blkstorage.IndexableAttrTxValidationCode,
}


func NewProvider(storeDir string, conf *pvtdatastorage.PrivateDataConfig, metricsProvider metrics.Provider) *Provider {
	
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	blockStoreProvider := fsblkstorage.NewProvider(
		fsblkstorage.NewConf(
			filepath.Join(storeDir, "chains"),
			maxBlockFileSize,
		),
		indexConfig,
		metricsProvider,
	)

	pvtStoreProvider := pvtdatastorage.NewProvider(conf)
	return &Provider{blockStoreProvider, pvtStoreProvider}
}


func (p *Provider) Open(ledgerid string) (*Store, error) {
	var blockStore blkstorage.BlockStore
	var pvtdataStore pvtdatastorage.Store
	var err error

	if blockStore, err = p.blkStoreProvider.OpenBlockStore(ledgerid); err != nil {
		return nil, err
	}
	if pvtdataStore, err = p.pvtdataStoreProvider.OpenStore(ledgerid); err != nil {
		return nil, err
	}
	store := &Store{
		BlockStore:   blockStore,
		pvtdataStore: pvtdataStore,
	}
	if err := store.init(); err != nil {
		return nil, err
	}

	info, err := blockStore.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}
	pvtstoreHeight, err := pvtdataStore.LastCommittedBlockHeight()
	if err != nil {
		return nil, err
	}
	store.isPvtstoreAheadOfBlockstore.Store(pvtstoreHeight > info.Height)

	return store, nil
}


func (p *Provider) Close() {
	p.blkStoreProvider.Close()
	p.pvtdataStoreProvider.Close()
}


func (p *Provider) Exists(ledgerID string) (bool, error) {
	return p.blkStoreProvider.Exists(ledgerID)
}


func (s *Store) Init(btlPolicy pvtdatapolicy.BTLPolicy) {
	s.pvtdataStore.Init(btlPolicy)
}


func (s *Store) CommitWithPvtData(blockAndPvtdata *ledger.BlockAndPvtData) error {
	blockNum := blockAndPvtdata.Block.Header.Number
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	pvtBlkStoreHt, err := s.pvtdataStore.LastCommittedBlockHeight()
	if err != nil {
		return err
	}

	if pvtBlkStoreHt < blockNum+1 { 
		
		
		logger.Debugf("Writing block [%d] to pvt block store", blockNum)
		
		
		
		
		
		
		
		
		pvtData, missingPvtData := constructPvtDataAndMissingData(blockAndPvtdata)
		if err := s.pvtdataStore.Commit(blockAndPvtdata.Block.Header.Number, pvtData, missingPvtData); err != nil {
			return err
		}
	} else {
		logger.Debugf("Skipping writing block [%d] to pvt block store as the store height is [%d]", blockNum, pvtBlkStoreHt)
	}

	if err := s.AddBlock(blockAndPvtdata.Block); err != nil {
		return err
	}

	if pvtBlkStoreHt == blockNum+1 {
		
		
		
		s.isPvtstoreAheadOfBlockstore.Store(false)
	}

	return nil
}

func constructPvtDataAndMissingData(blockAndPvtData *ledger.BlockAndPvtData) ([]*ledger.TxPvtData,
	ledger.TxMissingPvtDataMap) {

	var pvtData []*ledger.TxPvtData
	missingPvtData := make(ledger.TxMissingPvtDataMap)

	numTxs := uint64(len(blockAndPvtData.Block.Data.Data))

	
	for txNum := uint64(0); txNum < numTxs; txNum++ {
		if pvtdata, ok := blockAndPvtData.PvtData[txNum]; ok {
			pvtData = append(pvtData, pvtdata)
		}

		if missingData, ok := blockAndPvtData.MissingPvtData[txNum]; ok {
			for _, missing := range missingData {
				missingPvtData.Add(txNum, missing.Namespace,
					missing.Collection, missing.IsEligible)
			}
		}
	}
	return pvtData, missingPvtData
}


func (s *Store) CommitPvtDataOfOldBlocks(blocksPvtData map[uint64][]*ledger.TxPvtData) error {
	err := s.pvtdataStore.CommitPvtDataOfOldBlocks(blocksPvtData)
	if err != nil {
		return err
	}
	return nil
}



func (s *Store) GetPvtDataAndBlockByNum(blockNum uint64, filter ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error) {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()

	var block *common.Block
	var pvtdata []*ledger.TxPvtData
	var err error
	if block, err = s.RetrieveBlockByNumber(blockNum); err != nil {
		return nil, err
	}
	if pvtdata, err = s.getPvtDataByNumWithoutLock(blockNum, filter); err != nil {
		return nil, err
	}
	return &ledger.BlockAndPvtData{Block: block, PvtData: constructPvtdataMap(pvtdata)}, nil
}




func (s *Store) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	return s.getPvtDataByNumWithoutLock(blockNum, filter)
}




func (s *Store) getPvtDataByNumWithoutLock(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	var pvtdata []*ledger.TxPvtData
	var err error
	if pvtdata, err = s.pvtdataStore.GetPvtDataByBlockNum(blockNum, filter); err != nil {
		return nil, err
	}
	return pvtdata, nil
}






func (s *Store) DoesPvtDataInfoExist(blockNum uint64) (bool, error) {
	pvtStoreHt, err := s.pvtdataStore.LastCommittedBlockHeight()
	if err != nil {
		return false, err
	}
	return blockNum+1 <= pvtStoreHt, nil
}


func (s *Store) GetMissingPvtDataInfoForMostRecentBlocks(maxBlock int) (ledger.MissingPvtDataInfo, error) {
	
	
	
	
	return s.pvtdataStore.GetMissingPvtDataInfoForMostRecentBlocks(maxBlock)
}


func (s *Store) ProcessCollsEligibilityEnabled(committingBlk uint64, nsCollMap map[string][]string) error {
	return s.pvtdataStore.ProcessCollsEligibilityEnabled(committingBlk, nsCollMap)
}


func (s *Store) GetLastUpdatedOldBlocksPvtData() (map[uint64][]*ledger.TxPvtData, error) {
	return s.pvtdataStore.GetLastUpdatedOldBlocksPvtData()
}


func (s *Store) ResetLastUpdatedOldBlocksList() error {
	return s.pvtdataStore.ResetLastUpdatedOldBlocksList()
}



func (s *Store) IsPvtStoreAheadOfBlockStore() bool {
	return s.isPvtstoreAheadOfBlockstore.Load().(bool)
}













func (s *Store) init() error {
	var initialized bool
	var err error
	if initialized, err = s.initPvtdataStoreFromExistingBlockchain(); err != nil || initialized {
		return err
	}
	return nil
}








func (s *Store) initPvtdataStoreFromExistingBlockchain() (bool, error) {
	var bcInfo *common.BlockchainInfo
	var pvtdataStoreEmpty bool
	var err error

	if bcInfo, err = s.BlockStore.GetBlockchainInfo(); err != nil {
		return false, err
	}
	if pvtdataStoreEmpty, err = s.pvtdataStore.IsEmpty(); err != nil {
		return false, err
	}
	if pvtdataStoreEmpty && bcInfo.Height > 0 {
		if err = s.pvtdataStore.InitLastCommittedBlock(bcInfo.Height - 1); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func constructPvtdataMap(pvtdata []*ledger.TxPvtData) ledger.TxPvtDataMap {
	if pvtdata == nil {
		return nil
	}
	m := make(map[uint64]*ledger.TxPvtData)
	for _, pvtdatum := range pvtdata {
		m[pvtdatum.SeqInBlock] = pvtdatum
	}
	return m
}


func LoadPreResetHeight(blockstorePath string) (map[string]uint64, error) {
	return fsblkstorage.LoadPreResetHeight(blockstorePath)
}


func ResetBlockStore(blockstorePath string) error {
	return fsblkstorage.ResetBlockStore(blockstorePath)
}



func ValidateRollbackParams(blockstorePath, ledgerID string, blockNum uint64) error {
	return fsblkstorage.ValidateRollbackParams(blockstorePath, ledgerID, blockNum)
}


func Rollback(blockstorePath, ledgerID string, blockNum uint64) error {
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	return fsblkstorage.Rollback(blockstorePath, ledgerID, blockNum, indexConfig)
}
