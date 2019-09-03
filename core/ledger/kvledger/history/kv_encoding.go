/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package history

import (
	"bytes"

	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/pkg/errors"
)

type dataKey []byte
type rangeScan struct {
	startKey, endKey []byte
}

var (
	compositeKeySep = []byte{0x00} 
	dataKeyPrefix   = []byte{'d'}  
	savePointKey    = []byte{'s'}  
	emptyValue      = []byte{}     
)





func constructDataKey(ns string, key string, blocknum uint64, trannum uint64) dataKey {
	k := append([]byte(ns), compositeKeySep...)
	k = append(k, util.EncodeOrderPreservingVarUint64(uint64(len(key)))...)
	k = append(k, []byte(key)...)
	k = append(k, compositeKeySep...)
	k = append(k, util.EncodeOrderPreservingVarUint64(blocknum)...)
	k = append(k, util.EncodeOrderPreservingVarUint64(trannum)...)
	return dataKey(k)
}





func constructRangeScan(ns string, key string) *rangeScan {
	k := append([]byte(ns), compositeKeySep...)
	k = append(k, util.EncodeOrderPreservingVarUint64(uint64(len(key)))...)
	k = append(k, []byte(key)...)
	k = append(k, compositeKeySep...)

	return &rangeScan{
		startKey: k,
		endKey:   append(k, 0xff),
	}
}

func (r *rangeScan) decodeBlockNumTranNum(dataKey dataKey) (uint64, uint64, error) {
	blockNumTranNumBytes := bytes.TrimPrefix(dataKey, r.startKey)
	blockNum, blockBytesConsumed, err := util.DecodeOrderPreservingVarUint64(blockNumTranNumBytes)
	if err != nil {
		return 0, 0, err
	}

	tranNum, tranBytesConsumed, err := util.DecodeOrderPreservingVarUint64(blockNumTranNumBytes[blockBytesConsumed:])
	if err != nil {
		return 0, 0, err
	}

	
	if blockBytesConsumed+tranBytesConsumed != len(blockNumTranNumBytes) {
		return 0, 0, errors.Errorf("number of decoded bytes (%d) is not equal to the length of blockNumTranNumBytes (%d)",
			blockBytesConsumed+tranBytesConsumed, len(blockNumTranNumBytes))
	}
	return blockNum, tranNum, nil
}
