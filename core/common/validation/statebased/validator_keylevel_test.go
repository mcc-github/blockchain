/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

type mockPolicyEvaluator struct {
	EvaluateRV          error
	EvaluateResByPolicy map[string]error
}

func (m *mockPolicyEvaluator) Evaluate(policyBytes []byte, signatureSet []*protoutil.SignedData) error {
	if res, ok := m.EvaluateResByPolicy[string(policyBytes)]; ok {
		return res
	}

	return m.EvaluateRV
}

func buildBlockWithTxs(txs ...[]byte) *common.Block {
	return &common.Block{
		Header: &common.BlockHeader{
			Number: 1,
		},
		Data: &common.BlockData{
			Data: txs,
		},
	}
}

func buildTXWithRwset(rws []byte) []byte {
	return protoutil.MarshalOrPanic(&common.Envelope{
		Payload: protoutil.MarshalOrPanic(
			&common.Payload{
				Data: protoutil.MarshalOrPanic(
					&pb.Transaction{
						Actions: []*pb.TransactionAction{
							{
								Payload: protoutil.MarshalOrPanic(&pb.ChaincodeActionPayload{
									Action: &pb.ChaincodeEndorsedAction{
										ProposalResponsePayload: protoutil.MarshalOrPanic(
											&pb.ProposalResponsePayload{
												Extension: protoutil.MarshalOrPanic(&pb.ChaincodeAction{Results: rws}),
											},
										),
									},
								}),
							},
						},
					},
				),
			},
		),
	})
}

func rwsetBytes(t *testing.T, cc string) []byte {
	rwsb := rwsetutil.NewRWSetBuilder()
	rwsb.AddToWriteSet(cc, "key", []byte("value"))
	rws := rwsb.GetTxReadWriteSet()
	rwsetbytes, err := rws.ToProtoBytes()
	assert.NoError(t, err)

	return rwsetbytes
}

func TestKeylevelValidation(t *testing.T) {
	t.Parallel()

	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: []byte("EP")}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: []byte("EP")}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	pe := &mockPolicyEvaluator{}
	validator := NewKeyLevelValidator(NewV13Evaluator(pe, pm), pm)

	rwsb := rwsetBytes(t, "cc")
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	endorsements := []*pb.Endorsement{
		{
			Signature: []byte("signature"),
			Endorser:  []byte("endorser"),
		},
	}

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	err := validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), endorsements)
	assert.NoError(t, err)

	pe.EvaluateRV = fmt.Errorf("policy evaluation error")

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), endorsements)
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCEndorsementPolicyError{}, err)
}

func TestKeylevelValidationPvtData(t *testing.T) {
	t.Parallel()

	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: []byte("EP")}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: []byte("EP")}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	pe := &mockPolicyEvaluator{}
	validator := NewKeyLevelValidator(NewV13Evaluator(pe, pm), pm)

	rwsbu := rwsetutil.NewRWSetBuilder()
	rwsbu.AddToPvtAndHashedWriteSet("cc", "coll", "key", []byte("value"))
	rws := rwsbu.GetTxReadWriteSet()
	rwsb, err := rws.ToProtoBytes()
	assert.NoError(t, err)
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.NoError(t, err)

	pe.EvaluateRV = fmt.Errorf("policy evaluation error")

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCEndorsementPolicyError{}, err)
}

func TestKeylevelValidationMetaUpdate(t *testing.T) {
	t.Parallel()

	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: []byte("EP")}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: []byte("EP")}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	pe := &mockPolicyEvaluator{}
	validator := NewKeyLevelValidator(NewV13Evaluator(pe, pm), pm)

	rwsbu := rwsetutil.NewRWSetBuilder()
	rwsbu.AddToMetadataWriteSet("cc", "key", map[string][]byte{})
	rws := rwsbu.GetTxReadWriteSet()
	rwsb, err := rws.ToProtoBytes()
	assert.NoError(t, err)
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.NoError(t, err)

	pe.EvaluateRV = fmt.Errorf("policy evaluation error")

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCEndorsementPolicyError{}, err)
}

func TestKeylevelValidationPvtMetaUpdate(t *testing.T) {
	t.Parallel()

	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: []byte("EP")}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: []byte("EP")}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	pe := &mockPolicyEvaluator{}
	validator := NewKeyLevelValidator(NewV13Evaluator(pe, pm), pm)

	rwsbu := rwsetutil.NewRWSetBuilder()
	rwsbu.AddToHashedMetadataWriteSet("cc", "coll", "key", map[string][]byte{})
	rws := rwsbu.GetTxReadWriteSet()
	rwsb, err := rws.ToProtoBytes()
	assert.NoError(t, err)
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.NoError(t, err)

	pe.EvaluateRV = fmt.Errorf("policy evaluation error")

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCEndorsementPolicyError{}, err)
}

func TestKeylevelValidationPolicyRetrievalFailure(t *testing.T) {
	t.Parallel()

	
	
	
	

	mr := &mockState{GetStateMetadataErr: fmt.Errorf("metadata retrieval failure")}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	validator := NewKeyLevelValidator(NewV13Evaluator(&mockPolicyEvaluator{}, pm), pm)

	rwsb := rwsetBytes(t, "cc")
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	err := validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCExecutionFailureError{}, err)
}

