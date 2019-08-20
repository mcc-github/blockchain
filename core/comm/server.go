/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"sync"
	"sync/atomic"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCServer struct {
	
	address string
	
	listener net.Listener
	
	server *grpc.Server
	
	
	serverCertificate atomic.Value
	
	lock *sync.Mutex
	
	
	clientRootCAs map[string]*x509.Certificate
	
	tlsConfig *tls.Config
	
	healthServer *health.Server
}



func NewGRPCServer(address string, serverConfig ServerConfig) (*GRPCServer, error) {
	if address == "" {
		return nil, errors.New("missing address parameter")
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
	if secureConfig.UseTLS {
		
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

			if serverConfig.SecOpts.TimeShift > 0 {
				timeShift := serverConfig.SecOpts.TimeShift
				grpcServer.tlsConfig.Time = func() time.Time {
					return time.Now().Add((-1) * timeShift)
				}
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

			
			creds := NewServerTransportCredentials(grpcServer.tlsConfig, serverConfig.Logger)
			serverOpts = append(serverOpts, grpc.Creds(creds))
		} else {
			return nil, errors.New("serverConfig.SecOpts must contain both Key and Certificate when UseTLS is true")
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
	
	if len(serverConfig.StreamInterceptors) > 0 {
		serverOpts = append(
			serverOpts,
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(serverConfig.StreamInterceptors...)),
		)
	}
	if len(serverConfig.UnaryInterceptors) > 0 {
		serverOpts = append(
			serverOpts,
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(serverConfig.UnaryInterceptors...)),
		)
	}

	if serverConfig.MetricsProvider != nil {
		sh := NewServerStatsHandler(serverConfig.MetricsProvider)
		serverOpts = append(serverOpts, grpc.StatsHandler(sh))
	}

	grpcServer.server = grpc.NewServer(serverOpts...)

	if serverConfig.HealthCheckEnabled {
		grpcServer.healthServer = health.NewServer()
		healthpb.RegisterHealthServer(grpcServer.server, grpcServer.healthServer)
	}

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
	
	if gServer.healthServer != nil {
		for name := range gServer.server.GetServiceInfo() {
			gServer.healthServer.SetServingStatus(
				name,
				healthpb.HealthCheckResponse_SERVING,
			)
		}
		gServer.healthServer.SetServingStatus(
			"",
			healthpb.HealthCheckResponse_SERVING,
		)
	}
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

	certs, subjects, err := pemToX509Certs(clientRoot)
	if err != nil {
		return errors.WithMessage(err, "failed to append client root certificate(s)")
	}

	if len(certs) < 1 {
		return errors.New("no client root certificates found")
	}

	for i, cert := range certs {
		
		gServer.tlsConfig.ClientCAs.AddCert(cert)
		
		gServer.clientRootCAs[subjects[i]] = cert
	}
	return nil
}



func (gServer *GRPCServer) SetClientRootCAs(clientRoots [][]byte) error {
	gServer.lock.Lock()
	defer gServer.lock.Unlock()

	
	clientRootCAs := make(map[string]*x509.Certificate)
	for _, clientRoot := range clientRoots {
		certs, subjects, err := pemToX509Certs(clientRoot)
		if err != nil {
			return errors.WithMessage(err, "failed to set client root certificate(s)")
		}

		for i, cert := range certs {
			clientRootCAs[subjects[i]] = cert
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
