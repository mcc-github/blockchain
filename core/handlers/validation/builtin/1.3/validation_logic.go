/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package builtin1_3

import (
	"fmt"
	"regexp"

	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/ccmetadata"
	. "github.com/mcc-github/blockchain/core/handlers/validation/api/capabilities"
	. "github.com/mcc-github/blockchain/core/handlers/validation/api/identities"
	. "github.com/mcc-github/blockchain/core/handlers/validation/api/policies"
	. "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("vscc")

const (
	DUPLICATED_IDENTITY_ERROR = "Endorsement policy evaluation failure might be caused by duplicated identities"
)

var validCollectionNameRegex = regexp.MustCompile(ccmetadata.AllowedCharsCollectionName)








func New(c Capabilities, s StateFetcher, d IdentityDeserializer, pe PolicyEvaluator) *Validator {
	return &Validator{
		capabilities:    c,
		stateFetcher:    s,
		deserializer:    d,
		policyEvaluator: pe,
	}
}





type Validator struct {
	deserializer    IdentityDeserializer
	capabilities    Capabilities
	stateFetcher    StateFetcher
	policyEvaluator PolicyEvaluator
}



func (vscc *Validator) Validate(
	block *common.Block,
	namespace string,
	txPosition int,
	actionPosition int,
	policyBytes []byte,
) commonerrors.TxValidationError {
	
	env, err := utils.GetEnvelopeFromBlock(block.Data.Data[txPosition])
	if err != nil {
		logger.Errorf("VSCC error: GetEnvelope failed, err %s", err)
		return policyErr(err)
	}

	
	payl, err := utils.GetPayload(env)
	if err != nil {
		logger.Errorf("VSCC error: GetPayload failed, err %s", err)
		return policyErr(err)
	}

	chdr, err := utils.UnmarshalChannelHeader(payl.Header.ChannelHeader)
	if err != nil {
		return policyErr(err)
	}

	
	if common.HeaderType(chdr.Type) != common.HeaderType_ENDORSER_TRANSACTION {
		logger.Errorf("Only Endorser Transactions are supported, provided type %d", chdr.Type)
		return policyErr(fmt.Errorf("Only Endorser Transactions are supported, provided type %d", chdr.Type))
	}

	
	tx, err := utils.GetTransaction(payl.Data)
	if err != nil {
		logger.Errorf("VSCC error: GetTransaction failed, err %s", err)
		return policyErr(err)
	}

	cap, err := utils.GetChaincodeActionPayload(tx.Actions[actionPosition].Payload)
	if err != nil {
		logger.Errorf("VSCC error: GetChaincodeActionPayload failed, err %s", err)
		return policyErr(err)
	}

	
	signatureSet := []*common.SignedData{}
	for _, endorsement := range cap.Action.Endorsements {
		signatureSet = append(signatureSet, &common.SignedData{
			
			Data: append(append([]byte{}, cap.Action.ProposalResponsePayload...), endorsement.Endorser...),
			
			Identity: endorsement.Endorser,
			
			Signature: endorsement.Signature})
	}

	
	err = vscc.policyEvaluator.Evaluate(policyBytes, signatureSet)
	if err != nil {
		logger.Warningf("Endorsement policy failure for transaction txid=%s, err: %s", chdr.GetTxId(), err.Error())
		if len(signatureSet) < len(cap.Action.Endorsements) {
			
			return policyErr(errors.New(DUPLICATED_IDENTITY_ERROR))
		}
		return policyErr(fmt.Errorf("VSCC error: endorsement policy failure, err: %s", err))
	}

	
	if namespace == "lscc" {
		logger.Debugf("VSCC info: doing special validation for LSCC")
		err := vscc.ValidateLSCCInvocation(chdr.ChannelId, env, cap, payl, vscc.capabilities)
		if err != nil {
			logger.Errorf("VSCC error: ValidateLSCCInvocation failed, err %s", err)
			return err
		}
	}

	return nil
}

func policyErr(err error) *commonerrors.VSCCEndorsementPolicyError {
	return &commonerrors.VSCCEndorsementPolicyError{
		Err: err,
	}
}
