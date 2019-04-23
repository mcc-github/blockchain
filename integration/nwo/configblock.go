/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/internal/configtxlator/update"
	"github.com/mcc-github/blockchain/protos/common"
	protosorderer "github.com/mcc-github/blockchain/protos/orderer"
	ectdraft_protos "github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	"github.com/mcc-github/blockchain/protoutil"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)


func GetConfigBlock(n *Network, peer *Peer, orderer *Orderer, channel string) *common.Block {
	tempDir, err := ioutil.TempDir("", "getConfigBlock")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(tempDir)

	
	output := filepath.Join(tempDir, "config_block.pb")
	sess, err := n.OrdererAdminSession(orderer, peer, commands.ChannelFetch{
		ChannelID:  channel,
		Block:      "config",
		Orderer:    n.OrdererAddress(orderer, ListenPort),
		OutputFile: output,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess.Err).To(gbytes.Say("Received block: "))

	
	configBlock := UnmarshalBlockFromFile(output)
	return configBlock
}


func GetConfig(n *Network, peer *Peer, orderer *Orderer, channel string) *common.Config {
	configBlock := GetConfigBlock(n, peer, orderer, channel)
	
	envelope, err := protoutil.GetEnvelopeFromBlock(configBlock.Data.Data[0])
	Expect(err).NotTo(HaveOccurred())

	
	payload, err := protoutil.GetPayload(envelope)
	Expect(err).NotTo(HaveOccurred())

	
	configEnv := &common.ConfigEnvelope{}
	err = proto.Unmarshal(payload.Data, configEnv)
	Expect(err).NotTo(HaveOccurred())

	
	return configEnv.Config
}



func UpdateConfig(n *Network, orderer *Orderer, channel string, current, updated *common.Config, getConfigBlockFromOrderer bool, submitter *Peer, additionalSigners ...*Peer) {
	tempDir, err := ioutil.TempDir("", "updateConfig")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(tempDir)

	
	configUpdate, err := update.Compute(current, updated)
	Expect(err).NotTo(HaveOccurred())
	configUpdate.ChannelId = channel

	signedEnvelope, err := protoutil.CreateSignedEnvelope(
		common.HeaderType_CONFIG_UPDATE,
		channel,
		nil, 
		&common.ConfigUpdateEnvelope{ConfigUpdate: protoutil.MarshalOrPanic(configUpdate)},
		0, 
		0, 
	)
	Expect(err).NotTo(HaveOccurred())
	Expect(signedEnvelope).NotTo(BeNil())

	updateFile := filepath.Join(tempDir, "update.pb")
	err = ioutil.WriteFile(updateFile, protoutil.MarshalOrPanic(signedEnvelope), 0600)
	Expect(err).NotTo(HaveOccurred())

	for _, signer := range additionalSigners {
		sess, err := n.PeerAdminSession(signer, commands.SignConfigTx{File: updateFile})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	}

	var currentBlockNumber uint64
	
	if getConfigBlockFromOrderer {
		currentBlockNumber = CurrentConfigBlockNumber(n, submitter, orderer, channel)
	} else {
		currentBlockNumber = CurrentConfigBlockNumber(n, submitter, nil, channel)
	}

	sess, err := n.PeerAdminSession(submitter, commands.ChannelUpdate{
		ChannelID: channel,
		Orderer:   n.OrdererAddress(orderer, ListenPort),
		File:      updateFile,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess.Err).To(gbytes.Say("Successfully submitted channel update"))

	if getConfigBlockFromOrderer {
		ccb := func() uint64 { return CurrentConfigBlockNumber(n, submitter, orderer, channel) }
		Eventually(ccb, n.EventuallyTimeout).Should(BeNumerically(">", currentBlockNumber))
		return
	}
	
	
	for _, peer := range n.PeersWithChannel(channel) {
		ccb := func() uint64 { return CurrentConfigBlockNumber(n, peer, nil, channel) }
		Eventually(ccb, n.EventuallyTimeout).Should(BeNumerically(">", currentBlockNumber))
	}
}



func UpdateOrdererConfig(n *Network, orderer *Orderer, channel string, current, updated *common.Config, submitter *Peer, additionalSigners ...*Orderer) {
	tempDir, err := ioutil.TempDir("", "updateConfig")
	Expect(err).NotTo(HaveOccurred())
	updateFile := filepath.Join(tempDir, "update.pb")
	defer os.RemoveAll(tempDir)

	computeUpdateOrdererConfig(updateFile, n, channel, current, updated, submitter, additionalSigners...)

	currentBlockNumber := CurrentConfigBlockNumber(n, submitter, orderer, channel)

	Eventually(func() string {
		sess, err := n.OrdererAdminSession(orderer, submitter, commands.ChannelUpdate{
			ChannelID: channel,
			Orderer:   n.OrdererAddress(orderer, ListenPort),
			File:      updateFile,
		})
		if err != nil {
			return err.Error()
		}
		sess.Wait(n.EventuallyTimeout)
		if sess.ExitCode() != 0 {
			return fmt.Sprintf("exit code is %d", sess.ExitCode())
		}
		if strings.Contains(string(sess.Err.Contents()), "Successfully submitted channel update") {
			return ""
		}
		return fmt.Sprintf("channel update output: %s", string(sess.Err.Contents()))
	}, n.EventuallyTimeout).Should(BeEmpty())

	
	ccb := func() uint64 { return CurrentConfigBlockNumber(n, submitter, orderer, channel) }
	Eventually(ccb, n.EventuallyTimeout).Should(BeNumerically(">", currentBlockNumber))
}





func CurrentConfigBlockNumber(n *Network, peer *Peer, orderer *Orderer, channel string) uint64 {
	tempDir, err := ioutil.TempDir("", "currentConfigBlock")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(tempDir)

	
	output := filepath.Join(tempDir, "config_block.pb")
	if orderer == nil {
		return CurrentConfigBlockNumberFromPeer(n, peer, channel, output)
	}

	FetchConfigBlock(n, peer, orderer, channel, output)

	
	configBlock := UnmarshalBlockFromFile(output)

	return configBlock.Header.Number
}



func CurrentConfigBlockNumberFromPeer(n *Network, peer *Peer, channel, output string) uint64 {
	sess, err := n.PeerAdminSession(peer, commands.ChannelFetch{
		ChannelID:  channel,
		Block:      "config",
		OutputFile: output,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess.Err).To(gbytes.Say("Received block: "))

	configBlock := UnmarshalBlockFromFile(output)

	return configBlock.Header.Number
}


func FetchConfigBlock(n *Network, peer *Peer, orderer *Orderer, channel string, output string) {
	fetch := func() int {
		sess, err := n.OrdererAdminSession(orderer, peer, commands.ChannelFetch{
			ChannelID:  channel,
			Block:      "config",
			Orderer:    n.OrdererAddress(orderer, ListenPort),
			OutputFile: output,
		})
		Expect(err).NotTo(HaveOccurred())
		code := sess.Wait(n.EventuallyTimeout).ExitCode()
		if code == 0 {
			Expect(sess.Err).To(gbytes.Say("Received block: "))
		}
		return code
	}
	Eventually(fetch, n.EventuallyTimeout).Should(Equal(0))
}



func UpdateOrdererConfigFail(n *Network, orderer *Orderer, channel string, current, updated *common.Config, submitter *Peer, additionalSigners ...*Orderer) {
	tempDir, err := ioutil.TempDir("", "updateConfig")
	Expect(err).NotTo(HaveOccurred())
	updateFile := filepath.Join(tempDir, "update.pb")
	defer os.RemoveAll(tempDir)

	computeUpdateOrdererConfig(updateFile, n, channel, current, updated, submitter, additionalSigners...)

	
	sess, err := n.OrdererAdminSession(orderer, submitter, commands.ChannelUpdate{
		ChannelID: channel,
		Orderer:   n.OrdererAddress(orderer, ListenPort),
		File:      updateFile,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).ShouldNot(gexec.Exit(0))
	Expect(sess.Err).NotTo(gbytes.Say("Successfully submitted channel update"))
}

func computeUpdateOrdererConfig(updateFile string, n *Network, channel string, current, updated *common.Config, submitter *Peer, additionalSigners ...*Orderer) {
	
	configUpdate, err := update.Compute(current, updated)
	Expect(err).NotTo(HaveOccurred())
	configUpdate.ChannelId = channel

	signedEnvelope, err := protoutil.CreateSignedEnvelope(
		common.HeaderType_CONFIG_UPDATE,
		channel,
		nil, 
		&common.ConfigUpdateEnvelope{ConfigUpdate: protoutil.MarshalOrPanic(configUpdate)},
		0, 
		0, 
	)
	Expect(err).NotTo(HaveOccurred())
	Expect(signedEnvelope).NotTo(BeNil())

	err = ioutil.WriteFile(updateFile, protoutil.MarshalOrPanic(signedEnvelope), 0600)
	Expect(err).NotTo(HaveOccurred())

	for _, signer := range additionalSigners {
		sess, err := n.OrdererAdminSession(signer, submitter, commands.SignConfigTx{File: updateFile})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	}
}


func UnmarshalBlockFromFile(blockFile string) *common.Block {
	blockBytes, err := ioutil.ReadFile(blockFile)
	Expect(err).NotTo(HaveOccurred())

	block, err := protoutil.UnmarshalBlock(blockBytes)
	Expect(err).NotTo(HaveOccurred())

	return block
}


func AddConsenter(n *Network, peer *Peer, orderer *Orderer, channel string, consenter ectdraft_protos.Consenter) {
	UpdateEtcdRaftMetadata(n, peer, orderer, channel, func(metadata *ectdraft_protos.ConfigMetadata) {
		metadata.Consenters = append(metadata.Consenters, &consenter)
	})
}


func RemoveConsenter(n *Network, peer *Peer, orderer *Orderer, channel string, certificate []byte) {
	UpdateEtcdRaftMetadata(n, peer, orderer, channel, func(metadata *ectdraft_protos.ConfigMetadata) {
		var newConsenters []*ectdraft_protos.Consenter
		for _, consenter := range metadata.Consenters {
			if bytes.Equal(consenter.ClientTlsCert, certificate) || bytes.Equal(consenter.ServerTlsCert, certificate) {
				continue
			}
			newConsenters = append(newConsenters, consenter)
		}

		metadata.Consenters = newConsenters
	})
}


type ConsensusMetadataMutator func([]byte) []byte


func UpdateConsensusMetadata(network *Network, peer *Peer, orderer *Orderer, channel string, mutateMetadata ConsensusMetadataMutator) {
	config := GetConfig(network, peer, orderer, channel)
	updatedConfig := proto.Clone(config).(*common.Config)

	consensusTypeConfigValue := updatedConfig.ChannelGroup.Groups["Orderer"].Values["ConsensusType"]
	consensusTypeValue := &protosorderer.ConsensusType{}
	err := proto.Unmarshal(consensusTypeConfigValue.Value, consensusTypeValue)
	Expect(err).NotTo(HaveOccurred())

	consensusTypeValue.Metadata = mutateMetadata(consensusTypeValue.Metadata)

	updatedConfig.ChannelGroup.Groups["Orderer"].Values["ConsensusType"] = &common.ConfigValue{
		ModPolicy: "Admins",
		Value:     protoutil.MarshalOrPanic(consensusTypeValue),
	}

	UpdateOrdererConfig(network, orderer, channel, config, updatedConfig, peer, orderer)
}


func UpdateEtcdRaftMetadata(network *Network, peer *Peer, orderer *Orderer, channel string, f func(md *ectdraft_protos.ConfigMetadata)) {
	UpdateConsensusMetadata(network, peer, orderer, channel, func(originalMetadata []byte) []byte {
		metadata := &ectdraft_protos.ConfigMetadata{}
		err := proto.Unmarshal(originalMetadata, metadata)
		Expect(err).NotTo(HaveOccurred())

		f(metadata)

		newMetadata, err := proto.Marshal(metadata)
		Expect(err).NotTo(HaveOccurred())
		return newMetadata
	})
}
