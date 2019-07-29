/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/integration/nwo"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/discovery"
	pm "github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protoutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tedsuo/ifrit"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("DiscoveryService", func() {
	var (
		testDir   string
		client    *docker.Client
		network   *nwo.Network
		process   ifrit.Process
		orderer   *nwo.Orderer
		org1Peer0 *nwo.Peer
		org1Peer1 *nwo.Peer
		org2Peer0 *nwo.Peer
		org2Peer1 *nwo.Peer
		org3Peer0 *nwo.Peer
		org3Peer1 *nwo.Peer
	)

	BeforeEach(func() {
		var err error
		testDir, err = ioutil.TempDir("", "e2e-sd")
		Expect(err).NotTo(HaveOccurred())

		client, err = docker.NewClientFromEnv()
		Expect(err).NotTo(HaveOccurred())

		configBytes, err := ioutil.ReadFile(filepath.Join("testdata", "network.yaml"))
		Expect(err).NotTo(HaveOccurred())

		var networkConfig *nwo.Config
		err = yaml.Unmarshal(configBytes, &networkConfig)
		Expect(err).NotTo(HaveOccurred())

		network = nwo.New(networkConfig, testDir, client, StartPort(), components)
		network.GenerateConfigTree()
		network.Bootstrap()

		networkRunner := network.NetworkGroupRunner()
		process = ifrit.Invoke(networkRunner)
		Eventually(process.Ready()).Should(BeClosed())

		orderer = network.Orderer("orderer")
		network.CreateAndJoinChannel(orderer, "testchannel")
		network.UpdateChannelAnchors(orderer, "testchannel")

		org1Peer0 = network.Peer("org1", "peer0")
		org1Peer1 = network.Peer("org1", "peer1")
		org2Peer0 = network.Peer("org2", "peer0")
		org2Peer1 = network.Peer("org2", "peer1")
		org3Peer0 = network.Peer("org3", "peer0")
		org3Peer1 = network.Peer("org3", "peer1")
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

	It("discovers channel information", func() {
		org1Peer0 := network.Peer("org1", "peer0")
		org2Peer0 := network.Peer("org2", "peer0")
		org3Peer0 := network.Peer("org3", "peer0")

		By("discovering endorsers when missing chaincode")
		endorsers := commands.Endorsers{
			UserCert:  network.PeerUserCert(org1Peer0, "User1"),
			UserKey:   network.PeerUserKey(org1Peer0, "User1"),
			MSPID:     network.Organization(org1Peer0.Organization).MSPID,
			Server:    network.PeerAddress(org1Peer0, nwo.ListenPort),
			Channel:   "testchannel",
			Chaincode: "mycc",
		}
		sess, err := network.Discover(endorsers)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
		Expect(sess.Err).To(gbytes.Say(`failed constructing descriptor for chaincodes:<name:"mycc"`))

		By("installing and instantiating chaincode on org1.peer0")
		chaincode := nwo.Chaincode{
			Name:    "mycc",
			Version: "1.0",
			Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:    `{"Args":["init","a","100","b","200"]}`,
			Policy:  `OR (AND ('Org1MSP.member','Org2MSP.member'), AND ('Org1MSP.member','Org3MSP.member'), AND ('Org2MSP.member','Org3MSP.member'))`,
		}
		nwo.DeployChaincodeLegacy(network, "testchannel", orderer, chaincode, org1Peer0)

		By("discovering endorsers for chaincode that has not been installed to enough orgs to satisfy endorsement policy")
		sess, err = network.Discover(endorsers)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
		Expect(sess.Err).To(gbytes.Say(`failed constructing descriptor for chaincodes:<name:"mycc"`))

		By("installing chaincode to enough organizations to satisfy the endorsement policy")
		nwo.InstallChaincodeLegacy(network, chaincode, org2Peer0)

		By("discovering endorsers for chaincode that has been installed to some orgs")
		de := discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
		))
		discovered := de()
		Expect(discovered).To(HaveLen(1))
		Expect(discovered[0].Layouts).To(HaveLen(1))
		Expect(discovered[0].Layouts[0].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1)))

		By("installing chaincode to all orgs")
		nwo.InstallChaincodeLegacy(network, chaincode, org3Peer0)

		By("discovering endorsers for chaincode that has been installed to all orgs")
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org3Peer0)},
		))
		discovered = de()
		Expect(discovered).To(HaveLen(1))
		Expect(discovered[0].Layouts).To(HaveLen(3))
		Expect(discovered[0].Layouts[0].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1)))
		Expect(discovered[0].Layouts[1].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1)))
		Expect(discovered[0].Layouts[2].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1)))

		By("upgrading chaincode and adding a collections config")
		chaincode.Name = "mycc"
		chaincode.Version = "2.0"
		chaincode.CollectionsConfig = filepath.Join("testdata", "collections_config_org1_org2.json")
		nwo.UpgradeChaincodeLegacy(network, "testchannel", orderer, chaincode, org1Peer0, org2Peer0, org3Peer0)

		By("discovering endorsers for chaincode with a private collection")
		endorsers.Collection = "mycc:collectionMarbles"
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
		))
		discovered = de()
		Expect(discovered).To(HaveLen(1))
		Expect(discovered[0].Layouts).To(HaveLen(1))
		Expect(discovered[0].Layouts[0].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1)))

		By("installing chaincode to all peers")
		nwo.DeployChaincodeLegacy(network, "testchannel", orderer, nwo.Chaincode{
			Name:    "mycc2",
			Version: "1.0",
			Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:    `{"Args":["init","a","100","b","200"]}`,
			Policy:  `AND ('Org1MSP.member', 'Org2MSP.member', 'Org3MSP.member')`,
		})

		By("discovering endorsers for chaincode that has been installed to all peers")
		endorsers.Collection = ""
		endorsers.Chaincode = "mycc2"
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			ConsistOf(network.DiscoveredPeer(org1Peer0), network.DiscoveredPeer(org1Peer1)),
			ConsistOf(network.DiscoveredPeer(org2Peer0), network.DiscoveredPeer(org2Peer1)),
			ConsistOf(network.DiscoveredPeer(org3Peer0), network.DiscoveredPeer(org3Peer1)),
		))
		discovered = de()
		Expect(discovered).To(HaveLen(1))
		Expect(discovered[0].Layouts).To(HaveLen(1))
		Expect(discovered[0].Layouts[0].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1), uint32(1)))

		By("changing the channel policy")
		currentConfig := nwo.GetConfig(network, org3Peer0, orderer, "testchannel")
		updatedConfig := proto.Clone(currentConfig).(*common.Config)
		updatedConfig.ChannelGroup.Groups["Application"].Groups["org3"].Policies["Writers"].Policy.Value = protoutil.MarshalOrPanic(cauthdsl.SignedByMspAdmin("Org3MSP"))
		nwo.UpdateConfig(network, orderer, "testchannel", currentConfig, updatedConfig, true, org3Peer0)

		By("trying to discover peers as an org3 member")
		endorsers = commands.Endorsers{
			UserCert:  network.PeerUserCert(org3Peer0, "User1"),
			UserKey:   network.PeerUserKey(org3Peer0, "User1"),
			MSPID:     network.Organization(org3Peer0.Organization).MSPID,
			Server:    network.PeerAddress(org3Peer0, nwo.ListenPort),
			Channel:   "testchannel",
			Chaincode: "mycc",
		}
		sess, err = network.Discover(endorsers)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
		Expect(sess.Err).To(gbytes.Say(`access denied`))

		
		
		

		By("enabling V2_0 application capabilities on the channel")
		nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, org1Peer0, org2Peer0, org3Peer0)

		By("ensuring mycc is still discoverable after upgrading to V2_0 application capabilities")
		endorsers = commands.Endorsers{
			UserCert:  network.PeerUserCert(org1Peer0, "User1"),
			UserKey:   network.PeerUserKey(org1Peer0, "User1"),
			MSPID:     network.Organization(org1Peer0.Organization).MSPID,
			Server:    network.PeerAddress(org1Peer0, nwo.ListenPort),
			Channel:   "testchannel",
			Chaincode: "mycc",
		}
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org3Peer0)},
		))

		By("ensuring mycc's collection is still discoverable after upgrading to V2_0 application capabilities")
		endorsers.Collection = "mycc:collectionMarbles"
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
		))

		By("ensuring mycc2 is still discoverable after upgrading to V2_0 application capabilities")
		endorsers.Collection = ""
		endorsers.Chaincode = "mycc2"
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			ConsistOf(network.DiscoveredPeer(org1Peer0), network.DiscoveredPeer(org1Peer1)),
			ConsistOf(network.DiscoveredPeer(org2Peer0), network.DiscoveredPeer(org2Peer1)),
			ConsistOf(network.DiscoveredPeer(org3Peer0), network.DiscoveredPeer(org3Peer1)),
		))

		By("discovering endorsers when missing chaincode")
		endorsers = commands.Endorsers{
			UserCert:  network.PeerUserCert(org1Peer0, "User1"),
			UserKey:   network.PeerUserKey(org1Peer0, "User1"),
			MSPID:     network.Organization(org1Peer0.Organization).MSPID,
			Server:    network.PeerAddress(org1Peer0, nwo.ListenPort),
			Channel:   "testchannel",
			Chaincode: "mycc-lifecycle",
		}
		sess, err = network.Discover(endorsers)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
		Expect(sess.Err).To(gbytes.Say(`failed constructing descriptor for chaincodes:<name:"mycc-lifecycle"`))

		By("deploying chaincode using org1.peer0 and org2.peer0")
		chaincode = nwo.Chaincode{
			Name:                "mycc-lifecycle",
			Version:             "1.0",
			Lang:                "golang",
			PackageFile:         filepath.Join(testDir, "simplecc.tar.gz"),
			Path:                "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:                `{"Args":["init","a","100","b","200"]}`,
			ChannelConfigPolicy: "/Channel/Application/Endorsement",
			Sequence:            "1",
			InitRequired:        true,
			Label:               "simplecc",
		}

		By("packaging chaincode")
		nwo.PackageChaincode(network, chaincode, org1Peer0)

		
		filebytes, err := ioutil.ReadFile(chaincode.PackageFile)
		Expect(err).NotTo(HaveOccurred())
		hashStr := fmt.Sprintf("%x", util.ComputeSHA256(filebytes))
		chaincode.PackageID = chaincode.Label + ":" + hashStr

		By("installing chaincode to org1.peer0 and org2.peer0")
		nwo.InstallChaincode(network, chaincode, org1Peer0, org2Peer0)

		By("approving chaincode definition for org1 and org2")
		for _, org := range []string{"org1", "org2"} {
			nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, network.PeersInOrg(org)...)
		}

		By("committing chaincode definition using org1.peer0 and org2.peer0")
		nwo.CommitChaincode(network, "testchannel", orderer, chaincode, org1Peer0, org1Peer0, org2Peer0)
		nwo.InitChaincode(network, "testchannel", orderer, chaincode, org1Peer0, org2Peer0)

		By("discovering endorsers for chaincode that has been installed to some orgs")
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
		))
		discovered = de()
		Expect(discovered).To(HaveLen(1))
		Expect(discovered[0].Layouts).To(HaveLen(1))
		Expect(discovered[0].Layouts[0].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1)))

		By("installing chaincode to all orgs")
		nwo.InstallChaincode(network, chaincode, org3Peer0)

		By("discovering endorsers for chaincode that has been installed to all orgs but not yet approved by org3")
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
		))

		By("approving chaincode definition for org3")
		nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, network.PeersInOrg("org3")...)

		By("discovering endorsers for chaincode that has been installed and approved by all orgs")
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org3Peer0)},
		))

		By("updating the chaincode definition to sequence 2 to add a collections config")
		chaincode.Sequence = "2"
		chaincode.CollectionsConfig = filepath.Join("testdata", "collections_config_org1_org2.json")
		for _, org := range []string{"org1", "org2"} {
			nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, network.PeersInOrg(org)...)
		}

		By("committing the new chaincode definition using org1 and org2")
		nwo.CheckCommitReadinessUntilReady(network, "testchannel", chaincode, []*nwo.Organization{network.Organization("org1"), network.Organization("org2")}, org1Peer0, org2Peer0, org3Peer0)
		nwo.CommitChaincode(network, "testchannel", orderer, chaincode, org1Peer0, org1Peer0, org2Peer0)

		By("discovering endorsers for sequence 2 that has only been approved by org1 and org2")
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
		))

		By("approving the chaincode definition at sequence 2 by org3")
		maxLedgerHeight := nwo.GetMaxLedgerHeight(network, "testchannel", org1Peer0, org2Peer0, org3Peer0)
		nwo.ApproveChaincodeForMyOrg(network, "testchannel", orderer, chaincode, network.PeersInOrg("org3")...)
		nwo.WaitUntilEqualLedgerHeight(network, "testchannel", maxLedgerHeight+1, org1Peer0, org2Peer0, org3Peer0)

		By("discovering endorsers for sequence 2 that has been approved by all orgs")
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org3Peer0)},
		))

		By("discovering endorsers for chaincode with a private collection")
		endorsers.Collection = "mycc-lifecycle:collectionMarbles"
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org1Peer0)},
			[]nwo.DiscoveredPeer{network.DiscoveredPeer(org2Peer0)},
		))
		discovered = de()
		Expect(discovered).To(HaveLen(1))
		Expect(discovered[0].Layouts).To(HaveLen(1))
		Expect(discovered[0].Layouts[0].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1)))

		By("upgrading a legacy chaincode for all peers")
		nwo.DeployChaincode(network, "testchannel", orderer, nwo.Chaincode{
			Name:              "mycc",
			Version:           "2.0",
			Path:              "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			SignaturePolicy:   `AND ('Org1MSP.member', 'Org2MSP.member', 'Org3MSP.member')`,
			Lang:              "golang",
			PackageFile:       filepath.Join(testDir, "simplecc.tar.gz"),
			Sequence:          "1",
			CollectionsConfig: filepath.Join("testdata", "collections_config_org1_org2_org3.json"),
			Label:             "simplecc",
		})

		By("discovering endorsers for chaincode that has been installed to all peers")
		endorsers.Chaincode = "mycc"
		endorsers.Collection = ""
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			ConsistOf(network.DiscoveredPeer(org1Peer0), network.DiscoveredPeer(org1Peer1)),
			ConsistOf(network.DiscoveredPeer(org2Peer0), network.DiscoveredPeer(org2Peer1)),
			ConsistOf(network.DiscoveredPeer(org3Peer0), network.DiscoveredPeer(org3Peer1)),
		))
		discovered = de()
		Expect(discovered).To(HaveLen(1))
		Expect(discovered[0].Layouts).To(HaveLen(1))
		Expect(discovered[0].Layouts[0].QuantitiesByGroup).To(ConsistOf(uint32(1), uint32(1), uint32(1)))

		By("discovering endorsers for collection available on all peers")
		endorsers.Collection = "mycc:collectionMarbles"
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			ConsistOf(network.DiscoveredPeer(org1Peer0), network.DiscoveredPeer(org1Peer1)),
			ConsistOf(network.DiscoveredPeer(org2Peer0), network.DiscoveredPeer(org2Peer1)),
			ConsistOf(network.DiscoveredPeer(org3Peer0), network.DiscoveredPeer(org3Peer1)),
		))

		By("trying to discover endorsers as an org3 admin")
		endorsers = commands.Endorsers{
			UserCert:  network.PeerUserCert(org3Peer0, "Admin"),
			UserKey:   network.PeerUserKey(org3Peer0, "Admin"),
			MSPID:     network.Organization(org3Peer0.Organization).MSPID,
			Server:    network.PeerAddress(org3Peer0, nwo.ListenPort),
			Channel:   "testchannel",
			Chaincode: "mycc",
		}
		de = discoverEndorsers(network, endorsers)
		Eventually(endorsersByGroups(de), network.EventuallyTimeout).Should(ConsistOf(
			ConsistOf(network.DiscoveredPeer(org1Peer0), network.DiscoveredPeer(org1Peer1)),
			ConsistOf(network.DiscoveredPeer(org2Peer0), network.DiscoveredPeer(org2Peer1)),
			ConsistOf(network.DiscoveredPeer(org3Peer0), network.DiscoveredPeer(org3Peer1)),
		))

		By("trying to discover endorsers as an org3 member")
		endorsers = commands.Endorsers{
			UserCert:  network.PeerUserCert(org3Peer0, "User1"),
			UserKey:   network.PeerUserKey(org3Peer0, "User1"),
			MSPID:     network.Organization(org3Peer0.Organization).MSPID,
			Server:    network.PeerAddress(org3Peer0, nwo.ListenPort),
			Channel:   "testchannel",
			Chaincode: "mycc-lifecycle",
		}
		sess, err = network.Discover(endorsers)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(1))
		Expect(sess.Err).To(gbytes.Say(`access denied`))
	})

	It("discovers peer membership", func() {
		By("discovering peers")
		Eventually(nwo.DiscoverPeers(network, org1Peer0, "User1", "testchannel"), network.EventuallyTimeout).Should(ConsistOf(
			network.DiscoveredPeer(org1Peer0),
			network.DiscoveredPeer(org1Peer1),
			network.DiscoveredPeer(org2Peer0),
			network.DiscoveredPeer(org2Peer1),
			network.DiscoveredPeer(org3Peer0),
			network.DiscoveredPeer(org3Peer1),
		))

		By("installing and instantiating chaincode on a peer")
		chaincode := nwo.Chaincode{
			Name:    "mycc",
			Version: "1.0",
			Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:    `{"Args":["init","a","100","b","200"]}`,
			Policy:  `OR ('Org1MSP.member','Org2MSP.member', 'Org3MSP.member')`,
		}
		nwo.DeployChaincodeLegacy(network, "testchannel", orderer, chaincode, org1Peer0)

		By("discovering peers after installing and instantiating chaincode on a peer")
		dp := nwo.DiscoverPeers(network, org1Peer0, "User1", "testchannel")
		Eventually(peersWithChaincode(dp, "mycc"), network.EventuallyTimeout).Should(HaveLen(1))
		peersWithCC := peersWithChaincode(dp, "mycc")()
		Expect(peersWithCC).To(ConsistOf(network.DiscoveredPeer(org1Peer0, "mycc")))

		
		
		

		By("enabling V2_0 application capabilities on the channel")
		nwo.EnableCapabilities(network, "testchannel", "Application", "V2_0", orderer, org1Peer0, org2Peer0, org3Peer0)

		By("ensuring peer is still discoverable after enabling V2_0 application capabilities")
		dp = nwo.DiscoverPeers(network, org1Peer0, "User1", "testchannel")
		Eventually(peersWithChaincode(dp, "mycc"), network.EventuallyTimeout).Should(HaveLen(1))
		peersWithCC = peersWithChaincode(dp, "mycc")()
		Expect(peersWithCC).To(ConsistOf(network.DiscoveredPeer(org1Peer0, "mycc")))

		By("deploying chaincode using _lifecycle to one peer per org")
		chaincode = nwo.Chaincode{
			Name:                "mycc-lifecycle",
			Version:             "1.0",
			Lang:                "golang",
			PackageFile:         filepath.Join(testDir, "simplecc.tar.gz"),
			Path:                "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:                `{"Args":["init","a","100","b","200"]}`,
			ChannelConfigPolicy: "/Channel/Application/Endorsement",
			Sequence:            "1",
			InitRequired:        true,
			Label:               "simplecc",
		}

		nwo.DeployChaincode(network, "testchannel", orderer, chaincode, org1Peer1, org2Peer1, org3Peer1)

		By("ensuring only the peers that have installed the chaincode are discovered")
		Eventually(peersWithChaincode(dp, "mycc-lifecycle"), network.EventuallyTimeout).Should(HaveLen(3))
		peersWithCC = peersWithChaincode(dp, "mycc-lifecycle")()
		Expect(peersWithCC).To(ConsistOf(
			network.DiscoveredPeer(org1Peer1, "mycc-lifecycle"),
			network.DiscoveredPeer(org2Peer1, "mycc-lifecycle"),
			network.DiscoveredPeer(org3Peer1, "mycc-lifecycle"),
		))

		By("installing the chaincode to the remaining peers")
		nwo.PackageChaincode(network, chaincode, org1Peer0)

		
		filebytes, err := ioutil.ReadFile(chaincode.PackageFile)
		Expect(err).NotTo(HaveOccurred())
		hashStr := fmt.Sprintf("%x", util.ComputeSHA256(filebytes))
		chaincode.PackageID = chaincode.Label + ":" + hashStr

		nwo.InstallChaincode(network, chaincode, org1Peer0, org2Peer0, org3Peer0)

		By("ensuring all peers are discovered")
		Eventually(peersWithChaincode(dp, "mycc-lifecycle"), network.EventuallyTimeout).Should(HaveLen(6))
		peersWithCC = peersWithChaincode(dp, "mycc-lifecycle")()
		Expect(peersWithCC).To(ConsistOf(
			network.DiscoveredPeer(org1Peer0, "mycc-lifecycle", "mycc"),
			network.DiscoveredPeer(org1Peer1, "mycc-lifecycle"),
			network.DiscoveredPeer(org2Peer0, "mycc-lifecycle"),
			network.DiscoveredPeer(org2Peer1, "mycc-lifecycle"),
			network.DiscoveredPeer(org3Peer0, "mycc-lifecycle"),
			network.DiscoveredPeer(org3Peer1, "mycc-lifecycle"),
		))
	})

	It("discovers network configuration information", func() {
		By("retrieving the configuration")
		config := commands.Config{
			UserCert: network.PeerUserCert(org1Peer0, "User1"),
			UserKey:  network.PeerUserKey(org1Peer0, "User1"),
			MSPID:    network.Organization(org1Peer0.Organization).MSPID,
			Server:   network.PeerAddress(org1Peer0, nwo.ListenPort),
			Channel:  "testchannel",
		}
		sess, err := network.Discover(config)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))

		By("unmarshaling the response")
		discoveredConfig := &discovery.ConfigResult{}
		err = json.Unmarshal(sess.Out.Contents(), &discoveredConfig)
		Expect(err).NotTo(HaveOccurred())

		By("validating the membership data")
		Expect(discoveredConfig.Msps).To(HaveLen(len(network.Organizations)))
		for _, o := range network.Orderers {
			org := network.Organization(o.Organization)
			mspConfig, err := msp.GetVerifyingMspConfig(network.OrdererOrgMSPDir(org), org.MSPID, "bccsp")
			Expect(err).NotTo(HaveOccurred())
			Expect(discoveredConfig.Msps[org.MSPID]).To(Equal(unmarshalFabricMSPConfig(mspConfig)))
		}
		for _, p := range network.Peers {
			org := network.Organization(p.Organization)
			mspConfig, err := msp.GetVerifyingMspConfig(network.PeerOrgMSPDir(org), org.MSPID, "bccsp")
			Expect(err).NotTo(HaveOccurred())
			Expect(discoveredConfig.Msps[org.MSPID]).To(Equal(unmarshalFabricMSPConfig(mspConfig)))
		}

		By("validating the orderers")
		Expect(discoveredConfig.Orderers).To(HaveLen(len(network.Orderers)))
		for _, orderer := range network.Orderers {
			ordererMSPID := network.Organization(orderer.Organization).MSPID
			Expect(discoveredConfig.Orderers[ordererMSPID].Endpoint).To(ConsistOf(
				&discovery.Endpoint{Host: "127.0.0.1", Port: uint32(network.OrdererPort(orderer, nwo.ListenPort))},
			))
		}
	})
})

