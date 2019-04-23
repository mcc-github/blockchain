/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/common/ledger/blockledger/ram"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	mockchannelconfig "github.com/mcc-github/blockchain/common/mocks/config"
	mockpolicies "github.com/mcc-github/blockchain/common/mocks/policies"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/orderer/common/blockcutter"
	"github.com/mcc-github/blockchain/orderer/common/multichannel/mocks"
	"github.com/mcc-github/blockchain/orderer/consensus"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)



type signerSerializer interface {
	identity.SignerSerializer
}

func mockCrypto() *mocks.SignerSerializer {
	return &mocks.SignerSerializer{}
}

func newRAMLedgerAndFactory(maxSize int,
	chainID string, genesisBlockSys *cb.Block) (blockledger.Factory, blockledger.ReadWriter) {
	rlf := ramledger.New(maxSize)
	rl, err := rlf.GetOrCreate(chainID)
	if err != nil {
		panic(err)
	}
	err = rl.Append(genesisBlockSys)
	if err != nil {
		panic(err)
	}
	return rlf, rl
}

func newRAMLedgerAndFactory3Chan(maxSize int,
	chainID string, genesisBlockSys *cb.Block,
	chainID1 string, genesisBlockStd1 *cb.Block,
	chainID2 string, genesisBlockStd2 *cb.Block) (blockledger.Factory, []blockledger.ReadWriter) {
	var rls []blockledger.ReadWriter
	rlf, rl := newRAMLedgerAndFactory(maxSize, chainID, genesisBlockSys)
	rls = append(rls, rl)

	rl, err := rlf.GetOrCreate(chainID1)
	if err != nil {
		panic(err)
	}
	err = rl.Append(genesisBlockStd1)
	if err != nil {
		panic(err)
	}
	rls = append(rls, rl)

	rl, err = rlf.GetOrCreate(chainID2)
	if err != nil {
		panic(err)
	}
	err = rl.Append(genesisBlockStd2)
	if err != nil {
		panic(err)
	}
	rls = append(rls, rl)

	return rlf, rls
}

func testMessageOrderAndRetrieval(maxMessageCount uint32, chainID string, chainSupport *ChainSupport, lr blockledger.ReadWriter, t *testing.T) {
	messages := make([]*cb.Envelope, maxMessageCount)
	for i := uint32(0); i < maxMessageCount; i++ {
		messages[i] = makeNormalTx(chainID, int(i))
	}
	for _, message := range messages {
		chainSupport.Order(message, 0)
	}
	it, _ := lr.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: 1}}})
	defer it.Close()
	block, status := it.Next()
	assert.Equal(t, cb.Status_SUCCESS, status, "Could not retrieve block")
	for i := uint32(0); i < maxMessageCount; i++ {
		assert.True(t, proto.Equal(messages[i], protoutil.ExtractEnvelopeOrPanic(block, int(i))), "Block contents wrong at index %d", i)
	}
}

func TestConfigTx(t *testing.T) {
	
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()

	
	
	t.Run("GetConfigTx - ok", func(t *testing.T) {
		_, rl := newRAMLedgerAndFactory(10, genesisconfig.TestChainID, genesisBlockSys)
		for i := 0; i < 5; i++ {
			rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{makeNormalTx(genesisconfig.TestChainID, i)}))
		}
		rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{makeConfigTx(genesisconfig.TestChainID, 5)}))
		ctx := makeConfigTx(genesisconfig.TestChainID, 6)
		rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{ctx}))

		block := blockledger.CreateNextBlock(rl, []*cb.Envelope{makeNormalTx(genesisconfig.TestChainID, 7)})
		block.Metadata.Metadata[cb.BlockMetadataIndex_LAST_CONFIG] = protoutil.MarshalOrPanic(&cb.Metadata{Value: protoutil.MarshalOrPanic(&cb.LastConfig{Index: 7})})
		rl.Append(block)

		pctx := configTx(rl)
		assert.True(t, proto.Equal(pctx, ctx), "Did not select most recent config transaction")
	})

	
	
	t.Run("GetConfigTx - failure", func(t *testing.T) {
		_, rl := newRAMLedgerAndFactory(10, genesisconfig.TestChainID, genesisBlockSys)
		for i := 0; i < 10; i++ {
			rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{
				makeNormalTx(genesisconfig.TestChainID, i),
				makeConfigTx(genesisconfig.TestChainID, i),
			}))
		}
		rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{makeNormalTx(genesisconfig.TestChainID, 11)}))
		assert.Panics(t, func() { configTx(rl) }, "Should have panicked because there was no config tx")

		block := blockledger.CreateNextBlock(rl, []*cb.Envelope{makeNormalTx(genesisconfig.TestChainID, 12)})
		block.Metadata.Metadata[cb.BlockMetadataIndex_LAST_CONFIG] = []byte("bad metadata")
		assert.Panics(t, func() { configTx(rl) }, "Should have panicked because of bad last config metadata")
	})
}

