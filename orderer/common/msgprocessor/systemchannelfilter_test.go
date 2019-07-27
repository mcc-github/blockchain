/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	mockconfigtx "github.com/mcc-github/blockchain/common/mocks/configtx"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor/mocks"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)


type metadataValidator interface {
	MetadataValidator
}

var validConfig *cb.Config

func init() {
	cg, err := encoder.NewChannelGroup(configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile))
	if err != nil {
		panic(err)
	}
	validConfig = &cb.Config{
		Sequence:     1,
		ChannelGroup: cg,
	}
}

func mockCrypto() *mocks.SignerSerializer {
	return &mocks.SignerSerializer{}
}

func makeConfigTxFromConfigUpdateTx(configUpdateTx *cb.Envelope) *cb.Envelope {
	confUpdate := configtx.UnmarshalConfigUpdateOrPanic(
		configtx.UnmarshalConfigUpdateEnvelopeOrPanic(
			protoutil.UnmarshalPayloadOrPanic(configUpdateTx.Payload).Data,
		).ConfigUpdate,
	)
	res, err := protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG, confUpdate.ChannelId, nil, &cb.ConfigEnvelope{
		Config:     validConfig,
		LastUpdate: configUpdateTx,
	}, 0, 0)
	if err != nil {
		panic(err)
	}
	return res
}

func wrapConfigTx(env *cb.Envelope) *cb.Envelope {
	result, err := protoutil.CreateSignedEnvelope(cb.HeaderType_ORDERER_TRANSACTION, "foo", mockCrypto(), env, msgVersion, epoch)
	if err != nil {
		panic(err)
	}
	return result
}

type mockSupport struct {
	msc *mockconfig.Orderer
}

func newMockSupport() *mockSupport {
	return &mockSupport{
		msc: &mockconfig.Orderer{
			ConsensusMetadataVal: []byte("old consensus metadata"),
		},
	}
}

func (ms *mockSupport) OrdererConfig() (channelconfig.Orderer, bool) {
	return ms.msc, true
}

type mockChainCreator struct {
	ms                  *mockSupport
	newChains           []*cb.Envelope
	NewChannelConfigErr error
}

func newMockChainCreator() *mockChainCreator {
	mcc := &mockChainCreator{
		ms: newMockSupport(),
	}
	return mcc
}

func (mcc *mockChainCreator) ChannelsCount() int {
	return len(mcc.newChains)
}

func (mcc *mockChainCreator) CreateBundle(channelID string, config *cb.Config) (channelconfig.Resources, error) {
	return &mockconfig.Resources{
		ConfigtxValidatorVal: &mockconfigtx.Validator{
			ChainIDVal: channelID,
		},
		OrdererConfigVal: &mockconfig.Orderer{
			ConsensusMetadataVal: []byte("new consensus metadata"),
			CapabilitiesVal:      &mockconfig.OrdererCapabilities{},
		},
		ChannelConfigVal: &mockconfig.Channel{
			CapabilitiesVal: &mockconfig.ChannelCapabilities{},
		},
	}, nil
}

func (mcc *mockChainCreator) NewChannelConfig(envConfigUpdate *cb.Envelope) (channelconfig.Resources, error) {
	if mcc.NewChannelConfigErr != nil {
		return nil, mcc.NewChannelConfigErr
	}
	return &mockconfig.Resources{
		ConfigtxValidatorVal: &mockconfigtx.Validator{
			ProposeConfigUpdateVal: &cb.ConfigEnvelope{
				Config:     validConfig,
				LastUpdate: envConfigUpdate,
			},
		},
	}, nil
}

func TestGoodProposal(t *testing.T) {
	newChainID := "new-chain-id"

	mcc := newMockChainCreator()
	mv := &mocks.FakeMetadataValidator{}

	configUpdate, err := encoder.MakeChannelCreationTransaction(newChainID, nil, configtxgentest.Load(genesisconfig.SampleSingleMSPChannelProfile))
	assert.Nil(t, err, "Error constructing configtx")
	ingressTx := makeConfigTxFromConfigUpdateTx(configUpdate)

	wrapped := wrapConfigTx(ingressTx)

	assert.NoError(t, NewSystemChannelFilter(mcc.ms, mcc, mv).Apply(wrapped), "Did not accept valid transaction")
}

