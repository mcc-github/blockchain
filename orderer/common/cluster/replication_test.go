/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/mocks/crypto"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/cluster/mocks"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsReplicationNeeded(t *testing.T) {
	for _, testCase := range []struct {
		name                string
		bootBlock           *common.Block
		systemChannelHeight uint64
		systemChannelError  error
		expectedError       string
		replicationNeeded   bool
	}{
		{
			name:                "no replication needed",
			systemChannelHeight: 100,
			bootBlock:           &common.Block{Header: &common.BlockHeader{Number: 99}},
		},
		{
			name:                "replication is needed",
			systemChannelHeight: 99,
			bootBlock:           &common.Block{Header: &common.BlockHeader{Number: 99}},
			replicationNeeded:   true,
		},
		{
			name:               "IO error",
			systemChannelError: errors.New("IO error"),
			expectedError:      "IO error",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			ledgerWriter := &mocks.LedgerWriter{}
			ledgerWriter.On("Height").Return(testCase.systemChannelHeight)

			ledgerFactory := &mocks.LedgerFactory{}
			ledgerFactory.On("Close")
			ledgerFactory.On("GetOrCreate", "system").Return(ledgerWriter, testCase.systemChannelError)

			r := cluster.Replicator{
				Logger:        flogging.MustGetLogger("test"),
				BootBlock:     testCase.bootBlock,
				SystemChannel: "system",
				LedgerFactory: ledgerFactory,
			}

			ok, err := r.IsReplicationNeeded()
			if testCase.expectedError != "" {
				assert.EqualError(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.replicationNeeded, ok)
			}
			
			ledgerFactory.AssertCalled(t, "Close")
		})
	}
}

func TestReplicateChainsFailures(t *testing.T) {
	for _, testCase := range []struct {
		name                    string
		blocks                  []*common.Block
		expectedError           string
		latestBlockSeqInOrderer uint64
		ledgerFactoryError      error
		appendBlockError        error
		expectedPanic           string
		mutateBlocks            func([]*common.Block)
	}{
		{
			name:          "no block received",
			expectedError: "failed obtaining the latest block for channel system",
		},
		{
			name: "latest block seq is less than boot block seq",
			expectedError: "latest height found among system channel(system) orderers is 19," +
				" but the boot block's sequence is 21",
			latestBlockSeqInOrderer: 18,
		},
		{
			name: "hash chain mismatch",
			expectedError: "block header mismatch on sequence 11, " +
				"expected 9cd61b7e9a5ea2d128cc877e5304e7205888175a8032d40b97db7412dca41d9e, got 010203",
			latestBlockSeqInOrderer: 21,
			mutateBlocks: func(systemChannelBlocks []*common.Block) {
				systemChannelBlocks[len(systemChannelBlocks)/2].Header.PreviousHash = []byte{1, 2, 3}
			},
		},
		{
			name: "last pulled block doesn't match the boot block",
			expectedPanic: "Block header mismatch on last system channel block," +
				" expected 8ec93b2ef5ffdc302f0c0e24611be04ad2b17b099a1aeafd7cfb76a95923f146," +
				" got e428decfc78f8e4c97b26da9c16f9d0b73f886dafa80477a0dd9bac7eb14fe7a",
			latestBlockSeqInOrderer: 21,
			mutateBlocks: func(systemChannelBlocks []*common.Block) {
				systemChannelBlocks[21].Header.DataHash = nil
			},
		},
		{
			name:                    "failure in creating ledger",
			latestBlockSeqInOrderer: 21,
			ledgerFactoryError:      errors.New("IO error"),
			expectedPanic:           "Failed to create a ledger for channel system: IO error",
		},
		{
			name:                    "failure in appending a block to the ledger",
			latestBlockSeqInOrderer: 21,
			appendBlockError:        errors.New("IO error"),
			expectedPanic:           "Failed to write block 0: IO error",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			systemChannelBlocks := createBlockChain(0, 21)
			if testCase.mutateBlocks != nil {
				testCase.mutateBlocks(systemChannelBlocks)
			}

			lw := &mocks.LedgerWriter{}
			lw.On("Append", mock.Anything).Return(testCase.appendBlockError)

			lf := &mocks.LedgerFactory{}
			lf.On("GetOrCreate", "system").Return(lw, testCase.ledgerFactoryError)

			osn := newClusterNode(t)
			defer osn.stop()

			dialer := newCountingDialer()
			bp := newBlockPuller(dialer, osn.srv.Address())
			bp.FetchTimeout = time.Millisecond * 100

			r := cluster.Replicator{
				Logger:        flogging.MustGetLogger("test"),
				BootBlock:     systemChannelBlocks[21],
				SystemChannel: "system",
				LedgerFactory: lf,
				Puller:        bp,
			}

			osn.addExpectProbeAssert()
			osn.enqueueResponse(testCase.latestBlockSeqInOrderer)
			osn.addExpectProbeAssert()
			osn.enqueueResponse(testCase.latestBlockSeqInOrderer)
			osn.addExpectPullAssert(0)
			for _, block := range systemChannelBlocks {
				osn.blockResponses <- &orderer.DeliverResponse{
					Type: &orderer.DeliverResponse_Block{Block: block},
				}
			}

			if testCase.expectedPanic == "" {
				err := r.PullChannel("system")
				assert.EqualError(t, err, testCase.expectedError)
			} else {
				assert.PanicsWithValue(t, testCase.expectedPanic, func() {
					r.PullChannel("system")
				})
			}

			bp.Close()
			dialer.assertAllConnectionsClosed(t)
		})
	}
}