func TestNewRegistrar(t *testing.T) {
	
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()

	
	t.Run("No system chain - failure", func(t *testing.T) {
		lf := ramledger.New(10)

		consenters := make(map[string]consensus.Consenter)
		consenters[confSys.Orderer.OrdererType] = &mockConsenter{}

		assert.Panics(t, func() {
			NewRegistrar(lf, mockCrypto(), &disabled.Provider{}).Initialize(consenters)
		}, "Should have panicked when starting without a system chain")
	})

	
	t.Run("Multiple system chains - failure", func(t *testing.T) {
		lf := ramledger.New(10)

		for _, id := range []string{"foo", "bar"} {
			rl, err := lf.GetOrCreate(id)
			assert.NoError(t, err)

			err = rl.Append(encoder.New(confSys).GenesisBlockForChannel(id))
			assert.NoError(t, err)
		}

		consenters := make(map[string]consensus.Consenter)
		consenters[confSys.Orderer.OrdererType] = &mockConsenter{}

		assert.Panics(t, func() {
			NewRegistrar(lf, mockCrypto(), &disabled.Provider{}).Initialize(consenters)
		}, "Two system channels should have caused panic")
	})

	
	t.Run("Correct flow", func(t *testing.T) {
		lf, rl := newRAMLedgerAndFactory(10, genesisconfig.TestChainID, genesisBlockSys)

		consenters := make(map[string]consensus.Consenter)
		consenters[confSys.Orderer.OrdererType] = &mockConsenter{}

		manager := NewRegistrar(lf, mockCrypto(), &disabled.Provider{})
		manager.Initialize(consenters)

		chainSupport := manager.GetChain("Fake")
		assert.Nilf(t, chainSupport, "Should not have found a chain that was not created")

		chainSupport = manager.GetChain(genesisconfig.TestChainID)
		assert.NotNilf(t, chainSupport, "Should have gotten chain which was initialized by ramledger")

		testMessageOrderAndRetrieval(confSys.Orderer.BatchSize.MaxMessageCount, genesisconfig.TestChainID, chainSupport, rl, t)
	})
}

