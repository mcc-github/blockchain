/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"fmt"
)




type ValidationParameterUpdatedError struct {
	Key    string
	Height uint64
}

func (f *ValidationParameterUpdatedError) Error() string {
	return fmt.Sprintf("validation parameters for key %s have been changed in a transaction in block %d", f.Key, f.Height)
}









type KeyLevelValidationParameterManager interface {
	
	
	
	
	
	
	
	
	
	
	
	
	GetValidationParameterForKey(cc, coll, key string, blockNum, txNum uint64) ([]byte, error)

	
	
	
	
	ExtractValidationParameterDependency(blockNum, txNum uint64, rwset []byte)

	
	
	
	
	SetTxValidationResult(cc string, blockNum, txNum uint64, err error)
}
