/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	etcdraftproto "github.com/mcc-github/blockchain-protos-go/orderer/etcdraft"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/cluster/mocks"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/raft/raftpb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestIsConsenterOfChannel(t *testing.T) {
	certInsideConfigBlock, err := base64.StdEncoding.DecodeString("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNmekNDQWlhZ0F3SUJBZ0l" +
		"SQUo4bjFLYTVzS1ZaTXRMTHJ1dldERDB3Q2dZSUtvWkl6ajBFQXdJd2JERUwKTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnVENrTmhiR" +
		"2xtYjNKdWFXRXhGakFVQmdOVkJBY1REVk5oYmlCRwpjbUZ1WTJselkyOHhGREFTQmdOVkJBb1RDMlY0WVcxd2JHVXVZMjl0TVJvd0dBWUR" +
		"WUVFERXhGMGJITmpZUzVsCmVHRnRjR3hsTG1OdmJUQWVGdzB4T0RFeE1EWXdPVFE1TURCYUZ3MHlPREV4TURNd09UUTVNREJhTUZreEN6QU" +
		"oKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFJRXdwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSEV3MVRZVzRnUm5KaApibU5wYzJOdk1SMHdH" +
		"d1lEVlFRREV4UnZjbVJsY21WeU1TNWxlR0Z0Y0d4bExtTnZiVEJaTUJNR0J5cUdTTTQ5CkFnRUdDQ3FHU000OUF3RUhBMElBQkRUVlFZc0" +
		"ZKZWxUcFZDMDFsek5DSkx6OENRMFFGVDBvN1BmSnBwSkl2SXgKUCtRVjQvRGRCSnRqQ0cvcGsvMGFxZXRpSjhZRUFMYmMrOUhmWnExN2tJ" +
		"Q2pnYnN3Z2Jnd0RnWURWUjBQQVFILwpCQVFEQWdXZ01CMEdBMVVkSlFRV01CUUdDQ3NHQVFVRkJ3TUJCZ2dyQmdFRkJRY0RBakFNQmdOV" +
		"khSTUJBZjhFCkFqQUFNQ3NHQTFVZEl3UWtNQ0tBSUVBOHFrSVJRTVBuWkxBR2g0TXZla2gzZFpHTmNxcEhZZWlXdzE3Rmw0ZlMKTUV3R0" +
		"ExVWRFUVJGTUVPQ0ZHOXlaR1Z5WlhJeExtVjRZVzF3YkdVdVkyOXRnZ2h2Y21SbGNtVnlNWUlKYkc5agpZV3hvYjNOMGh3Ui9BQUFCaHh" +
		"BQUFBQUFBQUFBQUFBQUFBQUFBQUFCTUFvR0NDcUdTTTQ5QkFNQ0EwY0FNRVFDCklFckJZRFVzV0JwOHB0ZVFSaTZyNjNVelhJQi81Sn" +
		"YxK0RlTkRIUHc3aDljQWlCakYrM3V5TzBvMEdRclB4MEUKUWptYlI5T3BVREN2LzlEUkNXWU9GZitkVlE9PQotLS0tLUVORCBDRVJUSU" +
		"ZJQ0FURS0tLS0tCg==")
	assert.NoError(t, err)

	validBlock := func() *common.Block {
		b, err := ioutil.ReadFile(filepath.Join("testdata", "etcdraftgenesis.block"))
		assert.NoError(t, err)
		block := &common.Block{}
		err = proto.Unmarshal(b, block)
		assert.NoError(t, err)
		return block
	}
	for _, testCase := range []struct {
		name          string
		expectedError string
		configBlock   *common.Block
		certificate   []byte
	}{
		{
			name:          "nil block",
			expectedError: "nil block",
		},
		{
			name:          "no block data",
			expectedError: "block data is nil",
			configBlock:   &common.Block{},
		},
		{
			name: "invalid envelope inside block",
			expectedError: "failed to unmarshal payload from envelope:" +
				" error unmarshaling Payload: proto: common.Payload: illegal tag 0 (wire type 1)",
			configBlock: &common.Block{
				Data: &common.BlockData{
					Data: [][]byte{protoutil.MarshalOrPanic(&common.Envelope{
						Payload: []byte{1, 2, 3},
					})},
				},
			},
		},
		{
			name:          "valid config block with cert mismatch",
			configBlock:   validBlock(),
			certificate:   certInsideConfigBlock[2:],
			expectedError: cluster.ErrNotInChannel.Error(),
		},
		{
			name:        "valid config block with matching cert",
			configBlock: validBlock(),
			certificate: certInsideConfigBlock,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			err := ConsenterCertificate(testCase.certificate).IsConsenterOfChannel(testCase.configBlock)
			if testCase.expectedError != "" {
				assert.EqualError(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPeriodicCheck(t *testing.T) {
	t.Parallel()

	g := gomega.NewGomegaWithT(t)

	var cond uint32
	var checkNum uint32

	fiveChecks := func() bool {
		return atomic.LoadUint32(&checkNum) > uint32(5)
	}

	condition := func() bool {
		atomic.AddUint32(&checkNum, 1)
		return atomic.LoadUint32(&cond) == uint32(1)
	}

	reports := make(chan time.Duration, 1000)

	report := func(duration time.Duration) {
		reports <- duration
	}

	check := &PeriodicCheck{
		Logger:        flogging.MustGetLogger("test"),
		Condition:     condition,
		CheckInterval: time.Millisecond,
		Report:        report,
	}

	go check.Run()

	g.Eventually(fiveChecks, time.Minute, time.Millisecond).Should(gomega.BeTrue())
	
	atomic.StoreUint32(&cond, 1)
	g.Eventually(reports, time.Minute, time.Millisecond).Should(gomega.Not(gomega.BeEmpty()))
	
	firstReport := <-reports
	g.Eventually(reports, time.Minute, time.Millisecond).Should(gomega.Not(gomega.BeEmpty()))
	
	secondReport := <-reports
	
	g.Expect(secondReport).To(gomega.BeNumerically(">", firstReport))
	
	g.Eventually(func() int { return len(reports) }, time.Minute, time.Millisecond).Should(gomega.BeNumerically("==", 1000))

	
	atomic.StoreUint32(&cond, 0)

	var lastReport time.Duration
	
	for len(reports) > 0 {
		select {
		case report := <-reports:
			lastReport = report
		default:
			break
		}
	}

	
	checksDoneSoFar := atomic.LoadUint32(&checkNum)
	g.Consistently(reports, time.Second*2, time.Millisecond).Should(gomega.BeEmpty())
	checksDoneAfter := atomic.LoadUint32(&checkNum)
	g.Expect(checksDoneAfter).To(gomega.BeNumerically(">", checksDoneSoFar))
	
	g.Expect(reports).To(gomega.BeEmpty())

	
	atomic.StoreUint32(&cond, 1)
	g.Eventually(reports, time.Minute, time.Millisecond).Should(gomega.Not(gomega.BeEmpty()))
	
	
	firstReport = <-reports
	g.Expect(lastReport).To(gomega.BeNumerically(">", firstReport))
	
	check.Stop()
	checkCountAfterStop := atomic.LoadUint32(&checkNum)
	
	time.Sleep(check.CheckInterval * 50)
	
	g.Expect(atomic.LoadUint32(&checkNum)).To(gomega.BeNumerically("<", checkCountAfterStop+2))
}

func TestEvictionSuspector(t *testing.T) {
	configBlock := &common.Block{
		Header: &common.BlockHeader{Number: 9},
		Metadata: &common.BlockMetadata{
			Metadata: [][]byte{{}, {}, {}, {}},
		},
	}
	configBlock.Metadata.Metadata[common.BlockMetadataIndex_LAST_CONFIG] = protoutil.MarshalOrPanic(&common.Metadata{
		Value: protoutil.MarshalOrPanic(&common.LastConfig{Index: 9}),
	})

	puller := &mocks.ChainPuller{}
	puller.On("Close")
	puller.On("HeightsByEndpoints").Return(map[string]uint64{"foo": 10}, nil)
	puller.On("PullBlock", uint64(9)).Return(configBlock)

	for _, testCase := range []struct {
		description                 string
		expectedPanic               string
		expectedLog                 string
		expectedCommittedBlockCount int
		amIInChannelReturns         error
		evictionSuspicionThreshold  time.Duration
		blockPuller                 BlockPuller
		blockPullerErr              error
		height                      uint64
		halt                        func()
	}{
		{
			description:                "suspected time is lower than threshold",
			evictionSuspicionThreshold: 11 * time.Minute,
			halt:                       t.Fail,
		},
		{
			description:                "puller creation fails",
			evictionSuspicionThreshold: 10*time.Minute - time.Second,
			blockPullerErr:             errors.New("oops"),
			expectedPanic:              "Failed creating a block puller: oops",
			halt:                       t.Fail,
		},
		{
			description:                "our height is the highest",
			expectedLog:                "Our height is higher or equal than the height of the orderer we pulled the last block from, aborting",
			evictionSuspicionThreshold: 10*time.Minute - time.Second,
			blockPuller:                puller,
			height:                     10,
			halt:                       t.Fail,
		},
		{
			description:                "failed pulling the block",
			expectedLog:                "Cannot confirm our own eviction from the channel: bad block",
			evictionSuspicionThreshold: 10*time.Minute - time.Second,
			amIInChannelReturns:        errors.New("bad block"),
			blockPuller:                puller,
			height:                     9,
			halt:                       t.Fail,
		},
		{
			description:                "we are still in the channel",
			expectedLog:                "Cannot confirm our own eviction from the channel, our certificate was found in config block with sequence 9",
			evictionSuspicionThreshold: 10*time.Minute - time.Second,
			amIInChannelReturns:        nil,
			blockPuller:                puller,
			height:                     9,
			halt:                       t.Fail,
		},
		{
			description:                 "we are not in the channel",
			expectedLog:                 "Detected our own eviction from the channel in block [9]",
			evictionSuspicionThreshold:  10*time.Minute - time.Second,
			amIInChannelReturns:         cluster.ErrNotInChannel,
			blockPuller:                 puller,
			height:                      8,
			expectedCommittedBlockCount: 2,
			halt: func() {
				puller.On("PullBlock", uint64(8)).Return(&common.Block{
					Header: &common.BlockHeader{Number: 8},
					Metadata: &common.BlockMetadata{
						Metadata: [][]byte{{}, {}, {}, {}},
					},
				})
			},
		},
	} {
		testCase := testCase
		t.Run(testCase.description, func(t *testing.T) {
			committedBlocks := make(chan *common.Block, 2)

			commitBlock := func(block *common.Block) error {
				committedBlocks <- block
				return nil
			}

			es := &evictionSuspector{
				halt: testCase.halt,
				amIInChannel: func(_ *common.Block) error {
					return testCase.amIInChannelReturns
				},
				evictionSuspicionThreshold: testCase.evictionSuspicionThreshold,
				createPuller: func() (BlockPuller, error) {
					return testCase.blockPuller, testCase.blockPullerErr
				},
				writeBlock: commitBlock,
				height: func() uint64 {
					return testCase.height
				},
				logger:         flogging.MustGetLogger("test"),
				triggerCatchUp: func(sn *raftpb.Snapshot) { return },
			}

			foundExpectedLog := testCase.expectedLog == ""
			es.logger = es.logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
				if strings.Contains(entry.Message, testCase.expectedLog) {
					foundExpectedLog = true
				}
				return nil
			}))

			runTestCase := func() {
				es.confirmSuspicion(time.Minute * 10)
			}

			if testCase.expectedPanic != "" {
				assert.PanicsWithValue(t, testCase.expectedPanic, runTestCase)
			} else {
				runTestCase()
				
				
				
				
				
				runTestCase()
			}

			assert.True(t, foundExpectedLog, "expected to find %s but didn't", testCase.expectedLog)
			assert.Equal(t, testCase.expectedCommittedBlockCount, len(committedBlocks))
		})
	}
}

func TestCheckConfigMetadata(t *testing.T) {
	tlsCA, err := tlsgen.NewCA()
	if err != nil {
		panic(err)
	}
	serverPair, err := tlsCA.NewServerCertKeyPair("localhost")
	serverCert := serverPair.Cert
	if err != nil {
		panic(err)
	}
	clientPair, err := tlsCA.NewClientCertKeyPair()
	clientCert := clientPair.Cert
	if err != nil {
		panic(err)
	}
	validOptions := &etcdraftproto.Options{
		TickInterval:         "500ms",
		ElectionTick:         10,
		HeartbeatTick:        1,
		MaxInflightBlocks:    5,
		SnapshotIntervalSize: 20 * 1024 * 1024, 
	}
	singleConsenter := &etcdraftproto.Consenter{
		Host:          "host1",
		Port:          10001,
		ClientTlsCert: clientCert,
		ServerTlsCert: serverCert,
	}

	
	goodMetadata := &etcdraftproto.ConfigMetadata{
		Options: validOptions,
		Consenters: []*etcdraftproto.Consenter{
			singleConsenter,
		},
	}
	assert.Nil(t, CheckConfigMetadata(goodMetadata))

	
	for _, testCase := range []struct {
		description string
		metadata    *etcdraftproto.ConfigMetadata
		errRegex    string
	}{
		{
			description: "nil metadata",
			metadata:    nil,
			errRegex:    "nil Raft config metadata",
		},
		{
			description: "HeartbeatTick is 0",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: &etcdraftproto.Options{
					HeartbeatTick: 0,
				},
			},
			errRegex: "none of HeartbeatTick .* can be zero",
		},
		{
			description: "ElectionTick is 0",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: &etcdraftproto.Options{
					HeartbeatTick: validOptions.HeartbeatTick,
					ElectionTick:  0,
				},
			},
			errRegex: "none of .* ElectionTick .* can be zero",
		},
		{
			description: "MaxInflightBlocks is 0",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: &etcdraftproto.Options{
					HeartbeatTick:     validOptions.HeartbeatTick,
					ElectionTick:      validOptions.ElectionTick,
					MaxInflightBlocks: 0,
				},
			},
			errRegex: "none of .* MaxInflightBlocks .* can be zero",
		},
		{
			description: "ElectionTick is less than HeartbeatTick",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: &etcdraftproto.Options{
					HeartbeatTick:     10,
					ElectionTick:      1,
					MaxInflightBlocks: validOptions.MaxInflightBlocks,
				},
			},
			errRegex: "ElectionTick .* must be greater than HeartbeatTick",
		},
		{
			description: "TickInterval is not parsable",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: &etcdraftproto.Options{
					HeartbeatTick:     validOptions.HeartbeatTick,
					ElectionTick:      validOptions.ElectionTick,
					MaxInflightBlocks: validOptions.MaxInflightBlocks,
					TickInterval:      "abcd",
				},
			},
			errRegex: "failed to parse TickInterval .* to time duration",
		},
		{
			description: "TickInterval is 0",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: &etcdraftproto.Options{
					HeartbeatTick:     validOptions.HeartbeatTick,
					ElectionTick:      validOptions.ElectionTick,
					MaxInflightBlocks: validOptions.MaxInflightBlocks,
					TickInterval:      "0s",
				},
			},
			errRegex: "TickInterval cannot be zero",
		},
		{
			description: "consenter set is empty",
			metadata: &etcdraftproto.ConfigMetadata{
				Options:    validOptions,
				Consenters: []*etcdraftproto.Consenter{},
			},
			errRegex: "empty consenter set",
		},
		{
			description: "metadata has nil consenter",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: validOptions,
				Consenters: []*etcdraftproto.Consenter{
					nil,
				},
			},
			errRegex: "metadata has nil consenter",
		},
		{
			description: "consenter has invalid server cert",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: validOptions,
				Consenters: []*etcdraftproto.Consenter{
					{
						ServerTlsCert: []byte("invalid"),
						ClientTlsCert: clientCert,
					},
				},
			},
			errRegex: "server TLS certificate is not PEM encoded",
		},
		{
			description: "consenter has invalid client cert",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: validOptions,
				Consenters: []*etcdraftproto.Consenter{
					{
						ServerTlsCert: serverCert,
						ClientTlsCert: []byte("invalid"),
					},
				},
			},
			errRegex: "client TLS certificate is not PEM encoded",
		},
		{
			description: "metadata has duplicate consenters",
			metadata: &etcdraftproto.ConfigMetadata{
				Options: validOptions,
				Consenters: []*etcdraftproto.Consenter{
					singleConsenter,
					singleConsenter,
				},
			},
			errRegex: "duplicate consenter",
		},
	} {
		err := CheckConfigMetadata(testCase.metadata)
		assert.NotNil(t, err, testCase.description)
		assert.Regexp(t, testCase.errRegex, err)
	}
}
