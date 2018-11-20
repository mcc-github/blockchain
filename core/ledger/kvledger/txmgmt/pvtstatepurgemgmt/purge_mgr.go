/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtstatepurgemgmt

import (
	"math"
	"sync"

	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/mcc-github/blockchain/core/ledger/util"
)


type PurgeMgr interface {
	
	PrepareForExpiringKeys(expiringAtBlk uint64)
	
	WaitForPrepareToFinish()
	
	DeleteExpiredAndUpdateBookkeeping(
		pvtUpdates *privacyenabledstate.PvtUpdateBatch,
		hashedUpdates *privacyenabledstate.HashedUpdateBatch) error
	
	UpdateBookkeepingForPvtDataOfOldBlocks(pvtUpdates *privacyenabledstate.PvtUpdateBatch) error
	
	BlockCommitDone() error
}

type keyAndVersion struct {
	key             string
	committingBlock uint64
	purgeKeyOnly    bool
}

type expiryInfoMap map[privacyenabledstate.HashedCompositeKey]*keyAndVersion

type workingset struct {
	toPurge             expiryInfoMap
	toClearFromSchedule []*expiryInfoKey
	expiringBlk         uint64
	err                 error
}

type purgeMgr struct {
	btlPolicy pvtdatapolicy.BTLPolicy
	db        privacyenabledstate.DB
	expKeeper expiryKeeper

	lock    *sync.Mutex
	waitGrp *sync.WaitGroup

	workingset *workingset
}


func InstantiatePurgeMgr(ledgerid string, db privacyenabledstate.DB, btlPolicy pvtdatapolicy.BTLPolicy, bookkeepingProvider bookkeeping.Provider) (PurgeMgr, error) {
	return &purgeMgr{
		btlPolicy: btlPolicy,
		db:        db,
		expKeeper: newExpiryKeeper(ledgerid, bookkeepingProvider),
		lock:      &sync.Mutex{},
		waitGrp:   &sync.WaitGroup{},
	}, nil
}


func (p *purgeMgr) PrepareForExpiringKeys(expiringAtBlk uint64) {
	p.waitGrp.Add(1)
	go func() {
		p.lock.Lock()
		p.waitGrp.Done()
		defer p.lock.Unlock()
		p.workingset = p.prepareWorkingsetFor(expiringAtBlk)
	}()
	p.waitGrp.Wait()
}


func (p *purgeMgr) WaitForPrepareToFinish() {
	p.lock.Lock()
	p.lock.Unlock()
}

func (p *purgeMgr) UpdateBookkeepingForPvtDataOfOldBlocks(pvtUpdates *privacyenabledstate.PvtUpdateBatch) error {
	builder := newExpiryScheduleBuilder(p.btlPolicy)
	pvtUpdateCompositeKeyMap := pvtUpdates.ToCompositeKeyMap()
	for k, vv := range pvtUpdateCompositeKeyMap {
		builder.add(k.Namespace, k.CollectionName, k.Key, util.ComputeStringHash(k.Key), vv)
	}

	var updatedList []*expiryInfo
	for _, toAdd := range builder.getExpiryInfo() {
		toUpdate, err := p.expKeeper.retrieveByExpiryKey(toAdd.expiryInfoKey)
		if err != nil {
			return err
		}
		
		
		
		toUpdate.pvtdataKeys.addAll(toAdd.pvtdataKeys)
		updatedList = append(updatedList, toUpdate)
	}
	return p.expKeeper.updateBookkeeping(updatedList, nil)
}


func (p *purgeMgr) DeleteExpiredAndUpdateBookkeeping(
	pvtUpdates *privacyenabledstate.PvtUpdateBatch,
	hashedUpdates *privacyenabledstate.HashedUpdateBatch) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.workingset.err != nil {
		return p.workingset.err
	}

	listExpiryInfo, err := buildExpirySchedule(p.btlPolicy, pvtUpdates, hashedUpdates)
	if err != nil {
		return err
	}

	
	
	for compositeHashedKey, keyAndVersion := range p.workingset.toPurge {
		ns := compositeHashedKey.Namespace
		coll := compositeHashedKey.CollectionName
		keyHash := []byte(compositeHashedKey.KeyHash)
		key := keyAndVersion.key
		purgeKeyOnly := keyAndVersion.purgeKeyOnly
		hashUpdated := hashedUpdates.Contains(ns, coll, keyHash)
		pvtKeyUpdated := pvtUpdates.Contains(ns, coll, key)

		logger.Debugf("Checking whether the key [ns=%s, coll=%s, keyHash=%x, purgeKeyOnly=%t] "+
			"is updated in the update batch for the committing block - hashUpdated=%t, and pvtKeyUpdated=%t",
			ns, coll, keyHash, purgeKeyOnly, hashUpdated, pvtKeyUpdated)

		expiringTxVersion := version.NewHeight(p.workingset.expiringBlk, math.MaxUint64)
		if !hashUpdated && !purgeKeyOnly {
			logger.Debugf("Adding the hashed key to be purged to the delete list in the update batch")
			hashedUpdates.Delete(ns, coll, keyHash, expiringTxVersion)
		}
		if key != "" && !pvtKeyUpdated {
			logger.Debugf("Adding the pvt key to be purged to the delete list in the update batch")
			pvtUpdates.Delete(ns, coll, key, expiringTxVersion)
		}
	}
	return p.expKeeper.updateBookkeeping(listExpiryInfo, nil)
}









