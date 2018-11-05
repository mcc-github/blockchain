/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/














package peer

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"

	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/config"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)


var configurationCached = false



var localAddress string
var localAddressError error
var peerEndpoint *pb.PeerEndpoint
var peerEndpointError error






func CacheConfiguration() (err error) {
	
	getLocalAddress := func() (string, error) {
		peerAddress := viper.GetString("peer.address")
		if peerAddress == "" {
			return "", fmt.Errorf("peer.address isn't set")
		}
		host, port, err := net.SplitHostPort(peerAddress)
		if err != nil {
			return "", errors.Errorf("peer.address isn't in host:port format: %s", peerAddress)
		}

		autoDetectedIPAndPort := net.JoinHostPort(GetLocalIP(), port)
		peerLogger.Info("Auto-detected peer address:", autoDetectedIPAndPort)
		
		
		if ip := net.ParseIP(host); ip != nil && ip.IsUnspecified() {
			peerLogger.Info("Host is", host, ", falling back to auto-detected address:", autoDetectedIPAndPort)
			return autoDetectedIPAndPort, nil
		}

		if viper.GetBool("peer.addressAutoDetect") {
			peerLogger.Info("Auto-detect flag is set, returning", autoDetectedIPAndPort)
			return autoDetectedIPAndPort, nil
		}
		peerLogger.Info("Returning", peerAddress)
		return peerAddress, nil

	}

	
	getPeerEndpoint := func() (*pb.PeerEndpoint, error) {
		var peerAddress string
		peerAddress, err := getLocalAddress()
		if err != nil {
			return nil, err
		}
		return &pb.PeerEndpoint{Id: &pb.PeerID{Name: viper.GetString("peer.id")}, Address: peerAddress}, nil
	}

	localAddress, localAddressError = getLocalAddress()
	peerEndpoint, _ = getPeerEndpoint()

	configurationCached = true

	if localAddressError != nil {
		return localAddressError
	}
	return
}


func cacheConfiguration() {
	if err := CacheConfiguration(); err != nil {
		peerLogger.Errorf("Execution continues after CacheConfiguration() failure : %s", err)
	}
}




func GetLocalAddress() (string, error) {
	if !configurationCached {
		cacheConfiguration()
	}
	return localAddress, localAddressError
}


func GetPeerEndpoint() (*pb.PeerEndpoint, error) {
	if !configurationCached {
		cacheConfiguration()
	}
	return peerEndpoint, peerEndpointError
}


func GetServerConfig() (comm.ServerConfig, error) {
	secureOptions := &comm.SecureOptions{
		UseTLS: viper.GetBool("peer.tls.enabled"),
	}
	serverConfig := comm.ServerConfig{SecOpts: secureOptions}
	if secureOptions.UseTLS {
		
		serverKey, err := ioutil.ReadFile(config.GetPath("peer.tls.key.file"))
		if err != nil {
			return serverConfig, fmt.Errorf("error loading TLS key (%s)", err)
		}
		serverCert, err := ioutil.ReadFile(config.GetPath("peer.tls.cert.file"))
		if err != nil {
			return serverConfig, fmt.Errorf("error loading TLS certificate (%s)", err)
		}
		secureOptions.Certificate = serverCert
		secureOptions.Key = serverKey
		secureOptions.RequireClientCert = viper.GetBool("peer.tls.clientAuthRequired")
		if secureOptions.RequireClientCert {
			var clientRoots [][]byte
			for _, file := range viper.GetStringSlice("peer.tls.clientRootCAs.files") {
				clientRoot, err := ioutil.ReadFile(
					config.TranslatePath(filepath.Dir(viper.ConfigFileUsed()), file))
				if err != nil {
					return serverConfig,
						fmt.Errorf("error loading client root CAs (%s)", err)
				}
				clientRoots = append(clientRoots, clientRoot)
			}
			secureOptions.ClientRootCAs = clientRoots
		}
		
		if config.GetPath("peer.tls.rootcert.file") != "" {
			rootCert, err := ioutil.ReadFile(config.GetPath("peer.tls.rootcert.file"))
			if err != nil {
				return serverConfig, fmt.Errorf("error loading TLS root certificate (%s)", err)
			}
			secureOptions.ServerRootCAs = [][]byte{rootCert}
		}
	}
	
	serverConfig.KaOpts = comm.DefaultKeepaliveOptions
	
	if viper.IsSet("peer.keepalive.minInterval") {
		serverConfig.KaOpts.ServerMinInterval = viper.GetDuration("peer.keepalive.minInterval")
	}

	serverConfig.MetricsProvider = nil

	return serverConfig, nil
}



func GetClientCertificate() (tls.Certificate, error) {
	cert := tls.Certificate{}

	keyPath := viper.GetString("peer.tls.clientKey.file")
	certPath := viper.GetString("peer.tls.clientCert.file")

	if keyPath != "" || certPath != "" {
		
		if keyPath == "" || certPath == "" {
			return cert, errors.New("peer.tls.clientKey.file and " +
				"peer.tls.clientCert.file must both be set or must both be empty")
		}
		keyPath = config.GetPath("peer.tls.clientKey.file")
		certPath = config.GetPath("peer.tls.clientCert.file")

	} else {
		
		keyPath = viper.GetString("peer.tls.key.file")
		certPath = viper.GetString("peer.tls.key.file")

		if keyPath != "" || certPath != "" {
			
			if keyPath == "" || certPath == "" {
				return cert, errors.New("peer.tls.key.file and " +
					"peer.tls.cert.file must both be set or must both be empty")
			}
			keyPath = config.GetPath("peer.tls.key.file")
			certPath = config.GetPath("peer.tls.cert.file")
		} else {
			return cert, errors.New("must set either " +
				"[peer.tls.key.file and peer.tls.cert.file] or " +
				"[peer.tls.clientKey.file and peer.tls.clientCert.file]" +
				"when peer.tls.clientAuthEnabled is set to true")
		}
	}
	
	clientKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return cert, errors.WithMessage(err,
			"error loading client TLS key")
	}
	clientCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return cert, errors.WithMessage(err,
			"error loading client TLS certificate")
	}
	cert, err = tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return cert, errors.WithMessage(err,
			"error parsing client TLS key pair")
	}
	return cert, nil
}