func TestCreateChain(t *testing.T) {
	
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()

	t.Run("Create chain", func(t *testing.T) {
		lf, _ := newRAMLedgerAndFactory(10, genesisconfig.TestChainID, genesisBlockSys)

		consenters := make(map[string]consensus.Consenter)
		consenters[confSys.Orderer.OrdererType] = &mockConsenter{}

		manager := NewRegistrar(lf, mockCrypto(), &disabled.Provider{})
		manager.Initialize(consenters)

		ledger, err := lf.GetOrCreate("mychannel")
		assert.NoError(t, err)

		genesisBlock := encoder.New(confSys).GenesisBlockForChannel("mychannel")
		ledger.Append(genesisBlock)

		
		assert.Nil(t, manager.GetChain("mychannel"))
		
		manager.CreateChain("mychannel")
		chain := manager.GetChain("mychannel")
		assert.NotNil(t, chain)
		
		manager.CreateChain("mychannel")
		chain2 := manager.GetChain("mychannel")
		assert.NotNil(t, chain2)
		
		assert.NotEqual(t, chain, chain2)
		
		_, ok := <-chain.Chain.(*mockChain).queue
		assert.False(t, ok)
		
		close(chain2.Chain.(*mockChain).queue)
	})

	
	t.Run("New chain", func(t *testing.T) {
		expectedLastConfigBlockNumber := uint64(0)
		expectedLastConfigSeq := uint64(1)
		newChainID := "test-new-chain"

		lf, rl := newRAMLedgerAndFactory(10, genesisconfig.TestChainID, genesisBlockSys)

		consenters := make(map[string]consensus.Consenter)
		consenters[confSys.Orderer.OrdererType] = &mockConsenter{}

		manager := NewRegistrar(lf, mockCrypto(), &disabled.Provider{})
		manager.Initialize(consenters)
		orglessChannelConf := configtxgentest.Load(genesisconfig.SampleSingleMSPChannelProfile)
		orglessChannelConf.Application.Organizations = nil
		envConfigUpdate, err := encoder.MakeChannelCreationTransaction(newChainID, mockCrypto(), orglessChannelConf)
		assert.NoError(t, err, "Constructing chain creation tx")

		res, err := manager.NewChannelConfig(envConfigUpdate)
		assert.NoError(t, err, "Constructing initial channel config")

		configEnv, err := res.ConfigtxValidator().ProposeConfigUpdate(envConfigUpdate)
		assert.NoError(t, err, "Proposing initial update")
		assert.Equal(t, expectedLastConfigSeq, configEnv.GetConfig().Sequence, "Sequence of config envelope for new channel should always be set to %d", expectedLastConfigSeq)

		ingressTx, err := protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG, newChainID, mockCrypto(), configEnv, msgVersion, epoch)
		assert.NoError(t, err, "Creating ingresstx")

		wrapped := wrapConfigTx(ingressTx)

		chainSupport := manager.GetChain(manager.SystemChannelID())
		assert.NotNilf(t, chainSupport, "Could not find system channel")

		chainSupport.Configure(wrapped, 0)
		func() {
			it, _ := rl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: 1}}})
			defer it.Close()
			block, status := it.Next()
			if status != cb.Status_SUCCESS {
				t.Fatalf("Could not retrieve block")
			}
			if len(block.Data.Data) != 1 {
				t.Fatalf("Should have had only one message in the orderer transaction block")
			}

			assert.True(t, proto.Equal(wrapped, protoutil.UnmarshalEnvelopeOrPanic(block.Data.Data[0])), "Orderer config block contains wrong transaction")
		}()

		chainSupport = manager.GetChain(newChainID)
		if chainSupport == nil {
			t.Fatalf("Should have gotten new chain which was created")
		}

		messages := make([]*cb.Envelope, confSys.Orderer.BatchSize.MaxMessageCount)
		for i := 0; i < int(confSys.Orderer.BatchSize.MaxMessageCount); i++ {
			messages[i] = makeNormalTx(newChainID, i)
		}

		for _, message := range messages {
			chainSupport.Order(message, 0)
		}

		it, _ := chainSupport.Reader().Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: 0}}})
		defer it.Close()
		block, status := it.Next()
		if status != cb.Status_SUCCESS {
			t.Fatalf("Could not retrieve new chain genesis block")
		}
		testLastConfigBlockNumber(t, block, expectedLastConfigBlockNumber)
		if len(block.Data.Data) != 1 {
			t.Fatalf("Should have had only one message in the new genesis block")
		}

		assert.True(t, proto.Equal(ingressTx, protoutil.UnmarshalEnvelopeOrPanic(block.Data.Data[0])), "Genesis block contains wrong transaction")

		block, status = it.Next()
		if status != cb.Status_SUCCESS {
			t.Fatalf("Could not retrieve block on new chain")
		}
		testLastConfigBlockNumber(t, block, expectedLastConfigBlockNumber)
		for i := 0; i < int(confSys.Orderer.BatchSize.MaxMessageCount); i++ {
			if !proto.Equal(protoutil.ExtractEnvelopeOrPanic(block, i), messages[i]) {
				t.Errorf("Block contents wrong at index %d in new chain", i)
			}
		}

		rcs := newChainSupport(manager, chainSupport.ledgerResources, consenters, mockCrypto(), blockcutter.NewMetrics(&disabled.Provider{}))
		assert.Equal(t, expectedLastConfigSeq, rcs.lastConfigSeq, "On restart, incorrect lastConfigSeq")
	})
}

