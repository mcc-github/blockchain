/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/



package identity




type Signer interface {
	Sign(message []byte) ([]byte, error)
}




type Serializer interface {
	Serialize() ([]byte, error)
}


type SignerSerializer interface {
	Signer
	Serializer
}
