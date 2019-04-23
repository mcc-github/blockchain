/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"sync"

	commonerrors "github.com/mcc-github/blockchain/common/errors"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api/policies"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

type epEvaluator interface {
	CheckCCEPIfNotChecked(cc, coll string, blockNum, txNum uint64, sd []*protoutil.SignedData) commonerrors.TxValidationError
	CheckCCEPIfNoEPChecked(cc string, blockNum, txNum uint64, sd []*protoutil.SignedData) commonerrors.TxValidationError
	SBEPChecked()
}




type baseEvaluator struct {
	epEvaluator
	vpmgr         KeyLevelValidationParameterManager
	policySupport validation.PolicyEvaluator
}

func (p *baseEvaluator) checkSBAndCCEP(cc, coll, key string, blockNum, txNum uint64, signatureSet []*protoutil.SignedData) commonerrors.TxValidationError {
	
	vp, err := p.vpmgr.GetValidationParameterForKey(cc, coll, key, blockNum, txNum)
	if err != nil {
		
		switch err := errors.Cause(err).(type) {
		
		
		
		case *ValidationParameterUpdatedError:
			return policyErr(err)
		
		
		
		
		
		case *ledger.CollConfigNotDefinedError, *ledger.InvalidCollNameError:
			logger.Warningf(errors.WithMessage(err, "skipping key-level validation").Error())
			err = nil
		
		
		
		
		
		
		
		
		default:
			return &commonerrors.VSCCExecutionFailureError{
				Err: err,
			}
		}
	}

	
	if len(vp) == 0 {
		return p.CheckCCEPIfNotChecked(cc, coll, blockNum, txNum, signatureSet)
	}

	
	err = p.policySupport.Evaluate(vp, signatureSet)
	if err != nil {
		return policyErr(errors.Wrapf(err, "validation of key %s (coll'%s':ns'%s') in tx %d:%d failed", key, coll, cc, blockNum, txNum))
	}

	p.SBEPChecked()

	return nil
}

func (p *baseEvaluator) Evaluate(blockNum, txNum uint64, NsRwSets []*rwsetutil.NsRwSet, ns string, sd []*protoutil.SignedData) commonerrors.TxValidationError {
	
	for _, nsRWSet := range NsRwSets {
		
		if nsRWSet.NameSpace != ns {
			continue
		}

		
		
		
		for _, pubWrite := range nsRWSet.KvRwSet.Writes {
			err := p.checkSBAndCCEP(ns, "", pubWrite.Key, blockNum, txNum, sd)
			if err != nil {
				return err
			}
		}
		
		
		
		for _, pubMdWrite := range nsRWSet.KvRwSet.MetadataWrites {
			err := p.checkSBAndCCEP(ns, "", pubMdWrite.Key, blockNum, txNum, sd)
			if err != nil {
				return err
			}
		}
		
		
		
		for _, collRWSet := range nsRWSet.CollHashedRwSets {
			coll := collRWSet.CollectionName
			for _, hashedWrite := range collRWSet.HashedRwSet.HashedWrites {
				key := string(hashedWrite.KeyHash)
				err := p.checkSBAndCCEP(ns, coll, key, blockNum, txNum, sd)
				if err != nil {
					return err
				}
			}
		}
		
		
		
		for _, collRWSet := range nsRWSet.CollHashedRwSets {
			coll := collRWSet.CollectionName
			for _, hashedMdWrite := range collRWSet.HashedRwSet.MetadataWrites {
				key := string(hashedMdWrite.KeyHash)
				err := p.checkSBAndCCEP(ns, coll, key, blockNum, txNum, sd)
				if err != nil {
					return err
				}
			}
		}
	}

	
	return p.CheckCCEPIfNoEPChecked(ns, blockNum, txNum, sd)
}





type RWSetPolicyEvaluatorFactory interface {
	
	
	Evaluator(ccEP []byte) RWSetPolicyEvaluator
}


type RWSetPolicyEvaluator interface {
	Evaluate(blockNum, txNum uint64, NsRwSets []*rwsetutil.NsRwSet, ns string, sd []*protoutil.SignedData) commonerrors.TxValidationError
}




type blockDependency struct {
	mutex     sync.Mutex
	blockNum  uint64
	txDepOnce []sync.Once
}


type KeyLevelValidator struct {
	vpmgr    KeyLevelValidationParameterManager
	blockDep blockDependency
	pef      RWSetPolicyEvaluatorFactory
}

func NewKeyLevelValidator(evaluator RWSetPolicyEvaluatorFactory, vpmgr KeyLevelValidationParameterManager) *KeyLevelValidator {
	return &KeyLevelValidator{
		vpmgr:    vpmgr,
		blockDep: blockDependency{},
		pef:      evaluator,
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
	env, err := protoutil.GetEnvelopeFromBlock(envelopeBytes)
	if err != nil {
		logger.Warningf("while executing GetEnvelopeFromBlock got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	payl, err := protoutil.GetPayload(env)
	if err != nil {
		logger.Warningf("while executing GetPayload got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	tx, err := protoutil.GetTransaction(payl.Data)
	if err != nil {
		logger.Warningf("while executing GetTransaction got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	cap, err := protoutil.GetChaincodeActionPayload(tx.Actions[0].Payload)
	if err != nil {
		logger.Warningf("while executing GetChaincodeActionPayload got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	pRespPayload, err := protoutil.GetProposalResponsePayload(cap.Action.ProposalResponsePayload)
	if err != nil {
		logger.Warningf("while executing GetProposalResponsePayload got error '%s', skipping tx at height (%d,%d)", err, blockNum, txNum)
		return
	}

	respPayload, err := protoutil.GetChaincodeAction(pRespPayload.Extension)
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
	
	signatureSet := []*protoutil.SignedData{}
	for _, endorsement := range endorsements {
		data := make([]byte, len(prp)+len(endorsement.Endorser))
		copy(data, prp)
		copy(data[len(prp):], endorsement.Endorser)

		signatureSet = append(signatureSet, &protoutil.SignedData{
			
			Data: data,
			
			Identity: endorsement.Endorser,
			
			Signature: endorsement.Signature})
	}

	
	policyEvaluator := klv.pef.Evaluator(ccEP)

	
	rwset := &rwsetutil.TxRwSet{}
	if err := rwset.FromProtoBytes(rwsetBytes); err != nil {
		return policyErr(errors.WithMessagef(err, "txRWSet.FromProtoBytes failed on tx (%d,%d)", blockNum, txNum))
	}

	
	return policyEvaluator.Evaluate(blockNum, txNum, rwset.NsRwSets, cc, signatureSet)
}


func (klv *KeyLevelValidator) PostValidate(cc string, blockNum, txNum uint64, err error) {
	klv.vpmgr.SetTxValidationResult(cc, blockNum, txNum, err)
}

func policyErr(err error) *commonerrors.VSCCEndorsementPolicyError {
	return &commonerrors.VSCCEndorsementPolicyError{
		Err: err,
	}
}
