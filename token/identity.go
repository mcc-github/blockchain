/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package token




type Identity interface {
	Serialize() ([]byte, error)
}






type SigningIdentity interface {
	Identity 

	Sign(msg []byte) ([]byte, error)
}
