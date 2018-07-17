/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"fmt"

	"github.com/mcc-github/blockchain/protos/peer"
)




type ValidationParameterUpdatedErr struct {
	Key    string
	Height uint64
}

func (f ValidationParameterUpdatedErr) Error() string {
	return fmt.Sprintf("validation parameters for key %s have been changed in a transaction in block %d", f.Key, f.Height)
}









type KeyLevelValidationParameterManager interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	GetValidationParameterForKey(ch, cc, coll, key string, blockNum, txNum uint64) ([]byte, error)
	

	
	
	
	
	ExtractValidationParameterDependency(ch, cc string, blockNum, txNum uint64, rwset []byte)
	

	
	
	
	SetTxValidationCode(ch, cc string, blockNum, txNum uint64, vc peer.TxValidationCode)
	
}
