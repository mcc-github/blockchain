/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package peer

import (
	"crypto/tls"
	"fmt"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/comm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCacheConfigurationNegative(t *testing.T) {
	
	viper.Set("peer.addressAutoDetect", true)
	viper.Set("peer.address", "testing.com")
	_, err := GlobalConfig()
	assert.Error(t, err, "Expected error for bad configuration")

	viper.Set("peer.addressAutoDetect", false)
	viper.Set("peer.address", "")
	_, err = GlobalConfig()
	assert.Error(t, err, "Expected error for bad configuration")

	viper.Set("peer.address", "wrongAddress")
	_, err = GlobalConfig()
	assert.Error(t, err, "Expected error for bad configuration")

}

func TestConfiguration(t *testing.T) {
	
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		t.Fatal("Failed to get interface addresses")
	}

	var ips []string
	for _, address := range addresses {
		
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			ips = append(ips, ip.IP.String()+":7051")
			t.Logf("found interface address [%s]", ip.IP.String())
		}
	}

	
	localIP, err := GetLocalIP()
	assert.NoError(t, err)

	var tests = []struct {
		name                string
		settings            map[string]interface{}
		validAddresses      []string
		invalidAddresses    []string
		expectedPeerAddress string
	}{
		{
			name: "test1",
			settings: map[string]interface{}{
				"peer.addressAutoDetect": false,
				"peer.address":           "testing.com:7051",
				"peer.id":                "testPeer",
			},
			validAddresses:      []string{"testing.com:7051"},
			invalidAddresses:    ips,
			expectedPeerAddress: "testing.com:7051",
		},
		{
			name: "test2",
			settings: map[string]interface{}{
				"peer.addressAutoDetect": true,
				"peer.address":           "testing.com:7051",
				"peer.id":                "testPeer",
			},
			validAddresses:      ips,
			invalidAddresses:    []string{"testing.com:7051"},
			expectedPeerAddress: net.JoinHostPort(localIP, "7051"),
		},
		{
			name: "test3",
			settings: map[string]interface{}{
				"peer.addressAutoDetect": false,
				"peer.address":           "0.0.0.0:7051",
				"peer.id":                "testPeer",
			},
			validAddresses:      []string{fmt.Sprintf("%s:7051", localIP)},
			invalidAddresses:    []string{"0.0.0.0:7051"},
			expectedPeerAddress: net.JoinHostPort(localIP, "7051"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.settings {
				viper.Set(k, v)
			}
			
			coreConfig, err := GlobalConfig()
			assert.NoError(t, err, "GetPeerEndpoint returned unexpected error")
			assert.Equal(t, test.settings["peer.id"], coreConfig.PeerEndpoint.Id.Name, "GetPeerEndpoint returned the wrong peer ID")
			assert.Equal(t, test.expectedPeerAddress, coreConfig.PeerEndpoint.Address, "GetPeerEndpoint returned the wrong peer address")
		})
	}
}

