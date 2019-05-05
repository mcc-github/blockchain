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
	"github.com/mcc-github/blockchain/protoutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"
	"go.uber.org/zap"
)

const (
	interval            = 100 * time.Millisecond
	LongEventualTimeout = 10 * time.Second

	
	
	
	
	
	ELECTION_TICK  = 10
	HEARTBEAT_TICK = 1
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
		logger = flogging.NewFabricLogger(zap.NewExample())
		env = &common.Envelope{
			Payload: marshalOrPanic(&common.Payload{
				Header: &common.Header{ChannelHeader: marshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_MESSAGE), ChannelId: channelID})},
				Data:   []byte("TEST_MESSAGE"),
			}),
		}
	})

	Describe("Single Raft node", func() {
		var (
			configurator      *mocks.FakeConfigurator
			consenterMetadata *raftprotos.ConfigMetadata
			consenters        map[uint64]*raftprotos.Consenter
			clock             *fakeclock.FakeClock
			opts              etcdraft.Options
			support           *consensusmocks.FakeConsenterSupport
			cutter            *mockblockcutter.Receiver
			storage           *raft.MemoryStorage
			observeC          chan raft.SoftState
			chain             *etcdraft.Chain
			dataDir           string
			walDir            string
			snapDir           string
			err               error
			fakeFields        *fakeMetricsFields
		)

		BeforeEach(func() {
			configurator = &mocks.FakeConfigurator{}
			clock = fakeclock.NewFakeClock(time.Now())
			storage = raft.NewMemoryStorage()

			dataDir, err = ioutil.TempDir("", "wal-")
			Expect(err).NotTo(HaveOccurred())
			walDir = path.Join(dataDir, "wal")
			snapDir = path.Join(dataDir, "snapshot")

			observeC = make(chan raft.SoftState, 1)

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

			meta := &raftprotos.BlockMetadata{
				ConsenterIds:    make([]uint64, len(consenterMetadata.Consenters)),
				NextConsenterId: 1,
			}

			for i := range meta.ConsenterIds {
				meta.ConsenterIds[i] = meta.NextConsenterId
				meta.NextConsenterId++
			}

			consenters = map[uint64]*raftprotos.Consenter{}
			for i, c := range consenterMetadata.Consenters {
				consenters[meta.ConsenterIds[i]] = c
			}

			fakeFields = newFakeMetricsFields()

			opts = etcdraft.Options{
				RaftID:            1,
				Clock:             clock,
				TickInterval:      interval,
				ElectionTick:      ELECTION_TICK,
				HeartbeatTick:     HEARTBEAT_TICK,
				MaxSizePerMsg:     1024 * 1024,
				MaxInflightBlocks: 256,
				BlockMetadata:     meta,
				Consenters:        consenters,
				Logger:            logger,
				MemoryStorage:     storage,
				WALDir:            walDir,
				SnapDir:           snapDir,
				Metrics:           newFakeMetrics(fakeFields),
			}
		})

		campaign := func(c *etcdraft.Chain, observeC <-chan raft.SoftState) {
			Eventually(func() <-chan raft.SoftState {
				c.Consensus(&orderer.ConsensusRequest{Payload: protoutil.MarshalOrPanic(&raftpb.Message{Type: raftpb.MsgTimeoutNow})}, 0)
				return observeC
			}, LongEventualTimeout).Should(Receive(StateEqual(1, raft.StateLeader)))
		}

		JustBeforeEach(func() {
			chain, err = etcdraft.NewChain(support, opts, configurator, nil, noOpBlockPuller, observeC)
			Expect(err).NotTo(HaveOccurred())

			chain.Start()

			
			
			
			
			
			
			
			
			
			
			
			
			Eventually(func() error {
				_, err := storage.Entries(1, 1, 1)
				return err
			}, LongEventualTimeout).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			chain.Halt()
			Eventually(chain.Errored, LongEventualTimeout).Should(BeClosed())
			os.RemoveAll(dataDir)
		})

		Context("when a node starts up", func() {
			It("properly configures the communication layer", func() {
				expectedNodeConfig := nodeConfigFromMetadata(consenterMetadata)
				Eventually(configurator.ConfigureCallCount, LongEventualTimeout).Should(Equal(1))
				_, arg2 := configurator.ConfigureArgsForCall(0)
				Expect(arg2).To(Equal(expectedNodeConfig))
			})

			It("correctly sets the metrics labels and publishes requisite metrics", func() {
				type withImplementers interface {
					WithCallCount() int
					WithArgsForCall(int) []string
				}
				metricsList := []withImplementers{
					fakeFields.fakeClusterSize,
					fakeFields.fakeIsLeader,
					fakeFields.fakeCommittedBlockNumber,
					fakeFields.fakeSnapshotBlockNumber,
					fakeFields.fakeLeaderChanges,
					fakeFields.fakeProposalFailures,
					fakeFields.fakeDataPersistDuration,
					fakeFields.fakeNormalProposalsReceived,
					fakeFields.fakeConfigProposalsReceived,
				}
				for _, m := range metricsList {
					Expect(m.WithCallCount()).To(Equal(1))
					Expect(func() string {
						return m.WithArgsForCall(0)[1]
					}()).To(Equal(channelID))
				}

				Expect(fakeFields.fakeClusterSize.SetCallCount()).To(Equal(1))
				Expect(fakeFields.fakeClusterSize.SetArgsForCall(0)).To(Equal(float64(1)))
				Expect(fakeFields.fakeIsLeader.SetCallCount()).To(Equal(1))
				Expect(fakeFields.fakeIsLeader.SetArgsForCall(0)).To(Equal(float64(0)))
			})
		})

		Context("when no Raft leader is elected", func() {
			It("fails to order envelope", func() {
				err := chain.Order(env, 0)
				Expect(err).To(MatchError("no Raft leader"))
				Expect(fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(1))
				Expect(fakeFields.fakeNormalProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
				Expect(fakeFields.fakeConfigProposalsReceived.AddCallCount()).To(Equal(0))
				Expect(fakeFields.fakeProposalFailures.AddCallCount()).To(Equal(1))
				Expect(fakeFields.fakeProposalFailures.AddArgsForCall(0)).To(Equal(float64(1)))
			})

			It("starts proactive campaign", func() {
				
				
				for i := 0; i < ELECTION_TICK; i++ {
					clock.Increment(interval)
					time.Sleep(10 * time.Millisecond)
				}
				Eventually(observeC, LongEventualTimeout).Should(Receive(StateEqual(1, raft.StateLeader)))
			})
		})

		Context("when Raft leader is elected", func() {
			JustBeforeEach(func() {
				campaign(chain, observeC)
			})

			It("updates metrics upon leader election)", func() {
				Expect(fakeFields.fakeIsLeader.SetCallCount()).To(Equal(2))
				Expect(fakeFields.fakeIsLeader.SetArgsForCall(1)).To(Equal(float64(1)))
				Expect(fakeFields.fakeLeaderChanges.AddCallCount()).To(Equal(1))
				Expect(fakeFields.fakeLeaderChanges.AddArgsForCall(0)).To(Equal(float64(1)))
			})

			It("fails to order envelope if chain is halted", func() {
				chain.Halt()
				err := chain.Order(env, 0)
				Expect(err).To(MatchError("chain is stopped"))
				Expect(fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(1))
				Expect(fakeFields.fakeNormalProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
				Expect(fakeFields.fakeProposalFailures.AddCallCount()).To(Equal(1))
				Expect(fakeFields.fakeProposalFailures.AddArgsForCall(0)).To(Equal(float64(1)))
			})

			It("produces blocks following batch rules", func() {
				close(cutter.Block)

				By("cutting next batch directly")
				cutter.CutNext = true
				err := chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(1))
				Expect(fakeFields.fakeNormalProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
				Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
				Expect(fakeFields.fakeCommittedBlockNumber.SetCallCount()).Should(Equal(2)) 
				Expect(fakeFields.fakeCommittedBlockNumber.SetArgsForCall(1)).Should(Equal(float64(1)))

				
				
				
				
				
				
				Expect(fakeFields.fakeDataPersistDuration.ObserveCallCount()).Should(Equal(3))
				Expect(fakeFields.fakeDataPersistDuration.ObserveArgsForCall(0)).Should(Equal(float64(0)))
				Expect(fakeFields.fakeDataPersistDuration.ObserveArgsForCall(1)).Should(Equal(float64(0)))
				Expect(fakeFields.fakeDataPersistDuration.ObserveArgsForCall(2)).Should(Equal(float64(0)))

				By("respecting batch timeout")
				cutter.CutNext = false
				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})
				err = chain.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(2))
				Expect(fakeFields.fakeNormalProposalsReceived.AddArgsForCall(1)).To(Equal(float64(1)))

				clock.WaitForNWatchersAndIncrement(timeout, 2)
				Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				Expect(fakeFields.fakeCommittedBlockNumber.SetCallCount()).Should(Equal(3)) 
				Expect(fakeFields.fakeCommittedBlockNumber.SetArgsForCall(2)).Should(Equal(float64(2)))
				Expect(fakeFields.fakeDataPersistDuration.ObserveCallCount()).Should(Equal(4))
				Expect(fakeFields.fakeDataPersistDuration.ObserveArgsForCall(3)).Should(Equal(float64(0)))
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
					Eventually(clock.WatcherCount, LongEventualTimeout).Should(Equal(2))
				})

				It("does not enqueue if envelope is not valid", func() {
					support.ProcessNormalMsgReturns(1, errors.Errorf("Envelope is invalid"))

					err := chain.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
					Consistently(cutter.CurBatch).Should(HaveLen(0))
					Consistently(clock.WatcherCount).Should(Equal(1))
				})
			})

			It("unblocks Errored if chain is halted", func() {
				errorC := chain.Errored()
				Expect(errorC).NotTo(BeClosed())
				chain.Halt()
				Eventually(errorC, LongEventualTimeout).Should(BeClosed())
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
						Expect(err).To(MatchError("config transaction has unknown header type: CONFIG_UPDATE"))
						Expect(fakeFields.fakeConfigProposalsReceived.AddCallCount()).To(Equal(1))
						Expect(fakeFields.fakeConfigProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
						Expect(fakeFields.fakeProposalFailures.AddCallCount()).To(Equal(1))
						Expect(fakeFields.fakeProposalFailures.AddArgsForCall(0)).To(Equal(float64(1)))
					})
				})

				Context("when a type A config update comes", func() {

					Context("for existing channel", func() {

						
						BeforeEach(func() {
							newValues := map[string]*common.ConfigValue{
								"BatchTimeout": {
									Version: 1,
									Value: marshalOrPanic(&orderer.BatchTimeout{
										Timeout: "3ms",
									}),
								},
								"ConsensusType": {
									Version: 4,
								},
							}
							oldValues := map[string]*common.ConfigValue{
								"ConsensusType": {
									Version: 4,
								},
							}
							configEnv = newConfigEnv(channelID,
								common.HeaderType_CONFIG,
								newConfigUpdateEnv(channelID, oldValues, newValues),
							)
							configSeq = 0
						}) 

						Context("without revalidation (i.e. correct config sequence)", func() {

							Context("without pending normal envelope", func() {
								It("should create a config block and no normal block", func() {
									err := chain.Configure(configEnv, configSeq)
									Expect(err).NotTo(HaveOccurred())
									Expect(fakeFields.fakeConfigProposalsReceived.AddCallCount()).To(Equal(1))
									Expect(fakeFields.fakeConfigProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
									Eventually(support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
									Consistently(support.WriteBlockCallCount).Should(Equal(0))
									Expect(fakeFields.fakeCommittedBlockNumber.SetCallCount()).Should(Equal(2)) 
									Expect(fakeFields.fakeCommittedBlockNumber.SetArgsForCall(1)).Should(Equal(float64(1)))
								})
							})

							Context("with pending normal envelope", func() {
								It("should create a normal block and a config block", func() {
									
									close(cutter.Block)

									By("adding a normal envelope")
									err := chain.Order(env, 0)
									Expect(err).NotTo(HaveOccurred())
									Expect(fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(1))
									Expect(fakeFields.fakeNormalProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
									Eventually(cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

									

									By("adding a config envelope")
									err = chain.Configure(configEnv, configSeq)
									Expect(err).NotTo(HaveOccurred())
									Expect(fakeFields.fakeConfigProposalsReceived.AddCallCount()).To(Equal(1))
									Expect(fakeFields.fakeConfigProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))

									Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
									Eventually(support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
									Expect(fakeFields.fakeCommittedBlockNumber.SetCallCount()).Should(Equal(3)) 
									Expect(fakeFields.fakeCommittedBlockNumber.SetArgsForCall(2)).Should(Equal(float64(2)))
								})
							})
						})

						Context("with revalidation (i.e. incorrect config sequence)", func() {

							BeforeEach(func() {
								close(cutter.Block)
								support.SequenceReturns(1) 
							})

							It("should create config block upon correct revalidation", func() {
								support.ProcessConfigMsgReturns(configEnv, 1, nil) 

								Expect(chain.Configure(configEnv, configSeq)).To(Succeed())
								Consistently(clock.WatcherCount).Should(Equal(1))
								Eventually(support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
							})

							It("should not create config block upon incorrect revalidation", func() {
								support.ProcessConfigMsgReturns(configEnv, 1, errors.Errorf("Invalid config envelope at changed config sequence"))

								Expect(chain.Configure(configEnv, configSeq)).To(Succeed())
								Consistently(clock.WatcherCount).Should(Equal(1))
								Consistently(support.WriteConfigBlockCallCount).Should(Equal(0)) 
							})

							It("should not disturb current running timer upon incorrect revalidation", func() {
								support.ProcessNormalMsgReturns(1, nil)
								support.ProcessConfigMsgReturns(configEnv, 1, errors.Errorf("Invalid config envelope at changed config sequence"))

								Expect(chain.Order(env, configSeq)).To(Succeed())
								Eventually(clock.WatcherCount, LongEventualTimeout).Should(Equal(2))

								clock.Increment(30 * time.Minute)
								Consistently(support.WriteBlockCallCount).Should(Equal(0))

								Expect(chain.Configure(configEnv, configSeq)).To(Succeed())
								Consistently(clock.WatcherCount).Should(Equal(2))

								Consistently(support.WriteBlockCallCount).Should(Equal(0))
								Consistently(support.WriteConfigBlockCallCount).Should(Equal(0))

								clock.Increment(30 * time.Minute)
								Eventually(support.WriteBlockCallCount).Should(Equal(1))
							})
						})
					})

					Context("for creating a new channel", func() {

						
						BeforeEach(func() {
							chainID := "mychannel"
							values := make(map[string]*common.ConfigValue)
							configEnv = newConfigEnv(chainID,
								common.HeaderType_CONFIG,
								newConfigUpdateEnv(chainID, nil, values),
							)
							configSeq = 0
						}) 

						It("should be able to create a channel", func() {
							err := chain.Configure(configEnv, configSeq)
							Expect(err).NotTo(HaveOccurred())
							Eventually(support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
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
								newConfigUpdateEnv(channelID, nil, values))
							configSeq = 0

						}) 

						It("should be able to process config update of type B", func() {
							err := chain.Configure(configEnv, configSeq)
							Expect(err).NotTo(HaveOccurred())
							Expect(fakeFields.fakeConfigProposalsReceived.AddCallCount()).To(Equal(1))
							Expect(fakeFields.fakeConfigProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
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
								newConfigUpdateEnv(channelID, nil, values))
							configSeq = 0

						}) 

						It("should fail, since consenters set change is not supported", func() {
							err := chain.Configure(configEnv, configSeq)
							Expect(err).To(MatchError(ContainSubstring("update of more than one consenter at a time is not supported, requested changes: add 3 node(s), remove 1 node(s)")))
							Expect(fakeFields.fakeConfigProposalsReceived.AddCallCount()).To(Equal(1))
							Expect(fakeFields.fakeConfigProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
							Expect(fakeFields.fakeProposalFailures.AddCallCount()).To(Equal(1))
							Expect(fakeFields.fakeProposalFailures.AddArgsForCall(0)).To(Equal(float64(1)))
						})
					})

					Context("updating consenters set by exactly one node", func() {
						It("should be able to process config update adding single node", func() {
							metadata := proto.Clone(consenterMetadata).(*raftprotos.ConfigMetadata)
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
								newConfigUpdateEnv(channelID, nil, values))
							configSeq = 0

							err := chain.Configure(configEnv, configSeq)
							Expect(err).NotTo(HaveOccurred())
						})

						It("should not be able to remove node from single node cluster", func() {
							metadata := proto.Clone(consenterMetadata).(*raftprotos.ConfigMetadata)
							
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
								newConfigUpdateEnv(channelID, nil, values))
							configSeq = 0

							Expect(chain.Configure(configEnv, configSeq)).To(MatchError("empty consenter set"))
						})
					})
				})
			})

			Describe("Crash Fault Tolerance", func() {
				var (
					raftMetadata *raftprotos.BlockMetadata
				)

				BeforeEach(func() {
					raftMetadata = &raftprotos.BlockMetadata{
						ConsenterIds:    []uint64{1},
						NextConsenterId: 2,
					}
				})

				Describe("when a chain is started with existing WAL", func() {
					var (
						m1 *raftprotos.BlockMetadata
						m2 *raftprotos.BlockMetadata
					)
					JustBeforeEach(func() {
						
						
						close(cutter.Block)
						cutter.CutNext = true

						
						err := chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

						_, metadata := support.WriteBlockArgsForCall(0)
						m1 = &raftprotos.BlockMetadata{}
						proto.Unmarshal(metadata, m1)

						err = chain.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

						_, metadata = support.WriteBlockArgsForCall(1)
						m2 = &raftprotos.BlockMetadata{}
						proto.Unmarshal(metadata, m2)

						chain.Halt()
					})

					It("replays blocks from committed entries", func() {
						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)
						c.init()
						c.Start()
						defer c.Halt()

						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

						_, metadata := c.support.WriteBlockArgsForCall(0)
						m := &raftprotos.BlockMetadata{}
						proto.Unmarshal(metadata, m)
						Expect(m.RaftIndex).To(Equal(m1.RaftIndex))

						_, metadata = c.support.WriteBlockArgsForCall(1)
						m = &raftprotos.BlockMetadata{}
						proto.Unmarshal(metadata, m)
						Expect(m.RaftIndex).To(Equal(m2.RaftIndex))

						
						campaign(c.Chain, c.observe)

						c.cutter.CutNext = true

						err := c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))

					})

					It("only replays blocks after Applied index", func() {
						raftMetadata.RaftIndex = m1.RaftIndex
						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)
						c.support.WriteBlock(support.WriteBlockArgsForCall(0))

						c.init()
						c.Start()
						defer c.Halt()

						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

						_, metadata := c.support.WriteBlockArgsForCall(1)
						m := &raftprotos.BlockMetadata{}
						proto.Unmarshal(metadata, m)
						Expect(m.RaftIndex).To(Equal(m2.RaftIndex))

						
						campaign(c.Chain, c.observe)

						c.cutter.CutNext = true

						err := c.Order(env, uint64(0))
						Expect(err).NotTo(HaveOccurred())
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))
					})

					It("does not replay any block if already in sync", func() {
						raftMetadata.RaftIndex = m2.RaftIndex
						c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)
						c.init()
						c.Start()
						defer c.Halt()

						Consistently(c.support.WriteBlockCallCount).Should(Equal(0))

						
						campaign(c.Chain, c.observe)

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

							c, err := etcdraft.NewChain(support, opts, configurator, nil, noOpBlockPuller, observeC)
							Expect(c).To(BeNil())
							Expect(err).To(MatchError(ContainSubstring("permission denied")))
						})
					})
				})

				Describe("when snapshotting is enabled (snapshot interval is not zero)", func() {
					var (
						ledgerLock sync.Mutex
						ledger     map[uint64]*common.Block
					)

					countFiles := func() int {
						files, err := ioutil.ReadDir(snapDir)
						Expect(err).NotTo(HaveOccurred())
						return len(files)
					}

					BeforeEach(func() {
						opts.SnapshotCatchUpEntries = 2

						close(cutter.Block)
						cutter.CutNext = true

						ledgerLock.Lock()
						ledger = map[uint64]*common.Block{
							0: getSeedBlock(), 
						}
						ledgerLock.Unlock()

						support.WriteBlockStub = func(block *common.Block, meta []byte) {
							b := proto.Clone(block).(*common.Block)

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

					Context("Small SnapshotInterval", func() {
						BeforeEach(func() {
							opts.SnapshotIntervalSize = 1
						})

						It("writes snapshot file to snapDir", func() {
							
							

							i, _ := opts.MemoryStorage.FirstIndex()

							Expect(chain.Order(env, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
							Eventually(countFiles, LongEventualTimeout).Should(Equal(1))
							Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))
							Expect(fakeFields.fakeSnapshotBlockNumber.SetCallCount()).To(Equal(2)) 
							s, _ := opts.MemoryStorage.Snapshot()
							b := protoutil.UnmarshalBlockOrPanic(s.Data)
							Expect(fakeFields.fakeSnapshotBlockNumber.SetArgsForCall(1)).To(Equal(float64(b.Header.Number)))

							i, _ = opts.MemoryStorage.FirstIndex()

							Expect(chain.Order(env, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

							Eventually(countFiles, LongEventualTimeout).Should(Equal(2))
							Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))
							Expect(fakeFields.fakeSnapshotBlockNumber.SetCallCount()).To(Equal(3)) 
							s, _ = opts.MemoryStorage.Snapshot()
							b = protoutil.UnmarshalBlockOrPanic(s.Data)
							Expect(fakeFields.fakeSnapshotBlockNumber.SetArgsForCall(2)).To(Equal(float64(b.Header.Number)))
						})

						It("pauses chain if sync is in progress", func() {
							
							
							
							

							i, _ := opts.MemoryStorage.FirstIndex()

							Expect(chain.Order(env, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
							Eventually(countFiles, LongEventualTimeout).Should(Equal(1))
							Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))

							i, _ = opts.MemoryStorage.FirstIndex()

							Expect(chain.Order(env, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
							Eventually(countFiles, LongEventualTimeout).Should(Equal(2))
							Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))

							chain.Halt()

							c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)
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
							Eventually(done).Should(Receive(nil)) 
							Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
						})

						It("restores snapshot w/o extra entries", func() {
							
							
							
							
							
							

							Expect(chain.Order(env, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
							_, metadata := support.WriteBlockArgsForCall(0)
							m := &raftprotos.BlockMetadata{}
							proto.Unmarshal(metadata, m)

							Eventually(countFiles, LongEventualTimeout).Should(Equal(1))
							Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", 1))
							snapshot, err := opts.MemoryStorage.Snapshot() 
							Expect(err).NotTo(HaveOccurred())
							i, err := opts.MemoryStorage.FirstIndex() 
							Expect(err).NotTo(HaveOccurred())

							
							Expect(i).To(Equal(snapshot.Metadata.Index - opts.SnapshotCatchUpEntries + 1))

							chain.Halt()

							raftMetadata.RaftIndex = m.RaftIndex
							c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)
							c.opts.SnapshotIntervalSize = 1

							c.init()
							c.Start()

							
							
							Eventually(c.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(Equal(snapshot.Metadata.Index + 1))
							Eventually(c.opts.MemoryStorage.LastIndex, LongEventualTimeout).Should(Equal(snapshot.Metadata.Index))

							
							Eventually(func() <-chan raft.SoftState {
								c.clock.Increment(interval)
								return c.observe
							}, LongEventualTimeout).Should(Receive(StateEqual(1, raft.StateLeader)))

							c.cutter.CutNext = true
							err = c.Order(env, uint64(0))
							Expect(err).NotTo(HaveOccurred())
							Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

							Eventually(countFiles, LongEventualTimeout).Should(Equal(2))
							c.Halt()

							_, metadata = c.support.WriteBlockArgsForCall(0)
							m = &raftprotos.BlockMetadata{}
							proto.Unmarshal(metadata, m)
							raftMetadata.RaftIndex = m.RaftIndex
							cx := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)

							cx.init()
							cx.Start()
							defer cx.Halt()

							
							Eventually(func() <-chan raft.SoftState {
								cx.clock.Increment(interval)
								return cx.observe
							}, LongEventualTimeout).Should(Receive(StateEqual(1, raft.StateLeader)))
						})
					})

					Context("Large SnapshotInterval", func() {
						BeforeEach(func() {
							opts.SnapshotIntervalSize = 1024
						})

						It("restores snapshot w/ extra entries", func() {
							
							
							
							
							

							largeEnv := &common.Envelope{
								Payload: marshalOrPanic(&common.Payload{
									Header: &common.Header{ChannelHeader: marshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_MESSAGE), ChannelId: channelID})},
									Data:   make([]byte, 500),
								}),
							}

							By("Ordering two large envelopes to trigger snapshot")
							Expect(chain.Order(largeEnv, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

							Expect(chain.Order(largeEnv, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

							_, metadata := support.WriteBlockArgsForCall(1)
							m := &raftprotos.BlockMetadata{}
							proto.Unmarshal(metadata, m)

							
							Eventually(countFiles, LongEventualTimeout).Should(Equal(1))
							Eventually(opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", 1))
							snapshot, err := opts.MemoryStorage.Snapshot() 
							Expect(err).NotTo(HaveOccurred())
							i, err := opts.MemoryStorage.FirstIndex() 
							Expect(err).NotTo(HaveOccurred())

							
							Expect(i).To(Equal(snapshot.Metadata.Index - opts.SnapshotCatchUpEntries + 1))

							By("Ordering another envlope to append new data to memory after snaphost")
							Expect(chain.Order(env, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))

							lasti, _ := opts.MemoryStorage.LastIndex()

							chain.Halt()

							raftMetadata.RaftIndex = m.RaftIndex
							c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)
							cnt := support.WriteBlockCallCount()
							for i := 0; i < cnt; i++ {
								c.support.WriteBlock(support.WriteBlockArgsForCall(i))
							}

							By("Restarting the node")
							c.init()
							c.Start()
							defer c.Halt()

							By("Checking latest index is larger than index in snapshot")
							Eventually(c.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(Equal(snapshot.Metadata.Index + 1))
							Eventually(c.opts.MemoryStorage.LastIndex, LongEventualTimeout).Should(Equal(lasti))
						})

						When("local ledger is in sync with snapshot", func() {
							It("does not pull blocks and still respects snapshot interval", func() {
								
								
								
								
								
								
								
								
								
								

								largeEnv := &common.Envelope{
									Payload: marshalOrPanic(&common.Payload{
										Header: &common.Header{ChannelHeader: marshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_MESSAGE), ChannelId: channelID})},
										Data:   make([]byte, 500),
									}),
								}

								By("Ordering two large envelopes to trigger snapshot")
								Expect(chain.Order(largeEnv, uint64(0))).To(Succeed())
								Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

								Expect(chain.Order(largeEnv, uint64(0))).To(Succeed())
								Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

								Eventually(countFiles, LongEventualTimeout).Should(Equal(1))

								_, metadata := support.WriteBlockArgsForCall(1)
								m := &raftprotos.BlockMetadata{}
								proto.Unmarshal(metadata, m)

								By("Cutting block [3]")
								
								err = chain.Order(largeEnv, uint64(0))
								Expect(err).NotTo(HaveOccurred())
								Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))

								chain.Halt()

								raftMetadata.RaftIndex = m.RaftIndex
								c := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)
								
								c.support.WriteBlock(support.WriteBlockArgsForCall(0))
								c.support.WriteBlock(support.WriteBlockArgsForCall(1))

								c.opts.SnapshotIntervalSize = 1024

								By("Restarting node at block [2]")
								c.init()
								c.Start()
								defer c.Halt()

								
								campaign(c.Chain, c.observe)

								By("Ordering one more block to trigger snapshot")
								c.cutter.CutNext = true
								err = c.Order(largeEnv, uint64(0))
								Expect(err).NotTo(HaveOccurred())

								Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(4))
								Expect(c.puller.PullBlockCallCount()).Should(BeZero())
								
								Eventually(countFiles, LongEventualTimeout).Should(Equal(2))
							})
						})

						It("respects snapshot interval after reboot", func() {
							largeEnv := &common.Envelope{
								Payload: marshalOrPanic(&common.Payload{
									Header: &common.Header{ChannelHeader: marshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_MESSAGE), ChannelId: channelID})},
									Data:   make([]byte, 500),
								}),
							}

							Expect(chain.Order(largeEnv, uint64(0))).To(Succeed())
							Eventually(support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
							
							Consistently(countFiles).Should(Equal(0))

							_, metadata := support.WriteBlockArgsForCall(0)
							m := &raftprotos.BlockMetadata{}
							proto.Unmarshal(metadata, m)

							chain.Halt()

							raftMetadata.RaftIndex = m.RaftIndex
							c1 := newChain(10*time.Second, channelID, dataDir, 1, raftMetadata, consenters)
							cnt := support.WriteBlockCallCount()
							for i := 0; i < cnt; i++ {
								c1.support.WriteBlock(support.WriteBlockArgsForCall(i))
							}
							c1.cutter.CutNext = true
							c1.opts.SnapshotIntervalSize = 1024

							By("Restarting chain")
							c1.init()
							c1.Start()
							
							campaign(c1.Chain, c1.observe)

							Expect(c1.Order(largeEnv, uint64(0))).To(Succeed())
							Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
							
							Eventually(countFiles, LongEventualTimeout).Should(Equal(1))
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
								BlockMetadata: &raftprotos.BlockMetadata{},
								Metrics:       newFakeMetrics(newFakeMetricsFields()),
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
								BlockMetadata: &raftprotos.BlockMetadata{},
								Metrics:       newFakeMetrics(newFakeMetricsFields()),
							},
							nil,
							nil,
							noOpBlockPuller,
							nil)
						Expect(chain).NotTo(BeNil())
						Expect(err).NotTo(HaveOccurred())
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
								WALDir:        path.Join(d, "wal-dir"),
								SnapDir:       snapDir,
								Logger:        logger,
								BlockMetadata: &raftprotos.BlockMetadata{},
							},
							nil,
							nil,
							noOpBlockPuller,
							nil)
						Expect(chain).To(BeNil())
						Expect(err).To(MatchError(ContainSubstring("failed to initialize WAL: mkdir")))
					})
				})
			})
		})
	})

	Describe("2-node Raft cluster", func() {
		var (
			network      *network
			channelID    string
			timeout      time.Duration
			dataDir      string
			c1, c2       *chain
			raftMetadata *raftprotos.BlockMetadata
			consenters   map[uint64]*raftprotos.Consenter
			configEnv    *common.Envelope
		)
		BeforeEach(func() {
			var err error

			channelID = "multi-node-channel"
			timeout = 10 * time.Second

			dataDir, err = ioutil.TempDir("", "raft-test-")
			Expect(err).NotTo(HaveOccurred())

			raftMetadata = &raftprotos.BlockMetadata{
				ConsenterIds:    []uint64{1, 2},
				NextConsenterId: 3,
			}

			consenters = map[uint64]*raftprotos.Consenter{
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
			}

			metadata := &raftprotos.ConfigMetadata{
				Options: &raftprotos.Options{
					TickInterval:         "500ms",
					ElectionTick:         10,
					HeartbeatTick:        1,
					MaxInflightBlocks:    5,
					SnapshotIntervalSize: 200,
				},
				Consenters: []*raftprotos.Consenter{consenters[2]},
			}
			value := map[string]*common.ConfigValue{
				"ConsensusType": {
					Version: 1,
					Value: marshalOrPanic(&orderer.ConsensusType{
						Metadata: marshalOrPanic(metadata),
					}),
				},
			}
			
			configEnv = newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, value))

			network = createNetwork(timeout, channelID, dataDir, raftMetadata, consenters)
			c1, c2 = network.chains[1], network.chains[2]
			c1.cutter.CutNext = true
			network.init()
			network.start()
		})

		AfterEach(func() {
			network.stop()
			os.RemoveAll(dataDir)
		})

		It("can remove leader by reconfiguring cluster", func() {
			network.elect(1)

			Expect(c1.Configure(configEnv, 0)).To(Succeed())
			network.exec(func(c *chain) {
				Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
				Eventually(c.configurator.ConfigureCallCount, LongEventualTimeout).Should(Equal(2))
			})

			Consistently(c1.Chain.Errored).ShouldNot(BeClosed())
			c1.clock.WaitForNWatchersAndIncrement(ELECTION_TICK*interval, 2)
			Eventually(c1.Chain.Errored, LongEventualTimeout).Should(BeClosed())
			close(c1.stopped) 

			By("Electing 2 as new leader")
			network.elect(2)

			By("Asserting leader can still serve requests as single-node cluster")
			c2.cutter.CutNext = true
			Expect(c2.Order(env, 0)).To(Succeed())
			Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
		})

		It("can remove follower by reconfiguring cluster", func() {
			network.elect(2)

			Expect(c1.Configure(configEnv, 0)).To(Succeed())
			network.exec(func(c *chain) {
				Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
				Eventually(c.configurator.ConfigureCallCount, LongEventualTimeout).Should(Equal(2))
			})
			Eventually(c1.Chain.Errored, LongEventualTimeout).Should(BeClosed())

			By("Asserting leader can still serve requests as single-node cluster")
			c2.cutter.CutNext = true
			Expect(c2.Order(env, 0)).To(Succeed())
			Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
		})
	})

	Describe("3-node Raft cluster", func() {
		var (
			network      *network
			channelID    string
			timeout      time.Duration
			dataDir      string
			c1, c2, c3   *chain
			raftMetadata *raftprotos.BlockMetadata
			consenters   map[uint64]*raftprotos.Consenter
		)

		BeforeEach(func() {
			var err error

			channelID = "multi-node-channel"
			timeout = 10 * time.Second

			dataDir, err = ioutil.TempDir("", "raft-test-")
			Expect(err).NotTo(HaveOccurred())

			raftMetadata = &raftprotos.BlockMetadata{
				ConsenterIds:    []uint64{1, 2, 3},
				NextConsenterId: 4,
			}

			consenters = map[uint64]*raftprotos.Consenter{
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
			}

			network = createNetwork(timeout, channelID, dataDir, raftMetadata, consenters)
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
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() int { return c1.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
				Eventually(func() int { return c2.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
				Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))

				network.start(3)

				c1.clock.Increment(interval)
				Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))

				network.stop()
			})

			It("late node receives snapshot from leader", func() {
				c1.opts.SnapshotIntervalSize = 1
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
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() int { return c1.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
				Eventually(func() int { return c2.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
				Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))

				err = c1.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())

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
				options = &raftprotos.Options{
					TickInterval:         "500ms",
					ElectionTick:         10,
					HeartbeatTick:        1,
					MaxInflightBlocks:    5,
					SnapshotIntervalSize: 200,
				}
				updateRaftConfigValue = func(metadata *raftprotos.ConfigMetadata) map[string]*common.ConfigValue {
					return map[string]*common.ConfigValue{
						"ConsensusType": {
							Version: 1,
							Value: marshalOrPanic(&orderer.ConsensusType{
								Metadata: marshalOrPanic(metadata),
							}),
						},
					}
				}
				addConsenterConfigValue = func() map[string]*common.ConfigValue {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					newConsenter := &raftprotos.Consenter{
						Host:          "localhost",
						Port:          7050,
						ServerTlsCert: serverTLSCert(tlsCA),
						ClientTlsCert: clientTLSCert(tlsCA),
					}
					metadata.Consenters = append(metadata.Consenters, newConsenter)
					return updateRaftConfigValue(metadata)
				}
				removeConsenterConfigValue = func(id uint64) map[string]*common.ConfigValue {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for nodeID, consenter := range consenters {
						if nodeID == id {
							continue
						}
						metadata.Consenters = append(metadata.Consenters, consenter)
					}
					return updateRaftConfigValue(metadata)
				}
				createChannelEnv = func(metadata *raftprotos.ConfigMetadata) *common.Envelope {
					configEnv := newConfigEnv("another-channel",
						common.HeaderType_CONFIG,
						newConfigUpdateEnv(channelID, nil, updateRaftConfigValue(metadata)))

					
					return &common.Envelope{
						Payload: marshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: marshalOrPanic(&common.ChannelHeader{
									Type:      int32(common.HeaderType_ORDERER_TRANSACTION),
									ChannelId: channelID,
								}),
							},
							Data: marshalOrPanic(configEnv),
						}),
					}
				}
			)

			BeforeEach(func() {
				network.exec(func(c *chain) {
					c.opts.EvictionSuspicion = time.Millisecond * 100
					c.opts.LeaderCheckInterval = time.Millisecond * 100
				})

				network.init()
				network.start()
				network.elect(1)

				By("Submitting first tx to cut the block")
				c1.cutter.CutNext = true
				err := c1.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())

				c1.clock.Increment(interval)

				network.exec(
					func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, defaultTimeout).Should(Equal(1))
					})
			})

			AfterEach(func() {
				network.stop()
			})

			Context("channel creation", func() {
				It("succeeds with valid config metadata", func() {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(Succeed())
					network.exec(func(c *chain) {
						Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
					})
				})

				It("fails with invalid timeout", func() {
					metadata := &raftprotos.ConfigMetadata{Options: proto.Clone(options).(*raftprotos.Options)}
					metadata.Options.ElectionTick = metadata.Options.HeartbeatTick
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(MatchError(
						"ElectionTick (1) must be greater than HeartbeatTick (1)"))
				})

				It("fails with zero ElectionTick", func() {
					metadata := &raftprotos.ConfigMetadata{Options: proto.Clone(options).(*raftprotos.Options)}
					metadata.Options.ElectionTick = 0
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(MatchError(
						"none of HeartbeatTick (1), ElectionTick (0) and MaxInflightBlocks (5) can be zero"))
				})

				It("fails with zero max inflgith blocks", func() {
					metadata := &raftprotos.ConfigMetadata{Options: proto.Clone(options).(*raftprotos.Options)}
					metadata.Options.MaxInflightBlocks = 0
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(MatchError(
						"none of HeartbeatTick (1), ElectionTick (10) and MaxInflightBlocks (0) can be zero"))
				})

				It("fails with invalid tick interval", func() {
					metadata := &raftprotos.ConfigMetadata{Options: proto.Clone(options).(*raftprotos.Options)}
					metadata.Options.TickInterval = "invalid"
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(MatchError(ContainSubstring(
						"failed to parse TickInterval (invalid) to time duration")))
				})

				It("fails with zero tick interval", func() {
					metadata := &raftprotos.ConfigMetadata{Options: proto.Clone(options).(*raftprotos.Options)}
					metadata.Options.TickInterval = "0"
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(MatchError(ContainSubstring(
						"TickInterval cannot be zero")))
				})

				It("fails with empty consenter set", func() {
					metadata := &raftprotos.ConfigMetadata{Options: options}

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(MatchError(
						"empty consenter set"))
				})

				It("fails with invalid certificate for non PEM certificates", func() {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}
					metadata.Consenters[0].ClientTlsCert = []byte("hello")

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(MatchError(
						"client TLS certificate is not PEM encoded: hello"))
				})

				It("fails with invalid certificate for malformed certificates", func() {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}
					metadata.Consenters[0].ServerTlsCert = pem.EncodeToMemory(&pem.Block{Bytes: []byte("hello")})

					Expect(c1.Configure(createChannelEnv(metadata), 0).Error()).To(ContainSubstring(
						"server TLS certificate has invalid ASN1 structure"))
				})

				It("fails with extra consenter", func() {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for _, consenter := range consenters {
						metadata.Consenters = append(metadata.Consenters, consenter)
					}
					metadata.Consenters = append(
						metadata.Consenters,
						&raftprotos.Consenter{
							Host:          "localhost",
							Port:          7050,
							ServerTlsCert: serverTLSCert(tlsCA),
							ClientTlsCert: clientTLSCert(tlsCA),
						})

					Expect(c1.Configure(createChannelEnv(metadata), 0)).To(MatchError(
						"new channel has consenter that is not part of system consenter set"))
				})
			})

			Context("reconfiguration", func() {
				It("cannot change consenter set by more than 1 node", func() {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for id, consenter := range consenters {
						if id == 2 || id == 3 {
							
							continue
						}
						metadata.Consenters = append(metadata.Consenters, consenter)
					}

					value := map[string]*common.ConfigValue{
						"ConsensusType": {
							Version: 1,
							Value: marshalOrPanic(&orderer.ConsensusType{
								Metadata: marshalOrPanic(metadata),
							}),
						},
					}

					By("creating new configuration with removed node and new one")
					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, value))
					c1.cutter.CutNext = true

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).To(MatchError(ContainSubstring("update of more than one consenter at a time is not supported, requested changes: add 0 node(s), remove 2 node(s)")))
				})

				It("rejects invalid certificates", func() {
					configMetadata := &raftprotos.ConfigMetadata{Options: options}
					for _, consenter := range consenters {
						configMetadata.Consenters = append(configMetadata.Consenters, consenter)
					}
					configMetadata.Consenters[0].ServerTlsCert = []byte("hello")
					value := map[string]*common.ConfigValue{
						"ConsensusType": {
							Version: 1,
							Value: marshalOrPanic(&orderer.ConsensusType{
								Metadata: marshalOrPanic(configMetadata),
							}),
						},
					}

					By("creating new configuration with invalid certificate")
					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, value))
					c1.cutter.CutNext = true

					By("sending config transaction")
					Expect(c1.Configure(configEnv, 0)).To(MatchError("server TLS certificate is not PEM encoded: hello"))
				})

				It("can rotate certificate by adding and removing 1 node in one config update", func() {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for id, consenter := range consenters {
						if id == 2 {
							
							continue
						}
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
					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, value))
					c1.cutter.CutNext = true

					By("sending config transaction")
					Expect(c1.Configure(configEnv, 0)).To(Succeed())

					network.exec(func(c *chain) {
						Eventually(c.configurator.ConfigureCallCount, LongEventualTimeout).Should(Equal(2))
					})
				})

				It("rotates leader certificate and triggers leadership transfer", func() {
					metadata := &raftprotos.ConfigMetadata{Options: options}
					for id, consenter := range consenters {
						if id == 1 {
							
							continue
						}
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
					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, value))
					c1.cutter.CutNext = true

					By("sending config transaction")
					Expect(c1.Configure(configEnv, 0)).To(Succeed())

					Eventually(c1.observe, LongEventualTimeout).Should(Receive(BeFollower()))
					network.exec(func(c *chain) {
						Eventually(c.configurator.ConfigureCallCount, LongEventualTimeout).Should(Equal(2))
					})
				})

				When("Leader is disconnected after cert rotation", func() {
					It("still configures communication after failed leader transfer attempt", func() {
						metadata := &raftprotos.ConfigMetadata{Options: options}
						for id, consenter := range consenters {
							if id == 1 {
								
								continue
							}
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
						configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, value))
						c1.cutter.CutNext = true

						step1 := c1.getStepFunc()
						count := c1.rpc.SendConsensusCallCount() 
						c1.setStepFunc(func(dest uint64, msg *orderer.ConsensusRequest) error {
							
							
							
							if c1.rpc.SendConsensusCallCount() == count+4 {
								defer network.disconnect(1)
							}

							return step1(dest, msg)
						})

						By("sending config transaction")
						Expect(c1.Configure(configEnv, 0)).To(Succeed())

						Consistently(c1.observe).ShouldNot(Receive())
						network.exec(func(c *chain) {
							Eventually(c.configurator.ConfigureCallCount, LongEventualTimeout).Should(Equal(2))
						})
					})
				})

				When("Follower is disconnected while leader cert is being rotated", func() {
					It("still configures communication and transfer leader", func() {
						metadata := &raftprotos.ConfigMetadata{Options: options}
						for id, consenter := range consenters {
							if id == 1 {
								
								continue
							}
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

						cnt := c1.rpc.SendConsensusCallCount()
						network.disconnect(3)

						
						
						
						Eventually(func() int {
							c1.clock.Increment(interval)
							return c1.rpc.SendConsensusCallCount()
						}, LongEventualTimeout).Should(BeNumerically(">=", cnt+5))

						By("creating new configuration with removed node and new one")
						configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, value))
						c1.cutter.CutNext = true

						By("sending config transaction")
						Expect(c1.Configure(configEnv, 0)).To(Succeed())

						Eventually(c1.observe, LongEventualTimeout).Should(Receive(StateEqual(2, raft.StateFollower)))
						network.Lock()
						network.leader = 2 
						network.Unlock()
						network.disconnect(1)

						network.exec(func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
							Eventually(c.configurator.ConfigureCallCount, LongEventualTimeout).Should(Equal(2))
						}, 1, 2)

						network.join(3, true)
						Eventually(c3.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
						Eventually(c3.configurator.ConfigureCallCount, LongEventualTimeout).Should(Equal(2))

						By("Ordering normal transaction")
						c2.cutter.CutNext = true
						Expect(c3.Order(env, 0)).To(Succeed())
						network.exec(func(c *chain) {
							Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
						}, 2, 3)
					})
				})

				When("two type B config are sent back-to-back", func() {
					It("discards the second", func() {
						
						
						
						c1.cutter.CutNext = true
						configEnvAdd := newConfigEnv(channelID,
							common.HeaderType_CONFIG,
							newConfigUpdateEnv(channelID, nil, addConsenterConfigValue()))
						configEnvRm := newConfigEnv(channelID,
							common.HeaderType_CONFIG,
							newConfigUpdateEnv(channelID, nil, removeConsenterConfigValue(3)))

						By("Submitting two config tx back-to-back")
						c1.support.SequenceReturnsOnCall(1, 0)
						c1.support.SequenceReturnsOnCall(2, 1)
						c1.support.ProcessConfigMsgReturns(configEnvRm, 1, nil)

						Expect(c1.Configure(configEnvAdd, 0)).To(Succeed())
						Expect(c1.Configure(configEnvRm, 0)).To(Succeed())
						network.exec(func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
							Consistently(c.support.WriteConfigBlockCallCount).Should(Equal(1))
						})
					})
				})

				It("adding node to the cluster", func() {
					addConsenterUpdate := addConsenterConfigValue()
					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, addConsenterUpdate))
					c1.cutter.CutNext = true

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())
					Expect(c1.fakeFields.fakeConfigProposalsReceived.AddCallCount()).To(Equal(1))
					Expect(c1.fakeFields.fakeConfigProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))

					network.exec(func(c *chain) {
						Eventually(c.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))
						Eventually(c.fakeFields.fakeClusterSize.SetCallCount, LongEventualTimeout).Should(Equal(2))
						Expect(c.fakeFields.fakeClusterSize.SetArgsForCall(1)).To(Equal(float64(4)))
					})

					_, raftmetabytes := c1.support.WriteConfigBlockArgsForCall(0)
					meta := &common.Metadata{Value: raftmetabytes}
					raftmeta, err := etcdraft.ReadBlockMetadata(meta, nil)
					Expect(err).NotTo(HaveOccurred())

					c4 := newChain(timeout, channelID, dataDir, 4, raftmeta, consenters)
					
					
					c4.support.WriteBlock(c1.support.WriteBlockArgsForCall(0))
					c4.support.WriteConfigBlock(c1.support.WriteConfigBlockArgsForCall(0))
					c4.init()

					network.addChain(c4)
					c4.Start()

					
					
					
					
					Eventually(func() <-chan raft.SoftState {
						c1.clock.Increment(interval)
						return c4.observe
					}, defaultTimeout).Should(Receive(Equal(raft.SoftState{Lead: 1, RaftState: raft.StateFollower})))

					Eventually(c4.support.WriteBlockCallCount, defaultTimeout).Should(Equal(1))
					Eventually(c4.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))

					By("submitting new transaction to follower")
					c1.cutter.CutNext = true
					err = c4.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
					Expect(c4.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(1))
					Expect(c4.fakeFields.fakeNormalProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))

					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, defaultTimeout).Should(Equal(2))
					})

					By("resubmitting the config update by mistake")
					c1.cutter.CutNext = true
					consensusType := &orderer.ConsensusType{}
					proto.Unmarshal(addConsenterUpdate["ConsensusType"].Value, consensusType)
					duplicatedMetadata := &raftprotos.ConfigMetadata{Options: options}
					proto.Unmarshal(consensusType.Metadata, duplicatedMetadata)
					duplicatedMetadata.Consenters = append(duplicatedMetadata.Consenters, duplicatedMetadata.Consenters[1])
					configEnv = newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, map[string]*common.ConfigValue{
						"ConsensusType": {
							Version: 1,
							Value: marshalOrPanic(&orderer.ConsensusType{
								Metadata: marshalOrPanic(duplicatedMetadata),
							}),
						},
					}))

					err = c1.Configure(configEnv, 1)
					By("Expecting it to be rejected")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("duplicate consenter"))
					Expect(err.Error()).To(ContainSubstring(string(duplicatedMetadata.Consenters[1].ServerTlsCert)))
					Expect(err.Error()).To(ContainSubstring(string(duplicatedMetadata.Consenters[1].ClientTlsCert)))
				})

				It("does not reconfigure raft cluster if it's a channel creation tx", func() {
					configEnv := newConfigEnv("another-channel",
						common.HeaderType_CONFIG,
						newConfigUpdateEnv(channelID, nil, removeConsenterConfigValue(2)))

					
					channelCreationEnv := &common.Envelope{
						Payload: marshalOrPanic(&common.Payload{
							Header: &common.Header{
								ChannelHeader: marshalOrPanic(&common.ChannelHeader{
									Type:      int32(common.HeaderType_ORDERER_TRANSACTION),
									ChannelId: channelID,
								}),
							},
							Data: marshalOrPanic(configEnv),
						}),
					}

					c1.cutter.CutNext = true

					Expect(c1.Configure(channelCreationEnv, 0)).To(Succeed())
					network.exec(func(c *chain) {
						Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
					})

					
					Consistently(c2.Errored).ShouldNot(BeClosed())
					Expect(c2.Order(env, 0)).To(Succeed())

					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					})
				})

				It("adding node to the cluster of 2/3 available nodes", func() {
					
					
					
					

					
					network.disconnect(2)

					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, addConsenterConfigValue()))
					c1.cutter.CutNext = true

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())

					Eventually(c1.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))
					
					Eventually(c2.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(0))
					Eventually(c3.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))

					_, raftmetabytes := c1.support.WriteConfigBlockArgsForCall(0)
					meta := &common.Metadata{Value: raftmetabytes}
					raftmeta, err := etcdraft.ReadBlockMetadata(meta, nil)
					Expect(err).NotTo(HaveOccurred())

					c4 := newChain(timeout, channelID, dataDir, 4, raftmeta, consenters)
					
					
					c4.support.WriteBlock(c1.support.WriteBlockArgsForCall(0))
					c4.support.WriteConfigBlock(c1.support.WriteConfigBlockArgsForCall(0))

					c4.init()

					network.addChain(c4)
					c4.start()
					Expect(c4.WaitReady()).To(Succeed())
					network.join(4, true)

					Eventually(c4.support.WriteBlockCallCount, defaultTimeout).Should(Equal(1))
					Eventually(c4.support.WriteConfigBlockCallCount, defaultTimeout).Should(Equal(1))

					By("submitting new transaction to follower")
					c1.cutter.CutNext = true
					err = c4.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
					Eventually(c4.support.WriteBlockCallCount, defaultTimeout).Should(Equal(2))

					
					network.elect(4)

					c2.Consensus(&orderer.ConsensusRequest{Payload: protoutil.MarshalOrPanic(&raftpb.Message{Type: raftpb.MsgTimeoutNow})}, 0)
					Eventually(c2.observe, LongEventualTimeout).Should(Receive(StateEqual(0, raft.StateCandidate)))

					
					network.connect(2)

					
					Eventually(func() <-chan raft.SoftState {
						c4.clock.Increment(interval)
						return c2.observe
					}, defaultTimeout).Should(Receive(Equal(raft.SoftState{Lead: 4, RaftState: raft.StateFollower})))

					By("submitting new transaction to re-connected node")
					c4.cutter.CutNext = true
					err = c2.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, defaultTimeout).Should(Equal(3))
					})
				})

				It("stop leader and continue reconfiguration failing over to new leader", func() {
					
					
					
					
					
					
					
					
					

					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, addConsenterConfigValue()))
					c1.cutter.CutNext = true

					step1 := c1.getStepFunc()
					count := c1.rpc.SendConsensusCallCount() 
					c1.setStepFunc(func(dest uint64, msg *orderer.ConsensusRequest) error {
						
						
						
						if c1.rpc.SendConsensusCallCount() == count+4 {
							defer network.disconnect(1)
						}

						return step1(dest, msg)
					})

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())

					
					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
						})
					c1.setStepFunc(step1)

					
					i2, _ := c2.storage.LastIndex() 
					i3, _ := c3.storage.LastIndex()
					candidate := uint64(2)
					if i3 > i2 {
						candidate = 3
					}
					network.chains[candidate].cutter.CutNext = true
					network.elect(candidate)

					_, raftmetabytes := c1.support.WriteConfigBlockArgsForCall(0)
					meta := &common.Metadata{Value: raftmetabytes}
					raftmeta, err := etcdraft.ReadBlockMetadata(meta, nil)
					Expect(err).NotTo(HaveOccurred())

					c4 := newChain(timeout, channelID, dataDir, 4, raftmeta, consenters)
					
					
					c4.support.WriteBlock(c1.support.WriteBlockArgsForCall(0))
					c4.support.WriteConfigBlock(c1.support.WriteConfigBlockArgsForCall(0))
					c4.init()

					network.addChain(c4)
					c4.start()
					Expect(c4.WaitReady()).To(Succeed())
					network.join(4, true)

					Eventually(c4.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c4.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))

					By("submitting new transaction to follower")
					err = c4.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c4.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

					
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(1))

					network.join(1, true)
					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				})

				It("stop cluster quorum and continue reconfiguration after the restart", func() {
					
					
					

					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, addConsenterConfigValue()))
					c1.cutter.CutNext = true

					step1 := c1.getStepFunc()
					count := c1.rpc.SendConsensusCallCount() 
					c1.setStepFunc(func(dest uint64, msg *orderer.ConsensusRequest) error {
						
						
						
						if c1.rpc.SendConsensusCallCount() == count+4 {
							defer func() {
								network.disconnect(1)
								network.disconnect(2)
								network.disconnect(3)
							}()
						}

						return step1(dest, msg)
					})

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())

					
					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
						})

					_, raftmetabytes := c1.support.WriteConfigBlockArgsForCall(0)
					meta := &common.Metadata{Value: raftmetabytes}
					raftmeta, err := etcdraft.ReadBlockMetadata(meta, nil)
					Expect(err).NotTo(HaveOccurred())

					c4 := newChain(timeout, channelID, dataDir, 4, raftmeta, consenters)
					
					
					c4.support.WriteBlock(c1.support.WriteBlockArgsForCall(0))
					c4.support.WriteConfigBlock(c1.support.WriteConfigBlockArgsForCall(0))
					c4.init()

					network.addChain(c4)

					By("reconnecting nodes back")
					for i := uint64(1); i < 4; i++ {
						network.connect(i)
					}

					
					i2, _ := c2.storage.LastIndex() 
					i3, _ := c3.storage.LastIndex()
					candidate := uint64(2)
					if i3 > i2 {
						candidate = 3
					}
					network.chains[candidate].cutter.CutNext = true
					network.elect(candidate)

					c4.start()
					Expect(c4.WaitReady()).To(Succeed())
					network.join(4, false)

					Eventually(c4.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c4.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))

					By("submitting new transaction to follower")
					err = c4.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					
					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c4.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				})

				It("ensures that despite leader failure cluster continue to process configuration to remove the leader", func() {
					
					
					

					configEnv := newConfigEnv(channelID,
						common.HeaderType_CONFIG,
						newConfigUpdateEnv(channelID, nil, removeConsenterConfigValue(1))) 

					c1.cutter.CutNext = true

					step1 := c1.getStepFunc()
					count := c1.rpc.SendConsensusCallCount() 
					c1.setStepFunc(func(dest uint64, msg *orderer.ConsensusRequest) error {
						
						
						
						if c1.rpc.SendConsensusCallCount() == count+4 {
							defer network.disconnect(1)
						}

						return step1(dest, msg)
					})

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())

					network.exec(func(c *chain) {
						Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
					})

					
					i2, _ := c2.storage.LastIndex() 
					i3, _ := c3.storage.LastIndex()
					candidate := uint64(2)
					if i3 > i2 {
						candidate = 3
					}
					network.chains[candidate].cutter.CutNext = true
					network.elect(candidate)

					By("submitting new transaction to follower")
					err = c3.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				})

				It("removes leader from replica set", func() {
					
					
					
					
					
					

					configEnv := newConfigEnv(channelID,
						common.HeaderType_CONFIG,
						newConfigUpdateEnv(channelID, nil, removeConsenterConfigValue(1))) 

					c1.cutter.CutNext = true

					By("sending config transaction")
					err := c1.Configure(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())

					
					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
							Eventually(c.fakeFields.fakeClusterSize.SetCallCount, LongEventualTimeout).Should(Equal(2))
							Expect(c.fakeFields.fakeClusterSize.SetArgsForCall(1)).To(Equal(float64(2)))
						})

					
					c1.clock.WaitForNWatchersAndIncrement(ELECTION_TICK*interval, 2)
					Eventually(c1.Errored, LongEventualTimeout).Should(BeClosed())
					close(c1.stopped)

					By("making sure remaining two nodes will elect new leader")

					
					network.elect(2)

					By("submitting transaction to new leader")
					c2.cutter.CutNext = true
					err = c2.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

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

				It("does not deadlock if leader steps down while config block is in-flight", func() {
					configEnv := newConfigEnv(channelID, common.HeaderType_CONFIG, newConfigUpdateEnv(channelID, nil, addConsenterConfigValue()))
					c1.cutter.CutNext = true

					signal := make(chan struct{})
					stub := c1.support.WriteConfigBlockStub
					c1.support.WriteConfigBlockStub = func(b *common.Block, meta []byte) {
						signal <- struct{}{}
						<-signal
						stub(b, meta)
					}

					By("Sending config transaction")
					Expect(c1.Configure(configEnv, 0)).To(Succeed())

					Eventually(signal, LongEventualTimeout).Should(Receive())
					network.disconnect(1)

					By("Ticking leader till it steps down")
					Eventually(func() raft.SoftState {
						c1.clock.Increment(interval)
						return c1.Node.Status().SoftState
					}, LongEventualTimeout).Should(StateEqual(0, raft.StateFollower))

					close(signal)

					Eventually(c1.observe, LongEventualTimeout).Should(Receive(StateEqual(0, raft.StateFollower)))

					By("Re-electing 1 as leader")
					network.connect(1)
					network.elect(1)

					_, raftmetabytes := c1.support.WriteConfigBlockArgsForCall(0)
					meta := &common.Metadata{Value: raftmetabytes}
					raftmeta, err := etcdraft.ReadBlockMetadata(meta, nil)
					Expect(err).NotTo(HaveOccurred())

					c4 := newChain(timeout, channelID, dataDir, 4, raftmeta, consenters)
					
					
					c4.support.WriteBlock(c1.support.WriteBlockArgsForCall(0))
					c4.support.WriteConfigBlock(c1.support.WriteConfigBlockArgsForCall(0))
					c4.init()

					network.addChain(c4)
					c4.Start()

					Eventually(func() <-chan raft.SoftState {
						c1.clock.Increment(interval)
						return c4.observe
					}, LongEventualTimeout).Should(Receive(StateEqual(1, raft.StateFollower)))

					By("Submitting tx to confirm network is still working")
					Expect(c1.Order(env, 0)).To(Succeed())

					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
						Eventually(c.support.WriteConfigBlockCallCount, LongEventualTimeout).Should(Equal(1))
					})
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

			It("correctly sets the cluster size and leadership metrics", func() {
				
				network.exec(func(c *chain) {
					Expect(c.fakeFields.fakeLeaderChanges.AddCallCount()).Should(Equal(1))
					Expect(c.fakeFields.fakeLeaderChanges.AddArgsForCall(0)).Should(Equal(float64(1)))
					Expect(c.fakeFields.fakeClusterSize.SetCallCount()).Should(Equal(1))
					Expect(c.fakeFields.fakeClusterSize.SetArgsForCall(0)).To(Equal(float64(3)))
				})
				
				Expect(c1.fakeFields.fakeIsLeader.SetCallCount()).Should(Equal(2))
				Expect(c1.fakeFields.fakeIsLeader.SetArgsForCall(1)).Should(Equal(float64(1)))
				
				Expect(c2.fakeFields.fakeIsLeader.SetCallCount()).Should(Equal(1))
				Expect(c2.fakeFields.fakeIsLeader.SetArgsForCall(0)).Should(Equal(float64(0)))
				Expect(c3.fakeFields.fakeIsLeader.SetCallCount()).Should(Equal(1))
				Expect(c3.fakeFields.fakeIsLeader.SetArgsForCall(0)).Should(Equal(float64(0)))
			})

			It("orders envelope on leader", func() {
				By("instructed to cut next block")
				c1.cutter.CutNext = true
				err := c1.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(c1.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(1))
				Expect(c1.fakeFields.fakeNormalProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))

				network.exec(
					func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					})

				By("respect batch timeout")
				c1.cutter.CutNext = false

				err = c1.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(c1.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(2))
				Expect(c1.fakeFields.fakeNormalProposalsReceived.AddArgsForCall(1)).To(Equal(float64(1)))
				Eventually(c1.cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				c1.clock.WaitForNWatchersAndIncrement(timeout, 2)
				network.exec(
					func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					})
			})

			It("orders envelope on follower", func() {
				By("instructed to cut next block")
				c1.cutter.CutNext = true
				err := c2.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(c2.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(1))
				Expect(c2.fakeFields.fakeNormalProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
				Expect(c1.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(0))

				network.exec(
					func(c *chain) {
						Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
					})

				By("respect batch timeout")
				c1.cutter.CutNext = false

				err = c2.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(c2.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(2))
				Expect(c2.fakeFields.fakeNormalProposalsReceived.AddArgsForCall(1)).To(Equal(float64(1)))
				Expect(c1.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(0))
				Eventually(c1.cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

				c1.clock.WaitForNWatchersAndIncrement(timeout, 2)
				network.exec(
					func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					})
			})

			When("MaxInflightBlocks is reached", func() {
				BeforeEach(func() {
					network.exec(func(c *chain) { c.opts.MaxInflightBlocks = 1 })
				})

				It("waits for in flight blocks to be committed", func() {
					c1.cutter.CutNext = true
					
					network.disconnect(1)

					Expect(c1.Order(env, 0)).To(Succeed())

					doneProp := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						Expect(c1.Order(env, 0)).To(Succeed())
						close(doneProp)
					}()
					
					Consistently(doneProp).ShouldNot(BeClosed())
					network.exec(func(c *chain) {
						Consistently(c.support.WriteBlockCallCount).Should(BeZero())
					})

					network.connect(1)
					c1.clock.Increment(interval)

					Eventually(doneProp, LongEventualTimeout).Should(BeClosed())
					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					})
				})

				It("resets block in flight when steps down from leader", func() {
					c1.cutter.CutNext = true
					c2.cutter.CutNext = true
					
					network.disconnect(1)

					Expect(c1.Order(env, 0)).To(Succeed())

					doneProp := make(chan struct{})
					go func() {
						defer GinkgoRecover()

						Expect(c1.Order(env, 0)).To(Succeed())
						close(doneProp)
					}()
					
					Consistently(doneProp).ShouldNot(BeClosed())
					network.exec(func(c *chain) {
						Consistently(c.support.WriteBlockCallCount).Should(BeZero())
					})

					network.elect(2)
					Expect(c3.Order(env, 0)).To(Succeed())
					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

					network.connect(1)
					c2.clock.Increment(interval)

					Eventually(doneProp, LongEventualTimeout).Should(BeClosed())
					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
					})
				})
			})

			When("leader is disconnected", func() {
				It("proactively steps down to follower", func() {
					network.disconnect(1)

					By("Ticking leader until it steps down")
					Eventually(func() <-chan raft.SoftState {
						c1.clock.Increment(interval)
						return c1.observe
					}, LongEventualTimeout).Should(Receive(Equal(raft.SoftState{Lead: 0, RaftState: raft.StateFollower})))

					By("Ensuring it does not accept message due to the cluster being leaderless")
					err := c1.Order(env, 0)
					Expect(err).To(MatchError("no Raft leader"))

					network.elect(2)

					
					Expect(c1.fakeFields.fakeIsLeader.SetCallCount()).Should(Equal(3))
					Expect(c1.fakeFields.fakeIsLeader.SetArgsForCall(2)).Should(Equal(float64(0)))
					
					Expect(c2.fakeFields.fakeIsLeader.SetCallCount()).Should(Equal(2))
					Expect(c2.fakeFields.fakeIsLeader.SetArgsForCall(1)).Should(Equal(float64(1)))
					
					Expect(c3.fakeFields.fakeIsLeader.SetCallCount()).Should(Equal(1))

					network.join(1, true)
					network.exec(func(c *chain) {
						Expect(c.fakeFields.fakeLeaderChanges.AddCallCount()).Should(Equal(3))
						Expect(c.fakeFields.fakeLeaderChanges.AddArgsForCall(2)).Should(Equal(float64(1)))
					})

					err = c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
				})

				It("does not deadlock if propose is blocked", func() {
					signal := make(chan struct{})
					c1.cutter.CutNext = true
					c1.support.SequenceStub = func() uint64 {
						signal <- struct{}{}
						<-signal
						return 0
					}

					By("Sending a normal transaction")
					Expect(c1.Order(env, 0)).To(Succeed())

					Eventually(signal).Should(Receive())
					network.disconnect(1)

					By("Ticking leader till it steps down")
					Eventually(func() raft.SoftState {
						c1.clock.Increment(interval)
						return c1.Node.Status().SoftState
					}).Should(StateEqual(0, raft.StateFollower))

					close(signal)

					Eventually(c1.observe).Should(Receive(StateEqual(0, raft.StateFollower)))
					c1.support.SequenceStub = nil
					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))
					})

					By("Re-electing 1 as leader")
					network.connect(1)
					network.elect(1)

					By("Sending another normal transaction")
					Expect(c1.Order(env, 0)).To(Succeed())

					network.exec(func(c *chain) {
						Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					})
				})
			})

			When("follower is disconnected", func() {
				It("should return error when receiving an env", func() {
					network.disconnect(2)

					errorC := c2.Errored()
					Consistently(errorC).ShouldNot(BeClosed()) 

					By("Ticking node 2 until it becomes pre-candidate")
					Eventually(func() <-chan raft.SoftState {
						c2.clock.Increment(interval)
						return c2.observe
					}, LongEventualTimeout).Should(Receive(Equal(raft.SoftState{Lead: 0, RaftState: raft.StatePreCandidate})))

					Eventually(errorC).Should(BeClosed())
					err := c2.Order(env, 0)
					Expect(err).To(HaveOccurred())
					Expect(c2.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(1))
					Expect(c2.fakeFields.fakeNormalProposalsReceived.AddArgsForCall(0)).To(Equal(float64(1)))
					Expect(c1.fakeFields.fakeNormalProposalsReceived.AddCallCount()).To(Equal(0))

					network.connect(2)
					c1.clock.Increment(interval)
					Expect(errorC).To(BeClosed())

					Eventually(c2.Errored).ShouldNot(BeClosed())
				})
			})

			It("leader retransmits lost messages", func() {
				

				c1.cutter.CutNext = true

				network.disconnect(1) 

				err := c1.Order(env, 0)
				Expect(err).NotTo(HaveOccurred())

				network.exec(
					func(c *chain) {
						Consistently(func() int { return c.support.WriteBlockCallCount() }).Should(Equal(0))
					})

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
				for i := 0; i < 3; i++ {
					Expect(c1.Order(env, 0)).To(Succeed())
				}

				Consistently(c1.support.WriteBlockCallCount).Should(Equal(0))

				network.connect(1)

				
				
				
				
				
				
				
				m2 := &raftpb.Message{To: c1.id, From: c2.id, Type: raftpb.MsgHeartbeatResp}
				c1.Consensus(&orderer.ConsensusRequest{Channel: channelID, Payload: protoutil.MarshalOrPanic(m2)}, c2.id)
				m3 := &raftpb.Message{To: c1.id, From: c3.id, Type: raftpb.MsgHeartbeatResp}
				c1.Consensus(&orderer.ConsensusRequest{Channel: channelID, Payload: protoutil.MarshalOrPanic(m3)}, c3.id)

				network.exec(func(c *chain) {
					Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))
				})
			})

			It("new leader should wait for in-fight blocks to commit before accepting new env", func() {
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				c1.cutter.CutNext = true
				c2.cutter.CutNext = true

				step1 := c1.getStepFunc()
				c1.setStepFunc(func(dest uint64, msg *orderer.ConsensusRequest) error {
					stepMsg := &raftpb.Message{}
					Expect(proto.Unmarshal(msg.Payload, stepMsg)).NotTo(HaveOccurred())

					if dest == 3 {
						return nil
					}

					if stepMsg.Type == raftpb.MsgApp && len(stepMsg.Entries) == 0 {
						return nil
					}

					return step1(dest, msg)
				})

				Expect(c1.Order(env, 0)).NotTo(HaveOccurred())

				Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
				Consistently(c2.support.WriteBlockCallCount).Should(Equal(0))
				Consistently(c3.support.WriteBlockCallCount).Should(Equal(0))

				network.disconnect(1)

				step2 := c2.getStepFunc()
				c2.setStepFunc(func(dest uint64, msg *orderer.ConsensusRequest) error {
					stepMsg := &raftpb.Message{}
					Expect(proto.Unmarshal(msg.Payload, stepMsg)).NotTo(HaveOccurred())

					if stepMsg.Type == raftpb.MsgApp && len(stepMsg.Entries) != 0 && dest == 3 {
						for _, ent := range stepMsg.Entries {
							if len(ent.Data) != 0 {
								return nil
							}
						}
					}
					return step2(dest, msg)
				})

				network.elect(2)

				go func() {
					defer GinkgoRecover()
					Expect(c2.Order(env, 0)).NotTo(HaveOccurred())
				}()

				Consistently(c2.support.WriteBlockCallCount).Should(Equal(0))
				Consistently(c3.support.WriteBlockCallCount).Should(Equal(0))

				c2.setStepFunc(step2)
				c2.clock.Increment(interval)

				Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
				Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))

				b, _ := c2.support.WriteBlockArgsForCall(0)
				Expect(b.Header.Number).To(Equal(uint64(1)))
				b, _ = c2.support.WriteBlockArgsForCall(1)
				Expect(b.Header.Number).To(Equal(uint64(2)))
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
						newConfigUpdateEnv(channelID, nil, values),
					)
				})

				It("holds up block creation on leader once a config block has been created and not written out", func() {
					
					network.disconnect(1)

					c1.cutter.CutNext = true
					
					err := c1.Order(configEnv, 0)
					Expect(err).NotTo(HaveOccurred())

					
					tempEnv := env
					tempC1 := c1

					done := make(chan struct{})

					
					go func() {
						defer GinkgoRecover()

						
						err := tempC1.Order(tempEnv, 0)
						Expect(err).NotTo(HaveOccurred())

						close(done)
					}()

					Consistently(done).ShouldNot(BeClosed())

					network.connect(1)
					c1.clock.Increment(interval)

					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteConfigBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
						})

					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
						})
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
					Expect(err).NotTo(HaveOccurred())
					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
						})
				})
			})

			When("Snapshotting is enabled", func() {
				BeforeEach(func() {
					c1.opts.SnapshotIntervalSize = 1
					c1.opts.SnapshotCatchUpEntries = 1
				})

				It("keeps running if some entries in memory are purged", func() {
					
					

					i, err := c1.opts.MemoryStorage.FirstIndex()
					Expect(err).NotTo(HaveOccurred())
					Expect(i).To(Equal(uint64(1)))

					c1.cutter.CutNext = true

					err = c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(1))
						})

					Eventually(c1.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))
					i, err = c1.opts.MemoryStorage.FirstIndex()
					Expect(err).NotTo(HaveOccurred())

					err = c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(2))
						})

					Eventually(c1.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))
					i, err = c1.opts.MemoryStorage.FirstIndex()
					Expect(err).NotTo(HaveOccurred())

					err = c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					network.exec(
						func(c *chain) {
							Eventually(c.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(3))
						})

					Eventually(c1.opts.MemoryStorage.FirstIndex, LongEventualTimeout).Should(BeNumerically(">", i))
				})

				It("lagged node can catch up using snapshot", func() {
					network.disconnect(2)
					c1.cutter.CutNext = true

					c2Lasti, _ := c2.opts.MemoryStorage.LastIndex()
					var blockCnt int
					
					
					Eventually(func() bool {
						c1Firsti, _ := c1.opts.MemoryStorage.FirstIndex()
						if c1Firsti > c2Lasti+1 {
							return true
						}

						Expect(c1.Order(env, 0)).To(Succeed())
						blockCnt++
						Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(blockCnt))
						Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(blockCnt))
						return false
					}, LongEventualTimeout).Should(BeTrue())

					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))

					network.join(2, false)

					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(blockCnt))
					indices := etcdraft.ListSnapshots(logger, c2.opts.SnapDir)
					Expect(indices).To(HaveLen(1))
					gap := indices[0] - c2Lasti
					Expect(c2.puller.PullBlockCallCount()).To(Equal(int(gap)))

					
					Expect(c2.Order(env, 0)).To(Succeed())

					network.exec(
						func(c *chain) {
							Eventually(func() int { return c.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(blockCnt + 1))
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
					Expect(err).NotTo(HaveOccurred())

					
					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))

					
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

					By("order envelope on follower")
					err = c3.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					
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

					Eventually(c1.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(0))
					Eventually(c3.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(1))

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

					network.join(2, false)
					Eventually(c2.support.WriteBlockCallCount, LongEventualTimeout).Should(Equal(10))

					network.disconnect(1)
					network.elect(2)
				})

				It("purges blockcutter, stops timer and discards created blocks if leadership is lost", func() {
					
					c1.cutter.CutNext = false
					err := c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())
					Eventually(c1.cutter.CurBatch, LongEventualTimeout).Should(HaveLen(1))

					
					c1.clock.WaitForNWatchersAndIncrement(interval, 2)
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(0))

					network.disconnect(1)
					network.elect(2)
					network.join(1, true)

					Eventually(c1.clock.WatcherCount, LongEventualTimeout).Should(Equal(1)) 
					Eventually(c1.cutter.CurBatch, LongEventualTimeout).Should(HaveLen(0))
					
					Consistently(c1.support.WriteBlockCallCount).Should(Equal(0))

					network.disconnect(2)
					network.elect(1)

					err = c1.Order(env, 0)
					Expect(err).NotTo(HaveOccurred())

					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					
					

					c1.clock.WaitForNWatchersAndIncrement(timeout-interval, 2)
					Eventually(func() int { return c1.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))
					Eventually(func() int { return c3.support.WriteBlockCallCount() }, LongEventualTimeout).Should(Equal(0))

					c1.clock.Increment(interval)
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
					
					Eventually(c1.observe, LongEventualTimeout).Should(Receive(Equal(raft.SoftState{Lead: 2, RaftState: raft.StateFollower})))
				})
			})
		})
	})
})

