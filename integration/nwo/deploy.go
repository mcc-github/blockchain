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
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer/lifecycle"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
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

func (c *Chaincode) SetPackageIDFromPackageFile() {
	fileBytes, err := ioutil.ReadFile(c.PackageFile)
	Expect(err).NotTo(HaveOccurred())
	hashStr := fmt.Sprintf("%x", util.ComputeSHA256(fileBytes))
	c.PackageID = c.Label + ":" + hashStr
}







func DeployChaincode(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
	if len(peers) == 0 {
		peers = n.PeersWithChannel(channel)
	}
	if len(peers) == 0 {
		return
	}

	PackageAndInstallChaincode(n, chaincode, peers...)

	
	ApproveChaincodeForMyOrg(n, channel, orderer, chaincode, peers...)

	
	CheckCommitReadinessUntilReady(n, channel, chaincode, n.PeerOrgs(), peers...)
	CommitChaincode(n, channel, orderer, chaincode, peers[0], peers...)

	
	if chaincode.InitRequired {
		InitChaincode(n, channel, orderer, chaincode, peers...)
	}
}









func DeployChaincodeLegacy(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
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

	
	PackageChaincodeLegacy(n, chaincode, peers[0])

	
	InstallChaincodeLegacy(n, chaincode, peers...)

	
	InstantiateChaincodeLegacy(n, channel, orderer, chaincode, peers[0], peers...)
}

func PackageAndInstallChaincode(n *Network, chaincode Chaincode, peers ...*Peer) {
	
	if chaincode.PackageFile == "" {
		tempFile, err := ioutil.TempFile("", "chaincode-package")
		Expect(err).NotTo(HaveOccurred())
		tempFile.Close()
		defer os.Remove(tempFile.Name())
		chaincode.PackageFile = tempFile.Name()
	}

	
	if _, err := os.Stat(chaincode.PackageFile); os.IsNotExist(err) {
		
		PackageChaincode(n, chaincode, peers[0])
	}

	
	chaincode.SetPackageIDFromPackageFile()

	
	InstallChaincode(n, chaincode, peers...)
}

