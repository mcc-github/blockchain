/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugindispatcher_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/mocks/ledger"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher/mocks"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/testdata"
	"github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/msp"
	. "github.com/mcc-github/blockchain/msp/mocks"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateWithPlugin(t *testing.T) {
	pm := make(plugin.MapBasedMapper)
	qec := &mocks.QueryExecutorCreator{}
	deserializer := &mocks.IdentityDeserializer{}
	capabilites := &mocks.Capabilities{}

	mcpmg := &mocks.ChannelPolicyManagerGetter{}
	mcpmg.On("Manager", "").Return(nil, true)

	v := plugindispatcher.NewPluginValidator(pm, qec, deserializer, capabilites, mcpmg, nil)
	ctx := &plugindispatcher.Context{
		Namespace:  "mycc",
		PluginName: "vscc",
	}

	
	err := v.ValidateWithPlugin(ctx)
	assert.Contains(t, err.Error(), "plugin with name vscc wasn't found")

	
	factory := &mocks.PluginFactory{}
	plugin := &mocks.Plugin{}
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("foo")).Once()
	factory.On("New").Return(plugin)
	pm["vscc"] = factory
	err = v.ValidateWithPlugin(ctx)
	assert.Contains(t, err.(*validation.ExecutionFailureError).Error(), "failed initializing plugin: foo")

	
	
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	validationErr := &validation.ExecutionFailureError{
		Reason: "bar",
	}
	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(validationErr).Once()
	err = v.ValidateWithPlugin(ctx)
	assert.Equal(t, validationErr, err)

	
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	err = v.ValidateWithPlugin(ctx)
	assert.NoError(t, err)
}

func TestSamplePlugin(t *testing.T) {
	pm := make(plugin.MapBasedMapper)
	qec := &mocks.QueryExecutorCreator{}

	qec.On("NewQueryExecutor").Return(&ledger.MockQueryExecutor{
		State: map[string]map[string][]byte{
			"lscc": {
				"mycc": []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
	}, nil)

	deserializer := &mocks.IdentityDeserializer{}
	identity := &MockIdentity{}
	identity.On("GetIdentifier").Return(&msp.IdentityIdentifier{
		Mspid: "SampleOrg",
		Id:    "foo",
	})
	deserializer.On("DeserializeIdentity", []byte{7, 8, 9}).Return(identity, nil)
	capabilites := &mocks.Capabilities{}
	capabilites.On("PrivateChannelData").Return(true)
	factory := &mocks.PluginFactory{}
	factory.On("New").Return(testdata.NewSampleValidationPlugin(t))
	pm["vscc"] = factory

	mcpmg := &mocks.ChannelPolicyManagerGetter{}
	mcpmg.On("Manager", "mychannel").Return(nil, true)

	transaction := testdata.MarshaledSignedData{
		Data:      []byte{1, 2, 3},
		Signature: []byte{4, 5, 6},
		Identity:  []byte{7, 8, 9},
	}

	txnData, _ := proto.Marshal(&transaction)

	v := plugindispatcher.NewPluginValidator(pm, qec, deserializer, capabilites, mcpmg, nil)
	acceptAllPolicyBytes, _ := proto.Marshal(&peer.ApplicationPolicy{Type: &peer.ApplicationPolicy_SignaturePolicy{SignaturePolicy: cauthdsl.AcceptAllPolicy}})
	ctx := &plugindispatcher.Context{
		Namespace:  "mycc",
		PluginName: "vscc",
		Policy:     acceptAllPolicyBytes,
		Block: &common.Block{
			Header: &common.BlockHeader{},
			Data: &common.BlockData{
				Data: [][]byte{txnData},
			},
		},
		Channel: "mychannel",
	}
	assert.NoError(t, v.ValidateWithPlugin(ctx))
}
