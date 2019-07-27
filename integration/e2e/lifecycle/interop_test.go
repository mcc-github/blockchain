/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/integration/nwo"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/internal/peer/common"
	"github.com/mcc-github/blockchain/msp"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tedsuo/ifrit"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Release interoperability", func() {
	var (
		client  *docker.Client
		tempDir string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "nwo")
		Expect(err).NotTo(HaveOccurred())

		client, err = docker.NewClientFromEnv()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("solo network", func() {
		var network *nwo.Network
		var process ifrit.Process

		BeforeEach(func() {
			soloBytes, err := ioutil.ReadFile("solo.yaml")
			Expect(err).NotTo(HaveOccurred())

			var config *nwo.Config
			err = yaml.Unmarshal(soloBytes, &config)
			Expect(err).NotTo(HaveOccurred())

			network = nwo.New(config, tempDir, client, StartPort(), components)

			
			network.GenerateConfigTree()
			network.Bootstrap()

			
			networkRunner := network.NetworkGroupRunner()
			process = ifrit.Invoke(networkRunner)
			Eventually(process.Ready(), network.EventuallyTimeout).Should(BeClosed())
		})

		AfterEach(func() {
			
			process.Signal(syscall.SIGTERM)
			Eventually(process.Wait(), network.EventuallyTimeout).Should(Receive())
			network.Cleanup()
		})

		It("deploys and executes chaincode (simple), upgrades the channel application capabilities to V2_0 and uses _lifecycle to update the endorsement policy", func() {
			By("deploying the chaincode using LSCC on a channel with V1_4 application capabilities")
			orderer := network.Orderer("orderer0")
			peer := network.Peer("org1", "peer2")

			chaincode := nwo.Chaincode{
				Name:    "mycc",
				Version: "0.0",
				Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
				Ctor:    `{"Args":["init","a","100","b","200"]}`,
				Policy:  `AND ('Org1ExampleCom.member','Org2ExampleCom.member')`,
			}

			network.CreateAndJoinChannels(orderer)
			nwo.DeployChaincodeLegacy(network, "testchannel", orderer, chaincode)
			RunQueryInvokeQuery(network, orderer, peer, 100)

			By("enabling V2_0 application capabilities")
			nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, network.Peer("org1", "peer1"), network.Peer("org2", "peer1"))

			By("ensuring that the chaincode is still operational after the upgrade")
			RunQueryInvokeQuery(network, orderer, peer, 90)

			By("restarting the network from persistence")
			RestartNetwork(&process, network)

			By("ensuring that the chaincode is still operational after the upgrade and restart")
			RunQueryInvokeQuery(network, orderer, peer, 80)

			By("attempting to invoke the chaincode without sufficient endorsements")
			sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeInvoke{
				ChannelID: "testchannel",
				Orderer:   network.OrdererAddress(orderer, nwo.ListenPort),
				Name:      "mycc",
				Ctor:      `{"Args":["invoke","a","b","10"]}`,
				PeerAddresses: []string{
					network.PeerAddress(network.Peer("org1", "peer1"), nwo.ListenPort),
				},
				WaitForEvent: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
			Expect(sess.Err).To(gbytes.Say(`\Qcommitted with status (ENDORSEMENT_POLICY_FAILURE)\E`))

			By("upgrading the chaincode definition using _lifecycle")
			chaincode = nwo.Chaincode{
				Name:            "mycc",
				Version:         "0.0",
				Path:            "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
				Lang:            "golang",
				PackageFile:     filepath.Join(tempDir, "simplecc.tar.gz"),
				SignaturePolicy: `OR ('Org1ExampleCom.member','Org2ExampleCom.member')`,
				Sequence:        "1",
				InitRequired:    false,
				Label:           "my_simple_chaincode",
			}
			nwo.DeployChaincode(network, "testchannel", orderer, chaincode)

			By("querying/invoking/querying the chaincode with the new definition")
			RunQueryInvokeQueryWithAddresses(network, orderer, peer, 70, network.PeerAddress(network.Peer("org1", "peer2"), nwo.ListenPort))

			By("restarting the network from persistence")
			RestartNetwork(&process, network)

			By("querying/invoking/querying the chaincode with the new definition again")
			RunQueryInvokeQueryWithAddresses(network, orderer, peer, 60, network.PeerAddress(network.Peer("org1", "peer2"), nwo.ListenPort))
		})

		Describe("Interoperability scenarios", func() {
			var (
				orderer              *nwo.Orderer
				peer                 *nwo.Peer
				userSigner           msp.SigningIdentity
				serialisedUserSigner []byte
				endorserClient       pb.EndorserClient
				deliveryClient       pb.DeliverClient
				ordererClient        common.BroadcastClient
			)

			BeforeEach(func() {
				orderer = network.Orderer("orderer0")
				peer = network.Peer("org1", "peer2")
				userSigner, serialisedUserSigner = Signer(network.PeerUserMSPDir(peer, "User1"))
				endorserClient = EndorserClient(
					network.PeerAddress(peer, nwo.ListenPort),
					filepath.Join(network.PeerLocalTLSDir(peer), "ca.crt"),
				)
				deliveryClient = DeliverClient(
					network.PeerAddress(peer, nwo.ListenPort),
					filepath.Join(network.PeerLocalTLSDir(peer), "ca.crt"),
				)
				ordererClient = OrdererClient(
					network.OrdererAddress(orderer, nwo.ListenPort),
					filepath.Join(network.OrdererLocalTLSDir(orderer), "ca.crt"),
				)
			})

			It("deploys a chaincode with the legacy lifecycle, invokes it and the tx is committed only after the chaincode is upgraded via _lifecycle", func() {
				By("deploying the chaincode using the legacy lifecycle")
				chaincode := nwo.Chaincode{
					Name:    "mycc",
					Version: "0.0",
					Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
					Ctor:    `{"Args":["init","a","100","b","200"]}`,
					Policy:  `OR ('Org1ExampleCom.member','Org2ExampleCom.member')`,
				}

				network.CreateAndJoinChannels(orderer)
				nwo.DeployChaincodeLegacy(network, "testchannel", orderer, chaincode)
				RunQueryInvokeQuery(network, orderer, peer, 100)

				By("invoking the chaincode with the legacy definition and keeping the transaction")
				signedProp, prop, txid := SignedProposal(
					"testchannel",
					"mycc",
					userSigner,
					serialisedUserSigner,
					"invoke",
					"a",
					"b",
					"10",
				)
				presp, err := endorserClient.ProcessProposal(
					context.Background(),
					signedProp,
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(presp).NotTo(BeNil())
				env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
				Expect(err).NotTo(HaveOccurred())
				Expect(env).NotTo(BeNil())

				By("enabling V2_0 application capabilities")
				nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, network.Peer("org1", "peer1"), network.Peer("org2", "peer1"))

				By("upgrading the chaincode definition using _lifecycle")
				chaincode = nwo.Chaincode{
					Name:            "mycc",
					Version:         "0.0",
					Path:            "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
					Lang:            "golang",
					PackageFile:     filepath.Join(tempDir, "simplecc.tar.gz"),
					SignaturePolicy: `OR ('Org1ExampleCom.member','Org2ExampleCom.member')`,
					Sequence:        "1",
					InitRequired:    false,
					Label:           "my_simple_chaincode",
				}
				nwo.DeployChaincode(network, "testchannel", orderer, chaincode)

				By("committing the old transaction, expecting to hit an MVCC conflict")
				err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("transaction invalidated with status (MVCC_READ_CONFLICT)"))
			})

			It("deploys a chaincode with the new lifecycle, invokes it and the tx is committed only after the chaincode is upgraded via _lifecycle", func() {
				By("enabling V2_0 application capabilities")
				network.CreateAndJoinChannels(orderer)
				nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, network.Peer("org1", "peer1"), network.Peer("org2", "peer1"))

				By("deploying the chaincode definition using _lifecycle")
				chaincode := nwo.Chaincode{
					Name:            "mycc",
					Version:         "0.0",
					Path:            "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
					Lang:            "golang",
					PackageFile:     filepath.Join(tempDir, "simplecc.tar.gz"),
					SignaturePolicy: `AND ('Org1ExampleCom.member','Org2ExampleCom.member')`,
					Sequence:        "1",
					InitRequired:    true,
					Label:           "my_simple_chaincode",
					Ctor:            `{"Args":["init","a","100","b","200"]}`,
				}
				nwo.DeployChaincode(network, "testchannel", orderer, chaincode)

				By("Invoking the chaincode with the first definition and keeping the transaction")
				signedProp, prop, txid := SignedProposal(
					"testchannel",
					"mycc",
					userSigner,
					serialisedUserSigner,
					"invoke",
					"a",
					"b",
					"10",
				)
				presp, err := endorserClient.ProcessProposal(
					context.Background(),
					signedProp,
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(presp).NotTo(BeNil())
				env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
				Expect(err).NotTo(HaveOccurred())
				Expect(env).NotTo(BeNil())

				By("upgrading the chaincode definition using _lifecycle")
				chaincode = nwo.Chaincode{
					Name:            "mycc",
					Version:         "0.0",
					Path:            "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
					Lang:            "golang",
					PackageFile:     filepath.Join(tempDir, "simplecc.tar.gz"),
					SignaturePolicy: `OR ('Org1ExampleCom.member','Org2ExampleCom.member')`,
					Sequence:        "2",
					InitRequired:    false,
					Label:           "my_simple_chaincode",
				}
				nwo.DeployChaincode(network, "testchannel", orderer, chaincode)

				By("committing the old transaction, expecting to hit an MVCC conflict")
				err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("transaction invalidated with status (MVCC_READ_CONFLICT)"))
			})

			Describe("Chaincode-to-chaincode interoperability", func() {
				var (
					callerDefOld nwo.Chaincode
					callerDefNew nwo.Chaincode
					calleeDefOld nwo.Chaincode
					calleeDefNew nwo.Chaincode
				)

				BeforeEach(func() {
					ccEP := `OR ('Org1ExampleCom.member','Org2ExampleCom.member')`
					callerDefOld = nwo.Chaincode{
						Name:    "caller",
						Version: "0.0",
						Path:    "github.com/mcc-github/blockchain/integration/e2e/lifecycle/chaincode/caller/cmd",
						Ctor:    `{"Args":[""]}`,
						Policy:  ccEP,
					}
					calleeDefOld = nwo.Chaincode{
						Name:    "callee",
						Version: "0.0",
						Path:    "github.com/mcc-github/blockchain/integration/e2e/lifecycle/chaincode/callee/cmd",
						Ctor:    `{"Args":[""]}`,
						Policy:  ccEP,
					}
					callerDefNew = nwo.Chaincode{
						Name:            "caller",
						Version:         "0.0",
						Path:            "github.com/mcc-github/blockchain/integration/e2e/lifecycle/chaincode/caller/cmd",
						Lang:            "golang",
						PackageFile:     filepath.Join(tempDir, "caller.tar.gz"),
						SignaturePolicy: ccEP,
						Sequence:        "1",
						Label:           "my_caller_chaincode",
						InitRequired:    true,
						Ctor:            `{"Args":[""]}`,
					}
					calleeDefNew = nwo.Chaincode{
						Name:            "callee",
						Version:         "0.0",
						Path:            "github.com/mcc-github/blockchain/integration/e2e/lifecycle/chaincode/callee/cmd",
						Lang:            "golang",
						PackageFile:     filepath.Join(tempDir, "callee.tar.gz"),
						SignaturePolicy: ccEP,
						Sequence:        "1",
						Label:           "my_callee_chaincode",
						InitRequired:    true,
						Ctor:            `{"Args":[""]}`,
					}
					By("Creating and joining the channel")
					network.CreateAndJoinChannels(orderer)
				})

				It("Deploys two chaincodes with the new lifecycle and performs a successful cc2cc invocation", func() {
					By("enabling the 2.0 capability on the channel")
					nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, network.Peer("org1", "peer1"), network.Peer("org2", "peer1"))

					By("deploying the caller chaincode using _lifecycle")
					nwo.DeployChaincode(network, "testchannel", orderer, callerDefNew)

					By("deploying the callee chaincode using _lifecycle")
					nwo.DeployChaincode(network, "testchannel", orderer, calleeDefNew)

					By("invoking the chaincode and generating a transaction")
					signedProp, prop, txid := SignedProposal(
						"testchannel",
						"caller",
						userSigner,
						serialisedUserSigner,
						"INVOKE",
						"callee",
					)
					presp, err := endorserClient.ProcessProposal(
						context.Background(),
						signedProp,
					)
					Expect(err).NotTo(HaveOccurred())
					Expect(presp).NotTo(BeNil())
					env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
					Expect(err).NotTo(HaveOccurred())
					Expect(env).NotTo(BeNil())

					By("committing the transaction")
					err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
					Expect(err).NotTo(HaveOccurred())

					By("querying the caller chaincode")
					sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
						ChannelID: "testchannel",
						Name:      "caller",
						Ctor:      `{"Args":["QUERY"]}`,
					})
					Expect(err).NotTo(HaveOccurred())
					Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
					Expect(sess).To(gbytes.Say("bar"))

					By("querying the callee chaincode")
					sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
						ChannelID: "testchannel",
						Name:      "callee",
						Ctor:      `{"Args":["QUERY"]}`,
					})
					Expect(err).NotTo(HaveOccurred())
					Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
					Expect(sess).To(gbytes.Say("bar"))
				})

				When("the network starts with new definitions", func() {
					BeforeEach(func() {
						By("enabling the 2.0 capability on the channel")
						nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, network.Peer("org1", "peer1"), network.Peer("org2", "peer1"))

						By("upgrading the caller with the new definition")
						nwo.DeployChaincode(network, "testchannel", orderer, callerDefNew)

						By("upgrading the callee with the new definition")
						nwo.DeployChaincode(network, "testchannel", orderer, calleeDefNew)
					})

					It("performs a successful cc2cc invocation", func() {
						By("invoking the chaincode and generating a transaction")
						signedProp, prop, txid := SignedProposal(
							"testchannel",
							"caller",
							userSigner,
							serialisedUserSigner,
							"INVOKE",
							"callee",
						)
						presp, err := endorserClient.ProcessProposal(
							context.Background(),
							signedProp,
						)
						Expect(err).NotTo(HaveOccurred())
						Expect(presp).NotTo(BeNil())
						env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
						Expect(err).NotTo(HaveOccurred())
						Expect(env).NotTo(BeNil())

						By("committing the transaction")
						err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
						Expect(err).NotTo(HaveOccurred())

						By("querying the caller chaincode")
						sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "caller",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("bar"))

						By("querying the callee chaincode")
						sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "callee",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("bar"))
					})

					It("performs a successful cc2cc invocation which doesn't commit because the caller is further upgraded", func() {
						By("invoking the chaincode and generating a transaction")
						signedProp, prop, txid := SignedProposal(
							"testchannel",
							"caller",
							userSigner,
							serialisedUserSigner,
							"INVOKE",
							"callee",
						)
						presp, err := endorserClient.ProcessProposal(
							context.Background(),
							signedProp,
						)
						Expect(err).NotTo(HaveOccurred())
						Expect(presp).NotTo(BeNil())
						env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
						Expect(err).NotTo(HaveOccurred())
						Expect(env).NotTo(BeNil())

						By("further upgrading the caller with the new definition")
						callerDefNew.Sequence = "2"
						callerDefNew.InitRequired = false
						nwo.DeployChaincode(network, "testchannel", orderer, callerDefNew)

						By("committing the transaction")
						err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("transaction invalidated with status (MVCC_READ_CONFLICT)"))

						By("querying the caller chaincode")
						sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "caller",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("foo"))

						By("querying the callee chaincode")
						sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "callee",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("foo"))
					})

					It("performs a successful cc2cc invocation which doesn't commit because the callee is further upgraded", func() {
						By("invoking the chaincode and generating a transaction")
						signedProp, prop, txid := SignedProposal(
							"testchannel",
							"caller",
							userSigner,
							serialisedUserSigner,
							"INVOKE",
							"callee",
						)
						presp, err := endorserClient.ProcessProposal(
							context.Background(),
							signedProp,
						)
						Expect(err).NotTo(HaveOccurred())
						Expect(presp).NotTo(BeNil())
						env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
						Expect(err).NotTo(HaveOccurred())
						Expect(env).NotTo(BeNil())

						By("further upgrading the callee with the new definition")
						calleeDefNew.Sequence = "2"
						calleeDefNew.InitRequired = false
						nwo.DeployChaincode(network, "testchannel", orderer, calleeDefNew)

						By("committing the transaction")
						err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("transaction invalidated with status (MVCC_READ_CONFLICT)"))

						By("querying the caller chaincode")
						sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "caller",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("foo"))

						By("querying the callee chaincode")
						sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "callee",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("foo"))
					})
				})

				When("the network starts with legacy definitions and then upgrades to 2.0", func() {
					BeforeEach(func() {
						By("deploying the caller chaincode using the legacy lifecycle")
						nwo.DeployChaincodeLegacy(network, "testchannel", orderer, callerDefOld)

						By("deploying the callee chaincode using the legacy lifecycle")
						nwo.DeployChaincodeLegacy(network, "testchannel", orderer, calleeDefOld)

						By("enabling the 2.0 capability on the channel")
						nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, network.Peer("org1", "peer1"), network.Peer("org2", "peer1"))
					})

					It("upgrades the caller with the new and performs a successful cc2cc invocation", func() {
						By("upgrading the caller with the new definition")
						nwo.DeployChaincode(network, "testchannel", orderer, callerDefNew)

						By("invoking the chaincode and generating a transaction")
						signedProp, prop, txid := SignedProposal(
							"testchannel",
							"caller",
							userSigner,
							serialisedUserSigner,
							"INVOKE",
							"callee",
						)
						presp, err := endorserClient.ProcessProposal(
							context.Background(),
							signedProp,
						)
						Expect(err).NotTo(HaveOccurred())
						Expect(presp).NotTo(BeNil())
						env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
						Expect(err).NotTo(HaveOccurred())
						Expect(env).NotTo(BeNil())

						By("committing the transaction")
						err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
						Expect(err).NotTo(HaveOccurred())

						By("querying the caller chaincode")
						sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "caller",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("bar"))

						By("querying the callee chaincode")
						sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "callee",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("bar"))
					})

					It("upgrades the callee with the new and performs a successful cc2cc invocation", func() {
						By("upgrading the callee with the new definition")
						nwo.DeployChaincode(network, "testchannel", orderer, calleeDefNew)

						By("invoking the chaincode and generating a transaction")
						signedProp, prop, txid := SignedProposal(
							"testchannel",
							"caller",
							userSigner,
							serialisedUserSigner,
							"INVOKE",
							"callee",
						)
						presp, err := endorserClient.ProcessProposal(
							context.Background(),
							signedProp,
						)
						Expect(err).NotTo(HaveOccurred())
						Expect(presp).NotTo(BeNil())
						env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
						Expect(err).NotTo(HaveOccurred())
						Expect(env).NotTo(BeNil())

						By("committing the transaction")
						err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
						Expect(err).NotTo(HaveOccurred())

						By("querying the caller chaincode")
						sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "caller",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("bar"))

						By("querying the callee chaincode")
						sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "callee",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("bar"))
					})

					It("performs a cc2cc invocation which fails because in the meantime, the callee is upgraded with the new lifecycle", func() {
						By("invoking the chaincode and generating a transaction")
						signedProp, prop, txid := SignedProposal(
							"testchannel",
							"caller",
							userSigner,
							serialisedUserSigner,
							"INVOKE",
							"callee",
						)
						presp, err := endorserClient.ProcessProposal(
							context.Background(),
							signedProp,
						)
						Expect(err).NotTo(HaveOccurred())
						Expect(presp).NotTo(BeNil())
						env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
						Expect(err).NotTo(HaveOccurred())
						Expect(env).NotTo(BeNil())

						By("upgrading the callee with the new definition")
						nwo.DeployChaincode(network, "testchannel", orderer, calleeDefNew)

						By("committing the transaction")
						err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("transaction invalidated with status (MVCC_READ_CONFLICT)"))

						By("querying the caller chaincode")
						sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "caller",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("foo"))

						By("querying the callee chaincode")
						sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "callee",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("foo"))
					})

					It("performs a cc2cc invocation which fails because in the meantime, the caller is upgraded with the new lifecycle", func() {
						By("invoking the chaincode and generating a transaction")
						signedProp, prop, txid := SignedProposal(
							"testchannel",
							"caller",
							userSigner,
							serialisedUserSigner,
							"INVOKE",
							"callee",
						)
						presp, err := endorserClient.ProcessProposal(
							context.Background(),
							signedProp,
						)
						Expect(err).NotTo(HaveOccurred())
						Expect(presp).NotTo(BeNil())
						env, err := protoutil.CreateSignedTx(prop, userSigner, presp)
						Expect(err).NotTo(HaveOccurred())
						Expect(env).NotTo(BeNil())

						By("upgrading the caller with the new definition")
						nwo.DeployChaincode(network, "testchannel", orderer, callerDefNew)

						By("committing the transaction")
						err = CommitTx(network, env, peer, deliveryClient, ordererClient, userSigner, txid)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("transaction invalidated with status (MVCC_READ_CONFLICT)"))

						By("querying the caller chaincode")
						sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "caller",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("foo"))

						By("querying the callee chaincode")
						sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
							ChannelID: "testchannel",
							Name:      "callee",
							Ctor:      `{"Args":["QUERY"]}`,
						})
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
						Expect(sess).To(gbytes.Say("foo"))
					})
				})
			})
		})
	})
})
