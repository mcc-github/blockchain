/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package manager

import (
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/identity"
	"github.com/pkg/errors"
)


type AllIssuingValidator struct {
	Deserializer identity.Deserializer
}


func (p *AllIssuingValidator) Validate(creator identity.PublicInfo, tokenType string) error {
	
	identity, err := p.Deserializer.DeserializeIdentity(creator.Public())
	if err != nil {
		return errors.Wrapf(err, "identity [0x%x] cannot be deserialised", creator.Public())
	}

	
	if err := identity.Validate(); err != nil {
		return errors.Wrapf(err, "identity [0x%x] cannot be validated", creator.Public())
	}

	return nil
}


type FabricTokenOwnerValidator struct {
	Deserializer identity.Deserializer
}

func (v *FabricTokenOwnerValidator) Validate(owner *token.TokenOwner) error {
	if owner == nil {
		return errors.New("identity cannot be nil")
	}

	switch owner.Type {
	case token.TokenOwner_MSP_IDENTIFIER:
		
		id, err := v.Deserializer.DeserializeIdentity(owner.Raw)
		if err != nil {
			return errors.Wrapf(err, "identity [0x%x] cannot be deserialised", owner)
		}

		
		if err := id.Validate(); err != nil {
			return errors.Wrapf(err, "identity [0x%x] cannot be validated", owner)
		}
	default:
		return errors.Errorf("identity's type '%s' not recognized", owner.Type)
	}

	return nil
}