func TestGetServerConfig(t *testing.T) {
	
	viper.Set("peer.tls.enabled", false)
	sc, _ := GetServerConfig()
	assert.Equal(t, false, sc.SecOpts.UseTLS, "ServerConfig.SecOpts.UseTLS should be false")

	
	assert.Equal(t, comm.DefaultKeepaliveOptions, sc.KaOpts, "ServerConfig.KaOpts should be set to default values")
	viper.Set("peer.keepalive.interval", "60m")
	sc, _ = GetServerConfig()
	assert.Equal(t, time.Duration(60)*time.Minute, sc.KaOpts.ServerInterval, "ServerConfig.KaOpts.ServerInterval should be set to 60 min")
	viper.Set("peer.keepalive.timeout", "30s")
	sc, _ = GetServerConfig()
	assert.Equal(t, time.Duration(30)*time.Second, sc.KaOpts.ServerTimeout, "ServerConfig.KaOpts.ServerTimeout should be set to 30 sec")
	viper.Set("peer.keepalive.minInterval", "2m")
	sc, _ = GetServerConfig()
	assert.Equal(t, time.Duration(2)*time.Minute, sc.KaOpts.ServerMinInterval, "ServerConfig.KaOpts.ServerMinInterval should be set to 2 min")

	
	viper.Set("peer.tls.enabled", true)
	viper.Set("peer.tls.cert.file", filepath.Join("testdata", "Org1-server1-cert.pem"))
	viper.Set("peer.tls.key.file", filepath.Join("testdata", "Org1-server1-key.pem"))
	viper.Set("peer.tls.rootcert.file", filepath.Join("testdata", "Org1-cert.pem"))
	sc, _ = GetServerConfig()
	assert.Equal(t, true, sc.SecOpts.UseTLS, "ServerConfig.SecOpts.UseTLS should be true")
	assert.Equal(t, false, sc.SecOpts.RequireClientCert, "ServerConfig.SecOpts.RequireClientCert should be false")
	viper.Set("peer.tls.clientAuthRequired", true)
	viper.Set("peer.tls.clientRootCAs.files", []string{
		filepath.Join("testdata", "Org1-cert.pem"),
		filepath.Join("testdata", "Org2-cert.pem"),
	})
	sc, _ = GetServerConfig()
	assert.Equal(t, true, sc.SecOpts.RequireClientCert, "ServerConfig.SecOpts.RequireClientCert should be true")
	assert.Equal(t, 2, len(sc.SecOpts.ClientRootCAs), "ServerConfig.SecOpts.ClientRootCAs should contain 2 entries")

	
	viper.Set("peer.tls.rootcert.file", filepath.Join("testdata", "Org11-cert.pem"))
	_, err := GetServerConfig()
	assert.Error(t, err, "GetServerConfig should return error with bad root cert path")
	viper.Set("peer.tls.cert.file", filepath.Join("testdata", "Org11-cert.pem"))
	_, err = GetServerConfig()
	assert.Error(t, err, "GetServerConfig should return error with bad tls cert path")

	
	viper.Set("peer.tls.enabled", false)
	viper.Set("peer.tls.clientAuthRequired", false)
}

func TestGetClientCertificate(t *testing.T) {
	viper.Set("peer.tls.key.file", "")
	viper.Set("peer.tls.cert.file", "")
	viper.Set("peer.tls.clientKey.file", "")
	viper.Set("peer.tls.clientCert.file", "")

	
	_, err := GetClientCertificate()
	assert.Error(t, err)

	viper.Set("peer.tls.key.file", "")
	viper.Set("peer.tls.cert.file", filepath.Join("testdata", "Org1-server1-cert.pem"))
	
	_, err = GetClientCertificate()
	assert.Error(t, err)

	viper.Set("peer.tls.key.file", filepath.Join("testdata", "Org1-server1-key.pem"))
	viper.Set("peer.tls.cert.file", "")
	
	_, err = GetClientCertificate()
	assert.Error(t, err)

	
	
	viper.Set("peer.tls.key.file", filepath.Join("testdata", "Org1-server1-key.pem"))
	viper.Set("peer.tls.cert.file", filepath.Join("testdata", "Org1-server1-cert.pem"))

	
	viper.Set("peer.tls.clientKey.file", filepath.Join("testdata", "Org2-server1-key.pem"))
	_, err = GetClientCertificate()
	assert.Error(t, err)

	
	viper.Set("peer.tls.clientKey.file", "")
	viper.Set("peer.tls.clientCert.file", filepath.Join("testdata", "Org2-server1-cert.pem"))
	_, err = GetClientCertificate()
	assert.Error(t, err)

	
	expected, err := tls.LoadX509KeyPair(
		filepath.Join("testdata", "Org2-server1-cert.pem"),
		filepath.Join("testdata", "Org2-server1-key.pem"),
	)
	if err != nil {
		t.Fatalf("Failed to load test certificate (%s)", err)
	}
	viper.Set("peer.tls.clientKey.file", filepath.Join("testdata", "Org2-server1-key.pem"))
	cert, err := GetClientCertificate()
	assert.NoError(t, err)
	assert.Equal(t, expected, cert)

	
	
	viper.Set("peer.tls.clientKey.file", "")
	viper.Set("peer.tls.clientCert.file", "")
	expected, err = tls.LoadX509KeyPair(
		filepath.Join("testdata", "Org1-server1-cert.pem"),
		filepath.Join("testdata", "Org1-server1-key.pem"),
	)
	if err != nil {
		t.Fatalf("Failed to load test certificate (%s)", err)
	}
	cert, err = GetClientCertificate()
	assert.NoError(t, err)
	assert.Equal(t, expected, cert)
}

