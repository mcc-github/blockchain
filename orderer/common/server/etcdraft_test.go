/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var basePort = int32(8000)

func nextPort() int32 {
	return atomic.AddInt32(&basePort, 1)
}

func TestSpawnEtcdRaft(t *testing.T) {

	gt := NewGomegaWithT(t)

	cwd, err := filepath.Abs(".")
	gt.Expect(err).NotTo(HaveOccurred())

	
	blockchainRootDir, err := filepath.Abs(filepath.Join("..", "..", ".."))
	gt.Expect(err).NotTo(HaveOccurred())

	
	configtxgen, err := gexec.Build("github.com/mcc-github/blockchain/cmd/configtxgen")
	gt.Expect(err).NotTo(HaveOccurred())

	
	orderer, err := gexec.Build("github.com/mcc-github/blockchain/orderer")
	gt.Expect(err).NotTo(HaveOccurred())

	defer gexec.CleanupBuildArtifacts()

	t.Run("Bad", func(t *testing.T) {
		gt = NewGomegaWithT(t)
		tempDir, err := ioutil.TempDir("", "etcdraft-orderer-launch")
		gt.Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(tempDir)

		t.Run("Invalid bootstrap block", func(t *testing.T) {
			testEtcdRaftOSNFailureInvalidBootstrapBlock(NewGomegaWithT(t), tempDir, orderer, blockchainRootDir)
		})

		t.Run("TLS disabled single listener", func(t *testing.T) {
			testEtcdRaftOSNNoTLSSingleListener(NewGomegaWithT(t), tempDir, orderer, blockchainRootDir, configtxgen)
		})
	})

	t.Run("Good", func(t *testing.T) {
		
		
		t.Run("TLS disabled dual listener", func(t *testing.T) {
			gt = NewGomegaWithT(t)
			tempDir, err := ioutil.TempDir("", "etcdraft-orderer-launch")
			gt.Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tempDir)

			testEtcdRaftOSNNoTLSDualListener(gt, tempDir, orderer, blockchainRootDir, configtxgen)
		})

		t.Run("TLS enabled single listener", func(t *testing.T) {
			gt = NewGomegaWithT(t)
			tempDir, err := ioutil.TempDir("", "etcdraft-orderer-launch")
			gt.Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tempDir)

			testEtcdRaftOSNSuccess(gt, tempDir, configtxgen, cwd, orderer, blockchainRootDir)
		})
	})
}

func createBootstrapBlock(gt *GomegaWithT, tempDir, configtxgen, cwd string) string {
	
	genesisBlockPath := filepath.Join(tempDir, "genesis.block")
	cmd := exec.Command(configtxgen, "-channelID", "system", "-profile", "SampleDevModeEtcdRaft",
		"-outputBlock", genesisBlockPath)
	cmd.Env = append(cmd.Env, fmt.Sprintf("FABRIC_CFG_PATH=%s", filepath.Join(cwd, "testdata")))
	configtxgenProcess, err := gexec.Start(cmd, nil, nil)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Eventually(configtxgenProcess, time.Minute).Should(gexec.Exit(0))
	gt.Expect(configtxgenProcess.Err).To(gbytes.Say("Writing genesis block"))

	return genesisBlockPath
}

func testEtcdRaftOSNSuccess(gt *GomegaWithT, tempDir, configtxgen, cwd, orderer, blockchainRootDir string) {
	genesisBlockPath := createBootstrapBlock(gt, tempDir, configtxgen, cwd)

	
	ordererProcess := launchOrderer(gt, orderer, tempDir, genesisBlockPath, blockchainRootDir)
	defer ordererProcess.Kill()
	
	
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.DialTimeout = 5s"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.RPCTimeout = 7s"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.ReplicationBufferSize = 20971520"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.ReplicationPullTimeout = 5s"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.ReplicationRetryTimeout = 5s"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.ReplicationBackgroundRefreshInterval = 5m0s"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.ReplicationMaxRetries = 12"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.SendBufferSize = 10"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("General.Cluster.CertExpirationWarningThreshold = 168h0m0s"))

	
	
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("EvictionSuspicion not set, defaulting to 10m"))
	
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("Beginning to serve requests"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("becomeLeader"))
}

