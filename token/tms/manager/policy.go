/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package manager

import (
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
