/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package e2e

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/integration/nwo"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/protos/common"
	orderer2 "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	"github.com/mcc-github/blockchain/protos/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tedsuo/ifrit"
)

var _ = Describe("EndToEnd", func() {
	var (
		testDir   string
		client    *docker.Client
		network   *nwo.Network
		chaincode nwo.Chaincode
		process   ifrit.Process
	)

	BeforeEach(func() {
		var err error
		testDir, err = ioutil.TempDir("", "e2e")
		Expect(err).NotTo(HaveOccurred())

		client, err = docker.NewClientFromEnv()
		Expect(err).NotTo(HaveOccurred())

		chaincode = nwo.Chaincode{
			Name:    "mycc",
			Version: "0.0",
			Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:    `{"Args":["init","a","100","b","200"]}`,
			Policy:  `AND ('Org1MSP.member','Org2MSP.member')`,
		}
	})

	AfterEach(func() {
		if process != nil {
			process.Signal(syscall.SIGTERM)
			Eventually(process.Wait(), network.EventuallyTimeout).Should(Receive())
		}
		if network != nil {
			network.Cleanup()
		}
		os.RemoveAll(testDir)
	})

	Describe("basic solo network with 2 orgs", func() {
		BeforeEach(func() {
			network = nwo.New(nwo.BasicSolo(), testDir, client, 30000, components)
			network.GenerateConfigTree()
			network.Bootstrap()

			networkRunner := network.NetworkGroupRunner()
			process = ifrit.Invoke(networkRunner)
			Eventually(process.Ready()).Should(BeClosed())
		})

		It("executes a basic solo network with 2 orgs", func() {
			By("getting the orderer by name")
			orderer := network.Orderer("orderer")

			By("setting up the channel")
			network.CreateAndJoinChannel(orderer, "testchannel")

			By("deploying the chaincode")
			nwo.DeployChaincode(network, "testchannel", orderer, chaincode)

			By("getting the client peer by name")
			peer := network.Peer("Org1", "peer1")

			RunQueryInvokeQuery(network, orderer, peer)
			RunRespondWith(network, orderer, peer)
		})
	})

	Describe("basic kafka network with 2 orgs", func() {
		BeforeEach(func() {
			network = nwo.New(nwo.BasicKafka(), testDir, client, 31000, components)
			network.GenerateConfigTree()
			network.Bootstrap()

			networkRunner := network.NetworkGroupRunner()
			process = ifrit.Invoke(networkRunner)
			Eventually(process.Ready()).Should(BeClosed())
		})

		It("executes a basic kafka network with 2 orgs", func() {
			orderer := network.Orderer("orderer")
			peer := network.Peer("Org1", "peer1")

			network.CreateAndJoinChannel(orderer, "testchannel")
			nwo.DeployChaincode(network, "testchannel", orderer, chaincode)
			RunQueryInvokeQuery(network, orderer, peer)
		})
	})

	Describe("basic single node etcdraft network with 2 orgs", func() {
		BeforeEach(func() {
			network = nwo.New(nwo.BasicEtcdRaft(), testDir, client, 32000, components)
			network.GenerateConfigTree()
			network.Bootstrap()

			networkRunner := network.NetworkGroupRunner()
			process = ifrit.Invoke(networkRunner)
			Eventually(process.Ready()).Should(BeClosed())
		})

		It("executes a basic etcdraft network with 2 orgs and a single node", func() {
			orderer := network.Orderer("orderer")
			peer := network.Peer("Org1", "peer1")

			network.CreateAndJoinChannel(orderer, "testchannel")
			nwo.DeployChaincode(network, "testchannel", orderer, chaincode)
			RunQueryInvokeQuery(network, orderer, peer)
		})
	})

	Describe("three node etcdraft network with 2 orgs", func() {
		BeforeEach(func() {
			network = nwo.New(nwo.MultiNodeEtcdRaft(), testDir, client, 33000, components)
			network.GenerateConfigTree()
			network.Bootstrap()

			networkRunner := network.NetworkGroupRunner()
			process = ifrit.Invoke(networkRunner)
			Eventually(process.Ready()).Should(BeClosed())
		})

		
		
		
		
		
		It("executes an etcdraft network with 2 orgs and three orderer nodes", func() {
			orderer1 := network.Orderer("orderer1")
			orderer2 := network.Orderer("orderer2")
			orderer3 := network.Orderer("orderer3")
			peer := network.Peer("Org1", "peer1")
			org1Peer0 := network.Peer("Org1", "peer0")
			blockFile1 := filepath.Join(testDir, "newest_orderer1_block.pb")
			blockFile2 := filepath.Join(testDir, "newest_orderer2_block.pb")
			blockFile3 := filepath.Join(testDir, "newest_orderer3_block.pb")

			fetchLatestBlock := func(targetOrderer *nwo.Orderer, blockFile string) {
				c := commands.ChannelFetch{
					ChannelID:  "testchannel",
					Block:      "newest",
					OutputFile: blockFile,
				}
				if targetOrderer != nil {
					c.Orderer = network.OrdererAddress(targetOrderer, nwo.ListenPort)
				}
				sess, err := network.PeerAdminSession(org1Peer0, c)
				Expect(err).NotTo(HaveOccurred())
				Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
			}

			By("creating a new chain and having the peers join it to test for channel creation")
			network.CreateAndJoinChannel(orderer1, "testchannel")
			nwo.DeployChaincode(network, "testchannel", orderer1, chaincode)
			RunQueryInvokeQuery(network, orderer1, peer)

			
			
			
			By("fetching the latest blocks from all the orderer nodes and testing them for equality")
			fetchLatestBlock(orderer1, blockFile1)
			fetchLatestBlock(orderer2, blockFile2)
			fetchLatestBlock(orderer3, blockFile3)
			b1 := nwo.UnmarshalBlockFromFile(blockFile1)
			b2 := nwo.UnmarshalBlockFromFile(blockFile2)
			b3 := nwo.UnmarshalBlockFromFile(blockFile3)
			Expect(bytes.Equal(b1.Header.Bytes(), b2.Header.Bytes())).To(BeTrue())
			Expect(bytes.Equal(b2.Header.Bytes(), b3.Header.Bytes())).To(BeTrue())

			By("updating ACL policies to test for type A configuration updates")
			invokeChaincode := commands.ChaincodeInvoke{
				ChannelID:    "testchannel",
				Orderer:      network.OrdererAddress(orderer1, nwo.ListenPort),
				Name:         chaincode.Name,
				Ctor:         `{"Args":["invoke","a","b","10"]}`,
				WaitForEvent: true,
			}
			
			policyName := resources.Event_FilteredBlock
			policy := "/Channel/Application/org2/Admins"
			SetACLPolicy(network, "testchannel", policyName, policy, "orderer1")
			
			sess, err := network.PeerAdminSession(org1Peer0, invokeChaincode)
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess.Err, network.EventuallyTimeout).Should(gbytes.Say(`\Qdeliver completed with status (FORBIDDEN)\E`))
		})
	})

	Describe("etcd raft, checking valid configuration update of type B", func() {
		BeforeEach(func() {
			network = nwo.New(nwo.BasicEtcdRaft(), testDir, client, 32000, components)
			network.GenerateConfigTree()
			network.Bootstrap()

			networkRunner := network.NetworkGroupRunner()
			process = ifrit.Invoke(networkRunner)
			Eventually(process.Ready()).Should(BeClosed())
		})

		It("executes a basic etcdraft network with a single Raft node", func() {
			orderer := network.Orderer("orderer")
			peer := network.Peer("Org1", "peer1")

			channel := "testchannel"
			network.CreateAndJoinChannel(orderer, channel)
			nwo.DeployChaincode(network, "testchannel", orderer, chaincode)
			RunQueryInvokeQuery(network, orderer, peer)

			config := nwo.GetConfigBlock(network, peer, orderer, channel)
			updatedConfig := proto.Clone(config).(*common.Config)

			consensusTypeConfigValue := updatedConfig.ChannelGroup.Groups["Orderer"].Values["ConsensusType"]
			consensusTypeValue := &orderer2.ConsensusType{}
			err := proto.Unmarshal(consensusTypeConfigValue.Value, consensusTypeValue)
			Expect(err).NotTo(HaveOccurred())

			metadata := &etcdraft.Metadata{}
			err = proto.Unmarshal(consensusTypeValue.Metadata, metadata)
			Expect(err).NotTo(HaveOccurred())

			
			metadata.Options.MaxInflightMsgs = 1000
			metadata.Options.MaxSizePerMsg = 512

			
			consensusTypeValue.Metadata, err = proto.Marshal(metadata)
			Expect(err).NotTo(HaveOccurred())

			updatedConfig.ChannelGroup.Groups["Orderer"].Values["ConsensusType"] = &common.ConfigValue{
				ModPolicy: "Admins",
				Value:     utils.MarshalOrPanic(consensusTypeValue),
			}

			nwo.UpdateOrdererConfig(network, orderer, channel, config, updatedConfig, peer, orderer)
		})
	})

})