func TestProposalRejectedByConfig(t *testing.T) {
	newChainID := "NewChainID"

	mcc := newMockChainCreator()
	mcc.NewChannelConfigErr = fmt.Errorf("desired err text")
	mv := &mocks.FakeMetadataValidator{}

	configUpdate, err := encoder.MakeChannelCreationTransaction(newChainID, nil, configtxgentest.Load(genesisconfig.SampleSingleMSPChannelProfile))
	assert.Nil(t, err, "Error constructing configtx")
	ingressTx := makeConfigTxFromConfigUpdateTx(configUpdate)

	wrapped := wrapConfigTx(ingressTx)

	err = NewSystemChannelFilter(mcc.ms, mcc, mv).Apply(wrapped)

	assert.NotNil(t, err, "Did not accept valid transaction")
	assert.Regexp(t, mcc.NewChannelConfigErr.Error(), err)
	assert.Len(t, mcc.newChains, 0, "Proposal should not have created a new chain")
}

func TestNumChainsExceeded(t *testing.T) {
	newChainID := "NewChainID"

	mcc := newMockChainCreator()
	mcc.ms.msc.MaxChannelsCountVal = 1
	mcc.newChains = make([]*cb.Envelope, 2)
	mv := &mocks.FakeMetadataValidator{}

	configUpdate, err := encoder.MakeChannelCreationTransaction(newChainID, nil, configtxgentest.Load(genesisconfig.SampleSingleMSPChannelProfile))
	assert.Nil(t, err, "Error constructing configtx")
	ingressTx := makeConfigTxFromConfigUpdateTx(configUpdate)

	wrapped := wrapConfigTx(ingressTx)

	err = NewSystemChannelFilter(mcc.ms, mcc, mv).Apply(wrapped)

	assert.NotNil(t, err, "Transaction had created too many channels")
	assert.Regexp(t, "exceed maximimum number", err)
}

func TestMaintenanceMode(t *testing.T) {
	newChainID := "NewChainID"

	mcc := newMockChainCreator()
	mcc.ms.msc.ConsensusTypeStateVal = orderer.ConsensusType_STATE_MAINTENANCE
	mv := &mocks.FakeMetadataValidator{}

	configUpdate, err := encoder.MakeChannelCreationTransaction(newChainID, nil, configtxgentest.Load(genesisconfig.SampleSingleMSPChannelProfile))
	assert.Nil(t, err, "Error constructing configtx")
	ingressTx := makeConfigTxFromConfigUpdateTx(configUpdate)

	wrapped := wrapConfigTx(ingressTx)

	err = NewSystemChannelFilter(mcc.ms, mcc, mv).Apply(wrapped)
	assert.EqualError(t, err, "channel creation is not permitted: maintenance mode")
}

