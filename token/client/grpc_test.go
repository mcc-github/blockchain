/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client_test

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/mcc-github/blockchain/token/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var _ = Describe("GRPCClient", func() {
	var (
		connConfig *client.ConnectionConfig
		serverCert tls.Certificate
		endpoint   string
		listener   net.Listener
	)

	BeforeEach(func() {
		
		var err error
		listener, err = net.Listen("tcp", "127.0.0.1:")
		Expect(err).To(BeNil())
		endpoint = listener.Addr().String()

		serverCert, err = tls.LoadX509KeyPair(
			"./testdata/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt",
			"./testdata/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key",
		)
		Expect(err).NotTo(HaveOccurred())

		connConfig = &client.ConnectionConfig{
			Address:           endpoint,
			ConnectionTimeout: 30 * time.Second,
			TLSEnabled:        true,
			TLSRootCertFile:   "./testdata/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt",
		}
	})

	Describe("CreateGRPCClient", func() {
		AfterEach(func() {
			if listener != nil {
				listener.Close()
			}
		})

		It("creates a useful GRPCClient when TLS is enabled", func() {
			
			grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{serverCert},
			})))
			defer grpcServer.Stop()
			go grpcServer.Serve(listener)

			grpcClient, err := client.CreateGRPCClient(connConfig)
			Expect(err).NotTo(HaveOccurred())

			
			conn, err := grpcClient.NewConnection(endpoint, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(conn).NotTo(BeNil())
		})

		It("creates a useful GRPCClient when TLS is disabled", func() {
			grpcServer := grpc.NewServer()
			defer grpcServer.Stop()
			go grpcServer.Serve(listener)

			connConfig.TLSEnabled = false
			grpcClient, err := client.CreateGRPCClient(connConfig)
			Expect(err).NotTo(HaveOccurred())

			
			conn, err := grpcClient.NewConnection(connConfig.Address, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(conn).NotTo(BeNil())
		})

		Context("when it fails to connect to the endpoint", func() {
			BeforeEach(func() {
				
				connConfig.Address = "127.0.0.1:11111"
				connConfig.ConnectionTimeout = 1 * time.Second
			})

			It("returns an error", func() {
				grpcServer := grpc.NewServer()
				defer grpcServer.Stop()
				go grpcServer.Serve(listener)

				connConfig.TLSEnabled = false
				grpcClient, err := client.CreateGRPCClient(connConfig)
				Expect(err).NotTo(HaveOccurred())

				
				_, err = grpcClient.NewConnection(connConfig.Address, "")
				Expect(err.Error()).To(ContainSubstring("connection refused"))
			})
		})

		Context("when TLS root cert file is missing", func() {
			BeforeEach(func() {
				connConfig.TLSRootCertFile = ""
			})

			It("returns an error", func() {
				_, err := client.CreateGRPCClient(connConfig)
				Expect(err).To(MatchError("missing TLSRootCertFile in client config"))
			})
		})

		Context("when it failed to load root cert file", func() {
			BeforeEach(func() {
				connConfig.TLSRootCertFile = "./testdata/crypto/non-existent-file"
			})

			It("returns an error", func() {
				_, err := client.CreateGRPCClient(connConfig)
				Expect(err.Error()).To(ContainSubstring("unable to load TLS cert from " + connConfig.TLSRootCertFile))
			})
		})
	})
})
