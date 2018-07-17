/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valimpl

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/statebasedval"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/valinternal"
)

var logger = flogging.MustGetLogger("valimpl")





type DefaultImpl struct {
	txmgr txmgr.TxMgr
	db    privacyenabledstate.DB
	valinternal.InternalValidator
}



func NewStatebasedValidator(txmgr txmgr.TxMgr, db privacyenabledstate.DB) validator.Validator {
	return &DefaultImpl{txmgr, db, statebasedval.NewValidator(db)}
}


func (impl *DefaultImpl) ValidateAndPrepareBatch(blockAndPvtdata *ledger.BlockAndPvtData,
	doMVCCValidation bool) (*privacyenabledstate.UpdateBatch, error) {
	block := blockAndPvtdata.Block
	logger.Debugf("ValidateAndPrepareBatch() for block number = [%d]", block.Header.Number)
	var internalBlock *valinternal.Block
	var pubAndHashUpdates *valinternal.PubAndHashUpdates
	var pvtUpdates *privacyenabledstate.PvtUpdateBatch
	var err error

	logger.Debug("preprocessing ProtoBlock...")
	if internalBlock, err = preprocessProtoBlock(impl.txmgr, impl.db.ValidateKeyValue, block, doMVCCValidation); err != nil {
		return nil, err
	}

	if pubAndHashUpdates, err = impl.InternalValidator.ValidateAndPrepareBatch(internalBlock, doMVCCValidation); err != nil {
		return nil, err
	}
	logger.Debug("validating rwset...")
	if pvtUpdates, err = validateAndPreparePvtBatch(internalBlock, blockAndPvtdata.BlockPvtData); err != nil {
		return nil, err
	}
	logger.Debug("postprocessing ProtoBlock...")
	postprocessProtoBlock(block, internalBlock)
	logger.Debug("ValidateAndPrepareBatch() complete")
	return &privacyenabledstate.UpdateBatch{
		PubUpdates:  pubAndHashUpdates.PubUpdates,
		HashUpdates: pubAndHashUpdates.HashUpdates,
		PvtUpdates:  pvtUpdates,
	}, nil
}
