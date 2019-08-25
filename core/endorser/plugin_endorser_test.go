/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/mocks/ledger"
	"github.com/mcc-github/blockchain/core/endorser"
	"github.com/mcc-github/blockchain/core/endorser/mocks"
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	. "github.com/mcc-github/blockchain/core/handlers/endorsement/api/state"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/peer"
	transientstore2 "github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockTransientStoreRetriever = transientStoreRetriever()
	mockTransientStore          = &mocks.Store{}
)

func TestPluginEndorserNotFound(t *testing.T) {
	pluginMapper := &mocks.PluginMapper{}
	pluginMapper.On("PluginFactoryByName", endorser.PluginName("notfound")).Return(nil)
	pluginEndorser := endorser.NewPluginEndorser(&endorser.PluginSupport{
		PluginMapper: pluginMapper,
	})
	endorsement, prpBytes, err := pluginEndorser.EndorseWithPlugin("notfound", "", nil, nil)
	assert.Nil(t, endorsement)
	assert.Nil(t, prpBytes)
	assert.Contains(t, err.Error(), "plugin with name notfound wasn't found")
}

func TestPluginEndorserGreenPath(t *testing.T) {
	expectedSignature := []byte{5, 4, 3, 2, 1}
	expectedProposalResponsePayload := []byte{1, 2, 3}
	pluginMapper := &mocks.PluginMapper{}
	pluginFactory := &mocks.PluginFactory{}
	plugin := &mocks.Plugin{}
	plugin.On("Endorse", mock.Anything, mock.Anything).Return(&peer.Endorsement{Signature: expectedSignature}, expectedProposalResponsePayload, nil)
	pluginMapper.On("PluginFactoryByName", endorser.PluginName("plugin")).Return(pluginFactory)
	
	plugin.On("Init", mock.Anything, mock.Anything).Return(nil).Once()
	pluginFactory.On("New").Return(plugin).Once()
	sif := &mocks.SigningIdentityFetcher{}
	cs := &mocks.ChannelStateRetriever{}
	queryCreator := &mocks.QueryCreator{}
	cs.On("NewQueryCreator", "mychannel").Return(queryCreator, nil)
	pluginEndorser := endorser.NewPluginEndorser(&endorser.PluginSupport{
		ChannelStateRetriever:   cs,
		SigningIdentityFetcher:  sif,
		PluginMapper:            pluginMapper,
		TransientStoreRetriever: mockTransientStoreRetriever,
	})

	
	endorsement, prpBytes, err := pluginEndorser.EndorseWithPlugin("plugin", "mychannel", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, endorsement.Signature)
	assert.Equal(t, expectedProposalResponsePayload, prpBytes)
	
	plugin.AssertCalled(t, "Init", &endorser.ChannelState{QueryCreator: queryCreator, Store: mockTransientStore}, sif)

	
	
	
	
	endorsement, prpBytes, err = pluginEndorser.EndorseWithPlugin("plugin", "mychannel", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, endorsement.Signature)
	assert.Equal(t, expectedProposalResponsePayload, prpBytes)
	pluginFactory.AssertNumberOfCalls(t, "New", 1)
	plugin.AssertNumberOfCalls(t, "Init", 1)

	
	
	
	pluginFactory.On("New").Return(plugin).Once()
	plugin.On("Init", mock.Anything).Return(nil).Once()
	endorsement, prpBytes, err = pluginEndorser.EndorseWithPlugin("plugin", "", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, endorsement.Signature)
	assert.Equal(t, expectedProposalResponsePayload, prpBytes)
	plugin.AssertCalled(t, "Init", sif)
}

func TestPluginEndorserErrors(t *testing.T) {
	pluginMapper := &mocks.PluginMapper{}
	pluginFactory := &mocks.PluginFactory{}
	plugin := &mocks.Plugin{}
	plugin.On("Endorse", mock.Anything, mock.Anything)
	pluginMapper.On("PluginFactoryByName", endorser.PluginName("plugin")).Return(pluginFactory)
	pluginFactory.On("New").Return(plugin)
	sif := &mocks.SigningIdentityFetcher{}
	cs := &mocks.ChannelStateRetriever{}
	queryCreator := &mocks.QueryCreator{}
	cs.On("NewQueryCreator", "mychannel").Return(queryCreator, nil)
	pluginEndorser := endorser.NewPluginEndorser(&endorser.PluginSupport{
		ChannelStateRetriever:   cs,
		SigningIdentityFetcher:  sif,
		PluginMapper:            pluginMapper,
		TransientStoreRetriever: mockTransientStoreRetriever,
	})

	
	t.Run("PluginInitializationFailure", func(t *testing.T) {
		plugin.On("Init", mock.Anything, mock.Anything).Return(errors.New("plugin initialization failed")).Once()
		endorsement, prpBytes, err := pluginEndorser.EndorseWithPlugin("plugin", "mychannel", nil, nil)
		assert.Nil(t, endorsement)
		assert.Nil(t, prpBytes)
		assert.Contains(t, err.Error(), "plugin initialization failed")
	})

}

