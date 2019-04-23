/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"hash"

	"github.com/mcc-github/blockchain/bccsp"
)

type impl struct{}


func New(config map[string]interface{}) (bccsp.BCCSP, error) {
	return &impl{}, nil
}


func (csp *impl) KeyGen(opts bccsp.KeyGenOpts) (k bccsp.Key, err error) {
	return nil, nil
}



func (csp *impl) KeyDeriv(k bccsp.Key, opts bccsp.KeyDerivOpts) (dk bccsp.Key, err error) {
	return nil, nil
}



func (csp *impl) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (k bccsp.Key, err error) {
	return nil, nil
}



func (csp *impl) GetKey(ski []byte) (k bccsp.Key, err error) {
	return nil, nil
}



func (csp *impl) Hash(msg []byte, opts bccsp.HashOpts) (hash []byte, err error) {
	return nil, nil
}



func (csp *impl) GetHash(opts bccsp.HashOpts) (h hash.Hash, err error) {
	return nil, nil
}







func (csp *impl) Sign(k bccsp.Key, digest []byte, opts bccsp.SignerOpts) (signature []byte, err error) {
	return nil, nil
}



func (csp *impl) Verify(k bccsp.Key, signature, digest []byte, opts bccsp.SignerOpts) (valid bool, err error) {
	return true, nil
}



func (csp *impl) Encrypt(k bccsp.Key, plaintext []byte, opts bccsp.EncrypterOpts) (ciphertext []byte, err error) {
	return nil, nil
}



func (csp *impl) Decrypt(k bccsp.Key, ciphertext []byte, opts bccsp.DecrypterOpts) (plaintext []byte, err error) {
	return nil, nil
}
