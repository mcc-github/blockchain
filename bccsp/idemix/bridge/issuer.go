/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package bridge

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-amcl/amcl"
	"github.com/mcc-github/blockchain/bccsp/idemix"
	cryptolib "github.com/mcc-github/blockchain/idemix"
	"github.com/pkg/errors"
)


type IssuerPublicKey struct {
	PK *cryptolib.IssuerPublicKey
}

func (o *IssuerPublicKey) Bytes() ([]byte, error) {
	return proto.Marshal(o.PK)
}


type IssuerSecretKey struct {
	SK *cryptolib.IssuerKey
}

func (o *IssuerSecretKey) Bytes() ([]byte, error) {
	return proto.Marshal(o.SK)
}

func (o *IssuerSecretKey) Public() idemix.IssuerPublicKey {
	return &IssuerPublicKey{o.SK.Ipk}
}


type Issuer struct {
	NewRand func() *amcl.RAND
}


func (i *Issuer) NewKey(attributeNames []string) (res idemix.IssuerSecretKey, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = errors.Errorf("failure [%s]", r)
		}
	}()

	sk, err := cryptolib.NewIssuerKey(attributeNames, i.NewRand())
	if err != nil {
		return
	}

	res = &IssuerSecretKey{SK: sk}

	return
}
