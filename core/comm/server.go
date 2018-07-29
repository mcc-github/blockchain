/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	
	address string
	
	listener net.Listener
	
	server *grpc.Server
	
	
	serverCertificate atomic.Value
	
	serverKeyPEM []byte
	
	lock *sync.Mutex
	
	
	clientRootCAs map[string]*x509.Certificate
	
	tlsConfig *tls.Config
}



func NewGRPCServer(address string, serverConfig ServerConfig) (*GRPCServer, error) {
	if address == "" {
		return nil, errors.New("Missing address parameter")
	}
	
	lis, err := net.Listen("tcp", address)

	if err != nil {
		return nil, err
	}
	return NewGRPCServerFromListener(lis, serverConfig)
}



func NewGRPCServerFromListener(listener net.Listener, serverConfig ServerConfig) (*GRPCServer, error) {
	grpcServer := &GRPCServer{
		address:  listener.Addr().String(),
		listener: listener,
		lock:     &sync.Mutex{},
	}

	
	var serverOpts []grpc.ServerOption
	
	secureConfig := serverConfig.SecOpts
	if secureConfig != nil && secureConfig.UseTLS {
		
		if secureConfig.Key != nil && secureConfig.Certificate != nil {
			
			cert, err := tls.X509KeyPair(secureConfig.Certificate, secureConfig.Key)
			if err != nil {
				return nil, err
			}
			grpcServer.serverCertificate.Store(cert)

			
			if len(secureConfig.CipherSuites) == 0 {
				secureConfig.CipherSuites = DefaultTLSCipherSuites
			}
			getCert := func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert := grpcServer.serverCertificate.Load().(tls.Certificate)
				return &cert, nil
			}
			
			grpcServer.tlsConfig = &tls.Config{
				VerifyPeerCertificate:  secureConfig.VerifyCertificate,
				GetCertificate:         getCert,
				SessionTicketsDisabled: true,
				CipherSuites:           secureConfig.CipherSuites,
			}
			grpcServer.tlsConfig.ClientAuth = tls.RequestClientCert
			
			if secureConfig.RequireClientCert {
				
				grpcServer.tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
				
				if len(secureConfig.ClientRootCAs) > 0 {
					grpcServer.clientRootCAs = make(map[string]*x509.Certificate)
					grpcServer.tlsConfig.ClientCAs = x509.NewCertPool()
					for _, clientRootCA := range secureConfig.ClientRootCAs {
						err = grpcServer.appendClientRootCA(clientRootCA)
						if err != nil {
							return nil, err
						}
					}
				}
			}

			
			creds := NewServerTransportCredentials(grpcServer.tlsConfig)
			serverOpts = append(serverOpts, grpc.Creds(creds))
		} else {
			return nil, errors.New("serverConfig.SecOpts must contain both Key and " +
				"Certificate when UseTLS is true")
		}
	}
	
	serverOpts = append(serverOpts, grpc.MaxSendMsgSize(MaxSendMsgSize))
	serverOpts = append(serverOpts, grpc.MaxRecvMsgSize(MaxRecvMsgSize))
	
	serverOpts = append(serverOpts, ServerKeepaliveOptions(serverConfig.KaOpts)...)
	
	if serverConfig.ConnectionTimeout <= 0 {
		serverConfig.ConnectionTimeout = DefaultConnectionTimeout
	}
	serverOpts = append(
		serverOpts,
		grpc.ConnectionTimeout(serverConfig.ConnectionTimeout))

	grpcServer.server = grpc.NewServer(serverOpts...)

	return grpcServer, nil
}


func (gServer *GRPCServer) SetServerCertificate(cert tls.Certificate) {
	gServer.serverCertificate.Store(cert)
}


func (gServer *GRPCServer) Address() string {
	return gServer.address
}


func (gServer *GRPCServer) Listener() net.Listener {
	return gServer.listener
}


func (gServer *GRPCServer) Server() *grpc.Server {
	return gServer.server
}


func (gServer *GRPCServer) ServerCertificate() tls.Certificate {
	return gServer.serverCertificate.Load().(tls.Certificate)
}



func (gServer *GRPCServer) TLSEnabled() bool {
	return gServer.tlsConfig != nil
}



func (gServer *GRPCServer) MutualTLSRequired() bool {
	return gServer.tlsConfig != nil &&
		gServer.tlsConfig.ClientAuth ==
			tls.RequireAndVerifyClientCert
}


func (gServer *GRPCServer) Start() error {
	return gServer.server.Serve(gServer.listener)
}


func (gServer *GRPCServer) Stop() {
	gServer.server.Stop()
}



func (gServer *GRPCServer) AppendClientRootCAs(clientRoots [][]byte) error {
	gServer.lock.Lock()
	defer gServer.lock.Unlock()
	for _, clientRoot := range clientRoots {
		err := gServer.appendClientRootCA(clientRoot)
		if err != nil {
			return err
		}
	}
	return nil
}


func (gServer *GRPCServer) appendClientRootCA(clientRoot []byte) error {

	errMsg := "Failed to append client root certificate(s): %s"
	
	certs, subjects, err := pemToX509Certs(clientRoot)
	if err != nil {
		return fmt.Errorf(errMsg, err.Error())
	}

	if len(certs) < 1 {
		return fmt.Errorf(errMsg, "No client root certificates found")
	}

	for i, cert := range certs {
		
		gServer.tlsConfig.ClientCAs.AddCert(cert)
		
		gServer.clientRootCAs[subjects[i]] = cert
	}
	return nil
}



func (gServer *GRPCServer) RemoveClientRootCAs(clientRoots [][]byte) error {
	gServer.lock.Lock()
	defer gServer.lock.Unlock()
	
	for _, clientRoot := range clientRoots {
		err := gServer.removeClientRootCA(clientRoot)
		if err != nil {
			return err
		}
	}

	
	certPool := x509.NewCertPool()
	for _, clientRoot := range gServer.clientRootCAs {
		certPool.AddCert(clientRoot)
	}

	
	gServer.tlsConfig.ClientCAs = certPool
	return nil
}


func (gServer *GRPCServer) removeClientRootCA(clientRoot []byte) error {

	errMsg := "Failed to remove client root certificate(s): %s"
	
	certs, subjects, err := pemToX509Certs(clientRoot)
	if err != nil {
		return fmt.Errorf(errMsg, err.Error())
	}

	if len(certs) < 1 {
		return fmt.Errorf(errMsg, "No client root certificates found")
	}

	for i, subject := range subjects {
		
		
		if certs[i].Equal(gServer.clientRootCAs[subject]) {
			delete(gServer.clientRootCAs, subject)
		}
	}
	return nil
}



func (gServer *GRPCServer) SetClientRootCAs(clientRoots [][]byte) error {
	gServer.lock.Lock()
	defer gServer.lock.Unlock()

	errMsg := "Failed to set client root certificate(s): %s"

	
	clientRootCAs := make(map[string]*x509.Certificate)
	for _, clientRoot := range clientRoots {
		certs, subjects, err := pemToX509Certs(clientRoot)
		if err != nil {
			return fmt.Errorf(errMsg, err.Error())
		}
		if len(certs) >= 1 {
			for i, cert := range certs {
				
				clientRootCAs[subjects[i]] = cert
			}
		}
	}

	
	certPool := x509.NewCertPool()
	for _, clientRoot := range clientRootCAs {
		certPool.AddCert(clientRoot)
	}
	
	gServer.clientRootCAs = clientRootCAs
	
	gServer.tlsConfig.ClientCAs = certPool
	return nil
}
