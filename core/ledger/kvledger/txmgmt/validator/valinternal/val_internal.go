/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valinternal

import (
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/protos/peer"
)



type InternalValidator interface {
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


func (u *PubAndHashUpdates) ApplyWriteSet(txRWSet *rwsetutil.TxRwSet, txHeight *version.Height) {
	for _, nsRWSet := range txRWSet.NsRwSets {
		ns := nsRWSet.NameSpace
		for _, kvWrite := range nsRWSet.KvRwSet.Writes {
			if kvWrite.IsDelete {
				u.PubUpdates.Delete(ns, kvWrite.Key, txHeight)
			} else {
				u.PubUpdates.Put(ns, kvWrite.Key, kvWrite.Value, txHeight)
			}
		}

		for _, collHashRWset := range nsRWSet.CollHashedRwSets {
			coll := collHashRWset.CollectionName
			for _, hashedWrite := range collHashRWset.HashedRwSet.HashedWrites {
				if hashedWrite.IsDelete {
					u.HashUpdates.Delete(ns, coll, hashedWrite.KeyHash, txHeight)
				} else {
					u.HashUpdates.Put(ns, coll, hashedWrite.KeyHash, hashedWrite.ValueHash, txHeight)
				}
			}
		}
	}
}