func TestPullerConfigFromTopLevelConfig(t *testing.T) {
	signer := &crypto.LocalSigner{}
	expected := cluster.PullerConfig{
		Channel:             "system",
		MaxTotalBufferBytes: 100,
		Signer:              signer,
		TLSCert:             []byte{3, 2, 1},
		TLSKey:              []byte{1, 2, 3},
		Timeout:             time.Hour,
	}

	topLevelConfig := &localconfig.TopLevel{
		General: localconfig.General{
			SystemChannel: "system",
			Cluster: localconfig.Cluster{
				ReplicationBufferSize: 100,
				RPCTimeout:            time.Hour,
			},
		},
	}

	config := cluster.PullerConfigFromTopLevelConfig(topLevelConfig, []byte{1, 2, 3}, []byte{3, 2, 1}, signer)
	assert.Equal(t, expected, config)
}

func TestReplicateChainsChannelClassificationFailure(t *testing.T) {
	
	

	block30WithConfigBlockOf21 := common.NewBlock(30, nil)
	block30WithConfigBlockOf21.Metadata.Metadata[common.BlockMetadataIndex_LAST_CONFIG] = utils.MarshalOrPanic(&common.Metadata{
		Value: utils.MarshalOrPanic(&common.LastConfig{Index: 21}),
	})

	osn := newClusterNode(t)
	defer osn.stop()
	osn.blockResponses = make(chan *orderer.DeliverResponse, 1000)

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn.srv.Address())
	bp.FetchTimeout = time.Hour

	channelLister := &mocks.ChannelLister{}
	channelLister.On("Channels").Return([]string{"A"})
	channelLister.On("Close")

	
	osn.addExpectProbeAssert()
	osn.enqueueResponse(30)

	
	osn.addExpectProbeAssert()
	osn.enqueueResponse(30)
	osn.addExpectPullAssert(30)
	osn.blockResponses <- &orderer.DeliverResponse{
		Type: &orderer.DeliverResponse_Block{Block: block30WithConfigBlockOf21},
	}
	
	
	
	osn.blockResponses <- nil
	
	osn.addExpectProbeAssert()
	osn.enqueueResponse(30)
	
	osn.enqueueResponse(21)
	osn.addExpectPullAssert(21)

	r := cluster.Replicator{
		AmIPartOfChannel: func(configBlock *common.Block) error {
			return errors.New("oops")
		},
		Logger:        flogging.MustGetLogger("test"),
		SystemChannel: "system",
		ChannelLister: channelLister,
		Puller:        bp,
	}

	assert.PanicsWithValue(t, "Failed classifying whether I belong to channel A: oops, skipping chain retrieval", func() {
		r.ReplicateChains()
	})

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestReplicateChainsGreenPath(t *testing.T) {
	
	
	

	systemChannelBlocks := createBlockChain(0, 21)
	block30WithConfigBlockOf21 := common.NewBlock(30, nil)
	block30WithConfigBlockOf21.Metadata.Metadata[common.BlockMetadataIndex_LAST_CONFIG] = utils.MarshalOrPanic(&common.Metadata{
		Value: utils.MarshalOrPanic(&common.LastConfig{Index: 21}),
	})

	osn := newClusterNode(t)
	defer osn.stop()
	osn.blockResponses = make(chan *orderer.DeliverResponse, 1000)

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn.srv.Address())
	bp.FetchTimeout = time.Hour

	channelLister := &mocks.ChannelLister{}
	channelLister.On("Channels").Return([]string{"A", "B"})
	channelLister.On("Close")

	amIPartOfChannelMock := &mock.Mock{}
	
	amIPartOfChannelMock.On("func2").Return(nil).Once()
	
	amIPartOfChannelMock.On("func2").Return(cluster.NotInChannelError).Once()

	
	blocksCommittedToLedger := make(chan *common.Block, 22+31)

	lw := &mocks.LedgerWriter{}
	lw.On("Append", mock.Anything).Return(nil).Run(func(arg mock.Arguments) {
		blocksCommittedToLedger <- arg.Get(0).(*common.Block)
	})

	lf := &mocks.LedgerFactory{}
	lf.On("Close")
	lf.On("GetOrCreate", "A").Return(lw, nil)
	lf.On("GetOrCreate", "B").Return(lw, nil)
	lf.On("GetOrCreate", "system").Return(lw, nil)

	r := cluster.Replicator{
		LedgerFactory: lf,
		AmIPartOfChannel: func(configBlock *common.Block) error {
			return amIPartOfChannelMock.Called().Error(0)
		},
		Logger:        flogging.MustGetLogger("test"),
		SystemChannel: "system",
		ChannelLister: channelLister,
		Puller:        bp,
		BootBlock:     systemChannelBlocks[21],
	}

	for _, channel := range []string{"A", "B"} {
		channel := channel
		
		
		

		
		osn.seekAssertions <- func(info *orderer.SeekInfo, actualChannel string) {
			
			assert.NotNil(osn.t, info.GetStart().GetNewest())
			assert.Equal(t, channel, actualChannel)
		}

		
		
		osn.enqueueResponse(30)

		
		osn.addExpectProbeAssert()
		osn.enqueueResponse(30)

		
		osn.addExpectPullAssert(30)
		osn.blockResponses <- &orderer.DeliverResponse{
			Type: &orderer.DeliverResponse_Block{Block: block30WithConfigBlockOf21},
		}
		
		osn.blockResponses <- nil
		
		osn.addExpectProbeAssert()
		osn.enqueueResponse(30)
		
		osn.enqueueResponse(21)
		osn.addExpectPullAssert(21)
		
		osn.blockResponses <- nil
	}

	
	

	
	osn.seekAssertions <- func(info *orderer.SeekInfo, actualChannel string) {
		
		assert.NotNil(osn.t, info.GetStart().GetNewest())
		assert.Equal(t, "A", actualChannel)
	}
	osn.enqueueResponse(30)
	
	osn.enqueueResponse(30)
	osn.addExpectProbeAssert()
	osn.addExpectPullAssert(0)
	
	for _, block := range createBlockChain(0, 30) {
		osn.blockResponses <- &orderer.DeliverResponse{
			Type: &orderer.DeliverResponse_Block{Block: block},
		}
	}
	
	
	osn.blockResponses <- nil

	
	
	osn.seekAssertions <- func(info *orderer.SeekInfo, actualChannel string) {
		
		assert.NotNil(osn.t, info.GetStart().GetNewest())
		assert.Equal(t, "system", actualChannel)
	}
	osn.blockResponses <- &orderer.DeliverResponse{
		Type: &orderer.DeliverResponse_Block{Block: systemChannelBlocks[21]},
	}
	osn.addExpectProbeAssert()
	osn.enqueueResponse(21)
	osn.addExpectPullAssert(0)
	for _, block := range systemChannelBlocks {
		osn.blockResponses <- &orderer.DeliverResponse{
			Type: &orderer.DeliverResponse_Block{Block: block},
		}
	}

	
	
	
	r.ReplicateChains()

	
	
	
	close(blocksCommittedToLedger)
	assert.Len(t, blocksCommittedToLedger, cap(blocksCommittedToLedger))
	
	var expectedSequence uint64
	for block := range blocksCommittedToLedger {
		assert.Equal(t, expectedSequence, block.Header.Number)
		expectedSequence++
		if expectedSequence == 31 {
			break
		}
	}

	
	expectedSequence = uint64(0)
	for block := range blocksCommittedToLedger {
		assert.Equal(t, expectedSequence, block.Header.Number)
		expectedSequence++
	}

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
	lf.AssertNumberOfCalls(t, "Close", 1)
}