func testEtcdRaftOSNFailureInvalidBootstrapBlock(gt *GomegaWithT, tempDir, orderer, blockchainRootDir string) {
	
	genesisBlockPath := filepath.Join(filepath.Join("testdata", "mychannel.block"))
	genesisBlockBytes, err := ioutil.ReadFile(genesisBlockPath)
	gt.Expect(err).NotTo(HaveOccurred())

	
	genesisBlockPath = filepath.Join(tempDir, "genesis.block")
	err = ioutil.WriteFile(genesisBlockPath, genesisBlockBytes, 0644)
	gt.Expect(err).NotTo(HaveOccurred())

	
	ordererProcess := launchOrderer(gt, orderer, tempDir, genesisBlockPath, blockchainRootDir)
	defer ordererProcess.Kill()

	expectedErr := "Failed validating bootstrap block: the block isn't a system channel block because it lacks ConsortiumsConfig"
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say(expectedErr))
}

func testEtcdRaftOSNNoTLSSingleListener(gt *GomegaWithT, tempDir, orderer, blockchainRootDir string, configtxgen string) {
	cwd, err := filepath.Abs(".")
	gt.Expect(err).NotTo(HaveOccurred())

	genesisBlockPath := createBootstrapBlock(gt, tempDir, configtxgen, cwd)

	cmd := exec.Command(orderer)
	cmd.Env = []string{
		fmt.Sprintf("ORDERER_GENERAL_LISTENPORT=%d", nextPort()),
		"ORDERER_GENERAL_GENESISMETHOD=file",
		"ORDERER_GENERAL_SYSTEMCHANNEL=system",
		fmt.Sprintf("ORDERER_FILELEDGER_LOCATION=%s", filepath.Join(tempDir, "ledger")),
		fmt.Sprintf("ORDERER_GENERAL_GENESISFILE=%s", genesisBlockPath),
		fmt.Sprintf("FABRIC_CFG_PATH=%s", filepath.Join(blockchainRootDir, "sampleconfig")),
	}
	ordererProcess, err := gexec.Start(cmd, nil, nil)
	gt.Expect(err).NotTo(HaveOccurred())
	defer ordererProcess.Kill()

	expectedErr := "TLS is required for running ordering nodes of type etcdraft."
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say(expectedErr))
}

func testEtcdRaftOSNNoTLSDualListener(gt *GomegaWithT, tempDir, orderer, blockchainRootDir string, configtxgen string) {
	cwd, err := filepath.Abs(".")
	gt.Expect(err).NotTo(HaveOccurred())

	genesisBlockPath := createBootstrapBlock(gt, tempDir, configtxgen, cwd)

	cmd := exec.Command(orderer)
	cmd.Env = []string{
		fmt.Sprintf("ORDERER_GENERAL_LISTENPORT=%d", nextPort()),
		"ORDERER_GENERAL_GENESISMETHOD=file",
		"ORDERER_GENERAL_SYSTEMCHANNEL=system",
		"ORDERER_GENERAL_TLS_ENABLED=false",
		"ORDERER_OPERATIONS_TLS_ENABLED=false",
		fmt.Sprintf("ORDERER_FILELEDGER_LOCATION=%s", filepath.Join(tempDir, "ledger")),
		fmt.Sprintf("ORDERER_GENERAL_GENESISFILE=%s", genesisBlockPath),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_LISTENPORT=%d", nextPort()),
		"ORDERER_GENERAL_CLUSTER_LISTENADDRESS=127.0.0.1",
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_SERVERCERTIFICATE=%s", filepath.Join(cwd, "testdata", "tls", "server.crt")),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_SERVERPRIVATEKEY=%s", filepath.Join(cwd, "testdata", "tls", "server.key")),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=%s", filepath.Join(cwd, "testdata", "tls", "server.crt")),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=%s", filepath.Join(cwd, "testdata", "tls", "server.key")),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_ROOTCAS=[%s]", filepath.Join(cwd, "testdata", "tls", "ca.crt")),
		fmt.Sprintf("ORDERER_CONSENSUS_WALDIR=%s", filepath.Join(tempDir, "wal")),
		fmt.Sprintf("ORDERER_CONSENSUS_SNAPDIR=%s", filepath.Join(tempDir, "snapshot")),
		fmt.Sprintf("FABRIC_CFG_PATH=%s", filepath.Join(blockchainRootDir, "sampleconfig")),
		"ORDERER_OPERATIONS_LISTENADDRESS=127.0.0.1:0",
	}
	ordererProcess, err := gexec.Start(cmd, nil, nil)
	gt.Expect(err).NotTo(HaveOccurred())
	defer ordererProcess.Kill()

	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("Beginning to serve requests"))
	gt.Eventually(ordererProcess.Err, time.Minute).Should(gbytes.Say("becomeLeader"))
}

