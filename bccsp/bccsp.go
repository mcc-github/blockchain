/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bccsp

import (
	"crypto"
	"hash"
)


type Key interface {

	
	
	Bytes() ([]byte, error)

	
	SKI() []byte

	
	
	Symmetric() bool

	
	
	Private() bool

	
	
	PublicKey() (Key, error)
}


type KeyGenOpts interface {

	
	Algorithm() string

	
	
	Ephemeral() bool
}


type KeyDerivOpts interface {

	
	Algorithm() string

	
	
	Ephemeral() bool
}


type KeyImportOpts interface {

	
	Algorithm() string

	
	
	Ephemeral() bool
}


type HashOpts interface {

	
	Algorithm() string
}


type SignerOpts interface {
	crypto.SignerOpts
}


type EncrypterOpts interface{}


type DecrypterOpts interface{}



type BCCSP interface {

	
	KeyGen(opts KeyGenOpts) (k Key, err error)

	
	
	KeyDeriv(k Key, opts KeyDerivOpts) (dk Key, err error)

	
	
	KeyImport(raw interface{}, opts KeyImportOpts) (k Key, err error)

	
	
	GetKey(ski []byte) (k Key, err error)

	
	
	Hash(msg []byte, opts HashOpts) (hash []byte, err error)

	
	
	GetHash(opts HashOpts) (h hash.Hash, err error)

	
	
	
	
	
	
	Sign(k Key, digest []byte, opts SignerOpts) (signature []byte, err error)

	
	
	Verify(k Key, signature, digest []byte, opts SignerOpts) (valid bool, err error)

	
	
	Encrypt(k Key, plaintext []byte, opts EncrypterOpts) (ciphertext []byte, err error)

	
	
	Decrypt(k Key, ciphertext []byte, opts DecrypterOpts) (plaintext []byte, err error)
}
