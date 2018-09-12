/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"fmt"
)




type ValidationParameterUpdatedError struct {
	CC     string
	Coll   string
	Key    string
	Height uint64
	Txnum  uint64
}

func (f *ValidationParameterUpdatedError) Error() string {
	return fmt.Sprintf("validation parameters for key [%s] in namespace [%s:%s] have been changed in transaction %d of block %d", f.Key, f.CC, f.Coll, f.Txnum, f.Height)
}









type KeyLevelValidationParameterManager interface {
	
	
	
	
	
	
	
	
	
	
	
	
	GetValidationParameterForKey(cc, coll, key string, blockNum, txNum uint64) ([]byte, error)

	
	
	
	
	ExtractValidationParameterDependency(blockNum, txNum uint64, rwset []byte)

	
	
	
	
	SetTxValidationResult(cc string, blockNum, txNum uint64, err error)
}
