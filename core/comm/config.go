/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"crypto/tls"
	"crypto/x509"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)


var (
	
	MaxRecvMsgSize = 100 * 1024 * 1024
	MaxSendMsgSize = 100 * 1024 * 1024
	
	DefaultKeepaliveOptions = &KeepaliveOptions{
		ClientInterval:    time.Duration(1) * time.Minute,  
		ClientTimeout:     time.Duration(20) * time.Second, 
		ServerInterval:    time.Duration(2) * time.Hour,    
		ServerTimeout:     time.Duration(20) * time.Second, 
		ServerMinInterval: time.Duration(1) * time.Minute,  
	}
	
	DefaultTLSCipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	}
	
	DefaultConnectionTimeout = 5 * time.Second
)


type ServerConfig struct {
	
	
	ConnectionTimeout time.Duration
	
	SecOpts *SecureOptions
	
	KaOpts *KeepaliveOptions
	
	
	StreamInterceptors []grpc.StreamServerInterceptor
	
	
	UnaryInterceptors []grpc.UnaryServerInterceptor
}


type ClientConfig struct {
	
	SecOpts *SecureOptions
	
	KaOpts *KeepaliveOptions
	
	
	Timeout time.Duration
}



type SecureOptions struct {
	
	
	
	VerifyCertificate func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error
	
	Certificate []byte
	
	Key []byte
	
	
	ServerRootCAs [][]byte
	
	
	ClientRootCAs [][]byte
	
	UseTLS bool
	
	RequireClientCert bool
	
	CipherSuites []uint16
}



type KeepaliveOptions struct {
	
	
	ClientInterval time.Duration
	
	
	ClientTimeout time.Duration
	
	
	ServerInterval time.Duration
	
	
	ServerTimeout time.Duration
	
	
	ServerMinInterval time.Duration
}



func ServerKeepaliveOptions(ka *KeepaliveOptions) []grpc.ServerOption {
	
	if ka == nil {
		ka = DefaultKeepaliveOptions
	}
	var serverOpts []grpc.ServerOption
	kap := keepalive.ServerParameters{
		Time:    ka.ServerInterval,
		Timeout: ka.ServerTimeout,
	}
	serverOpts = append(serverOpts, grpc.KeepaliveParams(kap))
	kep := keepalive.EnforcementPolicy{
		MinTime: ka.ServerMinInterval,
		
		PermitWithoutStream: true,
	}
	serverOpts = append(serverOpts, grpc.KeepaliveEnforcementPolicy(kep))
	return serverOpts
}



func ClientKeepaliveOptions(ka *KeepaliveOptions) []grpc.DialOption {
	
	if ka == nil {
		ka = DefaultKeepaliveOptions
	}

	var dialOpts []grpc.DialOption
	kap := keepalive.ClientParameters{
		Time:                ka.ClientInterval,
		Timeout:             ka.ClientTimeout,
		PermitWithoutStream: true,
	}
	dialOpts = append(dialOpts, grpc.WithKeepaliveParams(kap))
	return dialOpts
}