func TestParticipant(t *testing.T) {
	for _, testCase := range []struct {
		name                      string
		heightsByEndpointsReturns map[string]uint64
		latestBlockSeq            uint64
		latestBlock               *common.Block
		latestConfigBlockSeq      uint64
		latestConfigBlock         *common.Block
		expectedError             string
		predicateReturns          error
	}{
		{
			name:          "No available orderer",
			expectedError: "no available orderer",
		},
		{
			name: "Pulled block has no metadata",
			heightsByEndpointsReturns: map[string]uint64{
				"orderer.example.com:7050": 100,
			},
			latestBlockSeq: uint64(99),
			latestBlock:    &common.Block{},
			expectedError:  "no metadata in block",
		},
		{
			name: "Pulled block has no last config sequence in metadata",
			heightsByEndpointsReturns: map[string]uint64{
				"orderer.example.com:7050": 100,
			},
			latestBlockSeq: uint64(99),
			latestBlock: &common.Block{
				Metadata: &common.BlockMetadata{
					Metadata: [][]byte{{1, 2, 3}},
				},
			},
			expectedError: "no metadata in block",
		},
		{
			name: "Pulled block's metadata is malformed",
			heightsByEndpointsReturns: map[string]uint64{
				"orderer.example.com:7050": 100,
			},
			latestBlockSeq: uint64(99),
			latestBlock: &common.Block{
				Metadata: &common.BlockMetadata{
					Metadata: [][]byte{{1, 2, 3}, {1, 2, 3}},
				},
			},
			expectedError: "error unmarshaling metadata from" +
				" block at index [LAST_CONFIG]: proto: common.Metadata: illegal tag 0 (wire type 1)",
		},
		{
			name: "Pulled block's metadata is valid and has a last config",
			heightsByEndpointsReturns: map[string]uint64{
				"orderer.example.com:7050": 100,
			},
			latestBlockSeq: uint64(99),
			latestBlock: &common.Block{
				Metadata: &common.BlockMetadata{
					Metadata: [][]byte{{1, 2, 3}, utils.MarshalOrPanic(&common.Metadata{
						Value: utils.MarshalOrPanic(&common.LastConfig{
							Index: 42,
						}),
					})},
				},
			},
			latestConfigBlockSeq: 42,
			latestConfigBlock:    &common.Block{Header: &common.BlockHeader{Number: 42}},
			predicateReturns:     cluster.NotInChannelError,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			configBlocks := make(chan *common.Block, 1)
			predicate := func(configBlock *common.Block) error {
				configBlocks <- configBlock
				return testCase.predicateReturns
			}
			puller := &mocks.ChainPuller{}
			puller.On("HeightsByEndpoints").Return(testCase.heightsByEndpointsReturns)
			puller.On("PullBlock", testCase.latestBlockSeq).Return(testCase.latestBlock)
			puller.On("PullBlock", testCase.latestConfigBlockSeq).Return(testCase.latestConfigBlock)
			puller.On("Close")

			err := cluster.Participant(puller, predicate)
			if testCase.expectedError != "" {
				assert.EqualError(t, err, testCase.expectedError)
				assert.Len(t, configBlocks, 0)
			} else {
				assert.Len(t, configBlocks, 1)
				assert.Equal(t, err, testCase.predicateReturns)
			}
		})
	}
}

