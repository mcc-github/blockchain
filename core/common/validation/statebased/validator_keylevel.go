/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"fmt"
	"sync"

	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/core/handlers/validation/api/policies"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)




type policyChecker struct {
	someEPChecked bool
	ccEPChecked   bool
	vpmgr         KeyLevelValidationParameterManager
	policySupport validation.PolicyEvaluator
	ccEP          []byte
	signatureSet  []*common.SignedData
}

func (p *policyChecker) checkCCEPIfCondition(cc string, blockNum, txNum uint64, condition bool) commonerrors.TxValidationError {
	if condition {
		return nil
	}

	
	err := p.policySupport.Evaluate(p.ccEP, p.signatureSet)
	if err != nil {
		return policyErr(errors.Wrapf(err, "validation of endorsement policy for chaincode %s in tx %d:%d failed", cc, blockNum, txNum))
	}

	p.ccEPChecked = true
	p.someEPChecked = true
	return nil
}

func (p *policyChecker) checkCCEPIfNotChecked(cc string, blockNum, txNum uint64) commonerrors.TxValidationError {
	return p.checkCCEPIfCondition(cc, blockNum, txNum, p.ccEPChecked)
}

func (p *policyChecker) checkCCEPIfNoEPChecked(cc string, blockNum, txNum uint64) commonerrors.TxValidationError {
	return p.checkCCEPIfCondition(cc, blockNum, txNum, p.someEPChecked)
}

func (p *policyChecker) checkSBAndCCEP(cc, coll, key string, blockNum, txNum uint64) commonerrors.TxValidationError {
	
	vp, err := p.vpmgr.GetValidationParameterForKey(cc, coll, key, blockNum, txNum)
	if err != nil {
		switch err := err.(type) {
		case *ValidationParameterUpdatedError:
			return policyErr(err)
		default:
			return &commonerrors.VSCCExecutionFailureError{
				Err: err,
			}
		}
	}

	
	if len(vp) == 0 {
		return p.checkCCEPIfNotChecked(cc, blockNum, txNum)
	}

	
	err = p.policySupport.Evaluate(vp, p.signatureSet)
	if err != nil {
		return policyErr(errors.Wrapf(err, "validation of key %s (coll'%s':ns'%s') in tx %d:%d failed", key, coll, cc, blockNum, txNum))
	}

	p.someEPChecked = true

	return nil
}




type blockDependency struct {
	mutex     sync.Mutex
	blockNum  uint64
	txDepOnce []sync.Once
}


type KeyLevelValidator struct {
	vpmgr         KeyLevelValidationParameterManager
	policySupport validation.PolicyEvaluator
	blockDep      blockDependency
}

func NewKeyLevelValidator(policySupport validation.PolicyEvaluator, vpmgr KeyLevelValidationParameterManager) *KeyLevelValidator {
	return &KeyLevelValidator{
		vpmgr:         vpmgr,
		policySupport: policySupport,
		blockDep:      blockDependency{},
	}
}

func (klv *KeyLevelValidator) invokeOnce(block *common.Block, txnum uint64) *sync.Once {
	klv.blockDep.mutex.Lock()
	defer klv.blockDep.mutex.Unlock()

	if klv.blockDep.blockNum != block.Header.Number {
		klv.blockDep.blockNum = block.Header.Number
		klv.blockDep.txDepOnce = make([]sync.Once, len(block.Data.Data))
	}

	return &klv.blockDep.txDepOnce[txnum]
}

