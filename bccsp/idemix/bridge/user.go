/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package bridge

import (
	"github.com/mcc-github/blockchain-amcl/amcl"
	"github.com/mcc-github/blockchain/bccsp/idemix"
	cryptolib "github.com/mcc-github/blockchain/idemix"
	"github.com/pkg/errors"
)


type User struct {
	NewRand func() *amcl.RAND
}


func (u *User) NewKey() (res idemix.Big, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = errors.Errorf("failure [%s]", r)
		}
	}()

	res = &Big{E: cryptolib.RandModOrder(u.NewRand())}

	return
}


func (u *User) MakeNym(sk idemix.Big, ipk idemix.IssuerPublicKey) (r1 idemix.Ecp, r2 idemix.Big, err error) {
	defer func() {
		if r := recover(); r != nil {
			r1 = nil
			r2 = nil
			err = errors.Errorf("failure [%s]", r)
		}
	}()

	isk, ok := sk.(*Big)
	if !ok {
		return nil, nil, errors.Errorf("invalid user secret key, expected *Big, got [%T]", sk)
	}
	iipk, ok := ipk.(*IssuerPublicKey)
	if !ok {
		return nil, nil, errors.Errorf("invalid issuer public key, expected *IssuerPublicKey, got [%T]", ipk)
	}

	ecp, big := cryptolib.MakeNym(isk.E, iipk.PK, u.NewRand())

	r1 = &Ecp{E: ecp}
	r2 = &Big{E: big}

	return
}
