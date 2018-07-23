/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/comm"
	testpb "github.com/mcc-github/blockchain/core/comm/testdata/grpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const testTimeout = 1 * time.Second 

type echoServer struct{}

func (es *echoServer) EchoCall(ctx context.Context,
	echo *testpb.Echo) (*testpb.Echo, error) {
	return echo, nil
}

func TestNewGRPCClient_GoodConfig(t *testing.T) {
	t.Parallel()
	testCerts := loadCerts(t)

	config := comm.ClientConfig{}
	client, err := comm.NewGRPCClient(config)
	assert.NoError(t, err)
	assert.Equal(t, tls.Certificate{}, client.Certificate())
	assert.False(t, client.TLSEnabled())
	assert.False(t, client.MutualTLSRequired())

	secOpts := &comm.SecureOptions{
		UseTLS: false,
	}
	config.SecOpts = secOpts
	client, err = comm.NewGRPCClient(config)
	assert.NoError(t, err)
	assert.Equal(t, tls.Certificate{}, client.Certificate())
	assert.False(t, client.TLSEnabled())
	assert.False(t, client.MutualTLSRequired())

	secOpts = &comm.SecureOptions{
		UseTLS:            true,
		ServerRootCAs:     [][]byte{testCerts.caPEM},
		RequireClientCert: false,
	}
	config.SecOpts = secOpts
	client, err = comm.NewGRPCClient(config)
	assert.NoError(t, err)
	assert.True(t, client.TLSEnabled())
	assert.False(t, client.MutualTLSRequired())

	secOpts = &comm.SecureOptions{
		Certificate:       testCerts.certPEM,
		Key:               testCerts.keyPEM,
		UseTLS:            true,
		ServerRootCAs:     [][]byte{testCerts.caPEM},
		RequireClientCert: true,
	}
	config.SecOpts = secOpts
	client, err = comm.NewGRPCClient(config)
	assert.NoError(t, err)
	assert.True(t, client.TLSEnabled())
	assert.True(t, client.MutualTLSRequired())
	assert.Equal(t, testCerts.clientCert, client.Certificate())
}

func TestNewGRPCClient_BadConfig(t *testing.T) {
	t.Parallel()
	testCerts := loadCerts(t)

	
	config := comm.ClientConfig{
		SecOpts: &comm.SecureOptions{
			UseTLS:        true,
			ServerRootCAs: [][]byte{[]byte(badPEM)},
		},
	}
	_, err := comm.NewGRPCClient(config)
	assert.Contains(t, err.Error(), "error adding root certificate")

	
	missing := "both Key and Certificate are required when using mutual TLS"
	config.SecOpts = &comm.SecureOptions{
		Certificate:       []byte("cert"),
		UseTLS:            true,
		RequireClientCert: true,
	}
	_, err = comm.NewGRPCClient(config)
	assert.Equal(t, missing, err.Error())

	
	config.SecOpts = &comm.SecureOptions{
		Key:               []byte("key"),
		UseTLS:            true,
		RequireClientCert: true,
	}
	_, err = comm.NewGRPCClient(config)
	assert.Equal(t, missing, err.Error())

	
	failed := "failed to load client certificate"
	config.SecOpts = &comm.SecureOptions{
		Certificate:       testCerts.certPEM,
		Key:               []byte(badPEM),
		UseTLS:            true,
		RequireClientCert: true,
	}
	_, err = comm.NewGRPCClient(config)
	assert.Contains(t, err.Error(), failed)

	
	config.SecOpts = &comm.SecureOptions{
		Certificate:       []byte(badPEM),
		Key:               testCerts.keyPEM,
		UseTLS:            true,
		RequireClientCert: true,
	}
	_, err = comm.NewGRPCClient(config)
	assert.Contains(t, err.Error(), failed)
}

func TestNewConnection_Timeout(t *testing.T) {
	t.Parallel()
	testAddress := "localhost:11111"
	config := comm.ClientConfig{
		Timeout: 1 * time.Second,
	}
	client, err := comm.NewGRPCClient(config)
	conn, err := client.NewConnection(testAddress, "")
	assert.Contains(t, err.Error(), "context deadline exceeded")
	t.Log(err)
	assert.Nil(t, conn)
}

