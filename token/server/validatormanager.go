/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"github.com/mcc-github/blockchain/token/identity"
	"github.com/mcc-github/blockchain/token/tms/manager"
	"github.com/pkg/errors"
)



type PeerTokenOwnerValidatorManager struct {
	IdentityDeserializerManager identity.DeserializerManager
}

func (p *PeerTokenOwnerValidatorManager) Get(channel string) (identity.TokenOwnerValidator, error) {
	identityDeserializerManager, err := p.IdentityDeserializerManager.Deserializer(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting identity deserialiser manager for channel '%s'", channel)
	}

	return &manager.FabricTokenOwnerValidator{Deserializer: identityDeserializerManager}, nil
}
