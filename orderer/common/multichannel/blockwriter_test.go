/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"testing"

	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/common/channelconfig"
	newchannelconfig "github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/common/ledger/blockledger/ramledger"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/orderer/common/blockcutter/mock"
	"github.com/mcc-github/blockchain/orderer/common/multichannel/mocks"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



type configtxValidator interface {
	configtx.Validator
}

type mockBlockWriterSupport struct {
	*mocks.ConfigTXValidator
	identity.SignerSerializer
	blockledger.ReadWriter
	fakeConfig *mock.OrdererConfig
	bccsp      bccsp.BCCSP
	sequence   uint64
}

func (mbws mockBlockWriterSupport) Update(bundle *newchannelconfig.Bundle) {
	return
}

func (mbws mockBlockWriterSupport) CreateBundle(channelID string, config *cb.Config) (*newchannelconfig.Bundle, error) {
	return channelconfig.NewBundle(channelID, config, mbws.bccsp)
}

func (mbws mockBlockWriterSupport) SharedConfig() newchannelconfig.Orderer {
	return mbws.fakeConfig
}

func TestCreateBlock(t *testing.T) {
	seedBlock := protoutil.NewBlock(7, []byte("lasthash"))
	seedBlock.Data.Data = [][]byte{[]byte("somebytes")}

	bw := &BlockWriter{lastBlock: seedBlock}
	block := bw.CreateNextBlock([]*cb.Envelope{
		{Payload: []byte("some other bytes")},
	})

	assert.Equal(t, seedBlock.Header.Number+1, block.Header.Number)
	assert.Equal(t, protoutil.BlockDataHash(block.Data), block.Header.DataHash)
	assert.Equal(t, protoutil.BlockHeaderHash(seedBlock.Header), block.Header.PreviousHash)
}

func TestBlockSignature(t *testing.T) {
	rlf := ramledger.New(2)
	l, err := rlf.GetOrCreate("mychannel")
	assert.NoError(t, err)
	lastBlock := protoutil.NewBlock(0, nil)
	l.Append(lastBlock)

	bw := &BlockWriter{
		lastConfigBlockNum: 42,
		support: &mockBlockWriterSupport{
			SignerSerializer:  mockCrypto(),
			ConfigTXValidator: &mocks.ConfigTXValidator{},
			ReadWriter:        l,
		},
		lastBlock: protoutil.NewBlock(1, protoutil.BlockHeaderHash(lastBlock.Header)),
	}

	consensusMetadata := []byte("bar")
	bw.commitBlock(consensusMetadata)

	it, seq := l.Iterator(&orderer.SeekPosition{Type: &orderer.SeekPosition_Newest{}})
	assert.Equal(t, uint64(1), seq)
	committedBlock, status := it.Next()
	assert.Equal(t, cb.Status_SUCCESS, status)

	md := protoutil.GetMetadataFromBlockOrPanic(committedBlock, cb.BlockMetadataIndex_SIGNATURES)

	expectedMetadataValue := protoutil.MarshalOrPanic(&cb.OrdererBlockMetadata{
		LastConfig:        &cb.LastConfig{Index: 42},
		ConsenterMetadata: protoutil.MarshalOrPanic(&cb.Metadata{Value: consensusMetadata}),
	})

	assert.Equal(t, expectedMetadataValue, md.Value, "Value contains the consensus metadata and the last config")
	assert.NotNil(t, md.Signatures, "Should have signature")
}

func TestBlockLastConfig(t *testing.T) {
	lastConfigSeq := uint64(6)
	newConfigSeq := lastConfigSeq + 1
	newBlockNum := uint64(9)

	mockValidator := &mocks.ConfigTXValidator{}
	mockValidator.SequenceReturns(newConfigSeq)
	bw := &BlockWriter{
		support: &mockBlockWriterSupport{
			SignerSerializer:  mockCrypto(),
			ConfigTXValidator: mockValidator,
		},
		lastConfigSeq: lastConfigSeq,
	}

	block := protoutil.NewBlock(newBlockNum, []byte("foo"))
	bw.addLastConfig(block)

	assert.Equal(t, newBlockNum, bw.lastConfigBlockNum)
	assert.Equal(t, newConfigSeq, bw.lastConfigSeq)

	md := protoutil.GetMetadataFromBlockOrPanic(block, cb.BlockMetadataIndex_LAST_CONFIG)
	assert.NotNil(t, md.Value, "Value not be empty in this case")
	assert.Nil(t, md.Signatures, "Should not have signature")

	lc := protoutil.GetLastConfigIndexFromBlockOrPanic(block)
	assert.Equal(t, newBlockNum, lc)
}

