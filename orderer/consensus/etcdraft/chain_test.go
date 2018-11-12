/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft_test

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
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



func skipIfRoot() {
	u, err := user.Current()
	Expect(err).NotTo(HaveOccurred())
	if u.Uid == "0" {
		Skip("you are running test as root, there's no way to make files unreadable")
	}
}

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
			walDir            string
		)

		BeforeEach(func() {
			var err error

			configurator = &mocks.Configurator{}
			configurator.On("Configure", mock.Anything, mock.Anything)
			clock = fakeclock.NewFakeClock(time.Now())
			storage = raft.NewMemoryStorage()
			walDir, err = ioutil.TempDir("", "wal-")
			Expect(err).NotTo(HaveOccurred())
			observeC = make(chan uint64, 1)

			support = &consensusmocks.FakeConsenterSupport{}
			support.ChainIDReturns(channelID)
			consenterMetadata = createMetadata(1, tlsCA)
			support.SharedConfigReturns(&mockconfig.Orderer{
				BatchTimeoutVal:      time.Hour,
				ConsensusMetadataVal: marshalOrPanic(consenterMetadata),
			})
			cutter = mockblockcutter.NewReceiver()
			support.BlockCutterReturns(cutter)

			membership := &raftprotos.RaftMetadata{
				Consenters:      map[uint64]*raftprotos.Consenter{},
				NextConsenterID: 1,
			}

			for _, c := range consenterMetadata.Consenters {
				membership.Consenters[membership.NextConsenterID] = c
				membership.NextConsenterID++
			}

			opts = etcdraft.Options{
				RaftID:          1,
				Clock:           clock,
				TickInterval:    interval,
				ElectionTick:    ELECTION_TICK,
				HeartbeatTick:   HEARTBEAT_TICK,
				MaxSizePerMsg:   1024 * 1024,
				MaxInflightMsgs: 256,
				RaftMetadata:    membership,
				Logger:          logger,
				MemoryStorage:   storage,
				WALDir:          walDir,
			}

			chain, err = etcdraft.NewChain(support, opts, configurator, nil, observeC)
			Expect(err).NotTo(HaveOccurred())
		})

		campaign := func(clock *fakeclock.FakeClock, observeC <-chan uint64) {
			Eventually(func() bool {
				clock.Increment(interval)
				select {
				case <-observeC:
					return true
				default:
					return false
				}
			}).Should(BeTrue())
		}

		JustBeforeEach(func() {
			chain.Start()

			
			
			
			
			
			
			
			
			
			
			
			
			Eventually(func() error {
				_, err := storage.Entries(1, 1, 1)
				return err
			}).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			chain.Halt()
			os.RemoveAll(walDir)
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
				campaign(clock, observeC)
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

					Context("updating protocol values", func() {
						
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

			Describe("Crash Fault Tolerance", func() {
				Describe("when a chain is started with existing WAL", func() {
					var (
						m1      *raftprotos.RaftMetadata
						m2      *raftprotos.RaftMetadata
						newOpts etcdraft.Options
					)

					BeforeEach(func() {
						newOpts = opts                                  
						newOpts.MemoryStorage = raft.NewMemoryStorage() 
					})

					JustBeforeEach(func() {
						
						
						close(cutter.Block)
						cutter.CutNext = true
						support.CreateNextBlockReturns(normalBlock)

						
						err := chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount).Should(Equal(1))

						_, metadata := support.WriteBlockArgsForCall(0)
						m1 = &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m1)

						err = chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount).Should(Equal(2))

						_, metadata = support.WriteBlockArgsForCall(1)
						m2 = &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m2)

						chain.Halt()
					})

					It("replays blocks from committed entries", func() {
						c := newChain(10*time.Second, channelID, walDir, 0, 1, []uint64{1})
						c.Start()
						defer c.Halt()

						Eventually(c.support.WriteBlockCallCount).Should(Equal(2))

						_, metadata := c.support.WriteBlockArgsForCall(0)
						m := &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m)
						Expect(m.RaftIndex).To(Equal(m1.RaftIndex))

						_, metadata = c.support.WriteBlockArgsForCall(1)
						m = &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m)
						Expect(m.RaftIndex).To(Equal(m2.RaftIndex))

						
						campaign(c.clock, c.observe)

						c.cutter.CutNext = true
						c.support.CreateNextBlockReturns(normalBlock)

						err := c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount).Should(Equal(3))
					})

					It("only replays blocks after Applied index", func() {
						c := newChain(10*time.Second, channelID, walDir, m1.RaftIndex, 1, []uint64{1})
						c.Start()
						defer c.Halt()

						Eventually(c.support.WriteBlockCallCount).Should(Equal(1))

						_, metadata := c.support.WriteBlockArgsForCall(0)
						m := &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m)
						Expect(m.RaftIndex).To(Equal(m2.RaftIndex))

						
						campaign(c.clock, c.observe)

						c.cutter.CutNext = true
						c.support.CreateNextBlockReturns(normalBlock)

						err := c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount).Should(Equal(2))
					})

					It("does not replay any block if already in sync", func() {
						c := newChain(10*time.Second, channelID, walDir, m2.RaftIndex, 1, []uint64{1})
						c.Start()
						defer c.Halt()

						Consistently(c.support.WriteBlockCallCount).Should(Equal(0))

						
						campaign(c.clock, c.observe)

						c.cutter.CutNext = true
						c.support.CreateNextBlockReturns(normalBlock)

						err := c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount).Should(Equal(1))
					})

					Context("WAL file is not readable", func() {
						It("fails to load wal", func() {
							skipIfRoot()

							files, err := ioutil.ReadDir(walDir)
							Expect(err).NotTo(HaveOccurred())
							for _, f := range files {
								os.Chmod(path.Join(walDir, f.Name()), 0300)
							}

							c, err := etcdraft.NewChain(support, opts, configurator, nil, observeC)
							Expect(c).To(BeNil())
							Expect(err).To(MatchError(ContainSubstring("failed to open existing WAL")))
						})
					})
				})
			})

			Context("Invalid WAL dir", func() {
				var support = &consensusmocks.FakeConsenterSupport{}

				When("WAL dir is a file", func() {
					It("replaces file with fresh WAL dir", func() {
						f, err := ioutil.TempFile("", "wal-")
						Expect(err).NotTo(HaveOccurred())
						defer os.RemoveAll(f.Name())

						chain, err := etcdraft.NewChain(
							support,
							etcdraft.Options{
								WALDir:        f.Name(),
								Logger:        logger,
								MemoryStorage: storage,
								RaftMetadata:  &raftprotos.RaftMetadata{},
							},
							configurator,
							nil,
							observeC)
						Expect(chain).NotTo(BeNil())
						Expect(err).NotTo(HaveOccurred())

						info, err := os.Stat(f.Name())
						Expect(err).NotTo(HaveOccurred())
						Expect(info.IsDir()).To(BeTrue())
					})
				})

				When("WAL dir is not writeable", func() {
					It("replace it with fresh WAL dir", func() {
						d, err := ioutil.TempDir("", "wal-")
						Expect(err).NotTo(HaveOccurred())
						defer os.RemoveAll(d)

						err = os.Chmod(d, 0500)
						Expect(err).NotTo(HaveOccurred())

						chain, err := etcdraft.NewChain(
							support,
							etcdraft.Options{
								WALDir:        d,
								Logger:        logger,
								MemoryStorage: storage,
								RaftMetadata:  &raftprotos.RaftMetadata{},
							},
							nil,
							nil,
							nil)
						Expect(chain).NotTo(BeNil())
						Expect(err).ToNot(HaveOccurred())
					})
				})

				When("WAL parent dir is not writeable", func() {
					It("fails to bootstrap fresh raft node", func() {
						skipIfRoot()

						d, err := ioutil.TempDir("", "wal-")
						Expect(err).NotTo(HaveOccurred())
						defer os.RemoveAll(d)

						err = os.Chmod(d, 0500)
						Expect(err).NotTo(HaveOccurred())

						chain, err := etcdraft.NewChain(
							support,
							etcdraft.Options{
								WALDir:       path.Join(d, "wal-dir"),
								Logger:       logger,
								RaftMetadata: &raftprotos.RaftMetadata{},
							},
							nil,
							nil,
							nil)
						Expect(chain).To(BeNil())
						Expect(err).To(MatchError(ContainSubstring("failed to initialize WAL: mkdir")))
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
			dataDir    string
			c1, c2, c3 *chain
		)

		BeforeEach(func() {
			var err error

			channelID = "multi-node-channel"
			timeout = 10 * time.Second

			dataDir, err = ioutil.TempDir("", "raft-test-")
			Expect(err).NotTo(HaveOccurred())

			network = createNetwork(timeout, channelID, dataDir, []uint64{1, 2, 3})
			c1 = network.chains[1]
			c2 = network.chains[2]
			c3 = network.chains[3]
		})

		AfterEach(func() {
			os.RemoveAll(dataDir)
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

				It("follower cannot be elected if its log is not up-to-date", func() {
					network.disconnect(2)

					c1.cutter.CutNext = true
					err := c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					Eventually(c1.support.WriteBlockCallCount).Should(Equal(1))
					Eventually(c2.support.WriteBlockCallCount).Should(Equal(0))
					Eventually(c3.support.WriteBlockCallCount).Should(Equal(1))

					network.disconnect(1)
					network.connect(2)

					
					for tick := 0; tick < 2*ELECTION_TICK-1; tick++ {
						c2.clock.Increment(interval)
						Consistently(c2.observe).ShouldNot(Receive(Equal(2)))
					}

					Eventually(c3.observe).Should(Receive(Equal(uint64(0))))
					network.elect(3) 
				})

				It("follower can catch up and then campaign with success", func() {
					network.disconnect(2)

					c1.cutter.CutNext = true
					for i := 0; i < 10; i++ {
						err := c1.Order(env, 0)
						Expect(err).NotTo(HaveOccurred())
					}

					Eventually(c1.support.WriteBlockCallCount).Should(Equal(10))
					Eventually(c2.support.WriteBlockCallCount).Should(Equal(0))
					Eventually(c3.support.WriteBlockCallCount).Should(Equal(10))

					network.rejoin(2, false)
					Eventually(c2.support.WriteBlockCallCount).Should(Equal(10))

					network.disconnect(1)
					network.elect(2)
				})

				It("purges blockcutter and stops timer if leadership is lost", func() {
					
					err := c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())
					Eventually(c1.cutter.CurBatch).Should(HaveLen(1))

					network.disconnect(1)
					network.elect(2)
					network.rejoin(1, true)

					Eventually(c1.clock.WatcherCount).Should(Equal(1)) 

					Expect(c1.clock.WatcherCount()).To(Equal(1)) 
					Eventually(c1.cutter.CurBatch).Should(HaveLen(0))

					network.disconnect(2)
					n := network.elect(1) 

					err = c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					

					c1.clock.WaitForNWatchersAndIncrement(timeout-time.Duration(n*int(interval/time.Millisecond)), 2)
					Eventually(func() int { return c1.support.WriteBlockCallCount() }).Should(Equal(0))
					Eventually(func() int { return c3.support.WriteBlockCallCount() }).Should(Equal(0))

					c1.clock.Increment(time.Duration(n * int(interval/time.Millisecond)))
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
	walDir       string
	clock        *fakeclock.FakeClock
	opts         etcdraft.Options

	observe   chan uint64
	unstarted chan struct{}

	*etcdraft.Chain
}

func newChain(timeout time.Duration, channel string, walDir string, applied uint64, id uint64, all []uint64) *chain {
	rpc := &mocks.FakeRPC{}
	clock := fakeclock.NewFakeClock(time.Now())
	storage := raft.NewMemoryStorage()
	tlsCA, _ := tlsgen.NewCA()

	meta := &raftprotos.RaftMetadata{
		Consenters:      map[uint64]*raftprotos.Consenter{},
		NextConsenterID: 1,
		RaftIndex:       applied,
	}

	for _, raftID := range all {
		meta.Consenters[uint64(raftID)] = &raftprotos.Consenter{
			Host:          "localhost",
			Port:          7051,
			ClientTlsCert: clientTLSCert(tlsCA),
			ServerTlsCert: serverTLSCert(tlsCA),
		}
		if uint64(raftID) > meta.NextConsenterID {
			meta.NextConsenterID = uint64(raftID)
		}
	}
	meta.NextConsenterID++

	opts := etcdraft.Options{
		RaftID:          uint64(id),
		Clock:           clock,
		TickInterval:    interval,
		ElectionTick:    ELECTION_TICK,
		HeartbeatTick:   HEARTBEAT_TICK,
		MaxSizePerMsg:   1024 * 1024,
		MaxInflightMsgs: 256,
		RaftMetadata:    meta,
		Logger:          flogging.NewFabricLogger(zap.NewNop()),
		MemoryStorage:   storage,
		WALDir:          walDir,
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
		opts:      opts,
		Chain:     c,
	}
}

type network struct {
	leader uint64
	chains map[uint64]*chain

	
	
	
	
	
	connLock     sync.RWMutex
	connectivity map[uint64]chan struct{}
}

func createNetwork(timeout time.Duration, channel string, dataDir string, ids []uint64) *network {
	n := &network{
		chains:       make(map[uint64]*chain),
		connectivity: make(map[uint64]chan struct{}),
	}

	for _, i := range ids {
		n.connectivity[i] = make(chan struct{})

		dir, err := ioutil.TempDir(dataDir, fmt.Sprintf("node-%d-", i))
		Expect(err).NotTo(HaveOccurred())

		c := newChain(timeout, channel, dir, 0, i, ids)

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
			defer GinkgoRecover()
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



func (n *network) rejoin(id uint64, wasLeader bool) {
	n.connect(id)
	n.chains[n.leader].clock.Increment(interval)

	if wasLeader {
		Eventually(n.chains[id].observe).Should(Receive(Equal(n.leader)))
	} else {
		Consistently(n.chains[id].observe).ShouldNot(Receive())
	}

	
	i, err := n.chains[n.leader].opts.MemoryStorage.LastIndex()
	Expect(err).NotTo(HaveOccurred())
	Eventually(n.chains[id].opts.MemoryStorage.LastIndex).Should(Equal(i))
}




func (n *network) elect(id uint64) (tick int) {
	
	
	
	
	
	t := 1000 * time.Millisecond

	c := n.chains[id]

	var elected bool
	for !elected {
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

	n.leader = id
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
