/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/mcc-github/blockchain/integration/nwo/commands"
)

type Chaincode struct {
	Name              string
	Version           string
	Path              string
	Ctor              string
	Policy            string
	Lang              string
	CollectionsConfig string 
	PackageFile       string
}




func DeployChaincode(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
	if len(peers) == 0 {
		peers = n.PeersWithChannel(channel)
	}
	if len(peers) == 0 {
		return
	}

	
	if chaincode.PackageFile == "" {
		tempFile, err := ioutil.TempFile("", "chaincode-package")
		Expect(err).NotTo(HaveOccurred())
		tempFile.Close()
		defer os.Remove(tempFile.Name())
		chaincode.PackageFile = tempFile.Name()
	}

	
	PackageChaincode(n, chaincode, peers[0])

	
	InstallChaincode(n, chaincode, peers...)

	
	InstantiateChaincode(n, channel, orderer, chaincode, peers[0], peers...)
}

func PackageChaincode(n *Network, chaincode Chaincode, peer *Peer) {
	sess, err := n.PeerAdminSession(peer, commands.ChaincodePackage{
		Name:       chaincode.Name,
		Version:    chaincode.Version,
		Path:       chaincode.Path,
		Lang:       chaincode.Lang,
		OutputFile: chaincode.PackageFile,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
}

func InstallChaincode(n *Network, chaincode Chaincode, peers ...*Peer) {
	for _, p := range peers {
		sess, err := n.PeerAdminSession(p, commands.ChaincodeInstall{
			Name:        chaincode.Name,
			Version:     chaincode.Version,
			Path:        chaincode.Path,
			Lang:        chaincode.Lang,
			PackageFile: chaincode.PackageFile,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

		sess, err = n.PeerAdminSession(p, commands.ChaincodeListInstalled{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
		Expect(sess).To(gbytes.Say(fmt.Sprintf("Name: %s, Version: %s,", chaincode.Name, chaincode.Version)))
	}
}

func InstantiateChaincode(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peer *Peer, checkPeers ...*Peer) {
	sess, err := n.PeerAdminSession(peer, commands.ChaincodeInstantiate{
		ChannelID:         channel,
		Orderer:           n.OrdererAddress(orderer, ListenPort),
		Name:              chaincode.Name,
		Version:           chaincode.Version,
		Ctor:              chaincode.Ctor,
		Policy:            chaincode.Policy,
		Lang:              chaincode.Lang,
		CollectionsConfig: chaincode.CollectionsConfig,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

	EnsureInstantiated(n, channel, chaincode.Name, chaincode.Version, checkPeers...)
}

func EnsureInstantiated(n *Network, channel, name, version string, peers ...*Peer) {
	for _, p := range peers {
		Eventually(listInstantiated(n, p, channel), n.EventuallyTimeout).Should(
			gbytes.Say(fmt.Sprintf("Name: %s, Version: %s,", name, version)),
		)
	}
}

func UpgradeChaincode(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
	if len(peers) == 0 {
		peers = n.PeersWithChannel(channel)
	}
	if len(peers) == 0 {
		return
	}

	
	InstallChaincode(n, chaincode, peers...)

	
	sess, err := n.PeerAdminSession(peers[0], commands.ChaincodeUpgrade{
		ChannelID:         channel,
		Orderer:           n.OrdererAddress(orderer, ListenPort),
		Name:              chaincode.Name,
		Version:           chaincode.Version,
		Ctor:              chaincode.Ctor,
		Policy:            chaincode.Policy,
		CollectionsConfig: chaincode.CollectionsConfig,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

	EnsureInstantiated(n, channel, chaincode.Name, chaincode.Version, peers...)
}

func listInstantiated(n *Network, peer *Peer, channel string) func() *gbytes.Buffer {
	return func() *gbytes.Buffer {
		sess, err := n.PeerAdminSession(peer, commands.ChaincodeListInstantiated{
			ChannelID: channel,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
		return sess.Buffer()
	}
}
