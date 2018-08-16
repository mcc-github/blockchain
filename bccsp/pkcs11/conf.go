/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pkcs11

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"fmt"
	"hash"

	"golang.org/x/crypto/sha3"
)

type config struct {
	ellipticCurve asn1.ObjectIdentifier
	hashFunction  func() hash.Hash
	aesBitLength  int
	rsaBitLength  int
}

func (conf *config) setSecurityLevel(securityLevel int, hashFamily string) (err error) {
	switch hashFamily {
	case "SHA2":
		err = conf.setSecurityLevelSHA2(securityLevel)
	case "SHA3":
		err = conf.setSecurityLevelSHA3(securityLevel)
	default:
		err = fmt.Errorf("Hash Family not supported [%s]", hashFamily)
	}
	return
}

func (conf *config) setSecurityLevelSHA2(level int) (err error) {
	switch level {
	case 256:
		conf.ellipticCurve = oidNamedCurveP256
		conf.hashFunction = sha256.New
		conf.rsaBitLength = 2048
		conf.aesBitLength = 32
	case 384:
		conf.ellipticCurve = oidNamedCurveP384
		conf.hashFunction = sha512.New384
		conf.rsaBitLength = 3072
		conf.aesBitLength = 32
	default:
		err = fmt.Errorf("Security level not supported [%d]", level)
	}
	return
}

func (conf *config) setSecurityLevelSHA3(level int) (err error) {
	switch level {
	case 256:
		conf.ellipticCurve = oidNamedCurveP256
		conf.hashFunction = sha3.New256
		conf.rsaBitLength = 2048
		conf.aesBitLength = 32
	case 384:
		conf.ellipticCurve = oidNamedCurveP384
		conf.hashFunction = sha3.New384
		conf.rsaBitLength = 3072
		conf.aesBitLength = 32
	default:
		err = fmt.Errorf("Security level not supported [%d]", level)
	}
	return
}


type PKCS11Opts struct {
	
	SecLevel   int    `mapstructure:"security" json:"security"`
	HashFamily string `mapstructure:"hash" json:"hash"`

	
	Ephemeral     bool               `mapstructure:"tempkeys,omitempty" json:"tempkeys,omitempty"`
	FileKeystore  *FileKeystoreOpts  `mapstructure:"filekeystore,omitempty" json:"filekeystore,omitempty"`
	DummyKeystore *DummyKeystoreOpts `mapstructure:"dummykeystore,omitempty" json:"dummykeystore,omitempty"`

	
	Library    string `mapstructure:"library" json:"library"`
	Label      string `mapstructure:"label" json:"label"`
	Pin        string `mapstructure:"pin" json:"pin"`
	SoftVerify bool   `mapstructure:"softwareverify,omitempty" json:"softwareverify,omitempty"`
}


type FileKeystoreOpts struct {
	KeyStorePath string `mapstructure:"keystore" json:"keystore" yaml:"KeyStore"`
}


type DummyKeystoreOpts struct{}
