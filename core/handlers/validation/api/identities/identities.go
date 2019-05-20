/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/protos/msp"
)



type IdentityDeserializer interface {
	validation.Dependency
	
	
	
	
	DeserializeIdentity(serializedIdentity []byte) (Identity, error)
}






type Identity interface {
	
	Validate() error

	
	
	
	
	SatisfiesPrincipal(principal *msp.MSPPrincipal) error

	
	Verify(msg []byte, sig []byte) error

	
	GetIdentityIdentifier() *IdentityIdentifier

	
	GetMSPIdentifier() string
}



type IdentityIdentifier struct {

	
	Mspid string

	
	Id string
}
