/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"time"

	"github.com/mcc-github/blockchain/protos/msp"
)


type IdentityDeserializer interface {
	
	
	
	
	DeserializeIdentity(serializedIdentity []byte) (Identity, error)

	
	IsWellFormed(identity *msp.SerializedIdentity) error
}




















type MSPManager interface {

	
	IdentityDeserializer

	
	Setup(msps []MSP) error

	
	GetMSPs() (map[string]MSP, error)
}



type MSP interface {

	
	IdentityDeserializer

	
	Setup(config *msp.MSPConfig) error

	
	GetVersion() MSPVersion

	
	GetType() ProviderType

	
	GetIdentifier() (string, error)

	
	GetSigningIdentity(identifier *IdentityIdentifier) (SigningIdentity, error)

	
	GetDefaultSigningIdentity() (SigningIdentity, error)

	
	GetTLSRootCerts() [][]byte

	
	GetTLSIntermediateCerts() [][]byte

	
	Validate(id Identity) error

	
	
	
	
	SatisfiesPrincipal(id Identity, principal *msp.MSPPrincipal) error
}



type OUIdentifier struct {
	
	
	CertifiersIdentifier []byte
	
	
	OrganizationalUnitIdentifier string
}









type Identity interface {

	
	
	
	
	ExpiresAt() time.Time

	
	GetIdentifier() *IdentityIdentifier

	
	GetMSPIdentifier() string

	
	
	
	
	Validate() error

	
	
	
	
	
	
	
	
	
	
	
	
	
	GetOrganizationalUnits() []*OUIdentifier

	
	Anonymous() bool

	
	Verify(msg []byte, sig []byte) error

	
	Serialize() ([]byte, error)

	
	
	
	
	SatisfiesPrincipal(principal *msp.MSPPrincipal) error
}





type SigningIdentity interface {

	
	Identity

	
	Sign(msg []byte) ([]byte, error)

	
	GetPublicVersion() Identity
}



type IdentityIdentifier struct {

	
	Mspid string

	
	Id string
}


type ProviderType int


const (
	FABRIC ProviderType = iota 
	IDEMIX                     
	OTHER                      

	
	
)

var mspTypeStrings = map[ProviderType]string{
	FABRIC: "bccsp",
	IDEMIX: "idemix",
}

var Options = map[string]NewOpts{
	ProviderTypeToString(FABRIC): &BCCSPNewOpts{NewBaseOpts: NewBaseOpts{Version: MSPv1_0}},
	ProviderTypeToString(IDEMIX): &IdemixNewOpts{NewBaseOpts: NewBaseOpts{Version: MSPv1_1}},
}


func ProviderTypeToString(id ProviderType) string {
	if res, found := mspTypeStrings[id]; found {
		return res
	}

	return ""
}
