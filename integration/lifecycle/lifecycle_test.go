/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gogo/protobuf/proto"
	"github.com/mcc-github/blockchain/common/tools/protolator"
	"github.com/mcc-github/blockchain/common/tools/protolator/protoext/ordererext"
	"github.com/mcc-github/blockchain/integration/nwo"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/protos/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tedsuo/ifrit"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Lifecycle", func() {
	var (
		client    *docker.Client
		tempDir   string
		network   *nwo.Network
		processes = map[string]ifrit.Process{}
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "nwo")
		Expect(err).NotTo(HaveOccurred())

		client, err = docker.NewClientFromEnv()
		Expect(err).NotTo(HaveOccurred())

		soloBytes, err := ioutil.ReadFile("solo.yaml")
		Expect(err).NotTo(HaveOccurred())

		var config *nwo.Config
		err = yaml.Unmarshal(soloBytes, &config)
		Expect(err).NotTo(HaveOccurred())

		network = nwo.New(config, tempDir, client, StartPort(), components)

		
		network.GenerateConfigTree()
		network.Bootstrap()

		for _, o := range network.Orderers {
			or := network.OrdererRunner(o)
			p := ifrit.Invoke(or)
			processes[o.ID()] = p
			Eventually(p.Ready(), network.EventuallyTimeout).Should(BeClosed())
		}

		for _, peer := range network.Peers {
			pr := network.PeerRunner(peer)
			p := ifrit.Invoke(pr)
			processes[peer.ID()] = p
			Eventually(p.Ready(), network.EventuallyTimeout).Should(BeClosed())
		}
	})

	AfterEach(func() {
		
		for _, p := range processes {
			p.Signal(syscall.SIGTERM)
			Eventually(p.Wait(), network.EventuallyTimeout).Should(Receive())
		}
		network.Cleanup()

		os.RemoveAll(tempDir)
	})

	It("deploys and executes chaincode using _lifecycle and upgrades it", func() {
		orderer := network.Orderer("orderer0")
		testPeers := network.PeersWithChannel("testchannel")
		org1peer2 := network.Peer("org1", "peer2")

		chaincode := nwo.Chaincode{
			Name:                "mycc",
			Version:             "0.0",
			Path:                "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Lang:                "golang",
			PackageFile:         filepath.Join(tempDir, "simplecc.tar.gz"),
			Ctor:                `{"Args":["init","a","100","b","200"]}`,
			ChannelConfigPolicy: "/Channel/Application/Endorsement",
			Sequence:            "1",
			InitRequired:        true,
			Label:               "my_simple_chaincode",
		}

		By("setting up the channel")
		network.CreateAndJoinChannels(orderer)
		network.UpdateChannelAnchors(orderer, "testchannel")
		network.VerifyMembership(network.PeersWithChannel("testchannel"), "testchannel")
		nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, network.Peer("org1", "peer1"), network.Peer("org2", "peer1"))

		By("deploying the chaincode")
		nwo.PackageChaincode(network, chaincode, testPeers[0])

		
		chaincode.SetPackageIDFromPackageFile()

		nwo.InstallChaincode(network, chaincode, testPeers...)

		By("verifying the installed chaincode package matches the one that was submitted")
		sess, err := network.PeerAdminSession(testPeers[0], commands.ChaincodeGetInstalledPackage{
			PackageID:       chaincode.PackageID,
			OutputDirectory: tempDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
		fileBytes, err := ioutil.ReadFile(chaincode.PackageFile)
		Expect(err).NotTo(HaveOccurred())
		fileBytesFromPeer, err := ioutil.ReadFile(filepath.Join(network.RootDir, chaincode.PackageID+".tar.gz"))
		Expect(err).NotTo(HaveOccurred())
		Expect(fileBytesFromPeer).To(Equal(fileBytes))

		nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, testPeers...)

		nwo.CheckCommitReadinessUntilReady(network, "testchannel", chaincode, network.PeerOrgs(), testPeers...)
		nwo.CommitChaincode(network, "testchannel", orderer, chaincode, testPeers[0], testPeers...)
		nwo.InitChaincode(network, "testchannel", orderer, chaincode, testPeers...)

		By("ensuring the chaincode can be invoked and queried")
		RunQueryInvokeQuery(network, orderer, org1peer2, 100)

		By("setting a bad package ID to temporarily disable endorsements on org1")
		savedPackageID := chaincode.PackageID
		
		
		
		chaincode.PackageID = "bad"
		nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, org1peer2)

		By("querying the chaincode and expecting the invocation to fail")
		sess, err = network.PeerUserSession(org1peer2, "User1", commands.ChaincodeQuery{
			ChannelID: "testchannel",
			Name:      "mycc",
			Ctor:      `{"Args":["query","a"]}`,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
		Expect(sess.Err).To(gbytes.Say("Error: endorsement failure during query. response: status:500 " +
			"message:\"make sure the chaincode mycc has been successfully defined on channel testchannel and try " +
			"again: chaincode definition for 'mycc' exists, but chaincode is not installed\""))

		By("setting the correct package ID to restore the chaincode")
		chaincode.PackageID = savedPackageID
		nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, org1peer2)

		By("querying the chaincode and expecting the invocation to succeed")
		sess, err = network.PeerUserSession(org1peer2, "User1", commands.ChaincodeQuery{
			ChannelID: "testchannel",
			Name:      "mycc",
			Ctor:      `{"Args":["query","a"]}`,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
		Expect(sess).To(gbytes.Say("90"))

		By("upgrading the chaincode to sequence 2")
		chaincode.Sequence = "2"

		nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, testPeers...)

		nwo.CheckCommitReadinessUntilReady(network, "testchannel", chaincode, network.PeerOrgs(), testPeers...)
		nwo.CommitChaincode(network, "testchannel", orderer, chaincode, testPeers[0], testPeers...)

		By("ensuring the chaincode can still be invoked and queried")
		RunQueryInvokeQuery(network, orderer, testPeers[0], 90)

		By("deploying another chaincode using the same chaincode package")
		nwo.DeployChaincode(network, "testchannel", orderer, nwo.Chaincode{
			Name:                "yourcc",
			Version:             "0.0",
			Path:                "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Lang:                "golang",
			PackageFile:         filepath.Join(tempDir, "simplecc.tar.gz"),
			Ctor:                `{"Args":["init","a","100","b","200"]}`,
			ChannelConfigPolicy: "/Channel/Application/Endorsement",
			Sequence:            "1",
			InitRequired:        true,
			Label:               "my_simple_chaincode",
		})

		By("listing the installed chaincodes and verifying the channel/chaincode definitions that are using the chaincode package")
		nwo.QueryInstalledReferences(network, "testchannel", chaincode.Label, chaincode.PackageID, network.Peer("org2", "peer1"), []string{"mycc", "0.0"}, []string{"yourcc", "0.0"})
		By("adding a new org")
		org3 := &nwo.Organization{
			MSPID:         "Org3ExampleCom",
			Name:          "org3",
			Domain:        "org3.example.com",
			EnableNodeOUs: true,
			Users:         2,
			CA: &nwo.CA{
				Hostname: "ca",
			},
		}

		org3peer1 := &nwo.Peer{
			Name:         "peer1",
			Organization: "org3",
			Channels:     testPeers[0].Channels,
		}
		org3peer2 := &nwo.Peer{
			Name:         "peer2",
			Organization: "org3",
			Channels:     testPeers[0].Channels,
		}
		org3Peers := []*nwo.Peer{org3peer1, org3peer2}

		network.AddOrg(org3, org3peer1, org3peer2)
		network.GenerateOrgUpdateMaterials(org3peer1, org3peer2)

		By("starting the org3 peers")
		for _, peer := range org3Peers {
			pr := network.PeerRunner(peer)
			p := ifrit.Invoke(pr)
			processes[peer.ID()] = p
			Eventually(p.Ready(), network.EventuallyTimeout).Should(BeClosed())
		}

		By("updating the channel config to include org3")
		
		currentConfig := nwo.GetConfig(network, testPeers[0], orderer, "testchannel")
		updatedConfig := proto.Clone(currentConfig).(*common.Config)

		
		sess, err = network.ConfigTxGen(commands.PrintOrg{
			ConfigPath: network.RootDir,
			PrintOrg:   "org3",
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
		org3Group := &ordererext.DynamicOrdererOrgGroup{ConfigGroup: &common.ConfigGroup{}}
		err = protolator.DeepUnmarshalJSON(bytes.NewBuffer(sess.Out.Contents()), org3Group)
		Expect(err).NotTo(HaveOccurred())

		
		updatedConfig.ChannelGroup.Groups["Application"].Groups["org3"] = org3Group.ConfigGroup
		nwo.UpdateConfig(network, orderer, "testchannel", currentConfig, updatedConfig, true, testPeers[0], testPeers...)

		By("joining the org3 peers to the channel")
		network.JoinChannel("testchannel", orderer, org3peer1, org3peer2)

		
		testPeers = network.PeersWithChannel("testchannel")

		
		
		maxLedgerHeight := nwo.GetMaxLedgerHeight(network, "testchannel", testPeers...)
		nwo.WaitUntilEqualLedgerHeight(network, "testchannel", maxLedgerHeight, testPeers...)

		By("querying definitions by org3 before performing any chaincode actions")
		sess, err = network.PeerAdminSession(network.Peer("org2", "peer1"), commands.ChaincodeListCommitted{
			ChannelID: "testchannel",
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))

		By("installing the chaincode to the org3 peers")
		nwo.InstallChaincode(network, chaincode, org3peer1, org3peer2)

		By("ensuring org3 peers do not execute the chaincode before approving the definition")
		org3AndOrg1PeerAddresses := []string{
			network.PeerAddress(org3peer1, nwo.ListenPort),
			network.PeerAddress(org1peer2, nwo.ListenPort),
		}

		sess, err = network.PeerUserSession(org3peer1, "User1", commands.ChaincodeInvoke{
			ChannelID:     "testchannel",
			Orderer:       network.OrdererAddress(orderer, nwo.ListenPort),
			Name:          "mycc",
			Ctor:          `{"Args":["invoke","a","b","10"]}`,
			PeerAddresses: org3AndOrg1PeerAddresses,
			WaitForEvent:  true,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
		Expect(sess.Err).To(gbytes.Say("chaincode definition for 'mycc' at sequence 2 on channel 'testchannel' has not yet been approved by this org"))

		By("org3 approving the chaincode definition")
		nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, network.PeersInOrg("org3")...)
		nwo.EnsureChaincodeCommitted(network, "testchannel", chaincode.Name, chaincode.Version, chaincode.Sequence, []*nwo.Organization{network.Organization("org1"), network.Organization("org2"), network.Organization("org3")}, org3peer1)

		By("ensuring chaincode can be invoked and queried by org3")
		RunQueryInvokeQueryWithAddresses(network, orderer, org3peer1, 80, org3AndOrg1PeerAddresses...)

		By("deploying a chaincode without an endorsement policy specified")
		chaincode = nwo.Chaincode{
			Name:         "defaultpolicycc",
			Version:      "0.0",
			Path:         "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Lang:         "golang",
			PackageFile:  filepath.Join(tempDir, "simplecc.tar.gz"),
			Ctor:         `{"Args":["init","a","100","b","200"]}`,
			Sequence:     "1",
			InitRequired: true,
			Label:        "my_simple_chaincode",
		}

		nwo.DeployChaincode(network, "testchannel", orderer, chaincode)

		By("attempting to invoke the chaincode without a majority")
		sess, err = network.PeerUserSession(org3peer1, "User1", commands.ChaincodeInvoke{
			ChannelID:    "testchannel",
			Orderer:      network.OrdererAddress(orderer, nwo.ListenPort),
			Name:         "defaultpolicycc",
			Ctor:         `{"Args":["invoke","a","b","10"]}`,
			WaitForEvent: true,
		})
		Expect(err).ToNot(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
		Expect(sess.Err).To(gbytes.Say(`\QError: transaction invalidated with status (ENDORSEMENT_POLICY_FAILURE)\E`))

		By("attempting to invoke the chaincode with a majority")
		sess, err = network.PeerUserSession(org3peer1, "User1", commands.ChaincodeInvoke{
			ChannelID:     "testchannel",
			Orderer:       network.OrdererAddress(orderer, nwo.ListenPort),
			Name:          "defaultpolicycc",
			Ctor:          `{"Args":["invoke","a","b","10"]}`,
			PeerAddresses: org3AndOrg1PeerAddresses,
			WaitForEvent:  true,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
		Expect(sess.Err).To(gbytes.Say(`\Qcommitted with status (VALID)\E`))
		Expect(sess.Err).To(gbytes.Say(`Chaincode invoke successful. result: status:200`))
	})
})
