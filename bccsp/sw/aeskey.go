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
	"errors"

	"crypto/sha256"

	"github.com/mcc-github/blockchain/bccsp"
)

type aesPrivateKey struct {
	privKey    []byte
	exportable bool
}



func (k *aesPrivateKey) Bytes() (raw []byte, err error) {
	if k.exportable {
		return k.privKey, nil
	}

	return nil, errors.New("Not supported.")
}


func (k *aesPrivateKey) SKI() (ski []byte) {
	hash := sha256.New()
	hash.Write([]byte{0x01})
	hash.Write(k.privKey)
	return hash.Sum(nil)
}



func (k *aesPrivateKey) Symmetric() bool {
	return true
}



func (k *aesPrivateKey) Private() bool {
	return true
}



func (k *aesPrivateKey) PublicKey() (bccsp.Key, error) {
	return nil, errors.New("Cannot call this method on a symmetric key.")
}