func (p *purgeMgr) BlockCommitDone() error {
	defer func() { p.workingset = nil }()
	return p.expKeeper.updateBookkeeping(nil, p.workingset.toClearFromSchedule)
}



func (p *purgeMgr) prepareWorkingsetFor(expiringAtBlk uint64) *workingset {
	logger.Debugf("Preparing potential purge list working-set for expiringAtBlk [%d]", expiringAtBlk)
	workingset := &workingset{expiringBlk: expiringAtBlk}
	
	expiryInfo, err := p.expKeeper.retrieve(expiringAtBlk)
	if err != nil {
		workingset.err = err
		return workingset
	}
	
	toPurge := transformToExpiryInfoMap(expiryInfo)
	
	p.preloadCommittedVersionsInCache(toPurge)
	var expiryInfoKeysToClear []*expiryInfoKey

	if len(toPurge) == 0 {
		logger.Debugf("No expiry entry found for expiringAtBlk [%d]", expiringAtBlk)
		return workingset
	}
	logger.Debugf("Total [%d] expiring entries found. Evaluaitng whether some of these keys have been overwritten in later blocks...", len(toPurge))

	for purgeEntryK, purgeEntryV := range toPurge {
		logger.Debugf("Evaluating for hashedKey [%s]", purgeEntryK)
		expiryInfoKeysToClear = append(expiryInfoKeysToClear, &expiryInfoKey{committingBlk: purgeEntryV.committingBlock, expiryBlk: expiringAtBlk})
		currentVersion, err := p.db.GetKeyHashVersion(purgeEntryK.Namespace, purgeEntryK.CollectionName, []byte(purgeEntryK.KeyHash))
		if err != nil {
			workingset.err = err
			return workingset
		}

		if sameVersion(currentVersion, purgeEntryV.committingBlock) {
			logger.Debugf(
				"The version of the hashed key in the committed state and in the expiry entry is same " +
					"hence, keeping the entry in the purge list")
			continue
		}

		logger.Debugf("The version of the hashed key in the committed state and in the expiry entry is different")
		if purgeEntryV.key != "" {
			logger.Debugf("The expiry entry also contains the raw key along with the key hash")
			committedPvtVerVal, err := p.db.GetPrivateData(purgeEntryK.Namespace, purgeEntryK.CollectionName, purgeEntryV.key)
			if err != nil {
				workingset.err = err
				return workingset
			}

			if sameVersionFromVal(committedPvtVerVal, purgeEntryV.committingBlock) {
				logger.Debugf(
					"The version of the pvt key in the committed state and in the expiry entry is same" +
						"Including only key in the purge list and not the hashed key")
				purgeEntryV.purgeKeyOnly = true
				continue
			}
		}

		
		logger.Debugf("Removing from purge list - the key hash and key (if present, in the expiry entry)")
		delete(toPurge, purgeEntryK)
	}
	
	workingset.toPurge = toPurge
	
	workingset.toClearFromSchedule = expiryInfoKeysToClear
	return workingset
}

func (p *purgeMgr) preloadCommittedVersionsInCache(expInfoMap expiryInfoMap) {
	if !p.db.IsBulkOptimizable() {
		return
	}
	var hashedKeys []*privacyenabledstate.HashedCompositeKey
	for k := range expInfoMap {
		hashedKeys = append(hashedKeys, &k)
	}
	p.db.LoadCommittedVersionsOfPubAndHashedKeys(nil, hashedKeys)
}

func transformToExpiryInfoMap(expiryInfo []*expiryInfo) expiryInfoMap {
	expinfoMap := make(expiryInfoMap)
	for _, expinfo := range expiryInfo {
		for ns, colls := range expinfo.pvtdataKeys.Map {
			for coll, keysAndHashes := range colls.Map {
				for _, keyAndHash := range keysAndHashes.List {
					compositeKey := privacyenabledstate.HashedCompositeKey{Namespace: ns, CollectionName: coll, KeyHash: string(keyAndHash.Hash)}
					expinfoMap[compositeKey] = &keyAndVersion{key: keyAndHash.Key, committingBlock: expinfo.expiryInfoKey.committingBlk}
				}
			}
		}
	}
	return expinfoMap
}

func sameVersion(version *version.Height, blockNum uint64) bool {
	return version != nil && version.BlockNum == blockNum
}

func sameVersionFromVal(vv *statedb.VersionedValue, blockNum uint64) bool {
	return vv != nil && sameVersion(vv.Version, blockNum)
}
