/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"fmt"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("vscc")




type ledgerKeyID struct {
	cc   string
	coll string
	key  string
}

func newLedgerKeyID(cc, coll, key string) *ledgerKeyID {
	return &ledgerKeyID{cc, coll, key}
}











type txDependency struct {
	mutex               sync.Mutex
	cond                *sync.Cond
	validationResultMap map[string]error
	depInserted         chan struct{}
}

func newTxDependency() *txDependency {
	txd := &txDependency{
		depInserted:         make(chan struct{}),
		validationResultMap: make(map[string]error),
	}
	txd.cond = sync.NewCond(&txd.mutex)
	return txd
}




func (d *txDependency) waitForDepInserted() {
	<-d.depInserted
}





func (d *txDependency) signalDepInserted() {
	close(d.depInserted)
}




func (d *txDependency) waitForAndRetrieveValidationResult(ns string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	err, ok := d.validationResultMap[ns]
	if ok {
		return err
	}

	for !ok {
		d.cond.Wait()
		err, ok = d.validationResultMap[ns]
	}

	return err
}





func (d *txDependency) signalValidationResult(ns string, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.validationResultMap[ns] = err
	d.cond.Broadcast()
}





type validationContext struct {
	
	
	
	mutex                sync.RWMutex
	blockHeight          uint64
	depsByTxnumMap       map[uint64]*txDependency
	depsByLedgerKeyIDMap map[ledgerKeyID]map[uint64]*txDependency
}

func (c *validationContext) forBlock(newHeight uint64) *validationContext {
	c.mutex.RLock()
	curHeight := c.blockHeight
	c.mutex.RUnlock()

	if curHeight > newHeight {
		logger.Panicf("programming error: block with number %d validated after block with number %d", newHeight, curHeight)
	}

	
	
	if curHeight < newHeight {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.blockHeight < newHeight {
			c.blockHeight = newHeight
			c.depsByLedgerKeyIDMap = map[ledgerKeyID]map[uint64]*txDependency{}
			c.depsByTxnumMap = map[uint64]*txDependency{}
		}
	}

	return c
}

func (c *validationContext) addDependency(kid *ledgerKeyID, txnum uint64, dep *txDependency) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	
	_, ok := c.depsByLedgerKeyIDMap[*kid]
	if !ok {
		c.depsByLedgerKeyIDMap[*kid] = map[uint64]*txDependency{}
	}

	c.depsByLedgerKeyIDMap[*kid][txnum] = dep
}

func (c *validationContext) dependenciesForTxnum(kid *ledgerKeyID, txnum uint64) []*txDependency {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var deps []*txDependency

	dl, in := c.depsByLedgerKeyIDMap[*kid]
	if in {
		deps = make([]*txDependency, 0, len(dl))
		for depTxnum, dep := range dl {
			if depTxnum < txnum {
				deps = append(deps, dep)
			}
		}
	}

	return deps
}

func (c *validationContext) getOrCreateDependencyByTxnum(txnum uint64) *txDependency {
	c.mutex.RLock()
	dep, ok := c.depsByTxnumMap[txnum]
	c.mutex.RUnlock()

	if !ok {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		dep, ok = c.depsByTxnumMap[txnum]
		if !ok {
			dep = newTxDependency()
			c.depsByTxnumMap[txnum] = dep
		}
	}

	return dep
}

func (c *validationContext) waitForValidationResults(kid *ledgerKeyID, blockNum uint64, txnum uint64) error {
	
	
	
	
	
	
	
	
	

	
	
	
	
	
	
	
	
	
	

	for _, dep := range c.dependenciesForTxnum(kid, txnum) {
		if valErr := dep.waitForAndRetrieveValidationResult(kid.cc); valErr == nil {
			return &ValidationParameterUpdatedError{
				CC:     kid.cc,
				Coll:   kid.coll,
				Key:    kid.key,
				Height: blockNum,
				Txnum:  txnum,
			}
		}
	}
	return nil
}




type KeyLevelValidationParameterManagerImpl struct {
	StateFetcher  validation.StateFetcher
	validationCtx validationContext
}






func (m *KeyLevelValidationParameterManagerImpl) ExtractValidationParameterDependency(blockNum, txNum uint64, rwsetBytes []byte) {
	vCtx := m.validationCtx.forBlock(blockNum)

	
	dep := vCtx.getOrCreateDependencyByTxnum(txNum)

	rwset := &rwsetutil.TxRwSet{}
	err := rwset.FromProtoBytes(rwsetBytes)
	
	
	if err == nil {
		
		
		
		
		for _, rws := range rwset.NsRwSets {
			for _, mw := range rws.KvRwSet.MetadataWrites {
				
				vCtx.addDependency(newLedgerKeyID(rws.NameSpace, "", mw.Key), txNum, dep)
			}

			for _, cw := range rws.CollHashedRwSets {
				for _, mw := range cw.HashedRwSet.MetadataWrites {
					
					vCtx.addDependency(newLedgerKeyID(rws.NameSpace, cw.CollectionName, string(mw.KeyHash)), txNum, dep)
				}
			}
		}
	} else {
		logger.Warningf("unmarshalling the read write set returned error '%s', skipping", err)
	}

	
	dep.signalDepInserted()
}



func (m *KeyLevelValidationParameterManagerImpl) GetValidationParameterForKey(cc, coll, key string, blockNum, txNum uint64) ([]byte, error) {
	vCtx := m.validationCtx.forBlock(blockNum)

	
	for i := int64(txNum) - 1; i >= 0; i-- {
		txdep := vCtx.getOrCreateDependencyByTxnum(uint64(i))
		txdep.waitForDepInserted()
	}

	
	
	err := vCtx.waitForValidationResults(newLedgerKeyID(cc, coll, key), blockNum, txNum)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	
	

	state, err := m.StateFetcher.FetchState()
	if err != nil {
		err = errors.WithMessage(err, "could not retrieve ledger")
		logger.Errorf(err.Error())
		return nil, err
	}
	defer state.Done()

	var mdMap map[string][]byte
	if coll == "" {
		mdMap, err = state.GetStateMetadata(cc, key)
		if err != nil {
			err = errors.WithMessage(err, fmt.Sprintf("could not retrieve metadata for %s:%s", cc, key))
			logger.Errorf(err.Error())
			return nil, err
		}
	} else {
		mdMap, err = state.GetPrivateDataMetadataByHash(cc, coll, []byte(key))
		if err != nil {
			err = errors.WithMessage(err, fmt.Sprintf("could not retrieve metadata for %s:%s:%x", cc, coll, []byte(key)))
			logger.Errorf(err.Error())
			return nil, err
		}
	}

	return mdMap[pb.MetaDataKeys_VALIDATION_PARAMETER.String()], nil
}





func (m *KeyLevelValidationParameterManagerImpl) SetTxValidationResult(ns string, blockNum, txNum uint64, err error) {
	vCtx := m.validationCtx.forBlock(blockNum)

	
	dep := vCtx.getOrCreateDependencyByTxnum(txNum)

	
	dep.signalValidationResult(ns, err)
}
