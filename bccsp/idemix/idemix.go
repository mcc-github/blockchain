/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package idemix

import (
	"crypto/ecdsa"

	"github.com/mcc-github/blockchain/bccsp"
)


type IssuerPublicKey interface {

	
	Bytes() ([]byte, error)
}


type IssuerSecretKey interface {

	
	Bytes() ([]byte, error)

	
	Public() IssuerPublicKey
}


type Issuer interface {
	
	NewKey(AttributeNames []string) (IssuerSecretKey, error)
}


type Big interface {
	
	Bytes() ([]byte, error)
}


type Ecp interface {
	
	Bytes() ([]byte, error)
}


type User interface {
	
	NewKey() (Big, error)

	
	MakeNym(sk Big, key IssuerPublicKey) (Ecp, Big, error)
}



type CredRequest interface {
	
	
	Sign(sk Big, ipk IssuerPublicKey) ([]byte, error)

	
	Verify(credRequest []byte, ipk IssuerPublicKey) error
}



type Credential interface {

	
	
	
	Sign(key IssuerSecretKey, credentialRequest []byte, attributes []bccsp.IdemixAttribute) ([]byte, error)

	
	
	Verify(sk Big, ipk IssuerPublicKey, credential []byte, attributes []bccsp.IdemixAttribute) error
}



type Revocation interface {

	
	NewKey() (*ecdsa.PrivateKey, error)

	
	
	
	
	Sign(key *ecdsa.PrivateKey, unrevokedHandles [][]byte, epoch int, alg bccsp.RevocationAlgorithm) ([]byte, error)

	
	
	
	
	
	Verify(pk *ecdsa.PublicKey, cri []byte, epoch int, alg bccsp.RevocationAlgorithm) error
}
