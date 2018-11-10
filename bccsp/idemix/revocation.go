/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package idemix

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"fmt"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/pkg/errors"
)



type revocationSecretKey struct {
	
	privKey *ecdsa.PrivateKey
	
	exportable bool
}

func NewRevocationSecretKey(sk *ecdsa.PrivateKey, exportable bool) *revocationSecretKey {
	return &revocationSecretKey{privKey: sk, exportable: exportable}
}



func (k *revocationSecretKey) Bytes() ([]byte, error) {
	if k.exportable {
		return k.privKey.D.Bytes(), nil
	}

	return nil, errors.New("not exportable")
}


func (k *revocationSecretKey) SKI() []byte {
	
	raw := elliptic.Marshal(k.privKey.Curve, k.privKey.PublicKey.X, k.privKey.PublicKey.Y)

	
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}



func (k *revocationSecretKey) Symmetric() bool {
	return false
}



func (k *revocationSecretKey) Private() bool {
	return true
}



func (k *revocationSecretKey) PublicKey() (bccsp.Key, error) {
	return &revocationPublicKey{&k.privKey.PublicKey}, nil
}

type revocationPublicKey struct {
	pubKey *ecdsa.PublicKey
}



func (k *revocationPublicKey) Bytes() (raw []byte, err error) {
	raw, err = x509.MarshalPKIXPublicKey(k.pubKey)
	if err != nil {
		return nil, fmt.Errorf("Failed marshalling key [%s]", err)
	}
	return
}


func (k *revocationPublicKey) SKI() []byte {
	
	raw := elliptic.Marshal(k.pubKey.Curve, k.pubKey.X, k.pubKey.Y)

	
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}



func (k *revocationPublicKey) Symmetric() bool {
	return false
}



func (k *revocationPublicKey) Private() bool {
	return false
}



func (k *revocationPublicKey) PublicKey() (bccsp.Key, error) {
	return k, nil
}


type RevocationKeyGen struct {
	
	
	Exportable bool
	
	Revocation Revocation
}

func (g *RevocationKeyGen) KeyGen(opts bccsp.KeyGenOpts) (k bccsp.Key, err error) {
	
	key, err := g.Revocation.NewKey()
	if err != nil {
		return nil, err
	}

	return &revocationSecretKey{exportable: g.Exportable, privKey: key}, nil
}
