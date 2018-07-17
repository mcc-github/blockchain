/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo

import (
	"fmt"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/mcc-github/blockchain/integration/nwo/commands"
)

type Chaincode struct {
	Name    string
	Version string
	Path    string
	Ctor    string
	Policy  string
}




func DeployChaincode(n *Network, channel string, orderer *Orderer, chaincode Chaincode) {
	peers := n.PeersWithChannel(channel)
	if len(peers) == 0 {
		return
	}

	
	n.InstallChaincode(peers, commands.ChaincodeInstall{
		Name:    chaincode.Name,
		Version: chaincode.Version,
		Path:    chaincode.Path,
	})

	
	n.InstantiateChaincode(peers[0], commands.ChaincodeInstantiate{
		ChannelID: channel,
		Orderer:   n.OrdererAddress(orderer, ListenPort),
		Name:      chaincode.Name,
		Version:   chaincode.Version,
		Ctor:      chaincode.Ctor,
		Policy:    chaincode.Policy,
	})

	
	for _, p := range peers[1:] {
		Eventually(listInstantiated(n, p, channel), time.Minute).Should(
			gbytes.Say(fmt.Sprintf("Name: %s, Version: %s,", chaincode.Name, chaincode.Version)),
		)
	}
}

func listInstantiated(n *Network, peer *Peer, channel string) func() *gbytes.Buffer {
	return func() *gbytes.Buffer {
		sess, err := n.PeerAdminSession(peer, commands.ChaincodeListInstantiated{
			ChannelID: channel,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, 10*time.Second).Should(gexec.Exit(0))
		return sess.Buffer()
	}
}
