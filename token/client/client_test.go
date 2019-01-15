/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client_test

import (
	"crypto/rand"
	"net"
	"time"

	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/client"
	"github.com/mcc-github/blockchain/token/client/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var _ = Describe("Client", func() {
	var (
		payload      *common.Payload
		payloadBytes []byte
		envelope     *common.Envelope
		expectedTxid string

		fakeSigningIdentity *mock.SigningIdentity
		fakeProver          *mock.Prover
		fakeTxSubmitter     *mock.FabricTxSubmitter

		tokenClient *client.Client
	)

	BeforeEach(func() {
		payload = &common.Payload{Data: []byte("tx-payload")}
		payloadBytes = ProtoMarshal(payload)
		envelope = &common.Envelope{Payload: payloadBytes, Signature: []byte("tx-signature")}
		expectedTxid = "dummy-tx-id"

		fakeProver = &mock.Prover{}
		fakeProver.RequestImportReturns(payload.Data, nil) 
		fakeProver.RequestTransferReturns(payload.Data, nil)

		fakeSigningIdentity = &mock.SigningIdentity{}
		fakeSigningIdentity.SerializeReturns([]byte("creator"), nil) 
		fakeSigningIdentity.SignReturns([]byte("tx-signature"), nil) 

		fakeTxSubmitter = &mock.FabricTxSubmitter{}
		ordererStatus := common.Status_SUCCESS
		fakeTxSubmitter.SubmitReturns(&ordererStatus, true, nil)
		fakeTxSubmitter.CreateTxEnvelopeReturns(envelope, expectedTxid, nil)

		tokenClient = &client.Client{
			SigningIdentity: fakeSigningIdentity,
			Prover:          fakeProver,
			TxSubmitter:     fakeTxSubmitter,
		}
	})

	Describe("Issue", func() {
		var (
			tokensToIssue []*token.TokenToIssue
		)

		BeforeEach(func() {
			
			tokensToIssue = []*token.TokenToIssue{
				{
					Type:      "type",
					Quantity:  1,
					Recipient: []byte("alice"),
				},
			}
		})

		It("returns tx envelope and valid status", func() {
			txEnvelope, txid, ordererStatus, committed, err := tokenClient.Issue(tokensToIssue, 10*time.Second)
			Expect(err).NotTo(HaveOccurred())
			Expect(txEnvelope).To(Equal(envelope))
			Expect(txid).To(Equal(expectedTxid))
			Expect(*ordererStatus).To(Equal(common.Status_SUCCESS))
			Expect(committed).To(Equal(true))

			Expect(fakeProver.RequestImportCallCount()).To(Equal(1))
			tokens, signingIdentity := fakeProver.RequestImportArgsForCall(0)
			Expect(tokens).To(Equal(tokensToIssue))
			Expect(signingIdentity).To(Equal(fakeSigningIdentity))

			Expect(fakeTxSubmitter.CreateTxEnvelopeCallCount()).To(Equal(1))
			txBytes := fakeTxSubmitter.CreateTxEnvelopeArgsForCall(0)
			Expect(txBytes).To(Equal(payload.Data))

			Expect(fakeTxSubmitter.SubmitCallCount()).To(Equal(1))
			txEnvelope, waitTime := fakeTxSubmitter.SubmitArgsForCall(0)
			Expect(txEnvelope).To(Equal(envelope))
			Expect(waitTime).To(Equal(10 * time.Second))
		})

		Context("when prover.RequestImport fails", func() {
			BeforeEach(func() {
				fakeProver.RequestImportReturns(nil, errors.New("wild-banana"))
			})

			It("returns an error", func() {
				envelope, txid, ordererStatus, committed, err := tokenClient.Issue(tokensToIssue, 0)
				Expect(err).To(MatchError("wild-banana"))
				Expect(envelope).To(BeNil())
				Expect(txid).To(Equal(""))
				Expect(ordererStatus).To(BeNil())
				Expect(committed).To(Equal(false))

				Expect(fakeProver.RequestImportCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.CreateTxEnvelopeCallCount()).To(Equal(0))
				Expect(fakeTxSubmitter.SubmitCallCount()).To(Equal(0))
			})
		})

		Context("when TxSubmitter CreateTxEnvelope fails", func() {
			BeforeEach(func() {
				fakeTxSubmitter.CreateTxEnvelopeReturns(nil, "", errors.New("wild-banana"))
			})

			It("returns an error", func() {
				envelope, txid, ordererStatus, committed, err := tokenClient.Issue(tokensToIssue, 0)
				Expect(err).To(MatchError("wild-banana"))
				Expect(envelope).To(BeNil())
				Expect(txid).To(Equal(""))
				Expect(ordererStatus).To(BeNil())
				Expect(committed).To(Equal(false))

				Expect(fakeProver.RequestImportCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.CreateTxEnvelopeCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.SubmitCallCount()).To(Equal(0))
			})
		})

		Context("when TxSubmitter Submit fails", func() {
			BeforeEach(func() {
				status := common.Status_BAD_REQUEST
				fakeTxSubmitter.SubmitReturns(&status, false, errors.New("wild-banana"))
			})

			It("returns an error", func() {
				txEnvelope, txid, ordererStatus, committed, err := tokenClient.Issue(tokensToIssue, 0)
				Expect(err).To(MatchError("wild-banana"))
				Expect(txEnvelope).To(Equal(envelope))
				Expect(txid).To(Equal(expectedTxid))
				Expect(*ordererStatus).To(Equal(common.Status_BAD_REQUEST))
				Expect(committed).To(Equal(false))

				Expect(fakeProver.RequestImportCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.CreateTxEnvelopeCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.SubmitCallCount()).To(Equal(1))
			})
		})
	})

	Describe("Transfer", func() {
		var (
			tokenIDs       [][]byte
			transferShares []*token.RecipientTransferShare
		)

		BeforeEach(func() {
			
			tokenIDs = [][]byte{[]byte("id1"), []byte("id2")}
			transferShares = []*token.RecipientTransferShare{
				{Recipient: []byte("alice"), Quantity: 100},
				{Recipient: []byte("Bob"), Quantity: 50},
			}
		})

		It("returns tx envelope and valid status", func() {
			txEnvelope, txid, ordererStatus, committed, err := tokenClient.Transfer(tokenIDs, transferShares, 10*time.Second)
			Expect(err).NotTo(HaveOccurred())
			Expect(txEnvelope).To(Equal(envelope))
			Expect(txid).To(Equal(expectedTxid))
			Expect(*ordererStatus).To(Equal(common.Status_SUCCESS))
			Expect(committed).To(Equal(true))

			Expect(fakeProver.RequestTransferCallCount()).To(Equal(1))
			tokens, shares, signingIdentity := fakeProver.RequestTransferArgsForCall(0)
			Expect(tokens).To(Equal(tokenIDs))
			Expect(shares).To(Equal(transferShares))
			Expect(signingIdentity).To(Equal(fakeSigningIdentity))

			Expect(fakeTxSubmitter.CreateTxEnvelopeCallCount()).To(Equal(1))
			txBytes := fakeTxSubmitter.CreateTxEnvelopeArgsForCall(0)
			Expect(txBytes).To(Equal(payload.Data))

			Expect(fakeTxSubmitter.SubmitCallCount()).To(Equal(1))
			txEnvelope, waitTime := fakeTxSubmitter.SubmitArgsForCall(0)
			Expect(txEnvelope).To(Equal(envelope))
			Expect(waitTime).To(Equal(10 * time.Second))
		})

		Context("when prover.RequestTransfer fails", func() {
			BeforeEach(func() {
				fakeProver.RequestTransferReturns(nil, errors.New("wild-banana"))
			})

			It("returns an error", func() {
				envelope, txid, ordererStatus, committed, err := tokenClient.Transfer(tokenIDs, transferShares, 0)
				Expect(err).To(MatchError("wild-banana"))
				Expect(envelope).To(BeNil())
				Expect(txid).To(Equal(""))
				Expect(ordererStatus).To(BeNil())
				Expect(committed).To(Equal(false))

				Expect(fakeProver.RequestTransferCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.CreateTxEnvelopeCallCount()).To(Equal(0))
				Expect(fakeTxSubmitter.SubmitCallCount()).To(Equal(0))
			})
		})

		Context("when TxSubmitter CreateTxEnvelope fails", func() {
			BeforeEach(func() {
				fakeTxSubmitter.CreateTxEnvelopeReturns(nil, "", errors.New("wild-banana"))
			})

			It("returns an error", func() {
				envelope, txid, ordererStatus, committed, err := tokenClient.Transfer(tokenIDs, transferShares, 0)
				Expect(err).To(MatchError("wild-banana"))
				Expect(envelope).To(BeNil())
				Expect(txid).To(Equal(""))
				Expect(ordererStatus).To(BeNil())
				Expect(committed).To(Equal(false))

				Expect(fakeProver.RequestTransferCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.CreateTxEnvelopeCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.SubmitCallCount()).To(Equal(0))
			})
		})

		Context("when TxSubmitter Submit fails", func() {
			BeforeEach(func() {
				fakeTxSubmitter.SubmitReturns(nil, false, errors.New("wild-banana"))
			})

			It("returns an error", func() {
				txEnvelope, txid, ordererStatus, committed, err := tokenClient.Transfer(tokenIDs, transferShares, 0)
				Expect(err).To(MatchError("wild-banana"))
				Expect(txEnvelope).To(Equal(envelope))
				Expect(txid).To(Equal(expectedTxid))
				Expect(ordererStatus).To(BeNil())
				Expect(committed).To(Equal(false))

				Expect(fakeProver.RequestTransferCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.CreateTxEnvelopeCallCount()).To(Equal(1))
				Expect(fakeTxSubmitter.SubmitCallCount()).To(Equal(1))
			})
		})
	})

	Describe("NewClient", func() {
		var (
			config          *client.ClientConfig
			ordererListener net.Listener
			deliverListener net.Listener
			ordererServer   *grpc.Server
			deliverServer   *grpc.Server
		)

		BeforeEach(func() {
			
			var err error
			ordererListener, err = net.Listen("tcp", "127.0.0.1:")
			Expect(err).To(BeNil())

			deliverListener, err = net.Listen("tcp", "127.0.0.1:")
			Expect(err).To(BeNil())

			ordererServer = grpc.NewServer()
			go ordererServer.Serve(ordererListener)

			deliverServer = grpc.NewServer()
			go deliverServer.Serve(deliverListener)

			ordererEndpoint := ordererListener.Addr().String()
			deliverEndpoint := deliverListener.Addr().String()
			config = getClientConfig(false, "test-channel", ordererEndpoint, deliverEndpoint, "dummy_endpoint")
		})

		AfterEach(func() {
			if ordererListener != nil {
				ordererListener.Close()
			}
			if deliverListener != nil {
				deliverListener.Close()
			}
			ordererServer.Stop()
			deliverServer.Stop()
		})

		It("creates a client", func() {
			c, err := client.NewClient(*config, fakeSigningIdentity)
			Expect(err).NotTo(HaveOccurred())

			prover, ok := c.Prover.(*client.ProverPeer)
			Expect(ok).To(Equal(true))
			Expect(prover.ChannelID).To(Equal(config.ChannelID))
			Expect(prover.RandomnessReader).To(Equal(rand.Reader))

			submitter, ok := c.TxSubmitter.(*client.TxSubmitter)
			Expect(ok).To(Equal(true))
			Expect(submitter.Config).To(Equal(config))
			Expect(submitter.SigningIdentity).To(Equal(c.SigningIdentity))
		})

		Context("when channel id is missing", func() {
			BeforeEach(func() {
				config.ChannelID = ""
			})

			It("returns an error", func() {
				_, err := client.NewClient(*config, fakeSigningIdentity)
				Expect(err).To(MatchError("missing channel id"))
			})
		})

		Context("when TLS root cert file is missing", func() {
			BeforeEach(func() {
				config.Orderer.TLSEnabled = true
				config.Orderer.TLSRootCertFile = ""
			})

			It("returns an error", func() {
				_, err := client.NewClient(*config, fakeSigningIdentity)
				Expect(err).To(MatchError("missing orderer TLSRootCertFile"))
			})
		})

		Context("when it fails to load TLS root cert file", func() {
			BeforeEach(func() {
				config.ProverPeer.TLSEnabled = true
				config.ProverPeer.TLSRootCertFile = "/non-file"
			})

			It("returns an error", func() {
				_, err := client.NewClient(*config, fakeSigningIdentity)
				Expect(err.Error()).To(ContainSubstring("unable to load TLS cert from %s", config.ProverPeer.TLSRootCertFile))
			})
		})

		Context("when NewTxSumitter fails to connect to orderer", func() {
			BeforeEach(func() {
				ordererServer.Stop()
			})

			It("returns an error", func() {
				_, err := client.NewClient(*config, fakeSigningIdentity)
				Expect(err.Error()).To(ContainSubstring("failed to connect to orderer"))
			})
		})

		Context("when NewTxSumitter fails to connect to commit peer", func() {
			BeforeEach(func() {
				deliverServer.Stop()
			})

			It("returns an error", func() {
				_, err := client.NewClient(*config, fakeSigningIdentity)
				Expect(err.Error()).To(ContainSubstring("failed to connect to commit peer"))
			})
		})

		Context("when SignIdentity fails to serialize", func() {
			BeforeEach(func() {
				fakeSigningIdentity.SerializeReturns(nil, errors.New("wild-banana"))
			})

			It("returns an error", func() {
				_, err := client.NewClient(*config, fakeSigningIdentity)
				Expect(err).To(MatchError("wild-banana"))
			})
		})
	})
})
