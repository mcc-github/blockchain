/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transientstore

import (
	"bytes"
	"errors"

	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
)

var (
	prwsetPrefix             = []byte("P")[0] 
	purgeIndexByHeightPrefix = []byte("H")[0] 
	purgeIndexByTxidPrefix   = []byte("T")[0] 
	compositeKeySep          = byte(0x00)
)



func createCompositeKeyForPvtRWSet(txid string, uuid string, blockHeight uint64) []byte {
	var compositeKey []byte
	compositeKey = append(compositeKey, prwsetPrefix)
	compositeKey = append(compositeKey, compositeKeySep)
	compositeKey = append(compositeKey, createCompositeKeyWithoutPrefixForTxid(txid, uuid, blockHeight)...)

	return compositeKey
}




func createCompositeKeyForPurgeIndexByTxid(txid string, uuid string, blockHeight uint64) []byte {
	var compositeKey []byte
	compositeKey = append(compositeKey, purgeIndexByTxidPrefix)
	compositeKey = append(compositeKey, compositeKeySep)
	compositeKey = append(compositeKey, createCompositeKeyWithoutPrefixForTxid(txid, uuid, blockHeight)...)

	return compositeKey
}


func createCompositeKeyWithoutPrefixForTxid(txid string, uuid string, blockHeight uint64) []byte {
	var compositeKey []byte
	compositeKey = append(compositeKey, []byte(txid)...)
	compositeKey = append(compositeKey, compositeKeySep)
	compositeKey = append(compositeKey, []byte(uuid)...)
	compositeKey = append(compositeKey, compositeKeySep)
	compositeKey = append(compositeKey, util.EncodeOrderPreservingVarUint64(blockHeight)...)

	return compositeKey
}




func createCompositeKeyForPurgeIndexByHeight(blockHeight uint64, txid string, uuid string) []byte {
	var compositeKey []byte
	compositeKey = append(compositeKey, purgeIndexByHeightPrefix)
	compositeKey = append(compositeKey, compositeKeySep)
	compositeKey = append(compositeKey, util.EncodeOrderPreservingVarUint64(blockHeight)...)
	compositeKey = append(compositeKey, compositeKeySep)
	compositeKey = append(compositeKey, []byte(txid)...)
	compositeKey = append(compositeKey, compositeKeySep)
	compositeKey = append(compositeKey, []byte(uuid)...)

	return compositeKey
}



func splitCompositeKeyOfPvtRWSet(compositeKey []byte) (uuid string, blockHeight uint64) {
	return splitCompositeKeyWithoutPrefixForTxid(compositeKey[2:])
}



func splitCompositeKeyOfPurgeIndexByTxid(compositeKey []byte) (uuid string, blockHeight uint64) {
	return splitCompositeKeyWithoutPrefixForTxid(compositeKey[2:])
}



func splitCompositeKeyOfPurgeIndexByHeight(compositeKey []byte) (txid string, uuid string, blockHeight uint64) {
	var n int
	blockHeight, n, _ = util.DecodeOrderPreservingVarUint64(compositeKey[2:])
	splits := bytes.Split(compositeKey[n+3:], []byte{compositeKeySep})
	txid = string(splits[0])
	uuid = string(splits[1])
	return
}



func splitCompositeKeyWithoutPrefixForTxid(compositeKey []byte) (uuid string, blockHeight uint64) {
	
	firstSepIndex := bytes.IndexByte(compositeKey, compositeKeySep)
	secondSepIndex := firstSepIndex + bytes.IndexByte(compositeKey[firstSepIndex+1:], compositeKeySep) + 1
	uuid = string(compositeKey[firstSepIndex+1 : secondSepIndex])
	blockHeight, _, _ = util.DecodeOrderPreservingVarUint64(compositeKey[secondSepIndex+1:])
	return
}


func createTxidRangeStartKey(txid string) []byte {
	var startKey []byte
	startKey = append(startKey, prwsetPrefix)
	startKey = append(startKey, compositeKeySep)
	startKey = append(startKey, []byte(txid)...)
	startKey = append(startKey, compositeKeySep)
	return startKey
}