func TestGlobalConfig(t *testing.T) {
	
	viper.Set("peer.addressAutoDetect", false)
	viper.Set("peer.address", "localhost:8080")
	viper.Set("peer.id", "testPeerID")
	viper.Set("peer.localMspId", "SampleOrg")
	viper.Set("peer.listenAddress", "0.0.0.0:7051")
	viper.Set("peer.authentication.timewindow", "15m")
	viper.Set("peer.tls.enabled", "false")
	viper.Set("peer.networkId", "testNetwork")
	viper.Set("peer.limits.concurrency.qscc", 5000)
	viper.Set("peer.discovery.enabled", true)
	viper.Set("peer.profile.enabled", false)
	viper.Set("peer.profile.listenAddress", "peer.authentication.timewindow")
	viper.Set("peer.discovery.orgMembersAllowedAccess", false)
	viper.Set("peer.discovery.authCacheEnabled", true)
	viper.Set("peer.discovery.authCacheMaxSize", 1000)
	viper.Set("peer.discovery.authCachePurgeRetentionRatio", 0.75)
	viper.Set("peer.chaincodeListenAddress", "0.0.0.0:7052")
	viper.Set("peer.chaincodeAddress", "0.0.0.0:7052")
	viper.Set("peer.adminService.listenAddress", "0.0.0.0:7055")

	viper.Set("vm.endpoint", "unix:/
	viper.Set("vm.docker.tls.enabled", false)
	viper.Set("vm.docker.attachStdout", false)

	viper.Set("chaincode.pull", false)

	coreConfig, err := GlobalConfig()
	assert.NoError(t, err)

	assert.Equal(t, coreConfig.LocalMspID, "SampleOrg")
	assert.Equal(t, coreConfig.ListenAddress, "0.0.0.0:7051")
	assert.Equal(t, coreConfig.AuthenticationTimeWindow, 15*time.Minute)
	assert.Equal(t, coreConfig.PeerTLSEnabled, false)
	assert.Equal(t, coreConfig.NetworkID, "testNetwork")
	assert.Equal(t, coreConfig.LimitsConcurrencyQSCC, 5000)
	assert.Equal(t, coreConfig.DiscoveryEnabled, true)
	assert.Equal(t, coreConfig.ProfileEnabled, false)
	assert.Equal(t, coreConfig.ProfileListenAddress, "peer.authentication.timewindow")
	assert.Equal(t, coreConfig.DiscoveryOrgMembersAllowed, false)
	assert.Equal(t, coreConfig.DiscoveryAuthCacheEnabled, true)
	assert.Equal(t, coreConfig.DiscoveryAuthCacheMaxSize, 1000)
	assert.Equal(t, coreConfig.DiscoveryAuthCachePurgeRetentionRatio, 0.75)
	assert.Equal(t, coreConfig.ChaincodeListenAddr, "0.0.0.0:7052")
	assert.Equal(t, coreConfig.ChaincodeAddr, "0.0.0.0:7052")
	assert.Equal(t, coreConfig.AdminListenAddr, "0.0.0.0:7055")

	assert.Equal(t, coreConfig.VMEndpoint, "unix:/
	assert.Equal(t, coreConfig.VMDockerTLSEnabled, false)
	assert.Equal(t, coreConfig.VMDockerAttachStdout, false)

	assert.Equal(t, coreConfig.ChaincodePull, false)

	assert.Equal(t, coreConfig.PeerAddress, "localhost:8080")
	assert.Equal(t, coreConfig.PeerID, "testPeerID")
	assert.Equal(t, coreConfig.PeerEndpoint.Id.Name, "testPeerID")
	assert.Equal(t, coreConfig.PeerEndpoint.Address, "localhost:8080")

}
