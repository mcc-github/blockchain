/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api/policies"
	s "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)






func NewV20Evaluator(
	vpmgr KeyLevelValidationParameterManager,
	policySupport validation.PolicyEvaluator,
	collRes CollectionResources,
	StateFetcher s.StateFetcher,
) *policyCheckerFactoryV20 {
	return &policyCheckerFactoryV20{
		vpmgr:         vpmgr,
		policySupport: policySupport,
		StateFetcher:  StateFetcher,
		collRes:       collRes,
	}
}

type policyCheckerFactoryV20 struct {
	vpmgr         KeyLevelValidationParameterManager
	policySupport validation.PolicyEvaluator
	collRes       CollectionResources
	StateFetcher  s.StateFetcher
}

func (p *policyCheckerFactoryV20) Evaluator(ccEP []byte) RWSetPolicyEvaluator {
	return &baseEvaluator{
		epEvaluator: &policyCheckerV20{
			ccEP:          ccEP,
			policySupport: p.policySupport,
			nsEPChecked:   map[string]bool{},
			collRes:       p.collRes,
			StateFetcher:  p.StateFetcher,
		},
		vpmgr:         p.vpmgr,
		policySupport: p.policySupport,
	}
}





type CollectionResources interface {
	
	
	
	
	
	
	CollectionValidationInfo(chaincodeName, collectionName string, state s.State) (args []byte, unexpectedErr error, validationErr error)
}




type policyCheckerV20 struct {
	someEPChecked bool
	vpmgr         KeyLevelValidationParameterManager
	policySupport validation.PolicyEvaluator
	ccEP          []byte
	nsEPChecked   map[string]bool
	collRes       CollectionResources
	StateFetcher  s.StateFetcher
}

func (p *policyCheckerV20) fetchCollEP(cc, coll string) ([]byte, commonerrors.TxValidationError) {
	state, err := p.StateFetcher.FetchState()
	if err != nil {
		return nil, &commonerrors.VSCCExecutionFailureError{
			Err: errors.WithMessage(err, "could not retrieve ledger"),
		}
	}
	defer state.Done()

	collEP, unexpectedErr, validationErr := p.collRes.CollectionValidationInfo(cc, coll, state)
	if unexpectedErr != nil {
		return nil, &commonerrors.VSCCExecutionFailureError{
			Err: unexpectedErr,
		}
	}
	if validationErr != nil {
		return nil, policyErr(validationErr)
	}

	return collEP, nil
}

func (p *policyCheckerV20) CheckCCEPIfNotChecked(cc, coll string, blockNum, txNum uint64, sd []*protoutil.SignedData) commonerrors.TxValidationError {
	if coll != "" {
		
		
		if p.nsEPChecked[coll] {
			return nil
		}

		
		collEP, err := p.fetchCollEP(cc, coll)
		if err != nil {
			return err
		}

		
		if len(collEP) != 0 {
			err := p.policySupport.Evaluate(collEP, sd)
			if err != nil {
				return policyErr(errors.Wrapf(err, "validation of endorsement policy for collection %s chaincode %s in tx %d:%d failed", coll, cc, blockNum, txNum))
			}

			p.nsEPChecked[coll] = true
			p.someEPChecked = true
			return nil
		}
	}

	
	
	if p.nsEPChecked[""] {
		return nil
	}

	
	err := p.policySupport.Evaluate(p.ccEP, sd)
	if err != nil {
		return policyErr(errors.Wrapf(err, "validation of endorsement policy for chaincode %s in tx %d:%d failed", cc, blockNum, txNum))
	}

	p.nsEPChecked[""] = true
	p.someEPChecked = true
	return nil
}

func (p *policyCheckerV20) CheckCCEPIfNoEPChecked(cc string, blockNum, txNum uint64, sd []*protoutil.SignedData) commonerrors.TxValidationError {
	if p.someEPChecked {
		return nil
	}

	
	err := p.policySupport.Evaluate(p.ccEP, sd)
	if err != nil {
		return policyErr(errors.Wrapf(err, "validation of endorsement policy for chaincode %s in tx %d:%d failed", cc, blockNum, txNum))
	}

	p.nsEPChecked[""] = true
	p.someEPChecked = true
	return nil
}

func (p *policyCheckerV20) SBEPChecked() {
	p.someEPChecked = true
}