func launchOrderer(gt *GomegaWithT, orderer, tempDir, genesisBlockPath, blockchainRootDir string) *gexec.Session {
	cwd, err := filepath.Abs(".")
	gt.Expect(err).NotTo(HaveOccurred())

	
	cmd := exec.Command(orderer)
	cmd.Env = []string{
		fmt.Sprintf("ORDERER_GENERAL_LISTENPORT=%d", nextPort()),
		"ORDERER_GENERAL_GENESISMETHOD=file",
		"ORDERER_GENERAL_SYSTEMCHANNEL=system",
		"ORDERER_GENERAL_TLS_CLIENTAUTHREQUIRED=true",
		"ORDERER_GENERAL_TLS_ENABLED=true",
		"ORDERER_OPERATIONS_TLS_ENABLED=false",
		fmt.Sprintf("ORDERER_FILELEDGER_LOCATION=%s", filepath.Join(tempDir, "ledger")),
		fmt.Sprintf("ORDERER_GENERAL_GENESISFILE=%s", genesisBlockPath),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_LISTENPORT=%d", nextPort()),
		"ORDERER_GENERAL_CLUSTER_LISTENADDRESS=127.0.0.1",
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_SERVERCERTIFICATE=%s", filepath.Join(cwd, "testdata", "tls", "server.crt")),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_SERVERPRIVATEKEY=%s", filepath.Join(cwd, "testdata", "tls", "server.key")),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=%s", filepath.Join(cwd, "testdata", "tls", "server.crt")),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=%s", filepath.Join(cwd, "testdata", "tls", "server.key")),
		fmt.Sprintf("ORDERER_GENERAL_CLUSTER_ROOTCAS=[%s]", filepath.Join(cwd, "testdata", "tls", "ca.crt")),
		fmt.Sprintf("ORDERER_GENERAL_TLS_ROOTCAS=[%s]", filepath.Join(cwd, "testdata", "tls", "ca.crt")),
		fmt.Sprintf("ORDERER_GENERAL_TLS_CERTIFICATE=%s", filepath.Join(cwd, "testdata", "tls", "server.crt")),
		fmt.Sprintf("ORDERER_GENERAL_TLS_PRIVATEKEY=%s", filepath.Join(cwd, "testdata", "tls", "server.key")),
		fmt.Sprintf("ORDERER_CONSENSUS_WALDIR=%s", filepath.Join(tempDir, "wal")),
		fmt.Sprintf("ORDERER_CONSENSUS_SNAPDIR=%s", filepath.Join(tempDir, "snapshot")),
		fmt.Sprintf("FABRIC_CFG_PATH=%s", filepath.Join(blockchainRootDir, "sampleconfig")),
	}
	sess, err := gexec.Start(cmd, nil, nil)
	gt.Expect(err).NotTo(HaveOccurred())
	return sess
}
