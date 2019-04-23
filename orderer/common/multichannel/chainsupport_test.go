/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"testing"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/deliver/mock"
	"github.com/mcc-github/blockchain/common/ledger/blockledger/mocks"
	"github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/common/mocks/configtx"
	mockpolicies "github.com/mcc-github/blockchain/common/mocks/policies"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	"github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestChainSupportBlock(t *testing.T) {
	ledger := &mocks.ReadWriter{}
	ledger.On("Height").Return(uint64(100))
	iterator := &mock.BlockIterator{}
	iterator.NextReturns(&common.Block{Header: &common.BlockHeader{Number: 99}}, common.Status_SUCCESS)
	ledger.On("Iterator", &orderer.SeekPosition{
		Type: &orderer.SeekPosition_Specified{
			Specified: &orderer.SeekSpecified{Number: 99},
		},
	}).Return(iterator, uint64(99))
	cs := &ChainSupport{ledgerResources: &ledgerResources{ReadWriter: ledger}}

	assert.Nil(t, cs.Block(100))
	assert.Equal(t, uint64(99), cs.Block(99).Header.Number)
}

type mutableResourcesMock struct {
	config.Resources
}

func (*mutableResourcesMock) Update(*channelconfig.Bundle) {
	panic("implement me")
}

func TestVerifyBlockSignature(t *testing.T) {
	policyMgr := &mockpolicies.Manager{
		PolicyMap: make(map[string]policies.Policy),
	}
	ms := &mutableResourcesMock{
		Resources: config.Resources{
			ConfigtxValidatorVal: &configtx.Validator{ChainIDVal: "mychannel"},
			PolicyManagerVal:     policyMgr,
		},
	}
	cs := &ChainSupport{
		ledgerResources: &ledgerResources{
			configResources: &configResources{
				mutableResources: ms,
			},
		},
	}

	
	
	err := cs.VerifyBlockSignature([]*protoutil.SignedData{}, nil)
	assert.EqualError(t, err, "policy /Channel/Orderer/BlockValidation wasn't found")

	
	
	policyMgr.PolicyMap["/Channel/Orderer/BlockValidation"] = &mockpolicies.Policy{
		Err: errors.New("invalid signature"),
	}
	err = cs.VerifyBlockSignature([]*protoutil.SignedData{}, nil)
	assert.EqualError(t, err, "block verification failed: invalid signature")

	
	policyMgr.PolicyMap["/Channel/Orderer/BlockValidation"] = &mockpolicies.Policy{
		Err: nil,
	}
	assert.NoError(t, cs.VerifyBlockSignature([]*protoutil.SignedData{}, nil))

	
	err = cs.VerifyBlockSignature([]*protoutil.SignedData{}, &common.ConfigEnvelope{})
	assert.EqualError(t, err, "channelconfig Config cannot be nil")

	
	assert.NoError(t, cs.VerifyBlockSignature([]*protoutil.SignedData{}, testConfigEnvelope(t)))

}

func testConfigEnvelope(t *testing.T) *common.ConfigEnvelope {
	config := configtxgentest.Load(localconfig.SampleInsecureSoloProfile)
	group, err := encoder.NewChannelGroup(config)
	assert.NoError(t, err)
	assert.NotNil(t, group)
	return &common.ConfigEnvelope{
		Config: &common.Config{
			ChannelGroup: group,
		},
	}
}
