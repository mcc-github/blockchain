/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tlsgen

import (
	"crypto"
	"crypto/x509"
)



type CertKeyPair struct {
	
	Cert []byte
	
	Key []byte

	crypto.Signer
	TLSCert *x509.Certificate
}



type CA interface {
	
	CertBytes() []byte

	
	
	
	NewClientCertKeyPair() (*CertKeyPair, error)

	
	
	
	
	NewServerCertKeyPair(host string) (*CertKeyPair, error)
}

type ca struct {
	caCert *CertKeyPair
}

func NewCA() (CA, error) {
	c := &ca{}
	var err error
	c.caCert, err = newCertKeyPair(true, false, "", nil, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}


func (c *ca) CertBytes() []byte {
	return c.caCert.Cert
}




func (c *ca) NewClientCertKeyPair() (*CertKeyPair, error) {
	return newCertKeyPair(false, false, "", c.caCert.Signer, c.caCert.TLSCert)
}




func (c *ca) NewServerCertKeyPair(host string) (*CertKeyPair, error) {
	keypair, err := newCertKeyPair(false, true, host, c.caCert.Signer, c.caCert.TLSCert)
	if err != nil {
		return nil, err
	}
	return keypair, nil
}
