/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package entities



type Entity interface {
	
	
	
	
	ID() string

	
	
	
	
	Equals(Entity) bool

	
	
	
	Public() (Entity, error)
}


type Signer interface {
	
	Sign(msg []byte) (signature []byte, err error)

	
	
	Verify(signature, msg []byte) (valid bool, err error)
}


type Encrypter interface {
	
	Encrypt(plaintext []byte) (ciphertext []byte, err error)

	
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
}


type EncrypterEntity interface {
	Entity
	Encrypter
}


type SignerEntity interface {
	Entity
	Signer
}



type EncrypterSignerEntity interface {
	Entity
	Encrypter
	Signer
}