type ChaincodeEndorsers struct {
	Chaincode         string
	EndorsersByGroups map[string][]nwo.DiscoveredPeer
	Layouts           []*discovery.Layout
}

func discoverEndorsers(n *nwo.Network, command commands.Endorsers) func() []ChaincodeEndorsers {
	return func() []ChaincodeEndorsers {
		sess, err := n.Discover(command)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit())
		if sess.ExitCode() != 0 {
			return nil
		}

		discovered := []ChaincodeEndorsers{}
		err = json.Unmarshal(sess.Out.Contents(), &discovered)
		Expect(err).NotTo(HaveOccurred())
		return discovered
	}
}

func endorsersByGroups(discover func() []ChaincodeEndorsers) func() map[string][]nwo.DiscoveredPeer {
	return func() map[string][]nwo.DiscoveredPeer {
		discovered := discover()
		if len(discovered) == 1 {
			return discovered[0].EndorsersByGroups
		}
		return map[string][]nwo.DiscoveredPeer{}
	}
}

func peersWithChaincode(discover func() []nwo.DiscoveredPeer, ccName string) func() []nwo.DiscoveredPeer {
	return func() []nwo.DiscoveredPeer {
		peers := []nwo.DiscoveredPeer{}
		for _, p := range discover() {
			for _, cc := range p.Chaincodes {
				if cc == ccName {
					peers = append(peers, p)
				}
			}
		}
		return peers
	}
}

func unmarshalFabricMSPConfig(c *pm.MSPConfig) *pm.FabricMSPConfig {
	blockchainConfig := &pm.FabricMSPConfig{}
	err := proto.Unmarshal(c.Config, blockchainConfig)
	Expect(err).NotTo(HaveOccurred())
	return blockchainConfig
}