func TestBlockPullerFromConfigBlockFailures(t *testing.T) {
	blockBytes, err := ioutil.ReadFile("testdata/mychannel.block")
	assert.NoError(t, err)

	validBlock := &common.Block{}
	assert.NoError(t, proto.Unmarshal(blockBytes, validBlock))

	for _, testCase := range []struct {
		name         string
		expectedErr  string
		pullerConfig cluster.PullerConfig
		block        *common.Block
	}{
		{
			name:        "nil block",
			expectedErr: "nil block",
		},
		{
			name:        "invalid block",
			expectedErr: "block data is nil",
			block:       &common.Block{},
		},
		{
			name: "bad envelope inside block",
			expectedErr: "failed extracting bundle from envelope: " +
				"failed to unmarshal payload from envelope: " +
				"error unmarshaling Payload: " +
				"proto: common.Payload: illegal tag 0 (wire type 1)",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: []byte{1, 2, 3},
					})},
				},
			},
		},
		{
			name:        "invalid TLS certificate",
			expectedErr: "unable to decode TLS certificate PEM: 
			block:       validBlock,
			pullerConfig: cluster.PullerConfig{
				TLSCert: []byte{255, 255, 255},
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			bp, err := cluster.BlockPullerFromConfigBlock(testCase.pullerConfig, testCase.block)
			assert.EqualError(t, err, testCase.expectedErr)
			assert.Nil(t, bp)
		})
	}
}

