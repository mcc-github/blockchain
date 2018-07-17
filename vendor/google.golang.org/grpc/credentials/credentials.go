





package credentials 

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"golang.org/x/net/context"
)


var alpnProtoStr = []string{"h2"}



type PerRPCCredentials interface {
	
	
	
	
	
	
	
	
	
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	
	
	RequireTransportSecurity() bool
}



type ProtocolInfo struct {
	
	ProtocolVersion string
	
	SecurityProtocol string
	
	SecurityVersion string
	
	ServerName string
}


type AuthInfo interface {
	AuthType() string
}



var ErrConnDispatched = errors.New("credentials: rawConn is dispatched out of gRPC")



type TransportCredentials interface {
	
	
	
	
	
	
	
	
	
	
	ClientHandshake(context.Context, string, net.Conn) (net.Conn, AuthInfo, error)
	
	
	
	
	
	ServerHandshake(net.Conn) (net.Conn, AuthInfo, error)
	
	Info() ProtocolInfo
	
	Clone() TransportCredentials
	
	
	
	OverrideServerName(string) error
}



type TLSInfo struct {
	State tls.ConnectionState
}


func (t TLSInfo) AuthType() string {
	return "tls"
}


type tlsCreds struct {
	
	config *tls.Config
}

func (c tlsCreds) Info() ProtocolInfo {
	return ProtocolInfo{
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
		ServerName:       c.config.ServerName,
	}
}

func (c *tlsCreds) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (_ net.Conn, _ AuthInfo, err error) {
	
	cfg := cloneTLSConfig(c.config)
	if cfg.ServerName == "" {
		colonPos := strings.LastIndex(authority, ":")
		if colonPos == -1 {
			colonPos = len(authority)
		}
		cfg.ServerName = authority[:colonPos]
	}
	conn := tls.Client(rawConn, cfg)
	errChannel := make(chan error, 1)
	go func() {
		errChannel <- conn.Handshake()
	}()
	select {
	case err := <-errChannel:
		if err != nil {
			return nil, nil, err
		}
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	return conn, TLSInfo{conn.ConnectionState()}, nil
}

func (c *tlsCreds) ServerHandshake(rawConn net.Conn) (net.Conn, AuthInfo, error) {
	conn := tls.Server(rawConn, c.config)
	if err := conn.Handshake(); err != nil {
		return nil, nil, err
	}
	return conn, TLSInfo{conn.ConnectionState()}, nil
}

func (c *tlsCreds) Clone() TransportCredentials {
	return NewTLS(c.config)
}

func (c *tlsCreds) OverrideServerName(serverNameOverride string) error {
	c.config.ServerName = serverNameOverride
	return nil
}


func NewTLS(c *tls.Config) TransportCredentials {
	tc := &tlsCreds{cloneTLSConfig(c)}
	tc.config.NextProtos = alpnProtoStr
	return tc
}




func NewClientTLSFromCert(cp *x509.CertPool, serverNameOverride string) TransportCredentials {
	return NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp})
}




func NewClientTLSFromFile(certFile, serverNameOverride string) (TransportCredentials, error) {
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp}), nil
}


func NewServerTLSFromCert(cert *tls.Certificate) TransportCredentials {
	return NewTLS(&tls.Config{Certificates: []tls.Certificate{*cert}})
}



func NewServerTLSFromFile(certFile, keyFile string) (TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}}), nil
}
