/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode_test

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/internal/peer/lifecycle/chaincode"
	"github.com/mcc-github/blockchain/internal/peer/lifecycle/chaincode/mock"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("SimulateCommit", func() {
	Describe("CommitSimulator", func() {
		var (
			mockProposalResponse *pb.ProposalResponse
			mockEndorserClient   *mock.EndorserClient
			mockSigner           *mock.Signer
			input                *chaincode.CommitSimulationInput
			commitSimulator      *chaincode.CommitSimulator
		)

		BeforeEach(func() {
			mockEndorserClient = &mock.EndorserClient{}
			mockResult := &lb.SimulateCommitChaincodeDefinitionResult{
				Approved: map[string]bool{
					"seemsfinetome":  true,
					"well...ok":      true,
					"absolutely-not": false,
				},
			}
			mockResultBytes, err := proto.Marshal(mockResult)
			Expect(err).NotTo(HaveOccurred())
			Expect(mockResultBytes).NotTo(BeNil())
			mockProposalResponse = &pb.ProposalResponse{
				Response: &pb.Response{
					Status:  200,
					Payload: mockResultBytes,
				},
				Endorsement: &pb.Endorsement{},
			}
			mockEndorserClient.ProcessProposalReturns(mockProposalResponse, nil)

			input = &chaincode.CommitSimulationInput{
				ChannelID: "testchannel",
				Name:      "testcc",
				Version:   "1.0",
				Sequence:  1,
			}

			mockSigner = &mock.Signer{}
			buffer := gbytes.NewBuffer()

			commitSimulator = &chaincode.CommitSimulator{
				Input:          input,
				EndorserClient: mockEndorserClient,
				Signer:         mockSigner,
				Writer:         buffer,
			}
		})

		It("simulates committing a chaincode definition and writes the output as human readable plain-text", func() {
			err := commitSimulator.Simulate()
			Expect(err).NotTo(HaveOccurred())
			Eventually(commitSimulator.Writer).Should(gbytes.Say("Chaincode definition for chaincode 'testcc', version '1.0', sequence '1' on channel 'testchannel' approval status by org"))
			Eventually(commitSimulator.Writer).Should(gbytes.Say("absolutely-not: false"))
			Eventually(commitSimulator.Writer).Should(gbytes.Say("seemsfinetome: true"))
			Eventually(commitSimulator.Writer).Should(gbytes.Say("well...ok: true"))
		})

		Context("when JSON-formatted output is requested", func() {
			BeforeEach(func() {
				commitSimulator.Input.OutputFormat = "json"
			})
			It("simulates committing a chaincode definition and writes the output as JSON", func() {
				err := commitSimulator.Simulate()
				Expect(err).NotTo(HaveOccurred())
				expectedOutput := &lb.SimulateCommitChaincodeDefinitionResult{
					Approved: map[string]bool{
						"absolutely-not": false,
						"well...ok":      true,
						"seemsfinetome":  true,
					},
				}
				json, err := json.MarshalIndent(expectedOutput, "", "\t")
				Eventually(commitSimulator.Writer).Should(gbytes.Say(fmt.Sprintf("%s", string(json))))
			})
		})

		Context("when the channel name is not provided", func() {
			BeforeEach(func() {
				commitSimulator.Input.ChannelID = ""
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("The required parameter 'channelID' is empty. Rerun the command with -C flag"))
			})
		})

		Context("when the chaincode name is not provided", func() {
			BeforeEach(func() {
				commitSimulator.Input.Name = ""
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("The required parameter 'name' is empty. Rerun the command with -n flag"))
			})
		})

		Context("when the chaincode version is not provided", func() {
			BeforeEach(func() {
				commitSimulator.Input.Version = ""
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("The required parameter 'version' is empty. Rerun the command with -v flag"))
			})
		})

		Context("when the sequence is not provided", func() {
			BeforeEach(func() {
				commitSimulator.Input.Sequence = 0
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("The required parameter 'sequence' is empty. Rerun the command with --sequence flag"))
			})
		})

		Context("when the signer cannot be serialized", func() {
			BeforeEach(func() {
				mockSigner.SerializeReturns(nil, errors.New("cafe"))
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create proposal: failed to serialize identity: cafe"))
			})
		})

		Context("when the signer fails to sign the proposal", func() {
			BeforeEach(func() {
				mockSigner.SignReturns(nil, errors.New("tea"))
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create signed proposal: tea"))
			})
		})

		Context("when the endorser fails to endorse the proposal", func() {
			BeforeEach(func() {
				mockEndorserClient.ProcessProposalReturns(nil, errors.New("latte"))
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to endorse proposal: latte"))
			})
		})

		Context("when the endorser returns a nil proposal response", func() {
			BeforeEach(func() {
				mockProposalResponse = nil
				mockEndorserClient.ProcessProposalReturns(mockProposalResponse, nil)
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("received nil proposal response"))
			})
		})

		Context("when the endorser returns a proposal response with a nil response", func() {
			BeforeEach(func() {
				mockProposalResponse.Response = nil
				mockEndorserClient.ProcessProposalReturns(mockProposalResponse, nil)
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("received proposal response with nil response"))
			})
		})

		Context("when the endorser returns a non-success status", func() {
			BeforeEach(func() {
				mockProposalResponse.Response = &pb.Response{
					Status:  500,
					Message: "capuccino",
				}
				mockEndorserClient.ProcessProposalReturns(mockProposalResponse, nil)
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("query failed with status: 500 - capuccino"))
			})
		})

		Context("when the endorser returns an unexpected result", func() {
			BeforeEach(func() {
				mockProposalResponse.Response = &pb.Response{
					Status:  200,
					Payload: []byte("jibberish"),
				}
				mockEndorserClient.ProcessProposalReturns(mockProposalResponse, nil)
			})

			It("returns an error", func() {
				err := commitSimulator.Simulate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to unmarshal proposal response's response payload"))
			})
		})
	})

	Describe("SimulateCommitCmd", func() {
		var (
			simulateCommitCmd *cobra.Command
		)

		BeforeEach(func() {
			simulateCommitCmd = chaincode.SimulateCommitCmd(nil)
			simulateCommitCmd.SetArgs([]string{
				"--channelID=testchannel",
				"--name=testcc",
				"--version=testversion",
				"--sequence=1",
				"--peerAddresses=querypeer1",
				"--tlsRootCertFiles=tls1",
				"--signature-policy=AND ('Org1MSP.member','Org2MSP.member')",
			})
		})

		AfterEach(func() {
			chaincode.ResetFlags()
		})

		It("sets up the commit simulator and attempts to simulate committing the chaincode definition", func() {
			err := simulateCommitCmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to retrieve endorser client"))
		})

		Context("when the policy is invalid", func() {
			BeforeEach(func() {
				simulateCommitCmd.SetArgs([]string{
					"--signature-policy=notapolicy",
					"--channelID=testchannel",
					"--name=testcc",
					"--version=testversion",
					"--sequence=1",
					"--peerAddresses=querypeer1",
					"--tlsRootCertFiles=tls1",
				})
			})

			It("returns an error", func() {
				err := simulateCommitCmd.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid signature policy: notapolicy"))
			})
		})

		Context("when the collections config is invalid", func() {
			BeforeEach(func() {
				simulateCommitCmd.SetArgs([]string{
					"--collections-config=idontexist.json",
					"--channelID=testchannel",
					"--name=testcc",
					"--version=testversion",
					"--sequence=1",
					"--peerAddresses=querypeer1",
					"--tlsRootCertFiles=tls1",
				})
			})

			It("returns an error", func() {
				err := simulateCommitCmd.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid collection configuration in file idontexist.json: could not read file 'idontexist.json': open idontexist.json: no such file or directory"))
			})
		})
	})
})