func TestBlockPullerFromConfigBlockGreenPath(t *testing.T) {
	caCert, err := ioutil.ReadFile(filepath.Join("testdata", "ca.crt"))
	assert.NoError(t, err)

	tlsCert, err := ioutil.ReadFile(filepath.Join("testdata", "server.crt"))
	assert.NoError(t, err)

	tlsKey, err := ioutil.ReadFile(filepath.Join("testdata", "server.key"))
	assert.NoError(t, err)

	osn := newClusterNode(t)
	osn.srv.Stop()
	
	osn.srv, err = comm.NewGRPCServer("127.0.0.1:0", comm.ServerConfig{
		SecOpts: &comm.SecureOptions{
			Key:               tlsKey,
			RequireClientCert: true,
			Certificate:       tlsCert,
			ClientRootCAs:     [][]byte{caCert},
			UseTLS:            true,
		},
	})
	assert.NoError(t, err)
	orderer.RegisterAtomicBroadcastServer(osn.srv.Server(), osn)
	
	go osn.srv.Start()
	defer osn.stop()

	
	blockBytes, err := ioutil.ReadFile(filepath.Join("testdata", "mychannel.block"))
	assert.NoError(t, err)

	validBlock := &common.Block{}
	assert.NoError(t, proto.Unmarshal(blockBytes, validBlock))

	
	injectTLSCACert(t, validBlock, caCert)
	injectOrdererEndpoint(t, validBlock, osn.srv.Address())
	validBlock.Header.DataHash = validBlock.Data.Hash()

	blockMsg := &orderer.DeliverResponse_Block{
		Block: validBlock,
	}

	osn.blockResponses <- &orderer.DeliverResponse{
		Type: blockMsg,
	}

	osn.blockResponses <- &orderer.DeliverResponse{
		Type: blockMsg,
	}

	bp, err := cluster.BlockPullerFromConfigBlock(cluster.PullerConfig{
		TLSCert:             tlsCert,
		TLSKey:              tlsKey,
		MaxTotalBufferBytes: 1,
		Channel:             "mychannel",
		Signer:              &crypto.LocalSigner{},
		Timeout:             time.Second,
	}, validBlock)
	assert.NoError(t, err)
	defer bp.Close()

	osn.addExpectProbeAssert()
	osn.addExpectPullAssert(0)

	block := bp.PullBlock(0)
	assert.Equal(t, uint64(0), block.Header.Number)
}

