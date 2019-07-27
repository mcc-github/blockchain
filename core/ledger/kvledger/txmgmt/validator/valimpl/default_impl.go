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
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/internal"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/validator/statebasedval"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
)

var logger = flogging.MustGetLogger("valimpl")





type DefaultImpl struct {
	txmgr              txmgr.TxMgr
	db                 privacyenabledstate.DB
	internalValidator  internal.Validator
	customTxProcessors map[common.HeaderType]ledger.CustomTxProcessor
}



func NewStatebasedValidator(
	txmgr txmgr.TxMgr,
	db privacyenabledstate.DB,
	customTxProcessors map[common.HeaderType]ledger.CustomTxProcessor,
) validator.Validator {
	return &DefaultImpl{txmgr, db, statebasedval.NewValidator(db), customTxProcessors}
}


func (impl *DefaultImpl) ValidateAndPrepareBatch(blockAndPvtdata *ledger.BlockAndPvtData,
	doMVCCValidation bool) (*privacyenabledstate.UpdateBatch, []*txmgr.TxStatInfo, error) {
	block := blockAndPvtdata.Block
	logger.Debugf("ValidateAndPrepareBatch() for block number = [%d]", block.Header.Number)
	var internalBlock *internal.Block
	var txsStatInfo []*txmgr.TxStatInfo
	var pubAndHashUpdates *internal.PubAndHashUpdates
	var pvtUpdates *privacyenabledstate.PvtUpdateBatch
	var err error

	logger.Debug("preprocessing ProtoBlock...")
	if internalBlock, txsStatInfo, err = preprocessProtoBlock(
		impl.txmgr,
		impl.db.ValidateKeyValue,
		block,
		doMVCCValidation,
		impl.customTxProcessors,
	); err != nil {
		return nil, nil, err
	}

	if pubAndHashUpdates, err = impl.internalValidator.ValidateAndPrepareBatch(internalBlock, doMVCCValidation); err != nil {
		return nil, nil, err
	}
	logger.Debug("validating rwset...")
	if pvtUpdates, err = validateAndPreparePvtBatch(
		internalBlock,
		impl.db,
		pubAndHashUpdates,
		blockAndPvtdata.PvtData,
		impl.customTxProcessors,
	); err != nil {
		return nil, nil, err
	}
	logger.Debug("postprocessing ProtoBlock...")
	postprocessProtoBlock(block, internalBlock)
	logger.Debug("ValidateAndPrepareBatch() complete")

	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	for i := range txsFilter {
		txsStatInfo[i].ValidationCode = txsFilter.Flag(i)
	}
	return &privacyenabledstate.UpdateBatch{
		PubUpdates:  pubAndHashUpdates.PubUpdates,
		HashUpdates: pubAndHashUpdates.HashUpdates,
		PvtUpdates:  pvtUpdates,
	}, txsStatInfo, nil
}