func TestKeylevelValidationLedgerFailures(t *testing.T) {
	
	
	
	
	

	rwsb := rwsetBytes(t, "cc")
	prp := []byte("barf")

	t.Run("CollConfigNotDefinedError", func(t *testing.T) {
		mr := &mockState{GetStateMetadataErr: &ledger.CollConfigNotDefinedError{Ns: "mycc"}}
		ms := &mockStateFetcher{FetchStateRv: mr}
		pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
		validator := NewKeyLevelValidator(NewV13Evaluator(&mockPolicyEvaluator{}, pm), pm)

		err := validator.Validate("cc", 1, 0, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
		assert.NoError(t, err)
	})

	t.Run("InvalidCollNameError", func(t *testing.T) {
		mr := &mockState{GetStateMetadataErr: &ledger.InvalidCollNameError{Ns: "mycc", Coll: "mycoll"}}
		ms := &mockStateFetcher{FetchStateRv: mr}
		pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
		validator := NewKeyLevelValidator(NewV13Evaluator(&mockPolicyEvaluator{}, pm), pm)

		err := validator.Validate("cc", 1, 0, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
		assert.NoError(t, err)
	})

	t.Run("I/O error", func(t *testing.T) {
		mr := &mockState{GetStateMetadataErr: fmt.Errorf("some I/O error")}
		ms := &mockStateFetcher{FetchStateRv: mr}
		pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
		validator := NewKeyLevelValidator(NewV13Evaluator(&mockPolicyEvaluator{}, pm), pm)

		err := validator.Validate("cc", 1, 0, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
		assert.Error(t, err)
		assert.IsType(t, &errors.VSCCExecutionFailureError{}, err)
	})
}

func TestCCEPValidation(t *testing.T) {
	t.Parallel()

	
	
	

	mr := &mockState{GetStateMetadataRv: map[string][]byte{}, GetPrivateDataMetadataByHashRv: map[string][]byte{}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	pe := &mockPolicyEvaluator{}
	validator := NewKeyLevelValidator(NewV13Evaluator(pe, pm), pm)

	rwsbu := rwsetutil.NewRWSetBuilder()
	rwsbu.AddToWriteSet("cc", "key", []byte("value"))
	rwsbu.AddToWriteSet("cc", "key1", []byte("value"))
	rwsbu.AddToReadSet("cc", "readkey", &version.Height{})
	rwsbu.AddToHashedReadSet("cc", "coll", "readpvtkey", &version.Height{})
	rws := rwsbu.GetTxReadWriteSet()
	rwsb, err := rws.ToProtoBytes()
	assert.NoError(t, err)
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.NoError(t, err)

	pe.EvaluateRV = fmt.Errorf("policy evaluation error")

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCEndorsementPolicyError{}, err)
}

func TestCCEPValidationReads(t *testing.T) {
	t.Parallel()

	
	
	

	mr := &mockState{GetStateMetadataRv: map[string][]byte{}, GetPrivateDataMetadataByHashRv: map[string][]byte{}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	pe := &mockPolicyEvaluator{}
	validator := NewKeyLevelValidator(NewV13Evaluator(pe, pm), pm)

	rwsbu := rwsetutil.NewRWSetBuilder()
	rwsbu.AddToReadSet("cc", "readkey", &version.Height{})
	rws := rwsbu.GetTxReadWriteSet()
	rwsb, err := rws.ToProtoBytes()
	assert.NoError(t, err)
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.NoError(t, err)

	pe.EvaluateRV = fmt.Errorf("policy evaluation error")

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCEndorsementPolicyError{}, err)
}

func TestOnlySBEPChecked(t *testing.T) {
	t.Parallel()

	
	
	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: []byte("SBEP")}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	pe := &mockPolicyEvaluator{}
	validator := NewKeyLevelValidator(NewV13Evaluator(pe, pm), pm)

	rwsb := rwsetBytes(t, "cc")
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	pe.EvaluateRV = fmt.Errorf("policy evaluation error")
	pe.EvaluateResByPolicy = map[string]error{
		"SBEP": nil,
	}

	err := validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.NoError(t, err)

	
	rwsbu := rwsetutil.NewRWSetBuilder()
	rwsbu.AddToWriteSet("cc", "key", []byte("value"))
	rwsbu.AddToReadSet("cc", "key", nil)
	rws := rwsbu.GetTxReadWriteSet()
	rwsb, _ = rws.ToProtoBytes()

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.NoError(t, err)
}

func TestCCEPValidationPvtReads(t *testing.T) {
	t.Parallel()

	
	
	

	mr := &mockState{GetStateMetadataRv: map[string][]byte{}, GetPrivateDataMetadataByHashRv: map[string][]byte{}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	pe := &mockPolicyEvaluator{}
	validator := NewKeyLevelValidator(NewV13Evaluator(pe, pm), pm)

	rwsbu := rwsetutil.NewRWSetBuilder()
	rwsbu.AddToHashedReadSet("cc", "coll", "readpvtkey", &version.Height{})
	rws := rwsbu.GetTxReadWriteSet()
	rwsb, err := rws.ToProtoBytes()
	assert.NoError(t, err)
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, fmt.Errorf(""))
	}()

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.NoError(t, err)

	pe.EvaluateRV = fmt.Errorf("policy evaluation error")

	err = validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCEndorsementPolicyError{}, err)
}

func TestKeylevelValidationFailure(t *testing.T) {
	t.Parallel()

	
	
	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: []byte("EP")}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: []byte("EP")}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{PolicyTranslator: &mockTranslator{}, StateFetcher: ms}
	validator := NewKeyLevelValidator(NewV13Evaluator(&mockPolicyEvaluator{}, pm), pm)

	rwsb := rwsetBytes(t, "cc")
	prp := []byte("barf")
	block := buildBlockWithTxs(buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")), buildTXWithRwset(rwsetUpdatingMetadataFor("cc", "key")))

	validator.PreValidate(1, block)

	go func() {
		validator.PostValidate("cc", 1, 0, nil)
	}()

	err := validator.Validate("cc", 1, 1, rwsb, prp, []byte("CCEP"), []*pb.Endorsement{})
	assert.Error(t, err)
	assert.IsType(t, &errors.VSCCEndorsementPolicyError{}, err)
}
