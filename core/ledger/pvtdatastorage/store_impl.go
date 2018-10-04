/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/willf/bitset"
)

var logger = flogging.MustGetLogger("pvtdatastorage")

type provider struct {
	dbProvider *leveldbhelper.Provider
}

type store struct {
	db        *leveldbhelper.DBHandle
	ledgerid  string
	btlPolicy pvtdatapolicy.BTLPolicy

	isEmpty            bool
	lastCommittedBlock uint64
	batchPending       bool
	purgerLock         sync.Mutex
	collElgProcSync    *collElgProcSync
}

type blkTranNumKey []byte

type dataEntry struct {
	key   *dataKey
	value *rwset.CollectionPvtReadWriteSet
}

type expiryEntry struct {
	key   *expiryKey
	value *ExpiryData
}

type expiryKey struct {
	expiringBlk   uint64
	committingBlk uint64
}

type nsCollBlk struct {
	ns, coll string
	blkNum   uint64
}

type dataKey struct {
	nsCollBlk
	txNum uint64
}

type missingDataKey struct {
	nsCollBlk
	isEligible bool
}

type storeEntries struct {
	dataEntries        []*dataEntry
	expiryEntries      []*expiryEntry
	missingDataEntries map[missingDataKey]*bitset.BitSet
}





func NewProvider() Provider {
	dbPath := ledgerconfig.GetPvtdataStorePath()
	dbProvider := leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: dbPath})
	return &provider{dbProvider: dbProvider}
}


func (p *provider) OpenStore(ledgerid string) (Store, error) {
	dbHandle := p.dbProvider.GetDBHandle(ledgerid)
	s := &store{db: dbHandle, ledgerid: ledgerid,
		collElgProcSync: &collElgProcSync{
			notification: make(chan bool, 1),
			procComplete: make(chan bool, 1),
		},
	}
	if err := s.initState(); err != nil {
		return nil, err
	}
	s.launchCollElgProc()
	logger.Debugf("Pvtdata store opened. Initial state: isEmpty [%t], lastCommittedBlock [%d], batchPending [%t]",
		s.isEmpty, s.lastCommittedBlock, s.batchPending)
	return s, nil
}


func (p *provider) Close() {
	p.dbProvider.Close()
}




func (s *store) initState() error {
	var err error
	if s.isEmpty, s.lastCommittedBlock, err = s.getLastCommittedBlockNum(); err != nil {
		return err
	}
	if s.batchPending, err = s.hasPendingCommit(); err != nil {
		return err
	}

	return nil
}

func (s *store) Init(btlPolicy pvtdatapolicy.BTLPolicy) {
	s.btlPolicy = btlPolicy
}


func (s *store) Prepare(blockNum uint64, pvtData []*ledger.TxPvtData, missingData *ledger.MissingPrivateDataList) error {
	if s.batchPending {
		return &ErrIllegalCall{`A pending batch exists as as result of last invoke to "Prepare" call.
			 Invoke "Commit" or "Rollback" on the pending batch before invoking "Prepare" function`}
	}
	expectedBlockNum := s.nextBlockNum()
	if expectedBlockNum != blockNum {
		return &ErrIllegalArgs{fmt.Sprintf("Expected block number=%d, recived block number=%d", expectedBlockNum, blockNum)}
	}

	batch := leveldbhelper.NewUpdateBatch()
	var err error
	var keyBytes, valBytes []byte

	storeEntries, err := prepareStoreEntries(blockNum, pvtData, s.btlPolicy, missingData)

	if err != nil {
		return err
	}

	for _, dataEntry := range storeEntries.dataEntries {
		keyBytes = encodeDataKey(dataEntry.key)
		if valBytes, err = encodeDataValue(dataEntry.value); err != nil {
			return err
		}
		batch.Put(keyBytes, valBytes)
	}

	for _, expiryEntry := range storeEntries.expiryEntries {
		keyBytes = encodeExpiryKey(expiryEntry.key)
		if valBytes, err = encodeExpiryValue(expiryEntry.value); err != nil {
			return err
		}
		batch.Put(keyBytes, valBytes)
	}

	for missingDataKey, missingDataValue := range storeEntries.missingDataEntries {
		keyBytes = encodeMissingDataKey(&missingDataKey)
		if valBytes, err = encodeMissingDataValue(missingDataValue); err != nil {
			return err
		}
		batch.Put(keyBytes, valBytes)
	}

	batch.Put(pendingCommitKey, emptyValue)
	if err := s.db.WriteBatch(batch, true); err != nil {
		return err
	}
	s.batchPending = true
	logger.Debugf("Saved %d private data write sets for block [%d]", len(pvtData), blockNum)
	return nil
}


