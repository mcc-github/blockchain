/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package builtin

import (
	"testing"

	"github.com/mcc-github/blockchain-protos-go/common"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	. "github.com/mcc-github/blockchain/core/handlers/validation/api"
	vmocks "github.com/mcc-github/blockchain/core/handlers/validation/builtin/mocks"
	"github.com/mcc-github/blockchain/core/handlers/validation/builtin/v12/mocks"
	v20mocks "github.com/mcc-github/blockchain/core/handlers/validation/builtin/v20/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInit(t *testing.T) {
	factory := &DefaultValidationFactory{}
	defValidation := factory.New()

	identityDeserializer := &mocks.IdentityDeserializer{}
	capabilities := &mocks.Capabilities{}
	stateFetcher := &mocks.StateFetcher{}
	polEval := &mocks.PolicyEvaluator{}
	colRes := &v20mocks.CollectionResources{}

	assert.Equal(t, "stateFetcher not passed in init", defValidation.Init(identityDeserializer, capabilities, polEval, colRes).Error())
	assert.Equal(t, "identityDeserializer not passed in init", defValidation.Init(capabilities, stateFetcher, polEval, colRes).Error())
	assert.Equal(t, "capabilities not passed in init", defValidation.Init(identityDeserializer, stateFetcher, polEval, colRes).Error())
	assert.Equal(t, "policy fetcher not passed in init", defValidation.Init(identityDeserializer, capabilities, stateFetcher, colRes).Error())
	assert.Equal(t, "collection resources not passed in init", defValidation.Init(identityDeserializer, capabilities, stateFetcher, polEval).Error())

	fullDeps := []Dependency{identityDeserializer, capabilities, stateFetcher, polEval, colRes}
	assert.NoError(t, defValidation.Init(fullDeps...))
}

func TestErrorConversion(t *testing.T) {
	validator := &vmocks.TransactionValidator{}
	capabilities := &mocks.Capabilities{}
	validation := &DefaultValidation{
		TxValidatorV1_2: validator,
		Capabilities:    capabilities,
	}
	block := &common.Block{
		Header: &common.BlockHeader{},
		Data: &common.BlockData{
			Data: [][]byte{{}},
		},
	}

	capabilities.On("V2_0Validation").Return(false)
	capabilities.On("V1_3Validation").Return(false)
	capabilities.On("V1_2Validation").Return(true)

	
	
	validator.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("bla bla")).Once()
	assert.Panics(t, func() {
		validation.Validate(block, "", 0, 0, plugin.SerializedPolicy("policy"))
	})

	
	validator.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&commonerrors.VSCCEndorsementPolicyError{Err: errors.New("foo")}).Once()
	err := validation.Validate(block, "", 0, 0, plugin.SerializedPolicy("policy"))
	assert.Equal(t, (&commonerrors.VSCCEndorsementPolicyError{Err: errors.New("foo")}).Error(), err.Error())

	
	validator.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&commonerrors.VSCCExecutionFailureError{Err: errors.New("bar")}).Once()
	err = validation.Validate(block, "", 0, 0, plugin.SerializedPolicy("policy"))
	assert.Equal(t, &ExecutionFailureError{Reason: "bar"}, err)

	
	validator.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	assert.NoError(t, validation.Validate(block, "", 0, 0, plugin.SerializedPolicy("policy")))
}

func TestValidateBadInput(t *testing.T) {
	validator := &vmocks.TransactionValidator{}
	validation := &DefaultValidation{
		TxValidatorV1_2: validator,
	}

	
	validator.On("Validate", mock.Anything, mock.Anything).Return(nil).Once()
	err := validation.Validate(nil, "", 0, 0, plugin.SerializedPolicy("policy"))
	assert.Equal(t, "empty block", err.Error())

	block := &common.Block{
		Header: &common.BlockHeader{},
		Data: &common.BlockData{
			Data: [][]byte{{}},
		},
	}
	
	validator.On("Validate", mock.Anything, mock.Anything).Return(nil).Once()
	err = validation.Validate(block, "", 1, 0, plugin.SerializedPolicy("policy"))
	assert.Equal(t, "block has only 1 transactions, but requested tx at position 1", err.Error())

	
	validator.On("Validate", mock.Anything, mock.Anything).Return(nil).Once()
	err = validation.Validate(&common.Block{
		Data: &common.BlockData{
			Data: [][]byte{{}},
		},
	}, "", 0, 0, plugin.SerializedPolicy("policy"))
	assert.Equal(t, "no block header", err.Error())

	
	assert.Panics(t, func() {
		validator.On("Validate", mock.Anything, mock.Anything).Return(nil).Once()
		err = validation.Validate(&common.Block{
			Header: &common.BlockHeader{},
			Data: &common.BlockData{
				Data: [][]byte{{}},
			},
		}, "", 0, 0)
	})

	
	assert.Panics(t, func() {
		validator.On("Validate", mock.Anything, mock.Anything).Return(nil).Once()
		err = validation.Validate(&common.Block{
			Header: &common.BlockHeader{},
			Data: &common.BlockData{
				Data: [][]byte{{}},
			},
		}, "", 0, 0, []byte("policy"))
	})

}