func TestNoopBlockVerifier(t *testing.T) {
	v := &cluster.NoopBlockVerifier{}
	assert.Nil(t, v.VerifyBlockSignature(nil, nil))
}

func injectOrdererEndpoint(t *testing.T, block *common.Block, endpoint string) {
	ordererAddresses := channelconfig.OrdererAddressesValue([]string{endpoint})
	
	env, err := utils.ExtractEnvelope(block, 0)
	assert.NoError(t, err)
	payload, err := utils.ExtractPayload(env)
	assert.NoError(t, err)
	confEnv, err := configtx.UnmarshalConfigEnvelope(payload.Data)
	assert.NoError(t, err)
	
	confEnv.Config.ChannelGroup.Values[ordererAddresses.Key()].Value = utils.MarshalOrPanic(ordererAddresses.Value())
	
	payload.Data = utils.MarshalOrPanic(confEnv)
	env.Payload = utils.MarshalOrPanic(payload)
	block.Data.Data[0] = utils.MarshalOrPanic(env)
}

func injectTLSCACert(t *testing.T, block *common.Block, tlsCA []byte) {
	
	env, err := utils.ExtractEnvelope(block, 0)
	assert.NoError(t, err)
	payload, err := utils.ExtractPayload(env)
	assert.NoError(t, err)
	confEnv, err := configtx.UnmarshalConfigEnvelope(payload.Data)
	assert.NoError(t, err)
	mspKey := confEnv.Config.ChannelGroup.Groups[channelconfig.OrdererGroupKey].Groups["OrdererOrg"].Values[channelconfig.MSPKey]
	rawMSPConfig := mspKey.Value
	mspConf := &msp.MSPConfig{}
	proto.Unmarshal(rawMSPConfig, mspConf)
	blockchainMSPConf := &msp.FabricMSPConfig{}
	proto.Unmarshal(mspConf.Config, blockchainMSPConf)
	
	blockchainMSPConf.TlsRootCerts = [][]byte{tlsCA}
	
	mspConf.Config = utils.MarshalOrPanic(blockchainMSPConf)
	mspKey.Value = utils.MarshalOrPanic(mspConf)
	payload.Data = utils.MarshalOrPanic(confEnv)
	env.Payload = utils.MarshalOrPanic(payload)
	block.Data.Data[0] = utils.MarshalOrPanic(env)
}