func TestNewConnection(t *testing.T) {
	t.Parallel()
	testCerts := loadCerts(t)

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(testCerts.caPEM)
	if !ok {
		t.Fatal("failed to create test root cert pool")
	}

	tests := []struct {
		name       string
		serverPort int
		clientPort int
		config     comm.ClientConfig
		serverTLS  *tls.Config
		success    bool
		errorMsg   string
	}{
		{
			name:       "client / server same port",
			serverPort: 8351,
			clientPort: 8351,
			config: comm.ClientConfig{
				Timeout: testTimeout,
			},
			success: true,
		},
		{
			name:       "client / server wrong port",
			serverPort: 8352,
			clientPort: 7352,
			config: comm.ClientConfig{
				Timeout: time.Second,
			},
			success:  false,
			errorMsg: "context deadline exceeded",
		},
		{
			name:       "client TLS / server no TLS",
			serverPort: 8353,
			clientPort: 8353,
			config: comm.ClientConfig{
				SecOpts: &comm.SecureOptions{
					Certificate:       testCerts.certPEM,
					Key:               testCerts.keyPEM,
					UseTLS:            true,
					ServerRootCAs:     [][]byte{testCerts.caPEM},
					RequireClientCert: true,
				},
				Timeout: time.Second,
			},
			success:  false,
			errorMsg: "context deadline exceeded",
		},
		{
			name:       "client TLS / server TLS match",
			serverPort: 8354,
			clientPort: 8354,
			config: comm.ClientConfig{
				SecOpts: &comm.SecureOptions{
					Certificate:   testCerts.certPEM,
					Key:           testCerts.keyPEM,
					UseTLS:        true,
					ServerRootCAs: [][]byte{testCerts.caPEM},
				},
				Timeout: testTimeout,
			},
			serverTLS: &tls.Config{
				Certificates: []tls.Certificate{testCerts.serverCert},
			},
			success: true,
		},
		{
			name:       "client TLS / server TLS no server roots",
			serverPort: 8355,
			clientPort: 8355,
			config: comm.ClientConfig{
				SecOpts: &comm.SecureOptions{
					Certificate:   testCerts.certPEM,
					Key:           testCerts.keyPEM,
					UseTLS:        true,
					ServerRootCAs: [][]byte{},
				},
				Timeout: testTimeout,
			},
			serverTLS: &tls.Config{
				Certificates: []tls.Certificate{testCerts.serverCert},
			},
			success:  false,
			errorMsg: "context deadline exceeded",
		},
		{
			name:       "client TLS / server TLS missing client cert",
			serverPort: 8356,
			clientPort: 8356,
			config: comm.ClientConfig{
				SecOpts: &comm.SecureOptions{
					Certificate:   testCerts.certPEM,
					Key:           testCerts.keyPEM,
					UseTLS:        true,
					ServerRootCAs: [][]byte{testCerts.caPEM},
				},
				Timeout: testTimeout,
			},
			serverTLS: &tls.Config{
				Certificates: []tls.Certificate{testCerts.serverCert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
			},
			success:  false,
			errorMsg: "context deadline exceeded",
		},
		{
			name:       "client TLS / server TLS client cert",
			serverPort: 8357,
			clientPort: 8357,
			config: comm.ClientConfig{
				SecOpts: &comm.SecureOptions{
					Certificate:       testCerts.certPEM,
					Key:               testCerts.keyPEM,
					UseTLS:            true,
					RequireClientCert: true,
					ServerRootCAs:     [][]byte{testCerts.caPEM}},
				Timeout: testTimeout},
			serverTLS: &tls.Config{
				Certificates: []tls.Certificate{testCerts.serverCert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    certPool,
			},
			success: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", test.serverPort))
			if err != nil {
				t.Fatalf("error creating server for test: %v", err)
			}
			defer lis.Close()
			serverOpts := []grpc.ServerOption{}
			if test.serverTLS != nil {
				serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(test.serverTLS)))
			}
			srv := grpc.NewServer(serverOpts...)
			defer srv.Stop()
			go srv.Serve(lis)
			client, err := comm.NewGRPCClient(test.config)
			if err != nil {
				t.Fatalf("error creating client for test: %v", err)
			}
			conn, err := client.NewConnection(fmt.Sprintf("localhost:%d", test.clientPort), "")
			if test.success {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
			} else {
				t.Log(errors.WithStack(err))
				assert.Contains(t, err.Error(), test.errorMsg)
			}
		})
	}
}

