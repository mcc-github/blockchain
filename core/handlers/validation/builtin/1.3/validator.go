/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package builtin1_3

import (
	"github.com/mcc-github/blockchain/protos/peer"
)











type StateBasedValidator interface {
	
	
	PreValidate(ch, cc string, blockNum, txNum uint64, rwset []byte) error
	

	
	
	
	Validate(ch, cc string, blockNum, txNum uint64, rwset, prp, ep []byte, endorsements []*peer.Endorsement) error
	

	
	
	PostValidate(ch, cc string, blockNum, txNum uint64, vc peer.TxValidationCode) error
	
}
