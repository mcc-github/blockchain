/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package handlers

import (
	"crypto/sha256"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/pkg/errors"
)



type issuerSecretKey struct {
	
	sk IssuerSecretKey
	
	exportable bool
}

func NewIssuerSecretKey(sk IssuerSecretKey, exportable bool) *issuerSecretKey {
	return &issuerSecretKey{sk: sk, exportable: exportable}
}

func (k *issuerSecretKey) Bytes() ([]byte, error) {
	if k.exportable {
		return k.sk.Bytes()
	}

	return nil, errors.New("not exportable")
}

func (k *issuerSecretKey) SKI() []byte {
	pk, err := k.PublicKey()
	if err != nil {
		return nil
	}

	return pk.SKI()
}

func (*issuerSecretKey) Symmetric() bool {
	return false
}

func (*issuerSecretKey) Private() bool {
	return true
}

func (k *issuerSecretKey) PublicKey() (bccsp.Key, error) {
	return &issuerPublicKey{k.sk.Public()}, nil
}



type issuerPublicKey struct {
	pk IssuerPublicKey
}

func NewIssuerPublicKey(pk IssuerPublicKey) *issuerPublicKey {
	return &issuerPublicKey{pk}
}

func (k *issuerPublicKey) Bytes() ([]byte, error) {
	return k.pk.Bytes()
}

func (k *issuerPublicKey) SKI() []byte {
	raw, err := k.Bytes()
	if err != nil {
		return nil
	}

	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

func (*issuerPublicKey) Symmetric() bool {
	return false
}

func (*issuerPublicKey) Private() bool {
	return false
}

func (k *issuerPublicKey) PublicKey() (bccsp.Key, error) {
	return k, nil
}


type IssuerKeyGen struct {
	
	
	Exportable bool
	
	Issuer Issuer
}

func (g *IssuerKeyGen) KeyGen(opts bccsp.KeyGenOpts) (k bccsp.Key, err error) {
	o, ok := opts.(*bccsp.IdemixIssuerKeyGenOpts)
	if !ok {
		return nil, errors.New("invalid options, expected *bccsp.IdemixIssuerKeyGenOpts")
	}

	
	key, err := g.Issuer.NewKey(o.AttributeNames)
	if err != nil {
		return nil, err
	}

	return &issuerSecretKey{exportable: g.Exportable, sk: key}, nil
}
