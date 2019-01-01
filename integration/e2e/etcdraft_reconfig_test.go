/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/integration/nwo"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	"github.com/mcc-github/blockchain/protos/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
	"github.com/tedsuo/ifrit/grouper"
)

var _ = Describe("EndToEnd reconfiguration and onboarding", func() {
	var (
		testDir string
		client  *docker.Client
		network *nwo.Network
		mycc    nwo.Chaincode
		mycc2   nwo.Chaincode
		mycc3   nwo.Chaincode
		peer    *nwo.Peer

		peerProcesses    ifrit.Process
		ordererProcesses []ifrit.Process
		ordererRunners   []*ginkgomon.Runner
	)

	BeforeEach(func() {
		var err error
		testDir, err = ioutil.TempDir("", "e2e-etcfraft_reconfig")
		Expect(err).NotTo(HaveOccurred())

		client, err = docker.NewClientFromEnv()
		Expect(err).NotTo(HaveOccurred())

		mycc = nwo.Chaincode{
			Name:    "mycc",
			Version: "0.0",
			Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:    `{"Args":["init","a","100","b","200"]}`,
			Policy:  `OR ('Org1MSP.member','Org2MSP.member')`,
		}
		mycc2 = nwo.Chaincode{
			Name:    "mycc2",
			Version: "0.0",
			Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:    `{"Args":["init","a","100","b","200"]}`,
			Policy:  `OR ('Org1MSP.member','Org2MSP.member')`,
		}
		mycc3 = nwo.Chaincode{
			Name:    "mycc3",
			Version: "0.0",
			Path:    "github.com/mcc-github/blockchain/integration/chaincode/simple/cmd",
			Ctor:    `{"Args":["init","a","100","b","200"]}`,
			Policy:  `OR ('Org1MSP.member','Org2MSP.member')`,
		}
	})

	AfterEach(func() {
		if peerProcesses != nil {
			peerProcesses.Signal(syscall.SIGTERM)
			Eventually(peerProcesses.Wait(), network.EventuallyTimeout).Should(Receive())
		}
		if network != nil {
			network.Cleanup()
		}
		for _, ordererInstance := range ordererProcesses {
			ordererInstance.Signal(syscall.SIGTERM)
			Eventually(ordererInstance.Wait(), network.EventuallyTimeout).Should(Receive())
		}
		os.RemoveAll(testDir)
	})

	When("the orderer certificates are all rotated", func() {
		It("is still possible to onboard new orderers", func() {
			
			
			

			
			
			
			

			
			
			
			
			
			
			
			
			
			
			
			
			
			
			

			layout := nwo.MultiNodeEtcdRaft()
			layout.Channels = append(layout.Channels, &nwo.Channel{
				Name:    "testchannel2",
				Profile: "TwoOrgsChannel",
			}, &nwo.Channel{
				Name:    "testchannel3",
				Profile: "TwoOrgsChannel",
			})
			network = nwo.New(layout, testDir, client, BasePort(), components)
			o1, o2, o3 := network.Orderer("orderer1"), network.Orderer("orderer2"), network.Orderer("orderer3")
			orderers := []*nwo.Orderer{o1, o2, o3}

			peer = network.Peer("Org1", "peer1")

			network.GenerateConfigTree()
			network.Bootstrap()

			By("Launching the orderers")
			for _, o := range orderers {
				runner := network.OrdererRunner(o)
				ordererRunners = append(ordererRunners, runner)
				process := ifrit.Invoke(runner)
				ordererProcesses = append(ordererProcesses, process)
			}

			for _, ordererProc := range ordererProcesses {
				Eventually(ordererProc.Ready()).Should(BeClosed())
			}

			By("Launching the peers")
			peerGroup := network.PeerGroupRunner()
			peerProcesses = ifrit.Invoke(peerGroup)
			Eventually(peerProcesses.Ready()).Should(BeClosed())

			By("Checking that all orderers are online")
			assertBlockReception(map[string]int{
				"systemchannel": 0,
			}, orderers, peer, network)

			By("Creating a channel and checking that all orderers got the channel creation")
			network.CreateChannel("testchannel", network.Orderers[0], peer)
			assertBlockReception(map[string]int{
				"systemchannel": 1,
				"testchannel":   0,
			}, orderers, peer, network)

			By("Preparing new certificates for the orderer nodes")
			certificateRotations := refreshOrdererPEMs(network)

			expectedBlockHeightsPerChannel := []map[string]int{
				{"systemchannel": 2, "testchannel": 1},
				{"systemchannel": 3, "testchannel": 2},
				{"systemchannel": 4, "testchannel": 3},
				{"systemchannel": 5, "testchannel": 4},
				{"systemchannel": 6, "testchannel": 5},
				{"systemchannel": 7, "testchannel": 6},
			}

			for i, rotation := range certificateRotations {
				o := network.Orderers[i]
				port := network.OrdererPort(o, nwo.ListenPort)

				By(fmt.Sprintf("Adding the future certificate of orderer node %d", i))
				for _, channelName := range []string{"systemchannel", "testchannel"} {
					nwo.AddConsenter(network, peer, o, channelName, etcdraft.Consenter{
						ServerTlsCert: rotation.newCert,
						ClientTlsCert: rotation.newCert,
						Host:          "127.0.0.1",
						Port:          uint32(port),
					})
				}

				By("Waiting for all orderers to sync")
				assertBlockReception(expectedBlockHeightsPerChannel[i*2], orderers, peer, network)

				By("Killing the orderer")
				ordererProcesses[i].Signal(syscall.SIGTERM)
				Eventually(ordererProcesses[i].Wait(), network.EventuallyTimeout).Should(Receive())

				By("Starting the orderer again")
				ordererRunner := network.OrdererRunner(orderers[i])
				ordererRunners = append(ordererRunners, ordererRunner)
				ordererProcesses[i] = ifrit.Invoke(grouper.Member{Name: orderers[i].ID(), Runner: ordererRunner})
				Eventually(ordererProcesses[i].Ready()).Should(BeClosed())

				By("And waiting for it to stabilize")
				assertBlockReception(expectedBlockHeightsPerChannel[i*2], orderers, peer, network)

				By("Removing the previous certificate of the old orderer")
				for _, channelName := range []string{"systemchannel", "testchannel"} {
					nwo.RemoveConsenter(network, peer, network.Orderers[(i+1)%len(network.Orderers)], channelName, rotation.oldCert)
				}

				By("Waiting for all orderers to sync")
				assertBlockReception(expectedBlockHeightsPerChannel[i*2+1], orderers, peer, network)
			}

			By("Creating testchannel2")
			network.CreateChannel("testchannel2", network.Orderers[0], peer)
			assertBlockReception(map[string]int{
				"systemchannel": 8,
			}, orderers, peer, network)

			By("Creating testchannel3")
			network.CreateChannel("testchannel3", network.Orderers[0], peer)
			assertBlockReception(map[string]int{
				"systemchannel": 9,
			}, orderers, peer, network)

			o4 := &nwo.Orderer{
				Name:         "orderer4",
				Organization: "OrdererOrg",
			}

			By("Configuring orderer4 in the network")
			ports := nwo.Ports{}
			for _, portName := range nwo.OrdererPortNames() {
				ports[portName] = network.ReservePort()
			}
			network.PortsByOrdererID[o4.ID()] = ports

			network.Orderers = append(network.Orderers, o4)
			network.GenerateOrdererConfig(network.Orderer("orderer4"))

			By("Adding orderer4 to the channels")
			orderer4CertificatePath := filepath.Join(testDir, "crypto", "ordererOrganizations", "example.com",
				"orderers", "orderer4.example.com", "tls", "server.crt")
			orderer4Certificate, err := ioutil.ReadFile(orderer4CertificatePath)
			Expect(err).NotTo(HaveOccurred())
			for _, channel := range []string{"systemchannel", "testchannel"} {
				nwo.AddConsenter(network, peer, o1, channel, etcdraft.Consenter{
					ServerTlsCert: orderer4Certificate,
					ClientTlsCert: orderer4Certificate,
					Host:          "127.0.0.1",
					Port:          uint32(network.OrdererPort(o4, nwo.ListenPort)),
				})
			}

			By("Ensuring all orderers know about orderer4's addition")
			assertBlockReception(map[string]int{
				"systemchannel": 10,
				"testchannel":   7,
			}, orderers, peer, network)

			By("Joining the peer to testchannel")
			network.JoinChannel("testchannel", o1, peer)
			By("Joining the peer to testchannel2")
			network.JoinChannel("testchannel2", o1, peer)
			By("Joining the peer to testchannel3")
			network.JoinChannel("testchannel3", o1, peer)

			By("Deploying mycc and mycc2 and mycc3 to testchannel and testchannel2 and testchannel3")
			deployChaincodes(network, peer, o2, mycc, mycc2, mycc3)

			By("Waiting for orderers to sync")
			assertBlockReception(map[string]int{
				"testchannel": 8,
			}, orderers, peer, network)

			By("Transacting on testchannel once more")
			assertInvoke(network, peer, o1, mycc.Name, "testchannel", "Chaincode invoke successful. result: status:200", 0)

			assertBlockReception(map[string]int{
				"testchannel": 9,
			}, orderers, peer, network)

			By("Corrupting the readers policy of testchannel3")
			revokeReaderAccess(network, "testchannel3", o3, peer)

			
			configBlock := nwo.GetConfigBlock(network, peer, o1, "systemchannel")
			
			err = ioutil.WriteFile(filepath.Join(testDir, "systemchannel_block.pb"), utils.MarshalOrPanic(configBlock), 06444)
			Expect(err).NotTo(HaveOccurred())

			By("Launching orderer4")
			orderers = append(orderers, o4)
			orderer4Runner := network.OrdererRunner(o4)
			ordererRunners = append(ordererRunners, orderer4Runner)
			
			o4process := ifrit.Invoke(grouper.Member{Name: o4.ID(), Runner: orderer4Runner})
			Eventually(o4process.Ready()).Should(BeClosed())
			ordererProcesses = append(ordererProcesses, o4process)

			By("And waiting for it to sync with the rest of the orderers")
			assertBlockReception(map[string]int{
				"systemchannel": 10,
				"testchannel":   9,
			}, orderers, peer, network)

			By("Ensuring orderer4 doesn't serve testchannel2 and testchannel3")
			assertInvoke(network, peer, o4, mycc2.Name, "testchannel2", "channel testchannel2 is not serviced by me", 1)
			assertInvoke(network, peer, o4, mycc3.Name, "testchannel3", "channel testchannel3 is not serviced by me", 1)
			Expect(string(orderer4Runner.Err().Contents())).To(ContainSubstring("I do not belong to channel testchannel2 or am forbidden pulling it (not in the channel), skipping chain retrieval"))
			Expect(string(orderer4Runner.Err().Contents())).To(ContainSubstring("I do not belong to channel testchannel3 or am forbidden pulling it (forbidden), skipping chain retrieval"))
			By("Ensuring that all orderers don't log errors to the log")
			assertNoErrorsAreLogged(ordererRunners)
		})
	})
})

