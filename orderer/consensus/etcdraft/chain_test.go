/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft_test

import (
	"bytes"
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
	"github.com/mcc-github/blockchain/bccsp/factory"
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
	interval            = time.Second
	LongEventualTimeout = 5 * time.Second
	ELECTION_TICK       = 2
	HEARTBEAT_TICK      = 1
)

func init() {
	factory.InitFactories(nil)
}



func skipIfRoot() {
	u, err := user.Current()
	Expect(err).NotTo(HaveOccurred())
	if u.Uid == "0" {
		Skip("you are running test as root, there's no way to make files unreadable")
	}
}

var _ = Describe("Chain", func() {
	var (
		env       *common.Envelope
		channelID string
		tlsCA     tlsgen.CA
		logger    *flogging.FabricLogger
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
			dataDir           string
			walDir            string
			snapDir           string
			err               error
		)

		BeforeEach(func() {
			configurator = &mocks.Configurator{}
			configurator.On("Configure", mock.Anything, mock.Anything)
			clock = fakeclock.NewFakeClock(time.Now())
			storage = raft.NewMemoryStorage()

			dataDir, err = ioutil.TempDir("", "wal-")
			Expect(err).NotTo(HaveOccurred())
			walDir = path.Join(dataDir, "wal")
			snapDir = path.Join(dataDir, "snapshot")

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

			
			support.HeightReturns(1)
			support.BlockReturns(getSeedBlock())

			meta := &raftprotos.RaftMetadata{
				Consenters:      map[uint64]*raftprotos.Consenter{},
				NextConsenterId: 1,
			}

			for _, c := range consenterMetadata.Consenters {
				meta.Consenters[meta.NextConsenterId] = c
				meta.NextConsenterId++
			}

			opts = etcdraft.Options{
				RaftID:          1,
				Clock:           clock,
				TickInterval:    interval,
				ElectionTick:    ELECTION_TICK,
				HeartbeatTick:   HEARTBEAT_TICK,
				MaxSizePerMsg:   1024 * 1024,
				MaxInflightMsgs: 256,
				RaftMetadata:    meta,
				Logger:          logger,
				MemoryStorage:   storage,
				WALDir:          walDir,
				SnapDir:         snapDir,
			}
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
			}, LongEventualTimeout).Should(BeTrue())
		}

		JustBeforeEach(func() {
			chain, err = etcdraft.NewChain(support, opts, configurator, nil, nil, observeC)
			Expect(err).NotTo(HaveOccurred())

			chain.Start()

			
			
			
			
			
			
			
			
			
			
			
			
			Eventually(func() error {
				_, err := storage.Entries(1, 1, 1)
				return err
			}, LongEventualTimeout).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			chain.Halt()
			os.RemoveAll(dataDir)
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
				Expect(err).To(MatchError("node bootstrapping has not finished"))
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

				By("cutting next batch directly")
				cutter.CutNext = true
				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

				By("respecting batch timeout")
				cutter.CutNext = false
				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())

				clock.WaitForNWatchersAndIncrement(timeout, 2)
				Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
			})

			It("does not reset timer for every envelope", func() {
				close(cutter.Block)

				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				clock.WaitForNWatchersAndIncrement(timeout/2, 2)

				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(2))

				
				
				clock.Increment(timeout / 2)
				Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
			})

			It("does not write a block if halted before timeout", func() {
				close(cutter.Block)
				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				
				Eventually(clock.WatcherCount, LongEventualTimeout).Should(Equal(2))

				chain.Halt()
				Consistently(support.WriteBlockCallCount).Should(Equal(0))
			})

			It("stops the timer if a batch is cut", func() {
				close(cutter.Block)

				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				clock.WaitForNWatchersAndIncrement(timeout/2, 2)

				By("force a batch to be cut before timer expires")
				cutter.CutNext = true
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())

				Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
				b, _ := support.WriteBlockArgsForCall(0)
				Expect(b.Data.Data).To(HaveLen(2))
				Expect(cutter.CurBatch()).To(HaveLen(0))

				
				cutter.CutNext = false
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				clock.WaitForNWatchersAndIncrement(timeout/2, 2)
				Consistently(support.WriteBlockCallCount).Should(Equal(1))

				clock.Increment(timeout / 2)

				Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				b, _ = support.WriteBlockArgsForCall(1)
				Expect(b.Data.Data).To(HaveLen(1))
			})

			It("cut two batches if incoming envelope does not fit into first batch", func() {
				close(cutter.Block)

				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				cutter.IsolatedTx = true
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())

				Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
			})

			Context("revalidation", func() {
				BeforeEach(func() {
					close(cutter.Block)

					timeout := time.Hour
					support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})
					support.SequenceReturns(1)
				})

				It("enqueue if envelope is still valid", func() {
					support.ProcessNormalMsgReturns(1, nil)

					err := chain.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
					Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))
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
					configEnv *common.Envelope
					configSeq uint64
				)

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
									Eventually(support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
									Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))
								})
							})

							Context("with pending normal envelope", func() {
								It("should create a normal block and a config block", func() {
									
									close(cutter.Block)

									By("adding a normal envelope")
									err := chain.Order(env, 0)
									Expect(err).NotTo(HaveOccurred())
									Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

									

									By("adding a config envelope")
									err = chain.Configure(configEnv, configSeq)
									Expect(err).NotTo(HaveOccurred())

									Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
									Eventually(support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
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
								Eventually(support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
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

					Context("updating consenters set by more than one node", func() {
						
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

						It("should fail, since consenters set change is not supported", func() {
							err := chain.Configure(configEnv, configSeq)
							Expect(err).To(MatchError("update of more than one consenters at a time is not supported"))
						})
					})

					Context("updating consenters set by exactly one node", func() {
						It("should be able to process config update adding single node", func() {
							metadata := proto.Clone(consenterMetadata).(*raftprotos.Metadata)
							metadata.Consenters = append(metadata.Consenters, &raftprotos.Consenter{
								Host:          "localhost",
								Port:          7050,
								ServerTlsCert: serverTLSCert(tlsCA),
								ClientTlsCert: clientTLSCert(tlsCA),
							})

							values := map[string]*common.ConfigValue{
								"ConsensusType": {
									Version: 1,
									Value: marshalOrPanic(&orderer.ConsensusType{
										Metadata: marshalOrPanic(metadata),
									}),
								},
							}
							configEnv = newConfigEnv(channelID,
								common.HeaderType_CONFIG,
								newConfigUpdateEnv(channelID, values))
							configSeq = 0

							err := chain.Configure(configEnv, configSeq)
							Expect(err).NotTo(HaveOccurred())
						})

						It("should be able to process config update removing single node", func() {
							metadata := proto.Clone(consenterMetadata).(*raftprotos.Metadata)
							
							metadata.Consenters = metadata.Consenters[1:]
							values := map[string]*common.ConfigValue{
								"ConsensusType": {
									Version: 1,
									Value: marshalOrPanic(&orderer.ConsensusType{
										Metadata: marshalOrPanic(metadata),
									}),
								},
							}
							configEnv = newConfigEnv(channelID,
								common.HeaderType_CONFIG,
								newConfigUpdateEnv(channelID, values))
							configSeq = 0

							err := chain.Configure(configEnv, configSeq)
							Expect(err).NotTo(HaveOccurred())
						})

						It("fail since not allowed to add and remove node at same change", func() {
							metadata := proto.Clone(consenterMetadata).(*raftprotos.Metadata)
							
							metadata.Consenters = metadata.Consenters[1:]
							metadata.Consenters = append(metadata.Consenters, &raftprotos.Consenter{
								Host:          "localhost",
								Port:          7050,
								ServerTlsCert: serverTLSCert(tlsCA),
								ClientTlsCert: clientTLSCert(tlsCA),
							})
							values := map[string]*common.ConfigValue{
								"ConsensusType": {
									Version: 1,
									Value: marshalOrPanic(&orderer.ConsensusType{
										Metadata: marshalOrPanic(metadata),
									}),
								},
							}
							configEnv = newConfigEnv(channelID,
								common.HeaderType_CONFIG,
								newConfigUpdateEnv(channelID, values))
							configSeq = 0

							err := chain.Configure(configEnv, configSeq)
							Expect(err).To(MatchError("update of more than one consenters at a time is not supported"))
						})
					})
				})
			})

			Describe("Crash Fault Tolerance", func() {
				var (
					raftMetadata *raftprotos.RaftMetadata
				)

				BeforeEach(func() {
					tlsCA, _ := tlsgen.NewCA()

					raftMetadata = &raftprotos.RaftMetadata{
						Consenters: map[uint64]*raftprotos.Consenter{
							1: {
								Host:          "localhost",
								Port:          7051,
								ClientTlsCert: clientTLSCert(tlsCA),
								ServerTlsCert: serverTLSCert(tlsCA),
							},
						},
						NextConsenterId: 2,
					}
				})

				Describe("when a chain is started with existing WAL", func() {
					var (
						m1 *raftprotos.RaftMetadata
						m2 *raftprotos.RaftMetadata
					)
					JustBeforeEach(func() {
						
						
						close(cutter.Block)
						cutter.CutNext = true

						
						err := chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

						_, metadata := support.WriteBlockArgsForCall(0)
						m1 = &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m1)

						err = chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

						_, metadata = support.WriteBlockArgsForCall(1)
						m2 = &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m2)

						chain.Halt()
					})

					It("replays blocks from committed entries", func() {
						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata)
						c.init()
						c.Start()
						defer c.Halt()

						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

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

						err := c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))

					})

					It("only replays blocks after Applied index", func() {
						raftMetadata.RaftIndex = m1.RaftIndex
						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata)
						c.init()
						c.Start()
						defer c.Halt()

						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

						_, metadata := c.support.WriteBlockArgsForCall(0)
						m := &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m)
						Expect(m.RaftIndex).To(Equal(m2.RaftIndex))

						
						campaign(c.clock, c.observe)

						c.cutter.CutNext = true

						err := c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					})

					It("does not replay any block if already in sync", func() {
						raftMetadata.RaftIndex = m2.RaftIndex
						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata)
						c.init()
						c.Start()
						defer c.Halt()

						Consistently(c.support.WriteBlockCallCount).Should(Equal(0))

						
						campaign(c.clock, c.observe)

						c.cutter.CutNext = true

						err := c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					})

					Context("WAL file is not readable", func() {
						It("fails to load wal", func() {
							skipIfRoot()

							files, err := ioutil.ReadDir(walDir)
							Expect(err).NotTo(HaveOccurred())
							for _, f := range files {
								os.Chmod(path.Join(walDir, f.Name()), 0300)
							}

							c, err := etcdraft.NewChain(support, opts, configurator, nil, nil, observeC)
							Expect(c).To(BeNil())
							Expect(err).To(MatchError(ContainSubstring("failed to open existing WAL")))
						})
					})
				})

				Describe("when snapshotting is enabled (snapshot interval is not zero)", func() {
					var (
						m *raftprotos.RaftMetadata

						ledgerLock sync.Mutex
						ledger     map[uint64]*common.Block
					)

					countFiles := func() int {
						files, err := ioutil.ReadDir(snapDir)
						Expect(err).NotTo(HaveOccurred())
						return len(files)
					}

					BeforeEach(func() {
						opts.SnapInterval = 2
						opts.SnapshotCatchUpEntries = 2

						close(cutter.Block)
						cutter.CutNext = true

						ledgerLock.Lock()
						ledger = map[uint64]*common.Block{
							0: getSeedBlock(), 
						}
						ledgerLock.Unlock()

						support.WriteBlockStub = func(b *common.Block, meta []byte) {
							bytes, err := proto.Marshal(&common.Metadata{Value: meta})
							Expect(err).NotTo(HaveOccurred())
							b.Metadata.Metadata[common.BlockMetadataIndex_ORDERER] = bytes

							ledgerLock.Lock()
							defer ledgerLock.Unlock()
							ledger[b.Header.Number] = b
						}

						support.HeightStub = func() uint64 {
							ledgerLock.Lock()
							defer ledgerLock.Unlock()
							return uint64(len(ledger))
						}
					})

					JustBeforeEach(func() {
						err = chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

						err = chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

						_, metadata := support.WriteBlockArgsForCall(1)
						m = &raftprotos.RaftMetadata{}
						proto.Unmarshal(metadata, m)
					})

					It("writes snapshot file to snapDir", func() {
						
						
						
						
						

						Eventually(countFiles, LongEventualTimeout).Should(Equal(1))
						Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", 1))

						
						err = chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))
					})

					It("pauses chain if sync is in progress", func() {
						
						
						
						

						
						Eventually(countFiles, LongEventualTimeout).Should(Equal(1))

						chain.Halt()

						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata)
						c.init()

						signal := make(chan struct{})

						c.puller.PullBlockStub = func(i uint64) *common.Block {
							<-signal 
							ledgerLock.Lock()
							defer ledgerLock.Unlock()
							if i >= uint64(len(ledger)) {
								return nil
							}

							return ledger[i]
						}

						err := c.WaitReady()
						Expect(err).To(MatchError("chain is not started"))

						c.Start()
						defer c.Halt()

						
						signal <- struct{}{}

						done := make(chan error)
						go func() {
							done <- c.WaitReady()
						}()

						Consistently(done).ShouldNot(Receive())
						close(signal) 

						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					})

					It("restores snapshot w/o extra entries", func() {
						
						
						
						
						
						

						Eventually(countFiles, LongEventualTimeout).Should(Equal(1))
						Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", 1))
						snapshot, err := opts.MemoryStorage.Snapshot() 
						Expect(err).NotTo(HaveOccurred())
						i, err := opts.MemoryStorage.FirstIndex() 
						Expect(err).NotTo(HaveOccurred())

						
						Expect(i).To(Equal(snapshot.Metadata.Index - opts.SnapshotCatchUpEntries + 1))

						chain.Halt()

						raftMetadata.RaftIndex = m.RaftIndex
						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata)

						c.init()
						c.Start()
						defer c.Halt()

						
						
						Eventually(c.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(Equal(snapshot.Metadata.Index + 1))
						Eventually(c.opts.MemoryStorage.LastIndex, LongEventualTimeout).Should(Equal(snapshot.Metadata.Index))

						
						Eventually(func() bool {
							c.clock.Increment(interval)
							select {
							case <-c.observe:
								return true
							default:
								return false
							}
						}, LongEventualTimeout).Should(BeTrue())

						c.cutter.CutNext = true
						err = c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					})

					It("restores snapshot w/ extra entries", func() {
						
						
						
						
						

						
						Eventually(countFiles, LongEventualTimeout).Should(Equal(1))
						Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", 1))
						snapshot, err := opts.MemoryStorage.Snapshot() 
						Expect(err).NotTo(HaveOccurred())
						i, err := opts.MemoryStorage.FirstIndex() 
						Expect(err).NotTo(HaveOccurred())

						
						Expect(i).To(Equal(snapshot.Metadata.Index - opts.SnapshotCatchUpEntries + 1))

						err = chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))

						lasti, _ := opts.MemoryStorage.LastIndex()

						chain.Halt()

						raftMetadata.RaftIndex = m.RaftIndex
						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata)
						c.support.HeightReturns(5)
						c.support.BlockReturns(&common.Block{
							Header:   &common.BlockHeader{},
							Data:     &common.BlockData{Data: [][]byte{[]byte("foo")}},
							Metadata: &common.BlockMetadata{Metadata: make([][]byte, 4)},
						})

						c.init()
						c.Start()
						defer c.Halt()

						Eventually(c.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(Equal(snapshot.Metadata.Index + 1))
						Eventually(c.opts.MemoryStorage.LastIndex, LongEventualTimeout).Should(Equal(lasti))

						
						Eventually(func() bool {
							c.clock.Increment(interval)
							select {
							case <-c.observe:
								return true
							default:
								return false
							}
						}, LongEventualTimeout).Should(BeTrue())

						c.cutter.CutNext = true
						err = c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					})

					When("local ledger is in sync with snapshot", func() {
						It("does not pull blocks and still respects snapshot interval", func() {
							
							
							
							
							
							
							
							
							
							

							
							Eventually(countFiles, LongEventualTimeout).Should(Equal(1))

							
							err = chain.Order(env, uint64(0))
							Expect(err).NotTo(HaveOccurred())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))

							chain.Halt()

							raftMetadata.RaftIndex = m.RaftIndex
							c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata)
							

							c.support.HeightReturns(3)
							c.support.BlockReturns(&common.Block{
								Header:   &common.BlockHeader{},
								Data:     &common.BlockData{Data: [][]byte{[]byte("foo")}},
								Metadata: &common.BlockMetadata{Metadata: make([][]byte, 4)},
							})
							c.opts.SnapInterval = 2

							c.init()
							c.Start()
							defer c.Halt()

							
							Eventually(func() bool {
								c.clock.Increment(interval)
								select {
								case <-c.observe:
									return true
								default:
									return false
								}
							}, LongEventualTimeout).Should(BeTrue())

							c.cutter.CutNext = true
							err = c.Order(env, uint64(0))
							Expect(err).NotTo(HaveOccurred())

							Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
							Expect(c.puller.PullBlockCallCount()).Should(BeZero())
							Eventually(countFiles, LongEventualTimeout).Should(Equal(2))
						})
					})
				})
			})

			Context("Invalid WAL dir", func() {
				var support = &consensusmocks.FakeConsenterSupport{}
				BeforeEach(func() {
					
					support.HeightReturns(1)
					support.BlockReturns(getSeedBlock())
				})

				When("WAL dir is a file", func() {
					It("replaces file with fresh WAL dir", func() {
						f, err := ioutil.TempFile("", "wal-")
						Expect(err).NotTo(HaveOccurred())
						defer os.RemoveAll(f.Name())

						chain, err := etcdraft.NewChain(
							support,
							etcdraft.Options{
								WALDir:        f.Name(),
								SnapDir:       snapDir,
								Logger:        logger,
								MemoryStorage: storage,
								RaftMetadata:  &raftprotos.RaftMetadata{},
							},
							configurator,
							nil,
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
								SnapDir:       snapDir,
								Logger:        logger,
								MemoryStorage: storage,
								RaftMetadata:  &raftprotos.RaftMetadata{},
							},
							nil,
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
								SnapDir:      snapDir,
								Logger:       logger,
								RaftMetadata: &raftprotos.RaftMetadata{},
							},
							nil,
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
			network      *network
			channelID    string
			timeout      time.Duration
			dataDir      string
			c1, c2, c3   *chain
			raftMetadata *raftprotos.RaftMetadata
		)

		BeforeEach(func() {
			var err error

			channelID = "multi-node-channel"
			timeout = 10 * time.Second

			dataDir, err = ioutil.TempDir("", "raft-test-")
			Expect(err).NotTo(HaveOccurred())

			raftMetadata = &raftprotos.RaftMetadata{
				Consenters: map[uint64]*raftprotos.Consenter{
					1: {
						Host:          "localhost",
						Port:          7051,
						ClientTlsCert: clientTLSCert(tlsCA),
						ServerTlsCert: serverTLSCert(tlsCA),
					},
					2: {
						Host:          "localhost",
						Port:          7051,
						ClientTlsCert: clientTLSCert(tlsCA),
						ServerTlsCert: serverTLSCert(tlsCA),
					},
					3: {
						Host:          "localhost",
						Port:          7051,
						ClientTlsCert: clientTLSCert(tlsCA),
						ServerTlsCert: serverTLSCert(tlsCA),
					},
				},
				NextConsenterId: 4,
			}

			network = createNetwork(timeout, channelID, dataDir, raftMetadata)
			c1 = network.chains[1]
			c2 = network.chains[2]
			c3 = network.chains[3]
		})

		AfterEach(func() {
			os.RemoveAll(dataDir)
		})

		When("2/3 nodes are running", func() {
			It("late node can catch up", func() {
				network.init()
				network.start(1, 2)
				network.elect(1)

				c1.cutter.CutNext = true
				err := c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() int { return c1.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
				Eventually(func() int { return c2.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
				Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))

				network.start(3)

				c1.clock.Increment(interval)
				Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))

				network.stop()
			})

			It("late node receives snapshot from leader", func() {
				c1.opts.SnapInterval = 1
				c1.opts.SnapshotCatchUpEntries = 1

				c1.cutter.CutNext = true

				var blocksLock sync.Mutex
				blocks := make(map[uint64]*common.Block) 

				c1.support.WriteBlockStub = func(b *common.Block, meta []byte) {
					blocksLock.Lock()
					defer blocksLock.Unlock()
					bytes, err := proto.Marshal(&common.Metadata{Value: meta})
					Expect(err).NotTo(HaveOccurred())
					b.Metadata.Metadata[common.BlockMetadataIndex_ORDERER] = bytes
					blocks[b.Header.Number] = b
				}

				c3.puller.PullBlockStub = func(i uint64) *common.Block {
					blocksLock.Lock()
					defer blocksLock.Unlock()
					b, exist := blocks[i]
					if !exist {
						return nil
					}

					return b
				}

				network.init()
				network.start(1, 2)
				network.elect(1)

				err := c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() int { return c1.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
				Eventually(func() int { return c2.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
				Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))

				err = c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() int { return c1.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(2))
				Eventually(func() int { return c2.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(2))
				Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))

				network.start(3)

				c1.clock.Increment(interval)
				Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(2))

				network.stop()
			})
		})

		When("reconfiguring raft cluster", func() {
			const (
				defaultTimeout = 5 * time.Second
			)
			var (
				addConsenterConfigValue = func() map[string]*common.ConfigValue {
					metadata := &raftprotos.Metadata{}
					for _, consenter := range raftMetadata.Consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					newConsenter := &raftprotos.Consenter{
						Host:          "localhost",
						Port:          7050,
						ServerTlsCert: serverTLSCert(tlsCA),
						ClientTlsCert: clientTLSCert(tlsCA),
					}
					metadata.Consenters = append(metadata.Consenters, newConsenter)

					return map[string]*common.ConfigValue{
						"ConsensusType": {
							Version: 1,
							Value: marshalOrPanic(&orderer.ConsensusType{
								Metadata: marshalOrPanic(metadata),
							}),
						},
					}
				}
				removeConsenterConfigValue = func(id uint64) map[string]*common.ConfigValue {
					newRaftMetadata := proto.Clone(raftMetadata).(*raftprotos.RaftMetadata)
					delete(newRaftMetadata.Consenters, id)

					metadata := &raftprotos.Metadata{}
					for _, consenter := range newRaftMetadata.Consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					return map[string]*common.ConfigValue{
						"ConsensusType": {
							Version: 1,
							Value: marshalOrPanic(&orderer.ConsensusType{
								Metadata: marshalOrPanic(metadata),
							}),
						},
					}
				}
			)

			BeforeEach(func() {
				network.init()
				network.start()
				network.elect(1)

				By("Submitting first tx to cut the block")
				c1.cutter.CutNext = true
				err := c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				c1.clock.Increment(interval)

				network.exec(
					func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, defaultTimeout).Should(Equal(1))
					})
			})

			AfterEach(func() {
				network.stop()
			})

			Context("reconfiguration", func() {
				It("trying to simultaneously add and remove nodes in one config update", func() {

					updatedRaftMetadata := proto.Clone(raftMetadata).(*raftprotos.RaftMetadata)
					
					delete(updatedRaftMetadata.Consenters, 2)

					metadata := &raftprotos.Metadata{}
					for _, consenter := range updatedRaftMetadata.Consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					
					newConsenter := &raftprotos.Consenter{
						Host:          "localhost",
						Port:          7050,
						ServerTlsCert: serverTLSCert(tlsCA),
						ClientTlsCert: clientTLSCert(tlsCA),
					}
					metadata.Consenters = append(metadata.Consenters, newConsenter)

					value := map[string]*common.ConfigValue{
						"ConsensusType": {
							Version: 1,
							Value: marshalOrPanic(&orderer.ConsensusType{
								Metadata: marshalOrPanic(metadata),
							}),
						},
					}

					By("creating new configuration with removed node and new one")
					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, value))
					c1.cutter.CutNext = true

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).To(MatchError("update of more than one consenters at a time is not supported"))
				})

				It("adding node to the cluster", func() {

					c4 := newChain(timeout, channelID, dataDir, 4, &raftprotos.RaftMetadata{
						Consenters: map[uint64]*raftprotos.Consenter{},
					})
					c4.init()

					By("adding new node to the network")
					Expect(c4.support.WriteBlockCallCount()).Should(Equal(0))
					Expect(c4.support.WriteConfigBlockCallCount()).Should(Equal(0))

					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, addConsenterConfigValue()))
					c1.cutter.CutNext = true

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).ToNot(HaveOccurred())

					network.exec(func(c *chain) {
						Eventually(c.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))
					})

					network.addChain(c4)
					c4.Start()

					
					
					
					
					Eventually(func() <-chan uint64 {
						c1.clock.Increment(interval)
						return c4.observe
					}, defaultTimeout).Should(Receive(Equal(uint64(1))))

					Eventually(c4.support.WriteBlockCallCount, defaultTimeout).Should(Equal(1))
					Eventually(c4.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))

					By("submitting new transaction to follower")
					c1.cutter.CutNext = true
					err = c4.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, defaultTimeout).Should(Equal(2))
					})
				})

				It("adding node to the cluster of 2/3 available nodes", func() {
					
					
					
					

					
					network.disconnect(2)

					c4 := newChain(timeout, channelID, dataDir, 4, &raftprotos.RaftMetadata{
						Consenters: map[uint64]*raftprotos.Consenter{},
					})
					c4.init()

					By("adding new node to the network")
					Eventually(c4.support.WriteBlockCallCount, defaultTimeout).Should(Equal(0))
					Eventually(c4.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(0))

					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, addConsenterConfigValue()))
					c1.cutter.CutNext = true

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).ToNot(HaveOccurred())

					Eventually(c1.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))
					
					Eventually(c2.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(0))
					Eventually(c3.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))

					network.addChain(c4)
					c4.Start()

					
					
					
					
					Eventually(func() <-chan uint64 {
						c1.clock.Increment(interval)
						return c4.observe
					}).Should(Receive(Equal(uint64(1))))

					Eventually(c4.support.WriteBlockCallCount, defaultTimeout).Should(Equal(1))
					Eventually(c4.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))

					By("submitting new transaction to follower")
					c1.cutter.CutNext = true
					err = c4.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					network.elect(4)

					
					network.connect(2)

					
					Eventually(func() <-chan uint64 {
						c4.clock.Increment(interval)
						return c2.observe
					}, defaultTimeout).Should(Receive(Equal(uint64(4))))

					By("submitting new transaction to re-connected node")
					c4.cutter.CutNext = true
					err = c2.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, defaultTimeout).Should(Equal(3))
					})
				})

				It("stop leader and continue reconfiguration failing over to new leader", func() {
					
					
					
					
					
					
					
					
					

					c4 := newChain(timeout, channelID, dataDir, 4, &raftprotos.RaftMetadata{
						Consenters: map[uint64]*raftprotos.Consenter{},
					})
					c4.init()

					By("adding new node to the network")
					Expect(c4.support.WriteBlockCallCount()).Should(Equal(0))
					Expect(c4.support.WriteConfigBlockCallCount()).Should(Equal(0))

					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, addConsenterConfigValue()))
					c1.cutter.CutNext = true

					stub := c1.support.WriteConfigBlockStub
					c1.support.WriteConfigBlockStub = func(block *common.Block, metadata []byte) {
						stub(block, metadata)
						
						network.disconnect(1)
						
						network.elect(2)
					}

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).ToNot(HaveOccurred())

					
					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
						})

					network.addChain(c4)
					c4.Start()
					
					
					
					
					Eventually(func() <-chan uint64 {
						c2.clock.Increment(interval)
						return c4.observe
					}, LongEventualTimeout).Should(Receive(Equal(uint64(2))))

					Eventually(c4.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c4.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))

					By("submitting new transaction to follower")
					c2.cutter.CutNext = true
					err = c4.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c4.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

					
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(1))

					network.connect(1)

					c2.clock.Increment(interval)
					
					
					Eventually(c1.observe, LongEventualTimeout).Should(Receive(Equal(uint64(2))))
					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				})

				It("stop cluster quorum and continue reconfiguration after the restart", func() {
					
					
					

					c4 := newChain(timeout, channelID, dataDir, 4, &raftprotos.RaftMetadata{
						Consenters: map[uint64]*raftprotos.Consenter{},
					})
					c4.init()

					By("adding new node to the network")
					Expect(c4.support.WriteBlockCallCount()).Should(Equal(0))
					Expect(c4.support.WriteConfigBlockCallCount()).Should(Equal(0))

					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, addConsenterConfigValue()))
					c1.cutter.CutNext = true

					stub1 := c1.support.WriteConfigBlockStub
					c1.support.WriteConfigBlockStub = func(block *common.Block, metadata []byte) {
						stub1(block, metadata)
						network.disconnect(1)
					}

					stub2 := c2.support.WriteConfigBlockStub
					c2.support.WriteConfigBlockStub = func(block *common.Block, metadata []byte) {
						stub2(block, metadata)
						network.disconnect(2)
					}

					stub3 := c3.support.WriteConfigBlockStub
					c3.support.WriteConfigBlockStub = func(block *common.Block, metadata []byte) {
						stub3(block, metadata)
						network.disconnect(3)
					}

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).ToNot(HaveOccurred())

					
					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
						})

					network.addChain(c4)

					By("reconnecting nodes back")
					for i := uint64(1); i < 4; i++ {
						network.connect(i)
					}

					c4.Start()

					By("re-elect node 2 to be a leader")
					network.elect(2)

					By("confirm new node observed a leader node")
					Eventually(func() <-chan uint64 {
						c2.clock.Increment(interval)
						return c4.observe
					}, LongEventualTimeout).Should(Receive(Equal(uint64(2))))

					Eventually(c4.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c4.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))

					By("submitting new transaction to follower")
					c2.cutter.CutNext = true
					err = c4.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					c2.clock.Increment(interval)

					
					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c4.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				})

				It("ensures that despite leader failure cluster continue to process configuration to remove the leader", func() {
					
					
					

					configEnv := newConfigEnv(channelID,
						common.HeaderType_CONFIG,
						newConfigUpdateEnv(channelID, removeConsenterConfigValue(1))) 

					c1.cutter.CutNext = true

					stub := c1.support.WriteConfigBlockStub
					c1.support.WriteConfigBlockStub = func(block *common.Block, metadata []byte) {
						stub(block, metadata)
						
						network.disconnect(1)
					}

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).ToNot(HaveOccurred())

					
					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
						})

					
					network.elect(2)

					By("submitting new transaction to follower")
					c2.cutter.CutNext = true
					err = c3.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(1))
				})

				It("removes leader from replica set", func() {
					
					
					
					
					
					

					configEnv := newConfigEnv(channelID,
						common.HeaderType_CONFIG,
						newConfigUpdateEnv(channelID, removeConsenterConfigValue(1))) 

					c1.cutter.CutNext = true

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).ToNot(HaveOccurred())

					
					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
						})

					By("by making sure remaining two nodes will elect new leader")

					
					network.elect(2)

					By("submitting transaction to new leader")
					c2.cutter.CutNext = true
					err = c2.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(1))

					By("trying to submit to new node, expected to fail")
					c1.cutter.CutNext = true
					err = c1.Order(env, 0)
					Expect(err).To(HaveOccurred())

					
					Consistently(c2.support.WriteBlockCallCount).Should(Equal(2))
					Consistently(c3.support.WriteBlockCallCount).Should(Equal(2))
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(1))
				})
			})
		})

		When("3/3 nodes are running", func() {
			JustBeforeEach(func() {
				network.init()
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
						Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
					})

				By("respect batch timeout")
				c1.cutter.CutNext = false

				err = c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())
				Eventually(c1.cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				c1.clock.WaitForNWatchersAndIncrement(timeout, 2)
				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(2))
					})
			})

			It("orders envelope on follower", func() {
				By("instructed to cut next block")
				c1.cutter.CutNext = true
				err := c2.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
					})

				By("respect batch timeout")
				c1.cutter.CutNext = false

				err = c2.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())
				Eventually(c1.cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				c1.clock.WaitForNWatchersAndIncrement(timeout, 2)
				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(2))
					})
			})

			It("leader retransmits lost messages", func() {
				

				c1.cutter.CutNext = true

				stepCnt := c1.rpc.StepCallCount()
				network.disconnect(1) 

				err := c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				network.exec(
					func(c *chain) {
						Consistently(func() int { return c.support.WriteBlockCallCount() }).Should(Equal(0))
					})

				
				Eventually(c1.rpc.StepCallCount).Should(Equal(stepCnt + 2))
				network.connect(1) 

				c1.clock.Increment(interval) 
				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
					})
			})

			It("allows the leader to create multiple normal blocks without having to wait for them to be written out", func() {
				
				network.disconnect(1)

				c1.cutter.CutNext = true
				for i := 0; i < 10; i++ {
					err := c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
				}

				Eventually(c1.BlockCreator.CreatedBlocks).Should(HaveLen(10))
			})

			It("calls BlockCreator.commitBlock on all the nodes' chains once a block is written", func() {
				normalBlock := &common.Block{
					Header:   &common.BlockHeader{},
					Data:     &common.BlockData{Data: [][]byte{[]byte("foo")}},
					Metadata: &common.BlockMetadata{Metadata: make([][]byte, 4)},
				}
				
				
				c2.BlockCreator.CreatedBlocks <- normalBlock

				c1.cutter.CutNext = true
				err := c1.Order(env, 0)
				Expect(err).ToNot(HaveOccurred())

				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
						b, _ := c.support.WriteBlockArgsForCall(0)
						Eventually(c.BlockCreator.CreatedBlocks, LongEventualTimeout).Should(HaveLen(0)) 
						
						Eventually(bytes.Equal(b.Header.Bytes(), c.BlockCreator.LastCreatedBlock.Header.Bytes()), LongEventualTimeout).Should(BeTrue())
					})
			})

			Context("handling config blocks", func() {
				var configEnv *common.Envelope
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
				})

				It("holds up block creation on leader once a config block has been created and not written out", func() {
					
					network.disconnect(1)

					c1.cutter.CutNext = true
					
					err := c1.Order(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())
					Eventually(c1.BlockCreator.CreatedBlocks).Should(HaveLen(1))

					
					tempEnv := env
					tempC1 := c1

					
					go func() {
						defer GinkgoRecover()
						err := tempC1.Order(tempEnv, 0)
						
						Expect(err).To(MatchError("chain is stopped"))
					}()

					
					Consistently(c1.BlockCreator.CreatedBlocks).Should(HaveLen(1))
				})

				It("continues creating blocks on leader after a config block has been successfully written out", func() {
					c1.cutter.CutNext = true
					
					err := c1.Configure(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())
					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteConfigBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
						})

					
					err = c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())
					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
						})
				})
			})

			When("Snapshotting is enabled", func() {
				BeforeEach(func() {
					c1.opts.SnapInterval = 1
					c1.opts.SnapshotCatchUpEntries = 1
				})

				It("keeps running if some entries in memory are purged", func() {
					
					

					i, err := c1.opts.MemoryStorage.FirstIndex()
					Expect(err).NotTo(HaveOccurred())
					Expect(i).To(Equal(uint64(1)))

					c1.cutter.CutNext = true

					err = c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
						})

					Eventually(c1.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))
					i, err = c1.opts.MemoryStorage.FirstIndex()
					Expect(err).NotTo(HaveOccurred())

					err = c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(2))
						})

					Eventually(c1.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))
					i, err = c1.opts.MemoryStorage.FirstIndex()
					Expect(err).NotTo(HaveOccurred())

					err = c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(3))
						})

					Eventually(c1.opts.MemoryStorage.FirstIndex).Should(BeNumerically(">", i))
				})

				It("lagged node can catch up using snapshot", func() {
					network.disconnect(2)

					c1.cutter.CutNext = true

					for i := 1; i <= 10; i++ {
						err := c1.Order(env, 0)
						Expect(err).NotTo(HaveOccurred())
						Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(i))
						Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(i))
					}

					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))

					network.rejoin(2, false)

					Eventually(c2.puller.PullBlockCallCount, LongEventualTimeout).Should(Equal(10))
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(10))

					files, err := ioutil.ReadDir(c2.opts.SnapDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(files).To(HaveLen(1)) 

					
					err = c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(11))
						})
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

					
					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))

					
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

					By("order envelope on follower")
					err = c3.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))

					
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				})

				It("follower cannot be elected if its log is not up-to-date", func() {
					network.disconnect(2)

					c1.cutter.CutNext = true
					err := c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

					network.disconnect(1)
					network.connect(2)

					
					for tick := 0; tick < 2*ELECTION_TICK-1; tick++ {
						c2.clock.Increment(interval)
						Consistently(c2.observe).ShouldNot(Receive(Equal(2)))
					}

					
					
					
					Consistently(c3.observe).ShouldNot(Receive())
					network.elect(3) 
				})

				It("PreVote prevents reconnected node from disturbing network", func() {
					network.disconnect(2)

					c1.cutter.CutNext = true
					err := c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					Eventually(c1.support.WriteBlockCallCount).Should(Equal(1))
					Eventually(c2.support.WriteBlockCallCount).Should(Equal(0))
					Eventually(c3.support.WriteBlockCallCount).Should(Equal(1))

					network.connect(2)

					for tick := 0; tick < 2*ELECTION_TICK-1; tick++ {
						c2.clock.Increment(interval)
						Consistently(c2.observe).ShouldNot(Receive(Equal(2)))
					}

					Consistently(c1.observe).ShouldNot(Receive())
					Consistently(c3.observe).ShouldNot(Receive())
				})

				It("follower can catch up and then campaign with success", func() {
					network.disconnect(2)

					c1.cutter.CutNext = true
					for i := 0; i < 10; i++ {
						err := c1.Order(env, 0)
						Expect(err).NotTo(HaveOccurred())
					}

					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(10))
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(10))

					network.rejoin(2, false)
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(10))

					network.disconnect(1)
					network.elect(2)
				})

				It("purges blockcutter, stops timer and discards created blocks if leadership is lost", func() {
					
					network.disconnect(1)
					normalBlock := &common.Block{
						Header:   &common.BlockHeader{},
						Data:     &common.BlockData{Data: [][]byte{[]byte("foo")}},
						Metadata: &common.BlockMetadata{Metadata: make([][]byte, 4)},
					}
					c1.BlockCreator.CreatedBlocks <- normalBlock
					Expect(len(c1.BlockCreator.CreatedBlocks)).To(Equal(1))

					
					c1.cutter.CutNext = false
					err := c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())
					Eventually(c1.cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

					
					
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(0))

					network.elect(2)
					network.rejoin(1, true)

					Eventually(c1.clock.WatcherCount, LongEventualTimeout).Should(Equal(1)) 
					Eventually(c1.cutter.CurBatch, LongEventualTimeout).Should(HaveLen(0))
					
					Eventually(c1.BlockCreator.CreatedBlocks).Should(HaveLen(0))
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(0))

					network.disconnect(2)
					n := network.elect(1) 

					err = c1.Order(env, 0)
					Expect(err).ToNot(HaveOccurred())

					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					

					c1.clock.WaitForNWatchersAndIncrement(timeout-time.Duration(n*int(interval/time.Millisecond)), 2)
					Eventually(func() int { return c1.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))
					Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))

					c1.clock.Increment(time.Duration(n * int(interval/time.Millisecond)))
					Eventually(func() int { return c1.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
					Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
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
					
					Eventually(c1.observe, LongEventualTimeout).Should(Receive(Equal(uint64(2))))
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
	puller       *mocks.FakeBlockPuller

	
	ledger       map[uint64]*common.Block
	ledgerHeight uint64
	ledgerLock   sync.RWMutex

	observe   chan uint64
	unstarted chan struct{}

	*etcdraft.Chain
}

func newChain(timeout time.Duration, channel string, dataDir string, id uint64, raftMetadata *raftprotos.RaftMetadata) *chain {
	rpc := &mocks.FakeRPC{}
	clock := fakeclock.NewFakeClock(time.Now())
	storage := raft.NewMemoryStorage()

	opts := etcdraft.Options{
		RaftID:          uint64(id),
		Clock:           clock,
		TickInterval:    interval,
		ElectionTick:    ELECTION_TICK,
		HeartbeatTick:   HEARTBEAT_TICK,
		MaxSizePerMsg:   1024 * 1024,
		MaxInflightMsgs: 256,
		RaftMetadata:    raftMetadata,
		Logger:          flogging.NewFabricLogger(zap.NewNop()),
		MemoryStorage:   storage,
		WALDir:          path.Join(dataDir, "wal"),
		SnapDir:         path.Join(dataDir, "snapshot"),
	}

	support := &consensusmocks.FakeConsenterSupport{}
	support.ChainIDReturns(channel)
	support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

	cutter := mockblockcutter.NewReceiver()
	close(cutter.Block)
	support.BlockCutterReturns(cutter)

	
	
	
	observe := make(chan uint64, 2)

	configurator := &mocks.Configurator{}
	configurator.On("Configure", mock.Anything, mock.Anything)

	puller := &mocks.FakeBlockPuller{}

	ch := make(chan struct{})
	close(ch)

	c := &chain{
		id:           id,
		support:      support,
		cutter:       cutter,
		rpc:          rpc,
		storage:      storage,
		observe:      observe,
		clock:        clock,
		opts:         opts,
		unstarted:    ch,
		configurator: configurator,
		puller:       puller,
		ledger: map[uint64]*common.Block{
			0: getSeedBlock(), 
		},
		ledgerHeight: 1,
	}

	
	
	appendBlockToLedger := func(b *common.Block, meta []byte) {
		bytes, err := proto.Marshal(&common.Metadata{Value: meta})
		Expect(err).NotTo(HaveOccurred())
		b.Metadata.Metadata[common.BlockMetadataIndex_ORDERER] = bytes
		c.ledgerLock.Lock()
		defer c.ledgerLock.Unlock()
		c.ledger[b.Header.Number] = b
		if c.ledgerHeight < b.Header.Number+1 {
			c.ledgerHeight = b.Header.Number + 1
		}
	}

	c.support.WriteBlockStub = appendBlockToLedger
	c.support.WriteConfigBlockStub = appendBlockToLedger

	
	c.support.HeightStub = func() uint64 {
		c.ledgerLock.RLock()
		defer c.ledgerLock.RUnlock()
		return c.ledgerHeight
	}

	
	c.support.BlockStub = func(number uint64) *common.Block {
		c.ledgerLock.RLock()
		defer c.ledgerLock.RUnlock()
		return c.ledger[number]
	}

	return c
}

func (c *chain) init() {
	ch, err := etcdraft.NewChain(c.support, c.opts, c.configurator, c.rpc, c.puller, c.observe)
	Expect(err).NotTo(HaveOccurred())
	c.Chain = ch
}

type network struct {
	leader uint64
	chains map[uint64]*chain

	
	
	
	
	
	connLock     sync.RWMutex
	connectivity map[uint64]chan struct{}
}

func (n *network) appendChain(c *chain) {
	n.connLock.Lock()
	n.chains[c.id] = c
	n.connLock.Unlock()
}

func (n *network) addConnection(id uint64) {
	n.connLock.Lock()
	n.connectivity[id] = make(chan struct{})
	n.connLock.Unlock()
}

func (n *network) addChain(c *chain) {
	n.addConnection(c.id)

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

	c.puller.PullBlockStub = func(i uint64) *common.Block {
		n.connLock.RLock()
		leaderChain := n.chains[n.leader]
		n.connLock.RUnlock()

		leaderChain.ledgerLock.RLock()
		defer leaderChain.ledgerLock.RUnlock()
		block := leaderChain.ledger[i]
		return block
	}

	n.appendChain(c)
}

func createNetwork(timeout time.Duration, channel string, dataDir string, raftMetadata *raftprotos.RaftMetadata) *network {
	n := &network{
		chains:       make(map[uint64]*chain),
		connectivity: make(map[uint64]chan struct{}),
	}

	for nodeID := range raftMetadata.Consenters {
		dir, err := ioutil.TempDir(dataDir, fmt.Sprintf("node-%d-", nodeID))
		Expect(err).NotTo(HaveOccurred())

		m := proto.Clone(raftMetadata).(*raftprotos.RaftMetadata)
		n.addChain(newChain(timeout, channel, dir, nodeID, m))
	}

	return n
}


func (n *network) init() {
	n.exec(func(c *chain) { c.init() })
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
			Eventually(n.chains[id].WaitReady).ShouldNot(HaveOccurred())
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
	
	
	
	
	
	t := 50 * time.Millisecond

	n.connLock.RLock()
	c := n.chains[id]
	n.connLock.RUnlock()

	var elected bool
	for !elected {
		c.clock.Increment(interval)
		tick++

		select {
		case <-time.After(t):
			
		case n := <-c.observe: 
			if n == 0 {
				
				
				
				
				
				Eventually(c.observe, LongEventualTimeout).Should(Receive(Equal(id)))
			} else {
				
				
				Expect(n).To(Equal(id))
			}
			elected = true
			break
		}
	}

	

	n.connLock.RLock()
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
	n.connLock.RUnlock()

	n.leader = id
	return tick
}

func (n *network) disconnect(i uint64) {
	n.connLock.Lock()
	defer n.connLock.Unlock()
	close(n.connectivity[i])
}

func (n *network) connect(i uint64) {
	n.connLock.Lock()
	defer n.connLock.Unlock()
	n.connectivity[i] = make(chan struct{})
}


func newConfigEnv(chainID string, headerType common.HeaderType, configUpdateEnv *common.ConfigUpdateEnvelope) *common.Envelope {
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

func newConfigUpdateEnv(chainID string, values map[string]*common.ConfigValue) *common.ConfigUpdateEnvelope {
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

func getSeedBlock() *common.Block {
	return &common.Block{
		Header:   &common.BlockHeader{},
		Data:     &common.BlockData{Data: [][]byte{[]byte("foo")}},
		Metadata: &common.BlockMetadata{Metadata: make([][]byte, 4)},
	}
}
