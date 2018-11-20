/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtstatepurgemgmt

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/bookkeeping"
)

var logger = flogging.MustGetLogger("pvtstatepurgemgmt")

const (
	expiryPrefix = '1'
)


type expiryInfoKey struct {
	committingBlk uint64
	expiryBlk     uint64
}





type expiryInfo struct {
	expiryInfoKey *expiryInfoKey
	pvtdataKeys   *PvtdataKeys
}


type expiryKeeper interface {
	
	updateBookkeeping(toTrack []*expiryInfo, toClear []*expiryInfoKey) error
	
	retrieve(expiringAtBlkNum uint64) ([]*expiryInfo, error)
	
	retrieveByExpiryKey(expiryKey *expiryInfoKey) (*expiryInfo, error)
}

func newExpiryKeeper(ledgerid string, provider bookkeeping.Provider) expiryKeeper {
	return &expKeeper{provider.GetDBHandle(ledgerid, bookkeeping.PvtdataExpiry)}
}

type expKeeper struct {
	db *leveldbhelper.DBHandle
}













func (ek *expKeeper) updateBookkeeping(toTrack []*expiryInfo, toClear []*expiryInfoKey) error {
	updateBatch := leveldbhelper.NewUpdateBatch()
	for _, expinfo := range toTrack {
		k, v, err := encodeKV(expinfo)
		if err != nil {
			return err
		}
		updateBatch.Put(k, v)
	}
	for _, expinfokey := range toClear {
		updateBatch.Delete(encodeExpiryInfoKey(expinfokey))
	}
	return ek.db.WriteBatch(updateBatch, true)
}

func (ek *expKeeper) retrieve(expiringAtBlkNum uint64) ([]*expiryInfo, error) {
	startKey := encodeExpiryInfoKey(&expiryInfoKey{expiryBlk: expiringAtBlkNum, committingBlk: 0})
	endKey := encodeExpiryInfoKey(&expiryInfoKey{expiryBlk: expiringAtBlkNum + 1, committingBlk: 0})
	itr := ek.db.GetIterator(startKey, endKey)
	defer itr.Release()

	var listExpinfo []*expiryInfo
	for itr.Next() {
		expinfo, err := decodeExpiryInfo(itr.Key(), itr.Value())
		if err != nil {
			return nil, err
		}
		listExpinfo = append(listExpinfo, expinfo)
	}
	return listExpinfo, nil
}

func (ek *expKeeper) retrieveByExpiryKey(expiryKey *expiryInfoKey) (*expiryInfo, error) {
	key := encodeExpiryInfoKey(expiryKey)
	value, err := ek.db.Get(key)
	if err != nil {
		return nil, err
	}
	return decodeExpiryInfo(key, value)
}

func encodeKV(expinfo *expiryInfo) (key []byte, value []byte, err error) {
	key = encodeExpiryInfoKey(expinfo.expiryInfoKey)
	value, err = encodeExpiryInfoValue(expinfo.pvtdataKeys)
	return
}

func encodeExpiryInfoKey(expinfoKey *expiryInfoKey) []byte {
	key := append([]byte{expiryPrefix}, util.EncodeOrderPreservingVarUint64(expinfoKey.expiryBlk)...)
	return append(key, util.EncodeOrderPreservingVarUint64(expinfoKey.committingBlk)...)
}

func encodeExpiryInfoValue(pvtdataKeys *PvtdataKeys) ([]byte, error) {
	return proto.Marshal(pvtdataKeys)
}

func decodeExpiryInfo(key []byte, value []byte) (*expiryInfo, error) {
	expiryBlk, n := util.DecodeOrderPreservingVarUint64(key[1:])
	committingBlk, _ := util.DecodeOrderPreservingVarUint64(key[n+1:])
	pvtdataKeys := &PvtdataKeys{}
	if err := proto.Unmarshal(value, pvtdataKeys); err != nil {
		return nil, err
	}
	return &expiryInfo{
			expiryInfoKey: &expiryInfoKey{committingBlk: committingBlk, expiryBlk: expiryBlk},
			pvtdataKeys:   pvtdataKeys},
		nil
}