func (s *store) Commit() error {
	if !s.batchPending {
		return &ErrIllegalCall{"No pending batch to commit"}
	}
	committingBlockNum := s.nextBlockNum()
	logger.Debugf("Committing private data for block [%d]", committingBlockNum)
	batch := leveldbhelper.NewUpdateBatch()
	batch.Delete(pendingCommitKey)
	batch.Put(lastCommittedBlkkey, encodeLastCommittedBlockVal(committingBlockNum))
	if err := s.db.WriteBatch(batch, true); err != nil {
		return err
	}
	s.batchPending = false
	s.isEmpty = false
	s.lastCommittedBlock = committingBlockNum
	logger.Debugf("Committed private data for block [%d]", committingBlockNum)
	s.performPurgeIfScheduled(committingBlockNum)
	return nil
}




func (s *store) Rollback() error {
	if !s.batchPending {
		return &ErrIllegalCall{"No pending batch to rollback"}
	}
	s.batchPending = false
	return nil
}




func (s *store) GetPvtDataByBlockNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	logger.Debugf("Get private data for block [%d], filter=%#v", blockNum, filter)
	if s.isEmpty {
		return nil, &ErrOutOfRange{"The store is empty"}
	}
	if blockNum > s.lastCommittedBlock {
		return nil, &ErrOutOfRange{fmt.Sprintf("Last committed block=%d, block requested=%d", s.lastCommittedBlock, blockNum)}
	}
	startKey, endKey := getDataKeysForRangeScanByBlockNum(blockNum)
	logger.Debugf("Querying private data storage for write sets using startKey=%#v, endKey=%#v", startKey, endKey)
	itr := s.db.GetIterator(startKey, endKey)
	defer itr.Release()

	var blockPvtdata []*ledger.TxPvtData
	var currentTxNum uint64
	var currentTxWsetAssember *txPvtdataAssembler
	firstItr := true

	for itr.Next() {
		dataKeyBytes := itr.Key()
		if v11Format(dataKeyBytes) {
			return v11RetrievePvtdata(itr, filter)
		}
		dataValueBytes := itr.Value()
		dataKey := decodeDatakey(dataKeyBytes)
		expired, err := isExpired(dataKey.nsCollBlk, s.btlPolicy, s.lastCommittedBlock)
		if err != nil {
			return nil, err
		}
		if expired || !passesFilter(dataKey, filter) {
			continue
		}
		dataValue, err := decodeDataValue(dataValueBytes)
		if err != nil {
			return nil, err
		}

		if firstItr {
			currentTxNum = dataKey.txNum
			currentTxWsetAssember = newTxPvtdataAssembler(blockNum, currentTxNum)
			firstItr = false
		}

		if dataKey.txNum != currentTxNum {
			blockPvtdata = append(blockPvtdata, currentTxWsetAssember.getTxPvtdata())
			currentTxNum = dataKey.txNum
			currentTxWsetAssember = newTxPvtdataAssembler(blockNum, currentTxNum)
		}
		currentTxWsetAssember.add(dataKey.ns, dataValue)
	}
	if currentTxWsetAssember != nil {
		blockPvtdata = append(blockPvtdata, currentTxWsetAssember.getTxPvtdata())
	}
	return blockPvtdata, nil
}


func (s *store) InitLastCommittedBlock(blockNum uint64) error {
	if !(s.isEmpty && !s.batchPending) {
		return &ErrIllegalCall{"The private data store is not empty. InitLastCommittedBlock() function call is not allowed"}
	}
	batch := leveldbhelper.NewUpdateBatch()
	batch.Put(lastCommittedBlkkey, encodeLastCommittedBlockVal(blockNum))
	if err := s.db.WriteBatch(batch, true); err != nil {
		return err
	}
	s.isEmpty = false
	s.lastCommittedBlock = blockNum
	logger.Debugf("InitLastCommittedBlock set to block [%d]", blockNum)
	return nil
}


