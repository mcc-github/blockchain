/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft_test

import (
	"encoding/pem"
	"sync"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"github.com/coreos/etcd/raft"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/flogging"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft/mocks"
	consensusmocks "github.com/mcc-github/blockchain/orderer/consensus/mocks"
	mockblockcutter "github.com/mcc-github/blockchain/orderer/mocks/common/blockcutter"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	raftprotos "github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

const (
	interval       = time.Second
	ELECTION_TICK  = 2
	HEARTBEAT_TICK = 1
)

var _ = Describe("Chain", func() {
	var (
		env         *common.Envelope
		normalBlock *common.Block
		channelID   string
		tlsCA       tlsgen.CA
		logger      *flogging.FabricLogger
	)

	BeforeEach(func() {
		tlsCA, _ = tlsgen.NewCA()
		channelID = "test-channel"
		logger = flogging.NewFabricLogger(zap.NewNop())
		env = &common.Envelope{
			Payload: marshalOrPanic(&common.Payload{
				Header: &common.Header{ChannelHeader: marshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_MESSAGE), ChannelId: channelID})},
				Data:   []byte("TEST_MESSAGE"),
			}),
		}
		normalBlock = &common.Block{Data: &common.BlockData{Data: [][]byte{[]byte("foo")}}}
	})

	Describe("Single Raft node", func() {
		var (
			configurator      *mocks.Configurator
			consenterMetadata *raftprotos.Metadata
			clock             *fakeclock.FakeClock
			opts              etcdraft.Options
			support           *consensusmocks.FakeConsenterSupport
			cutter            *mockblockcutter.Receiver
			storage           *raft.MemoryStorage
			observeC          chan uint64
			chain             *etcdraft.Chain
		)

		BeforeEach(func() {
			configurator = &mocks.Configurator{}
			configurator.On("Configure", mock.Anything, mock.Anything)
			clock = fakeclock.NewFakeClock(time.Now())
			storage = raft.NewMemoryStorage()
			observeC = make(chan uint64, 1)
			opts = etcdraft.Options{
				RaftID:          1,
				Clock:           clock,
				TickInterval:    interval,
				ElectionTick:    ELECTION_TICK,
				HeartbeatTick:   HEARTBEAT_TICK,
				MaxSizePerMsg:   1024 * 1024,
				MaxInflightMsgs: 256,
				Peers:           []raft.Peer{{ID: 1}},
				Logger:          logger,
				Storage:         storage,
			}
			support = &consensusmocks.FakeConsenterSupport{}
			support.ChainIDReturns(channelID)
			consenterMetadata = createMetadata(3, tlsCA)
			support.SharedConfigReturns(&mockconfig.Orderer{
				BatchTimeoutVal:      time.Hour,
				ConsensusMetadataVal: marshalOrPanic(consenterMetadata),
			})
			cutter = mockblockcutter.NewReceiver()
			support.BlockCutterReturns(cutter)

			var err error
			chain, err = etcdraft.NewChain(support, opts, configurator, nil, observeC)
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			chain.Start()

			
			
			
			
			
			
			
			
			
			
			
			
			Eventually(func() error {
				_, err := storage.Entries(1, 1, 1)
				return err
			}).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			chain.Halt()
		})

		Context("when a node starts up", func() {
			It("properly configures the communication layer", func() {
				expectedNodeConfig := nodeConfigFromMetadata(consenterMetadata)
				configurator.AssertCalled(testingInstance, "Configure", channelID, expectedNodeConfig)
			})
		})

		Context("when no Raft leader is elected", func() {
			It("fails to order envelope", func() {
				err := chain.Order(env, 0)
				Expect(err).To(MatchError("no Raft leader"))
			})
		})

		Context("when Raft leader is elected", func() {
			JustBeforeEach(func() {
				Eventually(func() bool {
					clock.Increment(interval)
					select {
					case <-observeC:
						return true
					default:
						return false
					}
				}).Should(BeTrue())
			})

			It("fails to order envelope if chain is halted", func() {
				chain.Halt()
				err := chain.Order(env, 0)
				Expect(err).To(MatchError("chain is stopped"))
			})

			It("produces blocks following batch rules", func() {
				close(cutter.Block)
				support.CreateNextBlockReturns(normalBlock)

				By("cutting next batch directly")
				cutter.CutNext = true
				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(support.WriteBlockCallCount).Should(Equal(1))

				By("respecting batch timeout")
				cutter.CutNext = false
				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())

				clock.WaitForNWatchersAndIncrement(timeout, 2)
				Eventually(support.WriteBlockCallCount).Should(Equal(2))
			})

			It("does not reset timer for every envelope", func() {
				close(cutter.Block)
				support.CreateNextBlockReturns(normalBlock)

				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(1))

				clock.WaitForNWatchersAndIncrement(timeout/2, 2)

				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(2))

				
				
				clock.Increment(timeout / 2)
				Eventually(support.WriteBlockCallCount).Should(Equal(1))
			})

			It("does not write a block if halted before timeout", func() {
				close(cutter.Block)
				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(1))

				
				Eventually(clock.WatcherCount).Should(Equal(2))

				chain.Halt()
				Consistently(support.WriteBlockCallCount).Should(Equal(0))
			})

			It("stops the timer if a batch is cut", func() {
				close(cutter.Block)
				support.CreateNextBlockReturns(normalBlock)

				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(1))

				clock.WaitForNWatchersAndIncrement(timeout/2, 2)

				By("force a batch to be cut before timer expires")
				cutter.CutNext = true
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(support.WriteBlockCallCount).Should(Equal(1))
				Expect(support.CreateNextBlockArgsForCall(0)).To(HaveLen(2))
				Expect(cutter.CurBatch()).To(HaveLen(0))

				
				cutter.CutNext = false
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(1))

				clock.WaitForNWatchersAndIncrement(timeout/2, 2)
				Consistently(support.WriteBlockCallCount).Should(Equal(1))

				clock.Increment(timeout / 2)
				Eventually(support.WriteBlockCallCount).Should(Equal(2))
				Expect(support.CreateNextBlockArgsForCall(1)).To(HaveLen(1))
			})

			It("cut two batches if incoming envelope does not fit into first batch", func() {
				close(cutter.Block)
				support.CreateNextBlockReturns(normalBlock)

				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(1))

				cutter.IsolatedTx = true
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())

				Eventually(support.CreateNextBlockCallCount).Should(Equal(2))
				Eventually(support.WriteBlockCallCount).Should(Equal(2))
			})

			Context("revalidation", func() {
				BeforeEach(func() {
					close(cutter.Block)
					support.CreateNextBlockReturns(normalBlock)

					timeout := time.Hour
					support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})
					support.SequenceReturns(1)
				})

				It("enqueue if envelope is still valid", func() {
					support.ProcessNormalMsgReturns(1, nil)

					err := chain.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
					Eventually(cutter.CurBatch).Should(HaveLen(1))
				})

				It("does not enqueue if envelope is not valid", func() {
					support.ProcessNormalMsgReturns(1, errors.Errorf("Envelope is invalid"))

					err := chain.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
					Consistently(cutter.CurBatch).Should(HaveLen(0))
				})
			})

			It("unblocks Errored if chain is halted", func() {
				Expect(chain.Errored()).NotTo(Receive())
				chain.Halt()
				Expect(chain.Errored()).Should(BeClosed())
			})

			Describe("Config updates", func() {
				var (
					configEnv   *common.Envelope
					configSeq   uint64
					configBlock *common.Block
				)

				
				newConfigEnv := func(chainID string, headerType common.HeaderType, configUpdateEnv *common.ConfigUpdateEnvelope) *common.Envelope {
					return &common.Envelope{
						Payload: marshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: marshalOrPanic(&common.ChannelHeader{
									Type:      int32(headerType),
									ChannelId: chainID,
								}),
							},
							Data: marshalOrPanic(&common.ConfigEnvelope{
								LastUpdate: &common.Envelope{
									Payload: marshalOrPanic(&common.Payload{
										Header: &common.Header{
											ChannelHeader: marshalOrPanic(&common.ChannelHeader{
												Type:      int32(common.HeaderType_CONFIG_UPDATE),
												ChannelId: chainID,
											}),
										},
										Data: marshalOrPanic(configUpdateEnv),
									}), 
								}, 
							}),
						}),
					}
				}

				newConfigUpdateEnv := func(chainID string, values map[string]*common.ConfigValue) *common.ConfigUpdateEnvelope {
					return &common.ConfigUpdateEnvelope{
						ConfigUpdate: marshalOrPanic(&common.ConfigUpdate{
							ChannelId: chainID,
							ReadSet:   &common.ConfigGroup{},
							WriteSet: &common.ConfigGroup{
								Groups: map[string]*common.ConfigGroup{
									"Orderer": {
										Values: values,
									},
								},
							}, 
						}),
					}
				}

				
				JustBeforeEach(func() {
					configBlock = &common.Block{
						Data: &common.BlockData{
							Data: [][]byte{marshalOrPanic(configEnv)},
						},
					}
					support.CreateNextBlockReturns(configBlock)
				})

				Context("when a config update with invalid header comes", func() {

					BeforeEach(func() {
						configEnv = newConfigEnv(channelID,
							common.HeaderType_CONFIG_UPDATE, 
							&common.ConfigUpdateEnvelope{ConfigUpdate: []byte("test invalid envelope")})
						configSeq = 0
					})

					It("should throw an error", func() {
						err := chain.Configure(configEnv, configSeq)
						Expect(err).To(MatchError("config transaction has unknown header type"))
					})
				})

				Context("when a type A config update comes", func() {

					Context("for existing channel", func() {

						
						BeforeEach(func() {
							values := map[string]*common.ConfigValue{
								"BatchTimeout": {
									Version: 1,
									Value: marshalOrPanic(&orderer.BatchTimeout{
										Timeout: "3ms",
									}),
								},
							}
							configEnv = newConfigEnv(channelID,
								common.HeaderType_CONFIG,
								newConfigUpdateEnv(channelID, values),
							)
							configSeq = 0
						}) 

						Context("without revalidation (i.e. correct config sequence)", func() {

							Context("without pending normal envelope", func() {
								It("should create a config block and no normal block", func() {
									err := chain.Configure(configEnv, configSeq)
									Expect(err).NotTo(HaveOccurred())
									Eventually(support.WriteConfigBlockCallCount).Should(Equal(1))
									Eventually(support.WriteBlockCallCount).Should(Equal(0))
								})
							})

							Context("with pending normal envelope", func() {
								It("should create a normal block and a config block", func() {
									
									close(cutter.Block)
									support.CreateNextBlockReturnsOnCall(0, normalBlock)
									support.CreateNextBlockReturnsOnCall(1, configBlock)

									By("adding a normal envelope")
									err := chain.Order(env, 0)
									Expect(err).NotTo(HaveOccurred())
									Eventually(cutter.CurBatch).Should(HaveLen(1))

									

									By("adding a config envelope")
									err = chain.Configure(configEnv, configSeq)
									Expect(err).NotTo(HaveOccurred())

									Eventually(support.CreateNextBlockCallCount).Should(Equal(2))
									Eventually(support.WriteBlockCallCount).Should(Equal(1))
									Eventually(support.WriteConfigBlockCallCount).Should(Equal(1))
								})
							})
						})

						Context("with revalidation (i.e. incorrect config sequence)", func() {

							BeforeEach(func() {
								support.SequenceReturns(1) 
							})

							It("should create config block upon correct revalidation", func() {
								support.ProcessConfigMsgReturns(configEnv, 1, nil) 

								err := chain.Configure(configEnv, configSeq)
								Expect(err).NotTo(HaveOccurred())
								Eventually(support.WriteConfigBlockCallCount).Should(Equal(1))
							})

							It("should not create config block upon incorrect revalidation", func() {
								support.ProcessConfigMsgReturns(configEnv, 1, errors.Errorf("Invalid config envelope at changed config sequence"))

								err := chain.Configure(configEnv, configSeq)
								Expect(err).NotTo(HaveOccurred())
								Consistently(support.WriteConfigBlockCallCount).Should(Equal(0)) 
							})
						})
					})

					Context("for creating a new channel", func() {

						
						BeforeEach(func() {
							chainID := "mychannel"
							configEnv = newConfigEnv(chainID,
								common.HeaderType_ORDERER_TRANSACTION,
								&common.ConfigUpdateEnvelope{ConfigUpdate: []byte("test channel creation envelope")})
							configSeq = 0
						}) 

						It("should be able to create a channel", func() {
							err := chain.Configure(configEnv, configSeq)
							Expect(err).NotTo(HaveOccurred())
						})
					})
				}) 

				Context("when a type B config update comes", func() {

					Context("for existing channel", func() {
						
						BeforeEach(func() {
							values := map[string]*common.ConfigValue{
								"ConsensusType": {
									Version: 1,
									Value: marshalOrPanic(&orderer.ConsensusType{
										Metadata: marshalOrPanic(consenterMetadata),
									}),
								},
							}
							configEnv = newConfigEnv(channelID,
								common.HeaderType_CONFIG,
								newConfigUpdateEnv(channelID, values))
							configSeq = 0
						}) 

						It("should be able to process config update of type B", func() {
							err := chain.Configure(configEnv, configSeq)
							Expect(err).NotTo(HaveOccurred())
						})
					})

					Context("updating consenters set", func() {
						
						BeforeEach(func() {
							values := map[string]*common.ConfigValue{
								"ConsensusType": {
									Version: 1,
									Value: marshalOrPanic(&orderer.ConsensusType{
										Metadata: marshalOrPanic(createMetadata(3, tlsCA)),
									}),
								},
							}
							configEnv = newConfigEnv(channelID,
								common.HeaderType_CONFIG,
								newConfigUpdateEnv(channelID, values))
							configSeq = 0
						}) 

						It("should fail, since consenters set change is not supported yet", func() {
							err := chain.Configure(configEnv, configSeq)
							Expect(err).To(MatchError("update of consenters set is not supported yet"))
						})
					})
				})
			})
		})
	})

	Describe("Multiple Raft nodes", func() {
		var (
			network    *network
			channelID  string
			timeout    time.Duration
			c1, c2, c3 *chain
		)

		BeforeEach(func() {
			channelID = "multi-node-channel"
			timeout = 10 * time.Second

			network = createNetwork(timeout, channelID, []uint64{1, 2, 3})
			c1 = network.chains[1]
			c2 = network.chains[2]
			c3 = network.chains[3]
		})

		When("2/3 nodes are running", func() {
			It("late node can catch up", func() {
				network.start(1, 2)
				network.elect(1)

				c1.cutter.CutNext = true
				err := c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() int { return c1.support.WriteBlockCallCount() }).Should(Equal(1))
				Eventually(func() int { return c2.support.WriteBlockCallCount() }).Should(Equal(1))
				Eventually(func() int { return c3.support.WriteBlockCallCount() }).Should(Equal(0))

				network.start(3)

				c1.clock.Increment(interval)
				Eventually(func() int { return c3.support.WriteBlockCallCount() }).Should(Equal(1))

				network.stop()
			})
		})

		When("3/3 nodes are running", func() {
			JustBeforeEach(func() {
				network.start()
				network.elect(1)
			})

			AfterEach(func() {
				network.stop()
			})

			It("orders envelope on leader", func() {
				By("instructed to cut next block")
				c1.cutter.CutNext = true
				err := c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }).Should(Equal(1))
					})

				By("respect batch timeout")
				c1.cutter.CutNext = false
				c1.support.CreateNextBlockReturns(normalBlock)

				err = c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())
				Eventually(c1.cutter.CurBatch).Should(HaveLen(1))

				c1.clock.WaitForNWatchersAndIncrement(timeout, 2)
				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }).Should(Equal(2))
					})
			})

			It("orders envelope on follower", func() {
				By("instructed to cut next block")
				c1.cutter.CutNext = true
				err := c2.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }).Should(Equal(1))
					})

				By("respect batch timeout")
				c1.cutter.CutNext = false
				c1.support.CreateNextBlockReturns(normalBlock)

				err = c2.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())
				Eventually(c1.cutter.CurBatch).Should(HaveLen(1))

				c1.clock.WaitForNWatchersAndIncrement(timeout, 2)
				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }).Should(Equal(2))
					})
			})

			Context("failover", func() {
				It("follower should step up as leader upon failover", func() {
					network.stop(1)
					network.elect(2)

					By("order envelope on new leader")
					c2.cutter.CutNext = true
					err := c2.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					Eventually(c1.support.WriteBlockCallCount).Should(Equal(0))

					
					Eventually(c2.support.WriteBlockCallCount).Should(Equal(1))
					Eventually(c3.support.WriteBlockCallCount).Should(Equal(1))

					By("order envelope on follower")
					err = c3.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					Eventually(c1.support.WriteBlockCallCount).Should(Equal(0))

					
					Eventually(c2.support.WriteBlockCallCount).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount).Should(Equal(2))
				})

				It("purges blockcutter and stops timer if leadership is lost", func() {
					
					err := c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())
					Eventually(c1.cutter.CurBatch).Should(HaveLen(1))

					network.disconnect(1)
					network.elect(2)
					network.connect(1)

					
					c2.clock.Increment(interval)                             
					Eventually(c1.observe).Should(Receive(Equal(uint64(2)))) 
					Expect(c1.clock.WatcherCount()).To(Equal(1))             

					Eventually(c1.cutter.CurBatch).Should(HaveLen(0))

					network.disconnect(2)
					n := network.elect(1) 

					err = c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					

					c1.clock.WaitForNWatchersAndIncrement(timeout-time.Duration(n*int(interval/time.Millisecond)), 2)
					Eventually(func() int { return c1.support.WriteBlockCallCount() }).Should(Equal(0))
					Eventually(func() int { return c3.support.WriteBlockCallCount() }).Should(Equal(0))

					c1.clock.Increment(3 * interval)
					Eventually(func() int { return c1.support.WriteBlockCallCount() }).Should(Equal(1))
					Eventually(func() int { return c3.support.WriteBlockCallCount() }).Should(Equal(1))
				})

				It("stale leader should not be able to propose block because of lagged term", func() {
					network.disconnect(1)
					network.elect(2)
					network.connect(1)

					c1.cutter.CutNext = true
					err := c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					network.exec(
						func(c *chain) {
							Consistently(c.support.WriteBlockCallCount).Should(Equal(0))
						})
				})

				It("aborts waiting for block to be committed upon leadership lost", func() {
					network.disconnect(1)

					c1.cutter.CutNext = true
					err := c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					network.exec(
						func(c *chain) {
							Consistently(c.support.WriteBlockCallCount).Should(Equal(0))
						})

					network.elect(2)
					network.connect(1)

					c2.clock.Increment(interval)
					
					Eventually(c1.observe).Should(Receive(Equal(uint64(2))))
				})
			})
		})
	})
})

