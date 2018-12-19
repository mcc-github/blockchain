/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"bytes"
	"crypto/x509"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)


type RemoteVerifier func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error




type SecureDialer interface {
	Dial(address string, verifyFunc RemoteVerifier) (*grpc.ClientConn, error)
}




type ConnectionMapper interface {
	Lookup(cert []byte) (*grpc.ClientConn, bool)
	Put(cert []byte, conn *grpc.ClientConn)
	Remove(cert []byte)
}


type ConnectionStore struct {
	lock        sync.RWMutex
	Connections ConnectionMapper
	dialer      SecureDialer
}


func NewConnectionStore(dialer SecureDialer) *ConnectionStore {
	connMapping := &ConnectionStore{
		Connections: make(ConnByCertMap),
		dialer:      dialer,
	}
	return connMapping
}



func (c *ConnectionStore) verifyHandshake(endpoint string, certificate []byte) RemoteVerifier {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if bytes.Equal(certificate, rawCerts[0]) {
			return nil
		}
		return errors.Errorf("certificate presented by %s doesn't match any authorized certificate", endpoint)
	}
}


func (c *ConnectionStore) Disconnect(expectedServerCert []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	conn, connected := c.Connections.Lookup(expectedServerCert)
	if !connected {
		return
	}
	conn.Close()
	c.Connections.Remove(expectedServerCert)
}



func (c *ConnectionStore) Connection(endpoint string, expectedServerCert []byte) (*grpc.ClientConn, error) {
	c.lock.RLock()
	conn, alreadyConnected := c.Connections.Lookup(expectedServerCert)
	c.lock.RUnlock()

	if alreadyConnected {
		return conn, nil
	}

	
	return c.connect(endpoint, expectedServerCert)
}



func (c *ConnectionStore) connect(endpoint string, expectedServerCert []byte) (*grpc.ClientConn, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	
	
	conn, alreadyConnected := c.Connections.Lookup(expectedServerCert)
	if alreadyConnected {
		return conn, nil
	}

	v := c.verifyHandshake(endpoint, expectedServerCert)
	conn, err := c.dialer.Dial(endpoint, v)
	if err != nil {
		return nil, err
	}

	c.Connections.Put(expectedServerCert, conn)
	return conn, nil
}