func testLastConfigBlockNumber(t *testing.T, block *cb.Block, expectedBlockNumber uint64) {
	metadataItem := &cb.Metadata{}
	err := proto.Unmarshal(block.Metadata.Metadata[cb.BlockMetadataIndex_LAST_CONFIG], metadataItem)
	assert.NoError(t, err, "Block should carry LAST_CONFIG metadata item")
	lastConfig := &cb.LastConfig{}
	err = proto.Unmarshal(metadataItem.Value, lastConfig)
	assert.NoError(t, err, "LAST_CONFIG metadata item should carry last config value")
	assert.Equal(t, expectedBlockNumber, lastConfig.Index, "LAST_CONFIG value should point to last config block")
}

func TestResourcesCheck(t *testing.T) {
	t.Run("GoodResources", func(t *testing.T) {
		err := checkResources(&mockchannelconfig.Resources{
			PolicyManagerVal: &mockpolicies.Manager{},
			OrdererConfigVal: &mockchannelconfig.Orderer{
				CapabilitiesVal: &mockchannelconfig.OrdererCapabilities{},
			},
			ChannelConfigVal: &mockchannelconfig.Channel{
				CapabilitiesVal: &mockchannelconfig.ChannelCapabilities{},
			},
		})

		assert.NoError(t, err)
	})

	t.Run("MissingOrdererConfigPanic", func(t *testing.T) {
		err := checkResources(&mockchannelconfig.Resources{
			PolicyManagerVal: &mockpolicies.Manager{},
		})

		assert.Error(t, err)
		assert.Regexp(t, "config does not contain orderer config", err.Error())
	})

	t.Run("MissingOrdererCapability", func(t *testing.T) {
		err := checkResources(&mockchannelconfig.Resources{
			PolicyManagerVal: &mockpolicies.Manager{},
			OrdererConfigVal: &mockchannelconfig.Orderer{
				CapabilitiesVal: &mockchannelconfig.OrdererCapabilities{
					SupportedErr: errors.New("An error"),
				},
			},
		})

		assert.Error(t, err)
		assert.Regexp(t, "config requires unsupported orderer capabilities:", err.Error())
	})

	t.Run("MissingChannelCapability", func(t *testing.T) {
		err := checkResources(&mockchannelconfig.Resources{
			PolicyManagerVal: &mockpolicies.Manager{},
			OrdererConfigVal: &mockchannelconfig.Orderer{
				CapabilitiesVal: &mockchannelconfig.OrdererCapabilities{},
			},
			ChannelConfigVal: &mockchannelconfig.Channel{
				CapabilitiesVal: &mockchannelconfig.ChannelCapabilities{
					SupportedErr: errors.New("An error"),
				},
			},
		})

		assert.Error(t, err)
		assert.Regexp(t, "config requires unsupported channel capabilities:", err.Error())
	})

	t.Run("MissingOrdererConfigPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			checkResourcesOrPanic(&mockchannelconfig.Resources{
				PolicyManagerVal: &mockpolicies.Manager{},
			})
		})
	})
}


func TestBroadcastChannelSupportRejection(t *testing.T) {
	
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()

	t.Run("Rejection", func(t *testing.T) {

		ledgerFactory, _ := newRAMLedgerAndFactory(10, genesisconfig.TestChainID, genesisBlockSys)
		mockConsenters := map[string]consensus.Consenter{confSys.Orderer.OrdererType: &mockConsenter{}}
		registrar := NewRegistrar(ledgerFactory, mockCrypto(), &disabled.Provider{})
		registrar.Initialize(mockConsenters)
		randomValue := 1
		configTx := makeConfigTx(genesisconfig.TestChainID, randomValue)
		_, _, _, err := registrar.BroadcastChannelSupport(configTx)
		assert.Error(t, err, "Messages of type HeaderType_CONFIG should return an error.")
	})
}