func (s *store) GetMissingPvtDataInfoForMostRecentBlocks(maxBlock int) (ledger.MissingPvtDataInfo, error) {
	
	
	if maxBlock < 1 {
		return nil, nil
	}

	missingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	numberOfBlockProcessed := 0
	lastProcessedBlock := uint64(0)
	isMaxBlockLimitReached := false
	
	
	
	
	lastCommittedBlock := atomic.LoadUint64(&s.lastCommittedBlock)

	startKey, endKey := createRangeScanKeysForEligibleMissingDataEntries(lastCommittedBlock)
	dbItr := s.db.GetIterator(startKey, endKey)
	defer dbItr.Release()

	for dbItr.Next() {
		missingDataKeyBytes := dbItr.Key()
		missingDataKey := decodeMissingDataKey(missingDataKeyBytes)

		if isMaxBlockLimitReached && (missingDataKey.blkNum != lastProcessedBlock) {
			
			
			break
		}

		
		expired, err := isExpired(missingDataKey.nsCollBlk, s.btlPolicy, lastCommittedBlock)
		if err != nil {
			return nil, err
		}
		if expired {
			continue
		}

		
		
		
		if _, ok := missingPvtDataInfo[missingDataKey.blkNum]; !ok {
			numberOfBlockProcessed++
			if numberOfBlockProcessed == maxBlock {
				isMaxBlockLimitReached = true
				
				
				lastProcessedBlock = missingDataKey.blkNum
			}
		}

		valueBytes := dbItr.Value()
		bitmap, err := decodeMissingDataValue(valueBytes)
		if err != nil {
			return nil, err
		}

		
		for index, isSet := bitmap.NextSet(0); isSet; index, isSet = bitmap.NextSet(index + 1) {
			txNum := uint64(index)
			missingPvtDataInfo.Add(missingDataKey.blkNum, txNum, missingDataKey.ns, missingDataKey.coll)
		}
	}

	return missingPvtDataInfo, nil
}

func (s *store) ProcessCollsEligibilityEnabled(committingBlk uint64, nsCollMap map[string][]string) error {
	key := encodeCollElgKey(committingBlk)
	m := newCollElgInfo(nsCollMap)
	val, err := encodeCollElgVal(m)
	if err != nil {
		return err
	}
	batch := leveldbhelper.NewUpdateBatch()
	batch.Put(key, val)
	if err = s.db.WriteBatch(batch, true); err != nil {
		return err
	}
	s.collElgProcSync.notify()
	return nil
}

func (s *store) performPurgeIfScheduled(latestCommittedBlk uint64) {
	if latestCommittedBlk%ledgerconfig.GetPvtdataStorePurgeInterval() != 0 {
		return
	}
	go func() {
		s.purgerLock.Lock()
		logger.Debugf("Purger started: Purging expired private data till block number [%d]", latestCommittedBlk)
		defer s.purgerLock.Unlock()
		err := s.purgeExpiredData(0, latestCommittedBlk)
		if err != nil {
			logger.Warningf("Could not purge data from pvtdata store:%s", err)
		}
		logger.Debug("Purger finished")
	}()
}

func (s *store) purgeExpiredData(minBlkNum, maxBlkNum uint64) error {
	batch := leveldbhelper.NewUpdateBatch()
	expiryEntries, err := s.retrieveExpiryEntries(minBlkNum, maxBlkNum)
	if err != nil || len(expiryEntries) == 0 {
		return err
	}
	for _, expiryEntry := range expiryEntries {
		
		
		batch.Delete(encodeExpiryKey(expiryEntry.key))
		dataKeys, missingDataKeys := deriveKeys(expiryEntry)
		for _, dataKey := range dataKeys {
			batch.Delete(encodeDataKey(dataKey))
		}
		for _, missingDataKey := range missingDataKeys {
			batch.Delete(encodeMissingDataKey(missingDataKey))
		}
		s.db.WriteBatch(batch, false)
	}
	logger.Infof("[%s] - [%d] Entries purged from private data storage till block number [%d]", s.ledgerid, len(expiryEntries), maxBlkNum)
	return nil
}

func (s *store) retrieveExpiryEntries(minBlkNum, maxBlkNum uint64) ([]*expiryEntry, error) {
	startKey, endKey := getExpiryKeysForRangeScan(minBlkNum, maxBlkNum)
	logger.Debugf("retrieveExpiryEntries(): startKey=%#v, endKey=%#v", startKey, endKey)
	itr := s.db.GetIterator(startKey, endKey)
	defer itr.Release()

	var expiryEntries []*expiryEntry
	for itr.Next() {
		expiryKeyBytes := itr.Key()
		expiryValueBytes := itr.Value()
		expiryKey := decodeExpiryKey(expiryKeyBytes)
		expiryValue, err := decodeExpiryValue(expiryValueBytes)
		if err != nil {
			return nil, err
		}
		expiryEntries = append(expiryEntries, &expiryEntry{key: expiryKey, value: expiryValue})
	}
	return expiryEntries, nil
}

func (s *store) launchCollElgProc() {
	go func() {
		s.processCollElgEvents() 
		for {
			logger.Debugf("Waiting for collection eligibility event")
			s.collElgProcSync.waitForNotification()
			s.processCollElgEvents()
			s.collElgProcSync.done()
		}
	}()
}

