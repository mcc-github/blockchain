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
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"crypto/sha256"

	"errors"

	"encoding/asn1"
	"math/big"

	"github.com/mcc-github/blockchain/bccsp"
)


type rsaPublicKeyASN struct {
	N *big.Int
	E int
}

type rsaPrivateKey struct {
	privKey *rsa.PrivateKey
}



func (k *rsaPrivateKey) Bytes() (raw []byte, err error) {
	return nil, errors.New("Not supported.")
}


func (k *rsaPrivateKey) SKI() (ski []byte) {
	if k.privKey == nil {
		return nil
	}

	
	raw, _ := asn1.Marshal(rsaPublicKeyASN{
		N: k.privKey.N,
		E: k.privKey.E,
	})

	
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}



func (k *rsaPrivateKey) Symmetric() bool {
	return false
}



func (k *rsaPrivateKey) Private() bool {
	return true
}



func (k *rsaPrivateKey) PublicKey() (bccsp.Key, error) {
	return &rsaPublicKey{&k.privKey.PublicKey}, nil
}

type rsaPublicKey struct {
	pubKey *rsa.PublicKey
}



func (k *rsaPublicKey) Bytes() (raw []byte, err error) {
	if k.pubKey == nil {
		return nil, errors.New("Failed marshalling key. Key is nil.")
	}
	raw, err = x509.MarshalPKIXPublicKey(k.pubKey)
	if err != nil {
		return nil, fmt.Errorf("Failed marshalling key [%s]", err)
	}
	return
}


func (k *rsaPublicKey) SKI() (ski []byte) {
	if k.pubKey == nil {
		return nil
	}

	
	raw, _ := asn1.Marshal(rsaPublicKeyASN{
		N: k.pubKey.N,
		E: k.pubKey.E,
	})

	
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}



func (k *rsaPublicKey) Symmetric() bool {
	return false
}



func (k *rsaPublicKey) Private() bool {
	return false
}



func (k *rsaPublicKey) PublicKey() (bccsp.Key, error) {
	return k, nil
}
