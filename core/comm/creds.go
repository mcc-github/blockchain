/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"context"
	"crypto/tls"
	"errors"
	"net"

	"github.com/mcc-github/blockchain/common/flogging"
	"google.golang.org/grpc/credentials"
)

var (
	ErrClientHandshakeNotImplemented = errors.New("core/comm: client handshakes are not implemented with serverCreds")
	ErrOverrideHostnameNotSupported  = errors.New("core/comm: OverrideServerName is not supported")

	
	alpnProtoStr = []string{"h2"}
)



func NewServerTransportCredentials(
	serverConfig *tls.Config,
	logger *flogging.FabricLogger) credentials.TransportCredentials {

	
	
	serverConfig.NextProtos = alpnProtoStr
	
	serverConfig.MinVersion = tls.VersionTLS12
	serverConfig.MaxVersion = tls.VersionTLS12
	return &serverCreds{
		serverConfig: serverConfig,
		logger:       logger}
}


type serverCreds struct {
	serverConfig *tls.Config
	logger       *flogging.FabricLogger
}


func (sc *serverCreds) ClientHandshake(context.Context,
	string, net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return nil, nil, ErrClientHandshakeNotImplemented
}


func (sc *serverCreds) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn := tls.Server(rawConn, sc.serverConfig)
	if err := conn.Handshake(); err != nil {
		if sc.logger != nil {
			sc.logger.With("remote address",
				conn.RemoteAddr().String()).Errorf("TLS handshake failed with error %s", err)
		}
		return nil, nil, err
	}
	return conn, credentials.TLSInfo{State: conn.ConnectionState()}, nil
}


func (sc *serverCreds) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
	}
}


func (sc *serverCreds) Clone() credentials.TransportCredentials {
	creds := NewServerTransportCredentials(sc.serverConfig, sc.logger)
	return creds
}



func (sc *serverCreds) OverrideServerName(string) error {
	return ErrOverrideHostnameNotSupported
}