func TestIsNewChannelBlock(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		expectedErr  string
		returnedName string
		block        *common.Block
	}{
		{
			name:        "nil block",
			expectedErr: "nil block",
		},
		{
			name:        "no data section in block",
			expectedErr: "block data is nil",
			block:       &common.Block{},
		},
		{
			name: "corrupt envelope in block",
			expectedErr: "block data does not carry an" +
				" envelope at index 0: error unmarshaling Envelope: proto: common.Envelope: illegal tag 0 (wire type 1)",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{{1, 2, 3}},
				},
			},
		},
		{
			name:        "corrupt payload in envelope",
			expectedErr: "no payload in envelope: proto: common.Payload: illegal tag 0 (wire type 1)",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: []byte{1, 2, 3},
					})},
				},
			},
		},
		{
			name:        "no header in block",
			expectedErr: "nil header in payload",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{}),
					})},
				},
			},
		},
		{
			name: "corrupt channel header",
			expectedErr: "error unmarshaling ChannelHeader:" +
				" proto: common.ChannelHeader: illegal tag 0 (wire type 1)",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: []byte{1, 2, 3},
							},
						}),
					})},
				},
			},
		},
		{
			name:        "not an orderer transaction",
			expectedErr: "",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
									Type: int32(common.HeaderType_CONFIG_UPDATE),
								}),
							},
						}),
					})},
				},
			},
		},
		{
			name:        "orderer transaction with corrupt inner envelope",
			expectedErr: "error unmarshaling Envelope: proto: common.Envelope: illegal tag 0 (wire type 1)",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
									Type: int32(common.HeaderType_ORDERER_TRANSACTION),
								}),
							},
							Data: []byte{1, 2, 3},
						}),
					})},
				},
			},
		},
		{
			name:        "orderer transaction with corrupt inner payload",
			expectedErr: "error unmarshaling Payload: proto: common.Payload: illegal tag 0 (wire type 1)",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
									Type: int32(common.HeaderType_ORDERER_TRANSACTION),
								}),
							},
							Data: utils.MarshalOrPanic(&common.Envelope{
								Payload: []byte{1, 2, 3},
							}),
						}),
					})},
				},
			},
		},
		{
			name:        "orderer transaction with nil inner header",
			expectedErr: "inner payload's header is nil",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
									Type: int32(common.HeaderType_ORDERER_TRANSACTION),
								}),
							},
							Data: utils.MarshalOrPanic(&common.Envelope{
								Payload: utils.MarshalOrPanic(&common.Payload{}),
							}),
						}),
					})},
				},
			},
		},
		{
			name:        "orderer transaction with corrupt inner channel header",
			expectedErr: "error unmarshaling ChannelHeader: proto: common.ChannelHeader: illegal tag 0 (wire type 1)",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
									Type: int32(common.HeaderType_ORDERER_TRANSACTION),
								}),
							},
							Data: utils.MarshalOrPanic(&common.Envelope{
								Payload: utils.MarshalOrPanic(&common.Payload{
									Header: &common.Header{
										ChannelHeader: []byte{1, 2, 3},
									},
								}),
							}),
						}),
					})},
				},
			},
		},
		{
			name:        "orderer transaction that is not a config, but a config update",
			expectedErr: "",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
									Type: int32(common.HeaderType_ORDERER_TRANSACTION),
								}),
							},
							Data: utils.MarshalOrPanic(&common.Envelope{
								Payload: utils.MarshalOrPanic(&common.Payload{
									Header: &common.Header{
										ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
											Type: int32(common.HeaderType_CONFIG_UPDATE),
										}),
									},
								}),
							}),
						}),
					})},
				},
			},
		},
		{
			expectedErr: "",
			name:        "orderer transaction that is a system channel config block",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
									ChannelId: "systemChannel",
									Type:      int32(common.HeaderType_ORDERER_TRANSACTION),
								}),
							},
							Data: utils.MarshalOrPanic(&common.Envelope{
								Payload: utils.MarshalOrPanic(&common.Payload{
									Header: &common.Header{
										ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
											Type:      int32(common.HeaderType_CONFIG),
											ChannelId: "systemChannel",
										}),
									},
								}),
							}),
						}),
					})},
				},
			},
		},
		{
			name:         "orderer transaction that creates a new application channel",
			expectedErr:  "",
			returnedName: "notSystemChannel",
			block: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
						Payload: utils.MarshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
									ChannelId: "systemChannel",
									Type:      int32(common.HeaderType_ORDERER_TRANSACTION),
								}),
							},
							Data: utils.MarshalOrPanic(&common.Envelope{
								Payload: utils.MarshalOrPanic(&common.Payload{
									Header: &common.Header{
										ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
											Type:      int32(common.HeaderType_CONFIG),
											ChannelId: "notSystemChannel",
										}),
									},
								}),
							}),
						}),
					})},
				},
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			channelName, err := cluster.IsNewChannelBlock(testCase.block)
			if testCase.expectedErr != "" {
				assert.EqualError(t, err, testCase.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.returnedName, channelName)
		})
	}
}

