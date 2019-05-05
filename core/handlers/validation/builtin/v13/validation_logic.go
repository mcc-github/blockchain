/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package v13

import (
	"fmt"
	"regexp"

	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/ccmetadata"
	"github.com/mcc-github/blockchain/core/common/validation/statebased"
	vc "github.com/mcc-github/blockchain/core/handlers/validation/api/capabilities"
	vi "github.com/mcc-github/blockchain/core/handlers/validation/api/identities"
	vp "github.com/mcc-github/blockchain/core/handlers/validation/api/policies"
	vs "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
)

var logger = flogging.MustGetLogger("vscc")

var validCollectionNameRegex = regexp.MustCompile(ccmetadata.AllowedCharsCollectionName)




type Capabilities interface {
	vc.Capabilities
}




type StateFetcher interface {
	vs.StateFetcher
}




type IdentityDeserializer interface {
	vi.IdentityDeserializer
}




type PolicyEvaluator interface {
	vp.PolicyEvaluator
}








type noopTranslator struct{}

func (n *noopTranslator) Translate(b []byte) ([]byte, error) {
	return b, nil
}



func New(c vc.Capabilities, s vs.StateFetcher, d vi.IdentityDeserializer, pe vp.PolicyEvaluator) *Validator {
	vpmgr := &statebased.KeyLevelValidationParameterManagerImpl{
		StateFetcher:     s,
		PolicyTranslator: &noopTranslator{},
	}
	eval := statebased.NewV13Evaluator(pe, vpmgr)
	sbv := statebased.NewKeyLevelValidator(eval, vpmgr)

	return &Validator{
		capabilities:        c,
		stateFetcher:        s,
		deserializer:        d,
		policyEvaluator:     pe,
		stateBasedValidator: sbv,
	}
}





type Validator struct {
	deserializer        vi.IdentityDeserializer
	capabilities        vc.Capabilities
	stateFetcher        vs.StateFetcher
	policyEvaluator     vp.PolicyEvaluator
	stateBasedValidator StateBasedValidator
}

type validationArtifacts struct {
	rwset        []byte
	prp          []byte
	endorsements []*peer.Endorsement
	chdr         *common.ChannelHeader
	env          *common.Envelope
	payl         *common.Payload
	cap          *peer.ChaincodeActionPayload
}

func (vscc *Validator) extractValidationArtifacts(
	block *common.Block,
	txPosition int,
	actionPosition int,
) (*validationArtifacts, error) {
	
	env, err := protoutil.GetEnvelopeFromBlock(block.Data.Data[txPosition])
	if err != nil {
		logger.Errorf("VSCC error: GetEnvelope failed, err %s", err)
		return nil, err
	}

	
	payl, err := protoutil.GetPayload(env)
	if err != nil {
		logger.Errorf("VSCC error: GetPayload failed, err %s", err)
		return nil, err
	}

	chdr, err := protoutil.UnmarshalChannelHeader(payl.Header.ChannelHeader)
	if err != nil {
		return nil, err
	}

	
	if common.HeaderType(chdr.Type) != common.HeaderType_ENDORSER_TRANSACTION {
		logger.Errorf("Only Endorser Transactions are supported, provided type %d", chdr.Type)
		err = fmt.Errorf("Only Endorser Transactions are supported, provided type %d", chdr.Type)
		return nil, err
	}

	
	tx, err := protoutil.GetTransaction(payl.Data)
	if err != nil {
		logger.Errorf("VSCC error: GetTransaction failed, err %s", err)
		return nil, err
	}

	cap, err := protoutil.GetChaincodeActionPayload(tx.Actions[actionPosition].Payload)
	if err != nil {
		logger.Errorf("VSCC error: GetChaincodeActionPayload failed, err %s", err)
		return nil, err
	}

	pRespPayload, err := protoutil.GetProposalResponsePayload(cap.Action.ProposalResponsePayload)
	if err != nil {
		err = fmt.Errorf("GetProposalResponsePayload error %s", err)
		return nil, err
	}
	if pRespPayload.Extension == nil {
		err = fmt.Errorf("nil pRespPayload.Extension")
		return nil, err
	}
	respPayload, err := protoutil.GetChaincodeAction(pRespPayload.Extension)
	if err != nil {
		err = fmt.Errorf("GetChaincodeAction error %s", err)
		return nil, err
	}

	return &validationArtifacts{
		rwset:        respPayload.Results,
		prp:          cap.Action.ProposalResponsePayload,
		endorsements: cap.Action.Endorsements,
		chdr:         chdr,
		env:          env,
		payl:         payl,
		cap:          cap,
	}, nil
}







func (vscc *Validator) Validate(
	block *common.Block,
	namespace string,
	txPosition int,
	actionPosition int,
	policyBytes []byte,
) commonerrors.TxValidationError {
	vscc.stateBasedValidator.PreValidate(uint64(txPosition), block)

	va, err := vscc.extractValidationArtifacts(block, txPosition, actionPosition)
	if err != nil {
		vscc.stateBasedValidator.PostValidate(namespace, block.Header.Number, uint64(txPosition), err)
		return policyErr(err)
	}

	txverr := vscc.stateBasedValidator.Validate(
		namespace,
		block.Header.Number,
		uint64(txPosition),
		va.rwset,
		va.prp,
		policyBytes,
		va.endorsements,
	)
	if txverr != nil {
		logger.Errorf("VSCC error: stateBasedValidator.Validate failed, err %s", txverr)
		vscc.stateBasedValidator.PostValidate(namespace, block.Header.Number, uint64(txPosition), txverr)
		return txverr
	}

	
	if namespace == "lscc" {
		logger.Debugf("VSCC info: doing special validation for LSCC")
		err := vscc.ValidateLSCCInvocation(va.chdr.ChannelId, va.env, va.cap, va.payl, vscc.capabilities)
		if err != nil {
			logger.Errorf("VSCC error: ValidateLSCCInvocation failed, err %s", err)
			vscc.stateBasedValidator.PostValidate(namespace, block.Header.Number, uint64(txPosition), err)
			return err
		}
	}

	vscc.stateBasedValidator.PostValidate(namespace, block.Header.Number, uint64(txPosition), nil)
	return nil
}

func policyErr(err error) *commonerrors.VSCCEndorsementPolicyError {
	return &commonerrors.VSCCEndorsementPolicyError{
		Err: err,
	}
}
