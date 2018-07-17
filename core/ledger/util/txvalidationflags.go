/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"github.com/mcc-github/blockchain/protos/peer"
)


type TxValidationFlags []uint8



func NewTxValidationFlags(size int) TxValidationFlags {
	return newTxValidationFlagsSetValue(size, peer.TxValidationCode_NOT_VALIDATED)
}



func NewTxValidationFlagsSetValue(size int, value peer.TxValidationCode) TxValidationFlags {
	return newTxValidationFlagsSetValue(size, value)
}

func newTxValidationFlagsSetValue(size int, value peer.TxValidationCode) TxValidationFlags {
	inst := make(TxValidationFlags, size)
	for i := range inst {
		inst[i] = uint8(value)
	}

	return inst
}


func (obj TxValidationFlags) SetFlag(txIndex int, flag peer.TxValidationCode) {
	obj[txIndex] = uint8(flag)
}


func (obj TxValidationFlags) Flag(txIndex int) peer.TxValidationCode {
	return peer.TxValidationCode(obj[txIndex])
}


func (obj TxValidationFlags) IsValid(txIndex int) bool {
	return obj.IsSetTo(txIndex, peer.TxValidationCode_VALID)
}


func (obj TxValidationFlags) IsInvalid(txIndex int) bool {
	return !obj.IsValid(txIndex)
}


func (obj TxValidationFlags) IsSetTo(txIndex int, flag peer.TxValidationCode) bool {
	return obj.Flag(txIndex) == flag
}