func RunQueryInvokeQuery(n *nwo.Network, orderer *nwo.Orderer, peer *nwo.Peer) {
	By("querying the chaincode")
	sess, err := n.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
		ChannelID: "testchannel",
		Name:      "mycc",
		Ctor:      `{"Args":["query","a"]}`,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess).To(gbytes.Say("100"))

	sess, err = n.PeerUserSession(peer, "User1", commands.ChaincodeInvoke{
		ChannelID: "testchannel",
		Orderer:   n.OrdererAddress(orderer, nwo.ListenPort),
		Name:      "mycc",
		Ctor:      `{"Args":["invoke","a","b","10"]}`,
		PeerAddresses: []string{
			n.PeerAddress(n.Peer("Org1", "peer0"), nwo.ListenPort),
			n.PeerAddress(n.Peer("Org2", "peer1"), nwo.ListenPort),
		},
		WaitForEvent: true,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess.Err).To(gbytes.Say("Chaincode invoke successful. result: status:200"))

	sess, err = n.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
		ChannelID: "testchannel",
		Name:      "mycc",
		Ctor:      `{"Args":["query","a"]}`,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess).To(gbytes.Say("90"))
}

func RunRespondWith(n *nwo.Network, orderer *nwo.Orderer, peer *nwo.Peer) {
	By("responding with a 300")
	sess, err := n.PeerUserSession(peer, "User1", commands.ChaincodeInvoke{
		ChannelID: "testchannel",
		Orderer:   n.OrdererAddress(orderer, nwo.ListenPort),
		Name:      "mycc",
		Ctor:      `{"Args":["respond","300","response-message","response-payload"]}`,
		PeerAddresses: []string{
			n.PeerAddress(n.Peer("Org2", "peer1"), nwo.ListenPort),
		},
		WaitForEvent: true,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, time.Minute).Should(gexec.Exit(0))
	Expect(sess.Err).To(gbytes.Say("Chaincode invoke successful. result: status:300"))

	By("responding with a 400")
	sess, err = n.PeerUserSession(peer, "User1", commands.ChaincodeInvoke{
		ChannelID: "testchannel",
		Orderer:   n.OrdererAddress(orderer, nwo.ListenPort),
		Name:      "mycc",
		Ctor:      `{"Args":["respond","400","response-message","response-payload"]}`,
		PeerAddresses: []string{
			n.PeerAddress(n.Peer("Org2", "peer1"), nwo.ListenPort),
		},
		WaitForEvent: true,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, time.Minute).Should(gexec.Exit(1))
	Expect(sess.Err).To(gbytes.Say(`Error: endorsement failure during invoke.`))
}
