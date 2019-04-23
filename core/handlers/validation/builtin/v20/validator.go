/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package v20

import (
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)











type StateBasedValidator interface {
	
	
	PreValidate(txNum uint64, block *common.Block)

	
	
	
	Validate(cc string, blockNum, txNum uint64, rwset, prp, ep []byte, endorsements []*peer.Endorsement) commonerrors.TxValidationError

	
	
	PostValidate(cc string, blockNum, txNum uint64, err error)
}
