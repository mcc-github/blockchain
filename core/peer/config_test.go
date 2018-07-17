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
	cacheConfiguration()
	err := CacheConfiguration()
	assert.Error(t, err, "Expected error for bad configuration")
}

func TestConfiguration(t *testing.T) {

	var ips []string
	
	if addresses, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addresses {
			
			if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
				ips = append(ips, ip.IP.String()+":7051")
				t.Logf("found interface address [%s]", ip.IP.String())
			}
		}
	} else {
		t.Fatal("Failed to get interface addresses")
	}

	var tests = []struct {
		name             string
		settings         map[string]interface{}
		validAddresses   []string
		invalidAddresses []string
	}{
		{
			name: "test1",
			settings: map[string]interface{}{
				"peer.addressAutoDetect": false,
				"peer.address":           "testing.com:7051",
				"peer.id":                "testPeer",
			},
			validAddresses:   []string{"testing.com:7051"},
			invalidAddresses: ips,
		},
		{
			name: "test2",
			settings: map[string]interface{}{
				"peer.addressAutoDetect": true,
				"peer.address":           "testing.com:7051",
				"peer.id":                "testPeer",
			},
			validAddresses:   ips,
			invalidAddresses: []string{"testing.com:7051"},
		},
		{
			name: "test3",
			settings: map[string]interface{}{
				"peer.addressAutoDetect": false,
				"peer.address":           "0.0.0.0:7051",
				"peer.id":                "testPeer",
			},
			validAddresses:   []string{fmt.Sprintf("%s:7051", GetLocalIP())},
			invalidAddresses: []string{"0.0.0.0:7051"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.settings {
				viper.Set(k, v)
			}
			
			configurationCached = false
			
			address, err := GetLocalAddress()
			assert.NoError(t, err, "GetLocalAddress returned unexpected error")
			assert.Contains(t, test.validAddresses, address,
				"GetLocalAddress returned unexpected address")
			assert.NotContains(t, test.invalidAddresses, address,
				"GetLocalAddress returned invalid address")
			
			configurationCached = false
			
			pe, err := GetPeerEndpoint()
			assert.NoError(t, err, "GetPeerEndpoint returned unexpected error")
			assert.Equal(t, test.settings["peer.id"], pe.Id.Name,
				"GetPeerEndpoint returned the wrong peer ID")
			assert.Equal(t, address, pe.Address,
				"GetPeerEndpoint returned the wrong peer address")

			
			err = CacheConfiguration()
			assert.NoError(t, err, "CacheConfiguration should not have returned an err")
			
			
			address, err = GetLocalAddress()
			assert.NoError(t, err, "GetLocalAddress should not have returned error")
			assert.Contains(t, test.validAddresses, address,
				"GetLocalAddress returned unexpected address")
			assert.NotContains(t, test.invalidAddresses, address,
				"GetLocalAddress returned invalid address")
			
			pe, err = GetPeerEndpoint()
			assert.NoError(t, err, "GetPeerEndpoint returned unexpected error")
			assert.Equal(t, test.settings["peer.id"], pe.Id.Name,
				"GetPeerEndpoint returned the wrong peer ID")
			assert.Equal(t, address, pe.Address,
				"GetPeerEndpoint returned the wrong peer address")
		})
	}
}

func TestGetServerConfig(t *testing.T) {

	
	viper.Set("peer.tls.enabled", false)
	sc, _ := GetServerConfig()
	assert.Equal(t, false, sc.SecOpts.UseTLS,
		"ServerConfig.SecOpts.UseTLS should be false")

	
	assert.Equal(t, comm.DefaultKeepaliveOptions, sc.KaOpts,
		"ServerConfig.KaOpts should be set to default values")
	viper.Set("peer.keepalive.minInterval", "2m")
	sc, _ = GetServerConfig()
	assert.Equal(t, time.Duration(2)*time.Minute, sc.KaOpts.ServerMinInterval,
		"ServerConfig.KaOpts.ServerMinInterval should be set to 2 min")

	
	viper.Set("peer.tls.enabled", true)
	viper.Set("peer.tls.cert.file", filepath.Join("testdata", "Org1-server1-cert.pem"))
	viper.Set("peer.tls.key.file", filepath.Join("testdata", "Org1-server1-key.pem"))
	viper.Set("peer.tls.rootcert.file", filepath.Join("testdata", "Org1-cert.pem"))
	sc, _ = GetServerConfig()
	assert.Equal(t, true, sc.SecOpts.UseTLS, "ServerConfig.SecOpts.UseTLS should be true")
	assert.Equal(t, false, sc.SecOpts.RequireClientCert,
		"ServerConfig.SecOpts.RequireClientCert should be false")
	viper.Set("peer.tls.clientAuthRequired", true)
	viper.Set("peer.tls.clientRootCAs.files",
		[]string{filepath.Join("testdata", "Org1-cert.pem"),
			filepath.Join("testdata", "Org2-cert.pem")})
	sc, _ = GetServerConfig()
	assert.Equal(t, true, sc.SecOpts.RequireClientCert,
		"ServerConfig.SecOpts.RequireClientCert should be true")
	assert.Equal(t, 2, len(sc.SecOpts.ClientRootCAs),
		"ServerConfig.SecOpts.ClientRootCAs should contain 2 entries")

	
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
	viper.Set("peer.tls.cert.file",
		filepath.Join("testdata", "Org1-server1-cert.pem"))
	
	_, err = GetClientCertificate()
	assert.Error(t, err)

	viper.Set("peer.tls.key.file",
		filepath.Join("testdata", "Org1-server1-key.pem"))
	viper.Set("peer.tls.cert.file", "")
	
	_, err = GetClientCertificate()
	assert.Error(t, err)

	
	
	viper.Set("peer.tls.key.file",
		filepath.Join("testdata", "Org1-server1-key.pem"))
	viper.Set("peer.tls.cert.file",
		filepath.Join("testdata", "Org1-server1-cert.pem"))

	
	viper.Set("peer.tls.clientKey.file",
		filepath.Join("testdata", "Org2-server1-key.pem"))
	_, err = GetClientCertificate()
	assert.Error(t, err)

	
	viper.Set("peer.tls.clientKey.file", "")
	viper.Set("peer.tls.clientCert.file",
		filepath.Join("testdata", "Org2-server1-cert.pem"))
	_, err = GetClientCertificate()
	assert.Error(t, err)

	
	expected, err := tls.LoadX509KeyPair(
		filepath.Join("testdata", "Org2-server1-cert.pem"),
		filepath.Join("testdata", "Org2-server1-key.pem"))
	if err != nil {
		t.Fatalf("Failed to load test certificate (%s)", err)
	}
	viper.Set("peer.tls.clientKey.file",
		filepath.Join("testdata", "Org2-server1-key.pem"))
	cert, err := GetClientCertificate()
	assert.NoError(t, err)
	assert.Equal(t, expected, cert)

	
	
	viper.Set("peer.tls.clientKey.file", "")
	viper.Set("peer.tls.clientCert.file", "")
	expected, err = tls.LoadX509KeyPair(
		filepath.Join("testdata", "Org1-server1-cert.pem"),
		filepath.Join("testdata", "Org1-server1-key.pem"))
	if err != nil {
		t.Fatalf("Failed to load test certificate (%s)", err)
	}
	cert, err = GetClientCertificate()
	assert.NoError(t, err)
	assert.Equal(t, expected, cert)
}
