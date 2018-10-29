/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package identity

import (
	"github.com/mcc-github/blockchain/msp"
)


type IssuingValidator interface {
	
	Validate(creator PublicInfo, tokenType string) error
}


type PublicInfo interface {
	Public() []byte
}


type DeserializerManager interface {
	
	
	Deserializer(channel string) (Deserializer, error)
}


type Deserializer interface {
	
	
	
	
	DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error)
}
