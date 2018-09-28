/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package internal

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/protos/peer"
)

var logger = flogging.MustGetLogger("valinternal")



type Validator interface {
	ValidateAndPrepareBatch(block *Block, doMVCCValidation bool) (*PubAndHashUpdates, error)
}



type Block struct {
	Num uint64
	Txs []*Transaction
}



type Transaction struct {
	IndexInBlock   int
	ID             string
	RWSet          *rwsetutil.TxRwSet
	ValidationCode peer.TxValidationCode
}



type PubAndHashUpdates struct {
	PubUpdates  *privacyenabledstate.PubUpdateBatch
	HashUpdates *privacyenabledstate.HashedUpdateBatch
}


func NewPubAndHashUpdates() *PubAndHashUpdates {
	return &PubAndHashUpdates{
		privacyenabledstate.NewPubUpdateBatch(),
		privacyenabledstate.NewHashedUpdateBatch(),
	}
}


func (t *Transaction) ContainsPvtWrites() bool {
	for _, ns := range t.RWSet.NsRwSets {
		for _, coll := range ns.CollHashedRwSets {
			if coll.PvtRwSetHash != nil {
				return true
			}
		}
	}
	return false
}



func (t *Transaction) RetrieveHash(ns string, coll string) []byte {
	if t.RWSet == nil {
		return nil
	}
	for _, nsData := range t.RWSet.NsRwSets {
		if nsData.NameSpace != ns {
			continue
		}

		for _, collData := range nsData.CollHashedRwSets {
			if collData.CollectionName == coll {
				return collData.PvtRwSetHash
			}
		}
	}
	return nil
}


func (u *PubAndHashUpdates) ApplyWriteSet(txRWSet *rwsetutil.TxRwSet, txHeight *version.Height, db privacyenabledstate.DB) error {
	txops, err := prepareTxOps(txRWSet, txHeight, u, db)
	logger.Debugf("txops=%#v", txops)
	if err != nil {
		return err
	}
	for compositeKey, keyops := range txops {
		if compositeKey.coll == "" {
			ns, key := compositeKey.ns, compositeKey.key
			if keyops.isDelete() {
				u.PubUpdates.Delete(ns, key, txHeight)
			} else {
				u.PubUpdates.PutValAndMetadata(ns, key, keyops.value, keyops.metadata, txHeight)
			}
		} else {
			ns, coll, keyHash := compositeKey.ns, compositeKey.coll, []byte(compositeKey.key)
			if keyops.isDelete() {
				u.HashUpdates.Delete(ns, coll, keyHash, txHeight)
			} else {
				u.HashUpdates.PutValHashAndMetadata(ns, coll, keyHash, keyops.value, keyops.metadata, txHeight)
			}
		}
	}
	return nil
}