func TestWriteConfigBlock(t *testing.T) {
	
	t.Run("EmptyBlock", func(t *testing.T) {
		assert.Panics(t, func() { (&BlockWriter{}).WriteConfigBlock(&cb.Block{}, nil) })
	})
	t.Run("BadPayload", func(t *testing.T) {
		assert.Panics(t, func() {
			(&BlockWriter{}).WriteConfigBlock(&cb.Block{
				Data: &cb.BlockData{
					Data: [][]byte{
						protoutil.MarshalOrPanic(&cb.Envelope{Payload: []byte("bad")}),
					},
				},
			}, nil)
		})
	})
	t.Run("MissingHeader", func(t *testing.T) {
		assert.Panics(t, func() {
			(&BlockWriter{}).WriteConfigBlock(&cb.Block{
				Data: &cb.BlockData{
					Data: [][]byte{
						protoutil.MarshalOrPanic(&cb.Envelope{
							Payload: protoutil.MarshalOrPanic(&cb.Payload{}),
						}),
					},
				},
			}, nil)
		})
	})
	t.Run("BadChannelHeader", func(t *testing.T) {
		assert.Panics(t, func() {
			(&BlockWriter{}).WriteConfigBlock(&cb.Block{
				Data: &cb.BlockData{
					Data: [][]byte{
						protoutil.MarshalOrPanic(&cb.Envelope{
							Payload: protoutil.MarshalOrPanic(&cb.Payload{
								Header: &cb.Header{
									ChannelHeader: []byte("bad"),
								},
							}),
						}),
					},
				},
			}, nil)
		})
	})
	t.Run("BadChannelHeaderType", func(t *testing.T) {
		assert.Panics(t, func() {
			(&BlockWriter{}).WriteConfigBlock(&cb.Block{
				Data: &cb.BlockData{
					Data: [][]byte{
						protoutil.MarshalOrPanic(&cb.Envelope{
							Payload: protoutil.MarshalOrPanic(&cb.Payload{
								Header: &cb.Header{
									ChannelHeader: protoutil.MarshalOrPanic(&cb.ChannelHeader{}),
								},
							}),
						}),
					},
				},
			}, nil)
		})
	})
}

func TestGoodWriteConfig(t *testing.T) {
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()
	_, l := newRAMLedgerAndFactory(10, genesisconfig.TestChannelID, genesisBlockSys)

	fakeConfig := &mock.OrdererConfig{}
	fakeConfig.ConsensusTypeReturns("solo")

	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)

	mockValidator := &mocks.ConfigTXValidator{}
	mockValidator.ChannelIDReturns(genesisconfig.TestChannelID)
	bw := newBlockWriter(genesisBlockSys, nil,
		&mockBlockWriterSupport{
			SignerSerializer:  mockCrypto(),
			ReadWriter:        l,
			ConfigTXValidator: mockValidator,
			fakeConfig:        fakeConfig,
			bccsp:             cryptoProvider,
		},
	)

	ctx := makeConfigTxFull(genesisconfig.TestChannelID, 1)
	block := protoutil.NewBlock(1, protoutil.BlockHeaderHash(genesisBlockSys.Header))
	block.Data.Data = [][]byte{protoutil.MarshalOrPanic(ctx)}
	consenterMetadata := []byte("foo")
	bw.WriteConfigBlock(block, consenterMetadata)

	
	bw.committingBlock.Lock()
	bw.committingBlock.Unlock()

	cBlock := blockledger.GetBlock(l, block.Header.Number)
	assert.Equal(t, block.Header, cBlock.Header)
	assert.Equal(t, block.Data, cBlock.Data)

	omd, err := protoutil.GetConsenterMetadataFromBlock(block)
	require.NoError(t, err)
	assert.Equal(t, consenterMetadata, omd.Value)
}

