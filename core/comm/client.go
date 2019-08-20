/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GRPCClient struct {
	
	tlsConfig *tls.Config
	
	dialOpts []grpc.DialOption
	
	timeout time.Duration
	
	maxRecvMsgSize int
	
	maxSendMsgSize int
}



func NewGRPCClient(config ClientConfig) (*GRPCClient, error) {
	client := &GRPCClient{}

	
	err := client.parseSecureOptions(config.SecOpts)
	if err != nil {
		return client, err
	}

	

	kap := keepalive.ClientParameters{
		Time:                config.KaOpts.ClientInterval,
		Timeout:             config.KaOpts.ClientTimeout,
		PermitWithoutStream: true,
	}
	
	client.dialOpts = append(client.dialOpts, grpc.WithKeepaliveParams(kap))
	
	if !config.AsyncConnect {
		client.dialOpts = append(client.dialOpts, grpc.WithBlock())
		client.dialOpts = append(client.dialOpts, grpc.FailOnNonTempDialError(true))
	}
	client.timeout = config.Timeout
	
	client.maxRecvMsgSize = MaxRecvMsgSize
	client.maxSendMsgSize = MaxSendMsgSize

	return client, nil
}

func (client *GRPCClient) parseSecureOptions(opts SecureOptions) error {
	
	if !opts.UseTLS {
		return nil
	}

	client.tlsConfig = &tls.Config{
		VerifyPeerCertificate: opts.VerifyCertificate,
		MinVersion:            tls.VersionTLS12} 
	if len(opts.ServerRootCAs) > 0 {
		client.tlsConfig.RootCAs = x509.NewCertPool()
		for _, certBytes := range opts.ServerRootCAs {
			err := AddPemToCertPool(certBytes, client.tlsConfig.RootCAs)
			if err != nil {
				commLogger.Debugf("error adding root certificate: %v", err)
				return errors.WithMessage(err,
					"error adding root certificate")
			}
		}
	}
	if opts.RequireClientCert {
		
		if opts.Key != nil &&
			opts.Certificate != nil {
			cert, err := tls.X509KeyPair(opts.Certificate,
				opts.Key)
			if err != nil {
				return errors.WithMessage(err, "failed to "+
					"load client certificate")
			}
			client.tlsConfig.Certificates = append(
				client.tlsConfig.Certificates, cert)
		} else {
			return errors.New("both Key and Certificate " +
				"are required when using mutual TLS")
		}
	}

	if opts.TimeShift > 0 {
		client.tlsConfig.Time = func() time.Time {
			return time.Now().Add((-1) * opts.TimeShift)
		}
	}

	return nil
}



func (client *GRPCClient) Certificate() tls.Certificate {
	cert := tls.Certificate{}
	if client.tlsConfig != nil && len(client.tlsConfig.Certificates) > 0 {
		cert = client.tlsConfig.Certificates[0]
	}
	return cert
}



func (client *GRPCClient) TLSEnabled() bool {
	return client.tlsConfig != nil
}



func (client *GRPCClient) MutualTLSRequired() bool {
	return client.tlsConfig != nil &&
		len(client.tlsConfig.Certificates) > 0
}


func (client *GRPCClient) SetMaxRecvMsgSize(size int) {
	client.maxRecvMsgSize = size
}


func (client *GRPCClient) SetMaxSendMsgSize(size int) {
	client.maxSendMsgSize = size
}



func (client *GRPCClient) SetServerRootCAs(serverRoots [][]byte) error {

	
	
	certPool := x509.NewCertPool()
	for _, root := range serverRoots {
		err := AddPemToCertPool(root, certPool)
		if err != nil {
			return errors.WithMessage(err, "error adding root certificate")
		}
	}
	client.tlsConfig.RootCAs = certPool
	return nil
}




func (client *GRPCClient) NewConnection(address string, serverNameOverride string) (
	*grpc.ClientConn, error) {

	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, client.dialOpts...)

	
	
	
	
	if client.tlsConfig != nil {
		client.tlsConfig.ServerName = serverNameOverride
		dialOpts = append(dialOpts,
			grpc.WithTransportCredentials(
				credentials.NewTLS(client.tlsConfig)))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(client.maxRecvMsgSize),
		grpc.MaxCallSendMsgSize(client.maxSendMsgSize)))

	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, dialOpts...)
	if err != nil {
		return nil, errors.WithMessage(errors.WithStack(err),
			"failed to create new connection")
	}
	return conn, nil
}
