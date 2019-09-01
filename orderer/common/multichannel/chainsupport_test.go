/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/deliver/mock"
	mockblockledger "github.com/mcc-github/blockchain/common/ledger/blockledger/mocks"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	"github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	msgprocessormocks "github.com/mcc-github/blockchain/orderer/common/msgprocessor/mocks"
	"github.com/mcc-github/blockchain/orderer/common/multichannel/mocks"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)



type policyManager interface {
	policies.Manager
}



type policy interface {
	policies.Policy
}

func TestChainSupportBlock(t *testing.T) {
	ledger := &mockblockledger.ReadWriter{}
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
	*mocks.Resources
	newConsensusMetadataVal []byte
}

func (*mutableResourcesMock) Update(*channelconfig.Bundle) {
	panic("implement me")
}

func (mrm *mutableResourcesMock) CreateBundle(channelID string, c *common.Config) (channelconfig.Resources, error) {
	mockOrderer := &mocks.OrdererConfig{}
	mockOrderer.ConsensusMetadataReturns(mrm.newConsensusMetadataVal)
	mockResources := &mocks.Resources{}
	mockResources.OrdererConfigReturns(mockOrderer, true)

	return mockResources, nil

}

func TestVerifyBlockSignature(t *testing.T) {
	mockResources := &mocks.Resources{}
	mockValidator := &mocks.ConfigTXValidator{}
	mockValidator.ChannelIDReturns("mychannel")
	mockResources.ConfigtxValidatorReturns(mockValidator)

	mockPolicy := &mocks.Policy{}
	mockPolicyManager := &mocks.PolicyManager{}
	mockResources.PolicyManagerReturns(mockPolicyManager)

	ms := &mutableResourcesMock{
		Resources: mockResources,
	}
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	cs := &ChainSupport{
		ledgerResources: &ledgerResources{
			configResources: &configResources{
				mutableResources: ms,
				bccsp:            cryptoProvider,
			},
		},
	}

	
	
	mockPolicyManager.GetPolicyReturns(nil, false)
	err = cs.VerifyBlockSignature([]*protoutil.SignedData{}, nil)
	assert.EqualError(t, err, "policy /Channel/Orderer/BlockValidation wasn't found")

	mockPolicyManager.GetPolicyReturns(mockPolicy, true)
	
	
	mockPolicy.EvaluateReturns(errors.New("invalid signature"))
	err = cs.VerifyBlockSignature([]*protoutil.SignedData{}, nil)
	assert.EqualError(t, err, "block verification failed: invalid signature")

	
	mockPolicy.EvaluateReturns(nil)
	assert.NoError(t, cs.VerifyBlockSignature([]*protoutil.SignedData{}, nil))

	
	err = cs.VerifyBlockSignature([]*protoutil.SignedData{}, &common.ConfigEnvelope{})
	assert.EqualError(t, err, "channelconfig Config cannot be nil")

	
	assert.NoError(t, cs.VerifyBlockSignature([]*protoutil.SignedData{}, testConfigEnvelope(t)))

}

func TestConsensusMetadataValidation(t *testing.T) {
	oldConsensusMetadata := []byte("old consensus metadata")
	newConsensusMetadata := []byte("new consensus metadata")
	mockValidator := &mocks.ConfigTXValidator{}
	mockValidator.ChannelIDReturns("mychannel")
	mockValidator.ProposeConfigUpdateReturns(testConfigEnvelope(t), nil)
	mockOrderer := &mocks.OrdererConfig{}
	mockOrderer.ConsensusMetadataReturns(oldConsensusMetadata)
	mockResources := &mocks.Resources{}
	mockResources.ConfigtxValidatorReturns(mockValidator)
	mockResources.OrdererConfigReturns(mockOrderer, true)

	ms := &mutableResourcesMock{
		Resources:               mockResources,
		newConsensusMetadataVal: newConsensusMetadata,
	}
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	mv := &msgprocessormocks.MetadataValidator{}
	cs := &ChainSupport{
		ledgerResources: &ledgerResources{
			configResources: &configResources{
				mutableResources: ms,
				bccsp:            cryptoProvider,
			},
		},
		MetadataValidator: mv,
	}

	
	_, err = cs.ProposeConfigUpdate(&common.Envelope{})
	assert.NoError(t, err)

	
	assert.Equal(t, 1, mv.ValidateConsensusMetadataCallCount())
	om, nm, nc := mv.ValidateConsensusMetadataArgsForCall(0)
	assert.False(t, nc)
	assert.Equal(t, oldConsensusMetadata, om)
	assert.Equal(t, newConsensusMetadata, nm)

	
	mv.ValidateConsensusMetadataReturns(errors.New("bananas"))
	_, err = cs.ProposeConfigUpdate(&common.Envelope{})
	assert.EqualError(t, err, "consensus metadata update for channel config update is invalid: bananas")
}

func testConfigEnvelope(t *testing.T) *common.ConfigEnvelope {
	conf := configtxgentest.Load(localconfig.SampleInsecureSoloProfile)
	group, err := encoder.NewChannelGroup(conf)
	assert.NoError(t, err)
	group.Groups["Orderer"].Values["ConsensusType"].Value, err = proto.Marshal(&orderer.ConsensusType{
		Metadata: []byte("new consensus metadata"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, group)
	return &common.ConfigEnvelope{
		Config: &common.Config{
			ChannelGroup: group,
		},
	}
}
