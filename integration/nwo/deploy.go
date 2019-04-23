/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

type Chaincode struct {
	Name                string
	Version             string
	Path                string
	Ctor                string
	Policy              string
	Lang                string
	CollectionsConfig   string 
	PackageFile         string
	PackageID           string 
	Sequence            string
	EndorsementPlugin   string
	ValidationPlugin    string
	InitRequired        bool
	Label               string
	SignaturePolicy     string
	ChannelConfigPolicy string
}







func DeployChaincodeNewLifecycle(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
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

	
	PackageChaincodeNewLifecycle(n, chaincode, peers[0])

	
	filebytes, err := ioutil.ReadFile(chaincode.PackageFile)
	Expect(err).NotTo(HaveOccurred())
	hashStr := fmt.Sprintf("%x", util.ComputeSHA256(filebytes))
	chaincode.PackageID = chaincode.Label + ":" + hashStr

	
	InstallChaincodeNewLifecycle(n, chaincode, peers...)

	
	
	maxLedgerHeight := GetMaxLedgerHeight(n, channel, peers...)

	
	ApproveChaincodeForMyOrgNewLifecycle(n, channel, orderer, chaincode, peers...)

	
	
	
	WaitUntilEqualLedgerHeight(n, channel, maxLedgerHeight+len(n.PeerOrgs()), peers...)

	
	CommitChaincodeNewLifecycle(n, channel, orderer, chaincode, peers[0], peers...)

	
	if chaincode.InitRequired {
		InitChaincodeNewLifecycle(n, channel, orderer, chaincode, peers...)
	}
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

func PackageChaincodeNewLifecycle(n *Network, chaincode Chaincode, peer *Peer) {
	sess, err := n.PeerAdminSession(peer, commands.ChaincodePackageLifecycle{
		Path:       chaincode.Path,
		Lang:       chaincode.Lang,
		Label:      chaincode.Label,
		OutputFile: chaincode.PackageFile,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
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

func InstallChaincodeNewLifecycle(n *Network, chaincode Chaincode, peers ...*Peer) {
	for _, p := range peers {
		sess, err := n.PeerAdminSession(p, commands.ChaincodeInstallLifecycle{
			PackageFile: chaincode.PackageFile,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

		sess, err = n.PeerAdminSession(p, commands.ChaincodeQueryInstalledLifecycle{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
		Expect(sess).To(gbytes.Say(fmt.Sprintf("Package ID: %s, Label: %s", chaincode.PackageID, chaincode.Label)))
	}
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

func ApproveChaincodeForMyOrgNewLifecycle(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
	if chaincode.PackageID == "" {
		pkgBytes, err := ioutil.ReadFile(chaincode.PackageFile)
		Expect(err).NotTo(HaveOccurred())
		hash := util.ComputeSHA256(pkgBytes)
		chaincode.PackageID = fmt.Sprintf("%s:%x", chaincode.Label, hash)
	}

	
	approvedOrgs := map[string]bool{}
	for _, p := range peers {
		if _, ok := approvedOrgs[p.Organization]; !ok {
			sess, err := n.PeerAdminSession(p, commands.ChaincodeApproveForMyOrgLifecycle{
				ChannelID:           channel,
				Orderer:             n.OrdererAddress(orderer, ListenPort),
				Name:                chaincode.Name,
				Version:             chaincode.Version,
				PackageID:           chaincode.PackageID,
				Sequence:            chaincode.Sequence,
				EndorsementPlugin:   chaincode.EndorsementPlugin,
				ValidationPlugin:    chaincode.ValidationPlugin,
				SignaturePolicy:     chaincode.SignaturePolicy,
				ChannelConfigPolicy: chaincode.ChannelConfigPolicy,
				InitRequired:        chaincode.InitRequired,
				CollectionsConfig:   chaincode.CollectionsConfig,
			})
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
			approvedOrgs[p.Organization] = true
			Eventually(sess.Err, n.EventuallyTimeout).Should(gbytes.Say(`\Qcommitted with status (VALID)\E`))
		}
	}
}

func CommitChaincodeNewLifecycle(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peer *Peer, checkPeers ...*Peer) {
	
	commitOrgs := map[string]bool{}
	var peerAddresses []string
	for _, p := range checkPeers {
		if exists := commitOrgs[p.Organization]; !exists {
			peerAddresses = append(peerAddresses, n.PeerAddress(p, ListenPort))
			commitOrgs[p.Organization] = true
		}
	}

	sess, err := n.PeerAdminSession(peer, commands.ChaincodeCommitLifecycle{
		ChannelID:           channel,
		Orderer:             n.OrdererAddress(orderer, ListenPort),
		Name:                chaincode.Name,
		Version:             chaincode.Version,
		Sequence:            chaincode.Sequence,
		EndorsementPlugin:   chaincode.EndorsementPlugin,
		ValidationPlugin:    chaincode.ValidationPlugin,
		SignaturePolicy:     chaincode.SignaturePolicy,
		ChannelConfigPolicy: chaincode.ChannelConfigPolicy,
		InitRequired:        chaincode.InitRequired,
		CollectionsConfig:   chaincode.CollectionsConfig,
		PeerAddresses:       peerAddresses,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	for i := 0; i < len(peerAddresses); i++ {
		Eventually(sess.Err, n.EventuallyTimeout).Should(gbytes.Say(`\Qcommitted with status (VALID)\E`))
	}
	EnsureCommitted(n, channel, chaincode.Name, chaincode.Version, chaincode.Sequence, checkPeers...)
}

func EnsureCommitted(n *Network, channel, name, version, sequence string, peers ...*Peer) {
	for _, p := range peers {
		Eventually(listCommitted(n, p, channel, name), n.EventuallyTimeout).Should(
			gbytes.Say(fmt.Sprintf("Committed chaincode definition for chaincode '%s' on channel '%s':\nVersion: %s, Sequence: %s", name, channel, version, sequence)),
		)
	}
}

func InitChaincodeNewLifecycle(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
	
	initOrgs := map[string]bool{}
	var peerAddresses []string
	for _, p := range peers {
		if exists := initOrgs[p.Organization]; !exists {
			peerAddresses = append(peerAddresses, n.PeerAddress(p, ListenPort))
			initOrgs[p.Organization] = true
		}
	}

	sess, err := n.PeerUserSession(peers[0], "User1", commands.ChaincodeInvoke{
		ChannelID:     channel,
		Orderer:       n.OrdererAddress(orderer, ListenPort),
		Name:          chaincode.Name,
		Ctor:          chaincode.Ctor,
		PeerAddresses: peerAddresses,
		WaitForEvent:  true,
		IsInit:        true,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	for i := 0; i < len(peerAddresses); i++ {
		Eventually(sess.Err, n.EventuallyTimeout).Should(gbytes.Say(`\Qcommitted with status (VALID)\E`))
	}
	Expect(sess.Err).To(gbytes.Say("Chaincode invoke successful. result: status:200"))
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

func listCommitted(n *Network, peer *Peer, channel, name string) func() *gbytes.Buffer {
	return func() *gbytes.Buffer {
		sess, err := n.PeerAdminSession(peer, commands.ChaincodeListCommittedLifecycle{
			ChannelID: channel,
			Name:      name,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
		return sess.Buffer()
	}
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





func EnableV2_0Capabilities(network *Network, channel string, orderer *Orderer, peers ...*Peer) {
	if len(peers) == 0 {
		return
	}

	config := GetConfig(network, peers[0], orderer, channel)
	updatedConfig := proto.Clone(config).(*common.Config)

	
	updatedConfig.ChannelGroup.Groups["Application"].Values["Capabilities"] = &common.ConfigValue{
		ModPolicy: "Admins",
		Value: protoutil.MarshalOrPanic(
			&common.Capabilities{
				Capabilities: map[string]*common.Capability{
					"V2_0": {},
				},
			},
		),
	}

	UpdateConfig(network, orderer, channel, config, updatedConfig, false, peers[0], peers...)
}



func WaitUntilEqualLedgerHeight(n *Network, channel string, height int, peers ...*Peer) {
	for _, peer := range peers {
		Eventually(func() int {
			return GetLedgerHeight(n, peer, channel)
		}, n.EventuallyTimeout).Should(Equal(height))
	}
}



func GetLedgerHeight(n *Network, peer *Peer, channel string) int {
	sess, err := n.PeerUserSession(peer, "User1", commands.ChannelInfo{
		ChannelID: channel,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

	channelInfoStr := strings.TrimPrefix(string(sess.Buffer().Contents()[:]), "Blockchain info:")
	var channelInfo = common.BlockchainInfo{}
	json.Unmarshal([]byte(channelInfoStr), &channelInfo)
	return int(channelInfo.Height)
}



func GetMaxLedgerHeight(n *Network, channel string, peers ...*Peer) int {
	var maxHeight int
	for _, peer := range peers {
		peerHeight := GetLedgerHeight(n, peer, channel)
		if peerHeight > maxHeight {
			maxHeight = peerHeight
		}
	}
	return maxHeight
}
