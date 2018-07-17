/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"io/ioutil"
	"time"

	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/pkg/errors"
)

type genTLSCertFunc func() (*tlsgen.CertKeyPair, error)


type Config struct {
	CertPath       string
	KeyPath        string
	PeerCACertPath string
	Timeout        time.Duration
}




func (conf Config) ToSecureOptions(newSelfSignedTLSCert genTLSCertFunc) (*comm.SecureOptions, error) {
	if conf.PeerCACertPath == "" {
		return &comm.SecureOptions{}, nil
	}
	caBytes, err := loadFile(conf.PeerCACertPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var keyBytes, certBytes []byte
	
	if conf.KeyPath == "" && conf.CertPath == "" {
		tlsCert, err := newSelfSignedTLSCert()
		if err != nil {
			return nil, err
		}
		keyBytes, certBytes = tlsCert.Key, tlsCert.Cert
	} else {
		keyBytes, err = loadFile(conf.KeyPath)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		certBytes, err = loadFile(conf.CertPath)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return &comm.SecureOptions{
		Key:               keyBytes,
		Certificate:       certBytes,
		UseTLS:            true,
		ServerRootCAs:     [][]byte{caBytes},
		RequireClientCert: true,
	}, nil
}

func loadFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf("Failed opening file %s: %v", path, err)
	}
	return b, nil
}