func nodeConfigFromMetadata(consenterMetadata *raftprotos.ConfigMetadata) []cluster.RemoteNode {
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

func createMetadata(nodeCount int, tlsCA tlsgen.CA) *raftprotos.ConfigMetadata {
	md := &raftprotos.ConfigMetadata{Options: &raftprotos.Options{
		TickInterval:      time.Duration(interval).String(),
		ElectionTick:      ELECTION_TICK,
		HeartbeatTick:     HEARTBEAT_TICK,
		MaxInflightBlocks: 5,
	}}
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


type stepFunc func(dest uint64, msg *orderer.ConsensusRequest) error

type chain struct {
	id uint64

	stepLock sync.RWMutex
	step     stepFunc

	support      *consensusmocks.FakeConsenterSupport
	cutter       *mockblockcutter.Receiver
	configurator *mocks.FakeConfigurator
	rpc          *mocks.FakeRPC
	storage      *raft.MemoryStorage
	clock        *fakeclock.FakeClock
	opts         etcdraft.Options
	puller       *mocks.FakeBlockPuller

	
	ledgerLock            sync.RWMutex
	ledger                map[uint64]*common.Block
	ledgerHeight          uint64
	lastConfigBlockNumber uint64

	observe   chan raft.SoftState
	unstarted chan struct{}
	stopped   chan struct{}

	fakeFields *fakeMetricsFields

	*etcdraft.Chain
}

func newChain(timeout time.Duration, channel string, dataDir string, id uint64, raftMetadata *raftprotos.BlockMetadata, consenters map[uint64]*raftprotos.Consenter) *chain {
	rpc := &mocks.FakeRPC{}
	clock := fakeclock.NewFakeClock(time.Now())
	storage := raft.NewMemoryStorage()

	fakeFields := newFakeMetricsFields()

	opts := etcdraft.Options{
		RaftID:              uint64(id),
		Clock:               clock,
		TickInterval:        interval,
		ElectionTick:        ELECTION_TICK,
		HeartbeatTick:       HEARTBEAT_TICK,
		MaxSizePerMsg:       1024 * 1024,
		MaxInflightBlocks:   256,
		BlockMetadata:       raftMetadata,
		LeaderCheckInterval: 500 * time.Millisecond,
		Consenters:          consenters,
		Logger:              flogging.NewFabricLogger(zap.NewExample()),
		MemoryStorage:       storage,
		WALDir:              path.Join(dataDir, "wal"),
		SnapDir:             path.Join(dataDir, "snapshot"),
		Metrics:             newFakeMetrics(fakeFields),
	}

	support := &consensusmocks.FakeConsenterSupport{}
	support.ChainIDReturns(channel)
	support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

	cutter := mockblockcutter.NewReceiver()
	close(cutter.Block)
	support.BlockCutterReturns(cutter)

	
	
	
	observe := make(chan raft.SoftState, 2)

	configurator := &mocks.FakeConfigurator{}
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
		stopped:      make(chan struct{}),
		configurator: configurator,
		puller:       puller,
		ledger: map[uint64]*common.Block{
			0: getSeedBlock(), 
		},
		ledgerHeight: 1,
		fakeFields:   fakeFields,
	}

	
	
	appendNormalBlockToLedger := func(b *common.Block, meta []byte) {
		c.ledgerLock.Lock()
		defer c.ledgerLock.Unlock()

		b = proto.Clone(b).(*common.Block)
		bytes, err := proto.Marshal(&common.Metadata{Value: meta})
		Expect(err).NotTo(HaveOccurred())
		b.Metadata.Metadata[common.BlockMetadataIndex_ORDERER] = bytes

		lastConfigValue := protoutil.MarshalOrPanic(&common.LastConfig{Index: c.lastConfigBlockNumber})
		b.Metadata.Metadata[common.BlockMetadataIndex_LAST_CONFIG] = protoutil.MarshalOrPanic(&common.Metadata{
			Value: lastConfigValue,
		})

		c.ledger[b.Header.Number] = b
		if c.ledgerHeight < b.Header.Number+1 {
			c.ledgerHeight = b.Header.Number + 1
		}
	}

	
	
	appendConfigBlockToLedger := func(b *common.Block, meta []byte) {
		c.ledgerLock.Lock()
		defer c.ledgerLock.Unlock()

		b = proto.Clone(b).(*common.Block)
		bytes, err := proto.Marshal(&common.Metadata{Value: meta})
		Expect(err).NotTo(HaveOccurred())
		b.Metadata.Metadata[common.BlockMetadataIndex_ORDERER] = bytes

		c.lastConfigBlockNumber = b.Header.Number

		lastConfigValue := protoutil.MarshalOrPanic(&common.LastConfig{Index: c.lastConfigBlockNumber})
		b.Metadata.Metadata[common.BlockMetadataIndex_LAST_CONFIG] = protoutil.MarshalOrPanic(&common.Metadata{
			Value: lastConfigValue,
		})

		c.ledger[b.Header.Number] = b
		if c.ledgerHeight < b.Header.Number+1 {
			c.ledgerHeight = b.Header.Number + 1
		}
	}

	c.support.WriteBlockStub = appendNormalBlockToLedger
	c.support.WriteConfigBlockStub = appendConfigBlockToLedger

	
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
	ch, err := etcdraft.NewChain(
		c.support,
		c.opts,
		c.configurator,
		c.rpc,
		func() (etcdraft.BlockPuller, error) { return c.puller, nil },
		c.observe,
	)
	Expect(err).NotTo(HaveOccurred())
	c.Chain = ch
}

func (c *chain) start() {
	c.unstarted = nil
	c.Start()
}

func (c *chain) setStepFunc(f stepFunc) {
	c.stepLock.Lock()
	c.step = f
	c.stepLock.Unlock()
}

func (c *chain) getStepFunc() stepFunc {
	c.stepLock.RLock()
	defer c.stepLock.RUnlock()
	return c.step
}

type network struct {
	sync.RWMutex

	leader uint64
	chains map[uint64]*chain

	
	
	links map[uint64]map[uint64]bool
	
	
	connectivity map[uint64]bool
}

func (n *network) link(from []uint64, to uint64) {
	links := make(map[uint64]bool)
	for _, id := range from {
		links[id] = true
	}

	n.Lock()
	defer n.Unlock()

	n.links[to] = links
}

func (n *network) linked(from, to uint64) bool {
	n.RLock()
	defer n.RUnlock()

	return n.links[to][from]
}

func (n *network) connect(id uint64) {
	n.Lock()
	defer n.Unlock()

	n.connectivity[id] = true
}

func (n *network) disconnect(id uint64) {
	n.Lock()
	defer n.Unlock()

	n.connectivity[id] = false
}

func (n *network) connected(id uint64) bool {
	n.RLock()
	defer n.RUnlock()

	return n.connectivity[id]
}

func (n *network) addChain(c *chain) {
	n.connect(c.id) 

	c.step = func(dest uint64, msg *orderer.ConsensusRequest) error {
		if !n.linked(c.id, dest) {
			return errors.Errorf("connection refused")
		}

		if !n.connected(c.id) || !n.connected(dest) {
			return errors.Errorf("connection lost")
		}

		n.RLock()
		target := n.chains[dest]
		n.RUnlock()
		go func() {
			defer GinkgoRecover()
			target.Consensus(msg, c.id)
		}()
		return nil
	}

	c.rpc.SendConsensusStub = func(dest uint64, msg *orderer.ConsensusRequest) error {
		c.stepLock.RLock()
		defer c.stepLock.RUnlock()
		return c.step(dest, msg)
	}

	c.rpc.SendSubmitStub = func(dest uint64, msg *orderer.SubmitRequest) error {
		if !n.linked(c.id, dest) {
			return errors.Errorf("connection refused")
		}

		if !n.connected(c.id) || !n.connected(dest) {
			return errors.Errorf("connection lost")
		}

		n.RLock()
		target := n.chains[dest]
		n.RUnlock()
		go func() {
			defer GinkgoRecover()
			target.Submit(msg, c.id)
		}()
		return nil
	}

	c.puller.PullBlockStub = func(i uint64) *common.Block {
		n.RLock()
		leaderChain := n.chains[n.leader]
		n.RUnlock()

		leaderChain.ledgerLock.RLock()
		defer leaderChain.ledgerLock.RUnlock()
		block := leaderChain.ledger[i]
		return block
	}

	c.puller.HeightsByEndpointsStub = func() (map[string]uint64, error) {
		n.RLock()
		leader := n.chains[n.leader]
		n.RUnlock()

		if leader == nil {
			return nil, errors.Errorf("ledger not available")
		}

		leader.ledgerLock.RLock()
		defer leader.ledgerLock.RUnlock()
		return map[string]uint64{"leader": leader.ledgerHeight}, nil
	}

	c.configurator.ConfigureCalls(func(channel string, nodes []cluster.RemoteNode) {
		var ids []uint64
		for _, node := range nodes {
			ids = append(ids, node.ID)
		}
		n.link(ids, c.id)
	})

	n.Lock()
	defer n.Unlock()
	n.chains[c.id] = c
}

func createNetwork(timeout time.Duration, channel string, dataDir string, raftMetadata *raftprotos.BlockMetadata, consenters map[uint64]*raftprotos.Consenter) *network {
	n := &network{
		chains:       make(map[uint64]*chain),
		connectivity: make(map[uint64]bool),
		links:        make(map[uint64]map[uint64]bool),
	}

	for _, nodeID := range raftMetadata.ConsenterIds {
		dir, err := ioutil.TempDir(dataDir, fmt.Sprintf("node-%d-", nodeID))
		Expect(err).NotTo(HaveOccurred())

		m := proto.Clone(raftMetadata).(*raftprotos.BlockMetadata)
		n.addChain(newChain(timeout, channel, dir, nodeID, m, consenters))
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

	for _, id := range nodes {
		n.chains[id].start()

		
		
		
		
		
		
		
		
		
		
		
		
		Eventually(func() error {
			_, err := n.chains[id].storage.Entries(1, 1, 1)
			return err
		}, LongEventualTimeout).ShouldNot(HaveOccurred())
		Eventually(n.chains[id].WaitReady, LongEventualTimeout).ShouldNot(HaveOccurred())
	}
}

func (n *network) stop(ids ...uint64) {
	nodes := ids
	if len(nodes) == 0 {
		for i := range n.chains {
			nodes = append(nodes, i)
		}
	}

	for _, id := range nodes {
		c := n.chains[id]
		c.Halt()
		Eventually(c.Errored).Should(BeClosed())
		select {
		case <-c.stopped:
		default:
			close(c.stopped)
		}
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









func (n *network) join(id uint64, expectLeaderChange bool) {
	n.connect(id)

	n.RLock()
	leader, follower := n.chains[n.leader], n.chains[id]
	n.RUnlock()

	step := leader.getStepFunc()
	signal := make(chan struct{})
	leader.setStepFunc(func(dest uint64, msg *orderer.ConsensusRequest) error {
		if dest == id {
			
			
			select {
			case <-signal:
			default:
				close(signal)
			}
		}

		return step(dest, msg)
	})

	
	
	
	Eventually(func() <-chan struct{} {
		leader.clock.Increment(interval)
		return signal
	}, LongEventualTimeout, 100*time.Millisecond).Should(BeClosed())

	leader.setStepFunc(step)

	if expectLeaderChange {
		Eventually(follower.observe, LongEventualTimeout).Should(Receive(Equal(raft.SoftState{Lead: n.leader, RaftState: raft.StateFollower})))
	}

	
	i, err := n.chains[n.leader].opts.MemoryStorage.LastIndex()
	Expect(err).NotTo(HaveOccurred())
	Eventually(n.chains[id].opts.MemoryStorage.LastIndex, LongEventualTimeout).Should(Equal(i))
}


func (n *network) elect(id uint64) {
	n.RLock()
	candidate := n.chains[id]
	var followers []*chain
	for _, c := range n.chains {
		if c.id != id {
			followers = append(followers, c)
		}
	}
	n.RUnlock()

	
	fmt.Fprintf(GinkgoWriter, "Send artificial MsgTimeoutNow to elect node %d\n", id)
	candidate.Consensus(&orderer.ConsensusRequest{Payload: protoutil.MarshalOrPanic(&raftpb.Message{Type: raftpb.MsgTimeoutNow})}, 0)
	Eventually(candidate.observe, LongEventualTimeout).Should(Receive(StateEqual(id, raft.StateLeader)))

	
	for _, c := range followers {
		if c.id == id {
			continue
		}

		select {
		case <-c.stopped: 
		case <-c.unstarted: 
		default:
			if n.linked(c.id, id) && n.connected(c.id) {
				Eventually(c.observe, LongEventualTimeout).Should(Receive(StateEqual(id, raft.StateFollower)))
			}
		}
	}

	n.Lock()
	n.leader = id
	n.Unlock()
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

func newConfigUpdateEnv(chainID string, oldValues, newValues map[string]*common.ConfigValue) *common.ConfigUpdateEnvelope {
	return &common.ConfigUpdateEnvelope{
		ConfigUpdate: marshalOrPanic(&common.ConfigUpdate{
			ChannelId: chainID,
			ReadSet: &common.ConfigGroup{
				Groups: map[string]*common.ConfigGroup{
					"Orderer": {
						Values: oldValues,
					},
				},
			},
			WriteSet: &common.ConfigGroup{
				Groups: map[string]*common.ConfigGroup{
					"Orderer": {
						Values: newValues,
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

func StateEqual(lead uint64, state raft.StateType) types.GomegaMatcher {
	return Equal(raft.SoftState{Lead: lead, RaftState: state})
}

func BeLeader() types.GomegaMatcher {
	return &StateMatcher{expect: raft.StateLeader}
}

func BeFollower() types.GomegaMatcher {
	return &StateMatcher{expect: raft.StateFollower}
}

type StateMatcher struct {
	expect raft.StateType
}

func (stmatcher *StateMatcher) Match(actual interface{}) (success bool, err error) {
	state, ok := actual.(raft.SoftState)
	if !ok {
		return false, errors.Errorf("StateMatcher expects a raft SoftState")
	}

	return state.RaftState == stmatcher.expect, nil
}

func (stmatcher *StateMatcher) FailureMessage(actual interface{}) (message string) {
	state, ok := actual.(raft.SoftState)
	if !ok {
		return "StateMatcher expects a raft SoftState"
	}

	return fmt.Sprintf("Expected %s to be %s", state.RaftState, stmatcher.expect)
}

func (stmatcher *StateMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	state, ok := actual.(raft.SoftState)
	if !ok {
		return "StateMatcher expects a raft SoftState"
	}

	return fmt.Sprintf("Expected %s not to be %s", state.RaftState, stmatcher.expect)
}

func noOpBlockPuller() (etcdraft.BlockPuller, error) {
	bp := &mocks.FakeBlockPuller{}
	return bp, nil
}
