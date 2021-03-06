/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package idemix

import (
	"crypto/sha256"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/pkg/errors"
)


type userSecretKey struct {
	
	sk Big
	
	exportable bool
}

func NewUserSecretKey(sk Big, exportable bool) *userSecretKey {
	return &userSecretKey{sk: sk, exportable: exportable}
}

func (k *userSecretKey) Bytes() ([]byte, error) {
	if k.exportable {
		return k.sk.Bytes()
	}

	return nil, errors.New("not exportable")
}

func (k *userSecretKey) SKI() []byte {
	raw, err := k.sk.Bytes()
	if err != nil {
		return nil
	}
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

func (*userSecretKey) Symmetric() bool {
	return true
}

func (*userSecretKey) Private() bool {
	return true
}

func (k *userSecretKey) PublicKey() (bccsp.Key, error) {
	return nil, errors.New("cannot call this method on a symmetric key")
}

type UserKeyGen struct {
	
	
	Exportable bool
	
	User User
}

func (g *UserKeyGen) KeyGen(opts bccsp.KeyGenOpts) (bccsp.Key, error) {
	sk, err := g.User.NewKey()
	if err != nil {
		return nil, err
	}

	return &userSecretKey{exportable: g.Exportable, sk: sk}, nil
}
