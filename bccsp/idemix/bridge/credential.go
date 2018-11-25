/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package bridge

import (
	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-amcl/amcl"
	"github.com/mcc-github/blockchain-amcl/amcl/FP256BN"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
	cryptolib "github.com/mcc-github/blockchain/idemix"
	"github.com/pkg/errors"
)




type Credential struct {
	NewRand func() *amcl.RAND
}





func (c *Credential) Sign(key handlers.IssuerSecretKey, credentialRequest []byte, attributes []bccsp.IdemixAttribute) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = errors.Errorf("failure [%s]", r)
		}
	}()

	iisk, ok := key.(*IssuerSecretKey)
	if !ok {
		return nil, errors.Errorf("invalid issuer secret key, expected *Big, got [%T]", key)
	}

	cr := &cryptolib.CredRequest{}
	err = proto.Unmarshal(credentialRequest, cr)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling credential request")
	}

	attrValues := make([]*FP256BN.BIG, len(attributes))
	for i := 0; i < len(attributes); i++ {
		switch attributes[i].Type {
		case bccsp.IdemixBytesAttribute:
			attrValues[i] = cryptolib.HashModOrder(attributes[i].Value.([]byte))
		case bccsp.IdemixIntAttribute:
			attrValues[i] = FP256BN.NewBIGint(attributes[i].Value.(int))
		default:
			return nil, errors.Errorf("attribute type not allowed or supported [%v] at position [%d]", attributes[i].Type, i)
		}
	}

	cred, err := cryptolib.NewCredential(iisk.SK, cr, attrValues, c.NewRand())
	if err != nil {
		return nil, errors.WithMessage(err, "failed creating new credential")
	}

	return proto.Marshal(cred)
}





func (*Credential) Verify(sk handlers.Big, ipk handlers.IssuerPublicKey, credential []byte, attributes []bccsp.IdemixAttribute) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("failure [%s]", r)
		}
	}()

	isk, ok := sk.(*Big)
	if !ok {
		return errors.Errorf("invalid user secret key, expected *Big, got [%T]", sk)
	}
	iipk, ok := ipk.(*IssuerPublicKey)
	if !ok {
		return errors.Errorf("invalid issuer public key, expected *IssuerPublicKey, got [%T]", sk)
	}

	cred := &cryptolib.Credential{}
	err = proto.Unmarshal(credential, cred)
	if err != nil {
		return err
	}

	for i := 0; i < len(attributes); i++ {
		switch attributes[i].Type {
		case bccsp.IdemixBytesAttribute:
			if !bytes.Equal(
				cryptolib.BigToBytes(cryptolib.HashModOrder(attributes[i].Value.([]byte))),
				cred.Attrs[i]) {
				return errors.Errorf("credential does not contain the correct attribute value at position [%d]", i)
			}
		case bccsp.IdemixIntAttribute:
			if !bytes.Equal(
				cryptolib.BigToBytes(FP256BN.NewBIGint(attributes[i].Value.(int))),
				cred.Attrs[i]) {
				return errors.Errorf("credential does not contain the correct attribute value at position [%d]", i)
			}
		case bccsp.IdemixHiddenAttribute:
			continue
		default:
			return errors.Errorf("attribute type not allowed or supported [%v] at position [%d]", attributes[i].Type, i)
		}
	}

	return cred.Ver(isk.E, iipk.PK)
}