func TestChannels(t *testing.T) {
	makeBlock := func(outerChannelName, innerChannelName string) *common.Block {
		return &common.Block{
			Header: &common.BlockHeader{},
			Data: &common.BlockData{
				Data: [][]byte{utils.MarshalOrPanic(&common.Envelope{
					Payload: utils.MarshalOrPanic(&common.Payload{
						Header: &common.Header{
							ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
								ChannelId: outerChannelName,
								Type:      int32(common.HeaderType_ORDERER_TRANSACTION),
							}),
						},
						Data: utils.MarshalOrPanic(&common.Envelope{
							Payload: utils.MarshalOrPanic(&common.Payload{
								Header: &common.Header{
									ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{
										Type:      int32(common.HeaderType_CONFIG),
										ChannelId: innerChannelName,
									}),
								},
							}),
						}),
					}),
				})},
			},
		}
	}

	for _, testCase := range []struct {
		name               string
		prepareSystemChain func(systemChain []*common.Block)
		assertion          func(t *testing.T, ci *cluster.ChainInspector)
	}{
		{
			name: "happy path - artificial blocks",
			prepareSystemChain: func(systemChain []*common.Block) {
				assignHashes(systemChain)
			},
			assertion: func(t *testing.T, ci *cluster.ChainInspector) {
				actual := ci.Channels()
				
				assert.Contains(t, [][]string{{"mychannel", "mychannel2"}, {"mychannel2", "mychannel"}}, actual)
			},
		},
		{
			name: "happy path - one block is not artificial but real",
			prepareSystemChain: func(systemChain []*common.Block) {
				blockbytes, err := ioutil.ReadFile(filepath.Join("testdata", "block3.pb"))
				assert.NoError(t, err)
				block := &common.Block{}
				err = proto.Unmarshal(blockbytes, block)
				assert.NoError(t, err)

				systemChain[len(systemChain)/2] = block
				assignHashes(systemChain)
			},
			assertion: func(t *testing.T, ci *cluster.ChainInspector) {
				actual := ci.Channels()
				
				assert.Contains(t, [][]string{{"mychannel", "bar"}, {"bar", "mychannel"}}, actual)
			},
		},
		{
			name:               "bad path - pulled chain's hash is mismatched",
			prepareSystemChain: func(_ []*common.Block) {},
			assertion: func(t *testing.T, ci *cluster.ChainInspector) {
				panicValue := "System channel pulled doesn't match the boot last config block:" +
					" block 4's hash (d8553eb97aa57e3c795a185f30efdbe8d88ae4b1e44b984b311159beac9bd5f4)" +
					" mismatches 3's prev block hash ()"
				assert.PanicsWithValue(t, panicValue, func() {
					ci.Channels()
				})
			},
		},
		{
			name: "bad path - a block cannot be classified",
			prepareSystemChain: func(systemChain []*common.Block) {
				assignHashes(systemChain)
				systemChain[len(systemChain)-2].Data.Data = [][]byte{{1, 2, 3}}
			},
			assertion: func(t *testing.T, ci *cluster.ChainInspector) {
				panicValue := "Failed classifying block 3 : block data does not carry" +
					" an envelope at index 0: error unmarshaling Envelope: " +
					"proto: common.Envelope: illegal tag 0 (wire type 1)"
				assert.PanicsWithValue(t, panicValue, func() {
					ci.Channels()
				})
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			systemChain := []*common.Block{
				makeBlock("systemChannel", "systemChannel"),
				makeBlock("systemChannel", "mychannel"),
				makeBlock("systemChannel", "mychannel2"),
				makeBlock("systemChannel", "systemChannel"),
			}

			for i := 0; i < len(systemChain); i++ {
				systemChain[i].Header.DataHash = systemChain[i].Data.Hash()
				systemChain[i].Header.Number = uint64(i + 1)
			}
			testCase.prepareSystemChain(systemChain)
			puller := &mocks.ChainPuller{}
			for seq := uint64(1); int(seq) <= len(systemChain); seq++ {
				puller.On("PullBlock", seq).Return(systemChain[int(seq)-1])
			}

			ci := &cluster.ChainInspector{
				Logger:          flogging.MustGetLogger("test"),
				Puller:          puller,
				LastConfigBlock: systemChain[len(systemChain)-1],
			}
			testCase.assertion(t, ci)
		})
	}
}