func TestMigrationWriteConfig(t *testing.T) {
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()
	_, l := newRAMLedgerAndFactory(10, genesisconfig.TestChannelID, genesisBlockSys)

	fakeConfig := &mock.OrdererConfig{}
	fakeConfig.ConsensusTypeReturns("solo")
	fakeConfig.ConsensusStateReturns(orderer.ConsensusType_STATE_MAINTENANCE)

	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)

	mockValidator := &mocks.ConfigTXValidator{}
	mockValidator.ChannelIDReturns(genesisconfig.TestChannelID)
	bw := newBlockWriter(genesisBlockSys, nil,
		&mockBlockWriterSupport{
			SignerSerializer:  mockCrypto(),
			ReadWriter:        l,
			ConfigTXValidator: mockValidator,
			fakeConfig:        fakeConfig,
			bccsp:             cryptoProvider,
		},
	)

	ctx := makeConfigTxMig(genesisconfig.TestChannelID, 1)
	block := protoutil.NewBlock(1, protoutil.BlockHeaderHash(genesisBlockSys.Header))
	block.Data.Data = [][]byte{protoutil.MarshalOrPanic(ctx)}
	consenterMetadata := []byte("foo")

	bw.WriteConfigBlock(block, consenterMetadata)

	
	bw.committingBlock.Lock()
	bw.committingBlock.Unlock()

	cBlock := blockledger.GetBlock(l, block.Header.Number)
	assert.Equal(t, block.Header, cBlock.Header)
	assert.Equal(t, block.Data, cBlock.Data)

	omd := protoutil.GetMetadataFromBlockOrPanic(block, cb.BlockMetadataIndex_ORDERER)
	assert.Equal(t, []byte(nil), omd.Value)
}

func TestRaceWriteConfig(t *testing.T) {
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()
	_, l := newRAMLedgerAndFactory(10, genesisconfig.TestChannelID, genesisBlockSys)

	fakeConfig := &mock.OrdererConfig{}
	fakeConfig.ConsensusTypeReturns("solo")

	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)

	mockValidator := &mocks.ConfigTXValidator{}
	bw := newBlockWriter(genesisBlockSys, nil,
		&mockBlockWriterSupport{
			SignerSerializer:  mockCrypto(),
			ReadWriter:        l,
			ConfigTXValidator: mockValidator,
			fakeConfig:        fakeConfig,
			bccsp:             cryptoProvider,
		},
	)

	ctx := makeConfigTxFull(genesisconfig.TestChannelID, 1)
	block1 := protoutil.NewBlock(1, protoutil.BlockHeaderHash(genesisBlockSys.Header))
	block1.Data.Data = [][]byte{protoutil.MarshalOrPanic(ctx)}
	consenterMetadata1 := []byte("foo")
	mockValidator.SequenceReturnsOnCall(1, 1)

	ctx = makeConfigTxFull(genesisconfig.TestChannelID, 1)
	block2 := protoutil.NewBlock(2, protoutil.BlockHeaderHash(block1.Header))
	block2.Data.Data = [][]byte{protoutil.MarshalOrPanic(ctx)}
	consenterMetadata2 := []byte("bar")
	mockValidator.SequenceReturnsOnCall(2, 2)

	bw.WriteConfigBlock(block1, consenterMetadata1)
	bw.WriteConfigBlock(block2, consenterMetadata2)

	
	bw.committingBlock.Lock()
	bw.committingBlock.Unlock()

	cBlock := blockledger.GetBlock(l, block1.Header.Number)
	assert.Equal(t, block1.Header, cBlock.Header)
	assert.Equal(t, block1.Data, cBlock.Data)
	expectedLastConfigBlockNumber := block1.Header.Number
	testLastConfigBlockNumber(t, block1, expectedLastConfigBlockNumber)

	cBlock = blockledger.GetBlock(l, block2.Header.Number)
	assert.Equal(t, block2.Header, cBlock.Header)
	assert.Equal(t, block2.Data, cBlock.Data)
	expectedLastConfigBlockNumber = block2.Header.Number
	testLastConfigBlockNumber(t, block2, expectedLastConfigBlockNumber)

	omd, err := protoutil.GetConsenterMetadataFromBlock(block1)
	require.NoError(t, err)
	assert.Equal(t, consenterMetadata1, omd.Value)
}
