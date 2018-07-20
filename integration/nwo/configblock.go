/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/tools/configtxlator/update"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
)



func GetConfigBlock(n *Network, peer *Peer, orderer *Orderer, channel string) *common.Config {
	tempDir, err := ioutil.TempDir("", "getConfigBlock")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(tempDir)

	
	output := filepath.Join(tempDir, "config_block.pb")
	sess, err := n.PeerAdminSession(peer, commands.ChannelFetch{
		ChannelID:  channel,
		Block:      "config",
		Orderer:    n.OrdererAddress(orderer, ListenPort),
		OutputFile: output,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess.Err).To(gbytes.Say("Received block: "))

	
	configBlock := UnmarshalBlockFromFile(output)

	
	envelope, err := utils.GetEnvelopeFromBlock(configBlock.Data.Data[0])
	Expect(err).NotTo(HaveOccurred())

	
	payload, err := utils.GetPayload(envelope)
	Expect(err).NotTo(HaveOccurred())

	
	configEnv := &common.ConfigEnvelope{}
	err = proto.Unmarshal(payload.Data, configEnv)
	Expect(err).NotTo(HaveOccurred())

	
	return configEnv.Config
}



func UpdateConfig(n *Network, orderer *Orderer, channel string, current, updated *common.Config, submitter *Peer, additionalSigners ...*Peer) {
	tempDir, err := ioutil.TempDir("", "updateConfig")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(tempDir)

	
	configUpdate, err := update.Compute(current, updated)
	Expect(err).NotTo(HaveOccurred())
	configUpdate.ChannelId = channel

	signedEnvelope, err := utils.CreateSignedEnvelope(
		common.HeaderType_CONFIG_UPDATE,
		channel,
		nil, 
		&common.ConfigUpdateEnvelope{ConfigUpdate: utils.MarshalOrPanic(configUpdate)},
		0, 
		0, 
	)
	Expect(err).NotTo(HaveOccurred())
	Expect(signedEnvelope).NotTo(BeNil())

	updateFile := filepath.Join(tempDir, "update.pb")
	err = ioutil.WriteFile(updateFile, utils.MarshalOrPanic(signedEnvelope), 0600)
	Expect(err).NotTo(HaveOccurred())

	for _, signer := range additionalSigners {
		sess, err := n.PeerAdminSession(signer, commands.SignConfigTx{File: updateFile})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	}

	
	currentBlockNumber := CurrentConfigBlockNumber(n, submitter, orderer, channel)

	sess, err := n.PeerAdminSession(submitter, commands.ChannelUpdate{
		ChannelID: channel,
		Orderer:   n.OrdererAddress(orderer, ListenPort),
		File:      updateFile,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess.Err).To(gbytes.Say("Successfully submitted channel update"))

	
	ccb := func() uint64 { return CurrentConfigBlockNumber(n, submitter, orderer, channel) }
	Eventually(ccb, n.EventuallyTimeout).Should(BeNumerically(">", currentBlockNumber))
}




func CurrentConfigBlockNumber(n *Network, peer *Peer, orderer *Orderer, channel string) uint64 {
	tempDir, err := ioutil.TempDir("", "currentConfigBlock")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(tempDir)

	
	output := filepath.Join(tempDir, "config_block.pb")
	sess, err := n.PeerAdminSession(peer, commands.ChannelFetch{
		ChannelID:  channel,
		Block:      "config",
		Orderer:    n.OrdererAddress(orderer, ListenPort),
		OutputFile: output,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess.Err).To(gbytes.Say("Received block: "))

	
	configBlock := UnmarshalBlockFromFile(output)
	return configBlock.Header.Number
}


func UnmarshalBlockFromFile(blockFile string) *common.Block {
	blockBytes, err := ioutil.ReadFile(blockFile)
	Expect(err).NotTo(HaveOccurred())

	block, err := utils.UnmarshalBlock(blockBytes)
	Expect(err).NotTo(HaveOccurred())

	return block
}
