/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package idemix


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