func (klv *KeyLevelValidator) extractDependenciesForTx(blockNum, txNum uint64, envelopeBytes []byte) {
	env, err := utils.GetEnvelopeFromBlock(envelopeBytes)
	if err != nil {
		logger.Warningf("while executing GetEnvelopeFromBlock got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	payl, err := utils.GetPayload(env)
	if err != nil {
		logger.Warningf("while executing GetPayload got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	tx, err := utils.GetTransaction(payl.Data)
	if err != nil {
		logger.Warningf("while executing GetTransaction got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	cap, err := utils.GetChaincodeActionPayload(tx.Actions[0].Payload)
	if err != nil {
		logger.Warningf("while executing GetChaincodeActionPayload got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	pRespPayload, err := utils.GetProposalResponsePayload(cap.Action.ProposalResponsePayload)
	if err != nil {
		logger.Warningf("while executing GetProposalResponsePayload got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	respPayload, err := utils.GetChaincodeAction(pRespPayload.Extension)
	if err != nil {
		logger.Warningf("while executing GetChaincodeAction got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	klv.vpmgr.ExtractValidationParameterDependency(blockNum, txNum, respPayload.Results)
}


func (klv *KeyLevelValidator) PreValidate(txNum uint64, block *common.Block) {
	for i := int64(txNum); i >= 0; i-- {
		txPosition := uint64(i)

		klv.invokeOnce(block, txPosition).Do(
			func() {
				klv.extractDependenciesForTx(block.Header.Number, txPosition, block.Data.Data[txPosition])
			})
	}
}


func (klv *KeyLevelValidator) Validate(cc string, blockNum, txNum uint64, rwsetBytes, prp, ccEP []byte, endorsements []*peer.Endorsement) commonerrors.TxValidationError {
	
	signatureSet := []*common.SignedData{}
	for _, endorsement := range endorsements {
		data := make([]byte, len(prp)+len(endorsement.Endorser))
		copy(data, prp)
		copy(data[len(prp):], endorsement.Endorser)

		signatureSet = append(signatureSet, &common.SignedData{
			
			Data: data,
			
			Identity: endorsement.Endorser,
			
			Signature: endorsement.Signature})
	}

	
	policyChecker := policyChecker{
		ccEP:          ccEP,
		policySupport: klv.policySupport,
		signatureSet:  signatureSet,
		vpmgr:         klv.vpmgr,
	}

	
	rwset := &rwsetutil.TxRwSet{}
	if err := rwset.FromProtoBytes(rwsetBytes); err != nil {
		return policyErr(errors.WithMessage(err, fmt.Sprintf("txRWSet.FromProtoBytes failed on tx (%d,%d)", blockNum, txNum)))
	}

	
	for _, nsRWSet := range rwset.NsRwSets {
		
		if nsRWSet.NameSpace != cc {
			continue
		}

		
		
		
		for _, pubWrite := range nsRWSet.KvRwSet.Writes {
			err := policyChecker.checkSBAndCCEP(cc, "", pubWrite.Key, blockNum, txNum)
			if err != nil {
				return err
			}
		}
		
		
		
		for _, pubMdWrite := range nsRWSet.KvRwSet.MetadataWrites {
			err := policyChecker.checkSBAndCCEP(cc, "", pubMdWrite.Key, blockNum, txNum)
			if err != nil {
				return err
			}
		}
		
		
		
		for _, collRWSet := range nsRWSet.CollHashedRwSets {
			coll := collRWSet.CollectionName
			for _, hashedWrite := range collRWSet.HashedRwSet.HashedWrites {
				key := string(hashedWrite.KeyHash)
				err := policyChecker.checkSBAndCCEP(cc, coll, key, blockNum, txNum)
				if err != nil {
					return err
				}
			}
		}
		
		
		
		for _, collRWSet := range nsRWSet.CollHashedRwSets {
			coll := collRWSet.CollectionName
			for _, hashedMdWrite := range collRWSet.HashedRwSet.MetadataWrites {
				key := string(hashedMdWrite.KeyHash)
				err := policyChecker.checkSBAndCCEP(cc, coll, key, blockNum, txNum)
				if err != nil {
					return err
				}
			}
		}
	}

	
	return policyChecker.checkCCEPIfNoEPChecked(cc, blockNum, txNum)
}


func (klv *KeyLevelValidator) PostValidate(cc string, blockNum, txNum uint64, err error) {
	klv.vpmgr.SetTxValidationResult(cc, blockNum, txNum, err)
}

func policyErr(err error) *commonerrors.VSCCEndorsementPolicyError {
	return &commonerrors.VSCCEndorsementPolicyError{
		Err: err,
	}
}
