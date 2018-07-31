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
package pkcs11

import (
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/mcc-github/blockchain/bccsp"
)

type ecdsaPrivateKey struct {
	ski []byte
	pub ecdsaPublicKey
}



func (k *ecdsaPrivateKey) Bytes() ([]byte, error) {
	return nil, errors.New("Not supported.")
}


func (k *ecdsaPrivateKey) SKI() []byte {
	return k.ski
}



func (k *ecdsaPrivateKey) Symmetric() bool {
	return false
}



func (k *ecdsaPrivateKey) Private() bool {
	return true
}



func (k *ecdsaPrivateKey) PublicKey() (bccsp.Key, error) {
	return &k.pub, nil
}

type ecdsaPublicKey struct {
	ski []byte
	pub *ecdsa.PublicKey
}



func (k *ecdsaPublicKey) Bytes() (raw []byte, err error) {
	raw, err = x509.MarshalPKIXPublicKey(k.pub)
	if err != nil {
		return nil, fmt.Errorf("Failed marshalling key [%s]", err)
	}
	return
}


func (k *ecdsaPublicKey) SKI() []byte {
	return k.ski
}



func (k *ecdsaPublicKey) Symmetric() bool {
	return false
}



func (k *ecdsaPublicKey) Private() bool {
	return false
}



func (k *ecdsaPublicKey) PublicKey() (bccsp.Key, error) {
	return k, nil
}