func createTxidRangeEndKey(txid string) []byte {
	var endKey []byte
	endKey = append(endKey, prwsetPrefix)
	endKey = append(endKey, compositeKeySep)
	endKey = append(endKey, []byte(txid)...)
	
	
	endKey = append(endKey, byte(0xff))
	return endKey
}



func createPurgeIndexByHeightRangeStartKey(blockHeight uint64) []byte {
	var startKey []byte
	startKey = append(startKey, purgeIndexByHeightPrefix)
	startKey = append(startKey, compositeKeySep)
	startKey = append(startKey, util.EncodeOrderPreservingVarUint64(blockHeight)...)
	startKey = append(startKey, compositeKeySep)
	return startKey
}



func createPurgeIndexByHeightRangeEndKey(blockHeight uint64) []byte {
	var endKey []byte
	endKey = append(endKey, purgeIndexByHeightPrefix)
	endKey = append(endKey, compositeKeySep)
	endKey = append(endKey, util.EncodeOrderPreservingVarUint64(blockHeight)...)
	endKey = append(endKey, byte(0xff))
	return endKey
}



func createPurgeIndexByTxidRangeStartKey(txid string) []byte {
	var startKey []byte
	startKey = append(startKey, purgeIndexByTxidPrefix)
	startKey = append(startKey, compositeKeySep)
	startKey = append(startKey, []byte(txid)...)
	startKey = append(startKey, compositeKeySep)
	return startKey
}



func createPurgeIndexByTxidRangeEndKey(txid string) []byte {
	var endKey []byte
	endKey = append(endKey, purgeIndexByTxidPrefix)
	endKey = append(endKey, compositeKeySep)
	endKey = append(endKey, []byte(txid)...)
	
	
	endKey = append(endKey, byte(0xff))
	return endKey
}



func trimPvtWSet(pvtWSet *rwset.TxPvtReadWriteSet, filter ledger.PvtNsCollFilter) *rwset.TxPvtReadWriteSet {
	if filter == nil {
		return pvtWSet
	}

	var filteredNsRwSet []*rwset.NsPvtReadWriteSet
	for _, ns := range pvtWSet.NsPvtRwset {
		var filteredCollRwSet []*rwset.CollectionPvtReadWriteSet
		for _, coll := range ns.CollectionPvtRwset {
			if filter.Has(ns.Namespace, coll.CollectionName) {
				filteredCollRwSet = append(filteredCollRwSet, coll)
			}
		}
		if filteredCollRwSet != nil {
			filteredNsRwSet = append(filteredNsRwSet,
				&rwset.NsPvtReadWriteSet{
					Namespace:          ns.Namespace,
					CollectionPvtRwset: filteredCollRwSet,
				},
			)
		}
	}
	var filteredTxPvtRwSet *rwset.TxPvtReadWriteSet
	if filteredNsRwSet != nil {
		filteredTxPvtRwSet = &rwset.TxPvtReadWriteSet{
			DataModel:  pvtWSet.GetDataModel(),
			NsPvtRwset: filteredNsRwSet,
		}
	}
	return filteredTxPvtRwSet
}

func trimPvtCollectionConfigs(configs map[string]*common.CollectionConfigPackage,
	filter ledger.PvtNsCollFilter) (map[string]*common.CollectionConfigPackage, error) {
	if filter == nil {
		return configs, nil
	}
	result := make(map[string]*common.CollectionConfigPackage)

	for ns, pkg := range configs {
		result[ns] = &common.CollectionConfigPackage{}
		for _, colConf := range pkg.GetConfig() {
			switch cconf := colConf.Payload.(type) {
			case *common.CollectionConfig_StaticCollectionConfig:
				if filter.Has(ns, cconf.StaticCollectionConfig.Name) {
					result[ns].Config = append(result[ns].Config, colConf)
				}
			default:
				return nil, errors.New("unexpected collection type")
			}
		}
	}
	return result, nil
}
