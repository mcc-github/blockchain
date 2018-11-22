/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package bridge

import (
	"github.com/mcc-github/blockchain-amcl/amcl"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/idemix"
	cryptolib "github.com/mcc-github/blockchain/idemix"
	"github.com/pkg/errors"
)



type NymSignatureScheme struct {
	NewRand func() *amcl.RAND
}



func (n *NymSignatureScheme) Sign(sk idemix.Big, Nym idemix.Ecp, RNym idemix.Big, ipk idemix.IssuerPublicKey, digest []byte) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = errors.Errorf("failure [%s]", r)
		}
	}()

	isk, ok := sk.(*Big)
	if !ok {
		return nil, errors.Errorf("invalid user secret key, expected *Big, got [%T]", sk)
	}
	inym, ok := Nym.(*Ecp)
	if !ok {
		return nil, errors.Errorf("invalid nym public key, expected *Ecp, got [%T]", Nym)
	}
	irnym, ok := RNym.(*Big)
	if !ok {
		return nil, errors.Errorf("invalid nym secret key, expected *Big, got [%T]", RNym)
	}
	iipk, ok := ipk.(*IssuerPublicKey)
	if !ok {
		return nil, errors.Errorf("invalid issuer public key, expected *IssuerPublicKey, got [%T]", ipk)
	}

	sig, err := cryptolib.NewNymSignature(
		isk.E,
		inym.E,
		irnym.E,
		iipk.PK,
		digest,
		n.NewRand())
	if err != nil {
		return nil, err
	}

	return proto.Marshal(sig)
}



func (*NymSignatureScheme) Verify(ipk idemix.IssuerPublicKey, Nym idemix.Ecp, signature, digest []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("failure [%s]", r)
		}
	}()

	iipk, ok := ipk.(*IssuerPublicKey)
	if !ok {
		return errors.Errorf("invalid issuer public key, expected *IssuerPublicKey, got [%T]", ipk)
	}
	inym, ok := Nym.(*Ecp)
	if !ok {
		return errors.Errorf("invalid nym public key, expected *Ecp, got [%T]", Nym)
	}

	sig := &cryptolib.NymSignature{}
	err = proto.Unmarshal(signature, sig)
	if err != nil {
		return err
	}

	return sig.Ver(inym.E, iipk.PK, digest)
}