func (s *store) processCollElgEvents() {
	logger.Debugf("Starting to process collection eligibility events")
	maxBatchSize := ledgerconfig.GetPvtdataStoreCollElgProcMaxDbBatchSize()
	batchesInterval := ledgerconfig.GetPvtdataStoreCollElgProcDbBatchesInterval()
	s.purgerLock.Lock()
	defer s.purgerLock.Unlock()
	collElgStartKey, collElgEndKey := createRangeScanKeysForCollElg()
	eventItr := s.db.GetIterator(collElgStartKey, collElgEndKey)
	defer eventItr.Release()
	batch := leveldbhelper.NewUpdateBatch()
	totalEntriesConverted := 0

	for eventItr.Next() {
		collElgKey, collElgVal := eventItr.Key(), eventItr.Value()
		blkNum := decodeCollElgKey(collElgKey)
		CollElgInfo, err := decodeCollElgVal(collElgVal)
		logger.Debugf("Processing collection eligibility event [blkNum=%d], CollElgInfo=%s", blkNum, CollElgInfo)
		if err != nil {
			logger.Errorf("This error is not expected %s", err)
			continue
		}
		for ns, colls := range CollElgInfo.NsCollMap {
			var coll string
			for _, coll = range colls.Entries {
				logger.Infof("Converting missing data entries from ineligible to eligible for [ns=%s, coll=%s]", ns, coll)
				startKey, endKey := createRangeScanKeysForIneligibleMissingData(blkNum, ns, coll)
				collItr := s.db.GetIterator(startKey, endKey)
				collEntriesConverted := 0

				for collItr.Next() { 
					originalKey, originalVal := collItr.Key(), collItr.Value()
					modifiedKey := decodeMissingDataKey(originalKey)
					modifiedKey.isEligible = true
					batch.Delete(originalKey)
					copyVal := make([]byte, len(originalVal))
					copy(copyVal, originalVal)
					batch.Put(encodeMissingDataKey(modifiedKey), copyVal)
					collEntriesConverted++
					if batch.Len() > maxBatchSize {
						s.db.WriteBatch(batch, true)
						batch = leveldbhelper.NewUpdateBatch()
						sleepTime := time.Duration(batchesInterval)
						logger.Infof("Going to sleep for %d milliseconds between batches. Entries for [ns=%s, coll=%s] converted so far = %d",
							sleepTime, ns, coll, collEntriesConverted)
						s.purgerLock.Unlock()
						time.Sleep(sleepTime * time.Millisecond)
						s.purgerLock.Lock()
					}
				} 

				collItr.Release()
				logger.Infof("Converted all [%d] entries for [ns=%s, coll=%s]", collEntriesConverted, ns, coll)
				totalEntriesConverted += collEntriesConverted
			} 
		} 
		batch.Delete(collElgKey) 
	} 

	s.db.WriteBatch(batch, true)
	logger.Debugf("Converted [%d] inelligible mising data entries to elligible", totalEntriesConverted)
}


func (s *store) LastCommittedBlockHeight() (uint64, error) {
	if s.isEmpty {
		return 0, nil
	}
	return s.lastCommittedBlock + 1, nil
}


func (s *store) HasPendingBatch() (bool, error) {
	return s.batchPending, nil
}


func (s *store) IsEmpty() (bool, error) {
	return s.isEmpty, nil
}


func (s *store) Shutdown() {
	
}

func (s *store) nextBlockNum() uint64 {
	if s.isEmpty {
		return 0
	}
	return s.lastCommittedBlock + 1
}

func (s *store) hasPendingCommit() (bool, error) {
	var v []byte
	var err error
	if v, err = s.db.Get(pendingCommitKey); err != nil {
		return false, err
	}
	return v != nil, nil
}

func (s *store) getLastCommittedBlockNum() (bool, uint64, error) {
	var v []byte
	var err error
	if v, err = s.db.Get(lastCommittedBlkkey); v == nil || err != nil {
		return true, 0, err
	}
	return false, decodeLastCommittedBlockVal(v), nil
}

type collElgProcSync struct {
	notification, procComplete chan bool
}

func (sync *collElgProcSync) notify() {
	select {
	case sync.notification <- true:
		logger.Debugf("Signaled to collection eligibility processing routine")
	default: 
		logger.Debugf("Previous signal still pending. Skipping new signal")
	}
}

func (sync *collElgProcSync) waitForNotification() {
	<-sync.notification
}

func (sync *collElgProcSync) done() {
	select {
	case sync.procComplete <- true:
	default:
	}
}

func (sync *collElgProcSync) waitForDone() {
	<-sync.procComplete
}