func transientStoreRetriever() *mocks.TransientStoreRetriever {
	storeRetriever := &mocks.TransientStoreRetriever{}
	storeRetriever.On("StoreForChannel", mock.Anything).Return(mockTransientStore)
	return storeRetriever
}

type fakeEndorsementPlugin struct {
	StateFetcher
}

func (fep *fakeEndorsementPlugin) Endorse(payload []byte, sp *peer.SignedProposal) (*peer.Endorsement, []byte, error) {
	state, _ := fep.StateFetcher.FetchState()
	txrws, _ := state.GetTransientByTXID("tx")
	b, _ := proto.Marshal(txrws[0])
	return nil, b, nil
}

func (fep *fakeEndorsementPlugin) Init(dependencies ...endorsement.Dependency) error {
	for _, dep := range dependencies {
		if state, isState := dep.(StateFetcher); isState {
			fep.StateFetcher = state
			return nil
		}
	}
	panic("could not find State dependency")
}

type rwsetScanner struct {
	mock.Mock
	data []*rwset.TxPvtReadWriteSet
}

func (*rwsetScanner) Next() (*transientstore.EndorserPvtSimulationResults, error) {
	panic("implement me")
}

func (rws *rwsetScanner) NextWithConfig() (*transientstore.EndorserPvtSimulationResultsWithConfig, error) {
	if len(rws.data) == 0 {
		return nil, nil
	}
	res := rws.data[0]
	rws.data = rws.data[1:]
	return &transientstore.EndorserPvtSimulationResultsWithConfig{
		PvtSimulationResultsWithConfig: &transientstore2.TxPvtReadWriteSetWithConfigInfo{
			PvtRwset: res,
		},
	}, nil
}

func (rws *rwsetScanner) Close() {
	rws.Called()
}

func TestTransientStore(t *testing.T) {
	plugin := &fakeEndorsementPlugin{}
	factory := &mocks.PluginFactory{}
	factory.On("New").Return(plugin)
	sif := &mocks.SigningIdentityFetcher{}
	cs := &mocks.ChannelStateRetriever{}
	queryCreator := &mocks.QueryCreator{}
	queryCreator.On("NewQueryExecutor").Return(&ledger.MockQueryExecutor{}, nil)
	cs.On("NewQueryCreator", "mychannel").Return(queryCreator, nil)

	transientStore := &mocks.Store{}
	storeRetriever := &mocks.TransientStoreRetriever{}
	storeRetriever.On("StoreForChannel", mock.Anything).Return(transientStore)

	pluginEndorser := endorser.NewPluginEndorser(&endorser.PluginSupport{
		ChannelStateRetriever:  cs,
		SigningIdentityFetcher: sif,
		PluginMapper: endorser.MapBasedPluginMapper{
			"plugin": factory,
		},
		TransientStoreRetriever: storeRetriever,
	})

	rws := &rwset.TxPvtReadWriteSet{
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{
			{
				Namespace: "ns",
				CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
					{
						CollectionName: "col",
					},
				},
			},
		},
	}
	scanner := &rwsetScanner{
		data: []*rwset.TxPvtReadWriteSet{rws},
	}
	scanner.On("Close")

	transientStore.On("GetTxPvtRWSetByTxid", mock.Anything, mock.Anything).Return(scanner, nil)

	_, prpBytes, err := pluginEndorser.EndorseWithPlugin("plugin", "mychannel", nil, nil)
	assert.NoError(t, err)

	txrws := &rwset.TxPvtReadWriteSet{}
	err = proto.Unmarshal(prpBytes, txrws)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(rws, txrws))
	scanner.AssertCalled(t, "Close")
}
