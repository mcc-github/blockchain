/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package etcdraft_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/clock/fakeclock"
	"github.com/coreos/etcd/raft"
	"github.com/mcc-github/blockchain/common/flogging"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
	consensusmocks "github.com/mcc-github/blockchain/orderer/consensus/mocks"
	mockblockcutter "github.com/mcc-github/blockchain/orderer/mocks/common/blockcutter"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var _ = Describe("Chain", func() {
	var (
		m           *common.Envelope
		normalBlock *common.Block
		interval    time.Duration
		channelID   string
	)

	BeforeEach(func() {
		channelID = "test-chain"
		m = &common.Envelope{
			Payload: utils.MarshalOrPanic(&common.Payload{
				Header: &common.Header{ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_MESSAGE), ChannelId: channelID})},
				Data:   []byte("TEST_MESSAGE"),
			}),
		}
		normalBlock = &common.Block{Data: &common.BlockData{Data: [][]byte{[]byte("foo")}}}
		interval = time.Second
	})

	Describe("Single raft node", func() {
		var (
			clock    *fakeclock.FakeClock
			opts     etcdraft.Options
			support  *consensusmocks.FakeConsenterSupport
			cutter   *mockblockcutter.Receiver
			storage  *raft.MemoryStorage
			observeC chan uint64
			chain    *etcdraft.Chain
			logger   *flogging.FabricLogger
		)

		BeforeEach(func() {
			clock = fakeclock.NewFakeClock(time.Now())
			storage = raft.NewMemoryStorage()
			logger = flogging.NewFabricLogger(zap.NewNop())
			observeC = make(chan uint64, 1)
			opts = etcdraft.Options{
				RaftID:          1,
				Clock:           clock,
				TickInterval:    interval,
				ElectionTick:    2,
				HeartbeatTick:   1,
				MaxSizePerMsg:   1024 * 1024,
				MaxInflightMsgs: 256,
				Peers:           []raft.Peer{{ID: 1}},
				Logger:          logger,
				Storage:         storage,
			}
			support = &consensusmocks.FakeConsenterSupport{}
			support.ChainIDReturns(channelID)
			support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: time.Hour})
			cutter = mockblockcutter.NewReceiver()
			support.BlockCutterReturns(cutter)

			var err error
			chain, err = etcdraft.NewChain(support, opts, observeC)
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

		Context("when no raft leader is elected", func() {
			It("fails to order envelope", func() {
				err := chain.Order(m, 0)
				Expect(err).To(MatchError("no raft leader"))
			})
		})

		Context("when raft leader is elected", func() {
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
				err := chain.Order(m, 0)
				Expect(err).To(MatchError("chain is stopped"))
			})

			It("produces blocks following batch rules", func() {
				close(cutter.Block)
				support.CreateNextBlockReturns(normalBlock)

				By("cutting next batch directly")
				cutter.CutNext = true
				err := chain.Order(m, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(support.WriteBlockCallCount).Should(Equal(1))

				By("respecting batch timeout")
				cutter.CutNext = false
				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})
				err = chain.Order(m, 0)
				Expect(err).NotTo(HaveOccurred())

				clock.WaitForNWatchersAndIncrement(timeout, 2)
				Eventually(support.WriteBlockCallCount).Should(Equal(2))
			})

			It("does not reset timer for every envelope", func() {
				close(cutter.Block)
				support.CreateNextBlockReturns(normalBlock)

				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(m, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(1))

				clock.WaitForNWatchersAndIncrement(timeout/2, 2)

				err = chain.Order(m, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(2))

				
				
				clock.Increment(timeout / 2)
				Eventually(support.WriteBlockCallCount).Should(Equal(1))
			})

			It("does not write a block if halted before timeout", func() {
				close(cutter.Block)
				timeout := time.Second
				support.SharedConfigReturns(&mockconfig.Orderer{BatchTimeoutVal: timeout})

				err := chain.Order(m, 0)
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

				err := chain.Order(m, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(1))

				clock.WaitForNWatchersAndIncrement(timeout/2, 2)

				By("force a batch to be cut before timer expires")
				cutter.CutNext = true
				err = chain.Order(m, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(support.WriteBlockCallCount).Should(Equal(1))
				Expect(support.CreateNextBlockArgsForCall(0)).To(HaveLen(2))
				Expect(cutter.CurBatch()).To(HaveLen(0))

				
				cutter.CutNext = false
				err = chain.Order(m, 0)
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

				err := chain.Order(m, 0)
				Expect(err).NotTo(HaveOccurred())
				Eventually(cutter.CurBatch).Should(HaveLen(1))

				cutter.IsolatedTx = true
				err = chain.Order(m, 0)
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

				It("enqueue if an envelope is still valid", func() {
					support.ProcessNormalMsgReturns(1, nil)

					err := chain.Order(m, 0)
					Expect(err).NotTo(HaveOccurred())
					Eventually(cutter.CurBatch).Should(HaveLen(1))
				})

				It("does not enqueue if an envelope is not valid", func() {
					support.ProcessNormalMsgReturns(1, errors.Errorf("Envelope is invalid"))

					err := chain.Order(m, 0)
					Expect(err).NotTo(HaveOccurred())
					Consistently(cutter.CurBatch).Should(HaveLen(0))
				})
			})

			It("unblocks Errored if chain is halted", func() {
				Expect(chain.Errored()).NotTo(Receive())
				chain.Halt()
				Expect(chain.Errored()).Should(BeClosed())
			})

			It("config message is not yet supported", func() {
				c := &common.Envelope{
					Payload: utils.MarshalOrPanic(&common.Payload{
						Header: &common.Header{ChannelHeader: utils.MarshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_CONFIG), ChannelId: channelID})},
						Data:   []byte("TEST_MESSAGE"),
					}),
				}

				Expect(func() {
					chain.Configure(c, 0)
				}).To(Panic())
			})
		})
	})
})
