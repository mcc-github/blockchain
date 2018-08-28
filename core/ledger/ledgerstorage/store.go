/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerstorage

import (
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/blkstorage/fsblkstorage"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatastorage"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("ledgerstorage")


type Provider struct {
	blkStoreProvider     blkstorage.BlockStoreProvider
	pvtdataStoreProvider pvtdatastorage.Provider
}


type Store struct {
	blkstorage.BlockStore
	pvtdataStore pvtdatastorage.Store
	rwlock       *sync.RWMutex
}


func NewProvider() *Provider {
	
	attrsToIndex := []blkstorage.IndexableAttr{
		blkstorage.IndexableAttrBlockHash,
		blkstorage.IndexableAttrBlockNum,
		blkstorage.IndexableAttrTxID,
		blkstorage.IndexableAttrBlockNumTranNum,
		blkstorage.IndexableAttrBlockTxID,
		blkstorage.IndexableAttrTxValidationCode,
	}
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	blockStoreProvider := fsblkstorage.NewProvider(
		fsblkstorage.NewConf(ledgerconfig.GetBlockStorePath(), ledgerconfig.GetMaxBlockfileSize()),
		indexConfig)

	pvtStoreProvider := pvtdatastorage.NewProvider()
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
	store := &Store{blockStore, pvtdataStore, &sync.RWMutex{}}
	if err := store.init(); err != nil {
		return nil, err
	}
	return store, nil
}


func (p *Provider) Close() {
	p.blkStoreProvider.Close()
	p.pvtdataStoreProvider.Close()
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

	writtenToPvtStore := false
	if pvtBlkStoreHt < blockNum+1 { 
		
		
		logger.Debugf("Writing block [%d] to pvt block store", blockNum)
		var pvtdata []*ledger.TxPvtData
		for _, v := range blockAndPvtdata.BlockPvtData {
			pvtdata = append(pvtdata, v)
		}
		if err := s.pvtdataStore.Prepare(blockAndPvtdata.Block.Header.Number, pvtdata, blockAndPvtdata.Missing); err != nil {
			return err
		}
		writtenToPvtStore = true
	} else {
		logger.Debugf("Skipping writing block [%d] to pvt block store as the store height is [%d]", blockNum, pvtBlkStoreHt)
	}

	if err := s.AddBlock(blockAndPvtdata.Block); err != nil {
		s.pvtdataStore.Rollback()
		return err
	}

	if writtenToPvtStore {
		return s.pvtdataStore.Commit()
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
	return &ledger.BlockAndPvtData{Block: block, BlockPvtData: constructPvtdataMap(pvtdata)}, nil
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

func (s *Store) GetMissingPvtDataInfoForMostRecentBlocks(maxBlock int) (ledger.MissingPvtDataInfo, error) {
	
	
	
	
	return s.pvtdataStore.GetMissingPvtDataInfoForMostRecentBlocks(maxBlock)
}






func (s *Store) init() error {
	var initialized bool
	var err error
	if initialized, err = s.initPvtdataStoreFromExistingBlockchain(); err != nil || initialized {
		return err
	}
	return s.syncPvtdataStoreWithBlockStore()
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








func (s *Store) syncPvtdataStoreWithBlockStore() error {
	var pendingPvtbatch bool
	var err error
	if pendingPvtbatch, err = s.pvtdataStore.HasPendingBatch(); err != nil {
		return err
	}
	if !pendingPvtbatch {
		return nil
	}
	var bcInfo *common.BlockchainInfo
	var pvtdataStoreHt uint64

	if bcInfo, err = s.GetBlockchainInfo(); err != nil {
		return err
	}
	if pvtdataStoreHt, err = s.pvtdataStore.LastCommittedBlockHeight(); err != nil {
		return err
	}

	if bcInfo.Height == pvtdataStoreHt {
		return s.pvtdataStore.Rollback()
	}

	if bcInfo.Height == pvtdataStoreHt+1 {
		return s.pvtdataStore.Commit()
	}

	return errors.Errorf("This is not expected. blockStoreHeight=%d, pvtdataStoreHeight=%d", bcInfo.Height, pvtdataStoreHt)
}

func constructPvtdataMap(pvtdata []*ledger.TxPvtData) map[uint64]*ledger.TxPvtData {
	if pvtdata == nil {
		return nil
	}
	m := make(map[uint64]*ledger.TxPvtData)
	for _, pvtdatum := range pvtdata {
		m[pvtdatum.SeqInBlock] = pvtdatum
	}
	return m
}
