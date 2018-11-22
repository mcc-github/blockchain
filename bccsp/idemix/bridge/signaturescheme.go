/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package bridge

import (
	"crypto/ecdsa"

	"github.com/mcc-github/blockchain-amcl/amcl"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-amcl/amcl/FP256BN"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/idemix"
	cryptolib "github.com/mcc-github/blockchain/idemix"
	"github.com/pkg/errors"
)


type SignatureScheme struct {
	NewRand func() *amcl.RAND
}




func (s *SignatureScheme) Sign(cred []byte, sk idemix.Big, Nym idemix.Ecp, RNym idemix.Big, ipk idemix.IssuerPublicKey, attributes []bccsp.IdemixAttribute,
	msg []byte, rhIndex int, criRaw []byte) (res []byte, err error) {
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

	credential := &cryptolib.Credential{}
	err = proto.Unmarshal(cred, credential)
	if err != nil {
		return nil, err
	}

	cri := &cryptolib.CredentialRevocationInformation{}
	err = proto.Unmarshal(criRaw, cri)
	if err != nil {
		return nil, err
	}

	disclosure := make([]byte, len(attributes))
	for i := 0; i < len(attributes); i++ {
		if attributes[i].Type == bccsp.IdemixHiddenAttribute {
			disclosure[i] = 0
		} else {
			disclosure[i] = 1
		}
	}

	sig, err := cryptolib.NewSignature(
		credential,
		isk.E,
		inym.E,
		irnym.E,
		iipk.PK,
		disclosure,
		msg,
		rhIndex,
		cri,
		s.NewRand())
	if err != nil {
		return nil, err
	}

	return proto.Marshal(sig)
}



func (*SignatureScheme) Verify(ipk idemix.IssuerPublicKey, signature, digest []byte, attributes []bccsp.IdemixAttribute, rhIndex int, revocationPublicKey *ecdsa.PublicKey, epoch int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("failure [%s]", r)
		}
	}()

	iipk, ok := ipk.(*IssuerPublicKey)
	if !ok {
		return errors.Errorf("invalid issuer public key, expected *IssuerPublicKey, got [%T]", ipk)
	}

	sig := &cryptolib.Signature{}
	err = proto.Unmarshal(signature, sig)
	if err != nil {
		return err
	}

	disclosure := make([]byte, len(attributes))
	attrValues := make([]*FP256BN.BIG, len(attributes))
	for i := 0; i < len(attributes); i++ {
		switch attributes[i].Type {
		case bccsp.IdemixHiddenAttribute:
			disclosure[i] = 0
			attrValues[i] = nil
		case bccsp.IdemixBytesAttribute:
			disclosure[i] = 1
			attrValues[i] = cryptolib.HashModOrder(attributes[i].Value.([]byte))
		case bccsp.IdemixIntAttribute:
			disclosure[i] = 1
			attrValues[i] = FP256BN.NewBIGint(attributes[i].Value.(int))
		default:
			err = errors.Errorf("attribute type not allowed or supported [%v] at position [%d]", attributes[i].Type, i)
		}
	}
	if err != nil {
		return
	}

	return sig.Ver(
		disclosure,
		iipk.PK,
		digest,
		attrValues,
		rhIndex,
		revocationPublicKey,
		epoch)
}