func TestSetServerRootCAs(t *testing.T) {
	t.Parallel()
	testCerts := loadCerts(t)

	config := comm.ClientConfig{
		SecOpts: &comm.SecureOptions{
			UseTLS:        true,
			ServerRootCAs: [][]byte{testCerts.caPEM},
		},
		Timeout: testTimeout,
	}
	client, err := comm.NewGRPCClient(config)
	if err != nil {
		t.Fatalf("error creating base client: %v", err)
	}

	
	address := "localhost:8358"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("failed to create listener for test server: %v", err)
	}
	defer lis.Close()
	srv := grpc.NewServer(grpc.Creds(credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{testCerts.serverCert},
	})))
	defer srv.Stop()
	go srv.Serve(lis)

	
	t.Log("running initial good config")
	conn, err := client.NewConnection(address, "")
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	if conn != nil {
		conn.Close()
	}

	
	t.Log("running bad config")
	err = client.SetServerRootCAs([][]byte{})
	assert.NoError(t, err)
	
	_, err = client.NewConnection(address, "")
	assert.Error(t, err)

	
	t.Log("running good config")
	err = client.SetServerRootCAs([][]byte{[]byte(testCerts.caPEM)})
	assert.NoError(t, err)
	
	conn, err = client.NewConnection(address, "")
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	if conn != nil {
		conn.Close()
	}

	
	t.Log("running bad root cert")
	err = client.SetServerRootCAs([][]byte{[]byte(badPEM)})
	assert.Contains(t, err.Error(), "error adding root certificate")
}

func TestSetMessageSize(t *testing.T) {
	t.Parallel()
	address := "localhost:8359"

	
	srv, err := comm.NewGRPCServer(address, comm.ServerConfig{})
	if err != nil {
		t.Fatalf("failed to create test server: %v", err)
	}
	testpb.RegisterEchoServiceServer(srv.Server(), &echoServer{})
	defer srv.Stop()
	go srv.Start()

	var tests = []struct {
		name        string
		maxRecvSize int
		maxSendSize int
		failRecv    bool
		failSend    bool
	}{
		{
			name:     "defaults should pass",
			failRecv: false,
			failSend: false,
		},
		{
			name:        "non-defaults should pass",
			failRecv:    false,
			failSend:    false,
			maxRecvSize: 20,
			maxSendSize: 20,
		},
		{
			name:        "recv should fail",
			failRecv:    true,
			failSend:    false,
			maxRecvSize: 1,
		},
		{
			name:        "send should fail",
			failRecv:    false,
			failSend:    true,
			maxSendSize: 1,
		},
	}

	
	client, err := comm.NewGRPCClient(comm.ClientConfig{
		Timeout: testTimeout,
	})
	if err != nil {
		t.Fatalf("error creating test client: %v", err)
	}
	
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Log(test.name)
			if test.maxRecvSize > 0 {
				client.SetMaxRecvMsgSize(test.maxRecvSize)
			}
			if test.maxSendSize > 0 {
				client.SetMaxSendMsgSize(test.maxSendSize)
			}
			conn, err := client.NewConnection(address, "")
			assert.NoError(t, err)
			defer conn.Close()
			
			svcClient := testpb.NewEchoServiceClient(conn)
			callCtx := context.Background()
			callCtx, cancel := context.WithTimeout(callCtx, timeout)
			defer cancel()
			
			echo := &testpb.Echo{
				Payload: []byte{0, 0, 0, 0, 0},
			}
			resp, err := svcClient.EchoCall(callCtx, echo)
			if !test.failRecv && !test.failSend {
				assert.NoError(t, err)
				assert.True(t, proto.Equal(echo, resp))
			}
			if test.failSend {
				t.Logf("send error: %v", err)
				assert.Contains(t, err.Error(), "trying to send message larger than max")
			}
			if test.failRecv {
				t.Logf("recv error: %v", err)
				assert.Contains(t, err.Error(), "received message larger than max")
			}
		})
	}
}

type testCerts struct {
	caPEM      []byte
	certPEM    []byte
	keyPEM     []byte
	serverKey  []byte
	serverPEM  []byte
	clientCert tls.Certificate
	serverCert tls.Certificate
}

func loadCerts(t *testing.T) testCerts {
	t.Helper()

	var certs testCerts
	var err error
	certs.caPEM, err = ioutil.ReadFile(filepath.Join("testdata", "certs", "Org1-cert.pem"))
	if err != nil {
		t.Fatalf("unexpected error reading root cert for test: %v", err)
	}
	certs.certPEM, err = ioutil.ReadFile(filepath.Join("testdata", "certs", "Org1-client1-cert.pem"))
	if err != nil {
		t.Fatalf("unexpected error reading cert for test: %v", err)
	}
	certs.keyPEM, err = ioutil.ReadFile(filepath.Join("testdata", "certs", "Org1-client1-key.pem"))
	if err != nil {
		t.Fatalf("unexpected error reading key for test: %v", err)
	}
	certs.clientCert, err = tls.X509KeyPair(certs.certPEM, certs.keyPEM)
	if err != nil {
		t.Fatalf("unexpected error loading certificate for test: %v", err)
	}
	certs.serverCert, err = tls.LoadX509KeyPair(
		filepath.Join("testdata", "certs", "Org1-server1-cert.pem"),
		filepath.Join("testdata", "certs", "Org1-server1-key.pem"),
	)

	return certs
}