func PackageChaincode(n *Network, chaincode Chaincode, peer *Peer) {
	sess, err := n.PeerAdminSession(peer, commands.ChaincodePackage{
		Path:       chaincode.Path,
		Lang:       chaincode.Lang,
		Label:      chaincode.Label,
		OutputFile: chaincode.PackageFile,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
}

func PackageChaincodeLegacy(n *Network, chaincode Chaincode, peer *Peer) {
	sess, err := n.PeerAdminSession(peer, commands.ChaincodePackageLegacy{
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
			PackageFile: chaincode.PackageFile,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

		EnsureInstalled(n, chaincode.Label, chaincode.PackageID, p)
	}
}

func InstallChaincodeLegacy(n *Network, chaincode Chaincode, peers ...*Peer) {
	for _, p := range peers {
		sess, err := n.PeerAdminSession(p, commands.ChaincodeInstallLegacy{
			Name:        chaincode.Name,
			Version:     chaincode.Version,
			Path:        chaincode.Path,
			Lang:        chaincode.Lang,
			PackageFile: chaincode.PackageFile,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

		sess, err = n.PeerAdminSession(p, commands.ChaincodeListInstalledLegacy{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
		Expect(sess).To(gbytes.Say(fmt.Sprintf("Name: %s, Version: %s,", chaincode.Name, chaincode.Version)))
	}
}

func ApproveChaincodeForMyOrg(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
	if chaincode.PackageID == "" {
		pkgBytes, err := ioutil.ReadFile(chaincode.PackageFile)
		Expect(err).NotTo(HaveOccurred())
		hash := util.ComputeSHA256(pkgBytes)
		chaincode.PackageID = fmt.Sprintf("%s:%x", chaincode.Label, hash)
	}

	
	approvedOrgs := map[string]bool{}
	for _, p := range peers {
		if _, ok := approvedOrgs[p.Organization]; !ok {
			sess, err := n.PeerAdminSession(p, commands.ChaincodeApproveForMyOrg{
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

func CheckCommitReadinessUntilReady(n *Network, channel string, chaincode Chaincode, checkOrgs []*Organization, peers ...*Peer) {
	for _, p := range peers {
		keys := Keys{}
		for _, org := range checkOrgs {
			keys[org.MSPID] = BeTrue()
		}
		Eventually(checkCommitReadiness(n, p, channel, chaincode), n.EventuallyTimeout).Should(
			MatchKeys(IgnoreExtras, keys),
		)
	}
}

func CommitChaincode(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peer *Peer, checkPeers ...*Peer) {
	
	commitOrgs := map[string]bool{}
	var peerAddresses []string
	for _, p := range checkPeers {
		if exists := commitOrgs[p.Organization]; !exists {
			peerAddresses = append(peerAddresses, n.PeerAddress(p, ListenPort))
			commitOrgs[p.Organization] = true
		}
	}

	sess, err := n.PeerAdminSession(peer, commands.ChaincodeCommit{
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
	checkOrgs := []*Organization{}
	for org := range commitOrgs {
		checkOrgs = append(checkOrgs, n.Organization(org))
	}
	EnsureChaincodeCommitted(n, channel, chaincode.Name, chaincode.Version, chaincode.Sequence, checkOrgs, checkPeers...)
}



func EnsureChaincodeCommitted(n *Network, channel, name, version, sequence string, checkOrgs []*Organization, peers ...*Peer) {
	for _, p := range peers {
		sequenceInt, err := strconv.ParseInt(sequence, 10, 64)
		Expect(err).NotTo(HaveOccurred())
		approvedKeys := Keys{}
		for _, org := range checkOrgs {
			approvedKeys[org.MSPID] = BeTrue()
		}
		Eventually(listCommitted(n, p, channel, name), n.EventuallyTimeout).Should(
			MatchFields(IgnoreExtras, Fields{
				"Version":   Equal(version),
				"Sequence":  Equal(sequenceInt),
				"Approvals": MatchKeys(IgnoreExtras, approvedKeys),
			}),
		)
	}
}

func InitChaincode(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
	
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

func InstantiateChaincodeLegacy(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peer *Peer, checkPeers ...*Peer) {
	sess, err := n.PeerAdminSession(peer, commands.ChaincodeInstantiateLegacy{
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

	EnsureInstantiatedLegacy(n, channel, chaincode.Name, chaincode.Version, checkPeers...)
}

func EnsureInstantiatedLegacy(n *Network, channel, name, version string, peers ...*Peer) {
	for _, p := range peers {
		Eventually(listInstantiatedLegacy(n, p, channel), n.EventuallyTimeout).Should(
			gbytes.Say(fmt.Sprintf("Name: %s, Version: %s,", name, version)),
		)
	}
}

func UpgradeChaincodeLegacy(n *Network, channel string, orderer *Orderer, chaincode Chaincode, peers ...*Peer) {
	if len(peers) == 0 {
		peers = n.PeersWithChannel(channel)
	}
	if len(peers) == 0 {
		return
	}

	
	InstallChaincodeLegacy(n, chaincode, peers...)

	
	sess, err := n.PeerAdminSession(peers[0], commands.ChaincodeUpgradeLegacy{
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

	EnsureInstantiatedLegacy(n, channel, chaincode.Name, chaincode.Version, peers...)
}

func EnsureInstalled(n *Network, label, packageID string, peers ...*Peer) {
	for _, p := range peers {
		Eventually(queryInstalled(n, p), n.EventuallyTimeout).Should(
			ContainElement(MatchFields(IgnoreExtras,
				Fields{
					"Label":     Equal(label),
					"PackageId": Equal(packageID),
				},
			)),
		)
	}
}

func QueryInstalledReferences(n *Network, channel, label, packageID string, checkPeer *Peer, nameVersions ...[]string) {
	chaincodes := make([]*lifecycle.QueryInstalledChaincodesResult_Chaincode, len(nameVersions))
	for i, nameVersion := range nameVersions {
		chaincodes[i] = &lifecycle.QueryInstalledChaincodesResult_Chaincode{
			Name:    nameVersion[0],
			Version: nameVersion[1],
		}
	}

	Expect(queryInstalled(n, checkPeer)()).To(
		ContainElement(MatchFields(IgnoreExtras,
			Fields{
				"Label":     Equal(label),
				"PackageId": Equal(packageID),
				"References": HaveKeyWithValue(channel, PointTo(MatchFields(IgnoreExtras,
					Fields{
						"Chaincodes": ConsistOf(chaincodes),
					},
				))),
			},
		)),
	)
}

type queryInstalledOutput struct {
	InstalledChaincodes []lifecycle.QueryInstalledChaincodesResult_InstalledChaincode `json:"installed_chaincodes"`
}

func queryInstalled(n *Network, peer *Peer) func() []lifecycle.QueryInstalledChaincodesResult_InstalledChaincode {
	return func() []lifecycle.QueryInstalledChaincodesResult_InstalledChaincode {
		sess, err := n.PeerAdminSession(peer, commands.ChaincodeQueryInstalled{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
		output := &queryInstalledOutput{}
		err = json.Unmarshal(sess.Out.Contents(), output)
		Expect(err).NotTo(HaveOccurred())
		return output.InstalledChaincodes
	}
}

type checkCommitReadinessOutput struct {
	Approvals map[string]bool `json:"approvals"`
}

func checkCommitReadiness(n *Network, peer *Peer, channel string, chaincode Chaincode) func() map[string]bool {
	return func() map[string]bool {
		sess, err := n.PeerAdminSession(peer, commands.ChaincodeCheckCommitReadiness{
			ChannelID:           channel,
			Name:                chaincode.Name,
			Version:             chaincode.Version,
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
		output := &checkCommitReadinessOutput{}
		err = json.Unmarshal(sess.Out.Contents(), output)
		Expect(err).NotTo(HaveOccurred())
		return output.Approvals
	}
}

type queryCommittedOutput struct {
	Sequence  int64           `json:"sequence"`
	Version   string          `json:"version"`
	Approvals map[string]bool `json:"approvals"`
}




func listCommitted(n *Network, peer *Peer, channel, name string) func() queryCommittedOutput {
	return func() queryCommittedOutput {
		sess, err := n.PeerAdminSession(peer, commands.ChaincodeListCommitted{
			ChannelID: channel,
			Name:      name,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit())
		output := &queryCommittedOutput{}
		if sess.ExitCode() == 1 {
			
			return *output
		}
		err = json.Unmarshal(sess.Out.Contents(), output)
		Expect(err).NotTo(HaveOccurred())
		return *output
	}
}

func listInstantiatedLegacy(n *Network, peer *Peer, channel string) func() *gbytes.Buffer {
	return func() *gbytes.Buffer {
		sess, err := n.PeerAdminSession(peer, commands.ChaincodeListInstantiatedLegacy{
			ChannelID: channel,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
		return sess.Buffer()
	}
}





func EnableCapabilities(network *Network, channel, capabilitiesGroup, capabilitiesVersion string, orderer *Orderer, peers ...*Peer) {
	if len(peers) == 0 {
		return
	}

	config := GetConfig(network, peers[0], orderer, channel)
	updatedConfig := proto.Clone(config).(*common.Config)

	updatedConfig.ChannelGroup.Groups[capabilitiesGroup].Values["Capabilities"] = &common.ConfigValue{
		ModPolicy: "Admins",
		Value: protoutil.MarshalOrPanic(
			&common.Capabilities{
				Capabilities: map[string]*common.Capability{
					capabilitiesVersion: {},
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
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit())

	if sess.ExitCode() == 1 {
		
		return -1
	}

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
