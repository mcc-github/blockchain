/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	cb "github.com/mcc-github/blockchain/protos/common"
)


type LocalSigner interface {
	SignatureHeaderMaker
	Signer
}


type Signer interface {
	
	Sign(message []byte) ([]byte, error)
}


type IdentitySerializer interface {
	
	Serialize() ([]byte, error)
}


type SignatureHeaderMaker interface {
	
	NewSignatureHeader() (*cb.SignatureHeader, error)
}


type SignatureHeaderCreator struct {
	SignerSupport
}


type SignerSupport interface {
	Signer
	IdentitySerializer
}


func NewSignatureHeaderCreator(ss SignerSupport) *SignatureHeaderCreator {
	return &SignatureHeaderCreator{ss}
}


func (bs *SignatureHeaderCreator) NewSignatureHeader() (*cb.SignatureHeader, error) {
	creator, err := bs.Serialize()
	if err != nil {
		return nil, err
	}
	nonce, err := GetRandomNonce()
	if err != nil {
		return nil, err
	}

	return &cb.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}, nil
}