var extendedCryptoConfig = `---
OrdererOrgs:
- Name: OrdererOrg
  Domain: example.com
  EnableNodeOUs: false
  CA:
    Hostname: ca
  Specs:
  - Hostname: orderer1
    SANS:
    - localhost
    - 127.0.0.1
    - ::1
  - Hostname: orderer1new
    SANS:
    - localhost
    - 127.0.0.1
    - ::1
  - Hostname: orderer2
    SANS:
    - localhost
    - 127.0.0.1
    - ::1
  - Hostname: orderer2new
    SANS:
    - localhost
    - 127.0.0.1
    - ::1
  - Hostname: orderer3
    SANS:
    - localhost
    - 127.0.0.1
    - ::1
  - Hostname: orderer3new
    SANS:
    - localhost
    - 127.0.0.1
    - ::1
  - Hostname: orderer4
    SANS:
    - localhost
    - 127.0.0.1
    - ::1
`

type certificateChange struct {
	srcFile string
	dstFile string
	oldCert []byte
	newCert []byte
}



func refreshOrdererPEMs(n *nwo.Network) []*certificateChange {
	var fileChanges []*certificateChange
	
	cryptoConfigYAML, err := ioutil.TempFile("", "crypto-config.yaml")
	Expect(err).NotTo(HaveOccurred())
	defer os.Remove(cryptoConfigYAML.Name())

	err = ioutil.WriteFile(cryptoConfigYAML.Name(), []byte(extendedCryptoConfig), 0644)
	Expect(err).NotTo(HaveOccurred())

	
	sess, err := n.Cryptogen(commands.Extend{
		Config: cryptoConfigYAML.Name(),
		Input:  n.CryptoPath(),
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

	
	filepath.Walk(filepath.Join(n.RootDir, "crypto"), func(path string, info os.FileInfo, err error) error {
		if !strings.Contains(path, "/tls/") {
			return nil
		}
		if strings.Contains(path, "new") {
			fileChanges = append(fileChanges, &certificateChange{
				srcFile: path,
				dstFile: strings.Replace(path, "new", "", -1),
			})
		}
		return nil
	})

	var serverCertChanges []*certificateChange

	
	for _, certChange := range fileChanges {
		previousCertBytes, err := ioutil.ReadFile(certChange.dstFile)
		Expect(err).NotTo(HaveOccurred())

		newCertBytes, err := ioutil.ReadFile(certChange.srcFile)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(certChange.dstFile, newCertBytes, 06444)
		Expect(err).NotTo(HaveOccurred())

		if !strings.Contains(certChange.dstFile, "server.crt") {
			continue
		}
		serverCertChanges = append(serverCertChanges, certChange)
		certChange.newCert = newCertBytes
		certChange.oldCert = previousCertBytes
	}
	return serverCertChanges
}


func assertBlockReception(expectedHeightsPerChannel map[string]int, orderers []*nwo.Orderer, p *nwo.Peer, n *nwo.Network) {
	assertReception := func(channelName string, blockSeq int) {
		var wg sync.WaitGroup
		wg.Add(len(orderers))
		for _, orderer := range orderers {
			go func(orderer *nwo.Orderer) {
				defer wg.Done()
				waitForBlockReception(orderer, p, n, channelName, blockSeq)
			}(orderer)
		}
		wg.Wait()
	}

	var wg sync.WaitGroup
	wg.Add(len(expectedHeightsPerChannel))

	for channelName, blockSeq := range expectedHeightsPerChannel {
		go func(channelName string, blockSeq int) {
			defer wg.Done()
			assertReception(channelName, blockSeq)
		}(channelName, blockSeq)
	}
	wg.Wait()
}

func waitForBlockReception(o *nwo.Orderer, submitter *nwo.Peer, network *nwo.Network, channelName string, blockSeq int) {
	c := commands.ChannelFetch{
		ChannelID:  channelName,
		Block:      "newest",
		OutputFile: "/dev/null",
		Orderer:    network.OrdererAddress(o, nwo.ListenPort),
	}
	Eventually(func() string {
		sess, err := network.OrdererAdminSession(o, submitter, c)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(0))
		sessErr := string(sess.Err.Contents())
		expected := fmt.Sprintf("Received block: %d", blockSeq)
		if strings.Contains(sessErr, expected) {
			return ""
		}
		return sessErr
	}, network.EventuallyTimeout).Should(BeEmpty())
}

func assertNoErrorsAreLogged(ordererRunners []*ginkgomon.Runner) {
	var wg sync.WaitGroup
	wg.Add(len(ordererRunners))

	assertNoErrors := func(runner *ginkgomon.Runner) {
		buff := runner.Err()
		readOutput := func() string {
			out := bytes.Buffer{}
			
			for {
				b := make([]byte, 1024)
				n, _ := buff.Read(b)
				if n == 0 {
					break
				}
				bytesRead := make([]byte, n)
				copy(bytesRead, b)
				out.Write(bytesRead)
			}
			return out.String()
		}
		Eventually(readOutput, time.Minute, time.Second*5).Should(Not(ContainSubstring("ERRO")))
	}

	for _, runner := range ordererRunners {
		go func(runner *ginkgomon.Runner) {
			defer wg.Done()
			assertNoErrors(runner)
		}(runner)
	}
	wg.Wait()
}

func deployChaincodes(n *nwo.Network, p *nwo.Peer, o *nwo.Orderer, mycc nwo.Chaincode, mycc2 nwo.Chaincode, mycc3 nwo.Chaincode) {
	var wg sync.WaitGroup
	wg.Add(3)
	for channel, chaincode := range map[string]nwo.Chaincode{
		"testchannel":  mycc,
		"testchannel2": mycc2,
		"testchannel3": mycc3,
	} {
		go func(channel string, cc nwo.Chaincode) {
			defer wg.Done()
			nwo.DeployChaincode(n, channel, o, cc, p)
		}(channel, chaincode)
	}

	wg.Wait()
}

func assertInvoke(network *nwo.Network, peer *nwo.Peer, o *nwo.Orderer, cc string, channel string, expectedOutput string, expectedStatus int) {
	sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeInvoke{
		ChannelID: channel,
		Orderer:   network.OrdererAddress(o, nwo.ListenPort),
		Name:      cc,
		Ctor:      `{"Args":["invoke","a","b","10"]}`,
		PeerAddresses: []string{
			network.PeerAddress(peer, nwo.ListenPort),
		},
		WaitForEvent: true,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, network.EventuallyTimeout).Should(gexec.Exit(expectedStatus))
	Expect(sess.Err).To(gbytes.Say(expectedOutput))
}

func revokeReaderAccess(network *nwo.Network, channel string, orderer *nwo.Orderer, peer *nwo.Peer) {
	config := nwo.GetConfig(network, peer, orderer, channel)
	updatedConfig := proto.Clone(config).(*common.Config)

	
	adminPolicy := utils.MarshalOrPanic(&common.ImplicitMetaPolicy{
		SubPolicy: "Admins",
		Rule:      common.ImplicitMetaPolicy_MAJORITY,
	})
	updatedConfig.ChannelGroup.Groups["Orderer"].Policies["Readers"].Policy.Value = adminPolicy
	nwo.UpdateOrdererConfig(network, orderer, channel, config, updatedConfig, peer, orderer)
}