func nodeConfigFromMetadata(consenterMetadata *raftprotos.Metadata) []cluster.RemoteNode {
	var nodes []cluster.RemoteNode
	for i, consenter := range consenterMetadata.Consenters {
		
		if i == 0 {
			continue
		}
		serverDER, _ := pem.Decode(consenter.ServerTlsCert)
		clientDER, _ := pem.Decode(consenter.ClientTlsCert)
		node := cluster.RemoteNode{
			ID:            uint64(i + 1),
			Endpoint:      "localhost:7050",
			ServerTLSCert: serverDER.Bytes,
			ClientTLSCert: clientDER.Bytes,
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func createMetadata(nodeCount int, tlsCA tlsgen.CA) *raftprotos.Metadata {
	md := &raftprotos.Metadata{}
	for i := 0; i < nodeCount; i++ {
		md.Consenters = append(md.Consenters, &raftprotos.Consenter{
			Host:          "localhost",
			Port:          7050,
			ServerTlsCert: serverTLSCert(tlsCA),
			ClientTlsCert: clientTLSCert(tlsCA),
		})
	}
	return md
}

func serverTLSCert(tlsCA tlsgen.CA) []byte {
	cert, err := tlsCA.NewServerCertKeyPair("localhost")
	if err != nil {
		panic(err)
	}
	return cert.Cert
}

func clientTLSCert(tlsCA tlsgen.CA) []byte {
	cert, err := tlsCA.NewClientCertKeyPair()
	if err != nil {
		panic(err)
	}
	return cert.Cert
}



func marshalOrPanic(pb proto.Message) []byte {
	data, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return data
}


type chain struct {
	id uint64

	support      *consensusmocks.FakeConsenterSupport
	cutter       *mockblockcutter.Receiver
	configurator *mocks.Configurator
	rpc          *mocks.FakeRPC
	storage      *raft.MemoryStorage
	clock        *fakeclock.FakeClock

	observe   chan uint64
	unstarted chan struct{}

	*etcdraft.Chain
}

func newChain(timeout time.Duration, channel string, id uint64, all []uint64) *chain {
	rpc := &mocks.FakeRPC{}
	clock := fakeclock.NewFakeClock(time.Now())
	storage := raft.NewMemoryStorage()

	peers := []raft.Peer{}
	for _, i := range all {
		peers = append(peers, raft.Peer{ID: i})
	}

	opts := etcdraft.Options{
		RaftID:          uint64(id),
		Clock:           clock,
		TickInterval:    interval,
		ElectionTick:    ELECTION_TICK,
		HeartbeatTick:   HEARTBEAT_TICK,
		MaxSizePerMsg:   1024 * 1024,
		MaxInflightMsgs: 256,
		Peers:           peers,
		Logger:          flogging.NewFabricLogger(zap.NewNop()),
		Storage:         storage,
	}

	support := &consensusmocks.FakeConsenterSupport{}
	support.ChainIDReturns(channel)
	support.CreateNextBlockReturns(&common.Block{Data: &common.BlockData{Data: [][]byte{[]byte("foo")}}})
	support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

	cutter := mockblockcutter.NewReceiver()
	close(cutter.Block)
	support.BlockCutterReturns(cutter)

	
	
	
	observe := make(chan uint64, 2)

	configurator := &mocks.Configurator{}
	configurator.On("Configure", mock.Anything, mock.Anything)

	c, err := etcdraft.NewChain(support, opts, configurator, rpc, observe)
	Expect(err).NotTo(HaveOccurred())

	ch := make(chan struct{})
	close(ch)
	return &chain{
		id:        id,
		support:   support,
		cutter:    cutter,
		rpc:       rpc,
		storage:   storage,
		observe:   observe,
		clock:     clock,
		unstarted: ch,
		Chain:     c,
	}
}

type network struct {
	chains map[uint64]*chain

	
	
	
	
	
	connLock     sync.RWMutex
	connectivity map[uint64]chan struct{}
}

func createNetwork(timeout time.Duration, channel string, ids []uint64) *network {
	n := &network{
		chains:       make(map[uint64]*chain),
		connectivity: make(map[uint64]chan struct{}),
	}

	for _, i := range ids {
		n.connectivity[i] = make(chan struct{})

		c := newChain(timeout, channel, i, ids)

		c.rpc.StepStub = func(dest uint64, msg *orderer.StepRequest) (*orderer.StepResponse, error) {
			n.connLock.RLock()
			defer n.connLock.RUnlock()

			select {
			case <-n.connectivity[dest]:
			case <-n.connectivity[c.id]:
			default:
				go n.chains[dest].Step(msg, c.id)
			}

			return nil, nil
		}

		c.rpc.SendSubmitStub = func(dest uint64, msg *orderer.SubmitRequest) error {
			n.connLock.RLock()
			defer n.connLock.RUnlock()

			select {
			case <-n.connectivity[dest]:
			case <-n.connectivity[c.id]:
			default:
				go n.chains[dest].Submit(msg, c.id)
			}

			return nil
		}

		n.chains[i] = c
	}

	return n
}

func (n *network) start(ids ...uint64) {
	nodes := ids
	if len(nodes) == 0 {
		for i := range n.chains {
			nodes = append(nodes, i)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for _, i := range nodes {
		go func(id uint64) {
			n.chains[id].Start()
			n.chains[id].unstarted = nil

			
			
			
			
			
			
			
			
			
			
			
			
			Eventually(func() error {
				_, err := n.chains[id].storage.Entries(1, 1, 1)
				return err
			}).ShouldNot(HaveOccurred())

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func (n *network) stop(ids ...uint64) {
	nodes := ids
	if len(nodes) == 0 {
		for i := range n.chains {
			nodes = append(nodes, i)
		}
	}

	for _, c := range nodes {
		n.chains[c].Halt()
		<-n.chains[c].Errored()
	}
}

func (n *network) exec(f func(c *chain), ids ...uint64) {
	if len(ids) == 0 {
		for _, c := range n.chains {
			f(c)
		}

		return
	}

	for _, i := range ids {
		f(n.chains[i])
	}
}




func (n *network) elect(id uint64) (tick int) {
	
	
	
	
	

	
	
	
	
	
	t := 100 * time.Millisecond

	c := n.chains[id]

	
	for tick < ELECTION_TICK-1 {
		c.clock.Increment(interval)
		tick++
		Consistently(c.observe, t).ShouldNot(Receive())
	}

	
	
	var elected bool
	for tick < 2*ELECTION_TICK-1 {
		c.clock.Increment(interval)
		tick++

		select {
		case <-time.After(t):
			
		case n := <-c.observe: 
			if n == 0 {
				
				
				
				
				
				Eventually(c.observe).Should(Receive(Equal(id)))
			} else {
				
				
				Expect(n).To(Equal(id))
			}
			elected = true
			break
		}
	}

	Expect(elected).To(BeTrue())

	
	for _, c := range n.chains {
		if c.id == id {
			continue
		}

		select {
		case <-c.Errored(): 
		case <-n.connectivity[c.id]: 
		case <-c.unstarted: 
		default:
			Eventually(c.observe).Should(Receive(Equal(id)))
		}
	}

	return tick
}

func (n *network) disconnect(i uint64) {
	close(n.connectivity[i])
}

func (n *network) connect(i uint64) {
	n.connLock.Lock()
	defer n.connLock.Unlock()
	n.connectivity[i] = make(chan struct{})
}
