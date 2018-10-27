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
package sw

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/mcc-github/blockchain/bccsp"
)

type ecdsaPrivateKey struct {
	privKey *ecdsa.PrivateKey
}



func (k *ecdsaPrivateKey) Bytes() ([]byte, error) {
	return nil, errors.New("Not supported.")
}


func (k *ecdsaPrivateKey) SKI() []byte {
	if k.privKey == nil {
		return nil
	}

	
	raw := elliptic.Marshal(k.privKey.Curve, k.privKey.PublicKey.X, k.privKey.PublicKey.Y)

	
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}



func (k *ecdsaPrivateKey) Symmetric() bool {
	return false
}



func (k *ecdsaPrivateKey) Private() bool {
	return true
}



func (k *ecdsaPrivateKey) PublicKey() (bccsp.Key, error) {
	return &ecdsaPublicKey{&k.privKey.PublicKey}, nil
}

type ecdsaPublicKey struct {
	pubKey *ecdsa.PublicKey
}



func (k *ecdsaPublicKey) Bytes() (raw []byte, err error) {
	raw, err = x509.MarshalPKIXPublicKey(k.pubKey)
	if err != nil {
		return nil, fmt.Errorf("Failed marshalling key [%s]", err)
	}
	return
}


func (k *ecdsaPublicKey) SKI() []byte {
	if k.pubKey == nil {
		return nil
	}

	
	raw := elliptic.Marshal(k.pubKey.Curve, k.pubKey.X, k.pubKey.Y)

	
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
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
