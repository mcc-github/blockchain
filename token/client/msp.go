/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package client



type Signer interface {
	
	Sign([]byte) ([]byte, error)
}


type SignerIdentity interface {
	Signer

	
	
	Serialize() ([]byte, error)
}