func TestBadProposal(t *testing.T) {
	mcc := newMockChainCreator()
	mv := &mocks.FakeMetadataValidator{}
	sysFilter := NewSystemChannelFilter(mcc.ms, mcc, mv)
	t.Run("BadPayload", func(t *testing.T) {
		err := sysFilter.Apply(&cb.Envelope{Payload: []byte("garbage payload")})
		assert.Regexp(t, "bad payload", err)
	})

	for _, tc := range []struct {
		name    string
		payload *cb.Payload
		regexp  string
	}{
		{
			"MissingPayloadHeader",
			&cb.Payload{},
			"missing payload header",
		},
		{
			"BadChannelHeader",
			&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: []byte("bad channel header"),
				},
			},
			"bad channel header",
		},
		{
			"BadConfigTx",
			&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: protoutil.MarshalOrPanic(
						&cb.ChannelHeader{
							Type: int32(cb.HeaderType_ORDERER_TRANSACTION),
						},
					),
				},
				Data: []byte("bad configTx"),
			},
			"payload data error unmarshaling to envelope",
		},
		{
			"BadConfigTxPayload",
			&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: protoutil.MarshalOrPanic(
						&cb.ChannelHeader{
							Type: int32(cb.HeaderType_ORDERER_TRANSACTION),
						},
					),
				},
				Data: protoutil.MarshalOrPanic(
					&cb.Envelope{
						Payload: []byte("bad payload"),
					},
				),
			},
			"error unmarshaling wrapped configtx envelope payload",
		},
		{
			"MissingConfigTxChannelHeader",
			&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: protoutil.MarshalOrPanic(
						&cb.ChannelHeader{
							Type: int32(cb.HeaderType_ORDERER_TRANSACTION),
						},
					),
				},
				Data: protoutil.MarshalOrPanic(
					&cb.Envelope{
						Payload: protoutil.MarshalOrPanic(
							&cb.Payload{},
						),
					},
				),
			},
			"wrapped configtx envelope missing header",
		},
		{
			"BadConfigTxChannelHeader",
			&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: protoutil.MarshalOrPanic(
						&cb.ChannelHeader{
							Type: int32(cb.HeaderType_ORDERER_TRANSACTION),
						},
					),
				},
				Data: protoutil.MarshalOrPanic(
					&cb.Envelope{
						Payload: protoutil.MarshalOrPanic(
							&cb.Payload{
								Header: &cb.Header{
									ChannelHeader: []byte("bad channel header"),
								},
							},
						),
					},
				),
			},
			"error unmarshaling wrapped configtx envelope channel header",
		},
		{
			"BadConfigTxChannelHeaderType",
			&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: protoutil.MarshalOrPanic(
						&cb.ChannelHeader{
							Type: int32(cb.HeaderType_ORDERER_TRANSACTION),
						},
					),
				},
				Data: protoutil.MarshalOrPanic(
					&cb.Envelope{
						Payload: protoutil.MarshalOrPanic(
							&cb.Payload{
								Header: &cb.Header{
									ChannelHeader: protoutil.MarshalOrPanic(
										&cb.ChannelHeader{
											Type: 0xBad,
										},
									),
								},
							},
						),
					},
				),
			},
			"wrapped configtx envelope not a config transaction",
		},
		{
			"BadConfigEnvelope",
			&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: protoutil.MarshalOrPanic(
						&cb.ChannelHeader{
							Type: int32(cb.HeaderType_ORDERER_TRANSACTION),
						},
					),
				},
				Data: protoutil.MarshalOrPanic(
					&cb.Envelope{
						Payload: protoutil.MarshalOrPanic(
							&cb.Payload{
								Header: &cb.Header{
									ChannelHeader: protoutil.MarshalOrPanic(
										&cb.ChannelHeader{
											Type: int32(cb.HeaderType_CONFIG),
										},
									),
								},
								Data: []byte("bad config update"),
							},
						),
					},
				),
			},
			"error unmarshalling wrapped configtx config envelope from payload",
		},
		{
			"MissingConfigEnvelopeLastUpdate",
			&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: protoutil.MarshalOrPanic(
						&cb.ChannelHeader{
							Type: int32(cb.HeaderType_ORDERER_TRANSACTION),
						},
					),
				},
				Data: protoutil.MarshalOrPanic(
					&cb.Envelope{
						Payload: protoutil.MarshalOrPanic(
							&cb.Payload{
								Header: &cb.Header{
									ChannelHeader: protoutil.MarshalOrPanic(
										&cb.ChannelHeader{
											Type: int32(cb.HeaderType_CONFIG),
										},
									),
								},
								Data: protoutil.MarshalOrPanic(
									&cb.ConfigEnvelope{},
								),
							},
						),
					},
				),
			},
			"updated config does not include a config update",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := sysFilter.Apply(&cb.Envelope{Payload: protoutil.MarshalOrPanic(tc.payload)})
			assert.NotNil(t, err)
			assert.Regexp(t, tc.regexp, err.Error())
		})
	}
}

func TestFailedMetadataValidation(t *testing.T) {
	newChainID := "NewChainID"

	errorString := "bananas"
	mv := &mocks.FakeMetadataValidator{}
	mv.ValidateConsensusMetadataReturns(errors.New(errorString))
	mcc := newMockChainCreator()

	configUpdate, err := encoder.MakeChannelCreationTransaction(newChainID, nil, configtxgentest.Load(genesisconfig.SampleSingleMSPChannelProfile))
	assert.NoError(t, err, "Error constructing configtx")
	ingressTx := makeConfigTxFromConfigUpdateTx(configUpdate)

	wrapped := wrapConfigTx(ingressTx)

	err = NewSystemChannelFilter(mcc.ms, mcc, mv).Apply(wrapped)

	
	assert.EqualError(t, err, "consensus metadata update for channel creation is invalid: bananas", "Transaction metadata validation fails")

	
	assert.Equal(t, 1, mv.ValidateConsensusMetadataCallCount())
	om, nm, nc := mv.ValidateConsensusMetadataArgsForCall(0)
	assert.True(t, nc)
	assert.Equal(t, []byte("old consensus metadata"), om)
	assert.Equal(t, []byte("new consensus metadata"), nm)
}
