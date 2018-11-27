/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo

import (
	"encoding/json"

	"github.com/mcc-github/blockchain/integration/nwo/commands"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)



type DiscoveredPeer struct {
	MSPID      string   `yaml:"mspid,omitempty"`
	Endpoint   string   `yaml:"endpoint,omitempty"`
	Identity   string   `yaml:"identity,omitempty"`
	Chaincodes []string `yaml:"chaincodes,omitempty"`
}



func DiscoverPeers(n *Network, p *Peer, user, channelName string) func() []DiscoveredPeer {
	return func() []DiscoveredPeer {
		peers := commands.Peers{
			UserCert: n.PeerUserCert(p, user),
			UserKey:  n.PeerUserKey(p, user),
			MSPID:    n.Organization(p.Organization).MSPID,
			Server:   n.PeerAddress(p, ListenPort),
			Channel:  channelName,
		}
		sess, err := n.Discover(peers)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess).Should(gexec.Exit(0))

		var discovered []DiscoveredPeer
		err = json.Unmarshal(sess.Out.Contents(), &discovered)
		Expect(err).NotTo(HaveOccurred())
		return discovered
	}
}
