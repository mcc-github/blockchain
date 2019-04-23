/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valinforetriever_test

import (
	"testing"

	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/valinforetriever"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/valinforetriever/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidationInfoRetrieverFromNew(t *testing.T) {
	cc := "cc"
	newPlugin := "new"
	newArgs := []byte("new")
	uerr := errors.New("unexpected error")
	verr := errors.New("validation error")

	new := &mocks.LifecycleResources{}
	legacy := &mocks.LifecycleResources{}
	shim := valinforetriever.ValidationInfoRetrieveShim{
		Legacy: legacy,
		New:    new,
	}

	
	new.On("ValidationInfo", "channel", cc, nil).Return(newPlugin, newArgs, nil, nil).Once()
	plugin, args, unexpectedErr, validationErr := shim.ValidationInfo("channel", "cc", nil)
	assert.NoError(t, unexpectedErr)
	assert.NoError(t, validationErr)
	assert.Equal(t, "new", plugin)
	assert.Equal(t, []byte("new"), args)
	legacy.AssertNotCalled(t, "ValidationInfo")

	
	new.On("ValidationInfo", "channel", cc, nil).Return("", nil, nil, verr).Once()
	plugin, args, unexpectedErr, validationErr = shim.ValidationInfo("channel", "cc", nil)
	assert.NoError(t, unexpectedErr)
	assert.Error(t, validationErr)
	assert.Contains(t, validationErr.Error(), "validation error")
	assert.Equal(t, "", plugin)
	assert.Equal(t, []byte(nil), args)
	legacy.AssertNotCalled(t, "ValidationInfo")

	
	new.On("ValidationInfo", "channel", cc, nil).Return("", nil, uerr, verr).Once()
	plugin, args, unexpectedErr, validationErr = shim.ValidationInfo("channel", "cc", nil)
	assert.Error(t, unexpectedErr)
	assert.Error(t, validationErr)
	assert.Contains(t, unexpectedErr.Error(), "unexpected error")
	assert.Contains(t, validationErr.Error(), "validation error")
	assert.Equal(t, "", plugin)
	assert.Equal(t, []byte(nil), args)
	legacy.AssertNotCalled(t, "ValidationInfo")
}

func TestValidationInfoRetrieverFromLegacy(t *testing.T) {
	cc := "cc"
	legacyPlugin := "legacy"
	legacyArgs := []byte("legacy")
	uerr := errors.New("unexpected error")
	verr := errors.New("validation error")

	new := &mocks.LifecycleResources{}
	legacy := &mocks.LifecycleResources{}
	shim := valinforetriever.ValidationInfoRetrieveShim{
		Legacy: legacy,
		New:    new,
	}

	
	new.On("ValidationInfo", "channel", cc, nil).Return("", nil, nil, nil)

	
	legacy.On("ValidationInfo", "channel", cc, nil).Return(legacyPlugin, legacyArgs, nil, nil).Once()
	plugin, args, unexpectedErr, validationErr := shim.ValidationInfo("channel", "cc", nil)
	assert.NoError(t, unexpectedErr)
	assert.NoError(t, validationErr)
	assert.Equal(t, "legacy", plugin)
	assert.Equal(t, []byte("legacy"), args)

	
	legacy.On("ValidationInfo", "channel", cc, nil).Return("", nil, nil, verr).Once()
	plugin, args, unexpectedErr, validationErr = shim.ValidationInfo("channel", "cc", nil)
	assert.NoError(t, unexpectedErr)
	assert.Error(t, validationErr)
	assert.Contains(t, validationErr.Error(), "validation error")
	assert.Equal(t, "", plugin)
	assert.Equal(t, []byte(nil), args)

	
	legacy.On("ValidationInfo", "channel", cc, nil).Return("", nil, uerr, nil).Once()
	plugin, args, unexpectedErr, validationErr = shim.ValidationInfo("channel", "cc", nil)
	assert.Error(t, unexpectedErr)
	assert.NoError(t, validationErr)
	assert.Contains(t, unexpectedErr.Error(), "unexpected error")
	assert.Equal(t, "", plugin)
	assert.Equal(t, []byte(nil), args)
}
