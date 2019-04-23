/*
Copyright Hitachi America, Ltd.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/viperutil"
	"github.com/mcc-github/blockchain/core/handlers/library"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/testutil"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestStartCmd(t *testing.T) {
	defer viper.Reset()
	g := NewGomegaWithT(t)

	tempDir, err := ioutil.TempDir("", "startcmd")
	g.Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(tempDir)

	viper.Set("peer.address", "localhost:6051")
	viper.Set("peer.listenAddress", "0.0.0.0:6051")
	viper.Set("peer.chaincodeListenAddress", "0.0.0.0:6052")
	viper.Set("peer.fileSystemPath", tempDir)
	viper.Set("chaincode.executetimeout", "30s")
	viper.Set("chaincode.mode", "dev")
	viper.Set("vm.endpoint", "unix:/

	msptesttools.LoadMSPSetupForTesting()

	go func() {
		cmd := startCmd()
		assert.NoError(t, cmd.Execute(), "expected to successfully start command")
	}()

	grpcProbe := func(addr string) bool {
		c, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithInsecure())
		if err == nil {
			c.Close()
			return true
		}
		return false
	}
	g.Eventually(grpcProbe("localhost:6051")).Should(BeTrue())
}

func TestAdminHasSeparateListener(t *testing.T) {
	assert.False(t, adminHasSeparateListener("0.0.0.0:7051", ""))

	assert.Panics(t, func() {
		adminHasSeparateListener("foo", "blabla")
	})

	assert.Panics(t, func() {
		adminHasSeparateListener("0.0.0.0:7051", "blabla")
	})

	assert.False(t, adminHasSeparateListener("0.0.0.0:7051", "0.0.0.0:7051"))
	assert.False(t, adminHasSeparateListener("0.0.0.0:7051", "127.0.0.1:7051"))
	assert.True(t, adminHasSeparateListener("0.0.0.0:7051", "0.0.0.0:7055"))
}

func TestHandlerMap(t *testing.T) {
	config1 := `
  peer:
    handlers:
      authFilters:
        -
          name: filter1
          library: /opt/lib/filter1.so
        -
          name: filter2
  `
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(config1)))
	assert.NoError(t, err)

	libConf := library.Config{}
	err = viperutil.EnhancedExactUnmarshalKey("peer.handlers", &libConf)
	assert.NoError(t, err)
	assert.Len(t, libConf.AuthFilters, 2, "expected two filters")
	assert.Equal(t, "/opt/lib/filter1.so", libConf.AuthFilters[0].Library)
	assert.Equal(t, "filter2", libConf.AuthFilters[1].Name)
}

func TestComputeChaincodeEndpoint(t *testing.T) {
	
	coreConfig := &peer.Config{}
	
	
	peerAddress0 := "0.0.0.0"
	ccEndpoint, err := computeChaincodeEndpoint(coreConfig, peerAddress0)
	assert.Error(t, err)
	assert.Equal(t, "", ccEndpoint)
	
	
	peerAddress := "127.0.0.1"
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig, peerAddress)
	assert.NoError(t, err)
	assert.Equal(t, peerAddress+":7052", ccEndpoint)

	
	
	chaincodeListenPort := "8052"
	settingChaincodeListenAddress0 := "0.0.0.0:" + chaincodeListenPort
	coreConfig.ChaincodeListenAddr = settingChaincodeListenAddress0
	coreConfig.ChaincodeAddr = ""
	
	
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig, peerAddress0)
	assert.Error(t, err)
	assert.Equal(t, "", ccEndpoint)
	
	
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig, peerAddress)
	assert.NoError(t, err)
	assert.Equal(t, peerAddress+":"+chaincodeListenPort, ccEndpoint)
	
	
	settingChaincodeListenAddress := "127.0.0.1:" + chaincodeListenPort
	coreConfig.ChaincodeListenAddr = settingChaincodeListenAddress
	coreConfig.ChaincodeAddr = ""
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig, peerAddress)
	assert.NoError(t, err)
	assert.Equal(t, settingChaincodeListenAddress, ccEndpoint)
	
	
	settingChaincodeListenAddressInvalid := "abc"
	coreConfig.ChaincodeListenAddr = settingChaincodeListenAddressInvalid
	coreConfig.ChaincodeAddr = ""
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig, peerAddress)

	assert.Error(t, err)
	assert.Equal(t, "", ccEndpoint)

	
	
	
	chaincodeAddressPort := "9052"
	settingChaincodeAddress0 := "0.0.0.0:" + chaincodeAddressPort
	coreConfig.ChaincodeListenAddr = ""
	coreConfig.ChaincodeAddr = settingChaincodeAddress0
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig, peerAddress)

	assert.Error(t, err)
	assert.Equal(t, "", ccEndpoint)
	
	
	settingChaincodeAddress := "127.0.0.2:" + chaincodeAddressPort
	coreConfig.ChaincodeListenAddr = ""
	coreConfig.ChaincodeAddr = settingChaincodeAddress
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig, peerAddress)
	assert.NoError(t, err)
	assert.Equal(t, settingChaincodeAddress, ccEndpoint)
	
	
	settingChaincodeAddressInvalid := "bcd"
	coreConfig.ChaincodeListenAddr = ""
	coreConfig.ChaincodeAddr = settingChaincodeAddressInvalid
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig, peerAddress)
	assert.Error(t, err)
	assert.Equal(t, "", ccEndpoint)

	
	
}

func TestGetDockerHostConfig(t *testing.T) {
	testutil.SetupTestConfig()
	hostConfig := getDockerHostConfig()
	assert.NotNil(t, hostConfig)
	assert.Equal(t, "host", hostConfig.NetworkMode)
	assert.Equal(t, "json-file", hostConfig.LogConfig.Type)
	assert.Equal(t, "50m", hostConfig.LogConfig.Config["max-size"])
	assert.Equal(t, "5", hostConfig.LogConfig.Config["max-file"])
	assert.Equal(t, int64(1024*1024*1024*2), hostConfig.Memory)
	assert.Equal(t, int64(0), hostConfig.CPUShares)
}
